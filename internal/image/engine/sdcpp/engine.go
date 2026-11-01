package sdcpp

import (
	"fmt"
	"howl-chat/internal/paths"
	"howl-chat/internal/proc"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/image/engine"
	"howl-chat/internal/image/scanner"
	"howl-chat/internal/image/types"
)

// TODO: Implement HuggingFace repo link support for image models
// - Add HuggingFace repo URL field to ModelInfo struct
// - Display repo link in model library UI (like professional apps)
// - Allow users to click through to view model details, documentation, and download alternatives
// - Consider adding model card preview or metadata fetching from HuggingFace API
// - Reference: https://huggingface.co/docs/huggingface_hub/guides/huggingface_hub_py

// Engine implements engine.ImageEngine using stable-diffusion.cpp sd-server.
type Engine struct {
	config       types.EngineConfig
	currentModel *types.ModelInfo
	serverProc   *os.Process
	httpClient   *client

	// Video runs go through sd-cli rather than the HTTP server.
	runMu     sync.Mutex
	videoMu   sync.Mutex
	videoProc *os.Process
	progress  types.ProgressInfo

	// sd-server stdout progress (see progress_stdout.go).
	srvMu       sync.Mutex
	srvProgress types.ProgressInfo
	srvUpdated  time.Time
	cancelled   bool
}

var _ engine.ImageEngine = (*Engine)(nil)

// NewEngine creates an uninitialized sdcpp engine.
func NewEngine() *Engine {
	return &Engine{}
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
		config.TempDir = filepath.Join(os.TempDir(), "howl-image-temp")
	}
	if config.Port == 0 {
		config.Port = 8081 // default, different from llama-server 8080
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
	target, err := e.resolveTarget(opts)
	if err != nil {
		return err
	}

	// Kill existing server
	if e.serverProc != nil {
		e.serverProc.Kill()
		e.serverProc.Wait()
		e.serverProc = nil
	}

	args := e.buildServerArgs(target, opts)
	fmt.Printf("[ImageEngine] Starting sd-server: %s %v\n", e.config.BinaryPath, args)
	cmd := exec.Command(e.config.BinaryPath, args...)
	stdoutPipe, pipeErr := cmd.StdoutPipe()
	if pipeErr != nil {
		return fmt.Errorf("sd-server stdout: %w", pipeErr)
	}
	stderrBuf := new(strings.Builder)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderrBuf)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start sd-server: %w", err)
	}
	e.serverProc = cmd.Process
	proc.Attach(cmd)
	go e.pumpServerOutput(stdoutPipe)
	e.currentModel = target

	// Wait for server to be ready
	if err := e.httpClient.waitForReady(120 * time.Second); err != nil {
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
		if e.serverProc != nil {
			e.serverProc.Kill()
		}
		e.serverProc = nil
		e.currentModel = nil
		return fmt.Errorf("sd-server did not become ready: %w", err)
	}
	return nil
}

// IsRunning returns true if the sd-server subprocess is alive.
func (e *Engine) IsRunning() bool {
	if e.serverProc == nil {
		return false
	}
	if e.httpClient == nil {
		return false
	}
	// Quick health check via HTTP
	if err := e.httpClient.waitForReady(3 * time.Second); err != nil {
		return false
	}
	return true
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

// GenerateTxt2Img creates an image from a text prompt via sd-server HTTP API.
func (e *Engine) GenerateTxt2Img(opts types.ImageOptions) (*types.ImageResult, error) {
	if e.httpClient == nil {
		return nil, fmt.Errorf("engine not initialized")
	}
	if e.currentModel == nil {
		fmt.Printf("[ImageEngine] GenerateTxt2Img: currentModel is nil (engine=%p, httpClient=%p)\n", e, e.httpClient)
		return nil, fmt.Errorf("no model loaded")
	}
	e.resetServerProgress()
	result, err := e.httpClient.txt2img(opts)
	if err != nil {
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
	return e.httpClient.img2img(opts)
}

// GenerateInpaint fills a masked region of an image.
func (e *Engine) GenerateInpaint(opts types.ImageOptions) (*types.ImageResult, error) {
	if e.httpClient == nil {
		return nil, fmt.Errorf("engine not initialized")
	}
	if e.currentModel == nil {
		return nil, fmt.Errorf("no model loaded")
	}
	return e.httpClient.inpaint(opts)
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

// CancelJob sends a cancellation signal to the backend for the active job.
func (e *Engine) CancelJob() error {
	if e.cancelVideo() {
		return nil
	}
	if e.httpClient == nil {
		return nil
	}
	return e.httpClient.cancelJob("")
}

// Shutdown kills the server process and resets client state.
func (e *Engine) Shutdown() error {
	e.cancelVideo()
	if e.serverProc != nil {
		proc.KillTree(e.serverProc.Pid)
		_ = e.serverProc.Kill()
		_, _ = e.serverProc.Wait()
		e.serverProc = nil
	}
	if e.httpClient != nil {
		e.httpClient.resetState()
	}
	e.currentModel = nil
	return nil
}

// buildServerArgs constructs sd-server CLI arguments for a model.
func (e *Engine) buildServerArgs(m *types.ModelInfo, opts types.ModelLoadOptions) []string {
	args := []string{
		"--listen-port", fmt.Sprintf("%d", e.config.Port),
		"--threads", fmt.Sprintf("%d", defaultThreads()),
	}

	// MODEL PATH LOGIC
	switch m.Family {
	case "flux", "chroma", "anima":
		// These models use --diffusion-model (split component style)
		if m.Files["model"] != "" {
			args = append(args, "--diffusion-model", m.Files["model"])
		}
	default:
		// SD3, SDXL, SD1.x, SD2.x, HiDream use standard -m flag (single-file checkpoint)
		if m.Files["model"] != "" {
			args = append(args, "-m", m.Files["model"])
		}
	}

	// COMMON ADDITIONS (skip companion models for HiDream - it's self-contained)
	// Chroma also skips CLIP-L (only uses T5XXL + VAE)
	// SD1.5/SD2.1 don't use T5XXL
	familyLower := strings.ToLower(m.Family)
	if familyLower != "hidream" {
		if m.Files["vae"] != "" {
			args = append(args, "--vae", m.Files["vae"])
		}
		// Chroma doesn't use CLIP-L
		if familyLower != "chroma" {
			if m.Files["clip_l"] != "" {
				args = append(args, "--clip_l", m.Files["clip_l"])
			}
		}
		if m.Files["clip_g"] != "" {
			args = append(args, "--clip_g", m.Files["clip_g"])
		}
		// Only add T5XXL for Flux, Chroma, SD3, SD3.5
		if familyLower == "flux" || familyLower == "chroma" || familyLower == "sd3" || familyLower == "sd3.5" {
			if m.Files["t5xxl"] != "" {
				args = append(args, "--t5xxl", m.Files["t5xxl"])
			}
		}
	}
	if m.Files["clip_vision"] != "" {
		args = append(args, "--clip_vision", m.Files["clip_vision"])
	}
	if m.Files["taesd"] != "" {
		args = append(args, "--taesd", m.Files["taesd"])
		// TAESD is only used for step previews here: the final image must still be
		// decoded by the real VAE or quality drops.
		args = append(args, "--taesd-preview-only")
	}
	if m.Files["control_net"] != "" {
		// sd-server spells this with hyphens; "--control_net" is rejected as an
		// unknown argument and the server refuses to start at all.
		args = append(args, "--control-net", m.Files["control_net"])
	}

	// FAMILY-SPECIFIC FLAGS
	if strings.ToLower(m.Family) == "chroma" {
		args = append(args, "--chroma-disable-dit-mask")
		args = append(args, "--clip-on-cpu")
	}

	// ANIMA-SPECIFIC: uses --llm for Qwen3 text encoder (NOT clip_l/t5xxl)
	if strings.ToLower(m.Family) == "anima" {
		if m.Files["llm"] != "" {
			args = append(args, "--llm", m.Files["llm"])
		}
	}

	// HIDREAM-SPECIFIC: self-contained model, no companion models needed
	// Skip loading clip_l, t5xxl, vae for HiDream

	// VRAM budget - auto-detect free VRAM at startup, offload to CPU/RAM for 8GB GPU
	args = append(args, "--max-vram", "-1")

	// STARTUP FLAGS from ModelLoadOptions
	if opts.Backend != "" {
		args = append(args, "--backend", opts.Backend)
	}
	if opts.ParamsBackend != "" {
		args = append(args, "--params-backend", opts.ParamsBackend)
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
	if opts.FlashAttention {
		args = append(args, "--fa")
	}
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
		// Absolute file path supplied — build ModelInfo directly, no ModelsDir scan.
		base := filepath.Base(opts.ModelPath)
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(base)), ".")
		name := opts.ModelName
		if name == "" {
			name = strings.TrimSuffix(base, filepath.Ext(base))
		}
		target = &types.ModelInfo{
			Name:   name,
			Path:   opts.ModelPath,
			Family: inferFamilyFromName(base),
			Format: ext,
			Files:  map[string]string{"model": opts.ModelPath},
		}
		// Auto-fill companions since we bypassed ListModels()
		s := scanner.New(e.config.ModelsDir)
		idx, _ := s.LoadIndex() // Use cached index if available
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
