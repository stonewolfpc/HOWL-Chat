package main

import (
	"testing"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/llama"
	appruntime "howl-chat/internal/runtime"
	appsettings "howl-chat/internal/settings"
)

// testApp returns an App with an empty session and nothing else: no boot, no
// disk, no services. The active conversation lives in the chat domain, so a test
// that needs turns, a world or lorebooks fills the session rather than App.
func testApp() *App { return &App{conversation: chat.NewActiveConversation(nil)} }

// newTestSettingsAuthority builds a settings owner rooted at dir, the same way
// boot does, without loading anything from disk.
func newTestSettingsAuthority(t *testing.T, dir string) *appruntime.Settings {
	t.Helper()
	return appruntime.NewSettings(appsettings.NewStore(dir), dir, nil)
}

// sessionEntriesFromMaps converts renderer-shaped messages into session turns.
func sessionEntriesFromMaps(messages []map[string]interface{}) []chat.SessionEntry {
	entries := make([]chat.SessionEntry, 0, len(messages))
	for _, m := range messages {
		var entry chat.SessionEntry
		entry.Role, _ = m["role"].(string)
		entry.Content, _ = m["content"].(string)
		entry.Sender, _ = m["sender"].(string)
		entry.Character, _ = m["character"].(string)
		entry.CharacterID, _ = m["characterId"].(string)
		entry.HasImage, _ = m["hasImage"].(bool)
		entry.Timestamp, _ = m["timestamp"].(string)
		entries = append(entries, entry)
	}
	return entries
}

// setTestMessages replaces the session's turns from renderer-shaped maps.
func (a *App) setTestMessages(messages []map[string]interface{}) {
	a.session().ReplaceEntries(sessionEntriesFromMaps(messages))
}

// promptMessages assembles the message list through the same route the send path
// uses. Assembly itself lives in the chat domain now.
func promptMessages(app *App) []llama.ChatMessage {
	return app.session().BuildMessages(app.promptSettings())
}
