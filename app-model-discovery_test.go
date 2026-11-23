package main

// app-model-discovery_test.go
//
// Detection must classify a folder from its contents, and the prompt preference
// must be a persisted setting the user can reverse — not a flag that resets.

import (
	"os"
	"path/filepath"
	"testing"

	appsettings "howl-chat/internal/settings"
)

func writeFixture(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte("fixture"), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestInspectModelDirectoryDetectsKindFromContents(t *testing.T) {
	tests := []struct {
		name string
		file string
		want string
	}{
		{"chat", "model-Q4_K_M.gguf", "llm"},
		{"voice", "voice.onnx", "tts"},
		{"image", "model.safetensors", "image"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFixture(t, dir, test.file)
			kind, files := inspectModelDirectory(dir)
			if kind != test.want || files != 1 {
				t.Fatalf("inspect = (%q,%d), want (%q,1)", kind, files, test.want)
			}
		})
	}
}

func TestInspectModelDirectoryIgnoresNonModelFolders(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, filepath.Join(dir, "notes"), "readme.txt")
	if kind, files := inspectModelDirectory(dir); kind != "" || files != 0 {
		t.Fatalf("inspect = (%q,%d), want no model folder", kind, files)
	}
}

func TestDiscoveryPreferencePersistsAndIsReversible(t *testing.T) {
	configDir := t.TempDir()
	app := &App{settings: newTestSettingsAuthority(t, configDir)}

	// "Never show again" suppresses prompts but must not disable discovery.
	if err := app.SetDirectoryDiscovery(false); err != nil {
		t.Fatalf("SetDirectoryDiscovery(false): %v", err)
	}
	if !app.discoveryEnabled() || !app.discoverySuppressed() {
		t.Fatalf("after suppression enabled=%v suppressed=%v, want discovery on and prompts off",
			app.discoveryEnabled(), app.discoverySuppressed())
	}

	stored, err := appsettings.NewStore(configDir).Load()
	if err != nil {
		t.Fatalf("reload settings: %v", err)
	}
	settings := newTestSettingsAuthority(t, configDir)
	settings.Adopt(stored)
	reloaded := &App{settings: settings}
	if !reloaded.discoverySuppressed() {
		t.Fatal("suppression did not survive a restart")
	}

	// Re-enabling in Settings reverses it and clears the suppression.
	if err := reloaded.SetDirectoryDiscovery(true); err != nil {
		t.Fatalf("SetDirectoryDiscovery(true): %v", err)
	}
	if reloaded.discoverySuppressed() || !reloaded.discoveryEnabled() {
		t.Fatalf("after re-enable enabled=%v suppressed=%v, want prompts on",
			reloaded.discoveryEnabled(), reloaded.discoverySuppressed())
	}
}
