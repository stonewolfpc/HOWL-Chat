package gguf

import "testing"

func TestClassifyProjectorMetadata(t *testing.T) {
	tests := []struct {
		name         string
		keys         map[string]any
		wantKind     ProjectorKind
		wantVision   bool
		wantAudio    bool
		wantEvidence string
	}{
		{
			name:       "empty metadata is unknown, not vision",
			keys:       map[string]any{},
			wantKind:   ProjectorUnknown,
			wantVision: false,
			wantAudio:  false,
		},
		{
			name: "audio flag only",
			keys: map[string]any{
				"clip.has_audio_encoder": true,
				"clip.audio.n_mel_bins":  uint32(128),
			},
			wantKind:     ProjectorAudio,
			wantVision:   false,
			wantAudio:    true,
			wantEvidence: "clip.has_audio_encoder",
		},
		{
			name: "audio key space without flag",
			keys: map[string]any{
				"clip.audio.sample_rate": float32(16000),
			},
			wantKind:     ProjectorAudio,
			wantVision:   false,
			wantAudio:    true,
			wantEvidence: "clip.audio.sample_rate",
		},
		{
			name: "vision flag only",
			keys: map[string]any{
				"clip.has_vision_encoder": true,
				"clip.vision.image_size":  uint32(336),
			},
			wantKind:     ProjectorVision,
			wantVision:   true,
			wantAudio:    false,
			wantEvidence: "clip.has_vision_encoder",
		},
		{
			name: "both encoders",
			keys: map[string]any{
				"clip.has_audio_encoder":  true,
				"clip.has_vision_encoder": true,
			},
			wantKind:   ProjectorHybrid,
			wantVision: true,
			wantAudio:  true,
		},
		{
			name: "explicit false flag vetoes an accidental prefix",
			keys: map[string]any{
				"clip.has_vision_encoder": false,
				"clip.vision.stray_key":   uint32(1),
			},
			wantKind:   ProjectorUnknown,
			wantVision: false,
			wantAudio:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info := ClassifyProjectorMetadata(test.keys)
			if info.Kind != test.wantKind {
				t.Fatalf("kind = %q, want %q", info.Kind, test.wantKind)
			}
			if info.HasVision() != test.wantVision {
				t.Fatalf("HasVision() = %t, want %t", info.HasVision(), test.wantVision)
			}
			if info.HasAudio() != test.wantAudio {
				t.Fatalf("HasAudio() = %t, want %t", info.HasAudio(), test.wantAudio)
			}
			if test.wantEvidence != "" {
				found := false
				for _, key := range info.Evidence {
					if key == test.wantEvidence {
						found = true
					}
				}
				if !found {
					t.Fatalf("evidence = %v, want it to contain %q", info.Evidence, test.wantEvidence)
				}
			}
		})
	}
}

func TestInspectProjectorMissingFileIsUnknown(t *testing.T) {
	info := InspectProjector("does-not-exist-mmproj.gguf")
	if info.Kind != ProjectorUnknown {
		t.Fatalf("kind = %q, want unknown", info.Kind)
	}
	if info.HasVision() {
		t.Fatal("unreadable projector claimed vision support")
	}
	if info.Err == "" {
		t.Fatal("unreadable projector reported no error")
	}
}
