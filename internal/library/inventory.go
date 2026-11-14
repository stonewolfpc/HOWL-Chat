package library

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// InventoryRoot describes one registered content root. Items may be a live
// scan or the last-known snapshot; Stale distinguishes the two. A cached
// snapshot is never presented as a live scan.
type InventoryRoot struct {
	Registration
	Available bool   `json:"available"`
	Scanned   bool   `json:"scanned"`
	Stale     bool   `json:"stale"`
	ScanError string `json:"scan_error,omitempty"`
	ScannedAt string `json:"scanned_at,omitempty"`
	Items     []Item `json:"items"`
}

type cachedInventory struct {
	ScannedAt string `json:"scanned_at"`
	Items     []Item `json:"items"`
}

type persistedInventory struct {
	Version int                        `json:"version"`
	Roots   map[string]cachedInventory `json:"roots"`
}

// InventoryStore persists last-known metadata for registered roots. It is a
// cache, never an authority: content files remain the source of truth.
type InventoryStore struct {
	mu    sync.RWMutex
	path  string
	roots map[string]cachedInventory
}

// NewInventoryStore loads the bounded inventory cache if it exists. A missing
// or unreadable cache starts empty; it never prevents library registration
// from loading.
func NewInventoryStore(configDir string) *InventoryStore {
	s := &InventoryStore{
		path:  filepath.Join(configDir, "content-library-inventory.json"),
		roots: map[string]cachedInventory{},
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		return s
	}
	var stored persistedInventory
	if err := json.Unmarshal(data, &stored); err != nil {
		return s
	}
	for id, root := range stored.Roots {
		if id == "" {
			continue
		}
		s.roots[id] = root
	}
	return s
}

// Record stores a successful live scan. It is called only after a bounded
// scan completes, never for an unavailable root.
func (s *InventoryStore) Record(rootID string, items []Item) error {
	if s == nil || rootID == "" {
		return nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	copied := make([]Item, len(items))
	copy(copied, items)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.roots[rootID] = cachedInventory{ScannedAt: now, Items: copied}
	return s.saveLocked()
}

// Cached returns the last successful scan without modifying it.
func (s *InventoryStore) Cached(rootID string) (cachedInventory, bool) {
	if s == nil {
		return cachedInventory{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.roots[rootID]
	if !ok {
		return cachedInventory{}, false
	}
	entry.Items = append([]Item(nil), entry.Items...)
	return entry, true
}

// Clear removes every cached root. It is used by a completed content wipe so
// stale character entries cannot reappear from the last-known snapshot.
func (s *InventoryStore) Clear() error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	old := s.roots
	s.roots = map[string]cachedInventory{}
	if err := s.saveLocked(); err != nil {
		s.roots = old
		return err
	}
	return nil
}

// Forget removes cache for a registration that no longer exists.
func (s *InventoryStore) Forget(rootID string) {
	if s == nil || rootID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.roots[rootID]; !ok {
		return
	}
	delete(s.roots, rootID)
	_ = s.saveLocked()
}

func (s *InventoryStore) saveLocked() error {
	if s.path == "" {
		return nil
	}
	data, err := json.MarshalIndent(persistedInventory{Version: 1, Roots: s.roots}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode content inventory: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return fmt.Errorf("create content inventory directory: %w", err)
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("write content inventory: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("commit content inventory: %w", err)
	}
	return nil
}
