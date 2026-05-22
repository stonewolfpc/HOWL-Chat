package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Scanner walks a models directory, classifies every file, generates .howl.json
// manifests (once per file), and returns a fully populated ModelIndex.
type Scanner struct {
	modelsDir string
	store     *IndexStore
}

// New creates a Scanner for the given models directory.
func New(modelsDir string) *Scanner {
	return &Scanner{
		modelsDir: modelsDir,
		store:     NewIndexStore(modelsDir),
	}
}

// Scan performs a full directory walk. It:
//  1. Checks each file's manifest — reuses it if present and up-to-date.
//  2. Classifies new files and writes manifests.
//  3. Handles directory bundles (FLUX-style multi-file models).
//  4. Saves and returns the updated ModelIndex.
func (s *Scanner) Scan() (ModelIndex, error) {
	if err := os.MkdirAll(s.modelsDir, 0755); err != nil {
		return ModelIndex{}, fmt.Errorf("scanner: cannot access models dir: %w", err)
	}

	var idx ModelIndex
	dirBundles := map[string]bool{} // track dirs already handled as bundles

	// Walk the entire tree
	err := filepath.WalkDir(s.modelsDir, func(path string, d os.DirEntry, werr error) error {
		if werr != nil {
			return nil // skip unreadable entries
		}

		// Skip hidden files and manifest sidecars
		base := d.Name()
		if strings.HasPrefix(base, ".") || strings.HasSuffix(base, manifestSuffix) || base == indexFile {
			return nil
		}

		// If this is a sub-directory (not the root), try to treat it as a model bundle
		if d.IsDir() && path != s.modelsDir {
			entry, ok := s.scanDirBundle(path)
			if ok {
				dirBundles[path] = true
				appendEntry(&idx, entry)
				return filepath.SkipDir // don't recurse into handled bundles
			}
			return nil // recurse into unrecognised dirs
		}

		if d.IsDir() {
			return nil
		}

		// Skip files inside already-handled bundles
		parent := filepath.Dir(path)
		if dirBundles[parent] {
			return nil
		}

		if !IsModelFile(base) {
			return nil
		}

		entry := s.scanFile(path, d)
		if entry != nil {
			appendEntry(&idx, *entry)
		}
		return nil
	})
	if err != nil {
		return idx, fmt.Errorf("scanner: walk failed: %w", err)
	}

	_ = s.store.Save(idx)
	fmt.Printf("[Scanner] Scan complete: %d models, %d VAEs, %d LoRAs, %d ControlNets, %d embeddings\n",
		len(idx.Models), len(idx.VAEs), len(idx.LoRAs), len(idx.ControlNets), len(idx.Embeddings))
	return idx, nil
}

// LoadIndex returns the cached index without rescanning.
// Falls back to a fresh Scan if no cache exists.
func (s *Scanner) LoadIndex() (ModelIndex, error) {
	idx, err := s.store.Load()
	if err != nil {
		return s.Scan()
	}
	if len(idx.Models)+len(idx.VAEs)+len(idx.LoRAs) == 0 {
		return s.Scan()
	}
	return idx, nil
}

// ---- file classification ---------------------------------------------------

func (s *Scanner) scanFile(path string, d os.DirEntry) *AssetEntry {
	info, err := d.Info()
	if err != nil {
		return nil
	}
	sizeMB := float64(info.Size()) / (1024 * 1024)
	name := d.Name()
	stem := strings.TrimSuffix(name, filepath.Ext(name))

	// Reuse existing manifest if present
	if m := LoadManifest(path); m != nil {
		return manifestToEntry(m, path)
	}

	nameType := ClassifyName(name)
	assetType := ClassifySize(nameType, sizeMB)

	var family ModelFamily
	switch assetType {
	case AssetDiffusion:
		family = GuessFamily(stem)
		if family == FamilyUnknown {
			return nil // not a diffusion model we can handle
		}
	case AssetLoRA, AssetControlNet, AssetVAE:
		family = GuessFamily(stem) // best-effort; FamilyUnknown is acceptable
	default:
		family = FamilyUnknown
	}

	rec := buildRecommended(assetType, family)
	m := &Manifest{
		Name:        stem,
		Type:        assetType,
		Family:      family,
		Format:      strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), "."),
		SizeMB:      sizeMB,
		Files:       map[string]string{"model": path},
		Recommended: rec,
	}
	_ = SaveManifest(path, m)
	return manifestToEntry(m, path)
}

// ---- directory bundle classification ---------------------------------------

func (s *Scanner) scanDirBundle(dir string) (AssetEntry, bool) {
	// Reuse existing manifest if present
	if m := LoadDirManifest(dir); m != nil {
		e := manifestToEntry(m, dir)
		if e != nil {
			return *e, true
		}
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return AssetEntry{}, false
	}

	fileMap := make(map[string]string)
	var totalMB float64

	for _, f := range files {
		if f.IsDir() || !IsModelFile(f.Name()) {
			continue
		}
		info, err := f.Info()
		if err != nil {
			continue
		}
		totalMB += float64(info.Size()) / (1024 * 1024)
		fp := filepath.Join(dir, f.Name())
		lower := strings.ToLower(f.Name())

		switch {
		case strings.HasPrefix(lower, "diffusion_model"), strings.HasPrefix(lower, "model") && strings.Contains(lower, "diffusion"):
			fileMap["diffusion_model"] = fp
		case strings.HasPrefix(lower, "clip_l"), strings.HasPrefix(lower, "clip.l"):
			fileMap["clip_l"] = fp
		case strings.HasPrefix(lower, "clip_g"), strings.HasPrefix(lower, "clip.g"):
			fileMap["clip_g"] = fp
		case strings.HasPrefix(lower, "t5xxl"), strings.HasPrefix(lower, "t5-xxl"), strings.HasPrefix(lower, "t5_xxl"):
			fileMap["t5xxl"] = fp
		case strings.HasPrefix(lower, "vae"):
			fileMap["vae"] = fp
		default:
			if _, exists := fileMap["model"]; !exists {
				fileMap["model"] = fp
			}
		}
	}

	if len(fileMap) == 0 {
		return AssetEntry{}, false
	}

	// Determine family from presence of component files
	var family ModelFamily
	dirName := filepath.Base(dir)
	switch {
	case hasKey(fileMap, "diffusion_model") && hasKey(fileMap, "t5xxl"):
		family = FamilyFlux
	case hasKey(fileMap, "clip_g") && !hasKey(fileMap, "t5xxl"):
		family = FamilySDXL
	default:
		family = GuessFamily(dirName)
		if family == FamilyUnknown {
			family = FamilySD15
		}
	}

	rec := buildRecommended(AssetDiffusion, family)
	m := &Manifest{
		Name:        dirName,
		Type:        AssetDiffusion,
		Family:      family,
		Format:      "bundle",
		SizeMB:      totalMB,
		Files:       fileMap,
		Recommended: rec,
	}
	_ = SaveDirManifest(dir, m)
	entry := manifestToEntry(m, dir)
	if entry == nil {
		return AssetEntry{}, false
	}
	return *entry, true
}

// ---- helpers ---------------------------------------------------------------

func appendEntry(idx *ModelIndex, e AssetEntry) {
	switch e.Type {
	case AssetDiffusion:
		idx.Models = append(idx.Models, e)
	case AssetVAE:
		idx.VAEs = append(idx.VAEs, e)
	case AssetLoRA:
		idx.LoRAs = append(idx.LoRAs, e)
	case AssetControlNet:
		idx.ControlNets = append(idx.ControlNets, e)
	case AssetEmbedding:
		idx.Embeddings = append(idx.Embeddings, e)
	case AssetUpscaler:
		idx.Upscalers = append(idx.Upscalers, e)
	}
}

func manifestToEntry(m *Manifest, path string) *AssetEntry {
	if m == nil {
		return nil
	}
	return &AssetEntry{
		Name:        m.Name,
		Path:        path,
		Type:        m.Type,
		Family:      m.Family,
		Format:      m.Format,
		SizeMB:      m.SizeMB,
		Files:       m.Files,
		Recommended: m.Recommended,
	}
}

func buildRecommended(t AssetType, family ModelFamily) Recommended {
	if t != AssetDiffusion {
		return Recommended{}
	}
	res := DefaultResolution(family)
	return Recommended{
		Sampler:   DefaultSampler(family),
		Scheduler: "karras",
		Steps:     20,
		CFG:       7.0,
		Width:     res,
		Height:    res,
	}
}

func hasKey(m map[string]string, k string) bool {
	_, ok := m[k]
	return ok
}
