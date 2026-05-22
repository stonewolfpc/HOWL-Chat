package resource

import (
	"fmt"
	"sort"
)

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

// Priority values for different resource types.
// Lower = more important, less likely to be evicted.
const (
	PriorityLLM   = 0  // Chat model — never evict automatically
	PriorityTTS   = 5  // Voice — can evict if needed
	PriorityImage = 8  // Image gen — evict first
)

// Safety margins (percent of reported free space to avoid edge-of-cliff).
const vramSafetyPct = 0.90 // Use 90% of reported free VRAM
const ramSafetyPct  = 0.95 // Use 95% of reported free RAM

// ---------------------------------------------------------------------------
// Allocate
// ---------------------------------------------------------------------------

// Allocate determines the best resource strategy for a new model load.
// It reads the current hardware snapshot and loaded models from rm.
func Allocate(rm *ResourceManager, req AllocationRequest) AllocationResult {
	rm.mu.RLock()
	hw := rm.snapshot
	models := make(map[string]*ModelState, len(rm.models))
	for k, v := range rm.models {
		models[k] = v
	}
	rm.mu.RUnlock()

	res := AllocationResult{
		UseGPU:       req.PreferGPU,
		GPULayers:    -1,
		GPUMemoryMB:  req.RequiredVRAMMB,
		SystemRAMMB:  req.RequiredRAMMB,
		Confidence:   1.0,
	}

	freeVRAM := int64(float64(hw.FreeVRAMBytes) * vramSafetyPct / 1024 / 1024)
	freeRAM  := int64(float64(hw.FreeRAMBytes)  * ramSafetyPct  / 1024 / 1024)

	// ---- GPU path --------------------------------------------------------
	if req.PreferGPU && req.RequiredVRAMMB > 0 {
		if freeVRAM >= req.RequiredVRAMMB {
			res.Reason = fmt.Sprintf("VRAM OK (%d MB free ≥ %d MB needed)", freeVRAM, req.RequiredVRAMMB)
			return res
		}

		// Not enough VRAM — try eviction of lower-priority models
		candidates := evictionCandidates(models, req.ResourceType, req.Priority)
		freed := int64(0)
		for _, c := range candidates {
			if freed >= req.RequiredVRAMMB {
				break
			}
			res.UnloadModels = append(res.UnloadModels, c.Name)
			freed += c.GPUMemoryMB
		}

		if freeVRAM+freed >= req.RequiredVRAMMB {
			res.Reason = fmt.Sprintf("Unload %d model(s) to free %d MB VRAM", len(res.UnloadModels), freed)
			res.Confidence = 0.75
			return res
		}

		// Still not enough — fall back to CPU if allowed
		if req.AllowCPU {
			res.UseGPU = false
			res.GPULayers = 0
			res.GPUMemoryMB = 0
			res.SystemRAMMB = req.RequiredVRAMMB + req.RequiredRAMMB
			res.Reason = fmt.Sprintf(
				"Insufficient VRAM (%d MB free, %d MB needed); falling back to CPU",
				freeVRAM, req.RequiredVRAMMB,
			)
			res.Confidence = 0.70
		} else {
			res.Reason = fmt.Sprintf(
				"Insufficient VRAM (%d MB free, %d MB needed) and CPU fallback disabled",
				freeVRAM, req.RequiredVRAMMB,
			)
			res.Confidence = 0.20
		}
	}

	// ---- RAM check -------------------------------------------------------
	if res.SystemRAMMB > freeRAM {
		res.ReduceBatch = true
		res.ReduceContext = true
		res.SystemRAMMB = freeRAM / 2 // conservative
		res.Reason += fmt.Sprintf("; RAM constrained (%d MB free), reducing batch/context", freeRAM)
		res.Confidence *= 0.80
	}

	return res
}

// ---------------------------------------------------------------------------
// Eviction helpers
// ---------------------------------------------------------------------------

// evictionCandidates returns loaded GPU-offloaded models that are lower
// priority than the incoming request, sorted by priority (lowest first).
func evictionCandidates(
	models map[string]*ModelState,
	incomingType ResourceType,
	incomingPriority int,
) []*ModelState {
	var candidates []*ModelState
	for _, m := range models {
		if !m.Loaded || !m.IsGPUOffload {
			continue
		}
		// Never evict same resource type (e.g. don't unload LLM to load LLM)
		if m.Type == incomingType {
			continue
		}
		if m.Priority < incomingPriority {
			candidates = append(candidates, m)
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		// Lower priority first, then least recently used
		if candidates[i].Priority != candidates[j].Priority {
			return candidates[i].Priority > candidates[j].Priority
		}
		return candidates[i].LastUsed.Before(candidates[j].LastUsed)
	})

	return candidates
}
