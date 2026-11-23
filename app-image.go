package main

import (
	"fmt"
	"os"

	imageengine "howl-chat/internal/image/engine"
	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/image/registry"
	imgsettings "howl-chat/internal/image/settings"
	imagetypes "howl-chat/internal/image/types"
)

// -----------------------------------------------------------------------
// The image session is owned by internal/media: the engine, gallery, job
// manager, settings store and last model selection all live in media.Image.
// App holds the owner (`images`) and the bindings; it keeps no copy of that
// state. Generation and model loading remain App clients of the owner and are
// recorded as retained in docs/architecture/HOC-H0-LEVEL0-2026-11-22.md §3.
// -----------------------------------------------------------------------

// saveImageSettings persists the current engine config and last-model state.
func (a *App) saveImageSettings() {
	// The image domain owns both halves of this record (engine config and last
	// model), so it owns the write.
	if err := a.images.SaveSettings(); err != nil {
		fmt.Printf("[ImageEngine] Failed to save settings: %v", err)
		fmt.Println()
	}
}

// GetImageSettings returns the persisted image engine settings for the frontend.
func (a *App) GetImageSettings() imgsettings.ImageSettings {
	if a.images.Settings() == nil {
		return imgsettings.ImageSettings{}
	}
	s, _ := a.images.Settings().Load()
	return s
}

// GetImageEngineConfig returns the current image engine configuration.
func (a *App) GetImageEngineConfig() imagetypes.EngineConfig {
	return a.images.EngineConfig()
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
	if a.registries.Models() != nil {
		if err := a.registries.Models().RegisterDir(dir, registry.KindImage); err != nil {
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
	// Engine construction, teardown and persistence are the image domain's
	// own work now; App only states what changed.
	return a.images.Reconfigure(mutate)
}

// UnloadImageModel shuts down the sd-server subprocess and frees model memory.
func (a *App) UnloadImageModel() error {
	if a.images.Engine() == nil {
		return fmt.Errorf("image engine not initialized")
	}
	fmt.Printf("[App] UnloadImageModel called\n")
	return a.images.Engine().Shutdown()
}

// Ensure interface compliance is checked at compile time.
var _ imageengine.ImageEngine = (*sdcpp.Engine)(nil)
