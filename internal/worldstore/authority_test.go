package worldstore

import (
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/library"
)

func addReadableRoot(t *testing.T, reg *library.Registry) library.Registration {
	t.Helper()
	registration, err := reg.Add(library.Registration{
		Category: library.Worlds,
		Path:     t.TempDir(),
		Role:     library.RoleRead,
		Origin:   library.OriginExternal,
		Format:   library.FormatUnknown,
	})
	if err != nil {
		t.Fatalf("register readable root: %v", err)
	}
	return registration
}

func TestUniqueCopyAcrossReadableRootsIsAuthoritative(t *testing.T) {
	store, reg, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("wld_authoritative"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	addReadableRoot(t, reg)

	search, err := store.Lookup("wld_authoritative")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !search.Found || !search.Complete {
		t.Fatalf("search = %+v, want found and complete", search)
	}
	if search.AuthorityIndeterminate {
		t.Fatal("single copy across readable roots was marked indeterminate")
	}
	if search.Copies != 1 {
		t.Fatalf("copies = %d, want 1", search.Copies)
	}
	if search.Item == nil || !search.Item.Editable {
		t.Fatal("authoritative writable copy was not editable")
	}

	result, err := store.Update("wld_authoritative", map[string]interface{}{"title": "Still Canonical"}, saved.Revision)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Status != "updated" {
		t.Fatalf("status = %q, want updated", result.Status)
	}
}

func TestDuplicateCopiesAreIndeterminate(t *testing.T) {
	store, reg, _, rootA := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("wld_duplicate"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	rootB := addReadableRoot(t, reg)

	payload, err := os.ReadFile(saved.Path)
	if err != nil {
		t.Fatalf("read copy: %v", err)
	}
	secondPath := filepath.Join(rootB.Path, "wld_duplicate", "world.json")
	if err := os.MkdirAll(filepath.Dir(secondPath), 0755); err != nil {
		t.Fatalf("create second folder: %v", err)
	}
	if err := os.WriteFile(secondPath, payload, 0644); err != nil {
		t.Fatalf("write second copy: %v", err)
	}

	search, err := store.Lookup("wld_duplicate")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !search.Found || search.Copies != 2 {
		t.Fatalf("search = %+v, want found with two copies", search)
	}
	if !search.AuthorityIndeterminate {
		t.Fatal("duplicate copies were treated as authoritative")
	}
	if search.Item == nil || search.Item.Editable {
		t.Fatal("indeterminate copy reported as editable")
	}

	before, err := os.ReadFile(saved.Path)
	if err != nil {
		t.Fatalf("read before: %v", err)
	}
	result, err := store.Update("wld_duplicate", map[string]interface{}{"title": "Fork"}, saved.Revision)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Status != "authority_indeterminate" {
		t.Fatalf("status = %q, want authority_indeterminate", result.Status)
	}
	afterA, err := os.ReadFile(saved.Path)
	if err != nil {
		t.Fatalf("read after: %v", err)
	}
	afterB, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatalf("read second after: %v", err)
	}
	if string(before) != string(afterA) || string(before) != string(afterB) {
		t.Fatal("an indeterminate update changed a copy")
	}
	if search.Item.Root != rootA {
		t.Fatalf("returned root = %q, want the store root %q", search.Item.Root, rootA)
	}
}

func TestUnavailableRootMakesFoundCopyIndeterminate(t *testing.T) {
	store, reg, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("wld_shadowed"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	addUnavailableRoot(t, reg)

	search, err := store.Lookup("wld_shadowed")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if !search.Found {
		t.Fatal("readable copy was not found")
	}
	if !search.AuthorityIndeterminate {
		t.Fatal("found copy was treated as authoritative while a root was unavailable")
	}
	if !search.Item.AuthorityIndeterminate {
		t.Fatal("item did not carry the indeterminate authority flag")
	}
	if search.Item.Editable {
		t.Fatal("possibly shadowed copy reported as editable")
	}

	item, err := store.Load("wld_shadowed")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !item.AuthorityIndeterminate {
		t.Fatal("Load dropped the authority flag")
	}
	if item.Revision != saved.Revision {
		t.Fatalf("revision = %q, want %q", item.Revision, saved.Revision)
	}

	before, err := os.ReadFile(saved.Path)
	if err != nil {
		t.Fatalf("read before: %v", err)
	}
	result, err := store.Update("wld_shadowed", map[string]interface{}{"title": "Stale Edit"}, saved.Revision)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Status != "authority_indeterminate" {
		t.Fatalf("status = %q, want authority_indeterminate", result.Status)
	}
	after, err := os.ReadFile(saved.Path)
	if err != nil {
		t.Fatalf("read after: %v", err)
	}
	if string(before) != string(after) {
		t.Fatal("update wrote through an indeterminate authority state")
	}
}
