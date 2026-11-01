package main

import (
	"strings"
	"testing"

	"howl-chat/internal/backend/lorebook"
)

func dragonLoreApp(chatMessages []map[string]interface{}) *App {
	return &App{
		lorebooks: []lorebook.Entry{{
			ID:               "lb1",
			Title:            "Dragons",
			Content:          "Dragons sleep beneath the city.",
			Enabled:          true,
			Scope:            lorebook.ScopeWorld,
			TriggerPhrases:   []string{"dragon"},
			TriggerMode:      lorebook.TriggerLoose,
			TriggerDirection: lorebook.TriggerBoth,
			TriggerFrequency: lorebook.FrequencyAlways,
			PriorityLevel:    3,
		}},
		chatMessages: chatMessages,
	}
}

func loreInPrompt(t *testing.T, app *App) string {
	t.Helper()
	messages := app.buildPromptMessages()
	if len(messages) == 0 {
		t.Fatal("no messages built")
	}
	system, _ := messages[0].Content.(string)
	return system
}

// A regenerate leaves the assistant reply last. Lore must still resolve then.
func TestLoreResolvesOnRegenerate(t *testing.T) {
	app := dragonLoreApp([]map[string]interface{}{
		{"role": "user", "content": "I search for the dragon"},
		{"role": "assistant", "content": "The cave is dark."},
	})
	system := loreInPrompt(t, app)
	if !strings.Contains(system, "Dragons sleep beneath the city.") {
		t.Fatalf("lore did not resolve when the last message was assistant:\n%s", system)
	}
}

// The ordinary case must keep working.
func TestLoreResolvesOnNormalTurn(t *testing.T) {
	app := dragonLoreApp([]map[string]interface{}{
		{"role": "user", "content": "I search for the dragon"},
	})
	system := loreInPrompt(t, app)
	if !strings.Contains(system, "Dragons sleep beneath the city.") {
		t.Fatalf("lore did not resolve on a normal user turn:\n%s", system)
	}
}

// No matching keyword means no lore, and no crash on an empty transcript.
func TestLoreAbsentWithoutKeyword(t *testing.T) {
	app := dragonLoreApp([]map[string]interface{}{
		{"role": "user", "content": "I look at the sky"},
	})
	if system := loreInPrompt(t, app); strings.Contains(system, "Dragons sleep beneath the city.") {
		t.Errorf("lore fired without a keyword:\n%s", system)
	}
	empty := &App{}
	if messages := empty.buildPromptMessages(); len(messages) != 1 {
		t.Errorf("expected only the system message, got %d", len(messages))
	}
}
