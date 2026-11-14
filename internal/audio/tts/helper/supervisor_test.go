package helper

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSupervisorUsesOneWaitObserver(t *testing.T) {
	h := newFakeHelper()
	defer h.p.Kill()
	path := filepath.Join(t.TempDir(), "helper.exe")
	if err := os.WriteFile(path, []byte("fake"), 0o644); err != nil {
		t.Fatal(err)
	}
	sup := NewSupervisor(SupervisorOptions{
		Discoverer: Discoverer{LookupEnv: func(string) (string, bool) { return path, true }},
		Factory: func(context.Context, LaunchSpec) (Process, error) {
			return h.p, nil
		},
	})
	if _, err := sup.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for h.p.waitCount.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if got := h.p.waitCount.Load(); got != 1 {
		t.Fatalf("Wait calls after start = %d, want 1", got)
	}
	if err := sup.Terminate(context.Background(), time.Second); err != nil {
		t.Fatal(err)
	}
	if got := h.p.waitCount.Load(); got != 1 {
		t.Fatalf("Wait calls after termination = %d, want 1", got)
	}
}

func TestSupervisorObservesCrash(t *testing.T) {
	h := newFakeHelper()
	path := filepath.Join(t.TempDir(), "helper.exe")
	if err := os.WriteFile(path, []byte("fake"), 0o644); err != nil {
		t.Fatal(err)
	}
	sup := NewSupervisor(SupervisorOptions{
		Discoverer: Discoverer{LookupEnv: func(string) (string, bool) { return path, true }},
		Factory: func(context.Context, LaunchSpec) (Process, error) {
			return h.p, nil
		},
	})
	if _, err := sup.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	_ = h.p.Kill()
	deadline := time.Now().Add(time.Second)
	for sup.Status().State != "crashed" && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	status := sup.Status()
	if status.State != "crashed" || status.ProcessAlive || status.ProcessState != "exited" {
		t.Fatalf("status after crash = %+v", status)
	}
	if got := h.p.waitCount.Load(); got != 1 {
		t.Fatalf("Wait calls after crash = %d, want 1", got)
	}
}

func TestSupervisorFailedStartCreatesNoWaiter(t *testing.T) {
	path := filepath.Join(t.TempDir(), "helper.exe")
	if err := os.WriteFile(path, []byte("fake"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := newFakeHelper()
	defer h.p.Kill()
	sup := NewSupervisor(SupervisorOptions{
		Discoverer: Discoverer{LookupEnv: func(string) (string, bool) { return path, true }},
		Factory: func(context.Context, LaunchSpec) (Process, error) {
			// A failed factory must transfer no process and therefore create no
			// supervisor waiter.
			return h.p, errors.New("start failed")
		},
	})
	if _, err := sup.Start(context.Background()); err == nil {
		t.Fatal("Start succeeded, want failure")
	}
	if got := h.p.waitCount.Load(); got != 0 {
		t.Fatalf("waiter count = %d, want 0", got)
	}
	status := sup.Status()
	if status.State != "unavailable" || status.ProcessAlive {
		t.Fatalf("status = %+v, want unavailable and not alive", status)
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("HOWL_HELPER_TEST_PROCESS") != "1" {
		return
	}
	marker := os.Getenv("HOWL_HELPER_TEST_MARKER")
	if marker != "" {
		_ = os.WriteFile(marker, []byte("started"), 0o644)
	}
	os.Exit(0)
}

func TestStartProcessUsesRealStartupPath(t *testing.T) {
	if os.Getenv("HOWL_HELPER_TEST_PROCESS") == "1" {
		return
	}
	marker := filepath.Join(t.TempDir(), "marker")
	p, err := StartProcess(context.Background(), LaunchSpec{
		Path: os.Args[0],
		Args: []string{"-test.run=^TestHelperProcess$"},
		Env:  append(os.Environ(), "HOWL_HELPER_TEST_PROCESS=1", "HOWL_HELPER_TEST_MARKER="+marker),
	})
	if err != nil {
		t.Fatalf("StartProcess: %v", err)
	}
	if err := p.Wait(); err != nil {
		t.Fatalf("Wait: %v", err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("helper process did not run after ownership startup: %v", err)
	}
	_ = p.Close()
}
