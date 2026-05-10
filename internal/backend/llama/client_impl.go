/**
 * Llama Client Implementation
 *
 * This package provides the actual implementation of the llama client
 * using the gollama.cpp library for real model loading and inference.
 *
 * @package llama
 */

package llama

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"howl-chat/internal/backend/gguf"
	"howl-chat/internal/backend/types"

	gollama "github.com/dianlight/gollama.cpp"
)

// LlamaClient is the actual implementation of the Client interface
type LlamaClient struct {
	mu           sync.RWMutex
	model        gollama.LlamaModel
	ctx          gollama.LlamaContext
	modelLoaded  bool
	modelPath    string
	progressCB   ProgressCallback
	modelProfile *types.ModelProfile // Store model metadata
}

// localLibsDir returns the absolute path to the bundled gollama.cpp libs for this platform.
func localLibsDir() string {
	// Walk up from this source file to find the project root (contains gollama.cpp/libs)
	_, selfPath, _, _ := runtime.Caller(0)
	dir := filepath.Dir(selfPath)
	for i := 0; i < 6; i++ {
		candidate := filepath.Join(dir, "gollama.cpp", "libs",
			fmt.Sprintf("windows_amd64_b6862"))
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		dir = filepath.Dir(dir)
	}
	// Fall back to CWD-relative path
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "gollama.cpp", "libs", "windows_amd64_b6862")
}

// initGollamaConfig configures gollama to use the bundled b6862 libs.
// Note: b6862 does NOT support Gemma-4. Use Llama 3, Gemma 2, Mistral, etc.
// Returns the libs directory path.
func initGollamaConfig() string {
	libsDir := localLibsDir()
	llaDLL := filepath.Join(libsDir, "llama.dll")
	if _, err := os.Stat(llaDLL); err != nil {
		// libs not found via source path; try exe-relative path (production)
		exeDir := filepath.Dir(os.Args[0])
		allDLLs := filepath.Join(exeDir, "llama.dll")
		if _, err2 := os.Stat(allDLLs); err2 == nil {
			libsDir = exeDir
			llaDLL = allDLLs
		}
	}
	cfg := gollama.GetGlobalConfig()
	if cfg == nil {
		cfg = gollama.DefaultConfig()
	}
	if _, err := os.Stat(llaDLL); err == nil {
		cfg.LibraryPath = llaDLL
		cfg.CacheDir = filepath.Dir(libsDir)
	}
	_ = gollama.SetGlobalConfig(cfg)
	fmt.Printf("INFO: gollama libs dir: %s\n", libsDir)
	return libsDir
}

// NewLlamaClient creates a new llama client with actual gollama.cpp implementation
func NewLlamaClient() *LlamaClient {
	libsDir := initGollamaConfig()
	if err := gollama.Backend_init(); err != nil {
		fmt.Printf("WARNING: gollama backend init: %v\n", err)
	}
	if err := gollama.Ggml_backend_load_all_from_path(libsDir); err != nil {
		fmt.Printf("WARNING: ggml backend load: %v\n", err)
	} else {
		fmt.Printf("INFO: ggml backends loaded from: %s\n", libsDir)
	}
	return &LlamaClient{
		modelLoaded: false,
		modelPath:   "",
	}
}

// resolveModelPath returns an absolute path for the model, searching the models dir if needed.
func resolveModelPath(modelPath string) string {
	if filepath.IsAbs(modelPath) {
		return modelPath
	}
	// Search next to exe, then CWD, then models/ subdir
	exeDir := filepath.Dir(os.Args[0])
	cwd, _ := os.Getwd()
	for _, base := range []string{cwd, exeDir, filepath.Join(cwd, "models"), filepath.Join(exeDir, "models")} {
		candidate := filepath.Join(base, modelPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return filepath.Join(cwd, "models", modelPath)
}

// LoadModel loads a model using gollama.cpp
func (c *LlamaClient) LoadModel(modelPath string, options *LoadOptions) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	modelPath = resolveModelPath(modelPath)

	if options == nil {
		options = NewLoadOptions(modelPath)
	}

	// Ensure backend is initialized
	if err := gollama.Backend_init(); err != nil {
		fmt.Printf("ERROR: Failed to initialize backend: %v\n", err)
		return types.WrapError(types.ErrorCodeModelLoad, "failed to initialize backend", err)
	}

	// Read GGUF metadata
	ggufReader := gguf.NewReader()
	fmt.Printf("INFO: Reading GGUF metadata from: %s\n", modelPath)
	profile, err := ggufReader.ReadProfile(modelPath)
	if err != nil {
		// Log warning but continue - metadata might not be critical
		fmt.Printf("WARNING: failed to read GGUF metadata: %v\n", err)
		profile = nil
	}
	c.modelProfile = profile

	// Configure model parameters
	modelParams := gollama.Model_default_params()
	modelParams.NGpuLayers = int32(options.GPULayers)
	modelParams.UseMmap = uint8(0)
	if options.UseMMap {
		modelParams.UseMmap = uint8(1)
	}
	modelParams.UseMlock = uint8(0)
	if options.UseMlock {
		modelParams.UseMlock = uint8(1)
	}
	modelParams.VocabOnly = uint8(0)
	if options.VocabOnly {
		modelParams.VocabOnly = uint8(1)
	}
	modelParams.SplitMode = gollama.LlamaSplitMode(options.SplitMode)
	modelParams.MainGpu = int32(options.MainGPU)

	// Load the model
	fmt.Printf("INFO: Loading model from: %s\n", modelPath)
	fmt.Printf("INFO: GPU layers: %d, UseMMap: %v, UseMlock: %v\n", options.GPULayers, options.UseMMap, options.UseMlock)
	model, err := gollama.Model_load_from_file(modelPath, modelParams)
	if err != nil {
		fmt.Printf("ERROR: Failed to load model from file: %v\n", err)
		return types.WrapError(types.ErrorCodeModelLoad, "failed to load model", err)
	}
	fmt.Printf("INFO: Model loaded successfully\n")

	// Build context params from scratch with safe values.
	// Do NOT call Context_default_params() - it reads struct layout from the DLL
	// which is misaligned against gollama's b6862 Go struct when using newer DLLs.
	contextParams := gollama.LlamaContextParams{
		NSeqMax:         1, // must be <= 256; 1 = single sequence chat
		PoolingType:     gollama.LLAMA_POOLING_TYPE_UNSPECIFIED,
		AttentionType:   gollama.LLAMA_ATTENTION_TYPE_CAUSAL,
		RopeScalingType: gollama.LLAMA_ROPE_SCALING_TYPE_UNSPECIFIED,
		DefragThold:     -1.0,
		Offload_kqv:     1,
		NoPerf:          0,
		OpOffload:       1,
	}

	// Use context size from options, or fall back to GGUF metadata, or default
	contextSize := options.ContextSize
	if contextSize <= 0 && profile != nil && profile.UsableContext > 0 {
		contextSize = profile.UsableContext
	}
	if contextSize <= 0 {
		contextSize = 4096 // default
	}

	contextParams.NCtx = uint32(contextSize)
	contextParams.NThreads = int32(options.Threads)
	contextParams.NBatch = uint32(512)  // batch size
	contextParams.NUbatch = uint32(512) // micro-batch size
	contextParams.NSeqMax = uint32(1)   // single sequence (must be <= 256)
	contextParams.NThreadsBatch = int32(options.Threads)
	contextParams.PoolingType = gollama.LLAMA_POOLING_TYPE_UNSPECIFIED // use model default
	contextParams.RopeScalingType = gollama.LlamaRopeScalingType(options.RopeScalingType)
	contextParams.RopeFreqBase = float32(options.RopeFreqBase)
	contextParams.RopeFreqScale = float32(options.RopeFreqScale)

	// Use GGUF metadata for RoPE settings if available
	if profile != nil {
		if profile.RopeBase > 0 {
			contextParams.RopeFreqBase = float32(profile.RopeBase)
		}
		if profile.RopeFactor > 0 {
			contextParams.RopeFreqScale = float32(profile.RopeFactor)
		}
	}

	contextParams.YarnExtFactor = float32(options.YarnExtFactor)
	contextParams.YarnAttnFactor = float32(options.YarnAttnFactor)
	contextParams.YarnBetaFast = float32(options.YarnBetaFast)
	contextParams.YarnBetaSlow = float32(options.YarnBetaSlow)
	contextParams.YarnOrigCtx = uint32(options.YarnOrigCtx)

	// Create context
	ctx, err := gollama.Init_from_model(model, contextParams)
	if err != nil {
		gollama.Model_free(model)
		return types.WrapError(types.ErrorCodeModelLoad, "failed to create context", err)
	}

	c.model = model
	c.ctx = ctx
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

	if c.ctx != 0 {
		// Context will be freed with the model
		c.ctx = 0
	}

	if c.model != 0 {
		gollama.Model_free(c.model)
	}

	c.model = 0
	c.ctx = 0
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

	if !c.modelLoaded || c.model == 0 {
		return nil, types.ErrModelNotFound("no model loaded")
	}

	return &types.Model{
		Name:     c.modelPath,
		FilePath: c.modelPath,
		Status:   types.ModelStatusLoaded,
	}, nil
}

// buildSamplerChain creates a sampler chain based on inference options
func (c *LlamaClient) buildSamplerChain(options *InferenceOptions) gollama.LlamaSampler {
	// Initialize sampler chain with default params
	samplerParams := gollama.Sampler_chain_default_params()
	sampler := gollama.Sampler_chain_init(samplerParams)

	// Add samplers based on options in the correct order
	// The order matters - typically: penalties -> top_k/top_p -> temperature -> greedy

	// 1. Add repeat penalty
	if options.RepeatPenalty > 1.0 {
		// Add repeat penalty sampler using internal function
	}

	// 2. Add frequency penalty
	if options.FrequencyPenaltyEnabled && options.FrequencyPenalty > 0 {
		// Add frequency penalty sampler
	}

	// 3. Add presence penalty
	if options.PresencePenaltyEnabled && options.PresencePenalty > 0 {
		// Add presence penalty sampler
	}

	// 4. Add top_k if enabled
	if options.TopK > 0 {
		topKSampler := gollama.Sampler_init_top_k(int32(options.TopK))
		gollama.Sampler_chain_add(sampler, topKSampler)
	}

	// 5. Add top_p if enabled
	if options.TopPEnabled && options.TopP > 0 && options.TopP < 1.0 {
		topPSampler := gollama.Sampler_init_top_p(float32(options.TopP), 1)
		gollama.Sampler_chain_add(sampler, topPSampler)
	}

	// 6. Add min_p
	if options.MinP > 0 {
		minPSampler := gollama.Sampler_init_min_p(float32(options.MinP), 1)
		gollama.Sampler_chain_add(sampler, minPSampler)
	}

	// 7. Add typical_p if enabled
	if options.TypicalPEnabled && options.TypicalP > 0 && options.TypicalP < 1.0 {
		typicalSampler := gollama.Sampler_init_typical(float32(options.TypicalP), 1)
		gollama.Sampler_chain_add(sampler, typicalSampler)
	}

	// 8. Add temperature if enabled and not 1.0
	if options.Temperature > 0 && options.Temperature != 1.0 {
		tempSampler := gollama.Sampler_init_temp(float32(options.Temperature))
		gollama.Sampler_chain_add(sampler, tempSampler)
	}

	// 9. Add mirostat if enabled
	if options.MirostatEnabled && options.Mirostat > 0 {
		if options.Mirostat == 1 {
			mirostatSampler := gollama.Sampler_init_mirostat(float32(options.MirostatTau), float32(options.MirostatETA), int32(options.Mirostat), uint32(options.Seed))
			gollama.Sampler_chain_add(sampler, mirostatSampler)
		} else if options.Mirostat == 2 {
			mirostatSampler := gollama.Sampler_init_mirostat_v2(float32(options.MirostatTau), float32(options.MirostatETA), uint32(options.Seed))
			gollama.Sampler_chain_add(sampler, mirostatSampler)
		}
	}

	// 10. Add greedy sampler at the end
	greedySampler := gollama.Sampler_init_greedy()
	gollama.Sampler_chain_add(sampler, greedySampler)

	return sampler
}

// Generate generates text using the loaded model
func (c *LlamaClient) Generate(prompt string, options *InferenceOptions) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.modelLoaded || c.ctx == 0 {
		return "", types.NewError(types.ErrorCodeModelNotFound, "no model loaded")
	}

	if options == nil {
		options = NewInferenceOptions()
	}

	// Tokenize the prompt
	tokens, err := gollama.Tokenize(c.model, prompt, true, false)
	if err != nil {
		return "", types.WrapError(types.ErrorCodeContext, "tokenization failed", err)
	}

	// Create batch for the prompt tokens
	batch := gollama.Batch_get_one(tokens)
	defer gollama.Batch_free(batch)

	// Process the prompt
	if err := gollama.Decode(c.ctx, batch); err != nil {
		return "", types.WrapError(types.ErrorCodeInference, "failed to decode prompt", err)
	}

	// Build sampler chain from options
	sampler := c.buildSamplerChain(options)
	defer gollama.Sampler_free(sampler)

	// Generate tokens
	var result string
	nCur := len(tokens)
	maxTokens := options.NPredict
	if maxTokens <= 0 {
		maxTokens = 512 // default
	}

	// Get context size from model metadata or use a reasonable default
	maxContext := 2048
	if c.modelProfile != nil && c.modelProfile.UsableContext > 0 {
		maxContext = c.modelProfile.UsableContext
	}

	for i := 0; i < maxTokens && nCur < maxContext; i++ {
		// Sample next token from the context
		newToken := gollama.Sampler_sample(sampler, c.ctx, -1)

		// Check for EOS token
		if newToken == 0 || newToken == 2 || newToken == 3 {
			break
		}

		// Check stop strings from options
		if c.shouldStop(result, options.StopStrings) {
			break
		}

		// Check stop strings from model metadata
		if c.modelProfile != nil && len(c.modelProfile.StopSequences) > 0 {
			if c.shouldStop(result, c.modelProfile.StopSequences) {
				break
			}
		}

		// Convert token to text
		piece := gollama.Token_to_piece(c.model, newToken, false)
		result += piece

		// Create a new batch with the single token
		batch = gollama.Batch_get_one([]gollama.LlamaToken{newToken})

		// Decode the new token
		if err := gollama.Decode(c.ctx, batch); err != nil {
			types.WrapError(types.ErrorCodeInference, "failed to decode token", err)
			break
		}

		nCur++
	}

	return result, nil
}

// shouldStop checks if generation should stop based on stop strings
func (c *LlamaClient) shouldStop(text string, stopStrings []string) bool {
	for _, stopStr := range stopStrings {
		if len(text) >= len(stopStr) && text[len(text)-len(stopStr):] == stopStr {
			return true
		}
	}
	return false
}

// GenerateStream generates text with streaming
func (c *LlamaClient) GenerateStream(prompt string, options *InferenceOptions, callback TokenCallback) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.modelLoaded || c.ctx == 0 {
		return types.ErrModelNotFound("no model loaded")
	}

	if options == nil {
		options = NewInferenceOptions()
	}

	// Tokenize the prompt
	tokens, err := gollama.Tokenize(c.model, prompt, true, false)
	if err != nil {
		return types.WrapError(types.ErrorCodeContext, "tokenization failed", err)
	}

	// Create batch for the prompt tokens
	batch := gollama.Batch_get_one(tokens)
	defer gollama.Batch_free(batch)

	// Process the prompt
	if err := gollama.Decode(c.ctx, batch); err != nil {
		return types.WrapError(types.ErrorCodeInference, "failed to decode prompt", err)
	}

	// Build sampler chain from options
	sampler := c.buildSamplerChain(options)
	defer gollama.Sampler_free(sampler)

	// Generate tokens
	nCur := len(tokens)
	maxTokens := options.NPredict
	if maxTokens <= 0 {
		maxTokens = 512 // default
	}

	// Get context size from model metadata or use a reasonable default
	maxContext := 2048
	if c.modelProfile != nil && c.modelProfile.UsableContext > 0 {
		maxContext = c.modelProfile.UsableContext
	}

	for i := 0; i < maxTokens && nCur < maxContext; i++ {
		// Sample next token from the context
		newToken := gollama.Sampler_sample(sampler, c.ctx, -1)

		// Check for EOS token
		if newToken == 0 || newToken == 2 || newToken == 3 {
			break
		}

		// Convert token to text
		piece := gollama.Token_to_piece(c.model, newToken, false)

		// Send token to callback
		callback(piece, false)

		// Create a new batch with the single token
		batch = gollama.Batch_get_one([]gollama.LlamaToken{newToken})

		// Decode the new token
		if err := gollama.Decode(c.ctx, batch); err != nil {
			types.WrapError(types.ErrorCodeInference, "failed to decode token", err)
			break
		}

		nCur++
	}

	// Signal completion
	callback("", true)
	return nil
}

// Tokenize converts text to token IDs
func (c *LlamaClient) Tokenize(text string, options *TokenizerOptions) ([]int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.model == 0 {
		return nil, types.ErrModelNotFound("no model loaded")
	}

	if options == nil {
		options = NewTokenizerOptions()
	}

	tokens, err := gollama.Tokenize(c.model, text, options.AddBOS, options.AddEOS)
	if err != nil {
		return nil, types.WrapError(types.ErrorCodeContext, "tokenization failed", err)
	}

	// Convert []gollama.LlamaToken to []int
	result := make([]int, len(tokens))
	for i, t := range tokens {
		result[i] = int(t)
	}

	return result, nil
}

// Detokenize converts token IDs back to text
func (c *LlamaClient) Detokenize(tokens []int) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.model == 0 {
		return "", types.ErrModelNotFound("no model loaded")
	}

	// Convert []int to []gollama.LlamaToken
	llamaTokens := make([]gollama.LlamaToken, len(tokens))
	for i, t := range tokens {
		llamaTokens[i] = gollama.LlamaToken(t)
	}

	var result string
	for _, token := range llamaTokens {
		piece := gollama.Token_to_piece(c.model, token, false)
		result += piece
	}

	return result, nil
}

// GetContextSize returns the context window size
func (c *LlamaClient) GetContextSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.ctx == 0 {
		return 0
	}

	// Placeholder - gollama doesn't have a direct context size getter
	// Would need to use Context_default_params or query the context
	return 2048 // Default context size
}

// GetUsedMemory returns memory usage
func (c *LlamaClient) GetUsedMemory() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.ctx == 0 {
		return 0
	}

	// Placeholder - gollama doesn't have direct memory getters
	return 0
}

// GetTotalMemory returns total memory
func (c *LlamaClient) GetTotalMemory() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.ctx == 0 {
		return 0
	}

	// Placeholder - gollama doesn't have direct memory getters
	return 0
}

// Embedding generates embeddings
func (c *LlamaClient) Embedding(text string) ([]float32, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modelLoaded || c.model == 0 {
		return nil, types.ErrModelNotFound("no model loaded")
	}

	// Placeholder - gollama doesn't have a direct embedding function
	// Would need to implement embedding generation manually
	return nil, types.NewError(types.ErrorCodeInference, "embedding generation not yet implemented with gollama.cpp")
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

	if c.ctx != 0 {
		// Context will be freed with the model
		c.ctx = 0
	}

	if c.model != 0 {
		gollama.Model_free(c.model)
	}

	c.model = 0
	c.ctx = 0
	c.modelLoaded = false
	c.modelPath = ""

	return nil
}
