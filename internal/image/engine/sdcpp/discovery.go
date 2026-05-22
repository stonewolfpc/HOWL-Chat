package sdcpp

import (
	"fmt"
	"strings"

	"howl-chat/internal/image/scanner"
	"howl-chat/internal/image/types"
)

// ListModels runs the scanner and returns ModelInfo for every diffusion model found.
func (e *Engine) ListModels() ([]types.ModelInfo, error) {
	fmt.Printf("[ImageEngine] ListModels scanning: %s\n", e.config.ModelsDir)
	s := scanner.New(e.config.ModelsDir)
	idx, err := s.Scan()
	if err != nil {
		return []types.ModelInfo{}, err
	}
	var out []types.ModelInfo
	for _, m := range idx.Models {
		out = append(out, entryToModelInfo(m))
	}
	fmt.Printf("[ImageEngine] ListModels returning %d models\n", len(out))
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
	res := fmt.Sprintf("%dx%d", e.Recommended.Width, e.Recommended.Height)
	if e.Recommended.Width == 0 {
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
	// Chroma architecture (SDXL-derivative) is not yet recognised by sd-server ≤0665a7f
	if e.Family == "chroma" {
		compat = "unsupported"
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
	case strings.Contains(lower, "flux"):
		return "flux"
	case strings.Contains(lower, "sdxl-turbo"), strings.Contains(lower, "sd_xl_turbo"):
		return "sdxl-turbo"
	case strings.Contains(lower, "chroma"):
		return "chroma"
	case strings.Contains(lower, "sdxl"), strings.Contains(lower, "sd_xl"):
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
	case "flux", "sdxl", "sd3", "sdxl-turbo", "chroma", "kolors", "auraflow":
		return "1024x1024"
	case "sd2.1", "cascade":
		return "768x768"
	default:
		return "512x512"
	}
}
