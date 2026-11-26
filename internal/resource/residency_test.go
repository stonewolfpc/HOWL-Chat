package resource

import "testing"

func TestLedgerPrefersTheLargerObservation(t *testing.T) {
	engine := NewResidencyObservation(ResourceLLM, "m", 6057, 4888, 0)
	if got := engine.LedgerVRAMMB(); got != 4888 {
		t.Fatalf("engine-only observation: ledger=%d, want 4888", got)
	}

	device := NewResidencyObservation(ResourceLLM, "m", 6057, 4888, 5600)
	if got := device.LedgerVRAMMB(); got != 5600 {
		t.Fatalf("device larger than engine: ledger=%d, want 5600", got)
	}

	engineWins := NewResidencyObservation(ResourceImage, "m", 2048, 6100, 5400)
	if got := engineWins.LedgerVRAMMB(); got != 6100 {
		t.Fatalf("engine larger than device: ledger=%d, want 6100", got)
	}
}

func TestUnmeasuredLoadKeepsThePlan(t *testing.T) {
	unmeasured := NewResidencyObservation(ResourceLLM, "m", 6057, 0, 0)
	if got := unmeasured.LedgerVRAMMB(); got != 6057 {
		t.Fatalf("unmeasured observation must keep the plan: ledger=%d, want 6057", got)
	}
	if !unmeasured.AgreesWithPlan() {
		t.Fatal("an unmeasured load cannot disagree with its plan")
	}
}

func TestNegativeMeasurementsAreDiscarded(t *testing.T) {
	shrunk := NewResidencyObservation(ResourceImage, "m", 3000, -50, -400)
	if shrunk.EngineReportedVRAMMB != 0 || shrunk.DeviceDeltaVRAMMB != 0 {
		t.Fatalf("negative measurements must not be recorded: %+v", shrunk)
	}
	if got := shrunk.LedgerVRAMMB(); got != 3000 {
		t.Fatalf("with nothing measured the plan stands: ledger=%d, want 3000", got)
	}
}

func TestObservationAbovePlanIsReported(t *testing.T) {
	over := NewResidencyObservation(ResourceImage, "hidream", 2048, 11249, 11000)
	if over.AgreesWithPlan() {
		t.Fatal("a footprint above the plan must not report agreement")
	}
	if got := over.LedgerVRAMMB(); got != 11249 {
		t.Fatalf("ledger=%d, want the larger engine figure 11249", got)
	}
}
