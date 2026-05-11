package runtime

import (
	"testing"

	"howl-chat/internal/audio/types"
)

func TestBuildSelectionPlanTinyTier(t *testing.T) {
	plan := BuildSelectionPlan(Constraints{
		Tier:             types.TierTiny,
		RequireStreaming: true,
		PreferredBackend: "local",
	})

	if len(plan.ASR) == 0 {
		t.Fatal("expected ASR candidates for tiny tier")
	}
	if plan.ASR[0] != types.ASRWhisperTiny {
		t.Fatalf("expected whisper tiny first, got %s", plan.ASR[0])
	}
	if len(plan.TTS) == 0 || plan.TTS[0] != types.TTSPiper {
		t.Fatalf("expected piper first for tiny tier, got %v", plan.TTS)
	}
	if len(plan.Multimodal) != 0 {
		t.Fatalf("expected no multimodal candidates for tiny tier, got %d", len(plan.Multimodal))
	}
}

func TestBuildSelectionPlanBackendPreference(t *testing.T) {
	plan := BuildSelectionPlan(Constraints{
		Tier:             types.TierSmall,
		PreferredBackend: "local",
		Languages:        []string{"en"},
	})

	if len(plan.ASR) == 0 {
		t.Fatal("expected ASR candidates for small tier")
	}
	topProfile, ok := types.GetASRProfile(plan.ASR[0])
	if !ok {
		t.Fatalf("missing profile for top model %s", plan.ASR[0])
	}
	if topProfile.Backend != "local" {
		t.Fatalf("expected local backend model first, got backend=%s model=%s", topProfile.Backend, plan.ASR[0])
	}
}
