package hw

import (
	"testing"

	"howl-chat/internal/audio/types"
)

func TestDetermineMemoryTier(t *testing.T) {
	tests := []struct {
		name     string
		ramBytes int64
		expected types.MemoryTier
	}{
		{"Tiny RAM", 1 * 1024 * 1024 * 1024, types.TierTiny},
		{"Small RAM", 4 * 1024 * 1024 * 1024, types.TierSmall},
		{"Small RAM upper boundary", (8 * 1024 * 1024 * 1024) - 1, types.TierSmall},
		{"Medium RAM", 12 * 1024 * 1024 * 1024, types.TierMedium},
		{"Medium RAM upper boundary", (16 * 1024 * 1024 * 1024) - 1, types.TierMedium},
		{"Large RAM", 32 * 1024 * 1024 * 1024, types.TierLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetermineMemoryTier(tt.ramBytes)
			if result != tt.expected {
				t.Errorf("DetermineMemoryTier(%d) = %v, want %v", tt.ramBytes, result, tt.expected)
			}
		})
	}
}

func TestSelectModelsForHardware(t *testing.T) {
	tests := []struct {
		name     string
		hardware *HardwareInfo
		expected []types.ASRModelType
	}{
		{
			name: "Tiny RAM",
			hardware: &HardwareInfo{
				MemoryTier: types.TierTiny,
				HasGPU:     false,
			},
			expected: []types.ASRModelType{types.ASRWhisperTiny},
		},
		{
			name: "Small RAM",
			hardware: &HardwareInfo{
				MemoryTier: types.TierSmall,
				HasGPU:     false,
			},
			expected: []types.ASRModelType{types.ASRWhisperBase, types.ASRWhisperTiny},
		},
		{
			name: "Medium RAM with GPU",
			hardware: &HardwareInfo{
				MemoryTier: types.TierMedium,
				HasGPU:     true,
			},
			expected: []types.ASRModelType{types.ASRQwen3ASR06B, types.ASRWhisperSmall, types.ASRWhisperBase},
		},
		{
			name: "Large RAM with GPU",
			hardware: &HardwareInfo{
				MemoryTier: types.TierLarge,
				HasGPU:     true,
			},
			expected: []types.ASRModelType{types.ASRQwen3ASR17B, types.ASRQwen3ASR06B, types.ASRWhisperSmall, types.ASRWhisperBase},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SelectModelsForHardware(tt.hardware)
			if len(result) != len(tt.expected) {
				t.Errorf("SelectModelsForHardware() returned %d models, want %d", len(result), len(tt.expected))
			}
			for i, model := range result {
				if i >= len(tt.expected) || model != tt.expected[i] {
					t.Errorf("SelectModelsForHardware()[%d] = %v, want %v", i, model, tt.expected[i])
				}
			}
		})
	}
}

func TestFilterModelsByTier(t *testing.T) {
	allModels := []types.ASRModelType{
		types.ASRWhisperTiny,
		types.ASRWhisperBase,
		types.ASRWhisperSmall,
		types.ASRQwen3ASR06B,
		types.ASRQwen3ASR17B,
	}

	tests := []struct {
		name     string
		models   []types.ASRModelType
		minTier  types.MemoryTier
		expected []types.ASRModelType
	}{
		{
			name:     "Filter by Tiny tier",
			models:   allModels,
			minTier:  types.TierTiny,
			expected: []types.ASRModelType{types.ASRWhisperTiny},
		},
		{
			name:     "Filter by Small tier",
			models:   allModels,
			minTier:  types.TierSmall,
			expected: []types.ASRModelType{types.ASRWhisperTiny, types.ASRQwen3ASR06B, types.ASRWhisperBase, types.ASRWhisperSmall},
		},
		{
			name:     "Filter by Medium tier",
			models:   allModels,
			minTier:  types.TierMedium,
			expected: []types.ASRModelType{types.ASRWhisperTiny, types.ASRQwen3ASR06B, types.ASRWhisperBase, types.ASRWhisperSmall, types.ASRQwen3ASR17B},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterModelsByTier(tt.models, tt.minTier)
			if len(result) != len(tt.expected) {
				t.Errorf("FilterModelsByTier() returned %d models, want %d", len(result), len(tt.expected))
				return
			}
			for i, model := range result {
				if i >= len(tt.expected) || model != tt.expected[i] {
					t.Errorf("FilterModelsByTier()[%d] = %v, want %v", i, model, tt.expected[i])
				}
			}
		})
	}
}

func TestDetectHardware(t *testing.T) {
	t.Run("Hardware detection", func(t *testing.T) {
		info := DetectHardware()
		if info.TotalRAMBytes <= 0 {
			t.Error("TotalRAMBytes should be positive")
		}
		if info.OS == "" {
			t.Error("OS should not be empty")
		}
		if info.Arch == "" {
			t.Error("Arch should not be empty")
		}
		if info.MemoryTier == "" {
			t.Error("MemoryTier should be set")
		}
		t.Logf("Detected hardware: RAM=%d, GPU=%v, VRAM=%d, OS=%s, Arch=%s, Tier=%v",
			info.TotalRAMBytes, info.HasGPU, info.GPUMemoryBytes, info.OS, info.Arch, info.MemoryTier)
	})
}
