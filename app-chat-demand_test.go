package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// Readiness is an answering engine, not a started process.
func TestWaitForHealthRequiresAnAnswer(t *testing.T) {
	start := time.Now()
	err := waitForHealth(func() error { return fmt.Errorf("connection refused") }, 300*time.Millisecond)
	if err == nil {
		t.Fatal("a process that never answers must not report ready")
	}
	if elapsed := time.Since(start); elapsed < 250*time.Millisecond {
		t.Fatalf("returned before the bound: %s", elapsed)
	}
}

func TestWaitForHealthReturnsWhenTheEngineAnswers(t *testing.T) {
	calls := 0
	start := time.Now()
	err := waitForHealth(func() error {
		calls++
		if calls >= 3 {
			return nil
		}
		return fmt.Errorf("connection refused")
	}, 5*time.Second)
	if err != nil {
		t.Fatalf("the engine answered and the wait still failed: %v", err)
	}
	if calls < 3 {
		t.Fatalf("returned before readiness: %d calls", calls)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("waited longer than the engine needed: %s", elapsed)
	}
}

// A send with nothing selected is a refusal with a reason, never a silent
// default load.
func TestDemandLoadRefusesWithoutASelection(t *testing.T) {
	app := &App{}
	err := app.EnsureChatModelReady("   ")
	if err == nil {
		t.Fatal("an empty selection must be refused")
	}
	if !strings.Contains(err.Error(), "no chat model selected") {
		t.Fatalf("refusal must say what is missing: %v", err)
	}
}
