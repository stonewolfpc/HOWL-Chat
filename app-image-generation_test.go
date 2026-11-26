package main

import (
	"testing"

	"howl-chat/internal/media"
	"howl-chat/internal/paths"
)

func generationNumber(value interface{}) int {
	switch number := value.(type) {
	case int:
		return number
	case float64:
		return int(number)
	default:
		return -1
	}
}

func TestImageGenerationBindingsRoundTripAndConflict(t *testing.T) {
	t.Setenv(paths.IsolationEnvVar, t.TempDir())
	app := &App{images: media.NewImage(t.TempDir())}

	created, err := app.SaveImageGenerationSettings(map[string]interface{}{"steps": 24, "width": 768}, "")
	if err != nil {
		t.Fatal(err)
	}
	if created["revision"] == "" {
		t.Fatal("save returned no revision")
	}

	loaded, err := app.GetImageGenerationSettings()
	if err != nil {
		t.Fatal(err)
	}
	if loaded["steps"] != float64(24) || loaded["revision"] != created["revision"] {
		t.Fatalf("loaded = %#v", loaded)
	}

	if _, err := app.SaveImageGenerationSettings(map[string]interface{}{"steps": 4}, "stale"); err == nil {
		t.Fatal("stale save was accepted")
	}
	unchanged, err := app.GetImageGenerationSettings()
	if err != nil {
		t.Fatal(err)
	}
	if unchanged["steps"] != float64(24) {
		t.Fatalf("stale save changed settings: %#v", unchanged)
	}
}

func TestImageGenerationImportNeverOverwritesBackend(t *testing.T) {
	t.Setenv(paths.IsolationEnvVar, t.TempDir())
	app := &App{images: media.NewImage(t.TempDir())}

	imported, err := app.ImportImageGenerationSettings(map[string]interface{}{"steps": 31, "prompt": "legacy"})
	if err != nil {
		t.Fatal(err)
	}
	if generationNumber(imported["steps"]) != 31 {
		t.Fatalf("import = %#v", imported)
	}

	second, err := app.ImportImageGenerationSettings(map[string]interface{}{"steps": 2, "prompt": "stale"})
	if err != nil {
		t.Fatal(err)
	}
	if generationNumber(second["steps"]) != 31 || second["prompt"] != "legacy" {
		t.Fatalf("second import overwrote backend: %#v", second)
	}
}

func TestImageGenerationImportIsJournaledAndRollbackClearsIt(t *testing.T) {
	t.Setenv(paths.IsolationEnvVar, t.TempDir())
	app := &App{images: media.NewImage(t.TempDir())}

	if _, err := app.ImportImageGenerationSettings(map[string]interface{}{"steps": 31}); err != nil {
		t.Fatal(err)
	}
	status := app.GetImageGenerationCutoverStatus()
	if !status.Cutover || !status.RollbackAvailable {
		t.Fatalf("cutover status after import = %+v", status)
	}
	result, err := app.RollbackImageGenerationImport()
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "rolled_back" {
		t.Fatalf("rollback = %+v", result)
	}
	loaded, err := app.GetImageGenerationSettings()
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 1 || loaded["revision"] != "" {
		t.Fatalf("rollback left settings: %#v", loaded)
	}
	if status := app.GetImageGenerationCutoverStatus(); status.Cutover {
		t.Fatalf("cutover marker survived rollback: %+v", status)
	}
}
