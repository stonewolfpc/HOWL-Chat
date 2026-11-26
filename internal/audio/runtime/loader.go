package runtime

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"howl-chat/internal/audio/asr"
	ttsruntime "howl-chat/internal/audio/tts"
	"howl-chat/internal/audio/types"
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
	nextTTSID   uint64
}

// ProviderInitPanic records a provider initialization panic after the loader
// has released the provider. CleanupErr and CleanupPanic preserve failures that
// occurred while releasing the panicking provider.
type ProviderInitPanic struct {
	ProviderID   string
	Value        any
	CleanupErr   error
	CleanupPanic any
}

func (e *ProviderInitPanic) Error() string {
	if e == nil {
		return "<nil>"
	}
	message := fmt.Sprintf("provider %s initialization panicked: %v", e.ProviderID, e.Value)
	if e.CleanupErr != nil {
		message += fmt.Sprintf("; cleanup failed: %v", e.CleanupErr)
	}
	if e.CleanupPanic != nil {
		message += fmt.Sprintf("; cleanup panicked: %v", e.CleanupPanic)
	}
	return message
}

// Unwrap exposes the original panic when it is an error and the cleanup error
// when one was observed. The cleanup panic remains available directly.
func (e *ProviderInitPanic) Unwrap() []error {
	if e == nil {
		return nil
	}
	errs := make([]error, 0, 2)
	if original, ok := e.Value.(error); ok {
		errs = append(errs, original)
	}
	if e.CleanupErr != nil {
		errs = append(errs, e.CleanupErr)
	}
	return errs
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
func (l *ModelLoader) LoadASR(ctx context.Context, modelType string, config types.ASRConfig) (types.ASRRecognizer, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	recognizer, err := l.asrFactory.Create(types.ASRModelType(modelType))
	if err != nil {
		return nil, err
	}
	id := "asr:" + modelType
	if err := initializeProvider(id, func() error {
		return recognizer.Initialize(ctx, config)
	}, recognizer.Release); err != nil {
		releaseErr := recognizer.Release()
		return nil, errors.Join(err, releaseErr)
	}
	l.trackLoaded(id, "asr", recognizer.Release)
	return recognizer, nil
}

// LoadTTS initializes a TTS model and tracks it for lifecycle management.
func (l *ModelLoader) LoadTTS(ctx context.Context, modelType string, config types.TTSConfig) (types.TTSSynthesizer, error) {
	synth, _, err := l.LoadTTSWithID(ctx, modelType, config)
	return synth, err
}

// LoadTTSWithID is LoadTTS that also returns the loader id, so the owner can
// later release exactly this provider instead of sweeping every loaded model.
func (l *ModelLoader) LoadTTSWithID(ctx context.Context, modelType string, config types.TTSConfig) (types.TTSSynthesizer, string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, "", err
	}
	synth, err := l.ttsFactory.Create(types.TTSModelType(modelType))
	if err != nil {
		return nil, "", err
	}
	if err := initializeProvider("tts:"+modelType, func() error {
		return synth.Initialize(ctx, config)
	}, synth.Release); err != nil {
		releaseErr := synth.Release()
		return nil, "", errors.Join(err, releaseErr)
	}
	id := l.trackTTSLoaded(modelType, synth.Release)
	return synth, id, nil
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

func (l *ModelLoader) trackTTSLoaded(modelType string, releaseFn func() error) string {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.nextTTSID++
	id := fmt.Sprintf("tts:%s:%d", modelType, l.nextTTSID)
	l.loaded[id] = releaseFn
	l.loadedKinds[id] = "tts"
	return id
}

func initializeProvider(providerID string, init func() error, release func() error) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			cleanupErr, cleanupPanic := safeRelease(release)
			panic(&ProviderInitPanic{
				ProviderID:   providerID,
				Value:        recovered,
				CleanupErr:   cleanupErr,
				CleanupPanic: cleanupPanic,
			})
		}
	}()
	return init()
}

func safeRelease(release func() error) (err error, recovered any) {
	defer func() {
		if r := recover(); r != nil {
			recovered = r
		}
	}()
	return release(), nil
}
