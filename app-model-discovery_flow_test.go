package main

// app-model-discovery_flow_test.go
//
// The discovery cache is only correct if every directory change reaches the
// picker. These tests drive the real scan -> registry -> llmInfos path with
// minimal GGUF files (real parser, no tensors) and check all four cases:
// added, removed, replaced (content changed) and modified (mtime only).
//
// The prompt tests then cover detection, deferral, suppression and the silent
// scan that keeps pickers accurate after "Never show again".

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
	"time"

	"howl-chat/internal/image/registry"
	appsettings "howl-chat/internal/settings"
)

// writeDiscoveryGGUF writes a minimal GGUF v3 header so classification runs
// through the real parser instead of a test-only decision table.
func writeDiscoveryGGUF(t *testing.T, dir, name string, kvs map[string]any) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create dir: %v", err)
	}
	path := filepath.Join(dir, name)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("create %s: %v", name, err)
	}
	defer file.Close()

	writeString := func(value string) {
		var size [8]byte
		binary.LittleEndian.PutUint64(size[:], uint64(len(value)))
		if _, err := file.Write(size[:]); err != nil {
			t.Fatalf("write string size: %v", err)
		}
		if _, err := file.WriteString(value); err != nil {
			t.Fatalf("write string: %v", err)
		}
	}

	if err := binary.Write(file, binary.LittleEndian, uint32(0x46554747)); err != nil {
		t.Fatalf("write magic: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(3)); err != nil {
		t.Fatalf("write version: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint64(0)); err != nil {
		t.Fatalf("write tensor count: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint64(len(kvs))); err != nil {
		t.Fatalf("write kv count: %v", err)
	}

	for key, value := range kvs {
		writeString(key)
		switch typed := value.(type) {
		case string:
			if err := binary.Write(file, binary.LittleEndian, uint32(8)); err != nil {
				t.Fatalf("write string type: %v", err)
			}
			writeString(typed)
		case []string:
			if err := binary.Write(file, binary.LittleEndian, uint32(9)); err != nil {
				t.Fatalf("write array type: %v", err)
			}
			if err := binary.Write(file, binary.LittleEndian, uint32(8)); err != nil {
				t.Fatalf("write array element type: %v", err)
			}
			if err := binary.Write(file, binary.LittleEndian, uint64(len(typed))); err != nil {
				t.Fatalf("write array count: %v", err)
			}
			for _, item := range typed {
				writeString(item)
			}
		default:
			t.Fatalf("unsupported test metadata type %T", value)
		}
	}
	return path
}

func pickerNames(t *testing.T, reg *registry.Registry, dir string) map[string]bool {
	t.Helper()
	entries, seen, err := scanDirForType(dir, registry.KindLLM, []string{".gguf"}, reg.Lookup)
	if err != nil {
		t.Fatalf("scanDirForType: %v", err)
	}
	if _, _, err := reg.SyncDir(dir, entries, seen); err != nil {
		t.Fatalf("SyncDir: %v", err)
	}
	names := map[string]bool{}
	for _, info := range llmInfos(reg.ByTypeInDirs("llm", registry.KindLLM)) {
		names[filepath.Base(info.Path)] = true
	}
	return names
}

func TestDiscoveryFlowUpdatesPickerForEveryDirectoryChange(t *testing.T) {
	models := t.TempDir()
	reg, err := registry.New(t.TempDir())
	if err != nil {
		t.Fatalf("registry.New: %v", err)
	}
	if err := reg.RegisterDir(models, registry.KindLLM); err != nil {
		t.Fatalf("RegisterDir: %v", err)
	}

	// 1. Added file.
	alpha := writeDiscoveryGGUF(t, models, "alpha.gguf", map[string]any{"general.architecture": "llama"})
	if names := pickerNames(t, reg, models); !names["alpha.gguf"] {
		t.Fatalf("picker = %v, want the added chat model", names)
	}

	// 2. Added file that must stay out of the chat picker.
	writeDiscoveryGGUF(t, models, "asr.gguf", map[string]any{
		"general.architecture": "qwen3vl",
		"general.tags":         []string{"automatic-speech-recognition"},
	})
	names := pickerNames(t, reg, models)
	if !names["alpha.gguf"] {
		t.Fatalf("picker = %v, lost an unchanged model", names)
	}
	if names["asr.gguf"] {
		t.Fatalf("picker = %v, admitted a speech-recognition model", names)
	}

	// 3. Modified file (mtime only): still a chat model, still listed.
	future := time.Now().Add(2 * time.Hour)
	if err := os.Chtimes(alpha, future, future); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	if names := pickerNames(t, reg, models); !names["alpha.gguf"] {
		t.Fatalf("picker = %v, a modified chat model disappeared", names)
	}

	// 4. Added second chat model.
	writeDiscoveryGGUF(t, models, "beta.gguf", map[string]any{"general.architecture": "llama"})
	if names := pickerNames(t, reg, models); !names["beta.gguf"] {
		t.Fatalf("picker = %v, want both chat models", names)
	}

	// 5. Replaced file: same path, new content that is a component, not chat.
	writeDiscoveryGGUF(t, models, "beta.gguf", map[string]any{
		"general.architecture":    "clip",
		"clip.has_vision_encoder": "yes",
	})
	names = pickerNames(t, reg, models)
	if names["beta.gguf"] {
		t.Fatalf("picker = %v, kept a replaced file under its old classification", names)
	}
	if !names["alpha.gguf"] {
		t.Fatalf("picker = %v, lost the untouched model", names)
	}

	// 6. Removed file.
	if err := os.Remove(alpha); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if names := pickerNames(t, reg, models); names["alpha.gguf"] {
		t.Fatalf("picker = %v, kept a deleted model", names)
	}
}

func TestDiscoveryPromptDetectionDeferSuppressAndSilentScan(t *testing.T) {
	parent := t.TempDir()
	registered := filepath.Join(parent, "Chat")
	writeDiscoveryGGUF(t, registered, "known.gguf", map[string]any{"general.architecture": "llama"})

	reg, err := registry.New(t.TempDir())
	if err != nil {
		t.Fatalf("registry.New: %v", err)
	}
	if err := reg.RegisterDir(registered, registry.KindLLM); err != nil {
		t.Fatalf("RegisterDir: %v", err)
	}

	candidate := filepath.Join(parent, "New Models")
	writeDiscoveryGGUF(t, candidate, "fresh.gguf", map[string]any{"general.architecture": "llama"})

	recordsDir := t.TempDir()
	app := &App{
		modelRegistry:        reg,
		discoveryDirOverride: recordsDir,
		settingsStore:        appsettings.NewStore(t.TempDir()),
	}

	state, err := app.GetModelDiscoveryState()
	if err != nil {
		t.Fatalf("GetModelDiscoveryState: %v", err)
	}
	if !state.Enabled || !state.Prompting || state.Candidate == nil {
		t.Fatalf("state = %+v, want an offered candidate", state)
	}
	if state.Candidate.Path != candidate {
		t.Fatalf("candidate = %q, want %q", state.Candidate.Path, candidate)
	}

	// Not now: the same folder is not offered again immediately, even after a
	// restart (the record is on disk).
	if err := app.DeferDiscoveredDirectory(candidate); err != nil {
		t.Fatalf("DeferDiscoveredDirectory: %v", err)
	}
	restarted := &App{
		modelRegistry:        reg,
		discoveryDirOverride: recordsDir,
		settingsStore:        appsettings.NewStore(t.TempDir()),
	}
	state, _ = restarted.GetModelDiscoveryState()
	if state.Candidate != nil || state.Deferred != 1 {
		t.Fatalf("state = %+v, want the deferred folder withheld", state)
	}

	// Never show again: prompts off, discovery still on, and a newly found
	// folder is scanned silently so the picker stays accurate.
	silent := filepath.Join(parent, "Newer Models")
	writeDiscoveryGGUF(t, silent, "silent.gguf", map[string]any{"general.architecture": "llama"})
	if err := restarted.SetDirectoryDiscovery(false); err != nil {
		t.Fatalf("SetDirectoryDiscovery(false): %v", err)
	}
	state, _ = restarted.GetModelDiscoveryState()
	if state.Prompting || !state.Enabled || !state.Suppressed {
		t.Fatalf("state = %+v, want suppressed prompts with discovery on", state)
	}
	restarted.RunSilentModelDiscovery()
	if kind := reg.KindForPath(silent); kind != registry.KindLLM {
		t.Fatalf("silent discovery kind = %q, want the folder registered as chat", kind)
	}
	if names := pickerNames(t, reg, silent); !names["silent.gguf"] {
		t.Fatalf("picker = %v, want the silently discovered model", names)
	}

	// Reversible in Settings: prompts come back.
	if err := restarted.SetDirectoryDiscovery(true); err != nil {
		t.Fatalf("SetDirectoryDiscovery(true): %v", err)
	}
	state, _ = restarted.GetModelDiscoveryState()
	if !state.Prompting || state.Suppressed {
		t.Fatalf("state = %+v, want prompts restored", state)
	}
}
