package runtime

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"howl-chat/internal/audio/asr"
	"howl-chat/internal/audio/types"
	ttsruntime "howl-chat/internal/audio/tts"
)

type multimodalConstructor func() types.MultimodalBridge

// ModelLoader is a modular runtime loader with explicit registration points.
type ModelLoader struct {
	asrFactory *asr.RecognizerFactory
	ttsFactory *ttsruntime.SynthesizerFactory

	multimodalRegistry map[types.MultimodalModelType]multimodalConstructor

	mu          sync.RWMutex
	loaded      map[string]func() error
	loadedKinds map[string]string
}

// NewModelLoader creates an empty, production-safe loader registry.
func NewModelLoader() *ModelLoader {
	return &ModelLoader{
		asrFactory:         asr.NewRecognizerFactory(),
		ttsFactory:         ttsruntime.NewSynthesizerFactory(),
		multimodalRegistry: make(map[types.MultimodalModelType]multimodalConstructor),
		loaded:             make(map[string]func() error),
		loadedKinds:        make(map[string]string),
	}
}

// RegisterASR registers an ASR constructor in the loader.
func (l *ModelLoader) RegisterASR(modelType types.ASRModelType, constructor func() types.ASRRecognizer) {
	l.asrFactory.Register(modelType, constructor)
}

// RegisterTTS registers a TTS constructor in the loader.
func (l *ModelLoader) RegisterTTS(modelType types.TTSModelType, constructor func() types.TTSSynthesizer) {
	l.ttsFactory.Register(modelType, constructor)
}

// RegisterMultimodal registers a multimodal bridge constructor in the loader.
func (l *ModelLoader) RegisterMultimodal(modelType types.MultimodalModelType, constructor func() types.MultimodalBridge) {
	l.multimodalRegistry[modelType] = constructor
}

// LoadASR initializes an ASR model and tracks it for lifecycle management.
func (l *ModelLoader) LoadASR(modelType string, config types.ASRConfig) (types.ASRRecognizer, error) {
	recognizer, err := l.asrFactory.Create(types.ASRModelType(modelType))
	if err != nil {
		return nil, err
	}
	if err := recognizer.Initialize(context.Background(), config); err != nil {
		return nil, err
	}
	id := "asr:" + modelType
	l.trackLoaded(id, "asr", recognizer.Release)
	return recognizer, nil
}

// LoadTTS initializes a TTS model and tracks it for lifecycle management.
func (l *ModelLoader) LoadTTS(modelType string, config types.TTSConfig) (types.TTSSynthesizer, error) {
	synth, err := l.ttsFactory.Create(types.TTSModelType(modelType))
	if err != nil {
		return nil, err
	}
	if err := synth.Initialize(context.Background(), config); err != nil {
		return nil, err
	}
	id := "tts:" + modelType
	l.trackLoaded(id, "tts", synth.Release)
	return synth, nil
}

// LoadMultimodal initializes a multimodal bridge and tracks it for lifecycle management.
func (l *ModelLoader) LoadMultimodal(modelType string, config types.MultimodalConfig) (types.MultimodalBridge, error) {
	constructor, ok := l.multimodalRegistry[types.MultimodalModelType(modelType)]
	if !ok {
		return nil, fmt.Errorf("no multimodal bridge registered for model type: %s", modelType)
	}
	bridge := constructor()
	if err := bridge.Initialize(context.Background(), config); err != nil {
		return nil, err
	}
	id := "multimodal:" + modelType
	l.trackLoaded(id, "multimodal", bridge.Release)
	return bridge, nil
}

// Unload releases a loaded model by its loader-managed id.
func (l *ModelLoader) Unload(modelID string) error {
	l.mu.Lock()
	release, ok := l.loaded[modelID]
	if ok {
		delete(l.loaded, modelID)
		delete(l.loadedKinds, modelID)
	}
	l.mu.Unlock()
	if !ok {
		return fmt.Errorf("model not loaded: %s", modelID)
	}
	return release()
}

// GetLoadedModels returns stable-sorted IDs for loaded runtime models.
func (l *ModelLoader) GetLoadedModels() []string {
	l.mu.RLock()
	models := make([]string, 0, len(l.loaded))
	for id := range l.loaded {
		models = append(models, id)
	}
	l.mu.RUnlock()
	sort.Strings(models)
	return models
}

// GetLoadedKinds returns loaded model kind metadata keyed by model id.
func (l *ModelLoader) GetLoadedKinds() map[string]string {
	l.mu.RLock()
	out := make(map[string]string, len(l.loadedKinds))
	for k, v := range l.loadedKinds {
		out[k] = v
	}
	l.mu.RUnlock()
	return out
}

func (l *ModelLoader) trackLoaded(id, kind string, releaseFn func() error) {
	l.mu.Lock()
	l.loaded[id] = releaseFn
	l.loadedKinds[id] = kind
	l.mu.Unlock()
}
