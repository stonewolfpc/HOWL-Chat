/**
 * Model Profile Types
 *
 * This package defines types for model profiles derived from GGUF metadata.
 * Model profiles contain auto-detected configuration from GGUF files.
 *
 * @package types
 */

package types

// ModelProfile is derived from GGUF metadata and contains auto-detected
// model configuration such as architecture, context length, RoPE settings,
// and chat templates.
type ModelProfile struct {
	// Path is the file path to the GGUF model
	Path string

	// Family is the model architecture family (llama, mistral, qwen, lfm, etc.)
	Family string

	// Variant is the model variant name (7B, 8x7B, 70B, etc.)
	Variant string

	// MaxContext is the maximum context window from GGUF metadata
	MaxContext int

	// UsableContext is the safe context window with safety margins applied
	UsableContext int

	// RopeMode is the RoPE scaling mode (none, linear, dynamic, yarn)
	RopeMode string

	// RopeFactor is the RoPE frequency scaling factor
	RopeFactor float64

	// RopeBase is the RoPE frequency base
	RopeBase float64

	// Template is the resolved chat template from metadata or fallback
	Template string

	// Tokenizer is the tokenizer.ggml.model value
	Tokenizer string

	// StopSequences are default stop sequences from metadata
	StopSequences []string
}

// RuntimeSettings are the settings passed to gollama.cpp for model loading.
// These are the actual runtime parameters that affect model performance.
type RuntimeSettings struct {
	// Threads is the number of CPU threads to use
	Threads int

	// BatchSize is the logical batch size for processing
	BatchSize int

	// ContextSize is the context window size (0 = use profile.UsableContext)
	ContextSize int

	// GPULayers is the number of layers to offload to GPU
	GPULayers int

	// RopeMode is the RoPE scaling mode (none, linear, dynamic, yarn)
	RopeMode string

	// RopeFactor is the RoPE frequency scaling factor
	RopeFactor float64

	// RopeBase is the RoPE frequency base
	RopeBase float64

	// FlashAttention enables flash attention optimization
	FlashAttention bool

	// TensorSplit is the tensor split for multi-GPU
	TensorSplit []float32

	// MainGPU is the main GPU device
	MainGPU int

	// OffloadKQV enables KV cache offload to GPU
	OffloadKQV bool

	// UseMMap enables memory mapping
	UseMMap bool

	// UseMLock enables memory locking
	UseMLock bool

	// VocabOnly loads only vocabulary without weights
	VocabOnly bool
}

// PromptSettings are application-level prompt configuration settings.
// These are not passed to gollama.cpp but used in the prompt pipeline.
type PromptSettings struct {
	// PromptTemplate is the custom prompt template (empty = use profile.Template)
	PromptTemplate string

	// CustomJinjaTemplate is a custom Jinja template for advanced formatting
	CustomJinjaTemplate string

	// SystemPromptOverride overrides the default system prompt
	SystemPromptOverride string

	// UserPrefix is prepended to user messages
	UserPrefix string

	// AssistantPrefix is prepended to assistant messages
	AssistantPrefix string

	// StopSequences overrides profile stop sequences if non-empty
	StopSequences []string
}

// VisionSettings are the settings for vision/multimodal capabilities.
type VisionSettings struct {
	// Enabled enables vision/multimodal support
	Enabled bool

	// ClipModelPath is the path to the CLIP model GGUF file
	ClipModelPath string

	// MaxImageSize is the maximum image size in pixels
	MaxImageSize int

	// AutoResize automatically resizes images to MaxImageSize
	AutoResize bool

	// Normalize normalizes pixel values for CLIP
	Normalize bool
}

// VisionProfile contains CLIP model metadata extracted from GGUF.
type VisionProfile struct {
	// ClipModelPath is the path to the CLIP model GGUF file
	ClipModelPath string

	// ImageSize is the CLIP image size (typically 336 for CLIP)
	ImageSize int

	// PatchSize is the CLIP patch size (typically 14)
	PatchSize int

	// EmbeddingDim is the CLIP embedding dimension
	EmbeddingDim int

	// ContextLength is the CLIP context length (typically 77)
	ContextLength int

	// ModelType is the CLIP model type (clip, siglip, etc.)
	ModelType string
}
