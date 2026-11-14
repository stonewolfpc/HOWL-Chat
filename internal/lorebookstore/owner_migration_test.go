package lorebookstore

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/library"
)

func characterEntry(id, ownerID, ownerName string) map[string]interface{} {
	m := entry(id, id, "body-"+id)
	m["scope"] = "character"
	m["visibility"] = "private"
	m["ownerId"] = ownerID
	if ownerName != "" {
		m["ownerName"] = ownerName
	}
	return m
}

func ownerIdentities(pairs ...string) []OwnerIdentity {
	out := make([]OwnerIdentity, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, OwnerIdentity{ID: pairs[i], Name: pairs[i+1]})
	}
	return out
}

func TestMigrateOwnersRewritesUniqueNameAndIsIdempotent(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, characterEntry("entry-1", "Aria", ""))

	res, err := f.store.MigrateOwners(ownerIdentities("char_aria", "Aria"), "v1")
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if res.Status != "completed" || res.Migrated != 1 {
		t.Fatalf("result = %+v", res)
	}
	list := f.store.List()
	got := list.Items[0].Entry
	if got["ownerId"] != "char_aria" {
		t.Fatalf("ownerId = %v, want char_aria", got["ownerId"])
	}
	if got["ownerName"] != "Aria" {
		t.Fatalf("ownerName = %v, want Aria", got["ownerName"])
	}
	if got["ownerMigratedFrom"] != "Aria" || got["ownerMigrationVersion"] != "v1" {
		t.Fatalf("migration metadata missing: %+v", got)
	}

	again, err := f.store.MigrateOwners(ownerIdentities("char_aria", "Aria"), "v1")
	if err != nil {
		t.Fatalf("re-run: %v", err)
	}
	if again.Migrated != 0 || again.Already != 1 || again.Status != "nothing_to_do" {
		t.Fatalf("re-run result = %+v", again)
	}
	if f.store.List().Items[0].Entry["ownerId"] != "char_aria" {
		t.Fatal("re-running changed the migrated identity")
	}
}

func TestMigrateOwnersLeavesAmbiguousAndUnknownNames(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, characterEntry("entry-dup", "Aria", ""), characterEntry("entry-unknown", "Nobody", ""))

	res, err := f.store.MigrateOwners(ownerIdentities(
		"char_aria_1", "Aria",
		"char_aria_2", "Aria",
	), "v1")
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if res.Ambiguous != 1 || res.Unknown != 1 || res.Migrated != 0 {
		t.Fatalf("result = %+v", res)
	}
	byID := map[string]map[string]interface{}{}
	for _, item := range f.store.List().Items {
		byID[item.ID] = item.Entry
	}
	if byID["entry-dup"]["ownerId"] != "Aria" {
		t.Fatalf("ambiguous name was guessed: %v", byID["entry-dup"]["ownerId"])
	}
	if byID["entry-unknown"]["ownerId"] != "Nobody" {
		t.Fatalf("unknown name was guessed: %v", byID["entry-unknown"]["ownerId"])
	}
}

func TestMigrateOwnersNeverOverwritesStableID(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, characterEntry("entry-1", "char_bram", "Bram"))

	res, err := f.store.MigrateOwners(ownerIdentities("char_aria", "Aria", "char_bram", "Bram"), "v1")
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if res.Already != 1 || res.Migrated != 0 {
		t.Fatalf("result = %+v", res)
	}
	if got := f.store.List().Items[0].Entry["ownerId"]; got != "char_bram" {
		t.Fatalf("stable id overwritten: %v", got)
	}
}

func TestMigrateOwnersRejectsMalformedIdentities(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, characterEntry("entry-1", "Aria", ""))

	for name, identities := range map[string][]OwnerIdentity{
		"empty":          {},
		"missing id":     {{Name: "Aria"}},
		"missing name":   {{ID: "char_aria"}},
		"conflicting id": {{ID: "char_aria", Name: "Aria"}, {ID: "char_aria", Name: "Bram"}},
	} {
		if _, err := f.store.MigrateOwners(identities, "v1"); !errors.Is(err, ErrInvalidIdentityMap) {
			t.Errorf("%s: error = %v, want ErrInvalidIdentityMap", name, err)
		}
	}
	if _, err := f.store.MigrateOwners(ownerIdentities("char_aria", "Aria"), ""); !errors.Is(err, ErrInvalidIdentityMap) {
		t.Fatalf("missing version error = %v", err)
	}
	if got := f.store.List().Items[0].Entry["ownerId"]; got != "Aria" {
		t.Fatalf("malformed input changed ownership: %v", got)
	}
}

func TestMigrateOwnersRequiresCutover(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	if _, err := f.store.MigrateOwners(ownerIdentities("char_aria", "Aria"), "v1"); !errors.Is(err, ErrNotCutover) {
		t.Fatalf("error = %v, want ErrNotCutover", err)
	}
}

func TestRollbackOwnersRestoresPriorValues(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, characterEntry("entry-1", "Aria", ""))

	if _, err := f.store.MigrateOwners(ownerIdentities("char_aria", "Aria"), "v1"); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	journalPath := filepath.Join(f.configDir, "lorebook-owner-migration-journal.json")
	if _, err := os.Stat(journalPath); err != nil {
		t.Fatalf("migration journal missing: %v", err)
	}

	res, err := f.store.RollbackOwners()
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if res.Status != "rolled_back" || res.Restored != 1 {
		t.Fatalf("rollback result = %+v", res)
	}
	got := f.store.List().Items[0].Entry
	if got["ownerId"] != "Aria" {
		t.Fatalf("ownerId = %v, want Aria", got["ownerId"])
	}
	if _, ok := got["ownerName"]; ok {
		t.Fatalf("rollback left ownerName: %v", got["ownerName"])
	}
	if _, ok := got["ownerMigratedFrom"]; ok {
		t.Fatal("rollback left migration metadata")
	}
	if got["content"] != "body-entry-1" {
		t.Fatalf("rollback changed unrelated content: %v", got["content"])
	}
	if _, err := os.Stat(filepath.Join(f.configDir, "lorebook-owner-migration.json")); err == nil {
		t.Fatal("rollback must clear the migration marker")
	}

	if _, err := f.store.MigrateOwners(ownerIdentities("char_aria", "Aria"), "v2"); err != nil {
		t.Fatalf("migrate after rollback: %v", err)
	}
	if got := f.store.List().Items[0].Entry["ownerId"]; got != "char_aria" {
		t.Fatalf("ownerId after re-migrate = %v", got)
	}
}

func TestRollbackOwnersSkipsEntriesChangedAfterMigration(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	mustImport(t, f, characterEntry("entry-1", "Aria", ""))
	if _, err := f.store.MigrateOwners(ownerIdentities("char_aria", "Aria"), "v1"); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	edited := characterEntry("entry-1", "char_bram", "Bram")
	if _, err := f.store.Save(edited, true); err != nil {
		t.Fatalf("edit: %v", err)
	}
	res, err := f.store.RollbackOwners()
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if res.Skipped != 1 || res.Restored != 0 {
		t.Fatalf("rollback result = %+v", res)
	}
	if got := f.store.List().Items[0].Entry["ownerId"]; got != "char_bram" {
		t.Fatalf("later edit was overwritten: %v", got)
	}
}

func TestSaveReplacePreservesOmittedOwnerNameAndClearsOnOwnerClear(t *testing.T) {
	f := newFixture(t, library.RoleSave|library.RoleWrite)
	owned := characterEntry("entry-1", "char_aria", "Aria")
	mustImport(t, f, owned)

	edited := characterEntry("entry-1", "char_aria", "")
	if _, err := f.store.Save(edited, true); err != nil {
		t.Fatalf("replace: %v", err)
	}
	if got := f.store.List().Items[0].Entry["ownerName"]; got != "Aria" {
		t.Fatalf("ownerName = %v, want preserved Aria", got)
	}

	cleared := characterEntry("entry-1", "", "")
	if _, err := f.store.Save(cleared, true); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if got := f.store.List().Items[0].Entry["ownerName"]; got != "" {
		t.Fatalf("ownerName after clear = %v, want empty", got)
	}
}
