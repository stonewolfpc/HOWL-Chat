package chatpolicy

import (
	"strings"
	"testing"

	"howl-chat/internal/backend/gguf"
)

func TestSpeechRecognitionModelsAreNeverChatEligible(t *testing.T) {
	tests := []struct {
		name string
		path string
		tags []string
	}{
		{
			name: "qwen3-asr-0.6b-tag",
			path: `C:\models\Qwen3-ASR-0.6B-Q8_0.gguf`,
			tags: []string{"automatic-speech-recognition"},
		},
		{
			name: "qwen3-asr-1.7b-identity",
			path: `C:\models\Qwen3-ASR-1.7B-Q8_0.gguf`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eligibility := Evaluate(test.path, gguf.Admission{
				Path:  test.path,
				Class: gguf.AdmissionMultimodal,
				Arch:  "qwen3vl",
				Tags:  test.tags,
			})
			if eligibility.Eligible || eligibility.TextChat {
				t.Fatalf("eligibility = %+v, want ASR refused for chat", eligibility)
			}
			if eligibility.Task != TaskSpeechRecognition {
				t.Fatalf("task = %q, want %q", eligibility.Task, TaskSpeechRecognition)
			}
			if eligibility.AllowsVision() || eligibility.AllowsAudio() {
				t.Fatalf("eligibility = %+v, want no advertised capability", eligibility)
			}
			message := eligibility.RejectionMessage(test.path)
			if !strings.Contains(message, "speech-recognition") || !strings.Contains(message, "ASR route") {
				t.Fatalf("message = %q, want ASR-specific explanation", message)
			}
		})
	}
}

func TestAudioSpecializedModelsAreNotOrdinaryChatModels(t *testing.T) {
	tests := []struct {
		name string
		path string
		meta string
		tags []string
	}{
		{
			name: "ultravox-identity",
			path: `C:\models\Ultravox-Llama-3.2-1B-Q4_K_M.gguf`,
			meta: "Llama 3.2 1B Instruct",
		},
		{
			name: "voxtral-tag",
			path: `C:\models\Voxtral-Mini-3B-2507-Q4_K_M.gguf`,
			tags: []string{"audio-text-to-text"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eligibility := Evaluate(test.path, gguf.Admission{
				Path:  test.path,
				Class: gguf.AdmissionChat,
				Arch:  "llama",
				Name:  test.meta,
				Tags:  test.tags,
			})
			if eligibility.Eligible || eligibility.TextChat {
				t.Fatalf("eligibility = %+v, want audio-specialized model refused", eligibility)
			}
			if eligibility.Task != TaskAudioUnderstanding {
				t.Fatalf("task = %q, want %q", eligibility.Task, TaskAudioUnderstanding)
			}
			if eligibility.AllowsAudio() {
				t.Fatal("an unverified audio capability was advertised")
			}
			message := eligibility.RejectionMessage(test.path)
			if !strings.Contains(message, "audio-specialized") {
				t.Fatalf("message = %q, want audio-specialized explanation", message)
			}
		})
	}
}

func TestOmniAndGemmaRemainTextEligibleWithoutVerifiedMultimodalClaims(t *testing.T) {
	tests := []struct {
		name string
		path string
		meta string
		arch string
		tags []string
	}{
		{
			name: "qwen2.5-omni-3b",
			path: `C:\models\Qwen2.5-Omni-3B-Q4_K_M.gguf`,
			arch: "qwen2vl",
			tags: []string{"multimodal", "any-to-any"},
		},
		{
			name: "gemma-4-e2b-it",
			path: `C:\models\Gemma-4-E2B-it-Q8_0.gguf`,
			meta: "gemma-4-E2B-it",
			arch: "gemma4",
			tags: []string{"any-to-any"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eligibility := Evaluate(test.path, gguf.Admission{
				Path:  test.path,
				Class: gguf.AdmissionChat,
				Arch:  test.arch,
				Name:  test.meta,
				Tags:  test.tags,
			})
			if !eligibility.Eligible || !eligibility.TextChat {
				t.Fatalf("eligibility = %+v, want text-chat eligible", eligibility)
			}
			if eligibility.AllowsVision() || eligibility.AllowsAudio() {
				t.Fatalf("eligibility = %+v, declared multimodal capabilities must not be advertised as verified", eligibility)
			}
			if eligibility.Vision != CapabilityUnverified || eligibility.AudioInput != CapabilityUnverified {
				t.Fatalf("capabilities = vision %q audio %q, want declared-unverified", eligibility.Vision, eligibility.AudioInput)
			}
		})
	}
}

func TestUnknownMultimodalIdentityFailsClosed(t *testing.T) {
	path := `C:\models\mystery-vl.gguf`
	eligibility := Evaluate(path, gguf.Admission{
		Path:  path,
		Class: gguf.AdmissionMultimodal,
		Arch:  "qwen3vl",
	})
	if eligibility.Eligible {
		t.Fatalf("eligibility = %+v, want unknown multimodal identity refused", eligibility)
	}
	if eligibility.Task != TaskUnknown || !strings.Contains(eligibility.Reason, "not established") {
		t.Fatalf("eligibility = %+v, want an explicit unestablished-task refusal", eligibility)
	}
}

func TestUnknownPlainChatArchitectureKeepsTextFallback(t *testing.T) {
	path := `C:\models\some-chat.gguf`
	eligibility := Evaluate(path, gguf.Admission{
		Path:  path,
		Class: gguf.AdmissionChat,
		Arch:  "llama",
	})
	if !eligibility.Eligible || !eligibility.TextChat {
		t.Fatalf("eligibility = %+v, want plain text fallback", eligibility)
	}
	if eligibility.Task != TaskTextChat || eligibility.Basis != "architecture-fallback" {
		t.Fatalf("task = %q basis = %q, want explicit architecture fallback", eligibility.Task, eligibility.Basis)
	}
	if eligibility.AllowsVision() || eligibility.AllowsAudio() {
		t.Fatal("plain text fallback advertised a multimodal capability")
	}
}

func TestLegacyNonGGUFIsDeniedByDefault(t *testing.T) {
	path := `C:\models\ggml-base.bin`
	eligibility := Evaluate(path, gguf.Admission{
		Path:  path,
		Class: gguf.AdmissionNonGGUF,
	})
	if eligibility.Eligible {
		t.Fatalf("eligibility = %+v, want legacy .bin denied", eligibility)
	}
	if !strings.Contains(eligibility.Reason, "denied") {
		t.Fatalf("reason = %q, want default-deny explanation", eligibility.Reason)
	}
}

func TestIdentityUsesMetadataWithoutTruncatingDottedNames(t *testing.T) {
	path := `C:\models\Qwen3-ASR-1.7B-Q8_0.gguf`
	eligibility := Evaluate(path, gguf.Admission{
		Path:  path,
		Class: gguf.AdmissionMultimodal,
		Arch:  "qwen3vl",
		Name:  "Llama 3.2 1B Instruct",
		Tags:  []string{"automatic-speech-recognition"},
	})
	if eligibility.Identity != "llama-3-2-1b-instruct" {
		t.Fatalf("identity = %q, want metadata identity", eligibility.Identity)
	}
	if len(eligibility.Aliases) < 2 || eligibility.Aliases[1] != "qwen3-asr-1-7b" {
		t.Fatalf("aliases = %v, want filename identity as an alias", eligibility.Aliases)
	}
}
