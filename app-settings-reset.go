package main

import (
	"fmt"

	imgsettings "howl-chat/internal/image/settings"
	imagetypes "howl-chat/internal/image/types"
)

// ResetApplicationSettings writes the caller-supplied complete defaults and
// resets the separate image settings store. It never deletes creative content,
// model files, libraries, or the application itself.
//
// The unified store is rolled back if the image store cannot be committed, so a
// failed reset is reported rather than leaving the two settings contracts
// disagreeing silently.
func (a *App) ResetApplicationSettings(defaults map[string]interface{}) error {
	if defaults == nil {
		return fmt.Errorf("reset settings: defaults are required")
	}
	if a.settingsStore == nil {
		return fmt.Errorf("reset settings: settings store not initialized")
	}
	if a.imgSettingsStore == nil {
		return fmt.Errorf("reset settings: image settings store not initialized")
	}

	previousSettings := a.appSettings
	if _, err := a.imgSettingsStore.Load(); err != nil {
		return fmt.Errorf("reset settings: read image settings: %w", err)
	}

	if err := a.SetSettings(defaults); err != nil {
		return fmt.Errorf("reset settings: write application settings: %w", err)
	}

	if a.imageEngine != nil {
		if err := a.imageEngine.Shutdown(); err != nil {
			_ = a.SetSettings(previousSettings)
			return fmt.Errorf("reset settings: stop image engine: %w", err)
		}
	}
	if err := a.imgSettingsStore.Save(imgsettings.ImageSettings{}); err != nil {
		_ = a.SetSettings(previousSettings)
		return fmt.Errorf("reset settings: write image settings: %w", err)
	}
	a.lastImageModel = ""
	a.lastImageModelOpts = imagetypes.ModelLoadOptions{}
	return nil
}

// ResetImageSettings is the explicit image-only reset. The application-wide
// reset calls the same storage path so there is one definition of "image
// settings defaults".
func (a *App) ResetImageSettings() error {
	if a.imgSettingsStore == nil {
		return fmt.Errorf("image settings store not initialized")
	}
	if a.imageEngine != nil {
		if err := a.imageEngine.Shutdown(); err != nil {
			return fmt.Errorf("stop image engine: %w", err)
		}
	}
	if err := a.imgSettingsStore.Save(imgsettings.ImageSettings{}); err != nil {
		return fmt.Errorf("write image settings: %w", err)
	}
	a.lastImageModel = ""
	a.lastImageModelOpts = imagetypes.ModelLoadOptions{}
	return nil
}
