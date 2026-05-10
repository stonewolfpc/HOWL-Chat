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
	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/model"
	"howl-chat/internal/backend/types"
	"sync"
)

// Service handles chat operations including message processing and AI interaction
type Service struct {
	mu            sync.RWMutex
	context       *types.Context
	modelManager  *model.Manager
	llamaClient   llama.Client
	inferenceOpts *llama.InferenceOptions
	contextOpts   *types.ContextConfig
}

// NewService creates a new chat service
func NewService(modelManager *model.Manager, client llama.Client) *Service {
	return &Service{
		context:       types.NewContext("", 4096),
		modelManager:  modelManager,
		llamaClient:   client,
		inferenceOpts: llama.NewInferenceOptions(),
		contextOpts:   types.DefaultContextConfig(),
	}
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
	
	// Build prompt from context
	prompt := s.buildPrompt()
	
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
	
	// Build prompt from context
	prompt := s.buildPrompt()
	
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

// buildPrompt builds a prompt from the context messages
func (s *Service) buildPrompt() string {
	messages := s.context.GetMessages()
	var prompt string
	
	for _, msg := range messages {
		switch msg.Role {
		case types.RoleSystem:
			prompt += "System: " + msg.Content + "\n"
		case types.RoleUser:
			prompt += "User: " + msg.Content + "\n"
		case types.RoleAssistant:
			prompt += "Assistant: " + msg.Content + "\n"
		}
	}
	
	prompt += "Assistant:"
	return prompt
}

// Close releases all resources
func (s *Service) Close() error {
	return nil
}
