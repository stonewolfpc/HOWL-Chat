package main

import (
	"fmt"
	llamabackend "howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/loader"
	"howl-chat/internal/proc"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/paths"
)

// ---------------------------------------------------------------------------
// llama-server discovery
// ---------------------------------------------------------------------------

// findLlamaServer returns the absolute path to llama-server.exe.
// The bundled llama.cpp runtime lives in <project>/runtime/llama so its
// ggml/llama DLLs never collide with the stable-diffusion ones.
// findLlamaServerForRuntime is the one resolver used by the model load path.
// Tests replace it with the test binary, which can act as a scripted
// llama-server without changing the spawn or monitor behavior.
var findLlamaServerForRuntime = findLlamaServer

func findLlamaServer() string {
	candidates := []string{
		filepath.Join(paths.ProjectRoot(), "runtime", "llama", "llama-server.exe"),
		filepath.Join(filepath.Dir(os.Args[0]), "llama-server.exe"),
		"llama-server.exe",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			if abs, err := filepath.Abs(c); err == nil {
				return abs
			}
			return c
		}
	}
	return ""
}

// getProjectRoot returns the project root directory.
func getProjectRoot() string {
	return paths.ProjectRoot()
}

// findSDServer searches common locations for the sd-server.exe binary.
func findSDServer() string {
	candidates := []string{
		filepath.Join(getProjectRoot(), "sd-server.exe"),
		filepath.Join(filepath.Dir(os.Args[0]), "sd-server.exe"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	// Last resort: allow sd-server on PATH so users can keep it elsewhere.
	if p, err := exec.LookPath("sd-server.exe"); err == nil {
		return p
	}
	return ""
}

// ---------------------------------------------------------------------------
// Settings map helpers
// ---------------------------------------------------------------------------

// intSetting extracts an int from a settings map, returning def if absent.
func intSetting(m map[string]interface{}, key string, def int) string {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return fmt.Sprintf("%d", int(n))
		case int:
			return fmt.Sprintf("%d", n)
		}
	}
	return fmt.Sprintf("%d", def)
}

// strSetting extracts a string from a settings map, returning def if absent or "auto".
func strSetting(m map[string]interface{}, key string, def string) string {
	if v, ok := m[key].(string); ok && v != "" && v != "auto" {
		return v
	}
	return def
}

// ---------------------------------------------------------------------------
// llama-server lifecycle
// ---------------------------------------------------------------------------

// buildLlamaServerArgs assembles the llama-server command line. The projector
// is passed only when the caller resolved one by exact identity; this function
// never searches the directory.
func buildLlamaServerArgs(modelPath, mmprojPath string, settings map[string]interface{}, fitParams bool, caps llamabackend.CapabilitySet) []string {
	args := []string{
		"--host", "127.0.0.1",
		"--port", "8080",
		"-m", modelPath,
		"-c", intSetting(settings, "context_size", 4096),
		"--threads", intSetting(settings, "threads", 8),
		"-b", intSetting(settings, "batch_size", 512),
		"-ngl", intSetting(settings, "gpu_layers", 0),
	}

	// Howl always supplies an explicit context, batch and layer plan (either the
	// user's or the admission-fitted one), so llama.cpp's fit pass has nothing
	// left to size. Bypass it explicitly instead of leaving it silently inert.
	if !fitParams && caps.Supports("--fit") {
		args = append(args, "--fit", "off")
	}

	// Placement evidence is emitted at log verbosity 4. Without this the engine
	// still loads, but never says how many layers or MiB it actually placed, so
	// the app could only report its own arithmetic as if it were the result.
	// The capture is bounded; this changes diagnostics, not execution.
	if caps.Supports("--log-verbosity") {
		args = append(args, "--log-verbosity", "4")
	}

	if mmprojPath != "" {
		args = append(args, "--mmproj", mmprojPath)
	}

	// Flash Attention is tri-state. Sending "auto" explicitly is different from
	// sending nothing only in that it records the user's choice; "off" really
	// disables it, which the old boolean could not express.
	if mode := flashAttentionMode(settings); mode != "" && caps.Supports("--flash-attn") {
		args = append(args, "--flash-attn", mode)
	}

	// KV cache precision. "standard" is f16; the quantized types trade a small
	// amount of precision for a large KV footprint reduction. Unknown/legacy
	// values (including the removed "paged" and "flash-attention" options) are
	// not sent: the UI no longer offers them.
	if kvType := kvCacheType(settings); kvType != "" &&
		caps.Supports("--cache-type-k") && caps.Supports("--cache-type-v") {
		args = append(args, "--cache-type-k", kvType, "--cache-type-v", kvType)
	}

	// KV cache residency. Both directions are sent explicitly so the checkbox
	// state is visible in the command line and cannot be read as "auto".
	if offloadKQV, ok := settings["offload_kqv"].(bool); ok {
		if offloadKQV && caps.Supports("--kv-offload") {
			args = append(args, "--kv-offload")
		}
		if !offloadKQV && caps.Supports("--no-kv-offload") {
			args = append(args, "--no-kv-offload")
		}
	}

	// Model loading mode. The legacy use_mmap/use_mlock booleans were persisted
	// but never reached the engine; map them to llama.cpp's one flag.
	if mode := llamaLoadMode(settings); mode != "" && caps.Supports("--load-mode") {
		args = append(args, "--load-mode", mode)
	}

	if mainGPU := intSetting(settings, "main_gpu", 0); mainGPU != "0" && caps.Supports("--main-gpu") {
		args = append(args, "--main-gpu", mainGPU)
	}
	if vocabOnly, ok := settings["vocab_only"].(bool); ok && vocabOnly && caps.Supports("--vocab-only") {
		args = append(args, "--vocab-only")
	}

	if rs := strSetting(settings, "rope_mode", ""); rs != "" && rs != "auto" {
		args = append(args, "--rope-scaling", rs)
	}
	if v, ok := settings["rope_factor"].(float64); ok && v > 0 && v != 1.0 {
		args = append(args, "--rope-scale", fmt.Sprintf("%g", v))
	}
	if v, ok := settings["rope_base"].(float64); ok && v > 0 {
		args = append(args, "--rope-freq-base", fmt.Sprintf("%g", v))
	}
	if ts := tensorSplitSetting(settings); ts != "" {
		args = append(args, "--tensor-split", ts)
	}
	if jinja := strSetting(settings, "custom_jinja_template", ""); jinja != "" {
		args = append(args, "--chat-template", jinja)
	}

	// Thinking models. "--reasoning-budget 0" ends thinking immediately; a
	// positive budget caps it. The budget MESSAGE is not optional in practice:
	// when the budget runs out without it the model keeps thinking inside the
	// answer instead of answering (llama.cpp issue #22717), which is exactly the
	// loop this is meant to prevent. Verified against this binary's --help.
	switch strSetting(settings, "think_mode", "hidden") {
	case "off":
		args = append(args, "--reasoning", "off", "--reasoning-budget", "0")
	case "hidden", "on":
		args = append(args,
			"--reasoning", "on",
			"--reasoning-format", "deepseek",
			"--reasoning-budget", intSetting(settings, "think_budget", 2048),
			"--reasoning-budget-message",
			strSetting(settings, "think_budget_message", "I have thought enough. Time to answer."),
		)
	case "auto":
		args = append(args, "--reasoning", "auto")
	}

	return args
}

// tensorSplitSetting maps the UI's single 0-100 balance slider to llama.cpp's
// two-device proportion list. A string or slice supplied by an API caller is
// passed through unchanged. The old builder used a string-only helper, so the
// slider's numeric value was persisted and then silently ignored.
func tensorSplitSetting(settings map[string]interface{}) string {
	switch v := settings["tensor_split"].(type) {
	case string:
		s := strings.TrimSpace(v)
		if s != "" && s != "0" {
			return s
		}
	case float64:
		if v > 0 {
			return fmt.Sprintf("%g,%g", v, 100-v)
		}
	case int:
		if v > 0 {
			return fmt.Sprintf("%d,%d", v, 100-v)
		}
	case []float32:
		if len(v) == 0 {
			return ""
		}
		parts := make([]string, 0, len(v))
		for _, n := range v {
			parts = append(parts, fmt.Sprintf("%g", n))
		}
		return strings.Join(parts, ",")
	case []interface{}:
		parts := make([]string, 0, len(v))
		for _, raw := range v {
			if n, ok := raw.(float64); ok && n >= 0 {
				parts = append(parts, fmt.Sprintf("%g", n))
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, ",")
		}
	}
	return ""
}

// flashAttentionMode resolves the tri-state control. The legacy boolean is
// kept for old settings documents: true meant "send --flash-attn" and false
// meant "send nothing", which the engine treated as auto. The UI now writes
// flash_attention_mode directly, so off is genuinely off.
func flashAttentionMode(settings map[string]interface{}) string {
	if raw, ok := settings["flash_attention_mode"].(string); ok {
		switch strings.ToLower(strings.TrimSpace(raw)) {
		case "on", "off", "auto":
			return strings.ToLower(strings.TrimSpace(raw))
		}
	}
	if enabled, ok := settings["flash_attention"].(bool); ok {
		if enabled {
			return "on"
		}
		return "auto"
	}
	return ""
}

// kvCacheType maps the offered KV types to llama.cpp's --cache-type-k/v values.
// The removed "paged" and "flash-attention" values are deliberately not mapped:
// they described features the shipped binary does not have or a different flag.
func kvCacheType(settings map[string]interface{}) string {
	raw, _ := settings["kv_cache_type"].(string)
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "standard", "f16":
		return "f16"
	case "q8_0":
		return "q8_0"
	case "q4_0":
		return "q4_0"
	default:
		return ""
	}
}

// kvCacheRequiresFlashAttention reports the engine's own rule: a quantized V
// cache cannot be built with Flash Attention disabled. This is not visible in
// --help; the shipped binary states it at context creation.
func kvCacheRequiresFlashAttention(kvType string) bool {
	return kvType == "q8_0" || kvType == "q4_0"
}

// resolveKVType returns the cache type to send and, when the requested combination
// cannot be guaranteed on this engine, a user-facing reason. On and auto are
// passed through so the engine can use the quantized cache when it can enable
// Flash Attention; the monitor handles the engine's rejection with one f16
// retry. Off and unset are downgraded before spawn because Howl does not rely
// on the engine default when Flash Attention is not explicitly requested.
func resolveKVType(settings map[string]interface{}) (string, string) {
	kvType := kvCacheType(settings)
	if !kvCacheRequiresFlashAttention(kvType) {
		return kvType, ""
	}
	mode := flashAttentionMode(settings)
	if mode == "on" || mode == "auto" {
		return kvType, ""
	}
	detail := "is off"
	if mode == "" {
		detail = "is not set to On"
	}
	return "f16", fmt.Sprintf("KV cache type %s requires Flash Attention, but Flash Attention %s; loading with f16 KV. Set Flash Attention On to use %s.", kvType, detail, kvType)
}

// quantizedKVRejection is the shipped engine's exact context-creation error when
// a quantized V cache is requested without Flash Attention. A retry is allowed
// only for this text and only once per load.
const quantizedKVRejection = "quantized V cache requires flash_attn to be enabled"

// prepareQuantizedKVRetry is the single decision point for recovering from the
// engine rejection. It mutates the per-load settings to f16 and returns the
// requested cache type and the user-facing explanation. The caller owns the
// retry spawn, so this cannot loop or spawn a parallel load path.
func prepareQuantizedKVRetry(captured string, settings map[string]interface{}, retries *int) (requested string, warning string, ok bool) {
	if settings == nil || retries == nil || *retries >= 1 || !strings.Contains(captured, quantizedKVRejection) {
		return "", "", false
	}
	requested = kvCacheType(settings)
	if !kvCacheRequiresFlashAttention(requested) {
		return "", "", false
	}
	*retries++
	settings["kv_cache_type"] = "f16"
	warning = fmt.Sprintf("Requested KV cache type %s, but llama-server rejected it because Flash Attention was not enabled (%s). Retrying once with f16 KV; the model will load, but the KV cache will use more memory.", requested, quantizedKVRejection)
	return requested, warning, true
}

// llamaLoadMode maps the persisted mapping/locking booleans to --load-mode.
// Empty means neither key is present, so the engine keeps its own default.
func llamaLoadMode(settings map[string]interface{}) string {
	mmap, mmapSet := settings["use_mmap"].(bool)
	mlock, mlockSet := settings["use_mlock"].(bool)
	if !mmapSet && !mlockSet {
		return ""
	}
	switch {
	case mlock && mmap:
		return "mmap+mlock"
	case mlock:
		return "mlock"
	case mmap:
		return "mmap"
	default:
		return "none"
	}
}

// spawnLlamaServerWithModel starts llama-server with a model and settings.
//
// output receives both stdout and stderr. It is the primary destination, so a
// packaged GUI whose console handles are unavailable still captures the
// process's diagnostics: the console is only a mirrored copy and its write
// failures are ignored (see internal/backend/loader).
//
// fitParams controls the --fit flag; when false, "--fit off" is appended to
// bypass the automatic parameter-fitting phase that can crash with
// GGML_ASSERT(n_inputs < GGML_SCHED_MAX_SPLIT_INPUTS) on complex models.
func spawnLlamaServerWithModel(modelPath, mmprojPath string, settings map[string]interface{}, output io.Writer, fitParams bool, caps llamabackend.CapabilitySet) (*exec.Cmd, error) {
	serverPath := findLlamaServerForRuntime()
	if serverPath == "" {
		return nil, fmt.Errorf("llama-server.exe not found")
	}

	args := buildLlamaServerArgs(modelPath, mmprojPath, settings, fitParams, caps)
	// The exact spawned command line is evidence, not an implementation detail:
	// a setting that never reaches this slice is not applied.
	fmt.Printf("[LoadModel] argv: %q\n", args)

	if mmprojPath != "" {
		fmt.Printf("INFO: using paired mmproj file: %s\n", mmprojPath)
	} else {
		fmt.Printf("INFO: no paired mmproj file for model: %s\n", modelPath)
	}

	cmd := exec.Command(serverPath, args...)
	applyHidden(cmd)
	cmd.Dir = filepath.Dir(serverPath)

	capture := output
	if capture == nil {
		capture = io.Discard
	}
	// Verbosity 4 is needed for load placement evidence, but after the model is
	// loaded the same stream can contain prompts and sampled text. The guard
	// forwards only the load phase into the capture and console mirrors.
	cmd.Stdout = loader.NewLoadPhaseWriter(io.MultiWriter(capture, loader.SafeConsole(os.Stdout)))
	cmd.Stderr = loader.NewLoadPhaseWriter(io.MultiWriter(capture, loader.SafeConsole(os.Stderr)))

	// Start suspended, assign to the application job, then resume. No child code
	// runs before ownership is established, so a crash cannot leave an unowned
	// llama-server behind.
	if err := proc.Start(cmd); err != nil {
		return nil, fmt.Errorf("failed to start llama-server: %w", err)
	}
	fmt.Printf("INFO: llama-server started (pid %d) with model %s\n", cmd.Process.Pid, modelPath)
	return cmd, nil
}

// ValidateLlamaServer runs llama-server --version to confirm it is executable.
func ValidateLlamaServer() error {
	serverPath := findLlamaServerForRuntime()
	if serverPath == "" {
		return fmt.Errorf("llama-server.exe not found")
	}
	cmd := exec.Command(serverPath, "--version")
	applyHidden(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "0xc0000135") {
			return fmt.Errorf("llama-server is missing a required DLL (VC++ Redistributable or CUDA runtime). Install the matching VC++ Redist from microsoft.com/download/details.aspx?id=48145")
		}
		return fmt.Errorf("llama-server cannot run: %w (output: %s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// waitForPortFree polls until localhost:port accepts a connection or timeout expires.
func waitForPortFree(port string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ln, err := net.Listen("tcp", "localhost:"+port)
		if err == nil {
			ln.Close()
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

// ---------------------------------------------------------------------------
// processMonitor — tracks a running server process
// ---------------------------------------------------------------------------

type processMonitor struct {
	cmd     *exec.Cmd
	done    chan struct{}
	exitErr error
}

// waitMonitor starts a goroutine that waits for the process and signals completion.
func waitMonitor(cmd *exec.Cmd) *processMonitor {
	pm := &processMonitor{cmd: cmd, done: make(chan struct{})}
	go func() {
		pm.exitErr = cmd.Wait()
		close(pm.done)
	}()
	return pm
}

// isAlive returns true if the process has not yet exited.
func (pm *processMonitor) isAlive() bool {
	select {
	case <-pm.done:
		return false
	default:
		return true
	}
}

// wait waits for the monitor's single Wait to complete. A killed child is a
// successful stop; only a timeout is an error.
func (pm *processMonitor) wait(timeout time.Duration) error {
	if pm == nil {
		return nil
	}
	if timeout <= 0 {
		<-pm.done
		return nil
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-pm.done:
		return nil
	case <-timer.C:
		return fmt.Errorf("process did not exit within %s", timeout)
	}
}
