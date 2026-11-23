package worldstore

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/library"
)

func TestListReturnsCanonicalDocuments(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	first, err := store.Save(map[string]interface{}{"id": "wld_one", "title": "One"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Save(map[string]interface{}{"id": "wld_two", "title": "Two"}); err != nil {
		t.Fatal(err)
	}

	items, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}
	if items[0].ID != "wld_one" || items[0].Document["title"] != "One" {
		t.Fatalf("first item = %+v", items[0])
	}
	if items[0].Revision != first.Revision {
		t.Fatalf("revision = %q, want %q", items[0].Revision, first.Revision)
	}
}

func TestDeleteRemovesOnlyAnAuthoritativeItem(t *testing.T) {
	store, _, _, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	if _, err := store.Save(map[string]interface{}{"id": "wld_keep", "title": "Keep"}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Save(map[string]interface{}{"id": "wld_drop", "title": "Drop"}); err != nil {
		t.Fatal(err)
	}

	res, err := store.Delete("wld_drop")
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "deleted" {
		t.Fatalf("status = %q message=%q", res.Status, res.Message)
	}
	if _, err := os.Stat(filepath.Join(root, "wld_drop")); !os.IsNotExist(err) {
		t.Fatalf("dropped folder still exists: %v", err)
	}
	if _, err := store.Load("wld_keep"); err != nil {
		t.Fatalf("kept item was removed: %v", err)
	}
}

func TestDeleteRefusesWithoutAWriteRole(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleRead)
	// A read-only root proves absence only through a complete search; either
	// not_found or permission_denied is acceptable, but never a silent removal.
	res, err := store.Delete("wld_missing")
	if err != nil {
		t.Fatal(err)
	}
	if res.Status == "deleted" {
		t.Fatal("read-only store reported a delete")
	}
}

func TestDeleteReportsNotFoundForaCompleteSearch(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Scenarios, library.RoleSave|library.RoleWrite)
	res, err := store.Delete("scn_missing")
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != "not_found" {
		t.Fatalf("status = %q", res.Status)
	}
	if errors.Is(err, ErrNotFound) {
		t.Fatal("Delete must report the status, not only an error")
	}
}
