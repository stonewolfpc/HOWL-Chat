//go:build windows

package main

import (
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"golang.org/x/sys/windows"

	"howl-chat/internal/proc"
)

const d3WaitTimeout = 0x102

func d3ProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	handle, err := windows.OpenProcess(0x00100000, false, uint32(pid))
	if err != nil {
		return false
	}
	defer windows.CloseHandle(handle)
	status, err := windows.WaitForSingleObject(handle, 0)
	if err != nil {
		return false
	}
	return status == d3WaitTimeout
}

// D3 chat regression: the chat spawn path uses exec.Command + proc.Attach, and
// the corrected Attach must now really place the server in the job. Gated by an
// environment variable so the normal suite does not load a model.
func TestD3ChatServerSpawnOwnership(t *testing.T) {
	model := os.Getenv("HOWL_D3_CHAT_MODEL")
	if model == "" {
		t.Skip("HOWL_D3_CHAT_MODEL not set")
	}
	if _, err := os.Stat(model); err != nil {
		t.Skipf("chat model unavailable: %v", err)
	}

	start := func() *exec.Cmd {
		cmd, err := spawnLlamaServerWithModel(model, "", nil, io.Discard, true, llamaCapsForTests())
		if err != nil {
			t.Fatalf("spawn chat server: %v", err)
		}
		return cmd
	}
	waitAlive := func(cmd *exec.Cmd) {
		deadline := time.Now().Add(10 * time.Second)
		for !d3ProcessAlive(cmd.Process.Pid) {
			if time.Now().After(deadline) {
				t.Fatalf("chat server pid %d did not stay alive", cmd.Process.Pid)
			}
			if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
				t.Fatalf("chat server exited early: %v", cmd.ProcessState)
			}
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(3 * time.Second)
		if !d3ProcessAlive(cmd.Process.Pid) {
			t.Fatalf("chat server pid %d exited while loading", cmd.Process.Pid)
		}
	}
	waitGone := func(cmd *exec.Cmd, limit time.Duration) {
		deadline := time.Now().Add(limit)
		for d3ProcessAlive(cmd.Process.Pid) {
			if time.Now().After(deadline) {
				t.Fatalf("chat server pid %d survived shutdown", cmd.Process.Pid)
			}
			time.Sleep(100 * time.Millisecond)
		}
		_, _ = cmd.Process.Wait()
	}

	// 1. Normal shutdown path used by OnShutdown.
	first := start()
	waitAlive(first)
	proc.KillTree(first.Process.Pid)
	waitGone(first, 15*time.Second)

	// 2. Ownership: closing the job must terminate a server started through the
	// chat path, which only works if Attach actually assigned it.
	second := start()
	waitAlive(second)
	proc.CloseJob()
	waitGone(second, 20*time.Second)
}
