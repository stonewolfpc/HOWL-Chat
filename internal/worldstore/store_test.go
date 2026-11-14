package worldstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/library"
)

func newTestStore(t *testing.T, category library.Category, role library.Role) (*Store, *library.Registry, library.Registration, string) {
	t.Helper()
	root := t.TempDir()
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatalf("library registry: %v", err)
	}
	registration, err := reg.Add(library.Registration{
		Category: category,
		Path:     root,
		Role:     role,
		Origin:   library.OriginHowl,
		Format:   library.FormatHowl,
	})
	if err != nil {
		t.Fatalf("register root: %v", err)
	}
	store, err := New(reg, category)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	return store, reg, registration, root
}

func sampleWorld(id string) map[string]interface{} {
	return map[string]interface{}{
		"id":          id,
		"title":       "Ashfall",
		"details":     "A city under falling ash",
		"description": "old description",
		"createdAt":   "2026-01-01T00:00:00Z",
		"updatedAt":   "2026-01-01T00:00:00Z",
		"characters": []interface{}{
			map[string]interface{}{"id": "char_aria", "name": "Aria"},
		},
		"howl": map[string]interface{}{
			"lorebooks": []interface{}{
				map[string]interface{}{"id": "1794319623163", "title": "Tide Oath"},
			},
			"notes":      "keep these notes",
			"customFlag": "keep",
		},
		"metadata": map[string]interface{}{
			"version": 1,
			"seed":    7,
			"extra":   "keep",
		},
		"legacyNotes": "keep me",
	}
}

func TestSaveCreatesIdentityFolderAndLoads(t *testing.T) {
	store, _, registration, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)

	result, err := store.Save(sampleWorld("1794319623163"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if result.Status != "saved" {
		t.Fatalf("status = %q, want saved", result.Status)
	}
	wantPath := filepath.Join(root, "1794319623163", "world.json")
	if result.Path != wantPath {
		t.Fatalf("path = %q, want %q", result.Path, wantPath)
	}
	if _, err := os.Stat(wantPath); err != nil {
		t.Fatalf("document not written: %v", err)
	}

	item, err := store.Load("1794319623163")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if item.Revision != result.Revision {
		t.Fatalf("revision = %q, want %q", item.Revision, result.Revision)
	}
	if item.RevisionMismatch {
		t.Fatal("freshly written revision reported as mismatched")
	}
	if !item.Editable {
		t.Fatal("text writable root should be editable")
	}
	if item.LibraryID != registration.ID {
		t.Fatalf("library id = %q, want %q", item.LibraryID, registration.ID)
	}
	if item.Document["legacyNotes"] != "keep me" {
		t.Fatalf("unknown field lost: %v", item.Document["legacyNotes"])
	}
}

func TestSaveScenarioUsesScenarioEnvelope(t *testing.T) {
	store, _, _, root := newTestStore(t, library.Scenarios, library.RoleSave|library.RoleWrite)
	id := "scenario-1794319623163-abcdef123"

	result, err := store.Save(map[string]interface{}{"id": id, "title": "Tide Watch"})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if result.Path != filepath.Join(root, id, "scenario.json") {
		t.Fatalf("path = %q", result.Path)
	}
	item, err := store.Load(id)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if item.Format != FormatScenario {
		t.Fatalf("format = %q, want %q", item.Format, FormatScenario)
	}
}

func TestSaveRejectsUnsafeIdentity(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	if _, err := store.Save(map[string]interface{}{"id": "../escape"}); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("error = %v, want ErrInvalidDocument", err)
	}
}

func TestSaveWithoutDestination(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleRead)
	if _, err := store.Save(sampleWorld("world_read_only")); !errors.Is(err, ErrNoDestination) {
		t.Fatalf("error = %v, want ErrNoDestination", err)
	}
}

func TestSaveCollisionKeepsExistingDocument(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	first, err := store.Save(sampleWorld("1794319623163"))
	if err != nil {
		t.Fatalf("first save: %v", err)
	}

	second := sampleWorld("1794319623163")
	second["title"] = "Overwritten"
	result, err := store.Save(second)
	if err != nil {
		t.Fatalf("second save: %v", err)
	}
	if result.Status != "collision" {
		t.Fatalf("status = %q, want collision", result.Status)
	}

	item, err := store.Load("1794319623163")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if item.Document["title"] != "Ashfall" {
		t.Fatalf("existing document changed: %v", item.Document["title"])
	}
	if item.Revision != first.Revision {
		t.Fatal("existing revision changed")
	}
}

func TestUpdatePreservesUnknownFields(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("1794319623163"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	patch := map[string]interface{}{
		"title": "Ashfall II",
		"howl": map[string]interface{}{
			"lorebooks": []interface{}{
				map[string]interface{}{"id": "1794319623770", "title": "Deep Current"},
			},
		},
	}
	result, err := store.Update("1794319623163", patch, saved.Revision)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Status != "updated" {
		t.Fatalf("status = %q, want updated", result.Status)
	}

	item, err := store.Load("1794319623163")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if item.Document["title"] != "Ashfall II" {
		t.Fatalf("known field not applied: %v", item.Document["title"])
	}
	if item.Document["legacyNotes"] != "keep me" {
		t.Fatal("unknown top-level field was dropped")
	}
	howl, ok := item.Document["howl"].(map[string]interface{})
	if !ok {
		t.Fatal("howl map missing")
	}
	if howl["notes"] != "keep these notes" || howl["customFlag"] != "keep" {
		t.Fatalf("unknown howl fields were dropped: %v", howl)
	}
	metadata, ok := item.Document["metadata"].(map[string]interface{})
	if !ok || metadata["extra"] != "keep" {
		t.Fatalf("unknown metadata field was dropped: %v", metadata)
	}
	if item.Revision == saved.Revision {
		t.Fatal("content change did not change the revision")
	}
}

func TestUpdateClearsAllowlistedNullOnly(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("1794319623163"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	result, err := store.Update("1794319623163", map[string]interface{}{
		"description": nil,
		"legacyNotes": nil,
	}, saved.Revision)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Status != "updated" {
		t.Fatalf("status = %q", result.Status)
	}

	item, err := store.Load("1794319623163")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if _, exists := item.Document["description"]; exists {
		t.Fatal("allowlisted field was not cleared")
	}
	if item.Document["legacyNotes"] != "keep me" {
		t.Fatal("null on a non-allowlisted field erased data")
	}
}

func TestUpdateRevisionConflictChangesNothing(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("1794319623163"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	before, err := os.ReadFile(saved.Path)
	if err != nil {
		t.Fatalf("read before: %v", err)
	}

	result, err := store.Update("1794319623163", map[string]interface{}{"title": "Race"}, "rv1-0000000000000000")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Status != "revision_conflict" {
		t.Fatalf("status = %q, want revision_conflict", result.Status)
	}
	if result.CurrentRevision != saved.Revision {
		t.Fatalf("current revision = %q, want %q", result.CurrentRevision, saved.Revision)
	}

	after, err := os.ReadFile(saved.Path)
	if err != nil {
		t.Fatalf("read after: %v", err)
	}
	if string(before) != string(after) {
		t.Fatal("document changed despite revision conflict")
	}
}

func TestUpdateRequiresBaseRevision(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	if _, err := store.Save(sampleWorld("1794319623163")); err != nil {
		t.Fatalf("save: %v", err)
	}

	result, err := store.Update("1794319623163", map[string]interface{}{"title": "No Base"}, "")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Status != "revision_conflict" {
		t.Fatalf("status = %q, want revision_conflict", result.Status)
	}
}

func TestUpdateRejectsImmutableFields(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("1794319623163"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	if _, err := store.Update("1794319623163", map[string]interface{}{"id": "other"}, saved.Revision); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("identity change error = %v, want ErrInvalidDocument", err)
	}
	if _, err := store.Update("1794319623163", map[string]interface{}{"createdAt": "2020-01-01T00:00:00Z"}, saved.Revision); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("creation-time change error = %v, want ErrInvalidDocument", err)
	}
}

func TestUpdateRequiresWriteRole(t *testing.T) {
	store, reg, registration, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("1794319623163"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	if err := reg.SetRole(registration.ID, library.RoleSave); err != nil {
		t.Fatalf("set role: %v", err)
	}
	result, err := store.Update("1794319623163", map[string]interface{}{"title": "Denied"}, saved.Revision)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Status != "permission_denied" {
		t.Fatalf("status = %q, want permission_denied", result.Status)
	}

	if err := reg.SetRole(registration.ID, library.RoleSave|library.RoleWrite); err != nil {
		t.Fatalf("restore role: %v", err)
	}
	result, err = store.Update("1794319623163", map[string]interface{}{"title": "Allowed"}, saved.Revision)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Status != "updated" {
		t.Fatalf("status = %q, want updated", result.Status)
	}
}

func TestLoadFlagsRevisionMismatch(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("1794319623163"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	payload, err := os.ReadFile(saved.Path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var envelope map[string]interface{}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	document, ok := envelope["document"].(map[string]interface{})
	if !ok {
		t.Fatal("envelope document missing")
	}
	document["title"] = "Tampered"
	tampered, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(saved.Path, tampered, 0644); err != nil {
		t.Fatalf("write tampered: %v", err)
	}

	item, err := store.Load("1794319623163")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !item.RevisionMismatch {
		t.Fatal("tampered document was not flagged")
	}
	if item.Document["title"] != "Tampered" {
		t.Fatal("mismatched document should stay readable")
	}
}

func TestLoadNotFound(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	if _, err := store.Load("wld_missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("error = %v, want ErrNotFound", err)
	}
}
