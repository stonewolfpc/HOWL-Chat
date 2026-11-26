package resource

import (
	"strings"
	"testing"
)

func TestImageActivationEstimateIsWorkloadAware(t *testing.T) {
	// Measured untiled VAE compute buffers on the GTX 1080: 480 MB at 256x256,
	// 1920 MB at 512x512 and 4320 MB at 768x768 (7325 MB per megapixel).
	for _, tc := range []struct {
		size int
		want int64
	}{
		{256, 480}, {512, 1920}, {768, 4320}, {1024, 7681},
	} {
		if got := imageDecodeActivationMB(tc.size, tc.size, 1); got != tc.want {
			t.Fatalf("%dx%d activation = %d MB, want %d MB", tc.size, tc.size, got, tc.want)
		}
	}

	authority := NewAuthority(fixedHardware(testSnapshot(8000, 16384)), fixedLedger())
	base := AdmissionRequest{
		Kind: ResourceImage, Name: "sdxl", Width: 1024, Height: 1024, Steps: 20,
		BatchCount: 1, Family: "sdxl", Mode: "txt2img",
		CheckpointBytes: 1000 * bytesPerMB, RequireGPU: true,
	}
	decision, lease := authority.Admit(base)
	if !decision.Admitted {
		t.Fatalf("small checkpoint should fit: %s", decision.Reason)
	}
	lease.Release()
	if !decision.Estimate.TilingPlanned {
		t.Fatalf("1024x1024 decode must plan tiling when the untiled peak exceeds the budget: %+v", decision.Estimate)
	}
	if decision.Estimate.ActivationVRAMMB != sdTiledActivationMB {
		t.Fatalf("tiled 1024x1024 activation = %d MB, want %d MB", decision.Estimate.ActivationVRAMMB, sdTiledActivationMB)
	}

	small := base
	small.Width, small.Height = 512, 512
	decision, lease = authority.Admit(small)
	if !decision.Admitted {
		t.Fatalf("512x512 should fit: %s", decision.Reason)
	}
	lease.Release()
	if decision.Estimate.TilingPlanned {
		t.Fatal("512x512 should fit untiled at this budget")
	}
	if got := decision.Estimate.ActivationVRAMMB; got != 1920 {
		t.Fatalf("512x512 activation = %d MB, want the measured 1920 MB", got)
	}

	// Howl submits batch items sequentially, so a batch must not multiply peak.
	batched := small
	batched.BatchCount = 4
	decision, lease = authority.Admit(batched)
	if !decision.Admitted {
		t.Fatalf("batch should fit: %s", decision.Reason)
	}
	lease.Release()
	if decision.Estimate.ActivationVRAMMB != 1920 {
		t.Fatalf("batch peak = %d MB, want the single-image peak 1920 MB", decision.Estimate.ActivationVRAMMB)
	}
}

func TestImageAdmissionProjectsMeasuredBasePeakAndPlansTiling(t *testing.T) {
	// Measured Illustrious-XL v2.0 q8_0 base stage on the 8 GiB card:
	// resident weights 4100 MB, untiled VAE decode 1920 MB, 512 MB engine
	// overhead = 6532 MB projected; the run measured 7579 MB used with a
	// ~1080 MB idle baseline (~6498 MB delta). With the 6181 MB post-reserve
	// budget (819 MB device floor on the 8 GiB card) the untiled projection
	// crosses the reserve, so admission must plan tiling.
	authority := NewAuthority(fixedHardware(testSnapshot(7000, 8192)), fixedLedger())
	req := AdmissionRequest{
		Kind: ResourceImage, Name: "illustrious", Width: 512, Height: 512,
		Steps: 4, BatchCount: 1, Family: "sdxl", Mode: "txt2img",
		WeightType:       "q8_0",
		CheckpointBytes:  4897 * bytesPerMB,
		ConditionerBytes: 1559 * bytesPerMB,
		VAEBytes:         159 * bytesPerMB,
		RequireGPU:       true,
	}
	decision, lease := authority.Admit(req)
	if !decision.Admitted {
		t.Fatalf("tiled base stage should fit: %s", decision.Reason)
	}
	defer lease.Release()
	if !decision.Estimate.TilingPlanned {
		t.Fatal("untiled projected peak exceeds available; admission must plan tiling")
	}
	if decision.Estimate.ResidentVRAMMB != 4100 {
		t.Fatalf("resident = %d MB, want 4100 MB", decision.Estimate.ResidentVRAMMB)
	}
	if decision.Estimate.ProjectedPeakVRAMMB != 5092 {
		t.Fatalf("tiled projected peak = %d MB, want 5092 MB", decision.Estimate.ProjectedPeakVRAMMB)
	}
	if decision.Estimate.AvailableVRAMMB != 6181 || decision.Estimate.ReserveVRAMMB != 819 {
		t.Fatalf("reserve/available = %d/%d, want 819/6181", decision.Estimate.ReserveVRAMMB, decision.Estimate.AvailableVRAMMB)
	}
	if !strings.Contains(decision.Reason, "projected peak 5092 MB") || !strings.Contains(decision.Reason, "VAE tiling planned") {
		t.Fatalf("reason must report the tiled projection and tiling decision: %s", decision.Reason)
	}
}

func TestImageAdmissionPlansControlNetCPUWhenStagingWouldCrossReserve(t *testing.T) {
	// The pinned engine loads the ControlNet params on CPU but stages the full
	// 1415.82 MiB set to CUDA for sampling. On this 8 GiB card, adding that to the
	// tiled base peak crosses the 6300 MB post-reserve budget, so admission must
	// keep ControlNet on CPU while leaving the base tiling decision intact.
	authority := NewAuthority(fixedHardware(testSnapshot(7000, 8192)), fixedLedger())
	req := AdmissionRequest{
		Kind: ResourceImage, Name: "illustrious", Width: 512, Height: 512,
		Steps: 4, BatchCount: 1, Family: "sdxl", Mode: "img2img",
		WeightType:       "q8_0",
		CheckpointBytes:  4897 * bytesPerMB,
		ConditionerBytes: 1559 * bytesPerMB,
		VAEBytes:         159 * bytesPerMB,
		ControlNetBytes:  4772 * bytesPerMB,
		RequireGPU:       true,
	}
	decision, lease := authority.Admit(req)
	if !decision.Admitted {
		t.Fatalf("base plus tiled ControlNet stage should fit with CPU ControlNet: %s", decision.Reason)
	}
	defer lease.Release()
	if !decision.Estimate.TilingPlanned {
		t.Fatal("base untiled peak exceeds the budget; tiling must still be planned")
	}
	if !decision.Estimate.ControlNetCPUPlanned {
		t.Fatal("staging ControlNet params crosses the budget; CPU placement must be planned")
	}
	if decision.Estimate.ControlNetVRAMMB != 1431 {
		t.Fatalf("ControlNet staged params = %d MB, want 1431 MB", decision.Estimate.ControlNetVRAMMB)
	}
	if decision.Estimate.ProjectedPeakVRAMMB != 5092 {
		t.Fatalf("projected peak = %d MB, want 5092 MB without staged ControlNet", decision.Estimate.ProjectedPeakVRAMMB)
	}
	if !strings.Contains(decision.Reason, "ControlNet CPU placement planned") {
		t.Fatalf("reason must report ControlNet CPU placement: %s", decision.Reason)
	}
}

func TestResidentLoadedDeltaDoesNotRechargeModuleReserve(t *testing.T) {
	// After stage 3 the run measured only 2003 MB free. The DGL stage reuses the
	// resident upscaler and needs 104 MB of compute, not another 2048 MB module
	// reserve; recharging it refused a stage the hardware could run.
	authority := NewAuthority(fixedHardware(testSnapshot(2003, 8192)), fixedLedger())
	decision, lease := authority.Admit(AdmissionRequest{
		Kind: ResourceImage, Name: "illustrious", Width: 512, Height: 512,
		Mode: "dgl_upscale", RequireGPU: true, ResidentLoaded: true, VAETiling: true,
	})
	if !decision.Admitted {
		t.Fatalf("resident-loaded DGL delta should fit: %s", decision.Reason)
	}
	defer lease.Release()
	if decision.Estimate.RequiredVRAMMB != 0 {
		t.Fatalf("delta required = %d MB, want 0 MB (only the engine overhead is charged)", decision.Estimate.RequiredVRAMMB)
	}
	if decision.Estimate.ActivationVRAMMB != 0 {
		t.Fatalf("DGL activation = %d MB, want 0 MB (no VAE decode)", decision.Estimate.ActivationVRAMMB)
	}
	if decision.Estimate.ProjectedPeakVRAMMB != 512 {
		t.Fatalf("delta projected peak = %d MB, want the 512 MB engine overhead", decision.Estimate.ProjectedPeakVRAMMB)
	}
	// A delta must still hold the image device floor (819 MB on an 8 GiB card);
	// a zero RequiredVRAMMB is not a reason to spend that headroom.
	if decision.Estimate.ReserveVRAMMB != 819 || decision.Estimate.AvailableVRAMMB != 1184 {
		t.Fatalf("delta reserve/available = %d/%d, want 819/1184", decision.Estimate.ReserveVRAMMB, decision.Estimate.AvailableVRAMMB)
	}
}

func TestImageAdmissionRefusesWhenEvenTiledProjectionExceedsBudget(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(5000, 8192)), fixedLedger())
	decision, _ := authority.Admit(AdmissionRequest{
		Kind: ResourceImage, Name: "illustrious", Width: 512, Height: 512,
		Steps: 4, BatchCount: 1, Family: "sdxl", Mode: "txt2img",
		WeightType:       "q8_0",
		CheckpointBytes:  4897 * bytesPerMB,
		ConditionerBytes: 1559 * bytesPerMB,
		VAEBytes:         159 * bytesPerMB,
		RequireGPU:       true,
	})
	if decision.Admitted {
		t.Fatal("even the tiled projection exceeds the budget; admission must refuse")
	}
	if decision.Estimate.TilingPlanned {
		t.Fatal("tiling must not be planned when it still exceeds the budget")
	}
	if decision.Estimate.ProjectedPeakVRAMMB != 6532 {
		t.Fatalf("refusal projection = %d MB, want 6532 MB", decision.Estimate.ProjectedPeakVRAMMB)
	}
	if !strings.Contains(decision.Reason, "projected peak 6532 MB") {
		t.Fatalf("refusal must report the projection: %s", decision.Reason)
	}
}

func TestDefaultReserveAllowsBoundedSoftDip(t *testing.T) {
	// 10% of 7000 is 700, so the preferred budget is 6300. A 6800 MB request
	// spends exactly the 500 MB soft allowance and keeps 200 MB free; 6801 MB
	// would leave less than the 128 MB hard floor and must refuse.
	authority := NewAuthority(fixedHardware(testSnapshot(7000, 8192)), fixedLedger())
	req := AdmissionRequest{Kind: ResourceASR, Name: "soft-fit", AudioVRAMMB: 6800, RequireGPU: true}
	decision, lease := authority.Admit(req)
	if !decision.Admitted {
		t.Fatalf("a request inside the soft allowance should run: %s", decision.Reason)
	}
	lease.Release()
	if decision.Estimate.SoftReserveAllowanceMB != 500 || decision.Estimate.SoftReserveUsedMB != 500 {
		t.Fatalf("soft allowance/use = %d/%d MB, want 500/500",
			decision.Estimate.SoftReserveAllowanceMB, decision.Estimate.SoftReserveUsedMB)
	}
	if decision.Estimate.ReserveVRAMMB != 700 || decision.Estimate.AvailableVRAMMB != 6300 {
		t.Fatalf("preferred reserve/available = %d/%d, want 700/6300",
			decision.Estimate.ReserveVRAMMB, decision.Estimate.AvailableVRAMMB)
	}
	if !strings.Contains(decision.Reason, "soft reserve") || !strings.Contains(decision.Reason, "500 MB") ||
		!strings.Contains(decision.Reason, "200 MB free") {
		t.Fatalf("grant must report the soft-reserve use and remaining headroom: %s", decision.Reason)
	}
	if decision.Warning == "" {
		t.Fatal("soft-reserve grant must carry the user-facing warning")
	}

	req.AudioVRAMMB = 6801
	decision, lease = authority.Admit(req)
	if decision.Admitted {
		t.Fatal("6801 MB spends more than the 500 MB allowance and must be a hard refusal")
	}
	if !strings.Contains(decision.Reason, "soft reserve allowance is not enough") {
		t.Fatalf("hard refusal must explain the soft cap: %s", decision.Reason)
	}
}

func TestExplicitReserveIsNotSoftened(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(7000, 8192)), fixedLedger())
	for _, explicit := range []bool{true, false} {
		decision, _ := authority.Admit(AdmissionRequest{
			Kind: ResourceASR, Name: "explicit", AudioVRAMMB: 6500, RequireGPU: true,
			ReserveOverrideMB: 700, ReserveExplicit: explicit,
		})
		if decision.Admitted {
			t.Fatalf("a measured reserve override is hard (explicit=%t); the request must not dip into it", explicit)
		}
		if decision.Estimate.SoftReserveAllowanceMB != 0 {
			t.Fatalf("override allowance (explicit=%t) = %d MB, want 0", explicit, decision.Estimate.SoftReserveAllowanceMB)
		}
	}
}

func TestImageStageUsesSoftAllowanceOnlyAfterTiling(t *testing.T) {
	// Resident 5580 MB + tiled activation 480 MB + 512 MB overhead = 6572 MB.
	// That is above the preferred 6181 MB budget but inside the 6681 MB soft
	// budget. Admission must tile first, then use 391 MB of the allowance.
	authority := NewAuthority(fixedHardware(testSnapshot(7000, 8192)), fixedLedger())
	decision, lease := authority.Admit(AdmissionRequest{
		Kind: ResourceImage, Name: "soft-stage", Width: 512, Height: 512,
		Steps: 4, BatchCount: 1, Family: "sdxl", Mode: "txt2img",
		WeightType: "q8_0", CheckpointBytes: 7000 * bytesPerMB,
		ConditionerBytes: 2000 * bytesPerMB, RequireGPU: true,
	})
	if !decision.Admitted {
		t.Fatalf("the tiled stage fits the soft allowance: %s", decision.Reason)
	}
	defer lease.Release()
	if !decision.Estimate.TilingPlanned {
		t.Fatal("admission must still prefer VAE tiling before spending soft reserve")
	}
	if decision.Estimate.ProjectedPeakVRAMMB != 6572 {
		t.Fatalf("projected peak = %d MB, want 6572 MB", decision.Estimate.ProjectedPeakVRAMMB)
	}
	if decision.Estimate.SoftReserveUsedMB != 391 {
		t.Fatalf("soft used = %d MB, want 391 MB", decision.Estimate.SoftReserveUsedMB)
	}
}

func TestImageReserveDoesNotShrinkBelowDeviceFloorWhenFreeIsTight(t *testing.T) {
	// 265 MB free on an 8 GiB card must reserve the 819 MB device floor, not
	// 26 MB (10% of the moment's free value). Otherwise a small stage could be
	// admitted into headroom the image reclaimer is meant to protect.
	authority := NewAuthority(fixedHardware(testSnapshot(265, 8192)), fixedLedger())
	decision, _ := authority.Admit(AdmissionRequest{
		Kind: ResourceImage, Name: "pressure", Width: 64, Height: 64,
		CheckpointBytes: 100 * bytesPerMB, RequireGPU: true,
	})
	if decision.Admitted {
		t.Fatal("a 265 MB-free card must not admit an image stage against a reserve below 819 MB")
	}
	if decision.Estimate.ReserveVRAMMB != 819 || decision.Estimate.AvailableVRAMMB != 0 {
		t.Fatalf("reserve/available = %d/%d, want 819/0", decision.Estimate.ReserveVRAMMB, decision.Estimate.AvailableVRAMMB)
	}
}

func TestImageDefaultReserveIsTenPercent(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(10000, 16384)), fixedLedger())
	decision, lease := authority.Admit(AdmissionRequest{
		Kind: ResourceImage, Name: "sdxl", Width: 512, Height: 512,
		CheckpointBytes: 1000 * bytesPerMB, RequireGPU: true,
	})
	if !decision.Admitted {
		t.Fatalf("small image should fit: %s", decision.Reason)
	}
	lease.Release()
	if decision.Estimate.ReserveVRAMMB != 1638 {
		t.Fatalf("default image reserve = %d MB, want 10%% of the 16384 MB device = 1638 MB", decision.Estimate.ReserveVRAMMB)
	}
}
