package fallback

import (
	"context"
	"fmt"

	"howl-chat/internal/audio/asr"
	"howl-chat/internal/audio/hw"
	"howl-chat/internal/audio/types"
)

// ChainConfig configures the fallback chain behavior
type ChainConfig struct {
	// Recognizers in priority order (highest priority first)
	Recognizers []types.ASRRecognizer

	// Model-to-recognizer mapping for hardware-based selection
	RecognizerMap map[types.ASRModelType]types.ASRRecognizer

	// Max attempts per recognizer before falling back
	MaxAttempts int

	// Enable hardware-based auto-selection
	AutoSelectByHardware bool

	// Minimum memory tier to use
	MinMemoryTier types.MemoryTier

	// Preferred language (empty for auto-detect)
	Language string
}

// FallbackChain implements ASRRecognizer with fallback support
type FallbackChain struct {
	asr.BaseRecognizer
	config ChainConfig
}

// NewChain creates a new fallback chain
func NewChain(config ChainConfig) *FallbackChain {
	if config.MaxAttempts == 0 {
		config.MaxAttempts = 2
	}

	return &FallbackChain{
		config: config,
	}
}

// Initialize sets up the fallback chain
func (c *FallbackChain) Initialize(ctx context.Context, config types.ASRConfig) error {
	// Initialize all recognizers
	for _, recognizer := range c.config.Recognizers {
		if err := recognizer.Initialize(ctx, config); err != nil {
			return fmt.Errorf("failed to initialize recognizer: %w", err)
		}
	}

	c.SetInitialized(config)
	return nil
}

// Transcribe performs speech recognition with fallback chain
func (c *FallbackChain) Transcribe(ctx context.Context, audioPath string, callback types.ProgressCallback) (*types.RecognitionResult, error) {
	// Ensure initialized
	if err := c.CheckInitialized(); err != nil {
		return nil, err
	}

	// Auto-select recognizers by hardware if enabled
	recognizers := c.selectRecognizers()

	var lastError error

	// Try each recognizer in priority order
	for i, recognizer := range recognizers {
		// Report progress
		if callback != nil {
			progress := 10 + (i * 20)
			callback(progress, fmt.Sprintf("Trying recognizer %d of %d...", i+1, len(recognizers)))
		}

		// Try multiple attempts for this recognizer
		for attempt := 0; attempt < c.config.MaxAttempts; attempt++ {
			result, err := recognizer.Transcribe(ctx, audioPath, callback)
			if err == nil {
				// Success - return result
				if callback != nil {
					callback(100, "Complete")
				}
				return result, nil
			}

			// Check if error is retryable
			if !c.isRetryableError(err) {
				break
			}

			lastError = err
		}
	}

	// All recognizers failed
	return nil, fmt.Errorf("all recognizers failed, last error: %w", lastError)
}

// TranscribeStream is not supported by fallback chain
func (c *FallbackChain) TranscribeStream(ctx context.Context, chunks <-chan types.AudioChunk, results chan<- types.RecognitionResult) error {
	return types.NewAudioError(
		asr.ErrCodeFormatNotSupported,
		"Streaming transcription not supported by fallback chain",
		nil,
	)
}

// Release cleans up all recognizers
func (c *FallbackChain) Release() error {
	var errors []error

	for _, recognizer := range c.config.Recognizers {
		if err := recognizer.Release(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors releasing recognizers: %v", errors)
	}

	return nil
}

// selectRecognizers returns the list of recognizers to try
func (c *FallbackChain) selectRecognizers() []types.ASRRecognizer {
	if !c.config.AutoSelectByHardware {
		return c.config.Recognizers
	}

	// Detect hardware capabilities
	hardwareInfo := hw.DetectHardware()

	// Select models based on hardware
	recommendedModels := hw.SelectModelsForHardware(hardwareInfo)

	// Filter by minimum memory tier if specified
	if c.config.MinMemoryTier != "" && c.config.MinMemoryTier != types.TierAuto {
		recommendedModels = hw.FilterModelsByTier(recommendedModels, c.config.MinMemoryTier)
	}

	// Map model types to recognizers using the recognizer map
	var selectedRecognizers []types.ASRRecognizer
	seen := make(map[types.ASRRecognizer]bool)

	for _, modelType := range recommendedModels {
		if recognizer, ok := c.config.RecognizerMap[modelType]; ok {
			// Avoid duplicates
			if !seen[recognizer] {
				selectedRecognizers = append(selectedRecognizers, recognizer)
				seen[recognizer] = true
			}
		}
	}

	// If no recognizers selected, fall back to default list
	if len(selectedRecognizers) == 0 {
		return c.config.Recognizers
	}

	return selectedRecognizers
}

// isRetryableError checks if an error is worth retrying
func (c *FallbackChain) isRetryableError(err error) bool {
	// Don't retry on configuration errors
	if asr.IsModelNotLoaded(err) {
		return false
	}

	// Don't retry on unsupported language
	if asr.IsLanguageNotSupported(err) {
		return false
	}

	// Retry on transient errors (timeout, network, etc.)
	return true
}
