package whisper

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"howl-chat/internal/audio/asr"
	"howl-chat/internal/audio/types"
)

// Recognizer implements the ASRRecognizer interface using whisper.cpp
// It shells out to the whisper-cli executable for transcription
type Recognizer struct {
	*asr.BaseRecognizer

	// Binary configuration
	whisperCliPath string
	modelPath      string

	// Runtime state
	tempDir string
}

// New creates a new Whisper recognizer for the specified model type
func New(modelType types.ASRModelType) *Recognizer {
	return &Recognizer{
		BaseRecognizer: asr.NewBaseRecognizer(modelType),
	}
}

// Initialize prepares the recognizer for use
// Downloads whisper-cli binary and model if not present
func (r *Recognizer) Initialize(ctx context.Context, config types.ASRConfig) error {
	// Validate model path
	if config.ModelPath == "" {
		return types.NewAudioError(
			types.ErrCodeModelNotFound,
			"Model path required for Whisper initialization",
			nil,
		)
	}

	// Check model file exists
	if _, err := os.Stat(config.ModelPath); err != nil {
		return types.NewAudioError(
			types.ErrCodeModelNotFound,
			fmt.Sprintf("Whisper model not found: %s", config.ModelPath),
			err,
		)
	}

	r.modelPath = config.ModelPath

	// Find or download whisper-cli binary
	cliPath, err := r.findOrDownloadWhisperCLI(ctx)
	if err != nil {
		return err
	}
	r.whisperCliPath = cliPath

	// Create temp directory for outputs
	r.tempDir = filepath.Join(os.TempDir(), "howl-whisper")
	if err := os.MkdirAll(r.tempDir, 0755); err != nil {
		return types.NewAudioError(
			types.ErrCodeModelLoadFailed,
			"Failed to create temp directory",
			err,
		)
	}

	// Validate hardware compatibility
	if err := r.CheckHardwareCompatibility(); err != nil {
		return err
	}

	// Mark as initialized
	r.SetInitialized(config)

	return nil
}

// Transcribe converts audio file to text using whisper-cli
func (r *Recognizer) Transcribe(ctx context.Context, audioPath string, callback types.ProgressCallback) (*types.RecognitionResult, error) {
	// Ensure initialized
	if err := r.CheckInitialized(); err != nil {
		return nil, err
	}

	// Validate audio file
	if err := r.ValidateAudioFile(audioPath); err != nil {
		return nil, err
	}

	// Validate language if specified
	config := r.GetConfig()
	if err := r.ValidateLanguage(config.Language); err != nil {
		return nil, err
	}

	// Report progress - starting
	if callback != nil {
		callback(5, "Loading audio...")
	}

	// Skip audio conversion for now - whisper-cli handles MP3 directly
	// This avoids temp directory path issues
	compatiblePath := audioPath

	// Report progress - processing
	if callback != nil {
		callback(20, "Running speech recognition...")
	}

	// Build and execute whisper-cli command
	startTime := time.Now()
	outputPath, err := r.runWhisperCLI(ctx, compatiblePath, callback)
	if err != nil {
		return nil, err
	}
	defer r.cleanupTempFile(outputPath)

	// Report progress - parsing
	if callback != nil {
		callback(80, "Parsing results...")
	}

	// Parse JSON output
	result, err := parseWhisperOutput(outputPath)
	if err != nil {
		return nil, err
	}

	// Set processing duration
	result.Duration = time.Since(startTime).Seconds()

	// Validate result has content
	if result.Text == "" {
		return nil, asr.NewEmptyResultError()
	}

	// Report completion
	if callback != nil {
		callback(100, "Complete")
	}

	return result, nil
}

// TranscribeStream is not supported by whisper.cpp in standard builds
func (r *Recognizer) TranscribeStream(ctx context.Context, chunks <-chan types.AudioChunk, results chan<- types.RecognitionResult) error {
	return types.NewAudioError(
		types.ErrCodeFormatNotSupported,
		"Streaming transcription not supported by Whisper",
		nil,
	)
}

// Release frees resources and cleans up temp files
func (r *Recognizer) Release() error {
	// Clean up temp directory
	if r.tempDir != "" {
		if err := os.RemoveAll(r.tempDir); err != nil {
			return err
		}
		r.tempDir = ""
	}

	return nil
}

// findOrDownloadWhisperCLI locates the whisper-cli binary or downloads it
func (r *Recognizer) findOrDownloadWhisperCLI(ctx context.Context) (string, error) {
	// Search in standard locations
	searchPaths := []string{
		"whisper-cli.exe",                       // Current directory (Windows)
		"whisper-cli",                           // Current directory (Unix)
		filepath.Join("bin", "whisper-cli.exe"), // bin/ subdirectory
		filepath.Join("bin", "whisper-cli"),
		"C:\\Program Files\\whisper\\whisper-cli.exe", // Windows install path
		"/usr/local/bin/whisper-cli",                  // Unix install path
	}

	for _, path := range searchPaths {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// TODO: Implement download logic
	// For now, return error with instructions
	return "", types.NewAudioError(
		types.ErrCodeModelNotFound,
		"whisper-cli binary not found. Please download from https://github.com/ggerganov/whisper.cpp/releases "+
			"and place in PATH or bin/ directory",
		nil,
	)
}

// ensureWAVFormat converts audio to 16kHz mono WAV if needed
func (r *Recognizer) ensureWAVFormat(ctx context.Context, inputPath string) (string, error) {
	// Check if already in correct format
	metadata, err := types.GetAudioMetadata(inputPath)
	if err != nil {
		return "", err
	}

	// If already 16kHz mono WAV, return as-is
	if metadata.Format == types.FormatWAV &&
		metadata.SampleRate == 16000 &&
		metadata.Channels == 1 {
		return inputPath, nil
	}

	// Need to convert - use ffmpeg
	tempPath := filepath.Join(r.tempDir, fmt.Sprintf("whisper_input_%d.wav", time.Now().Unix()))

	if err := types.ConvertAudioFormat(inputPath, types.FormatWAV, tempPath); err != nil {
		return "", asr.NewTranscodingError(metadata.Format, types.FormatWAV, err)
	}

	return tempPath, nil
}

// runWhisperCLI executes the whisper-cli command and returns output JSON path
func (r *Recognizer) runWhisperCLI(ctx context.Context, audioPath string, callback types.ProgressCallback) (string, error) {
	// whisper-cli appends .json to the full filename (including original extension)
	// e.g., input.mp3 becomes input.mp3.json
	inputDir := filepath.Dir(audioPath)
	inputBase := filepath.Base(audioPath)
	expectedOutput := filepath.Join(inputDir, inputBase+".json")

	// Build command arguments - let whisper-cli handle audio format
	args := []string{
		"--model", r.modelPath,
		"--file", audioPath,
		"--output-json",
	}

	// Add language if specified (not auto)
	config := r.GetConfig()
	if config.Language != "" && config.Language != "auto" {
		args = append(args, "--language", config.Language)
	}

	// Add processing parameters
	if config.BeamSize > 0 {
		args = append(args, "--beam-size", fmt.Sprintf("%d", config.BeamSize))
	}
	if config.BestOf > 0 {
		args = append(args, "--best-of", fmt.Sprintf("%d", config.BestOf))
	}
	if config.NumThreads > 0 {
		args = append(args, "--threads", fmt.Sprintf("%d", config.NumThreads))
	}

	// Suppress tokens (standard set for speech)
	args = append(args, "--suppress-non-speech-tokens")

	// Execute command with context support
	cmd := execCommand(r.whisperCliPath, args...)

	// Run command with context cancellation support
	type result struct {
		err error
	}

	resultChan := make(chan result, 1)
	go func() {
		_, err := cmd.CombinedOutput()
		resultChan <- result{err}
	}()

	var err error

	select {
	case <-ctx.Done():
		// Context cancelled - kill process
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", types.NewAudioError(
			types.ErrCodeTimeout,
			"Transcription cancelled by user",
			ctx.Err(),
		)
	case res := <-resultChan:
		err = res.err
	}

	if err != nil {
		return "", types.NewAudioError(
			types.ErrCodeModelLoadFailed,
			fmt.Sprintf("whisper-cli failed: %v", err),
			err,
		)
	}

	// Verify output file was created
	if _, err := os.Stat(expectedOutput); err != nil {
		return "", types.NewAudioError(
			types.ErrCodeModelLoadFailed,
			fmt.Sprintf("whisper-cli did not produce output file at %s", expectedOutput),
			err,
		)
	}

	return expectedOutput, nil
}

// cleanupTempFile removes a temporary file if it's in our temp directory
func (r *Recognizer) cleanupTempFile(path string) {
	if path == "" {
		return
	}
	// Only delete if it's in our temp directory (safety check)
	if filepath.Dir(path) == r.tempDir {
		os.Remove(path)
	}
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

// execCommand is a helper to create exec.Command for testing purposes
var execCommand = func(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
