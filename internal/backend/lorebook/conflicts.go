package lorebook

import (
	"sort"
	"strings"
)

type ConflictRule string

const (
	ConflictHighestPriority ConflictRule = "highest_priority"
	ConflictNewerWins       ConflictRule = "newer_wins"
	ConflictUseBoth         ConflictRule = "use_both"
	ConflictMergeSummaries  ConflictRule = "merge_summaries"
)

type ConflictResolution struct {
	Rule           ConflictRule
	ConflictingIDs []string
	Resolution     []ResolvedEntry
}

// ResolveConflicts handles multiple entries that match the same trigger
// Only applies conflict resolution if ConflictRule is explicitly set on entries
func ResolveConflicts(matched []ResolvedEntry, defaultRule ConflictRule) []ResolvedEntry {
	if len(matched) <= 1 {
		return matched
	}

	// Group by trigger phrase
	groups := groupByTrigger(matched)
	var result []ResolvedEntry

	for _, group := range groups {
		if len(group) == 1 {
			result = append(result, group[0])
			continue
		}

		// Only apply conflict resolution if entries have conflicting rules
		hasConflictRules := false
		for _, entry := range group {
			if entry.ConflictRule != "" {
				hasConflictRules = true
				break
			}
		}

		// If no explicit conflict rules, just return all entries
		if !hasConflictRules {
			result = append(result, group...)
			continue
		}

		// Determine which rule to use (entry's rule overrides default)
		rule := defaultRule
		if group[0].ConflictRule != "" {
			rule = ConflictRule(group[0].ConflictRule)
		}

		resolved := applyConflictRule(group, rule)
		result = append(result, resolved...)
	}

	return result
}

// hasConflict checks if entries in a group are actually conflicting
// (different IDs, same trigger phrase)
func hasConflict(group []ResolvedEntry) bool {
	if len(group) <= 1 {
		return false
	}

	// Check if all are from same entry (not a conflict)
	firstID := group[0].ID
	for _, entry := range group[1:] {
		if entry.ID != firstID {
			return true
		}
	}
	return false
}

// groupByTrigger groups entries by their matched phrase
func groupByTrigger(entries []ResolvedEntry) [][]ResolvedEntry {
	groups := make(map[string][]ResolvedEntry)
	var keys []string

	for _, entry := range entries {
		key := entry.MatchedPhrase
		if _, exists := groups[key]; !exists {
			keys = append(keys, key)
		}
		groups[key] = append(groups[key], entry)
	}

	var result [][]ResolvedEntry
	for _, key := range keys {
		result = append(result, groups[key])
	}
	return result
}

// applyConflictRule applies the specified conflict resolution rule
func applyConflictRule(conflicting []ResolvedEntry, rule ConflictRule) []ResolvedEntry {
	switch rule {
	case ConflictNewerWins:
		return newerWinsStrategy(conflicting)
	case ConflictUseBoth:
		return useBothStrategy(conflicting)
	case ConflictMergeSummaries:
		return mergeSummariesStrategy(conflicting)
	default:
		return highestPriorityStrategy(conflicting)
	}
}

// highestPriorityStrategy keeps only the highest priority entry
func highestPriorityStrategy(conflicting []ResolvedEntry) []ResolvedEntry {
	if len(conflicting) == 0 {
		return nil
	}

	sort.SliceStable(conflicting, func(i, j int) bool {
		if conflicting[i].PriorityLevel != conflicting[j].PriorityLevel {
			return conflicting[i].PriorityLevel < conflicting[j].PriorityLevel
		}
		return conflicting[i].UpdatedAt > conflicting[j].UpdatedAt
	})

	return []ResolvedEntry{conflicting[0]}
}

// newerWinsStrategy keeps the most recently updated entry
func newerWinsStrategy(conflicting []ResolvedEntry) []ResolvedEntry {
	if len(conflicting) == 0 {
		return nil
	}

	sort.SliceStable(conflicting, func(i, j int) bool {
		return conflicting[i].UpdatedAt > conflicting[j].UpdatedAt
	})

	return []ResolvedEntry{conflicting[0]}
}

// useBothStrategy keeps all conflicting entries
func useBothStrategy(conflicting []ResolvedEntry) []ResolvedEntry {
	return conflicting
}

// mergeSummariesStrategy merges content from conflicting entries
func mergeSummariesStrategy(conflicting []ResolvedEntry) []ResolvedEntry {
	if len(conflicting) == 0 {
		return nil
	}

	merged := conflicting[0]
	var contents []string
	var ids []string

	for _, entry := range conflicting {
		if strings.TrimSpace(entry.Content) != "" {
			contents = append(contents, entry.Content)
		}
		ids = append(ids, entry.ID)
	}

	merged.Content = strings.Join(contents, " | ")
	merged.Title = "Merged: " + strings.Join(ids, ", ")

	return []ResolvedEntry{merged}
}

// DetermineConflictRule selects the appropriate conflict rule
// Priority: entry rule > scope-based default > global default
func DetermineConflictRule(entry Entry, globalDefault ConflictRule) ConflictRule {
	if entry.ConflictRule != "" {
		rule := ConflictRule(entry.ConflictRule)
		if isValidConflictRule(rule) {
			return rule
		}
	}

	// Scope-based defaults
	switch entry.Scope {
	case ScopeCharacter:
		return ConflictNewerWins
	case ScopeScenario:
		return ConflictUseBoth
	default:
		return globalDefault
	}
}

func isValidConflictRule(rule ConflictRule) bool {
	switch rule {
	case ConflictHighestPriority, ConflictNewerWins, ConflictUseBoth, ConflictMergeSummaries:
		return true
	default:
		return false
	}
}

// FindConflicts identifies entries that would conflict with each other
func FindConflicts(entries []ResolvedEntry) map[string][]ResolvedEntry {
	conflicts := make(map[string][]ResolvedEntry)

	for _, entry := range entries {
		key := entry.MatchedPhrase
		conflicts[key] = append(conflicts[key], entry)
	}

	// Filter to only actual conflicts
	result := make(map[string][]ResolvedEntry)
	for key, group := range conflicts {
		if len(group) > 1 {
			result[key] = group
		}
	}

	return result
}
