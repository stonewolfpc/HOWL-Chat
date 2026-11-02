package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"howl-chat/internal/errors"
	"howl-chat/internal/image/deps"
	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/image/registry"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/resource"
)

// -----------------------------------------------------------------------
// Model loading
// -----------------------------------------------------------------------

// ensureEngineReady auto-initializes the engine and loads the model if needed.
// If genOpts is provided, its startup flags are merged into lastImageModelOpts.
func (a *App) ensureEngineReady(modelName, jobID string, genOpts ...imagetypes.ImageOptions) error {
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
	loadOpts.Family = a.registryFamily(modelName)
	a.autofillCompanionModels(&loadOpts)
	if err := a.validateModelLoadOptions(loadOpts); err != nil {
		return err
	}

	a.lastImageModelOpts = loadOpts
	// The job id belongs to this request, not to the remembered model selection:
	// it is what lets one job's cancellation stop its own load and nothing else.
	loadOpts.JobID = jobID

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
//
// Not cancellable: a direct load has no job identity, so CancelImageJob cannot
// address it. Cancellable loads run inside StartImageJob, where the job owns the
// load it started.
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
	loadOpts.Family = a.registryFamily(loadOpts.ModelName)
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
	var paths []string
	for _, e := range a.modelRegistry.All() {
		if e.Name == name && e.Path != "" {
			paths = append(paths, e.Path)
		}
	}
	return preferModelPath(paths, a.imageModelsDir())
}

// imageModelsDir is the configured image models directory, cleaned for prefix
// comparison. Empty when nothing is configured.
func (a *App) imageModelsDir() string {
	cfg := a.GetImageEngineConfig()
	if dir := strings.TrimSpace(cfg.ModelsDir); dir != "" {
		return filepath.Clean(dir)
	}
	return ""
}

// preferModelPath picks one path when several registry entries share a name.
//
// Names are not unique in practice: an InvokeAI tree and the G: models folder
// both hold krea2_turbo-Q4_K_M. Returning the first match picked the InvokeAI
// copy, sd-server answered "model metadata validation failed", and the load
// burned its full ten-minute timeout without ever trying the file that works.
// Selection is by usefulness: a real file in the configured models directory,
// then any real file, then a models-directory path that has gone missing, then
// path order so the answer is stable.
func preferModelPath(paths []string, modelsDir string) string {
	best := ""
	bestRank := 4
	for _, p := range paths {
		clean := filepath.Clean(p)
		st, err := os.Stat(clean)
		exists := err == nil && !st.IsDir()
		inModelsDir := modelsDir != "" && strings.HasPrefix(strings.ToLower(clean), strings.ToLower(modelsDir))
		rank := 3
		switch {
		case exists && inModelsDir:
			rank = 0
		case exists:
			rank = 1
		case inModelsDir:
			rank = 2
		}
		if rank < bestRank || (rank == bestRank && (best == "" || clean < best)) {
			best, bestRank = clean, rank
		}
	}
	return best
}

func (a *App) modelEntryByName(name string) *registry.Entry {
	if a.modelRegistry == nil || name == "" {
		return nil
	}
	// Resolve the same path the loader will use, so the family and companions
	// always describe the file that is actually opened.
	want := a.resolveModelPath(name)
	for _, e := range a.modelRegistry.All() {
		if e.Name != name || e.Type != "diffusion" || a.modelRegistry.KindForPath(e.Path) != registry.KindImage {
			continue
		}
		if want != "" && filepath.Clean(e.Path) != filepath.Clean(want) {
			continue
		}
		c := e
		return &c
	}
	return nil
}

// registryFamily returns the architecture the image registry recorded for a
// model, which is the same value the dependency report and the validator use.
// Empty when the model is not registered.
func (a *App) registryFamily(modelName string) string {
	e := a.modelEntryByName(modelName)
	if e == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(e.Family))
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
	usesFluxAE := sdcpp.UsesFluxAE(fam)

	// Companions resolved for the previous model must not follow a new one in.
	// Loading flux1-schnell and then switching to an SDXL checkpoint carried
	// ae.safetensors across and the load failed with "VAE 'ae.safetensors' is for
	// flux models, but ... is sdxl" — the paths live on lastImageModelOpts, and
	// autofill only ever fills empty fields, so a stale one is never revisited.
	pruneStaleCompanions(opts, fam, a.registryFamilyForPath)

	// Registry-recorded companions are candidates, not decisions: they are
	// written at scan time, before anything has checked the file is the right
	// kind or still on disk. A wrong one falls through to the picker below.
	if opts.VAEPath == "" && usesFluxAE {
		opts.VAEPath = usableCompanion(fam, "vae", e.Files["vae"])
	}
	if opts.ClipLPath == "" {
		opts.ClipLPath = usableCompanion(fam, "clip_l", e.Files["clip_l"])
	}
	if opts.ClipGPath == "" {
		opts.ClipGPath = usableCompanion(fam, "clip_g", e.Files["clip_g"])
	}
	if opts.T5XXLPath == "" {
		opts.T5XXLPath = usableCompanion(fam, "t5xxl", e.Files["t5xxl"])
	}
	if opts.DiffusionModel == "" && e.Files != nil && e.Files["diffusion_model"] != "" {
		opts.DiffusionModel = e.Files["diffusion_model"]
	}
	if opts.ClipLPath == "" {
		opts.ClipLPath = a.findCompatibleCompanion(fam, "clip_l")
	}
	if opts.T5XXLPath == "" {
		opts.T5XXLPath = a.findCompatibleCompanion(fam, "t5xxl")
	}
	// Families that encode text with an LLM (z_image, qwen_image, flux2, anima)
	// were passed to sd-server with no --llm at all: the engine only emits the
	// flag when the option is set, and nothing ever set it here.
	if opts.LLMPath == "" {
		opts.LLMPath = usableCompanion(fam, "llm", e.Files["llm"])
	}
	if opts.LLMPath == "" {
		opts.LLMPath = a.findCompatibleCompanion(fam, "llm")
	}
	if opts.VAEPath == "" && usesFluxAE {
		opts.VAEPath = a.findCompatibleCompanion(fam, "vae")
	}
	// Same gap for the VAE: only Flux-style families were auto-filled, so a
	// z_image or qwen_image checkpoint loaded with no decoder.
	if opts.VAEPath == "" && !usesFluxAE {
		if cap, ok := sdcpp.Capability(fam); ok && cap.Uses("vae") {
			opts.VAEPath = a.findCompatibleCompanion(fam, "vae")
		}
	}
}

// usableCompanion returns a registry-provided companion path only when it is the
// right kind of file for the role and still exists on disk.
//
// The registry records a model's companions at scan time. Nothing checked them
// then, so flux1-schnell carried
// H:\InvokeAI\...\text_encoder_2\bnb_llm_int8_model.safetensors as its T5 and
// sd-server aborted on GGML_ASSERT(weight_scale->type == GGML_TYPE_F32). A path
// that fails the catalogue check or has since moved is rejected here, which lets
// autofill fall through to the picker instead of into a failed load.
func usableCompanion(family, role, path string) string {
	p := strings.TrimSpace(path)
	if p == "" {
		return ""
	}
	if ok, _ := deps.Compatible(family, role, filepath.Base(p)); !ok {
		return ""
	}
	if st, err := os.Stat(p); err != nil || st.IsDir() {
		return ""
	}
	return p
}

// registryFamilyForPath returns the family the image registry recorded for a
// file. Empty when the file was never scanned, which means it cannot be judged
// and is left alone.
func (a *App) registryFamilyForPath(path string) string {
	if a.modelRegistry == nil || path == "" {
		return ""
	}
	want := filepath.Clean(path)
	for _, e := range a.modelRegistry.All() {
		if filepath.Clean(e.Path) == want {
			return strings.ToLower(strings.TrimSpace(e.Family))
		}
	}
	return ""
}

// pruneStaleCompanions clears companion paths that do not belong to the family
// about to be loaded.
//
// It is a pure function so the rule can be tested without an App: familyFor may
// return "" for a file the registry does not know, and an unknown file is kept
// rather than guessed at. Two rules decide: the catalogue's identity check, and
// the file's own recorded family (an SDXL checkpoint cannot decode Flux latents
// no matter how valid the file is).
func pruneStaleCompanions(opts *imagetypes.ModelLoadOptions, family string, familyFor func(string) string) {
	if opts == nil {
		return
	}
	target := strings.ToLower(strings.TrimSpace(family))
	if target == "" {
		return
	}
	roles := []struct {
		role string
		path *string
	}{
		{"vae", &opts.VAEPath},
		{"clip_l", &opts.ClipLPath},
		{"clip_g", &opts.ClipGPath},
		{"t5xxl", &opts.T5XXLPath},
		{"llm", &opts.LLMPath},
	}
	for _, r := range roles {
		p := strings.TrimSpace(*r.path)
		if p == "" {
			continue
		}
		if ok, _ := deps.Compatible(target, r.role, filepath.Base(p)); !ok {
			*r.path = ""
			continue
		}
		if f := familyFor(p); f != "" && !companionFamilyUsable(target, f) {
			*r.path = ""
		}
	}
}

func (a *App) validateModelLoadOptions(opts imagetypes.ModelLoadOptions) error {
	e := a.modelEntryByName(opts.ModelName)
	if e == nil {
		return nil
	}
	fam := strings.ToLower(strings.TrimSpace(e.Family))

	// A file we could not identify is not a checkpoint sd-server can build a
	// graph for. The message names what this build implements, taken from the
	// capability table instead of a hand-written list that drifts.
	if fam == "" || fam == "unknown" {
		return fmt.Errorf("'%s' could not be identified as a diffusion architecture, so it cannot be loaded. This build supports: %s.",
			opts.ModelName, strings.Join(sdcpp.SupportedFamilies(), ", "))
	}
	cap, known := sdcpp.Capability(fam)
	if !known {
		return fmt.Errorf("'%s' is %s, which this stable-diffusion.cpp build does not implement. Supported: %s.",
			opts.ModelName, e.Family, strings.Join(sdcpp.SupportedFamilies(), ", "))
	}
	if missing := missingCompanions(cap, opts); len(missing) > 0 {
		return fmt.Errorf("%s model '%s' requires %s. Select the %s model first.",
			e.Family, opts.ModelName, strings.Join(missing, ", "), missing[0])
	}

	// Companion files must belong to the same family as the checkpoint. A SD1.5
	// VAE under an SDXL model decodes garbage, and the failure surfaces only
	// after a full sampling run.
	if err := a.validateCompanionFamilies(fam, opts); err != nil {
		return err
	}
	return nil
}

// companionFamilyUsable reports whether a companion file may be used under the
// target family: the same architecture, a file the scanner could not classify,
// or a family that shares the autoencoder (Z-Image and Flux both decode with
// ae.safetensors). One rule, because the picker and the validator disagreeing is
// what produced "VAE 'model.safetensors' is for sd1.5 models" on a model the UI
// had just reported as fully equipped.
func companionFamilyUsable(targetFamily, fileFamily string) bool {
	target := strings.ToLower(strings.TrimSpace(targetFamily))
	file := strings.ToLower(strings.TrimSpace(fileFamily))
	if file == "" || file == "unknown" || file == target {
		return true
	}
	return sdcpp.SharesAutoencoder(target, file)
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
		if companionFamilyUsable(fam, vfam) {
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
	case "gguf", "bundle", "diffusers":
		return "supported"
	}
	return "unknown"
}

// resolutionHintFor returns the default resolution string for a model family.
func resolutionHintFor(family string) string {
	// Base resolution comes from the capability table: the UI must not keep its
	// own copy of which architecture is trained at 1024.
	if res := sdcpp.FamilyResolution(family); res > 0 {
		return fmt.Sprintf("%dx%d", res, res)
	}
	return "512x512"
}

// companionPath maps a required role onto the load option that carries it.
func companionPath(role string, o imagetypes.ModelLoadOptions) string {
	switch role {
	case "vae":
		return o.VAEPath
	case "clip_l":
		return o.ClipLPath
	case "clip_g":
		return o.ClipGPath
	case "t5xxl":
		return o.T5XXLPath
	case "llm":
		return o.LLMPath
	}
	return ""
}

// missingCompanions lists the required roles the caller has not supplied.
// The requirement lives in the shared capability table, so the engine and this
// check cannot disagree about what a family needs.
func missingCompanions(c sdcpp.FamilyCapability, o imagetypes.ModelLoadOptions) []string {
	var out []string
	for _, role := range c.Companions {
		if companionPath(role, o) == "" {
			out = append(out, role)
		}
	}
	return out
}
