package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/errors"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ---------------------------------------------------------------------------
// Character / World / Scenario Wails Bindings
// ---------------------------------------------------------------------------

// OpenCharacterSheet navigates the app window to the character editor.
func (a *App) OpenCharacterSheet(characterId string) error {
	runtime.WindowExecJS(a.ctx, fmt.Sprintf("window.location.href='character.html?id=%s'", characterId))
	return nil
}

// AddCharacterToChat adds or replaces a character in the current chat session.
func (a *App) AddCharacterToChat(character map[string]interface{}) error {
	parsed, err := a.parseCharacterContextFromMap(character)
	if err != nil {
		return fmt.Errorf("AddCharacterToChat: failed to parse character: %w", err)
	}
	fmt.Printf("[TRACE:AddCharacterToChat] id=%s name=%s (total=%d)\n", parsed.ID, parsed.Name, len(a.characters))
	for i, c := range a.characters {
		if c.ID == parsed.ID {
			a.characters[i] = &parsed
			return nil
		}
	}
	a.characters = append(a.characters, &parsed)
	return nil
}

// RemoveCharacterFromChat removes a character from the current chat session.
func (a *App) RemoveCharacterFromChat(characterId string) error {
	filtered := make([]*chat.CharacterContext, 0, len(a.characters))
	for _, c := range a.characters {
		if c != nil && c.ID != characterId {
			filtered = append(filtered, c)
		}
	}
	a.characters = filtered
	return nil
}

// ResetSessionCharacters clears the in-memory character list for the current session.
func (a *App) ResetSessionCharacters() error {
	a.characters = []*chat.CharacterContext{}
	return nil
}

// CreateCharacter creates a new character and writes it to Characters/<name>/character.json.
func (a *App) CreateCharacter(character map[string]interface{}) (map[string]interface{}, error) {
	if character["id"] == nil || character["id"] == "" {
		character["id"] = fmt.Sprintf("char_%d", time.Now().Unix())
	}

	name, _ := character["name"].(string)
	if name == "" {
		name, _ = character["id"].(string)
	}

	cwd, _ := os.Getwd()
	charDir := filepath.Join(cwd, "Characters", name)
	if _, err := os.Stat(filepath.Join(filepath.Dir(os.Args[0]), "Characters")); err == nil {
		charDir = filepath.Join(filepath.Dir(os.Args[0]), "Characters", name)
	}
	if err := os.MkdirAll(charDir, 0755); err != nil {
		a.logError(errors.CHARDirFail, errors.High, "character.create", "app-characters.go:59", err)
		return nil, fmt.Errorf("CreateCharacter: failed to create directory: %w", err)
	}

	data, err := json.Marshal(character)
	if err != nil {
		a.logError(errors.CHARCreateFail, errors.High, "character.create", "app-characters.go:68", err)
		return nil, fmt.Errorf("CreateCharacter: failed to marshal character: %w", err)
	}
	charPath := filepath.Join(charDir, "character.json")
	if err := os.WriteFile(charPath, data, 0644); err != nil {
		a.logError(errors.CHARSaveFail, errors.High, "character.create", "app-characters.go:73", err)
		return nil, fmt.Errorf("CreateCharacter: failed to write character: %w", err)
	}
	return character, nil
}

// getCharactersDir resolves the Characters directory path (cwd or exe dir).
func getCharactersDir() string {
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, "Characters")
	if _, err := os.Stat(dir); err != nil {
		dir = filepath.Join(filepath.Dir(os.Args[0]), "Characters")
	}
	return dir
}

// getWorldsDir resolves the Worlds directory path.
func getWorldsDir() string {
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, "Worlds")
	if _, err := os.Stat(dir); err != nil {
		dir = filepath.Join(filepath.Dir(os.Args[0]), "Worlds")
	}
	return dir
}

// getScenariosDir resolves the Scenarios directory path.
func getScenariosDir() string {
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, "Scenarios")
	if _, err := os.Stat(dir); err != nil {
		dir = filepath.Join(filepath.Dir(os.Args[0]), "Scenarios")
	}
	return dir
}

// listDirNames returns the names of subdirectories in a given path.
func listDirNames(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names
}

// characterExistsOnDisk returns true if a character's JSON file still exists.
func characterExistsOnDisk(name string) bool {
	if name == "" {
		return false
	}
	dir := getCharactersDir()
	path := filepath.Join(dir, name, "character.json")
	_, err := os.Stat(path)
	exists := err == nil
	if !exists {
		fmt.Printf("[DEBUG characterExistsOnDisk] name=%q path=%q exists=%v err=%v\n", name, path, exists, err)
		fmt.Printf("[DEBUG] Characters dir contents: %v\n", listDirNames(dir))
	}
	return exists
}

// loadCharacterFromDiskByID scans the Characters directory for a character
// whose JSON file contains the matching ID.  Returns nil if not found.
func loadCharacterFromDiskByID(id string) *chat.CharacterContext {
	if id == "" {
		return nil
	}
	dir := getCharactersDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name(), "character.json")
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var raw map[string]interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			continue
		}
		if raw["id"] != id {
			continue
		}
		cc, err := chat.ParseCharacterContextFromMap(raw)
		if err != nil {
			continue
		}
		return &cc
	}
	return nil
}

// worldExistsOnDisk returns true if a world's JSON file still exists.
func worldExistsOnDisk(name string) bool {
	if name == "" {
		return false
	}
	dir := getWorldsDir()
	path := filepath.Join(dir, name, "world.json")
	_, err := os.Stat(path)
	exists := err == nil
	if !exists {
		fmt.Printf("[DEBUG worldExistsOnDisk] name=%q path=%q exists=%v err=%v\n", name, path, exists, err)
		fmt.Printf("[DEBUG] Worlds dir contents: %v\n", listDirNames(dir))
	}
	return exists
}

// scenarioExistsOnDisk returns true if a scenario's JSON file still exists.
func scenarioExistsOnDisk(name string) bool {
	if name == "" {
		return false
	}
	dir := getScenariosDir()
	path := filepath.Join(dir, name, "scenario.json")
	_, err := os.Stat(path)
	exists := err == nil
	if !exists {
		fmt.Printf("[DEBUG scenarioExistsOnDisk] name=%q path=%q exists=%v err=%v\n", name, path, exists, err)
		fmt.Printf("[DEBUG] Scenarios dir contents: %v\n", listDirNames(dir))
	}
	return exists
}

// GetCharacters loads all characters from the Characters directory.
func (a *App) GetCharacters() ([]map[string]interface{}, error) {
	var characters []map[string]interface{}

	cwd, _ := os.Getwd()
	charactersDir := filepath.Join(cwd, "Characters")
	if _, err := os.Stat(charactersDir); err != nil {
		charactersDir = filepath.Join(filepath.Dir(os.Args[0]), "Characters")
	}
	_ = os.MkdirAll(charactersDir, 0755)

	entries, err := os.ReadDir(charactersDir)
	if err != nil {
		a.logError(errors.CHARLoadFail, errors.High, "character.load.all", "app-characters.go:90", err)
		return characters, fmt.Errorf("failed to read characters directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		charName := entry.Name()
		charPath := filepath.Join(charactersDir, charName, "character.json")
		data, err := os.ReadFile(charPath)
		if err != nil {
			continue
		}
		if len(data) == 0 || string(data) == "{}" {
			characters = append(characters, map[string]interface{}{
				"id": charName, "name": charName, "howl": map[string]interface{}{},
			})
			continue
		}
		var charData map[string]interface{}
		if err := json.Unmarshal(data, &charData); err != nil {
			a.logErrorSilent(errors.CHARParseFail, errors.Medium, "character.load.parse", "app-characters.go:112", err)
			charData = map[string]interface{}{
				"id": charName, "name": charName, "howl": map[string]interface{}{},
			}
		}
		if _, ok := charData["name"]; !ok {
			charData["name"] = charName
		}
		if _, ok := charData["id"]; !ok {
			charData["id"] = charName
		}
		characters = append(characters, charData)
	}
	return characters, nil
}

// AddWorldToChat sets the world context for the current chat session.
func (a *App) AddWorldToChat(worldMap map[string]interface{}) error {
	world, err := a.parseWorldContextFromMap(worldMap)
	if err != nil {
		return fmt.Errorf("AddWorldToChat: failed to parse world: %w", err)
	}
	a.world = &world
	if a.chatPersistence != nil {
		if chatID := a.chatPersistence.GetLastChatID(); chatID != "" {
			_ = a.performSave()
		}
	}
	return nil
}

// AddScenarioToChat sets the scenario context for the current chat session.
func (a *App) AddScenarioToChat(scenario map[string]interface{}) error {
	parsed, err := a.parseScenarioContextFromMap(scenario)
	if err != nil {
		return fmt.Errorf("AddScenarioToChat: failed to parse scenario: %w", err)
	}
	a.scenario = &parsed
	if a.chatPersistence != nil {
		if chatID := a.chatPersistence.GetLastChatID(); chatID != "" {
			_ = a.performSave()
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// State sync helpers
// ---------------------------------------------------------------------------

// loadChatDataFromStorage loads chat state from the active persisted chat.
// Used as a fallback by LoadChatData when no in-memory state is present.
func (a *App) loadChatDataFromStorage() (map[string]interface{}, error) {
	if a.chatPersistence == nil {
		return nil, fmt.Errorf("chat persistence not initialized")
	}
	chatID := a.chatPersistence.GetLastChatID()
	if chatID == "" {
		return nil, fmt.Errorf("no active chat")
	}
	chatData, err := a.chatPersistence.LoadChat(chatID)
	if err != nil {
		a.logError(errors.CHATLoadFail, errors.High, "chat.load", "app-characters.go:172", err)
		return nil, err
	}

	charactersSlice := make([]interface{}, 0, len(chatData.Characters))
	for _, c := range chatData.Characters {
		if c == nil {
			continue
		}
		charactersSlice = append(charactersSlice, map[string]interface{}{
			"id": c.ID, "name": c.Name,
			"personality": c.Personality, "backstory": c.Backstory,
		})
	}

	var worldMap interface{}
	if chatData.World != nil {
		worldMap = map[string]interface{}{
			"id": chatData.World.ID, "name": chatData.World.Name,
			"description": chatData.World.Description,
		}
	}

	var scenarioMap interface{}
	if chatData.Scenario != nil {
		scenarioMap = map[string]interface{}{
			"id": chatData.Scenario.ID, "name": chatData.Scenario.Name,
			"description": chatData.Scenario.Description,
		}
	}

	return map[string]interface{}{
		"characters": charactersSlice,
		"world":      worldMap,
		"scenario":   scenarioMap,
		"lorebooks":  chatData.Lorebooks,
		"chatName":   chatData.Meta.Name,
		"meta": map[string]interface{}{
			"id":                chatData.Meta.ID,
			"name":              chatData.Meta.Name,
			"worldId":           chatData.Meta.WorldID,
			"scenarioId":        chatData.Meta.ScenarioID,
			"activeCharacters":  chatData.Meta.ActiveCharacters,
			"removedCharacters": chatData.Meta.RemovedCharacters,
			"addedCharacters":   chatData.Meta.AddedCharacters,
			"activeLorebooks":   chatData.Meta.ActiveLorebooks,
			"removedLorebooks":  chatData.Meta.RemovedLorebooks,
			"addedLorebooks":    chatData.Meta.AddedLorebooks,
		},
	}, nil
}

// saveChatDataToStorage persists the supplied chat state map via chatPersistence.
func (a *App) saveChatDataToStorage(data map[string]interface{}) error {
	if a.chatPersistence == nil {
		return nil
	}
	chatData, err := a.convertFrontendDataToChatData(data)
	if err != nil {
		return fmt.Errorf("saveChatDataToStorage: %w", err)
	}
	return a.chatPersistence.SaveChat(chatData)
}

// syncAppStateFromChatService pulls the latest world/scenario/character/lorebook
// state out of the chat.Service into the App fields.
func (a *App) syncAppStateFromChatService() {
	a.world = a.chatService.GetWorldData()
	a.scenario = a.chatService.GetScenarioData()
	a.characters = a.chatService.GetCharacterData()
	a.lorebooks = a.chatService.GetAllLorebookEntries()
}

// applyChatDataToState keeps backend state in sync with frontend chat state.
func (a *App) applyChatDataToState(data map[string]interface{}) {
	if worldMap, ok := data["world"].(map[string]interface{}); ok && worldMap != nil {
		world, err := a.parseWorldContextFromMap(worldMap)
		if err != nil {
			a.logErrorSilent(errors.WORLDParseFail, errors.Medium, "world.parse", "app-characters.go:237", err)
			a.world = nil
		} else {
			a.world = &world
		}
	} else {
		a.world = nil
	}

	if scenarioMap, ok := data["scenario"].(map[string]interface{}); ok && scenarioMap != nil {
		scenario, err := a.parseScenarioContextFromMap(scenarioMap)
		if err != nil {
			a.logErrorSilent(errors.SCENParseFail, errors.Medium, "scenario.parse", "app-characters.go:249", err)
			a.scenario = nil
		} else {
			a.scenario = &scenario
		}
	} else {
		a.scenario = nil
	}

	if charsSlice, ok := data["characters"].([]interface{}); ok {
		var characters []*chat.CharacterContext
		for _, charItem := range charsSlice {
			if charMap, ok := charItem.(map[string]interface{}); ok {
				character, err := a.parseCharacterContextFromMap(charMap)
				if err != nil {
					a.logErrorSilent(errors.CHARParseFail, errors.Medium, "character.parse", "app-characters.go:264", err)
					continue
				}
				characters = append(characters, &character)
			}
		}
		a.characters = characters
	} else {
		a.characters = []*chat.CharacterContext{}
	}

	if lorebooks, ok := data["lorebooks"].([]interface{}); ok {
		a.lorebooks = make([]lorebook.Entry, 0, len(lorebooks))
		for _, l := range lorebooks {
			if lorebookMap, ok := l.(map[string]interface{}); ok {
				entry, err := lorebook.ParseLorebookEntryFromMap(lorebookMap)
				if err != nil {
					a.logErrorSilent(errors.LOREParseFail, errors.Medium, "lorebook.parse", "app-characters.go:281", err)
					continue
				}
				a.lorebooks = append(a.lorebooks, entry)
			}
		}
	} else {
		a.lorebooks = []lorebook.Entry{}
	}

	if chatName, ok := data["chatName"].(string); ok {
		_ = chatName
	}
}
