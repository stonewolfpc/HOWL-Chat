package audio

import (
	"context"
	"testing"

	"howl-chat/internal/audio/types"
)

func TestPipelineForwardsWhisperDecodeConfig(t *testing.T) {
	p, set := newScriptedPipeline(providerBehavior{})
	cfg := testPipelineConfig(true, false)
	cfg.ASREngine = string(types.ASRWhisperBase)
	cfg.ASRLanguage = "de"
	cfg.ASRThreads = 8
	cfg.ASRBeamSize = 3
	cfg.ASRBestOf = 2

	if err := p.Initialize(context.Background(), cfg); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	t.Cleanup(func() { _ = p.Shutdown(context.Background()) })

	provider := set.asr[types.ASRWhisperBase]
	if provider.initCalls.Load() != 1 {
		t.Fatalf("Whisper initialize calls = %d, want 1", provider.initCalls.Load())
	}
	if provider.config.Language != "de" || provider.config.NumThreads != 8 ||
		provider.config.BeamSize != 3 || provider.config.BestOf != 2 {
		t.Fatalf("Whisper config = language %q threads %d beam %d best-of %d, want de/8/3/2",
			provider.config.Language, provider.config.NumThreads, provider.config.BeamSize, provider.config.BestOf)
	}
}

func TestASRDecodeConfigIsOnlyForExplicitWhisper(t *testing.T) {
	cfg := types.PipelineConfig{
		ASREngine:   string(types.ASRWhisperBase),
		ASRLanguage: "de",
		ASRThreads:  8,
		ASRBeamSize: 3,
		ASRBestOf:   2,
	}
	got := asrDecodeConfigFor(cfg)
	if got.Language != "de" || got.NumThreads != 8 || got.BeamSize != 3 || got.BestOf != 2 {
		t.Fatalf("explicit Whisper decode config = %+v, want de/8/3/2", got)
	}

	for _, engine := range []string{"", "auto", string(types.ASRQwen3ASR06B), "qwen3-asr-1.7b-local"} {
		cfg.ASREngine = engine
		got := asrDecodeConfigFor(cfg)
		if got.Language != "auto" || got.NumThreads != 6 || got.BeamSize != 0 || got.BestOf != 0 {
			t.Fatalf("%q engine received wrong values: %+v, want legacy auto/6/0/0", engine, got)
		}
	}
}

func TestPipelineKeepsLegacyASRDefaultsForLocalLlama(t *testing.T) {
	p, set := newScriptedPipeline(providerBehavior{})
	cfg := testPipelineConfig(true, false)
	cfg.ASREngine = string(types.ASRQwen3ASR17BLocal)
	cfg.ASRLanguage = "de"
	cfg.ASRThreads = 2
	cfg.ASRBeamSize = 9
	cfg.ASRBestOf = 9

	if err := p.Initialize(context.Background(), cfg); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	t.Cleanup(func() { _ = p.Shutdown(context.Background()) })

	provider := set.asr[types.ASRQwen3ASR17BLocal]
	if provider.initCalls.Load() != 1 {
		t.Fatalf("local llama initialize calls = %d, want 1", provider.initCalls.Load())
	}
	if provider.config.Language != "auto" || provider.config.NumThreads != 6 ||
		provider.config.BeamSize != 0 || provider.config.BestOf != 0 {
		t.Fatalf("local llama config = language %q threads %d beam %d best-of %d, want legacy auto/6/0/0",
			provider.config.Language, provider.config.NumThreads, provider.config.BeamSize, provider.config.BestOf)
	}
}
