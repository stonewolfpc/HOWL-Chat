/**
 * Cross-Scenario Character Continuity
 *
 * Enables character memories to persist across different scenarios.
 * Implements inheritance rules for memory migration between chats.
 *
 * @package memory
 * @version 1.0.0
 */

package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"howl-chat/internal/backend/types"
)

// InheritanceRule defines how memories transfer between scenarios
type InheritanceRule string

const (
	InheritAll        InheritanceRule = "all"         // Transfer all character memories
	InheritPinned     InheritanceRule = "pinned"      // Only pinned memories
	InheritHighImpact InheritanceRule = "high_impact" // Priority 1-3 only
	InheritTagged     InheritanceRule = "tagged"      // Only memories with continuity tag
)

// ContinuityConfig controls cross-scenario inheritance
type ContinuityConfig struct {
	Enabled         bool            `json:"enabled"`
	Rule            InheritanceRule `json:"rule"`
	MaxMemories     int             `json:"max_memories"`
	MinConfidence   float64         `json:"min_confidence"`
	DecayOnTransfer bool            `json:"decay_on_transfer"` // Slightly reduce priority
}

// DefaultContinuityConfig returns sensible defaults
func DefaultContinuityConfig() ContinuityConfig {
	return ContinuityConfig{
		Enabled:         true,
		Rule:            InheritHighImpact,
		MaxMemories:     15,
		MinConfidence:   50,
		DecayOnTransfer: true,
	}
}

// CharacterContinuity manages cross-scenario memory persistence
type CharacterContinuity struct {
	projectRoot string
}

// NewCharacterContinuity creates a continuity manager
func NewCharacterContinuity() *CharacterContinuity {
	root, _ := os.Getwd()
	return &CharacterContinuity{projectRoot: root}
}

// ExportCharacterMemories extracts transferable memories for a character
func (cc *CharacterContinuity) ExportCharacterMemories(characterName string, allMemories []types.MemoryRecord, config ContinuityConfig) []types.MemoryRecord {
	if !config.Enabled {
		return nil
	}

	// Filter to character-scoped memories
	var characterMemories []types.MemoryRecord
	for _, m := range allMemories {
		if m.Scope != "character" || m.OwnerID != characterName {
			continue
		}
		if m.Confidence < config.MinConfidence {
			continue
		}
		characterMemories = append(characterMemories, m)
	}

	// Apply inheritance rule
	var transferable []types.MemoryRecord
	for _, m := range characterMemories {
		switch config.Rule {
		case InheritAll:
			transferable = append(transferable, m)
		case InheritPinned:
			if m.Pinned {
				transferable = append(transferable, m)
			}
		case InheritHighImpact:
			// Use priority inference from tags: high priority = low decay rate
			if isHighImpact(m) {
				transferable = append(transferable, m)
			}
		case InheritTagged:
			if hasTag(m.Tags, "continuity") || hasTag(m.Tags, "cross-scenario") {
				transferable = append(transferable, m)
			}
		}
	}

	// Limit count
	if len(transferable) > config.MaxMemories {
		transferable = transferable[:config.MaxMemories]
	}

	return transferable
}

// ImportCharacterMemories merges exported memories into a new scenario
func (cc *CharacterContinuity) ImportCharacterMemories(characterName string, imported []types.MemoryRecord, existing []types.MemoryRecord, config ContinuityConfig) []types.MemoryRecord {
	if !config.Enabled || len(imported) == 0 {
		return existing
	}

	seen := make(map[string]bool)
	for _, e := range existing {
		seen[e.ID] = true
	}

	var result []types.MemoryRecord
	result = append(result, existing...)

	for _, imp := range imported {
		if seen[imp.ID] {
			continue
		}

		// Apply decay on transfer if configured
		if config.DecayOnTransfer {
			imp = applyTransferDecay(imp)
		}

		// Mark as inherited
		imp.Tags = append(imp.Tags, "inherited")
		imp.UpdatedAt = time.Now().Format(time.RFC3339)

		result = append(result, imp)
		seen[imp.ID] = true
	}

	return result
}

// SaveContinuityConfig persists continuity settings
func (cc *CharacterContinuity) SaveContinuityConfig(chatID string, config ContinuityConfig) error {
	chatDir := filepath.Join(cc.projectRoot, "Chats", chatID)
	if err := os.MkdirAll(chatDir, 0755); err != nil {
		return fmt.Errorf("failed to create chat directory: %w", err)
	}

	path := filepath.Join(chatDir, "continuity_config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal continuity config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// LoadContinuityConfig retrieves continuity settings
func (cc *CharacterContinuity) LoadContinuityConfig(chatID string) (ContinuityConfig, error) {
	path := filepath.Join(cc.projectRoot, "Chats", chatID, "continuity_config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultContinuityConfig(), nil
		}
		return DefaultContinuityConfig(), err
	}

	var config ContinuityConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultContinuityConfig(), err
	}

	return config, nil
}

// --- Helpers ---

func isHighImpact(m types.MemoryRecord) bool {
	// Death, betrayal, romance are always high impact
	for _, tag := range m.Tags {
		low := tag
		switch low {
		case "death", "betrayal", "romance", "promise", "first", "shared":
			return true
		}
	}
	// Pinned memories are high impact
	if m.Pinned {
		return true
	}
	// High confidence indicates importance
	if m.Confidence >= 80 {
		return true
	}
	return false
}

func applyTransferDecay(m types.MemoryRecord) types.MemoryRecord {
	// Slightly reduce confidence to reflect age/context change
	if m.Confidence > 10 {
		m.Confidence -= 5
	}
	return m
}
