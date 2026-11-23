package characterstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func readEnvelope(t *testing.T, path string) envelope {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var env envelope
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("envelope is not JSON: %v", err)
	}
	return env
}

func TestListReadsLegacyFoldersAndKeepsTheDirectoryIdentity(t *testing.T) {
	dir := t.TempDir()
	legacyDir := filepath.Join(dir, "Mira Vale")
	if err := os.MkdirAll(legacyDir, 0755); err != nil {
		t.Fatal(err)
	}
	legacy := map[string]interface{}{"name": "Mira Vale", "description": "legacy"}
	data, _ := json.Marshal(legacy)
	if err := os.WriteFile(filepath.Join(legacyDir, "character.json"), data, 0644); err != nil {
		t.Fatal(err)
	}

	items, err := New(dir).List()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("items = %d, want 1", len(items))
	}
	if items[0].ID != "Mira Vale" || !items[0].Legacy {
		t.Fatalf("legacy item = %+v", items[0])
	}
	if items[0].Document["description"] != "legacy" {
		t.Fatalf("document = %+v", items[0].Document)
	}
}

func TestSaveCreatesEnvelopeAndUpdatesLegacyFolderInPlace(t *testing.T) {
	dir := t.TempDir()
	legacyDir := filepath.Join(dir, "Mira Vale")
	if err := os.MkdirAll(legacyDir, 0755); err != nil {
		t.Fatal(err)
	}
	legacy := map[string]interface{}{"id": "mira", "name": "Mira Vale", "description": "old", "howl": map[string]interface{}{"voice": map[string]interface{}{"enabled": true}}}
	data, _ := json.Marshal(legacy)
	if err := os.WriteFile(filepath.Join(legacyDir, "character.json"), data, 0644); err != nil {
		t.Fatal(err)
	}

	store := New(dir)
	item, err := store.Get("mira")
	if err != nil {
		t.Fatal(err)
	}
	updated, err := store.Save(map[string]interface{}{
		"id":          "mira",
		"name":        "Mira Vale",
		"description": "new",
	}, item.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Revision == item.Revision {
		t.Fatal("revision did not change")
	}
	if filepath.Dir(updated.Path) != legacyDir {
		t.Fatalf("update moved the character: %s", updated.Path)
	}
	env := readEnvelope(t, updated.Path)
	if env.Format != Format || env.Version != Version || env.ID != "mira" {
		t.Fatalf("envelope = %+v", env)
	}
	if env.Document["description"] != "new" {
		t.Fatalf("description = %v", env.Document["description"])
	}
	howl, _ := env.Document["howl"].(map[string]interface{})
	if howl == nil || howl["voice"] == nil {
		t.Fatalf("unknown section was not preserved: %+v", env.Document)
	}
}

func TestSaveRejectsAStaleRevision(t *testing.T) {
	dir := t.TempDir()
	store := New(dir)
	first, err := store.Save(map[string]interface{}{"id": "chr_test", "name": "One"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Save(map[string]interface{}{"id": "chr_test", "name": "Two"}, first.Revision); err != nil {
		t.Fatal(err)
	}
	_, err = store.Save(map[string]interface{}{"id": "chr_test", "name": "Three"}, first.Revision)
	if !errors.Is(err, ErrRevisionConflict) {
		t.Fatalf("err = %v, want ErrRevisionConflict", err)
	}
}

func TestDeleteRemovesOnlyTheRequestedCharacter(t *testing.T) {
	dir := t.TempDir()
	store := New(dir)
	for _, id := range []string{"chr_one", "chr_two"} {
		if _, err := store.Save(map[string]interface{}{"id": id, "name": id}, ""); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.Delete("chr_one"); err != nil {
		t.Fatal(err)
	}
	items, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != "chr_two" {
		t.Fatalf("items = %+v", items)
	}
}

func TestImportIsIdempotentAndRollbackRestoresWhatItReplaced(t *testing.T) {
	dir := t.TempDir()
	store := New(dir)
	if _, err := store.Save(map[string]interface{}{"id": "chr_keep", "name": "Keep", "description": "original"}, ""); err != nil {
		t.Fatal(err)
	}
	imported := []map[string]interface{}{
		{"id": "chr_keep", "name": "Keep", "description": "imported"},
		{"id": "chr_new", "name": "New"},
	}
	results, err := store.Import(imported)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 || results[0].Status != "updated" || results[1].Status != "created" {
		t.Fatalf("results = %+v", results)
	}
	again, err := store.Import(imported)
	if err != nil {
		t.Fatal(err)
	}
	if again[0].Status != "updated" || again[1].Status != "updated" {
		t.Fatalf("second import = %+v", again)
	}

	if err := store.Rollback(results); err != nil {
		t.Fatal(err)
	}
	kept, err := store.Get("chr_keep")
	if err != nil {
		t.Fatal(err)
	}
	if kept.Document["description"] != "original" {
		t.Fatalf("rollback description = %v", kept.Document["description"])
	}
	if _, err := store.Get("chr_new"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("created character survived rollback: %v", err)
	}
}

func TestSaveRejectsUnsafeNewIdentities(t *testing.T) {
	store := New(t.TempDir())
	if _, err := store.Save(map[string]interface{}{"id": "../escape", "name": "bad"}, ""); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("err = %v, want ErrInvalidDocument", err)
	}
}
