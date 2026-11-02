package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Two registry entries shared the name krea2_turbo-Q4_K_M: an InvokeAI tree and
// the G: models folder. Picking the first match sent the load at the InvokeAI
// file and sd-server rejected it with "model metadata validation failed".
func TestPreferModelPathPrefersTheModelsDirectory(t *testing.T) {
	modelsDir := t.TempDir()
	local := filepath.Join(modelsDir, "krea2_turbo-Q4_K_M.gguf")
	if err := os.WriteFile(local, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	otherDir := t.TempDir()
	remote := filepath.Join(otherDir, "krea2_turbo-Q4_K_M.gguf")
	if err := os.WriteFile(remote, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := preferModelPath([]string{remote, local}, modelsDir)
	if got != filepath.Clean(local) {
		t.Errorf("preferModelPath = %q, want the copy in %q", got, modelsDir)
	}
}

// A registered path whose file is gone must never beat one that exists.
func TestPreferModelPathSkipsMissingFiles(t *testing.T) {
	modelsDir := t.TempDir()
	present := filepath.Join(t.TempDir(), "model.gguf")
	if err := os.WriteFile(present, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	missingLocal := filepath.Join(modelsDir, "model.gguf")

	got := preferModelPath([]string{missingLocal, present}, modelsDir)
	if got != filepath.Clean(present) {
		t.Errorf("preferModelPath = %q, want the file that exists", got)
	}
}

func TestPreferModelPathIsStableAndHandlesEmpty(t *testing.T) {
	if got := preferModelPath(nil, ""); got != "" {
		t.Errorf("preferModelPath(nil) = %q, want empty", got)
	}
	dir := t.TempDir()
	a := filepath.Join(dir, "a.gguf")
	b := filepath.Join(dir, "b.gguf")
	for _, p := range []string{a, b} {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	first := preferModelPath([]string{b, a}, "")
	second := preferModelPath([]string{a, b}, "")
	if first != second {
		t.Errorf("order changed the pick: %q vs %q", first, second)
	}
}
