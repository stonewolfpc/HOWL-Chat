# Stage 3 — Backend Request-Scoped Character Engine / Voice Resolution

**Date:** 2026-11-06  
**Status:** Implemented and verified with focused and regression tests; ready for independent review.  
**Scope:** Backend request-scoped character TTS resolution, with the minimum renderer routing needed to carry character identity and an explicit inheritance decision.  
**Not performed:** Commit, push, Stage 4 cache work, Character Tab UI work, chat-bubble menus, model downloads, model discovery, or resource-management redesign.

## 1. Outcome

Character TTS requests are now resolved per synthesis request. A character can select a supported engine and voice without replacing the pipeline's active global engine or the global voice binding. The resulting configuration is carried on the `SynthesisRequest` to the selected engine.

The implementation preserves the original behavior for callers with no character override: they continue to use the active global generation, global voice default, and Group 2A speed policy.

A failure while resolving an explicit character engine/voice returns an error through the existing managed audio path. It does not fall back to the previous character, a different engine, or the browser TTS path.

## 2. Original behavior and source evidence

The following was established before the implementation:

- `frontend/dist/js/audio-service.js` read `character.howl.voice`, merged it with global audio settings, and called `SynthesizeSpeechWithVoice` when a character voice configuration existed. It did not include a character identity field.
- `app-audio.go` parsed the override JSON. Malformed JSON was non-fatal in the original path, and an empty/blank override could degrade to a normal synthesis request.
- `Pipeline.ProcessAudioMessage` called `ensureTTSVoice` for an override voice. That helper always loaded a Piper model and installed the new synthesizer into `p.ttsSynth`/`p.loadedVoice`, regardless of the requested engine.
- `VoiceOverride.Engine` was parsed but not used to select an engine.
- `SynthesisRequest` had no engine, character identity, or Piper model-path field.
- Piper's `runPiper` used the synthesizer's initialized `s.modelPath` and ignored `req.Voice`.
- Kokoro passed `req.Voice` to its bridge as `--voice`.
- Melo passed a normalized speaker ID as the request's speaker and normalized language, but the pipeline did not carry a selected language into the request.
- The global persisted `ttsVoice` was not included in `DesiredSettings.Backend()`, `audio.Settings`, or `PipelineConfig`; normal `SynthesizeSpeech` used the synthesizer's first voice.
- The TTS cache key is `text|voice|speed` (`internal/audio/types/utils.go`). Cache design was explicitly outside this stage.

This is a source-derived account of the pre-change path; it is not an acoustic measurement.

## 3. Resulting request flow

```text
character.howl.voice
  -> AudioService.playTTS
  -> SynthesizeSpeechWithVoice(text, JSON)
  -> parseVoiceOverrideJSON (strict)
  -> VoiceOverride { character_id, engine, voice, speed }
  -> ProcessAudioMessage
  -> resolveTTSRequest
     -> active global binding for blank/auto engine
     -> separate per-engine synthesizer for an explicit engine
     -> Piper VoicePath, Kokoro Voice, or Melo Voice + Language
  -> synthesizeWithCache
  -> selected engine Synthesize
```

The character identity remains on `VoiceOverride.CharacterID` and is copied to `SynthesisRequest.CharacterID` for the selected engine call. It is also included in resolution error messages as `character "... voice: ..."`.

### Inheritance and explicit selections

- **Absent override:** `VoiceOverride == nil` inherits the active global engine, the configured global `TTSVoice` when present, the established global speed policy, and the first voice when no global voice is set.
- **Explicit inheritance:** a non-nil override with blank or `auto` `Engine` and blank `Voice` follows the same global path, while retaining `CharacterID`.
- **Explicit engine/voice:** a non-blank, non-`auto` engine is validated and loaded or reused independently. An explicit character voice is never replaced because the active engine is different.
- **Explicit engine with blank voice:** the selected engine's first/default voice is used. The global Piper voice is not forced onto a different engine.
- **Invalid configuration:** unsupported engine, renderer-only `browser`, unavailable Piper model file, or unsupported voice returns an error. No fallback synthesis is attempted.

Blank/`auto` engine and explicit inheritance are intentionally represented differently in the renderer so an explicit inheritance request still carries a character identity.

## 4. Engine-specific behavior

### Piper

`resolveTTSRequest` resolves `voiceModelPath(voice)` and checks the model file before synthesis. The path is carried as `SynthesisRequest.VoicePath`. `runPiper` now prefers `req.VoicePath` over the synthesizer's `s.modelPath`. This keeps Piper's model-file semantics without replacing the active Piper synthesizer.

The selected engine map key includes the Piper model path (`piper|<path>`), so two Piper requests with different model files do not share a synthesizer binding.

### Kokoro

The existing Kokoro bridge contract is preserved: `req.Voice` is passed as `--voice` by `runKokoro`. The resolver passes the requested voice through the request; it does not install a new active Kokoro synthesizer.

### MeloTTS

The existing Melo contract uses a speaker ID and a language hint. The resolver passes the voice as `Voice` and, for Melo, also sets `Language` to the selected voice/speaker ID. `normalizeMeloLanguage` therefore receives the speaker's language family instead of defaulting every selection to English. The HTTP adapter continues to send the normalized speaker and language parameters.

Melo's supported speed guard remains enforced when the request does not provide an override speed.

## 5. Shared-state protection

The pipeline now has:

- the existing active global binding: `p.ttsSynth`;
- the active global engine name: `p.activeTTSEngine`;
- a request-only engine map: `p.ttsEngines`, protected by `p.mu`;
- `p.ttsLoadMu`, which serializes explicit engine loading.

For an explicit character engine different from the active engine, `loadTTSForRequest` loads the requested provider into the request engine map. It never assigns `p.ttsSynth`, `p.activeTTSEngine`, or `p.config`. A failure leaves those fields untouched.

A request for a different character therefore cannot change the next character's active engine. `Shutdown` clears the request engine map; resources loaded through the existing runtime loader remain tracked and are released by the existing shutdown path.

Group 2A speed policy is preserved by using the configured global engine, not the automatic provider selected at initialization, when deciding whether global speed applies. In particular, an `auto` global selection continues to use the established `1.0` backend speed. `VoiceOverride.Speed > 0` still takes precedence.

## 6. Global voice defaults

`TTSVoice` is now part of:

- `DesiredSettings.Backend()`;
- `audio.Settings` and the production factory;
- `types.PipelineConfig`;
- the settings report as a backend-applied setting when the active TTS generation matches;
- the inheritance branch of `resolveTTSRequest`.

A no-override request and an explicit inheritance request use the configured global voice when it is non-empty. If no global voice is configured, the active synthesizer's first voice remains the established fallback. Renderer characters without an enabled per-character voice remain `nil` overrides, preserving the existing no-override call path and browser fallback behavior.

## 7. Failure behavior

- Malformed voice-override JSON now returns a parsing error from `SynthesizeSpeechWithVoice`.
- Unsupported engines and `browser` return explicit backend errors.
- A missing Piper model file returns a model-not-found error before synthesis.
- An unsupported non-Piper voice is rejected when the adapter enumerates voices.
- A failed explicit engine load does not change the active binding; the next request resolves from the unchanged global state.
- The renderer does not silently switch a configured character request to browser TTS after a backend synthesis error.

## 8. Tests and exact results

### New focused tests

`internal/audio/pipeline_character_voice_test.go` covers:

1. Character A selects Kokoro and `af_heart`.
2. Character B selects Melo and `ZH`.
3. Sequential A -> B -> A resolves each engine/voice correctly.
4. A no-override request inherits the active Piper default and the configured global voice.
5. Explicit inheritance behaves as inheritance, not as an invalid empty voice.
6. Unsupported `browser`, unknown engine, and unknown Kokoro voice fail without fallback.
7. Engine-specific parameters: Piper `VoicePath`, Kokoro voice, and Melo voice + language.
8. A failed request does not leave the next character on the failed configuration.
9. The existing no-override path still works.
10. Group 2A speed precedence and the `auto` speed policy remain correct.
11. Concurrent explicit engine requests reach the correct fake engine with the correct character ID and leave the active Piper binding unchanged.

`app_audio_voice_test.go` covers strict malformed JSON rejection, empty payload as no override, explicit inheritance preservation, and character identity preservation.

### Existing/regression tests run

```text
go test ./... -count=1
```

Result: **PASS**. This includes the root app tests, all `internal/audio/...` packages, backend packages, and image packages.

```text
go test -race ./internal/audio -run 'TestCharacter|TestPiperCharacter|TestConcurrentCharacter' -count=1
```

Result: **PASS** (no race reports).

```text
go test ./internal/audio -run 'TestCharacter|TestPiperCharacter|TestConcurrentCharacter|TestPipeline|TestSettingsReport|TestProductionFactory' -count=1
```

Result: **PASS**.

```text
go test . -run 'TestParseVoiceOverride|TestAllAudioSettings|TestFailedTTSSpeed|TestSettingsReport' -count=1
```

Result: **PASS**.

```text
node --check frontend/dist/js/audio-service.js
```

Result: **PASS**.

Scoped `git diff --check` for the Stage 3 Go/JS files: **PASS**. The full worktree contains unrelated pre-existing trailing whitespace and modified files; those were not changed for this stage.

### Not tested

- Real Piper/Kokoro/Melo acoustic output or voice identity by listening.
- Real engine startup with a full model set; the focused tests use `scriptedTTS` fakes. One Piper test uses a temporary stand-in model file to verify `VoicePath` propagation, not acoustic synthesis.
- Frontend Playwright suites; no frontend test suite was run for this stage.
- Cache correctness across engines for identical voice strings; see limitations.
- Long-running resource pressure or repeated model-load timing.

## 9. Files changed for this stage

Production/runtime:

- `app-audio.go` — strict override parsing; character identity and error propagation.
- `frontend/dist/js/audio-service.js` — sends character identity; preserves blank engine inheritance; avoids cross-engine global voice substitution; no browser fallback for configured character failures.
- `internal/audio/pipeline.go` — request resolution, per-engine map, explicit engine loading, Piper path, Melo language, global voice, speed policy.
- `internal/audio/types/interfaces.go` — `TTSVoice`, `SynthesisRequest.Engine`, `CharacterID`, `VoicePath`, and `VoiceOverride.CharacterID`.
- `internal/audio/tts/piper/local.go` — request-scoped Piper model path.
- `internal/audio/manager.go` — operational `TTSVoice` setting.
- `internal/audio/manager_factory.go` — forwards `TTSVoice` to `PipelineConfig`.
- `internal/audio/settings_contract.go` — globally applied `ttsVoice` and report status.

Tests:

- `internal/audio/pipeline_test.go` — configurable fake voices used by the new tests.
- `internal/audio/pipeline_character_voice_test.go` — new Stage 3 coverage.
- `app_audio_voice_test.go` — new Wails-boundary parsing coverage.
- `app_audio_contract_test.go` — updated backend-applied setting count and global voice assertion.
- `internal/audio/settings_contract_test.go` — updated report count for the newly supported `ttsVoice`.

## 10. Limitations and review notes

1. **Cache identity remains text/voice/speed.** `GenerateCacheKey` was not changed, per the Stage 3 cache boundary. A same text, same voice string, and same speed requested through two different engines could therefore share a cache entry. This was not exercised because cache-key work is explicitly out of scope; it is a Stage 4 decision.
2. **Explicit engine instances persist until shutdown.** The per-engine map avoids active-state mutation but does not add per-request unloading or a resource-pressure policy.
3. **Legacy `ensureTTSVoice` remains.** It is retained for existing pipeline regression tests and the old lazy Piper helper contract, but it is no longer called by `ProcessAudioMessage`. It still mutates the active Piper binding if invoked directly; it is not part of the character request path.
4. **Global voice is engine-scoped in practice.** A persisted global voice that is valid for Piper is not automatically valid for Kokoro or Melo. The resolver fails accurately when the active/requested engine does not advertise that voice rather than silently choosing another one.
5. **No acoustic verification.** Engine routing and request parameters are verified; actual voice quality and whether the engine's native voice model produces the expected timbre are not.

## 11. Protected boundaries confirmation

- Stage 1 character identity/persistence and Stage 2 configuration integrity files were not redesigned or reworked in this stage.
- Group 2A speed behavior remains protected by existing tests and the new precedence/auto-policy tests.
- No cache-key redesign, cache-size work, performance benchmark, VRAM/RAM measurement, model discovery, automatic download, background loading, Character Tab redesign, chat-bubble menu, or unrelated audio-settings feature was added.
- The only renderer changes are request routing for character voice data and the existing failure handling around that path.

No commit or push was performed.
