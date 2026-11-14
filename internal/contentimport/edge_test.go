package contentimport

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/library"
	"howl-chat/internal/worldstore"
)

func TestMissingInvalidAndExplicitNewIdentity(t *testing.T) {
	service, _, registry, _, _ := newImportTest(t, library.Worlds)
	destination := destinationID(t, registry, library.Worlds)

	missing := Request{
		Category:             "worlds",
		DestinationLibraryID: destination,
		BrowserRecords:       []map[string]interface{}{{"title": "No Identity"}},
	}
	plan := mustPreview(t, service, missing)
	if plan.Items[0].Reason != "missing_document_id" {
		t.Fatalf("missing reason = %q", plan.Items[0].Reason)
	}

	invalid := Request{
		Category:             "worlds",
		DestinationLibraryID: destination,
		BrowserRecords:       []map[string]interface{}{{"id": "../escape", "title": "Bad Identity"}},
	}
	plan = mustPreview(t, service, invalid)
	if plan.Items[0].Reason != "invalid_document_id" {
		t.Fatalf("invalid reason = %q", plan.Items[0].Reason)
	}

	resolved := Request{
		Category:             "worlds",
		DestinationLibraryID: destination,
		BrowserRecords:       []map[string]interface{}{{"id": "../escape", "title": "Resolved Identity"}},
		Resolutions: []Resolution{{
			SourceKey: "browser:0", NewID: true, Actor: "test-user",
			Evidence: "user explicitly chose a new identity",
		}},
	}
	plan = mustPreview(t, service, resolved)
	if plan.Ready != 1 || !worldstore.ValidID(plan.Items[0].ProposedID) {
		t.Fatalf("resolution plan = %+v", plan)
	}
	if plan.Items[0].Evidence != "user explicitly chose a new identity" {
		t.Fatalf("evidence = %q", plan.Items[0].Evidence)
	}
}

func TestLegacyUnmappedOwnerQuarantines(t *testing.T) {
	service, _, registry, _, _ := newImportTest(t, library.Worlds)
	source := t.TempDir()
	folder := filepath.Join(source, "Wrong Owner")
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(folder, "world.json"), []byte(`{"id":"wld_owner","title":"Wrong Owner"}`), 0644); err != nil {
		t.Fatalf("doc: %v", err)
	}
	if err := os.WriteFile(filepath.Join(folder, "world_lore.json"), []byte(`{"entries":[{"id":"lb_owner","ownerId":"somebody-else"}]}`), 0644); err != nil {
		t.Fatalf("lore: %v", err)
	}
	req := Request{Category: "worlds", DestinationLibraryID: destinationID(t, registry, library.Worlds), LegacyRoots: []string{source}}
	plan := mustPreview(t, service, req)
	if plan.Quarantined != 1 || plan.Items[0].Reason != "unmapped_owner" {
		t.Fatalf("plan = %+v", plan)
	}
}

func TestSourceDestinationOverlapAndUnsafeDestination(t *testing.T) {
	service, _, registry, root, _ := newImportTest(t, library.Worlds)
	overlap := Request{
		Category:             "worlds",
		DestinationLibraryID: destinationID(t, registry, library.Worlds),
		LegacyRoots:          []string{root},
	}
	plan := mustPreview(t, service, overlap)
	if plan.Quarantined != 1 || plan.Items[0].Reason != "source_destination_overlap" {
		t.Fatalf("overlap plan = %+v", plan)
	}

	view := registry.List(library.Worlds)[0]
	if err := registry.SetRole(view.ID, library.RoleSave); err != nil {
		t.Fatalf("set role: %v", err)
	}
	req := browserRequest("wld_unsafe", "Unsafe")
	req.DestinationLibraryID = view.ID
	plan = mustPreview(t, service, req)
	if plan.Items[0].Reason != "destination_not_rollback_safe" {
		t.Fatalf("unsafe reason = %q", plan.Items[0].Reason)
	}
}

func TestIncompleteStagingRequiresExplicitAmbiguity(t *testing.T) {
	service, _, registry, root, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_incomplete", "Incomplete")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	req.OperationID = "imp_incomplete"
	plan := mustPreview(t, service, req)
	service.SetFaultHook(func(name string) error {
		if name == "after_stage" {
			path := filepath.Join(root, worldstore.ImportStagingDir, "imp_incomplete", "wld_incomplete", "world.json")
			if err := os.Remove(path); err != nil {
				t.Fatalf("remove staged document: %v", err)
			}
			return errors.New("crash during staging")
		}
		return nil
	})
	result, err := service.Execute(req, plan.Hash)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	service.SetFaultHook(nil)
	resumed, err := service.Resume("worlds", result.OperationID)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if resumed.Status != OperationAmbiguous {
		t.Fatalf("status = %q, want ambiguous; %+v", resumed.Status, resumed.Items)
	}
	if _, err := os.Stat(filepath.Join(root, "wld_incomplete")); !os.IsNotExist(err) {
		t.Fatalf("incomplete staging produced canonical content: %v", err)
	}
}

func TestPlanHashIsDeterministicAndSensitive(t *testing.T) {
	service, _, registry, _, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_hash", "Hash")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	one := mustPreview(t, service, req)
	two := mustPreview(t, service, req)
	if one.Hash == "" || one.Hash != two.Hash {
		t.Fatalf("hashes = %q %q", one.Hash, two.Hash)
	}
	req.BrowserRecords[0]["details"] = "changed"
	three := mustPreview(t, service, req)
	if three.Hash == one.Hash || !strings.HasPrefix(three.Hash, "") {
		t.Fatalf("changed source did not change hash: %q %q", one.Hash, three.Hash)
	}
}
