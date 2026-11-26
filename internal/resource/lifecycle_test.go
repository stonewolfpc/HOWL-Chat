package resource

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// The three races Milestone 2 names are all gate conflicts: an unload during a
// start, an unload during a reload, and two reloads. These tests model each one
// deterministically — no sleeps that decide the outcome, no real engines.

func TestUnloadDuringStartIsRefused(t *testing.T) {
	gate := NewLifecycleGate()
	start, err := gate.Begin(ResourceImage, "loading")
	if err != nil {
		t.Fatalf("start transition: %v", err)
	}

	_, err = gate.Begin(ResourceImage, "stop")
	if err == nil {
		t.Fatal("an unload during a start must be refused")
	}
	busy, ok := err.(*LifecycleBusyError)
	if !ok {
		t.Fatalf("refusal type = %T, want *LifecycleBusyError", err)
	}
	if busy.Current != "loading" || busy.Kind != ResourceImage {
		t.Fatalf("refusal says kind=%q current=%q, want image/loading", busy.Kind, busy.Current)
	}
	if !strings.Contains(err.Error(), "loading") {
		t.Fatalf("refusal text does not surface the conflicting transition: %q", err.Error())
	}
	start.Release()
}

func TestStopWaitsForStartThenProceeds(t *testing.T) {
	gate := NewLifecycleGate()
	start, err := gate.Begin(ResourceLLM, "loading")
	if err != nil {
		t.Fatalf("start transition: %v", err)
	}

	released := make(chan struct{})
	go func() {
		time.Sleep(20 * time.Millisecond)
		start.Release()
		close(released)
	}()

	stop, err := gate.BeginWithin(ResourceLLM, "stop", 2*time.Second)
	if err != nil {
		t.Fatalf("stop should wait for the start to finish: %v", err)
	}
	<-released
	stop.Release()
}

func TestBusyStartIsRefusedAfterBoundedWait(t *testing.T) {
	gate := NewLifecycleGate()
	first, err := gate.Begin(ResourceLLM, "loading")
	if err != nil {
		t.Fatalf("first start: %v", err)
	}
	defer first.Release()

	began := time.Now()
	_, err = gate.BeginWithin(ResourceLLM, "loading", 30*time.Millisecond)
	if err == nil {
		t.Fatal("a competing start must not silently queue forever")
	}
	if elapsed := time.Since(began); elapsed < 25*time.Millisecond || elapsed > time.Second {
		t.Fatalf("bounded wait took %s, want about 30ms", elapsed)
	}
}

func TestDifferentWorkloadsDoNotBlockEachOther(t *testing.T) {
	gate := NewLifecycleGate()
	chat, err := gate.Begin(ResourceLLM, "loading")
	if err != nil {
		t.Fatalf("chat transition: %v", err)
	}
	defer chat.Release()

	image, err := gate.Begin(ResourceImage, "loading")
	if err != nil {
		t.Fatalf("image must not block on a chat transition: %v", err)
	}
	image.Release()
}

func TestExecutionHoldsNoTransition(t *testing.T) {
	gate := NewLifecycleGate()

	// A stage execution acquires its transition, finishes the start, and
	// releases before generating. Everything after that point is free, which is
	// what lets an unload during long execution be a decision instead of a hang.
	start, err := gate.Begin(ResourceImage, "loading")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	start.Release()

	if _, active := gate.Active(ResourceImage); active {
		t.Fatal("execution must not leave a transition held")
	}

	stop, err := gate.Begin(ResourceImage, "stop")
	if err != nil {
		t.Fatalf("stop during execution must be allowed to decide: %v", err)
	}
	stop.Release()
}

func TestCompetingReloadsSerialize(t *testing.T) {
	gate := NewLifecycleGate()
	reload, err := gate.Begin(ResourceImage, "reloading")
	if err != nil {
		t.Fatalf("first reload: %v", err)
	}

	if _, err := gate.Begin(ResourceImage, "reloading"); err == nil {
		t.Fatal("a second reload must not run alongside the first")
	}
	reload.Release()

	second, err := gate.Begin(ResourceImage, "reloading")
	if err != nil {
		t.Fatalf("reload after release: %v", err)
	}
	second.Release()
}

func TestReleaseIsIdempotentAndNilSafe(t *testing.T) {
	gate := NewLifecycleGate()
	transition, err := gate.Begin(ResourceLLM, "loading")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	transition.Release()
	transition.Release()

	// A release from a stale transition must not clear a newer holder's slot.
	current, err := gate.Begin(ResourceLLM, "stop")
	if err != nil {
		t.Fatalf("second transition: %v", err)
	}
	transition.Release()
	if purpose, active := gate.Active(ResourceLLM); !active || purpose != "stop" {
		t.Fatalf("stale release cleared the live slot: active=%t purpose=%q", active, purpose)
	}
	current.Release()

	var nilTransition *Transition
	nilTransition.Release()
}

func TestConcurrentTransitionsLeaveNoSlotBehind(t *testing.T) {
	gate := NewLifecycleGate()
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transition, err := gate.BeginWithin(ResourceLLM, "loading", 2*time.Second)
			if err != nil {
				t.Errorf("transition %d: %v", i, err)
				return
			}
			transition.Release()
		}()
	}
	wg.Wait()
	if purpose, active := gate.Active(ResourceLLM); active {
		t.Fatalf("a slot survived every release: %q", purpose)
	}
}
