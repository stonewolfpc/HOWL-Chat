package qwen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"howl-chat/internal/audio/asr"
	"howl-chat/internal/audio/types"
)

// Client implements ASR recognizer for Qwen-ASR HTTP API
type Client struct {
	asr.BaseRecognizer
	apiKey     string
	endpoint   string
	httpClient *http.Client
	model      string
}

// NewClient creates a new Qwen-ASR HTTP client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:   apiKey,
		endpoint: "https://dashscope-intl.aliyuncs.com/api/v1/services/audio/asr/transcription",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		model: "qwen3-asr-flash-filetrans",
	}
}

// Initialize sets up the client configuration
func (c *Client) Initialize(ctx context.Context, config types.ASRConfig) error {
	// Validate API key
	if c.apiKey == "" {
		return types.NewAudioError(
			asr.ErrCodeInvalidConfig,
			"Qwen-ASR API key is required",
			nil,
		)
	}

	// Set model from config if specified
	if config.ModelPath != "" {
		c.model = config.ModelPath
	}

	c.SetInitialized(config)
	return nil
}

// Transcribe performs speech recognition using Qwen-ASR HTTP API
func (c *Client) Transcribe(ctx context.Context, audioPath string, callback types.ProgressCallback) (*types.RecognitionResult, error) {
	// Ensure initialized
	if err := c.CheckInitialized(); err != nil {
		return nil, err
	}

	// Report progress
	if callback != nil {
		callback(10, "Uploading audio...")
	}

	// Submit transcription task
	taskID, err := c.submitTask(ctx, audioPath)
	if err != nil {
		return nil, err
	}

	// Report progress
	if callback != nil {
		callback(30, "Processing transcription...")
	}

	// Poll for result
	result, err := c.pollForResult(ctx, taskID, callback)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// TranscribeStream is not supported by Qwen-ASR HTTP API
func (c *Client) TranscribeStream(ctx context.Context, chunks <-chan types.AudioChunk, results chan<- types.RecognitionResult) error {
	return types.NewAudioError(
		types.ErrCodeFormatNotSupported,
		"Streaming transcription not supported by Qwen-ASR HTTP API",
		nil,
	)
}

// Release cleans up resources
func (c *Client) Release() error {
	// No resources to clean up for HTTP client
	return nil
}

// submitTask submits an audio file for transcription
func (c *Client) submitTask(ctx context.Context, audioPath string) (string, error) {
	// Build request body
	reqBody := map[string]interface{}{
		"model": c.model,
		"input": map[string]string{
			"file_url": audioPath,
		},
		"parameters": map[string]interface{}{
			"channel_id":   []int{0},
			"enable_itn":   false,
			"enable_words": true,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-Async", "enable")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to submit task: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", types.NewAudioError(
			asr.ErrCodeAPIError,
			fmt.Sprintf("Qwen-ASR API error: %s", string(body)),
			nil,
		)
	}

	// Parse response
	var submitResp submitResponse
	if err := json.NewDecoder(resp.Body).Decode(&submitResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if submitResp.Output.TaskID == "" {
		return "", types.NewAudioError(
			asr.ErrCodeAPIError,
			"No task ID returned from Qwen-ASR API",
			nil,
		)
	}

	return submitResp.Output.TaskID, nil
}

// pollForResult polls the API for transcription result
func (c *Client) pollForResult(ctx context.Context, taskID string, callback types.ProgressCallback) (*types.RecognitionResult, error) {
	pollURL := fmt.Sprintf("https://dashscope-intl.aliyuncs.com/api/v1/tasks/%s", taskID)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return nil, types.NewAudioError(
				types.ErrCodeTimeout,
				"Transcription cancelled by user",
				ctx.Err(),
			)
		case <-timeout:
			return nil, types.NewAudioError(
				types.ErrCodeTimeout,
				"Transcription timed out after 5 minutes",
				nil,
			)
		case <-ticker.C:
			// Poll for result
			result, err := c.getTaskResult(ctx, pollURL)
			if err != nil {
				return nil, err
			}

			// Update progress
			if callback != nil {
				progress := 30 + (70 * 5 / 100) // Simple progress estimate
				callback(progress, "Processing...")
			}

			// Check if task is complete
			if result.Output.TaskStatus == "SUCCEEDED" {
				return c.parseTranscriptionResult(result)
			}

			if result.Output.TaskStatus == "FAILED" {
				return nil, types.NewAudioError(
					types.ErrCodeModelLoadFailed,
					fmt.Sprintf("Qwen-ASR task failed: %s", result.Output.Message),
					nil,
				)
			}
		}
	}
}

// getTaskResult retrieves the task result from API
func (c *Client) getTaskResult(ctx context.Context, url string) (*taskResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-Async", "enable")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get task result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, types.NewAudioError(
			asr.ErrCodeAPIError,
			fmt.Sprintf("Qwen-ASR API error: %s", string(body)),
			nil,
		)
	}

	var result taskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// parseTranscriptionResult converts Qwen-ASR response to RecognitionResult
func (c *Client) parseTranscriptionResult(result *taskResponse) (*types.RecognitionResult, error) {
	// If transcription_url is provided, fetch and parse the full result
	if result.Output.TranscriptionURL != "" {
		return c.fetchAndParseTranscription(result.Output.TranscriptionURL)
	}

	// Otherwise, try to parse from result field (synchronous mode)
	return c.parseInlineResult(result.Output.Result)
}

// fetchAndParseTranscription fetches the full transcription from URL and parses it
func (c *Client) fetchAndParseTranscription(url string) (*types.RecognitionResult, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transcription: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("transcription fetch failed with status %d", resp.StatusCode)
	}

	var transResult transcriptionResult
	if err := json.NewDecoder(resp.Body).Decode(&transResult); err != nil {
		return nil, fmt.Errorf("failed to decode transcription: %w", err)
	}

	return c.convertToRecognitionResult(&transResult)
}

// parseInlineResult parses transcription from inline result field
func (c *Client) parseInlineResult(result string) (*types.RecognitionResult, error) {
	if result == "" {
		return nil, types.NewAudioError(
			asr.ErrCodeEmptyResult,
			"No transcription result provided",
			nil,
		)
	}

	// Try to parse as JSON first
	var transResult transcriptionResult
	if err := json.Unmarshal([]byte(result), &transResult); err == nil {
		return c.convertToRecognitionResult(&transResult)
	}

	// If not JSON, treat as plain text
	return &types.RecognitionResult{
		Text:       result,
		Language:   "en",
		Confidence: 0.8,
	}, nil
}

// convertToRecognitionResult converts Qwen-ASR result to RecognitionResult
func (c *Client) convertToRecognitionResult(result *transcriptionResult) (*types.RecognitionResult, error) {
	if len(result.Transcripts) == 0 {
		return nil, asr.NewEmptyResultError()
	}

	// Use first transcript (primary channel)
	transcript := result.Transcripts[0]

	recognitionResult := &types.RecognitionResult{
		Text:       transcript.Text,
		Language:   "en", // Will be detected from sentences if available
		Confidence: 0.8,  // Qwen-ASR doesn't provide confidence scores
	}

	// Extract segments from sentences
	for _, sent := range transcript.Sentences {
		segment := types.Segment{
			StartTime:  float64(sent.BeginTime) / 1000.0, // Convert ms to seconds
			EndTime:    float64(sent.EndTime) / 1000.0,
			Text:       sent.Text,
			Confidence: 0.8,
		}
		recognitionResult.Segments = append(recognitionResult.Segments, segment)

		// Detect language from first sentence
		if sent.Language != "" && recognitionResult.Language == "en" {
			recognitionResult.Language = sent.Language
		}
	}

	// Calculate duration from last segment
	if len(recognitionResult.Segments) > 0 {
		recognitionResult.Duration = recognitionResult.Segments[len(recognitionResult.Segments)-1].EndTime
	}

	return recognitionResult, nil
}

// submitResponse represents the response from task submission
type submitResponse struct {
	Output struct {
		TaskID string `json:"task_id"`
	} `json:"output"`
}

// taskResponse represents the response from task status check
type taskResponse struct {
	Output struct {
		TaskStatus       string `json:"task_status"`
		Message          string `json:"message"`
		Result           string `json:"result"`
		TranscriptionURL string `json:"transcription_url"`
	} `json:"output"`
}

// transcriptionResult represents the full transcription result from Qwen-ASR
type transcriptionResult struct {
	FileURL     string       `json:"file_url"`
	AudioInfo   audioInfo    `json:"audio_info"`
	Transcripts []transcript `json:"transcripts"`
}

type audioInfo struct {
	Format     string `json:"format"`
	SampleRate int    `json:"sample_rate"`
}

type transcript struct {
	ChannelID int        `json:"channel_id"`
	Text      string     `json:"text"`
	Sentences []sentence `json:"sentences"`
}

type sentence struct {
	SentenceID int    `json:"sentence_id"`
	BeginTime  int    `json:"begin_time"`
	EndTime    int    `json:"end_time"`
	Language   string `json:"language"`
	Emotion    string `json:"emotion"`
	Text       string `json:"text"`
	Words      []word `json:"words,omitempty"`
}

type word struct {
	BeginTime   int    `json:"begin_time"`
	EndTime     int    `json:"end_time"`
	Text        string `json:"text"`
	Punctuation string `json:"punctuation"`
}
