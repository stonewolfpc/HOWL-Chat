package main

import (
	"fmt"
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

// findMMProj searches for an mmproj file next to the model file.
// It tries specific naming patterns first, then falls back to a directory scan.
func findMMProj(modelPath string) string {
	dir := filepath.Dir(modelPath)
	baseName := filepath.Base(modelPath)
	modelName := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	patterns := []string{
		// Quant-matched names (most common: mmproj uses same quant as model)
		filepath.Join(dir, "mmproj-IQ4_XS.gguf"),
		filepath.Join(dir, "mmproj-IQ3_XXS.gguf"),
		filepath.Join(dir, "mmproj-IQ2_XS.gguf"),
		filepath.Join(dir, "mmproj-IQ2_XXS.gguf"),
		filepath.Join(dir, "mmproj-IQ2_S.gguf"),
		filepath.Join(dir, "mmproj-Q2_K.gguf"),
		filepath.Join(dir, "mmproj-Q3_K.gguf"),
		filepath.Join(dir, "mmproj-Q3_K_S.gguf"),
		filepath.Join(dir, "mmproj-Q3_K_M.gguf"),
		filepath.Join(dir, "mmproj-Q3_K_L.gguf"),
		filepath.Join(dir, "mmproj-Q4_0.gguf"),
		filepath.Join(dir, "mmproj-Q4_K.gguf"),
		filepath.Join(dir, "mmproj-Q4_K_S.gguf"),
		filepath.Join(dir, "mmproj-Q4_K_M.gguf"),
		filepath.Join(dir, "mmproj-Q4_K_L.gguf"),
		filepath.Join(dir, "mmproj-Q5_0.gguf"),
		filepath.Join(dir, "mmproj-Q5_K.gguf"),
		filepath.Join(dir, "mmproj-Q5_K_S.gguf"),
		filepath.Join(dir, "mmproj-Q5_K_M.gguf"),
		filepath.Join(dir, "mmproj-Q6_K.gguf"),
		filepath.Join(dir, "mmproj-Q8_0.gguf"),
		filepath.Join(dir, "mmproj-F16.gguf"),
		filepath.Join(dir, "mmproj-F32.gguf"),
		filepath.Join(dir, "mmproj-BF16.gguf"),
		// Lowercase variants
		filepath.Join(dir, "mmproj-iq4_xs.gguf"),
		filepath.Join(dir, "mmproj-iq3_xxs.gguf"),
		filepath.Join(dir, "mmproj-q4_k_m.gguf"),
		filepath.Join(dir, "mmproj-q5_k_m.gguf"),
		filepath.Join(dir, "mmproj-q8_0.gguf"),
		filepath.Join(dir, "mmproj-f16.gguf"),
		filepath.Join(dir, "mmproj-bf16.gguf"),
		// Generic
		filepath.Join(dir, "mmproj.gguf"),
		filepath.Join(dir, "mmproj-model.gguf"),
		filepath.Join(dir, "mmproj-vision.gguf"),
		filepath.Join(dir, "mmproj-image.gguf"),
		filepath.Join(dir, "mmproj-proj.gguf"),
		// Model-specific
		filepath.Join(dir, modelName+"-mmproj.gguf"),
		filepath.Join(dir, modelName+".mmproj.gguf"),
	}
	for _, pattern := range patterns {
		if _, err := os.Stat(pattern); err == nil {
			fmt.Printf("[findMMProj] found: %s\n", pattern)
			return pattern
		}
	}

	// Fallback: scan directory for any file starting with "mmproj" and ending with ".gguf"
	entries, err := os.ReadDir(dir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := strings.ToLower(entry.Name())
			if strings.HasPrefix(name, "mmproj") && strings.HasSuffix(name, ".gguf") {
				found := filepath.Join(dir, entry.Name())
				fmt.Printf("[findMMProj] found via wildcard scan: %s\n", found)
				return found
			}
		}
	}

	fmt.Printf("[findMMProj] no mmproj found in %s\n", dir)
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

// spawnLlamaServerWithModel starts llama-server with a model and settings.
// stderrW, if non-nil, tees the process's stderr for diagnostic capture.
// fitParams controls the --fit flag; when false, "--fit off" is appended to
// bypass the automatic parameter-fitting phase that can crash with
// GGML_ASSERT(n_inputs < GGML_SCHED_MAX_SPLIT_INPUTS) on complex models.
func spawnLlamaServerWithModel(modelPath string, settings map[string]interface{}, stderrW io.Writer, fitParams bool) (*exec.Cmd, error) {
	serverPath := findLlamaServer()
	if serverPath == "" {
		return nil, fmt.Errorf("llama-server.exe not found")
	}

	args := []string{
		"--host", "127.0.0.1",
		"--port", "8080",
		"-m", modelPath,
		"-c", intSetting(settings, "context_size", 4096),
		"--threads", intSetting(settings, "threads", 8),
		"-b", intSetting(settings, "batch_size", 512),
		"-ngl", intSetting(settings, "gpu_layers", 0),
	}

	if !fitParams {
		args = append(args, "--fit", "off")
	}

	if mmprojPath := findMMProj(modelPath); mmprojPath != "" {
		args = append(args, "--mmproj", mmprojPath)
		fmt.Printf("INFO: Auto-detected mmproj file: %s\n", mmprojPath)
	} else {
		fmt.Printf("INFO: No mmproj file found for model: %s\n", modelPath)
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
	if v, ok := settings["flash_attention"].(bool); ok && v {
		args = append(args, "--flash-attn")
	}
	if ts := strSetting(settings, "tensor_split", ""); ts != "" && ts != "0" {
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

	cmd := exec.Command(serverPath, args...)
	applyHidden(cmd)
	cmd.Dir = filepath.Dir(serverPath)
	cmd.Stdout = os.Stdout
	if stderrW != nil {
		cmd.Stderr = io.MultiWriter(os.Stderr, stderrW)
	} else {
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start llama-server: %w", err)
	}
	proc.Attach(cmd)
	fmt.Printf("INFO: llama-server started (pid %d) with model %s\n", cmd.Process.Pid, modelPath)
	return cmd, nil
}

// ValidateLlamaServer runs llama-server --version to confirm it is executable.
func ValidateLlamaServer() error {
	serverPath := findLlamaServer()
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
