package fakehelper

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	modeNormal              = "normal"
	modeHoldSynthesis       = "hold-synth"
	modeMalformedHello      = "malformed-hello"
	modeIncompatible        = "incompatible"
	modeCrashOnHealth       = "crash-on-health"
	modeMalformedHealth     = "malformed-health"
	modeBadWAV              = "bad-wav"
	modeHangSynthesis       = "hang-synth"
	modePartialCapabilities = "partial-capabilities"
	fakeProviderID          = "fake-go"
	fakeEngine              = "fake"
	fakeProviderVersion     = "0.1.0"
	fakeHelperVersion       = "fake-1.0"
	fakeVoice               = "fake-voice"
	fakeUnavailableVoice    = "fake-unavailable"
)

type workerRequest struct {
	env    requestEnvelope
	cancel chan struct{}
}

type server struct {
	in        *bufio.Reader
	out       *bufio.Writer
	diag      io.Writer
	mode      string
	markerDir string
	started   time.Time

	writeMu sync.Mutex

	mu             sync.Mutex
	handshake      bool
	stopping       bool
	loaded         bool
	modelIdentity  string
	busy           bool
	activeID       string
	activeKind     string
	cancel         chan struct{}
	inflight       map[string]struct{}
	workerCh       chan workerRequest
	workerDone     chan struct{}
	stopCh         chan struct{}
	stopOnce       sync.Once
	cancelOnce     sync.Once
	shutdownWaited bool
}

// Run executes one Protocol v1 helper session over the supplied stdio. It
// returns the process exit code defined by the frozen contract.
func Run(in io.Reader, out, diag io.Writer, mode, markerDir string) int {
	if mode == "" {
		mode = modeNormal
	}
	s := &server{
		in:         bufio.NewReaderSize(in, 64*1024),
		out:        bufio.NewWriterSize(out, 64*1024),
		diag:       diag,
		mode:       mode,
		markerDir:  markerDir,
		started:    time.Now(),
		inflight:   make(map[string]struct{}),
		workerCh:   make(chan workerRequest, 1),
		workerDone: make(chan struct{}),
		stopCh:     make(chan struct{}),
	}
	go s.worker()
	code := s.loop()
	s.stop()
	s.waitWorker(2 * time.Second)
	return code
}

func (s *server) loop() int {
	for {
		line, frameErr := readFrame(s.in)
		if frameErr != nil {
			s.writeError("", frameErr)
			return 4
		}
		if line == nil {
			return 0
		}
		env, reqErr, fatal := decodeRequestFrame(line)
		if reqErr != nil {
			if fatal {
				s.writeError(env.ID, reqErr)
				return 4
			}
			s.respondError(env.ID, reqErr)
			continue
		}

		if s.duplicate(env.ID) {
			s.writeError(env.ID, newWireError("duplicate_request_id", "request ID is already in flight"))
			continue
		}
		s.markInflight(env.ID)

		if env.Op == "hello" {
			if code := s.handleHello(env); code != 0 {
				return code
			}
			continue
		}
		if !s.handshaken() {
			s.respondError(env.ID, newWireError("invalid_request", "hello must be the first operation"))
			continue
		}

		switch env.Op {
		case "health":
			if code := s.handleHealth(env); code != 0 {
				return code
			}
		case "cancel":
			if !s.hasServiceCapability(CapabilityCancel) {
				s.respondError(env.ID, unsupportedCapability(CapabilityCancel))
				continue
			}
			s.handleCancel(env)
		case "shutdown":
			return s.handleShutdown(env)
		case "catalog", "load", "unload", "synthesize":
			if !s.hasServiceCapability(serviceCapabilityForOperation(env.Op)) {
				s.respondError(env.ID, unsupportedCapability(serviceCapabilityForOperation(env.Op)))
				continue
			}
			if !s.enqueue(env) {
				s.respondError(env.ID, newWireError("busy", "one non-control operation is already active"))
			}
		default:
			s.respondError(env.ID, newWireError("unsupported_operation", "operation is not supported"))
		}
	}
}

func (s *server) handleHello(env requestEnvelope) int {
	if s.handshaken() {
		s.respondError(env.ID, newWireError("invalid_request", "hello was already completed"))
		return 0
	}
	var p helloParams
	if e := decodeParams(env.Params, &p); e != nil {
		s.respondError(env.ID, e)
		return 0
	}
	if p.Client != "howl-chat" || p.ClientVersion == "" || len(p.ClientVersion) > 64 ||
		!p.ProtocolMin.valid() || !p.ProtocolMax.valid() ||
		p.ProtocolMin.Major > p.ProtocolMax.Major ||
		(p.ProtocolMin.Major == p.ProtocolMax.Major && p.ProtocolMin.Minor > p.ProtocolMax.Minor) {
		s.respondError(env.ID, newWireError("invalid_request", "hello parameters are invalid"))
		return 0
	}
	if s.mode == modeIncompatible {
		s.respondError(env.ID, newWireError("protocol_incompatible", "fake helper only supports major 2"))
		time.Sleep(10 * time.Millisecond)
		return 4
	}
	if s.mode == modeMalformedHello {
		_ = s.writeRaw([]byte("not a protocol frame"))
		time.Sleep(20 * time.Millisecond)
		return 4
	}
	if !(p.ProtocolMin.Major <= protocolMajor && protocolMajor <= p.ProtocolMax.Major) {
		s.respondError(env.ID, newWireError("protocol_incompatible", "no common protocol major"))
		return 4
	}
	if p.ProtocolMin.Major == protocolMajor && p.ProtocolMin.Minor > protocolMinor {
		s.respondError(env.ID, newWireError("protocol_minor_unsupported", "no common protocol minor"))
		return 4
	}

	s.mu.Lock()
	s.handshake = true
	s.mu.Unlock()
	s.respondResult(env.ID, helloResult{
		HelperVersion: fakeHelperVersion,
		Providers: []providerDescriptor{{
			ProviderID: fakeProviderID, ProviderVersion: fakeProviderVersion, Engine: fakeEngine,
			DisplayName: "Fake TTS", Runtime: map[string]string{"onnxruntime": "unreported", "g2p": "unreported"},
			Capabilities: s.providerCapabilities(),
		}},
		Capabilities: s.serviceCapabilities(),
		Extensions:   []extension{},
		Limits:       limits{MaxFrameBytes: maxFrameBytes, MaxTextBytes: maxTextBytes},
	})
	return 0
}

func (s *server) serviceCapabilities() []capability {
	return []capability{
		{ID: CapabilityCatalog, Version: 1}, {ID: CapabilityLoad, Version: 1},
		{ID: CapabilityUnload, Version: 1}, {ID: CapabilitySynthesize, Version: 1},
		{ID: CapabilityCancel, Version: 1},
	}
}

func (s *server) providerCapabilities() []capability {
	if s.mode == modePartialCapabilities {
		return []capability{{ID: CapabilitySynthesize, Version: 1}}
	}
	return []capability{
		{ID: CapabilitySynthesize, Version: 1}, {ID: CapabilitySpeed, Version: 1},
		{ID: CapabilityLanguage, Version: 1},
	}
}

func (s *server) hasServiceCapability(id string) bool {
	return hasCapability(s.serviceCapabilities(), id)
}

func (s *server) hasProviderCapability(id string) bool {
	return hasCapability(s.providerCapabilities(), id)
}

func hasCapability(caps []capability, id string) bool {
	if id == "" {
		return false
	}
	for _, c := range caps {
		if c.ID == id {
			return true
		}
	}
	return false
}

func serviceCapabilityForOperation(op string) string {
	switch op {
	case "catalog":
		return CapabilityCatalog
	case "load":
		return CapabilityLoad
	case "unload":
		return CapabilityUnload
	case "synthesize":
		return CapabilitySynthesize
	default:
		return ""
	}
}

func unsupportedCapability(id string) *wireError {
	return &wireError{
		Code: "unsupported_capability", Message: "helper does not advertise the required capability",
		Detail: map[string]string{"capability": id},
	}
}

func (s *server) handleHealth(env requestEnvelope) int {
	if s.mode == modeCrashOnHealth {
		return 7
	}
	if s.mode == modeMalformedHealth {
		_ = s.writeRaw([]byte("health is not a frame"))
		return 4
	}
	var p healthParams
	if e := decodeParams(env.Params, &p); e != nil {
		s.respondError(env.ID, e)
		return 0
	}
	s.mu.Lock()
	state := "unloaded"
	ready, accepting := false, true
	if s.loaded {
		state, ready = "ready", true
	}
	if s.busy {
		accepting = false
		switch s.activeKind {
		case "load":
			state = "loading"
		case "synthesize":
			state = "synthesizing"
		}
	}
	if s.stopping {
		state, ready, accepting = "stopping", false, false
	}
	var active *string
	if state == "synthesizing" && s.activeID != "" {
		value := s.activeID
		active = &value
	}
	s.mu.Unlock()
	s.respondResult(env.ID, healthResult{
		Process: healthProcess{Alive: true, PID: os.Getpid(), UptimeMS: time.Since(s.started).Milliseconds()},
		State:   state, Ready: ready, AcceptingWork: accepting,
		ActiveRequestID: active, QueueDepth: 0, LastError: nil,
	})
	return 0
}

func (s *server) handleCancel(env requestEnvelope) {
	var p cancelParams
	if e := decodeParams(env.Params, &p); e != nil {
		s.respondError(env.ID, e)
		return
	}
	if p.RequestID == "" || !requestIDPattern.MatchString(p.RequestID) {
		s.respondError(env.ID, newWireError("invalid_request", "cancel target is invalid").withField("params.request_id"))
		return
	}
	s.mu.Lock()
	activeID, activeKind := s.activeID, s.activeKind
	cancel := s.cancel
	if activeID == p.RequestID && activeKind == "synthesize" && cancel != nil {
		s.cancelOnce.Do(func() { close(cancel) })
		s.mu.Unlock()
		s.respondResult(env.ID, cancelResult{RequestID: p.RequestID, CancelRequested: true})
		return
	}
	s.mu.Unlock()
	if activeID == p.RequestID && activeKind != "synthesize" {
		s.respondError(env.ID, newWireError("request_not_cancellable", "target is not a synthesis"))
		return
	}
	s.respondError(env.ID, newWireError("request_not_found", "cancel target is not active"))
}

func (s *server) handleShutdown(env requestEnvelope) int {
	var p struct{}
	if e := decodeParams(env.Params, &p); e != nil {
		s.respondError(env.ID, e)
		return 0
	}
	s.mu.Lock()
	s.stopping = true
	cancel := s.cancel
	s.mu.Unlock()
	if cancel != nil {
		s.cancelOnce.Do(func() { close(cancel) })
	}
	s.respondResult(env.ID, shutdownResult{Accepted: true})
	s.stop()
	s.waitWorker(2 * time.Second)
	return 0
}

func (s *server) enqueue(env requestEnvelope) bool {
	s.mu.Lock()
	if s.stopping || s.busy {
		s.mu.Unlock()
		return false
	}
	cancel := make(chan struct{})
	s.busy = true
	s.activeID = env.ID
	s.activeKind = env.Op
	s.cancel = cancel
	s.mu.Unlock()

	select {
	case s.workerCh <- workerRequest{env: env, cancel: cancel}:
		return true
	case <-s.stopCh:
		s.finishWork()
		return false
	}
}

func (s *server) worker() {
	defer close(s.workerDone)
	for {
		select {
		case req := <-s.workerCh:
			s.handleWork(req)
			s.finishWork()
		case <-s.stopCh:
			return
		}
	}
}

func (s *server) handleWork(req workerRequest) {
	switch req.env.Op {
	case "catalog":
		s.handleCatalog(req.env)
	case "load":
		s.handleLoad(req.env)
	case "unload":
		s.handleUnload(req.env)
	case "synthesize":
		s.handleSynthesize(req.env, req.cancel)
	}
}

func (s *server) handleCatalog(env requestEnvelope) {
	var p catalogParams
	if e := decodeParams(env.Params, &p); e != nil {
		s.respondError(env.ID, e)
		return
	}
	if p.ProviderID != fakeProviderID || p.Engine != fakeEngine {
		s.respondError(env.ID, newWireError("provider_unknown", "fake provider was not requested"))
		return
	}
	voices := []voiceInfo{{
		VoiceID: fakeVoice, Name: "Fake Voice", Language: "en", Available: true,
	}}
	if p.IncludeUnavailable {
		voices = append(voices, voiceInfo{
			VoiceID: fakeUnavailableVoice, Name: "Unavailable Fake Voice", Language: "en",
			Available: false, Detail: "intentionally unavailable test voice",
		})
	}
	s.respondResult(env.ID, catalogResult{
		ProviderID: fakeProviderID, Engine: fakeEngine, CatalogVersion: 1, Available: true,
		Detail:       "deterministic fake catalog; no model is loaded",
		DefaultVoice: fakeVoice, Languages: []string{"en"},
		Speed:         speedRange{Supported: true, Min: 0.5, Max: 2.0},
		SSMLSupported: false, StreamingSupported: false, Voices: voices,
	})
}

func (s *server) handleLoad(env requestEnvelope) {
	var p loadParams
	if e := decodeParams(env.Params, &p); e != nil {
		s.respondError(env.ID, e)
		return
	}
	if p.ProviderID != fakeProviderID || p.Engine != fakeEngine {
		s.respondError(env.ID, newWireError("provider_unknown", "fake provider was not requested"))
		return
	}
	if p.ModelIdentity == "" || len(p.ModelIdentity) > 512 {
		s.respondError(env.ID, newWireError("invalid_request", "model identity is invalid").withField("params.model_identity"))
		return
	}
	if e := validateLoadAssets(p.Assets); e != nil {
		s.respondError(env.ID, e)
		return
	}

	s.mu.Lock()
	already := false
	switch {
	case !s.loaded:
		s.loaded = true
		s.modelIdentity = p.ModelIdentity
	case s.modelIdentity == p.ModelIdentity:
		already = true
	default:
		s.mu.Unlock()
		s.respondError(env.ID, newWireError("provider_already_loaded", "unload the current fake provider first"))
		return
	}
	s.mu.Unlock()
	s.respondResult(env.ID, loadResult{
		ProviderID: fakeProviderID, Engine: fakeEngine, ProviderVersion: fakeProviderVersion,
		ModelIdentity: p.ModelIdentity, Loaded: true, AlreadyLoaded: already,
		Runtime: map[string]string{"onnxruntime": "unreported", "g2p": "unreported"},
	})
}

func (s *server) handleUnload(env requestEnvelope) {
	var p struct{}
	if e := decodeParams(env.Params, &p); e != nil {
		s.respondError(env.ID, e)
		return
	}
	s.mu.Lock()
	s.loaded = false
	s.modelIdentity = ""
	s.mu.Unlock()
	s.respondResult(env.ID, unloadResult{Unloaded: true})
}

func (s *server) handleSynthesize(env requestEnvelope, cancel <-chan struct{}) {
	var p synthesisParams
	if e := decodeParams(env.Params, &p); e != nil {
		s.respondError(env.ID, e)
		return
	}
	s.mu.Lock()
	loaded, identity := s.loaded, s.modelIdentity
	s.mu.Unlock()
	if !loaded {
		s.respondError(env.ID, newWireError("provider_not_loaded", "load the fake provider first"))
		return
	}
	if p.ModelIdentity != identity {
		s.respondError(env.ID, newWireError("model_identity_mismatch", "synthesis identity does not match the loaded identity"))
		return
	}
	if !s.hasServiceCapability(CapabilitySynthesize) || !s.hasProviderCapability(CapabilitySynthesize) {
		s.respondError(env.ID, unsupportedCapability(CapabilitySynthesize))
		return
	}
	if p.Speed != 1.0 && !s.hasProviderCapability(CapabilitySpeed) {
		s.respondError(env.ID, unsupportedCapability(CapabilitySpeed))
		return
	}
	if p.Language != "" && !s.hasProviderCapability(CapabilityLanguage) {
		s.respondError(env.ID, unsupportedCapability(CapabilityLanguage))
		return
	}
	if e := validateText(p.Text); e != nil {
		s.respondError(env.ID, e)
		return
	}
	if p.VoiceID == fakeUnavailableVoice {
		s.respondError(env.ID, newWireError("voice_missing", "fake voice is intentionally unavailable").withField("params.voice_id"))
		return
	}
	if p.VoiceID != fakeVoice {
		s.respondError(env.ID, newWireError("voice_unsupported", "voice is not in the fake catalog").withField("params.voice_id"))
		return
	}
	if p.Language != "" && p.Language != "en" {
		s.respondError(env.ID, newWireError("language_unsupported", "fake provider only advertises en").withField("params.language"))
		return
	}
	if !finite(p.Speed) || p.Speed < 0.5 || p.Speed > 2.0 {
		s.respondError(env.ID, newWireError("speed_unsupported", "speed is outside the fake range").withField("params.speed"))
		return
	}
	if e := validateOutputPath(p.Output.Directory, p.Output.Filename); e != nil {
		s.respondError(env.ID, e)
		return
	}

	s.mark("synth-started")
	if s.mode == modeHoldSynthesis || s.mode == modeHangSynthesis {
		if !s.waitForReleaseOrCancel(cancel) {
			s.respondError(env.ID, newWireError("cancelled", "synthesis was cancelled"))
			return
		}
	}

	final := filepath.Join(filepath.Clean(p.Output.Directory), p.Output.Filename)
	part, err := createPart(final)
	if err != nil {
		s.respondError(env.ID, newWireError("output_failed", "could not create temporary output"))
		return
	}
	hash, size, duration, err := writeSyntheticWAV(part)
	if err != nil {
		_ = os.Remove(part)
		s.respondError(env.ID, newWireError("output_failed", "could not write fake WAV"))
		return
	}
	if err := validateSyntheticWAV(part, hash, size); err != nil {
		_ = os.Remove(part)
		s.respondError(env.ID, newWireError("output_failed", "fake WAV validation failed"))
		return
	}
	if canceled(cancel) {
		_ = os.Remove(part)
		s.respondError(env.ID, newWireError("cancelled", "synthesis was cancelled"))
		return
	}
	if err := publishNoReplace(part, final); err != nil {
		_ = os.Remove(part)
		s.respondError(env.ID, newWireError("output_failed", "could not publish fake WAV"))
		return
	}
	// Publication has completed. If cancellation is observed before the
	// terminal, return cancelled without mutating the published file; the
	// parent owns cleanup after any live helper is terminated.
	if canceled(cancel) {
		s.respondError(env.ID, newWireError("cancelled", "synthesis was cancelled after publication"))
		return
	}
	resultHash := hash
	if s.mode == modeBadWAV {
		resultHash = strings.Repeat("0", 64)
	}
	s.respondResult(env.ID, synthesisResult{
		OutputPath: final, Format: "wav", Codec: "pcm_s16le", SampleRateHz: syntheticRate,
		Channels: 1, BitDepth: 16, DurationSeconds: duration, Bytes: size, SHA256: resultHash,
		ModelIdentity: p.ModelIdentity,
	})
}

func (s *server) waitForReleaseOrCancel(cancel <-chan struct{}) bool {
	if s.mode == modeHangSynthesis {
		select {
		case <-s.stopCh:
			return false
		case <-time.After(5 * time.Second):
			return true
		}
	}
	if s.markerDir == "" {
		select {
		case <-cancel:
			return false
		case <-s.stopCh:
			return false
		case <-time.After(2 * time.Second):
			return true
		}
	}
	release := filepath.Join(s.markerDir, "synth-release")
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-cancel:
			return false
		case <-s.stopCh:
			return false
		case <-ticker.C:
			if _, err := os.Stat(release); err == nil {
				return true
			}
		}
	}
}

func validateLoadAssets(a loadAssets) *wireError {
	if !filepath.IsAbs(a.Root) || filepath.Clean(a.Root) != a.Root {
		return newWireError("invalid_request", "assets.root must be absolute and cleaned").withField("params.assets.root")
	}
	if info, err := os.Stat(a.Root); err != nil || !info.IsDir() {
		return newWireError("model_missing", "assets root is unavailable").withField("params.assets.root")
	}
	if reparse, err := hasReparseComponent(a.Root); err != nil || reparse {
		return newWireError("model_missing", "assets root contains a reparse point").withField("params.assets.root")
	}
	for name, value := range map[string]string{"model": a.Model, "voices": a.Voices, "config": a.Config, "data": a.Data} {
		if value == "" {
			continue
		}
		if filepath.IsAbs(value) || filepath.Clean(value) != value || value == ".." ||
			strings.HasPrefix(value, ".."+string(filepath.Separator)) {
			return newWireError("invalid_request", "asset path escapes root").withField("params.assets." + name)
		}
		full := filepath.Join(a.Root, value)
		info, err := os.Lstat(full)
		if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
			return newWireError("model_missing", "fake asset file is missing").withField("params.assets." + name)
		}
		reparse, err := hasReparseComponent(full)
		if err != nil || reparse {
			return newWireError("model_missing", "fake asset path contains a reparse point").withField("params.assets." + name)
		}
	}
	return nil
}

func validateOutputPath(dir, filename string) *wireError {
	if !filepath.IsAbs(dir) || filepath.Clean(dir) != dir || len(dir) > maxDirectory {
		return newWireError("invalid_request", "output directory must be absolute, cleaned, and bounded").withField("params.output.directory")
	}
	if len(filename) > maxFilename || !filenamePattern.MatchString(filename) ||
		filepath.Base(filename) != filename || strings.Contains(filename, "..") {
		return newWireError("invalid_request", "output filename is invalid").withField("params.output.filename")
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return newWireError("output_failed", "output directory is unavailable").withField("params.output.directory")
	}
	if reparse, err := hasReparseComponent(dir); err != nil || reparse {
		return newWireError("output_failed", "output directory contains a reparse point").withField("params.output.directory")
	}
	final := filepath.Join(dir, filename)
	rel, err := filepath.Rel(dir, final)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return newWireError("output_failed", "output path escapes the request directory")
	}
	if _, err := os.Lstat(final); err == nil {
		return newWireError("output_failed", "final output already exists")
	} else if !os.IsNotExist(err) {
		return newWireError("output_failed", "could not inspect final output")
	}
	return nil
}

func createPart(final string) (string, error) {
	for i := 0; i < 8; i++ {
		var b [8]byte
		if _, err := rand.Read(b[:]); err != nil {
			return "", err
		}
		part := final + ".part-" + hex.EncodeToString(b[:])
		if _, err := os.Lstat(part); os.IsNotExist(err) {
			return part, nil
		}
	}
	return "", fmt.Errorf("could not allocate a unique temporary file")
}

func canceled(ch <-chan struct{}) bool {
	if ch == nil {
		return false
	}
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func (s *server) mark(name string) {
	if s.markerDir == "" {
		return
	}
	_ = os.WriteFile(filepath.Join(s.markerDir, name), []byte("1"), 0o600)
}

func (s *server) respondResult(id string, result any) {
	line, err := successFrame(id, result)
	if err != nil {
		s.respondError(id, newWireError("internal_error", "could not encode fake response"))
		return
	}
	_ = s.writeRaw(line)
	s.finishRequest(id)
}

func (s *server) respondError(id string, e *wireError) {
	if e == nil {
		e = newWireError("internal_error", "fake helper failure")
	}
	s.writeError(id, e)
	s.finishRequest(id)
}

func (s *server) writeError(id string, e *wireError) {
	line, err := errorFrame(id, e)
	if err != nil {
		return
	}
	_ = s.writeRaw(line)
}

func (s *server) writeRaw(line []byte) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	if _, err := s.out.Write(line); err != nil {
		return err
	}
	if err := s.out.WriteByte('\n'); err != nil {
		return err
	}
	return s.out.Flush()
}

func (s *server) markInflight(id string) {
	s.mu.Lock()
	s.inflight[id] = struct{}{}
	s.mu.Unlock()
}

func (s *server) duplicate(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.inflight[id]
	return ok
}

func (s *server) finishRequest(id string) {
	if id == "" {
		return
	}
	s.mu.Lock()
	delete(s.inflight, id)
	s.mu.Unlock()
}

func (s *server) handshaken() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.handshake
}

func (s *server) finishWork() {
	s.mu.Lock()
	s.busy = false
	s.activeID = ""
	s.activeKind = ""
	s.cancel = nil
	s.mu.Unlock()
}

func (s *server) stop() {
	s.stopOnce.Do(func() { close(s.stopCh) })
}

func (s *server) waitWorker(timeout time.Duration) {
	s.mu.Lock()
	if s.shutdownWaited {
		s.mu.Unlock()
		return
	}
	s.shutdownWaited = true
	s.mu.Unlock()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-s.workerDone:
	case <-timer.C:
	}
}
