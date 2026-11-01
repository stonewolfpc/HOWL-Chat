package chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PersistenceManager handles chat persistence with the specified directory structure
type PersistenceManager struct {
	rootDir string
	config  AppConfig
}

// NewPersistenceManagerWithDir creates a persistence manager rooted at an arbitrary directory.
// Intended for tests; production code should use NewPersistenceManager.
func NewPersistenceManagerWithDir(rootDir string) (*PersistenceManager, error) {
	pm := &PersistenceManager{rootDir: rootDir, config: AppConfig{}}
	if err := pm.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize persistence: %w", err)
	}
	return pm, nil
}

// NewPersistenceManager creates a new persistence manager
func NewPersistenceManager() (*PersistenceManager, error) {
	// Get app data directory
	rootDir, err := getAppDataDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get app data directory: %w", err)
	}

	pm := &PersistenceManager{
		rootDir: rootDir,
		config:  AppConfig{},
	}

	// Initialize directories and load config
	if err := pm.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize persistence: %w", err)
	}

	return pm, nil
}

// getAppDataDir returns the app data directory based on OS
func getAppDataDir() (string, error) {
	// For Windows: C:\Users\<User>\AppData\Roaming\HOWLChat
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, "AppData", "Roaming", "HOWLChat"), nil
}

// initialize creates directories and loads config
func (pm *PersistenceManager) initialize() error {
	// Create root directory
	if err := os.MkdirAll(pm.rootDir, 0755); err != nil {
		return fmt.Errorf("failed to create root directory: %w", err)
	}

	// Create chats directory
	chatsDir := filepath.Join(pm.rootDir, "chats")
	if err := os.MkdirAll(chatsDir, 0755); err != nil {
		return fmt.Errorf("failed to create chats directory: %w", err)
	}

	// Load or create config
	if err := pm.loadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return nil
}

// loadConfig loads the application config
func (pm *PersistenceManager) loadConfig() error {
	configPath := filepath.Join(pm.rootDir, "config.json")

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		pm.config = AppConfig{LastChatID: ""}
		return pm.saveConfig()
	}

	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, &pm.config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// saveConfig saves the application config
func (pm *PersistenceManager) saveConfig() error {
	configPath := filepath.Join(pm.rootDir, "config.json")

	data, err := json.MarshalIndent(pm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// generateChatID generates a new chat ID based on timestamp
func (pm *PersistenceManager) generateChatID() string {
	now := time.Now().UTC()
	return fmt.Sprintf("%d-%02d-%02dT%02d-%02d-%02dZ",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
}

// NewChat creates a new chat with the given name
func (pm *PersistenceManager) NewChat(name string) (*ChatData, error) {
	chatID := pm.generateChatID()
	chatDir := filepath.Join(pm.rootDir, "chats", chatID)

	// Create chat directory
	if err := os.MkdirAll(chatDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create chat directory: %w", err)
	}

	// Create chat data
	now := time.Now().UTC()
	chatData := &ChatData{
		Meta: ChatMeta{
			ID:             chatID,
			Name:           name,
			CreatedAt:      now.Format(time.RFC3339),
			UpdatedAt:      now.Format(time.RFC3339),
			WorldName:      "",
			ScenarioName:   "",
			CharacterNames: []string{},
			LorebookNames:  []string{},
		},
		History:    []ChatHistory{},
		World:      nil,
		Scenario:   nil,
		Characters: []*CharacterContext{},
		Lorebooks:  []interface{}{},
		Memories:   []interface{}{},
	}

	// Save initial chat data
	if err := pm.SaveChat(chatData); err != nil {
		return nil, fmt.Errorf("failed to save new chat: %w", err)
	}

	// Update config to set as last chat
	pm.config.LastChatID = chatID
	if err := pm.saveConfig(); err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	return chatData, nil
}

// SaveChat saves all chat data to files
func (pm *PersistenceManager) SaveChat(chatData *ChatData) error {
	chatDir := filepath.Join(pm.rootDir, "chats", chatData.Meta.ID)

	// Ensure chat directory exists
	if err := os.MkdirAll(chatDir, 0755); err != nil {
		return fmt.Errorf("failed to create chat directory: %w", err)
	}

	// Update updated timestamp
	chatData.Meta.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	// Save meta.json
	if err := pm.saveJSON(filepath.Join(chatDir, "meta.json"), chatData.Meta); err != nil {
		return fmt.Errorf("failed to save meta.json: %w", err)
	}

	// Save canonical history.json snapshot (atomic write) for fallback/backup only.
	// The append log (history.ndjson) is the true source of truth and must NOT
	// be deleted — frontend DOM scraping can corrupt history.json roles/ids.
	if err := pm.saveJSON(filepath.Join(chatDir, "history.json"), chatData.History); err != nil {
		return fmt.Errorf("failed to save history.json: %w", err)
	}

	// Save world.json
	if chatData.World != nil {
		if err := pm.saveJSON(filepath.Join(chatDir, "world.json"), chatData.World); err != nil {
			return fmt.Errorf("failed to save world.json: %w", err)
		}
	}

	// Save scenario.json
	if chatData.Scenario != nil {
		if err := pm.saveJSON(filepath.Join(chatDir, "scenario.json"), chatData.Scenario); err != nil {
			return fmt.Errorf("failed to save scenario.json: %w", err)
		}
	}

	// Save characters.json
	if err := pm.saveJSON(filepath.Join(chatDir, "characters.json"), chatData.Characters); err != nil {
		return fmt.Errorf("failed to save characters.json: %w", err)
	}

	// Save lorebooks.json
	if err := pm.saveJSON(filepath.Join(chatDir, "lorebooks.json"), chatData.Lorebooks); err != nil {
		return fmt.Errorf("failed to save lorebooks.json: %w", err)
	}

	// Save memories.json
	if err := pm.saveJSON(filepath.Join(chatDir, "memories.json"), chatData.Memories); err != nil {
		return fmt.Errorf("failed to save memories.json: %w", err)
	}

	// Update config
	pm.config.LastChatID = chatData.Meta.ID
	if err := pm.saveConfig(); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}

// LoadChat loads chat data by ID. Every piece is best-effort:
// missing or corrupt file → log and continue with empty/default.
// Only meta.json is required (it IS the chat identity).
func (pm *PersistenceManager) LoadChat(chatID string) (*ChatData, error) {
	chatDir := filepath.Join(pm.rootDir, "chats", chatID)

	// Check if chat directory exists
	if _, err := os.Stat(chatDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("chat not found: %s", chatID)
	}

	chatData := &ChatData{}

	// Load meta.json (required)
	if err := pm.loadJSON(filepath.Join(chatDir, "meta.json"), &chatData.Meta); err != nil {
		return nil, fmt.Errorf("failed to load meta.json: %w", err)
	}
	if chatData.Meta.ID == "" {
		chatData.Meta.ID = chatID
	}

	// Load history (optional — empty if missing)
	chatData.History = pm.loadHistoryFile(chatDir)

	// Load world.json (optional)
	if worldBytes, err := os.ReadFile(filepath.Join(chatDir, "world.json")); err == nil && len(worldBytes) > 0 {
		if world, err := decodeWorldContext(worldBytes); err == nil {
			chatData.World = world
		} else {
			fmt.Printf("[LoadChat] skipping corrupt world.json for %s: %v\n", chatID, err)
		}
	}

	// Load scenario.json (optional)
	if scenarioBytes, err := os.ReadFile(filepath.Join(chatDir, "scenario.json")); err == nil && len(scenarioBytes) > 0 {
		if scenario, err := decodeScenarioContext(scenarioBytes); err == nil {
			chatData.Scenario = scenario
		} else {
			fmt.Printf("[LoadChat] skipping corrupt scenario.json for %s: %v\n", chatID, err)
		}
	}

	// Load characters.json (optional)
	if charactersBytes, err := os.ReadFile(filepath.Join(chatDir, "characters.json")); err == nil && len(charactersBytes) > 0 {
		var rawCharacters []json.RawMessage
		if err := json.Unmarshal(charactersBytes, &rawCharacters); err == nil {
			chatData.Characters = make([]*CharacterContext, 0, len(rawCharacters))
			for _, raw := range rawCharacters {
				if character, err := decodeCharacterContext(raw); err == nil && character != nil {
					chatData.Characters = append(chatData.Characters, character)
				}
			}
		} else {
			fmt.Printf("[LoadChat] skipping corrupt characters.json for %s: %v\n", chatID, err)
		}
	}

	// Load lorebooks.json (optional)
	_ = pm.loadJSON(filepath.Join(chatDir, "lorebooks.json"), &chatData.Lorebooks)

	// Load memories.json (optional)
	_ = pm.loadJSON(filepath.Join(chatDir, "memories.json"), &chatData.Memories)

	// Migrate legacy meta format → source-of-truth tracking arrays
	pm.migrateLegacyMeta(chatData)

	// Filter loaded characters/lorebooks by tracking arrays.
	// Final list = (active - removed) + added
	pm.filterChatDataByTracking(chatData)

	return chatData, nil
}

// filterChatDataByTracking filters Characters and Lorebooks so only
// entities that pass (active - removed) + added remain.
func (pm *PersistenceManager) filterChatDataByTracking(chatData *ChatData) {
	meta := chatData.Meta

	// Characters
	removedCharSet := make(map[string]struct{}, len(meta.RemovedCharacters))
	for _, id := range meta.RemovedCharacters {
		removedCharSet[id] = struct{}{}
	}
	addedCharSet := make(map[string]struct{}, len(meta.AddedCharacters))
	for _, id := range meta.AddedCharacters {
		addedCharSet[id] = struct{}{}
	}
	finalCharIds := make(map[string]struct{})
	for _, id := range meta.ActiveCharacters {
		if _, removed := removedCharSet[id]; !removed {
			finalCharIds[id] = struct{}{}
		}
	}
	for id := range addedCharSet {
		finalCharIds[id] = struct{}{}
	}

	filteredChars := make([]*CharacterContext, 0, len(chatData.Characters))
	for _, c := range chatData.Characters {
		if c != nil && c.ID != "" {
			if _, ok := finalCharIds[c.ID]; ok {
				filteredChars = append(filteredChars, c)
			}
		}
	}
	chatData.Characters = filteredChars

	// Lorebooks
	removedBookSet := make(map[string]struct{}, len(meta.RemovedLorebooks))
	for _, id := range meta.RemovedLorebooks {
		removedBookSet[id] = struct{}{}
	}
	addedBookSet := make(map[string]struct{}, len(meta.AddedLorebooks))
	for _, id := range meta.AddedLorebooks {
		addedBookSet[id] = struct{}{}
	}
	finalBookIds := make(map[string]struct{})
	for _, id := range meta.ActiveLorebooks {
		if _, removed := removedBookSet[id]; !removed {
			finalBookIds[id] = struct{}{}
		}
	}
	for id := range addedBookSet {
		finalBookIds[id] = struct{}{}
	}

	filteredBooks := make([]interface{}, 0, len(chatData.Lorebooks))
	for _, lb := range chatData.Lorebooks {
		if m, ok := lb.(map[string]interface{}); ok {
			if id, ok := m["id"].(string); ok && id != "" {
				if _, ok := finalBookIds[id]; ok {
					filteredBooks = append(filteredBooks, lb)
				}
				continue
			}
		}
		// If no id, keep it (legacy or inline lorebook)
		filteredBooks = append(filteredBooks, lb)
	}
	chatData.Lorebooks = filteredBooks
}

// migrateLegacyMeta populates new tracking fields from legacy characterNames /
// lorebookNames when the new arrays are empty.  This is a one-time migration
// for chats created before the tracking-array architecture.
func (pm *PersistenceManager) migrateLegacyMeta(chatData *ChatData) {
	meta := &chatData.Meta

	// Characters
	if len(meta.ActiveCharacters) == 0 {
		if len(chatData.Characters) > 0 {
			for _, c := range chatData.Characters {
				if c != nil && c.ID != "" {
					meta.ActiveCharacters = append(meta.ActiveCharacters, c.ID)
				}
			}
		}
		// Fallback: use legacy names as pseudo-IDs
		if len(meta.ActiveCharacters) == 0 && len(meta.CharacterNames) > 0 {
			meta.ActiveCharacters = append([]string{}, meta.CharacterNames...)
		}
	}

	// Lorebooks
	if len(meta.ActiveLorebooks) == 0 {
		for _, lb := range chatData.Lorebooks {
			if m, ok := lb.(map[string]interface{}); ok {
				if id, ok := m["id"].(string); ok && id != "" {
					meta.ActiveLorebooks = append(meta.ActiveLorebooks, id)
				}
			}
		}
		// Fallback: use legacy names as pseudo-IDs
		if len(meta.ActiveLorebooks) == 0 && len(meta.LorebookNames) > 0 {
			meta.ActiveLorebooks = append([]string{}, meta.LorebookNames...)
		}
	}

	// World / Scenario IDs (fallback to names for now)
	if meta.WorldID == "" && meta.WorldName != "" {
		meta.WorldID = meta.WorldName
	}
	if meta.ScenarioID == "" && meta.ScenarioName != "" {
		meta.ScenarioID = meta.ScenarioName
	}
}

// SaveChatAs creates a copy of the current chat with a new name
func (pm *PersistenceManager) SaveChatAs(currentChatID, newName string) (*ChatData, error) {
	// Load current chat
	currentChat, err := pm.LoadChat(currentChatID)
	if err != nil {
		return nil, fmt.Errorf("failed to load current chat: %w", err)
	}

	// Create new chat
	newChat, err := pm.NewChat(newName)
	if err != nil {
		return nil, fmt.Errorf("failed to create new chat: %w", err)
	}

	// Copy all data except ID and timestamps
	newChat.Meta.WorldName = currentChat.Meta.WorldName
	newChat.Meta.ScenarioName = currentChat.Meta.ScenarioName
	newChat.Meta.CharacterNames = currentChat.Meta.CharacterNames
	newChat.Meta.LorebookNames = currentChat.Meta.LorebookNames
	newChat.History = currentChat.History
	newChat.World = currentChat.World
	newChat.Scenario = currentChat.Scenario
	newChat.Characters = currentChat.Characters
	newChat.Lorebooks = currentChat.Lorebooks
	newChat.Memories = currentChat.Memories

	// Save the new chat
	if err := pm.SaveChat(newChat); err != nil {
		return nil, fmt.Errorf("failed to save new chat: %w", err)
	}

	return newChat, nil
}

// DeleteChat deletes a chat directory and all its contents.
func (pm *PersistenceManager) DeleteChat(chatID string) error {
	chatDir := filepath.Join(pm.rootDir, "chats", chatID)

	if err := os.RemoveAll(chatDir); err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	if pm.config.LastChatID == chatID {
		pm.config.LastChatID = ""
		_ = pm.saveConfig()
	}
	return nil
}

// GetChatList returns a list of all chat summaries
func (pm *PersistenceManager) GetChatList() ([]ChatSummary, error) {
	chatsDir := filepath.Join(pm.rootDir, "chats")

	// Read all chat directories
	entries, err := os.ReadDir(chatsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read chats directory: %w", err)
	}

	var summaries []ChatSummary

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		chatID := entry.Name()
		chatDir := filepath.Join(chatsDir, chatID)

		// Load meta.json
		var meta ChatMeta
		metaPath := filepath.Join(chatDir, "meta.json")
		if err := pm.loadJSON(metaPath, &meta); err != nil {
			// Skip invalid chat directories
			continue
		}

		// Load history to get message count (ignore error)
		var history []ChatHistory
		pm.loadJSON(filepath.Join(chatDir, "history.json"), &history)

		// Parse timestamps
		createdAt, _ := time.Parse(time.RFC3339, meta.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339, meta.UpdatedAt)

		summary := ChatSummary{
			ID:             meta.ID,
			Name:           meta.Name,
			CreatedAt:      createdAt.Format(time.RFC3339),
			UpdatedAt:      updatedAt.Format(time.RFC3339),
			WorldName:      meta.WorldName,
			ScenarioName:   meta.ScenarioName,
			CharacterNames: meta.CharacterNames,
			LorebookNames:  meta.LorebookNames,
			MessageCount:   len(history),
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// GetLastChatID returns the last chat ID from config
func (pm *PersistenceManager) GetLastChatID() string {
	return pm.config.LastChatID
}

// SetLastChatID sets the last chat ID in config
func (pm *PersistenceManager) SetLastChatID(chatID string) error {
	pm.config.LastChatID = chatID
	return pm.saveConfig()
}

// GetActiveChatMeta returns the metadata for the active chat by reading only
// meta.json — does NOT load history, world, characters etc.
func (pm *PersistenceManager) GetActiveChatMeta() (*ChatMeta, error) {
	if pm.config.LastChatID == "" {
		return nil, fmt.Errorf("no active chat")
	}
	metaPath := filepath.Join(pm.rootDir, "chats", pm.config.LastChatID, "meta.json")
	var meta ChatMeta
	if err := pm.loadJSON(metaPath, &meta); err != nil {
		return nil, fmt.Errorf("failed to load meta.json for active chat: %w", err)
	}
	if meta.ID == "" {
		meta.ID = pm.config.LastChatID
	}
	return &meta, nil
}

// Helper functions

// saveJSON atomically saves data to a JSON file.
// It writes to a temp file first then renames so a crash mid-write never corrupts the target.
func (pm *PersistenceManager) saveJSON(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	tmp := filename + ".tmp"
	if err := os.WriteFile(tmp, jsonData, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, filename)
}

// loadHistoryFile reads chat history, preferring the append log over the JSON snapshot.
// Missing or corrupt history → log and return empty slice. Never errors.
func (pm *PersistenceManager) loadHistoryFile(chatDir string) []ChatHistory {
	ndjsonPath := filepath.Join(chatDir, "history.ndjson")
	jsonPath := filepath.Join(chatDir, "history.json")

	// Prefer the append log when it exists and is non-empty
	if info, err := os.Stat(ndjsonPath); err == nil && info.Size() > 0 {
		f, err := os.Open(ndjsonPath)
		if err != nil {
			fmt.Printf("[loadHistoryFile] cannot open ndjson, returning empty: %v\n", err)
			return []ChatHistory{}
		}
		defer f.Close()

		var history []ChatHistory
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)
		for scanner.Scan() {
			line := bytes.TrimSpace(scanner.Bytes())
			if len(line) == 0 {
				continue
			}
			var entry ChatHistory
			if err := json.Unmarshal(line, &entry); err != nil {
				fmt.Printf("[loadHistoryFile] corrupt ndjson line, stopping: %v\n", err)
				break // return what we have so far
			}
			history = append(history, entry)
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("[loadHistoryFile] scanner error, returning partial: %v\n", err)
		}
		return history
	}

	// Fall back to the canonical JSON snapshot
	var history []ChatHistory
	if err := pm.loadJSON(jsonPath, &history); err != nil {
		// Missing or corrupt — just return empty
		return []ChatHistory{}
	}
	if history == nil {
		return []ChatHistory{}
	}
	return history
}

// RewriteHistory atomically replaces history.ndjson with entries, so deletions and
// edits actually remove messages. The previous delete path cleared in-memory state
// and re-appended the surviving messages on top of the log, which left the deleted
// lines in place for the next load to replay.
func (pm *PersistenceManager) RewriteHistory(chatID string, entries []ChatHistory) error {
	chatDir := filepath.Join(pm.rootDir, "chats", chatID)
	if err := os.MkdirAll(chatDir, 0755); err != nil {
		return fmt.Errorf("failed to create chat directory: %w", err)
	}

	var buf bytes.Buffer
	for _, e := range entries {
		line, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal history entry: %w", err)
		}
		buf.Write(line)
		buf.WriteByte(0x0A)
	}

	ndjsonPath := filepath.Join(chatDir, "history.ndjson")
	tmp := ndjsonPath + ".tmp"
	if err := os.WriteFile(tmp, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write history.ndjson: %w", err)
	}
	if err := os.Rename(tmp, ndjsonPath); err != nil {
		return fmt.Errorf("failed to replace history.ndjson: %w", err)
	}

	// Keep the JSON snapshot identical so a fallback load sees the same history.
	return pm.saveJSON(filepath.Join(chatDir, "history.json"), entries)
}

// AppendHistoryEntry appends a single entry to history.ndjson using an O(1) file append.
// This is the hot path called after every message — no full file rewrite.
// If history.ndjson does not yet exist but history.json does, we migrate the
// snapshot entries first so old messages are not lost.
func (pm *PersistenceManager) AppendHistoryEntry(chatID string, entry ChatHistory) error {
	chatDir := filepath.Join(pm.rootDir, "chats", chatID)
	ndjsonPath := filepath.Join(chatDir, "history.ndjson")
	jsonPath := filepath.Join(chatDir, "history.json")

	// Migrate existing history.json into history.ndjson on first append
	if _, err := os.Stat(ndjsonPath); os.IsNotExist(err) {
		if _, err := os.Stat(jsonPath); err == nil {
			var history []ChatHistory
			if loadErr := pm.loadJSON(jsonPath, &history); loadErr == nil && len(history) > 0 {
				f, openErr := os.OpenFile(ndjsonPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if openErr == nil {
					for _, h := range history {
						line, _ := json.Marshal(h)
						f.Write(append(line, '\n'))
					}
					f.Close()
				}
			}
		}
	}

	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal history entry: %w", err)
	}
	line = append(line, '\n')

	f, err := os.OpenFile(ndjsonPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open history.ndjson for append: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(line); err != nil {
		return fmt.Errorf("failed to append history entry: %w", err)
	}
	return nil
}

// loadJSON loads data from a JSON file
func (pm *PersistenceManager) loadJSON(filename string, data interface{}) error {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, data)
}

// RenameChat renames a chat
func (pm *PersistenceManager) RenameChat(chatID, newName string) error {
	chat, err := pm.LoadChat(chatID)
	if err != nil {
		return fmt.Errorf("failed to load chat: %w", err)
	}

	chat.Meta.Name = newName
	chat.Meta.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	return pm.SaveChat(chat)
}
