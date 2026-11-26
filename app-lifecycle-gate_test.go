package main

import (
	"strings"
	"testing"

	"howl-chat/internal/resource"
)

// A chat unload must not reach the process registry while a chat load is
// starting: the gate refuses first, and the refusal is what the user sees.
func TestUnloadModelRefusesWhileChatIsStarting(t *testing.T) {
	app := &App{}
	start, err := app.lifecycleTransitions().Begin(resource.ResourceLLM, "loading")
	if err != nil {
		t.Fatalf("start transition: %v", err)
	}
	defer start.Release()

	err = app.UnloadModel()
	if err == nil {
		t.Fatal("unloading a starting chat model must be refused")
	}
	if !strings.Contains(err.Error(), "loading") {
		t.Fatalf("refusal does not name the conflicting transition: %v", err)
	}
}

// The image reclaimer holds the image slot as "reloading". A concurrent unload
// is a different transition on the same workload and must be refused, not
// interleaved.
func TestUnloadImageModelRefusesWhileImageIsReloading(t *testing.T) {
	app := &App{}
	reload, err := app.lifecycleTransitions().Begin(resource.ResourceImage, "reloading")
	if err != nil {
		t.Fatalf("reload transition: %v", err)
	}
	defer reload.Release()

	err = app.UnloadImageModel()
	if err == nil {
		t.Fatal("unloading during a reload must be refused")
	}
	if !strings.Contains(err.Error(), "reloading") {
		t.Fatalf("refusal does not name the conflicting transition: %v", err)
	}
}

// An active image workflow owns the engine for its stages. Freeing the card for
// chat must not tear that engine down mid-stage; the refusal names the owner.
func TestUnloadImageModelRefusesDuringActiveWorkflow(t *testing.T) {
	app := &App{}
	app.imageWorkflow = resource.NewWorkflowLease("engine-under-test")
	defer app.imageWorkflow.Release()

	err := app.UnloadImageModel()
	if err == nil {
		t.Fatal("unloading the engine during a workflow must be refused")
	}
	if !strings.Contains(err.Error(), "workflow") {
		t.Fatalf("refusal does not name the workflow owner: %v", err)
	}
}

// After the workflow ends, the same call no longer reports the workflow; with
// no engine it reports that instead, which proves the guard is what refused
// earlier rather than an unconditional error.
func TestUnloadImageModelFallsThroughAfterWorkflowRelease(t *testing.T) {
	app := &App{}
	lease := resource.NewWorkflowLease("engine-under-test")
	app.imageWorkflow = lease
	lease.Release()

	err := app.UnloadImageModel()
	if err == nil {
		t.Fatal("no image engine is configured, so the call should still fail")
	}
	if strings.Contains(err.Error(), "workflow") {
		t.Fatalf("a released workflow must not keep blocking the unload: %v", err)
	}
}
