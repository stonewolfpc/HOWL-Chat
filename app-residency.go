package main

// app-residency.go — the measurement record behind the resident ledger.
//
// Each load records three numbers (planned, engine-reported, device delta) and
// the ledger value derived from them. The record is for evidence and for the
// next decision's blocker list; it never feeds back into admission's estimate,
// so an observation cannot quietly become a planning assumption.

import (
	"fmt"
	"sync"

	"howl-chat/internal/image/engine/sdcpp"
	"howl-chat/internal/resource"
)

type residencyRecorder struct {
	mu       sync.Mutex
	last     map[resource.ResourceType]resource.ResidencyObservation
	baseline map[resource.ResourceType]residencyBaseline
}

// residencyBaseline is what a workload's admission committed and what the device
// already held before its process started. Image residency can only be measured
// after the first generation (sd-server loads weights lazily and reports its
// placement then), so the baseline has to survive until that moment.
type residencyBaseline struct {
	PlannedMB    int64
	UsedBeforeMB int64
}

func (a *App) noteResidencyBaseline(kind resource.ResourceType, plannedMB, usedBeforeMB int64) {
	a.residencyOnce.Do(func() {
		if a.residency == nil {
			a.residency = &residencyRecorder{last: make(map[resource.ResourceType]resource.ResidencyObservation), baseline: make(map[resource.ResourceType]residencyBaseline)}
		}
	})
	recorder := a.residency
	recorder.mu.Lock()
	defer recorder.mu.Unlock()
	if recorder.baseline == nil {
		recorder.baseline = make(map[resource.ResourceType]residencyBaseline)
	}
	recorder.baseline[kind] = residencyBaseline{PlannedMB: plannedMB, UsedBeforeMB: usedBeforeMB}
}

func (a *App) residencyBaselineFor(kind resource.ResourceType) (residencyBaseline, bool) {
	if a.residency == nil {
		return residencyBaseline{}, false
	}
	a.residency.mu.Lock()
	defer a.residency.mu.Unlock()
	baseline, ok := a.residency.baseline[kind]
	return baseline, ok
}

// reconcileRegisteredResidency rewrites the resident ledger entry from a
// measurement. The registry stays the single ledger owner; this is its one
// measured-update path, and it never invents a number when nothing was observed.
func (a *App) reconcileRegisteredResidency(kind resource.ResourceType, name string, vramMB int64) {
	if vramMB <= 0 {
		return
	}
	for _, model := range resource.GetManager().ModelsByType(kind) {
		if model.Name != name {
			continue
		}
		resource.GetManager().RegisterModel(kind, model.Name, model.Path, vramMB, model.SystemRAMMB, model.Priority, vramMB > 0)
	}
}

func (a *App) recordResidency(observation resource.ResidencyObservation) {
	agreement := "within plan"
	if !observation.AgreesWithPlan() {
		agreement = "above plan"
	}
	fmt.Printf("[Residency] kind=%s model=%q planned=%dMB engine=%dMB device_delta=%dMB ledger=%dMB %s\n",
		observation.Kind, observation.Model, observation.PlannedVRAMMB,
		observation.EngineReportedVRAMMB, observation.DeviceDeltaVRAMMB,
		observation.LedgerVRAMMB(), agreement)

	a.residencyOnce.Do(func() {
		if a.residency == nil {
			a.residency = &residencyRecorder{last: make(map[resource.ResourceType]resource.ResidencyObservation)}
		}
	})
	recorder := a.residency
	recorder.mu.Lock()
	recorder.last[observation.Kind] = observation
	recorder.mu.Unlock()
}

// recordImageResidencyFromGeneration is the image observation point. sd-server
// loads weights lazily, so a load-time record measured nothing; after the first
// generation both the engine's placement and the device's usage are real.
func (a *App) recordImageResidencyFromGeneration(engine *sdcpp.Engine, model string) {
	if engine == nil || model == "" {
		return
	}
	baseline, ok := a.residencyBaselineFor(resource.ResourceImage)
	if !ok {
		return
	}
	// The pinned sd-server (3b12781) reports no per-component placement: its
	// only placement statement is unknown. So the engine's VRAM claim is zero
	// here by fact, not by omission, and the measured device change is the
	// evidence. When a build does report a total, it is used as the other bound.
	var engineReportedMB int64
	if placement := engine.GetPlacement(); placement.Known {
		engineReportedMB = int64(placement.VRAMMB)
	}
	snap, err := resource.GetManager().RefreshSnapshot()
	if err != nil {
		return
	}
	const mb = int64(1024 * 1024)
	afterMB := (snap.TotalVRAMBytes - snap.FreeVRAMBytes) / mb
	observation := resource.NewResidencyObservation(
		resource.ResourceImage, model,
		baseline.PlannedMB,
		engineReportedMB,
		afterMB-baseline.UsedBeforeMB,
	)
	if observation.ObservedVRAMMB() == 0 {
		return
	}
	a.recordResidency(observation)
	a.reconcileRegisteredResidency(resource.ResourceImage, model, observation.LedgerVRAMMB())
}

// GetResidencyReport returns the last observation per workload, for the same
// kind of read-back the placement report offers. Empty before any load.
func (a *App) GetResidencyReport() map[string]interface{} {
	out := map[string]interface{}{"observations": map[string]interface{}{}}
	if a.residency == nil {
		return out
	}
	a.residency.mu.Lock()
	defer a.residency.mu.Unlock()
	observations := map[string]interface{}{}
	for kind, observation := range a.residency.last {
		observations[string(kind)] = map[string]interface{}{
			"model":                  observation.Model,
			"planned_vram_mb":        observation.PlannedVRAMMB,
			"engine_reported_vram_mb": observation.EngineReportedVRAMMB,
			"device_delta_vram_mb":   observation.DeviceDeltaVRAMMB,
			"ledger_vram_mb":         observation.LedgerVRAMMB(),
			"agrees_with_plan":       observation.AgreesWithPlan(),
			"at":                     observation.At.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	out["observations"] = observations
	return out
}
