package gguf

import "testing"

func TestKindOfArch(t *testing.T) {
	cases := map[string]Kind{
		// Chat models
		"llama":   KindLLM,
		"qwen3":   KindLLM,
		"qwen3vl": KindLLM,
		"gemma3":  KindLLM,
		"mistral": KindLLM,
		// Image models
		"flux":       KindDiffusion,
		"sdxl":       KindDiffusion,
		"sd3.5":      KindDiffusion,
		"wan":        KindDiffusion,
		"chroma":     KindDiffusion,
		"qwen_image": KindDiffusion,
		// Not a chat model nor an image model: leave the decision to the caller
		"":             KindUnknown,
		"whisper":      KindUnknown,
		"t5":           KindUnknown,
		"clip":         KindUnknown,
		"brandnewarch": KindUnknown,
	}

	for arch, want := range cases {
		if got := KindOfArch(arch); got != want {
			t.Errorf("KindOfArch(%q) = %q, want %q", arch, got, want)
		}
	}
}

func TestKindOfArchNormalisesInput(t *testing.T) {
	if got := KindOfArch("  Qwen3  "); got != KindLLM {
		t.Errorf("expected whitespace and case to be ignored, got %q", got)
	}
	if got := KindOfArch("FLUX"); got != KindDiffusion {
		t.Errorf("expected uppercase FLUX to classify as diffusion, got %q", got)
	}
}
