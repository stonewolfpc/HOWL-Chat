package scanner

import "testing"

func TestGuessFamilyTreatsAEAsFlux(t *testing.T) {
	if got := GuessFamily("ae"); got != FamilyFlux {
		t.Fatalf("ae.safetensors classified as %q, want %q", got, FamilyFlux)
	}
}

func TestUpgradeLargeCheckpoint(t *testing.T) {
	cases := []struct {
		family ModelFamily
		sizeMB float64
		want   ModelFamily
	}{
		{FamilySD15, 6616, FamilySDXL}, // unnamed 6.6 GB merge
		{FamilySD15, 2100, FamilySD15}, // genuine SD1.5 checkpoint
		{FamilySDXL, 6616, FamilySDXL},
		{FamilyWan, 9165, FamilyWan}, // video families are never promoted
	}
	for _, c := range cases {
		if got := upgradeLargeCheckpoint(c.family, c.sizeMB); got != c.want {
			t.Errorf("upgradeLargeCheckpoint(%q, %.0f) = %q, want %q", c.family, c.sizeMB, got, c.want)
		}
	}
}

func TestFluxVAENotCompatibleWithSDXL(t *testing.T) {
	if CompatibleVAE(FamilySDXL, FamilyFlux) {
		t.Fatal("flux VAE must not be compatible with an SDXL checkpoint")
	}
	if !CompatibleVAE(FamilyFlux, FamilyFlux) {
		t.Fatal("flux VAE must stay compatible with flux models")
	}
}
