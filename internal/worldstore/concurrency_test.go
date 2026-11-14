package worldstore

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"howl-chat/internal/library"
)

func TestConcurrentUpdatesWithSameBaseOnlyOneSucceeds(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("wld_concurrent"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	const workers = 8
	results := make([]SaveResult, workers)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			<-start
			result, updateErr := store.Update("wld_concurrent", map[string]interface{}{
				"title": fmt.Sprintf("Version %d", index),
			}, saved.Revision)
			if updateErr != nil {
				t.Errorf("worker %d: %v", index, updateErr)
				return
			}
			results[index] = result
		}(i)
	}
	close(start)
	wg.Wait()

	updated, conflicts, winner := 0, 0, -1
	for index, result := range results {
		switch result.Status {
		case "updated":
			updated++
			winner = index
		case "revision_conflict":
			conflicts++
		default:
			t.Errorf("worker %d: unexpected status %q (%s)", index, result.Status, result.Message)
		}
	}
	if updated != 1 || conflicts != workers-1 {
		t.Fatalf("updated=%d conflicts=%d, want 1 and %d", updated, conflicts, workers-1)
	}

	item, err := store.Load("wld_concurrent")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if want := fmt.Sprintf("Version %d", winner); item.Document["title"] != want {
		t.Fatalf("winning document = %v, want %q", item.Document["title"], want)
	}
	if item.Revision != results[winner].Revision {
		t.Fatalf("stored revision = %q, want winner %q", item.Revision, results[winner].Revision)
	}
	if item.RevisionMismatch {
		t.Fatal("concurrent update left a document whose revision does not match its content")
	}
}

func TestConcurrentUpdatesAcrossStoreInstances(t *testing.T) {
	storeA, reg, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	storeB, err := New(reg, library.Worlds)
	if err != nil {
		t.Fatalf("second store: %v", err)
	}
	saved, err := storeA.Save(sampleWorld("wld_two_stores"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	start := make(chan struct{})
	results := make([]SaveResult, 2)
	var wg sync.WaitGroup
	for index, candidate := range []*Store{storeA, storeB} {
		wg.Add(1)
		go func(index int, store *Store) {
			defer wg.Done()
			<-start
			result, updateErr := store.Update("wld_two_stores", map[string]interface{}{
				"title": fmt.Sprintf("Store %d", index),
			}, saved.Revision)
			if updateErr != nil {
				t.Errorf("store %d: %v", index, updateErr)
				return
			}
			results[index] = result
		}(index, candidate)
	}
	close(start)
	wg.Wait()

	updated, conflicts := 0, 0
	for _, result := range results {
		switch result.Status {
		case "updated":
			updated++
		case "revision_conflict":
			conflicts++
		}
	}
	if updated != 1 || conflicts != 1 {
		t.Fatalf("updated=%d conflicts=%d, want 1 and 1", updated, conflicts)
	}
}

func TestConcurrentCreatesDoNotOverwrite(t *testing.T) {
	store, _, _, _ := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)

	const workers = 8
	results := make([]SaveResult, workers)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			<-start
			document := sampleWorld("wld_create_race")
			document["title"] = fmt.Sprintf("Creator %d", index)
			result, saveErr := store.Save(document)
			if saveErr != nil {
				t.Errorf("worker %d: %v", index, saveErr)
				return
			}
			results[index] = result
		}(i)
	}
	close(start)
	wg.Wait()

	saved, collisions, winner := 0, 0, -1
	for index, result := range results {
		switch result.Status {
		case "saved":
			saved++
			winner = index
		case "collision":
			collisions++
		default:
			t.Errorf("worker %d: unexpected status %q", index, result.Status)
		}
	}
	if saved != 1 || collisions != workers-1 {
		t.Fatalf("saved=%d collisions=%d, want 1 and %d", saved, collisions, workers-1)
	}

	item, err := store.Load("wld_create_race")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if want := fmt.Sprintf("Creator %d", winner); item.Document["title"] != want {
		t.Fatalf("document = %v, want the winning creator %q", item.Document["title"], want)
	}
}

func TestInterruptedWriteLeavesPriorDocumentIntact(t *testing.T) {
	store, _, _, root := newTestStore(t, library.Worlds, library.RoleSave|library.RoleWrite)
	saved, err := store.Save(sampleWorld("wld_atomic"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	before, err := os.ReadFile(saved.Path)
	if err != nil {
		t.Fatalf("read before: %v", err)
	}

	// Make the atomic temp target unusable so the registry's write fails before
	// the canonical file is replaced.
	blocker := filepath.Join(root, "wld_atomic", "world.json.tmp")
	if err := os.Mkdir(blocker, 0755); err != nil {
		t.Fatalf("create blocker: %v", err)
	}

	result, err := store.Update("wld_atomic", map[string]interface{}{"title": "Should Not Land"}, saved.Revision)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if result.Status != "write_failed" {
		t.Fatalf("status = %q, want write_failed (%s)", result.Status, result.Message)
	}

	after, err := os.ReadFile(saved.Path)
	if err != nil {
		t.Fatalf("read after: %v", err)
	}
	if string(before) != string(after) {
		t.Fatal("failed write changed the canonical document")
	}

	item, err := store.Load("wld_atomic")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if item.Document["title"] != "Ashfall" {
		t.Fatalf("title = %v, want original", item.Document["title"])
	}
	if item.RevisionMismatch {
		t.Fatal("failed write left an inconsistent document")
	}
}
