package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

// Z-Image is told apart by its refiner blocks, which no other supported
// architecture has. This is the model Invoke stores as z_image_turbo-Q4_K.gguf
// next to a Qwen3 encoder.
func TestProbeIdentifiesZImageFromStructure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "z_image_turbo.safetensors")
	writeSafetensors(t, path, nil, map[string][]int{
		"noise_refiner.0.attention.to_q.weight": {3840, 3840},
		"layers.0.attention.to_q.weight":        {3840, 3840},
	})

	probe := ProbeWeight(path)
	if got, ok := probe.Family(); !ok || got != FamilyZImage {
		t.Fatalf("family = %q (ok=%v), want %q", got, ok, FamilyZImage)
	}
	if got, ok := probe.AssetType(); !ok || got != AssetDiffusion {
		t.Fatalf("type = %q (ok=%v), want %q", got, ok, AssetDiffusion)
	}
}

// Flux.2 must not be mistaken for Flux.1: it uses a different text encoder and
// a different block layout, and the engine loads them differently.
func TestProbeSeparatesFlux2FromFlux(t *testing.T) {
	dir := t.TempDir()
	flux2 := filepath.Join(dir, "flux2.safetensors")
	writeSafetensors(t, flux2, map[string]string{
		"modelspec.architecture": "flux-2-dev",
	}, map[string][]int{
		"double_blocks.0.img_attn.proj.weight": {3072, 3072},
	})
	if got, _ := ProbeWeight(flux2).Family(); got != FamilyFlux2 {
		t.Fatalf("flux2 family = %q, want %q", got, FamilyFlux2)
	}

	flux1 := filepath.Join(dir, "flux1.safetensors")
	writeSafetensors(t, flux1, map[string]string{
		"modelspec.architecture": "flux-1-dev",
	}, map[string][]int{
		"double_blocks.0.img_attn.proj.weight": {3072, 3072},
	})
	if got, _ := ProbeWeight(flux1).Family(); got != FamilyFlux {
		t.Fatalf("flux1 family = %q, want %q", got, FamilyFlux)
	}
}

// Qwen-Image's GGUF metadata says "qwen", which used to classify the diffusion
// weights as a language model.
func TestProbeIdentifiesQwenImageAsDiffusion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "qwen-image.safetensors")
	writeSafetensors(t, path, map[string]string{
		"modelspec.architecture": "qwen_image",
	}, map[string][]int{
		"img_in.weight": {3840, 64},
		"txt_in.weight": {3840, 3584},
	})

	probe := ProbeWeight(path)
	if got, ok := probe.AssetType(); !ok || got != AssetDiffusion {
		t.Fatalf("type = %q (ok=%v), want %q", got, ok, AssetDiffusion)
	}
	if got, _ := probe.Family(); got != FamilyQwenImage {
		t.Fatalf("family = %q, want %q", got, FamilyQwenImage)
	}
}

// A UUID folder holding z_image_turbo-Q4_K.gguf named "model" still gets a
// usable name in the menu.
func TestDisplayNameForNewFamilies(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "b4f4d36f-ac8b-4141-80f4-47a8e089685d")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if got := DisplayName("model", dir, WeightProbe{}, FamilyZImage); got != "Z-Image b4f4d36f" {
		t.Fatalf("display name = %q, want %q", got, "Z-Image b4f4d36f")
	}
}
