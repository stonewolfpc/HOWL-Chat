package hw

import (
	"sort"

	"howl-chat/internal/audio/types"
)

// SelectModelsForHardware returns recommended ASR models for detected hardware.
func SelectModelsForHardware(info *HardwareInfo) []types.ASRModelType {
	switch info.MemoryTier {
	case types.TierTiny:
		return []types.ASRModelType{types.ASRWhisperTiny}
	case types.TierSmall:
		return []types.ASRModelType{types.ASRWhisperBase, types.ASRWhisperTiny}
	case types.TierMedium:
		if info.HasGPU {
			return []types.ASRModelType{types.ASRQwen3ASR06B, types.ASRWhisperSmall, types.ASRWhisperBase}
		}
		return []types.ASRModelType{types.ASRWhisperSmall, types.ASRWhisperBase, types.ASRWhisperTiny}
	case types.TierLarge:
		if info.HasGPU {
			return []types.ASRModelType{types.ASRQwen3ASR17B, types.ASRQwen3ASR06B, types.ASRWhisperSmall, types.ASRWhisperBase}
		}
		return []types.ASRModelType{types.ASRWhisperSmall, types.ASRWhisperBase, types.ASRWhisperTiny}
	default:
		return []types.ASRModelType{types.ASRWhisperBase, types.ASRWhisperTiny}
	}
}

// FilterModelsByTier filters models by minimum memory tier requirement.
func FilterModelsByTier(models []types.ASRModelType, minTier types.MemoryTier) []types.ASRModelType {
	var filtered []types.ASRModelType
	for _, modelType := range models {
		if profile, ok := types.ASRModelRegistry[modelType]; ok {
			if memoryTierRank(profile.MinMemoryTier) <= memoryTierRank(minTier) {
				filtered = append(filtered, modelType)
			}
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		pi, _ := types.ASRModelRegistry[filtered[i]]
		pj, _ := types.ASRModelRegistry[filtered[j]]
		ri := memoryTierRank(pi.MinMemoryTier)
		rj := memoryTierRank(pj.MinMemoryTier)
		if ri != rj {
			return ri < rj
		}
		return string(filtered[i]) < string(filtered[j])
	})
	return filtered
}

func memoryTierRank(tier types.MemoryTier) int {
	switch tier {
	case types.TierTiny:
		return 0
	case types.TierSmall:
		return 1
	case types.TierMedium:
		return 2
	case types.TierLarge, types.TierAuto:
		return 3
	default:
		return -1
	}
}
