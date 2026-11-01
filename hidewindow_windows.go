//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

// applyHidden starts child processes without flashing a console window.
// Without this, every probe (sd-server --version, llama-server --version)
// pops up a console window on a GUI build.
func applyHidden(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true
}
