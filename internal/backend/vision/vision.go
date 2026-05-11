/**
 * Vision / Multimodal Support
 *
 * This package provides CLIP-based vision capabilities for multimodal models.
 * It handles image preprocessing, CLIP encoding, and vision embedding injection.
 *
 * @package vision
 */

package vision

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strconv"

	"howl-chat/internal/backend/types"

	"github.com/disintegration/imaging"
	gguf_parser "github.com/gpustack/gguf-parser-go"
)

// BuildVisionProfile extracts CLIP metadata from a GGUF file.
func BuildVisionProfile(path string) (*types.VisionProfile, error) {
	gguf, err := gguf_parser.ParseGGUFFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GGUF: %w", err)
	}

	meta := gguf.Header.MetadataKV

	vp := &types.VisionProfile{
		ClipModelPath: path,
		ImageSize:     getInt(meta, "clip.image_size", 336),
		PatchSize:     getInt(meta, "clip.patch_size", 14),
		EmbeddingDim:  getInt(meta, "clip.projection_dim", 1024),
		ContextLength: getInt(meta, "clip.context_length", 77),
		ModelType:     getString(meta, "clip.model_type", "clip"),
	}

	return vp, nil
}

// getInt extracts an integer value from GGUF metadata with a default.
func getInt(meta gguf_parser.GGUFMetadataKVs, key string, defaultValue int) int {
	for _, kv := range meta {
		if kv.Key == key {
			switch t := kv.Value.(type) {
			case int:
				return t
			case int64:
				return int(t)
			case float64:
				return int(t)
			case string:
				if i, err := strconv.Atoi(t); err == nil {
					return i
				}
			}
		}
	}
	return defaultValue
}

// getString extracts a string value from GGUF metadata with a default.
func getString(meta gguf_parser.GGUFMetadataKVs, key string, defaultValue string) string {
	for _, kv := range meta {
		if kv.Key == key {
			if s, ok := kv.Value.(string); ok {
				return s
			}
		}
	}
	return defaultValue
}

// PreprocessImage resizes and normalizes an image for CLIP encoding.
func PreprocessImage(img image.Image, size int, normalize bool) ([]float32, error) {
	// Resize image to target size
	resized := imaging.Resize(img, size, size, imaging.Lanczos)

	bounds := resized.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Convert to float32 tensor (RGB)
	tensor := make([]float32, w*h*3)

	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()

			if normalize {
				// Normalize to [0, 1] range
				tensor[i+0] = float32(r) / 65535.0
				tensor[i+1] = float32(g) / 65535.0
				tensor[i+2] = float32(b) / 65535.0
			} else {
				// Use raw pixel values
				tensor[i+0] = float32(r >> 8)
				tensor[i+1] = float32(g >> 8)
				tensor[i+2] = float32(b >> 8)
			}
			i += 3
		}
	}

	return tensor, nil
}

// BuildVisionPrompt adds vision tags to the prompt template for multimodal models.
func BuildVisionPrompt(template string, hasImage bool) string {
	if hasImage {
		// Replace {{image}} placeholder with <image> tag
		template = replacePlaceholder(template, "{{image}}", "<image>")
	} else {
		// Remove image placeholder
		template = replacePlaceholder(template, "{{image}}", "")
	}

	return template
}

// replacePlaceholder replaces a placeholder in a template string.
func replacePlaceholder(template, placeholder, replacement string) string {
	// Simple string replacement
	result := template
	for {
		idx := findSubstring(result, placeholder)
		if idx == -1 {
			break
		}
		result = result[:idx] + replacement + result[idx+len(placeholder):]
	}
	return result
}

// findSubstring finds a substring (case-insensitive).
func findSubstring(s, substr string) int {
	// Simple case-insensitive search
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(s[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// toLower converts a byte to lowercase.
func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}
