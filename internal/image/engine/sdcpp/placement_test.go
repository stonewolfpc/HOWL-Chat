package sdcpp

import "testing"

// These are the real lines sd-server printed while loading Z-Image-Turbo Q4_K_M
// on an 8 GB GTX 1080 — the case where every weight ended up in system RAM.
const (
	realPlacementTotal = `[I] total params memory size = 6083.61MB (VRAM 0.00MB, RAM 6083.61MB): text_encoders 1137.00MB(RAM), diffusion_model 4786.61MB(RAM), vae 160.00MB(RAM), controlnet 0.00MB(N/A), extensions 0.00MB(N/A) --- diffusion_engine.cpp:1438`
	realPlacementPlan  = `[I] auto-fit: --backend "diffusion=CUDA0,te=CUDA0,vae=CUDA0" --params-backend "diffusion=cpu,te=cpu,vae=cpu" --- backend_fit.cpp:477`
)

func TestParsePlacementReadsTheRealSummary(t *testing.T) {
	p, ok := parsePlacement(realPlacementTotal)
	if !ok {
		t.Fatal("summary line not recognised")
	}
	if p.VRAMMB != 0 || p.RAMMB != 6083.61 {
		t.Errorf("vram=%v ram=%v, want 0 / 6083.61", p.VRAMMB, p.RAMMB)
	}
	if p.Components["diffusion_model"] != "RAM" || p.Components["text_encoders"] != "RAM" {
		t.Errorf("components = %v", p.Components)
	}
	if p.Components["controlnet"] != "N/A" {
		t.Errorf("unused component should read N/A, got %q", p.Components["controlnet"])
	}
	if !p.Offloaded() {
		t.Error("0 MB VRAM and 6 GB RAM must count as offloaded")
	}
	if s := placementSummary(p); s == "" {
		t.Error("an offloaded model needs a summary sentence")
	}
}

func TestParsePlacementHandlesAPartiallyOffloadedModel(t *testing.T) {
	line := `[I] total params memory size = 6000.00MB (VRAM 4096.00MB, RAM 1904.00MB): text_encoders 1904.00MB(RAM), diffusion_model 4096.00MB(VRAM), vae 0.00MB(N/A)`
	p, ok := parsePlacement(line)
	if !ok {
		t.Fatal("line not recognised")
	}
	if p.Offloaded() {
		t.Error("weights resident in VRAM is not offloaded")
	}
	if p.Components["diffusion_model"] != "VRAM" || p.Components["text_encoders"] != "RAM" {
		t.Errorf("components = %v", p.Components)
	}
}

func TestPlacementIgnoresUnrelatedOutput(t *testing.T) {
	for _, line := range []string{
		`[I] listening on: http://127.0.0.1:8081 --- main.cpp:170`,
		`[I] Version: Z-Image  --- diffusion_engine.cpp:975`,
		``,
	} {
		if _, ok := parsePlacement(line); ok {
			t.Errorf("%q should not be read as a placement", line)
		}
	}
	// An unknown placement must produce no sentence rather than a guess.
	if s := placementSummary(Placement{}); s != "" {
		t.Errorf("unknown placement produced %q", s)
	}
}

func TestPlacementRecordsEngineVAERetry(t *testing.T) {
	var p Placement
	for _, line := range []string{
		`[I] VAE decode compute failed --- vae.cpp:1`,
		`[I] VAE decode ran out of memory; retrying with spatial tiling --- vae.cpp:2`,
		`[I] Reducing VAE decode tiles from 1024x1024 to 512x512 image pixels --- vae.cpp:3`,
	} {
		next, ok := parsePlacementEvent(line)
		if !ok {
			t.Fatalf("event not recognised: %s", line)
		}
		p = p.merge(next)
	}
	if !p.VAETiled || p.VAETilingNote == "" {
		t.Fatalf("engine retry was not recorded: %+v", p)
	}
	if _, ok := parsePlacementEvent(`[I] listening on: http://127.0.0.1:8081`); ok {
		t.Fatal("ordinary output must not be read as a VAE retry")
	}
}

// The plan and the totals arrive as separate lines, in either order, and the
// engine must keep both instead of letting the later one erase the earlier.
func TestPlacementMergeKeepsPlanAndTotals(t *testing.T) {
	var p Placement
	for _, line := range []string{realPlacementTotal, realPlacementPlan} {
		next, ok := parsePlacement(line)
		if !ok {
			t.Fatalf("not recognised: %s", line)
		}
		p = p.merge(next)
	}
	if p.Backends == "" || p.RAMMB == 0 {
		t.Fatalf("merge lost data: %+v", p)
	}
	var q Placement
	for _, line := range []string{realPlacementPlan, realPlacementTotal} {
		next, _ := parsePlacement(line)
		q = q.merge(next)
	}
	if q.Backends != p.Backends || q.RAMMB != p.RAMMB {
		t.Errorf("order changed the result: %+v vs %+v", p, q)
	}
}
