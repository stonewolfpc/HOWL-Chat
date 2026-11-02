package sdcpp

import (
	"strings"
	"testing"

	"howl-chat/internal/image/scanner"
	"howl-chat/internal/image/types"
)

func TestCapabilityTableCoversTheShippedArchitectures(t *testing.T) {
	for _, family := range []string{
		"sd1.5", "sd2.1", "sdxl", "sdxl-turbo", "sd3", "sd3.5",
		"flux", "flux2", "chroma", "auraflow", "anima", "z_image", "qwen_image", "hidream",
	} {
		if _, ok := Capability(family); !ok {
			t.Errorf("%s missing from the capability table", family)
		}
	}
	if _, ok := Capability("not_a_family"); ok {
		t.Error("unknown family must not resolve")
	}
	if _, ok := Capability(""); ok {
		t.Error("empty family must not resolve")
	}
}

func TestCapabilityCompanionRequirements(t *testing.T) {
	cases := []struct {
		family string
		split  bool
		need   []string
	}{
		// LLM-encoder architectures: Qwen3/Qwen2VL/Mistral instead of CLIP+T5.
		{"z_image", true, []string{"llm", "vae"}},
		{"qwen_image", true, []string{"llm", "vae"}},
		{"flux2", true, []string{"llm", "vae"}},
		{"flux", true, []string{"t5xxl", "clip_l", "vae"}},
		{"sd1.5", false, nil},
	}
	for _, c := range cases {
		cap, ok := Capability(c.family)
		if !ok {
			t.Fatalf("%s not in table", c.family)
		}
		if cap.SplitComponents != c.split {
			t.Errorf("%s SplitComponents = %v, want %v", c.family, cap.SplitComponents, c.split)
		}
		if len(cap.Companions) != len(c.need) {
			t.Errorf("%s companions = %v, want %v", c.family, cap.Companions, c.need)
			continue
		}
		for i, role := range c.need {
			if cap.Companions[i] != role {
				t.Errorf("%s companion[%d] = %q, want %q", c.family, i, cap.Companions[i], role)
			}
		}
	}
}

// A fully wired Z-Image model must reach sd-server with --diffusion-model and
// --llm. Before the table, --llm was only emitted for anima, so this model
// loaded without its text encoder and failed inside the engine.
func TestBuildServerArgsUsesCapabilityTable(t *testing.T) {
	e := &Engine{config: types.EngineConfig{Port: 8081}}
	m := &types.ModelInfo{
		Name:   "z-image-turbo",
		Family: "z_image",
		Files: map[string]string{
			"model": `C:\models\z_image_turbo-Q4_K.gguf`,
			"llm":   `C:\models\Qwen_3_4b-Q6_K.gguf`,
			"vae":   `C:\models\qwen_image_vae.safetensors`,
		},
	}
	args := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{}), " ")

	for _, want := range []string{"--diffusion-model", "--llm", "--vae"} {
		if !strings.Contains(args, want) {
			t.Errorf("args missing %s: %s", want, args)
		}
	}
	if strings.Contains(args, " -m ") {
		t.Errorf("split-component family must not use -m: %s", args)
	}
}

// Chroma's always-on flags live in the table now, so they cannot be dropped by
// an edit to the argument builder.
func TestBuildServerArgsAppliesFamilyFlags(t *testing.T) {
	e := &Engine{config: types.EngineConfig{Port: 8081}}
	m := &types.ModelInfo{
		Name:   "chroma",
		Family: "chroma",
		Files: map[string]string{
			"model": `C:\models\chroma.safetensors`,
			"t5xxl": `C:\models\t5xxl.safetensors`,
			"vae":   `C:\models\ae.safetensors`,
		},
	}
	args := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{}), " ")
	for _, want := range []string{"--chroma-disable-dit-mask", "--clip-on-cpu", "--t5xxl"} {
		if !strings.Contains(args, want) {
			t.Errorf("chroma args missing %s: %s", want, args)
		}
	}
	if strings.Contains(args, "--clip_l") {
		t.Errorf("chroma has no CLIP-L and must not receive the flag: %s", args)
	}
}

func TestModelLoadabilityFollowsTable(t *testing.T) {
	// Z-Image needs both an LLM text encoder and a VAE.
	complete := scanner.ModelIndex{
		LLMs: []scanner.AssetEntry{{Name: "qwen3", Path: `C:\m\qwen3.safetensors`}},
		VAEs: []scanner.AssetEntry{{Name: "qwen_image_vae", Path: `C:\m\vae.safetensors`}},
	}
	ready := &types.ModelInfo{Family: "z_image", Files: map[string]string{"model": `C:\m\z.gguf`}}
	if !isModelLoadable(ready, complete) {
		t.Error("z_image with an LLM and a VAE in the index must be loadable")
	}

	// An LLM alone is not enough, and the reason must name what is missing.
	partial := &types.ModelInfo{Family: "z_image", Files: map[string]string{"model": `C:\m\z.gguf`}}
	if isModelLoadable(partial, scanner.ModelIndex{LLMs: complete.LLMs}) {
		t.Error("z_image without a VAE must not be loadable")
	}
	if !strings.Contains(partial.Compatibility, "vae") {
		t.Errorf("compatibility should name the missing VAE, got %q", partial.Compatibility)
	}
}

// Roles a family never uses are stripped, so a stale companion from an earlier
// classification cannot be attached to a load.
func TestAutoFillDropsUnusedRoles(t *testing.T) {
	info := &types.ModelInfo{
		Family: "sd1.5",
		Files:  map[string]string{"model": `C:\m\x.safetensors`, "t5xxl": `C:\m\t5.safetensors`},
	}
	autoFillCompanions(info, scanner.ModelIndex{})
	if info.Files["t5xxl"] != "" {
		t.Error("sd1.5 must not carry a T5XXL encoder")
	}
}
