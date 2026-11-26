package melotts

import (
	"context"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"howl-chat/internal/process"
)

// TestMeloServiceHelper is the local MeloTTS stand-in Howl owns in tests.
func TestMeloServiceHelper(t *testing.T) {
	mode := os.Getenv("HOWL_MELO_HELPER")
	if mode == "" {
		t.Skip("helper process only")
	}
	if mode == "exit" {
		t.Fatalf("helper exits immediately by request")
	}
	port := os.Getenv("HOWL_MELO_HELPER_PORT")
	if port == "" {
		t.Fatal("helper requires HOWL_MELO_HELPER_PORT")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	time.Sleep(30 * time.Second)
}

func configureManagedHelper(t *testing.T, port, mode string) {
	t.Helper()
	t.Setenv("HOWL_MELOTTS_COMMAND", os.Args[0])
	t.Setenv("HOWL_MELOTTS_ARGS", "-test.run=TestMeloServiceHelper")
	t.Setenv("HOWL_MELO_HELPER", mode)
	t.Setenv("HOWL_MELO_HELPER_PORT", port)
}

func loopbackURL(t *testing.T) (string, string) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve port: %v", err)
	}
	defer listener.Close()
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener address %T", listener.Addr())
	}
	port := strconv.Itoa(addr.Port)
	return port, "http://127.0.0.1:" + port
}

// TestServiceStartupReadinessAndReap proves Howl starts, owns, registers, and
// reaps the local MeloTTS service.
func TestServiceStartupReadinessAndReap(t *testing.T) {
	registry := process.NewRegistry()
	port, url := loopbackURL(t)
	configureManagedHelper(t, port, "serve")

	service := NewService(registry)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := service.Ensure(ctx, url); err != nil {
		t.Fatalf("ensure service: %v", err)
	}
	if !service.Owned() {
		t.Fatal("a started service must report owned")
	}
	identity, ok := service.Identity()
	if !ok || identity.Role != process.RoleTTSService || identity.PID == 0 || identity.Token == "" {
		t.Fatalf("service identity = %+v ok=%v", identity, ok)
	}
	if got := registry.Children(process.RoleTTSService); len(got) != 1 {
		t.Fatalf("registered services = %d, want 1", len(got))
	}

	if err := service.Terminate(context.Background(), 3*time.Second); err != nil {
		t.Fatalf("terminate service: %v", err)
	}
	if service.Owned() {
		t.Fatal("terminated service still reports owned")
	}
	if got := registry.Children(process.RoleTTSService); len(got) != 0 {
		t.Fatalf("terminated service still registered: %+v", got)
	}
}

// TestServiceFailedStartupIsNotReady proves a launcher that exits is a real
// failure: no readiness, no retained registration.
func TestServiceFailedStartupIsNotReady(t *testing.T) {
	registry := process.NewRegistry()
	port, url := loopbackURL(t)
	configureManagedHelper(t, port, "exit")

	service := NewService(registry)
	err := service.Ensure(context.Background(), url)
	if err == nil {
		t.Fatal("a service that exits before readiness must fail the load")
	}
	if !strings.Contains(err.Error(), "exited before becoming ready") {
		t.Fatalf("startup failure is not explained: %v", err)
	}
	if service.Owned() {
		t.Fatal("failed service reports owned")
	}
	if got := registry.Children(process.RoleTTSService); len(got) != 0 {
		t.Fatalf("failed service left a registration: %+v", got)
	}
}

// TestServiceUnownedEndpointIsNotAdopted proves a loopback listener Howl did
// not start is not silently claimed as an owned service.
func TestServiceUnownedEndpointIsNotAdopted(t *testing.T) {
	registry := process.NewRegistry()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	url := "http://" + listener.Addr().String()

	service := NewService(registry)
	err = service.Ensure(context.Background(), url)
	if err == nil || !strings.Contains(err.Error(), "did not start it") {
		t.Fatalf("unowned endpoint must be refused, got %v", err)
	}
	if service.Owned() {
		t.Fatal("unowned endpoint reported as owned")
	}
	if got := registry.Children(process.RoleTTSService); len(got) != 0 {
		t.Fatalf("unowned endpoint created an ownership record: %+v", got)
	}
}

// TestServiceRefusesNonLoopback proves the managed service cannot be pointed at
// a remote endpoint, so no cloud dependency can be introduced through it.
func TestServiceRefusesNonLoopback(t *testing.T) {
	service := NewService(process.NewRegistry())
	err := service.Ensure(context.Background(), "http://203.0.113.10:8888")
	if err == nil || !strings.Contains(err.Error(), "loopback") {
		t.Fatalf("non-loopback endpoint must be refused, got %v", err)
	}
}
