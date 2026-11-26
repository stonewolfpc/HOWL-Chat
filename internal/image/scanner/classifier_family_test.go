package scanner

import "testing"

func TestGuessFamilyTreatsAEAsFlux(t *testing.T) {
	if got := GuessFamily("ae"); got != FamilyFlux {
		t.Fatalf("ae.safetensors classified as %q, want %q", got, FamilyFlux)
	}
}

// The official ControlNet-v1-1 files are named control_v11p_sd15_openpose.pth,
// without the word "controlnet". Size alone then filed the 1.38 GB model as a
// diffusion checkpoint, so the advertised pose stage had no ControlNet.
func TestClassifyNameRecognisesControlNetV11Files(t *testing.T) {
	for _, name := range []string{
		"control_v11p_sd15_openpose.pth",
		"control_v11f1e_sd15_tile.pth",
		"control_v11p_sd15_canny.pth",
	} {
		if got := ClassifyName(name); got != AssetControlNet {
			t.Errorf("ClassifyName(%q) = %q, want %q", name, got, AssetControlNet)
		}
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

// Converted GGUFs keep their base model's metadata, so the weights name only the
// base: flux for Chroma, sd3 for SD3.5, sd1 for SD2.1. A filename that names the
// member of that base decides which flag set the engine gets. The refinement must
// never cross an architecture, and an uninformative filename must change nothing.
func TestFamilyRefinementStaysInsideTheArchitecture(t *testing.T) {
	cases := []struct {
		name   string
		probed ModelFamily
		named  ModelFamily
		want   ModelFamily
	}{
		{"chroma converted to gguf reads as flux", FamilyFlux, FamilyChroma, FamilyChroma},
		{"sd3.5 converted to gguf reads as sd3", FamilySD3, FamilySD35, FamilySD35},
		{"sd2.1 converted to gguf reads as sd1.5", FamilySD15, FamilySD21, FamilySD21},
		{"an sdxl file is never turned into flux", FamilySDXL, FamilyFlux, FamilySDXL},
		{"an sd1.5 file is never turned into sdxl", FamilySD15, FamilySDXL, FamilySD15},
		{"a flux file named sd3.5 stays flux", FamilyFlux, FamilySD35, FamilyFlux},
		{"an unnamed file keeps the weight verdict", FamilyFlux, FamilyUnknown, FamilyFlux},
		{"a real sd3.5 header stays sd3.5", FamilySD35, FamilySD35, FamilySD35},
		{"a real chroma header stays chroma", FamilyChroma, FamilyFlux, FamilyChroma},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := refineFamilyFromName(tc.probed, tc.named); got != tc.want {
				t.Fatalf("refineFamilyFromName(%q, %q) = %q, want %q", tc.probed, tc.named, got, tc.want)
			}
		})
	}
}

// Krea2 declares the generic qwen_image base tag in its header, and the engine
// still runs it through the Krea2 loader. The name is the only evidence that
// separates them, and the wrong answer hands it the Qwen2.5-VL encoder, which
// sd-server rejects with "cannot load component".
func TestKrea2IsItsOwnFamily(t *testing.T) {
	if got := GuessFamily("Krea-2-Turbo-Q4_K_M"); got != FamilyKrea2 {
		t.Fatalf("Krea-2 name = %q, want krea2", got)
	}
	if got := refineFamilyFromName(FamilyQwenImage, FamilyKrea2); got != FamilyKrea2 {
		t.Fatalf("qwen_image header + Krea2 name = %q, want krea2", got)
	}
	// Any other cross-architecture name still loses to the header.
	if got := refineFamilyFromName(FamilySDXL, FamilyFlux); got != FamilySDXL {
		t.Fatalf("sdxl header + flux name = %q, want sdxl", got)
	}
}
