package whisper

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
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
	limiter *Limiter

	// Test seams for process discovery, execution, and hardware checks.
	cliSearchPaths []string
	runCLI         func(context.Context, []string) ([]byte, error)
	checkHardware  func() error
}

// New creates a new Whisper recognizer for the specified model type.
func New(modelType types.ASRModelType) *Recognizer {
	return NewWithLimiter(modelType, defaultLimiter)
}

// NewWithLimiter creates a recognizer with an isolated limiter. It is intended
// for tests and explicit startup wiring; callers must not mutate the package
// default from tests.
func NewWithLimiter(modelType types.ASRModelType, limiter *Limiter) *Recognizer {
	if limiter == nil {
		limiter = defaultLimiter
	}
	return &Recognizer{
		BaseRecognizer: asr.NewBaseRecognizer(modelType),
		limiter:        limiter,
	}
}

// Initialize prepares the recognizer for use. It never downloads or builds
// whisper-cli; a missing binary is a model-not-found error.
func (r *Recognizer) Initialize(ctx context.Context, config types.ASRConfig) error {
	r.LockCalls()
	defer r.UnlockCalls()

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

	cliPath, err := r.findWhisperCLI()
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", "howl-whisper-*")
	if err != nil {
		return types.NewAudioError(
			types.ErrCodeModelLoadFailed,
			"Failed to create temp directory",
			err,
		)
	}
	r.tempDir = tempDir

	keepTempDir := false
	defer func() {
		if !keepTempDir {
			_ = os.RemoveAll(tempDir)
			if r.tempDir == tempDir {
				r.tempDir = ""
			}
		}
	}()

	// Validate hardware compatibility
	if err := r.checkHardwareCompatibility(); err != nil {
		return err
	}

	r.modelPath = config.ModelPath
	r.whisperCliPath = cliPath
	r.SetInitialized(config)
	keepTempDir = true
	return nil
}

// Transcribe converts audio file to text using whisper-cli
func (r *Recognizer) Transcribe(ctx context.Context, audioPath string, callback types.ProgressCallback) (*types.RecognitionResult, error) {
	r.LockCalls()
	defer r.UnlockCalls()

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

	if r.limiter != nil {
		if err := r.limiter.Acquire(ctx); err != nil {
			return nil, types.NewAudioError(
				types.ErrCodeCancelled,
				"transcription cancelled while waiting for whisper capacity",
				err,
			)
		}
		defer r.limiter.Release()
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
	r.LockCalls()
	defer r.UnlockCalls()

	if r.tempDir != "" {
		dir := r.tempDir
		r.tempDir = ""
		r.ClearInitialized()
		return os.RemoveAll(dir)
	}
	r.ClearInitialized()
	return nil
}

func whisperInstallDir() string {
	if v := os.Getenv("HOWL_WHISPER_HOME"); v != "" {
		return filepath.Clean(v)
	}
	return filepath.Join("bin", "whisper")
}

// findWhisperCLI locates an existing whisper-cli. It never downloads or builds
// one.
func (r *Recognizer) findWhisperCLI() (string, error) {
	var searchPaths []string
	if r.cliSearchPaths != nil {
		searchPaths = r.cliSearchPaths
	} else {
		if path, err := exec.LookPath("whisper-cli"); err == nil {
			searchPaths = append(searchPaths, path)
		}
		if path, err := exec.LookPath("whisper-cli.exe"); err == nil {
			searchPaths = append(searchPaths, path)
		}
		searchPaths = append(searchPaths,
			"whisper-cli.exe",
			"whisper-cli",
			filepath.Join("bin", "whisper-cli.exe"),
			filepath.Join("bin", "whisper-cli"),
			filepath.Join(whisperInstallDir(), "whisper-cli.exe"),
			filepath.Join(whisperInstallDir(), "whisper-cli"),
			"C:\\Program Files\\whisper\\whisper-cli.exe",
			"/usr/local/bin/whisper-cli",
		)
	}

	for _, path := range searchPaths {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", types.NewAudioError(
		types.ErrCodeModelNotFound,
		"whisper-cli not found; install it or put it on PATH\n"+GetInstallationInstructions(),
		nil,
	)
}

func (r *Recognizer) checkHardwareCompatibility() error {
	if r.checkHardware != nil {
		return r.checkHardware()
	}
	return r.CheckHardwareCompatibility()
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

	var err error
	if r.runCLI != nil {
		_, err = r.runCLI(ctx, args)
	} else {
		// Execute command with context support
		cmd := execCommand(r.whisperCliPath, args...)

		// Run command with context cancellation support
		resultChan := make(chan error, 1)
		go func() {
			_, runErr := cmd.CombinedOutput()
			resultChan <- runErr
		}()

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
		case runErr := <-resultChan:
			err = runErr
		}
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

// Limiter bounds how many whisper-cli processes may run at once. It is a
// resource limit, not the provider correctness lock.
type Limiter struct {
	mu      sync.Mutex
	limit   int
	sem     chan struct{}
	waiters int
}

// NewLimiter creates an isolated limiter with the requested concurrency.
func NewLimiter(limit int) *Limiter {
	if limit < 1 {
		limit = 1
	}
	return &Limiter{
		limit: limit,
		sem:   make(chan struct{}, limit),
	}
}

// Acquire reserves one slot. It returns ctx.Err() if the context ends first.
func (l *Limiter) Acquire(ctx context.Context) error {
	l.mu.Lock()
	l.waiters++
	sem := l.sem
	l.mu.Unlock()

	select {
	case sem <- struct{}{}:
		l.mu.Lock()
		l.waiters--
		l.mu.Unlock()
		return nil
	case <-ctx.Done():
		l.mu.Lock()
		l.waiters--
		l.mu.Unlock()
		return ctx.Err()
	}
}

// Release returns one slot. A caller must only release an acquired slot.
func (l *Limiter) Release() {
	l.mu.Lock()
	sem := l.sem
	l.mu.Unlock()
	select {
	case <-sem:
	default:
	}
}

// SetLimit changes the limit only while no call is active or waiting.
func (l *Limiter) SetLimit(limit int) error {
	if limit < 1 {
		return fmt.Errorf("whisper concurrency limit must be at least 1")
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.limit == limit {
		return nil
	}
	if len(l.sem) > 0 || l.waiters > 0 {
		return fmt.Errorf("whisper concurrency limit cannot change while calls are active")
	}
	l.limit = limit
	l.sem = make(chan struct{}, limit)
	return nil
}

// Limit returns the configured concurrency.
func (l *Limiter) Limit() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.limit
}

var defaultLimiter = NewLimiter(1)

// SetDefaultConcurrencyLimit configures the package default at startup. It
// rejects changes while any recognizer is using the default limiter.
func SetDefaultConcurrencyLimit(limit int) error {
	return defaultLimiter.SetLimit(limit)
}

// execCommand is a helper to create exec.Command for testing purposes
var execCommand = func(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
