package worldstore

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/library"
)

func TestLoadLoreReadsCanonicalItemFolder(t *testing.T) {
	store, _, _, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("wld_lore"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	lore := `{"entries":[{"id":"lb_manual","content":"ash falls","scope":"world","ownerId":"wld_lore"}]}`
	if err := os.WriteFile(filepath.Join(root, saved.ID, "world_lore.json"), []byte(lore), 0644); err != nil {
		t.Fatalf("write lore: %v", err)
	}

	file, err := store.LoadLore(saved.ID, "world_lore.json")
	if err != nil {
		t.Fatalf("load lore: %v", err)
	}
	if file.AuthorityIndeterminate {
		t.Fatal("single readable copy reported indeterminate")
	}
	if len(file.Entries) != 1 {
		t.Fatalf("entries = %d, want 1", len(file.Entries))
	}
	if file.Item.Root != root {
		t.Fatalf("item root = %q, want %q", file.Item.Root, root)
	}
}

func TestLoadLoreMissingFileIsEmptyNotAbsentItem(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("wld_lore_absent"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	file, err := store.LoadLore(saved.ID, "world_lore.json")
	if err != nil {
		t.Fatalf("load missing lore: %v", err)
	}
	if len(file.Entries) != 0 {
		t.Fatalf("entries = %d, want 0", len(file.Entries))
	}
}

func TestLoadLoreUnavailableSiblingIsIndeterminate(t *testing.T) {
	store, reg, _, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("wld_lore_indeterminate"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, saved.ID, "world_lore.json"), []byte(`{"entries":[{"id":"lb_1"}]}`), 0644); err != nil {
		t.Fatalf("write lore: %v", err)
	}
	addUnavailableRoot(t, reg)

	file, err := store.LoadLore(saved.ID, "world_lore.json")
	if err != nil {
		t.Fatalf("load lore: %v", err)
	}
	if !file.AuthorityIndeterminate {
		t.Fatal("unavailable sibling did not mark lore indeterminate")
	}
	if len(file.Entries) != 0 {
		t.Fatalf("indeterminate lore returned %d entries", len(file.Entries))
	}
	if len(file.UnavailableRoots) != 1 {
		t.Fatalf("unavailable roots = %v, want one", file.UnavailableRoots)
	}
}

func TestLoadLoreDoesNotModifyTheItemFolder(t *testing.T) {
	store, _, _, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("wld_lore_readonly"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	itemDir := filepath.Join(root, saved.ID)
	payload := []byte(`{"entries":[{"id":"lb_1","content":"ash"}]}`)
	if err := os.WriteFile(filepath.Join(itemDir, "world_lore.json"), payload, 0644); err != nil {
		t.Fatalf("write lore: %v", err)
	}
	before, err := os.ReadDir(itemDir)
	if err != nil {
		t.Fatalf("read item folder: %v", err)
	}

	if _, err := store.LoadLore(saved.ID, "automated_world_lore.json"); err != nil {
		t.Fatalf("load lore: %v", err)
	}

	after, err := os.ReadFile(filepath.Join(itemDir, "world_lore.json"))
	if err != nil {
		t.Fatalf("re-read lore: %v", err)
	}
	if !bytes.Equal(after, payload) {
		t.Fatalf("lore file changed on read: %s", after)
	}
	entries, err := os.ReadDir(itemDir)
	if err != nil {
		t.Fatalf("re-read item folder: %v", err)
	}
	if len(entries) != len(before) {
		t.Fatalf("item folder gained or lost files on read: before=%d after=%d", len(before), len(entries))
	}
	for i := range entries {
		if entries[i].Name() != before[i].Name() {
			t.Fatalf("item folder changed on read: before=%v after=%v", names(before), names(entries))
		}
	}
}

func names(entries []os.DirEntry) []string {
	out := make([]string, 0, len(entries))
	for _, entry := range entries {
		out = append(out, entry.Name())
	}
	return out
}

func TestLoadLoreRejectsFilesOutsideTheContract(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("wld_lore_reject"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	for _, filename := range []string{"../world.json", "scenario_lore.json", "world.json"} {
		if _, err := store.LoadLore(saved.ID, filename); !errors.Is(err, ErrInvalidDocument) {
			t.Errorf("LoadLore(%q) error = %v, want ErrInvalidDocument", filename, err)
		}
	}
}
