package registry

import (
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/image/scanner"
)

// A file the scanner reclassifies must have its registry entry rewritten. A
// GGUF once treated as a checkpoint and later identified as a language model
// kept type "diffusion" forever, so llama-q2_k.gguf stayed in the image model
// dropdown: the directory sync deliberately keeps entries whose file still
// exists, and nothing ever corrected the type.
func TestMergeRewritesTypeWhenFileIsReclassified(t *testing.T) {
	dir := t.TempDir()
	model := filepath.Join(dir, "llama-q2_k.gguf")
	if err := os.WriteFile(model, []byte("weights"), 0o644); err != nil {
		t.Fatal(err)
	}

	r := &Registry{configDir: dir, entries: map[string]Entry{}, dirs: map[string]string{}}
	first := scanner.ModelIndex{Models: []scanner.AssetEntry{{
		Name: "llama-q2_k", Path: model, Type: scanner.AssetDiffusion, Family: scanner.FamilySD15,
	}}}
	if err := r.Merge(first); err != nil {
		t.Fatal(err)
	}
	if e, _ := r.Lookup(model); e.Type != string(scanner.AssetDiffusion) {
		t.Fatalf("setup: type = %q, want diffusion", e.Type)
	}

	second := scanner.ModelIndex{LLMs: []scanner.AssetEntry{{
		Name: "llama-q2_k", Path: model, Type: scanner.AssetLLM,
	}}}
	if err := r.Merge(second); err != nil {
		t.Fatal(err)
	}

	e, ok := r.Lookup(model)
	if !ok {
		t.Fatal("entry vanished instead of being reclassified")
	}
	if e.Type != string(scanner.AssetLLM) {
		t.Fatalf("type = %q, want llm", e.Type)
	}
	if e.Family != "" {
		t.Fatalf("family = %q, want empty for a language model", e.Family)
	}
}
