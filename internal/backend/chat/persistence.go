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

	// Save canonical history.json snapshot (atomic write)
	if err := pm.saveJSON(filepath.Join(chatDir, "history.json"), chatData.History); err != nil {
		return fmt.Errorf("failed to save history.json: %w", err)
	}
	// Remove the append log — the canonical snapshot is now authoritative.
	// On the next turn, AppendHistoryEntry recreates history.ndjson from scratch.
	// This keeps the log from diverging from the full snapshot.
	_ = os.Remove(filepath.Join(chatDir, "history.ndjson"))

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

	return chatData, nil
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

// AppendHistoryEntry appends a single entry to history.ndjson using an O(1) file append.
// This is the hot path called after every message — no full file rewrite.
func (pm *PersistenceManager) AppendHistoryEntry(chatID string, entry ChatHistory) error {
	chatDir := filepath.Join(pm.rootDir, "chats", chatID)
	ndjsonPath := filepath.Join(chatDir, "history.ndjson")

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
