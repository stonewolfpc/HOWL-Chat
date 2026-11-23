package main

import (
	"testing"

	"howl-chat/internal/contentimport"
	"howl-chat/internal/library"
	"howl-chat/internal/registries"
)

func contentRegistriesForTest(t *testing.T) *registries.Set {
	t.Helper()
	libs, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return registries.Compose(registries.Parts{Libraries: libs})
}

// A build that has not reached the registries stage (tests, probes) must keep
// answering its bindings with the same "not initialised" refusals the old
// per-field nil checks produced, rather than panicking on a nil owner.
func TestRegistryBindingsRefuseWhenTheOwnerIsAbsent(t *testing.T) {
	app := &App{}

	if _, err := app.GetContentLibraries(""); err == nil {
		t.Error("GetContentLibraries accepted an absent registry")
	}
	if _, err := app.GetContentLibraryCharacters(); err == nil {
		t.Error("GetContentLibraryCharacters accepted an absent registry")
	}
	if _, err := app.GetContentLibraryInventory(""); err == nil {
		t.Error("GetContentLibraryInventory accepted an absent registry")
	}
	if err := app.RemoveContentLibrary("id"); err == nil {
		t.Error("RemoveContentLibrary accepted an absent registry")
	}
	if _, err := app.GetLLMModels(); err == nil {
		t.Error("GetLLMModels accepted an absent registry")
	}
	if _, err := app.GetRegisteredLLMDirs(); err == nil {
		t.Error("GetRegisteredLLMDirs accepted an absent registry")
	}
	if _, err := app.GetModelCatalog(); err == nil {
		t.Error("GetModelCatalog accepted an absent registry")
	}
	if _, err := app.GetModelIntelligence("model.gguf"); err == nil {
		t.Error("GetModelIntelligence accepted an absent cache")
	}
	if _, err := app.ListLorebooks(); err == nil {
		t.Error("ListLorebooks accepted an absent store")
	}
	if _, err := app.GetLorebookStoreStatus(); err == nil {
		t.Error("GetLorebookStoreStatus accepted an absent store")
	}
	if _, err := app.PreviewContentImport(contentimport.Request{Category: "worlds"}); err == nil {
		t.Error("PreviewContentImport accepted an absent import service")
	}
	if _, err := app.ListContentImportJournal("worlds"); err == nil {
		t.Error("ListContentImportJournal accepted an absent import service")
	}
	if _, err := app.ListContentImportJournal("characters"); err == nil {
		t.Error("ListContentImportJournal accepted an unknown category")
	}
}

// The owner's parts are independent: a binding whose part exists works even
// when the other registries never loaded.
func TestRegistryBindingsUseOnlyTheirOwnPart(t *testing.T) {
	app := &App{registries: contentRegistriesForTest(t)}
	if _, err := app.GetContentLibraries(""); err != nil {
		t.Fatalf("GetContentLibraries: %v", err)
	}
	// Models and lorebooks were not composed, so those routes still refuse.
	if _, err := app.GetLLMModels(); err == nil {
		t.Fatal("GetLLMModels found a registry that was never composed")
	}
	if _, err := app.ListLorebooks(); err == nil {
		t.Fatal("ListLorebooks found a store that was never composed")
	}
}
