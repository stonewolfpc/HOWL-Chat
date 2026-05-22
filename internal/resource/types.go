// Package resource provides cross-component resource coordination for HOWL.
// It tracks hardware state, loaded model memory usage, and recommends
// allocation strategies to prevent OOM when multiple AI backends run concurrently.
package resource

import "time"

// ---------------------------------------------------------------------------
// Resource classification
// ---------------------------------------------------------------------------

// ResourceType identifies the category of an AI model asset.
type ResourceType string

const (
	ResourceLLM   ResourceType = "llm"
	ResourceTTS   ResourceType = "tts"
	ResourceImage ResourceType = "image"
)

// ---------------------------------------------------------------------------
// Model state
// ---------------------------------------------------------------------------

// ModelState represents a loaded model's runtime footprint.
type ModelState struct {
	Type         ResourceType `json:"type"`
	Name         string       `json:"name"`
	Path         string       `json:"path"`
	Loaded       bool         `json:"loaded"`
	GPUMemoryMB  int64        `json:"gpu_memory_mb"`
	SystemRAMMB  int64        `json:"system_ram_mb"`
	Priority     int          `json:"priority"`     // 0 = highest (chat), 10 = lowest (image)
	LastUsed     time.Time    `json:"last_used"`
	IsGPUOffload bool         `json:"is_gpu_offload"`
}

// ---------------------------------------------------------------------------
// Hardware snapshot
// ---------------------------------------------------------------------------

// HardwareSnapshot captures a point-in-time view of system resources.
type HardwareSnapshot struct {
	Timestamp      time.Time `json:"timestamp"`
	GPUModel       string    `json:"gpu_model"`
	GPUVendor      string    `json:"gpu_vendor"`
	TotalVRAMBytes int64     `json:"total_vram_bytes"`
	UsedVRAMBytes  int64     `json:"used_vram_bytes"`
	FreeVRAMBytes  int64     `json:"free_vram_bytes"`
	TotalRAMBytes  int64     `json:"total_ram_bytes"`
	UsedRAMBytes   int64     `json:"used_ram_bytes"`
	FreeRAMBytes   int64     `json:"free_ram_bytes"`
	AvailableRAM   int64     `json:"available_ram_bytes"`
	MemoryTier     string    `json:"memory_tier"`
	OS             string    `json:"os"`
	Arch           string    `json:"arch"`
}

// ---------------------------------------------------------------------------
// Allocation
// ---------------------------------------------------------------------------

// AllocationRequest describes what a backend wants to load.
type AllocationRequest struct {
	ResourceType   ResourceType `json:"resource_type"`
	ModelName      string       `json:"model_name"`
	RequiredVRAMMB int64        `json:"required_vram_mb"`
	RequiredRAMMB  int64        `json:"required_ram_mb"`
	Priority       int          `json:"priority"`
	PreferGPU      bool         `json:"prefer_gpu"`
	AllowCPU       bool         `json:"allow_cpu"`
}

// AllocationResult is the recommended strategy returned by the allocator.
type AllocationResult struct {
	UseGPU        bool     `json:"use_gpu"`
	GPULayers     int      `json:"gpu_layers"`     // -1 = auto/default
	GPUMemoryMB   int64    `json:"gpu_memory_mb"`
	SystemRAMMB   int64    `json:"system_ram_mb"`
	ReduceBatch   bool     `json:"reduce_batch"`
	ReduceContext bool     `json:"reduce_context"`
	UnloadModels  []string `json:"unload_models"` // Models to unload first
	Reason        string   `json:"reason"`
	Confidence    float64  `json:"confidence"`    // 0.0–1.0
}
