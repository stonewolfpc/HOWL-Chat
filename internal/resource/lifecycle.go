package resource

// lifecycle.go — the cross-workload lifecycle gate.
//
// One card, several engines, and each engine used to have its own start lock.
// Two locks cannot express "do not stop the image engine while it is starting",
// so a chat load could shut sd-server down mid-admission and an image
// reclamation could kill llama-server mid-spawn.
//
// This gate serializes *transitions*, not execution. A transition is the
// bounded span in which a workload starts, stops, recovers or reloads: exactly
// the moments where free memory is read, a process is spawned or killed, and a
// ledger entry is written. Generation, chat inference and stage execution hold a
// lease of their own (imageStartMu, session state) and never hold this gate, so
// a long-running workload cannot block another workload's start or stop for the
// duration of its work.
//
// Ordering rule: a caller may hold a transition for one kind and then request
// one for another (chat start -> image stop), but never the reverse order in the
// same call chain. Every current caller follows chat-then-image, which is what
// keeps the two locks from forming a cycle.

import (
	"fmt"
	"sync"
	"time"
)

// LifecycleGate tracks at most one in-flight transition per resource kind.
type LifecycleGate struct {
	mu     sync.Mutex
	active map[ResourceType]*Transition
}

// Transition is one held lifecycle slot. Release is idempotent, so every exit
// path (including a monitor goroutine that outlives the caller) can release it.
type Transition struct {
	gate    *LifecycleGate
	kind    ResourceType
	purpose string
	once    sync.Once
}

// LifecycleBusyError is a refusal, not a wait result: the requested transition
// conflicts with one that is already running for the same kind.
type LifecycleBusyError struct {
	Kind    ResourceType
	Purpose string
	Current string
	Waited  time.Duration
}

func (e *LifecycleBusyError) Error() string {
	if e == nil {
		return "lifecycle transition refused"
	}
	verb := e.Purpose
	if verb == "" {
		verb = "start"
	}
	if e.Waited > 0 {
		return fmt.Sprintf("cannot %s the %s engine: it is still %s and did not finish within %s",
			verb, e.Kind, e.Current, e.Waited)
	}
	return fmt.Sprintf("cannot %s the %s engine: it is already %s", verb, e.Kind, e.Current)
}

// NewLifecycleGate builds an empty gate.
func NewLifecycleGate() *LifecycleGate {
	return &LifecycleGate{active: make(map[ResourceType]*Transition)}
}

// Begin takes the slot for kind without waiting. It fails when that kind is
// already transitioning.
func (g *LifecycleGate) Begin(kind ResourceType, purpose string) (*Transition, error) {
	return g.BeginWithin(kind, purpose, 0)
}

// BeginWithin waits up to wait for the kind to become free, then takes the slot.
// A zero wait means "fail if busy". It never waits forever: a caller that cannot
// get the slot must surface the refusal instead of hanging a UI action.
func (g *LifecycleGate) BeginWithin(kind ResourceType, purpose string, wait time.Duration) (*Transition, error) {
	if g == nil {
		return nil, fmt.Errorf("lifecycle gate is not configured")
	}
	deadline := time.Now().Add(wait)
	for {
		g.mu.Lock()
		if current, busy := g.active[kind]; busy {
			g.mu.Unlock()
			if wait <= 0 || !time.Now().Before(deadline) {
				return nil, &LifecycleBusyError{Kind: kind, Purpose: purpose, Current: current.purpose, Waited: wait}
			}
			time.Sleep(2 * time.Millisecond)
			continue
		}
		transition := &Transition{gate: g, kind: kind, purpose: purpose}
		g.active[kind] = transition
		g.mu.Unlock()
		return transition, nil
	}
}

// Active reports the purpose of the transition currently held for kind.
func (g *LifecycleGate) Active(kind ResourceType) (string, bool) {
	if g == nil {
		return "", false
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	current, ok := g.active[kind]
	if !ok {
		return "", false
	}
	return current.purpose, true
}

// Release returns the slot. It is safe to call twice and safe on a nil
// transition, because failure paths and monitor goroutines both release.
func (t *Transition) Release() {
	if t == nil || t.gate == nil {
		return
	}
	t.once.Do(func() {
		t.gate.mu.Lock()
		if current, ok := t.gate.active[t.kind]; ok && current == t {
			delete(t.gate.active, t.kind)
		}
		t.gate.mu.Unlock()
	})
}

// Kind reports which workload the transition belongs to.
func (t *Transition) Kind() ResourceType {
	if t == nil {
		return ""
	}
	return t.kind
}

// Purpose reports what the transition is doing, for logs and refusals.
func (t *Transition) Purpose() string {
	if t == nil {
		return ""
	}
	return t.purpose
}
