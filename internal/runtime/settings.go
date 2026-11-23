package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	appsettings "howl-chat/internal/settings"
)

// Legacy per-section files. They predate the unified store; boot still reads
// them when the unified document has no section for them, and the section
// bindings still write them so a downgrade keeps the user's values.
const (
	samplerSettingsFile = "sampler-settings.json"
	modelSettingsFile   = "model-settings.json"
)

// Settings is the settings authority. It owns the unified document, the live
// sampler and model sections the runtime reads on every turn, and persistence
// for all three. App holds this handle and the Wails bindings; it keeps no
// settings maps of its own.
//
// Sections are REPLACED, not merged, when a caller supplies a whole settings
// object: a key missing from it means the user's value is gone. Merging kept
// stale sampler values alive and made a cleared section look applied, which
// silently overrode per-family defaults.
//
// A nil *Settings is an uninitialized authority: reads return empty maps,
// writes are dropped, and Replace/Load report that it is not initialized.
type Settings struct {
	mu      sync.RWMutex
	store   *appsettings.Store
	dir     string
	doc     map[string]interface{}
	sampler map[string]interface{}
	model   map[string]interface{}
}

// NewSettings builds the settings authority rooted at dir. samplerDefaults is
// copied, so mutating the caller's map cannot change the authority later.
func NewSettings(store *appsettings.Store, dir string, samplerDefaults map[string]interface{}) *Settings {
	return &Settings{
		store:   store,
		dir:     dir,
		doc:     make(map[string]interface{}),
		sampler: copyMap(samplerDefaults),
		model:   make(map[string]interface{}),
	}
}

// Store returns the persistence store behind this authority.
func (s *Settings) Store() *appsettings.Store {
	if s == nil {
		return nil
	}
	return s.store
}

// Dir returns the directory the legacy section files live in.
func (s *Settings) Dir() string {
	if s == nil {
		return ""
	}
	return s.dir
}

// Load reads the unified document and restores the sampler and model sections.
// A section present in the document wins; when it is absent (or empty) the
// legacy per-section file is read instead. The first failure is returned, but
// later sections are still restored: one unreadable file must not discard the
// rest of the user's settings.
func (s *Settings) Load() error {
	if s == nil {
		return fmt.Errorf("settings not initialized")
	}
	if s.store == nil {
		return fmt.Errorf("settings store not initialized")
	}
	doc, err := s.store.Load()
	var firstErr error
	if err != nil {
		firstErr = err
		doc = make(map[string]interface{})
	}

	s.mu.Lock()
	s.doc = doc
	s.mu.Unlock()

	samplerSection, hasSampler := doc["sampler"].(map[string]interface{})
	modelSection, hasModel := doc["model"].(map[string]interface{})

	if hasSampler && len(samplerSection) > 0 {
		s.MergeSampler(samplerSection)
	} else if err := s.LoadSamplerFile(); err != nil && firstErr == nil {
		firstErr = err
	}

	if hasModel && len(modelSection) > 0 {
		s.MergeModel(modelSection)
	} else if err := s.LoadModelFile(); err != nil && firstErr == nil {
		firstErr = err
	}

	return firstErr
}

// All returns a snapshot of the unified document.
func (s *Settings) All() map[string]interface{} {
	if s == nil {
		return map[string]interface{}{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return copyMap(s.doc)
}

// Section returns a snapshot of one top-level document section.
func (s *Settings) Section(name string) map[string]interface{} {
	if s == nil {
		return map[string]interface{}{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	section, _ := s.doc[name].(map[string]interface{})
	return copyMap(section)
}

// Adopt installs a document and its sections without writing them. Boot uses it
// after an external load; tests use it to seed a store.
func (s *Settings) Adopt(doc map[string]interface{}) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.adoptLocked(doc)
}

// Replace installs the caller's whole settings document, replaces the sampler
// and model sections it carries, and persists the document.
func (s *Settings) Replace(doc map[string]interface{}) error {
	if s == nil || s.store == nil {
		return fmt.Errorf("settings store not initialized")
	}
	if doc == nil {
		return fmt.Errorf("settings document is required")
	}
	s.Adopt(doc)
	return s.store.Save(doc)
}

// adoptLocked installs doc and its sections. Callers hold the write lock.
func (s *Settings) adoptLocked(doc map[string]interface{}) {
	if doc == nil {
		doc = make(map[string]interface{})
	}
	s.doc = doc
	if sampler, ok := doc["sampler"].(map[string]interface{}); ok {
		s.sampler = copyMap(sampler)
	}
	if model, ok := doc["model"].(map[string]interface{}); ok {
		s.model = copyMap(model)
	}
}

// Sampler returns a snapshot of the live sampler section.
func (s *Settings) Sampler() map[string]interface{} {
	if s == nil {
		return map[string]interface{}{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return copyMap(s.sampler)
}

// SetSampler replaces the live sampler section.
func (s *Settings) SetSampler(values map[string]interface{}) {
	if s == nil {
		return
	}
	s.mu.Lock()
	s.sampler = copyMap(values)
	s.mu.Unlock()
}

// MergeSampler merges values into the live sampler section.
func (s *Settings) MergeSampler(values map[string]interface{}) {
	if s == nil {
		return
	}
	s.mu.Lock()
	for k, v := range values {
		s.sampler[k] = v
	}
	s.mu.Unlock()
}

// Model returns a snapshot of the live model section.
func (s *Settings) Model() map[string]interface{} {
	if s == nil {
		return map[string]interface{}{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return copyMap(s.model)
}

// SetModel replaces the live model section.
func (s *Settings) SetModel(values map[string]interface{}) {
	if s == nil {
		return
	}
	s.mu.Lock()
	s.model = copyMap(values)
	s.mu.Unlock()
}

// MergeModel merges values into the live model section.
func (s *Settings) MergeModel(values map[string]interface{}) {
	if s == nil {
		return
	}
	s.mu.Lock()
	for k, v := range values {
		s.model[k] = v
	}
	s.mu.Unlock()
}

// SaveSamplerFile writes the legacy sampler file from the live section.
func (s *Settings) SaveSamplerFile() error {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	values := copyMap(s.sampler)
	dir := s.dir
	s.mu.RUnlock()
	return writeSettingsFile(filepath.Join(dir, samplerSettingsFile), values, "sampler settings")
}

// LoadSamplerFile replaces the live sampler section from the legacy file. A
// missing file is not an error: the sampler keeps its defaults.
func (s *Settings) LoadSamplerFile() error {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	dir := s.dir
	s.mu.RUnlock()
	values, err := readSettingsFile(filepath.Join(dir, samplerSettingsFile), "sampler settings")
	if err != nil || values == nil {
		return err
	}
	s.SetSampler(values)
	return nil
}

// SaveModelFile writes the legacy model file from the live section.
func (s *Settings) SaveModelFile() error {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	values := copyMap(s.model)
	dir := s.dir
	s.mu.RUnlock()
	return writeSettingsFile(filepath.Join(dir, modelSettingsFile), values, "model settings")
}

// LoadModelFile replaces the live model section from the legacy file. A missing
// file is not an error: the section stays as it is.
func (s *Settings) LoadModelFile() error {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	dir := s.dir
	s.mu.RUnlock()
	values, err := readSettingsFile(filepath.Join(dir, modelSettingsFile), "model settings")
	if err != nil || values == nil {
		return err
	}
	s.SetModel(values)
	return nil
}

func readSettingsFile(path, label string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", label, err)
	}
	var values map[string]interface{}
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", label, err)
	}
	return values, nil
}

func writeSettingsFile(path string, values map[string]interface{}, label string) error {
	if values == nil {
		return nil
	}
	data, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", label, err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write %s: %w", label, err)
	}
	return nil
}

func copyMap(source map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(source))
	for k, v := range source {
		out[k] = v
	}
	return out
}
