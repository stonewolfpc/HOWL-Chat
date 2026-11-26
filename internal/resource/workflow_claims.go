package resource

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// ComponentKey is the stable identity of a shared resident component. Display
// names and file sizes are deliberately absent: the same upscaler reached
// through different spellings is one component only when its engine, role,
// normalized path, weight type and companion set all agree.
type ComponentKey struct {
	EngineID   string
	Role       string
	Path       string
	WeightType string
	Companions string
}

// NewComponentKey builds a normalized key. Companion order is not identity.
func NewComponentKey(engineID, role, path, weightType string, companions []string) ComponentKey {
	normalized := make([]string, 0, len(companions))
	for _, companion := range companions {
		if strings.TrimSpace(companion) == "" {
			continue
		}
		normalized = append(normalized, filepath.Clean(companion))
	}
	sort.Strings(normalized)
	return ComponentKey{
		EngineID:   engineID,
		Role:       role,
		Path:       filepath.Clean(path),
		WeightType: strings.ToLower(strings.TrimSpace(weightType)),
		Companions: strings.Join(normalized, "\x00"),
	}
}

func (k ComponentKey) String() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", k.EngineID, k.Role, k.Path, k.WeightType, k.Companions)
}

// WorkflowLease is the parent execution scope of one workflow request. It
// records the engine instance and owns the component claims made by its stages.
// A batch remains one parent lease; each item releases its own component claim.
type WorkflowLease struct {
	mu       sync.Mutex
	engineID string
	claims   map[ComponentKey]int
	released bool
}

func NewWorkflowLease(engineID string) *WorkflowLease {
	return &WorkflowLease{engineID: engineID, claims: make(map[ComponentKey]int)}
}

func (w *WorkflowLease) EngineID() string {
	if w == nil {
		return ""
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.engineID
}

func (w *WorkflowLease) Released() bool {
	if w == nil {
		return true
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.released
}

// Acquire increments the claim count for key. The same component claimed by a
// later stage is the same resident component, not a new allocation.
func (w *WorkflowLease) Acquire(key ComponentKey) (*ComponentClaim, error) {
	if w == nil {
		return nil, fmt.Errorf("workflow lease is not configured")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.released {
		return nil, fmt.Errorf("workflow lease has already been released")
	}
	w.claims[key]++
	return &ComponentClaim{lease: w, key: key}, nil
}

func (w *WorkflowLease) Verify(key ComponentKey) bool {
	if w == nil {
		return false
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return !w.released && w.claims[key] > 0
}

func (w *WorkflowLease) ClaimCount(key ComponentKey) int {
	if w == nil {
		return 0
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.claims[key]
}

func (w *WorkflowLease) ActiveClaimCount() int {
	if w == nil {
		return 0
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	total := 0
	for _, count := range w.claims {
		total += count
	}
	return total
}

// Release ends the parent scope and clears any claim a failed path did not
// release. It is idempotent.
func (w *WorkflowLease) Release() {
	if w == nil {
		return
	}
	w.mu.Lock()
	w.released = true
	w.claims = nil
	w.mu.Unlock()
}

func (w *WorkflowLease) releaseClaim(key ComponentKey) {
	if w == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.claims == nil {
		return
	}
	if w.claims[key] <= 1 {
		delete(w.claims, key)
		return
	}
	w.claims[key]--
}

// ComponentClaim is one stage's reference to a retained component.
type ComponentClaim struct {
	lease *WorkflowLease
	key   ComponentKey
	once  sync.Once
}

// Release decrements the claim exactly once.
func (c *ComponentClaim) Release() {
	if c == nil || c.lease == nil {
		return
	}
	c.once.Do(func() { c.lease.releaseClaim(c.key) })
}
