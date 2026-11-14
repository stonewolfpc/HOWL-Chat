package helper

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"howl-chat/internal/audio/types"
)

// ClientOptions sets the internal defaults from the frozen numeric policy.
type ClientOptions struct {
	ClientVersion    string
	HandshakeTimeout time.Duration
	OperationTimeout time.Duration
	CancelGrace      time.Duration
	ShutdownTimeout  time.Duration
	MaxPending       int
}

// SynthesisRequest is the Stage 1 file-path request. It deliberately excludes
// pipeline, character, and cache policy; those belong to later stages.
type SynthesisRequest struct {
	ModelIdentity   string
	Text            string
	VoiceID         string
	Language        string
	Speed           float64
	OutputDirectory string
	Filename        string
}

// LoadRequest describes one provider load without provider-specific fields.
type LoadRequest struct {
	ProviderID    string
	Engine        string
	ModelIdentity string
	Assets        LoadAssets
}

type LoadAssets struct {
	Root   string
	Model  string
	Voices string
	Config string
	Data   string
}

type pendingOutcome struct {
	resp    *responseEnvelope
	failure *Failure
}

type pendingRequest struct {
	op       string
	ch       chan pendingOutcome
	terminal bool

	deliverMu sync.Mutex
	delivered bool
}

// deliver publishes the first outcome without blocking the protocol reader.
// The channel is capacity-one; a consumer may have retired the request after
// already taking a terminal, leaving the buffer full. That must never stop the
// fatal path from reaching process termination.
func (r *pendingRequest) deliver(out pendingOutcome) {
	if r == nil {
		return
	}
	r.deliverMu.Lock()
	defer r.deliverMu.Unlock()
	if r.delivered {
		return
	}
	r.delivered = true
	select {
	case r.ch <- out:
	default:
		// A terminal is already queued. Preserve it; the fatal session failure
		// is still recorded in Client.readFailure and can terminate the process.
	}
}

type requestAcceptance struct {
	mu     sync.Mutex
	state  string
	done   chan struct{}
	closed bool
}

func newRequestAcceptance() *requestAcceptance {
	return &requestAcceptance{state: "active", done: make(chan struct{})}
}

func (r *requestAcceptance) set(state string) {
	r.mu.Lock()
	if !r.closed {
		r.state = state
	}
	r.mu.Unlock()
}

func (r *requestAcceptance) snapshot() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state
}

func (r *requestAcceptance) close() {
	r.mu.Lock()
	if !r.closed {
		r.closed = true
		close(r.done)
	}
	r.mu.Unlock()
}

// Client is the parent-side Protocol v1 client. The process reader is separate
// from operation execution, so health and cancellation remain responsive while
// a synthesis request is outstanding.
type Client struct {
	sup  *Supervisor
	opts ClientOptions

	startMu    sync.Mutex
	transport  Process
	started    bool
	session    string
	readerDone chan struct{}

	mu          sync.Mutex
	next        uint64
	pending     map[string]*pendingRequest
	retired     map[string]struct{}
	readFailure *Failure
	closing     bool

	workMu sync.Mutex

	hello      helloResult
	providerID string
	engine     string

	// validationHook is an internal test seam for the narrow window between a
	// received terminal and the A1 acceptance re-check. Production leaves it nil.
	validationHook func()
}

// NewClient returns an unstarted client. It does not spawn the helper.
func NewClient(sup *Supervisor, opts ClientOptions) (*Client, error) {
	if sup == nil {
		return nil, newFailure("helper_start_failed", "supervisor is nil")
	}
	if opts.ClientVersion == "" {
		opts.ClientVersion = "0.1.0"
	}
	if opts.HandshakeTimeout <= 0 {
		opts.HandshakeTimeout = defaultHandshakeTimeout
	}
	if opts.OperationTimeout <= 0 {
		opts.OperationTimeout = defaultOperationTimeout
	}
	if opts.CancelGrace <= 0 {
		opts.CancelGrace = defaultCancelGrace
	}
	if opts.ShutdownTimeout <= 0 {
		opts.ShutdownTimeout = defaultShutdownTimeout
	}
	if opts.MaxPending <= 0 {
		opts.MaxPending = 32
	}
	session, err := newSessionID()
	if err != nil {
		return nil, newFailure("helper_start_failed", "could not create protocol session")
	}
	return &Client{
		sup:     sup,
		opts:    opts,
		session: session,
		pending: make(map[string]*pendingRequest),
		retired: make(map[string]struct{}),
	}, nil
}

// Start performs lazy discovery, startup, and the hello handshake. It never
// loads a provider or a model.
func (c *Client) Start(ctx context.Context) error {
	if c == nil {
		return newFailure("helper_start_failed", "client is nil").AudioError()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	c.startMu.Lock()
	defer c.startMu.Unlock()
	c.mu.Lock()
	if c.started && c.transport != nil {
		if c.sup.ProcessAlive() {
			c.mu.Unlock()
			return nil
		}
		// A crashed helper is retried only by an explicit request, never by an
		// autonomous restart loop. Rotate the session and start a new observer.
		c.started = false
		c.transport = nil
		c.readFailure = nil
		c.retired = make(map[string]struct{})
	}
	if c.closing {
		c.mu.Unlock()
		return newFailure("shutdown", "TTS helper client is closing")
	}
	if !c.started && c.transport == nil {
		c.readFailure = nil
		c.retired = make(map[string]struct{})
		if session, err := newSessionID(); err == nil {
			c.session = session
		}
	}
	c.mu.Unlock()

	p, err := c.sup.Start(ctx)
	if err != nil {
		return err
	}
	readerDone := make(chan struct{})
	c.mu.Lock()
	c.transport = p
	c.started = true
	c.readerDone = readerDone
	c.mu.Unlock()
	go func() {
		defer close(readerDone)
		c.readLoop(p)
	}()

	h, failure := c.helloHandshake(ctx)
	if failure != nil {
		// A failed handshake must not leave the old reader alive or make a
		// later retry share its stream. Terminate, observe exit, then clear the
		// transport before returning the failure. Preserve the typed failure
		// state after termination so an incompatible helper is not reported as
		// merely absent.
		_ = c.sup.Terminate(context.Background(), c.opts.ShutdownTimeout)
		waitForReader(readerDone, c.opts.ShutdownTimeout)
		if failure.Code == "protocol_incompatible" || failure.Code == "protocol_minor_unsupported" {
			c.sup.SetState(types.TTSStateIncompatible, "", failure)
		} else {
			c.sup.SetState(types.TTSStateUnavailable, "", failure)
		}
		c.mu.Lock()
		c.started = false
		c.transport = nil
		c.mu.Unlock()
		return failure
	}
	c.mu.Lock()
	c.hello = h
	var selected *providerDescriptor
	if len(h.Providers) > 0 {
		selected = &h.Providers[0]
		c.providerID = selected.ProviderID
		c.engine = selected.Engine
	}
	c.mu.Unlock()
	c.sup.MarkHandshake(h, selected)
	return nil
}

func (c *Client) helloHandshake(ctx context.Context) (helloResult, *Failure) {
	var h helloResult
	params := helloParams{
		Client:        "howl-chat",
		ClientVersion: c.opts.ClientVersion,
		ProtocolMin:   ProtocolVersion,
		ProtocolMax:   ProtocolVersion,
	}
	resp, failure := c.callOn(ctx, "hello", params, c.opts.HandshakeTimeout)
	if failure != nil {
		return h, failure
	}
	if err := decodeResult(resp.Result, &h); err != nil {
		return h, asFailure(err)
	}
	if err := validateHello(h); err != nil {
		return h, asFailure(err)
	}
	return h, nil
}

func (c *Client) readLoop(p Process) {
	for {
		line, err := p.Receive()
		if err != nil {
			f := asFailure(err)
			if f.Code == "internal_error" {
				f = newFailure("helper_crashed", "TTS helper closed its protocol stream")
			}
			c.failAll(f)
			if f.Code == "invalid_frame" || f.Code == "frame_too_large" || f.Code == "protocol_violation" {
				_ = c.sup.Kill(context.Background())
			}
			return
		}
		resp, err := decodeResponseFrame(line)
		if err != nil {
			c.failAll(asFailure(err))
			_ = c.sup.Kill(context.Background())
			return
		}
		if !c.dispatch(resp) {
			return
		}
	}
}

func (c *Client) dispatch(resp responseEnvelope) bool {
	if resp.ID == "" {
		c.protocolViolation(newFailure("protocol_violation", "frame-level response is not a request terminal"))
		return false
	}
	c.mu.Lock()
	req, ok := c.pending[resp.ID]
	if ok {
		if req.terminal {
			c.mu.Unlock()
			c.protocolViolation(newFailure("protocol_violation", "duplicate terminal response"))
			return false
		}
		req.terminal = true
		c.mu.Unlock()
		req.deliver(pendingOutcome{resp: &resp})
		return true
	}
	if _, retired := c.retired[resp.ID]; retired {
		// A late terminal is evidence only: never match it to another request.
		c.mu.Unlock()
		return true
	}
	c.mu.Unlock()
	c.protocolViolation(newFailure("protocol_violation", "response has an unknown request ID"))
	return false
}

func waitForReader(done <-chan struct{}, timeout time.Duration) {
	if done == nil {
		return
	}
	if timeout <= 0 {
		timeout = defaultShutdownTimeout
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-done:
	case <-timer.C:
	}
}

func (c *Client) protocolViolation(f *Failure) {
	c.failAll(f)
	_ = c.sup.Kill(context.Background())
}

func (c *Client) failAll(f *Failure) {
	if f == nil {
		f = newFailure("helper_crashed", "TTS helper failed")
	}
	c.mu.Lock()
	if c.readFailure == nil {
		c.readFailure = f
	}
	pending := make(map[string]*pendingRequest, len(c.pending))
	for id, req := range c.pending {
		pending[id] = req
	}
	c.pending = make(map[string]*pendingRequest)
	c.mu.Unlock()
	for _, req := range pending {
		req.deliver(pendingOutcome{failure: f})
	}
}

func (c *Client) call(ctx context.Context, op string, params any, timeout time.Duration) (*responseEnvelope, *Failure) {
	if err := c.Start(ctx); err != nil {
		if f, ok := AsFailure(err); ok {
			return nil, f
		}
		return nil, newFailure("helper_start_failed", err.Error())
	}
	return c.callOn(ctx, op, params, timeout)
}

func (c *Client) callOn(ctx context.Context, op string, params any, timeout time.Duration) (*responseEnvelope, *Failure) {
	if ctx == nil {
		ctx = context.Background()
	}
	id, req, failure := c.begin(op)
	if failure != nil {
		return nil, failure
	}
	cleanup := func() { c.retire(id) }
	defer cleanup()

	frame, err := encodeRequestFrame(id, op, params)
	if err != nil {
		if f, ok := AsFailure(err); ok {
			return nil, f
		}
		return nil, newFailure("invalid_request", err.Error())
	}
	transport := c.currentTransport()
	if transport == nil {
		return nil, newFailure("helper_start_failed", "TTS helper is not started")
	}
	if err := transport.Send(ctx, frame); err != nil {
		return nil, newFailure("helper_crashed", "could not write to TTS helper")
	}

	if timeout <= 0 {
		timeout = c.opts.OperationTimeout
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case out := <-req.ch:
		if out.failure != nil {
			return nil, out.failure
		}
		if out.resp.Error != nil {
			return nil, failureFromWire(out.resp.Error)
		}
		return out.resp, nil
	case <-timer.C:
		return nil, newFailure("timeout", "TTS helper request timed out").withRetryable(true)
	case <-ctx.Done():
		return nil, newFailure("cancelled", "TTS helper request was cancelled")
	}
}

func (c *Client) begin(op string) (string, *pendingRequest, *Failure) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.readFailure != nil {
		return "", nil, c.readFailure
	}
	if c.transport == nil {
		return "", nil, newFailure("helper_start_failed", "TTS helper is not started")
	}
	if len(c.pending) >= c.opts.MaxPending {
		return "", nil, newFailure("busy", "too many in-flight TTS helper requests").withRetryable(true)
	}
	c.next++
	id := newRequestID(c.session, c.next)
	req := &pendingRequest{op: op, ch: make(chan pendingOutcome, 1)}
	c.pending[id] = req
	return id, req, nil
}

func (c *Client) retire(id string) {
	c.mu.Lock()
	delete(c.pending, id)
	c.retired[id] = struct{}{}
	c.mu.Unlock()
}

func failureFromWire(e *wireError) *Failure {
	if e == nil {
		return newFailure("protocol_violation", "missing helper error")
	}
	return &Failure{
		Code:      e.Code,
		Message:   e.Message,
		Retryable: e.Retryable,
		Field:     e.Field,
		Detail:    cloneDetail(e.Detail),
	}
}

func asFailure(err error) *Failure {
	if err == nil {
		return nil
	}
	if f, ok := AsFailure(err); ok {
		return f
	}
	var audioErr *types.AudioError
	if errors.As(err, &audioErr) {
		if f, ok := audioErr.Cause.(*Failure); ok {
			return f
		}
	}
	return newFailure("internal_error", err.Error())
}

// Health is a read-only snapshot and is permitted while synthesis is active.
func (c *Client) Health(ctx context.Context) (healthResult, error) {
	var h healthResult
	resp, f := c.call(ctx, "health", healthParams{}, c.opts.OperationTimeout)
	if f != nil {
		return h, f.AudioError()
	}
	if err := decodeResult(resp.Result, &h); err != nil {
		return h, asFailure(err).AudioError()
	}
	if err := validateHealth(h); err != nil {
		return h, asFailure(err).AudioError()
	}
	return h, nil
}

// Catalog returns provider metadata without loading the provider.
func (c *Client) Catalog(ctx context.Context, req types.TTSProviderCatalogRequest) (types.TTSProviderCatalog, error) {
	c.workMu.Lock()
	defer c.workMu.Unlock()
	if err := c.Start(ctx); err != nil {
		if f, ok := AsFailure(err); ok {
			return types.TTSProviderCatalog{}, f.AudioError()
		}
		return types.TTSProviderCatalog{}, err
	}

	if failure := c.requireServiceCapability(CapabilityCatalog); failure != nil {
		return types.TTSProviderCatalog{}, failure.AudioError()
	}
	provider, engine, failure := c.selectedProvider(req.Engine)
	if failure != nil {
		return types.TTSProviderCatalog{}, failure.AudioError()
	}
	params := catalogParams{ProviderID: provider, Engine: engine, IncludeUnavailable: req.IncludeUnavailable}
	resp, f := c.callOn(ctx, "catalog", params, c.opts.OperationTimeout)
	if f != nil {
		return types.TTSProviderCatalog{}, f.AudioError()
	}
	var result catalogResult
	if err := decodeResult(resp.Result, &result); err != nil {
		return types.TTSProviderCatalog{}, asFailure(err).AudioError()
	}
	if err := validateCatalog(result, params); err != nil {
		return types.TTSProviderCatalog{}, asFailure(err).AudioError()
	}
	c.sup.SetSelected(result.ProviderID, result.Engine, "", nil)
	return mapCatalog(result), nil
}

// Load explicitly loads a provider. It does not synthesize.
func (c *Client) Load(ctx context.Context, req LoadRequest) error {
	c.workMu.Lock()
	defer c.workMu.Unlock()
	if err := c.Start(ctx); err != nil {
		if f, ok := AsFailure(err); ok {
			return f.AudioError()
		}
		return err
	}
	if failure := c.requireServiceCapability(CapabilityLoad); failure != nil {
		return failure.AudioError()
	}
	if err := validateAssets(loadAssets{
		Root: req.Assets.Root, Model: req.Assets.Model, Voices: req.Assets.Voices,
		Config: req.Assets.Config, Data: req.Assets.Data,
	}); err != nil {
		return asFailure(err).AudioError()
	}
	if req.ProviderID == "" || !providerPattern.MatchString(req.ProviderID) ||
		req.Engine == "" || !enginePattern.MatchString(req.Engine) ||
		req.ModelIdentity == "" || len(req.ModelIdentity) > 512 {
		return newFailure("invalid_request", "load identity is invalid").AudioError()
	}
	params := loadParams{
		ProviderID: req.ProviderID, Engine: req.Engine, ModelIdentity: req.ModelIdentity,
		Assets: loadAssets{Root: req.Assets.Root, Model: req.Assets.Model, Voices: req.Assets.Voices, Config: req.Assets.Config, Data: req.Assets.Data},
	}
	c.sup.SetState(types.TTSStateLoading, "", nil)
	resp, f := c.callOn(ctx, "load", params, c.opts.OperationTimeout)
	if f != nil {
		c.sup.SetState(types.TTSStateUnavailable, "", f)
		return f.AudioError()
	}
	var result loadResult
	if err := decodeResult(resp.Result, &result); err != nil {
		c.sup.SetState(types.TTSStateUnavailable, "", asFailure(err))
		return asFailure(err).AudioError()
	}
	if err := validateLoad(result, params); err != nil {
		c.sup.SetState(types.TTSStateUnavailable, "", asFailure(err))
		return asFailure(err).AudioError()
	}
	c.sup.SetSelected(result.ProviderID, result.Engine, result.ProviderVersion, result.Runtime)
	c.sup.SetState(types.TTSStateReady, "", nil)
	return nil
}

// Unload releases a loaded provider and is idempotent at the helper.
func (c *Client) Unload(ctx context.Context) error {
	c.workMu.Lock()
	defer c.workMu.Unlock()
	if failure := c.requireServiceCapability(CapabilityUnload); failure != nil {
		return failure.AudioError()
	}
	resp, f := c.call(ctx, "unload", struct{}{}, c.opts.OperationTimeout)
	if f != nil {
		return f.AudioError()
	}
	var result unloadResult
	if err := decodeResult(resp.Result, &result); err != nil {
		return asFailure(err).AudioError()
	}
	if !result.Unloaded {
		return newFailure("protocol_violation", "unload result is false").AudioError()
	}
	c.sup.SetState(types.TTSStateUnloaded, "", nil)
	return nil
}

// Cancel requests cooperative cancellation of a target synthesis. Its own
// request ID is independent of the target.
func (c *Client) Cancel(ctx context.Context, target string) error {
	if !requestIDPattern.MatchString(target) {
		return newFailure("invalid_request", "cancel target is invalid").withField("params.request_id").AudioError()
	}
	if failure := c.requireServiceCapability(CapabilityCancel); failure != nil {
		return failure.AudioError()
	}
	resp, f := c.call(ctx, "cancel", cancelParams{RequestID: target}, c.opts.CancelGrace)
	if f != nil {
		return f.AudioError()
	}
	var result cancelResult
	if err := decodeResult(resp.Result, &result); err != nil {
		return asFailure(err).AudioError()
	}
	if result.RequestID != target || !result.CancelRequested {
		return newFailure("protocol_violation", "cancel acknowledgement is invalid").AudioError()
	}
	return nil
}

// Shutdown sends shutdown, waits for the acknowledgement, and then terminates
// the process with the frozen grace. It does not kill an already absent helper.
func (c *Client) Shutdown(ctx context.Context) error {
	if c == nil {
		return nil
	}
	if c.currentTransport() == nil {
		return nil
	}
	resp, f := c.call(ctx, "shutdown", struct{}{}, c.opts.ShutdownTimeout)
	if f != nil && f.Code != "helper_crashed" && f.Code != "timeout" {
		_ = c.sup.Terminate(ctx, c.opts.ShutdownTimeout)
		return f.AudioError()
	}
	if f == nil {
		var result shutdownResult
		if err := decodeResult(resp.Result, &result); err != nil || !result.Accepted {
			if err != nil {
				_ = c.sup.Terminate(ctx, c.opts.ShutdownTimeout)
				return asFailure(err).AudioError()
			}
			_ = c.sup.Terminate(ctx, c.opts.ShutdownTimeout)
			return newFailure("protocol_violation", "shutdown was not accepted").AudioError()
		}
	}
	err := c.sup.Terminate(ctx, c.opts.ShutdownTimeout)
	c.mu.Lock()
	c.closing = true
	c.started = false
	c.transport = nil
	c.mu.Unlock()
	return err
}

// Synthesize executes one file-path synthesis and applies the A1 acceptance
// barrier before returning a usable result.
func (c *Client) Synthesize(ctx context.Context, req SynthesisRequest) (*types.SynthesisResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if textErr := validateText(req.Text); textErr != nil {
		return nil, asFailure(textErr).AudioError()
	}
	if outputErr := validateOutput(req.OutputDirectory, req.Filename); outputErr != nil {
		return nil, asFailure(outputErr).AudioError()
	}
	if !finite(req.Speed) || req.Speed < 0.5 || req.Speed > 2.0 {
		return nil, newFailure("invalid_request", "speed must be finite and between 0.5 and 2.0").withField("params.speed").AudioError()
	}
	if req.ModelIdentity == "" || len(req.ModelIdentity) > 512 || req.VoiceID == "" || len(req.VoiceID) > 128 {
		return nil, newFailure("invalid_request", "synthesis identity or voice is invalid").AudioError()
	}
	c.workMu.Lock()
	defer c.workMu.Unlock()
	if err := c.Start(ctx); err != nil {
		if f, ok := AsFailure(err); ok {
			return nil, f.AudioError()
		}
		return nil, err
	}
	if failure := c.requireServiceCapability(CapabilitySynthesize); failure != nil {
		return nil, failure.AudioError()
	}
	providerID, _, failure := c.selectedProvider("")
	if failure != nil {
		return nil, failure.AudioError()
	}
	if failure := c.requireProviderCapability(providerID, CapabilitySynthesize); failure != nil {
		return nil, failure.AudioError()
	}
	if req.Speed != 1.0 {
		if failure := c.requireProviderCapability(providerID, CapabilitySpeed); failure != nil {
			return nil, failure.AudioError()
		}
	}
	if req.Language != "" {
		if failure := c.requireProviderCapability(providerID, CapabilityLanguage); failure != nil {
			return nil, failure.AudioError()
		}
	}

	params := synthesisParams{
		ModelIdentity: req.ModelIdentity,
		Text:          req.Text,
		VoiceID:       req.VoiceID,
		Language:      req.Language,
		Speed:         req.Speed,
		Output:        outputSpec{Directory: filepath.Clean(req.OutputDirectory), Filename: req.Filename},
	}
	id, pending, failure := c.begin("synthesize")
	if failure != nil {
		return nil, failure.AudioError()
	}
	state := newRequestAcceptance()
	c.sup.SetState(types.TTSStateBusy, id, nil)
	go c.watchCancellation(ctx, id, state)
	defer func() {
		c.retire(id)
		state.close()
	}()

	frame, err := encodeRequestFrame(id, "synthesize", params)
	if err != nil {
		c.retire(id)
		state.close()
		c.sup.SetState(types.TTSStateReady, "", nil)
		return nil, asFailure(err).AudioError()
	}
	transport := c.currentTransport()
	if transport == nil {
		return nil, newFailure("helper_start_failed", "TTS helper is not started").AudioError()
	}
	if err := transport.Send(ctx, frame); err != nil {
		c.sup.SetState(types.TTSStateCrashed, "", newFailure("helper_crashed", "synthesis write failed"))
		c.cleanupRejected(req.OutputDirectory, req.Filename, true)
		return nil, newFailure("helper_crashed", "could not write synthesis request").AudioError()
	}

	timeout := time.NewTimer(c.opts.OperationTimeout)
	defer timeout.Stop()
	var outcome pendingOutcome
	select {
	case outcome = <-pending.ch:
	case <-timeout.C:
		state.set("timeout")
		c.retire(id)
		state.close()
		c.cleanupRejected(req.OutputDirectory, req.Filename, true)
		c.sup.SetState(types.TTSStateReady, "", newFailure("timeout", "synthesis timed out").withRetryable(true))
		return nil, newFailure("timeout", "synthesis timed out").withRetryable(true).AudioError()
	case <-ctx.Done():
		state.set("cancel_requested")
		// The watcher sends cancel; wait for its grace here before killing.
		select {
		case outcome = <-pending.ch:
		case <-time.After(c.opts.CancelGrace):
			state.set("cancelled")
			c.retire(id)
			state.close()
			c.cleanupRejected(req.OutputDirectory, req.Filename, true)
			c.sup.SetState(types.TTSStateReady, "", newFailure("cancelled", "synthesis was cancelled"))
			return nil, newFailure("cancelled", "synthesis was cancelled").AudioError()
		}
	}

	c.retire(id)
	defer state.close()
	if outcome.failure != nil {
		c.cleanupRejected(req.OutputDirectory, req.Filename, true)
		if state.snapshot() == "cancel_requested" {
			c.sup.SetState(types.TTSStateReady, "", outcome.failure)
			return nil, newFailure("cancelled", "synthesis was cancelled").AudioError()
		}
		c.sup.SetState(types.TTSStateDegraded, "", outcome.failure)
		return nil, outcome.failure.AudioError()
	}
	if outcome.resp.Error != nil {
		f := failureFromWire(outcome.resp.Error)
		c.cleanupRejected(req.OutputDirectory, req.Filename, false)
		if state.snapshot() == "cancel_requested" {
			c.sup.SetState(types.TTSStateReady, "", f)
			return nil, newFailure("cancelled", "synthesis was cancelled").AudioError()
		}
		c.sup.SetState(types.TTSStateDegraded, "", f)
		return nil, f.AudioError()
	}

	var wire synthesisResult
	if err := decodeResult(outcome.resp.Result, &wire); err != nil {
		f := asFailure(err)
		c.cleanupRejected(req.OutputDirectory, req.Filename, false)
		c.sup.SetState(types.TTSStateDegraded, "", f)
		return nil, f.AudioError()
	}
	if err := validateSynthesisResult(wire, params); err != nil {
		f := asFailure(err)
		c.cleanupRejected(req.OutputDirectory, req.Filename, false)
		c.sup.SetState(types.TTSStateDegraded, "", f)
		return nil, f.AudioError()
	}
	state.set("validating")
	metadata, wavErr := validateWAVFile(wire, req.OutputDirectory, req.Filename)
	if wavErr != nil {
		f := asFailure(wavErr)
		c.cleanupRejected(req.OutputDirectory, req.Filename, false)
		c.sup.SetState(types.TTSStateDegraded, "", f)
		return nil, f.AudioError()
	}
	if c.validationHook != nil {
		c.validationHook()
	}

	state.mu.Lock()
	if state.closed || state.state == "cancel_requested" || state.state == "cancelled" || state.state == "timeout" {
		state.mu.Unlock()
		c.cleanupRejected(req.OutputDirectory, req.Filename, false)
		c.sup.SetState(types.TTSStateReady, "", newFailure("cancelled", "synthesis was cancelled before acceptance"))
		return nil, newFailure("cancelled", "synthesis was cancelled before acceptance").AudioError()
	}
	state.state = "accepted"
	state.closed = true
	close(state.done)
	state.mu.Unlock()

	c.sup.SetState(types.TTSStateReady, "", nil)
	return &types.SynthesisResult{
		AudioPath: wire.OutputPath,
		Metadata:  *metadata,
		Duration:  wire.DurationSeconds,
	}, nil
}

// cleanupRejected enforces the contract ordering before deleting request-owned
// output. If the helper is still alive and output may exist, it terminates the
// process and waits for the single exit observer first. force is used when no
// terminal arrived, so the helper could still create output later.
func (c *Client) cleanupRejected(dir, filename string, force bool) {
	if c != nil && c.sup != nil && c.sup.ProcessAlive() && (force || requestOutputExists(dir, filename)) {
		_ = c.sup.Terminate(context.Background(), c.opts.ShutdownTimeout)
	}
	cleanupExact(dir, filename)
}

func (c *Client) watchCancellation(ctx context.Context, target string, state *requestAcceptance) {
	select {
	case <-ctx.Done():
		state.set("cancel_requested")
		go func() {
			_ = c.Cancel(context.Background(), target)
		}()
	case <-state.done:
	}
}

func (c *Client) requireServiceCapability(id string) *Failure {
	c.mu.Lock()
	caps := c.hello.Capabilities
	c.mu.Unlock()
	if !hasCapability(caps, id) {
		return newFailure("unsupported_capability", "TTS helper does not advertise "+id)
	}
	return nil
}

func (c *Client) requireProviderCapability(providerID, id string) *Failure {
	c.mu.Lock()
	caps := []capability(nil)
	for _, p := range c.hello.Providers {
		if p.ProviderID == providerID {
			caps = p.Capabilities
			break
		}
	}
	c.mu.Unlock()
	if !hasCapability(caps, id) {
		return newFailure("unsupported_capability", "TTS provider does not advertise "+id)
	}
	return nil
}

func (c *Client) selectedProvider(requestedEngine string) (providerID, engine string, failure *Failure) {
	c.mu.Lock()
	providerID, engine = c.providerID, c.engine
	hello := c.hello
	c.mu.Unlock()
	if requestedEngine != "" {
		for _, p := range hello.Providers {
			if p.Engine == requestedEngine {
				providerID, engine = p.ProviderID, p.Engine
				break
			}
		}
	}
	if engine == "" {
		return "", "", newFailure("provider_unknown", "no helper provider was advertised")
	}
	if providerID == "" || !providerPattern.MatchString(providerID) || !enginePattern.MatchString(engine) {
		return "", "", newFailure("provider_unknown", "provider identity is invalid")
	}
	if requestedEngine != "" && engine != requestedEngine {
		return "", "", newFailure("provider_unknown", "requested engine is not advertised")
	}
	return providerID, engine, nil
}

func (c *Client) currentTransport() Process {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.transport
}

func mapCatalog(c catalogResult) types.TTSProviderCatalog {
	out := types.TTSProviderCatalog{
		Engine: c.Engine, ProviderID: c.ProviderID, CatalogVersion: c.CatalogVersion,
		Available: c.Available, Detail: c.Detail, DefaultVoice: c.DefaultVoice,
		Languages: append([]string(nil), c.Languages...), SSMLSupported: c.SSMLSupported,
		StreamSupported: c.StreamingSupported,
		Speed:           &types.TTSSpeedRange{Supported: c.Speed.Supported, Min: c.Speed.Min, Max: c.Speed.Max},
	}
	for _, v := range c.Voices {
		out.Voices = append(out.Voices, types.TTSVoiceInfo{
			Voice:     types.Voice{ID: v.VoiceID, Name: v.Name, Language: v.Language, Gender: v.Gender, Description: v.Description},
			Available: v.Available, Detail: v.Detail,
		})
	}
	return out
}

// ProviderStatus may handshake/start but never loads a model.
func (c *Client) ProviderStatus(ctx context.Context) (types.TTSProviderStatus, error) {
	if err := c.Start(ctx); err != nil {
		return c.sup.Status(), err
	}
	return c.sup.Status(), nil
}

// ProviderCatalog is the additive inspector catalog path. It may start and
// handshake, but it does not load a model.
func (c *Client) ProviderCatalog(ctx context.Context, req types.TTSProviderCatalogRequest) (types.TTSProviderCatalog, error) {
	return c.Catalog(ctx, req)
}

func (c *Client) String() string {
	return fmt.Sprintf("helper client (provider=%s engine=%s)", c.providerID, c.engine)
}
