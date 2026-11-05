package tts

import (
	"fmt"
	"sync"

	"howl-chat/internal/audio/types"
)

// BaseSynthesizer provides shared behavior for TTS implementations.
type BaseSynthesizer struct {
	modelType   types.TTSModelType
	initialized bool

	// callMu serializes operations on the underlying synthesizer resource.
	// It is acquired before stateMu; never the other way around.
	callMu  sync.Mutex
	stateMu sync.RWMutex
	config  types.TTSConfig
}

// NewBaseSynthesizer creates a base synthesizer wrapper.
func NewBaseSynthesizer(modelType types.TTSModelType) *BaseSynthesizer {
	return &BaseSynthesizer{
		modelType: modelType,
	}
}

// GetModelType returns the configured model type.
func (b *BaseSynthesizer) GetModelType() types.TTSModelType {
	return b.modelType
}

// LockCalls acquires the provider operation lock.
func (b *BaseSynthesizer) LockCalls() { b.callMu.Lock() }

// UnlockCalls releases the provider operation lock.
func (b *BaseSynthesizer) UnlockCalls() { b.callMu.Unlock() }

// IsInitialized checks if the synthesizer has been initialized.
func (b *BaseSynthesizer) IsInitialized() bool {
	b.stateMu.RLock()
	defer b.stateMu.RUnlock()
	return b.initialized
}

// SetInitialized marks the synthesizer as initialized.
func (b *BaseSynthesizer) SetInitialized(config types.TTSConfig) {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()
	b.initialized = true
	b.config = config
}

// ClearInitialized returns the synthesizer to its uninitialized state.
func (b *BaseSynthesizer) ClearInitialized() {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()
	b.initialized = false
	b.config = types.TTSConfig{}
}

// GetConfig returns the active synthesizer configuration.
func (b *BaseSynthesizer) GetConfig() types.TTSConfig {
	b.stateMu.RLock()
	defer b.stateMu.RUnlock()
	return b.config
}

// CheckInitialized verifies the synthesizer is ready.
func (b *BaseSynthesizer) CheckInitialized() error {
	b.stateMu.RLock()
	defer b.stateMu.RUnlock()
	if !b.initialized {
		return types.NewAudioError(
			types.ErrCodeModelNotFound,
			fmt.Sprintf("tts model not initialized: %s", b.modelType),
			nil,
		)
	}
	return nil
}

// ValidateRequest performs baseline request checks shared by implementations.
func (b *BaseSynthesizer) ValidateRequest(req types.SynthesisRequest) error {
	if req.Text == "" {
		return types.NewAudioError(types.ErrCodeSynthesisFailed, "synthesis text is empty", nil)
	}
	if req.Speed > 0 && (req.Speed < 0.5 || req.Speed > 2.0) {
		return types.NewAudioError(types.ErrCodeSynthesisFailed, "synthesis speed must be between 0.5 and 2.0", nil)
	}
	return nil
}

// SynthesizerFactory creates TTS synthesizers based on registered model type.
type SynthesizerFactory struct {
	registry map[types.TTSModelType]func() types.TTSSynthesizer
}

// NewSynthesizerFactory creates an empty TTS synthesizer factory.
func NewSynthesizerFactory() *SynthesizerFactory {
	return &SynthesizerFactory{
		registry: make(map[types.TTSModelType]func() types.TTSSynthesizer),
	}
}

// Register binds a constructor to a model type.
func (f *SynthesizerFactory) Register(modelType types.TTSModelType, constructor func() types.TTSSynthesizer) {
	f.registry[modelType] = constructor
}

// Create builds a TTS synthesizer from a registered model type.
func (f *SynthesizerFactory) Create(modelType types.TTSModelType) (types.TTSSynthesizer, error) {
	constructor, ok := f.registry[modelType]
	if !ok {
		return nil, fmt.Errorf("no synthesizer registered for model type: %s", modelType)
	}
	return constructor(), nil
}

// ListAvailable returns all registered synthesizer model types.
func (f *SynthesizerFactory) ListAvailable() []types.TTSModelType {
	models := make([]types.TTSModelType, 0, len(f.registry))
	for t := range f.registry {
		models = append(models, t)
	}
	return models
}

// RecommendForHardware suggests a baseline TTS model for the target tier.
func RecommendForHardware(tier types.MemoryTier) types.TTSModelType {
	return types.RecommendTTS(tier)
}
