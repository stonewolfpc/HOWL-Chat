package types

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateAudioFileRejectsWebM(t *testing.T) {
	path := filepath.Join(t.TempDir(), "recording.webm")
	data := []byte{0x1A, 0x45, 0xDF, 0xA3, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write WebM fixture: %v", err)
	}

	err := ValidateAudioFile(path)
	if err == nil {
		t.Fatal("ValidateAudioFile accepted WebM, want rejection")
	}
	var audioErr *AudioError
	if !errors.As(err, &audioErr) || audioErr.Code != ErrCodeFormatNotSupported {
		t.Fatalf("error = %v, want ErrCodeFormatNotSupported", err)
	}
}

func TestValidateAudioFileAcceptsWAV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "recording.wav")
	data := make([]byte, 44)
	copy(data[0:4], "RIFF")
	binary.LittleEndian.PutUint32(data[4:8], 36)
	copy(data[8:12], "WAVE")
	copy(data[12:16], "fmt ")
	binary.LittleEndian.PutUint32(data[16:20], 16)
	binary.LittleEndian.PutUint16(data[20:22], 1)
	binary.LittleEndian.PutUint16(data[22:24], 1)
	binary.LittleEndian.PutUint32(data[24:28], 16000)
	binary.LittleEndian.PutUint32(data[28:32], 32000)
	binary.LittleEndian.PutUint16(data[32:34], 2)
	binary.LittleEndian.PutUint16(data[34:36], 16)
	copy(data[36:40], "data")
	binary.LittleEndian.PutUint32(data[40:44], 0)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write WAV fixture: %v", err)
	}

	if err := ValidateAudioFile(path); err != nil {
		t.Fatalf("ValidateAudioFile rejected valid WAV: %v", err)
	}
}
