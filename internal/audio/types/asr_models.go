package types

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

// ASRModelRegistry contains all ASR models
var ASRModelRegistry = map[ASRModelType]ASRModelProfile{
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
}
