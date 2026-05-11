package types

import (
	"context"
	"io"
)

// AudioFormat represents supported audio codec formats
type AudioFormat string

const (
	FormatWAV  AudioFormat = "wav"
	FormatMP3  AudioFormat = "mp3"
	FormatOGG  AudioFormat = "ogg"
	FormatFLAC AudioFormat = "flac"
	FormatM4A  AudioFormat = "m4a"
)

// AudioMetadata contains properties of an audio stream
type AudioMetadata struct {
	Format       AudioFormat
	SampleRate   int     // Hz (e.g., 16000, 44100)
	Channels     int     // 1=mono, 2=stereo
	BitDepth     int     // 16, 24, 32
	Duration     float64 // seconds
	BitRate      int     // kbps (for compressed formats)
}

// AudioChunk represents a segment of audio data
type AudioChunk struct {
	Data     []byte
	Metadata AudioMetadata
	Index    int // Position in sequence for streaming
	IsLast   bool
}

// RecognitionResult contains transcribed text and metadata
type RecognitionResult struct {
	Text       string
	Confidence float64
	Language   string
	Segments   []Segment
	Duration   float64
}

// Segment represents a timestamped portion of transcription
type Segment struct {
	StartTime  float64
	EndTime    float64
	Text       string
	Confidence float64
}

// SynthesisRequest contains text to be spoken
type SynthesisRequest struct {
	Text     string
	Voice    string
	Speed    float64 // 0.5 - 2.0
	Language string
}

// SynthesisResult contains generated audio
type SynthesisResult struct {
	AudioPath string
	Metadata  AudioMetadata
	Duration  float64
}

// ProgressCallback reports operation progress
type ProgressCallback func(percent int, message string)

// ASRRecognizer defines the contract for speech-to-text engines
type ASRRecognizer interface {
	// Initialize prepares the recognizer for use
	Initialize(ctx context.Context, config ASRConfig) error

	// Transcribe converts audio file to text
	Transcribe(ctx context.Context, audioPath string, callback ProgressCallback) (*RecognitionResult, error)

	// TranscribeStream processes audio chunks in real-time
	TranscribeStream(ctx context.Context, chunks <-chan AudioChunk, results chan<- RecognitionResult) error

	// SupportsLanguage checks if language is supported
	SupportsLanguage(lang string) bool

	// GetModelInfo returns recognizer capabilities
	GetModelInfo() ASRModelInfo

	// Release frees resources
	Release() error
}

// ASRConfig contains recognizer configuration
type ASRConfig struct {
	ModelPath    string
	Language     string // "auto" or specific code
	Device       string // "cpu", "cuda", "auto"
	NumThreads   int
	BeamSize     int
	BestOf       int
	Temperature  float64
}

// ASRModelInfo describes recognizer capabilities
type ASRModelInfo struct {
	Name           string
	Version        string
	Languages      []string
	SupportsStream bool
	MaxDuration    float64
	VRAMBytes      int64
}

// TTSSynthesizer defines the contract for text-to-speech engines
type TTSSynthesizer interface {
	// Initialize prepares the synthesizer for use
	Initialize(ctx context.Context, config TTSConfig) error

	// Synthesize converts text to audio file
	Synthesize(ctx context.Context, req SynthesisRequest, callback ProgressCallback) (*SynthesisResult, error)

	// SynthesizeStream produces audio chunks for streaming playback
	SynthesizeStream(ctx context.Context, req SynthesisRequest, chunks chan<- AudioChunk) error

	// GetVoices returns available voice options
	GetVoices() []Voice

	// SupportsSSML checks if SSML markup is supported
	SupportsSSML() bool

	// Release frees resources
	Release() error
}

// TTSConfig contains synthesizer configuration
type TTSConfig struct {
	ModelPath    string
	Device       string
	CacheDir     string
	CacheEnabled bool
	MaxCacheSize int64 // MB
}

// Voice represents a TTS voice option
type Voice struct {
	ID          string
	Name        string
	Language    string
	Gender      string
	Description string
}

// AudioCache defines the contract for audio caching
type AudioCache interface {
	// Get retrieves cached audio if available
	Get(key string) (*SynthesisResult, bool)

	// Store saves audio to cache
	Store(key string, result *SynthesisResult) error

	// Invalidate removes specific entry
	Invalidate(key string)

	// Clear removes all entries
	Clear()

	// GetStats returns cache statistics
	GetStats() CacheStats
}

// CacheStats contains cache performance metrics
type CacheStats struct {
	HitCount    int64
	MissCount   int64
	EntryCount  int
	SizeBytes   int64
	MaxSizeBytes int64
}

// VADEngine defines voice activity detection interface
type VADEngine interface {
	// Initialize prepares VAD for use
	Initialize(ctx context.Context, config VADConfig) error

	// Process analyzes audio chunk for voice activity
	Process(chunk AudioChunk) VADResult

	// Reset clears internal state
	Reset()

	// Release frees resources
	Release() error
}

// VADConfig contains voice activity detection settings
type VADConfig struct {
	Threshold     float64 // 0.0 - 1.0
	MinSilenceMs  int
	MinSpeechMs   int
	PaddingMs     int
}

// VADResult contains voice detection outcome
type VADResult struct {
	IsSpeech    bool
	Confidence  float64
	SpeechStart bool // Transition to speech detected
	SpeechEnd   bool // Transition to silence detected
}

// MultimodalBridge combines audio, vision, and text
type MultimodalBridge interface {
	// Initialize prepares the multimodal processor
	Initialize(ctx context.Context, config MultimodalConfig) error

	// ProcessAudioText handles audio input + text context
	ProcessAudioText(ctx context.Context, audio AudioChunk, textContext string) (string, error)

	// ProcessAudioVision handles audio + image input
	ProcessAudioVision(ctx context.Context, audio AudioChunk, imageData []byte) (string, error)

	// GetTokenBudget returns remaining token budget for audio
	GetTokenBudget(audioDuration float64) int

	// Release frees resources
	Release() error
}

// MultimodalConfig contains multimodal settings
type MultimodalConfig struct {
	ModelPath      string
	MaxAudioTokens int
	MaxVisionTokens int
	TextContextTokens int
}

// AudioPipeline orchestrates the full audio flow
type AudioPipeline interface {
	// Initialize sets up the complete pipeline
	Initialize(ctx context.Context, config PipelineConfig) error

	// ProcessAudioMessage handles complete audio message flow
	ProcessAudioMessage(ctx context.Context, msg AudioMessage) (*AudioResponse, error)

	// EnableASR activates speech recognition
	EnableASR(recognizer ASRRecognizer) error

	// EnableTTS activates speech synthesis
	EnableTTS(synthesizer TTSSynthesizer) error

	// EnableMultimodal activates multimodal processing
	EnableMultimodal(bridge MultimodalBridge) error

	// GetStatus returns pipeline operational status
	GetStatus() PipelineStatus

	// Shutdown gracefully stops all components
	Shutdown(ctx context.Context) error
}

// PipelineConfig contains full pipeline configuration
type PipelineConfig struct {
	ASREnabled       bool
	TTSEnabled       bool
	MultimodalEnabled bool
	RealtimeEnabled   bool
	HardwareTier    string
}

// AudioMessage represents an audio input message
type AudioMessage struct {
	Type      string  // "audio_input", "audio_output"
	Text      string  // For TTS requests
	AudioPath string  // For ASR requests
	AudioData []byte  // Raw audio data
	Metadata  AudioMetadata
}

// AudioResponse contains processed audio result
type AudioResponse struct {
	TranscribedText string
	ResponseText    string
	AudioOutputPath string
	ProcessingTime  float64
}

// PipelineStatus reports operational state
type PipelineStatus struct {
	ASRReady        bool
	TTSReady        bool
	MultimodalReady bool
	RealtimeReady   bool
	ActiveSessions  int
	ErrorCount      int64
}

// FormatConverter handles audio format transformations
type FormatConverter interface {
	// Convert transforms audio from one format to another
	Convert(inputPath string, outputFormat AudioFormat, metadata AudioMetadata) (string, error)

	// ConvertStream transforms audio chunks on-the-fly
	ConvertStream(input <-chan AudioChunk, output chan<- AudioChunk, targetFormat AudioFormat) error

	// Resample changes sample rate
	Resample(inputPath string, targetSampleRate int) (string, error)

	// ExtractMetadata reads audio file properties
	ExtractMetadata(path string) (*AudioMetadata, error)
}

// ModelLoader handles audio model lifecycle
type ModelLoader interface {
	// LoadASR initializes an ASR model
	LoadASR(modelType string, config ASRConfig) (ASRRecognizer, error)

	// LoadTTS initializes a TTS model
	LoadTTS(modelType string, config TTSConfig) (TTSSynthesizer, error)

	// LoadMultimodal initializes multimodal model
	LoadMultimodal(modelType string, config MultimodalConfig) (MultimodalBridge, error)

	// Unload releases a specific model
	Unload(modelID string) error

	// GetLoadedModels returns currently active models
	GetLoadedModels() []string
}

// Error types for audio operations
type AudioErrorCode int

const (
	ErrCodeNone AudioErrorCode = iota
	ErrCodeFormatNotSupported
	ErrCodeModelNotFound
	ErrCodeModelLoadFailed
	ErrCodeTranscriptionFailed
	ErrCodeSynthesisFailed
	ErrCodeCacheError
	ErrCodeHardwareInsufficient
	ErrCodeTimeout
	ErrCodeCancelled
)

// AudioError provides structured error information
type AudioError struct {
	Code    AudioErrorCode
	Message string
	Cause   error
	Recoverable bool
}

func (e *AudioError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Utility functions for error handling

func NewAudioError(code AudioErrorCode, message string, cause error) *AudioError {
	return &AudioError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Recoverable: code != ErrCodeHardwareInsufficient && code != ErrCodeModelNotFound,
	}
}

// IsRecoverable checks if error allows fallback
func IsRecoverable(err error) bool {
	if audioErr, ok := err.(*AudioError); ok {
		return audioErr.Recoverable
	}
	return true
}

// ReadSeekCloser combines io.Reader, io.Seeker, and io.Closer
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
