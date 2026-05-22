package gguf

import (
	"howl-chat/internal/backend/types"

	gguf_parser "github.com/gpustack/gguf-parser-go"
)

// DetectArchitecture extracts a complete ArchitectureSpec from GGUF metadata.
// It is family-agnostic: it reads the architecture prefix from metadata,
// then reads all {family}.* keys automatically.
func (r *Reader) DetectArchitecture(metadata gguf_parser.GGUFMetadataKVs, fileSize int64, paramCount uint64) *types.ArchitectureSpec {
	family := r.getString(metadata, "general.architecture", "llama")
	if family == "" {
		family = "llama"
	}

	spec := &types.ArchitectureSpec{
		Family:                 family,
		FileSizeBytes:          fileSize,
		ParameterCount:         paramCount,
		QuantType:              r.getFileType(metadata),
		BitsPerWeight:          r.getBitsPerWeight(metadata),
		KVCacheType:            "fp16",
		KVCacheBytesPerElement: 2,
	}

	// Read all family-prefixed keys
	spec.Layers = r.getFamilyInt(metadata, family, "block_count", 0)
	if spec.Layers == 0 {
		spec.Layers = r.getFamilyInt(metadata, family, "n_layers", 0)
	}

	spec.DModel = r.getFamilyInt(metadata, family, "embedding_length", 0)
	if spec.DModel == 0 {
		spec.DModel = r.getFamilyInt(metadata, family, "hidden_size", 0)
	}

	spec.FFNSize = r.getFamilyInt(metadata, family, "feed_forward_length", 0)
	if spec.FFNSize == 0 {
		spec.FFNSize = r.getFamilyInt(metadata, family, "intermediate_size", 0)
	}

	spec.NumHeads = r.getFamilyInt(metadata, family, "attention.head_count", 0)
	spec.NumKVHeads = r.getFamilyInt(metadata, family, "attention.head_count_kv", 0)
	if spec.NumKVHeads == 0 {
		// Fallback: some models use key_count, value_count separately
		kvCount := r.getFamilyInt(metadata, family, "attention.head_count_k", 0)
		if kvCount == 0 {
			kvCount = spec.NumHeads // Standard MHA
		}
		spec.NumKVHeads = kvCount
	}

	spec.VocabSize = r.getFamilyInt(metadata, family, "vocab_size", 0)
	if spec.VocabSize == 0 {
		spec.VocabSize = r.getFamilyInt(metadata, family, "token_count", 0)
	}

	spec.ContextLength = r.getFamilyInt(metadata, family, "context_length", 0)
	if spec.ContextLength == 0 {
		spec.ContextLength = r.getFamilyInt(metadata, family, "max_position_embeddings", 0)
	}

	spec.RoPEDimension = r.getFamilyInt(metadata, family, "rope.dimension_count", 0)
	spec.RoPEBase = r.getFamilyFloat(metadata, family, "rope.freq_base", 10000.0)
	spec.RoPEScale = r.getFamilyFloat(metadata, family, "rope.freq_scale", 1.0)

	// Sliding window — try multiple common keys
	spec.SlidingWindowSize = r.getFamilyInt(metadata, family, "attention.sliding_window", 0)
	if spec.SlidingWindowSize == 0 {
		spec.SlidingWindowSize = r.getFamilyInt(metadata, family, "sliding_window", 0)
	}
	if spec.SlidingWindowSize == 0 {
		spec.SlidingWindowSize = r.getFamilyInt(metadata, family, "attention.window_size", 0)
	}
	if spec.SlidingWindowSize == 0 {
		spec.SlidingWindowSize = r.getInt(metadata, "sliding_window", 0)
	}
	if spec.SlidingWindowSize == 0 {
		spec.SlidingWindowSize = r.getInt(metadata, "attention.sliding_window", 0)
	}
	spec.UseSlidingWindow = spec.SlidingWindowSize > 0
	// Family-based inference: gemma models use sliding window attention
	if !spec.UseSlidingWindow && (family == "gemma" || family == "gemma2" || family == "gemma4") {
		spec.UseSlidingWindow = true
		spec.SlidingWindowSize = 4096
	}

	// MoE detection
	spec.ExpertCount = r.getFamilyInt(metadata, family, "expert_count", 0)
	spec.ExpertUsedCount = r.getFamilyInt(metadata, family, "expert_used_count", 0)
	spec.IsMoE = spec.ExpertCount > 0

	// Compute head_dim
	if spec.DModel > 0 && spec.NumHeads > 0 {
		spec.HeadDim = spec.DModel / spec.NumHeads
	} else {
		spec.HeadDim = inferHeadDim(family, paramCount, spec.DModel)
	}

	// Detect GQA
	spec.UseGQA = spec.NumKVHeads > 0 && spec.NumKVHeads < spec.NumHeads

	// Track missing fields
	spec.MissingFields = r.findMissingFields(spec)
	spec.HasMetadata = len(spec.MissingFields) == 0

	// If critical fields missing, estimate from parameter count or file size
	if !spec.HasCriticalFields() {
		if paramCount > 0 {
			r.inferFromParams(spec, paramCount)
		} else if spec.FileSizeBytes > 0 {
			// No metadata param count — estimate from file size
			estimatedParams := estimateParamsFromFileSize(spec.FileSizeBytes, spec.QuantType)
			r.inferFromParams(spec, estimatedParams)
		}
	}

	// Validate: if context is still 0, use family defaults
	if spec.ContextLength == 0 {
		spec.ContextLength = defaultContextLength(family)
	}

	return spec
}

// getFamilyInt reads {family}.{key} from metadata.
func (r *Reader) getFamilyInt(metadata gguf_parser.GGUFMetadataKVs, family, key string, def int) int {
	val := r.getInt(metadata, family+"."+key, 0)
	if val > 0 {
		return val
	}
	// Some models store under just the key without prefix
	return r.getInt(metadata, key, def)
}

// getFamilyFloat reads {family}.{key} from metadata.
func (r *Reader) getFamilyFloat(metadata gguf_parser.GGUFMetadataKVs, family, key string, def float64) float64 {
	val := r.getFloat(metadata, family+"."+key, 0)
	if val > 0 {
		return val
	}
	return r.getFloat(metadata, key, def)
}

// findMissingFields reports which critical fields were not found.
func (r *Reader) findMissingFields(spec *types.ArchitectureSpec) []string {
	var missing []string
	if spec.Layers <= 0 {
		missing = append(missing, "block_count")
	}
	if spec.DModel <= 0 {
		missing = append(missing, "embedding_length")
	}
	if spec.NumHeads <= 0 {
		missing = append(missing, "attention.head_count")
	}
	if spec.NumKVHeads <= 0 {
		missing = append(missing, "attention.head_count_kv")
	}
	if spec.FFNSize <= 0 {
		missing = append(missing, "feed_forward_length")
	}
	if spec.ContextLength <= 0 {
		missing = append(missing, "context_length")
	}
	return missing
}

// inferFromParams estimates missing architecture fields from parameter count.
func (r *Reader) inferFromParams(spec *types.ArchitectureSpec, paramCount uint64) {
	b := float64(paramCount)
	if spec.Layers <= 0 {
		spec.Layers = EstimateLayers(b)
	}
	if spec.DModel <= 0 {
		spec.DModel = EstimateDModel(b)
	}
	if spec.NumHeads <= 0 {
		spec.NumHeads = EstimateNumHeads(spec.DModel)
	}
	if spec.NumKVHeads <= 0 {
		spec.NumKVHeads = spec.NumHeads // assume MHA if unknown
	}
	if spec.FFNSize <= 0 {
		spec.FFNSize = EstimateFFNSize(spec.DModel)
	}
	if spec.HeadDim <= 0 && spec.DModel > 0 && spec.NumHeads > 0 {
		spec.HeadDim = spec.DModel / spec.NumHeads
	}
}

// EstimateLayers infers layer count from parameter count (billions).
func EstimateLayers(params float64) int {
	switch {
	case params < 1_000_000_000:
		return 24
	case params < 3_000_000_000:
		return 28
	case params < 7_000_000_000:
		return 32
	case params < 15_000_000_000:
		return 40
	case params < 40_000_000_000:
		return 48
	case params < 80_000_000_000:
		return 64
	default:
		return 80
	}
}

// estimateDModel infers embedding dimension from parameter count.
func EstimateDModel(params float64) int {
	switch {
	case params < 1_000_000_000:
		return 2048
	case params < 3_000_000_000:
		return 2560
	case params < 7_000_000_000:
		return 4096
	case params < 15_000_000_000:
		return 5120
	case params < 40_000_000_000:
		return 6144
	case params < 80_000_000_000:
		return 8192
	default:
		return 8192
	}
}

// estimateNumHeads infers head count from d_model.
func EstimateNumHeads(dModel int) int {
	switch {
	case dModel <= 2048:
		return 32
	case dModel <= 4096:
		return 32
	case dModel <= 5120:
		return 40
	case dModel <= 6144:
		return 48
	case dModel <= 8192:
		return 64
	default:
		return 64
	}
}

// estimateFFNSize infers intermediate dimension from d_model.
func EstimateFFNSize(dModel int) int {
	// Most architectures use 4x or ~2.7x (SwiGLU)
	if dModel <= 0 {
		return 0
	}
	// SwiGLU uses ~2/3 * 4 * d_model ≈ 2.7 * d_model
	// But for memory estimation, use 4x as conservative upper bound
	return dModel * 4
}

// inferHeadDim returns a sensible head_dim when we can't compute DModel/NumHeads.
func inferHeadDim(family string, paramCount uint64, dModel int) int {
	if dModel > 0 {
		switch {
		case dModel <= 2048:
			return 64
		case dModel <= 2560:
			return 80
		case dModel <= 4096:
			return 128
		case dModel <= 5120:
			return 128
		default:
			return 128
		}
	}
	// Fallback by family
	switch family {
	case "gemma", "gemma2", "gemma3", "gemma4":
		return 256
	case "qwen":
		return 128
	default:
		return 128
	}
}

// defaultContextLength returns the typical context length for a model family.
func defaultContextLength(family string) int {
	switch family {
	case "llama3":
		return 128000
	case "llama":
		return 4096
	case "mistral":
		return 32768
	case "qwen":
		return 32768
	case "gemma", "gemma2", "gemma3", "gemma4":
		return 8192
	case "phi3":
		return 128000
	default:
		return 4096
	}
}

// KVCacheMemoryMB returns the KV cache memory in MB for the given context size.
func KVCacheMemoryMB(spec *types.ArchitectureSpec, contextSize int) float64 {
	if spec == nil || spec.Layers <= 0 || spec.HeadDim <= 0 || contextSize <= 0 {
		return 0
	}

	// KV cache: for each layer, for each token:
	//   2 tensors (K, V) × n_kv_heads × head_dim × bytes_per_element
	bytesPerTokenPerLayer := 2 * int64(spec.NumKVHeads) * int64(spec.HeadDim) * int64(spec.KVCacheBytesPerElement)
	totalBytes := bytesPerTokenPerLayer * int64(spec.Layers) * int64(contextSize)
	return float64(totalBytes) / (1024 * 1024)
}

// WeightMemoryMB returns the weight memory estimate.
// For mmap, this is essentially the file size (the OS maps it directly).
func WeightMemoryMB(spec *types.ArchitectureSpec) float64 {
	if spec == nil || spec.FileSizeBytes <= 0 {
		return 0
	}
	return float64(spec.FileSizeBytes) / (1024 * 1024)
}

// ActivationMemoryMB estimates temporary activation memory for a single forward pass.
func ActivationMemoryMB(spec *types.ArchitectureSpec, batchSize int) float64 {
	if spec == nil || spec.DModel <= 0 {
		return 0
	}
	if batchSize <= 0 {
		batchSize = 512
	}

	// Activation memory ≈ (d_model + ffn_size) * batch_size * 4 bytes (fp32 intermediate)
	// This is a rough upper bound; actual usage depends on the compute graph.
	activationBytes := int64(spec.DModel+spec.FFNSize) * int64(batchSize) * 4
	overhead := int64(100 * 1024 * 1024) // 100MB base
	return float64(activationBytes+overhead) / (1024 * 1024)
}
