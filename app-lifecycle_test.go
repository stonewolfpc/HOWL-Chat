package main

import (
	"context"
	stderrors "errors"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"howl-chat/internal/host"
)

func TestShutdownCoordinatorRunsOnceUnderConcurrency(t *testing.T) {
	gate := host.NewGate()
	var calls atomic.Int32
	started := make(chan struct{})
	release := make(chan struct{})

	fn := func(context.Context) error {
		calls.Add(1)
		close(started)
		<-release
		return nil
	}

	const callers = 8
	errs := make(chan error, callers)
	for i := 0; i < callers; i++ {
		go func() {
			errs <- gate.Once(context.Background(), fn)
		}()
	}

	<-started
	time.Sleep(50 * time.Millisecond)
	close(release)

	for i := 0; i < callers; i++ {
		if err := <-errs; err != nil {
			t.Fatalf("caller %d: %v", i, err)
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("shutdown body ran %d times, want 1", got)
	}
}

func TestShutdownCoordinatorReturnsTheSameError(t *testing.T) {
	gate := host.NewGate()
	want := stderrors.New("shutdown failed")
	if err := gate.Once(context.Background(), func(context.Context) error { return want }); !stderrors.Is(err, want) {
		t.Fatalf("first error = %v, want %v", err, want)
	}
	if err := gate.Once(context.Background(), func(context.Context) error { return nil }); !stderrors.Is(err, want) {
		t.Fatalf("repeated error = %v, want the first result %v", err, want)
	}
}

func TestShutdownCoordinatorRecoversPanicAndLetsWaitersFinish(t *testing.T) {
	gate := host.NewGate()
	started := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		close(started)
		done <- gate.Once(context.Background(), func(context.Context) error {
			panic("cleanup panic")
		})
	}()
	<-started

	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "shutdown panicked") {
			t.Fatalf("panic result = %v, want a shutdown panic error", err)
		}
	case <-time.After(time.Second):
		t.Fatal("panicking shutdown did not return")
	}

	second := make(chan error, 1)
	go func() { second <- gate.Once(context.Background(), func(context.Context) error { return nil }) }()
	select {
	case err := <-second:
		if err == nil || !strings.Contains(err.Error(), "shutdown panicked") {
			t.Fatalf("second result = %v, want the stored panic result", err)
		}
	case <-time.After(time.Second):
		t.Fatal("second shutdown caller did not observe the completed first call")
	}
}

func TestShutdownWithNilComponentsIsIdempotent(t *testing.T) {
	app := &App{}
	if err := app.shutdown(context.Background()); err != nil {
		t.Fatalf("first shutdown: %v", err)
	}
	if err := app.shutdown(context.Background()); err != nil {
		t.Fatalf("second shutdown: %v", err)
	}
	if !app.level0().Gate.Stopping() {
		t.Fatal("shutdown did not close the work gate")
	}
}

func TestMainRegistersShutdownCallback(t *testing.T) {
	data, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}
	if !strings.Contains(string(data), "OnShutdown: app.OnShutdown") {
		t.Fatal("main.go does not register OnShutdown: app.OnShutdown")
	}
}
