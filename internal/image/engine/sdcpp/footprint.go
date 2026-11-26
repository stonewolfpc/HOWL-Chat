package sdcpp

import (
	"howl-chat/internal/proc"
	"time"
)

// Footprint is one actual allocation observation of the sd-server child.
//
// This is deliberately separate from Placement: Placement is the engine's own
// per-component claim, while a working-set measurement is evidence Howl takes
// from the operating system. The pinned sd-server build does not publish a
// total-params line, so without this the UI could only ever say "unverified".
type Footprint struct {
	Known        bool      `json:"known"`
	WorkingSetMB float64   `json:"workingSetMB"`
	Phase        string    `json:"phase"`
	CapturedAt   time.Time `json:"capturedAt"`
}

// recordFootprint measures the live sd-server process. A failed measurement
// leaves the previous observation intact and never records a zero.
func (e *Engine) recordFootprint(phase string) {
	p := e.getServerProc()
	if p == nil {
		return
	}
	mb, ok := proc.WorkingSetMB(p.Pid)
	if !ok {
		return
	}
	e.footprintMu.Lock()
	if e.footprints == nil {
		e.footprints = map[string]Footprint{}
	}
	e.footprints[phase] = Footprint{
		Known:        true,
		WorkingSetMB: mb,
		Phase:        phase,
		CapturedAt:   time.Now().UTC(),
	}
	e.footprintMu.Unlock()
}

// GetFootprints returns a copy of the observations, keyed by "load" and
// "generation".
func (e *Engine) GetFootprints() map[string]Footprint {
	e.footprintMu.Lock()
	defer e.footprintMu.Unlock()
	out := make(map[string]Footprint, len(e.footprints))
	for phase, fp := range e.footprints {
		out[phase] = fp
	}
	return out
}

// resetFootprints drops observations of a process that no longer exists.
func (e *Engine) resetFootprints() {
	e.footprintMu.Lock()
	e.footprints = nil
	e.footprintMu.Unlock()
}
