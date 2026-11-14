package worldstore

import (
	"errors"
	"os"
	"testing"

	"howl-chat/internal/library"
)

func addUnavailableRoot(t *testing.T, reg *library.Registry) string {
	t.Helper()
	path := t.TempDir()
	registration, err := reg.Add(library.Registration{
		Category: library.Worlds,
		Path:     path,
		Role:     library.RoleRead,
		Origin:   library.OriginExternal,
		Format:   library.FormatUnknown,
	})
	if err != nil {
		t.Fatalf("register unavailable root: %v", err)
	}
	if err := os.RemoveAll(path); err != nil {
		t.Fatalf("remove root: %v", err)
	}
	return registration.ID
}

func TestLookupDistinguishesMissingFromUnavailable(t *testing.T) {
	store, reg, _, _ := newTestStore(t, library.Worlds, library.RoleRead)
	unavailableID := addUnavailableRoot(t, reg)

	search, err := store.Lookup("wld_absent")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if search.Found {
		t.Fatal("absent identity reported as found")
	}
	if search.Complete {
		t.Fatal("search with a disconnected root reported complete")
	}
	if len(search.UnavailableRoots) != 1 {
		t.Fatalf("unavailable roots = %v, want one", search.UnavailableRoots)
	}
	loaded, loadErr := store.Load("wld_absent")
	if !errors.Is(loadErr, ErrUnavailable) {
		t.Fatalf("load error = %v, want ErrUnavailable", loadErr)
	}
	if errors.Is(loadErr, ErrNotFound) {
		t.Fatal("incomplete search was reported as proven absence")
	}
	if loaded.ID != "" {
		t.Fatalf("unavailable load returned an item: %+v", loaded)
	}

	if err := reg.Remove(unavailableID); err != nil {
		t.Fatalf("remove registration: %v", err)
	}
	search, err = store.Lookup("wld_absent")
	if err != nil {
		t.Fatalf("lookup after removing root: %v", err)
	}
	if search.Found || !search.Complete {
		t.Fatalf("search after removing root = %+v, want absent and complete", search)
	}
	if _, err := store.Load("wld_absent"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("load error = %v, want ErrNotFound", err)
	}
}

func TestLookupFindsItemDespiteUnavailableRoot(t *testing.T) {
	store, reg, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("wld_present"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	addUnavailableRoot(t, reg)

	search, err := store.Lookup("wld_present")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !search.Found || search.Item == nil {
		t.Fatal("item in a readable root was not found")
	}
	if search.Item.Revision != saved.Revision {
		t.Fatalf("revision = %q, want %q", search.Item.Revision, saved.Revision)
	}
	if !search.AuthorityIndeterminate {
		t.Fatal("found copy was treated as authoritative while another root was unavailable")
	}
	if search.Item.Editable {
		t.Fatal("indeterminate item reported as editable")
	}

	item, err := store.Load("wld_present")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if item.ID != "wld_present" {
		t.Fatalf("id = %q", item.ID)
	}
}

func TestUpdateDoesNotTreatUnavailableAsMissing(t *testing.T) {
	store, reg, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	addUnavailableRoot(t, reg)

	result, err := store.Update("wld_absent", map[string]interface{}{"title": "Recreate"}, "rv1-0000000000000000")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Status != "destination_unavailable" {
		t.Fatalf("status = %q, want destination_unavailable", result.Status)
	}
	if result.Status == "not_found" {
		t.Fatal("incomplete search was treated as proven absence")
	}
}

func TestLookupWithoutRegisteredRootsIsUnavailable(t *testing.T) {
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	store, err := New(reg, library.Worlds)
	if err != nil {
		t.Fatalf("store: %v", err)
	}

	search, err := store.Lookup("wld_absent")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if search.Found || search.Complete {
		t.Fatalf("search = %+v, want not found and incomplete", search)
	}
	if _, err := store.Load("wld_absent"); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("load error = %v, want ErrUnavailable", err)
	}
}
