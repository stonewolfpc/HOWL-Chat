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

// loadWorldData loads world data from a given JSON file path.
func loadWorldData(filePath string) (*WorldContext, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // Return nil if file doesn't exist
		}
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to read world file", err)
	}

	world, err := decodeWorldContext(data)
	if err != nil {
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to unmarshal world JSON", err)
	}
	return world, nil
}

// loadScenarioData loads scenario data from a given JSON file path.
func loadScenarioData(filePath string) (*ScenarioContext, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // Return nil if file doesn't exist
		}
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to read scenario file", err)
	}

	scenario, err := decodeScenarioContext(data)
	if err != nil {
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to unmarshal scenario JSON", err)
	}
	return scenario, nil
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
// and converts them to lorebook.Entry format for resolution
func loadMemoriesFromService(memService *memory.Service, worldName, scenarioName string, characterNames []string) ([]lorebook.Entry, error) {
	// Get all memories from service
	filter := types.MemoryFilter{Status: "active"}
	records, err := memService.GetMemories(filter)
	if err != nil {
		return nil, err
	}

	// Convert to lorebook entries
	entries := memory.ToLorebookEntries(records)

	// Filter to only include memories relevant to current context
	var filtered []lorebook.Entry
	for _, entry := range entries {
		switch entry.Scope {
		case lorebook.ScopeWorld:
			if entry.OwnerID == "" || entry.OwnerID == worldName {
				filtered = append(filtered, entry)
			}
		case lorebook.ScopeScenario:
			if entry.OwnerID == "" || entry.OwnerID == scenarioName {
				filtered = append(filtered, entry)
			}
		case lorebook.ScopeCharacter:
			if contains(characterNames, entry.OwnerID) {
				filtered = append(filtered, entry)
			}
		default:
			// Global scope - include all
			filtered = append(filtered, entry)
		}
	}

	return filtered, nil
}
