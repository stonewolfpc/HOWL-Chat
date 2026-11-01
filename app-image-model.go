package main

import (
	"fmt"
	"howl-chat/internal/image/engine/sdcpp"
	"path/filepath"
	"strings"

	"howl-chat/internal/errors"
	"howl-chat/internal/image/registry"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/resource"
)

// -----------------------------------------------------------------------
// Model loading
// -----------------------------------------------------------------------

// ensureEngineReady auto-initializes the engine and loads the model if needed.
// If genOpts is provided, its startup flags are merged into lastImageModelOpts.
func (a *App) ensureEngineReady(modelName string, genOpts ...imagetypes.ImageOptions) error {
	if len(genOpts) > 0 {
		o := genOpts[0]
		if o.VAETiling {
			a.lastImageModelOpts.VAETiling = true
		}
		if o.TemporalTiling {
			a.lastImageModelOpts.TemporalTiling = true
		}
		if o.FlashAttention {
			a.lastImageModelOpts.FlashAttention = true
		}
		if o.OffloadToCPU {
			a.lastImageModelOpts.OffloadToCPU = true
		}
		if o.VAEPath != "" {
			a.lastImageModelOpts.VAEPath = o.VAEPath
		}
		if o.ClipLPath != "" {
			a.lastImageModelOpts.ClipLPath = o.ClipLPath
		}
		if o.ClipGPath != "" {
			a.lastImageModelOpts.ClipGPath = o.ClipGPath
		}
		if o.T5XXLPath != "" {
			a.lastImageModelOpts.T5XXLPath = o.T5XXLPath
		}
	}
	if a.imageEngine == nil {
		if err := a.lazyInitEngine(); err != nil {
			return err
		}
	}
	if modelName == "" {
		modelName = a.lastImageModel
	}
	if modelName == "" {
		return fmt.Errorf("no model selected")
	}
	loadOpts := a.lastImageModelOpts
	loadOpts.ModelName = modelName
	loadOpts.ModelPath = a.resolveModelPath(modelName)
	a.autofillCompanionModels(&loadOpts)
	if err := a.validateModelLoadOptions(loadOpts); err != nil {
		return err
	}

	a.lastImageModelOpts = loadOpts

	// Video renders through sd-cli, so resolve the model without spawning
	// sd-server: only one process should own the GPU at a time.
	if len(genOpts) > 0 && genOpts[0].Mode == "video" {
		if eng, ok := a.imageEngine.(*sdcpp.Engine); ok {
			return eng.SelectModel(loadOpts)
		}
	}
	return a.imageEngine.EnsureReady(modelName, loadOpts)
}

// lazyInitEngine initialises the engine on-demand during generation if it was
// not available at startup (e.g. binary configured later).
func (a *App) lazyInitEngine() error {
	cfg := a.GetImageEngineConfig()
	if cfg.BinaryPath == "" {
		cfg.BinaryPath = findSDServerBinary()
	}
	if cfg.BinaryPath == "" {
		return fmt.Errorf("sd-server binary not found")
	}
	if cfg.ModelsDir == "" {
		cfg.ModelsDir = defaultModelsDir()
	}
	if cfg.Port == 0 {
		cfg.Port = 8081
	}
	return a.reinitEngine(func(c *imagetypes.EngineConfig) { *c = cfg })
}

// SetActiveImageModel loads a model into sd-server and persists the selection.
func (a *App) SetActiveImageModel(opts interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[PANIC] SetActiveImageModel: %v\n", r)
			err = fmt.Errorf("internal panic: %v", r)
		}
	}()
	if a.imageEngine == nil {
		return fmt.Errorf("image engine not initialized")
	}

	loadOpts, parseErr := parseModelLoadOpts(opts)
	if parseErr != nil {
		return parseErr
	}

	a.lastImageModel = loadOpts.ModelName
	loadOpts.ModelPath = a.resolveModelPath(loadOpts.ModelName)
	a.autofillCompanionModels(&loadOpts)
	if err := a.validateModelLoadOptions(loadOpts); err != nil {
		return err
	}
	a.lastImageModelOpts = loadOpts
	fmt.Printf("[App] SetActiveImageModel: name=%q path=%q\n", loadOpts.ModelName, loadOpts.ModelPath)
	if err := a.imageEngine.LoadModel(loadOpts); err != nil {
		if loadOpts.OffloadToCPU {
			a.logError(errors.IMGCPUOffloadFail, errors.High, "image.load", "app-image-model.go", err)
		}
		a.logImageLoadError(err, "app-image-model.go")
		return err
	}

	resource.GetManager().RegisterModel(
		resource.ResourceImage,
		loadOpts.ModelName,
		loadOpts.ModelPath,
		4000,
		2000,
		resource.PriorityImage,
		!loadOpts.OffloadToCPU,
	)
	a.saveImageSettings()
	return nil
}

// parseModelLoadOpts converts the loosely-typed Wails opts into ModelLoadOptions.
func parseModelLoadOpts(opts interface{}) (imagetypes.ModelLoadOptions, error) {
	var out imagetypes.ModelLoadOptions
	switch v := opts.(type) {
	case string:
		out.ModelName = v
	case map[string]interface{}:
		if m, ok := v["model"].(string); ok {
			out.ModelName = m
		}
		if s, ok := v["vae"].(string); ok {
			out.VAEPath = s
		}
		if s, ok := v["clip_l"].(string); ok {
			out.ClipLPath = s
		}
		if s, ok := v["clip_g"].(string); ok {
			out.ClipGPath = s
		}
		if s, ok := v["t5xxl"].(string); ok {
			out.T5XXLPath = s
		}
		if s, ok := v["diffusion_model"].(string); ok {
			out.DiffusionModel = s
		}
		if s, ok := v["backend"].(string); ok {
			out.Backend = s
		}
		if s, ok := v["params_backend"].(string); ok {
			out.ParamsBackend = s
		}
		if b, ok := v["vae_tiling"].(bool); ok {
			out.VAETiling = b
		}
		if b, ok := v["temporal_tiling"].(bool); ok {
			out.TemporalTiling = b
		}
		if b, ok := v["circular"].(bool); ok {
			out.Circular = b
		}
		if b, ok := v["flash_attention"].(bool); ok {
			out.FlashAttention = b
		}
		if b, ok := v["offload_to_cpu"].(bool); ok {
			out.OffloadToCPU = b
		}
		if b, ok := v["disable_metadata"].(bool); ok {
			out.DisableMetadata = b
		}
		if s, ok := v["cache_mode"].(string); ok {
			out.CacheMode = s
		}
		if s, ok := v["cache_option"].(string); ok {
			out.CacheOption = s
		}
		if f, ok := v["max_vram"].(float64); ok {
			out.MaxVRAM = int(f)
		}
	default:
		return out, fmt.Errorf("invalid options type: %T", opts)
	}
	return out, nil
}

// -----------------------------------------------------------------------
// Model discovery
// -----------------------------------------------------------------------

// ListImageModels merges the engine scan with registry entries.
func (a *App) ListImageModels() ([]imagetypes.ModelInfo, error) {
	var models []imagetypes.ModelInfo
	if a.imageEngine != nil {
		m, err := a.imageEngine.ListModels()
		if err == nil && m != nil {
			models = m
		}
	}

	if a.modelRegistry != nil {
		seen := make(map[string]struct{}, len(models))
		for _, m := range models {
			seen[m.Path] = struct{}{}
		}
		for _, e := range a.modelRegistry.All() {
			if e.Type != "diffusion" {
				continue
			}
			// Menu boundary: only models living in an image directory.
			if a.modelRegistry.KindForPath(e.Path) != registry.KindImage {
				continue
			}
			if _, ok := seen[e.Path]; ok {
				continue
			}
			files := e.Files
			if files == nil {
				files = map[string]string{"model": e.Path}
			}
			models = append(models, imagetypes.ModelInfo{
				Name:           e.Name,
				Path:           e.Path,
				Family:         e.Family,
				Capabilities:   []string{"txt2img", "img2img", "inpaint"},
				ResolutionHint: resolutionHintFor(e.Family),
				Files:          files,
				Format:         e.Format,
				Compatibility:  assetCompatForInfo(e.Family, e.Format),
			})
		}
	}

	if models == nil {
		return []imagetypes.ModelInfo{}, nil
	}
	return models, nil
}

// ListImageLoRAs scans the models directory for LoRA files.
func (a *App) ListImageLoRAs() ([]string, error) {
	if a.imageEngine == nil {
		return []string{}, fmt.Errorf("image engine not initialized")
	}
	loras, err := a.imageEngine.ListLoRAs()
	if err != nil {
		return []string{}, err
	}
	if loras == nil {
		return []string{}, nil
	}
	return loras, nil
}

// resolveModelPath looks up the absolute file path for a model name from the registry.
func (a *App) resolveModelPath(name string) string {
	if a.modelRegistry == nil || name == "" {
		return ""
	}
	for _, e := range a.modelRegistry.All() {
		if e.Name == name {
			return e.Path
		}
	}
	return ""
}

func (a *App) modelEntryByName(name string) *registry.Entry {
	if a.modelRegistry == nil || name == "" {
		return nil
	}
	for _, e := range a.modelRegistry.All() {
		if e.Name == name && e.Type == "diffusion" && a.modelRegistry.KindForPath(e.Path) == registry.KindImage {
			c := e
			return &c
		}
	}
	return nil
}

func (a *App) findCompanionPath(companionType string) string {
	if a.modelRegistry == nil {
		return ""
	}
	// Companions must come from image directories: borrowing a VAE or T5 from
	// a chat directory is exactly the cross-contamination this keeps out.
	var all []registry.Entry
	for _, e := range a.modelRegistry.All() {
		if a.modelRegistry.KindForPath(e.Path) == registry.KindImage {
			all = append(all, e)
		}
	}
	for _, e := range all {
		if e.Type == companionType && e.Path != "" {
			return e.Path
		}
	}
	for _, e := range all {
		if e.Path == "" {
			continue
		}
		name := strings.ToLower(filepath.Base(e.Path))
		switch companionType {
		case "flux_ae":
			if name == "ae.safetensors" || strings.HasPrefix(name, "ae.") {
				return e.Path
			}
		case "clip_l":
			if strings.Contains(name, "clip_l") || strings.Contains(name, "clip-l") || strings.Contains(name, "clip.l") {
				return e.Path
			}
		case "t5xxl":
			if strings.Contains(name, "t5xxl") || strings.Contains(name, "t5-xxl") || strings.Contains(name, "t5_xxl") {
				return e.Path
			}
		case "vae":
			if strings.Contains(name, "vae") || strings.HasPrefix(name, "ae.") || strings.Contains(name, "first_stage") {
				return e.Path
			}
		}
	}
	return ""
}

func (a *App) autofillCompanionModels(opts *imagetypes.ModelLoadOptions) {
	if opts == nil {
		return
	}
	e := a.modelEntryByName(opts.ModelName)
	if e == nil {
		return
	}
	fam := strings.ToLower(e.Family)
	usesFluxAE := fam == "flux" || fam == "flux2" || fam == "chroma" || fam == "auraflow" || fam == "hidream" || fam == "kolors"

	if opts.VAEPath == "" && usesFluxAE && e.Files != nil && e.Files["vae"] != "" {
		opts.VAEPath = e.Files["vae"]
	}
	if opts.ClipLPath == "" && e.Files != nil && e.Files["clip_l"] != "" {
		opts.ClipLPath = e.Files["clip_l"]
	}
	if opts.ClipGPath == "" && e.Files != nil && e.Files["clip_g"] != "" {
		opts.ClipGPath = e.Files["clip_g"]
	}
	if opts.T5XXLPath == "" && e.Files != nil && e.Files["t5xxl"] != "" {
		opts.T5XXLPath = e.Files["t5xxl"]
	}
	if opts.DiffusionModel == "" && e.Files != nil && e.Files["diffusion_model"] != "" {
		opts.DiffusionModel = e.Files["diffusion_model"]
	}
	if opts.ClipLPath == "" {
		opts.ClipLPath = a.findCompanionPath("clip_l")
	}
	if opts.T5XXLPath == "" {
		opts.T5XXLPath = a.findCompanionPath("t5xxl")
	}
	if opts.VAEPath == "" && usesFluxAE {
		opts.VAEPath = a.findCompanionPath("flux_ae")
	}
}

func (a *App) validateModelLoadOptions(opts imagetypes.ModelLoadOptions) error {
	e := a.modelEntryByName(opts.ModelName)
	if e == nil {
		return nil
	}
	fam := strings.ToLower(e.Family)

	// A file whose family we could not identify is not a Stable Diffusion
	// checkpoint: sd-server cannot build a graph for it, and loading it used to
	// fail deep inside the engine with an opaque message.
	if fam == "" || fam == "unknown" {
		return fmt.Errorf("'%s' is not a recognised Stable Diffusion family, so it cannot be loaded. Supported: SD1.5/SD2.x, SDXL (incl. Turbo/Lightning), SD3/3.5, Flux, Chroma, Hidream, and Wan for video.", opts.ModelName)
	}

	switch fam {
	case "flux", "flux2":
		if opts.T5XXLPath == "" {
			return fmt.Errorf("Flux model '%s' requires T5XXL. Select a T5XXL model first.", opts.ModelName)
		}
		if opts.ClipLPath == "" {
			return fmt.Errorf("Flux model '%s' requires CLIP-L. Add/select a clip_l model file.", opts.ModelName)
		}
		if opts.VAEPath == "" {
			return fmt.Errorf("Flux model '%s' requires a VAE/AE file (for example ae.safetensors). Add/select one first.", opts.ModelName)
		}
	case "chroma":
		if opts.T5XXLPath == "" {
			return fmt.Errorf("Chroma model '%s' requires T5XXL. Add/select a T5XXL model first.", opts.ModelName)
		}
		if opts.VAEPath == "" {
			return fmt.Errorf("Chroma model '%s' requires ae.safetensors. Add/select the Flux AE first.", opts.ModelName)
		}
	case "auraflow", "hidream":
		if opts.T5XXLPath == "" {
			return fmt.Errorf("%s model '%s' requires T5XXL. Select a T5XXL model first.", e.Family, opts.ModelName)
		}
		if opts.VAEPath == "" {
			return fmt.Errorf("%s model '%s' requires ae.safetensors. Add/select the Flux AE first.", e.Family, opts.ModelName)
		}
	}

	// Companion files must belong to the same family as the checkpoint. A SD1.5
	// VAE under an SDXL model decodes garbage, and the failure surfaces only
	// after a full sampling run.
	if err := a.validateCompanionFamilies(fam, opts); err != nil {
		return err
	}
	return nil
}

// validateCompanionFamilies rejects a VAE picked from a different family than
// the loaded checkpoint. Family data comes from the image registry, so files
// that were never scanned (manual paths) are left alone rather than guessed at.
func (a *App) validateCompanionFamilies(fam string, opts imagetypes.ModelLoadOptions) error {
	if a.modelRegistry == nil || opts.VAEPath == "" {
		return nil
	}
	idx := a.modelRegistry.ToModelIndexForKind(registry.KindImage)
	want := filepath.Clean(opts.VAEPath)
	for _, v := range idx.VAEs {
		if filepath.Clean(v.Path) != want {
			continue
		}
		vfam := strings.ToLower(string(v.Family))
		if vfam == "" || vfam == "unknown" || vfam == fam {
			return nil
		}
		return fmt.Errorf("VAE '%s' is for %s models, but '%s' is %s. Pick a %s VAE or leave VAE empty to use the built-in one.",
			filepath.Base(v.Path), vfam, opts.ModelName, fam, fam)
	}
	return nil
}

// -----------------------------------------------------------------------
// Shared helpers (also used by app-image-scan.go)
// -----------------------------------------------------------------------

// assetCompatForInfo returns a compatibility string for a given family/format pair.
func assetCompatForInfo(family, format string) string {
	switch format {
	case "safetensors", "ckpt", "sft":
		return "verified"
	case "gguf", "bundle":
		return "supported"
	}
	return "unknown"
}

// resolutionHintFor returns the default resolution string for a model family.
func resolutionHintFor(family string) string {
	switch family {
	case "flux", "sdxl", "sd3", "sdxl-turbo", "chroma", "kolors", "auraflow", "hidream":
		return "1024x1024"
	case "sd2.1", "cascade":
		return "768x768"
	default:
		return "512x512"
	}
}
