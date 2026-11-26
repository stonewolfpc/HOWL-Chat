package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// minimalWAV is enough for RIFF/WAVE detection without a real audio payload.
var minimalWAV = []byte("RIFF\x24\x00\x00\x00WAVEfmt \x10\x00\x00\x00")

func TestGetTTSAudioDataReturnsCacheWAV(t *testing.T) {
	cache := t.TempDir()
	t.Setenv("HOWL_AUDIO_CACHE_DIR", cache)

	path := filepath.Join(cache, "tts_1.wav")
	if err := os.WriteFile(path, minimalWAV, 0600); err != nil {
		t.Fatalf("write cache wav: %v", err)
	}

	app := &App{}
	got, err := app.GetTTSAudioData(path)
	if err != nil {
		t.Fatalf("GetTTSAudioData: %v", err)
	}

	const prefix = "data:audio/wav;base64,"
	if !strings.HasPrefix(got, prefix) {
		t.Fatalf("payload is not a wav data URL: %q", got)
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(got, prefix))
	if err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if string(decoded) != string(minimalWAV) {
		t.Fatalf("payload bytes changed: %d bytes", len(decoded))
	}
}

func TestGetTTSAudioDataRefusesPathOutsideCache(t *testing.T) {
	cache := t.TempDir()
	t.Setenv("HOWL_AUDIO_CACHE_DIR", cache)

	outside := filepath.Join(t.TempDir(), "private.wav")
	if err := os.WriteFile(outside, minimalWAV, 0600); err != nil {
		t.Fatalf("write outside wav: %v", err)
	}

	app := &App{}
	if _, err := app.GetTTSAudioData(outside); err == nil {
		t.Fatal("expected a refusal for a path outside the tts cache")
	}
}

func TestGetTTSAudioDataRefusesNonWAV(t *testing.T) {
	cache := t.TempDir()
	t.Setenv("HOWL_AUDIO_CACHE_DIR", cache)

	path := filepath.Join(cache, "tts_2.txt")
	if err := os.WriteFile(path, minimalWAV, 0600); err != nil {
		t.Fatalf("write non-wav: %v", err)
	}

	app := &App{}
	if _, err := app.GetTTSAudioData(path); err == nil {
		t.Fatal("expected a refusal for a non-wav cache path")
	}
}
