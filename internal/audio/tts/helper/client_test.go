package helper

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

func fastClientOptions() ClientOptions {
	return ClientOptions{
		HandshakeTimeout: time.Second,
		OperationTimeout: time.Second,
		CancelGrace:      50 * time.Millisecond,
		ShutdownTimeout:  50 * time.Millisecond,
	}
}

func TestClientLifecycleAndCatalog(t *testing.T) {
	h := newFakeHelper()
	c, sup := newFakeClient(t, h, fastClientOptions())
	defer c.Shutdown(context.Background())

	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	status, err := c.ProviderStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if status.State != types.TTSStateUnloaded || status.ProcessState != types.TTSProcessAlive {
		t.Fatalf("status = %+v", status)
	}
	if h.loads() != 0 {
		t.Fatalf("status loaded %d models, want 0", h.loads())
	}

	catalog, err := c.ProviderCatalog(context.Background(), types.TTSProviderCatalogRequest{Engine: "fake"})
	if err != nil {
		t.Fatal(err)
	}
	if catalog.ProviderID != "fake-go" || catalog.DefaultVoice != "fake-voice" || catalog.Speed == nil || !catalog.Speed.Supported {
		t.Fatalf("catalog = %+v", catalog)
	}
	if h.loads() != 0 {
		t.Fatalf("catalog loaded %d models, want 0", h.loads())
	}

	assets := t.TempDir()
	if err := c.Load(context.Background(), LoadRequest{
		ProviderID: "fake-go", Engine: "fake", ModelIdentity: "model:1",
		Assets: LoadAssets{Root: assets, Model: "model.onnx"},
	}); err != nil {
		t.Fatal(err)
	}
	if h.loads() != 1 {
		t.Fatalf("load count = %d, want 1", h.loads())
	}
	if status := sup.Status(); status.State != types.TTSStateReady || !status.Ready {
		t.Fatalf("status after load = %+v", status)
	}
	health, err := c.Health(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !health.Process.Alive || health.State != "ready" {
		t.Fatalf("health = %+v", health)
	}
	if err := c.Unload(context.Background()); err != nil {
		t.Fatal(err)
	}
	if sup.Status().State != types.TTSStateUnloaded {
		t.Fatalf("status after unload = %+v", sup.Status())
	}
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestClientRejectsMissingServiceCapability(t *testing.T) {
	h := newFakeHelper()
	c, _ := newFakeClient(t, h, fastClientOptions())
	defer c.Shutdown(context.Background())
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	c.mu.Lock()
	var kept []capability
	for _, cap := range c.hello.Capabilities {
		if cap.ID != CapabilityCatalog {
			kept = append(kept, cap)
		}
	}
	c.hello.Capabilities = kept
	c.mu.Unlock()

	before := h.p.sentCount()
	_, err := c.Catalog(context.Background(), types.TTSProviderCatalogRequest{Engine: "fake"})
	f, ok := AsFailure(err)
	if !ok || f.Code != "unsupported_capability" {
		t.Fatalf("error = %v, want unsupported_capability", err)
	}
	if got := h.p.sentCount(); got != before {
		t.Fatalf("catalog frame was sent without the catalog capability: %d -> %d", before, got)
	}
}

func TestClientRejectsUnadvertisedProviderFeatures(t *testing.T) {
	tests := []struct {
		name         string
		providerCaps []capability
		speed        float64
		language     string
	}{
		{
			name: "synthesis",
			providerCaps: []capability{
				{ID: CapabilitySpeed, Version: 1},
			},
			speed: 1,
		},
		{
			name: "speed",
			providerCaps: []capability{
				{ID: CapabilitySynthesize, Version: 1},
			},
			speed: 1.25,
		},
		{
			name: "language",
			providerCaps: []capability{
				{ID: CapabilitySynthesize, Version: 1},
			},
			speed:    1,
			language: "en",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newFakeHelper()
			h.mu.Lock()
			h.hello.Providers[0].Capabilities = tt.providerCaps
			h.mu.Unlock()
			c, _ := newFakeClient(t, h, fastClientOptions())
			defer c.Shutdown(context.Background())
			if err := c.Load(context.Background(), LoadRequest{
				ProviderID: "fake-go", Engine: "fake", ModelIdentity: "model:1",
				Assets: LoadAssets{Root: t.TempDir(), Model: "model.onnx"},
			}); err != nil {
				t.Fatal(err)
			}

			before := h.p.sentCount()
			_, err := c.Synthesize(context.Background(), SynthesisRequest{
				ModelIdentity: "model:1", Text: "hello", VoiceID: "fake-voice",
				Language: tt.language, Speed: tt.speed,
				OutputDirectory: t.TempDir(), Filename: "unsupported.wav",
			})
			f, ok := AsFailure(err)
			if !ok || f.Code != "unsupported_capability" {
				t.Fatalf("error = %v, want unsupported_capability", err)
			}
			if got := h.p.sentCount(); got != before {
				t.Fatalf("synthesis frame was sent without the required capability: %d -> %d", before, got)
			}
		})
	}
}

func TestConstructionDoesNotStartHelper(t *testing.T) {
	h := newFakeHelper()
	c, _ := newFakeClient(t, h, fastClientOptions())
	if got := h.p.waitCount.Load(); got != 0 {
		t.Fatalf("Wait calls before explicit Start = %d, want 0", got)
	}
	if got := h.p.sentCount(); got != 0 {
		t.Fatalf("frames before explicit Start = %d, want 0", got)
	}
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := h.p.waitCount.Load(); got != 1 {
		t.Fatalf("Wait calls after explicit Start = %d, want 1", got)
	}
	if got := h.p.sentCount(); got == 0 {
		t.Fatal("expected handshake traffic after explicit Start")
	}
}

func TestShutdownWithActiveRequest(t *testing.T) {
	h := newFakeHelper()
	h.setHold()
	c, _ := newFakeClient(t, h, fastClientOptions())
	dir := t.TempDir()
	done := make(chan error, 1)
	go func() {
		_, err := c.Synthesize(context.Background(), SynthesisRequest{
			ModelIdentity: "model:1", Text: "hello", VoiceID: "fake-voice", Speed: 1,
			OutputDirectory: dir, Filename: "active-shutdown.wav",
		})
		done <- err
	}()
	select {
	case <-h.started:
	case <-time.After(time.Second):
		t.Fatal("synthesis did not start")
	}
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("active synthesis did not stop after shutdown")
	}
}

func TestSynthesisSuccess(t *testing.T) {
	h := newFakeHelper()
	c, _ := newFakeClient(t, h, fastClientOptions())
	dir := t.TempDir()
	result, err := c.Synthesize(context.Background(), SynthesisRequest{
		ModelIdentity: "model:1", Text: "hello", VoiceID: "fake-voice", Speed: 1,
		OutputDirectory: dir, Filename: "out.wav",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.AudioPath != filepath.Join(dir, "out.wav") {
		t.Fatalf("path = %s", result.AudioPath)
	}
	if result.Metadata.SampleRate != 24000 || result.Metadata.Channels != 1 || result.Metadata.BitDepth != 16 {
		t.Fatalf("metadata = %+v", result.Metadata)
	}
	if _, err := os.Stat(result.AudioPath); err != nil {
		t.Fatal(err)
	}
}

func TestCancellationBeforeTerminal(t *testing.T) {
	h := newFakeHelper()
	h.setHold()
	c, _ := newFakeClient(t, h, fastClientOptions())
	dir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := c.Synthesize(ctx, SynthesisRequest{
			ModelIdentity: "model:1", Text: "hello", VoiceID: "fake-voice", Speed: 1,
			OutputDirectory: dir, Filename: "cancelled.wav",
		})
		done <- err
	}()
	select {
	case <-h.started:
	case <-time.After(time.Second):
		t.Fatal("synthesis did not start")
	}
	cancel()
	select {
	case err := <-done:
		f, ok := AsFailure(err)
		if !ok || f.Code != "cancelled" {
			t.Fatalf("error = %v, want cancelled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("synthesis did not stop")
	}
	if _, err := os.Stat(filepath.Join(dir, "cancelled.wav")); !os.IsNotExist(err) {
		t.Fatalf("cancelled output exists: %v", err)
	}
}

func TestCancellationAfterTerminalBeforeAcceptance(t *testing.T) {
	h := newFakeHelper()
	c, _ := newFakeClient(t, h, fastClientOptions())
	dir := t.TempDir()
	validating := make(chan struct{})
	releaseValidation := make(chan struct{})
	c.validationHook = func() {
		close(validating)
		<-releaseValidation
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := c.Synthesize(ctx, SynthesisRequest{
			ModelIdentity: "model:1", Text: "hello", VoiceID: "fake-voice", Speed: 1,
			OutputDirectory: dir, Filename: "late-cancel.wav",
		})
		done <- err
	}()
	select {
	case <-validating:
	case <-time.After(time.Second):
		t.Fatal("synthesis did not reach validation")
	}
	cancel()
	select {
	case <-h.cancelSeen:
	case <-time.After(time.Second):
		t.Fatal("cancellation was not sent")
	}
	close(releaseValidation)
	select {
	case err := <-done:
		f, ok := AsFailure(err)
		if !ok || f.Code != "cancelled" {
			t.Fatalf("error = %v, want cancelled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("synthesis did not finish")
	}
	if _, err := os.Stat(filepath.Join(dir, "late-cancel.wav")); !os.IsNotExist(err) {
		t.Fatalf("cancelled accepted output exists: %v", err)
	}
}

func TestAcceptedResultIsNotRetractedByLaterCancellation(t *testing.T) {
	h := newFakeHelper()
	c, _ := newFakeClient(t, h, fastClientOptions())
	dir := t.TempDir()
	validating := make(chan struct{})
	releaseValidation := make(chan struct{})
	c.validationHook = func() {
		close(validating)
		<-releaseValidation
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := c.Synthesize(ctx, SynthesisRequest{
			ModelIdentity: "model:1", Text: "hello", VoiceID: "fake-voice", Speed: 1,
			OutputDirectory: dir, Filename: "accepted.wav",
		})
		done <- err
	}()
	select {
	case <-validating:
	case <-time.After(time.Second):
		t.Fatal("synthesis did not reach validation")
	}
	close(releaseValidation)
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("accepted synthesis failed: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("synthesis did not finish")
	}
	cancel()
	if _, err := os.Stat(filepath.Join(dir, "accepted.wav")); err != nil {
		t.Fatalf("accepted output removed: %v", err)
	}
}

func TestTimeoutForcesTermination(t *testing.T) {
	h := newFakeHelper()
	h.setHold()
	h.p.closeStops = false
	opts := fastClientOptions()
	opts.OperationTimeout = 50 * time.Millisecond
	opts.ShutdownTimeout = 20 * time.Millisecond
	c, _ := newFakeClient(t, h, opts)
	dir := t.TempDir()
	_, err := c.Synthesize(context.Background(), SynthesisRequest{
		ModelIdentity: "model:1", Text: "hello", VoiceID: "fake-voice", Speed: 1,
		OutputDirectory: dir, Filename: "timeout.wav",
	})
	f, ok := AsFailure(err)
	if !ok || f.Code != "timeout" {
		t.Fatalf("error = %v, want timeout", err)
	}
	if h.p.killCount.Load() == 0 {
		t.Fatal("timed-out helper was not killed")
	}
}

func TestControlMessagesDuringSynthesis(t *testing.T) {
	h := newFakeHelper()
	release := h.setHold()
	c, _ := newFakeClient(t, h, fastClientOptions())
	dir := t.TempDir()
	done := make(chan error, 1)
	go func() {
		_, err := c.Synthesize(context.Background(), SynthesisRequest{
			ModelIdentity: "model:1", Text: "hello", VoiceID: "fake-voice", Speed: 1,
			OutputDirectory: dir, Filename: "control.wav",
		})
		done <- err
	}()
	var id string
	select {
	case id = <-h.started:
	case <-time.After(time.Second):
		t.Fatal("synthesis did not start")
	}
	if _, err := c.Health(context.Background()); err != nil {
		t.Fatalf("health during synthesis: %v", err)
	}
	if err := c.Cancel(context.Background(), id); err != nil {
		t.Fatalf("cancel during synthesis: %v", err)
	}
	select {
	case <-release:
	default:
		close(release)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("synthesis did not finish")
	}
}

func TestMalformedFrameIsFatal(t *testing.T) {
	h := newFakeHelper()
	c, _ := newFakeClient(t, h, fastClientOptions())
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := h.p.pushRaw([]byte("not json\n")); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for {
		c.mu.Lock()
		f := c.readFailure
		c.mu.Unlock()
		if f != nil {
			if f.Code != "invalid_frame" {
				t.Fatalf("failure = %+v, want invalid_frame", f)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("malformed frame did not fail the session")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestUnknownResponseIDIsProtocolViolation(t *testing.T) {
	h := newFakeHelper()
	c, _ := newFakeClient(t, h, fastClientOptions())
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	ok := true
	frame, _ := json.Marshal(responseEnvelope{
		Protocol: ProtocolVersion, ID: "req-unknown", OK: &ok, Result: json.RawMessage(`{}`),
	})
	if err := h.p.pushRaw(append(frame, '\n')); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for {
		c.mu.Lock()
		f := c.readFailure
		c.mu.Unlock()
		if f != nil {
			if f.Code != "protocol_violation" {
				t.Fatalf("failure = %+v, want protocol_violation", f)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("unknown ID did not fail the session")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestMissingHelperIsOptional(t *testing.T) {
	sup := NewSupervisor(SupervisorOptions{
		Discoverer: Discoverer{
			LookupEnv:  func(string) (string, bool) { return "", false },
			Executable: func() (string, error) { return filepath.Join(t.TempDir(), "howl.exe"), nil },
		},
	})
	c, err := NewClient(sup, fastClientOptions())
	if err != nil {
		t.Fatal(err)
	}
	err = c.Start(context.Background())
	f, ok := AsFailure(err)
	if !ok || f.Code != "helper_missing" {
		t.Fatalf("error = %v, want helper_missing", err)
	}
	status, statusErr := c.ProviderStatus(context.Background())
	if statusErr == nil {
		t.Fatal("status unexpectedly succeeded")
	}
	if status.State != types.TTSStateUnavailable || status.ProcessAlive {
		t.Fatalf("status = %+v", status)
	}
}

func TestNoGoroutineLeakAfterShutdown(t *testing.T) {
	baseline := runtime.NumGoroutine()
	h := newFakeHelper()
	c, _ := newFakeClient(t, h, fastClientOptions())
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > baseline && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if got := runtime.NumGoroutine(); got > baseline {
		t.Fatalf("goroutines after shutdown = %d, baseline = %d", got, baseline)
	}
}

func TestProtocolIncompatibleStopsHandshake(t *testing.T) {
	h := newFakeHelper()
	h.setHelloError("protocol_incompatible", "helper supports only major 2")
	c, _ := newFakeClient(t, h, fastClientOptions())
	err := c.Start(context.Background())
	f, ok := AsFailure(err)
	if !ok || f.Code != "protocol_incompatible" {
		t.Fatalf("error = %v, want protocol_incompatible", err)
	}
}
