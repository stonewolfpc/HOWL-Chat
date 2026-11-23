package chat

import (
	entities "howl-chat/internal/backend/types/entities"
	"strings"
	"testing"

	"howl-chat/internal/backend/lorebook"
)

func TestBuildMessagesReadsTheSessionDirectly(t *testing.T) {
	s := NewSession()
	s.SetWorld(&WorldContext{World: entities.World{ID: "w1", Name: "Ashfall", Description: "Ash on the wind.", KeyFacts: []string{"no sun"}}})
	s.SetScenario(&ScenarioContext{Scenario: entities.Scenario{ID: "s1", Name: "The Siege", Description: "The walls are failing."}})
	s.SetCharacters([]*CharacterContext{{Character: entities.Character{ID: "char_aria", Name: "Aria", Personality: "wry", Backstory: "A soldier."}}})
	s.ReplaceEntries([]SessionEntry{
		{Role: "user", Content: "we hold the gate", Sender: "User"},
		{Role: "assistant", Content: "then we hold it together", Character: "Aria", CharacterID: "char_aria"},
		{Role: "user", Content: "what do you see?", Sender: "User"},
	})

	messages := s.BuildMessages(PromptSettings{ContextSize: 4096})
	if len(messages) != 4 {
		t.Fatalf("assembled %d messages, want system + 3 turns", len(messages))
	}
	if messages[0].Role != "system" {
		t.Fatalf("first message role = %q", messages[0].Role)
	}
	system, _ := messages[0].Content.(string)
	for _, want := range []string{"[WORLD: Ashfall]", "[CHARACTER SHEETS]", "Aria", "[SCENARIO: The Siege]", "Start your response immediately as the character."} {
		if !strings.Contains(system, want) {
			t.Fatalf("system prompt missing %q:\n%s", want, system)
		}
	}
	// Character sheets precede the scenario on purpose.
	if strings.Index(system, "[CHARACTER SHEETS]") > strings.Index(system, "[SCENARIO:") {
		t.Fatal("scenario was placed before the character sheets")
	}
	if messages[1].Role != "user" || messages[1].Content != "we hold the gate" {
		t.Fatalf("first turn = %+v", messages[1])
	}
	if messages[2].Role != "assistant" || messages[2].Content != "then we hold it together" {
		t.Fatalf("assistant turn = %+v", messages[2])
	}
	// Depth-0 rules ride on the final user turn, not on a trailing system message.
	last, _ := messages[3].Content.(string)
	if !strings.Contains(last, "what do you see?") || !strings.Contains(last, "[System note:") {
		t.Fatalf("final user turn = %q", last)
	}
}

func TestBuildMessagesAppliesPrefixesAndSystemOverride(t *testing.T) {
	s := NewSession()
	s.ReplaceEntries([]SessionEntry{
		{Role: "user", Content: "hello"},
		{Role: "assistant", Content: "hi"},
		{Role: "user", Content: "again"},
	})
	messages := s.BuildMessages(PromptSettings{
		ContextSize: 4096,
		ModelSettings: map[string]interface{}{
			"user_prefix":            "U> ",
			"assistant_prefix":       "A> ",
			"system_prompt_override": "You are terse.",
		},
	})
	system, _ := messages[0].Content.(string)
	if !strings.HasPrefix(system, "You are terse.") {
		t.Fatalf("system override not applied: %q", system)
	}
	if messages[1].Content != "U> hello" {
		t.Fatalf("user prefix = %q", messages[1].Content)
	}
	if messages[2].Content != "A> hi" {
		t.Fatalf("assistant prefix = %q", messages[2].Content)
	}
}

func TestBuildMessagesResolvesLoreFromTheSession(t *testing.T) {
	s := NewSession()
	s.SetLorebooks([]lorebook.Entry{{
		ID:               "lb1",
		Title:            "Sigil",
		Content:          "The sigil answers.",
		Enabled:          true,
		Scope:            lorebook.ScopeWorld,
		TriggerPhrases:   []string{"sigil"},
		TriggerMode:      lorebook.TriggerLoose,
		TriggerDirection: lorebook.TriggerBoth,
		PriorityLevel:    3,
	}})
	s.ReplaceEntries([]SessionEntry{{Role: "user", Content: "I touch the sigil"}})

	system, _ := s.BuildMessages(PromptSettings{ContextSize: 4096})[0].Content.(string)
	if !strings.Contains(system, "The sigil answers.") {
		t.Fatalf("lore did not reach the prompt:\n%s", system)
	}
}

func TestBuildMessagesTruncatesOlderTurnsToFitContext(t *testing.T) {
	s := NewSession()
	s.ReplaceEntries([]SessionEntry{
		{Role: "user", Content: strings.Repeat("old ", 400)},
		{Role: "assistant", Content: strings.Repeat("reply ", 100)},
		{Role: "user", Content: "newest"},
	})
	// 200 tokens = 800 characters: the system prompt and the newest turn fit, the
	// oldest cannot.
	messages := s.BuildMessages(PromptSettings{ContextSize: 200})
	if len(messages) < 2 {
		t.Fatalf("truncation removed the newest turn: %d messages", len(messages))
	}
	last, _ := messages[len(messages)-1].Content.(string)
	if !strings.Contains(last, "newest") {
		t.Fatalf("newest turn was dropped: %q", last)
	}
	for _, m := range messages {
		if content, ok := m.Content.(string); ok && strings.Contains(content, "old old old") {
			t.Fatal("oldest turn survived truncation")
		}
	}
}

func TestStripThinkBlocksKeepsTheAnswerOnly(t *testing.T) {
	cases := map[string]string{
		"visible":                             "visible",
		"<think>hidden</think>visible":        "visible",
		"<|think|>hidden<|/think|>visible":    "visible",
		"visible</think>":                     "visible",
		"answer <thinkingcut off mid-thought": "answer",
	}
	for in, want := range cases {
		if got := StripThinkBlocks(in); got != want {
			t.Fatalf("StripThinkBlocks(%q) = %q, want %q", in, got, want)
		}
	}
}
