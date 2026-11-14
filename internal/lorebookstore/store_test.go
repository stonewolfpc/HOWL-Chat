package lorebookstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/library"
)

type fixture struct {
	store     *Store
	reg       *library.Registry
	root      string
	configDir string
}

func newFixture(t *testing.T, role library.Role) *fixture {
	t.Helper()
	root := t.TempDir()
	configDir := t.TempDir()
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatalf("library.New: %v", err)
	}
	if _, err := reg.Add(library.Registration{
		Category: library.Lorebooks,
		Path:     root,
		Role:     role,
		Origin:   library.OriginHowl,
		Format:   library.FormatHowl,
	}); err != nil {
		t.Fatalf("register lorebook root: %v", err)
	}
	return &fixture{store: New(reg, configDir), reg: reg, root: root, configDir: configDir}
}

func entry(id, title, content string) map[string]interface{} {
	return map[string]interface{}{
		"id":        id,
		"title":     title,
		"content":   content,
		"scope":     "world",
		"enabled":   true,
		"updatedAt": "2026-11-10T00:00:00Z",
	}
}

func mustImport(t *testing.T, f *fixture, entries ...map[string]interface{}) ImportResult {
	t.Helper()
	res, err := f.store.Import(entries)
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if res.Status != "completed" && res.Status != "partial" && res.Status != "nothing_to_do" {
		t.Fatalf("import status = %q (%s)", res.Status, res.Message)
	}
	return res
}

func TestSaveRequiresExplicitCutover(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	if f.store.Cutover() {
		t.Fatal("a new store must not be cut over")
	}
	if _, err := f.store.Save(entry("a", "A", "x"), false); !errors.Is(err, ErrNotCutover) {
		t.Fatalf("save before import error = %v, want ErrNotCutover", err)
	}
	if _, err := os.Stat(filepath.Join(f.root, "a.json")); !os.IsNotExist(err) {
		t.Fatal("save before import must not create a file")
	}
}

func TestImportCreatesCanonicalFilesAndMarker(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	res := mustImport(t, f, entry("legacy-1", "Mana Flow", "Magic flows through ley lines."))
	if res.Created != 1 || res.Status != "completed" {
		t.Fatalf("import result = %+v", res)
	}
	if !f.store.Cutover() {
		t.Fatal("import must set the cutover marker")
	}
	path := filepath.Join(f.root, "legacy-1.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("canonical file missing: %v", err)
	}
	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("decode document: %v", err)
	}
	if doc.Format != Format || doc.ID != "legacy-1" || doc.Entry["title"] != "Mana Flow" {
		t.Fatalf("unexpected document: %+v", doc)
	}
}

func TestSaveReplaceKeepsIdentityAcrossTitleChange(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, entry("legacy-1", "Old Title", "body"))

	renamed := entry("legacy-1", "New Title", "body")
	res, err := f.store.Save(renamed, true)
	if err != nil {
		t.Fatalf("replace: %v", err)
	}
	if res.Status != "replaced" || res.ID != "legacy-1" {
		t.Fatalf("replace result = %+v", res)
	}
	list := f.store.List()
	if list.Count != 1 {
		t.Fatalf("list count = %d, want 1", list.Count)
	}
	if got := list.Items[0].Entry["title"]; got != "New Title" {
		t.Fatalf("title = %v, want New Title", got)
	}
	if list.Items[0].ID != "legacy-1" {
		t.Fatalf("identity changed: %q", list.Items[0].ID)
	}
	if _, err := os.Stat(filepath.Join(f.root, "New Title.json")); !os.IsNotExist(err) {
		t.Fatal("a rename must not create a title-based file")
	}
}

func TestSaveRejectsReplacementOnSaveOnlyRoot(t *testing.T) {
	f := newFixture(t, library.RoleSave)
	// Import cannot run without write permission, so create the file directly.
	doc := Document{Format: Format, Version: Version, ID: "legacy-1", Entry: entry("legacy-1", "A", "a")}
	payload, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.root, "legacy-1.json"), payload, 0644); err != nil {
		t.Fatal(err)
	}
	// Simulate a completed cutover without writing through the store.
	if err := os.WriteFile(filepath.Join(f.configDir, "lorebook-cutover.json"), []byte(`{"version":1}`), 0644); err != nil {
		t.Fatal(err)
	}
	res, err := f.store.Save(entry("legacy-1", "B", "b"), true)
	if err != nil {
		t.Fatalf("replace on save-only root: %v", err)
	}
	if res.Status != "read_only" {
		t.Fatalf("status = %q, want read_only", res.Status)
	}
	if _, err := os.Stat(filepath.Join(f.root, "legacy-1.json")); err != nil {
		t.Fatal("the existing file must survive a refused replace")
	}
}

func TestNewEntryGetsStableIdentityAndSaves(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, entry("legacy-1", "Seed", "seed"))

	created := map[string]interface{}{"title": "New Entry", "content": "body", "scope": "world"}
	res, err := f.store.Save(created, false)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if res.Status != "saved" || !validID.MatchString(res.ID) {
		t.Fatalf("create result = %+v", res)
	}
	if !strings.HasPrefix(res.ID, "lb_") {
		t.Fatalf("new identity = %q, want an lb_ identity", res.ID)
	}
	list := f.store.List()
	if list.Count != 2 {
		t.Fatalf("list count = %d, want 2", list.Count)
	}
}

func TestDeleteRemovesOnlyTheCanonicalEntry(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, entry("a", "A", "a"), entry("b", "B", "b"))

	res, err := f.store.Delete("a")
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if res.Status != "deleted" {
		t.Fatalf("delete status = %q", res.Status)
	}
	if _, err := os.Stat(filepath.Join(f.root, "a.json")); !os.IsNotExist(err) {
		t.Fatal("deleted file still exists")
	}
	if _, err := os.Stat(filepath.Join(f.root, "b.json")); err != nil {
		t.Fatal("unrelated entry was removed")
	}
}

func TestDeleteRejectsPathAndIdentityAbuse(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, entry("a", "A", "a"))
	outside := filepath.Join(t.TempDir(), "outside.json")
	if err := os.WriteFile(outside, []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"../outside", "a/../../outside", "", "bad id"} {
		res, err := f.store.Delete(id)
		if err != nil {
			t.Fatalf("delete(%q) returned error: %v", id, err)
		}
		if res.Status == "deleted" {
			t.Fatalf("delete(%q) must not report deletion", id)
		}
	}
	if _, err := os.Stat(outside); err != nil {
		t.Fatal("an outside file was removed")
	}
}

func TestDeleteRequiresWritableRoot(t *testing.T) {
	f := newFixture(t, library.RoleRead)
	doc := Document{Format: Format, Version: Version, ID: "a", Entry: entry("a", "A", "a")}
	payload, _ := json.Marshal(doc)
	if err := os.WriteFile(filepath.Join(f.root, "a.json"), payload, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.configDir, "lorebook-cutover.json"), []byte(`{"version":1}`), 0644); err != nil {
		t.Fatal(err)
	}
	res, err := f.store.Delete("a")
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if res.Status != "read_only" {
		t.Fatalf("status = %q, want read_only", res.Status)
	}
	if _, err := os.Stat(filepath.Join(f.root, "a.json")); err != nil {
		t.Fatal("a read-only entry was removed")
	}
}

func TestImportIsIdempotentAndReportsConflicts(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	first := mustImport(t, f, entry("a", "A", "a"), entry("b", "B", "b"))
	if first.Created != 2 {
		t.Fatalf("first import created = %d, want 2", first.Created)
	}
	second := mustImport(t, f, entry("a", "A", "a"), entry("b", "B", "b"))
	if second.Created != 0 || second.Skipped != 2 || second.Status != "nothing_to_do" {
		t.Fatalf("second import = %+v", second)
	}

	// A different entry claiming the same identity is a conflict, not an overwrite.
	conflict := mustImport(t, f, entry("a", "A", "different body"))
	if conflict.Conflicts != 1 || conflict.Status != "partial" {
		t.Fatalf("conflict import = %+v", conflict)
	}
	list := f.store.List()
	for _, item := range list.Items {
		if item.ID == "a" && item.Entry["content"] != "a" {
			t.Fatalf("conflicting import overwrote the canonical entry: %v", item.Entry["content"])
		}
	}
}

func TestImportReassignsDuplicateLegacyIdentities(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	res := mustImport(t, f, entry("dup", "First", "one"), entry("dup", "Second", "two"))
	if res.Created != 2 {
		t.Fatalf("created = %d, want 2", res.Created)
	}
	reassigned := 0
	for _, item := range res.Items {
		if item.Action == "reassigned" || strings.Contains(item.Detail, "imported as") {
			reassigned++
		}
	}
	if reassigned == 0 {
		t.Fatal("duplicate identities must be reported")
	}
	list := f.store.List()
	if list.Count != 2 {
		t.Fatalf("list count = %d, want 2", list.Count)
	}
	for _, item := range list.Items {
		if !validID.MatchString(item.ID) {
			t.Fatalf("invalid identity after import: %q", item.ID)
		}
	}
}

func TestRollbackRemovesOnlyImportedFiles(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	manual := Document{Format: Format, Version: Version, ID: "keep", Entry: entry("keep", "Keep", "keep")}
	payload, _ := json.Marshal(manual)
	if err := os.WriteFile(filepath.Join(f.root, "keep.json"), payload, 0644); err != nil {
		t.Fatal(err)
	}
	mustImport(t, f, entry("a", "A", "a"))

	res, err := f.store.RollbackImport()
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if res.Status != "rolled_back" {
		t.Fatalf("rollback status = %q", res.Status)
	}
	if _, err := os.Stat(filepath.Join(f.root, "a.json")); !os.IsNotExist(err) {
		t.Fatal("imported file was not removed")
	}
	if _, err := os.Stat(filepath.Join(f.root, "keep.json")); err != nil {
		t.Fatal("rollback removed a file it did not create")
	}
	if f.store.Cutover() {
		t.Fatal("rollback must clear the cutover marker")
	}
	again, err := f.store.RollbackImport()
	if err != nil {
		t.Fatalf("second rollback: %v", err)
	}
	if again.Status != "nothing_to_do" {
		t.Fatalf("second rollback status = %q", again.Status)
	}
}

func TestListIgnoresForeignFormatsAndReportsUnavailableRoots(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, entry("a", "A", "a"))
	foreign := map[string]interface{}{"format": "howl-content-export-v1", "content": map[string]interface{}{"id": "x"}}
	data, _ := json.Marshal(foreign)
	if err := os.WriteFile(filepath.Join(f.root, "export.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
	list := f.store.List()
	if list.Count != 1 {
		t.Fatalf("list count = %d, want 1", list.Count)
	}
	if list.Ignored != 1 {
		t.Fatalf("ignored = %d, want 1", list.Ignored)
	}

	if err := os.RemoveAll(f.root); err != nil {
		t.Fatal(err)
	}
	gone := f.store.List()
	if gone.Count != 1 {
		t.Fatalf("items should remain known from the marker, got %d", gone.Count)
	}
	if len(gone.Roots) != 1 || gone.Roots[0].Available {
		t.Fatalf("unavailable root not reported: %+v", gone.Roots)
	}
}

func TestImportWithoutSaveDestinationFails(t *testing.T) {
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	store := New(reg, t.TempDir())
	res, err := store.Import([]map[string]interface{}{entry("a", "A", "a")})
	if !errors.Is(err, ErrNoDestination) {
		t.Fatalf("error = %v, want ErrNoDestination", err)
	}
	if res.Status == "completed" {
		t.Fatal("import must not report success without a destination")
	}
	if store.Cutover() {
		t.Fatal("no marker may be written without a destination")
	}
}

func TestListSkipsSymlinkedEntry(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, entry("a", "A", "a"))
	outside := filepath.Join(t.TempDir(), "outside.json")
	doc := Document{Format: Format, Version: Version, ID: "outside", Entry: entry("outside", "Outside", "x")}
	payload, _ := json.Marshal(doc)
	if err := os.WriteFile(outside, payload, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(f.root, "link.json")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	list := f.store.List()
	if list.Count != 1 {
		t.Fatalf("list count = %d, want 1", list.Count)
	}
	for _, item := range list.Items {
		if item.ID == "outside" {
			t.Fatal("a symlinked entry was listed")
		}
	}
}

func TestSaveReplacePreservesOmittedOwner(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	owned := entry("owned-1", "Aria secret", "body")
	owned["scope"] = "character"
	owned["visibility"] = "private"
	owned["ownerId"] = "Aria"
	mustImport(t, f, owned)

	// A legacy client that does not know about ownership omits the key entirely.
	edited := entry("owned-1", "Aria secret", "edited body")
	edited["scope"] = "character"
	edited["visibility"] = "private"
	res, err := f.store.Save(edited, true)
	if err != nil || res.Status != "replaced" {
		t.Fatalf("replace = %+v, %v", res, err)
	}
	list := f.store.List()
	if len(list.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(list.Items))
	}
	if got := list.Items[0].Entry["ownerId"]; got != "Aria" {
		t.Fatalf("ownerId = %v, want Aria (omitted field must be preserved)", got)
	}
}

func TestSaveReplaceHonorsExplicitOwnerChangeAndClear(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	owned := entry("owned-2", "Bram secret", "body")
	owned["scope"] = "character"
	owned["visibility"] = "private"
	owned["ownerId"] = "Aria"
	mustImport(t, f, owned)

	changed := entry("owned-2", "Bram secret", "body")
	changed["scope"] = "character"
	changed["visibility"] = "private"
	changed["ownerId"] = "Bram"
	if _, err := f.store.Save(changed, true); err != nil {
		t.Fatalf("owner change: %v", err)
	}
	if got := f.store.List().Items[0].Entry["ownerId"]; got != "Bram" {
		t.Fatalf("ownerId after change = %v, want Bram", got)
	}

	cleared := entry("owned-2", "Bram secret", "body")
	cleared["scope"] = "character"
	cleared["visibility"] = "private"
	cleared["ownerId"] = ""
	if _, err := f.store.Save(cleared, true); err != nil {
		t.Fatalf("owner clear: %v", err)
	}
	if got := f.store.List().Items[0].Entry["ownerId"]; got != "" {
		t.Fatalf("ownerId after clear = %v, want empty", got)
	}
}
