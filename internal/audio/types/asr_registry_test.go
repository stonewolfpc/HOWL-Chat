package types

import "testing"

// DEBT-T2: a model may not advertise a streaming mode its adapter rejects.
// Every current adapter returns "not supported" from TranscribeStream, so every
// registry entry must declare false.
func TestASRRegistryStreamingMatchesAdapters(t *testing.T) {
	adapters := map[ASRModelType]bool{
		ASRWhisperTiny:      false, // whisper-cli has no streaming path
		ASRWhisperBase:      false,
		ASRWhisperSmall:     false,
		ASRQwen3ASR06B:      false, // DashScope async task API cannot stream
		ASRQwen3ASR17B:      false,
		ASRQwen3ASR17BLocal: false, // file transcription only; no streaming adapter
	}

	for modelType, profile := range ASRModelRegistry {
		adapterStreams, known := adapters[modelType]
		if !known {
			t.Fatalf("registry entry %q has no adapter streaming declaration", modelType)
		}
		if profile.SupportsStreaming != adapterStreams {
			t.Fatalf("%q declares SupportsStreaming=%t but its adapter streams=%t",
				modelType, profile.SupportsStreaming, adapterStreams)
		}
	}
	if len(adapters) != len(ASRModelRegistry) {
		t.Fatalf("adapter declarations = %d, registry entries = %d", len(adapters), len(ASRModelRegistry))
	}
}

// Every entry must declare the facts the capability contract reports, so no
// dimension is guessed from a name at runtime.
func TestASRRegistryDeclaresCapabilityFacts(t *testing.T) {
	for modelType, profile := range ASRModelRegistry {
		t.Run(string(modelType), func(t *testing.T) {
			if profile.Runtime == "" {
				t.Fatal("no runtime declared")
			}
			if profile.Verification == "" {
				t.Fatal("no verification grade declared")
			}
			switch profile.Verification {
			case VerificationUnverified, VerificationImplementation, VerificationEndToEnd:
			default:
				t.Fatalf("unknown verification grade %q", profile.Verification)
			}
			if profile.Local && profile.Backend != "local" {
				t.Fatalf("local model declares backend %q", profile.Backend)
			}
			if !profile.Local && profile.Backend == "local" {
				t.Fatal("remote model declares a local backend")
			}
			if profile.Local && profile.ApprovedLocal && profile.LocalFile == "" {
				t.Fatal("admitted local model declares no model file")
			}
		})
	}
}

// The local GGUF candidate must not be the same identity as the DashScope
// provider, and is admitted only after E4 and with a registered adapter.
func TestLocalQwenIsSeparateAndAdmitted(t *testing.T) {
	remote, ok := ASRModelRegistry[ASRQwen3ASR17B]
	if !ok {
		t.Fatal("remote Qwen 1.7B entry is missing")
	}
	local, ok := ASRModelRegistry[ASRQwen3ASR17BLocal]
	if !ok {
		t.Fatal("local Qwen 1.7B entry is missing")
	}
	if remote.Local || remote.Backend != "http" {
		t.Fatalf("remote entry is not a remote provider: %+v", remote.ModelProfile)
	}
	if !local.Local || local.Backend != "local" {
		t.Fatalf("local entry is not local: %+v", local.ModelProfile)
	}
	if !local.ApprovedLocal {
		t.Fatal("local Qwen is not admitted after the E4 transcription")
	}
	if local.LocalFile == "" || local.ProjectorFile == "" {
		t.Fatalf("local Qwen must name its model and projector: %+v", local.ModelProfile)
	}
	if local.Verification != VerificationEndToEnd {
		t.Fatalf("local Qwen verification = %q, want end_to_end", local.Verification)
	}
	if local.ProjectorKind != ProjectorKindAudio {
		t.Fatalf("local Qwen projector kind = %q, want audio", local.ProjectorKind)
	}
	if local.SupportsStreaming {
		t.Fatal("local Qwen advertises streaming")
	}
}

// Automatic recommendation must never point at a model that cannot run here.
func TestRecommendASRReturnsOnlyAdmittedLocalModels(t *testing.T) {
	for _, tier := range []MemoryTier{TierTiny, TierSmall, TierMedium, TierLarge} {
		model := RecommendASR(tier, nil)
		profile, ok := ASRModelRegistry[model]
		if !ok {
			t.Fatalf("tier %q recommended unknown model %q", tier, model)
		}
		if !profile.Local || !profile.ApprovedLocal {
			t.Fatalf("tier %q recommended non-admitted model %q", tier, model)
		}
	}
}
