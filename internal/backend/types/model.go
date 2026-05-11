/**
 * Model Types
 * 
 * This package defines all model-related types for the chat system.
 * Models represent LLM models that can be loaded and used for inference.
 * 
 * @package types
 */

package types

import "time"

// ModelStatus represents the current state of a model
type ModelStatus string

const (
	// ModelStatusUnloaded represents a model that is not loaded
	ModelStatusUnloaded ModelStatus = "unloaded"
	// ModelStatusLoading represents a model that is currently loading
	ModelStatusLoading ModelStatus = "loading"
	// ModelStatusLoaded represents a model that is successfully loaded
	ModelStatusLoaded ModelStatus = "loaded"
	// ModelStatusError represents a model that failed to load
	ModelStatusError ModelStatus = "error"
)

// Model represents a loaded or loadable LLM model
type Model struct {
	ID          string      `json:"id"`          // Unique model identifier
	Name        string      `json:"name"`        // Human-readable model name
	FilePath    string      `json:"filepath"`    // Path to the model file
	Format      ModelFormat `json:"format"`      // Model file format (gguf/ggml)
	Size        int64       `json:"size"`        // Model file size in bytes
	Status      ModelStatus `json:"status"`      // Current model status
	LoadedAt    time.Time   `json:"loaded_at"`   // When the model was loaded
	Parameters  ModelParams `json:"parameters"`  // Model parameters
	Metadata    ModelMeta   `json:"metadata"`    // Model metadata
}

// ModelFormat represents the format of a model file
type ModelFormat string

const (
	// FormatGGUF represents the GGUF format (current standard)
	FormatGGUF ModelFormat = "gguf"
	// FormatGGML represents the legacy GGML format
	FormatGGML ModelFormat = "ggml"
)

// ModelParams contains model configuration parameters
type ModelParams struct {
	ContextSize      int     `json:"context_size"`       // Context window size
	EmbeddingSize    int     `json:"embedding_size"`     // Embedding dimension
	FeedForwardSize  int     `json:"feed_forward_size"`  // Feed-forward network size
	AttentionHeads   int     `json:"attention_heads"`    // Number of attention heads
	Layers           int     `json:"layers"`             // Number of transformer layers
	Quantization     string  `json:"quantization"`       // Quantization level (q4_0, q4_k, etc.)
	VocabSize        int     `json:"vocab_size"`         // Vocabulary size
	RopeDimensions   int     `json:"rope_dimensions"`    // RoPE dimensionality
	MaxSequence      int     `json:"max_sequence"`       // Maximum sequence length
}

// ModelMeta contains additional model metadata
type ModelMeta struct {
	Architecture string    `json:"architecture"` // Model architecture (llama, mistral, etc.)
	License      string    `json:"license"`       // Model license
	Version      string    `json:"version"`       // Model version
	Author       string    `json:"author"`        // Model author/creator
	Description  string    `json:"description"`   // Model description
	Tags         []string  `json:"tags"`          // Model tags for categorization
	LastModified time.Time `json:"last_modified"` // Last modification time
}

// ModelConfig represents configuration for model loading
type ModelConfig struct {
	GPULayers      int     `json:"gpu_layers"`       // Number of layers to offload to GPU
	Threads        int     `json:"threads"`          // Number of CPU threads to use
	UseMMap        bool    `json:"use_mmap"`         // Whether to use memory mapping
	UseMlock       bool    `json:"use_mlock"`        // Whether to use mlock for pinned memory
	VocabOnly      bool    `json:"vocab_only"`       // Whether to load only vocabulary
	RMSNormEPS     float64 `json:"rms_norm_eps"`     // RMS normalization epsilon
	RopeFreqBase   float64 `json:"rope_freq_base"`   // RoPE frequency base
	RopeFreqScale  float64 `json:"rope_freq_scale"`  // RoPE frequency scale
	MemoryFraction float64 `json:"memory_fraction"` // Fraction of memory to use
}

// ModelInfo represents information about a model file
type ModelInfo struct {
	FilePath    string      `json:"filepath"`    // Path to the model file
	Format      ModelFormat `json:"format"`      // Model format
	Size        int64       `json:"size"`        // File size in bytes
	Modified    time.Time   `json:"modified"`    // Last modification time
	Checksum    string      `json:"checksum"`    // File checksum (SHA256)
	Compatible  bool        `json:"compatible"`  // Whether the model is compatible
}

// ModelList represents a collection of models
type ModelList struct {
	Models []Model `json:"models"` // List of models
	Total  int     `json:"total"`  // Total number of models
}

// NewModel creates a new model instance
func NewModel(name, filePath string, format ModelFormat, size int64) *Model {
	return &Model{
		ID:       generateModelID(),
		Name:     name,
		FilePath: filePath,
		Format:   format,
		Size:     size,
		Status:   ModelStatusUnloaded,
		LoadedAt: time.Time{},
		Parameters: ModelParams{
			ContextSize:     2048,
			EmbeddingSize:   0,
			FeedForwardSize: 0,
			AttentionHeads:  0,
			Layers:          0,
			Quantization:    "",
			VocabSize:       0,
			RopeDimensions:  0,
			MaxSequence:     0,
		},
		Metadata: ModelMeta{
			Architecture: "",
			License:      "",
			Version:      "",
			Author:       "",
			Description:  "",
			Tags:         []string{},
			LastModified: time.Now(),
		},
	}
}

// DefaultModelConfig returns default model loading configuration
func DefaultModelConfig() *ModelConfig {
	return &ModelConfig{
		GPULayers:      0,
		Threads:        4,
		UseMMap:        true,
		UseMlock:       false,
		VocabOnly:      false,
		RMSNormEPS:     1e-5,
		RopeFreqBase:   10000.0,
		RopeFreqScale:  1.0,
		MemoryFraction: 0.9,
	}
}

// generateModelID generates a unique identifier for a model
func generateModelID() string {
	return "model_" + time.Now().Format("20060102150405")
}

// IsLoaded returns true if the model is currently loaded
func (m *Model) IsLoaded() bool {
	return m.Status == ModelStatusLoaded
}

// IsLoading returns true if the model is currently loading
func (m *Model) IsLoading() bool {
	return m.Status == ModelStatusLoading
}

// HasError returns true if the model is in an error state
func (m *Model) HasError() bool {
	return m.Status == ModelStatusError
}

// GenerateImagePayload represents the payload for image generation from the UI
type GenerateImagePayload struct {
	Prompt       string  `json:"prompt"`
	NegativePrompt string  `json:"negative_prompt"`
	Model        string  `json:"model"`
	Seed         int     `json:"seed"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
	Steps        int     `json:"steps"`
	Cfg          float32 `json:"cfg"`
	Style        string  `json:"style"`
	Character    string  `json:"character"`
	World        string  `json:"world"`
	Scenario     string  `json:"scenario"`
}

// ImageMetadata represents the metadata for a generated image
type ImageMetadata struct {
	File        string    `json:"file"`
	Seed        int       `json:"seed"`
	Prompt      string    `json:"prompt"`
	NegativePrompt string    `json:"negative_prompt"`
	Model       string    `json:"model"`
	Character   string    `json:"character"`
	World       string    `json:"world"`
	Scenario    string    `json:"scenario"`
	Timestamp   time.Time `json:"timestamp"`
}