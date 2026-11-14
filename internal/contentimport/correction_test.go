package contentimport

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/library"
	"howl-chat/internal/worldstore"
)

func writeLegacyItem(t *testing.T, root, folder, document string, files map[string]string) {
	t.Helper()
	dir := filepath.Join(root, folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir %s: %v", folder, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "world.json"), []byte(document), 0644); err != nil {
		t.Fatalf("write document: %v", err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
}

func TestImportNormalizesReferencesLosslesslyAndResolvesByID(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	req := Request{
		Category:             "worlds",
		DestinationLibraryID: destinationID(t, registry, library.Worlds),
		BrowserRecords: []map[string]interface{}{{
			"id":          "wld_refs",
			"title":       "References",
			"unknownTop":  "keep",
			"scenarios":   []interface{}{map[string]interface{}{"id": "scn_1", "title": "Scenario", "custom": "keep"}, map[string]interface{}{"name": "Legacy Scenario", "legacyUnknown": "keep"}},
			"characters":  []interface{}{map[string]interface{}{"id": "char_1", "name": "Aria", "custom": "keep"}, "Old Aria"},
			"howl":        map[string]interface{}{"custom": "keep", "lorebooks": []interface{}{map[string]interface{}{"id": "lb_1", "title": "Lore", "custom": "keep"}, map[string]interface{}{"name": "Old Lore", "legacyUnknown": "keep"}}},
			"lorebooks":   []interface{}{map[string]interface{}{"id": "lb_1", "title": "Lore duplicate", "duplicateUnknown": "merged"}},
			"description": "keep",
		}},
	}
	plan := mustPreview(t, service, req)
	if plan.ReferenceIssues < 3 || len(plan.Items[0].References) < 3 {
		t.Fatalf("reference report = %+v", plan.Items[0].References)
	}
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("execute: %v %+v", err, result)
	}
	item, err := store.Load("wld_refs")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	scenarios := item.Document["scenarios"].([]interface{})
	scenario := scenarios[0].(map[string]interface{})
	if scenario["kind"] != "scenario" || scenario["id"] != "scn_1" || scenario["custom"] != "keep" {
		t.Fatalf("typed scenario = %+v", scenario)
	}
	legacyScenario := scenarios[1].(map[string]interface{})
	if legacyScenario["kind"] != "scenario" || legacyScenario["id"] != "" || legacyScenario["legacyName"] != "Legacy Scenario" {
		t.Fatalf("legacy scenario = %+v", legacyScenario)
	}
	characters := item.Document["characters"].([]interface{})
	character := characters[0].(map[string]interface{})
	if character["kind"] != "character" || character["id"] != "char_1" || character["custom"] != "keep" {
		t.Fatalf("typed character = %+v", character)
	}
	legacyCharacter := characters[1].(map[string]interface{})
	if legacyCharacter["legacyName"] != "Old Aria" || legacyCharacter["id"] != "" {
		t.Fatalf("legacy character = %+v", legacyCharacter)
	}
	if _, exists := item.Document["lorebooks"]; exists {
		t.Fatal("duplicate top-level lorebooks list was written")
	}
	howl := item.Document["howl"].(map[string]interface{})
	lorebooks := howl["lorebooks"].([]interface{})
	if len(lorebooks) != 2 {
		t.Fatalf("lorebooks = %+v", lorebooks)
	}
	lore := lorebooks[0].(map[string]interface{})
	if lore["kind"] != "lorebook" || lore["id"] != "lb_1" || lore["custom"] != "keep" || lore["duplicateUnknown"] != "merged" {
		t.Fatalf("typed lorebook = %+v", lore)
	}
	legacyLore := lorebooks[1].(map[string]interface{})
	if legacyLore["legacyName"] != "Old Lore" || legacyLore["legacyUnknown"] != "keep" {
		t.Fatalf("legacy lorebook = %+v", legacyLore)
	}
	if item.Document["unknownTop"] != "keep" || howl["custom"] != "keep" {
		t.Fatal("unknown document fields were lost")
	}
	refs := worldstore.NormalizeReferences(worldstore.RefLorebook, lorebooks)
	if len(refs) != 2 || refs[0].ID != "lb_1" || refs[1].LegacyName != "Old Lore" {
		t.Fatalf("normalized refs = %+v", refs)
	}
	resolutions := worldstore.ResolveReferences(refs, worldstore.Catalog{worldstore.RefLorebook: {"lb_1": true}})
	if !resolutions[0].Resolved || resolutions[1].Resolved || resolutions[1].Reason != "legacy_name_unresolved" {
		t.Fatalf("resolutions = %+v", resolutions)
	}

	updated, err := store.Update("wld_refs", map[string]interface{}{"title": "References Updated"}, item.Revision)
	if err != nil || updated.Status != "updated" {
		t.Fatalf("update: %v %+v", err, updated)
	}
	reloaded, err := store.Load("wld_refs")
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	reloadedScenarios := reloaded.Document["scenarios"].([]interface{})
	if len(reloadedScenarios) != 2 || reloadedScenarios[0].(map[string]interface{})["custom"] != "keep" {
		t.Fatal("lossless update changed normalized reference fields")
	}
	reloadedLorebooks := reloaded.Document["howl"].(map[string]interface{})["lorebooks"].([]interface{})
	if len(reloadedLorebooks) != 2 || reloadedLorebooks[0].(map[string]interface{})["custom"] != "keep" || reloadedLorebooks[0].(map[string]interface{})["duplicateUnknown"] != "merged" {
		t.Fatal("lossless update changed normalized references")
	}
}

func TestImportClassifiesAllReferenceEdgeCases(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	req := Request{
		Category:             "worlds",
		DestinationLibraryID: destinationID(t, registry, library.Worlds),
		BrowserRecords: []map[string]interface{}{{
			"id":    "wld_ref_edges",
			"title": "Reference Edges",
			"scenarios": []interface{}{
				map[string]interface{}{"kind": "character", "id": "char_mismatch", "custom": "mismatch-keep"},
				map[string]interface{}{"kind": "alien", "id": "alien_1", "custom": "unknown-keep"},
				map[string]interface{}{"custom": "missing-id-keep"},
				7,
				"Legacy Scenario",
			},
		}},
	}
	plan := mustPreview(t, service, req)
	if plan.ReferenceIssues != 5 || len(plan.Items[0].References) != 5 {
		t.Fatalf("reference report = %d/%d", plan.ReferenceIssues, len(plan.Items[0].References))
	}
	gotReasons := map[string]bool{}
	for _, issue := range plan.Items[0].References {
		gotReasons[issue.Reason] = true
	}
	for _, reason := range []string{"kind_mismatch", "unknown_kind", "missing_id", "malformed", "legacy_name"} {
		if !gotReasons[reason] {
			t.Fatalf("missing reason %q in %+v", reason, plan.Items[0].References)
		}
	}
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("execute: %v %+v", err, result)
	}
	item, err := store.Load("wld_ref_edges")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	scenarios := item.Document["scenarios"].([]interface{})
	if len(scenarios) != 5 {
		t.Fatalf("edge scenarios = %+v", scenarios)
	}
	mismatch := scenarios[0].(map[string]interface{})
	if mismatch["kind"] != "character" || mismatch["custom"] != "mismatch-keep" {
		t.Fatalf("kind mismatch was not preserved: %+v", mismatch)
	}
	unknown := scenarios[1].(map[string]interface{})
	if unknown["kind"] != "alien" || unknown["custom"] != "unknown-keep" {
		t.Fatalf("unknown kind was not preserved: %+v", unknown)
	}
	missingID := scenarios[2].(map[string]interface{})
	if missingID["id"] != "" || missingID["custom"] != "missing-id-keep" {
		t.Fatalf("missing-id reference was not preserved: %+v", missingID)
	}
	malformed, ok := scenarios[3].(json.Number)
	if !ok || malformed.String() != "7" {
		t.Fatalf("malformed reference was not preserved: %#v", scenarios[3])
	}
	legacy := scenarios[4].(map[string]interface{})
	if legacy["kind"] != "scenario" || legacy["legacyName"] != "Legacy Scenario" || legacy["id"] != "" {
		t.Fatalf("legacy-name reference was not typed: %+v", legacy)
	}

	updated, err := store.Update("wld_ref_edges", map[string]interface{}{"title": "Reference Edges Updated"}, item.Revision)
	if err != nil || updated.Status != "updated" {
		t.Fatalf("update: %v %+v", err, updated)
	}
	reloaded, err := store.Load("wld_ref_edges")
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	reloadedScenarios := reloaded.Document["scenarios"].([]interface{})
	if reloadedScenarios[0].(map[string]interface{})["custom"] != "mismatch-keep" || reloadedScenarios[1].(map[string]interface{})["custom"] != "unknown-keep" || reloadedScenarios[2].(map[string]interface{})["custom"] != "missing-id-keep" {
		t.Fatalf("unknown per-reference fields did not survive update: %+v", reloadedScenarios)
	}
}

func TestLegacyRootRejectsRegisteredEquivalentAndNestedPaths(t *testing.T) {
	service, _, registry, root, _ := newImportTest(t, library.Worlds)
	destination := destinationID(t, registry, library.Worlds)
	for _, source := range []string{root, root + string(filepath.Separator) + "."} {
		plan := mustPreview(t, service, Request{Category: "worlds", DestinationLibraryID: destination, LegacyRoots: []string{source}})
		if len(plan.Items) != 1 || plan.Items[0].Reason != "source_destination_overlap" {
			t.Fatalf("source %q plan = %+v", source, plan)
		}
	}
	other := t.TempDir()
	if _, err := registry.Add(library.Registration{Category: library.Worlds, Path: other, Role: library.RoleRead, Origin: library.OriginExternal, Format: library.FormatHowl}); err != nil {
		t.Fatalf("register other root: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(other, "Nested"), 0755); err != nil {
		t.Fatalf("nested dir: %v", err)
	}
	for _, source := range []string{other, filepath.Join(other, "Nested")} {
		plan := mustPreview(t, service, Request{Category: "worlds", DestinationLibraryID: destination, LegacyRoots: []string{source}})
		if len(plan.Items) != 1 || plan.Items[0].Reason != "source_registered" {
			t.Fatalf("registered source %q plan = %+v", source, plan)
		}
	}

	valid := t.TempDir()
	writeLegacyItem(t, valid, "Valid", `{"id":"wld_valid_root","title":"Valid"}`, nil)
	validPlan := mustPreview(t, service, Request{Category: "worlds", DestinationLibraryID: destination, LegacyRoots: []string{valid}})
	if validPlan.Ready != 1 || validPlan.Items[0].Reason != "" {
		t.Fatalf("valid unregistered root plan = %+v", validPlan)
	}
}

func TestMediaAdoptionRequiresIdenticalManifest(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	destination := destinationID(t, registry, library.Worlds)
	document := `{"id":"wld_media_adopt","title":"Media"}`
	firstRoot := t.TempDir()
	writeLegacyItem(t, firstRoot, "Same", document, map[string]string{
		"media.json": `{"format":"howl-media-refs-v1","items":[{"mediaId":"m1","role":"thumbnail","provider":"legacy-indexeddb","custom":"keep"}],"unknownTop":"keep"}`,
	})
	first := Request{Category: "worlds", DestinationLibraryID: destination, LegacyRoots: []string{firstRoot}}
	plan := mustPreview(t, service, first)
	if plan.Ready != 1 {
		t.Fatalf("first plan = %+v", plan)
	}
	result, err := service.Execute(first, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("first execute: %v %+v", err, result)
	}
	second := mustPreview(t, service, first)
	if second.Adoptable != 1 {
		t.Fatalf("identical media plan = %+v", second)
	}
	item, err := store.Load("wld_media_adopt")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(item.Root, item.ID, "media.json"))
	if err != nil {
		t.Fatalf("read media: %v", err)
	}
	var media map[string]interface{}
	if err := json.Unmarshal(data, &media); err != nil || media["unknownTop"] != "keep" {
		t.Fatalf("media unknown fields = %v err=%v", media, err)
	}
	items := media["items"].([]interface{})
	if items[0].(map[string]interface{})["custom"] != "keep" {
		t.Fatalf("per-item media unknown field was lost: %+v", items[0])
	}

	differentRoot := t.TempDir()
	writeLegacyItem(t, differentRoot, "Different", document, map[string]string{
		"media.json": `{"format":"howl-media-refs-v1","items":[{"mediaId":"m2","role":"thumbnail"}]}`,
	})
	different := mustPreview(t, service, Request{Category: "worlds", DestinationLibraryID: destination, LegacyRoots: []string{differentRoot}})
	if different.Quarantined != 1 || different.Items[0].Reason != "media_conflict" {
		t.Fatalf("different media plan = %+v", different)
	}

	noMediaRoot := t.TempDir()
	writeLegacyItem(t, noMediaRoot, "NoMedia", document, nil)
	noMedia := mustPreview(t, service, Request{Category: "worlds", DestinationLibraryID: destination, LegacyRoots: []string{noMediaRoot}})
	if noMedia.Quarantined != 1 || noMedia.Items[0].Reason != "media_conflict" {
		t.Fatalf("missing media plan = %+v", noMedia)
	}
}

func TestLegacyMediaEnvelopeAdoptsWithUnknownFields(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	destination := destinationID(t, registry, library.Worlds)
	document := `{"id":"wld_legacy_media","title":"Legacy Media"}`
	firstRoot := t.TempDir()
	writeLegacyItem(t, firstRoot, "Same", document, map[string]string{
		"media.json": `{"images":[{"id":"img_1","stored":true,"custom":"keep"}],"unknownTop":"keep"}`,
	})
	req := Request{Category: "worlds", DestinationLibraryID: destination, LegacyRoots: []string{firstRoot}}
	plan := mustPreview(t, service, req)
	if plan.Ready != 1 || plan.Items[0].Reason != "" {
		t.Fatalf("legacy media plan = %+v", plan)
	}
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("legacy media execute: %v %+v", err, result)
	}
	media, exists, err := store.LoadMedia("wld_legacy_media")
	if err != nil || !exists {
		t.Fatalf("load media: %v exists=%v", err, exists)
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(media, &doc); err != nil {
		t.Fatalf("decode media: %v", err)
	}
	images := doc["images"].([]interface{})
	if doc["unknownTop"] != "keep" || images[0].(map[string]interface{})["custom"] != "keep" {
		t.Fatalf("legacy media unknown fields = %+v", doc)
	}
	second := mustPreview(t, service, req)
	if second.Adoptable != 1 {
		t.Fatalf("legacy media adoption = %+v", second)
	}

	differentRoot := t.TempDir()
	writeLegacyItem(t, differentRoot, "Different", document, map[string]string{
		"media.json": `{"images":[{"id":"img_2","stored":true}]}`,
	})
	different := mustPreview(t, service, Request{Category: "worlds", DestinationLibraryID: destination, LegacyRoots: []string{differentRoot}})
	if different.Quarantined != 1 || different.Items[0].Reason != "media_conflict" {
		t.Fatalf("different legacy media plan = %+v", different)
	}
}

func TestEmptyAndWhitespaceOwnerIDsAreStamped(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	source := t.TempDir()
	lore := `{"entries":[
		{"id":"lb_absent","content":"absent"},
		{"id":"lb_empty","ownerId":"","content":"empty"},
		{"id":"lb_ws","ownerId":"   ","content":"whitespace"},
		{"id":"lb_match","ownerId":"wld_owner","content":"matching"}
	]}`
	writeLegacyItem(t, source, "Owners", `{"id":"wld_owner","title":"Owners"}`, map[string]string{"world_lore.json": lore})
	req := Request{Category: "worlds", DestinationLibraryID: destinationID(t, registry, library.Worlds), LegacyRoots: []string{source}}
	plan := mustPreview(t, service, req)
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompleted {
		t.Fatalf("execute: %v %+v", err, result)
	}
	loaded, err := store.LoadLore("wld_owner", "world_lore.json")
	if err != nil || len(loaded.Entries) != 4 {
		t.Fatalf("lore: %+v err=%v", loaded, err)
	}
	for _, raw := range loaded.Entries {
		var entry map[string]interface{}
		if err := json.Unmarshal(raw, &entry); err != nil {
			t.Fatalf("entry: %v", err)
		}
		if entry["ownerId"] != "wld_owner" {
			t.Fatalf("owner was not stamped: %+v", entry)
		}
	}
}

func TestResultConflictsExcludeNonConflictQuarantine(t *testing.T) {
	service, _, registry, root, _ := newImportTest(t, library.Worlds)
	source := t.TempDir()
	writeLegacyItem(t, source, "Unsupported", `{"id":"wld_unsupported_result","title":"Unsupported"}`, map[string]string{"notes.txt": "keep"})
	req := Request{Category: "worlds", DestinationLibraryID: destinationID(t, registry, library.Worlds), LegacyRoots: []string{source}}
	plan := mustPreview(t, service, req)
	if plan.Quarantined != 1 || plan.Unsupported != 1 || plan.Conflicted != 0 {
		t.Fatalf("plan = %+v", plan)
	}
	result, err := service.Execute(req, plan.Hash)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result.Status != OperationCompletedPartial || result.Plan.Quarantined != 1 || result.Plan.Unsupported != 1 || result.Plan.Conflicted != 0 {
		t.Fatalf("result = %+v", result)
	}
	if entries, _ := os.ReadDir(root); len(entries) != 0 {
		t.Fatalf("unsupported item wrote content: %v", entries)
	}
}

func TestMediaManifestValidationAndUnknownExtensions(t *testing.T) {
	service, _, registry, root, _ := newImportTest(t, library.Worlds)
	source := t.TempDir()
	writeLegacyItem(t, source, "Bad", `{"id":"wld_bad_media","title":"Bad"}`, map[string]string{
		"media.json": `{"format":"howl-media-refs-v1"}`,
	})
	writeLegacyItem(t, source, "Unsupported", `{"id":"wld_unsupported_media","title":"Unsupported"}`, map[string]string{
		"media.json": `{"format":"howl-media-refs-v1","items":[{"mediaId":"m1","role":"hero"}]}`,
	})
	writeLegacyItem(t, source, "Inline", `{"id":"wld_inline_media","title":"Inline"}`, map[string]string{
		"media.json": `{"format":"howl-media-refs-v1","items":[],"dataUrl":"data:image/png;base64,AAAA"}`,
	})
	writeLegacyItem(t, source, "Extended", `{"id":"wld_ext_media","title":"Extended"}`, map[string]string{
		"media.json": `{"format":"howl-media-refs-v1","items":[{"mediaId":"m1","role":"gallery","provider":"legacy-indexeddb","custom":"keep"}],"unknownTop":"keep"}`,
	})
	req := Request{Category: "worlds", DestinationLibraryID: destinationID(t, registry, library.Worlds), LegacyRoots: []string{source}}
	plan := mustPreview(t, service, req)
	reasons := map[string]string{}
	for _, item := range plan.Items {
		reasons[item.SourceID] = item.Reason
	}
	if reasons["wld_bad_media"] != "malformed_media_manifest" {
		t.Fatalf("bad media reason = %q plan=%+v", reasons["wld_bad_media"], plan)
	}
	if reasons["wld_unsupported_media"] != "unsupported_media_manifest" {
		t.Fatalf("unsupported media reason = %q plan=%+v", reasons["wld_unsupported_media"], plan)
	}
	if reasons["wld_inline_media"] != "inline_media_unmigrated" {
		t.Fatalf("inline media reason = %q plan=%+v", reasons["wld_inline_media"], plan)
	}
	if reasons["wld_ext_media"] != "" {
		t.Fatalf("extended media reason = %q plan=%+v", reasons["wld_ext_media"], plan)
	}
	result, err := service.Execute(req, plan.Hash)
	if err != nil || result.Status != OperationCompletedPartial {
		t.Fatalf("execute: %v %+v", err, result)
	}
	if result.Plan.Quarantined != 3 || result.Plan.Unsupported != 1 || result.Plan.Conflicted != 0 {
		t.Fatalf("result counts = %+v", result.Plan)
	}
	if entries, _ := os.ReadDir(root); len(entries) != 1 {
		t.Fatalf("only extended media should commit, got %v", entries)
	}
}

func TestCleanupFailureRemainsVisibleAndRollbackLeavesStaging(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	registration := registry.List(library.Worlds)[0]
	staged, err := store.StageImport(worldstore.ImportItem{
		ID:                   "wld_cleanup_fail",
		DestinationLibraryID: registration.ID,
		Document:             map[string]interface{}{"id": "wld_cleanup_fail", "title": "Cleanup"},
		Origin:               worldstore.Origin{Kind: "imported"},
		PlanHash:             "cleanup-plan",
	}, "imp_cleanup_fail")
	if err != nil {
		t.Fatalf("stage: %v", err)
	}
	if err := os.WriteFile(filepath.Join(staged.StagingPath, "unexpected.txt"), []byte("do-not-delete"), 0644); err != nil {
		t.Fatalf("unexpected file: %v", err)
	}
	operation := Operation{
		Version:              1,
		Category:             "worlds",
		OperationID:          "imp_cleanup_fail",
		PlanHash:             "cleanup-plan",
		DestinationLibraryID: registration.ID,
		DestinationRoot:      registration.Path,
		Status:               OperationCompletedPartial,
		StartedAt:            nowRFC3339(),
		UpdatedAt:            nowRFC3339(),
		Items: []JournalItem{{
			ItemID:    staged.ItemID,
			LibraryID: staged.LibraryID,
			Root:      staged.Root,
			PlanHash:  staged.PlanHash,
			State:     StateQuarantined,
			Reason:    "write_failed",
		}},
	}
	if err := service.writeOperation(operation); err != nil {
		t.Fatalf("write operation: %v", err)
	}
	if err := service.writeMarker(operation); err != nil {
		t.Fatalf("write completion marker: %v", err)
	}
	rollback, err := service.Rollback("worlds", operation.OperationID)
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if rollback.Status != "rollback_conflict" {
		t.Fatalf("rollback status = %q result=%+v", rollback.Status, rollback)
	}
	if _, err := os.Stat(staged.StagingPath); err != nil {
		t.Fatalf("staging was deleted despite unexpected file: %v", err)
	}
	if data, err := os.ReadFile(filepath.Join(staged.StagingPath, "unexpected.txt")); err != nil || string(data) != "do-not-delete" {
		t.Fatalf("unexpected file was changed: %q err=%v", data, err)
	}
}

func TestCleanupFailureAggregatesAsQuarantine(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	registration := registry.List(library.Worlds)[0]
	stage := func(id string) worldstore.StagedImport {
		t.Helper()
		staged, err := store.StageImport(worldstore.ImportItem{
			ID:                   id,
			DestinationLibraryID: registration.ID,
			Document:             map[string]interface{}{"id": id, "title": id},
			Origin:               worldstore.Origin{Kind: "imported"},
			PlanHash:             "cleanup-plan",
		}, "imp_cleanup_mixed")
		if err != nil {
			t.Fatalf("stage %s: %v", id, err)
		}
		return staged
	}
	clean := stage("wld_cleanup_clean")
	failed := stage("wld_cleanup_failed")
	if err := os.WriteFile(filepath.Join(failed.StagingPath, "unexpected.txt"), []byte("do-not-delete"), 0644); err != nil {
		t.Fatalf("unexpected file: %v", err)
	}
	operation := Operation{
		Version:              1,
		Category:             "worlds",
		OperationID:          "imp_cleanup_mixed",
		PlanHash:             "cleanup-plan",
		DestinationLibraryID: registration.ID,
		DestinationRoot:      registration.Path,
		Status:               OperationExecuting,
		StartedAt:            nowRFC3339(),
		UpdatedAt:            nowRFC3339(),
		Items: []JournalItem{
			{ItemID: clean.ItemID, LibraryID: clean.LibraryID, Root: clean.Root, PlanHash: clean.PlanHash, State: StateQuarantined, Reason: "write_failed"},
			{ItemID: failed.ItemID, LibraryID: failed.LibraryID, Root: failed.Root, PlanHash: failed.PlanHash, State: StateQuarantined, Reason: "write_failed"},
		},
	}
	if err := service.writeOperation(operation); err != nil {
		t.Fatalf("write operation: %v", err)
	}

	result, err := service.Resume("worlds", operation.OperationID)
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if result.Status != OperationAmbiguous {
		t.Fatalf("status = %q, want %q", result.Status, OperationAmbiguous)
	}
	if result.Plan.Quarantined != 2 || result.Plan.Conflicted != 0 || result.Plan.Unsupported != 0 {
		t.Fatalf("aggregate counts = %+v", result.Plan)
	}
	states := map[string]string{}
	for _, item := range result.Items {
		states[item.ItemID] = item.State
	}
	if states[clean.ItemID] != StateCleanupDone || states[failed.ItemID] != StateCleanupFailed {
		t.Fatalf("item states = %+v", states)
	}
	if len(result.Plan.Items) != 2 || result.Plan.Items[0].Action != ActionQuarantine || result.Plan.Items[1].Action != ActionQuarantine {
		t.Fatalf("item actions = %+v", result.Plan.Items)
	}
	if result.Plan.Items[0].Reason != "write_failed" || result.Plan.Items[1].Reason != "cleanup_failed" {
		t.Fatalf("item reasons = %+v", result.Plan.Items)
	}
}

func TestPreparedRollbackIsNotAggregatedAsQuarantine(t *testing.T) {
	service, store, registry, _, _ := newImportTest(t, library.Worlds)
	registration := registry.List(library.Worlds)[0]
	staged, err := store.StageImport(worldstore.ImportItem{
		ID:                   "wld_prepared_rollback",
		DestinationLibraryID: registration.ID,
		Document:             map[string]interface{}{"id": "wld_prepared_rollback", "title": "Prepared"},
		Origin:               worldstore.Origin{Kind: "imported"},
		PlanHash:             "prepared-rollback-plan",
	}, "imp_prepared_rollback")
	if err != nil {
		t.Fatalf("stage: %v", err)
	}
	operation := Operation{
		Version:              1,
		Category:             "worlds",
		OperationID:          "imp_prepared_rollback",
		PlanHash:             "prepared-rollback-plan",
		DestinationLibraryID: registration.ID,
		DestinationRoot:      registration.Path,
		Status:               OperationCompletedPartial,
		StartedAt:            nowRFC3339(),
		UpdatedAt:            nowRFC3339(),
		Items: []JournalItem{{
			ItemID:       staged.ItemID,
			LibraryID:    staged.LibraryID,
			Root:         staged.Root,
			FinalPath:    staged.FinalPath,
			StagingPath:  staged.StagingPath,
			OwnerPath:    staged.OwnerPath,
			Revision:     staged.Revision,
			PlanHash:     staged.PlanHash,
			State:        StatePrepared,
			DocumentHash: staged.DocumentHash,
		}},
	}
	if err := service.writeOperation(operation); err != nil {
		t.Fatalf("write operation: %v", err)
	}
	if err := service.writeMarker(operation); err != nil {
		t.Fatalf("write marker: %v", err)
	}

	rollback, err := service.Rollback("worlds", operation.OperationID)
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if rollback.Status != OperationRolledBack {
		t.Fatalf("rollback status = %q result=%+v", rollback.Status, rollback)
	}
	if len(rollback.Items) != 1 || rollback.Items[0].State != StateRolledBack {
		t.Fatalf("rollback items = %+v", rollback.Items)
	}

	result, err := service.Resume("worlds", operation.OperationID)
	if err != nil {
		t.Fatalf("resume after rollback: %v", err)
	}
	if result.Plan.Quarantined != 0 {
		t.Fatalf("rolled-back item counted as quarantined: %+v", result.Plan)
	}
}

func TestResultHasNoStringlyConflictForInlineMedia(t *testing.T) {
	service, _, registry, _, _ := newImportTest(t, library.Worlds)
	req := browserRequest("wld_inline_result", "Inline Result")
	req.BrowserRecords[0]["image"] = map[string]interface{}{"data": strings.Repeat("A", 128)}
	req.DestinationLibraryID = destinationID(t, registry, library.Worlds)
	plan := mustPreview(t, service, req)
	if plan.Quarantined != 1 || plan.Conflicted != 0 {
		t.Fatalf("plan = %+v", plan)
	}
	result, err := service.Execute(req, plan.Hash)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if result.Plan.Conflicted != 0 || result.Plan.Quarantined != 1 {
		t.Fatalf("result = %+v", result)
	}
}
