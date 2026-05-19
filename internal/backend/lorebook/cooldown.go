/**
 * Trigger Cooldown and Rotation Manager
 *
 * Prevents memory spam by tracking trigger frequency and enforcing cooldowns.
 * Implements trigger rotation to ensure variety in triggered memories.
 *
 * @package lorebook
 * @version 1.0.0
 */

package lorebook

import (
	"time"
)

// Frequency aliases for local use (maps to TriggerFrequency constants)
const (
	FreqAlways       = FrequencyAlways
	FreqOncePerChat  = FrequencyOncePerChat
	FreqOncePerScene = FrequencyOncePerScene
	FreqOncePerDay   = TriggerFrequency("once_per_day")
	FreqCooldown     = TriggerFrequency("cooldown")
)

// CooldownManager tracks trigger state for all entries
type CooldownManager struct {
	triggerHistory map[string]time.Time // entryID -> last triggered
	triggerCounts  map[string]int       // entryID -> times triggered this session
	sceneMarkers   []time.Time          // When scenes were marked
}

// NewCooldownManager creates a new cooldown tracker
func NewCooldownManager() *CooldownManager {
	return &CooldownManager{
		triggerHistory: make(map[string]time.Time),
		triggerCounts:  make(map[string]int),
		sceneMarkers:   []time.Time{time.Now()}, // Start with initial scene
	}
}

// CanTrigger checks if an entry is eligible to trigger now
func (cm *CooldownManager) CanTrigger(entry Entry, now time.Time) bool {
	freq := TriggerFrequency(entry.TriggerFrequency)
	lastTriggered, hasHistory := cm.triggerHistory[entry.ID]
	_ = cm.triggerCounts[entry.ID] // Used for reinforcement elsewhere

	switch freq {
	case FreqAlways:
		return true

	case FreqOncePerChat:
		return !hasHistory

	case FreqOncePerScene:
		if !hasHistory {
			return true
		}
		// Check if a new scene started since last trigger
		return cm.isNewSceneSince(lastTriggered)

	case FreqOncePerDay:
		if !hasHistory {
			return true
		}
		return now.Sub(lastTriggered).Hours() >= 24

	case FreqCooldown:
		// Default cooldown period: 5 minutes
		cooldown := 5 * time.Minute
		if !hasHistory {
			return true
		}
		return now.Sub(lastTriggered) >= cooldown

	default:
		// Unknown frequency defaults to once per chat
		return !hasHistory
	}
}

// RecordTrigger marks an entry as triggered
func (cm *CooldownManager) RecordTrigger(entryID string, now time.Time) {
	cm.triggerHistory[entryID] = now
	cm.triggerCounts[entryID]++
}

// MarkScene notes a scene transition, resetting scene-scoped frequencies
func (cm *CooldownManager) MarkScene(now time.Time) {
	cm.sceneMarkers = append(cm.sceneMarkers, now)
}

// GetTriggerCount returns how many times an entry has triggered
func (cm *CooldownManager) GetTriggerCount(entryID string) int {
	return cm.triggerCounts[entryID]
}

// GetLastTriggered returns when an entry last triggered (zero if never)
func (cm *CooldownManager) GetLastTriggered(entryID string) time.Time {
	if t, ok := cm.triggerHistory[entryID]; ok {
		return t
	}
	return time.Time{}
}

// isNewSceneSince checks if any scene marker is newer than the given time
func (cm *CooldownManager) isNewSceneSince(t time.Time) bool {
	for _, marker := range cm.sceneMarkers {
		if marker.After(t) {
			return true
		}
	}
	return false
}

// RotationSet manages which entries to trigger when multiple match
type RotationSet struct {
	cooldownMgr   *CooldownManager
	priorityQueue []string // entryIDs sorted by priority
	lastRotation  time.Time
	rotationIndex int
}

// NewRotationSet creates a rotation manager for triggered entries
func NewRotationSet(cm *CooldownManager) *RotationSet {
	return &RotationSet{
		cooldownMgr:   cm,
		priorityQueue: []string{},
		lastRotation:  time.Now(),
		rotationIndex: 0,
	}
}

// SelectFromMultiple chooses which entry to trigger when multiple match
// Implements round-robin with priority weighting
func (rs *RotationSet) SelectFromMultiple(entries []Entry, now time.Time) Entry {
	if len(entries) == 0 {
		return Entry{}
	}
	if len(entries) == 1 {
		return entries[0]
	}

	// Filter to only entries that can actually trigger
	var available []Entry
	for _, e := range entries {
		if rs.cooldownMgr.CanTrigger(e, now) {
			available = append(available, e)
		}
	}

	if len(available) == 0 {
		return Entry{}
	}
	if len(available) == 1 {
		return available[0]
	}

	// Sort by priority (lower = better)
	// Simple selection: pick highest priority, with rotation for ties
	sorted := SortByPriority(available, now)

	// If top entries have same priority, rotate among them
	topPriority := CalculateEffectivePriority(sorted[0], now)
	var tied []Entry
	for _, e := range sorted {
		if CalculateEffectivePriority(e, now) == topPriority {
			tied = append(tied, e)
		}
	}

	if len(tied) == 1 {
		return tied[0]
	}

	// Rotate through tied entries
	rs.rotationIndex = (rs.rotationIndex + 1) % len(tied)
	return tied[rs.rotationIndex]
}

// EntryCooldownInfo provides human-readable cooldown status
type EntryCooldownInfo struct {
	CanTrigger     bool
	LastTriggered  time.Time
	NextTrigger    time.Time
	TimesTriggered int
	Reason         string
}

// GetCooldownInfo returns detailed cooldown status for an entry
func (cm *CooldownManager) GetCooldownInfo(entry Entry, now time.Time) EntryCooldownInfo {
	info := EntryCooldownInfo{
		CanTrigger:     cm.CanTrigger(entry, now),
		LastTriggered:  cm.GetLastTriggered(entry.ID),
		TimesTriggered: cm.GetTriggerCount(entry.ID),
	}

	if !info.CanTrigger {
		freq := TriggerFrequency(entry.TriggerFrequency)
		switch freq {
		case FreqOncePerChat:
			info.Reason = "Already triggered in this chat"
		case FreqOncePerScene:
			info.Reason = "Already triggered this scene"
			info.NextTrigger = cm.getNextSceneTime()
		case FreqOncePerDay:
			info.Reason = "Daily cooldown"
			info.NextTrigger = info.LastTriggered.Add(24 * time.Hour)
		case FreqCooldown:
			cooldown := 5 * time.Minute
			info.Reason = "In cooldown period"
			info.NextTrigger = info.LastTriggered.Add(cooldown)
		default:
			info.Reason = "Frequency limited"
		}
	}

	return info
}

// getNextSceneTime returns a placeholder for when next scene might occur
// In practice, scenes are marked manually or by heuristics
func (cm *CooldownManager) getNextSceneTime() time.Time {
	return time.Now().Add(30 * time.Minute) // Assume 30 min estimate
}
