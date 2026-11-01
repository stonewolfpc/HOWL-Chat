/**
 * Memory Service
 *
 * This package provides comprehensive memory management functionality
 * for storing, retrieving, and managing automated lorebook entries
 * captured during chat sessions. Integrates with existing lorebook system.
 *
 * @package memory
 * @version 1.0.0
 */

package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/backend/types"
	"howl-chat/internal/paths"
)

// Service handles memory operations and persistence
type Service struct {
	projectRoot string
	chatID      string
}

// SetChatID switches memory storage scope to a specific chat session.
func (s *Service) SetChatID(chatID string) {
	if strings.TrimSpace(chatID) == "" {
		return
	}
	s.chatID = chatID
}

// NewService creates a new memory service instance
func NewService(chatID string) *Service {
	return &Service{
		projectRoot: paths.ProjectRoot(),
		chatID:      chatID,
	}
}

// MemoryStorage handles file-based memory persistence
type MemoryStorage struct {
	projectRoot string
	chatID      string
}

// NewMemoryStorage creates a new storage instance
func NewMemoryStorage(chatID string) *MemoryStorage {
	return &MemoryStorage{
		projectRoot: paths.ProjectRoot(),
		chatID:      chatID,
	}
}

// GetMemories retrieves memories with optional filtering
func (s *Service) GetMemories(filter types.MemoryFilter) ([]types.MemoryRecord, error) {
	storage := NewMemoryStorage(s.chatID)

	// Load memories from file
	memories, err := storage.loadMemories()
	if err != nil {
		return nil, fmt.Errorf("failed to load memories: %w", err)
	}

	// Apply filters
	filtered := storage.applyFilters(memories, filter)

	return filtered, nil
}

// GetMemoryDetail retrieves a specific memory by ID
func (s *Service) GetMemoryDetail(memoryID string) (*types.MemoryRecord, error) {
	storage := NewMemoryStorage(s.chatID)

	memories, err := storage.loadMemories()
	if err != nil {
		return nil, fmt.Errorf("failed to load memories: %w", err)
	}

	// Find the memory
	for _, memory := range memories {
		if memory.ID == memoryID {
			return &memory, nil
		}
	}

	return nil, fmt.Errorf("memory not found: %s", memoryID)
}

// AddMemory creates a new memory record and appends it to storage.
func (s *Service) AddMemory(record types.MemoryRecord) (*types.MemoryRecord, error) {
	storage := NewMemoryStorage(s.chatID)

	memories, err := storage.loadMemories()
	if err != nil {
		return nil, fmt.Errorf("failed to load memories: %w", err)
	}

	// Ensure ID is set
	if record.ID == "" {
		record.ID = types.NewMemoryID()
	}
	now := time.Now().Format(time.RFC3339)
	if record.CreatedAt == "" {
		record.CreatedAt = now
	}
	if record.UpdatedAt == "" {
		record.UpdatedAt = now
	}
	if record.Status == "" {
		record.Status = "active"
	}

	memories = append(memories, record)
	if err := storage.saveMemories(memories); err != nil {
		return nil, fmt.Errorf("failed to save memories: %w", err)
	}

	return &record, nil
}

// EditMemory updates an existing memory
func (s *Service) EditMemory(memoryID string, edits types.MemoryEdit) (*types.MemoryRecord, error) {
	storage := NewMemoryStorage(s.chatID)

	memories, err := storage.loadMemories()
	if err != nil {
		return nil, fmt.Errorf("failed to load memories: %w", err)
	}

	// Find and update the memory
	for i, memory := range memories {
		if memory.ID == memoryID {
			// Apply edits
			if edits.Title != "" {
				memory.Title = edits.Title
			}
			if edits.Summary != "" {
				memory.Summary = edits.Summary
			}
			if len(edits.Tags) > 0 {
				memory.Tags = edits.Tags
			}
			memory.Pinned = edits.Pinned
			memory.UpdatedAt = time.Now().Format(time.RFC3339)

			// Save updated memory
			memories[i] = memory
			err = storage.saveMemories(memories)
			if err != nil {
				return nil, fmt.Errorf("failed to save memories: %w", err)
			}
			return &memory, nil
		}
	}

	return nil, fmt.Errorf("memory not found: %s", memoryID)
}

// PinMemory toggles the pinned status of a memory
func (s *Service) PinMemory(memoryID string, pinned bool) error {
	storage := NewMemoryStorage(s.chatID)

	memories, err := storage.loadMemories()
	if err != nil {
		return fmt.Errorf("failed to load memories: %w", err)
	}

	// Find and update the memory
	for i, memory := range memories {
		if memory.ID == memoryID {
			memory.Pinned = pinned
			memory.UpdatedAt = time.Now().Format(time.RFC3339)

			// Save updated memory
			memories[i] = memory
			return storage.saveMemories(memories)
		}
	}

	return fmt.Errorf("memory not found: %s", memoryID)
}

// DeleteMemory removes a memory from storage
func (s *Service) DeleteMemory(memoryID string) error {
	storage := NewMemoryStorage(s.chatID)

	memories, err := storage.loadMemories()
	if err != nil {
		return fmt.Errorf("failed to load memories: %w", err)
	}

	// Find and remove the memory
	var updatedMemories []types.MemoryRecord
	found := false

	for _, memory := range memories {
		if memory.ID != memoryID {
			updatedMemories = append(updatedMemories, memory)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("memory not found: %s", memoryID)
	}

	// Save updated memories
	err = storage.saveMemories(updatedMemories)
	if err != nil {
		return fmt.Errorf("failed to save memories: %w", err)
	}

	return nil
}

// ExportMemories exports memories with optional filtering
func (s *Service) ExportMemories(filter types.MemoryFilter) ([]types.MemoryRecord, error) {
	storage := NewMemoryStorage(s.chatID)

	// Load memories
	memories, err := storage.loadMemories()
	if err != nil {
		return nil, fmt.Errorf("failed to load memories: %w", err)
	}

	// Apply filters
	filtered := storage.applyFilters(memories, filter)

	return filtered, nil
}

// defaultAutoCaptureConfig is used for chats that have never saved settings.
func defaultAutoCaptureConfig() types.AutoCaptureConfig {
	return types.AutoCaptureConfig{
		Enabled:       true,
		MinConfidence: 75,
		MaxPerMessage: 3,
		TriggerWords:  []string{"remember", "note", "important", "key"},
	}
}

// parseAutoCaptureConfig decodes saved settings, falling back to defaults for
// fields an older config file may not contain.
func parseAutoCaptureConfig(raw []byte) (types.AutoCaptureConfig, error) {
	config := defaultAutoCaptureConfig()
	if err := json.Unmarshal(raw, &config); err != nil {
		return config, err
	}
	if len(config.TriggerWords) == 0 {
		config.TriggerWords = defaultAutoCaptureConfig().TriggerWords
	}
	return config, nil
}

// GetAutoCaptureConfig retrieves automatic memory capture settings.
// Settings live in Chats/<chatID>/auto_capture_config.json; a missing file means
// the chat has never been configured, so the defaults apply.
func (s *Service) GetAutoCaptureConfig() (*types.AutoCaptureConfig, error) {
	storage := NewMemoryStorage(s.chatID)
	configPath := filepath.Join(storage.projectRoot, "Chats", s.chatID, "auto_capture_config.json")

	raw, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			config := defaultAutoCaptureConfig()
			return &config, nil
		}
		return nil, fmt.Errorf("failed to read auto capture config: %w", err)
	}

	config, err := parseAutoCaptureConfig(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse auto capture config: %w", err)
	}
	return &config, nil
}

// SetAutoCaptureConfig updates automatic memory capture settings
func (s *Service) SetAutoCaptureConfig(config types.AutoCaptureConfig) error {
	storage := NewMemoryStorage(s.chatID)

	// Ensure chat directory exists
	chatDir := filepath.Join(storage.projectRoot, "Chats", s.chatID)
	if err := os.MkdirAll(chatDir, 0755); err != nil {
		return fmt.Errorf("failed to create chat directory: %w", err)
	}

	// Save config to file
	configPath := filepath.Join(chatDir, "auto_capture_config.json")
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, configData, 0644)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// ProcessMessageForEvents runs the auto-capture pipeline on a single message
// Detects events, classifies them, and queues candidates for review.
// Returns the number of new candidates queued.
func (s *Service) ProcessMessageForEvents(message string, history []string, characterNames []string) (int, error) {
	config, err := s.GetAutoCaptureConfig()
	if err != nil {
		return 0, err
	}
	if !config.Enabled {
		return 0, nil
	}

	// Step 1: Detect events
	detections := DetectEvents(message, history, *config)
	if len(detections) == 0 {
		return 0, nil
	}

	// Step 2: Classify and build candidates
	candidates := BatchClassify(detections, message, characterNames, s.chatID)
	if len(candidates) == 0 {
		return 0, nil
	}

	// Step 3: Queue candidates
	if err := s.AddCandidates(candidates); err != nil {
		return 0, err
	}

	return len(candidates), nil
}

// SaveChatData saves chat-related memory data
func (s *Service) SaveChatData(chatName string, syncAutomatedLore bool) error {
	storage := NewMemoryStorage(s.chatID)

	// Load existing memories
	memories, err := storage.loadMemories()
	if err != nil {
		return fmt.Errorf("failed to load memories: %w", err)
	}

	// For now, just ensure memories file exists
	// In a full implementation, this would sync with automated lore
	return storage.saveMemories(memories)
}

// MemoryStorage methods

// loadMemories loads memories from JSON file
func (ms *MemoryStorage) loadMemories() ([]types.MemoryRecord, error) {
	memoriesPath := filepath.Join(ms.projectRoot, "Chats", ms.chatID, "memories.json")

	data, err := os.ReadFile(memoriesPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return empty array
			return []types.MemoryRecord{}, nil
		}
		return nil, fmt.Errorf("failed to read memories file: %w", err)
	}

	var memories []types.MemoryRecord
	if err := json.Unmarshal(data, &memories); err != nil {
		return nil, fmt.Errorf("failed to parse memories: %w", err)
	}

	return memories, nil
}

// --- Phase 4: Clustering, Dedup, Continuity ---

// ClusterMemories returns semantic clusters of active memories
func (s *Service) ClusterMemories() ([]Cluster, error) {
	storage := NewMemoryStorage(s.chatID)
	memories, err := storage.loadMemories()
	if err != nil {
		return nil, err
	}

	// Only cluster active memories
	var active []types.MemoryRecord
	for _, m := range memories {
		if m.Status == "active" {
			active = append(active, m)
		}
	}

	return ClusterMemories(active, DefaultClusterConfig()), nil
}

// DeduplicateMemories removes duplicate memories and returns merge count
func (s *Service) DeduplicateMemories() (int, error) {
	storage := NewMemoryStorage(s.chatID)
	memories, err := storage.loadMemories()
	if err != nil {
		return 0, err
	}

	result := DeduplicateMemories(memories, DefaultDedupConfig())
	if result.Merged > 0 {
		if err := storage.saveMemories(result.Kept); err != nil {
			return 0, err
		}
	}

	return result.Merged, nil
}

// ExportCharacterMemoriesForContinuity exports memories for cross-scenario use
func (s *Service) ExportCharacterMemoriesForContinuity(characterName string) ([]types.MemoryRecord, error) {
	storage := NewMemoryStorage(s.chatID)
	memories, err := storage.loadMemories()
	if err != nil {
		return nil, err
	}

	cc := NewCharacterContinuity()
	config, _ := cc.LoadContinuityConfig(s.chatID)
	return cc.ExportCharacterMemories(characterName, memories, config), nil
}

// ImportCharacterMemoriesForContinuity imports memories into current chat
func (s *Service) ImportCharacterMemoriesForContinuity(characterName string, imported []types.MemoryRecord) error {
	storage := NewMemoryStorage(s.chatID)
	memories, err := storage.loadMemories()
	if err != nil {
		return err
	}

	cc := NewCharacterContinuity()
	config, _ := cc.LoadContinuityConfig(s.chatID)
	merged := cc.ImportCharacterMemories(characterName, imported, memories, config)

	return storage.saveMemories(merged)
}

// saveMemories saves memories to JSON file
func (ms *MemoryStorage) saveMemories(memories []types.MemoryRecord) error {
	chatDir := filepath.Join(ms.projectRoot, "Chats", ms.chatID)
	memoriesPath := filepath.Join(chatDir, "memories.json")

	// Ensure chat directory exists
	if err := os.MkdirAll(chatDir, 0755); err != nil {
		return fmt.Errorf("failed to create chat directory: %w", err)
	}

	// Create backup
	backupPath := filepath.Join(chatDir, "memories_backup.json")
	os.WriteFile(backupPath, []byte("{}"), 0644)

	data, err := json.MarshalIndent(memories, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal memories: %w", err)
	}

	err = os.WriteFile(memoriesPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save memories: %w", err)
	}

	return nil
}

// applyFilters applies filtering logic to memories
func (ms *MemoryStorage) applyFilters(memories []types.MemoryRecord, filter types.MemoryFilter) []types.MemoryRecord {
	if filter.Query == "" && filter.Scope == "" && filter.Status == "" {
		return memories
	}

	var filtered []types.MemoryRecord

	for _, memory := range memories {
		include := true

		// Query filter
		if filter.Query != "" {
			if !strings.Contains(strings.ToLower(memory.Title), strings.ToLower(filter.Query)) &&
				!strings.Contains(strings.ToLower(memory.Summary), strings.ToLower(filter.Query)) {
				include = false
			}
		}

		// Scope filter
		if filter.Scope != "" && memory.Scope != filter.Scope {
			include = false
		}

		// Status filter
		if filter.Status != "" && memory.Status != filter.Status {
			include = false
		}

		if include {
			filtered = append(filtered, memory)
		}
	}

	return filtered
}
