package scanner

import (
	"path/filepath"
	"testing"
)

// Anima is a DiT whose blocks are named net.blocks.N and modulated with adaLN.
// Without that marker the file fell through to the SD1.5 default, so the UI
// offered it as a 512x512 SD1.5 checkpoint.
func TestProbeIdentifiesAnimaFromStructure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "anima-base-v1.0.safetensors")
	writeSafetensors(t, path, nil, map[string][]int{
		"net.blocks.0.adaln_modulation_cross_attn.1.weight": {3072, 18432},
		"net.blocks.0.cross_attn.k_proj.weight":             {3072, 3072},
	})

	probe := ProbeWeight(path)
	if got, ok := probe.Family(); !ok || got != FamilyAnima {
		t.Fatalf("family = %q (ok=%v), want %q", got, ok, FamilyAnima)
	}
	if got, ok := probe.AssetType(); !ok || got != AssetDiffusion {
		t.Fatalf("type = %q (ok=%v), want %q", got, ok, AssetDiffusion)
	}
}

// A GGUF header that declares a language architecture is a language model,
// whatever the file is called and whatever folder it was dropped in. This is
// llama-q2_k.gguf, which the image page listed as an SD1.5 checkpoint.
func TestLanguageGGUFIsNotADiffusionModel(t *testing.T) {
	p := WeightProbe{Metadata: map[string]string{"general.architecture": "llama"}}
	if got, ok := p.AssetType(); !ok || got != AssetLLM {
		t.Fatalf("type = %q (ok=%v), want %q", got, ok, AssetLLM)
	}
	if fam, ok := p.Family(); ok {
		t.Fatalf("a language model must not claim an image family, got %q", fam)
	}
}
