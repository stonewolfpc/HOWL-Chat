// Package llama implements file transcription through the bundled
// llama-server's OpenAI-compatible audio route. The adapter owns one server
// process for the model/projector pair and releases it with the recognizer.
package llama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/audio/asr"
	"howl-chat/internal/audio/types"
	"howl-chat/internal/proc"
)

const startupTimeout = 120 * time.Second

// Recognizer is the llama-server ASR adapter. It does not download or discover
// models: the pipeline supplies resolved model/projector/server paths.
type Recognizer struct {
	*asr.BaseRecognizer

	cmd        *exec.Cmd
	waitCh     chan error
	baseURL    string
	httpClient *http.Client
}

// New creates a llama-server ASR recognizer for the local Qwen3-ASR identity.
func New() *Recognizer {
	return &Recognizer{
		BaseRecognizer: asr.NewBaseRecognizer(types.ASRQwen3ASR17BLocal),
		httpClient:     &http.Client{Timeout: 15 * time.Minute},
	}
}

// Initialize starts the server with the model and its matching audio projector.
// A missing pair, binary, or failed startup is a model-load failure; nothing is
// substituted.
func (r *Recognizer) Initialize(ctx context.Context, config types.ASRConfig) error {
	r.LockCalls()
	defer r.UnlockCalls()

	_ = r.stopLocked()

	if strings.TrimSpace(config.ModelPath) == "" {
		return types.NewAudioError(types.ErrCodeModelNotFound, "llama-server ASR model path is required", nil)
	}
	if _, err := os.Stat(config.ModelPath); err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "ASR model not found: "+config.ModelPath, err)
	}
	if strings.TrimSpace(config.ProjectorPath) == "" {
		return types.NewAudioError(types.ErrCodeModelNotFound, "audio projector path is required for llama-server ASR", nil)
	}
	if _, err := os.Stat(config.ProjectorPath); err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "audio projector not found: "+config.ProjectorPath, err)
	}
	serverPath := strings.TrimSpace(config.ServerPath)
	if serverPath == "" {
		return types.NewAudioError(types.ErrCodeModelNotFound, "llama-server path is required for ASR", nil)
	}
	if _, err := os.Stat(serverPath); err != nil {
		return types.NewAudioError(types.ErrCodeModelNotFound, "llama-server not found: "+serverPath, err)
	}

	port, err := freePort()
	if err != nil {
		return types.NewAudioError(types.ErrCodeModelLoadFailed, "cannot reserve a local port for llama-server ASR", err)
	}

	threads := config.NumThreads
	if threads <= 0 {
		threads = 6
	}
	args := []string{
		"-m", config.ModelPath,
		"--mmproj", config.ProjectorPath,
		"--host", "127.0.0.1",
		"--port", strconv.Itoa(port),
		"-c", "2048",
		"-ngl", "99",
		"--no-warmup",
		"-t", strconv.Itoa(threads),
	}

	var logs tailBuffer
	cmd := exec.Command(serverPath, args...)
	applyHidden(cmd)
	cmd.Dir = filepath.Dir(serverPath)
	cmd.Stdout = &logs
	cmd.Stderr = &logs
	// proc.Start creates the server suspended, assigns it to the application
	// job, and resumes it only after ownership succeeds. KILL_ON_JOB_CLOSE then
	// owns the server even if the app is force-killed before OnShutdown runs;
	// no code can run in the window before assignment. If ownership cannot be
	// established, the start fails instead of leaving an unowned server.
	if err := proc.Start(cmd); err != nil {
		return types.NewAudioError(types.ErrCodeModelLoadFailed, "failed to start llama-server for ASR under job ownership", err)
	}

	waitCh := make(chan error, 1)
	go func() { waitCh <- cmd.Wait() }()

	r.cmd = cmd
	r.waitCh = waitCh
	r.baseURL = fmt.Sprintf("http://127.0.0.1:%d", port)

	deadline := time.Now().Add(startupTimeout)
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			_ = r.stopLocked()
			return types.NewAudioError(types.ErrCodeModelLoadFailed, "llama-server ASR startup cancelled", ctx.Err())
		case err := <-waitCh:
			// The server exited on its own; make sure a descendant cannot
			// outlive the process that failed to become ready.
			if cmd.Process != nil {
				proc.KillTree(cmd.Process.Pid)
			}
			r.cmd = nil
			r.waitCh = nil
			r.baseURL = ""
			return types.NewAudioError(types.ErrCodeModelLoadFailed, "llama-server exited before becoming ready: "+logs.tail(2000), err)
		case <-ticker.C:
			if time.Now().After(deadline) {
				_ = r.stopLocked()
				return types.NewAudioError(types.ErrCodeModelLoadFailed, "llama-server did not become ready within 120s: "+logs.tail(2000), nil)
			}
			if r.healthy(ctx) {
				r.SetInitialized(config)
				return nil
			}
		}
	}
}

// Transcribe posts one audio file to /v1/audio/transcriptions and returns clean
// text with the server's language prefix removed.
func (r *Recognizer) Transcribe(ctx context.Context, audioPath string, callback types.ProgressCallback) (*types.RecognitionResult, error) {
	r.LockCalls()
	defer r.UnlockCalls()

	if err := r.CheckInitialized(); err != nil {
		return nil, err
	}
	if r.baseURL == "" {
		return nil, types.NewAudioError(types.ErrCodeModelLoadFailed, "llama-server ASR is not running", nil)
	}
	file, err := os.Open(audioPath)
	if err != nil {
		return nil, types.NewAudioError(types.ErrCodeFormatNotSupported, "cannot open audio file: "+audioPath, err)
	}
	defer file.Close()

	if callback != nil {
		callback(5, "Preparing audio...")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	_ = writer.WriteField("model", "qwen3-asr")
	_ = writer.WriteField("response_format", "json")
	if err := writer.Close(); err != nil {
		return nil, err
	}

	if callback != nil {
		callback(35, "Transcribing...")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/v1/audio/transcriptions", &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, types.NewAudioError(types.ErrCodeModelLoadFailed, "llama-server transcription request failed", err)
	}
	defer resp.Body.Close()
	payload, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode != http.StatusOK {
		return nil, types.NewAudioError(
			types.ErrCodeModelLoadFailed,
			fmt.Sprintf("llama-server transcription failed (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(payload))),
			nil,
		)
	}

	var decoded struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return nil, types.NewAudioError(types.ErrCodeFormatNotSupported, "invalid transcription response", err)
	}
	text, language := CleanTranscript(decoded.Text)
	if text == "" {
		return nil, types.NewAudioError(types.ErrCodeFormatNotSupported, "llama-server returned an empty transcript", nil)
	}
	if callback != nil {
		callback(100, "Transcription complete")
	}
	return &types.RecognitionResult{Text: text, Language: language}, nil
}

// TranscribeStream is not implemented: E4 proved file transcription only.
func (r *Recognizer) TranscribeStream(ctx context.Context, chunks <-chan types.AudioChunk, results chan<- types.RecognitionResult) error {
	return types.NewAudioError(types.ErrCodeFormatNotSupported, "streaming transcription is not implemented for llama-server ASR", nil)
}

// Release stops the server and clears initialized state. It is idempotent.
func (r *Recognizer) Release() error {
	r.LockCalls()
	defer r.UnlockCalls()
	return r.stopLocked()
}

func (r *Recognizer) stopLocked() error {
	if r.cmd == nil {
		r.ClearInitialized()
		return nil
	}
	cmd := r.cmd
	waitCh := r.waitCh
	pid := 0
	if cmd.Process != nil {
		pid = cmd.Process.Pid
		// Kill the tree, not just the direct child: a descendant must not
		// survive the recognizer that owns it. Process.Kill is kept as a
		// fallback for the (unexpected) case where taskkill cannot run.
		proc.KillTree(pid)
		_ = cmd.Process.Kill()
	}
	var stopErr error
	if waitCh != nil {
		select {
		case <-waitCh:
		case <-time.After(5 * time.Second):
			// Do not report a clean release when the process is still alive.
			stopErr = fmt.Errorf("llama-server ASR (pid %d) did not exit within 5s", pid)
		}
	}
	r.cmd = nil
	r.waitCh = nil
	r.baseURL = ""
	r.ClearInitialized()
	return stopErr
}

func (r *Recognizer) healthy(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL+"/health", nil)
	if err != nil {
		return false
	}
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode == http.StatusOK
}

func freePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

// tailBuffer keeps the tail of server output for startup diagnostics without
// letting a long-running server grow memory.
type tailBuffer struct {
	mu   sync.Mutex
	data []byte
}

func (b *tailBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = append(b.data, p...)
	if len(b.data) > 16*1024 {
		b.data = b.data[len(b.data)-16*1024:]
	}
	return len(p), nil
}

func (b *tailBuffer) tail(limit int) string {
	b.mu.Lock()
	defer b.mu.Unlock()
	data := b.data
	if len(data) > limit {
		data = data[len(data)-limit:]
	}
	return strings.TrimSpace(string(data))
}
