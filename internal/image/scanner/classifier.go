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
	AssetIPAdapter  AssetType = "ip_adapter"
	AssetClipVision AssetType = "clip_vision"
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

	// Families the shipped stable-diffusion.cpp build implements but the app
	// previously had no name for. Without them these models were refused before
	// they could reach the engine.
	FamilyFlux2     ModelFamily = "flux2"
	FamilyZImage    ModelFamily = "z_image"
	FamilyQwenImage ModelFamily = "qwen_image"
	// FamilyKrea2 has a loader path of its own. Its GGUF header carries the
	// generic base tag qwen_image, but sd-server reports
	// "Detected model family: Krea2" and then refuses the Qwen2.5-VL encoder, so
	// classifying it as qwen_image hands it the wrong text component.
	FamilyKrea2 ModelFamily = "krea2"
	// Animagine-style Anima DiT: net.blocks with adaLN modulation.
	FamilyAnima ModelFamily = "anima"

	FamilyHiDream ModelFamily = "hidream"
	FamilyWan     ModelFamily = "wan"
	FamilyLTX     ModelFamily = "ltx2"
	FamilyHunyuan ModelFamily = "hunyuanvideo"
	FamilyAnimate ModelFamily = "animatediff"
	FamilyUnknown ModelFamily = ""
)

// refineFamilyFromName lets a filename pick between two families that share one
// tensor layout when the weights themselves only identify the base one.
//
// Converted GGUFs are the case that needs it: a Chroma, SD3.5 or SD2.1 checkpoint
// keeps its base model's metadata (flux / sd3 / sd1) and, for SD2.1 and SD3.5,
// loses the projection width that told them apart, so the layout can name only
// the base. Treating that as the whole answer handed the engine the wrong
// family's flags: flux with clip_l for Chroma, sd3 for SD3.5, sd1.5 for SD2.1.
//
// The refinement is deliberately narrow. It never crosses an architecture: a
// file called "flux" whose weights are SDXL stays SDXL, and an unnamed file
// keeps the weight verdict.
func refineFamilyFromName(probed, named ModelFamily) ModelFamily {
	switch probed {
	case FamilySD15:
		if named == FamilySD21 {
			return FamilySD21
		}
	case FamilySD3:
		if named == FamilySD35 {
			return FamilySD35
		}
	case FamilyFlux:
		if named == FamilyChroma {
			return FamilyChroma
		}
	case FamilyQwenImage:
		// The one refinement that crosses the header tag, because the loader
		// overrides it: Krea2 files declare qwen_image and run as Krea2.
		if named == FamilyKrea2 {
			return FamilyKrea2
		}
	}
	return probed
}

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
	// CLIP-Vision is an image encoder, not a text encoder. It must be tested
	// before the generic clip_l/clip_g rules and before ControlNet.
	case contains(lower, "clip_vision", "clipvision", "clip-vision", "image_encoder", "image-encoder"):
		return AssetClipVision
	case contains(lower, "ip-adapter", "ipadapter", "ip_adapter"):
		return AssetIPAdapter
	case contains(lower, "controlnet", "control_net", "control-net", "control_v11", "control-v11"):
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
		case contains(lower, "qwen_image", "qwen-image", "qwenimage"):
			// Anima and Qwen-Image ship the same autoencoder. Without this the
			// conservative default below labels it sd1.5, and Anima then refuses
			// the only VAE that can decode it.
			return FamilyQwenImage
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

	// CLIP-Vision H is shared by SD1.5 and SDXL IP-Adapters; it is not a
	// text-encoder family. Leave it unknown so compatibility is permissive.
	if contains(lower, "clip_vision", "clipvision", "clip-vision", "image_encoder", "image-encoder") {
		return FamilyUnknown
	}

	// IP-Adapter family is the adapter family, not the base model's filename.
	if contains(lower, "ip-adapter", "ipadapter", "ip_adapter") {
		switch {
		case contains(lower, "sdxl", "sd_xl", "sd-xl", "xl"):
			return FamilySDXL
		case contains(lower, "sd15", "sd_1.5", "sd-1.5", "1.5"):
			return FamilySD15
		}
		return FamilyUnknown
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
	case contains(lower, "anima"):
		return FamilyAnima
	// Flux.2 and Flux.1 differ in text encoder and block layout, so the
	// more specific name has to be tested first.
	case contains(lower, "flux2", "flux.2", "flux-2", "flux_2"):
		return FamilyFlux2
	case contains(lower, "z-image", "z_image", "zimage"):
		return FamilyZImage
	case contains(lower, "krea2", "krea-2", "krea_2", "krea2_turbo"):
		return FamilyKrea2
	case contains(lower, "qwen-image", "qwen_image", "qwenimage"):
		return FamilyQwenImage
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
	case FamilyFlux, FamilyFlux2, FamilyZImage, FamilyQwenImage, FamilyAnima, FamilySDXL, FamilySD35, FamilySD3, FamilyTurbo:
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
	case FamilyFlux, FamilyFlux2, FamilyZImage, FamilyQwenImage, FamilyAnima, FamilyWan, FamilyLTX, FamilyHunyuan, FamilyAnimate:
		return "euler"
	case FamilySDXL, FamilySD35, FamilySD3, FamilyTurbo:
		return "dpmpp_2m"
	default:
		return "euler_a"
	}
}

// DefaultSteps returns the step count a family is built for. A distilled model
// run at the generic 20 steps and CFG 7 does not get better — it burns the image,
// and a burned image is indistinguishable from a broken engine or a bad offload.
func DefaultSteps(family ModelFamily) int {
	switch family {
	case FamilyZImage:
		return 8 // Z-Image-Turbo: 9 scheduler steps, 8 DiT forwards
	case FamilyFlux:
		return 4 // FLUX.1-schnell
	case FamilyFlux2:
		return 20
	case FamilyQwenImage:
		return 20
	case FamilyAnima:
		return 30
	default:
		return 20
	}
}

// DefaultCFG returns the classifier-free guidance a family expects. Distilled
// models bake guidance in and must run at 1.0; anything else oversaturates.
func DefaultCFG(family ModelFamily) float64 {
	switch family {
	case FamilyZImage, FamilyFlux:
		return 1.0
	case FamilyQwenImage:
		return 2.5
	case FamilyAnima:
		return 6.0
	case FamilyFlux2:
		return 4.0
	default:
		return 7.0
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
	// Anima ships the Qwen-Image autoencoder (circlestone-labs/Anima), so the two
	// families interchange here; an sd1.5 or Flux VAE is still rejected.
	if (modelFamily == FamilyAnima && vaeFamily == FamilyQwenImage) ||
		(modelFamily == FamilyQwenImage && vaeFamily == FamilyAnima) {
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

// CompatibleIPAdapter returns true only for the SD1.5 and SDXL families this
// integration supports. An unknown adapter family is not accepted: guessing
// would load weights that cannot project into the selected UNet.
func CompatibleIPAdapter(modelFamily, adapterFamily ModelFamily) bool {
	if adapterFamily != FamilySD15 && adapterFamily != FamilySDXL {
		return false
	}
	return modelFamily == adapterFamily
}

// CompatibleClipVision returns true for the universal ViT-H encoder. A known
// family must match the adapter family; unknown is accepted because the encoder
// itself has no SD-family metadata.
func CompatibleClipVision(modelFamily, visionFamily ModelFamily) bool {
	if visionFamily == FamilyUnknown {
		return true
	}
	return modelFamily == visionFamily
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
