package asr

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

type serializedRecognizer struct {
	*BaseRecognizer
	active  int32
	max     int32
	entered chan struct{}
	release chan struct{}
}

func newSerializedRecognizer() *serializedRecognizer {
	return &serializedRecognizer{
		BaseRecognizer: NewBaseRecognizer(types.ASRWhisperTiny),
		entered:        make(chan struct{}, 4),
		release:        make(chan struct{}),
	}
}

func (r *serializedRecognizer) Initialize(_ context.Context, config types.ASRConfig) error {
	r.SetInitialized(config)
	return nil
}

func (r *serializedRecognizer) Transcribe(ctx context.Context, _ string, _ types.ProgressCallback) (*types.RecognitionResult, error) {
	r.LockCalls()
	defer r.UnlockCalls()

	n := atomic.AddInt32(&r.active, 1)
	for {
		old := atomic.LoadInt32(&r.max)
		if n <= old || atomic.CompareAndSwapInt32(&r.max, old, n) {
			break
		}
	}
	r.entered <- struct{}{}
	select {
	case <-r.release:
	case <-ctx.Done():
		atomic.AddInt32(&r.active, -1)
		return nil, ctx.Err()
	}
	atomic.AddInt32(&r.active, -1)
	return &types.RecognitionResult{Text: "ok"}, nil
}

func (r *serializedRecognizer) TranscribeStream(context.Context, <-chan types.AudioChunk, chan<- types.RecognitionResult) error {
	return nil
}

func (r *serializedRecognizer) SupportsLanguage(string) bool { return true }

func (r *serializedRecognizer) GetModelInfo() types.ASRModelInfo {
	return types.ASRModelInfo{Name: "serialized fake"}
}

func (r *serializedRecognizer) Release() error { return nil }

func TestBaseRecognizerSerializesCalls(t *testing.T) {
	r := newSerializedRecognizer()
	if err := r.Initialize(context.Background(), types.ASRConfig{}); err != nil {
		t.Fatalf("initialize: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	transcribe := func() {
		defer wg.Done()
		if _, err := r.Transcribe(context.Background(), "fixture", nil); err != nil {
			t.Errorf("transcribe: %v", err)
		}
	}
	go transcribe()
	select {
	case <-r.entered:
	case <-time.After(time.Second):
		t.Fatal("first call did not enter")
	}
	go transcribe()
	select {
	case <-r.entered:
		t.Fatal("second call entered while first was active")
	case <-time.After(75 * time.Millisecond):
	}
	if got := atomic.LoadInt32(&r.max); got != 1 {
		t.Fatalf("max concurrent calls = %d, want 1", got)
	}
	close(r.release)
	wg.Wait()
}
