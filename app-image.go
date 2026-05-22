package main

import (
	"encoding/base64"
	"fmt"
	"howl-chat/internal/errors"
	"howl-chat/internal/resource"
	"os"
	"path/filepath"
	goruntime "runtime"
	"time"

	imageengine "howl-chat/internal/image/engine"
	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/image/gallery"
	"howl-chat/internal/image/job"
	"howl-chat/internal/image/registry"
	imgsettings "howl-chat/internal/image/settings"
	imagetypes "howl-chat/internal/image/types"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// -----------------------------------------------------------------------
// Image engine fields are declared in App (app.go):
//   imageEngine        imageengine.ImageEngine
//   gallerySvc         *gallery.Service
//   jobManager         *job.Manager
//   lastImageModel     string
//   lastImageModelOpts imagetypes.ModelLoadOptions
//   imgSettingsStore   *imgsettings.Store
// -----------------------------------------------------------------------

// initImageEngine initialises the image subsystem from persisted settings.
// Called once from NewApp. Non-fatal: if the binary is absent the engine starts nil.
func (a *App) initImageEngine() {
	store := imgsettings.NewStore(a.getSettingsDir())
	a.imgSettingsStore = store

	// Initialise the persistent model registry (non-fatal).
	if reg, err := registry.New(a.getSettingsDir()); err != nil {
		fmt.Printf("[Registry] Could not load model registry: %v\n", err)
	} else {
		a.modelRegistry = reg
		// Re-walk every remembered directory to pick up models added since last run.
		go a.RescanAllDirs()
	}

	saved, err := store.Load()
	if err != nil {
		fmt.Printf("[ImageEngine] Could not load saved settings: %v\n", err)
	}

	cfg := saved.EngineConfig
	if cfg.BinaryPath == "" {
		cfg.BinaryPath = sdcpp.FindSDServer()
	}
	if cfg.ModelsDir == "" {
		cfg.ModelsDir = filepath.Join(filepath.Dir(os.Args[0]), "models")
	}
	if cfg.OutputDir == "" {
		cfg.OutputDir = filepath.Join(filepath.Dir(os.Args[0]), "images", "gallery")
	}
	if cfg.TempDir == "" {
		cfg.TempDir = filepath.Join(os.TempDir(), "howl-image-temp")
	}
	if cfg.Port == 0 {
		cfg.Port = 8081
	}

	a.gallerySvc = gallery.NewService(cfg.OutputDir)
	imgEngine := sdcpp.NewEngine()
	if err := imgEngine.Initialize(cfg); err != nil {
		fmt.Printf("[ImageEngine] Not ready at startup (paths not configured yet): %v\n", err)
		a.jobManager = job.NewManager(nil)
		return
	}
	a.imageEngine = imgEngine
	a.jobManager = job.NewManager(imgEngine)

	if saved.LastModel != "" {
		a.lastImageModel = saved.LastModel
		a.lastImageModelOpts = saved.LastModelOpts
		fmt.Printf("[ImageEngine] Restored last model: %q\n", saved.LastModel)
	}
}

// saveImageSettings persists the current engine config and last-model state.
func (a *App) saveImageSettings() {
	if a.imgSettingsStore == nil {
		return
	}
	cfg := a.GetImageEngineConfig()
	if err := a.imgSettingsStore.Save(imgsettings.ImageSettings{
		EngineConfig:  cfg,
		LastModel:     a.lastImageModel,
		LastModelOpts: a.lastImageModelOpts,
	}); err != nil {
		fmt.Printf("[ImageEngine] Failed to save settings: %v\n", err)
	}
}

// GetImageSettings returns the persisted image engine settings for the frontend.
func (a *App) GetImageSettings() imgsettings.ImageSettings {
	if a.imgSettingsStore == nil {
		return imgsettings.ImageSettings{}
	}
	s, _ := a.imgSettingsStore.Load()
	return s
}

// -----------------------------------------------------------------------
// Engine lifecycle
// -----------------------------------------------------------------------

// GetImageEngineConfig returns the current image engine configuration.
func (a *App) GetImageEngineConfig() imagetypes.EngineConfig {
	if e, ok := a.imageEngine.(*sdcpp.Engine); ok {
		return e.GetConfig()
	}
	return imagetypes.EngineConfig{}
}

// SetImageModelsDir updates the models directory, re-initializes the engine, and persists.
func (a *App) SetImageModelsDir(dir string) error {
	if dir == "" {
		return fmt.Errorf("models directory cannot be empty")
	}
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("models directory not found: %w", err)
	}
	oldEngine := a.imageEngine
	imgEngine := sdcpp.NewEngine()
	cfg := a.GetImageEngineConfig()
	cfg.ModelsDir = dir
	if err := imgEngine.Initialize(cfg); err != nil {
		return fmt.Errorf("failed to re-initialize image engine: %w", err)
	}
	if oldEngine != nil {
		_ = oldEngine.Shutdown()
	}
	a.imageEngine = imgEngine
	if a.jobManager != nil {
		a.jobManager.SetEngine(imgEngine)
	}
	a.saveImageSettings()
	return nil
}

// SetSDServerPath updates the sd-server binary path, re-initializes the engine, and persists.
func (a *App) SetSDServerPath(path string) error {
	if path == "" {
		return fmt.Errorf("sd-server path cannot be empty")
	}
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("sd-server binary not found: %w", err)
	}
	oldEngine := a.imageEngine
	imgEngine := sdcpp.NewEngine()
	cfg := a.GetImageEngineConfig()
	cfg.BinaryPath = path
	if err := imgEngine.Initialize(cfg); err != nil {
		return fmt.Errorf("failed to re-initialize image engine: %w", err)
	}
	if oldEngine != nil {
		_ = oldEngine.Shutdown()
	}
	a.imageEngine = imgEngine
	if a.jobManager != nil {
		a.jobManager.SetEngine(imgEngine)
	}
	a.saveImageSettings()
	return nil
}

// UnloadImageModel shuts down the sd-server subprocess and frees model memory.
func (a *App) UnloadImageModel() error {
	if a.imageEngine == nil {
		return fmt.Errorf("image engine not initialized")
	}
	fmt.Printf("[App] UnloadImageModel called\n")
	return a.imageEngine.Shutdown()
}

// -----------------------------------------------------------------------
// Model loading
// -----------------------------------------------------------------------

// ensureEngineReady auto-initializes the engine and loads the model if needed.
// If genOpts is provided, its startup flags are merged into lastImageModelOpts.
func (a *App) ensureEngineReady(modelName string, genOpts ...imagetypes.ImageOptions) error {
	if len(genOpts) > 0 {
		o := genOpts[0]
		if o.VAETiling {
			a.lastImageModelOpts.VAETiling = true
		}
		if o.TemporalTiling {
			a.lastImageModelOpts.TemporalTiling = true
		}
		if o.FlashAttention {
			a.lastImageModelOpts.FlashAttention = true
		}
		if o.OffloadToCPU {
			a.lastImageModelOpts.OffloadToCPU = true
		}
		if o.VAEPath != "" {
			a.lastImageModelOpts.VAEPath = o.VAEPath
		}
		if o.ClipLPath != "" {
			a.lastImageModelOpts.ClipLPath = o.ClipLPath
		}
		if o.ClipGPath != "" {
			a.lastImageModelOpts.ClipGPath = o.ClipGPath
		}
		if o.T5XXLPath != "" {
			a.lastImageModelOpts.T5XXLPath = o.T5XXLPath
		}
	}
	if a.imageEngine == nil {
		imgEngine := sdcpp.NewEngine()
		cfg := a.GetImageEngineConfig()
		if cfg.BinaryPath == "" {
			cfg.BinaryPath = sdcpp.FindSDServer()
		}
		if cfg.BinaryPath == "" {
			return fmt.Errorf("sd-server binary not found")
		}
		if cfg.ModelsDir == "" {
			cfg.ModelsDir = filepath.Join(filepath.Dir(os.Args[0]), "models")
		}
		if cfg.Port == 0 {
			cfg.Port = 8081
		}
		if err := imgEngine.Initialize(cfg); err != nil {
			return fmt.Errorf("failed to initialize image engine: %w", err)
		}
		a.imageEngine = imgEngine
		if a.jobManager != nil {
			a.jobManager.SetEngine(imgEngine)
		}
	}
	if modelName == "" {
		modelName = a.lastImageModel
	}
	if modelName == "" {
		return fmt.Errorf("no model selected")
	}
	// Resolve the absolute path so the engine can load from any directory.
	a.lastImageModelOpts.ModelPath = a.resolveModelPath(modelName)
	return a.imageEngine.EnsureReady(modelName, a.lastImageModelOpts)
}

// SetActiveImageModel loads a model into sd-server and persists the selection.
func (a *App) SetActiveImageModel(opts interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[PANIC] SetActiveImageModel: %v\n", r)
			err = fmt.Errorf("internal panic: %v", r)
		}
	}()
	if a.imageEngine == nil {
		return fmt.Errorf("image engine not initialized")
	}

	var loadOpts imagetypes.ModelLoadOptions
	switch v := opts.(type) {
	case string:
		loadOpts.ModelName = v
	case map[string]interface{}:
		if m, ok := v["model"].(string); ok {
			loadOpts.ModelName = m
		}
		if val, ok := v["vae"].(string); ok {
			loadOpts.VAEPath = val
		}
		if cl, ok := v["clip_l"].(string); ok {
			loadOpts.ClipLPath = cl
		}
		if cg, ok := v["clip_g"].(string); ok {
			loadOpts.ClipGPath = cg
		}
		if t5, ok := v["t5xxl"].(string); ok {
			loadOpts.T5XXLPath = t5
		}
		if diff, ok := v["diffusion_model"].(string); ok {
			loadOpts.DiffusionModel = diff
		}
		if b, ok := v["backend"].(string); ok {
			loadOpts.Backend = b
		}
		if pb, ok := v["params_backend"].(string); ok {
			loadOpts.ParamsBackend = pb
		}
		if vt, ok := v["vae_tiling"].(bool); ok {
			loadOpts.VAETiling = vt
		}
		if tt, ok := v["temporal_tiling"].(bool); ok {
			loadOpts.TemporalTiling = tt
		}
		if c, ok := v["circular"].(bool); ok {
			loadOpts.Circular = c
		}
		if fa, ok := v["flash_attention"].(bool); ok {
			loadOpts.FlashAttention = fa
		}
		if oc, ok := v["offload_to_cpu"].(bool); ok {
			loadOpts.OffloadToCPU = oc
		}
		if dm, ok := v["disable_metadata"].(bool); ok {
			loadOpts.DisableMetadata = dm
		}
	default:
		return fmt.Errorf("invalid options type: %T", opts)
	}

	a.lastImageModel = loadOpts.ModelName
	loadOpts.ModelPath = a.resolveModelPath(loadOpts.ModelName)
	a.lastImageModelOpts = loadOpts
	fmt.Printf("[App] SetActiveImageModel: name=%q path=%q\n", loadOpts.ModelName, loadOpts.ModelPath)
	if err := a.imageEngine.LoadModel(loadOpts); err != nil {
		if loadOpts.OffloadToCPU {
			a.logError(errors.IMGCPUOffloadFail, errors.High, "image.load", "app-image.go:319", err)
		}
		a.logImageLoadError(err, "app-image.go:319")
		return err
	}

	// Register with resource manager
	resource.GetManager().RegisterModel(
		resource.ResourceImage,
		loadOpts.ModelName,
		loadOpts.ModelPath,
		4000,
		2000,
		resource.PriorityImage,
		!loadOpts.OffloadToCPU,
	)
	a.saveImageSettings()
	return nil
}

// -----------------------------------------------------------------------
// Model discovery
// -----------------------------------------------------------------------

// ListImageModels returns all image models available for selection.
// It merges two sources:
//  1. The engine's scan of the configured ModelsDir (current directory)
//  2. Every diffusion model in the persistent registry (any directory)
//
// Registry entries whose files no longer exist are pruned automatically.
// Duplicates (same absolute path) are deduplicated, preferring the engine's
// fresher metadata.
func (a *App) ListImageModels() ([]imagetypes.ModelInfo, error) {
	// Start with the engine scan (may be empty if engine not initialised).
	var models []imagetypes.ModelInfo
	if a.imageEngine != nil {
		m, err := a.imageEngine.ListModels()
		if err == nil && m != nil {
			models = m
		}
	}

	// Merge in registry entries not already covered by the engine scan.
	if a.modelRegistry != nil {
		seen := make(map[string]struct{}, len(models))
		for _, m := range models {
			seen[m.Path] = struct{}{}
		}
		for _, e := range a.modelRegistry.All() {
			if e.Type != "diffusion" {
				continue // only diffusion models go in the model dropdown
			}
			if _, ok := seen[e.Path]; ok {
				continue // already covered by engine scan
			}
			compat := assetCompatForInfo(e.Family, e.Format)
			files := e.Files
			if files == nil {
				files = map[string]string{"model": e.Path}
			}
			models = append(models, imagetypes.ModelInfo{
				Name:           e.Name,
				Path:           e.Path,
				Family:         e.Family,
				Capabilities:   []string{"txt2img", "img2img", "inpaint"},
				ResolutionHint: resolutionHintFor(e.Family),
				Files:          files,
				Format:         e.Format,
				Compatibility:  compat,
			})
		}
	}

	if models == nil {
		return []imagetypes.ModelInfo{}, nil
	}
	return models, nil
}

// resolveModelPath looks up the absolute file path for a model name from the
// registry. Returns empty string if not found (engine falls back to ModelsDir scan).
func (a *App) resolveModelPath(name string) string {
	if a.modelRegistry == nil || name == "" {
		return ""
	}
	for _, e := range a.modelRegistry.All() {
		if e.Name == name {
			return e.Path
		}
	}
	return ""
}

// assetCompatForInfo mirrors the logic in sdcpp/discovery.go.
func assetCompatForInfo(family, format string) string {
	compat := "unknown"
	switch format {
	case "safetensors", "ckpt", "sft":
		compat = "verified"
	case "gguf", "bundle":
		compat = "supported"
	}
	if family == "chroma" {
		compat = "unsupported"
	}
	return compat
}

// resolutionHintFor returns a default resolution string for the given family.
func resolutionHintFor(family string) string {
	switch family {
	case "flux", "sdxl", "sd3", "sdxl-turbo", "chroma", "kolors", "auraflow":
		return "1024x1024"
	case "sd2.1", "cascade":
		return "768x768"
	default:
		return "512x512"
	}
}

// ListImageLoRAs scans the models directory for LoRA files.
func (a *App) ListImageLoRAs() ([]string, error) {
	if a.imageEngine == nil {
		return []string{}, fmt.Errorf("image engine not initialized")
	}
	loras, err := a.imageEngine.ListLoRAs()
	if err != nil {
		return []string{}, err
	}
	if loras == nil {
		return []string{}, nil
	}
	return loras, nil
}

// -----------------------------------------------------------------------
// Generation
// -----------------------------------------------------------------------

// StartImageJob queues an async generation job and returns the job ID.
func (a *App) StartImageJob(opts imagetypes.ImageOptions) (string, error) {
	if a.jobManager == nil {
		return "", fmt.Errorf("job manager not initialized")
	}
	if err := a.ensureEngineReady(opts.ModelName, opts); err != nil {
		return "", fmt.Errorf("engine ready failed: %w", err)
	}
	id := a.jobManager.Start(opts)
	return id, nil
}

// GetImageJobStatus returns the current status of a job by ID.
func (a *App) GetImageJobStatus(id string) (*imagetypes.ImageJob, error) {
	if a.jobManager == nil {
		return nil, fmt.Errorf("job manager not initialized")
	}
	j := a.jobManager.Status(id)
	if j == nil {
		return nil, fmt.Errorf("job not found: %s", id)
	}
	return j, nil
}

// CancelImageJob marks a job as cancelled.
func (a *App) CancelImageJob(id string) error {
	if a.jobManager == nil {
		return fmt.Errorf("job manager not initialized")
	}
	if !a.jobManager.Cancel(id) {
		return fmt.Errorf("job not found or already finished: %s", id)
	}
	return nil
}

// GenerateImage creates an image using the active model (lazy-loads if needed).
func (a *App) GenerateImage(opts imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	if err := a.ensureEngineReady(opts.ModelName, opts); err != nil {
		return nil, fmt.Errorf("engine ready failed: %w", err)
	}

	var result *imagetypes.ImageResult
	var genErr error

	switch opts.Mode {
	case "img2img", "upscale", "outpaint":
		result, genErr = a.imageEngine.GenerateImg2Img(opts)
	case "inpaint":
		result, genErr = a.imageEngine.GenerateInpaint(opts)
	case "video":
		result, genErr = a.imageEngine.GenerateVideo(opts)
	default:
		result, genErr = a.imageEngine.GenerateTxt2Img(opts)
	}
	if genErr != nil {
		a.logImageGenError(genErr, "app-image.go:500")
		return nil, genErr
	}
	if a.gallerySvc != nil && result != nil {
		if err := a.gallerySvc.Save(result); err != nil {
			fmt.Printf("Warning: failed to save to gallery: %v\n", err)
		}
	}
	return result, nil
}

// -----------------------------------------------------------------------
// Gallery
// -----------------------------------------------------------------------

// ListGalleryImages returns all persisted gallery entries.
func (a *App) ListGalleryImages() ([]imagetypes.ImageMetadata, error) {
	if a.gallerySvc == nil {
		return nil, fmt.Errorf("gallery service not initialized")
	}
	return a.gallerySvc.List()
}

// DeleteGalleryImage removes an image from the gallery.
func (a *App) DeleteGalleryImage(baseName string) error {
	if a.gallerySvc == nil {
		return fmt.Errorf("gallery service not initialized")
	}
	return a.gallerySvc.Delete(baseName)
}

// GetGalleryImage reads a gallery image file and returns base64 PNG data.
func (a *App) GetGalleryImage(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read image: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// getGalleryDir returns the gallery output directory.
func (a *App) getGalleryDir() string {
	if a.gallerySvc != nil {
		return a.gallerySvc.BaseDir()
	}
	return filepath.Join(filepath.Dir(os.Args[0]), "images", "gallery")
}

// -----------------------------------------------------------------------
// Image file helpers
// -----------------------------------------------------------------------

// SaveImageAs shows a save dialog and writes the base64 PNG image to the chosen path.
func (a *App) SaveImageAs(imageDataB64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(imageDataB64)
	if err != nil {
		return "", fmt.Errorf("decode image data: %w", err)
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "Save Image As",
		DefaultDirectory: a.getGalleryDir(),
		DefaultFilename:  "generated-image.png",
		Filters: []runtime.FileFilter{
			{DisplayName: "PNG Images (*.png)", Pattern: "*.png"},
			{DisplayName: "JPEG Images (*.jpg)", Pattern: "*.jpg;*.jpeg"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("save dialog: %w", err)
	}
	if path == "" {
		return "", nil
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("write image: %w", err)
	}
	return path, nil
}

// SaveTempImage decodes a base64 image and writes it to the engine temp directory.
func (a *App) SaveTempImage(imageDataB64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(imageDataB64)
	if err != nil {
		return "", fmt.Errorf("decode image data: %w", err)
	}
	var tempDir string
	if e, ok := a.imageEngine.(*sdcpp.Engine); ok {
		tempDir = e.GetConfig().TempDir
	}
	if tempDir == "" {
		tempDir = filepath.Join(os.TempDir(), "howl-image-temp")
	}
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}
	timestamp := time.Now().UTC().Format("20060102_150405")
	path := filepath.Join(tempDir, fmt.Sprintf("init_%s.png", timestamp))
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("write temp image: %w", err)
	}
	return path, nil
}

// SaveEntityImage saves a base64 image to the entity-images directory.
func (a *App) SaveEntityImage(imageDataB64, entityType, entityID string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(imageDataB64)
	if err != nil {
		return "", fmt.Errorf("decode image data: %w", err)
	}
	safeType := filepath.Clean(entityType)
	if safeType == "." || safeType == "" {
		safeType = "unknown"
	}
	entityDir := filepath.Join(filepath.Dir(os.Args[0]), "images", "entities", safeType)
	if err := os.MkdirAll(entityDir, 0755); err != nil {
		return "", fmt.Errorf("create entity dir: %w", err)
	}
	safeID := filepath.Clean(entityID)
	if safeID == "." || safeID == "" {
		safeID = "new"
	}
	timestamp := time.Now().UTC().Format("20060102_150405")
	path := filepath.Join(entityDir, fmt.Sprintf("%s_%s.png", safeID, timestamp))
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("write entity image: %w", err)
	}
	return path, nil
}

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

// Ensure unused imports are referenced (imageengine is used via interface field).
var _ imageengine.ImageEngine = (*sdcpp.Engine)(nil)
