package resource

// chat_placement.go — the upward half of chat placement.
//
// The estimator (internal/backend/gguf.MemoryModel) reports what a given
// context/batch/layer plan will use. Until now Howl only ever asked it to clamp
// a plan downward, so a request with gpu_layers 0 stayed on the CPU even on an
// idle card. This file is the one place that searches upward: given the
// post-reserve budget admission has already computed, it returns the largest
// layer count whose estimate fits. It does not admit, reserve or spawn
// anything; admission still owns that decision.

import (
	"fmt"
	"strings"

	"howl-chat/internal/backend/gguf"
	"howl-chat/internal/backend/types"
)

// KVCacheBytesPerElement maps the engine cache type Howl sends through
// --cache-type-k/v to the bytes one KV element occupies in the estimator.
//
// f16 is the engine default and the estimator's own default (2 bytes). q8_0 is
// 1 byte per element plus block scales; 1 is used because the estimator has no
// block-scale term and the difference (about 6%) is smaller than the 10%
// reserve. q4_0 is a 4-bit value plus scales; it is charged as 1 byte for the
// same reason, which over-reserves rather than under-reserves. An unknown or
// empty type keeps the 2-byte default, so the estimate can only be too high.
func KVCacheBytesPerElement(cacheType string) int {
	switch strings.ToLower(strings.TrimSpace(cacheType)) {
	case "q8_0":
		return 1
	case "q4_0":
		return 1
	default:
		return 2
	}
}

// ChatPlan is the result of an upward placement search.
type ChatPlan struct {
	// GPULayers is the chosen layer count. Zero means no layer count fits the
	// VRAM budget, which is a CPU plan and never selected while another model
	// is blocking the device (see Authority.Admit).
	GPULayers int
	// Fits reports that the returned plan fits both budgets. When false the
	// caller must refuse or fall back explicitly; it must not present the plan
	// as safe.
	Fits bool
	// Estimate is the estimator's own answer for the returned layer count.
	Estimate types.MemoryEstimate
	// Notes are the estimator's adjustments plus the search's own statement.
	Notes []string
}

// PlanChatLayers searches downward from every layer to zero and returns the
// first plan that fits the given post-reserve budgets. The estimator is the
// existing architecture-aware one; this function adds no new memory model.
//
// kvBytesPerElement is applied to a copy of the spec so the shared profile is
// never mutated and a quantized KV cache is costed for what the engine will
// actually build.
func PlanChatLayers(
	spec *types.ArchitectureSpec,
	ctxSize, batchSize, kvBytesPerElement int,
	availableVRAMMB, availableRAMMB int64,
) ChatPlan {
	if spec == nil || spec.Layers <= 0 {
		return ChatPlan{Notes: []string{"no architecture spec; no placement can be planned"}}
	}
	if ctxSize <= 0 {
		ctxSize = 4096
	}
	if batchSize <= 0 {
		batchSize = 512
	}

	costed := spec
	if kvBytesPerElement > 0 && spec.KVCacheBytesPerElement != kvBytesPerElement {
		clone := *spec
		clone.KVCacheBytesPerElement = kvBytesPerElement
		costed = &clone
	}

	model := gguf.NewMemoryModel(costed)
	var last types.MemoryEstimate
	for layers := spec.Layers; layers >= 0; layers-- {
		est := model.Estimate(ctxSize, layers, batchSize)
		// A zero-layer plan is CPU-only: there is no CUDA context or device
		// buffer, which is exactly how estimateChat charges it. The estimator's
		// fixed 400 MB GPU overhead does not apply to it.
		vramNeed := est.GPUMB
		if layers == 0 {
			vramNeed = 0
			est.GPUMB = 0
			est.TotalMB = est.RAMMB
		}
		last = est
		if vramNeed <= availableVRAMMB && est.RAMMB <= availableRAMMB {
			notes := append([]string{}, model.Adjustments...)
			notes = append(notes, fmt.Sprintf(
				"auto placement: %d/%d GPU layers fit %d MB VRAM and %d MB RAM after reserve (estimate %d MB VRAM / %d MB RAM, KV %d bytes/element)",
				layers, spec.Layers, availableVRAMMB, availableRAMMB,
				est.GPUMB, est.RAMMB, costed.KVCacheBytesPerElement))
			return ChatPlan{GPULayers: layers, Fits: true, Estimate: est, Notes: notes}
		}
	}

	_ = last
	notes := append([]string{}, model.Adjustments...)
	notes = append(notes, fmt.Sprintf(
		"auto placement: no layer count fits %d MB VRAM and %d MB RAM after reserve, and the model's weights do not fit RAM either",
		availableVRAMMB, availableRAMMB))
	return ChatPlan{GPULayers: 0, Fits: false, Estimate: last, Notes: notes}
}
