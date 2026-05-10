/**
 * Message Types
 * 
 * This package defines all message-related types for the chat system.
 * Messages represent individual chat interactions between users and the AI.
 * 
 * @package types
 */

package types

import "time"

// Role represents the role of a message sender
type Role string

const (
	// RoleUser represents a message from the user
	RoleUser Role = "user"
	// RoleAssistant represents a message from the AI assistant
	RoleAssistant Role = "assistant"
	// RoleSystem represents a system message
	RoleSystem Role = "system"
)

// Message represents a single chat message
type Message struct {
	ID        string    `json:"id"`        // Unique message identifier
	Role      Role      `json:"role"`      // Sender role (user/assistant/system)
	Content   string    `json:"content"`   // Message content
	Timestamp time.Time `json:"timestamp"` // When the message was created
	Tokens    int       `json:"tokens"`    // Number of tokens in the message
	Metadata  Metadata  `json:"metadata"`  // Additional message metadata
}

// Metadata contains optional message metadata
type Metadata struct {
	Model        string            `json:"model,omitempty"`        // Model used for generation
	Temperature  float64           `json:"temperature,omitempty"`  // Generation temperature
	TopP         float64           `json:"top_p,omitempty"`         // Top P sampling
	TopK         int              `json:"top_k,omitempty"`         // Top K sampling
	Images       []ImageReference `json:"images,omitempty"`       // Attached images
	Audio        AudioReference    `json:"audio,omitempty"`        // Attached audio
	Documents    []DocumentRef     `json:"documents,omitempty"`    // Attached documents
}

// ImageReference represents an attached image
type ImageReference struct {
	ID       string `json:"id"`       // Image identifier
	FilePath string `json:"filepath"` // Path to the image file
	Thumbnail string `json:"thumbnail"` // Thumbnail file path
	MimeType string `json:"mimetype"` // Image MIME type
}

// AudioReference represents attached audio
type AudioReference struct {
	ID       string `json:"id"`       // Audio identifier
	FilePath string `json:"filepath"` // Path to the audio file
	Duration int    `json:"duration"` // Duration in seconds
	MimeType string `json:"mimetype"` // Audio MIME type
}

// DocumentReference represents an attached document
type DocumentRef struct {
	ID       string `json:"id"`       // Document identifier
	FilePath string `json:"filepath"` // Path to the document file
	Title    string `json:"title"`    // Document title
	MimeType string `json:"mimetype"` // Document MIME type
}

// MessageList represents a collection of messages
type MessageList struct {
	Messages []Message `json:"messages"` // List of messages
	Total    int       `json:"total"`    // Total number of messages
}

// NewMessage creates a new message with the given role and content
func NewMessage(role Role, content string) *Message {
	return &Message{
		ID:        generateID(),
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
		Tokens:    0, // Will be calculated by tokenizer
		Metadata:  Metadata{},
	}
}

// generateID generates a unique identifier for a message
func generateID() string {
	// Simple ID generation - in production use UUID or similar
	return "msg_" + time.Now().Format("20060102150405")
}
