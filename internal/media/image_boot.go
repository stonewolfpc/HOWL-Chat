package media

import (
	"fmt"
	"path/filepath"

	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/image/gallery"
	"howl-chat/internal/image/job"
	imgsettings "howl-chat/internal/image/settings"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/paths"
)

// BootOptions is what the composition root supplies to the image domain: the
// directories it resolves, and the one cross-domain hook it calls. Recording the
// image models directory in the shared model registry is a Runtime & Resources
// concern, so the domain asks for it instead of owning the registry.
type BootOptions struct {
	TempDir           string
	RegisterModelsDir func(dir string) error
	Logf              func(format string, args ...interface{})
}

// Boot constructs the image session this domain owns: gallery, sd-server
// engine, job manager, the job→gallery sink and the restored last model.
//
// App used to do this in app-image.go. Nothing here is App's business: the
// domain decides what an engine is, when it is ready, and how a finished job is
// stored. A missing binary is not fatal — the engine stays nil and the bindings
// report it — because a first run has no paths configured yet.
func (i *Image) Boot(opts BootOptions) error {
	if i == nil {
		return fmt.Errorf("image domain is nil")
	}
	logf := opts.Logf
	if logf == nil {
		logf = func(string, ...interface{}) {}
	}
	if i.Settings() == nil {
		return fmt.Errorf("image settings store not initialized")
	}

	saved, err := i.Settings().Load()
	if err != nil {
		logf("[ImageEngine] Could not load saved settings: %v", err)
	}

	cfg := saved.EngineConfig
	if cfg.BinaryPath == "" {
		cfg.BinaryPath = sdcpp.FindSDServer()
	}
	if cfg.ModelsDir == "" {
		cfg.ModelsDir = filepath.Join(paths.ProjectRoot(), "models")
	}
	if cfg.OutputDir == "" {
		cfg.OutputDir = paths.StatePath("images", "gallery")
	}
	if cfg.TempDir == "" {
		cfg.TempDir = filepath.Join(opts.TempDir, "howl-image-temp")
	}
	if cfg.Port == 0 {
		cfg.Port = 8081
	}

	// The configured models directory is an image directory: recording its kind
	// keeps its contents out of the chat model menu.
	if cfg.ModelsDir != "" && opts.RegisterModelsDir != nil {
		if err := opts.RegisterModelsDir(cfg.ModelsDir); err != nil {
			logf("[Registry] Could not record image models dir: %v", err)
		}
	}

	i.SetGallery(gallery.NewService(cfg.OutputDir))

	engine := sdcpp.NewEngine()
	if err := engine.Initialize(cfg); err != nil {
		logf("[ImageEngine] Not ready at startup (paths not configured yet): %v", err)
		i.SetJobs(job.NewManager(nil))
		return nil
	}
	i.SetEngine(engine)
	i.SetJobs(job.NewManager(engine))

	version := sdcpp.BinaryVersion(cfg.BinaryPath)
	logf("[ImageEngine] sd-server commit: %s", version)
	if sdcpp.IsKnownOldVersion(version) {
		logf("[ImageEngine] WARNING: sd-server commit %s is outdated. Flux/SD3/SDXL GGUF models may fail. Please update.", version)
	}

	if saved.LastModel != "" {
		i.SetLastModel(saved.LastModel, saved.LastModelOpts)
		logf("[ImageEngine] Restored last model: %q", saved.LastModel)
	}

	i.installJobSink(logf)
	return nil
}

// installJobSink makes this domain the one writer of finished jobs: a completed
// job is stored in the gallery here, not by the application.
func (i *Image) installJobSink(logf func(string, ...interface{})) {
	jobs := i.Jobs()
	if jobs == nil {
		return
	}
	jobs.SetOnComplete(func(id string, result *imagetypes.ImageResult) {
		if result == nil || i.Gallery() == nil {
			return
		}
		if err := i.Gallery().Save(result); err != nil {
			logf("[Image] Could not save gallery from job %s: %v", id, err)
			return
		}
		logf("[Image] Job %s result saved to gallery", id)
	})
}

// Reconfigure tears the current engine down and starts a new one with a mutated
// configuration, then persists the result. This replaced App's reinitEngine:
// constructing an engine is the image domain's own work.
func (i *Image) Reconfigure(mutate func(*imagetypes.EngineConfig)) error {
	if i == nil {
		return fmt.Errorf("image domain is nil")
	}
	if mutate == nil {
		return fmt.Errorf("reconfigure: no configuration change supplied")
	}

	cfg := i.EngineConfig()
	mutate(&cfg)

	engine := sdcpp.NewEngine()
	if err := engine.Initialize(cfg); err != nil {
		return fmt.Errorf("failed to re-initialize image engine: %w", err)
	}
	if previous := i.Engine(); previous != nil {
		_ = previous.Shutdown()
	}
	i.SetEngine(engine)
	if jobs := i.Jobs(); jobs != nil {
		jobs.SetEngine(engine)
	}
	return i.SaveSettings()
}

// EngineConfig reports the running engine's configuration, or an empty one when
// no engine is ready.
func (i *Image) EngineConfig() imagetypes.EngineConfig {
	engine := i.Engine()
	if engine == nil {
		return imagetypes.EngineConfig{}
	}
	if concrete, ok := engine.(*sdcpp.Engine); ok {
		return concrete.GetConfig()
	}
	return imagetypes.EngineConfig{}
}

// SaveSettings persists the engine configuration together with the last model
// selection. The domain owns both halves, so it owns the write.
func (i *Image) SaveSettings() error {
	if i == nil || i.Settings() == nil {
		return fmt.Errorf("image settings store not initialized")
	}
	last, opts := i.LastModel()
	return i.Settings().Save(imgsettings.ImageSettings{
		EngineConfig:  i.EngineConfig(),
		LastModel:     last,
		LastModelOpts: opts,
	})
}
