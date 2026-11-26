package resource

import "testing"

// A workflow upscale loads an adapter that the native checkpoint admission did
// not include. Without this charge a stage could start against free VRAM that
// the upscaler immediately consumes.
func TestImageAdmissionChargesUpscalerModule(t *testing.T) {
	authority := NewAuthority(fixedHardware(testSnapshot(5000, 8192)), fixedLedger())
	base := AdmissionRequest{
		Kind: ResourceImage, Name: "workflow", Width: 512, Height: 512,
		Mode: "dgl_upscale", RequireGPU: true,
	}
	without, lease := authority.Admit(base)
	if !without.Admitted {
		t.Fatalf("baseline upscale admission refused: %s", without.Reason)
	}
	lease.Release()

	with := base
	with.UpscalerBytes = 64 * bytesPerMB
	decision, lease := authority.Admit(with)
	if !decision.Admitted {
		t.Fatalf("64 MB upscaler admission refused at 5 GB free: %s", decision.Reason)
	}
	defer lease.Release()
	if got, want := decision.Estimate.RequiredVRAMMB, sdModuleReserveMB+64; got != want {
		t.Fatalf("upscaler requirement = %d MB, want %d MB (module reserve + measured adapter)", got, want)
	}
	if decision.Estimate.RequiredVRAMMB <= without.Estimate.RequiredVRAMMB {
		t.Fatalf("upscaler did not increase the requirement: %d <= %d",
			decision.Estimate.RequiredVRAMMB, without.Estimate.RequiredVRAMMB)
	}
}
