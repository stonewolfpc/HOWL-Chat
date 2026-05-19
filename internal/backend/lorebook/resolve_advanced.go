/**
 * Advanced Resolution with Priority, Decay, Cooldown, and Budget
 *
 * Integrates Phase 2 anti-flooding features into the lorebook resolution pipeline.
 * This is the recommended entry point for memory-enabled chats.
 *
 * @package lorebook
 * @version 1.0.0
 */

package lorebook

import (
	"time"
)

// ResolveRequestAdvanced extends ResolveRequest with Phase 2 features
type ResolveRequestAdvanced struct {
	ResolveRequest
	Now           time.Time        // Current time for decay calculations
	CooldownMgr   *CooldownManager // For trigger frequency enforcement
	BudgetMgr     *BudgetManager   // For per-scope budget limits
	MinPriority   int              // Evict entries below this effective priority
}

// ResolveWithMemorySupport runs full resolution with anti-flooding protection
func ResolveWithMemorySupport(entries []Entry, req ResolveRequestAdvanced) []ResolvedEntry {
	if req.Now.IsZero() {
		req.Now = time.Now()
	}

	// Step 1: Pre-filter with priority decay and budget allocation
	entries = preFilterEntries(entries, req)

	// Step 2: Run standard resolution
	standardReq := req.ResolveRequest
	resolved := Resolve(entries, standardReq)

	// Step 3: Record triggers for cooldown tracking
	if req.CooldownMgr != nil {
		for _, r := range resolved {
			req.CooldownMgr.RecordTrigger(r.ID, req.Now)
			// Reinforce memory (boosts priority, resets decay)
			for i, e := range entries {
				if e.ID == r.ID {
					entries[i] = ReinforceMemory(e, req.Now)
					break
				}
			}
		}
	}

	return resolved
}

// preFilterEntries applies priority decay, cooldown, and budget before matching
func preFilterEntries(entries []Entry, req ResolveRequestAdvanced) []Entry {
	now := req.Now

	// Filter 1: Remove decayed entries (those with effective priority too low)
	var nonDecayed []Entry
	for _, e := range entries {
		effectivePriority := CalculateEffectivePriority(e, now)
		if effectivePriority <= req.MinPriority || req.MinPriority == 0 {
			nonDecayed = append(nonDecayed, e)
		}
	}
	entries = nonDecayed

	// Filter 2: Apply cooldown/frequency limits
	if req.CooldownMgr != nil {
		var available []Entry
		for _, e := range entries {
			if req.CooldownMgr.CanTrigger(e, now) {
				available = append(available, e)
			}
		}
		entries = available
	}

	// Filter 3: Apply budget limits per scope
	if req.BudgetMgr != nil {
		entries = req.BudgetMgr.AllocateWithinBudget(entries, now)
	}

	// Sort by effective priority (best first) for consistent results
	return SortByPriority(entries, now)
}

// GetMemorySystemStatus returns diagnostic info about memory system state
func GetMemorySystemStatus(entries []Entry, req ResolveRequestAdvanced) map[string]interface{} {
	now := req.Now
	if now.IsZero() {
		now = time.Now()
	}

	status := make(map[string]interface{})

	// Priority/Decay stats
	priorityDist := make(map[int]int)
	decayedCount := 0
	for _, e := range entries {
		p := CalculateEffectivePriority(e, now)
		priorityDist[p]++
		if p != e.PriorityLevel {
			decayedCount++
		}
	}
	status["priorityDistribution"] = priorityDist
	status["decayedEntries"] = decayedCount

	// Budget status
	if req.BudgetMgr != nil {
		status["budgetStatus"] = req.BudgetMgr.GetTotalBudget()
		// Calculate per-scope usage
		scopeCounts := make(map[BudgetScope]int)
		for _, e := range entries {
			scopeCounts[BudgetScope(e.Scope)]++
		}
		status["scopeCounts"] = scopeCounts
	}

	// Cooldown status
	if req.CooldownMgr != nil {
		triggered := 0
		for _, e := range entries {
			if !req.CooldownMgr.CanTrigger(e, now) {
				triggered++
			}
		}
		status["cooldownBlocked"] = triggered
	}

	return status
}
