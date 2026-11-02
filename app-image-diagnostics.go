package main

import (
	"fmt"
	"howl-chat/internal/paths"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	goruntime "runtime"
	"strings"
	"sync"

	"howl-chat/internal/image/engine/sdcpp"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// -----------------------------------------------------------------------
// Diagnostics
// -----------------------------------------------------------------------

// DebugImageEngine returns raw directory contents and engine config for diagnostics.
func (a *App) DebugImageEngine() map[string]interface{} {
	result := map[string]interface{}{"engine_nil": a.imageEngine == nil}
	if a.imageEngine == nil {
		return result
	}
	cfg := a.GetImageEngineConfig()
	result["binary_path"] = cfg.BinaryPath
	result["models_dir"] = cfg.ModelsDir
	result["port"] = cfg.Port
	if cfg.ModelsDir != "" {
		entries, err := os.ReadDir(cfg.ModelsDir)
		if err != nil {
			result["read_dir_error"] = err.Error()
		} else {
			var files []map[string]interface{}
			for _, e := range entries {
				info, _ := e.Info()
				f := map[string]interface{}{"name": e.Name(), "is_dir": e.IsDir()}
				if info != nil {
					f["size"] = info.Size()
				}
				files = append(files, f)
			}
			result["files"] = files
			result["file_count"] = len(files)
		}
	}
	return result
}

// GetAppDiagnostics returns runtime memory stats and engine health.
func (a *App) GetAppDiagnostics() map[string]interface{} {
	var m goruntime.MemStats
	goruntime.ReadMemStats(&m)
	result := map[string]interface{}{
		"go_heap_alloc_mb":  fmt.Sprintf("%.1f", float64(m.HeapAlloc)/1024/1024),
		"go_heap_sys_mb":    fmt.Sprintf("%.1f", float64(m.HeapSys)/1024/1024),
		"go_total_alloc_mb": fmt.Sprintf("%.1f", float64(m.TotalAlloc)/1024/1024),
		"go_num_gc":         m.NumGC,
		"go_goroutines":     goruntime.NumGoroutine(),
		"engine_nil":        a.imageEngine == nil,
		"model_loaded":      false,
		"model_name":        "-",
		"model_family":      "-",
	}
	if e, ok := a.imageEngine.(*sdcpp.Engine); ok {
		model := e.GetLoadedModel()
		if model != nil {
			result["model_loaded"] = true
			result["model_name"] = model.Name
			result["model_family"] = model.Family
		}
	}
	return result
}

// GetImagePlacement reports where the loaded model's weights actually live.
// On a card that cannot hold the model, sd-server keeps them in system RAM and
// streams them per step; the app used to report "loaded" either way, so a model
// running entirely from RAM looked the same as one resident in VRAM.
func (a *App) GetImagePlacement() map[string]interface{} {
	out := map[string]interface{}{"known": false}
	e, ok := a.imageEngine.(*sdcpp.Engine)
	if !ok {
		return out
	}
	p := e.GetPlacement()
	out["known"] = p.Known
	out["vramMB"] = p.VRAMMB
	out["ramMB"] = p.RAMMB
	out["components"] = p.Components
	out["backends"] = p.Backends
	out["offloaded"] = p.Offloaded()
	if model := e.GetLoadedModel(); model != nil {
		out["model"] = model.Name
		out["family"] = model.Family
	}
	return out
}

// -----------------------------------------------------------------------
// Version checking
// -----------------------------------------------------------------------

// GetSDServerVersion runs the binary with --help and parses the commit hash.
// sdVersionCache memoises the sd-server commit per binary path: the image
// page polls GetSDServerVersion every 2s, and each call used to spawn a
// console process (visible window flash + CPU churn) instead of reusing
// the answer.
var (
	sdVersionMu    sync.Mutex
	sdVersionCache = map[string]string{}
)

// GetSDServerVersion runs the binary once and parses the commit hash.
func (a *App) GetSDServerVersion() string {
	cfg := a.GetImageEngineConfig()
	binary := cfg.BinaryPath
	if binary == "" {
		binary = sdcpp.FindSDServer()
	}
	if binary == "" {
		return "not found"
	}

	sdVersionMu.Lock()
	cached, ok := sdVersionCache[binary]
	sdVersionMu.Unlock()
	if ok {
		return cached
	}

	version := "unknown"
	cmd := exec.Command(binary, "--help")
	applyHidden(cmd)
	if out, err := cmd.CombinedOutput(); err == nil {
		re := regexp.MustCompile(`commit\s+([0-9a-f]+)`)
		if m := re.FindStringSubmatch(string(out)); len(m) > 1 {
			version = m[1]
		}
	}

	if version != "unknown" {
		sdVersionMu.Lock()
		sdVersionCache[binary] = version
		sdVersionMu.Unlock()
	}
	return version
}

// isOldSDServer returns true for known commits that fail on modern GGUFs.
// SD3.5 support was added in stable-diffusion.cpp around May 2026.
// Commits before this cannot load SD3.5 GGUF models.
func isOldSDServer(commit string) bool {
	switch commit {
	case "0665a7f", "3c2d18b", "5f3e8a1":
		return true
	}
	// Any commit before May 2026 likely doesn't support SD3.5
	// This is a conservative check - users should update to latest
	return false
}

// isFamilySupportedByBackend checks if the current sd-server version supports a model family.
// SD3.5 requires a very recent stable-diffusion.cpp build (May 2026+).
func isFamilySupportedByBackend(family string, commit string) bool {
	familyLower := strings.ToLower(family)

	// SD3.5 requires latest stable-diffusion.cpp
	if familyLower == "sd3.5" {
		// For now, assume only "unknown" or very recent commits support SD3.5
		// Users should update to latest to get SD3.5 support
		if commit == "unknown" {
			return false // conservative: assume unknown means old
		}
		// If we can verify the commit is from May 2026 or later, return true
		// For now, be conservative and require user to update
		return false
	}

	// All other families (SD1.5, SD2.1, SDXL, SD3, Flux, Chroma, etc.) are supported
	return true
}

// -----------------------------------------------------------------------
// File pickers (Wails bindings)
// -----------------------------------------------------------------------

// OpenImageModelsDirPicker opens a native folder dialog for the image models directory.
func (a *App) OpenImageModelsDirPicker() (string, error) {
	modelsDir := a.GetImageEngineConfig().ModelsDir
	if modelsDir == "" {
		modelsDir = filepath.Dir(os.Args[0])
	}
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Select Image Models Directory",
		DefaultDirectory: modelsDir,
	})
	if err != nil {
		return "", fmt.Errorf("directory dialog: %w", err)
	}
	return dir, nil
}

// OpenSDServerFilePicker opens a file dialog for selecting the sd-server binary.
func (a *App) OpenSDServerFilePicker() (string, error) {
	exeDir := filepath.Dir(os.Args[0])
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Select sd-server Binary",
		DefaultDirectory: exeDir,
		Filters: []runtime.FileFilter{
			{DisplayName: "Executable (*.exe)", Pattern: "*.exe"},
			{DisplayName: "All Files", Pattern: "*"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("file dialog: %w", err)
	}
	return path, nil
}

// -----------------------------------------------------------------------
// Path helpers shared across image subsystem files
// -----------------------------------------------------------------------

// findSDServerBinary locates the sd-server binary via sdcpp discovery.
func findSDServerBinary() string {
	return sdcpp.FindSDServer()
}

// defaultModelsDir returns the default models directory relative to the executable.
func defaultModelsDir() string {
	return filepath.Join(paths.ProjectRoot(), "models")
}
