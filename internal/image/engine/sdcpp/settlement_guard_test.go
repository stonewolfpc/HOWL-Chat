package sdcpp

import (
	"errors"
	"testing"

	"howl-chat/internal/image/types"
)

// Review-1: a handle settles exactly once. Settling it again must not move the
// slot, and must not re-describe a context that cleanup already cancelled.
func TestSettledHandleCannotSettleTwice(t *testing.T) {
	e := NewEngine()
	e.currentModel = &types.ModelInfo{Name: "m", Family: "z-image"}

	h, err := e.beginGeneration("job1")
	if err != nil {
		t.Fatalf("beginGeneration: %v", err)
	}

	if err := e.settleGeneration(h, nil); err != nil {
		t.Fatalf("first settlement should succeed: %v", err)
	}
	if got := e.GenerationState(); got != "idle" {
		t.Fatalf("state = %q after completion, want idle", got)
	}

	// The defect: with only genHandle/genCancel cleared, this still matched the
	// token and pushed the idle engine back to abandoned.
	secondErr := e.settleGeneration(h, errors.New("late failure"))
	if got := e.GenerationState(); got != "idle" {
		t.Fatalf("a settled handle moved the slot to %q", got)
	}
	if !errors.Is(secondErr, errGenerationSettled) {
		t.Fatalf("want errGenerationSettled from the second settlement, got %v", secondErr)
	}

	// The slot must still be genuinely free afterwards.
	next, err := e.beginGeneration("job2")
	if err != nil {
		t.Fatalf("the slot should be claimable after a double settlement: %v", err)
	}
	if err := e.settleGeneration(next, nil); err != nil {
		t.Fatalf("settling the next generation: %v", err)
	}
	if got := e.GenerationState(); got != "idle" {
		t.Fatalf("state = %q after the following generation, want idle", got)
	}
}

// The guard must not change a legitimate first settlement: each outcome still
// classifies and releases (or retains) exactly as the frozen contract says.
func TestFirstSettlementStillClassifiesEachOutcome(t *testing.T) {
	cases := []struct {
		name     string
		cancel   bool
		err      error
		want     string
		wantSent error
	}{
		{"completed", false, nil, "idle", nil},
		{"cancelled", true, errors.New("context canceled"), "abandoned", errClientCancelled},
		{"ambiguous", false, errors.New("backend said no"), "abandoned", errEngineUncertain},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := NewEngine()
			e.currentModel = &types.ModelInfo{Name: "m", Family: "z-image"}

			h, err := e.beginGeneration("job1")
			if err != nil {
				t.Fatalf("beginGeneration: %v", err)
			}
			if tc.cancel && !e.cancelGeneration("job1") {
				t.Fatal("the owner's cancel did not match")
			}

			err = e.settleGeneration(h, tc.err)
			if tc.wantSent == nil {
				if err != nil {
					t.Fatalf("want no error, got %v", err)
				}
			} else if !errors.Is(err, tc.wantSent) {
				t.Fatalf("want %v, got %v", tc.wantSent, err)
			}
			if got := e.GenerationState(); got != tc.want {
				t.Fatalf("state = %q, want %q", got, tc.want)
			}
		})
	}
}
