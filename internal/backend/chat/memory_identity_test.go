package chat

import (
	"testing"

	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/backend/types"
)

func TestSnapshotFilterDropsUnmarkedCharacterMemory(t *testing.T) {
	entries := []lorebook.Entry{
		{ID: "legacy-name", Scope: lorebook.ScopeCharacter, OwnerID: "char_aria"},
		{ID: "stable", Scope: lorebook.ScopeCharacter, OwnerID: "char_aria", OwnerKind: types.OwnerKindCharacterID},
		{ID: "foreign", Scope: lorebook.ScopeCharacter, OwnerID: "char_other", OwnerKind: types.OwnerKindCharacterID},
		{ID: "world", Scope: lorebook.ScopeWorld, OwnerID: "world_1"},
		{ID: "scenario", Scope: lorebook.ScopeScenario, OwnerID: "scenario_1"},
	}

	filtered := filterStableCharacterMemoryEntries(entries, []string{"char_aria"})
	if len(filtered) != 3 {
		t.Fatalf("filtered %d entries, want stable + world + scenario: %+v", len(filtered), filtered)
	}
	seen := map[string]bool{}
	for _, entry := range filtered {
		seen[entry.ID] = true
	}
	if seen["legacy-name"] {
		t.Fatal("a legacy name-owned character snapshot reached the prompt pool")
	}
	if seen["foreign"] {
		t.Fatal("a foreign character snapshot reached the prompt pool")
	}
	if !seen["stable"] || !seen["world"] || !seen["scenario"] {
		t.Fatalf("safe entries were removed: %+v", seen)
	}
}
