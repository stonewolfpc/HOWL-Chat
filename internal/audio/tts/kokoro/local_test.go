package kokoro

import (
	"context"
	"os"
	"testing"

	"howl-chat/internal/audio/types"
)

func TestKokoroListVoices(t *testing.T) {
	s := New()
	cfg := types.TTSConfig{
		ModelPath:  `D:\AI GGUF\TTS\kokoro-v1.0.onnx`,
		VoicesPath: `D:\AI GGUF\TTS\voices-v1.0.bin`,
		ConfigPath: `D:\AI GGUF\TTS\hf_cache\config.json`,
		CacheDir:   `D:\AI GGUF\TTS\cache`,
		ScriptPath: `D:\HOWL_Chat\scripts\kokoro-tts.py`,
	}
	if err := s.Initialize(context.Background(), cfg); err != nil {
		t.Skipf("Kokoro initialization skipped (model not available): %v", err)
		return
	}
	voices := s.GetVoices()
	if len(voices) == 0 {
		t.Error("Expected voices from kokoro archive, got none")
	}
	for _, v := range voices[:5] {
		t.Logf("Voice: %s", v.ID)
	}
}

func TestKokoroSynthesize(t *testing.T) {
	s := New()
	cfg := types.TTSConfig{
		ModelPath:  `D:\AI GGUF\TTS\kokoro-v1.0.onnx`,
		VoicesPath: `D:\AI GGUF\TTS\voices-v1.0.bin`,
		ConfigPath: `D:\AI GGUF\TTS\hf_cache\config.json`,
		CacheDir:   `D:\AI GGUF\TTS\cache`,
		ScriptPath: `D:\HOWL_Chat\scripts\kokoro-tts.py`,
	}
	if err := s.Initialize(context.Background(), cfg); err != nil {
		t.Skipf("Kokoro initialization skipped: %v", err)
		return
	}
	result, err := s.Synthesize(context.Background(), types.SynthesisRequest{
		Text:  "Hello, this is Kokoro TTS.",
		Voice: "af_heart",
		Speed: 1.0,
	}, nil)
	if err != nil {
		t.Fatalf("Synthesis failed: %v", err)
	}
	if result.AudioPath == "" {
		t.Error("Expected audio path, got empty string")
	}
	if _, err := os.Stat(result.AudioPath); err != nil {
		t.Errorf("Audio file not found: %v", err)
	}
	t.Logf("Synthesized: %s (%.1fs)", result.AudioPath, result.Duration)
}
