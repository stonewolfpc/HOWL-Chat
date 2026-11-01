package registry

import (
	"path/filepath"
	"strings"

	"howl-chat/internal/backend/gguf"
)

// imageAssetTypes are the asset types produced by the image pipeline.
var imageAssetTypes = map[string]bool{
	"diffusion": true, "vae": true, "lora": true, "controlnet": true,
	"embedding": true, "upscaler": true, "clip_l": true, "clip_g": true,
	"t5xxl": true,
}

// entryGGUFFile returns the .gguf file that represents an entry, if it has one.
// Bundle entries point at a directory, so the model slot is checked first. The
// separate llm slot is only consulted for entries that are themselves chat
// models: an image bundle that carries a Qwen text encoder would otherwise look
// like a chat model and be thrown away.
func entryGGUFFile(e Entry) string {
	candidates := []string{e.Files["model"], e.Path}
	if e.Type == KindLLM {
		candidates = append(candidates, e.Files["llm"])
	}
	for _, c := range candidates {
		if strings.EqualFold(filepath.Ext(c), ".gguf") {
			return c
		}
	}
	return ""
}

// entryKindMatches reports whether an entry's registered type agrees with the
// model's architecture. Only GGUF carries a readable architecture, so anything
// else is accepted unchanged.
//
// Used when the registry is loaded: a chat model that was indexed as an image
// asset (or the reverse) is dropped so the next scan can file it correctly.
func entryKindMatches(e Entry) bool {
	file := entryGGUFFile(e)
	if file == "" {
		return true
	}
	switch gguf.KindOf(file) {
	case gguf.KindLLM:
		return !imageAssetTypes[e.Type]
	case gguf.KindDiffusion:
		return e.Type != KindLLM
	default:
		return true
	}
}
