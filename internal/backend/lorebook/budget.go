/**
 * Memory Budget Manager
 *
 * Enforces per-scope memory limits to prevent context flooding.
 * Critical for MVP - prevents LLM context window exhaustion.
 *
 * @package lorebook
 * @version 1.0.0
 */

package lorebook

import (
	"sort"
	"time"
)

// BudgetScope defines memory scope for budget allocation
type BudgetScope string

const (
	BudgetGlobal    BudgetScope = "global"
	BudgetWorld     BudgetScope = "world"
	BudgetScenario  BudgetScope = "scenario"
	BudgetCharacter BudgetScope = "character"
	BudgetChat      BudgetScope = "chat"
)

// BudgetLimits defines default max memories per scope
var BudgetLimits = map[BudgetScope]int{
	BudgetGlobal:    5,  // Rare, high-impact only
	BudgetWorld:     2,  // World-level context is heavy
	BudgetScenario:  3,  // Scenario context
	BudgetCharacter: 6,  // Per-character allowance
	BudgetChat:      10, // Session-specific memories
}

// BudgetManager enforces memory limits per scope
type BudgetManager struct {
	limits    map[BudgetScope]int
	allocated map[string]bool // entryIDs currently in budget
}

// NewBudgetManager creates a budget manager with default limits
func NewBudgetManager() *BudgetManager {
	return &BudgetManager{
		limits:    BudgetLimits,
		allocated: make(map[string]bool),
	}
}

// NewBudgetManagerWithLimits creates a budget manager with custom limits
func NewBudgetManagerWithLimits(limits map[BudgetScope]int) *BudgetManager {
	bm := NewBudgetManager()
	for scope, limit := range limits {
		bm.limits[scope] = limit
	}
	return bm
}

// AllocateWithinBudget selects entries that fit within scope budgets
// Prioritizes by: pinned > high priority > recently triggered > age
func (bm *BudgetManager) AllocateWithinBudget(entries []Entry, now time.Time) []Entry {
	// Group by scope
	byScope := make(map[BudgetScope][]Entry)
	for _, e := range entries {
		scope := BudgetScope(e.Scope)
		byScope[scope] = append(byScope[scope], e)
	}

	var selected []Entry

	// Process each scope with its budget
	for scope, scopeEntries := range byScope {
		limit := bm.getLimit(scope)
		scopeSelected := bm.selectTopEntries(scopeEntries, limit, now)
		selected = append(selected, scopeSelected...)
	}

	return selected
}

// selectTopEntries picks best entries up to limit using priority scoring
func (bm *BudgetManager) selectTopEntries(entries []Entry, limit int, now time.Time) []Entry {
	if len(entries) <= limit {
		return entries
	}

	// Score all entries (lower score = better)
	type scoredEntry struct {
		entry Entry
		score float64
	}

	var scored []scoredEntry
	for _, e := range entries {
		score := ScoreForBudget(e, now)
		scored = append(scored, scoredEntry{entry: e, score: score})
	}

	// Sort by score ascending (best first)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score < scored[j].score
	})

	// Take top N
	var result []Entry
	for i := 0; i < limit && i < len(scored); i++ {
		result = append(result, scored[i].entry)
		bm.allocated[scored[i].entry.ID] = true
	}

	return result
}

// getLimit returns the budget limit for a scope
func (bm *BudgetManager) getLimit(scope BudgetScope) int {
	if limit, ok := bm.limits[scope]; ok {
		return limit
	}
	return 3 // Default fallback
}

// GetBudgetStatus returns current allocation status per scope
func (bm *BudgetManager) GetBudgetStatus(entries []Entry) map[BudgetScope]BudgetStatus {
	status := make(map[BudgetScope]BudgetStatus)

	// Count entries per scope
	counts := make(map[BudgetScope]int)
	for _, e := range entries {
		counts[BudgetScope(e.Scope)]++
	}

	// Build status for each scope with a limit
	for scope, limit := range bm.limits {
		used := counts[scope]
		status[scope] = BudgetStatus{
			Scope:     scope,
			Limit:     limit,
			Used:      used,
			Available: limit - used,
			Percent:   float64(used) / float64(limit) * 100,
		}
	}

	return status
}

// BudgetStatus shows allocation for a scope
type BudgetStatus struct {
	Scope     BudgetScope
	Limit     int
	Used      int
	Available int
	Percent   float64
}

// IsOverBudget checks if adding another entry would exceed limit
func (bm *BudgetManager) IsOverBudget(scope BudgetScope, currentCount int) bool {
	limit := bm.getLimit(scope)
	return currentCount >= limit
}

// EmergencyPrune removes lowest priority entries when severely over budget
// Returns pruned entries and remaining entries
func (bm *BudgetManager) EmergencyPrune(entries []Entry, now time.Time, maxTotal int) (pruned, remaining []Entry) {
	if len(entries) <= maxTotal {
		return []Entry{}, entries
	}

	// Score and sort all
	type scoredEntry struct {
		entry Entry
		score float64
	}

	var scoredList []scoredEntry
	for _, e := range entries {
		scoredList = append(scoredList, scoredEntry{entry: e, score: ScoreForBudget(e, now)})
	}

	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].score < scoredList[j].score
	})

	// Keep top maxTotal, prune the rest
	keepCount := maxTotal
	for i, s := range scoredList {
		if i < keepCount {
			remaining = append(remaining, s.entry)
		} else {
			pruned = append(pruned, s.entry)
		}
	}

	return pruned, remaining
}

// SetScopeLimit updates budget for a specific scope
func (bm *BudgetManager) SetScopeLimit(scope BudgetScope, limit int) {
	bm.limits[scope] = limit
}

// GetTotalBudget calculates total memory slots across all scopes
func (bm *BudgetManager) GetTotalBudget() int {
	total := 0
	for _, limit := range bm.limits {
		total += limit
	}
	return total
}

// DefaultBudgetConfig returns standard MVP budget allocation
func DefaultBudgetConfig() map[BudgetScope]int {
	return map[BudgetScope]int{
		BudgetGlobal:    5,
		BudgetWorld:     2,
		BudgetScenario:  3,
		BudgetCharacter: 6,
		BudgetChat:      10,
	}
}

// ConservativeBudgetConfig returns stricter limits for constrained contexts
func ConservativeBudgetConfig() map[BudgetScope]int {
	return map[BudgetScope]int{
		BudgetGlobal:    3,
		BudgetWorld:     1,
		BudgetScenario:  2,
		BudgetCharacter: 3,
		BudgetChat:      5,
	}
}
