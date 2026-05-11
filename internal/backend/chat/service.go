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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/backend/model"
	"howl-chat/internal/backend/types"
	entities "howl-chat/internal/backend/types/entities"
)

// Service handles chat operations including message processing and AI interaction
type Service struct {
	mu            sync.RWMutex
	context       *types.Context
	modelManager  *model.Manager
	llamaClient   llama.Client
	inferenceOpts *llama.InferenceOptions
	contextOpts   *types.ContextConfig

	worldName      string
	scenarioName   string
	characterNames []string

	worldData     *entities.World
	scenarioData  *entities.Scenario
	characterData map[string]*entities.Character // Map character ID to character data

	// Lorebook and Memory Management
	manualLore        []lorebook.Entry
	automatedLore     []lorebook.Entry
	memorySnapshot    []lorebook.Entry // This will represent memories_snapshot.json
	triggeredMemories map[string]bool  // To track once_per_chat, once_per_scene, once_per_message
}

// NewService creates a new chat service
func NewService(modelManager *model.Manager, client llama.Client, worldName, scenarioName string, characterNames []string, chatName string) (*Service, error) {
	s := &Service{
		context:           types.NewContext("", 4096),
		modelManager:      modelManager,
		llamaClient:       client,
		inferenceOpts:     llama.NewInferenceOptions(),
		contextOpts:       types.DefaultContextConfig(),
		worldData:         nil,
		scenarioData:      nil,
		characterData:     make(map[string]*entities.Character),
		manualLore:        []lorebook.Entry{},
		automatedLore:     []lorebook.Entry{},
		memorySnapshot:    []lorebook.Entry{},
		triggeredMemories: make(map[string]bool),
	}

	if err := s.LoadChatData(worldName, scenarioName, characterNames, chatName); err != nil {
		return nil, err
	}

	return s, nil
}

// LoadChatData loads all relevant lorebook and memory data for a chat session.
func (s *Service) LoadChatData(worldName, scenarioName string, characterNames []string, chatName string) error {
	projectRoot := "d:\\HOWL_Chat"
	s.worldName = worldName
	s.scenarioName = scenarioName
	s.characterNames = characterNames

	// Load world, scenario, and character data
	worldFilePath := filepath.Join(projectRoot, "Worlds", worldName, "world.json")
	worldData, err := loadWorldData(worldFilePath)
	if err != nil {
		return err
	}
	s.worldData = worldData

	scenarioFilePath := filepath.Join(projectRoot, "Scenarios", scenarioName, "scenario.json")
	scenarioData, err := loadScenarioData(scenarioFilePath)
	if err != nil {
		return err
	}
	s.scenarioData = scenarioData

	s.characterData = make(map[string]*entities.Character)
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

	// Load manual lore
	worldLorePath := filepath.Join(projectRoot, "Worlds", worldName, "world_lore.json")
	scenarioLorePath := filepath.Join(projectRoot, "Scenarios", scenarioName, "scenario_lore.json")
	var charLorePaths []string
	for _, char := range characterNames {
		charLorePaths = append(charLorePaths, filepath.Join(projectRoot, "Characters", char, "character_lore.json"))
	}

	var allManualLore []lorebook.Entry

	worldManualLore, err := loadEntriesFromFile(worldLorePath)
	if err != nil {
		return err
	}
	allManualLore = append(allManualLore, worldManualLore...)

	scenarioManualLore, err := loadEntriesFromFile(scenarioLorePath)
	if err != nil {
		return err
	}
	allManualLore = append(allManualLore, scenarioManualLore...)

	for _, p := range charLorePaths {
		charManualLore, err := loadEntriesFromFile(p)
		if err != nil {
			return err
		}
		allManualLore = append(allManualLore, charManualLore...)
	}
	s.manualLore = allManualLore

	// Load automated lore (global)
	automatedWorldLorePath := filepath.Join(projectRoot, "Worlds", worldName, "automated_world_lore.json")
	automatedScenarioLorePath := filepath.Join(projectRoot, "Scenarios", scenarioName, "automated_scenario_lore.json")
	var automatedCharLorePaths []string
	for _, char := range characterNames {
		automatedCharLorePaths = append(automatedCharLorePaths, filepath.Join(projectRoot, "Characters", char, "automated_character_lore.json"))
	}

	var allAutomatedLore []lorebook.Entry

	automatedWorldLore, err := loadEntriesFromFile(automatedWorldLorePath)
	if err != nil {
		return err
	}
	allAutomatedLore = append(allAutomatedLore, automatedWorldLore...)

	automatedScenarioLore, err := loadEntriesFromFile(automatedScenarioLorePath)
	if err != nil {
		return err
	}
	allAutomatedLore = append(allAutomatedLore, automatedScenarioLore...)

	for _, p := range automatedCharLorePaths {
		automatedCharLore, err := loadEntriesFromFile(p)
		if err != nil {
			return err
		}
		allAutomatedLore = append(allAutomatedLore, automatedCharLore...)
	}
	s.automatedLore = allAutomatedLore

	// Load memory snapshot for this chat
	memoriesSnapshotPath := filepath.Join(projectRoot, "Chats", chatName, "memories_snapshot.json")
	existingMemories, err := loadMemoriesSnapshot(memoriesSnapshotPath)
	if err != nil {
		return err
	}

	if len(existingMemories) == 0 {
		// If no existing memories snapshot, create one from automated lore
		var chatMemories []lorebook.Entry
		for _, entry := range s.automatedLore {
			// Only copy automated lore that is relevant to the current chat's scope
			if (entry.Scope == lorebook.ScopeWorld && entry.OwnerID == worldName) ||
				(entry.Scope == lorebook.ScopeScenario && entry.OwnerID == scenarioName) ||
				(entry.Scope == lorebook.ScopeCharacter && contains(characterNames, entry.OwnerID)) {
				chatMemories = append(chatMemories, entry)
			}
		}
		s.memorySnapshot = chatMemories
	} else {
		s.memorySnapshot = existingMemories
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// SaveChatData saves all relevant chat data to JSON files.
func (s *Service) SaveChatData(chatName string, syncAutomatedLore bool) error {
	projectRoot := "d:\\HOWL_Chat"

	// Save chat.json
	chatFilePath := filepath.Join(projectRoot, "Chats", chatName, "chat.json")
	chatData := map[string]interface{}{
		"chat_name":  chatName,
		"seed":       987654321, // Placeholder for now
		"world":      s.worldName,
		"scenario":   s.scenarioName,
		"characters": s.characterNames,
		"version":    1,
	}
	chatJSON, err := json.MarshalIndent(chatData, "", "  ")
	if err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to marshal chat.json", err)
	}
	if err := os.WriteFile(chatFilePath, chatJSON, 0644); err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to write chat.json", err)
	}

	// Save history.json
	historyFilePath := filepath.Join(projectRoot, "Chats", chatName, "history.json")
	historyJSON, err := json.MarshalIndent(s.context.GetMessages(), "", "  ")
	if err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to marshal history.json", err)
	}
	if err := os.WriteFile(historyFilePath, historyJSON, 0644); err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to write history.json", err)
	}

	// Save memories_snapshot.json
	memoriesSnapshotFilePath := filepath.Join(projectRoot, "Chats", chatName, "memories_snapshot.json")
	memoriesSnapshotData := struct {
		WorldMemories     []lorebook.Entry
		ScenarioMemories  []lorebook.Entry
		CharacterMemories map[string][]lorebook.Entry
	}{
		WorldMemories:     []lorebook.Entry{},
		ScenarioMemories:  []lorebook.Entry{},
		CharacterMemories: make(map[string][]lorebook.Entry),
	}

	for _, mem := range s.memorySnapshot {
		switch mem.Scope {
		case lorebook.ScopeWorld:
			memoriesSnapshotData.WorldMemories = append(memoriesSnapshotData.WorldMemories, mem)
		case lorebook.ScopeScenario:
			memoriesSnapshotData.ScenarioMemories = append(memoriesSnapshotData.ScenarioMemories, mem)
		case lorebook.ScopeCharacter:
			if mem.OwnerID != "" {
				memoriesSnapshotData.CharacterMemories[mem.OwnerID] = append(memoriesSnapshotData.CharacterMemories[mem.OwnerID], mem)
			}
		}
	}

	memoriesSnapshotJSON, err := json.MarshalIndent(memoriesSnapshotData, "", "  ")
	if err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to marshal memories_snapshot.json", err)
	}
	if err := os.WriteFile(memoriesSnapshotFilePath, memoriesSnapshotJSON, 0644); err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to write memories_snapshot.json", err)
	}

	// Optional: sync selected memories back to global automated_*_lore
	if syncAutomatedLore {
		for _, mem := range s.memorySnapshot {
			// For now, sync all automated memories back if they are marked as automated
			if mem.Automated {
				switch mem.Scope {
				case lorebook.ScopeWorld:
					// Load, update, and save automated_world_lore.json
					filePath := filepath.Join(projectRoot, "Worlds", s.worldName, "automated_world_lore.json")
					updateAutomatedLoreFile(filePath, mem)
				case lorebook.ScopeScenario:
					// Load, update, and save automated_scenario_lore.json
					filePath := filepath.Join(projectRoot, "Scenarios", s.scenarioName, "automated_scenario_lore.json")
					updateAutomatedLoreFile(filePath, mem)
				case lorebook.ScopeCharacter:
					if mem.OwnerID != "" {
						filePath := filepath.Join(projectRoot, "Characters", mem.OwnerID, "automated_character_lore.json")
						updateAutomatedLoreFile(filePath, mem)
					}
				}
			}
		}
	}

	return nil
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
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return types.WrapError(types.ErrorCodeInternal, "failed to write automated lore file", err)
	}
	return nil
}

// SendMessage sends a user message and generates an AI response
func (s *Service) SendMessage(message string) (*types.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if model is loaded
	if !s.modelManager.IsModelLoaded() {
		return nil, types.ErrModelNotFound("no model loaded")
	}

	// Create user message
	userMsg := types.NewMessage(types.RoleUser, message)
	s.context.AddMessage(*userMsg)

	// Trim context if needed
	if s.context.GetMessageCount() > s.contextOpts.MaxHistory {
		s.context.TrimToMaxTokens(s.contextOpts.MaxTokens)
	}

	// Resolve lore and memories
	allEntries := append(s.manualLore, s.memorySnapshot...)
	resolveReq := lorebook.ResolveRequest{
		Message:         userMsg.Content,
		History:         convertToLorebookMessages(s.context.GetMessages()),
		ActiveWorld:     s.worldName,
		ActiveScenario:  s.scenarioName,
		ActiveCharacter: s.characterNames[0], // Assuming the first character is the active one for now
		Actor:           lorebook.TriggerUser,
		Triggered:       s.triggeredMemories,
	}

	resolvedEntries := lorebook.Resolve(allEntries, resolveReq)

	// Update triggered memories for frequency control
	for _, entry := range resolvedEntries {
		if entry.TriggerFrequency == lorebook.FrequencyOncePerChat ||
			entry.TriggerFrequency == lorebook.FrequencyOncePerScene ||
			entry.TriggerFrequency == lorebook.FrequencyOncePerMessage {
			s.triggeredMemories[entry.ID] = true
			// Update LastTriggeredAt for the entry in memorySnapshot
			for i, mem := range s.memorySnapshot {
				if mem.ID == entry.ID {
					now := time.Now()
					s.memorySnapshot[i].LastTriggeredAt = &now
					break
				}
			}
		}
	}

	// Build prompt from context and resolved entries
	prompt := s.buildPrompt(s.characterNames[0], s.context.GetMessageCount() == 1, resolvedEntries)

	// Generate response
	response, err := s.llamaClient.Generate(prompt, s.inferenceOpts)
	if err != nil {
		return nil, types.WrapError(types.ErrorCodeInference,
			"failed to generate response", err)
	}

	// Create assistant message
	assistantMsg := types.NewMessage(types.RoleAssistant, response)
	s.context.AddMessage(*assistantMsg)

	return assistantMsg, nil
}

// SendMessageStream sends a message with streaming response
func (s *Service) SendMessageStream(message string, callback llama.TokenCallback) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if model is loaded
	if !s.modelManager.IsModelLoaded() {
		return types.ErrModelNotFound("no model loaded")
	}

	// Create user message
	userMsg := types.NewMessage(types.RoleUser, message)
	s.context.AddMessage(*userMsg)

	// Trim context if needed
	if s.context.GetMessageCount() > s.contextOpts.MaxHistory {
		s.context.TrimToMaxTokens(s.contextOpts.MaxTokens)
	}

	// Resolve lore and memories
	allEntries := append(s.manualLore, s.memorySnapshot...)
	resolveReq := lorebook.ResolveRequest{
		Message:         userMsg.Content,
		History:         convertToLorebookMessages(s.context.GetMessages()),
		ActiveWorld:     s.worldName,
		ActiveScenario:  s.scenarioName,
		ActiveCharacter: s.characterNames[0], // Assuming the first character is the active one for now
		Actor:           lorebook.TriggerUser,
		Triggered:       s.triggeredMemories,
	}

	resolvedEntries := lorebook.Resolve(allEntries, resolveReq)

	// Update triggered memories for frequency control
	for _, entry := range resolvedEntries {
		if entry.TriggerFrequency == lorebook.FrequencyOncePerChat ||
			entry.TriggerFrequency == lorebook.FrequencyOncePerScene ||
			entry.TriggerFrequency == lorebook.FrequencyOncePerMessage {
			s.triggeredMemories[entry.ID] = true
			// Update LastTriggeredAt for the entry in memorySnapshot
			for i, mem := range s.memorySnapshot {
				if mem.ID == entry.ID {
					now := time.Now()
					s.memorySnapshot[i].LastTriggeredAt = &now
					break
				}
			}
		}
	}

	// Build prompt from context and resolved entries
	prompt := s.buildPrompt(s.characterNames[0], s.context.GetMessageCount() == 1, resolvedEntries)

	// Generate streaming response
	err := s.llamaClient.GenerateStream(prompt, s.inferenceOpts, callback)
	if err != nil {
		return types.WrapError(types.ErrorCodeInference,
			"failed to generate streaming response", err)
	}

	return nil
}

// AddSystemMessage adds a system message to the context
func (s *Service) AddSystemMessage(message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	systemMsg := types.NewMessage(types.RoleSystem, message)
	s.context.AddMessage(*systemMsg)

	return nil
}

// GetContext returns the current conversation context
func (s *Service) GetContext() *types.Context {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.context
}

// GetMessages returns all messages in the context
func (s *Service) GetMessages() []types.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.context.GetMessages()
}

// GetMessageCount returns the number of messages in the context
func (s *Service) GetMessageCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.context.GetMessageCount()
}

// ClearContext clears all messages from the context
func (s *Service) ClearContext() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.context = types.NewContext(s.context.ModelID, s.context.MaxTokens)
	return nil
}

// SetInferenceOptions sets the inference options
func (s *Service) SetInferenceOptions(options *llama.InferenceOptions) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := options.Validate(); err != nil {
		return types.WrapError(types.ErrorCodeInvalidInput,
			"invalid inference options", err)
	}

	s.inferenceOpts = options
	return nil
}

// GetInferenceOptions returns the current inference options
func (s *Service) GetInferenceOptions() *llama.InferenceOptions {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.inferenceOpts
}

// SetContextConfig sets the context configuration
func (s *Service) SetContextConfig(config *types.ContextConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.contextOpts = config
	s.context.MaxTokens = config.MaxTokens

	return nil
}

// GetContextConfig returns the current context configuration
func (s *Service) GetContextConfig() *types.ContextConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.contextOpts
}

// GetContextStats returns statistics about the current context
func (s *Service) GetContextStats() *types.ContextStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.context.GetStats()
}

// LoadContext loads a context from a saved state
func (s *Service) LoadContext(context *types.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.context = context
	return nil
}

// ExportContext exports the current context
func (s *Service) ExportContext() *types.Context {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.context
}

// convertToLorebookMessages converts a slice of types.Message to a slice of lorebook.Message
func convertToLorebookMessages(tMsgs []types.Message) []lorebook.Message {
	lMsgs := make([]lorebook.Message, len(tMsgs))
	for i, tMsg := range tMsgs {
		lMsgs[i] = lorebook.Message{Role: string(tMsg.Role), Content: tMsg.Content}
	}
	return lMsgs
}

// buildPrompt builds a prompt from the context messages and resolved lore entries
func (s *Service) buildPrompt(activeCharacterName string, isFirstTurn bool, resolvedEntries []lorebook.ResolvedEntry) string {
	var prompt strings.Builder

	// 1. System framing (who this AI is, what this chat is)
	prompt.WriteString(fmt.Sprintf("You are %s, a character in HOWL Chat, a text-based adventure. Reply as %s only.\n\n", activeCharacterName, activeCharacterName))

	// 2. World context
	if s.worldData != nil {
		if isFirstTurn {
			prompt.WriteString("[WORLD]\n")
			prompt.WriteString(fmt.Sprintf("You are in the world of %s. %s\n", s.worldData.Name, s.worldData.Description))
			if len(s.worldData.KeyFacts) > 0 {
				prompt.WriteString("Key facts: " + strings.Join(s.worldData.KeyFacts, ", ") + "\n")
			}
			prompt.WriteString("\n")
		} else {
			// Summarized world context
			prompt.WriteString("[WORLD SUMMARY]\n")
			prompt.WriteString(fmt.Sprintf("World: %s. %s\n", s.worldData.Name, s.worldData.Description)) // Simplified summary for now
			prompt.WriteString("\n")
		}
	}

	// 3. Scenario context
	if s.scenarioData != nil {
		if isFirstTurn {
			prompt.WriteString("[SCENARIO]\n")
			prompt.WriteString(fmt.Sprintf("Current scenario: %s. %s\n", s.scenarioData.Name, s.scenarioData.Description))
			if len(s.scenarioData.KeyDetails) > 0 {
				prompt.WriteString("Key details: " + strings.Join(s.scenarioData.KeyDetails, ", ") + "\n")
			}
			prompt.WriteString("\n")
		} else {
			// Summarized scenario context
			prompt.WriteString("[SCENARIO SUMMARY]\n")
			prompt.WriteString(fmt.Sprintf("Scenario: %s. %s\n", s.scenarioData.Name, s.scenarioData.Description)) // Simplified summary for now
			prompt.WriteString("\n")
		}
	}

	// 4. Character sheet for the current AI
	if activeChar, ok := s.characterData[activeCharacterName]; ok {
		prompt.WriteString(fmt.Sprintf("[YOU ARE: %s]\n", activeCharacterName))
		prompt.WriteString(fmt.Sprintf("Personality: %s\n", activeChar.Personality))
		prompt.WriteString(fmt.Sprintf("Backstory: %s\n", activeChar.Backstory))
		if len(activeChar.Goals) > 0 {
			prompt.WriteString("Goals: " + strings.Join(activeChar.Goals, ", ") + "\n")
		}
		if len(activeChar.Relationships) > 0 {
			prompt.WriteString("Relationships: " + strings.Join(activeChar.Relationships, ", ") + "\n")
		}
		prompt.WriteString("\n")
	}

	// 5. Other characters’ summaries
	if len(s.characterNames) > 1 {
		prompt.WriteString("[OTHER PARTICIPANTS]\n")
		for _, charName := range s.characterNames {
			if charName != activeCharacterName {
				if otherChar, ok := s.characterData[charName]; ok {
					// For now, a very basic summary. This could be enhanced.
					prompt.WriteString(fmt.Sprintf("%s: %s\n", charName, otherChar.Personality))
				}
			}
		}
		prompt.WriteString("\n")
	}

	// 6. Memories / lore (resolved for this turn)
	if len(resolvedEntries) > 0 {
		prompt.WriteString("[RELEVANT MEMORIES]\n")
		for _, entry := range resolvedEntries {
			prompt.WriteString(fmt.Sprintf("- %s\n", entry.Content))
		}
		prompt.WriteString("\n")
	}

	// 7. Conversation history (group-style, with roles)
	prompt.WriteString("[CONVERSATION SO FAR]\n")
	messages := s.context.GetMessages()
	for _, msg := range messages {
		speakerName := ""
		switch msg.Role {
		case types.RoleUser:
			speakerName = "User"
		case types.RoleAssistant:
			// Need to determine which assistant spoke. Assuming message metadata might have this.
			// For now, if activeCharacterName is in chat, assume it's them. Otherwise, a generic "Assistant".
			// A better approach would be to store speaker name in message metadata.
			if contains(s.characterNames, activeCharacterName) {
				speakerName = activeCharacterName
			} else {
				speakerName = "Assistant"
			}
		case types.RoleSystem:
			speakerName = "System"
		default:
			speakerName = string(msg.Role)
		}
		prompt.WriteString(fmt.Sprintf("%s: %s\n", speakerName, msg.Content))
	}
	prompt.WriteString("\n")

	// 8. User’s latest message + instruction
	prompt.WriteString("[YOUR TASK]\n")
	prompt.WriteString(fmt.Sprintf("You are %s. Reply as %s only.\n", activeCharacterName, activeCharacterName))
	prompt.WriteString("Do not speak for the User or other characters.\n")

	return prompt.String()
}

// Close releases all resources
func (s *Service) Close() error {
	return nil
}
