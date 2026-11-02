package scanner

import "testing"

// A distilled model at 20 steps and CFG 7 produces a burned image, and a burned
// image looks the same whether the engine, the offload or the settings are wrong.
// These defaults are what the UI applies when a model is selected.
func TestDefaultGenerationSettingsMatchTheFamily(t *testing.T) {
	cases := []struct {
		family   ModelFamily
		steps    int
		cfg      float64
	}{
		{FamilyZImage, 8, 1.0},
		{FamilyFlux, 4, 1.0},
		{FamilyAnima, 30, 6.0},
		{FamilyQwenImage, 20, 2.5},
		{FamilyFlux2, 20, 4.0},
		{FamilySDXL, 20, 7.0},
		{FamilySD15, 20, 7.0},
	}
	for _, c := range cases {
		if got := DefaultSteps(c.family); got != c.steps {
			t.Errorf("DefaultSteps(%s) = %d, want %d", c.family, got, c.steps)
		}
		if got := DefaultCFG(c.family); got != c.cfg {
			t.Errorf("DefaultCFG(%s) = %v, want %v", c.family, got, c.cfg)
		}
		rec := buildRecommended(AssetDiffusion, c.family)
		if rec.Steps != c.steps || rec.CFG != c.cfg {
			t.Errorf("buildRecommended(%s) = %d steps / CFG %v, want %d / %v",
				c.family, rec.Steps, rec.CFG, c.steps, c.cfg)
		}
	}
	// Companions carry no generation settings at all.
	if got := buildRecommended(AssetVAE, FamilyFlux); got != (Recommended{}) {
		t.Errorf("a VAE got generation defaults: %+v", got)
	}
}
