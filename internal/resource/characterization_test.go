package resource

import "testing"

// TestAdmissionReleaseIsReservationOnly characterizes the current meaning of
// AdmissionLease.Release: it removes an in-flight capacity reservation, but it
// does not unregister a resident model. This test is expected to pass today.
func TestAdmissionReleaseIsReservationOnly(t *testing.T) {
	loaded := []ModelState{{
		Type: ResourceImage, Name: "resident-image", Loaded: true, GPUMemoryMB: 3000,
	}}
	authority := NewAuthority(fixedHardware(testSnapshot(8000, 8192)), fixedLedger(loaded...))

	decision, lease := authority.Admit(AdmissionRequest{
		Kind: ResourceTTS, Name: "resident-lease", AudioVRAMMB: 100, RequireGPU: true,
	})
	if !decision.Admitted {
		t.Fatalf("admission refused a request that fits: %s", decision.Reason)
	}
	if lease == nil {
		t.Fatal("admitted request returned no lease")
	}
	if got := len(authority.inFlight); got != 1 {
		t.Fatalf("in-flight reservations = %d, want 1 before release", got)
	}

	lease.Release()

	if got := len(authority.inFlight); got != 0 {
		t.Fatalf("in-flight reservations = %d, want 0 after release", got)
	}
	if got := len(loaded); got != 1 {
		t.Fatalf("release removed the resident ledger entry: got %d models", got)
	}
}

// TestPendingReservationBlocksSecondAdmission characterizes the current
// in-flight ledger: a second request is judged against physical free memory
// minus outstanding reservations, so it cannot spend capacity twice while the
// first lease is open. This test is expected to pass today.
func TestPendingReservationBlocksSecondAdmission(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(10000, 24000)), fixedLedger())

	first, firstLease := authority.Admit(AdmissionRequest{
		Kind: ResourceTTS, Name: "first", AudioVRAMMB: 6000,
	})
	if !first.Admitted {
		t.Fatalf("first request should fit: %s", first.Reason)
	}

	second, secondLease := authority.Admit(AdmissionRequest{
		Kind: ResourceTTS, Name: "second", AudioVRAMMB: 5000,
	})
	if second.Admitted {
		secondLease.Release()
		firstLease.Release()
		t.Fatal("second request spent capacity already reserved by the first")
	}
	if secondLease != nil {
		firstLease.Release()
		t.Fatal("a refused request returned a lease")
	}

	firstLease.Release()

	third, thirdLease := authority.Admit(AdmissionRequest{
		Kind: ResourceTTS, Name: "second-after-release", AudioVRAMMB: 5000,
	})
	if !third.Admitted {
		t.Fatalf("request should fit after the first reservation is released: %s", third.Reason)
	}
	thirdLease.Release()
}

// TestResidentStageIsNotChargedResidentWeightsAgain characterizes the current
// stage-delta path: physical free memory already excludes the resident image
// model, so a ResidentLoaded stage is judged on its incremental peak rather
// than on the resident weights again. This test is expected to pass today.
func TestResidentStageIsNotChargedResidentWeightsAgain(t *testing.T) {
	loaded := []ModelState{{
		Type: ResourceImage, Name: "resident-image", Loaded: true, GPUMemoryMB: 6000,
	}}
	authority := NewAuthority(fixedHardware(testSnapshot(2000, 8192)), fixedLedger(loaded...))

	decision, lease := authority.Admit(AdmissionRequest{
		Kind:           ResourceImage,
		Name:           "resident-upscale",
		Mode:           "dgl_upscale",
		Width:          512,
		Height:         512,
		ResidentLoaded: true,
		VAETiling:      true,
		RequireGPU:     true,
	})
	if !decision.Admitted {
		t.Fatalf("resident stage should fit on its incremental peak: %s", decision.Reason)
	}
	defer lease.Release()

	if decision.Estimate.RequiredVRAMMB != 0 {
		t.Fatalf("resident stage required = %d MB, want 0 resident weights", decision.Estimate.RequiredVRAMMB)
	}
	if decision.Estimate.ProjectedPeakVRAMMB != 512 {
		t.Fatalf("resident stage projected peak = %d MB, want 512 MB engine overhead", decision.Estimate.ProjectedPeakVRAMMB)
	}
}
