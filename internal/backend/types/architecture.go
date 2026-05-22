package types

// ArchitectureSpec captures all architecture-specific parameters extracted
// from GGUF metadata. It is model-family agnostic — any transformer-based
// GGUF model can be described by these fields.
type ArchitectureSpec struct {
	// Family is the model architecture family (llama, mistral, qwen, gemma, etc.)
	Family string

	// Layers is the total number of transformer blocks (block_count / n_layers)
	Layers int

	// DModel is the embedding dimension (embedding_length / hidden_size)
	DModel int

	// FFNSize is the feed-forward intermediate dimension (feed_forward_length)
	FFNSize int

	// NumHeads is the number of attention heads (attention.head_count)
	NumHeads int

	// NumKVHeads is the number of KV heads (attention.head_count_kv).
	// For standard MHA, NumKVHeads == NumHeads.
	// For GQA/MQA, NumKVHeads < NumHeads.
	NumKVHeads int

	// HeadDim is the dimension per head: DModel / NumHeads.
	// If DModel or NumHeads is 0, HeadDim is inferred from defaults.
	HeadDim int

	// VocabSize is the tokenizer vocabulary size
	VocabSize int

	// ContextLength is the native maximum context from metadata
	ContextLength int

	// KVCacheType describes how KV cache is stored (fp16, fp8, q8_0, etc.)
	KVCacheType string

	// KVCacheBytesPerElement is the size of one KV element in bytes.
	// Default is 2 (fp16). Quantized KV caches may use 1 or less.
	KVCacheBytesPerElement int

	// UseGQA is true when NumKVHeads < NumHeads (grouped-query attention)
	UseGQA bool

	// UseSlidingWindow is true for models with sliding window attention
	UseSlidingWindow bool

	// SlidingWindowSize is the attention window size if applicable
	SlidingWindowSize int

	// IsMoE is true for mixture-of-experts architectures
	IsMoE bool

	// ExpertCount is the total number of experts (for MoE)
	ExpertCount int

	// ExpertUsedCount is the number of active experts per token (for MoE)
	ExpertUsedCount int

	// RoPEDimension is the rotary embedding dimension count
	RoPEDimension int

	// RoPEBase is the rotary frequency base
	RoPEBase float64

	// RoPEScale is the rotary scaling factor
	RoPEScale float64

	// QuantType is the weight quantization (Q4_K_M, IQ4_XS, etc.)
	QuantType string

	// BitsPerWeight is the average bits per weight
	BitsPerWeight float64

	// FileSizeBytes is the actual GGUF file size on disk
	FileSizeBytes int64

	// ParameterCount is the total number of parameters
	ParameterCount uint64

	// HasMetadata is true when all critical fields were read from GGUF.
	// If false, some fields were estimated from parameter count / file size.
	HasMetadata bool

	// MissingFields lists which critical fields were not found in metadata.
	MissingFields []string
}

// Confidence returns how much we trust this spec (0.0–1.0).
// Higher when all critical fields are present in metadata.
func (a *ArchitectureSpec) Confidence() float64 {
	critical := 0
	present := 0
	if a.Layers > 0 {
		critical++
		present++
	}
	if a.DModel > 0 {
		critical++
		present++
	}
	if a.NumHeads > 0 {
		critical++
		present++
	}
	if a.NumKVHeads > 0 {
		critical++
		present++
	}
	if a.FFNSize > 0 {
		critical++
		present++
	}
	if a.ContextLength > 0 {
		critical++
		present++
	}
	if critical == 0 {
		return 0.0
	}
	return float64(present) / float64(critical)
}

// HasCriticalFields returns true when all fields needed for memory
// estimation are present.
func (a *ArchitectureSpec) HasCriticalFields() bool {
	return a.Layers > 0 && a.DModel > 0 && a.NumHeads > 0 && a.NumKVHeads > 0 && a.ContextLength > 0
}
