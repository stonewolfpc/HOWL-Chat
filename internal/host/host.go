// Package host is HOWL Chat's Level 0 application authority.
//
// It owns lifecycle and coordination only: admission of new work, single-entry
// shutdown sequencing, joining of startup work, the one active stream token,
// supervision of child processes this application started, and the ordered
// reporting of shutdown-stage failures.
//
// It owns no domain state. Conversation, model, audio, image, memory and
// configuration values live with their own owners; the host decides *when*
// work may start or stop and *reports* what happened. A field here that stores
// those values would recreate the god-object this package exists to undo.
package host

// Authority is the set of coordination primitives available to every domain.
// It is a container for ownership boundaries, not a state store: each field
// owns one coordination concern and nothing else.
type Authority struct {
	// Lifecycle runs the application's ordered sequence: boot and wiring,
	// readiness, shutdown, and the restart that re-runs it.
	Lifecycle  *Lifecycle
	Gate       *Gate
	Background *Background
	Streams    *Streams
	Processes  *Processes
}

// NewAuthority builds the application's coordination authority. One gate is
// shared by the lifecycle, background work and every domain's admission check,
// so there is a single place that decides when new work may start.
func NewAuthority() *Authority {
	gate := NewGate()
	return &Authority{
		Lifecycle:  NewLifecycle(gate),
		Gate:       gate,
		Background: NewBackground(gate),
		Streams:    NewStreams(),
		Processes:  NewProcesses(),
	}
}
