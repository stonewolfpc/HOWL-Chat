package audio

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"howl-chat/internal/audio/types"
)

// DesiredSettings is the complete persisted audio settings object. It is the
// request, not the effective configuration: most fields are preserved and
// reported but are not applied by the current backend.
type DesiredSettings struct {
	ASREnabled bool `json:"asrEnabled"`

	TTSEnabled bool `json:"ttsEnabled"`

	ASREngine string `json:"asrEngine"`

	ASRLanguage           string  `json:"asrLanguage"`
	ASRDevice             string  `json:"asrDevice"`
	ASRThreads            int     `json:"asrThreads"`
	ASRBeamSize           int     `json:"asrBeamSize"`
	ASRBestOf             int     `json:"asrBestOf"`
	ASRTemperature        float64 `json:"asrTemperature"`
	ASRTranslateToEnglish bool    `json:"asrTranslateToEnglish"`
	ASRNoTimestamps       bool    `json:"asrNoTimestamps"`
	ASRFlashAttention     bool    `json:"asrFlashAttention"`

	VADEnabled      bool    `json:"vadEnabled"`
	VADThreshold    float64 `json:"vadThreshold"`
	VADMinSpeechMs  int     `json:"vadMinSpeechMs"`
	VADMinSilenceMs int     `json:"vadMinSilenceMs"`
	VADPaddingMs    int     `json:"vadPaddingMs"`

	TTSEngine string `json:"ttsEngine"`

	TTSVoice           string  `json:"ttsVoice"`
	TTSSpeed           float64 `json:"ttsSpeed"`
	TTSVolume          float64 `json:"ttsVolume"`
	TTSNoiseScale      float64 `json:"ttsNoiseScale"`
	TTSNoiseW          float64 `json:"ttsNoiseW"`
	TTSSentenceSilence float64 `json:"ttsSentenceSilence"`
	TTSAutoplayMessage bool    `json:"ttsAutoplayMessage"`
	TTSCacheEnabled    bool    `json:"ttsCacheEnabled"`
	TTSCacheSizeMB     int     `json:"ttsCacheSizeMB"`

	AudioInputDevice  string  `json:"audioInputDevice"`
	AudioOutputDevice string  `json:"audioOutputDevice"`
	AudioSampleRate   int     `json:"audioSampleRate"`
	AudioChannels     string  `json:"audioChannels"`
	AudioBitDepth     int     `json:"audioBitDepth"`
	MicSensitivity    float64 `json:"micSensitivity"`

	PushToTalk bool `json:"pushToTalk"`
}

// DefaultDesiredSettings mirrors the renderer defaults. It is used when a
// stored settings object predates a key, not as a fallback for an invalid
// requested value.
func DefaultDesiredSettings() DesiredSettings {
	return DesiredSettings{
		ASREnabled:            false,
		TTSEnabled:            false,
		ASREngine:             "auto",
		ASRLanguage:           "auto",
		ASRDevice:             "auto",
		ASRThreads:            4,
		ASRBeamSize:           5,
		ASRBestOf:             5,
		ASRTemperature:        0,
		ASRTranslateToEnglish: false,
		ASRNoTimestamps:       false,
		ASRFlashAttention:     true,
		VADEnabled:            false,
		VADThreshold:          0.5,
		VADMinSpeechMs:        250,
		VADMinSilenceMs:       100,
		VADPaddingMs:          30,
		TTSEngine:             "piper",
		TTSVoice:              "en_US-lessac-medium",
		TTSSpeed:              1,
		TTSVolume:             0.8,
		TTSNoiseScale:         0.667,
		TTSNoiseW:             0.8,
		TTSSentenceSilence:    0.2,
		TTSAutoplayMessage:    true,
		TTSCacheEnabled:       true,
		TTSCacheSizeMB:        256,
		AudioInputDevice:      "default",
		AudioOutputDevice:     "default",
		AudioSampleRate:       16000,
		AudioChannels:         "mono",
		AudioBitDepth:         16,
		MicSensitivity:        0.5,
		PushToTalk:            true,
	}
}

// DesiredSettingsFromMap overlays only keys present in raw onto the defaults.
// Absent keys therefore retain the same defaults the renderer uses.
func DesiredSettingsFromMap(raw map[string]interface{}) DesiredSettings {
	out := DefaultDesiredSettings()
	if len(raw) == 0 {
		return out
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return out
	}
	_ = json.Unmarshal(data, &out)
	return out
}

// Backend returns the operational settings that the audio manager can apply
// today. The other desired fields are intentionally excluded until they are
// wired and validated by their owning component.
func (d DesiredSettings) Backend() Settings {
	return Settings{
		ASREnabled: d.ASREnabled,
		TTSEnabled: d.TTSEnabled,
		ASREngine:  d.ASREngine,
		TTSEngine:  d.TTSEngine,
		TTSVoice:   d.TTSVoice,
		TTSSpeed:   d.TTSSpeed,

		ASRLanguage: d.ASRLanguage,
		ASRThreads:  d.ASRThreads,
		ASRBeamSize: d.ASRBeamSize,
		ASRBestOf:   d.ASRBestOf,
	}
}

// Validate returns per-key reasons for malformed values. It does not decide
// whether a setting is implemented; implementation status is separate.
func (d DesiredSettings) Validate() map[string]string {
	issues := make(map[string]string)

	if d.ASRLanguage != "" && !oneOf(d.ASRLanguage, "auto", "en", "es", "fr", "de", "it", "pt", "ja", "zh", "ko", "ru", "nl", "pl") {
		issues["asrLanguage"] = fmt.Sprintf("unsupported spoken language %q", d.ASRLanguage)
	} else if d.ASREnabled && isWhisperASREngine(d.ASREngine) && d.ASRLanguage != "" && d.ASRLanguage != "auto" && !asrLanguageSupportedByEngine(d.ASREngine, d.ASRLanguage) {
		issues["asrLanguage"] = fmt.Sprintf("ASR engine %q does not support spoken language %q", d.ASREngine, d.ASRLanguage)
	}
	if d.ASRDevice != "" && !oneOf(d.ASRDevice, "auto", "cpu", "cuda") {
		issues["asrDevice"] = fmt.Sprintf("unsupported compute device %q", d.ASRDevice)
	}
	if d.ASRThreads < 1 || d.ASRThreads > 16 {
		issues["asrThreads"] = fmt.Sprintf("CPU threads must be between 1 and 16, got %d", d.ASRThreads)
	}
	if d.ASRBeamSize < 1 || d.ASRBeamSize > 10 {
		issues["asrBeamSize"] = fmt.Sprintf("beam size must be between 1 and 10, got %d", d.ASRBeamSize)
	}
	if d.ASRBestOf < 1 || d.ASRBestOf > 10 {
		issues["asrBestOf"] = fmt.Sprintf("best-of must be between 1 and 10, got %d", d.ASRBestOf)
	}
	if d.ASRTemperature < 0 || d.ASRTemperature > 1 {
		issues["asrTemperature"] = fmt.Sprintf("ASR temperature must be between 0 and 1, got %g", d.ASRTemperature)
	}
	if d.ASREnabled {
		engine := strings.TrimSpace(d.ASREngine)
		if engine != "" && engine != "auto" {
			if _, ok := types.GetASRProfile(types.ASRModelType(engine)); !ok {
				issues["asrEngine"] = fmt.Sprintf("unknown ASR engine %q", d.ASREngine)
			}
		}
	}

	if d.VADThreshold < 0 || d.VADThreshold > 1 {
		issues["vadThreshold"] = fmt.Sprintf("VAD threshold must be between 0 and 1, got %g", d.VADThreshold)
	}
	if d.VADMinSpeechMs < 0 || d.VADMinSpeechMs > 1000 {
		issues["vadMinSpeechMs"] = fmt.Sprintf("minimum speech must be between 0 and 1000 ms, got %d", d.VADMinSpeechMs)
	}
	if d.VADMinSilenceMs < 0 || d.VADMinSilenceMs > 1000 {
		issues["vadMinSilenceMs"] = fmt.Sprintf("minimum silence must be between 0 and 1000 ms, got %d", d.VADMinSilenceMs)
	}
	if d.VADPaddingMs < 0 || d.VADPaddingMs > 500 {
		issues["vadPaddingMs"] = fmt.Sprintf("speech padding must be between 0 and 500 ms, got %d", d.VADPaddingMs)
	}

	if d.TTSSpeed < 0.5 || d.TTSSpeed > 2 {
		issues["ttsSpeed"] = fmt.Sprintf("speech speed must be between 0.5 and 2.0, got %g", d.TTSSpeed)
	} else if d.TTSEnabled {
		if min, max, supported, explicit := ttsSpeedRangeFor(d.TTSEngine); explicit && supported && (d.TTSSpeed < min || d.TTSSpeed > max) {
			issues["ttsSpeed"] = fmt.Sprintf("speech speed for %s must be between %g and %g, got %g", strings.TrimSpace(d.TTSEngine), min, max, d.TTSSpeed)
		}
	}
	if d.TTSVolume < 0 || d.TTSVolume > 1 {
		issues["ttsVolume"] = fmt.Sprintf("TTS volume must be between 0 and 1, got %g", d.TTSVolume)
	}
	if d.TTSNoiseScale < 0 || d.TTSNoiseScale > 1 {
		issues["ttsNoiseScale"] = fmt.Sprintf("noise scale must be between 0 and 1, got %g", d.TTSNoiseScale)
	}
	if d.TTSNoiseW < 0 || d.TTSNoiseW > 1 {
		issues["ttsNoiseW"] = fmt.Sprintf("noise width must be between 0 and 1, got %g", d.TTSNoiseW)
	}
	if d.TTSSentenceSilence < 0 || d.TTSSentenceSilence > 2 {
		issues["ttsSentenceSilence"] = fmt.Sprintf("sentence silence must be between 0 and 2 seconds, got %g", d.TTSSentenceSilence)
	}
	if d.TTSCacheSizeMB < 50 || d.TTSCacheSizeMB > 1024 {
		issues["ttsCacheSizeMB"] = fmt.Sprintf("cache size must be between 50 and 1024 MB, got %d", d.TTSCacheSizeMB)
	}

	if d.AudioSampleRate != 16000 && d.AudioSampleRate != 22050 && d.AudioSampleRate != 44100 && d.AudioSampleRate != 48000 {
		issues["audioSampleRate"] = fmt.Sprintf("unsupported sample rate %d", d.AudioSampleRate)
	}
	if d.AudioChannels != "" && d.AudioChannels != "mono" && d.AudioChannels != "stereo" {
		issues["audioChannels"] = fmt.Sprintf("unsupported channel layout %q", d.AudioChannels)
	}
	if d.AudioBitDepth != 16 && d.AudioBitDepth != 24 && d.AudioBitDepth != 32 {
		issues["audioBitDepth"] = fmt.Sprintf("unsupported bit depth %d", d.AudioBitDepth)
	}
	if d.MicSensitivity < 0 || d.MicSensitivity > 1 {
		issues["micSensitivity"] = fmt.Sprintf("microphone sensitivity must be between 0 and 1, got %g", d.MicSensitivity)
	}

	return issues
}

func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}

// BackendSettings are the only settings with an actual Go-side application
// path today.
var backendSettingKeys = map[string]bool{
	"asrEnabled":  true,
	"ttsEnabled":  true,
	"asrEngine":   true,
	"ttsEngine":   true,
	"ttsVoice":    true,
	"asrLanguage": true,
	"asrThreads":  true,
	"asrBeamSize": true,
	"asrBestOf":   true,
	"ttsSpeed":    true,
}

var asrDecodeSettingKeys = []string{"asrLanguage", "asrThreads", "asrBeamSize", "asrBestOf"}

// HasBackendValidationIssue reports whether an invalid value belongs to a
// setting the backend would actually apply. Invalid unsupported fields are
// preserved and reported, but they must not prevent the supported operational
// fields from being applied. The four Whisper decode values become blocking
// only when the request names an explicit Whisper engine.
func HasBackendValidationIssue(desired DesiredSettings, issues map[string]string) bool {
	for key := range issues {
		if key == "ttsSpeed" && !desired.TTSEnabled {
			// A saved speed is not operational while TTS is off.
			continue
		}
		if isASRDecodeSetting(key) {
			continue
		}
		if backendSettingKeys[key] {
			return true
		}
	}
	if desired.ASREnabled && isWhisperASREngine(desired.ASREngine) {
		for _, key := range asrDecodeSettingKeys {
			if issues[key] != "" {
				return true
			}
		}
	}
	return false
}

func isASRDecodeSetting(key string) bool {
	for _, candidate := range asrDecodeSettingKeys {
		if key == candidate {
			return true
		}
	}
	return false
}

// isWhisperASREngine identifies an explicit engine whose adapter emits the
// decode flags this slice owns. Automatic selection is deliberately excluded:
// it does not identify the provider that will actually decode.
func isWhisperASREngine(engine string) bool {
	profile, ok := types.GetASRProfile(types.ASRModelType(strings.TrimSpace(engine)))
	return ok && profile.Local && profile.Runtime == types.RuntimeWhisperCLI
}

// ttsSpeedRangeFor identifies a TTS engine whose adapter accepts a speech-rate
// multiplier. Automatic selection deliberately returns explicit=false: the
// effective provider is unknown until initialization, so speed is not applied.
// The range comes from the same model profile the adapter is built around;
// MeloTTS's 0.8-1.2 is a product guard, not an upstream guarantee.
func ttsSpeedRangeFor(engine string) (min, max float64, supported, explicit bool) {
	engine = strings.TrimSpace(engine)
	if engine == "" || engine == "auto" {
		return 0, 0, false, false
	}
	profile, ok := types.GetTTSProfile(types.TTSModelType(engine))
	if !ok {
		return 0, 0, false, true
	}
	if !profile.SupportsSpeed || profile.MinSpeed <= 0 || profile.MaxSpeed <= 0 {
		return 0, 0, false, true
	}
	return profile.MinSpeed, profile.MaxSpeed, true, true
}

func ttsSpeedOrDefault(speed float64) float64 {
	if speed <= 0 {
		return 1.0
	}
	return speed
}

func effectiveTTSSpeedFor(settings Settings) (float64, bool) {
	if !settings.TTSEnabled {
		return 0, false
	}
	min, max, supported, explicit := ttsSpeedRangeFor(settings.TTSEngine)
	if !explicit || !supported {
		return 0, false
	}
	speed := ttsSpeedOrDefault(settings.TTSSpeed)
	if speed < min || speed > max {
		return 0, false
	}
	return speed, true
}

func asrLanguageSupportedByEngine(engine, language string) bool {
	profile, ok := types.GetASRProfile(types.ASRModelType(strings.TrimSpace(engine)))
	if !ok {
		return false
	}
	for _, supported := range profile.Languages {
		if supported == language {
			return true
		}
	}
	return false
}

// SettingsValidationError reports malformed requested values. It never
// substitutes a different value.
type SettingsValidationError struct {
	Issues map[string]string
}

func (e *SettingsValidationError) Error() string {
	if e == nil || len(e.Issues) == 0 {
		return "audio settings validation failed"
	}
	keys := make([]string, 0, len(e.Issues))
	for key := range e.Issues {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, key+": "+e.Issues[key])
	}
	return "invalid audio settings: " + strings.Join(parts, "; ")
}

type SettingOwner string

const (
	SettingOwnerBackend  SettingOwner = "backend"
	SettingOwnerRenderer SettingOwner = "renderer"
	SettingOwnerNone     SettingOwner = "none"
)

type SettingState string

const (
	SettingApplied     SettingState = "applied"
	SettingPending     SettingState = "pending"
	SettingRejected    SettingState = "rejected"
	SettingInvalid     SettingState = "invalid"
	SettingUnsupported SettingState = "unsupported"
	SettingUnavailable SettingState = "unavailable"
	// SettingRenderer marks a value owned and applied by the renderer. The Go
	// backend can preserve and report it, but must not claim it applied it.
	SettingRenderer SettingState = "renderer"
)

// SettingStatus is one setting's truthful lifecycle state. Applied is set only
// when that exact value is effective in its owning component.
type SettingStatus struct {
	Key     string       `json:"key"`
	Owner   SettingOwner `json:"owner"`
	State   SettingState `json:"state"`
	Applied interface{}  `json:"applied,omitempty"`
	Reason  string       `json:"reason,omitempty"`
}

// SettingsCounts is a compact summary of the per-setting states.
type SettingsCounts struct {
	Total       int `json:"total"`
	Applied     int `json:"applied"`
	Pending     int `json:"pending"`
	Rejected    int `json:"rejected"`
	Invalid     int `json:"invalid"`
	Unsupported int `json:"unsupported"`
	Unavailable int `json:"unavailable"`
	Renderer    int `json:"renderer"`
}

// SettingsReport is the authoritative contract returned with audio status.
// Applied contains backend-effective values only. Renderer-owned values are
// marked by owner/state and are not included in Applied.
type SettingsReport struct {
	Desired  DesiredSettings          `json:"desired"`
	Applied  map[string]interface{}   `json:"applied"`
	Statuses map[string]SettingStatus `json:"statuses"`
	Counts   SettingsCounts           `json:"counts"`
}

// BuildSettingsReport derives the per-setting lifecycle from the desired
// request and the manager's actual active generation. It never promotes an
// unsupported field into Applied.
func BuildSettingsReport(desired DesiredSettings, applied Settings, appliedOK bool, status Status, applyErr string) SettingsReport {
	issues := desired.Validate()
	report := SettingsReport{
		Desired:  desired,
		Applied:  make(map[string]interface{}),
		Statuses: make(map[string]SettingStatus),
	}

	for key, reason := range unsupportedSettings() {
		report.Statuses[key] = SettingStatus{
			Key:    key,
			Owner:  SettingOwnerNone,
			State:  SettingUnsupported,
			Reason: reason,
		}
	}
	for key, reason := range rendererSettings() {
		report.Statuses[key] = SettingStatus{
			Key:    key,
			Owner:  SettingOwnerRenderer,
			State:  SettingRenderer,
			Reason: reason,
		}
	}
	for key := range backendSettingKeys {
		report.Statuses[key] = SettingStatus{
			Key:   key,
			Owner: SettingOwnerBackend,
			State: SettingPending,
		}
	}

	setBool := func(key string, want, got bool) {
		if issues[key] != "" {
			return
		}
		st := report.Statuses[key]
		if appliedOK {
			// Applied is the last effective configuration, even when the newest
			// request differs. The status below says whether this field still
			// matches what the user asked for.
			report.Applied[key] = got
		}
		switch {
		case appliedOK && want == got:
			st.State = SettingApplied
			st.Applied = want
		case !appliedOK && !want:
			// No active generation and the capability is not requested: the
			// disabled value is the effective state.
			st.State = SettingApplied
			st.Applied = false
			report.Applied[key] = false
		case appliedOK:
			st.State, st.Reason = notEffectiveState(status, applyErr)
		default:
			st.State, st.Reason = notEffectiveState(status, applyErr)
		}
		report.Statuses[key] = st
	}

	setEngine := func(key string, enabled, appliedEnabled bool, want, got string, ready bool, capErr string) {
		st := report.Statuses[key]
		if appliedOK && appliedEnabled {
			// Preserve the active generation's engine even after a failed
			// replacement or a failed disable so Applied never loses it.
			report.Applied[key] = got
		}
		if issues[key] != "" {
			return
		}
		if !enabled {
			if appliedOK && appliedEnabled {
				st.State, st.Reason = notEffectiveState(status, applyErr)
			} else {
				st.State = SettingUnavailable
				st.Reason = "the capability is disabled; this selection is saved but not loaded"
			}
			report.Statuses[key] = st
			return
		}
		switch {
		case !appliedOK:
			st.State, st.Reason = notEffectiveState(status, applyErr)
		case strings.TrimSpace(got) != strings.TrimSpace(want):
			st.State, st.Reason = notEffectiveState(status, applyErr)
		case !ready:
			st.State = SettingUnavailable
			st.Reason = capErr
			if st.Reason == "" {
				st.Reason = "the capability is applied but its provider is not initialized"
			}
		default:
			st.State = SettingApplied
			st.Applied = want
			report.Applied[key] = want
		}
		report.Statuses[key] = st
	}

	setBool("asrEnabled", desired.ASREnabled, applied.ASREnabled)
	setBool("ttsEnabled", desired.TTSEnabled, applied.TTSEnabled)
	setEngine("asrEngine", desired.ASREnabled, applied.ASREnabled, desired.ASREngine, applied.ASREngine, status.ASR.Ready, status.ASR.Error)
	setEngine("ttsEngine", desired.TTSEnabled, applied.TTSEnabled, desired.TTSEngine, applied.TTSEngine, status.TTS.Ready, status.TTS.Error)

	setTTSVoice := func() {
		key := "ttsVoice"
		st := report.Statuses[key]
		if appliedOK && applied.TTSEnabled {
			report.Applied[key] = applied.TTSVoice
		}
		if !desired.TTSEnabled {
			st.State = SettingUnavailable
			st.Reason = "TTS is disabled; this value is saved but not applied"
			report.Statuses[key] = st
			return
		}
		if !appliedOK || !applied.TTSEnabled {
			st.State, st.Reason = notEffectiveState(status, applyErr)
			report.Statuses[key] = st
			return
		}
		if strings.TrimSpace(applied.TTSVoice) != strings.TrimSpace(desired.TTSVoice) {
			st.State, st.Reason = notEffectiveState(status, applyErr)
			report.Statuses[key] = st
			return
		}
		if !status.TTS.Ready {
			st.State = SettingUnavailable
			st.Reason = status.TTS.Error
			if st.Reason == "" {
				st.Reason = "the TTS capability is applied but its provider is not initialized"
			}
			report.Statuses[key] = st
			return
		}
		st.State = SettingApplied
		st.Applied = desired.TTSVoice
		report.Applied[key] = desired.TTSVoice
		report.Statuses[key] = st
	}

	setTTSVoice()

	setTTSSpeed := func() {
		key := "ttsSpeed"
		st := report.Statuses[key]
		// Preserve the active generation's speed even when the newest request
		// failed or selected a different engine. The state below says whether
		// the current request owns that value.
		if appliedOK {
			// Do not publish a prior value for auto-select: that request has no
			// engine-specific speed contract. An explicit engine (including a
			// failed disable/replacement) may preserve its last effective value.
			if _, _, _, explicit := ttsSpeedRangeFor(desired.TTSEngine); explicit {
				if got, ok := effectiveTTSSpeedFor(applied); ok {
					report.Applied[key] = got
				}
			}
		}
		if issues[key] != "" {
			report.Statuses[key] = st
			return
		}
		if !desired.TTSEnabled {
			st.State = SettingUnavailable
			st.Reason = "TTS is disabled; this value is saved but not applied"
			report.Statuses[key] = st
			return
		}
		min, max, supported, explicit := ttsSpeedRangeFor(desired.TTSEngine)
		if !explicit {
			st.State = SettingUnavailable
			st.Reason = "auto-select does not identify a speed-capable TTS engine; choose an explicit engine"
			report.Statuses[key] = st
			return
		}
		if !supported {
			st.State = SettingUnsupported
			st.Reason = fmt.Sprintf("the selected TTS engine %q does not support speech speed", strings.TrimSpace(desired.TTSEngine))
			report.Statuses[key] = st
			return
		}
		if desired.TTSSpeed < min || desired.TTSSpeed > max {
			st.State = SettingInvalid
			st.Reason = fmt.Sprintf("speech speed for %s must be between %g and %g, got %g", strings.TrimSpace(desired.TTSEngine), min, max, desired.TTSSpeed)
			report.Statuses[key] = st
			return
		}
		if !appliedOK {
			st.State, st.Reason = notEffectiveState(status, applyErr)
			report.Statuses[key] = st
			return
		}
		if !applied.TTSEnabled {
			st.State, st.Reason = notEffectiveState(status, applyErr)
			report.Statuses[key] = st
			return
		}
		got, effective := effectiveTTSSpeedFor(applied)
		if !effective {
			st.State, st.Reason = notEffectiveState(status, applyErr)
			report.Statuses[key] = st
			return
		}
		if got != desired.TTSSpeed {
			st.State, st.Reason = notEffectiveState(status, applyErr)
			report.Statuses[key] = st
			return
		}
		if !status.TTS.Ready {
			st.State = SettingUnavailable
			st.Reason = status.TTS.Error
			if st.Reason == "" {
				st.Reason = "the TTS capability is applied but its provider is not initialized"
			}
			report.Statuses[key] = st
			return
		}
		st.State = SettingApplied
		st.Applied = desired.TTSSpeed
		report.Applied[key] = desired.TTSSpeed
		report.Statuses[key] = st
	}

	setTTSSpeed()

	setASRDecode := func(key string, want, got interface{}, matches bool) {
		st := report.Statuses[key]
		// Preserve the active generation's values before considering the
		// newest request. An invalid or rejected replacement must not erase the
		// last effective Whisper configuration from the report.
		if appliedOK && applied.ASREnabled && isWhisperASREngine(applied.ASREngine) {
			report.Applied[key] = got
		}
		if issues[key] != "" {
			report.Statuses[key] = st
			return
		}
		if !desired.ASREnabled {
			st.State = SettingUnavailable
			st.Reason = "ASR is disabled; this value is saved but not applied"
			report.Statuses[key] = st
			return
		}
		if !isWhisperASREngine(desired.ASREngine) {
			engine := strings.TrimSpace(desired.ASREngine)
			if engine == "" || engine == "auto" {
				st.State = SettingUnavailable
				st.Reason = "auto-select does not identify a Whisper engine; choose an explicit Whisper engine"
			} else {
				st.State = SettingUnsupported
				st.Reason = fmt.Sprintf("the selected ASR engine %q does not support Whisper decode options", engine)
			}
			report.Statuses[key] = st
			return
		}
		if !appliedOK {
			st.State, st.Reason = notEffectiveState(status, applyErr)
			report.Statuses[key] = st
			return
		}
		if !applied.ASREnabled || !isWhisperASREngine(applied.ASREngine) {
			st.State, st.Reason = notEffectiveState(status, applyErr)
			report.Statuses[key] = st
			return
		}
		if !status.ASR.Ready {
			st.State = SettingUnavailable
			st.Reason = status.ASR.Error
			if st.Reason == "" {
				st.Reason = "the ASR capability is applied but its Whisper provider is not initialized"
			}
			report.Statuses[key] = st
			return
		}
		if matches {
			st.State = SettingApplied
			st.Applied = want
		} else {
			st.State, st.Reason = notEffectiveState(status, applyErr)
		}
		report.Statuses[key] = st
	}

	setASRDecode("asrLanguage", desired.ASRLanguage, applied.ASRLanguage, desired.ASRLanguage == applied.ASRLanguage)
	setASRDecode("asrThreads", desired.ASRThreads, applied.ASRThreads, desired.ASRThreads == applied.ASRThreads)
	setASRDecode("asrBeamSize", desired.ASRBeamSize, applied.ASRBeamSize, desired.ASRBeamSize == applied.ASRBeamSize)
	setASRDecode("asrBestOf", desired.ASRBestOf, applied.ASRBestOf, desired.ASRBestOf == applied.ASRBestOf)

	for key, reason := range issues {
		st, ok := report.Statuses[key]
		if !ok {
			st = SettingStatus{Key: key, Owner: SettingOwnerNone}
		}
		st.State = SettingInvalid
		st.Applied = nil
		st.Reason = reason
		report.Statuses[key] = st
	}

	report.Counts = countSettingStates(report.Statuses)
	return report
}

func notEffectiveState(status Status, applyErr string) (SettingState, string) {
	if applyErr != "" {
		return SettingRejected, applyErr
	}
	if status.Builds > 0 {
		return SettingPending, "an audio activation is still in progress"
	}
	return SettingRejected, "the requested value is not in the active generation"
}

func countSettingStates(statuses map[string]SettingStatus) SettingsCounts {
	counts := SettingsCounts{Total: len(statuses)}
	for _, st := range statuses {
		switch st.State {
		case SettingApplied:
			counts.Applied++
		case SettingPending:
			counts.Pending++
		case SettingRejected:
			counts.Rejected++
		case SettingInvalid:
			counts.Invalid++
		case SettingUnsupported:
			counts.Unsupported++
		case SettingUnavailable:
			counts.Unavailable++
		case SettingRenderer:
			counts.Renderer++
		}
	}
	return counts
}

func unsupportedSettings() map[string]string {
	return map[string]string{
		"asrDevice":             "Compute-device selection is not implemented by the current ASR adapters.",
		"asrTemperature":        "ASR decode options are not wired in this build.",
		"asrTranslateToEnglish": "This ASR option is not implemented in this build.",
		"asrNoTimestamps":       "This ASR option has no observable output in this build.",
		"asrFlashAttention":     "This ASR option is not implemented in this build.",
		"vadEnabled":            "Voice activity detection is not implemented in this build.",
		"vadThreshold":          "Voice activity detection is not implemented in this build.",
		"vadMinSpeechMs":        "Voice activity detection is not implemented in this build.",
		"vadMinSilenceMs":       "Voice activity detection is not implemented in this build.",
		"vadPaddingMs":          "Voice activity detection is not implemented in this build.",
		"ttsNoiseScale":         "This TTS option is not wired in this build.",
		"ttsNoiseW":             "This TTS option is not wired in this build.",
		"ttsSentenceSilence":    "This TTS option is not wired in this build.",
		"ttsCacheEnabled":       "The TTS cache control is not implemented in this build.",
		"ttsCacheSizeMB":        "The TTS cache control is not implemented in this build.",
		"audioInputDevice":      "Capture device selection is not applied in this build.",
		"audioOutputDevice":     "Output device selection is not applied in this build.",
		"audioSampleRate":       "Capture format selection is not applied in this build.",
		"audioChannels":         "Capture format selection is not applied in this build.",
		"audioBitDepth":         "Capture format selection is not applied in this build.",
		"micSensitivity":        "Microphone sensitivity is not implemented in this build.",
		"pushToTalk":            "Push-to-talk input mode is not implemented in this build.",
	}
}

func rendererSettings() map[string]string {
	return map[string]string{
		"ttsVolume":          "Applied by the renderer; not applied by the Go backend.",
		"ttsAutoplayMessage": "Applied by the renderer; not applied by the Go backend.",
	}
}
