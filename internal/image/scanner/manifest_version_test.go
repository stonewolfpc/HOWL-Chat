package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

// Files scanned by the previous classifier already have sidecars on disk. The
// user's Invoke AI folders held 81 of them, all labelled sd1.5. A sidecar whose
// logic version is missing must be re-read rather than trusted, or the fix
// never reaches files that were already scanned.
func TestStaleManifestIsReclassified(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "model.safetensors")
	writeSafetensors(t, path, nil, map[string][]int{
		"double_blocks.0.img_attn.proj.weight": {3072, 3072},
		"single_blocks.0.linear1.weight":       {12288, 3072},
	})

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}

	// Exactly what the old classifier wrote: no Logic field, family sd1.5.
	stale := &Manifest{
		Name:        "model",
		Type:        AssetDiffusion,
		Family:      FamilySD15,
		Format:      "safetensors",
		SizeMB:      1,
		FileSize:    info.Size(),
		ModTimeNS:   info.ModTime().UnixNano(),
		Files:       map[string]string{"model": path},
		Recommended: buildRecommended(AssetDiffusion, FamilySD15),
	}
	if err := SaveManifest(path, stale); err != nil {
		t.Fatalf("save stale manifest: %v", err)
	}

	idx, err := New(dir).Scan()
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(idx.Models) != 1 {
		t.Fatalf("got %d models, want 1", len(idx.Models))
	}
	if got := idx.Models[0].Family; got != FamilyFlux {
		t.Fatalf("stale sidecar was trusted: family = %q, want %q", got, FamilyFlux)
	}

	// And the reclassification is written back, carrying the current version.
	fresh := LoadManifest(path)
	if fresh == nil {
		t.Fatal("manifest missing after rescan")
	}
	if fresh.Logic != manifestLogicVersion {
		t.Fatalf("rewritten manifest logic = %d, want %d", fresh.Logic, manifestLogicVersion)
	}
	if fresh.Family != FamilyFlux {
		t.Fatalf("rewritten manifest family = %q, want %q", fresh.Family, FamilyFlux)
	}
}
