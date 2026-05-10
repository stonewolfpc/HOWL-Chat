/**
 * Llama Options
 *
 * This package defines configuration options for the llama client.
 * These options control model loading, inference parameters, and system behavior.
 *
 * @package llama
 */

package llama

import "fmt"

// LoadOptions represents options for loading a model
type LoadOptions struct {
	ModelPath        string           `json:"model_path"`        // Path to the model file
	GPULayers        int              `json:"gpu_layers"`        // Number of layers to offload to GPU
	Threads          int              `json:"threads"`           // Number of CPU threads
	UseMMap          bool             `json:"use_mmap"`          // Whether to use memory mapping
	UseMlock         bool             `json:"use_mlock"`         // Whether to use mlock for pinned memory
	VocabOnly        bool             `json:"vocab_only"`        // Whether to load only vocabulary
	RMSNormEPS       float64          `json:"rms_norm_eps"`      // RMS normalization epsilon
	RopeFreqBase     float64          `json:"rope_freq_base"`    // RoPE frequency base
	RopeFreqScale    float64          `json:"rope_freq_scale"`   // RoPE frequency scale
	MemoryFraction   float64          `json:"memory_fraction"`   // Fraction of memory to use
	ProgressCallback ProgressCallback `json:"-"`                 // Callback for loading progress
	ContextSize      int              `json:"context_size"`      // Context window size
	Seed             int              `json:"seed"`              // Random seed
	MainGPU          int              `json:"main_gpu"`          // Main GPU device
	TensorSplit      string           `json:"tensor_split"`      // Tensor split for multi-GPU
	SplitMode        int              `json:"split_mode"`        // Split mode for multi-GPU
	EmbeddingMode    bool             `json:"embedding_mode"`    // Embedding mode
	NumGQA           int              `json:"num_gqa"`           // Number of GQA groups
	RopeScalingType  int              `json:"rope_scaling_type"` // RoPE scaling type
	YarnExtFactor    float64          `json:"yarn_ext_factor"`   // YARN extension factor
	YarnAttnFactor   float64          `json:"yarn_attn_factor"`  // YARN attention factor
	YarnBetaFast     float64          `json:"yarn_beta_fast"`    // YARN beta fast
	YarnBetaSlow     float64          `json:"yarn_beta_slow"`    // YARN beta slow
	YarnOrigCtx      int              `json:"yarn_orig_ctx"`     // YARN original context
}

// InferenceOptions represents options for text generation
type InferenceOptions struct {
	Temperature             float64  `json:"temperature"`                // Sampling temperature (0.0 = deterministic)
	TopPEnabled             bool     `json:"top_p_enabled"`              // Whether Top-P sampling is enabled
	TopP                    float64  `json:"top_p"`                      // Nucleus sampling parameter
	TopK                    int      `json:"top_k"`                      // Top-K sampling parameter
	MinP                    float64  `json:"min_p"`                      // Minimum probability threshold
	RepeatPenalty           float64  `json:"repeat_penalty"`             // Repetition penalty
	RepeatLastN             int      `json:"repeat_last_n"`              // How many recent tokens repeat penalty applies to
	FrequencyPenaltyEnabled bool     `json:"frequency_penalty_enabled"`  // Whether frequency penalty is enabled
	FrequencyPenalty        float64  `json:"frequency_penalty"`          // Frequency penalty
	PresencePenaltyEnabled  bool     `json:"presence_penalty_enabled"`   // Whether presence penalty is enabled
	PresencePenalty         float64  `json:"presence_penalty"`           // Presence penalty
	TFS                     float64  `json:"tfs"`                        // Tail Free Sampling parameter
	TypicalPEnabled         bool     `json:"typical_p_enabled"`          // Whether Typical-P sampling is enabled
	TypicalP                float64  `json:"typical_p"`                  // Typical sampling parameter
	MirostatEnabled         bool     `json:"mirostat_enabled"`           // Whether Mirostat sampling is enabled
	Mirostat                int      `json:"mirostat"`                   // Mirostat sampling mode (0=disabled, 1=mirostat, 2=mirostat_v2)
	MirostatTau             float64  `json:"mirostat_tau"`               // Mirostat target entropy
	MirostatETA             float64  `json:"mirostat_eta"`               // Mirostat learning rate
	DynamicTempRangeEnabled bool     `json:"dynamic_temp_range_enabled"` // Whether dynamic temperature is enabled
	DynamicTempRange        float64  `json:"dynamic_temp_range"`         // Dynamic temperature range
	DynamicTempExponent     float64  `json:"dynamic_temp_exponent"`      // Dynamic temperature exponent
	DRYMultiplier           float64  `json:"dry_multiplier"`             // DRY multiplier
	DRYAllowedLength        int      `json:"dry_allowed_length"`         // DRY allowed length
	DRYBase                 float64  `json:"dry_base"`                   // DRY base penalty
	SmoothingFactor         float64  `json:"smoothing_factor"`           // Smoothing factor
	SmoothingCurve          float64  `json:"smoothing_curve"`            // Smoothing curve
	TopAEnabled             bool     `json:"top_a_enabled"`              // Whether Top-A sampling is enabled
	TopA                    float64  `json:"top_a"`                      // Top-A sampling parameter
	EpsilonCutoff           float64  `json:"epsilon_cutoff"`             // Epsilon cutoff
	EtaCutoff               float64  `json:"eta_cutoff"`                 // Eta cutoff
	EncoderRepeatPenalty    float64  `json:"encoder_repeat_penalty"`     // Encoder repetition penalty
	NoRepeatNGramSize       int      `json:"no_repeat_ngram_size"`       // No-repeat n-gram size
	Seed                    int      `json:"seed"`                       // Random seed for generation
	NPredict                int      `json:"n_predict"`                  // Maximum tokens to generate
	NKeep                   int      `json:"n_keep"`                     // Number of tokens to keep from prompt
	StopStrings             []string `json:"stop_strings"`               // Strings that stop generation
}

// ContextOptions represents options for context management
type ContextOptions struct {
	ContextSize     int  `json:"context_size"`      // Context window size
	BatchSize       int  `json:"batch_size"`        // Batch size for processing
	FlashAttention  bool `json:"flash_attention"`   // Whether to use flash attention
	ThreadsBatch    int  `json:"threads_batch"`     // Batch processing threads
	OffloadKv       bool `json:"offload_kv"`        // Whether to offload KV cache to GPU
	RopeScalingType int  `json:"rope_scaling_type"` // RoPE scaling type
}

// TokenizerOptions represents options for tokenization
type TokenizerOptions struct {
	AddBOS     bool `json:"add_bos"`      // Whether to add BOS token
	AddEOS     bool `json:"add_eos"`      // Whether to add EOS token
	PadTokenID int  `json:"pad_token_id"` // Padding token ID
}

// NewLoadOptions creates default load options
func NewLoadOptions(modelPath string) *LoadOptions {
	return &LoadOptions{
		ModelPath:        modelPath,
		GPULayers:        0,
		Threads:          4,
		UseMMap:          true,
		UseMlock:         false,
		VocabOnly:        false,
		RMSNormEPS:       1e-5,
		RopeFreqBase:     10000.0,
		RopeFreqScale:    1.0,
		MemoryFraction:   0.9,
		ProgressCallback: nil,
		MainGPU:          -1, // -1 means no specific GPU (CPU-only safe)
		SplitMode:        -1, // -1 means no split mode (CPU-only safe)
	}
}

// NewInferenceOptions creates default inference options
func NewInferenceOptions() *InferenceOptions {
	return &InferenceOptions{
		Temperature:             0.7,
		TopPEnabled:             true,
		TopP:                    0.9,
		TopK:                    40,
		MinP:                    0.05,
		RepeatPenalty:           1.1,
		RepeatLastN:             64,
		FrequencyPenaltyEnabled: true,
		FrequencyPenalty:        0.0,
		PresencePenaltyEnabled:  true,
		PresencePenalty:         0.0,
		TFS:                     1.0,
		TypicalPEnabled:         true,
		TypicalP:                1.0,
		MirostatEnabled:         true,
		Mirostat:                0,
		MirostatTau:             5.0,
		MirostatETA:             0.1,
		DynamicTempRangeEnabled: true,
		DynamicTempRange:        0.0,
		DynamicTempExponent:     1.0,
		DRYMultiplier:           0.0,
		DRYAllowedLength:        2,
		DRYBase:                 1.0,
		SmoothingFactor:         0.0,
		SmoothingCurve:          1.0,
		TopAEnabled:             true,
		TopA:                    0.0,
		EpsilonCutoff:           0.0,
		EtaCutoff:               0.0,
		EncoderRepeatPenalty:    1.0,
		NoRepeatNGramSize:       0,
		Seed:                    -1,
		NPredict:                512,
		NKeep:                   0,
		StopStrings:             []string{},
	}
}

// NewContextOptions creates default context options
func NewContextOptions() *ContextOptions {
	return &ContextOptions{
		ContextSize:    2048,
		BatchSize:      512,
		FlashAttention: false,
	}
}

// NewTokenizerOptions creates default tokenizer options
func NewTokenizerOptions() *TokenizerOptions {
	return &TokenizerOptions{
		AddBOS:     true,
		AddEOS:     false,
		PadTokenID: 0,
	}
}

// Validate validates the load options
func (o *LoadOptions) Validate() error {
	if o.ModelPath == "" {
		return ErrInvalidModelPath
	}
	if o.Threads < 1 {
		return ErrInvalidThreads
	}
	if o.GPULayers < 0 {
		return ErrInvalidGPULayers
	}
	if o.MemoryFraction <= 0 || o.MemoryFraction > 1 {
		return ErrInvalidMemoryFraction
	}
	return nil
}

// Validate validates the inference options
func (o *InferenceOptions) Validate() error {
	if o.Temperature < 0 {
		return ErrInvalidTemperature
	}
	if o.TopP < 0 || o.TopP > 1 {
		return ErrInvalidTopP
	}
	if o.TopK < 0 {
		return ErrInvalidTopK
	}
	if o.NPredict < 1 {
		return ErrInvalidNPredict
	}
	return nil
}

// ProgressCallback is a function called during model loading to report progress
type ProgressCallback func(progress float64, stage LoadingStage)

// LoadingStage represents a stage in the model loading process
type LoadingStage string

const (
	// StageModelLoadStart indicates model loading has started
	StageModelLoadStart LoadingStage = "model_load_start"
	// StageTokenizerLoad indicates tokenizer is being loaded
	StageTokenizerLoad LoadingStage = "tokenizer_load"
	// StageTensorAllocation indicates tensors are being allocated
	StageTensorAllocation LoadingStage = "tensor_allocation"
	// StageKVCacheInit indicates KV cache is being initialized
	StageKVCacheInit LoadingStage = "kv_cache_init"
	// StageModelReady indicates model is ready for inference
	StageModelReady LoadingStage = "model_ready"
)

// Error types for options validation
var (
	ErrInvalidModelPath      = &ValidationError{Field: "model_path", Message: "model path cannot be empty"}
	ErrInvalidThreads        = &ValidationError{Field: "threads", Message: "threads must be at least 1"}
	ErrInvalidGPULayers      = &ValidationError{Field: "gpu_layers", Message: "gpu_layers cannot be negative"}
	ErrInvalidMemoryFraction = &ValidationError{Field: "memory_fraction", Message: "memory_fraction must be between 0 and 1"}
	ErrInvalidTemperature    = &ValidationError{Field: "temperature", Message: "temperature cannot be negative"}
	ErrInvalidTopP           = &ValidationError{Field: "top_p", Message: "top_p must be between 0 and 1"}
	ErrInvalidTopK           = &ValidationError{Field: "top_k", Message: "top_k cannot be negative"}
	ErrInvalidNPredict       = &ValidationError{Field: "n_predict", Message: "n_predict must be at least 1"}
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}
