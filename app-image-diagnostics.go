package main

import (
	"encoding/base64"
	"fmt"
	"howl-chat/internal/host"
	"howl-chat/internal/paths"
	"log"
	"net/http"
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

// GetImagePlacement reports the theoretical requirement, the plan Howl sent,
// and where sd-server actually put the loaded model's weights. The actual half
// is diagnostic only; it never feeds back into admission.
func (a *App) GetImagePlacement() map[string]interface{} {
	return a.buildImagePlacementReport("status")
}

// -----------------------------------------------------------------------
// Version checking
// -----------------------------------------------------------------------

// ImageEngineCapabilities answers the four capability questions in one place,
// each from its own authority:
//
//	CLI   — cli_* : parsed from the binary's own --help (does it take the flag?)
//	API   — api_* : the server's /sdcpp/v1/capabilities (does it implement it?)
//	Model — families : loading contracts this build knows
//
// Whether a capability was actually activated is answered by the generation
// that used it, not by this list.
func (a *App) ImageEngineCapabilities() map[string]interface{} {
	cfg := a.images.EngineConfig()
	binary := cfg.BinaryPath
	if binary == "" {
		binary = sdcpp.FindSDServer()
	}
	cli := sdcpp.ProbeCLI(binary)
	out := map[string]interface{}{
		"binary":         binary,
		"commit":         sdcpp.BinaryVersion(binary),
		"pinned_tag":     sdcpp.EngineTag,
		"pinned_commit":  sdcpp.EngineCommit,
		"cli_probed":     cli.Probed,
		"cli_flag_count": len(cli.Flags),
		"families":       sdcpp.SupportedFamilies(),
	}
	port := cfg.Port
	if port == 0 {
		port = 8081
	}
	api := sdcpp.ProbeAPI(fmt.Sprintf("http://127.0.0.1:%d", port))
	out["api_reached"] = api.Reached
	out["api_features"] = api.Features
	out["api_modes"] = api.Modes
	out["api_samplers"] = api.Samplers
	out["api_schedulers"] = api.Schedulers
	return out
}

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

// Version-based family gates used to live here: a commit blacklist and a rule
// that reported SD 3.5 as unsupported on every build. Neither had a caller, and
// both contradicted observed behaviour (SD 3.5 generates on the pinned engine).
// Capability is answered by sdcpp.ProbeCLI against the binary in use.

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

// OpenIPAdapterImagePicker opens the ownerless native image chooser and returns
// the selected file as a data URL. It is the packaged host's real reference-image
// path: the HTML file input routes through WebView2's owned file chooser, which
// disables the Wails window and can end in the same fatal Chromium.Focus exit as
// the old directory picker. The browser harness keeps its HTML input.
func (a *App) OpenIPAdapterImagePicker() (string, error) {
	path, err := host.ChooseFile(
		"Select IP-Adapter Reference Image",
		dialogDefaultDir(a.GetImageEngineConfig().ModelsDir),
		"Images (*.png;*.jpg;*.jpeg;*.webp;*.bmp)\x00*.png;*.jpg;*.jpeg;*.webp;*.bmp\x00All files (*.*)\x00*.*\x00",
		"png",
	)
	if err != nil {
		return "", fmt.Errorf("reference-image dialog: %w", err)
	}
	if path == "" {
		log.Printf("[ImagePicker] reference image selection cancelled")
		return "", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read reference image: %w", err)
	}
	mime := http.DetectContentType(data)
	if !strings.HasPrefix(mime, "image/") {
		return "", fmt.Errorf("reference image %q is not an image (%s)", path, mime)
	}
	log.Printf("[ImagePicker] reference image selected path=%q bytes=%d mime=%s", path, len(data), mime)
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data), nil
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
