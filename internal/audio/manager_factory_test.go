package audio

import (
	"context"
	"errors"
	"testing"
)

func TestProductionFactoryMapsCapabilityOutcomes(t *testing.T) {
	asrErr := errors.New("ASR init failed")
	pipeline, _ := newScriptedPipeline(providerBehavior{asrInitErr: asrErr})

	candidate := &Candidate{}
	err := newProductionFactory(func() *Pipeline { return pipeline })(context.Background(), settingsBoth, candidate)
	if err != nil {
		t.Fatalf("factory returned top-level error: %v", err)
	}
	if candidate.Pipeline != pipeline {
		t.Fatal("candidate did not receive the constructed pipeline")
	}
	if candidate.ASR.Ready {
		t.Fatal("candidate ASR is ready, want failed")
	}
	if !errors.Is(candidate.ASR.Err, asrErr) {
		t.Fatalf("candidate ASR error = %v, want %v", candidate.ASR.Err, asrErr)
	}
	if !candidate.TTS.Ready {
		t.Fatalf("candidate TTS is not ready: %v", candidate.TTS.Err)
	}
	if candidate.TTS.Err != nil {
		t.Fatalf("candidate TTS error = %v, want nil", candidate.TTS.Err)
	}
	if err := pipeline.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestProductionFactoryAssignsCandidateBeforeInitialize(t *testing.T) {
	candidate := &Candidate{}
	var recovered any
	func() {
		defer func() { recovered = recover() }()
		_ = newProductionFactory(func() *Pipeline { return &Pipeline{} })(context.Background(), settingsASR, candidate)
	}()

	if recovered == nil {
		t.Fatal("factory did not panic from the nil-loader initialization")
	}
	if candidate.Pipeline == nil {
		t.Fatal("candidate pipeline was not assigned before Initialize")
	}
}

func TestNewProductionManagerStartsIdleAndProviderFree(t *testing.T) {
	manager := NewProductionManager()
	if manager == nil {
		t.Fatal("NewProductionManager returned nil")
	}
	if status := manager.Status(); status.Phase != "idle" || status.Builds != 0 {
		t.Fatalf("initial status = %+v, want idle with no builds", status)
	}
	if err := mustClose(t, manager); err != nil {
		t.Fatalf("Close: %v", err)
	}
}
