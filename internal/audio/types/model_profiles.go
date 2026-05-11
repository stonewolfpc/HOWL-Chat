package types

import (
	"fmt"
	"strings"
)

// MemoryTier represents hardware capability levels
type MemoryTier string

const (
	TierTiny   MemoryTier = "tiny"   // <2GB
	TierSmall  MemoryTier = "small"  // 4-8GB
	TierMedium MemoryTier = "medium" // 8-16GB
	TierLarge  MemoryTier = "large"  // 16GB+
	TierAuto   MemoryTier = "auto"
)

// ASRModelType identifies speech recognition models
type ASRModelType string

const (
	// Whisper models (OpenAI)
	ASRWhisperTiny  ASRModelType = "whisper-tiny"
	ASRWhisperBase  ASRModelType = "whisper-base"
	ASRWhisperSmall ASRModelType = "whisper-small"
	ASRWhisperMed   ASRModelType = "whisper-medium"
	ASRWhisperLarge ASRModelType = "whisper-large-v3"

	// Qwen ASR models (Alibaba)
	ASRQwen3ASR06B ASRModelType = "qwen3-asr-0.6b"
	ASRQwen3ASR17B ASRModelType = "qwen3-asr-1.7b"

	// Fallback
	ASRFallback ASRModelType = "fallback"
)

// TTSModelType identifies text-to-speech models
type TTSModelType string

const (
	// Piper TTS (fast, lightweight)
	TTSPiper     TTSModelType = "piper"
	TTSPiperFast TTSModelType = "piper-fast"

	// MeloTTS (better quality)
	TTSMeloTTS TTSModelType = "melotts"

	// Bark (highest quality, slowest)
	TTSBarkSmall TTSModelType = "bark-small"
	TTSBarkLarge TTSModelType = "bark-large"

	// Fallback
	TTSFallback TTSModelType = "system-fallback"
)

// MultimodalModelType identifies audio+vision+text models
type MultimodalModelType string

const (
	// Qwen Omni models
	MultimodalQwen25Omni3B MultimodalModelType = "qwen2.5-omni-3b"
	MultimodalQwen25Omni7B MultimodalModelType = "qwen2.5-omni-7b"
	MultimodalQwen3Omni30B MultimodalModelType = "qwen3-omni-30b"

	// Gemma 4 (multimodal)
	MultimodalGemma4E2B MultimodalModelType = "gemma-4-e2b"
	MultimodalGemma4E4B MultimodalModelType = "gemma-4-e4b"
	MultimodalGemma426B MultimodalModelType = "gemma-4-26b"
	MultimodalGemma431B MultimodalModelType = "gemma-4-31b"
)

// ModelProfile contains complete model metadata
type ModelProfile struct {
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Version        string   `json:"version"`
	VRAMBytes      int64    `json:"vram_bytes"`
	RAMBytes       int64    `json:"ram_bytes"`
	ModelSizeBytes int64    `json:"model_size_bytes"`
	Languages      []string `json:"languages"`
	License        string   `json:"license"`
	SourceURL      string   `json:"source_url"`
	GGUFAvailable  bool     `json:"gguf_available"`

	// Capabilities
	SupportsStreaming bool `json:"supports_streaming"`
	SupportsGPU       bool `json:"supports_gpu"`
	SupportsCPU       bool `json:"supports_cpu"`
	SupportsQuantized bool `json:"supports_quantized"`

	// Performance characteristics
	RealTimeFactor float64 `json:"real_time_factor"` // 1.0 = realtime, 0.5 = 2x faster
	QualityScore   float64 `json:"quality_score"`    // 0.0 - 1.0

	// Hardware requirements
	MinMemoryTier   MemoryTier `json:"min_memory_tier"`
	RecommendedTier MemoryTier `json:"recommended_tier"`

	// Integration
	Backend string `json:"backend"` // "local", "http", "grpc"
}

// ASRModelProfile extends ModelProfile for ASR-specific configs
type ASRModelProfile struct {
	ModelProfile
	MaxAudioDuration        float64 `json:"max_audio_duration"` // seconds
	BeamSize                int     `json:"beam_size"`
	BestOf                  int     `json:"best_of"`
	Temperature             float64 `json:"temperature"`
	Patience                float64 `json:"patience"`
	SuppressTokens          string  `json:"suppress_tokens"`
	InitialPrompt           string  `json:"initial_prompt"`
	ConditionOnPreviousText bool    `json:"condition_on_previous_text"`
}

// TTSModelProfile extends ModelProfile for TTS-specific configs
type TTSModelProfile struct {
	ModelProfile
	SampleRate      int      `json:"sample_rate"`
	SupportedVoices []string `json:"supported_voices"`
	SupportsSSML    bool     `json:"supports_ssml"`
	SupportsSpeed   bool     `json:"supports_speed"`
	MinSpeed        float64  `json:"min_speed"`
	MaxSpeed        float64  `json:"max_speed"`
}

// MultimodalModelProfile extends ModelProfile for multimodal configs
type MultimodalModelProfile struct {
	ModelProfile
	MaxAudioTokens     int  `json:"max_audio_tokens"`
	MaxVisionTokens    int  `json:"max_vision_tokens"`
	MaxTextTokens      int  `json:"max_text_tokens"`
	AudioSampleRate    int  `json:"audio_sample_rate"`
	SupportsVideo      bool `json:"supports_video"`
	StreamingLatencyMs int  `json:"streaming_latency_ms"`
}

// ModelRegistry contains all supported models
var ModelRegistry = struct {
	ASRModels        map[ASRModelType]ASRModelProfile
	TTSModels        map[TTSModelType]TTSModelProfile
	MultimodalModels map[MultimodalModelType]MultimodalModelProfile
}{
	ASRModels: map[ASRModelType]ASRModelProfile{
		ASRWhisperTiny: {
			ModelProfile: ModelProfile{
				Type:              "asr",
				Name:              "Whisper Tiny",
				Description:       "Fastest ASR, 39MB, realtime on CPU",
				Version:           "openai/whisper-v3",
				VRAMBytes:         0,                 // CPU-only
				RAMBytes:          100 * 1024 * 1024, // 100MB
				ModelSizeBytes:    39 * 1024 * 1024,
				Languages:         []string{"en", "zh", "de", "es", "fr", "it", "ja", "ko", "pt", "ru"},
				License:           "MIT",
				SourceURL:         "https://huggingface.co/ggerganov/whisper.cpp",
				GGUFAvailable:     true,
				SupportsStreaming: true,
				SupportsGPU:       true,
				SupportsCPU:       true,
				SupportsQuantized: true,
				RealTimeFactor:    0.3,
				QualityScore:      0.75,
				MinMemoryTier:     TierTiny,
				RecommendedTier:   TierTiny,
				Backend:           "local",
			},
			MaxAudioDuration: 300,
			BeamSize:         5,
			BestOf:           5,
			Temperature:      0.0,
			Patience:         -1.0,
		},
		ASRWhisperBase: {
			ModelProfile: ModelProfile{
				Type:              "asr",
				Name:              "Whisper Base",
				Description:       "Balanced ASR, 74MB, good quality",
				Version:           "openai/whisper-v3",
				VRAMBytes:         512 * 1024 * 1024,
				RAMBytes:          200 * 1024 * 1024,
				ModelSizeBytes:    74 * 1024 * 1024,
				Languages:         []string{"en", "zh", "de", "es", "fr", "it", "ja", "ko", "pt", "ru", "ar", "hi"},
				License:           "MIT",
				SourceURL:         "https://huggingface.co/ggerganov/whisper.cpp",
				GGUFAvailable:     true,
				SupportsStreaming: true,
				SupportsGPU:       true,
				SupportsCPU:       true,
				SupportsQuantized: true,
				RealTimeFactor:    0.5,
				QualityScore:      0.85,
				MinMemoryTier:     TierSmall,
				RecommendedTier:   TierSmall,
				Backend:           "local",
			},
			MaxAudioDuration: 600,
			BeamSize:         5,
			BestOf:           5,
			Temperature:      0.0,
			Patience:         -1.0,
		},
		ASRWhisperSmall: {
			ModelProfile: ModelProfile{
				Type:              "asr",
				Name:              "Whisper Small",
				Description:       "High quality ASR, 244MB",
				Version:           "openai/whisper-v3",
				VRAMBytes:         1024 * 1024 * 1024,
				RAMBytes:          512 * 1024 * 1024,
				ModelSizeBytes:    244 * 1024 * 1024,
				Languages:         []string{"en", "zh", "de", "es", "fr", "it", "ja", "ko", "pt", "ru", "ar", "hi", "pl", "tr", "vi"},
				License:           "MIT",
				SourceURL:         "https://huggingface.co/ggerganov/whisper.cpp",
				GGUFAvailable:     true,
				SupportsStreaming: false,
				SupportsGPU:       true,
				SupportsCPU:       true,
				SupportsQuantized: true,
				RealTimeFactor:    1.0,
				QualityScore:      0.92,
				MinMemoryTier:     TierSmall,
				RecommendedTier:   TierMedium,
				Backend:           "local",
			},
			MaxAudioDuration: 1200,
			BeamSize:         5,
			BestOf:           5,
			Temperature:      0.0,
			Patience:         -1.0,
		},
		ASRQwen3ASR06B: {
			ModelProfile: ModelProfile{
				Type:              "asr",
				Name:              "Qwen3 ASR 0.6B",
				Description:       "Tiny multilingual ASR, 600MB",
				Version:           "qwen3-asr-0.6b",
				VRAMBytes:         int64(2 * 1024 * 1024 * 1024),
				RAMBytes:          int64(1 * 1024 * 1024 * 1024),
				ModelSizeBytes:    int64(600 * 1024 * 1024),
				Languages:         []string{"en", "zh", "ja", "ko", "de", "fr", "es", "it", "pt", "ru"},
				License:           "qwen",
				SourceURL:         "https://huggingface.co/Qwen/Qwen3-ASR-0.6B",
				GGUFAvailable:     true,
				SupportsStreaming: true,
				SupportsGPU:       true,
				SupportsCPU:       true,
				SupportsQuantized: true,
				RealTimeFactor:    0.8,
				QualityScore:      0.88,
				MinMemoryTier:     TierSmall,
				RecommendedTier:   TierSmall,
				Backend:           "http",
			},
			MaxAudioDuration: 300,
			BeamSize:         1,
			BestOf:           1,
			Temperature:      0.0,
			Patience:         0.0,
		},
		ASRQwen3ASR17B: {
			ModelProfile: ModelProfile{
				Type:              "asr",
				Name:              "Qwen3 ASR 1.7B",
				Description:       "Best multilingual ASR, 1.7GB",
				Version:           "qwen3-asr-1.7b",
				VRAMBytes:         int64(4 * 1024 * 1024 * 1024),
				RAMBytes:          int64(2 * 1024 * 1024 * 1024),
				ModelSizeBytes:    int64(1.75 * 1024 * 1024 * 1024),
				Languages:         []string{"en", "zh", "ja", "ko", "de", "fr", "es", "it", "pt", "ru", "ar", "hi", "th", "vi"},
				License:           "qwen",
				SourceURL:         "https://huggingface.co/Qwen/Qwen3-ASR-1.7B",
				GGUFAvailable:     true,
				SupportsStreaming: true,
				SupportsGPU:       true,
				SupportsCPU:       false, // Too slow on CPU
				SupportsQuantized: true,
				RealTimeFactor:    1.2,
				QualityScore:      0.95,
				MinMemoryTier:     TierMedium,
				RecommendedTier:   TierMedium,
				Backend:           "http",
			},
			MaxAudioDuration: 600,
			BeamSize:         1,
			BestOf:           1,
			Temperature:      0.0,
			Patience:         0.0,
		},
	},
	TTSModels: map[TTSModelType]TTSModelProfile{
		TTSPiper: {
			ModelProfile: ModelProfile{
				Type:              "tts",
				Name:              "Piper",
				Description:       "Fast local TTS, ~100MB",
				Version:           "piper-2023",
				VRAMBytes:         0,
				RAMBytes:          200 * 1024 * 1024,
				ModelSizeBytes:    100 * 1024 * 1024,
				Languages:         []string{"en", "de", "es", "fr", "it", "nl", "pl", "pt", "ru", "zh"},
				License:           "MIT",
				SourceURL:         "https://github.com/rhasspy/piper",
				GGUFAvailable:     false,
				SupportsStreaming: true,
				SupportsGPU:       false, // CPU-only
				SupportsCPU:       true,
				SupportsQuantized: false,
				RealTimeFactor:    0.1,
				QualityScore:      0.80,
				MinMemoryTier:     TierTiny,
				RecommendedTier:   TierTiny,
				Backend:           "local",
			},
			SampleRate:      22050,
			SupportedVoices: []string{"en_US-lessac-medium", "en_US-amy-medium", "en_US-danny-low"},
			SupportsSSML:    false,
			SupportsSpeed:   true,
			MinSpeed:        0.5,
			MaxSpeed:        2.0,
		},
		TTSMeloTTS: {
			ModelProfile: ModelProfile{
				Type:              "tts",
				Name:              "MeloTTS",
				Description:       "Quality TTS, ~300MB",
				Version:           "melo-2024",
				VRAMBytes:         int64(1 * 1024 * 1024 * 1024),
				RAMBytes:          int64(500 * 1024 * 1024),
				ModelSizeBytes:    int64(300 * 1024 * 1024),
				Languages:         []string{"en", "zh", "ja", "ko", "es", "fr", "de", "it"},
				License:           "MIT",
				SourceURL:         "https://github.com/myshell-ai/MeloTTS",
				GGUFAvailable:     false,
				SupportsStreaming: false,
				SupportsGPU:       true,
				SupportsCPU:       true,
				SupportsQuantized: false,
				RealTimeFactor:    0.5,
				QualityScore:      0.90,
				MinMemoryTier:     TierSmall,
				RecommendedTier:   TierSmall,
				Backend:           "http",
			},
			SampleRate:      44100,
			SupportedVoices: []string{"EN", "ZH", "JP"},
			SupportsSSML:    false,
			SupportsSpeed:   true,
			MinSpeed:        0.8,
			MaxSpeed:        1.2,
		},
	},
	MultimodalModels: map[MultimodalModelType]MultimodalModelProfile{
		MultimodalQwen25Omni7B: {
			ModelProfile: ModelProfile{
				Type:              "multimodal",
				Name:              "Qwen2.5-Omni 7B",
				Description:       "Audio+Vision+Text, 7B params",
				Version:           "qwen2.5-omni-7b",
				VRAMBytes:         int64(8 * 1024 * 1024 * 1024),
				RAMBytes:          int64(4 * 1024 * 1024 * 1024),
				ModelSizeBytes:    int64(14 * 1024 * 1024 * 1024), // 7B params ≈ 14GB fp16
				Languages:         []string{"en", "zh", "ja", "ko", "de", "fr", "es"},
				License:           "qwen",
				SourceURL:         "https://huggingface.co/Qwen/Qwen2.5-Omni-7B",
				GGUFAvailable:     true,
				SupportsStreaming: true,
				SupportsGPU:       true,
				SupportsCPU:       false,
				SupportsQuantized: true,
				RealTimeFactor:    2.0,
				QualityScore:      0.92,
				MinMemoryTier:     TierLarge,
				RecommendedTier:   TierLarge,
				Backend:           "http",
			},
			MaxAudioTokens:     1500,
			MaxVisionTokens:    1024,
			MaxTextTokens:      2048,
			AudioSampleRate:    16000,
			SupportsVideo:      false,
			StreamingLatencyMs: 500,
		},
	},
}

// RecommendASR suggests best ASR model for hardware tier
func RecommendASR(tier MemoryTier, languages []string) ASRModelType {
	// Check language support
	supportsLang := func(modelLangs []string) bool {
		for _, lang := range languages {
			found := false
			for _, ml := range modelLangs {
				if strings.EqualFold(ml, lang) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}

	// Tier-based recommendations with language filter
	switch tier {
	case TierTiny, TierSmall:
		if supportsLang(ModelRegistry.ASRModels[ASRQwen3ASR06B].Languages) {
			return ASRQwen3ASR06B
		}
		return ASRWhisperBase
	case TierMedium:
		if supportsLang(ModelRegistry.ASRModels[ASRQwen3ASR17B].Languages) {
			return ASRQwen3ASR17B
		}
		return ASRWhisperSmall
	case TierLarge:
		if supportsLang(ModelRegistry.ASRModels[ASRQwen3ASR17B].Languages) {
			return ASRQwen3ASR17B
		}
		return ASRWhisperLarge
	default:
		return ASRWhisperBase
	}
}

// RecommendTTS suggests best TTS model for hardware tier
func RecommendTTS(tier MemoryTier) TTSModelType {
	switch tier {
	case TierTiny:
		return TTSPiper
	case TierSmall, TierMedium:
		return TTSPiper // Fast, good quality
	case TierLarge:
		return TTSMeloTTS // Better quality
	default:
		return TTSPiper
	}
}

// GetModelInfo retrieves profile by model type string
func GetModelInfo(modelType string) (*ModelProfile, error) {
	// Try ASR models
	if profile, ok := ModelRegistry.ASRModels[ASRModelType(modelType)]; ok {
		return &profile.ModelProfile, nil
	}

	// Try TTS models
	if profile, ok := ModelRegistry.TTSModels[TTSModelType(modelType)]; ok {
		return &profile.ModelProfile, nil
	}

	// Try multimodal models
	if profile, ok := ModelRegistry.MultimodalModels[MultimodalModelType(modelType)]; ok {
		return &profile.ModelProfile, nil
	}

	return nil, fmt.Errorf("model type not found: %s", modelType)
}

// GetASRProfile retrieves ASR profile by type
func GetASRProfile(modelType ASRModelType) (*ASRModelProfile, bool) {
	profile, ok := ModelRegistry.ASRModels[modelType]
	if !ok {
		return nil, false
	}
	return &profile, true
}

// GetTTSProfile retrieves TTS profile by type
func GetTTSProfile(modelType TTSModelType) (*TTSModelProfile, bool) {
	profile, ok := ModelRegistry.TTSModels[modelType]
	if !ok {
		return nil, false
	}
	return &profile, true
}

// ListModelsByTier returns all models suitable for a memory tier
func ListModelsByTier(tier MemoryTier) []string {
	var models []string

	// ASR models
	for t, profile := range ModelRegistry.ASRModels {
		if profile.MinMemoryTier <= tier {
			models = append(models, string(t))
		}
	}

	// TTS models
	for t, profile := range ModelRegistry.TTSModels {
		if profile.MinMemoryTier <= tier {
			models = append(models, string(t))
		}
	}

	// Multimodal models
	for t, profile := range ModelRegistry.MultimodalModels {
		if profile.MinMemoryTier <= tier {
			models = append(models, string(t))
		}
	}

	return models
}
