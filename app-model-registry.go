package main

// app-model-registry.go
//
// Wails bindings for the unified model registry:
//   - ScanLLMDir    — walk a directory for .gguf/.bin files, register + remember it
//   - ScanTTSDir    — walk a directory for .onnx piper voice files, register + remember it
//   - GetLLMModels  — return all registered LLM models from any directory
//   - GetTTSModels  — return all registered TTS voice models from any directory
//   - RescanAllDirs — re-walk every remembered directory to pick up newly added files
//
// On startup, RescanAllDirs is called automatically from initImageEngine so new
// models dropped into any previously registered directory appear without any
// manual action from the user.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"howl-chat/internal/backend/gguf"
	"howl-chat/internal/errors"
	"howl-chat/internal/image/registry"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ---- public model info types ------------------------------------------------

// LLMModelInfo describes a registered LLM (chat) model.
type LLMModelInfo struct {
	Name   string  `json:"name"`
	Path   string  `json:"path"`
	Format string  `json:"format"`
	SizeMB float64 `json:"size_mb"`
}

// LLMModelDetail provides enriched metadata + memory estimates for a single model.
type LLMModelDetail struct {
	Name           string  `json:"name"`
	Path           string  `json:"path"`
	Format         string  `json:"format"`
	SizeMB         float64 `json:"size_mb"`
	Quant          string  `json:"quant"`
	ParameterCount string  `json:"parameter_count"`
	LayerCount     int     `json:"layer_count"`
	MaxContext     int     `json:"max_context"`
	Family         string  `json:"family"`
	GPUMB          int64   `json:"gpu_mb"`
	RAMMB          int64   `json:"ram_mb"`
	Notes          string  `json:"notes"`
}

// TTSModelInfo describes a registered TTS voice model.
type TTSModelInfo struct {
	Name   string  `json:"name"`
	Path   string  `json:"path"`
	Format string  `json:"format"`
	SizeMB float64 `json:"size_mb"`
}

// ---- Wails bindings ---------------------------------------------------------

// ScanLLMDir walks dir for .gguf and .bin files, registers them in the
// persistent registry, and remembers the directory for future auto-rescans.
func (a *App) ScanLLMDir(dir string) ([]LLMModelInfo, error) {
	if a.modelRegistry == nil {
		return nil, fmt.Errorf("model registry not initialised")
	}
	entries, err := scanDirForType(dir, "llm", []string{".gguf", ".bin", ".mmproj"})
	if err != nil {
		return nil, err
	}
	if err := a.modelRegistry.MergeEntries(entries); err != nil {
		a.logError(errors.REGModelAddFail, errors.High, "registry.llm.merge", "app-model-registry.go:74", err)
		return nil, err
	}
	if err := a.modelRegistry.RegisterDir(dir, "llm"); err != nil {
		a.logErrorSilent(errors.REGModelAddFail, errors.Medium, "registry.llm.dir", "app-model-registry.go:77", err)
	}
	return llmInfos(a.modelRegistry.ByType("llm")), nil
}

// ScanTTSDir walks dir for Piper .onnx voice model files, registers them,
// and remembers the directory for future auto-rescans.
func (a *App) ScanTTSDir(dir string) ([]TTSModelInfo, error) {
	if a.modelRegistry == nil {
		return nil, fmt.Errorf("model registry not initialised")
	}
	entries, err := scanDirForType(dir, "tts", []string{".onnx"})
	if err != nil {
		return nil, err
	}
	if err := a.modelRegistry.MergeEntries(entries); err != nil {
		a.logError(errors.REGModelAddFail, errors.High, "registry.tts.merge", "app-model-registry.go:93", err)
		return nil, err
	}
	if err := a.modelRegistry.RegisterDir(dir, "tts"); err != nil {
		a.logErrorSilent(errors.REGModelAddFail, errors.Medium, "registry.tts.dir", "app-model-registry.go:96", err)
	}
	return ttsInfos(a.modelRegistry.ByType("tts")), nil
}

// GetLLMModels returns all registered LLM models from every known directory,
// pruning any whose files have been deleted from disk.
func (a *App) GetLLMModels() ([]LLMModelInfo, error) {
	if a.modelRegistry == nil {
		return []LLMModelInfo{}, fmt.Errorf("model registry not initialised")
	}
	return llmInfos(a.modelRegistry.ByType("llm")), nil
}

// GetLLMModelDetails returns enriched metadata and memory estimates for a single model path.
func (a *App) GetLLMModelDetails(path string) (*LLMModelDetail, error) {
	profile, err := gguf.ExtractModelProfile(path)
	if err != nil {
		a.logErrorSilent(errors.LLMGGUFReadFail, errors.Medium, "registry.llm.detail", "app-model-registry.go:113", err)
		return nil, fmt.Errorf("failed to read model metadata: %w", err)
	}

	gpuLayers := 0
	ctxSize := 4096
	if a.modelSettings != nil {
		if gl, ok := a.modelSettings["gpu_layers"].(int); ok {
			gpuLayers = gl
		}
		if cs, ok := a.modelSettings["context_size"].(int); ok {
			ctxSize = cs
		}
	}
	est := gguf.EstimateMemory(profile, gpuLayers, ctxSize)

	paramStr := ""
	if profile.ParameterCount > 0 {
		b := float64(profile.ParameterCount)
		switch {
		case b >= 1_000_000_000:
			paramStr = fmt.Sprintf("%.1fB", b/1_000_000_000)
		case b >= 1_000_000:
			paramStr = fmt.Sprintf("%.1fM", b/1_000_000)
		default:
			paramStr = fmt.Sprintf("%.0f", b)
		}
	}

	return &LLMModelDetail{
		Name:           filepath.Base(path),
		Path:           path,
		Format:         filepath.Ext(path),
		SizeMB:         float64(profile.FileSizeBytes) / (1024 * 1024),
		Quant:          profile.FileType,
		ParameterCount: paramStr,
		LayerCount:     profile.LayerCount,
		MaxContext:     profile.MaxContext,
		Family:         profile.Family,
		GPUMB:          est.GPUMB,
		RAMMB:          est.RAMMB,
		Notes:          est.Notes,
	}, nil
}

// GetTTSModels returns all registered TTS voice models from every known
// directory, pruning any whose files have been deleted from disk.
func (a *App) GetTTSModels() ([]TTSModelInfo, error) {
	if a.modelRegistry == nil {
		return []TTSModelInfo{}, fmt.Errorf("model registry not initialised")
	}
	return ttsInfos(a.modelRegistry.ByType("tts")), nil
}

// RescanAllDirs re-walks every previously registered directory and merges
// any new files into the registry. Called automatically on startup.
// Non-fatal: errors per directory are logged, not returned.
func (a *App) RescanAllDirs() {
	if a.modelRegistry == nil {
		return
	}
	dirs := a.modelRegistry.RegisteredDirs()
	for dir, assetType := range dirs {
		if _, err := os.Stat(dir); err != nil {
			continue // directory no longer exists — skip silently
		}
		var exts []string
		switch assetType {
		case "llm":
			exts = []string{".gguf", ".bin", ".mmproj"}
		case "tts":
			exts = []string{".onnx"}
		case "image":
			// Image directories are rescanned via ScanImageModels; skip here.
			continue
		default:
			continue
		}
		entries, err := scanDirForType(dir, assetType, exts)
		if err != nil {
			a.logErrorSilent(errors.REGModelListFail, errors.Medium, "registry.rescan.walk", "app-model-registry.go:192", err)
			continue
		}
		if err := a.modelRegistry.MergeEntries(entries); err != nil {
			a.logErrorSilent(errors.REGModelAddFail, errors.Medium, "registry.rescan.merge", "app-model-registry.go:196", err)
		}
	}
}

// OpenDirectoryPicker opens a native directory-selection dialog and returns
// the chosen path. Used by both ScanLLMDir and ScanTTSDir flows in the UI.
func (a *App) OpenDirectoryPicker() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Model Directory",
	})
	if err != nil {
		return "", fmt.Errorf("directory dialog: %w", err)
	}
	return dir, nil
}

// ---- helpers ----------------------------------------------------------------

// scanDirForType walks dir recursively and returns registry.Entry values for
// every file whose extension is in exts. For .gguf files, it validates the
// architecture to prevent diffusion models (flux, sd3, etc.) from being added
// as LLMs, and skips .gguf files that are clearly LLM models when scanning
// for image assets.
func scanDirForType(dir, assetType string, exts []string) ([]registry.Entry, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot access directory %q: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", dir)
	}

	// Diffusion architectures that may be packaged as .gguf
	diffusionFamilies := map[string]bool{
		"flux": true, "flux1": true, "sd3": true, "sd3.5": true,
		"sdxl": true, "stable-diffusion": true, "stable_diffusion": true,
		"auraflow": true, "pixart": true, "lumina": true,
	}

	var entries []registry.Entry
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		for _, want := range exts {
			if ext != want {
				continue
			}

			info, statErr := d.Info()
			sizeMB := 0.0
			if statErr == nil {
				sizeMB = float64(info.Size()) / (1024 * 1024)
			}
			name := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))

			// For .gguf files, validate architecture to prevent cross-contamination
			if ext == ".gguf" || ext == ".bin" {
				profile, parseErr := gguf.ExtractModelProfile(path)
				if parseErr == nil && profile != nil && profile.Family != "" {
					family := strings.ToLower(profile.Family)
					// If scanning for LLM but file is a diffusion model, skip it
					if assetType == "llm" && diffusionFamilies[family] {
						return nil // skip diffusion .gguf in LLM scan
					}
					// If scanning for image but file is an LLM model, skip it
					if assetType == "image" && !diffusionFamilies[family] {
						return nil // skip LLM .gguf in image scan
					}
					// Use the actual family from metadata
					entries = append(entries, registry.Entry{
						Name:   name,
						Path:   path,
						Type:   assetType,
						Family: profile.Family,
						Format: strings.TrimPrefix(ext, "."),
						SizeMB: sizeMB,
					})
					return nil
				}
			}

			entries = append(entries, registry.Entry{
				Name:   name,
				Path:   path,
				Type:   assetType,
				Family: assetType,
				Format: strings.TrimPrefix(ext, "."),
				SizeMB: sizeMB,
			})
			break
		}
		return nil
	})
	return entries, err
}

func llmInfos(entries []registry.Entry) []LLMModelInfo {
	out := make([]LLMModelInfo, 0, len(entries))
	for _, e := range entries {
		out = append(out, LLMModelInfo{
			Name:   e.Name,
			Path:   e.Path,
			Format: e.Format,
			SizeMB: e.SizeMB,
		})
	}
	return out
}

func ttsInfos(entries []registry.Entry) []TTSModelInfo {
	out := make([]TTSModelInfo, 0, len(entries))
	for _, e := range entries {
		out = append(out, TTSModelInfo{
			Name:   e.Name,
			Path:   e.Path,
			Format: e.Format,
			SizeMB: e.SizeMB,
		})
	}
	return out
}
