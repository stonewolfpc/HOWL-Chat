package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/library"
	"howl-chat/internal/lorebookstore"
)

func newLorebookApp(t *testing.T) (*App, string) {
	t.Helper()
	root := t.TempDir()
	reg, err := library.New(t.TempDir())
	if err != nil {
		t.Fatalf("library.New: %v", err)
	}
	if _, err := reg.Add(library.Registration{
		Category: library.Lorebooks,
		Path:     root,
		Role:     library.RoleSave | library.RoleWrite,
		Origin:   library.OriginHowl,
		Format:   library.FormatHowl,
	}); err != nil {
		t.Fatalf("register root: %v", err)
	}
	app := &App{
		contentLibraries: reg,
		lorebookStore:    lorebookstore.New(reg, t.TempDir()),
	}
	return app, root
}

func TestSaveLorebookRequiresImportFirst(t *testing.T) {
	app, root := newLorebookApp(t)
	_, err := app.SaveLorebook(map[string]interface{}{"title": "A", "content": "a"}, false)
	if !errors.Is(err, lorebookstore.ErrNotCutover) {
		t.Fatalf("error = %v, want ErrNotCutover", err)
	}
	entries, readErr := os.ReadDir(root)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("files written before import: %d", len(entries))
	}
}

func TestImportThenSaveListAndDelete(t *testing.T) {
	app, root := newLorebookApp(t)
	imported, err := app.ImportLorebooks([]map[string]interface{}{
		{"id": "legacy-1", "title": "Mana", "content": "Ley lines.", "scope": "world"},
	})
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if imported.Created != 1 || !imported.Cutover {
		t.Fatalf("import result = %+v", imported)
	}
	if _, err := os.Stat(filepath.Join(root, "legacy-1.json")); err != nil {
		t.Fatalf("canonical file missing: %v", err)
	}

	saved, err := app.SaveLorebook(map[string]interface{}{"title": "New", "content": "Body"}, false)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if saved.Status != "saved" {
		t.Fatalf("save status = %+v", saved)
	}

	list, err := app.ListLorebooks()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if list.Count != 2 {
		t.Fatalf("count = %d, want 2", list.Count)
	}

	deleted, err := app.DeleteLorebook("legacy-1")
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if deleted.Status != "deleted" {
		t.Fatalf("delete status = %+v", deleted)
	}
	after, err := app.ListLorebooks()
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if after.Count != 1 {
		t.Fatalf("count after delete = %d, want 1", after.Count)
	}
}

func TestRollbackLorebookImportThroughBinding(t *testing.T) {
	app, root := newLorebookApp(t)
	if _, err := app.ImportLorebooks([]map[string]interface{}{{"id": "legacy-1", "title": "A", "content": "a"}}); err != nil {
		t.Fatalf("import: %v", err)
	}
	rolled, err := app.RollbackLorebookImport()
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if rolled.Status != "rolled_back" {
		t.Fatalf("rollback status = %+v", rolled)
	}
	if _, err := os.Stat(filepath.Join(root, "legacy-1.json")); !os.IsNotExist(err) {
		t.Fatal("imported file survived rollback")
	}
	list, err := app.ListLorebooks()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if list.Cutover {
		t.Fatal("cutover must be cleared after rollback")
	}
}

// A chat persists its own lorebook snapshot, and the loader re-sends it on open.
// After cutover the canonical files own the text, so the prompt must resolve
// against canonical content while the chat keeps its on/off selection.
func TestSetLorebooksResolvesCanonicalContent(t *testing.T) {
	app, _ := newLorebookApp(t)
	if _, err := app.ImportLorebooks([]map[string]interface{}{
		{"id": "lore-1", "title": "Sigil", "content": "The sigil burns cold.", "triggerPhrases": []string{"sigil"}},
	}); err != nil {
		t.Fatalf("import: %v", err)
	}
	if _, err := app.SaveLorebook(map[string]interface{}{
		"id": "lore-1", "title": "Sigil", "content": "The sigil burns white-hot.", "triggerPhrases": []string{"sigil"},
	}, true); err != nil {
		t.Fatalf("canonical edit: %v", err)
	}

	app.chatMessages = []map[string]interface{}{{"role": "user", "content": "I touch the sigil."}}
	snapshot := []map[string]interface{}{
		{"id": "lore-1", "title": "Sigil", "content": "stale text", "triggerPhrases": []string{"sigil"}, "enabled": false},
	}
	if err := app.SetLorebooks(snapshot); err != nil {
		t.Fatalf("set: %v", err)
	}

	system := promptSystem(t, app)
	if strings.Contains(system, "white-hot") {
		t.Fatalf("a disabled entry was injected:\n%s", system)
	}

	if err := app.SetLorebooks([]map[string]interface{}{
		{"id": "lore-1", "title": "Sigil", "content": "stale text", "triggerPhrases": []string{"sigil"}, "enabled": true},
	}); err != nil {
		t.Fatalf("set: %v", err)
	}
	system = promptSystem(t, app)
	if !strings.Contains(system, "white-hot") {
		t.Fatalf("canonical content missing from the prompt:\n%s", system)
	}
	if strings.Contains(system, "stale text") {
		t.Fatalf("chat snapshot text reached the prompt:\n%s", system)
	}
}

// Deleting a canonical entry must not let a chat snapshot resurrect it, while a
// disconnected root must not be treated as a deletion.
func TestDeletedCanonicalEntryDoesNotResurrect(t *testing.T) {
	app, _ := newLorebookApp(t)
	if _, err := app.ImportLorebooks([]map[string]interface{}{
		{"id": "lore-1", "title": "Old", "content": "old canon", "triggerPhrases": []string{"token"}},
		{"id": "lore-2", "title": "Kept", "content": "kept canon", "triggerPhrases": []string{"token"}},
	}); err != nil {
		t.Fatalf("import: %v", err)
	}

	app.chatMessages = []map[string]interface{}{{"role": "user", "content": "the token"}}
	snapshot := []map[string]interface{}{
		{"id": "lore-1", "title": "Old", "content": "old snapshot", "triggerPhrases": []string{"token"}, "enabled": true},
		{"id": "lore-2", "title": "Kept", "content": "kept snapshot", "triggerPhrases": []string{"token"}, "enabled": true},
	}
	if err := app.SetLorebooks(snapshot); err != nil {
		t.Fatalf("set: %v", err)
	}
	if len(app.lorebooks) != 2 {
		t.Fatalf("pool = %d entries, want 2", len(app.lorebooks))
	}

	if _, err := app.DeleteLorebook("lore-1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := app.SetLorebooks(snapshot); err != nil {
		t.Fatalf("set after delete: %v", err)
	}
	if len(app.lorebooks) != 1 || app.lorebooks[0].ID != "lore-2" {
		t.Fatalf("deleted entry survived through the snapshot: %+v", app.lorebooks)
	}
	system := promptSystem(t, app)
	if strings.Contains(system, "old canon") || strings.Contains(system, "old snapshot") {
		t.Fatalf("deleted entry reached the prompt:\n%s", system)
	}
}

func TestDisconnectedRootKeepsLastKnownEntry(t *testing.T) {
	app, root := newLorebookApp(t)
	if _, err := app.ImportLorebooks([]map[string]interface{}{
		{"id": "lore-1", "title": "Offline", "content": "offline canon", "triggerPhrases": []string{"token"}},
	}); err != nil {
		t.Fatalf("import: %v", err)
	}
	if err := os.RemoveAll(root); err != nil {
		t.Fatal(err)
	}

	app.chatMessages = []map[string]interface{}{{"role": "user", "content": "the token"}}
	if err := app.SetLorebooks([]map[string]interface{}{
		{"id": "lore-1", "title": "Offline", "content": "snapshot", "triggerPhrases": []string{"token"}, "enabled": true},
	}); err != nil {
		t.Fatalf("set: %v", err)
	}
	if len(app.lorebooks) != 1 {
		t.Fatalf("incomplete canonical view dropped the last-known entry: %+v", app.lorebooks)
	}
	if !app.lorebooks[0].Stale {
		t.Fatal("last-known entry from an unavailable root must be marked stale")
	}
	if system := promptSystem(t, app); strings.Contains(system, "offline canon") || strings.Contains(system, "snapshot") {
		t.Fatalf("stale lore was injected into the prompt:\n%s", system)
	}
}

// Reconnecting the root clears the stale marker and restores canonical content.
func TestReconnectedRootRestoresCanonicalContent(t *testing.T) {
	app, root := newLorebookApp(t)
	if _, err := app.ImportLorebooks([]map[string]interface{}{
		{"id": "lore-1", "title": "Offline", "content": "canon text", "triggerPhrases": []string{"token"}},
	}); err != nil {
		t.Fatalf("import: %v", err)
	}

	offline := root + "-offline"
	if err := os.Rename(root, offline); err != nil {
		t.Fatal(err)
	}
	snapshot := []map[string]interface{}{
		{"id": "lore-1", "title": "Offline", "content": "snapshot", "triggerPhrases": []string{"token"}, "enabled": true},
	}
	if err := app.SetLorebooks(snapshot); err != nil {
		t.Fatalf("set offline: %v", err)
	}
	if len(app.lorebooks) != 1 || !app.lorebooks[0].Stale {
		t.Fatalf("offline root entry = %+v, want one stale entry", app.lorebooks)
	}

	if err := os.Rename(offline, root); err != nil {
		t.Fatal(err)
	}
	if err := app.SetLorebooks(snapshot); err != nil {
		t.Fatalf("set online: %v", err)
	}
	if len(app.lorebooks) != 1 || app.lorebooks[0].Stale {
		t.Fatalf("reconnected entry = %+v, want one current entry", app.lorebooks)
	}
	if app.lorebooks[0].Content != "canon text" {
		t.Fatalf("content = %q, want canonical text", app.lorebooks[0].Content)
	}
}

// A lorebook created after the chat was saved is not in its snapshot; enabling
// it must attach the canonical entry so it can reach the prompt.
func TestEnablingLorebookCreatedAfterChatSave(t *testing.T) {
	app, _ := newLorebookApp(t)
	if _, err := app.ImportLorebooks([]map[string]interface{}{
		{"id": "lore-1", "title": "Existing", "content": "existing canon"},
	}); err != nil {
		t.Fatalf("import: %v", err)
	}
	if err := app.SetLorebooks(nil); err != nil {
		t.Fatalf("set: %v", err)
	}

	saved, err := app.SaveLorebook(map[string]interface{}{
		"title": "Beacon", "content": "A pale beacon answers.", "triggerPhrases": []string{"beacon"},
	}, false)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if saved.Status != "saved" || saved.ID == "" {
		t.Fatalf("save result = %+v", saved)
	}

	app.chatMessages = []map[string]interface{}{{"role": "user", "content": "I light the beacon."}}
	if err := app.UpdateLorebookContext([]string{saved.ID}); err != nil {
		t.Fatalf("update context: %v", err)
	}
	if len(app.lorebooks) != 1 || !app.lorebooks[0].Enabled {
		t.Fatalf("new canonical entry was not attached and enabled: %+v", app.lorebooks)
	}
	if system := promptSystem(t, app); !strings.Contains(system, "pale beacon") {
		t.Fatalf("new canonical entry did not reach the prompt:\n%s", system)
	}
}

// Before cutover nothing changes: the browser store stays the source.
func TestSetLorebooksKeepsSnapshotBeforeCutover(t *testing.T) {
	app := &App{}
	if err := app.SetLorebooks([]map[string]interface{}{
		{"id": "lore-1", "title": "Legacy", "content": "legacy text"},
	}); err != nil {
		t.Fatalf("set: %v", err)
	}
	if len(app.lorebooks) != 1 || app.lorebooks[0].Content != "legacy text" {
		t.Fatalf("snapshot was rewritten before cutover: %+v", app.lorebooks)
	}
}

func promptSystem(t *testing.T, app *App) string {
	t.Helper()
	messages := app.buildPromptMessages()
	if len(messages) == 0 {
		t.Fatal("no prompt messages")
	}
	system, _ := messages[0].Content.(string)
	return system
}

func TestLorebookBindingsWithoutStoreReportErrors(t *testing.T) {
	app := &App{}
	if _, err := app.ListLorebooks(); err == nil {
		t.Fatal("list without a store must fail")
	}
	if res, err := app.ImportLorebooks(nil); err != nil || res.Status != "failed" {
		t.Fatalf("import without a store must report a failed result, got %+v (err %v)", res, err)
	}
}
