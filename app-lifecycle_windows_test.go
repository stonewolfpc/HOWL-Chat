//go:build windows

package main

import (
	"os/exec"
	"testing"
	"time"
)

func TestStopLlamaServerKillsAndWaitsThroughMonitor(t *testing.T) {
	cmd := exec.Command("cmd", "/c", "ping", "-n", "30", "127.0.0.1")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start bounded child: %v", err)
	}
	monitor := waitMonitor(cmd)

	app := &App{}
	app.level0().Processes.Adopt(cmd, monitor.wait)

	start := time.Now()
	if err := app.level0().Processes.Stop(5 * time.Second); err != nil {
		t.Fatalf("stopLlamaServer: %v", err)
	}
	if monitor.isAlive() {
		t.Fatal("child survived stopLlamaServer")
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("stop took %s, want bounded exit", elapsed)
	}
	if cmd := app.level0().Processes.Current(); cmd != nil {
		t.Fatal("server reference was not released")
	}
}
