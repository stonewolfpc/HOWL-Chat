package main

// app-lifecycle-gate.go — App's access to the cross-workload lifecycle gate.
//
// The gate lives in internal/resource because it is a resource-ownership
// concern, not an App one. These helpers only decide how long a caller waits
// before the refusal surfaces. Starts wait long enough for an in-flight stop to
// finish; stops refuse quickly so a user action is never left hanging behind a
// spawn that may take a minute.
//
// This is deliberately separate from app-lifecycle.go, which owns the Level 0
// authority (boot, shutdown, accepting work). This file adds no authority.

import (
	"time"

	"howl-chat/internal/resource"
)

// Vars, not consts, so tests can shorten a wait instead of sleeping for real.
var (
	lifecycleStartWait = 30 * time.Second
	lifecycleStopWait  = 5 * time.Second
)

// lifecycleTransitions returns the process-wide gate, creating it on first use.
// Tests build App values directly, so the zero value must be usable.
func (a *App) lifecycleTransitions() *resource.LifecycleGate {
	a.lifecycleOnce.Do(func() {
		if a.lifecycleGate == nil {
			a.lifecycleGate = resource.NewLifecycleGate()
		}
	})
	return a.lifecycleGate
}

// beginLifecycleTransition takes the slot for one workload's transition. The
// caller must release it; every call site defers the release.
func (a *App) beginLifecycleTransition(kind resource.ResourceType, purpose string, wait time.Duration) (*resource.Transition, error) {
	return a.lifecycleTransitions().BeginWithin(kind, purpose, wait)
}
