package main

import (
	"fmt"
	"strings"
	"time"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/errors"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ---------------------------------------------------------------------------
// Chat Persistence Wails Bindings
// ---------------------------------------------------------------------------

// GetChatList returns a list of all chat summaries.
func (a *App) GetChatList() ([]chat.ChatSummary, error) {
	if a.chatPersistence == nil {
		return []chat.ChatSummary{}, fmt.Errorf("chat persistence not initialized")
	}
	return a.chatPersistence.GetChatList()
}

// GetActiveChatMeta returns the metadata for the active chat.
func (a *App) GetActiveChatMeta() (*chat.ChatMeta, error) {
	if a.chatPersistence == nil {
		return nil, fmt.Errorf("chat persistence not initialized")
	}
	return a.chatPersistence.GetActiveChatMeta()
}

// NewChat creates a new chat with the given name.
func (a *App) NewChat(name string) (*chat.ChatMeta, error) {
	if a.chatPersistence == nil {
		return nil, fmt.Errorf("chat persistence not initialized")
	}
	chatData, err := a.chatPersistence.NewChat(name)
	if err != nil {
		a.logError(errors.CHATLoadFail, errors.High, "chat.new", "app-chat.go:39", err)
		return nil, err
	}
	a.chatMessages = []map[string]interface{}{}
	a.triggeredLore = map[string]bool{}
	if a.memoryService != nil {
		a.memoryService.SetChatID(chatData.Meta.ID)
	}
	return &chatData.Meta, nil
}

// SaveChat consolidates the append log into a canonical snapshot.
func (a *App) SaveChat() error {
	if a.chatPersistence == nil {
		return fmt.Errorf("chat persistence not initialized")
	}
	chatID := a.chatPersistence.GetLastChatID()
	if chatID == "" {
		return fmt.Errorf("no active chat to save")
	}

	history := a.buildChatHistory()
	chatMeta, err := a.chatPersistence.GetActiveChatMeta()
	if err != nil {
		a.logErrorSilent(errors.CHATSaveFail, errors.Medium, "chat.save", "app-chat.go:62", err)
		return a.performSave()
	}

	lorebookEntries := make([]interface{}, 0, len(a.lorebooks))
	for _, lb := range a.lorebooks {
		lorebookEntries = append(lorebookEntries, map[string]interface{}{
			"id": lb.ID, "title": lb.Title, "content": lb.Content,
		})
	}

	chatData := &chat.ChatData{
		Meta:       *chatMeta,
		History:    history,
		World:      a.world,
		Scenario:   a.scenario,
		Characters: a.characters,
		Lorebooks:  lorebookEntries,
	}
	chatData.Meta.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return a.chatPersistence.SaveChat(chatData)
}

// SaveChatAs creates a copy of the current chat with a new name.
func (a *App) SaveChatAs(name string) (*chat.ChatMeta, error) {
	if a.chatPersistence == nil {
		return nil, fmt.Errorf("chat persistence not initialized")
	}
	lastChatID := a.chatPersistence.GetLastChatID()
	if lastChatID == "" {
		return nil, fmt.Errorf("no active chat to save")
	}
	chatData, err := a.chatPersistence.SaveChatAs(lastChatID, name)
	if err != nil {
		a.logError(errors.CHATSaveFail, errors.High, "chat.saveas", "app-chat.go:95", err)
		return nil, err
	}
	if a.memoryService != nil {
		a.memoryService.SetChatID(chatData.Meta.ID)
	}
	return &chatData.Meta, nil
}

// LoadChat loads a chat by ID, saving the current session first.
func (a *App) LoadChat(chatID string) (*chat.ChatData, error) {
	if a.chatPersistence == nil {
		return nil, fmt.Errorf("chat persistence not initialized")
	}
	if lastChatID := a.chatPersistence.GetLastChatID(); lastChatID != "" && lastChatID != chatID {
		if err := a.performSave(); err != nil {
			fmt.Printf("Warning: Failed to save current chat: %v\n", err)
		}
	}
	chatData, err := a.chatPersistence.LoadChat(chatID)
	if err != nil {
		a.logError(errors.CHATLoadFail, errors.High, "chat.load", "app-chat.go:115", err)
		return nil, err
	}
	a.updateAppStateFromChat(chatData)
	if err := a.chatPersistence.SetLastChatID(chatID); err != nil {
		return nil, fmt.Errorf("failed to set active chat: %w", err)
	}
	if a.memoryService != nil {
		a.memoryService.SetChatID(chatID)
	}
	return chatData, nil
}

// DeleteChat deletes a chat by ID and wipes its directory from disk.
func (a *App) DeleteChat(chatID string) error {
	if a.chatPersistence == nil {
		return fmt.Errorf("chat persistence not initialized")
	}
	if a.chatPersistence.GetLastChatID() == chatID {
		a.world = nil
		a.scenario = nil
		a.characters = []*chat.CharacterContext{}
		a.lorebooks = nil
		a.chatMessages = []map[string]interface{}{}
		a.triggeredLore = map[string]bool{}
	}
	return a.chatPersistence.DeleteChat(chatID)
}

// ClearWorldFromSession removes the active world and persists.
func (a *App) ClearWorldFromSession() error {
	a.world = nil
	if a.chatPersistence != nil {
		if chatID := a.chatPersistence.GetLastChatID(); chatID != "" {
			_ = a.performSave()
		}
	}
	return nil
}

// ClearScenarioFromSession removes the active scenario and persists.
func (a *App) ClearScenarioFromSession() error {
	a.scenario = nil
	if a.chatPersistence != nil {
		if chatID := a.chatPersistence.GetLastChatID(); chatID != "" {
			_ = a.performSave()
		}
	}
	return nil
}

// ExitApp performs a final save and quits the application.
func (a *App) ExitApp() error {
	if a.chatPersistence != nil {
		a.performFinalSave()
	}
	runtime.Quit(a.ctx)
	return nil
}

// ---------------------------------------------------------------------------
// Message send bindings
// ---------------------------------------------------------------------------

// SendMessage sends a message and returns the AI response synchronously.
func (a *App) SendMessage(message string) (string, error) {
	if a.httpClient == nil {
		return "", fmt.Errorf("llama client not initialized")
	}
	a.chatMessages = append(a.chatMessages, map[string]interface{}{
		"role": "user", "content": message, "timestamp": time.Now().Unix(),
	})
	response, err := a.httpClient.Generate(a.buildPrompt(), a.buildInferenceOptions())
	if err != nil {
		return "", err
	}
	a.chatMessages = append(a.chatMessages, map[string]interface{}{
		"role": "assistant", "content": response, "timestamp": time.Now().Unix(),
	})
	if a.chatPersistence != nil {
		go func() {
			if err := a.performSave(); err != nil {
				fmt.Printf("Autosave failed: %v\n", err)
			}
		}()
	}
	return response, nil
}

// SendMessageStream sends a message and streams the response via Wails events.
func (a *App) SendMessageStream(message string) error {
	return a.SendMessageStreamWithImage(message, "")
}

// ---------------------------------------------------------------------------
// Chat data Wails bindings (frontend ↔ backend state sync)
// ---------------------------------------------------------------------------

// LoadChatData returns the current backend chat state as a frontend-compatible map.
func (a *App) LoadChatData() (map[string]interface{}, error) {
	if a.world != nil || a.scenario != nil || len(a.characters) > 0 || len(a.lorebooks) > 0 {
		var worldMap map[string]interface{}
		if a.world != nil {
			imagePath := ""
			if a.world.Image != nil {
				imagePath = a.world.Image.Path
			}
			worldMap = map[string]interface{}{
				"id": a.world.ID, "name": a.world.Name,
				"description": a.world.Description, "image": imagePath,
			}
		}

		var scenarioMap map[string]interface{}
		if a.scenario != nil {
			imagePath := ""
			if a.scenario.Image != nil {
				imagePath = a.scenario.Image.Path
			}
			scenarioMap = map[string]interface{}{
				"id": a.scenario.ID, "name": a.scenario.Name,
				"description": a.scenario.Description, "image": imagePath,
			}
		}

		charactersSlice := make([]map[string]interface{}, len(a.characters))
		for i, char := range a.characters {
			imagePath := ""
			if char.Image != nil {
				imagePath = char.Image.Path
			}
			charactersSlice[i] = map[string]interface{}{
				"id": char.ID, "name": char.Name,
				"personality": char.Personality, "backstory": char.Backstory,
				"image": imagePath,
			}
		}

		lorebooksInterface := make([]interface{}, len(a.lorebooks))
		for i, lb := range a.lorebooks {
			lorebooksInterface[i] = map[string]interface{}{
				"id": lb.ID, "name": lb.Name, "title": lb.Title,
				"description": lb.Description, "content": lb.Content,
				"tags": lb.Tags, "scope": lb.Scope, "type": lb.Type,
				"ownerId": lb.OwnerID, "enabled": lb.Enabled,
				"triggerPhrases": lb.TriggerPhrases, "secondaryTriggers": lb.SecondaryTriggers,
				"triggerMode": lb.TriggerMode, "triggerDirection": lb.TriggerDirection,
				"triggerFrequency": lb.TriggerFrequency, "priorityLevel": lb.PriorityLevel,
				"maxLength": lb.MaxLength, "contextBudget": lb.ContextBudget,
				"scanWindow": lb.ScanWindow, "automated": lb.Automated,
			}
		}
		return map[string]interface{}{
			"characters": charactersSlice, "world": worldMap,
			"scenario": scenarioMap, "lorebooks": lorebooksInterface, "chatName": "New Chat",
		}, nil
	}
	if data, err := a.loadChatDataFromStorage(); err == nil && data != nil {
		a.applyChatDataToState(data)
		return data, nil
	}
	return map[string]interface{}{
		"characters": []interface{}{}, "world": nil,
		"scenario": nil, "lorebooks": []interface{}{}, "chatName": "New Chat",
	}, nil
}

// SaveChatData persists chat session data from the frontend.
func (a *App) SaveChatData(data map[string]interface{}) error {
	a.applyChatDataToState(data)
	chatData, err := a.convertFrontendDataToChatData(data)
	if err != nil {
		return fmt.Errorf("failed to convert frontend data: %w", err)
	}
	if a.chatPersistence == nil {
		return nil
	}
	return a.chatPersistence.SaveChat(chatData)
}

// ---------------------------------------------------------------------------
// Internal save helpers
// ---------------------------------------------------------------------------

func (a *App) performSave() error {
	lastChatID := a.chatPersistence.GetLastChatID()
	if lastChatID == "" {
		return fmt.Errorf("no active chat to save")
	}
	chatData, err := a.chatPersistence.LoadChat(lastChatID)
	if err != nil {
		return err
	}
	a.updateChatDataFromAppState(chatData)
	return a.chatPersistence.SaveChat(chatData)
}

func (a *App) performFinalSave() {
	if err := a.performSave(); err != nil {
		fmt.Printf("Final save failed: %v\n", err)
	}
}

func (a *App) buildChatHistory() []chat.ChatHistory {
	history := make([]chat.ChatHistory, 0, len(a.chatMessages))
	for _, msg := range a.chatMessages {
		role, _ := msg["role"].(string)
		content, _ := msg["content"].(string)
		if strings.TrimSpace(content) == "" {
			continue
		}
		sender := "System"
		switch role {
		case "user":
			if s, ok := msg["sender"].(string); ok && s != "" {
				sender = s
			} else {
				sender = "User"
			}
		case "assistant":
			if c, ok := msg["character"].(string); ok && c != "" {
				sender = c
			} else if s, ok := msg["sender"].(string); ok && s != "" {
				sender = s
			} else {
				sender = "Assistant"
			}
		}
		history = append(history, chat.ChatHistory{
			ID:        fmt.Sprintf("%d", len(history)),
			Sender:    sender,
			Text:      content,
			Role:      role,
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}
	return history
}

func (a *App) updateChatDataFromAppState(chatData *chat.ChatData) {
	history := a.buildChatHistory()
	if len(history) > 0 || len(chatData.History) == 0 {
		chatData.History = history
	}
	chatData.Meta.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if a.world != nil {
		chatData.World = a.world
	}
	if a.scenario != nil {
		chatData.Scenario = a.scenario
	}
	if a.characters != nil {
		chatData.Characters = a.characters
	}

	lorebookItems := make([]interface{}, 0, len(a.lorebooks))
	lorebookNames := make([]string, 0, len(a.lorebooks))
	for _, lb := range a.lorebooks {
		lorebookItems = append(lorebookItems, map[string]interface{}{
			"id": lb.ID, "title": lb.Title, "content": lb.Content, "tags": lb.Tags,
		})
		if lb.Title != "" {
			lorebookNames = append(lorebookNames, lb.Title)
		}
	}
	chatData.Lorebooks = lorebookItems
	if chatData.Memories == nil {
		chatData.Memories = []interface{}{}
	}
	if chatData.Meta.WorldName == "" && a.world != nil {
		chatData.Meta.WorldName = a.world.Name
	}
	if chatData.Meta.ScenarioName == "" && a.scenario != nil {
		chatData.Meta.ScenarioName = a.scenario.Name
	}
	if len(a.characters) > 0 {
		charNames := make([]string, 0, len(a.characters))
		for _, c := range a.characters {
			if c != nil && c.Name != "" {
				charNames = append(charNames, c.Name)
			}
		}
		chatData.Meta.CharacterNames = charNames
	}
	chatData.Meta.LorebookNames = lorebookNames
}

func (a *App) updateAppStateFromChat(chatData *chat.ChatData) {
	a.chatMessages = []map[string]interface{}{}
	for _, h := range chatData.History {
		msg := map[string]interface{}{
			"role": h.Role, "content": h.Text,
		}
		// Preserve sender/character info so frontend can display who said what
		if h.Sender != "" && h.Sender != "User" && h.Sender != "Assistant" && h.Sender != "System" {
			if h.Role == "assistant" {
				msg["character"] = h.Sender
			} else {
				msg["sender"] = h.Sender
			}
		} else if h.Sender != "" && h.Sender != "System" {
			msg["sender"] = h.Sender
		}
		a.chatMessages = append(a.chatMessages, msg)
	}

	a.world = chatData.World
	a.scenario = chatData.Scenario
	a.characters = chatData.Characters
	a.triggeredLore = map[string]bool{}
}

// ---------------------------------------------------------------------------
// Frontend data conversion helpers
// ---------------------------------------------------------------------------

func (a *App) convertFrontendDataToChatData(data map[string]interface{}) (*chat.ChatData, error) {
	chatData := &chat.ChatData{}

	if meta, ok := data["meta"].(map[string]interface{}); ok {
		chatData.Meta = chat.ChatMeta{
			ID:             a.getStringFromMap(meta, "id"),
			Name:           a.getStringFromMap(meta, "name"),
			CreatedAt:      a.getStringFromMap(meta, "createdAt"),
			UpdatedAt:      a.getStringFromMap(meta, "updatedAt"),
			WorldName:      a.getStringFromMap(meta, "worldName"),
			ScenarioName:   a.getStringFromMap(meta, "scenarioName"),
			CharacterNames: a.getStringSliceFromMap(meta, "characterNames"),
			LorebookNames:  a.getStringSliceFromMap(meta, "lorebookNames"),
		}
	}
	if historyData, ok := data["history"].([]interface{}); ok {
		chatData.History = make([]chat.ChatHistory, 0, len(historyData))
		for _, item := range historyData {
			if m, ok := item.(map[string]interface{}); ok {
				chatData.History = append(chatData.History, chat.ChatHistory{
					ID: a.getStringFromMap(m, "id"), Sender: a.getStringFromMap(m, "sender"),
					Text: a.getStringFromMap(m, "text"), Role: a.getStringFromMap(m, "role"),
					Timestamp: a.getStringFromMap(m, "timestamp"), HasImage: a.getBoolFromMap(m, "hasImage"),
				})
			}
		}
	}
	if worldDataMap, ok := data["world"].(map[string]interface{}); ok && worldDataMap != nil {
		world, err := a.parseWorldContextFromMap(worldDataMap)
		if err != nil {
			return nil, fmt.Errorf("failed to parse world data: %w", err)
		}
		chatData.World = &world
	}
	if scenarioDataMap, ok := data["scenario"].(map[string]interface{}); ok && scenarioDataMap != nil {
		scenario, err := a.parseScenarioContextFromMap(scenarioDataMap)
		if err != nil {
			return nil, fmt.Errorf("failed to parse scenario data: %w", err)
		}
		chatData.Scenario = &scenario
	}
	if charactersDataSlice, ok := data["characters"].([]interface{}); ok {
		var characters []*chat.CharacterContext
		for _, charItem := range charactersDataSlice {
			if charMap, ok := charItem.(map[string]interface{}); ok {
				character, err := a.parseCharacterContextFromMap(charMap)
				if err != nil {
					return nil, fmt.Errorf("failed to parse character data: %w", err)
				}
				characters = append(characters, &character)
			}
		}
		chatData.Characters = characters
	}
	if lorebooksData, ok := data["lorebooks"].([]interface{}); ok {
		chatData.Lorebooks = lorebooksData
	}
	if memoriesData, ok := data["memories"].([]interface{}); ok {
		chatData.Memories = memoriesData
	}
	if chatData.Meta.ID == "" {
		return nil, fmt.Errorf("chat data missing required meta.id field")
	}
	return chatData, nil
}

// ---------------------------------------------------------------------------
// Map extraction helpers
// ---------------------------------------------------------------------------

func (a *App) getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func (a *App) getIntFromMap(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		}
	}
	return 0
}

func (a *App) getBoolFromMap(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func (a *App) getStringSliceFromMap(m map[string]interface{}, key string) []string {
	if slice, ok := m[key].([]interface{}); ok {
		result := make([]string, 0, len(slice))
		for _, item := range slice {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return []string{}
}

func (a *App) parseImageFromValue(imageVal interface{}) *chat.Image {
	return chat.ParseImageFromValue(imageVal)
}

func (a *App) parseWorldContextFromMap(worldMap map[string]interface{}) (chat.WorldContext, error) {
	return chat.ParseWorldContextFromMap(worldMap)
}

func (a *App) parseScenarioContextFromMap(scenarioMap map[string]interface{}) (chat.ScenarioContext, error) {
	return chat.ParseScenarioContextFromMap(scenarioMap)
}

func (a *App) parseCharacterContextFromMap(characterMap map[string]interface{}) (chat.CharacterContext, error) {
	return chat.ParseCharacterContextFromMap(characterMap)
}

// ---------------------------------------------------------------------------
// Lorebook helpers used by app-characters.go
// ---------------------------------------------------------------------------

func (a *App) parseLorebookEntryFromMap(m map[string]interface{}) (lorebook.Entry, error) {
	return lorebook.ParseLorebookEntryFromMap(m)
}
