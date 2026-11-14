package main

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"howl-chat/internal/audio"
	appsettings "howl-chat/internal/settings"
)

func contractDesiredSettings() audio.DesiredSettings {
	desired := audio.DefaultDesiredSettings()
	desired.ASREnabled = true
	desired.TTSEnabled = true
	desired.ASREngine = "whisper-base"
	desired.ASRLanguage = "en"
	desired.ASRDevice = "cuda"
	desired.ASRThreads = 6
	desired.ASRBeamSize = 4
	desired.ASRBestOf = 4
	desired.ASRTemperature = 0.1
	desired.ASRTranslateToEnglish = true
	desired.ASRNoTimestamps = true
	desired.ASRFlashAttention = false
	desired.VADEnabled = true
	desired.VADThreshold = 0.55
	desired.VADMinSpeechMs = 300
	desired.VADMinSilenceMs = 120
	desired.VADPaddingMs = 40
	desired.TTSEngine = "piper"
	desired.TTSVoice = "voice-test"
	desired.TTSSpeed = 1.25
	desired.TTSVolume = 0.55
	desired.TTSNoiseScale = 0.5
	desired.TTSNoiseW = 0.6
	desired.TTSSentenceSilence = 0.3
	desired.TTSAutoplayMessage = false
	desired.TTSCacheEnabled = false
	desired.TTSCacheSizeMB = 512
	desired.AudioInputDevice = "microphone-1"
	desired.AudioOutputDevice = "speaker-1"
	desired.AudioSampleRate = 48000
	desired.AudioChannels = "stereo"
	desired.AudioBitDepth = 24
	desired.MicSensitivity = 0.7
	desired.PushToTalk = false
	return desired
}

func desiredSettingsMap(t *testing.T, desired audio.DesiredSettings) map[string]interface{} {
	t.Helper()
	data, err := json.Marshal(desired)
	if err != nil {
		t.Fatalf("marshal desired settings: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal desired settings: %v", err)
	}
	return raw
}

func newContractTestApp(t *testing.T, factory audio.PipelineFactory) (*App, *audio.Manager) {
	t.Helper()
	manager := audio.NewManager(factory)
	app := &App{
		settingsStore: appsettings.NewStore(t.TempDir()),
		appSettings:   map[string]interface{}{},
		audioManager:  manager,
	}
	app.setAudioDesired(audio.DefaultDesiredSettings())
	t.Cleanup(func() { _ = manager.Close(context.Background()) })
	return app, manager
}

func TestAllAudioSettingsPersistWithoutApplyingAndReportTruthfully(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateReady, "", 1)))
	pipeline := &appFakePipeline{}
	app, _ := newContractTestApp(t, func(_ context.Context, desired audio.Settings, candidate *audio.Candidate) error {
		*candidate = *appCandidate(pipeline, true, true, nil, nil)
		return nil
	})
	want := contractDesiredSettings()

	if err := app.SetSettings(map[string]interface{}{"audio": desiredSettingsMap(t, want)}); err != nil {
		t.Fatalf("SetSettings: %v", err)
	}
	if got := app.currentAudioDesired(); !reflect.DeepEqual(got, want) {
		t.Fatalf("in-memory desired settings differ:\n got: %+v\nwant: %+v", got, want)
	}
	loaded, err := app.settingsStore.Load()
	if err != nil {
		t.Fatalf("reload settings: %v", err)
	}
	rawAudio, ok := loaded["audio"].(map[string]interface{})
	if !ok {
		t.Fatalf("loaded audio section = %T, want map", loaded["audio"])
	}
	if got := audio.DesiredSettingsFromMap(rawAudio); !reflect.DeepEqual(got, want) {
		t.Fatalf("persisted desired settings differ:\n got: %+v\nwant: %+v", got, want)
	}
	if got := app.audioManager.Status().DesiredGen; got != 0 {
		t.Fatalf("SetSettings advanced the manager desired generation to %d", got)
	}

	if err := app.ApplyAudioSettings(want); err != nil {
		t.Fatalf("ApplyAudioSettings: %v", err)
	}
	reportRaw, ok := app.GetAudioStatus()["settings"].(audio.SettingsReport)
	if !ok {
		t.Fatalf("settings report has type %T", app.GetAudioStatus()["settings"])
	}
	if reportRaw.Counts.Total != 34 {
		t.Fatalf("report total = %d, want 34", reportRaw.Counts.Total)
	}
	if reportRaw.Counts.Applied != 10 {
		t.Fatalf("backend applied = %d, want 10", reportRaw.Counts.Applied)
	}
	if reportRaw.Counts.Unsupported != 22 {
		t.Fatalf("unsupported = %d, want 22", reportRaw.Counts.Unsupported)
	}
	if reportRaw.Counts.Renderer != 2 {
		t.Fatalf("renderer = %d, want 2", reportRaw.Counts.Renderer)
	}
	if len(reportRaw.Applied) != 10 {
		t.Fatalf("applied map length = %d, want 10", len(reportRaw.Applied))
	}
	if got := reportRaw.Applied["ttsVoice"]; got != "voice-test" {
		t.Fatalf("applied ttsVoice = %v, want voice-test", got)
	}
	if got := reportRaw.Applied["asrLanguage"]; got != "en" {
		t.Fatalf("applied asrLanguage = %v, want en", got)
	}
	if got := reportRaw.Applied["asrThreads"]; got != 6 {
		t.Fatalf("applied asrThreads = %v, want 6", got)
	}
	if got := reportRaw.Applied["asrBeamSize"]; got != 4 {
		t.Fatalf("applied asrBeamSize = %v, want 4", got)
	}
	if got := reportRaw.Applied["asrBestOf"]; got != 4 {
		t.Fatalf("applied asrBestOf = %v, want 4", got)
	}
	if got := reportRaw.Applied["ttsSpeed"]; got != 1.25 {
		t.Fatalf("applied ttsSpeed = %v, want 1.25", got)
	}
	if _, ok := reportRaw.Applied["ttsVolume"]; ok {
		t.Fatal("renderer-owned ttsVolume appeared in the backend applied map")
	}
	if got := reportRaw.Statuses["asrLanguage"].State; got != audio.SettingApplied {
		t.Fatalf("asrLanguage state = %q, want applied", got)
	}
	if got := reportRaw.Statuses["ttsVolume"].Owner; got != audio.SettingOwnerRenderer {
		t.Fatalf("ttsVolume owner = %q, want renderer", got)
	}
}

func TestFailedTTSSpeedReplacementKeepsLastAppliedSpeed(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateReady, "", 1)))
	call := 0
	replacementErr := errors.New("replacement TTS failed")
	app, _ := newContractTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		call++
		if call == 1 {
			*candidate = *appCandidate(&appFakePipeline{}, false, true, nil, nil)
			return nil
		}
		return replacementErr
	})

	first := audio.DefaultDesiredSettings()
	first.TTSEnabled = true
	first.TTSEngine = "piper"
	first.TTSSpeed = 1.25
	if err := app.ApplyAudioSettings(first); err != nil {
		t.Fatalf("initial apply: %v", err)
	}

	changed := first
	changed.TTSSpeed = 1.5
	if err := app.ApplyAudioSettings(changed); err == nil {
		t.Fatal("replacement apply unexpectedly succeeded")
	}

	report := app.GetAudioStatus()["settings"].(audio.SettingsReport)
	if got := report.Applied["ttsSpeed"]; got != 1.25 {
		t.Fatalf("applied ttsSpeed = %v, want previous 1.25", got)
	}
	if got := report.Statuses["ttsSpeed"].State; got != audio.SettingRejected {
		t.Fatalf("ttsSpeed state = %q, want rejected", got)
	}
	if got := report.Statuses["ttsSpeed"].Reason; got == "" {
		t.Fatal("failed speed replacement has no explicit reason")
	}
}

func TestFailedApplyKeepsLastAppliedConfiguration(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateReady, "", 1)))
	call := 0
	replacementErr := errors.New("replacement ASR failed")
	app, _ := newContractTestApp(t, func(_ context.Context, desired audio.Settings, candidate *audio.Candidate) error {
		call++
		if call == 1 {
			*candidate = *appCandidate(&appFakePipeline{}, true, false, nil, nil)
			return nil
		}
		return replacementErr
	})

	first := desiredWithBackend(audio.Settings{
		ASREnabled: true, ASREngine: "whisper-base",
		ASRLanguage: "de", ASRThreads: 8, ASRBeamSize: 3, ASRBestOf: 2,
	})
	if err := app.ApplyAudioSettings(first); err != nil {
		t.Fatalf("initial apply: %v", err)
	}

	changed := first
	changed.ASREngine = "whisper-small"
	if err := app.ApplyAudioSettings(changed); err == nil {
		t.Fatal("replacement apply unexpectedly succeeded")
	}

	status := app.GetAudioStatus()
	report, ok := status["settings"].(audio.SettingsReport)
	if !ok {
		t.Fatalf("settings report has type %T", status["settings"])
	}
	if got := report.Applied["asrEngine"]; got != "whisper-base" {
		t.Fatalf("applied asrEngine = %v, want the last successful whisper-base", got)
	}
	if got := report.Statuses["asrEngine"].State; got != audio.SettingRejected {
		t.Fatalf("asrEngine state = %q, want rejected", got)
	}
	if got := report.Statuses["asrEngine"].Reason; got == "" {
		t.Fatal("failed replacement has no explicit reason")
	}
	if got := report.Applied["asrEnabled"]; got != true {
		t.Fatalf("unchanged asrEnabled = %v, want true", got)
	}
	if got := report.Applied["asrLanguage"]; got != "de" {
		t.Fatalf("applied asrLanguage = %v, want previous de", got)
	}
	if got := report.Applied["asrThreads"]; got != 8 {
		t.Fatalf("applied asrThreads = %v, want previous 8", got)
	}
	if got := report.Applied["asrBeamSize"]; got != 3 {
		t.Fatalf("applied asrBeamSize = %v, want previous 3", got)
	}
	if got := report.Applied["asrBestOf"]; got != 2 {
		t.Fatalf("applied asrBestOf = %v, want previous 2", got)
	}
}

func TestInvalidWhisperDecodeValueFailsBeforeActivation(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateReady, "", 1)))
	app, _ := newContractTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		*candidate = *appCandidate(&appFakePipeline{}, true, false, nil, nil)
		return nil
	})
	desired := audio.DefaultDesiredSettings()
	desired.ASREnabled = true
	desired.ASREngine = "whisper-base"
	desired.ASRThreads = 8
	if err := app.ApplyAudioSettings(desired); err != nil {
		t.Fatalf("initial apply: %v", err)
	}

	desired.ASRThreads = 0
	err := app.ApplyAudioSettings(desired)
	var validationErr *audio.SettingsValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("ApplyAudioSettings error = %v, want SettingsValidationError", err)
	}
	if validationErr.Issues["asrThreads"] == "" {
		t.Fatalf("validation error has no thread reason: %#v", validationErr.Issues)
	}
	report := app.GetAudioStatus()["settings"].(audio.SettingsReport)
	if got := report.Statuses["asrThreads"].State; got != audio.SettingInvalid {
		t.Fatalf("asrThreads state = %q, want invalid", got)
	}
	if got := report.Applied["asrThreads"]; got != 8 {
		t.Fatalf("applied asrThreads = %v, want the last effective 8", got)
	}
}

func TestWhisperLanguageOutsideProfileFailsBeforeActivation(t *testing.T) {
	app, _ := newContractTestApp(t, func(_ context.Context, _ audio.Settings, _ *audio.Candidate) error {
		t.Fatal("factory ran for a language the selected Whisper profile does not support")
		return nil
	})
	desired := audio.DefaultDesiredSettings()
	desired.ASREnabled = true
	desired.ASREngine = "whisper-tiny"
	desired.ASRLanguage = "nl"

	err := app.ApplyAudioSettings(desired)
	var validationErr *audio.SettingsValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("ApplyAudioSettings error = %v, want SettingsValidationError", err)
	}
	if validationErr.Issues["asrLanguage"] == "" {
		t.Fatalf("validation error has no language reason: %#v", validationErr.Issues)
	}
	report := app.GetAudioStatus()["settings"].(audio.SettingsReport)
	if got := report.Statuses["asrLanguage"].State; got != audio.SettingInvalid {
		t.Fatalf("asrLanguage state = %q, want invalid", got)
	}
}

func TestInvalidDecodeValueForUnsupportedEngineDoesNotBlockCoreApply(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateReady, "", 1)))
	app, _ := newContractTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		*candidate = *appCandidate(&appFakePipeline{}, true, false, nil, nil)
		return nil
	})
	desired := audio.DefaultDesiredSettings()
	desired.ASREnabled = true
	desired.ASREngine = "qwen3-asr-1.7b"
	desired.ASRThreads = 0

	if err := app.ApplyAudioSettings(desired); err != nil {
		t.Fatalf("ApplyAudioSettings: %v", err)
	}
	report := app.GetAudioStatus()["settings"].(audio.SettingsReport)
	if got := report.Statuses["asrThreads"].State; got != audio.SettingInvalid {
		t.Fatalf("asrThreads state = %q, want invalid", got)
	}
	if got := report.Statuses["asrEngine"].State; got != audio.SettingApplied {
		t.Fatalf("asrEngine state = %q, want applied", got)
	}
	if _, ok := report.Applied["asrThreads"]; ok {
		t.Fatal("invalid decode value for an unsupported engine appeared in applied")
	}
}

func TestInvalidBackendValueFailsBeforeActivation(t *testing.T) {
	app, _ := newContractTestApp(t, func(_ context.Context, _ audio.Settings, _ *audio.Candidate) error {
		t.Fatal("factory ran for an invalid backend value")
		return nil
	})
	desired := audio.DefaultDesiredSettings()
	desired.ASREnabled = true
	desired.ASREngine = "not-a-real-engine"

	err := app.ApplyAudioSettings(desired)
	var validationErr *audio.SettingsValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("ApplyAudioSettings error = %v, want SettingsValidationError", err)
	}
	if validationErr.Issues["asrEngine"] == "" {
		t.Fatalf("validation error has no engine reason: %#v", validationErr.Issues)
	}
	report := app.GetAudioStatus()["settings"].(audio.SettingsReport)
	if got := report.Statuses["asrEngine"].State; got != audio.SettingInvalid {
		t.Fatalf("asrEngine state = %q, want invalid", got)
	}
	if _, ok := report.Applied["asrEngine"]; ok {
		t.Fatal("invalid engine appeared in the applied map")
	}
}

func TestSavedSettingsDoNotActivateAtStartup(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateReady, "", 1)))
	app, _ := newContractTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		t.Fatal("startup status activated a persisted selection")
		return nil
	})
	desired := contractDesiredSettings()
	if err := app.SetSettings(map[string]interface{}{"audio": desiredSettingsMap(t, desired)}); err != nil {
		t.Fatalf("SetSettings: %v", err)
	}

	_ = app.GetAudioStatus()
	if _, applied := app.audioManager.AppliedSettings(); applied {
		t.Fatal("persisted settings were applied without an explicit activation")
	}
}
