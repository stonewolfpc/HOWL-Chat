package resource

import (
	"strings"
	"testing"

	backendtypes "howl-chat/internal/backend/types"
)

func testSnapshot(freeMB, totalMB int64) HardwareSnapshot {
	return HardwareSnapshot{
		GPUModel:       "Test GPU",
		TotalVRAMBytes: totalMB * bytesPerMB,
		FreeVRAMBytes:  freeMB * bytesPerMB,
		UsedVRAMBytes:  (totalMB - freeMB) * bytesPerMB,
		TotalRAMBytes:  32 * 1024 * bytesPerMB,
		FreeRAMBytes:   16 * 1024 * bytesPerMB,
		AvailableRAM:   16 * 1024 * bytesPerMB,
		VRAMLive:       true,
		VRAMSource:     "test-live",
		GPUAvailable:   true,
	}
}

func fixedHardware(snap HardwareSnapshot) func() (HardwareSnapshot, error) {
	return func() (HardwareSnapshot, error) { return snap, nil }
}

func fixedLedger(models ...ModelState) func() []ModelState {
	return func() []ModelState { return models }
}

func TestAdmissionReserveBoundary(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(10000, 24000)), fixedLedger())

	// 10% of 10000 is 1000, leaving 9000. Exactly 9000 is admitted.
	decision, lease := authority.Admit(AdmissionRequest{
		Kind: ResourceTTS, Name: "boundary", AudioVRAMMB: 9000,
	})
	if !decision.Admitted {
		t.Fatalf("9000 MB should fit in 9000 MB available after reserve: %s", decision.Reason)
	}
	if decision.Estimate.ReserveVRAMMB != 1000 || decision.Estimate.AvailableVRAMMB != 9000 {
		t.Fatalf("unexpected boundary numbers: reserve=%d available=%d",
			decision.Estimate.ReserveVRAMMB, decision.Estimate.AvailableVRAMMB)
	}
	lease.Release()

	// The default reserve has a 500 MB soft allowance: 9500 MB spends all of
	// it and keeps 500 MB free. One more megabyte is the hard boundary.
	decision, lease = authority.Admit(AdmissionRequest{
		Kind: ResourceTTS, Name: "soft-boundary", AudioVRAMMB: 9500,
	})
	if !decision.Admitted {
		t.Fatalf("9500 MB is exactly the soft boundary and should run: %s", decision.Reason)
	}
	if decision.Estimate.SoftReserveUsedMB != 500 {
		t.Fatalf("soft used = %d MB, want 500 MB", decision.Estimate.SoftReserveUsedMB)
	}
	lease.Release()

	decision, _ = authority.Admit(AdmissionRequest{
		Kind: ResourceTTS, Name: "over-boundary", AudioVRAMMB: 9501,
	})
	if decision.Admitted {
		t.Fatal("9501 MB exceeds the 500 MB soft allowance and must be refused")
	}
	if !strings.Contains(decision.Reason, "9501 MB") ||
		!strings.Contains(decision.Reason, "1000 MB reserve") ||
		!strings.Contains(decision.Reason, "9000 MB available") ||
		!strings.Contains(decision.Reason, "soft reserve allowance is not enough") {
		t.Fatalf("refusal must name exact required/reserve/available numbers and the soft cap: %s", decision.Reason)
	}
}

func TestAdmissionFailsClosedOnDetectionFailure(t *testing.T) {
	snap := testSnapshot(10000, 24000)
	snap.DetectionError = "nvidia-smi failed"
	authority := NewAuthority(fixedHardware(snap), fixedLedger())

	decision, _ := authority.Admit(AdmissionRequest{Kind: ResourceTTS, Name: "detect", AudioVRAMMB: 100})
	if decision.Admitted {
		t.Fatal("admission must refuse when hardware detection reports an error")
	}
	if !strings.Contains(decision.Reason, "hardware detection failed") {
		t.Fatalf("refusal must expose the detection failure: %s", decision.Reason)
	}

	nonLive := testSnapshot(10000, 24000)
	nonLive.VRAMLive = false
	nonLive.VRAMSource = "registry"
	authority = NewAuthority(fixedHardware(nonLive), fixedLedger())
	decision, _ = authority.Admit(AdmissionRequest{Kind: ResourceTTS, Name: "stale", AudioVRAMMB: 100})
	if decision.Admitted {
		t.Fatal("admission must refuse a GPU request whose free memory was not read live")
	}
	if !strings.Contains(decision.Reason, "could not be read live") {
		t.Fatalf("refusal must name the non-live source: %s", decision.Reason)
	}
}

func TestAdmissionRefreshesHardwareForEveryDecision(t *testing.T) {
	calls := 0
	authority := NewAuthority(func() (HardwareSnapshot, error) {
		calls++
		if calls == 1 {
			return testSnapshot(5000, 24000), nil
		}
		return testSnapshot(20000, 24000), nil
	}, fixedLedger())

	req := AdmissionRequest{Kind: ResourceTTS, Name: "refresh", AudioVRAMMB: 9000}
	first, _ := authority.Admit(req)
	if first.Admitted {
		t.Fatal("stale 5000 MB snapshot must not admit a 9000 MB request")
	}
	if calls != 1 {
		t.Fatalf("authority must refresh once per decision, calls=%d", calls)
	}

	second, lease := authority.Admit(req)
	if !second.Admitted {
		t.Fatalf("refreshed 20000 MB snapshot should admit the request: %s", second.Reason)
	}
	if calls != 2 {
		t.Fatalf("authority must not reuse the stale snapshot, calls=%d", calls)
	}
	lease.Release()
}

func TestChatAdmissionRefusedWhenImageModelLoaded(t *testing.T) {
	loaded := []ModelState{{
		Type: ResourceImage, Name: "sdxl", Loaded: true, GPUMemoryMB: 6000, SystemRAMMB: 2000,
	}}
	authority := NewAuthority(fixedHardware(testSnapshot(8000, 24000)), fixedLedger(loaded...))

	req := AdmissionRequest{
		Kind:        ResourceLLM,
		Name:        "chat",
		ContextSize: 4096,
		GPULayers:   32,
		BatchSize:   512,
		RequireGPU:  true,
		Arch: &backendtypes.ArchitectureSpec{
			Family: "llama", Layers: 32, DModel: 4096, NumHeads: 32, NumKVHeads: 8,
			HeadDim: 128, FFNSize: 11008, ContextLength: 8192,
			FileSizeBytes: 20 * 1024 * 1024 * 1024, ParameterCount: 20_000_000_000,
		},
	}
	decision, _ := authority.Admit(req)
	if decision.Admitted {
		t.Fatal("chat load must be refused while the image model holds the card")
	}
	if !strings.Contains(decision.Reason, "image") || !strings.Contains(decision.Reason, "sdxl") ||
		!strings.Contains(decision.Reason, "6000 MB") {
		t.Fatalf("refusal must name the blocking image model: %s", decision.Reason)
	}

	// Once the image load is gone, the same request is admissible.
	authority = NewAuthority(fixedHardware(testSnapshot(48000, 48000)), fixedLedger())
	decision, lease := authority.Admit(req)
	if !decision.Admitted {
		t.Fatalf("chat should load once the image model is gone: %s", decision.Reason)
	}
	lease.Release()
}

func TestASRAdmissionRefusedWhenChatLoaded(t *testing.T) {
	loaded := []ModelState{{
		Type: ResourceLLM, Name: "chat", Loaded: true, GPUMemoryMB: 3000, SystemRAMMB: 4000,
	}}
	authority := NewAuthority(fixedHardware(testSnapshot(4500, 24000)), fixedLedger(loaded...))

	decision, _ := authority.Admit(AdmissionRequest{
		Kind: ResourceASR, Name: "qwen-asr", AudioVRAMMB: 4400, AudioRAMMB: 1024, RequireGPU: true,
	})
	if decision.Admitted {
		t.Fatal("ASR start must be refused when the chat model leaves less than the reserve-aware requirement")
	}
	if !strings.Contains(decision.Reason, "asr") || !strings.Contains(decision.Reason, "chat") {
		t.Fatalf("refusal must name the ASR request and the blocking chat model: %s", decision.Reason)
	}
}

// A checkpoint's own weights are part of its GPU requirement. Before this was
// counted, an 8 GB card with a chat model loaded was told that a 6617 MB SDXL
// checkpoint plus its 1559 MB text encoder "fits", no unload confirmation was
// raised, and sd-server silently fell back to CPU params.
func TestImageEstimateCountsCheckpointAndConditionerWeights(t *testing.T) {
	chat := ModelState{Type: ResourceLLM, Name: "chat.gguf", GPUMemoryMB: 3300, IsGPUOffload: true}
	authority := NewAuthority(fixedHardware(testSnapshot(3500, 8192)), fixedLedger(chat))

	sdxl := AdmissionRequest{
		Kind: ResourceImage, Name: "Illustrious-XL-v2.0",
		Width: 1024, Height: 1024, Steps: 8, BatchCount: 1,
		Family: "sdxl", Mode: "txt2img",
		CheckpointBytes:  6617 * bytesPerMB,
		ConditionerBytes: 1559 * bytesPerMB,
		RequireGPU:       true, ReserveOverrideMB: 1024,
	}
	decision, lease := authority.Admit(sdxl)
	if decision.Admitted {
		t.Fatalf("a 6617 MB checkpoint must not be predicted to fit beside a resident chat model: %s", decision.Reason)
	}
	if lease != nil {
		t.Fatal("a refused image request must not return a lease")
	}
	if decision.Estimate.RequiredVRAMMB < 6617 {
		t.Fatalf("requirement must include the checkpoint bytes, got %d MB", decision.Estimate.RequiredVRAMMB)
	}
	if len(decision.Blocking) == 0 || decision.Blocking[0].Name != "chat.gguf" {
		t.Fatalf("refusal must name the resident chat model as blocking: %+v", decision.Blocking)
	}

	// The same checkpoint cannot fit on this 8 GB card even with the chat model
	// gone, because 8176 MB of weights exceeds the 1 GiB-reserved free memory.
	empty := NewAuthority(fixedHardware(testSnapshot(7200, 8192)), fixedLedger())
	if d, _ := empty.Admit(sdxl); d.Admitted {
		t.Fatalf("8176 MB of weights must not be admitted against 7200 MB free: %s", d.Reason)
	}

	// A checkpoint that does fit still admits, so the change is not a blanket
	// refusal: a Q4 SDXL bundle with the same text encoder fits once the card is
	// free.
	small := sdxl
	small.Name = "Pony-Realism-V23ULTRA-Q4_K_M"
	small.CheckpointBytes = 1557 * bytesPerMB
	decision, lease = empty.Admit(small)
	if !decision.Admitted {
		t.Fatalf("a 3116 MB weight set should fit in 7200 MB free with a 1 GiB reserve: %s", decision.Reason)
	}
	lease.Release()
}

func TestAdmissionImageOverrideUsesOneGiB(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(3000, 24000)), fixedLedger())
	decision, _ := authority.Admit(AdmissionRequest{
		Kind: ResourceImage, Name: "sdxl", Width: 1024, Height: 1024, Steps: 20,
		// A real image request states the checkpoint it will load; without one
		// the requirement would be only engine overhead.
		CheckpointBytes:   2000 * bytesPerMB,
		ReserveOverrideMB: 1024, RequireGPU: true,
	})
	if decision.Admitted {
		t.Fatal("image request must be refused when the 1 GiB sd-server reserve leaves too little")
	}
	if decision.Estimate.ReserveVRAMMB != 1024 {
		t.Fatalf("image reserve override should be 1024 MB, got %d", decision.Estimate.ReserveVRAMMB)
	}
}

// Admission and sd-server must judge the same budget. sd.cpp's backend_fit
// places a module only when weights + 2048 MB reserve fit the --max-vram budget
// (1024 MB for the VAE); the policy authority now applies that rule, so the
// default 1 GiB reserve refuses an 8.6 GB SDXL on an 8 GB card while
// --max-vram 0 admits it - the same placement sd.cpp then reports as CUDA0.
func TestImageAdmissionFollowsSdServerBudget(t *testing.T) {
	sdxl := AdmissionRequest{
		Kind: ResourceImage, Name: "Illustrious-XL-v2.0",
		Width: 1024, Height: 1024, Steps: 8, BatchCount: 1,
		Family: "sdxl", Mode: "txt2img",
		CheckpointBytes:   6617 * bytesPerMB,
		ConditionerBytes:  1559 * bytesPerMB,
		VAEBytes:          319 * bytesPerMB,
		RequireGPU:        true,
		ReserveOverrideMB: 1024,
		ReserveExplicit:   true,
	}
	authority := NewAuthority(fixedHardware(testSnapshot(7200, 8192)), fixedLedger())
	if d, _ := authority.Admit(sdxl); d.Admitted {
		t.Fatalf("with the 1 GiB reserve an 8665 MB DiT requirement must be refused at 7200 MB free: %s", d.Reason)
	}

	// Without quantization the whole checkpoint is one module: admission stays
	// deliberately conservative here because it does not open the file to find
	// sd.cpp's internal DiT/encoder split. The default reserve refuses it too.
	all := sdxl
	all.ReserveOverrideMB = 0 // --max-vram 0
	if d, _ := authority.Admit(all); d.Admitted {
		t.Fatalf("an unquantized 6617 MB checkpoint plus 2048 MB reserve exceeds 7200 MB: %s", d.Reason)
	}

	// With --type q8_0 sd.cpp converts at load and keeps ~53% of the bytes, so
	// policy and placement now agree that the module fits the same budget.
	quantized := all
	quantized.WeightType = "q8_0"
	decision, lease := authority.Admit(quantized)
	if !decision.Admitted {
		t.Fatalf("a q8_0-converted checkpoint fits 7200 MB with --max-vram 0: %s", decision.Reason)
	}
	if decision.Estimate.ReserveVRAMMB != 0 {
		t.Fatalf("an explicit zero reserve must not become the 10%% default, got %d", decision.Estimate.ReserveVRAMMB)
	}
	q8Scale := 0.62 // measured: 3996.05 MB of 6624.11 MB under --type q8_0
	if got, want := decision.Estimate.RequiredVRAMMB, int64(float64(6617)*q8Scale)+2048; got != want {
		t.Fatalf("requirement must be the converted weights plus the 2048 MB reserve (%d), got %d", want, got)
	}
	lease.Release()

	// A positive --max-vram is a ceiling on what sd-server may spend.
	capped := all
	capped.VRAMBudgetMB = 4096
	if d, _ := authority.Admit(capped); d.Admitted {
		t.Fatalf("a 4 GiB --max-vram ceiling must refuse an 8665 MB requirement: %s", d.Reason)
	}
}

// A real SD1.5 IP + ControlNet process measured 3129 MB VRAM at 6 GiB. Lowering
// the budget to 5 GiB must reuse that measurement instead of the pessimistic
// module+2048 MB projection, which would otherwise select CPU diffusion even
// though the actual resident set fits the new cap.
func TestImageAdmissionUsesMeasuredResidentForBudgetChange(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(7100, 8192)), fixedLedger())
	req := AdmissionRequest{
		Kind: ResourceImage, Name: "sd15-ip",
		Width: 256, Height: 256, Steps: 4, BatchCount: 1,
		Family: "sd15", Mode: "txt2img",
		CheckpointBytes: 3278 * bytesPerMB, ConditionerBytes: 469 * bytesPerMB,
		VAEBytes: 319 * bytesPerMB, ControlNetBytes: 1378 * bytesPerMB,
		IPAdapterBytes: 42 * bytesPerMB, ClipVisionBytes: 2411 * bytesPerMB,
		ClipVisionCPU: true, RequireGPU: true, VRAMBudgetMB: 5120,
	}
	if d, _ := authority.Admit(req); d.Admitted {
		t.Fatalf("the conservative module projection must not fit 5 GiB: %s", d.Reason)
	}
	req.MeasuredResidentVRAMMB = 3129
	d, lease := authority.Admit(req)
	if !d.Admitted {
		t.Fatalf("the measured 3129 MB resident set fits the 5 GiB cap: %s", d.Reason)
	}
	if d.Estimate.RequiredVRAMMB != 3129 {
		t.Fatalf("measured requirement must be 3129 MB, got %d", d.Estimate.RequiredVRAMMB)
	}
	if want := int64(3129 + 480 + ImageFixedOverheadMB); d.Estimate.ProjectedPeakVRAMMB != want {
		t.Fatalf("projected peak must use measured resident %d, got %d", want, d.Estimate.ProjectedPeakVRAMMB)
	}
	lease.Release()
}

// An image run the user asked to offload to CPU keeps its weights in system RAM
// and needs only the compute working set on the device. Admission has to judge
// that placement: judging it as a fully resident load refused every model bigger
// than the card even though the engine was already told to offload it.
func TestOffloadedImageAdmissionChargesWeightsToRAM(t *testing.T) {
	snap := testSnapshot(6000, 8192)
	snap.TotalRAMBytes = 64 * 1024 * bytesPerMB
	snap.FreeRAMBytes = 48 * 1024 * bytesPerMB
	snap.AvailableRAM = 48 * 1024 * bytesPerMB
	authority := NewAuthority(fixedHardware(snap), fixedLedger())

	base := AdmissionRequest{
		Kind:             ResourceImage,
		Name:             "big-checkpoint",
		CheckpointBytes:  10 * 1024 * bytesPerMB,
		ConditionerBytes: 5 * 1024 * bytesPerMB,
		VAEBytes:         300 * bytesPerMB,
		RequireGPU:       true,
	}
	resident, _ := authority.Admit(base)
	if resident.Admitted {
		t.Fatalf("a 15 GB workload was admitted as VRAM-resident on a 6 GB free card: %+v", resident.Estimate)
	}

	offloaded := base
	offloaded.RequireGPU = false
	offloaded.WeightsOffloaded = true
	decision, lease := authority.Admit(offloaded)
	if !decision.Admitted {
		t.Fatalf("the same workload was refused after asking to offload to CPU: %s", decision.Reason)
	}
	lease.Release()
	if decision.Estimate.RequiredVRAMMB != sdModuleReserveMB {
		t.Fatalf("offloaded VRAM requirement = %d MB, want only the %d MB compute reserve",
			decision.Estimate.RequiredVRAMMB, sdModuleReserveMB)
	}
	if decision.Estimate.RequiredRAMMB < 15000 {
		t.Fatalf("offloaded RAM requirement = %d MB, want the staged weights (~15 GB)", decision.Estimate.RequiredRAMMB)
	}
}
