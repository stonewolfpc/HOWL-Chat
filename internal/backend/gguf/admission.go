package gguf

// Model admission: decide whether a GGUF file may be loaded as a chat model.
//
// WHY THIS EXISTS
//
// The model registry recorded whatever file it found, so a vocoder (an audio
// component with no chat ability at all) appeared in the chat model picker and
// was handed to llama-server, which refused it with a one-line error. The file
// extension and the file name are not evidence of capability; the GGUF
// architecture is. This classifier reads general.architecture from the file
// itself and fails closed for everything it cannot positively identify.
//
// EVIDENCE BASIS (observed 2026-11-05 against the local inventory)
//
//   - vocoder-LFM2.5-Audio-1.5B-Q4_0.gguf declares the literal string
//     "this model cannot be used as LLM, use it via --model-vocoder in TTS
//     examples" as general.architecture. llama-server then exits 1.
//   - Qwen3-ASR-1.7B-Q8_0.gguf declares qwen3vl; its projector declares
//     clip.has_audio_encoder=true and clip.audio.projection_dim=2048.
//   - Qwen2.5-Omni-3B-Q4_K_M.gguf declares qwen2vl; its projector declares
//     both clip.has_vision_encoder and clip.has_audio_encoder.
//   - Projectors (mmproj-*.gguf) declare general.architecture=clip, which is
//     an encoder container, never a chat LLM.
//
// LIMITATIONS (deliberate, and why failure is closed)
//
//   - llmArchs is a curated list, not a copy of llama.cpp's architecture
//     table, and this process cannot ask llama-server which architectures its
//     build supports. An architecture missing from the list is therefore
//     reported as unknown and refused instead of being guessed as chat
//     capable. Adding a family is a one-line change once its chat ability has
//     been verified against the shipped runtime; the runtime load remains the
//     final authority for admitted models.
//   - Registry labels are ignored on purpose: a registry entry cannot promote
//     a component into a chat model.
//   - Non-GGUF files (legacy .bin/GGML) carry no GGUF metadata to classify.
//     They are denied for chat by default. A .bin extension is not evidence of
//     chat capability; any future legacy support needs an explicit validated
//     allowlist rather than the old blanket admission.

import (
	"fmt"
	"path/filepath"
	"strings"

	gguf_parser "github.com/gpustack/gguf-parser-go"
)

// AdmissionClass is the role a GGUF file can play in the application.
type AdmissionClass string

const (
	// AdmissionChat is a text chat architecture llama-server can serve.
	AdmissionChat AdmissionClass = "chat"
	// AdmissionMultimodal is a chat architecture that only makes sense with a
	// companion projector attached.
	AdmissionMultimodal AdmissionClass = "multimodal-chat"
	// AdmissionComponent is an asset that is part of another model: a
	// projector, text encoder, vocoder or speech encoder. It is never a chat
	// LLM, even if it sits in a chat model directory.
	AdmissionComponent AdmissionClass = "component"
	// AdmissionUnsupported is a recognised but non-chat architecture, such as
	// a diffusion/image model.
	AdmissionUnsupported AdmissionClass = "unsupported"
	// AdmissionUnknown is a file whose architecture is missing, unreadable or
	// not in the verified chat set. Unknown is never treated as chat.
	AdmissionUnknown AdmissionClass = "unknown"
	// AdmissionNonGGUF is a legacy format that has no GGUF metadata to read.
	// It is denied for chat by default; see the limitation note above.
	AdmissionNonGGUF AdmissionClass = "non-gguf"
)

// ProjectorRequirement states what a model needs beside it to be usable.
type ProjectorRequirement string

const (
	// ProjectorNotRequired is a text-only chat architecture.
	ProjectorNotRequired ProjectorRequirement = ""
	// ProjectorAny is a multimodal architecture that needs a projector but
	// does not restrict it to one encoder kind (Qwen VL families accept
	// vision, audio or hybrid projectors depending on the converted model).
	ProjectorAny ProjectorRequirement = "any"
	// ProjectorVisionRequired needs a projector with a vision encoder.
	ProjectorVisionRequired ProjectorRequirement = "vision"
)

// multimodalArchs are chat architectures that are only usable with a
// projector. The value is the minimum requirement, not a filename hint.
var multimodalArchs = map[string]ProjectorRequirement{
	// Qwen2-VL / Qwen2.5-Omni and Qwen3-VL / Qwen3-ASR. The same architecture
	// string covers vision and audio conversions, so any declared encoder is
	// acceptable; the exact companion is still chosen by identity.
	"qwen2vl": ProjectorAny,
	"qwen3vl": ProjectorAny,
	// These carry a vision tower in every published conversion.
	"mllama": ProjectorVisionRequired,
	"llava":  ProjectorVisionRequired,
}

// componentArchs are GGUF architectures that exist only as part of a larger
// model. They must never be offered as chat LLMs.
var componentArchs = map[string]bool{
	"clip": true, "mtmd": true,
	"t5": true, "t5encoder": true, "t5-decoder": true, "umt5": true,
	"bert": true, "roberta": true, "nomic-bert": true, "jina-bert-v2": true,
	"whisper": true, "wav2vec2": true, "wavtokenizer": true, "wavlm": true,
	"encodec": true, "dac": true, "bigvgan": true, "vocoder": true,
	"mimi": true, "speecht5": true, "nemo-conformer": true,
}

// Admission is the classification of one file, plus the evidence used.
type Admission struct {
	Path  string
	Class AdmissionClass
	Arch  string
	Name  string
	// EmbeddingLength is the model's text hidden size, used with the
	// projector's projection dimension to report compatibility.
	EmbeddingLength int
	// Projector records what the model needs beside it when loaded.
	Projector ProjectorRequirement
	// Tags are the model's declared general.tags. They are evidence for a
	// separate task policy; this classifier never promotes a tag to chat
	// capability on its own.
	Tags []string
	// Evidence lists the metadata keys that decided the classification.
	Evidence []string
	// Reason explains a refusal in actionable terms.
	Reason string
}

// AllowedAsChat reports whether the file's format and architecture may be
// loaded by the chat loader. Task eligibility is decided separately by the
// chatpolicy package; this method only answers the admission question.
func (a Admission) AllowedAsChat() bool {
	return a.Class == AdmissionChat || a.Class == AdmissionMultimodal
}

// RequiredModality maps the admission requirement to a projector kind filter.
func (a Admission) RequiredModality() ProjectorKind {
	switch a.Projector {
	case ProjectorVisionRequired:
		return ProjectorVision
	default:
		return ProjectorUnknown
	}
}

// RequiresProjector reports whether loading without a projector is invalid.
func (a Admission) RequiresProjector() bool {
	return a.Projector != ProjectorNotRequired
}

// ChatRejectionMessage returns the user-facing refusal for a non-chat file.
func (a Admission) ChatRejectionMessage() string {
	base := filepath.Base(a.Path)
	arch := a.Arch
	if arch == "" {
		arch = "(missing)"
	}
	switch a.Class {
	case AdmissionComponent:
		if isProjectorArch(arch) {
			return fmt.Sprintf("%s is a multimodal projector/encoder asset (GGUF architecture %q), not a chat model. "+
				"A projector is loaded automatically as the companion of the model it belongs to and cannot be selected as an LLM.",
				base, arch)
		}
		return fmt.Sprintf("%s cannot be used as a chat model: its GGUF architecture %q identifies a component-only asset, not a chat LLM. %s",
			base, arch, a.Reason)
	case AdmissionUnsupported:
		return fmt.Sprintf("%s cannot be used as a chat model: GGUF architecture %q is an unsupported non-chat architecture. %s",
			base, arch, a.Reason)
	case AdmissionUnknown:
		return fmt.Sprintf("%s cannot be used as a chat model: its GGUF architecture %q is not a verified chat architecture, so it was refused instead of risking a crash. %s",
			base, arch, a.Reason)
	case AdmissionNonGGUF:
		return fmt.Sprintf("%s cannot be used as a chat model: legacy .bin/.ggml files are denied for chat by default. "+
			"A .bin extension alone does not establish chat capability, and no explicit validated legacy allowlist matched.", base)
	default:
		return fmt.Sprintf("%s cannot be used as a chat model (classification %q).", base, a.Class)
	}
}

// AdmitModel classifies a file from its own GGUF metadata. Registry labels and
// file names are never consulted.
func AdmitModel(path string) Admission {
	if !strings.EqualFold(filepath.Ext(path), ".gguf") {
		return Admission{
			Path:   path,
			Class:  AdmissionNonGGUF,
			Reason: "not a GGUF file; GGUF architecture admission does not apply",
		}
	}

	parsed, err := gguf_parser.ParseGGUFFile(path)
	if err != nil {
		return Admission{
			Path:   path,
			Class:  AdmissionUnknown,
			Reason: fmt.Sprintf("GGUF metadata is unreadable (%v), so the file cannot be admitted as a chat model", err),
		}
	}
	return ClassifyModelMetadata(path, metadataKeys(parsed))
}

// ClassifyModelMetadata classifies a model from its GGUF metadata. It is
// separated from file IO so every rule is testable without a real file.
func ClassifyModelMetadata(path string, keys map[string]any) Admission {
	arch := strings.TrimSpace(metadataString(keys, "general.architecture"))
	name := strings.TrimSpace(metadataString(keys, "general.name"))

	admission := Admission{
		Path: path,
		Arch: arch,
		Name: name,
		Tags: metadataStrings(keys, "general.tags"),
	}
	if arch == "" {
		admission.Class = AdmissionUnknown
		admission.Reason = "the file declares no general.architecture"
		return admission
	}

	lower := strings.ToLower(arch)
	admission.Evidence = append(admission.Evidence, "general.architecture="+arch)

	// The LFM2.5 vocoder literally declares a refusal sentence as its
	// architecture. Treat any such self-description as a component.
	if strings.Contains(lower, "cannot be used as llm") ||
		strings.Contains(lower, "not an llm") ||
		strings.Contains(lower, "vocoder") {
		admission.Class = AdmissionComponent
		admission.Reason = "the file declares that it cannot be used as an LLM and must be loaded through the audio/TTS pipeline instead"
		return admission
	}

	if componentArchs[lower] {
		admission.Class = AdmissionComponent
		admission.Reason = "this architecture is an encoder/projector component (for example an mmproj whose general.architecture is clip), not a chat model"
		return admission
	}

	if diffusionArchs[lower] {
		admission.Class = AdmissionUnsupported
		admission.Reason = "diffusion and image architectures are served by the image engine, not by the chat loader"
		return admission
	}

	if !llmArchs[lower] {
		admission.Class = AdmissionUnknown
		admission.Reason = "add the architecture to the verified chat set only after confirming the shipped llama-server build supports it"
		return admission
	}

	admission.EmbeddingLength = metadataInt(keys, lower+".embedding_length")
	if admission.EmbeddingLength > 0 {
		admission.Evidence = append(admission.Evidence, fmt.Sprintf("%s.embedding_length=%d", lower, admission.EmbeddingLength))
	}

	if requirement, ok := multimodalArchs[lower]; ok {
		admission.Class = AdmissionMultimodal
		admission.Projector = requirement
		admission.Reason = "the architecture requires a companion projector"
		return admission
	}

	admission.Class = AdmissionChat
	return admission
}

// isProjectorArch reports whether an architecture names an encoder container.
func isProjectorArch(arch string) bool {
	lower := strings.ToLower(strings.TrimSpace(arch))
	return lower == "clip" || lower == "mtmd"
}

// metadataKeys converts the parser's key list into a map.
func metadataKeys(parsed *gguf_parser.GGUFFile) map[string]any {
	keys := make(map[string]any, len(parsed.Header.MetadataKV))
	for _, kv := range parsed.Header.MetadataKV {
		keys[kv.Key] = kv.Value
	}
	return keys
}

func metadataString(keys map[string]any, want string) string {
	for key, value := range keys {
		if !strings.EqualFold(key, want) {
			continue
		}
		if s, ok := value.(string); ok {
			return s
		}
	}
	return ""
}

// metadataStrings reads a string-array metadata value. The parser represents
// arrays as GGUFMetadataKVArrayValue rather than []string, so both forms are
// handled; the tag order is preserved.
func metadataStrings(keys map[string]any, want string) []string {
	for key, value := range keys {
		if !strings.EqualFold(key, want) {
			continue
		}
		switch v := value.(type) {
		case []string:
			return append([]string(nil), v...)
		case []any:
			out := make([]string, 0, len(v))
			for _, item := range v {
				if s, ok := item.(string); ok {
					out = append(out, s)
				}
			}
			return out
		case gguf_parser.GGUFMetadataKVArrayValue:
			out := make([]string, 0, len(v.Array))
			for _, item := range v.Array {
				if s, ok := item.(string); ok {
					out = append(out, s)
				}
			}
			return out
		}
	}
	return nil
}

func metadataInt(keys map[string]any, want string) int {
	for key, value := range keys {
		if !strings.EqualFold(key, want) {
			continue
		}
		switch v := value.(type) {
		case int:
			return v
		case int8:
			return int(v)
		case int16:
			return int(v)
		case int32:
			return int(v)
		case int64:
			return int(v)
		case uint:
			return int(v)
		case uint8:
			return int(v)
		case uint16:
			return int(v)
		case uint32:
			return int(v)
		case uint64:
			return int(v)
		}
	}
	return 0
}
