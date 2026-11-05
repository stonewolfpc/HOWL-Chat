package tts_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"howl-chat/internal/audio/asr"
	"howl-chat/internal/audio/tts"
	"howl-chat/internal/audio/types"
)

type blockingASR struct {
	*asr.BaseRecognizer
	entered chan struct{}
	release chan struct{}
}

func (f *blockingASR) Initialize(_ context.Context, config types.ASRConfig) error {
	f.SetInitialized(config)
	return nil
}

func (f *blockingASR) Transcribe(context.Context, string, types.ProgressCallback) (*types.RecognitionResult, error) {
	f.LockCalls()
	defer f.UnlockCalls()
	f.entered <- struct{}{}
	<-f.release
	return &types.RecognitionResult{Text: "asr"}, nil
}

func (f *blockingASR) TranscribeStream(context.Context, <-chan types.AudioChunk, chan<- types.RecognitionResult) error {
	return nil
}

func (f *blockingASR) SupportsLanguage(string) bool { return true }

func (f *blockingASR) GetModelInfo() types.ASRModelInfo {
	return types.ASRModelInfo{Name: "blocking ASR"}
}

func (f *blockingASR) Release() error { return nil }

type blockingTTS struct {
	*tts.BaseSynthesizer
	entered chan struct{}
	release chan struct{}
}

func (f *blockingTTS) Initialize(_ context.Context, config types.TTSConfig) error {
	f.SetInitialized(config)
	return nil
}

func (f *blockingTTS) Synthesize(context.Context, types.SynthesisRequest, types.ProgressCallback) (*types.SynthesisResult, error) {
	f.LockCalls()
	defer f.UnlockCalls()
	f.entered <- struct{}{}
	<-f.release
	return &types.SynthesisResult{AudioPath: "tts.wav"}, nil
}

func (f *blockingTTS) SynthesizeStream(context.Context, types.SynthesisRequest, chan<- types.AudioChunk) error {
	return nil
}

func (f *blockingTTS) GetVoices() []types.Voice { return nil }

func (f *blockingTTS) SupportsSSML() bool { return false }

func (f *blockingTTS) Release() error { return nil }

func TestASRAndTTSRunConcurrently(t *testing.T) {
	release := make(chan struct{})
	asrFixture := &blockingASR{
		BaseRecognizer: asr.NewBaseRecognizer(types.ASRWhisperTiny),
		entered:        make(chan struct{}, 1),
		release:        release,
	}
	ttsFixture := &blockingTTS{
		BaseSynthesizer: tts.NewBaseSynthesizer(types.TTSPiper),
		entered:         make(chan struct{}, 1),
		release:         release,
	}
	defer func() {
		select {
		case <-release:
		default:
			close(release)
		}
	}()

	if err := asrFixture.Initialize(context.Background(), types.ASRConfig{}); err != nil {
		t.Fatalf("initialize ASR: %v", err)
	}
	if err := ttsFixture.Initialize(context.Background(), types.TTSConfig{}); err != nil {
		t.Fatalf("initialize TTS: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if _, err := asrFixture.Transcribe(context.Background(), "fixture.wav", nil); err != nil {
			t.Errorf("ASR: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := ttsFixture.Synthesize(context.Background(), types.SynthesisRequest{Text: "hello"}, nil); err != nil {
			t.Errorf("TTS: %v", err)
		}
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-asrFixture.entered:
		case <-ttsFixture.entered:
		case <-time.After(2 * time.Second):
			t.Fatal("ASR and TTS calls did not both enter; they serialized each other")
		}
	}

	close(release)
	wg.Wait()
}
