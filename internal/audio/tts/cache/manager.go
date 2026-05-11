package cache

import (
	"container/list"
	"sync"

	"howl-chat/internal/audio/types"
)

type entry struct {
	key    string
	result *types.SynthesisResult
	size   int64
}

// Manager implements a bounded in-memory LRU cache for synthesized audio metadata.
type Manager struct {
	mu          sync.Mutex
	maxSize     int64
	currentSize int64
	items       map[string]*list.Element
	order       *list.List
	stats       types.CacheStats
}

// NewManager creates a cache with a max size in MB.
func NewManager(maxSizeMB int64) *Manager {
	if maxSizeMB <= 0 {
		maxSizeMB = 100
	}
	maxBytes := maxSizeMB * 1024 * 1024
	return &Manager{
		maxSize: maxBytes,
		items:   make(map[string]*list.Element),
		order:   list.New(),
		stats: types.CacheStats{
			MaxSizeBytes: maxBytes,
		},
	}
}

// Get retrieves a cached synthesis result by key.
func (m *Manager) Get(key string) (*types.SynthesisResult, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	elem, ok := m.items[key]
	if !ok {
		m.stats.MissCount++
		return nil, false
	}
	m.order.MoveToFront(elem)
	m.stats.HitCount++
	return cloneResult(elem.Value.(*entry).result), true
}

// Store stores a synthesis result in cache.
func (m *Manager) Store(key string, result *types.SynthesisResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	size := estimateResultSize(result)

	if elem, ok := m.items[key]; ok {
		existing := elem.Value.(*entry)
		m.currentSize -= existing.size
		existing.result = cloneResult(result)
		existing.size = size
		m.currentSize += size
		m.order.MoveToFront(elem)
		m.rebalance()
		m.refreshStats()
		return nil
	}

	record := &entry{
		key:    key,
		result: cloneResult(result),
		size:   size,
	}
	elem := m.order.PushFront(record)
	m.items[key] = elem
	m.currentSize += size
	m.rebalance()
	m.refreshStats()
	return nil
}

// Invalidate removes a specific cache entry.
func (m *Manager) Invalidate(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if elem, ok := m.items[key]; ok {
		m.removeElement(elem)
	}
	m.refreshStats()
}

// Clear removes all cache entries.
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[string]*list.Element)
	m.order.Init()
	m.currentSize = 0
	m.refreshStats()
}

// GetStats returns cache counters and usage information.
func (m *Manager) GetStats() types.CacheStats {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.stats
}

func (m *Manager) rebalance() {
	for m.currentSize > m.maxSize && m.order.Len() > 0 {
		elem := m.order.Back()
		m.removeElement(elem)
	}
}

func (m *Manager) removeElement(elem *list.Element) {
	if elem == nil {
		return
	}
	ent := elem.Value.(*entry)
	delete(m.items, ent.key)
	m.order.Remove(elem)
	m.currentSize -= ent.size
	if m.currentSize < 0 {
		m.currentSize = 0
	}
}

func (m *Manager) refreshStats() {
	m.stats.EntryCount = len(m.items)
	m.stats.SizeBytes = m.currentSize
}

func estimateResultSize(result *types.SynthesisResult) int64 {
	if result == nil {
		return 0
	}
	size := int64(len(result.AudioPath) + 128)
	if result.Metadata.BitRate > 0 && result.Duration > 0 {
		size += int64(result.Duration * float64(result.Metadata.BitRate*125))
	}
	return size
}

func cloneResult(result *types.SynthesisResult) *types.SynthesisResult {
	if result == nil {
		return nil
	}
	copy := *result
	return &copy
}
