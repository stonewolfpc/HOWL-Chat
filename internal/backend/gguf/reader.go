/**
 * GGUF Metadata Reader
 *
 * This package provides functionality for reading GGUF model metadata
 * and building model profiles for auto-configuration.
 *
 * @package gguf
 */

package gguf

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"howl-chat/internal/backend/types"

	gguf_parser "github.com/gpustack/gguf-parser-go"
)

// Reader handles GGUF metadata extraction and profile building
type Reader struct{}

// NewReader creates a new GGUF metadata reader
func NewReader() *Reader {
	return &Reader{}
}

// ReadProfile reads GGUF metadata from a file and builds a ModelProfile
func (r *Reader) ReadProfile(path string) (*types.ModelProfile, error) {
	// Parse GGUF file using gpustack/gguf-parser-go
	gguf, err := gguf_parser.ParseGGUFFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GGUF file: %w", err)
	}

	// Build profile from metadata
	profile := &types.ModelProfile{
		Path:           path,
		Family:         r.getString(gguf.Header.MetadataKV, "general.architecture", "llama"),
		Variant:        r.getString(gguf.Header.MetadataKV, "general.name", ""),
		MaxContext:     r.getInt(gguf.Header.MetadataKV, "llama.context_length", 4096),
		RopeMode:       r.getString(gguf.Header.MetadataKV, "llama.rope.scaling", "none"),
		RopeFactor:     r.getFloat(gguf.Header.MetadataKV, "llama.rope.freq_scale", 1.0),
		RopeBase:       r.getFloat(gguf.Header.MetadataKV, "llama.rope.freq_base", 10000.0),
		Template:       r.resolveTemplate(gguf.Header.MetadataKV),
		Tokenizer:      r.getString(gguf.Header.MetadataKV, "tokenizer.ggml.model", ""),
		StopSequences:  r.resolveStopSequences(gguf.Header.MetadataKV),
		LayerCount:     r.getLayerCount(gguf.Header.MetadataKV),
		ParameterCount: r.getParamCount(gguf.Header.MetadataKV),
		FileType:       r.getFileType(gguf.Header.MetadataKV),
		BitsPerWeight:  r.getBitsPerWeight(gguf.Header.MetadataKV),
	}

	// Read actual file size
	if fi, err := os.Stat(path); err == nil {
		profile.FileSizeBytes = fi.Size()
	}

	// Fallback: estimate parameter count from file size when metadata is missing
	if profile.ParameterCount == 0 && profile.FileSizeBytes > 0 {
		profile.ParameterCount = estimateParamsFromFileSize(profile.FileSizeBytes, profile.FileType)
	}

	// Fallback: try to detect quantization from filename if metadata is missing
	if profile.FileType == "" {
		profile.FileType = detectQuantFromFilename(path)
		if profile.BitsPerWeight <= 0 {
			profile.BitsPerWeight = bitsForQuant(profile.FileType)
		}
	}

	// Try alternative context length fields if llama.context_length is 0
	if profile.MaxContext == 0 {
		profile.MaxContext = r.getInt(gguf.Header.MetadataKV, "mistral.context_length", 4096)
	}
	if profile.MaxContext == 0 {
		profile.MaxContext = r.getInt(gguf.Header.MetadataKV, "qwen.context_length", 4096)
	}

	// Compute usable context with safety margins
	profile.UsableContext = computeUsableContext(profile.MaxContext)

	// Detect full architecture spec (the authoritative source for memory estimation)
	profile.ArchitectureSpec = r.DetectArchitecture(
		gguf.Header.MetadataKV,
		profile.FileSizeBytes,
		profile.ParameterCount,
	)
	// Sync legacy fields from architecture spec
	profile.LayerCount = profile.ArchitectureSpec.Layers
	profile.Family = profile.ArchitectureSpec.Family

	fmt.Printf("[ReadProfile] final: Family=%s Layers=%d Params=%d Size=%.1fMB Quant=%s DModel=%d Heads=%d KVHeads=%d\n",
		profile.Family, profile.LayerCount, profile.ParameterCount,
		float64(profile.FileSizeBytes)/(1024*1024), profile.FileType,
		profile.ArchitectureSpec.DModel, profile.ArchitectureSpec.NumHeads,
		profile.ArchitectureSpec.NumKVHeads)

	return profile, nil
}

// getString extracts a string value from metadata with a default fallback
func (r *Reader) getString(metadata gguf_parser.GGUFMetadataKVs, key, defaultValue string) string {
	for _, kv := range metadata {
		if kv.Key == key {
			if s, ok := kv.Value.(string); ok {
				return s
			}
		}
	}
	return defaultValue
}

// getInt extracts an integer value from metadata with a default fallback.
// Handles all numeric types the GGUF parser may emit.
func (r *Reader) getInt(metadata gguf_parser.GGUFMetadataKVs, key string, defaultValue int) int {
	for _, kv := range metadata {
		if kv.Key == key {
			switch t := kv.Value.(type) {
			case int:
				return t
			case int8:
				return int(t)
			case int16:
				return int(t)
			case int32:
				return int(t)
			case int64:
				return int(t)
			case uint:
				return int(t)
			case uint8:
				return int(t)
			case uint16:
				return int(t)
			case uint32:
				return int(t)
			case uint64:
				return int(t)
			case float32:
				return int(t)
			case float64:
				return int(t)
			case string:
				if i, err := strconv.Atoi(t); err == nil {
					return i
				}
			}
		}
	}
	return defaultValue
}

// getFloat extracts a float value from metadata with a default fallback.
// Handles all numeric types the GGUF parser may emit.
func (r *Reader) getFloat(metadata gguf_parser.GGUFMetadataKVs, key string, defaultValue float64) float64 {
	for _, kv := range metadata {
		if kv.Key == key {
			switch t := kv.Value.(type) {
			case float32:
				return float64(t)
			case float64:
				return t
			case int:
				return float64(t)
			case int8:
				return float64(t)
			case int16:
				return float64(t)
			case int32:
				return float64(t)
			case int64:
				return float64(t)
			case uint:
				return float64(t)
			case uint8:
				return float64(t)
			case uint16:
				return float64(t)
			case uint32:
				return float64(t)
			case uint64:
				return float64(t)
			case string:
				if f, err := strconv.ParseFloat(t, 64); err == nil {
					return f
				}
			}
		}
	}
	return defaultValue
}

// getLayerCount extracts the total layer count from GGUF metadata
func (r *Reader) getLayerCount(metadata gguf_parser.GGUFMetadataKVs) int {
	// Try common layer count metadata keys
	layerKeys := []string{
		"llama.block_count",
		"llama.n_layers",
		"general.block_count",
		"general.n_layers",
		"mistral.block_count",
		"mistral.n_layers",
		"qwen.block_count",
		"qwen.n_layers",
		"gemma.block_count",
		"gemma.n_layers",
		"gemma2.block_count",
		"gemma2.n_layers",
		"gemma4.block_count",
		"gemma4.n_layers",
		"gemma3.block_count",
		"gemma3.n_layers",
		"phi3.block_count",
		"phi3.n_layers",
		"falcon.block_count",
		"falcon.n_layers",
		"gpt_neox.block_count",
		"gpt_neox.n_layers",
		"cohere.block_count",
		"cohere.n_layers",
	}

	for _, key := range layerKeys {
		if count := r.getInt(metadata, key, 0); count > 0 {
			return count
		}
	}

	// Fallback: estimate from parameter count for dense transformers
	paramCount := r.getParamCount(metadata)
	if paramCount > 0 {
		// Typical dense transformer: ~0.1B params per layer (rough)
		estLayers := int(paramCount / 100_000_000)
		if estLayers < 1 {
			estLayers = 1
		}
		if estLayers > 80 {
			estLayers = 80
		}
		return estLayers
	}
	// Default to 0 if not found (will be handled by caller)
	return 0
}

// resolveTemplate extracts the chat template from metadata with fallbacks
func (r *Reader) resolveTemplate(metadata gguf_parser.GGUFMetadataKVs) string {
	// Priority order for template keys
	templateKeys := []string{
		"chat_template",
		"lfm2.chat_template",
		"mistral.chat_template",
		"qwen.chat_template",
	}

	for _, key := range templateKeys {
		if v := r.getString(metadata, key, ""); v != "" {
			// Convert Handlebars to Jinja2 if needed
			return ConvertHandlebarsToJinja(v)
		}
	}

	// Fallback by family
	family := r.getString(metadata, "general.architecture", "llama")
	return ConvertHandlebarsToJinja(r.getFallbackTemplate(family))
}

// getFallbackTemplate returns a default template based on model family
func (r *Reader) getFallbackTemplate(family string) string {
	switch family {
	case "llama":
		return "{{ system }}\n{{#each messages}}{{ role }}: {{ content }}\n{{/each}}assistant:"
	case "mistral":
		return "<s>[INST] {{ system }} {{#each messages}}{{ role }}: {{ content }} [/INST]\n{{/each}}"
	case "qwen":
		return "<|im_start|>system\n{{ system }}<|im_end|>\n{{#each messages}}<|im_start|>{{ role }}\n{{ content }}<|im_end|>\n{{/each}}<|im_start|>assistant"
	default:
		return "{{#each messages}}{{ role }}: {{ content }}\n{{/each}}assistant:"
	}
}

// resolveStopSequences extracts stop sequences from metadata
func (r *Reader) resolveStopSequences(metadata gguf_parser.GGUFMetadataKVs) []string {
	// Try to extract stop sequences from tokenizer.stop or similar metadata keys
	stopKeys := []string{
		"tokenizer.stop",
		"tokenizer.stop_strings",
		"stop",
		"stop_strings",
	}

	for _, key := range stopKeys {
		if v := r.getString(metadata, key, ""); v != "" {
			// If the value is a string, try to parse it as a comma-separated list
			if strings.Contains(v, ",") {
				return strings.Split(v, ",")
			}
			// If it's a single string, return it as a single-element slice
			return []string{v}
		}
	}

	// Try to extract from array-type metadata
	for _, key := range stopKeys {
		for _, kv := range metadata {
			if kv.Key == key {
				if arr, ok := kv.Value.([]any); ok {
					sequences := make([]string, 0, len(arr))
					for _, item := range arr {
						if s, ok := item.(string); ok {
							sequences = append(sequences, s)
						}
					}
					if len(sequences) > 0 {
						return sequences
					}
				}
			}
		}
	}

	// Fallback to common stop sequences for different model families
	family := r.getString(metadata, "general.architecture", "llama")
	return r.getFallbackStopSequences(family)
}

// getFallbackStopSequences returns common stop sequences for different model families
func (r *Reader) getFallbackStopSequences(family string) []string {
	switch family {
	case "llama", "llama2", "llama3":
		return []string{"<|endoftext|>", "<|im_end|>", "<|eot_id|>"}
	case "mistral":
		return []string{"</s>", "[INST]", "[/INST]"}
	case "qwen":
		return []string{"<|im_end|>", "<|endoftext|>"}
	default:
		return []string{"<|endoftext|>"}
	}
}

// getParamCount extracts the parameter count from metadata.
// Handles all numeric types the GGUF parser may emit.
func (r *Reader) getParamCount(metadata gguf_parser.GGUFMetadataKVs) uint64 {
	for _, key := range []string{
		"general.parameter_count",
		"general.params",
		"model.params",
		"general.parameter_count_trainable",
		"general.count_params",
	} {
		for _, kv := range metadata {
			if kv.Key == key {
				switch t := kv.Value.(type) {
				case uint64:
					return t
				case uint32:
					return uint64(t)
				case uint16:
					return uint64(t)
				case uint8:
					return uint64(t)
				case uint:
					return uint64(t)
				case int:
					return uint64(t)
				case int8:
					return uint64(t)
				case int16:
					return uint64(t)
				case int32:
					return uint64(t)
				case int64:
					return uint64(t)
				case float32:
					return uint64(t)
				case float64:
					return uint64(t)
				case string:
					if v, err := strconv.ParseUint(t, 10, 64); err == nil {
						return v
					}
				}
			}
		}
	}
	return 0
}

// getFileType extracts the quantization type from metadata
func (r *Reader) getFileType(metadata gguf_parser.GGUFMetadataKVs) string {
	ft := r.getInt(metadata, "general.file_type", -1)
	if ft < 0 {
		return ""
	}
	// GGUF file_type values from llama.cpp
	quantTypes := map[int]string{
		1: "Q4_0", 2: "Q4_1", 3: "Q4_2", 4: "Q4_3", 5: "Q5_0", 6: "Q5_1",
		7: "Q8_0", 8: "Q8_1", 9: "Q2_K", 10: "Q3_K", 11: "Q4_K", 12: "Q5_K",
		13: "Q6_K", 14: "Q8_K", 15: "IQ2_XXS", 16: "IQ2_XS", 17: "IQ2_S",
		18: "IQ3_XXS", 19: "IQ1_S", 20: "IQ4_NL", 21: "IQ3_S", 22: "IQ3_M",
		23: "IQ2_M", 24: "IQ1_M", 25: "IQ4_XS", 26: "IQ4_NL", 27: "IQ3_XXL",
		28: "IQ4_XS_R", 29: "BF16", 30: "Q4_K_R", 31: "Q6_K_R",
	}
	if name, ok := quantTypes[ft]; ok {
		return name
	}
	return fmt.Sprintf("FT%d", ft)
}

// getBitsPerWeight returns the average bits per weight for a quantization type
func (r *Reader) getBitsPerWeight(metadata gguf_parser.GGUFMetadataKVs) float64 {
	ft := r.getInt(metadata, "general.file_type", -1)
	switch ft {
	case 1:
		return 4.5 // Q4_0
	case 2:
		return 5.0 // Q4_1
	case 5:
		return 5.5 // Q5_0
	case 6:
		return 6.0 // Q5_1
	case 7:
		return 8.5 // Q8_0
	case 8:
		return 9.0 // Q8_1
	case 9:
		return 2.625 // Q2_K
	case 10:
		return 3.4375 // Q3_K
	case 11:
		return 4.5 // Q4_K
	case 12:
		return 5.5 // Q5_K
	case 13:
		return 6.5625 // Q6_K
	case 14:
		return 8.5 // Q8_K
	case 15:
		return 2.375 // IQ2_XXS
	case 16:
		return 2.5 // IQ2_XS
	case 17:
		return 2.75 // IQ2_S
	case 18:
		return 3.0625 // IQ3_XXS
	case 19:
		return 1.5625 // IQ1_S
	case 20, 26:
		return 4.5 // IQ4_NL
	case 21:
		return 3.4375 // IQ3_S
	case 22:
		return 3.5 // IQ3_M
	case 23:
		return 2.875 // IQ2_M
	case 24:
		return 1.75 // IQ1_M
	case 25, 28:
		return 4.25 // IQ4_XS / IQ4_XS_R
	case 27:
		return 3.0 // IQ3_XXL
	case 29:
		return 16.0 // BF16
	case 30:
		return 4.5 // Q4_K_R
	case 31:
		return 6.5625 // Q6_K_R
	default:
		// Fallback: if we have file size and parameter count, derive from those
		return 0
	}
}

// computeUsableContext calculates safe context window with safety margins
// Following LM Studio's approach to avoid slowdown and hallucination
func computeUsableContext(max int) int {
	if max <= 0 {
		return 4096
	}

	margin := 1000
	switch {
	case max <= 32000:
		margin = 1000
	case max <= 64000:
		margin = 2000
	case max <= 128000:
		margin = 3000
	default:
		margin = 4000
	}

	// Ensure we don't go below minimum safe context
	if max-margin < 1024 {
		return max
	}

	return max - margin
}

// EstimateMemory computes estimated GPU and RAM requirements for a model load.
// It delegates to the architecture-based estimator for accuracy.
func EstimateMemory(profile *types.ModelProfile, gpuLayers, contextSize int) types.MemoryEstimate {
	if profile == nil {
		return types.MemoryEstimate{Notes: "no profile available"}
	}

	// Prefer the new architecture-aware estimator
	if profile.ArchitectureSpec != nil && profile.ArchitectureSpec.HasCriticalFields() {
		model := NewMemoryModel(profile.ArchitectureSpec)
		return model.Estimate(contextSize, gpuLayers, 512)
	}

	// Legacy fallback when ArchitectureSpec is unavailable
	fileSizeMB := float64(profile.FileSizeBytes) / (1024 * 1024)
	if fileSizeMB <= 0 {
		if profile.ParameterCount > 0 && profile.BitsPerWeight > 0 {
			fileSizeMB = float64(profile.ParameterCount) * profile.BitsPerWeight / 8.0 / (1024 * 1024)
		} else {
			fileSizeMB = 4000
		}
	}

	const overheadMB = 300
	totalLayers := profile.LayerCount
	if totalLayers <= 0 {
		totalLayers = 32
	}
	perLayerWeightMB := fileSizeMB / float64(totalLayers+1)

	nEmbd := 4096
	switch {
	case profile.ParameterCount < 3_000_000_000:
		nEmbd = 2048
	case profile.ParameterCount < 8_000_000_000:
		nEmbd = 4096
	case profile.ParameterCount < 20_000_000_000:
		nEmbd = 5120
	default:
		nEmbd = 8192
	}

	kvPerTokenPerLayer := float64(nEmbd) * 4.0
	kvCacheTotalMB := float64(contextSize) * kvPerTokenPerLayer * float64(totalLayers) / (1024 * 1024)

	gpuWeightMB := perLayerWeightMB * float64(gpuLayers)
	gpuKVShare := 0.0
	if gpuLayers > 0 {
		gpuKVShare = kvCacheTotalMB * (float64(gpuLayers) / float64(totalLayers))
	}
	gpuTotalMB := gpuWeightMB + gpuKVShare + overheadMB*0.5

	ramWeightMB := fileSizeMB - gpuWeightMB
	if ramWeightMB < 0 {
		ramWeightMB = 0
	}
	ramKVShare := kvCacheTotalMB - gpuKVShare
	ramTotalMB := ramWeightMB + ramKVShare + overheadMB

	return types.MemoryEstimate{
		GPUMB:           int64(gpuTotalMB),
		RAMMB:           int64(ramTotalMB),
		TotalMB:         int64(gpuTotalMB + ramTotalMB),
		ContextSizeUsed: contextSize,
		GPULayersUsed:   gpuLayers,
		Notes:           fmt.Sprintf("LEGACY: file=%.0fMB layers=%d gpu_layers=%d ctx=%d n_embd=%d", fileSizeMB, totalLayers, gpuLayers, contextSize, nEmbd),
	}
}

// ExtractModelProfile reads a GGUF file and extracts the model profile including chat template
func ExtractModelProfile(path string) (*types.ModelProfile, error) {
	reader := NewReader()
	return reader.ReadProfile(path)
}

// ConvertHandlebarsToJinja converts Handlebars-style templates to Jinja2 syntax
// Some GGUF files use Handlebars ({{#each}}) instead of Jinja2 ({%- for %})
func ConvertHandlebarsToJinja(template string) string {
	// Check if this is Handlebars syntax
	if !strings.Contains(template, "{{#each") && !strings.Contains(template, "{{/each") {
		// Already Jinja2 or no template
		return template
	}

	// Convert Handlebars to Jinja2
	jinja := template

	// Replace {{#each messages}} with {%- for message in messages %}
	jinja = strings.ReplaceAll(jinja, "{{#each messages}}", "{%- for message in messages %}")

	// Replace {{/each}} with {%- endfor %}
	jinja = strings.ReplaceAll(jinja, "{{/each}}", "{%- endfor %}")

	// Replace {{ role }} with {{ message.role }}
	jinja = strings.ReplaceAll(jinja, "{{ role }}", "{{ message.role }}")

	// Replace {{ content }} with {{ message.content }}
	jinja = strings.ReplaceAll(jinja, "{{ content }}", "{{ message.content }}")

	return jinja
}

// estimateParamsFromFileSize guesses parameter count from file size when metadata is absent.
// It assumes typical quantization overhead and uses conservative estimates.
func estimateParamsFromFileSize(fileSizeBytes int64, quantType string) uint64 {
	fileSizeMB := float64(fileSizeBytes) / (1024 * 1024)

	// Estimate bits per weight from quant type, default to 4-bit guess
	bitsPerWeight := 4.0
	switch {
	case strings.Contains(quantType, "Q8") || strings.Contains(quantType, "BF16"):
		bitsPerWeight = 16.0 // BF16 or Q8 with overhead
	case strings.Contains(quantType, "Q6"):
		bitsPerWeight = 6.5
	case strings.Contains(quantType, "Q5"):
		bitsPerWeight = 5.5
	case strings.Contains(quantType, "Q4") || strings.Contains(quantType, "IQ4"):
		bitsPerWeight = 4.5
	case strings.Contains(quantType, "Q3") || strings.Contains(quantType, "IQ3"):
		bitsPerWeight = 3.5
	case strings.Contains(quantType, "Q2") || strings.Contains(quantType, "IQ2"):
		bitsPerWeight = 2.6
	case strings.Contains(quantType, "FP16") || strings.Contains(quantType, "F16"):
		bitsPerWeight = 16.0
	case strings.Contains(quantType, "FP32") || strings.Contains(quantType, "F32"):
		bitsPerWeight = 32.0
	}

	// params = file_size_bytes * 8 / bits_per_weight
	params := fileSizeMB * 8.0 / bitsPerWeight * 1_000_000

	// Apply overhead factor (GGUF has metadata, vocab, tensors overhead ~15-25%)
	params = params / 1.2

	return uint64(params)
}

// detectQuantFromFilename extracts quantization type from common filename patterns.
func detectQuantFromFilename(path string) string {
	lower := strings.ToLower(path)
	patterns := []struct {
		pat  string
		name string
	}{
		{"IQ4_XS", "IQ4_XS"}, {"IQ4_NL", "IQ4_NL"}, {"IQ4_XS_R", "IQ4_XS_R"},
		{"IQ3_XXL", "IQ3_XXL"}, {"IQ3_XXS", "IQ3_XXS"}, {"IQ3_S", "IQ3_S"}, {"IQ3_M", "IQ3_M"},
		{"IQ2_XXS", "IQ2_XXS"}, {"IQ2_XS", "IQ2_XS"}, {"IQ2_S", "IQ2_S"}, {"IQ2_M", "IQ2_M"},
		{"IQ1_S", "IQ1_S"}, {"IQ1_M", "IQ1_M"},
		{"Q4_K_M", "Q4_K_M"}, {"Q4_K_S", "Q4_K_S"}, {"Q4_K_R", "Q4_K_R"},
		{"Q5_K_M", "Q5_K_M"}, {"Q5_K_S", "Q5_K_S"},
		{"Q6_K", "Q6_K"}, {"Q6_K_R", "Q6_K_R"},
		{"Q8_0", "Q8_0"}, {"Q8_K", "Q8_K"},
		{"Q4_0", "Q4_0"}, {"Q4_1", "Q4_1"},
		{"Q5_0", "Q5_0"}, {"Q5_1", "Q5_1"},
		{"Q2_K", "Q2_K"}, {"Q3_K", "Q3_K"},
		{"BF16", "BF16"}, {"FP16", "FP16"}, {"F16", "FP16"},
		{"FP32", "FP32"}, {"F32", "FP32"},
	}
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p.pat)) {
			return p.name
		}
	}
	return ""
}

// bitsForQuant returns approximate bits per weight for a quantization name.
func bitsForQuant(quant string) float64 {
	switch {
	case strings.Contains(quant, "Q8"):
		return 8.5
	case strings.Contains(quant, "Q6"):
		return 6.5
	case strings.Contains(quant, "Q5"):
		return 5.5
	case strings.Contains(quant, "Q4") || strings.Contains(quant, "IQ4"):
		return 4.5
	case strings.Contains(quant, "Q3") || strings.Contains(quant, "IQ3"):
		return 3.5
	case strings.Contains(quant, "Q2") || strings.Contains(quant, "IQ2"):
		return 2.6
	case strings.Contains(quant, "BF16") || strings.Contains(quant, "FP16") || strings.Contains(quant, "F16"):
		return 16.0
	case strings.Contains(quant, "FP32") || strings.Contains(quant, "F32"):
		return 32.0
	default:
		return 4.0
	}
}
