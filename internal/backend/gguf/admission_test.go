package gguf

import (
	"strings"
	"testing"
)

// The exact architecture string observed in the LFM2.5 vocoder GGUF.
const observedVocoderArch = "this model cannot be used as LLM, use it via --model-vocoder in TTS examples"

func TestAdmitModelRejectsVocoderAsChat(t *testing.T) {
	path := writeGGUFFile(t, t.TempDir(), "vocoder-LFM2.5-Audio-1.5B-Q4_0.gguf", map[string]any{
		"general.architecture": observedVocoderArch,
	})

	admission := AdmitModel(path)
	if admission.Class != AdmissionComponent {
		t.Fatalf("class = %q, want component", admission.Class)
	}
	if admission.AllowedAsChat() {
		t.Fatal("the vocoder was admitted as a chat model")
	}
	message := admission.ChatRejectionMessage()
	if !strings.Contains(message, "cannot be used as a chat model") {
		t.Fatalf("message %q does not explain the refusal", message)
	}
	if !strings.Contains(message, "vocoder") {
		t.Fatalf("message %q does not identify the asset", message)
	}
}

func TestAdmitModelRejectsProjectorComponent(t *testing.T) {
	path := writeGGUFFile(t, t.TempDir(), "mmproj-Qwen3-ASR-1.7B-Q8_0.gguf", map[string]any{
		"general.architecture":   "clip",
		"clip.has_audio_encoder": true,
	})

	admission := AdmitModel(path)
	if admission.Class != AdmissionComponent {
		t.Fatalf("class = %q, want component", admission.Class)
	}
	if admission.AllowedAsChat() {
		t.Fatal("an mmproj was admitted as a chat model")
	}
	if !strings.Contains(admission.ChatRejectionMessage(), "projector") {
		t.Fatalf("message %q does not call out the projector role", admission.ChatRejectionMessage())
	}
}

func TestAdmitModelFailsClosedForUnknownArchitecture(t *testing.T) {
	path := writeGGUFFile(t, t.TempDir(), "mystery.gguf", map[string]any{
		"general.architecture": "brandnewarch",
	})

	admission := AdmitModel(path)
	if admission.Class != AdmissionUnknown {
		t.Fatalf("class = %q, want unknown", admission.Class)
	}
	if admission.AllowedAsChat() {
		t.Fatal("an unknown architecture was admitted")
	}
	if !strings.Contains(admission.ChatRejectionMessage(), "brandnewarch") {
		t.Fatalf("message %q does not name the architecture", admission.ChatRejectionMessage())
	}
}

func TestAdmitModelRejectsMissingArchitecture(t *testing.T) {
	path := writeGGUFFile(t, t.TempDir(), "no-arch.gguf", map[string]any{
		"general.name": "nameless",
	})

	admission := AdmitModel(path)
	if admission.Class != AdmissionUnknown || admission.AllowedAsChat() {
		t.Fatalf("admission = %+v, want unknown and refused", admission)
	}
}

func TestAdmitModelRejectsDiffusionArchitecture(t *testing.T) {
	path := writeGGUFFile(t, t.TempDir(), "flux.gguf", map[string]any{
		"general.architecture": "flux",
	})

	admission := AdmitModel(path)
	if admission.Class != AdmissionUnsupported {
		t.Fatalf("class = %q, want unsupported", admission.Class)
	}
	if admission.AllowedAsChat() {
		t.Fatal("a diffusion model was admitted as chat")
	}
}

func TestAdmitModelAdmitsRegisteredChatFamilies(t *testing.T) {
	for _, arch := range []string{"llama", "mistral", "gemma4", "qwen35", "gpt-oss", "lfm2moe", "cohere2", "chatglm"} {
		t.Run(arch, func(t *testing.T) {
			path := writeGGUFFile(t, t.TempDir(), arch+"-model.gguf", map[string]any{
				"general.architecture": arch,
			})
			admission := AdmitModel(path)
			if admission.Class != AdmissionChat {
				t.Fatalf("class = %q, want chat", admission.Class)
			}
			if !admission.AllowedAsChat() {
				t.Fatal("a registered chat family was refused")
			}
			if admission.RequiresProjector() {
				t.Fatal("a text chat family unexpectedly requires a projector")
			}
		})
	}
}

func TestAdmitModelMarksMultimodalArchitectures(t *testing.T) {
	for _, arch := range []string{"qwen2vl", "qwen3vl"} {
		t.Run(arch, func(t *testing.T) {
			path := writeGGUFFile(t, t.TempDir(), arch+"-model.gguf", map[string]any{
				"general.architecture":     arch,
				arch + ".embedding_length": uint32(2048),
				arch + ".block_count":      uint32(28),
			})
			admission := AdmitModel(path)
			if admission.Class != AdmissionMultimodal {
				t.Fatalf("class = %q, want multimodal-chat", admission.Class)
			}
			if !admission.RequiresProjector() {
				t.Fatal("a multimodal architecture did not require a projector")
			}
			if admission.EmbeddingLength != 2048 {
				t.Fatalf("embedding length = %d, want 2048", admission.EmbeddingLength)
			}
			if !admission.AllowedAsChat() {
				t.Fatal("a supported multimodal architecture was refused")
			}
		})
	}
}

func TestAdmitModelDeniesLegacyNonGGUFByDefault(t *testing.T) {
	admission := AdmitModel("legacy-model.bin")
	if admission.Class != AdmissionNonGGUF || admission.AllowedAsChat() {
		t.Fatalf("admission = %+v, want non-gguf and denied", admission)
	}
	if !strings.Contains(admission.ChatRejectionMessage(), "denied for chat by default") {
		t.Fatalf("message = %q, want default-deny explanation", admission.ChatRejectionMessage())
	}
}

func TestAdmitModelReadsGeneralTagsForTaskPolicy(t *testing.T) {
	path := writeGGUFFile(t, t.TempDir(), "Qwen3-ASR-1.7B-Q8_0.gguf", map[string]any{
		"general.architecture": "qwen3vl",
		"general.tags":         []string{"automatic-speech-recognition"},
	})

	admission := AdmitModel(path)
	if admission.Class != AdmissionMultimodal {
		t.Fatalf("class = %q, want multimodal", admission.Class)
	}
	if len(admission.Tags) != 1 || admission.Tags[0] != "automatic-speech-recognition" {
		t.Fatalf("tags = %v, want the declared ASR tag", admission.Tags)
	}
}
