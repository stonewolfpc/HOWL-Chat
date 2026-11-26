package media

import (
	"sync"

	imageengine "howl-chat/internal/image/engine"
	"howl-chat/internal/image/gallery"
	"howl-chat/internal/image/job"
	imgsettings "howl-chat/internal/image/settings"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/process"
)

// Image owns the image engine session: the active engine, its gallery, its job
// manager, the persisted settings store, and the last model selection.
//
// It does not own model admission, the shared model registry, or the chat
// model profile — those belong to Runtime & Resources. Generation and model
// loading are currently clients of this owner and are recorded as retained.
type Image struct {
	mu       sync.RWMutex
	engine   imageengine.ImageEngine
	gallery  *gallery.Service
	jobs     *job.Manager
	settings *imgsettings.Store
	last     string
	lastOpts imagetypes.ModelLoadOptions
	// loadedOpts is what the running process actually started with, which can
	// differ from the user's requested options when admission selected an
	// automatic q8_0 fallback. It is runtime truth, not saved user intent.
	loadedOpts imagetypes.ModelLoadOptions
	// processes is the shared process authority, when the composition root
	// supplies one. Engines record sd-server and sd-cli there.
	processes *process.Registry
	// onEngineExit invalidates registrations after the engine ends. The engine
	// reports exit; it does not write the resource registry.
	onEngineExit func()
}

// NewImage builds the owner with a settings store rooted at configDir.
func NewImage(configDir string) *Image {
	return &Image{settings: imgsettings.NewStore(configDir)}
}

// SetProcessRegistry records the shared process authority for engines this
// domain constructs.
func (i *Image) SetProcessRegistry(registry *process.Registry) {
	if i == nil {
		return
	}
	i.mu.Lock()
	i.processes = registry
	i.mu.Unlock()
}

// SetOnEngineExit records the registration invalidation observer for engines
// this domain constructs.
func (i *Image) SetOnEngineExit(fn func()) {
	if i == nil {
		return
	}
	i.mu.Lock()
	i.onEngineExit = fn
	i.mu.Unlock()
}

// OnEngineExit returns the registration invalidation observer, or nil.
func (i *Image) OnEngineExit() func() {
	if i == nil {
		return nil
	}
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.onEngineExit
}

// ProcessRegistry returns the shared process authority, or nil.
func (i *Image) ProcessRegistry() *process.Registry {
	if i == nil {
		return nil
	}
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.processes
}

// Engine returns the active engine, or nil when none is ready.
func (i *Image) Engine() imageengine.ImageEngine {
	if i == nil {
		return nil
	}
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.engine
}

// SetEngine records the active engine.
func (i *Image) SetEngine(engine imageengine.ImageEngine) {
	if i == nil {
		return
	}
	i.mu.Lock()
	i.engine = engine
	i.mu.Unlock()
}

// Gallery returns the gallery service.
func (i *Image) Gallery() *gallery.Service {
	if i == nil {
		return nil
	}
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.gallery
}

// SetGallery records the gallery service.
func (i *Image) SetGallery(service *gallery.Service) {
	if i == nil {
		return
	}
	i.mu.Lock()
	i.gallery = service
	i.mu.Unlock()
}

// Jobs returns the job manager.
func (i *Image) Jobs() *job.Manager {
	if i == nil {
		return nil
	}
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.jobs
}

// SetJobs records the job manager.
func (i *Image) SetJobs(manager *job.Manager) {
	if i == nil {
		return
	}
	i.mu.Lock()
	i.jobs = manager
	i.mu.Unlock()
}

// Settings returns the persisted settings store.
func (i *Image) Settings() *imgsettings.Store {
	if i == nil {
		return nil
	}
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.settings
}

// LastModel returns the last model name and the options it was loaded with.
func (i *Image) LastModel() (string, imagetypes.ModelLoadOptions) {
	if i == nil {
		return "", imagetypes.ModelLoadOptions{}
	}
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.last, i.lastOpts
}

// LastModelName returns only the last model name.
func (i *Image) LastModelName() string {
	name, _ := i.LastModel()
	return name
}

// SetLastModel records the last model selection.
func (i *Image) SetLastModel(name string, opts imagetypes.ModelLoadOptions) {
	if i == nil {
		return
	}
	i.mu.Lock()
	i.last = name
	i.lastOpts = opts
	i.mu.Unlock()
}

// SetActiveModel records the selection without options.
func (i *Image) SetActiveModel(name string) {
	i.SetLastModel(name, i.LoadOptions())
}

// LoadOptions returns the options the last selection was loaded with.
func (i *Image) LoadOptions() imagetypes.ModelLoadOptions {
	if i == nil {
		return imagetypes.ModelLoadOptions{}
	}
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.lastOpts
}

// SetLoadOptions records the options for the current selection.
func (i *Image) SetLoadOptions(opts imagetypes.ModelLoadOptions) {
	if i == nil {
		return
	}
	i.mu.Lock()
	i.lastOpts = opts
	i.mu.Unlock()
}

// LoadedOptions returns the options the running process was actually started
// with. It is empty until a model has loaded.
func (i *Image) LoadedOptions() imagetypes.ModelLoadOptions {
	if i == nil {
		return imagetypes.ModelLoadOptions{}
	}
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.loadedOpts
}

// SetLoadedOptions records the actual successful start options.
func (i *Image) SetLoadedOptions(opts imagetypes.ModelLoadOptions) {
	if i == nil {
		return
	}
	i.mu.Lock()
	i.loadedOpts = opts
	i.mu.Unlock()
}

// Close shuts the engine down. The gallery and job manager are process-local
// and hold no resource of their own.
func (i *Image) Close() error {
	engine := i.Engine()
	if engine == nil {
		return nil
	}
	return engine.Shutdown()
}
