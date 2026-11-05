package audio

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

type appliedPipeline struct{}

func (appliedPipeline) ProcessAudioMessage(context.Context, types.AudioMessage) (*types.AudioResponse, error) {
	return nil, nil
}
func (appliedPipeline) GetStatus() types.PipelineStatus { return types.PipelineStatus{} }
func (appliedPipeline) GetTTSVoices() []types.Voice     { return nil }
func (appliedPipeline) Shutdown(context.Context) error  { return nil }

func appliedFactory(_ context.Context, desired Settings, candidate *Candidate) error {
	candidate.Pipeline = appliedPipeline{}
	candidate.ASR = CapabilityResult{Ready: desired.ASREnabled}
	candidate.TTS = CapabilityResult{Ready: desired.TTSEnabled}
	return nil
}

// AppliedSettings must describe the installed generation, not the newest
// request, or the panel could confirm a configuration that was never activated.
func TestAppliedSettingsReportsActiveGeneration(t *testing.T) {
	manager := NewManager(appliedFactory)
	defer func() { _ = manager.Close(context.Background()) }()

	if _, ok := manager.AppliedSettings(); ok {
		t.Fatal("AppliedSettings reported a generation before any apply")
	}

	installed := Settings{ASREnabled: true, ASREngine: "qwen3-asr-1.7b-local", TTSEngine: "piper"}
	if _, _, err := manager.Apply(context.Background(), installed); err != nil {
		t.Fatalf("Apply returned %v", err)
	}

	got, ok := manager.AppliedSettings()
	if !ok {
		t.Fatal("AppliedSettings reported no generation after a successful apply")
	}
	if got != installed {
		t.Fatalf("AppliedSettings = %+v, want %+v", got, installed)
	}
}

// A failed activation must not publish itself as the effective configuration.
func TestAppliedSettingsIgnoresFailedActivation(t *testing.T) {
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		candidate.ASR = CapabilityResult{Ready: false}
		candidate.TTS = CapabilityResult{Ready: false}
		return nil
	})
	defer func() { _ = manager.Close(context.Background()) }()

	_, _, err := manager.Apply(context.Background(), Settings{ASREnabled: true, ASREngine: "whisper-base"})
	if err == nil {
		t.Fatal("Apply returned nil for an activation with no ready capability")
	}
	if settings, ok := manager.AppliedSettings(); ok {
		t.Fatalf("AppliedSettings reported %+v after a failed activation", settings)
	}
}

// A caller queued on applyMu has already advanced DesiredGen, so it must be
// visible in Status even before any build for it begins.
func TestStatusExposesQueuedGeneration(t *testing.T) {
	started := make(chan struct{})
	gate := make(chan struct{})
	var calls int32
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		if atomic.AddInt32(&calls, 1) == 1 {
			close(started)
			<-gate
		}
		candidate.Pipeline = appliedPipeline{}
		candidate.ASR = CapabilityResult{Ready: desired.ASREnabled}
		candidate.TTS = CapabilityResult{Ready: desired.TTSEnabled}
		return nil
	})
	defer func() { _ = manager.Close(context.Background()) }()

	firstDone := make(chan struct{})
	go func() {
		defer close(firstDone)
		_, _, _ = manager.Apply(context.Background(), Settings{ASREnabled: true, ASREngine: "whisper-base"})
	}()
	<-started

	secondDone := make(chan struct{})
	go func() {
		defer close(secondDone)
		_, _, _ = manager.Apply(context.Background(), Settings{ASREnabled: true, ASREngine: "whisper-small"})
	}()

	deadline := time.Now().Add(2 * time.Second)
	for {
		st := manager.Status()
		if st.DesiredGen > st.AppliedGen {
			if st.Builds != 1 {
				t.Fatalf("Builds = %d while one build is gated, want 1", st.Builds)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("queued generation never became visible: %+v", st)
		}
		time.Sleep(5 * time.Millisecond)
	}

	close(gate)
	<-firstDone
	<-secondDone

	final := manager.Status()
	if final.Builds != 0 {
		t.Fatalf("Builds = %d after settlement, want 0", final.Builds)
	}
	if final.DesiredGen != final.AppliedGen {
		t.Fatalf("generations did not converge: desired=%d applied=%d", final.DesiredGen, final.AppliedGen)
	}
}

// Characterization only: Apply reports a superseded caller through its bool
// while returning a nil error, and the app-level binding discards that bool.
// This test records the behavior; it does not change the contract.
func TestAppliedSettingsSupersededGeneration(t *testing.T) {
	gate := make(chan struct{})
	var calls int32
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		if atomic.AddInt32(&calls, 1) == 1 {
			<-gate
		}
		candidate.Pipeline = appliedPipeline{}
		candidate.ASR = CapabilityResult{Ready: desired.ASREnabled}
		return nil
	})
	defer func() { _ = manager.Close(context.Background()) }()

	first := Settings{ASREnabled: true, ASREngine: "whisper-base"}
	second := Settings{ASREnabled: true, ASREngine: "whisper-small"}

	var wg sync.WaitGroup
	var firstErr error
	var firstSuperseded bool
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, firstSuperseded, firstErr = manager.Apply(context.Background(), first)
	}()

	deadline := time.Now().Add(2 * time.Second)
	for manager.Status().Builds == 0 {
		if time.Now().After(deadline) {
			t.Fatal("first build never started")
		}
		time.Sleep(5 * time.Millisecond)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _, _ = manager.Apply(context.Background(), second)
	}()

	for manager.Status().DesiredGen < 2 {
		if time.Now().After(deadline) {
			t.Fatal("second request never advanced the desired generation")
		}
		time.Sleep(5 * time.Millisecond)
	}

	close(gate)
	wg.Wait()

	if firstErr != nil {
		t.Fatalf("superseded Apply returned %v, want nil", firstErr)
	}
	if !firstSuperseded {
		t.Fatal("first Apply was not reported as superseded")
	}
	got, ok := manager.AppliedSettings()
	if !ok || got != second {
		t.Fatalf("AppliedSettings = %+v (ok=%v), want %+v", got, ok, second)
	}
	if builds := manager.Status().Builds; builds != 0 {
		t.Fatalf("Builds = %d after settlement, want 0", builds)
	}
}
