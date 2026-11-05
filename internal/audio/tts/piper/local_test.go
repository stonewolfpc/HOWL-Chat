package piper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

func piperTestWAV(t *testing.T, path string) {
	t.Helper()
	data := make([]byte, 64)
	copy(data[0:4], "RIFF")
	copy(data[8:12], "WAVE")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write wav: %v", err)
	}
}

func piperUnitSynth(t *testing.T) *Synthesizer {
	t.Helper()
	s := New()
	s.SetInitialized(types.TTSConfig{})
	s.piperPath = "piper"
	s.modelPath = "model"
	s.outputDir = t.TempDir()
	return s
}

type piperRunner struct {
	active  int32
	max     int32
	entered chan struct{}
	release chan struct{}
}

func newPiperRunner() *piperRunner {
	return &piperRunner{entered: make(chan struct{}, 8), release: make(chan struct{})}
}

func (r *piperRunner) run(ctx context.Context, _ string, args []string, _ string) ([]byte, error) {
	output := ""
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "--output_file" {
			output = args[i+1]
			break
		}
	}
	if output == "" {
		return nil, fmt.Errorf("missing --output_file")
	}

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
	piperTestWAVRaw(output)
	return nil, nil
}

func piperTestWAVRaw(path string) {
	data := make([]byte, 64)
	copy(data[0:4], "RIFF")
	copy(data[8:12], "WAVE")
	_ = os.WriteFile(path, data, 0o644)
}

func TestPiperSerializesSynthesisAndVoices(t *testing.T) {
	s := piperUnitSynth(t)
	runner := newPiperRunner()
	s.runCommand = runner.run

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	startSynthesis := func() {
		defer wg.Done()
		_, err := s.Synthesize(context.Background(), types.SynthesisRequest{Text: "hello"}, nil)
		errs <- err
	}

	wg.Add(1)
	go startSynthesis()
	select {
	case <-runner.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("first synthesis did not start")
	}

	wg.Add(1)
	go startSynthesis()
	select {
	case <-runner.entered:
		t.Fatal("second synthesis entered while the first was active")
	case <-time.After(75 * time.Millisecond):
	}

	voicesDone := make(chan struct{})
	go func() {
		s.GetVoices()
		close(voicesDone)
	}()
	select {
	case <-voicesDone:
		t.Fatal("GetVoices returned while synthesis held the call lock")
	case <-time.After(75 * time.Millisecond):
	}

	close(runner.release)
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
	if got := atomic.LoadInt32(&runner.max); got != 1 {
		t.Fatalf("max concurrent provider calls = %d, want 1", got)
	}
}

func TestPiperStreamHoldsCallLockThroughRead(t *testing.T) {
	s := piperUnitSynth(t)
	runner := newPiperRunner()
	runner.release = closedChannel()
	s.runCommand = runner.run

	readEntered := make(chan struct{})
	allowRead := make(chan struct{})
	s.readFile = func(path string) ([]byte, error) {
		close(readEntered)
		<-allowRead
		return os.ReadFile(path)
	}

	chunks := make(chan types.AudioChunk, 8)
	streamDone := make(chan error, 1)
	go func() {
		streamDone <- s.SynthesizeStream(context.Background(), types.SynthesisRequest{Text: "hello"}, chunks)
	}()
	go func() {
		for range chunks {
		}
	}()

	select {
	case <-runner.entered:
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
		_, err := s.Synthesize(context.Background(), types.SynthesisRequest{Text: "second"}, nil)
		secondDone <- err
	}()
	select {
	case <-runner.entered:
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

func TestPiperGetVoicesSerializedWithInitialize(t *testing.T) {
	s := New()
	model := filepath.Join(t.TempDir(), "model.onnx")
	piperTestWAVRaw(model)
	entered := make(chan struct{})
	release := make(chan struct{})
	s.resolveBinary = func() (string, error) {
		close(entered)
		<-release
		return "piper", nil
	}

	initDone := make(chan error, 1)
	go func() {
		initDone <- s.Initialize(context.Background(), types.TTSConfig{
			ModelPath: model,
			CacheDir:  t.TempDir(),
		})
	}()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("initialize did not reach binary resolution")
	}

	voicesDone := make(chan struct{})
	go func() {
		s.GetVoices()
		close(voicesDone)
	}()
	select {
	case <-voicesDone:
		t.Fatal("GetVoices returned while Initialize held the call lock")
	case <-time.After(75 * time.Millisecond):
	}

	close(release)
	select {
	case err := <-initDone:
		if err != nil {
			t.Fatalf("initialize: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("initialize did not finish")
	}
	select {
	case <-voicesDone:
	case <-time.After(time.Second):
		t.Fatal("GetVoices did not return after initialize")
	}
}

func closedChannel() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}
