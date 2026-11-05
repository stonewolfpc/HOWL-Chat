package fallback_test

import (
	"context"
	"testing"

	"howl-chat/internal/audio/asr/fallback"
	"howl-chat/internal/audio/runtime"
	"howl-chat/internal/audio/types"
)

func loadOwnershipProvider(t *testing.T, loader *runtime.ModelLoader, model types.ASRModelType, fake *ownershipRecognizer) {
	t.Helper()
	loader.RegisterASR(model, func() types.ASRRecognizer { return fake })
	if _, err := loader.LoadASR(context.Background(), string(model), types.ASRConfig{ModelPath: "fixture"}); err != nil {
		t.Fatalf("LoadASR %s: %v", model, err)
	}
}

func TestChainDoesNotInitializeLoaderProvider(t *testing.T) {
	loader := runtime.NewModelLoader()
	fake := &ownershipRecognizer{}
	loadOwnershipProvider(t, loader, types.ASRWhisperTiny, fake)

	chain := fallback.NewChain(fallback.ChainConfig{
		Recognizers: []types.ASRRecognizer{fake},
		MaxAttempts: 1,
	})
	if err := chain.Initialize(context.Background(), types.ASRConfig{ModelPath: "fixture"}); err != nil {
		t.Fatalf("chain Initialize: %v", err)
	}
	if fake.initCalls != 1 {
		t.Fatalf("initialize calls = %d, want 1 (loader only)", fake.initCalls)
	}
	if _, err := chain.Transcribe(context.Background(), "fixture", nil); err != nil {
		t.Fatalf("chain Transcribe: %v", err)
	}
	if fake.initCalls != 1 {
		t.Fatalf("initialize calls after transcribe = %d, want 1", fake.initCalls)
	}
	if fake.releaseCalls != 0 {
		t.Fatalf("release calls after transcribe = %d, want 0", fake.releaseCalls)
	}
}

func TestChainDoesNotReleaseLoaderProvider(t *testing.T) {
	loader := runtime.NewModelLoader()
	first := &ownershipRecognizer{}
	second := &ownershipRecognizer{}
	loadOwnershipProvider(t, loader, types.ASRWhisperTiny, first)
	loadOwnershipProvider(t, loader, types.ASRWhisperBase, second)

	chain := fallback.NewChain(fallback.ChainConfig{
		Recognizers: []types.ASRRecognizer{first, second},
		MaxAttempts: 1,
	})
	if err := chain.Initialize(context.Background(), types.ASRConfig{ModelPath: "fixture"}); err != nil {
		t.Fatalf("chain Initialize: %v", err)
	}
	if err := chain.Release(); err != nil {
		t.Fatalf("first chain Release: %v", err)
	}
	if err := chain.Release(); err != nil {
		t.Fatalf("second chain Release: %v", err)
	}
	if _, err := chain.Transcribe(context.Background(), "fixture", nil); err != nil {
		t.Fatalf("chain Transcribe: %v", err)
	}
	if first.releaseCalls != 0 || second.releaseCalls != 0 {
		t.Fatalf("chain released providers: first=%d second=%d", first.releaseCalls, second.releaseCalls)
	}

	if err := loader.Unload("asr:" + string(types.ASRWhisperTiny)); err != nil {
		t.Fatalf("unload first: %v", err)
	}
	if first.releaseCalls != 1 {
		t.Fatalf("first release calls = %d, want 1", first.releaseCalls)
	}
	if second.releaseCalls != 0 {
		t.Fatalf("second released early: %d", second.releaseCalls)
	}
	if err := loader.Unload("asr:" + string(types.ASRWhisperTiny)); err == nil {
		t.Fatal("repeated unload succeeded")
	}
	if first.releaseCalls != 1 {
		t.Fatalf("repeated unload released again: %d", first.releaseCalls)
	}

	if err := loader.Unload("asr:" + string(types.ASRWhisperBase)); err != nil {
		t.Fatalf("unload second: %v", err)
	}
	if second.releaseCalls != 1 {
		t.Fatalf("second release calls = %d, want 1", second.releaseCalls)
	}
	if err := loader.Unload("asr:" + string(types.ASRWhisperBase)); err == nil {
		t.Fatal("repeated second unload succeeded")
	}
	if second.releaseCalls != 1 {
		t.Fatalf("repeated second unload released again: %d", second.releaseCalls)
	}
}

func TestChainUnloadOrderIsIndependent(t *testing.T) {
	for _, order := range [][]int{{0, 1}, {1, 0}} {
		order := order
		t.Run(orderName(order), func(t *testing.T) {
			loader := runtime.NewModelLoader()
			first := &ownershipRecognizer{}
			second := &ownershipRecognizer{}
			models := []types.ASRModelType{types.ASRWhisperTiny, types.ASRWhisperBase}
			fakes := []*ownershipRecognizer{first, second}
			loadOwnershipProvider(t, loader, models[0], first)
			loadOwnershipProvider(t, loader, models[1], second)

			chain := fallback.NewChain(fallback.ChainConfig{
				Recognizers: []types.ASRRecognizer{first, second},
				MaxAttempts: 1,
			})
			if err := chain.Initialize(context.Background(), types.ASRConfig{ModelPath: "fixture"}); err != nil {
				t.Fatalf("chain Initialize: %v", err)
			}
			if err := chain.Release(); err != nil {
				t.Fatalf("chain Release: %v", err)
			}
			for _, index := range order {
				if err := loader.Unload("asr:" + string(models[index])); err != nil {
					t.Fatalf("unload index %d: %v", index, err)
				}
				if fakes[index].releaseCalls != 1 {
					t.Fatalf("provider %d release calls = %d, want 1", index, fakes[index].releaseCalls)
				}
				if err := loader.Unload("asr:" + string(models[index])); err == nil {
					t.Fatalf("repeated unload index %d succeeded", index)
				}
				if fakes[index].releaseCalls != 1 {
					t.Fatalf("provider %d released again on repeated unload: %d", index, fakes[index].releaseCalls)
				}
			}
		})
	}
}

func orderName(order []int) string {
	if len(order) == 2 && order[0] == 0 {
		return "first-then-second"
	}
	return "second-then-first"
}

type ownershipRecognizer struct {
	initCalls      int
	releaseCalls   int
	transcribeCall int
}

func (f *ownershipRecognizer) Initialize(context.Context, types.ASRConfig) error {
	f.initCalls++
	return nil
}

func (f *ownershipRecognizer) Transcribe(context.Context, string, types.ProgressCallback) (*types.RecognitionResult, error) {
	f.transcribeCall++
	return &types.RecognitionResult{Text: "not used"}, nil
}

func (f *ownershipRecognizer) TranscribeStream(context.Context, <-chan types.AudioChunk, chan<- types.RecognitionResult) error {
	return nil
}

func (f *ownershipRecognizer) SupportsLanguage(string) bool { return true }

func (f *ownershipRecognizer) GetModelInfo() types.ASRModelInfo {
	return types.ASRModelInfo{Name: "ownership fixture"}
}

func (f *ownershipRecognizer) Release() error {
	f.releaseCalls++
	return nil
}
