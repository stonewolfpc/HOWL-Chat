package chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"howl-chat/internal/durable"
	"howl-chat/internal/paths"
)

// PersistenceManager handles chat persistence with the specified directory structure
type PersistenceManager struct {
	rootDir string
	config  AppConfig
	// membership is the single authority for collection membership and deletion
	// intent. It is loaded from the lifecycle log at initialize.
	membership *Membership
}

// RootDir reports the application-owned chat root. Callers must use this
// instead of guessing an executable-relative "chats" directory.
func (pm *PersistenceManager) RootDir() string {
	if pm == nil {
		return ""
	}
	return pm.rootDir
}

// ChatDirectory reports the on-disk directory for a conversation. It uses the same
// constructor every mutation uses, so the path shown to a user is the path the
// application would have written to.
func (pm *PersistenceManager) ChatDirectory(chatID string) (string, error) {
	return pm.conversationDir(chatID)
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
	// An isolated native test run keeps this legacy store inside its own root.
	if paths.IsolationActive() {
		return paths.StatePath("persisted-chats"), nil
	}

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
	if err := ensureDir(pm.rootDir); err != nil {
		return fmt.Errorf("failed to create root directory: %w", err)
	}

	// Create chats directory
	chatsDir := filepath.Join(pm.rootDir, "chats")
	if err := ensureDir(chatsDir); err != nil {
		return fmt.Errorf("failed to create chats directory: %w", err)
	}

	// Locks live outside chats/ so a lock file can never be a candidate
	// conversation, and so deletion never removes a directory holding an open
	// handle.
	if err := ensureDir(pm.locksDir()); err != nil {
		return fmt.Errorf("failed to create lock directory: %w", err)
	}

	// Load the membership authority. This is also where a legacy deleted.ndjson
	// is absorbed into the lifecycle log, once.
	membership, err := LoadMembership(pm.rootDir)
	if err != nil {
		return fmt.Errorf("failed to load conversation membership: %w", err)
	}
	pm.membership = membership

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

	// LastChatID has one owner (this config) and one serialised write. Two
	// instances pointing at the same root would otherwise interleave a
	// read-modify-write: the lock makes the last intent win whole, rather than
	// leaving a torn file or a silently lost update.
	unlock, err := durable.LockFile(filepath.Join(pm.rootDir, "config.lock"), 2*time.Second)
	if err != nil {
		return fmt.Errorf("failed to lock config: %w", err)
	}
	defer unlock()

	data, err := json.MarshalIndent(pm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := writeFileAtomic(configPath, data, 0644); err != nil {
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

// newChatData builds a conversation's initial value. It performs no I/O.
func newChatData(chatID, name string) *ChatData {
	now := time.Now().UTC()
	return &ChatData{
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
}

// errConversationExists means the generated identifier was already claimed. The
// caller retries with a new identifier instead of adopting the existing directory.
var errConversationExists = errors.New("conversation identifier already claimed")

// NewChat creates a new chat with the given name.
//
// Creation claims its directory exclusively. Two creations in the same second no
// longer share one directory: the second claims a suffixed identifier, so neither
// can overwrite the other's metadata, companions or history (the data loss row 580
// observed).
func (pm *PersistenceManager) NewChat(name string) (*ChatData, error) {
	baseID := pm.generateChatID()
	for attempt := 0; attempt < 10; attempt++ {
		chatID := baseID
		if attempt > 0 {
			chatID = fmt.Sprintf("%s-%d", baseID, attempt+1)
		}

		var created *ChatData
		err := pm.withChatLock("create", chatID, func() error {
			chatDir, err := pm.conversationDir(chatID)
			if err != nil {
				return err
			}
			if err := claimDir(chatDir); err != nil {
				if errors.Is(err, fs.ErrExist) {
					return errConversationExists
				}
				return fmt.Errorf("failed to claim chat directory: %w", err)
			}
			// The directory exists; now it becomes a conversation. A crash between
			// the claim and this record leaves an orphaned directory, which
			// reconciliation reports instead of quietly adopting.
			if err := pm.membership.Create(chatID, name); err != nil {
				return err
			}
			chatData := newChatData(chatID, name)
			if err := pm.saveChatFilesLocked(chatData, true); err != nil {
				return fmt.Errorf("failed to save new chat: %w", err)
			}
			created = chatData
			return nil
		})
		if errors.Is(err, errConversationExists) {
			continue
		}
		if err != nil {
			return nil, err
		}
		return created, nil
	}
	return nil, fmt.Errorf("failed to create chat: no unused identifier after 10 attempts")
}

// SaveChat saves all chat data to files.
func (pm *PersistenceManager) SaveChat(chatData *ChatData) error {
	if chatData == nil {
		return &Refusal{Code: RefusedUnbound, Op: "save", Detail: "no chat data"}
	}
	return pm.withChatLock("save", chatData.Meta.ID, func() error {
		return pm.saveChatLocked(chatData)
	})
}

// saveChatLocked is SaveChat inside the held lock. A conversation the user
// deleted is refused rather than recreated, so a late autosave cannot resurrect
// it.
func (pm *PersistenceManager) saveChatLocked(chatData *ChatData) error {
	return pm.saveChatFilesLocked(chatData, true)
}

// saveChatFilesLocked writes metadata and companions. The history snapshot is
// written only when writeSnapshot is set; the final commit writes it last, after
// the append log, so an interruption cannot leave a snapshot newer than the log.
func (pm *PersistenceManager) saveChatFilesLocked(chatData *ChatData, writeSnapshot bool) error {
	chatID := chatData.Meta.ID
	if pm.deletedSet()[chatID] {
		return &Refusal{Code: RefusedUnbound, Op: "save", ChatID: chatID, Detail: "conversation was deleted"}
	}
	chatDir, err := pm.conversationDir(chatID)
	if err != nil {
		return err
	}

	// A save never recreates a conversation the collection already knows: if its
	// directory is gone, the save is refused rather than materialising the chat
	// again around whatever is in memory. An identifier the collection has never
	// recorded is a first save (import, fixture), so it is recorded as a member
	// before the write — it can never end up as an unclaimed directory.
	if statErr := requireExistingDir(chatDir); statErr != nil {
		if _, known := pm.membership.Member(chatID); known {
			return &Refusal{Code: RefusedUnbound, Op: "save", ChatID: chatID, Detail: "conversation directory is missing", Err: statErr}
		}
		if err := pm.membership.Create(chatID, chatData.Meta.Name); err != nil {
			return err
		}
		if err := ensureDir(chatDir); err != nil {
			return fmt.Errorf("failed to create chat directory: %w", err)
		}
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
	if writeSnapshot {
		if err := pm.saveJSON(filepath.Join(chatDir, "history.json"), chatData.History); err != nil {
			return fmt.Errorf("failed to save history.json: %w", err)
		}
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
	if pm.deletedSet()[chatID] {
		return nil, &Refusal{Code: RefusedUnbound, Op: "load", ChatID: chatID, Detail: "conversation was deleted"}
	}
	chatDir, err := pm.conversationDir(chatID)
	if err != nil {
		return nil, err
	}

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

	// World / Scenario IDs stay unresolved when the record carries only a
	// display name. Copying the name into an ID field invents identity and lets a
	// legacy chat look canonical; the consumer must report it unset instead.
}

// SaveChatAs creates a copy of the current chat with a new name.
func (pm *PersistenceManager) SaveChatAs(currentChatID, newName string) (*ChatData, error) {
	// Load current chat
	currentChat, err := pm.LoadChat(currentChatID)
	if err != nil {
		return nil, fmt.Errorf("failed to load current chat: %w", err)
	}

	baseID := pm.generateChatID()
	for attempt := 0; attempt < 10; attempt++ {
		chatID := baseID
		if attempt > 0 {
			chatID = fmt.Sprintf("%s-%d", baseID, attempt+1)
		}

		var created *ChatData
		err = pm.withChatLock("save-as", chatID, func() error {
			chatDir, err := pm.conversationDir(chatID)
			if err != nil {
				return err
			}
			if err := claimDir(chatDir); err != nil {
				if errors.Is(err, fs.ErrExist) {
					return errConversationExists
				}
				return fmt.Errorf("failed to claim chat directory: %w", err)
			}
			if err := pm.membership.Create(chatID, newName); err != nil {
				return err
			}

			newChat := newChatData(chatID, newName)

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

			// One write of the finished copy: the destination is never published in
			// the empty intermediate state the previous implementation wrote first.
			if err := pm.saveChatFilesLocked(newChat, true); err != nil {
				return fmt.Errorf("failed to save new chat: %w", err)
			}
			created = newChat
			return nil
		})
		if errors.Is(err, errConversationExists) {
			continue
		}
		if err != nil {
			return nil, err
		}
		return created, nil
	}
	return nil, fmt.Errorf("failed to copy chat: no unused identifier after 10 attempts")
}

// DeleteChat removes a conversation from the visible collection and then makes a
// best-effort attempt to remove its directory.
//
// Deletion is deliberately not a verified operation. Intent is recorded first, so
// a directory that survives a blocked delete can never reappear as a conversation
// and a late writer cannot recreate it. No proof of physical destruction is
// awaited, and an incomplete removal does not restore the item or fail the call.
func (pm *PersistenceManager) DeleteChat(chatID string) error {
	chatDir, err := pm.conversationDir(chatID)
	if err != nil {
		return err
	}

	if err := pm.membership.Delete(chatID, chatDir); err != nil {
		return fmt.Errorf("failed to record deletion intent: %w", err)
	}

	if pm.config.LastChatID == chatID {
		pm.config.LastChatID = ""
		_ = pm.saveConfig()
	}

	// The Recycle Bin is the target, and there is no permanent-delete fallback in
	// this code: if the shell refuses, the files stay and are reported.
	err = pm.withChatLock("delete", chatID, func() error {
		if err := recycleTree(chatDir); err != nil && !os.IsNotExist(err) {
			return &DeletionIncomplete{ChatID: chatID, Directory: chatDir, Err: err}
		}
		return nil
	})
	if err == nil {
		return nil
	}

	var incomplete *DeletionIncomplete
	if errors.As(err, &incomplete) || IsRefusal(err, RefusedBusy) {
		fmt.Printf("[DeleteChat] %s left the collection; its directory may remain at %s: %v\n", chatID, chatDir, err)
		return nil
	}
	return err
}

// Present reports whether the collection claims this conversation and its
// directory exists. It is the check every write-into-a-conversation path uses:
// only a present member may be written to.
func (pm *PersistenceManager) Present(chatID string) bool {
	if pm == nil || pm.membership == nil || chatID == "" {
		return false
	}
	if pm.membership.Deleted(chatID) {
		return false
	}
	if _, ok := pm.membership.Member(chatID); !ok {
		return false
	}
	dir, err := pm.conversationDir(chatID)
	if err != nil {
		return false
	}
	return requireExistingDir(dir) == nil
}

// GetChatList returns the collection as the membership authority reconciles it.
//
// Read-only by construction: present conversations carry their metadata, a
// claimed conversation whose directory is gone is reported as missing, and a
// directory nobody claimed is reported as orphaned. Nothing is adopted, repaired
// or removed here.
func (pm *PersistenceManager) GetChatList() ([]ChatSummary, error) {
	chatsDir, err := pm.chatsRoot()
	if err != nil {
		return nil, err
	}
	conversations, err := pm.membership.Reconcile(chatsDir)
	if err != nil {
		return nil, err
	}

	summaries := make([]ChatSummary, 0, len(conversations))
	for _, conversation := range conversations {
		summary := ChatSummary{
			ID:        conversation.ID,
			Name:      conversation.Name,
			State:     string(conversation.State),
			Directory: conversation.Directory,
		}
		if summary.Name == "" {
			summary.Name = conversation.ID
		}
		if conversation.State != StatePresent {
			summaries = append(summaries, summary)
			continue
		}

		var meta ChatMeta
		if err := pm.loadJSON(filepath.Join(conversation.Directory, "meta.json"), &meta); err != nil {
			// A claimed directory whose metadata cannot be read is still a member;
			// report it as present with what the authority knows rather than
			// hiding the conversation from its owner.
			summaries = append(summaries, summary)
			continue
		}
		if meta.Name != "" {
			summary.Name = meta.Name
		}

		createdAt, _ := time.Parse(time.RFC3339, meta.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339, meta.UpdatedAt)
		summary.CreatedAt = createdAt.Format(time.RFC3339)
		summary.UpdatedAt = updatedAt.Format(time.RFC3339)
		summary.WorldName = meta.WorldName
		summary.ScenarioName = meta.ScenarioName
		summary.CharacterNames = meta.CharacterNames
		summary.LorebookNames = meta.LorebookNames

		var history []ChatHistory
		pm.loadJSON(filepath.Join(conversation.Directory, "history.json"), &history)
		summary.MessageCount = len(history)

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
	chatDir, err := pm.conversationDir(pm.config.LastChatID)
	if err != nil {
		return nil, err
	}
	metaPath := filepath.Join(chatDir, "meta.json")
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
	return writeFileAtomic(filename, jsonData, 0644)
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
	return pm.withChatLock("rewrite-history", chatID, func() error {
		return pm.rewriteHistoryLocked(chatID, entries)
	})
}

func (pm *PersistenceManager) rewriteHistoryLocked(chatID string, entries []ChatHistory) error {
	if pm.deletedSet()[chatID] {
		return &Refusal{Code: RefusedUnbound, Op: "rewrite-history", ChatID: chatID, Detail: "conversation was deleted"}
	}
	chatDir, err := pm.conversationDir(chatID)
	if err != nil {
		return err
	}
	// A rewrite is never a creation: the conversation exists because the
	// lifecycle created it, not because something wrote about its identifier.
	if err := requireExistingDir(chatDir); err != nil {
		return &Refusal{Code: RefusedUnbound, Op: "rewrite-history", ChatID: chatID, Detail: "conversation directory is missing", Err: err}
	}

	if err := pm.writeHistoryLogLocked(chatID, entries); err != nil {
		return err
	}

	// Keep the JSON snapshot identical so a fallback load sees the same history.
	return pm.saveJSON(filepath.Join(chatDir, "history.json"), entries)
}

// writeHistoryLogLocked replaces history.ndjson with the supplied entries.
func (pm *PersistenceManager) writeHistoryLogLocked(chatID string, entries []ChatHistory) error {
	chatDir, err := pm.conversationDir(chatID)
	if err != nil {
		return err
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

	if err := writeFileAtomic(filepath.Join(chatDir, "history.ndjson"), buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write history.ndjson: %w", err)
	}
	return nil
}

// CommitChat persists metadata, companions and history as one operation, writing
// the append log before the JSON snapshot.
//
// Ordering is the point. The loader prefers a non-empty log, so an interruption
// between the two writes can only leave the log newer than the snapshot - the
// state that loses nothing. The previous shutdown path wrote the snapshot first and
// then rewrote both, which is why every clean exit passed through "snapshot newer
// than log" (the window recorded in row 570).
func (pm *PersistenceManager) CommitChat(chatData *ChatData) error {
	if chatData == nil {
		return &Refusal{Code: RefusedUnbound, Op: "commit", Detail: "no chat data"}
	}
	// A commit writes into an existing conversation. An identifier the authority
	// does not claim is a cached session or a stale client — never a reason to
	// create one.
	if !pm.Present(chatData.Meta.ID) {
		return &Refusal{Code: RefusedUnbound, Op: "commit", ChatID: chatData.Meta.ID, Detail: "conversation is not present in the collection"}
	}
	return pm.withChatLock("commit", chatData.Meta.ID, func() error {
		if err := pm.saveChatFilesLocked(chatData, false); err != nil {
			return err
		}
		if err := pm.writeHistoryLogLocked(chatData.Meta.ID, chatData.History); err != nil {
			return err
		}
		chatDir, err := pm.conversationDir(chatData.Meta.ID)
		if err != nil {
			return err
		}
		return pm.saveJSON(filepath.Join(chatDir, "history.json"), chatData.History)
	})
}

// AppendHistoryEntry appends a single entry to history.ndjson using an O(1) file append.
// This is the hot path called after every message — no full file rewrite.
// If history.ndjson does not yet exist but history.json does, we migrate the
// snapshot entries first so old messages are not lost.
func (pm *PersistenceManager) AppendHistoryEntry(chatID string, entry ChatHistory) error {
	return pm.withChatLock("append-history", chatID, func() error {
		return pm.appendHistoryLocked(chatID, entry)
	})
}

func (pm *PersistenceManager) appendHistoryLocked(chatID string, entry ChatHistory) error {
	if pm.deletedSet()[chatID] {
		return &Refusal{Code: RefusedUnbound, Op: "append-history", ChatID: chatID, Detail: "conversation was deleted"}
	}
	chatDir, err := pm.conversationDir(chatID)
	if err != nil {
		return err
	}
	// An append is never a creation: a missing directory is reported, never
	// materialised by the write that mentions the identifier.
	if err := requireExistingDir(chatDir); err != nil {
		return &Refusal{Code: RefusedUnbound, Op: "append-history", ChatID: chatID, Detail: "conversation directory is missing", Err: err}
	}
	ndjsonPath := filepath.Join(chatDir, "history.ndjson")
	jsonPath := filepath.Join(chatDir, "history.json")

	// Migrate existing history.json into history.ndjson on first append
	if _, err := os.Stat(ndjsonPath); os.IsNotExist(err) {
		if _, err := os.Stat(jsonPath); err == nil {
			var history []ChatHistory
			if loadErr := pm.loadJSON(jsonPath, &history); loadErr == nil && len(history) > 0 {
				var buf bytes.Buffer
				for _, h := range history {
					line, marshalErr := json.Marshal(h)
					if marshalErr != nil {
						return fmt.Errorf("failed to marshal migrated history entry: %w", marshalErr)
					}
					buf.Write(line)
					buf.WriteByte(0x0A)
				}
				if err := writeFileAtomic(ndjsonPath, buf.Bytes(), 0644); err != nil {
					return fmt.Errorf("failed to migrate history.json: %w", err)
				}
			}
		}
	}

	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal history entry: %w", err)
	}
	line = append(line, '\n')

	if err := appendLine(ndjsonPath, line); err != nil {
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
	return pm.withChatLock("rename", chatID, func() error {
		chat, err := pm.LoadChat(chatID)
		if err != nil {
			return fmt.Errorf("failed to load chat: %w", err)
		}

		chat.Meta.Name = newName
		chat.Meta.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

		return pm.saveChatLocked(chat)
	})
}
