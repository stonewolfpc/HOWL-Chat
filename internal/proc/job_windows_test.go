//go:build windows

package proc

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

const (
	waitTimeout = 0x102
	synchronize = 0x00100000
)

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

// The child writes a file as its very first action. While it is suspended that
// file must not appear; after Attach and resume it must appear, and closing the
// job must kill it. This proves the ownership handshake, not just a cleanup
// call.
func TestSuspendedChildRunsOnlyAfterOwnership(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "ran.txt")

	cmd := exec.Command("cmd", "/c", "echo ran>"+marker)
	HideWindow(cmd)
	if cmd.SysProcAttr == nil {
		t.Fatal("HideWindow did not configure SysProcAttr")
	}
	cmd.SysProcAttr.CreationFlags |= windows.CREATE_SUSPENDED
	if err := cmd.Start(); err != nil {
		t.Fatalf("start suspended: %v", err)
	}
	pid := cmd.Process.Pid

	if _, err := os.Stat(marker); err == nil {
		t.Fatalf("child executed before ownership: %s exists", marker)
	}
	time.Sleep(400 * time.Millisecond)
	if _, err := os.Stat(marker); err == nil {
		t.Fatalf("child executed while suspended: %s exists", marker)
	}
	if !processAlive(pid) {
		t.Fatal("suspended child exited instead of waiting for resume")
	}

	if err := Attach(cmd); err != nil {
		t.Fatalf("Attach: %v", err)
	}
	if err := resumeProcess(pid); err != nil {
		t.Fatalf("resumeProcess: %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(marker); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("resumed child never ran")
		}
		time.Sleep(25 * time.Millisecond)
	}

	// Job membership is the ownership statement: closing the job (the shutdown
	// path) must terminate the child.
	CloseJob()
	deadline = time.Now().Add(5 * time.Second)
	for processAlive(pid) {
		if time.Now().After(deadline) {
			t.Fatalf("pid %d survived CloseJob", pid)
		}
		time.Sleep(25 * time.Millisecond)
	}
	_, _ = cmd.Process.Wait()
}

func TestStartFailsClosedWhenJobUnavailable(t *testing.T) {
	// CloseJob leaves the singleton without a job, which is the only realistic
	// state in which ownership cannot be established after shutdown.
	CloseJob()
	cmd := exec.Command("cmd", "/c", "exit 0")
	HideWindow(cmd)
	if err := Start(cmd); err == nil {
		t.Fatal("Start succeeded without a job object")
	}
}
