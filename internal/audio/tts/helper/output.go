package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"howl-chat/internal/audio/types"
)

// validateWAVFile validates the actual bytes, metadata, hash, and exact-path
// containment before a result may cross the acceptance barrier.
func validateWAVFile(w synthesisResult, dir, filename string) (*types.AudioMetadata, *Failure) {
	final := filepath.Join(filepath.Clean(dir), filename)
	info, err := os.Lstat(final)
	if err != nil {
		return nil, newFailure("output_failed", "synthesis output is missing")
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return nil, newFailure("output_failed", "synthesis output is not a regular file")
	}
	reparse, reparseErr := hasReparseComponent(final)
	if reparseErr != nil || reparse {
		return nil, newFailure("output_failed", "synthesis output path contains a reparse point")
	}
	rel, err := filepath.Rel(filepath.Clean(dir), final)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return nil, newFailure("output_failed", "synthesis output escapes the request directory")
	}
	raw, err := os.ReadFile(final)
	if err != nil {
		return nil, newFailure("output_failed", "synthesis output is unreadable")
	}
	if int64(len(raw)) != w.Bytes {
		return nil, newFailure("output_failed", "synthesis byte count does not match the file")
	}
	if len(raw) < 44 || string(raw[0:4]) != "RIFF" || string(raw[8:12]) != "WAVE" ||
		string(raw[12:16]) != "fmt " || string(raw[36:40]) != "data" {
		return nil, newFailure("output_failed", "synthesis output is not canonical RIFF/WAVE")
	}
	if little16(raw[20:22]) != 1 || little16(raw[22:24]) != 1 || little16(raw[34:36]) != 16 {
		return nil, newFailure("output_failed", "synthesis output is not PCM16 mono")
	}
	rate := int(little32(raw[24:28]))
	dataBytes := int(little32(raw[40:44]))
	if rate < 8000 || rate > 48000 || little32(raw[28:32]) != uint32(rate*2) ||
		little16(raw[32:34]) != 2 || int(little32(raw[16:20])) != 16 {
		return nil, newFailure("output_failed", "synthesis WAV format fields are invalid")
	}
	if dataBytes <= 0 || dataBytes%2 != 0 || len(raw) != 44+dataBytes {
		return nil, newFailure("output_failed", "synthesis WAV data length is invalid")
	}
	duration := float64(dataBytes) / float64(rate*2)
	if abs(duration-w.DurationSeconds) > 0.001 {
		return nil, newFailure("output_failed", "synthesis duration does not match PCM frames")
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != w.SHA256 {
		return nil, newFailure("output_failed", "synthesis hash does not match the file")
	}
	parts, err := filepath.Glob(final + ".part-*")
	if err != nil {
		return nil, newFailure("output_failed", "could not check temporary output")
	}
	if len(parts) != 0 {
		return nil, newFailure("output_failed", "synthesis left temporary output")
	}
	return &types.AudioMetadata{
		Format: types.FormatWAV, SampleRate: rate, Channels: 1, BitDepth: 16,
		Duration: duration, BitRate: rate * 16 / 1000,
	}, nil
}

func requestOutputExists(dir, filename string) bool {
	if dir == "" || filename == "" || filepath.Base(filename) != filename || strings.Contains(filename, "..") {
		return false
	}
	final := filepath.Join(filepath.Clean(dir), filename)
	if _, err := os.Lstat(final); err == nil {
		return true
	}
	parts, err := filepath.Glob(final + ".part-*")
	return err == nil && len(parts) > 0
}

func cleanupExact(dir, filename string) {
	if dir == "" || filename == "" {
		return
	}
	if filepath.Base(filename) != filename || strings.Contains(filename, "..") {
		return
	}
	if reparse, err := hasReparseComponent(filepath.Clean(dir)); err != nil || reparse {
		return
	}
	final := filepath.Join(filepath.Clean(dir), filename)
	_ = os.Remove(final)
	parts, err := filepath.Glob(final + ".part-*")
	if err != nil {
		return
	}
	for _, part := range parts {
		if filepath.Dir(part) == filepath.Clean(dir) {
			_ = os.Remove(part)
		}
	}
}

func little16(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}

func little32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
