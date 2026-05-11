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
		Path:          path,
		Family:        r.getString(gguf.Header.MetadataKV, "general.architecture", "llama"),
		Variant:       r.getString(gguf.Header.MetadataKV, "general.name", ""),
		MaxContext:    r.getInt(gguf.Header.MetadataKV, "llama.context_length", 4096),
		RopeMode:      r.getString(gguf.Header.MetadataKV, "llama.rope.scaling", "none"),
		RopeFactor:    r.getFloat(gguf.Header.MetadataKV, "llama.rope.freq_scale", 1.0),
		RopeBase:      r.getFloat(gguf.Header.MetadataKV, "llama.rope.freq_base", 10000.0),
		Template:      r.resolveTemplate(gguf.Header.MetadataKV),
		Tokenizer:     r.getString(gguf.Header.MetadataKV, "tokenizer.ggml.model", ""),
		StopSequences: r.resolveStopSequences(gguf.Header.MetadataKV),
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

// getInt extracts an integer value from metadata with a default fallback
func (r *Reader) getInt(metadata gguf_parser.GGUFMetadataKVs, key string, defaultValue int) int {
	for _, kv := range metadata {
		if kv.Key == key {
			switch t := kv.Value.(type) {
			case int:
				return t
			case int64:
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

// getFloat extracts a float value from metadata with a default fallback
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
			case int64:
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
			return v
		}
	}

	// Fallback by family
	family := r.getString(metadata, "general.architecture", "llama")
	return r.getFallbackTemplate(family)
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
