package main

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	"howl-chat/internal/audio"
)

// The capture port is a precondition check, not a capture capability. If it ever
// runs the initialize-level probe, Howl opens a shared-mode capture client just
// to answer "is a microphone there?", which flips Bluetooth headsets from A2DP
// to HFP and changes every other application's playback until the client is
// released. These tests therefore fail on invocation of that seam rather than on
// the shape of the resulting status.

// fatalInitializeSeam replaces the initialize-level probe with one that fails the
// test if any path asks for it, and counts invocations for the assertion after
// the calls under test.
func fatalInitializeSeam(t *testing.T) *int32 {
	t.Helper()
	var calls int32
	prev := microphoneProbe
	microphoneProbe = func() MicrophoneStatus {
		atomic.AddInt32(&calls, 1)
		t.Errorf("initialize-level capture probe invoked: the capture port and every status/apply path must only check presence")
		return newMicrophoneStatus(MicStateCaptureInitFailed, "initialize seam must not run", 0)
	}
	t.Cleanup(func() { microphoneProbe = prev })
	return &calls
}

// presenceSeam replaces the presence-level probe with a fixed answer.
func presenceSeam(t *testing.T, status MicrophoneStatus) {
	t.Helper()
	prev := microphonePresenceProbe
	microphonePresenceProbe = func() MicrophoneStatus { return status }
	t.Cleanup(func() { microphonePresenceProbe = prev })
}

func readyPresenceStatus() MicrophoneStatus {
	return newMicrophoneStatus(MicStateReady,
		"the default console capture endpoint is present and active in Windows; capture was not opened (no stream, no recording)", 1)
}

func TestCapturePortNeverInitializesTheMicrophone(t *testing.T) {
	initializations := fatalInitializeSeam(t)
	presenceSeam(t, readyPresenceStatus())

	var factoryCalls int32
	manager := audio.NewManager(func(_ context.Context, desired audio.Settings, candidate *audio.Candidate) error {
		atomic.AddInt32(&factoryCalls, 1)
		if candidate != nil {
			candidate.Pipeline = &micTestPipeline{}
			if desired.ASREnabled {
				candidate.ASR = audio.CapabilityResult{Ready: true}
			}
		}
		return nil
	})
	t.Cleanup(func() { _ = manager.Close(context.Background()) })
	app := newAudioTestAppWithManager(manager)

	// 1. The capture port itself.
	status, usable, detail := app.captureProbe()
	if !usable {
		t.Fatalf("capture port reported unusable for a present-and-enabled endpoint: %#v (%s)", status, detail)
	}

	// 2. The status read that runs on every page load and tab return.
	report := app.GetAudioStatus()
	mic, ok := report["mic"].(map[string]interface{})
	if !ok {
		t.Fatalf("GetAudioStatus has no mic record: %#v", report["mic"])
	}
	if got := mic["state"]; got != string(MicStateReady) {
		t.Fatalf("status-read mic state = %v, want %q", got, MicStateReady)
	}
	// The rendered copy must describe what the presence probe establishes and
	// nothing more: no claim of initialization, of a device in use, or of
	// recording working.
	message, _ := mic["message"].(string)
	if !strings.Contains(message, "present and enabled") {
		t.Fatalf("ready message = %q, want it to say the default microphone is present and enabled", message)
	}
	for _, forbidden := range []string{"initialized", "in use", "will work", "recording will"} {
		if strings.Contains(strings.ToLower(message), forbidden) {
			t.Fatalf("ready message = %q, must not claim %q", message, forbidden)
		}
	}

	// 3. The explicit apply that used to open the endpoint twice.
	if err := app.ApplyAudioSettings(desiredWithBackend(audio.Settings{ASREnabled: true, ASREngine: "whisper-base"})); err != nil {
		t.Fatalf("ApplyAudioSettings with a present-and-enabled endpoint: %v", err)
	}
	if got := atomic.LoadInt32(&factoryCalls); got == 0 {
		t.Fatal("apply never built a pipeline; the successful-apply half of this test proved nothing")
	}

	if got := atomic.LoadInt32(initializations); got != 0 {
		t.Fatalf("initialize-level probe ran %d time(s) from the capture port, GetAudioStatus or ApplyAudioSettings", got)
	}
}

func TestPresenceNoDeviceStillBlocksASRApply(t *testing.T) {
	initializations := fatalInitializeSeam(t)
	presenceSeam(t, newMicrophoneStatus(MicStateNoDevice,
		"GetDefaultAudioEndpoint(eCapture, eConsole) returned E_NOTFOUND: no default console capture endpoint exists", 0))

	var factoryCalls int32
	manager := audio.NewManager(func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		atomic.AddInt32(&factoryCalls, 1)
		if candidate != nil {
			candidate.ASR = audio.CapabilityResult{Ready: true}
		}
		return nil
	})
	t.Cleanup(func() { _ = manager.Close(context.Background()) })
	app := newAudioTestAppWithManager(manager)

	err := app.applyAudioSettings(desiredWithBackend(audio.Settings{ASREnabled: true, ASREngine: "whisper-base"}))
	if err == nil {
		t.Fatal("ASR apply succeeded with no capture endpoint")
	}
	var micErr *MicrophoneUnavailableError
	if !errors.As(err, &micErr) {
		t.Fatalf("apply error = %T (%v), want *MicrophoneUnavailableError", err, err)
	}
	if micErr.Status.State != MicStateNoDevice {
		t.Fatalf("refusal state = %q, want %q", micErr.Status.State, MicStateNoDevice)
	}
	if micErr.Status.Message != MicrophoneNoDeviceMessage {
		t.Fatalf("refusal message = %q, want the frozen %q", micErr.Status.Message, MicrophoneNoDeviceMessage)
	}
	if err.Error() != MicrophoneNoDeviceMessage {
		t.Fatalf("error text = %q, want the frozen %q", err.Error(), MicrophoneNoDeviceMessage)
	}
	if got := atomic.LoadInt32(&factoryCalls); got != 0 {
		t.Fatalf("pipeline factory ran %d time(s); ASR must be refused before allocation", got)
	}
	if got := atomic.LoadInt32(initializations); got != 0 {
		t.Fatalf("initialize-level probe ran %d time(s) on the refusal path", got)
	}
}
