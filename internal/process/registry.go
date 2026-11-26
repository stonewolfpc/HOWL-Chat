// Package process records the identity of child processes this application
// starts and the one goroutine that owns each child's Wait.
//
// It is deliberately about ownership, not scheduling: callers still decide
// whether to start an engine. A registered handle is the authority for the
// child's role, PID and start token, so a reused PID cannot be mistaken for the
// process a caller started.
package process

import (
	"fmt"
	"os/exec"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"howl-chat/internal/proc"
)

// Role names the one lifecycle slot a child occupies.
type Role string

const (
	RoleChat        Role = "llama-server"
	RoleImageServer Role = "sd-server"
	RoleVideoCLI    Role = "sd-cli"
	// RoleTTSChild is one short-lived synthesis or voice-enumeration child.
	// Several may run at once, so each occupies a scoped slot rather than the
	// single per-role slot chat/image/video use.
	RoleTTSChild Role = "tts-child"
	// RoleTTSService is a long-lived local speech service (MeloTTS).
	RoleTTSService Role = "tts-service"
)

// Identity distinguishes a child from an earlier process that had the same PID.
// The token is assigned when the child is adopted and never reused.
type Identity struct {
	Role  Role
	PID   int
	Token string
}

func (id Identity) String() string {
	return fmt.Sprintf("%s pid=%d token=%s", id.Role, id.PID, id.Token)
}

// Handle is one adopted child. Its waiter is the only goroutine allowed to
// call Wait on the command; every other observer waits on that result.
type Handle struct {
	Identity Identity
	Cmd      *exec.Cmd

	mu       sync.Mutex
	waiter   func(time.Duration) error
	registry *Registry
	instance string
}

// Registry owns the active child for each role. Singleton roles keep exactly
// one child; scoped slots (TTS children) allow several concurrent children of
// one role, each with its own stable instance identity.
type Registry struct {
	mu      sync.Mutex
	current map[Role]*Handle
	scoped  map[string]*Handle
	seq     atomic.Uint64
}

// NewRegistry returns an empty process registry.
func NewRegistry() *Registry {
	return &Registry{current: make(map[Role]*Handle), scoped: make(map[string]*Handle)}
}

var (
	defaultRegistryOnce sync.Once
	defaultRegistry     *Registry
)

// DefaultRegistry returns the process-wide registry used when a caller has no
// injected owner. Production injects the composition root's registry; this is
// the fallback for a package constructed without one, never a second authority
// alongside it.
func DefaultRegistry() *Registry {
	defaultRegistryOnce.Do(func() { defaultRegistry = NewRegistry() })
	return defaultRegistry
}

func scopedKey(role Role, instance string) string {
	return string(role) + "\x00" + instance
}

// Register adopts an already-started command under role. The waiter is the
// one Wait owner; this package never starts a second one while it is present.
func (r *Registry) Register(role Role, cmd *exec.Cmd, waiter func(time.Duration) error) (*Handle, error) {
	if r == nil {
		return nil, fmt.Errorf("process: nil registry")
	}
	if role == "" {
		return nil, fmt.Errorf("process: empty role")
	}
	if cmd == nil || cmd.Process == nil {
		return nil, fmt.Errorf("process: cannot adopt an unstarted command")
	}
	pid := cmd.Process.Pid
	token := fmt.Sprintf("%d-%d", time.Now().UnixNano(), r.seq.Add(1))
	h := &Handle{
		Identity: Identity{Role: role, PID: pid, Token: token},
		Cmd:      cmd,
		waiter:   waiter,
		registry: r,
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if existing := r.current[role]; existing != nil {
		return nil, fmt.Errorf("process: %s already owns %s", existing.Identity, h.Identity)
	}
	r.current[role] = h
	return h, nil
}

// RegisterChild adopts a started command into a scoped slot. Concurrent
// children of the same role occupy different instances; the instance is part
// of the child's stable identity, so a reused PID cannot be mistaken for the
// process that occupied the slot before.
func (r *Registry) RegisterChild(role Role, instance string, cmd *exec.Cmd, waiter func(time.Duration) error) (*Handle, error) {
	if r == nil {
		return nil, fmt.Errorf("process: nil registry")
	}
	if role == "" {
		return nil, fmt.Errorf("process: empty role")
	}
	if instance == "" {
		return nil, fmt.Errorf("process: empty instance")
	}
	if cmd == nil || cmd.Process == nil {
		return nil, fmt.Errorf("process: cannot adopt an unstarted command")
	}
	pid := cmd.Process.Pid
	h := &Handle{
		Identity: Identity{Role: role, PID: pid, Token: fmt.Sprintf("%d-%d", time.Now().UnixNano(), r.seq.Add(1))},
		Cmd:      cmd,
		waiter:   waiter,
		registry: r,
		instance: instance,
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.scoped == nil {
		r.scoped = make(map[string]*Handle)
	}
	key := scopedKey(role, instance)
	if existing := r.scoped[key]; existing != nil {
		return nil, fmt.Errorf("process: %s already owns %s", existing.Identity, h.Identity)
	}
	r.scoped[key] = h
	return h, nil
}

// Children returns every live child registered under role, singleton and
// scoped alike, in registration order.
func (r *Registry) Children(role Role) []*Handle {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*Handle, 0, 1+len(r.scoped))
	if h := r.current[role]; h != nil {
		out = append(out, h)
	}
	for _, h := range r.scoped {
		if h.Identity.Role == role {
			out = append(out, h)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Identity.Token < out[j].Identity.Token })
	return out
}

// Current returns the live handle for role, or nil.
func (r *Registry) Current(role Role) *Handle {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.current[role]
}

// Lookup returns the handle only when identity is still the active child.
func (r *Registry) Lookup(identity Identity) *Handle {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if h := r.current[identity.Role]; h != nil && h.Identity == identity {
		return h
	}
	for _, h := range r.scoped {
		if h.Identity == identity {
			return h
		}
	}
	return nil
}

// Remove releases handle only if it is still the active child for its role.
func (r *Registry) Remove(handle *Handle) {
	if r == nil || handle == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if handle.instance != "" {
		key := scopedKey(handle.Identity.Role, handle.instance)
		if h := r.scoped[key]; h == handle {
			delete(r.scoped, key)
		}
		return
	}
	if h := r.current[handle.Identity.Role]; h == handle {
		delete(r.current, handle.Identity.Role)
	}
}

// Stop kills and waits for the active child in role. A failed wait keeps the
// handle registered so a retry targets the same process.
func (r *Registry) Stop(role Role, timeout time.Duration) error {
	h := r.Current(role)
	if h == nil {
		return nil
	}
	return h.Stop(timeout)
}

// Stop kills and waits for this child. It is the only stop route for an
// adopted process; a timeout leaves the handle registered for a later retry.
func (h *Handle) Stop(timeout time.Duration) error {
	if h == nil || h.Cmd == nil || h.Cmd.Process == nil {
		if h != nil {
			h.Release()
		}
		return nil
	}
	pid := h.Identity.PID
	proc.KillTree(pid)
	_ = h.Cmd.Process.Kill()

	waiter := h.waiterForStop()
	if waiter == nil {
		return fmt.Errorf("%s pid %d has no exit observer", h.Identity.Role, pid)
	}
	if err := waiter(timeout); err != nil {
		return fmt.Errorf("%s pid %d did not exit within %s: %w", h.Identity.Role, pid, timeout, err)
	}
	h.Release()
	return nil
}

// Waiter returns the recorded exit observer, or nil when the caller provided
// none and no stop has installed a fallback.
func (h *Handle) Waiter() func(time.Duration) error {
	if h == nil {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.waiter
}

// Release removes this handle from its registry. It is idempotent.
func (h *Handle) Release() {
	if h == nil {
		return
	}
	h.registry.Remove(h)
}

// waiterForStop installs a fallback observer only for a hand-built test command
// that arrived without one. Production spawn paths always supply a monitor.
func (h *Handle) waiterForStop() func(time.Duration) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.waiter != nil {
		return h.waiter
	}
	done := make(chan struct{})
	pid := h.Identity.PID
	h.waiter = func(waitFor time.Duration) error {
		select {
		case <-done:
			return nil
		case <-time.After(waitFor):
			return fmt.Errorf("pid %d still running", pid)
		}
	}
	go func() {
		_, _ = h.Cmd.Process.Wait()
		close(done)
	}()
	return h.waiter
}
