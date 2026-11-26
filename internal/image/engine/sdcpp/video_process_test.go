package sdcpp

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"howl-chat/internal/process"
	"howl-chat/internal/resource"
)

// recordingRegistry observes the engine's stop route without changing the
// registry's behavior.
type recordingRegistry struct {
	*process.Registry
	stopCalls int
	stopErr   error
}

func (r *recordingRegistry) Stop(role process.Role, timeout time.Duration) error {
	r.stopCalls++
	if r.stopErr != nil {
		return r.stopErr
	}
	return r.Registry.Stop(role, timeout)
}

func TestStopServerUsesOneRegisteredWaiter(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("spawns cmd.exe")
	}
	e := NewEngine()
	registry := &recordingRegistry{Registry: process.NewRegistry()}
	e.processRegistry = registry

	cmd := exec.Command("cmd", "/c", "ping -n 30 127.0.0.1 >NUL")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	waiterCalls := 0
	if _, err := registry.Register(process.RoleImageServer, cmd, func(time.Duration) error {
		waiterCalls++
		_ = cmd.Wait()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	e.setServerProc(cmd.Process)

	if err := e.stopServer(); err != nil {
		t.Fatalf("stopServer: %v", err)
	}
	if registry.stopCalls != 1 {
		t.Fatalf("registry Stop calls = %d, want 1", registry.stopCalls)
	}
	if waiterCalls != 1 {
		t.Fatalf("registered waiter calls = %d, want exactly 1", waiterCalls)
	}
	if got := e.getServerProc(); got != nil {
		t.Fatalf("serverProc = %p after confirmed stop, want nil", got)
	}
}

// A video run stops sd-server before starting sd-cli. That process exit must
// invalidate the image registration through the app's exit hook; leaving the
// registration behind would claim a live image process that is gone.
func TestVideoLeavesNoImageRegistration(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("spawns cmd.exe")
	}
	const name = "video-stops-image-registration"
	resource.GetManager().UnregisterModel(resource.ResourceImage, name)
	t.Cleanup(func() {
		resource.GetManager().UnregisterModel(resource.ResourceImage, name)
	})
	resource.GetManager().RegisterModel(
		resource.ResourceImage, name, "/image/model.safetensors", 100, 0, resource.PriorityImage, true,
	)

	e := NewEngine()
	registry := &recordingRegistry{Registry: process.NewRegistry()}
	e.processRegistry = registry
	e.SetShutdownHook(func() {
		resource.GetManager().UnregisterModel(resource.ResourceImage, name)
	})

	cmd := exec.Command("cmd", "/c", "ping -n 30 127.0.0.1 >NUL")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	if _, err := registry.Register(process.RoleImageServer, cmd, func(time.Duration) error {
		_ = cmd.Wait()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	e.setServerProc(cmd.Process)

	if err := e.stopServer(); err != nil {
		t.Fatalf("stopServer: %v", err)
	}
	for _, m := range resource.GetManager().ModelsByType(resource.ResourceImage) {
		if m.Name == name {
			t.Fatalf("image registration %q outlived its stopped process", name)
		}
	}
}

// sd-cli is adopted as a Howl child for the whole run and removed after its
// recorded waiter observes exit.
func TestVideoProcessRegisteredAndReleased(t *testing.T) {
	e := NewEngine()
	t.Setenv("HOWL_VIDEO_CLI_HELPER", "1")

	done := make(chan error, 1)
	go func() {
		done <- e.runVideoCLI(os.Args[0], []string{"-test.run=TestVideoCLIHelperProcess", "--"}, filepath.Join(t.TempDir(), "out.mp4"), "test-owner")
	}()

	var handle *process.Handle
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		handle = e.processes().Current(process.RoleVideoCLI)
		if handle != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if handle == nil {
		t.Fatal("sd-cli was not registered while it was running")
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("runVideoCLI: %v", err)
		}
	case <-time.After(8 * time.Second):
		t.Fatal("sd-cli did not finish")
	}
	if got := e.processes().Current(process.RoleVideoCLI); got != nil {
		t.Fatalf("sd-cli registration survived exit: %v", got.Identity)
	}
}

// Cancellation also releases the registry slot exactly once: the kill makes
// the recorded waiter return, and the deferred Remove follows it.
func TestVideoProcessReleasedOnCancel(t *testing.T) {
	e := NewEngine()
	t.Setenv("HOWL_VIDEO_CLI_HELPER", "1")

	done := make(chan error, 1)
	go func() {
		done <- e.runVideoCLI(os.Args[0], []string{"-test.run=TestVideoCLIHelperProcess", "--"}, filepath.Join(t.TempDir(), "out.mp4"), "test-owner")
	}()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) && e.processes().Current(process.RoleVideoCLI) == nil {
		time.Sleep(10 * time.Millisecond)
	}
	if e.processes().Current(process.RoleVideoCLI) == nil {
		t.Fatal("sd-cli was not registered before cancel")
	}
	if !e.cancelVideo() {
		t.Fatal("cancelVideo did not find the running sd-cli")
	}

	select {
	case <-done:
	case <-time.After(8 * time.Second):
		t.Fatal("sd-cli did not finish after cancellation")
	}
	if got := e.processes().Current(process.RoleVideoCLI); got != nil {
		t.Fatalf("sd-cli registration survived cancellation: %v", got.Identity)
	}
}

// TestVideoCLIHelperProcess is re-executed as the sd-cli child above. It holds
// the process open long enough for the parent to observe its registration.
func TestVideoCLIHelperProcess(t *testing.T) {
	if os.Getenv("HOWL_VIDEO_CLI_HELPER") != "1" {
		return
	}
	time.Sleep(time.Second)
}

func TestStopServerRetainsOwnershipWhenRegisteredStopFails(t *testing.T) {
	e := NewEngine()
	registry := &recordingRegistry{
		Registry: process.NewRegistry(),
		stopErr:  errors.New("exit not observed"),
	}
	e.processRegistry = registry

	cmd := exec.Command("cmd", "/c", "ping -n 30 127.0.0.1 >NUL")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	if _, err := registry.Register(process.RoleImageServer, cmd, nil); err != nil {
		t.Fatal(err)
	}
	e.setServerProc(cmd.Process)

	if err := e.stopServer(); err == nil {
		t.Fatal("stopServer returned nil when the registry did not observe exit")
	}
	if got := e.getServerProc(); got != cmd.Process {
		t.Fatalf("serverProc = %p after failed stop, want retained %p", got, cmd.Process)
	}
	if registry.stopCalls != 1 {
		t.Fatalf("registry Stop calls = %d, want 1", registry.stopCalls)
	}
}
