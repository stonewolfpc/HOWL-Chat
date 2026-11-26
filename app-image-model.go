package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"howl-chat/internal/errors"
	"howl-chat/internal/image/deps"
	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/image/registry"
	"howl-chat/internal/image/scanner"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/resource"
)

// -----------------------------------------------------------------------
// Model loading
// -----------------------------------------------------------------------

// ensureEngineReady auto-initializes the engine and loads the model if needed.
// If genOpts is provided, its startup flags are merged into lastImageModelOpts.
func (a *App) ensureEngineReady(modelName, jobID string, genOpts ...imagetypes.ImageOptions) error {
	a.imageStartMu.Lock()
	defer a.imageStartMu.Unlock()
	return a.ensureEngineReadyLocked(modelName, jobID, genOpts...)
}

// ensureEngineReadyLocked is the body for callers that already hold
// imageStartMu. The sequential workflow holds it across all stages so a queued
// job cannot interleave a load or a reclamation between them.
func (a *App) ensureEngineReadyLocked(modelName, jobID string, genOpts ...imagetypes.ImageOptions) error {
	return a.ensureEngineReadyExecutionLocked(modelName, jobID, nil, genOpts...)
}

// ensureVideoExecution keeps the video admission lease open until the caller
// releases it, after sd-cli has exited. The prepare phase resolves the model and
// starts no sd-server.
func (a *App) ensureVideoExecution(modelName, jobID string, genOpts ...imagetypes.ImageOptions) (func(), error) {
	a.imageStartMu.Lock()
	defer a.imageStartMu.Unlock()
	return a.ensureVideoExecutionLocked(modelName, jobID, genOpts...)
}

func (a *App) ensureVideoExecutionLocked(modelName, jobID string, genOpts ...imagetypes.ImageOptions) (func(), error) {
	var release func()
	err := a.ensureEngineReadyExecutionLocked(modelName, jobID, &release, genOpts...)
	return release, err
}

func (a *App) ensureEngineReadyExecutionLocked(modelName, jobID string, holdVideoLease *func(), genOpts ...imagetypes.ImageOptions) error {
	// The start is a transition: admission, spawn and registration happen under
	// the gate, and the gate is released before any generation runs. This is the
	// body every caller reaches, including the sequential workflow, so there is
	// no bypass around it.
	transition, gateErr := a.beginLifecycleTransition(resource.ResourceImage, "loading", lifecycleStartWait)
	if gateErr != nil {
		return gateErr
	}
	defer transition.Release()

	var genOptsValue imagetypes.ImageOptions
	if len(genOpts) > 0 {
		genOptsValue = genOpts[0]
	}
	if a.images.Engine() == nil {
		if err := a.lazyInitEngine(); err != nil {
			return err
		}
	}
	if modelName == "" {
		modelName = genOptsValue.ModelName
	}
	if modelName == "" {
		modelName = a.images.LastModelName()
	}
	if modelName == "" {
		return fmt.Errorf("no model selected")
	}
	loadOpts, err := a.imageLoadOptionsForGeneration(modelName, genOptsValue)
	if err != nil {
		return err
	}
	// EnsureReady is a no-op when this model is already running. In that case
	// the process kept its original startup flags, so re-registering the
	// request would overwrite the actual loaded type with the saved one and make
	// the next admission charge f16 for a q8 process. A changed startup flag,
	// however, must actually restart the process: otherwise the control is inert.
	alreadyLoaded := a.imageEngineLoadedSame(modelName)
	reload := false
	if alreadyLoaded {
		loadOpts = preserveImageLoadedOptions(loadOpts, a.images.LoadedOptions())
		reload = imageStartupOptionsDiffer(loadOpts, a.images.LoadedOptions())
	}

	a.images.SetLoadOptions(loadOpts)
	// The preflight plans automatic placement without warning (so the start does
	// not warn twice). The actual start announces the selected trade once.
	if (!alreadyLoaded || reload) && loadOpts.AutoPlacement {
		a.announceImageAutoPlacement(loadOpts)
	}
	// The job id belongs to this request, not to the remembered model selection:
	// it is what lets one job's cancellation stop its own load and nothing else.
	loadOpts.JobID = jobID

	// Video renders through sd-cli, so resolve the model without spawning
	// sd-server: only one process should own the GPU at a time. sd-cli still
	// starts a model process, so the same admission check applies.
	if genOptsValue.Mode == "video" {
		if eng, ok := a.images.Engine().(*sdcpp.Engine); ok {
			_, lease, chosen, admitErr := a.admitImageLoad(loadOpts, genOptsValue, true)
			if admitErr != nil {
				return admitErr
			}
			if err := eng.SelectModel(chosen); err != nil {
				lease.Release()
				return err
			}
			if holdVideoLease != nil {
				// The video execution scope owns this reservation until sd-cli exits.
				*holdVideoLease = lease.Release
			} else {
				lease.Release()
			}
			return nil
		}
	}

	decision, lease, chosen, admitErr := a.admitImageLoad(loadOpts, genOptsValue, reload)
	if admitErr != nil {
		a.emitImagePlacement("failed")
		return admitErr
	}
	if alreadyLoaded && reload {
		err = a.images.Engine().LoadModel(chosen)
	} else {
		err = a.images.Engine().EnsureReady(modelName, chosen)
	}
	if err != nil {
		lease.Release()
		// A failed start can have stopped the previous child. A registration
		// may not outlive the process it described.
		a.invalidateImageRegistrationsIfStopped()
		return err
	}
	if a.level0().Gate.Stopping() {
		lease.Release()
		_ = a.images.Engine().Shutdown()
		return fmt.Errorf("image model load cancelled because Howl is shutting down")
	}
	if !alreadyLoaded || reload {
		a.registerImageModelWithObservation(chosen, decision)
	}
	lease.Release()
	a.emitImagePlacement("load")
	return nil
}

// lazyInitEngine initialises the engine on-demand during generation if it was
// not available at startup (e.g. binary configured later).
func (a *App) lazyInitEngine() error {
	cfg := a.GetImageEngineConfig()
	if cfg.BinaryPath == "" {
		cfg.BinaryPath = findSDServerBinary()
	}
	if cfg.BinaryPath == "" {
		return fmt.Errorf("sd-server binary not found")
	}
	if cfg.ModelsDir == "" {
		cfg.ModelsDir = defaultModelsDir()
	}
	if cfg.Port == 0 {
		cfg.Port = 8081
	}
	return a.reinitEngine(func(c *imagetypes.EngineConfig) { *c = cfg })
}

// SetActiveImageModel loads a model into sd-server and persists the selection.
//
// Not cancellable: a direct load has no job identity, so CancelImageJob cannot
// address it. Cancellable loads run inside StartImageJob, where the job owns the
// load it started.
func (a *App) SetActiveImageModel(opts interface{}) (err error) {
	if err := a.ensureAcceptingWork("load image model"); err != nil {
		return err
	}
	a.imageStartMu.Lock()
	defer a.imageStartMu.Unlock()
	transition, gateErr := a.beginLifecycleTransition(resource.ResourceImage, "loading", lifecycleStartWait)
	if gateErr != nil {
		return gateErr
	}
	defer transition.Release()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[PANIC] SetActiveImageModel: %v\n", r)
			err = fmt.Errorf("internal panic: %v", r)
		}
	}()
	if a.images.Engine() == nil {
		return fmt.Errorf("image engine not initialized")
	}

	loadOpts, parseErr := parseModelLoadOpts(opts)
	if parseErr != nil {
		return parseErr
	}

	a.images.SetActiveModel(loadOpts.ModelName)
	loadOpts.ModelPath = a.resolveModelPath(loadOpts.ModelName)
	loadOpts.Family = a.registryFamily(loadOpts.ModelName)
	a.autofillCompanionModels(&loadOpts)
	if err := a.validateModelLoadOptions(loadOpts); err != nil {
		return err
	}
	a.images.SetLoadOptions(loadOpts)
	fmt.Printf("[App] SetActiveImageModel: name=%q path=%q\n", loadOpts.ModelName, loadOpts.ModelPath)
	decision, lease, chosen, admitErr := a.admitImageLoad(loadOpts, imagetypes.ImageOptions{ModelName: loadOpts.ModelName}, true)
	if admitErr != nil {
		return admitErr
	}
	if err := a.images.Engine().LoadModel(chosen); err != nil {
		lease.Release()
		if chosen.OffloadToCPU {
			a.logError(errors.IMGCPUOffloadFail, errors.High, "image.load", "app-image-model.go", err)
		}
		a.logImageLoadError(err, "app-image-model.go")
		a.invalidateImageRegistrationsIfStopped()
		return err
	}

	if a.level0().Gate.Stopping() {
		lease.Release()
		_ = a.images.Engine().Shutdown()
		return fmt.Errorf("image model load cancelled because Howl is shutting down")
	}

	a.registerImageModelWithObservation(chosen, decision)
	lease.Release()
	a.emitImagePlacement("load")
	a.saveImageSettings()
	return nil
}

// admitImageLoad is the one admission point for image model starts. It credits
// an already-loaded image model only when the engine will stop it, and holds the
// reservation until the load has completed. A generation against an already
// loaded model is checked for its additional activation footprint.
func (a *App) admitImageLoad(loadOpts imagetypes.ModelLoadOptions, gen imagetypes.ImageOptions, forceStartup bool) (resource.AdmissionDecision, *resource.AdmissionLease, imagetypes.ModelLoadOptions, error) {
	if a.admission == nil {
		a.admission = resource.DefaultAuthority()
	}
	chosen := loadOpts
	// Startup flags only affect a process that is about to start. A generation
	// against the already-loaded model cannot change them, so do not claim a
	// fallback was applied. An explicit load/reload always restarts the process
	// and may change them even when the model name is unchanged.
	generationRequest := gen.Width > 0 || gen.Height > 0 || gen.Mode != ""
	startupWillApply := forceStartup || !generationRequest || !a.imageEngineLoadedSame(chosen.ModelName)
	// Apply the constrained/refusal decision at the one admission point as well
	// as in imageLoadOptionsForGeneration. Direct loads and preflight reach this
	// function without the generation wrapper; without this call a small budget
	// could still acquire a lease before the spawn-time check refused it.
	if startupWillApply {
		if err := a.applyConstrainedOffload(&chosen, gen); err != nil {
			return resource.AdmissionDecision{}, nil, chosen, err
		}
	}
	req := a.imageAdmissionRequest(chosen, gen, !forceStartup && generationRequest && !startupWillApply)
	decision, lease := a.admission.Admit(req)
	fmt.Printf("[App] image admission: admitted=%t activation=%d MB %s\n", decision.Admitted, decision.Estimate.ActivationVRAMMB, decision.Reason)

	// The stored type did not fit. q8_0 is the only automatic conversion: it is
	// near-lossless, and silently dropping to q4/q5 would change quality without
	// the user asking. Never touch an explicit type, and never re-quantize a
	// GGUF whose stored type is already the user's choice of file.
	if !decision.Admitted && startupWillApply && imageAutoQuantizeEligible(chosen) {
		quantized := chosen
		quantized.WeightType = "q8_0"
		q8Req := a.imageAdmissionRequest(quantized, gen, false)
		q8Decision, q8Lease := a.admission.Admit(q8Req)
		if q8Decision.Admitted {
			chosen = quantized
			decision, lease = q8Decision, q8Lease
			a.warnImageHardware("Model did not fit at its stored precision; loading q8_0 (near-lossless) to fit the GPU. Set Weight Type explicitly to override.")
		} else if q8Decision.Reason != "" {
			decision.Reason += "; q8_0 fallback also refused: " + q8Decision.Reason
		}
	}

	if !decision.Admitted {
		a.logError(errors.IMGModelLoadFail, errors.Critical, "image.load.admission", "app-image-model.go", fmt.Errorf("%s", decision.Reason))
		return decision, nil, chosen, fmt.Errorf("%s", decision.Reason)
	}

	// sd.cpp's decode is the largest workload-dependent allocation. The engine
	// can recover by tiling, but it pays a failed decode first. A CPU-only run
	// has no VRAM decode budget to protect. The authority plans tiling when the
	// measured untiled peak would cross the reserve and the measured tiled peak
	// fits; start the same engine flag so the projection is what actually runs.
	if startupWillApply && !chosen.VAETiling && !imageCPUOnly(chosen) && decision.Estimate.TilingPlanned {
		chosen.VAETiling = true
		a.warnImageHardware(fmt.Sprintf("VAE tiling planned for %dx%d: tiled projection %d MB fits the %d MB post-reserve budget. Decode will be slower and is tiled by sd.cpp.", gen.Width, gen.Height, decision.Estimate.ProjectedPeakVRAMMB, decision.Estimate.AvailableVRAMMB))
	}
	// The pinned engine loads ControlNet params on CPU but stages the full
	// 1415.82 MiB set to CUDA for sampling. When admission projects that staging
	// would cross the reserve, keep both params and compute on CPU. The process
	// stays the same and no resident base component is unloaded.
	if startupWillApply && !chosen.ControlNetCPU && decision.Estimate.ControlNetCPUPlanned {
		chosen.ControlNetCPU = true
		a.warnImageHardware(fmt.Sprintf("ControlNet CPU placement planned: staging %d MB of ControlNet params would exceed the %d MB post-reserve budget. Diffusion stays on GPU; ControlNet compute will be slower.", decision.Estimate.ControlNetVRAMMB, decision.Estimate.AvailableVRAMMB))
	}
	// Print the final placement after the q8 fallback and the tiling/ControlNet
	// plans, so the startup decision is auditable even when the first estimate
	// was refused at the stored precision.
	fmt.Printf("[App] image admission final: admitted=%t projected=%dMB resident=%dMB activation=%dMB required=%dMB tiling=%t controlnet_cpu=%t free=%dMB reserve=%dMB %s\n",
		decision.Admitted, decision.Estimate.ProjectedPeakVRAMMB, decision.Estimate.ResidentVRAMMB,
		decision.Estimate.ActivationVRAMMB, decision.Estimate.RequiredVRAMMB,
		decision.Estimate.TilingPlanned, decision.Estimate.ControlNetCPUPlanned,
		decision.Estimate.FreeVRAMMB, decision.Estimate.ReserveVRAMMB, decision.Reason)

	// A blank Max VRAM is Auto, and Auto must still keep the 10% reserve. sd-server's
	// own default (-1) reserves only 1 GiB, which is less than 10% on a large card.
	// An automatic post-reserve cap makes the policy unconditional; an explicit user
	// --max-vram always wins. This also keeps --offload-to-cpu from spending the
	// reserve when it opportunistically loads modules into VRAM.
	if startupWillApply && !imageCPUOnly(chosen) && chosen.MaxVRAM == nil {
		budgetGiB := decision.Estimate.AvailableVRAMMB / 1024
		if budgetGiB < 1 {
			budgetGiB = 1
		}
		budget := int(budgetGiB)
		chosen.MaxVRAM = &budget
		chosen.AutoMaxVRAM = true
		a.warnImageHardware(fmt.Sprintf("Max VRAM Auto uses live free memory; Howl capped it at %d GiB so the 10%% reserve stays free.", budgetGiB))
	}
	// sd.cpp documents 0 as "live free VRAM without an explicit budget", which
	// would let the engine spend the reserve. Keep the user's intent (use live
	// free memory) but express it as a post-reserve positive cap.
	if startupWillApply && !imageCPUOnly(chosen) && chosen.MaxVRAM != nil && *chosen.MaxVRAM > 0 && !chosen.AutoMaxVRAM {
		availableGiB := decision.Estimate.AvailableVRAMMB / 1024
		requested := *chosen.MaxVRAM
		if availableGiB >= 1 && int64(requested) > availableGiB {
			clamped := int(availableGiB)
			chosen.MaxVRAM = &clamped
			chosen.AutoMaxVRAM = true
			a.warnImageHardware(fmt.Sprintf("Max VRAM %d GiB exceeds the post-reserve budget; capped at %d GiB.", requested, clamped))
		}
	}
	if startupWillApply && !imageCPUOnly(chosen) && chosen.MaxVRAM != nil && *chosen.MaxVRAM == 0 && !chosen.AutoMaxVRAM {
		budgetGiB := decision.Estimate.AvailableVRAMMB / 1024
		if budgetGiB < 1 {
			budgetGiB = 1
		}
		budget := int(budgetGiB)
		chosen.MaxVRAM = &budget
		chosen.AutoMaxVRAM = true
		a.warnImageHardware(fmt.Sprintf("Max VRAM 0 uses live free memory; Howl capped it at %d GiB so the 10%% reserve stays free.", budgetGiB))
	}

	// Residency TTL is never automated: that would discard reusable weights.
	// Staging TTL is not exposed either: the engine already releases staged GPU
	// copies when a component finishes, and measured 0/2s values changed neither
	// post-generation VRAM nor RAM on this build.

	return decision, lease, chosen, nil
}

func imageAutoQuantizeEligible(loadOpts imagetypes.ModelLoadOptions) bool {
	if strings.TrimSpace(loadOpts.WeightType) != "" {
		return false
	}
	switch strings.ToLower(filepath.Ext(loadOpts.ModelPath)) {
	case ".safetensors", ".sft", ".ckpt":
		return true
	default:
		return false
	}
}

// measuredImageResidentVRAMMB returns the engine's own last reported VRAM
// total when the resident process is the same model and component set as the
// request. It is evidence from the placement authority, not an estimate. A
// different model, adapter, encoder or ControlNet invalidates it.
func (a *App) measuredImageResidentVRAMMB(loadOpts imagetypes.ModelLoadOptions, gen imagetypes.ImageOptions) int64 {
	if a.images == nil || a.images.Engine() == nil || !a.images.Engine().IsRunning() {
		return 0
	}
	engine, ok := a.images.Engine().(*sdcpp.Engine)
	if !ok || engine == nil {
		return 0
	}
	loadedModel := engine.GetLoadedModel()
	if loadedModel == nil || loadedModel.Name == "" || loadedModel.Name != loadOpts.ModelName {
		return 0
	}
	actual := engine.GetPlacement()
	if !actual.Known || actual.VRAMMB <= 0 {
		return 0
	}
	loaded := a.images.LoadedOptions()
	samePath := func(a, b string) bool {
		if a == "" || b == "" {
			return a == b
		}
		return filepath.Clean(a) == filepath.Clean(b)
	}
	if !samePath(loaded.ModelPath, loadOpts.ModelPath) ||
		!samePath(loaded.ControlNetPath, loadOpts.ControlNetPath) ||
		!samePath(loaded.IPAdapterPath, loadOpts.IPAdapterPath) ||
		!samePath(loaded.ClipVisionPath, loadOpts.ClipVisionPath) ||
		!strings.EqualFold(strings.TrimSpace(loaded.WeightType), strings.TrimSpace(loadOpts.WeightType)) {
		return 0
	}
	return int64(actual.VRAMMB + 0.5)
}

func (a *App) imageEngineLoadedSame(modelName string) bool {
	if a.images == nil || modelName == "" {
		return false
	}
	eng := a.images.Engine()
	if eng == nil || !eng.IsRunning() {
		return false
	}
	loaded := eng.GetLoadedModel()
	return loaded != nil && loaded.Name == modelName
}

func imageCPUOnly(loadOpts imagetypes.ModelLoadOptions) bool {
	return strings.EqualFold(strings.TrimSpace(loadOpts.Backend), "cpu")
}

func imageWeightsOffloaded(loadOpts imagetypes.ModelLoadOptions) bool {
	return loadOpts.OffloadToCPU ||
		strings.EqualFold(strings.TrimSpace(loadOpts.ParamsBackend), "cpu") ||
		strings.EqualFold(strings.TrimSpace(loadOpts.ParamsBackend), "disk")
}

// imageConstrainedOffloadNeedMB is the measured process need for the hard
// low-VRAM mode below. A 512x512 CPU-diffusion run peaked at 3526 MB total GPU
// against an 892 MB desktop baseline, i.e. about 2636 MB for the process; a
// 2 GiB budget cannot hold that, while 3 GiB does (and a 1024 tiled run was
// lower, 2075 MB total).
const imageConstrainedOffloadNeedMB = int64(2636)

// imageConstrainedOffloadFloorMB is the conservative admission floor Howl uses
// rather than admitting a budget close to the measured 2636 MB need. It leaves
// room for driver/allocator overhead and the tiled VAE path, which measured
// below this on this card; it is a policy floor, not a claim that 3072 MB was
// itself the process peak.
const imageConstrainedOffloadFloorMB = int64(3072)

// imageConstrainedBackend is the measured per-module placement that keeps
// diffusion and CLIP on CPU while a tiled VAE remains on CUDA. That is the only
// combination measured to bound this process under a 3 GiB budget.
const imageConstrainedBackend = "diffusion=cpu,clip=cpu,vae=cuda0"

// applyConstrainedOffload turns a budget that is too small for a resident GPU
// diffusion model into a hard low-VRAM execution mode: diffusion runs on CPU,
// while CLIP and the VAE stay on CUDA and the VAE is tiled. The measured direct
// engine runs show this is the only advertised combination that bounds the
// process under a 3 GiB budget (512: 3526 MB total, 1024+tiling: 2075 MB total).
// Ordinary --offload-to-cpu still peaks near 6.2 GiB because sd.cpp loads the
// fitting modules into VRAM when a component runs. This mode is slower and the
// warning says so; it is never silent and it is not used when the budget can
// hold the resident weights plus the decode working set.
func (a *App) applyConstrainedOffload(loadOpts *imagetypes.ModelLoadOptions, gen imagetypes.ImageOptions) error {
	return a.constrainOffload(loadOpts, gen, true)
}

// planConstrainedOffload is the same decision without the user-facing warning.
// The read-only preflight uses it so the start that follows emits the trade
// once, not twice.
func (a *App) planConstrainedOffload(loadOpts *imagetypes.ModelLoadOptions, gen imagetypes.ImageOptions) error {
	return a.constrainOffload(loadOpts, gen, false)
}

func (a *App) constrainOffload(loadOpts *imagetypes.ModelLoadOptions, gen imagetypes.ImageOptions, announce bool) error {
	if loadOpts == nil || imageCPUOnly(*loadOpts) {
		return nil
	}
	// The same decision is made by preflight and by admission. Once Howl has
	// selected a bounded mode, a second pass must not warn twice or rewrite it.
	if loadOpts.AutoPlacement {
		return nil
	}
	// An automatic post-reserve cap is Howl's own policy expression, not a
	// user-requested hard budget. Treating it as one made the next request
	// switch a just-proven resident workflow to CPU diffusion whenever the
	// conservative module estimate exceeded the cap.
	if loadOpts.AutoMaxVRAM {
		return nil
	}
	if loadOpts.MaxVRAM == nil || *loadOpts.MaxVRAM <= 0 {
		return nil
	}
	// An explicit backend or params backend is the user's placement choice; do
	// not rewrite it behind their back. Auto (both empty) is Howl's to fit.
	if strings.TrimSpace(loadOpts.Backend) != "" || strings.TrimSpace(loadOpts.ParamsBackend) != "" {
		return nil
	}
	req := a.imageAdmissionRequest(*loadOpts, gen, false)
	weights, activation, ok := resource.ImageResidentNeedMB(req)
	if !ok {
		return nil
	}
	budgetMB := int64(*loadOpts.MaxVRAM) * 1024
	if budgetMB <= 0 || weights+activation <= budgetMB {
		return nil
	}
	// Reuse the last engine-measured resident set for the same model and
	// components. The module+reserve estimate is deliberately conservative and
	// can call for CPU diffusion at 5 GiB even when the real process needs 3.1
	// GiB. Measurement wins only for an identical composition; any changed
	// adapter, encoder or ControlNet discards it.
	if measured := a.measuredImageResidentVRAMMB(*loadOpts, gen); measured > 0 {
		if _, act, ok := resource.ImageResidentNeedMB(req); ok && measured+act+resource.ImageFixedOverheadMB <= budgetMB {
			return nil
		}
	}

	// The measured GPU-assisted bounded mode is the best use of this budget
	// when it can hold the process. It keeps CLIP and the VAE on CUDA while the
	// diffusion weights run from system RAM. Below its measured process need,
	// the honest fallback is CPU-only: no device compute, so no way to exceed
	// the VRAM budget. Both are explicit trades and both are announced.
	if budgetMB >= imageConstrainedOffloadFloorMB {
		loadOpts.Backend = imageConstrainedBackend
		loadOpts.OffloadToCPU = true
		loadOpts.VAETiling = true
		loadOpts.AutoPlacement = true
		loadOpts.AutoPlacementMode = "cpu-diffusion"
		if announce {
			a.warnImageHardware(fmt.Sprintf(
				"VRAM Budget %d GiB is below the %d MB resident weights plus %d MB decode working set. Automatic partial placement: diffusion and CLIP on CPU, VAE on CUDA and tiled. This stays inside the budget but is much slower.",
				*loadOpts.MaxVRAM, weights, activation))
		}
		return nil
	}

	loadOpts.Backend = "cpu"
	loadOpts.ParamsBackend = ""
	loadOpts.OffloadToCPU = false
	loadOpts.VAETiling = true
	loadOpts.AutoPlacement = true
	loadOpts.AutoPlacementMode = "cpu-only"
	if announce {
		a.warnImageHardware(fmt.Sprintf(
			"VRAM Budget %d GiB is below the measured %d MB GPU-assisted CPU-diffusion process need (conservative %d MB floor). Automatic CPU-only placement selected: no weights or compute on the GPU, so the budget is respected; expect tens of seconds per image.",
			*loadOpts.MaxVRAM, imageConstrainedOffloadNeedMB, imageConstrainedOffloadFloorMB))
	}
	return nil
}

// preserveImageLoadedOptions carries the actual fallback decisions forward when
// the request left them blank. Without this, any unrelated reload would silently
// undo an automatic q8_0 conversion or auto tiling, and an empty weight type
// would be charged as f16 again.
func preserveImageLoadedOptions(req, actual imagetypes.ModelLoadOptions) imagetypes.ModelLoadOptions {
	if req.ControlNetPath == "" {
		req.ControlNetPath = actual.ControlNetPath
	}
	if req.UpscalerDir == "" {
		req.UpscalerDir = actual.UpscalerDir
	}
	if req.WeightType == "" {
		req.WeightType = actual.WeightType
	}
	if !req.VAETiling {
		req.VAETiling = actual.VAETiling
	}
	if !req.TemporalTiling {
		req.TemporalTiling = actual.TemporalTiling
	}
	if req.FlashAttentionMode == "" {
		req.FlashAttentionMode = actual.FlashAttentionMode
	}
	if req.WeightResidencyTTL == nil {
		req.WeightResidencyTTL = actual.WeightResidencyTTL
	}
	// An automatic bounded placement is runtime state of the running process.
	// A later generation whose controls still say Auto must reuse it, or every
	// request would look like a changed startup flag and restart the engine.
	if actual.AutoPlacement && imageIntPtrEqual(req.MaxVRAM, actual.MaxVRAM) && strings.TrimSpace(req.Backend) == "" && strings.TrimSpace(req.ParamsBackend) == "" {
		req.Backend = actual.Backend
		req.ParamsBackend = actual.ParamsBackend
		req.OffloadToCPU = actual.OffloadToCPU
		req.VAETiling = actual.VAETiling || req.VAETiling
		req.AutoPlacement = true
		req.AutoPlacementMode = actual.AutoPlacementMode
	}
	if req.MaxVRAM == nil && imageWeightsOffloaded(req) {
		req.MaxVRAM = actual.MaxVRAM
		req.AutoMaxVRAM = actual.AutoMaxVRAM
	}
	// A blank Max VRAM control is Auto, and an automatic cap is how that Auto
	// choice was expressed to the running process. Treat the blank request as
	// the actual auto cap; otherwise the new MaxVRAM startup comparison would
	// restart the engine before every generation while the UI correctly says
	// Auto. An actual explicit cap (AutoMaxVRAM=false) still differs and forces
	// the restart that clears it.
	if req.MaxVRAM == nil && actual.AutoMaxVRAM {
		req.MaxVRAM = actual.MaxVRAM
		req.AutoMaxVRAM = true
	}
	// 0 means "live free memory" but Howl turned it into a positive auto cap at
	// startup. Keep comparing against that cap instead of restarting every
	// generation because the saved field is still 0.
	if req.MaxVRAM != nil && *req.MaxVRAM == 0 && actual.AutoMaxVRAM {
		req.MaxVRAM = actual.MaxVRAM
		req.AutoMaxVRAM = true
	}
	// A saved explicit cap that Howl had to clamp is still the user's saved value;
	// compare against the actual clamp so it does not restart every generation.
	if req.MaxVRAM != nil && actual.AutoMaxVRAM && actual.MaxVRAM != nil && *req.MaxVRAM > *actual.MaxVRAM {
		req.MaxVRAM = actual.MaxVRAM
		req.AutoMaxVRAM = true
	}
	return req
}

// imageStartupOptionsDiffer reports whether a running process must be restarted
// for the request to mean what its controls say. Empty strings and nil pointers
// mean "no explicit change", not "use a default", so they never force a reload.
func imageStartupOptionsDiffer(req, actual imagetypes.ModelLoadOptions) bool {
	if req.EngineProfile != actual.EngineProfile {
		return true
	}
	if req.EngineBinaryPath != "" && req.EngineBinaryPath != actual.EngineBinaryPath {
		return true
	}
	if req.IPAdapterPath != actual.IPAdapterPath || req.ClipVisionPath != actual.ClipVisionPath || req.ClipVisionCPU != actual.ClipVisionCPU {
		return true
	}
	if req.WeightType != "" && req.WeightType != actual.WeightType {
		return true
	}
	if req.Backend != "" && req.Backend != actual.Backend {
		return true
	}
	// An empty control means Auto, which is a real choice. Restart when the
	// resident process still carries an explicit backend/params flag, so
	// clearing the control cannot silently leave the old placement running.
	// Auto that resolved to no flag stores an empty actual value, so the next
	// automatic generation does not restart.
	if req.Backend != actual.Backend && (req.Backend == "" || actual.Backend == "") {
		return true
	}
	if req.ParamsBackend != actual.ParamsBackend && (req.ParamsBackend == "" || actual.ParamsBackend == "") {
		return true
	}
	if req.OffloadToCPU != actual.OffloadToCPU {
		return true
	}
	if req.VAETiling != actual.VAETiling || req.TemporalTiling != actual.TemporalTiling {
		return true
	}
	if req.FlashAttentionMode != "" && req.FlashAttentionMode != actual.FlashAttentionMode {
		return true
	}
	if req.Circular != actual.Circular || req.MMap != actual.MMap {
		return true
	}
	if req.CacheMode != actual.CacheMode || req.CacheOption != actual.CacheOption {
		return true
	}
	if req.ConditioningCacheSize != actual.ConditioningCacheSize {
		return true
	}
	// A ControlNet is a startup slot, and the DGL allowlist is fixed when the
	// process starts. A changed path must restart the process or the workflow
	// would run stage 3 with no pose adapter (or stage 2 outside the allowlist).
	if req.ControlNetPath != "" && req.ControlNetPath != actual.ControlNetPath {
		return true
	}
	if req.UpscalerDir != "" && req.UpscalerDir != actual.UpscalerDir {
		return true
	}
	// --max-vram is a startup flag too. A changed cap is the user's decision and
	// must reach a new process instead of sitting inert on the resident one.
	if !imageIntPtrEqual(req.MaxVRAM, actual.MaxVRAM) {
		return true
	}
	if !imageIntPtrEqual(req.WeightResidencyTTL, actual.WeightResidencyTTL) {
		return true
	}
	return req.VAEPath != actual.VAEPath || req.ClipLPath != actual.ClipLPath ||
		req.ClipGPath != actual.ClipGPath || req.T5XXLPath != actual.T5XXLPath
}

func imageIntPtrEqual(a, b *int) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func imageDecodeExceedsHeadroom(est resource.AdmissionEstimate) bool {
	if est.TilingPlanned {
		return true
	}
	if est.ProjectedPeakVRAMMB > 0 && est.AvailableVRAMMB > 0 {
		return est.ProjectedPeakVRAMMB > est.AvailableVRAMMB
	}
	if est.ActivationVRAMMB <= 0 {
		return false
	}
	headroom := est.AvailableVRAMMB - est.RequiredVRAMMB
	if headroom < 0 {
		headroom = 0
	}
	return est.ActivationVRAMMB > headroom
}

func (a *App) warnImageHardware(message string) {
	fmt.Printf("[App] image hardware: %s\n", message)
	a.emitEvent("image:warning", message)
}

// announceImageAutoPlacement names the automatic trade once, on the request
// that actually starts or reloads the engine. A repeated generation at the same
// budget reuses the resident plan and does not nag again.
func (a *App) announceImageAutoPlacement(opts imagetypes.ModelLoadOptions) {
	budget := 0
	if opts.MaxVRAM != nil {
		budget = *opts.MaxVRAM
	}
	switch opts.AutoPlacementMode {
	case "cpu-diffusion":
		a.warnImageHardware(fmt.Sprintf(
			"VRAM Budget %d GiB selected automatic partial placement: diffusion and CLIP on CPU, VAE on CUDA and tiled. This honors the budget but is much slower.", budget))
	case "cpu-only":
		a.warnImageHardware(fmt.Sprintf(
			"VRAM Budget %d GiB selected automatic CPU-only placement: no GPU weights or compute. This honors the budget but is much slower.", budget))
	}
}

// imageLoadOptionsForGeneration is the one place that turns a generation
// request into the load options admission will see. The UI preflight and the
// backend start path both call it, so model identity, startup flags, companions
// and the GPU requirement cannot drift between the two.
func (a *App) imageLoadOptionsForGeneration(modelName string, gen imagetypes.ImageOptions) (imagetypes.ModelLoadOptions, error) {
	loadOpts := a.images.LoadOptions()
	loadOpts.VAETiling = gen.VAETiling
	loadOpts.TemporalTiling = gen.TemporalTiling
	loadOpts.FlashAttention = gen.FlashAttention
	loadOpts.FlashAttentionMode = gen.FlashAttentionMode
	loadOpts.OffloadToCPU = gen.OffloadToCPU
	// The generation panel is authoritative for these controls: an empty value
	// means Auto, not "keep whatever the previous request used". Conditional
	// assignment here made the UI show Auto while a resident engine still ran
	// the old backend.
	loadOpts.Backend = gen.Backend
	loadOpts.ParamsBackend = gen.ParamsBackend
	loadOpts.Circular = gen.Circular
	loadOpts.WeightResidencyTTL = gen.WeightResidencyTTL
	// Hardware/engine flags the generation panel sends. They reach the same
	// engine argument builder the explicit load path uses, so the preflight and
	// the spawn cannot disagree about them.
	loadOpts.CacheMode = gen.CacheMode
	loadOpts.CacheOption = gen.CacheOption
	loadOpts.WeightType = gen.WeightType
	loadOpts.ConditioningCacheSize = gen.ConditioningCacheSize
	loadOpts.MMap = gen.MMap
	// A runtime auto-cap belongs to the process that derived it, not to the
	// remembered panel state. Keeping the flag when the user supplies an
	// explicit cap made that cap look derived: the next blank Max VRAM copied the
	// old 3 GiB cap forward and admission refused a 6 GiB request.
	loadOpts.AutoMaxVRAM = false
	loadOpts.MaxVRAM = gen.MaxVRAM
	// Automatic placement is runtime state of one process at one budget. It is
	// not a saved user choice: carrying it into a different budget or model made
	// the 2 GiB request inherit the 3 GiB partial plan, skip re-planning, and get
	// refused by admission instead of falling back to CPU-only. Keep it only when
	// the same model is running at the same budget; otherwise derive it again.
	loaded := a.images.LoadedOptions()
	if !a.imageEngineLoadedSame(modelName) || !imageIntPtrEqual(loadOpts.MaxVRAM, loaded.MaxVRAM) {
		loadOpts.AutoPlacement = false
		loadOpts.AutoPlacementMode = ""
	}
	if gen.VAEPath != "" {
		loadOpts.VAEPath = gen.VAEPath
	}
	if gen.ClipLPath != "" {
		loadOpts.ClipLPath = gen.ClipLPath
	}
	if gen.ClipGPath != "" {
		loadOpts.ClipGPath = gen.ClipGPath
	}
	if gen.T5XXLPath != "" {
		loadOpts.T5XXLPath = gen.T5XXLPath
	}
	if gen.ControlNetPath != "" {
		loadOpts.ControlNetPath = gen.ControlNetPath
	}
	if gen.UpscalerDir != "" {
		loadOpts.UpscalerDir = gen.UpscalerDir
	} else if gen.UpscalerPath != "" {
		loadOpts.UpscalerDir = filepath.Dir(gen.UpscalerPath)
	}
	if modelName == "" {
		modelName = gen.ModelName
	}
	if modelName == "" {
		return loadOpts, fmt.Errorf("no model selected")
	}
	loadOpts.ModelName = modelName
	loadOpts.ModelPath = a.resolveModelPath(modelName)
	loadOpts.Family = a.registryFamily(modelName)
	// Engine profile follows this request, never the resident process or the
	// saved last-model record. An empty profile means the pinned default.
	loadOpts.EngineProfile = imagetypes.EngineProfileDefault
	if gen.EngineProfile != "" {
		loadOpts.EngineProfile = gen.EngineProfile
	}
	if gen.IPAdapterPath != "" || gen.IPAdapterImagePath != "" {
		loadOpts.EngineProfile = imagetypes.EngineProfileIPAdapter
	}
	// IP slots are request intent, exactly like the engine profile. An empty
	// selection means no IP-Adapter; carrying the resident adapter forward made
	// the next plain generation demand the IP profile and fail preflight.
	loadOpts.IPAdapterPath = gen.IPAdapterPath
	loadOpts.ClipVisionPath = gen.ClipVisionPath
	loadOpts.ClipVisionCPU = gen.ClipVisionCPU
	a.autofillCompanionModels(&loadOpts)
	if err := a.validateIPAdapterOptions(loadOpts); err != nil {
		return loadOpts, err
	}
	a.planIPAdapterPlacement(&loadOpts, gen)
	if err := a.validateModelLoadOptions(loadOpts); err != nil {
		return loadOpts, err
	}
	if err := a.planConstrainedOffload(&loadOpts, gen); err != nil {
		return loadOpts, err
	}
	return loadOpts, nil
}

// validateIPAdapterOptions keeps an IP request from silently degrading into a
// plain prompt. Adapter and vision are both required, and the adapter family
// must be one this integration implements.
func (a *App) validateIPAdapterOptions(opts imagetypes.ModelLoadOptions) error {
	if opts.IPAdapterPath == "" && opts.ClipVisionPath == "" {
		return nil
	}
	if opts.EngineProfile != imagetypes.EngineProfileIPAdapter {
		return fmt.Errorf("IP-Adapter requires the native-IP engine profile; the selected profile does not implement it")
	}
	if opts.IPAdapterPath == "" {
		return fmt.Errorf("IP-Adapter needs an adapter file")
	}
	if opts.ClipVisionPath == "" {
		return fmt.Errorf("IP-Adapter needs a CLIP-Vision ViT-H encoder; select clip_vision_h.safetensors")
	}
	for name, path := range map[string]string{"IP-Adapter": opts.IPAdapterPath, "CLIP-Vision": opts.ClipVisionPath} {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("%s is not readable: %w", name, err)
		}
		if info.IsDir() {
			return fmt.Errorf("%s is a directory, not a weight file: %s", name, path)
		}
	}
	adapterFamily := a.imageAssetFamily(opts.IPAdapterPath)
	if !scanner.CompatibleIPAdapter(scanner.ModelFamily(opts.Family), adapterFamily) {
		return fmt.Errorf("IP-Adapter '%s' is %s, but '%s' is %s. Pick a %s adapter.",
			filepath.Base(opts.IPAdapterPath), adapterFamily, opts.ModelName, opts.Family, opts.Family)
	}
	visionFamily := a.imageAssetFamily(opts.ClipVisionPath)
	if !scanner.CompatibleClipVision(scanner.ModelFamily(opts.Family), visionFamily) {
		return fmt.Errorf("CLIP-Vision '%s' is %s, which does not match '%s' (%s)",
			filepath.Base(opts.ClipVisionPath), visionFamily, opts.ModelName, opts.Family)
	}
	return nil
}

// imageAssetFamily returns the scanned family for an IP-Adapter or CLIP-Vision
// path. InvokeAI exports live under UUID directories named model.safetensors, so
// the filename alone cannot say whether an adapter is SD1.5 or SDXL.
func (a *App) imageAssetFamily(path string) scanner.ModelFamily {
	want := filepath.Clean(path)
	if a.registries.Models() != nil {
		idx := a.registries.Models().ToModelIndexForKind(registry.KindImage)
		for _, e := range append(append([]scanner.AssetEntry{}, idx.IPAdapters...), idx.ClipVisions...) {
			if filepath.Clean(e.Path) == want && e.Family != scanner.FamilyUnknown {
				return e.Family
			}
		}
	}
	return scanner.GuessFamily(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
}

// planIPAdapterPlacement keeps CLIP-Vision on the GPU when the whole stage fits
// under the free-VRAM reserve and moves only that encoder to CPU when it does
// not. The adapter stays a GPU adapter; if the stage still cannot fit, admission
// refuses instead of silently changing the request.
func (a *App) planIPAdapterPlacement(opts *imagetypes.ModelLoadOptions, gen imagetypes.ImageOptions) {
	if opts == nil || opts.IPAdapterPath == "" || opts.ClipVisionCPU {
		return
	}
	snap, err := resource.GetManager().RefreshSnapshot()
	if err != nil {
		return // admission reports the live-memory refusal
	}
	if snap.FreeVRAMBytes <= 0 || snap.TotalVRAMBytes <= 0 {
		return
	}
	req := a.imageAdmissionRequest(*opts, gen, false)
	req.ClipVisionCPU = false
	// The decision needs the converted total resident set, not the largest
	// module: an IP-Adapter and CLIP-Vision are additive.
	weights := resource.ImageResidentVRAMMB(req)
	_, activation, ok := resource.ImageResidentNeedMB(req)
	if !ok {
		return
	}
	freeMB := snap.FreeVRAMBytes / (1024 * 1024)
	reserveMB := int64(float64(snap.TotalVRAMBytes/(1024*1024)) * resource.DefaultReservePercent)
	const fixedOverheadMB = 512
	if weights+activation+fixedOverheadMB > freeMB-reserveMB {
		opts.ClipVisionCPU = true
	}
}

func (a *App) imageAdmissionRequest(loadOpts imagetypes.ModelLoadOptions, gen imagetypes.ImageOptions, useLoaded bool) resource.AdmissionRequest {
	// A generation against the model that is already running must judge the
	// precision and placement that process actually has, not the user's saved
	// request. Without this, an automatic q8_0 load was credited at q8 but charged
	// again at f16, and the next 512x512 generation was refused by 6 MB.
	if useLoaded {
		actual := a.images.LoadedOptions()
		loadOpts.WeightType = actual.WeightType
		loadOpts.Backend = actual.Backend
		loadOpts.ParamsBackend = actual.ParamsBackend
		loadOpts.OffloadToCPU = actual.OffloadToCPU
		loadOpts.VAETiling = actual.VAETiling
		loadOpts.ControlNetCPU = actual.ControlNetCPU
	}
	// Weights sd-server will hold: the checkpoint it opens plus every text
	// encoder loaded from a separate file. These are measured file bytes, the
	// same files the load options above hand to the engine.
	checkpointBytes := a.fileSizeBytes(loadOpts.ModelPath)
	conditionerBytes := a.fileSizeBytes(loadOpts.ClipLPath) + a.fileSizeBytes(loadOpts.ClipGPath) +
		a.fileSizeBytes(loadOpts.T5XXLPath) + a.fileSizeBytes(loadOpts.LLMPath)
	if loadOpts.DiffusionModel != "" && filepath.Clean(loadOpts.DiffusionModel) != filepath.Clean(loadOpts.ModelPath) {
		conditionerBytes += a.fileSizeBytes(loadOpts.DiffusionModel)
	}
	// One ControlNet slot exists per sd-server process. Charging the request path
	// and the model-bound path separately double-counted the same 4.7 GB adapter
	// whenever both named it. Prefer the explicit startup/request path and fall
	// back to the model bundle.
	controlNetPath := loadOpts.ControlNetPath
	if controlNetPath == "" {
		controlNetPath = gen.ControlNetPath
	}
	controlNetBytes := a.boundControlNetBytes(loadOpts.ModelPath)
	if controlNetPath != "" {
		controlNetBytes = a.fileSizeBytes(controlNetPath)
	}
	vaeBytes := a.fileSizeBytes(loadOpts.VAEPath)
	if vaeBytes == 0 {
		vaeBytes = a.fileSizeBytes(gen.VAEPath)
	}
	// A full safetensors checkpoint carries its conditioner and VAE inside the
	// same file. sd.cpp places those modules independently, so charging the whole
	// file to the diffusion module overstates the largest module by gigabytes.
	// The header reader is exact and reads no tensor payload; formats it cannot
	// split (GGUF) keep the whole-file fallback.
	if parts, ok := sdcpp.ReadComponentBytes(loadOpts.ModelPath); ok {
		checkpointBytes = parts.Diffusion + parts.Other
		conditionerBytes += parts.Conditioner
		vaeBytes += parts.VAE
	}
	cpuOnly := imageCPUOnly(loadOpts)
	weightsOffloaded := imageWeightsOffloaded(loadOpts)
	// The 10% free-VRAM reserve is a product requirement, not an option that a
	// --max-vram cap can trade away. A positive --max-vram only limits how much
	// the engine may spend; admission still keeps the reserve free. The request
	// therefore leaves ReserveOverrideMB at zero and lets the authority apply
	// DefaultReservePercent.
	ipAdapterBytes := a.fileSizeBytes(loadOpts.IPAdapterPath)
	if ipAdapterBytes == 0 {
		ipAdapterBytes = a.fileSizeBytes(gen.IPAdapterPath)
	}
	clipVisionBytes := a.fileSizeBytes(loadOpts.ClipVisionPath)
	if clipVisionBytes == 0 {
		clipVisionBytes = a.fileSizeBytes(gen.ClipVisionPath)
	}
	req := resource.AdmissionRequest{
		Kind:             resource.ResourceImage,
		Name:             loadOpts.ModelName,
		Path:             loadOpts.ModelPath,
		Width:            gen.Width,
		Height:           gen.Height,
		BatchCount:       gen.BatchCount,
		Steps:            gen.Steps,
		Family:           loadOpts.Family,
		Mode:             gen.Mode,
		CheckpointBytes:  checkpointBytes,
		ConditionerBytes: conditionerBytes,
		ControlNetBytes:  controlNetBytes,
		IPAdapterBytes:   ipAdapterBytes,
		ClipVisionBytes:  clipVisionBytes,
		ClipVisionCPU:    loadOpts.ClipVisionCPU,
		LoRABytes:        a.imageLoRABytes(gen.LoRAs),
		VAEBytes:         vaeBytes,
		VAETiling:        loadOpts.VAETiling,
		ControlNetCPU:    loadOpts.ControlNetCPU,
		WeightType:       loadOpts.WeightType,
		RequireGPU:       !weightsOffloaded && !cpuOnly,
		// The engine passes --offload-to-cpu for the first flag and
		// --params-backend cpu/disk for the others. All keep the weight modules
		// out of VRAM, so admission must judge the placement that will happen.
		WeightsOffloaded: weightsOffloaded,
		CPUOnly:          cpuOnly,
		DeviceIndex:      0,
		// A positive --max-vram is a spend cap, not a reserve override.
		VRAMBudgetMB: imageVRAMBudgetMB(loadOpts.MaxVRAM),
	}
	// Credit a replaced model only when its process is actually running. A
	// registered model whose server was killed reserves no VRAM, and crediting it
	// let recovery choose an 8 GiB budget on an 8 GiB card.
	if eng := a.images.Engine(); eng != nil && eng.IsRunning() {
		if loaded := eng.GetLoadedModel(); loaded != nil && loaded.Name != "" {
			req.ReplacingKind = resource.ResourceImage
			req.ReplacingName = loaded.Name
		}
	}
	// Print the measured bytes, not just the verdict: a wrong companion set is
	// what makes an otherwise correct estimate huge.
	// A generation against the running model is a stage delta. The resident
	// weights were admitted when the process started; charging them again would
	// refuse a healthy resident engine under a budget that its measured placement
	// already fits (for example the second pass of a 5 GiB batch).
	req.ResidentLoaded = useLoaded
	if !useLoaded {
		req.MeasuredResidentVRAMMB = a.measuredImageResidentVRAMMB(loadOpts, gen)
		// A post-generation device delta for this same model and an equal or
		// larger workload shape is stronger evidence than the metadata
		// projection. It raises the commitment; it never replaces the resident
		// requirement, because that delta already contains the activation and
		// allocator retention of the measured run.
		batch := gen.BatchCount
		if batch <= 0 {
			batch = 1
		}
		req.MinimumCommitmentVRAMMB = a.observedCommitmentFloorMB(resource.ResourceImage, loadOpts.ModelName, gen.Width, gen.Height, batch, loadOpts.WeightType)
	}
	fmt.Printf("[App] image request: checkpoint=%dMB conditioner=%dMB vae=%dMB control=%dMB ip_adapter=%dMB clip_vision=%dMB cpu=%t lora=%dMB measured_resident=%dMB observed_floor=%dMB model=%q\n",
		req.CheckpointBytes>>20, req.ConditionerBytes>>20, req.VAEBytes>>20, req.ControlNetBytes>>20, req.IPAdapterBytes>>20, req.ClipVisionBytes>>20, req.ClipVisionCPU, req.LoRABytes>>20, req.MeasuredResidentVRAMMB, req.MinimumCommitmentVRAMMB, req.Name)
	return req
}

// imageVRAMBudgetMB is a positive --max-vram cap in MB; 0 means no explicit cap
// beyond the device's current free memory.
func imageVRAMBudgetMB(maxVRAM *int) int64 {
	if maxVRAM == nil || *maxVRAM <= 0 {
		return 0
	}
	return int64(*maxVRAM) * 1024
}

// invalidateImageRegistrationsIfStopped removes image registrations only when
// the engine confirms no live process owns them. A refused or failed load that
// left the previous child running keeps its registration.
func (a *App) invalidateImageRegistrationsIfStopped() {
	if a.images == nil || a.images.Engine() == nil || !a.images.Engine().IsRunning() {
		a.unregisterImageModels()
	}
}

// registerImageModelWithObservation records what the engine and the device
// reported, then registers the conservative figure. The admitted commitment
// (projected peak) is the planning bound; RequiredVRAMMB alone understated a
// running engine by its activation and overhead.
func (a *App) registerImageModelWithObservation(loadOpts imagetypes.ModelLoadOptions, decision resource.AdmissionDecision) {
	plannedMB := decision.Estimate.CommittedVRAMMB()
	vramMB := plannedMB
	var engineReportedMB int64
	if eng, ok := a.images.Engine().(*sdcpp.Engine); ok && eng != nil {
		if placement := eng.GetPlacement(); placement.Known {
			engineReportedMB = int64(placement.VRAMMB)
		}
	}
	beforeMB := decision.Estimate.TotalVRAMMB - decision.Estimate.FreeVRAMMB
	// sd-server loads weights lazily and reports placement only once a
	// generation has run, so the honest observation point is post-generation.
	// Record the baseline here and do not invent a zero-footprint observation.
	a.noteResidencyBaseline(resource.ResourceImage, plannedMB, beforeMB)
	if engineReportedMB > 0 {
		const mb = int64(1024 * 1024)
		if snap, snapErr := resource.GetManager().RefreshSnapshot(); snapErr == nil {
			afterMB := (snap.TotalVRAMBytes - snap.FreeVRAMBytes) / mb
			observation := resource.NewResidencyObservation(resource.ResourceImage, loadOpts.ModelName, plannedMB, engineReportedMB, afterMB-beforeMB)
			vramMB = observation.LedgerVRAMMB()
			a.recordResidency(observation)
		}
	}
	a.registerImageModel(loadOpts, decision, vramMB)
}

func (a *App) registerImageModel(loadOpts imagetypes.ModelLoadOptions, decision resource.AdmissionDecision, vramMB int64) {
	a.images.SetLoadedOptions(loadOpts)
	a.unregisterImageModels()
	name := loadOpts.ModelName
	if name == "" {
		name = filepath.Base(loadOpts.ModelPath)
	}
	resource.GetManager().RegisterModel(
		resource.ResourceImage,
		name,
		loadOpts.ModelPath,
		vramMB,
		decision.Estimate.RequiredRAMMB,
		resource.PriorityImage,
		vramMB > 0,
	)
}

func (a *App) unregisterImageModels() {
	for _, m := range resource.GetManager().ModelsByType(resource.ResourceImage) {
		resource.GetManager().UnregisterModel(resource.ResourceImage, m.Name)
	}
}

func (a *App) fileSizeBytes(path string) int64 {
	if path == "" {
		return 0
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return 0
	}
	return info.Size()
}

func (a *App) imageLoRABytes(loras []imagetypes.LoRAStrength) int64 {
	var total int64
	for _, lora := range loras {
		if !lora.Enabled || lora.Strength == 0 || lora.Name == "" {
			continue
		}
		path := lora.Name
		if !filepath.IsAbs(path) {
			if resolved := a.resolveModelPath(path); resolved != "" {
				path = resolved
			} else if dir := a.imageModelsDir(); dir != "" {
				path = filepath.Join(dir, path)
			}
		}
		total += a.fileSizeBytes(path)
	}
	return total
}

// parseModelLoadOpts converts the loosely-typed Wails opts into ModelLoadOptions.
func parseModelLoadOpts(opts interface{}) (imagetypes.ModelLoadOptions, error) {
	var out imagetypes.ModelLoadOptions
	switch v := opts.(type) {
	case string:
		out.ModelName = v
	case map[string]interface{}:
		if m, ok := v["model"].(string); ok {
			out.ModelName = m
		}
		if s, ok := v["vae"].(string); ok {
			out.VAEPath = s
		}
		if s, ok := v["clip_l"].(string); ok {
			out.ClipLPath = s
		}
		if s, ok := v["clip_g"].(string); ok {
			out.ClipGPath = s
		}
		if s, ok := v["t5xxl"].(string); ok {
			out.T5XXLPath = s
		}
		if s, ok := v["diffusion_model"].(string); ok {
			out.DiffusionModel = s
		}
		if s, ok := v["backend"].(string); ok {
			out.Backend = s
		}
		if s, ok := v["params_backend"].(string); ok {
			out.ParamsBackend = s
		}
		if b, ok := v["vae_tiling"].(bool); ok {
			out.VAETiling = b
		}
		if b, ok := v["temporal_tiling"].(bool); ok {
			out.TemporalTiling = b
		}
		if b, ok := v["circular"].(bool); ok {
			out.Circular = b
		}
		if s, ok := v["flash_attention_mode"].(string); ok {
			out.FlashAttentionMode = s
		}
		if b, ok := v["flash_attention"].(bool); ok {
			out.FlashAttention = b
			if out.FlashAttentionMode == "" {
				if b {
					out.FlashAttentionMode = "on"
				} else {
					out.FlashAttentionMode = "auto"
				}
			}
		}
		if b, ok := v["offload_to_cpu"].(bool); ok {
			out.OffloadToCPU = b
		}
		if b, ok := v["disable_metadata"].(bool); ok {
			out.DisableMetadata = b
		}
		if s, ok := v["cache_mode"].(string); ok {
			out.CacheMode = s
		}
		if s, ok := v["cache_option"].(string); ok {
			out.CacheOption = s
		}
		if f, ok := v["max_vram"].(float64); ok {
			n := int(f)
			out.MaxVRAM = &n
		}
		if s, ok := v["weight_type"].(string); ok {
			out.WeightType = s
		}
		if f, ok := v["conditioning_cache_size"].(float64); ok {
			out.ConditioningCacheSize = int(f)
		}
		if f, ok := v["weight_residency_ttl"].(float64); ok {
			n := int(f)
			out.WeightResidencyTTL = &n
		}
		if b, ok := v["mmap"].(bool); ok {
			out.MMap = b
		}
	default:
		return out, fmt.Errorf("invalid options type: %T", opts)
	}
	return out, nil
}

// -----------------------------------------------------------------------
// Model discovery
// -----------------------------------------------------------------------

// ListImageModels merges the engine scan with registry entries.
func (a *App) ListImageModels() ([]imagetypes.ModelInfo, error) {
	var models []imagetypes.ModelInfo
	if a.images.Engine() != nil {
		m, err := a.images.Engine().ListModels()
		if err == nil && m != nil {
			models = m
		}
	}

	if a.registries.Models() != nil {
		seen := make(map[string]struct{}, len(models))
		for _, m := range models {
			seen[m.Path] = struct{}{}
		}
		for _, e := range a.registries.Models().All() {
			if e.Type != "diffusion" {
				continue
			}
			// Menu boundary: only models living in an image directory.
			if a.registries.Models().KindForPath(e.Path) != registry.KindImage {
				continue
			}
			if _, ok := seen[e.Path]; ok {
				continue
			}
			files := e.Files
			if files == nil {
				files = map[string]string{"model": e.Path}
			}
			models = append(models, imagetypes.ModelInfo{
				Name:           e.Name,
				Path:           e.Path,
				Family:         e.Family,
				Capabilities:   []string{"txt2img", "img2img", "inpaint"},
				ResolutionHint: resolutionHintFor(e.Family),
				Files:          files,
				Format:         e.Format,
				Compatibility:  assetCompatForInfo(e.Family, e.Format),
			})
		}
	}

	if models == nil {
		return []imagetypes.ModelInfo{}, nil
	}
	return models, nil
}

// ListImageLoRAs scans the models directory for LoRA files.
func (a *App) ListImageLoRAs() ([]string, error) {
	if a.images.Engine() == nil {
		return []string{}, fmt.Errorf("image engine not initialized")
	}
	loras, err := a.images.Engine().ListLoRAs()
	if err != nil {
		return []string{}, err
	}
	if loras == nil {
		return []string{}, nil
	}
	return loras, nil
}

// resolveModelPath looks up the absolute file path for a model name from the registry.
func (a *App) resolveModelPath(name string) string {
	if a.registries.Models() == nil || name == "" {
		return ""
	}
	var paths []string
	for _, e := range a.registries.Models().All() {
		if e.Name == name && e.Path != "" {
			paths = append(paths, e.Path)
		}
	}
	return preferModelPath(paths, a.imageModelsDir())
}

// imageModelsDir is the configured image models directory, cleaned for prefix
// comparison. Empty when nothing is configured.
func (a *App) imageModelsDir() string {
	cfg := a.GetImageEngineConfig()
	if dir := strings.TrimSpace(cfg.ModelsDir); dir != "" {
		return filepath.Clean(dir)
	}
	return ""
}

// preferModelPath picks one path when several registry entries share a name.
//
// Names are not unique in practice: an InvokeAI tree and the G: models folder
// both hold krea2_turbo-Q4_K_M. Returning the first match picked the InvokeAI
// copy, sd-server answered "model metadata validation failed", and the load
// burned its full ten-minute timeout without ever trying the file that works.
// Selection is by usefulness: a real file in the configured models directory,
// then any real file, then a models-directory path that has gone missing, then
// path order so the answer is stable.
func preferModelPath(paths []string, modelsDir string) string {
	best := ""
	bestRank := 4
	for _, p := range paths {
		clean := filepath.Clean(p)
		st, err := os.Stat(clean)
		exists := err == nil && !st.IsDir()
		inModelsDir := modelsDir != "" && strings.HasPrefix(strings.ToLower(clean), strings.ToLower(modelsDir))
		rank := 3
		switch {
		case exists && inModelsDir:
			rank = 0
		case exists:
			rank = 1
		case inModelsDir:
			rank = 2
		}
		if rank < bestRank || (rank == bestRank && (best == "" || clean < best)) {
			best, bestRank = clean, rank
		}
	}
	return best
}

func (a *App) modelEntryByName(name string) *registry.Entry {
	if a.registries.Models() == nil || name == "" {
		return nil
	}
	// Resolve the same path the loader will use, so the family and companions
	// always describe the file that is actually opened.
	want := a.resolveModelPath(name)
	for _, e := range a.registries.Models().All() {
		if e.Name != name || e.Type != "diffusion" || a.registries.Models().KindForPath(e.Path) != registry.KindImage {
			continue
		}
		if want != "" && filepath.Clean(e.Path) != filepath.Clean(want) {
			continue
		}
		c := e
		return &c
	}
	return nil
}

// registryFamily returns the architecture the image registry recorded for a
// model, which is the same value the dependency report and the validator use.
// Empty when the model is not registered.
func (a *App) registryFamily(modelName string) string {
	e := a.modelEntryByName(modelName)
	if e == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(e.Family))
}

func (a *App) autofillCompanionModels(opts *imagetypes.ModelLoadOptions) {
	if opts == nil {
		return
	}
	e := a.modelEntryByName(opts.ModelName)
	if e == nil {
		return
	}
	fam := strings.ToLower(e.Family)
	usesFluxAE := sdcpp.UsesFluxAE(fam)

	// Companions resolved for the previous model must not follow a new one in.
	// Loading flux1-schnell and then switching to an SDXL checkpoint carried
	// ae.safetensors across and the load failed with "VAE 'ae.safetensors' is for
	// flux models, but ... is sdxl" — the paths live on lastImageModelOpts, and
	// autofill only ever fills empty fields, so a stale one is never revisited.
	pruneStaleCompanions(opts, fam, a.registryFamilyForPath)

	// Registry-recorded companions are candidates, not decisions: they are
	// written at scan time, before anything has checked the file is the right
	// kind or still on disk. A wrong one falls through to the picker below.
	if opts.VAEPath == "" && usesFluxAE {
		opts.VAEPath = usableCompanion(fam, "vae", e.Files["vae"])
	}
	if opts.ClipLPath == "" {
		opts.ClipLPath = usableCompanion(fam, "clip_l", e.Files["clip_l"])
	}
	if opts.ClipGPath == "" {
		opts.ClipGPath = usableCompanion(fam, "clip_g", e.Files["clip_g"])
	}
	if opts.T5XXLPath == "" {
		opts.T5XXLPath = usableCompanion(fam, "t5xxl", e.Files["t5xxl"])
	}
	if opts.DiffusionModel == "" && e.Files != nil && e.Files["diffusion_model"] != "" {
		opts.DiffusionModel = e.Files["diffusion_model"]
	}
	if opts.ClipLPath == "" {
		opts.ClipLPath = a.findCompatibleCompanion(fam, "clip_l")
	}
	// clip_g had no fallback at all: only a checkpoint whose own manifest named
	// a clip_g could ever fill it, and a converted diffusion GGUF never does.
	// SD3.5 therefore reported "Generation cancelled — the model is missing
	// required files" with clip_g.safetensors sitting in the models folder.
	if opts.ClipGPath == "" {
		opts.ClipGPath = a.findCompatibleCompanion(fam, "clip_g")
	}
	if opts.T5XXLPath == "" {
		opts.T5XXLPath = a.findCompatibleCompanion(fam, "t5xxl")
	}
	// Families that encode text with an LLM (z_image, qwen_image, flux2, anima)
	// were passed to sd-server with no --llm at all: the engine only emits the
	// flag when the option is set, and nothing ever set it here.
	if opts.LLMPath == "" {
		opts.LLMPath = usableCompanion(fam, "llm", e.Files["llm"])
	}
	if opts.LLMPath == "" {
		opts.LLMPath = a.findCompatibleCompanion(fam, "llm")
	}
	if opts.VAEPath == "" && usesFluxAE {
		opts.VAEPath = a.findCompatibleCompanion(fam, "vae")
	}
	// Same gap for the VAE: only Flux-style families were auto-filled, so a
	// z_image or qwen_image checkpoint loaded with no decoder.
	if opts.VAEPath == "" && !usesFluxAE {
		if cap, ok := sdcpp.Capability(fam); ok && cap.Uses("vae") {
			opts.VAEPath = a.findCompatibleCompanion(fam, "vae")
		}
	}

	// A checkpoint that carries a role inside itself must not be handed a
	// borrowed file for that role. Every full SDXL .safetensors was given the
	// index's CLIP-L and a separate VAE, so sd-server loaded the file's UNet
	// without its CLIP-G and 768x768 renders came back as coloured noise. This
	// runs last: the fills above are the ones that invent those paths.
	checkpoint := e.Files["model"]
	if checkpoint == "" {
		checkpoint = e.Path
	}
	clearEmbeddedCompanions(opts, checkpoint)
	pruneUnsupportedCompanions(opts, fam)
}

// pruneUnsupportedCompanions drops companion paths for roles the family does not
// declare. Autofill (and a previous model's remembered options) can attach an
// SD3/Flux text encoder or LLM to an SDXL checkpoint: that both overstates the
// admission requirement by gigabytes and hands sd-server files the graph never
// reads. Roles a capability lists as required or optional are kept.
func pruneUnsupportedCompanions(opts *imagetypes.ModelLoadOptions, family string) {
	if opts == nil {
		return
	}
	cap, ok := sdcpp.Capability(family)
	if !ok {
		return
	}
	roles := []struct {
		role string
		path *string
	}{
		{"vae", &opts.VAEPath},
		{"clip_l", &opts.ClipLPath},
		{"clip_g", &opts.ClipGPath},
		{"t5xxl", &opts.T5XXLPath},
		{"llm", &opts.LLMPath},
	}
	for _, r := range roles {
		if !cap.Uses(r.role) {
			*r.path = ""
		}
	}
}

// clearEmbeddedCompanions drops companion paths for roles the checkpoint file
// already carries inside itself. A full SDXL .safetensors holds its VAE and both
// text encoders; handing sd-server a separate CLIP-L and VAE instead replaced
// them, left CLIP-G unloaded, and every render came back as coloured noise.
func clearEmbeddedCompanions(opts *imagetypes.ModelLoadOptions, checkpoint string) {
	if opts == nil || checkpoint == "" {
		return
	}
	embedded := sdcpp.EmbeddedRoles(checkpoint)
	if embedded["vae"] {
		opts.VAEPath = ""
	}
	if embedded["clip_l"] {
		opts.ClipLPath = ""
	}
	if embedded["clip_g"] {
		opts.ClipGPath = ""
	}
}

// usableCompanion returns a registry-provided companion path only when it is the
// right kind of file for the role and still exists on disk.
//
// The registry records a model's companions at scan time. Nothing checked them
// then, so flux1-schnell carried
// H:\InvokeAI\...\text_encoder_2\bnb_llm_int8_model.safetensors as its T5 and
// sd-server aborted on GGML_ASSERT(weight_scale->type == GGML_TYPE_F32). A path
// that fails the catalogue check or has since moved is rejected here, which lets
// autofill fall through to the picker instead of into a failed load.
func usableCompanion(family, role, path string) string {
	p := strings.TrimSpace(path)
	if p == "" {
		return ""
	}
	if ok, _ := deps.Compatible(family, role, filepath.Base(p)); !ok {
		return ""
	}
	if st, err := os.Stat(p); err != nil || st.IsDir() {
		return ""
	}
	return p
}

// registryFamilyForPath returns the family the image registry recorded for a
// file. Empty when the file was never scanned, which means it cannot be judged
// and is left alone.
func (a *App) registryFamilyForPath(path string) string {
	if a.registries.Models() == nil || path == "" {
		return ""
	}
	want := filepath.Clean(path)
	for _, e := range a.registries.Models().All() {
		if filepath.Clean(e.Path) == want {
			return strings.ToLower(strings.TrimSpace(e.Family))
		}
	}
	return ""
}

// pruneStaleCompanions clears companion paths that do not belong to the family
// about to be loaded.
//
// It is a pure function so the rule can be tested without an App: familyFor may
// return "" for a file the registry does not know, and an unknown file is kept
// rather than guessed at. Two rules decide: the catalogue's identity check, and
// the file's own recorded family (an SDXL checkpoint cannot decode Flux latents
// no matter how valid the file is).
func pruneStaleCompanions(opts *imagetypes.ModelLoadOptions, family string, familyFor func(string) string) {
	if opts == nil {
		return
	}
	target := strings.ToLower(strings.TrimSpace(family))
	if target == "" {
		return
	}
	roles := []struct {
		role string
		path *string
	}{
		{"vae", &opts.VAEPath},
		{"clip_l", &opts.ClipLPath},
		{"clip_g", &opts.ClipGPath},
		{"t5xxl", &opts.T5XXLPath},
		{"llm", &opts.LLMPath},
	}
	for _, r := range roles {
		p := strings.TrimSpace(*r.path)
		if p == "" {
			continue
		}
		if ok, _ := deps.Compatible(target, r.role, filepath.Base(p)); !ok {
			*r.path = ""
			continue
		}
		if f := familyFor(p); f != "" && !companionFamilyUsable(target, f) {
			*r.path = ""
		}
	}
}

func (a *App) validateModelLoadOptions(opts imagetypes.ModelLoadOptions) error {
	e := a.modelEntryByName(opts.ModelName)
	if e == nil {
		return nil
	}
	fam := strings.ToLower(strings.TrimSpace(e.Family))

	// A file we could not identify is not a checkpoint sd-server can build a
	// graph for. The message names what this build implements, taken from the
	// capability table instead of a hand-written list that drifts.
	if fam == "" || fam == "unknown" {
		return fmt.Errorf("'%s' could not be identified as a diffusion architecture, so it cannot be loaded. This build supports: %s.",
			opts.ModelName, strings.Join(sdcpp.SupportedFamilies(), ", "))
	}
	cap, known := sdcpp.Capability(fam)
	if !known {
		return fmt.Errorf("'%s' is %s, which this stable-diffusion.cpp build does not implement. Supported: %s.",
			opts.ModelName, e.Family, strings.Join(sdcpp.SupportedFamilies(), ", "))
	}
	if missing := missingCompanions(cap, opts); len(missing) > 0 {
		return fmt.Errorf("%s model '%s' requires %s. Select the %s model first.",
			e.Family, opts.ModelName, strings.Join(missing, ", "), missing[0])
	}

	// Companion files must belong to the same family as the checkpoint. A SD1.5
	// VAE under an SDXL model decodes garbage, and the failure surfaces only
	// after a full sampling run.
	if err := a.validateCompanionFamilies(fam, opts); err != nil {
		return err
	}
	return nil
}

// companionFamilyUsable reports whether a companion file may be used under the
// target family: the same architecture, a file the scanner could not classify,
// or a family that shares the autoencoder (Z-Image and Flux both decode with
// ae.safetensors). One rule, because the picker and the validator disagreeing is
// what produced "VAE 'model.safetensors' is for sd1.5 models" on a model the UI
// had just reported as fully equipped.
func companionFamilyUsable(targetFamily, fileFamily string) bool {
	target := strings.ToLower(strings.TrimSpace(targetFamily))
	file := strings.ToLower(strings.TrimSpace(fileFamily))
	if file == "" || file == "unknown" || file == target {
		return true
	}
	return sdcpp.SharesAutoencoder(target, file)
}

// validateCompanionFamilies rejects a VAE picked from a different family than
// the loaded checkpoint. Family data comes from the image registry, so files
// that were never scanned (manual paths) are left alone rather than guessed at.
func (a *App) validateCompanionFamilies(fam string, opts imagetypes.ModelLoadOptions) error {
	if a.registries.Models() == nil || opts.VAEPath == "" {
		return nil
	}
	idx := a.registries.Models().ToModelIndexForKind(registry.KindImage)
	want := filepath.Clean(opts.VAEPath)
	for _, v := range idx.VAEs {
		if filepath.Clean(v.Path) != want {
			continue
		}
		vfam := strings.ToLower(string(v.Family))
		if companionFamilyUsable(fam, vfam) {
			return nil
		}
		return fmt.Errorf("VAE '%s' is for %s models, but '%s' is %s. Pick a %s VAE or leave VAE empty to use the built-in one.",
			filepath.Base(v.Path), vfam, opts.ModelName, fam, fam)
	}
	return nil
}

// -----------------------------------------------------------------------
// Shared helpers (also used by app-image-scan.go)
// -----------------------------------------------------------------------

// assetCompatForInfo returns a compatibility string for a given family/format pair.
func assetCompatForInfo(family, format string) string {
	switch format {
	case "safetensors", "ckpt", "sft":
		return "verified"
	case "gguf", "bundle", "diffusers":
		return "supported"
	}
	return "unknown"
}

// resolutionHintFor returns the default resolution string for a model family.
func resolutionHintFor(family string) string {
	// Base resolution comes from the capability table: the UI must not keep its
	// own copy of which architecture is trained at 1024.
	if res := sdcpp.FamilyResolution(family); res > 0 {
		return fmt.Sprintf("%dx%d", res, res)
	}
	return "512x512"
}

// companionPath maps a required role onto the load option that carries it.
func companionPath(role string, o imagetypes.ModelLoadOptions) string {
	switch role {
	case "vae":
		return o.VAEPath
	case "clip_l":
		return o.ClipLPath
	case "clip_g":
		return o.ClipGPath
	case "t5xxl":
		return o.T5XXLPath
	case "llm":
		return o.LLMPath
	}
	return ""
}

// missingCompanions lists the required roles the caller has not supplied.
// The requirement lives in the shared capability table, so the engine and this
// check cannot disagree about what a family needs.
func missingCompanions(c sdcpp.FamilyCapability, o imagetypes.ModelLoadOptions) []string {
	var out []string
	for _, role := range c.Companions {
		if companionPath(role, o) == "" {
			out = append(out, role)
		}
	}
	return out
}
