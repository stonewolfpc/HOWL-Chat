package sdcpp

import (
	"fmt"
	"strings"

	"howl-chat/internal/image/scanner"
	"howl-chat/internal/image/types"
)

// ListModels runs the scanner and returns ModelInfo for every diffusion model found.
// Auto-fills companion files (VAE, CLIP, T5XXL, LLM) from all scanned directories.
// Filters out models that cannot load due to missing required companions.
func (e *Engine) ListModels() ([]types.ModelInfo, error) {
	fmt.Printf("[ImageEngine] ListModels scanning: %s\n", e.config.ModelsDir)
	s := scanner.New(e.config.ModelsDir)
	idx, err := s.Scan()
	if err != nil {
		return []types.ModelInfo{}, err
	}

	var out []types.ModelInfo
	for _, m := range idx.Models {
		info := entryToModelInfo(m)
		// Auto-fill missing companion files from the full index
		autoFillCompanions(&info, idx)
		// Only include models that are loadable (have required companions)
		if isModelLoadable(&info, idx) {
			out = append(out, info)
		} else {
			fmt.Printf("[ImageEngine] Skipping model '%s' (family: %s) - missing required companions\n",
				info.Name, info.Family)
		}
	}
	fmt.Printf("[ImageEngine] ListModels returning %d models (filtered from %d total)\n", len(out), len(idx.Models))
	return out, nil
}

// isModelLoadable checks if a model has all required companions for its family.
// Returns false if the model cannot load due to missing dependencies.
// Marks unknown family models as unsupported for user clarity.
// SD3.5 models are filtered out if the backend doesn't support them.
func isModelLoadable(info *types.ModelInfo, idx scanner.ModelIndex) bool {
	family := strings.ToLower(info.Family)
	files := info.Files

	// SD3.5 GGUF is currently broken in stable-diffusion.cpp (issue #1105)
	// Safetensors version works, GGUF crashes on load
	if family == "sd3.5" && info.Format == "gguf" {
		info.Compatibility = "unsupported: SD3.5 GGUF is currently broken in stable-diffusion.cpp. Use the .safetensors version instead."
		return false
	}

	switch family {
	case "flux":
		// Flux requires T5XXL, CLIP-L, and VAE
		if files["t5xxl"] == "" && len(idx.T5XXLs) == 0 {
			return false
		}
		if files["clip_l"] == "" && len(idx.ClipLs) == 0 {
			return false
		}
		if files["vae"] == "" && len(idx.VAEs) == 0 {
			return false
		}
	case "chroma":
		// Chroma only requires T5XXL and VAE (CLIP-L not used)
		if files["t5xxl"] == "" && len(idx.T5XXLs) == 0 {
			return false
		}
		if files["vae"] == "" && len(idx.VAEs) == 0 {
			return false
		}
	case "sd3", "sd3.5":
		// SD3/SD3.5 require CLIP-L, CLIP-G, and T5XXL
		if files["clip_l"] == "" && len(idx.ClipLs) == 0 {
			return false
		}
		if files["clip_g"] == "" && len(idx.ClipGs) == 0 {
			return false
		}
		if files["t5xxl"] == "" && len(idx.T5XXLs) == 0 {
			return false
		}
	case "auraflow", "hidream":
		// AuraFlow/HiDream require T5XXL and VAE
		if files["t5xxl"] == "" && len(idx.T5XXLs) == 0 {
			return false
		}
		if files["vae"] == "" && len(idx.VAEs) == 0 {
			return false
		}
	case "anima":
		// Anima requires LLM and VAE
		if files["llm"] == "" && len(idx.LLMs) == 0 {
			return false
		}
		if files["vae"] == "" && len(idx.VAEs) == 0 {
			return false
		}
	case "wan", "wan2.1", "wan2.2", "wan21", "wan22", "ltx", "ltx2", "ltx-video",
		"hunyuan", "hunyuanvideo", "hunyuan-video", "minimax", "minimax-h3":
		// Video families need a text encoder and their matching VAE.
		if files["t5xxl"] == "" && len(idx.T5XXLs) == 0 {
			return false
		}
		if files["vae"] == "" && len(idx.VAEs) == 0 {
			return false
		}
	case "":
		// Unknown family - include but mark as unsupported for user awareness
		info.Compatibility = "unsupported: unknown model family"
		return true
	default:
		// SD1.5, SD2.1, SDXL have baked-in encoders - always loadable
		return true
	}
	return true
}

// autoFillCompanions searches the full index for companion files when a model is missing them.
// This enables models to find their companions across multiple scanned directories.
func autoFillCompanions(info *types.ModelInfo, idx scanner.ModelIndex) {
	family := strings.ToLower(info.Family)

	// Debug: log what we're working with
	fmt.Printf("[ImageEngine] autoFillCompanions for '%s' (family: %s, files: %v)\n", info.Name, info.Family, info.Files)
	fmt.Printf("[ImageEngine] Available - VAEs: %d, ClipL: %d, ClipG: %d, T5XXL: %d, LLMs: %d\n",
		len(idx.VAEs), len(idx.ClipLs), len(idx.ClipGs), len(idx.T5XXLs), len(idx.LLMs))

	// Determine what companion files this model family needs
	needsVAE := false // Default: most models have baked-in VAEs
	needsClipL := true
	needsClipG := false
	needsT5XXL := false
	needsLLM := false

	switch family {
	case "flux":
		needsVAE = true // Flux requires external VAE
		needsT5XXL = true
		needsClipL = true
	case "chroma":
		// Chroma requires external VAE (Flux VAE)
		needsVAE = true
		needsT5XXL = true
		needsClipL = false
		needsClipG = false
		// Explicitly clear CLIP encoders from files map
		info.Files["clip_l"] = ""
		info.Files["clip_g"] = ""
	case "sd3", "sd3.5":
		needsVAE = true // SD3/SD3.5 require external VAE
		needsClipL = true
		needsClipG = true
		needsT5XXL = true
	case "sdxl":
		needsClipL = true
		// SDXL models may or may not need VAE - leave to user choice
	case "sd1.5", "sd2.1":
		needsClipL = true
		// SD1.5/SD2.1 have baked-in VAEs - don't auto-fill
		// SD1.5/SD2.1 don't use T5XXL
		info.Files["t5xxl"] = ""
	case "anima":
		needsLLM = true
		needsVAE = true // qwen_image_vae
	case "wan", "wan2.1", "wan2.2", "wan21", "wan22", "ltx", "ltx2", "ltx-video",
		"hunyuan", "hunyuanvideo", "hunyuan-video", "minimax", "minimax-h3":
		needsVAE = true
		needsT5XXL = true
		needsClipL = false
		needsClipG = false
	case "hidream", "qwen_image":
		// Self-contained, no companions
		needsVAE = false
		needsClipL = false
	}

	// Auto-fill VAE if missing (only if we find a compatible one)
	if needsVAE && info.Files["vae"] == "" && len(idx.VAEs) > 0 {
		// Try to find a compatible VAE by family
		for _, vae := range idx.VAEs {
			// Skip Flux VAEs (ae.safetensors) for non-Flux/Chroma models
			// Check both name and path for ae.safetensors or files starting with ae.
			vaeLower := strings.ToLower(vae.Name)
			pathLower := strings.ToLower(vae.Path)
			fmt.Printf("[DEBUG] Checking VAE: name=%s path=%s family=%s\n", vae.Name, vae.Path, vae.Family)
			if strings.HasPrefix(vaeLower, "ae.") ||
				strings.HasPrefix(pathLower, "ae.") ||
				vaeLower == "ae.safetensors" ||
				pathLower == "ae.safetensors" ||
				strings.Contains(vaeLower, "ae.safetensors") ||
				strings.Contains(pathLower, "ae.safetensors") {
				fmt.Printf("[DEBUG] Skipping Flux VAE: %s\n", vae.Name)
				continue
			}
			if scanner.CompatibleVAE(scanner.ModelFamily(family), vae.Family) {
				fmt.Printf("[DEBUG] Selected VAE: %s (family: %s)\n", vae.Name, vae.Family)
				info.Files["vae"] = vae.Path
				break
			}
		}
		// If no compatible VAE found, leave empty (model may have baked-in VAE)
	}

	// Auto-fill CLIP-L if missing
	if needsClipL && info.Files["clip_l"] == "" && len(idx.ClipLs) > 0 {
		info.Files["clip_l"] = idx.ClipLs[0].Path
	}

	// Auto-fill CLIP-G if missing
	if needsClipG && info.Files["clip_g"] == "" && len(idx.ClipGs) > 0 {
		info.Files["clip_g"] = idx.ClipGs[0].Path
	}

	// Auto-fill T5XXL if missing
	if needsT5XXL && info.Files["t5xxl"] == "" && len(idx.T5XXLs) > 0 {
		info.Files["t5xxl"] = idx.T5XXLs[0].Path
	}

	// Auto-fill LLM if missing (for Anima)
	if needsLLM && info.Files["llm"] == "" && len(idx.LLMs) > 0 {
		info.Files["llm"] = idx.LLMs[0].Path
		fmt.Printf("[ImageEngine] Auto-filled LLM for %s: %s\n", info.Name, idx.LLMs[0].Path)
	}

	// Debug logging for Anima models
	if family == "anima" {
		fmt.Printf("[ImageEngine] Anima model '%s' companions - llm: %v, vae: %v\n",
			info.Name, info.Files["llm"] != "", info.Files["vae"] != "")
	}
}

// ListLoRAs returns the names of all LoRA files found in the models directory.
func (e *Engine) ListLoRAs() ([]string, error) {
	s := scanner.New(e.config.ModelsDir)
	idx, err := s.Scan()
	if err != nil {
		return []string{}, err
	}
	var out []string
	for _, l := range idx.LoRAs {
		out = append(out, l.Name)
	}
	return out, nil
}

// entryToModelInfo converts a scanner.AssetEntry to the types.ModelInfo used by the rest of the app.
func entryToModelInfo(e scanner.AssetEntry) types.ModelInfo {
	caps := []string{"txt2img", "img2img", "inpaint"}
	if IsVideoFamily(string(e.Family)) {
		caps = []string{"txt2video", "img2video"}
	}
	res := fmt.Sprintf("%dx%d", e.Recommended.Width, e.Recommended.Height)
	if e.Recommended.Width == 0 || IsVideoFamily(string(e.Family)) {
		res = defaultResolutionStr(string(e.Family))
	}
	format := e.Format
	compat := "unknown"
	switch format {
	case "safetensors", "ckpt", "sft":
		compat = "verified"
	case "gguf":
		compat = "supported" // merged GGUF checkpoints are supported; split FLUX-style use --diffusion-model
	case "bundle":
		compat = "supported"
	}
	files := e.Files
	if files == nil {
		files = map[string]string{"model": e.Path}
	}
	return types.ModelInfo{
		Name:           e.Name,
		Path:           e.Path,
		Family:         string(e.Family),
		Capabilities:   caps,
		ResolutionHint: res,
		Files:          files,
		Format:         format,
		Compatibility:  compat,
	}
}

// inferFamilyFromName returns a best-guess model family string from a filename.
// Mirrors scanner.GuessFamily without importing the scanner package.
func inferFamilyFromName(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "wan2"), strings.Contains(lower, "wan_2"), strings.Contains(lower, "wan-"), strings.Contains(lower, "-wan"):
		return "wan"
	case strings.Contains(lower, "ltx"):
		return "ltx2"
	case strings.Contains(lower, "hunyuan"), strings.Contains(lower, "hunyuanvideo"):
		return "hunyuanvideo"
	case strings.Contains(lower, "animatediff"), strings.Contains(lower, "motion-module"), strings.Contains(lower, "motion_module"):
		return "animatediff"
	case strings.Contains(lower, "minimax"):
		return "minimax-h3"
	case strings.Contains(lower, "flux"):
		return "flux"
	case strings.Contains(lower, "sdxl-turbo"), strings.Contains(lower, "sd_xl_turbo"):
		return "sdxl-turbo"
	case strings.Contains(lower, "chroma"):
		return "chroma"
	case strings.Contains(lower, "anima"):
		return "anima"
	case strings.Contains(lower, "auraflow"), strings.Contains(lower, "aura-flow"), strings.Contains(lower, "aura_flow"):
		return "auraflow"
	case strings.Contains(lower, "hidream"), strings.Contains(lower, "hi-dream"), strings.Contains(lower, "hi_dream"):
		return "hidream"
	case strings.Contains(lower, "sd3.5"), strings.Contains(lower, "sd_3.5"), strings.Contains(lower, "sd35"), strings.Contains(lower, "stable-diffusion-3.5"):
		return "sd3.5"
	case strings.Contains(lower, "sdxl"), strings.Contains(lower, "sd_xl"),
		strings.Contains(lower, "juggernaut_xl"), strings.Contains(lower, "juggernaut-xl"),
		strings.Contains(lower, "realisticvision_xl"), strings.Contains(lower, "realisticvision-xl"),
		strings.Contains(lower, "dreamshaperxl"), strings.Contains(lower, "dreamshaper-xl"),
		strings.Contains(lower, "epicrealism"):
		return "sdxl"
	case strings.Contains(lower, "sd3"), strings.Contains(lower, "stable-diffusion-3"):
		return "sd3"
	case strings.Contains(lower, "sd2.1"), strings.Contains(lower, "sd21"):
		return "sd2.1"
	default:
		return "sd1.5"
	}
}

func defaultResolutionStr(family string) string {
	switch strings.ToLower(family) {
	case "wan", "wan2.1", "wan2.2", "wan21", "wan22", "animatediff":
		return "832x480"
	case "ltx", "ltx2", "ltx-video":
		return "768x512"
	case "hunyuan", "hunyuanvideo", "hunyuan-video":
		return "848x480"
	case "minimax", "minimax-h3":
		return "768x432"
	case "flux", "sdxl", "sd3.5", "sd3", "sdxl-turbo", "chroma", "anima", "kolors", "auraflow", "hidream":
		return "1024x1024"
	case "sd2.1", "cascade":
		return "768x768"
	default:
		return "512x512"
	}
}
