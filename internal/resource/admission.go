package resource

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"howl-chat/internal/backend/gguf"
	backendtypes "howl-chat/internal/backend/types"
)

// DefaultReservePercent is the baseline headroom kept free for the operating
// system, drivers, background work and allocation fluctuation. It is not a
// universal hard rule: a caller may supply ReserveOverrideMB when a measured,
// device-specific requirement is known (sd-server's --max-vram reserve is the
// current example).
const DefaultReservePercent = 0.10

// softReserveAllowanceMB lets a request use part of the default 10% reserve
// when it is the difference between running and refusing. The reserve is still
// the target: admission plans tiling and CPU placement against it first, and
// only dips into this allowance for the final fit. A request that would use
// more than this much of the reserve is refused. Explicit reserves are never
// softened.
const softReserveAllowanceMB = 500

// hardReserveFloorMB is the minimum headroom kept free even when the soft
// allowance is used. It prevents a small reserve on an almost-full card from
// being spent entirely by the allowance.
const hardReserveFloorMB = 128

const bytesPerMB = int64(1024 * 1024)

// AdmissionRequest is the complete workload a caller wants to start. The
// authority owns the estimate; callers only state the model and what it will do.
type AdmissionRequest struct {
	Kind ResourceType
	Name string
	Path string

	// Chat workload.
	Arch        *backendtypes.ArchitectureSpec
	ContextSize int
	GPULayers   int
	BatchSize   int

	// Image workload.
	Width      int
	Height     int
	BatchCount int
	Steps      int
	Family     string
	Mode       string
	// Component sizes are measured file sizes where available.
	//
	// CheckpointBytes is the diffusion checkpoint / DiT file sd-server opens, and
	// ConditionerBytes is the text-encoder set it loads from separate files
	// (CLIP-L/G, T5, LLM). Both are resident weights, not optional extras: an
	// estimate that omitted them predicted a 6617 MB checkpoint as 2189 MB and
	// told the user a chat model could safely coexist with it.
	CheckpointBytes  int64
	ConditionerBytes int64
	ControlNetBytes  int64
	LoRABytes        int64
	VAEBytes         int64
	// UpscalerBytes is an image-upscaler model loaded by a workflow stage. It is
	// a resident adapter in the same sense as ControlNet: the engine may keep it
	// between stages, so admission must reserve for it before the stage runs.
	UpscalerBytes int64
	// IPAdapterBytes is the IP-Adapter weight set. It is a resident adapter,
	// separate from ControlNet and Upscaler; charge it once per engine start.
	IPAdapterBytes int64
	// ClipVisionBytes is the CLIP-Vision encoder used by IP-Adapter. It is a
	// conditioner module; ClipVisionCPU keeps it in RAM/CPU instead of VRAM.
	ClipVisionBytes int64
	ClipVisionCPU   bool
	// VAETiling says the engine will be started with --vae-tiling, so the VAE
	// working set is the measured tiled buffer instead of the full decode.
	VAETiling bool
	// ControlNetCPU says the engine will compute the ControlNet on CPU. Its
	// params then stay in system RAM instead of being staged to CUDA for the
	// sampling phase.
	ControlNetCPU bool
	// ResidentLoaded marks a workflow-stage delta against a model that is
	// already admitted and resident. In that case the fixed module reserve is
	// already held; only new adapters and the stage working set are charged.
	ResidentLoaded bool
	// WeightType is sd-server's --type. File bytes are converted to this type
	// while loading, so admission sizes the converted weights, not the file.
	WeightType string

	// Audio workload.
	AudioVRAMMB int64
	AudioRAMMB  int64

	// RequireGPU means the start asks for GPU placement. When the fitted plan
	// cannot keep any layer on the GPU, the plan is CPU-only instead of a
	// refusal (for a manual chat start).
	RequireGPU bool
	// AllowManual lets a chat load proceed through the fitted/manual fallback
	// when the automatic plan cannot fit or telemetry cannot be established.
	// It is set only by the user-initiated load path; preflight and handoff
	// checks leave it false so they can report the true blocker.
	AllowManual bool
	// WeightsOffloaded means the start asks the engine to keep the weights in
	// system memory (sd-server --offload-to-cpu, --params-backend cpu, or disk)
	// and stream them for compute. Then the device holds only the compute working
	// set, so the weight bytes belong to the RAM requirement, not the VRAM
	// requirement. Without this the app refused every model bigger than the card
	// even though the engine was already being told to offload it.
	WeightsOffloaded bool
	// CPUOnly means --backend cpu: there is no device compute at all, so the
	// request needs RAM for weights and activations and no VRAM. It is explicit
	// user intent, never selected automatically.
	CPUOnly bool
	// DeviceIndex selects the target GPU. Zero is the first device.
	DeviceIndex int

	// KVCacheBytesPerElement is the cost of one KV element under the cache type
	// the engine will be told to build (--cache-type-k/v). Zero keeps the
	// estimator's own f16 default, so an unknown type can only over-reserve.
	KVCacheBytesPerElement int

	// ReserveOverrideMB replaces the default percentage only when a concrete,
	// measured reserve applies. sd-server's --max-vram reserve is the current
	// model-specific case.
	ReserveOverrideMB int64
	// ReserveExplicit marks ReserveOverrideMB as deliberate, including an
	// explicit zero (sd-server --max-vram 0 reserves nothing). Without it a zero
	// override falls back to DefaultReservePercent.
	ReserveExplicit bool
	// VRAMBudgetMB is a hard ceiling on usable GPU memory, the same shape as a
	// positive sd-server --max-vram. Zero means no explicit cap.
	VRAMBudgetMB int64
	// MinimumCommitmentVRAMMB is a previously observed peak for this model at
	// an equal or larger workload shape. Admission may raise the commitment to
	// it, never lower it, and never treats it as the resident requirement: a
	// post-generation observation already contains activation and allocator
	// retention.
	MinimumCommitmentVRAMMB int64
	// MeasuredResidentVRAMMB is a previously observed resident total for the
	// same model and component set. Zero means no measurement is available.
	MeasuredResidentVRAMMB int64

	// ReplacingKind/Name let a preflight account for a model that the caller
	// will stop before spawning the replacement. The final admission before
	// spawning must not use this.
	ReplacingKind ResourceType
	ReplacingName string
}

// AdmissionEstimate carries every number used in the decision, so a refusal can
// name the exact requirement, free value, reserve and available value.
type AdmissionEstimate struct {
	RequiredVRAMMB int64
	// ActivationVRAMMB is the estimated peak decode/execution working set for
	// the requested workload. It is separate from RequiredVRAMMB (resident
	// weights) so a caller can decide whether the engine should tile instead of
	// refusing a load whose weights fit.
	ActivationVRAMMB int64
	// ResidentVRAMMB is the converted model set expected to stay on the device
	// while the workload runs. ProjectedPeakVRAMMB adds the selected activation
	// and the fixed engine overhead to it; it is the number compared with
	// AvailableVRAMMB when the workload is an image stage.
	ResidentVRAMMB      int64
	ProjectedPeakVRAMMB int64
	// CommitmentFloorMB is a previously measured peak that raised the
	// commitment above the projection. Source names the evidence.
	CommitmentFloorMB     int64
	CommitmentFloorSource string
	// TilingPlanned reports that VAE tiling is required to keep the projected
	// peak inside the available budget. The caller must start the engine with
	// --vae-tiling for this decision to hold.
	TilingPlanned bool
	// ControlNetVRAMMB is the measured ControlNet params set the pinned engine
	// would stage to CUDA during sampling. ControlNetCPUPlanned reports that
	// admission keeps it on CPU because the staged set would cross the budget.
	ControlNetVRAMMB     int64
	ControlNetCPUPlanned bool
	RequiredRAMMB        int64
	ReserveVRAMMB        int64
	// SoftReserveAllowanceMB is the part of the default reserve this decision
	// may use without refusing (zero for an explicit reserve).
	SoftReserveAllowanceMB int64
	// SoftReserveUsedMB is how much of that allowance the projected peak needs.
	SoftReserveUsedMB int64
	ReserveRAMMB      int64
	FreeVRAMMB        int64
	FreeRAMMB         int64
	AvailableVRAMMB   int64
	AvailableRAMMB    int64
	TotalVRAMMB       int64
	TotalRAMMB        int64
	GPUName           string
	DeviceIndex       int
}

// AdmissionAdjustment is the plan the authority fitted when the requested
// settings exceeded the measured budget, or the plan it selected when the
// request asked for automatic placement. It is applied to the spawned engine so
// admission and the command line describe the same workload.
type AdmissionAdjustment struct {
	ContextSize int `json:"context_size"`
	GPULayers   int `json:"gpu_layers"`
	BatchSize   int `json:"batch_size"`
	// Auto is true when GPULayers came from the upward placement search rather
	// than from a clamp of a layer count the caller supplied.
	Auto  bool     `json:"auto"`
	Notes []string `json:"notes,omitempty"`
}

// AdmissionDecision is the authority's answer. Blocking contains the loaded
// models that make the request impossible, largest GPU footprint first.
type AdmissionDecision struct {
	Admitted bool
	Reason   string
	Estimate AdmissionEstimate
	Blocking []ModelState
	// Adjustment is non-nil when the authority fitted a smaller workload.
	Adjustment *AdmissionAdjustment
	// Manual is true when the load proceeds outside a verified fit, because
	// the user asked to load and no automatic plan could be established.
	Manual bool
	// Warning carries the user-facing numbers for an adjusted or manual plan.
	Warning string
}

// AdmissionLease holds an admitted allocation until the caller either registers
// the model after a confirmed start or releases the lease after a failure. One
// lease prevents two concurrent starts from both passing against the same free
// value.
type AdmissionLease struct {
	authority *Authority
	id        uint64
}

// Release returns the reservation to the ledger. It is idempotent.
func (l *AdmissionLease) Release() {
	if l == nil || l.authority == nil || l.id == 0 {
		return
	}
	l.authority.release(l.id)
	l.id = 0
}

type inFlightAllocation struct {
	id   uint64
	kind ResourceType
	name string
	vram int64
	ram  int64
}

// Authority is the one admission decision point. It refreshes hardware on every
// decision, estimates the requested workload, applies the reserve, and records
// in-flight allocations until their leases are released.
type Authority struct {
	mu       sync.Mutex
	hardware func() (HardwareSnapshot, error)
	ledger   func() []ModelState
	nextID   uint64
	inFlight []inFlightAllocation
}

// NewAuthority builds an authority with injectable hardware and ledger
// providers. Production uses DefaultAuthority.
func NewAuthority(hardware func() (HardwareSnapshot, error), ledger func() []ModelState) *Authority {
	return &Authority{hardware: hardware, ledger: ledger}
}

var defaultAuthority = NewAuthority(
	func() (HardwareSnapshot, error) { return GetManager().RefreshSnapshot() },
	func() []ModelState { return GetManager().LoadedModels() },
)

// DefaultAuthority returns the process-wide admission authority.
func DefaultAuthority() *Authority { return defaultAuthority }

// Admit refreshes hardware, estimates the request, and returns a decision and a
// lease. A refused request never returns a lease.
func (a *Authority) Admit(req AdmissionRequest) (AdmissionDecision, *AdmissionLease) {
	if a == nil {
		return AdmissionDecision{Reason: "admission authority is not configured; refusing to start a model"}, nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	hardware := a.hardware
	if hardware == nil {
		hardware = func() (HardwareSnapshot, error) { return GetManager().RefreshSnapshot() }
	}
	snap, err := hardware()
	if err != nil {
		return a.refuseOrManual(req, snap, AdmissionEstimate{}, nil,
			"hardware detection failed: "+err.Error()+"; no measured free-memory value is available")
	}
	if snap.DetectionError != "" {
		return a.refuseOrManual(req, snap, AdmissionEstimate{}, nil,
			"hardware detection failed: "+snap.DetectionError+"; no measured free-memory value is available")
	}

	requiredVRAM, requiredRAM, activationVRAM, notes, err := estimateAdmission(req)
	if err != nil {
		return refused(req, AdmissionEstimate{}, nil, err.Error()), nil
	}

	loaded := a.loadedLocked()
	blocking := blockingModels(loaded, req, requiredRAM > 0)

	totalVRAM, usedVRAM, freeVRAM, gpuName, est, err := targetMemory(snap, req.DeviceIndex)
	if err != nil {
		return refused(req, est, blocking, err.Error()), nil
	}
	_ = usedVRAM

	replaceVRAM, replaceRAM := replacementFootprint(loaded, req)
	pendingVRAM, pendingRAM := a.pendingLocked()

	effectiveFreeVRAMMB := freeVRAM/bytesPerMB + replaceVRAM - pendingVRAM
	if effectiveFreeVRAMMB < 0 {
		effectiveFreeVRAMMB = 0
	}
	// A replaced model's registered footprint is an estimate; when the old
	// process is still alive, free memory already excludes it and adding the
	// whole estimate can claim more free VRAM than the card has. Never let the
	// post-reclamation estimate exceed the device total (observed 8693 MB on an
	// 8192 MB card during a blank-Max-VRAM restart).
	if totalVRAMMB := totalVRAM / bytesPerMB; totalVRAMMB > 0 && effectiveFreeVRAMMB > totalVRAMMB {
		effectiveFreeVRAMMB = totalVRAMMB
	}
	effectiveFreeRAMMB := snap.FreeRAMBytes/bytesPerMB + replaceRAM - pendingRAM
	if effectiveFreeRAMMB < 0 {
		effectiveFreeRAMMB = 0
	}

	reserveVRAM := req.ReserveOverrideMB
	if reserveVRAM < 0 {
		reserveVRAM = 0
	}
	if reserveVRAM == 0 && (requiredVRAM > 0 || req.ResidentLoaded) && !req.ReserveExplicit {
		reserveVRAM = int64(float64(effectiveFreeVRAMMB) * DefaultReservePercent)
		if req.Kind == ResourceImage {
			// The image reclaimer already holds a fixed 10%-of-device floor
			// (819 MB on the 8 GiB card). Deriving the admission reserve only
			// from the moment's free value shrank it exactly when memory was
			// tight (26 MB on a 265 MB-free card), letting admission spend the
			// headroom the reclaimer is meant to protect. Image stages must use
			// the stronger of the two; chat and audio keep the existing rule.
			if totalVRAMMB := totalVRAM / bytesPerMB; totalVRAMMB > 0 {
				if deviceReserve := int64(float64(totalVRAMMB) * DefaultReservePercent); deviceReserve > reserveVRAM {
					reserveVRAM = deviceReserve
				}
			}
		}
	}
	availableVRAM := effectiveFreeVRAMMB - reserveVRAM
	if availableVRAM < 0 {
		availableVRAM = 0
	}
	// MeasuredResidentVRAMMB is a previously observed resident set for the same
	// model and component composition. It lets a budget change reuse measurement
	// instead of charging the pessimistic module+reserve projection again; it is
	// not a user override and never replaces the reserve or the live free value.
	if req.MeasuredResidentVRAMMB > 0 && !req.CPUOnly && !req.WeightsOffloaded {
		requiredVRAM = req.MeasuredResidentVRAMMB
	}

	// A positive --max-vram is a ceiling, not a reservation: sd-server will not
	// spend more than that even when the device has more free memory.
	if req.VRAMBudgetMB > 0 && availableVRAM > req.VRAMBudgetMB {
		availableVRAM = req.VRAMBudgetMB
	}

	// Soften only the default percentage reserve. An explicit measured reserve
	// (ReserveOverrideMB) is the caller's hard requirement and stays intact.
	softAllowance := int64(0)
	softAvailableVRAM := availableVRAM
	if req.ReserveOverrideMB == 0 && !req.ReserveExplicit && reserveVRAM > 0 {
		softReserve := reserveVRAM - softReserveAllowanceMB
		if softReserve < hardReserveFloorMB {
			softReserve = hardReserveFloorMB
		}
		if softReserve > reserveVRAM {
			softReserve = reserveVRAM
		}
		softAllowance = reserveVRAM - softReserve
		softAvailableVRAM = effectiveFreeVRAMMB - softReserve
		if softAvailableVRAM < 0 {
			softAvailableVRAM = 0
		}
		if req.VRAMBudgetMB > 0 && softAvailableVRAM > req.VRAMBudgetMB {
			softAvailableVRAM = req.VRAMBudgetMB
		}
	}

	reserveRAM := int64(float64(effectiveFreeRAMMB) * DefaultReservePercent)
	availableRAM := effectiveFreeRAMMB - reserveRAM
	if availableRAM < 0 {
		availableRAM = 0
	}

	// gpu_layers -1 is automatic placement: choose the largest layer count that
	// fits the post-reserve budget with the shared estimator. Zero stays an
	// explicit CPU-only choice and never reaches this branch.
	var autoPlan ChatPlan
	autoPlanned := false
	if req.Kind == ResourceLLM && req.GPULayers < 0 && req.Arch != nil {
		ctxSize := req.ContextSize
		if ctxSize <= 0 {
			ctxSize = 4096
		}
		batchSize := req.BatchSize
		if batchSize <= 0 {
			batchSize = 512
		}
		autoPlan = PlanChatLayers(req.Arch, ctxSize, batchSize, req.KVCacheBytesPerElement, availableVRAM, availableRAM)
		req.GPULayers = autoPlan.GPULayers
		requiredVRAM = int64(autoPlan.Estimate.GPUMB)
		requiredRAM = int64(autoPlan.Estimate.RAMMB)
		if autoPlan.GPULayers == 0 {
			// A CPU plan has no CUDA context or device buffers.
			requiredVRAM = 0
		}
		est.RequiredVRAMMB = requiredVRAM
		est.RequiredRAMMB = requiredRAM
		autoPlanned = true
	}

	est.RequiredVRAMMB = requiredVRAM
	est.ActivationVRAMMB = activationVRAM
	// Image stages allocate the resident modules plus a decode/compute working
	// set. Judge that projected peak, not only the largest module: at 512x512 the
	// measured VAE decode alone is 1920 MB, and the old decode-only estimate let
	// the base stage cross the reserve. If the untiled peak does not fit, plan
	// VAE tiling (the measured tiled decode buffer is 480 MB) and admit only when
	// the tiled peak fits. If neither fits, the projection below refuses it.
	if req.Kind == ResourceImage && (requiredVRAM > 0 || req.ResidentLoaded) && totalVRAM > 0 && !req.CPUOnly && !req.WeightsOffloaded {
		resident := imageResidentVRAMMB(req)
		controlNet := imageControlNetVRAMMB(req)
		if req.MeasuredResidentVRAMMB > 0 && !req.ResidentLoaded {
			// The measured total already includes every module the engine kept.
			resident = req.MeasuredResidentVRAMMB
			controlNet = 0
		}
		if req.ResidentLoaded {
			// The base weights are already held; a delta stage only adds its own
			// adapter bytes and working set.
			resident = 0
			controlNet = 0
		} else if req.ControlNetCPU {
			controlNet = 0
		}
		activation := activationVRAM
		if req.VAETiling {
			activation = minInt64(activation, sdTiledActivationMB)
		} else {
			untiledPeak := resident + activation + sdFixedOverheadMB
			if untiledPeak > availableVRAM {
				tiledPeak := resident + sdTiledActivationMB + sdFixedOverheadMB
				// Prefer the smaller working set even when it only fits through
				// the soft allowance: tiling is the safer way to spend that space.
				if tiledPeak <= softAvailableVRAM {
					est.TilingPlanned = true
					activation = sdTiledActivationMB
				}
			}
		}
		projected := resident + controlNet + activation + sdFixedOverheadMB
		// A resident stage has no base weights in this projection. Its newly
		// added modules (for example an upscaler) are incremental capacity in
		// addition to the execution working set, not a replacement for it.
		if req.ResidentLoaded {
			projected += requiredVRAM
		}
		if controlNet > 0 && projected > availableVRAM {
			basePeak := resident + activation + sdFixedOverheadMB
			if basePeak <= softAvailableVRAM {
				est.ControlNetCPUPlanned = true
				projected = basePeak
			}
		}
		est.ResidentVRAMMB = resident
		est.ControlNetVRAMMB = imageControlNetVRAMMB(req)
		est.ActivationVRAMMB = activation
		est.ProjectedPeakVRAMMB = projected
	}
	est.RequiredRAMMB = requiredRAM
	est.ReserveVRAMMB = reserveVRAM
	est.SoftReserveAllowanceMB = softAllowance
	est.ReserveRAMMB = reserveRAM
	est.FreeVRAMMB = effectiveFreeVRAMMB
	est.FreeRAMMB = effectiveFreeRAMMB
	est.AvailableVRAMMB = availableVRAM
	est.AvailableRAMMB = availableRAM
	est.TotalVRAMMB = totalVRAM / bytesPerMB
	est.TotalRAMMB = snap.TotalRAMBytes / bytesPerMB
	est.GPUName = gpuName
	est.DeviceIndex = req.DeviceIndex
	_ = notes

	if requiredVRAM > 0 {
		if totalVRAM <= 0 {
			if req.RequireGPU {
				reason := fmt.Sprintf("cannot start %s %q: it needs %d MB of GPU memory but no GPU was detected", req.Kind, req.Name, requiredVRAM)
				return a.refuseOrManual(req, snap, est, blocking, reason)
			}
			// A CPU-capable workload may still declare a GPU profile; without a
			// measured device there is no GPU allocation to check.
			requiredVRAM = 0
			est.RequiredVRAMMB = 0
			reserveVRAM = 0
			est.ReserveVRAMMB = 0
			est.AvailableVRAMMB = 0
		}
		if totalVRAM > 0 && !snap.VRAMLive {
			source := strings.TrimSpace(snap.VRAMSource)
			if source == "" {
				source = "unknown"
			}
			return a.refuseOrManual(req, snap, est, blocking,
				fmt.Sprintf("cannot start %s %q: GPU memory use could not be read live (source: %s)", req.Kind, req.Name, source))
		}
	}

	// A measured peak for this model at an equal or larger workload shape is
	// stronger evidence than the metadata projection. It can only raise the
	// commitment: the observation already contains activation and allocator
	// retention, so it may not replace the resident requirement or lower it.
	if req.MinimumCommitmentVRAMMB > 0 && !req.CPUOnly && !req.WeightsOffloaded {
		if measured := req.MinimumCommitmentVRAMMB; measured > est.committedVRAMMB() {
			est.CommitmentFloorMB = measured
			est.CommitmentFloorSource = CoveragePostGenerationDelta
		}
	}

	needVRAM := est.committedVRAMMB()
	if autoPlanned {
		// The search already proved the plan fits, or proved nothing fits. A
		// zero-layer plan next to a loaded model is the silent CPU demotion this
		// milestone exists to remove: refuse and name the blocker instead.
		if autoPlan.Fits && autoPlan.GPULayers > 0 {
			warning := fmt.Sprintf("automatic placement: %d of %d layers fit %d MB of post-reserve VRAM (estimate %d MB VRAM / %d MB RAM, KV %d bytes per element).",
				autoPlan.GPULayers, req.Arch.Layers, availableVRAM, requiredVRAM, requiredRAM, req.KVCacheBytesPerElement)
			return a.grant(req, est, blocking, &AdmissionAdjustment{
				ContextSize: req.ContextSize,
				GPULayers:   autoPlan.GPULayers,
				BatchSize:   req.BatchSize,
				Auto:        true,
				Notes:       autoPlan.Notes,
			}, warning, false)
		}
		if len(blocking) > 0 {
			return refused(req, est, blocking, fmt.Sprintf(
				"automatic placement found no GPU space while %s is loaded; unload it or set GPU layers to 0 explicitly to run on the CPU",
				blocking[0].Name)), nil
		}
		if autoPlan.Fits && autoPlan.GPULayers == 0 {
			warning := fmt.Sprintf("automatic placement: no layer fits the %d MB post-reserve VRAM; starting on the CPU (estimate %d MB RAM). Set GPU layers explicitly to override.",
				availableVRAM, requiredRAM)
			return a.grant(req, est, blocking, &AdmissionAdjustment{
				ContextSize: req.ContextSize,
				GPULayers:   0,
				BatchSize:   req.BatchSize,
				Auto:        true,
				Notes:       autoPlan.Notes,
			}, warning, false)
		}
		if req.AllowManual {
			warning := fmt.Sprintf("automatic placement could not fit %d MB VRAM / %d MB RAM after reserve; starting with a reduced CPU plan because the model load was requested.", availableVRAM, availableRAM)
			return a.manualChat(req, snap, est, blocking, warning)
		}
		return refused(req, est, blocking, fmt.Sprintf(
			"automatic placement found no plan that fits %d MB post-reserve VRAM and %d MB post-reserve RAM",
			availableVRAM, availableRAM)), nil
	}
	if needVRAM > availableVRAM || requiredRAM > availableRAM {
		// A lone chat load is fitted rather than refused: the unused clamp
		// ladder is the one authority for the reduced context / layer / batch
		// plan, and the caller receives the numbers that produced it.
		if req.Kind == ResourceLLM && len(blocking) == 0 && req.Arch != nil {
			adjusted, adjVRAM, adjRAM, notes, fits := fitChatRequest(req, availableVRAM, availableRAM)
			if fits {
				warning := fittedWarning(req, requiredVRAM, requiredRAM, adjusted, adjVRAM, adjRAM, availableVRAM, availableRAM, notes)
				est.RequiredVRAMMB = adjVRAM
				est.RequiredRAMMB = adjRAM
				return a.grant(adjusted, est, blocking, &AdmissionAdjustment{
					ContextSize: adjusted.ContextSize,
					GPULayers:   adjusted.GPULayers,
					BatchSize:   adjusted.BatchSize,
					Notes:       notes,
				}, warning, false)
			}
			if req.AllowManual {
				warning := fmt.Sprintf("requested context %d / %d GPU layers / batch %d needs %d MB VRAM and %d MB RAM; only %d MB VRAM and %d MB RAM are available after reserve, and no fitted plan fits. Starting with a reduced CPU/mixed plan because the model load was requested.",
					req.ContextSize, req.GPULayers, req.BatchSize, requiredVRAM, requiredRAM, availableVRAM, availableRAM)
				return a.manualChat(req, snap, est, blocking, warning)
			}
		}
		// If tiling and reduced plans could not bring the stage under the
		// preferred reserve, the bounded soft allowance is the last fit. It is
		// never used for RAM, and it cannot cross the hard headroom floor.
		if needVRAM > availableVRAM && requiredRAM <= availableRAM && needVRAM <= softAvailableVRAM {
			est.SoftReserveUsedMB = needVRAM - availableVRAM
			warning := fmt.Sprintf("using %d MB of the %d MB soft reserve allowance; %d MB remains free above the %d MB hard floor",
				est.SoftReserveUsedMB, est.SoftReserveAllowanceMB, effectiveFreeVRAMMB-needVRAM, hardReserveFloorMB)
			return a.grant(req, est, blocking, nil, warning, false)
		}
		if needVRAM > availableVRAM {
			return refused(req, est, blocking, vramReason(req, est, blocking)), nil
		}
		return refused(req, est, blocking, ramReason(req, est, blocking)), nil
	}

	return a.grant(req, est, blocking, nil, "", false)
}

// grant records the admitted allocation and returns its lease.
func (a *Authority) grant(req AdmissionRequest, est AdmissionEstimate, blocking []ModelState, adjustment *AdmissionAdjustment, warning string, manual bool) (AdmissionDecision, *AdmissionLease) {
	a.nextID++
	id := a.nextID
	committedVRAM := est.committedVRAMMB()
	a.inFlight = append(a.inFlight, inFlightAllocation{
		id:   id,
		kind: req.Kind,
		name: req.Name,
		vram: committedVRAM,
		ram:  est.RequiredRAMMB,
	})

	reason := fmt.Sprintf("admitted: %d MB VRAM required, %d MB available after %d MB reserve", committedVRAM, est.AvailableVRAMMB, est.ReserveVRAMMB)
	if est.ProjectedPeakVRAMMB > 0 && est.ProjectedPeakVRAMMB >= est.RequiredVRAMMB {
		if req.ResidentLoaded {
			reason = fmt.Sprintf("admitted: projected stage peak %d MB VRAM (%d MB new modules + activation %d MB + %d MB engine overhead), %d MB available after %d MB reserve",
				est.ProjectedPeakVRAMMB, est.RequiredVRAMMB, est.ActivationVRAMMB, sdFixedOverheadMB, est.AvailableVRAMMB, est.ReserveVRAMMB)
		} else {
			reason = fmt.Sprintf("admitted: projected peak %d MB VRAM (resident %d MB + activation %d MB + %d MB engine overhead), %d MB available after %d MB reserve",
				est.ProjectedPeakVRAMMB, est.ResidentVRAMMB, est.ActivationVRAMMB, sdFixedOverheadMB, est.AvailableVRAMMB, est.ReserveVRAMMB)
		}
		if est.TilingPlanned {
			reason += "; VAE tiling planned"
		}
		if est.ControlNetCPUPlanned {
			reason += "; ControlNet CPU placement planned"
		}
	}
	if adjustment != nil {
		reason += fmt.Sprintf("; fitted plan context=%d gpu_layers=%d batch=%d", adjustment.ContextSize, adjustment.GPULayers, adjustment.BatchSize)
	}
	if est.CommitmentFloorMB > 0 {
		reason += fmt.Sprintf("; measured floor %d MB (%s) raised the commitment above the %d MB projection",
			est.CommitmentFloorMB, est.CommitmentFloorSource, est.ProjectedPeakVRAMMB)
	}
	if est.SoftReserveUsedMB > 0 {
		need := est.committedVRAMMB()
		reason += fmt.Sprintf("; soft reserve: used %d MB of the %d MB allowance, keeping %d MB free above the %d MB hard floor",
			est.SoftReserveUsedMB, est.SoftReserveAllowanceMB, est.FreeVRAMMB-need, hardReserveFloorMB)
	}
	if manual {
		reason = "manual fallback: " + warning
	}
	return AdmissionDecision{
		Admitted:   true,
		Reason:     reason,
		Estimate:   est,
		Blocking:   blocking,
		Adjustment: adjustment,
		Manual:     manual,
		Warning:    warning,
	}, &AdmissionLease{authority: a, id: id}
}

// refuseOrManual keeps the historical refusal for every workload except a
// user-initiated chat load, which falls back to a reduced CPU/mixed plan and
// reports why instead of locking the user out.
func (a *Authority) refuseOrManual(req AdmissionRequest, snap HardwareSnapshot, est AdmissionEstimate, blocking []ModelState, reason string) (AdmissionDecision, *AdmissionLease) {
	if req.Kind == ResourceLLM && req.AllowManual && len(blocking) == 0 {
		return a.manualChat(req, snap, est, blocking, reason)
	}
	return refused(req, est, blocking, reason), nil
}

// manualChat builds a CPU-first reduced plan. It is the last rung of the
// fallback ladder, used only when the user asked for the load and no verified
// fit exists. The warning always carries the measured numbers.
func (a *Authority) manualChat(req AdmissionRequest, snap HardwareSnapshot, est AdmissionEstimate, blocking []ModelState, warning string) (AdmissionDecision, *AdmissionLease) {
	if req.Arch == nil {
		return refused(req, est, blocking, warning), nil
	}
	freeRAMMB := snap.FreeRAMBytes / bytesPerMB
	reserveRAM := int64(float64(freeRAMMB) * DefaultReservePercent)
	availableRAM := freeRAMMB - reserveRAM
	if availableRAM < 0 {
		availableRAM = 0
	}
	adjusted, vram, ram, notes, _ := fitChatRequest(req, 0, availableRAM)

	est.TotalVRAMMB = snap.TotalVRAMBytes / bytesPerMB
	est.FreeVRAMMB = snap.FreeVRAMBytes / bytesPerMB
	est.AvailableVRAMMB = 0
	est.TotalRAMMB = snap.TotalRAMBytes / bytesPerMB
	est.FreeRAMMB = freeRAMMB
	est.ReserveRAMMB = reserveRAM
	est.AvailableRAMMB = availableRAM
	est.RequiredVRAMMB = vram
	est.RequiredRAMMB = ram

	if len(notes) > 0 {
		warning += " " + strings.Join(notes, "; ")
	}
	return a.grant(adjusted, est, blocking, &AdmissionAdjustment{
		ContextSize: adjusted.ContextSize,
		GPULayers:   adjusted.GPULayers,
		BatchSize:   adjusted.BatchSize,
		Notes:       notes,
	}, warning, true)
}

// fitChatRequest runs the existing estimator's clamp ladder against the budget
// that remains after the reserve. It reports whether the fitted plan fits; the
// caller decides whether to admit it or to fall back to the manual rung.
func fitChatRequest(req AdmissionRequest, availableVRAMMB, availableRAMMB int64) (AdmissionRequest, int64, int64, []string, bool) {
	if req.Arch == nil {
		return req, 0, 0, nil, false
	}
	if availableVRAMMB < 0 {
		availableVRAMMB = 0
	}
	if availableRAMMB < 0 {
		availableRAMMB = 0
	}
	ctx := req.ContextSize
	if ctx <= 0 {
		ctx = 4096
	}
	layers := req.GPULayers
	batch := req.BatchSize

	model := gguf.NewMemoryModel(req.Arch)
	newCtx, newLayers, newBatch, _, notes := model.ClampToHardware(ctx, layers, batch, availableVRAMMB, availableRAMMB)
	if newCtx <= 0 {
		newCtx = ctx
	}

	adjusted := req
	adjusted.ContextSize = newCtx
	adjusted.GPULayers = newLayers
	adjusted.BatchSize = newBatch
	adjusted.RequireGPU = newLayers > 0

	vram, ram, _, _, err := estimateChat(adjusted)
	if err != nil {
		return req, 0, 0, nil, false
	}
	return adjusted, vram, ram, notes, vram <= availableVRAMMB && ram <= availableRAMMB
}

func fittedWarning(original AdmissionRequest, requiredVRAM, requiredRAM int64, adjusted AdmissionRequest, adjVRAM, adjRAM, availableVRAM, availableRAM int64, notes []string) string {
	warning := fmt.Sprintf("requested context %d / %d GPU layers / batch %d needs %d MB VRAM and %d MB RAM; %d MB VRAM and %d MB RAM are available after reserve. Loaded with context %d / %d GPU layers / batch %d (%d MB VRAM, %d MB RAM).",
		original.ContextSize, original.GPULayers, original.BatchSize, requiredVRAM, requiredRAM,
		availableVRAM, availableRAM, adjusted.ContextSize, adjusted.GPULayers, adjusted.BatchSize, adjVRAM, adjRAM)
	if len(notes) > 0 {
		warning += " " + strings.Join(notes, "; ")
	}
	return warning
}

func (a *Authority) release(id uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := a.inFlight[:0]
	for _, pending := range a.inFlight {
		if pending.id != id {
			out = append(out, pending)
		}
	}
	a.inFlight = out
}

func (a *Authority) loadedLocked() []ModelState {
	if a.ledger == nil {
		return nil
	}
	return a.ledger()
}

func (a *Authority) pendingLocked() (vram, ram int64) {
	for _, pending := range a.inFlight {
		vram += pending.vram
		ram += pending.ram
	}
	return vram, ram
}

// committedVRAMMB is the amount this admission path approved. A full image
// load's projected peak is the total resident + activation + overhead; a
// resident stage's projected peak is its incremental execution peak. Chat and
// audio have no projection and commit RequiredVRAMMB. The commitment is never
// Required + Projected, and a replacement credit is applied to effective free
// separately rather than subtracted from this commitment.
func (est AdmissionEstimate) committedVRAMMB() int64 { return est.CommittedVRAMMB() }

// CommittedVRAMMB is the amount this admission approved. A full image load's
// projected peak is the total resident + activation + overhead; a resident
// stage's projected peak is its incremental execution peak. Chat and audio have
// no projection and commit RequiredVRAMMB. The commitment is never Required +
// Projected, and a replacement credit is applied to effective free separately
// rather than subtracted from this commitment.
func (est AdmissionEstimate) CommittedVRAMMB() int64 {
	committed := est.RequiredVRAMMB
	if est.ProjectedPeakVRAMMB > committed {
		committed = est.ProjectedPeakVRAMMB
	}
	if est.CommitmentFloorMB > committed {
		committed = est.CommitmentFloorMB
	}
	return committed
}

func replacementFootprint(loaded []ModelState, req AdmissionRequest) (vram, ram int64) {
	if req.ReplacingKind == "" {
		return 0, 0
	}
	for _, model := range loaded {
		if model.Type != req.ReplacingKind {
			continue
		}
		if req.ReplacingName != "" && model.Name != req.ReplacingName {
			continue
		}
		vram += model.GPUMemoryMB
		ram += model.SystemRAMMB
	}
	return vram, ram
}

// blockingModels returns the resident models relevant to this decision. A
// GPU-only request keeps the historical GPU-resident list. A request with a
// RAM requirement also includes RAM-only residents, which otherwise disappear
// because their GPUMemoryMB is zero.
func blockingModels(loaded []ModelState, req AdmissionRequest, includeRAMResidents bool) []ModelState {
	out := make([]ModelState, 0, len(loaded))
	for _, model := range loaded {
		if model.Type == req.ReplacingKind && (req.ReplacingName == "" || model.Name == req.ReplacingName) {
			continue
		}
		if model.GPUMemoryMB <= 0 && (!includeRAMResidents || model.SystemRAMMB <= 0) {
			continue
		}
		out = append(out, model)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].GPUMemoryMB != out[j].GPUMemoryMB {
			return out[i].GPUMemoryMB > out[j].GPUMemoryMB
		}
		if out[i].SystemRAMMB != out[j].SystemRAMMB {
			return out[i].SystemRAMMB > out[j].SystemRAMMB
		}
		return out[i].Name < out[j].Name
	})
	return out
}

func targetMemory(snap HardwareSnapshot, index int) (total, used, free int64, name string, est AdmissionEstimate, err error) {
	if len(snap.GPUs) == 0 {
		return snap.TotalVRAMBytes, snap.UsedVRAMBytes, snap.FreeVRAMBytes, snap.GPUModel, AdmissionEstimate{}, nil
	}
	if index < 0 {
		index = 0
	}
	if index >= len(snap.GPUs) {
		return 0, 0, 0, "", AdmissionEstimate{}, fmt.Errorf("GPU device %d was requested but only %d device(s) were detected", index, len(snap.GPUs))
	}
	gpu := snap.GPUs[index]
	return gpu.TotalVRAMBytes, gpu.UsedVRAMBytes, gpu.FreeVRAMBytes, gpu.Name, AdmissionEstimate{DeviceIndex: index}, nil
}

// ImageResidentNeedMB reports the VRAM a resident image model needs plus the
// workload's decode working set, without crediting CPU offload. Howl uses it to
// decide whether a configured budget is too small to honor with opportunistic
// offload; the answer is derived from the same component-aware estimate as
// admission, not from a UI constant.
func ImageResidentNeedMB(req AdmissionRequest) (weights, activation int64, ok bool) {
	if req.Kind != ResourceImage {
		return 0, 0, false
	}
	req.WeightsOffloaded = false
	req.CPUOnly = false
	weights, _, activation, _, err := estimateImage(req)
	if err != nil || weights <= 0 {
		return 0, 0, false
	}
	return weights, activation, true
}

func estimateAdmission(req AdmissionRequest) (vram, ram, activation int64, notes []string, err error) {
	switch req.Kind {
	case ResourceLLM:
		return estimateChat(req)
	case ResourceImage:
		return estimateImage(req)
	case ResourceASR, ResourceTTS:
		return req.AudioVRAMMB, req.AudioRAMMB, 0, nil, nil
	default:
		return 0, 0, 0, nil, fmt.Errorf("unknown model kind %q", req.Kind)
	}
}

func estimateChat(req AdmissionRequest) (vram, ram, activation int64, notes []string, err error) {
	if req.Arch == nil {
		return 0, 0, 0, nil, fmt.Errorf("cannot estimate chat model %q: architecture metadata is unavailable; refusing to load without a memory estimate", req.Name)
	}
	ctx := req.ContextSize
	if ctx <= 0 {
		ctx = 4096
	}
	batch := req.BatchSize
	if batch <= 0 {
		batch = 512
	}
	layers := req.GPULayers
	if layers < 0 {
		layers = req.Arch.Layers
	}
	// Cost the KV cache at the type the engine will build. The shared profile is
	// never mutated; a copy carries the caller's cache type.
	spec := req.Arch
	if req.KVCacheBytesPerElement > 0 && spec.KVCacheBytesPerElement != req.KVCacheBytesPerElement {
		clone := *spec
		clone.KVCacheBytesPerElement = req.KVCacheBytesPerElement
		spec = &clone
	}
	est := gguf.NewMemoryModel(spec).Estimate(ctx, layers, batch)
	vram = est.GPUMB
	ram = est.RAMMB
	if layers == 0 {
		// A CPU-only load does not need a CUDA context or GPU buffers.
		vram = 0
	}
	return vram, ram, 0, []string{est.Notes}, nil
}

// sd-server's own placement rule, read from its backend_fit log on this machine:
// each module is judged on its own weights plus a fixed runner reserve against
// the --max-vram budget (2048 MB for the DiT and the text encoders, 1024 MB for
// the VAE). Admission uses the same rule so the policy authority and the
// placement authority (sd.cpp backend_fit) cannot disagree about whether the
// largest module can be resident in VRAM.
const (
	sdModuleReserveMB = int64(2048)
	sdVAEReserveMB    = int64(1024)
	// ImageModuleReserveMB exposes the same per-module reserve to the image
	// placement planner. One rule, so planning and admission cannot disagree.
	ImageModuleReserveMB = sdModuleReserveMB
)

func estimateImage(req AdmissionRequest) (vram, ram, activation int64, notes []string, err error) {
	activation = imageDecodeActivationMB(req.Width, req.Height, req.BatchCount)
	mode := strings.ToLower(req.Mode)
	if strings.Contains(mode, "dgl") || (strings.Contains(mode, "upscale") && !strings.Contains(mode, "latent")) {
		// ESRGAN upscaling (DGL or the native route) uses the upscaler graph,
		// not the VAE decode/sampling pass.
		activation = 0
	}
	weights := req.CheckpointBytes + req.ConditionerBytes + req.VAEBytes + req.ControlNetBytes + req.LoRABytes + req.UpscalerBytes + req.IPAdapterBytes + req.ClipVisionBytes
	if req.CPUOnly {
		ramMB := weights/bytesPerMB + 500 + activation
		return 0, ramMB, activation, []string{
			fmt.Sprintf("image estimate (CPU compute): %d MB of weights and %d MB decode activation in RAM, no VRAM", weights/bytesPerMB, activation),
		}, nil
	}
	// A CPU-offloaded start keeps every weight module in system memory. sd-server
	// then holds only the compute working set on the device, so the estimate must
	// charge the weights to RAM and leave VRAM at that reserve.
	if req.WeightsOffloaded {
		ramMB := weights/bytesPerMB + 500
		return sdModuleReserveMB, ramMB, activation, []string{
			fmt.Sprintf("image estimate (offload-to-cpu): %d MB of weights staged in RAM, %d MB VRAM for the compute working set, %d MB decode activation", weights/bytesPerMB, sdModuleReserveMB, activation),
		}, nil
	}

	scale := weightTypeScale(req.WeightType)
	conditionerBytes := req.ConditionerBytes
	if !req.ClipVisionCPU {
		conditionerBytes += req.ClipVisionBytes
	}
	modules := []struct {
		name    string
		bytes   int64
		reserve int64
	}{
		{"checkpoint", int64(float64(req.CheckpointBytes) * scale), sdModuleReserveMB},
		{"conditioner", int64(float64(conditionerBytes) * scale), sdModuleReserveMB},
		{"vae", int64(float64(req.VAEBytes) * scale), sdVAEReserveMB},
		{"adapter", int64(float64(req.ControlNetBytes+req.LoRABytes+req.UpscalerBytes+req.IPAdapterBytes) * scale), sdModuleReserveMB},
	}

	required := sdModuleReserveMB
	if req.ResidentLoaded {
		// The fixed reserve belongs to the already-loaded base process. A delta
		// stage must be judged on the modules it adds, not on that reserve again.
		required = 0
	}
	largest := "engine"
	for _, module := range modules {
		if module.bytes <= 0 {
			continue
		}
		need := module.bytes / bytesPerMB
		if !req.ResidentLoaded {
			need += module.reserve
		}
		if need > required {
			required, largest = need, module.name
		}
	}

	// RAM is not the GPU question: it keeps the old components-plus-overhead
	// shape for whichever modules sd-server stages in system memory.
	componentMB := float64(req.ControlNetBytes+req.LoRABytes+req.UpscalerBytes+req.VAEBytes+req.IPAdapterBytes) / float64(bytesPerMB)
	if req.ClipVisionCPU {
		componentMB += float64(req.ClipVisionBytes) / float64(bytesPerMB)
	}
	return required, int64(componentMB + 500), activation, []string{
		fmt.Sprintf("image estimate: %s needs %d MB (weights + %d MB reserve, sd-server backend_fit rule, type=%q), %d MB decode activation", largest, required, sdModuleReserveMB, req.WeightType, activation),
	}, nil
}

// sdActivationMBPerMegapixel is the measured untiled VAE decode coefficient for
// the SDXL-family VAE on this 8 GiB GTX 1080. Verbose sd-server runs measured the
// compute buffer at 480 MB for 256x256, 1920 MB for 512x512 and 4320 MB for
// 768x768: exactly 7325 MB per megapixel. The old 2100 MB/MP estimate (551 MB at
// 512) let the base stage exceed the free-VRAM reserve. This is a planning
// estimate; the engine's own buffer line remains the capacity proof.
const sdActivationMBPerMegapixel = 7325

// sdTiledActivationMB is the measured VAE compute buffer at 512x512 with
// --vae-tiling (3x3 tiles). Tiling cuts the decode working set from 1920 MB to
// 480 MB and is what lets the base stage fit while holding the reserve.
const sdTiledActivationMB = 480

// sdFixedOverheadMB covers the engine's non-module allocations during a stage
// (contexts, temporary buffers and allocator slack). The verbose base run
// measured 7579 MB used minus a ~1080 MB idle baseline = ~6498 MB delta; the
// 4100 MB converted resident set plus the 1920 MB decode leaves 478 MB. Use 512
// as the measured upper bound so a 26 MB planning margin cannot hide a reserve
// breach on a system whose idle baseline shifts.
const sdFixedOverheadMB = 512

// imageDecodeActivationMB estimates the peak decode working set. Howl's batch
// path submits one image at a time (generateSequentialBatch), so batch count
// affects elapsed time, not peak memory; steps affect time, not peak memory.
func imageDecodeActivationMB(width, height, batchCount int) int64 {
	if width <= 0 || height <= 0 {
		return 0
	}
	// Peak is one image even when the user asked for a batch: the client splits
	// it into sequential single-image requests.
	_ = batchCount
	mp := float64(width) * float64(height) / 1_000_000
	return int64(mp*sdActivationMBPerMegapixel + 0.5)
}

// imageControlNetVRAMMB is the ControlNet params set the pinned engine stages
// to CUDA during sampling. The 5,004,167,829-byte OpenPoseXL2.safetensors loaded
// as 1415.82 MiB (ratio 0.297); 0.30 is the measured planning ratio. The load
// itself reports controlnet::load -> CPU, but model_manager then stages all
// 1415.92 MB of compute params to CUDA0 for each sampling phase.
func imageControlNetVRAMMB(req AdmissionRequest) int64 {
	if req.ControlNetBytes <= 0 || req.ControlNetCPU {
		return 0
	}
	return int64(float64(req.ControlNetBytes)*0.30) / bytesPerMB
}

// ImageResidentVRAMMB exposes the converted resident module total to the image
// placement planner. It is the additive set, unlike the largest-module figure
// admission uses for its fit rule.
func ImageResidentVRAMMB(req AdmissionRequest) int64 { return imageResidentVRAMMB(req) }

// ImageFixedOverheadMB is the engine overhead included in every projected peak.
// It is exported so the preflight budget decision uses the same fixed number as
// admission instead of duplicating it.
const ImageFixedOverheadMB = sdFixedOverheadMB

// imageResidentVRAMMB sizes the converted modules that stay on the device while
// an image stage runs. ControlNet weights stay in system RAM and are charged as
// their own module; LoRA weights are applied to the resident model.
func imageResidentVRAMMB(req AdmissionRequest) int64 {
	scale := weightTypeScale(req.WeightType)
	// Sum the modules after conversion, as estimateImage does. Rounding the
	// combined byte count once changes the 8 GiB planning number by a megabyte
	// and makes the recorded projection unnecessarily hard to compare.
	var total int64
	conditioner := req.ConditionerBytes
	if !req.ClipVisionCPU {
		conditioner += req.ClipVisionBytes
	}
	for _, bytes := range []int64{req.CheckpointBytes, conditioner, req.VAEBytes, req.LoRABytes, req.IPAdapterBytes} {
		if bytes > 0 {
			total += int64(float64(bytes)*scale) / bytesPerMB
		}
	}
	return total
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// weightTypeScale is the fraction of stored bytes a module keeps after sd-server
// converts it to --type while loading.
//
// q8_0 is measured: Illustrious-XL-v2.0 is 6624.11 MB of f16 weights and loaded
// as 3996.05 MB (0.603) under --type q8_0. 0.62 keeps a small margin on the safe
// side of that measurement. q5/q4 are the same family's block sizes, scaled down
// from the q8 ratio and rounded up; they are NOT measured on this machine.
func weightTypeScale(t string) float64 {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "q8_0":
		return 0.62
	case "q5_0", "q5_1":
		return 0.40
	case "q4_0", "q4_1":
		return 0.32
	case "f32":
		return 2.0
	default: // f16, bf16, empty and anything sd.cpp treats as no-op
		return 1.0
	}
}

func refused(req AdmissionRequest, est AdmissionEstimate, blocking []ModelState, reason string) AdmissionDecision {
	return AdmissionDecision{Admitted: false, Reason: reason, Estimate: est, Blocking: blocking}
}

func vramReason(req AdmissionRequest, est AdmissionEstimate, blocking []ModelState) string {
	need := est.committedVRAMMB()
	reason := fmt.Sprintf(
		"cannot start %s %q: needs %d MB VRAM; %d MB free on device %d (%s), %d MB reserve (%.0f%%), %d MB available",
		req.Kind, req.Name, need, est.FreeVRAMMB, est.DeviceIndex, displayGPU(est.GPUName), est.ReserveVRAMMB, reservePercent(est), est.AvailableVRAMMB,
	)
	if est.ProjectedPeakVRAMMB > 0 {
		if req.ResidentLoaded {
			reason += fmt.Sprintf("; projected stage peak %d MB = new modules %d MB + activation %d MB + %d MB engine overhead",
				est.ProjectedPeakVRAMMB, est.RequiredVRAMMB, est.ActivationVRAMMB, sdFixedOverheadMB)
		} else {
			reason += fmt.Sprintf("; projected peak %d MB = resident %d MB + activation %d MB + %d MB engine overhead",
				est.ProjectedPeakVRAMMB, est.ResidentVRAMMB, est.ActivationVRAMMB, sdFixedOverheadMB)
		}
	}
	if est.CommitmentFloorMB > 0 {
		reason += fmt.Sprintf("; a measured %d MB peak for this model and workload shape raised the commitment above the %d MB projection",
			est.CommitmentFloorMB, est.ProjectedPeakVRAMMB)
	}
	if est.SoftReserveAllowanceMB > 0 && need > est.AvailableVRAMMB {
		reason += fmt.Sprintf("; the %d MB soft reserve allowance is not enough (it would leave less than %d MB free)",
			est.SoftReserveAllowanceMB, hardReserveFloorMB)
	}
	return reason + blockingSuffix(blocking)
}

func ramReason(req AdmissionRequest, est AdmissionEstimate, blocking []ModelState) string {
	reason := fmt.Sprintf(
		"cannot start %s %q: needs %d MB RAM; %d MB free, %d MB reserve (%.0f%%), %d MB available",
		req.Kind, req.Name, est.RequiredRAMMB, est.FreeRAMMB, est.ReserveRAMMB, reservePercent(est), est.AvailableRAMMB,
	)
	return reason + blockingSuffix(blocking)
}

func reservePercent(est AdmissionEstimate) float64 {
	if est.FreeVRAMMB <= 0 && est.FreeRAMMB <= 0 {
		return 0
	}
	base := est.FreeVRAMMB
	reserve := est.ReserveVRAMMB
	if base <= 0 {
		base = est.FreeRAMMB
		reserve = est.ReserveRAMMB
	}
	if base <= 0 {
		return 0
	}
	return float64(reserve) / float64(base) * 100
}

func displayGPU(name string) string {
	if strings.TrimSpace(name) == "" {
		return "detected GPU"
	}
	return name
}

func blockingSuffix(blocking []ModelState) string {
	if len(blocking) == 0 {
		return ""
	}
	first := blocking[0]
	switch {
	case first.GPUMemoryMB > 0 && first.SystemRAMMB > 0:
		return fmt.Sprintf("; blocking model: %s %q holds %d MB VRAM and %d MB RAM", first.Type, first.Name, first.GPUMemoryMB, first.SystemRAMMB)
	case first.GPUMemoryMB > 0:
		return fmt.Sprintf("; blocking model: %s %q holds %d MB VRAM", first.Type, first.Name, first.GPUMemoryMB)
	default:
		return fmt.Sprintf("; blocking model: %s %q holds %d MB RAM", first.Type, first.Name, first.SystemRAMMB)
	}
}
