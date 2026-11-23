package main

import (
	"testing"

	"howl-chat/internal/library"
	"howl-chat/internal/registries"
)

func worldTestApp(t *testing.T, category library.Category) *App {
	t.Helper()
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
