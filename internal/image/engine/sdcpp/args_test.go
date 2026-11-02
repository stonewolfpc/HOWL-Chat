package sdcpp

import (
	"strings"
	"testing"

	"howl-chat/internal/image/types"
)

// Anima and Z-Image were started with --t5xxl and --clip_l because the caller
// supplied them; neither family reads either encoder, and at the time t5xxl was
// still being downloaded, so a truncated file was handed to the loader. A family
// that declares its roles gets exactly those roles.
func TestBuildServerArgsOnlySendsCompanionsTheFamilyReads(t *testing.T) {
	e := NewEngine()
	e.config.Port = 8081
	m := &types.ModelInfo{
		Name:   "anima-base-v1.0",
		Family: "anima",
		Files: map[string]string{
			"model":  `H:\models\anima-base-v1.0.safetensors`,
			"llm":    `H:\models\qwen_3_06b_base.safetensors`,
			"vae":    `E:\models\qwen_image_vae.safetensors`,
			"t5xxl":  `E:\models\t5xxl_fp8_e4m3fn.safetensors`,
			"clip_l": `E:\models\clip_l.safetensors`,
		},
	}
	args := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{}), " ")

	if !strings.Contains(args, "--diffusion-model") {
		t.Errorf("split-component family must use --diffusion-model: %s", args)
	}
	if !strings.Contains(args, "--llm") || !strings.Contains(args, "--vae") {
		t.Errorf("required companions missing: %s", args)
	}
	if strings.Contains(args, "--t5xxl") || strings.Contains(args, "--clip_l") {
		t.Errorf("family does not read T5 or CLIP-L but they were passed: %s", args)
	}
}

// The scoping must not undo the case it exists for: a family this build does not
// know has no role list, so a supplied encoder still reaches the loader.
func TestBuildServerArgsKeepsSuppliedCompanionsForUnknownFamily(t *testing.T) {
	e := NewEngine()
	e.config.Port = 8081
	m := &types.ModelInfo{
		Name:   "mystery-model",
		Family: "not_a_family",
		Files: map[string]string{
			"model": `E:\models\mystery.safetensors`,
			"llm":   `E:\models\encoder.safetensors`,
			"vae":   `E:\models\vae.safetensors`,
		},
	}
	args := strings.Join(e.buildServerArgs(m, types.ModelLoadOptions{}), " ")
	if !strings.Contains(args, "--llm") || !strings.Contains(args, "--vae") {
		t.Errorf("unknown family dropped its supplied companions: %s", args)
	}
}
