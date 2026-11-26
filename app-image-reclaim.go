package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/image/engine/sdcpp"
	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/resource"
)

// Image reclamation is the resident engine's own post-generation policy. The
// existing authorities stay in place: internal/resource admits a load and owns
// the budget and the reserve, sd.cpp's backend_fit chooses placement and retries
// the VAE, and Image.LoadedOptions records the process that is actually running.
// This policy only decides whether that already-admitted process should be
// recycled after a generation; it never chooses a model, budget or placement.
//
// The pinned sd-server (3b12781) has no generation-boundary flush, unload or
// reset route. Its only runtime release surfaces are startup flags (for example
// --weight-residency-ttl) and process teardown. A recycle is therefore the
// strongest lever that actually exists. It is gated on measurement so the
// normal case keeps reusable weights and the same PID.
const imageReclaimPolicyName = "measurement-gated process recycle"

// The band is deliberately not a fixed timer or a fixed number of generations.
// It grows with the work that plausibly accompanies a generation: decode
// activation, the number of loaded model components, batch count, and the
// placement mode. A 512x512 one-image GPU-resident run and a 1344x1024
// CPU-diffusion batch therefore get different allowances.
const (
	imageReclaimMinimumBandMB       = 512.0
	imageReclaimActivationShare     = 0.50
	imageReclaimComponentShare      = 0.02
	imageReclaimComponentStepMB     = 32.0
	imageReclaimBatchStepMB         = 128.0
	imageReclaimMaxBatchAllowanceMB = 512.0
	imageReclaimGPUResidentMB       = 256.0
	imageReclaimCPUDiffusionMB      = 512.0
	imageReclaimCPUOnlyMB           = 768.0
)

type imageReclaimEvent struct {
	At                 string  `json:"at"`
	Reclaimed          bool    `json:"reclaimed"`
	Reason             string  `json:"reason"`
	PIDBefore          int     `json:"pidBefore"`
	PIDAfter           int     `json:"pidAfter"`
	BaselineMB         float64 `json:"baselineMB"`
	BandMB             float64 `json:"bandMB"`
	ThresholdMB        float64 `json:"thresholdMB"`
	MeasuredMB         float64 `json:"measuredMB"`
	GrowthMB           float64 `json:"growthMB"`
	InjectedGrowthMB   float64 `json:"injectedGrowthMB,omitempty"`
	WorkingSetBeforeMB float64 `json:"workingSetBeforeMB"`
	WorkingSetAfterMB  float64 `json:"workingSetAfterMB,omitempty"`
	FreeBeforeMB       int64   `json:"freeBeforeMB"`
	FreeAfterMB        int64   `json:"freeAfterMB,omitempty"`
	ReserveMB          int64   `json:"reserveMB"`
	Workload           string  `json:"workload"`
	ReloadError        string  `json:"reloadError,omitempty"`
	// Blocked is a refusal, not a retain decision: another owner forbade the
	// reclamation evaluation. It is reported so the UI does not present stale
	// evidence as if the policy had simply chosen to keep the weights.
	Blocked              bool   `json:"blocked,omitempty"`
	BlockedReason        string `json:"blockedReason,omitempty"`
	WorkflowEngine       string `json:"workflowEngine,omitempty"`
	WorkflowActiveClaims int    `json:"workflowActiveClaims,omitempty"`
}

type imageReclaimPolicy struct {
	mu             sync.Mutex
	baselinePID    int
	baselineMB     float64
	baselineSource string
	calibrated     bool
	bandMB         float64
	thresholdMB    float64
	evaluations    int
	testGrowthUsed bool
	last           *imageReclaimEvent
	events         []imageReclaimEvent
}

// reclaimImageEngineAfterGeneration is called by both the synchronous binding
// and the job manager's all-outcome callback. It returns nothing: reclamation is
// maintenance after the result, never a reason to turn a successful generation
// into a failure. A failed reload is reported, and the next generation's normal
// ensure-ready path can still recover.
func (a *App) reclaimImageEngineAfterGeneration(name string, opts imagetypes.ImageOptions, genErr error) {
	// This is the same lock ensureEngineReady takes. A queued job cannot begin
	// its admission/load decision while the post-generation recycle is running.
	a.imageStartMu.Lock()
	defer a.imageStartMu.Unlock()
	a.reclaimImageEngineAfterGenerationLocked(name, opts, genErr)
}

// reclaimImageEngineAfterGenerationLocked is the policy body for callers that
// already hold imageStartMu (the sequential workflow). Keeping one body means
// the workflow does not add a second reclamation decision or unlock and re-lock
// between stages, which would let a queued job interleave.
func (a *App) reclaimImageEngineAfterGenerationLocked(_ string, opts imagetypes.ImageOptions, genErr error) {
	if a == nil || a.images == nil || opts.Mode == "video" {
		return
	}
	if a.workflowReclaimBlocked() {
		lease := a.imageWorkflow
		event := &imageReclaimEvent{
			At:                   time.Now().UTC().Format(time.RFC3339),
			Reclaimed:            false,
			Blocked:              true,
			BlockedReason:        "workflow parent lease active",
			Reason:               "workflow parent lease active",
			WorkflowEngine:       lease.EngineID(),
			WorkflowActiveClaims: lease.ActiveClaimCount(),
			Workload:             nonEmpty(opts.Mode, "generation"),
		}
		state := &a.imageReclaim
		state.mu.Lock()
		state.evaluations++
		state.last = event
		state.mu.Unlock()
		a.logImageReclaim(event)
		return
	}
	engine, ok := a.images.Engine().(*sdcpp.Engine)
	if !ok || engine == nil || !engine.IsRunning() {
		return
	}

	if !engine.IsRunning() {
		return
	}
	footprints := engine.GetFootprints()
	load, loadKnown := footprints["load"]
	measured, measuredKnown := footprints["generation"]
	if !measuredKnown || !measured.Known {
		// A caller that never reached generation has no post-generation evidence
		// to act on. Do not treat the readiness footprint as a generation.
		return
	}
	if !loadKnown || !load.Known {
		load = measured
	}

	loaded := a.images.LoadedOptions()
	if loaded.ModelName == "" {
		if info := engine.GetLoadedModel(); info != nil {
			loaded.ModelName = info.Name
		}
	}
	if loaded.ModelName == "" {
		return
	}

	// The engine has run a real generation and reported its placement: this is
	// the first honest moment to measure image residency and reconcile the
	// ledger with it.
	a.recordImageResidencyFromGeneration(engine, loaded.ModelName, opts, loaded.WeightType)

	req := a.imageAdmissionRequest(loaded, opts, false)
	componentBytes := req.CheckpointBytes + req.ConditionerBytes + req.VAEBytes + req.ControlNetBytes + req.LoRABytes
	componentMB := float64(componentBytes) / (1024 * 1024)
	components := 0
	for _, b := range []int64{req.CheckpointBytes, req.ConditionerBytes, req.VAEBytes, req.ControlNetBytes, req.LoRABytes} {
		if b > 0 {
			components++
		}
	}
	_, activationMB, estimateKnown := resource.ImageResidentNeedMB(req)
	if !estimateKnown {
		activationMB = 0
	}
	batch := opts.BatchCount
	if batch <= 0 {
		batch = 1
	}
	band := imageReclaimBandMB(activationMB, componentMB, components, batch, loaded)
	snap, snapErr := resource.GetManager().RefreshSnapshot()
	freeBeforeMB := snap.FreeVRAMBytes / (1024 * 1024)
	totalVRAMMB := snap.TotalVRAMBytes / (1024 * 1024)
	reserveMB := int64(float64(totalVRAMMB) * resource.DefaultReservePercent)
	if snapErr != nil {
		// The snapshot is live hardware evidence. If it failed, do not invent a
		// reserve breach; the measured working-set band still applies.
		freeBeforeMB, totalVRAMMB, reserveMB = 0, 0, 0
	}

	workload := fmt.Sprintf("%s %dx%d batch=%d steps=%d components=%d activation=%dMB",
		nonEmpty(opts.Mode, "generation"), opts.Width, opts.Height, batch, opts.Steps, components, activationMB)

	state := &a.imageReclaim
	state.mu.Lock()
	pid := engine.ServerPID()
	if state.baselinePID != pid || state.baselineMB <= 0 {
		state.baselinePID = pid
		state.baselineMB = load.WorkingSetMB
		state.baselineSource = "after readiness (lazy-load floor)"
		state.calibrated = false
	}
	// sd-server loads weights lazily, so the readiness footprint is roughly the
	// loader's floor (171 MB here) rather than the resident model. The first
	// successful generation is the first honest resident baseline. Without this
	// calibration every first generation would look like a multi-gigabyte leak.
	reserveBreach := reserveMB > 0 && freeBeforeMB <= reserveMB
	if !state.calibrated {
		state.bandMB = band
		state.thresholdMB = state.baselineMB + band
		if genErr == nil {
			state.baselineMB = measured.WorkingSetMB
			state.baselineSource = "first successful generation (resident weights)"
			state.thresholdMB = state.baselineMB + band
			state.calibrated = true
			if !reserveBreach {
				state.evaluations++
				event := &imageReclaimEvent{
					At:                 time.Now().UTC().Format(time.RFC3339),
					Reclaimed:          false,
					Reason:             "resident baseline established; weights retained",
					PIDBefore:          pid,
					PIDAfter:           pid,
					BaselineMB:         state.baselineMB,
					BandMB:             band,
					ThresholdMB:        state.thresholdMB,
					MeasuredMB:         measured.WorkingSetMB,
					WorkingSetBeforeMB: measured.WorkingSetMB,
					FreeBeforeMB:       freeBeforeMB,
					ReserveMB:          reserveMB,
					Workload:           workload,
				}
				state.last = event
				state.mu.Unlock()
				a.logImageReclaim(event)
				return
			}
		}
	}

	growth := measured.WorkingSetMB - state.baselineMB
	if growth < 0 {
		growth = 0
	}
	injected := 0.0
	// Test-only hook. It changes a real decision (there is no inert UI control);
	// it is consumed once and is never exposed as a product setting. A normal
	// process cannot produce this growth and therefore never reclaims.
	if v := imageReclaimTestGrowthMB(); v > 0 && !state.testGrowthUsed {
		injected = v
		state.testGrowthUsed = true
	}
	effectiveGrowth := growth + injected

	reason := imageReclaimReason(genErr, engine.GenerationState() == "abandoned", freeBeforeMB, reserveMB, effectiveGrowth, band)

	event := &imageReclaimEvent{
		At:                 time.Now().UTC().Format(time.RFC3339),
		Reclaimed:          reason != "",
		Reason:             reason,
		PIDBefore:          pid,
		PIDAfter:           pid,
		BaselineMB:         state.baselineMB,
		BandMB:             band,
		ThresholdMB:        state.baselineMB + band,
		MeasuredMB:         measured.WorkingSetMB,
		GrowthMB:           growth,
		InjectedGrowthMB:   injected,
		WorkingSetBeforeMB: measured.WorkingSetMB,
		FreeBeforeMB:       freeBeforeMB,
		ReserveMB:          reserveMB,
		Workload:           workload,
	}
	state.bandMB = band
	state.thresholdMB = state.baselineMB + band
	state.evaluations++
	state.last = event
	reclaim := event.Reclaimed
	state.mu.Unlock()

	if !reclaim {
		a.logImageReclaim(event)
		return
	}

	// The only real release lever in the pinned build is process teardown. The
	// reload is a new start of the same process, so it needs a new decision: an
	// earlier admission was a statement about that moment, not a permanent
	// permission. It also takes the reload transition so a concurrent start or
	// unload cannot interleave with it.
	reclaimOpts := loaded
	reclaimOpts.JobID = ""
	reloadErr := func() error {
		transition, gateErr := a.beginLifecycleTransition(resource.ResourceImage, "reloading", lifecycleStopWait)
		if gateErr != nil {
			return gateErr
		}
		defer transition.Release()

		authority := a.admission
		if authority == nil {
			authority = resource.DefaultAuthority()
		}
		decision, lease := authority.Admit(req)
		if !decision.Admitted {
			return fmt.Errorf("reload refused by admission: %s", decision.Reason)
		}
		defer lease.Release()
		return engine.LoadModel(reclaimOpts)
	}()

	afterWS := 0.0
	if loadFP, ok := engine.GetFootprints()["load"]; ok && loadFP.Known {
		afterWS = loadFP.WorkingSetMB
	}
	afterFree := int64(0)
	if afterSnap, err := resource.GetManager().RefreshSnapshot(); err == nil {
		afterFree = afterSnap.FreeVRAMBytes / (1024 * 1024)
	}

	state.mu.Lock()
	event.PIDAfter = engine.ServerPID()
	event.WorkingSetAfterMB = afterWS
	event.FreeAfterMB = afterFree
	if reloadErr != nil {
		event.ReloadError = reloadErr.Error()
	}
	state.events = append(state.events, *event)
	if len(state.events) > 10 {
		state.events = state.events[len(state.events)-10:]
	}
	state.last = event
	if reloadErr == nil {
		state.baselinePID = event.PIDAfter
		state.baselineMB = afterWS
		state.baselineSource = "after reload readiness"
		state.bandMB = band
		state.thresholdMB = afterWS + band
		state.calibrated = false
	} else {
		// The process is gone or failed to start. Leave no baseline claiming
		// otherwise; the next ensure-ready load will establish one.
		state.baselinePID = 0
		state.baselineMB = 0
		state.bandMB = band
		state.thresholdMB = 0
		state.calibrated = false
	}
	state.mu.Unlock()

	a.logImageReclaim(event)
	a.emitEvent("image:reclaim", event)
	a.emitImagePlacement("generation")
}

// imageReclaimReason is the policy's one decision table. Keeping it pure makes
// the reserve branch directly assertable without replacing the production
// snapshot or admission path.
func imageReclaimReason(genErr error, abandoned bool, freeMB, reserveMB int64, growthMB, bandMB float64) string {
	switch {
	case genErr != nil && abandoned:
		return fmt.Sprintf("failed generation left the engine state abandoned: %v", genErr)
	case reserveMB > 0 && freeMB <= reserveMB:
		return fmt.Sprintf("free VRAM %dMB fell to or below the %dMB (10%%) reserve", freeMB, reserveMB)
	case growthMB > bandMB:
		return fmt.Sprintf("working-set growth %.0fMB exceeded the workload band %.0fMB", growthMB, bandMB)
	}
	return ""
}

func (a *App) logImageReclaim(event *imageReclaimEvent) {

	decision := "retain"
	if event.Reclaimed {
		decision = "reclaim"
	} else if event.Blocked {
		decision = "blocked"
	}
	fmt.Printf("[ImageReclaim] policy=%q decision=%s reason=%q baseline=%.0fMB band=%.0fMB threshold=%.0fMB measured=%.0fMB growth=%.0fMB injected=%.0fMB free=%dMB reserve=%dMB pid=%d->%d workload=%q\n",
		imageReclaimPolicyName, decision, event.Reason, event.BaselineMB, event.BandMB, event.ThresholdMB,
		event.MeasuredMB, event.GrowthMB, event.InjectedGrowthMB, event.FreeBeforeMB, event.ReserveMB,
		event.PIDBefore, event.PIDAfter, event.Workload)
	if event.Reclaimed {
		fmt.Printf("[ImageReclaim] recycled process and reloaded weights: before=%0.fMB free=%dMB after=%.0fMB free=%dMB reload_error=%q\n",
			event.WorkingSetBeforeMB, event.FreeBeforeMB, event.WorkingSetAfterMB, event.FreeAfterMB, event.ReloadError)
	}
}

func (p *imageReclaimPolicy) report() map[string]interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := map[string]interface{}{
		"policy":         imageReclaimPolicyName,
		"evaluated":      p.evaluations > 0,
		"baselineMB":     p.baselineMB,
		"baselineSource": p.baselineSource,
		"bandMB":         p.bandMB,
		"thresholdMB":    p.thresholdMB,
		"reclaims":       len(p.events),
	}
	if p.last != nil {
		out["measuredMB"] = p.last.MeasuredMB
		out["growthMB"] = p.last.GrowthMB
		out["last"] = *p.last
		out["summary"] = p.summaryLocked()
		if p.last.Blocked {
			out["blocked"] = true
			out["blockedReason"] = p.last.BlockedReason
			out["workflowEngine"] = p.last.WorkflowEngine
			out["workflowActiveClaims"] = p.last.WorkflowActiveClaims
		}
	} else {
		out["summary"] = "Measurement-gated reclamation armed; no completed generation measured yet."
	}
	if len(p.events) > 0 {
		out["events"] = append([]imageReclaimEvent(nil), p.events...)
		out["lastReclaim"] = p.events[len(p.events)-1]
	}
	return out
}

func (p *imageReclaimPolicy) summaryLocked() string {
	if p.last == nil {
		return "Measurement-gated reclamation armed; no completed generation measured yet."
	}
	if p.last.Blocked {
		return fmt.Sprintf("Reclamation blocked: %s (workflow engine %s, %d active claim(s)); no reclaim decision was made.",
			p.last.BlockedReason, p.last.WorkflowEngine, p.last.WorkflowActiveClaims)
	}
	if p.last.Reclaimed {
		return fmt.Sprintf("Reclaimed %d time(s); last: %s (baseline %.0fMB, band %.0fMB, measured %.0fMB, free %dMB, reserve %dMB; working set %.0f -> %.0fMB).",
			len(p.events), p.last.Reason, p.last.BaselineMB, p.last.BandMB, p.last.MeasuredMB,
			p.last.FreeBeforeMB, p.last.ReserveMB, p.last.WorkingSetBeforeMB, p.last.WorkingSetAfterMB)
	}
	if len(p.events) > 0 {
		lastReclaim := p.events[len(p.events)-1]
		return fmt.Sprintf("Weights retained now: growth %.0fMB within band %.0fMB (baseline %.0fMB, threshold %.0fMB, measured %.0fMB); %d earlier reclaim(s), last: %s (free %d -> %d MB, reserve %d MB, working set %.0f -> %.0f MB).",
			p.last.GrowthMB, p.last.BandMB, p.last.BaselineMB, p.last.ThresholdMB, p.last.MeasuredMB,
			len(p.events), lastReclaim.Reason, lastReclaim.FreeBeforeMB, lastReclaim.FreeAfterMB,
			lastReclaim.ReserveMB, lastReclaim.WorkingSetBeforeMB, lastReclaim.WorkingSetAfterMB)
	}
	return fmt.Sprintf("Weights retained: growth %.0fMB within band %.0fMB (baseline %.0fMB, threshold %.0fMB, measured %.0fMB).",
		p.last.GrowthMB, p.last.BandMB, p.last.BaselineMB, p.last.ThresholdMB, p.last.MeasuredMB)
}

// imageReclaimBandMB is the workload allowance above the resident baseline. It
// is a policy band, not an admission estimate: a larger image, batch or
// CPU/partial placement gets a larger allowance because its transient working
// set is larger, while a single GPU-resident image gets the tightest band.
func imageReclaimBandMB(activationMB int64, componentMB float64, components, batch int, loaded imagetypes.ModelLoadOptions) float64 {
	if batch < 1 {
		batch = 1
	}
	placement := imageReclaimGPUResidentMB
	switch {
	case strings.EqualFold(loaded.Backend, "cpu"):
		placement = imageReclaimCPUOnlyMB
	case loaded.AutoPlacementMode == "cpu-only":
		placement = imageReclaimCPUOnlyMB
	case loaded.AutoPlacementMode == "cpu-diffusion" || loaded.OffloadToCPU:
		placement = imageReclaimCPUDiffusionMB
	}
	componentAllowance := componentMB*imageReclaimComponentShare + float64(components)*imageReclaimComponentStepMB
	batchAllowance := float64(batch-1) * imageReclaimBatchStepMB
	if batchAllowance > imageReclaimMaxBatchAllowanceMB {
		batchAllowance = imageReclaimMaxBatchAllowanceMB
	}
	band := float64(activationMB)*imageReclaimActivationShare + componentAllowance + batchAllowance + placement
	if band < imageReclaimMinimumBandMB {
		band = imageReclaimMinimumBandMB
	}
	return band
}

func imageReclaimTestGrowthMB() float64 {
	raw := strings.TrimSpace(os.Getenv("HOWL_IMAGE_RECLAIM_TEST_GROWTH_MB"))
	if raw == "" {
		return 0
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || v <= 0 {
		return 0
	}
	return v
}

func nonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
