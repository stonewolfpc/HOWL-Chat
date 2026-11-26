package melotts

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

type blockingMeloTransport struct {
	active  int32
	max     int32
	entered chan struct{}
	release chan struct{}
}

func newBlockingMeloTransport() *blockingMeloTransport {
	return &blockingMeloTransport{
		entered: make(chan struct{}, 8),
		release: make(chan struct{}),
	}
}

func (b *blockingMeloTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	n := atomic.AddInt32(&b.active, 1)
	for {
		old := atomic.LoadInt32(&b.max)
		if n <= old || atomic.CompareAndSwapInt32(&b.max, old, n) {
			break
		}
	}
	b.entered <- struct{}{}
	select {
	case <-b.release:
	case <-req.Context().Done():
		atomic.AddInt32(&b.active, -1)
		return nil, req.Context().Err()
	}
	atomic.AddInt32(&b.active, -1)

	body := make([]byte, 64)
	copy(body[0:4], "RIFF")
	copy(body[8:12], "WAVE")
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     http.Header{"Content-Type": []string{"audio/wav"}},
		Body:       io.NopCloser(strings.NewReader(string(body))),
		Request:    req,
	}, nil
}

func meloUnitClient(t *testing.T) *Client {
	t.Helper()
	c := NewUnownedForTest()
	c.SetInitialized(types.TTSConfig{CacheDir: t.TempDir()})
	u, err := url.Parse("http://127.0.0.1:9999")
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	c.baseURL = u
	return c
}

func TestMeloTTSSerializesSynthesisAndVoices(t *testing.T) {
	c := meloUnitClient(t)
	transport := newBlockingMeloTransport()
	c.httpClient = &http.Client{Transport: transport}

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	startSynthesis := func() {
		defer wg.Done()
		_, err := c.Synthesize(context.Background(), types.SynthesisRequest{Text: "hello"}, nil)
		errs <- err
	}

	wg.Add(1)
	go startSynthesis()
	select {
	case <-transport.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("first synthesis did not start")
	}

	wg.Add(1)
	go startSynthesis()
	select {
	case <-transport.entered:
		t.Fatal("second synthesis entered while the first was active")
	case <-time.After(75 * time.Millisecond):
	}

	voicesDone := make(chan struct{})
	go func() {
		c.GetVoices()
		close(voicesDone)
	}()
	select {
	case <-voicesDone:
		t.Fatal("GetVoices returned while synthesis held the call lock")
	case <-time.After(75 * time.Millisecond):
	}

	close(transport.release)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("synthesis: %v", err)
		}
	}
	select {
	case <-voicesDone:
	case <-time.After(time.Second):
		t.Fatal("GetVoices did not return after synthesis released the lock")
	}
	if got := atomic.LoadInt32(&transport.max); got != 1 {
		t.Fatalf("max concurrent HTTP calls = %d, want 1", got)
	}
}

func TestMeloTTSStreamHoldsCallLockThroughRead(t *testing.T) {
	c := meloUnitClient(t)
	transport := newBlockingMeloTransport()
	transport.release = closedChan()
	c.httpClient = &http.Client{Transport: transport}

	readEntered := make(chan struct{})
	allowRead := make(chan struct{})
	c.readFile = func(path string) ([]byte, error) {
		close(readEntered)
		<-allowRead
		return os.ReadFile(path)
	}

	chunks := make(chan types.AudioChunk, 8)
	streamDone := make(chan error, 1)
	go func() {
		streamDone <- c.SynthesizeStream(context.Background(), types.SynthesisRequest{Text: "hello"}, chunks)
	}()
	go func() {
		for range chunks {
		}
	}()

	select {
	case <-transport.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("stream synthesis did not start")
	}
	select {
	case <-readEntered:
	case <-time.After(2 * time.Second):
		t.Fatal("stream did not reach the output read")
	}

	secondDone := make(chan error, 1)
	go func() {
		_, err := c.Synthesize(context.Background(), types.SynthesisRequest{Text: "second"}, nil)
		secondDone <- err
	}()
	select {
	case <-transport.entered:
		t.Fatal("second synthesis entered while the stream read held the call lock")
	case <-time.After(75 * time.Millisecond):
	}

	close(allowRead)
	select {
	case err := <-streamDone:
		if err != nil {
			t.Fatalf("stream: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("stream did not finish")
	}
	select {
	case err := <-secondDone:
		if err != nil {
			t.Fatalf("second synthesis: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("second synthesis did not finish")
	}
}

type recordedMeloRequest struct {
	path string
	body map[string]any
}

type meloRecordingTransport struct {
	mu       sync.Mutex
	requests []recordedMeloRequest
	fallback bool
}

func (r *meloRecordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var body map[string]any
	_ = json.NewDecoder(req.Body).Decode(&body)
	r.mu.Lock()
	r.requests = append(r.requests, recordedMeloRequest{path: req.URL.Path, body: body})
	r.mu.Unlock()

	if r.fallback && strings.HasSuffix(req.URL.Path, "/synthesize") {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "404 Not Found",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	}

	wav := make([]byte, 64)
	copy(wav[0:4], "RIFF")
	copy(wav[8:12], "WAVE")
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     http.Header{"Content-Type": []string{"audio/wav"}},
		Body:       io.NopCloser(strings.NewReader(string(wav))),
		Request:    req,
	}, nil
}

func TestMeloTTSLegacyPayloadSendsSpeedWithoutLengthScale(t *testing.T) {
	c := meloUnitClient(t)
	transport := &meloRecordingTransport{}
	c.httpClient = &http.Client{Transport: transport}

	if _, err := c.Synthesize(context.Background(), types.SynthesisRequest{Text: "hello", Speed: 1.25}, nil); err != nil {
		t.Fatalf("Synthesize: %v", err)
	}
	if len(transport.requests) != 1 {
		t.Fatalf("requests = %d, want 1 legacy request", len(transport.requests))
	}
	req := transport.requests[0]
	if !strings.HasSuffix(req.path, "/synthesize") {
		t.Fatalf("path = %q, want /synthesize", req.path)
	}
	if got := req.body["speed"]; got != 1.25 {
		t.Fatalf("legacy speed = %v, want 1.25", got)
	}
	if _, ok := req.body["length_scale"]; ok {
		t.Fatal("legacy /synthesize payload still contains length_scale")
	}
}

func TestMeloTTSOpenAIPayloadIsUnchanged(t *testing.T) {
	c := meloUnitClient(t)
	c.SetInitialized(types.TTSConfig{ModelName: "configured-model", CacheDir: t.TempDir()})
	transport := &meloRecordingTransport{fallback: true}
	c.httpClient = &http.Client{Transport: transport}

	if _, err := c.Synthesize(context.Background(), types.SynthesisRequest{Text: "hello", Speed: 1.25}, nil); err != nil {
		t.Fatalf("Synthesize: %v", err)
	}
	if len(transport.requests) != 2 {
		t.Fatalf("requests = %d, want legacy probe plus OpenAI request", len(transport.requests))
	}
	req := transport.requests[1]
	if !strings.HasSuffix(req.path, "/v1/audio/speech") {
		t.Fatalf("path = %q, want /v1/audio/speech", req.path)
	}
	if _, ok := req.body["length_scale"]; ok {
		t.Fatal("OpenAI-compatible payload unexpectedly contains length_scale")
	}
	for _, key := range []string{"model", "input", "voice", "response_format", "speed"} {
		if _, ok := req.body[key]; !ok {
			t.Fatalf("OpenAI-compatible payload missing %q: %#v", key, req.body)
		}
	}
	if got := req.body["speed"]; got != 1.25 {
		t.Fatalf("OpenAI-compatible speed = %v, want 1.25", got)
	}
	if got := req.body["model"]; got != "configured-model" {
		t.Fatalf("OpenAI-compatible model = %v, want configured-model", got)
	}
}

func closedChan() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}
