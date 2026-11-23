package host

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// Gate is the application's admission authority. It answers two questions with
// exactly one answer each:
//
//	Stopping() — may new work start?
//	Once(fn)   — who runs the one shutdown sequence, and what did it return?
//
// Every domain asks the gate before starting user-visible work, so shutdown has
// a single place to close the door rather than a flag each subsystem copies.
type Gate struct {
	stopping atomic.Bool

	mu      sync.Mutex
	started bool
	done    chan struct{}
	err     error
}

// NewGate returns an open gate: work is admitted until Stop is called.
func NewGate() *Gate { return &Gate{} }

// Stopping reports whether shutdown has begun.
func (g *Gate) Stopping() bool { return g.stopping.Load() }

// Stop closes admission. It is idempotent and safe to call from any goroutine.
func (g *Gate) Stop() { g.stopping.Store(true) }

// Ensure refuses an operation once shutdown has begun. The returned error names
// the operation, so a refusal is always attributable to the caller that asked.
func (g *Gate) Ensure(operation string) error {
	if g.stopping.Load() {
		return fmt.Errorf("%s is unavailable because Howl is shutting down", operation)
	}
	return nil
}

// Reopen returns the gate to admitting work after a completed shutdown, so a
// restart re-runs the same boot sequence instead of starting a second, parallel
// one. A shutdown still in progress is never interrupted: reopening then would
// race the cleanup that is closing resources.
func (g *Gate) Reopen() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.started && g.done != nil {
		select {
		case <-g.done:
		default:
			return // shutdown still running
		}
	}
	g.stopping.Store(false)
	g.started = false
	g.done = nil
	g.err = nil
}

// Once runs fn exactly once. Concurrent and later callers wait for the same
// result instead of running cleanup a second time. A panic in fn is converted
// to an error and stored, so a waiting caller is never stranded.
func (g *Gate) Once(ctx context.Context, fn func(context.Context) error) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	g.mu.Lock()
	if g.started {
		done := g.done
		g.mu.Unlock()
		select {
		case <-done:
			return g.err
		case <-ctx.Done():
			return fmt.Errorf("shutdown already in progress: %w", ctx.Err())
		}
	}
	g.started = true
	g.done = make(chan struct{})
	g.mu.Unlock()

	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("shutdown panicked: %v", recovered)
		}
		g.mu.Lock()
		g.err = err
		close(g.done)
		g.mu.Unlock()
	}()

	return fn(ctx)
}
