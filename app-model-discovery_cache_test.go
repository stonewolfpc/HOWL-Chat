package main

// app-model-discovery_cache_test.go
//
// Discovery must not re-read a GGUF whose size and mtime still match the
// registry. These tests use files that are deliberately NOT valid GGUF: if a
// path parses the file, classification fails and the entry disappears, so a
// surviving entry proves the cache was used and the file was never read.

import (
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/backend/gguf"
	"howl-chat/internal/image/registry"
)

// cachedGarbageGGUF writes a file that cannot be parsed as GGUF and returns the
// registry entry the scanner would have recorded for it.
func cachedGarbageGGUF(t *testing.T, dir string) (string, registry.Entry) {
	t.Helper()
	path := filepath.Join(dir, "cached-chat.gguf")
	if err := os.WriteFile(path, []byte("this is not a gguf header"), 0o644); err != nil {
		t.Fatalf("write fake gguf: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat fake gguf: %v", err)
	}
	return path, registry.Entry{
		Name:           "cached-chat",
		Path:           path,
		Type:           registry.KindLLM,
		Family:         "llama",
		Format:         "gguf",
		SizeMB:         float64(info.Size()) / (1024 * 1024),
		FileSize:       info.Size(),
		ModTimeNS:      info.ModTime().UnixNano(),
		AdmissionClass: string(gguf.AdmissionChat),
		Arch:           "llama",
		AdmissionV:     registry.AdmissionCacheVersion,
	}
}

func knownOnly(entry registry.Entry) func(string) (registry.Entry, bool) {
	return func(path string) (registry.Entry, bool) {
		if path == entry.Path {
			return entry, true
		}
		return registry.Entry{}, false
	}
}

func TestScanDirForTypeReusesCachedAdmissionWithoutParsing(t *testing.T) {
	dir := t.TempDir()
	path, entry := cachedGarbageGGUF(t, dir)

	entries, seen, err := scanDirForType(dir, registry.KindLLM, []string{".gguf"}, knownOnly(entry))
	if err != nil {
		t.Fatalf("scanDirForType: %v", err)
	}
	if len(entries) != 1 || entries[0].Path != path {
		t.Fatalf("entries = %+v, want the cached entry", entries)
	}
	if entries[0].Family != "llama" || entries[0].AdmissionClass != string(gguf.AdmissionChat) {
		t.Fatalf("entry lost its classification: %+v", entries[0])
	}
	if _, ok := seen[path]; !ok {
		t.Fatal("cached entry was not marked seen")
	}
}

func TestScanDirForTypeReclassifiesChangedFile(t *testing.T) {
	dir := t.TempDir()
	_, entry := cachedGarbageGGUF(t, dir)
	entry.ModTimeNS-- // the file changed since it was recorded

	entries, _, err := scanDirForType(dir, registry.KindLLM, []string{".gguf"}, knownOnly(entry))
	if err != nil {
		t.Fatalf("scanDirForType: %v", err)
	}
	// The cache no longer applies, the file is parsed, and an unreadable
	// header is refused for chat.
	if len(entries) != 0 {
		t.Fatalf("entries = %+v, want a re-classification that refuses the file", entries)
	}
}

func TestScanDirForTypeReparsesEntryWhoseCacheIsStale(t *testing.T) {
	dir := t.TempDir()
	_, entry := cachedGarbageGGUF(t, dir)
	entry.AdmissionV = registry.AdmissionCacheVersion - 1

	entries, _, err := scanDirForType(dir, registry.KindLLM, []string{".gguf"}, knownOnly(entry))
	if err != nil {
		t.Fatalf("scanDirForType: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("entries = %+v, want the stale cache ignored", entries)
	}
}

func TestLLMInfosServesCachedAdmissionWithoutReadingTheFile(t *testing.T) {
	dir := t.TempDir()
	path, entry := cachedGarbageGGUF(t, dir)

	infos := llmInfos([]registry.Entry{entry})
	if len(infos) != 1 || infos[0].Path != path {
		t.Fatalf("picker entries = %+v, want the cached chat model", infos)
	}
}

func TestLLMInfosFallsBackWhenCacheIsStale(t *testing.T) {
	dir := t.TempDir()
	_, entry := cachedGarbageGGUF(t, dir)
	entry.AdmissionV = registry.AdmissionCacheVersion - 1

	if infos := llmInfos([]registry.Entry{entry}); len(infos) != 0 {
		t.Fatalf("picker entries = %+v, want stale cache ignored", infos)
	}
}

func TestRegistryPersistsAdmissionCache(t *testing.T) {
	dir := t.TempDir()
	_, entry := cachedGarbageGGUF(t, dir)
	entry.ModelTags = []string{"any-to-any"}
	entry.EmbeddingLen = 3072

	configDir := t.TempDir()
	first, err := registry.New(configDir)
	if err != nil {
		t.Fatalf("registry.New: %v", err)
	}
	if err := first.MergeEntries([]registry.Entry{entry}); err != nil {
		t.Fatalf("MergeEntries: %v", err)
	}

	reloaded, err := registry.New(configDir)
	if err != nil {
		t.Fatalf("reload registry: %v", err)
	}
	got, ok := reloaded.Lookup(entry.Path)
	if !ok {
		t.Fatal("entry did not survive the reload")
	}
	if got.AdmissionV != registry.AdmissionCacheVersion || got.Arch != "llama" ||
		got.AdmissionClass != string(gguf.AdmissionChat) || got.EmbeddingLen != 3072 ||
		len(got.ModelTags) != 1 || got.ModelTags[0] != "any-to-any" {
		t.Fatalf("cached admission did not survive the reload: %+v", got)
	}
}
