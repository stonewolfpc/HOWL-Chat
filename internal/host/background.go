package host

import (
	"fmt"
	"sync"
	"time"
)

// Background owns work that outlives the call which starts it — directory
// scans, model discovery, anything the window does not wait for.
//
// It joins that work before durable state is written, so a scan is never
// mid-write when a store is finalised. A panic is reported and contained: a
// goroutine nobody is watching must not take the process down.
type Background struct {
	gate *Gate
	wg   sync.WaitGroup
}

// NewBackground binds background work to the gate that decides whether new
// work may still start.
func NewBackground(gate *Gate) *Background { return &Background{gate: gate} }

// Start runs fn in a goroutine that Join will wait for.
func (b *Background) Start(name string, fn func()) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		defer func() {
			if recovered := recover(); recovered != nil {
				fmt.Printf("[background] %s panicked: %v\n", name, recovered)
			}
		}()
		if b.gate != nil && b.gate.Stopping() {
			return
		}
		fn()
	}()
}

// Join waits for started work with a bound. Work that refuses to finish is
// reported rather than waited on forever; the process must still exit.
func (b *Background) Join(timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("background work still running after %s", timeout)
	}
}
