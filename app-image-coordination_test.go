package main

import (
	"strings"
	"testing"

	"howl-chat/internal/image/types"
	"howl-chat/internal/resource"
)

// Cross-workload coordination: a resident workload is a handoff decision, not a
// reason to quietly run the requested model on the CPU. The refusal must name
// the blocker so the user-approved reclaim path can act on it.
func TestResidentBlockerIsNamedRatherThanDegradedToCPU(t *testing.T) {
	path := testImageSafetensors(t, 4897, 1559, 159)
	app := testImageApp(6901, 8192)
	// Chat holds enough of the card that the image model cannot fit beside it.
	app.admission = resource.NewAuthority(
		func() (resource.HardwareSnapshot, error) {
			return resource.HardwareSnapshot{
				GPUModel: "Test GPU", TotalVRAMBytes: 8192 * 1024 * 1024, FreeVRAMBytes: 1860 * 1024 * 1024,
				UsedVRAMBytes: (8192 - 1860) * 1024 * 1024, TotalRAMBytes: 32 * 1024 * 1024 * 1024,
				FreeRAMBytes: 16 * 1024 * 1024 * 1024, VRAMLive: true, VRAMSource: "test-live", GPUAvailable: true,
			}, nil
		},
		func() []resource.ModelState {
			return []resource.ModelState{{
				Type: resource.ResourceLLM, Name: "chat", Loaded: true,
				GPUMemoryMB: 5209, SystemRAMMB: 729,
			}}
		},
	)
	gen := types.ImageOptions{ModelName: "illustrious", Width: 512, Height: 512, BatchCount: 1, Steps: 1}

	_, _, chosen, err := app.admitImageLoad(types.ModelLoadOptions{ModelName: "illustrious", ModelPath: path}, gen, true)
	if err == nil {
		t.Fatal("an image model must not start on the CPU just because chat holds the GPU")
	}
	if !strings.Contains(err.Error(), "chat") {
		t.Fatalf("the refusal must name the blocking workload: %v", err)
	}
	if chosen.AutoPlacement || chosen.AutoPlacementMode != "" {
		t.Fatalf("a blocked start must not acquire a bounded placement: %+v", chosen)
	}
}

// The same model with no resident blocker still gets the M5 treatment.
func TestBoundedPlacementStillAppliesWithNoBlocker(t *testing.T) {
	path := testImageSafetensors(t, 4897, 4667, 159)
	app := testImageApp(6901, 8192)
	gen := types.ImageOptions{ModelName: "hidream", Width: 512, Height: 512, BatchCount: 1, Steps: 1}

	_, lease, chosen, err := app.admitImageLoad(types.ModelLoadOptions{ModelName: "hidream", ModelPath: path}, gen, true)
	if err != nil {
		t.Fatalf("bounded placement should still start an oversized model: %v", err)
	}
	lease.Release()
	if !chosen.AutoPlacement {
		t.Fatal("expected a bounded placement for a model larger than the card")
	}
}
