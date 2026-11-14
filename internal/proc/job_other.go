//go:build !windows

// Package proc manages child processes started by the app so they never
// outlive it.
package proc

import "os/exec"

// Attach is a no-op off Windows.
func Attach(cmd *exec.Cmd) error { return nil }

// Start launches the command directly off Windows; there is no job object to
// assign. Callers that need crash-safe ownership on those platforms must use a
// different mechanism.
func Start(cmd *exec.Cmd) error {
	if cmd == nil {
		return nil
	}
	return cmd.Start()
}

// HideWindow is a no-op off Windows.
func HideWindow(cmd *exec.Cmd) {}

// KillTree is a no-op off Windows.
func KillTree(pid int) {}

// CloseJob is a no-op off Windows.
func CloseJob() {}
