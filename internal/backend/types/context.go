/**
 * Context Types
 * 
 * This package defines all context-related types for the chat system.
 * Context represents the conversation history and state for a chat session.
 * 
 * @package types
 */

package types

import "time"

// Context represents a conversation context with message history
type Context struct {
	ID        string      `json:"id"`        // Unique context identifier
	Messages  []Message   `json:"messages"`  // Message history
	ModelID   string      `json:"model_id"`  // Currently loaded model ID
	MaxTokens int         `json:"max_tokens"` // Maximum context window size
	CreatedAt time.Time   `json:"created_at"` // When the context was created
	UpdatedAt time.Time   `json:"updated_at"` // When the context was last updated
	Metadata  ContextMeta `json:"metadata"`  // Additional context metadata
}

// ContextMeta contains optional context metadata
type ContextMeta struct {
	Title       string    `json:"title,omitempty"`       // Context title
	Description string    `json:"description,omitempty"` // Context description
	Tags        []string  `json:"tags,omitempty"`        // Context tags
	SystemPrompt string  `json:"system_prompt,omitempty"` // System prompt for the context
	WorldID     string   `json:"world_id,omitempty"`     // Associated world ID
	ScenarioID  string   `json:"scenario_id,omitempty"`  // Associated scenario ID
	CharacterID string   `json:"character_id,omitempty"` // Associated character ID
	LorebookIDs []string `json:"lorebook_ids,omitempty"` // Associated lorebook IDs
}

// ContextConfig represents configuration for context management
type ContextConfig struct {
	MaxHistory      int     `json:"max_history"`       // Maximum messages to keep in history
	MaxTokens       int     `json:"max_tokens"`        // Maximum tokens in context
	TruncateStrategy string `json:"truncate_strategy"` // Strategy for truncating context (oldest, smart)
	SystemPrompt    string  `json:"system_prompt"`     // Default system prompt
	Temperature     float64 `json:"temperature"`       // Default temperature
	TopP            float64 `json:"top_p"`             // Default top P
	TopK            int     `json:"top_k"`             // Default top K
	RepeatPenalty   float64 `json:"repeat_penalty"`    // Repetition penalty
}

// ContextStats represents statistics about a context
type ContextStats struct {
	MessageCount    int       `json:"message_count"`    // Total number of messages
	TokenCount      int       `json:"token_count"`      // Total token count
	UserMessages    int       `json:"user_messages"`    // Number of user messages
	AssistantMessages int     `json:"assistant_messages"` // Number of assistant messages
	SystemMessages  int       `json:"system_messages"`  // Number of system messages
	FirstMessage    time.Time `json:"first_message"`    // Timestamp of first message
	LastMessage     time.Time `json:"last_message"`     // Timestamp of last message
	AverageTokens   float64   `json:"average_tokens"`   // Average tokens per message
}

// ContextUpdate represents an update to a context
type ContextUpdate struct {
	ContextID string   `json:"context_id"` // Context to update
	AddMessages []Message `json:"add_messages,omitempty"` // Messages to add
	RemoveMessageIDs []string `json:"remove_message_ids,omitempty"` // Message IDs to remove
	UpdateMetadata *ContextMeta `json:"update_metadata,omitempty"` // Metadata to update
}

// NewContext creates a new context with the given model ID
func NewContext(modelID string, maxTokens int) *Context {
	now := time.Now()
	return &Context{
		ID:        generateContextID(),
		Messages:  []Message{},
		ModelID:   modelID,
		MaxTokens: maxTokens,
		CreatedAt: now,
		UpdatedAt: now,
		Metadata: ContextMeta{
			Title:       "",
			Description: "",
			Tags:        []string{},
			SystemPrompt: "",
			WorldID:     "",
			ScenarioID:  "",
			CharacterID: "",
			LorebookIDs: []string{},
		},
	}
}

// DefaultContextConfig returns default context configuration
func DefaultContextConfig() *ContextConfig {
	return &ContextConfig{
		MaxHistory:      100,
		MaxTokens:       4096,
		TruncateStrategy: "oldest",
		SystemPrompt:    "",
		Temperature:     0.7,
		TopP:            0.9,
		TopK:            40,
		RepeatPenalty:   1.1,
	}
}

// generateContextID generates a unique identifier for a context
func generateContextID() string {
	return "ctx_" + time.Now().Format("20060102150405")
}

// AddMessage adds a message to the context
func (c *Context) AddMessage(message Message) {
	c.Messages = append(c.Messages, message)
	c.UpdatedAt = time.Now()
}

// GetMessages returns all messages in the context
func (c *Context) GetMessages() []Message {
	return c.Messages
}

// GetLastNMessages returns the last N messages from the context
func (c *Context) GetLastNMessages(n int) []Message {
	if len(c.Messages) <= n {
		return c.Messages
	}
	return c.Messages[len(c.Messages)-n:]
}

// GetMessageCount returns the total number of messages in the context
func (c *Context) GetMessageCount() int {
	return len(c.Messages)
}

// GetStats returns statistics about the context
func (c *Context) GetStats() *ContextStats {
	stats := &ContextStats{
		MessageCount:     len(c.Messages),
		TokenCount:       0, // Will be calculated by tokenizer
		UserMessages:     0,
		AssistantMessages: 0,
		SystemMessages:   0,
	}

	for _, msg := range c.Messages {
		stats.TokenCount += msg.Tokens
		
		switch msg.Role {
		case RoleUser:
			stats.UserMessages++
		case RoleAssistant:
			stats.AssistantMessages++
		case RoleSystem:
			stats.SystemMessages++
		}
	}

	if len(c.Messages) > 0 {
		stats.FirstMessage = c.Messages[0].Timestamp
		stats.LastMessage = c.Messages[len(c.Messages)-1].Timestamp
		stats.AverageTokens = float64(stats.TokenCount) / float64(stats.MessageCount)
	}

	return stats
}

// TrimToMaxTokens trims the context to fit within the max token limit
func (c *Context) TrimToMaxTokens(maxTokens int) {
	// This is a placeholder - actual implementation would use tokenizer
	// to calculate token count and trim messages accordingly
	if len(c.Messages) == 0 {
		return
	}
	
	currentTokens := 0
	for i := len(c.Messages) - 1; i >= 0; i-- {
		currentTokens += c.Messages[i].Tokens
		if currentTokens > maxTokens {
			c.Messages = c.Messages[i+1:]
			c.UpdatedAt = time.Now()
			return
		}
	}
}
