package main

import (
	"os"
	"path/filepath"
	"testing"

	audiotypes "howl-chat/internal/audio/types"
	"howl-chat/internal/image/registry"
)

func TestBuildPiperVoiceCatalogMarksInstalledAndMissingVoices(t *testing.T) {
	dir := t.TempDir()
	// buildPiperVoiceCatalog always merges the configured default voice, which
	// otherwise comes from the user's models/ directory. Point that input at a
	// path inside this test's temp dir so the fixture owns every voice and a
	// real installed model cannot decide the default.
	t.Setenv("HOWL_TTS_PIPER_MODEL", filepath.Join(dir, "default-voice.onnx"))
	installed := filepath.Join(dir, "present.onnx")
	if err := os.WriteFile(installed, []byte("model"), 0o600); err != nil {
		t.Fatal(err)
	}

	catalog := buildPiperVoiceCatalog(TTSVoiceCatalog{}, []registry.Entry{
		{Name: "present", Path: installed},
	}, audiotypes.TTSModelProfile{SupportedVoices: []string{"missing"}})

	byID := make(map[string]TTSVoiceOption)
	for _, voice := range catalog.Voices {
		byID[voice.ID] = voice
	}
	if !byID["present"].Available {
		t.Fatalf("present voice = %+v, want available", byID["present"])
	}
	if byID["missing"].Available {
		t.Fatalf("missing voice = %+v, want unavailable", byID["missing"])
	}
	if !catalog.Available || catalog.DefaultVoice != "present" {
		t.Fatalf("catalog = %+v, want available with present default", catalog)
	}
}

func TestVoiceLanguageUsesLocalePrefix(t *testing.T) {
	if got := voiceLanguage("en_GB-northern_english_male-medium"); got != "en-gb" {
		t.Fatalf("voiceLanguage = %q, want en-gb", got)
	}
}
