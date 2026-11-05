package audio

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	audioruntime "howl-chat/internal/audio/runtime"
	"howl-chat/internal/audio/types"
)

type scriptedASR struct {
	initErr      error
	initPanic    any
	releaseErr   error
	transcribe   string
	config       types.ASRConfig
	initCalls    atomic.Int32
	releaseCalls atomic.Int32
	initStarted  chan struct{}
	initRelease  <-chan struct{}
	startOnce    *sync.Once
}

func (s *scriptedASR) Initialize(ctx context.Context, config types.ASRConfig) error {
	s.config = config
	s.initCalls.Add(1)
	if s.initStarted != nil {
		if s.startOnce != nil {
			s.startOnce.Do(func() { close(s.initStarted) })
		} else {
			close(s.initStarted)
		}
	}
	if s.initRelease != nil {
		<-s.initRelease
	}
	if s.initPanic != nil {
		panic(s.initPanic)
	}
	return s.initErr
}

func (s *scriptedASR) Transcribe(context.Context, string, types.ProgressCallback) (*types.RecognitionResult, error) {
	text := s.transcribe
	if text == "" {
		text = "transcribed"
	}
	return &types.RecognitionResult{Text: text}, nil
}

func (s *scriptedASR) TranscribeStream(context.Context, <-chan types.AudioChunk, chan<- types.RecognitionResult) error {
	return nil
}

func (s *scriptedASR) SupportsLanguage(string) bool { return true }
func (s *scriptedASR) GetModelInfo() types.ASRModelInfo {
	return types.ASRModelInfo{Name: "scripted ASR"}
}
func (s *scriptedASR) Release() error {
	s.releaseCalls.Add(1)
	return s.releaseErr
}

type scriptedTTS struct {
	initErr      error
	initPanic    any
	releaseErr   error
	releasePanic any
	config       types.TTSConfig
	initCalls    atomic.Int32
	releaseCalls atomic.Int32
	synthCalls   atomic.Int32
	initStarted  chan struct{}
	initRelease  <-chan struct{}
	startOnce    *sync.Once
}

func (s *scriptedTTS) Initialize(ctx context.Context, config types.TTSConfig) error {
	s.config = config
	s.initCalls.Add(1)
	if s.initStarted != nil {
		if s.startOnce != nil {
			s.startOnce.Do(func() { close(s.initStarted) })
		} else {
			close(s.initStarted)
		}
	}
	if s.initRelease != nil {
		<-s.initRelease
	}
	if s.initPanic != nil {
		panic(s.initPanic)
	}
	return s.initErr
}

func (s *scriptedTTS) Synthesize(context.Context, types.SynthesisRequest, types.ProgressCallback) (*types.SynthesisResult, error) {
	s.synthCalls.Add(1)
	return &types.SynthesisResult{AudioPath: "scripted.wav"}, nil
}

func (s *scriptedTTS) SynthesizeStream(context.Context, types.SynthesisRequest, chan<- types.AudioChunk) error {
	return nil
}

func (s *scriptedTTS) GetVoices() []types.Voice { return []types.Voice{{ID: "scripted"}} }
func (s *scriptedTTS) SupportsSSML() bool       { return false }
func (s *scriptedTTS) Release() error {
	s.releaseCalls.Add(1)
	if s.releasePanic != nil {
		panic(s.releasePanic)
	}
	return s.releaseErr
}

type providerBehavior struct {
	asrInitErr      error
	asrInitPanic    any
	ttsInitErr      error
	ttsInitPanic    any
	ttsReleaseErr   error
	ttsReleasePanic any
}

type providerSet struct {
	asr map[types.ASRModelType]*scriptedASR
	tts map[types.TTSModelType]*scriptedTTS
}

func newProviderSet(behavior providerBehavior) *providerSet {
	set := &providerSet{
		asr: make(map[types.ASRModelType]*scriptedASR, len(types.ASRModelRegistry)),
		tts: make(map[types.TTSModelType]*scriptedTTS, len(types.ModelRegistry.TTSModels)),
	}
	for model := range types.ASRModelRegistry {
		set.asr[model] = &scriptedASR{
			initErr:   behavior.asrInitErr,
			initPanic: behavior.asrInitPanic,
			startOnce: &sync.Once{},
		}
	}
	for model := range types.ModelRegistry.TTSModels {
		set.tts[model] = &scriptedTTS{
			initErr:      behavior.ttsInitErr,
			initPanic:    behavior.ttsInitPanic,
			releaseErr:   behavior.ttsReleaseErr,
			releasePanic: behavior.ttsReleasePanic,
			startOnce:    &sync.Once{},
		}
	}
	return set
}

func (s *providerSet) install(p *Pipeline) {
	for model, provider := range s.asr {
		model := model
		provider := provider
		p.loader.RegisterASR(model, func() types.ASRRecognizer { return provider })
	}
	for model, provider := range s.tts {
		model := model
		provider := provider
		p.loader.RegisterTTS(model, func() types.TTSSynthesizer { return provider })
	}
}

func (s *providerSet) asrReleases() int {
	total := 0
	for _, provider := range s.asr {
		total += int(provider.releaseCalls.Load())
	}
	return total
}

func (s *providerSet) ttsReleases() int {
	total := 0
	for _, provider := range s.tts {
		total += int(provider.releaseCalls.Load())
	}
	return total
}

func (s *providerSet) activeASR() []*scriptedASR {
	var out []*scriptedASR
	for _, provider := range s.asr {
		if provider.initCalls.Load() > 0 && provider.initErr == nil && provider.initPanic == nil {
			out = append(out, provider)
		}
	}
	return out
}

func testPipelineConfig(asrEnabled, ttsEnabled bool) types.PipelineConfig {
	return types.PipelineConfig{
		ASREnabled:   asrEnabled,
		TTSEnabled:   ttsEnabled,
		HardwareTier: "large",
	}
}

func newScriptedPipeline(behavior providerBehavior) (*Pipeline, *providerSet) {
	p := NewPipeline()
	set := newProviderSet(behavior)
	set.install(p)
	return p, set
}

func TestPipelineR1CapabilityOutcomes(t *testing.T) {
	asrErr := errors.New("asr unavailable")
	ttsErr := errors.New("tts unavailable")

	tests := []struct {
		name     string
		asr      bool
		tts      bool
		behavior providerBehavior
		wantASR  bool
		wantTTS  bool
		wantErrs []error
	}{
		{name: "ASR only", asr: true, wantASR: true},
		{name: "TTS only", tts: true, wantTTS: true},
		{name: "both success", asr: true, tts: true, wantASR: true, wantTTS: true},
		{name: "ASR failure with TTS success", asr: true, tts: true, behavior: providerBehavior{asrInitErr: asrErr}, wantTTS: true, wantErrs: []error{asrErr}},
		{name: "TTS failure with ASR success", asr: true, tts: true, behavior: providerBehavior{ttsInitErr: ttsErr}, wantASR: true, wantErrs: []error{ttsErr}},
		{name: "both fail", asr: true, tts: true, behavior: providerBehavior{asrInitErr: asrErr, ttsInitErr: ttsErr}, wantErrs: []error{asrErr, ttsErr}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p, _ := newScriptedPipeline(tc.behavior)
			err := p.Initialize(context.Background(), testPipelineConfig(tc.asr, tc.tts))
			outcome := p.initializationResult()
			status := p.GetStatus()

			if outcome.asr.requested != tc.asr || outcome.tts.requested != tc.tts {
				t.Fatalf("requested = ASR:%v TTS:%v, want ASR:%v TTS:%v", outcome.asr.requested, outcome.tts.requested, tc.asr, tc.tts)
			}
			if outcome.asr.ready != tc.wantASR || outcome.tts.ready != tc.wantTTS {
				t.Fatalf("ready = ASR:%v TTS:%v, want ASR:%v TTS:%v", outcome.asr.ready, outcome.tts.ready, tc.wantASR, tc.wantTTS)
			}
			if status.ASRReady != tc.wantASR || status.TTSReady != tc.wantTTS {
				t.Fatalf("status = ASR:%v TTS:%v, want ASR:%v TTS:%v", status.ASRReady, status.TTSReady, tc.wantASR, tc.wantTTS)
			}
			for _, want := range tc.wantErrs {
				if !errors.Is(err, want) {
					t.Fatalf("Initialize error = %v, want %v", err, want)
				}
			}
			if len(tc.wantErrs) == 0 && err != nil {
				t.Fatalf("Initialize error = %v, want nil", err)
			}
			if err := p.Shutdown(context.Background()); err != nil {
				t.Fatalf("Shutdown: %v", err)
			}
		})
	}
}

func TestPipelineR6NoImplicitTTS(t *testing.T) {
	p, set := newScriptedPipeline(providerBehavior{})
	if err := p.Initialize(context.Background(), testPipelineConfig(true, true)); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	resp, err := p.ProcessAudioMessage(context.Background(), types.AudioMessage{Type: "audio_input", AudioPath: "input.wav"})
	if err != nil {
		t.Fatalf("ProcessAudioMessage: %v", err)
	}
	if resp.TranscribedText != "transcribed" || resp.AudioOutputPath != "" {
		t.Fatalf("response = %+v, want transcription without audio", resp)
	}
	for _, synth := range set.tts {
		if synth.synthCalls.Load() != 0 {
			t.Fatalf("implicit TTS synthesis calls = %d, want 0", synth.synthCalls.Load())
		}
	}
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestPipelineR7ChainOwnership(t *testing.T) {
	p, set := newScriptedPipeline(providerBehavior{})
	if err := p.Initialize(context.Background(), testPipelineConfig(true, false)); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	for model, provider := range set.asr {
		if provider.initCalls.Load() != 1 {
			continue
		}
		wantPath := ASRModelPathFor(model)
		if provider.config.ModelPath != wantPath {
			t.Fatalf("ASR %s config path = %q, want %q", model, provider.config.ModelPath, wantPath)
		}
	}
	if _, err := p.ProcessAudioMessage(context.Background(), types.AudioMessage{Type: "audio_input", AudioPath: "input.wav"}); err != nil {
		t.Fatalf("ProcessAudioMessage: %v", err)
	}
	for model, provider := range set.asr {
		if provider.initCalls.Load() > 0 && provider.initCalls.Load() != 1 {
			t.Fatalf("ASR %s initialized %d times, want 1", model, provider.initCalls.Load())
		}
		if provider.releaseCalls.Load() != 0 {
			t.Fatalf("ASR %s released before shutdown", model)
		}
	}
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	for model, provider := range set.asr {
		if provider.initCalls.Load() > 0 && provider.releaseCalls.Load() != 1 {
			t.Fatalf("ASR %s release calls = %d, want 1", model, provider.releaseCalls.Load())
		}
	}
}

func TestPipelineR8LazyTTSOwnership(t *testing.T) {
	p, set := newScriptedPipeline(providerBehavior{})
	if err := p.Initialize(context.Background(), testPipelineConfig(false, true)); err != nil {
		t.Fatalf("Initialize: %v", err)
	}

	dir := t.TempDir()
	voices := []string{"test-voice-one", "test-voice-two"}
	paths := make([]string, len(voices))
	for i, voice := range voices {
		paths[i] = filepath.Join(dir, voice+".onnx")
		if err := os.WriteFile(paths[i], []byte("model"), 0o600); err != nil {
			t.Fatalf("write model: %v", err)
		}
		t.Setenv("HOWL_TTS_PIPER_MODEL", paths[i])
		if err := p.ensureTTSVoice(context.Background(), voice); err != nil {
			t.Fatalf("ensureTTSVoice(%s): %v", voice, err)
		}
	}

	ids := ttsIDs(p.loader.GetLoadedModels())
	if len(ids) != 3 {
		t.Fatalf("tracked TTS ids = %v, want 3", ids)
	}
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("duplicate TTS id %q", id)
		}
		seen[id] = true
	}

	failedVoice := "test-voice-three"
	failedPath := filepath.Join(dir, failedVoice+".onnx")
	if err := os.WriteFile(failedPath, []byte("model"), 0o600); err != nil {
		t.Fatalf("write failed model: %v", err)
	}
	t.Setenv("HOWL_TTS_PIPER_MODEL", failedPath)
	loadErr := errors.New("replacement failed")
	set.tts[types.TTSPiper].initErr = loadErr
	if err := p.ensureTTSVoice(context.Background(), failedVoice); err == nil {
		t.Fatal("ensureTTSVoice succeeded, want failure")
	}
	if got := len(ttsIDs(p.loader.GetLoadedModels())); got != 3 {
		t.Fatalf("tracked TTS ids after failed replacement = %d, want 3", got)
	}
	if p.loadedVoice != voices[1] {
		t.Fatalf("active voice = %q, want %q", p.loadedVoice, voices[1])
	}

	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	if got := len(ttsIDs(p.loader.GetLoadedModels())); got != 0 {
		t.Fatalf("tracked TTS ids after shutdown = %d, want 0", got)
	}
	if set.tts[types.TTSPiper].releaseCalls.Load() != 3 {
		t.Fatalf("piper release calls = %d, want 3", set.tts[types.TTSPiper].releaseCalls.Load())
	}
}

func ttsIDs(ids []string) []string {
	var out []string
	for _, id := range ids {
		if strings.HasPrefix(id, "tts:") {
			out = append(out, id)
		}
	}
	return out
}

func scriptedManager(t *testing.T, builds []providerBehavior) (*Manager, *[]*providerSet) {
	t.Helper()
	sets := make([]*providerSet, 0, len(builds))
	index := 0
	manager := NewManager(func(ctx context.Context, desired Settings, candidate *Candidate) error {
		if index >= len(builds) {
			t.Fatalf("unexpected factory call %d", index+1)
		}
		behavior := builds[index]
		index++
		p := NewPipeline()
		set := newProviderSet(behavior)
		set.install(p)
		sets = append(sets, set)
		candidate.Pipeline = p
		_ = p.Initialize(ctx, testPipelineConfig(desired.ASREnabled, desired.TTSEnabled))
		outcome := p.initializationResult()
		if desired.ASREnabled {
			candidate.ASR = CapabilityResult{Ready: outcome.asr.ready, Err: outcome.asr.err}
		}
		if desired.TTSEnabled {
			candidate.TTS = CapabilityResult{Ready: outcome.tts.ready, Err: outcome.tts.err}
		}
		return nil
	})
	return manager, &sets
}

func TestPipelineManagerR2StrictReplacement(t *testing.T) {
	asrErr := errors.New("candidate ASR failed")
	manager, sets := scriptedManager(t, []providerBehavior{{}, {asrInitErr: asrErr}})
	mustApply(t, manager, settingsBoth)

	status, _, err := manager.Apply(context.Background(), settingsBoth)
	if !errors.Is(err, asrErr) {
		t.Fatalf("replacement error = %v, want ASR error", err)
	}
	if status.Phase != "failed" || !status.ASR.Ready || !status.TTS.Ready {
		t.Fatalf("status = %+v, want failed with old generation ready", status)
	}
	if (*sets)[1].ttsReleases() != 1 {
		t.Fatalf("candidate TTS releases = %d, want 1", (*sets)[1].ttsReleases())
	}
	if (*sets)[0].ttsReleases() != 0 {
		t.Fatalf("old TTS released before close: %d", (*sets)[0].ttsReleases())
	}
	for model, provider := range (*sets)[1].asr {
		if provider.initCalls.Load() > 0 && provider.releaseCalls.Load() != 1 {
			t.Fatalf("failed ASR %s release calls = %d, want 1", model, provider.releaseCalls.Load())
		}
	}
	if status.Builds != 0 {
		t.Fatalf("builds = %d, want 0", status.Builds)
	}

	lease := mustAcquire(t, manager)
	if lease.Pipeline() == nil {
		t.Fatal("old generation is not acquirable")
	}
	lease.Release()
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if (*sets)[0].ttsReleases() != 1 {
		t.Fatalf("old TTS releases after close = %d, want 1", (*sets)[0].ttsReleases())
	}
}

func TestPipelineManagerR3DegradedStartup(t *testing.T) {
	asrErr := errors.New("ASR unavailable")
	manager, sets := scriptedManager(t, []providerBehavior{{asrInitErr: asrErr}})
	status, _, err := manager.Apply(context.Background(), settingsBoth)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if status.Phase != "degraded" {
		t.Fatalf("phase = %q, want degraded", status.Phase)
	}
	if status.ASR.Ready || status.ASR.Error == "" || !status.TTS.Ready {
		t.Fatalf("status = %+v, want ASR failed and TTS ready", status)
	}
	lease := mustAcquire(t, manager)
	if lease.Pipeline() == nil {
		t.Fatal("degraded generation is not acquirable")
	}
	lease.Release()
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if (*sets)[0].ttsReleases() != 1 {
		t.Fatalf("TTS releases = %d, want 1", (*sets)[0].ttsReleases())
	}
}

func TestPipelineManagerR4OrdinaryFailure(t *testing.T) {
	asrErr := errors.New("ASR failed")
	ttsErr := errors.New("TTS failed")
	manager, sets := scriptedManager(t, []providerBehavior{{asrInitErr: asrErr, ttsInitErr: ttsErr}})
	_, _, err := manager.Apply(context.Background(), settingsBoth)
	if !errors.Is(err, asrErr) || !errors.Is(err, ttsErr) {
		t.Fatalf("Apply error = %v, want both capability errors", err)
	}
	if status := manager.Status(); status.Builds != 0 || status.Phase != "failed" {
		t.Fatalf("status = %+v, want failed with zero builds", status)
	}
	for model, provider := range (*sets)[0].asr {
		if provider.initCalls.Load() > 0 && provider.releaseCalls.Load() != 1 {
			t.Fatalf("ASR %s release calls = %d, want 1", model, provider.releaseCalls.Load())
		}
	}
	for model, provider := range (*sets)[0].tts {
		if provider.releaseCalls.Load() != 1 {
			t.Fatalf("TTS %s release calls = %d, want 1", model, provider.releaseCalls.Load())
		}
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestPipelineManagerR5PanicCleanup(t *testing.T) {
	initPanic := errors.New("TTS initialization panic")
	cleanupErr := errors.New("TTS release failed")
	manager, sets := scriptedManager(t, []providerBehavior{{ttsInitPanic: initPanic, ttsReleaseErr: cleanupErr}})
	_, _, err := manager.Apply(context.Background(), settingsBoth)
	var failure *BuildFailure
	if !errors.As(err, &failure) {
		t.Fatalf("Apply error = %T %v, want *BuildFailure", err, err)
	}
	record, ok := failure.Recovered.(*audioruntime.ProviderInitPanic)
	if !ok {
		t.Fatalf("Recovered = %T %v, want *runtime.ProviderInitPanic", failure.Recovered, failure.Recovered)
	}
	if record.Value != initPanic {
		t.Fatalf("Value = %v, want original init panic", record.Value)
	}
	if !errors.Is(record, initPanic) || !errors.Is(record, cleanupErr) {
		t.Fatalf("wrapper does not preserve both errors: %v", record)
	}
	if record.CleanupErr != cleanupErr {
		t.Fatalf("CleanupErr = %v, want %v", record.CleanupErr, cleanupErr)
	}
	if !strings.HasPrefix(record.ProviderID, "tts:") {
		t.Fatalf("ProviderID = %q, want tts: prefix", record.ProviderID)
	}
	if failure.CleanupErr != nil {
		t.Fatalf("manager cleanup error = %v, want nil", failure.CleanupErr)
	}
	if !strings.Contains(failure.Error(), initPanic.Error()) || !strings.Contains(failure.Error(), cleanupErr.Error()) {
		t.Fatalf("failure text does not expose panic and cleanup: %v", failure.Error())
	}
	if (*sets)[0].ttsReleases() != 1 {
		t.Fatalf("panicking TTS releases = %d, want 1", (*sets)[0].ttsReleases())
	}
	for model, provider := range (*sets)[0].asr {
		if provider.initCalls.Load() > 0 && provider.releaseCalls.Load() != 1 {
			t.Fatalf("ASR %s release calls = %d, want 1", model, provider.releaseCalls.Load())
		}
	}
	if status := manager.Status(); status.Builds != 0 {
		t.Fatalf("builds = %d, want 0", status.Builds)
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestPipelineManagerR5CleanupPanic(t *testing.T) {
	initPanic := errors.New("TTS initialization panic")
	releasePanic := errors.New("TTS release panic")
	manager, _ := scriptedManager(t, []providerBehavior{{ttsInitPanic: initPanic, ttsReleasePanic: releasePanic}})
	_, _, err := manager.Apply(context.Background(), settingsBoth)
	var failure *BuildFailure
	if !errors.As(err, &failure) {
		t.Fatalf("Apply error = %T %v, want *BuildFailure", err, err)
	}
	record, ok := failure.Recovered.(*audioruntime.ProviderInitPanic)
	if !ok {
		t.Fatalf("Recovered = %T %v, want *runtime.ProviderInitPanic", failure.Recovered, failure.Recovered)
	}
	if record.Value != initPanic {
		t.Fatalf("Value = %v, want original init panic", record.Value)
	}
	if record.CleanupPanic != releasePanic {
		t.Fatalf("CleanupPanic = %v, want release panic", record.CleanupPanic)
	}
	if !strings.Contains(failure.Error(), releasePanic.Error()) {
		t.Fatalf("failure text does not expose cleanup panic: %v", failure.Error())
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestPipelineLifecycleShutdownWaitsForInitialize(t *testing.T) {
	p, set := newScriptedPipeline(providerBehavior{})
	started := make(chan struct{})
	release := make(chan struct{})
	sharedStart := &sync.Once{}
	for _, provider := range set.asr {
		provider.initStarted = started
		provider.initRelease = release
		provider.startOnce = sharedStart
	}

	initDone := make(chan error, 1)
	go func() {
		initDone <- p.Initialize(context.Background(), testPipelineConfig(true, false))
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("initialization did not start")
	}

	shutdownDone := make(chan error, 1)
	go func() { shutdownDone <- p.Shutdown(context.Background()) }()
	select {
	case <-shutdownDone:
		t.Fatal("Shutdown completed while Initialize was still running")
	case <-time.After(50 * time.Millisecond):
	}
	close(release)
	if err := <-initDone; err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	if err := <-shutdownDone; err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestPipelineConcurrentLazyTTSReplacement(t *testing.T) {
	p, set := newScriptedPipeline(providerBehavior{})
	if err := p.Initialize(context.Background(), testPipelineConfig(false, true)); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	voice := "concurrent-voice"
	dir := t.TempDir()
	path := filepath.Join(dir, voice+".onnx")
	if err := os.WriteFile(path, []byte("model"), 0o600); err != nil {
		t.Fatalf("write model: %v", err)
	}
	t.Setenv("HOWL_TTS_PIPER_MODEL", path)
	piper := set.tts[types.TTSPiper]
	before := piper.initCalls.Load()
	started := make(chan struct{})
	release := make(chan struct{})
	piper.initStarted = started
	piper.initRelease = release

	first := make(chan error, 1)
	go func() { first <- p.ensureTTSVoice(context.Background(), voice) }()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("replacement load did not start")
	}
	second := make(chan error, 1)
	go func() { second <- p.ensureTTSVoice(context.Background(), voice) }()
	close(release)
	if err := <-first; err != nil {
		t.Fatalf("first replacement: %v", err)
	}
	if err := <-second; err != nil {
		t.Fatalf("second replacement: %v", err)
	}
	if piper.initCalls.Load() != before+1 {
		t.Fatalf("piper init calls = %d, want %d", piper.initCalls.Load(), before+1)
	}
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestPipelineScriptedManagerUsesRealPipeline(t *testing.T) {
	manager, sets := scriptedManager(t, []providerBehavior{{}})
	if _, _, err := manager.Apply(context.Background(), settingsBoth); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(*sets) != 1 || (*sets)[0].asrReleases() != 0 {
		t.Fatal("real pipeline was not retained by the manager")
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
}
