package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestScanReusesUnchangedFiles proves the incremental image scan: a recorded
// entry whose size and mtime still match the file is reused without being
// re-classified, a changed file is not reused, and a vanished file is dropped.
func TestScanReusesUnchangedFiles(t *testing.T) {
	dir := t.TempDir()
	modelPath := filepath.Join(dir, "test-model.safetensors")
	if err := os.WriteFile(modelPath, []byte("not a real model"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	info, err := os.Stat(modelPath)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}

	recorded := KnownEntry{
		Path:      modelPath,
		Type:      AssetDiffusion,
		Family:    FamilySD15,
		Format:    "safetensors",
		SizeMB:    0.0001,
		FileSize:  info.Size(),
		ModTimeNS: info.ModTime().UnixNano(),
	}

	sc := New(dir)
	sc.SetKnown(func(path string) (KnownEntry, bool) {
		if path != modelPath {
			return KnownEntry{}, false
		}
		return recorded, true
	})

	idx, err := sc.Scan()
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(idx.Models) != 1 {
		t.Fatalf("unchanged file should be reused, got %d models", len(idx.Models))
	}
	if idx.Models[0].Path != modelPath || idx.Models[0].Family != FamilySD15 {
		t.Fatalf("reused entry wrong: %+v", idx.Models[0])
	}

	// Same name, different bytes and mtime: the recorded identity no longer matches.
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(modelPath, []byte("different content now"), 0644); err != nil {
		t.Fatalf("rewrite: %v", err)
	}
	if err := os.Chtimes(modelPath, time.Now(), time.Now()); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	idx2, err := sc.Scan()
	if err != nil {
		t.Fatalf("scan 2: %v", err)
	}
	if len(idx2.Models) != 0 {
		t.Fatalf("changed file must not be reused, got %d models", len(idx2.Models))
	}

	// Deleted file: no longer reported at all.
	if err := os.Remove(modelPath); err != nil {
		t.Fatalf("remove: %v", err)
	}
	idx3, err := sc.Scan()
	if err != nil {
		t.Fatalf("scan 3: %v", err)
	}
	if len(idx3.Models) != 0 {
		t.Fatalf("deleted file still reported: %d models", len(idx3.Models))
	}
}
