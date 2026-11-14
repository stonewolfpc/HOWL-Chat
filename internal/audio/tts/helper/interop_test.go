package helper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

var stage2HelperPath string

func TestMain(m *testing.M) {
	if os.Getenv("HOWL_HELPER_TEST_PROCESS") == "1" {
		os.Exit(m.Run())
	}
	root, err := moduleRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, "stage2 module root:", err)
		os.Exit(1)
	}
	dir, err := os.MkdirTemp("", "howl-stage2-helper-")
	if err != nil {
		fmt.Fprintln(os.Stderr, "stage2 temp:", err)
		os.Exit(1)
	}
	name := "howl-tts-helper-fake"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	stage2HelperPath = filepath.Join(dir, name)
	cmd := exec.Command("go", "build", "-o", stage2HelperPath, "./cmd/howl-tts-helper-fake")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "stage2 helper build failed: %v\n%s", err, out)
		_ = os.RemoveAll(dir)
		os.Exit(1)
	}
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

func moduleRoot() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime caller unavailable")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		return "", fmt.Errorf("go.mod not found at %s: %w", root, err)
	}
	return root, nil
}

func stage2Options() ClientOptions {
	return ClientOptions{
		HandshakeTimeout: 5 * time.Second,
		OperationTimeout: 5 * time.Second,
		CancelGrace:      500 * time.Millisecond,
		ShutdownTimeout:  3 * time.Second,
	}
}

func newStage2Unstarted(t *testing.T, mode string) (*Client, *Supervisor, int, string) {
	t.Helper()
	t.Setenv("HOWL_TTS_HELPER_PATH", stage2HelperPath)
	marker := t.TempDir()
	sup := NewSupervisor(SupervisorOptions{
		Discoverer:  Discoverer{},
		StderrLimit: 8192,
		Factory: func(ctx context.Context, spec LaunchSpec) (Process, error) {
			spec.Env = append(spec.Env,
				"HOWL_FAKE_HELPER_MODE="+mode,
				"HOWL_FAKE_HELPER_MARKER_DIR="+marker,
			)
			return StartProcess(ctx, spec)
		},
	})
	p, err := sup.Start(context.Background())
	if err != nil {
		t.Fatalf("real helper start: %v", err)
	}
	pid := p.PID()
	c, err := NewClient(sup, stage2Options())
	if err != nil {
		_ = sup.Terminate(context.Background(), time.Second)
		t.Fatalf("NewClient: %v", err)
	}
	t.Cleanup(func() {
		_ = c.Shutdown(context.Background())
		_ = sup.Terminate(context.Background(), time.Second)
		if !waitProcessGone(pid, 3*time.Second) {
			t.Errorf("stage2 helper pid %d survived cleanup", pid)
		}
	})
	return c, sup, pid, marker
}

func newStage2Client(t *testing.T, mode string) (*Client, *Supervisor, int, string) {
	t.Helper()
	c, sup, pid, marker := newStage2Unstarted(t, mode)
	if err := c.Start(context.Background()); err != nil {
		t.Fatalf("real helper handshake: %v", err)
	}
	return c, sup, pid, marker
}

func loadStage2(t *testing.T, c *Client) {
	t.Helper()
	root := t.TempDir()
	model := filepath.Join(root, "model.onnx")
	if err := os.WriteFile(model, []byte("synthetic fake model marker"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := c.Load(context.Background(), LoadRequest{
		ProviderID: "fake-go", Engine: "fake", ModelIdentity: "fake:model:1",
		Assets: LoadAssets{Root: root, Model: "model.onnx"},
	}); err != nil {
		t.Fatalf("load: %v", err)
	}
}

func waitMarker(t *testing.T, dir, name string) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	path := filepath.Join(dir, name)
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("marker %s not observed", name)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func waitStage(t *testing.T, sup *Supervisor, predicate func(types.TTSProviderStatus) bool) types.TTSProviderStatus {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for {
		status := sup.Status()
		if predicate(status) {
			return status
		}
		if time.Now().After(deadline) {
			t.Fatalf("stage2 status did not converge: %+v", status)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestStage2RealProcessHandshakeHealthCatalog(t *testing.T) {
	c, sup, _, _ := newStage2Client(t, "normal")
	status := sup.Status()
	if status.State != types.TTSStateUnloaded || status.ProcessState != types.TTSProcessAlive ||
		status.HelperVersion != "fake-1.0" || status.ProtocolVersion != "1.0" ||
		status.ProviderID != "fake-go" || len(status.Capabilities) != 5 {
		t.Fatalf("handshake status = %+v", status)
	}
	wantCaps := map[string]bool{
		CapabilityCatalog: true, CapabilityLoad: true, CapabilityUnload: true,
		CapabilitySynthesize: true, CapabilityCancel: true,
	}
	for _, cap := range status.Capabilities {
		if !wantCaps[cap.ID] || cap.Version != 1 {
			t.Fatalf("handshake advertised unexpected capability = %+v", cap)
		}
	}
	if stderr := sup.Stderr(); stderr != "" {
		t.Fatalf("helper wrote diagnostics on the normal path: %q", stderr)
	}
	health, err := c.Health(context.Background())
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	if !health.Process.Alive || health.State != "unloaded" || !health.AcceptingWork || health.QueueDepth != 0 {
		t.Fatalf("health = %+v", health)
	}
	catalog, err := c.Catalog(context.Background(), types.TTSProviderCatalogRequest{Engine: "fake"})
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	if catalog.ProviderID != "fake-go" || catalog.DefaultVoice != "fake-voice" ||
		catalog.CatalogVersion != 1 || len(catalog.Voices) != 1 {
		t.Fatalf("catalog = %+v", catalog)
	}
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	waitStage(t, sup, func(s types.TTSProviderStatus) bool { return !s.ProcessAlive })
}

func TestStage2RealProcessLoadSynthesizeUnload(t *testing.T) {
	c, sup, _, _ := newStage2Client(t, "normal")
	loadStage2(t, c)
	if got := sup.Status().State; got != types.TTSStateReady {
		t.Fatalf("state after load = %s", got)
	}
	dir := t.TempDir()
	result, err := c.Synthesize(context.Background(), SynthesisRequest{
		ModelIdentity: "fake:model:1", Text: "synthetic interoperability text",
		VoiceID: "fake-voice", Language: "en", Speed: 1, OutputDirectory: dir,
		Filename: "stage2.wav",
	})
	if err != nil {
		t.Fatalf("synthesize: %v", err)
	}
	final := filepath.Join(dir, "stage2.wav")
	if result.AudioPath != final || result.Metadata.SampleRate != 24000 ||
		result.Metadata.Channels != 1 || result.Metadata.BitDepth != 16 ||
		result.Metadata.Duration <= 0 {
		t.Fatalf("result = %+v", result)
	}
	raw, err := os.ReadFile(final)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if len(raw) != 44+1200*2 || string(raw[0:4]) != "RIFF" || string(raw[8:12]) != "WAVE" {
		t.Fatalf("output is not canonical synthetic WAV: %d bytes", len(raw))
	}
	sum := sha256.Sum256(raw)
	if len(hex.EncodeToString(sum[:])) != 64 {
		t.Fatal("hash length")
	}
	parts, err := filepath.Glob(final + ".part-*")
	if err != nil || len(parts) != 0 {
		t.Fatalf("temporary output remains: %v %v", parts, err)
	}
	if err := c.Unload(context.Background()); err != nil {
		t.Fatalf("unload: %v", err)
	}
	if got := sup.Status().State; got != types.TTSStateUnloaded {
		t.Fatalf("state after unload = %s", got)
	}
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}

func TestStage2CancellationAndControlResponsiveness(t *testing.T) {
	c, sup, _, marker := newStage2Client(t, "hold-synth")
	loadStage2(t, c)
	dir := t.TempDir()
	done := make(chan error, 1)
	go func() {
		_, err := c.Synthesize(context.Background(), SynthesisRequest{
			ModelIdentity: "fake:model:1", Text: "hold for cancellation", VoiceID: "fake-voice",
			Speed: 1, OutputDirectory: dir, Filename: "cancel.wav",
		})
		done <- err
	}()
	waitMarker(t, marker, "synth-started")
	health, err := c.Health(context.Background())
	if err != nil {
		t.Fatalf("health during synthesis: %v", err)
	}
	if health.State != "synthesizing" || health.AcceptingWork || health.ActiveRequestID == nil {
		t.Fatalf("health during synthesis = %+v", health)
	}
	if err := c.Cancel(context.Background(), *health.ActiveRequestID); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	select {
	case err := <-done:
		f, ok := AsFailure(err)
		if !ok || f.Code != "cancelled" {
			t.Fatalf("synthesis error = %v, want cancelled", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("cancelled synthesis did not terminate")
	}
	if _, err := os.Stat(filepath.Join(dir, "cancel.wav")); !os.IsNotExist(err) {
		t.Fatalf("cancelled output exists: %v", err)
	}
	parts, _ := filepath.Glob(filepath.Join(dir, "cancel.wav.part-*"))
	if len(parts) != 0 {
		t.Fatalf("cancelled temporary output remains: %v", parts)
	}
	if !sup.Status().ProcessAlive {
		t.Fatalf("helper died on cooperative cancel: %+v", sup.Status())
	}
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}

func TestStage2ForcedCancellationKillsAndCleans(t *testing.T) {
	c, sup, pid, marker := newStage2Client(t, "hang-synth")
	loadStage2(t, c)
	dir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := c.Synthesize(ctx, SynthesisRequest{
			ModelIdentity: "fake:model:1", Text: "force cancellation", VoiceID: "fake-voice",
			Speed: 1, OutputDirectory: dir, Filename: "forced.wav",
		})
		done <- err
	}()
	waitMarker(t, marker, "synth-started")
	cancel()
	select {
	case err := <-done:
		f, ok := AsFailure(err)
		if !ok || f.Code != "cancelled" {
			t.Fatalf("forced cancel error = %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("forced cancellation did not terminate")
	}
	if !waitProcessGone(pid, 3*time.Second) {
		t.Fatalf("helper pid %d survived forced cancellation", pid)
	}
	if _, err := os.Stat(filepath.Join(dir, "forced.wav")); !os.IsNotExist(err) {
		t.Fatalf("forced-cancel output exists: %v", err)
	}
	if sup.Status().ProcessAlive {
		t.Fatalf("supervisor reports alive after forced cancellation: %+v", sup.Status())
	}
}

func TestStage2TimeoutKillsAndCleans(t *testing.T) {
	c, sup, pid, marker := newStage2Client(t, "hang-synth")
	c.opts.OperationTimeout = 100 * time.Millisecond
	loadStage2(t, c)
	dir := t.TempDir()
	_, err := c.Synthesize(context.Background(), SynthesisRequest{
		ModelIdentity: "fake:model:1", Text: "timeout", VoiceID: "fake-voice",
		Speed: 1, OutputDirectory: dir, Filename: "timeout.wav",
	})
	f, ok := AsFailure(err)
	if !ok || f.Code != "timeout" {
		t.Fatalf("timeout error = %v", err)
	}
	waitMarker(t, marker, "synth-started")
	if !waitProcessGone(pid, 3*time.Second) {
		t.Fatalf("timed-out helper pid %d survived", pid)
	}
	if _, err := os.Stat(filepath.Join(dir, "timeout.wav")); !os.IsNotExist(err) {
		t.Fatalf("timed-out output exists: %v", err)
	}
	parts, _ := filepath.Glob(filepath.Join(dir, "timeout.wav.part-*"))
	if len(parts) != 0 {
		t.Fatalf("timed-out temporary output remains: %v", parts)
	}
	if sup.Status().ProcessAlive {
		t.Fatalf("supervisor reports alive after timeout: %+v", sup.Status())
	}
}

func TestStage2ShutdownObservesExit(t *testing.T) {
	c, sup, pid, _ := newStage2Client(t, "normal")
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	if !waitProcessGone(pid, 3*time.Second) {
		t.Fatalf("helper pid %d survived shutdown", pid)
	}
	status := waitStage(t, sup, func(s types.TTSProviderStatus) bool { return !s.ProcessAlive })
	if status.State != types.TTSStateAbsent && status.State != types.TTSStateUnavailable {
		t.Fatalf("state after shutdown = %+v", status)
	}
}

func TestStage2UnexpectedTermination(t *testing.T) {
	c, sup, pid, _ := newStage2Client(t, "crash-on-health")
	if _, err := c.Health(context.Background()); err == nil {
		t.Fatal("health unexpectedly succeeded after helper crash")
	}
	status := waitStage(t, sup, func(s types.TTSProviderStatus) bool {
		return s.State == types.TTSStateCrashed && !s.ProcessAlive
	})
	if status.ProcessState != types.TTSProcessExited {
		t.Fatalf("crash status = %+v", status)
	}
	if !waitProcessGone(pid, 3*time.Second) {
		t.Fatalf("crashed helper pid %d still exists", pid)
	}
}

func TestStage2MalformedHelperResponseFailsSafely(t *testing.T) {
	c, sup, pid, _ := newStage2Unstarted(t, "malformed-hello")
	err := c.Start(context.Background())
	f, ok := AsFailure(err)
	if !ok || f.Code != "invalid_frame" {
		t.Fatalf("start error = %v, want invalid_frame", err)
	}
	status := waitStage(t, sup, func(s types.TTSProviderStatus) bool { return !s.ProcessAlive })
	if status.State != types.TTSStateUnavailable {
		t.Fatalf("state after malformed response = %+v", status)
	}
	if !waitProcessGone(pid, 3*time.Second) {
		t.Fatalf("malformed helper pid %d still exists", pid)
	}
}

func TestStage2IncompatibleHelperFailsSafely(t *testing.T) {
	c, sup, pid, _ := newStage2Unstarted(t, "incompatible")
	err := c.Start(context.Background())
	f, ok := AsFailure(err)
	if !ok || f.Code != "protocol_incompatible" {
		t.Fatalf("start error = %v, want protocol_incompatible", err)
	}
	status := waitStage(t, sup, func(s types.TTSProviderStatus) bool { return !s.ProcessAlive })
	if status.State != types.TTSStateIncompatible {
		t.Fatalf("state after incompatible helper = %+v", status)
	}
	if !waitProcessGone(pid, 3*time.Second) {
		t.Fatalf("incompatible helper pid %d still exists", pid)
	}
}

func TestStage2ParentRejectsBadWAVHash(t *testing.T) {
	c, sup, _, _ := newStage2Client(t, "bad-wav")
	loadStage2(t, c)
	dir := t.TempDir()
	_, err := c.Synthesize(context.Background(), SynthesisRequest{
		ModelIdentity: "fake:model:1", Text: "bad hash", VoiceID: "fake-voice",
		Speed: 1, OutputDirectory: dir, Filename: "bad.wav",
	})
	f, ok := AsFailure(err)
	if !ok || f.Code != "output_failed" {
		t.Fatalf("bad-hash error = %v, want output_failed", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "bad.wav")); !os.IsNotExist(err) {
		t.Fatalf("rejected bad-hash output was not cleaned: %v", err)
	}
	if err := c.Shutdown(context.Background()); err != nil && !strings.Contains(err.Error(), "shutdown") {
		t.Fatalf("shutdown after bad hash: %v", err)
	}
	waitStage(t, sup, func(s types.TTSProviderStatus) bool { return !s.ProcessAlive })
}

func TestStage2MissingHelperIsNonfatal(t *testing.T) {
	dir := t.TempDir()
	sup := NewSupervisor(SupervisorOptions{
		Discoverer: Discoverer{
			LookupEnv:  func(string) (string, bool) { return "", false },
			Executable: func() (string, error) { return filepath.Join(dir, "howl.exe"), nil },
		},
	})
	c, err := NewClient(sup, stage2Options())
	if err != nil {
		t.Fatal(err)
	}
	err = c.Start(context.Background())
	f, ok := AsFailure(err)
	if !ok || f.Code != "helper_missing" {
		t.Fatalf("missing helper error = %v", err)
	}
	status := sup.Status()
	if status.ProcessAlive || status.State != types.TTSStateUnavailable {
		t.Fatalf("missing helper status = %+v", status)
	}
}

func TestStage2ReaderDoesNotLeakAfterShutdown(t *testing.T) {
	baseline := runtime.NumGoroutine()
	c, sup, pid, _ := newStage2Client(t, "normal")
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !waitProcessGone(pid, 3*time.Second) {
		t.Fatalf("helper pid %d survived shutdown", pid)
	}
	deadline := time.Now().Add(3 * time.Second)
	for runtime.NumGoroutine() > baseline+2 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if got := runtime.NumGoroutine(); got > baseline+2 {
		t.Fatalf("goroutines after shutdown = %d, baseline = %d", got, baseline)
	}
	if sup.Status().ProcessAlive {
		t.Fatalf("supervisor still reports alive: %+v", sup.Status())
	}
}
