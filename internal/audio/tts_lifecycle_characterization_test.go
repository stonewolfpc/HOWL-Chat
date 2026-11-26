package audio

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/audio/tts/cache"
	"howl-chat/internal/audio/types"
	"howl-chat/internal/resource"
)

// TestTTSStartConsultsAdmissionWhileChatLoaded is the TTS counterpart of the
// ASR admission baseline. It proves the pipeline asks the shared
// resource.Authority before constructing or initializing any TTS provider, and
// that a refusal names the resident model instead of failing anonymously.
//
// This is a characterization of the CURRENT policy, not a policy change: the
// fixture is chosen so the documented policy genuinely refuses (a 300 MB Piper
// RAM request against 180 MB available after the 10% reserve of 200 MB free).
// The ASR baseline fixture still admits; its disposition remains open.
func TestTTSStartConsultsAdmissionWhileChatLoaded(t *testing.T) {
	const mb = int64(1024 * 1024)
	loaded := []resource.ModelState{{
		Type: resource.ResourceLLM, Name: "chat", Loaded: true, GPUMemoryMB: 3000, SystemRAMMB: 30000,
	}}
	authority := resource.NewAuthority(
		func() (resource.HardwareSnapshot, error) {
			return resource.HardwareSnapshot{
				GPUModel:       "Test GPU",
				TotalVRAMBytes: 24000 * mb,
				FreeVRAMBytes:  4500 * mb,
				UsedVRAMBytes:  19500 * mb,
				TotalRAMBytes:  32 * 1024 * mb,
				FreeRAMBytes:   200 * mb,
				VRAMLive:       true,
				VRAMSource:     "test-live",
				GPUAvailable:   true,
			}, nil
		},
		func() []resource.ModelState { return loaded },
	)

	p := NewPipeline()
	p.admission = authority
	_, _, _, outcome := p.initTTS(context.Background(), []types.TTSModelType{types.TTSPiper})
	if outcome.ready {
		t.Fatal("TTS must not become ready when admission refuses the start")
	}
	if outcome.err == nil || !strings.Contains(outcome.err.Error(), "blocking model: llm \"chat\"") {
		t.Fatalf("TTS refusal must name the chat model: %v", outcome.err)
	}
	// The refusal must also precede provider ownership: no synthesizer may be
	// installed and no runtime model may be tracked after a refused start.
	if p.GetStatus().TTSReady {
		t.Fatal("TTS status must stay not-ready after a refused start")
	}
	if models := p.loader.GetLoadedModels(); len(models) != 0 {
		t.Fatalf("a refused TTS start must not leave a loaded provider: %v", models)
	}
}

// blockingTTSSynth stands in for a provider that owns a child process. It only
// returns after its caller cancels, so a shutdown that reports completion
// without joining the synthesis is observable as an early return.
type blockingTTSSynth struct {
	entered   chan struct{}
	finished  chan struct{}
	once      sync.Once
	cancelled atomic.Bool
}

func newBlockingTTSSynth() *blockingTTSSynth {
	return &blockingTTSSynth{entered: make(chan struct{}), finished: make(chan struct{})}
}

func (s *blockingTTSSynth) Initialize(context.Context, types.TTSConfig) error { return nil }

func (s *blockingTTSSynth) Synthesize(ctx context.Context, _ types.SynthesisRequest, _ types.ProgressCallback) (*types.SynthesisResult, error) {
	s.once.Do(func() { close(s.entered) })
	<-ctx.Done()
	s.cancelled.Store(true)
	close(s.finished)
	return nil, ctx.Err()
}

func (s *blockingTTSSynth) SynthesizeStream(context.Context, types.SynthesisRequest, chan<- types.AudioChunk) error {
	return nil
}

func (s *blockingTTSSynth) GetVoices() []types.Voice {
	return []types.Voice{{ID: "en_US-lessac-medium", Name: "en_US-lessac-medium", Language: "en"}}
}

func (s *blockingTTSSynth) SupportsSSML() bool { return false }
func (s *blockingTTSSynth) Release() error     { return nil }

// TestManagedAudioCloseCancelsInFlightSynthesis proves the audio execution
// boundary is the same shape as the image workflow's: Close cancels the active
// operation's context, the cancellation reaches the provider call inside the
// pipeline, and Close does not report completion before that call returns.
//
// It does NOT prove child-process termination. Piper and Kokoro spawn their
// process with exec.CommandContext outside job ownership (see the HRA-P1
// implementation plan §7.6.1); that gap is characterized, not closed, here.
func TestManagedAudioCloseCancelsInFlightSynthesis(t *testing.T) {
	pipeline := NewPipeline()
	synth := newBlockingTTSSynth()
	pipeline.mu.Lock()
	pipeline.config = types.PipelineConfig{TTSEnabled: true, TTSEngine: "piper"}
	pipeline.ttsSynth = synth
	pipeline.activeTTSEngine = "piper"
	pipeline.ttsCache = cache.NewManager(10)
	pipeline.status.TTSReady = true
	pipeline.mu.Unlock()

	manager := NewManager(func(_ context.Context, _ Settings, candidate *Candidate) error {
		candidate.Pipeline = pipeline
		candidate.TTS = CapabilityResult{Ready: true}
		return nil
	})
	if _, _, err := manager.Apply(context.Background(), Settings{TTSEnabled: true, TTSEngine: "piper"}); err != nil {
		t.Fatalf("apply TTS generation: %v", err)
	}

	op, opCtx := manager.BeginOperation(context.Background(), 0)
	lease, err := manager.Acquire(opCtx)
	if err != nil {
		op.End()
		t.Fatalf("acquire TTS generation: %v", err)
	}

	requestDone := make(chan struct{})
	go func() {
		defer close(requestDone)
		_, _ = lease.Pipeline().ProcessAudioMessage(opCtx, types.AudioMessage{
			Type: "audio_output",
			Text: "hello from the audio lifecycle test",
		})
		lease.Release()
		op.End()
	}()

	select {
	case <-synth.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("synthesis did not start")
	}

	closeDone := make(chan error, 1)
	go func() { closeDone <- manager.Close(context.Background()) }()

	var closeErr error
	select {
	case closeErr = <-closeDone:
	case <-time.After(5 * time.Second):
		t.Fatal("Close did not complete after cancelling the in-flight synthesis")
	}
	if closeErr != nil {
		t.Fatalf("Close after joined synthesis: %v", closeErr)
	}
	// The ordering guarantee: when Close reports completion, the provider call
	// has already returned. A Close that reported success without cancelling
	// would never reach this branch because the request would still be blocked
	// until the 5 s timeout above.
	select {
	case <-synth.finished:
	default:
		t.Fatal("Close reported completion before the in-flight synthesis returned")
	}

	select {
	case <-requestDone:
	case <-time.After(2 * time.Second):
		t.Fatal("audio request did not finish releasing its lease")
	}
	if !synth.cancelled.Load() {
		t.Fatal("provider did not observe context cancellation")
	}
	if status := manager.Status(); status.InFlight != 0 {
		t.Fatalf("in-flight audio operations after close = %d, want 0", status.InFlight)
	}
}
