package main

// app-family-defaults.go — per-family sampling defaults.
//
// Defaults stay "suggested general" but vary by family, because the publisher
// recommendations genuinely differ. Only families with a *published* number are
// listed here; everything else keeps the app's generic defaults. Each row carries
// the source it came from so a future edit can be checked against the original
// statement instead of somebody's memory.
//
// The user's explicit choices always win: these are the base that
// buildInferenceOptions() starts from before sampler settings are applied.

import (
	"strings"

	"howl-chat/internal/backend/llama"
)

type familySampling struct {
	// family is the display label used in logs and in GetEffectiveSampler.
	family      string
	temperature float64
	topP        float64
	minP        float64
	hasMinP     bool // MinP=0 is a published value for some families, not "unset"
	topK        int
	// Zero means "not published for this family" — the generic default stays.
	presencePenalty float64
	repeatPenalty   float64
}

var familySamplingTable = []struct {
	match string // lower-case substring of the GGUF general.architecture value
	samp  familySampling
	note  string
}{
	{
		// "For thinking mode (enable_thinking=True), use Temperature=0.6,
		// TopP=0.95, TopK=20, and MinP=0. DO NOT use greedy decoding, as it can
		// lead to performance degradation and endless repetitions."
		// https://huggingface.co/Qwen/Qwen3-8B/raw/main/README.md
		match: "qwen3",
		samp:  familySampling{family: "qwen3", temperature: 0.6, topP: 0.95, topK: 20, minP: 0, hasMinP: true, presencePenalty: 0},
		note:  "Qwen3 model card (thinking mode)",
	},
	{
		// "Use the following standardized sampling configuration across all use
		// cases: temperature=1.0, top_p=0.95, top_k=64"
		// https://ai.google.dev/gemma/docs/core/model_card_4
		// Gemma 3 ships the same top_k 64 / top_p 0.95 in its generation_config.
		match: "gemma",
		samp:  familySampling{family: "gemma", temperature: 1.0, topP: 0.95, topK: 64, minP: 0},
		note:  "Gemma 4 model card (Gemma 3 ships same top_k/top_p)",
	},
	{
		// "Set the temperature within the range of 0.5-0.7 (0.6 is recommended)
		// to prevent endless repetitions or incoherent outputs."
		// https://huggingface.co/deepseek-ai/DeepSeek-R1-Distill-Llama-8B/raw/main/README.md
		match: "deepseek",
		samp:  familySampling{family: "deepseek-r1", temperature: 0.6, topP: 0.95, minP: 0},
		note:  "DeepSeek-R1 model card (temp 0.5-0.7)",
	},
}

// familySamplingDefaults applies the published defaults for a family to opts and
// returns a short description of what was applied ("" when the family has none).
func familySamplingDefaults(opts *llama.InferenceOptions, family string) string {
	if opts == nil {
		return ""
	}
	f := strings.ToLower(family)
	if f == "" {
		return ""
	}
	for _, row := range familySamplingTable {
		if !strings.Contains(f, row.match) {
			continue
		}
		s := row.samp
		if s.temperature > 0 {
			opts.Temperature = s.temperature
		}
		if s.topP > 0 {
			opts.TopP = s.topP
			opts.TopPEnabled = true
		}
		if s.topK > 0 {
			opts.TopK = s.topK
		}
		if s.hasMinP {
			opts.MinP = s.minP
		}
		if s.presencePenalty > 0 {
			opts.PresencePenalty = s.presencePenalty
			opts.PresencePenaltyEnabled = true
		}
		if s.repeatPenalty > 0 {
			opts.RepeatPenalty = s.repeatPenalty
		}
		return s.family + " defaults (" + row.note + ")"
	}
	// Llama 3.x, Mistral-Nemo and most others publish no numbers at all, so the
	// generic defaults are the honest choice rather than invented ones.
	return ""
}

// familyName reports the loaded model's family for logs and the sampler binding.
func (a *App) familyName() string {
	if a.currentModelProfile != nil && a.currentModelProfile.Family != "" {
		return a.currentModelProfile.Family
	}
	return ""
}

// GetEffectiveSampler reports the values actually used for the loaded model:
// the per-family base merged with the user's explicit sampler settings. The UI
// reads this so its sliders show what the engine is really doing instead of
// whatever the last saved value happened to be.
func (a *App) GetEffectiveSampler() map[string]interface{} {
	opts := a.buildInferenceOptions()
	return map[string]interface{}{
		"family":            a.familyName(),
		"temperature":       opts.Temperature,
		"top_p":             opts.TopP,
		"top_k":             float64(opts.TopK),
		"min_p":             opts.MinP,
		"repeat_penalty":    opts.RepeatPenalty,
		"repeat_last_n":     float64(opts.RepeatLastN),
		"presence_penalty":  opts.PresencePenalty,
		"frequency_penalty": opts.FrequencyPenalty,
		"n_predict":         float64(opts.NPredict),
	}
}
