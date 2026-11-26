// These tests graduated from the knownfailures characterization gate on
// 2026-11-25. Their original assertions are unchanged and now run in the
// default suite. D-1 and D-9 close only because these assertions pass here.
package resource

import (
	"strings"
	"testing"
)

// TestGrantReservesProjectedPeak pins D-1 after Gate 1: admission judges the
// projected peak, and grant must record that same commitment. A resident DGL
// stage has a zero resident requirement and a non-zero projected peak; before
// Gate 1 the pending ledger reserved less than the authority approved.
func TestGrantReservesProjectedPeak(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(8000, 8192)), fixedLedger())

	decision, lease := authority.Admit(AdmissionRequest{
		Kind:           ResourceImage,
		Name:           "resident-stage",
		Mode:           "dgl_upscale",
		ResidentLoaded: true,
		RequireGPU:     true,
	})
	if !decision.Admitted {
		t.Fatalf("test setup: resident stage was refused: %s", decision.Reason)
	}
	defer lease.Release()

	if got := len(authority.inFlight); got != 1 {
		t.Fatalf("in-flight reservations = %d, want 1", got)
	}
	if got, want := authority.inFlight[0].vram, decision.Estimate.ProjectedPeakVRAMMB; got != want {
		t.Fatalf("pending reservation = %d MB, want the projected peak %d MB", got, want)
	}
}

// TestRAMOnlyResidentBlocksRAMAdmission pins D-9 after Gate 1: blockingModels
// used to ignore a resident with GPUMemoryMB <= 0. A RAM-only resident can
// still be the reason a RAM request is refused, so it must be named.
func TestRAMOnlyResidentBlocksRAMAdmission(t *testing.T) {
	loaded := []ModelState{{
		Type: ResourceLLM, Name: "ram-only-chat", Loaded: true, GPUMemoryMB: 0, SystemRAMMB: 4096,
	}}
	authority := NewAuthority(fixedHardware(testSnapshot(8000, 8192)), fixedLedger(loaded...))

	decision, _ := authority.Admit(AdmissionRequest{
		Kind: ResourceASR, Name: "large-asr", AudioRAMMB: 20000, RequireGPU: true,
	})
	if decision.Admitted {
		t.Fatal("test setup: RAM request should have been refused")
	}
	if len(decision.Blocking) != 1 || decision.Blocking[0].Name != "ram-only-chat" {
		t.Fatalf("RAM-only resident was not reported as the blocker: %+v", decision.Blocking)
	}
}

// TestCommitmentUsesComparableProjectedPeakAcrossPaths pins the D-1 accounting
// contract for all three admission shapes. A full image load commits its total
// projected peak, a resident stage commits its incremental projected peak, and
// a replacement commits the new model's peak. None of them commits
// Required + Projected, and the replacement credit is applied to effective free
// rather than subtracted from the commitment.
func TestCommitmentUsesComparableProjectedPeakAcrossPaths(t *testing.T) {
	const gib = int64(1024 * 1024 * 1024)
	old := ModelState{
		Type: ResourceImage, Name: "old-image", Loaded: true, GPUMemoryMB: 3000,
	}

	tests := []struct {
		name     string
		req      AdmissionRequest
		loaded   []ModelState
		freeMB   int64
		wantPeak int64
	}{
		{
			name: "full-load",
			req: AdmissionRequest{
				Kind: ResourceImage, Name: "full-image",
				CheckpointBytes: 4 * gib, Width: 512, Height: 512,
				BatchCount: 1, Mode: "txt2img", RequireGPU: true,
			},
			freeMB:   8000,
			wantPeak: 6528,
		},
		{
			name: "resident-stage",
			req: AdmissionRequest{
				Kind: ResourceImage, Name: "resident-stage",
				Mode: "dgl_upscale", Width: 512, Height: 512,
				ResidentLoaded: true, RequireGPU: true,
			},
			freeMB:   8000,
			wantPeak: 512,
		},
		{
			name: "resident-stage-with-adapter",
			req: AdmissionRequest{
				Kind: ResourceImage, Name: "resident-stage-with-adapter",
				Mode: "dgl_upscale", Width: 512, Height: 512,
				ResidentLoaded: true, UpscalerBytes: gib, RequireGPU: true,
			},
			freeMB:   8000,
			wantPeak: 1536,
		},
		{
			name: "replacement",
			req: AdmissionRequest{
				Kind: ResourceImage, Name: "new-image",
				CheckpointBytes: 4 * gib, Width: 512, Height: 512,
				BatchCount: 1, Mode: "txt2img", RequireGPU: true,
				ReplacingKind: ResourceImage, ReplacingName: "old-image",
			},
			loaded:   []ModelState{old},
			freeMB:   5000,
			wantPeak: 6528,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authority := NewAuthority(
				fixedHardware(testSnapshot(tt.freeMB, 8192)),
				fixedLedger(tt.loaded...),
			)
			decision, lease := authority.Admit(tt.req)
			if !decision.Admitted {
				t.Fatalf("admission refused the comparable-peak case: %s", decision.Reason)
			}
			defer lease.Release()

			if got := decision.Estimate.ProjectedPeakVRAMMB; got != tt.wantPeak {
				t.Fatalf("projected peak = %d MB, want %d MB", got, tt.wantPeak)
			}
			if len(authority.inFlight) != 1 {
				t.Fatalf("in-flight commitments = %d, want 1", len(authority.inFlight))
			}
			got := authority.inFlight[0].vram
			if got != tt.wantPeak {
				t.Fatalf("committed VRAM = %d MB, want the projected peak %d MB", got, tt.wantPeak)
			}
			if decision.Estimate.RequiredVRAMMB > 0 && decision.Estimate.ProjectedPeakVRAMMB > 0 {
				doubleCounted := decision.Estimate.RequiredVRAMMB + decision.Estimate.ProjectedPeakVRAMMB
				if got == doubleCounted {
					t.Fatalf("commitment double-counted required + projected: %d MB", got)
				}
			}
		})
	}
}

// TestRAMOnlyResidentReasonNamesRAM makes the blocker useful to the user, not
// only present in the structured decision. The RAM-only model must be named
// with its RAM footprint, not reported as a zero-MB GPU resident.
func TestRAMOnlyResidentReasonNamesRAM(t *testing.T) {
	loaded := []ModelState{{
		Type: ResourceLLM, Name: "ram-only-chat", Loaded: true, GPUMemoryMB: 0, SystemRAMMB: 4096,
	}}
	authority := NewAuthority(fixedHardware(testSnapshot(8000, 8192)), fixedLedger(loaded...))

	decision, _ := authority.Admit(AdmissionRequest{
		Kind: ResourceASR, Name: "large-asr", AudioRAMMB: 20000, RequireGPU: true,
	})
	if decision.Admitted {
		t.Fatal("test setup: RAM request should have been refused")
	}
	if !strings.Contains(decision.Reason, "ram-only-chat") {
		t.Fatalf("refusal must name the RAM-only resident: %s", decision.Reason)
	}
	if !strings.Contains(decision.Reason, "4096 MB RAM") {
		t.Fatalf("refusal must name the RAM footprint, not a zero-MB GPU footprint: %s", decision.Reason)
	}
}

// TestRAMOnlyAndGPUResidentsBothAppear checks that adding the RAM-only blocker
// did not replace the GPU blocker. The decision must describe both residents
// when both consume capacity relevant to the request.
func TestRAMOnlyAndGPUResidentsBothAppear(t *testing.T) {
	loaded := []ModelState{
		{Type: ResourceImage, Name: "gpu-image", Loaded: true, GPUMemoryMB: 3000},
		{Type: ResourceLLM, Name: "ram-only-chat", Loaded: true, GPUMemoryMB: 0, SystemRAMMB: 4096},
	}
	authority := NewAuthority(fixedHardware(testSnapshot(8000, 8192)), fixedLedger(loaded...))

	decision, _ := authority.Admit(AdmissionRequest{
		Kind: ResourceASR, Name: "large-asr", AudioRAMMB: 20000, RequireGPU: true,
	})
	if decision.Admitted {
		t.Fatal("test setup: RAM request should have been refused")
	}
	seen := map[string]bool{}
	for _, model := range decision.Blocking {
		seen[model.Name] = true
	}
	if !seen["gpu-image"] || !seen["ram-only-chat"] {
		t.Fatalf("blocking must include both residents: %+v", decision.Blocking)
	}
}

// TestResidentStageCommitmentBlocksSecondAdmission proves that the D-1 fix is
// not only descriptive: the resident stage's 512 MB commitment must reduce the
// capacity a second admission can spend. Before the fix the pending value was
// zero and the second 7500 MB request fit through the soft-reserve allowance.
func TestResidentStageCommitmentBlocksSecondAdmission(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(8000, 8192)), fixedLedger())

	stage, stageLease := authority.Admit(AdmissionRequest{
		Kind: ResourceImage, Name: "resident-stage",
		Mode: "dgl_upscale", Width: 512, Height: 512,
		ResidentLoaded: true, RequireGPU: true,
	})
	if !stage.Admitted {
		t.Fatalf("resident stage should fit: %s", stage.Reason)
	}
	defer stageLease.Release()

	second, _ := authority.Admit(AdmissionRequest{
		Kind: ResourceTTS, Name: "second", AudioVRAMMB: 7500,
	})
	if second.Admitted {
		t.Fatal("second request spent the resident stage's committed peak")
	}

	stageLease.Release()
	third, thirdLease := authority.Admit(AdmissionRequest{
		Kind: ResourceTTS, Name: "after-release", AudioVRAMMB: 7500,
	})
	if !third.Admitted {
		t.Fatalf("request should fit after the stage commitment is released: %s", third.Reason)
	}
	thirdLease.Release()
}
