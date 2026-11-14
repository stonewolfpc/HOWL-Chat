package worldstore

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"howl-chat/internal/library"
)

func stageWorld(t *testing.T, store *Store, registration library.Registration, id, operation, plan string) StagedImport {
	t.Helper()
	staged, err := store.StageImport(ImportItem{
		ID:                   id,
		DestinationLibraryID: registration.ID,
		Document:             sampleWorld(id),
		Origin:               Origin{Kind: "imported"},
		Lore:                 map[string][]byte{"world_lore.json": []byte(`{"entries":[{"id":"lb_1","ownerId":"` + id + `"}]}`)},
		PlanHash:             plan,
	}, operation)
	if err != nil {
		t.Fatalf("stage: %v", err)
	}
	return staged
}

func TestStageIsHiddenAndCommitPublishes(t *testing.T) {
	store, _, registration, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	staged := stageWorld(t, store, registration, "wld_stage", "imp_stage", "plan-1")

	if _, err := os.Stat(filepath.Join(root, "wld_stage")); !os.IsNotExist(err) {
		t.Fatalf("staging created a canonical item path: %v", err)
	}
	search, err := store.Lookup("wld_stage")
	if err != nil || search.Found || !search.Complete {
		t.Fatalf("staging was discoverable: %+v err=%v", search, err)
	}
	if _, err := os.Stat(staged.OwnerPath); err != nil {
		t.Fatalf("owner marker missing: %v", err)
	}

	result, err := store.CommitImport(staged)
	if err != nil || result.Status != "saved" {
		t.Fatalf("commit: status=%q err=%v", result.Status, err)
	}
	item, err := store.Load("wld_stage")
	if err != nil {
		t.Fatalf("load committed: %v", err)
	}
	if item.Document["title"] != "Ashfall" || item.Origin.Kind != "imported" {
		t.Fatalf("imported item = %+v", item)
	}
	if _, err := os.Stat(filepath.Join(root, ImportStagingDir, "imp_stage")); !os.IsNotExist(err) {
		t.Fatalf("staging parent remained after commit: %v", err)
	}
}

func TestCommitNeverOverwritesExistingFinal(t *testing.T) {
	store, _, registration, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	staged := stageWorld(t, store, registration, "wld_collision", "imp_collision", "plan-2")
	if err := os.MkdirAll(filepath.Join(root, "wld_collision"), 0755); err != nil {
		t.Fatalf("create final: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "wld_collision", "world.json"), []byte("existing"), 0644); err != nil {
		t.Fatalf("write final: %v", err)
	}
	result, err := store.CommitImport(staged)
	if err != nil {
		t.Fatalf("commit: %v", err)
	}
	if result.Status != "collision" {
		t.Fatalf("status = %q, want collision", result.Status)
	}
	data, _ := os.ReadFile(filepath.Join(root, "wld_collision", "world.json"))
	if string(data) != "existing" {
		t.Fatal("existing canonical content was overwritten")
	}
}

func TestCommitRefusesUnavailableSiblingAfterStage(t *testing.T) {
	store, reg, registration, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	staged := stageWorld(t, store, registration, "wld_unavailable", "imp_unavailable", "plan-3")
	addUnavailableRoot(t, reg)
	result, err := store.CommitImport(staged)
	if err != nil {
		t.Fatalf("commit: %v", err)
	}
	if result.Status != "destination_unavailable" {
		t.Fatalf("status = %q, want destination_unavailable", result.Status)
	}
	if _, err := os.Stat(filepath.Join(root, "wld_unavailable")); !os.IsNotExist(err) {
		t.Fatalf("commit wrote despite unavailable sibling: %v", err)
	}
}

func TestRemoveImportedRequiresEveryHash(t *testing.T) {
	store, _, registration, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	staged := stageWorld(t, store, registration, "wld_remove", "imp_remove", "plan-4")
	result, err := store.CommitImport(staged)
	if err != nil || result.Status != "saved" {
		t.Fatalf("commit: %v %v", result, err)
	}
	lorePath := filepath.Join(root, "wld_remove", "world_lore.json")
	if err := os.WriteFile(lorePath, []byte(`{"entries":[{"id":"lb_1","ownerId":"changed"}]}`), 0644); err != nil {
		t.Fatalf("edit lore: %v", err)
	}
	removed, err := store.RemoveImported("wld_remove", staged.Revision, map[string]string{
		"world.json":      staged.DocumentHash,
		"world_lore.json": staged.LoreHashes["world_lore.json"],
	})
	if err != nil {
		t.Fatalf("remove: %v", err)
	}
	if removed.Status != "rollback_conflict" {
		t.Fatalf("status = %q, want rollback_conflict", removed.Status)
	}
	if _, err := os.Stat(lorePath); err != nil {
		t.Fatalf("changed lore was removed: %v", err)
	}
}

func TestAbortImportRefusesUnexpectedStagingContents(t *testing.T) {
	store, _, registration, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	staged := stageWorld(t, store, registration, "wld_abort_unexpected", "imp_abort_unexpected", "plan-abort")
	if err := os.WriteFile(filepath.Join(staged.StagingPath, "unexpected.txt"), []byte("do-not-delete"), 0644); err != nil {
		t.Fatalf("write unexpected: %v", err)
	}
	if err := store.AbortImport(staged); err == nil {
		t.Fatal("abort accepted unexpected staging contents")
	}
	if data, err := os.ReadFile(filepath.Join(staged.StagingPath, "unexpected.txt")); err != nil || string(data) != "do-not-delete" {
		t.Fatalf("unexpected staging file changed: %q err=%v", data, err)
	}
	if _, err := os.Stat(filepath.Join(root, "wld_abort_unexpected")); !os.IsNotExist(err) {
		t.Fatalf("abort created a canonical item: %v", err)
	}
}

func TestAbortAndCommitSerializeToOneOutcome(t *testing.T) {
	store, _, registration, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	staged := stageWorld(t, store, registration, "wld_abort_commit", "imp_abort_commit", "plan-abort-commit")
	var wg sync.WaitGroup
	var abortErr error
	var commitResult SaveResult
	var commitErr error
	wg.Add(2)
	go func() { defer wg.Done(); abortErr = store.AbortImport(staged) }()
	go func() { defer wg.Done(); commitResult, commitErr = store.CommitImport(staged) }()
	wg.Wait()
	if abortErr != nil && commitErr != nil {
		t.Fatalf("both abort and commit failed: %v / %v", abortErr, commitErr)
	}
	if commitErr == nil && commitResult.Status == "saved" {
		if _, err := os.Stat(filepath.Join(root, "wld_abort_commit")); err != nil {
			t.Fatalf("committed item missing: %v", err)
		}
		if _, err := os.Stat(staged.StagingPath); !os.IsNotExist(err) {
			t.Fatalf("staging remained after commit: %v", err)
		}
		return
	}
	if abortErr != nil {
		t.Fatalf("abort failed while commit did not save: abort=%v commit=%+v err=%v", abortErr, commitResult, commitErr)
	}
	if _, err := os.Stat(filepath.Join(root, "wld_abort_commit")); !os.IsNotExist(err) {
		t.Fatalf("aborted import produced a canonical item: %v", err)
	}
}

func TestConcurrentCommitsHaveOneWinner(t *testing.T) {
	store, _, registration, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	one := stageWorld(t, store, registration, "wld_race_import", "imp_race_a", "plan-5")
	two := stageWorld(t, store, registration, "wld_race_import", "imp_race_b", "plan-5")
	var wg sync.WaitGroup
	results := make([]SaveResult, 2)
	wg.Add(2)
	go func() { defer wg.Done(); results[0], _ = store.CommitImport(one) }()
	go func() { defer wg.Done(); results[1], _ = store.CommitImport(two) }()
	wg.Wait()
	saved, collision := 0, 0
	for _, result := range results {
		switch result.Status {
		case "saved":
			saved++
		case "collision":
			collision++
		}
	}
	if saved != 1 || collision != 1 {
		t.Fatalf("results = %+v, want exactly one saved and one collision", results)
	}
}
