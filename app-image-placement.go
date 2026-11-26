package main

import (
	"fmt"
	"strings"
	"time"

	"howl-chat/internal/image/engine/sdcpp"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/resource"
)

// imagePlacementReport keeps the three answers separate and names the one the
// engine actually gave:
//
//	theoretical: what the measured file bytes and workload imply
//	plan:        what Howl asked sd.cpp to do (budget and placement flags)
//	actual:      what sd-server reported after loading the model
//
// Its actual half is diagnostic only. Admission remains the sole capacity
// authority; this report never feeds a number back into admission.
func (a *App) buildImagePlacementReport(phase string) map[string]interface{} {
	now := time.Now().UTC().Format(time.RFC3339)
	out := map[string]interface{}{
		"known":      false,
		"phase":      phase,
		"capturedAt": now,
	}
	if a.images == nil {
		return out
	}
	e, ok := a.images.Engine().(*sdcpp.Engine)
	if !ok || e == nil || !e.IsRunning() {
		return out
	}
	model := e.GetLoadedModel()
	if model == nil {
		return out
	}
	loadOpts := a.images.LoadedOptions()
	actual := e.GetPlacement()
	width, height, batch, steps, workloadAt := e.LastWorkload()

	// The request uses the options this process actually received, not the
	// panel's saved intent. That is what makes the plan auditable.
	req := a.imageAdmissionRequest(loadOpts, imagetypes.ImageOptions{
		Width: width, Height: height, BatchCount: batch, Steps: steps,
	}, false)
	weightBytes := req.CheckpointBytes + req.ConditionerBytes + req.VAEBytes + req.ControlNetBytes + req.LoRABytes + req.IPAdapterBytes + req.ClipVisionBytes
	residentNeedMB, activationMB, measured := resource.ImageResidentNeedMB(req)
	if !measured {
		residentNeedMB, activationMB = 0, 0
	}

	snap, snapErr := resource.GetManager().RefreshSnapshot()
	if snapErr != nil {
		snap = resource.GetManager().HardwareSnapshot()
	}
	totalMB := snap.TotalVRAMBytes / (1024 * 1024)
	freeMB := snap.FreeVRAMBytes / (1024 * 1024)
	usedMB := snap.UsedVRAMBytes / (1024 * 1024)
	reserveMB := int64(float64(totalMB) * resource.DefaultReservePercent)

	budgetGiB := 0
	if loadOpts.MaxVRAM != nil {
		budgetGiB = *loadOpts.MaxVRAM
	}
	planMode := imagePlanDescription(loadOpts)
	out["known"] = true
	out["model"] = model.Name
	out["family"] = model.Family
	out["theoretical"] = map[string]interface{}{
		"weightsMB":      weightBytes / (1024 * 1024),
		"activationMB":   activationMB,
		"residentNeedMB": residentNeedMB,
		"weightType":     loadOpts.WeightType,
		"width":          width,
		"height":         height,
		"batch":          batch,
		"steps":          steps,
		"measured":       measured,
	}
	out["plan"] = map[string]interface{}{
		"maxVRAMGiB":        budgetGiB,
		"autoMaxVRAM":       loadOpts.AutoMaxVRAM,
		"backend":           loadOpts.Backend,
		"paramsBackend":     loadOpts.ParamsBackend,
		"offloadToCPU":      loadOpts.OffloadToCPU,
		"vaeTiling":         loadOpts.VAETiling,
		"weightType":        loadOpts.WeightType,
		"autoPlacement":     loadOpts.AutoPlacement,
		"autoPlacementMode": loadOpts.AutoPlacementMode,
		"engineProfile":     loadOpts.EngineProfile,
		"ipAdapter":         loadOpts.IPAdapterPath,
		"clipVisionCPU":     loadOpts.ClipVisionCPU,
		"description":       planMode,
	}
	out["actual"] = map[string]interface{}{
		"known":         actual.Known,
		"vramMB":        actual.VRAMMB,
		"ramMB":         actual.RAMMB,
		"components":    actual.Components,
		"componentMB":   actual.ComponentMB,
		"backends":      actual.Backends,
		"offloaded":     actual.Offloaded(),
		"vaeTiled":      actual.VAETiled,
		"vaeTilingNote": actual.VAETilingNote,
		"process":       actual.Process,
	}
	out["live"] = map[string]interface{}{
		"totalVRAMMB": totalMB,
		"usedVRAMMB":  usedMB,
		"freeVRAMMB":  freeMB,
		"reserveMB":   reserveMB,
		"reserveHeld": freeMB >= reserveMB,
		"source":      snap.VRAMSource,
		"live":        snap.VRAMLive,
	}
	out["workloadAt"] = ""
	if !workloadAt.IsZero() {
		out["workloadAt"] = workloadAt.UTC().Format(time.RFC3339)
	}

	divergence := imagePlacementDivergence(loadOpts, actual, freeMB, reserveMB)
	out["divergence"] = divergence

	// Keep the original flat fields: current bindings and the image page read
	// them, and removing them would make an older panel silently lose placement.
	out["vramMB"] = actual.VRAMMB
	out["ramMB"] = actual.RAMMB
	out["components"] = actual.Components
	out["componentMB"] = actual.ComponentMB
	out["backends"] = actual.Backends
	out["offloaded"] = actual.Offloaded()
	out["vaeTiled"] = actual.VAETiled
	out["vaeTilingNote"] = actual.VAETilingNote
	if model != nil {
		out["model"] = model.Name
		out["family"] = model.Family
	}

	out["summary"] = imagePlacementSummary(planMode, actual, divergence)
	// Reclamation is reported from the lifecycle policy's own state, not inferred
	// from the placement. It states the baseline and band actually used and the
	// last decision, including the result of a recycle when one happened.
	reclaim := a.imageReclaim.report()
	out["reclaim"] = reclaim
	if sentence, ok := reclaim["summary"].(string); ok && sentence != "" {
		out["summary"] = out["summary"].(string) + " " + sentence
	}
	return out
}

// emitImagePlacement publishes the same report the binding returns. The load
// path emits it so the card appears without polling, and the generation path
// emits it again after the engine has run the workload.
func (a *App) emitImagePlacement(phase string) {
	if a == nil {
		return
	}
	a.emitEvent("image:placement", a.buildImagePlacementReport(phase))
}

// imagePlanDescription states the placement Howl asked for, including an
// automatic bounded mode. It never claims the engine used it.
func imagePlanDescription(opts imagetypes.ModelLoadOptions) string {
	if opts.EngineProfile == imagetypes.EngineProfileIPAdapter {
		vision := "CLIP-Vision on GPU"
		if opts.ClipVisionCPU {
			vision = "CLIP-Vision on CPU"
		}
		return "native IP-Adapter engine profile: adapter resident, " + vision + "; sd.cpp auto-fit places the rest within the admitted budget"
	}
	if opts.AutoPlacementMode == "cpu-diffusion" {
		return "automatic bounded placement: diffusion and CLIP on CPU, VAE on CUDA and tiled"
	}
	if opts.AutoPlacementMode == "cpu-only" || strings.EqualFold(strings.TrimSpace(opts.Backend), "cpu") {
		return "CPU-only placement: no GPU weights or compute"
	}
	if strings.TrimSpace(opts.Backend) != "" || strings.TrimSpace(opts.ParamsBackend) != "" || opts.OffloadToCPU {
		parts := []string{}
		if opts.Backend != "" {
			parts = append(parts, "backend "+opts.Backend)
		}
		if opts.ParamsBackend != "" {
			parts = append(parts, "params "+opts.ParamsBackend)
		}
		if opts.OffloadToCPU {
			parts = append(parts, "offload-to-cpu")
		}
		return "explicit placement: " + strings.Join(parts, ", ")
	}
	return "sd.cpp auto-fit chooses module placement within the admitted budget"
}

// imagePlacementDivergence compares the plan with the engine's own placement
// summary. It reports only differences the evidence supports: absent components
// (N/A) are not a mismatch, and an auto-fit plan is not compared module by
// module because it made no module claim.
func imagePlacementDivergence(opts imagetypes.ModelLoadOptions, actual sdcpp.Placement, freeMB, reserveMB int64) string {
	if !actual.Known {
		if sentence := imageFootprintSentence(actual); sentence != "" {
			return "this sd-server build did not report per-component placement. " + sentence + ". Component residency remains unverified."
		}
		return "sd-server has not reported its placement; the plan is unverified."
	}
	if opts.MaxVRAM != nil && *opts.MaxVRAM > 0 {
		budgetMB := int64(*opts.MaxVRAM) * 1024
		if actual.VRAMMB > float64(budgetMB)+32 {
			return fmt.Sprintf("engine reports %.0f MB of weights in VRAM, above the planned %d MB budget.", actual.VRAMMB, budgetMB)
		}
	}
	if opts.AutoPlacementMode == "cpu-only" && actual.VRAMMB > 32 {
		return fmt.Sprintf("plan is CPU-only but the engine reports %.0f MB in VRAM.", actual.VRAMMB)
	}
	if opts.AutoPlacementMode == "cpu-diffusion" {
		if location := actual.Components["diffusion_model"]; location == "VRAM" {
			return "plan puts diffusion on CPU, but the engine reports diffusion_model in VRAM."
		}
		if location := actual.Components["text_encoders"]; location == "VRAM" {
			return "plan puts CLIP on CPU, but the engine reports text_encoders in VRAM."
		}
		if location := actual.Components["vae"]; location == "RAM" {
			return "plan puts the VAE on CUDA, but the engine reports the VAE in RAM."
		}
	}
	if strings.TrimSpace(opts.Backend) != "" || strings.TrimSpace(opts.ParamsBackend) != "" {
		if mismatch := plannedPlacementMismatch(opts, actual); mismatch != "" {
			return mismatch
		}
	}
	if reserveMB > 0 && freeMB < reserveMB {
		return fmt.Sprintf("the 10%% reserve is not held at this observation: %d MB free, %d MB required free.", freeMB, reserveMB)
	}
	return ""
}

func plannedPlacementMismatch(opts imagetypes.ModelLoadOptions, actual sdcpp.Placement) string {
	checks := []struct {
		label string
		spec  string
	}{
		{"compute backend", opts.Backend},
		{"params backend", opts.ParamsBackend},
	}
	for _, check := range checks {
		for module, target := range parseBackendAssignments(check.spec) {
			component := placementComponentName(module)
			if component == "" {
				continue
			}
			location, ok := actual.Components[component]
			if !ok || location == "N/A" {
				continue
			}
			cpuTarget := strings.EqualFold(target, "cpu")
			if cpuTarget && location == "VRAM" {
				return fmt.Sprintf("plan puts %s on CPU (%s), but the engine reports %s in VRAM.", module, check.label, component)
			}
			if !cpuTarget && location == "RAM" {
				return fmt.Sprintf("plan puts %s on %s (%s), but the engine reports %s in RAM.", module, target, check.label, component)
			}
		}
	}
	return ""
}

func parseBackendAssignments(spec string) map[string]string {
	out := map[string]string{}
	for _, part := range strings.Split(spec, ",") {
		key, value, ok := strings.Cut(strings.TrimSpace(part), "=")
		if ok && key != "" && value != "" {
			out[strings.ToLower(key)] = value
		}
	}
	return out
}

func placementComponentName(module string) string {
	switch strings.ToLower(module) {
	case "diffusion", "diffusion-model", "diffusion_model":
		return "diffusion_model"
	case "clip", "te", "text-encoder", "text_encoder", "text-encoders":
		return "text_encoders"
	case "vae":
		return "vae"
	}
	return ""
}

func imagePlacementSummary(plan string, actual sdcpp.Placement, divergence string) string {
	if !actual.Known {
		out := "Placement planned: " + plan + ". This sd-server build did not report per-component residency, so the plan is unverified."
		if sentence := imageFootprintSentence(actual); sentence != "" {
			out += " " + sentence + "."
		}
		return out
	}
	actualSentence := ""
	switch {
	case actual.Offloaded():
		actualSentence = fmt.Sprintf("engine reports all %.0f MB of weights in RAM and 0 MB in VRAM", actual.RAMMB)
	case actual.RAMMB > 0:
		actualSentence = fmt.Sprintf("engine reports %.0f MB in VRAM and %.0f MB in RAM (partial placement)", actual.VRAMMB, actual.RAMMB)
	default:
		actualSentence = fmt.Sprintf("engine reports all %.0f MB of weights in VRAM", actual.VRAMMB)
	}
	out := "Planned: " + plan + ". Actual: " + actualSentence + "."
	if sentence := imageFootprintSentence(actual); sentence != "" {
		out += " " + sentence + "."
	}
	if divergence != "" {
		out += " Divergence: " + divergence
	}
	return out
}

// imageFootprintSentence names only measurements Howl actually took. It never
// presents a working-set number as if the engine had reported component
// residency.
func imageFootprintSentence(actual sdcpp.Placement) string {
	load, hasLoad := actual.Process["load"]
	gen, hasGen := actual.Process["generation"]
	switch {
	case hasLoad && load.Known && hasGen && gen.Known:
		return fmt.Sprintf("Measured sd-server process working set: %.0f MB after load, %.0f MB after generation", load.WorkingSetMB, gen.WorkingSetMB)
	case hasLoad && load.Known:
		return fmt.Sprintf("Measured sd-server process working set: %.0f MB after load; no completed generation measured yet", load.WorkingSetMB)
	case hasGen && gen.Known:
		return fmt.Sprintf("Measured sd-server process working set: %.0f MB after generation", gen.WorkingSetMB)
	}
	return ""
}
