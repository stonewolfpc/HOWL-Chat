package deps

import "testing"

// Identity rules that were wrong in the field. Z-Image was handed a Qwen3-0.6B
// encoder — 1024-dim where its DiT expects 2560 — and Anima was handed the Flux
// autoencoder. Both got far enough to report a loaded model, which is exactly the
// failure this check exists to prevent: loading is not the same as working.

func TestCompatibleRejectsTheWrongEncoderSize(t *testing.T) {
	if ok, why := Compatible("z_image", "llm", "qwen_3_06b_base.safetensors"); ok {
		t.Errorf("0.6B encoder accepted for z_image (%s)", why)
	}
	if ok, _ := Compatible("z_image", "llm", "qwen_3_4b.safetensors"); !ok {
		t.Error("4B encoder rejected for z_image")
	}
	// The file stable-diffusion.cpp's own Z-Image docs tell you to download.
	if ok, _ := Compatible("z_image", "llm", "Qwen3-4B-Instruct-2507-Q4_K_M.gguf"); !ok {
		t.Error("sd.cpp's recommended GGUF encoder rejected for z_image")
	}
	if ok, _ := Compatible("anima", "llm", "qwen_3_06b_base.safetensors"); !ok {
		t.Error("0.6B encoder rejected for anima")
	}
	if ok, why := Compatible("anima", "llm", "qwen_3_4b.safetensors"); ok {
		t.Errorf("4B encoder accepted for anima, whose DiT is driven by Qwen3-0.6B (%s)", why)
	}
	// A chat GGUF from the local library must never pass as an image encoder.
	if ok, _ := Compatible("z_image", "llm", "llama-q2_k.gguf"); ok {
		t.Error("llama GGUF accepted as a text encoder")
	}
}

func TestCompatibleSeparatesTheAutoencoders(t *testing.T) {
	cases := []struct {
		family, role, file string
		want               bool
	}{
		{"anima", "vae", "qwen_image_vae.safetensors", true},
		{"anima", "vae", "ae.safetensors", false},
		{"qwen_image", "vae", "qwen_image_vae.safetensors", true},
		{"qwen_image", "vae", "ae.safetensors", false},
		{"flux", "vae", "ae.safetensors", true},
		{"z_image", "vae", "ae.safetensors", true}, // Z-Image's VAE is the Flux one
	}
	for _, c := range cases {
		if ok, why := Compatible(c.family, c.role, c.file); ok != c.want {
			t.Errorf("Compatible(%q,%q,%q) = %v, want %v (%s)", c.family, c.role, c.file, ok, c.want, why)
		}
	}
}

// Every family that needs an encoder or an autoencoder must name its own file.
// The old catalogue fell back to another family's entry, which is how Flux's VAE
// reached Anima; an unknown family must now get no guidance rather than a guess.
func TestCatalogueNamesEveryCompanionExplicitly(t *testing.T) {
	for _, family := range []string{"z_image", "anima", "qwen_image", "flux2"} {
		art := artifactFor(family, "llm")
		if art.FileName == "" {
			t.Errorf("%s has no llm entry of its own", family)
			continue
		}
		if len(art.Match) == 0 {
			t.Errorf("%s llm entry cannot judge identity", family)
		}
		if art.URL == "" {
			t.Errorf("%s llm entry has no download URL", family)
		}
	}
	for _, family := range []string{"flux", "z_image", "anima", "qwen_image", "flux2"} {
		art := artifactFor(family, "vae")
		if art.FileName == "" || len(art.Match) == 0 {
			t.Errorf("%s has no usable vae entry of its own", family)
		}
	}
	if got := artifactFor("made_up_family", "llm"); got.FileName != "" {
		t.Errorf("unknown family inherited %q as an encoder", got.FileName)
	}
	if got := artifactFor("made_up_family", "vae"); got.FileName != "" {
		t.Errorf("unknown family inherited %q as a VAE", got.FileName)
	}
}
