# Audio Runtime Model Loading

## Purpose

This module defines a production-grade contract for selecting, loading, and unloading audio-capable models in a way that is:

- modular (providers are registered, not hardcoded),
- hardware-aware (tier + capability filtering),
- testable (pure planner + mockable loader),
- and safe for long-running sessions (explicit lifecycle management).

## Runtime Components

### `internal/audio/runtime/planner.go`

`BuildSelectionPlan()` computes ranked model candidates for:

- ASR
- TTS
- multimodal

The planner accepts constraints such as:

- memory tier (`tiny`, `small`, `medium`, `large`),
- streaming requirement,
- GPU requirement,
- preferred backend (`local`, `http`, etc.),
- language needs.

The planner is intentionally stateless and side-effect free, so it can be reused by:

- startup auto-configuration,
- per-conversation routing,
- fallback strategy builders.

### `internal/audio/runtime/loader.go`

`ModelLoader` is a registry-driven loader implementing the `types.ModelLoader` contract.

It provides:

- constructor registration for each model family,
- initialization hooks for loaded instances,
- in-memory tracking of loaded model IDs,
- explicit unload support for resource cleanup.

Model IDs are namespaced by capability type:

- `asr:<model>`
- `tts:<model>`
- `multimodal:<model>`

## Extension Contract

To add a new model provider:

1. implement the relevant interface in `internal/audio/types/interfaces.go`;
2. register constructor with `ModelLoader`:
   - `RegisterASR()`,
   - `RegisterTTS()`,
   - or `RegisterMultimodal()`;
3. ensure `Initialize()` is deterministic and returns actionable errors;
4. add planner-compatible metadata to model profile registry.

## Design Rules

- Keep each file under 500 lines.
- Keep business logic out of transport and provider adapters.
- Keep tier comparisons numeric, never lexical string comparison.
- Keep planner deterministic so behavior is predictable for users.
- Keep unload paths explicit to avoid leaked long-lived processes.

## Recommended Next Layer

1. Add provider adapters for Piper and MeloTTS under `internal/audio/tts/`.
2. Add a concrete `AudioPipeline` implementation that composes:
   - model planning,
   - loader lifecycle,
   - ASR/TTS routing.
3. Add integration tests for:
   - low-tier fallback behavior,
   - mixed backend routing (`local` + `http`),
   - unload/reload stability for multi-character sessions.
