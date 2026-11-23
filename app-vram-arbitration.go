package main

// app-vram-arbitration.go — VRAM arbitration between the chat and image models.
//
// One 8 GB card cannot hold a chat model and a diffusion checkpoint at the same
// time. The app owns the load state, so it is the app's job to say so before a
// load fails quietly: report what is in the way, let the user confirm, unload
// the other model, re-read the freed VRAM, then proceed.
//
// Free VRAM is always re-read live (nvidia-smi through the resource manager)
// rather than inferred from the registry, because a killed subprocess can take
// a moment to release its allocation.

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"howl-chat/internal/backend/gguf"
	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/resource"
)

// ModelFitResult answers "will this model fit right now, and what is in the way?"
type ModelFitResult struct {
	Kind        string            `json:"kind"`
	CanFit      bool              `json:"can_fit"`
	RequiredMB  int64             `json:"required_mb"`
	FreeVRAMMB  int64             `json:"free_vram_mb"`
	TotalVRAMMB int64             `json:"total_vram_mb"`
	Blocking    []LoadedModelInfo `json:"blocking"`
	Message     string            `json:"message"`
}

const (
	// imageBaseVRAMMB is the floor a diffusion checkpoint needs on the GPU even
	// when most weights live in RAM: text encoders, VAE and sampling buffers.
	// Measured on the GTX 1080 with SDXL at 512x512 (see STAGE_STATUS.md): the
	// sd-server process held ~2.1 GB with a ControlNet bound and ~1.5 GB without.
	imageBaseVRAMMB = 1600
	// imageOverheadMB covers the control image, TAESD preview and CUDA context.
	imageOverheadMB = 500
)

// CheckModelFit reports whether kind ("llm" or "image") can load with the VRAM
// that is free right now, and which loaded model is in the way if it cannot.
func (a *App) CheckModelFit(kind, path string) (ModelFitResult, error) {
	switch kind {
	case "llm", "image":
	default:
		return ModelFitResult{}, fmt.Errorf("unknown model kind %q (want llm or image)", kind)
	}

	resource.GetManager().RefreshHardware()
	hw := resource.GetManager().HardwareSnapshot()

	res := ModelFitResult{
		Kind:        kind,
		Blocking:    []LoadedModelInfo{},
		FreeVRAMMB:  hw.FreeVRAMBytes / 1024 / 1024,
		TotalVRAMMB: hw.TotalVRAMBytes / 1024 / 1024,
	}

	switch kind {
	case "llm":
		res.RequiredMB = a.estimateLLMVRAM(path)
		if name, vram := a.loadedImageModelInfo(); name != "" {
			res.Blocking = append(res.Blocking, LoadedModelInfo{
				Type: "image", Name: name, GPUMemoryMB: vram, IsGPUOffload: true,
			})
		}
	case "image":
		res.RequiredMB = a.estimateImageVRAM(path)
		if a.IsModelLoaded() {
			name := a.GetLoadedModelName()
			vram := int64(0)
			for _, m := range resource.GetManager().ModelsByType(resource.ResourceLLM) {
				if m.GPUMemoryMB > vram {
					vram = m.GPUMemoryMB
				}
			}
			res.Blocking = append(res.Blocking, LoadedModelInfo{
				Type: "llm", Name: name, GPUMemoryMB: vram, IsGPUOffload: true,
			})
		}
	}

	res.CanFit = res.RequiredMB <= res.FreeVRAMMB
	res.Message = fitMessage(res)
	return res, nil
}

// FreeVRAMFor unloads whatever stands in the way of kind, waits for the driver
// to reclaim the allocation, and re-reports the fit.
func (a *App) FreeVRAMFor(kind, path string) (ModelFitResult, error) {
	switch kind {
	case "llm":
		if name, _ := a.loadedImageModelInfo(); name != "" {
			fmt.Printf("[VRAM] unloading image model %q to make room for the chat model\n", name)
			if err := a.UnloadImageModel(); err != nil {
				return ModelFitResult{}, fmt.Errorf("could not unload the image model: %w", err)
			}
		}
	case "image":
		if a.IsModelLoaded() {
			fmt.Println("[VRAM] unloading chat model to make room for the image model")
			if err := a.UnloadModel(); err != nil {
				return ModelFitResult{}, fmt.Errorf("could not unload the chat model: %w", err)
			}
		}
	default:
		return ModelFitResult{}, fmt.Errorf("unknown model kind %q (want llm or image)", kind)
	}

	// A killed llama-server / sd-server releases VRAM asynchronously; reading
	// immediately reports the old number and would produce a false refusal.
	time.Sleep(1500 * time.Millisecond)
	return a.CheckModelFit(kind, path)
}

// loadedImageModelInfo returns the image model's display name and the VRAM it
// is holding. The image engine is not tracked by the resource registry, so its
// state is read from the engine itself.
func (a *App) loadedImageModelInfo() (string, int64) {
	e, ok := a.images.Engine().(*sdcpp.Engine)
	if !ok || e == nil {
		return "", 0
	}
	m := e.GetLoadedModel()
	if m == nil || m.Name == "" {
		return "", 0
	}
	return m.Name, a.estimateImageVRAM("")
}

// estimateLLMVRAM estimates the GPU footprint of a chat model with the current
// settings, using the same estimator LoadModel uses to register the model.
func (a *App) estimateLLMVRAM(path string) int64 {
	if path == "" {
		return 0
	}
	profile, err := gguf.ExtractModelProfile(path)
	if err != nil || profile == nil {
		return 0
	}
	ctxSize, _ := strconv.Atoi(intSetting(a.settings.Model(), "context_size", 4096))
	gpuLayers, _ := strconv.Atoi(intSetting(a.settings.Model(), "gpu_layers", 0))
	est := gguf.EstimateMemory(profile, gpuLayers, ctxSize)
	if est.GPUMB > 0 {
		return est.GPUMB
	}
	return profile.FileSizeBytes / 1024 / 1024
}

// estimateImageVRAM estimates what the diffusion engine needs on the GPU.
// sd-server offloads most weights to RAM, so the checkpoint's file size is not
// the requirement; the floor is text encoders, VAE and sampling buffers.
func (a *App) estimateImageVRAM(path string) int64 {
	need := int64(imageBaseVRAMMB)
	// A bound ControlNet is loaded into VRAM alongside the checkpoint, so it
	// has to be counted or the estimate is 2.4 GB short on this machine.
	if cn := a.boundControlNetBytes(path); cn > 0 {
		need += cn / 1024 / 1024
	}
	return need + imageOverheadMB
}

// boundControlNetBytes returns the size of the ControlNet that ships inside the
// selected checkpoint's bundle directory, if there is one.
func (a *App) boundControlNetBytes(modelPath string) int64 {
	dir := ""
	if modelPath != "" {
		dir = filepath.Dir(modelPath)
	} else if a.images.LastModelName() != "" {
		dir = filepath.Dir(a.images.LastModelName())
	}
	if dir == "" || dir == "." {
		dir = a.GetImageEngineConfig().ModelsDir
	}
	if dir == "" {
		return 0
	}
	files, err := filepath.Glob(filepath.Join(dir, "*control*net*.safetensors"))
	if err != nil || len(files) == 0 {
		return 0
	}
	var size int64
	for _, f := range files {
		if fi, err := os.Stat(f); err == nil && fi.Size() > size {
			size = fi.Size()
		}
	}
	return size
}

// fitMessage is the user-facing sentence shown in the confirm dialog.
func fitMessage(r ModelFitResult) string {
	card := fmt.Sprintf("%d MB", r.TotalVRAMMB)
	if r.CanFit {
		return fmt.Sprintf("Fits: needs about %d MB and %d MB is free of %s.", r.RequiredMB, r.FreeVRAMMB, card)
	}
	if len(r.Blocking) == 0 {
		return fmt.Sprintf("This model needs about %d MB and only %d MB of %s is free. It may fall back to system RAM and run slowly.",
			r.RequiredMB, r.FreeVRAMMB, card)
	}
	b := r.Blocking[0]
	what := "chat model"
	if b.Type == "image" {
		what = "image model"
	}
	headline := fmt.Sprintf("Your %s %q is loaded and this one needs about %d MB, with %d MB free of %s.",
		what, b.Name, r.RequiredMB, r.FreeVRAMMB, card)
	if b.GPUMemoryMB > 0 {
		headline = fmt.Sprintf("Your %s %q is using about %d MB, and this one needs about %d MB — more than the %d MB free of %s.",
			what, b.Name, b.GPUMemoryMB, r.RequiredMB, r.FreeVRAMMB, card)
	}
	return headline + "\n\nUnload it and continue?"
}
