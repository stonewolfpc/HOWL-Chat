/**
 * Model Manager
 *
 * This package provides model lifecycle management including loading, unloading,
 * caching, and metadata handling. Ensures thread-safe model operations.
 *
 * @package model
 */

package model

import (
	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/types"
	"sync"
)

// Manager handles model lifecycle management
type Manager struct {
	mu              sync.RWMutex
	currentModel    *types.Model
	modelCache      map[string]*types.Model
	llamaClient     llama.Client
	loadOptions     *llama.LoadOptions
	progressTracker *llama.ProgressTracker
}

// NewManager creates a new model manager
func NewManager(client llama.Client) *Manager {
	return &Manager{
		currentModel:    nil,
		modelCache:      make(map[string]*types.Model),
		llamaClient:     client,
		loadOptions:     nil, // Will be set per-load
		progressTracker: llama.NewProgressTracker(),
	}
}

// LoadModel loads a model from the given path
func (m *Manager) LoadModel(modelPath string, options *llama.LoadOptions) (*types.Model, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Unload current model if loaded
	if m.currentModel != nil {
		if err := m.UnloadModel(); err != nil {
			return nil, types.WrapError(types.ErrorCodeModelUnload,
				"failed to unload current model", err)
		}
	}

	// Set up progress tracking
	m.progressTracker.Reset()
	m.progressTracker.SetCurrentStage(llama.StageModelLoadStart)

	// Configure load options
	if options == nil {
		options = llama.NewLoadOptions(modelPath)
	}

	// Validate options
	if err := options.Validate(); err != nil {
		return nil, types.WrapError(types.ErrorCodeModelLoad,
			"invalid load options", err)
	}

	// Update progress
	m.progressTracker.SetCurrentStage(llama.StageTokenizerLoad)
	m.progressTracker.UpdateStageProgress(1.0)

	// Load model using llama client
	m.progressTracker.SetCurrentStage(llama.StageTensorAllocation)
	m.progressTracker.UpdateStageProgress(0.5)

	if err := m.llamaClient.LoadModel(modelPath, options); err != nil {
		m.progressTracker.SetCurrentStage(llama.StageModelReady)
		m.progressTracker.UpdateStageProgress(0.0)
		return nil, types.WrapError(types.ErrorCodeModelLoad,
			"failed to load model", err)
	}

	m.progressTracker.UpdateStageProgress(1.0)
	m.progressTracker.SetCurrentStage(llama.StageKVCacheInit)
	m.progressTracker.UpdateStageProgress(0.5)

	// Create model instance
	model := types.NewModel("", modelPath, types.FormatGGUF, 0)
	model.Status = types.ModelStatusLoaded

	m.currentModel = model
	m.modelCache[modelPath] = model

	m.progressTracker.SetCurrentStage(llama.StageModelReady)
	m.progressTracker.UpdateStageProgress(1.0)

	return model, nil
}

// UnloadModel unloads the currently loaded model
func (m *Manager) UnloadModel() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentModel == nil {
		return types.ErrModelNotFound("no model currently loaded")
	}

	if err := m.llamaClient.UnloadModel(); err != nil {
		m.currentModel.Status = types.ModelStatusError
		return types.WrapError(types.ErrorCodeModelUnload,
			"failed to unload model", err)
	}

	m.currentModel.Status = types.ModelStatusUnloaded
	m.currentModel = nil

	return nil
}

// GetCurrentModel returns the currently loaded model
func (m *Manager) GetCurrentModel() (*types.Model, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.currentModel == nil {
		return nil, types.ErrModelNotFound("no model currently loaded")
	}

	return m.currentModel, nil
}

// IsModelLoaded returns true if a model is currently loaded
func (m *Manager) IsModelLoaded() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.currentModel != nil && m.currentModel.IsLoaded()
}

// GetLoadedModelName returns the name of the currently loaded model
func (m *Manager) GetLoadedModelName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.currentModel == nil {
		return "No model loaded"
	}
	return m.currentModel.Name
}

// GetLoadingProgress returns the current loading progress
func (m *Manager) GetLoadingProgress() float64 {
	return m.progressTracker.GetTotalProgress()
}

// GetLoadingStage returns the current loading stage
func (m *Manager) GetLoadingStage() llama.LoadingStage {
	return m.progressTracker.GetCurrentStage()
}

// SetProgressCallback sets a callback for loading progress
func (m *Manager) SetProgressCallback(callback llama.ProgressCallback) {
	m.llamaClient.SetProgressCallback(callback)
}

// GetModelInfo returns information about the currently loaded model
func (m *Manager) GetModelInfo() (*types.Model, error) {
	return m.GetCurrentModel()
}

// ClearCache clears the model cache
func (m *Manager) ClearCache() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.modelCache = make(map[string]*types.Model)
}

// GetCachedModels returns all cached models
func (m *Manager) GetCachedModels() []*types.Model {
	m.mu.RLock()
	defer m.mu.RUnlock()

	models := make([]*types.Model, 0, len(m.modelCache))
	for _, model := range m.modelCache {
		models = append(models, model)
	}

	return models
}

// Close releases all resources
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentModel != nil {
		if err := m.UnloadModel(); err != nil {
			return err
		}
	}

	if err := m.llamaClient.Close(); err != nil {
		return types.WrapError(types.ErrorCodeInternal,
			"failed to close llama client", err)
	}

	return nil
}
