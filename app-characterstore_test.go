package main

import (
	"testing"

	"howl-chat/internal/characterstore"
	"howl-chat/internal/paths"
)

// TestCharacterBindingsAreBackedByTheStore proves the host path the renderer
// uses: save returns the canonical document and revision, list reads it back,
// a stale revision is refused, and delete removes it.
func TestCharacterBindingsAreBackedByTheStore(t *testing.T) {
	app := &App{characters: characterstore.New(t.TempDir())}

	saved, err := app.SaveCharacter(map[string]interface{}{
		"name":        "Mira Vale",
		"description": "first",
	})
	if err != nil {
		t.Fatalf("SaveCharacter: %v", err)
	}
	id, _ := saved["id"].(string)
	if id == "" {
		t.Fatal("SaveCharacter returned no identity")
	}
	revision, _ := saved["revision"].(string)
	if revision == "" {
		t.Fatal("SaveCharacter returned no revision")
	}

	list, err := app.GetCharacters()
	if err != nil {
		t.Fatalf("GetCharacters: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("characters = %d, want 1", len(list))
	}
	if list[0]["id"] != id || list[0]["name"] != "Mira Vale" {
		t.Fatalf("listed character = %+v", list[0])
	}
	if list[0]["revision"] != revision {
		t.Fatalf("listed revision = %v, want %s", list[0]["revision"], revision)
	}

	updated, err := app.SaveCharacter(map[string]interface{}{
		"id":          id,
		"name":        "Mira Vale",
		"description": "second",
		"revision":    revision,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated["description"] != "second" {
		t.Fatalf("updated description = %v", updated["description"])
	}
	if _, err := app.SaveCharacter(map[string]interface{}{
		"id":          id,
		"name":        "Mira Vale",
		"description": "stale",
		"revision":    revision,
	}); err == nil {
		t.Fatal("a stale revision was accepted")
	}

	if err := app.DeleteCharacter(id); err != nil {
		t.Fatalf("DeleteCharacter: %v", err)
	}
	list, err = app.GetCharacters()
	if err != nil {
		t.Fatalf("GetCharacters after delete: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("characters after delete = %+v", list)
	}
}

// TestCharacterImportAndRollbackBindings proves the explicit import/rollback
// path the renderer uses for pre-cutover data.
func TestCharacterImportAndRollbackBindings(t *testing.T) {
	t.Setenv(paths.IsolationEnvVar, t.TempDir())
	app := &App{characters: characterstore.New(t.TempDir())}
	imported, err := app.ImportCharacters([]interface{}{
		map[string]interface{}{"id": "chr_imported", "name": "Imported"},
	})
	if err != nil {
		t.Fatalf("ImportCharacters: %v", err)
	}
	if len(imported) != 1 || imported[0]["status"] != "created" {
		t.Fatalf("import results = %+v", imported)
	}
	if status := app.GetCharacterCutoverStatus(); !status.Cutover || !status.RollbackAvailable {
		t.Fatalf("cutover status after import = %+v", status)
	}
	if err := app.RollbackCharacters(); err != nil {
		t.Fatalf("RollbackCharacters: %v", err)
	}
	if status := app.GetCharacterCutoverStatus(); status.Cutover || status.RollbackAvailable {
		t.Fatalf("cutover status after rollback = %+v", status)
	}
	list, err := app.GetCharacters()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("rollback left characters: %+v", list)
	}
}
