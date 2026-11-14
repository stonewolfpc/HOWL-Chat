package helper

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

func newFakeClientWithHelpers(t *testing.T, opts ClientOptions, helpers ...*fakeHelper) (*Client, *Supervisor) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "howl-tts-helper.exe")
	if err := os.WriteFile(path, []byte("test double"), 0o644); err != nil {
		t.Fatal(err)
	}
	next := 0
	sup := NewSupervisor(SupervisorOptions{
		Discoverer: Discoverer{LookupEnv: func(string) (string, bool) { return path, true }},
		Factory: func(context.Context, LaunchSpec) (Process, error) {
			if next >= len(helpers) {
				return nil, errors.New("unexpected helper start")
			}
			h := helpers[next]
			next++
			return h.p, nil
		},
	})
	c, err := NewClient(sup, opts)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sup.Terminate(context.Background(), 50*time.Millisecond) })
	return c, sup
}

func TestFailAllFullOutcomeDoesNotBlockOrSkipKill(t *testing.T) {
	h := newFakeHelper()
	c, sup := newFakeClientWithHelpers(t, fastClientOptions(), h)
	if _, err := sup.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	ok := true
	full := &pendingRequest{ch: make(chan pendingOutcome, 1)}
	full.ch <- pendingOutcome{resp: &responseEnvelope{
		Protocol: ProtocolVersion, ID: "req-full", OK: &ok, Result: []byte(`{}`),
	}}
	empty := &pendingRequest{ch: make(chan pendingOutcome, 1)}
	c.mu.Lock()
	c.pending = map[string]*pendingRequest{"req-full": full, "req-empty": empty}
	c.mu.Unlock()

	done := make(chan struct{})
	go func() {
		c.protocolViolation(newFailure("protocol_violation", "second fatal event"))
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("failAll/protocolViolation blocked on a full outcome channel")
	}

	if h.p.waitCount.Load() != 1 {
		t.Fatalf("Wait calls = %d, want 1", h.p.waitCount.Load())
	}
	select {
	case <-h.p.closed:
	default:
		t.Fatal("process termination did not complete")
	}
	c.mu.Lock()
	readFailure := c.readFailure
	c.mu.Unlock()
	if readFailure == nil || readFailure.Code != "protocol_violation" {
		t.Fatalf("read failure = %+v, want protocol_violation", readFailure)
	}
	select {
	case out := <-empty.ch:
		if out.failure == nil || out.failure.Code != "protocol_violation" {
			t.Fatalf("empty request outcome = %+v, want protocol_violation", out)
		}
	default:
		t.Fatal("empty request did not receive the fatal outcome")
	}
	select {
	case out := <-full.ch:
		if out.resp == nil || out.failure != nil {
			t.Fatalf("full request outcome = %+v, want preserved success", out)
		}
	default:
		t.Fatal("full request's already-delivered outcome was lost")
	}
}

func TestFailedHandshakeStopsOldReaderAndAllowsRetry(t *testing.T) {
	h1 := newFakeHelper()
	h1.setHelloError("invalid_request", "handshake rejected")
	h2 := newFakeHelper()
	c, _ := newFakeClientWithHelpers(t, fastClientOptions(), h1, h2)

	err := c.Start(context.Background())
	f, ok := AsFailure(err)
	if !ok || f.Code != "invalid_request" {
		t.Fatalf("first Start error = %v, want invalid_request", err)
	}
	if h1.p.waitCount.Load() != 1 {
		t.Fatalf("first helper Wait calls = %d, want 1", h1.p.waitCount.Load())
	}
	select {
	case <-h1.p.closed:
	default:
		t.Fatal("failed handshake left the first process alive")
	}
	if got := c.sup.Status().State; got != types.TTSStateUnavailable {
		t.Fatalf("failed handshake state = %s, want %s", got, types.TTSStateUnavailable)
	}
	c.mu.Lock()
	readerDone := c.readerDone
	c.mu.Unlock()
	select {
	case <-readerDone:
	default:
		t.Fatal("failed handshake left the first reader running")
	}

	if err := c.Start(context.Background()); err != nil {
		t.Fatalf("retry Start: %v", err)
	}
	if h2.p.waitCount.Load() != 1 {
		t.Fatalf("retry helper Wait calls = %d, want 1", h2.p.waitCount.Load())
	}
}

func TestRejectedCleanupTerminatesBeforeRemovingOutput(t *testing.T) {
	h := newFakeHelper()
	c, sup := newFakeClientWithHelpers(t, fastClientOptions(), h)
	if _, err := sup.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	final := filepath.Join(dir, "rejected.wav")
	if _, _, err := writeTestWAV(final, 24000, 240); err != nil {
		t.Fatal(err)
	}
	observed := make(chan bool, 1)
	h.p.waitHook = func() {
		_, err := os.Stat(final)
		observed <- err == nil
	}
	c.cleanupRejected(dir, "rejected.wav", false)
	select {
	case present := <-observed:
		if !present {
			t.Fatal("cleanup removed output before the helper exit was observed")
		}
	case <-time.After(time.Second):
		t.Fatal("helper Wait observer did not run")
	}
	if _, err := os.Stat(final); !os.IsNotExist(err) {
		t.Fatalf("rejected output still exists: %v", err)
	}
}
