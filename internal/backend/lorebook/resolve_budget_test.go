package lorebook

import (
	"strings"
	"testing"
)

// Equal trigger scores fall back to priority, 1 being the highest.
func TestResolveOrdersEqualScoresByPriority(t *testing.T) {
	makeEntry := func(id string, priority int) Entry {
		return withDefaults(Entry{
			ID:               id,
			Title:            id,
			Content:          id + " content",
			Enabled:          true,
			Scope:            ScopeWorld,
			TriggerPhrases:   []string{"sigil"},
			TriggerMode:      TriggerExact,
			TriggerDirection: TriggerBoth,
			TriggerFrequency: FrequencyAlways,
			PriorityLevel:    priority,
		})
	}

	resolved := Resolve([]Entry{makeEntry("low", 5), makeEntry("high", 1)}, ResolveRequest{Message: "the sigil"})
	if len(resolved) != 2 {
		t.Fatalf("resolved %d entries, want 2: %+v", len(resolved), resolved)
	}
	if resolved[0].ID != "high" || resolved[1].ID != "low" {
		t.Fatalf("order = %s, %s; want high, low", resolved[0].ID, resolved[1].ID)
	}

	block := BuildPromptBlock(resolved)
	if strings.Index(block, "high") > strings.Index(block, "low") {
		t.Fatalf("prompt block is out of priority order:\n%s", block)
	}
}

// The same entries and budget must always produce the same selection: no map
// iteration order or timestamp drift may change what a character receives.
func TestResolveBudgetIsDeterministic(t *testing.T) {
	makeEntry := func(id string) Entry {
		return withDefaults(Entry{
			ID:               id,
			Title:            id,
			Content:          strings.Repeat("x", 20),
			Enabled:          true,
			Scope:            ScopeWorld,
			TriggerPhrases:   []string{"sigil"},
			TriggerMode:      TriggerExact,
			TriggerDirection: TriggerBoth,
			TriggerFrequency: FrequencyAlways,
			PriorityLevel:    3,
		})
	}

	entries := []Entry{makeEntry("first"), makeEntry("second")}
	request := ResolveRequest{Message: "the sigil", MaxCharacters: 25}

	var previous string
	for i := 0; i < 5; i++ {
		resolved := Resolve(entries, request)
		if len(resolved) != 1 {
			t.Fatalf("run %d resolved %d entries with a 25-char budget: %+v", i, len(resolved), resolved)
		}
		if resolved[0].ID != "first" {
			t.Fatalf("run %d selected %s, want first", i, resolved[0].ID)
		}
		block := BuildPromptBlock(resolved)
		if i > 0 && block != previous {
			t.Fatalf("budget selection changed between runs:\n%s\n---\n%s", previous, block)
		}
		previous = block
	}
}
