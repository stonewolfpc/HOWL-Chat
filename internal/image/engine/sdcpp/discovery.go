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
	f := strings.ToLower(family)
	switch f {
	case "wan", "wan2.1", "wan2.2", "wan21", "wan22", "animatediff":
		return "832x480"
	case "ltx", "ltx2", "ltx-video":
		return "768x512"
	case "hunyuan", "hunyuanvideo", "hunyuan-video":
		return "848x480"
	case "minimax", "minimax-h3":
		return "768x432"
	}
	// Image architectures read their base resolution from the capability table.
	if res := FamilyResolution(f); res > 0 {
		return fmt.Sprintf("%dx%d", res, res)
	}
	return "512x512"
}
