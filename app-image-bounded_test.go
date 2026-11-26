package main

import (
	"testing"

	"howl-chat/internal/image/types"
)

// HiDream-class: the diffusion file plus its text encoders is larger than the
// whole card. Refusing is honest but final; with placement left on Auto, Howl
// can bound the start instead — weights in RAM, device spend capped by the
// budget admission just computed.

func hidreamShapedFixture(t *testing.T) string {
	t.Helper()
	return testImageSafetensors(t, 4897, 4667, 159)
}

func TestBoundedPlacementStartsAModelLargerThanTheCard(t *testing.T) {
	path := hidreamShapedFixture(t)
	app := testImageApp(6901, 8192)
	gen := types.ImageOptions{ModelName: "hidream", Width: 512, Height: 512, BatchCount: 1, Steps: 1}

	// The full resident set is refused first: that refusal is the evidence for
	// why a bounded placement is needed, and it stays the reported outcome when
	// no bounded mode fits.
	decision, lease, chosen, err := app.admitImageLoad(types.ModelLoadOptions{ModelName: "hidream", ModelPath: path}, gen, true)
	if err != nil {
		t.Fatalf("a bounded placement should start this model instead of refusing it: %v", err)
	}
	if !decision.Admitted || lease == nil {
		t.Fatal("bounded admission must hold a lease")
	}
	lease.Release()
	if !chosen.AutoPlacement || chosen.AutoPlacementMode != "cpu-diffusion" {
		t.Fatalf("bounded mode = %q (auto=%t), want cpu-diffusion", chosen.AutoPlacementMode, chosen.AutoPlacement)
	}
	if !chosen.OffloadToCPU || chosen.Backend != imageConstrainedBackend || !chosen.VAETiling {
		t.Fatalf("bounded placement flags were not applied: %+v", chosen)
	}
	if chosen.MaxVRAM == nil || *chosen.MaxVRAM <= 0 || !chosen.AutoMaxVRAM {
		t.Fatalf("bounded placement must carry a positive automatic cap, got %v", chosen.MaxVRAM)
	}
	if got := int64(*chosen.MaxVRAM) * 1024; got > decision.Estimate.AvailableVRAMMB {
		t.Fatalf("cap %d MB exceeds the post-reserve budget %d MB", got, decision.Estimate.AvailableVRAMMB)
	}
}

func TestBoundedPlacementNeverOverridesAnExplicitBackend(t *testing.T) {
	path := hidreamShapedFixture(t)
	app := testImageApp(6901, 8192)
	gen := types.ImageOptions{ModelName: "hidream", Width: 512, Height: 512, BatchCount: 1, Steps: 1}

	_, _, _, err := app.admitImageLoad(types.ModelLoadOptions{
		ModelName: "hidream", ModelPath: path, Backend: "cuda",
	}, gen, true)
	if err == nil {
		t.Fatal("an explicit backend must not be rewritten into a bounded placement")
	}
}

func TestBoundedPlacementFallsBackToCPUBelowTheMeasuredFloor(t *testing.T) {
	path := hidreamShapedFixture(t)
	// 1500 MB free leaves roughly 860 MB after reserve: below the measured
	// GPU-assisted process need, so the honest bounded mode is CPU-only.
	app := testImageApp(1500, 8192)
	gen := types.ImageOptions{ModelName: "hidream", Width: 512, Height: 512, BatchCount: 1, Steps: 1}

	decision, lease, chosen, err := app.admitImageLoad(types.ModelLoadOptions{ModelName: "hidream", ModelPath: path}, gen, true)
	if err != nil {
		t.Fatalf("CPU-only bounded placement should start instead of refusing: %v", err)
	}
	lease.Release()
	if !decision.Admitted || chosen.AutoPlacementMode != "cpu-only" || chosen.Backend != "cpu" {
		t.Fatalf("expected cpu-only bounded placement, got mode=%q backend=%q", chosen.AutoPlacementMode, chosen.Backend)
	}
}
