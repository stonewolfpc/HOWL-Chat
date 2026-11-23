package chat

import (
	"errors"
	"os"
	"testing"
)

// A write that merely mentions an identifier must never bring a conversation
// into existence. Append and rewrite are not the creation route, so a missing
// directory is refused and nothing is created.

func TestAppendDoesNotCreateAMissingConversation(t *testing.T) {
	pm := newBoundaryTestPM(t)
	id := "2026-11-22T17-15-55Z"

	err := pm.AppendHistoryEntry(id, ChatHistory{Text: "late", Role: "assistant"})

	var refusal *Refusal
	if !errors.As(err, &refusal) {
		t.Fatalf("append to a missing conversation returned %v, want a refusal", err)
	}
	if refusal.Code != RefusedUnbound {
		t.Fatalf("append refusal code = %s, want %s", refusal.Code, RefusedUnbound)
	}

	dir, dirErr := pm.conversationDir(id)
	if dirErr != nil {
		t.Fatalf("conversation dir: %v", dirErr)
	}
	if _, statErr := os.Stat(dir); !os.IsNotExist(statErr) {
		t.Fatalf("the refused append created %q (stat error: %v)", dir, statErr)
	}
}

func TestRewriteDoesNotCreateAMissingConversation(t *testing.T) {
	pm := newBoundaryTestPM(t)
	id := "2026-11-22T17-15-56Z"

	err := pm.RewriteHistory(id, []ChatHistory{{Text: "first", Role: "user"}})

	var refusal *Refusal
	if !errors.As(err, &refusal) {
		t.Fatalf("rewrite of a missing conversation returned %v, want a refusal", err)
	}
	if refusal.Code != RefusedUnbound {
		t.Fatalf("rewrite refusal code = %s, want %s", refusal.Code, RefusedUnbound)
	}

	dir, dirErr := pm.conversationDir(id)
	if dirErr != nil {
		t.Fatalf("conversation dir: %v", dirErr)
	}
	if _, statErr := os.Stat(dir); !os.IsNotExist(statErr) {
		t.Fatalf("the refused rewrite created %q (stat error: %v)", dir, statErr)
	}
}

func TestSaveRefusesToRecreateAKnownConversationWhoseDirectoryIsGone(t *testing.T) {
	pm := newBoundaryTestPM(t)
	chatData, err := pm.NewChat("vanishing")
	if err != nil {
		t.Fatalf("create chat: %v", err)
	}
	dir, err := pm.conversationDir(chatData.Meta.ID)
	if err != nil {
		t.Fatalf("conversation dir: %v", err)
	}
	// Something outside this program removed the folder (cleanup tool, sync
	// client, hand deletion). A later save must report it, not rebuild it.
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("remove conversation directory: %v", err)
	}

	err = pm.SaveChat(chatData)
	var refusal *Refusal
	if !errors.As(err, &refusal) || refusal.Code != RefusedUnbound {
		t.Fatalf("save of a known conversation with no directory returned %v, want %s", err, RefusedUnbound)
	}
	if _, statErr := os.Stat(dir); !os.IsNotExist(statErr) {
		t.Fatalf("the refused save recreated %q", dir)
	}
}
