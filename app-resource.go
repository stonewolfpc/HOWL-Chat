package main

import (
	"fmt"

	"howl-chat/internal/errors"
	"howl-chat/internal/resource"
)

// ---------------------------------------------------------------------------
// Wails bindings: Resource Manager
// ---------------------------------------------------------------------------

// ResourceSnapshot exposes the current hardware snapshot to the frontend.
type ResourceSnapshot struct {
	GPUModel    string `json:"gpu_model"`
	GPUVendor   string `json:"gpu_vendor"`
	TotalVRAMMB int64  `json:"total_vram_mb"`
	UsedVRAMMB  int64  `json:"used_vram_mb"`
	FreeVRAMMB  int64  `json:"free_vram_mb"`
	TotalRAMMB  int64  `json:"total_ram_mb"`
	UsedRAMMB   int64  `json:"used_ram_mb"`
	FreeRAMMB   int64  `json:"free_ram_mb"`
	MemoryTier  string `json:"memory_tier"`
	OS          string `json:"os"`
}

// LoadedModelInfo exposes a loaded model's footprint to the frontend.
type LoadedModelInfo struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	GPUMemoryMB  int64  `json:"gpu_memory_mb"`
	SystemRAMMB  int64  `json:"system_ram_mb"`
	Priority     int    `json:"priority"`
	IsGPUOffload bool   `json:"is_gpu_offload"`
}

// AllocationAdvice exposes the result of an allocation request.
type AllocationAdvice struct {
	UseGPU        bool     `json:"use_gpu"`
	GPULayers     int      `json:"gpu_layers"`
	GPUMemoryMB   int64    `json:"gpu_memory_mb"`
	SystemRAMMB   int64    `json:"system_ram_mb"`
	ReduceBatch   bool     `json:"reduce_batch"`
	ReduceContext bool     `json:"reduce_context"`
	UnloadModels  []string `json:"unload_models"`
	Reason        string   `json:"reason"`
	Confidence    float64  `json:"confidence"`
}

// ---------------------------------------------------------------------------
// Hardware queries
// ---------------------------------------------------------------------------

// GetResourceSnapshot returns the current hardware state.
func (a *App) GetResourceSnapshot() ResourceSnapshot {
	s := resource.GetManager().HardwareSnapshot()
	return ResourceSnapshot{
		GPUModel:    s.GPUModel,
		GPUVendor:   s.GPUVendor,
		TotalVRAMMB: s.TotalVRAMBytes / 1024 / 1024,
		UsedVRAMMB:  s.UsedVRAMBytes / 1024 / 1024,
		FreeVRAMMB:  s.FreeVRAMBytes / 1024 / 1024,
		TotalRAMMB:  s.TotalRAMBytes / 1024 / 1024,
		UsedRAMMB:   s.UsedRAMBytes / 1024 / 1024,
		FreeRAMMB:   s.FreeRAMBytes / 1024 / 1024,
		MemoryTier:  s.MemoryTier,
		OS:          s.OS,
	}
}

// RefreshResourceSnapshot forces an immediate hardware re-scan.
func (a *App) RefreshResourceSnapshot() {
	resource.GetManager().RefreshHardware()
	s := resource.GetManager().HardwareSnapshot()
	if s.GPUModel == "" && s.TotalVRAMBytes == 0 {
		a.logErrorSilent(errors.RESHWDetectFail, errors.Medium, "resource.hwrefresh", "app-resource.go:72", fmt.Errorf("GPU detection returned empty"))
	}
}

// ---------------------------------------------------------------------------
// Loaded model queries
// ---------------------------------------------------------------------------

// GetLoadedModels returns all currently tracked loaded models.
func (a *App) GetLoadedModels() []LoadedModelInfo {
	models := resource.GetManager().LoadedModels()
	out := make([]LoadedModelInfo, len(models))
	for i, m := range models {
		out[i] = LoadedModelInfo{
			Type:         string(m.Type),
			Name:         m.Name,
			GPUMemoryMB:  m.GPUMemoryMB,
			SystemRAMMB:  m.SystemRAMMB,
			Priority:     m.Priority,
			IsGPUOffload: m.IsGPUOffload,
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Model registration (called by backend loaders)
// ---------------------------------------------------------------------------

// RegisterLoadedModel tells the resource manager a model is now loaded.
func (a *App) RegisterLoadedModel(modelType, name, path string, gpuMB, ramMB int64, priority int, isGPU bool) {
	rt := resource.ResourceType(modelType)
	resource.GetManager().RegisterModel(rt, name, path, gpuMB, ramMB, priority, isGPU)
}

// UnregisterLoadedModel tells the resource manager a model was unloaded.
func (a *App) UnregisterLoadedModel(modelType, name string) {
	rt := resource.ResourceType(modelType)
	resource.GetManager().UnregisterModel(rt, name)
}

// ---------------------------------------------------------------------------
// Allocation advice (frontend can call before loading)
// ---------------------------------------------------------------------------

// RequestAllocation asks the resource manager how to safely load a model.
func (a *App) RequestAllocation(
	modelType, modelName string,
	requiredVRAMMB, requiredRAMMB int64,
	priority int,
	preferGPU, allowCPU bool,
) AllocationAdvice {
	req := resource.AllocationRequest{
		ResourceType:   resource.ResourceType(modelType),
		ModelName:      modelName,
		RequiredVRAMMB: requiredVRAMMB,
		RequiredRAMMB:  requiredRAMMB,
		Priority:       priority,
		PreferGPU:      preferGPU,
		AllowCPU:       allowCPU,
	}
	res := resource.Allocate(resource.GetManager(), req)
	if res.Confidence < 0.3 {
		a.logError(errors.RESAllocFail, errors.High, "resource.allocate", "app-resource.go:136", fmt.Errorf("low confidence allocation: %s", res.Reason))
	}
	return AllocationAdvice{
		UseGPU:        res.UseGPU,
		GPULayers:     res.GPULayers,
		GPUMemoryMB:   res.GPUMemoryMB,
		SystemRAMMB:   res.SystemRAMMB,
		ReduceBatch:   res.ReduceBatch,
		ReduceContext: res.ReduceContext,
		UnloadModels:  res.UnloadModels,
		Reason:        res.Reason,
		Confidence:    res.Confidence,
	}
}

// GetResourceSummary returns a one-line human-readable resource status.
func (a *App) GetResourceSummary() string {
	return resource.GetManager().Summary()
}
