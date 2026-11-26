package main

// Sequential image workflows.
//
// The engine has no workflow runner: native txt2img/img2img are HTTP requests,
// and DGL is a separate job surface. Howl therefore owns the sequence. Each
// stage is one operation, waits for its predecessor, obtains its own admission
// decision, and reports the placement the engine actually chose. The workflow
// holds imageStartMu for the whole sequence so a queued job cannot interleave a
// load or a reclamation between stages.
//
// Reclamation stays with app-image-reclaim.go. The workflow only calls its
// locked helper once at the terminal outcome, exactly like a single generation.

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/image/engine/sdcpp"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/resource"
)

// ImageWorkflowOptions is one sequential pipeline request.
type ImageWorkflowOptions struct {
	ModelName      string  `json:"model_name"`
	Prompt         string  `json:"prompt"`
	NegativePrompt string  `json:"negative_prompt"`
	Seed           int64   `json:"seed"`
	Sampler        string  `json:"sampler"`
	Scheduler      string  `json:"scheduler"`
	CfgScale       float64 `json:"cfg_scale"`

	BaseWidth  int `json:"base_width"`
	BaseHeight int `json:"base_height"`
	BaseSteps  int `json:"base_steps"`

	UpscalerPath string `json:"upscaler_path"`
	UpscalerTile int    `json:"upscaler_tile"`

	ControlNetPath  string  `json:"control_net_path"`
	PoseImagePath   string  `json:"pose_image_path"`
	ControlStrength float64 `json:"control_strength"`

	// IP-Adapter conditioning. When IPAdapterPath is set the workflow runs on
	// the native-IP profile: native txt2img, native upscale, native ControlNet
	// img2img, native upscale. It does not mix in the DGL upscaler.
	IPAdapterPath      string  `json:"ip_adapter_path"`
	ClipVisionPath     string  `json:"clip_vision_path"`
	IPAdapterImagePath string  `json:"ip_adapter_image_path"`
	IPAdapterStrength  float64 `json:"ip_adapter_strength"`
	ClipVisionCPU      bool    `json:"clip_vision_cpu"`
	EngineProfile      string  `json:"engine_profile"`

	RefineWidth    int     `json:"refine_width"`
	RefineHeight   int     `json:"refine_height"`
	RefineSteps    int     `json:"refine_steps"`
	RefineStrength float64 `json:"refine_strength"`

	BatchCount int `json:"batch_count"`
	// MaxVRAM is the panel's VRAM budget in GiB. The workflow forwards it to
	// the same startup path as a single generation; otherwise a constrained
	// pipeline would silently run on the resident engine's old budget.
	MaxVRAM *int `json:"max_vram,omitempty"`
}

// ImageWorkflowStage is the evidence for one executed stage.
type ImageWorkflowStage struct {
	Index           int                    `json:"index"`
	Name            string                 `json:"name"`
	Status          string                 `json:"status"` // running, completed, failed
	Executor        string                 `json:"executor"`
	StartedAt       string                 `json:"started_at"`
	CompletedAt     string                 `json:"completed_at"`
	DurationMS      int64                  `json:"duration_ms"`
	Steps           int                    `json:"steps,omitempty"`
	Admission       string                 `json:"admission,omitempty"`
	Placement       map[string]interface{} `json:"placement,omitempty"`
	OutputPath      string                 `json:"output_path,omitempty"`
	OutputWidth     int                    `json:"output_width,omitempty"`
	OutputHeight    int                    `json:"output_height,omitempty"`
	JobID           string                 `json:"job_id,omitempty"`
	Error           string                 `json:"error,omitempty"`
	ResourcesBefore map[string]interface{} `json:"resources_before,omitempty"`
	ResourcesAfter  map[string]interface{} `json:"resources_after,omitempty"`
}

// ImageWorkflowResult is returned even when a stage fails: the stages already
// completed are the evidence of what ran, and the error carries the failure.
type ImageWorkflowResult struct {
	Stages         []ImageWorkflowStage   `json:"stages"`
	FinalImagePath string                 `json:"final_image_path,omitempty"`
	FinalImageB64  string                 `json:"final_image_b64,omitempty"`
	Width          int                    `json:"width,omitempty"`
	Height         int                    `json:"height,omitempty"`
	BatchRuns      []ImageWorkflowRun     `json:"batch_runs,omitempty"`
	Reclaim        map[string]interface{} `json:"reclaim,omitempty"`
	Error          string                 `json:"error,omitempty"`
}

// ImageWorkflowRun is one pipeline pass inside a batch.
type ImageWorkflowRun struct {
	Index          int                    `json:"index"`
	Stages         []ImageWorkflowStage   `json:"stages"`
	FinalImagePath string                 `json:"final_image_path,omitempty"`
	Reclaim        map[string]interface{} `json:"reclaim,omitempty"`
	Width          int                    `json:"width,omitempty"`
	Height         int                    `json:"height,omitempty"`
	Error          string                 `json:"error,omitempty"`
}

const (
	workflowStageBase    = "base generation"
	workflowStageUpscale = "upscale"
	workflowStagePose    = "pose/controlnet generation"
	workflowStageFinal   = "final upscale"
)

// RunImageWorkflow executes base txt2img, ESRGAN upscale, ControlNet img2img and
// a second ESRGAN upscale in that order. It is synchronous by design: the caller
// sees the stage evidence, and the workflow owns the image-start lock for the
// whole sequence.
func (a *App) RunImageWorkflow(opts ImageWorkflowOptions) (*ImageWorkflowResult, error) {
	if err := a.ensureAcceptingWork("image workflow"); err != nil {
		return nil, err
	}
	if a.images == nil || a.images.Engine() == nil {
		return nil, fmt.Errorf("image engine not initialized")
	}
	a.imageStartMu.Lock()
	defer a.imageStartMu.Unlock()

	run := newImageWorkflowExecution()
	a.setImageWorkflowRun(run)
	defer func() {
		run.complete()
		a.clearImageWorkflowRun(run)
	}()

	engine, _ := a.images.Engine().(*sdcpp.Engine)
	parent := a.beginImageWorkflow(engine)
	result, lastGen, runErr := a.runImageWorkflowLocked(run.ctx, run, opts)
	// No cancellation is accepted after the stage loop returns. A cancel that
	// already won the race is handled here; a later one is refused as complete.
	run.markFinished()
	cancelled := run.cancelRequested()
	if cancelled {
		// The workflow goroutine returning is not the termination boundary for an
		// active engine request. Shutdown cancels and joins the generation handle,
		// then kills/reaps the owned process before this result is published.
		if shutdownErr := a.images.Engine().Shutdown(); shutdownErr != nil {
			if runErr == nil {
				runErr = fmt.Errorf("%w; engine termination failed: %v", errImageWorkflowCancelled, shutdownErr)
			} else {
				runErr = fmt.Errorf("%w; engine termination failed: %v", runErr, shutdownErr)
			}
		} else if runErr == nil {
			runErr = errImageWorkflowCancelled
		}
	}
	// The parent scope ends before the terminal reclamation decision. Until
	// then, reclamation is refused because the workflow still owns the resident
	// engine and its component claims.
	a.endImageWorkflow(parent)
	// One terminal reclamation decision, using the same policy and owner as a
	// single generation. A failure is still evidence: a failed run can leave the
	// process abandoned or holding transient buffers.
	a.reclaimImageEngineAfterGenerationLocked("", lastGen, runErr)
	result.Reclaim = a.imageReclaim.report()
	if runErr != nil {
		// Return the partial evidence with the error in the result. Returning a
		// non-nil Go error makes the Wails bridge discard the stage list, so the UI
		// could say "failed" but never show which stage failed or what completed
		// before it. The structured result is the user-facing failure report.
		result.Error = runErr.Error()
		return result, nil
	}
	return result, nil
}

func (a *App) runImageWorkflowLocked(ctx context.Context, exec *imageWorkflowExecution, opts ImageWorkflowOptions) (*ImageWorkflowResult, imagetypes.ImageOptions, error) {
	opts = normalizeImageWorkflowOptions(opts)
	baseGen := imagetypes.ImageOptions{
		ModelName:          opts.ModelName,
		Prompt:             opts.Prompt,
		NegativePrompt:     opts.NegativePrompt,
		Width:              opts.BaseWidth,
		Height:             opts.BaseHeight,
		Steps:              opts.BaseSteps,
		Seed:               opts.Seed,
		Sampler:            opts.Sampler,
		Scheduler:          opts.Scheduler,
		CfgScale:           opts.CfgScale,
		Mode:               "txt2img",
		BatchCount:         1,
		IPAdapterPath:      opts.IPAdapterPath,
		ClipVisionPath:     opts.ClipVisionPath,
		IPAdapterImagePath: opts.IPAdapterImagePath,
		IPAdapterStrength:  workflowIPAdapterStrength(opts),
		ClipVisionCPU:      opts.ClipVisionCPU,
		EngineProfile:      opts.EngineProfile,
		MaxVRAM:            opts.MaxVRAM,
	}
	lastGen := baseGen
	result := &ImageWorkflowResult{Stages: []ImageWorkflowStage{}}

	if err := validateImageWorkflowOptions(opts); err != nil {
		return result, lastGen, err
	}
	modelName := opts.ModelName
	if modelName == "" {
		modelName = a.images.LastModelName()
	}
	if modelName == "" {
		return result, lastGen, fmt.Errorf("image workflow needs a model")
	}
	opts.ModelName = modelName
	baseGen.ModelName = modelName

	tempDir, err := a.workflowTempDir()
	if err != nil {
		return result, lastGen, err
	}

	batch := opts.BatchCount
	if batch < 1 {
		batch = 1
	}
	for runIndex := 0; runIndex < batch; runIndex++ {
		if !exec.acceptStage() {
			return result, lastGen, workflowStageRefused(ctx)
		}
		run, runGen, err := a.runOneWorkflowLocked(ctx, exec, opts, baseGen, tempDir, runIndex)
		lastGen = runGen
		result.BatchRuns = append(result.BatchRuns, *run)
		result.Stages = run.Stages
		result.FinalImagePath = run.FinalImagePath
		result.Width = run.Width
		result.Height = run.Height
		if err != nil {
			return result, lastGen, err
		}
		// Evaluate the same single reclamation owner at an item boundary. The
		// batch is still one workflow request, but this records whether the
		// resident process survived the completed item before the next item
		// starts. The terminal call in runImageWorkflow covers the last item.
		if runIndex < batch-1 {
			a.reclaimImageEngineAfterGenerationLocked("", runGen, nil)
			result.BatchRuns[len(result.BatchRuns)-1].Reclaim = a.imageReclaim.report()
		}
	}
	// The UI shows the final pass; a batch's earlier passes stay in BatchRuns as
	// evidence without doubling the response with every base64 image.
	if len(result.BatchRuns) > 0 {
		lastPath := result.BatchRuns[len(result.BatchRuns)-1].FinalImagePath
		if data, err := os.ReadFile(lastPath); err == nil {
			result.FinalImageB64 = base64.StdEncoding.EncodeToString(data)
		}
	}
	return result, lastGen, nil
}

func (a *App) runOneWorkflowLocked(ctx context.Context, exec *imageWorkflowExecution, opts ImageWorkflowOptions, baseGen imagetypes.ImageOptions, tempDir string, runIndex int) (*ImageWorkflowRun, imagetypes.ImageOptions, error) {
	run := &ImageWorkflowRun{Index: runIndex}
	lastGen := baseGen
	// A fixed seed advances per pass; -1 stays random for each pass.
	baseGen.Seed = workflowSeed(opts.Seed, runIndex)
	lastGen = baseGen
	if !exec.acceptStage() {
		return run, lastGen, workflowStageRefused(ctx)
	}

	// Stage 1: native txt2img. The ControlNet and upscaler directory are startup
	// options: sd-server only accepts them when the process starts, and DGL's
	// file allowlist is fixed then too. They are sent through a copy of the
	// generation request so this stage's HTTP body does not claim to do control.
	startupGen := baseGen
	startupGen.ControlNetPath = opts.ControlNetPath
	startupGen.UpscalerDir = filepath.Dir(opts.UpscalerPath)
	if opts.IPAdapterPath != "" && startupGen.EngineProfile == "" {
		startupGen.EngineProfile = imagetypes.EngineProfileIPAdapter
	}
	stage := newWorkflowStage(1, workflowStageBase, "sdcpp-native")
	stage.ResourcesBefore = a.workflowResourceSample()
	run.Stages = append(run.Stages, stage)
	a.emitWorkflowStage(runIndex, &run.Stages[0])
	started := time.Now()
	reusedBase := a.imageEngineLoadedSame(opts.ModelName)
	if err := a.ensureEngineReadyLocked(opts.ModelName, "", startupGen); err != nil {
		a.finishWorkflowStage(&run.Stages[0], started, "", 0, 0, "", err)
		return run, lastGen, fmt.Errorf("stage %d %s: %w", stage.Index, stage.Name, err)
	}
	if err := workflowContextError(ctx); err != nil {
		a.finishWorkflowStage(&run.Stages[0], started, "", 0, 0, "", err)
		return run, lastGen, fmt.Errorf("stage %d %s: %w", stage.Index, stage.Name, err)
	}
	var baseLease *resource.AdmissionLease
	if reusedBase {
		admission, lease, admitErr := a.admitWorkflowDelta(baseGen, 0, true)
		run.Stages[0].Admission = admission
		if admitErr != nil {
			a.finishWorkflowStage(&run.Stages[0], started, "", 0, 0, "", admitErr)
			return run, lastGen, fmt.Errorf("stage %d %s: %w", stage.Index, stage.Name, admitErr)
		}
		baseLease = lease
	}
	stage1Engine, ok := a.images.Engine().(*sdcpp.Engine)
	if !ok || stage1Engine == nil {
		err := fmt.Errorf("image engine does not support the base generation route")
		a.finishWorkflowStage(&run.Stages[0], started, "", 0, 0, "", err)
		return run, lastGen, err
	}
	baseResult, err := stage1Engine.GenerateTxt2ImgContext(ctx, baseGen)
	if baseLease != nil {
		baseLease.Release()
	}
	if err != nil {
		a.finishWorkflowStage(&run.Stages[0], started, "", 0, 0, "", err)
		return run, lastGen, fmt.Errorf("stage %d %s: %w", stage.Index, stage.Name, err)
	}
	stage1Path, width, height, err := a.writeWorkflowStageImage(tempDir, runIndex, 1, baseResult.ImageData, nil)
	if err != nil {
		a.finishWorkflowStage(&run.Stages[0], started, "", 0, 0, "", err)
		return run, lastGen, fmt.Errorf("stage %d %s: %w", stage.Index, stage.Name, err)
	}
	a.finishWorkflowStage(&run.Stages[0], started, stage1Path, width, height, "", nil)
	a.emitWorkflowStage(runIndex, &run.Stages[0])

	engine, ok := a.images.Engine().(*sdcpp.Engine)
	if !ok || engine == nil {
		err := fmt.Errorf("image engine does not support the selected upscale route")
		return run, lastGen, err
	}
	if !exec.acceptStage() {
		return run, lastGen, workflowStageRefused(ctx)
	}
	ipMode := opts.IPAdapterPath != ""
	upscaleExecutor := "dgl"
	if ipMode {
		upscaleExecutor = "sdcpp-native-upscale"
	}
	// Stage 2: ESRGAN upscale. DGL and the native route both add the upscaler
	// as a resident adapter, so the delta admission charges its measured bytes
	// once; stage 4 reuses it only through this explicit claim.
	stage2 := newWorkflowStage(2, workflowStageUpscale, upscaleExecutor)
	stage2.ResourcesBefore = a.workflowResourceSample()
	run.Stages = append(run.Stages, stage2)
	a.emitWorkflowStage(runIndex, &run.Stages[1])
	started = time.Now()
	if err := a.requireWorkflowEngine("stage 2 " + stage2.Name); err != nil {
		a.finishWorkflowStage(&run.Stages[1], started, "", 0, 0, "", err)
		return run, lastGen, err
	}
	upscalerKey := a.workflowUpscalerKey(engine, opts.UpscalerPath)
	upscalerClaim, claimErr := a.claimWorkflowComponent(upscalerKey)
	if claimErr != nil {
		a.finishWorkflowStage(&run.Stages[1], started, "", 0, 0, "", claimErr)
		return run, lastGen, fmt.Errorf("stage 2 %s: %w", stage2.Name, claimErr)
	}
	defer upscalerClaim.Release()
	upGen := imagetypes.ImageOptions{Width: width, Height: height, Mode: "dgl_upscale"}
	if ipMode {
		upGen.Mode = "native_upscale"
	}
	admission, lease, err := a.admitWorkflowDelta(upGen, a.fileSizeBytes(opts.UpscalerPath), true)
	run.Stages[1].Admission = admission
	if err != nil {
		a.finishWorkflowStage(&run.Stages[1], started, "", 0, 0, "", err)
		return run, lastGen, fmt.Errorf("stage 2 %s: %w", stage2.Name, err)
	}
	var upscaled []byte
	var upscaleJobID string
	if ipMode {
		upscaled, upscaleJobID, err = engine.NativeUpscale(ctx, stage1Path, opts.UpscalerPath, opts.UpscalerTile)
	} else {
		var dglReport sdcpp.DGLUpscaleReport
		upscaled, dglReport, err = engine.DGLUpscale(ctx, stage1Path, opts.UpscalerPath, opts.UpscalerTile)
		upscaleJobID = dglReport.JobID
	}
	lease.Release()
	if err != nil {
		a.finishWorkflowStage(&run.Stages[1], started, "", 0, 0, upscaleJobID, err)
		return run, lastGen, fmt.Errorf("stage 2 %s: %w", stage2.Name, err)
	}
	stage2Path, width, height, err := a.writeWorkflowStageImage(tempDir, runIndex, 2, "", upscaled)
	if err != nil {
		a.finishWorkflowStage(&run.Stages[1], started, "", 0, 0, upscaleJobID, err)
		return run, lastGen, fmt.Errorf("stage 2 %s: %w", stage2.Name, err)
	}
	a.finishWorkflowStage(&run.Stages[1], started, stage2Path, width, height, upscaleJobID, nil)
	a.emitWorkflowStage(runIndex, &run.Stages[1])

	// Stage 3: native ControlNet img2img against the upscaled base. The engine is
	// already loaded with the same startup options, so this must not reload the
	// checkpoint. If a previous stage's failure caused a restart, load it now.
	refineGen := imagetypes.ImageOptions{
		ModelName:          opts.ModelName,
		Prompt:             opts.Prompt,
		NegativePrompt:     opts.NegativePrompt,
		Width:              opts.RefineWidth,
		Height:             opts.RefineHeight,
		Steps:              opts.RefineSteps,
		Seed:               baseGen.Seed,
		Sampler:            opts.Sampler,
		Scheduler:          opts.Scheduler,
		CfgScale:           opts.CfgScale,
		Mode:               "img2img",
		InitImagePath:      stage2Path,
		ControlImagePath:   opts.PoseImagePath,
		ControlNetPath:     opts.ControlNetPath,
		ControlStrength:    opts.ControlStrength,
		Strength:           opts.RefineStrength,
		BatchCount:         1,
		IPAdapterPath:      opts.IPAdapterPath,
		ClipVisionPath:     opts.ClipVisionPath,
		IPAdapterImagePath: opts.IPAdapterImagePath,
		IPAdapterStrength:  workflowIPAdapterStrength(opts),
		ClipVisionCPU:      opts.ClipVisionCPU,
		EngineProfile:      opts.EngineProfile,
		MaxVRAM:            opts.MaxVRAM,
	}
	lastGen = refineGen
	if !exec.acceptStage() {
		return run, lastGen, workflowStageRefused(ctx)
	}
	stage3 := newWorkflowStage(3, workflowStagePose, "sdcpp-native")
	stage3.ResourcesBefore = a.workflowResourceSample()
	run.Stages = append(run.Stages, stage3)
	a.emitWorkflowStage(runIndex, &run.Stages[2])
	started = time.Now()
	if !a.imageEngineLoadedSame(opts.ModelName) {
		loadGen := refineGen
		loadGen.UpscalerDir = filepath.Dir(opts.UpscalerPath)
		if err := a.ensureEngineReadyLocked(opts.ModelName, "", loadGen); err != nil {
			a.finishWorkflowStage(&run.Stages[2], started, "", 0, 0, "", err)
			return run, lastGen, fmt.Errorf("stage 3 %s: %w", stage3.Name, err)
		}
		if err := workflowContextError(ctx); err != nil {
			a.finishWorkflowStage(&run.Stages[2], started, "", 0, 0, "", err)
			return run, lastGen, fmt.Errorf("stage 3 %s: %w", stage3.Name, err)
		}
	}
	admission, lease, err = a.admitWorkflowDelta(refineGen, 0, false)
	run.Stages[2].Admission = admission
	if err != nil {
		a.finishWorkflowStage(&run.Stages[2], started, "", 0, 0, "", err)
		return run, lastGen, fmt.Errorf("stage 3 %s: %w", stage3.Name, err)
	}
	// Native route: the compatibility route silently drops ControlNet fields
	// (control and no-control requests produced identical output), while the
	// native route loads and applies the ControlNet and honours steps nested
	// under sample_params. Both properties are required by this stage.
	refined, err := engine.GenerateImg2ImgNativeControlContext(ctx, refineGen)
	lease.Release()
	if err != nil {
		a.finishWorkflowStage(&run.Stages[2], started, "", 0, 0, "", err)
		return run, lastGen, fmt.Errorf("stage 3 %s: %w", stage3.Name, err)
	}
	stage3Path, width, height, err := a.writeWorkflowStageImage(tempDir, runIndex, 3, refined.ImageData, nil)
	if err != nil {
		a.finishWorkflowStage(&run.Stages[2], started, "", 0, 0, "", err)
		return run, lastGen, fmt.Errorf("stage 3 %s: %w", stage3.Name, err)
	}
	a.finishWorkflowStage(&run.Stages[2], started, stage3Path, width, height, "", nil)
	a.emitWorkflowStage(runIndex, &run.Stages[2])

	if !exec.acceptStage() {
		return run, lastGen, workflowStageRefused(ctx)
	}
	// Stage 4: reuse the upscaler. It is already resident from stage 2, so the
	// delta admission charges only the decode/execution working set; charging the
	// file again would double-count the same resident adapter.
	stage4 := newWorkflowStage(4, workflowStageFinal, upscaleExecutor)
	stage4.ResourcesBefore = a.workflowResourceSample()
	run.Stages = append(run.Stages, stage4)
	a.emitWorkflowStage(runIndex, &run.Stages[3])
	started = time.Now()
	if err := a.requireWorkflowEngine("stage 4 " + stage4.Name); err != nil {
		a.finishWorkflowStage(&run.Stages[3], started, "", 0, 0, "", err)
		return run, lastGen, err
	}
	if err := a.requireWorkflowComponent(upscalerKey, "stage 4 "+stage4.Name); err != nil {
		a.finishWorkflowStage(&run.Stages[3], started, "", 0, 0, "", err)
		return run, lastGen, err
	}
	finalGen := imagetypes.ImageOptions{Width: width, Height: height, Mode: "dgl_upscale"}
	if ipMode {
		finalGen.Mode = "native_upscale"
	}
	admission, lease, err = a.admitWorkflowDelta(finalGen, 0, true)
	run.Stages[3].Admission = admission
	if err != nil {
		a.finishWorkflowStage(&run.Stages[3], started, "", 0, 0, "", err)
		return run, lastGen, fmt.Errorf("stage 4 %s: %w", stage4.Name, err)
	}
	var finalBytes []byte
	var finalJobID string
	if ipMode {
		finalBytes, finalJobID, err = engine.NativeUpscale(ctx, stage3Path, opts.UpscalerPath, opts.UpscalerTile)
	} else {
		var dglReport sdcpp.DGLUpscaleReport
		finalBytes, dglReport, err = engine.DGLUpscale(ctx, stage3Path, opts.UpscalerPath, opts.UpscalerTile)
		finalJobID = dglReport.JobID
	}
	lease.Release()
	if err != nil {
		a.finishWorkflowStage(&run.Stages[3], started, "", 0, 0, finalJobID, err)
		return run, lastGen, fmt.Errorf("stage 4 %s: %w", stage4.Name, err)
	}
	finalPath, width, height, err := a.writeWorkflowStageImage(tempDir, runIndex, 4, "", finalBytes)
	if err != nil {
		a.finishWorkflowStage(&run.Stages[3], started, "", 0, 0, finalJobID, err)
		return run, lastGen, fmt.Errorf("stage 4 %s: %w", stage4.Name, err)
	}
	a.finishWorkflowStage(&run.Stages[3], started, finalPath, width, height, finalJobID, nil)
	a.emitWorkflowStage(runIndex, &run.Stages[3])

	run.FinalImagePath = finalPath
	run.Width = width
	run.Height = height
	return run, lastGen, nil
}

func normalizeImageWorkflowOptions(opts ImageWorkflowOptions) ImageWorkflowOptions {
	if opts.BaseWidth <= 0 {
		opts.BaseWidth = 512
	}
	if opts.BaseHeight <= 0 {
		opts.BaseHeight = 512
	}
	if opts.BaseSteps <= 0 {
		opts.BaseSteps = 20
	}
	if opts.RefineWidth <= 0 {
		opts.RefineWidth = 768
	}
	if opts.RefineHeight <= 0 {
		opts.RefineHeight = 768
	}
	if opts.RefineSteps <= 0 {
		opts.RefineSteps = opts.BaseSteps
	}
	if opts.RefineStrength <= 0 {
		opts.RefineStrength = 0.5
	}
	if opts.ControlStrength <= 0 {
		opts.ControlStrength = 0.9
	}
	if opts.IPAdapterPath != "" && opts.IPAdapterStrength <= 0 {
		opts.IPAdapterStrength = 0.7
	}
	if opts.IPAdapterPath == "" {
		// IP-only parameters must not leak into a plain profile request; the
		// pinned non-IP engine rejects ip_adapter_strength without a DGL script.
		opts.IPAdapterStrength = 0
	}
	if opts.IPAdapterPath != "" && opts.EngineProfile == "" {
		opts.EngineProfile = imagetypes.EngineProfileIPAdapter
	}
	if opts.UpscalerTile <= 0 {
		opts.UpscalerTile = 128
	}
	if opts.MaxVRAM != nil && *opts.MaxVRAM == 0 {
		// 0 means Auto to the panel; keep it explicit here so the load path
		// applies the post-reserve cap instead of inheriting a stale budget.
		opts.MaxVRAM = nil
	}
	if opts.CfgScale <= 0 {
		opts.CfgScale = 7
	}
	if opts.BatchCount < 1 {
		opts.BatchCount = 1
	}
	return opts
}

func validateImageWorkflowOptions(opts ImageWorkflowOptions) error {
	if strings.TrimSpace(opts.Prompt) == "" {
		return fmt.Errorf("image workflow needs a prompt")
	}
	if strings.TrimSpace(opts.UpscalerPath) == "" {
		return fmt.Errorf("image workflow needs an ESRGAN upscaler path")
	}
	if strings.TrimSpace(opts.ControlNetPath) == "" {
		return fmt.Errorf("image workflow needs an OpenPose ControlNet path")
	}
	if strings.TrimSpace(opts.PoseImagePath) == "" {
		return fmt.Errorf("image workflow needs a pose image; pose extraction is an external OpenPose/DWPose step, not a ControlNet feature")
	}
	// IP-Adapter is all-or-nothing: an adapter without its CLIP-Vision encoder
	// and reference image would run as a plain prompt and look successful.
	if opts.IPAdapterPath != "" || opts.ClipVisionPath != "" || opts.IPAdapterImagePath != "" {
		if opts.IPAdapterPath == "" || opts.ClipVisionPath == "" || opts.IPAdapterImagePath == "" {
			return fmt.Errorf("IP-Adapter needs adapter, CLIP-Vision encoder and reference image; missing %s",
				strings.Join(missingIPWorkflowInputs(opts), ", "))
		}
	}
	files := map[string]string{
		"upscaler":   opts.UpscalerPath,
		"controlnet": opts.ControlNetPath,
		"pose image": opts.PoseImagePath,
	}
	if opts.IPAdapterPath != "" {
		files["IP-Adapter"] = opts.IPAdapterPath
		files["CLIP-Vision"] = opts.ClipVisionPath
		files["IP-Adapter reference image"] = opts.IPAdapterImagePath
	}
	for name, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("image workflow %s is not readable: %w", name, err)
		}
		if info.IsDir() {
			return fmt.Errorf("image workflow %s is a directory: %s", name, path)
		}
	}
	return nil
}

func missingIPWorkflowInputs(opts ImageWorkflowOptions) []string {
	var missing []string
	if opts.IPAdapterPath == "" {
		missing = append(missing, "adapter")
	}
	if opts.ClipVisionPath == "" {
		missing = append(missing, "CLIP-Vision encoder")
	}
	if opts.IPAdapterImagePath == "" {
		missing = append(missing, "reference image")
	}
	return missing
}

func workflowIPAdapterStrength(opts ImageWorkflowOptions) float64 {
	if opts.IPAdapterPath == "" {
		return 0
	}
	return opts.IPAdapterStrength
}

func workflowSeed(seed int64, runIndex int) int64 {
	if seed < 0 {
		return -1
	}
	if seed == 0 {
		return 0
	}
	return seed + int64(runIndex)
}

func newWorkflowStage(index int, name, executor string) ImageWorkflowStage {
	return ImageWorkflowStage{
		Index:     index,
		Name:      name,
		Status:    "running",
		Executor:  executor,
		StartedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func (a *App) finishWorkflowStage(stage *ImageWorkflowStage, started time.Time, path string, width, height int, jobID string, err error) {
	stage.DurationMS = time.Since(started).Milliseconds()
	stage.CompletedAt = time.Now().UTC().Format(time.RFC3339)
	stage.OutputPath = path
	stage.OutputWidth = width
	stage.OutputHeight = height
	stage.JobID = jobID
	if err != nil {
		stage.Status = "failed"
		stage.Error = err.Error()
	} else {
		stage.Status = "completed"
	}
	stage.Placement = a.buildImagePlacementReport("workflow")
	stage.ResourcesAfter = a.workflowResourceSample()
	pid := 0
	if engine, ok := a.images.Engine().(*sdcpp.Engine); ok && engine != nil {
		pid = engine.ServerPID()
	}
	beforeUsed, beforeFree := workflowSampleValue(stage.ResourcesBefore, "vram_used_mb"), workflowSampleValue(stage.ResourcesBefore, "vram_free_mb")
	afterUsed, afterFree := workflowSampleValue(stage.ResourcesAfter, "vram_used_mb"), workflowSampleValue(stage.ResourcesAfter, "vram_free_mb")
	beforeRAM, beforeAvail := workflowSampleValue(stage.ResourcesBefore, "ram_used_mb"), workflowSampleValue(stage.ResourcesBefore, "ram_available_mb")
	afterRAM, afterAvail := workflowSampleValue(stage.ResourcesAfter, "ram_used_mb"), workflowSampleValue(stage.ResourcesAfter, "ram_available_mb")
	fmt.Printf("[ImageWorkflow] stage=%d name=%q status=%s duration=%dms pid=%d vram_before=%vMB free_before=%vMB vram_after=%vMB free_after=%vMB ram_before=%vMB avail_before=%vMB ram_after=%vMB avail_after=%vMB output=%q admission=%q error=%q\n",
		stage.Index, stage.Name, stage.Status, stage.DurationMS, pid, beforeUsed, beforeFree, afterUsed, afterFree, beforeRAM, beforeAvail, afterRAM, afterAvail, stage.OutputPath, stage.Admission, stage.Error)
}

func workflowSampleValue(sample map[string]interface{}, key string) interface{} {
	if sample == nil {
		return nil
	}
	return sample[key]
}

// workflowResourceSample is the live hardware view at a stage boundary. It is
// deliberately read-only: it records the coordinator's measurement and never
// feeds a new admission number back into the authority.
func (a *App) workflowResourceSample() map[string]interface{} {
	snap, err := resource.GetManager().RefreshSnapshot()
	if err != nil {
		snap = resource.GetManager().HardwareSnapshot()
	}
	const mb = 1024 * 1024
	out := map[string]interface{}{
		"captured_at":      time.Now().UTC().Format(time.RFC3339),
		"vram_used_mb":     snap.UsedVRAMBytes / mb,
		"vram_free_mb":     snap.FreeVRAMBytes / mb,
		"vram_total_mb":    snap.TotalVRAMBytes / mb,
		"ram_used_mb":      snap.UsedRAMBytes / mb,
		"ram_available_mb": snap.AvailableRAM / mb,
		"vram_live":        snap.VRAMLive,
		"vram_source":      snap.VRAMSource,
		"gpu":              snap.GPUModel,
	}
	if snap.DetectionError != "" {
		out["detection_error"] = snap.DetectionError
	}
	return out
}

func (a *App) emitWorkflowStage(runIndex int, stage *ImageWorkflowStage) {
	a.emitEvent("image:workflow", map[string]interface{}{
		"run":   runIndex,
		"stage": *stage,
	})
}

// admitWorkflowDelta asks the one admission authority about the additional
// working set of a stage. It deliberately does not re-charge resident weights:
// the checkpoint was admitted when the engine loaded, and a reused upscaler was
// admitted in the first upscale. The lease is held for the stage and released
// after it, so a concurrent start cannot pass against the same free value.
func (a *App) admitWorkflowDelta(gen imagetypes.ImageOptions, extraBytes int64, requireGPU bool) (string, *resource.AdmissionLease, error) {
	if a.admission == nil {
		a.admission = resource.DefaultAuthority()
	}
	loaded := a.images.LoadOptions()
	req := resource.AdmissionRequest{
		Kind:           resource.ResourceImage,
		Name:           loaded.ModelName,
		Path:           loaded.ModelPath,
		VAETiling:      a.images.LoadedOptions().VAETiling,
		ControlNetCPU:  a.images.LoadedOptions().ControlNetCPU,
		ResidentLoaded: true,
		Width:          gen.Width,
		Height:         gen.Height,
		BatchCount:     gen.BatchCount,
		Steps:          gen.Steps,
		Family:         loaded.Family,
		Mode:           gen.Mode,
		UpscalerBytes:  extraBytes,
		RequireGPU:     requireGPU && !imageCPUOnly(loaded),
		CPUOnly:        imageCPUOnly(loaded),
		DeviceIndex:    0,
	}
	decision, lease := a.admission.Admit(req)
	summary := fmt.Sprintf("admitted=%t projected=%dMB required=%dMB activation=%dMB tiling=%t free=%dMB reserve=%dMB %s",
		decision.Admitted, decision.Estimate.ProjectedPeakVRAMMB, decision.Estimate.RequiredVRAMMB,
		decision.Estimate.ActivationVRAMMB, decision.Estimate.TilingPlanned,
		decision.Estimate.FreeVRAMMB, decision.Estimate.ReserveVRAMMB, decision.Reason)
	if !decision.Admitted {
		lease.Release()
		return summary, nil, fmt.Errorf("%s", decision.Reason)
	}
	return summary, lease, nil
}

func (a *App) workflowTempDir() (string, error) {
	base := ""
	if engine, ok := a.images.Engine().(*sdcpp.Engine); ok {
		base = engine.GetConfig().TempDir
	}
	if base == "" {
		base = isolatedTempDir()
	}
	dir := filepath.Join(base, "workflow-"+time.Now().UTC().Format("20060102_150405"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create workflow temp directory: %w", err)
	}
	return dir, nil
}

func (a *App) writeWorkflowStageImage(dir string, runIndex, stage int, b64 string, raw []byte) (string, int, int, error) {
	data := raw
	if len(data) == 0 {
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(b64))
		if err != nil {
			return "", 0, 0, fmt.Errorf("decode stage image: %w", err)
		}
		data = decoded
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", 0, 0, fmt.Errorf("decode stage image dimensions: %w", err)
	}
	path := filepath.Join(dir, fmt.Sprintf("run-%d-stage-%d.png", runIndex+1, stage))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", 0, 0, fmt.Errorf("write stage image: %w", err)
	}
	return path, cfg.Width, cfg.Height, nil
}
