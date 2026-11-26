package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestUnloadModelFailedStopRetainsTheProcessAndReportsNoSuccess(t *testing.T) {
	app := &App{}
	cmd := exec.Command("cmd", "/c", "ping -n 30 127.0.0.1 >NUL")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})
	calls := 0
	app.level0().Processes.Adopt(cmd, func(time.Duration) error {
		calls++
		return errors.New("llama-server still running")
	})

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	unloadErr := app.UnloadModel()
	_ = w.Close()
	os.Stdout = oldStdout
	out, readErr := io.ReadAll(r)
	_ = r.Close()
	if readErr != nil {
		t.Fatal(readErr)
	}

	if unloadErr == nil {
		t.Fatal("UnloadModel returned nil after a failed process stop")
	}
	if !strings.Contains(unloadErr.Error(), "still running") {
		t.Fatalf("UnloadModel error = %q, want the stop failure", unloadErr)
	}
	if strings.Contains(string(out), "memory freed") {
		t.Fatalf("UnloadModel reported success while the process was still alive: %q", out)
	}
	if got := app.level0().Processes.Current(); got != cmd {
		t.Fatalf("Current() = %p after failed unload, want retained handle %p", got, cmd)
	}

	// The retry must still reach the real waiter, not succeed on a discarded handle.
	if retryErr := app.UnloadModel(); retryErr == nil {
		t.Fatal("retry UnloadModel returned nil while the process was still alive")
	}
	if calls != 2 {
		t.Fatalf("waiter called %d times, want 2 (initial unload + retry)", calls)
	}
	if got := app.level0().Processes.Current(); got != cmd {
		t.Fatalf("Current() = %p after failed retry, want retained handle %p", got, cmd)
	}
}
