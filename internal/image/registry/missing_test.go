package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMissingModelStaysRegisteredButUnavailable(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "unplugged", "model.gguf")
	raw := []Entry{{Name: "model", Path: missing, Type: "llm", Family: "llama", Format: "gguf"}}
	data, err := json.Marshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, registryFile), data, 0644); err != nil {
		t.Fatal(err)
	}
	r, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	all := r.All()
	if len(all) != 1 {
		t.Fatalf("missing model was forgotten: %+v", all)
	}
	if all[0].Available {
		t.Fatal("missing model was reported available")
	}
	if got := r.ByType("llm"); len(got) != 0 {
		t.Fatalf("missing model was offered by type query: %+v", got)
	}
}
