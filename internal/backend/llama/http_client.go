package llama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"howl-chat/internal/backend/types"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPClient handles communication with llama.cpp server via HTTP API
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	enabled    bool
}

// NewHTTPClient creates a new llama.cpp HTTP client
func NewHTTPClient(host string, port int, enabled bool) *HTTPClient {
	if !enabled {
		return &HTTPClient{
			enabled: false,
		}
	}

	return &HTTPClient{
		baseURL: fmt.Sprintf("http://%s:%d", host, port),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		enabled: true,
	}
}

// ChatMessage represents a message in the chat
type ChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // Can be string or array for vision
}

// ImageContent represents OpenAI vision format
type ImageContent struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// ImageURL represents image URL in OpenAI format
type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// ChatCompletionRequest represents a chat completion request to llama.cpp (OpenAI-compatible)
type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
	Stop     []string      `json:"stop,omitempty"`

	// Standard OpenAI fields
	MaxTokens        int     `json:"max_tokens,omitempty"`
	Temperature      float64 `json:"temperature,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	Seed             *int    `json:"seed,omitempty"`

	// llama.cpp extensions
	TopK              int     `json:"top_k,omitempty"`
	MinP              float64 `json:"min_p,omitempty"`
	TypicalP          float64 `json:"typical_p,omitempty"`
	RepeatPenalty     float64 `json:"repeat_penalty,omitempty"`
	RepeatLastN       int     `json:"repeat_last_n,omitempty"`
	Mirostat          int     `json:"mirostat,omitempty"`
	MirostatTau       float64 `json:"mirostat_tau,omitempty"`
	MirostatEta       float64 `json:"mirostat_eta,omitempty"`
	DynTempRange      float64 `json:"dynatemp_range,omitempty"`
	DynTempExponent   float64 `json:"dynatemp_exponent,omitempty"`
	DryMultiplier     float64 `json:"dry_multiplier,omitempty"`
	DryAllowedLength  int     `json:"dry_allowed_length,omitempty"`
	DryBase           float64 `json:"dry_base,omitempty"`
	SmoothingFactor   float64 `json:"smoothing_factor,omitempty"`
	SmoothingCurve    float64 `json:"smoothing_curve,omitempty"`
	TopA              float64 `json:"top_a,omitempty"`
	EpsilonCutoff     float64 `json:"epsilon_cutoff,omitempty"`
	EtaCutoff         float64 `json:"eta_cutoff,omitempty"`
	NoRepeatNgramSize int     `json:"no_repeat_ngram_size,omitempty"`
}

// ChatCompletionResponse represents a chat completion response from llama.cpp (OpenAI-compatible)
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Health checks if llama.cpp server is accessible
func (c *HTTPClient) Health() error {
	if !c.enabled {
		return fmt.Errorf("llama.cpp HTTP client is disabled")
	}

	url := fmt.Sprintf("%s/health", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to llama.cpp: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("llama.cpp health check failed: %s", resp.Status)
	}

	return nil
}

// StreamCallback is called for each chunk in a streaming response
type StreamCallback func(chunk string)

// GenerateChatCompletionStreaming sends a streaming chat completion request to llama.cpp (OpenAI-compatible)
func (c *HTTPClient) GenerateChatCompletionStreaming(req *ChatCompletionRequest, callback StreamCallback) error {
	if !c.enabled {
		return fmt.Errorf("llama.cpp HTTP client is disabled")
	}

	// Enable streaming
	req.Stream = true

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/chat/completions", c.baseURL)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("llama.cpp request failed: %s - %s", resp.Status, string(body))
	}

	// Handle streaming response
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line == "[DONE]" {
			continue
		}
		if strings.HasPrefix(line, "data: ") {
			jsonStr := strings.TrimPrefix(line, "data: ")
			if jsonStr == "[DONE]" {
				break
			}
			var chunk struct {
				Choices []struct {
					Delta struct {
						Content          string `json:"content"`
						ReasoningContent string `json:"reasoning_content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(jsonStr), &chunk); err == nil {
				if len(chunk.Choices) > 0 {
					content := chunk.Choices[0].Delta.Content
					reasoning := chunk.Choices[0].Delta.ReasoningContent
					if content != "" {
						callback(content)
					}
					if reasoning != "" {
						callback(reasoning)
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// LoadModel loads a model into llama.cpp (implements Client interface)
func (c *HTTPClient) LoadModel(modelPath string, options *LoadOptions) error {
	if !c.enabled {
		return fmt.Errorf("llama.cpp HTTP client is not enabled")
	}

	// Use llama.cpp API to load the model
	loadURL := fmt.Sprintf("%s/models/load", c.baseURL)

	requestBody := map[string]interface{}{
		"model": modelPath, // full absolute path
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal load request: %w", err)
	}

	resp, err := c.httpClient.Post(loadURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to send load request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to load model: %s", string(body))
	}

	return nil
}

// LoadModelWithMMProj loads a model with optional mmproj file for vision models
func (c *HTTPClient) LoadModelWithMMProj(modelPath string, mmprojPath string) error {
	// For HTTP client, this is the same as LoadModel since llama-server handles mmproj auto-discovery
	return c.LoadModel(modelPath, nil)
}

// GetLoadedModel returns information about the currently loaded model
func (c *HTTPClient) GetLoadedModel() (map[string]interface{}, error) {
	if !c.enabled {
		return nil, fmt.Errorf("llama.cpp HTTP client is disabled")
	}

	url := fmt.Sprintf("%s/v1/models", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get loaded model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("llama.cpp get model failed: %s - %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// Close releases all resources associated with the client
func (c *HTTPClient) Close() error {
	// HTTP client doesn't need explicit cleanup
	return nil
}

// UnloadModel unloads the currently loaded model (implements Client interface)
func (c *HTTPClient) UnloadModel() error {
	// llama-server doesn't have an unload API in router mode
	return nil
}

// IsLoaded returns true if a model is currently loaded (implements Client interface)
func (c *HTTPClient) IsLoaded() bool {
	result, err := c.GetLoadedModel()
	if err != nil || result == nil {
		return false
	}
	data, ok := result["data"].([]interface{})
	return ok && len(data) > 0
}

// GetModelInfo returns information about the currently loaded model (implements Client interface)
func (c *HTTPClient) GetModelInfo() (*types.Model, error) {
	// This would need to be implemented based on the HTTP API
	return nil, fmt.Errorf("GetModelInfo not implemented for HTTP client")
}

// Generate generates text completion for the given prompt (implements Client interface)
func (c *HTTPClient) Generate(prompt string, options *InferenceOptions) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("llama.cpp HTTP client is disabled")
	}

	req := buildRequest("model", prompt, options, false)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/chat/completions", c.baseURL)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("llama-server error: %s - %s", resp.Status, string(body))
	}

	var result ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	return result.Choices[0].Message.Content, nil
}

// buildRequest constructs a ChatCompletionRequest from InferenceOptions
func buildRequest(model string, prompt string, options *InferenceOptions, stream bool) *ChatCompletionRequest {
	req := &ChatCompletionRequest{
		Model:    model,
		Messages: []ChatMessage{{Role: "user", Content: prompt}},
		Stream:   stream,

		MaxTokens:         options.NPredict,
		Temperature:       options.Temperature,
		TopK:              options.TopK,
		MinP:              options.MinP,
		RepeatPenalty:     options.RepeatPenalty,
		RepeatLastN:       options.RepeatLastN,
		Mirostat:          options.Mirostat,
		MirostatTau:       options.MirostatTau,
		MirostatEta:       options.MirostatETA,
		DynTempRange:      options.DynamicTempRange,
		DynTempExponent:   options.DynamicTempExponent,
		DryMultiplier:     options.DRYMultiplier,
		DryAllowedLength:  options.DRYAllowedLength,
		DryBase:           options.DRYBase,
		SmoothingFactor:   options.SmoothingFactor,
		SmoothingCurve:    options.SmoothingCurve,
		TopA:              options.TopA,
		EpsilonCutoff:     options.EpsilonCutoff,
		EtaCutoff:         options.EtaCutoff,
		NoRepeatNgramSize: options.NoRepeatNGramSize,
	}

	// Conditionally enabled fields
	if options.TopPEnabled {
		req.TopP = options.TopP
	}
	if options.TypicalPEnabled {
		req.TypicalP = options.TypicalP
	}
	if options.FrequencyPenaltyEnabled {
		req.FrequencyPenalty = options.FrequencyPenalty
	}
	if options.PresencePenaltyEnabled {
		req.PresencePenalty = options.PresencePenalty
	}
	if options.Seed >= 0 {
		s := options.Seed
		req.Seed = &s
	}
	if len(options.StopStrings) > 0 {
		req.Stop = options.StopStrings
	}

	return req
}

// GenerateStream generates text completion with streaming output (implements Client interface)
func (c *HTTPClient) GenerateStream(prompt string, options *InferenceOptions, callback TokenCallback) error {
	req := buildRequest("model", prompt, options, true)
	return c.GenerateChatCompletionStreaming(req, func(chunk string) {
		callback(chunk, false)
	})
}

// Tokenize converts text to token IDs (implements Client interface)
func (c *HTTPClient) Tokenize(text string, options *TokenizerOptions) ([]int, error) {
	// This would need to be implemented via HTTP API
	return nil, fmt.Errorf("tokenize not implemented for HTTP client")
}

// Detokenize converts token IDs back to text (implements Client interface)
func (c *HTTPClient) Detokenize(tokens []int) (string, error) {
	// This would need to be implemented via HTTP API
	return "", fmt.Errorf("detokenize not implemented for HTTP client")
}

// GetContextSize returns the maximum context window size (implements Client interface)
func (c *HTTPClient) GetContextSize() int {
	// This would need to be implemented via HTTP API
	return 0
}

// GetUsedMemory returns the amount of memory used by the model (implements Client interface)
func (c *HTTPClient) GetUsedMemory() int64 {
	// This would need to be implemented via HTTP API
	return 0
}

// GetTotalMemory returns the total memory available (implements Client interface)
func (c *HTTPClient) GetTotalMemory() int64 {
	// This would need to be implemented via HTTP API
	return 0
}

// Embedding generates embeddings for the given text (implements Client interface)
func (c *HTTPClient) Embedding(text string) ([]float32, error) {
	// This would need to be implemented via HTTP API
	return nil, fmt.Errorf("embedding not implemented for HTTP client")
}

// SetProgressCallback sets a callback for loading progress (implements Client interface)
func (c *HTTPClient) SetProgressCallback(callback ProgressCallback) {
	// Not needed for HTTP client
}
