package audio

import (
	"context"
	"os"
	"strings"
	"testing"

	"howl-chat/internal/audio/capability"
	"howl-chat/internal/audio/types"
)

// TestQwenASRPipelineE2E exercises the application pipeline against the real
// bundled llama-server and the installed Qwen3-ASR pair. It is opt-in because
// it loads a 2.1 GB model and needs the real assets:
//
//	HOWL_QWEN_ASR_E2E=1 HOWL_QWEN_ASR_WAV=<16k wav> \
//	HOWL_ASR_MODEL=<model.gguf> HOWL_ASR_MMPROJ=<mmproj.gguf> go test ./internal/audio -run TestQwenASRPipelineE2E -v
func TestQwenASRPipelineE2E(t *testing.T) {
	if os.Getenv("HOWL_QWEN_ASR_E2E") != "1" {
		t.Skip("set HOWL_QWEN_ASR_E2E=1 and HOWL_QWEN_ASR_WAV to run the real pipeline transcription")
	}
	wav := strings.TrimSpace(os.Getenv("HOWL_QWEN_ASR_WAV"))
	if wav == "" {
		t.Skip("HOWL_QWEN_ASR_WAV is not set")
	}

	previous := pipelineCapabilityCatalog
	pipelineCapabilityCatalog = func() []capability.Record {
		return capability.Catalog(capability.ProductionEnv(ASRModelPathFor, ASRProjectorPathFor, TTSModelPathFor))
	}
	t.Cleanup(func() { pipelineCapabilityCatalog = previous })

	pipeline := NewPipeline()
	if err := pipeline.Initialize(context.Background(), types.PipelineConfig{
		ASREnabled:   true,
		ASREngine:    string(types.ASRQwen3ASR17BLocal),
		HardwareTier: "large",
	}); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	stopped := false
	t.Cleanup(func() {
		if !stopped {
			_ = pipeline.Shutdown(context.Background())
		}
	})

	response, err := pipeline.ProcessAudioMessage(context.Background(), types.AudioMessage{
		Type:      "audio_input",
		AudioPath: wav,
	})
	if err != nil {
		t.Fatalf("ProcessAudioMessage: %v", err)
	}
	text := response.TranscribedText
	lower := strings.ToLower(text)
	if !strings.Contains(lower, "quick brown fox") {
		t.Fatalf("transcript = %q, want the known sentence", text)
	}
	if strings.Contains(lower, "<asr_text>") || strings.Contains(lower, "language english") {
		t.Fatalf("transcript still contains the server wrapper: %q", text)
	}

	stopped = true
	if err := pipeline.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	if loaded := pipeline.loader.GetLoadedModels(); len(loaded) != 0 {
		t.Fatalf("loaded models after shutdown = %v, want none", loaded)
	}
}
