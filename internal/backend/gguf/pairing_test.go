package gguf

import (
	"errors"
	"path/filepath"
	"testing"
)

func audioProjector(dim int) map[string]any {
	return map[string]any{
		"general.architecture":      "clip",
		"clip.has_audio_encoder":    true,
		"clip.has_vision_encoder":   false,
		"clip.audio.projection_dim": uint32(dim),
		"clip.projector_type":       "test-audio",
	}
}

func visionProjector(dim int) map[string]any {
	return map[string]any{
		"general.architecture":       "clip",
		"clip.has_audio_encoder":     false,
		"clip.has_vision_encoder":    true,
		"clip.vision.projection_dim": uint32(dim),
	}
}

func hybridProjector(dim int) map[string]any {
	return map[string]any{
		"general.architecture":       "clip",
		"clip.has_audio_encoder":     true,
		"clip.has_vision_encoder":    true,
		"clip.audio.projection_dim":  uint32(dim),
		"clip.vision.projection_dim": uint32(dim),
	}
}

func chatModel(arch string, embedding int) map[string]any {
	return map[string]any{
		"general.architecture":     arch,
		arch + ".embedding_length": uint32(embedding),
	}
}

func TestProjectorIdentityNormalisesQuantVersionAndMmproj(t *testing.T) {
	pairs := [][2]string{
		{"Qwen3-ASR-1.7B-Q8_0.gguf", "mmproj-Qwen3-ASR-1.7B-Q8_0.gguf"},
		{"Qwen2.5-Omni-3B-Q4_K_M.gguf", "Qwen2.5-Omni-3B-mmproj-Q8_0.gguf"},
		{"Ultravox-Llama-3.2-1B-Q4_K_M.gguf", "mmproj-Ultravox-v0_5-llama-3_2-1b-f16.gguf"},
		{"Gemma-4-E2B-it-Q8_0.gguf", "mmproj-gemma-4-E2B-it-Q8_0.gguf"},
		{"Voxtral-Mini-3B-2507-Q4_K_M.gguf", "mmproj-Voxtral-Mini-3B-2507-Q8_0.gguf"},
	}
	for _, pair := range pairs {
		if got, want := ProjectorIdentity(pair[0]), ProjectorIdentity(pair[1]); got != want {
			t.Errorf("identity(%q)=%q, identity(%q)=%q; want equal", pair[0], got, pair[1], want)
		}
	}

	mismatches := [][2]string{
		{"Qwen3-ASR-1.7B-Q8_0.gguf", "mmproj-Qwen3-ASR-0.6B-Q8_0.gguf"},
		{"Gemma-4-E2B-it-Q8_0.gguf", "mmproj-gemma-4-E4B-it-Q8_0.gguf"},
	}
	for _, pair := range mismatches {
		if ProjectorIdentity(pair[0]) == ProjectorIdentity(pair[1]) {
			t.Errorf("identity(%q) unexpectedly equals identity(%q)", pair[0], pair[1])
		}
	}
}

func TestModelIdentityDoesNotTruncateDottedMetadataNames(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Llama 3.2 1B Instruct", "llama-3-2-1b-instruct"},
		{"gemma-4-E2B-it", "gemma-4-e2b-it"},
		{"model.gguf", "model"},
	}
	for _, test := range tests {
		if got := ModelIdentity(test.name); got != test.want {
			t.Errorf("ModelIdentity(%q) = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestResolveProjectorPairsKnownAudioInventoryExactly(t *testing.T) {
	dir := t.TempDir()

	type inventoryEntry struct {
		model     string
		arch      string
		embedding int
		projector string
		kv        map[string]any
	}
	inventory := []inventoryEntry{
		{"Gemma-4-E2B-it-Q8_0.gguf", "gemma4", 1536, "mmproj-gemma-4-E2B-it-Q8_0.gguf", hybridProjector(1536)},
		{"Qwen2.5-Omni-3B-Q4_K_M.gguf", "qwen2vl", 2048, "Qwen2.5-Omni-3B-mmproj-Q8_0.gguf", hybridProjector(2048)},
		{"Qwen3-ASR-0.6B-Q8_0.gguf", "qwen3vl", 1024, "mmproj-Qwen3-ASR-0.6B-Q8_0.gguf", audioProjector(1024)},
		{"Qwen3-ASR-1.7B-Q8_0.gguf", "qwen3vl", 2048, "mmproj-Qwen3-ASR-1.7B-Q8_0.gguf", audioProjector(2048)},
		{"Ultravox-Llama-3.2-1B-Q4_K_M.gguf", "llama", 2048, "mmproj-Ultravox-v0_5-llama-3_2-1b-f16.gguf", audioProjector(4096)},
		{"Voxtral-Mini-3B-2507-Q4_K_M.gguf", "llama", 3072, "mmproj-Voxtral-Mini-3B-2507-Q8_0.gguf", audioProjector(3072)},
	}

	modelPaths := make(map[string]string, len(inventory))
	for _, entry := range inventory {
		modelPaths[entry.model] = writeGGUFFile(t, dir, entry.model, chatModel(entry.arch, entry.embedding))
		writeGGUFFile(t, dir, entry.projector, entry.kv)
	}

	for _, entry := range inventory {
		t.Run(entry.model, func(t *testing.T) {
			modelPath := modelPaths[entry.model]
			admission := AdmitModel(modelPath)
			if !admission.AllowedAsChat() {
				t.Fatalf("inventory model was refused: %+v", admission)
			}

			pairing, err := ResolveProjector(PairingRequest{
				ModelPath:        modelPath,
				Arch:             admission.Arch,
				EmbeddingLength:  admission.EmbeddingLength,
				Require:          admission.RequiresProjector(),
				RequiredModality: admission.RequiredModality(),
			})
			if err != nil {
				t.Fatalf("ResolveProjector: %v", err)
			}
			if filepath.Base(pairing.Path) != entry.projector {
				t.Fatalf("paired %q, want %q", filepath.Base(pairing.Path), entry.projector)
			}
		})
	}
}

func TestResolveProjectorNeverChoosesAWrongSibling(t *testing.T) {
	dir := t.TempDir()
	model := writeGGUFFile(t, dir, "Qwen3-ASR-1.7B-Q8_0.gguf", chatModel("qwen3vl", 2048))
	writeGGUFFile(t, dir, "mmproj-Qwen3-ASR-0.6B-Q8_0.gguf", audioProjector(1024))
	correct := writeGGUFFile(t, dir, "mmproj-Qwen3-ASR-1.7B-Q8_0.gguf", audioProjector(2048))

	pairing, err := ResolveProjector(PairingRequest{
		ModelPath:        model,
		Arch:             "qwen3vl",
		EmbeddingLength:  2048,
		Require:          true,
		RequiredModality: ProjectorAudio,
	})
	if err != nil {
		t.Fatalf("ResolveProjector: %v", err)
	}
	if pairing.Path != correct {
		t.Fatalf("paired %q, want %q", pairing.Path, correct)
	}
}

func TestResolveProjectorOmniNamingMatchesByImplicitIdentity(t *testing.T) {
	dir := t.TempDir()
	model := writeGGUFFile(t, dir, "Qwen2.5-Omni-3B-Q4_K_M.gguf", chatModel("qwen2vl", 2048))
	projector := writeGGUFFile(t, dir, "Qwen2.5-Omni-3B-mmproj-Q8_0.gguf", hybridProjector(2048))

	pairing, err := ResolveProjector(PairingRequest{
		ModelPath:        model,
		Arch:             "qwen2vl",
		EmbeddingLength:  2048,
		Require:          true,
		RequiredModality: ProjectorUnknown,
	})
	if err != nil {
		t.Fatalf("ResolveProjector: %v", err)
	}
	if pairing.Path != projector {
		t.Fatalf("paired %q, want %q", pairing.Path, projector)
	}
	if pairing.Info.Kind != ProjectorHybrid {
		t.Fatalf("kind = %q, want hybrid", pairing.Info.Kind)
	}
}

func TestResolveProjectorReportsMissingPair(t *testing.T) {
	dir := t.TempDir()
	model := writeGGUFFile(t, dir, "multimodal-Q4_K_M.gguf", chatModel("qwen2vl", 2048))
	writeGGUFFile(t, dir, "mmproj-unrelated-Q8_0.gguf", hybridProjector(2048))

	_, err := ResolveProjector(PairingRequest{
		ModelPath:        model,
		Arch:             "qwen2vl",
		EmbeddingLength:  2048,
		Require:          true,
		RequiredModality: ProjectorUnknown,
	})
	if !errors.Is(err, ErrProjectorMissing) {
		t.Fatalf("error = %v, want ErrProjectorMissing", err)
	}
}

func TestResolveProjectorRejectsAmbiguousPairing(t *testing.T) {
	dir := t.TempDir()
	model := writeGGUFFile(t, dir, "model-Q4_K_M.gguf", chatModel("llama", 4096))
	writeGGUFFile(t, dir, "mmproj-model-Q8_0.gguf", audioProjector(4096))
	writeGGUFFile(t, dir, "mmproj-model-F16.gguf", audioProjector(4096))

	_, err := ResolveProjector(PairingRequest{
		ModelPath: model,
		Arch:      "llama",
		Require:   true,
	})
	if !errors.Is(err, ErrProjectorAmbiguous) {
		t.Fatalf("error = %v, want ErrProjectorAmbiguous", err)
	}
}

func TestResolveProjectorRejectsWrongModality(t *testing.T) {
	dir := t.TempDir()
	model := writeGGUFFile(t, dir, "vision-model-Q4_K_M.gguf", chatModel("mllama", 4096))
	writeGGUFFile(t, dir, "mmproj-vision-model-Q8_0.gguf", audioProjector(4096))

	_, err := ResolveProjector(PairingRequest{
		ModelPath:        model,
		Arch:             "mllama",
		EmbeddingLength:  4096,
		Require:          true,
		RequiredModality: ProjectorVision,
	})
	if !errors.Is(err, ErrProjectorMismatch) {
		t.Fatalf("error = %v, want ErrProjectorMismatch", err)
	}
}

func TestResolveProjectorExplicitNameIsAuthoritative(t *testing.T) {
	dir := t.TempDir()
	model := writeGGUFFile(t, dir, "Qwen3-ASR-1.7B-Q8_0.gguf", chatModel("qwen3vl", 2048))
	writeGGUFFile(t, dir, "mmproj-Qwen3-ASR-0.6B-Q8_0.gguf", audioProjector(1024))
	want := writeGGUFFile(t, dir, "mmproj-Qwen3-ASR-1.7B-Q8_0.gguf", audioProjector(2048))

	pairing, err := ResolveProjector(PairingRequest{
		ModelPath:        model,
		ExpectedNames:    []string{"mmproj-Qwen3-ASR-1.7B-Q8_0.gguf"},
		Require:          true,
		RequiredModality: ProjectorAudio,
	})
	if err != nil {
		t.Fatalf("ResolveProjector: %v", err)
	}
	if pairing.Path != want {
		t.Fatalf("paired %q, want %q", pairing.Path, want)
	}

	if _, err := ResolveProjector(PairingRequest{
		ModelPath:        model,
		ExpectedNames:    []string{"mmproj-Qwen3-ASR-2.0B-Q8_0.gguf"},
		Require:          true,
		RequiredModality: ProjectorAudio,
	}); !errors.Is(err, ErrProjectorMissing) {
		t.Fatalf("missing explicit projector error = %v, want ErrProjectorMissing", err)
	}
}

func TestResolveProjectorTextModelWithoutCompanionIsAllowed(t *testing.T) {
	dir := t.TempDir()
	model := writeGGUFFile(t, dir, "text-Q4_K_M.gguf", chatModel("llama", 4096))
	writeGGUFFile(t, dir, "mmproj-other-Q4_K_M.gguf", visionProjector(4096))

	pairing, err := ResolveProjector(PairingRequest{ModelPath: model, Arch: "llama"})
	if err != nil {
		t.Fatalf("ResolveProjector: %v", err)
	}
	if pairing.Path != "" {
		t.Fatalf("text model paired with %q, want no projector", pairing.Path)
	}
}

func TestResolveProjectorRecordsDimensionEvidence(t *testing.T) {
	dir := t.TempDir()
	model := writeGGUFFile(t, dir, "Ultravox-Llama-3.2-1B-Q4_K_M.gguf", chatModel("llama", 2048))
	writeGGUFFile(t, dir, "mmproj-Ultravox-v0_5-llama-3_2-1b-f16.gguf", audioProjector(4096))

	pairing, err := ResolveProjector(PairingRequest{
		ModelPath:       model,
		Arch:            "llama",
		EmbeddingLength: 2048,
	})
	if err != nil {
		t.Fatalf("ResolveProjector: %v", err)
	}
	if pairing.ProjectionDim != 4096 || pairing.ModelEmbedding != 2048 {
		t.Fatalf("pairing dimensions = %d/%d, want 4096/2048", pairing.ProjectionDim, pairing.ModelEmbedding)
	}
	if !pairing.DimensionMismatch {
		t.Fatal("expected the known Ultravox dimension difference to be recorded")
	}
}
