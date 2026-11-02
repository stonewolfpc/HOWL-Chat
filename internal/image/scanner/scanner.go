package scanner

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Scanner walks a models directory, classifies every file, generates .howl.json
// manifests (once per file), and returns a fully populated ModelIndex.
type Scanner struct {
	modelsDir string
	store     *IndexStore
	known     KnownFunc
}

// New creates a Scanner for the given models directory.
func New(modelsDir string) *Scanner {
	return &Scanner{
		modelsDir: modelsDir,
		store:     NewIndexStore(modelsDir),
	}
}

// KnownEntry is a previously recorded asset handed back to the scanner when the
// file on disk is unchanged, so it never has to be re-read or re-classified.
type KnownEntry struct {
	Path        string
	Type        AssetType
	Family      ModelFamily
	Format      string
	SizeMB      float64
	Files       map[string]string
	Recommended Recommended
	FileSize    int64
	ModTimeNS   int64
}

// KnownFunc looks up a recorded asset by absolute path.
type KnownFunc func(path string) (KnownEntry, bool)

// SetKnown installs the lookup used to skip unchanged files during a scan.
func (s *Scanner) SetKnown(f KnownFunc) { s.known = f }

func (k KnownEntry) toAssetEntry() AssetEntry {
	name := strings.TrimSuffix(filepath.Base(k.Path), filepath.Ext(k.Path))
	return AssetEntry{
		Name:        name,
		Path:        k.Path,
		Type:        k.Type,
		Family:      k.Family,
		Format:      k.Format,
		SizeMB:      k.SizeMB,
		Files:       k.Files,
		Recommended: k.Recommended,
		FileSize:    k.FileSize,
		ModTimeNS:   k.ModTimeNS,
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

		// Unchanged files are never re-classified: a recorded entry whose size
		// and mtime still match this file is reused as-is.
		if s.known != nil {
			if fi, statErr := d.Info(); statErr == nil && fi.Size() > 0 {
				if k, ok := s.known(path); ok && k.FileSize == fi.Size() && k.ModTimeNS == fi.ModTime().UnixNano() {
					appendEntry(&idx, k.toAssetEntry())
					return nil
				}
			}
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
	fmt.Printf("[Scanner] Scan complete: %d models, %d VAEs, %d LoRAs, %d ControlNets, %d embeddings, %d LLMs\n",
		len(idx.Models), len(idx.VAEs), len(idx.LoRAs), len(idx.ControlNets), len(idx.Embeddings), len(idx.LLMs))
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

// largeCheckpointMB is the size above which an unnamed single-file checkpoint is
// SDXL-class. SD1.5 checkpoints stay under ~4 GB even in fp32.
const largeCheckpointMB = 5000.0

// upgradeLargeCheckpoint promotes a name-less SD1.5 classification to SDXL for
// SDXL-scale weights. Without it a 6.6 GB SDXL merge is labelled SD1.5, which
// picks a 512x512 default and the wrong sampler.
func upgradeLargeCheckpoint(family ModelFamily, sizeMB float64) ModelFamily {
	if family == FamilySD15 && sizeMB >= largeCheckpointMB {
		return FamilySDXL
	}
	return family
}

func (s *Scanner) scanFile(path string, d os.DirEntry) *AssetEntry {
	info, err := d.Info()
	if err != nil {
		return nil
	}
	sizeMB := float64(info.Size()) / (1024 * 1024)
	name := d.Name()
	stem := strings.TrimSuffix(name, filepath.Ext(name))
	// Read the weights first: an Invoke AI export is called model.safetensors
	// inside a UUID folder, so the filename knows nothing about it.
	probe := ProbeWeight(path)
	nameType := ClassifyName(name)
	probeType, probeKnowsType := probe.AssetType()
	// For these two verdicts the contents outrank the filename. A language
	// model dropped in an image folder was listed as an SD1.5 checkpoint, and
	// a VAE whose name said nothing was offered as a diffusion model; both
	// only failed once the engine tried to load them. Family selection has
	// always trusted the weights first, and the asset type now agrees.
	if probeKnowsType && (probeType == AssetLLM || probeType == AssetVAE) {
		nameType = probeType
	}
	if nameType == AssetUnknown && probeKnowsType {
		nameType = probeType
	}
	if nameType == AssetUnknown {
		nameType = ClassifySafetensorsHeader(path)
	}

	// Reuse the sidecar only while it still describes this exact file. A model
	// replaced under the same name has a new size or mtime, so its stale
	// classification is discarded instead of being shown as truth.
	if m := LoadManifest(path); m != nil && m.Logic == manifestLogicVersion {
		fresh := m.FileSize == info.Size() && m.ModTimeNS == info.ModTime().UnixNano() && m.FileSize != 0
		if fresh && (nameType == AssetUnknown || m.Type == nameType) {
			return manifestToEntry(m, path)
		}
	}

	assetType := ClassifySize(nameType, sizeMB)

	var family ModelFamily
	switch assetType {
	case AssetDiffusion:
		// The weights are the authority: a file can be renamed, but its tensor
		// layout still says Flux or SDXL. The filename is only a fallback for
		// formats whose contents we cannot read.
		if f, ok := probe.Family(); ok {
			family = f
		} else {
			family = GuessFamily(stem)
			if family == FamilyUnknown {
				return nil // not a diffusion model we can handle
			}
			family = upgradeLargeCheckpoint(family, sizeMB)
		}
	case AssetLoRA, AssetControlNet, AssetVAE:
		if f, ok := probe.Family(); ok {
			family = f
		} else {
			family = GuessFamily(stem) // best-effort; FamilyUnknown is acceptable
		}
	default:
		family = FamilyUnknown
	}

	rec := buildRecommended(assetType, family)
	m := &Manifest{
		Name:        DisplayName(stem, filepath.Dir(path), probe, family),
		Logic:       manifestLogicVersion,
		Type:        assetType,
		Family:      family,
		Format:      strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), "."),
		SizeMB:      sizeMB,
		FileSize:    info.Size(),
		ModTimeNS:   info.ModTime().UnixNano(),
		Files:       map[string]string{"model": path},
		Recommended: rec,
	}
	_ = SaveManifest(path, m)
	return manifestToEntry(m, path)
}

// ClassifySafetensorsHeader catches large LoRAs and companion files whose
// filenames are ambiguous enough for size-based classification to get wrong.
func ClassifySafetensorsHeader(path string) AssetType {
	if strings.ToLower(filepath.Ext(path)) != ".safetensors" {
		return AssetUnknown
	}
	f, err := os.Open(path)
	if err != nil {
		return AssetUnknown
	}
	defer f.Close()

	var headerLen uint64
	if err := binary.Read(f, binary.LittleEndian, &headerLen); err != nil {
		return AssetUnknown
	}
	if headerLen == 0 {
		return AssetUnknown
	}
	if headerLen > 1<<20 {
		headerLen = 1 << 20
	}
	data := make([]byte, headerLen)
	if _, err := io.ReadFull(f, data); err != nil {
		return AssetUnknown
	}
	header := strings.ToLower(string(data))
	switch {
	case strings.Contains(header, "lora_unet") || strings.Contains(header, "lora_te") || strings.Contains(header, "lycoris"):
		return AssetLoRA
	case strings.Contains(header, "first_stage_model.decoder") && !strings.Contains(header, "model.diffusion_model"):
		return AssetVAE
	default:
		return AssetUnknown
	}
}

// ---- directory bundle classification ---------------------------------------

func (s *Scanner) scanDirBundle(dir string) (AssetEntry, bool) {
	// A Diffusers pipeline is one model made of component folders, not a pile
	// of loose unet/vae/text_encoder fragments.
	if pipeline := loadModelIndexClass(dir); pipeline != "" {
		return s.scanDiffusersBundle(dir, pipeline)
	}

	// One weight and no companions is a single model that happens to live in a
	// folder. Calling it a bundle hid its real format and named it after the
	// folder, which for Invoke is a UUID. Let the walker classify the file.
	if !hasBundleShape(dir) {
		return AssetEntry{}, false
	}
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

	// A bundle's model slot must hold the checkpoint itself: sidecars such as
	// TAESD decoders sort earlier alphabetically and would otherwise take it.
	var modelSize int64
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

		// Classify by content-independent name rules so bundled video packs
		// (umt5-*-encoder, wan_*_vae, *-T2V-*.gguf) land in the right slot.
		switch ClassifyName(f.Name()) {
		case AssetVAE:
			fileMap["vae"] = fp
		case AssetT5XXL:
			fileMap["t5xxl"] = fp
		case AssetClipL:
			fileMap["clip_l"] = fp
		case AssetClipG:
			fileMap["clip_g"] = fp
		case AssetLLM:
			fileMap["llm"] = fp
		case AssetControlNet:
			fileMap["control_net"] = fp
		case AssetLoRA:
			// LoRAs are indexed on their own, never as bundle components.
		case AssetDiffusion:
			if _, exists := fileMap["model"]; !exists || info.Size() > modelSize {
				fileMap["model"] = fp
				modelSize = info.Size()
			}
		default:
			switch {
			case strings.Contains(lower, "taesd"):
				// TAESD/TAESDXL decoders drive step previews, never the model slot.
				fileMap["taesd"] = fp
			case strings.HasPrefix(lower, "diffusion_model"),
				strings.HasPrefix(lower, "model") && strings.Contains(lower, "diffusion"):
				fileMap["diffusion_model"] = fp
			default:
				if _, exists := fileMap["model"]; !exists || info.Size() > modelSize {
					fileMap["model"] = fp
					modelSize = info.Size()
				}
			}
		}

	}

	if len(fileMap) == 0 {
		return AssetEntry{}, false
	}

	// Determine family from the component files, then the model filename.
	dirName := filepath.Base(dir)
	modelBase := dirName
	if p := fileMap["model"]; p != "" {
		modelBase = strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
	}
	guessed := GuessFamily(modelBase)
	var family ModelFamily
	switch {
	case hasKey(fileMap, "clip_g") && !hasKey(fileMap, "t5xxl"):
		family = FamilySDXL
	case guessed == FamilyWan || guessed == FamilyLTX || guessed == FamilyHunyuan || guessed == FamilyAnimate:
		family = guessed
	case hasKey(fileMap, "diffusion_model") && hasKey(fileMap, "t5xxl"):
		family = FamilyFlux
	case guessed != FamilyUnknown:
		family = guessed
	default:
		family = FamilySD15
	}

	// Same rule as single-file checkpoints: a multi-GB bundle whose model name
	// carries no architecture hint is SDXL, not SD1.5.
	if modelPath := fileMap["model"]; modelPath != "" {
		if fi, err := os.Stat(modelPath); err == nil {
			family = upgradeLargeCheckpoint(family, float64(fi.Size())/(1024*1024))
		}
	}

	// A bundle is just a directory, so a Flux VAE sitting beside an SDXL
	// checkpoint would otherwise be attached to it and render garbage.
	if vaePath := fileMap["vae"]; vaePath != "" {
		vaeFamily := GuessFamily(strings.TrimSuffix(filepath.Base(vaePath), filepath.Ext(vaePath)))
		if !CompatibleVAE(family, vaeFamily) {
			delete(fileMap, "vae")
		}
	}

	// The bundle's own weights override component-name guesses.
	probe := ProbeWeight(fileMap["model"])
	if f, ok := probe.Family(); ok {
		family = f
	}

	rec := buildRecommended(AssetDiffusion, family)
	m := &Manifest{
		Name:        DisplayName(dirName, dir, probe, family),
		Logic:       manifestLogicVersion,
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
	case AssetT5XXL:
		idx.T5XXLs = append(idx.T5XXLs, e)
	case AssetClipL:
		idx.ClipLs = append(idx.ClipLs, e)
	case AssetClipG:
		idx.ClipGs = append(idx.ClipGs, e)
	case AssetLLM:
		idx.LLMs = append(idx.LLMs, e)
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
		FileSize:    m.FileSize,
		ModTimeNS:   m.ModTimeNS,
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
		Steps:     DefaultSteps(family),
		CFG:       DefaultCFG(family),
		Width:     res,
		Height:    res,
	}
}

func hasKey(m map[string]string, k string) bool {
	_, ok := m[k]
	return ok
}
