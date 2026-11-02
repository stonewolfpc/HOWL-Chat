package deps

import "testing"

// A companion of the right kind but the wrong model must be reported as a
// mismatch. This is the difference between "you have what you need" and "this
// will fail inside sd-server": a Llama GGUF is an LLM, it is not Qwen2.5-VL.
func TestResolveReportsWrongCompanionAsMismatch(t *testing.T) {
	present := map[string]string{
		"llm": `E:\models\llama-q2_k.gguf`,
		"vae": `E:\models\qwen_image_vae.safetensors`,
	}
	r := Resolve("krea2_turbo-Q4_K_M", "qwen_image", present, "/models/companions")

	if !r.BlocksLoad() {
		t.Fatal("a wrong text encoder must block the load")
	}
	if len(r.Missing) != 1 || r.Missing[0] != "llm" {
		t.Fatalf("missing = %v, want [llm]", r.Missing)
	}
	for _, it := range r.Items {
		if it.Role != "llm" {
			continue
		}
		if it.State != StateMismatch {
			t.Errorf("state = %q, want mismatch", it.State)
		}
		if it.Mismatch == "" {
			t.Error("mismatch reason is empty; the UI has nothing to explain")
		}
	}
}

// The right file for the family is accepted, so a correctly installed setup does
// not get nagged about a download it does not need.
func TestCompatibleAcceptsTheFileTheCatalogueNames(t *testing.T) {
	cases := []struct {
		family, role, name string
		want               bool
	}{
		{"z_image", "llm", "qwen_3_4b.safetensors", true},
		{"z_image", "llm", "llama-q2_k.gguf", false},
		{"qwen_image", "llm", "qwen_2.5_vl_7b_fp8_scaled.safetensors", true},
		{"qwen_image", "llm", "Qwen-Image-Lightning-8steps-V2.0-bf16.safetensors", false},
		{"flux2", "llm", "mistral_3_small_flux2_fp8.safetensors", true},
		{"flux", "vae", "ae.safetensors", true},
		{"flux", "vae", "qwen_image_vae.safetensors", false},
		{"flux", "t5xxl", "t5xxl_fp8_e4m3fn.safetensors", true},
		{"flux", "clip_l", "clip_l.safetensors", true},
		{"flux2", "vae", "flux2-vae.safetensors", true},
	}
	for _, c := range cases {
		got, why := Compatible(c.family, c.role, c.name)
		if got != c.want {
			t.Errorf("Compatible(%s,%s,%s) = %v (%s), want %v", c.family, c.role, c.name, got, why, c.want)
		}
		if !got && why == "" {
			t.Errorf("%s/%s: rejection has no reason", c.family, c.role)
		}
	}
}

// A role this package cannot judge is accepted rather than rejected on a guess.
// Every shipped family now names its own encoders and autoencoders, so the
// unjudgeable case is a family the catalogue has never heard of.
func TestCompatibleDoesNotGuessWhenIdentityIsUnknown(t *testing.T) {
	if ok, why := Compatible("future_family", "vae", "something_vae.safetensors"); !ok {
		t.Fatalf("unknown identity should be accepted, got %v (%s)", ok, why)
	}
}

// A mismatch is downloadable when the catalogue knows the file, so the user can
// replace the wrong file from the same modal that reported it.
func TestMismatchIsOfferedAsADownload(t *testing.T) {
	present := map[string]string{"llm": `E:\models\llama-q2_k.gguf`, "vae": `E:\models\ae.safetensors`}
	r := Resolve("z_image_turbo", "z_image", present, "/models/companions")
	d := r.Downloadable()
	if len(d) != 1 || d[0].Role != "llm" {
		t.Fatalf("downloadable = %+v, want the llm only", d)
	}
	if d[0].URL == "" || d[0].FileName == "" {
		t.Errorf("downloadable entry is incomplete: %+v", d[0])
	}
}
