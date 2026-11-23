package main

// app-model-registry.go
//
// Wails bindings for the unified model registry:
//   - ScanLLMDir    — walk a directory for .gguf chat models, register + remember it
//   - ScanTTSDir    — walk a directory for .onnx piper voice files, register + remember it
//   - GetLLMModels  — return all registered LLM models from any directory
//   - GetTTSModels  — return all registered TTS voice models from any directory
//   - RescanAllDirs — re-walk every remembered directory to pick up newly added files
//
// On startup, RescanAllDirs runs as Level 0 background work so new
// models dropped into any previously registered directory appear without any
// manual action from the user.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"howl-chat/internal/backend/chatpolicy"
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

// ScanLLMDir walks dir for .gguf files, registers them in the persistent
// registry, and remembers the directory for future auto-rescans. Legacy .bin
// files are not registered as chat models; Whisper and other legacy assets
// stay on their own ASR/inventory routes.
func (a *App) ScanLLMDir(dir string) ([]LLMModelInfo, error) {
	if a.registries.Models() == nil {
		return nil, fmt.Errorf("model registry not initialised")
	}
	entries, seen, err := scanDirForType(dir, registry.KindLLM, []string{".gguf"}, a.registries.Models().Lookup)
	if err != nil {
		return nil, err
	}
	if _, _, err := a.registries.Models().SyncDir(dir, entries, seen); err != nil {
		a.logError(errors.REGModelAddFail, errors.High, "registry.llm.merge", "app-model-registry.go:74", err)
		return nil, err
	}
	if err := a.registries.Models().RegisterDir(dir, registry.KindLLM); err != nil {
		a.logErrorSilent(errors.REGModelAddFail, errors.Medium, "registry.llm.dir", "app-model-registry.go", err)
	}
	return llmInfos(a.registries.Models().ByTypeInDirs("llm", registry.KindLLM)), nil
}

// ScanTTSDir walks dir for Piper .onnx voice model files, registers them,
// and remembers the directory for future auto-rescans.
func (a *App) ScanTTSDir(dir string) ([]TTSModelInfo, error) {
	if a.registries.Models() == nil {
		return nil, fmt.Errorf("model registry not initialised")
	}
	entries, seen, err := scanDirForType(dir, registry.KindTTS, []string{".onnx"}, a.registries.Models().Lookup)
	if err != nil {
		return nil, err
	}
	if _, _, err := a.registries.Models().SyncDir(dir, entries, seen); err != nil {
		a.logError(errors.REGModelAddFail, errors.High, "registry.tts.merge", "app-model-registry.go:93", err)
		return nil, err
	}
	if err := a.registries.Models().RegisterDir(dir, registry.KindTTS); err != nil {
		a.logErrorSilent(errors.REGModelAddFail, errors.Medium, "registry.tts.dir", "app-model-registry.go", err)
	}
	return ttsInfos(a.registries.Models().ByTypeInDirs("tts", registry.KindTTS)), nil
}

// GetRegisteredLLMDirs returns all directories currently registered for LLM scanning.
func (a *App) GetRegisteredLLMDirs() ([]string, error) {
	if a.registries.Models() == nil {
		return []string{}, fmt.Errorf("model registry not initialised")
	}
	return a.registries.Models().Directories(registry.KindLLM), nil
}

// GetLLMModels returns all registered LLM models from every known directory,
// pruning any whose files have been deleted from disk.
func (a *App) GetLLMModels() ([]LLMModelInfo, error) {
	if a.registries.Models() == nil {
		return []LLMModelInfo{}, fmt.Errorf("model registry not initialised")
	}
	return llmInfos(a.registries.Models().ByTypeInDirs("llm", registry.KindLLM)), nil
}

// RemoveLLMModel removes a single model entry from the registry by its path.
// Does NOT delete the underlying file from disk.
func (a *App) RemoveLLMModel(path string) error {
	if a.registries.Models() == nil {
		return fmt.Errorf("model registry not initialised")
	}
	return a.registries.Models().Remove(path)
}

// UnregisterLLMDir removes a directory from the LLM scan list and also
// removes all registered models whose path lives inside that directory.
func (a *App) UnregisterLLMDir(dir string) error {
	if a.registries.Models() == nil {
		return fmt.Errorf("model registry not initialised")
	}
	return a.registries.Models().UnregisterDir(dir)
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
	if a.settings.Model() != nil {
		if gl, ok := a.settings.Model()["gpu_layers"].(int); ok {
			gpuLayers = gl
		}
		if cs, ok := a.settings.Model()["context_size"].(int); ok {
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
	if a.registries.Models() == nil {
		return []TTSModelInfo{}, fmt.Errorf("model registry not initialised")
	}
	return ttsInfos(a.registries.Models().ByTypeInDirs("tts", registry.KindTTS)), nil
}

// RescanAllDirs re-walks every previously registered directory and merges
// any new files into the registry. Called automatically on startup.
// Non-fatal: errors per directory are logged, not returned.
func (a *App) RescanAllDirs() {
	if a.registries.Models() == nil {
		return
	}
	// Image directories keep their persisted entries across restarts and are
	// rescanned on demand (ScanImageModels / RefreshModelDirectories), which
	// keeps startup cheap for users with very large model folders.
	startupMark("rescan: start")
	for _, kind := range []string{registry.KindLLM, registry.KindTTS} {
		for _, dir := range a.registries.Models().Directories(kind) {
			if _, err := os.Stat(dir); err != nil {
				continue // drive unplugged or directory gone
			}
			startupMark("rescan: scanning " + dir)
			atomic.StoreInt64(&scanLiveParses, 0)
			n, scanErr := a.scanDirectoryIntoRegistry(dir, kind)
			if scanErr != nil {
				a.logErrorSilent(errors.REGModelListFail, errors.Medium, "registry.rescan.walk", "app-model-registry.go", scanErr)
			}
			startupMark(fmt.Sprintf("rescan: done %s %d entries parsed=%d", dir, n, atomic.LoadInt64(&scanLiveParses)))
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

// ---- admission cache --------------------------------------------------------

// scanLiveParses counts GGUF headers parsed during scanning. It is diagnostic
// only: the startup log reports it per directory so a regression back to
// re-parsing unchanged files is visible instead of silent.
var scanLiveParses int64

// entryAdmission rebuilds the admission decision recorded when the file was
// classified. A miss means the entry predates the cache or was classified by an
// older rule set, and the caller must read the file exactly once.
func entryAdmission(e registry.Entry) (gguf.Admission, bool) {
	if e.AdmissionV != registry.AdmissionCacheVersion || e.AdmissionClass == "" {
		return gguf.Admission{}, false
	}
	return gguf.Admission{
		Path:            e.Path,
		Class:           gguf.AdmissionClass(e.AdmissionClass),
		Arch:            e.Arch,
		Name:            e.ModelName,
		Tags:            e.ModelTags,
		EmbeddingLength: e.EmbeddingLen,
	}, true
}

// cacheAdmission records the admission facts on an entry so neither a later
// scan nor a picker read has to parse the same file again.
func cacheAdmission(e registry.Entry, admission gguf.Admission) registry.Entry {
	e.AdmissionClass = string(admission.Class)
	e.Arch = admission.Arch
	e.ModelName = admission.Name
	e.ModelTags = append([]string(nil), admission.Tags...)
	e.EmbeddingLen = admission.EmbeddingLength
	e.AdmissionV = registry.AdmissionCacheVersion
	return e
}

// ---- helpers ----------------------------------------------------------------

// scanDirForType walks dir recursively and returns registry.Entry values for
// every file whose extension is in exts. For .gguf files, it validates the
// architecture to prevent diffusion models (flux, sd3, etc.) from being added
// as LLMs, and skips .gguf files that are clearly LLM models when scanning
// for image assets.
func scanDirForType(dir, assetType string, exts []string, known func(string) (registry.Entry, bool)) ([]registry.Entry, map[string]struct{}, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot access directory %q: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, nil, fmt.Errorf("%q is not a directory", dir)
	}

	var entries []registry.Entry
	seen := make(map[string]struct{})

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
			sizeBytes := int64(0)
			modNS := int64(0)
			sizeMB := 0.0
			if statErr == nil {
				sizeBytes = info.Size()
				modNS = info.ModTime().UnixNano()
				sizeMB = float64(sizeBytes) / (1024 * 1024)
			}
			name := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))

			// A chat directory accepts GGUF only. The old extension list let a
			// Whisper .bin into the chat picker and gave every component a free
			// pass around architecture admission.
			if assetType == registry.KindLLM && ext != ".gguf" {
				return nil
			}

			// .gguf carries its own architecture, which is authoritative. Use it
			// to keep chat models and image models out of each other's
			// directories instead of guessing from filenames.
			if ext == ".gguf" {
				// Unchanged files are never re-parsed. This check sits before any
				// GGUF parse and requires a current cached admission, so an entry
				// classified before the cache existed is re-read exactly once.
				//
				// Ordering is the whole point. Reading one GGUF header costs
				// 1-2s for a large model while the walk and stat that found it
				// cost milliseconds, so a scan that admits first and compares
				// size/mtime afterwards re-reads the entire library on every
				// launch (measured: 61s for 27 chat models).
				if known != nil && sizeBytes > 0 {
					if prev, ok := known(path); ok && prev.FileSize == sizeBytes && prev.ModTimeNS == modNS {
						if _, cached := entryAdmission(prev); cached {
							entries = append(entries, prev)
							seen[path] = struct{}{}
							return nil
						}
					}
				}
				// A chat directory must not adopt a component asset (vocoder,
				// projector, encoder) or an architecture HOWL has not verified.
				// The decision comes from the file's own metadata.
				admission := gguf.AdmitModel(path)
				atomic.AddInt64(&scanLiveParses, 1)
				if assetType == registry.KindLLM {
					if !admission.AllowedAsChat() {
						fmt.Printf("[scan] excluding non-chat GGUF from LLM registry: %s (%s)\n", path, admission.Class)
						return nil
					}
					eligibility := chatpolicy.Evaluate(path, admission)
					if !eligibility.Eligible {
						// Inventory and chat eligibility are separate: the model stays
						// known to the registry, but the picker/loader will refuse it.
						fmt.Printf("[scan] chat-ineligible model kept in inventory: %s (task=%s basis=%s)\n",
							path, eligibility.Task, eligibility.Basis)
					}
				}
				if admission.Arch != "" {
					kind := gguf.KindOfArch(admission.Arch)
					if assetType == registry.KindLLM && kind == gguf.KindDiffusion {
						return nil // diffusion model in a chat directory
					}
					if assetType == registry.KindImage && kind == gguf.KindLLM {
						return nil // chat model in an image directory
					}
					entries = append(entries, cacheAdmission(registry.Entry{
						Name:      name,
						Path:      path,
						Type:      assetType,
						Family:    admission.Arch,
						Format:    strings.TrimPrefix(ext, "."),
						SizeMB:    sizeMB,
						FileSize:  sizeBytes,
						ModTimeNS: modNS,
					}, admission))
					seen[path] = struct{}{}
					return nil
				}
			}

			entries = append(entries, registry.Entry{
				Name:      name,
				Path:      path,
				Type:      assetType,
				Family:    assetType,
				Format:    strings.TrimPrefix(ext, "."),
				SizeMB:    sizeMB,
				FileSize:  sizeBytes,
				ModTimeNS: modNS,
			})
			seen[path] = struct{}{}
			break
		}
		return nil
	})
	return entries, seen, err
}

func llmInfos(entries []registry.Entry) []LLMModelInfo {
	out := make([]LLMModelInfo, 0, len(entries))
	excluded := 0
	liveParsed := 0
	for _, e := range entries {
		// The registry may already hold entries written before admission or
		// task eligibility existed (the LFM vocoder, Qwen3-ASR, Whisper .bin).
		// Filter at read time so the picker can never offer one, even when an
		// old registry file is loaded unchanged.
		if strings.Contains(strings.ToLower(filepath.Base(e.Path)), "mmproj") {
			excluded++
			continue
		}
		if !strings.EqualFold(filepath.Ext(e.Path), ".gguf") {
			excluded++
			continue
		}
		admission, cached := entryAdmission(e)
		if cached {
			// The cache describes the file as it was scanned. If the file was
			// replaced since then, its classification may no longer hold, so it
			// is re-read once instead of being advertised from stale facts.
			if info, statErr := os.Stat(e.Path); statErr == nil &&
				(e.FileSize != info.Size() || e.ModTimeNS != info.ModTime().UnixNano()) {
				cached = false
			}
		}
		if !cached {
			// Entries written before the admission cache existed, or whose file
			// changed since it was written, are read once; the next scan records
			// the facts and this path goes cold again.
			admission = gguf.AdmitModel(e.Path)
			liveParsed++
		}
		if !admission.AllowedAsChat() {
			excluded++
			continue
		}
		eligibility := chatpolicy.Evaluate(e.Path, admission)
		if !eligibility.Eligible {
			excluded++
			continue
		}
		out = append(out, LLMModelInfo{
			Name:   e.Name,
			Path:   e.Path,
			Format: e.Format,
			SizeMB: e.SizeMB,
		})
	}
	if liveParsed > 0 {
		fmt.Printf("[GetLLMModels] live-parsed %d entr(ies) with no cached admission\n", liveParsed)
	}
	if excluded > 0 {
		fmt.Printf("[GetLLMModels] excluded %d non-chat or task-ineligible asset(s) from the chat picker\n", excluded)
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
