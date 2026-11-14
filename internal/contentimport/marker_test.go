package contentimport

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"howl-chat/internal/library"
)

func TestCompletedMarkerResumeIsNoOp(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_marker", "Marker")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("execute: %v %+v", err, result)
	}
	marker, err := service.readMarker(result.OperationID)
	if err != nil {
		t.Fatalf("read marker: %v", err)
	}
	if marker.OperationID != result.OperationID || marker.Status != OperationCompleted || len(marker.ItemIDs) != 1 {
		t.Fatalf("marker = %+v", marker)
	}
	resumed, err := service.Resume("worlds", result.OperationID)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if resumed.Status != OperationCompleted {
		t.Fatalf("resume status = %q", resumed.Status)
	}
	item, err := store.Load("wld_marker")
	if err != nil || item.Document["title"] != "Marker" {
		t.Fatalf("terminal resume changed item: %+v err=%v", item, err)
	}
	if _, err := os.Stat(markerPathForTest(service, result.OperationID)); err != nil {
		t.Fatalf("completion marker missing: %v", err)
	}
}

func TestMissingMarkerBlocksCompletedResume(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_marker_missing", "Missing Marker")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("execute: %v %+v", err, result)
	}
	if err := os.Remove(markerPathForTest(service, result.OperationID)); err != nil {
		t.Fatalf("remove marker: %v", err)
	}

	resumed, err := service.Resume("worlds", result.OperationID)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if resumed.Status != OperationIntegrityFailure {
		t.Fatalf("resume status = %q, want %q", resumed.Status, OperationIntegrityFailure)
	}
	if item, err := store.Load("wld_marker_missing"); err != nil || item.Document["title"] != "Missing Marker" {
		t.Fatalf("canonical content was changed: %+v err=%v", item, err)
	}
}

func TestMalformedMarkerBlocksCompletedResume(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_marker_malformed", "Malformed Marker")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("execute: %v %+v", err, result)
	}
	if err := os.WriteFile(markerPathForTest(service, result.OperationID), []byte(`{"version":2,"operation_id":`), 0644); err != nil {
		t.Fatalf("write malformed marker: %v", err)
	}

	resumed, err := service.Resume("worlds", result.OperationID)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if resumed.Status != OperationIntegrityFailure {
		t.Fatalf("resume status = %q, want %q", resumed.Status, OperationIntegrityFailure)
	}
	if item, err := store.Load("wld_marker_malformed"); err != nil || item.Document["title"] != "Malformed Marker" {
		t.Fatalf("canonical content was changed: %+v err=%v", item, err)
	}
}

func TestContradictoryMarkerBlocksCompletedResume(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_marker_contradiction", "Contradiction")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("execute: %v %+v", err, result)
	}
	marker, err := service.readMarker(result.OperationID)
	if err != nil {
		t.Fatalf("read marker: %v", err)
	}
	marker.Status = OperationCompletedPartial
	data, err := json.MarshalIndent(marker, "", "  ")
	if err != nil {
		t.Fatalf("encode marker: %v", err)
	}
	if err := writeAtomicFile(markerPathForTest(service, result.OperationID), data); err != nil {
		t.Fatalf("write contradictory marker: %v", err)
	}

	resumed, err := service.Resume("worlds", result.OperationID)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if resumed.Status != OperationIntegrityFailure {
		t.Fatalf("resume status = %q, want %q", resumed.Status, OperationIntegrityFailure)
	}
	if item, err := store.Load("wld_marker_contradiction"); err != nil || item.Document["title"] != "Contradiction" {
		t.Fatalf("canonical content was changed: %+v err=%v", item, err)
	}
}

func TestStaleMarkerBlocksResume(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_marker_stale", "Stale Marker")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("execute: %v %+v", err, result)
	}
	operation, ok := service.loadOperation(result.OperationID)
	if !ok {
		t.Fatal("operation missing")
	}
	operation.Status = OperationExecuting
	if err := service.writeOperation(operation); err != nil {
		t.Fatalf("write stale journal: %v", err)
	}

	resumed, err := service.Resume("worlds", result.OperationID)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if resumed.Status != OperationIntegrityFailure {
		t.Fatalf("resume status = %q, want %q", resumed.Status, OperationIntegrityFailure)
	}
	if item, err := store.Load("wld_marker_stale"); err != nil || item.Document["title"] != "Stale Marker" {
		t.Fatalf("canonical content was changed: %+v err=%v", item, err)
	}
}

func TestMarkerWriteRefusesForeignOwner(t *testing.T) {
	service, _, _, _, _ := newImportTest(t, library.Worlds)
	operation := Operation{
		Version:              1,
		Category:             "worlds",
		OperationID:          "imp_marker_owner",
		PlanHash:             "marker-plan",
		DestinationLibraryID: "library",
		DestinationRoot:      "root",
		Status:               OperationCompleted,
		Items:                []JournalItem{{ItemID: "wld_marker_owner"}},
	}
	foreign := markerForOperation(operation)
	foreign.OperationID = "imp_marker_foreign"
	data, err := json.MarshalIndent(foreign, "", "  ")
	if err != nil {
		t.Fatalf("encode foreign marker: %v", err)
	}
	path := markerPathForTest(service, operation.OperationID)
	if err := writeAtomicFile(path, data); err != nil {
		t.Fatalf("write foreign marker: %v", err)
	}

	err = service.writeMarker(operation)
	if !errors.Is(err, ErrMarkerConflict) {
		t.Fatalf("writeMarker error = %v, want ErrMarkerConflict", err)
	}
	got, err := service.readMarker(operation.OperationID)
	if err != nil {
		t.Fatalf("read marker after collision: %v", err)
	}
	if got.OperationID != foreign.OperationID {
		t.Fatalf("foreign marker was overwritten: %+v", got)
	}
}

func TestRollbackPartialRemovesCompletionMarker(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_marker_rollback", "Rollback Marker")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("execute: %v %+v", err, result)
	}
	item, err := store.Load("wld_marker_rollback")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	path := item.Path
	if err := os.WriteFile(path, []byte(`{"id":"wld_marker_rollback","title":"changed"}`), 0644); err != nil {
		t.Fatalf("mutate canonical item: %v", err)
	}

	rollback, err := service.Rollback("worlds", result.OperationID)
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if rollback.Status != "rollback_conflict" {
		t.Fatalf("rollback status = %q, want rollback_conflict", rollback.Status)
	}
	if _, err := os.Stat(markerPathForTest(service, result.OperationID)); !os.IsNotExist(err) {
		t.Fatalf("partial rollback left a marker: %v", err)
	}
	resumed, err := service.Resume("worlds", result.OperationID)
	if err != nil {
		t.Fatalf("resume after partial rollback: %v", err)
	}
	if resumed.Status != OperationRollbackPartial {
		t.Fatalf("resume status = %q, want %q", resumed.Status, OperationRollbackPartial)
	}
}

func markerPathForTest(service *Service, operationID string) string {
	return service.markerPath(operationID)
}
