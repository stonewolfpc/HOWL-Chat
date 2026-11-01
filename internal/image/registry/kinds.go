package registry

import (
	"path/filepath"
	"sort"
	"strings"

	"howl-chat/internal/image/scanner"
)

// Directory kinds. A registered directory belongs to exactly one kind, and that
// kind is the boundary between menus: chat models, image models and voices are
// kept apart by where they live, not by what they are called.
const (
	KindLLM   = "llm"
	KindTTS   = "tts"
	KindImage = "image"
)

// ValidKind reports whether kind is a recognised directory kind.
func ValidKind(kind string) bool {
	switch kind {
	case KindLLM, KindTTS, KindImage:
		return true
	}
	return false
}

// Directories returns the registered directories of one kind, sorted so the UI
// and the diagnostic tools always walk them in the same order. An empty kind
// returns every directory.
func (r *Registry) Directories(kind string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.dirs))
	for dir, k := range r.dirs {
		if kind == "" || k == kind {
			out = append(out, dir)
		}
	}
	sort.Strings(out)
	return out
}

// IsUnderDir reports whether path is dir itself or lives beneath it.
//
// The separator check is what stops "E:\AI GGUF\GGUF2" from being treated as
// living inside "E:\AI GGUF\GGUF". Comparison is case-insensitive because HOWL
// ships for Windows.
func IsUnderDir(path, dir string) bool {
	if path == "" || dir == "" {
		return false
	}
	p := filepath.Clean(path)
	d := filepath.Clean(dir)
	if strings.EqualFold(p, d) {
		return true
	}
	if !strings.HasSuffix(d, string(filepath.Separator)) {
		d += string(filepath.Separator)
	}
	return strings.HasPrefix(strings.ToLower(p), strings.ToLower(d))
}

// KindForPath returns the kind of the deepest registered directory containing
// path, or "" when path is not inside any registered directory.
func (r *Registry) KindForPath(path string) string {
	if path == "" {
		return ""
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	bestKind := ""
	bestLen := -1
	for dir, kind := range r.dirs {
		if len(dir) > bestLen && IsUnderDir(path, dir) {
			bestKind = kind
			bestLen = len(dir)
		}
	}
	return bestKind
}

// EntriesUnderDir returns live entries whose path lives inside dir.
func (r *Registry) EntriesUnderDir(dir string) []Entry {
	var out []Entry
	for _, e := range r.All() {
		if IsUnderDir(e.Path, dir) {
			out = append(out, e)
		}
	}
	return out
}

// ByTypeInDirs returns live entries of assetType that live under a directory of
// the given kind. This is the menu boundary: an asset only appears where the
// kind of its directory says it belongs.
func (r *Registry) ByTypeInDirs(assetType, kind string) []Entry {
	var out []Entry
	for _, e := range r.All() {
		if e.Type != assetType {
			continue
		}
		if r.KindForPath(e.Path) != kind {
			continue
		}
		out = append(out, e)
	}
	return out
}

// ToModelIndexForKind converts the entries under one directory kind back into a
// scanner.ModelIndex, so image menus never see assets from chat directories.
func (r *Registry) ToModelIndexForKind(kind string) scanner.ModelIndex {
	var keep []Entry
	for _, e := range r.All() {
		if r.KindForPath(e.Path) == kind {
			keep = append(keep, e)
		}
	}
	return buildModelIndex(keep)
}
