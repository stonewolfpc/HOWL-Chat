package scanner

import "testing"

// qwen_image_vae.safetensors is the autoencoder Anima and Qwen-Image both need.
// The VAE branch of GuessFamily ended in "conservative default: sd1.5", so the
// file was labelled sd1.5 and Anima refused to load with its own VAE.
func TestQwenImageVAEIsNotClassifiedAsSD15(t *testing.T) {
	for _, name := range []string{"qwen_image_vae", "qwen_image_vae.safetensors", "qwen-image-vae.safetensors"} {
		if got := GuessFamily(name); got != FamilyQwenImage {
			t.Errorf("GuessFamily(%q) = %q, want %q", name, got, FamilyQwenImage)
		}
	}
	// The Flux autoencoder stays Flux, and a genuinely unknown VAE keeps the old
	// conservative answer rather than silently becoming something else.
	if got := GuessFamily("ae"); got != FamilyFlux {
		t.Errorf("GuessFamily(ae) = %q, want flux", got)
	}
	if got := GuessFamily("mystery_vae"); got != FamilySD15 {
		t.Errorf("GuessFamily(mystery_vae) = %q, want the sd1.5 default", got)
	}
}

func TestCompatibleVAEAcceptsAnimaAndQwenImageTogether(t *testing.T) {
	if !CompatibleVAE(FamilyAnima, FamilyQwenImage) || !CompatibleVAE(FamilyQwenImage, FamilyAnima) {
		t.Error("Anima and Qwen-Image share an autoencoder and must accept each other")
	}
	if CompatibleVAE(FamilyAnima, FamilySD15) {
		t.Error("Anima accepted an sd1.5 VAE")
	}
	if CompatibleVAE(FamilyAnima, FamilyFlux) {
		t.Error("Anima accepted the Flux autoencoder")
	}
}
