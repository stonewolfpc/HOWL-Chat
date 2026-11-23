package chat

import (
	"testing"

	"howl-chat/internal/backend/lorebook"
	entities "howl-chat/internal/backend/types/entities"
)

func TestSessionAppendUserTurnDeduplicatesTheSameTurn(t *testing.T) {
	s := NewSession()
	s.AppendUserTurn("hello")
	s.AppendUserTurn("hello")
	if s.Len() != 1 {
		t.Fatalf("turn recorded %d times, want 1", s.Len())
	}
	s.AppendUserTurn("again")
	if s.Len() != 2 {
		t.Fatalf("a different turn was dropped: %d", s.Len())
	}
}

func TestSessionLoadKeepsAttributionAndResetsTriggers(t *testing.T) {
	world := &WorldContext{}
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
		TriggerFrequency: lorebook.FrequencyOncePerChat,
		PriorityLevel:    3,
	}})
	resolved := s.ResolveLore(lorebook.ResolveRequest{Message: "the sigil", Actor: lorebook.TriggerUser})
	if len(resolved) != 1 {
		t.Fatalf("resolved %d entries, want 1", len(resolved))
	}

	s.Load(&ChatData{
		World: world,
		History: []ChatHistory{
			{Role: "assistant", Sender: "Aria", Text: "line", CharacterID: "char_aria"},
			{Role: "user", Sender: "User", Text: "reply"},
		},
	})
	entries := s.Entries()
	if len(entries) != 2 {
		t.Fatalf("loaded %d turns, want 2", len(entries))
	}
	if entries[0].Character != "Aria" || entries[0].CharacterID != "char_aria" {
		t.Fatalf("assistant attribution lost: %+v", entries[0])
	}
	if entries[1].Sender != "User" {
		t.Fatalf("user sender lost: %+v", entries[1])
	}
	if s.World() != world {
		t.Fatal("world not adopted")
	}
	// The trigger bookkeeping belongs to the conversation that produced it: a load
	// starts a different one, so a once-per-chat entry may fire again.
	if again := s.ResolveLore(lorebook.ResolveRequest{Message: "the sigil", Actor: lorebook.TriggerUser}); len(again) != 1 {
		t.Fatalf("trigger state survived a load: %d entries", len(again))
	}
}

func TestSessionProjectKeepsStoredHistoryWhenTheSessionIsEmpty(t *testing.T) {
	stored := &ChatData{History: []ChatHistory{{Role: "user", Text: "on disk"}}}
	NewSession().Project(stored)
	if len(stored.History) != 1 || stored.History[0].Text != "on disk" {
		t.Fatalf("an unloaded session rewrote history: %+v", stored.History)
	}

	s := NewSession()
	s.AppendEntry(SessionEntry{Role: "user", Content: "new line"})
	s.Project(stored)
	if len(stored.History) != 1 || stored.History[0].Text != "new line" {
		t.Fatalf("a loaded session did not write history: %+v", stored.History)
	}
}

func TestSessionProjectRecomputesMembershipTracking(t *testing.T) {
	s := NewSession()
	s.SetCharacters([]*CharacterContext{
		{Character: entities.Character{ID: "char_aria", Name: "Aria"}},
	})
	stored := &ChatData{Meta: ChatMeta{RemovedCharacters: []string{"char_aria", "char_gone"}}}
	s.Project(stored)
	if len(stored.Meta.ActiveCharacters) != 1 || stored.Meta.ActiveCharacters[0] != "char_aria" {
		t.Fatalf("active characters = %v", stored.Meta.ActiveCharacters)
	}
	if len(stored.Meta.RemovedCharacters) != 1 || stored.Meta.RemovedCharacters[0] != "char_gone" {
		t.Fatalf("removed characters = %v, want only the one that is no longer active", stored.Meta.RemovedCharacters)
	}
}
