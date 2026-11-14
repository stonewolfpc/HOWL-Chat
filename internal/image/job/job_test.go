package job

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/image/engine"
	"howl-chat/internal/image/types"
)

// fakeEngine reports a permanently broken progress endpoint, and lets a test
// decide when the backend counts as stopped and when a generation finishes.
type fakeEngine struct {
	engine.ImageEngine
	alive     atomic.Bool
	progress  atomic.Int32
	generate  chan struct{}
	generated atomic.Int32

	mu        sync.Mutex
	cancelIDs []string
}

func (f *fakeEngine) GetProgress() (types.ProgressInfo, error) {
	f.progress.Add(1)
	return types.ProgressInfo{}, errors.New("progress endpoint unavailable")
}

func (f *fakeEngine) IsRunning() bool { return f.alive.Load() }

// CancelJob records which job the cancel belonged to: the engine scopes
// cancellation by that identity.
func (f *fakeEngine) CancelJob(jobID string) error {
	f.mu.Lock()
	f.cancelIDs = append(f.cancelIDs, jobID)
	f.mu.Unlock()
	return nil
}

func (f *fakeEngine) cancelledJobs() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.cancelIDs...)
}

func (f *fakeEngine) GenerateTxt2Img(types.ImageOptions) (*types.ImageResult, error) {
	<-f.generate
	f.generated.Add(1)
	return &types.ImageResult{Path: "x.png"}, nil
}

func newFakeEngine() *fakeEngine {
	f := &fakeEngine{generate: make(chan struct{})}
	f.alive.Store(true)
	return f
}

// The first progress error used to set a flag that skipped the liveness check
// for the rest of the job, so a dead sd-server left the UI on "Generating...
// Step 8 of 8" forever with no failure ever recorded.
func TestJobFailsWhenTheEngineDiesEvenIfProgressErrors(t *testing.T) {
	eng := newFakeEngine()
	m := NewManager(eng)
	id := m.Start(types.ImageOptions{Mode: "txt2img"}, nil)

	waitUntil(t, 3*time.Second, func() bool { return eng.progress.Load() > 0 })
	eng.alive.Store(false)

	var job *types.ImageJob
	waitUntil(t, 5*time.Second, func() bool {
		job = m.Status(id)
		return job != nil && job.Status == "failed"
	})
	if !strings.Contains(job.Error, "stopped before finishing") {
		t.Fatalf("want a stopped-engine failure, got %q", job.Error)
	}
}

// A slow load used to block the StartImageJob RPC for as long as the model took
// to load. The id must come back while loading is still running.
func TestStartReturnsWhileLoadingIsStillRunning(t *testing.T) {
	eng := newFakeEngine()
	m := NewManager(eng)

	release := make(chan struct{})
	started := make(chan struct{})
	prepare := func(string) error {
		close(started)
		<-release
		return nil
	}

	returned := make(chan string, 1)
	go func() { returned <- m.Start(types.ImageOptions{Mode: "txt2img"}, prepare) }()

	var id string
	select {
	case id = <-returned:
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return while prepare was still blocked")
	}
	if id == "" {
		t.Fatal("empty job id")
	}

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("prepare never ran")
	}
	if job := m.Status(id); job == nil || job.Status != "loading" {
		t.Fatalf("want status loading while the model loads, got %+v", job)
	}

	close(release)
	close(eng.generate)
	waitUntil(t, 3*time.Second, func() bool { return m.Status(id).Status == "completed" })
}

func TestLoadingSucceedsAndAdvancesToGeneration(t *testing.T) {
	eng := newFakeEngine()
	close(eng.generate)
	m := NewManager(eng)

	id := m.Start(types.ImageOptions{Mode: "txt2img"}, func(string) error { return nil })
	waitUntil(t, 3*time.Second, func() bool { return m.Status(id).Status == "completed" })

	if eng.generated.Load() != 1 {
		t.Fatalf("generation ran %d times, want 1", eng.generated.Load())
	}
	if got := m.Status(id).Result; got == nil || got.Path != "x.png" {
		t.Fatalf("want the generated artifact on the job, got %+v", got)
	}
}

// A startup failure is a terminal job failure, not an error thrown out of the
// call that started the job.
func TestLoadingFailureIsTerminal(t *testing.T) {
	eng := newFakeEngine()
	m := NewManager(eng)

	boom := errors.New("model file is not a checkpoint")
	id := m.Start(types.ImageOptions{Mode: "txt2img"}, func(string) error { return boom })

	var job *types.ImageJob
	waitUntil(t, 3*time.Second, func() bool {
		job = m.Status(id)
		return job != nil && job.Status == "failed"
	})
	if !strings.Contains(job.Error, boom.Error()) {
		t.Fatalf("want the load error on the job, got %q", job.Error)
	}
	if eng.generated.Load() != 0 {
		t.Fatal("generation ran after the load failed")
	}
}

// Cancelling during loading must stop the job there, even when the load itself
// then succeeds: a late readiness may not revive an abandoned job.
func TestCancelDuringLoadingPreventsGeneration(t *testing.T) {
	eng := newFakeEngine()
	close(eng.generate)
	m := NewManager(eng)

	release := make(chan struct{})
	started := make(chan struct{})
	prepare := func(string) error {
		close(started)
		<-release
		return nil
	}

	id := m.Start(types.ImageOptions{Mode: "txt2img"}, prepare)
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("prepare never ran")
	}

	if !m.Cancel(id) {
		t.Fatal("Cancel refused a job that was still loading")
	}
	if got := eng.cancelledJobs(); len(got) == 0 {
		t.Fatal("cancellation never reached the engine")
	}

	close(release) // the load "succeeds" after the cancel
	waitUntil(t, 3*time.Second, func() bool {
		return m.Status(id).Status == "cancelled" && eng.generated.Load() == 0
	})
}

// A completion that lands after a cancellation must not overwrite it.
func TestCancelDuringGenerationKeepsTerminalState(t *testing.T) {
	eng := newFakeEngine()
	m := NewManager(eng)

	id := m.Start(types.ImageOptions{Mode: "txt2img"}, nil)
	waitUntil(t, 3*time.Second, func() bool { return m.Status(id).Status == "generating" })

	if !m.Cancel(id) {
		t.Fatal("Cancel refused a generating job")
	}
	close(eng.generate) // generation finishes anyway, as a real engine might
	waitUntil(t, 2*time.Second, func() bool { return eng.generated.Load() == 1 })

	time.Sleep(50 * time.Millisecond)
	if got := m.Status(id).Status; got != "cancelled" {
		t.Fatalf("status = %q, want cancelled; a late completion revived the job", got)
	}
}

// The engine scopes cancellation by job identity, so the manager has to send it.
func TestCancelCarriesTheJobID(t *testing.T) {
	eng := newFakeEngine()
	m := NewManager(eng)

	id := m.Start(types.ImageOptions{Mode: "txt2img"}, nil)
	waitUntil(t, 3*time.Second, func() bool { return m.Status(id).Status == "generating" })
	if !m.Cancel(id) {
		t.Fatal("Cancel refused a generating job")
	}

	got := eng.cancelledJobs()
	if len(got) != 1 || got[0] != id {
		t.Fatalf("engine was told to cancel %v, want exactly [%s]", got, id)
	}
}

// Cancellation state must not leak into the next job: the previous one is
// terminal, the new one runs normally.
func TestNextJobCompletesAfterACancel(t *testing.T) {
	eng := newFakeEngine()
	m := NewManager(eng)

	first := m.Start(types.ImageOptions{Mode: "txt2img"}, nil)
	waitUntil(t, 3*time.Second, func() bool { return m.Status(first).Status == "generating" })
	if !m.Cancel(first) {
		t.Fatal("Cancel refused a generating job")
	}
	close(eng.generate) // both generations return from here on

	second := m.Start(types.ImageOptions{Mode: "txt2img"}, nil)
	waitUntil(t, 3*time.Second, func() bool { return m.Status(second).Status == "completed" })

	if got := m.Status(first).Status; got != "cancelled" {
		t.Fatalf("first job = %q, want cancelled", got)
	}
}

func TestCancelAllClosesAdmissionAndCancelsActiveJob(t *testing.T) {
	eng := newFakeEngine()
	m := NewManager(eng)

	id := m.Start(types.ImageOptions{Mode: "txt2img"}, nil)
	waitUntil(t, 3*time.Second, func() bool { return m.Status(id).Status == "generating" })

	m.CancelAll()

	if got := m.Status(id).Status; got != "cancelled" {
		t.Fatalf("active job status = %q, want cancelled", got)
	}
	if got := eng.cancelledJobs(); len(got) != 1 || got[0] != id {
		t.Fatalf("engine cancellations = %v, want [%s]", got, id)
	}
	if next := m.Start(types.ImageOptions{Mode: "txt2img"}, nil); next != "" {
		t.Fatalf("Start returned %q after CancelAll, want refusal", next)
	}
}

func waitUntil(t *testing.T, d time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("condition not met within %v", d)
}
