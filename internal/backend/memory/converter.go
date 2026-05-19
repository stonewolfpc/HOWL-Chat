/**
 * Memory Converter
 *
 * Converts MemoryRecord types to lorebook.Entry for resolution pipeline.
 * Auto-generates trigger phrases from title/summary.
 *
 * @package memory
 * @version 1.0.0
 */

package memory

import (
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/backend/types"
)

// ToLorebookEntry converts a MemoryRecord to a lorebook.Entry with auto-generated triggers
func ToLorebookEntry(record types.MemoryRecord) lorebook.Entry {
	entry := lorebook.Entry{
		ID:                record.ID,
		Title:             record.Title,
		Description:       record.Summary,
		Content:           record.Content,
		Scope:             lorebook.Scope(record.Scope),
		OwnerID:           record.OwnerID,
		Tags:              record.Tags,
		Enabled:           record.Status == "active",
		PriorityLevel:     calculatePriority(record),
		DecayRate:         calculateDecay(record),
		TriggerFrequency:  calculateFrequency(record),
		TriggerMode:       lorebook.TriggerLoose,
		TriggerDirection:  lorebook.TriggerBoth,
		TriggerPhrases:    generateTriggerPhrases(record),
		SecondaryTriggers: generateSecondaryTriggers(record),
		Automated:         record.Source == "auto",
		MemoryType:        determineMemoryType(record),
		TriggerConfidence: "medium",
		VisibleToUser:     true,
		VisibleToOthers:   true,
		MaxLength:         "medium",
		ContextBudget:     700,
		CreatedAt:         record.CreatedAt,
		LastTriggeredAt:   record.UpdatedAt,
		UpdatedAt:         record.UpdatedAt,
	}

	return entry
}

// ToLorebookEntries converts multiple MemoryRecords
func ToLorebookEntries(records []types.MemoryRecord) []lorebook.Entry {
	entries := make([]lorebook.Entry, 0, len(records))
	for _, record := range records {
		if record.Status != "archived" {
			entries = append(entries, ToLorebookEntry(record))
		}
	}
	return entries
}

// calculatePriority assigns priority based on source and content
func calculatePriority(record types.MemoryRecord) int {
	if record.Source == "auto" && hasTag(record.Tags, "death") {
		return 1 // Highest priority for deaths
	}
	if record.Pinned {
		return 2
	}
	return 3 // Default medium priority
}

// calculateDecay assigns decay rate
func calculateDecay(record types.MemoryRecord) string {
	if record.Source == "auto" && hasTag(record.Tags, "death") {
		return "none" // Death memories never decay
	}
	return "medium"
}

// calculateFrequency assigns trigger frequency
func calculateFrequency(record types.MemoryRecord) lorebook.TriggerFrequency {
	if record.Source == "auto" && hasTag(record.Tags, "death") {
		return lorebook.FrequencyOncePerChat
	}
	return lorebook.FrequencyAlways
}

// determineMemoryType categorizes the memory
func determineMemoryType(record types.MemoryRecord) string {
	for _, tag := range record.Tags {
		if tag == "death" || tag == "event" || tag == "promise" || tag == "location" {
			return tag
		}
	}
	return "general"
}

// hasTag checks if a tag exists in the slice
func hasTag(tags []string, target string) bool {
	for _, tag := range tags {
		if tag == target {
			return true
		}
	}
	return false
}
