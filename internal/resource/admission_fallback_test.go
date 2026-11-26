package resource

import (
	"strings"
	"testing"

	backendtypes "howl-chat/internal/backend/types"
)

func fallbackSpec() *backendtypes.ArchitectureSpec {
	return &backendtypes.ArchitectureSpec{
		Family: "llama", Layers: 16, DModel: 2048, NumHeads: 16, NumKVHeads: 8,
		HeadDim: 128, FFNSize: 8192, ContextLength: 8192,
		FileSizeBytes: 8 * 1024 * bytesPerMB, ParameterCount: 8_000_000_000,
	}
}

func fallbackRequest() AdmissionRequest {
	return AdmissionRequest{
		Kind:        ResourceLLM,
		Name:        "big-chat",
		ContextSize: 8192,
		GPULayers:   16,
		BatchSize:   512,
		RequireGPU:  true,
		AllowManual: true,
		Arch:        fallbackSpec(),
	}
}

func TestChatAdmissionFitsInsteadOfRefusing(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(3000, 8192)), fixedLedger())

	decision, lease := authority.Admit(fallbackRequest())
	if !decision.Admitted {
		t.Fatalf("an over-budget chat load must be fitted, not refused: %s", decision.Reason)
	}
	if lease == nil {
		t.Fatal("an admitted fitted load must return a lease")
	}
	defer lease.Release()
	if decision.Adjustment == nil {
		t.Fatal("the decision must carry the fitted context/layers/batch plan")
	}
	if decision.Manual {
		t.Fatal("a plan that fits must not be reported as a manual fallback")
	}
	if decision.Adjustment.GPULayers >= 16 {
		t.Fatalf("layers were not reduced: %+v", decision.Adjustment)
	}
	if decision.Estimate.RequiredVRAMMB > decision.Estimate.AvailableVRAMMB {
		t.Fatalf("fitted plan still exceeds VRAM: required=%d available=%d",
			decision.Estimate.RequiredVRAMMB, decision.Estimate.AvailableVRAMMB)
	}
	for _, want := range []string{"needs", "Loaded with context"} {
		if !strings.Contains(decision.Warning, want) {
			t.Fatalf("warning must explain the numbers (%q missing): %s", want, decision.Warning)
		}
	}
}

func TestChatAdmissionWithoutManualStillRefusesHonestly(t *testing.T) {
	// No schedule fits: the weights alone exceed the remaining RAM, so a
	// preflight without permission to force a manual load reports the blocker.
	authority := NewAuthority(fixedHardware(testSnapshot(3000, 8192)), fixedLedger())
	req := fallbackRequest()
	req.AllowManual = false
	req.Arch = fallbackSpec()
	req.Arch.FileSizeBytes = 24 * 1024 * bytesPerMB
	req.Arch.ParameterCount = 24_000_000_000

	decision, lease := authority.Admit(req)
	if decision.Admitted {
		t.Fatal("a preflight without manual permission must still report the true blocker")
	}
	if lease != nil {
		t.Fatal("a refused request must not return a lease")
	}
	if decision.Estimate.RequiredVRAMMB == 0 || decision.Estimate.AvailableVRAMMB == 0 {
		t.Fatalf("refusal must carry measurable numbers: %+v", decision.Estimate)
	}
}

func TestChatAdmissionManualFallbackWhenVRAMIsNotLive(t *testing.T) {
	snap := testSnapshot(3000, 8192)
	snap.VRAMLive = false
	snap.VRAMSource = "registry"
	authority := NewAuthority(fixedHardware(snap), fixedLedger())

	decision, lease := authority.Admit(fallbackRequest())
	if !decision.Admitted {
		t.Fatalf("a user-initiated load must not be locked out when telemetry is not live: %s", decision.Reason)
	}
	if lease == nil {
		t.Fatal("manual fallback must still hold a lease")
	}
	defer lease.Release()
	if !decision.Manual {
		t.Fatal("a load without verified VRAM must be reported as a manual fallback")
	}
	if decision.Adjustment == nil || decision.Adjustment.GPULayers != 0 {
		t.Fatalf("manual fallback must be CPU-first: %+v", decision.Adjustment)
	}
	if !strings.Contains(decision.Warning, "could not be read live") {
		t.Fatalf("warning must name the telemetry failure: %s", decision.Warning)
	}
	if !strings.Contains(decision.Warning, "source: registry") {
		t.Fatalf("warning must name the fallback source: %s", decision.Warning)
	}
}

func TestChatAdmissionManualFallbackStillDoesNotEvictAnotherModel(t *testing.T) {
	loaded := []ModelState{{Type: ResourceImage, Name: "sdxl", Loaded: true, GPUMemoryMB: 6000}}
	authority := NewAuthority(fixedHardware(testSnapshot(8000, 24000)), fixedLedger(loaded...))

	decision, lease := authority.Admit(fallbackRequest())
	if decision.Admitted {
		t.Fatal("a manual fallback must not silently start beside a blocking model")
	}
	if lease != nil {
		t.Fatal("a refused blocked load must not hold a lease")
	}
	if !strings.Contains(decision.Reason, "sdxl") {
		t.Fatalf("refusal must name the blocking model: %s", decision.Reason)
	}
}
