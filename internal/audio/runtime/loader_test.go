package runtime

import (
	"context"
	"testing"

	"howl-chat/internal/audio/types"
)

func TestModelLoaderLifecycle(t *testing.T) {
	loader := NewModelLoader()

	loader.RegisterASR(types.ASRWhisperTiny, func() types.ASRRecognizer { return &mockASR{} })
	loader.RegisterTTS(types.TTSPiper, func() types.TTSSynthesizer { return &mockTTS{} })
	loader.RegisterMultimodal(types.MultimodalQwen25Omni7B, func() types.MultimodalBridge { return &mockMultimodal{} })

	if _, err := loader.LoadASR(string(types.ASRWhisperTiny), types.ASRConfig{}); err != nil {
		t.Fatalf("LoadASR failed: %v", err)
	}
	if _, err := loader.LoadTTS(string(types.TTSPiper), types.TTSConfig{}); err != nil {
		t.Fatalf("LoadTTS failed: %v", err)
	}
	if _, err := loader.LoadMultimodal(string(types.MultimodalQwen25Omni7B), types.MultimodalConfig{}); err != nil {
		t.Fatalf("LoadMultimodal failed: %v", err)
	}

	loaded := loader.GetLoadedModels()
	if len(loaded) != 3 {
		t.Fatalf("expected 3 loaded models, got %d", len(loaded))
	}

	if err := loader.Unload("asr:" + string(types.ASRWhisperTiny)); err != nil {
		t.Fatalf("Unload ASR failed: %v", err)
	}
	if err := loader.Unload("tts:" + string(types.TTSPiper)); err != nil {
		t.Fatalf("Unload TTS failed: %v", err)
	}
	if err := loader.Unload("multimodal:" + string(types.MultimodalQwen25Omni7B)); err != nil {
		t.Fatalf("Unload multimodal failed: %v", err)
	}
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
func (m *mockASR) GetModelInfo() types.ASRModelInfo { return types.ASRModelInfo{Name: "mock"} }
func (m *mockASR) Release() error { return nil }

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

func (m *mockMultimodal) Initialize(ctx context.Context, config types.MultimodalConfig) error { return nil }
func (m *mockMultimodal) ProcessAudioText(ctx context.Context, audio types.AudioChunk, textContext string) (string, error) {
	return "ok", nil
}
func (m *mockMultimodal) ProcessAudioVision(ctx context.Context, audio types.AudioChunk, imageData []byte) (string, error) {
	return "ok", nil
}
func (m *mockMultimodal) GetTokenBudget(audioDuration float64) int { return 1024 }
func (m *mockMultimodal) Release() error                            { return nil }
