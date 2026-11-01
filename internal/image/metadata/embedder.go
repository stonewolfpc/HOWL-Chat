package metadata

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"image"
	"image/png"
	"strconv"
	"strings"
)

// EmbedMetadata embeds A1111-style parameters into a PNG image.
// The parameters are stored in a tEXt chunk with the keyword "parameters".
// Returns the PNG data with embedded metadata.
func EmbedMetadata(imageData []byte, params MetadataParams) ([]byte, error) {
	// Check if it's a PNG
	if len(imageData) < 8 || !bytes.HasPrefix(imageData, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return nil, fmt.Errorf("not a PNG image")
	}

	// Format the parameters string
	paramsStr := FormatParameters(params)

	// Create the tEXt chunk
	textChunk := createTEXtChunk("parameters", paramsStr)

	// Insert the tEXt chunk after the IHDR chunk
	return insertTEXtChunk(imageData, textChunk)
}

// MetadataParams contains the generation parameters to embed.
type MetadataParams struct {
	Prompt         string
	NegativePrompt string
	Steps          int
	Sampler        string
	CFGScale       float64
	Seed           int64
	Width          int
	Height         int
	Scheduler      string
}

// FormatParameters formats metadata into A1111-style parameter string.
// Format: "prompt\nNegative prompt: negative\nSteps: 20, Sampler: Euler a, CFG scale: 7, Seed: 123, Size: 512x512"
func FormatParameters(p MetadataParams) string {
	var sb strings.Builder

	// First line: positive prompt
	sb.WriteString(p.Prompt)
	sb.WriteString("\n")

	// Second line: negative prompt (if present)
	if p.NegativePrompt != "" {
		sb.WriteString("Negative prompt: ")
		sb.WriteString(p.NegativePrompt)
		sb.WriteString("\n")
	}

	// Third line: parameters
	sb.WriteString("Steps: ")
	sb.WriteString(strconv.Itoa(p.Steps))
	sb.WriteString(", Sampler: ")
	sb.WriteString(p.Sampler)
	sb.WriteString(", CFG scale: ")
	sb.WriteString(strconv.FormatFloat(p.CFGScale, 'f', 1, 64))
	sb.WriteString(", Seed: ")
	sb.WriteString(strconv.FormatInt(p.Seed, 10))
	sb.WriteString(", Size: ")
	sb.WriteString(strconv.Itoa(p.Width))
	sb.WriteString("x")
	sb.WriteString(strconv.Itoa(p.Height))

	if p.Scheduler != "" {
		sb.WriteString(", Scheduler: ")
		sb.WriteString(p.Scheduler)
	}

	return sb.String()
}

// createTEXtChunk creates a PNG tEXt chunk with the given keyword and text.
// tEXt chunk format: length (4 bytes) | type "tEXt" (4 bytes) | data | CRC (4 bytes)
// Data format: keyword (1-79 chars) + null byte + text
func createTEXtChunk(keyword, text string) []byte {
	// tEXt data: keyword + null + text
	data := []byte(keyword)
	data = append(data, 0)
	data = append(data, []byte(text)...)

	// Chunk: length (4 bytes, big-endian) + type (4 bytes) + data + CRC (4 bytes)
	length := uint32(len(data))
	chunk := make([]byte, 4+4+len(data)+4)

	// Write length
	binary.BigEndian.PutUint32(chunk[0:4], length)

	// Write type
	copy(chunk[4:8], []byte("tEXt"))

	// Write data
	copy(chunk[8:], data)

	// Calculate CRC over type + data
	crc := crc32.ChecksumIEEE(chunk[4 : 8+len(data)])
	binary.BigEndian.PutUint32(chunk[8+len(data):], crc)

	return chunk
}

// insertTEXtChunk inserts the tEXt chunk after the IHDR chunk in a PNG.
func insertTEXtChunk(imageData, textChunk []byte) ([]byte, error) {
	// PNG structure: signature (8) | IHDR chunk | ... | IEND chunk
	// We want to insert tEXt after IHDR

	result := make([]byte, 0, len(imageData)+len(textChunk))

	// Copy PNG signature
	if len(imageData) < 8 {
		return nil, fmt.Errorf("invalid PNG: too short")
	}
	result = append(result, imageData[:8]...)
	offset := 8

	// Read and copy chunks until we find IHDR
	for offset < len(imageData) {
		if offset+8 > len(imageData) {
			break
		}

		// Read chunk length (4 bytes)
		length := binary.BigEndian.Uint32(imageData[offset : offset+4])

		if offset+8+int(length)+4 > len(imageData) {
			break
		}

		// Read chunk type (4 bytes)
		chunkType := string(imageData[offset+4 : offset+8])

		// Copy this chunk
		chunkEnd := offset + 8 + int(length) + 4 // length + type + data + CRC
		result = append(result, imageData[offset:chunkEnd]...)
		offset = chunkEnd

		// After IHDR, insert our tEXt chunk
		if chunkType == "IHDR" {
			result = append(result, textChunk...)
		}

		// Stop at IEND
		if chunkType == "IEND" {
			break
		}

		// Safety check
		if length > 100_000_000 { // 100MB chunk limit
			break
		}
	}

	return result, nil
}

// EmbedMetadataFromBase64 embeds metadata into a base64-encoded PNG.
func EmbedMetadataFromBase64(imageB64 string, params MetadataParams) (string, error) {
	imageData, err := base64.StdEncoding.DecodeString(imageB64)
	if err != nil {
		return "", fmt.Errorf("decode base64: %w", err)
	}

	result, err := EmbedMetadata(imageData, params)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(result), nil
}

// EncodePNGWithMetadata encodes an image.Image to PNG with embedded metadata.
// This is used when we have an image.Image object rather than raw PNG bytes.
func EncodePNGWithMetadata(img image.Image, params MetadataParams) ([]byte, error) {
	// First encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encode PNG: %w", err)
	}

	// Then embed metadata
	return EmbedMetadata(buf.Bytes(), params)
}
