package diagnosis

import (
	"errors"
	"strings"
	"testing"
)

func TestCategoriesComeFromTheEngineMessage(t *testing.T) {
	cases := []struct {
		name  string
		stage Stage
		err   string
		want  Category
	}{
		{"unknown architecture, safetensors", StageGeneration, `sd-server txt2img error 400: {"error":"cannot identify model family from 'H:\\models\\mystery.safetensors'"}`, UnsupportedArchitecture},
		{"supported architecture, unsupported format", StageGeneration, `sd-server txt2img error 400: {"error":"cannot identify model family from 'H:\\models\\hidream-i1-fast-Q2_K.gguf'"}`, UnsupportedWeightFormat},
		{"torch checkpoint", StageGeneration, `sd-server txt2img error 400: {"error":"cannot identify model family from 'H:\\models\\HiDream.pt'"}`, UnsupportedWeightFormat},
		{"missing file", StageStartup, "failed to load model: no such file or directory", InvalidModelFile},
		{"admission refusal", StageStartup, `cannot start image "x": needs 11761 MB VRAM; 6233 MB available`, InsufficientResources},
		{"bad backend", StageStartup, "error: unknown argument: --backend diffusion=nowhere", PlacementFailure},
		{"engine died", StageGeneration, "image engine stopped before finishing: sd-server exited", EngineStartupFailure},
		{"opaque generation", StageGeneration, "sd-server txt2img error 500: internal error", GenerationFailure},
		{"opaque startup", StageStartup, "something went wrong", EngineStartupFailure},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := Diagnose(tt.stage, errors.New(tt.err))
			if got.Category != tt.want {
				t.Fatalf("category = %q, want %q", got.Category, tt.want)
			}
			if got.Summary == "" || got.Suggestion == "" {
				t.Fatal("every diagnosis must explain itself and say what to try")
			}
			if !strings.Contains(got.Raw, strings.TrimSpace(tt.err)) {
				t.Fatalf("diagnosis dropped the engine's own words: %q", got.Raw)
			}
		})
	}
}

// The headline case: the engine does not know the architecture. The user must
// not be sent to the VRAM settings for it.
func TestUnsupportedArchitectureDoesNotLookLikeAMemoryProblem(t *testing.T) {
	d := Diagnose(StageGeneration, errors.New("cannot identify model family from 'mystery.safetensors'"))
	if d.Category != UnsupportedArchitecture {
		t.Fatalf("category = %q, want unsupported_architecture", d.Category)
	}
	if d.Category == InsufficientResources {
		t.Fatal("an unidentifiable architecture must not be classified as a memory problem")
	}
	// Resource remedies must not appear: this category cannot be fixed by
	// spending less memory.
	for _, remedy := range []string{"lower the resolution", "unload another", "more quantized"} {
		if strings.Contains(strings.ToLower(d.Suggestion), remedy) {
			t.Fatalf("suggestion offers a memory remedy for an architecture problem: %s", d.Suggestion)
		}
	}
	if !strings.Contains(strings.ToLower(d.Suggestion), "engine") {
		t.Fatalf("suggestion must point at the engine or the model file: %s", d.Suggestion)
	}
	if !strings.Contains(d.Describe(), "unsupported_architecture") {
		t.Fatalf("description must name the category: %s", d.Describe())
	}
}

// A rejected GGUF must not be reported as an architecture the engine lacks:
// the same family is documented as supported in Safetensors, and the user's
// next action is a different file, not a different engine.
func TestRejectedFormatIsNotAnUnsupportedArchitecture(t *testing.T) {
	gguf := Diagnose(StageGeneration, errors.New("cannot identify model family from 'H:\\models\\hidream-i1-fast-Q2_K.gguf'"))
	safe := Diagnose(StageGeneration, errors.New("cannot identify model family from 'H:\\models\\hidream-i1-dev.safetensors'"))

	if gguf.Category != UnsupportedWeightFormat {
		t.Fatalf("GGUF rejection category = %q, want unsupported_weight_format", gguf.Category)
	}
	if safe.Category != UnsupportedArchitecture {
		t.Fatalf("safetensors rejection category = %q, want unsupported_architecture", safe.Category)
	}
	if gguf.Category == safe.Category {
		t.Fatal("format rejection and unknown architecture must not share a category")
	}
	if !strings.Contains(strings.ToLower(gguf.Suggestion), "safetensors") {
		t.Fatalf("format rejection must point at the supported release: %s", gguf.Suggestion)
	}
	if gguf.Retryable || safe.Retryable {
		t.Fatal("retrying the same file cannot fix either case")
	}
}

func TestWrapIsIdempotent(t *testing.T) {
	first := Wrap(StageGeneration, errors.New("cannot identify model family from 'x.gguf'"))
	second := Wrap(StageStartup, first)
	if first.Error() != second.Error() {
		t.Fatalf("re-wrapping changed the diagnosis:\n%s\n%s", first.Error(), second.Error())
	}
	d, ok := As(second)
	if !ok || d.Stage != StageGeneration {
		t.Fatalf("stage was rewritten by a later wrap: %+v", d)
	}
	if Wrap(StageGeneration, nil) != nil {
		t.Fatal("wrapping nil must stay nil")
	}
}
