package main

import (
	"fmt"
	"howl-chat/internal/host"
	"howl-chat/internal/paths"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"

	"howl-chat/internal/image/engine/sdcpp"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// -----------------------------------------------------------------------
// Diagnostics
// -----------------------------------------------------------------------

// DebugImageEngine returns raw directory contents and engine config for diagnostics.
func (a *App) DebugImageEngine() map[string]interface{} {
	result := map[string]interface{}{"engine_nil": a.images.Engine() == nil}
	if a.images.Engine() == nil {
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
		"engine_nil":        a.images.Engine() == nil,
		"model_loaded":      false,
		"model_name":        "-",
		"model_family":      "-",
	}
	if e, ok := a.images.Engine().(*sdcpp.Engine); ok {
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
	e, ok := a.images.Engine().(*sdcpp.Engine)
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

// GetSDServerVersion reports the sd-server commit the engine package detects.
// Probing lives in sdcpp.BinaryVersion so the image domain and this binding
// share one cache and one implementation instead of spawning the binary twice.
func (a *App) GetSDServerVersion() string {
	cfg := a.images.EngineConfig()
	binary := cfg.BinaryPath
	if binary == "" {
		binary = sdcpp.FindSDServer()
	}
	if binary == "" {
		return "not found"
	}
	return sdcpp.BinaryVersion(binary)
}

// isOldSDServer reports known commits that fail on modern GGUFs.
func isOldSDServer(commit string) bool { return sdcpp.IsKnownOldVersion(commit) }

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

// dialogDefaultDir returns preferred when it is an existing directory, otherwise
// the running executable's directory. Wails refuses to open a dialog whose
// DefaultDirectory does not exist, and a stored models dir is exactly what is
// missing on a first run or after a drive is unplugged — the picker must still
// open so the user can point at the models.
func dialogDefaultDir(preferred string) string {
	if preferred != "" {
		if info, err := os.Stat(preferred); err == nil && info.IsDir() {
			return preferred
		}
	}
	if exe, err := os.Executable(); err == nil {
		return filepath.Dir(exe)
	}
	return filepath.Dir(os.Args[0])
}

// OpenImageModelsDirPicker opens the shell folder chooser for the image models
// directory.
//
// It deliberately does not use runtime.OpenDirectoryDialog: that runs the dialog
// with the main window as parent, Windows disables the window while it is open,
// and the focus event on the way back hits go-webview2's Chromium.Focus, whose
// error callback calls os.Exit(1). Add Folder therefore closed the whole app
// (observed twice on 2026-11-22). host.ChooseFolder opens a chooser with no
// owner window, so the app window is never disabled.
func (a *App) OpenImageModelsDirPicker() (string, error) {
	dir, err := host.ChooseFolder("Select Image Models Directory", dialogDefaultDir(a.GetImageEngineConfig().ModelsDir))
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
