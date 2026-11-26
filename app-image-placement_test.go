package main

import (
	"strings"
	"testing"

	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/image/types"
)

func TestImagePlacementDivergenceUsesEngineReport(t *testing.T) {
	opts := types.ModelLoadOptions{
		MaxVRAM:           intPtrForTest(3),
		Backend:           imageConstrainedBackend,
		OffloadToCPU:      true,
		AutoPlacement:     true,
		AutoPlacementMode: "cpu-diffusion",
	}
	actual := sdcpp.Placement{
		Known:  true,
		VRAMMB: 160,
		RAMMB:  4937,
		Components: map[string]string{
			"diffusion_model": "RAM",
			"text_encoders":   "RAM",
			"vae":             "VRAM",
		},
	}
	if got := imagePlacementDivergence(opts, actual, 5000, 819); got != "" {
		t.Fatalf("matching engine placement reported divergence: %s", got)
	}

	actual.Components["diffusion_model"] = "VRAM"
	if got := imagePlacementDivergence(opts, actual, 5000, 819); !strings.Contains(got, "diffusion_model in VRAM") {
		t.Fatalf("planned CPU diffusion with actual VRAM was not reported: %q", got)
	}

	cpuOnly := opts
	cpuOnly.AutoPlacementMode = "cpu-only"
	cpuOnly.Backend = "cpu"
	cpuOnly.OffloadToCPU = false
	actual = sdcpp.Placement{Known: true, VRAMMB: 100, Components: map[string]string{"vae": "VRAM"}}
	if got := imagePlacementDivergence(cpuOnly, actual, 5000, 819); !strings.Contains(got, "CPU-only") {
		t.Fatalf("CPU-only plan with VRAM use was not reported: %q", got)
	}

	overBudget := sdcpp.Placement{Known: true, VRAMMB: 4096, Components: map[string]string{}}
	if got := imagePlacementDivergence(opts, overBudget, 5000, 819); !strings.Contains(got, "above the planned") {
		t.Fatalf("actual VRAM above the budget was not reported: %q", got)
	}

	if got := imagePlacementDivergence(opts, sdcpp.Placement{}, 5000, 819); !strings.Contains(got, "unverified") {
		t.Fatalf("missing engine report must say unverified: %q", got)
	}

	measured := sdcpp.Placement{Process: map[string]sdcpp.Footprint{
		"load":       {Known: true, WorkingSetMB: 6200},
		"generation": {Known: true, WorkingSetMB: 6900},
	}}
	got := imagePlacementDivergence(opts, measured, 5000, 819)
	if !strings.Contains(got, "6200 MB after load") || !strings.Contains(got, "unverified") {
		t.Fatalf("measured footprint must be reported without claiming component residency: %q", got)
	}
	summary := imagePlacementSummary(imagePlanDescription(opts), measured, got)
	if !strings.Contains(summary, "6900 MB after generation") {
		t.Fatalf("summary omitted measured post-generation footprint: %q", summary)
	}
}

func TestAutoPlacementIsCarriedToTheNextGenerationOnlyForTheSameBudget(t *testing.T) {
	actual := types.ModelLoadOptions{
		MaxVRAM:           intPtrForTest(3),
		Backend:           imageConstrainedBackend,
		OffloadToCPU:      true,
		VAETiling:         true,
		AutoPlacement:     true,
		AutoPlacementMode: "cpu-diffusion",
	}
	req := types.ModelLoadOptions{MaxVRAM: intPtrForTest(3)}
	got := preserveImageLoadedOptions(req, actual)
	if !got.AutoPlacement || got.Backend != imageConstrainedBackend || !got.OffloadToCPU || !got.VAETiling {
		t.Fatalf("automatic placement was not carried into the resident generation: %+v", got)
	}

	changed := types.ModelLoadOptions{MaxVRAM: intPtrForTest(6)}
	got = preserveImageLoadedOptions(changed, actual)
	if got.AutoPlacement || got.Backend != "" {
		t.Fatalf("a changed budget must be re-planned, not carried: %+v", got)
	}
}

func intPtrForTest(v int) *int { return &v }
