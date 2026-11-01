package kokoro

import (
	"context"
	"os"
	"testing"

	"howl-chat/internal/audio/types"
)

// kokoroTestConfig builds a TTS config from the environment so the integration
// tests run wherever the model set actually lives. Voices and config are
// derived from the model directory by Initialize, and the cache goes to a temp
// dir so a test run never writes into a real cache.
func kokoroTestConfig(t *testing.T) types.TTSConfig {
	t.Helper()
	model := os.Getenv("HOWL_TTS_KOKORO_MODEL")
	if model == "" {
		t.Skip("set HOWL_TTS_KOKORO_MODEL (optionally HOWL_TTS_KOKORO_VOICES / HOWL_TTS_KOKORO_CONFIG) to run Kokoro integration tests")
	}
	return types.TTSConfig{
		ModelPath:  model,
		VoicesPath: os.Getenv("HOWL_TTS_KOKORO_VOICES"),
		ConfigPath: os.Getenv("HOWL_TTS_KOKORO_CONFIG"),
		CacheDir:   t.TempDir(),
		ScriptPath: defaultKokoroScriptPath(),
	}
}

func TestKokoroListVoices(t *testing.T) {
	s := New()
	cfg := kokoroTestConfig(t)
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
	cfg := kokoroTestConfig(t)
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
