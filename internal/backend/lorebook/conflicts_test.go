package lorebook

import (
    "strings"
    "testing"
)

func TestResolveConflictsEmpty(t *testing.T) {
    result := ResolveConflicts([]ResolvedEntry{}, ConflictHighestPriority)
    if len(result) != 0 {
        t.Fatal("empty input should produce empty output")
    }
}

func TestResolveConflictsSingle(t *testing.T) {
    entries := []ResolvedEntry{
        {
            Entry: Entry{ID: "single", Content: "content", PriorityLevel: 1},
        },
    }

    result := ResolveConflicts(entries, ConflictHighestPriority)
    if len(result) != 1 {
        t.Fatalf("single entry should pass through, got %d", len(result))
    }
}

func TestHighestPriorityStrategyLowNumberHigher(t *testing.T) {
    entries := []ResolvedEntry{
        {
            Entry:         Entry{ID: "low_num", PriorityLevel: 1, UpdatedAt: "2024-01-01"},
            MatchedPhrase: "test",
        },
        {
            Entry:         Entry{ID: "high_num", PriorityLevel: 5, UpdatedAt: "2024-01-01"},
            MatchedPhrase: "test",
        },
    }

    result := highestPriorityStrategy(entries)
    if len(result) != 1 {
        t.Fatalf("should select one entry, got %d", len(result))
    }
    if result[0].ID != "low_num" {
        t.Fatalf("lower priority number should win, got %s", result[0].ID)
    }
}

func TestHighestPriorityStrategyByDate(t *testing.T) {
    entries := []ResolvedEntry{
        {
            Entry:         Entry{ID: "old", PriorityLevel: 3, UpdatedAt: "2024-01-01"},
            MatchedPhrase: "test",
        },
        {
            Entry:         Entry{ID: "new", PriorityLevel: 3, UpdatedAt: "2024-12-31"},
            MatchedPhrase: "test",
        },
    }

    result := highestPriorityStrategy(entries)
    if len(result) != 1 {
        t.Fatalf("should select one entry, got %d", len(result))
    }
    // When priority is equal, newer should win as tiebreaker
    if result[0].ID != "new" {
        t.Fatalf("newer date should win as tiebreaker, got %s", result[0].ID)
    }
}

func TestNewerWinsStrategy(t *testing.T) {
    entries := []ResolvedEntry{
        {
            Entry:         Entry{ID: "old", UpdatedAt: "2024-01-01"},
            MatchedPhrase: "test",
        },
        {
            Entry:         Entry{ID: "new", UpdatedAt: "2024-12-31"},
            MatchedPhrase: "test",
        },
        {
            Entry:         Entry{ID: "middle", UpdatedAt: "2024-06-15"},
            MatchedPhrase: "test",
        },
    }

    result := newerWinsStrategy(entries)
    if len(result) != 1 {
        t.Fatalf("newer wins should select one entry, got %d", len(result))
    }
    if result[0].ID != "new" {
        t.Fatalf("newest entry should win, got %s", result[0].ID)
    }
}

func TestUseBothStrategy(t *testing.T) {
    entries := []ResolvedEntry{
        {
            Entry:         Entry{ID: "first", Content: "First."},
            MatchedPhrase: "test",
        },
        {
            Entry:         Entry{ID: "second", Content: "Second."},
            MatchedPhrase: "test",
        },
        {
            Entry:         Entry{ID: "third", Content: "Third."},
            MatchedPhrase: "test",
        },
    }

    result := useBothStrategy(entries)
    if len(result) != 3 {
        t.Fatalf("use both should keep all entries, got %d", len(result))
    }
}

func TestMergeSummariesStrategy(t *testing.T) {
    entries := []ResolvedEntry{
        {
            Entry:         Entry{ID: "part1", Content: "First part"},
            MatchedPhrase: "story",
        },
        {
            Entry:         Entry{ID: "part2", Content: "Second part"},
            MatchedPhrase: "story",
        },
    }

    result := mergeSummariesStrategy(entries)
    if len(result) != 1 {
        t.Fatalf("merge should create one entry, got %d", len(result))
    }

    // Should have separator
    if !strings.Contains(result[0].Content, "|") {
        t.Fatalf("merged content should contain separator, got: %s", result[0].Content)
    }

    // Should contain both contents
    if !strings.Contains(result[0].Content, "First part") {
        t.Fatal("merged content should contain first part")
    }
    if !strings.Contains(result[0].Content, "Second part") {
        t.Fatal("merged content should contain second part")
    }
}

func TestMergeSummariesEmptyContent(t *testing.T) {
    entries := []ResolvedEntry{
        {
            Entry:         Entry{ID: "empty", Content: "   "},
            MatchedPhrase: "test",
        },
        {
            Entry:         Entry{ID: "full", Content: "Content here"},
            MatchedPhrase: "test",
        },
    }

    result := mergeSummariesStrategy(entries)
    if len(result) != 1 {
        t.Fatalf("should create one merged entry, got %d", len(result))
    }

    // Should only have the non-empty content in final (plus separator)
    if !strings.Contains(result[0].Content, "Content here") {
        t.Fatal("merged should contain non-empty content")
    }
}

func TestGroupByTrigger(t *testing.T) {
    entries := []ResolvedEntry{
        {
            Entry:         Entry{ID: "a"},
            MatchedPhrase: "magic",
        },
        {
            Entry:         Entry{ID: "b"},
            MatchedPhrase: "spell",
        },
        {
            Entry:         Entry{ID: "c"},
            MatchedPhrase: "magic",
        },
    }

    groups := groupByTrigger(entries)
    if len(groups) != 2 {
        t.Fatalf("should have 2 groups, got %d", len(groups))
    }

    // Find magic group
    var magicGroup []ResolvedEntry
    for _, group := range groups {
        if group[0].MatchedPhrase == "magic" {
            magicGroup = group
        }
    }

    if len(magicGroup) != 2 {
        t.Fatalf("magic group should have 2 entries, got %d", len(magicGroup))
    }
}

func TestFindConflicts(t *testing.T) {
    entries := []ResolvedEntry{
        {
            Entry:         Entry{ID: "a"},
            MatchedPhrase: "trigger1",
        },
        {
            Entry:         Entry{ID: "b"},
            MatchedPhrase: "trigger1",
        },
        {
            Entry:         Entry{ID: "c"},
            MatchedPhrase: "trigger2",
        },
    }

    conflicts := FindConflicts(entries)
    if len(conflicts) != 1 {
        t.Fatalf("should find 1 conflict group, got %d", len(conflicts))
    }

    if len(conflicts["trigger1"]) != 2 {
        t.Fatalf("trigger1 conflict should have 2 entries, got %d", len(conflicts["trigger1"]))
    }

    // trigger2 should not be in conflicts (no conflict)
    if _, exists := conflicts["trigger2"]; exists {
        t.Fatal("single-entry trigger should not be in conflicts")
    }
}

func TestDetermineConflictRulePriority(t *testing.T) {
    entry := Entry{
        ConflictRule: string(ConflictMergeSummaries),
        Scope:        ScopeCharacter,
    }

    rule := DetermineConflictRule(entry, ConflictHighestPriority)
    if rule != ConflictMergeSummaries {
        t.Fatalf("entry rule should take priority, got %s", rule)
    }
}

func TestDetermineConflictRuleScopeDefault(t *testing.T) {
    entry := Entry{
        ConflictRule: "",
        Scope:        ScopeCharacter,
    }

    rule := DetermineConflictRule(entry, ConflictHighestPriority)
    if rule != ConflictNewerWins {
        t.Fatalf("character scope should default to newer wins, got %s", rule)
    }
}

func TestDetermineConflictRuleScenarioScope(t *testing.T) {
    entry := Entry{
        ConflictRule: "",
        Scope:        ScopeScenario,
    }

    rule := DetermineConflictRule(entry, ConflictHighestPriority)
    if rule != ConflictUseBoth {
        t.Fatalf("scenario scope should default to use both, got %s", rule)
    }
}

func TestDetermineConflictRuleGlobalDefault(t *testing.T) {
    entry := Entry{
        ConflictRule: "",
        Scope:        ScopeWorld,
    }

    rule := DetermineConflictRule(entry, ConflictMergeSummaries)
    if rule != ConflictMergeSummaries {
        t.Fatalf("world scope should use global default, got %s", rule)
    }
}

func TestIsValidConflictRule(t *testing.T) {
    tests := []struct {
        rule  ConflictRule
        valid bool
    }{
        {ConflictHighestPriority, true},
        {ConflictNewerWins, true},
        {ConflictUseBoth, true},
        {ConflictMergeSummaries, true},
        {ConflictRule("invalid"), false},
    }

    for _, tt := range tests {
        result := isValidConflictRule(tt.rule)
        if result != tt.valid {
            t.Fatalf("rule %q: expected %v, got %v", tt.rule, tt.valid, result)
        }
    }
}

func TestResolveConflictsWithEntryOverride(t *testing.T) {
    entries := []ResolvedEntry{
        {
            Entry: Entry{
                ID:           "override",
                PriorityLevel: 5,
                ConflictRule: string(ConflictNewerWins),
                UpdatedAt:    "2024-12-31",
            },
            MatchedPhrase: "test",
        },
        {
            Entry: Entry{
                ID:           "default",
                PriorityLevel: 1,
                ConflictRule: string(ConflictHighestPriority),
                UpdatedAt:    "2024-01-01",
            },
            MatchedPhrase: "test",
        },
    }

    // Even though "default" has higher priority, the group's rule is newer_wins
    // So "override" (newer) should win
    result := ResolveConflicts(entries, ConflictHighestPriority)
    if len(result) != 1 {
        t.Fatalf("should resolve to 1 entry, got %d", len(result))
    }
    if result[0].ID != "override" {
        t.Fatalf("newer_wins rule should prefer newer entry, got %s", result[0].ID)
    }
}
