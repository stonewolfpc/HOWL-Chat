package main

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestWorkflowCancelledScopeRefusesNextStage(t *testing.T) {
	run := newImageWorkflowExecution()
	if err := workflowContextError(run.ctx); err != nil {
		t.Fatalf("fresh workflow scope refused its first stage: %v", err)
	}
	if !run.requestCancel() {
		t.Fatal("cancellation was refused while the workflow was still active")
	}
	if run.acceptStage() {
		t.Fatal("a cancelled execution accepted a subsequent stage")
	}
	err := workflowContextError(run.ctx)
	if err == nil {
		t.Fatal("a cancelled scope allowed the next stage")
	}
	if !strings.Contains(err.Error(), "cancelled") {
		t.Fatalf("cancellation error does not identify the cause: %v", err)
	}
	run.markFinished()
	run.complete()
}

func TestWorkflowCancellationAfterCompletionIsRefused(t *testing.T) {
	run := newImageWorkflowExecution()
	run.markFinished()
	if run.requestCancel() {
		t.Fatal("cancellation was accepted after the workflow reached completion")
	}
	if run.cancelRequested() {
		t.Fatal("a refused cancellation marked the execution cancelled")
	}
	run.complete()
}

func TestWorkflowExecutionCompletionClosesOnce(t *testing.T) {
	run := newImageWorkflowExecution()
	run.complete()
	run.complete()
	select {
	case <-run.done:
	case <-time.After(time.Second):
		t.Fatal("completion did not close the execution scope")
	}
}

func TestWorkflowCancellationCompletionRace(t *testing.T) {
	for i := 0; i < 100; i++ {
		run := newImageWorkflowExecution()
		var wg sync.WaitGroup
		wg.Add(2)
		var accepted int
		go func() {
			defer wg.Done()
			if run.requestCancel() {
				accepted++
			}
		}()
		go func() {
			defer wg.Done()
			run.markFinished()
		}()
		wg.Wait()
		if accepted > 1 {
			t.Fatalf("iteration %d accepted cancellation more than once", i)
		}
		run.complete()
		select {
		case <-run.done:
		default:
			t.Fatalf("iteration %d left completion open", i)
		}
	}
}

func TestCancelImageWorkflowWaitsForExecutionCleanup(t *testing.T) {
	app := &App{}
	run := newImageWorkflowExecution()
	app.setImageWorkflowRun(run)

	release := make(chan struct{})
	cleanupReached := make(chan struct{})
	go func() {
		<-run.ctx.Done()
		run.markFinished()
		<-release
		close(cleanupReached)
		run.complete()
	}()

	callDone := make(chan error, 1)
	go func() { callDone <- app.CancelImageWorkflow() }()

	select {
	case err := <-callDone:
		t.Fatalf("cancellation returned before execution cleanup: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	close(release)
	select {
	case err := <-callDone:
		if err != nil {
			t.Fatalf("CancelImageWorkflow: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("CancelImageWorkflow did not return after execution cleanup")
	}
	select {
	case <-cleanupReached:
	default:
		t.Fatal("cancellation returned without reaching cleanup")
	}
	app.clearImageWorkflowRun(run)
}

func TestCancelImageWorkflowWithoutActiveRun(t *testing.T) {
	app := &App{}
	if err := app.CancelImageWorkflow(); err == nil {
		t.Fatal("cancelling with no active workflow returned success")
	}
}

func TestWorkflowCancellationContextIsChildNotBackground(t *testing.T) {
	run := newImageWorkflowExecution()
	child, cancel := context.WithCancel(run.ctx)
	defer cancel()
	run.requestCancel()
	select {
	case <-child.Done():
	case <-time.After(time.Second):
		t.Fatal("cancelling the workflow scope did not reach a stage context")
	}
	run.markFinished()
	run.complete()
}
