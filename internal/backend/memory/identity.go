package memory

import (
	"fmt"
	"sort"

	"howl-chat/internal/backend/types"
)

// CharacterIdentity is the stable identity supplied to memory capture. Name is
// only a matching label; ID is what ownership is allowed to store.
type CharacterIdentity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// IsStableCharacterOwnership reports whether a memory may be treated as
// character-private lore for characterID. Records written before owner_kind
// existed have an empty kind and therefore fail closed.
func IsStableCharacterOwnership(record types.MemoryRecord, characterID string) bool {
	return characterID != "" &&
		record.Scope == "character" &&
		record.OwnerKind == types.OwnerKindCharacterID &&
		record.OwnerID == characterID
}

// IsStableCharacterOwnershipForAny reports whether a memory is owned by one of
// the supplied stable character IDs. It is the prompt-boundary filter.
func IsStableCharacterOwnershipForAny(record types.MemoryRecord, characterIDs []string) bool {
	for _, characterID := range characterIDs {
		if IsStableCharacterOwnership(record, characterID) {
			return true
		}
	}
	return false
}

// IdentityDiagnostics reports character-scoped memory ownership without
// exposing memory content. Unresolved records are retained and never reassigned.
type IdentityDiagnostics struct {
	Total            int      `json:"total"`
	Stable           int      `json:"stable"`
	Unresolved       int      `json:"unresolved"`
	UnresolvedOwners []string `json:"unresolvedOwners,omitempty"`
}

// summarizeMemoryIdentities produces the ownership report without touching
// storage. Unresolved owners are listed for display/audit, not reassigned.
func summarizeMemoryIdentities(memories []types.MemoryRecord) IdentityDiagnostics {
	diag := IdentityDiagnostics{}
	unresolved := map[string]bool{}
	for _, record := range memories {
		if record.Scope != "character" {
			continue
		}
		diag.Total++
		if record.OwnerKind == types.OwnerKindCharacterID && record.OwnerID != "" {
			diag.Stable++
			continue
		}
		diag.Unresolved++
		if record.OwnerID != "" {
			unresolved[record.OwnerID] = true
		}
	}
	for owner := range unresolved {
		diag.UnresolvedOwners = append(diag.UnresolvedOwners, owner)
	}
	sort.Strings(diag.UnresolvedOwners)
	return diag
}

// GetMemoryIdentityDiagnostics summarizes memory ownership for the current chat.
func (s *Service) GetMemoryIdentityDiagnostics() (IdentityDiagnostics, error) {
	storage := NewMemoryStorage(s.chatID)
	memories, err := storage.loadMemories()
	if err != nil {
		return IdentityDiagnostics{}, fmt.Errorf("failed to load memories: %w", err)
	}
	return summarizeMemoryIdentities(memories), nil
}
