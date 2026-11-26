package resource

// residency.go — what a loaded model actually costs.
//
// Admission estimates a footprint; the engine reports one; the device shows one.
// They disagree, and the disagreements are systematic:
//
//   - a chat estimate charges the whole model set plus KV, while the engine's
//     own total omits the CUDA context and allocator slack;
//   - an image ledger entry registered RequiredVRAMMB while the admitted
//     commitment was the projected peak (resident + activation + overhead).
//
// This type keeps the three numbers apart. It exists to be measured, not to
// change a decision: the ledger takes the conservative observation (the larger
// of what the engine reports and what the device actually gained), and the
// planned value stays a planning bound. No estimate is rewritten here.

import "time"

// ResidencyObservation is one load's footprint, as three separate claims.
type ResidencyObservation struct {
	Kind   ResourceType `json:"kind"`
	Model  string       `json:"model"`
	Source string       `json:"source,omitempty"`
	At     time.Time    `json:"at"`

	// PlannedVRAMMB is what admission charged for this workload.
	PlannedVRAMMB int64 `json:"planned_vram_mb"`
	// EngineReportedVRAMMB is the engine's own resident total (weights + KV +
	// compute). Zero means the engine reported nothing usable.
	EngineReportedVRAMMB int64 `json:"engine_reported_vram_mb"`
	// DeviceDeltaVRAMMB is the live device-used change across the load. It can
	// include other processes, so it is an upper bound on this model's share,
	// never a lower one.
	DeviceDeltaVRAMMB int64 `json:"device_delta_vram_mb"`
}

// NewResidencyObservation builds a record. Negative device deltas are recorded
// as zero: a load that appeared to free memory says nothing about residency.
func NewResidencyObservation(kind ResourceType, model string, planned, engineReported, deviceDelta int64) ResidencyObservation {
	if deviceDelta < 0 {
		deviceDelta = 0
	}
	if engineReported < 0 {
		engineReported = 0
	}
	return ResidencyObservation{
		Kind:                 kind,
		Model:                model,
		At:                   time.Now(),
		PlannedVRAMMB:        planned,
		EngineReportedVRAMMB: engineReported,
		DeviceDeltaVRAMMB:    deviceDelta,
	}
}

// ObservedVRAMMB is the conservative resident figure: the larger of the engine's
// claim and the device's own change. Zero means nothing was observed.
func (o ResidencyObservation) ObservedVRAMMB() int64 {
	if o.EngineReportedVRAMMB > o.DeviceDeltaVRAMMB {
		return o.EngineReportedVRAMMB
	}
	return o.DeviceDeltaVRAMMB
}

// LedgerVRAMMB is what the resident ledger should hold for this load. It prefers
// the observation and falls back to the plan only when nothing was measured, so
// an unmeasured load keeps the old, conservative entry instead of claiming zero.
func (o ResidencyObservation) LedgerVRAMMB() int64 {
	if observed := o.ObservedVRAMMB(); observed > 0 {
		return observed
	}
	return o.PlannedVRAMMB
}

// AgreesWithPlan reports whether the conservative observation stayed inside the
// plan. A false value is not automatically a defect — it is the signal that an
// estimate needs revisiting with this model's measurements.
func (o ResidencyObservation) AgreesWithPlan() bool {
	observed := o.ObservedVRAMMB()
	return observed == 0 || o.PlannedVRAMMB <= 0 || observed <= o.PlannedVRAMMB
}
