package audio

import (
	"context"

	"howl-chat/internal/audio/types"
	"howl-chat/internal/process"
)

// NewProductionManager constructs the single production manager. Construction
// is intentionally provider-free: no model, binary, or network I/O happens
// until Apply asks the factory to build a generation.
func NewProductionManager(registries ...*process.Registry) *Manager {
	return NewManager(NewProductionFactory(registries...))
}

// NewProductionFactory returns the production pipeline factory.
func NewProductionFactory(registries ...*process.Registry) PipelineFactory {
	return newProductionFactory(func() *Pipeline { return NewManagedPipeline(registries...) })
}

// newProductionFactory keeps the pipeline constructor injectable for tests
// without mutable package state. The candidate receives the pipeline before
// Initialize so the manager can recover and clean up an initialization panic.
// Capability outcomes are mapped individually; aggregate initialization
// failures stay per-capability rather than becoming a factory-wide error.
func newProductionFactory(create func() *Pipeline) PipelineFactory {
	return func(ctx context.Context, desired Settings, candidate *Candidate) error {
		pipeline := create()
		if candidate != nil {
			candidate.Pipeline = pipeline
		}
		if pipeline == nil {
			return nil
		}

		_ = pipeline.Initialize(ctx, types.PipelineConfig{
			ASREnabled:        desired.ASREnabled,
			TTSEnabled:        desired.TTSEnabled,
			MultimodalEnabled: false,
			RealtimeEnabled:   false,
			HardwareTier:      "auto",
			ASREngine:         desired.ASREngine,
			TTSEngine:         desired.TTSEngine,
			TTSVoice:          desired.TTSVoice,
			TTSSpeed:          desired.TTSSpeed,
			ASRLanguage:       desired.ASRLanguage,
			ASRThreads:        desired.ASRThreads,
			ASRBeamSize:       desired.ASRBeamSize,
			ASRBestOf:         desired.ASRBestOf,
		})

		outcome := pipeline.initializationResult()
		if candidate == nil {
			return nil
		}
		if desired.ASREnabled {
			candidate.ASR = CapabilityResult{
				Ready: outcome.asr.ready,
				Err:   outcome.asr.err,
			}
		}
		if desired.TTSEnabled {
			candidate.TTS = CapabilityResult{
				Ready: outcome.tts.ready,
				Err:   outcome.tts.err,
			}
		}
		return nil
	}
}
