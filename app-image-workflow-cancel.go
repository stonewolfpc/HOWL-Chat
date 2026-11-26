package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// errImageWorkflowCancelled is returned in the workflow result when cancellation
// was accepted before the workflow reached its terminal stage. It is distinct
// from a stage failure so the UI can report the user's action.
var errImageWorkflowCancelled = errors.New("image workflow cancelled")

// imageWorkflowExecution is the cancellation scope for one workflow request.
// The context is the parent of every active stage generation. Completion is
// separated from cancellation acceptance so a race with normal completion
// closes done exactly once and refuses a late cancel.
type imageWorkflowExecution struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	mu        sync.Mutex
	cancelled bool
	finished  bool
	doneOnce  sync.Once
}

func workflowContextError(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: %v", errImageWorkflowCancelled, err)
	}
	return nil
}

func workflowStageRefused(ctx context.Context) error {
	if err := workflowContextError(ctx); err != nil {
		return err
	}
	return errImageWorkflowCancelled
}

func newImageWorkflowExecution() *imageWorkflowExecution {
	ctx, cancel := context.WithCancel(context.Background())
	return &imageWorkflowExecution{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}
}

func (r *imageWorkflowExecution) requestCancel() bool {
	if r == nil {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.finished {
		return false
	}
	if !r.cancelled {
		r.cancelled = true
		r.cancel()
	}
	return true
}

// acceptStage is the linearization point for beginning a stage. If cancellation
// wins the lock first, the stage is refused; if the stage wins, cancellation is
// observed by that stage's execution context. This closes the gap between a
// plain context check and the stage's first side effect.
func (r *imageWorkflowExecution) acceptStage() bool {
	if r == nil {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return !r.cancelled && !r.finished
}

func (r *imageWorkflowExecution) cancelRequested() bool {
	if r == nil {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.cancelled
}

func (r *imageWorkflowExecution) markFinished() {
	if r == nil {
		return
	}
	r.mu.Lock()
	r.finished = true
	r.mu.Unlock()
}

func (r *imageWorkflowExecution) complete() {
	if r == nil {
		return
	}
	r.doneOnce.Do(func() { close(r.done) })
}

func (a *App) setImageWorkflowRun(run *imageWorkflowExecution) {
	a.imageWorkflowRunMu.Lock()
	a.imageWorkflowRun = run
	a.imageWorkflowRunMu.Unlock()
}

func (a *App) clearImageWorkflowRun(run *imageWorkflowExecution) {
	a.imageWorkflowRunMu.Lock()
	if a.imageWorkflowRun == run {
		a.imageWorkflowRun = nil
	}
	a.imageWorkflowRunMu.Unlock()
}

func (a *App) currentImageWorkflowRun() *imageWorkflowExecution {
	a.imageWorkflowRunMu.Lock()
	defer a.imageWorkflowRunMu.Unlock()
	return a.imageWorkflowRun
}

// CancelImageWorkflow asks the active workflow to stop and waits until its
// parent, claims, active engine operation and process termination boundary have
// all completed. A cancelled stage context alone is not reported as complete.
func (a *App) CancelImageWorkflow() error {
	run := a.currentImageWorkflowRun()
	if run == nil || !run.requestCancel() {
		return fmt.Errorf("no active image workflow")
	}
	<-run.done
	return nil
}
