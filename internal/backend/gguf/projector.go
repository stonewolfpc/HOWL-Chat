package gguf

import (
	"fmt"
	"sort"
	"strings"

	gguf_parser "github.com/gpustack/gguf-parser-go"
)

// ProjectorKind is the modality a multimodal projector declares in its own
// metadata.
//
// It is never derived from the file name. "mmproj" names a file family, not a
// capability: an audio projector and a vision projector are both called
// mmproj-*.gguf. The llama.cpp tooling makes the same distinction
// (mtmd_gen_audio_get_info / mtmd_cap.audio_gen vs the vision encoder flags),
// and the GGUF file records it in clip.has_vision_encoder / clip.has_audio_encoder
// plus the clip.vision.* / clip.audio.* key spaces.
type ProjectorKind string

const (
	// ProjectorUnknown means the file did not declare an encoder either way.
	// Unknown is not vision: callers must fail closed.
	ProjectorUnknown ProjectorKind = "unknown"
	// ProjectorVision declares a vision encoder only.
	ProjectorVision ProjectorKind = "vision"
	// ProjectorAudio declares an audio encoder only.
	ProjectorAudio ProjectorKind = "audio"
	// ProjectorHybrid declares both a vision and an audio encoder.
	ProjectorHybrid ProjectorKind = "hybrid"
)

// ProjectorInfo is the classification of one projector file.
type ProjectorInfo struct {
	Path string
	Kind ProjectorKind
	// Evidence lists the metadata keys that decided the classification.
	Evidence []string
	// Err records why classification could not be completed. Kind stays
	// ProjectorUnknown when Err is set.
	Err string
	// ProjectionDim is clip.<modality>.projection_dim when the file declares
	// it. It is diagnostic evidence for compatibility, never a filename guess.
	ProjectionDim int
	// ProjectorType is clip.projector_type when the file declares it.
	ProjectorType string
}

// HasVision reports whether the projector declares a vision encoder.
func (p ProjectorInfo) HasVision() bool {
	return p.Kind == ProjectorVision || p.Kind == ProjectorHybrid
}

// HasAudio reports whether the projector declares an audio encoder.
func (p ProjectorInfo) HasAudio() bool {
	return p.Kind == ProjectorAudio || p.Kind == ProjectorHybrid
}

// InspectProjector reads a projector's GGUF metadata and classifies it. A file
// that cannot be read is reported as unknown with Err set; it is never assumed
// to be a vision projector.
func InspectProjector(path string) ProjectorInfo {
	if strings.TrimSpace(path) == "" {
		return ProjectorInfo{Kind: ProjectorUnknown, Err: "no projector path"}
	}

	parsed, err := gguf_parser.ParseGGUFFile(path)
	if err != nil {
		return ProjectorInfo{
			Path: path,
			Kind: ProjectorUnknown,
			Err:  fmt.Sprintf("projector metadata unreadable: %v", err),
		}
	}

	keys := make(map[string]any, len(parsed.Header.MetadataKV))
	for _, kv := range parsed.Header.MetadataKV {
		keys[kv.Key] = kv.Value
	}

	info := ClassifyProjectorMetadata(keys)
	info.Path = path
	return info
}

// ClassifyProjectorMetadata decides the projector kind from GGUF metadata keys.
// It is separated from file IO so the decision is testable without a GGUF file.
func ClassifyProjectorMetadata(keys map[string]any) ProjectorInfo {
	info := ProjectorInfo{Kind: ProjectorUnknown}
	if len(keys) == 0 {
		return info
	}

	vision, visionDeclared := explicitBool(keys, "clip.has_vision_encoder")
	audio, audioDeclared := explicitBool(keys, "clip.has_audio_encoder")

	visionKeys := prefixKeys(keys, "clip.vision.")
	audioKeys := prefixKeys(keys, "clip.audio.")

	// An explicit flag wins; the key space is the fallback for projectors whose
	// converter did not write the flag.
	if !visionDeclared {
		vision = len(visionKeys) > 0
	}
	if !audioDeclared {
		audio = len(audioKeys) > 0
	}

	switch {
	case vision && audio:
		info.Kind = ProjectorHybrid
	case vision:
		info.Kind = ProjectorVision
	case audio:
		info.Kind = ProjectorAudio
	default:
		info.Kind = ProjectorUnknown
	}

	evidence := make([]string, 0, 4)
	if visionDeclared {
		evidence = append(evidence, "clip.has_vision_encoder")
	}
	if audioDeclared {
		evidence = append(evidence, "clip.has_audio_encoder")
	}
	if vision && len(visionKeys) > 0 {
		evidence = append(evidence, visionKeys[0])
	}
	if audio && len(audioKeys) > 0 {
		evidence = append(evidence, audioKeys[0])
	}
	info.Evidence = evidence
	info.ProjectorType = metadataString(keys, "clip.projector_type")
	if audio {
		info.ProjectionDim = metadataInt(keys, "clip.audio.projection_dim")
	}
	if info.ProjectionDim == 0 && vision {
		info.ProjectionDim = metadataInt(keys, "clip.vision.projection_dim")
	}
	return info
}

// explicitBool reads a boolean metadata key, tolerating the numeric encodings a
// GGUF value can take. The second return reports whether the key was present.
func explicitBool(keys map[string]any, want string) (bool, bool) {
	for key, value := range keys {
		if !strings.EqualFold(key, want) {
			continue
		}
		switch v := value.(type) {
		case bool:
			return v, true
		case uint8:
			return v != 0, true
		case uint16:
			return v != 0, true
		case uint32:
			return v != 0, true
		case uint64:
			return v != 0, true
		case int8:
			return v != 0, true
		case int16:
			return v != 0, true
		case int32:
			return v != 0, true
		case int64:
			return v != 0, true
		case int:
			return v != 0, true
		case float32:
			return v != 0, true
		case float64:
			return v != 0, true
		}
		return false, true
	}
	return false, false
}

// prefixKeys returns the matching keys in sorted order so evidence is stable.
func prefixKeys(keys map[string]any, prefix string) []string {
	var out []string
	for key := range keys {
		if strings.HasPrefix(strings.ToLower(key), prefix) {
			out = append(out, key)
		}
	}
	sort.Strings(out)
	return out
}
