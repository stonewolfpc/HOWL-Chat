//go:build !windows

package main

import "os/exec"

// applyHidden is a no-op off Windows.
func applyHidden(cmd *exec.Cmd) {}
