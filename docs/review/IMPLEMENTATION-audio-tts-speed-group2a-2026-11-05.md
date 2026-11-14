# Group 2A — TTS Speed (`ttsSpeed`) Implementation

**Date:** 2026-11-05

**Scope:** `ttsSpeed` only. Group 2B (TTS voice) and later audio groups were not started.

**Status:** Implementation complete; pending independent review. No commit or push was performed.

## 1. Source-to-invocation trace

The persisted speed now follows this path:

1. `DesiredSettings.TTSSpeed` (`json:"ttsSpeed"`, default `1.0`) is loaded from the persisted audio contract.
2. `DesiredSettings.Backend()` copies `TTSSpeed` into `audio.Settings`.
3. `Manager.Apply` records the desired generation and installs the active `Settings` only after the candidate pipeline succeeds.
4. `newProductionFactory` forwards `Settings.TTSSpeed` into `types.PipelineConfig.TTSSpeed`.
5. `Pipeline.Initialize` stores the pipeline configuration, including `TTSEngine` and `TTSSpeed`.
6. `Pipeline.ProcessAudioMessage` reads the configured values and resolves them with `ttsSpeedForConfig`.
7. `synthesizeWithCache` sets `SynthesisRequest.Speed` from the resolved global speed. A nonzero per-request `VoiceOverride.Speed` still takes precedence.
8. `GenerateCacheKey(text, voice, speed)` receives the final request speed, so generation and lookup use the same speed.
9. The selected synthesizer receives the final request:
   - **Piper:** `runPiper` retains `length_scale = 1 / speed` and passes `--length_scale`.
   - **Kokoro:** `runKokoro` receives the direct speed and passes `--speed`.
   - **MeloTTS:** the legacy `/synthesize` payload sends `speed` only; `length_scale` was removed. The OpenAI-compatible `/v1/audio/speech` payload was left unchanged.
   - **Browser TTS:** the renderer already used `audio.ttsSpeed` as `SpeechSynthesisUtterance.rate`; no Group 2A change was needed there.

## 2. Engine-aware status contract

`BuildSettingsReport` now treats `ttsSpeed` as a backend-owned setting and classifies it explicitly:

- **Explicit Piper / Kokoro / MeloTTS:** may report `applied` only when the selected engine matches the active generation, the active speed is within that engine profile’s range, and TTS is ready.
- **Automatic (`auto` or empty):** renders `unavailable` with a reason requiring an explicit engine. Auto selection is not forwarded a global speed; the pipeline uses the established `1.0` default because the effective provider is unresolved at configuration time.
- **Unsupported / unknown engine:** renders `unsupported`; it never receives an applied speed.
- **TTS disabled:** renders `unavailable`; the value remains saved but is not applied.
- **Invalid value:** renders `invalid`. The global range is `0.5–2.0`; MeloTTS keeps its existing `0.8–1.2` product guard.
- **Failed explicit replacement:** the report preserves the previous active applied speed and marks the new speed `rejected` with the activation error.

The `Applied` map only publishes a prior speed for an explicit desired engine. Auto-select cannot claim an applied speed.

`HasBackendValidationIssue` does not block ASR activation for a saved speed while TTS is disabled. Such a value is still reported as `invalid` if it is outside the global contract.

## 3. Changed-file inventory

### Backend contract and manager

- `internal/audio/settings_contract.go`
  - Added `TTSSpeed` to `Backend()`.
  - Added MeloTTS/global speed validation.
  - Moved `ttsSpeed` from `unsupportedSettings()` to `backendSettingKeys`.
  - Added explicit-engine range helpers.
  - Added engine-aware `setTTSSpeed` reporting.
- `internal/audio/settings_contract_test.go`
  - Updated unsupported/applied counts.
  - Added engine-specific speed, MeloTTS guard, failed-replacement, auto-status, and factory-forwarding tests.
- `internal/audio/manager.go`
  - Added `Settings.TTSSpeed`.
- `internal/audio/manager_factory.go`
  - Forwarded `Settings.TTSSpeed` into `PipelineConfig`.
- `internal/audio/types/interfaces.go`
  - Added `PipelineConfig.TTSSpeed`.

### Pipeline

- `internal/audio/pipeline.go`
  - Read `TTSSpeed` from pipeline configuration.
  - Added `ttsSpeedForConfig` for explicit-engine range enforcement and auto-select fallback.
  - Passed the resolved speed into `synthesizeWithCache` and `SynthesisRequest`.
  - Preserved `VoiceOverride.Speed` precedence.
- `internal/audio/pipeline_test.go`
  - Added request capture and tests for forwarding, auto-select behavior, MeloTTS guard rejection, and speed-sensitive cache keys.

### Synthesizer adapters

- `internal/audio/tts/melotts/client.go`
  - Removed `length_scale` from the legacy `/synthesize` payload.
  - Left the OpenAI-compatible payload and its existing clamping unchanged.
- `internal/audio/tts/melotts/client_test.go`
  - Added legacy payload and OpenAI payload assertions.
- `internal/audio/tts/piper/local_test.go`
  - Added inverse `length_scale` assertion for `1.25 -> 0.8000`.
- `internal/audio/tts/kokoro/local_test.go`
  - Added direct `--speed 1.25` assertion.

### App-level tests

- `app_audio_contract_test.go`
  - Updated supported counts and applied-map expectations.
  - Added failed TTS-speed replacement preservation coverage.
- `app_audio_ownership_test.go`
  - Preserved a nonzero `TTSSpeed` when a test helper builds desired backend settings.

### Renderer aggregate and fixtures

- `frontend/dist/js/audio-settings.js`
  - Requires `ttsSpeed` in the supported backend aggregate only for enabled, explicit Piper/Kokoro/MeloTTS configurations.
- `tests/e2e/specs/iframe-bridge.spec.js`
  - Updated fixture status modeling and unsupported count.
  - Added an enabled-TTS `9 of 9 backend settings applied` assertion.
- `tests/e2e/specs/activation-lock.spec.js`
  - Updated fixture status modeling for `ttsSpeed`.
- `tests/e2e/specs/audio-settings-lifecycle.spec.js`
  - Updated fixture status modeling and remote-engine unsupported count.
- `docs/review/IMPLEMENTATION-audio-tts-speed-group2a-2026-11-05.md`
  - This report.

## 4. Tests and exact results

### Formatting

```text
gofmt -w internal/audio/manager.go internal/audio/manager_factory.go internal/audio/settings_contract.go internal/audio/settings_contract_test.go internal/audio/pipeline.go internal/audio/pipeline_test.go internal/audio/types/interfaces.go internal/audio/tts/melotts/client.go internal/audio/tts/melotts/client_test.go internal/audio/tts/piper/local_test.go internal/audio/tts/kokoro/local_test.go app_audio_contract_test.go app_audio_ownership_test.go
```

Result: completed without error.

### Focused Go suites

```text
GOTOOLCHAIN=go1.23.8 go test -tags "desktop,production" ./internal/audio ./internal/audio/tts/piper ./internal/audio/tts/kokoro ./internal/audio/tts/melotts .
```

Result:

```text
ok  howl-chat/internal/audio
ok  howl-chat/internal/audio/tts/piper
ok  howl-chat/internal/audio/tts/kokoro
ok  howl-chat/internal/audio/tts/melotts
ok  howl-chat
```

### Full audio Go tree

```text
GOTOOLCHAIN=go1.23.8 go test -tags "desktop,production" ./internal/audio/...
```

Result: all packages reported `ok`, including ASR, fallback, Whisper, capability, runtime, TTS, cache, and types.

### Audio Playwright suites

```text
cd tests/e2e && npx --no-install playwright test --config=playwright.audio.config.js specs/iframe-bridge.spec.js specs/audio-settings-lifecycle.spec.js specs/activation-lock.spec.js
```

Result: **32 passed (16.9s)**.

The new frontend assertion is:

```text
explicit TTS speed is included in the supported backend aggregate
```

It sets `ttsEnabled`, `ttsEngine: "piper"`, and `ttsSpeed: 1.25`, then asserts `Applied (supported)` and `9 of 9 backend settings applied`.

### Microphone-gate regression suite

```text
cd tests/e2e && npx --no-install playwright test --config=playwright.audio.config.js specs/audio-mic-gate.spec.js
```

Result: **3 passed (2.0s)**.

This verifies the persisted ASR selection is still restored but not applied at startup, capture is blocked with the exact no-device copy, and the browser `NotFoundError` maps to the same boundary.

### JavaScript syntax

```text
node --check frontend/dist/js/audio-settings.js
```

Result: exit 0.

### Production Wails build

```text
GOTOOLCHAIN=go1.23.8 CGO_ENABLED=0 "$HOME/go/bin/wails.exe" build -devtools -s -nopackage -m -nosyncgomod -tags "desktop,production"
```

Result: exit 0; built `build/bin/howl-chat.exe` in 8.017s.

### Scoped diff check

```text
git diff --check -- internal/audio/manager.go internal/audio/manager_factory.go internal/audio/pipeline.go internal/audio/pipeline_test.go internal/audio/types/interfaces.go internal/audio/tts/melotts/client.go internal/audio/tts/melotts/client_test.go internal/audio/tts/piper/local_test.go internal/audio/tts/kokoro/local_test.go app_audio_ownership_test.go tests/e2e/specs/iframe-bridge.spec.js tests/e2e/specs/activation-lock.spec.js tests/e2e/specs/audio-settings-lifecycle.spec.js frontend/dist/js/audio-settings.js
```

Result: `DIFF_CHECK_EXIT=0`.

The repository’s known pre-existing whitespace in `IMAGE_PAGE_CHANGE_LEDGER.md` and `frontend/wailsjs/go/models.ts` was not included in this scoped check.

## 5. Required behavior evidence

The focused tests prove:

- **Persisted setting reaches the final request:** `TestProductionFactoryForwardsTTSSpeed` and `TestPipelineForwardsConfiguredTTSSpeed`.
- **Piper inverse conversion:** `TestPiperConvertsSpeedToInverseLengthScale` expects `--length_scale 0.8000` for speed `1.25`.
- **Kokoro direct speed:** `TestKokoroPassesDirectSpeed` expects `--speed 1.25`.
- **MeloTTS legacy payload:** `TestMeloTTSLegacyPayloadSendsSpeedWithoutLengthScale` checks `speed` is present and `length_scale` is absent.
- **MeloTTS OpenAI payload unchanged:** `TestMeloTTSOpenAIPayloadIsUnchanged` checks the expected fields remain and no `length_scale` was added.
- **Engine ranges and invalid values:** `TestSettingsReportMeloTTSSpeedGuardIsExplicit` and the pipeline MeloTTS guard test.
- **Auto/unresolved status:** `TestSettingsReportTTSSpeedIsEngineSpecific` and `TestPipelineAutomaticTTSDoesNotApplyConfiguredSpeed`.
- **Failed replacement preserves prior speed:** `TestSettingsReportFailedTTSSpeedReplacementPreservesPriorValue` and `TestFailedTTSSpeedReplacementKeepsLastAppliedSpeed`.
- **Distinct cache keys:** `TestPipelineSpeedParticipatesInCacheKey` proves two speeds generate two synthesis calls.
- **Frontend aggregate:** the enabled-TTS Playwright test above.

## 6. Limitations

- No live MeloTTS deployment was available. `HOWL_MELOTTS_URL` was unset, nothing listened on `8888`, and no local Melo service was available. The MeloTTS tests use in-process fake HTTP transports; they do not claim live server verification.
- No live Piper or Kokoro synthesis was required. Adapter tests assert the final argument/payload contract through the existing command and transport seams.
- The UI was exercised through the Playwright DOM assertions above, not through a manual visual review or a packaged-app screenshot.
- No commit or push was performed.

## 7. Deviations and decisions

- **Auto selection is intentionally not forwarded a global speed.** This avoids claiming that an unresolved provider accepted a speed. The pipeline keeps `1.0` for auto and the report returns `unavailable` with an explicit-engine reason.
- **TTS-disabled speed is `unavailable`, not `unsupported`.** The setting is saved but has no active TTS capability to apply it to.
- **MeloTTS `0.8–1.2` remains a product guard.** The implementation does not claim that upstream MeloTTS or the server enforces that range.
- **No cache redesign.** The existing cache key already includes speed; the implementation verifies the normal synthesis path feeds it the resolved request speed.
- **Group 0/1 boundaries preserved.** No TTS voice selection, microphone gate, startup no-auto-load behavior, lifecycle handling, headphone-only playback, or other audio setting was intentionally changed.
