package main

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/audio"
	audiotypes "howl-chat/internal/audio/types"
)

// micTestPipeline is a managed pipeline that records shutdowns so abort cleanup
// can be observed without any provider or hardware.
type micTestPipeline struct {
	shutdowns int32
}

func (p *micTestPipeline) ProcessAudioMessage(context.Context, audiotypes.AudioMessage) (*audiotypes.AudioResponse, error) {
	return &audiotypes.AudioResponse{}, nil
}

func (p *micTestPipeline) GetStatus() audiotypes.PipelineStatus { return audiotypes.PipelineStatus{} }

func (p *micTestPipeline) GetTTSVoices() []audiotypes.Voice { return nil }

func (p *micTestPipeline) Shutdown(context.Context) error {
	atomic.AddInt32(&p.shutdowns, 1)
	return nil
}

// useMicProbe replaces both microphone seams with one probe, so a test that only
// cares about the reported state does not need to know which level the code
// under test asks for.
func useMicProbe(t *testing.T, probe func() MicrophoneStatus) {
	t.Helper()
	prevProbe, prevPresence := microphoneProbe, microphonePresenceProbe
	microphoneProbe, microphonePresenceProbe = probe, probe
	t.Cleanup(func() {
		microphoneProbe, microphonePresenceProbe = prevProbe, prevPresence
	})
}

// useMicPresenceProbe replaces only the status-read seam. Call it after
// useMicProbe to observe which level a path actually uses.
func useMicPresenceProbe(t *testing.T, probe func() MicrophoneStatus) {
	t.Helper()
	prev := microphonePresenceProbe
	microphonePresenceProbe = probe
	t.Cleanup(func() { microphonePresenceProbe = prev })
}

func constantMicProbe(status MicrophoneStatus) func() MicrophoneStatus {
	return func() MicrophoneStatus { return status }
}

// queuedMicProbe returns each status once, then repeats the last one.
func queuedMicProbe(statuses ...MicrophoneStatus) func() MicrophoneStatus {
	var calls int32
	return func() MicrophoneStatus {
		index := int(atomic.AddInt32(&calls, 1)) - 1
		if index >= len(statuses) {
			index = len(statuses) - 1
		}
		return statuses[index]
	}
}

func newMicTestApp(t *testing.T) (*App, *int32, *[]*micTestPipeline) {
	t.Helper()
	var factoryCalls int32
	pipelines := &[]*micTestPipeline{}
	factory := func(_ context.Context, desired audio.Settings, candidate *audio.Candidate) error {
		atomic.AddInt32(&factoryCalls, 1)
		pipeline := &micTestPipeline{}
		*pipelines = append(*pipelines, pipeline)
		if candidate != nil {
			candidate.Pipeline = pipeline
			if desired.ASREnabled {
				candidate.ASR = audio.CapabilityResult{Ready: true}
			}
			if desired.TTSEnabled {
				candidate.TTS = audio.CapabilityResult{Ready: true}
			}
		}
		return nil
	}
	return newAudioTestAppWithManager(audio.NewManager(factory)), &factoryCalls, pipelines
}

func micStatusFromReport(t *testing.T, report map[string]interface{}) map[string]interface{} {
	t.Helper()
	raw, ok := report["mic"].(map[string]interface{})
	if !ok {
		t.Fatalf("GetAudioStatus has no mic record: %#v", report["mic"])
	}
	return raw
}

func TestMicrophoneGateRefusesEveryUnavailableState(t *testing.T) {
	tests := []struct {
		name    string
		status  MicrophoneStatus
		wantMsg string
	}{
		{"no device", newMicrophoneStatus(MicStateNoDevice, "GetDefaultAudioEndpoint returned E_NOTFOUND", 0), MicrophoneNoDeviceMessage},
		{"disconnected", newMicrophoneStatus(MicStateUnavailable, "GetState reports disabled", 1), ""},
		{"access denied", newMicrophoneStatus(MicStateAccessDenied, "Initialize returned E_ACCESSDENIED", 1), ""},
		{"capture init failed", newMicrophoneStatus(MicStateCaptureInitFailed, "Initialize returned E_INVALIDARG", 1), ""},
		{"unsupported platform", unsupportedMicrophoneStatus(), ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useMicProbe(t, constantMicProbe(test.status))
			app, factoryCalls, pipelines := newMicTestApp(t)

			err := app.applyAudioSettings(desiredWithBackend(audio.Settings{ASREnabled: true, ASREngine: "whisper-base"}))
			if err == nil {
				t.Fatal("ASR activation succeeded with an unusable microphone")
			}
			if got := atomic.LoadInt32(factoryCalls); got != 0 {
				t.Fatalf("pipeline factory ran %d time(s); ASR must be refused before allocation", got)
			}
			if len(*pipelines) != 0 {
				t.Fatalf("a pipeline was built despite the microphone gate: %d", len(*pipelines))
			}
			if err.Error() != test.status.Message {
				t.Fatalf("error = %q, want %q", err.Error(), test.status.Message)
			}
			if test.wantMsg != "" && err.Error() != test.wantMsg {
				t.Fatalf("no-device error = %q, want exact copy %q", err.Error(), test.wantMsg)
			}

			report := app.GetAudioStatus()
			if ready, _ := report["asr_ready"].(bool); ready {
				t.Fatal("asr_ready is true while the microphone is unusable")
			}
			if readiness, _ := report["asr_readiness"].(string); readiness != "mic-unavailable" {
				t.Fatalf("asr_readiness = %q, want mic-unavailable", readiness)
			}
			mic := micStatusFromReport(t, report)
			if state, _ := mic["state"].(string); state != string(test.status.State) {
				t.Fatalf("mic.state = %q, want %q", state, test.status.State)
			}
			if message, _ := mic["message"].(string); message != test.status.Message {
				t.Fatalf("mic.message = %q, want %q", message, test.status.Message)
			}
		})
	}
}

func TestUnsupportedPlatformStatusDoesNotClaimReady(t *testing.T) {
	status := unsupportedMicrophoneStatus()
	if status.State != MicStateUnsupported {
		t.Fatalf("state = %q, want unsupported", status.State)
	}
	if status.Usable() {
		t.Fatal("unsupported status is usable; ASR would be permitted without a probe")
	}
	if status.Message == "" {
		t.Fatal("unsupported status has no user-facing message")
	}

	useMicProbe(t, constantMicProbe(status))
	app, factoryCalls, _ := newMicTestApp(t)
	if err := app.applyAudioSettings(desiredWithBackend(audio.Settings{ASREnabled: true, ASREngine: "whisper-base"})); err == nil {
		t.Fatal("ASR activation succeeded on an unsupported platform")
	}
	if got := atomic.LoadInt32(factoryCalls); got != 0 {
		t.Fatalf("pipeline factory calls = %d, want 0", got)
	}
}

func TestMicrophoneReadyAllowsASRAndReportsLocalInitialized(t *testing.T) {
	ready := newMicrophoneStatus(MicStateReady, "", 1)
	ready.Role = "console"
	useMicProbe(t, constantMicProbe(ready))
	app, factoryCalls, _ := newMicTestApp(t)

	if err := app.applyAudioSettings(desiredWithBackend(audio.Settings{ASREnabled: true, ASREngine: "whisper-base"})); err != nil {
		t.Fatalf("ASR activation failed with a ready microphone: %v", err)
	}
	if got := atomic.LoadInt32(factoryCalls); got != 1 {
		t.Fatalf("pipeline factory calls = %d, want 1", got)
	}

	report := app.GetAudioStatus()
	if ready, _ := report["asr_ready"].(bool); !ready {
		t.Fatal("asr_ready is false after initialization with a ready microphone")
	}
	if readiness, _ := report["asr_readiness"].(string); readiness != "local-initialized" {
		t.Fatalf("asr_readiness = %q, want local-initialized", readiness)
	}
	mic := micStatusFromReport(t, report)
	if state, _ := mic["state"].(string); state != string(MicStateReady) {
		t.Fatalf("mic.state = %q, want ready", state)
	}
	if role, _ := mic["role"].(string); role != "console" {
		t.Fatalf("mic.role = %q, want console", role)
	}
}

func TestMicrophoneDisappearingDuringInitializationAbortsGeneration(t *testing.T) {
	useMicProbe(t, queuedMicProbe(
		newMicrophoneStatus(MicStateReady, "", 1),
		newMicrophoneStatus(MicStateUnavailable, "device removed during load", 0),
	))
	app, factoryCalls, pipelines := newMicTestApp(t)

	err := app.applyAudioSettings(desiredWithBackend(audio.Settings{ASREnabled: true, ASREngine: "whisper-base"}))
	if err == nil {
		t.Fatal("activation succeeded after the microphone disappeared during initialization")
	}
	if got := atomic.LoadInt32(factoryCalls); got != 1 {
		t.Fatalf("pipeline factory calls = %d, want 1", got)
	}

	if _, applied := app.audio.Manager().AppliedSettings(); applied {
		t.Fatal("an ASR generation remained active after the microphone disappeared")
	}
	if desired := app.audio.Manager().Desired(); desired.ASREnabled {
		t.Fatal("desired state still requests ASR after the abort")
	}

	deadline := time.Now().Add(2 * time.Second)
	for {
		var shutdowns int32
		if len(*pipelines) > 0 {
			shutdowns = atomic.LoadInt32(&(*pipelines)[0].shutdowns)
		}
		if shutdowns > 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the aborted pipeline was not shut down")
		}
		time.Sleep(10 * time.Millisecond)
	}

	report := app.GetAudioStatus()
	if ready, _ := report["asr_ready"].(bool); ready {
		t.Fatal("asr_ready is true after the abort")
	}
	mic := micStatusFromReport(t, report)
	if state, _ := mic["state"].(string); state != string(MicStateUnavailable) {
		t.Fatalf("mic.state = %q, want unavailable", state)
	}
}

func TestMicrophoneRetrySucceedsWithoutRestart(t *testing.T) {
	useMicProbe(t, queuedMicProbe(
		newMicrophoneStatus(MicStateNoDevice, "", 0),
		newMicrophoneStatus(MicStateReady, "", 1),
	))
	app, factoryCalls, _ := newMicTestApp(t)

	first := app.applyAudioSettings(desiredWithBackend(audio.Settings{ASREnabled: true, ASREngine: "whisper-base"}))
	if first == nil {
		t.Fatal("first activation unexpectedly succeeded without a microphone")
	}
	if first.Error() != MicrophoneNoDeviceMessage {
		t.Fatalf("first error = %q, want exact no-device copy", first.Error())
	}
	if got := atomic.LoadInt32(factoryCalls); got != 0 {
		t.Fatalf("factory calls after refusal = %d, want 0", got)
	}

	if err := app.applyAudioSettings(desiredWithBackend(audio.Settings{ASREnabled: true, ASREngine: "whisper-base"})); err != nil {
		t.Fatalf("retry failed after the microphone became available: %v", err)
	}
	if got := atomic.LoadInt32(factoryCalls); got != 1 {
		t.Fatalf("factory calls after retry = %d, want 1", got)
	}
	report := app.GetAudioStatus()
	if ready, _ := report["asr_ready"].(bool); !ready {
		t.Fatal("asr_ready is false after a successful retry")
	}
}

// TestStatusReadNeverOpensCaptureStream is the regression guard for the audio
// disturbance reported at startup and on tab return: reading status must use the
// presence probe, which never activates the endpoint or initializes a capture
// client. It covers both status branches, including the pre-manager one.
func TestStatusReadNeverOpensCaptureStream(t *testing.T) {
	var presenceCalls, initializeCalls int32
	ready := newMicrophoneStatus(MicStateReady, "", 1)
	useMicProbe(t, func() MicrophoneStatus {
		atomic.AddInt32(&initializeCalls, 1)
		return ready
	})
	useMicPresenceProbe(t, func() MicrophoneStatus {
		atomic.AddInt32(&presenceCalls, 1)
		return ready
	})

	app, factoryCalls, _ := newMicTestApp(t)
	for i := 0; i < 3; i++ {
		if report := app.GetAudioStatus(); report == nil {
			t.Fatal("GetAudioStatus returned nil")
		}
	}
	// The manager-less branch: status must still answer, and still without
	// opening a capture stream.
	if report := (&App{}).GetAudioStatus(); report == nil {
		t.Fatal("GetAudioStatus without a manager returned nil")
	}

	if got := atomic.LoadInt32(&initializeCalls); got != 0 {
		t.Fatalf("status reads ran the initialize-level probe %d time(s); a status read must not open a capture stream", got)
	}
	if got := atomic.LoadInt32(&presenceCalls); got != 4 {
		t.Fatalf("presence probe calls = %d, want 4", got)
	}
	if got := atomic.LoadInt32(factoryCalls); got != 0 {
		t.Fatalf("status reads built %d generation(s); a status read must not activate audio", got)
	}
}

// TestExplicitActivationStillOpensCaptureStream keeps the other half of the
// contract: the initialize-level probe (the one that proves capture works) is
// what gates ASR activation, so it must still run exactly twice per apply.
func TestExplicitActivationStillOpensCaptureStream(t *testing.T) {
	var initializeCalls, presenceCalls int32
	ready := newMicrophoneStatus(MicStateReady, "", 1)
	useMicProbe(t, func() MicrophoneStatus {
		atomic.AddInt32(&initializeCalls, 1)
		return ready
	})
	useMicPresenceProbe(t, func() MicrophoneStatus {
		atomic.AddInt32(&presenceCalls, 1)
		return ready
	})

	app, factoryCalls, _ := newMicTestApp(t)
	if err := app.applyAudioSettings(desiredWithBackend(audio.Settings{ASREnabled: true, ASREngine: "whisper-base"})); err != nil {
		t.Fatalf("ASR activation failed with a ready microphone: %v", err)
	}
	if got := atomic.LoadInt32(&initializeCalls); got != 2 {
		t.Fatalf("initialize-level probe calls during apply = %d, want 2 (before and after the build)", got)
	}
	if got := atomic.LoadInt32(&presenceCalls); got != 0 {
		t.Fatalf("apply used the presence-only probe %d time(s); activation must prove capture works", got)
	}
	if got := atomic.LoadInt32(factoryCalls); got != 1 {
		t.Fatalf("pipeline factory calls = %d, want 1", got)
	}
}

func TestPersistedSelectionDoesNotLoadWithoutApply(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateReady, "", 1)))
	app, factoryCalls, _ := newMicTestApp(t)

	// Startup restores settings and reads status; it must not build a generation.
	_ = app.GetAudioStatus()
	time.Sleep(50 * time.Millisecond)

	if got := atomic.LoadInt32(factoryCalls); got != 0 {
		t.Fatalf("factory calls before any explicit apply = %d, want 0", got)
	}
	if _, applied := app.audio.Manager().AppliedSettings(); applied {
		t.Fatal("a generation was applied without an explicit activation request")
	}
}

func TestTTSEngineIsNotBlockedByMicrophone(t *testing.T) {
	useMicProbe(t, constantMicProbe(newMicrophoneStatus(MicStateNoDevice, "", 0)))
	app, factoryCalls, _ := newMicTestApp(t)

	if err := app.applyAudioSettings(desiredWithBackend(audio.Settings{TTSEnabled: true, TTSEngine: "piper"})); err != nil {
		t.Fatalf("TTS activation was blocked by the microphone gate: %v", err)
	}
	if got := atomic.LoadInt32(factoryCalls); got != 1 {
		t.Fatalf("factory calls = %d, want 1", got)
	}
}
