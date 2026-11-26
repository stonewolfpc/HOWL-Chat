// Package installer registers scanned model assets into the HOWL model index.
//
// "Install" means: read the file's metadata, write a .howl.json sidecar
// manifest so future scans are instant, and persist the index — nothing more.
//
// Files are NEVER moved, copied, or renamed. The user's folder structure is
// left completely untouched. A model at D:\AI GGUF\flux1-dev.gguf stays
// exactly there; installing it simply means HOWL has catalogued it and knows
// its type, family, and generation defaults.
package installer

import (
	"fmt"

	"howl-chat/internal/image/scanner"
)

// InstalledAsset records the outcome of registering a single asset.
type InstalledAsset struct {
	Name          string  `json:"name"`
	OriginalPath  string  `json:"original_path"`
	InstalledPath string  `json:"installed_path"` // same as OriginalPath — files never move
	Type          string  `json:"type"`
	Family        string  `json:"family"`
	Format        string  `json:"format"`
	SizeMB        float64 `json:"size_mb"`
	Error         string  `json:"error,omitempty"`
}

// InstallResult is the aggregate outcome of a registration pass.
type InstallResult struct {
	Installed      []InstalledAsset `json:"installed"`
	Skipped        []InstalledAsset `json:"skipped"` // manifest already existed
	Failed         []InstalledAsset `json:"failed"`
	TotalCount     int              `json:"total_count"`
	InstalledCount int              `json:"installed_count"`
	SkippedCount   int              `json:"skipped_count"`
	FailedCount    int              `json:"failed_count"`
}

// Installer registers assets from a scan result into the model index.
type Installer struct {
	modelsDir string
}

// New creates an Installer for the given models root directory.
func New(modelsDir string, _ bool) *Installer {
	return &Installer{modelsDir: modelsDir}
}

// Install iterates every entry in idx, writes a .howl.json sidecar for each
// file that does not already have one, and saves the updated index.
// It never touches the actual model files.
func (inst *Installer) Install(idx scanner.ModelIndex) (InstallResult, error) {
	var result InstallResult
	all := collectAllEntries(idx)
	result.TotalCount = len(all)

	for _, e := range all {
		ia := registerEntry(e)
		if ia.Error != "" {
			result.Failed = append(result.Failed, ia)
			result.FailedCount++
		} else if ia.InstalledPath == "" {
			// manifest already existed — skipped
			ia.InstalledPath = ia.OriginalPath
			result.Skipped = append(result.Skipped, ia)
			result.SkippedCount++
		} else {
			result.Installed = append(result.Installed, ia)
			result.InstalledCount++
		}
	}

	// Persist the index so the next LoadIndex call returns instantly.
	store := scanner.NewIndexStore(inst.modelsDir)
	if err := store.Save(idx); err != nil {
		return result, fmt.Errorf("installer: failed to save index: %w", err)
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// internal helpers
// ---------------------------------------------------------------------------

// registerEntry ensures the asset has a .howl.json sidecar manifest.
// Returns ia with InstalledPath="" when the manifest already existed (skip),
// or InstalledPath=OriginalPath when a new manifest was written (installed).
func registerEntry(e scanner.AssetEntry) InstalledAsset {
	ia := InstalledAsset{
		Name:         e.Name,
		OriginalPath: e.Path,
		Type:         string(e.Type),
		Family:       string(e.Family),
		Format:       e.Format,
		SizeMB:       e.SizeMB,
	}

	// If a sidecar already exists the asset is already registered.
	if existing := scanner.LoadManifest(e.Path); existing != nil {
		return ia // InstalledPath intentionally empty → caller counts as skipped
	}

	// Write a fresh manifest so future scans never re-classify this file.
	m := &scanner.Manifest{
		Name:        e.Name,
		Type:        e.Type,
		Family:      e.Family,
		Format:      e.Format,
		SizeMB:      e.SizeMB,
		Files:       e.Files,
		Recommended: e.Recommended,
	}
	if err := scanner.SaveManifest(e.Path, m); err != nil {
		ia.Error = err.Error()
		return ia
	}

	ia.InstalledPath = e.Path // file is untouched; path is identical
	return ia
}

func collectAllEntries(idx scanner.ModelIndex) []scanner.AssetEntry {
	var out []scanner.AssetEntry
	out = append(out, idx.Models...)
	out = append(out, idx.VAEs...)
	out = append(out, idx.LoRAs...)
	out = append(out, idx.ControlNets...)
	out = append(out, idx.IPAdapters...)
	out = append(out, idx.ClipVisions...)
	out = append(out, idx.Embeddings...)
	out = append(out, idx.Upscalers...)
	return out
}
