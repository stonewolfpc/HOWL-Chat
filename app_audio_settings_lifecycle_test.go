package main

import (
	"context"
	"errors"
	"testing"

	audio "howl-chat/internal/audio"
	appsettings "howl-chat/internal/settings"
)

// The settings panel must be able to save a selection and a fresh store must
// read it back unchanged. This is the controlled restart portion of Unit B;
// the UI calls SetSettings before ApplyAudioSettings.
func TestAudioSettingsPersistAcrossStoreRestart(t *testing.T) {
	dir := t.TempDir()
	app := &App{
		settingsStore: appsettings.NewStore(dir),
		appSettings:   map[string]interface{}{},
	}
	payload := map[string]interface{}{
		"audio": map[string]interface{}{
			"asrEnabled": true,
			"ttsEnabled": false,
			"asrEngine":  "whisper-base",
			"ttsEngine":  "piper",
		},
	}
	if err := app.SetSettings(payload); err != nil {
		t.Fatalf("SetSettings: %v", err)
	}

	// A new Store over the same directory stands in for the next application
	// process opening the same settings file.
	restarted := appsettings.NewStore(dir)
	loaded, err := restarted.Load()
	if err != nil {
		t.Fatalf("Load after restart: %v", err)
	}
	audioSection, ok := loaded["audio"].(map[string]interface{})
	if !ok {
		t.Fatalf("loaded audio section = %T, want map", loaded["audio"])
	}
	if got := audioSection["asrEnabled"]; got != true {
		t.Fatalf("asrEnabled after restart = %v, want true", got)
	}
	if got := audioSection["asrEngine"]; got != "whisper-base" {
		t.Fatalf("asrEngine after restart = %v, want whisper-base", got)
	}
}

func TestRemoteASRApplyNeverReportsLocalReadiness(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateReady, "", 1)))
	for _, engine := range []string{"qwen3-asr-0.6b", "qwen3-asr-1.7b"} {
		t.Run(engine, func(t *testing.T) {
			pipeline := &appFakePipeline{}
			app, manager := newAudioTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
				*candidate = *appCandidate(pipeline, true, false, nil, nil)
				return nil
			})
			if _, _, err := manager.Apply(context.Background(), audio.Settings{ASREnabled: true, ASREngine: engine}); err != nil {
				t.Fatalf("Apply: %v", err)
			}

			status := app.GetAudioStatus()
			if status["asr_ready"] != false {
				t.Fatalf("asr_ready = %v, want false for a remote engine", status["asr_ready"])
			}
			if status["asr_readiness"] != "remote-unverified" {
				t.Fatalf("asr_readiness = %v, want remote-unverified", status["asr_readiness"])
			}
			if status["asr_error"] == "" {
				t.Fatal("remote readiness is missing its unverified explanation")
			}
		})
	}
}

func TestLocalASRReadinessIsInitializedOnly(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateReady, "", 1)))
	pipeline := &appFakePipeline{}
	app, manager := newAudioTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		*candidate = *appCandidate(pipeline, true, false, nil, nil)
		return nil
	})
	if _, _, err := manager.Apply(context.Background(), audio.Settings{ASREnabled: true, ASREngine: "whisper-base"}); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	status := app.GetAudioStatus()
	if status["asr_ready"] != true {
		t.Fatalf("asr_ready = %v, want true for an initialized local engine", status["asr_ready"])
	}
	if status["asr_readiness"] != "local-initialized" {
		t.Fatalf("asr_readiness = %v, want local-initialized", status["asr_readiness"])
	}
}

func TestASRReadinessReportsUnconfiguredAndDisabled(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateReady, "", 1)))
	pipeline := &appFakePipeline{}
	app, manager := newAudioTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		*candidate = *appCandidate(pipeline, true, false, nil, nil)
		return nil
	})

	status := app.GetAudioStatus()
	if status["asr_readiness"] != "disabled" {
		t.Fatalf("disabled readiness = %v, want disabled", status["asr_readiness"])
	}

	if _, _, err := manager.Apply(context.Background(), audio.Settings{ASREnabled: true, ASREngine: "auto"}); err != nil {
		t.Fatalf("Apply auto: %v", err)
	}
	status = app.GetAudioStatus()
	if status["asr_readiness"] != "unconfigured" {
		t.Fatalf("auto readiness = %v, want unconfigured", status["asr_readiness"])
	}
	if status["asr_ready"] != false {
		t.Fatalf("auto asr_ready = %v, want false", status["asr_ready"])
	}
}

func TestFailedActivationDoesNotEraseSavedSelection(t *testing.T) {
	dir := t.TempDir()
	app := &App{
		settingsStore: appsettings.NewStore(dir),
		appSettings:   map[string]interface{}{},
	}
	payload := map[string]interface{}{
		"audio": map[string]interface{}{
			"asrEnabled": true,
			"asrEngine":  "whisper-base",
		},
	}
	if err := app.SetSettings(payload); err != nil {
		t.Fatalf("SetSettings: %v", err)
	}

	activationErr := errors.New("whisper is not installed")
	app.audioManager = audio.NewManager(func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		*candidate = *appCandidate(&appFakePipeline{}, false, false, activationErr, nil)
		return nil
	})
	t.Cleanup(func() { _ = app.audioManager.Close(context.Background()) })
	if err := app.ApplyAudioSettings(desiredWithBackend(audio.Settings{ASREnabled: true, ASREngine: "whisper-base", TTSEngine: "piper"})); err == nil {
		t.Fatal("ApplyAudioSettings returned nil for a failed activation")
	}

	loaded, err := appsettings.NewStore(dir).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	audioSection, ok := loaded["audio"].(map[string]interface{})
	if !ok {
		t.Fatalf("loaded audio section = %T, want map", loaded["audio"])
	}
	if got := audioSection["asrEngine"]; got != "whisper-base" {
		t.Fatalf("saved asrEngine after failed apply = %v, want whisper-base", got)
	}
}
