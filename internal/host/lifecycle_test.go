package host

import (
	"context"
	stderrors "errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func stagesOf(order *[]string, names ...string) []Stage {
	stages := make([]Stage, 0, len(names))
	for _, name := range names {
		name := name
		stages = append(stages, Stage{Name: name, Run: func(context.Context) error {
			*order = append(*order, name)
			return nil
		}})
	}
	return stages
}

func TestLifecycleBootRunsStagesInOrderAndReachesReadiness(t *testing.T) {
	var order []string
	lifecycle := NewLifecycle(NewGate())

	if err := lifecycle.Boot(context.Background(), stagesOf(&order, "core", "settings", "image")); err != nil {
		t.Fatalf("Boot: %v", err)
	}
	if got := strings.Join(order, ","); got != "core,settings,image" {
		t.Fatalf("stage order = %q, want core,settings,image", got)
	}
	if !lifecycle.Ready() {
		t.Fatal("lifecycle is not ready after every stage completed")
	}
	if err := lifecycle.BootError(); err != nil {
		t.Fatalf("BootError = %v, want nil", err)
	}
	if got := lifecycle.Epoch(); got != 1 {
		t.Fatalf("epoch = %d, want 1", got)
	}
}

func TestLifecycleBootStopsAtTheFirstFailureAndStaysUnready(t *testing.T) {
	var order []string
	boom := stderrors.New("engine binary missing")
	lifecycle := NewLifecycle(NewGate())

	err := lifecycle.Boot(context.Background(), []Stage{
		{Name: "core", Run: func(context.Context) error { order = append(order, "core"); return nil }},
		{Name: "image", Run: func(context.Context) error { order = append(order, "image"); return boom }},
		{Name: "session", Run: func(context.Context) error { order = append(order, "session"); return nil }},
	})

	if err == nil || !stderrors.Is(err, boom) {
		t.Fatalf("Boot error = %v, want the stage failure", err)
	}
	if !strings.Contains(err.Error(), "image") {
		t.Fatalf("Boot error %q does not name the failing stage", err)
	}
	if got := strings.Join(order, ","); got != "core,image" {
		t.Fatalf("stages after the failure ran anyway: %q", got)
	}
	if lifecycle.Ready() {
		t.Fatal("lifecycle reported ready after a failed stage")
	}
	if lifecycle.BootError() == nil {
		t.Fatal("BootError is nil after a failed stage")
	}
}

func TestLifecycleShutdownRunsEveryStageOnceAndReportsEachFailure(t *testing.T) {
	var order []string
	var reported []string
	lifecycle := NewLifecycle(NewGate())

	if err := lifecycle.Boot(context.Background(), stagesOf(&order, "boot")); err != nil {
		t.Fatalf("Boot: %v", err)
	}

	stages := []Stage{
		{Name: "save", Run: func(context.Context) error { order = append(order, "save"); return stderrors.New("disk full") }},
		{Name: "engine", Run: func(context.Context) error { order = append(order, "engine"); return nil }},
		{Name: "child", Run: func(context.Context) error { order = append(order, "child"); panic("hard") }},
	}
	opts := ShutdownOptions{OnFailure: func(stage string, _ error, panicked bool) {
		suffix := ""
		if panicked {
			suffix = "(panicked)"
		}
		reported = append(reported, stage+suffix)
	}}

	err := lifecycle.Shutdown(context.Background(), stages, opts)
	if err == nil {
		t.Fatal("Shutdown returned nil after stage failures")
	}
	for _, want := range []string{"save: disk full", "child panicked: hard"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("shutdown error %q missing %q", err, want)
		}
	}
	if got := strings.Join(order, ","); got != "boot,save,engine,child" {
		t.Fatalf("shutdown order = %q, want save,engine,child after boot", got)
	}
	if got := strings.Join(reported, ","); got != "save,child(panicked)" {
		t.Fatalf("reported failures = %q, want save,child(panicked)", got)
	}

	_, shutdownReport := lifecycle.Report()
	if len(shutdownReport) != 3 || shutdownReport[2].Name != "child" || !shutdownReport[2].Panicked {
		t.Fatalf("shutdown report = %+v, want three stages with child panicking", shutdownReport)
	}

	// A second call shares the first result and does not run cleanup again.
	var reran atomic.Bool
	again := lifecycle.Shutdown(context.Background(), []Stage{{Name: "late", Run: func(context.Context) error {
		reran.Store(true)
		return nil
	}}}, ShutdownOptions{})
	if again == nil || !strings.Contains(again.Error(), "disk full") {
		t.Fatalf("second Shutdown = %v, want the shared result", again)
	}
	if reran.Load() {
		t.Fatal("second Shutdown ran its stages")
	}
}

func TestLifecycleBootAfterShutdownReopensAdmission(t *testing.T) {
	lifecycle := NewLifecycle(NewGate())
	var boots atomic.Int32
	boot := []Stage{{Name: "wire", Run: func(context.Context) error {
		boots.Add(1)
		return lifecycle.Gate().Ensure("wire domains")
	}}}

	if err := lifecycle.Boot(context.Background(), boot); err != nil {
		t.Fatalf("first Boot: %v", err)
	}
	if err := lifecycle.Shutdown(context.Background(), stagesOf(new([]string), "close"), ShutdownOptions{}); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	if lifecycle.Gate().Ensure("after shutdown") == nil {
		t.Fatal("gate admitted work after shutdown")
	}

	// This is the restart path: the same composition runs again through the
	// same authority, and admission is open while it does.
	if err := lifecycle.Boot(context.Background(), boot); err != nil {
		t.Fatalf("reboot: %v", err)
	}
	if got := boots.Load(); got != 2 {
		t.Fatalf("boot stage ran %d times, want 2", got)
	}
	if !lifecycle.Ready() {
		t.Fatal("lifecycle is not ready after the reboot")
	}
	if got := lifecycle.Epoch(); got != 2 {
		t.Fatalf("epoch = %d, want 2", got)
	}
	if err := lifecycle.Gate().Ensure("after reboot"); err != nil {
		t.Fatalf("gate refused work after reboot: %v", err)
	}
}

func TestLifecycleShutdownHonoursACancelledCaller(t *testing.T) {
	lifecycle := NewLifecycle(NewGate())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// The first caller runs the sequence; a cancelled second caller must not be
	// stranded waiting for it.
	started := make(chan struct{})
	release := make(chan struct{})
	go func() {
		close(started)
		_ = lifecycle.Shutdown(context.Background(), []Stage{{Name: "slow", Run: func(context.Context) error {
			<-release
			return nil
		}}}, ShutdownOptions{})
	}()
	<-started
	time.Sleep(20 * time.Millisecond)

	err := lifecycle.Shutdown(ctx, nil, ShutdownOptions{})
	close(release)
	if err == nil || !strings.Contains(err.Error(), "shutdown already in progress") {
		t.Fatalf("cancelled caller error = %v, want the in-progress refusal", err)
	}
}
