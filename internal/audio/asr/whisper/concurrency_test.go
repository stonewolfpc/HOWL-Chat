package whisper

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

func whisperOutputPath(args []string) string {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "--file" {
			return args[i+1] + ".json"
		}
	}
	return ""
}

func writeWhisperJSON(path string) error {
	const body = `{"result":{"language":"en"},"transcription":[{"timestamps":{"from":"00:00:00,000","to":"00:00:01,000"},"text":"hello"}]}`
	return os.WriteFile(path, []byte(body), 0o644)
}

type blockingWhisperRunner struct {
	active  int32
	max     int32
	entered chan struct{}
	release chan struct{}
}

func newBlockingWhisperRunner() *blockingWhisperRunner {
	return &blockingWhisperRunner{
		entered: make(chan struct{}, 8),
		release: make(chan struct{}),
	}
}

func (b *blockingWhisperRunner) run(ctx context.Context, args []string) ([]byte, error) {
	output := whisperOutputPath(args)
	if output == "" {
		return nil, fmt.Errorf("missing --file argument")
	}

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
	case <-ctx.Done():
		atomic.AddInt32(&b.active, -1)
		return nil, ctx.Err()
	}
	atomic.AddInt32(&b.active, -1)
	if err := writeWhisperJSON(output); err != nil {
		return nil, err
	}
	return nil, nil
}

func TestWhisperSerializesConcurrentTranscriptions(t *testing.T) {
	model := whisperTestModel(t)
	cli := whisperTestCLI(t)
	audio := writeWhisperWAV(t)
	r := newWhisperTestRecognizer(t, cli, NewLimiter(4))
	if err := r.Initialize(context.Background(), types.ASRConfig{ModelPath: model}); err != nil {
		t.Fatalf("initialize: %v", err)
	}

	runner := newBlockingWhisperRunner()
	r.runCLI = runner.run

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	transcribe := func() {
		defer wg.Done()
		_, err := r.Transcribe(context.Background(), audio, nil)
		errs <- err
	}

	wg.Add(1)
	go transcribe()
	select {
	case <-runner.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("first transcription did not start")
	}

	wg.Add(1)
	go transcribe()
	select {
	case <-runner.entered:
		t.Fatal("second transcription entered while the first was active")
	case <-time.After(100 * time.Millisecond):
	}
	if got := atomic.LoadInt32(&runner.max); got != 1 {
		t.Fatalf("max concurrent provider calls = %d, want 1", got)
	}

	close(runner.release)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("transcription: %v", err)
		}
	}
	if got := atomic.LoadInt32(&runner.max); got != 1 {
		t.Fatalf("max concurrent provider calls = %d, want 1", got)
	}
}

func TestWhisperLimiterSerializesAcrossRecognizers(t *testing.T) {
	model := whisperTestModel(t)
	cli := whisperTestCLI(t)
	audio := writeWhisperWAV(t)
	limiter := NewLimiter(1)

	newRecognizer := func() *Recognizer {
		r := newWhisperTestRecognizer(t, cli, limiter)
		if err := r.Initialize(context.Background(), types.ASRConfig{ModelPath: model}); err != nil {
			t.Fatalf("initialize: %v", err)
		}
		return r
	}
	r1, r2 := newRecognizer(), newRecognizer()

	runner := newBlockingWhisperRunner()
	r1.runCLI = runner.run
	r2.runCLI = runner.run

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := r1.Transcribe(context.Background(), audio, nil)
		errs <- err
	}()
	select {
	case <-runner.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("first recognizer did not start")
	}
	go func() {
		defer wg.Done()
		_, err := r2.Transcribe(context.Background(), audio, nil)
		errs <- err
	}()

	select {
	case <-runner.entered:
		t.Fatal("second recognizer entered while the limiter slot was held")
	case <-time.After(100 * time.Millisecond):
	}
	if got := atomic.LoadInt32(&runner.max); got != 1 {
		t.Fatalf("max concurrent whisper calls = %d, want 1", got)
	}

	close(runner.release)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("transcription: %v", err)
		}
	}
	if got := atomic.LoadInt32(&runner.max); got != 1 {
		t.Fatalf("max concurrent whisper calls = %d, want 1", got)
	}
}

func TestWhisperLimiterLifecycle(t *testing.T) {
	if got := defaultLimiter.Limit(); got != 1 {
		t.Fatalf("default limiter = %d, want 1", got)
	}
	l := NewLimiter(2)
	if got := l.Limit(); got != 2 {
		t.Fatalf("limit = %d, want 2", got)
	}
	if err := l.SetLimit(3); err != nil {
		t.Fatalf("set limit while idle: %v", err)
	}
	if err := l.Acquire(context.Background()); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	if err := l.SetLimit(1); err == nil {
		t.Fatal("SetLimit succeeded while a call was active")
	}
	l.Release()
	if err := l.SetLimit(1); err != nil {
		t.Fatalf("set limit after release: %v", err)
	}
	if got := l.Limit(); got != 1 {
		t.Fatalf("limit = %d, want 1", got)
	}
}
