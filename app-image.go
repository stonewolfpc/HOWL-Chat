package main

import (
	"fmt"
	"os"

	"howl-chat/internal/domaincutover"
	imageengine "howl-chat/internal/image/engine"
	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/image/registry"
	"howl-chat/internal/resource"
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

// GetImageGenerationSettings returns the backend-owned generation document and
// the revision the next save must cite.
func (a *App) GetImageGenerationSettings() (map[string]interface{}, error) {
	if a.images == nil || a.images.Settings() == nil {
		return nil, fmt.Errorf("image settings store not initialised")
	}
	document, revision, err := a.images.Settings().LoadGeneration()
	if err != nil {
		return nil, fmt.Errorf("GetImageGenerationSettings: %w", err)
	}
	document["revision"] = revision
	return document, nil
}

// SaveImageGenerationSettings is the only generation-settings write path. A
// stale or missing base revision is refused instead of overwriting newer values.
func (a *App) SaveImageGenerationSettings(settings map[string]interface{}, baseRevision string) (map[string]interface{}, error) {
	if a.images == nil || a.images.Settings() == nil {
		return nil, fmt.Errorf("image settings store not initialised")
	}
	result, err := a.images.Settings().SaveGeneration(settings, baseRevision)
	if err != nil {
		return nil, fmt.Errorf("SaveImageGenerationSettings: %w", err)
	}
	if result.Status != "saved" && result.Status != "updated" {
		message := result.Message
		if message == "" {
			message = result.Status
		}
		return nil, fmt.Errorf("SaveImageGenerationSettings: %s", message)
	}
	document := result.Document
	document["revision"] = result.Revision
	return document, nil
}

// ImportImageGenerationSettings seeds the backend from pre-cutover renderer
// settings only when no generation document exists. It never replaces a saved
// backend document on a later load.
func (a *App) ImportImageGenerationSettings(settings map[string]interface{}) (map[string]interface{}, error) {
	if a.images == nil || a.images.Settings() == nil {
		return nil, fmt.Errorf("image settings store not initialised")
	}
	current, revision, err := a.images.Settings().LoadGeneration()
	if err != nil {
		return nil, fmt.Errorf("ImportImageGenerationSettings: %w", err)
	}
	if revision != "" || len(current) > 0 {
		item := domaincutover.Item{ID: "generation", Action: "adopted", Revision: revision, Previous: current}
		if err := a.cutoverStore("image-generation").Record(0, []domaincutover.Item{item}, "completed", "Existing image generation settings were adopted explicitly."); err != nil {
			return nil, fmt.Errorf("ImportImageGenerationSettings: %w", err)
		}
		current["revision"] = revision
		return current, nil
	}
	saved, err := a.SaveImageGenerationSettings(settings, "")
	if err != nil {
		return nil, err
	}
	savedRevision, _ := saved["revision"].(string)
	item := domaincutover.Item{ID: "generation", Action: "created", Revision: savedRevision}
	if err := a.cutoverStore("image-generation").Record(1, []domaincutover.Item{item}, "completed", "Image generation settings were imported explicitly."); err != nil {
		return nil, fmt.Errorf("ImportImageGenerationSettings: %w", err)
	}
	return saved, nil
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
//
// It is a lifecycle transition: it serializes against image starts, reloads and
// other unloads, and it refuses while an image workflow still owns the engine.
// Freeing the card for chat must not tear down a stage that is mid-flight.
func (a *App) UnloadImageModel() error {
	transition, err := a.beginLifecycleTransition(resource.ResourceImage, "stop", lifecycleStopWait)
	if err != nil {
		return err
	}
	defer transition.Release()
	if a.workflowReclaimBlocked() {
		return fmt.Errorf("the image model cannot be unloaded while an image workflow is running; cancel the workflow first")
	}
	if a.images.Engine() == nil {
		return fmt.Errorf("image engine not initialized")
	}
	fmt.Printf("[App] UnloadImageModel called\n")
	err = a.images.Engine().Shutdown()
	if err != nil {
		return fmt.Errorf("unload image model: %w", err)
	}
	a.unregisterImageModels()
	a.emitImagePlacement("unloaded")
	return nil
}

// Ensure interface compliance is checked at compile time.
var _ imageengine.ImageEngine = (*sdcpp.Engine)(nil)
