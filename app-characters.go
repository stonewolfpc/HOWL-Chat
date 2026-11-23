package main

import (
	"fmt"
	"os"
	"path/filepath"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/characterstore"
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
	fmt.Printf("[TRACE:AddCharacterToChat] id=%s name=%s (total=%d)\n", parsed.ID, parsed.Name, len(a.sessionCharacters()))
	a.session().UpsertCharacter(&parsed)
	return nil
}

// RemoveCharacterFromChat removes a character from the current chat session.
func (a *App) RemoveCharacterFromChat(characterId string) error {
	a.session().RemoveCharacter(characterId)
	return nil
}

// ResetSessionCharacters clears the in-memory character list for the current session.
func (a *App) ResetSessionCharacters() error {
	a.session().SetCharacters(nil)
	return nil
}

// characterStore returns the one canonical character store for this process.
func (a *App) characterStore() *characterstore.Store {
	a.characterMu.Lock()
	defer a.characterMu.Unlock()
	if a.characters == nil {
		a.characters = characterstore.New(getCharactersDir())
	}
	return a.characters
}

// CreateCharacter is the legacy binding name for SaveCharacter. The store owns the
// identity, folder and envelope; this path no longer writes raw JSON itself.
func (a *App) CreateCharacter(character map[string]interface{}) (map[string]interface{}, error) {
	return a.SaveCharacter(character)
}

// SaveCharacter creates or updates a character through the canonical store. A
// supplied revision is enforced, so a stale editor cannot overwrite a newer save.
func (a *App) SaveCharacter(character map[string]interface{}) (map[string]interface{}, error) {
	if character == nil {
		return nil, fmt.Errorf("SaveCharacter: character is required")
	}
	base, _ := character["revision"].(string)
	item, err := a.characterStore().Save(character, base)
	if err != nil {
		a.logError(errors.CHARSaveFail, errors.High, "character.save", "app-characters.go", err)
		return nil, fmt.Errorf("SaveCharacter: %w", err)
	}
	// Return the canonical document plus the revision the next save must cite.
	document := item.Document
	document["revision"] = item.Revision
	return document, nil
}

// DeleteCharacter removes one character through the canonical store.
func (a *App) DeleteCharacter(characterID string) error {
	if err := a.characterStore().Delete(characterID); err != nil {
		a.logError(errors.CHARSaveFail, errors.High, "character.delete", "app-characters.go", err)
		return fmt.Errorf("DeleteCharacter: %w", err)
	}
	return nil
}

// ImportCharacters upserts legacy/external character documents. It is
// idempotent and returns the previous document per item for rollback.
func (a *App) ImportCharacters(characters []interface{}) ([]map[string]interface{}, error) {
	documents := make([]map[string]interface{}, 0, len(characters))
	for _, raw := range characters {
		if document, ok := raw.(map[string]interface{}); ok && document != nil {
			documents = append(documents, document)
		}
	}
	results, err := a.characterStore().Import(documents)
	if err != nil {
		a.logError(errors.CHARSaveFail, errors.High, "character.import", "app-characters.go", err)
		return nil, fmt.Errorf("ImportCharacters: %w", err)
	}
	out := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		entry := map[string]interface{}{
			"id":       result.ID,
			"status":   result.Status,
			"revision": result.Revision,
		}
		if result.Previous != nil {
			entry["previous"] = result.Previous
		}
		out = append(out, entry)
	}
	return out, nil
}

// RollbackCharacters restores what an ImportCharacters call replaced.
func (a *App) RollbackCharacters(results []interface{}) error {
	entries := make([]characterstore.ImportResult, 0, len(results))
	for _, raw := range results {
		result, ok := raw.(map[string]interface{})
		if !ok || result == nil {
			continue
		}
		entry := characterstore.ImportResult{}
		entry.ID, _ = result["id"].(string)
		entry.Status, _ = result["status"].(string)
		entry.Revision, _ = result["revision"].(string)
		entry.Previous, _ = result["previous"].(map[string]interface{})
		entries = append(entries, entry)
	}
	if err := a.characterStore().Rollback(entries); err != nil {
		a.logError(errors.CHARSaveFail, errors.High, "character.rollback", "app-characters.go", err)
		return fmt.Errorf("RollbackCharacters: %w", err)
	}
	return nil
}

// getCharactersDir resolves the Characters directory path (cwd or exe dir).
func getCharactersDir() string {
	if dir := isolatedStateDir("Characters"); dir != "" {
		return dir
	}
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, "Characters")
	if _, err := os.Stat(dir); err != nil {
		dir = filepath.Join(filepath.Dir(os.Args[0]), "Characters")
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

// characterExistsOnDisk reports whether the canonical store holds an identity.
func characterExistsOnDisk(name string) bool {
	if name == "" {
		return false
	}
	_, err := characterstore.New(getCharactersDir()).Get(name)
	return err == nil
}

// loadCharacterFromDiskByID reads one character through the canonical store.
func loadCharacterFromDiskByID(id string) *chat.CharacterContext {
	if id == "" {
		return nil
	}
	item, err := characterstore.New(getCharactersDir()).Get(id)
	if err != nil {
		return nil
	}
	cc, err := chat.ParseCharacterContextFromMap(item.Document)
	if err != nil {
		return nil
	}
	return &cc
}

// GetCharacters returns the canonical character documents with the revision
// each one must cite on its next save.
func (a *App) GetCharacters() ([]map[string]interface{}, error) {
	items, err := a.characterStore().List()
	if err != nil {
		a.logError(errors.CHARLoadFail, errors.High, "character.load.all", "app-characters.go", err)
		return nil, fmt.Errorf("GetCharacters: %w", err)
	}
	characters := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		document := item.Document
		if _, ok := document["id"]; !ok {
			document["id"] = item.ID
		}
		if _, ok := document["name"]; !ok {
			document["name"] = item.ID
		}
		document["revision"] = item.Revision
		characters = append(characters, document)
	}
	return characters, nil
}

// AddWorldToChat sets the world context for the current chat session.
func (a *App) AddWorldToChat(worldMap map[string]interface{}) error {
	world, err := a.parseWorldContextFromMap(worldMap)
	if err != nil {
		return fmt.Errorf("AddWorldToChat: failed to parse world: %w", err)
	}
	a.session().SetWorld(&world)
	if a.conversation.ActiveID() != "" {
		a.performSaveLogged("app-characters.go")
	}
	return nil
}

// AddScenarioToChat sets the scenario context for the current chat session.
func (a *App) AddScenarioToChat(scenario map[string]interface{}) error {
	parsed, err := a.parseScenarioContextFromMap(scenario)
	if err != nil {
		return fmt.Errorf("AddScenarioToChat: failed to parse scenario: %w", err)
	}
	a.session().SetScenario(&parsed)
	if a.conversation.ActiveID() != "" {
		a.performSaveLogged("app-characters.go")
	}
	return nil
}

// ---------------------------------------------------------------------------
// State sync helpers
// ---------------------------------------------------------------------------

// loadChatDataFromStorage loads chat state from the active persisted chat.
// Used as a fallback by LoadChatData when no in-memory state is present.
func (a *App) loadChatDataFromStorage() (map[string]interface{}, error) {
	if !a.conversation.Bound() {
		return nil, fmt.Errorf("chat persistence not initialized")
	}
	chatID := a.conversation.ActiveID()
	if chatID == "" {
		return nil, fmt.Errorf("no active chat")
	}
	chatData, err := a.conversation.Store().LoadChat(chatID)
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

// saveChatDataToStorage persists the supplied chat state map through the
// conversation owner, which refuses a record the collection does not claim.
func (a *App) saveChatDataToStorage(data map[string]interface{}) error {
	chatData, err := a.convertFrontendDataToChatData(data)
	if err != nil {
		return fmt.Errorf("saveChatDataToStorage: %w", err)
	}
	return a.conversation.SaveRecord(chatData)
}

// syncAppStateFromChatService pulls the latest world/scenario/character/lorebook
// state out of the chat.Service into the App fields.
func (a *App) syncAppStateFromChatService() {
	a.session().SetWorld(a.chatService.GetWorldData())
	a.session().SetScenario(a.chatService.GetScenarioData())
	a.session().SetCharacters(a.chatService.GetCharacterData())
	// Same overlay as SetLorebooks: whichever sync runs last, the prompt pool keeps
	// canonical text after cutover instead of the chat service's legacy copies.
	a.session().SetLorebooks(a.mergeCanonicalLore(a.chatService.GetAllLorebookEntries(), false))
}

// applyChatDataToState keeps backend state in sync with frontend chat state.
func (a *App) applyChatDataToState(data map[string]interface{}) {
	if worldMap, ok := data["world"].(map[string]interface{}); ok && worldMap != nil {
		world, err := a.parseWorldContextFromMap(worldMap)
		if err != nil {
			a.logErrorSilent(errors.WORLDParseFail, errors.Medium, "world.parse", "app-characters.go:237", err)
			a.session().ClearWorld()
		} else {
			a.session().SetWorld(&world)
		}
	} else {
		a.session().ClearWorld()
	}

	if scenarioMap, ok := data["scenario"].(map[string]interface{}); ok && scenarioMap != nil {
		scenario, err := a.parseScenarioContextFromMap(scenarioMap)
		if err != nil {
			a.logErrorSilent(errors.SCENParseFail, errors.Medium, "scenario.parse", "app-characters.go:249", err)
			a.session().ClearScenario()
		} else {
			a.session().SetScenario(&scenario)
		}
	} else {
		a.session().ClearScenario()
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
		a.session().SetCharacters(characters)
	} else {
		a.session().SetCharacters(nil)
	}

	if lorebooks, ok := data["lorebooks"].([]interface{}); ok {
		parsed := make([]lorebook.Entry, 0, len(lorebooks))
		for _, l := range lorebooks {
			if lorebookMap, ok := l.(map[string]interface{}); ok {
				entry, err := lorebook.ParseLorebookEntryFromMap(lorebookMap)
				if err != nil {
					a.logErrorSilent(errors.LOREParseFail, errors.Medium, "lorebook.parse", "app-characters.go:281", err)
					continue
				}
				parsed = append(parsed, entry)
			}
		}
		a.session().SetLorebooks(parsed)
	} else {
		a.session().SetLorebooks(nil)
	}

	if chatName, ok := data["chatName"].(string); ok {
		_ = chatName
	}
}
