/**
 * Chat Service
 *
 * This package provides the main chat service for handling conversation management,
 * message processing, and AI interaction. Serves as the primary interface for chat operations.
 *
 * @package chat
 */

package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/backend/memory"
	"howl-chat/internal/backend/model"
	"howl-chat/internal/backend/rag"
	"howl-chat/internal/backend/types"
	"howl-chat/internal/library"
	"howl-chat/internal/paths"
	"howl-chat/internal/worldstore"
)

// Context status values report whether canonical world/scenario data could be
// used. They are explicit so a consumer never has to infer unavailability from
// nil data.
const (
	ContextNotSet                 = "not_set"
	ContextAvailable              = "available"
	ContextNotFound               = "not_found"
	ContextUnavailable            = "unavailable"
	ContextAuthorityIndeterminate = "authority_indeterminate"
	ContextRevisionMismatch       = "revision_mismatch"
)

// ContextStatus is the per-source result of the last canonical load.
type ContextStatus struct {
	World    string `json:"world"`
	Scenario string `json:"scenario"`
}

// Service handles chat operations including message processing and AI interaction
type Service struct {
	mu           sync.RWMutex
	modelManager *model.Manager
	llamaClient  llama.Client

	worldName      string
	scenarioName   string
	characterNames []string

	worldData     *WorldContext
	scenarioData  *ScenarioContext
	characterData map[string]*CharacterContext // Map character ID to character data

	// Canonical content access. Documents and world/scenario lore are read only
	// through worldstore; this service never writes item folders.
	worldStore    *worldstore.Store
	scenarioStore *worldstore.Store
	worldID       string
	scenarioID    string
	contextStatus ContextStatus

	// Lorebook and Memory Management
	manualLore        []lorebook.Entry
	automatedLore     []lorebook.Entry
	memorySnapshot    []lorebook.Entry // This will represent memories_snapshot.json
	triggeredMemories map[string]bool  // To track once_per_chat, once_per_scene, once_per_message

	// Memory Service Integration
	memoryService *memory.Service

	// RAG Service Integration
	ragService *rag.Service
	ragEnabled bool
}

// NewService creates a new chat service over the canonical content stores.
//
// It owns world/scenario/character/lore reads and the memory snapshot; the
// conversation itself — its turns, its prompt and its saved history — belongs to
// chat.Conversation, so there is exactly one in-memory history in the process.
func NewService(modelManager *model.Manager, client llama.Client) (*Service, error) {
	s := &Service{
		modelManager:      modelManager,
		llamaClient:       client,
		triggeredMemories: make(map[string]bool),

		// Initialize memory service
		memoryService: memory.NewService(""),
	}

	if err := s.LoadChatData("", "", "", "", []string{}, ""); err != nil {
		return nil, err
	}

	return s, nil
}

// SetContentLibraries attaches the content-library registry used to read
// canonical world and scenario documents. Until this is called, canonical data
// is reported unavailable rather than read from a name-keyed fallback.
func (s *Service) SetContentLibraries(reg *library.Registry) error {
	if reg == nil {
		s.worldStore = nil
		s.scenarioStore = nil
		return nil
	}
	worlds, err := worldstore.New(reg, library.Worlds)
	if err != nil {
		return err
	}
	scenarios, err := worldstore.New(reg, library.Scenarios)
	if err != nil {
		return err
	}
	s.worldStore = worlds
	s.scenarioStore = scenarios
	return nil
}

// GetContextStatus reports the result of the last canonical load for each
// source. A caller must display unavailability rather than substitute data.
func (s *Service) GetContextStatus() ContextStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.contextStatus
}

// LoadChatData loads the chat's canonical world/scenario documents and lore,
// plus character-owned and chat-owned data. World and scenario content is
// addressed by stable ID only: the name parameters are display labels and are
// never used to locate a folder. Unavailable or authority-indeterminate content
// is not loaded into the prompt pool, and no name-keyed fallback is attempted.
func (s *Service) LoadChatData(worldID, scenarioID, worldName, scenarioName string, characterNames []string, chatName string) error {
	projectRoot := paths.StateRoot()
	s.worldID = strings.TrimSpace(worldID)
	s.scenarioID = strings.TrimSpace(scenarioID)
	s.worldName = worldName
	s.scenarioName = scenarioName
	s.characterNames = characterNames

	worldData, worldManualLore, worldAutomatedLore, worldStatus := s.loadCanonicalWorld(s.worldID)
	scenarioData, scenarioManualLore, scenarioAutomatedLore, scenarioStatus := s.loadCanonicalScenario(s.scenarioID)
	s.worldData = worldData
	s.scenarioData = scenarioData
	s.contextStatus = ContextStatus{World: worldStatus, Scenario: scenarioStatus}

	s.characterData = make(map[string]*CharacterContext)
	for _, charName := range characterNames {
		charFilePath := filepath.Join(projectRoot, "Characters", charName, "character.json")
		charData, err := loadCharacterData(charFilePath)
		if err != nil {
			return err
		}
		if charData != nil {
			s.characterData[charName] = charData
		}
	}

	// Character lore stays character-owned until F3-D. World and scenario lore
	// comes only from the canonical item folder loaded above.
	var allManualLore []lorebook.Entry
	allManualLore = append(allManualLore, worldManualLore...)
	allManualLore = append(allManualLore, scenarioManualLore...)
	for _, charName := range characterNames {
		charManualLore, err := loadEntriesFromFile(filepath.Join(projectRoot, "Characters", charName, "character_lore.json"))
		if err != nil {
			return err
		}
		allManualLore = append(allManualLore, charManualLore...)
	}
	s.manualLore = allManualLore

	var allAutomatedLore []lorebook.Entry
	allAutomatedLore = append(allAutomatedLore, worldAutomatedLore...)
	allAutomatedLore = append(allAutomatedLore, scenarioAutomatedLore...)
	for _, charName := range characterNames {
		charAutomatedLore, err := loadEntriesFromFile(filepath.Join(projectRoot, "Characters", charName, "automated_character_lore.json"))
		if err != nil {
			return err
		}
		allAutomatedLore = append(allAutomatedLore, charAutomatedLore...)
	}
	s.automatedLore = allAutomatedLore

	// Load memory snapshot for this chat
	memoriesSnapshotPath := filepath.Join(projectRoot, "Chats", chatName, "memories_snapshot.json")
	existingMemories, err := loadMemoriesSnapshot(memoriesSnapshotPath)
	if err != nil {
		return err
	}

	// Load user-created memories from memory service. Character-scoped memory
	// ownership is matched by durable ID; a display name must not authorize it.
	characterIDs := make([]string, 0, len(s.characterData))
	for _, character := range s.characterData {
		if character != nil && character.ID != "" {
			characterIDs = append(characterIDs, character.ID)
		}
	}
	userMemories, err := loadMemoriesFromService(s.memoryService, s.worldID, s.scenarioID, characterIDs)
	if err != nil {
		// Log warning but don't fail - user memories are additive
		fmt.Printf("Warning: failed to load user memories: %v\n", err)
		userMemories = []lorebook.Entry{}
	}
	// Snapshots written before owner_kind existed carry character entries whose
	// OwnerID may be a display name. Keep them in the file, but do not inject
	// them into a character prompt without a stable ownership marker.
	existingMemories = filterStableCharacterMemoryEntries(existingMemories, characterIDs)

	if len(existingMemories) == 0 {
		// If no existing memories snapshot, create one from automated lore + user memories
		var chatMemories []lorebook.Entry
		for _, entry := range s.automatedLore {
			// Only copy automated lore owned by this chat's stable IDs.
			if (entry.Scope == lorebook.ScopeWorld && s.worldID != "" && entry.OwnerID == s.worldID) ||
				(entry.Scope == lorebook.ScopeScenario && s.scenarioID != "" && entry.OwnerID == s.scenarioID) ||
				(entry.Scope == lorebook.ScopeCharacter && entry.OwnerKind == types.OwnerKindCharacterID && contains(characterIDs, entry.OwnerID)) {
				chatMemories = append(chatMemories, entry)
			}
		}
		// Merge user memories, avoiding duplicates
		chatMemories = mergeMemories(chatMemories, userMemories)
		s.memorySnapshot = chatMemories
	} else {
		// Merge user memories into existing snapshot, avoiding duplicates
		s.memorySnapshot = mergeMemories(existingMemories, userMemories)
	}

	return nil
}

// loadCanonicalWorld reads one world document and its lore through the store.
// Storage uncertainty is returned as a status, never as a reason to read a
// legacy name folder. A corrupt or indeterminate document contributes nothing.
func (s *Service) loadCanonicalWorld(id string) (*WorldContext, []lorebook.Entry, []lorebook.Entry, string) {
	if id == "" {
		return nil, nil, nil, ContextNotSet
	}
	if s.worldStore == nil {
		return nil, nil, nil, ContextUnavailable
	}
	item, err := s.worldStore.Load(id)
	if err != nil {
		if errors.Is(err, worldstore.ErrNotFound) {
			return nil, nil, nil, ContextNotFound
		}
		return nil, nil, nil, ContextUnavailable
	}
	if item.AuthorityIndeterminate {
		return nil, nil, nil, ContextAuthorityIndeterminate
	}
	payload, err := json.Marshal(item.Document)
	if err != nil {
		return nil, nil, nil, ContextUnavailable
	}
	world, err := decodeWorldContext(payload)
	if err != nil || world == nil {
		return nil, nil, nil, ContextUnavailable
	}
	manual, err := s.loadItemLore(s.worldStore, id, "world_lore.json", lorebook.ScopeWorld, false)
	if err != nil {
		return nil, nil, nil, ContextUnavailable
	}
	automated, err := s.loadItemLore(s.worldStore, id, "automated_world_lore.json", lorebook.ScopeWorld, true)
	if err != nil {
		return nil, nil, nil, ContextUnavailable
	}
	status := ContextAvailable
	if item.RevisionMismatch {
		status = ContextRevisionMismatch
	}
	return world, manual, automated, status
}

// loadCanonicalScenario is the scenario counterpart of loadCanonicalWorld.
func (s *Service) loadCanonicalScenario(id string) (*ScenarioContext, []lorebook.Entry, []lorebook.Entry, string) {
	if id == "" {
		return nil, nil, nil, ContextNotSet
	}
	if s.scenarioStore == nil {
		return nil, nil, nil, ContextUnavailable
	}
	item, err := s.scenarioStore.Load(id)
	if err != nil {
		if errors.Is(err, worldstore.ErrNotFound) {
			return nil, nil, nil, ContextNotFound
		}
		return nil, nil, nil, ContextUnavailable
	}
	if item.AuthorityIndeterminate {
		return nil, nil, nil, ContextAuthorityIndeterminate
	}
	payload, err := json.Marshal(item.Document)
	if err != nil {
		return nil, nil, nil, ContextUnavailable
	}
	scenario, err := decodeScenarioContext(payload)
	if err != nil || scenario == nil {
		return nil, nil, nil, ContextUnavailable
	}
	manual, err := s.loadItemLore(s.scenarioStore, id, "scenario_lore.json", lorebook.ScopeScenario, false)
	if err != nil {
		return nil, nil, nil, ContextUnavailable
	}
	automated, err := s.loadItemLore(s.scenarioStore, id, "automated_scenario_lore.json", lorebook.ScopeScenario, true)
	if err != nil {
		return nil, nil, nil, ContextUnavailable
	}
	status := ContextAvailable
	if item.RevisionMismatch {
		status = ContextRevisionMismatch
	}
	return scenario, manual, automated, status
}

// loadItemLore reads one canonical lore file. Entries owned by another stable
// ID are excluded rather than injected; a missing owner is stamped from the
// canonical folder that holds it. Legacy generated entries keep their
// unattributable origin and stay prompt-eligible.
func (s *Service) loadItemLore(store *worldstore.Store, id, filename string, scope lorebook.Scope, legacyGenerated bool) ([]lorebook.Entry, error) {
	file, err := store.LoadLore(id, filename)
	if err != nil {
		return nil, err
	}
	if file.AuthorityIndeterminate {
		return nil, worldstore.ErrUnavailable
	}
	entries := make([]lorebook.Entry, 0, len(file.Entries))
	for _, raw := range file.Entries {
		var entry lorebook.Entry
		if err := json.Unmarshal(raw, &entry); err != nil {
			return nil, fmt.Errorf("decode %s: %w", filename, err)
		}
		entry.Enabled = true
		if entry.Scope == "" {
			entry.Scope = scope
			entry.Type = scope
		}
		if legacyGenerated {
			entry.Automated = true
			entry.OriginKind = lorebook.OriginLegacyGeneratedUnattributed
		}
		if entry.Scope == lorebook.ScopeWorld || entry.Scope == lorebook.ScopeScenario {
			if entry.OwnerID == "" {
				entry.OwnerID = id
			} else if entry.OwnerID != id {
				// Unmapped legacy owner: preserve the entry in the file, exclude it
				// from prompts until ownership is repaired explicitly.
				continue
			}
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// filterStableCharacterMemoryEntries removes character-scope snapshot entries
// that lack a stable character_id owner kind. Other scopes are untouched.
func filterStableCharacterMemoryEntries(entries []lorebook.Entry, characterIDs []string) []lorebook.Entry {
	filtered := make([]lorebook.Entry, 0, len(entries))
	for _, entry := range entries {
		if entry.Scope == lorebook.ScopeCharacter {
			if entry.OwnerKind != types.OwnerKindCharacterID || !contains(characterIDs, entry.OwnerID) {
				continue
			}
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

// mergeMemories combines existing and user memories, avoiding duplicates by ID
func mergeMemories(existing, user []lorebook.Entry) []lorebook.Entry {
	seen := make(map[string]bool)
	for _, e := range existing {
		seen[e.ID] = true
	}

	result := make([]lorebook.Entry, len(existing))
	copy(result, existing)

	for _, u := range user {
		if !seen[u.ID] {
			result = append(result, u)
			seen[u.ID] = true
		}
	}

	return result
}

// isRetryableError checks if an error is transient and should be retried
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Check for common retryable error patterns
	retryablePatterns := []string{
		"resource_exhausted",
		"rate limit",
		"429",
		"timeout",
		"connection refused",
		"temporary",
		"unavailable",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(strings.ToLower(errStr), pattern) {
			return true
		}
	}

	return false
}

func (s *Service) GetWorldData() *WorldContext {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.worldData
}

func (s *Service) GetScenarioData() *ScenarioContext {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.scenarioData
}

func (s *Service) GetCharacterData() []*CharacterContext {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Convert map to slice
	chars := make([]*CharacterContext, 0, len(s.characterData))
	for _, char := range s.characterData {
		chars = append(chars, char)
	}
	return chars
}

func (s *Service) GetAllLorebookEntries() []lorebook.Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Combine manual and automated lore
	allLore := make([]lorebook.Entry, 0, len(s.manualLore)+len(s.automatedLore))
	allLore = append(allLore, s.manualLore...)
	allLore = append(allLore, s.automatedLore...)
	return allLore
}

// SaveChatData saves all relevant chat data to JSON files.
func (s *Service) SaveChatData(chatName string, syncAutomatedLore bool) error {
	projectRoot := paths.StateRoot()

	// Save chat.json
	chatFilePath := filepath.Join(projectRoot, "Chats", chatName, "chat.json")
	chatData := map[string]interface{}{
		"chat_name":  chatName,
		"seed":       s.currentSeed(), // -1 means random, otherwise the pinned seed
		"world":      s.worldName,
		"scenario":   s.scenarioName,
		"characters": s.characterNames,
		"version":    1,
	}
	chatJSON, err := json.MarshalIndent(chatData, "", "  ")
	if err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to marshal chat.json", err)
	}
	if err := writeFileAtomic(chatFilePath, chatJSON, 0644); err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to write chat.json", err)
	}

	// The conversation's history is not written here. This service has no
	// in-memory history any more; chat.Conversation owns the turns, the append log
	// and the snapshot, and writing a second history.json from an empty context
	// would recreate exactly the competing representation this slice retires.

	// Save memories_snapshot.json with memory service integration
	if err := s.memoryService.SaveChatData(chatName, syncAutomatedLore); err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to save memories", err)
	}

	// Sync is character-only until F3-D. World/scenario generated lore is
	// chat-owned now; writing it into an item folder would rebuild the second
	// writer this slice removes.
	if syncAutomatedLore {
		for _, mem := range s.memorySnapshot {
			if filePath, ok := s.automatedLoreSyncTarget(projectRoot, mem); ok {
				updateAutomatedLoreFile(filePath, mem)
			}
		}
	}

	return nil
}

// automatedLoreSyncTarget is the single allowlist for generated-lore writes
// back to content folders. World and scenario scopes have no target: those
// writes are removed, not redirected.
func (s *Service) automatedLoreSyncTarget(projectRoot string, entry lorebook.Entry) (string, bool) {
	if !entry.Automated || entry.Scope != lorebook.ScopeCharacter || strings.TrimSpace(entry.OwnerID) == "" {
		return "", false
	}
	return filepath.Join(projectRoot, "Characters", entry.OwnerID, "automated_character_lore.json"), true
}

// Helper to update automated lore files
func updateAutomatedLoreFile(filePath string, newEntry lorebook.Entry) error {
	// Read existing entries
	existingEntries, err := loadEntriesFromFile(filePath)
	if err != nil {
		return err
	}

	found := false
	for i, entry := range existingEntries {
		if entry.ID == newEntry.ID {
			existingEntries[i] = newEntry // Update existing entry
			found = true
			break
		}
	}

	if !found {
		existingEntries = append(existingEntries, newEntry) // Add new entry
	}

	// Write back to file
	lb := struct {
		Entries []lorebook.Entry `json:"entries"`
	}{
		Entries: existingEntries,
	}
	data, err := json.MarshalIndent(lb, "", "  ")
	if err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to marshal automated lore JSON", err)
	}
	if err := writeFileAtomic(filePath, data, 0644); err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to write automated lore file", err)
	}
	return nil
}

// SetRAGService attaches the RAG service to the chat service
func (s *Service) SetRAGService(ragSvc *rag.Service) {
	s.ragService = ragSvc
}

// SetRAGEnabled toggles RAG retrieval on/off
func (s *Service) SetRAGEnabled(enabled bool) {
	s.ragEnabled = enabled
}

// IsRAGEnabled returns whether RAG is enabled
func (s *Service) IsRAGEnabled() bool {
	return s.ragEnabled
}

// Close releases all resources
func (s *Service) Close() error {
	return nil
}
