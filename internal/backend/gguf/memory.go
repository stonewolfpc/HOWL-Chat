package gguf

import (
	"fmt"
	"math"

	"howl-chat/internal/backend/types"
)

// ---------------------------------------------------------------------------
// MemoryModel — unified, architecture-agnostic memory estimator
// ---------------------------------------------------------------------------

// MemoryModel holds the parameters and results of a memory estimation.
type MemoryModel struct {
	Spec        *types.ArchitectureSpec
	Settings    types.RuntimeSettings
	Result      types.MemoryEstimate
	Adjustments []string
}

// NewMemoryModel creates a memory estimator for the given architecture.
func NewMemoryModel(spec *types.ArchitectureSpec) *MemoryModel {
	return &MemoryModel{
		Spec:        spec,
		Adjustments: make([]string, 0, 8),
	}
}

// Estimate computes GPU and RAM requirements using the current settings.
// It does NOT clamp — it reports what the settings WILL use.
func (m *MemoryModel) Estimate(ctxSize, gpuLayers, batchSize int) types.MemoryEstimate {
	if m.Spec == nil {
		return types.MemoryEstimate{Notes: "no architecture spec"}
	}

	// Ensure sane defaults
	if ctxSize <= 0 {
		ctxSize = 4096
	}
	if gpuLayers < 0 {
		gpuLayers = m.Spec.Layers
	}
	if batchSize <= 0 {
		batchSize = 512
	}
	if gpuLayers > m.Spec.Layers {
		gpuLayers = m.Spec.Layers
	}

	weightMB := WeightMemoryMB(m.Spec)
	kvMB := KVCacheMemoryMB(m.Spec, ctxSize)
	actMB := ActivationMemoryMB(m.Spec, batchSize)

	// GPU memory breakdown
	gpuWeightMB := weightMB * float64(gpuLayers) / float64(m.Spec.Layers+1)
	gpuKVMB := kvMB * float64(gpuLayers) / float64(m.Spec.Layers)
	if m.Spec.Layers <= 1 {
		gpuKVMB = kvMB
	}
	gpuActivationMB := actMB * 0.5 // roughly half of activations on GPU
	gpuOverheadMB := 400.0         // CUDA context + compute buffers

	gpuTotalMB := gpuWeightMB + gpuKVMB + gpuActivationMB + gpuOverheadMB

	// RAM memory breakdown
	ramWeightMB := weightMB - gpuWeightMB
	ramKVMB := kvMB - gpuKVMB
	ramActivationMB := actMB * 0.5
	ramOverheadMB := 500.0 // llama.cpp workspace + OS overhead

	ramTotalMB := ramWeightMB + ramKVMB + ramActivationMB + ramOverheadMB

	// MoE adjustment: if not all experts are active, reduce weight estimate
	if m.Spec.IsMoE && m.Spec.ExpertCount > 0 && m.Spec.ExpertUsedCount > 0 {
		moeFactor := float64(m.Spec.ExpertUsedCount) / float64(m.Spec.ExpertCount)
		// Only the expert-router and active experts contribute to GPU
		moeReduction := weightMB * (1.0 - moeFactor) * 0.3 // ~30% of weights are experts
		gpuTotalMB -= moeReduction
		m.Adjustments = append(m.Adjustments,
			fmt.Sprintf("MoE: %d/%d active experts, reduced weight estimate by %.0f MB",
				m.Spec.ExpertUsedCount, m.Spec.ExpertCount, moeReduction))
	}

	// Sliding window reduces KV cache
	if m.Spec.UseSlidingWindow && m.Spec.SlidingWindowSize > 0 && m.Spec.SlidingWindowSize < ctxSize {
		swReduction := kvMB * (1.0 - float64(m.Spec.SlidingWindowSize)/float64(ctxSize))
		gpuTotalMB -= swReduction * (float64(gpuLayers) / float64(m.Spec.Layers))
		ramTotalMB -= swReduction * (1.0 - float64(gpuLayers)/float64(m.Spec.Layers))
		m.Adjustments = append(m.Adjustments,
			fmt.Sprintf("sliding_window=%d < ctx=%d, reduced KV by %.0f MB",
				m.Spec.SlidingWindowSize, ctxSize, swReduction))
	}

	notes := fmt.Sprintf(
		"arch=%s layers=%d d_model=%d heads=%d/%d ctx=%d gpu_layers=%d "+
			"weights=%.0fMB kv=%.0fMB act=%.0fMB overhead=%.0fMB confidence=%.2f",
		m.Spec.Family, m.Spec.Layers, m.Spec.DModel, m.Spec.NumHeads, m.Spec.NumKVHeads,
		ctxSize, gpuLayers, weightMB, kvMB, actMB, gpuOverheadMB+ramOverheadMB,
		m.Spec.Confidence(),
	)

	if len(m.Adjustments) > 0 {
		notes += " adjustments=" + fmt.Sprintf("%v", m.Adjustments)
	}

	return types.MemoryEstimate{
		GPUMB:           int64(math.Ceil(gpuTotalMB)),
		RAMMB:           int64(math.Ceil(ramTotalMB)),
		TotalMB:         int64(math.Ceil(gpuTotalMB + ramTotalMB)),
		ContextSizeUsed: ctxSize,
		GPULayersUsed:   gpuLayers,
		Notes:           notes,
	}
}

// ClampToHardware reduces settings until they fit within the given hardware limits.
// Returns the clamped ctxSize, gpuLayers, batchSize, the adjusted estimate,
// and a slice of adjustment descriptions.
func (m *MemoryModel) ClampToHardware(ctxSize, gpuLayers, batchSize int, freeVRAMMB, freeRAMMB int64) (int, int, int, types.MemoryEstimate, []string) {
	adjustments := make([]string, 0, 8)

	if m.Spec == nil || m.Spec.Layers <= 0 {
		return ctxSize, gpuLayers, batchSize, types.MemoryEstimate{Notes: "no spec, cannot clamp"}, adjustments
	}

	// Step 0: Ensure gpuLayers is within bounds
	if gpuLayers > m.Spec.Layers {
		gpuLayers = m.Spec.Layers
	}
	if gpuLayers < 0 {
		gpuLayers = m.Spec.Layers
	}

	// Step 1: Reduce GPU layers until VRAM fits
	for gpuLayers > 0 {
		est := m.Estimate(ctxSize, gpuLayers, batchSize)
		if est.GPUMB <= freeVRAMMB {
			break
		}
		old := gpuLayers
		gpuLayers = ReduceGPULayers(gpuLayers)
		adjustments = append(adjustments,
			fmt.Sprintf("GPU VRAM exceeded (%d MB > %d MB), reduced gpu_layers %d -> %d",
				est.GPUMB, freeVRAMMB, old, gpuLayers))
	}

	if gpuLayers == 0 {
		adjustments = append(adjustments,
			fmt.Sprintf("GPU VRAM insufficient even at 0 layers (%d MB free), using CPU only",
				freeVRAMMB))
	}

	// Step 2: Reduce context until RAM fits
	est := m.Estimate(ctxSize, gpuLayers, batchSize)
	minCtx := 2048
	if m.Spec.ContextLength > 0 && m.Spec.ContextLength < minCtx {
		minCtx = m.Spec.ContextLength / 2
		if minCtx < 512 {
			minCtx = 512
		}
	}

	for est.RAMMB > freeRAMMB && ctxSize > minCtx {
		old := ctxSize
		ctxSize = ReduceContext(ctxSize)
		if ctxSize < minCtx {
			ctxSize = minCtx
		}
		est = m.Estimate(ctxSize, gpuLayers, batchSize)
		adjustments = append(adjustments,
			fmt.Sprintf("RAM exceeded (%d MB > %d MB), reduced context %d -> %d",
				est.RAMMB, freeRAMMB, old, ctxSize))
		if ctxSize == old {
			break
		}
	}

	// Step 3: If RAM still exceeded, reduce batch size
	for est.RAMMB > freeRAMMB && batchSize > 128 {
		old := batchSize
		batchSize = batchSize / 2
		est = m.Estimate(ctxSize, gpuLayers, batchSize)
		adjustments = append(adjustments,
			fmt.Sprintf("RAM still exceeded (%d MB > %d MB), reduced batch %d -> %d",
				est.RAMMB, freeRAMMB, old, batchSize))
	}

	// Step 4: If RAM still exceeded and we're on GPU, move more layers to CPU
	for est.RAMMB > freeRAMMB && gpuLayers > 0 {
		old := gpuLayers
		gpuLayers = ReduceGPULayers(gpuLayers)
		est = m.Estimate(ctxSize, gpuLayers, batchSize)
		adjustments = append(adjustments,
			fmt.Sprintf("RAM still exceeded, moved layers to CPU: %d -> %d",
				old, gpuLayers))
		if gpuLayers == 0 {
			break
		}
	}

	// Step 5: For models with Sliding Window Attention (e.g. Gemma-4),
	// partial GPU offloading on Vulkan triggers GGML_SCHED_MAX_SPLIT_INPUTS.
	// Force either all layers on GPU (if they fit) or CPU-only.
	if m.Spec.UseSlidingWindow && gpuLayers > 0 && gpuLayers < m.Spec.Layers {
		allLayersEst := m.Estimate(ctxSize, m.Spec.Layers, batchSize)
		if allLayersEst.GPUMB <= freeVRAMMB {
			adjustments = append(adjustments,
				fmt.Sprintf("SWA model: partial offload unsafe, offloading ALL %d layers (GPU %d MB <= %d MB free)",
					m.Spec.Layers, allLayersEst.GPUMB, freeVRAMMB))
			gpuLayers = m.Spec.Layers
		} else {
			adjustments = append(adjustments,
				fmt.Sprintf("SWA model: partial offload unsafe, falling back to CPU-only (all layers would need %d MB > %d MB free)",
					allLayersEst.GPUMB, freeVRAMMB))
			gpuLayers = 0
		}
	}

	// Final estimate with clamped settings
	est = m.Estimate(ctxSize, gpuLayers, batchSize)
	if len(adjustments) > 0 {
		est.Notes = "CLAMPED: " + fmt.Sprintf("%v", adjustments) + " | " + est.Notes
	} else {
		est.Notes = "SAFE: fits within hardware | " + est.Notes
	}

	return ctxSize, gpuLayers, batchSize, est, adjustments
}

// ReduceGPULayers steps down to a lower GPU layer count.
func ReduceGPULayers(current int) int {
	switch {
	case current > 64:
		return current - 16
	case current > 32:
		return current - 8
	case current > 16:
		return current - 4
	case current > 4:
		return current - 2
	default:
		return current - 1
	}
}

// ReduceContext halves context, with a floor.
func ReduceContext(current int) int {
	next := current / 2
	if next < 512 {
		return 512
	}
	return next
}

// ValidateSpawn checks whether the given settings are safe to start a server with.
// It returns an error if the settings are impossible (would OOM or use invalid values).
func ValidateSpawn(spec *types.ArchitectureSpec, ctxSize, gpuLayers, batchSize int, freeVRAMMB, freeRAMMB int64) error {
	if spec == nil {
		return fmt.Errorf("no architecture spec available — cannot validate spawn")
	}
	if spec.Layers <= 0 {
		// Metadata missing layers — estimate from file size so we don't reject valid models
		if spec.FileSizeBytes > 0 {
			estParams := estimateParamsFromFileSize(spec.FileSizeBytes, spec.QuantType)
			spec.Layers = EstimateLayers(float64(estParams))
			if spec.DModel <= 0 {
				spec.DModel = EstimateDModel(float64(estParams))
			}
			if spec.NumHeads <= 0 {
				spec.NumHeads = EstimateNumHeads(spec.DModel)
			}
			if spec.NumKVHeads <= 0 {
				spec.NumKVHeads = spec.NumHeads
			}
			if spec.FFNSize <= 0 {
				spec.FFNSize = EstimateFFNSize(spec.DModel)
			}
			if spec.HeadDim <= 0 && spec.DModel > 0 && spec.NumHeads > 0 {
				spec.HeadDim = spec.DModel / spec.NumHeads
			}
			if spec.ContextLength <= 0 {
				spec.ContextLength = defaultContextLength(spec.Family)
			}
			fmt.Printf("[ValidateSpawn] estimated from file size: layers=%d dmodel=%d heads=%d\n",
				spec.Layers, spec.DModel, spec.NumHeads)
		}
	}
	if spec.Layers <= 0 {
		return fmt.Errorf("architecture spec has zero layers — cannot estimate memory (file_size=%d)", spec.FileSizeBytes)
	}
	if ctxSize <= 0 {
		return fmt.Errorf("context size must be > 0, got %d", ctxSize)
	}
	if gpuLayers < 0 || gpuLayers > spec.Layers {
		return fmt.Errorf("gpu_layers %d out of range [0, %d]", gpuLayers, spec.Layers)
	}
	if batchSize <= 0 {
		return fmt.Errorf("batch_size must be > 0, got %d", batchSize)
	}

	model := NewMemoryModel(spec)
	est := model.Estimate(ctxSize, gpuLayers, batchSize)

	// Hard reject if any single component exceeds 90% of available memory
	// (leaving headroom for OS and other processes)
	if est.GPUMB > int64(float64(freeVRAMMB)*1.1) {
		return fmt.Errorf(
			"spawn would exceed VRAM: estimated %d MB GPU > %d MB free (%d layers, %d context). "+
				"Reduce gpu_layers or context_size.",
			est.GPUMB, freeVRAMMB, gpuLayers, ctxSize)
	}
	if est.RAMMB > int64(float64(freeRAMMB)*1.1) {
		return fmt.Errorf(
			"spawn would exceed RAM: estimated %d MB RAM > %d MB free (%d layers, %d context). "+
				"Reduce context_size or use GPU offloading.",
			est.RAMMB, freeRAMMB, gpuLayers, ctxSize)
	}

	return nil
}
