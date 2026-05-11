package hw

import "howl-chat/internal/audio/types"

// HardwareInfo contains detected hardware capabilities.
type HardwareInfo struct {
	TotalRAMBytes     int64
	AvailableRAMBytes int64
	HasGPU            bool
	GPUMemoryBytes    int64
	MemoryTier        types.MemoryTier
	OS                string
	Arch              string
	Error             error
}
