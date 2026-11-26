package host

import (
	"errors"
	"os/exec"
	"testing"
	"time"
)

func TestStopRetainsOwnershipWhenExitIsNotObserved(t *testing.T) {
	p := NewProcesses()
	cmd := exec.Command("cmd", "/c", "ping -n 30 127.0.0.1 >NUL")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})

	calls := 0
	stillAlive := true
	if _, err := p.Adopt(cmd, func(time.Duration) error {
		calls++
		if stillAlive {
			return errors.New("process still running")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if err := p.Stop(10 * time.Millisecond); err == nil {
		t.Fatal("Stop returned nil when the process did not exit")
	}
	if got := p.Current(); got != cmd {
		t.Fatalf("Current() = %p after failed Stop, want retained handle %p", got, cmd)
	}

	// A retry must re-enter the real waiter; it must not take the nil fast path.
	if err := p.Stop(10 * time.Millisecond); err == nil {
		t.Fatal("retry Stop returned nil while the process was still tracked as alive")
	}
	if calls != 2 {
		t.Fatalf("waiter called %d times, want 2 (initial + retry)", calls)
	}
	if got := p.Current(); got != cmd {
		t.Fatalf("Current() = %p after failed retry, want retained handle %p", got, cmd)
	}

	stillAlive = false
	if err := p.Stop(10 * time.Millisecond); err != nil {
		t.Fatalf("Stop after observed exit: %v", err)
	}
	if calls != 3 {
		t.Fatalf("waiter called %d times, want 3", calls)
	}
	if got := p.Current(); got != nil {
		t.Fatalf("Current() = %p after confirmed exit, want nil", got)
	}
}
