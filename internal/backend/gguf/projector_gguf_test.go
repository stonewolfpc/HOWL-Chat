package gguf

import "testing"

// Minimal GGUF v3 writer: header plus metadata, no tensors. It lets the
// projector regression run against a real file so the metadata path is
// exercised, not just the classifier.
func writeMinimalGGUF(t *testing.T, kvs map[string]any) string {
	t.Helper()
	return writeGGUFFile(t, t.TempDir(), "mmproj-test.gguf", kvs)
}

// DEBT-T1 regression: an audio projector must not enable vision.
func TestInspectProjectorAudioFileDoesNotEnableVision(t *testing.T) {
	path := writeMinimalGGUF(t, map[string]any{
		"general.architecture":    "clip",
		"clip.has_audio_encoder":  true,
		"clip.has_vision_encoder": false,
		"clip.audio.n_mel_bins":   uint32(128),
	})

	info := InspectProjector(path)
	if info.Err != "" {
		t.Fatalf("InspectProjector: %s", info.Err)
	}
	if info.Kind != ProjectorAudio {
		t.Fatalf("kind = %q, want audio", info.Kind)
	}
	if info.HasVision() {
		t.Fatal("an audio projector enabled vision")
	}
	if !info.HasAudio() {
		t.Fatal("an audio projector did not report audio")
	}
}

func TestInspectProjectorVisionFileKeepsVision(t *testing.T) {
	path := writeMinimalGGUF(t, map[string]any{
		"general.architecture":    "clip",
		"clip.has_vision_encoder": true,
		"clip.has_audio_encoder":  false,
		"clip.vision.image_size":  uint32(336),
	})

	info := InspectProjector(path)
	if info.Err != "" {
		t.Fatalf("InspectProjector: %s", info.Err)
	}
	if info.Kind != ProjectorVision || !info.HasVision() {
		t.Fatalf("kind = %q vision = %t, want vision", info.Kind, info.HasVision())
	}
	if info.HasAudio() {
		t.Fatal("a vision projector reported audio input")
	}
}
