package process

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"howl-chat/internal/proc"
)

// DefaultGrace bounds how long a child may take to exit after a graceful
// request (stdin close) and again after a forced tree kill. It is a wait bound,
// never a success claim: a child that does not exit inside it retains its
// registry slot and is reported as a failure.
const DefaultGrace = 5 * time.Second

// ErrTerminationUnconfirmed reports a child that survived forced termination.
// The registry slot is deliberately retained so a later stop targets the same
// process instead of pretending it is gone.
var ErrTerminationUnconfirmed = errors.New("child did not exit after forced termination; ownership retained")

// ChildSpec describes one child Howl owns for the life of an operation.
type ChildSpec struct {
	// Role is the lifecycle family. Scoped children of one role may run
	// concurrently; Instance distinguishes them.
	Role Role
	// Instance is the stable slot identity within the role.
	Instance string

	Path string
	Args []string
	Dir  string
	Env  []string

	// Stdin is written to the child and the pipe closed. Synthesis engines read
	// text until EOF, so closing it is also the graceful stop request.
	Stdin []byte

	// Capture merges stdout and stderr into a bounded diagnostic buffer.
	Capture bool
	// OutputLimit bounds the captured bytes (default 64 KiB).
	OutputLimit int
}

// Child is a job-owned child with exactly one registered waiter. Construction
// either returns a running, registered child or an error; there is no state in
// which work is considered running without an owner.
type Child struct {
	handle *Handle
	cmd    *exec.Cmd

	stdin     io.WriteCloser
	stdinOnce sync.Once

	output *childBuffer

	done    chan struct{}
	waitErr error

	settleOnce sync.Once
}

// StartChild starts spec under job ownership and registers it before the
// child's first instruction runs. The single Wait observer is created here and
// is the only Wait call for this process.
func StartChild(registry *Registry, spec ChildSpec) (*Child, error) {
	if registry == nil {
		registry = DefaultRegistry()
	}
	if spec.Role == "" {
		return nil, fmt.Errorf("process: child role is required")
	}
	if spec.Instance == "" {
		return nil, fmt.Errorf("process: child instance is required")
	}
	if spec.Path == "" {
		return nil, fmt.Errorf("process: child path is required")
	}

	cmd := exec.Command(spec.Path, spec.Args...)
	if spec.Dir != "" {
		cmd.Dir = spec.Dir
	}
	if len(spec.Env) > 0 {
		cmd.Env = spec.Env
	}
	proc.HideWindow(cmd)

	c := &Child{cmd: cmd, done: make(chan struct{}), output: newChildBuffer(spec.OutputLimit)}
	if spec.Capture {
		cmd.Stdout = c.output
		cmd.Stderr = c.output
	}
	if spec.Stdin != nil {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}
		c.stdin = stdin
	}

	if err := proc.StartWithOwnership(cmd, func() error {
		handle, err := registry.RegisterChild(spec.Role, spec.Instance, cmd, c.waitFor)
		if err != nil {
			return err
		}
		c.handle = handle
		return nil
	}); err != nil {
		return nil, err
	}

	go func() {
		c.waitErr = cmd.Wait()
		close(c.done)
	}()

	if c.stdin != nil {
		if len(spec.Stdin) > 0 {
			data := append([]byte(nil), spec.Stdin...)
			go func() {
				_, _ = c.stdin.Write(data)
				c.closeStdin()
			}()
		} else {
			c.closeStdin()
		}
	}
	return c, nil
}

// Identity is the registered role, PID and start token.
func (c *Child) Identity() Identity {
	if c == nil || c.handle == nil {
		return Identity{}
	}
	return c.handle.Identity
}

// PID is the registered process id, or zero when unregistered.
func (c *Child) PID() int {
	return c.Identity().PID
}

// Output returns the bounded captured diagnostics.
func (c *Child) Output() string {
	if c == nil || c.output == nil {
		return ""
	}
	return c.output.String()
}

// Exited reports that the waiter observed exit.
func (c *Child) Exited() bool {
	if c == nil {
		return true
	}
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}

// Err returns the waiter's exit error, if any.
func (c *Child) Err() error {
	if c == nil {
		return nil
	}
	return c.waitErr
}

// Wait waits for exit. On context cancellation it escalates to process-tree
// termination and only returns after exit is confirmed; if confirmation fails
// the registry slot is retained and the error says so.
func (c *Child) Wait(ctx context.Context, grace time.Duration) error {
	if c == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-c.done:
		c.settle()
		return c.waitErr
	case <-ctx.Done():
		if err := c.Terminate(ctx, grace); err != nil {
			return err
		}
		return ctx.Err()
	}
}

// Terminate asks gracefully (stdin close), waits grace, then force-kills the
// process tree and waits for the registered waiter again. It reports success
// only after exit is confirmed.
func (c *Child) Terminate(ctx context.Context, grace time.Duration) error {
	if c == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if grace <= 0 {
		grace = DefaultGrace
	}
	if c.Exited() {
		c.settle()
		return nil
	}

	// Step 1: graceful request. Closing stdin is what Piper and the MeloTTS
	// launcher observe as "finish and stop". A caller that is itself cancelled
	// skips the remainder of the grace period and escalates immediately.
	c.closeStdin()

	select {
	case <-c.done:
		c.settle()
		return nil
	case <-ctx.Done():
		// Escalate now.
	case <-time.After(grace):
	}

	// Step 2: escalate to process-tree termination. Howl owns this child, so it
	// owns everything the child spawned.
	if pid := c.PID(); pid > 0 {
		proc.KillTree(pid)
	}
	if c.cmd != nil && c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
	}

	// Step 3: wait for the registered waiter to observe exit. This wait is
	// bounded by the grace interval and deliberately NOT by the caller's
	// context: in the cancellation path that context is already cancelled, and
	// treating that as "termination unconfirmed" would abandon a process that
	// died microseconds later.
	if waitFor(c.done, grace) {
		c.settle()
		return nil
	}
	return fmt.Errorf("%w: %s", ErrTerminationUnconfirmed, c.Identity())
}

// waitFor reports whether done closed within d.
func waitFor(done <-chan struct{}, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-done:
		return true
	case <-timer.C:
		return false
	}
}

// waitFor is the Handle waiter: the only consumer of the child's exit.
func (c *Child) waitFor(timeout time.Duration) error {
	if timeout <= 0 {
		timeout = DefaultGrace
	}
	if waitFor(c.done, timeout) {
		return nil
	}
	return fmt.Errorf("%s still running after %s", c.Identity(), timeout)
}

func (c *Child) closeStdin() {
	if c == nil || c.stdin == nil {
		return
	}
	c.stdinOnce.Do(func() { _ = c.stdin.Close() })
}

// settle releases the registry slot exactly once, after exit was observed.
func (c *Child) settle() {
	c.settleOnce.Do(func() {
		if c.handle != nil {
			c.handle.Release()
		}
	})
}

type childBuffer struct {
	mu    sync.Mutex
	limit int
	data  []byte
}

func newChildBuffer(limit int) *childBuffer {
	if limit <= 0 {
		limit = 64 * 1024
	}
	return &childBuffer{limit: limit}
}

func (b *childBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = append(b.data, p...)
	if len(b.data) > b.limit {
		b.data = append([]byte(nil), b.data[len(b.data)-b.limit:]...)
	}
	return len(p), nil
}

func (b *childBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return string(append([]byte(nil), b.data...))
}

// StopRole terminates every live child registered under role, escalating each
// one to confirmed exit. Any child whose termination cannot be confirmed stays
// registered and its error is returned; the caller must not report clean
// shutdown in that case.
func StopRole(registry *Registry, role Role, grace time.Duration) error {
	if registry == nil {
		return nil
	}
	var errs []error
	for _, handle := range registry.Children(role) {
		if err := handle.Cmd.Process.Kill(); err != nil {
			// The process may have exited between Children and Kill; the waiter
			// below decides whether exit is real.
		}
		if pid := handle.Identity.PID; pid > 0 {
			proc.KillTree(pid)
		}
		waiter := handle.Waiter()
		if waiter == nil {
			errs = append(errs, fmt.Errorf("%s pid %d has no exit observer", handle.Identity.Role, handle.Identity.PID))
			continue
		}
		if err := waiter(grace); err != nil {
			errs = append(errs, fmt.Errorf("%w: %s", ErrTerminationUnconfirmed, handle.Identity))
			continue
		}
		handle.Release()
	}
	return errors.Join(errs...)
}
