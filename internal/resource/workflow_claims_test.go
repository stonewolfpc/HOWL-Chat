package resource

import (
	"path/filepath"
	"testing"
)

func TestWorkflowRetainedUpscalerIsClaimed(t *testing.T) {
	lease := NewWorkflowLease("engine-a")
	key := NewComponentKey("engine-a:42", "upscaler", "/models/upscaler.pth", "f16", []string{"vae.safetensors"})

	if lease.Verify(key) {
		t.Fatal("unclaimed upscaler verified as resident")
	}
	claim, err := lease.Acquire(key)
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if !lease.Verify(key) {
		t.Fatal("stage-2 upscaler claim is not visible to stage 4")
	}
	if got := lease.ClaimCount(key); got != 1 {
		t.Fatalf("claim count = %d, want 1", got)
	}

	claim.Release()
	if lease.Verify(key) {
		t.Fatal("released upscaler claim still verified as resident")
	}
	if got := lease.ClaimCount(key); got != 0 {
		t.Fatalf("claim count after release = %d, want 0", got)
	}
}

func TestSharedComponentRefcountAcrossStages(t *testing.T) {
	lease := NewWorkflowLease("engine-a")
	key := NewComponentKey("engine-a:42", "upscaler", "/models/upscaler.pth", "f16", nil)

	stage2, err := lease.Acquire(key)
	if err != nil {
		t.Fatalf("stage 2 Acquire: %v", err)
	}
	stage4, err := lease.Acquire(key)
	if err != nil {
		t.Fatalf("stage 4 Acquire: %v", err)
	}
	if got := lease.ClaimCount(key); got != 2 {
		t.Fatalf("shared component count = %d, want 2", got)
	}

	stage2.Release()
	if got := lease.ClaimCount(key); got != 1 {
		t.Fatalf("count after first release = %d, want 1", got)
	}
	if !lease.Verify(key) {
		t.Fatal("component lost its remaining stage-4 claim")
	}

	stage4.Release()
	stage4.Release() // idempotent: a second release cannot go negative
	if got := lease.ClaimCount(key); got != 0 {
		t.Fatalf("count after final release = %d, want 0", got)
	}
}

// Characterization: releasing the parent closes ownership; it does not stop a
// stage that is already running. Cancellation needs a separate execution scope.
func TestWorkflowLeaseReleaseDoesNotStopActiveStage(t *testing.T) {
	lease := NewWorkflowLease("engine-a")
	started := make(chan struct{})
	finish := make(chan struct{})
	done := make(chan struct{})
	go func() {
		close(started)
		<-finish
		close(done)
	}()
	<-started

	lease.Release()
	select {
	case <-done:
		t.Fatal("parent release stopped the active stage")
	default:
	}

	close(finish)
	<-done
}

func TestWorkflowParentReleaseIsIdempotent(t *testing.T) {
	lease := NewWorkflowLease("engine-a")
	key := NewComponentKey("engine-a:42", "upscaler", "/models/upscaler.pth", "f16", nil)
	claim, err := lease.Acquire(key)
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}

	lease.Release()
	lease.Release()
	claim.Release()
	claim.Release()

	if !lease.Released() || lease.ActiveClaimCount() != 0 {
		t.Fatalf("released=%t active claims=%d", lease.Released(), lease.ActiveClaimCount())
	}
	if _, err := lease.Acquire(key); err == nil {
		t.Fatal("released parent accepted a new claim")
	}
}

func TestWorkflowReusesEngineAcrossStages(t *testing.T) {
	// The same component is reached through a different spelling and companion
	// order. The key must remain stable across stages.
	stage2 := NewComponentKey("engine-a:42", "upscaler", filepath.Join("models", "upscaler.pth"), "F16",
		[]string{"vae.safetensors", "clip_l.safetensors"})
	stage4 := NewComponentKey("engine-a:42", "upscaler", filepath.Join("models", ".", "upscaler.pth"), "f16",
		[]string{"clip_l.safetensors", "vae.safetensors"})
	if stage2 != stage4 {
		t.Fatalf("same resident upscaler produced different keys: %v vs %v", stage2, stage4)
	}

	reloaded := NewComponentKey("engine-a:43", "upscaler", filepath.Join("models", "upscaler.pth"), "f16",
		[]string{"vae.safetensors", "clip_l.safetensors"})
	if stage2 == reloaded {
		t.Fatal("a replacement engine process reused the previous component key")
	}
}
