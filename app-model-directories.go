package main

// app-model-directories.go
//
// Storage-level API for model directories. Every directory HOWL scans is
// recorded with exactly one kind (llm, tts, image), and that kind is what keeps
// the menus apart: a chat model in a chat directory can never appear in the
// image model list, and a diffusion model can never appear in the chat list.
//
// Adding a directory scans it immediately and remembers it. Removing one
// forgets the directory and its entries but never touches the files on disk.
//
// This layer is what the settings UI will drive: list, add, refresh, remove.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"howl-chat/internal/errors"
	"howl-chat/internal/image/registry"
	"howl-chat/internal/image/scanner"
)

// ModelDirectory describes one remembered model directory.
type ModelDirectory struct {
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	Exists bool   `json:"exists"`
	Models int    `json:"models"`
}

// GetModelDirectories lists every remembered directory with its kind, whether
// it is still on disk, and how many models are indexed under it.
func (a *App) GetModelDirectories() ([]ModelDirectory, error) {
	if a.modelRegistry == nil {
		return []ModelDirectory{}, fmt.Errorf("model registry not initialised")
	}
	dirs := a.modelRegistry.Directories("")
	out := make([]ModelDirectory, 0, len(dirs))
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		out = append(out, ModelDirectory{
			Path:   dir,
			Kind:   a.modelRegistry.KindForPath(dir),
			Exists: err == nil && info.IsDir(),
			Models: len(a.modelRegistry.EntriesUnderDir(dir)),
		})
	}
	return out, nil
}

// AddModelDirectory validates dir, scans it for the given kind, and remembers
// it. A directory can only have one kind: re-adding it as a different kind is
// refused rather than silently moving models between menus.
func (a *App) AddModelDirectory(dir, kind string) (ModelDirectory, error) {
	if a.modelRegistry == nil {
		return ModelDirectory{}, fmt.Errorf("model registry not initialised")
	}
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return ModelDirectory{}, fmt.Errorf("directory path is empty")
	}
	if !registry.ValidKind(kind) {
		return ModelDirectory{}, fmt.Errorf("unknown directory kind %q (want llm, tts or image)", kind)
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = filepath.Clean(dir)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return ModelDirectory{}, fmt.Errorf("cannot access %q: %w", abs, err)
	}
	if !info.IsDir() {
		return ModelDirectory{}, fmt.Errorf("%q is not a directory", abs)
	}
	if existing := a.modelRegistry.KindForPath(abs); existing != "" && existing != kind {
		return ModelDirectory{}, fmt.Errorf("%q is already added as %q; remove it before adding it as %q", abs, existing, kind)
	}

	count, scanErr := a.scanDirectoryIntoRegistry(abs, kind)
	if err := a.modelRegistry.RegisterDir(abs, kind); err != nil {
		return ModelDirectory{}, err
	}
	if scanErr != nil {
		a.logErrorSilent(errors.REGModelAddFail, errors.Medium, "registry.dir.add", "app-model-directories.go", scanErr)
		return ModelDirectory{}, scanErr
	}
	return ModelDirectory{Path: abs, Kind: kind, Exists: true, Models: count}, nil
}

// RemoveModelDirectory forgets a directory and every model indexed under it.
// Files on disk are never deleted, moved or modified.
func (a *App) RemoveModelDirectory(dir string) error {
	if a.modelRegistry == nil {
		return fmt.Errorf("model registry not initialised")
	}
	if strings.TrimSpace(dir) == "" {
		return fmt.Errorf("directory path is empty")
	}
	return a.modelRegistry.UnregisterDir(dir)
}

// RefreshModelDirectories rescans every remembered directory (or only one kind
// when kind is non-empty) so newly dropped-in files appear and deleted ones
// disappear. Missing directories are skipped, not forgotten: an unplugged drive
// keeps its place in the list.
func (a *App) RefreshModelDirectories(kind string) ([]ModelDirectory, error) {
	if a.modelRegistry == nil {
		return []ModelDirectory{}, fmt.Errorf("model registry not initialised")
	}
	if kind != "" && !registry.ValidKind(kind) {
		return []ModelDirectory{}, fmt.Errorf("unknown directory kind %q", kind)
	}
	for _, dir := range a.modelRegistry.Directories(kind) {
		if _, err := os.Stat(dir); err != nil {
			continue
		}
		if _, err := a.scanDirectoryIntoRegistry(dir, a.modelRegistry.KindForPath(dir)); err != nil {
			a.logErrorSilent(errors.REGModelListFail, errors.Medium, "registry.dir.refresh", "app-model-directories.go", err)
		}
	}
	return a.GetModelDirectories()
}

// scanDirectoryIntoRegistry walks one directory for one kind and merges what it
// finds. Image directories go through the image scanner so their assets keep
// their families and companion file slots; chat and voice directories only need
// extension matching plus the GGUF architecture check.
func (a *App) scanDirectoryIntoRegistry(dir, kind string) (int, error) {
	if a.modelRegistry == nil {
		return 0, fmt.Errorf("model registry not initialised")
	}
	switch kind {
	case registry.KindLLM, registry.KindTTS:
		exts := []string{".gguf"}
		if kind == registry.KindTTS {
			exts = []string{".onnx"}
		}
		entries, seen, err := scanDirForType(dir, kind, exts, a.modelRegistry.Lookup)
		if err != nil {
			return 0, err
		}
		added, removed, err := a.modelRegistry.SyncDir(dir, entries, seen)
		if err != nil {
			return 0, err
		}
		if added > 0 || removed > 0 {
			fmt.Println("[scan]", dir, "->", added, "new or updated,", removed, "removed")
		}
		return len(entries), nil

	case registry.KindImage:
		sc := scanner.New(dir)
		sc.SetKnown(func(path string) (scanner.KnownEntry, bool) {
			e, ok := a.modelRegistry.Lookup(path)
			if !ok {
				return scanner.KnownEntry{}, false
			}
			return scanner.KnownEntry{
				Path:        e.Path,
				Type:        scanner.AssetType(e.Type),
				Family:      scanner.ModelFamily(e.Family),
				Format:      e.Format,
				SizeMB:      e.SizeMB,
				Files:       e.Files,
				Recommended: e.Recommended,
				FileSize:    e.FileSize,
				ModTimeNS:   e.ModTimeNS,
			}, true
		})
		idx, err := sc.Scan()
		if err != nil {
			return 0, err
		}
		added, removed, err := a.modelRegistry.SyncIndex(dir, idx)
		if err != nil {
			return 0, err
		}
		if added > 0 || removed > 0 {
			fmt.Println("[scan-image]", dir, "->", added, "new or updated,", removed, "removed")
		}
		count := len(idx.Models) + len(idx.VAEs) + len(idx.LoRAs) +
			len(idx.ControlNets) + len(idx.Embeddings) + len(idx.Upscalers) +
			len(idx.T5XXLs) + len(idx.ClipLs) + len(idx.ClipGs)
		return count, nil
	}
	return 0, fmt.Errorf("unknown directory kind %q", kind)
}
