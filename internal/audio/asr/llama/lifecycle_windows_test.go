//go:build windows

package llama

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/sys/windows"

	"howl-chat/internal/audio/types"
	"howl-chat/internal/proc"
)

// These tests use real child processes and assert on their OS-level exit via
// OpenProcess/WaitForSingleObject. A cleanup call returning nil is never treated
// as proof that the process is gone.

const (
	waitTimeout = 0x102      // WAIT_TIMEOUT
	synchronize = 0x00100000 // SYNCHRONIZE access right
)

var (
	fakeOnce sync.Once
	fakeBin  string
	fakeErr  error
)

func TestMain(m *testing.M) {
	code := m.Run()
	if fakeBin != "" {
		_ = os.RemoveAll(filepath.Dir(fakeBin))
	}
	os.Exit(code)
}

// fakeServerSource is a controllable stand-in for llama-server. It honors only
// --port, serves /health and a delayed transcription response, and can spawn a
// long-lived descendant so tree cleanup is verifiable.
const fakeServerSource = `package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func argValue(args []string, name string) string {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == name {
			return args[i+1]
		}
	}
	return ""
}

func main() {
	mode := os.Getenv("FAKEASR_MODE")
	if mode == "idle" {
		time.Sleep(10 * time.Minute)
		return
	}
	if file := os.Getenv("FAKEASR_PID_FILE"); file != "" {
		_ = os.WriteFile(file, []byte(strconv.Itoa(os.Getpid())), 0o644)
	}
	if mode == "noserve" {
		time.Sleep(10 * time.Minute)
		return
	}
	if file := os.Getenv("FAKEASR_CHILD_PID_FILE"); file != "" {
		self, _ := os.Executable()
		child := exec.Command(self)
		child.Env = append(os.Environ(), "FAKEASR_MODE=idle")
		if err := child.Start(); err == nil {
			_ = os.WriteFile(file, []byte(strconv.Itoa(child.Process.Pid)), 0o644)
		}
	}
	port := argValue(os.Args[1:], "--port")
	if port == "" {
		port = "0"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/v1/audio/transcriptions", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"text": "hello from fake asr"})
	})
	ln, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		fmt.Println("listen error:", err)
		os.Exit(2)
	}
	fmt.Println("listening", port)
	_ = http.Serve(ln, mux)
}
`

func fakeServerPath(t *testing.T) string {
	t.Helper()
	fakeOnce.Do(func() {
		dir, err := os.MkdirTemp("", "howl-fake-asr-")
		if err != nil {
			fakeErr = err
			return
		}
		if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(fakeServerSource), 0o644); err != nil {
			fakeErr = err
			return
		}
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module fakeasr\n\ngo 1.21\n"), 0o644); err != nil {
			fakeErr = err
			return
		}
		out := filepath.Join(dir, "fakeasr.exe")
		build := exec.Command("go", "build", "-o", out, ".")
		build.Dir = dir
		if data, err := build.CombinedOutput(); err != nil {
			fakeErr = fmt.Errorf("build fake ASR server: %v\n%s", err, data)
			return
		}
		fakeBin = out
	})
	if fakeErr != nil {
		t.Fatalf("fake ASR server unavailable: %v", fakeErr)
	}
	return fakeBin
}

func writeDummyAssets(t *testing.T) (model, projector string) {
	t.Helper()
	dir := t.TempDir()
	model = filepath.Join(dir, "model.gguf")
	projector = filepath.Join(dir, "mmproj.gguf")
	for _, path := range []string{model, projector} {
		if err := os.WriteFile(path, []byte("stub"), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	return model, projector
}

func testConfig(t *testing.T) types.ASRConfig {
	t.Helper()
	model, projector := writeDummyAssets(t)
	return types.ASRConfig{
		ModelPath:     model,
		ProjectorPath: projector,
		ServerPath:    fakeServerPath(t),
		NumThreads:    1,
		Language:      "auto",
	}
}

func processAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	handle, err := windows.OpenProcess(synchronize, false, uint32(pid))
	if err != nil {
		return false
	}
	defer windows.CloseHandle(handle)
	status, err := windows.WaitForSingleObject(handle, 0)
	if err != nil {
		return false
	}
	return status == waitTimeout
}

func assertAlive(t *testing.T, pid int) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !processAlive(pid) {
		if time.Now().After(deadline) {
			t.Fatalf("pid %d is not running", pid)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func waitGone(t *testing.T, pid int, limit time.Duration) {
	t.Helper()
	deadline := time.Now().Add(limit)
	for processAlive(pid) {
		if time.Now().After(deadline) {
			t.Fatalf("pid %d is still running after %s", pid, limit)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func waitPidFile(t *testing.T, path string, limit time.Duration) int {
	t.Helper()
	deadline := time.Now().Add(limit)
	for {
		if data, err := os.ReadFile(path); err == nil {
			if pid, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil && pid > 0 {
				return pid
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("pid file %s was not written within %s", path, limit)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func waitFile(t *testing.T, path string, limit time.Duration) {
	t.Helper()
	deadline := time.Now().Add(limit)
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("file %s was not written within %s", path, limit)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

// 1. Normal startup followed by release: server and descendant must both exit.
func TestReleaseTerminatesServerAndDescendant(t *testing.T) {
	childFile := filepath.Join(t.TempDir(), "child.pid")
	t.Setenv("FAKEASR_MODE", "serve")
	t.Setenv("FAKEASR_CHILD_PID_FILE", childFile)

	r := New()
	if err := r.Initialize(context.Background(), testConfig(t)); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	serverPid := r.cmd.Process.Pid
	childPid := waitPidFile(t, childFile, 5*time.Second)
	assertAlive(t, serverPid)
	assertAlive(t, childPid)

	if err := r.Release(); err != nil {
		t.Fatalf("Release: %v", err)
	}
	waitGone(t, serverPid, 10*time.Second)
	waitGone(t, childPid, 10*time.Second)
}

// 2. Initialization failure after the child started (context cancelled while
// waiting for health) must not leave the server or its descendant behind.
func TestCancelledStartupTerminatesStartedServer(t *testing.T) {
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "server.pid")
	childFile := filepath.Join(dir, "child.pid")
	t.Setenv("FAKEASR_MODE", "noserve")
	t.Setenv("FAKEASR_PID_FILE", pidFile)
	t.Setenv("FAKEASR_CHILD_PID_FILE", childFile)

	r := New()
	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()
	if err := r.Initialize(ctx, testConfig(t)); err == nil {
		t.Fatal("Initialize succeeded without a health endpoint")
	}

	pid := waitPidFile(t, pidFile, 5*time.Second)
	waitGone(t, pid, 10*time.Second)
}

// 3. Release while a transcription is in flight: the call is serialized behind
// the provider lock, still terminates the server, and repeated Release is safe.
func TestReleaseDuringActiveWorkAndRepeatedRelease(t *testing.T) {
	t.Setenv("FAKEASR_MODE", "serve")

	r := New()
	if err := r.Initialize(context.Background(), testConfig(t)); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	pid := r.cmd.Process.Pid
	audioPath := filepath.Join(t.TempDir(), "audio.wav")
	if err := os.WriteFile(audioPath, []byte("RIFF0000WAVE"), 0o644); err != nil {
		t.Fatalf("write audio: %v", err)
	}

	type outcome struct {
		text string
		err  error
	}
	done := make(chan outcome, 1)
	go func() {
		result, err := r.Transcribe(context.Background(), audioPath, nil)
		if err != nil {
			done <- outcome{err: err}
			return
		}
		done <- outcome{text: result.Text}
	}()
	time.Sleep(250 * time.Millisecond)

	if err := r.Release(); err != nil {
		t.Fatalf("Release during active work: %v", err)
	}
	if err := r.Release(); err != nil {
		t.Fatalf("repeated Release: %v", err)
	}

	select {
	case got := <-done:
		if got.err != nil {
			t.Fatalf("Transcribe: %v", got.err)
		}
		if !strings.Contains(got.text, "hello from fake asr") {
			t.Fatalf("transcript = %q", got.text)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("Transcribe did not return")
	}

	waitGone(t, pid, 10*time.Second)
}

// 4. Abrupt exit of the owning process (no Release, no CloseJob) must still take
// the server and its descendant down through the job object.
func TestAbruptOwnerExitKillsServer(t *testing.T) {
	fake := fakeServerPath(t)
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "server.pid")
	childFile := filepath.Join(dir, "child.pid")

	cmd := exec.Command(os.Args[0], "-test.run=TestAbruptExitHelper")
	cmd.Env = append(os.Environ(),
		"HOWL_D3_ABRUPT=1",
		"HOWL_D3_FAKE="+fake,
		"HOWL_D3_PID_FILE="+pidFile,
		"HOWL_D3_CHILD_PID_FILE="+childFile,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("helper process: %v\n%s", err, out)
	}

	serverPid := waitPidFile(t, pidFile, 5*time.Second)
	childPid := waitPidFile(t, childFile, 5*time.Second)
	t.Cleanup(func() {
		proc.KillTree(serverPid)
		proc.KillTree(childPid)
	})
	waitGone(t, serverPid, 10*time.Second)
	waitGone(t, childPid, 10*time.Second)
}

// 5. CloseJob while the owner is still alive must terminate the attached tree;
// this is the OnShutdown path that runs after the manager's bounded drain.
func TestCloseJobKillsAttachedServer(t *testing.T) {
	fake := fakeServerPath(t)
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "server.pid")
	childFile := filepath.Join(dir, "child.pid")
	marker := filepath.Join(dir, "closed.marker")

	cmd := exec.Command(os.Args[0], "-test.run=TestCloseJobHelper")
	cmd.Env = append(os.Environ(),
		"HOWL_D3_CLOSEJOB=1",
		"HOWL_D3_FAKE="+fake,
		"HOWL_D3_PID_FILE="+pidFile,
		"HOWL_D3_CHILD_PID_FILE="+childFile,
		"HOWL_D3_MARKER="+marker,
	)
	if err := cmd.Start(); err != nil {
		t.Fatalf("helper process: %v", err)
	}
	defer func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}()

	serverPid := waitPidFile(t, pidFile, 5*time.Second)
	childPid := waitPidFile(t, childFile, 5*time.Second)
	assertAlive(t, serverPid)
	assertAlive(t, childPid)

	waitFile(t, marker, 10*time.Second)
	waitGone(t, serverPid, 10*time.Second)
	waitGone(t, childPid, 10*time.Second)
}

// Helper: starts the fake server, attaches it to the job, then exits abruptly
// without any cleanup, proving crash-safety rather than graceful shutdown.
func TestAbruptExitHelper(t *testing.T) {
	if os.Getenv("HOWL_D3_ABRUPT") != "1" {
		t.Skip("subprocess helper")
	}
	cmd := exec.Command(os.Getenv("HOWL_D3_FAKE"))
	cmd.Env = append(os.Environ(),
		"FAKEASR_MODE=serve",
		"FAKEASR_PID_FILE="+os.Getenv("HOWL_D3_PID_FILE"),
		"FAKEASR_CHILD_PID_FILE="+os.Getenv("HOWL_D3_CHILD_PID_FILE"),
	)
	if err := proc.Start(cmd); err != nil {
		t.Fatalf("start fake server under ownership: %v", err)
	}
	time.Sleep(1500 * time.Millisecond)
	os.Exit(0)
}

// Helper: attaches the fake server, calls CloseJob, then stays alive so the
// parent can observe the kill rather than a process-exit side effect.
func TestCloseJobHelper(t *testing.T) {
	if os.Getenv("HOWL_D3_CLOSEJOB") != "1" {
		t.Skip("subprocess helper")
	}
	cmd := exec.Command(os.Getenv("HOWL_D3_FAKE"))
	cmd.Env = append(os.Environ(),
		"FAKEASR_MODE=serve",
		"FAKEASR_PID_FILE="+os.Getenv("HOWL_D3_PID_FILE"),
		"FAKEASR_CHILD_PID_FILE="+os.Getenv("HOWL_D3_CHILD_PID_FILE"),
	)
	if err := proc.Start(cmd); err != nil {
		t.Fatalf("start fake server under ownership: %v", err)
	}
	time.Sleep(1500 * time.Millisecond)
	proc.CloseJob()
	if path := os.Getenv("HOWL_D3_MARKER"); path != "" {
		_ = os.WriteFile(path, []byte("closed"), 0o644)
	}
	time.Sleep(20 * time.Second)
}
