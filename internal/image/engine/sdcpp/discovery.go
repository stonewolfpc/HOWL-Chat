package sdcpp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"howl-chat/internal/image/types"
)

// ListModels scans the models directory and detects model families.
func (e *Engine) ListModels() ([]types.ModelInfo, error) {
	fmt.Printf("[ImageEngine] ListModels scanning: %s\n", e.config.ModelsDir)
	entries, err := os.ReadDir(e.config.ModelsDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("[ImageEngine] ModelsDir does not exist: %s\n", e.config.ModelsDir)
			return []types.ModelInfo{}, nil
		}
		return []types.ModelInfo{}, fmt.Errorf("read models dir: %w", err)
	}
	fmt.Printf("[ImageEngine] Found %d entries in models dir\n", len(entries))

	var models []types.ModelInfo
	seen := make(map[string]bool)

	// First: scan subdirectories (multi-file bundles like FLUX)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		subDir := filepath.Join(e.config.ModelsDir, entry.Name())
		info, err := e.scanModelFolder(entry.Name(), subDir)
		if err != nil {
			fmt.Printf("[ImageEngine] Skipping folder %q: %v\n", entry.Name(), err)
			continue
		}
		if info != nil {
			fmt.Printf("[ImageEngine] Found model folder: %s (%s)\n", info.Name, info.Family)
			models = append(models, *info)
			seen[info.Name] = true
		}
	}

	// Second: scan root-level model files (single-file models)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		lower := strings.ToLower(name)
		fmt.Printf("[ImageEngine] Checking file: %q (lower=%q)\n", name, lower)
		if !strings.HasSuffix(lower, ".gguf") && !strings.HasSuffix(lower, ".safetensors") &&
			!strings.HasSuffix(lower, ".ckpt") && !strings.HasSuffix(lower, ".pth") &&
			!strings.HasSuffix(lower, ".pt") && !strings.HasSuffix(lower, ".bin") {
			fmt.Printf("[ImageEngine]   -> skipped: not a recognized model extension\n")
			continue
		}
		base := strings.TrimSuffix(name, filepath.Ext(name))
		if seen[base] {
			fmt.Printf("[ImageEngine]   -> skipped: already seen as %q\n", base)
			continue
		}

		// Size-based classification: anything under 800MB is not a full model
		info, _ := entry.Info()
		sizeMB := float64(info.Size()) / (1024 * 1024)
		classification := classifyFile(name, info.Size())
		fmt.Printf("[ImageEngine]   -> size=%.1f MB, classified as %s\n", sizeMB, classification)

		switch classification {
		case "embedding":
			fmt.Printf("[ImageEngine]   -> skipped: embedding (< 50 MB)\n")
			continue
		case "LoRA":
			fmt.Printf("[ImageEngine]   -> skipped: LoRA (50-300 MB or name contains 'lora')\n")
			continue
		case "VAE":
			fmt.Printf("[ImageEngine]   -> skipped: VAE/adapter (300-800 MB or name contains 'vae/controlnet')\n")
			continue
		case "unknown_small":
			fmt.Printf("[ImageEngine]   -> skipped: unknown small file (< 800 MB, not a full model)\n")
			continue
		}

		// For files ≥800MB, allow even if name has 'vae' or 'lora' (some models include those words)
		family := guessFamily(base)
		if family == "" {
			fmt.Printf("[ImageEngine]   -> skipped: guessFamily returned empty (non-SD model)\n")
			continue
		}
		fmt.Printf("[ImageEngine]   -> ACCEPTED as %s (family=%s)\n", base, family)
		format, compat := detectFormatCompatibility(filepath.Join(e.config.ModelsDir, name))
		models = append(models, types.ModelInfo{
			Name:           base,
			Path:           filepath.Join(e.config.ModelsDir, name),
			Family:         family,
			Capabilities:   []string{"txt2img", "img2img", "inpaint"},
			ResolutionHint: defaultResolution(family),
			Files: map[string]string{
				"model": filepath.Join(e.config.ModelsDir, name),
			},
			Format:        format,
			Compatibility: compat,
		})
	}
	fmt.Printf("[ImageEngine] ListModels returning %d models\n", len(models))

	return models, nil
}

// scanModelFolder inspects a single model directory.
func (e *Engine) scanModelFolder(name, dir string) (*types.ModelInfo, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileMap := make(map[string]string)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fp := filepath.Join(dir, f.Name())
		lower := strings.ToLower(f.Name())
		if strings.HasPrefix(lower, "diffusion") || strings.HasPrefix(lower, "model") && strings.Contains(lower, "diffusion") {
			fileMap["diffusion_model"] = fp
		} else if strings.HasPrefix(lower, "clip_l") || strings.HasPrefix(lower, "clip.l") {
			fileMap["clip_l"] = fp
		} else if strings.HasPrefix(lower, "clip_g") || strings.HasPrefix(lower, "clip.g") {
			fileMap["clip_g"] = fp
		} else if strings.HasPrefix(lower, "t5") || strings.HasPrefix(lower, "t5xxl") {
			fileMap["t5xxl"] = fp
		} else if strings.HasPrefix(lower, "vae") {
			fileMap["vae"] = fp
		} else if strings.HasSuffix(lower, ".gguf") || strings.HasSuffix(lower, ".safetensors") ||
			strings.HasSuffix(lower, ".ckpt") || strings.HasSuffix(lower, ".pth") ||
			strings.HasSuffix(lower, ".pt") || strings.HasSuffix(lower, ".bin") {
			fileMap["model"] = fp
		}
	}

	// Detect family
	switch {
	case hasAll(fileMap, "diffusion_model", "clip_l", "t5xxl", "vae"):
		return &types.ModelInfo{
			Name:           name,
			Path:           dir,
			Family:         "flux",
			Capabilities:   []string{"txt2img", "img2img", "inpaint"},
			ResolutionHint: "1024x1024",
			Files:          fileMap,
		}, nil
	case hasAll(fileMap, "model"):
		family := guessFamily(name)
		if family == "" {
			return nil, nil // skip non-SD models
		}
		return &types.ModelInfo{
			Name:           name,
			Path:           dir,
			Family:         family,
			Capabilities:   []string{"txt2img", "img2img", "inpaint"},
			ResolutionHint: defaultResolution(family),
			Files:          fileMap,
		}, nil
	}
	return nil, nil // not a recognized model folder
}

func hasAll(m map[string]string, keys ...string) bool {
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			return false
		}
	}
	return true
}

func guessFamily(name string) string {
	lower := strings.ToLower(name)
	// Exclude known ASR / embedding / non-SD model names
	if strings.Contains(lower, "whisper") || strings.Contains(lower, "ggml-base") ||
		strings.Contains(lower, "ggml-small") || strings.Contains(lower, "ggml-medium") ||
		strings.Contains(lower, "ggml-large") || strings.Contains(lower, "espeak") ||
		strings.Contains(lower, "bert") || strings.Contains(lower, "embed") {
		return ""
	}
	switch {
	case strings.Contains(lower, "flux"):
		return "flux"
	case strings.Contains(lower, "sdxl"):
		return "sdxl"
	case strings.Contains(lower, "sd3"):
		return "sd3"
	case strings.Contains(lower, "sd1.5") || strings.Contains(lower, "sd_1.5"):
		return "sd1.5"
	default:
		return "sd1.5" // default assumption
	}
}

func defaultResolution(family string) string {
	switch family {
	case "flux", "sdxl", "sd3":
		return "1024x1024"
	case "sd15":
		return "512x512"
	default:
		return "1024x1024"
	}
}

// classifyFile uses size and name to determine if a file is a full model, LoRA, VAE, or embedding.
// Returns: "model", "LoRA", "VAE", "embedding", "unknown_small"
func classifyFile(name string, size int64) string {
	lower := strings.ToLower(name)
	mb := float64(size) / (1024 * 1024)

	// Name-based overrides (most reliable)
	if strings.Contains(lower, "lora") {
		return "LoRA"
	}
	if strings.Contains(lower, "controlnet") || strings.Contains(lower, "ip-adapter") {
		return "VAE"
	}
	if strings.Contains(lower, "embedding") || strings.Contains(lower, "textual") || strings.Contains(lower, "ti_") || strings.Contains(lower, "hypernetwork") {
		return "embedding"
	}
	if strings.Contains(lower, "neg") && mb < 2 {
		return "embedding"
	}

	// Size-based classification
	if mb < 50 {
		return "embedding"
	}
	if mb < 300 {
		return "LoRA"
	}
	if mb < 800 {
		// Could be VAE, LoRA, or small model - check name hints
		if strings.Contains(lower, "vae") || strings.Contains(lower, "upscaler") ||
			strings.Contains(lower, "esrgan") || strings.Contains(lower, "swinir") ||
			strings.Contains(lower, "4x") || strings.Contains(lower, "clip_vision") {
			return "VAE"
		}
		return "unknown_small"
	}
	return "model"
}

// isVAEOrUpscaler returns true if the filename looks like a VAE, upscaler, or embedding file.
func isVAEOrUpscaler(name string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, "vae") || strings.Contains(lower, "upscaler") ||
		strings.Contains(lower, "esrgan") || strings.Contains(lower, "swinir") ||
		strings.Contains(lower, "4x") || strings.Contains(lower, "embedding") ||
		strings.Contains(lower, "textual") || strings.Contains(lower, "ti_") ||
		strings.Contains(lower, "hypernetwork") || strings.Contains(lower, "controlnet") ||
		strings.Contains(lower, "clip_vision") || strings.Contains(lower, "ip-adapter")
}

// isLoRA returns true if the filename looks like a LoRA.
func isLoRA(name string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, "lora")
}

// ListLoRAs scans the models directory for LoRA files.
func (e *Engine) ListLoRAs() ([]string, error) {
	if e.config.ModelsDir == "" {
		return []string{}, nil
	}
	entries, err := os.ReadDir(e.config.ModelsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return []string{}, fmt.Errorf("read models dir: %w", err)
	}
	var loras []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		lower := strings.ToLower(name)
		if !strings.HasSuffix(lower, ".gguf") && !strings.HasSuffix(lower, ".safetensors") &&
			!strings.HasSuffix(lower, ".ckpt") {
			continue
		}
		if isLoRA(name) {
			loras = append(loras, strings.TrimSuffix(name, filepath.Ext(name)))
		}
	}
	return loras, nil
}

// detectFormatCompatibility examines a model file and returns (format, compatibility).
// .safetensors and .ckpt are native formats → "verified"
// .gguf requires specific converter metadata → "unknown" (may or may not work)
func detectFormatCompatibility(path string) (format, compat string) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".safetensors":
		return "safetensors", "verified"
	case ".ckpt", ".pth", ".pt":
		return "checkpoint", "verified"
	case ".gguf":
		// Quick header check: read first 64 bytes to look for GGUF magic and SD metadata
		if f, err := os.Open(path); err == nil {
			defer f.Close()
			header := make([]byte, 64)
			if n, _ := f.Read(header); n >= 4 {
				// GGUF magic is 'GGUF' at offset 0
				if string(header[0:4]) != "GGUF" {
					return "gguf", "incompatible"
				}
			}
		}
		return "gguf", "unknown"
	default:
		return ext, "unknown"
	}
}
