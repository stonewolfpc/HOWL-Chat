/**
 * Llama Client Implementation
 *
 * This package provides the actual implementation of the llama client
 * using the go-llama.cpp library for real model loading and inference.
 *
 * @package llama
 */

package llama

import (
	"howl-chat/internal/backend/types"
	"sync"

	llama "github.com/go-skynet/go-llama.cpp"
)

// LlamaClient is the actual implementation of the Client interface
type LlamaClient struct {
	mu          sync.RWMutex
	model       *llama.LLama
	modelLoaded bool
	modelPath   string
	progressCB  ProgressCallback
}

// NewLlamaClient creates a new llama client with actual go-llama.cpp implementation
func NewLlamaClient() *LlamaClient {
	return &LlamaClient{
		modelLoaded: false,
		modelPath:   "",
	}
}

// LoadModel loads a model using go-llama.cpp
func (c *LlamaClient) LoadModel(modelPath string, options *LoadOptions) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if options == nil {
		options = NewLoadOptions(modelPath)
	}

	// Configure llama.cpp options
	llamaOptions := []llama.ModelOption{
		llama.SetContext(options.ContextSize),
		llama.SetModelSeed(options.Seed),
		llama.SetGPULayers(options.GPULayers),
		llama.SetMMap(options.UseMMap),
	}

	// Load the model
	model, err := llama.New(modelPath, llamaOptions...)
	if err != nil {
		return types.WrapError(types.ErrorCodeModelLoad, "failed to load model", err)
	}

	c.model = model
	c.modelLoaded = true
	c.modelPath = modelPath

	return nil
}

// UnloadModel unloads the current model
func (c *LlamaClient) UnloadModel() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.modelLoaded {
		return types.ErrModelNotFound("no model loaded")
	}

	if c.model != nil {
		c.model.Close()
	}

	c.model = nil
	c.modelLoaded = false
	c.modelPath = ""

	return nil
}

// IsLoaded returns true if a model is loaded
func (c *LlamaClient) IsLoaded() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.modelLoaded
}

// GetModelInfo returns information about the loaded model
func (c *LlamaClient) GetModelInfo() (*types.Model, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.model == nil {
		return nil, types.ErrModelNotFound("no model loaded")
	}

	return &types.Model{
		Name:     c.modelPath,
		FilePath: c.modelPath,
		Status:   types.ModelStatusLoaded,
	}, nil
}

// Generate generates text completion
func (c *LlamaClient) Generate(prompt string, options *InferenceOptions) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.modelLoaded || c.model == nil {
		return "", types.ErrModelNotFound("no model loaded")
	}

	if options == nil {
		options = NewInferenceOptions()
	}

	// Configure generation options
	llamaOptions := []llama.PredictOption{
		llama.SetTemperature(float32(options.Temperature)),
		llama.SetTopP(float32(options.TopP)),
		llama.SetTopK(options.TopK),
		llama.SetRepeatPenalty(float32(options.RepeatPenalty)),
		llama.SetFrequencyPenalty(float32(options.FrequencyPenalty)),
		llama.SetPresencePenalty(float32(options.PresencePenalty)),
		llama.SetSeed(options.Seed),
		llama.SetTokens(options.NPredict),
		llama.SetStopWords(options.StopStrings...),
	}

	// Generate response
	response, err := c.model.Generate(prompt, llamaOptions...)
	if err != nil {
		return "", types.WrapError(types.ErrorCodeInference, "generation failed", err)
	}

	return response, nil
}

// GenerateStream generates text with streaming
func (c *LlamaClient) GenerateStream(prompt string, options *InferenceOptions, callback TokenCallback) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.modelLoaded || c.model == nil {
		return types.ErrModelNotFound("no model loaded")
	}

	if options == nil {
		options = NewInferenceOptions()
	}

	// Configure generation options
	llamaOptions := []llama.PredictOption{
		llama.SetTemperature(float32(options.Temperature)),
		llama.SetTopP(float32(options.TopP)),
		llama.SetTopK(options.TopK),
		llama.SetRepeatPenalty(float32(options.RepeatPenalty)),
		llama.SetFrequencyPenalty(float32(options.FrequencyPenalty)),
		llama.SetPresencePenalty(float32(options.PresencePenalty)),
		llama.SetSeed(options.Seed),
		llama.SetTokens(options.NPredict),
		llama.SetStopWords(options.StopStrings...),
	}

	// Generate with streaming
	err := c.model.GenerateStream(prompt, func(token string) {
		callback(token, false)
	}, llamaOptions...)

	if err != nil {
		return types.WrapError(types.ErrorCodeInference, "streaming generation failed", err)
	}

	callback("", true)

	return nil
}

// Tokenize converts text to token IDs
func (c *LlamaClient) Tokenize(text string, options *TokenizerOptions) ([]int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.model == nil {
		return nil, types.ErrModelNotFound("no model loaded")
	}

	if options == nil {
		options = NewTokenizerOptions()
	}

	tokens, err := c.model.Tokenize(text, options.AddBOS, options.AddEOS)
	if err != nil {
		return nil, types.WrapError(types.ErrorCodeContext, "tokenization failed", err)
	}

	return tokens, nil
}

// Detokenize converts token IDs back to text
func (c *LlamaClient) Detokenize(tokens []int) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.model == nil {
		return "", types.ErrModelNotFound("no model loaded")
	}

	text, err := c.model.DeTokenize(tokens)
	if err != nil {
		return "", types.WrapError(types.ErrorCodeContext, "detokenization failed", err)
	}

	return text, nil
}

// GetContextSize returns the context window size
func (c *LlamaClient) GetContextSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.model == nil {
		return 0
	}

	return c.model.ContextSize()
}

// GetUsedMemory returns memory usage
func (c *LlamaClient) GetUsedMemory() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.model == nil {
		return 0
	}

	return c.model.MemUsed()
}

// GetTotalMemory returns total memory
func (c *LlamaClient) GetTotalMemory() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.model == nil {
		return 0
	}

	return c.model.MemTotal()
}

// Embedding generates embeddings
func (c *LlamaClient) Embedding(text string) ([]float32, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.model == nil {
		return nil, types.ErrModelNotFound("no model loaded")
	}

	embeddings, err := c.model.Embedding(text)
	if err != nil {
		return nil, types.WrapError(types.ErrorCodeInference, "embedding generation failed", err)
	}

	return embeddings, nil
}

// SetProgressCallback sets a progress callback
func (c *LlamaClient) SetProgressCallback(callback ProgressCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.progressCB = callback
}

// Close releases resources
func (c *LlamaClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.model != nil {
		err := c.model.Close()
		if err != nil {
			return err
		}
	}

	c.model = nil
	c.modelLoaded = false
	c.modelPath = ""

	return nil
}
