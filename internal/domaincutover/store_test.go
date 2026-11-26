package domaincutover

import (
	"testing"
)

func TestMarkerRequiresRecordAndRollbackRetiresIt(t *testing.T) {
	store := New(t.TempDir(), "character")

	if status := store.Status(); status.Cutover || status.RollbackAvailable {
		t.Fatalf("fresh status = %+v", status)
	}

	items := []Item{{ID: "chr_one", Action: "created", Revision: "rev-one"}}
	if err := store.Record(1, items, "completed", "explicit import"); err != nil {
		t.Fatal(err)
	}
	status := store.Status()
	if !status.Cutover || !status.RollbackAvailable || status.Imported != 1 {
		t.Fatalf("status after import = %+v", status)
	}
	journal, ok := store.ReadJournal()
	if !ok || len(journal.Items) != 1 || journal.Items[0].ID != "chr_one" {
		t.Fatalf("journal = %+v, ok = %v", journal, ok)
	}

	if err := store.MarkRolledBack("undone"); err != nil {
		t.Fatal(err)
	}
	status = store.Status()
	if status.Cutover || status.RollbackAvailable {
		t.Fatalf("status after rollback = %+v", status)
	}
	journal, ok = store.ReadJournal()
	if !ok || journal.Status != "rolled_back" {
		t.Fatalf("journal after rollback = %+v, ok = %v", journal, ok)
	}
}

func TestRecordIsScopedToItsDomain(t *testing.T) {
	dir := t.TempDir()
	characters := New(dir, "character")
	worlds := New(dir, "world")
	if err := characters.Record(1, []Item{{ID: "chr_one"}}, "completed", ""); err != nil {
		t.Fatal(err)
	}
	if worlds.Status().Cutover {
		t.Fatal("a character marker was read as a world cutover")
	}
}
