package audio

import (
	"context"
	"sync/atomic"
	"testing"
)

// Re-applying the settings that are already installed must not build a
// replacement generation. Building one would re-admit models the active
// generation owns; the active registration refuses the duplicate, so the
// apply failed against itself.
func TestReapplyInstalledSettingsDoesNotRebuild(t *testing.T) {
	var calls int32
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		atomic.AddInt32(&calls, 1)
		candidate.Pipeline = appliedPipeline{}
		candidate.ASR = CapabilityResult{Ready: desired.ASREnabled}
		candidate.TTS = CapabilityResult{Ready: desired.TTSEnabled}
		return nil
	})
	defer func() { _ = manager.Close(context.Background()) }()

	installed := Settings{ASREnabled: true, ASREngine: "qwen3-asr-1.7b-local"}
	if _, _, err := manager.Apply(context.Background(), installed); err != nil {
		t.Fatalf("first Apply returned %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("factory calls after first Apply = %d, want 1", got)
	}

	status, superseded, err := manager.Apply(context.Background(), installed)
	if err != nil {
		t.Fatalf("re-Apply of the installed configuration returned %v", err)
	}
	if superseded {
		t.Fatal("re-Apply of the installed configuration was reported as superseded")
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("factory calls after re-Apply = %d, want 1", got)
	}
	if !status.ASR.Ready {
		t.Fatalf("re-Apply lost ASR readiness: %+v", status)
	}
	if status.DesiredGen != status.AppliedGen {
		t.Fatalf("re-Apply did not converge: desired=%d applied=%d", status.DesiredGen, status.AppliedGen)
	}
	if status.Builds != 0 {
		t.Fatalf("Builds after re-Apply = %d, want 0", status.Builds)
	}
	got, ok := manager.AppliedSettings()
	if !ok || got != installed {
		t.Fatalf("AppliedSettings = %+v (ok=%v), want %+v", got, ok, installed)
	}
}

// A degraded install is not a healthy instance to reuse. An identical request
// must be allowed to rebuild, otherwise a transient failure would pin the
// capability down until the settings changed.
func TestReapplyDegradedActivationRetries(t *testing.T) {
	var calls int32
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		n := atomic.AddInt32(&calls, 1)
		candidate.Pipeline = appliedPipeline{}
		candidate.ASR = CapabilityResult{Ready: n > 1}
		return nil
	})
	defer func() { _ = manager.Close(context.Background()) }()

	desired := Settings{ASREnabled: true, ASREngine: "qwen3-asr-1.7b-local"}
	if _, _, err := manager.Apply(context.Background(), desired); err == nil {
		t.Fatal("first Apply succeeded although ASR was not ready")
	}
	if _, _, err := manager.Apply(context.Background(), desired); err != nil {
		t.Fatalf("retry of a degraded install returned %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("factory calls = %d, want 2 (degraded install must be retried)", got)
	}
}

// A changed configuration is a genuine replacement. It must still build, so
// admission checks for the new engines continue to run.
func TestChangedSettingsStillRebuild(t *testing.T) {
	var calls int32
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		atomic.AddInt32(&calls, 1)
		candidate.Pipeline = appliedPipeline{}
		candidate.ASR = CapabilityResult{Ready: desired.ASREnabled}
		candidate.TTS = CapabilityResult{Ready: desired.TTSEnabled}
		return nil
	})
	defer func() { _ = manager.Close(context.Background()) }()

	first := Settings{ASREnabled: true, ASREngine: "qwen3-asr-1.7b-local"}
	second := Settings{ASREnabled: true, ASREngine: "whisper-base"}
	if _, _, err := manager.Apply(context.Background(), first); err != nil {
		t.Fatalf("first Apply returned %v", err)
	}
	if _, _, err := manager.Apply(context.Background(), second); err != nil {
		t.Fatalf("changed Apply returned %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("factory calls = %d, want 2 for a changed configuration", got)
	}
	got, ok := manager.AppliedSettings()
	if !ok || got != second {
		t.Fatalf("AppliedSettings = %+v (ok=%v), want %+v", got, ok, second)
	}
}

// Disabling both capabilities must still retire the active generation; the
// idempotence rule must not treat a disabled request as satisfied by a running
// install.
func TestIdempotenceDoesNotBlockRetirement(t *testing.T) {
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		candidate.Pipeline = appliedPipeline{}
		candidate.ASR = CapabilityResult{Ready: desired.ASREnabled}
		return nil
	})
	defer func() { _ = manager.Close(context.Background()) }()

	if _, _, err := manager.Apply(context.Background(), Settings{ASREnabled: true, ASREngine: "whisper-base"}); err != nil {
		t.Fatalf("Apply returned %v", err)
	}
	status, _, err := manager.Apply(context.Background(), Settings{})
	if err != nil {
		t.Fatalf("disable Apply returned %v", err)
	}
	if status.ASR.Ready || status.ASR.Requested {
		t.Fatalf("disable Apply kept ASR visible: %+v", status)
	}
	if _, ok := manager.AppliedSettings(); ok {
		t.Fatal("AppliedSettings still reports a generation after both capabilities were disabled")
	}
}
