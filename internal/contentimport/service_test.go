package contentimport

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"howl-chat/internal/library"
	"howl-chat/internal/worldstore"
)

func newImportTest(t *testing.T, category library.Category) (*Service, *worldstore.Store, *library.Registry, string, string) {
	t.Helper()
	registry, err := library.New(t.TempDir())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	root := t.TempDir()
	_, err = registry.Add(library.Registration{
		Category: category,
		Path:     root,
		Role:     library.RoleSave | library.RoleWrite,
		Origin:   library.OriginHowl,
		Format:   library.FormatHowl,
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	config := t.TempDir()
	service, err := New(registry, category, config)
	if err != nil {
		t.Fatalf("service: %v", err)
	}
	store, err := worldstore.New(registry, category)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	return service, store, registry, root, config
}

func browserRequest(id, title string) Request {
	return Request{
		Category:             "worlds",
		DestinationLibraryID: "",
		BrowserRecords: []map[string]interface{}{{
			"id":          id,
			"title":       title,
			"details":     "browser details",
			"legacyNotes": "must survive",
			"howl": map[string]interface{}{
				"notes":      "nested note",
				"customFlag": "nested value",
			},
			"characters": []interface{}{map[string]interface{}{
				"id": "char_1", "name": "Aria", "custom": "character field",
			}},
		}},
	}
}

func destinationID(t *testing.T, registry *library.Registry, category library.Category) string {
	t.Helper()
	for _, view := range registry.List(category) {
		if view.Role.CanSave() {
			return view.ID
		}
	}
	t.Fatal("no destination")
	return ""
}

func mustPreview(t *testing.T, service *Service, req Request) Plan {
	t.Helper()
	plan, err := service.Preview(req)
	if err != nil {
		t.Fatalf("preview: %v", err)
	}
	return plan
}

func TestBrowserPreviewWritesNothingAndImportPreservesUnknowns(t *testing.T) {
	service, store, registry, root, config := newImportTest(t, library.Worlds)
	req := browserRequest("wld_browser", "Ashfall")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	beforeRecord, err := normalizeMap(req.BrowserRecords[0])
	if err != nil {
		t.Fatalf("capture browser record: %v", err)
	}

	plan := mustPreview(t, service, req)
	if plan.Ready != 1 || plan.Writes != 1 {
		t.Fatalf("plan = %+v", plan)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read destination: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("preview created destination content: %v", entries)
	}
	if _, err := os.Stat(filepath.Join(config, "worlds-import-journal.json")); !os.IsNotExist(err) {
		t.Fatalf("preview created a journal: %v", err)
	}
	afterRecord, err := normalizeMap(req.BrowserRecords[0])
	if err != nil {
		t.Fatalf("re-read browser record: %v", err)
	}
	if !reflect.DeepEqual(beforeRecord, afterRecord) {
		t.Fatalf("browser record changed during preview: before=%v after=%v", beforeRecord, afterRecord)
	}

	result, err := service.Execute(req, plan.Hash)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result.Status != OperationCompleted {
		t.Fatalf("status = %q, result=%+v", result.Status, result)
	}
	item, err := store.Load("wld_browser")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if item.Document["legacyNotes"] != "must survive" {
		t.Fatal("unknown top-level field was lost")
	}
	howl := item.Document["howl"].(map[string]interface{})
	if howl["customFlag"] != "nested value" {
		t.Fatal("unknown nested field was lost")
	}
	chars := item.Document["characters"].([]interface{})
	if chars[0].(map[string]interface{})["custom"] != "character field" {
		t.Fatal("unknown nested character field was lost")
	}
	if item.Origin.Kind != "imported" {
		t.Fatalf("origin = %+v", item.Origin)
	}
	if _, err := os.Stat(filepath.Join(root, worldstore.ImportStagingDir)); !os.IsNotExist(err) {
		t.Fatalf("staging remained after successful import: %v", err)
	}
}

func TestInlineMediaIsQuarantinedNotDropped(t *testing.T) {
	service, _, registry, root, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_media", "Media World")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	req.BrowserRecords[0]["images"] = []interface{}{map[string]interface{}{
		"id": "img_1", "name": "cover", "data": "data:image/png;base64,AAAA",
	}}
	plan := mustPreview(t, service, req)
	if plan.Quarantined != 1 || plan.Items[0].Reason != "inline_media_unmigrated" || plan.Writes != 0 {
		t.Fatalf("plan = %+v", plan)
	}
	result, err := service.Execute(req, plan.Hash)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result.Status != OperationCompletedPartial {
		t.Fatalf("status = %q", result.Status)
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 0 {
		t.Fatalf("quarantined item wrote content: %v", entries)
	}
}

func TestDuplicateIdentityAndNameAreQuarantined(t *testing.T) {
	service, _, registry, _, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_dup", "Same Identity")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	req.BrowserRecords = append(req.BrowserRecords, map[string]interface{}{"id": "wld_dup", "title": "Other"})
	plan := mustPreview(t, service, req)
	if plan.Quarantined != 2 {
		t.Fatalf("plan = %+v", plan)
	}
	for _, item := range plan.Items {
		if item.Reason != "duplicate_identity" {
			t.Fatalf("reason = %q", item.Reason)
		}
	}

	req2 := browserRequest("wld_name_a", "Same Name")
	req2.DestinationLibraryID = req.DestinationLibraryID
	req2.BrowserRecords = append(req2.BrowserRecords, map[string]interface{}{"id": "wld_name_b", "title": "Same Name"})
	plan2 := mustPreview(t, service, req2)
	if plan2.Quarantined != 2 {
		t.Fatalf("name plan = %+v", plan2)
	}
	for _, item := range plan2.Items {
		if item.Reason != "duplicate_name" {
			t.Fatalf("name reason = %q", item.Reason)
		}
	}
}

func TestExecuteRejectsStaleSourceAndDestination(t *testing.T) {
	service, store, registry, root, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_stale", "Stale")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	req.BrowserRecords[0]["title"] = "Changed After Preview"
	result, err := service.Execute(req, plan.Hash)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result.Status != "plan_stale" {
		t.Fatalf("status = %q", result.Status)
	}
	if entries, _ := os.ReadDir(root); len(entries) != 0 {
		t.Fatalf("stale source wrote content: %v", entries)
	}

	changed := browserRequest("wld_stale_dest", "Destination Race")
	changed.DestinationLibraryID = req.DestinationLibraryID
	plan = mustPreview(t, service, changed)
	if _, err := store.Save(map[string]interface{}{"id": "wld_stale_dest", "title": "Existing"}); err != nil {
		t.Fatalf("create destination item: %v", err)
	}
	result, err = service.Execute(changed, plan.Hash)
	if err != nil {
		t.Fatalf("execute after destination change: %v", err)
	}
	if result.Status != "plan_stale" {
		t.Fatalf("destination status = %q", result.Status)
	}
	item, err := store.Load("wld_stale_dest")
	if err != nil || item.Document["title"] != "Existing" {
		t.Fatalf("existing item changed: %+v err=%v", item, err)
	}
}

func TestLegacyImportPreservesLoreAndLeavesSource(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	source := t.TempDir()
	folder := filepath.Join(source, "Old Ashfall")
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	write := func(name, data string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(folder, name), []byte(data), 0644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	write("world.json", `{"id":"wld_legacy","title":"Old Ashfall","details":"legacy details","legacyDocField":{"keep":true},"createdAt":"2025-01-02T03:04:05Z"}`)
	write("world_lore.json", `{"entries":[{"id":"lb_manual","content":"manual","unknownLoreField":"keep"}]}`)
	write("automated_world_lore.json", `{"entries":[{"id":"lb_auto","content":"generated","chatId":"do-not-invent"}]}`)
	write("media.json", `{"images":[{"id":"img_ref","stored":true,"name":"cover"}]}`)
	before := treeHashes(t, source)

	req := Request{Category: "worlds", DestinationLibraryID: destinationID(t, registry, library.Worlds), LegacyRoots: []string{source}}
	plan := mustPreview(t, service, req)
	if plan.Ready != 1 {
		t.Fatalf("plan = %+v", plan)
	}
	result, err := service.Execute(req, plan.Hash)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result.Status != OperationCompleted {
		t.Fatalf("status = %q result=%+v", result.Status, result)
	}
	item, err := store.Load("wld_legacy")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if item.Document["legacyDocField"].(map[string]interface{})["keep"] != true {
		t.Fatal("legacy document unknown field was lost")
	}
	manual, err := store.LoadLore("wld_legacy", "world_lore.json")
	if err != nil || len(manual.Entries) != 1 {
		t.Fatalf("manual lore: %+v err=%v", manual, err)
	}
	var manualEntry map[string]interface{}
	_ = json.Unmarshal(manual.Entries[0], &manualEntry)
	if manualEntry["ownerId"] != "wld_legacy" || manualEntry["unknownLoreField"] != "keep" {
		t.Fatalf("manual entry = %+v", manualEntry)
	}
	automated, err := store.LoadLore("wld_legacy", "automated_world_lore.json")
	if err != nil || len(automated.Entries) != 1 {
		t.Fatalf("automated lore: %+v err=%v", automated, err)
	}
	var autoEntry map[string]interface{}
	_ = json.Unmarshal(automated.Entries[0], &autoEntry)
	if autoEntry["automated"] != true || autoEntry["originKind"] != "legacy_generated_unattributed" {
		t.Fatalf("automated entry = %+v", autoEntry)
	}
	if autoEntry["chatId"] != "do-not-invent" {
		t.Fatal("automated unknown chat field was not preserved")
	}
	if _, err := os.Stat(filepath.Join(item.Root, item.ID, "media.json")); err != nil {
		t.Fatalf("media reference missing: %v", err)
	}
	after := treeHashes(t, source)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("legacy source changed: before=%v after=%v", before, after)
	}
}

func TestLegacyUnsupportedFileQuarantinesWholeFolder(t *testing.T) {
	service, _, registry, root, _ := newImportTest(t, library.Worlds)
	source := t.TempDir()
	folder := filepath.Join(source, "Unsupported")
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(folder, "world.json"), []byte(`{"id":"wld_unsupported","title":"Unsupported"}`), 0644); err != nil {
		t.Fatalf("document: %v", err)
	}
	if err := os.WriteFile(filepath.Join(folder, "notes.txt"), []byte("do not drop silently"), 0644); err != nil {
		t.Fatalf("extra file: %v", err)
	}
	req := Request{Category: "worlds", DestinationLibraryID: destinationID(t, registry, library.Worlds), LegacyRoots: []string{source}}
	plan := mustPreview(t, service, req)
	if plan.Quarantined != 1 || plan.Items[0].Reason != "unsupported_source_file" {
		t.Fatalf("plan = %+v", plan)
	}
	if entries, _ := os.ReadDir(root); len(entries) != 0 {
		t.Fatalf("unsupported folder wrote content: %v", entries)
	}
}

func TestCrashBoundariesResumeWithoutSourceMutation(t *testing.T) {
	boundaries := []struct {
		name       string
		wantStatus string
	}{
		{"after_journal", OperationCompletedPartial},
		{"after_stage", OperationCompleted},
		{"after_prepared", OperationCompleted},
		{"after_rename", OperationCompleted},
		{"after_committed", OperationCompleted},
	}
	for _, boundary := range boundaries {
		t.Run(boundary.name, func(t *testing.T) {
			service, store, registry, _, _ := newImportTest(t, library.Worlds)
			source := t.TempDir()
			folder := filepath.Join(source, "Crash World")
			if err := os.MkdirAll(folder, 0755); err != nil {
				t.Fatalf("mkdir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(folder, "world.json"), []byte(`{"id":"wld_crash","title":"Crash World","legacyField":"keep"}`), 0644); err != nil {
				t.Fatalf("write: %v", err)
			}
			before := treeHashes(t, source)
			req := Request{Category: "worlds", DestinationLibraryID: destinationID(t, registry, library.Worlds), LegacyRoots: []string{source}}
			plan := mustPreview(t, service, req)
			service.SetFaultHook(func(name string) error {
				if name == boundary.name {
					return errors.New("injected crash at " + name)
				}
				return nil
			})
			result, err := service.Execute(req, plan.Hash)
			if err != nil {
				t.Fatalf("execute: %v", err)
			}
			if result.Status != "resume_required" {
				t.Fatalf("status = %q", result.Status)
			}
			service.SetFaultHook(nil)
			resumed, err := service.Resume("worlds", result.OperationID)
			if err != nil {
				t.Fatalf("resume: %v", err)
			}
			if resumed.Status != boundary.wantStatus {
				t.Fatalf("resume status = %q, want %q (%+v)", resumed.Status, boundary.wantStatus, resumed.Items)
			}
			if boundary.name != "after_journal" {
				item, err := store.Load("wld_crash")
				if err != nil {
					t.Fatalf("load resumed item: %v", err)
				}
				if item.Document["legacyField"] != "keep" {
					t.Fatal("unknown field lost during recovery")
				}
			}
			if after := treeHashes(t, source); !reflect.DeepEqual(before, after) {
				t.Fatalf("source changed during crash recovery: before=%v after=%v", before, after)
			}
		})
	}
}

func TestConcurrentImportsClaimIdentityOnce(t *testing.T) {
	service, _, registry, _, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_concurrent", "Concurrent")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	var wg sync.WaitGroup
	results := make([]Result, 2)
	errorsOut := make([]error, 2)
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(index int) {
			defer wg.Done()
			results[index], errorsOut[index] = service.Execute(req, plan.Hash)
		}(i)
	}
	wg.Wait()
	if errorsOut[0] != nil || errorsOut[1] != nil {
		t.Fatalf("errors = %v %v", errorsOut[0], errorsOut[1])
	}
	committed := 0
	for _, result := range results {
		for _, item := range result.Items {
			if item.State == StateCommitted {
				committed++
			}
		}
	}
	if committed != 1 {
		t.Fatalf("committed = %d, want 1; results=%+v", committed, results)
	}
}

func TestImportVersusOrdinaryCreateHasOneWinner(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_create_race", "Create Race")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	var wg sync.WaitGroup
	var saveResult worldstore.SaveResult
	var importResult Result
	var importErr error
	wg.Add(2)
	go func() {
		defer wg.Done()
		saveResult, _ = store.Save(map[string]interface{}{"id": "wld_create_race", "title": "Ordinary"})
	}()
	go func() {
		defer wg.Done()
		importResult, importErr = service.Execute(req, plan.Hash)
	}()
	wg.Wait()
	if importErr != nil {
		t.Fatalf("import: %v", importErr)
	}
	if saveResult.Status != "saved" && saveResult.Status != "collision" {
		t.Fatalf("save status = %q", saveResult.Status)
	}
	item, err := store.Load("wld_create_race")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if item.Document["title"] == "" || importResult.Status == "" {
		t.Fatalf("race left no explicit outcome: item=%+v import=%+v", item, importResult)
	}
}

func TestAuthorityChangeAfterStagingStopsResume(t *testing.T) {
	service, _, registry, root, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_authority", "Authority")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	service.SetFaultHook(func(name string) error {
		if name == "after_stage" {
			other := t.TempDir()
			registration, err := registry.Add(library.Registration{
				Category: library.Worlds, Path: other, Role: library.RoleRead,
				Origin: library.OriginExternal, Format: library.FormatUnknown,
			})
			if err != nil {
				t.Fatalf("add sibling: %v", err)
			}
			_ = registration
			if err := os.RemoveAll(other); err != nil {
				t.Fatalf("remove sibling: %v", err)
			}
			return errors.New("authority changed")
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
	if _, err := os.Stat(filepath.Join(root, "wld_authority")); !os.IsNotExist(err) {
		t.Fatalf("resume wrote despite unavailable authority: %v", err)
	}
}

func TestRollbackRemovesCleanImportAndRefusesChangedLore(t *testing.T) {
	service, _, registry, root, _ := newImportTest(t, library.Worlds)
	source := t.TempDir()
	folder := filepath.Join(source, "Rollback World")
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(folder, "world.json"), []byte(`{"id":"wld_rollback","title":"Rollback"}`), 0644); err != nil {
		t.Fatalf("write doc: %v", err)
	}
	if err := os.WriteFile(filepath.Join(folder, "world_lore.json"), []byte(`{"entries":[{"id":"lb_rb","content":"original"}]}`), 0644); err != nil {
		t.Fatalf("write lore: %v", err)
	}
	before := treeHashes(t, source)
	req := Request{Category: "worlds", DestinationLibraryID: destinationID(t, registry, library.Worlds), LegacyRoots: []string{source}}
	plan := mustPreview(t, service, req)
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("execute: %v %+v", err, result)
	}
	lorePath := filepath.Join(root, "wld_rollback", "world_lore.json")
	if err := os.WriteFile(lorePath, []byte(`{"entries":[{"id":"lb_rb","content":"edited after import"}]}`), 0644); err != nil {
		t.Fatalf("edit lore: %v", err)
	}
	rollback, err := service.Rollback("worlds", result.OperationID)
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if rollback.Status != "rollback_conflict" {
		t.Fatalf("status = %q, want rollback_conflict", rollback.Status)
	}
	if _, err := os.Stat(lorePath); err != nil {
		t.Fatalf("changed item was removed: %v", err)
	}
	if after := treeHashes(t, source); !reflect.DeepEqual(before, after) {
		t.Fatal("rollback changed source")
	}

	cleanService, cleanStore, cleanReg, cleanRoot, _ := newImportTest(t, library.Worlds)
	cleanReq := browserRequest("wld_clean", "Clean")
	cleanReq.DestinationLibraryID = destinationID(t, cleanReg, library.Worlds)
	cleanPlan := mustPreview(t, cleanService, cleanReq)
	cleanResult, err := cleanService.Execute(cleanReq, cleanPlan.Hash)
	if err != nil || cleanResult.Status != OperationCompleted {
		t.Fatalf("clean import: %v %+v", err, cleanResult)
	}
	cleanRollback, err := cleanService.Rollback("worlds", cleanResult.OperationID)
	if err != nil || cleanRollback.Status != OperationRolledBack {
		t.Fatalf("clean rollback: %v %+v", err, cleanRollback)
	}
	if _, err := cleanStore.Load("wld_clean"); !errors.Is(err, worldstore.ErrNotFound) {
		t.Fatalf("clean item remained: %v", err)
	}
	if entries, _ := os.ReadDir(cleanRoot); len(entries) != 0 {
		t.Fatalf("clean rollback left content: %v", entries)
	}
}

func TestAdoptedItemBlocksRollback(t *testing.T) {
	service, _, registry, root, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_adopted", "Adopted")
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	first, err := service.Execute(req, plan.Hash)
	if err != nil || first.Status != OperationCompleted {
		t.Fatalf("first: %v %+v", err, first)
	}
	secondPlan := mustPreview(t, service, req)
	if secondPlan.Adoptable != 1 {
		t.Fatalf("second plan = %+v", secondPlan)
	}
	second, err := service.Execute(req, secondPlan.Hash)
	if err != nil || second.Status != OperationCompleted {
		t.Fatalf("second: %v %+v", err, second)
	}
	rollback, err := service.Rollback("worlds", second.OperationID)
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if rollback.Status != "rollback_conflict" {
		t.Fatalf("status = %q, want rollback_conflict", rollback.Status)
	}
	if _, err := os.Stat(filepath.Join(root, "wld_adopted", "world.json")); err != nil {
		t.Fatalf("adopted item was removed: %v", err)
	}
}

func treeHashes(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		rel, _ := filepath.Rel(root, path)
		out[filepath.ToSlash(rel)] = hex.EncodeToString(sum[:])
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return out
}
