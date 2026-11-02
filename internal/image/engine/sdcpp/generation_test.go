package sdcpp

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"howl-chat/internal/image/types"
)

// blockUntilClientGone blocks a handler until the client closes the connection.
// The body is drained first on purpose: the Go server only starts watching the
// connection for an abort once the request body has been consumed, so a handler
// that never reads it can never observe the client leaving.
func blockUntilClientGone(w http.ResponseWriter, r *http.Request, msg string) bool {
	_, _ = io.Copy(io.Discard, r.Body)
	select {
	case <-r.Context().Done():
		return true
	case <-time.After(10 * time.Second):
		http.Error(w, msg, http.StatusGatewayTimeout)
		return false
	}
}

// genEngine builds an engine pointed at a test server with a model loaded, so
// the generation paths get past their guards without touching a real backend.
func genEngine(t *testing.T, srv *httptest.Server) *Engine {
	t.Helper()
	e := NewEngine()
	e.httpClient = newClient(srv.URL)
	e.currentModel = &types.ModelInfo{Name: "test-model", Family: "z-image"}
	return e
}

// 1 + 2: the request is cancelled through its context, and the result never
// claims the backend stopped.
func TestTxt2ImgIsCancelledThroughItsRequestContext(t *testing.T) {
	entered := make(chan struct{})
	sawEnd := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		if blockUntilClientGone(w, r, "still generating") {
			close(sawEnd)
		}
	}))
	defer srv.Close()

	e := genEngine(t, srv)
	errCh := make(chan error, 1)
	go func() {
		_, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job1", Prompt: "x"})
		errCh <- err
	}()

	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("request never reached the test server")
	}
	if got := e.GenerationState(); got != "active" {
		t.Fatalf("state = %q, want active", got)
	}
	if err := e.CancelJob("job1"); err != nil {
		t.Fatalf("CancelJob: %v", err)
	}

	select {
	case err := <-errCh:
		if !errors.Is(err, errClientCancelled) {
			t.Fatalf("want errClientCancelled, got %v", err)
		}
		if !strings.Contains(err.Error(), "may still be running") {
			t.Fatalf("message must not claim the backend stopped: %q", err.Error())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("GenerateTxt2Img did not return after cancellation")
	}

	select {
	case <-sawEnd:
	case <-time.After(5 * time.Second):
		t.Fatal("the server never saw the request context end, so the connection stayed open")
	}

	// 3: the abandoned slot refuses a second request instead of overlapping it.
	if got := e.GenerationState(); got != "abandoned" {
		t.Fatalf("state = %q, want abandoned", got)
	}
	if _, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job2"}); !errors.Is(err, errEngineUnavailable) {
		t.Fatalf("second generation should be refused while abandoned, got %v", err)
	}
}

// 3: an overlapping request is refused, and the first still completes.
func TestSecondGenerationIsRefusedWhileTheFirstIsActive(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() { close(entered) })
		<-release
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"images":["aGk="]}`)
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
		t.Fatal("first request never reached the server")
	}

	if _, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job2"}); !errors.Is(err, errEngineBusy) {
		t.Fatalf("want errEngineBusy for an overlapping request, got %v", err)
	}

	close(release)
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("first generation failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("first generation never finished")
	}

	// 4: a confirmed normal completion releases ownership.
	if got := e.GenerationState(); got != "idle" {
		t.Fatalf("state = %q after completion, want idle", got)
	}
	if _, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job3"}); err != nil {
		t.Fatalf("a fresh generation after completion should be allowed, got %v", err)
	}
}

// 5: a confirmed process exit is the one completion accepted without the backend.
func TestOwnershipIsReleasedAfterAConfirmedProcessExit(t *testing.T) {
	e := NewEngine()
	e.httpClient = newClient("http://127.0.0.1:1")
	e.currentModel = &types.ModelInfo{Name: "m", Family: "z-image"}

	// An abandoned slot with no exit yet: unavailable.
	e.genMu.Lock()
	e.genState = genAbandoned
	e.genOwner = "job1"
	e.genMu.Unlock()
	if _, err := e.beginGeneration("job2"); !errors.Is(err, errEngineUnavailable) {
		t.Fatalf("want errEngineUnavailable before any exit, got %v", err)
	}

	// The child is confirmed gone, so its request died with it.
	done := make(chan struct{})
	close(done)
	e.procMu.Lock()
	e.procDone = done
	e.procMu.Unlock()

	h, err := e.beginGeneration("job2")
	if err != nil {
		t.Fatalf("a confirmed exit must release ownership, got %v", err)
	}
	e.endGeneration(h, genOutcomeCompleted)
	if got := e.GenerationState(); got != "idle" {
		t.Fatalf("state = %q, want idle", got)
	}
}

// 7: one job's cancel never touches another job's request.
func TestCancelIsScopedToOneJob(t *testing.T) {
	entered := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		blockUntilClientGone(w, r, "still generating")
	}))
	defer srv.Close()

	e := genEngine(t, srv)
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

	if err := e.CancelJob("jobB"); err != nil {
		t.Fatalf("CancelJob for a job that owns nothing: %v", err)
	}
	if got := e.GenerationState(); got != "active" {
		t.Fatalf("jobB's cancel changed jobA's state to %q", got)
	}

	if err := e.CancelJob("jobA"); err != nil {
		t.Fatalf("CancelJob: %v", err)
	}
	select {
	case err := <-errCh:
		if !errors.Is(err, errClientCancelled) {
			t.Fatalf("want errClientCancelled after the owner's cancel, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("owner's cancel did not end the request")
	}
}

// 8: repeating a cancel is safe and does not reopen the slot.
func TestRepeatedCancelIsSafe(t *testing.T) {
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

	for i := 0; i < 3; i++ {
		if err := e.CancelJob("job1"); err != nil {
			t.Fatalf("cancel #%d: %v", i+1, err)
		}
	}
	<-errCh

	if got := e.GenerationState(); got != "abandoned" {
		t.Fatalf("state = %q after repeated cancels, want abandoned", got)
	}
	if _, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job2"}); !errors.Is(err, errEngineUnavailable) {
		t.Fatalf("repeated cancels reopened the slot: %v", err)
	}
}

// 6: a cancelled async job must not poison the next one. This is the
// resetCancel/c.cancelled defect the audit found, exercised end to end through
// the img2img job-id path.
func TestAsyncCancellationDoesNotPoisonTheNextJob(t *testing.T) {
	var mu sync.Mutex
	n := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/sdcpp/v1/img_gen":
			mu.Lock()
			n++
			id := fmt.Sprintf("j%d", n)
			mu.Unlock()
			w.WriteHeader(http.StatusAccepted)
			_, _ = fmt.Fprintf(w, `{"id":%q}`, id)
		case strings.HasPrefix(r.URL.Path, "/sdcpp/v1/jobs/"):
			id := strings.TrimPrefix(r.URL.Path, "/sdcpp/v1/jobs/")
			if id == "j1" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{"status":"running"}`)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"status":"completed","result":{"images":["aGk="]}}`)
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

	// Job 1 polls forever until it is cancelled.
	errCh := make(chan error, 1)
	go func() {
		_, err := e.GenerateImg2Img(types.ImageOptions{JobID: "job1", InitImagePath: initImage})
		errCh <- err
	}()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		started := n >= 1
		mu.Unlock()
		if started {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := e.CancelJob("job1"); err != nil {
		t.Fatalf("CancelJob: %v", err)
	}
	select {
	case err := <-errCh:
		if !errors.Is(err, errClientCancelled) {
			t.Fatalf("want errClientCancelled, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cancelled async job did not return")
	}

	// A fresh child is what ends the abandoned slot in production; here the same
	// engine is made reusable the way a confirmed exit would.
	e.resetGeneration(closeForReload)

	if _, err := e.GenerateImg2Img(types.ImageOptions{JobID: "job2", InitImagePath: initImage}); err != nil {
		t.Fatalf("the next job inherited cancellation state: %v", err)
	}
}
