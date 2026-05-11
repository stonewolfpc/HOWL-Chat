package chat

import (
	"encoding/json"
	"errors"
	"os"

	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/backend/types"
	entities "howl-chat/internal/backend/types/entities"
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
func loadWorldData(filePath string) (*entities.World, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // Return nil if file doesn't exist
		}
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to read world file", err)
	}

	var world entities.World
	if err := json.Unmarshal(data, &world); err != nil {
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to unmarshal world JSON", err)
	}
	return &world, nil
}

// loadScenarioData loads scenario data from a given JSON file path.
func loadScenarioData(filePath string) (*entities.Scenario, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // Return nil if file doesn't exist
		}
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to read scenario file", err)
	}

	var scenario entities.Scenario
	if err := json.Unmarshal(data, &scenario); err != nil {
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to unmarshal scenario JSON", err)
	}
	return &scenario, nil
}

// loadCharacterData loads character data from a given JSON file path.
func loadCharacterData(filePath string) (*entities.Character, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // Return nil if file doesn't exist
		}
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to read character file", err)
	}

	var character entities.Character
	if err := json.Unmarshal(data, &character); err != nil {
		return nil, types.WrapError(types.ErrorCodeInternal, "failed to unmarshal character JSON", err)
	}
	return &character, nil
}
