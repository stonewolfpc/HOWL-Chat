package main

import (
	"os"
	"path/filepath"
	"testing"

	imgsettings "howl-chat/internal/image/settings"
	appsettings "howl-chat/internal/settings"
)

func TestResetApplicationSettingsResetsUnifiedAndImageStores(t *testing.T) {
	dir := t.TempDir()
	imageStore := imgsettings.NewStore(dir)
	if err := imageStore.Save(imgsettings.ImageSettings{LastModel: "old-model"}); err != nil {
		t.Fatal(err)
	}
	app := &App{
		settingsStore:    appsettings.NewStore(dir),
		imgSettingsStore: imageStore,
		appSettings:      map[string]interface{}{"general": map[string]interface{}{"theme": "light"}},
		lastImageModel:   "old-model",
	}
	defaults := map[string]interface{}{
		"general": map[string]interface{}{"theme": "dark"},
	}
	if err := app.ResetApplicationSettings(defaults); err != nil {
		t.Fatal(err)
	}
	loaded, err := appsettings.NewStore(dir).Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded["general"].(map[string]interface{})["theme"] != "dark" {
		t.Fatalf("unified settings = %+v", loaded)
	}
	image, err := imageStore.Load()
	if err != nil {
		t.Fatal(err)
	}
	if image.LastModel != "" || image.EngineConfig.ModelsDir != "" {
		t.Fatalf("image settings = %+v, want zero", image)
	}
	if app.lastImageModel != "" {
		t.Fatalf("lastImageModel = %q, want empty", app.lastImageModel)
	}
}

func TestResetApplicationSettingsDoesNotTouchCreativeFiles(t *testing.T) {
	dir := t.TempDir()
	imageStore := imgsettings.NewStore(dir)
	app := &App{settingsStore: appsettings.NewStore(dir), imgSettingsStore: imageStore, appSettings: map[string]interface{}{}}
	content := filepath.Join(dir, "character.json")
	if err := os.WriteFile(content, []byte(`{"name":"keep"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := app.ResetApplicationSettings(map[string]interface{}{"general": map[string]interface{}{"theme": "dark"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(content); err != nil {
		t.Fatalf("creative file was removed: %v", err)
	}
}
