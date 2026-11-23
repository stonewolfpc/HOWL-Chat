package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestSetSettingsPersistsBeforeReturning is the host half of the settings
// authority contract: when the binding returns, settings.json already holds
// the document the renderer sent.
func TestSetSettingsPersistsBeforeReturning(t *testing.T) {
	dir := t.TempDir()
	authority := newTestSettingsAuthority(t, dir)
	app := &App{settings: authority}

	doc := map[string]interface{}{
		"general": map[string]interface{}{"theme": "light", "uiScale": 1.25},
		"sampler": map[string]interface{}{"temperature": 0.6},
	}
	if err := app.SetSettings(doc); err != nil {
		t.Fatalf("SetSettings: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "settings.json"))
	if err != nil {
		t.Fatalf("settings.json was not written before SetSettings returned: %v", err)
	}
	var onDisk map[string]interface{}
	if err := json.Unmarshal(data, &onDisk); err != nil {
		t.Fatalf("settings.json is not valid JSON: %v", err)
	}
	general, _ := onDisk["general"].(map[string]interface{})
	if general["theme"] != "light" || general["uiScale"] != 1.25 {
		t.Fatalf("persisted general section = %+v", general)
	}

	reloaded := newTestSettingsAuthority(t, dir)
	if err := reloaded.Load(); err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := reloaded.Section("general")["theme"]; got != "light" {
		t.Fatalf("reloaded theme = %v, want light", got)
	}
	if got := reloaded.Sampler()["temperature"]; got != 0.6 {
		t.Fatalf("reloaded sampler temperature = %v, want 0.6", got)
	}
}

// TestGetSettingsIncludesLiveSamplerAndModelSections makes GetSettings a
// complete snapshot: the renderer must not need the legacy GetSamplerSettings /
// GetModelOptions bindings to see what the runtime is actually using.
func TestGetSettingsIncludesLiveSamplerAndModelSections(t *testing.T) {
	dir := t.TempDir()
	authority := newTestSettingsAuthority(t, dir)
	authority.SetSampler(map[string]interface{}{"temperature": 0.42})
	authority.SetModel(map[string]interface{}{"think_mode": "visible"})
	app := &App{settings: authority}

	snapshot := app.GetSettings()
	sampler, ok := snapshot["sampler"].(map[string]interface{})
	if !ok || sampler["temperature"] != 0.42 {
		t.Fatalf("snapshot sampler = %+v", snapshot["sampler"])
	}
	model, ok := snapshot["model"].(map[string]interface{})
	if !ok || model["think_mode"] != "visible" {
		t.Fatalf("snapshot model = %+v", snapshot["model"])
	}
}
