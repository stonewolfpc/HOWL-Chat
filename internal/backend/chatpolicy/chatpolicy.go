package chatpolicy

// Chat task eligibility is deliberately separate from GGUF architecture
// admission.
//
// AdmitModel answers "what is this file?" — a chat architecture, a multimodal
// architecture, a component, or something unverified. It must not answer "is
// this model suitable for the chat picker?". Qwen3-ASR declares qwen3vl, the
// same architecture used by vision-language chat models, and Ultravox/Voxtral
// declare llama, the same architecture used by text chat models. Architecture
// alone therefore cannot establish the task.
//
// This package adds a second, explicit decision based on declared model tags,
// stable model identity and the architecture fallback. It never imports the
// ASR registry: ASR admission remains its own authority. A model can stay in
// the registry as inventory while being ineligible for the chat picker.
//
// EVIDENCE BASIS (observed 2026-11-05 in the local inventory)
//
//   - Qwen3-ASR 0.6B/1.7B: general.tags=automatic-speech-recognition, qwen3vl.
//   - Qwen2.5-Omni 3B: general.tags={multimodal,any-to-any}, qwen2vl.
//   - Gemma 4 E2B: general.tags={any-to-any}, gemma4, general.name present.
//   - Voxtral Mini 3B: general.tags={transformers,vllm,audio-text-to-text},
//     llama, general.name absent.
//   - Ultravox 0.5 1B: general.tags={...text-generation}, llama, and
//     general.name is the generic base name "Llama 3.2 1B Instruct"; its
//     identity is recovered from the quantisation-stripped filename.
//
// LIMITS
//
//   - Unrecognized plain text-chat architectures remain text-chat eligible as
//     an architecture fallback, with no verified capability and no multimodal
//     claim. Otherwise every ordinary chat GGUF would disappear.
//   - Unrecognized multimodal architectures fail closed: their task could be
//     ASR, vision chat or audio chat, and none of those is established.
//   - Non-GGUF legacy files are denied by default. The old blanket .bin
//     admission let a Whisper component into the chat picker; any future legacy
//     support must be an explicit, validated allowlist, not an extension rule.
//   - A capability is advertised only when this policy marks it verified. A
//     projector declaration alone is evidence of a component, not proof that
//     the chat request path can use it.

import (
	"fmt"
	"path/filepath"
	"strings"

	"howl-chat/internal/backend/gguf"
)

// Task is the model's declared task for chat-admission purposes.
type Task string

const (
	// TaskTextChat is ordinary text chat without a multimodal claim.
	TaskTextChat Task = "text-chat"
	// TaskMultimodalChat is a chat model whose multimodal input is declared but
	// not verified in this request path.
	TaskMultimodalChat Task = "multimodal-chat"
	// TaskSpeechRecognition is an ASR model. It belongs to the ASR route.
	TaskSpeechRecognition Task = "speech-recognition"
	// TaskAudioUnderstanding is an audio-specialized model whose direct chat
	// audio path is not implemented.
	TaskAudioUnderstanding Task = "audio-understanding"
	// TaskUnknown is an ineligible or unestablished task.
	TaskUnknown Task = "unknown"
)

// Capability records whether a capability may be advertised. A projector
// declaration is only CapabilityUnverified until a policy marks it verified.
type Capability string

const (
	// CapabilityNone means the policy makes no capability claim.
	CapabilityNone Capability = "none"
	// CapabilityUnverified is declared by metadata but not proven in the chat
	// request path. It is never advertised.
	CapabilityUnverified Capability = "declared-unverified"
	// CapabilityVerified is a capability that has been exercised end to end.
	CapabilityVerified Capability = "verified"
)

// Eligibility is the chat-admission decision for one model file.
type Eligibility struct {
	Identity     string
	Aliases      []string
	Architecture string
	Task         Task
	Eligible     bool
	TextChat     bool
	Vision       Capability
	AudioInput   Capability
	// Verified is false unless the task itself was exercised end to end. The
	// ordinary chat path has no such verification yet.
	Verified bool
	// Basis names the evidence that decided the task.
	Basis string
	// Reason explains the decision for logs and user-facing refusals.
	Reason string
}

// AllowsVision reports whether image input may be advertised.
func (e Eligibility) AllowsVision() bool {
	return e.Vision == CapabilityVerified
}

// AllowsAudio reports whether audio input may be advertised.
func (e Eligibility) AllowsAudio() bool {
	return e.AudioInput == CapabilityVerified
}

// RejectionMessage is the user-facing refusal for an ineligible model.
func (e Eligibility) RejectionMessage(path string) string {
	base := filepath.Base(path)
	switch e.Task {
	case TaskSpeechRecognition:
		return fmt.Sprintf("%s is a speech-recognition model, not a chat model. It remains available through the ASR route and cannot be selected in the chat picker.", base)
	case TaskAudioUnderstanding:
		return fmt.Sprintf("%s is an audio-specialized model. Direct chat audio is not implemented and ordinary text chat is not verified, so it cannot be selected in the chat picker.", base)
	default:
		return fmt.Sprintf("%s is not eligible for the chat picker: %s", base, e.Reason)
	}
}

// Evaluate decides chat eligibility from architecture admission plus declared
// metadata. It does not read files and is safe to call at scan, picker and load
// boundaries.
func Evaluate(path string, admission gguf.Admission) Eligibility {
	identity, aliases := identityCandidates(path, admission)
	tags := tagSet(admission.Tags)

	eligibility := Eligibility{
		Identity:     identity,
		Aliases:      aliases,
		Architecture: admission.Arch,
		Task:         TaskUnknown,
		Vision:       CapabilityNone,
		AudioInput:   CapabilityNone,
	}

	if admission.Class == gguf.AdmissionNonGGUF {
		eligibility.Basis = "non-gguf-default-deny"
		eligibility.Reason = "legacy non-GGUF files are denied for chat by default; no explicit validated legacy allowlist matched"
		return eligibility
	}
	if admission.Class != gguf.AdmissionChat && admission.Class != gguf.AdmissionMultimodal {
		eligibility.Basis = "not-a-chat-architecture"
		eligibility.Reason = fmt.Sprintf("GGUF architecture %q is not a chat architecture", admission.Arch)
		return eligibility
	}

	if tags["automatic-speech-recognition"] || hasAlias(aliases, speechRecognitionAliases) {
		eligibility.Task = TaskSpeechRecognition
		eligibility.Basis = basisFor("automatic-speech-recognition", aliases, speechRecognitionAliases)
		eligibility.Reason = "the model declares speech recognition, not chat; ASR has its own registry and request path"
		return eligibility
	}

	if tags["audio-text-to-text"] || hasAlias(aliases, audioUnderstandingAliases) {
		eligibility.Task = TaskAudioUnderstanding
		eligibility.AudioInput = CapabilityUnverified
		eligibility.Basis = basisFor("audio-text-to-text", aliases, audioUnderstandingAliases)
		eligibility.Reason = "the model is audio-specialized; direct audio understanding is blocked and ordinary text chat is unverified"
		return eligibility
	}

	if tags["any-to-any"] || tags["multimodal"] || tags["image-text-to-text"] ||
		hasAlias(aliases, multimodalChatAliases) {
		eligibility.Task = TaskMultimodalChat
		eligibility.Eligible = true
		eligibility.TextChat = true
		eligibility.Vision = CapabilityUnverified
		eligibility.AudioInput = CapabilityUnverified
		eligibility.Basis = basisFor("any-to-any/multimodal", aliases, multimodalChatAliases)
		eligibility.Reason = "text chat is eligible; vision and audio are declared but not verified in the chat request path and are not advertised"
		return eligibility
	}

	if admission.Class == gguf.AdmissionMultimodal {
		eligibility.Basis = "unrecognized-multimodal-identity"
		eligibility.Reason = fmt.Sprintf("multimodal architecture %q has no recognized task identity; task capability is not established", admission.Arch)
		return eligibility
	}

	// A plain text-chat architecture without a recognized special-task identity
	// remains eligible for text chat. This is an explicit architecture fallback,
	// not a claim that the model has been exercised.
	eligibility.Task = TaskTextChat
	eligibility.Eligible = true
	eligibility.TextChat = true
	eligibility.Basis = "architecture-fallback"
	eligibility.Reason = "unrecognized model identity; plain text chat is eligible at architecture level only and no multimodal capability is advertised"
	return eligibility
}

var (
	speechRecognitionAliases = []string{
		"qwen3-asr-0-6b",
		"qwen3-asr-1-7b",
	}
	audioUnderstandingAliases = []string{
		"ultravox-llama-3-2-1b",
		"ultravox-llama-3-2-1b-instruct",
		"voxtral-mini-3b-2507",
		"voxtral-mini-4b-realtime-2602",
	}
	multimodalChatAliases = []string{
		"qwen2-5-omni-3b",
		"qwen2-5-omni-7b",
		"qwen3-omni-30b",
		"qwen3-omni-30b-a3b",
		"gemma-4-e2b",
		"gemma-4-e2b-it",
		"gemma-4-e4b",
		"gemma-4-e4b-it",
	}
)

func identityCandidates(path string, admission gguf.Admission) (string, []string) {
	var candidates []string
	if name := strings.TrimSpace(admission.Name); name != "" {
		candidates = append(candidates, gguf.ModelIdentity(name))
	}
	candidates = append(candidates, gguf.ProjectorIdentity(path))

	out := make([]string, 0, len(candidates))
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if _, ok := seen[candidate]; ok {
			continue
		}
		seen[candidate] = struct{}{}
		out = append(out, candidate)
	}
	if len(out) == 0 {
		return "", nil
	}
	return out[0], out
}

func tagSet(tags []string) map[string]bool {
	out := make(map[string]bool, len(tags))
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag != "" {
			out[tag] = true
		}
	}
	return out
}

func hasAlias(aliases, wanted []string) bool {
	for _, alias := range aliases {
		for _, want := range wanted {
			if alias == want {
				return true
			}
		}
	}
	return false
}

func basisFor(tag string, aliases, wanted []string) string {
	for _, alias := range aliases {
		for _, want := range wanted {
			if alias == want {
				return "identity=" + alias
			}
		}
	}
	return "general.tags=" + tag
}
