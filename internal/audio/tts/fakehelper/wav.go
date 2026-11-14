package fakehelper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

const (
	syntheticRate    = 24000
	syntheticSamples = 1200
)

// writeSyntheticWAV writes deterministic, clearly synthetic PCM16 mono audio.
// It is deliberately not speech and has no external asset dependency.
func writeSyntheticWAV(path string) (string, int64, float64, error) {
	header, data, err := syntheticWAVBytes(syntheticRate, syntheticSamples)
	if err != nil {
		return "", 0, 0, err
	}
	raw := append(header, data...)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return "", 0, 0, err
	}
	writeErr := writeAll(f, raw)
	closeErr := f.Close()
	if writeErr != nil {
		_ = os.Remove(path)
		return "", 0, 0, writeErr
	}
	if closeErr != nil {
		_ = os.Remove(path)
		return "", 0, 0, closeErr
	}
	sum := sha256.Sum256(raw)
	duration := float64(len(data)) / float64(syntheticRate*2)
	return hex.EncodeToString(sum[:]), int64(len(raw)), duration, nil
}

func syntheticWAVBytes(rate, samples int) ([]byte, []byte, error) {
	if rate < 8000 || rate > 48000 || samples <= 0 {
		return nil, nil, fmt.Errorf("invalid synthetic WAV parameters")
	}
	dataBytes := samples * 2
	header := make([]byte, 44)
	copy(header[0:4], "RIFF")
	put32(header[4:8], uint32(36+dataBytes))
	copy(header[8:12], "WAVE")
	copy(header[12:16], "fmt ")
	put32(header[16:20], 16)
	put16(header[20:22], 1)
	put16(header[22:24], 1)
	put32(header[24:28], uint32(rate))
	put32(header[28:32], uint32(rate*2))
	put16(header[32:34], 2)
	put16(header[34:36], 16)
	copy(header[36:40], "data")
	put32(header[40:44], uint32(dataBytes))

	data := make([]byte, dataBytes)
	for i := 0; i < samples; i++ {
		// A repeating low-amplitude ramp is unambiguously synthetic and has no
		// voice, recording, or external asset content.
		sample := int16((i*37)%1200 - 600)
		put16(data[i*2:i*2+2], uint16(sample))
	}
	return header, data, nil
}

func validateSyntheticWAV(path string, wantHash string, wantSize int64) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if int64(len(raw)) != wantSize || len(raw) < 44 {
		return fmt.Errorf("WAV size mismatch")
	}
	if string(raw[0:4]) != "RIFF" || string(raw[8:12]) != "WAVE" ||
		string(raw[12:16]) != "fmt " || string(raw[36:40]) != "data" {
		return fmt.Errorf("WAV header is not canonical")
	}
	if little16(raw[20:22]) != 1 || little16(raw[22:24]) != 1 || little16(raw[34:36]) != 16 {
		return fmt.Errorf("WAV is not PCM16 mono")
	}
	rate := int(little32(raw[24:28]))
	if rate != syntheticRate || int(little32(raw[16:20])) != 16 {
		return fmt.Errorf("WAV format fields are invalid")
	}
	dataBytes := int(little32(raw[40:44]))
	if dataBytes <= 0 || dataBytes%2 != 0 || len(raw) != 44+dataBytes {
		return fmt.Errorf("WAV data length is invalid")
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != wantHash {
		return fmt.Errorf("WAV hash mismatch")
	}
	return nil
}

func little16(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}

func little32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func put16(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func put32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func writeAll(f *os.File, raw []byte) error {
	for len(raw) > 0 {
		n, err := f.Write(raw)
		if err != nil {
			return err
		}
		raw = raw[n:]
	}
	return nil
}
