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

// AdmissionCacheVersion versions the cached GGUF admission facts stored on an
// entry. Bump it whenever AdmitModel's rules change, so entries classified by
// an older rule set are re-read exactly once instead of trusting stale facts.
const AdmissionCacheVersion = 1

// Entry is one record in the global registry.
type Entry struct {
	Name         string              `json:"name"`
	Path         string              `json:"path"`
	Type         string              `json:"type"`
	Family       string              `json:"family"`
	Format       string              `json:"format"`
	SizeMB       float64             `json:"size_mb"`
	FileSize     int64               `json:"file_size,omitempty"`
	ModTimeNS    int64               `json:"mod_time_ns,omitempty"`
	Files        map[string]string   `json:"files,omitempty"`
	Recommended  scanner.Recommended `json:"recommended"`
	RegisteredAt time.Time           `json:"registered_at"`

	// Cached GGUF admission facts, recorded once when the file was classified.
	// Picker reads and later scans reuse them instead of parsing the same
	// multi-gigabyte header again; AdmissionV marks which classifier rules
	// produced them.
	AdmissionClass string   `json:"admission,omitempty"`
	Arch           string   `json:"arch,omitempty"`
	ModelName      string   `json:"model_name,omitempty"`
	ModelTags      []string `json:"model_tags,omitempty"`
	EmbeddingLen   int      `json:"embedding_length,omitempty"`
	AdmissionV     int      `json:"admission_version,omitempty"`
}

// Registry is the in-memory + on-disk store.
type Registry struct {
	mu        sync.RWMutex
	configDir string
	entries   map[string]Entry  // keyed on absolute path
	dirs      map[string]string // dir path → asset type ("llm", "tts", "image")
	// needsSave is set when load() dropped entries, so the healing is written
	// back instead of being repeated (and re-injected) on every start.
	needsSave bool
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
	if r.needsSave {
		if err := r.save(); err != nil {
			return nil, err
		}
		r.needsSave = false
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

// Lookup returns the registered entry for an absolute path, if present.
func (r *Registry) Lookup(path string) (Entry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[path]
	return e, ok
}

// SyncDir merges the entries just scanned out of dir and drops any registered
// entry under that directory whose file was not seen on disk. That makes a scan
// a true directory sync: new files are added, vanished files are removed, and
// unchanged files are reused without re-reading their GGUF headers.
func (r *Registry) SyncDir(dir string, entries []Entry, seen map[string]struct{}) (int, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	added := 0
	for _, e := range entries {
		if e.RegisteredAt.IsZero() {
			e.RegisteredAt = time.Now()
		}
		if prev, ok := r.entries[e.Path]; !ok || !sameIdentity(prev, e) {
			added++
		}
		r.entries[e.Path] = e
	}

	removed := 0
	for path := range r.entries {
		if !IsUnderDir(path, dir) {
			continue
		}
		if _, ok := seen[path]; ok {
			continue
		}
		// Only drop entries whose file is actually gone. A file the scanner
		// chose not to classify (unknown family, unreadable header) is still on
		// disk and must stay listed rather than silently disappearing.
		if info, err := os.Stat(path); err == nil {
			// If the path is a directory, it was a bundle entry. Prune it if the
			// scanner no longer produces an entry for this exact path (e.g. the
			// bundle was split into individual ControlNet files).
			if info.IsDir() {
				// Check if any scanned entry has this exact path (a bundle that
				// was NOT split). If not, prune the stale directory entry.
				if _, ok := seen[path]; !ok {
					delete(r.entries, path)
					removed++
				}
				continue
			}
			// The file still exists but the scanner produced no entry for it.
			// If its identity changed, it was re-classified out (a chat model
			// replaced by a projector/component, or by an unreadable file).
			// Keeping the old entry would leave it advertised with a stale
			// classification, so the picker would still offer it. A file whose
			// identity is unchanged (locked, temporarily unreadable) keeps its
			// place.
			if prev, ok := r.entries[path]; ok && prev.FileSize != 0 &&
				(prev.FileSize != info.Size() || prev.ModTimeNS != info.ModTime().UnixNano()) {
				delete(r.entries, path)
				removed++
			}
			continue
		}
		delete(r.entries, path)
		removed++
	}

	if err := r.save(); err != nil {
		return added, removed, err
	}
	return added, removed, nil
}

// sameIdentity reports whether two entries describe the same unchanged file.
func sameIdentity(a, b Entry) bool {
	return a.FileSize != 0 && a.FileSize == b.FileSize &&
		a.ModTimeNS == b.ModTimeNS && a.Type == b.Type && a.Family == b.Family
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

// UnregisterDir removes a directory from the remembered scan locations and
// also removes all registry entries whose path starts with that directory.
func (r *Registry) UnregisterDir(dir string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Match regardless of case or trailing separators so a directory added by
	// the picker can always be removed by path. Only the registry forgets it —
	// the files on disk are never touched.
	for registered := range r.dirs {
		if IsUnderDir(registered, dir) && IsUnderDir(dir, registered) {
			delete(r.dirs, registered)
		}
	}
	for path := range r.entries {
		if IsUnderDir(path, dir) {
			delete(r.entries, path)
		}
	}
	if err := r.saveDirs(); err != nil {
		return err
	}
	return r.save()
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
			FileSize:     e.FileSize,
			ModTimeNS:    e.ModTimeNS,
			Files:        e.Files,
			Recommended:  e.Recommended,
			RegisteredAt: time.Now(),
		}
	}
	return r.save()
}

// SyncIndex merges a scanned image index for dir and drops registry entries
// under that directory whose files are gone from disk, mirroring the behaviour
// of a chat-model directory sync.
func (r *Registry) SyncIndex(dir string, idx scanner.ModelIndex) (int, int, error) {
	entries := make([]Entry, 0)
	seen := make(map[string]struct{})
	for _, ae := range collectEntries(idx) {
		entries = append(entries, Entry{
			Name:         ae.Name,
			Path:         ae.Path,
			Type:         string(ae.Type),
			Family:       string(ae.Family),
			Format:       ae.Format,
			SizeMB:       ae.SizeMB,
			FileSize:     ae.FileSize,
			ModTimeNS:    ae.ModTimeNS,
			Files:        ae.Files,
			Recommended:  ae.Recommended,
			RegisteredAt: time.Now(),
		})
		seen[ae.Path] = struct{}{}
	}
	return r.SyncDir(dir, entries, seen)
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
	return buildModelIndex(r.All())
}

// buildModelIndex groups registry entries into the scanner's index buckets.
func buildModelIndex(entries []Entry) scanner.ModelIndex {
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
		case scanner.AssetT5XXL:
			idx.T5XXLs = append(idx.T5XXLs, ae)
		case scanner.AssetClipL:
			idx.ClipLs = append(idx.ClipLs, ae)
		case scanner.AssetClipG:
			idx.ClipGs = append(idx.ClipGs, ae)
		case scanner.AssetLLM:
			idx.LLMs = append(idx.LLMs, ae)
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
		if _, err := os.Stat(e.Path); err != nil {
			r.needsSave = true // file is gone
			continue
		}
		// Kind validation deliberately does NOT run here. It used to call
		// entryKindMatches, which parses the GGUF header of every registered
		// model on every launch (measured 49s for 27 models) on the goroutine
		// that must finish before the window is created. Mis-filed entries are
		// re-filed by the background rescan instead, so startup pays only the
		// os.Stat above.
		r.entries[e.Path] = e
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
	out = append(out, idx.T5XXLs...)
	out = append(out, idx.ClipLs...)
	out = append(out, idx.ClipGs...)
	// Language models belong here too: a file can be reclassified from
	// checkpoint to text encoder, and the merge must rewrite its stored
	// type or the image dropdown keeps listing it forever.
	out = append(out, idx.LLMs...)
	return out
}
