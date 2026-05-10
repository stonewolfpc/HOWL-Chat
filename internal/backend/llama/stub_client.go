/**
 * Stub Llama Client
 *
 * This package provides a stub implementation of the llama client interface
 * for testing purposes. This allows the frontend-backend integration to be tested
 * before the actual go-llama.cpp implementation is complete.
 *
 * @package llama
 */

package llama

import (
	"howl-chat/internal/backend/types"
	"time"
)

// StubClient is a stub implementation of the Client interface for testing
type StubClient struct {
	modelLoaded bool
	modelPath   string
}

// NewStubClient creates a new stub client
func NewStubClient() *StubClient {
	return &StubClient{
		modelLoaded: false,
		modelPath:   "",
	}
}

// LoadModel simulates loading a model
func (c *StubClient) LoadModel(modelPath string, options *LoadOptions) error {
	c.modelLoaded = true
	c.modelPath = modelPath
	return nil
}

// UnloadModel simulates unloading the model
func (c *StubClient) UnloadModel() error {
	c.modelLoaded = false
	c.modelPath = ""
	return nil
}

// IsLoaded returns true if a model is loaded
func (c *StubClient) IsLoaded() bool {
	return c.modelLoaded
}

// GetModelInfo returns mock model information
func (c *StubClient) GetModelInfo() (*types.Model, error) {
	if !c.modelLoaded {
		return nil, types.ErrModelNotFound("no model loaded")
	}

	return &types.Model{
		Name:     "Stub Model",
		FilePath: c.modelPath,
		Status:   types.ModelStatusLoaded,
		LoadedAt: time.Now(),
	}, nil
}

// Generate simulates text generation with a mock response
func (c *StubClient) Generate(prompt string, options *InferenceOptions) (string, error) {
	if !c.modelLoaded {
		return "", types.ErrModelNotFound("no model loaded")
	}

	// Return a mock response for testing
	return "This is a stub response from the llama client. The actual go-llama.cpp implementation will be added later.", nil
}

// GenerateStream simulates streaming text generation
func (c *StubClient) GenerateStream(prompt string, options *InferenceOptions, callback TokenCallback) error {
	if !c.modelLoaded {
		return types.ErrModelNotFound("no model loaded")
	}

	// Simulate streaming by calling the callback with tokens
	tokens := []string{"This", " is", " a", " stub", " response", " from", " the", " llama", " client.", " The", " actual", " go-llama.cpp", " implementation", " will", " be", " added", " later."}

	for i, token := range tokens {
		callback(token, i == len(tokens)-1)
	}

	return nil
}

// Tokenize simulates tokenization
func (c *StubClient) Tokenize(text string, options *TokenizerOptions) ([]int, error) {
	if !c.modelLoaded {
		return nil, types.ErrModelNotFound("no model loaded")
	}

	// Return mock token IDs
	tokens := make([]int, len(text))
	for i := range text {
		tokens[i] = int(text[i]) + 1000
	}
	return tokens, nil
}

// Detokenize simulates detokenization
func (c *StubClient) Detokenize(tokens []int) (string, error) {
	if !c.modelLoaded {
		return "", types.ErrModelNotFound("no model loaded")
	}

	// Return mock text from tokens
	text := ""
	for _, token := range tokens {
		if token >= 1000 {
			text += string(rune(token - 1000))
		}
	}
	return text, nil
}

// GetContextSize returns mock context size
func (c *StubClient) GetContextSize() int {
	return 2048
}

// GetUsedMemory returns mock memory usage
func (c *StubClient) GetUsedMemory() int64 {
	if c.modelLoaded {
		return 1024 * 1024 * 1024 // 1GB
	}
	return 0
}

// GetTotalMemory returns mock total memory
func (c *StubClient) GetTotalMemory() int64 {
	return 8 * 1024 * 1024 * 1024 // 8GB
}

// Embedding simulates embedding generation
func (c *StubClient) Embedding(text string) ([]float32, error) {
	if !c.modelLoaded {
		return nil, types.ErrModelNotFound("no model loaded")
	}

	// Return mock embedding vector
	embedding := make([]float32, 512)
	for i := range embedding {
		embedding[i] = 0.1
	}
	return embedding, nil
}

// SetProgressCallback sets a progress callback
func (c *StubClient) SetProgressCallback(callback ProgressCallback) {
	// Stub implementation - no actual progress tracking
}

// Close releases resources
func (c *StubClient) Close() error {
	c.modelLoaded = false
	c.modelPath = ""
	return nil
}
