package host

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"howl-chat/internal/proc"
)

// Processes supervises the child process this application started — today the
// llama-server. It records the process and the single goroutine that owns its
// Wait, and it kills and waits in dependency order at shutdown.
//
// It does not decide *whether* to start an engine; that is runtime policy. It
// owns only the lifetime of a process already started, so exactly one place can
// wait on it and exactly one place can release it.
type Processes struct {
	mu     sync.Mutex
	cmd    *exec.Cmd
	waiter func(time.Duration) error
}

// NewProcesses returns an owner with no child process.
func NewProcesses() *Processes { return &Processes{} }

// Adopt records the process and the waiter that owns its single Wait. No other
// code may call Wait on this command.
func (p *Processes) Adopt(cmd *exec.Cmd, waiter func(time.Duration) error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cmd = cmd
	p.waiter = waiter
}

// Current reports the supervised process without releasing it.
func (p *Processes) Current() *exec.Cmd {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.cmd
}

// Take releases the process and its waiter from this owner.
func (p *Processes) Take() (*exec.Cmd, func(time.Duration) error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	cmd, waiter := p.cmd, p.waiter
	p.cmd, p.waiter = nil, nil
	return cmd, waiter
}

// Stop kills the owned process tree and waits for the existing monitor to
// observe exit. A timeout means the process did not exit; it is reported,
// never disguised as success.
func (p *Processes) Stop(timeout time.Duration) error {
	cmd, waiter := p.Take()
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	pid := cmd.Process.Pid
	proc.KillTree(pid)
	_ = cmd.Process.Kill()

	if waiter != nil {
		return waiter(timeout)
	}

	// A monitor is always installed by the production spawn path. This fallback
	// exists for partially initialized apps and test-owned commands: one
	// goroutine owns Wait, and shutdown only waits for its signal.
	done := make(chan struct{})
	go func() {
		_, _ = cmd.Process.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("llama-server pid %d did not exit within %s", pid, timeout)
	}
}
