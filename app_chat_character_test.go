package main

import (
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/backend/chat"
)

func newCharacterOwnershipTestApp(t *testing.T) (*App, *chat.PersistenceManager, string, string) {
	t.Helper()
	rootDir := t.TempDir()
	pm, err := chat.NewPersistenceManagerWithDir(rootDir)
	if err != nil {
		t.Fatalf("NewPersistenceManagerWithDir: %v", err)
	}
	chatData, err := pm.NewChat("Character ownership")
	if err != nil {
		t.Fatalf("NewChat: %v", err)
	}
	return &App{chatPersistence: pm}, pm, rootDir, chatData.Meta.ID
}

func TestCharacterOwnershipSurvivesAppendSaveRestart(t *testing.T) {
	app, pm, _, chatID := newCharacterOwnershipTestApp(t)

	if err := app.appendChatMessage("assistant", "Aria line", "Aria", "char_aria"); err != nil {
		t.Fatalf("append Aria: %v", err)
	}
	if err := app.appendChatMessage("assistant", "Doren line", "Doren", "char_doren"); err != nil {
		t.Fatalf("append Doren: %v", err)
	}
	if err := app.performSave(); err != nil {
		t.Fatalf("performSave: %v", err)
	}

	// Simulate a fresh application instance restoring the persisted chat.
	restarted := &App{chatPersistence: pm}
	loaded, err := restarted.LoadChat(chatID)
	if err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	if len(loaded.History) != 2 {
		t.Fatalf("History=%d, want 2", len(loaded.History))
	}
	if loaded.History[0].CharacterID != "char_aria" || loaded.History[1].CharacterID != "char_doren" {
		t.Fatalf("restored owners=%q/%q, want char_aria/char_doren", loaded.History[0].CharacterID, loaded.History[1].CharacterID)
	}
	if got, _ := restarted.chatMessages[0]["characterId"].(string); got != "char_aria" {
		t.Fatalf("in-memory owner=%q, want char_aria", got)
	}
	if got, _ := restarted.chatMessages[1]["characterId"].(string); got != "char_doren" {
		t.Fatalf("in-memory owner=%q, want char_doren", got)
	}

	// Final-save rewrites history.ndjson from the restored state and must keep owners.
	restarted.performFinalSave()
	again, err := pm.LoadChat(chatID)
	if err != nil {
		t.Fatalf("LoadChat after final save: %v", err)
	}
	if len(again.History) != 2 || again.History[0].CharacterID != "char_aria" || again.History[1].CharacterID != "char_doren" {
		t.Fatalf("owners after final save=%+v", again.History)
	}
}

func TestReplaceChatHistoryPreservesCharacterOwner(t *testing.T) {
	app, pm, _, chatID := newCharacterOwnershipTestApp(t)

	err := app.ReplaceChatHistory([]map[string]interface{}{
		{"role": "assistant", "text": "Aria line", "sender": "Aria", "characterId": "char_aria", "timestamp": "2026-01-01T00:00:00Z", "hasImage": true},
		{"role": "user", "text": "reply", "sender": "User"},
		{"role": "assistant", "text": "Doren line", "sender": "Doren", "characterId": "char_doren"},
	})
	if err != nil {
		t.Fatalf("ReplaceChatHistory: %v", err)
	}

	loaded, err := pm.LoadChat(chatID)
	if err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	if len(loaded.History) != 3 {
		t.Fatalf("History=%d, want 3", len(loaded.History))
	}
	if loaded.History[0].CharacterID != "char_aria" || loaded.History[2].CharacterID != "char_doren" {
		t.Fatalf("replace owners=%q/%q", loaded.History[0].CharacterID, loaded.History[2].CharacterID)
	}
	if loaded.History[0].Timestamp != "2026-01-01T00:00:00Z" || !loaded.History[0].HasImage {
		t.Fatalf("replace metadata=%+v", loaded.History[0])
	}
	if got, _ := app.chatMessages[0]["characterId"].(string); got != "char_aria" {
		t.Fatalf("in-memory replace owner=%q", got)
	}
	if err := app.performSave(); err != nil {
		t.Fatalf("performSave: %v", err)
	}
	saved, err := pm.LoadChat(chatID)
	if err != nil {
		t.Fatalf("LoadChat after save: %v", err)
	}
	if saved.History[0].Timestamp != "2026-01-01T00:00:00Z" || !saved.History[0].HasImage {
		t.Fatalf("replace save lost metadata: %+v", saved.History[0])
	}
}

func TestLegacyHistoryRemainsUnassignedThroughAppSave(t *testing.T) {
	app, pm, rootDir, chatID := newCharacterOwnershipTestApp(t)

	legacy := `{"id":"0","sender":"Aria","text":"legacy line","role":"assistant","timestamp":"2026-01-01T00:00:00Z"}` + "\n"
	ndjsonPath := filepath.Join(rootDir, "chats", chatID, "history.ndjson")
	if err := os.WriteFile(ndjsonPath, []byte(legacy), 0644); err != nil {
		t.Fatalf("write legacy history: %v", err)
	}

	if _, err := app.LoadChat(chatID); err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	if _, exists := app.chatMessages[0]["characterId"]; exists {
		t.Fatalf("legacy owner was synthesized: %+v", app.chatMessages[0])
	}
	if err := app.performSave(); err != nil {
		t.Fatalf("performSave: %v", err)
	}
	app.performFinalSave()

	loaded, err := pm.LoadChat(chatID)
	if err != nil {
		t.Fatalf("LoadChat after save: %v", err)
	}
	if len(loaded.History) != 1 || loaded.History[0].CharacterID != "" {
		t.Fatalf("legacy ownership changed: %+v", loaded.History)
	}
	if loaded.History[0].Text != "legacy line" || loaded.History[0].Sender != "Aria" {
		t.Fatalf("legacy content changed: %+v", loaded.History[0])
	}
}

func TestConvertFrontendHistoryPreservesCharacterOwner(t *testing.T) {
	app := &App{}
	converted, err := app.convertFrontendDataToChatData(map[string]interface{}{
		"meta": map[string]interface{}{"id": "chat_owner_test"},
		"history": []interface{}{
			map[string]interface{}{
				"id": "0", "sender": "Aria", "text": "hello", "role": "assistant",
				"timestamp": "2026-01-01T00:00:00Z", "characterId": "char_aria",
			},
			map[string]interface{}{
				"id": "1", "sender": "User", "text": "reply", "role": "user",
				"timestamp": "2026-01-01T00:00:01Z",
			},
		},
	})
	if err != nil {
		t.Fatalf("convertFrontendDataToChatData: %v", err)
	}
	if len(converted.History) != 2 {
		t.Fatalf("History=%d, want 2", len(converted.History))
	}
	if converted.History[0].CharacterID != "char_aria" {
		t.Fatalf("converted owner=%q, want char_aria", converted.History[0].CharacterID)
	}
	if converted.History[1].CharacterID != "" {
		t.Fatalf("user owner=%q, want empty", converted.History[1].CharacterID)
	}
}
