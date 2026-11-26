package resource

import (
	"strings"
	"testing"
)

// An observed peak is stronger evidence than a metadata projection, so it may
// raise the commitment. These tests pin the direction and the limits of that
// precedence: only upward, only for a comparable workload shape, and never by
// pretending the observation is a resident-only figure.

func imageRequestWithFloor(floor int64) AdmissionRequest {
	const gib = int64(1024 * 1024 * 1024)
	return AdmissionRequest{
		Kind: ResourceImage, Name: "illustrious",
		CheckpointBytes: 4 * gib, Width: 512, Height: 512,
		BatchCount: 1, Mode: "txt2img", RequireGPU: true,
		MinimumCommitmentVRAMMB: floor,
	}
}

func TestObservedFloorRaisesTheCommitment(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(8000, 8192)), fixedLedger())

	decision, lease := authority.Admit(imageRequestWithFloor(7000))
	if !decision.Admitted {
		t.Fatalf("floor case was refused: %s", decision.Reason)
	}
	defer lease.Release()

	if got := decision.Estimate.ProjectedPeakVRAMMB; got != 6528 {
		t.Fatalf("the projection itself must stay an estimate, got %d MB", got)
	}
	if got, want := decision.Estimate.CommitmentFloorMB, int64(7000); got != want {
		t.Fatalf("commitment floor = %d MB, want %d MB", got, want)
	}
	if got, want := decision.Estimate.CommittedVRAMMB(), int64(7000); got != want {
		t.Fatalf("commitment = %d MB, want the measured floor %d MB", got, want)
	}
	if got, want := authority.inFlight[0].vram, int64(7000); got != want {
		t.Fatalf("in-flight reservation = %d MB, want the measured floor %d MB", got, want)
	}
	if !strings.Contains(decision.Reason, "7000 MB") || !strings.Contains(decision.Reason, "measured") {
		t.Fatalf("grant must name the measured floor: %s", decision.Reason)
	}
	if got := decision.Estimate.RequiredVRAMMB; got == 7000 {
		t.Fatal("the observation must not replace the resident requirement")
	}
}

func TestObservedFloorNeverLowersThePlan(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(8000, 8192)), fixedLedger())

	decision, lease := authority.Admit(imageRequestWithFloor(4000))
	if !decision.Admitted {
		t.Fatalf("case was refused: %s", decision.Reason)
	}
	defer lease.Release()

	if got, want := decision.Estimate.CommittedVRAMMB(), int64(6528); got != want {
		t.Fatalf("commitment = %d MB, want the projection %d MB", got, want)
	}
	if decision.Estimate.CommitmentFloorMB != 0 {
		t.Fatalf("a smaller observation must not be recorded as a binding floor: %d MB", decision.Estimate.CommitmentFloorMB)
	}
}

// The floor is only useful if it can refuse. Here the projection fits with
// reserve left over, and the measured peak does not.
func TestObservedFloorCanRefuseWhatTheProjectionAdmits(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(7500, 8192)), fixedLedger())

	projectionOnly, projectionLease := authority.Admit(imageRequestWithFloor(0))
	if !projectionOnly.Admitted {
		t.Fatalf("test setup: the projection should fit on its own: %s", projectionOnly.Reason)
	}
	projectionLease.Release()

	measured, _ := authority.Admit(imageRequestWithFloor(7300))
	if measured.Admitted {
		t.Fatalf("a measured peak above the budget was admitted: %s", measured.Reason)
	}
	if !strings.Contains(measured.Reason, "measured 7300 MB peak") {
		t.Fatalf("refusal must name the measurement that changed the verdict: %s", measured.Reason)
	}
}

func TestObservationShapeDecidesWhetherItIsReused(t *testing.T) {
	observation := NewResidencyObservation(ResourceImage, "illustrious", 5093, 0, 5679).
		WithGeneration(512, 512, 1, "q8_0")

	if observation.Coverage != CoveragePostGenerationDelta {
		t.Fatalf("coverage = %q, want %q", observation.Coverage, CoveragePostGenerationDelta)
	}
	cases := []struct {
		name                 string
		width, height, batch int
		weight               string
		want                 bool
	}{
		{"same shape and precision", 512, 512, 1, "q8_0", true},
		{"smaller shape", 256, 256, 1, "q8_0", true},
		{"larger width", 768, 512, 1, "q8_0", false},
		{"larger batch", 512, 512, 2, "q8_0", false},
		{"different precision", 512, 512, 1, "f16", false},
		{"unknown precision", 512, 512, 1, "", false},
	}
	for _, tt := range cases {
		if got := observation.AppliesToGeneration(tt.width, tt.height, tt.batch, tt.weight); got != tt.want {
			t.Fatalf("%s: applies = %t, want %t", tt.name, got, tt.want)
		}
	}
}

func TestObservationWithoutShapeIsNotReusable(t *testing.T) {
	observation := NewResidencyObservation(ResourceImage, "illustrious", 5093, 0, 5679)
	if observation.AppliesToGeneration(512, 512, 1, "q8_0") {
		t.Fatal("an observation with no measured shape must not bound a request")
	}
}
