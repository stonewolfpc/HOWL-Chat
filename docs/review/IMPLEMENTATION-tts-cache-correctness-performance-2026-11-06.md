# Stage 4 — TTS Cache Correctness and Performance

**Date:** 2026-11-06  
**Status:** Targeted corrections implemented; Stage 4 remains HOLD pending review.  
**Authorization:** Stage 4 corrections only. No commit, push, deployment, Character Tab UI, chat-bubble controls, resource-management changes, or later-stage work was performed.

## 1. Initial source findings

### Cache storage and lifetime

- `internal/audio/tts/cache/manager.go` implements a bounded in-memory LRU:
  - `Get` moves an entry to the front and increments hit/miss counters.
  - `Store` clones the result, updates size accounting, and evicts from the back when the byte limit is exceeded.
  - `Invalidate` and `Clear` remove entries; `GetStats` reports hits, misses, entries, current bytes, and the limit.
- `Pipeline.Initialize` creates `cache.NewManager(100)` and resets the TTS engine map. The cache is not persisted across pipeline initialization or process restart.
- `ProcessAudioMessage` snapshots the cache pointer before dispatch. `synthesizeWithCache` is the only lookup/store path for `audio_output`.
- The cache-size and cache-enabled settings remain reported as not implemented by the existing settings contract; Stage 4 did not change that contract.

### Original key and request path

The original lookup was:

```go
key := types.GenerateCacheKey(req.Text, req.Voice, req.Speed)
```

and the original key was:

```go
sha256(fmt.Sprintf("%s|%s|%.6f", voice, text, speed))
```

The original request path was:

```text
character.howl.voice
  -> AudioService.playTTS
  -> SynthesizeSpeechWithVoice
  -> parseVoiceOverrideJSON
  -> ProcessAudioMessage
  -> resolveTTSRequest
  -> synthesizeWithCache
  -> selected synthesizer
```

The original key omitted the resolved engine, model assets, language, and request model path. A same text/voice/speed request through two engines could therefore reuse the other engine's audio.

### Material synthesis inputs observed in source

- Exact `AudioMessage.Text`.
- Resolved engine (`piper`, `kokoro`, `melotts`, or the inherited active engine).
- Piper model path from `VoiceOverride.Voice` / the active model.
- Kokoro model path, voices archive, config path, and bridge script path.
- MeloTTS base URL and model setting.
- Resolved voice/speaker ID.
- MeloTTS language, which the resolver derives from the selected speaker/voice.
- Effective validated speed.
- `TTSConfig.Device` where supplied; no Stage 4 call path changes device selection.

`VoiceOverride.Volume`, `NoiseScale`, `NoiseW`, and `SentenceSilence` are parsed at the boundary but are not forwarded into `SynthesisRequest` or consumed by the current adapters. They are not part of this cache identity because they do not currently affect synthesized bytes. This is recorded as an existing unimplemented-settings boundary, not silently treated as applied.

### Engine asset facts

- **Piper:** `piper.Initialize` validates and stores `TTSConfig.ModelPath`; `runPiper` uses `SynthesisRequest.VoicePath` when supplied.
- **Kokoro:** `kokoro.Initialize` uses `ModelPath`, `VoicesPath`, `ConfigPath`, and `ScriptPath`; `runKokoro` passes the request voice as `--voice`.
- **MeloTTS:** `melotts.Initialize` uses `TTSConfig.ModelPath` as the HTTP base URL; the adapter sends the resolved speaker and language, and uses the request speed.
- All three current adapters embed `tts.BaseSynthesizer`, whose `GetConfig()` exposes the initialized configuration to the resolver.

### F-01 finding

`frontend/dist/js/audio-service.js` contained:

```js
voice: voice.voice || (engine ? '' : (global.ttsVoice || 'en_US-lessac-medium')),
```

The exact source used `en_US-lessac-medium`. For a character with `voice.enabled`, a blank engine, no character voice, and no global voice, this injected a Piper-specific voice and could fail against a non-Piper active engine.

### F-02 finding

`resolveTTSRequest` did:

```go
if override != nil && override.Speed > 0 {
    speed = override.Speed
} else {
    speed, err = ttsSpeedForConfig(speedEngine, speed)
}
```

A positive renderer speed bypassed the engine-specific guard, including MeloTTS's product range of `0.8`–`1.2`.

## 2. Final cache identity

`SynthesisRequest` now carries:

```go
Engine        string
CharacterID   string
VoicePath     string
ModelIdentity string
```

The cache key is produced by `types.GenerateCacheKeyForRequest` and includes, in order:

1. `TTSCacheKeyVersion = "howl-tts-cache-v2"`.
2. Exact input text.
3. Resolved engine, normalized for comparison.
4. Canonical model identity.
5. Resolved voice/speaker.
6. Resolved language.
7. Request voice/model path.
8. Effective speed at full float precision.

The fields are length-prefixed before hashing, so delimiters embedded in text, paths, or voice labels cannot collapse distinct requests onto the same serialization.

`CharacterID` is deliberately excluded. Two characters resolving to identical text, engine, model, voice, language, path, and speed can share audio. The character identity remains available to the selected synthesizer and in resolution diagnostics, but is not a cache partition.

### Model identity

`types.TTSModelIdentity` canonicalizes:

- engine;
- `TTSConfig.ModelPath`;
- `TTSConfig.ModelName` (MeloTTS OpenAI-compatible model, default `tts-1`);
- `TTSConfig.VoicesPath`;
- `TTSConfig.ConfigPath`;
- `TTSConfig.ScriptPath`;
- `TTSConfig.Device`;
- request `VoicePath`.

MeloTTS has two independently output-affecting values: `ModelPath` is its HTTP base URL, while `ModelName` is the `model` sent to `/v1/audio/speech` and comes from `HOWL_MELOTTS_OPENAI_MODEL`. The model is resolved into `TTSConfig` before loading, included in the cache identity, and read from that config by the adapter; the environment fallback remains for clients constructed outside the pipeline.

`CacheDir`, `CacheEnabled`, and `MaxCacheSize` are excluded because they do not change the generated audio bytes.

### Versioning and invalidation

- The version prefix prevents the old text/voice/speed entries from being treated as valid Stage 4 entries.
- `Initialize` continues to reset the in-memory cache manager.
- The legacy `GenerateCacheKey` remains available for compatibility but is no longer used by the pipeline.
- A regression stores an entry under the legacy key and verifies a normal request synthesizes fresh audio rather than returning it.

## 3. Implementation changes

- `internal/audio/types/interfaces.go`
  - Added `SynthesisRequest.ModelIdentity` and `TTSConfig.ModelName`.
- `internal/audio/types/utils.go`
  - Added the cache-key version, `GenerateCacheKeyForRequest`, and `TTSModelIdentity`.
  - Includes `ModelName` in the model identity.
  - Kept the legacy `GenerateCacheKey`.
- `internal/audio/pipeline.go`
  - Resolves the model identity for the selected synthesizer, including explicit-engine loads.
  - Captures the MeloTTS OpenAI-compatible model name in the loaded configuration.
  - Uses the complete request identity for lookup and storage.
  - Adds a bounded per-key in-flight map: the first identical request synthesizes and stores; concurrent followers wait for that result instead of invoking the engine again.
  - Validates positive override speeds with the base `0.5`–`2.0` guard and `ttsSpeedForConfig` for the resolved engine, including an inherited active engine.
- `internal/audio/tts/melotts/client.go`
  - Uses the configured Melo model name for the OpenAI-compatible request, with the existing environment fallback.
- `frontend/dist/js/audio-service.js`
  - Removed the hardcoded `en_US-lessac-medium` fallback for blank-engine inheritance. A configured global voice is still inherited; with no global voice, the backend receives a blank voice and uses the active engine's default.
- `internal/audio/pipeline_test.go`
  - Added `GetConfig()` and a deterministic synthesis delay to the scripted fake.
- New tests:
  - `internal/audio/pipeline_cache_test.go`
  - `internal/audio/types/cache_key_test.go`
  - `tests/e2e/specs/audio-character-voice.spec.js`

### Exact Stage 4 changed files

Production/runtime:

- `internal/audio/pipeline.go`
- `internal/audio/types/interfaces.go`
- `internal/audio/types/utils.go`
- `internal/audio/tts/melotts/client.go`
- `frontend/dist/js/audio-service.js`

Tests and fixtures:

- `internal/audio/pipeline_test.go`
- `internal/audio/pipeline_cache_test.go`
- `internal/audio/types/cache_key_test.go`
- `internal/audio/tts/melotts/client_test.go`
- `tests/e2e/specs/audio-character-voice.spec.js`

Report:

- `docs/review/IMPLEMENTATION-tts-cache-correctness-performance-2026-11-06.md`

No changes were made to cache storage, engine residency, model loading/unloading, UI, or character persistence.

## 4. F-01 disposition and evidence

**Disposition: fixed.**

The renderer now produces:

```js
voice: voice.voice || (engine ? '' : (global.ttsVoice || ''))
```

The regression test loads the real `audio-service.js` in Chromium and calls `playTTS` with:

- active global engine `kokoro`;
- blank character engine;
- blank character voice;
- blank global voice.

Observed request after the fix:

```json
{ "engine": "", "voice": "", "character_id": "character-inherit" }
```

A second test verifies that a non-empty global voice is still inherited. Both tests pass.

## 5. F-02 disposition and evidence

**Disposition: fixed; inherited-engine path revalidated.**

The resolver assigns `engine = activeEngine` when the request engine is blank or `auto`. Positive override validation runs after that assignment and calls `ttsSpeedForConfig(engine, speed)`, so validation follows the resolved engine rather than whether the character selected it explicitly. The earlier report phrase "resolved explicit engine" was imprecise; the production path was already resolved-engine based. A regression now covers the inherited case.

Focused results:

- Inherited MeloTTS: global/active MeloTTS, blank character engine, speed `1.5` → rejected with the `0.8`–`1.2` product guard; zero synthesis calls.
- Inherited MeloTTS: same configuration, speed `1.1` → accepted; request speed is `1.1`.
- Inherited Kokoro: global/active Kokoro, blank character engine, speed `0.4` → rejected with the `0.5`–`2.0` range; zero synthesis calls.
- Explicit MeloTTS override `1.5` → rejected with the `0.8`–`1.2` product guard; zero synthesis calls.
- Explicit MeloTTS override `1.1` → accepted; request speed is `1.1`.
- Kokoro override `0.4` → rejected with the `0.5`–`2.0` range; zero synthesis calls.
- Kokoro override `1.25` → accepted; request speed is `1.25`.
- Piper override `2.5` → rejected with the `0.5`–`2.0` range.
- Existing configured-speed and Stage 2A precedence tests continue to pass.

## 6. Performance methodology and raw measurements

All measurements use the deterministic scripted fake synthesizer, not real acoustic engines. A fixed `20 ms` synthesis delay was used for the cold workload. Timing uses Go's `time.Now`/`time.Since`; Go records a monotonic reading, so the measurements are not affected by wall-clock adjustments. The high-resolution timer is not rounded to the `0s` display that the earlier single warm sample produced.

### Cold vs warm and repeated requests

Workload: one cold request, then 10,000 identical warm requests. The warm timer wraps the full `ProcessAudioMessage` path (resolution, cache-key construction, lookup, and response path), not `cache.Get` in isolation; the average therefore includes per-request pipeline overhead.

Observed in the corrected run:

```text
cold=20.4333ms
warm_batch_requests=10000
warm_batch_total=21.8843ms
warm_batch_average=2.188µs
synthesis_calls=1
cache_hits=10000
cache_misses=1
entries=1
size_bytes=140
```

Average warm latency is `warm_batch_total / 10,000 = 2.188µs`. The earlier `warm=0s` observation was a single sample below the test's displayed timer resolution and is not used as a latency result. Synthesis evidence remains separate: one cold miss and one synthesis, followed by 10,000 cache hits.

### Alternating two characters with different engines

Workload: 10 alternating Piper and Kokoro requests, same text, distinct engines, both using a `shared` voice label.

Observed:

```text
requests=20
synthesis_calls=2
cache_hits=18
cache_misses=2
```

Each engine synthesized exactly once.

### Same text and voice across engines

Workload: Piper and Kokoro both receive `same text` and `shared`; each request is repeated once.

Observed:

```text
piper synthesis_calls=1
kokoro synthesis_calls=1
cache_hits=2
cache_misses=2
```

This is the direct regression for the old cross-engine collision.

### Character A -> B -> A

Observed:

```text
piper synthesis_calls=1
kokoro synthesis_calls=1
cache_hits=1
cache_misses=2
```

### Concurrent identical requests

Workload: 8 concurrent requests, same text/config, distinct `CharacterID`s, 30 ms fake synthesis delay.

Observed:

```text
requests=8
synthesis_calls=1
cache_hits=0
cache_misses=1
```

The one miss is the leader's cache lookup. Followers wait on the per-key result and do not perform a cache lookup of their own, so they are not counted as hits. This is the intended coalescing behavior; character identity did not partition the cache.

### Concurrent different configurations

Workload: 12 concurrent requests, six Kokoro and six Melo, same text and voice label, distinct engines.

Observed:

```text
kokoro synthesis_calls=1
melo synthesis_calls=1
cache_hits=0
cache_misses=2
```

The two complete identities remained isolated while identical requests within each identity coalesced.

### Memory

The synthetic result is 140 bytes in the cache manager's accounting. The manager remains bounded by its existing 10 MB test limit. No real-audio memory profiling was performed; no claim is made about real WAV payload retention or process RSS.

## 7. Tests and actual results

### Passed

```text
go test ./... -count=1
```

Exit `0`.

```text
go test -race ./internal/audio ./internal/audio/types ./internal/audio/tts/melotts -count=1
```

Exit `0`; no race reports.

```text
go test ./internal/audio -run 'TestInheritedEngine|TestMeloTTSModelIdentity|TestTTSCacheColdWarm|TestVoiceOverrideSpeed' -v -count=1
```

Exit `0`. The inherited MeloTTS and Kokoro invalid-speed cases pass with zero synthesis calls; the batch log reports `warm_batch_average=2.188µs`.

```text
go test ./internal/audio/types -run 'TestGenerateCacheKey|TestTTSModelIdentity' -v -count=1
```

Exit `0`; the model-identity table now includes the `model_name` separation.

```text
go test ./internal/audio/tts/melotts -run 'TestMeloTTSOpenAI' -v -count=1
```

Exit `0`; the OpenAI-compatible payload test confirms the configured model name is sent.

```text
node --check frontend/dist/js/audio-service.js
```

Exit `0`.

```text
gofmt -l <changed Go files>
```

No files listed.

```text
git diff --check -- <scoped Stage 4 files>
```

Exit `0`.

Focused F-01 Playwright spec:

```text
npx playwright test --config=playwright.audio.config.js specs/audio-character-voice.spec.js
```

`2 passed`.

All audio frontend specs, including the new spec, `audio-mic-gate.spec.js`, and `audio-settings-lifecycle.spec.js`, passed in the broad config run. The broad config also collected unrelated `html-audit`, `save-load-persistence`, and `scenario-load` specs:

```text
37 passed
22 failed
```

Those 22 failures are unrelated to Stage 4 and report `page.goto`/`Cannot navigate to invalid URL` because `playwright.audio.config.js` has no `baseURL` for full-app URLs. They are not claimed as Stage 4 passes.

### Not tested

- Real Piper/Kokoro/Melo acoustic output, real model timing, or listening checks.
- Real cache memory/RSS growth with WAV payloads.
- Real model or bridge-file replacement at the same path; see limitations.
- Cache-enabled/cache-size settings, which remain outside this stage.

## 8. Remaining limitations and untested conditions

1. **Synthetic performance only.** The latency figures are fake synthesizer timings and must not be presented as real Piper/Kokoro/Melo performance.
2. **Path/model identity is not content identity.** The model identity includes configured paths, engine, the Melo OpenAI model name, bridge/config paths, device, and request voice path. Replacing a model or bridge script in place at the same path, or changing the model served under the same Melo model name, is not detected by a content hash; a cache-version change is required for that case.
3. **Cache statistics count leader lookups.** Coalesced followers do not increment hit counters because they wait on the in-flight leader rather than calling `Get`.
4. **No persistence.** The cache remains in-memory and is reset by `Pipeline.Initialize`; there is no cross-process invalidation.
5. **Cache settings remain unsupported.** `ttsCacheEnabled` and `ttsCacheSizeMB` are not wired by this stage; the manager uses its existing construction size.
6. **No real-engine or UI performance claim.** The audio frontend regression was exercised in Chromium; it is a behavior test, not an acoustic test.

### Findings requiring separate authorization

- Wiring `ttsCacheEnabled` and `ttsCacheSizeMB` into cache construction/lifecycle would require separate settings/resource-management work; Stage 4 intentionally left the existing manager construction unchanged.
- Content hashing of model and bridge files would require a separate asset-identity decision; Stage 4 uses configured path/engine identity plus the cache version.
- Real-engine latency, memory, and acoustic validation require installed model assets and a separate measurement authorization.

## 9. Final status

**Stage 4: HOLD — corrections implemented; review pending.**

**Cache correctness: corrections implemented, acceptance pending.**  
Engine, model, voice, language, path, speed, and exact text are separated; the Melo OpenAI model name is included; equivalent requests reuse; CharacterID does not partition; legacy entries are rejected; concurrent identical requests coalesce. The inherited-engine speed and Melo model-name regressions pass.

**Performance: PARTIAL.**  
The synthetic warm batch is measured over 10,000 requests with a monotonic timer; synthesis-count evidence remains separate. Real-engine performance is still INCONCLUSIVE/not measured.

**F-01: PASS.** Hardcoded Piper fallback removed; Chromium regression covers blank-engine inheritance against Kokoro.

**F-02: PASS/revalidated.** Override speeds receive engine-specific validation for explicit and inherited engines, including MeloTTS and Kokoro, with no synthesis on invalid values.

No commit or push was performed.
