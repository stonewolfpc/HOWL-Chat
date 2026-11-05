package gguf

import (
	"path/filepath"
	"strings"
)

// Kind describes what a GGUF file is actually for.
//
// HOWL stores chat models, diffusion models and image text encoders in the same
// container format, so filename heuristics are not enough to keep the menus
// apart: a 12B chat model named "Captain-Eris" looks like a large checkpoint to
// a size-based classifier. general.architecture is authoritative.
type Kind string

const (
	KindLLM       Kind = "llm"
	KindDiffusion Kind = "diffusion"
	KindUnknown   Kind = "unknown"
)

// diffusionArchs are general.architecture values emitted by image models.
// qwen_image is deliberately listed: it starts with "qwen" but is a diffusion
// model and must never be offered as a chat model.
var diffusionArchs = map[string]bool{
	"flux": true, "flux1": true,
	"sd1": true, "sd1.5": true, "sd15": true,
	"sdxl": true, "sd3": true, "sd3.5": true, "sd35": true,
	"stable-diffusion": true, "stable_diffusion": true,
	"chroma": true, "auraflow": true, "hidream": true, "kolors": true,
	"cascade": true, "pixart": true, "lumina": true, "cosmos": true,
	"wan": true, "wan2": true, "wan21": true, "wan22": true,
	"ltx": true, "ltxv": true, "hunyuan": true, "hunyuanvideo": true,
	"animatediff": true, "qwen_image": true, "qwen-image": true,
}

// llmArchs are general.architecture values emitted by chat/completion models.
var llmArchs = map[string]bool{
	"llama": true, "llama4": true,
	"qwen": true, "qwen2": true, "qwen2vl": true, "qwen2.5": true,
	"qwen3": true, "qwen3moe": true, "qwen3vl": true,
	"mistral": true, "mixtral": true,
	"phi2": true, "phi3": true, "phi4": true,
	"gemma": true, "gemma2": true, "gemma3": true,
	"deepseek": true, "deepseek2": true,
	"falcon": true, "mpt": true, "gpt2": true, "gptneox": true,
	"baichuan": true, "internlm2": true, "minicpm": true, "olmo": true,
	"exaone": true, "granite": true, "glm4": true, "chatglm": true,
	"starcoder": true, "starcoder2": true, "mllama": true, "llava": true,
	"command-r": true, "dbrx": true, "jamba": true, "nemotron": true,
	"cohere2": true, "seed_oss": true,
	// Families registered by the shipped model inventory that were missing
	// here, which made them look unknown to every architecture consumer.
	"gemma4": true, "qwen35": true, "gpt-oss": true,
	"lfm2": true, "lfm2moe": true,
}

// KindOfArch maps a general.architecture value to a Kind.
//
// Anything not recognised — whisper, t5, clip, bert, brand new architectures —
// returns KindUnknown so callers keep their existing behaviour rather than
// guessing and mis-filing a model.
func KindOfArch(arch string) Kind {
	a := strings.ToLower(strings.TrimSpace(arch))
	if a == "" {
		return KindUnknown
	}
	if diffusionArchs[a] {
		return KindDiffusion
	}
	if llmArchs[a] {
		return KindLLM
	}
	return KindUnknown
}

// KindOf reads a GGUF file's architecture and reports what it is for.
func KindOf(path string) Kind {
	if !strings.EqualFold(filepath.Ext(path), ".gguf") {
		return KindUnknown
	}
	profile, err := ExtractModelProfile(path)
	if err != nil || profile == nil {
		return KindUnknown
	}
	return KindOfArch(profile.Family)
}
