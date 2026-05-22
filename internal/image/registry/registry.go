// Package registry maintains a persistent, directory-agnostic catalogue of
// every model asset that has been registered with HOWL.
//
// Design goals:
//   - Spans multiple scan directories — registering from D:\AI GGUF\Image and
//     later from E:\More Models both end up in the same registry.
//   - Survives restarts — written atomically to config/model-registry.json.
//   - Self-healing — stale entries (file deleted from disk) are pruned on every
//     Load so the UI never shows phantom models.
//   - No duplicates — keyed on absolute file path; re-registering the same
//     path updates its metadata in place.
//   - Thread-safe — all public methods hold a mutex.
package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"howl-chat/internal/image/scanner"
)

const registryFile = "model-registry.json"
const dirsFile = "model-scan-dirs.json"

// Entry is one record in the global registry.
type Entry struct {
	Name         string              `json:"name"`
	Path         string              `json:"path"`
	Type         string              `json:"type"`
	Family       string              `json:"family"`
	Format       string              `json:"format"`
	SizeMB       float64             `json:"size_mb"`
	Files        map[string]string   `json:"files,omitempty"`
	Recommended  scanner.Recommended `json:"recommended"`
	RegisteredAt time.Time           `json:"registered_at"`
}

// Registry is the in-memory + on-disk store.
type Registry struct {
	mu        sync.RWMutex
	configDir string
	entries   map[string]Entry  // keyed on absolute path
	dirs      map[string]string // dir path → asset type ("llm", "tts", "image")
}

// New creates a Registry that persists to configDir/model-registry.json.
// It loads any existing registry from disk immediately, pruning missing files.
func New(configDir string) (*Registry, error) {
	r := &Registry{
		configDir: configDir,
		entries:   make(map[string]Entry),
		dirs:      make(map[string]string),
	}
	if err := r.load(); err != nil {
		return nil, err
	}
	if err := r.loadDirs(); err != nil {
		return nil, err
	}
	return r, nil
}

// MergeEntries adds or updates a slice of plain Entry values.
// Used by LLM/TTS scanners that do not go through scanner.ModelIndex.
func (r *Registry) MergeEntries(entries []Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range entries {
		if e.RegisteredAt.IsZero() {
			e.RegisteredAt = time.Now()
		}
		r.entries[e.Path] = e
	}
	return r.save()
}

// ByType returns a snapshot of all live entries of the given asset type
// (e.g. "llm", "tts", "diffusion"). Files no longer on disk are pruned.
func (r *Registry) ByType(assetType string) []Entry {
	all := r.All() // All() already prunes + holds lock internally
	var out []Entry
	for _, e := range all {
		if e.Type == assetType {
			out = append(out, e)
		}
	}
	return out
}

// RegisterDir records a directory as a known scan location for the given
// asset type. Persisted so startup auto-rescan knows where to look.
func (r *Registry) RegisterDir(dir, assetType string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.dirs[dir] = assetType
	return r.saveDirs()
}

// RegisteredDirs returns all remembered scan directories as a map of
// dir → assetType ("llm", "tts", "image").
func (r *Registry) RegisteredDirs() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]string, len(r.dirs))
	for k, v := range r.dirs {
		out[k] = v
	}
	return out
}

// Merge adds or updates entries from a scanner.ModelIndex.
// Existing entries for the same path are overwritten with fresh metadata.
func (r *Registry) Merge(idx scanner.ModelIndex) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, e := range collectEntries(idx) {
		r.entries[e.Path] = Entry{
			Name:         e.Name,
			Path:         e.Path,
			Type:         string(e.Type),
			Family:       string(e.Family),
			Format:       e.Format,
			SizeMB:       e.SizeMB,
			Files:        e.Files,
			Recommended:  e.Recommended,
			RegisteredAt: time.Now(),
		}
	}
	return r.save()
}

// All returns a snapshot of all known entries, pruning any whose file has
// since been deleted from disk.
func (r *Registry) All() []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()

	var pruned []string
	out := make([]Entry, 0, len(r.entries))
	for path, e := range r.entries {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			pruned = append(pruned, path)
			continue
		}
		out = append(out, e)
	}
	if len(pruned) > 0 {
		for _, p := range pruned {
			delete(r.entries, p)
		}
		_ = r.save() // best-effort; ignore error on prune
	}
	return out
}

// Remove deletes a single entry from the registry (does NOT delete the file).
func (r *Registry) Remove(path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, path)
	return r.save()
}

// ToModelIndex converts the current registry snapshot back into a ModelIndex
// so existing engine/scanner code can consume it without changes.
func (r *Registry) ToModelIndex() scanner.ModelIndex {
	entries := r.All()
	var idx scanner.ModelIndex
	for _, e := range entries {
		ae := scanner.AssetEntry{
			Name:        e.Name,
			Path:        e.Path,
			Type:        scanner.AssetType(e.Type),
			Family:      scanner.ModelFamily(e.Family),
			Format:      e.Format,
			SizeMB:      e.SizeMB,
			Files:       e.Files,
			Recommended: e.Recommended,
		}
		switch ae.Type {
		case scanner.AssetDiffusion:
			idx.Models = append(idx.Models, ae)
		case scanner.AssetVAE:
			idx.VAEs = append(idx.VAEs, ae)
		case scanner.AssetLoRA:
			idx.LoRAs = append(idx.LoRAs, ae)
		case scanner.AssetControlNet:
			idx.ControlNets = append(idx.ControlNets, ae)
		case scanner.AssetEmbedding:
			idx.Embeddings = append(idx.Embeddings, ae)
		case scanner.AssetUpscaler:
			idx.Upscalers = append(idx.Upscalers, ae)
		}
	}
	return idx
}

// ---------------------------------------------------------------------------
// persistence helpers
// ---------------------------------------------------------------------------

func (r *Registry) path() string     { return filepath.Join(r.configDir, registryFile) }
func (r *Registry) dirsPath() string { return filepath.Join(r.configDir, dirsFile) }

// load reads the registry file, creating an empty one if absent,
// and prunes entries whose files no longer exist.
func (r *Registry) load() error {
	data, err := os.ReadFile(r.path())
	if err != nil {
		if os.IsNotExist(err) {
			return nil // fresh start
		}
		return err
	}

	var raw []Entry
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, e := range raw {
		if _, err := os.Stat(e.Path); err == nil {
			r.entries[e.Path] = e // only keep entries whose files still exist
		}
	}
	return nil
}

// loadDirs reads the remembered scan-directory map from disk.
func (r *Registry) loadDirs() error {
	data, err := os.ReadFile(r.dirsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &r.dirs)
}

// saveDirs writes the scan-directory map atomically. Caller must hold r.mu.
func (r *Registry) saveDirs() error {
	if err := os.MkdirAll(r.configDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(r.dirs, "", "  ")
	if err != nil {
		return err
	}
	tmp := r.dirsPath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, r.dirsPath())
}

// save writes the current entries atomically. Caller must hold r.mu.
func (r *Registry) save() error {
	if err := os.MkdirAll(r.configDir, 0755); err != nil {
		return err
	}

	list := make([]Entry, 0, len(r.entries))
	for _, e := range r.entries {
		list = append(list, e)
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	tmp := r.path() + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, r.path())
}

func collectEntries(idx scanner.ModelIndex) []scanner.AssetEntry {
	var out []scanner.AssetEntry
	out = append(out, idx.Models...)
	out = append(out, idx.VAEs...)
	out = append(out, idx.LoRAs...)
	out = append(out, idx.ControlNets...)
	out = append(out, idx.Embeddings...)
	out = append(out, idx.Upscalers...)
	return out
}
