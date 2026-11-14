# Group 2A — `ttsSpeed`: Trace and Blocker

**Date:** 2026-11-05  
**Status:** Stopped before implementation. No production or test source was changed.  
**Scope:** `ttsSpeed` only.

## Executive Finding

The value-loss point is confirmed: `synthesizeWithCache` hardcodes `Speed: 1.0`; the persisted global `ttsSpeed` never enters `DesiredSettings.Backend()`, `audio.Settings`, `PipelineConfig`, or the normal backend synthesis request. Only a per-character `VoiceOverride.Speed` can currently change the backend request.

Implementation is blocked before editing because the engines do not share one truthful speed range or one unambiguous HTTP control:

- Piper: `0.5–2.0`, transformed to the inverse `--length_scale`.
- Kokoro: `0.5–2.0`, passed directly as `--speed`.
- MeloTTS: profile says `0.8–1.2`, but the adapter accepts the global `0.5–2.0` and sends both `speed` and inverse `length_scale` to the legacy endpoint.

That makes it impossible to claim one global `0.5–2.0` value is applied to every supported engine without either silently clamping, overclaiming support, or introducing an engine-aware range policy. This is the Group 2A stop condition.

## Existing Path

### 1. Persisted setting

- Key: `ttsSpeed`
- Type: `float64`
- Default: `1.0`
- Validation: `0.5 <= speed <= 2.0`
- Renderer slider: `0.5–2.0`, step `0.05`
- Current status: `unsupported` (“This TTS option is not wired in this build.”)

Source references:

- `internal/audio/settings_contract.go:41` — persisted field.
- `internal/audio/settings_contract.go:84` — default `1`.
- `internal/audio/settings_contract.go:181-182` — global range validation.
- `frontend/dist/js/audio-settings.js:322` — slider bounds.
- `internal/audio/settings_contract.go:602` — current unsupported status.

### 2. Renderer-to-Go request

The full 34-field settings object is persisted and sent through `ApplyAudioSettings`. The browser-TTS path already uses the global value as `SpeechSynthesisUtterance.rate`.

The backend path differs:

- Normal backend playback calls `SynthesizeSpeech(text)` with text only.
- Character playback calls `SynthesizeSpeechWithVoice(text, voiceJSON)` and can carry a per-character `VoiceOverride.Speed`.
- The global `ttsSpeed` is not included in the normal backend request.

Source references:

- `frontend/dist/js/audio-service.js:446` — merged character/global speed.
- `frontend/dist/js/audio-service.js:524` — browser `utterance.rate`.
- `app-audio.go` — `SynthesizeSpeech` and `SynthesizeSpeechWithVoice` construct `AudioMessage` without a global speed field.

### 3. Desired-settings contract

`DesiredSettings.Backend()` currently adopts only:

- ASR/TTS enable flags
- ASR/TTS engine names
- the four Whisper decode values

`TTSSpeed` is excluded. `Settings` and `PipelineConfig` likewise have no TTS speed field.

Source references:

- `internal/audio/settings_contract.go:120-133` — `Backend()`.
- `internal/audio/manager.go:60-74` — operational `Settings`.
- `internal/audio/types/interfaces.go:270-288` — `PipelineConfig`.

### 4. TTS manager and active generation

The manager stores the adopted `Settings` in the active generation and exposes it through `AppliedSettings()`. A successful replacement carries those settings forward; a failed replacement preserves the previous generation.

No TTS speed value exists in that chain, so there is nothing for the status report to compare.

### 5. Loss point: `synthesizeWithCache`

The normal synthesis request is built as:

```go
req := types.SynthesisRequest{
    Text:  text,
    Voice: voice,
    Speed: 1.0,
}
if override != nil && override.Speed > 0 {
    req.Speed = override.Speed
}
```

Thus a global value is lost at `internal/audio/pipeline.go:539-558`. A character override can still reach the synthesizer.

### 6. Final engine behavior

#### Piper

- Profile: `SupportsSpeed=true`, `MinSpeed=0.5`, `MaxSpeed=2.0`.
- Final invocation: `--length_scale <1/speed>` when `Speed > 0`.
- Unit: inverse duration scale; larger speed means smaller length scale.
- Speed is clamped to `0.5–2.0` before inversion.

Source: `internal/audio/tts/piper/local.go:189-197`; `internal/audio/types/model_profiles.go:199-201`.

#### Kokoro

- Profile: `SupportsSpeed=true`, `MinSpeed=0.5`, `MaxSpeed=2.0`.
- Final invocation: `--speed <value>` formatted to two decimals.
- Unit: direct speaking-speed multiplier consumed by `scripts/kokoro-tts.py` as the ONNX `speed` input.
- Missing/zero speed falls back to `1.0`.

Source: `internal/audio/tts/kokoro/local.go:238-243`; `scripts/kokoro-tts.py:101-118`; `internal/audio/types/model_profiles.go:261-263`.

#### MeloTTS

- Profile: `SupportsSpeed=true`, `MinSpeed=0.8`, `MaxSpeed=1.2`.
- Legacy `POST /synthesize` body sends both:
  - `speed: req.Speed`
  - `length_scale: 1/clampSpeechSpeed(req.Speed)`
- OpenAI-compatible `POST /v1/audio/speech` sends `speed`, clamped by the adapter to `0.25–4.0`.
- The adapter’s shared request validation accepts `0.5–2.0`; it does not enforce the profile’s `0.8–1.2` bounds.
- The repository does not establish which field the server treats as authoritative, or whether its accepted range matches the profile.

Source: `internal/audio/tts/melotts/client.go:222-232`, `:285-292`; `internal/audio/tts/tts.go:85-91`; `internal/audio/types/model_profiles.go:230-232`.

### 7. Cache identity

`GenerateCacheKey(text, voice, speed)` includes speed at six-decimal precision:

```go
sha256(fmt.Sprintf("%s|%s|%.6f", voice, text, speed))
```

`synthesizeWithCache` already uses that key. Different requested speeds therefore produce different keys and cannot accidentally reuse materially different audio.

The adapters round at their final boundary (Kokoro to two decimals, Piper’s inverse length scale to four decimals), so the cache may over-distinguish values that round to the same final argument. That causes extra synthesis misses, not incorrect reuse. No cache redesign is required for Group 2A.

Source: `internal/audio/types/utils.go:266-269`; `internal/audio/pipeline.go:554`.

## Blocking Ambiguity

The current global contract has one range: `0.5–2.0`.

The engine profiles do not agree:

| Engine | Profile range | Adapter behavior |
|---|---|---|
| Piper | `0.5–2.0` | Clamps to `0.5–2.0`, sends inverse `--length_scale` |
| Kokoro | `0.5–2.0` | Sends `--speed`; script consumes direct multiplier |
| MeloTTS | `0.8–1.2` | Shared validation accepts `0.5–2.0`; sends `speed` plus inverse `length_scale`; OpenAI fallback clamps to `0.25–4.0` |

Consequences:

- Reporting `ttsSpeed` applied for MeloTTS across the whole global range would overclaim support.
- Silently clamping MeloTTS to `0.8–1.2` would change the user’s requested value and violate the truthfulness contract.
- Applying MeloTTS values only within `0.8–1.2` requires engine-aware validation/status and a policy for auto-selected TTS, which is a narrow contract change that the task permits only after the semantics are resolved.
- The MeloTTS dual-field HTTP request leaves the authoritative speed control and actual accepted range unverified.

This is not a reason to redesign the cache or change voice selection. It is a bounded engine-capability decision that must be made before wiring.

## Decision Required

One of the following must be resolved before implementation:

1. **Engine-aware range policy** — keep `0.5–2.0` globally, accept it for Piper/Kokoro, and reject/mark unavailable for MeloTTS outside `0.8–1.2`. This also needs an explicit policy for `ttsEngine=auto` when the chosen engine is unknown.
2. **Verify MeloTTS semantics** — establish from the actual MeloTTS deployment/API that the full `0.5–2.0` range is accepted and whether `speed` or `length_scale` is authoritative, then decide whether the profile bounds are stale.
3. **Narrow the global setting** — change the global range to the intersection, `0.8–1.2`; this regresses Piper/Kokoro and should not be done without explicit approval.

## Not Done

- No production code changed.
- No test code changed.
- No cache changes.
- No TTS voice selection changes.
- No microphone work was needed; the task does not depend on a connected microphone.
- No Go tests, Playwright tests, or Wails build were run because no implementation was made.

## Evidence Grades

- **Observed:** source code and profile values listed above.
- **Source-derived:** the loss point, final Piper/Kokoro arguments, MeloTTS request fields, and cache-key composition.
- **Unestablished:** MeloTTS server’s authoritative speed field and accepted range; real audible speed changes for any engine.
- **Blocked:** truthful status for a single global speed range across all TTS engines.
