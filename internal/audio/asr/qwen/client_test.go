package qwen

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

type blockingQwenTransport struct {
	active  int32
	max     int32
	entered chan struct{}
}

func newBlockingQwenTransport() *blockingQwenTransport {
	return &blockingQwenTransport{entered: make(chan struct{}, 4)}
}

func (b *blockingQwenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	n := atomic.AddInt32(&b.active, 1)
	for {
		old := atomic.LoadInt32(&b.max)
		if n <= old || atomic.CompareAndSwapInt32(&b.max, old, n) {
			break
		}
	}
	b.entered <- struct{}{}
	<-req.Context().Done()
	atomic.AddInt32(&b.active, -1)
	return nil, req.Context().Err()
}

func TestQwenSerializesConcurrentTranscriptions(t *testing.T) {
	client := NewClient("test-key")
	transport := newBlockingQwenTransport()
	client.httpClient = &http.Client{Transport: transport}
	if err := client.Initialize(context.Background(), types.ASRConfig{}); err != nil {
		t.Fatalf("initialize: %v", err)
	}

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	firstErr := make(chan error, 1)
	go func() {
		_, err := client.Transcribe(ctx1, "audio.wav", nil)
		firstErr <- err
	}()
	select {
	case <-transport.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("first transcription did not start")
	}

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	secondErr := make(chan error, 1)
	go func() {
		_, err := client.Transcribe(ctx2, "audio.wav", nil)
		secondErr <- err
	}()

	select {
	case <-transport.entered:
		t.Fatal("second transcription entered while the first was active")
	case <-time.After(100 * time.Millisecond):
	}
	if got := atomic.LoadInt32(&transport.max); got != 1 {
		t.Fatalf("max concurrent provider calls = %d, want 1", got)
	}

	cancel1()
	select {
	case err := <-firstErr:
		if err == nil {
			t.Fatal("first transcription succeeded after cancellation")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("first transcription did not return")
	}

	select {
	case <-transport.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("second transcription did not enter after the first released")
	}
	if got := atomic.LoadInt32(&transport.max); got != 1 {
		t.Fatalf("max concurrent provider calls = %d, want 1", got)
	}

	cancel2()
	select {
	case err := <-secondErr:
		if err == nil {
			t.Fatal("second transcription succeeded after cancellation")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("second transcription did not return")
	}
}

func TestQwenInitializeDoesNotRequireNetwork(t *testing.T) {
	client := NewClient("test-key")
	client.httpClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		t.Fatalf("unexpected HTTP request during initialize: %s", req.URL)
		return nil, nil
	})}
	if err := client.Initialize(context.Background(), types.ASRConfig{}); err != nil {
		t.Fatalf("initialize: %v", err)
	}
	if !client.IsInitialized() {
		t.Fatal("client not initialized")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
