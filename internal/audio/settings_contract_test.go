package audio

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

func TestAllDesiredAudioSettingsSurviveMapRoundTrip(t *testing.T) {
	want := DefaultDesiredSettings()
	want.ASREnabled = true
	want.TTSEnabled = true
	want.ASREngine = "whisper-base"
	want.ASRLanguage = "en"
	want.ASRDevice = "cuda"
	want.ASRThreads = 8
	want.ASRBeamSize = 3
	want.ASRBestOf = 4
	want.ASRTemperature = 0.2
	want.ASRTranslateToEnglish = true
	want.ASRNoTimestamps = true
	want.ASRFlashAttention = false
	want.VADEnabled = true
	want.VADThreshold = 0.6
	want.VADMinSpeechMs = 300
	want.VADMinSilenceMs = 150
	want.VADPaddingMs = 40
	want.TTSEngine = "piper"
	want.TTSVoice = "voice-test"
	want.TTSSpeed = 1.25
	want.TTSVolume = 0.55
	want.TTSNoiseScale = 0.5
	want.TTSNoiseW = 0.6
	want.TTSSentenceSilence = 0.3
	want.TTSAutoplayMessage = false
	want.TTSCacheEnabled = false
	want.TTSCacheSizeMB = 512
	want.AudioInputDevice = "microphone-1"
	want.AudioOutputDevice = "speaker-1"
	want.AudioSampleRate = 48000
	want.AudioChannels = "stereo"
	want.AudioBitDepth = 24
	want.MicSensitivity = 0.7
	want.PushToTalk = false

	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal desired settings: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal raw settings: %v", err)
	}
	if len(raw) != 34 {
		t.Fatalf("serialized settings keys = %d, want 34", len(raw))
	}

	got := DesiredSettingsFromMap(raw)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round-tripped settings differ:\n got: %+v\nwant: %+v", got, want)
	}
}

func TestSettingsReportSeparatesBackendRendererAndUnsupported(t *testing.T) {
	desired := DefaultDesiredSettings()
	desired.ASREnabled = true
	desired.ASREngine = "whisper-base"
	desired.TTSEnabled = false
	desired.TTSEngine = "piper"

	report := BuildSettingsReport(
		desired,
		Settings{
			ASREnabled: true, TTSEnabled: false, ASREngine: "whisper-base", TTSEngine: "piper",
			ASRLanguage: "auto", ASRThreads: 4, ASRBeamSize: 5, ASRBestOf: 5,
		},
		true,
		Status{ASR: Capability{Requested: true, Ready: true}, TTS: Capability{Requested: false}},
		"",
	)

	if report.Counts.Total != 34 {
		t.Fatalf("total statuses = %d, want 34", report.Counts.Total)
	}
	if report.Counts.Unsupported != 22 {
		t.Fatalf("unsupported statuses = %d, want 22", report.Counts.Unsupported)
	}
	if report.Counts.Renderer != 2 {
		t.Fatalf("renderer statuses = %d, want 2", report.Counts.Renderer)
	}
	if report.Counts.Applied != 7 {
		t.Fatalf("backend applied statuses = %d, want 7 (ASR enabled, ASR engine, TTS disabled, four decode settings)", report.Counts.Applied)
	}
	if got := report.Statuses["asrLanguage"].State; got != SettingApplied {
		t.Fatalf("asrLanguage state = %q, want applied", got)
	}
	if got := report.Applied["asrLanguage"]; got != "auto" {
		t.Fatalf("applied asrLanguage = %v, want auto", got)
	}
	if got := report.Statuses["ttsVolume"].Owner; got != SettingOwnerRenderer {
		t.Fatalf("ttsVolume owner = %q, want renderer", got)
	}
	if got := report.Statuses["ttsVolume"].State; got != SettingRenderer {
		t.Fatalf("ttsVolume state = %q, want renderer", got)
	}
	if _, ok := report.Applied["ttsVolume"]; ok {
		t.Fatal("renderer-owned ttsVolume appeared in the backend applied map")
	}
	if got := report.Statuses["asrEngine"].State; got != SettingApplied {
		t.Fatalf("asrEngine state = %q, want applied", got)
	}
	if got := report.Statuses["ttsEngine"].State; got != SettingUnavailable {
		t.Fatalf("disabled ttsEngine state = %q, want unavailable", got)
	}
	if got := report.Statuses["ttsSpeed"].State; got != SettingUnavailable {
		t.Fatalf("disabled ttsSpeed state = %q, want unavailable", got)
	}
}

func TestSettingsReportPartialFailureUsesLastAppliedConfiguration(t *testing.T) {
	desired := DefaultDesiredSettings()
	desired.ASREnabled = true
	desired.ASREngine = "whisper-small"
	desired.TTSVolume = 0.3

	report := BuildSettingsReport(
		desired,
		Settings{
			ASREnabled: true, TTSEnabled: false, ASREngine: "whisper-base", TTSEngine: "piper",
			ASRLanguage: "auto", ASRThreads: 4, ASRBeamSize: 5, ASRBestOf: 5,
		},
		true,
		Status{ASR: Capability{Requested: true, Ready: true}},
		"replacement ASR failed",
	)

	if got := report.Applied["asrEnabled"]; got != true {
		t.Fatalf("applied asrEnabled = %v, want true", got)
	}
	if got := report.Applied["asrEngine"]; got != "whisper-base" {
		t.Fatalf("applied asrEngine = %v, want previous whisper-base", got)
	}
	if got := report.Statuses["asrEngine"].State; got != SettingRejected {
		t.Fatalf("asrEngine state = %q, want rejected", got)
	}
	if got := report.Statuses["asrEngine"].Reason; got != "replacement ASR failed" {
		t.Fatalf("asrEngine reason = %q, want the activation failure", got)
	}
	if got := report.Statuses["ttsVolume"].State; got != SettingRenderer {
		t.Fatalf("ttsVolume state = %q, want renderer", got)
	}
	if report.Counts.Rejected == 0 {
		t.Fatal("partial failure did not report a rejected setting")
	}
}

func TestSettingsReportFailedDisableKeepsOldAppliedEngine(t *testing.T) {
	desired := DefaultDesiredSettings()
	desired.ASREnabled = false
	desired.ASREngine = "whisper-base"
	report := BuildSettingsReport(
		desired,
		Settings{ASREnabled: true, ASREngine: "whisper-base"},
		true,
		Status{ASR: Capability{Requested: true, Ready: true}},
		"disable failed",
	)

	if got := report.Applied["asrEngine"]; got != "whisper-base" {
		t.Fatalf("applied asrEngine = %v, want previous whisper-base", got)
	}
	if got := report.Statuses["asrEnabled"].State; got != SettingRejected {
		t.Fatalf("asrEnabled state = %q, want rejected", got)
	}
	if got := report.Statuses["asrEngine"].State; got != SettingRejected {
		t.Fatalf("asrEngine state = %q, want rejected", got)
	}
}

func TestSettingsReportWhisperDecodeStatusIsEngineSpecific(t *testing.T) {
	desired := DefaultDesiredSettings()
	desired.ASREnabled = true
	desired.ASREngine = "qwen3-asr-1.7b"
	desired.ASRLanguage = "de"
	desired.ASRThreads = 8
	desired.ASRBeamSize = 3
	desired.ASRBestOf = 2

	report := BuildSettingsReport(
		desired,
		Settings{ASREnabled: true, ASREngine: "qwen3-asr-1.7b"},
		true,
		Status{ASR: Capability{Requested: true, Ready: true}},
		"",
	)

	for _, key := range asrDecodeSettingKeys {
		if got := report.Statuses[key].State; got != SettingUnsupported {
			t.Fatalf("%s state = %q, want unsupported for a non-Whisper engine", key, got)
		}
		if _, ok := report.Applied[key]; ok {
			t.Fatalf("unsupported %s appeared in the applied map", key)
		}
	}

	desired.ASREngine = "auto"
	report = BuildSettingsReport(desired, Settings{}, false, Status{}, "")
	for _, key := range asrDecodeSettingKeys {
		if got := report.Statuses[key].State; got != SettingUnavailable {
			t.Fatalf("%s state = %q, want unavailable for auto-select", key, got)
		}
		if _, ok := report.Applied[key]; ok {
			t.Fatalf("auto-select %s appeared in the applied map", key)
		}
	}
}

func TestSettingsReportTTSSpeedIsEngineSpecific(t *testing.T) {
	desired := DefaultDesiredSettings()
	desired.TTSEnabled = true
	desired.TTSEngine = "piper"
	desired.TTSSpeed = 1.25

	applied := Settings{TTSEnabled: true, TTSEngine: "piper", TTSSpeed: 1.25}
	report := BuildSettingsReport(
		desired,
		applied,
		true,
		Status{TTS: Capability{Requested: true, Ready: true}},
		"",
	)
	if got := report.Statuses["ttsSpeed"].State; got != SettingApplied {
		t.Fatalf("piper ttsSpeed state = %q, want applied", got)
	}
	if got := report.Applied["ttsSpeed"]; got != 1.25 {
		t.Fatalf("applied piper speed = %v, want 1.25", got)
	}

	desired.TTSEngine = "auto"
	report = BuildSettingsReport(
		desired,
		applied,
		true,
		Status{TTS: Capability{Requested: true, Ready: true}},
		"",
	)
	if got := report.Statuses["ttsSpeed"].State; got != SettingUnavailable {
		t.Fatalf("auto ttsSpeed state = %q, want unavailable", got)
	}
	if _, ok := report.Applied["ttsSpeed"]; ok {
		t.Fatal("auto-select claimed an applied speech speed")
	}

	desired.TTSEngine = "not-a-real-engine"
	report = BuildSettingsReport(desired, Settings{}, false, Status{}, "")
	if got := report.Statuses["ttsSpeed"].State; got != SettingUnsupported {
		t.Fatalf("unsupported engine ttsSpeed state = %q, want unsupported", got)
	}
	if _, ok := report.Applied["ttsSpeed"]; ok {
		t.Fatal("unsupported engine claimed an applied speech speed")
	}
}

func TestSettingsReportMeloTTSSpeedGuardIsExplicit(t *testing.T) {
	desired := DefaultDesiredSettings()
	desired.TTSEnabled = true
	desired.TTSEngine = "melotts"
	desired.TTSSpeed = 1.25

	issues := desired.Validate()
	if issues["ttsSpeed"] == "" {
		t.Fatal("MeloTTS outside its product guard has no validation reason")
	}
	if !HasBackendValidationIssue(desired, issues) {
		t.Fatal("MeloTTS outside its product guard did not block backend activation")
	}
	report := BuildSettingsReport(desired, Settings{}, false, Status{}, "")
	if got := report.Statuses["ttsSpeed"].State; got != SettingInvalid {
		t.Fatalf("MeloTTS out-of-guard state = %q, want invalid", got)
	}

	desired.TTSEnabled = false
	issues = desired.Validate()
	if issues["ttsSpeed"] != "" {
		t.Fatalf("disabled TTS unexpectedly blocked a saved speed: %q", issues["ttsSpeed"])
	}
	if HasBackendValidationIssue(desired, issues) {
		t.Fatal("disabled TTS speed blocked backend activation")
	}

	desired.TTSEnabled = true
	desired.TTSSpeed = 1.0
	report = BuildSettingsReport(
		desired,
		Settings{TTSEnabled: true, TTSEngine: "melotts", TTSSpeed: 1.0},
		true,
		Status{TTS: Capability{Requested: true, Ready: true}},
		"",
	)
	if got := report.Statuses["ttsSpeed"].State; got != SettingApplied {
		t.Fatalf("MeloTTS in-guard state = %q, want applied", got)
	}
}

func TestSettingsReportFailedTTSSpeedReplacementPreservesPriorValue(t *testing.T) {
	desired := DefaultDesiredSettings()
	desired.TTSEnabled = true
	desired.TTSEngine = "piper"
	desired.TTSSpeed = 1.5

	report := BuildSettingsReport(
		desired,
		Settings{TTSEnabled: true, TTSEngine: "piper", TTSSpeed: 1.25},
		true,
		Status{TTS: Capability{Requested: true, Ready: true}},
		"replacement failed",
	)
	if got := report.Applied["ttsSpeed"]; got != 1.25 {
		t.Fatalf("applied ttsSpeed = %v, want previous 1.25", got)
	}
	if got := report.Statuses["ttsSpeed"].State; got != SettingRejected {
		t.Fatalf("ttsSpeed state = %q, want rejected", got)
	}
	if got := report.Statuses["ttsSpeed"].Reason; got != "replacement failed" {
		t.Fatalf("ttsSpeed reason = %q, want replacement failure", got)
	}
}

func TestProductionFactoryForwardsTTSSpeed(t *testing.T) {
	built := NewPipeline()
	factory := newProductionFactory(func() *Pipeline { return built })
	desired := Settings{TTSEngine: "piper", TTSSpeed: 1.25}
	if err := factory(context.Background(), desired, &Candidate{}); err != nil {
		t.Fatalf("factory: %v", err)
	}
	if got := built.config.TTSSpeed; got != 1.25 {
		t.Fatalf("pipeline config TTSSpeed = %g, want 1.25", got)
	}
}

func TestSettingsReportValidationIsExplicitAndNeverApplied(t *testing.T) {
	desired := DefaultDesiredSettings()
	desired.ASRThreads = 0
	desired.TTSVolume = 2

	report := BuildSettingsReport(desired, Settings{}, false, Status{}, "")

	if got := report.Statuses["asrThreads"].State; got != SettingInvalid {
		t.Fatalf("asrThreads state = %q, want invalid", got)
	}
	if got := report.Statuses["asrThreads"].Reason; got == "" {
		t.Fatal("asrThreads has no validation reason")
	}
	if got := report.Statuses["ttsVolume"].State; got != SettingInvalid {
		t.Fatalf("ttsVolume state = %q, want invalid", got)
	}
	if _, ok := report.Applied["asrThreads"]; ok {
		t.Fatal("invalid unsupported setting appeared in the applied map")
	}
	if _, ok := report.Applied["ttsVolume"]; ok {
		t.Fatal("invalid renderer setting appeared in the backend applied map")
	}
	if report.Counts.Invalid != 2 {
		t.Fatalf("invalid count = %d, want 2", report.Counts.Invalid)
	}

	backendInvalid := DefaultDesiredSettings()
	backendInvalid.ASREnabled = true
	backendInvalid.ASREngine = "not-a-real-engine"
	issues := backendInvalid.Validate()
	if !HasBackendValidationIssue(backendInvalid, issues) {
		t.Fatalf("unknown backend engine was not a backend validation issue: %#v", issues)
	}
	if issues["asrEngine"] == "" {
		t.Fatal("unknown backend engine has no explicit reason")
	}
}
