package main

import (
	"fmt"
	"os/exec"
	goruntime "runtime"

	"howl-chat/internal/audio"
	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/errors"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// TODO: Implement HuggingFace repo link support for chat models
// - Add HuggingFace repo URL field to model metadata
// - Display repo link in model selection UI (like professional apps)
// - Allow users to click through to view model details, documentation, and download alternatives
// - Consider adding model card preview or metadata fetching from HuggingFace API
// - Reference: https://huggingface.co/docs/huggingface_hub/guides/huggingface_hub_py

// ---------------------------------------------------------------------------
// Chat Persistence Wails Bindings
//
// These bindings are thin: the collection, its membership and the append log are
// owned by chat.Conversation, and every route below asks that owner. App holds the
// handle and translates the renderer's shapes.
// ---------------------------------------------------------------------------

// GetChatList returns a list of all chat summaries.
func (a *App) GetChatList() ([]chat.ChatSummary, error) {
	if !a.conversation.Bound() {
		return []chat.ChatSummary{}, fmt.Errorf("chat persistence not initialized")
	}
	return a.conversation.List()
}

// GetActiveChatMeta returns the metadata for the active chat.
func (a *App) GetActiveChatMeta() (*chat.ChatMeta, error) {
	if !a.conversation.Bound() {
		return nil, fmt.Errorf("chat persistence not initialized")
	}
	return a.conversation.Meta()
}

// NewChat creates a new chat with the given name.
func (a *App) NewChat(name string) (*chat.ChatMeta, error) {
	if !a.conversation.Bound() {
		return nil, fmt.Errorf("chat persistence not initialized")
	}
	meta, err := a.conversation.New(name)
	if err != nil {
		a.logError(errors.CHATLoadFail, errors.High, "chat.new", "app-chat.go:39", err)
		return nil, err
	}
	if a.memoryService != nil {
		a.memoryService.SetChatID(meta.ID)
	}
	return meta, nil
}

// SaveChat consolidates the append log into a canonical snapshot. Membership and
// tracking arrays are recomputed by the conversation from its session — never read
// from stale meta on disk.
func (a *App) SaveChat() error {
	if !a.conversation.Bound() {
		return fmt.Errorf("chat persistence not initialized")
	}
	return a.performSave()
}

// SaveChatAs creates a copy of the current chat with a new name.
func (a *App) SaveChatAs(name string) (*chat.ChatMeta, error) {
	if !a.conversation.Bound() {
		return nil, fmt.Errorf("chat persistence not initialized")
	}
	meta, err := a.conversation.SaveAs(name)
	if err != nil {
		a.logError(errors.CHATSaveFail, errors.High, "chat.saveas", "app-chat.go:95", err)
		return nil, err
	}
	if a.memoryService != nil {
		a.memoryService.SetChatID(meta.ID)
	}
	return meta, nil
}

// LoadChat loads a chat by ID, saving the current session first.
func (a *App) LoadChat(chatID string) (*chat.ChatData, error) {
	if !a.conversation.Bound() {
		return nil, fmt.Errorf("chat persistence not initialized")
	}
	chatData, err := a.conversation.Load(chatID)
	if err != nil {
		a.logError(errors.CHATLoadFail, errors.High, "chat.load", "app-chat.go:115", err)
		return nil, err
	}
	if a.memoryService != nil {
		a.memoryService.SetChatID(chatID)
	}
	return chatData, nil
}

// DeleteChat deletes a chat by ID and wipes its directory from disk.
func (a *App) DeleteChat(chatID string) error {
	if !a.conversation.Bound() {
		return fmt.Errorf("chat persistence not initialized")
	}
	return a.conversation.Delete(chatID)
}

// GetChatDirectory reports where a conversation's files live, so the deletion
// confirmation can tell the user where to look if Windows blocks the delete.
func (a *App) GetChatDirectory(chatID string) (string, error) {
	if !a.conversation.Bound() {
		return "", fmt.Errorf("chat persistence not initialized")
	}
	return a.conversation.Directory(chatID)
}

// OpenChatDirectory opens a conversation's directory in the system file browser.
// It is the clickable half of that notice; a failure to open is reported to the
// caller but never blocks deletion.
func (a *App) OpenChatDirectory(chatID string) error {
	dir, err := a.GetChatDirectory(chatID)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "windows":
		cmd = exec.Command("explorer.exe", dir)
	case "darwin":
		cmd = exec.Command("open", dir)
	default:
		cmd = exec.Command("xdg-open", dir)
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("[OpenChatDirectory] failed to open %s: %v\n", dir, err)
		return fmt.Errorf("failed to open %s: %w", dir, err)
	}
	fmt.Printf("[OpenChatDirectory] opened %s\n", dir)
	// Explorer reports a non-zero exit code even on success, so the process is
	// reaped without treating its status as a result.
	go func() { _ = cmd.Wait() }()
	return nil
}

// ClearWorldFromSession removes the active world and persists.
func (a *App) ClearWorldFromSession() error {
	a.session().ClearWorld()
	if a.conversation.ActiveID() != "" {
		a.performSaveLogged("app-chat.go")
	}
	return nil
}

// ClearScenarioFromSession removes the active scenario and persists.
func (a *App) ClearScenarioFromSession() error {
	a.session().ClearScenario()
	if a.conversation.ActiveID() != "" {
		a.performSaveLogged("app-chat.go")
	}
	return nil
}

// ExitApp performs the same idempotent shutdown as window close, then quits.
func (a *App) ExitApp() error {
	if err := a.shutdown(nil); err != nil {
		fmt.Printf("Shutdown failed: %v\n", err)
	}
	runtime.Quit(a.ctx)
	return nil
}

// ---------------------------------------------------------------------------
// Message send bindings
// ---------------------------------------------------------------------------

// SendMessage sends a message and returns the AI response synchronously.
func (a *App) SendMessage(message string) (string, error) {
	if err := a.ensureAcceptingWork("send message"); err != nil {
		return "", err
	}
	if a.engine.Client() == nil {
		return "", fmt.Errorf("llama client not initialized")
	}
	response, err := a.turn().Complete(nil, message, chat.Attribution{Sender: "AI"})
	if err != nil {
		return "", err
	}
	if a.conversation.Bound() {
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
	world := a.session().World()
	scenario := a.session().Scenario()
	characters := a.session().Characters()
	lorebooks := a.session().Lorebooks()
	if world != nil || scenario != nil || len(characters) > 0 || len(lorebooks) > 0 {
		var worldMap map[string]interface{}
		if world != nil {
			imagePath := ""
			if world.Image != nil {
				imagePath = world.Image.Path
			}
			worldMap = map[string]interface{}{
				"id": world.ID, "name": world.Name,
				"description": world.Description, "image": imagePath,
			}
		}

		var scenarioMap map[string]interface{}
		if scenario != nil {
			imagePath := ""
			if scenario.Image != nil {
				imagePath = scenario.Image.Path
			}
			scenarioMap = map[string]interface{}{
				"id": scenario.ID, "name": scenario.Name,
				"description": scenario.Description, "image": imagePath,
			}
		}

		charactersSlice := make([]map[string]interface{}, len(characters))
		for i, char := range characters {
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

		lorebooksInterface := make([]interface{}, len(lorebooks))
		for i, lb := range lorebooks {
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
				"stale": lb.Stale,
			}
		}
		meta := map[string]interface{}{
			"id":                "",
			"name":              "New Chat",
			"worldId":           "",
			"scenarioId":        "",
			"activeCharacters":  []string{},
			"removedCharacters": []string{},
			"addedCharacters":   []string{},
			"activeLorebooks":   []string{},
			"removedLorebooks":  []string{},
			"addedLorebooks":    []string{},
		}
		return map[string]interface{}{
			"characters": charactersSlice, "world": worldMap,
			"scenario": scenarioMap, "lorebooks": lorebooksInterface,
			"chatName": "New Chat", "meta": meta,
		}, nil
	}
	if data, err := a.loadChatDataFromStorage(); err == nil && data != nil {
		a.applyChatDataToState(data)
		return data, nil
	}
	return map[string]interface{}{
		"characters": []interface{}{}, "world": nil,
		"scenario": nil, "lorebooks": []interface{}{}, "chatName": "New Chat",
		"meta": map[string]interface{}{
			"activeCharacters":  []string{},
			"removedCharacters": []string{},
			"addedCharacters":   []string{},
			"activeLorebooks":   []string{},
			"removedLorebooks":  []string{},
			"addedLorebooks":    []string{},
		},
	}, nil
}

// SaveChatData persists chat session data from the frontend.
func (a *App) SaveChatData(data map[string]interface{}) error {
	a.applyChatDataToState(data)
	chatData, err := a.convertFrontendDataToChatData(data)
	if err != nil {
		return fmt.Errorf("failed to convert frontend data: %w", err)
	}
	// The frontend payload is not evidence that a conversation exists. If the
	// membership authority does not claim this id and its folder, the payload is a
	// cached or stale session and is refused instead of becoming a chat.
	return a.conversation.SaveRecord(chatData)
}

// ---------------------------------------------------------------------------
// Internal save helpers
// ---------------------------------------------------------------------------

func (a *App) performSave() error {
	if !a.conversation.Bound() {
		return fmt.Errorf("chat persistence not initialized")
	}
	return a.conversation.Save()
}

// performSaveLogged persists and reports a failure instead of discarding it.
// A refusal (busy, re-entrant, deleted) is a real outcome, and silence about it
// is exactly what V3-27 makes non-conforming.
func (a *App) performSaveLogged(site string) {
	if err := a.performSave(); err != nil {
		a.logError(errors.CHATSaveFail, errors.High, "chat.save", site, err)
	}
}

// performFinalSave persists the active conversation once, on shutdown.
//
// Companions, then the log, then the snapshot are written under a single lock
// hold, so stale log lines, deleted messages and legacy DOM-scrape entries cannot
// survive the rewrite, and a chat that was never loaded keeps its stored history
// instead of being truncated.
func (a *App) performFinalSave() error {
	count, err := a.conversation.FinalSave()
	if err != nil {
		return err
	}
	if count > 0 || a.conversation.ActiveID() != "" {
		fmt.Println("[FinalSave] committed history with", count, "entries")
	}
	return nil
}

// sessionEntriesToMaps projects session turns into the renderer's message shape.
// The binding contract is unchanged; only its source moved into the domain.
func sessionEntriesToMaps(entries []chat.SessionEntry) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		msg := map[string]interface{}{
			"role":    entry.Role,
			"content": entry.Content,
		}
		if entry.Sender != "" {
			msg["sender"] = entry.Sender
		}
		if entry.Character != "" {
			msg["character"] = entry.Character
		}
		if entry.Timestamp != "" {
			msg["timestamp"] = entry.Timestamp
		}
		if entry.HasImage {
			msg["hasImage"] = true
		}
		if entry.CharacterID != "" {
			msg["characterId"] = entry.CharacterID
		}
		out = append(out, msg)
	}
	return out
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
			// Source-of-truth tracking arrays
			WorldID:           a.getStringFromMap(meta, "worldId"),
			ScenarioID:        a.getStringFromMap(meta, "scenarioId"),
			ActiveCharacters:  a.getStringSliceFromMap(meta, "activeCharacters"),
			RemovedCharacters: a.getStringSliceFromMap(meta, "removedCharacters"),
			AddedCharacters:   a.getStringSliceFromMap(meta, "addedCharacters"),
			ActiveLorebooks:   a.getStringSliceFromMap(meta, "activeLorebooks"),
			RemovedLorebooks:  a.getStringSliceFromMap(meta, "removedLorebooks"),
			AddedLorebooks:    a.getStringSliceFromMap(meta, "addedLorebooks"),
		}
	}
	if historyData, ok := data["history"].([]interface{}); ok {
		chatData.History = make([]chat.ChatHistory, 0, len(historyData))
		for _, item := range historyData {
			if m, ok := item.(map[string]interface{}); ok {
				characterID := a.getStringFromMap(m, "characterId")
				if characterID == "" {
					characterID = a.getStringFromMap(m, "CharacterId")
				}
				chatData.History = append(chatData.History, chat.ChatHistory{
					ID: a.getStringFromMap(m, "id"), Sender: a.getStringFromMap(m, "sender"),
					Text: a.getStringFromMap(m, "text"), Role: a.getStringFromMap(m, "role"),
					Timestamp: a.getStringFromMap(m, "timestamp"), HasImage: a.getBoolFromMap(m, "hasImage"),
					CharacterID: characterID,
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

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
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

// ---------------------------------------------------------------------------
// Unified Settings Store
// ---------------------------------------------------------------------------

func (a *App) GetSettings() map[string]interface{} {
	return a.settings.All()
}

func (a *App) SetSettings(settings map[string]interface{}) error {
	if a.settings == nil {
		return fmt.Errorf("settings store not initialized")
	}
	// The authority replaces the document and the sampler/model sections it
	// carries, then persists it. The sections are REPLACED, not merged: callers
	// send the whole settings object, so a key missing from it means the user's
	// value is gone.
	if err := a.settings.Replace(settings); err != nil {
		return err
	}
	// Saving is not applying. Keep the complete requested contract current so
	// status can distinguish a saved value from the active generation.
	if rawAudio, ok := settings["audio"].(map[string]interface{}); ok {
		a.audio.SetDesired(audio.DesiredSettingsFromMap(rawAudio))
		a.audio.SetApplyError("")
	}

	return nil
}
