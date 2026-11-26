package main

// app-vram-arbitration_test.go — the preflight handoff, proven without hardware.
//
// The chat <-> image handoff has four steps that can each fail silently: the
// user's confirmation, the other model's exit, the telemetry read after it, and
// the decision that follows. These tests drive the production handoff
// (App.reclaimOther) with a stubbed engine and stubbed telemetry so the order
// and the data flow are asserted, not assumed:
//
//   - no confirmation          -> nothing is unloaded, nothing is measured
//   - unload refuses/incomplete -> no admission is attempted afterwards
//   - telemetry never recovers  -> the handoff times out instead of reporting fit
//   - freed memory appears late -> the admitting decision reads it after the stop
//
// A pass here is a statement about wiring. It is not a statement that a real GPU
// gave the memory back; only the controlled hardware run can say that.

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	backendtypes "howl-chat/internal/backend/types"
	"howl-chat/internal/resource"
)

const testMB = int64(1024 * 1024)

// handoffRig stands in for the device and the engines. Telemetry only rises to
// the post-unload value once the stubbed unload has completed, which is what the
// production sequence is supposed to wait for.
type handoffRig struct {
	mu sync.Mutex

	otherLoaded bool
	preFreeMB   int64
	postFreeMB  int64
	// stuck models a stop that returns success while the allocation is still
	// held: telemetry never reaches postFreeMB.
	stuck     bool
	unloadErr error

	released bool
	reads    int
	unloads  int
	events   []string
	ledger   []resource.ModelState
}

func newRig(preFreeMB, postFreeMB int64, otherLoaded bool, ledger ...resource.ModelState) *handoffRig {
	return &handoffRig{
		otherLoaded: otherLoaded,
		preFreeMB:   preFreeMB,
		postFreeMB:  postFreeMB,
		ledger:      ledger,
	}
}

func (r *handoffRig) snapshot(freeMB int64) resource.HardwareSnapshot {
	return resource.HardwareSnapshot{
		GPUModel:       "stub-8g",
		GPUVendor:      "stub",
		TotalVRAMBytes: 8192 * testMB,
		UsedVRAMBytes:  (8192 - freeMB) * testMB,
		FreeVRAMBytes:  freeMB * testMB,
		TotalRAMBytes:  64 << 30,
		UsedRAMBytes:   16 << 30,
		FreeRAMBytes:   48 << 30,
		AvailableRAM:   48 << 30,
		VRAMLive:       true,
		VRAMSource:     "stub-telemetry",
		GPUAvailable:   true,
	}
}

func (r *handoffRig) hardware() (resource.HardwareSnapshot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reads++
	free := r.preFreeMB
	if r.released {
		free = r.postFreeMB
	}
	r.events = append(r.events, fmt.Sprintf("telemetry:%d", free))
	return r.snapshot(free), nil
}

func (r *handoffRig) models() []resource.ModelState {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]resource.ModelState, len(r.ledger))
	copy(out, r.ledger)
	return out
}

func (r *handoffRig) loaded() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.otherLoaded
}

func (r *handoffRig) unload() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, "unload")
	r.unloads++
	if r.unloadErr != nil {
		return r.unloadErr
	}
	r.otherLoaded = false
	if !r.stuck {
		r.released = true
		r.ledger = nil
	}
	return nil
}

func (r *handoffRig) eventIndex(event string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, e := range r.events {
		if e == event {
			return i
		}
	}
	return -1
}

func (r *handoffRig) lastEvent() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.events) == 0 {
		return ""
	}
	return r.events[len(r.events)-1]
}

func (r *handoffRig) readCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reads
}

func (r *handoffRig) unloadCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.unloads
}

func newHandoffApp(r *handoffRig) *App {
	return &App{admission: resource.NewAuthority(r.hardware, r.models)}
}

// imageWorkload is a measured-size checkpoint at the quantized type the panel
// sends, requiring the GPU.
func imageWorkload() resource.AdmissionRequest {
	return resource.AdmissionRequest{
		Kind:            resource.ResourceImage,
		Name:            "stub-sdxl.safetensors",
		Path:            "stub/stub-sdxl.safetensors",
		Width:           1024,
		Height:          1024,
		Steps:           30,
		BatchCount:      1,
		Family:          "sdxl",
		Mode:            "txt2img",
		CheckpointBytes: 6 << 30,
		WeightType:      "q8_0",
		RequireGPU:      true,
	}
}

// chatWorkload needs more than the few hundred MB left beside a loaded image
// model, and far less than the card once that model is gone.
func chatWorkload() resource.AdmissionRequest {
	return resource.AdmissionRequest{
		Kind: resource.ResourceLLM,
		Name: "stub-chat.gguf",
		Path: "stub/stub-chat.gguf",
		Arch: &backendtypes.ArchitectureSpec{
			Family:        "qwen",
			Layers:        28,
			DModel:        3584,
			FFNSize:       18944,
			NumHeads:      28,
			NumKVHeads:    4,
			VocabSize:     151936,
			FileSizeBytes: 5 << 30,
			BitsPerWeight: 4.5,
			HasMetadata:   true,
		},
		GPULayers:  28,
		RequireGPU: true,
	}
}

func chatEntry() resource.ModelState {
	return resource.ModelState{
		Type: resource.ResourceLLM, Name: "stub-chat.gguf", Loaded: true,
		GPUMemoryMB: 5600, IsGPUOffload: true,
	}
}

// shortenReleaseWait keeps the timeout tests fast without changing production
// timing.
func shortenReleaseWait(t *testing.T) {
	t.Helper()
	previous := vramReleaseTimeout
	vramReleaseTimeout = 150 * time.Millisecond
	t.Cleanup(func() { vramReleaseTimeout = previous })
}

// The preflight the panel calls before generation: a loaded chat model holding
// most of the card must produce a refusal that names it and asks the user to
// confirm the unload. This is the dialog text the confirm step consumes.
func TestLoadedChatRefusesImagePreflight(t *testing.T) {
	rig := newRig(600, 7200, true, chatEntry())
	app := newHandoffApp(rig)

	fit, err := app.fitAdmissionRequest("image", imageWorkload())
	if err != nil {
		t.Fatalf("preflight failed to produce a decision: %v", err)
	}
	if fit.CanFit {
		t.Fatalf("preflight admitted an image that cannot coexist with the loaded chat model: %+v", fit)
	}
	if len(fit.Blocking) != 1 || fit.Blocking[0].Type != "llm" {
		t.Fatalf("preflight did not report the chat model as the blocker: %+v", fit.Blocking)
	}
	if !strings.Contains(fit.Message, chatEntry().Name) || !strings.Contains(fit.Message, "Unload it and continue?") {
		t.Fatalf("refusal does not name the blocker and ask for the unload: %q", fit.Message)
	}
	if fit.FreeVRAMMB != 600 {
		t.Fatalf("preflight reports %d MB free, want the measured 600 MB", fit.FreeVRAMMB)
	}
}

// A declined or absent confirmation must not unload anything, and must not even
// read the device: the UI answer gates the whole handoff.
func TestHandoffRefusesWithoutConfirmation(t *testing.T) {
	rig := newRig(600, 7200, true, chatEntry())
	app := newHandoffApp(rig)

	fit, err := app.reclaimOther("image", imageWorkload(), false, reclaimTarget{
		Label:  "chat",
		Loaded: rig.loaded,
		Unload: rig.unload,
	})
	if err == nil {
		t.Fatal("handoff without confirmation returned no error")
	}
	if !strings.Contains(err.Error(), "was not confirmed") {
		t.Fatalf("error does not name the missing confirmation: %v", err)
	}
	if fit.CanFit {
		t.Fatal("handoff without confirmation reported the workload as fitting")
	}
	if n := rig.unloadCount(); n != 0 {
		t.Fatalf("unload ran %d time(s) without confirmation", n)
	}
	if n := rig.readCount(); n != 0 {
		t.Fatalf("device telemetry was read %d time(s) before consent", n)
	}
}

// The admitting decision must consume telemetry produced after the stop, not the
// pre-unload free value that made the request impossible.
func TestHandoffWaitsForPostUnloadTelemetry(t *testing.T) {
	rig := newRig(600, 7200, true, chatEntry())
	app := newHandoffApp(rig)

	fit, err := app.reclaimOther("image", imageWorkload(), true, reclaimTarget{
		Label:  "chat",
		Loaded: rig.loaded,
		Unload: rig.unload,
	})
	if err != nil {
		t.Fatalf("handoff failed after a completed unload: %v", err)
	}
	if !fit.CanFit {
		t.Fatalf("workload did not fit after the unload: %+v", fit)
	}
	if rig.unloadCount() != 1 {
		t.Fatalf("unload ran %d time(s), want 1", rig.unloadCount())
	}
	if rig.lastEvent() != "telemetry:7200" {
		t.Fatalf("the accepting decision did not use post-unload telemetry; events: %v", rig.events)
	}
	unloadAt := rig.eventIndex("unload")
	freedAt := rig.eventIndex("telemetry:7200")
	if unloadAt < 0 || freedAt < 0 || freedAt < unloadAt {
		t.Fatalf("freed memory was read before the unload completed; events: %v", rig.events)
	}
	if rig.readCount() != 1 {
		t.Fatalf("the handoff should measure the device once, after the unload; saw %d read(s)", rig.readCount())
	}
	if fit.FreeVRAMMB != 7200 {
		t.Fatalf("fit reports %d MB free, want the fresh 7200 MB", fit.FreeVRAMMB)
	}
}

// A stop that succeeds while the memory stays held must time out as a failure.
// This is the case a fire-and-forget unload would report as success.
func TestStaleTelemetryAfterUnloadIsRefused(t *testing.T) {
	shortenReleaseWait(t)
	rig := newRig(600, 7200, true, chatEntry())
	rig.stuck = true
	app := newHandoffApp(rig)

	fit, err := app.reclaimOther("image", imageWorkload(), true, reclaimTarget{
		Label:  "chat",
		Loaded: rig.loaded,
		Unload: rig.unload,
	})
	if err == nil {
		t.Fatal("handoff succeeded although the freed VRAM never appeared in telemetry")
	}
	if !strings.Contains(err.Error(), "did not complete") {
		t.Fatalf("error does not report an incomplete reclamation: %v", err)
	}
	if fit.CanFit {
		t.Fatal("a stale-telemetry handoff reported the workload as fitting")
	}
	if rig.unloadCount() != 1 {
		t.Fatalf("unload ran %d time(s), want 1", rig.unloadCount())
	}
}

// An unload that refuses to complete blocks admission: no decision is taken, and
// no fit is reported.
func TestFailedUnloadBlocksAdmission(t *testing.T) {
	rig := newRig(600, 7200, true, chatEntry())
	rig.unloadErr = errors.New("llama-server pid 1234 did not exit within 5s")
	app := newHandoffApp(rig)

	fit, err := app.reclaimOther("image", imageWorkload(), true, reclaimTarget{
		Label:  "chat",
		Loaded: rig.loaded,
		Unload: rig.unload,
	})
	if err == nil {
		t.Fatal("handoff reported success after the unload failed")
	}
	if !strings.Contains(err.Error(), "could not unload the chat model") {
		t.Fatalf("error does not name the failed unload: %v", err)
	}
	if fit.CanFit {
		t.Fatal("a failed unload reported the workload as fitting")
	}
	if n := rig.readCount(); n != 0 {
		t.Fatalf("admission was attempted %d time(s) after a failed unload", n)
	}
}

// The reverse handoff (image -> chat) must re-check the device after the image
// model is gone; a chat "ready" answer that never looked at fresh telemetry is
// the failure this asserts against.
func TestReverseHandoffRechecksDeviceBeforeReportingChatReady(t *testing.T) {
	imageEntry := resource.ModelState{
		Type: resource.ResourceImage, Name: "stub-sdxl.safetensors", Loaded: true,
		GPUMemoryMB: 4000, IsGPUOffload: true,
	}
	rig := newRig(600, 7200, true, imageEntry)
	app := newHandoffApp(rig)

	fit, err := app.reclaimOther("llm", chatWorkload(), true, reclaimTarget{
		Label:  "image",
		Loaded: rig.loaded,
		Unload: rig.unload,
	})
	if err != nil {
		t.Fatalf("reverse handoff failed: %v", err)
	}
	if !fit.CanFit {
		t.Fatalf("chat did not fit after the image model was unloaded: %+v", fit)
	}
	if len(fit.Blocking) != 0 {
		t.Fatalf("chat fit still reports blockers after the unload: %+v", fit.Blocking)
	}
	unloadAt := rig.eventIndex("unload")
	freedAt := rig.eventIndex("telemetry:7200")
	if unloadAt < 0 || freedAt < 0 || freedAt < unloadAt {
		t.Fatalf("chat readiness was decided without post-unload telemetry; events: %v", rig.events)
	}
	if rig.lastEvent() != "telemetry:7200" {
		t.Fatalf("the chat-ready decision did not use the final telemetry read; events: %v", rig.events)
	}
}

// The reverse handoff must fail the same way when the image model's memory never
// comes back: no chat-ready answer may be reported from stale numbers.
func TestReverseHandoffFailsOnStaleDeviceState(t *testing.T) {
	shortenReleaseWait(t)
	imageEntry := resource.ModelState{
		Type: resource.ResourceImage, Name: "stub-sdxl.safetensors", Loaded: true,
		GPUMemoryMB: 4000, IsGPUOffload: true,
	}
	rig := newRig(600, 7200, true, imageEntry)
	rig.stuck = true
	app := newHandoffApp(rig)

	fit, err := app.reclaimOther("llm", chatWorkload(), true, reclaimTarget{
		Label:  "image",
		Loaded: rig.loaded,
		Unload: rig.unload,
	})
	if err == nil {
		t.Fatal("reverse handoff reported chat ready although the image allocation never came back")
	}
	if fit.CanFit {
		t.Fatal("reverse handoff reported a fit from stale device state")
	}
}
