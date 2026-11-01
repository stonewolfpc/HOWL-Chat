// Package scanner provides model discovery, classification, manifest I/O,
// and index management for the HOWL image generation pipeline.
package scanner

import (
	"strings"
)

// AssetType classifies what a file actually is.
type AssetType string

const (
	AssetDiffusion  AssetType = "diffusion"
	AssetVAE        AssetType = "vae"
	AssetLoRA       AssetType = "lora"
	AssetControlNet AssetType = "controlnet"
	AssetEmbedding  AssetType = "embedding"
	AssetUpscaler   AssetType = "upscaler"
	AssetClipL      AssetType = "clip_l"
	AssetClipG      AssetType = "clip_g"
	AssetT5XXL      AssetType = "t5xxl"
	AssetLLM        AssetType = "llm" // For Anima's Qwen3 text encoder
	AssetUnknown    AssetType = "unknown"
)

// ModelFamily identifies the base architecture.
type ModelFamily string

const (
	FamilyFlux     ModelFamily = "flux"
	FamilySDXL     ModelFamily = "sdxl"
	FamilySD35     ModelFamily = "sd3.5"
	FamilySD3      ModelFamily = "sd3"
	FamilySD15     ModelFamily = "sd1.5"
	FamilySD21     ModelFamily = "sd2.1"
	FamilyTurbo    ModelFamily = "sdxl-turbo"
	FamilyAuraFlow ModelFamily = "auraflow"
	FamilyKolors   ModelFamily = "kolors"
	FamilyCascade  ModelFamily = "cascade"
	FamilyChroma   ModelFamily = "chroma"

	FamilyHiDream ModelFamily = "hidream"
	FamilyWan     ModelFamily = "wan"
	FamilyLTX     ModelFamily = "ltx2"
	FamilyHunyuan ModelFamily = "hunyuanvideo"
	FamilyAnimate ModelFamily = "animatediff"
	FamilyUnknown ModelFamily = ""
)

// modelExts is the set of file extensions that contain model weights.
var modelExts = map[string]bool{
	".gguf": true, ".safetensors": true, ".ckpt": true,
	".sft": true, ".pth": true, ".pt": true, ".bin": true,
}

// IsModelFile returns true if the extension is a recognised weight format.
func IsModelFile(name string) bool {
	return modelExts[strings.ToLower(fileExt(name))]
}

// ClassifyName uses filename heuristics (before size is known) to determine
// the likely asset type. Returns AssetUnknown when the name gives no signal.
func ClassifyName(name string) AssetType {
	lower := strings.ToLower(name)
	switch {
	case contains(lower, "controlnet", "control_net", "control-net"):
		return AssetControlNet
	case contains(lower, "ip-adapter", "ipadapter"):
		return AssetControlNet
	case contains(lower, "upscale", "esrgan", "swinir", "4x-", "_4x"):
		return AssetUpscaler
	case contains(lower, "vae"):
		return AssetVAE
	case lower == "ae.safetensors", strings.HasPrefix(lower, "ae."):
		return AssetVAE
	case contains(lower, "first_stage_model", "first-stage-model"):
		return AssetVAE
	case contains(lower, "lora", "loha", "locon"):
		return AssetLoRA
	case contains(lower, "embed", "textual_inversion", "ti_", "hypernetwork"):
		return AssetEmbedding
	case contains(lower, "clip_l", "clip.l", "clip-l"):
		return AssetClipL
	case contains(lower, "clip_g", "clip.g", "clip-g"):
		return AssetClipG
	case contains(lower, "t5xxl", "t5-xxl", "t5_xxl"):
		return AssetT5XXL

	case contains(lower, "qwen"):
		// Qwen text encoders (e.g., Qwen3-0.6B-Base for Anima)
		return AssetLLM
	case contains(lower, "wan2", "wan_2", "ltx", "hunyuan", "animatediff", "minimax"):
		return AssetDiffusion
	default:
		return AssetUnknown
	}
}

// ClassifySize refines an AssetUnknown based on file size.
// Order: known name signal wins, then size thresholds.
func ClassifySize(nameType AssetType, sizeMB float64) AssetType {
	if nameType != AssetUnknown {
		return nameType
	}
	switch {
	case sizeMB < 20:
		return AssetEmbedding
	case sizeMB < 350:
		return AssetLoRA
	case sizeMB < 900:
		return AssetVAE // most standalone VAEs are 300-900 MB
	default:
		return AssetDiffusion
	}
}

// GuessFamily infers the model architecture from the filename.
// Returns FamilyUnknown when the name is clearly not an SD model.
func GuessFamily(name string) ModelFamily {
	lower := strings.ToLower(name)

	// The Flux VAE ships as ae.safetensors (stem "ae"). Classify it as Flux so it
	// is never attached to an SD1.5/SDXL checkpoint that carries its own VAE.
	if lower == "ae" {
		return FamilyFlux
	}

	// Explicit non-SD model names — skip them entirely
	nonSD := []string{
		"whisper", "ggml-base", "ggml-small", "ggml-medium", "ggml-large",
		"espeak", "bert", "embed", "vocos", "piper", "bark", "kokoro",
	}
	for _, s := range nonSD {
		if strings.Contains(lower, s) {
			return FamilyUnknown
		}
	}

	// VAE-specific family detection
	if strings.Contains(lower, "vae") || strings.HasPrefix(lower, "ae.") {
		switch {
		case contains(lower, "flux", "sdxl", "xl", "juggernaut_xl", "juggernaut-xl",
			"realisticvision_xl", "realisticvision-xl", "dreamshaperxl", "dreamshaper-xl",
			"flatpiecexl", "piecexl"): // Explicitly handle flatpiecexl VAEs
			return FamilySDXL // Flux uses SDXL VAE
		case contains(lower, "sd3", "sd_3", "stable-diffusion-3", "sd3.5", "sd_3.5", "sd35"):
			return FamilySD3
		case contains(lower, "sd2.1", "sd_2.1", "sd21"):
			return FamilySD21

		case contains(lower, "wan"):
			return FamilyWan
		case contains(lower, "sd1.5", "sd_1.5", "sd15", "v1-5"):
			return FamilySD15
		default:
			return FamilySD15 // conservative default for VAEs
		}
	}

	// CLIP-specific family detection
	if strings.Contains(lower, "clip") {
		switch {
		case contains(lower, "flux", "sdxl", "xl"):
			return FamilySDXL
		case contains(lower, "sd3", "sd_3", "stable-diffusion-3", "sd3.5", "sd_3.5", "sd35"):
			return FamilySD3
		default:
			return FamilySD15
		}
	}

	// T5XXL-specific family detection
	if contains(lower, "t5xxl", "t5-xxl", "t5_xxl") {
		switch {
		case contains(lower, "flux", "chroma"):
			return FamilyFlux
		case contains(lower, "sd3", "sd_3", "stable-diffusion-3", "sd3.5", "sd_3.5", "sd35"):
			return FamilySD3
		default:
			return FamilyFlux // T5XXL is primarily used by Flux/Chroma/SD3
		}
	}

	switch {
	case contains(lower, "wan2", "wan_2", "wan-", "-wan"):
		return FamilyWan
	case contains(lower, "ltx"):
		return FamilyLTX
	case contains(lower, "hunyuan"):
		return FamilyHunyuan
	case contains(lower, "animatediff", "motion-module", "motion_module"):
		return FamilyAnimate
	case contains(lower, "flux"):
		return FamilyFlux
	case contains(lower, "chroma"):
		// Chroma is a distinct SDXL-derivative architecture requiring specific flags
		return FamilyChroma
	case contains(lower, "sdxl-turbo", "sd_xl_turbo"):
		return FamilyTurbo
	case contains(lower, "auraflow", "aura-flow", "aura_flow"):
		return FamilyAuraFlow
	case contains(lower, "hidream", "hi-dream", "hi_dream"):
		return FamilyHiDream
	case contains(lower, "kolors"):
		return FamilyKolors
	case contains(lower, "cascade", "wurstchen", "würstchen"):
		return FamilyCascade
	case contains(lower, "sd3.5", "sd_3.5", "sd35", "stable-diffusion-3.5"):
		return FamilySD35
	case contains(lower, "sdxl", "sd_xl", "sd-xl", "_xl_", "-xl-",
		"juggernaut_xl", "juggernaut-xl",
		"realisticvision_xl", "realisticvision-xl",
		"dreamshaperxl", "dreamshaper-xl",
		"epicrealism"):
		return FamilySDXL
	case contains(lower, "sd3", "sd_3", "stable-diffusion-3"):
		return FamilySD3
	case contains(lower, "sd2.1", "sd_2.1", "sd21", "2-1-"):
		return FamilySD21
	case contains(lower, "sd1.5", "sd_1.5", "sd15", "1-5-", "v1-5",
		"realisticvision", "dreamshaper", "deliberate", "absolutereality"):
		return FamilySD15
	default:
		return FamilySD15 // conservative default for unrecognised diffusion models
	}
}

// DefaultResolution returns the canonical resolution for a family.
func DefaultResolution(family ModelFamily) int {
	switch family {
	case FamilyFlux, FamilySDXL, FamilySD35, FamilySD3, FamilyTurbo:
		return 1024
	case FamilySD21:
		return 768
	default:
		return 512
	}
}

// DefaultSampler returns a reasonable sampler for the family.
func DefaultSampler(family ModelFamily) string {
	switch family {
	case FamilyFlux, FamilyWan, FamilyLTX, FamilyHunyuan, FamilyAnimate:
		return "euler"
	case FamilySDXL, FamilySD35, FamilySD3, FamilyTurbo:
		return "dpmpp_2m"
	default:
		return "euler_a"
	}
}

// CompatibleVAE returns true if the vaeFamily can be used with modelFamily.
func CompatibleVAE(modelFamily, vaeFamily ModelFamily) bool {
	if vaeFamily == FamilyUnknown {
		return true // unclassified VAE — allow it
	}
	// SD1.5 and SD2.1 can share VAEs, but NOT with SDXL VAEs
	if (modelFamily == FamilySD15 || modelFamily == FamilySD21) &&
		(vaeFamily == FamilySD15 || vaeFamily == FamilySD21) {
		return true
	}
	// SDXL models can use SDXL VAEs
	if modelFamily == FamilySDXL && vaeFamily == FamilySDXL {
		return true
	}
	// Flux and Chroma can share VAEs (both use SDXL VAE)
	if (modelFamily == FamilyFlux || modelFamily == FamilyChroma) &&
		(vaeFamily == FamilyFlux || vaeFamily == FamilyChroma || vaeFamily == FamilySDXL) {
		return true
	}
	// SD3 and SD3.5 can share VAEs
	if (modelFamily == FamilySD3 || modelFamily == FamilySD35) &&
		(vaeFamily == FamilySD3 || vaeFamily == FamilySD35) {
		return true
	}
	// Otherwise require exact match
	return modelFamily == vaeFamily
}

// CompatibleLoRA returns true if the loRA family works with the model family.
func CompatibleLoRA(modelFamily, loraFamily ModelFamily) bool {
	if loraFamily == FamilyUnknown {
		return true
	}
	return modelFamily == loraFamily
}

// CompatibleControlNet returns true if the ControlNet works with the model family.
func CompatibleControlNet(modelFamily, cnFamily ModelFamily) bool {
	if cnFamily == FamilyUnknown {
		return true
	}
	return modelFamily == cnFamily
}

// ---- helpers ---------------------------------------------------------------

func contains(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func fileExt(name string) string {
	idx := strings.LastIndex(name, ".")
	if idx < 0 {
		return ""
	}
	return name[idx:]
}
