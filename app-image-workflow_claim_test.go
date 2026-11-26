package main

import (
	"testing"

	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/media"
	"howl-chat/internal/resource"
)

func TestWorkflowParentLeaseRefusesReclamation(t *testing.T) {
	app := &App{}
	lease := resource.NewWorkflowLease("engine-a")
	app.imageWorkflow = lease

	if !app.workflowReclaimBlocked() {
		t.Fatal("reclamation was allowed while the workflow parent lease was active")
	}
	claim, err := lease.Acquire(resource.NewComponentKey("engine-a:42", "upscaler", "/models/upscaler.pth", "f16", nil))
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if !app.workflowReclaimBlocked() {
		t.Fatal("reclamation was allowed while a component claim was active")
	}
	claim.Release()
	if !app.workflowReclaimBlocked() {
		t.Fatal("reclamation was allowed before the parent lease ended")
	}

	app.endImageWorkflow(lease)
	if app.workflowReclaimBlocked() {
		t.Fatal("reclamation stayed blocked after the parent lease ended")
	}
}

func TestWorkflowEngineReplacementInvalidatesScope(t *testing.T) {
	images := media.NewImage(t.TempDir())
	engine := sdcpp.NewEngine()
	images.SetEngine(engine)
	app := &App{images: images}
	lease := app.beginImageWorkflow(engine)
	t.Cleanup(func() { app.endImageWorkflow(lease) })

	if err := app.requireWorkflowEngine("stage 2 upscale"); err != nil {
		t.Fatalf("current engine rejected: %v", err)
	}
	images.SetEngine(sdcpp.NewEngine())
	if err := app.requireWorkflowEngine("stage 2 upscale"); err == nil {
		t.Fatal("workflow continued after its engine instance was replaced")
	}
}

// Characterization: closing the parent between stages is the only current
// cancellation-shaped action. It prevents a later stage from acquiring a claim,
// but the workflow loop has no cancellation signal of its own to observe.
func TestWorkflowCancelBetweenStages(t *testing.T) {
	app := &App{}
	lease := app.beginImageWorkflow(nil)
	key := resource.NewComponentKey("engine:1", "upscaler", "/models/upscaler.pth", "f16", nil)
	claim, err := app.claimWorkflowComponent(key)
	if err != nil {
		t.Fatalf("stage-1 claim: %v", err)
	}
	claim.Release()
	app.endImageWorkflow(lease)

	if _, err := app.claimWorkflowComponent(key); err == nil {
		t.Fatal("a stage after parent close acquired a component claim")
	}
	if err := app.requireWorkflowEngine("stage 2 upscale"); err == nil {
		t.Fatal("a stage after parent close accepted the stale workflow scope")
	}
}

func TestWorkflowStage4RefusesAbsentUpscalerClaim(t *testing.T) {
	app := &App{}
	lease := resource.NewWorkflowLease("engine-a")
	app.imageWorkflow = lease
	t.Cleanup(func() { app.endImageWorkflow(lease) })

	key := resource.NewComponentKey("engine-a:42", "upscaler", "/models/upscaler.pth", "f16", nil)
	if err := app.requireWorkflowComponent(key, "stage 4 final upscale"); err == nil {
		t.Fatal("stage 4 accepted an upscaler claim the workflow never acquired")
	}

	claim, err := app.claimWorkflowComponent(key)
	if err != nil {
		t.Fatalf("claim stage-2 upscaler: %v", err)
	}
	if err := app.requireWorkflowComponent(key, "stage 4 final upscale"); err != nil {
		t.Fatalf("stage 4 rejected the active stage-2 claim: %v", err)
	}
	claim.Release()

	if err := app.requireWorkflowComponent(key, "stage 4 final upscale"); err == nil {
		t.Fatal("stage 4 accepted a claim released before reuse")
	}
}
