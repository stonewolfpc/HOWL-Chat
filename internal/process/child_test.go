package process

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestChildProcessHelper is the real child used by the ownership tests below.
// It is a normal test function that behaves as a helper only when the parent
// sets HOWL_CHILD_HELPER; otherwise it skips.
func TestChildProcessHelper(t *testing.T) {
	mode := os.Getenv("HOWL_CHILD_HELPER")
	if mode == "" {
		t.Skip("helper process only")
	}
	tickFile := os.Getenv("HOWL_CHILD_HELPER_FILE")
	if tickFile == "" {
		t.Fatalf("helper requires HOWL_CHILD_HELPER_FILE")
	}

	if mode == "parent" && os.Getenv("HOWL_CHILD_HELPER_GRANDCHILD") == "1" {
		cmd := exec.Command(os.Args[0], "-test.run=TestChildProcessHelper")
		cmd.Env = append(os.Environ(), "HOWL_CHILD_HELPER=grandchild")
		if err := cmd.Start(); err != nil {
			t.Fatalf("spawn grandchild: %v", err)
		}
		fmt.Printf("grandchild %d\n", cmd.Process.Pid)
	}

	fmt.Printf("ready %d\n", os.Getpid())
	deadline := time.Now().Add(30 * time.Second)
	if raw := os.Getenv("HOWL_CHILD_HELPER_LIFETIME_MS"); raw != "" {
		if ms, err := strconv.Atoi(raw); err == nil {
			deadline = time.Now().Add(time.Duration(ms) * time.Millisecond)
		}
	}
	for time.Now().Before(deadline) {
		f, err := os.OpenFile(tickFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		if err == nil {
			_, _ = fmt.Fprintf(f, "%d %d\n", os.Getpid(), time.Now().UnixNano())
			_ = f.Close()
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func childEnv(tickFile string, extra ...string) []string {
	env := append(os.Environ(),
		"HOWL_CHILD_HELPER=parent",
		"HOWL_CHILD_HELPER_FILE="+tickFile,
	)
	return append(env, extra...)
}

func helperSpec(registry *Registry, instance, tickFile string, extra ...string) ChildSpec {
	return ChildSpec{
		Role:     RoleTTSChild,
		Instance: instance,
		Path:     os.Args[0],
		Args:     []string{"-test.run=TestChildProcessHelper"},
		Env:      childEnv(tickFile, extra...),
		Capture:  true,
	}
}

// TestChildRegistersBeforeWorkIsRunning proves the child is registered with a
// stable identity while it is still alive, and that the slot disappears only
// after the waiter observed exit.
func TestChildRegistersBeforeWorkIsRunning(t *testing.T) {
	registry := NewRegistry()
	child, err := StartChild(registry, helperSpec(registry, "one", filepath.Join(t.TempDir(), "ticks")))
	if err != nil {
		t.Fatalf("start child: %v", err)
	}
	t.Cleanup(func() { _ = child.Terminate(context.Background(), time.Second) })

	if child.PID() == 0 || child.Identity().Token == "" {
		t.Fatalf("child has no registered identity: %+v", child.Identity())
	}
	current := registry.Children(RoleTTSChild)
	if len(current) != 1 || current[0].Identity != child.Identity() {
		t.Fatalf("registry does not hold the live child: %+v", current)
	}
	if !strings.Contains(child.Output(), "ready") && !waitForOutput(child, "ready", 5*time.Second) {
		t.Fatalf("child did not report readiness: %q", child.Output())
	}
	if err := child.Terminate(context.Background(), 5*time.Second); err != nil {
		t.Fatalf("terminate: %v", err)
	}
	if !child.Exited() {
		t.Fatal("terminate returned before exit was observed")
	}
	if got := registry.Children(RoleTTSChild); len(got) != 0 {
		t.Fatalf("registry kept a settled child: %+v", got)
	}
}

// TestChildCancellationKillsProcessTree proves cancellation terminates the
// child and the process it spawned, not merely the direct child.
func TestChildCancellationKillsProcessTree(t *testing.T) {
	registry := NewRegistry()
	tickFile := filepath.Join(t.TempDir(), "ticks")
	child, err := StartChild(registry, helperSpec(registry, "tree", tickFile, "HOWL_CHILD_HELPER_GRANDCHILD=1"))
	if err != nil {
		t.Fatalf("start child: %v", err)
	}
	if !waitForOutput(child, "grandchild", 5*time.Second) {
		t.Fatalf("helper did not spawn a grandchild: %q", child.Output())
	}
	if !waitForFile(tickFile, 5*time.Second) {
		t.Fatal("helper never wrote its heartbeat")
	}

	ctx, cancel := context.WithCancel(context.Background())
	waitErr := make(chan error, 1)
	go func() { waitErr <- child.Wait(ctx, 2*time.Second) }()
	cancel()

	select {
	case err := <-waitErr:
		if err == nil {
			t.Fatal("cancelled wait must report the cancellation, not clean exit")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("cancelled wait did not return")
	}
	if !child.Exited() {
		t.Fatal("cancellation returned before the child exited")
	}

	// Both the child and its grandchild must stop writing heartbeats.
	size := fileSize(tickFile)
	time.Sleep(300 * time.Millisecond)
	if grew := fileSize(tickFile); grew != size {
		t.Fatalf("a process in the tree survived cancellation: %d -> %d bytes", size, grew)
	}
	if got := registry.Children(RoleTTSChild); len(got) != 0 {
		t.Fatalf("cancelled child still registered: %+v", got)
	}
}

// TestChildTerminationWithoutObserverRetainsOwnership proves an unconfirmed
// exit keeps the registry slot and reports a failure instead of claiming the
// process is gone.
func TestChildTerminationWithoutObserverRetainsOwnership(t *testing.T) {
	registry := NewRegistry()
	cmd := exec.Command(os.Args[0], "-test.run=TestChildProcessHelper")
	cmd.Env = childEnv(filepath.Join(t.TempDir(), "ticks"), "HOWL_CHILD_HELPER_LIFETIME_MS=30000")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child: %v", err)
	}
	handle, err := registry.RegisterChild(RoleTTSChild, "no-observer", cmd, func(time.Duration) error {
		return fmt.Errorf("no exit observed")
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	// The command has no Wait observer by design: this is the failure path a
	// broken observer would create, not a normal one.
	child := &Child{cmd: cmd, handle: handle, done: make(chan struct{}), output: newChildBuffer(1024)}

	err = child.Terminate(context.Background(), 150*time.Millisecond)
	if !errors.Is(err, ErrTerminationUnconfirmed) {
		t.Fatalf("unconfirmed termination must be reported: %v", err)
	}
	if got := registry.Children(RoleTTSChild); len(got) != 1 {
		t.Fatalf("ownership must be retained for a retry, got %+v", got)
	}
	_ = cmd.Process.Kill()
	_, _ = cmd.Process.Wait()
}

// TestChildShutdownWaitsForConfirmedExit proves Terminate does not return
// before the process has actually exited.
func TestChildShutdownWaitsForConfirmedExit(t *testing.T) {
	registry := NewRegistry()
	child, err := StartChild(registry, helperSpec(registry, "slow",
		filepath.Join(t.TempDir(), "ticks"), "HOWL_CHILD_HELPER_LIFETIME_MS=400"))
	if err != nil {
		t.Fatalf("start child: %v", err)
	}
	start := time.Now()
	if err := child.Terminate(context.Background(), 5*time.Second); err != nil {
		t.Fatalf("terminate: %v", err)
	}
	if !child.Exited() {
		t.Fatal("terminate reported success without a confirmed exit")
	}
	// The helper ignores the graceful request (it does not read stdin) and
	// exits on its own; either way, exit must be observed, not assumed.
	if elapsed := time.Since(start); elapsed < 50*time.Millisecond {
		t.Logf("terminate returned in %s", elapsed)
	}
}

// TestStopRoleStopsEveryScopedChild proves several concurrent children of one
// role can be owned and stopped without colliding slots.
func TestStopRoleStopsEveryScopedChild(t *testing.T) {
	registry := NewRegistry()
	tickA := filepath.Join(t.TempDir(), "a")
	tickB := filepath.Join(t.TempDir(), "b")
	childA, err := StartChild(registry, helperSpec(registry, "a", tickA))
	if err != nil {
		t.Fatalf("start a: %v", err)
	}
	childB, err := StartChild(registry, helperSpec(registry, "b", tickB))
	if err != nil {
		t.Fatalf("start b: %v", err)
	}
	if got := registry.Children(RoleTTSChild); len(got) != 2 {
		t.Fatalf("scoped children = %d, want 2", len(got))
	}
	if err := StopRole(registry, RoleTTSChild, 5*time.Second); err != nil {
		t.Fatalf("stop role: %v", err)
	}
	if !childA.Exited() || !childB.Exited() {
		t.Fatal("stop role returned before both children exited")
	}
	if got := registry.Children(RoleTTSChild); len(got) != 0 {
		t.Fatalf("stop role left children registered: %+v", got)
	}
}

func waitForOutput(child *Child, want string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if strings.Contains(child.Output(), want) {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return strings.Contains(child.Output(), want)
}

func waitForFile(path string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if fileSize(path) > 0 {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return fileSize(path) > 0
}

func fileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}
