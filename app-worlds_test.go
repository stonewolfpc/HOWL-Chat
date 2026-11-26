package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/library"
	"howl-chat/internal/paths"
	"howl-chat/internal/registries"
)

func worldTestApp(t *testing.T, category library.Category) *App {
	t.Helper()
	t.Setenv(paths.IsolationEnvVar, t.TempDir())
	libs, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := libs.Add(library.Registration{
		Category: category,
		Path:     t.TempDir(),
		Role:     library.RoleRead | library.RoleWrite | library.RoleSave,
		Origin:   library.OriginHowl,
		Format:   library.FormatHowl,
	}); err != nil {
		t.Fatal(err)
	}
	return &App{registries: registries.Compose(registries.Parts{Libraries: libs})}
}

func TestWorldBindingsRoundTripThroughTheStore(t *testing.T) {
	app := worldTestApp(t, library.Worlds)
	created, err := app.SaveWorld(map[string]interface{}{"id": "wld_app", "title": "App World"})
	if err != nil {
		t.Fatal(err)
	}
	if created["revision"] == "" || created["revision"] == nil {
		t.Fatal("save did not return a revision")
	}
	list, err := app.GetWorlds()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0]["title"] != "App World" {
		t.Fatalf("list = %#v", list)
	}

	created["title"] = "Renamed"
	updated, err := app.SaveWorld(created)
	if err != nil {
		t.Fatal(err)
	}
	if updated["title"] != "Renamed" {
		t.Fatalf("updated = %#v", updated)
	}

	// The old revision is stale after the update and must be refused.
	if _, err := app.SaveWorld(map[string]interface{}{"id": "wld_app", "title": "Stale", "revision": created["revision"]}); err == nil {
		t.Fatal("stale revision was accepted")
	}

	if err := app.DeleteWorld("wld_app"); err != nil {
		t.Fatal(err)
	}
	list, err = app.GetWorlds()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("list after delete = %#v", list)
	}
}

func TestImportWorldsSeedsAbsentItemsAndNeverOverwritesCanon(t *testing.T) {
	app := worldTestApp(t, library.Worlds)
	if _, err := app.SaveWorld(map[string]interface{}{"id": "wld_canon", "title": "Canonical"}); err != nil {
		t.Fatal(err)
	}
	imported, err := app.ImportWorlds([]interface{}{
		map[string]interface{}{"id": "wld_canon", "title": "Stale"},
		map[string]interface{}{"id": "wld_legacy", "title": "Legacy"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if imported != 1 {
		t.Fatalf("imported = %d, want 1", imported)
	}
	list, err := app.GetWorlds()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("list = %#v", list)
	}
	for _, item := range list {
		if item["id"] == "wld_canon" && item["title"] != "Canonical" {
			t.Fatalf("canonical item was overwritten: %#v", item)
		}
	}
}

func TestScenarioBindingsAreIndependentOfWorlds(t *testing.T) {
	app := worldTestApp(t, library.Scenarios)
	if _, err := app.SaveScenario(map[string]interface{}{"id": "scn_app", "title": "Scene"}); err != nil {
		t.Fatal(err)
	}
	scenarios, err := app.GetScenarios()
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarios) != 1 || scenarios[0]["id"] != "scn_app" {
		t.Fatalf("scenarios = %#v", scenarios)
	}
	if err := app.DeleteScenario("scn_missing"); err == nil {
		t.Fatal("deleting an absent scenario was reported as success")
	}
}

func TestWorldImportIsJournaledAndRollsBackOnlyImportedDocuments(t *testing.T) {
	app := worldTestApp(t, library.Worlds)
	if _, err := app.SaveWorld(map[string]interface{}{"id": "wld_canon", "title": "Canonical"}); err != nil {
		t.Fatal(err)
	}
	if _, err := app.ImportWorlds([]interface{}{
		map[string]interface{}{"id": "wld_legacy", "title": "Legacy"},
	}); err != nil {
		t.Fatal(err)
	}
	status := app.GetWorldCutoverStatus()
	if !status.Cutover || !status.RollbackAvailable {
		t.Fatalf("cutover status after import = %+v", status)
	}
	result, err := app.RollbackWorlds()
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "rolled_back" || result.Removed != 1 {
		t.Fatalf("rollback = %+v", result)
	}
	list, err := app.GetWorlds()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0]["id"] != "wld_canon" {
		t.Fatalf("rollback touched canonical content: %#v", list)
	}
	if status := app.GetWorldCutoverStatus(); status.Cutover {
		t.Fatalf("cutover marker survived rollback: %+v", status)
	}
}

func TestUpdateInWritableSiblingSurvivesBlockedDefaultSaveRoot(t *testing.T) {
	t.Setenv(paths.IsolationEnvVar, t.TempDir())
	libs, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	sibling, err := libs.Add(library.Registration{
		Category: library.Worlds, Path: t.TempDir(),
		Role:   library.RoleRead | library.RoleWrite | library.RoleSave,
		Origin: library.OriginHowl, Format: library.FormatHowl,
	})
	if err != nil {
		t.Fatal(err)
	}
	app := &App{registries: registries.Compose(registries.Parts{Libraries: libs})}
	created, err := app.SaveWorld(map[string]interface{}{"title": "Sibling World"})
	if err != nil {
		t.Fatalf("create in writable sibling: %v", err)
	}

	// The default save root is now a file: creation there is impossible while
	// the sibling that owns the document stays writable.
	blockedPath := filepath.Join(t.TempDir(), "Worlds")
	if err := os.MkdirAll(blockedPath, 0o755); err != nil {
		t.Fatal(err)
	}
	blocked, err := libs.Add(library.Registration{
		Category: library.Worlds, Path: blockedPath,
		Role:   library.RoleRead | library.RoleWrite | library.RoleSave,
		Origin: library.OriginHowl, Format: library.FormatHowl,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := libs.SetDefaultSave(blocked.ID); err != nil {
		t.Fatal(err)
	}
	// Replace the directory with a file after registration so the registry is
	// valid while the destination cannot be created into.
	if err := os.RemoveAll(blockedPath); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(blockedPath, []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	app.recordProbe(blocked.ID, ProvisionReport{
		Category: string(library.Worlds), Path: blockedPath,
		Status: ContentUnavailable, Step: "create library directory",
		Operation: "create library directory",
		Reason:    "The library path is a file, not a directory.", Retryable: false,
		Err: fs.ErrInvalid,
	})
	app.recordProbe(sibling.ID, readyProbe(library.Worlds, sibling.Path))

	created["title"] = "Sibling World Edited"
	updated, err := app.SaveWorld(created)
	if err != nil {
		t.Fatalf("update in the writable sibling was refused by a blocked default root: %v", err)
	}
	if updated["title"] != "Sibling World Edited" {
		t.Fatalf("updated = %#v", updated)
	}

	if _, err := app.SaveWorld(map[string]interface{}{"title": "New World"}); err == nil {
		t.Fatal("create was accepted although the default save root is unusable")
	}
}
