package sdcpp

// One generation request at a time, per engine process.
//
// There is no established way to make sd-server stop work it has already
// accepted (build/logs/cancelled-generation-contract-audit.md), and two
// requests overlapping on one process is the state in which an abandoned
// generation and its successor collide. So ownership belongs to the process,
// not to whoever happens to be waiting:
//
//	idle      a request may start
//	active    the owner is running; a second request is refused
//	abandoned the backend's state is unknown (client cancellation, transport
//	          loss, poll timeout, server-side failure), so the process stays
//	          unavailable until completion or exit is established
//
// Nothing in here claims backend computation stopped. The only completions this
// file accepts without the backend's cooperation are a successful response and
// a confirmed process exit: a dead process cannot own a generation.

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type generationState int

const (
	genIdle generationState = iota
	genActive
	genAbandoned
)

// genClose records why a request was closed from this side. It is the difference
// between a user cancelling, a model reload, a crash, and an unload — reasons
// that used to collapse into "client request cancelled", telling the user they
// had done something they had not.
//
// The value belongs to one generation (genHandle.close), not to the engine: a
// later request must not be able to rewrite how an earlier one ended.
type genClose int

const (
	closeNone genClose = iota
	closeForReload
	closeForShutdown
	closeProcessExited
)

// genOutcome is what a finished request proves about the backend.
//
//	completed      a successful response was read: the backend finished this work.
//	notDispatched  the failure happened before a request was sent (marshal,
//	               request build, missing or unreadable input): no backend work
//	               existed, so nothing can still be running.
//	cancelled      our own request was closed; the audit established this is not
//	               proof that the backend stopped.
//	ambiguous      the backend was contacted and its completion was not observed:
//	               transport loss, poll timeout, an error status, an empty or
//	               malformed result. Whether work continues is unestablished, so
//	               the slot keeps its owner.
type genOutcome int

const (
	genOutcomeCompleted genOutcome = iota
	genOutcomeNotDispatched
	genOutcomeCancelled
	genOutcomeAmbiguous
)

var (
	// errEngineBusy reports an overlapping request, not a failure of the first.
	errEngineBusy = errors.New("image engine is already generating")
	// errEngineUnavailable is the honest state after work was abandoned: it may
	// still be running, so a second request must not start.
	errEngineUnavailable = errors.New("image engine is unavailable until it is restarted: the previous generation was abandoned and its completion could not be confirmed")
	// errClientCancelled marks a request a job closed itself, never a stopped
	// backend.
	errClientCancelled = errors.New("client request cancelled")
	// errEngineUncertain reports a failed request whose backend state is unknown.
	// It is distinct from a cancellation because the user did not ask for it.
	errEngineUncertain = errors.New("image generation failed and the engine may still be working on it, so it stays unavailable until it is restarted")
	// errEngineStopped reports work ended by an unload, a model switch or a
	// process exit, never by the user cancelling.
	errEngineStopped = errors.New("image engine stopped before the request finished")
	// errNotDispatched marks a failure that never reached the backend.
	errNotDispatched = errors.New("request was never sent to the backend")
	// errGenerationSettled reports a second settlement of a generation that has
	// already reported its outcome: a programming error, never a new result.
	errGenerationSettled = errors.New("generation was already settled")
)

// notDispatched marks a local failure as one that never left this process, so
// the generation slot can be released without guessing about the backend.
//
// The inner error is wrapped with %w, not printed: callers keep errors.Is/As for
// the file, encoding or transport failure underneath.
func notDispatched(err error) error {
	return fmt.Errorf("%w: %w", errNotDispatched, err)
}

// genHandle is one generation's ownership token. The token, not the context, is
// what identifies it: function values cannot be compared.
type genHandle struct {
	token  uint64
	ctx    context.Context
	owner  string
	cancel context.CancelFunc
	// close is why this request was closed from this side, and it belongs to the
	// generation. Guarded by Engine.genMu.
	close genClose
	// settled records that this generation has already reported its outcome, so a
	// second settlement cannot transition the slot or re-describe the context that
	// cleanup cancelled. Guarded by Engine.genMu.
	settled bool
	// done closes when the request has returned and the generation has settled.
	// Shutdown joins this, not the backend's work: cancellation is not proof the
	// backend stopped, but it is proof this execution scope is no longer running.
	done     chan struct{}
	doneOnce sync.Once
}

func (h *genHandle) finish() {
	if h == nil {
		return
	}
	h.doneOnce.Do(func() { close(h.done) })
}

// processConfirmedGone reports whether the watched child has exited, which is
// the only termination this engine can establish on its own.
func (e *Engine) processConfirmedGone() bool {
	done := e.processDone()
	if done == nil {
		return false
	}
	select {
	case <-done:
		return true
	default:
		return false
	}
}

// beginGeneration claims the slot for jobID, or explains why it cannot.
//
// An empty jobID is allowed for the two legacy direct callers that have no job
// (GenerateImage and SetActiveImageModel), but it can never be cancelled: see
// cancelGeneration.
func (e *Engine) beginGeneration(jobID string) (*genHandle, error) {
	return e.beginGenerationWithContext(jobID, context.Background())
}

// beginGenerationWithContext binds a generation's execution scope to the
// caller's cancellation. A cancelled parent is refused before a slot is taken;
// otherwise the generation context is a child, so the caller can stop the
// request while Shutdown still owns the final join and process termination.
func (e *Engine) beginGenerationWithContext(jobID string, parent context.Context) (*genHandle, error) {
	if parent == nil {
		parent = context.Background()
	}
	if err := parent.Err(); err != nil {
		return nil, err
	}

	// A confirmed exit ends whatever the process was doing: release the slot so a
	// fresh child can generate without a restart of the app.
	if e.generationBusy() && e.processConfirmedGone() {
		e.resetGeneration(closeProcessExited)
	}

	e.genMu.Lock()
	defer e.genMu.Unlock()

	switch e.genState {
	case genActive:
		return nil, fmt.Errorf("%w (owner %q)", errEngineBusy, e.genOwner)
	case genAbandoned:
		return nil, fmt.Errorf("%w (owner %q)", errEngineUnavailable, e.genOwner)
	}

	ctx, cancel := context.WithCancel(parent)
	e.genSeq++
	h := &genHandle{token: e.genSeq, ctx: ctx, owner: jobID, cancel: cancel, close: closeNone, done: make(chan struct{})}
	e.genState = genActive
	e.genOwner = jobID
	e.genCancel = cancel
	e.genHandle = h
	return h, nil
}

// endGeneration runs when the request returned, with the outcome already
// classified from facts captured before this call. It never re-derives
// cancellation here: cleanup must not be able to turn an ordinary failure into a
// client cancellation.
//
// It returns the close reason captured for this generation while the slot was
// still held, so describing the outcome afterwards cannot be rewritten by
// whatever claims the slot next.
func (e *Engine) endGeneration(h *genHandle, outcome genOutcome) genClose {
	if h == nil {
		return closeNone
	}
	retain := outcome == genOutcomeCancelled || outcome == genOutcomeAmbiguous
	gone := false
	if retain {
		gone = e.processConfirmedGone()
	}

	e.genMu.Lock()
	reason := h.close
	if e.genSeq == h.token {
		if retain && !gone {
			e.genState = genAbandoned
			e.genOwner = h.owner
		} else {
			e.genState, e.genOwner = genIdle, ""
		}
		e.genCancel = nil
		e.genHandle = nil
		h.settled = true
	}
	// A handle settles exactly once. Invalidate its token under the same lock that
	// decides the transition: clearing genHandle/genCancel was not enough, because a
	// second settlement of the same handle still matched genSeq and could push an
	// idle engine back to abandoned.
	h.token = 0
	e.genMu.Unlock()
	h.cancel()
	h.finish()
	return reason
}

// settleGeneration classifies one finished request, moves the slot, and reports
// what the caller should surface. Every generation path calls this instead of
// reading the context after cleanup.
func (e *Engine) settleGeneration(h *genHandle, err error) error {
	if h == nil {
		return err
	}
	outcome := genOutcomeCompleted
	if err != nil {
		outcome = classifyGeneration(h.ctx.Err() != nil, err)
	}
	return e.settleGenerationAs(h, err, outcome)
}

// settleGenerationKnown is for a backend that explicitly reported a terminal
// state even when the result is an error. A DGL job that says "failed" has
// finished; treating it as an abandoned generation would recycle a healthy
// process for no reason and hide the definite backend answer behind "unknown".
func (e *Engine) settleGenerationKnown(h *genHandle, err error) error {
	return e.settleGenerationAs(h, err, genOutcomeCompleted)
}

func (e *Engine) settleGenerationAs(h *genHandle, err error, outcome genOutcome) error {
	if e.handleSettled(h) {
		// This generation already reported. Describing it again would read the
		// context cleanup cancelled and call an ordinary late error a client
		// cancellation — the mislabelling this file exists to prevent.
		h.finish()
		if err == nil {
			return fmt.Errorf("%w", errGenerationSettled)
		}
		return fmt.Errorf("%w (later error not reported: %v)", errGenerationSettled, err)
	}
	// Measure after the request returned, while the process that did the work is
	// still the same one that loaded the weights. A failed request is recorded
	// too: its allocation is exactly what a later diagnosis needs.
	e.recordFootprint("generation")
	reason := e.endGeneration(h, outcome)
	if err == nil {
		return nil
	}
	return describeGenError(outcome, reason, err)
}

// handleSettled reports whether this generation has already reported its outcome.
func (e *Engine) handleSettled(h *genHandle) bool {
	e.genMu.Lock()
	defer e.genMu.Unlock()
	return h.settled
}

// classifyGeneration decides what an ending request proves about the backend.
func classifyGeneration(cancelled bool, err error) genOutcome {
	switch {
	case err == nil:
		return genOutcomeCompleted
	case cancelled:
		return genOutcomeCancelled
	case errors.Is(err, errNotDispatched):
		return genOutcomeNotDispatched
	default:
		return genOutcomeAmbiguous
	}
}

// describeGenError states the outcome without overstating it: a closed request is
// not a stopped backend, a failed request is not a cancellation, and an unload is
// not something the user asked for by cancelling.
//
// Every branch wraps the underlying error with %w so the file, HTTP or transport
// error underneath stays reachable with errors.Is/As.
func describeGenError(outcome genOutcome, reason genClose, err error) error {
	switch outcome {
	case genOutcomeCancelled:
		switch reason {
		case closeForShutdown:
			return fmt.Errorf("%w (image engine unloaded): %w", errEngineStopped, err)
		case closeForReload:
			return fmt.Errorf("%w (image engine reloaded for another model): %w", errEngineStopped, err)
		case closeProcessExited:
			return fmt.Errorf("%w (engine process exited): %w", errEngineStopped, err)
		}
		return fmt.Errorf("%w: sd-server may still be running the abandoned work, so the engine stays unavailable until it is restarted (%w)", errClientCancelled, err)
	case genOutcomeAmbiguous:
		return fmt.Errorf("%w: %w", errEngineUncertain, err)
	default:
		// Completed, or a local failure that never reached the backend: report it
		// as itself.
		return err
	}
}

// cancelGeneration closes the client request owned by jobID.
//
// It reports whether jobID actually owned a running generation. A cancel that
// belongs to some other job — or to no job at all — is a no-op: identity is
// required, so an empty id can never mean "whichever job is active".
func (e *Engine) cancelGeneration(jobID string) bool {
	if jobID == "" {
		return false
	}
	e.genMu.Lock()
	defer e.genMu.Unlock()
	if e.genState != genActive || e.genCancel == nil || e.genOwner != jobID {
		return false
	}
	// A job cancel leaves the reason at closeNone: this is the user asking.
	e.genCancel()
	return true
}

// resetGeneration returns the slot to idle because the process that owned it is
// being replaced, shut down, or is confirmed gone. The reason is recorded on the
// generation being closed, so it survives the reset and cannot be overwritten by
// the next request.
func (e *Engine) resetGeneration(reason genClose) {
	e.genMu.Lock()
	cancel := e.genCancel
	h := e.genHandle
	// Bump the token so a request still unwinding cannot write its state back over
	// this reset: its handle no longer identifies the current owner.
	e.genSeq++
	if cancel != nil && h != nil {
		h.close = reason
	}
	e.genState, e.genOwner, e.genCancel, e.genHandle = genIdle, "", nil, nil
	e.genMu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (e *Engine) generationBusy() bool {
	e.genMu.Lock()
	defer e.genMu.Unlock()
	return e.genState != genIdle
}

// GenerationState reports the slot as "idle", "active" or "abandoned". It exists
// so the UI and the tests can tell "cancelled the request" from "stopped the
// work", which are not the same thing.
func (e *Engine) GenerationState() string {
	e.genMu.Lock()
	defer e.genMu.Unlock()
	switch e.genState {
	case genActive:
		return "active"
	case genAbandoned:
		return "abandoned"
	default:
		return "idle"
	}
}
