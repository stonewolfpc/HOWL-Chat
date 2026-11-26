package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"howl-chat/internal/backend/llama"
)

const fakeLlamaHelp = `--help
--fit
--flash-attn
--cache-type-k
--cache-type-v
--kv-offload
--no-kv-offload
--load-mode
--main-gpu
--tensor-split
--chat-template
--reasoning
--reasoning-format
--reasoning-budget
--reasoning-budget-message
--mmproj
--vocab-only
`

const fakeLlamaVersion = "version: 0.8.0-dev\nbuild 12015\ncommit a34de8618\n"

type fakeLlamaAttempt struct {
	Args []string `json:"args"`
}

// runFakeLlamaServerIfRequested lets the test binary act as the spawned
// llama-server. Only the child processes started by the runtime test have
// HOWL_FAKE_LLAMA=1; normal `go test` processes never enter this branch.
func runFakeLlamaServerIfRequested(args []string) (bool, int) {
	if os.Getenv("HOWL_FAKE_LLAMA") != "1" {
		return false, 0
	}
	for _, arg := range args[1:] {
		switch arg {
		case "--help":
			fmt.Print(fakeLlamaHelp)
			return true, 0
		case "--version":
			fmt.Print(fakeLlamaVersion)
			return true, 0
		}
	}

	attempt := recordFakeAttempt(args[1:])
	mode := os.Getenv("HOWL_FAKE_LLAMA_MODE")
	cacheK := fakeArg(args[1:], "--cache-type-k")
	cacheV := fakeArg(args[1:], "--cache-type-v")

	switch mode {
	case "unrelated":
		fmt.Fprintln(os.Stderr, "0.00.100.000 E llama_init_from_model: failed to initialize: synthetic unrelated startup failure")
		return true, 1
	case "near_miss":
		fmt.Fprintln(os.Stderr, "0.00.100.000 E llama_init_from_model: quantized V cache requires flash_attn")
		return true, 1
	case "f16_reject":
		fmt.Fprintln(os.Stderr, "0.00.572.074 E llama_init_from_model: quantized V cache requires flash_attn to be enabled")
		return true, 1
	case "retry_fails":
		if attempt == 1 {
			fmt.Fprintln(os.Stderr, "0.00.572.074 E llama_init_from_model: quantized V cache requires flash_attn to be enabled")
			return true, 1
		}
		fmt.Fprintln(os.Stderr, "0.00.700.000 E llama_init_from_model: failed to create context with model")
		return true, 1
	case "f16":
		if cacheK != "f16" || cacheV != "f16" {
			fmt.Fprintln(os.Stderr, "0.00.572.074 E llama_init_from_model: quantized V cache requires flash_attn to be enabled")
			return true, 1
		}
		return serveFakeHealth()
	default: // reject_once
		if attempt == 1 {
			fmt.Fprintln(os.Stderr, "0.00.572.074 E llama_init_from_model: quantized V cache requires flash_attn to be enabled")
			return true, 1
		}
		if cacheK != "f16" || cacheV != "f16" {
			fmt.Fprintln(os.Stderr, "0.00.572.074 E llama_init_from_model: quantized V cache requires flash_attn to be enabled")
			return true, 1
		}
		return serveFakeHealth()
	}
}

func recordFakeAttempt(args []string) int {
	path := os.Getenv("HOWL_FAKE_LLAMA_STATE")
	if path == "" {
		return 1
	}
	data, _ := os.ReadFile(path)
	attempt := bytes.Count(data, []byte("\n")) + 1
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fake llama-server: attempt log:", err)
		return attempt
	}
	defer file.Close()
	_ = json.NewEncoder(file).Encode(fakeLlamaAttempt{Args: args})
	return attempt
}

func fakeArg(args []string, name string) string {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == name {
			return args[i+1]
		}
	}
	return ""
}

func serveFakeHealth() (bool, int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"status":"ok"}`)
	})
	server := &http.Server{Addr: "127.0.0.1:8080", Handler: mux}
	fmt.Println("fake llama-server: listening on http://127.0.0.1:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintln(os.Stderr, "fake llama-server:", err)
		return true, 1
	}
	return true, 0
}

type capturedLoadEvent struct {
	Name string
	Args []interface{}
}

func runFakeLoad(t *testing.T, cacheType, mode, flashMode string) ([]capturedLoadEvent, []fakeLlamaAttempt) {
	t.Helper()

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("test executable: %v", err)
	}
	statePath := filepath.Join(t.TempDir(), "attempts.jsonl")
	modelPath := writeLoadTestGGUF(t, t.TempDir(), "runtime-kv-test.gguf", map[string]any{
		"general.architecture": "llama",
		"general.name":         "Runtime KV test model",
	})

	previousResolver := findLlamaServerForRuntime
	previousSink := devEventSink
	findLlamaServerForRuntime = func() string { return exe }
	t.Setenv("HOWL_FAKE_LLAMA", "1")
	t.Setenv("HOWL_FAKE_LLAMA_MODE", mode)
	t.Setenv("HOWL_FAKE_LLAMA_STATE", statePath)

	app := NewApp()
	app.loaderDiagnosticsDir = t.TempDir()
	if app.settings == nil {
		app.settings = newTestSettingsAuthority(t, t.TempDir())
	}
	modelSettings := map[string]interface{}{
		"context_size":  2048,
		"gpu_layers":    0,
		"batch_size":    256,
		"threads":       2,
		"kv_cache_type": cacheType,
		"offload_kqv":   false,
		"use_mmap":      true,
		"use_mlock":     false,
	}
	if flashMode != "" {
		modelSettings["flash_attention_mode"] = flashMode
	}
	app.settings.SetModel(modelSettings)
	app.engine.SetClient(llama.NewHTTPClient("127.0.0.1", 8080, true))

	events := make(chan capturedLoadEvent, 32)
	devEventSink = func(name string, args ...interface{}) {
		select {
		case events <- capturedLoadEvent{Name: name, Args: args}:
		default:
		}
	}
	t.Cleanup(func() {
		findLlamaServerForRuntime = previousResolver
		devEventSink = previousSink
		_ = app.UnloadModel()
	})

	if err := app.LoadModel(modelPath); err != nil {
		t.Fatalf("LoadModel returned before monitoring: %v", err)
	}

	var got []capturedLoadEvent
	deadline := time.After(20 * time.Second)
	for {
		select {
		case event := <-events:
			got = append(got, event)
			if event.Name == "model:loaded" || event.Name == "model:error" {
				return got, readFakeAttempts(t, statePath)
			}
		case <-deadline:
			t.Fatalf("timed out waiting for model:loaded/model:error; events=%v", got)
		}
	}
}

func readFakeAttempts(t *testing.T, path string) []fakeLlamaAttempt {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		t.Fatalf("read attempts: %v", err)
	}
	var attempts []fakeLlamaAttempt
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var attempt fakeLlamaAttempt
		if err := json.Unmarshal(line, &attempt); err != nil {
			t.Fatalf("decode attempt %q: %v", line, err)
		}
		attempts = append(attempts, attempt)
	}
	return attempts
}

func terminalEvent(events []capturedLoadEvent, name string) capturedLoadEvent {
	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Name == name {
			return events[i]
		}
	}
	return capturedLoadEvent{}
}

func warningText(events []capturedLoadEvent) string {
	for _, event := range events {
		if event.Name == "model:warning" && len(event.Args) > 0 {
			return fmt.Sprint(event.Args[0])
		}
	}
	return ""
}

func withoutKVCacheArgs(args []string) string {
	var out []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--cache-type-k" || args[i] == "--cache-type-v" {
			i++
			continue
		}
		out = append(out, args[i])
	}
	return strings.Join(out, "|")
}

func requireOnlyKVCacheChanged(t *testing.T, before, after []string) {
	t.Helper()
	if withoutKVCacheArgs(before) != withoutKVCacheArgs(after) {
		t.Fatalf("retry changed an unrelated argument\nbefore: %v\nafter:  %v", before, after)
	}
}

func TestLoadModelQuantizedKVRecoveryRuntime(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("test executable: %v", err)
	}
	cmd := exec.Command(exe, "-test.run=^TestRuntimeQuantizedKV", "-test.v", "-test.count=1")
	cmd.Env = append(os.Environ(), "HOWL_RUNTIME_CHILD=1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 120 {
			lines = lines[len(lines)-120:]
		}
		t.Fatalf("isolated runtime child failed: %v\n%s", err, strings.Join(lines, "\n"))
	}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "--- PASS") || strings.Contains(line, "--- FAIL") {
			t.Log(line)
		}
	}
}

func requireRuntimeChild(t *testing.T) {
	t.Helper()
	if os.Getenv("HOWL_RUNTIME_CHILD") != "1" {
		t.Skip("runs through TestLoadModelQuantizedKVRecoveryRuntime in a fresh process")
	}
}

func TestRuntimeQuantizedKVRecoveryReachesHealthyRetry(t *testing.T) {
	requireRuntimeChild(t)
	for _, cacheType := range []string{"q4_0", "q8_0"} {
		t.Run(cacheType, func(t *testing.T) {
			events, attempts := runFakeLoad(t, cacheType, "reject_once", "auto")
			if terminalEvent(events, "model:loaded").Name == "" {
				t.Fatalf("retry did not reach healthy: %v", events)
			}
			if len(attempts) != 2 {
				t.Fatalf("attempts = %d, want initial + one retry: %v", len(attempts), attempts)
			}
			first, retry := attempts[0].Args, attempts[1].Args
			t.Logf("initial args: %v", first)
			t.Logf("retry args: %v", retry)
			if fakeArg(first, "--cache-type-k") != cacheType || fakeArg(first, "--cache-type-v") != cacheType {
				t.Fatalf("initial spawn did not carry %s: %v", cacheType, first)
			}
			if fakeArg(retry, "--cache-type-k") != "f16" || fakeArg(retry, "--cache-type-v") != "f16" {
				t.Fatalf("retry did not switch both KV caches to f16: %v", retry)
			}
			requireOnlyKVCacheChanged(t, first, retry)
			warning := warningText(events)
			t.Logf("warning: %s", warning)
			t.Logf("terminal event: %s", terminalEvent(events, "model:loaded").Name)
			if !strings.Contains(warning, cacheType) || !strings.Contains(warning, "f16") || !strings.Contains(warning, "quantized V cache requires flash_attn") {
				t.Fatalf("warning missing request/replacement/reason: %q", warning)
			}
			if !strings.Contains(strings.Join(first, " "), "--flash-attn auto") {
				t.Fatalf("initial spawn lost Flash Attention auto: %v", first)
			}
		})
	}
}

func TestRuntimeQuantizedKVF16RequestNeverRetries(t *testing.T) {
	requireRuntimeChild(t)
	events, attempts := runFakeLoad(t, "f16", "f16_reject", "auto")
	if terminalEvent(events, "model:error").Name == "" {
		t.Fatalf("f16 rejection did not report model:error: %v", events)
	}
	if terminalEvent(events, "model:loaded").Name != "" {
		t.Fatalf("f16 rejection reported success: %v", events)
	}
	if len(attempts) != 1 {
		t.Fatalf("f16 request attempted %d spawns, want 1: %v", len(attempts), attempts)
	}
	if warningText(events) != "" {
		t.Fatalf("f16 request emitted a recovery warning: %q", warningText(events))
	}
}

func TestRuntimeQuantizedKVNearMissRejectionDoesNotRetry(t *testing.T) {
	requireRuntimeChild(t)
	events, attempts := runFakeLoad(t, "q4_0", "near_miss", "auto")
	if terminalEvent(events, "model:error").Name == "" {
		t.Fatalf("near-miss failure did not report model:error: %v", events)
	}
	if len(attempts) != 1 {
		t.Fatalf("near-miss rejection attempted %d spawns, want 1: %v", len(attempts), attempts)
	}
	if warningText(events) != "" {
		t.Fatalf("near-miss rejection emitted a quantized-KV warning: %q", warningText(events))
	}
}

func TestRuntimeQuantizedKVOffAndUnsetDowngradeBeforeSpawn(t *testing.T) {
	requireRuntimeChild(t)
	for _, tc := range []struct {
		name      string
		flashMode string
		want      string
	}{
		{name: "off", flashMode: "off", want: "Flash Attention is off"},
		{name: "unset", flashMode: "", want: "Flash Attention is not set to On"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			events, attempts := runFakeLoad(t, "q4_0", "f16", tc.flashMode)
			if terminalEvent(events, "model:loaded").Name == "" {
				t.Fatalf("pre-spawn f16 load did not reach healthy: %v", events)
			}
			if len(attempts) != 1 {
				t.Fatalf("pre-spawn downgrade attempted %d spawns, want 1: %v", len(attempts), attempts)
			}
			args := attempts[0].Args
			if fakeArg(args, "--cache-type-k") != "f16" || fakeArg(args, "--cache-type-v") != "f16" {
				t.Fatalf("pre-spawn downgrade did not send f16: %v", args)
			}
			warning := warningText(events)
			if !strings.Contains(warning, tc.want) || !strings.Contains(warning, "q4_0") {
				t.Fatalf("pre-spawn warning = %q, want %q naming q4_0", warning, tc.want)
			}
		})
	}
}

func TestRuntimeQuantizedKVUnrelatedFailureDoesNotRetry(t *testing.T) {
	requireRuntimeChild(t)
	events, attempts := runFakeLoad(t, "q4_0", "unrelated", "auto")
	if terminalEvent(events, "model:error").Name == "" {
		t.Fatalf("unrelated failure did not report model:error: %v", events)
	}
	if len(attempts) != 1 {
		t.Fatalf("unrelated failure attempted %d spawns, want 1: %v", len(attempts), attempts)
	}
	if warningText(events) != "" {
		t.Fatalf("unrelated failure emitted a quantized-KV warning: %q", warningText(events))
	}
}

func TestRuntimeQuantizedKVFailedRetryReportsFailureWithoutLooping(t *testing.T) {
	requireRuntimeChild(t)
	events, attempts := runFakeLoad(t, "q4_0", "retry_fails", "auto")
	if terminalEvent(events, "model:error").Name == "" {
		t.Fatalf("failed retry did not report model:error: %v", events)
	}
	if terminalEvent(events, "model:loaded").Name != "" {
		t.Fatalf("failed retry reported success: %v", events)
	}
	if len(attempts) != 2 {
		t.Fatalf("failed retry attempted %d spawns, want exactly 2: %v", len(attempts), attempts)
	}
	warning := warningText(events)
	t.Logf("warning: %s", warning)
	if event := terminalEvent(events, "model:error"); event.Name != "" {
		t.Logf("error event: %v", event.Args)
	}
	if !strings.Contains(warning, "Retrying once with f16 KV") {
		t.Fatalf("failed retry did not surface the recovery warning: %q", warning)
	}
}
