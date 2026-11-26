package audio

import (
	"context"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/audio/tts/cache"
	"howl-chat/internal/audio/tts/melotts"
	"howl-chat/internal/audio/types"
	"howl-chat/internal/process"
	"howl-chat/internal/resource"
)

type stubTTSSynth struct {
	voices []types.Voice
}

func (s *stubTTSSynth) Initialize(context.Context, types.TTSConfig) error { return nil }

func (s *stubTTSSynth) Synthesize(context.Context, types.SynthesisRequest, types.ProgressCallback) (*types.SynthesisResult, error) {
	return &types.SynthesisResult{AudioPath: "stub.wav"}, nil
}

func (s *stubTTSSynth) SynthesizeStream(context.Context, types.SynthesisRequest, chan<- types.AudioChunk) error {
	return nil
}

func (s *stubTTSSynth) GetVoices() []types.Voice { return s.voices }
func (s *stubTTSSynth) SupportsSSML() bool       { return false }
func (s *stubTTSSynth) Release() error           { return nil }

func permissiveAuthority() *resource.Authority {
	const mb = int64(1024 * 1024)
	return resource.NewAuthority(
		func() (resource.HardwareSnapshot, error) {
			return resource.HardwareSnapshot{
				GPUModel:       "Test GPU",
				TotalVRAMBytes: 24000 * mb,
				FreeVRAMBytes:  20000 * mb,
				UsedVRAMBytes:  4000 * mb,
				TotalRAMBytes:  64 * 1024 * mb,
				FreeRAMBytes:   32 * 1024 * mb,
				VRAMLive:       true,
				VRAMSource:     "test-live",
				GPUAvailable:   true,
			}, nil
		},
		func() []resource.ModelState { return resource.GetManager().LoadedModels() },
	)
}

// installResidentProvider performs what initTTS does once a provider is loaded:
// admit, convert the lease into a registration, and install the provider as the
// active one. It returns the provider and its execution release func.
func installResidentProvider(t *testing.T, p *Pipeline, key, engine string, synth types.TTSSynthesizer, register bool) *ttsProvider {
	t.Helper()
	provider := &ttsProvider{key: key, synth: synth, engine: engine, modelName: key}
	if register {
		req, err := ttsAdmissionRequest(engine, "")
		if err != nil {
			t.Fatalf("admission request: %v", err)
		}
		lease, err := p.admitAudio(req)
		if err != nil {
			t.Fatalf("admit %s: %v", engine, err)
		}
		provider.registered = p.registerAudio(resource.ResourceTTS, key, "", req, lease)
	}
	p.mu.Lock()
	p.registerTTSResidencyLocked(provider)
	p.ttsSynth = synth
	p.activeTTSEngine = engine
	p.activeTTSKey = key
	p.globalTTSKey = key
	p.mu.Unlock()
	return provider
}

// TestTTSResidencySurvivesInitializationAndShutdownIsIdempotent proves a
// provider's residency is a registration, not a returned initialization call,
// and that shutdown releases it exactly once.
func TestTTSResidencySurvivesInitializationAndShutdownIsIdempotent(t *testing.T) {
	registry := process.NewRegistry()
	p := NewPipeline(registry)
	p.admission = permissiveAuthority()
	p.ttsCache = cache.NewManager(10)

	provider := installResidentProvider(t, p, ttsEngineKey("piper", ""), "piper", &stubTTSSynth{}, true)
	if p.activeProvider() != provider {
		t.Fatal("the installed provider is not the active provider")
	}
	if got := resource.GetManager().ModelsByType(resource.ResourceTTS); len(got) == 0 {
		t.Fatal("provider residency was released when initialization returned")
	}

	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("repeated shutdown must be idempotent: %v", err)
	}
	if got := resource.GetManager().ModelsByType(resource.ResourceTTS); len(got) != 0 {
		t.Fatalf("shutdown left registrations behind: %v", got)
	}
	if len(p.ttsProviders) != 0 {
		t.Fatalf("shutdown left providers behind: %d", len(p.ttsProviders))
	}
}

// TestTTSReclamationRefusesActiveProvider proves an executing provider is never
// reclaimed, and that shared references keep it valid until the last release.
func TestTTSReclamationRefusesActiveProvider(t *testing.T) {
	p := NewPipeline(process.NewRegistry())
	provider := &ttsProvider{key: "piper|shared", synth: &stubTTSSynth{}, engine: "piper"}
	p.mu.Lock()
	p.ttsProviders[provider.key] = provider
	p.mu.Unlock()

	first := p.acquireTTSExecution(provider)
	second := p.acquireTTSExecution(provider)

	// The guard itself refuses, naming the active executions.
	if err := p.unloadTTSProvider(provider); err == nil || !strings.Contains(err.Error(), "executing") {
		t.Fatalf("unload of an executing provider must be refused: %v", err)
	}
	// Bulk reclamation skips it without failing unrelated work.
	if err := p.ReclaimIdleTTSProviders(); err != nil {
		t.Fatalf("bulk reclamation must skip executing providers, not fail: %v", err)
	}
	if _, ok := p.ttsProviders[provider.key]; !ok {
		t.Fatal("an executing provider was reclaimed")
	}

	first()
	if err := p.unloadTTSProvider(provider); err == nil {
		t.Fatal("unload succeeded while a second execution lease was still held")
	}

	second()
	if err := p.ReclaimIdleTTSProviders(); err != nil {
		t.Fatalf("an idle provider must be reclaimable: %v", err)
	}
	if _, ok := p.ttsProviders[provider.key]; ok {
		t.Fatal("reclaimed provider is still resident")
	}
}

type blockingEnumerator struct {
	entered   chan struct{}
	finished  chan struct{}
	cancelled atomic.Bool
	once      sync.Once
}

func (b *blockingEnumerator) Initialize(context.Context, types.TTSConfig) error { return nil }

func (b *blockingEnumerator) Synthesize(context.Context, types.SynthesisRequest, types.ProgressCallback) (*types.SynthesisResult, error) {
	return nil, types.NewAudioError(types.ErrCodeSynthesisFailed, "not used", nil)
}

func (b *blockingEnumerator) SynthesizeStream(context.Context, types.SynthesisRequest, chan<- types.AudioChunk) error {
	return nil
}

func (b *blockingEnumerator) GetVoices() []types.Voice { return nil }
func (b *blockingEnumerator) SupportsSSML() bool       { return false }
func (b *blockingEnumerator) Release() error           { return nil }

func (b *blockingEnumerator) VoicesContext(ctx context.Context) ([]types.Voice, error) {
	b.once.Do(func() { close(b.entered) })
	<-ctx.Done()
	b.cancelled.Store(true)
	close(b.finished)
	return nil, ctx.Err()
}

// TestTTSEnumerationIsCancellable proves enumeration runs under an execution
// lease and that cancellation reaches it and waits for cleanup.
func TestTTSEnumerationIsCancellable(t *testing.T) {
	p := NewPipeline(process.NewRegistry())
	enumerator := &blockingEnumerator{entered: make(chan struct{}), finished: make(chan struct{})}
	key := ttsEngineKey("kokoro", "")
	p.mu.Lock()
	p.ttsSynth = enumerator
	p.activeTTSEngine = "kokoro"
	p.activeTTSKey = key
	p.ttsProviders[key] = &ttsProvider{key: key, synth: enumerator, engine: "kokoro"}
	p.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := p.GetTTSVoicesContext(ctx)
		result <- err
	}()

	select {
	case <-enumerator.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("enumeration did not start")
	}
	p.mu.RLock()
	refs := p.ttsProviders[key].refs
	p.mu.RUnlock()
	if refs != 1 {
		t.Fatalf("enumeration execution leases = %d, want 1", refs)
	}

	cancel()
	select {
	case err := <-result:
		if err == nil {
			t.Fatal("cancelled enumeration must report the cancellation")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cancelled enumeration did not return")
	}
	select {
	case <-enumerator.finished:
	default:
		t.Fatal("enumeration returned before its cleanup finished")
	}
	p.mu.RLock()
	refs = p.ttsProviders[key].refs
	p.mu.RUnlock()
	if refs != 0 {
		t.Fatalf("execution lease leaked: %d", refs)
	}
}

// TestAudioChildHelper is the real child used by the shutdown test.
func TestAudioChildHelper(t *testing.T) {
	if os.Getenv("HOWL_AUDIO_CHILD") == "" {
		t.Skip("helper process only")
	}
	time.Sleep(30 * time.Second)
}

// TestPipelineShutdownStopsOwnedTTSChild proves shutdown escalates a live owned
// child to confirmed exit instead of returning with the process still running.
func TestPipelineShutdownStopsOwnedTTSChild(t *testing.T) {
	registry := process.NewRegistry()
	p := NewPipeline(registry)
	child, err := process.StartChild(registry, process.ChildSpec{
		Role:     process.RoleTTSChild,
		Instance: "shutdown-test",
		Path:     os.Args[0],
		Args:     []string{"-test.run=TestAudioChildHelper"},
		Env:      append(os.Environ(), "HOWL_AUDIO_CHILD=1"),
		Capture:  true,
	})
	if err != nil {
		t.Fatalf("start owned child: %v", err)
	}

	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	if !child.Exited() {
		t.Fatal("shutdown returned before the owned child exited")
	}
	if got := registry.Children(process.RoleTTSChild); len(got) != 0 {
		t.Fatalf("shutdown left owned children registered: %+v", got)
	}
}

// TestPipelineShutdownDoesNotClaimSuccessWithExecutingProvider proves a provider
// that is still executing is reported as a failure and keeps its residency,
// instead of being silently unloaded and reported clean.
func TestPipelineShutdownDoesNotClaimSuccessWithExecutingProvider(t *testing.T) {
	p := NewPipeline(process.NewRegistry())
	p.admission = permissiveAuthority()
	p.ttsCache = cache.NewManager(10)
	provider := installResidentProvider(t, p, ttsEngineKey("piper", ""), "piper", &stubTTSSynth{}, true)
	release := p.acquireTTSExecution(provider)

	err := p.Shutdown(context.Background())
	if err == nil || !strings.Contains(err.Error(), "executing") {
		t.Fatalf("shutdown with an executing provider must be reported: %v", err)
	}
	if got := resource.GetManager().ModelsByType(resource.ResourceTTS); len(got) == 0 {
		t.Fatal("a retained provider lost its residency registration")
	}

	release()
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown after the provider finished: %v", err)
	}
	if got := resource.GetManager().ModelsByType(resource.ResourceTTS); len(got) != 0 {
		t.Fatalf("registrations left after clean shutdown: %v", got)
	}
}

// TestAudioCrossWorkloadCapacityAndOwnership proves audio residency participates
// in the same capacity accounting as chat: it appears as a blocker while
// resident, and disappears from the decision when it is released.
func TestAudioCrossWorkloadCapacityAndOwnership(t *testing.T) {
	const mb = int64(1024 * 1024)
	authority := resource.NewAuthority(
		func() (resource.HardwareSnapshot, error) {
			return resource.HardwareSnapshot{
				GPUModel:       "Test GPU",
				TotalVRAMBytes: 24000 * mb,
				FreeVRAMBytes:  3500 * mb,
				UsedVRAMBytes:  20500 * mb,
				TotalRAMBytes:  32 * 1024 * mb,
				FreeRAMBytes:   8000 * mb,
				VRAMLive:       true,
				VRAMSource:     "test-live",
				GPUAvailable:   true,
			}, nil
		},
		func() []resource.ModelState { return resource.GetManager().LoadedModels() },
	)

	resource.GetManager().RegisterModel(resource.ResourceLLM, "chat", "", 3000, 3000, 1, false)
	t.Cleanup(func() { resource.GetManager().UnregisterModel(resource.ResourceLLM, "chat") })

	p := NewPipeline(process.NewRegistry())
	p.admission = authority
	p.ttsCache = cache.NewManager(10)
	installResidentProvider(t, p, ttsEngineKey("piper", ""), "piper", &stubTTSSynth{}, true)

	req, err := asrAdmissionRequest(types.ASRQwen3ASR17BLocal)
	if err != nil {
		t.Fatalf("asr request: %v", err)
	}
	decision, lease := authority.Admit(req)
	if decision.Admitted {
		lease.Release()
		t.Fatal("ASR was admitted beside chat and a resident TTS provider on 3500 MB free VRAM")
	}
	if !blockingNames(decision.Blocking, "piper") {
		t.Fatalf("TTS residency is missing from the refusal: %s", decision.Reason)
	}
	if !blockingNames(decision.Blocking, "chat") {
		t.Fatalf("chat residency is missing from the refusal: %s", decision.Reason)
	}

	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}

	decision, lease = authority.Admit(req)
	if decision.Admitted {
		lease.Release()
		t.Fatal("ASR still refused after the TTS provider was released")
	}
	if blockingNames(decision.Blocking, "piper") {
		t.Fatalf("released TTS residency still appears in capacity accounting: %s", decision.Reason)
	}
	if !blockingNames(decision.Blocking, "chat") {
		t.Fatalf("chat must stay the remaining blocker: %s", decision.Reason)
	}
}

func blockingNames(blocking []resource.ModelState, name string) bool {
	for _, model := range blocking {
		if model.Name == name {
			return true
		}
	}
	return false
}

// TestMelottsResidencyIsNotDoubleCounted proves a service this process already
// owns is not charged a second admission reservation.
func TestMelottsResidencyIsNotDoubleCounted(t *testing.T) {
	registry := process.NewRegistry()
	port := freeLoopbackPort(t)
	url := "http://127.0.0.1:" + port
	t.Setenv("HOWL_MELOTTS_COMMAND", os.Args[0])
	t.Setenv("HOWL_MELOTTS_ARGS", "-test.run=TestAudioMeloServiceHelper")
	t.Setenv("HOWL_MELO_HELPER", "serve")
	t.Setenv("HOWL_MELO_HELPER_PORT", port)

	service, err := melotts.AcquireService(registry, url)
	if err != nil {
		t.Skipf("managed MeloTTS helper could not start in this environment: %v", err)
	}
	t.Cleanup(func() { _ = melotts.ReleaseService(service) })

	p := NewPipeline(registry)
	p.admission = permissiveAuthority()
	if _, lease, err := p.ttsResidencyRequest("melotts", url); err != nil {
		t.Fatalf("residency request: %v", err)
	} else if lease != nil {
		lease.Release()
		t.Fatal("an already-owned MeloTTS service took a second admission reservation")
	}
}

// TestAudioMeloServiceHelper is the local service Howl can own in tests: it
// listens on the configured loopback port until terminated.
func TestAudioMeloServiceHelper(t *testing.T) {
	if os.Getenv("HOWL_MELO_HELPER") == "" {
		t.Skip("helper process only")
	}
	port := os.Getenv("HOWL_MELO_HELPER_PORT")
	if port == "" {
		t.Fatal("helper requires HOWL_MELO_HELPER_PORT")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	time.Sleep(30 * time.Second)
}

func freeLoopbackPort(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve port: %v", err)
	}
	defer listener.Close()
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener address %T", listener.Addr())
	}
	return strconv.Itoa(addr.Port)
}
