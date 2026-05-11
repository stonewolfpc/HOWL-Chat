package asr

import (
	"fmt"

	"howl-chat/internal/audio/types"
)

// ASR-specific error codes extend the base audio error codes
const (
	ErrCodeModelNotLoaded types.AudioErrorCode = iota + 100
	ErrCodeAudioLoadFailed
	ErrCodeTranscodingFailed
	ErrCodeRecognitionTimeout
	ErrCodeEmptyResult
	ErrCodeLanguageNotSupported
	ErrCodeModelIncompatible
	ErrCodeInvalidConfig
	ErrCodeAPIError
	ErrCodeFormatNotSupported
)

// NewModelNotLoadedError creates error when model hasn't been initialized
func NewModelNotLoadedError(modelType string) *types.AudioError {
	return types.NewAudioError(
		ErrCodeModelNotLoaded,
		fmt.Sprintf("ASR model '%s' not loaded - call Initialize() first", modelType),
		nil,
	)
}

// NewAudioLoadError creates error when audio file can't be loaded
func NewAudioLoadError(path string, cause error) *types.AudioError {
	return types.NewAudioError(
		ErrCodeAudioLoadFailed,
		fmt.Sprintf("Failed to load audio file: %s", path),
		cause,
	)
}

// NewTranscodingError creates error when format conversion fails
func NewTranscodingError(fromFormat, toFormat types.AudioFormat, cause error) *types.AudioError {
	return types.NewAudioError(
		ErrCodeTranscodingFailed,
		fmt.Sprintf("Failed to convert %s to %s", fromFormat, toFormat),
		cause,
	)
}

// NewRecognitionTimeoutError creates error when transcription takes too long
func NewRecognitionTimeoutError(duration float64) *types.AudioError {
	return types.NewAudioError(
		ErrCodeRecognitionTimeout,
		fmt.Sprintf("Recognition timed out after %.1f seconds", duration),
		nil,
	)
}

// NewEmptyResultError creates error when no speech was detected
func NewEmptyResultError() *types.AudioError {
	return types.NewAudioError(
		ErrCodeEmptyResult,
		"No speech detected in audio",
		nil,
	)
}

// NewLanguageNotSupportedError creates error for unsupported languages
func NewLanguageNotSupportedError(lang string, supported []string) *types.AudioError {
	return types.NewAudioError(
		ErrCodeLanguageNotSupported,
		fmt.Sprintf("Language '%s' not supported. Available: %v", lang, supported),
		nil,
	)
}

// NewModelIncompatibleError creates error for model/hardware mismatch
func NewModelIncompatibleError(modelType string, requiredRAM, availableRAM int64) *types.AudioError {
	return types.NewAudioError(
		ErrCodeModelIncompatible,
		fmt.Sprintf("Model '%s' requires %s RAM but only %s available",
			modelType,
			types.FormatBytes(requiredRAM),
			types.FormatBytes(availableRAM),
		),
		nil,
	)
}

// IsModelNotLoaded checks if error is uninitialized model
func IsModelNotLoaded(err error) bool {
	if audioErr, ok := err.(*types.AudioError); ok {
		return audioErr.Code == ErrCodeModelNotLoaded
	}
	return false
}

// IsEmptyResult checks if error is no speech detected
func IsEmptyResult(err error) bool {
	if audioErr, ok := err.(*types.AudioError); ok {
		return audioErr.Code == ErrCodeEmptyResult
	}
	return false
}

// IsLanguageNotSupported checks if error is unsupported language
func IsLanguageNotSupported(err error) bool {
	if audioErr, ok := err.(*types.AudioError); ok {
		return audioErr.Code == ErrCodeLanguageNotSupported
	}
	return false
}
