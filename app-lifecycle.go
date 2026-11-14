package main

import (
	"context"
	stderrors "errors"
	"fmt"
	"os/exec"
	"sync"
	"time"

	apperrors "howl-chat/internal/errors"
	"howl-chat/internal/proc"
)

// llamaShutdownTimeout bounds the wait for a killed llama-server. The process
// is killed first; this only prevents shutdown from blocking forever on a
// driver or child-process failure.
const llamaShutdownTimeout = 10 * time.Second

// shutdownCoordinator makes application shutdown single-entry and idempotent.
// Every caller after the first waits for the same result instead of running
// cleanup a second time.
type shutdownCoordinator struct {
	mu      sync.Mutex
	started bool
	done    chan struct{}
	err     error
}

func (c *shutdownCoordinator) run(ctx context.Context, fn func(context.Context) error) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	c.mu.Lock()
	if c.started {
		done := c.done
		c.mu.Unlock()
		select {
		case <-done:
			return c.err
		case <-ctx.Done():
			return fmt.Errorf("shutdown already in progress: %w", ctx.Err())
		}
	}
	c.started = true
	c.done = make(chan struct{})
	c.mu.Unlock()

	// A panic in one cleanup stage must not strand later shutdown callers. The
	// recover converts it to the same stored error a normal failure returns.
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("shutdown panicked: %v", recovered)
		}
		c.mu.Lock()
		c.err = err
		close(c.done)
		c.mu.Unlock()
	}()

	return fn(ctx)
}

// shutdown runs the one authoritative application cleanup sequence. It is safe
// to call from both the Wails shutdown callback and an explicit Exit action.
func (a *App) shutdown(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	a.shutdownMu.Lock()
	if a.shutdownCoord == nil {
		a.shutdownCoord = &shutdownCoordinator{}
	}
	coord := a.shutdownCoord
	a.shutdownMu.Unlock()
	return coord.run(ctx, a.runShutdown)
}

// runShutdown stops new work first, then saves durable state and tears down
// owned resources in dependency order. Failures are accumulated so one broken
// component cannot prevent the remaining cleanup from running.
func (a *App) runShutdown(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	a.shuttingDown.Store(true)
	stopNativeTestSeam()

	var failures []error
	record := func(stage string, err error) {
		if err == nil {
			return
		}
		failures = append(failures, fmt.Errorf("%s: %w", stage, err))
		a.logError(apperrors.SYSShutdownFail, apperrors.High, "app.shutdown", "app-lifecycle.go:runShutdown", err)
	}
	runStage := func(stage string, fn func() error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				err := fmt.Errorf("%s panicked: %v", stage, recovered)
				failures = append(failures, err)
				a.logError(apperrors.SYSShutdownFail, apperrors.Critical, "app.shutdown", "app-lifecycle.go:runShutdown", err)
			}
		}()
		record(stage, fn())
	}

	// Stop issuing inference work before saving. The HTTP client's request
	// cancellation is separate from the stream-owned context; both are needed
	// for a blocking request and for a stream that has not reached its callback.
	a.cancelActiveStream()
	if a.httpClient != nil {
		a.httpClient.CancelRequest()
	}
	if a.jobManager != nil {
		a.jobManager.CancelAll()
	}

	if a.chatPersistence != nil {
		runStage("final chat save", func() error { return a.performFinalSave() })
	}

	runStage("llama-server stop", func() error { return a.stopLlamaServer(llamaShutdownTimeout) })

	if a.imageEngine != nil {
		runStage("image engine shutdown", func() error { return a.imageEngine.Shutdown() })
	}

	if a.audioManager != nil {
		runStage("audio manager shutdown", func() error {
			err := a.audioManager.Close(ctx)
			if err != nil {
				// Preserve the existing user-visible diagnostic expected by the
				// audio ownership test and by the desktop console.
				fmt.Printf("Audio manager shutdown error: %v\n", err)
			}
			return err
		})
	}

	// Closing the job is the final backstop. All persistent children are
	// assigned at start, so KILL_ON_JOB_CLOSE also covers a crash before this
	// point; this closes the handle deliberately on normal shutdown.
	runStage("process job close", func() error {
		proc.CloseJob()
		return nil
	})

	if len(failures) > 0 {
		return stderrors.Join(failures...)
	}
	return nil
}

// ensureAcceptingWork rejects new user operations once shutdown has begun.
func (a *App) ensureAcceptingWork(operation string) error {
	if a.shuttingDown.Load() {
		return fmt.Errorf("%s is unavailable because Howl is shutting down", operation)
	}
	return nil
}

// beginStream records the cancel function for one generation and returns a
// token. The token prevents an older generation from clearing a newer one.
func (a *App) beginStream(cancel context.CancelFunc) uint64 {
	a.streamMu.Lock()
	defer a.streamMu.Unlock()
	a.streamGen++
	a.streamCancel = cancel
	return a.streamGen
}

func (a *App) endStream(token uint64) {
	a.streamMu.Lock()
	defer a.streamMu.Unlock()
	if a.streamGen == token {
		a.streamCancel = nil
	}
}

func (a *App) cancelActiveStream() {
	a.streamMu.Lock()
	cancel := a.streamCancel
	a.streamGen++
	a.streamCancel = nil
	a.streamMu.Unlock()
	if cancel != nil {
		cancel()
	}
}

// adoptLlamaServer records the process and the single goroutine that owns its
// Wait. No other code may call Wait on this command.
func (a *App) adoptLlamaServer(cmd *exec.Cmd, monitor *processMonitor) {
	a.serverMu.Lock()
	defer a.serverMu.Unlock()
	a.serverProc = cmd
	a.serverMonitor = monitor
}

func (a *App) currentLlamaServer() (*exec.Cmd, *processMonitor) {
	a.serverMu.Lock()
	defer a.serverMu.Unlock()
	return a.serverProc, a.serverMonitor
}

func (a *App) takeLlamaServer() (*exec.Cmd, *processMonitor) {
	a.serverMu.Lock()
	defer a.serverMu.Unlock()
	cmd := a.serverProc
	monitor := a.serverMonitor
	a.serverProc = nil
	a.serverMonitor = nil
	return cmd, monitor
}

// stopLlamaServer kills the owned process tree and waits for the existing
// monitor to observe exit. A timeout means the process did not exit; it is
// reported, never disguised as success.
func (a *App) stopLlamaServer(timeout time.Duration) error {
	cmd, monitor := a.takeLlamaServer()
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	pid := cmd.Process.Pid
	proc.KillTree(pid)
	_ = cmd.Process.Kill()

	if monitor != nil {
		return monitor.wait(timeout)
	}

	// A monitor is always installed by the production spawn path. This fallback
	// exists for partially initialized apps and test-owned commands: one
	// goroutine owns Wait, and shutdown only waits for its signal.
	done := make(chan struct{})
	go func() {
		_, _ = cmd.Process.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("llama-server pid %d did not exit within %s", pid, timeout)
	}
}
