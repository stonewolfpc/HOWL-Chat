package host

import (
	"os/exec"
	"sync"
	"time"

	"howl-chat/internal/process"
)

// Processes supervises the child processes this application started. The
// process registry is the authority for each child's role, PID and start token;
// this owner keeps the chat slot and its one recorded Wait observer.
//
// It does not decide *whether* to start an engine; that is runtime policy. It
// owns only the lifetime of a process already started, so exactly one place can
// wait on it and exactly one place can release it.
type Processes struct {
	mu       sync.Mutex
	registry *process.Registry
	handle   *process.Handle
}

// NewProcesses returns an owner with its own process registry.
func NewProcesses() *Processes {
	return NewProcessesWithRegistry(process.NewRegistry())
}

// NewProcessesWithRegistry returns an owner that records adopted children in
// the supplied registry. The composition root shares one registry across every
// process-owning domain.
func NewProcessesWithRegistry(registry *process.Registry) *Processes {
	if registry == nil {
		registry = process.NewRegistry()
	}
	return &Processes{registry: registry}
}

// Registry returns the process registry this owner records into.
func (p *Processes) Registry() *process.Registry {
	if p == nil {
		return nil
	}
	return p.registry
}

// Adopt records the process and the waiter that owns its single Wait. No other
// code may call Wait on this command.
func (p *Processes) Adopt(cmd *exec.Cmd, waiter func(time.Duration) error) (*process.Handle, error) {
	handle, err := p.registry.Register(process.RoleChat, cmd, waiter)
	if err != nil {
		return nil, err
	}
	p.mu.Lock()
	p.handle = handle
	p.mu.Unlock()
	return handle, nil
}

// Current reports the supervised process without releasing it.
func (p *Processes) Current() *exec.Cmd {
	handle := p.Handle()
	if handle == nil {
		return nil
	}
	return handle.Cmd
}

// Handle reports the current chat process identity.
func (p *Processes) Handle() *process.Handle {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.handle
}

// Owns reports whether handle is still the registered chat child.
func (p *Processes) Owns(handle *process.Handle) bool {
	if p == nil || handle == nil {
		return false
	}
	return p.registry.Lookup(handle.Identity) == handle
}

// Take releases the process and its waiter from this owner.
func (p *Processes) Take() (*exec.Cmd, func(time.Duration) error) {
	if p == nil {
		return nil, nil
	}
	p.mu.Lock()
	handle := p.handle
	p.handle = nil
	p.mu.Unlock()
	if handle == nil {
		return nil, nil
	}
	handle.Release()
	return handle.Cmd, handle.Waiter()
}

// Stop kills the owned process tree and waits for the existing monitor to
// observe exit. A timeout means the process did not exit; it is reported,
// never disguised as success, and ownership is retained for a retry.
func (p *Processes) Stop(timeout time.Duration) error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	handle := p.handle
	p.mu.Unlock()
	if handle == nil {
		return nil
	}
	if err := handle.Stop(timeout); err != nil {
		return err
	}
	p.clear(handle)
	return nil
}

// clear releases the handle only if it is still the process this Stop owned.
func (p *Processes) clear(handle *process.Handle) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.handle == handle {
		p.handle = nil
	}
}
