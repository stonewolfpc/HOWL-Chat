//go:build !windows

// Package proc manages child processes started by the app so they never
// outlive it.
package proc

import "os/exec"

// Attach is a no-op off Windows.
func Attach(cmd *exec.Cmd) {}

// HideWindow is a no-op off Windows.
func HideWindow(cmd *exec.Cmd) {}

// KillTree is a no-op off Windows.
func KillTree(pid int) {}

// KillTreeByName is a no-op off Windows.
func KillTreeByName(name string) {}

// CloseJob is a no-op off Windows.
func CloseJob() {}
