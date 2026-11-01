package main

import (
	"fmt"
	"howl-chat/internal/paths"
	"os"
	"path/filepath"

	imageengine "howl-chat/internal/image/engine"
	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/image/gallery"
	"howl-chat/internal/image/job"
	"howl-chat/internal/image/registry"
	imgsettings "howl-chat/internal/image/settings"
	imagetypes "howl-chat/internal/image/types"
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
	startupMark("image init: begin")
	store := imgsettings.NewStore(a.getSettingsDir())
	a.imgSettingsStore = store

	if reg, err := registry.New(a.getSettingsDir()); err != nil {
		fmt.Printf("[Registry] Could not load model registry: %v\n", err)
	} else {
		a.modelRegistry = reg
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
		cfg.ModelsDir = filepath.Join(paths.ProjectRoot(), "models")
	}
	if cfg.OutputDir == "" {
		cfg.OutputDir = filepath.Join(paths.ProjectRoot(), "images", "gallery")
	}
	if cfg.TempDir == "" {
		cfg.TempDir = filepath.Join(os.TempDir(), "howl-image-temp")
	}
	if cfg.Port == 0 {
		cfg.Port = 8081
	}

	// The configured models directory is an image directory. Recording its kind
	// is what keeps its contents out of the chat model menu.
	if cfg.ModelsDir != "" && a.modelRegistry != nil {
		if err := a.modelRegistry.RegisterDir(cfg.ModelsDir, registry.KindImage); err != nil {
			fmt.Printf("[Registry] Could not record image models dir: %v\n", err)
		}
	}

	startupMark("image init: config ready")
	a.gallerySvc = gallery.NewService(cfg.OutputDir)
	imgEngine := sdcpp.NewEngine()
	if err := imgEngine.Initialize(cfg); err != nil {
		fmt.Printf("[ImageEngine] Not ready at startup (paths not configured yet): %v\n", err)
		a.jobManager = job.NewManager(nil)
		return
	}
	a.imageEngine = imgEngine
	a.jobManager = job.NewManager(imgEngine)

	startupMark("image init: engine initialized")
	ver := a.GetSDServerVersion()
	fmt.Printf("[ImageEngine] sd-server commit: %s\n", ver)
	if isOldSDServer(ver) {
		fmt.Printf("[ImageEngine] WARNING: sd-server commit %s is outdated. Flux/SD3/SDXL GGUF models may fail. Please update.\n", ver)
	}

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
	// Remember it as an image directory so its models stay out of the chat menu.
	if a.modelRegistry != nil {
		if err := a.modelRegistry.RegisterDir(dir, registry.KindImage); err != nil {
			fmt.Printf("[Registry] Could not record image models dir: %v\n", err)
		}
	}
	return a.reinitEngine(func(cfg *imagetypes.EngineConfig) { cfg.ModelsDir = dir })
}

// SetSDServerPath updates the sd-server binary path, re-initializes the engine, and persists.
func (a *App) SetSDServerPath(path string) error {
	if path == "" {
		return fmt.Errorf("sd-server path cannot be empty")
	}
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("sd-server binary not found: %w", err)
	}
	return a.reinitEngine(func(cfg *imagetypes.EngineConfig) { cfg.BinaryPath = path })
}

// reinitEngine tears down the current engine, applies a config mutation, and starts fresh.
func (a *App) reinitEngine(mutate func(*imagetypes.EngineConfig)) error {
	cfg := a.GetImageEngineConfig()
	mutate(&cfg)
	imgEngine := sdcpp.NewEngine()
	if err := imgEngine.Initialize(cfg); err != nil {
		return fmt.Errorf("failed to re-initialize image engine: %w", err)
	}
	if a.imageEngine != nil {
		_ = a.imageEngine.Shutdown()
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

// Ensure interface compliance is checked at compile time.
var _ imageengine.ImageEngine = (*sdcpp.Engine)(nil)
