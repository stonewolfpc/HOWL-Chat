package main

import (
	"errors"
	"testing"

	imageengine "howl-chat/internal/image/engine"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/media"
	"howl-chat/internal/resource"
)

// reloadFailingEngine models the ordered failure in Engine.LoadModel: the old
// process is already gone when the replacement fails.
type reloadFailingEngine struct {
	running bool
	loaded  *imagetypes.ModelInfo
}

func (e *reloadFailingEngine) Initialize(imagetypes.EngineConfig) error { return nil }
func (e *reloadFailingEngine) ListModels() ([]imagetypes.ModelInfo, error) {
	return nil, nil
}
func (e *reloadFailingEngine) ListLoRAs() ([]string, error) { return nil, nil }
func (e *reloadFailingEngine) LoadModel(imagetypes.ModelLoadOptions) error {
	e.running = false
	e.loaded = nil
	return errors.New("forced reload failure")
}
func (e *reloadFailingEngine) IsRunning() bool { return e.running }
func (e *reloadFailingEngine) EnsureReady(string, imagetypes.ModelLoadOptions) error {
	return errors.New("EnsureReady must not be used for a reload")
}
func (e *reloadFailingEngine) GenerateTxt2Img(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, errors.New("not used")
}
func (e *reloadFailingEngine) GenerateImg2Img(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, errors.New("not used")
}
func (e *reloadFailingEngine) GenerateInpaint(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, errors.New("not used")
}
func (e *reloadFailingEngine) GenerateVideo(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, errors.New("not used")
}
func (e *reloadFailingEngine) GetProgress() (imagetypes.ProgressInfo, error) {
	return imagetypes.ProgressInfo{}, nil
}
func (e *reloadFailingEngine) GetLoadedModel() *imagetypes.ModelInfo { return e.loaded }
func (e *reloadFailingEngine) CancelJob(string) error                { return nil }
func (e *reloadFailingEngine) Shutdown() error {
	e.running = false
	e.loaded = nil
	return nil
}

// reloadSucceedingEngine models the same replacement path with a successful
// load, so the registration matrix has both outcomes.
type reloadSucceedingEngine struct {
	running bool
	loaded  *imagetypes.ModelInfo
}

func (e *reloadSucceedingEngine) Initialize(imagetypes.EngineConfig) error { return nil }
func (e *reloadSucceedingEngine) ListModels() ([]imagetypes.ModelInfo, error) {
	return nil, nil
}
func (e *reloadSucceedingEngine) ListLoRAs() ([]string, error) { return nil, nil }
func (e *reloadSucceedingEngine) LoadModel(opts imagetypes.ModelLoadOptions) error {
	e.running = true
	e.loaded = &imagetypes.ModelInfo{Name: opts.ModelName, Family: opts.Family}
	return nil
}
func (e *reloadSucceedingEngine) IsRunning() bool { return e.running }
func (e *reloadSucceedingEngine) EnsureReady(string, imagetypes.ModelLoadOptions) error {
	return errors.New("EnsureReady must not be used for a reload")
}
func (e *reloadSucceedingEngine) GenerateTxt2Img(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, errors.New("not used")
}
func (e *reloadSucceedingEngine) GenerateImg2Img(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, errors.New("not used")
}
func (e *reloadSucceedingEngine) GenerateInpaint(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, errors.New("not used")
}
func (e *reloadSucceedingEngine) GenerateVideo(imagetypes.ImageOptions) (*imagetypes.ImageResult, error) {
	return nil, errors.New("not used")
}
func (e *reloadSucceedingEngine) GetProgress() (imagetypes.ProgressInfo, error) {
	return imagetypes.ProgressInfo{}, nil
}
func (e *reloadSucceedingEngine) GetLoadedModel() *imagetypes.ModelInfo { return e.loaded }
func (e *reloadSucceedingEngine) CancelJob(string) error                { return nil }
func (e *reloadSucceedingEngine) Shutdown() error {
	e.running = false
	e.loaded = nil
	return nil
}

func imageReloadApp(t *testing.T, oldName string, engine imageengine.ImageEngine) *App {
	t.Helper()
	images := media.NewImage(t.TempDir())
	images.SetEngine(engine)
	oldOpts := imagetypes.ModelLoadOptions{
		ModelName: oldName,
		ModelPath: "/old/model.safetensors",
		Family:    "sdxl",
		MaxVRAM:   intPtrForTest(4),
	}
	images.SetLoadOptions(oldOpts)
	images.SetLoadedOptions(oldOpts)
	app := testImageApp(8000, 8192)
	app.images = images
	return app
}

func imageRegistrationCount(name string) int {
	count := 0
	for _, m := range resource.GetManager().ModelsByType(resource.ResourceImage) {
		if m.Name == name {
			count++
		}
	}
	return count
}

// A failed reload has no live process and therefore no registration.
func TestImageReloadFailureInvalidatesPreviousRegistration(t *testing.T) {
	const oldName = "characterization-old-image"
	resource.GetManager().UnregisterModel(resource.ResourceImage, oldName)
	t.Cleanup(func() {
		resource.GetManager().UnregisterModel(resource.ResourceImage, oldName)
	})
	resource.GetManager().RegisterModel(
		resource.ResourceImage, oldName, "/old/model.safetensors", 3000, 0, resource.PriorityImage, true,
	)

	engine := &reloadFailingEngine{
		running: true,
		loaded:  &imagetypes.ModelInfo{Name: oldName, Family: "sdxl"},
	}
	app := imageReloadApp(t, oldName, engine)

	err := app.ensureEngineReadyLocked(oldName, "", imagetypes.ImageOptions{
		ModelName: oldName,
		MaxVRAM:   intPtrForTest(6),
	})
	if err == nil {
		t.Fatal("test setup: reload should have failed")
	}
	if got := imageRegistrationCount(oldName); got != 0 {
		t.Fatalf("failed reload left %d registrations, want 0", got)
	}
}

// A successful replacement has exactly one registration, and the process
// reports the model it just loaded.
func TestImageReloadSuccessLeavesExactlyOneRegistration(t *testing.T) {
	const name = "characterization-reloaded-image"
	resource.GetManager().UnregisterModel(resource.ResourceImage, name)
	t.Cleanup(func() {
		resource.GetManager().UnregisterModel(resource.ResourceImage, name)
	})
	resource.GetManager().RegisterModel(
		resource.ResourceImage, name, "/old/model.safetensors", 3000, 0, resource.PriorityImage, true,
	)

	engine := &reloadSucceedingEngine{
		running: true,
		loaded:  &imagetypes.ModelInfo{Name: name, Family: "sdxl"},
	}
	app := imageReloadApp(t, name, engine)

	if err := app.ensureEngineReadyLocked(name, "", imagetypes.ImageOptions{
		ModelName: name,
		MaxVRAM:   intPtrForTest(6),
	}); err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := imageRegistrationCount(name); got != 1 {
		t.Fatalf("successful reload left %d registrations, want 1", got)
	}
	if engine.loaded == nil || engine.loaded.Name != name {
		t.Fatalf("engine loaded model = %+v, want %q", engine.loaded, name)
	}
}
