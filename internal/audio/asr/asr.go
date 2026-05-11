package asr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"howl-chat/internal/audio/hw"
	"howl-chat/internal/audio/types"
)

// BaseRecognizer provides common functionality for all ASR implementations
type BaseRecognizer struct {
	modelType   types.ASRModelType
	config      types.ASRConfig
	initialized bool
	info        types.ASRModelInfo
}

// NewBaseRecognizer creates a base recognizer with model metadata
func NewBaseRecognizer(modelType types.ASRModelType) *BaseRecognizer {
	profile, ok := types.GetASRProfile(modelType)
	if !ok {
		// Use default tiny profile if not found
		profile, _ = types.GetASRProfile(types.ASRWhisperTiny)
	}

	return &BaseRecognizer{
		modelType: modelType,
		info: types.ASRModelInfo{
			Name:           profile.Name,
			Version:        profile.Version,
			Languages:      profile.Languages,
			SupportsStream: profile.SupportsStreaming,
			MaxDuration:    profile.MaxAudioDuration,
			VRAMBytes:      profile.VRAMBytes,
		},
	}
}

// GetModelType returns the ASR model type
func (b *BaseRecognizer) GetModelType() types.ASRModelType {
	return b.modelType
}

// IsInitialized checks if recognizer has been initialized
func (b *BaseRecognizer) IsInitialized() bool {
	return b.initialized
}

// SetInitialized marks recognizer as initialized
func (b *BaseRecognizer) SetInitialized(config types.ASRConfig) {
	b.initialized = true
	b.config = config
}

// GetConfig returns the current configuration
func (b *BaseRecognizer) GetConfig() types.ASRConfig {
	return b.config
}

// GetModelInfo returns model capabilities
func (b *BaseRecognizer) GetModelInfo() types.ASRModelInfo {
	return b.info
}

// SupportsLanguage checks if the recognizer can process a language code.
func (b *BaseRecognizer) SupportsLanguage(lang string) bool {
	if lang == "" || lang == "auto" {
		return true
	}
	for _, supported := range b.info.Languages {
		if supported == lang {
			return true
		}
	}
	return false
}

// CheckInitialized verifies recognizer is ready for use
func (b *BaseRecognizer) CheckInitialized() error {
	if !b.initialized {
		return NewModelNotLoadedError(string(b.modelType))
	}
	return nil
}

// ValidateLanguage checks if language is supported
func (b *BaseRecognizer) ValidateLanguage(lang string) error {
	if lang == "" || lang == "auto" {
		return nil // Auto-detection is always supported
	}

	for _, supported := range b.info.Languages {
		if supported == lang {
			return nil
		}
	}

	return NewLanguageNotSupportedError(lang, b.info.Languages)
}

// ValidateAudioFile checks file format and constraints
func (b *BaseRecognizer) ValidateAudioFile(path string) error {
	// Check basic file validity
	if err := types.ValidateAudioFile(path); err != nil {
		return err
	}

	// Get metadata for duration check - skip if ffprobe not available
	metadata, err := types.GetAudioMetadata(path)
	if err != nil {
		// ffprobe not available - skip duration check
		fmt.Printf("WARN: Could not get audio metadata (ffprobe not found), skipping validation\n")
		return nil
	}

	// Check duration limits
	if metadata.Duration > b.info.MaxDuration {
		return types.NewAudioError(
			types.ErrCodeFormatNotSupported,
			fmt.Sprintf("Audio duration %.1fs exceeds maximum %.1fs for %s",
				metadata.Duration, b.info.MaxDuration, b.modelType),
			nil,
		)
	}

	return nil
}

// EnsureCompatibleFormat converts audio to processing format
func (b *BaseRecognizer) EnsureCompatibleFormat(ctx context.Context, inputPath string) (string, error) {
	metadata, err := types.GetAudioMetadata(inputPath)
	if err != nil {
		// ffprobe not available - skip format check, assume whisper can handle it
		fmt.Printf("WARN: Could not get audio metadata (ffprobe not found), assuming format is compatible\n")
		return inputPath, nil
	}

	// Check if already in target format
	targetFormat := types.GetProcessingFormat()
	if metadata.Format == targetFormat.Format &&
		metadata.SampleRate == targetFormat.SampleRate &&
		metadata.Channels == targetFormat.Channels {
		return inputPath, nil // Already compatible
	}

	dir := os.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "howl-asr"), 0o755); err != nil {
		return "", NewTranscodingError(metadata.Format, targetFormat.Format, err)
	}
	outPath := filepath.Join(dir, "howl-asr", fmt.Sprintf("compat_%d.wav", time.Now().UnixNano()))
	if err := types.ConvertAudioFormat(inputPath, targetFormat.Format, outPath); err != nil {
		return "", NewTranscodingError(metadata.Format, targetFormat.Format, err)
	}
	return outPath, nil
}

// CheckHardwareCompatibility verifies sufficient resources
func (b *BaseRecognizer) CheckHardwareCompatibility() error {
	profile, ok := types.GetASRProfile(b.modelType)
	if !ok {
		return nil
	}
	h := hw.DetectHardware()
	minRAM := profile.RAMBytes + profile.ModelSizeBytes
	if profile.VRAMBytes > 0 && !h.HasGPU {
		return types.NewAudioError(types.ErrCodeHardwareInsufficient,
			fmt.Sprintf("model %s requires a GPU (%d bytes VRAM profile); none detected",
				b.modelType, profile.VRAMBytes),
			h.Error)
	}
	if profile.VRAMBytes > 0 && h.GPUMemoryBytes > 0 && h.GPUMemoryBytes < profile.VRAMBytes {
		return types.NewAudioError(types.ErrCodeHardwareInsufficient,
			fmt.Sprintf("model %s needs about %s VRAM but reported usable GPU memory is about %s",
				b.modelType, types.FormatBytes(profile.VRAMBytes), types.FormatBytes(h.GPUMemoryBytes)),
			nil)
	}
	if minRAM <= 0 {
		return nil
	}
	available := h.AvailableRAMBytes
	const headroomMul = float64(1.15)
	need := float64(minRAM) * headroomMul
	if float64(available) < need {
		return types.NewAudioError(types.ErrCodeHardwareInsufficient,
			fmt.Sprintf("insufficient RAM for %s (need approximately %s free, ~%s with headroom; have ~%s free)",
				b.modelType,
				types.FormatBytes(minRAM),
				types.FormatBytes(int64(need)),
				types.FormatBytes(available)),
			h.Error)
	}
	return nil
}

// CalculateProgress computes percentage based on audio duration
func CalculateProgress(elapsedTime, totalDuration float64) int {
	if totalDuration <= 0 {
		return 0
	}

	// Assume processing takes roughly real-time or faster
	rtf := 1.0 // Real-time factor (will vary by model)
	expectedDuration := totalDuration * rtf

	if expectedDuration <= 0 {
		return 0
	}

	progress := int((elapsedTime / expectedDuration) * 100)
	if progress > 99 {
		progress = 99 // Leave 100% for finalization
	}
	if progress < 0 {
		progress = 0
	}

	return progress
}

// RecognizerFactory creates ASR recognizers based on type
type RecognizerFactory struct {
	registry map[types.ASRModelType]func() types.ASRRecognizer
}

// NewRecognizerFactory creates a factory with registered recognizers
func NewRecognizerFactory() *RecognizerFactory {
	return &RecognizerFactory{
		registry: make(map[types.ASRModelType]func() types.ASRRecognizer),
	}
}

// Register adds a recognizer constructor to the factory
func (f *RecognizerFactory) Register(modelType types.ASRModelType, constructor func() types.ASRRecognizer) {
	f.registry[modelType] = constructor
}

// Create instantiates a recognizer by type
func (f *RecognizerFactory) Create(modelType types.ASRModelType) (types.ASRRecognizer, error) {
	constructor, ok := f.registry[modelType]
	if !ok {
		return nil, fmt.Errorf("no recognizer registered for model type: %s", modelType)
	}
	return constructor(), nil
}

// ListAvailable returns all registered model types
func (f *RecognizerFactory) ListAvailable() []types.ASRModelType {
	var out []types.ASRModelType
	for t := range f.registry {
		out = append(out, t)
	}
	return out
}

// RecommendForHardware suggests best model for current hardware
func RecommendForHardware(preferredLang string) types.ASRModelType {
	h := hw.DetectHardware()
	tier := h.MemoryTier

	langs := []string{}
	if preferredLang != "" && preferredLang != "auto" {
		langs = append(langs, preferredLang)
	}

	return types.RecommendASR(tier, langs)
}
