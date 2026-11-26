package sdcpp

// The goroutine started by watchProcess owns the single Wait() for the
// sd-server child. Everything that needs to know whether the process is alive
// reads that result instead of asking HTTP, because readiness and liveness are
// different questions and sd-server answers only the first one.
//
// This existed because of krea2: sd-server printed "new_sd_ctx_t failed" within
// seconds, and the app went on polling an endpoint that could never come up for
// the whole load timeout, then reported a version problem instead of the actual
// exit.

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"howl-chat/internal/process"
)

// errProcessExited marks a load that ended because the child died, so a caller
// never describes it as a timeout that did not happen.
var errProcessExited = errors.New("engine process exited while loading")

// errLoadCancelled reports a load stopped on request, so a cancellation is not
// explained to the user as a failure.
var errLoadCancelled = errors.New("model load cancelled")

// setServerProc and getServerProc guard the child handle. It is written by the
// load, killed by cancellation, and read by shutdown and liveness, and those can
// run on different goroutines at the same time.
func (e *Engine) setServerProc(p *os.Process) {
	e.procMu.Lock()
	e.serverProc = p
	e.procMu.Unlock()
}

func (e *Engine) getServerProc() *os.Process {
	e.procMu.Lock()
	defer e.procMu.Unlock()
	return e.serverProc
}

// clearServerProc releases the handle only if it is still the process the
// caller stopped. A replacement's handle must never be cleared by an older
// stop that finished late.
func (e *Engine) clearServerProc(p *os.Process) {
	e.procMu.Lock()
	defer e.procMu.Unlock()
	if e.serverProc == p {
		e.serverProc = nil
	}
}

// beginLoad marks a load in progress, records which job owns it, and clears any
// abort left by the previous one so a stale cancellation cannot kill the next
// model's process.
func (e *Engine) beginLoad(jobID string) {
	e.loadMu.Lock()
	e.loading = true
	e.loadOwner = jobID
	e.aborting = false
	e.loadMu.Unlock()
}

// endLoad marks the load finished, however it ended.
func (e *Engine) endLoad() {
	e.loadMu.Lock()
	e.loading = false
	e.loadMu.Unlock()
}

// loadInProgress reports whether a load is running right now.
func (e *Engine) loadInProgress() bool {
	e.loadMu.Lock()
	defer e.loadMu.Unlock()
	return e.loading
}

// loadAborted reports whether the running load was cancelled.
func (e *Engine) loadAborted() bool {
	e.loadMu.Lock()
	defer e.loadMu.Unlock()
	return e.aborting
}

// abortLoad stops a load owned by jobID, and reports whether it owned one. A
// child that is still reading weights only stops when the child stops: killing it
// is what makes waitReadyOrExit return at once, and LoadModel's error path then
// clears the half-built state.
//
// A cancel from any other job — or from no job — is a no-op, and returning false
// lets the caller route the cancel to whatever the job does own. One job must
// never abort another job's load, which is what a plain "is a load running"
// check used to do.
func (e *Engine) abortLoad(jobID string) bool {
	if jobID == "" {
		return false
	}
	e.loadMu.Lock()
	if !e.loading || e.loadOwner != jobID {
		e.loadMu.Unlock()
		return false
	}
	e.aborting = true
	e.loadMu.Unlock()
	if p := e.getServerProc(); p != nil {
		_ = p.Kill()
	}
	return true
}

// startProcessReaper reads the child's exit status. Exactly one goroutine may
// call Wait(), so this is the only caller in the package.
//
// The done channel is registered before the goroutine starts, not inside it: an
// observer that read the field first would otherwise get nil, park on a select
// over a never-closed channel, and miss the exit entirely — which is how this
// was first written, and what TestWaitReadyOrExitReturnsTheExitWhenTheChildDies
// caught by taking the full thirty-second timeout.
func (e *Engine) startProcessReaper(cmd *exec.Cmd) <-chan struct{} {
	return e.startProcessReaperWithHandle(cmd, nil)
}

// startProcessReaperWithHandle starts the one goroutine that calls Wait. When a
// registry handle is supplied, exit releases that exact identity, so a later
// child that reused the PID cannot be released by this observer.
func (e *Engine) startProcessReaperWithHandle(cmd *exec.Cmd, handle *process.Handle) <-chan struct{} {
	done := make(chan struct{})
	e.procMu.Lock()
	e.procDone = done
	e.procExit = nil
	e.procMu.Unlock()

	go func() {
		err := cmd.Wait()
		e.procMu.Lock()
		e.procExit = err
		e.procMu.Unlock()
		// No live process owns the placement or footprint anymore. Clearing here
		// covers a crash and a direct kill, not just the orderly Shutdown path.
		e.resetPlacement()
		e.resetFootprints()
		if handle != nil {
			handle.Release()
		}
		close(done)
	}()
	return done
}

// adoptServerProcess records the identity and single Wait owner for sd-server.
// It is called before readiness, so a health check can never refer to a child
// the engine no longer owns.
func (e *Engine) adoptServerProcess(cmd *exec.Cmd, done <-chan struct{}) (*process.Handle, error) {
	return e.processes().Register(process.RoleImageServer, cmd, func(timeout time.Duration) error {
		select {
		case <-done:
			return nil
		case <-time.After(timeout):
			return fmt.Errorf("sd-server pid %d still running", cmd.Process.Pid)
		}
	})
}

// processDone returns the channel closed when the watched child exits, or nil
// when no child has been started. Receive on nil deliberately blocks, so a nil
// result parks a select rather than reporting a phantom exit.
func (e *Engine) processDone() <-chan struct{} {
	e.procMu.Lock()
	defer e.procMu.Unlock()
	return e.procDone
}

// exitError describes how the child ended, for a message the user can act on.
func (e *Engine) exitError() error {
	e.procMu.Lock()
	defer e.procMu.Unlock()
	if e.procExit != nil {
		return e.procExit
	}
	// Wait() returned no error but the process is gone: a clean exit that still
	// left us without a server is a failure of the load, not of the binary.
	return fmt.Errorf("process exited cleanly before it was ready")
}

// waitReadyOrExit waits for sd-server to answer, and stops the moment the child
// exits. Polling readiness alone meant a dead process was waited on for the
// entire load timeout before the error even mentioned the real cause.
func (e *Engine) waitReadyOrExit(timeout time.Duration) error {
	if e.httpClient == nil {
		return fmt.Errorf("engine HTTP client not initialized")
	}
	done := e.processDone()
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	tick := time.NewTicker(250 * time.Millisecond)
	defer tick.Stop()

	for {
		// A cancellation can land after LoadModel's pre-start check, while the child
		// is already reading weights. Checking here closes that window: the load
		// stops within one tick instead of loading a whole model for a job the user
		// already abandoned.
		if e.loadAborted() {
			return errLoadCancelled
		}
		if err := e.httpClient.probe(); err == nil {
			return nil
		}
		select {
		case <-done:
			// Killing the child is how a cancellation stops a load, so the exit it
			// produces is a cancellation, not a crash to explain.
			if e.loadAborted() {
				return errLoadCancelled
			}
			return fmt.Errorf("%w: %v", errProcessExited, e.exitError())
		case <-deadline.C:
			if e.loadAborted() {
				return errLoadCancelled
			}
			return fmt.Errorf("sd-server not ready after %v", timeout.Round(time.Second))
		case <-tick.C:
		}
	}
}

// reapStoppedProcess waits for the observer to finish reaping a child that is
// being killed. Only the observer may call Wait(), so the rest of the engine
// waits on its signal instead. A timeout is a failed unload, never a silent
// success: the caller must keep ownership of the process handle and report it.
func (e *Engine) reapStoppedProcess(timeout time.Duration) error {
	done := e.processDone()
	if done == nil {
		return nil
	}
	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("process did not exit within %s", timeout)
	}
}
