package runtime

import (
	"sort"
	"strings"

	"howl-chat/internal/audio/types"
)

// Constraints define runtime requirements used for model planning.
type Constraints struct {
	Tier             types.MemoryTier
	RequireStreaming bool
	RequireGPU       bool
	PreferredBackend string
	Languages        []string
}

// SelectionPlan contains ranked candidates by subsystem.
type SelectionPlan struct {
	ASR        []types.ASRModelType
	TTS        []types.TTSModelType
	Multimodal []types.MultimodalModelType
}

// BuildSelectionPlan returns ranked model candidates that satisfy constraints.
func BuildSelectionPlan(c Constraints) SelectionPlan {
	return SelectionPlan{
		ASR:        planASR(c),
		TTS:        planTTS(c),
		Multimodal: planMultimodal(c),
	}
}

func planASR(c Constraints) []types.ASRModelType {
	ranked := make([]rankedASR, 0, len(types.ASRModelRegistry))
	for modelType, profile := range types.ASRModelRegistry {
		if !tierSupports(profile.MinMemoryTier, normalizeTier(c.Tier)) {
			continue
		}
		if c.RequireStreaming && !profile.SupportsStreaming {
			continue
		}
		if c.RequireGPU && !profile.SupportsGPU {
			continue
		}
		if !supportsAllLanguages(profile.Languages, c.Languages) {
			continue
		}
		ranked = append(ranked, rankedASR{
			Model: modelType,
			Score: scoreModel(profile.ModelProfile, c),
		})
	}
	sort.Slice(ranked, func(i, j int) bool { return ranked[i].Score > ranked[j].Score })
	out := make([]types.ASRModelType, 0, len(ranked))
	for _, entry := range ranked {
		out = append(out, entry.Model)
	}
	return out
}

func planTTS(c Constraints) []types.TTSModelType {
	ranked := make([]rankedTTS, 0, len(types.ModelRegistry.TTSModels))
	for modelType, profile := range types.ModelRegistry.TTSModels {
		if !tierSupports(profile.MinMemoryTier, normalizeTier(c.Tier)) {
			continue
		}
		if c.RequireStreaming && !profile.SupportsStreaming {
			continue
		}
		if c.RequireGPU && !profile.SupportsGPU {
			continue
		}
		if !supportsAnyLanguage(profile.Languages, c.Languages) {
			continue
		}
		ranked = append(ranked, rankedTTS{
			Model: modelType,
			Score: scoreModel(profile.ModelProfile, c),
		})
	}
	sort.Slice(ranked, func(i, j int) bool { return ranked[i].Score > ranked[j].Score })
	out := make([]types.TTSModelType, 0, len(ranked))
	for _, entry := range ranked {
		out = append(out, entry.Model)
	}
	return out
}

func planMultimodal(c Constraints) []types.MultimodalModelType {
	ranked := make([]rankedMultimodal, 0, len(types.ModelRegistry.MultimodalModels))
	for modelType, profile := range types.ModelRegistry.MultimodalModels {
		if !tierSupports(profile.MinMemoryTier, normalizeTier(c.Tier)) {
			continue
		}
		if c.RequireStreaming && !profile.SupportsStreaming {
			continue
		}
		if c.RequireGPU && !profile.SupportsGPU {
			continue
		}
		if !supportsAnyLanguage(profile.Languages, c.Languages) {
			continue
		}
		ranked = append(ranked, rankedMultimodal{
			Model: modelType,
			Score: scoreModel(profile.ModelProfile, c),
		})
	}
	sort.Slice(ranked, func(i, j int) bool { return ranked[i].Score > ranked[j].Score })
	out := make([]types.MultimodalModelType, 0, len(ranked))
	for _, entry := range ranked {
		out = append(out, entry.Model)
	}
	return out
}

func scoreModel(profile types.ModelProfile, c Constraints) float64 {
	score := profile.QualityScore*100 - profile.RealTimeFactor*10
	if strings.EqualFold(c.PreferredBackend, profile.Backend) && c.PreferredBackend != "" {
		score += 20
	}
	if profile.SupportsQuantized {
		score += 5
	}
	return score
}

func supportsAllLanguages(modelLangs []string, required []string) bool {
	for _, lang := range required {
		found := false
		for _, candidate := range modelLangs {
			if strings.EqualFold(lang, candidate) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func supportsAnyLanguage(modelLangs []string, preferred []string) bool {
	if len(preferred) == 0 {
		return true
	}
	for _, lang := range preferred {
		for _, candidate := range modelLangs {
			if strings.EqualFold(lang, candidate) {
				return true
			}
		}
	}
	return false
}

func normalizeTier(tier types.MemoryTier) types.MemoryTier {
	if tier == "" || tier == types.TierAuto {
		return types.TierLarge
	}
	return tier
}

func tierSupports(required types.MemoryTier, available types.MemoryTier) bool {
	return tierRank(required) <= tierRank(available)
}

func tierRank(tier types.MemoryTier) int {
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

type rankedASR struct {
	Model types.ASRModelType
	Score float64
}

type rankedTTS struct {
	Model types.TTSModelType
	Score float64
}

type rankedMultimodal struct {
	Model types.MultimodalModelType
	Score float64
}
