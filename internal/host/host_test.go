package host

import (
	"context"
	stderrors "errors"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGateRefusesNewWorkAfterStopAndNamesTheOperation(t *testing.T) {
	g := NewGate()
	if err := g.Ensure("send message"); err != nil {
		t.Fatalf("open gate refused work: %v", err)
	}
	g.Stop()
	if !g.Stopping() {
		t.Fatal("Stopping() is false after Stop()")
	}
	err := g.Ensure("send message")
	if err == nil {
		t.Fatal("closed gate admitted work")
	}
	if !strings.Contains(err.Error(), "send message") || !strings.Contains(err.Error(), "shutting down") {
		t.Fatalf("refusal = %q, want the operation and the reason", err)
	}
}

func TestGateOnceRunsTheSequenceOnceAndSharesItsError(t *testing.T) {
	g := NewGate()
	want := stderrors.New("shutdown failed")
	var calls atomic.Int32
	release := make(chan struct{})
	started := make(chan struct{})

	fn := func(context.Context) error {
		calls.Add(1)
		close(started)
		<-release
		return want
	}

	const callers = 8
	errs := make(chan error, callers)
	for i := 0; i < callers; i++ {
		go func() { errs <- g.Once(context.Background(), fn) }()
	}
	<-started
	time.Sleep(20 * time.Millisecond)
	close(release)

	for i := 0; i < callers; i++ {
		if err := <-errs; !stderrors.Is(err, want) {
			t.Fatalf("caller %d error = %v, want the shared result", i, err)
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("shutdown body ran %d times, want 1", got)
	}
}

func TestGateOnceRecoversAPanicForWaiters(t *testing.T) {
	g := NewGate()
	err := g.Once(context.Background(), func(context.Context) error { panic("cleanup panic") })
	if err == nil || !strings.Contains(err.Error(), "shutdown panicked") {
		t.Fatalf("panic result = %v, want a shutdown panic error", err)
	}
	again := g.Once(context.Background(), func(context.Context) error { return nil })
	if again == nil || !strings.Contains(again.Error(), "shutdown panicked") {
		t.Fatalf("later caller = %v, want the stored panic result", again)
	}
}

func TestBackgroundJoinsStartedWork(t *testing.T) {
	g := NewGate()
	b := NewBackground(g)
	var finished atomic.Int32
	for i := 0; i < 3; i++ {
		b.Start("scan", func() {
			time.Sleep(20 * time.Millisecond)
			finished.Add(1)
		})
	}
	if err := b.Join(2 * time.Second); err != nil {
		t.Fatalf("Join: %v", err)
	}
	if got := finished.Load(); got != 3 {
		t.Fatalf("finished work = %d, want 3", got)
	}
}

func TestBackgroundReportsWorkThatOutlivesTheBound(t *testing.T) {
	g := NewGate()
	b := NewBackground(g)
	release := make(chan struct{})
	b.Start("stuck", func() { <-release })

	err := b.Join(50 * time.Millisecond)
	close(release)
	if err == nil || !strings.Contains(err.Error(), "background work still running") {
		t.Fatalf("Join = %v, want a bounded-work failure", err)
	}
}

func TestBackgroundSkipsWorkStartedAfterStop(t *testing.T) {
	g := NewGate()
	b := NewBackground(g)
	g.Stop()
	var ran atomic.Bool
	b.Start("late", func() { ran.Store(true) })
	if err := b.Join(time.Second); err != nil {
		t.Fatalf("Join: %v", err)
	}
	if ran.Load() {
		t.Fatal("work started after Stop ran anyway")
	}
}

func TestStreamsTokenPreventsAnOlderTurnClearingANewerOne(t *testing.T) {
	s := NewStreams()
	var first, second atomic.Int32
	firstToken := s.Begin(func() { first.Add(1) })
	s.Begin(func() { second.Add(1) })

	s.End(firstToken) // the first turn finishes after the second began
	s.Cancel()

	if got := second.Load(); got != 1 {
		t.Fatalf("active stream cancelled %d times, want 1", got)
	}
	if got := first.Load(); got != 0 {
		t.Fatalf("finished stream cancelled %d times, want 0", got)
	}
}

func TestProcessesAdoptTakeAndStopWithoutAProcess(t *testing.T) {
	p := NewProcesses()
	if got := p.Current(); got != nil {
		t.Fatalf("fresh owner holds a process: %v", got)
	}
	if err := p.Stop(time.Second); err != nil {
		t.Fatalf("Stop on an empty owner: %v", err)
	}

	cmd := exec.Command("cmd", "/c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Wait() })
	if _, err := p.Adopt(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if got := p.Current(); got != cmd {
		t.Fatal("Adopt did not record the process")
	}
	taken, waiter := p.Take()
	if taken != cmd || waiter != nil {
		t.Fatal("Take did not return the adopted process and waiter")
	}
	if got := p.Current(); got != nil {
		t.Fatal("Take left the process in the owner")
	}
}

func TestStagesRecordEveryFailureInOrderAndReportSeverity(t *testing.T) {
	type failure struct {
		stage    string
		panicked bool
	}
	var seen []failure
	stages := NewStages(func(stage string, _ error, panicked bool) {
		seen = append(seen, failure{stage: stage, panicked: panicked})
	})

	stages.Run("first", func() error { return stderrors.New("boom") })
	stages.Run("second", func() error { return nil })
	stages.Run("third", func() error { panic("hard") })
	stages.Run("fourth", func() error { return stderrors.New("last") })

	err := stages.Err()
	if err == nil {
		t.Fatal("Err returned nil after failures")
	}
	for _, want := range []string{"first: boom", "third panicked: hard", "fourth: last"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("joined error %q missing %q", err, want)
		}
	}
	if len(seen) != 3 {
		t.Fatalf("reported %d failures, want 3", len(seen))
	}
	if seen[0].stage != "first" || seen[1].stage != "third" || !seen[1].panicked || seen[2].stage != "fourth" {
		t.Fatalf("failure order = %+v, want first/third(panic)/fourth", seen)
	}
}

func TestStagesWithNoFailuresReturnsNil(t *testing.T) {
	stages := NewStages(nil)
	var wg sync.WaitGroup
	wg.Add(1)
	stages.Run("ok", func() error { defer wg.Done(); return nil })
	wg.Wait()
	if err := stages.Err(); err != nil {
		t.Fatalf("Err = %v, want nil", err)
	}
}
