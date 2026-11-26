package audio

import (
	"context"
	"strings"
	"testing"

	"howl-chat/internal/audio/types"
	"howl-chat/internal/resource"
)

const asrAdmissionTestMB = int64(1024 * 1024)

// asrAdmissionFixture returns the authority both ASR cases consult: a 24 GB
// device with freeMB free and a 3000 MB chat model already resident.
func asrAdmissionFixture(freeMB int64) *resource.Authority {
	loaded := []resource.ModelState{{
		Type: resource.ResourceLLM, Name: "chat", Loaded: true, GPUMemoryMB: 3000, SystemRAMMB: 4000,
	}}
	return resource.NewAuthority(
		func() (resource.HardwareSnapshot, error) {
			return resource.HardwareSnapshot{
				GPUModel:       "Test GPU",
				TotalVRAMBytes: 24000 * asrAdmissionTestMB,
				FreeVRAMBytes:  freeMB * asrAdmissionTestMB,
				UsedVRAMBytes:  (24000 - freeMB) * asrAdmissionTestMB,
				TotalRAMBytes:  32 * 1024 * asrAdmissionTestMB,
				FreeRAMBytes:   16 * 1024 * asrAdmissionTestMB,
				VRAMLive:       true,
				VRAMSource:     "test-live",
				GPUAvailable:   true,
			}, nil
		},
		func() []resource.ModelState { return loaded },
	)
}

// asrQwenRequest pins the profile the accounting below depends on.
func asrQwenRequest(t *testing.T) resource.AdmissionRequest {
	t.Helper()
	req, err := asrAdmissionRequest(types.ASRQwen3ASR17BLocal)
	if err != nil {
		t.Fatalf("asr admission request: %v", err)
	}
	if req.AudioVRAMMB != 4096 {
		t.Fatalf("ASR profile VRAM = %d MB, want the 4096 MB this accounting pins", req.AudioVRAMMB)
	}
	return req
}

// TestASRStartConsultsAdmissionWhileChatLoaded proves the pipeline does not
// spawn the local Qwen ASR child before the shared authority admits it. The
// fixture is genuinely refused: 4000 MB free keeps a 400 MB reserve, the 272 MB
// soft allowance leaves 3872 MB usable, and the 4096 MB request does not fit.
func TestASRStartConsultsAdmissionWhileChatLoaded(t *testing.T) {
	authority := asrAdmissionFixture(4000)
	req := asrQwenRequest(t)

	decision, lease := authority.Admit(req)
	if decision.Admitted || lease != nil {
		t.Fatalf("admission at 4000 MB free = admitted, want refusal; reason: %s", decision.Reason)
	}
	est := decision.Estimate
	if est.FreeVRAMMB != 4000 || est.ReserveVRAMMB != 400 || est.SoftReserveAllowanceMB != 272 || est.AvailableVRAMMB != 3600 {
		t.Fatalf("estimate free/reserve/allowance/available = %d/%d/%d/%d, want 4000/400/272/3600",
			est.FreeVRAMMB, est.ReserveVRAMMB, est.SoftReserveAllowanceMB, est.AvailableVRAMMB)
	}
	if !strings.Contains(decision.Reason, `blocking model: llm "chat"`) {
		t.Fatalf("refusal must name the chat model: %s", decision.Reason)
	}

	p := NewPipeline()
	p.admission = authority
	_, outcome := p.initASR(context.Background(), []types.ASRModelType{types.ASRQwen3ASR17BLocal}, types.PipelineConfig{})
	if outcome.ready {
		t.Fatal("ASR must not become ready when admission refuses the start")
	}
	if outcome.err == nil || !strings.Contains(outcome.err.Error(), `blocking model: llm "chat"`) {
		t.Fatalf("ASR refusal must name the chat model: %v", outcome.err)
	}
}

// TestASRAdmissionAdmitsThroughSoftReserveWhileChatLoaded pins the other side of
// the same boundary. The documented policy treats the 10% reserve as a target:
// at 4500 MB free the 4096 MB request is admitted using 46 MB of the 322 MB soft
// allowance, leaving 404 MB free above the 128 MB hard floor. The test stops at
// the admission boundary; it does not spawn a recognizer.
func TestASRAdmissionAdmitsThroughSoftReserveWhileChatLoaded(t *testing.T) {
	authority := asrAdmissionFixture(4500)
	req := asrQwenRequest(t)

	decision, lease := authority.Admit(req)
	if !decision.Admitted || lease == nil {
		t.Fatalf("admission at 4500 MB free = refused, want the documented soft-reserve admission; reason: %s", decision.Reason)
	}
	lease.Release()

	est := decision.Estimate
	if est.FreeVRAMMB != 4500 || est.ReserveVRAMMB != 450 || est.AvailableVRAMMB != 4050 {
		t.Fatalf("estimate free/reserve/available = %d/%d/%d, want 4500/450/4050",
			est.FreeVRAMMB, est.ReserveVRAMMB, est.AvailableVRAMMB)
	}
	if est.SoftReserveAllowanceMB != 322 || est.SoftReserveUsedMB != 46 {
		t.Fatalf("soft allowance used = %d of %d MB, want 46 of 322", est.SoftReserveUsedMB, est.SoftReserveAllowanceMB)
	}
	if headroom := est.FreeVRAMMB - est.RequiredVRAMMB; headroom != 404 {
		t.Fatalf("free VRAM above the request = %d MB, want 404 MB", headroom)
	}
	if !strings.Contains(decision.Reason, "soft reserve") {
		t.Fatalf("admission reason must report the soft reserve: %s", decision.Reason)
	}

	p := NewPipeline()
	p.admission = authority
	pipelineLease, err := p.admitAudio(req)
	if err != nil {
		t.Fatalf("pipeline admission must accept the same request: %v", err)
	}
	pipelineLease.Release()
}
