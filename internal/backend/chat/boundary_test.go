package chat

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/durable"
)

func newBoundaryTestPM(t *testing.T) *PersistenceManager {
	t.Helper()
	pm, err := NewPersistenceManagerWithDir(t.TempDir())
	if err != nil {
		t.Fatalf("new persistence manager: %v", err)
	}
	return pm
}

func TestConversationDirRefusesNonComponentIdentifiers(t *testing.T) {
	pm := newBoundaryTestPM(t)

	for _, id := range []string{"", ".", "..", "a/b", `a\b`, "C:evil"} {
		if _, err := pm.conversationDir(id); err == nil {
			t.Fatalf("identifier %q was accepted; it must be refused", id)
		} else if !IsRefusal(err, RefusedContainment) {
			t.Fatalf("identifier %q produced %v, want refusal %s", id, err, RefusedContainment)
		}
	}

	dir, err := pm.conversationDir("2026-11-22T17-15-55Z")
	if err != nil {
		t.Fatalf("a normal conversation identifier was refused: %v", err)
	}
	if !strings.HasPrefix(dir, pm.RootDir()+string(filepath.Separator)) {
		t.Fatalf("conversation directory %q is not inside %q", dir, pm.RootDir())
	}
}

func TestAppendTerminatesUnterminatedTail(t *testing.T) {
	pm := newBoundaryTestPM(t)
	chatData, err := pm.NewChat("tail")
	if err != nil {
		t.Fatalf("create chat: %v", err)
	}

	chatDir, err := pm.conversationDir(chatData.Meta.ID)
	if err != nil {
		t.Fatalf("conversation dir: %v", err)
	}
	ndjsonPath := filepath.Join(chatDir, "history.ndjson")

	// A complete record with no trailing newline: the shape row 579 accepts.
	seeded := []byte(`{"text":"first","role":"user"}`)
	if err := os.WriteFile(ndjsonPath, seeded, 0o644); err != nil {
		t.Fatalf("seed unterminated tail: %v", err)
	}

	if err := pm.AppendHistoryEntry(chatData.Meta.ID, ChatHistory{Text: "second", Role: "assistant"}); err != nil {
		t.Fatalf("append after unterminated tail: %v", err)
	}

	data, err := os.ReadFile(ndjsonPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines after the append, got %d: %q", len(lines), string(data))
	}
	var first, second ChatHistory
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("first line was not a standalone record (the append was concatenated): %v in %q", err, lines[0])
	}
	if first.Text != "first" {
		t.Fatalf("first line changed: %q", first.Text)
	}
	if err := json.Unmarshal([]byte(lines[1]), &second); err != nil {
		t.Fatalf("second line is not a record: %v in %q", err, lines[1])
	}
	if second.Text != "second" {
		t.Fatalf("second line is %q, want %q", second.Text, "second")
	}
}

func TestSaveRefusedWhileAnotherWriterHoldsTheConversation(t *testing.T) {
	pm := newBoundaryTestPM(t)
	chatData, err := pm.NewChat("locked")
	if err != nil {
		t.Fatalf("create chat: %v", err)
	}

	unlock, err := durable.LockFile(pm.lockPath(chatData.Meta.ID), 0)
	if err != nil {
		t.Fatalf("take the lock as a competing writer: %v", err)
	}
	defer unlock()

	if err := pm.SaveChat(chatData); err == nil {
		t.Fatal("save succeeded while another writer held the lock")
	} else if !IsRefusal(err, RefusedBusy) {
		t.Fatalf("save produced %v, want refusal %s", err, RefusedBusy)
	}
}

func TestReentrantSaveIsRefusedWithItsOwnCode(t *testing.T) {
	pm := newBoundaryTestPM(t)
	chatData, err := pm.NewChat("reentrant")
	if err != nil {
		t.Fatalf("create chat: %v", err)
	}

	lockPath := pm.lockPath(chatData.Meta.ID)
	processChatLocks.mu.Lock()
	processChatLocks.held[lockPath] = true
	processChatLocks.mu.Unlock()
	defer func() {
		processChatLocks.mu.Lock()
		delete(processChatLocks.held, lockPath)
		processChatLocks.mu.Unlock()
	}()

	if err := pm.SaveChat(chatData); err == nil {
		t.Fatal("re-entrant save was accepted")
	} else if !IsRefusal(err, RefusedReentrant) {
		t.Fatalf("re-entrant save produced %v, want refusal %s", err, RefusedReentrant)
	}
}

func TestDeletedConversationStaysHiddenAndCannotBeRecreated(t *testing.T) {
	pm := newBoundaryTestPM(t)
	chatData, err := pm.NewChat("doomed")
	if err != nil {
		t.Fatalf("create chat: %v", err)
	}
	id := chatData.Meta.ID

	chatDir, err := pm.conversationDir(id)
	if err != nil {
		t.Fatalf("conversation dir: %v", err)
	}

	if err := pm.DeleteChat(id); err != nil {
		t.Fatalf("delete chat: %v", err)
	}

	summaries, err := pm.GetChatList()
	if err != nil {
		t.Fatalf("list chats: %v", err)
	}
	if len(summaries) != 0 {
		t.Fatalf("deleted conversation is still listed: %+v", summaries)
	}

	// A blocked delete leaves the directory behind. Recreating it by hand must not
	// bring the conversation back, and nothing may write into it again.
	if err := os.MkdirAll(chatDir, 0o755); err != nil {
		t.Fatalf("recreate leftover directory: %v", err)
	}
	meta := []byte(`{"id":"` + id + `","name":"doomed"}`)
	if err := os.WriteFile(filepath.Join(chatDir, "meta.json"), meta, 0o644); err != nil {
		t.Fatalf("write leftover meta: %v", err)
	}

	summaries, err = pm.GetChatList()
	if err != nil {
		t.Fatalf("list chats after leftover: %v", err)
	}
	// The leftover is no longer silently hidden: it is reported as orphaned. What
	// it must never become is a present, loadable conversation.
	for _, summary := range summaries {
		if summary.State != string(StateOrphaned) {
			t.Fatalf("a leftover directory resurrected the deleted conversation: %+v", summary)
		}
	}
	if len(summaries) != 1 {
		t.Fatalf("leftover directory was not reported as one orphan: %+v", summaries)
	}

	if err := pm.SaveChat(chatData); err == nil {
		t.Fatal("save into a deleted conversation was accepted")
	} else if !IsRefusal(err, RefusedUnbound) {
		t.Fatalf("save into a deleted conversation produced %v, want %s", err, RefusedUnbound)
	}
	if _, err := pm.LoadChat(id); err == nil {
		t.Fatal("load of a deleted conversation was accepted")
	}
	if err := pm.AppendHistoryEntry(id, ChatHistory{Text: "late", Role: "assistant"}); err == nil {
		t.Fatal("append into a deleted conversation was accepted")
	}
}

// Two creations in the same second must not share a directory. The second claims
// a suffixed identifier instead of adopting the first one's metadata and history.
func TestConcurrentCreationsClaimDistinctDirectories(t *testing.T) {
	pm := newBoundaryTestPM(t)

	seen := map[string]string{}
	for i := 0; i < 3; i++ {
		chatData, err := pm.NewChat(fmt.Sprintf("name-%d", i))
		if err != nil {
			t.Fatalf("create chat %d: %v", i, err)
		}
		if prior, exists := seen[chatData.Meta.ID]; exists {
			t.Fatalf("identifier %q was claimed twice (%q and %q)", chatData.Meta.ID, prior, chatData.Meta.Name)
		}
		seen[chatData.Meta.ID] = chatData.Meta.Name
	}

	// Every conversation must still be loadable under its own name: no creation
	// overwrote another's metadata.
	for id, name := range seen {
		loaded, err := pm.LoadChat(id)
		if err != nil {
			t.Fatalf("load %q: %v", id, err)
		}
		if loaded.Meta.Name != name {
			t.Fatalf("conversation %q holds name %q, want %q", id, loaded.Meta.Name, name)
		}
	}
}

func TestDeletionIsIdempotent(t *testing.T) {
	pm := newBoundaryTestPM(t)
	chatData, err := pm.NewChat("twice")
	if err != nil {
		t.Fatalf("create chat: %v", err)
	}
	if err := pm.DeleteChat(chatData.Meta.ID); err != nil {
		t.Fatalf("first delete: %v", err)
	}
	if err := pm.DeleteChat(chatData.Meta.ID); err != nil {
		t.Fatalf("second delete must not fail the item that is already gone: %v", err)
	}
}
