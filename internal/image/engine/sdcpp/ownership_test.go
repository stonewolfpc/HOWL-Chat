package sdcpp

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"howl-chat/internal/image/types"
)

// F1 + F2: an error status is a failure, not a cancellation, and it does not hand
// the engine to the next request while the backend's state is unknown.
func TestHTTPErrorIsNotReportedAsAClientCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		http.Error(w, `{"error":"generate_image returned no results"}`, http.StatusInternalServerError)
	}))
	defer srv.Close()

	e := genEngine(t, srv)
	_, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job1"})
	if err == nil {
		t.Fatal("expected the 500 to fail the request")
	}
	if errors.Is(err, errClientCancelled) {
		t.Fatalf("an HTTP 500 must not be reported as a client cancellation: %v", err)
	}
	if !errors.Is(err, errEngineUncertain) {
		t.Fatalf("want errEngineUncertain for an unconfirmed backend, got %v", err)
	}
	if !strings.Contains(err.Error(), "500") {
		t.Fatalf("the server's status should survive into the message: %v", err)
	}
	if got := e.GenerationState(); got != "abandoned" {
		t.Fatalf("state = %q, want abandoned", got)
	}
	if _, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job2"}); !errors.Is(err, errEngineUnavailable) {
		t.Fatalf("a second generation must be refused after an unconfirmed failure, got %v", err)
	}
}

// F2: a 200 with no images is the historical "finished but produced nothing"
// state. It must not look like a cancellation, and it must not look like proof
// that the backend is free.
func TestEmptyResultKeepsTheEngineUnavailable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"images":[]}`)
	}))
	defer srv.Close()

	e := genEngine(t, srv)
	_, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job1"})
	if err == nil {
		t.Fatal("expected an empty result to fail")
	}
	if errors.Is(err, errClientCancelled) {
		t.Fatalf("an empty result is not a cancellation: %v", err)
	}
	if !errors.Is(err, errEngineUncertain) {
		t.Fatalf("want errEngineUncertain, got %v", err)
	}
	if got := e.GenerationState(); got != "abandoned" {
		t.Fatalf("state = %q, want abandoned", got)
	}
	if _, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job2"}); !errors.Is(err, errEngineUnavailable) {
		t.Fatalf("a second generation must be refused, got %v", err)
	}
}

// F2: losing the connection says nothing about whether the server is still
// working, so the slot keeps its owner.
func TestTransportFailureKeepsTheEngineUnavailable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	deadURL := srv.URL
	srv.Close()

	e := NewEngine()
	e.httpClient = newClient(deadURL)
	e.currentModel = &types.ModelInfo{Name: "m", Family: "z-image"}

	_, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job1"})
	if err == nil {
		t.Fatal("expected a transport failure")
	}
	if errors.Is(err, errClientCancelled) {
		t.Fatalf("a transport failure is not a cancellation: %v", err)
	}
	if !errors.Is(err, errEngineUncertain) {
		t.Fatalf("want errEngineUncertain, got %v", err)
	}
	if got := e.GenerationState(); got != "abandoned" {
		t.Fatalf("state = %q, want abandoned", got)
	}
}

// F2: a polling timeout means completion was never observed.
func TestPollTimeoutKeepsTheEngineUnavailable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/sdcpp/v1/img_gen":
			_, _ = io.Copy(io.Discard, r.Body)
			w.WriteHeader(http.StatusAccepted)
			_, _ = io.WriteString(w, `{"id":"j1"}`)
		case strings.HasPrefix(r.URL.Path, "/sdcpp/v1/jobs/"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"status":"running"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	initImage := filepath.Join(t.TempDir(), "init.png")
	if err := os.WriteFile(initImage, []byte("not really a png"), 0o644); err != nil {
		t.Fatal(err)
	}

	e := genEngine(t, srv)
	e.httpClient.pollTimeout = 100 * time.Millisecond

	_, err := e.GenerateImg2Img(types.ImageOptions{JobID: "job1", InitImagePath: initImage})
	if err == nil {
		t.Fatal("expected the poll timeout to fail the request")
	}
	if errors.Is(err, errClientCancelled) {
		t.Fatalf("a poll timeout is not a cancellation: %v", err)
	}
	if !errors.Is(err, errEngineUncertain) {
		t.Fatalf("want errEngineUncertain, got %v", err)
	}
	if !strings.Contains(err.Error(), "poll timeout") {
		t.Fatalf("the message should name the timeout: %v", err)
	}
	if got := e.GenerationState(); got != "abandoned" {
		t.Fatalf("state = %q, want abandoned", got)
	}
}

// F3: job B's load must not swallow job A's cancellation.
func TestCancellingAGenerationIsNotSwallowedByAnotherJobsLoad(t *testing.T) {
	entered := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		blockUntilClientGone(w, r, "still generating")
	}))
	defer srv.Close()

	e := genEngine(t, srv)
	// Job B owns a load while job A generates.
	e.beginLoad("jobB")
	defer e.endLoad()

	errCh := make(chan error, 1)
	go func() {
		_, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "jobA"})
		errCh <- err
	}()

	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("request never reached the server")
	}

	if err := e.CancelJob("jobA"); err != nil {
		t.Fatalf("CancelJob: %v", err)
	}

	select {
	case err := <-errCh:
		if !errors.Is(err, errClientCancelled) {
			t.Fatalf("job A's request was not cancelled, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the generating job's cancel was swallowed by another job's load")
	}

	if e.loadAborted() {
		t.Fatal("job A's cancel aborted job B's load")
	}
}

// F5: an empty id belongs to the legacy direct callers and can cancel nothing.
func TestEmptyJobIdentityCannotCancelWork(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		close(entered)
		<-release
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"images":["aGk="]}`)
	}))
	defer srv.Close()

	e := genEngine(t, srv)
	errCh := make(chan error, 1)
	go func() {
		// GenerateImage's direct path submits with no job id at all.
		_, err := e.GenerateTxt2Img(types.ImageOptions{})
		errCh <- err
	}()

	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("request never reached the server")
	}

	if err := e.CancelJob(""); err != nil {
		t.Fatalf("CancelJob with no identity: %v", err)
	}
	if got := e.GenerationState(); got != "active" {
		t.Fatalf("an empty job id cancelled work it does not own: state = %q", got)
	}

	close(release)
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("the direct request should have completed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("request did not return")
	}
	if got := e.GenerationState(); got != "idle" {
		t.Fatalf("state = %q after completion, want idle", got)
	}
}

// F6: a request unwinding late may only clear the progress id it owns.
func TestStaleProgressOwnerCannotClearANewerJob(t *testing.T) {
	c := newClient("http://127.0.0.1:1")

	oldToken := c.setActiveJob("old-job")
	newToken := c.setActiveJob("new-job")

	// The old request finally unwinds, after the new one claimed progress.
	c.resetActiveJob(oldToken)
	if got := c.activeJob(); got != "new-job" {
		t.Fatalf("a stale cleanup cleared the newer job's progress id: %q", got)
	}

	c.resetActiveJob(newToken)
	if got := c.activeJob(); got != "" {
		t.Fatalf("the owner could not clear its own id: %q", got)
	}
}

// F7: an unload is not something the user cancelled.
func TestShutdownIsNotReportedAsAClientCancellation(t *testing.T) {
	entered := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		blockUntilClientGone(w, r, "still generating")
	}))
	defer srv.Close()

	e := genEngine(t, srv)
	errCh := make(chan error, 1)
	go func() {
		_, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job1"})
		errCh <- err
	}()

	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("request never reached the server")
	}

	if err := e.Shutdown(); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}

	select {
	case err := <-errCh:
		if errors.Is(err, errClientCancelled) {
			t.Fatalf("unload reported as a client cancellation: %v", err)
		}
		if !errors.Is(err, errEngineStopped) {
			t.Fatalf("want errEngineStopped after unload, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("unload did not end the request")
	}
}

// F2: a failure that never reached the backend releases the slot, and it is
// reported as itself rather than as a cancellation or an unknown state.
func TestLocalFailureReleasesTheSlotWithoutBlame(t *testing.T) {
	e := genEngine(t, httptest.NewServer(http.NotFoundHandler()))

	_, err := e.GenerateInpaint(types.ImageOptions{JobID: "job1"})
	if err == nil {
		t.Fatal("expected the missing init image to fail")
	}
	if errors.Is(err, errClientCancelled) || errors.Is(err, errEngineUncertain) || errors.Is(err, errEngineStopped) {
		t.Fatalf("a local input failure is none of those: %v", err)
	}
	if !strings.Contains(err.Error(), "init image") {
		t.Fatalf("the local reason should survive: %v", err)
	}
	if got := e.GenerationState(); got != "idle" {
		t.Fatalf("state = %q, want idle: nothing was sent", got)
	}
	// The slot is free: a fresh request can claim it. (The test server answers 404,
	// so actually generating here would land in the ambiguous branch instead.)
	h, err := e.beginGeneration("job2")
	if err != nil {
		t.Fatalf("a local failure must not wedge the engine: %v", err)
	}
	e.endGeneration(h, genOutcomeCompleted)
}
