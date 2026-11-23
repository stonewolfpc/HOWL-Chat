package scanner

import (
	"fmt"
	"strings"

	gguf_parser "github.com/gpustack/gguf-parser-go"
)

// GGUF image weights (Flux, SD3, Z-Image, Wan) carry the same header chat models
// do: a declared architecture plus the tensor list. Reading both is what lets a
// quantised image model be classified from its contents instead of its filename,
// the same bar safetensors already gets. The filename is still the fallback when
// the header has nothing to say.
func (p *WeightProbe) inspectGGUF(path string) {
	// The parser allocates from lengths read out of the file, so a corrupt or
	// non-GGUF file with a .gguf name panics on the string length it finds. One
	// such file in a library used to abort the whole scan (the Add Folder call
	// never returned and the UI sat on "Scanning…"). ProbeWeight promises an
	// empty probe instead, so a bad header yields no probe and filename rules
	// stay in charge for that file.
	defer func() {
		if r := recover(); r != nil {
			*p = WeightProbe{Metadata: map[string]string{}}
			fmt.Printf("[Scanner] skipping malformed GGUF header %s: %v\n", path, r)
		}
	}()

	f, err := gguf_parser.ParseGGUFFile(path)
	if err != nil || f == nil {
		return
	}

	for _, kv := range f.Header.MetadataKV {
		s, ok := kv.Value.(string)
		if !ok || strings.TrimSpace(s) == "" {
			continue
		}
		switch kv.Key {
		case "general.architecture", "general.name", "general.basename", "general.size_label":
			p.Metadata[kv.Key] = s
		}
	}

	// Tensor names carry the architecture signature (Z-Image's noise_refiner,
	// Flux's double_blocks) and the asset type (VAE / LoRA / ControlNet), so the
	// same structural rules apply to every format.
	for _, t := range f.TensorInfos {
		p.TensorCount++
		if p.TensorCount > maxProbeKeys {
			return
		}
		p.noteKey(t.Name)
	}
}
