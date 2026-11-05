package audio

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"howl-chat/internal/audio/types"
)

const defaultCloseTimeout = 5 * time.Second

var (
	// ErrDrainTimeout reports that Close observed outstanding manager work
	// before its deadline. The manager remains closing and a later Close may
	// complete normally.
	ErrDrainTimeout = errors.New("audio manager drain timeout")

	// ErrManagerClosing is returned by operations that cannot begin after Close
	// has started.
	ErrManagerClosing = errors.New("audio manager is closing")

	// ErrManagerClosed is returned after the manager has finalized.
	ErrManagerClosed = errors.New("audio manager is closed")

	// ErrNoActivePipeline reports that no generation is available to acquire.
	ErrNoActivePipeline = errors.New("audio manager has no active pipeline")

	// ErrNoFactory reports a manager constructed without a pipeline factory.
	ErrNoFactory = errors.New("audio manager factory is nil")

	// ErrNoReadyCapability reports that a requested set produced no ready
	// capability and therefore cannot be activated.
	ErrNoReadyCapability = errors.New("no requested capability is ready")

	// ErrInvalidCandidate reports a candidate that claims a ready requested
	// capability without transferring a pipeline handle.
	ErrInvalidCandidate = errors.New("candidate reports a ready capability without a pipeline")

	// ErrCapabilityNotReady is wrapped when a requested capability reports
	// Ready false without its own error.
	ErrCapabilityNotReady = errors.New("capability not ready")
)

// ManagedPipeline is the manager-owned pipeline surface. It deliberately does
// not depend on *Pipeline so tests and later wiring can supply fakes.
type ManagedPipeline interface {
	ProcessAudioMessage(context.Context, types.AudioMessage) (*types.AudioResponse, error)
	GetStatus() types.PipelineStatus
	GetTTSVoices() []types.Voice
	Shutdown(context.Context) error
}

// Settings is the minimal T2.9.2 desired-capability set. T2.9.6 extends it for
// the full D5 settings contract.
type Settings struct {
	ASREnabled bool
	TTSEnabled bool
	// ASREngine and TTSEngine carry an explicit provider selection. "" and
	// "auto" mean automatic selection. They are part of the desired state, so
	// changing only the engine still rebuilds the generation.
	ASREngine string
	TTSEngine string
}

// CapabilityResult is the factory's per-capability outcome for a candidate.
type CapabilityResult struct {
	Ready bool
	Err   error
}

// Candidate transfers ownership of a factory-built pipeline and its
// per-capability outcomes to the manager.
type Candidate struct {
	Pipeline ManagedPipeline
	ASR      CapabilityResult
	TTS      CapabilityResult
}

// PipelineFactory builds one candidate synchronously. Ownership of every
// resource it creates transfers through candidate when the factory returns or
// panics.
type PipelineFactory func(ctx context.Context, desired Settings, candidate *Candidate) error

// Capability is the minimal T2.9.2 status view of one capability.
type Capability struct {
	Requested bool
	Ready     bool
	Error     string
}

// Status is the minimal T2.9.2 manager status. T2.9.6 extends it with the
// remaining D5 fields.
type Status struct {
	ASR      Capability
	TTS      Capability
	InFlight int
	Builds   int
	// DesiredGen advances when an Apply enters, before it waits on applyMu;
	// AppliedGen advances only when a generation is installed. DesiredGen ==
	// AppliedGen with Builds == 0 means the published state has no outstanding
	// generation or active build. It does not prove that no caller is still
	// inside Apply after installation, nor that a call has reached Go at all.
	DesiredGen uint64
	AppliedGen uint64
	Phase      string
}

// BuildFailure reports a recovered panic during a synchronous build.
type BuildFailure struct {
	Gen        uint64
	Recovered  any
	CleanupErr error
}

func (e *BuildFailure) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.CleanupErr != nil {
		return fmt.Sprintf("audio build generation %d panicked: %v (cleanup: %v)", e.Gen, e.Recovered, e.CleanupErr)
	}
	return fmt.Sprintf("audio build generation %d panicked: %v", e.Gen, e.Recovered)
}

func (e *BuildFailure) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.CleanupErr
}

// ShutdownPanic reports a recovered panic from a pipeline shutdown.
type ShutdownPanic struct {
	Gen       uint64
	Recovered any
}

func (e *ShutdownPanic) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("generation %d shutdown panicked: %v", e.Gen, e.Recovered)
}

type slot struct {
	p   ManagedPipeline
	gen uint64

	settings Settings
	asrReady bool
	ttsReady bool

	refs            int
	retired         bool
	shutdownStarted bool
	shutdownDone    chan struct{}
	shutdownErr     error
}

type closeResult struct {
	err   error
	phase string
}

// Manager owns desired state, active and retired generations, leases, build
// accounting, coalesced applies, and bounded drain.
type Manager struct {
	mu            sync.Mutex
	active        *slot
	retired       map[*slot]struct{}
	builds        int
	closing       bool
	drained       chan struct{}
	closeErrs     []error
	closeTimeouts int
	terminal      *closeResult

	desired    Settings
	desiredGen uint64
	appliedGen uint64
	status     Status

	lastBuildFailure *BuildFailure

	applyMu     sync.Mutex
	buildCtx    context.Context
	buildCancel context.CancelFunc
	ops         *opRegistry
	factory     PipelineFactory
}

// NewManager returns an empty manager with no active generation.
func NewManager(factory PipelineFactory) *Manager {
	buildCtx, buildCancel := context.WithCancel(context.Background())
	return &Manager{
		retired:     make(map[*slot]struct{}),
		status:      Status{Phase: "idle"},
		buildCtx:    buildCtx,
		buildCancel: buildCancel,
		ops:         newOpRegistry(),
		factory:     factory,
	}
}

// Apply records desired state, coalesces queued callers onto the newest
// generation, and builds synchronously on this goroutine. It returns the current
// status and whether this caller was superseded by a newer desired generation.
func (m *Manager) Apply(ctx context.Context, desired Settings) (Status, bool, error) {
	if m == nil {
		return Status{}, false, ErrManagerClosed
	}
	if ctx == nil {
		ctx = context.Background()
	}

	m.mu.Lock()
	if m.terminal != nil {
		st := m.statusLocked()
		err := m.terminal.err
		if err == nil {
			err = ErrManagerClosed
		}
		m.mu.Unlock()
		return st, false, err
	}
	if m.closing {
		st := m.statusLocked()
		m.mu.Unlock()
		return st, false, ErrManagerClosing
	}
	m.desired = desired
	m.desiredGen++
	requestGen := m.desiredGen
	m.status.ASR.Requested = desired.ASREnabled
	m.status.TTS.Requested = desired.TTSEnabled
	m.mu.Unlock()

	m.applyMu.Lock()
	defer m.applyMu.Unlock()

	for {
		m.mu.Lock()
		if m.terminal != nil {
			st := m.statusLocked()
			err := m.terminal.err
			if err == nil {
				err = ErrManagerClosed
			}
			m.mu.Unlock()
			return st, requestGen < m.desiredGen, err
		}
		if m.closing {
			st := m.statusLocked()
			m.mu.Unlock()
			return st, requestGen < m.desiredGen, ErrManagerClosing
		}

		latest := m.desired
		buildGen := m.desiredGen
		if m.canSkipLocked(latest, buildGen) {
			st := m.statusLocked()
			m.mu.Unlock()
			return st, requestGen < buildGen, nil
		}

		if !latest.ASREnabled && !latest.TTSEnabled {
			claimed := m.retireActiveLocked()
			m.appliedGen = buildGen
			m.status = Status{
				ASR:   Capability{Requested: false},
				TTS:   Capability{Requested: false},
				Phase: "idle",
			}
			m.evaluateDrainLocked()
			st := m.statusLocked()
			m.mu.Unlock()
			if claimed != nil {
				go m.runShutdown(claimed)
			}
			return st, requestGen < buildGen, nil
		}

		m.builds++
		m.status = Status{
			ASR:   Capability{Requested: latest.ASREnabled},
			TTS:   Capability{Requested: latest.TTSEnabled},
			Phase: "applying",
		}
		m.mu.Unlock()

		buildCtx, cancel := m.newBuildContext(ctx)
		candidate := &Candidate{}
		var factoryErr error
		var factoryPanic any
		if m.factory == nil {
			factoryErr = ErrNoFactory
		} else {
			func() {
				defer func() {
					if recovered := recover(); recovered != nil {
						factoryPanic = recovered
					}
				}()
				factoryErr = m.factory(buildCtx, latest, candidate)
			}()
		}
		cancel()

		m.mu.Lock()
		obsolete := m.closing || m.terminal != nil || m.desiredGen != buildGen
		m.mu.Unlock()
		if obsolete {
			cleanupErr, cleanupPanic := m.disposeCandidate(candidate)
			m.mu.Lock()
			_ = m.recordCleanupLocked(buildGen, cleanupErr, cleanupPanic)
			retry := !m.closing && m.terminal == nil && m.desiredGen != buildGen
			st := m.statusLocked()
			superseded := requestGen < m.desiredGen
			terminal := m.terminal
			m.mu.Unlock()
			if retry {
				continue
			}
			if terminal != nil {
				err := terminal.err
				if err == nil {
					err = ErrManagerClosed
				}
				return st, superseded, err
			}
			return st, superseded, ErrManagerClosing
		}

		if factoryPanic != nil {
			buildFailure := &BuildFailure{Gen: buildGen, Recovered: factoryPanic}
			m.mu.Lock()
			m.status = m.failureStatusLocked(latest, buildFailure)
			m.mu.Unlock()

			cleanupErr, cleanupPanic := m.disposeCandidate(candidate)
			m.mu.Lock()
			cleanupErr = m.recordCleanupLocked(buildGen, cleanupErr, cleanupPanic)
			buildFailure.CleanupErr = cleanupErr
			m.lastBuildFailure = buildFailure
			retry := !m.closing && m.terminal == nil && m.desiredGen != buildGen
			st := m.statusLocked()
			superseded := requestGen < m.desiredGen
			m.mu.Unlock()
			if retry {
				continue
			}
			return st, superseded, buildFailure
		}

		if factoryErr != nil {
			m.mu.Lock()
			m.status = m.failureStatusLocked(latest, factoryErr)
			m.mu.Unlock()

			cleanupErr, cleanupPanic := m.disposeCandidate(candidate)
			m.mu.Lock()
			cleanupErr = m.recordCleanupLocked(buildGen, cleanupErr, cleanupPanic)
			retry := !m.closing && m.terminal == nil && m.desiredGen != buildGen
			st := m.statusLocked()
			superseded := requestGen < m.desiredGen
			m.mu.Unlock()
			if retry {
				continue
			}
			if cleanupErr != nil {
				return st, superseded, errors.Join(factoryErr, cleanupErr)
			}
			return st, superseded, factoryErr
		}

		requestedCount, readyCount := candidateCounts(latest, candidate)
		capErr := candidateCapabilityErrors(latest, candidate)
		if readyCount == 0 && capErr == nil {
			capErr = ErrNoReadyCapability
		}

		switch {
		case candidate.Pipeline == nil && readyCount > 0:
			m.mu.Lock()
			if m.active != nil {
				m.status = m.rejectedStatusLocked(latest, candidate)
			} else {
				m.status = m.failureStatusLocked(latest, ErrInvalidCandidate)
			}
			m.mu.Unlock()

			cleanupErr, cleanupPanic := m.disposeCandidate(candidate)
			m.mu.Lock()
			cleanupErr = m.recordCleanupLocked(buildGen, cleanupErr, cleanupPanic)
			retry := !m.closing && m.terminal == nil && m.desiredGen != buildGen
			st := m.statusLocked()
			superseded := requestGen < m.desiredGen
			m.mu.Unlock()
			if retry {
				continue
			}
			if cleanupErr != nil {
				return st, superseded, errors.Join(ErrInvalidCandidate, cleanupErr)
			}
			return st, superseded, ErrInvalidCandidate

		case readyCount == 0:
			m.mu.Lock()
			if m.active != nil {
				m.status = m.rejectedStatusLocked(latest, candidate)
			} else {
				m.status = candidateStatus(latest, candidate, "failed")
			}
			m.mu.Unlock()

			cleanupErr, cleanupPanic := m.disposeCandidate(candidate)
			m.mu.Lock()
			cleanupErr = m.recordCleanupLocked(buildGen, cleanupErr, cleanupPanic)
			retry := !m.closing && m.terminal == nil && m.desiredGen != buildGen
			st := m.statusLocked()
			superseded := requestGen < m.desiredGen
			m.mu.Unlock()
			if retry {
				continue
			}
			if cleanupErr != nil {
				return st, superseded, errors.Join(capErr, cleanupErr)
			}
			return st, superseded, capErr

		case readyCount < requestedCount:
			if m.active == nil {
				m.mu.Lock()
				claimed := m.activateLocked(latest, candidate, buildGen, "degraded")
				st := m.statusLocked()
				superseded := requestGen < m.desiredGen
				m.mu.Unlock()
				if claimed != nil {
					go m.runShutdown(claimed)
				}
				return st, superseded, nil
			}

			m.mu.Lock()
			m.status = m.rejectedStatusLocked(latest, candidate)
			m.mu.Unlock()

			cleanupErr, cleanupPanic := m.disposeCandidate(candidate)
			m.mu.Lock()
			cleanupErr = m.recordCleanupLocked(buildGen, cleanupErr, cleanupPanic)
			retry := !m.closing && m.terminal == nil && m.desiredGen != buildGen
			st := m.statusLocked()
			superseded := requestGen < m.desiredGen
			m.mu.Unlock()
			if retry {
				continue
			}
			if cleanupErr != nil {
				return st, superseded, errors.Join(capErr, cleanupErr)
			}
			return st, superseded, capErr

		case candidate.Pipeline == nil:
			m.mu.Lock()
			if m.active != nil {
				m.status = m.rejectedStatusLocked(latest, candidate)
			} else {
				m.status = m.failureStatusLocked(latest, ErrInvalidCandidate)
			}
			m.mu.Unlock()

			cleanupErr, cleanupPanic := m.disposeCandidate(candidate)
			m.mu.Lock()
			cleanupErr = m.recordCleanupLocked(buildGen, cleanupErr, cleanupPanic)
			retry := !m.closing && m.terminal == nil && m.desiredGen != buildGen
			st := m.statusLocked()
			superseded := requestGen < m.desiredGen
			m.mu.Unlock()
			if retry {
				continue
			}
			if cleanupErr != nil {
				return st, superseded, errors.Join(ErrInvalidCandidate, cleanupErr)
			}
			return st, superseded, ErrInvalidCandidate

		default:
			m.mu.Lock()
			claimed := m.activateLocked(latest, candidate, buildGen, "idle")
			st := m.statusLocked()
			superseded := requestGen < m.desiredGen
			m.mu.Unlock()
			if claimed != nil {
				go m.runShutdown(claimed)
			}
			return st, superseded, nil
		}
	}
}

// Acquire returns a lease bound to the current active generation. The lease
// never follows a later replacement.
func (m *Manager) Acquire(ctx context.Context) (*Lease, error) {
	if m == nil {
		return nil, ErrManagerClosed
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.terminal != nil {
		return nil, ErrManagerClosed
	}
	if m.closing {
		return nil, ErrManagerClosing
	}
	if m.active == nil || m.active.p == nil {
		return nil, ErrNoActivePipeline
	}
	m.active.refs++
	return &Lease{m: m, slot: m.active}, nil
}

// Desired returns the current desired settings. A toggle changes one
// capability through this snapshot so the other capability's explicit engine
// is not silently discarded.
func (m *Manager) Desired() Settings {
	if m == nil {
		return Settings{}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.desired
}

// Status returns a consistent snapshot of the manager state.
func (m *Manager) Status() Status {
	if m == nil {
		return Status{}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.statusLocked()
}

// AppliedSettings returns the settings of the ACTIVE generation and whether one
// is installed. It deliberately does not return m.desired: desired is the newest
// request, which may not have been activated. Callers use this to decide whether
// a requested configuration is the effective configuration.
func (m *Manager) AppliedSettings() (Settings, bool) {
	if m == nil {
		return Settings{}, false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.active == nil {
		return Settings{}, false
	}
	return m.active.settings, true
}

// Close performs one bounded wait for the manager drain. A timeout is local to
// this call; a later Close may complete after outstanding work finishes.
func (m *Manager) Close(ctx context.Context) error {
	if m == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultCloseTimeout)
		defer cancel()
	}

	m.mu.Lock()
	if m.terminal != nil {
		err := m.terminal.err
		m.mu.Unlock()
		return err
	}

	var claimed *slot
	if !m.closing {
		m.closing = true
		m.status.Phase = "closing"
		claimed = m.retireActiveLocked()
		m.evaluateDrainLocked()
	}
	drained := m.drained
	if m.drainConditionLocked() {
		result := m.finalizeLocked()
		err := result.err
		m.mu.Unlock()
		return err
	}
	if drained == nil {
		m.drained = make(chan struct{})
		drained = m.drained
	}
	m.mu.Unlock()

	m.buildCancel()
	m.ops.cancelAll()
	m.interruptAll()
	if claimed != nil {
		go m.runShutdown(claimed)
	}

	select {
	case <-drained:
		m.mu.Lock()
		if m.terminal != nil {
			err := m.terminal.err
			m.mu.Unlock()
			return err
		}
		if !m.drainConditionLocked() {
			outstanding := m.outstandingLocked()
			m.mu.Unlock()
			return fmt.Errorf("%w: %d outstanding generation(s)", ErrDrainTimeout, outstanding)
		}
		result := m.finalizeLocked()
		err := result.err
		m.mu.Unlock()
		return err
	case <-ctx.Done():
		m.mu.Lock()
		if m.terminal != nil {
			err := m.terminal.err
			m.mu.Unlock()
			return err
		}
		if m.drainConditionLocked() {
			result := m.finalizeLocked()
			err := result.err
			m.mu.Unlock()
			return err
		}
		outstanding := m.outstandingLocked()
		m.closeTimeouts++
		m.mu.Unlock()
		return fmt.Errorf("%w: %d outstanding generation(s)", ErrDrainTimeout, outstanding)
	}
}

// BeginOperation registers a cancellable operation before it acquires a lease.
// End must be called after the lease is released. If the manager is already
// closing, the returned context is cancelled and End is a no-op.
func (m *Manager) BeginOperation(parent context.Context, timeout time.Duration) (*Operation, context.Context) {
	if parent == nil {
		parent = context.Background()
	}
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(parent, timeout)
	} else {
		ctx, cancel = context.WithCancel(parent)
	}

	if m == nil {
		cancel()
		return &Operation{}, ctx
	}

	m.mu.Lock()
	if m.closing || m.terminal != nil {
		m.mu.Unlock()
		cancel()
		return &Operation{}, ctx
	}
	id := m.ops.add(cancel)
	m.mu.Unlock()
	return &Operation{m: m, registry: m.ops, id: id, cancel: cancel}, ctx
}

// Operation is one registered manager operation. End is idempotent.
type Operation struct {
	m        *Manager
	registry *opRegistry
	id       uint64
	cancel   context.CancelFunc
	once     sync.Once
}

func (o *Operation) End() {
	if o == nil || o.registry == nil {
		return
	}
	o.once.Do(func() {
		if o.cancel != nil {
			o.cancel()
		}
		if o.registry.remove(o.id) && o.m != nil {
			o.m.evaluateDrain()
		}
	})
}

// Lease pins one generation and must be released exactly once (Release is
// idempotent).
type Lease struct {
	m    *Manager
	slot *slot
	once sync.Once
}

// Pipeline returns the leased generation. It remains valid until Release,
// even if the manager has installed a replacement.
func (l *Lease) Pipeline() ManagedPipeline {
	if l == nil || l.slot == nil {
		return nil
	}
	return l.slot.p
}

// Release drops this lease. A retired slot is shut down only after its last
// reference is released.
func (l *Lease) Release() {
	if l == nil || l.m == nil || l.slot == nil {
		return
	}
	l.once.Do(func() {
		m := l.m
		s := l.slot
		m.mu.Lock()
		if s.refs > 0 {
			s.refs--
		}
		start := s.refs == 0 && s.retired && !s.shutdownStarted
		if start {
			s.shutdownStarted = true
		}
		m.evaluateDrainLocked()
		m.mu.Unlock()
		if start {
			go m.runShutdown(s)
		}
	})
}

func (m *Manager) canSkipLocked(desired Settings, gen uint64) bool {
	// D5 idempotence belongs to T2.9.6. T2.9.2 only avoids rebuilding when this
	// exact desired generation has already been installed by the caller that
	// held applyMu ahead of this one.
	if m.appliedGen < gen {
		return false
	}
	if !desired.ASREnabled && !desired.TTSEnabled {
		return m.active == nil
	}
	// An applied generation is authoritative even when it was degraded; a
	// queued same-generation caller must not rebuild it.
	return m.active != nil && m.active.settings == desired
}

func (m *Manager) activateLocked(desired Settings, candidate *Candidate, gen uint64, phase string) *slot {
	var claimed *slot
	if m.active != nil {
		old := m.active
		if m.retireSlotLocked(old) {
			claimed = old
		}
	}

	m.active = &slot{
		p:            candidate.Pipeline,
		gen:          gen,
		settings:     desired,
		asrReady:     desired.ASREnabled && candidate.ASR.Ready,
		ttsReady:     desired.TTSEnabled && candidate.TTS.Ready,
		shutdownDone: make(chan struct{}),
	}
	m.builds--
	m.appliedGen = gen
	m.status = candidateStatus(desired, candidate, phase)
	m.evaluateDrainLocked()
	return claimed
}

func (m *Manager) retireActiveLocked() *slot {
	active := m.active
	if active == nil {
		return nil
	}
	if m.retireSlotLocked(active) {
		return active
	}
	return nil
}

func (m *Manager) retireSlotLocked(s *slot) bool {
	if s == nil || s.retired {
		return false
	}
	if m.active == s {
		m.active = nil
	}
	s.retired = true
	if s.shutdownDone == nil {
		s.shutdownDone = make(chan struct{})
	}
	m.retired[s] = struct{}{}
	if s.refs == 0 && !s.shutdownStarted {
		s.shutdownStarted = true
		return true
	}
	return false
}

func (m *Manager) runShutdown(s *slot) {
	var shutdownErr error
	var recovered any
	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = r
			}
		}()
		shutdownErr = s.p.Shutdown(m.buildCtx)
	}()

	m.mu.Lock()
	s.shutdownErr = shutdownErr
	switch {
	case recovered != nil:
		panicErr := &ShutdownPanic{Gen: s.gen, Recovered: recovered}
		s.shutdownErr = panicErr
		m.closeErrs = append(m.closeErrs, panicErr)
	case shutdownErr != nil:
		wrapped := fmt.Errorf("generation %d shutdown: %w", s.gen, shutdownErr)
		s.shutdownErr = wrapped
		m.closeErrs = append(m.closeErrs, wrapped)
	}
	close(s.shutdownDone)
	delete(m.retired, s)
	m.evaluateDrainLocked()
	m.mu.Unlock()
}

func (m *Manager) disposeCandidate(candidate *Candidate) (err error, recovered any) {
	if candidate == nil || candidate.Pipeline == nil {
		return nil, nil
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = r
			}
		}()
		err = candidate.Pipeline.Shutdown(m.buildCtx)
	}()
	return err, recovered
}

func (m *Manager) recordCleanupLocked(gen uint64, cleanupErr error, cleanupPanic any) error {
	var retained error
	switch {
	case cleanupPanic != nil:
		panicErr := &ShutdownPanic{Gen: gen, Recovered: cleanupPanic}
		m.closeErrs = append(m.closeErrs, panicErr)
		retained = panicErr
	case cleanupErr != nil:
		wrapped := fmt.Errorf("generation %d shutdown: %w", gen, cleanupErr)
		m.closeErrs = append(m.closeErrs, wrapped)
		retained = wrapped
	}
	m.builds--
	m.evaluateDrainLocked()
	return retained
}

func (m *Manager) evaluateDrain() {
	m.mu.Lock()
	m.evaluateDrainLocked()
	m.mu.Unlock()
}

func (m *Manager) evaluateDrainLocked() {
	if !m.closing {
		return
	}
	if m.drainConditionLocked() {
		if m.drained != nil {
			close(m.drained)
			m.drained = nil
		}
		return
	}
	if m.drained == nil {
		m.drained = make(chan struct{})
	}
}

func (m *Manager) drainConditionLocked() bool {
	return m.closing &&
		m.active == nil &&
		m.builds == 0 &&
		len(m.retired) == 0 &&
		m.ops.len() == 0
}

func (m *Manager) finalizeLocked() *closeResult {
	if m.terminal != nil {
		return m.terminal
	}
	result := &closeResult{phase: "closed"}
	if len(m.closeErrs) > 0 {
		result.err = errors.Join(m.closeErrs...)
		result.phase = "closed_with_errors"
	}
	m.terminal = result
	m.status.Phase = result.phase
	if m.drained != nil {
		close(m.drained)
		m.drained = nil
	}
	return result
}

func (m *Manager) outstandingLocked() int {
	count := m.builds + len(m.retired) + m.ops.len()
	if m.active != nil {
		count++
	}
	return count
}

func (m *Manager) statusLocked() Status {
	st := m.status
	st.ASR.Requested = m.desired.ASREnabled
	st.TTS.Requested = m.desired.TTSEnabled
	if m.active == nil {
		st.ASR.Ready = false
		st.TTS.Ready = false
	}
	st.InFlight = m.ops.len()
	if m.active != nil {
		st.InFlight += m.active.refs
	}
	for retired := range m.retired {
		st.InFlight += retired.refs
	}
	st.Builds = m.builds
	st.DesiredGen = m.desiredGen
	st.AppliedGen = m.appliedGen
	switch {
	case m.terminal != nil:
		st.Phase = m.terminal.phase
	case m.closing:
		st.Phase = "closing"
	case st.Phase == "":
		st.Phase = "idle"
	}
	return st
}

func (m *Manager) failureStatusLocked(desired Settings, err error) Status {
	st := Status{
		ASR:   Capability{Requested: desired.ASREnabled},
		TTS:   Capability{Requested: desired.TTSEnabled},
		Phase: "failed",
	}
	if m.active != nil {
		st.ASR.Ready = desired.ASREnabled && m.active.asrReady
		st.TTS.Ready = desired.TTSEnabled && m.active.ttsReady
	}
	if err != nil {
		if desired.ASREnabled {
			st.ASR.Error = err.Error()
		}
		if desired.TTSEnabled {
			st.TTS.Error = err.Error()
		}
	}
	return st
}

func (m *Manager) rejectedStatusLocked(desired Settings, candidate *Candidate) Status {
	st := Status{
		ASR:   Capability{Requested: desired.ASREnabled},
		TTS:   Capability{Requested: desired.TTSEnabled},
		Phase: "failed",
	}
	if m.active != nil {
		st.ASR.Ready = desired.ASREnabled && m.active.asrReady
		st.TTS.Ready = desired.TTSEnabled && m.active.ttsReady
	}
	if desired.ASREnabled && !candidate.ASR.Ready {
		st.ASR.Error = capabilityError("ASR", candidate.ASR.Err).Error()
	}
	if desired.TTSEnabled && !candidate.TTS.Ready {
		st.TTS.Error = capabilityError("TTS", candidate.TTS.Err).Error()
	}
	return st
}

func candidateStatus(desired Settings, candidate *Candidate, phase string) Status {
	return Status{
		ASR:   capabilityStatus("ASR", desired.ASREnabled, candidate.ASR),
		TTS:   capabilityStatus("TTS", desired.TTSEnabled, candidate.TTS),
		Phase: phase,
	}
}

func capabilityStatus(name string, requested bool, result CapabilityResult) Capability {
	status := Capability{Requested: requested}
	if !requested {
		return status
	}
	status.Ready = result.Ready
	if !result.Ready {
		status.Error = capabilityError(name, result.Err).Error()
	}
	return status
}

func candidateCounts(desired Settings, candidate *Candidate) (requested, ready int) {
	if desired.ASREnabled {
		requested++
		if candidate.ASR.Ready {
			ready++
		}
	}
	if desired.TTSEnabled {
		requested++
		if candidate.TTS.Ready {
			ready++
		}
	}
	return requested, ready
}

func candidateCapabilityErrors(desired Settings, candidate *Candidate) error {
	var errs []error
	if desired.ASREnabled && !candidate.ASR.Ready {
		errs = append(errs, fmt.Errorf("ASR: %w", capabilityError("ASR", candidate.ASR.Err)))
	}
	if desired.TTSEnabled && !candidate.TTS.Ready {
		errs = append(errs, fmt.Errorf("TTS: %w", capabilityError("TTS", candidate.TTS.Err)))
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

func capabilityError(name string, err error) error {
	if err != nil {
		return err
	}
	return fmt.Errorf("%s: %w", name, ErrCapabilityNotReady)
}

func (m *Manager) newBuildContext(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	stop := context.AfterFunc(m.buildCtx, cancel)
	return ctx, func() {
		stop()
		cancel()
	}
}

func (m *Manager) interruptAll() {
	m.mu.Lock()
	pipelines := make([]ManagedPipeline, 0, 1+len(m.retired))
	if m.active != nil && m.active.p != nil {
		pipelines = append(pipelines, m.active.p)
	}
	for s := range m.retired {
		if s.p != nil {
			pipelines = append(pipelines, s.p)
		}
	}
	m.mu.Unlock()

	for _, pipeline := range pipelines {
		interruptible, ok := pipeline.(interface{ Interrupt() })
		if !ok {
			continue
		}
		func() {
			defer func() {
				_ = recover()
			}()
			interruptible.Interrupt()
		}()
	}
}

type opRegistry struct {
	mu     sync.Mutex
	next   uint64
	active map[uint64]context.CancelFunc
}

func newOpRegistry() *opRegistry {
	return &opRegistry{active: make(map[uint64]context.CancelFunc)}
}

func (r *opRegistry) add(cancel context.CancelFunc) uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	r.active[r.next] = cancel
	return r.next
}

func (r *opRegistry) remove(id uint64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.active[id]; !ok {
		return false
	}
	delete(r.active, id)
	return true
}

func (r *opRegistry) cancelAll() {
	r.mu.Lock()
	cancels := make([]context.CancelFunc, 0, len(r.active))
	for _, cancel := range r.active {
		cancels = append(cancels, cancel)
	}
	r.mu.Unlock()
	for _, cancel := range cancels {
		cancel()
	}
}

func (r *opRegistry) len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.active)
}
