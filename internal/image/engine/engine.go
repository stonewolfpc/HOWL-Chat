package engine

import (
	"howl-chat/internal/image/types"
)

// ImageEngine is the universal image generation interface.
type ImageEngine interface {
	// Initialize validates paths and prepares the engine.
	Initialize(config types.EngineConfig) error

	// ListModels scans the models directory and returns discovered bundles.
	ListModels() ([]types.ModelInfo, error)

	// ListLoRAs scans the models directory for LoRA files.
	ListLoRAs() ([]string, error)

	// LoadModel selects a model for subsequent generations.
	LoadModel(opts types.ModelLoadOptions) error

	// IsRunning returns true if the engine backend is alive.
	IsRunning() bool

	// EnsureReady starts the engine and loads the model if needed.
	EnsureReady(modelName string, opts types.ModelLoadOptions) error

	// GenerateTxt2Img creates an image from a text prompt.
	GenerateTxt2Img(opts types.ImageOptions) (*types.ImageResult, error)

	// GenerateImg2Img transforms an existing image with a prompt.
	GenerateImg2Img(opts types.ImageOptions) (*types.ImageResult, error)

	// GenerateInpaint fills masked regions of an image.
	GenerateInpaint(opts types.ImageOptions) (*types.ImageResult, error)

	// GenerateVideo creates a series of frames (SVD/WAN).
	GenerateVideo(opts types.ImageOptions) (*types.ImageResult, error)

	// GetProgress returns the current generation progress from the backend.
	GetProgress() (types.ProgressInfo, error)

	// GetLoadedModel returns the currently loaded model, or nil.
	GetLoadedModel() *types.ModelInfo

	// Shutdown releases all resources.
	Shutdown() error
}
