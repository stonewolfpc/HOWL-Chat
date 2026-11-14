package helper

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"howl-chat/internal/audio/types"
	"howl-chat/internal/proc"
)

const (
	defaultHandshakeTimeout = 5 * time.Second
	defaultOperationTimeout = 120 * time.Second
	defaultCancelGrace      = 2 * time.Second
	defaultShutdownTimeout  = 5 * time.Second
	defaultStderrLimit      = 8 * 1024
)

// Process is one started helper process and its parent-owned pipes. Wait may
// be called exactly once, and only by Supervisor.
type Process interface {
	Send(ctx context.Context, frame []byte) error
	Receive() ([]byte, error)
	Close() error
	Stderr() string
	Wait() error
	Kill() error
	PID() int
}

// ProcessFactory creates a helper process. It must not call Wait itself: a
// successful factory result transfers the single Wait ownership to Supervisor.
// A failed factory call must follow the existing proc.Start failure cleanup and
// must not create a second waiter.
type ProcessFactory func(ctx context.Context, spec LaunchSpec) (Process, error)

// LaunchSpec is the minimal process launch contract; helper arguments are
// deliberately not provider-specific.
type LaunchSpec struct {
	Path        string
	Dir         string
	Args        []string
	Env         []string
	StderrLimit int
}

// SupervisorOptions controls lazy discovery, startup, and lifecycle limits.
type SupervisorOptions struct {
	Discoverer  Discoverer
	Factory     ProcessFactory
	StderrLimit int
	Now         func() time.Time
}

// Supervisor owns the process lifecycle, one Wait observer, stderr diagnostics,
// and the process/provider status snapshot. It never loads a model.
type Supervisor struct {
	opts SupervisorOptions
	mu   sync.Mutex

	startMu  sync.Mutex
	path     string
	proc     Process
	procDone chan struct{}
	waitErr  error

	processState  types.TTSProcessState
	providerState types.TTSProviderState

	engine          string
	providerID      string
	providerVersion string
	helperVersion   string
	protocolVersion string
	capabilities    []types.TTSCapability
	runtime         map[string]string
	activeRequestID string
	lastFailure     *Failure

	lastStderr string
}

// NewSupervisor constructs a supervisor with the frozen provisional defaults.
func NewSupervisor(opts SupervisorOptions) *Supervisor {
	if opts.Factory == nil {
		opts.Factory = StartProcess
	}
	if opts.StderrLimit <= 0 {
		opts.StderrLimit = defaultStderrLimit
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	return &Supervisor{
		opts:          opts,
		processState:  types.TTSProcessAbsent,
		providerState: types.TTSStateAbsent,
	}
}

// Start lazily discovers and starts the helper. A missing helper is a TTS-only
// failure; it never affects the rest of the application.
func (s *Supervisor) Start(ctx context.Context) (Process, error) {
	if s == nil {
		return nil, newFailure("helper_start_failed", "supervisor is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	s.startMu.Lock()
	defer s.startMu.Unlock()

	s.mu.Lock()
	if s.proc != nil {
		p := s.proc
		s.mu.Unlock()
		return p, nil
	}
	s.mu.Unlock()

	discovered, failure := s.opts.Discoverer.Discover()
	if failure != nil {
		s.setUnavailable(failure)
		return nil, failure
	}

	s.mu.Lock()
	s.path = discovered.Path
	s.processState = types.TTSProcessStarting
	s.providerState = types.TTSStateStarting
	s.mu.Unlock()

	p, err := s.opts.Factory(ctx, LaunchSpec{
		Path:        discovered.Path,
		Dir:         filepath.Dir(discovered.Path),
		Env:         os.Environ(),
		StderrLimit: s.opts.StderrLimit,
	})
	if err != nil {
		f := newFailure("helper_start_failed", "TTS helper could not be started")
		f.Detail = map[string]string{"source": discovered.Source}
		s.setUnavailable(f)
		return nil, f
	}

	done := make(chan struct{})
	s.mu.Lock()
	s.proc = p
	s.procDone = done
	s.waitErr = nil
	s.processState = types.TTSProcessStarting
	s.providerState = types.TTSStateHandshaking
	s.mu.Unlock()

	go s.observeExit(p, done)
	return p, nil
}

func (s *Supervisor) observeExit(p Process, done chan struct{}) {
	err := p.Wait()
	stderr := p.Stderr()
	_ = p.Close()

	s.mu.Lock()
	if s.proc == p {
		s.proc = nil
		s.waitErr = err
		s.lastStderr = stderr
		if s.providerState == types.TTSStateStopping {
			if err == nil {
				s.processState = types.TTSProcessAbsent
			} else {
				s.processState = types.TTSProcessExited
			}
			s.providerState = types.TTSStateAbsent
		} else {
			s.processState = types.TTSProcessExited
			s.providerState = types.TTSStateCrashed
			f := newFailure("helper_crashed", "TTS helper exited unexpectedly")
			if err != nil {
				f.Message = fmt.Sprintf("TTS helper exited: %v", err)
			}
			s.lastFailure = f
		}
	}
	close(done)
	s.mu.Unlock()
}

func (s *Supervisor) setUnavailable(f *Failure) {
	s.mu.Lock()
	s.processState = types.TTSProcessAbsent
	s.providerState = types.TTSStateUnavailable
	s.lastFailure = f
	s.mu.Unlock()
}

// MarkHandshake records a validated hello result without loading a provider.
func (s *Supervisor) MarkHandshake(h helloResult, selected *providerDescriptor) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.helperVersion = h.HelperVersion
	s.protocolVersion = ProtocolVersion.String()
	s.capabilities = capabilityList(h.Capabilities)
	s.processState = types.TTSProcessAlive
	s.providerState = types.TTSStateUnloaded
	if selected != nil {
		s.engine = selected.Engine
		s.providerID = selected.ProviderID
		s.providerVersion = selected.ProviderVersion
		s.runtime = cloneDetail(selected.Runtime)
	}
}

// SetSelected records the provider selected by catalog/load.
func (s *Supervisor) SetSelected(providerID, engine, providerVersion string, runtime map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.providerID = providerID
	s.engine = engine
	s.providerVersion = providerVersion
	s.runtime = cloneDetail(runtime)
}

// SetState updates provider state and the active request diagnostic.
func (s *Supervisor) SetState(state types.TTSProviderState, activeRequestID string, failure *Failure) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.providerState = state
	s.activeRequestID = activeRequestID
	s.lastFailure = failure
	if s.proc != nil && state != types.TTSStateStopping {
		s.processState = types.TTSProcessAlive
	}
}

// Terminate closes stdin, waits for the single observer, and kills the tree if
// the grace expires. It never deletes output itself.
func (s *Supervisor) Terminate(ctx context.Context, grace time.Duration) error {
	if s == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	s.mu.Lock()
	p := s.proc
	done := s.procDone
	if p != nil {
		s.processState = types.TTSProcessStopping
		s.providerState = types.TTSStateStopping
	}
	s.mu.Unlock()
	if p == nil {
		return nil
	}
	_ = p.Close()

	if grace <= 0 {
		grace = defaultShutdownTimeout
	}
	timer := time.NewTimer(grace)
	defer timer.Stop()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		_ = p.Kill()
		return ctx.Err()
	case <-timer.C:
		_ = p.Kill()
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(grace):
			return newFailure("timeout", "TTS helper did not terminate")
		}
	}
}

// Kill terminates promptly and waits for the single exit observer.
func (s *Supervisor) Kill(ctx context.Context) error {
	return s.Terminate(ctx, defaultShutdownTimeout)
}

// Stderr returns the bounded diagnostics ring from the current/last process.
func (s *Supervisor) Stderr() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.proc != nil {
		return s.proc.Stderr()
	}
	return s.lastStderr
}

// Status returns process/provider state without starting or loading anything.
func (s *Supervisor) Status() types.TTSProviderStatus {
	if s == nil {
		return types.TTSProviderStatus{ProcessState: types.TTSProcessAbsent, State: types.TTSStateAbsent}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	st := types.TTSProviderStatus{
		ProcessState:    s.processState,
		State:           s.providerState,
		Engine:          s.engine,
		ProviderID:      s.providerID,
		ProviderVersion: s.providerVersion,
		HelperVersion:   s.helperVersion,
		ProtocolVersion: s.protocolVersion,
		Capabilities:    append([]types.TTSCapability(nil), s.capabilities...),
		Runtime:         cloneDetail(s.runtime),
		ActiveRequestID: s.activeRequestID,
	}
	st.ProcessAlive = processAlive(st.ProcessState, st.State)
	st.Available = providerAvailable(st.State)
	st.Ready = providerReady(st.State)
	if s.lastFailure != nil {
		f := s.lastFailure.status()
		st.LastError = &f
	}
	return st
}

// ProcessAlive reports whether the supervisor currently owns a live process.
func (s *Supervisor) ProcessAlive() bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.proc != nil && s.processState != types.TTSProcessExited
}

// Done is closed when the current process exits; it is nil before a start.
func (s *Supervisor) Done() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.procDone
}

func processAlive(process types.TTSProcessState, provider types.TTSProviderState) bool {
	if process == types.TTSProcessAlive || process == types.TTSProcessStarting || process == types.TTSProcessStopping {
		return provider != types.TTSStateCrashed && provider != types.TTSStateUnavailable
	}
	return false
}

func providerAvailable(state types.TTSProviderState) bool {
	switch state {
	case types.TTSStateUnloaded, types.TTSStateLoading, types.TTSStateReady, types.TTSStateBusy, types.TTSStateDegraded:
		return true
	default:
		return false
	}
}

func providerReady(state types.TTSProviderState) bool {
	switch state {
	case types.TTSStateReady, types.TTSStateBusy, types.TTSStateDegraded:
		return true
	default:
		return false
	}
}

func capabilityList(in []capability) []types.TTSCapability {
	if len(in) == 0 {
		return nil
	}
	out := make([]types.TTSCapability, 0, len(in))
	for _, c := range in {
		out = append(out, types.TTSCapability{ID: c.ID, Version: c.Version})
	}
	return out
}

// StartProcess is the production process factory. It creates pipes before
// proc.Start, which establishes job ownership before the child executes, and
// transfers the only Wait call to Supervisor.
func StartProcess(_ context.Context, spec LaunchSpec) (Process, error) {
	if spec.Path == "" {
		return nil, fmt.Errorf("helper path is empty")
	}
	cmd := exec.Command(spec.Path, spec.Args...)
	if spec.Dir != "" {
		cmd.Dir = spec.Dir
	}
	if len(spec.Env) > 0 {
		cmd.Env = spec.Env
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = stdin.Close()
		_ = stdout.Close()
		return nil, err
	}
	proc.HideWindow(cmd)
	if err := proc.Start(cmd); err != nil {
		_ = stdin.Close()
		_ = stdout.Close()
		_ = stderr.Close()
		// proc.Start's failure path already waited for the child. Do not call
		// cmd.Wait here and do not create a second observer.
		return nil, err
	}

	limit := spec.StderrLimit
	if limit <= 0 {
		limit = defaultStderrLimit
	}
	p := &osProcess{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReaderSize(stdout, 64*1024),
		stderr: stderr,
		pid:    cmd.Process.Pid,
		ring:   newRing(limit),
	}
	go p.drainStderr()
	return p, nil
}

type osProcess struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	stderr io.ReadCloser
	pid    int
	ring   *ringBuffer

	writeMu   sync.Mutex
	closeOnce sync.Once
	waitOnce  sync.Once
	waitErr   error
	waitDone  chan struct{}
}

func (p *osProcess) Send(_ context.Context, frame []byte) error {
	if p == nil || p.stdin == nil {
		return io.ErrClosedPipe
	}
	p.writeMu.Lock()
	defer p.writeMu.Unlock()
	for len(frame) > 0 {
		n, err := p.stdin.Write(frame)
		if err != nil {
			return err
		}
		frame = frame[n:]
	}
	return nil
}

func (p *osProcess) Receive() ([]byte, error) {
	if p == nil || p.stdout == nil {
		return nil, io.EOF
	}
	return readFrame(p.stdout)
}

func (p *osProcess) Close() error {
	if p == nil {
		return nil
	}
	p.closeOnce.Do(func() {
		if p.stdin != nil {
			_ = p.stdin.Close()
		}
		if p.stderr != nil {
			_ = p.stderr.Close()
		}
	})
	return nil
}

func (p *osProcess) Stderr() string {
	if p == nil || p.ring == nil {
		return ""
	}
	return p.ring.String()
}

func (p *osProcess) Wait() error {
	if p == nil {
		return nil
	}
	p.waitOnce.Do(func() {
		p.waitErr = p.cmd.Wait()
	})
	return p.waitErr
}

func (p *osProcess) Kill() error {
	if p == nil {
		return nil
	}
	proc.KillTree(p.pid)
	if p.cmd != nil && p.cmd.Process != nil {
		return p.cmd.Process.Kill()
	}
	return nil
}

func (p *osProcess) PID() int { return p.pid }

func (p *osProcess) drainStderr() {
	if p == nil || p.stderr == nil || p.ring == nil {
		return
	}
	buf := make([]byte, 2048)
	for {
		n, err := p.stderr.Read(buf)
		if n > 0 {
			p.ring.Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

type ringBuffer struct {
	mu    sync.Mutex
	limit int
	data  []byte
}

func newRing(limit int) *ringBuffer {
	return &ringBuffer{limit: limit}
}

func (r *ringBuffer) Write(p []byte) {
	if r == nil || r.limit <= 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = append(r.data, p...)
	if len(r.data) > r.limit {
		r.data = append([]byte(nil), r.data[len(r.data)-r.limit:]...)
	}
}

func (r *ringBuffer) String() string {
	if r == nil {
		return ""
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return string(append([]byte(nil), r.data...))
}
