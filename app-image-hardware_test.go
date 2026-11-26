package main

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"howl-chat/internal/image/types"
	"howl-chat/internal/resource"
)

func testImageSafetensors(t *testing.T, diffusionMB, conditionerMB, vaeMB int64) string {
	t.Helper()
	const mb = int64(1024 * 1024)
	tensors := map[string]interface{}{
		"model.diffusion_model.weight": map[string]interface{}{
			"dtype": "F16", "shape": []int{1}, "data_offsets": []int64{0, diffusionMB * mb},
		},
		"conditioner.weight": map[string]interface{}{
			"dtype": "F16", "shape": []int{1}, "data_offsets": []int64{diffusionMB * mb, (diffusionMB + conditionerMB) * mb},
		},
		"first_stage_model.weight": map[string]interface{}{
			"dtype": "F16", "shape": []int{1}, "data_offsets": []int64{(diffusionMB + conditionerMB) * mb, (diffusionMB + conditionerMB + vaeMB) * mb},
		},
	}
	raw, err := json.Marshal(tensors)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "sdxl.safetensors")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := binary.Write(f, binary.LittleEndian, uint64(len(raw))); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(raw); err != nil {
		t.Fatal(err)
	}
	return path
}

func testImageApp(freeMB, totalMB int64) *App {
	snapshot := resource.HardwareSnapshot{
		GPUModel:       "Test GPU",
		TotalVRAMBytes: totalMB * 1024 * 1024,
		FreeVRAMBytes:  freeMB * 1024 * 1024,
		UsedVRAMBytes:  (totalMB - freeMB) * 1024 * 1024,
		TotalRAMBytes:  32 * 1024 * 1024 * 1024,
		FreeRAMBytes:   16 * 1024 * 1024 * 1024,
		AvailableRAM:   16 * 1024 * 1024 * 1024,
		VRAMLive:       true,
		VRAMSource:     "test-live",
		GPUAvailable:   true,
	}
	return &App{admission: resource.NewAuthority(
		func() (resource.HardwareSnapshot, error) { return snapshot, nil },
		func() []resource.ModelState { return nil },
	)}
}

func TestImageAdmissionAutoQuantizesOnlyWhenStoredPrecisionDoesNotFit(t *testing.T) {
	path := testImageSafetensors(t, 4897, 1559, 159)
	app := testImageApp(6901, 8192)
	loadOpts := types.ModelLoadOptions{ModelName: "sdxl", ModelPath: path}
	gen := types.ImageOptions{Width: 1024, Height: 1024, BatchCount: 1, Steps: 20}

	decision, lease, chosen, err := app.admitImageLoad(loadOpts, gen, true)
	if err != nil {
		t.Fatalf("q8 fallback should admit the model: %v", err)
	}
	if lease == nil {
		t.Fatal("an admitted load must hold a lease")
	}
	lease.Release()
	if chosen.WeightType != "q8_0" {
		t.Fatalf("chosen type = %q, want q8_0", chosen.WeightType)
	}
	if loadOpts.WeightType != "" {
		t.Fatalf("automatic quantization must not mutate the caller's explicit type: %q", loadOpts.WeightType)
	}
	if decision.Estimate.ActivationVRAMMB <= 0 {
		t.Fatal("admission must report the workload activation estimate")
	}
	if !chosen.VAETiling {
		t.Fatal("1024x1024 decode exceeds the post-q8 headroom; tiling should be pre-enabled")
	}

	if imageAutoQuantizeEligible(types.ModelLoadOptions{ModelPath: path, WeightType: "f16"}) {
		t.Fatal("an explicit f16 choice must never be auto-converted")
	}
	if imageAutoQuantizeEligible(types.ModelLoadOptions{ModelPath: filepath.Join(t.TempDir(), "model.gguf")}) {
		t.Fatal("a GGUF must not be automatically re-quantized")
	}
}

func TestImageAdmissionZeroMaxVRAMKeepsReserve(t *testing.T) {
	path := testImageSafetensors(t, 4897, 1559, 159)
	app := testImageApp(6901, 8192)
	zero := 0

	_, lease, chosen, err := app.admitImageLoad(
		types.ModelLoadOptions{ModelName: "sdxl", ModelPath: path, MaxVRAM: &zero},
		types.ImageOptions{Width: 512, Height: 512, BatchCount: 1, Steps: 20},
		true,
	)
	if err != nil {
		t.Fatalf("zero max-vram should still admit through q8: %v", err)
	}
	lease.Release()
	if chosen.MaxVRAM == nil || *chosen.MaxVRAM == 0 || !chosen.AutoMaxVRAM {
		t.Fatalf("max-vram 0 must become a positive post-reserve cap, got %v (auto=%t)", chosen.MaxVRAM, chosen.AutoMaxVRAM)
	}
	if *chosen.MaxVRAM*1024 > 6901-690 {
		t.Fatalf("auto cap %d GiB exceeds the post-reserve budget", *chosen.MaxVRAM)
	}

	eight := 8
	_, lease, chosen, err = app.admitImageLoad(
		types.ModelLoadOptions{ModelName: "sdxl", ModelPath: path, MaxVRAM: &eight},
		types.ImageOptions{Width: 512, Height: 512, BatchCount: 1, Steps: 20},
		true,
	)
	if err != nil {
		t.Fatalf("an oversized explicit cap should be clamped, not refused: %v", err)
	}
	lease.Release()
	if chosen.MaxVRAM == nil || *chosen.MaxVRAM >= 8 || !chosen.AutoMaxVRAM {
		t.Fatalf("explicit max-vram 8 must be clamped to the post-reserve budget, got %v (auto=%t)", chosen.MaxVRAM, chosen.AutoMaxVRAM)
	}
}

func TestImageBudgetAutomaticallySelectsBoundedPlacement(t *testing.T) {
	path := testImageSafetensors(t, 4897, 1559, 159)
	app := testImageApp(6901, 8192)
	gen := types.ImageOptions{ModelName: "sdxl", Width: 512, Height: 512, BatchCount: 1, Steps: 1}

	// Offload is not a prerequisite: a budget below the resident need is enough
	// for Howl to select the measured GPU-assisted bounded mode itself.
	three := 3
	opts := types.ModelLoadOptions{ModelName: "sdxl", ModelPath: path, MaxVRAM: &three}
	if err := app.applyConstrainedOffload(&opts, gen); err != nil {
		t.Fatalf("3 GiB should select the bounded mode: %v", err)
	}
	if opts.Backend != imageConstrainedBackend || !opts.OffloadToCPU || !opts.VAETiling || !opts.AutoPlacement {
		t.Fatalf("3 GiB did not select automatic CPU diffusion + VAE tiling: %+v", opts)
	}

	decision, lease, chosen, err := app.admitImageLoad(opts, gen, true)
	if err != nil {
		t.Fatalf("3 GiB bounded admission should succeed: %v", err)
	}
	if !decision.Admitted || lease == nil {
		t.Fatal("3 GiB bounded admission must hold a lease")
	}
	lease.Release()
	if chosen.Backend != imageConstrainedBackend || !chosen.VAETiling || !chosen.AutoPlacement {
		t.Fatalf("admission dropped the bounded placement: %+v", chosen)
	}

	// Below the measured GPU-assisted process need, the honest way to honor a
	// hard budget is no GPU compute at all. It is slower, but it starts instead
	// of refusing a model that can run.
	two := 2
	cpuOnly := types.ModelLoadOptions{ModelName: "sdxl", ModelPath: path, MaxVRAM: &two}
	decision, lease, chosen, err = app.admitImageLoad(cpuOnly, gen, true)
	if err != nil {
		t.Fatalf("2 GiB should select CPU-only, not refuse: %v", err)
	}
	if !decision.Admitted || lease == nil {
		t.Fatal("2 GiB CPU-only admission must hold a lease")
	}
	lease.Release()
	if chosen.Backend != "cpu" || chosen.AutoPlacementMode != "cpu-only" || chosen.OffloadToCPU {
		t.Fatalf("2 GiB did not select CPU-only placement: %+v", chosen)
	}
	if decision.Estimate.RequiredVRAMMB != 0 {
		t.Fatalf("CPU-only must require no VRAM, got %d MB", decision.Estimate.RequiredVRAMMB)
	}
}

// An automatic post-reserve cap is Howl's policy, not a user budget. It must
// not trigger the hard-budget CPU fallback on the next generation: that was
// making a resident IP workflow restart with diffusion on CPU despite a just-
// measured GPU placement that fit.
func TestImageAutoMaxVRAMDoesNotSelectConstrainedOffload(t *testing.T) {
	path := testImageSafetensors(t, 3278, 469, 319)
	app := testImageApp(6901, 8192)
	gen := types.ImageOptions{ModelName: "sd15", Width: 512, Height: 512, BatchCount: 1, Steps: 8}
	six := 6
	opts := types.ModelLoadOptions{ModelName: "sd15", ModelPath: path, MaxVRAM: &six, AutoMaxVRAM: true}
	if err := app.applyConstrainedOffload(&opts, gen); err != nil {
		t.Fatalf("auto cap should not error: %v", err)
	}
	if opts.Backend != "" || opts.OffloadToCPU || opts.AutoPlacementMode != "" {
		t.Fatalf("automatic cap was treated as a hard budget: %+v", opts)
	}
}

func TestImageAdmissionKeepsResidencyManualWhenOffloaded(t *testing.T) {
	path := testImageSafetensors(t, 1000, 200, 100)
	app := testImageApp(6901, 8192)

	_, lease, chosen, err := app.admitImageLoad(
		types.ModelLoadOptions{ModelName: "sdxl", ModelPath: path, OffloadToCPU: true},
		types.ImageOptions{Width: 512, Height: 512, BatchCount: 1, Steps: 20},
		true,
	)
	if err != nil {
		t.Fatalf("offloaded load should admit: %v", err)
	}
	lease.Release()
	if chosen.WeightResidencyTTL != nil {
		t.Fatalf("residency TTL must stay manual to preserve reusable weights, got %v", *chosen.WeightResidencyTTL)
	}
	if chosen.MaxVRAM == nil || !chosen.AutoMaxVRAM {
		t.Fatal("offload without an explicit budget must cap GPU memory automatically")
	}
	if *chosen.MaxVRAM <= 0 || *chosen.MaxVRAM*1024 > 6901-690 {
		t.Fatalf("automatic offload budget %d GiB exceeds the post-reserve budget", *chosen.MaxVRAM)
	}
	if chosen.VAETiling {
		t.Fatal("512x512 activation fits the headroom and should not force tiling")
	}
}

func TestParseModelLoadOptsCarriesHardwareTruth(t *testing.T) {
	opts, err := parseModelLoadOpts(map[string]interface{}{
		"model":                "sdxl",
		"backend":              "cpu",
		"params_backend":       "disk",
		"flash_attention_mode": "off",
		"weight_residency_ttl": float64(0),
		"circular":             true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Backend != "cpu" || opts.ParamsBackend != "disk" {
		t.Fatalf("placement flags lost: %+v", opts)
	}
	if opts.FlashAttentionMode != "off" {
		t.Fatalf("flash attention mode = %q, want off", opts.FlashAttentionMode)
	}
	if opts.WeightResidencyTTL == nil || *opts.WeightResidencyTTL != 0 {
		t.Fatalf("residency TTL lost: %v", opts.WeightResidencyTTL)
	}
	if !opts.Circular {
		t.Fatal("circular flag lost")
	}
}

func TestImageDecodeHeadroomRule(t *testing.T) {
	if !imageDecodeExceedsHeadroom(resource.AdmissionEstimate{RequiredVRAMMB: 5000, AvailableVRAMMB: 6000, ActivationVRAMMB: 1200}) {
		t.Fatal("activation above headroom must request tiling")
	}
	if imageDecodeExceedsHeadroom(resource.AdmissionEstimate{RequiredVRAMMB: 5000, AvailableVRAMMB: 6000, ActivationVRAMMB: 900}) {
		t.Fatal("activation within headroom must not force tiling")
	}
	if imageDecodeExceedsHeadroom(resource.AdmissionEstimate{RequiredVRAMMB: 5000, AvailableVRAMMB: 6000}) {
		t.Fatal("a load with no generation workload must not force tiling")
	}
}
