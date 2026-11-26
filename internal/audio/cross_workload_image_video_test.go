package audio

// Cross-workload resource and lifecycle coverage: audio (TTS residency and
// execution) against Howl's image and video workloads.
//
// Every decision below goes through the one production authority
// (resource.Authority) and the one process registry. Nothing here replaces the
// admission policy, the reserve formula or the ownership model with a test-only
// variant: the audio residency is established by the pipeline's own
// registerAudio path, the image/video workload is expressed as the same
// resource.AdmissionRequest the app's image and video start paths build, and
// the process assertions use real helper children.
//
// Image and video share one admission path in production:
//   - app-image-model.go: ensureEngineReadyLocked -> admitImageLoad
//   - app-image-model.go: ensureVideoExecutionLocked -> admitImageLoad, with the
//     lease held until sd-cli exits
// so a video request is an image request whose Mode is "video"; it is modelled
// that way here rather than as a separate memory reservation.

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"howl-chat/internal/audio/tts/cache"
	"howl-chat/internal/audio/types"
	"howl-chat/internal/process"
	"howl-chat/internal/resource"
)

const crossWorkloadMB = int64(1024 * 1024)

// mutableHardware is physical telemetry the test can refresh. It exists so a
// test can distinguish "the owner released its claim" from "the card actually
// reported the memory back", which admission must not conflate.
type mutableHardware struct {
	mu          sync.Mutex
	totalVRAMMB int64
	freeVRAMMB  int64
	freeRAMMB   int64
}

func (h *mutableHardware) snapshot() (resource.HardwareSnapshot, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return resource.HardwareSnapshot{
		GPUModel:       "Cross-workload test GPU",
		TotalVRAMBytes: h.totalVRAMMB * crossWorkloadMB,
		FreeVRAMBytes:  h.freeVRAMMB * crossWorkloadMB,
		UsedVRAMBytes:  (h.totalVRAMMB - h.freeVRAMMB) * crossWorkloadMB,
		TotalRAMBytes:  64 * 1024 * crossWorkloadMB,
		FreeRAMBytes:   h.freeRAMMB * crossWorkloadMB,
		VRAMLive:       true,
		VRAMSource:     "test-live",
		GPUAvailable:   true,
	}, nil
}

func (h *mutableHardware) setFreeVRAMMB(mb int64) {
	h.mu.Lock()
	h.freeVRAMMB = mb
	h.mu.Unlock()
}

func newCrossWorkloadAuthority(freeVRAMMB, totalVRAMMB int64) (*resource.Authority, *mutableHardware) {
	hw := &mutableHardware{totalVRAMMB: totalVRAMMB, freeVRAMMB: freeVRAMMB, freeRAMMB: 16000}
	authority := resource.NewAuthority(
		hw.snapshot,
		func() []resource.ModelState { return resource.GetManager().LoadedModels() },
	)
	return authority, hw
}

// crossWorkloadImageRequest is the shape the image start path builds: a
// measured checkpoint, conditioner and VAE, plus the workload's decode size.
func crossWorkloadImageRequest(mode string, width, height int) resource.AdmissionRequest {
	return resource.AdmissionRequest{
		Kind:             resource.ResourceImage,
		Name:             "sdxl-cross",
		Path:             `C:\models\cross\sdxl.safetensors`,
		Width:            width,
		Height:           height,
		BatchCount:       1,
		Steps:            20,
		Family:           "sdxl",
		Mode:             mode,
		CheckpointBytes:  6600 * crossWorkloadMB,
		ConditionerBytes: 1500 * crossWorkloadMB,
		VAEBytes:         300 * crossWorkloadMB,
	}
}

// crossWorkloadVideoRequest is the same workload as a video render: one that
// sd-cli performs, still admitted as an image-family start.
func crossWorkloadVideoRequest() resource.AdmissionRequest {
	req := crossWorkloadImageRequest("video", 832, 480)
	req.Name = "wan-cross"
	req.Family = "wan"
	return req
}

// installResidentAudioProvider is what initTTS does once a provider is loaded:
// admit against the one authority, convert the lease into a residency
// registration, and install the provider as the active one. The request is the
// provider profile's own, not a number invented for the test.
func installResidentAudioProvider(t *testing.T, p *Pipeline, engine string, synth types.TTSSynthesizer) *ttsProvider {
	t.Helper()
	req, err := ttsAdmissionRequest(engine, "")
	if err != nil {
		t.Fatalf("admission request for %s: %v", engine, err)
	}
	key := ttsEngineKey(engine, "")
	lease, err := p.admitAudio(req)
	if err != nil {
		t.Fatalf("admit %s: %v", engine, err)
	}
	provider := &ttsProvider{key: key, synth: synth, engine: engine, modelName: key}
	provider.registered = p.registerAudio(resource.ResourceTTS, key, "", req, lease)
	t.Cleanup(func() {
		// A failing assertion must not leak a residency into another test.
		resource.GetManager().UnregisterModel(resource.ResourceTTS, key)
		p.mu.Lock()
		delete(p.ttsProviders, key)
		p.mu.Unlock()
	})
	p.mu.Lock()
	p.registerTTSResidencyLocked(provider)
	p.ttsSynth = synth
	p.activeTTSEngine = engine
	p.activeTTSKey = key
	p.globalTTSKey = key
	p.mu.Unlock()
	return provider
}

// countingNames reports how many ledger entries carry name. Audio residency must
// appear exactly once: a provider counted as both a resident registration and as
// an execution reservation would report itself twice and spend capacity twice.
func countingNames(blocking []resource.ModelState, name string) int {
	count := 0
	for _, model := range blocking {
		if model.Name == name {
			count++
		}
	}
	return count
}

// TestAudioResidencyBlocksImageAdmissionThenReleases is workstream item 2: an
// audio residency appears in the image decision, the release is the pipeline's
// own, physical telemetry is refreshed before the retry, and an admitted lease
// is the pending commitment that stops a second start from spending the same
// free value.
func TestAudioResidencyBlocksImageAdmissionThenReleases(t *testing.T) {
	authority, hw := newCrossWorkloadAuthority(8500, 24000)

	p := NewPipeline(process.NewRegistry())
	p.admission = authority
	p.ttsCache = cache.NewManager(10)
	provider := installResidentAudioProvider(t, p, "melotts", &stubTTSSynth{})

	req := crossWorkloadImageRequest("txt2img", 512, 512)
	decision, lease := authority.Admit(req)
	if decision.Admitted {
		lease.Release()
		t.Fatalf("image admitted on 8500 MB free beside a resident audio provider: %s", decision.Reason)
	}
	// The refusal is actionable: it carries the numbers, not just a verdict.
	if decision.Estimate.ProjectedPeakVRAMMB == 0 || decision.Estimate.AvailableVRAMMB == 0 {
		t.Fatalf("refusal lacks its estimate: %+v", decision.Estimate)
	}
	if decision.Estimate.ProjectedPeakVRAMMB <= decision.Estimate.AvailableVRAMMB {
		t.Fatalf("refusal numbers do not show the shortfall: projected=%d available=%d",
			decision.Estimate.ProjectedPeakVRAMMB, decision.Estimate.AvailableVRAMMB)
	}
	if !blockingNames(decision.Blocking, provider.key) {
		t.Fatalf("audio residency is missing from the image refusal: %s", decision.Reason)
	}
	if got := countingNames(decision.Blocking, provider.key); got != 1 {
		t.Fatalf("audio residency appears %d times in the ledger, want 1", got)
	}
	if !strings.Contains(decision.Reason, "cannot start image") {
		t.Fatalf("refusal does not identify the workload: %s", decision.Reason)
	}

	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	if got := resource.GetManager().ModelsByType(resource.ResourceTTS); len(got) != 0 {
		t.Fatalf("audio residency survived shutdown: %+v", got)
	}

	// The owner is gone, but the card has not reported the memory back. A
	// process exit is not evidence of recovered capacity.
	decision, lease = authority.Admit(req)
	if decision.Admitted {
		lease.Release()
		t.Fatal("image admission trusted a released owner while physical telemetry still reported the old free value")
	}

	hw.setFreeVRAMMB(12500)
	decision, lease = authority.Admit(req)
	if !decision.Admitted || lease == nil {
		t.Fatalf("image refused after the audio residency was released and telemetry refreshed: %s", decision.Reason)
	}
	if blockingNames(decision.Blocking, provider.key) {
		t.Fatalf("released audio residency still appears in the ledger: %s", decision.Reason)
	}

	// The admitted start has not registered its model yet, so its lease is the
	// only record of the commitment. A second start must not pass against it.
	second, secondLease := authority.Admit(req)
	if second.Admitted {
		secondLease.Release()
		t.Fatal("a second image start spent the same free value the first admission lease had committed")
	}
	lease.Release()
	lease.Release() // idempotent
	third, thirdLease := authority.Admit(req)
	if !third.Admitted {
		t.Fatalf("image refused after the admission lease was released: %s", third.Reason)
	}
	thirdLease.Release()
}

// TestAudioExecutionLeaseSurvivesImageWorkAndReclaimsOnlyWhenIdle is workstream
// item 3 with a shared provider: two legitimate owners hold execution leases,
// neither release unloads the provider, the image workload sees one residency
// entry rather than two, and reclamation is refused until the last lease ends.
func TestAudioExecutionLeaseSurvivesImageWorkAndReclaimsOnlyWhenIdle(t *testing.T) {
	authority, hw := newCrossWorkloadAuthority(8500, 24000)

	p := NewPipeline(process.NewRegistry())
	p.admission = authority
	provider := installResidentAudioProvider(t, p, "melotts", &stubTTSSynth{})

	first := p.acquireTTSExecution(provider)
	second := p.acquireTTSExecution(provider)

	if got := resource.GetManager().ModelsByType(resource.ResourceTTS); len(got) != 1 {
		t.Fatalf("shared provider has %d registrations, want 1", len(got))
	}

	req := crossWorkloadImageRequest("txt2img", 512, 512)
	decision, lease := authority.Admit(req)
	if decision.Admitted {
		lease.Release()
		t.Fatal("image admitted while an audio provider was executing and the card was full")
	}
	if got := countingNames(decision.Blocking, provider.key); got != 1 {
		t.Fatalf("an executing shared provider appears %d times, want 1 (double counted): %s", got, decision.Reason)
	}

	if err := p.unloadTTSProvider(provider); err == nil || !strings.Contains(err.Error(), "executing") {
		t.Fatalf("reclamation of an executing provider was not refused: %v", err)
	}
	first()
	if err := p.unloadTTSProvider(provider); err == nil {
		t.Fatal("one of two execution leases released the still-required residency")
	}
	second()
	if err := p.unloadTTSProvider(provider); err != nil {
		t.Fatalf("idle provider reclamation: %v", err)
	}
	if got := resource.GetManager().ModelsByType(resource.ResourceTTS); len(got) != 0 {
		t.Fatalf("unloaded provider kept its registration: %+v", got)
	}

	hw.setFreeVRAMMB(12500)
	decision, lease = authority.Admit(req)
	if !decision.Admitted {
		t.Fatalf("image refused after the audio residency and its execution leases ended: %s", decision.Reason)
	}
	lease.Release()
}

// TestImageResidencyBlocksAudioAdmission is the opposite direction: a resident
// image model and its execution lease are visible to the audio workload, and the
// release order matters. A test that only proves audio-blocks-image would miss a
// decision that ignores one side of the ledger.
func TestImageResidencyBlocksAudioAdmission(t *testing.T) {
	authority, hw := newCrossWorkloadAuthority(3500, 24000)

	const imageName = "sdxl-resident-cross"
	resource.GetManager().RegisterModel(resource.ResourceImage, imageName, `C:\models\cross\sdxl.safetensors`, 9000, 2000, resource.PriorityImage, true)
	t.Cleanup(func() { resource.GetManager().UnregisterModel(resource.ResourceImage, imageName) })

	p := NewPipeline(process.NewRegistry())
	p.admission = authority

	req, err := asrAdmissionRequest(types.ASRQwen3ASR17BLocal)
	if err != nil {
		t.Fatalf("asr request: %v", err)
	}
	lease, err := p.admitAudio(req)
	if err == nil {
		lease.Release()
		t.Fatal("audio was admitted beside a resident image model on 3500 MB free VRAM")
	}
	if !strings.Contains(err.Error(), imageName) {
		t.Fatalf("audio refusal does not name the image residency: %v", err)
	}

	// The image owner releases its registration, and the card reports the
	// memory back. Both are required: the ledger entry only describes a resident
	// model, it never invents free VRAM.
	resource.GetManager().UnregisterModel(resource.ResourceImage, imageName)
	hw.setFreeVRAMMB(9000)
	lease, err = p.admitAudio(req)
	if err != nil {
		t.Fatalf("audio refused after the image residency was released and telemetry refreshed: %v", err)
	}
	lease.Release()
}

// TestVideoRenderReservationObservesAudioResidency is workstream item 4 at the
// admission boundary the video start path uses. Video is admitted as an
// image-family workload whose Mode is "video"; the reservation is held for the
// whole render (that is what ensureVideoExecution returns to the caller after
// sd-cli exits), and audio ownership is unaffected by the video start.
func TestVideoRenderReservationObservesAudioResidency(t *testing.T) {
	authority, hw := newCrossWorkloadAuthority(8500, 24000)

	p := NewPipeline(process.NewRegistry())
	p.admission = authority
	p.ttsCache = cache.NewManager(10)
	provider := installResidentAudioProvider(t, p, "melotts", &stubTTSSynth{})

	video := crossWorkloadVideoRequest()
	decision, lease := authority.Admit(video)
	if decision.Admitted {
		lease.Release()
		t.Fatalf("video render admitted while audio was resident on a full card: %s", decision.Reason)
	}
	if !blockingNames(decision.Blocking, provider.key) {
		t.Fatalf("video refusal does not name the audio residency: %s", decision.Reason)
	}
	if decision.Estimate.ActivationVRAMMB <= 0 {
		t.Fatalf("video request has no decode/render working set: %+v", decision.Estimate)
	}

	// Releasing audio does not hand the video start capacity the card has not
	// reported; only the refreshed measurement admits it.
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	decision, lease = authority.Admit(video)
	if decision.Admitted {
		lease.Release()
		t.Fatal("video admission passed on stale telemetry after the audio owner was released")
	}
	hw.setFreeVRAMMB(12500)

	decision, lease = authority.Admit(video)
	if !decision.Admitted || lease == nil {
		t.Fatalf("video refused after audio was released and telemetry refreshed: %s", decision.Reason)
	}

	// The render reservation is held until sd-cli exits, so a second start
	// (image or video) cannot spend the same free value meanwhile.
	other, otherLease := authority.Admit(crossWorkloadImageRequest("txt2img", 512, 512))
	if other.Admitted {
		otherLease.Release()
		t.Fatal("a concurrent image start spent capacity reserved for the active video render")
	}
	lease.Release()
	lease.Release() // ensureVideoExecution releases exactly once; the release is idempotent

	after, afterLease := authority.Admit(crossWorkloadImageRequest("txt2img", 512, 512))
	if !after.Admitted {
		t.Fatalf("image refused after the video reservation was released: %s", after.Reason)
	}
	afterLease.Release()
}

// TestVideoProcessStopLeavesAudioOwnerAndRegistrationIntact is workstream item 4
// with real processes: an owned TTS child and an owned sd-cli share the one
// registry. Stopping the video child must not stop the audio child, must not
// remove the audio registration, and must leave no sd-cli registration behind.
// Audio shutdown then reaps its own child and clears exactly its own entry.
func TestVideoProcessStopLeavesAudioOwnerAndRegistrationIntact(t *testing.T) {
	registry := process.NewRegistry()
	p := NewPipeline(registry)
	p.admission = permissiveAuthority()
	provider := installResidentAudioProvider(t, p, "piper", &stubTTSSynth{})

	audioChild, err := process.StartChild(registry, process.ChildSpec{
		Role:     process.RoleTTSChild,
		Instance: "cross-workload-piper",
		Path:     os.Args[0],
		Args:     []string{"-test.run=TestAudioChildHelper"},
		Env:      append(os.Environ(), "HOWL_AUDIO_CHILD=1"),
		Capture:  true,
	})
	if err != nil {
		t.Fatalf("start owned TTS child: %v", err)
	}
	t.Cleanup(func() {
		if !audioChild.Exited() {
			_ = process.StopRole(registry, process.RoleTTSChild, process.DefaultGrace)
		}
	})

	videoCmd := exec.Command(os.Args[0], "-test.run=TestAudioVideoHelper")
	videoCmd.Env = append(os.Environ(), "HOWL_VIDEO_HELPER=1")
	if err := videoCmd.Start(); err != nil {
		t.Fatalf("start sd-cli helper: %v", err)
	}
	videoDone := make(chan struct{})
	var waitErr error
	go func() { waitErr = videoCmd.Wait(); close(videoDone) }()
	defer func() { _ = waitErr }()
	if _, err := registry.Register(process.RoleVideoCLI, videoCmd, func(timeout time.Duration) error {
		select {
		case <-videoDone:
			return nil
		case <-time.After(timeout):
			return fmt.Errorf("sd-cli helper pid %d still running after %s", videoCmd.Process.Pid, timeout)
		}
	}); err != nil {
		t.Fatalf("register sd-cli: %v", err)
	}
	t.Cleanup(func() {
		if registry.Current(process.RoleVideoCLI) != nil {
			_ = registry.Stop(process.RoleVideoCLI, process.DefaultGrace)
		}
		select {
		case <-videoDone:
		case <-time.After(process.DefaultGrace):
		}
	})

	if got := registry.Current(process.RoleVideoCLI); got == nil {
		t.Fatal("sd-cli was not registered while it was running")
	}
	if got := registry.Children(process.RoleTTSChild); len(got) != 1 {
		t.Fatalf("TTS child registrations = %d, want 1", len(got))
	}

	// The video workload finishes or is cancelled: its own slot must close.
	if err := registry.Stop(process.RoleVideoCLI, process.DefaultGrace); err != nil {
		t.Fatalf("stop sd-cli: %v", err)
	}
	if got := registry.Current(process.RoleVideoCLI); got != nil {
		t.Fatalf("sd-cli registration survived its exit: %s", got.Identity)
	}
	select {
	case <-videoDone:
	case <-time.After(process.DefaultGrace):
		t.Fatalf("sd-cli exit was not observed by its waiter: %v", waitErr)
	}

	if audioChild.Exited() {
		t.Fatal("stopping the video child also stopped the audio synthesis child")
	}
	if got := registry.Children(process.RoleTTSChild); len(got) != 1 {
		t.Fatalf("video stop removed the audio child's registration: %+v", got)
	}
	if got := resource.GetManager().ModelsByType(resource.ResourceTTS); len(got) != 1 {
		t.Fatalf("video stop removed the audio residency: %+v", got)
	}
	if p.providerByKey(provider.key) != provider {
		t.Fatal("video stop removed the resident audio provider")
	}

	// Audio shutdown owns exactly its own resources.
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("audio shutdown: %v", err)
	}
	if !audioChild.Exited() {
		t.Fatal("audio shutdown returned before its owned child exited")
	}
	if got := registry.Children(process.RoleTTSChild); len(got) != 0 {
		t.Fatalf("audio shutdown left its child registered: %+v", got)
	}
	if got := resource.GetManager().ModelsByType(resource.ResourceTTS); len(got) != 0 {
		t.Fatalf("audio shutdown left its registration: %+v", got)
	}
	if got := registry.Current(process.RoleVideoCLI); got != nil {
		t.Fatalf("audio shutdown resurrected or retained the video slot: %s", got.Identity)
	}
}

// TestAudioCancellationDoesNotReleaseImageExecutionLease is workstream item 5:
// cancelling an audio operation reaches that operation's child and cleanup, and
// must not disturb another workload's execution lease. The image lease is still
// the pending commitment after the cancellation, and releasing it once (twice,
// for idempotence) is what returns the capacity.
func TestAudioCancellationDoesNotReleaseImageExecutionLease(t *testing.T) {
	authority, _ := newCrossWorkloadAuthority(12500, 24000)

	imageLeaseDecision, imageLease := authority.Admit(crossWorkloadImageRequest("txt2img", 512, 512))
	if !imageLeaseDecision.Admitted || imageLease == nil {
		t.Fatalf("image admission: %s", imageLeaseDecision.Reason)
	}
	released := false
	t.Cleanup(func() {
		if !released {
			imageLease.Release()
		}
	})

	p := NewPipeline(process.NewRegistry())
	p.admission = authority
	enumerator := &blockingEnumerator{entered: make(chan struct{}), finished: make(chan struct{})}
	key := ttsEngineKey("kokoro", "")
	p.mu.Lock()
	p.ttsSynth = enumerator
	p.activeTTSEngine = "kokoro"
	p.activeTTSKey = key
	p.ttsProviders[key] = &ttsProvider{key: key, synth: enumerator, engine: "kokoro"}
	p.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := p.GetTTSVoicesContext(ctx)
		result <- err
	}()
	select {
	case <-enumerator.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("audio enumeration did not start")
	}
	cancel()
	select {
	case err := <-result:
		if err == nil {
			t.Fatal("cancelled audio enumeration reported success")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cancelled audio enumeration did not return")
	}
	select {
	case <-enumerator.finished:
	default:
		t.Fatal("audio enumeration returned before its cleanup finished")
	}
	p.mu.RLock()
	refs := p.ttsProviders[key].refs
	p.mu.RUnlock()
	if refs != 0 {
		t.Fatalf("cancelled audio operation leaked %d execution lease(s)", refs)
	}

	if other, otherLease := authority.Admit(crossWorkloadImageRequest("txt2img", 512, 512)); other.Admitted {
		otherLease.Release()
		t.Fatal("audio cancellation released the image workload's execution lease")
	}

	imageLease.Release()
	released = true
	imageLease.Release()

	if after, afterLease := authority.Admit(crossWorkloadImageRequest("txt2img", 512, 512)); !after.Admitted {
		t.Fatalf("image refused after its own lease was released: %s", after.Reason)
	} else {
		afterLease.Release()
	}
}

// TestAudioVideoHelper is re-executed as the sd-cli stand-in above. It holds the
// process open long enough for the parent to observe its registration and stop.
func TestAudioVideoHelper(t *testing.T) {
	if os.Getenv("HOWL_VIDEO_HELPER") != "1" {
		t.Skip("helper process only")
	}
	time.Sleep(30 * time.Second)
}
