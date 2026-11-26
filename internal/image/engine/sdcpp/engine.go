package sdcpp

import (
	"context"
	"errors"
	"fmt"
	"howl-chat/internal/paths"
	"howl-chat/internal/proc"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/image/engine"
	"howl-chat/internal/image/scanner"
	"howl-chat/internal/image/types"
	"howl-chat/internal/process"
)

// TODO: Implement HuggingFace repo link support for image models
// - Add HuggingFace repo URL field to ModelInfo struct
// - Display repo link in model library UI (like professional apps)
// - Allow users to click through to view model details, documentation, and download alternatives
// - Consider adding model card preview or metadata fetching from HuggingFace API
// - Reference: https://huggingface.co/docs/huggingface_hub/guides/huggingface_hub_py

// childProcesses is the process authority this engine needs. It is an
// interface so a test can observe the stop route without a production test
// hook; the composition root always supplies the shared registry.
type childProcesses interface {
	Register(role process.Role, cmd *exec.Cmd, waiter func(time.Duration) error) (*process.Handle, error)
	Current(role process.Role) *process.Handle
	Lookup(identity process.Identity) *process.Handle
	Remove(handle *process.Handle)
	Stop(role process.Role, timeout time.Duration) error
}

// Engine implements engine.ImageEngine using stable-diffusion.cpp sd-server.
type Engine struct {
	config       types.EngineConfig
	currentModel *types.ModelInfo
	serverProc   *os.Process

	// processRegistry records the identity and single Wait owner for sd-server
	// and sd-cli. It may be shared with the composition root.
	processRegistry childProcesses
	registryOnce    sync.Once

	// The goroutine in supervisor.go owns Wait() for serverProc; these carry its
	// verdict to everything that needs to know whether the engine is alive.
	procMu   sync.Mutex
	procDone chan struct{}
	procExit error

	// onShutdown reports that this engine's process and execution scope have
	// ended. The app/domain owns registration invalidation; the engine never
	// mutates the resource registry itself.
	shutdownMu sync.Mutex
	onShutdown func()

	// loadRun serializes loads: two jobs asking for different models would
	// otherwise kill and replace each other's process mid-read.
	loadRun sync.Mutex
	// loadMu guards loading/aborting, which tell CancelJob whether stopping the
	// child cancels a load or interrupts a finished generation.
	loadMu    sync.Mutex
	loading   bool
	aborting  bool
	loadOwner string

	// Generation ownership (see generation.go): one request at a time per process.
	genMu     sync.Mutex
	genState  generationState
	genOwner  string
	genCancel context.CancelFunc
	genSeq    uint64
	// genHandle is the generation that currently owns the slot. The close reason
	// lives on the generation, so the next request cannot rewrite it.
	genHandle *genHandle

	httpClient *client

	// Video runs go through sd-cli rather than the HTTP server.
	runMu     sync.Mutex
	videoMu   sync.Mutex
	videoProc *os.Process
	// videoOwner is the job whose sd-cli render is running. Cancellation is
	// routed by it; the process handling itself is unchanged.
	videoOwner string
	progress   types.ProgressInfo

	// sd-server stdout progress (see progress_stdout.go).
	srvMu       sync.Mutex
	srvProgress types.ProgressInfo
	srvUpdated  time.Time
	cancelled   bool

	// placement is where sd-server put this model's weights (see placement.go).
	placement   Placement
	placementMu sync.Mutex
	// footprints are Howl's own measurements of the child's allocation (see
	// footprint.go), needed because this build does not report component totals.
	footprints  map[string]Footprint
	footprintMu sync.Mutex

	// workload is the last image request this engine ran. Placement is only
	// interpretable next to the resolution/batch that caused its decode and
	// activation footprint, so the report must carry both.
	workload   imageWorkload
	workloadMu sync.Mutex
}

type imageWorkload struct {
	Width, Height, Batch, Steps int
	RecordedAt                  time.Time
}

// RecordWorkload stores the request dimensions for the placement report.
func (e *Engine) RecordWorkload(width, height, batch, steps int) {
	e.workloadMu.Lock()
	e.workload = imageWorkload{Width: width, Height: height, Batch: batch, Steps: steps, RecordedAt: time.Now()}
	e.workloadMu.Unlock()
}

// LastWorkload returns the last request dimensions and when they were recorded.
func (e *Engine) LastWorkload() (width, height, batch, steps int, at time.Time) {
	e.workloadMu.Lock()
	defer e.workloadMu.Unlock()
	w := e.workload
	return w.Width, w.Height, w.Batch, w.Steps, w.RecordedAt
}

var _ engine.ImageEngine = (*Engine)(nil)

// NewEngine creates an uninitialized sdcpp engine with its own process registry.
func NewEngine() *Engine {
	return &Engine{processRegistry: process.NewRegistry()}
}

// NewEngineWithRegistry creates an engine that records its children in the
// supplied registry, so the composition root can own one process authority.
func NewEngineWithRegistry(registry *process.Registry) *Engine {
	if registry == nil {
		registry = process.NewRegistry()
	}
	return &Engine{processRegistry: registry}
}

// processes returns the engine's process registry, creating one for a
// hand-built Engine in tests.
func (e *Engine) processes() childProcesses {
	e.registryOnce.Do(func() {
		if e.processRegistry == nil {
			e.processRegistry = process.NewRegistry()
		}
	})
	return e.processRegistry
}

// Initialize validates the binary and directories.
func (e *Engine) Initialize(config types.EngineConfig) error {
	if config.BinaryPath == "" {
		config.BinaryPath = findSDServer()
	}
	if config.BinaryPath == "" {
		return fmt.Errorf("sd-server binary not found")
	}
	if _, err := os.Stat(config.BinaryPath); err != nil {
		return fmt.Errorf("sd-server binary not accessible: %w", err)
	}
	if config.ModelsDir == "" {
		config.ModelsDir = filepath.Join(filepath.Dir(os.Args[0]), "models")
	}
	if config.OutputDir == "" {
		config.OutputDir = filepath.Join(filepath.Dir(os.Args[0]), "images", "gallery")
	}
	if config.TempDir == "" {
		config.TempDir = filepath.Join(paths.TempDir(), "howl-image-temp")
	}
	if config.Port == 0 {
		config.Port = 8081 // default, different from llama-server 8080
	}
	if config.IPAdapterBinaryPath == "" {
		config.IPAdapterBinaryPath = FindIPAdapterServer()
	}

	for _, d := range []string{config.OutputDir, config.TempDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}

	e.config = config
	e.httpClient = newClient(fmt.Sprintf("http://localhost:%d", config.Port))
	fmt.Printf("[ImageEngine] Initialized: binary=%s models=%s port=%d\n",
		config.BinaryPath, config.ModelsDir, config.Port)
	return nil
}

// LoadModel spawns sd-server with the selected model and optional modular overrides.
func (e *Engine) LoadModel(opts types.ModelLoadOptions) error {
	if e.httpClient == nil {
		return fmt.Errorf("engine HTTP client not initialized")
	}
	if e.config.BinaryPath == "" || e.config.ModelsDir == "" {
		return fmt.Errorf("engine config incomplete")
	}

	// One load at a time. Loading now runs in a job's goroutine instead of in the
	// RPC, so two callers can genuinely overlap here.
	e.loadRun.Lock()
	defer e.loadRun.Unlock()

	e.beginLoad(opts.JobID)
	defer e.endLoad()

	target, err := e.resolveTarget(opts)
	if err != nil {
		return err
	}

	// Stop the previous server. The observer goroutine owns Wait(), so this waits
	// for its signal instead of calling Wait() on the same command twice.
	if prev := e.getServerProc(); prev != nil {
		_ = prev.Kill()
		if err := e.reapStoppedProcess(5 * time.Second); err != nil {
			return fmt.Errorf("stop previous sd-server: %w", err)
		}
		e.clearServerProc(prev)
		// The process that could have owned a generation is gone, which is the one
		// completion this engine can confirm without the backend's help.
		e.resetGeneration(closeForReload)
	}

	// Fail fast on a checkpoint sd-server cannot identify, instead of spending
	// the whole readiness timeout and then blaming the binary version.
	if err := validateCheckpoint(target); err != nil {
		return err
	}
	dropEmbeddedCompanions(target, opts)

	// A new model gets new numbers; the previous model's placement and process
	// footprint would be a lie.
	e.resetPlacement()
	e.resetFootprints()

	// The user may have cancelled during the disk work above, before a process
	// existed to kill.
	if e.loadAborted() {
		return errLoadCancelled
	}

	binary, err := e.engineBinaryFor(opts)
	if err != nil {
		return err
	}
	args := e.buildServerArgs(target, opts)
	// Ask the binary, not a table, whether it takes this argv. sd-server exits on
	// the first unknown argument, which used to surface as a load timeout and a
	// blamed binary version; naming the flag here is the refusal itself.
	if unsupported := UnsupportedFlags(binary, args); len(unsupported) > 0 {
		return fmt.Errorf("this sd-server build cannot start model '%s': %s. Nothing was started.",
			target.Name, strings.Join(unsupported, "; "))
	}
	fmt.Printf("[ImageEngine] Starting sd-server (%s): %s %v\n", effectiveEngineProfile(opts), binary, args)
	cmd := exec.Command(binary, args...)
	// sd-server is a console binary. Started without this it opened a Windows
	// Terminal window over the app for the whole generation (seen on the real
	// profile at 1129x635 over the WebView); the video engine already hides its
	// child this way.
	proc.HideWindow(cmd)
	stdoutPipe, pipeErr := cmd.StdoutPipe()
	if pipeErr != nil {
		return fmt.Errorf("sd-server stdout: %w", pipeErr)
	}
	stderrBuf := new(strings.Builder)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderrBuf)
	// Start suspended and assign job ownership before any engine code runs.
	if err := proc.Start(cmd); err != nil {
		return fmt.Errorf("failed to start sd-server: %w", err)
	}
	e.setServerProc(cmd.Process)
	// Identity is established before the reaper starts. A second sd-server
	// cannot overwrite this slot, and a reused PID cannot inherit this
	// registration. The reaper is still the only goroutine that calls Wait.
	done := make(chan struct{})
	handle, err := e.adoptServerProcess(cmd, done)
	if err != nil {
		proc.KillTree(cmd.Process.Pid)
		_ = cmd.Process.Kill()
		return fmt.Errorf("failed to establish sd-server ownership: %w", err)
	}
	e.startProcessReaperWithHandle(cmd, handle)
	go e.pumpServerOutput(stdoutPipe)
	e.currentModel = target

	// Wait for the server to answer, or for it to die trying. Polling readiness
	// alone meant a dead child was waited on for the entire load timeout: the
	// failed krea2 load burned 10m0.405s after printing the real error in seconds.
	timeout := loadTimeoutFor(target.Files)
	if err := e.waitReadyOrExit(timeout); err != nil {
		stderr := stderrBuf.String()
		if strings.Contains(stderr, "get sd version from file failed") {
			extra := ""
			if target.Family == "chroma" || target.Family == "sd3.5" {
				extra = " Chroma/SD3.5-family GGUF models may not be supported by older stable-diffusion.cpp builds — try the .safetensors version or update sd-server."
			}
			return fmt.Errorf(
				"model '%s' (%s) cannot be loaded: this sd-server build failed to detect the architecture from the file. "+
					"GGUF single-file checkpoints for Flux/SD3/SD3.5/SDXL often use metadata formats that older sd.cpp builds cannot read.%s "+
					"Please update your sd-server binary to the latest stable-diffusion.cpp release. "+
					"Alternative: use the .safetensors version of this model: %w",
				target.Name, target.Format, extra, err)
		}
		if cur := e.getServerProc(); cur != nil {
			_ = cur.Kill()
		}
		e.setServerProc(nil)
		e.currentModel = nil
		// A cancelled load is not a broken one: say so before the error text tries
		// to explain an exit that only happened because we asked for it.
		if e.loadAborted() {
			return errLoadCancelled
		}
		// A child that exited already said why. Reporting that as "did not become
		// ready in 10m0s" sent the reader looking for a slow load instead of at
		// the error in the stderr tail printed beside it.
		if errors.Is(err, errProcessExited) {
			return fmt.Errorf("%s failed to load: %w%s", target.Name, err, stderrTail(stderrBuf, 900))
		}
		return fmt.Errorf("sd-server did not become ready in %s: %w%s",
			timeout.Round(time.Second), err, stderrTail(stderrBuf, 900))
	}
	// Capture the resident allocation now that loading is finished, before any
	// sampling buffers exist. This is the actual after-load evidence.
	e.recordFootprint("load")
	return nil
}

// ServerPID reports the current sd-server child pid, or 0 when no process is
// owned. Reclamation uses it to tie a measurement baseline to one process.
func (e *Engine) ServerPID() int {
	if p := e.getServerProc(); p != nil {
		return p.Pid
	}
	return 0
}

// IsRunning reports whether the sd-server child still exists. This is a process
// question, not a readiness question: it used to run a three-second HTTP probe,
// which the job poller called twice a second and which answered "not running"
// for a server that was merely busy loading weights.
func (e *Engine) IsRunning() bool {
	if e.httpClient == nil {
		return false
	}
	if e.getServerProc() == nil {
		return false
	}
	done := e.processDone()
	if done == nil {
		return true // started, but not yet watched
	}
	select {
	case <-done:
		return false
	default:
		return true
	}
}

// EnsureReady starts the engine and loads the model if needed.
func (e *Engine) EnsureReady(modelName string, opts types.ModelLoadOptions) error {
	if e.httpClient == nil {
		return fmt.Errorf("engine not initialized")
	}
	// If server is not running, start it with the target model
	if !e.IsRunning() {
		if modelName == "" {
			return fmt.Errorf("no model selected")
		}
		opts.ModelName = modelName
		return e.LoadModel(opts)
	}
	// Server is running; check if the right model is loaded
	if e.currentModel == nil || e.currentModel.Name != modelName {
		if modelName == "" {
			return fmt.Errorf("no model selected")
		}
		opts.ModelName = modelName
		return e.LoadModel(opts)
	}
	return nil
}

// GetConfig returns the current engine configuration.
func (e *Engine) GetConfig() types.EngineConfig {
	return e.config
}

func effectiveEngineProfile(opts types.ModelLoadOptions) string {
	if opts.EngineProfile == types.EngineProfileIPAdapter {
		return types.EngineProfileIPAdapter
	}
	return types.EngineProfileDefault
}

// engineBinaryFor resolves the capability-routed profile. Only one engine runs
// at a time; the IP profile replaces the default process through the same
// LoadModel path, admission and reclamation owner.
func (e *Engine) engineBinaryFor(opts types.ModelLoadOptions) (string, error) {
	if opts.EngineProfile != types.EngineProfileIPAdapter {
		return e.config.BinaryPath, nil
	}
	binary := strings.TrimSpace(e.config.IPAdapterBinaryPath)
	if binary == "" {
		binary = FindIPAdapterServer()
	}
	if binary == "" {
		return "", fmt.Errorf("native IP-Adapter needs the c6dfa9d engine profile, but no IP-capable sd-server was found. Set image engine config ip_adapter_binary_path or place it at runtime/stablediffusion-ip/sd-server.exe.")
	}
	if _, err := os.Stat(binary); err != nil {
		return "", fmt.Errorf("native IP-Adapter engine is not accessible at %s: %w", binary, err)
	}
	return binary, nil
}

// GenerateTxt2Img creates an image from a text prompt via sd-server HTTP API.
func (e *Engine) GenerateTxt2Img(opts types.ImageOptions) (*types.ImageResult, error) {
	return e.GenerateTxt2ImgContext(context.Background(), opts)
}

// GenerateTxt2ImgContext is the workflow path: its generation scope is a child
// of ctx, so cancellation reaches the active request instead of being ignored.
func (e *Engine) GenerateTxt2ImgContext(ctx context.Context, opts types.ImageOptions) (*types.ImageResult, error) {
	if e.httpClient == nil {
		return nil, fmt.Errorf("engine not initialized")
	}
	if e.currentModel == nil {
		fmt.Printf("[ImageEngine] GenerateTxt2Img: currentModel is nil (engine=%p, httpClient=%p)\n", e, e.httpClient)
		return nil, fmt.Errorf("no model loaded")
	}
	e.RecordWorkload(opts.Width, opts.Height, opts.BatchCount, opts.Steps)
	h, err := e.beginGenerationWithContext(opts.JobID, ctx)
	if err != nil {
		return nil, err
	}
	e.resetServerProgress()
	var result *types.ImageResult
	if opts.IPAdapterImagePath != "" || opts.IPAdapterPath != "" {
		result, err = e.httpClient.txt2imgNative(h.ctx, opts)
	} else {
		result, err = e.httpClient.txt2img(h.ctx, opts)
	}
	if err := e.settleGeneration(h, err); err != nil {
		return nil, err
	}
	result.Metadata.Model = e.currentModel.Name
	result.Metadata.Family = e.currentModel.Family
	return result, nil
}

// GenerateImg2Img transforms an existing image with a prompt.
func (e *Engine) GenerateImg2Img(opts types.ImageOptions) (*types.ImageResult, error) {
	if e.httpClient == nil {
		return nil, fmt.Errorf("engine not initialized")
	}
	if e.currentModel == nil {
		return nil, fmt.Errorf("no model loaded")
	}
	e.RecordWorkload(opts.Width, opts.Height, opts.BatchCount, opts.Steps)
	h, err := e.beginGeneration(opts.JobID)
	if err != nil {
		return nil, err
	}
	result, err := e.httpClient.img2img(h.ctx, opts)
	if err := e.settleGeneration(h, err); err != nil {
		return nil, err
	}
	return result, nil
}

// GenerateImg2ImgSync is the workflow path for img2img: it uses the
// synchronous compatibility route because the async route ignores the requested
// steps in this build. It has no live preview, but a sequential stage reports
// its result when it completes anyway.
func (e *Engine) GenerateImg2ImgSync(opts types.ImageOptions) (*types.ImageResult, error) {
	if e.httpClient == nil {
		return nil, fmt.Errorf("engine not initialized")
	}
	if e.currentModel == nil {
		return nil, fmt.Errorf("no model loaded")
	}
	e.RecordWorkload(opts.Width, opts.Height, opts.BatchCount, opts.Steps)
	h, err := e.beginGeneration(opts.JobID)
	if err != nil {
		return nil, err
	}
	result, err := e.httpClient.img2imgSync(h.ctx, opts)
	if err := e.settleGeneration(h, err); err != nil {
		return nil, err
	}
	result.Metadata.Model = e.currentModel.Name
	result.Metadata.Family = e.currentModel.Family
	return result, nil
}

// GenerateImg2ImgNativeControl is the workflow path for a ControlNet-conditioned
// img2img stage. The compatibility route accepts the same JSON fields but does
// not load or apply the ControlNet (verified: control and no-control requests
// produced the identical output hash). The native async route does, and honours
// steps when they are sent under sample_params.
func (e *Engine) GenerateImg2ImgNativeControl(opts types.ImageOptions) (*types.ImageResult, error) {
	return e.GenerateImg2ImgNativeControlContext(context.Background(), opts)
}

// GenerateImg2ImgNativeControlContext binds the ControlNet stage to the
// workflow's cancellation scope.
func (e *Engine) GenerateImg2ImgNativeControlContext(ctx context.Context, opts types.ImageOptions) (*types.ImageResult, error) {
	if e.httpClient == nil {
		return nil, fmt.Errorf("engine not initialized")
	}
	if e.currentModel == nil {
		return nil, fmt.Errorf("no model loaded")
	}
	e.RecordWorkload(opts.Width, opts.Height, opts.BatchCount, opts.Steps)
	h, err := e.beginGenerationWithContext(opts.JobID, ctx)
	if err != nil {
		return nil, err
	}
	result, err := e.httpClient.img2imgNativeControl(h.ctx, opts)
	if err := e.settleGeneration(h, err); err != nil {
		return nil, err
	}
	result.Metadata.Model = e.currentModel.Name
	result.Metadata.Family = e.currentModel.Family
	return result, nil
}

// NativeUpscale runs the IP-profile engine's native ESRGAN upscale. The
// workflow calls this only when the selected profile advertises the route; the
// pinned engine has no such endpoint and remains on the DGL path.
func (e *Engine) NativeUpscale(ctx context.Context, imagePath, upscalerPath string, tile int) ([]byte, string, error) {
	if e.httpClient == nil {
		return nil, "", fmt.Errorf("engine not initialized")
	}
	if e.currentModel == nil {
		return nil, "", fmt.Errorf("no model loaded")
	}
	e.RecordWorkload(0, 0, 0, 0)
	h, err := e.beginGenerationWithContext("", ctx)
	if err != nil {
		return nil, "", err
	}
	data, jobID, err := e.httpClient.nativeUpscale(h.ctx, imagePath, upscalerPath, tile)
	if err := e.settleGeneration(h, err); err != nil {
		return nil, jobID, err
	}
	return data, jobID, nil
}

// GenerateInpaint fills a masked region of an image.
func (e *Engine) GenerateInpaint(opts types.ImageOptions) (*types.ImageResult, error) {
	if e.httpClient == nil {
		return nil, fmt.Errorf("engine not initialized")
	}
	if e.currentModel == nil {
		return nil, fmt.Errorf("no model loaded")
	}
	e.RecordWorkload(opts.Width, opts.Height, opts.BatchCount, opts.Steps)
	h, err := e.beginGeneration(opts.JobID)
	if err != nil {
		return nil, err
	}
	result, err := e.httpClient.inpaint(h.ctx, opts)
	if err := e.settleGeneration(h, err); err != nil {
		return nil, err
	}
	return result, nil
}

// GenerateVideo creates a series of frames (SVD/WAN).
func (e *Engine) GenerateVideo(opts types.ImageOptions) (*types.ImageResult, error) {
	return e.generateVideo(opts)
}

// GetProgress polls the sd-server /sdapi/v1/progress endpoint.
func (e *Engine) GetProgress() (types.ProgressInfo, error) {
	if p, ok := e.videoProgressInfo(); ok {
		return p, nil
	}
	if e.httpClient == nil {
		return types.ProgressInfo{}, fmt.Errorf("engine not initialized")
	}
	// Two independent sources, each better at something different:
	//   - stdout carries the real sampling step (the endpoint's step is not
	//     reliable here)
	//   - /sdapi/v1/progress carries the decoded TAESD frame in current_image
	// Merging them is what makes a live preview possible; taking either one
	// alone throws the other away.
	httpProgress, httpErr := e.httpClient.progress()

	if stdoutProgress, ok := e.recentServerProgress(); ok {
		if httpErr == nil && httpProgress.PreviewImage != "" {
			stdoutProgress.PreviewImage = httpProgress.PreviewImage
		}
		return stdoutProgress, nil
	}

	if httpErr == nil && (httpProgress.Step > 0 || httpProgress.Progress > 0 || httpProgress.PreviewImage != "") {
		return httpProgress, nil
	}
	return httpProgress, httpErr
}

// GetLoadedModel returns the currently loaded model.
func (e *Engine) GetLoadedModel() *types.ModelInfo {
	return e.currentModel
}

// CancelJob cancels the work belonging to jobID.
//
// Routing follows ownership, not "what happens to be running": a job that owns a
// generation gets its request closed even while another job is loading a model,
// and a cancel for a job that owns nothing changes nothing. A load is cancelled
// by stopping the child that is loading it; a generation by closing its request.
// The backend is never claimed to have stopped: backend interruption is
// unestablished for this build.
func (e *Engine) CancelJob(jobID string) error {
	// Identity is required. The legacy direct callers (GenerateImage,
	// SetActiveImageModel) have no job, and an empty id must not mean "whichever
	// job is active".
	if jobID == "" {
		return nil
	}
	if e.cancelVideoOwned(jobID) {
		return nil
	}
	if e.abortLoad(jobID) {
		return nil
	}
	e.cancelGeneration(jobID)
	return nil
}

// SetShutdownHook records the one observer that invalidates registration after
// this engine has ended. It is set before the engine is used.
func (e *Engine) SetShutdownHook(fn func()) {
	e.shutdownMu.Lock()
	e.onShutdown = fn
	e.shutdownMu.Unlock()
}

func (e *Engine) shutdownHook() func() {
	e.shutdownMu.Lock()
	defer e.shutdownMu.Unlock()
	return e.onShutdown
}

// Shutdown cancels and joins the active generation scope, kills the server
// process, and reports that the engine ended. It does not mutate resource
// registrations; that is the app/domain's transition in response to this exit.
func (e *Engine) Shutdown() error {
	e.cancelVideo()

	// Capture the active scope before resetting the slot, cancel it, then join
	// the request goroutine. The transport wrapper makes cancellation sufficient
	// to return even when the transport itself ignores the context.
	e.genMu.Lock()
	var generationDone <-chan struct{}
	if e.genHandle != nil {
		generationDone = e.genHandle.done
	}
	e.genMu.Unlock()
	e.resetGeneration(closeForShutdown)
	if generationDone != nil {
		<-generationDone
	}

	if cur := e.getServerProc(); cur != nil {
		proc.KillTree(cur.Pid)
		_ = cur.Kill()
		if err := e.reapStoppedProcess(5 * time.Second); err != nil {
			return fmt.Errorf("sd-server pid %d did not exit: %w", cur.Pid, err)
		}
		e.clearServerProc(cur)
	}
	e.resetPlacement()
	e.resetFootprints()
	if e.httpClient != nil {
		e.httpClient.resetState()
	}
	e.currentModel = nil
	if hook := e.shutdownHook(); hook != nil {
		hook()
	}
	return nil
}

// buildServerArgs constructs sd-server CLI arguments for a model.
func (e *Engine) buildServerArgs(m *types.ModelInfo, opts types.ModelLoadOptions) []string {
	args := []string{
		"--listen-port", fmt.Sprintf("%d", e.config.Port),
		"--threads", fmt.Sprintf("%d", defaultThreads()),
	}

	// Loading contract for this architecture: which flag carries the weights,
	// which companion roles it needs, and any flags it always requires. One
	// table entry, shared with the discovery pass and the validator.
	cap, known := Capability(m.Family)
	if !known {
		cap = DefaultCapability
	}
	if m.Files["model"] != "" {
		if cap.SplitComponents {
			args = append(args, "--diffusion-model", m.Files["model"])
		} else {
			args = append(args, "-m", m.Files["model"])
		}
	}
	emittedRole := map[string]bool{}
	for _, role := range cap.Roles() {
		p := m.Files[role]
		flag := companionFlag[role]
		if p == "" || flag == "" {
			continue
		}
		args = append(args, flag, p)
		emittedRole[role] = true
	}
	// A companion the caller supplied is never silently dropped by a family this
	// build knows nothing about: it has no role list, so without this a fully
	// equipped model reached sd-server without its text encoder. It is scoped to
	// unknown families on purpose — a declared family states which encoders it
	// reads, and passing the rest put --t5xxl and --clip_l on Anima and Z-Image,
	// which never open them (and would have handed them a half-downloaded file).
	if !known {
		for _, role := range []string{"llm", "t5xxl", "clip_l", "clip_g", "vae"} {
			if emittedRole[role] {
				continue
			}
			if p, flag := m.Files[role], companionFlag[role]; p != "" && flag != "" {
				args = append(args, flag, p)
			}
		}
	}
	args = append(args, cap.ExtraFlags...)
	clipVisionPath := m.Files["clip_vision"]
	if opts.ClipVisionPath != "" {
		clipVisionPath = opts.ClipVisionPath
	}
	if clipVisionPath != "" {
		args = append(args, "--clip_vision", clipVisionPath)
	}
	if opts.IPAdapterPath != "" {
		args = append(args, "--ip-adapter", opts.IPAdapterPath)
	}
	if m.Files["taesd"] != "" {
		args = append(args, "--taesd", m.Files["taesd"])
		// TAESD is only used for step previews here: the final image must still be
		// decoded by the real VAE or quality drops.
		args = append(args, "--taesd-preview-only")
	}
	controlNetPath := m.Files["control_net"]
	if opts.ControlNetPath != "" {
		controlNetPath = opts.ControlNetPath
	}
	if controlNetPath != "" {
		// sd-server spells this with hyphens; "--control_net" is rejected as an
		// unknown argument and the server refuses to start at all.
		args = append(args, "--control-net", controlNetPath)
	}
	if opts.UpscalerDir != "" {
		// The installed sd-server authorizes DGL file access through its model
		// roots. --hires-upscalers-dir alone did not make the upscaler readable
		// from a DGL script (verified 2026-11-24: "file model read denied"), while
		// adding the directory as a model root did. Emit both: the dedicated flag
		// is what hires-fix uses, and the model root is what this build's DGL
		// allowlist consults.
		args = append(args, "--model-dir", opts.UpscalerDir)
		args = append(args, "--hires-upscalers-dir", opts.UpscalerDir)
	}

	// Family flags (Chroma's DiT-mask switches) come from the table above.

	// Placement is left to sd-server's auto-fit, which puts weights on the
	// compute GPU, then RAM, then disk in that order and only takes what fits.
	// In this build --max-vram is a per-device budget in GiB and a negative value
	// reserves that much free VRAM, so -1 keeps 1 GiB back for sampling buffers.
	// An explicit --params-backend would disable auto-fit, which is why it is
	// only sent when the user asks for it.
	//
	// Ownership: admission decides POLICY (may this workload start, given the
	// budget we are about to hand sd-server); sd-server's backend_fit decides
	// PLACEMENT (which module's params live in VRAM). The budget below is the
	// number admission was told to assume, so the two cannot drift.
	maxVRAM := -1
	if opts.MaxVRAM != nil {
		maxVRAM = *opts.MaxVRAM
	}
	args = append(args, "--max-vram", strconv.Itoa(maxVRAM))

	// STARTUP FLAGS from ModelLoadOptions
	backend := opts.Backend
	if opts.ControlNetCPU {
		if backend == "" {
			backend = "controlnet=cpu"
		} else if !strings.Contains(backend, "controlnet=") {
			backend += ",controlnet=cpu"
		}
	}
	if opts.ClipVisionCPU {
		if backend == "" {
			backend = "clip_vision=cpu"
		} else if !strings.Contains(backend, "clip_vision=") {
			backend += ",clip_vision=cpu"
		}
	}
	if backend != "" {
		args = append(args, "--backend", backend)
	}
	if opts.ParamsBackend != "" {
		args = append(args, "--params-backend", opts.ParamsBackend)
	}
	// On-load quantization: sd-server converts weights to this type while
	// loading, so a full f16 checkpoint can be resident when it otherwise is not.
	if opts.WeightType != "" {
		args = append(args, "--type", opts.WeightType)
	}
	if opts.ConditioningCacheSize > 0 {
		args = append(args, "--conditioning-cache-size", strconv.Itoa(opts.ConditioningCacheSize))
	}
	if opts.MMap {
		args = append(args, "--mmap")
	}
	// Idle-release timer (opt-in). sd-server unloads model components from RAM
	// and VRAM after this many idle seconds. -1 disables the timer and is the
	// engine's own default, so absence never changes behaviour: a value is only
	// sent when a caller asks for it. Staging TTL is not forwarded: the engine
	// already releases staged copies when a component finishes.
	//
	// The c6dfa9d IP profile predates this flag. Omitting it leaves the engine's
	// own default (disabled), which is the same -1 the default profile uses; the
	// IP profile's process is still reclaimed by Howl's single owner.
	if opts.WeightResidencyTTL != nil && opts.EngineProfile != types.EngineProfileIPAdapter {
		args = append(args, "--weight-residency-ttl", strconv.Itoa(*opts.WeightResidencyTTL))
	}
	if opts.VAETiling {
		args = append(args, "--vae-tiling")
	}
	if opts.TemporalTiling {
		args = append(args, "--temporal-tiling")
	}
	if opts.Circular {
		args = append(args, "--circular")
	}
	// Always send the mode, including auto. The old boolean could not express
	// "off": unchecked meant no flag, which the engine treats as auto. A mode is
	// the only truthful representation of the control.
	args = append(args, "--flash-attn", flashAttentionMode(opts.FlashAttentionMode, opts.FlashAttention))
	if opts.OffloadToCPU {
		args = append(args, "--offload-to-cpu")
	}
	if opts.CacheMode != "" {
		args = append(args, "--cache-mode", opts.CacheMode)
		if opts.CacheOption != "" {
			args = append(args, "--cache-option", opts.CacheOption)
		}
	}
	if opts.DisableMetadata {
		args = append(args, "--disable-image-metadata")
	}

	return args
}

// flashAttentionMode resolves the tri-state control. Legacy callers that only
// set the boolean keep their old intent: true meant on, false meant the engine's
// auto default.
func flashAttentionMode(mode string, legacy bool) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "auto":
		return "auto"
	case "on":
		return "on"
	case "off":
		return "off"
	}
	if legacy {
		return "on"
	}
	return "auto"
}

func defaultThreads() int {
	if n := runtime.NumCPU(); n > 0 {
		return n
	}
	return 4
}

// FindSDServer locates the sd-server binary, checking bundled runtime first.
func FindSDServer() string {
	candidates := []string{
		// Bundled runtime path (preferred)
		filepath.Join(filepath.Dir(os.Args[0]), "runtime", "stablediffusion", "sd-server.exe"),
		filepath.Join(paths.ProjectRoot(), "sd-server.exe"),
		// PATH / relative to executable
		filepath.Join(filepath.Dir(os.Args[0]), "sd-server.exe"),
		"sd-server.exe",
		filepath.Join(filepath.Dir(os.Args[0]), "sd.exe"),
		"sd.exe",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

func findSDServer() string { return FindSDServer() }

// FindIPAdapterServer locates the native-IP engine profile. It is a compatibility
// bridge: one alternate binary, replaced through the normal LoadModel path, not
// a second resident engine and not a growing set of per-feature binaries.
func FindIPAdapterServer() string {
	if p := strings.TrimSpace(os.Getenv("HOWL_IP_ENGINE")); p != "" {
		return p
	}
	exeDir := filepath.Dir(os.Args[0])
	candidates := []string{
		filepath.Join(exeDir, "runtime", "stablediffusion-ip", "sd-server.exe"),
		filepath.Join(paths.ProjectRoot(), "runtime", "stablediffusion-ip", "sd-server.exe"),
		filepath.Join(paths.ProjectRoot(), "build", "engine-baseline", "backup-c6dfa9d", "sd-server.exe"),
		filepath.Join(paths.ProjectRoot(), "build", "engine-baseline", "backup-c6dfa9d", "sd-server"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

// resolveTarget finds the model described by opts and applies manual
// companion overrides. It starts no process, so a video run can resolve its
// model without paying for an sd-server load first.
func (e *Engine) resolveTarget(opts types.ModelLoadOptions) (*types.ModelInfo, error) {
	var target *types.ModelInfo

	if opts.ModelPath != "" {
		// A directory is a bundle, not a checkpoint. sd-server must be handed the
		// checkpoint inside it: given the folder itself it assumes a diffusers
		// layout and dies with "init model loader from file failed", which reached
		// the user as a generation that finished "loading" and never produced an
		// image. Resolve bundles through the same scan the model dropdown used, so
		// family and companion files come from one source of truth.
		if fi, statErr := os.Stat(opts.ModelPath); statErr == nil && fi.IsDir() {
			if models, listErr := e.ListModels(); listErr == nil {
				for i := range models {
					if models[i].Path == opts.ModelPath {
						target = &models[i]
						break
					}
				}
			}
			if target == nil {
				if m := scanner.LoadDirManifest(opts.ModelPath); m != nil && m.Files["model"] != "" {
					target = &types.ModelInfo{
						Name:   m.Name,
						Path:   opts.ModelPath,
						Family: string(m.Family),
						Format: m.Format,
						Files:  m.Files,
					}
				}
			}
			if target == nil {
				return nil, fmt.Errorf("model bundle %s has no usable checkpoint inside it", opts.ModelPath)
			}
			if fi, statErr := os.Stat(target.Files["model"]); statErr != nil || fi.IsDir() {
				return nil, fmt.Errorf("model bundle %s points at %q, which is not a model file", opts.ModelPath, target.Files["model"])
			}
		}
	}

	if target == nil && opts.ModelPath != "" {
		base := filepath.Base(opts.ModelPath)
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(base)), ".")
		name := opts.ModelName
		if name == "" {
			name = strings.TrimSuffix(base, filepath.Ext(base))
		}
		s := scanner.New(e.config.ModelsDir)
		idx, _ := s.LoadIndex() // Use cached index if available
		// Prefer the scanner's record for this exact path: it carries the
		// architecture parsed from the file and the bundle's companion files.
		// A filename guess carries neither, so a registered Z-Image checkpoint
		// reached sd-server as a bare -m with no --diffusion-model, --llm or
		// --vae, and could not load for a reason nothing in the UI explained.
		for i := range idx.Models {
			m := idx.Models[i]
			if !samePath(m.Path, opts.ModelPath) {
				continue
			}
			files := m.Files
			if files == nil || files["model"] == "" {
				files = map[string]string{"model": m.Path}
			}
			target = &types.ModelInfo{
				Name:   name,
				Path:   m.Path,
				Family: string(m.Family),
				Format: m.Format,
				Files:  files,
			}
			break
		}
		if target == nil {
			target = &types.ModelInfo{
				Name:   name,
				Path:   opts.ModelPath,
				Family: inferFamilyFromName(base),
				Format: ext,
				Files:  map[string]string{"model": opts.ModelPath},
			}
		}
		autoFillCompanions(target, idx)
	} else {
		models, err := e.ListModels()
		if err != nil {
			return nil, err
		}
		for i := range models {
			if models[i].Name == opts.ModelName {
				target = &models[i]
				break
			}
		}
		if target == nil {
			return nil, fmt.Errorf("model not found: %s (not in %s)", opts.ModelName, e.config.ModelsDir)
		}
	}

	// The app resolved this family from the image registry, which covers every
	// scanned directory; the engine only sees its own ModelsDir. Prefer that over
	// the filename guess made above, because the family decides whether the file
	// is passed as --diffusion-model (split components) or as -m (full
	// checkpoint). Getting that wrong fails architecture detection before a
	// single byte of weights is read.
	if target != nil && opts.Family != "" {
		target.Family = opts.Family
	}

	// Apply manual overrides if provided
	if opts.VAEPath != "" {
		target.Files["vae"] = opts.VAEPath
	}
	if opts.ClipLPath != "" {
		target.Files["clip_l"] = opts.ClipLPath
	}
	if opts.ClipGPath != "" {
		target.Files["clip_g"] = opts.ClipGPath
	}
	if opts.T5XXLPath != "" {
		target.Files["t5xxl"] = opts.T5XXLPath
	}
	if opts.LLMPath != "" {
		target.Files["llm"] = opts.LLMPath
	}
	if opts.DiffusionModel != "" {
		target.Files["diffusion_model"] = opts.DiffusionModel
	}

	return target, nil
}

// samePath compares two file paths the way Windows does: case-insensitively and
// after cleaning.
func samePath(a, b string) bool {
	return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
}
