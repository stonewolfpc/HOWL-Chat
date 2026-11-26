package main

import (
	"errors"
	"strings"
	"testing"

	imagetypes "howl-chat/internal/image/types"
	"howl-chat/internal/media"
	"howl-chat/internal/resource"
)

func TestImageReclaimBandIsWorkloadAware(t *testing.T) {
	gpu := imagetypes.ModelLoadOptions{}
	partial := imagetypes.ModelLoadOptions{AutoPlacementMode: "cpu-diffusion", OffloadToCPU: true}
	cpu := imagetypes.ModelLoadOptions{Backend: "cpu"}

	small := imageReclaimBandMB(551, 6624, 3, 1, gpu)
	large := imageReclaimBandMB(2880, 6624, 3, 1, gpu)
	if large <= small {
		t.Fatalf("larger resolution band = %.0fMB, want greater than %.0fMB", large, small)
	}

	single := imageReclaimBandMB(551, 6624, 3, 1, gpu)
	batch := imageReclaimBandMB(551, 6624, 3, 2, gpu)
	if batch <= single {
		t.Fatalf("batch band = %.0fMB, want greater than single band %.0fMB", batch, single)
	}
	if batch-single > imageReclaimMaxBatchAllowanceMB+0.01 {
		t.Fatalf("batch allowance %.0fMB exceeds cap %.0fMB", batch-single, imageReclaimMaxBatchAllowanceMB)
	}

	if imageReclaimBandMB(551, 6624, 3, 1, partial) <= small {
		t.Fatal("partial CPU placement must have a larger band than GPU-resident placement")
	}
	if imageReclaimBandMB(551, 6624, 3, 1, cpu) <= imageReclaimBandMB(551, 6624, 3, 1, partial) {
		t.Fatal("CPU-only placement must have a larger band than partial placement")
	}
	if small < imageReclaimMinimumBandMB {
		t.Fatalf("minimum band not applied: %.0fMB", small)
	}
}

func TestImageReclaimReserveReasonWinsOverInsideBand(t *testing.T) {
	reason := imageReclaimReason(nil, false, 800, 819, 12, 1016)
	if reason == "" {
		t.Fatal("free VRAM below reserve did not reclaim")
	}
	if !strings.Contains(reason, "free VRAM") {
		t.Fatalf("reason = %q, want reserve attribution", reason)
	}
	if strings.Contains(reason, "working-set growth") {
		t.Fatalf("reason = %q, must not blame growth when growth is inside the band", reason)
	}
}

func TestImageReclaimReserveReasonAtExactReserve(t *testing.T) {
	reason := imageReclaimReason(nil, false, 819, 819, 0, 1016)
	if !strings.Contains(reason, "free VRAM") {
		t.Fatalf("reason = %q, want reclaim at the exact reserve boundary", reason)
	}
}

func TestImageReclaimGrowthReasonWhenReserveIsHeld(t *testing.T) {
	reason := imageReclaimReason(nil, false, 2000, 819, 1200, 1016)
	if !strings.Contains(reason, "working-set growth") {
		t.Fatalf("reason = %q, want growth attribution", reason)
	}
}

func TestImageReclaimFailureReasonWins(t *testing.T) {
	reason := imageReclaimReason(errors.New("engine abandoned"), true, 2000, 819, 0, 1016)
	if !strings.Contains(reason, "abandoned") {
		t.Fatalf("reason = %q, want failure attribution", reason)
	}
}

func TestImageReclaimSummaryKeepsLastReserveResult(t *testing.T) {
	p := &imageReclaimPolicy{
		last: &imageReclaimEvent{Reclaimed: false, GrowthMB: 8, BandMB: 1272, BaselineMB: 6892, ThresholdMB: 8164, MeasuredMB: 6900},
		events: []imageReclaimEvent{{
			Reclaimed: true, Reason: "free VRAM 484MB fell to or below the 819MB (10%) reserve",
			FreeBeforeMB: 484, FreeAfterMB: 7091, ReserveMB: 819,
			WorkingSetBeforeMB: 6895, WorkingSetAfterMB: 171,
		}},
	}
	summary := p.summaryLocked()
	for _, want := range []string{"484", "7091", "819", "6895", "171"} {
		if !strings.Contains(summary, want) {
			t.Fatalf("summary %q does not retain %s", summary, want)
		}
	}
}

func TestWorkflowParentLeaseRefusalIsStructured(t *testing.T) {
	app := &App{images: media.NewImage(t.TempDir())}
	lease := resource.NewWorkflowLease("engine-a")
	claim, err := lease.Acquire(resource.NewComponentKey(
		"engine-a:42", "upscaler", "/models/upscaler.pth", "f16", nil,
	))
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	app.imageWorkflow = lease
	t.Cleanup(func() {
		claim.Release()
		app.endImageWorkflow(lease)
	})

	app.reclaimImageEngineAfterGenerationLocked("", imagetypes.ImageOptions{Mode: "txt2img"}, nil)

	last := app.imageReclaim.last
	if last == nil || !last.Blocked {
		t.Fatalf("reclaim refusal was not recorded: %+v", last)
	}
	if last.BlockedReason == "" || last.WorkflowEngine != "engine-a" || last.WorkflowActiveClaims != 1 {
		t.Fatalf("blocked event lacks context: %+v", last)
	}
	report := app.imageReclaim.report()
	if report["blocked"] != true {
		t.Fatalf("report blocked = %v, want true: %+v", report["blocked"], report)
	}
	if report["blockedReason"] != "workflow parent lease active" {
		t.Fatalf("report blockedReason = %v", report["blockedReason"])
	}
	summary, _ := report["summary"].(string)
	if !strings.Contains(summary, "Reclamation blocked") || !strings.Contains(summary, "engine-a") {
		t.Fatalf("blocked summary = %q", summary)
	}
}

func TestImageReclaimBandAccountsForComponents(t *testing.T) {
	gpu := imagetypes.ModelLoadOptions{}
	one := imageReclaimBandMB(551, 3300, 1, 1, gpu)
	three := imageReclaimBandMB(551, 6600, 3, 1, gpu)
	if three <= one {
		t.Fatalf("more active components must raise the band: one=%.0fMB three=%.0fMB", one, three)
	}
}
