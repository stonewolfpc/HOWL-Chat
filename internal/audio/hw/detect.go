package hw

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/mem"

	"howl-chat/internal/audio/types"
)

// DetectHardware queries system memory, GPU presence, and estimated VRAM.
func DetectHardware() *HardwareInfo {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		const fallbackRAM = 8 * 1024 * 1024 * 1024
		return &HardwareInfo{
			TotalRAMBytes:     fallbackRAM,
			AvailableRAMBytes: fallbackRAM,
			HasGPU:            false,
			GPUMemoryBytes:    0,
			OS:                runtime.GOOS,
			Arch:              runtime.GOARCH,
			MemoryTier:        types.TierSmall,
			Error:             fmt.Errorf("hardware detection failed: %w", err),
		}
	}

	hasGPU := detectGPUPresent()
	vram := int64(0)
	if hasGPU {
		vram = detectVRAMBytes()
	}

	return &HardwareInfo{
		TotalRAMBytes:     int64(vmStat.Total),
		AvailableRAMBytes: int64(vmStat.Available),
		HasGPU:            hasGPU,
		GPUMemoryBytes:    vram,
		OS:                runtime.GOOS,
		Arch:              runtime.GOARCH,
		MemoryTier:        DetermineMemoryTier(int64(vmStat.Total)),
	}
}

// DetermineMemoryTier maps installed RAM to a tier used for model policy.
func DetermineMemoryTier(ramBytes int64) types.MemoryTier {
	const (
		tinyRAM  = 2 * 1024 * 1024 * 1024
		smallRAM = 8 * 1024 * 1024 * 1024
		medRAM   = 16 * 1024 * 1024 * 1024
	)
	switch {
	case ramBytes < tinyRAM:
		return types.TierTiny
	case ramBytes < smallRAM:
		return types.TierSmall
	case ramBytes < medRAM:
		return types.TierMedium
	default:
		return types.TierLarge
	}
}

// SystemRAM returns total physical RAM or an error if unavailable.
func SystemRAM() (int64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("memory detection failed: %w", err)
	}
	return int64(vmStat.Total), nil
}
