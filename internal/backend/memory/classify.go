/**
 * Event Classification and Candidate Generation
 *
 * Converts raw EventDetectionResult into MemoryCandidate objects.
 * Handles scope assignment, title/summary generation, tag creation,
 * and deduplication of overlapping detections.
 *
 * @package memory
 * @version 1.0.0
 */

package memory

import (
	"fmt"
	"strings"

	"howl-chat/internal/backend/types"
)

// ClassifyAndBuildCandidate converts a detection result into a reviewable candidate
func ClassifyAndBuildCandidate(result types.EventDetectionResult, sourceMessage string, identities []CharacterIdentity, chatID string) types.MemoryCandidate {
	scope := inferScope(result.EventType)
	ownerID, ownerKind := inferOwner(result, identities, chatID)
	title := generateTitle(result)
	summary := generateSummary(result)
	content := buildContent(result, sourceMessage)
	tags := generateTags(result, identities)

	return types.MemoryCandidate{
		ID:            types.NewMemoryID(),
		EventType:     result.EventType,
		Title:         title,
		Summary:       summary,
		Content:       content,
		SourceMessage: sourceMessage,
		Scope:         scope,
		OwnerID:       ownerID,
		OwnerKind:     ownerKind,
		Tags:          tags,
		Confidence:    result.Confidence,
		Status:        "pending",
		CreatedAt:     types.NewMemoryID()[4:], // Reuse timestamp format
	}
}

// inferScope maps event types to default lorebook scopes
func inferScope(eventType types.EventType) string {
	switch eventType {
	case types.EventDeath, types.EventCharacterChange:
		return "character"
	case types.EventDiscovery, types.EventLocation, types.EventItem:
		return "scenario"
	case types.EventBetrayal, types.EventRomance:
		return "character"
	case types.EventConflict:
		return "scenario"
	// All micro-memories are personal/character-scoped
	case types.EventPreference, types.EventLikeDislike, types.EventEmotion,
		types.EventAssociation, types.EventAesthetic, types.EventSensory,
		types.EventSocial, types.EventFirst, types.EventShared, types.EventPromise:
		return "character"
	default:
		return "chat"
	}
}

// inferOwner determines which entity owns this memory. A display name may match
// a character, but ownership is always the stable ID. If a name matches more
// than one identity (duplicates, or several characters mentioned), the memory is
// chat-owned rather than guessed.
func inferOwner(result types.EventDetectionResult, identities []CharacterIdentity, chatID string) (string, string) {
	lowerContext := strings.ToLower(result.Context)
	matches := map[string]bool{}
	for _, identity := range identities {
		id := strings.TrimSpace(identity.ID)
		name := strings.TrimSpace(identity.Name)
		if id == "" || name == "" {
			continue
		}
		if strings.Contains(lowerContext, strings.ToLower(name)) {
			matches[id] = true
		}
	}
	if len(matches) == 1 {
		for id := range matches {
			return id, types.OwnerKindCharacterID
		}
	}

	// Default to chat session for unowned or ambiguous events.
	if chatID != "" {
		return chatID, types.OwnerKindChat
	}
	return "default", types.OwnerKindChat
}

// generateTitle creates a human-readable title from event type and evidence
func generateTitle(result types.EventDetectionResult) string {
	switch result.EventType {
	case types.EventDeath:
		return "Death Event"
	case types.EventBetrayal:
		return "Betrayal"
	case types.EventRomance:
		return "Romantic Development"
	case types.EventDiscovery:
		return "Discovery"
	case types.EventConflict:
		return "Conflict"
	case types.EventCharacterChange:
		return "Character Change"
	case types.EventLocation:
		return "Location Event"
	case types.EventItem:
		return "Item Acquisition"
	// Micro-memory titles
	case types.EventPreference:
		return "Personal Preference"
	case types.EventLikeDislike:
		return "Like or Dislike"
	case types.EventEmotion:
		return "Emotional Moment"
	case types.EventAssociation:
		return "Personal Association"
	case types.EventAesthetic:
		return "Aesthetic Reaction"
	case types.EventSensory:
		return "Sensory Experience"
	case types.EventSocial:
		return "Social Moment"
	case types.EventFirst:
		return "First Experience"
	case types.EventShared:
		return "Shared Memory"
	case types.EventPromise:
		return "Promise or Commitment"
	default:
		return "Notable Event"
	}
}

// generateSummary creates a one-line summary from the detection
func generateSummary(result types.EventDetectionResult) string {
	if len(result.Evidence) > 0 {
		return fmt.Sprintf("Detected %s event (keywords: %s)",
			result.EventType, strings.Join(result.Evidence, ", "))
	}
	return fmt.Sprintf("Detected %s event with %.0f%% confidence", result.EventType, result.Confidence)
}

// buildContent constructs detailed memory content from context
func buildContent(result types.EventDetectionResult, sourceMessage string) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Event Type: %s", result.EventType))
	parts = append(parts, fmt.Sprintf("Confidence: %.0f%%", result.Confidence))
	if len(result.Evidence) > 0 {
		parts = append(parts, fmt.Sprintf("Trigger Words: %s", strings.Join(result.Evidence, ", ")))
	}
	parts = append(parts, "")
	parts = append(parts, "Source Message:")
	parts = append(parts, sourceMessage)
	if result.Context != "" && result.Context != sourceMessage {
		parts = append(parts, "")
		parts = append(parts, "Context:")
		parts = append(parts, result.Context)
	}
	return strings.Join(parts, "\n")
}

// generateTags creates relevant tags from event type and character names
func generateTags(result types.EventDetectionResult, identities []CharacterIdentity) []string {
	tags := []string{string(result.EventType), "auto-captured"}

	// Names are descriptive tags only; ownership uses the stable ID.
	lowerContext := strings.ToLower(result.Context)
	for _, identity := range identities {
		name := strings.TrimSpace(identity.Name)
		if name != "" && strings.Contains(lowerContext, strings.ToLower(name)) {
			tag := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
			tags = append(tags, tag)
		}
	}

	return uniqueStrings(tags)
}

// DeduplicateResults removes overlapping detections on the same message
// Keeps the highest-confidence result per event type
func DeduplicateResults(results []types.EventDetectionResult) []types.EventDetectionResult {
	best := make(map[types.EventType]types.EventDetectionResult)

	for _, r := range results {
		existing, ok := best[r.EventType]
		if !ok || r.Confidence > existing.Confidence {
			best[r.EventType] = r
		}
	}

	var out []types.EventDetectionResult
	for _, r := range best {
		out = append(out, r)
	}
	return out
}

// MergeRelatedEvents combines death + conflict or romance + character_change
// when they fire on the same context, keeping the more specific one
func MergeRelatedEvents(results []types.EventDetectionResult) []types.EventDetectionResult {
	hasDeath := false
	for _, r := range results {
		if r.EventType == types.EventDeath {
			hasDeath = true
			break
		}
	}

	if !hasDeath {
		return results
	}

	// If death is present, suppress generic conflict on same context
	var filtered []types.EventDetectionResult
	for _, r := range results {
		if r.EventType == types.EventConflict {
			continue // Death is more specific
		}
		filtered = append(filtered, r)
	}
	return filtered
}

// BatchClassify runs full pipeline: detect -> dedup -> merge -> build candidates
func BatchClassify(results []types.EventDetectionResult, sourceMessage string, identities []CharacterIdentity, chatID string) []types.MemoryCandidate {
	results = DeduplicateResults(results)
	results = MergeRelatedEvents(results)

	var candidates []types.MemoryCandidate
	for _, r := range results {
		candidates = append(candidates, ClassifyAndBuildCandidate(r, sourceMessage, identities, chatID))
	}
	return candidates
}
