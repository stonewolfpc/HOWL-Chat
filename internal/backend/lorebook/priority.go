/**
 * Priority and Decay Engine
 *
 * Manages memory priority levels, decay rates, and reinforcement.
 * Prevents memory flooding by self-pruning low-priority, aged memories.
 *
 * @package lorebook
 * @version 1.0.0
 */

package lorebook

import (
	"math"
	"time"
)

// DecayRate defines how quickly a memory loses priority over time
type DecayRate string

const (
	DecayNone     DecayRate = "none"     // Never decay (pinned, critical)
	DecaySlow     DecayRate = "slow"     // Decay over months
	DecayNormal   DecayRate = "normal"   // Decay over weeks
	DecayFast     DecayRate = "fast"     // Decay over days
	DecayVeryFast DecayRate = "veryfast" // Decay over hours
)

// PriorityLevel constants for semantic clarity
const (
	PriorityCritical = 1  // Deaths, major plot points - never evict
	PriorityHigh     = 3  // Important character development
	PriorityNormal   = 5  // Standard memories
	PriorityLow      = 7  // Minor details
	PriorityTrivial  = 10 // Filler, first to evict
)

// PriorityConfig holds priority calculation settings
type PriorityConfig struct {
	BasePriority  int
	DecayRate     DecayRate
	LastTriggered time.Time
	TriggerCount  int
}

// CalculateEffectivePriority computes current priority considering decay
func CalculateEffectivePriority(entry Entry, now time.Time) int {
	// If no decay or pinned, return base priority
	if entry.DecayRate == string(DecayNone) || entry.DecayRate == "never" {
		return entry.PriorityLevel
	}

	// Calculate age in days
	age := now.Sub(parseTime(entry.LastTriggeredAt)).Hours() / 24
	if age < 0 {
		age = 0
	}

	// Get decay factor based on rate
	decayFactor := getDecayFactor(entry.DecayRate)

	// Calculate decay penalty
	penalty := int(age * decayFactor)

	// Effective priority increases as base gets worse with decay
	effective := entry.PriorityLevel + penalty

	// Reinforcement bonus: each trigger reduces effective priority (improves it)
	reinforcementBonus := entry.TriggerCount / 3 // Every 3 triggers = -1 priority
	effective -= reinforcementBonus

	// Clamp to valid range
	if effective < 1 {
		effective = 1
	}
	if effective > 10 {
		effective = 10
	}

	return effective
}

// ShouldEvict determines if a memory should be removed based on age and priority
func ShouldEvict(entry Entry, now time.Time, minPriority int) bool {
	// Never evict pinned or critical memories
	if entry.DecayRate == string(DecayNone) || entry.DecayRate == "never" {
		return false
	}

	effectivePriority := CalculateEffectivePriority(entry, now)

	// Evict if priority has decayed below threshold
	return effectivePriority > minPriority
}

// ReinforceMemory boosts priority and resets decay when memory triggers
func ReinforceMemory(entry Entry, now time.Time) Entry {
	entry.TriggerCount++
	entry.LastTriggeredAt = now.Format(time.RFC3339)

	// Boost priority slightly (lower number = higher priority)
	if entry.PriorityLevel > 1 {
		boost := 1
		if entry.TriggerCount > 5 {
			boost = 2 // Stronger boost for frequently used memories
		}
		entry.PriorityLevel -= boost
		if entry.PriorityLevel < 1 {
			entry.PriorityLevel = 1
		}
	}

	return entry
}

// SortByPriority sorts entries by effective priority (best first)
func SortByPriority(entries []Entry, now time.Time) []Entry {
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)

	// Simple bubble sort for clarity (small arrays expected)
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			p1 := CalculateEffectivePriority(sorted[j], now)
			p2 := CalculateEffectivePriority(sorted[j+1], now)

			// Higher priority (lower number) comes first
			if p1 > p2 {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

// FilterAndPrioritize removes evicted memories and sorts by priority
func FilterAndPrioritize(entries []Entry, now time.Time, minPriority int) []Entry {
	var filtered []Entry

	for _, entry := range entries {
		if !ShouldEvict(entry, now, minPriority) {
			filtered = append(filtered, entry)
		}
	}

	return SortByPriority(filtered, now)
}

// getDecayFactor returns daily decay rate for a DecayRate
func getDecayFactor(rate string) float64 {
	switch DecayRate(rate) {
	case DecaySlow:
		return 0.1 // 0.1 priority per day
	case DecayNormal:
		return 0.3 // 0.3 priority per day
	case DecayFast:
		return 0.7 // 0.7 priority per day
	case DecayVeryFast:
		return 2.0 // 2.0 priority per day
	default:
		return 0.0 // No decay
	}
}

// parseTime safely parses RFC3339 time string
func parseTime(s string) time.Time {
	if s == "" {
		return time.Now()
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now()
	}
	return t
}

// AgeInDays returns how many days since memory was created/triggered
func AgeInDays(entry Entry, now time.Time) float64 {
	created := parseTime(entry.CreatedAt)
	age := now.Sub(created).Hours() / 24
	if age < 0 {
		return 0
	}
	return age
}

// ScoreForBudget calculates a composite score for budget allocation
// Lower score = more important, should be included
func ScoreForBudget(entry Entry, now time.Time) float64 {
	priority := float64(CalculateEffectivePriority(entry, now))

	// Age penalty (older = slightly less relevant unless reinforced)
	age := AgeInDays(entry, now)
	agePenalty := math.Log1p(age) * 0.5

	// Reinforcement bonus (frequently triggered = more relevant)
	reinforcementBonus := float64(entry.TriggerCount) * 0.3

	score := priority + agePenalty - reinforcementBonus

	// Pinned entries get massive bonus
	if isPinned(entry) {
		score -= 5.0
	}

	return score
}

func isPinned(entry Entry) bool {
	// Check various ways pinned might be stored
	// The Entry struct doesn't have Pinned directly, but may have it in metadata
	return entry.DecayRate == string(DecayNone) || entry.DecayRate == "never"
}
