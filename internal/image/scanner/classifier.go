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
	AssetClip       AssetType = "clip"
	AssetT5         AssetType = "t5"
	AssetUnknown    AssetType = "unknown"
)

// ModelFamily identifies the base architecture.
type ModelFamily string

const (
	FamilyFlux     ModelFamily = "flux"
	FamilySDXL     ModelFamily = "sdxl"
	FamilySD3      ModelFamily = "sd3"
	FamilySD15     ModelFamily = "sd1.5"
	FamilySD21     ModelFamily = "sd2.1"
	FamilyTurbo    ModelFamily = "sdxl-turbo"
	FamilyAuraFlow ModelFamily = "auraflow"
	FamilyKolors   ModelFamily = "kolors"
	FamilyCascade  ModelFamily = "cascade"
	FamilyChroma   ModelFamily = "chroma"
	FamilyUnknown  ModelFamily = ""
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
	case contains(lower, "lora", "loha", "locon"):
		return AssetLoRA
	case contains(lower, "embed", "textual_inversion", "ti_", "hypernetwork"):
		return AssetEmbedding
	case contains(lower, "clip_l", "clip.l", "clip-l"):
		return AssetClip
	case contains(lower, "clip_g", "clip.g", "clip-g"):
		return AssetClip
	case contains(lower, "t5xxl", "t5-xxl", "t5_xxl"):
		return AssetT5
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

	switch {
	case contains(lower, "flux"):
		return FamilyFlux
	case contains(lower, "chroma"):
		// Chroma is a distinct SDXL-derivative architecture requiring specific flags
		return FamilyChroma
	case contains(lower, "sdxl-turbo", "sd_xl_turbo"):
		return FamilyTurbo
	case contains(lower, "auraflow", "aura-flow", "aura_flow"):
		return FamilyAuraFlow
	case contains(lower, "kolors"):
		return FamilyKolors
	case contains(lower, "cascade", "wurstchen", "würstchen"):
		return FamilyCascade
	case contains(lower, "sdxl", "sd_xl", "sd-xl", "_xl_", "-xl-",
		"juggernaut_xl", "juggernaut-xl", "realisticvision_xl", "dreamshaperxl"):
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
	case FamilyFlux, FamilySDXL, FamilySD3, FamilyTurbo:
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
	case FamilyFlux:
		return "euler"
	case FamilySDXL, FamilySD3, FamilyTurbo:
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
