package process

import (
	"errors"
	"os/exec"
	"testing"
	"time"
)

func TestRegisterAssignsIdentityAndRejectsSecondChildInOneRole(t *testing.T) {
	r := NewRegistry()
	cmd := exec.Command("cmd", "/c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})

	handle, err := r.Register(RoleChat, cmd, nil)
	if err != nil {
		t.Fatal(err)
	}
	if handle.Identity.PID != cmd.Process.Pid || handle.Identity.Role != RoleChat || handle.Identity.Token == "" {
		t.Fatalf("identity = %+v, want role, pid and non-empty token", handle.Identity)
	}
	if got := r.Current(RoleChat); got != handle {
		t.Fatalf("Current(chat) = %p, want %p", got, handle)
	}
	if _, err := r.Register(RoleChat, cmd, nil); err == nil {
		t.Fatal("second registration in one role succeeded; role must have one owner")
	}
}

func TestStartTokensDistinguishSuccessiveRegistrations(t *testing.T) {
	r := NewRegistry()
	cmd := exec.Command("cmd", "/c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})

	first, err := r.Register(RoleChat, cmd, nil)
	if err != nil {
		t.Fatal(err)
	}
	first.Release()
	second, err := r.Register(RoleChat, cmd, nil)
	if err != nil {
		t.Fatal(err)
	}
	if first.Identity == second.Identity {
		t.Fatalf("successive registrations reused identity %+v", first.Identity)
	}
}

func TestOldHandleReleaseDoesNotReleaseAReplacement(t *testing.T) {
	r := NewRegistry()
	cmd := exec.Command("cmd", "/c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})

	first, err := r.Register(RoleChat, cmd, nil)
	if err != nil {
		t.Fatal(err)
	}
	first.Release()
	second, err := r.Register(RoleChat, cmd, nil)
	if err != nil {
		t.Fatal(err)
	}
	first.Release()
	if got := r.Lookup(second.Identity); got != second {
		t.Fatal("an old handle released the replacement process")
	}
	second.Release()
	if got := r.Current(RoleChat); got != nil {
		t.Fatalf("Current(chat) = %p after matching release, want nil", got)
	}
}

func TestStopKeepsOwnershipWhenExitIsNotObserved(t *testing.T) {
	r := NewRegistry()
	cmd := exec.Command("cmd", "/c", "ping -n 30 127.0.0.1 >NUL")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})

	handle, err := r.Register(RoleChat, cmd, func(time.Duration) error {
		return errors.New("exit not observed")
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := handle.Stop(10 * time.Millisecond); err == nil {
		t.Fatal("Stop returned nil when the exit was not observed")
	}
	if got := r.Lookup(handle.Identity); got != handle {
		t.Fatal("failed Stop released the child instead of retaining ownership")
	}
}
