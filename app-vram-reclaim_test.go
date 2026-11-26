package main

import (
	"strings"
	"testing"
	"time"

	"howl-chat/internal/image/types"
	"howl-chat/internal/resource"
)

// The reclaim wait must judge the same ladder the load path will try. Here the
// stored precision does not fit, the automatic q8_0 conversion does, and the
// handoff must report success rather than timing out on a condition the start
// would never ask for.
func TestReclaimSucceedsOnTheSameLadderTheLoadUses(t *testing.T) {
	path := testImageSafetensors(t, 4897, 1559, 159)
	// 7000 MB free: the stored precision does not fit, the q8_0 conversion does.
	app := testImageApp(7000, 8192)
	gen := types.ImageOptions{ModelName: "sdxl", Width: 512, Height: 512, BatchCount: 1, Steps: 1}

	loadOpts := types.ModelLoadOptions{ModelName: "sdxl", ModelPath: path}
	full := app.imageAdmissionRequest(loadOpts, gen, false)
	q8, ok := imageAutoQuantizeRequest(full, loadOpts, true)
	if !ok {
		t.Fatal("test setup: q8 fallback should be eligible for a safetensors file")
	}
	if decision, lease := app.admission.Admit(full); decision.Admitted {
		lease.Release()
		t.Fatal("test setup: the full-precision request should not fit")
	}

	fit, err := app.waitForAdmissionFit("image", []resource.AdmissionRequest{full, q8})
	if err != nil {
		t.Fatalf("the handoff must not report a timeout when the load's own fallback fits: %v", err)
	}
	if !fit.CanFit || fit.WeightType != "q8_0" {
		t.Fatalf("reported placement = %+v, want the admitted q8_0 fallback", fit)
	}
	if !strings.Contains(fit.Message, "q8_0") {
		t.Fatalf("the handoff must say which placement it admitted: %s", fit.Message)
	}
}

// Strictness is preserved: when nothing on the ladder fits, the handoff still
// fails and says so. The timeout was never the thing being fixed.
func TestReclaimStillTimesOutWhenNothingFits(t *testing.T) {
	previous := vramReleaseTimeout
	vramReleaseTimeout = 200 * time.Millisecond
	t.Cleanup(func() { vramReleaseTimeout = previous })

	path := testImageSafetensors(t, 4897, 1559, 159)
	app := testImageApp(1860, 8192)
	gen := types.ImageOptions{ModelName: "sdxl", Width: 512, Height: 512, BatchCount: 1, Steps: 1}
	req := app.imageAdmissionRequest(types.ModelLoadOptions{ModelName: "sdxl", ModelPath: path}, gen, false)

	if _, err := app.waitForAdmissionFit("image", []resource.AdmissionRequest{req}); err == nil {
		t.Fatal("a handoff that frees too little must still report failure")
	}
}
