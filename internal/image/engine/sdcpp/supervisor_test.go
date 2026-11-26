package sdcpp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

func testEngine(baseURL string) *Engine {
	e := NewEngine()
	e.httpClient = newClient(baseURL)
	return e
}

// krea2's load printed "new_sd_ctx_t failed" within seconds while the app kept
// polling until the ten-minute load timeout and then blamed the binary version.
func TestWaitReadyOrExitReturnsTheExitWhenTheChildDies(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("spawns cmd.exe")
	}
	e := testEngine("http://127.0.0.1:1") // nothing listens here
	cmd := exec.Command("cmd", "/c", "exit", "3")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	e.serverProc = cmd.Process
	e.startProcessReaper(cmd)

	start := time.Now()
	err := e.waitReadyOrExit(30 * time.Second)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected the child's exit to be reported, got ready")
	}
	if elapsed > 10*time.Second {
		t.Fatalf("took %v: a dead child must not be waited on for the whole load timeout", elapsed)
	}
	if !strings.Contains(err.Error(), "exited") {
		t.Fatalf("want a message about the exit, got %v", err)
	}
	if !errors.Is(err, errProcessExited) {
		t.Fatalf("callers must be able to tell an exit from a timeout, got %v", err)
	}
}

func TestWaitReadyOrExitSucceedsWhenTheServerAnswers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdapi/v1/samplers" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte("[]"))
	}))
	defer srv.Close()

	e := testEngine(srv.URL)
	if err := e.waitReadyOrExit(5 * time.Second); err != nil {
		t.Fatalf("want ready, got %v", err)
	}
}

// Liveness is a process question. It used to be a three-second HTTP readiness
// probe, which the progress poller called every 500ms.
func TestIsRunningIsAProcessQuestion(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("spawns cmd.exe")
	}
	e := testEngine("http://127.0.0.1:1")
	cmd := exec.Command("cmd", "/c", "exit", "0")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	e.serverProc = cmd.Process
	e.startProcessReaper(cmd)

	if !e.IsRunning() {
		t.Fatal("reported not running while the child was still alive")
	}
	deadline := time.Now().Add(10 * time.Second)
	for e.IsRunning() && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if e.IsRunning() {
		t.Fatal("still reported running after the child exited")
	}
}

// Cancelling a load has no HTTP job to cancel: only stopping the child ends it.
func TestCancelJobStopsALoadInProgress(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("spawns cmd.exe")
	}
	e := testEngine("http://127.0.0.1:1")
	cmd := exec.Command("cmd", "/c", "ping", "-n", "30", "127.0.0.1")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	e.setServerProc(cmd.Process)
	e.startProcessReaper(cmd)
	e.beginLoad("job_test")
	defer e.endLoad()

	if err := e.CancelJob("job_test"); err != nil {
		t.Fatalf("CancelJob: %v", err)
	}
	if !e.loadAborted() {
		t.Fatal("cancellation did not mark the load aborted")
	}
	deadline := time.Now().Add(5 * time.Second)
	for e.IsRunning() && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if e.IsRunning() {
		t.Fatal("the loading child survived cancellation")
	}

	e.beginLoad("job_test") // a fresh load must not inherit the abort
	if e.loadAborted() {
		t.Fatal("abort flag leaked into the next load")
	}
}

// A cancel that lands once the child is already reading weights must stop the
// load, not let it finish for a job nobody is waiting for.
func TestWaitReadyOrExitStopsWhenTheLoadIsCancelled(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("spawns cmd.exe")
	}
	e := testEngine("http://127.0.0.1:1")
	cmd := exec.Command("cmd", "/c", "ping", "-n", "30", "127.0.0.1")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	e.setServerProc(cmd.Process)
	e.startProcessReaper(cmd)
	e.beginLoad("job_test")
	defer e.endLoad()

	go func() {
		time.Sleep(150 * time.Millisecond)
		e.abortLoad("job_test")
	}()

	start := time.Now()
	err := e.waitReadyOrExit(30 * time.Second)
	if !errors.Is(err, errLoadCancelled) {
		t.Fatalf("want errLoadCancelled, got %v", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("cancel took %v to end the load", elapsed)
	}
}

// A process that is killed but not observed to exit must make the unload fail.
// Returning nil here was how UnloadImageModel could unregister a model whose
// sd-server was still alive.
func TestReapStoppedProcessReportsAnUnreapedExit(t *testing.T) {
	e := NewEngine()
	e.procDone = make(chan struct{})

	err := e.reapStoppedProcess(10 * time.Millisecond)
	if err == nil {
		t.Fatal("an unobserved process exit must be reported as an error")
	}
	if !strings.Contains(err.Error(), "did not exit") {
		t.Fatalf("timeout error = %q, want a did-not-exit explanation", err)
	}
}
