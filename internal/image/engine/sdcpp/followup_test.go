package sdcpp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"howl-chat/internal/image/types"
)

// R1: every pre-dispatch input failure is one contract, whatever the mode.
// Before this, img2img returned plain errors for a missing or unreadable init
// image while inpaint wrapped its equivalent checks, so the same user mistake
// either freed the engine or wedged it depending on the tab.
func TestPreDispatchInputFailuresShareOneContract(t *testing.T) {
	e := genEngine(t, httptest.NewServer(http.NotFoundHandler()))

	dir := t.TempDir()
	validInit := filepath.Join(dir, "init.png")
	if err := os.WriteFile(validInit, []byte("not really a png"), 0o644); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(dir, "missing.png")

	cases := []struct {
		name string
		run  func() error
	}{
		{"img2img without an init image", func() error {
			_, err := e.GenerateImg2Img(types.ImageOptions{JobID: "j1"})
			return err
		}},
		{"img2img with an unreadable init image", func() error {
			_, err := e.GenerateImg2Img(types.ImageOptions{JobID: "j1", InitImagePath: missing})
			return err
		}},
		{"inpaint without an init image", func() error {
			_, err := e.GenerateInpaint(types.ImageOptions{JobID: "j1"})
			return err
		}},
		{"inpaint without a mask", func() error {
			_, err := e.GenerateInpaint(types.ImageOptions{JobID: "j1", InitImagePath: validInit})
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run()
			if err == nil {
				t.Fatal("expected the input failure to be reported")
			}
			if !errors.Is(err, errNotDispatched) {
				t.Fatalf("want errNotDispatched, got %v", err)
			}
			if errors.Is(err, errEngineUncertain) {
				t.Fatalf("a local input failure must not claim the backend is busy: %v", err)
			}
			if errors.Is(err, errClientCancelled) {
				t.Fatalf("a local input failure is not a cancellation: %v", err)
			}
			if got := e.GenerationState(); got != "idle" {
				t.Fatalf("state = %q, want idle: nothing was dispatched", got)
			}
		})
	}
}

// R2: the close reason belongs to the generation being settled. A request that
// starts in the window between an earlier generation releasing the slot and
// describing its error must not rewrite that error.
func TestSettledGenerationKeepsItsOwnCloseReason(t *testing.T) {
	e := NewEngine()
	e.httpClient = newClient("http://127.0.0.1:1")
	e.currentModel = &types.ModelInfo{Name: "m", Family: "z-image"}

	a, err := e.beginGeneration("jobA")
	if err != nil {
		t.Fatalf("beginGeneration: %v", err)
	}
	// The unload closes A's request and frees the slot.
	e.resetGeneration(closeForShutdown)

	// The exact race: B claims the slot before A has described itself.
	b, err := e.beginGeneration("jobB")
	if err != nil {
		t.Fatalf("a request after the reset should be allowed: %v", err)
	}
	defer e.endGeneration(b, genOutcomeCompleted)

	err = e.settleGeneration(a, errors.New("context canceled"))
	if errors.Is(err, errClientCancelled) {
		t.Fatalf("B's start rewrote A's reason: %v", err)
	}
	if !errors.Is(err, errEngineStopped) {
		t.Fatalf("want errEngineStopped for the unload, got %v", err)
	}
}

// R2 support: an ordinary client cancellation must stay distinguishable from an
// engine stop after the reason moved onto the handle.
func TestClientCancellationStaysDistinctFromAnEngineStop(t *testing.T) {
	e := NewEngine()
	e.currentModel = &types.ModelInfo{Name: "m", Family: "z-image"}

	h, err := e.beginGeneration("jobA")
	if err != nil {
		t.Fatalf("beginGeneration: %v", err)
	}
	if !e.cancelGeneration("jobA") {
		t.Fatal("the owner's cancel did not match")
	}
	err = e.settleGeneration(h, errors.New("context canceled"))
	if !errors.Is(err, errClientCancelled) {
		t.Fatalf("want errClientCancelled, got %v", err)
	}
	if errors.Is(err, errEngineStopped) {
		t.Fatalf("a user cancel was reported as an engine stop: %v", err)
	}
}

// R4: the transport error underneath an uncertain outcome stays reachable.
func TestUncertainErrorKeepsTheTransportError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	deadURL := srv.URL
	srv.Close()

	e := NewEngine()
	e.httpClient = newClient(deadURL)
	e.currentModel = &types.ModelInfo{Name: "m", Family: "z-image"}

	_, err := e.GenerateTxt2Img(types.ImageOptions{JobID: "job1"})
	if !errors.Is(err, errEngineUncertain) {
		t.Fatalf("want errEngineUncertain, got %v", err)
	}
	var opErr *net.OpError
	if !errors.As(err, &opErr) {
		t.Fatalf("the transport error is no longer unwrappable: %v", err)
	}
}

// R4: a local file failure keeps its *os.PathError through both wrappers.
func TestNotDispatchedErrorKeepsTheFileError(t *testing.T) {
	dir := t.TempDir()
	initImage := filepath.Join(dir, "init.png")
	if err := os.WriteFile(initImage, []byte("not really a png"), 0o644); err != nil {
		t.Fatal(err)
	}
	missingMask := filepath.Join(dir, "missing-mask.png")

	e := genEngine(t, httptest.NewServer(http.NotFoundHandler()))
	_, err := e.GenerateInpaint(types.ImageOptions{JobID: "job1", InitImagePath: initImage, MaskPath: missingMask})
	if !errors.Is(err, errNotDispatched) {
		t.Fatalf("want errNotDispatched, got %v", err)
	}
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("the file error is no longer unwrappable: %v", err)
	}
	if errors.Is(err, errEngineUncertain) || errors.Is(err, errClientCancelled) {
		t.Fatalf("a local file failure is neither: %v", err)
	}
	if got := e.GenerationState(); got != "idle" {
		t.Fatalf("state = %q, want idle", got)
	}
}

// R6: a batch whose first item fails locally is notDispatched, and nothing was
// sent to the backend.
func TestBatchFirstItemLocalFailureIsNotDispatched(t *testing.T) {
	var requests atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		http.NotFound(w, r)
	}))
	defer srv.Close()

	e := genEngine(t, srv)
	missing := filepath.Join(t.TempDir(), "missing.png")
	_, err := e.GenerateImg2Img(types.ImageOptions{JobID: "job1", BatchCount: 2, InitImagePath: missing})

	if !errors.Is(err, errNotDispatched) {
		t.Fatalf("want errNotDispatched through the batch wrapper, got %v", err)
	}
	if got := requests.Load(); got != 0 {
		t.Fatalf("nothing should have been dispatched, got %d requests", got)
	}
	if got := e.GenerationState(); got != "idle" {
		t.Fatalf("state = %q, want idle", got)
	}
}

// R6: when an earlier item succeeded and a later item fails locally before
// dispatch, the batch keeps its earlier result, the slot is released, and the
// engine is not left abandoned for work that never existed.
func TestBatchLaterLocalFailureKeepsTheEarlierResultAndTheSlot(t *testing.T) {
	initImage := filepath.Join(t.TempDir(), "init.png")
	if err := os.WriteFile(initImage, []byte("not really a png"), 0o644); err != nil {
		t.Fatal(err)
	}

	var generations atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/sdcpp/v1/img_gen":
			_, _ = io.Copy(io.Discard, r.Body)
			n := generations.Add(1)
			if n == 1 {
				// Item 2 will find nothing to read and fail before dispatch.
				_ = os.Remove(initImage)
			}
			w.WriteHeader(http.StatusAccepted)
			_, _ = fmt.Fprintf(w, `{"id":"j%d"}`, n)
		case strings.HasPrefix(r.URL.Path, "/sdcpp/v1/jobs/"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"status":"completed","result":{"images":["aGk="]}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	e := genEngine(t, srv)
	result, err := e.GenerateImg2Img(types.ImageOptions{JobID: "job1", BatchCount: 2, InitImagePath: initImage})
	if err != nil {
		t.Fatalf("a later local failure must not discard the image already produced: %v", err)
	}
	if result == nil || len(result.Images) != 1 {
		t.Fatalf("want the first item's image, got %+v", result)
	}
	if got := generations.Load(); got != 1 {
		t.Fatalf("expected exactly one dispatched item, got %d", got)
	}
	if got := e.GenerationState(); got != "idle" {
		t.Fatalf("state = %q, want idle: the batch completed with a result", got)
	}
}
