package lorebookstore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/library"
)

// A crash after writing files but before the journal and marker must not leave
// ownership permanently incomplete: a retry verifies the files and adopts them.
func TestInterruptedImportAfterFileWriteIsRecovered(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	legacy := entry("legacy-1", "Mana", "body")
	payload, err := encodeDocument("legacy-1", legacy)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(f.root, "legacy-1.json")
	if err := os.WriteFile(path, payload, 0644); err != nil {
		t.Fatal(err)
	}
	if f.store.Cutover() {
		t.Fatal("the fixture must look like an interrupted import")
	}

	res, err := f.store.Import([]map[string]interface{}{legacy})
	if err != nil {
		t.Fatalf("retry import: %v", err)
	}
	if res.Adopted != 1 || res.Created != 0 {
		t.Fatalf("retry import = %+v, want one adopted entry", res)
	}
	if res.Status != "completed" {
		t.Fatalf("status = %q, want completed", res.Status)
	}
	if !f.store.Cutover() {
		t.Fatal("the retry must complete the ownership transition")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("adopted file missing: %v", err)
	}
	if string(after) != string(payload) {
		t.Fatal("adoption must not rewrite the existing file")
	}

	// Rollback returns ownership to the browser store but leaves adopted files,
	// because this store did not create them in this journal.
	rolled, err := f.store.RollbackImport()
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if rolled.Status != "rolled_back" {
		t.Fatalf("rollback status = %q", rolled.Status)
	}
	if f.store.Cutover() {
		t.Fatal("rollback must clear the marker")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("rollback removed an adopted file it did not create")
	}

	// A retry after rollback must be able to adopt again and re-complete.
	again, err := f.store.Import([]map[string]interface{}{legacy})
	if err != nil {
		t.Fatalf("second retry: %v", err)
	}
	if again.Adopted != 1 || !f.store.Cutover() {
		t.Fatalf("second retry = %+v", again)
	}
}

// Duplicate identities across roots must resolve the same way on every read and
// must be reported, so an attachment pointing at the identity is predictable.
func TestDuplicateIdentitiesAcrossRootsAreReportedAndDeterministic(t *testing.T) {
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	rootA := t.TempDir()
	rootB := t.TempDir()
	for _, root := range []string{rootA, rootB} {
		if _, err := reg.Add(library.Registration{
			Category: library.Lorebooks,
			Path:     root,
			Role:     library.RoleSave | library.RoleWrite,
			Origin:   library.OriginHowl,
			Format:   library.FormatHowl,
		}); err != nil {
			t.Fatalf("register %s: %v", root, err)
		}
	}
	write := func(root, title string) {
		doc := Document{Format: Format, Version: Version, ID: "shared", Entry: entry("shared", title, title)}
		data, err := json.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "shared.json"), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	write(rootA, "From A")
	write(rootB, "From B")

	store := New(reg, t.TempDir())
	first := store.List()
	second := store.List()
	if first.Count != 1 || second.Count != 1 {
		t.Fatalf("duplicate identity produced %d/%d items, want 1", first.Count, second.Count)
	}
	if len(first.Conflicts) != 1 {
		t.Fatalf("conflicts = %v, want one reported conflict", first.Conflicts)
	}
	winner := first.Items[0]
	if winner.ID != "shared" {
		t.Fatalf("winner id = %q", winner.ID)
	}
	if second.Items[0].Root != winner.Root || second.Items[0].Entry["title"] != winner.Entry["title"] {
		t.Fatalf("resolution is not deterministic: %q/%q then %q/%q",
			winner.Root, winner.Entry["title"], second.Items[0].Root, second.Items[0].Entry["title"])
	}
}
