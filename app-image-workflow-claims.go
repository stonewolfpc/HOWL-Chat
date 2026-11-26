package main

import (
	"fmt"

	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/resource"
)

// workflowEngineInstance identifies the engine object a workflow belongs to.
// Replacing that engine invalidates the parent scope even when the model name
// is unchanged.
func workflowEngineInstance(engine *sdcpp.Engine) string {
	if engine == nil {
		return ""
	}
	return fmt.Sprintf("%p", engine)
}

// workflowEngineIdentity identifies the running process for a component claim.
// The PID is part of the key, so a reload cannot let stage 4 reuse a claim that
// belonged to the previous process.
func workflowEngineIdentity(engine *sdcpp.Engine) string {
	return fmt.Sprintf("%s:%d", workflowEngineInstance(engine), engine.ServerPID())
}

func (a *App) beginImageWorkflow(engine *sdcpp.Engine) *resource.WorkflowLease {
	lease := resource.NewWorkflowLease(workflowEngineInstance(engine))
	a.imageWorkflow = lease
	return lease
}

func (a *App) endImageWorkflow(lease *resource.WorkflowLease) {
	if lease == nil {
		return
	}
	if a.imageWorkflow == lease {
		a.imageWorkflow = nil
	}
	lease.Release()
}

// claimWorkflowComponent records a stage's reference to a retained component.
func (a *App) claimWorkflowComponent(key resource.ComponentKey) (*resource.ComponentClaim, error) {
	if a.imageWorkflow == nil {
		return nil, fmt.Errorf("workflow parent lease is not active")
	}
	return a.imageWorkflow.Acquire(key)
}

func (a *App) verifyWorkflowComponent(key resource.ComponentKey) bool {
	return a.imageWorkflow != nil && a.imageWorkflow.Verify(key)
}

func (a *App) workflowReclaimBlocked() bool {
	lease := a.imageWorkflow
	return lease != nil && !lease.Released()
}

// requireWorkflowEngine refuses to continue after the image engine instance was
// replaced. A stage must not silently keep using a stale engine pointer.
func (a *App) requireWorkflowEngine(stage string) error {
	lease := a.imageWorkflow
	if lease == nil || lease.Released() {
		return fmt.Errorf("%s: workflow parent lease is not active", stage)
	}
	current, ok := a.images.Engine().(*sdcpp.Engine)
	if !ok || current == nil || lease.EngineID() != workflowEngineInstance(current) {
		return fmt.Errorf("%s: image engine was replaced; the workflow scope is stale", stage)
	}
	return nil
}

func (a *App) requireWorkflowComponent(key resource.ComponentKey, stage string) error {
	if a.verifyWorkflowComponent(key) {
		return nil
	}
	return fmt.Errorf("%s: upscaler claim is absent; refusing to assume the previous stage's adapter is resident", stage)
}

// workflowUpscalerKey is the stable component identity for a retained upscaler.
// The path is normalized by resource.NewComponentKey; companion order is not
// identity.
func (a *App) workflowUpscalerKey(engine *sdcpp.Engine, upscalerPath string) resource.ComponentKey {
	loaded := a.images.LoadedOptions()
	companions := []string{
		loaded.VAEPath,
		loaded.ClipLPath,
		loaded.ClipGPath,
		loaded.T5XXLPath,
		loaded.LLMPath,
		loaded.DiffusionModel,
		loaded.ControlNetPath,
		loaded.IPAdapterPath,
		loaded.ClipVisionPath,
	}
	return resource.NewComponentKey(workflowEngineIdentity(engine), "upscaler", upscalerPath, loaded.WeightType, companions)
}
