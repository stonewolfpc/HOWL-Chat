package kokoro

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

// kokoroTestConfig builds a TTS config from the environment so the integration
// tests run wherever the model set actually lives. Voices and config are
// derived from the model directory by Initialize, and the cache goes to a temp
// dir so a test run never writes into a real cache.
func kokoroTestConfig(t *testing.T) types.TTSConfig {
	t.Helper()
	model := os.Getenv("HOWL_TTS_KOKORO_MODEL")
	if model == "" {
		t.Skip("set HOWL_TTS_KOKORO_MODEL (optionally HOWL_TTS_KOKORO_VOICES / HOWL_TTS_KOKORO_CONFIG) to run Kokoro integration tests")
	}
	return types.TTSConfig{
		ModelPath:  model,
		VoicesPath: os.Getenv("HOWL_TTS_KOKORO_VOICES"),
		ConfigPath: os.Getenv("HOWL_TTS_KOKORO_CONFIG"),
		CacheDir:   t.TempDir(),
		ScriptPath: defaultKokoroScriptPath(),
	}
}

func TestKokoroListVoices(t *testing.T) {
	s := New()
	cfg := kokoroTestConfig(t)
	if err := s.Initialize(context.Background(), cfg); err != nil {
		t.Skipf("Kokoro initialization skipped (model not available): %v", err)
		return
	}
	voices := s.GetVoices()
	if len(voices) == 0 {
		t.Error("Expected voices from kokoro archive, got none")
	}
	for _, v := range voices[:5] {
		t.Logf("Voice: %s", v.ID)
	}
}

func TestKokoroSynthesize(t *testing.T) {
	s := New()
	cfg := kokoroTestConfig(t)
	if err := s.Initialize(context.Background(), cfg); err != nil {
		t.Skipf("Kokoro initialization skipped: %v", err)
		return
	}
	result, err := s.Synthesize(context.Background(), types.SynthesisRequest{
		Text:  "Hello, this is Kokoro TTS.",
		Voice: "af_heart",
		Speed: 1.0,
	}, nil)
	if err != nil {
		t.Fatalf("Synthesis failed: %v", err)
	}
	if result.AudioPath == "" {
		t.Error("Expected audio path, got empty string")
	}
	if _, err := os.Stat(result.AudioPath); err != nil {
		t.Errorf("Audio file not found: %v", err)
	}
	t.Logf("Synthesized: %s (%.1fs)", result.AudioPath, result.Duration)
}

func kokoroUnitSynth(t *testing.T) *Synthesizer {
	t.Helper()
	s := New()
	s.SetInitialized(types.TTSConfig{})
	s.pythonPath = "python"
	s.scriptPath = "script"
	s.modelPath = "model"
	s.voicesPath = "voices"
	s.configPath = "config"
	s.outputDir = t.TempDir()
	return s
}

func kokoroTestWAVRaw(path string) {
	data := make([]byte, 64)
	copy(data[0:4], "RIFF")
	copy(data[8:12], "WAVE")
	_ = os.WriteFile(path, data, 0o644)
}

func containsArg(args []string, target string) bool {
	for _, arg := range args {
		if arg == target {
			return true
		}
	}
	return false
}

type kokoroRunner struct {
	active  int32
	max     int32
	entered chan struct{}
	release chan struct{}
}

func newKokoroRunner() *kokoroRunner {
	return &kokoroRunner{entered: make(chan struct{}, 8), release: make(chan struct{})}
}

func (r *kokoroRunner) run(ctx context.Context, _ string, args []string, _ string) ([]byte, error) {
	if containsArg(args, "--list-voices") {
		return []byte("af_heart\nam_adam\n"), nil
	}

	output := ""
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "--output" {
			output = args[i+1]
			break
		}
	}
	if output == "" {
		return nil, fmt.Errorf("missing --output")
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
	kokoroTestWAVRaw(output)
	return nil, nil
}

func TestKokoroSerializesSynthesisAndVoices(t *testing.T) {
	s := kokoroUnitSynth(t)
	runner := newKokoroRunner()
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

func TestKokoroStreamHoldsCallLockThroughRead(t *testing.T) {
	s := kokoroUnitSynth(t)
	runner := newKokoroRunner()
	runner.release = closedKokoroChan()
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

func TestKokoroGetVoicesSerializedWithInitialize(t *testing.T) {
	dir := t.TempDir()
	model := filepath.Join(dir, "model.onnx")
	script := filepath.Join(dir, "kokoro-tts.py")
	voices := filepath.Join(dir, "voices-v1.0.bin")
	config := filepath.Join(dir, "config.json")
	for _, path := range []string{model, script, voices, config} {
		kokoroTestWAVRaw(path)
	}

	s := New()
	s.runCommand = func(_ context.Context, _ string, args []string, _ string) ([]byte, error) {
		if containsArg(args, "--list-voices") {
			return []byte("af_heart\n"), nil
		}
		return nil, fmt.Errorf("unexpected command: %v", args)
	}
	entered := make(chan struct{})
	release := make(chan struct{})
	s.resolvePython = func() (string, error) {
		close(entered)
		<-release
		return "python", nil
	}

	initDone := make(chan error, 1)
	go func() {
		initDone <- s.Initialize(context.Background(), types.TTSConfig{
			ModelPath:  model,
			ScriptPath: script,
			VoicesPath: voices,
			ConfigPath: config,
			CacheDir:   t.TempDir(),
		})
	}()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("initialize did not reach interpreter resolution")
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

func closedKokoroChan() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}
