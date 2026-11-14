/**
 * Memory Types
 *
 * This package defines all memory-related types for the chat system.
 * Memories represent automated lorebook entries captured during chat sessions.
 *
 * @package types
 */

package types

import "time"

// Memory ownership kinds. An empty kind is legacy/unresolved and must not be
// treated as a stable character identity.
const (
	OwnerKindCharacterID = "character_id"
	OwnerKindChat        = "chat"
)

// MemoryRecord represents a single automated memory entry
type MemoryRecord struct {
	ID           string   `json:"id"`                   // Unique memory identifier
	Title        string   `json:"title"`                // Memory title
	Summary      string   `json:"summary"`              // Memory summary
	Content      string   `json:"content"`              // Full memory content
	Scope        string   `json:"scope"`                // Memory scope (global/world/scenario/character)
	OwnerID      string   `json:"owner_id"`             // Owner ID for scoped memories
	OwnerKind    string   `json:"owner_kind,omitempty"` // character_id, chat, or empty (legacy)
	Tags         []string `json:"tags"`                 // Memory tags for categorization
	Pinned       bool     `json:"pinned"`               // Whether memory is pinned
	Status       string   `json:"status"`               // Memory status (active/archived)
	Confidence   float64  `json:"confidence"`           // Confidence score (0-100)
	Source       string   `json:"source"`               // Source of the memory (auto/manual)
	TriggerCount int      `json:"trigger_count"`        // Times this memory has triggered
	CreatedAt    string   `json:"created_at"`           // When the memory was created
	UpdatedAt    string   `json:"updated_at"`           // When the memory was last updated
}

// MemoryFilter represents filtering options for memory queries
type MemoryFilter struct {
	Query  string `json:"query"`  // Search query string
	Scope  string `json:"scope"`  // Filter by scope
	Status string `json:"status"` // Filter by status
}

// MemoryEdit represents edits to apply to a memory
type MemoryEdit struct {
	Title   string   `json:"title"`   // New title (empty to keep current)
	Summary string   `json:"summary"` // New summary (empty to keep current)
	Tags    []string `json:"tags"`    // New tags (empty to keep current)
	Pinned  bool     `json:"pinned"`  // New pinned status
}

// AutoCaptureConfig represents configuration for automatic memory capture
type AutoCaptureConfig struct {
	Enabled       bool     `json:"enabled"`         // Whether auto-capture is enabled
	MinConfidence float64  `json:"min_confidence"`  // Minimum confidence threshold
	MaxPerMessage int      `json:"max_per_message"` // Max memories to capture per message
	TriggerWords  []string `json:"trigger_words"`   // Words that trigger memory capture
}

// NewMemoryRecord creates a new memory record
func NewMemoryRecord(title, summary, content, scope, ownerID string) *MemoryRecord {
	now := time.Now().Format(time.RFC3339)
	return &MemoryRecord{
		ID:         NewMemoryID(),
		Title:      title,
		Summary:    summary,
		Content:    content,
		Scope:      scope,
		OwnerID:    ownerID,
		Tags:       []string{},
		Pinned:     false,
		Status:     "active",
		Confidence: 0.0,
		Source:     "auto",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// NewMemoryID generates a unique identifier for a memory
func NewMemoryID() string {
	return "mem_" + time.Now().Format("20060102150405")
}

// EventType represents categories of memory-worthy events
type EventType string

const (
	EventDeath           EventType = "death"
	EventBetrayal        EventType = "betrayal"
	EventRomance         EventType = "romance"
	EventDiscovery       EventType = "discovery"
	EventConflict        EventType = "conflict"
	EventCharacterChange EventType = "character_change"
	EventLocation        EventType = "location"
	EventItem            EventType = "item"
	EventCustom          EventType = "custom"

	// Micro-memory event types (small, personal, emotional moments)
	EventPreference  EventType = "preference"   // favorite food, music, color
	EventLikeDislike EventType = "like_dislike" // "I love this", "I hate crowds"
	EventEmotion     EventType = "emotion"      // "this makes me happy", "I feel safe"
	EventAssociation EventType = "association"  // "reminds me of childhood"
	EventAesthetic   EventType = "aesthetic"    // "this is beautiful", "love the colors"
	EventSensory     EventType = "sensory"      // "that smell is comforting"
	EventSocial      EventType = "social"       // inside jokes, compliments, gratitude
	EventFirst       EventType = "first"        // first time they experienced something
	EventShared      EventType = "shared"       // "you showed me", "we did this together"
	EventPromise     EventType = "promise"      // "I'll remember", "won't forget"
)

// MemoryCandidate represents a detected event awaiting user review
type MemoryCandidate struct {
	ID            string    `json:"id"`
	EventType     EventType `json:"event_type"`
	Title         string    `json:"title"`
	Summary       string    `json:"summary"`
	Content       string    `json:"content"`
	SourceMessage string    `json:"source_message"`
	Scope         string    `json:"scope"`
	OwnerID       string    `json:"owner_id"`
	OwnerKind     string    `json:"owner_kind,omitempty"`
	Tags          []string  `json:"tags"`
	Confidence    float64   `json:"confidence"`
	Status        string    `json:"status"` // pending | approved | rejected
	CreatedAt     string    `json:"created_at"`
	ReviewedAt    string    `json:"reviewed_at,omitempty"`
}

// EventDetectionResult is the raw output of event detection on a message
type EventDetectionResult struct {
	EventType  EventType `json:"event_type"`
	Confidence float64   `json:"confidence"`
	Evidence   []string  `json:"evidence"`
	Context    string    `json:"context"`
}
