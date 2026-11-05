package audio

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

var (
	settingsNone = Settings{}
	settingsASR  = Settings{ASREnabled: true}
	settingsTTS  = Settings{TTSEnabled: true}
	settingsBoth = Settings{ASREnabled: true, TTSEnabled: true}
)

type fakePipeline struct {
	name string

	shutdownStarted chan struct{}
	startedOnce     sync.Once
	shutdownRelease <-chan struct{}
	shutdownErr     error
	shutdownPanic   any
	shutdownCalls   atomic.Int32

	interrupts  []chan struct{}
	interruptMu sync.Mutex
}

func newFakePipeline(name string) *fakePipeline {
	return &fakePipeline{
		name:            name,
		shutdownStarted: make(chan struct{}),
	}
}

func (f *fakePipeline) ProcessAudioMessage(context.Context, types.AudioMessage) (*types.AudioResponse, error) {
	return &types.AudioResponse{}, nil
}

func (f *fakePipeline) GetStatus() types.PipelineStatus {
	return types.PipelineStatus{}
}

func (f *fakePipeline) GetTTSVoices() []types.Voice {
	return nil
}

func (f *fakePipeline) Shutdown(context.Context) error {
	f.shutdownCalls.Add(1)
	f.startedOnce.Do(func() {
		close(f.shutdownStarted)
	})
	if f.shutdownRelease != nil {
		<-f.shutdownRelease
	}
	if f.shutdownPanic != nil {
		panic(f.shutdownPanic)
	}
	return f.shutdownErr
}

func (f *fakePipeline) Interrupt() {
	f.interruptMu.Lock()
	for _, ch := range f.interrupts {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
	f.interruptMu.Unlock()
}

func (f *fakePipeline) addInterruptWatcher() <-chan struct{} {
	ch := make(chan struct{}, 4)
	f.interruptMu.Lock()
	f.interrupts = append(f.interrupts, ch)
	f.interruptMu.Unlock()
	return ch
}

func (f *fakePipeline) calls() int {
	return int(f.shutdownCalls.Load())
}

func candidateWith(p *fakePipeline, desired Settings, asrReady, ttsReady bool, asrErr, ttsErr error) *Candidate {
	candidate := &Candidate{Pipeline: p}
	if desired.ASREnabled {
		candidate.ASR = CapabilityResult{Ready: asrReady, Err: asrErr}
	}
	if desired.TTSEnabled {
		candidate.TTS = CapabilityResult{Ready: ttsReady, Err: ttsErr}
	}
	return candidate
}

func mustApply(t *testing.T, manager *Manager, desired Settings) (Status, bool) {
	t.Helper()
	status, superseded, err := manager.Apply(context.Background(), desired)
	if err != nil {
		t.Fatalf("Apply(%+v): %v", desired, err)
	}
	return status, superseded
}

func mustAcquire(t *testing.T, manager *Manager) *Lease {
	t.Helper()
	lease, err := manager.Acquire(context.Background())
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	return lease
}

func mustClose(t *testing.T, manager *Manager) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return manager.Close(ctx)
}

func waitClosed(t *testing.T, ch <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(3 * time.Second):
		t.Fatalf("timed out waiting for %s", what)
	}
}

func waitCalls(t *testing.T, p *fakePipeline, want int) {
	t.Helper()
	waitUntil(t, "pipeline shutdown count", func() bool {
		return p.calls() == want
	})
}

func waitUntil(t *testing.T, what string, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
}

func assertNoResult(t *testing.T, ch <-chan error, what string) {
	t.Helper()
	select {
	case err := <-ch:
		t.Fatalf("%s completed early with %v", what, err)
	default:
	}
}

// closeWaitContext signals when Close evaluates ctx.Done(), which happens only
// after it has entered the drain-wait select. Tests still need a blocked
// shutdown or operation to keep that select from completing immediately.
type closeWaitContext struct {
	context.Context
	waitEntered chan struct{}
	once        sync.Once
}

func newCloseWaitContext(timeout time.Duration) (*closeWaitContext, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &closeWaitContext{Context: ctx, waitEntered: make(chan struct{})}, cancel
}

func (c *closeWaitContext) Done() <-chan struct{} {
	c.once.Do(func() {
		close(c.waitEntered)
	})
	return c.Context.Done()
}

func closeWithTimeout(manager *Manager, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return manager.Close(ctx)
}

func releaseOnce(ch chan struct{}) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			close(ch)
		})
	}
}

func TestManagerM1GenerationLeasesAndOutOfOrderReleases(t *testing.T) {
	p1 := newFakePipeline("p1")
	p2 := newFakePipeline("p2")
	p3 := newFakePipeline("p3")
	built := []*fakePipeline{p1, p2, p3}
	next := 0
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		pipeline := built[next]
		next++
		*candidate = *candidateWith(pipeline, desired, true, true, nil, nil)
		return nil
	})

	mustApply(t, manager, settingsASR)
	lease1 := mustAcquire(t, manager)
	if lease1.Pipeline() != p1 {
		t.Fatalf("first lease pipeline = %T, want p1", lease1.Pipeline())
	}

	mustApply(t, manager, settingsTTS)
	lease2 := mustAcquire(t, manager)
	if lease2.Pipeline() != p2 {
		t.Fatalf("second lease pipeline = %T, want p2", lease2.Pipeline())
	}

	mustApply(t, manager, settingsBoth)
	lease3 := mustAcquire(t, manager)
	if lease3.Pipeline() != p3 {
		t.Fatalf("third lease pipeline = %T, want p3", lease3.Pipeline())
	}

	lease2.Release()
	lease2.Release()
	waitClosed(t, p2.shutdownStarted, "p2 shutdown after out-of-order release")
	if p2.calls() != 1 {
		t.Fatalf("p2 shutdown calls = %d, want 1", p2.calls())
	}
	if p1.calls() != 0 {
		t.Fatalf("p1 shutdown calls = %d before its release, want 0", p1.calls())
	}

	lease1.Release()
	waitClosed(t, p1.shutdownStarted, "p1 shutdown")
	if p1.calls() != 1 {
		t.Fatalf("p1 shutdown calls = %d, want 1", p1.calls())
	}

	lease3.Release()
	if p3.calls() != 0 {
		t.Fatalf("active p3 shutdown calls = %d before Close, want 0", p3.calls())
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
	waitCalls(t, p3, 1)
	for _, pipeline := range built {
		if pipeline.calls() != 1 {
			t.Fatalf("%s shutdown calls = %d, want exactly 1", pipeline.name, pipeline.calls())
		}
	}
}

func TestManagerM2StrictReplacementKeepsOldGeneration(t *testing.T) {
	asrErr := errors.New("asr failed")
	p1 := newFakePipeline("p1")
	p2 := newFakePipeline("p2")
	call := 0
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		call++
		switch call {
		case 1:
			*candidate = *candidateWith(p1, desired, true, true, nil, nil)
		case 2:
			*candidate = *candidateWith(p2, desired, false, true, asrErr, nil)
		default:
			t.Fatalf("unexpected factory call %d", call)
		}
		return nil
	})

	if _, _, err := manager.Apply(context.Background(), settingsBoth); err != nil {
		t.Fatalf("initial Apply: %v", err)
	}
	status, superseded, err := manager.Apply(context.Background(), settingsBoth)
	if err == nil {
		t.Fatal("replacement Apply returned nil error, want ASR rejection")
	}
	if superseded {
		t.Fatal("replacement Apply reported superseded, want current generation")
	}
	if !errors.Is(err, asrErr) {
		t.Fatalf("Apply error = %v, want wrapped ASR error", err)
	}
	if !status.ASR.Requested || !status.ASR.Ready {
		t.Fatalf("ASR status = %+v, want requested and old ready", status.ASR)
	}
	if status.ASR.Error == "" {
		t.Fatal("ASR status error is empty, want candidate ASR error")
	}
	if !status.TTS.Requested || !status.TTS.Ready {
		t.Fatalf("TTS status = %+v, want requested and ready", status.TTS)
	}
	if status.Phase != "failed" {
		t.Fatalf("phase = %q, want failed", status.Phase)
	}
	if p1.calls() != 0 {
		t.Fatalf("old pipeline shutdown calls = %d, want 0", p1.calls())
	}
	if p2.calls() != 1 {
		t.Fatalf("candidate shutdown calls = %d, want 1", p2.calls())
	}

	lease := mustAcquire(t, manager)
	if lease.Pipeline() != p1 {
		t.Fatal("Acquire returned the rejected candidate instead of the old generation")
	}
	lease.Release()
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
	waitCalls(t, p1, 1)
}

func TestManagerM3DegradedInitialActivation(t *testing.T) {
	asrErr := errors.New("asr unavailable")
	p1 := newFakePipeline("p1")
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		*candidate = *candidateWith(p1, desired, false, true, asrErr, nil)
		return nil
	})

	status, superseded, err := manager.Apply(context.Background(), settingsBoth)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if superseded {
		t.Fatal("Apply reported superseded")
	}
	if status.Phase != "degraded" {
		t.Fatalf("phase = %q, want degraded", status.Phase)
	}
	if !status.ASR.Requested || status.ASR.Ready || status.ASR.Error == "" {
		t.Fatalf("ASR status = %+v, want requested/not-ready/error", status.ASR)
	}
	if !status.TTS.Requested || !status.TTS.Ready || status.TTS.Error != "" {
		t.Fatalf("TTS status = %+v, want requested/ready/no-error", status.TTS)
	}

	lease := mustAcquire(t, manager)
	if lease.Pipeline() != p1 {
		t.Fatal("degraded activation did not install the ready candidate")
	}
	lease.Release()
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
	waitCalls(t, p1, 1)
}

func TestManagerM3bZeroReadyInitialActivationDisposesCandidate(t *testing.T) {
	asrErr := errors.New("asr failed")
	ttsErr := errors.New("tts failed")
	releaseShutdown := make(chan struct{})
	release := releaseOnce(releaseShutdown)
	t.Cleanup(release)

	p1 := newFakePipeline("p1")
	p1.shutdownRelease = releaseShutdown
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		*candidate = *candidateWith(p1, desired, false, false, asrErr, ttsErr)
		return nil
	})

	type outcome struct {
		status     Status
		superseded bool
		err        error
	}
	result := make(chan outcome, 1)
	go func() {
		status, superseded, err := manager.Apply(context.Background(), settingsBoth)
		result <- outcome{status: status, superseded: superseded, err: err}
	}()

	waitClosed(t, p1.shutdownStarted, "zero-ready candidate shutdown")
	if builds := manager.Status().Builds; builds != 1 {
		t.Fatalf("builds while cleanup is blocked = %d, want 1", builds)
	}
	if _, err := manager.Acquire(context.Background()); !errors.Is(err, ErrNoActivePipeline) {
		t.Fatalf("Acquire before cleanup = %v, want ErrNoActivePipeline", err)
	}
	blockedStatus := manager.Status()
	if blockedStatus.Phase != "failed" {
		t.Fatalf("phase while cleanup is blocked = %q, want failed", blockedStatus.Phase)
	}

	release()
	var got outcome
	select {
	case got = <-result:
	case <-time.After(3 * time.Second):
		t.Fatal("Apply did not return after candidate cleanup")
	}
	if got.err == nil {
		t.Fatal("Apply returned nil error for zero-ready requested activation")
	}
	if !errors.Is(got.err, asrErr) || !errors.Is(got.err, ttsErr) {
		t.Fatalf("Apply error = %v, want both capability errors", got.err)
	}
	if got.status.Phase != "failed" {
		t.Fatalf("phase = %q, want failed", got.status.Phase)
	}
	if !got.status.ASR.Requested || got.status.ASR.Ready || got.status.ASR.Error == "" {
		t.Fatalf("ASR status = %+v, want requested/not-ready/error", got.status.ASR)
	}
	if !got.status.TTS.Requested || got.status.TTS.Ready || got.status.TTS.Error == "" {
		t.Fatalf("TTS status = %+v, want requested/not-ready/error", got.status.TTS)
	}
	if builds := manager.Status().Builds; builds != 0 {
		t.Fatalf("builds after cleanup = %d, want 0", builds)
	}
	if _, err := manager.Acquire(context.Background()); !errors.Is(err, ErrNoActivePipeline) {
		t.Fatalf("Acquire after cleanup = %v, want ErrNoActivePipeline", err)
	}
	if p1.calls() != 1 {
		t.Fatalf("candidate shutdown calls = %d, want exactly 1", p1.calls())
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestManagerM4ExplicitDisableWithActiveGeneration(t *testing.T) {
	p1 := newFakePipeline("p1")
	factoryCalls := 0
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		factoryCalls++
		*candidate = *candidateWith(p1, desired, true, true, nil, nil)
		return nil
	})

	mustApply(t, manager, settingsBoth)
	lease := mustAcquire(t, manager)
	status, superseded, err := manager.Apply(context.Background(), settingsNone)
	if err != nil {
		t.Fatalf("disable Apply: %v", err)
	}
	if superseded {
		t.Fatal("disable Apply reported superseded")
	}
	if status.ASR.Requested || status.TTS.Requested {
		t.Fatalf("status requested after disable = %+v", status)
	}
	if status.ASR.Ready || status.TTS.Ready {
		t.Fatalf("status ready after disable = %+v", status)
	}
	if status.Phase != "idle" {
		t.Fatalf("phase = %q, want idle", status.Phase)
	}
	if status.Builds != 0 {
		t.Fatalf("builds = %d, want 0", status.Builds)
	}
	if factoryCalls != 1 {
		t.Fatalf("factory calls = %d, want no build for disable", factoryCalls)
	}
	if p1.calls() != 0 {
		t.Fatalf("shutdown calls while leased = %d, want 0", p1.calls())
	}
	if _, err := manager.Acquire(context.Background()); !errors.Is(err, ErrNoActivePipeline) {
		t.Fatalf("Acquire after disable = %v, want ErrNoActivePipeline", err)
	}

	lease.Release()
	waitClosed(t, p1.shutdownStarted, "disabled pipeline shutdown")
	if p1.calls() != 1 {
		t.Fatalf("shutdown calls = %d, want 1", p1.calls())
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestManagerM4bExplicitDisableFromEmptyManager(t *testing.T) {
	factoryCalls := 0
	manager := NewManager(func(context.Context, Settings, *Candidate) error {
		factoryCalls++
		return nil
	})

	status, superseded, err := manager.Apply(context.Background(), settingsNone)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if superseded {
		t.Fatal("Apply reported superseded")
	}
	if status.ASR.Requested || status.TTS.Requested || status.ASR.Ready || status.TTS.Ready {
		t.Fatalf("status = %+v, want both capabilities unrequested and not ready", status)
	}
	if status.Phase != "idle" {
		t.Fatalf("phase = %q, want idle", status.Phase)
	}
	if status.Builds != 0 {
		t.Fatalf("builds = %d, want 0", status.Builds)
	}
	if factoryCalls != 0 {
		t.Fatalf("factory calls = %d, want 0", factoryCalls)
	}
	if _, err := manager.Acquire(context.Background()); !errors.Is(err, ErrNoActivePipeline) {
		t.Fatalf("Acquire = %v, want ErrNoActivePipeline", err)
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestManagerM5ZeroReferenceRetirementViaApply(t *testing.T) {
	releaseOld := make(chan struct{})
	release := releaseOnce(releaseOld)
	t.Cleanup(release)

	p1 := newFakePipeline("p1")
	p1.shutdownRelease = releaseOld
	p2 := newFakePipeline("p2")
	call := 0
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		call++
		if call == 1 {
			*candidate = *candidateWith(p1, desired, true, true, nil, nil)
		} else {
			*candidate = *candidateWith(p2, desired, true, true, nil, nil)
		}
		return nil
	})

	mustApply(t, manager, settingsBoth)
	mustApply(t, manager, settingsBoth)
	waitClosed(t, p1.shutdownStarted, "retired p1 shutdown start")
	if builds := manager.Status().Builds; builds != 0 {
		t.Fatalf("builds while old shutdown is blocked = %d, want 0", builds)
	}
	lease := mustAcquire(t, manager)
	if lease.Pipeline() != p2 {
		t.Fatal("new generation is not active while old shutdown is blocked")
	}
	lease.Release()
	release()
	waitCalls(t, p1, 1)
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
	waitCalls(t, p2, 1)
	if p1.calls() != 1 {
		t.Fatalf("retired p1 shutdown calls after Close = %d, want exactly 1", p1.calls())
	}
	if p2.calls() != 1 {
		t.Fatalf("replacement p2 shutdown calls after Close = %d, want exactly 1", p2.calls())
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("repeated Close: %v", err)
	}
	if p1.calls() != 1 {
		t.Fatalf("retired p1 shutdown calls after repeated Close = %d, want exactly 1", p1.calls())
	}
	if p2.calls() != 1 {
		t.Fatalf("replacement p2 shutdown calls after repeated Close = %d, want exactly 1", p2.calls())
	}
}

func TestManagerM6ZeroReferenceRetirementViaClose(t *testing.T) {
	releaseShutdown := make(chan struct{})
	release := releaseOnce(releaseShutdown)
	t.Cleanup(release)

	p1 := newFakePipeline("p1")
	p1.shutdownRelease = releaseShutdown
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		*candidate = *candidateWith(p1, desired, true, false, nil, nil)
		return nil
	})
	mustApply(t, manager, settingsASR)

	ctx, cancel := newCloseWaitContext(3 * time.Second)
	defer cancel()
	result := make(chan error, 1)
	go func() {
		result <- manager.Close(ctx)
	}()
	waitClosed(t, p1.shutdownStarted, "Close-initiated shutdown")
	waitClosed(t, ctx.waitEntered, "Close drain wait")
	assertNoResult(t, result, "Close before shutdown completion")

	release()
	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("Close: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Close did not complete after shutdown")
	}
	if status := manager.Status(); status.Phase != "closed" {
		t.Fatalf("phase = %q, want closed", status.Phase)
	}
	if p1.calls() != 1 {
		t.Fatalf("shutdown calls = %d, want 1", p1.calls())
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("repeated Close: %v", err)
	}
	if p1.calls() != 1 {
		t.Fatalf("shutdown calls after repeated Close = %d, want 1", p1.calls())
	}
}

func TestManagerM7CloseOverlapsFinalRelease(t *testing.T) {
	releaseShutdown := make(chan struct{})
	release := releaseOnce(releaseShutdown)
	t.Cleanup(release)

	p1 := newFakePipeline("p1")
	p1.shutdownRelease = releaseShutdown
	interrupted := p1.addInterruptWatcher()
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		*candidate = *candidateWith(p1, desired, true, false, nil, nil)
		return nil
	})
	mustApply(t, manager, settingsASR)
	lease := mustAcquire(t, manager)

	ctx, cancel := newCloseWaitContext(3 * time.Second)
	defer cancel()
	result := make(chan error, 1)
	go func() {
		result <- manager.Close(ctx)
	}()
	waitClosed(t, interrupted, "Close start")
	waitClosed(t, ctx.waitEntered, "Close drain wait with an outstanding lease")
	assertNoResult(t, result, "Close with an outstanding lease")
	if p1.calls() != 0 {
		t.Fatalf("shutdown calls while lease outstanding = %d, want 0", p1.calls())
	}

	lease.Release()
	waitClosed(t, p1.shutdownStarted, "final release shutdown start")
	assertNoResult(t, result, "Close before blocked shutdown returns")

	release()
	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("Close: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Close did not complete after final release")
	}
	if p1.calls() != 1 {
		t.Fatalf("shutdown calls = %d, want 1", p1.calls())
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("repeated Close: %v", err)
	}
	if p1.calls() != 1 {
		t.Fatalf("shutdown calls after repeated Close = %d, want exactly 1", p1.calls())
	}
}

func TestManagerM8TimeoutThenLaterDrain(t *testing.T) {
	releaseShutdown := make(chan struct{})
	release := releaseOnce(releaseShutdown)
	t.Cleanup(release)

	p1 := newFakePipeline("p1")
	p1.shutdownRelease = releaseShutdown
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		*candidate = *candidateWith(p1, desired, true, false, nil, nil)
		return nil
	})
	mustApply(t, manager, settingsASR)

	first := closeWithTimeout(manager, 50*time.Millisecond)
	if !errors.Is(first, ErrDrainTimeout) {
		t.Fatalf("first Close = %v, want ErrDrainTimeout", first)
	}
	waitClosed(t, p1.shutdownStarted, "first Close shutdown start")
	if status := manager.Status(); status.Phase != "closing" {
		t.Fatalf("phase after timeout = %q, want closing", status.Phase)
	}

	release()
	second := closeWithTimeout(manager, 3*time.Second)
	if second != nil {
		t.Fatalf("second Close = %v, want nil", second)
	}
	if status := manager.Status(); status.Phase != "closed" {
		t.Fatalf("phase = %q, want closed", status.Phase)
	}
	if p1.calls() != 1 {
		t.Fatalf("shutdown calls = %d, want 1", p1.calls())
	}
}

func TestManagerM9TimeoutThenDrainWithShutdownError(t *testing.T) {
	shutdownErr := errors.New("shutdown failed")
	releaseShutdown := make(chan struct{})
	release := releaseOnce(releaseShutdown)
	t.Cleanup(release)

	p1 := newFakePipeline("p1")
	p1.shutdownRelease = releaseShutdown
	p1.shutdownErr = shutdownErr
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		*candidate = *candidateWith(p1, desired, true, false, nil, nil)
		return nil
	})
	mustApply(t, manager, settingsASR)

	first := closeWithTimeout(manager, 50*time.Millisecond)
	if !errors.Is(first, ErrDrainTimeout) {
		t.Fatalf("first Close = %v, want ErrDrainTimeout", first)
	}
	release()
	second := closeWithTimeout(manager, 3*time.Second)
	if second == nil {
		t.Fatal("second Close returned nil, want shutdown error")
	}
	if !errors.Is(second, shutdownErr) {
		t.Fatalf("second Close = %v, want wrapped shutdown error", second)
	}
	if status := manager.Status(); status.Phase != "closed_with_errors" {
		t.Fatalf("phase = %q, want closed_with_errors", status.Phase)
	}
	if p1.calls() != 1 {
		t.Fatalf("shutdown calls = %d, want 1", p1.calls())
	}
}

func TestManagerM10DrainSignal(t *testing.T) {
	releaseShutdown := make(chan struct{})
	release := releaseOnce(releaseShutdown)
	t.Cleanup(release)

	p1 := newFakePipeline("p1")
	p1.shutdownRelease = releaseShutdown
	interrupted := p1.addInterruptWatcher()
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		*candidate = *candidateWith(p1, desired, true, false, nil, nil)
		return nil
	})
	mustApply(t, manager, settingsASR)

	first := make(chan error, 1)
	second := make(chan error, 1)
	go func() {
		first <- closeWithTimeout(manager, 3*time.Second)
	}()
	waitClosed(t, interrupted, "first Close start")
	go func() {
		second <- closeWithTimeout(manager, 3*time.Second)
	}()
	waitClosed(t, interrupted, "second Close start")

	manager.mu.Lock()
	drained := manager.drained
	manager.mu.Unlock()
	if drained == nil {
		t.Fatal("drain channel is nil while work is outstanding")
	}
	select {
	case <-drained:
		t.Fatal("drain channel closed before work completed")
	default:
	}

	release()
	for name, result := range map[string]<-chan error{"first": first, "second": second} {
		select {
		case err := <-result:
			if err != nil {
				t.Fatalf("%s Close: %v", name, err)
			}
		case <-time.After(3 * time.Second):
			t.Fatalf("%s Close did not complete", name)
		}
	}

	manager.mu.Lock()
	drainedAfter := manager.drained
	manager.mu.Unlock()
	if drainedAfter != nil {
		t.Fatalf("drain channel after finalization = %v, want nil", drainedAfter)
	}
	select {
	case <-drained:
	default:
		t.Fatal("captured drain channel was not closed")
	}
	if p1.calls() != 1 {
		t.Fatalf("shutdown calls = %d, want 1", p1.calls())
	}
}

func TestManagerOperationDrain(t *testing.T) {
	manager := NewManager(nil)
	op, ctx := manager.BeginOperation(context.Background(), 3*time.Second)
	if op == nil {
		t.Fatal("BeginOperation returned nil operation")
	}
	defer op.End()
	if status := manager.Status(); status.InFlight != 1 {
		t.Fatalf("in-flight count with one operation = %d, want 1", status.InFlight)
	}

	closeCtx, cancel := newCloseWaitContext(3 * time.Second)
	defer cancel()
	result := make(chan error, 1)
	go func() {
		result <- manager.Close(closeCtx)
	}()
	waitUntil(t, "Close to observe the registered operation", func() bool {
		return manager.Status().Phase == "closing"
	})
	waitClosed(t, closeCtx.waitEntered, "Close drain wait with a registered operation")
	select {
	case <-ctx.Done():
	case <-time.After(3 * time.Second):
		t.Fatal("Close did not cancel the registered operation context")
	}
	assertNoResult(t, result, "Close with a registered operation")

	op.End()
	op.End()
	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("Close: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Close did not complete after operation end")
	}
	if status := manager.Status(); status.InFlight != 0 {
		t.Fatalf("in-flight count after operation end = %d, want 0", status.InFlight)
	}
}

func TestManagerM11CoalescingAndStaleBuilds(t *testing.T) {
	p1 := newFakePipeline("p1")
	p2 := newFakePipeline("p2")
	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	call := 0
	manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
		call++
		switch call {
		case 1:
			close(firstStarted)
			<-releaseFirst
			*candidate = *candidateWith(p1, desired, true, false, nil, nil)
		case 2:
			*candidate = *candidateWith(p2, desired, false, true, nil, nil)
		default:
			t.Fatalf("unexpected factory call %d", call)
		}
		return nil
	})

	type result struct {
		status     Status
		superseded bool
		err        error
	}
	first := make(chan result, 1)
	go func() {
		status, superseded, err := manager.Apply(context.Background(), settingsASR)
		first <- result{status: status, superseded: superseded, err: err}
	}()
	waitClosed(t, firstStarted, "first build start")

	second := make(chan result, 1)
	go func() {
		status, superseded, err := manager.Apply(context.Background(), settingsTTS)
		second <- result{status: status, superseded: superseded, err: err}
	}()
	waitUntil(t, "second desired generation", func() bool {
		manager.mu.Lock()
		defer manager.mu.Unlock()
		return manager.desiredGen >= 2
	})
	close(releaseFirst)

	var firstResult result
	select {
	case firstResult = <-first:
	case <-time.After(3 * time.Second):
		t.Fatal("first Apply did not return")
	}
	var secondResult result
	select {
	case secondResult = <-second:
	case <-time.After(3 * time.Second):
		t.Fatal("second Apply did not return")
	}
	if firstResult.err != nil {
		t.Fatalf("first Apply: %v", firstResult.err)
	}
	if secondResult.err != nil {
		t.Fatalf("second Apply: %v", secondResult.err)
	}
	if !firstResult.superseded {
		t.Fatal("first Apply was not reported superseded")
	}
	if secondResult.superseded {
		t.Fatal("newest Apply was reported superseded")
	}
	if firstResult.status.TTS.Requested != true || firstResult.status.ASR.Requested != false {
		t.Fatalf("first status = %+v, want latest TTS-only desire", firstResult.status)
	}
	if call != 2 {
		t.Fatalf("factory calls = %d, want 2", call)
	}
	if p1.calls() != 1 {
		t.Fatalf("obsolete candidate shutdown calls = %d, want 1", p1.calls())
	}
	if p2.calls() != 0 {
		t.Fatalf("active candidate shutdown calls = %d, want 0", p2.calls())
	}
	if builds := manager.Status().Builds; builds != 0 {
		t.Fatalf("builds after coalescing = %d, want 0", builds)
	}
	lease := mustAcquire(t, manager)
	if lease.Pipeline() != p2 {
		t.Fatal("final active generation is not p2")
	}
	lease.Release()
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
	waitCalls(t, p2, 1)
	if p1.calls() != 1 {
		t.Fatalf("obsolete candidate shutdown calls after Close = %d, want exactly 1", p1.calls())
	}
	if p2.calls() != 1 {
		t.Fatalf("final candidate shutdown calls after Close = %d, want exactly 1", p2.calls())
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("repeated Close: %v", err)
	}
	if p1.calls() != 1 {
		t.Fatalf("obsolete candidate shutdown calls after repeated Close = %d, want exactly 1", p1.calls())
	}
	if p2.calls() != 1 {
		t.Fatalf("final candidate shutdown calls after repeated Close = %d, want exactly 1", p2.calls())
	}
}

func TestManagerM12CandidateAccounting(t *testing.T) {
	t.Run("accepted activation", func(t *testing.T) {
		p1 := newFakePipeline("p1")
		var m *Manager
		duringFactory := -1
		manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
			duringFactory = m.Status().Builds
			*candidate = *candidateWith(p1, desired, true, false, nil, nil)
			return nil
		})
		m = manager
		if _, _, err := manager.Apply(context.Background(), settingsASR); err != nil {
			t.Fatalf("Apply: %v", err)
		}
		if duringFactory != 1 {
			t.Fatalf("builds during factory = %d, want 1", duringFactory)
		}
		if builds := manager.Status().Builds; builds != 0 {
			t.Fatalf("builds after activation = %d, want 0", builds)
		}
		if err := mustClose(t, manager); err != nil {
			t.Fatalf("Close: %v", err)
		}
		waitCalls(t, p1, 1)
	})

	t.Run("capability rejection cleanup", func(t *testing.T) {
		p1 := newFakePipeline("p1")
		var manager *Manager
		duringFactory := -1
		manager = NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
			duringFactory = manager.Status().Builds
			*candidate = *candidateWith(p1, desired, false, false, errors.New("no ASR"), errors.New("no TTS"))
			return nil
		})
		if _, _, err := manager.Apply(context.Background(), settingsBoth); err == nil {
			t.Fatal("Apply returned nil, want rejection")
		}
		if duringFactory != 1 {
			t.Fatalf("builds during factory = %d, want 1", duringFactory)
		}
		if builds := manager.Status().Builds; builds != 0 {
			t.Fatalf("builds after rejection cleanup = %d, want 0", builds)
		}
		if p1.calls() != 1 {
			t.Fatalf("candidate shutdown calls = %d, want 1", p1.calls())
		}
		if err := mustClose(t, manager); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	t.Run("factory error with handle", func(t *testing.T) {
		p1 := newFakePipeline("p1")
		factoryErr := errors.New("factory failed")
		var manager *Manager
		duringFactory := -1
		manager = NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
			duringFactory = manager.Status().Builds
			*candidate = *candidateWith(p1, desired, true, false, nil, nil)
			return factoryErr
		})
		if _, _, err := manager.Apply(context.Background(), settingsASR); !errors.Is(err, factoryErr) {
			t.Fatalf("Apply error = %v, want factory error", err)
		}
		if duringFactory != 1 {
			t.Fatalf("builds during factory = %d, want 1", duringFactory)
		}
		if builds := manager.Status().Builds; builds != 0 {
			t.Fatalf("builds after factory cleanup = %d, want 0", builds)
		}
		if p1.calls() != 1 {
			t.Fatalf("candidate shutdown calls = %d, want 1", p1.calls())
		}
		if err := mustClose(t, manager); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	t.Run("panic before handle", func(t *testing.T) {
		var manager *Manager
		duringFactory := -1
		manager = NewManager(func(context.Context, Settings, *Candidate) error {
			duringFactory = manager.Status().Builds
			panic("construction panic")
		})
		_, _, err := manager.Apply(context.Background(), settingsASR)
		var buildFailure *BuildFailure
		if !errors.As(err, &buildFailure) {
			t.Fatalf("Apply error = %T %v, want *BuildFailure", err, err)
		}
		if duringFactory != 1 {
			t.Fatalf("builds during factory = %d, want 1", duringFactory)
		}
		if builds := manager.Status().Builds; builds != 0 {
			t.Fatalf("builds after panic = %d, want 0", builds)
		}
		if err := mustClose(t, manager); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	t.Run("panic after handle", func(t *testing.T) {
		p1 := newFakePipeline("p1")
		var manager *Manager
		duringFactory := -1
		manager = NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
			duringFactory = manager.Status().Builds
			*candidate = *candidateWith(p1, desired, true, false, nil, nil)
			panic("initialization panic")
		})
		_, _, err := manager.Apply(context.Background(), settingsASR)
		var buildFailure *BuildFailure
		if !errors.As(err, &buildFailure) {
			t.Fatalf("Apply error = %T %v, want *BuildFailure", err, err)
		}
		if duringFactory != 1 {
			t.Fatalf("builds during factory = %d, want 1", duringFactory)
		}
		if builds := manager.Status().Builds; builds != 0 {
			t.Fatalf("builds after panic cleanup = %d, want 0", builds)
		}
		if p1.calls() != 1 {
			t.Fatalf("candidate shutdown calls = %d, want 1", p1.calls())
		}
		if err := mustClose(t, manager); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})
}

func TestManagerM13PanicBoundary(t *testing.T) {
	t.Run("construction panic", func(t *testing.T) {
		manager := NewManager(func(context.Context, Settings, *Candidate) error {
			panic("construct boom")
		})
		_, _, err := manager.Apply(context.Background(), settingsASR)
		var buildFailure *BuildFailure
		if !errors.As(err, &buildFailure) {
			t.Fatalf("Apply error = %T %v, want *BuildFailure", err, err)
		}
		if buildFailure.Gen != 1 {
			t.Fatalf("BuildFailure generation = %d, want 1", buildFailure.Gen)
		}
		if buildFailure.Recovered != "construct boom" {
			t.Fatalf("recovered = %v, want construct boom", buildFailure.Recovered)
		}
		if buildFailure.CleanupErr != nil {
			t.Fatalf("cleanup error = %v, want nil", buildFailure.CleanupErr)
		}
		if builds := manager.Status().Builds; builds != 0 {
			t.Fatalf("builds = %d, want 0", builds)
		}
		if err := mustClose(t, manager); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	t.Run("initialization panic after handle", func(t *testing.T) {
		p1 := newFakePipeline("p1")
		manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
			*candidate = *candidateWith(p1, desired, true, false, nil, nil)
			panic("initialize boom")
		})
		_, _, err := manager.Apply(context.Background(), settingsASR)
		var buildFailure *BuildFailure
		if !errors.As(err, &buildFailure) {
			t.Fatalf("Apply error = %T %v, want *BuildFailure", err, err)
		}
		if buildFailure.CleanupErr != nil {
			t.Fatalf("cleanup error = %v, want nil", buildFailure.CleanupErr)
		}
		if p1.calls() != 1 {
			t.Fatalf("candidate shutdown calls = %d, want 1", p1.calls())
		}
		if builds := manager.Status().Builds; builds != 0 {
			t.Fatalf("builds = %d, want 0", builds)
		}
		if err := mustClose(t, manager); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	t.Run("cleanup panic", func(t *testing.T) {
		p1 := newFakePipeline("p1")
		p1.shutdownPanic = "cleanup boom"
		manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
			*candidate = *candidateWith(p1, desired, true, false, nil, nil)
			panic("initialize boom")
		})
		_, _, err := manager.Apply(context.Background(), settingsASR)
		var buildFailure *BuildFailure
		if !errors.As(err, &buildFailure) {
			t.Fatalf("Apply error = %T %v, want *BuildFailure", err, err)
		}
		var shutdownPanic *ShutdownPanic
		if !errors.As(buildFailure.CleanupErr, &shutdownPanic) {
			t.Fatalf("cleanup error = %T %v, want *ShutdownPanic", buildFailure.CleanupErr, buildFailure.CleanupErr)
		}
		if p1.calls() != 1 {
			t.Fatalf("candidate shutdown calls = %d, want 1", p1.calls())
		}
		if builds := manager.Status().Builds; builds != 0 {
			t.Fatalf("builds = %d, want 0", builds)
		}
		closeErr := mustClose(t, manager)
		if closeErr == nil {
			t.Fatal("Close returned nil, want retained cleanup panic")
		}
		if !errors.As(closeErr, &shutdownPanic) {
			t.Fatalf("Close error = %T %v, want *ShutdownPanic", closeErr, closeErr)
		}
		if status := manager.Status(); status.Phase != "closed_with_errors" {
			t.Fatalf("phase = %q, want closed_with_errors", status.Phase)
		}
	})

	t.Run("retired shutdown panic", func(t *testing.T) {
		p1 := newFakePipeline("p1")
		p1.shutdownPanic = "retired boom"
		p2 := newFakePipeline("p2")
		call := 0
		manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
			call++
			if call == 1 {
				*candidate = *candidateWith(p1, desired, true, true, nil, nil)
			} else {
				*candidate = *candidateWith(p2, desired, true, true, nil, nil)
			}
			return nil
		})
		mustApply(t, manager, settingsBoth)
		mustApply(t, manager, settingsBoth)
		waitCalls(t, p1, 1)

		err := mustClose(t, manager)
		if err == nil {
			t.Fatal("Close returned nil, want retained retired shutdown panic")
		}
		var shutdownPanic *ShutdownPanic
		if !errors.As(err, &shutdownPanic) {
			t.Fatalf("Close error = %T %v, want *ShutdownPanic", err, err)
		}
		if status := manager.Status(); status.Phase != "closed_with_errors" {
			t.Fatalf("phase = %q, want closed_with_errors", status.Phase)
		}
		if p1.calls() != 1 {
			t.Fatalf("p1 shutdown calls = %d, want 1", p1.calls())
		}
		if p2.calls() != 1 {
			t.Fatalf("p2 shutdown calls = %d, want 1", p2.calls())
		}
	})
}

func TestManagerM14RepeatedClose(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		p1 := newFakePipeline("p1")
		manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
			*candidate = *candidateWith(p1, desired, true, false, nil, nil)
			return nil
		})
		mustApply(t, manager, settingsASR)
		if err := mustClose(t, manager); err != nil {
			t.Fatalf("first Close: %v", err)
		}
		if status := manager.Status(); status.Phase != "closed" {
			t.Fatalf("phase = %q, want closed", status.Phase)
		}
		if err := mustClose(t, manager); err != nil {
			t.Fatalf("second Close: %v", err)
		}
		if p1.calls() != 1 {
			t.Fatalf("shutdown calls = %d, want 1", p1.calls())
		}
	})

	t.Run("closed with errors", func(t *testing.T) {
		shutdownErr := errors.New("shutdown failed")
		p1 := newFakePipeline("p1")
		p1.shutdownErr = shutdownErr
		manager := NewManager(func(_ context.Context, desired Settings, candidate *Candidate) error {
			*candidate = *candidateWith(p1, desired, true, false, nil, nil)
			return nil
		})
		mustApply(t, manager, settingsASR)
		first := mustClose(t, manager)
		if first == nil {
			t.Fatal("first Close returned nil, want shutdown error")
		}
		if status := manager.Status(); status.Phase != "closed_with_errors" {
			t.Fatalf("phase = %q, want closed_with_errors", status.Phase)
		}
		second := mustClose(t, manager)
		if second != first {
			t.Fatalf("second Close = %v, want identical terminal error %v", second, first)
		}
		if p1.calls() != 1 {
			t.Fatalf("shutdown calls = %d, want 1", p1.calls())
		}
	})
}
