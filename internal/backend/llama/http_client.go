package llama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"howl-chat/internal/backend/chatpolicy"
	"howl-chat/internal/backend/gguf"
	"howl-chat/internal/backend/types"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// HTTPClient handles communication with llama.cpp server via HTTP API
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	enabled    bool
	currentReq *http.Request
	cancelReq  context.CancelFunc
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
			Timeout: 300 * time.Second, // 5 minutes for vision/image processing
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
	TopK             int     `json:"top_k,omitempty"`
	MinP             float64 `json:"min_p,omitempty"`
	TypicalP         float64 `json:"typical_p,omitempty"`
	RepeatPenalty    float64 `json:"repeat_penalty,omitempty"`
	RepeatLastN      int     `json:"repeat_last_n,omitempty"`
	Mirostat         int     `json:"mirostat,omitempty"`
	MirostatTau      float64 `json:"mirostat_tau,omitempty"`
	MirostatEta      float64 `json:"mirostat_eta,omitempty"`
	DynTempRange     float64 `json:"dynatemp_range,omitempty"`
	DynTempExponent  float64 `json:"dynatemp_exponent,omitempty"`
	DryMultiplier    float64 `json:"dry_multiplier,omitempty"`
	DryAllowedLength int     `json:"dry_allowed_length,omitempty"`
	DryBase          float64 `json:"dry_base,omitempty"`

	// OnReasoning receives reasoning/thinking deltas separately from content.
	// Never marshalled: a local delivery hook, not a server field.
	OnReasoning func(string) `json:"-"`
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

	// Create request with context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	c.cancelReq = cancel

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	c.currentReq = httpReq

	resp, err := c.httpClient.Do(httpReq)
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
	chunkCount := 0
	totalContent := 0
	reasoningChars := 0
	reasoningChunks := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line == "[DONE]" {
			continue
		}
		if strings.HasPrefix(line, "data: ") {
			jsonStr := strings.TrimPrefix(line, "data: ")
			if jsonStr == "[DONE]" {
				fmt.Printf("DEBUG: Stream received [DONE] signal, chunks=%d, totalContent=%d\n", chunkCount, totalContent)
				break
			}
			// Log raw chunk for debugging vision/multimodal responses
			fmt.Printf("DEBUG: raw chunk: %s\n", jsonStr)

			var chunk struct {
				Choices []struct {
					Delta struct {
						Role             string `json:"role"`
						Content          string `json:"content"`
						ReasoningContent string `json:"reasoning_content"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
					Message      struct {
						Content string `json:"content"`
					} `json:"message"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(jsonStr), &chunk); err == nil {
				if len(chunk.Choices) > 0 {
					choice := chunk.Choices[0]
					content := choice.Delta.Content
					reasoning := choice.Delta.ReasoningContent
					role := choice.Delta.Role
					finish := choice.FinishReason
					msgContent := choice.Message.Content

					// Log all fields for diagnostics
					if role != "" || finish != "" || msgContent != "" {
						fmt.Printf("DEBUG: choice fields role=%q finish=%q msgContent=%q\n", role, finish, msgContent)
					}

					// Prefer delta.content, fallback to message.content
					if content == "" && msgContent != "" {
						content = msgContent
						fmt.Printf("DEBUG: using message.content fallback: %q\n", content)
					}

					if content != "" {
						callback(content)
						chunkCount++
						totalContent += len(content)
						if chunkCount <= 5 || chunkCount%50 == 0 {
							fmt.Printf("DEBUG: Stream chunk %d, content len=%d: %q\n", chunkCount, len(content), content[:min(50, len(content))])
						}
					}
					if reasoning != "" {
						// Reasoning must never travel through the content callback:
						// whatever lands there becomes the assistant reply, is saved
						// to history, and is fed back to the model next turn. Count it
						// for diagnostics only.
						reasoningChars += len(reasoning)
						reasoningChunks++
						if req.OnReasoning != nil {
							req.OnReasoning(reasoning)
						}
					}
				}
			} else {
				fmt.Printf("DEBUG: Failed to unmarshal chunk: %v\n", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}
	fmt.Printf("DEBUG: Stream completed, total chunks=%d, total content=%d bytes\n", chunkCount, totalContent)
	fmt.Printf("DEBUG: Stream reasoning: %d chunks, %d chars (never fed back)\n", reasoningChunks, reasoningChars)

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// LoadModel loads a model into llama.cpp (implements Client interface)
func (c *HTTPClient) LoadModel(modelPath string, options *LoadOptions) error {
	return c.loadModelWithProjector(modelPath, "")
}

// loadModelWithProjector loads a model and, when the model requires one,
// passes the exact companion projector. It never searches the directory for
// the first mmproj: directory order is not evidence of pairing.
func (c *HTTPClient) loadModelWithProjector(modelPath, explicitMMProj string) error {
	if !c.enabled {
		return fmt.Errorf("llama.cpp HTTP client is not enabled")
	}

	// Use llama.cpp API to load the model
	loadURL := fmt.Sprintf("%s/models/load", c.baseURL)

	requestBody := map[string]interface{}{
		"model": modelPath, // full absolute path
	}

	admission := gguf.AdmitModel(modelPath)
	if !admission.AllowedAsChat() {
		return fmt.Errorf("%s", admission.ChatRejectionMessage())
	}
	eligibility := chatpolicy.Evaluate(modelPath, admission)
	if !eligibility.Eligible {
		return fmt.Errorf("%s", eligibility.RejectionMessage(modelPath))
	}

	mmprojPath := explicitMMProj
	if mmprojPath == "" {
		pairing, err := gguf.ResolveProjector(gguf.PairingRequest{
			ModelPath:        modelPath,
			Arch:             admission.Arch,
			EmbeddingLength:  admission.EmbeddingLength,
			Require:          admission.RequiresProjector(),
			RequiredModality: admission.RequiredModality(),
		})
		if err != nil {
			return fmt.Errorf("cannot pair a companion projector for %s: %w", filepath.Base(modelPath), err)
		}
		mmprojPath = pairing.Path
	}
	if mmprojPath != "" {
		requestBody["mmproj"] = mmprojPath
		fmt.Printf("INFO: using paired mmproj file: %s\n", mmprojPath)
	} else {
		fmt.Printf("INFO: no paired mmproj file for model: %s\n", modelPath)
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal load request: %w", err)
	}

	fmt.Printf("DEBUG: Load request body: %s\n", string(jsonBody))

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

// LoadModelWithMMProj loads a model with an explicitly paired mmproj file.
func (c *HTTPClient) LoadModelWithMMProj(modelPath string, mmprojPath string) error {
	return c.loadModelWithProjector(modelPath, mmprojPath)
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

// CancelRequest cancels the current ongoing request
func (c *HTTPClient) CancelRequest() {
	if c.cancelReq != nil {
		c.cancelReq()
		c.cancelReq = nil
	}
	c.currentReq = nil
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

// GetModelInfo returns information about the currently loaded model (implements Client interface).
// llama-server exposes the loaded model through /props (model_path, n_ctx) and /v1/models (id),
// so the info is read from the server rather than the local model registry.
func (c *HTTPClient) GetModelInfo() (*types.Model, error) {
	props, err := c.serverProps()
	if err != nil {
		return nil, err
	}

	modelPath, _ := props["model_path"].(string)
	if modelPath == "" {
		if loaded, loadErr := c.GetLoadedModel(); loadErr == nil {
			if data, ok := loaded["data"].([]interface{}); ok && len(data) > 0 {
				if first, ok := data[0].(map[string]interface{}); ok {
					modelPath, _ = first["id"].(string)
				}
			}
		}
	}
	if modelPath == "" {
		return nil, fmt.Errorf("llama.cpp server did not report a loaded model")
	}

	cleanPath := strings.ReplaceAll(modelPath, "\\", "/")
	format := types.FormatGGUF
	if strings.EqualFold(filepath.Ext(cleanPath), ".ggml") {
		format = types.FormatGGML
	}

	model := &types.Model{
		ID:       modelPath,
		Name:     filepath.Base(cleanPath),
		FilePath: modelPath,
		Format:   format,
		LoadedAt: time.Now().Format(time.RFC3339),
	}
	model.Parameters.ContextSize = contextSizeFromProps(props)
	return model, nil
}

// serverProps fetches /props from llama-server.
func (c *HTTPClient) serverProps() (map[string]interface{}, error) {
	if !c.enabled {
		return nil, fmt.Errorf("llama.cpp HTTP client is disabled")
	}

	resp, err := c.httpClient.Get(c.baseURL + "/props")
	if err != nil {
		return nil, fmt.Errorf("failed to get server props: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("llama.cpp props failed: %s - %s", resp.Status, string(body))
	}

	var props map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&props); err != nil {
		return nil, fmt.Errorf("failed to decode props: %w", err)
	}
	return props, nil
}

// contextSizeFromProps reads n_ctx, which current builds nest under
// default_generation_settings and older builds report at the top level.
func contextSizeFromProps(props map[string]interface{}) int {
	if gen, ok := props["default_generation_settings"].(map[string]interface{}); ok {
		if n, ok := gen["n_ctx"].(float64); ok && n > 0 {
			return int(n)
		}
	}
	if n, ok := props["n_ctx"].(float64); ok && n > 0 {
		return int(n)
	}
	return 0
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
	return buildRequestWithImage(model, prompt, options, stream, nil)
}

// buildRequestWithImage constructs a ChatCompletionRequest with optional image data
func buildRequestWithImage(model string, prompt string, options *InferenceOptions, stream bool, imageData []byte) *ChatCompletionRequest {
	var content interface{}
	if imageData != nil {
		// OpenAI vision format with image
		content = []ImageContent{
			{Type: "text", Text: prompt},
			{Type: "image_url", ImageURL: &ImageURL{
				URL:    fmt.Sprintf("data:image/jpeg;base64,%s", base64Encode(imageData)),
				Detail: "auto",
			}},
		}
	} else {
		content = prompt
	}

	req := &ChatCompletionRequest{
		Model:    model,
		Messages: []ChatMessage{{Role: "user", Content: content}},
		Stream:   stream,

		MaxTokens:        options.NPredict,
		Temperature:      options.Temperature,
		TopK:             options.TopK,
		MinP:             options.MinP,
		RepeatPenalty:    options.RepeatPenalty,
		RepeatLastN:      options.RepeatLastN,
		Mirostat:         options.Mirostat,
		MirostatTau:      options.MirostatTau,
		MirostatEta:      options.MirostatETA,
		DynTempRange:     options.DynamicTempRange,
		DynTempExponent:  options.DynamicTempExponent,
		DryMultiplier:    options.DRYMultiplier,
		DryAllowedLength: options.DRYAllowedLength,
		DryBase:          options.DRYBase,
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

// base64Encode encodes bytes to base64 string
func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// buildRequestFromMessages constructs a ChatCompletionRequest from a messages array.
func buildRequestFromMessages(messages []ChatMessage, options *InferenceOptions, stream bool) *ChatCompletionRequest {
	req := &ChatCompletionRequest{
		Model:    "model",
		Messages: messages,
		Stream:   stream,

		MaxTokens:        options.NPredict,
		Temperature:      options.Temperature,
		TopK:             options.TopK,
		MinP:             options.MinP,
		RepeatPenalty:    options.RepeatPenalty,
		RepeatLastN:      options.RepeatLastN,
		Mirostat:         options.Mirostat,
		MirostatTau:      options.MirostatTau,
		MirostatEta:      options.MirostatETA,
		DynTempRange:     options.DynamicTempRange,
		DynTempExponent:  options.DynamicTempExponent,
		DryMultiplier:    options.DRYMultiplier,
		DryAllowedLength: options.DRYAllowedLength,
		DryBase:          options.DRYBase,
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

// GenerateChatStream sends a proper messages array to the chat completions endpoint.
func (c *HTTPClient) GenerateChatStream(messages []ChatMessage, options *InferenceOptions, callback TokenCallback) error {
	req := buildRequestFromMessages(messages, options, true)
	req.OnReasoning = options.OnReasoning
	err := c.GenerateChatCompletionStreaming(req, func(chunk string) {
		callback(chunk, false)
	})
	if err != nil {
		return err
	}
	callback("", true)
	return nil
}

// GenerateChatStreamWithImage sends a proper messages array with optional image.
func (c *HTTPClient) GenerateChatStreamWithImage(messages []ChatMessage, imageData []byte, options *InferenceOptions, callback TokenCallback) error {
	req := buildRequestFromMessages(messages, options, true)
	req.OnReasoning = options.OnReasoning
	if imageData != nil {
		// Replace the last user message content with vision format
		for i := len(req.Messages) - 1; i >= 0; i-- {
			if req.Messages[i].Role == "user" {
				req.Messages[i].Content = []ImageContent{
					{Type: "text", Text: fmt.Sprintf("%v", req.Messages[i].Content)},
					{Type: "image_url", ImageURL: &ImageURL{
						URL:    fmt.Sprintf("data:image/jpeg;base64,%s", base64Encode(imageData)),
						Detail: "auto",
					}},
				}
				break
			}
		}
	}
	err := c.GenerateChatCompletionStreaming(req, func(chunk string) {
		callback(chunk, false)
	})
	if err != nil {
		return err
	}
	callback("", true)
	return nil
}

// GenerateChat generates text completion from a messages array (non-streaming).
func (c *HTTPClient) GenerateChat(messages []ChatMessage, options *InferenceOptions) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("llama.cpp HTTP client is disabled")
	}

	req := buildRequestFromMessages(messages, options, false)

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

// GenerateStreamWithImage generates text completion with streaming output and optional image
func (c *HTTPClient) GenerateStreamWithImage(prompt string, imageData []byte, options *InferenceOptions, callback TokenCallback) error {
	req := buildRequestWithImage("model", prompt, options, true, imageData)
	err := c.GenerateChatCompletionStreaming(req, func(chunk string) {
		callback(chunk, false)
	})
	if err != nil {
		return err
	}
	// Signal completion
	callback("", true)
	return nil
}

// GenerateStream generates text completion with streaming output (implements Client interface)
func (c *HTTPClient) GenerateStream(prompt string, options *InferenceOptions, callback TokenCallback) error {
	req := buildRequest("model", prompt, options, true)
	err := c.GenerateChatCompletionStreaming(req, func(chunk string) {
		callback(chunk, false)
	})
	if err != nil {
		return err
	}
	// Signal completion
	callback("", true)
	return nil
}

// Tokenize converts text to token IDs (implements Client interface).
// Uses llama-server's POST /tokenize; the server returns {"tokens": [...]}.
func (c *HTTPClient) Tokenize(text string, options *TokenizerOptions) ([]int, error) {
	if !c.enabled {
		return nil, fmt.Errorf("llama.cpp HTTP client is disabled")
	}

	payload := map[string]interface{}{"content": text}
	if options != nil && (options.AddBOS || options.AddEOS) {
		payload["add_special"] = true
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tokenize request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/tokenize", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("llama.cpp tokenize failed: %s - %s", resp.Status, string(raw))
	}

	var decoded struct {
		Tokens []interface{} `json:"tokens"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("failed to decode tokenize response: %w", err)
	}

	ids := make([]int, 0, len(decoded.Tokens))
	for _, token := range decoded.Tokens {
		switch v := token.(type) {
		case float64:
			ids = append(ids, int(v))
		case map[string]interface{}:
			if id, ok := v["id"].(float64); ok {
				ids = append(ids, int(id))
			}
		}
	}
	return ids, nil
}

// Detokenize converts token IDs back to text (implements Client interface).
// Uses llama-server's POST /detokenize; the server returns {"content": "..."}.
func (c *HTTPClient) Detokenize(tokens []int) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("llama.cpp HTTP client is disabled")
	}

	body, err := json.Marshal(map[string]interface{}{"tokens": tokens})
	if err != nil {
		return "", fmt.Errorf("failed to marshal detokenize request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/detokenize", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to detokenize: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("llama.cpp detokenize failed: %s - %s", resp.Status, string(raw))
	}

	var decoded struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", fmt.Errorf("failed to decode detokenize response: %w", err)
	}
	return decoded.Content, nil
}

// GetContextSize returns the maximum context window size (implements Client interface).
func (c *HTTPClient) GetContextSize() int {
	props, err := c.serverProps()
	if err != nil {
		return 0
	}
	return contextSizeFromProps(props)
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

// EmbeddingRequest represents an embedding request to llama.cpp (OpenAI-compatible)
type EmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// EmbeddingResponse represents an embedding response from llama.cpp
type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// Embedding generates embeddings for the given text (implements Client interface)
func (c *HTTPClient) Embedding(text string) ([]float32, error) {
	if !c.enabled {
		return nil, fmt.Errorf("llama.cpp HTTP client is disabled")
	}

	req := EmbeddingRequest{
		Model: "model",
		Input: text,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embedding request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/embeddings", c.baseURL)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send embedding request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("llama-server embedding error: %s - %s", resp.Status, string(body))
	}

	var result EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode embedding response: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding data in response")
	}
	return result.Data[0].Embedding, nil
}

// SetProgressCallback sets a callback for loading progress (implements Client interface)
func (c *HTTPClient) SetProgressCallback(callback ProgressCallback) {
	// Not needed for HTTP client
}
