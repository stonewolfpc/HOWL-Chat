package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/library"
)

func TestExportContentItemUsesStableIdentityAndDoesNotOverwrite(t *testing.T) {
	root := t.TempDir()
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := reg.Add(library.Registration{
		Category: library.Characters,
		Path:     root,
		Role:     library.RoleSave,
		Origin:   library.OriginHowl,
		Format:   library.FormatHowl,
	}); err != nil {
		t.Fatal(err)
	}
	app := &App{contentLibraries: reg}
	first, err := app.ExportContentItem("characters", "Example World", "world-1", map[string]interface{}{"title": "Example World"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if first.Status != library.StatusExportSucceeded || first.Target == "" {
		t.Fatalf("first export result = %+v", first)
	}
	if _, err := os.Stat(first.Target); err != nil {
		t.Fatalf("export file missing: %v", err)
	}
	second, err := app.ExportContentItem("characters", "Example World", "world-1", map[string]interface{}{"title": "Changed"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if second.Status != library.StatusNameCollision {
		t.Fatalf("second export status = %q, want collision", second.Status)
	}
	renamed, err := app.ExportContentItem("characters", "Renamed World", "world-1", map[string]interface{}{"title": "Renamed"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if renamed.Status != library.StatusNameCollision {
		t.Fatalf("renamed identity status = %q, want collision", renamed.Status)
	}
	data, err := os.ReadFile(first.Target)
	if err != nil {
		t.Fatal(err)
	}
	var envelope map[string]interface{}
	if err := json.Unmarshal(data, &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope["format"] != "howl-content-export-v1" || envelope["authoritative"] != false {
		t.Fatalf("export envelope = %+v", envelope)
	}
	content, ok := envelope["content"].(map[string]interface{})
	if !ok || content["title"] != "Example World" {
		t.Fatalf("original export was overwritten: %s", data)
	}
}

func TestExportContentItemRejectsMissingIdentityAndData(t *testing.T) {
	root := t.TempDir()
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := reg.Add(library.Registration{
		Category: library.Characters,
		Path:     root,
		Role:     library.RoleSave,
		Origin:   library.OriginHowl,
		Format:   library.FormatHowl,
	}); err != nil {
		t.Fatal(err)
	}
	app := &App{contentLibraries: reg}
	if _, err := app.ExportContentItem("characters", "", "", nil, false); err == nil {
		t.Fatal("export without identity or data was accepted")
	}
}

func TestExportContentItemRequiresWritePermissionForReplacement(t *testing.T) {
	root := t.TempDir()
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	added, err := reg.Add(library.Registration{
		Category: library.Characters,
		Path:     root,
		Role:     library.RoleSave,
		Origin:   library.OriginHowl,
		Format:   library.FormatHowl,
	})
	if err != nil {
		t.Fatal(err)
	}
	app := &App{contentLibraries: reg}
	result, err := app.ExportContentItem("characters", "Scene", "scene-1", map[string]interface{}{"title": "Scene"}, true)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != library.StatusPermissionDenied {
		t.Fatalf("save-only replacement status = %q, want permission_denied", result.Status)
	}
	if err := reg.SetRole(added.ID, library.RoleRead|library.RoleWrite|library.RoleSave); err != nil {
		t.Fatal(err)
	}
	result, err = app.ExportContentItem("characters", "Scene", "scene-1", map[string]interface{}{"title": "Scene two"}, true)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != library.StatusReplaced {
		t.Fatalf("writable replacement status = %q, want replaced", result.Status)
	}
}

func TestExportContentItemRefusesCanonicalWorldAndScenarioRoots(t *testing.T) {
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	app := &App{contentLibraries: reg}
	for _, category := range []string{"worlds", "scenarios"} {
		result, err := app.ExportContentItem(category, "Canonical", "item-1", map[string]interface{}{"title": "Canonical"}, false)
		if err == nil {
			t.Fatalf("%s export bypassed the canonical store: %+v", category, result)
		}
		if result.Target != "" {
			t.Fatalf("%s export produced target %q", category, result.Target)
		}
	}
}

func TestExportContentItemWithoutDestinationReportsUnavailable(t *testing.T) {
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	app := &App{contentLibraries: reg}
	result, err := app.ExportContentItem("characters", "Nobody", "nobody", map[string]interface{}{"name": "Nobody"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != library.StatusDestinationUnavailable || result.Target != "" {
		t.Fatalf("result = %+v", result)
	}
}

func TestSafeContentSegmentNeutralizesWindowsNamesAndTraversal(t *testing.T) {
	for input, want := range map[string]string{
		"CON":        "Item_CON",
		"../escape":  "_escape",
		"a:b*c?.txt": "a_b_c_.txt",
		"":           "Untitled",
	} {
		if got := safeContentSegment(input); got != want {
			t.Errorf("safeContentSegment(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestUnavailableRootKeepsLastKnownCharacters(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "card.json"), []byte(`{"name":"Persisted Card"}`), 0644); err != nil {
		t.Fatal(err)
	}
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := reg.Add(library.Registration{
		Category: library.Characters,
		Path:     root,
		Role:     library.RoleRead,
		Origin:   library.OriginExternal,
		Format:   library.FormatSillyTavern,
	}); err != nil {
		t.Fatal(err)
	}
	app := &App{contentLibraries: reg, contentInventory: library.NewInventoryStore(t.TempDir())}
	inventory, err := app.GetContentLibraryInventory("characters")
	if err != nil {
		t.Fatal(err)
	}
	if len(inventory) != 1 || !inventory[0].Scanned || inventory[0].Stale || len(inventory[0].Items) != 1 {
		t.Fatalf("live inventory = %+v", inventory)
	}
	if err := os.RemoveAll(root); err != nil {
		t.Fatal(err)
	}
	inventory, err = app.GetContentLibraryInventory("characters")
	if err != nil {
		t.Fatal(err)
	}
	if len(inventory) != 1 || inventory[0].Available || inventory[0].Scanned || !inventory[0].Stale || len(inventory[0].Items) != 1 {
		t.Fatalf("unavailable inventory = %+v", inventory)
	}
	if inventory[0].Items[0].Available {
		t.Fatal("stale item was reported available")
	}
	characters, err := app.GetContentLibraryCharacters()
	if err != nil {
		t.Fatal(err)
	}
	if len(characters) != 1 || characters[0].Available {
		t.Fatalf("unavailable characters = %+v", characters)
	}
}

func TestRemoveContentLibraryDoesNotDeleteFiles(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "keep.json")
	if err := os.WriteFile(file, []byte(`{"name":"Keep"}`), 0644); err != nil {
		t.Fatal(err)
	}
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	added, err := reg.Add(library.Registration{
		Category: library.Lorebooks,
		Path:     root,
		Role:     library.RoleRead,
		Origin:   library.OriginExternal,
		Format:   library.FormatHowl,
	})
	if err != nil {
		t.Fatal(err)
	}
	app := &App{contentLibraries: reg, contentInventory: library.NewInventoryStore(t.TempDir())}
	if err := app.RemoveContentLibrary(added.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(file); err != nil {
		t.Fatalf("external file was removed: %v", err)
	}
}

func TestExportStatusHelpers(t *testing.T) {
	if !errors.Is(library.ErrCollision, library.ErrCollision) {
		t.Fatal("collision sentinel is broken")
	}
	if (library.OperationResult{Status: library.StatusPartial}).OK() {
		t.Fatal("partial operation reported OK")
	}
}
