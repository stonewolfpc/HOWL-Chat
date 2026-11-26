package resource

import (
	"testing"

	"howl-chat/internal/backend/types"
)

// qwen35ish is a 33-layer, 4096-wide model with grouped KV heads, the shape the
// audit's 9B Q4 model reported.
func qwen35ish() *types.ArchitectureSpec {
	return &types.ArchitectureSpec{
		Family:                 "qwen35",
		Layers:                 33,
		DModel:                 4096,
		FFNSize:                14336,
		NumHeads:               32,
		NumKVHeads:             8,
		HeadDim:                128,
		FileSizeBytes:          5366 * 1024 * 1024,
		KVCacheBytesPerElement: 2,
	}
}

// An idle card must place layers, not zero. This is the exact regression the
// milestone exists to fix: before the planner, a request with gpu_layers 0 (or
// auto) could only ever clamp downward.
func TestPlanChatLayersUsesIdleCard(t *testing.T) {
	plan := PlanChatLayers(qwen35ish(), 4096, 512, 2, 7200, 30000)
	if !plan.Fits {
		t.Fatalf("plan does not fit an idle card: %+v", plan)
	}
	if plan.GPULayers <= 0 {
		t.Fatalf("idle card produced %d GPU layers", plan.GPULayers)
	}
	if plan.GPULayers > 33 {
		t.Fatalf("plan exceeds the model's layer count: %d", plan.GPULayers)
	}
	if plan.Estimate.GPUMB > 7200 {
		t.Fatalf("plan needs %d MB, above the 7200 MB budget", plan.Estimate.GPUMB)
	}
}

// Without VRAM the honest plan is CPU, and it must say it fits as a CPU plan —
// not pretend the GPU plan is safe.
func TestPlanChatLayersFallsToCPUWhenNoVRAM(t *testing.T) {
	plan := PlanChatLayers(qwen35ish(), 4096, 512, 2, 128, 30000)
	if plan.GPULayers != 0 {
		t.Fatalf("expected a CPU plan, got %d layers", plan.GPULayers)
	}
	if !plan.Fits {
		t.Fatalf("a CPU plan that fits RAM should report Fits: %+v", plan)
	}
	if plan.Estimate.GPUMB > 0 {
		t.Fatalf("CPU plan still charges %d MB of GPU", plan.Estimate.GPUMB)
	}
}

// A quantized KV cache must be cheaper than f16 for the same plan, and a tight
// budget must therefore allow at least as many layers.
func TestPlanChatLayersCostsQuantizedKV(t *testing.T) {
	f16 := PlanChatLayers(qwen35ish(), 4096, 512, 2, 3000, 30000)
	q8 := PlanChatLayers(qwen35ish(), 4096, 512, KVCacheBytesPerElement("q8_0"), 3000, 30000)

	if q8.GPULayers < f16.GPULayers {
		t.Fatalf("q8_0 placed fewer layers (%d) than f16 (%d) under the same budget", q8.GPULayers, f16.GPULayers)
	}
	if q8.GPULayers == f16.GPULayers && q8.Estimate.GPUMB >= f16.Estimate.GPUMB {
		t.Fatalf("q8_0 KV was not costed lower: q8 %d MB vs f16 %d MB", q8.Estimate.GPUMB, f16.Estimate.GPUMB)
	}
}

func TestKVCacheBytesPerElement(t *testing.T) {
	cases := map[string]int{"": 2, "f16": 2, "standard": 2, "q8_0": 1, "q4_0": 1, "nonsense": 2}
	for cacheType, want := range cases {
		if got := KVCacheBytesPerElement(cacheType); got != want {
			t.Fatalf("KVCacheBytesPerElement(%q) = %d, want %d", cacheType, got, want)
		}
	}
}

// A missing spec cannot be planned, and must not be presented as fitting.
func TestPlanChatLayersWithoutSpec(t *testing.T) {
	plan := PlanChatLayers(nil, 4096, 512, 2, 7200, 30000)
	if plan.Fits || plan.GPULayers != 0 {
		t.Fatalf("nil spec produced %+v", plan)
	}
}
