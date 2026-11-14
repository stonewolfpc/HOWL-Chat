package whisper

import (
	"context"
	"strconv"
	"testing"

	"howl-chat/internal/audio/types"
)

func whisperArgValue(args []string, flag string) (string, bool) {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == flag {
			return args[i+1], true
		}
	}
	return "", false
}

func transcribeWithDecodeConfig(t *testing.T, cfg types.ASRConfig) []string {
	t.Helper()
	model := whisperTestModel(t)
	cli := whisperTestCLI(t)
	audio := writeWhisperWAV(t)
	cfg.ModelPath = model

	r := newWhisperTestRecognizer(t, cli, NewLimiter(1))
	var got []string
	r.runCLI = func(_ context.Context, args []string) ([]byte, error) {
		got = append([]string(nil), args...)
		if err := writeWhisperJSON(whisperOutputPath(args)); err != nil {
			return nil, err
		}
		return nil, nil
	}
	if err := r.Initialize(context.Background(), cfg); err != nil {
		t.Fatalf("initialize: %v", err)
	}
	t.Cleanup(func() { _ = r.Release() })
	if _, err := r.Transcribe(context.Background(), audio, nil); err != nil {
		t.Fatalf("transcribe: %v", err)
	}
	return got
}

func TestWhisperDecodeSettingsReachFinalCLI(t *testing.T) {
	tests := []struct {
		name       string
		language   string
		threads    int
		beamSize   int
		bestOf     int
		wantLang   string
		wantNoLang bool
	}{
		{
			name: "automatic language and explicit decode values", language: "auto",
			threads: 3, beamSize: 2, bestOf: 7, wantNoLang: true,
		},
		{
			name: "explicit language and different decode values", language: "de",
			threads: 11, beamSize: 9, bestOf: 3, wantLang: "de",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := transcribeWithDecodeConfig(t, types.ASRConfig{
				Language: tc.language, NumThreads: tc.threads,
				BeamSize: tc.beamSize, BestOf: tc.bestOf,
			})

			if tc.wantNoLang {
				if v, ok := whisperArgValue(got, "--language"); ok {
					t.Fatalf("--language = %q, want argument omitted for automatic detection", v)
				}
			} else if v, ok := whisperArgValue(got, "--language"); !ok || v != tc.wantLang {
				t.Fatalf("--language = %q (present=%v), want %q", v, ok, tc.wantLang)
			}
			if v, ok := whisperArgValue(got, "--threads"); !ok || v != strconv.Itoa(tc.threads) {
				t.Fatalf("--threads = %q (present=%v), want %d", v, ok, tc.threads)
			}
			if v, ok := whisperArgValue(got, "--beam-size"); !ok || v != strconv.Itoa(tc.beamSize) {
				t.Fatalf("--beam-size = %q (present=%v), want %d", v, ok, tc.beamSize)
			}
			if v, ok := whisperArgValue(got, "--best-of"); !ok || v != strconv.Itoa(tc.bestOf) {
				t.Fatalf("--best-of = %q (present=%v), want %d", v, ok, tc.bestOf)
			}
		})
	}
}

func TestWhisperThreadChangeChangesInvocation(t *testing.T) {
	first := transcribeWithDecodeConfig(t, types.ASRConfig{
		Language: "auto", NumThreads: 2, BeamSize: 5, BestOf: 5,
	})
	second := transcribeWithDecodeConfig(t, types.ASRConfig{
		Language: "auto", NumThreads: 12, BeamSize: 5, BestOf: 5,
	})

	firstThreads, firstOK := whisperArgValue(first, "--threads")
	secondThreads, secondOK := whisperArgValue(second, "--threads")
	if !firstOK || firstThreads != "2" {
		t.Fatalf("first --threads = %q (present=%v), want 2", firstThreads, firstOK)
	}
	if !secondOK || secondThreads != "12" {
		t.Fatalf("second --threads = %q (present=%v), want 12", secondThreads, secondOK)
	}
	if firstThreads == secondThreads {
		t.Fatal("changing the configured thread count did not change the Whisper invocation")
	}
}
