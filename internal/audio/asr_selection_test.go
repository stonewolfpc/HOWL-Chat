package audio

import (
	"strings"
	"testing"

	"howl-chat/internal/audio/capability"
	"howl-chat/internal/audio/types"
)

// selectableRecord builds a capability record whose every dimension permits
// use, so explicit-selection precedence can be tested without machine assets.
func selectableTTSRecord(id string) capability.Record {
	record := selectableRecord(id)
	record.Task = capability.TaskSpeechSynthesis
	return record
}

func selectableRecord(id string) capability.Record {
	return capability.Record{
		ID:            id,
		Task:          capability.TaskFileTranscription,
		Local:         true,
		ProviderState: capability.Avail("test provider"),
		ModelState:    capability.Avail("test model"),
		Projector:     capability.NoRequirement("test projector"),
		Runtime:       capability.Avail("test runtime"),
		TaskSupport:   capability.Avail("test task"),
		Hardware:      capability.Avail("test hardware"),
		Verification:  types.VerificationImplementation,
	}
}

// DEBT-T3: automatic selection must not leave the local, admitted set, so a
// missing local model can never silently become a remote API call.
func TestAutomaticASRStaysLocalAndAdmitted(t *testing.T) {
	planned := []types.ASRModelType{
		types.ASRQwen3ASR17B,      // remote API
		types.ASRQwen3ASR17BLocal, // local and admitted after E4
		types.ASRWhisperBase,
		types.ASRWhisperTiny,
	}

	records := []capability.Record{
		selectableRecord(string(types.ASRWhisperBase)),
		selectableRecord(string(types.ASRWhisperTiny)),
		selectableRecord(string(types.ASRQwen3ASR17B)),
		selectableRecord(string(types.ASRQwen3ASR17BLocal)),
	}
	candidates, err := asrCandidatesForRecords(planned, "auto", records)
	if err != nil {
		t.Fatalf("asrCandidatesForRecords: %v", err)
	}
	if len(candidates) != 3 {
		t.Fatalf("candidates = %v, want three admitted local models", candidates)
	}
	for _, candidate := range candidates {
		profile, ok := types.GetASRProfile(candidate)
		if !ok {
			t.Fatalf("candidate %q has no profile", candidate)
		}
		if !profile.Local || !profile.ApprovedLocal {
			t.Fatalf("automatic selection included %q", candidate)
		}
		if profile.Backend != "local" {
			t.Fatalf("automatic selection crossed into backend %q", profile.Backend)
		}
	}
}

func TestExplicitASRIsHonouredOrRefused(t *testing.T) {
	records := []capability.Record{
		selectableRecord("whisper-base"),
		selectableRecord("qwen3-asr-1.7b"),
	}

	tests := []struct {
		name    string
		engine  string
		wantID  types.ASRModelType
		wantErr string
	}{
		{name: "explicit whisper base", engine: "whisper-base", wantID: types.ASRWhisperBase},
		{name: "explicit remote qwen", engine: "qwen3-asr-1.7b", wantID: types.ASRQwen3ASR17B},
		{name: "unknown engine", engine: "whisper-medium", wantErr: "whisper-medium"},
		{name: "blocked local qwen without assets", engine: "qwen3-asr-1.7b-local", wantErr: "qwen3-asr-1.7b-local"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidates, err := asrCandidatesForRecords(nil, test.engine, records)
			if test.wantErr != "" {
				if err == nil {
					t.Fatalf("asrCandidatesForRecords(%q) = %v, want an error", test.engine, candidates)
				}
				if !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("error = %v, want it to name %q", err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("asrCandidatesForRecords(%q): %v", test.engine, err)
			}
			if len(candidates) != 1 || candidates[0] != test.wantID {
				t.Fatalf("candidates = %v, want [%s]", candidates, test.wantID)
			}
		})
	}
}

// An explicit engine whose dimensions are blocked is refused with the blocker,
// not replaced by another provider.
func TestExplicitASRWithBlockedDimensionsIsRefused(t *testing.T) {
	blocked := selectableRecord("whisper-tiny")
	blocked.ProviderState = capability.Blocked("whisper-cli is not installed")

	if _, err := asrCandidatesForRecords(nil, "whisper-tiny", []capability.Record{blocked}); err == nil {
		t.Fatal("a blocked explicit engine was accepted")
	} else if !strings.Contains(err.Error(), "whisper-cli is not installed") {
		t.Fatalf("error = %v, want the blocker reported", err)
	}
}

func TestExplicitTTSIsHonouredOrRefused(t *testing.T) {
	planned := []types.TTSModelType{types.TTSPiper, types.TTSKokoro}
	records := []capability.Record{selectableTTSRecord("piper"), selectableTTSRecord("kokoro")}

	candidates, err := ttsCandidatesForRecords(planned, "", records)
	if err != nil || len(candidates) != 2 {
		t.Fatalf("blank engine = %v, %v; want the plan unchanged", candidates, err)
	}

	candidates, err = ttsCandidatesForRecords(planned, "kokoro", records)
	if err != nil {
		t.Fatalf("explicit kokoro: %v", err)
	}
	if len(candidates) != 1 || candidates[0] != types.TTSKokoro {
		t.Fatalf("candidates = %v, want [kokoro]", candidates)
	}

	if _, err := ttsCandidatesForRecords(planned, "not-an-engine", records); err == nil {
		t.Fatal("unknown TTS engine was accepted")
	}
}

// An explicit TTS engine whose voice model or service is missing is refused at
// selection with the contract's blocker, not at load time.
func TestExplicitTTSWithBlockedDimensionsIsRefused(t *testing.T) {
	blocked := selectableTTSRecord("kokoro")
	blocked.ModelState = capability.Blocked("voice model is not installed")

	if _, err := ttsCandidatesForRecords(nil, "kokoro", []capability.Record{blocked}); err == nil {
		t.Fatal("an explicit TTS engine with a missing voice model was accepted")
	} else if !strings.Contains(err.Error(), "voice model is not installed") {
		t.Fatalf("error = %v, want the blocker reported", err)
	}
}
