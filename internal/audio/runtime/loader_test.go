package runtime

import (
	"context"
	"errors"
	"strings"
	"testing"

	"howl-chat/internal/audio/types"
)

func TestModelLoaderLifecycle(t *testing.T) {
	loader := NewModelLoader()

	loader.RegisterASR(types.ASRWhisperTiny, func() types.ASRRecognizer { return &mockASR{} })
	loader.RegisterTTS(types.TTSPiper, func() types.TTSSynthesizer { return &mockTTS{} })
	loader.RegisterMultimodal(types.MultimodalQwen25Omni7B, func() types.MultimodalBridge { return &mockMultimodal{} })

	if _, err := loader.LoadASR(context.Background(), string(types.ASRWhisperTiny), types.ASRConfig{}); err != nil {
		t.Fatalf("LoadASR failed: %v", err)
	}
	if _, err := loader.LoadTTS(context.Background(), string(types.TTSPiper), types.TTSConfig{}); err != nil {
		t.Fatalf("LoadTTS failed: %v", err)
	}
	if _, err := loader.LoadMultimodal(string(types.MultimodalQwen25Omni7B), types.MultimodalConfig{}); err != nil {
		t.Fatalf("LoadMultimodal failed: %v", err)
	}

	loaded := loader.GetLoadedModels()
	if len(loaded) != 3 {
		t.Fatalf("expected 3 loaded models, got %d", len(loaded))
	}

	ttsID := ""
	for _, id := range loaded {
		if strings.HasPrefix(id, "tts:"+string(types.TTSPiper)+":") {
			ttsID = id
		}
	}
	if ttsID == "" {
		t.Fatalf("no unique TTS id in %v", loaded)
	}

	if err := loader.Unload("asr:" + string(types.ASRWhisperTiny)); err != nil {
		t.Fatalf("Unload ASR failed: %v", err)
	}
	if err := loader.Unload(ttsID); err != nil {
		t.Fatalf("Unload TTS failed: %v", err)
	}
	if err := loader.Unload("multimodal:" + string(types.MultimodalQwen25Omni7B)); err != nil {
		t.Fatalf("Unload multimodal failed: %v", err)
	}
}

func TestModelLoaderContextForwardingAndCancellation(t *testing.T) {
	loader := NewModelLoader()
	type contextKey struct{}
	captured := make(chan context.Context, 1)
	created := 0
	loader.RegisterASR(types.ASRWhisperTiny, func() types.ASRRecognizer {
		created++
		return &contextASR{captured: captured}
	})

	ctx := context.WithValue(context.Background(), contextKey{}, "forwarded")
	if _, err := loader.LoadASR(ctx, string(types.ASRWhisperTiny), types.ASRConfig{}); err != nil {
		t.Fatalf("LoadASR: %v", err)
	}
	got := <-captured
	if got.Value(contextKey{}) != "forwarded" {
		t.Fatal("loader did not forward the caller context")
	}
	if err := loader.Shutdown(); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := loader.LoadASR(cancelled, string(types.ASRWhisperTiny), types.ASRConfig{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled LoadASR = %v, want context.Canceled", err)
	}
	if created != 1 {
		t.Fatalf("cancelled context created %d providers, want 1 total", created)
	}
}

func TestModelLoaderUniqueTTSIDs(t *testing.T) {
	loader := NewModelLoader()
	loader.RegisterTTS(types.TTSPiper, func() types.TTSSynthesizer { return &mockTTS{} })
	for i := 0; i < 2; i++ {
		if _, err := loader.LoadTTS(context.Background(), string(types.TTSPiper), types.TTSConfig{}); err != nil {
			t.Fatalf("LoadTTS %d: %v", i, err)
		}
	}
	seen := map[string]bool{}
	for _, id := range loader.GetLoadedModels() {
		if strings.HasPrefix(id, "tts:"+string(types.TTSPiper)+":") {
			if seen[id] {
				t.Fatalf("duplicate TTS id %q", id)
			}
			seen[id] = true
		}
	}
	if len(seen) != 2 {
		t.Fatalf("unique TTS ids = %d, want 2", len(seen))
	}
	if err := loader.Shutdown(); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestModelLoaderP6FailedProviderRelease(t *testing.T) {
	t.Run("ASR", func(t *testing.T) {
		loader := NewModelLoader()
		initErr := errors.New("asr init failed")
		releaseErr := errors.New("asr release failed")
		provider := &failingASR{initErr: initErr, releaseErr: releaseErr}
		loader.RegisterASR(types.ASRWhisperTiny, func() types.ASRRecognizer { return provider })
		_, err := loader.LoadASR(context.Background(), string(types.ASRWhisperTiny), types.ASRConfig{})
		if !errors.Is(err, initErr) || !errors.Is(err, releaseErr) {
			t.Fatalf("LoadASR error = %v, want init and release errors", err)
		}
		if provider.releaseCalls != 1 {
			t.Fatalf("release calls = %d, want 1", provider.releaseCalls)
		}
		if len(loader.GetLoadedModels()) != 0 {
			t.Fatalf("failed provider was tracked: %v", loader.GetLoadedModels())
		}
	})

	t.Run("TTS", func(t *testing.T) {
		loader := NewModelLoader()
		initErr := errors.New("tts init failed")
		releaseErr := errors.New("tts release failed")
		provider := &failingTTS{initErr: initErr, releaseErr: releaseErr}
		loader.RegisterTTS(types.TTSPiper, func() types.TTSSynthesizer { return provider })
		_, err := loader.LoadTTS(context.Background(), string(types.TTSPiper), types.TTSConfig{})
		if !errors.Is(err, initErr) || !errors.Is(err, releaseErr) {
			t.Fatalf("LoadTTS error = %v, want init and release errors", err)
		}
		if provider.releaseCalls != 1 {
			t.Fatalf("release calls = %d, want 1", provider.releaseCalls)
		}
		if len(loader.GetLoadedModels()) != 0 {
			t.Fatalf("failed provider was tracked: %v", loader.GetLoadedModels())
		}
	})
}

func TestModelLoaderP6PanickingProviderRelease(t *testing.T) {
	t.Run("cleanup error", func(t *testing.T) {
		loader := NewModelLoader()
		initPanic := errors.New("tts init panic")
		releaseErr := errors.New("tts release error")
		provider := &panicTTS{initPanic: initPanic, releaseErr: releaseErr}
		loader.RegisterTTS(types.TTSPiper, func() types.TTSSynthesizer { return provider })

		var recovered any
		func() {
			defer func() { recovered = recover() }()
			_, _ = loader.LoadTTS(context.Background(), string(types.TTSPiper), types.TTSConfig{})
		}()
		panicRecord, ok := recovered.(*ProviderInitPanic)
		if !ok {
			t.Fatalf("recovered = %T %v, want *ProviderInitPanic", recovered, recovered)
		}
		if panicRecord.Value != initPanic {
			t.Fatalf("Value = %v, want original init panic", panicRecord.Value)
		}
		if !errors.Is(panicRecord, initPanic) || !errors.Is(panicRecord, releaseErr) {
			t.Fatalf("wrapper does not preserve both errors: %v", panicRecord)
		}
		if panicRecord.CleanupErr != releaseErr {
			t.Fatalf("CleanupErr = %v, want %v", panicRecord.CleanupErr, releaseErr)
		}
		if provider.releaseCalls != 1 {
			t.Fatalf("release calls = %d, want 1", provider.releaseCalls)
		}
		if len(loader.GetLoadedModels()) != 0 {
			t.Fatalf("panicking provider was tracked: %v", loader.GetLoadedModels())
		}
	})

	t.Run("cleanup panic", func(t *testing.T) {
		loader := NewModelLoader()
		initPanic := errors.New("tts init panic")
		releasePanic := errors.New("tts release panic")
		provider := &panicTTS{initPanic: initPanic, releasePanic: releasePanic}
		loader.RegisterTTS(types.TTSPiper, func() types.TTSSynthesizer { return provider })

		var recovered any
		func() {
			defer func() { recovered = recover() }()
			_, _ = loader.LoadTTS(context.Background(), string(types.TTSPiper), types.TTSConfig{})
		}()
		panicRecord, ok := recovered.(*ProviderInitPanic)
		if !ok {
			t.Fatalf("recovered = %T %v, want *ProviderInitPanic", recovered, recovered)
		}
		if panicRecord.Value != initPanic {
			t.Fatalf("Value = %v, want original init panic", panicRecord.Value)
		}
		if panicRecord.CleanupPanic != releasePanic {
			t.Fatalf("CleanupPanic = %v, want release panic", panicRecord.CleanupPanic)
		}
		if provider.releaseCalls != 1 {
			t.Fatalf("release calls = %d, want 1", provider.releaseCalls)
		}
	})
}

// Shutdown is a test helper that unloads every tracked provider.
func (l *ModelLoader) Shutdown() error {
	var errs []error
	for _, id := range l.GetLoadedModels() {
		if err := l.Unload(id); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

type contextASR struct {
	mockASR
	captured chan context.Context
}

func (c *contextASR) Initialize(ctx context.Context, config types.ASRConfig) error {
	c.captured <- ctx
	return nil
}

type failingASR struct {
	mockASR
	initErr      error
	releaseErr   error
	releaseCalls int
}

func (f *failingASR) Initialize(context.Context, types.ASRConfig) error { return f.initErr }
func (f *failingASR) Release() error {
	f.releaseCalls++
	return f.releaseErr
}

type failingTTS struct {
	mockTTS
	initErr      error
	releaseErr   error
	releaseCalls int
}

func (f *failingTTS) Initialize(context.Context, types.TTSConfig) error { return f.initErr }
func (f *failingTTS) Release() error {
	f.releaseCalls++
	return f.releaseErr
}

type panicTTS struct {
	mockTTS
	initPanic    any
	releaseErr   error
	releasePanic any
	releaseCalls int
}

func (p *panicTTS) Initialize(context.Context, types.TTSConfig) error { panic(p.initPanic) }
func (p *panicTTS) Release() error {
	p.releaseCalls++
	if p.releasePanic != nil {
		panic(p.releasePanic)
	}
	return p.releaseErr
}

type mockASR struct{}

func (m *mockASR) Initialize(ctx context.Context, config types.ASRConfig) error { return nil }
func (m *mockASR) Transcribe(ctx context.Context, audioPath string, callback types.ProgressCallback) (*types.RecognitionResult, error) {
	return &types.RecognitionResult{Text: "ok"}, nil
}
func (m *mockASR) TranscribeStream(ctx context.Context, chunks <-chan types.AudioChunk, results chan<- types.RecognitionResult) error {
	return nil
}
func (m *mockASR) SupportsLanguage(lang string) bool { return true }
func (m *mockASR) GetModelInfo() types.ASRModelInfo  { return types.ASRModelInfo{Name: "mock"} }
func (m *mockASR) Release() error                    { return nil }

type mockTTS struct{}

func (m *mockTTS) Initialize(ctx context.Context, config types.TTSConfig) error { return nil }
func (m *mockTTS) Synthesize(ctx context.Context, req types.SynthesisRequest, callback types.ProgressCallback) (*types.SynthesisResult, error) {
	return &types.SynthesisResult{AudioPath: "mock.wav"}, nil
}
func (m *mockTTS) SynthesizeStream(ctx context.Context, req types.SynthesisRequest, chunks chan<- types.AudioChunk) error {
	return nil
}
func (m *mockTTS) GetVoices() []types.Voice { return []types.Voice{{ID: "mock"}} }
func (m *mockTTS) SupportsSSML() bool       { return false }
func (m *mockTTS) Release() error           { return nil }

type mockMultimodal struct{}

func (m *mockMultimodal) Initialize(ctx context.Context, config types.MultimodalConfig) error {
	return nil
}
func (m *mockMultimodal) ProcessAudioText(ctx context.Context, audio types.AudioChunk, textContext string) (string, error) {
	return "ok", nil
}
func (m *mockMultimodal) ProcessAudioVision(ctx context.Context, audio types.AudioChunk, imageData []byte) (string, error) {
	return "ok", nil
}
func (m *mockMultimodal) GetTokenBudget(audioDuration float64) int { return 1024 }
func (m *mockMultimodal) Release() error                           { return nil }
