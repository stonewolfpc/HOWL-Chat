package fallback

import (
	"runtime"
	"runtime/debug"

	"howl-chat/internal/audio/types"
)

// HardwareInfo contains detected hardware capabilities
type HardwareInfo struct {
	TotalRAMBytes     int64
	AvailableRAMBytes int64
	HasGPU            bool
	GPUMemoryBytes    int64
	MemoryTier        types.MemoryTier
	OS                string
	Arch              string
}

// DetectHardware queries system hardware capabilities
func DetectHardware() *HardwareInfo {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Get total system memory (this is an approximation)
	// On Windows, this may not be accurate, but we can use golang.org/x/sys/... for better detection
	totalRAM := getSystemRAM()

	info := &HardwareInfo{
		TotalRAMBytes:     totalRAM,
		AvailableRAMBytes: int64(memStats.Sys), // Available memory
		HasGPU:            detectGPU(),
		GPUMemoryBytes:    0, // TODO: Implement GPU memory detection
		OS:                runtime.GOOS,
		Arch:              runtime.GOARCH,
	}

	// Determine memory tier
	info.MemoryTier = determineMemoryTier(totalRAM)

	return info
}

// getSystemRAM attempts to detect total system RAM
func getSystemRAM() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// This is not accurate for total system RAM on most systems
	// For now, use a heuristic based on SysAlloc
	// In production, use golang.org/x/sys/... for accurate detection

	// Fallback: assume 8GB if we can't detect
	const fallbackRAM = 8 * 1024 * 1024 * 1024

	// Try to get from build info if available
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "GOOS" || setting.Key == "GOARCH" {
				// We could use this for platform-specific detection
			}
		}
	}

	// For now, return a reasonable default
	// TODO: Implement proper platform-specific detection
	return fallbackRAM
}

// detectGPU checks if a GPU is available
func detectGPU() bool {
	// Check for common GPU indicators
	// This is a simple heuristic - in production, use proper detection libraries

	// Check build tags or environment variables
	// For now, assume no GPU unless explicitly detected
	return false
}

// determineMemoryTier maps RAM to memory tier
func determineMemoryTier(ramBytes int64) types.MemoryTier {
	const (
		tinyRAM  = 2 * 1024 * 1024 * 1024  // 2GB
		smallRAM = 8 * 1024 * 1024 * 1024  // 8GB
		medRAM   = 16 * 1024 * 1024 * 1024 // 16GB
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

// SelectModelsForHardware returns recommended ASR models for detected hardware
func SelectModelsForHardware(info *HardwareInfo) []types.ASRModelType {
	switch info.MemoryTier {
	case types.TierTiny:
		// Very limited RAM - use only smallest models
		return []types.ASRModelType{
			types.ASRWhisperTiny,
		}
	case types.TierSmall:
		// Small RAM - can use Whisper Base
		return []types.ASRModelType{
			types.ASRWhisperBase,
			types.ASRWhisperTiny,
		}
	case types.TierMedium:
		// Medium RAM - can use Qwen-ASR or Whisper Small
		if info.HasGPU {
			return []types.ASRModelType{
				types.ASRQwen3ASR06B,
				types.ASRWhisperSmall,
				types.ASRWhisperBase,
			}
		}
		return []types.ASRModelType{
			types.ASRWhisperSmall,
			types.ASRWhisperBase,
			types.ASRWhisperTiny,
		}
	case types.TierLarge:
		// Large RAM - can use any model
		if info.HasGPU {
			return []types.ASRModelType{
				types.ASRQwen3ASR17B,
				types.ASRQwen3ASR06B,
				types.ASRWhisperSmall,
				types.ASRWhisperBase,
			}
		}
		return []types.ASRModelType{
			types.ASRWhisperSmall,
			types.ASRWhisperBase,
			types.ASRWhisperTiny,
		}
	default:
		// Auto - use safe defaults
		return []types.ASRModelType{
			types.ASRWhisperBase,
			types.ASRWhisperTiny,
		}
	}
}

// FilterModelsByTier filters models by minimum memory tier requirement
func FilterModelsByTier(models []types.ASRModelType, minTier types.MemoryTier) []types.ASRModelType {
	var filtered []types.ASRModelType

	for _, modelType := range models {
		if profile, ok := types.ASRModelRegistry[modelType]; ok {
			if profile.MinMemoryTier <= minTier {
				filtered = append(filtered, modelType)
			}
		}
	}

	return filtered
}
