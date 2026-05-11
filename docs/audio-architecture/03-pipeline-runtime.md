# Audio Pipeline Runtime

## Scope

This stage introduces concrete runtime composition for immersion-grade audio:

- model selection via planner constraints,
- model lifecycle via loader registry,
- ASR fallback chaining,
- TTS synthesis with bounded cache.

## Runtime Modules

### `internal/audio/pipeline.go`

`Pipeline` implements `types.AudioPipeline` and orchestrates:

1. planner-driven candidate selection (`internal/audio/runtime/planner.go`),
2. modular provider loading (`internal/audio/runtime/loader.go`),
3. ASR fallback chain setup (`internal/audio/asr/fallback`),
4. cache-aware TTS synthesis (`internal/audio/tts/cache`).

### `internal/audio/tts/piper/local.go`

`piper` is the first concrete TTS backend:

- validates model and binary availability,
- synthesizes WAV output via Piper CLI,
- supports chunked streaming from generated output.

### `internal/audio/tts/cache/manager.go`

Bounded LRU metadata cache implementing `types.AudioCache`:

- max-size enforcement,
- deterministic hit/miss counters,
- explicit invalidation and clear operations.

## Operational Defaults

- ASR model paths are resolved from built-in defaults by model type.
- Piper model path uses `HOWL_TTS_PIPER_MODEL` when provided.
- Audio cache directory uses `HOWL_AUDIO_CACHE_DIR` when provided.
- Qwen ASR API key uses `HOWL_QWEN_API_KEY`.

## Design Guarantees

- File-size constraint maintained (<500 lines per file).
- Separation of concerns:
  - planner is pure selection logic,
  - loader is lifecycle management,
  - pipeline is orchestration only,
  - providers handle transport/runtime specifics.
- Errors are surfaced as structured `types.AudioError` where applicable.

## Next Hardening Steps

1. Add MeloTTS backend and backend selection policy by tier.
2. Add real multimodal bridge registration in loader.
3. Add integration tests with real binaries under build tags.
4. Add pipeline health telemetry (load times, synthesis latency, fallback frequency).
