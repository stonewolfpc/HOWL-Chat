//go:build !windows

package llama

import "os/exec"

func applyHidden(cmd *exec.Cmd) {}
