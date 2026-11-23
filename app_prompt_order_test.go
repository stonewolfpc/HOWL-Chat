package main

import (
	"strings"
	"testing"

	"howl-chat/internal/backend/chat"
	entities "howl-chat/internal/backend/types/entities"
)

// The order below is the one the research converges on: instructions, then
// impersonal world facts, then the character card, then the scene, with the
// rules that must survive a long chat restated at depth 0.
func TestBuildPromptMessagesOrder(t *testing.T) {
	app := testApp()
	app.session().SetWorld(&chat.WorldContext{World: entities.World{
		Name: "Aethermoor", Description: "WorldDesc", KeyFacts: []string{"WorldFact"},
	}})
	app.session().SetScenario(&chat.ScenarioContext{Scenario: entities.Scenario{
		Name: "Revolt", Description: "ScenarioDesc", KeyDetails: []string{"ScenarioDetail"},
	}})
	app.session().SetCharacters([]*chat.CharacterContext{
		{Character: entities.Character{
			Name: "Doren", Personality: "CharPersonality", Backstory: "CharBackstory",
			Goals: []string{"CharGoal"}, Relationships: []string{"CharRelation"},
		}},
	})
	app.setTestMessages([]map[string]interface{}{
		{"role": "user", "content": "first turn"},
		{"role": "assistant", "content": "reply"},
		{"role": "user", "content": "second turn"},
	})

	messages := promptMessages(app)
	if len(messages) == 0 || messages[0].Role != "system" {
		t.Fatalf("expected a leading system message, got %#v", messages)
	}
	system, _ := messages[0].Content.(string)

	for _, marker := range []string{"You are in an immersive roleplay", "[WORLD: Aethermoor]", "Character: Doren", "[SCENARIO: Revolt]"} {
		if !strings.Contains(system, marker) {
			t.Fatalf("system prompt missing %q:\n%s", marker, system)
		}
	}
	prev := -1
	for _, marker := range []string{"You are in an immersive roleplay", "[WORLD: Aethermoor]", "Character: Doren", "[SCENARIO: Revolt]"} {
		idx := strings.Index(system, marker)
		if idx <= prev {
			t.Errorf("%q is out of order (offset %d, previous marker at %d)", marker, idx, prev)
		}
		prev = idx
	}

	last := messages[len(messages)-1]
	if last.Role != "user" {
		t.Fatalf("expected the final message to be the user turn, got %q", last.Role)
	}
	lastContent, _ := last.Content.(string)
	if !strings.Contains(lastContent, "[System note:") {
		t.Errorf("depth-0 instructions missing from the final turn: %q", lastContent)
	}
	if !strings.Contains(lastContent, "Do not write the user's words") {
		t.Errorf("depth-0 instructions are incomplete: %q", lastContent)
	}

	firstContent, _ := messages[1].Content.(string)
	if strings.Contains(firstContent, "[System note:") {
		t.Errorf("depth-0 note leaked into an earlier turn: %q", firstContent)
	}
	if app.session().Len() != 3 {
		t.Errorf("session history was mutated, length = %d", app.session().Len())
	}
	if stored := app.session().Entries()[2].Content; strings.Contains(stored, "[System note:") {
		t.Errorf("depth-0 note was written into stored history: %q", stored)
	}
}

// A prompt with no roleplay context loaded must not grow roleplay scaffolding.
func TestBuildPromptMessagesWithoutContext(t *testing.T) {
	app := testApp()
	app.setTestMessages([]map[string]interface{}{{"role": "user", "content": "hi"}})
	messages := promptMessages(app)
	system, _ := messages[0].Content.(string)
	if !strings.Contains(system, "helpful, knowledgeable assistant") {
		t.Errorf("expected assistant system prompt, got %q", system)
	}
	if strings.Contains(system, "[WORLD") || strings.Contains(system, "[SCENARIO") {
		t.Errorf("context blocks present without context: %q", system)
	}
}
