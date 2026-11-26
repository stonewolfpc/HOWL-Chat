package main

import (
	"os/exec"
	"testing"
	"time"

	"howl-chat/internal/resource"
)

// A health check is evidence for the adopted process. Once that process has
// been stopped, the same monitor must not turn the result into a registration.
func TestStaleChatProcessCannotRegisterAfterStop(t *testing.T) {
	const name = "identity-stale-chat"
	resource.GetManager().UnregisterModel(resource.ResourceLLM, name)
	t.Cleanup(func() {
		resource.GetManager().UnregisterModel(resource.ResourceLLM, name)
	})

	app := &App{}
	cmd := exec.Command("cmd", "/c", "ping -n 30 127.0.0.1 >NUL")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	pm := waitMonitor(cmd)
	handle, err := app.level0().Processes.Adopt(cmd, pm.wait)
	if err != nil {
		t.Fatal(err)
	}
	if err := app.level0().Processes.Stop(5 * time.Second); err != nil {
		t.Fatalf("stop adopted process: %v", err)
	}

	if app.registerChatModelAfterHealth(name, "/stale/model.gguf", 100, 0, cmd, handle) {
		t.Fatal("stale health result registered a model for a stopped process")
	}
	for _, m := range resource.GetManager().ModelsByType(resource.ResourceLLM) {
		if m.Name == name {
			t.Fatalf("stale process left registration %+v", m)
		}
	}
}

// The same route must still register when the health-checked process is the
// one this load currently owns.
func TestCurrentChatProcessRegistersAfterHealth(t *testing.T) {
	const name = "identity-current-chat"
	resource.GetManager().UnregisterModel(resource.ResourceLLM, name)
	t.Cleanup(func() {
		resource.GetManager().UnregisterModel(resource.ResourceLLM, name)
	})

	app := &App{}
	cmd := exec.Command("cmd", "/c", "ping", "-n", "30", "127.0.0.1")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	pm := waitMonitor(cmd)
	handle, err := app.level0().Processes.Adopt(cmd, pm.wait)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = app.level0().Processes.Stop(5 * time.Second)
	})

	if !app.registerChatModelAfterHealth(name, "/current/model.gguf", 100, 0, cmd, handle) {
		t.Fatal("current health result did not register the model")
	}
	found := false
	for _, m := range resource.GetManager().ModelsByType(resource.ResourceLLM) {
		if m.Name == name {
			found = true
		}
	}
	if !found {
		t.Fatalf("registration %q missing after current health result", name)
	}
}
