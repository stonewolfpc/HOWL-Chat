package chat

import (
	"encoding/json"
	"errors"
	"os"

	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/backend/memory"
	"howl-chat/internal/backend/types"
)

// loadEntriesFromFile loads lorebook entries from a given JSON file path.
func loadEntriesFromFile(filePath string) ([]lorebook.Entry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		// If file does not exist, return empty slice and no error.
		if errors.Is(err, os.ErrNotExist) {
			return []lorebook.Entry{}, nil
		}
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to read lorebook file", err)
	}

	var lb struct {
		Entries []lorebook.Entry `json:"entries"`
	}

	if err := json.Unmarshal(data, &lb); err != nil {
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to unmarshal lorebook JSON", err)
	}

	// JSON lore files don't contain an "enabled" field; default all entries to enabled.
	for i := range lb.Entries {
		lb.Entries[i].Enabled = true
	}

	return lb.Entries, nil
}

// loadMemoriesSnapshot loads memories from memories_snapshot.json.
func loadMemoriesSnapshot(filePath string) ([]lorebook.Entry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []lorebook.Entry{}, nil
		}
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to read memories snapshot file", err)
	}

	var ms struct {
		WorldMemories     []lorebook.Entry            `json:"world_memories"`
		ScenarioMemories  []lorebook.Entry            `json:"scenario_memories"`
		CharacterMemories map[string][]lorebook.Entry `json:"character_memories"`
	}

	if err := json.Unmarshal(data, &ms); err != nil {
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to unmarshal memories snapshot JSON", err)
	}

	// Combine all memories into a single slice for easier processing
	allMemories := make([]lorebook.Entry, 0)
	allMemories = append(allMemories, ms.WorldMemories...)
	allMemories = append(allMemories, ms.ScenarioMemories...)
	for _, charMemories := range ms.CharacterMemories {
		allMemories = append(allMemories, charMemories...)
	}

	return allMemories, nil
}

// loadCharacterData loads character data from a given JSON file path.
func loadCharacterData(filePath string) (*CharacterContext, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // Return nil if file doesn't exist
		}
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to read character file", err)
	}

	character, err := decodeCharacterContext(data)
	if err != nil {
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to unmarshal character JSON", err)
	}
	return character, nil
}

// loadMemoriesFromService loads user-created memories via the memory service
// and converts them to lorebook.Entry format for resolution. World and scenario
// records are matched by stable ID only; a display name, or an empty owner, is
// unresolved and excluded from the prompt pool.
func loadMemoriesFromService(memService *memory.Service, worldID, scenarioID string, characterIDs []string) ([]lorebook.Entry, error) {
	// Get all memories from service
	filter := types.MemoryFilter{Status: "active"}
	records, err := memService.GetMemories(filter)
	if err != nil {
		return nil, err
	}

	// Filter records before conversion so a legacy name-owned memory can never be
	// promoted to lorebook entry ownership merely by matching a character ID.
	var filteredRecords []types.MemoryRecord
	for _, record := range records {
		switch record.Scope {
		case "world":
			if worldID != "" && record.OwnerID == worldID {
				filteredRecords = append(filteredRecords, record)
			}
		case "scenario":
			if scenarioID != "" && record.OwnerID == scenarioID {
				filteredRecords = append(filteredRecords, record)
			}
		case "character":
			if memory.IsStableCharacterOwnershipForAny(record, characterIDs) {
				filteredRecords = append(filteredRecords, record)
			}
		default:
			filteredRecords = append(filteredRecords, record)
		}
	}

	return memory.ToLorebookEntries(filteredRecords), nil
}
