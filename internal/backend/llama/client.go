/**
 * Llama Client Interface
 * 
 * This package defines the interface for llama model operations.
 * Provides abstraction for model loading, inference, and management.
 * 
 * @package llama
 */

package llama

import (
	"howl-chat/internal/backend/types"
)

// Client represents the interface for llama model operations
type Client interface {
	// LoadModel loads a model from the given path with the specified options
	LoadModel(modelPath string, options *LoadOptions) error
	
	// UnloadModel unloads the currently loaded model
	UnloadModel() error
	
	// IsLoaded returns true if a model is currently loaded
	IsLoaded() bool
	
	// GetModelInfo returns information about the currently loaded model
	GetModelInfo() (*types.Model, error)
	
	// Generate generates text completion for the given prompt
	Generate(prompt string, options *InferenceOptions) (string, error)
	
	// GenerateStream generates text completion with streaming output
	GenerateStream(prompt string, options *InferenceOptions, callback TokenCallback) error
	
	// Tokenize converts text to token IDs
	Tokenize(text string, options *TokenizerOptions) ([]int, error)
	
	// Detokenize converts token IDs back to text
	Detokenize(tokens []int) (string, error)
	
	// GetContextSize returns the maximum context window size
	GetContextSize() int
	
	// GetUsedMemory returns the amount of memory used by the model
	GetUsedMemory() int64
	
	// GetTotalMemory returns the total memory available
	GetTotalMemory() int64
	
	// Embedding generates embeddings for the given text
	Embedding(text string) ([]float32, error)
	
	// SetProgressCallback sets a callback for loading progress
	SetProgressCallback(callback ProgressCallback)
	
	// Close releases all resources associated with the client
	Close() error
}

// TokenCallback is called for each generated token during streaming
type TokenCallback func(token string, done bool)

// ClientConfig represents configuration for creating a client
type ClientConfig struct {
	// ModelPath is the path to the model file
	ModelPath string
	
	// LoadOptions are the options for loading the model
	LoadOptions *LoadOptions
	
	// ContextOptions are the options for context management
	ContextOptions *ContextOptions
	
	// EnableGPU enables GPU acceleration if available
	EnableGPU bool
	
	// EnableFlashAttention enables flash attention for faster inference
	EnableFlashAttention bool
}

// NewClient creates a new llama client with the given configuration
func NewClient(config *ClientConfig) (Client, error) {
	// This is a placeholder - actual implementation will use go-llama.cpp
	return nil, types.ErrInternal("client not implemented")
}

// Validate validates the client configuration
func (c *ClientConfig) Validate() error {
	if c.ModelPath == "" {
		return types.ErrInvalidInput("model path cannot be empty")
	}
	
	if c.LoadOptions == nil {
		c.LoadOptions = NewLoadOptions(c.ModelPath)
	}
	
	if c.ContextOptions == nil {
		c.ContextOptions = NewContextOptions()
	}
	
	return c.LoadOptions.Validate()
}

// DefaultClientConfig returns a default client configuration
func DefaultClientConfig(modelPath string) *ClientConfig {
	return &ClientConfig{
		ModelPath:            modelPath,
		LoadOptions:          NewLoadOptions(modelPath),
		ContextOptions:       NewContextOptions(),
		EnableGPU:            false,
		EnableFlashAttention: false,
	}
}
