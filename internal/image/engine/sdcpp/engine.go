package sdcpp

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"howl-chat/internal/image/engine"
	"howl-chat/internal/image/types"
)

// Engine implements engine.ImageEngine using stable-diffusion.cpp sd-server.
type Engine struct {
	config       types.EngineConfig
	currentModel *types.ModelInfo
	serverProc   *os.Process
	httpClient   *client
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
	var target *types.ModelInfo

	if opts.ModelPath != "" {
		// Absolute path supplied — build ModelInfo directly, no ModelsDir scan.
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
	} else {
		models, err := e.ListModels()
		if err != nil {
			return err
		}
		for i := range models {
			if models[i].Name == opts.ModelName {
				target = &models[i]
				break
			}
		}
		if target == nil {
			return fmt.Errorf("model not found: %s (not in %s)", opts.ModelName, e.config.ModelsDir)
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
	if opts.DiffusionModel != "" {
		target.Files["diffusion_model"] = opts.DiffusionModel
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
	cmd.Stdout = os.Stdout
	stderrBuf := new(strings.Builder)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderrBuf)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start sd-server: %w", err)
	}
	e.serverProc = cmd.Process
	e.currentModel = target

	// Wait for server to be ready
	if err := e.httpClient.waitForReady(120 * time.Second); err != nil {
		stderr := stderrBuf.String()
		if strings.Contains(stderr, "get sd version from file failed") {
			return fmt.Errorf(
				"model '%s' (%s) cannot be loaded: this sd-server build failed to detect the architecture from the file. "+
					"GGUF single-file checkpoints for Flux/SD3 often use metadata formats that older sd.cpp builds cannot read. "+
					"Solutions: (1) Use the .safetensors version of this model, (2) Update sd-server to a newer build, "+
					"(3) Re-convert the GGUF with a sd.cpp-compatible converter: %w",
				target.Name, target.Format, err)
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
	if e.httpClient == nil {
		return nil, fmt.Errorf("engine not initialized")
	}
	if e.currentModel == nil {
		return nil, fmt.Errorf("no model loaded")
	}
	return e.httpClient.video(opts)
}

// GetProgress polls the sd-server /sdapi/v1/progress endpoint.
func (e *Engine) GetProgress() (types.ProgressInfo, error) {
	if e.httpClient == nil {
		return types.ProgressInfo{}, fmt.Errorf("engine not initialized")
	}
	return e.httpClient.progress()
}

// GetLoadedModel returns the currently loaded model.
func (e *Engine) GetLoadedModel() *types.ModelInfo {
	return e.currentModel
}

// Shutdown kills the server process.
func (e *Engine) Shutdown() error {
	if e.serverProc != nil {
		e.serverProc.Kill()
		e.serverProc.Wait()
		e.serverProc = nil
	}
	return nil
}

// buildServerArgs constructs sd-server CLI arguments for a model.
func (e *Engine) buildServerArgs(m *types.ModelInfo, opts types.ModelLoadOptions) []string {
	args := []string{
		"--listen-port", fmt.Sprintf("%d", e.config.Port),
		"--threads", fmt.Sprintf("%d", defaultThreads()),
	}

	// Model path logic:
	// --diffusion-model is used ONLY when the model is a split component
	// (i.e. a separate diffusion transformer alongside CLIP/T5 encoders).
	// All single-file formats (safetensors, ckpt, gguf merged checkpoints)
	// must use -m / --model regardless of family.
	//
	// NOTE: Some GGUF single-file checkpoints for Flux/SD3 (e.g. from third-party
	// converters) write metadata keys that older sd.cpp builds cannot read,
	// causing "get sd version from file failed". --prediction does not fix this
	// because version detection happens before prediction type is applied.
	// The only workarounds are: use .safetensors, update sd.cpp, or re-convert.
	if m.Files["diffusion_model"] != "" && m.Format != "gguf" {
		// Split component bundle — use the diffusion transformer flag
		args = append(args, "--diffusion-model", m.Files["diffusion_model"])
	} else if m.Files["model"] != "" {
		// Single-file checkpoint (safetensors, ckpt, or merged GGUF)
		args = append(args, "-m", m.Files["model"])
	}

	// COMMON ADDITIONS
	if m.Files["vae"] != "" {
		args = append(args, "--vae", m.Files["vae"])
	}
	if m.Files["clip_l"] != "" {
		args = append(args, "--clip-l", m.Files["clip_l"])
	}
	if m.Files["clip_g"] != "" {
		args = append(args, "--clip-g", m.Files["clip_g"])
	}
	if m.Files["t5xxl"] != "" {
		args = append(args, "--t5xxl", m.Files["t5xxl"])
	}
	if m.Files["clip_vision"] != "" {
		args = append(args, "--clip-vision", m.Files["clip_vision"])
	}
	if m.Files["taesd"] != "" {
		args = append(args, "--taesd", m.Files["taesd"])
	}
	if m.Files["control_net"] != "" {
		args = append(args, "--control-net", m.Files["control_net"])
	}

	// PREDICTION TYPE — must be set at startup for GGUF single-file checkpoints
	// where sd.cpp cannot auto-detect the architecture from the GGUF metadata.
	// Explicit value from opts takes priority; otherwise infer from model family.
	prediction := opts.Prediction
	if prediction == "" && m.Format == "gguf" {
		switch strings.ToLower(m.Family) {
		case "flux":
			prediction = "flux_flow"
		case "flux2":
			prediction = "flux2_flow"
		case "sd3":
			prediction = "sd3_flow"
		}
	}
	if prediction != "" {
		args = append(args, "--prediction", prediction)
	}

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
		// Windows-specific known locations
		`D:\StableDiffusion\sd-server.exe`,
		`D:\stable-diffusion.cpp\build\bin\Release\sd-server.exe`,
		`D:\stable-diffusion.cpp\build\bin\sd-server.exe`,
		`D:\stable-diffusion.cpp\sd-server.exe`,
		`D:\sd.cpp\sd-server.exe`,
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
