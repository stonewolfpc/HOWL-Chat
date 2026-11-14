package helper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type fakeProcess struct {
	requests   chan []byte
	raw        chan []byte
	closed     chan struct{}
	closeOnce  sync.Once
	waitOnce   sync.Once
	waitCount  atomic.Int32
	killCount  atomic.Int32
	closeStops bool
	waitErr    error
	waitHook   func()
	stderr     string

	mu   sync.Mutex
	sent [][]byte
}

func newFakeProcess() *fakeProcess {
	return &fakeProcess{
		requests:   make(chan []byte, 64),
		raw:        make(chan []byte, 64),
		closed:     make(chan struct{}),
		closeStops: true,
	}
}

func (p *fakeProcess) Send(ctx context.Context, frame []byte) error {
	select {
	case <-p.closed:
		return io.ErrClosedPipe
	default:
	}
	p.mu.Lock()
	p.sent = append(p.sent, append([]byte(nil), frame...))
	p.mu.Unlock()
	select {
	case p.requests <- append([]byte(nil), frame...):
		return nil
	case <-p.closed:
		return io.ErrClosedPipe
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *fakeProcess) Receive() ([]byte, error) {
	select {
	case line := <-p.raw:
		return line, nil
	case <-p.closed:
		return nil, io.EOF
	}
}

func (p *fakeProcess) Close() error {
	if p.closeStops {
		p.closeOnce.Do(func() { close(p.closed) })
	}
	return nil
}

func (p *fakeProcess) Stderr() string { return p.stderr }

func (p *fakeProcess) Wait() error {
	p.waitCount.Add(1)
	p.waitOnce.Do(func() {
		<-p.closed
		if p.waitHook != nil {
			p.waitHook()
		}
	})
	return p.waitErr
}

func (p *fakeProcess) Kill() error {
	p.killCount.Add(1)
	p.closeOnce.Do(func() { close(p.closed) })
	return nil
}

func (p *fakeProcess) PID() int { return 4242 }

func (p *fakeProcess) pushRaw(line []byte) error {
	select {
	case p.raw <- append([]byte(nil), line...):
		return nil
	case <-p.closed:
		return io.ErrClosedPipe
	}
}

func (p *fakeProcess) sentCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.sent)
}

type fakeHelper struct {
	p *fakeProcess

	mu         sync.Mutex
	hello      helloResult
	helloErr   *wireError
	loadCount  int
	loaded     bool
	active     string
	hold       chan struct{}
	release    chan struct{}
	cancelCh   chan struct{}
	started    chan string
	cancelSeen chan string
}

func newFakeHelper() *fakeHelper {
	h := &fakeHelper{
		p:          newFakeProcess(),
		hello:      defaultHello(),
		started:    make(chan string, 8),
		cancelSeen: make(chan string, 8),
	}
	go h.loop()
	return h
}

func defaultHello() helloResult {
	return helloResult{
		HelperVersion: "fake-1.0",
		Providers: []providerDescriptor{{
			ProviderID: "fake-go", ProviderVersion: "1.0", Engine: "fake",
			DisplayName: "Fake TTS", Runtime: map[string]string{"onnxruntime": "unreported", "g2p": "unreported"},
			Capabilities: []capability{
				{ID: CapabilitySynthesize, Version: 1}, {ID: CapabilitySpeed, Version: 1},
				{ID: CapabilityLanguage, Version: 1},
			},
		}},
		Capabilities: []capability{
			{ID: CapabilityCatalog, Version: 1}, {ID: CapabilityLoad, Version: 1},
			{ID: CapabilityUnload, Version: 1}, {ID: CapabilitySynthesize, Version: 1},
			{ID: CapabilityCancel, Version: 1},
		},
		Extensions: []extension{},
		Limits:     limits{MaxFrameBytes: MaxFrameBytes, MaxTextBytes: MaxTextBytes},
	}
}

func (h *fakeHelper) loop() {
	for {
		select {
		case frame := <-h.p.requests:
			go h.handle(frame)
		case <-h.p.closed:
			return
		}
	}
}

func (h *fakeHelper) handle(frame []byte) {
	env, err := decodeRequestFrame(frame)
	if err != nil {
		h.sendError("", "invalid_request", "request frame is invalid")
		return
	}
	switch env.Op {
	case "hello":
		h.mu.Lock()
		errObj := h.helloErr
		result := h.hello
		h.mu.Unlock()
		if errObj != nil {
			h.sendWireError(env.ID, errObj)
			return
		}
		h.send(env.ID, result)
	case "health":
		h.send(env.ID, h.health())
	case "catalog":
		var params catalogParams
		if err := decodeParams(env.Params, &params); err != nil {
			h.sendError(env.ID, "invalid_request", "catalog params invalid")
			return
		}
		h.send(env.ID, catalogResult{
			ProviderID: params.ProviderID, Engine: params.Engine, CatalogVersion: 1, Available: true,
			DefaultVoice: "fake-voice", Languages: []string{"en"},
			Speed:  speedRange{Supported: true, Min: 0.5, Max: 2.0},
			Voices: []voiceInfo{{VoiceID: "fake-voice", Name: "Fake Voice", Language: "en", Available: true}},
		})
	case "load":
		var params loadParams
		if err := decodeParams(env.Params, &params); err != nil {
			h.sendError(env.ID, "invalid_request", "load params invalid")
			return
		}
		h.mu.Lock()
		h.loadCount++
		h.loaded = true
		h.mu.Unlock()
		h.send(env.ID, loadResult{
			ProviderID: params.ProviderID, Engine: params.Engine, ProviderVersion: "1.0",
			ModelIdentity: params.ModelIdentity, Loaded: true, AlreadyLoaded: false,
			Runtime: map[string]string{"onnxruntime": "unreported"},
		})
	case "unload":
		h.mu.Lock()
		h.loaded = false
		h.mu.Unlock()
		h.send(env.ID, unloadResult{Unloaded: true})
	case "cancel":
		var params cancelParams
		if err := decodeParams(env.Params, &params); err != nil {
			h.sendError(env.ID, "invalid_request", "cancel params invalid")
			return
		}
		h.mu.Lock()
		cancelCh := h.cancelCh
		h.mu.Unlock()
		select {
		case h.cancelSeen <- params.RequestID:
		default:
		}
		if cancelCh != nil {
			select {
			case <-cancelCh:
			default:
				close(cancelCh)
			}
		}
		h.send(env.ID, cancelResult{RequestID: params.RequestID, CancelRequested: true})
	case "shutdown":
		h.send(env.ID, shutdownResult{Accepted: true})
		go func() {
			time.Sleep(5 * time.Millisecond)
			_ = h.p.Close()
		}()
	case "synthesize":
		h.handleSynthesis(env)
	default:
		h.sendError(env.ID, "unsupported_operation", "operation is not supported")
	}
}

func (h *fakeHelper) health() healthResult {
	h.mu.Lock()
	loaded, active := h.loaded, h.active
	h.mu.Unlock()
	state, ready, accepting := "unloaded", false, true
	if loaded {
		state, ready = "ready", true
	}
	if active != "" {
		state, ready, accepting = "synthesizing", true, false
	}
	var activePtr *string
	if active != "" {
		value := active
		activePtr = &value
	}
	return healthResult{
		Process: healthProcess{Alive: true, PID: 4242, UptimeMS: 1},
		State:   state, Ready: ready, AcceptingWork: accepting,
		ActiveRequestID: activePtr, QueueDepth: 0, LastError: nil,
	}
}

func (h *fakeHelper) handleSynthesis(env requestEnvelope) {
	var params synthesisParams
	if err := decodeParams(env.Params, &params); err != nil {
		h.sendError(env.ID, "invalid_request", "synthesis params invalid")
		return
	}
	cancelCh := make(chan struct{})
	h.mu.Lock()
	h.active = env.ID
	h.cancelCh = cancelCh
	hold, release := h.hold, h.release
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		if h.active == env.ID {
			h.active = ""
			h.cancelCh = nil
		}
		h.mu.Unlock()
	}()
	select {
	case h.started <- env.ID:
	default:
	}

	if hold != nil {
		select {
		case <-release:
		case <-cancelCh:
			h.sendError(env.ID, "cancelled", "synthesis was cancelled")
			return
		case <-h.p.closed:
			return
		}
	}

	path := filepath.Join(params.Output.Directory, params.Output.Filename)
	hash, size, err := writeTestWAV(path, 24000, 240)
	if err != nil {
		h.sendError(env.ID, "output_failed", "fake output failed")
		return
	}
	h.send(env.ID, synthesisResult{
		OutputPath: path, Format: "wav", Codec: "pcm_s16le", SampleRateHz: 24000,
		Channels: 1, BitDepth: 16, DurationSeconds: 0.01, Bytes: size, SHA256: hash,
		ModelIdentity: params.ModelIdentity,
	})
}

func (h *fakeHelper) setHold() chan struct{} {
	release := make(chan struct{})
	h.mu.Lock()
	h.hold = release
	h.release = release
	h.mu.Unlock()
	return release
}

func (h *fakeHelper) setHelloError(code, message string) {
	h.mu.Lock()
	h.helloErr = &wireError{Code: code, Message: message, Retryable: false}
	h.mu.Unlock()
}

func (h *fakeHelper) loads() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.loadCount
}

func (h *fakeHelper) send(id string, result any) {
	raw, _ := json.Marshal(result)
	ok := true
	frame, _ := json.Marshal(responseEnvelope{Protocol: ProtocolVersion, ID: id, OK: &ok, Result: raw})
	_ = h.p.pushRaw(append(frame, '\n'))
}

func (h *fakeHelper) sendError(id, code, message string) {
	h.sendWireError(id, &wireError{Code: code, Message: message, Retryable: false})
}

func (h *fakeHelper) sendWireError(id string, e *wireError) {
	ok := false
	frame, _ := json.Marshal(responseEnvelope{Protocol: ProtocolVersion, ID: id, OK: &ok, Error: e})
	_ = h.p.pushRaw(append(frame, '\n'))
}

func newFakeClient(t *testing.T, h *fakeHelper, opts ClientOptions) (*Client, *Supervisor) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "howl-tts-helper.exe")
	if err := os.WriteFile(path, []byte("test double"), 0o644); err != nil {
		t.Fatalf("write fake path: %v", err)
	}
	sup := NewSupervisor(SupervisorOptions{
		Discoverer: Discoverer{LookupEnv: func(string) (string, bool) { return path, true }},
		Factory: func(context.Context, LaunchSpec) (Process, error) {
			return h.p, nil
		},
	})
	c, err := NewClient(sup, opts)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	t.Cleanup(func() {
		_ = sup.Terminate(context.Background(), 50*time.Millisecond)
	})
	return c, sup
}

func testWAV(t *testing.T, path string) (string, int64) {
	t.Helper()
	hash, size, err := writeTestWAV(path, 24000, 240)
	if err != nil {
		t.Fatalf("write WAV: %v", err)
	}
	return hash, size
}

func writeTestWAV(path string, sampleRate, samples int) (string, int64, error) {
	if sampleRate <= 0 {
		sampleRate = 24000
	}
	if samples <= 0 {
		samples = 240
	}
	dataBytes := samples * 2
	raw := make([]byte, 44+dataBytes)
	copy(raw[0:4], "RIFF")
	put32(raw[4:8], uint32(36+dataBytes))
	copy(raw[8:12], "WAVE")
	copy(raw[12:16], "fmt ")
	put32(raw[16:20], 16)
	put16(raw[20:22], 1)
	put16(raw[22:24], 1)
	put32(raw[24:28], uint32(sampleRate))
	put32(raw[28:32], uint32(sampleRate*2))
	put16(raw[32:34], 2)
	put16(raw[34:36], 16)
	copy(raw[36:40], "data")
	put32(raw[40:44], uint32(dataBytes))
	sum := sha256.Sum256(raw)
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(sum[:]), int64(len(raw)), nil
}

func put16(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func put32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

var _ = bytes.Equal
