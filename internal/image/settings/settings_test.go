package settings

import (
	"testing"

	imagetypes "howl-chat/internal/image/types"
)

func TestGenerationSaveAndRevisionConflict(t *testing.T) {
	store := NewStore(t.TempDir())

	document, revision, err := store.LoadGeneration()
	if err != nil {
		t.Fatal(err)
	}
	if len(document) != 0 || revision != "" {
		t.Fatalf("fresh generation = %#v rev=%q", document, revision)
	}

	first, err := store.SaveGeneration(map[string]interface{}{"steps": 30, "cfg": 7.5}, "")
	if err != nil {
		t.Fatal(err)
	}
	if first.Status != "saved" || first.Revision == "" {
		t.Fatalf("first save = %#v", first)
	}

	stale, err := store.SaveGeneration(map[string]interface{}{"steps": 4}, first.Revision+"-stale")
	if err != nil {
		t.Fatal(err)
	}
	if stale.Status != "revision_conflict" {
		t.Fatalf("stale status = %q", stale.Status)
	}

	loaded, current, err := store.LoadGeneration()
	if err != nil {
		t.Fatal(err)
	}
	if current != first.Revision || loaded["steps"] != float64(30) {
		t.Fatalf("stale save changed settings: %#v rev=%q", loaded, current)
	}

	next, err := store.SaveGeneration(map[string]interface{}{"steps": 12, "sampler": "euler"}, current)
	if err != nil {
		t.Fatal(err)
	}
	if next.Status != "updated" || next.Document["cfg"] != float64(7.5) {
		t.Fatalf("merge lost an unknown field: %#v", next)
	}
}

func TestEngineSettingsSavePreservesGeneration(t *testing.T) {
	store := NewStore(t.TempDir())
	saved, err := store.SaveGeneration(map[string]interface{}{"width": 1024, "steps": 25}, "")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveEngineSettings(imagetypes.EngineConfig{ModelsDir: "models"}, "model.gguf", imagetypes.ModelLoadOptions{}); err != nil {
		t.Fatal(err)
	}
	document, revision, err := store.LoadGeneration()
	if err != nil {
		t.Fatal(err)
	}
	if revision != saved.Revision || document["width"] != float64(1024) {
		t.Fatalf("engine save dropped generation: %#v rev=%q", document, revision)
	}
	full, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if full.LastModel != "model.gguf" || full.EngineConfig.ModelsDir != "models" {
		t.Fatalf("engine settings = %#v", full)
	}
}

func TestLegacyGenerationWithoutRevisionIsStillAuthoritative(t *testing.T) {
	store := NewStore(t.TempDir())
	if err := store.Save(ImageSettings{Generation: map[string]interface{}{"steps": float64(18)}}); err != nil {
		t.Fatal(err)
	}

	document, revision, err := store.LoadGeneration()
	if err != nil {
		t.Fatal(err)
	}
	if revision == "" || document["steps"] != float64(18) {
		t.Fatalf("legacy generation = %#v rev=%q", document, revision)
	}

	refused, err := store.SaveGeneration(map[string]interface{}{"steps": float64(2)}, "")
	if err != nil {
		t.Fatal(err)
	}
	if refused.Status != "revision_conflict" {
		t.Fatalf("missing base status = %q", refused.Status)
	}
}
