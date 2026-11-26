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
	"time"

	"howl-chat/internal/backend/gguf"
	"howl-chat/internal/image/engine/sdcpp"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/resource"
)

// vramReleaseTimeout bounds how long a handoff waits for released memory to
// appear in live telemetry. It is a variable so tests can shorten it instead of
// sleeping for ten seconds to prove a timeout.
var vramReleaseTimeout = 10 * time.Second

const vramReleasePoll = 250 * time.Millisecond

// reclaimTarget is the model that has to go before the requested one can start.
// The handoff reads its state and its stop through these two calls so the
// sequence can be exercised without a real engine.
type reclaimTarget struct {
	Label  string
	Loaded func() bool
	Unload func() error
}

// reclaimOther is the one preflight handoff: refuse without consent, stop the
// other model, then wait for fresh telemetry to admit the same request. Every
// caller that unloads for the other workload goes through it, so a missing
// confirmation or an unobserved exit cannot be papered over at one call site.
func (a *App) reclaimOther(kind string, req resource.AdmissionRequest, confirmed bool, other reclaimTarget) (ModelFitResult, error) {
	if !confirmed {
		return ModelFitResult{}, fmt.Errorf("unloading the %s model to start the %s model was not confirmed; nothing was unloaded", other.Label, kind)
	}
	if other.Loaded() {
		fmt.Printf("[VRAM] unloading %s model to make room for the %s model\n", other.Label, kind)
		if err := other.Unload(); err != nil {
			return ModelFitResult{}, fmt.Errorf("could not unload the %s model: %w", other.Label, err)
		}
	}
	return a.waitForAdmissionFit(kind, req)
}

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

// CheckModelFit reports whether a chat model can load with the VRAM that is
// free right now. It is a thin, read-only view over the same admission
// authority every start uses, so the preview cannot disagree with the decision.
//
// Image requests must use CheckImageFit: an image verdict needs the actual
// generation workload, and guessing one here is how the preflight and the
// backend used to disagree.
func (a *App) CheckModelFit(kind, path string) (ModelFitResult, error) {
	switch kind {
	case "llm":
		req, err := a.chatAdmissionRequestFor(path)
		if err != nil {
			return ModelFitResult{}, err
		}
		return a.fitAdmissionRequest(kind, req)
	case "image":
		return ModelFitResult{}, fmt.Errorf("image fit requires the generation workload; use CheckImageFit")
	default:
		return ModelFitResult{}, fmt.Errorf("unknown model kind %q (want llm or image)", kind)
	}
}

// CheckImageFit reports whether this exact image generation request can start
// with the VRAM that is free right now.
func (a *App) CheckImageFit(opts imagetypes.ImageOptions) (ModelFitResult, error) {
	req, err := a.imageAdmissionFor(opts)
	if err != nil {
		return ModelFitResult{}, err
	}
	return a.fitAdmissionRequest("image", req)
}

func (a *App) chatAdmissionRequestFor(path string) (resource.AdmissionRequest, error) {
	profile, err := gguf.ExtractModelProfile(path)
	if err != nil {
		return resource.AdmissionRequest{}, fmt.Errorf("cannot inspect chat model %q: %w", path, err)
	}
	return a.chatAdmissionRequest(profile, path, false), nil
}

// imageAdmissionFor resolves the selected model exactly as the generation load
// does, then builds the one admission request shared by the preflight and the
// backend. No caller on the image generation path may build a request from a
// default or empty ImageOptions value.
func (a *App) imageAdmissionFor(opts imagetypes.ImageOptions) (resource.AdmissionRequest, error) {
	loadOpts, err := a.imageLoadOptionsForGeneration(opts.ModelName, opts)
	if err != nil {
		return resource.AdmissionRequest{}, err
	}
	same := a.imageEngineLoadedSame(opts.ModelName)
	reload := false
	if same {
		loadOpts = preserveImageLoadedOptions(loadOpts, a.images.LoadedOptions())
		reload = imageStartupOptionsDiffer(loadOpts, a.images.LoadedOptions())
	}
	req := a.imageAdmissionRequest(loadOpts, opts, same && !reload)

	// The preflight must judge the same fallback the start path will use. If the
	// stored precision does not fit but an automatic q8_0 conversion does, report
	// the fit rather than cancelling a generation that would have succeeded. The
	// lease is released here: this is a read-only preview, not a reservation.
	if q8, ok := imageAutoQuantizeRequest(req, loadOpts, !same || reload); ok {
		authority := a.admission
		if authority == nil {
			authority = resource.DefaultAuthority()
		}
		decision, lease := authority.Admit(q8)
		lease.Release()
		if decision.Admitted {
			return q8, nil
		}
	}
	return req, nil
}

// imageAutoQuantizeRequest builds the q8_0 admission request when the same
// policy that guards the start path would allow it. It never changes an
// explicit type or a GGUF.
func imageAutoQuantizeRequest(req resource.AdmissionRequest, loadOpts imagetypes.ModelLoadOptions, startupWillApply bool) (resource.AdmissionRequest, bool) {
	if !startupWillApply || !imageAutoQuantizeEligible(loadOpts) || req.WeightType != "" {
		return resource.AdmissionRequest{}, false
	}
	q8 := req
	q8.WeightType = "q8_0"
	return q8, true
}

func (a *App) fitAdmissionRequest(kind string, req resource.AdmissionRequest) (ModelFitResult, error) {
	authority := a.admission
	if authority == nil {
		authority = resource.DefaultAuthority()
	}

	decision, lease := authority.Admit(req)
	lease.Release()

	res := ModelFitResult{
		Kind:        kind,
		CanFit:      decision.Admitted,
		RequiredMB:  decision.Estimate.RequiredVRAMMB,
		FreeVRAMMB:  decision.Estimate.FreeVRAMMB,
		TotalVRAMMB: decision.Estimate.TotalVRAMMB,
		Message:     decision.Reason,
		Blocking:    []LoadedModelInfo{},
	}
	for _, model := range decision.Blocking {
		res.Blocking = append(res.Blocking, LoadedModelInfo{
			Type:         string(model.Type),
			Name:         model.Name,
			GPUMemoryMB:  model.GPUMemoryMB,
			IsGPUOffload: model.IsGPUOffload,
		})
	}
	// The dialog text belongs to fitMessage whenever there is a loaded model to
	// unload, admitted or not: the raw admission reason names the blocker but
	// never asks the user to confirm the unload.
	if decision.Admitted || len(res.Blocking) > 0 {
		res.Message = fitMessage(res)
	}
	return res, nil
}

// FreeVRAMFor unloads the image model for a chat load, then polls live
// telemetry until the same chat admission request fits. It never treats a
// successful kill as evidence that the allocation was released.
//
// confirmed is the user's answer from the fit dialog. A false value unloads
// nothing: the unload is only ever a consequence of that answer.
func (a *App) FreeVRAMFor(kind, path string, confirmed bool) (ModelFitResult, error) {
	switch kind {
	case "llm":
		req, err := a.chatAdmissionRequestFor(path)
		if err != nil {
			return ModelFitResult{}, err
		}
		return a.reclaimOther("llm", req, confirmed, reclaimTarget{
			Label:  "image",
			Loaded: func() bool { name, _ := a.loadedImageModelInfo(); return name != "" },
			Unload: a.UnloadImageModel,
		})
	case "image":
		return ModelFitResult{}, fmt.Errorf("image reclamation requires the generation workload; use FreeVRAMForImage")
	default:
		return ModelFitResult{}, fmt.Errorf("unknown model kind %q (want llm or image)", kind)
	}
}

// FreeVRAMForImage unloads the chat model for an image generation, then polls
// live telemetry until the same image admission request fits. confirmed is the
// user's answer from the fit dialog; without it nothing is unloaded.
func (a *App) FreeVRAMForImage(opts imagetypes.ImageOptions, confirmed bool) (ModelFitResult, error) {
	req, err := a.imageAdmissionFor(opts)
	if err != nil {
		return ModelFitResult{}, err
	}
	return a.reclaimOther("image", req, confirmed, reclaimTarget{
		Label:  "chat",
		Loaded: a.IsModelLoaded,
		Unload: a.UnloadModel,
	})
}

// waitForAdmissionFit keeps the same immutable request and re-admits it against
// fresh hardware snapshots. The only success is a decision that admits that
// request; a timeout returns the last decision and an explicit error.
func (a *App) waitForAdmissionFit(kind string, req resource.AdmissionRequest) (ModelFitResult, error) {
	deadline := time.NewTimer(vramReleaseTimeout)
	defer deadline.Stop()
	ticker := time.NewTicker(vramReleasePoll)
	defer ticker.Stop()

	var last ModelFitResult
	for {
		fit, err := a.fitAdmissionRequest(kind, req)
		if err != nil {
			return fit, err
		}
		last = fit
		if fit.CanFit {
			return fit, nil
		}
		select {
		case <-deadline.C:
			return last, fmt.Errorf("VRAM reclamation for %s did not complete within %s: %s", kind, vramReleaseTimeout, last.Message)
		case <-ticker.C:
		}
	}
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
	var vram int64
	for _, loaded := range resource.GetManager().ModelsByType(resource.ResourceImage) {
		if loaded.GPUMemoryMB > vram {
			vram = loaded.GPUMemoryMB
		}
	}
	return m.Name, vram
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
		return fmt.Sprintf("This model needs about %d MB and only %d MB of %s is free, so it cannot start.",
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
