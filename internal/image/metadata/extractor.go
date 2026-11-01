package metadata

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"image"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ExtractedMetadata holds generation parameters extracted from an image.
type ExtractedMetadata struct {
	Prompt         string  `json:"prompt"`
	NegativePrompt string  `json:"negative_prompt"`
	Steps          int     `json:"steps"`
	Sampler        string  `json:"sampler"`
	CFGScale       float64 `json:"cfg_scale"`
	Seed           int64   `json:"seed"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	ModelHash      string  `json:"model_hash,omitempty"`
	ModelName      string  `json:"model_name,omitempty"`
	Scheduler      string  `json:"scheduler,omitempty"`
	// Raw contains the unparsed parameters string for debugging
	Raw string `json:"raw,omitempty"`
}

// ExtractFromFile extracts metadata from an image file (PNG or JPEG).
func ExtractFromFile(path string) (*ExtractedMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return ExtractFromBytes(data)
}

// ExtractFromBase64 extracts metadata from base64-encoded image data.
func ExtractFromBase64(b64 string) (*ExtractedMetadata, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}
	return ExtractFromBytes(data)
}

// ExtractFromBytes extracts metadata from raw image bytes.
func ExtractFromBytes(data []byte) (*ExtractedMetadata, error) {
	// Detect format by magic bytes
	if len(data) < 8 {
		return nil, fmt.Errorf("file too small")
	}

	// PNG magic: 89 50 4E 47 0D 0A 1A 0A
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return extractFromPNG(data)
	}

	// JPEG magic: FF D8 FF
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return extractFromJPEG(data)
	}

	return nil, fmt.Errorf("unsupported image format")
}

// extractFromPNG extracts metadata from PNG tEXt chunks.
func extractFromPNG(data []byte) (*ExtractedMetadata, error) {
	// Look for tEXt chunks manually since image/png doesn't expose them easily
	const tEXt = "tEXt"
	const parametersKey = "parameters"

	// Skip PNG signature (8 bytes)
	offset := 8
	var params string

	for offset < len(data) {
		if offset+8 > len(data) {
			break
		}

		// Read chunk length (4 bytes, big-endian)
		length := binary.BigEndian.Uint32(data[offset : offset+4])
		offset += 4

		if offset+4 > len(data) {
			break
		}

		// Read chunk type (4 bytes)
		chunkType := string(data[offset : offset+4])
		offset += 4

		if offset+int(length) > len(data) {
			break
		}

		// Read chunk data
		chunkData := data[offset : offset+int(length)]
		offset += int(length)

		// Skip CRC (4 bytes)
		offset += 4

		if chunkType == tEXt {
			// tEXt format: keyword\x00text
			// Keyword is 1-79 characters, null-terminated
			nullIdx := bytes.IndexByte(chunkData, 0)
			if nullIdx > 0 {
				keyword := string(chunkData[:nullIdx])
				text := string(chunkData[nullIdx+1:])
				if keyword == parametersKey {
					params = text
					break
				}
			}
		}

		// Stop at IEND
		if chunkType == "IEND" {
			break
		}

		// Safety check
		if length > 10_000_000 { // 10MB chunk size limit
			break
		}
	}

	if params == "" {
		return &ExtractedMetadata{}, nil // No metadata found, return empty
	}

	return parseParameters(params), nil
}

// extractFromJPEG attempts to extract metadata from JPEG EXIF.
// Note: Standard A1111 doesn't typically embed in JPEG EXIF, but we check for UserComment.
func extractFromJPEG(data []byte) (*ExtractedMetadata, error) {
	// For now, return empty - JPEG EXIF parsing is complex and A1111 uses PNG
	// TODO: Implement EXIF UserComment extraction if needed
	return &ExtractedMetadata{}, nil
}

// parseParameters parses the A1111 parameters format.
// Format: "prompt\nNegative prompt: negative\nSteps: 20, Sampler: Euler a, CFG scale: 7, Seed: 123, Size: 512x512"
func parseParameters(params string) *ExtractedMetadata {
	meta := &ExtractedMetadata{Raw: params}

	lines := strings.Split(params, "\n")
	if len(lines) == 0 {
		return meta
	}

	// First line is always the positive prompt
	meta.Prompt = strings.TrimSpace(lines[0])

	// Look for negative prompt line
	var paramLine string
	for i, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Negative prompt:") {
			meta.NegativePrompt = strings.TrimSpace(strings.TrimPrefix(line, "Negative prompt:"))
		} else if strings.Contains(line, "Steps:") || strings.Contains(line, "Sampler:") {
			// This is the parameters line
			paramLine = line
			// Also check remaining lines after this one
			if i+2 < len(lines) {
				paramLine += " " + strings.Join(lines[i+2:], " ")
			}
			break
		}
	}

	// If we didn't find a param line, the second line might be it
	if paramLine == "" && len(lines) > 1 && !strings.HasPrefix(lines[1], "Negative prompt:") {
		paramLine = strings.TrimSpace(lines[1])
	}

	// Parse the parameters line
	if paramLine != "" {
		parseParamLine(paramLine, meta)
	}

	return meta
}

// parseParamLine extracts individual parameters from the comma-separated line.
func parseParamLine(line string, meta *ExtractedMetadata) {
	// Split by comma, but be careful about prompts that might contain commas
	parts := splitParams(line)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Try to parse "Key: Value" format
		colonIdx := strings.Index(part, ":")
		if colonIdx < 0 {
			continue
		}

		key := strings.TrimSpace(strings.ToLower(part[:colonIdx]))
		value := strings.TrimSpace(part[colonIdx+1:])

		switch key {
		case "steps":
			if v, err := strconv.Atoi(value); err == nil {
				meta.Steps = v
			}
		case "sampler":
			meta.Sampler = value
		case "cfg scale", "cfg_scale", "cfg":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				meta.CFGScale = v
			}
		case "seed":
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				meta.Seed = v
			}
		case "size":
			// Format: 512x512 or 512 x 512
			parseSize(value, meta)
		case "model hash", "model_hash":
			meta.ModelHash = value
		case "model", "model name":
			meta.ModelName = value
		case "scheduler":
			meta.Scheduler = value
		}
	}
}

// parseSize parses size strings like "512x512" or "512 x 512".
func parseSize(sizeStr string, meta *ExtractedMetadata) {
	// Remove spaces
	sizeStr = strings.ReplaceAll(sizeStr, " ", "")

	// Match numbers around x
	re := regexp.MustCompile(`(\d+)\s*x\s*(\d+)`)
	matches := re.FindStringSubmatch(sizeStr)
	if len(matches) == 3 {
		if w, err := strconv.Atoi(matches[1]); err == nil {
			meta.Width = w
		}
		if h, err := strconv.Atoi(matches[2]); err == nil {
			meta.Height = h
		}
	}
}

// splitParams splits a parameter line by commas, being careful about nested structures.
func splitParams(line string) []string {
	var parts []string
	var current strings.Builder
	depth := 0

	for _, r := range line {
		switch r {
		case '(':
			depth++
			current.WriteRune(r)
		case ')':
			depth--
			current.WriteRune(r)
		case ',':
			if depth == 0 {
				parts = append(parts, current.String())
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// OpenImage opens an image for display (returns image.Image and format).
func OpenImage(path string) (image.Image, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}

// GetImageDimensions returns the width and height of an image file.
func GetImageDimensions(path string) (width, height int, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}

	// For JPEG/PNG, the dimensions are in the config
	if format == "jpeg" || format == "png" {
		return cfg.Width, cfg.Height, nil
	}

	// For other formats, we might need to decode fully
	return cfg.Width, cfg.Height, nil
}
