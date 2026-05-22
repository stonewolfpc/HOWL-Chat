package main

import (
	"fmt"

	"howl-chat/internal/errors"
	"howl-chat/internal/image/installer"
	"howl-chat/internal/image/registry"
	"howl-chat/internal/image/scanner"
)

// ScanResult is the full model index returned to the frontend.
// It mirrors scanner.ModelIndex but uses concrete types so Wails serialises it correctly.
type ScanResult struct {
	Models      []AssetInfo `json:"models"`
	VAEs        []AssetInfo `json:"vaes"`
	LoRAs       []AssetInfo `json:"loras"`
	ControlNets []AssetInfo `json:"controlnets"`
	Embeddings  []AssetInfo `json:"embeddings"`
	Upscalers   []AssetInfo `json:"upscalers"`
}

// AssetInfo is the frontend-visible description of a scanned asset.
type AssetInfo struct {
	Name          string            `json:"name"`
	Path          string            `json:"path"`
	Type          string            `json:"type"`
	Family        string            `json:"family"`
	Format        string            `json:"format"`
	SizeMB        float64           `json:"size_mb"`
	Compatibility string            `json:"compatibility"` // "verified", "supported", "unsupported", "unknown"
	Files         map[string]string `json:"files,omitempty"`
	Recommended   RecommendedInfo   `json:"recommended"`
}

// RecommendedInfo carries generation defaults for a model.
type RecommendedInfo struct {
	Sampler   string  `json:"sampler"`
	Scheduler string  `json:"scheduler"`
	Steps     int     `json:"steps"`
	CFG       float64 `json:"cfg"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
}

// FilteredAssets is the subset of VAEs / LoRAs / ControlNets compatible with a specific model family.
type FilteredAssets struct {
	VAEs        []AssetInfo `json:"vaes"`
	LoRAs       []AssetInfo `json:"loras"`
	ControlNets []AssetInfo `json:"controlnets"`
}

// InstalledAssetInfo is the frontend-visible result for one installed asset.
type InstalledAssetInfo struct {
	Name          string  `json:"name"`
	OriginalPath  string  `json:"original_path"`
	InstalledPath string  `json:"installed_path"`
	Type          string  `json:"type"`
	Family        string  `json:"family"`
	Format        string  `json:"format"`
	SizeMB        float64 `json:"size_mb"`
	Error         string  `json:"error,omitempty"`
}

// InstallSummary is the aggregate result of a model registration pass.
// Files are never moved — "installed" means a .howl.json manifest was written.
type InstallSummary struct {
	Registered        []InstalledAssetInfo `json:"registered"`
	AlreadyKnown      []InstalledAssetInfo `json:"already_known"`
	Failed            []InstalledAssetInfo `json:"failed"`
	TotalCount        int                  `json:"total_count"`
	RegisteredCount   int                  `json:"registered_count"`
	AlreadyKnownCount int                  `json:"already_known_count"`
	FailedCount       int                  `json:"failed_count"`
}

// ScanImageModels performs a full recursive scan of the models directory,
// classifies every asset, writes .howl.json manifests, and returns the index.
func (a *App) ScanImageModels() (ScanResult, error) {
	cfg := a.GetImageEngineConfig()
	if cfg.ModelsDir == "" {
		return ScanResult{}, fmt.Errorf("models directory not configured")
	}
	s := scanner.New(cfg.ModelsDir)
	idx, err := s.Scan()
	if err != nil {
		a.logError(errors.IMGRegistryFail, errors.High, "image.scan", "app-image-scan.go:84", err)
		return ScanResult{}, fmt.Errorf("scan failed: %w", err)
	}
	a.mergeIntoRegistry(idx)
	return a.mergeRegistryIntoResult(indexToResult(idx)), nil
}

// GetModelIndex returns the cached scan index without a full rescan.
// Falls back to a fresh scan if no cache exists.
func (a *App) GetModelIndex() (ScanResult, error) {
	cfg := a.GetImageEngineConfig()
	if cfg.ModelsDir == "" {
		return ScanResult{}, fmt.Errorf("models directory not configured")
	}
	s := scanner.New(cfg.ModelsDir)
	idx, err := s.LoadIndex()
	if err != nil {
		a.logErrorSilent(errors.IMGRegistryFail, errors.Medium, "image.index.load", "app-image-scan.go:100", err)
		return ScanResult{}, fmt.Errorf("load index failed: %w", err)
	}
	return a.mergeRegistryIntoResult(indexToResult(idx)), nil
}

// GetCompatibleAssets returns VAEs, LoRAs, and ControlNets that are compatible
// with the given model family string (e.g. "sdxl", "flux", "sd1.5").
// It sources assets from the full persistent registry (all directories), not
// just the currently configured ModelsDir.
func (a *App) GetCompatibleAssets(family string) (FilteredAssets, error) {
	// Build a combined index from the registry (all registered directories).
	var idx scanner.ModelIndex
	if a.modelRegistry != nil {
		idx = a.modelRegistry.ToModelIndex()
	} else {
		// Fallback: scan ModelsDir only when registry is unavailable.
		cfg := a.GetImageEngineConfig()
		if cfg.ModelsDir == "" {
			return FilteredAssets{}, fmt.Errorf("models directory not configured")
		}
		s := scanner.New(cfg.ModelsDir)
		var err error
		idx, err = s.LoadIndex()
		if err != nil {
			a.logErrorSilent(errors.IMGRegistryFail, errors.Medium, "image.index.load", "app-image-scan.go:124", err)
			return FilteredAssets{}, fmt.Errorf("load index failed: %w", err)
		}
	}

	f := scanner.ModelFamily(family)
	vaes := scanner.FilterVAEsForFamily(idx.VAEs, f)
	loras := scanner.FilterLoRAsForFamily(idx.LoRAs, f)
	cns := scanner.FilterControlNetsForFamily(idx.ControlNets, f)

	return FilteredAssets{
		VAEs:        entriesToInfos(vaes),
		LoRAs:       entriesToInfos(loras),
		ControlNets: entriesToInfos(cns),
	}, nil
}

// RescanImageModels deletes all cached .howl.json sidecar manifests in the
// models directory and runs a fresh classification scan. Use this after the
// classifier logic has been updated or when scan results look wrong.
func (a *App) RescanImageModels() (ScanResult, error) {
	cfg := a.GetImageEngineConfig()
	if cfg.ModelsDir == "" {
		return ScanResult{}, fmt.Errorf("models directory not configured")
	}
	if err := scanner.ClearManifests(cfg.ModelsDir); err != nil {
		a.logError(errors.IMGRegistryFail, errors.Medium, "image.manifest.clear", "app-image-scan.go:150", err)
		return ScanResult{}, fmt.Errorf("clear manifests failed: %w", err)
	}
	s := scanner.New(cfg.ModelsDir)
	idx, err := s.Scan()
	if err != nil {
		a.logError(errors.IMGRegistryFail, errors.High, "image.rescan", "app-image-scan.go:154", err)
		return ScanResult{}, fmt.Errorf("rescan failed: %w", err)
	}
	a.mergeIntoRegistry(idx)
	return a.mergeRegistryIntoResult(indexToResult(idx)), nil
}

// InstallAssets scans the models directory, writes .howl.json sidecar manifests
// for every discovered asset, and merges them into the persistent registry.
// Files are never moved — "install" means catalogue, not relocate.
func (a *App) InstallAssets(move bool) (InstallSummary, error) {
	cfg := a.GetImageEngineConfig()
	if cfg.ModelsDir == "" {
		return InstallSummary{}, fmt.Errorf("models directory not configured")
	}

	s := scanner.New(cfg.ModelsDir)
	idx, err := s.Scan()
	if err != nil {
		a.logError(errors.IMGRegistryFail, errors.High, "image.scan", "app-image-scan.go:172", err)
		return InstallSummary{}, fmt.Errorf("scan failed: %w", err)
	}

	inst := installer.New(cfg.ModelsDir, move)
	res, err := inst.Install(idx)
	if err != nil {
		a.logError(errors.IMGRegistryFail, errors.High, "image.install", "app-image-scan.go:178", err)
		return InstallSummary{}, fmt.Errorf("install failed: %w", err)
	}
	a.mergeIntoRegistry(idx)
	return convertInstallResult(res), nil
}

// GetRegisteredModels returns every model known to the persistent registry,
// across all scan directories, pruning any whose files have been deleted.
func (a *App) GetRegisteredModels() (ScanResult, error) {
	if a.modelRegistry == nil {
		return ScanResult{}, fmt.Errorf("model registry not initialised")
	}
	idx := a.modelRegistry.ToModelIndex()
	return indexToResult(idx), nil
}

// UnregisterModel removes a single model from the persistent registry by its
// absolute file path. The file itself is not touched.
func (a *App) UnregisterModel(path string) error {
	if a.modelRegistry == nil {
		return fmt.Errorf("model registry not initialised")
	}
	return a.modelRegistry.Remove(path)
}

// mergeIntoRegistry merges a scan result into the persistent registry.
// Non-fatal: logs and continues if the registry is unavailable.
func (a *App) mergeIntoRegistry(idx scanner.ModelIndex) {
	if a.modelRegistry == nil {
		return
	}
	if err := a.modelRegistry.Merge(idx); err != nil {
		a.logErrorSilent(errors.REGModelAddFail, errors.Medium, "registry.merge", "app-image-scan.go:211", err)
	}
}

// mergeRegistryIntoResult adds any registry entries not already present in
// result (matched by absolute path). This ensures every scan path returns
// the full cross-directory catalogue to the frontend.
func (a *App) mergeRegistryIntoResult(result ScanResult) ScanResult {
	if a.modelRegistry == nil {
		return result
	}
	// Build seen-path sets per category.
	seenModel := make(map[string]struct{})
	seenVAE := make(map[string]struct{})
	seenLoRA := make(map[string]struct{})
	seenCN := make(map[string]struct{})
	seenEmb := make(map[string]struct{})
	seenUp := make(map[string]struct{})
	for _, a := range result.Models {
		seenModel[a.Path] = struct{}{}
	}
	for _, a := range result.VAEs {
		seenVAE[a.Path] = struct{}{}
	}
	for _, a := range result.LoRAs {
		seenLoRA[a.Path] = struct{}{}
	}
	for _, a := range result.ControlNets {
		seenCN[a.Path] = struct{}{}
	}
	for _, a := range result.Embeddings {
		seenEmb[a.Path] = struct{}{}
	}
	for _, a := range result.Upscalers {
		seenUp[a.Path] = struct{}{}
	}

	for _, e := range a.modelRegistry.All() {
		info := registryEntryToAssetInfo(e)
		switch scanner.AssetType(e.Type) {
		case scanner.AssetDiffusion:
			if _, ok := seenModel[e.Path]; !ok {
				result.Models = append(result.Models, info)
				seenModel[e.Path] = struct{}{}
			}
		case scanner.AssetVAE:
			if _, ok := seenVAE[e.Path]; !ok {
				result.VAEs = append(result.VAEs, info)
				seenVAE[e.Path] = struct{}{}
			}
		case scanner.AssetLoRA:
			if _, ok := seenLoRA[e.Path]; !ok {
				result.LoRAs = append(result.LoRAs, info)
				seenLoRA[e.Path] = struct{}{}
			}
		case scanner.AssetControlNet:
			if _, ok := seenCN[e.Path]; !ok {
				result.ControlNets = append(result.ControlNets, info)
				seenCN[e.Path] = struct{}{}
			}
		case scanner.AssetEmbedding:
			if _, ok := seenEmb[e.Path]; !ok {
				result.Embeddings = append(result.Embeddings, info)
				seenEmb[e.Path] = struct{}{}
			}
		case scanner.AssetUpscaler:
			if _, ok := seenUp[e.Path]; !ok {
				result.Upscalers = append(result.Upscalers, info)
				seenUp[e.Path] = struct{}{}
			}
		}
	}
	return result
}

func registryEntryToAssetInfo(e registry.Entry) AssetInfo {
	return AssetInfo{
		Name:   e.Name,
		Path:   e.Path,
		Type:   e.Type,
		Family: e.Family,
		Format: e.Format,
		SizeMB: e.SizeMB,
		Files:  e.Files,
		Recommended: RecommendedInfo{
			Sampler:   e.Recommended.Sampler,
			Scheduler: e.Recommended.Scheduler,
			Width:     e.Recommended.Width,
			Height:    e.Recommended.Height,
			CFG:       e.Recommended.CFG,
			Steps:     e.Recommended.Steps,
		},
		Compatibility: assetCompatForInfo(e.Family, e.Format),
	}
}

// ---- conversion helpers ----------------------------------------------------

func indexToResult(idx scanner.ModelIndex) ScanResult {
	return ScanResult{
		Models:      entriesToInfos(idx.Models),
		VAEs:        entriesToInfos(idx.VAEs),
		LoRAs:       entriesToInfos(idx.LoRAs),
		ControlNets: entriesToInfos(idx.ControlNets),
		Embeddings:  entriesToInfos(idx.Embeddings),
		Upscalers:   entriesToInfos(idx.Upscalers),
	}
}

func convertInstallResult(r installer.InstallResult) InstallSummary {
	return InstallSummary{
		Registered:        convertInstalledAssets(r.Installed),
		AlreadyKnown:      convertInstalledAssets(r.Skipped),
		Failed:            convertInstalledAssets(r.Failed),
		TotalCount:        r.TotalCount,
		RegisteredCount:   r.InstalledCount,
		AlreadyKnownCount: r.SkippedCount,
		FailedCount:       r.FailedCount,
	}
}

func convertInstalledAssets(assets []installer.InstalledAsset) []InstalledAssetInfo {
	out := make([]InstalledAssetInfo, 0, len(assets))
	for _, a := range assets {
		out = append(out, InstalledAssetInfo{
			Name: a.Name, OriginalPath: a.OriginalPath,
			InstalledPath: a.InstalledPath, Type: a.Type,
			Family: a.Family, Format: a.Format,
			SizeMB: a.SizeMB, Error: a.Error,
		})
	}
	return out
}

func entriesToInfos(entries []scanner.AssetEntry) []AssetInfo {
	out := make([]AssetInfo, 0, len(entries))
	for _, e := range entries {
		out = append(out, AssetInfo{
			Name:          e.Name,
			Path:          e.Path,
			Type:          string(e.Type),
			Family:        string(e.Family),
			Format:        e.Format,
			SizeMB:        e.SizeMB,
			Compatibility: assetCompat(string(e.Family), e.Format),
			Files:         e.Files,
			Recommended: RecommendedInfo{
				Sampler:   e.Recommended.Sampler,
				Scheduler: e.Recommended.Scheduler,
				Steps:     e.Recommended.Steps,
				CFG:       e.Recommended.CFG,
				Width:     e.Recommended.Width,
				Height:    e.Recommended.Height,
			},
		})
	}
	return out
}

// assetCompat derives a compatibility rating from the family and format,
// matching the logic in sdcpp/discovery.go so the frontend can warn users
// before they attempt to load an unsupported model.
func assetCompat(family, format string) string {
	compat := "unknown"
	switch format {
	case "safetensors", "ckpt", "sft":
		compat = "verified"
	case "gguf":
		compat = "supported"
	case "bundle":
		compat = "supported"
	}
	if family == "chroma" {
		compat = "unsupported"
	}
	return compat
}
