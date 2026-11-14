package chat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChatHistoryCharacterIDJSONRoundTrip(t *testing.T) {
	original := ChatHistory{
		ID:          "0",
		Sender:      "Aria",
		Text:        "Hello",
		Role:        "assistant",
		Timestamp:   "2026-01-01T00:00:00Z",
		CharacterID: "char_aria",
	}

	raw, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded ChatHistory
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.CharacterID != "char_aria" {
		t.Fatalf("CharacterID=%q, want char_aria", decoded.CharacterID)
	}
}

func TestHistoryCharacterIDSurvivesAppendLoadAndRewrite(t *testing.T) {
	pm, err := NewPersistenceManagerWithDir(t.TempDir())
	if err != nil {
		t.Fatalf("NewPersistenceManagerWithDir: %v", err)
	}
	chatData, err := pm.NewChat("Ownership")
	if err != nil {
		t.Fatalf("NewChat: %v", err)
	}

	entries := []ChatHistory{
		{ID: "0", Sender: "Aria", Text: "first", Role: "assistant", Timestamp: "2026-01-01T00:00:00Z", CharacterID: "char_aria"},
		{ID: "1", Sender: "User", Text: "reply", Role: "user", Timestamp: "2026-01-01T00:00:01Z"},
		{ID: "2", Sender: "Doren", Text: "second", Role: "assistant", Timestamp: "2026-01-01T00:00:02Z", CharacterID: "char_doren"},
	}
	for _, entry := range entries {
		if err := pm.AppendHistoryEntry(chatData.Meta.ID, entry); err != nil {
			t.Fatalf("AppendHistoryEntry: %v", err)
		}
	}

	assertHistoryOwners(t, pm, chatData.Meta.ID, []string{"char_aria", "", "char_doren"})

	loaded, err := pm.LoadChat(chatData.Meta.ID)
	if err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	if err := pm.RewriteHistory(chatData.Meta.ID, loaded.History); err != nil {
		t.Fatalf("RewriteHistory: %v", err)
	}
	assertHistoryOwners(t, pm, chatData.Meta.ID, []string{"char_aria", "", "char_doren"})
}

func TestLegacyHistoryWithoutCharacterIDLoadsUnassigned(t *testing.T) {
	rootDir := t.TempDir()
	pm, err := NewPersistenceManagerWithDir(rootDir)
	if err != nil {
		t.Fatalf("NewPersistenceManagerWithDir: %v", err)
	}
	chatData, err := pm.NewChat("Legacy")
	if err != nil {
		t.Fatalf("NewChat: %v", err)
	}

	legacy := `{"id":"0","sender":"Aria","text":"legacy line","role":"assistant","timestamp":"2026-01-01T00:00:00Z"}` + "\n"
	ndjsonPath := filepath.Join(rootDir, "chats", chatData.Meta.ID, "history.ndjson")
	if err := os.WriteFile(ndjsonPath, []byte(legacy), 0644); err != nil {
		t.Fatalf("write legacy history: %v", err)
	}

	loaded, err := pm.LoadChat(chatData.Meta.ID)
	if err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	if len(loaded.History) != 1 {
		t.Fatalf("History=%d, want 1", len(loaded.History))
	}
	if loaded.History[0].CharacterID != "" {
		t.Fatalf("legacy CharacterID=%q, want empty", loaded.History[0].CharacterID)
	}
	if loaded.History[0].Text != "legacy line" {
		t.Fatalf("legacy text=%q, want legacy line", loaded.History[0].Text)
	}

	// Rewriting legacy history must be deterministic, safe to repeat, and must
	// not invent an owner.
	for i := 0; i < 2; i++ {
		if err := pm.RewriteHistory(chatData.Meta.ID, loaded.History); err != nil {
			t.Fatalf("RewriteHistory pass %d: %v", i+1, err)
		}
		again, err := pm.LoadChat(chatData.Meta.ID)
		if err != nil {
			t.Fatalf("LoadChat after rewrite pass %d: %v", i+1, err)
		}
		if len(again.History) != 1 || again.History[0].CharacterID != "" {
			t.Fatalf("rewrite pass %d invented ownership: %+v", i+1, again.History)
		}
		if again.History[0].Text != "legacy line" || again.History[0].Timestamp != "2026-01-01T00:00:00Z" {
			t.Fatalf("rewrite pass %d changed legacy content: %+v", i+1, again.History[0])
		}
		raw, err := os.ReadFile(ndjsonPath)
		if err != nil {
			t.Fatalf("read rewritten history pass %d: %v", i+1, err)
		}
		if strings.Contains(string(raw), "characterId") {
			t.Fatalf("unassigned legacy entry gained a characterId on pass %d: %s", i+1, raw)
		}
		loaded = again
	}
}

func TestSaveChatAsPreservesCharacterID(t *testing.T) {
	rootDir := t.TempDir()
	pm, err := NewPersistenceManagerWithDir(rootDir)
	if err != nil {
		t.Fatalf("NewPersistenceManagerWithDir: %v", err)
	}
	chatData, err := pm.NewChat("Original")
	if err != nil {
		t.Fatalf("NewChat: %v", err)
	}
	if err := pm.AppendHistoryEntry(chatData.Meta.ID, ChatHistory{
		ID: "0", Sender: "Aria", Text: "copy me", Role: "assistant",
		Timestamp: "2026-01-01T00:00:00Z", CharacterID: "char_aria",
	}); err != nil {
		t.Fatalf("AppendHistoryEntry: %v", err)
	}

	copyData, err := pm.SaveChatAs(chatData.Meta.ID, "Copy")
	if err != nil {
		t.Fatalf("SaveChatAs: %v", err)
	}
	if len(copyData.History) != 1 || copyData.History[0].CharacterID != "char_aria" {
		t.Fatalf("SaveChatAs history=%+v, want char_aria", copyData.History)
	}
	assertHistoryOwners(t, pm, copyData.Meta.ID, []string{"char_aria"})
}

func assertHistoryOwners(t *testing.T, pm *PersistenceManager, chatID string, want []string) {
	t.Helper()
	loaded, err := pm.LoadChat(chatID)
	if err != nil {
		t.Fatalf("LoadChat(%s): %v", chatID, err)
	}
	if len(loaded.History) != len(want) {
		t.Fatalf("History=%d, want %d: %+v", len(loaded.History), len(want), loaded.History)
	}
	for i, entry := range loaded.History {
		if entry.CharacterID != want[i] {
			t.Fatalf("history[%d].CharacterID=%q, want %q", i, entry.CharacterID, want[i])
		}
	}
}
