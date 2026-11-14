# AUDIT — Audio settings contract and field capability map

**Date:** 2026-11-05  
**Repository:** `C:\Users\StoneWolfPC\HOWL_Chat`, branch `master`, working tree  
**Mode:** read-only source audit. No application source, settings, or test files were changed. This document does not authorize implementation.  
**Scope:** The 34 persisted `audio` keys in the settings panel, with particular attention to the 28 keys that have no backend effect, and the current Saved / Applied / Ready contract.

---

## 1. Bottom line

The 28-field count is accurate but the contract problem is larger than “some controls are unwired.”

- The panel persists **34** audio keys.
- Only **4** reach the audio backend: `asrEnabled`, `ttsEnabled`, `asrEngine`, `ttsEngine`.
- **2** are effective only in the renderer: `ttsVolume`, `ttsAutoplayMessage`.
- **28** are saved but have no effect on what the backend applies or on readiness.

The important consequence is not just that those 28 controls do nothing. It is that **“Applied” is currently a four-field claim while “Saved” is a 34-field claim**:

- `audio-settings.js` snapshots only the four capability fields for apply.
- `ApplyAudioSettings` accepts only those four values.
- `audio.Settings` and `AppliedSettings()` contain only those four values.
- `GetAudioStatus().applied` publishes only those four values.
- `_confirmApplied()` compares only those four values.

Therefore a user can change `asrLanguage`, `ttsSpeed`, `micSensitivity`, or any other one of the 28, save successfully, and still see **“Saved” + “Applied”** because the unchanged four capability fields still match. The field never left the frontend, but the panel’s aggregate outcome does not distinguish that from a field that was genuinely applied.

That is the contract defect to fix before wiring more controls. One of two things must become true:

1. **Expand the applied contract** to every supported field, with an explicit applied value and comparison for each; or
2. **Scope the claim** so the panel says “Capabilities applied” and does not imply that the entire saved audio object is effective.

The 28 controls should not be wired one by one into the existing four-field contract. The contract has to grow first.

The 28 divide as follows:

| Class | Count | Meaning |
|---|---:|---|
| Supported now, plumbing/validation only | 6 | The provider already accepts the value; the typed contract and UI-to-backend path are missing. |
| Implementation required | 6 | 4 backend/provider additions and 2 renderer behaviors. |
| Product decision required first | 16 | The current label has no truthful implementation, or the feature is absent/out of scope. |

The 6 supported fields are not all supported by every engine. Engine-specific applicability is part of the map; a single global control must not imply support on a provider that ignores it.

---

## 2. What Saved, Applied, and Ready mean today

### 2.1 Saved

`frontend/dist/js/settings-state.js` holds one `currentSettings` object. Its `audio` default block contains all 34 keys. `saveSettings()` writes the whole object to `localStorage` immediately and schedules a 150 ms debounced backend write. `saveSettingsNow()` / `flushSettings()` clears the timer and awaits the backend write, so the audio panel can prove the file write before applying.

`app-chat.go:SetSettings` stores the entire map in `app.appSettings` and writes it through `internal/settings/store.go`. `SetSettings` mirrors only the `sampler` and `model` sections into their legacy in-memory maps; it does not project the `audio` section into `audio.Settings`.

**Saved meaning:** the whole audio object, including all 28, reached `settings.json` (subject to the normal write-failure handling). It does not mean any of it was applied.

### 2.2 Applied

The apply path is deliberately narrow:

- `audio-settings.js:_audioSnapshot()` returns `asrEnabled`, `ttsEnabled`, `asrEngine`, `ttsEngine`.
- `app-audio.go:ApplyAudioSettings(...)` accepts those four and builds `audio.Settings` with those four fields.
- `internal/audio/manager.go:Settings` has only those four fields.
- `Manager.AppliedSettings()` returns the active generation’s four-field settings.
- `app-audio.go:GetAudioStatus()` publishes `applied` with `asr_enabled`, `tts_enabled`, `asr_engine`, `tts_engine` only.
- `audio-settings.js:_confirmApplied()` compares only those four.

**Applied meaning:** the active audio generation matches the four chosen capability fields. It does not mean the panel’s other 30 persisted values are effective, and it does not include the two renderer-only values either.

### 2.3 Ready

`GetAudioStatus()` derives readiness from the manager’s ASR/TTS capability state, the microphone probe, and the selected ASR engine’s local/remote classification. The UI shows `Ready (local initialized)` only for a local engine that initialized, and the microphone message takes precedence when capture is unusable.

**Ready meaning:** the capability/provider is initialized and usable to the extent stated. It is not a per-field validation. None of the 28 fields can make a capability ready, and none of them have an independent ready/error state.

### 2.4 Why the outcome pill overclaims

`_showApplied()` sets “Applied” when `appliedConfirmed === true && applyInProgress === false`. `appliedConfirmed` comes from comparing the four-field snapshot with the four-field `status.applied`. It never asks whether the 28 changed fields were applied.

Examples that currently show “Saved” and “Applied” despite no backend effect:

- changing Spoken Language from `auto` to `en`;
- changing Speech Speed from `1.0` to `1.35` with no character voice override;
- changing CPU Threads;
- changing Mic Sensitivity;
- changing any VAD or device/format control.

This is the main reason to treat the Saved / Applied / Ready contract as a prerequisite, not as cosmetic copy.

---

## 3. Capability map

Legend:

- **SUPPORTED** — the provider or backend already accepts this value; the app needs typed plumbing and validation. Engine-specific limits are stated.
- **IMPLEMENT** — backend/provider work is required (or frontend work, marked separately).
- **DECIDE** — the current control has no truthful implementation. Choose remove/disable, or define the feature before coding.
- **FRONTEND** — no backend field is needed; the behavior belongs in the renderer.

### 3.1 Speech recognition (9)

| # | Persisted key / UI | Current evidence | Class | Disposition |
|---:|---|---|---|---|
| 1 | `asrLanguage` — Spoken Language | `ASRConfig.Language` exists and Whisper appends `--language` when it is not `auto`. `Pipeline.initASR` currently hardcodes `Language: "auto"` and initializes the fallback chain with `Language: "auto"`. Qwen remote and local llama ASR do not consume a language setting. | **SUPPORTED (Whisper only)** | Add the value to `audio.Settings` / `PipelineConfig`, pass it into `ASRConfig`, validate it against the selected profile, and disable/annotate it for engines that do not honor it. |
| 2 | `asrDevice` — Compute Device | `ASRConfig.Device` exists but no provider reads it. Whisper does not pass a device flag; local llama ASR hardcodes `-ngl 99`; Qwen remote has no local device. The UI’s Auto/CPU/CUDA promise is not backed by any adapter path. | **DECIDE** | Remove the control, or redefine it as a per-engine acceleration mode with tested semantics. Do not advertise CUDA merely because the field exists. |
| 3 | `asrThreads` — CPU Threads | `ASRConfig.NumThreads` is used by Whisper (`--threads`) and local llama (`-t`). `Pipeline.initASR` hardcodes `NumThreads: 6`. Remote Qwen has no local thread count. | **SUPPORTED (local engines)** | Pass through for Whisper and local llama; disable for remote Qwen. Validate the advertised 1–16 range. |
| 4 | `asrBeamSize` — Beam Size | `ASRConfig.BeamSize` exists and Whisper appends `--beam-size` when positive. The pipeline never sets it. Qwen profiles use beam size 1 and remote Qwen does not accept this control. | **SUPPORTED (Whisper only)** | Pass through for Whisper; disable/annotate for engines where beam search is not applicable. |
| 5 | `asrBestOf` — Best-of | `ASRConfig.BestOf` exists and Whisper appends `--best-of` when positive. The pipeline never sets it. Qwen profiles use best-of 1. | **SUPPORTED (Whisper only)** | Same as Beam Size. |
| 6 | `asrTemperature` — Temperature | `ASRConfig.Temperature` exists, but `whisper.runWhisperCLI` never appends a temperature argument. Whisper CLI supports temperature control; the adapter has no path for it. | **IMPLEMENT** | Add a tested Whisper argument and validation (the UI says 0–1). Scope to Whisper; do not send it to providers that do not decode with it. |
| 7 | `asrTranslateToEnglish` — Translate to English | The Whisper JSON parser reads `params.translate`, but the recognizer never requests `--translate`. Qwen/llama ASR do not expose this option. The existing design note D5 proposes disabling unsupported decode flags with a reason. | **DECIDE** | Either disable with a truthful reason, or implement Whisper-only translation and rename the control to say “Whisper only.” |
| 8 | `asrNoTimestamps` — Suppress Timestamps | Whisper’s parser collects segments, but the pipeline returns only `TranscribedText` and the UI never displays timestamps. Setting this true or false produces the same text. | **DECIDE** | Remove the control or expose segment timestamps first. Wiring a flag that changes no observable output would create a false success. |
| 9 | `asrFlashAttention` — Flash Attention | No ASR field or argument exists. Whisper CLI supports flash attention, but it is hardware/backend dependent and the existing design note proposes disabling it. | **DECIDE** | Disable with a reason, or implement a Whisper-only, capability-probed option. Do not leave a global toggle that silently does nothing. |

### 3.2 Text-to-speech (7)

| # | Persisted key / UI | Current evidence | Class | Disposition |
|---:|---|---|---|---|
| 10 | `ttsVoice` — TTS Voice Model | `SynthesisRequest.Voice` exists, `SynthesizeSpeechWithVoice` accepts a per-request override, and `ensureTTSVoice` can lazily load a Piper voice. But normal `playTTS` calls `SynthesizeSpeech` unless a character voice override is enabled, so the global selection is not sent. `ensureTTSVoice` always loads `types.TTSPiper`, even when the selected engine is Kokoro or MeloTTS. | **SUPPORTED (partial)** | Send the global voice through the backend and make voice switching engine-aware. Validate installed voices; do not silently use the first voice or switch a Kokoro/Melo session to Piper. |
| 11 | `ttsSpeed` — Speech Speed | `SynthesisRequest.Speed` is used by Piper, Kokoro, and MeloTTS, and the browser path uses the global value. `synthesizeWithCache` hardcodes `Speed: 1.0`; only a character override changes it. | **SUPPORTED** | Add the value to the applied settings and pass it into `SynthesisRequest`. It is already part of the cache key. |
| 12 | `ttsNoiseScale` — Noise Scale | `VoiceOverride` carries `NoiseScale`, but `synthesizeWithCache` copies only voice and speed. `SynthesisRequest` has no noise field. Piper CLI supports `--noise_scale`; Kokoro/Melo do not expose the same control. | **IMPLEMENT (Piper only)** | Add the field to the request and Piper argument builder, or remove/hide this advanced Piper-only control. |
| 13 | `ttsNoiseW` — Noise Width | Same path as Noise Scale; Piper CLI supports `--noise_w`, but no request field or argument exists. | **IMPLEMENT (Piper only)** | Same as Noise Scale. |
| 14 | `ttsSentenceSilence` — Sentence Silence | `VoiceOverride` carries `SentenceSilence`, but the pipeline ignores it and `SynthesisRequest` has no field. Piper CLI supports `--sentence_silence`. | **IMPLEMENT (Piper only)** | Add the field and Piper argument, or remove the control. |
| 15 | `ttsCacheEnabled` — Enable TTS Cache | The pipeline always creates `cache.NewManager(100)` and always uses it. `TTSConfig.CacheEnabled` is hardcoded true and no provider reads it. The manager is an **in-memory LRU of `SynthesisResult` metadata**, not an audio-byte disk cache; it does not evict files from disk. | **DECIDE** | Either implement the disk cache the UI describes (bytes, eviction, persistence) or remove the switch and correct the copy. Do not wire a boolean onto an always-on metadata cache. |
| 16 | `ttsCacheSizeMB` — Cache Size Limit | `cache.NewManager(100)` is hardcoded; `MaxCacheSize: 100` is hardcoded in every `ttsConfigForModel` branch. The manager’s limit applies to estimated metadata size, not disk files. | **DECIDE** | Same decision as Cache Enabled. If a real disk cache is built, make this the actual byte/MB eviction limit and validate it. |

### 3.3 Voice activity detection (5)

| # | Persisted key / UI | Current evidence | Class | Disposition |
|---:|---|---|---|---|
| 17 | `vadEnabled` — Enable VAD | Only a `types.VADEngine` interface and `VADConfig` exist. No implementation is registered and the pipeline never invokes VAD. The architecture notes list VAD as a pending stage. | **DECIDE** | Disable/remove the VAD tab with a reason, or implement VAD as a complete feature. Do not wire one toggle into no engine. |
| 18 | `vadThreshold` — Speech Threshold | No VAD consumer. | **DECIDE** | Same VAD group decision. |
| 19 | `vadMinSpeechMs` — Min Speech | No VAD consumer. | **DECIDE** | Same VAD group decision. |
| 20 | `vadMinSilenceMs` — Min Silence | No VAD consumer. | **DECIDE** | Same VAD group decision. |
| 21 | `vadPaddingMs` — Speech Padding | No VAD consumer. | **DECIDE** | Same VAD group decision. |

### 3.4 Capture and device format (6)

| # | Persisted key / UI | Current evidence | Class | Disposition |
|---:|---|---|---|---|
| 22 | `audioInputDevice` — Input Device | The settings panel populates the select from `navigator.mediaDevices.enumerateDevices()`, but `getUserMedia` does not pass a `deviceId`. The Windows preflight probes the console-default endpoint, not the selected device. The existing design note calls browser device selection independent/out of scope. | **DECIDE** | Remove/hide the selector, or define how a chosen browser device interacts with the default-endpoint WASAPI preflight. A selected device that differs from the probe would make the gate and capture disagree. |
| 23 | `audioOutputDevice` — Output Device | No `setSinkId` call exists anywhere. Backend TTS returns a file and frontend playback uses `new Audio(path)` on the system default. Browser TTS uses the system default as well. The audit’s headphone-only requirement is not implemented. | **DECIDE** | Decide the headphone/output policy first: explicit endpoint, endpoint-type preference, or unsupported with a clear message. Then implement `setSinkId` where the WebView supports it; do not silently fall back to speakers. |
| 24 | `audioSampleRate` — Sample Rate | Capture requests a 16 kHz `AudioContext` (falling back to the native rate), and `AudioWav.encodeWav16` always resamples to 16 kHz and writes a 16 kHz header. The 22.05/44.1/48 kHz options are never read. | **DECIDE** | Remove the options, or implement an end-to-end rate policy and prove the ASR path accepts the resulting WAV. The current path is fixed 16 kHz. |
| 25 | `audioChannels` — Audio Channels | `getUserMedia` requests `channelCount: 1`; the WAV encoder writes mono; ASR consumes mono. The mono/stereo control is not read. | **DECIDE** | Remove the control, or implement stereo capture plus a documented downmix. ASR does not need stereo today. |
| 26 | `audioBitDepth` — Bit Depth | `AudioWav.encodeWav16` always writes 16-bit PCM. The 24/32-bit options have no encoder and no consumer. | **DECIDE** | Remove the control, or implement additional encoders and prove the ASR path accepts them. 16-bit is the current contract. |
| 27 | `micSensitivity` — Mic Sensitivity | No reader in the capture path. No `GainNode` is inserted between source and processor. | **FRONTEND / IMPLEMENT** | Add a bounded gain stage and define the 0–1 mapping and clipping behavior, or remove the slider. No backend field is required. |

### 3.5 Input mode (1)

| # | Persisted key / UI | Current evidence | Class | Disposition |
|---:|---|---|---|---|
| 28 | `pushToTalk` — Push-to-Talk | No reader. The mic button calls `toggleRecording()` on click; there are no pointer-down/up/cancel handlers, no minimum-duration path, and the tooltip is static. The existing design note D6 describes pointer capture, blur/visibility cancellation, and a minimum duration. | **FRONTEND / IMPLEMENT** | Implement the pointer/keyboard behavior in the renderer and update the tooltip. No backend field is involved. |

---

## 4. Count reconciliation

| Group | Keys |
|---|---:|
| Backend-applied | 4 |
| Renderer-effective, not backend | 2 |
| No backend effect | 28 |
| **Total persisted audio keys** | **34** |

The 28 are exactly the audit’s earlier grouping:

- 9 ASR decode controls;
- 5 VAD controls;
- 6 device/format controls;
- 7 TTS voice/output/cache controls;
- 1 Push-to-Talk control.

No additional persisted audio key was found outside that set.

---

## 5. Recommended implementation groups

These are proposed boundaries, not authorized work.

### Group 0 — Contract and outcome truthfulness (prerequisite)

- Define a typed audio settings contract that can represent every field actually supported.
- Send the full desired set in one apply request.
- Publish an applied value for every supported field.
- Compare the applied value per field, or change the aggregate label to “Capabilities applied.”
- Add validation with per-field errors before any provider is loaded.
- Keep Saved independent: a failed apply must not erase the saved selection.

Without Group 0, wiring individual controls will keep producing an “Applied” pill that does not mean what it says.

### Group 1 — ASR Whisper decode (bounded, provider-supported)

`asrLanguage`, `asrThreads`, `asrBeamSize`, `asrBestOf` are supported by existing Whisper paths. Add typed fields, pass them to `ASRConfig`, and test the generated `whisper-cli` arguments. Gate the controls to Whisper/local engines and disable them for providers that do not consume them.

`asrTemperature` belongs in the next slice only if it is kept; it needs a new Whisper argument. `asrTranslateToEnglish`, `asrNoTimestamps`, and `asrFlashAttention` should remain disabled until their product decision is made.

### Group 2 — TTS output (bounded)

`ttsSpeed` is already supported by every backend synthesizer; wire it globally and confirm the cache key includes it. `ttsVoice` needs global routing and engine-aware voice switching before it is truthful. Piper noise/width/sentence-silence belong in a separate Piper-only advanced slice, not the general TTS contract.

TTS cache enable/size must wait for the disk-cache decision. The current cache is an in-memory metadata cache and cannot honor “Enable cache” or “Cache Size Limit” as described.

### Group 3 — Capture UX (frontend-only)

`pushToTalk` and `micSensitivity` can be implemented without changing `audio.Settings`. They still need real interaction tests: pointer cancellation, blur/visibility, minimum duration, gain clamping, and no recording left active.

### Group 4 — Decisions and removal

Do not implement these as if they were ordinary fields:

- VAD (five controls) — remove/disable or build the engine.
- Compute Device — define per-engine semantics or remove.
- Audio input/output devices — define the selected-endpoint and headphone policies first.
- Sample rate, channels, bit depth — remove or implement end-to-end capture/encode/resample support.
- Unsupported ASR flags — disable with reasons unless the product explicitly wants Whisper-only behavior.
- Cache enable/size — implement a real disk cache or remove the controls.

---

## 6. Product decisions required before coding

| ID | Decision | Why it cannot be settled by plumbing |
|---|---|---|
| D1 | VAD: implement or remove/disable the tab? | There is no VAD engine and no pipeline call. |
| D2 | ASR compute device: remove, or define a tested per-engine acceleration mode? | `ASRConfig.Device` has no consumer; local llama hardcodes GPU layers. |
| D3 | Translation, no-timestamps, flash attention: keep any of them? | Provider support is Whisper-specific; no-timestamps has no observable output today. |
| D4 | TTS cache: real disk cache with eviction, or remove the switch and correct the copy? | The current manager stores metadata in memory and never evicts disk files. |
| D5 | Piper noise/width/sentence silence: expose as Piper-only advanced controls, or remove? | These are real Piper options but not part of the generic TTS request. |
| D6 | Input device: remove, or align browser capture with a selected endpoint and the WASAPI preflight? | The probe currently checks the console default, while capture could target another device. |
| D7 | Output device / headphone policy: explicit endpoint, endpoint-type preference, or refusal? | No `setSinkId` exists; the fallback is always system default. |
| D8 | Sample rate/channels/bit depth: remove, or commit to full capture/encode/resample support? | The current renderer is fixed at 16 kHz mono 16-bit. |
| D9 | Applied pill: expand the applied contract or rename the aggregate claim? | The pill currently compares only four fields. |

---

## 7. Evidence boundaries

- This audit is source inspection only. It did not run the 28 controls, synthesize audio, capture a microphone, or exercise VAD.
- It did not alter the real `%APPDATA%\HOWLChat` profile or any settings file.
- The existing lifecycle tests cover persistence and the four capability fields; they do not cover the 28 mapped keys.
- The renderer paths cited are the production files under `frontend/dist/`; no build step was run.
- No code change, test change, commit, or push was made.

---

## 8. Source references

Primary files inspected:

- `frontend/dist/js/settings-state.js` — 34-key default audio schema, deep merge, persistence, `flushSettings`.
- `frontend/dist/js/audio-settings.js` — hydration, four-field `_audioSnapshot`, `_confirmApplied`, outcome pills.
- `frontend/dist/audio-settings.html` — the controls and their labels.
- `frontend/dist/js/audio-service.js` — `getUserMedia` capture, 16 kHz context, TTS playback, autoplay.
- `frontend/dist/js/audio-wav.js` — fixed 16 kHz mono 16-bit WAV encoder.
- `app-chat.go` — `GetSettings` / `SetSettings`, whole-object persistence.
- `internal/settings/store.go` — `settings.json` load/save.
- `app-audio.go` — `ApplyAudioSettings`, `GetAudioStatus().applied`, microphone gate.
- `internal/audio/manager.go` — four-field `Settings`, `AppliedSettings`.
- `internal/audio/pipeline.go` — hardcoded ASR language/threads, `synthesizeWithCache`, always-on cache, TTS config.
- `internal/audio/types/interfaces.go` — `ASRConfig`, `TTSConfig`, `SynthesisRequest`, `VoiceOverride`, `VADConfig`.
- `internal/audio/asr/whisper/recognizer.go` — Whisper arguments (language, beam, best-of, threads).
- `internal/audio/asr/whisper/parser.go` — timestamps parsed into segments but not surfaced.
- `internal/audio/asr/llama/llama.go` — local ASR threads and hardcoded GPU layers.
- `internal/audio/tts/piper/local.go` — Piper speed; no noise/width/sentence-silence request fields.
- `internal/audio/tts/kokoro/local.go`, `internal/audio/tts/melotts/client.go` — engine-specific speed/voice behavior.
- `internal/audio/tts/cache/manager.go` — in-memory metadata LRU, no disk byte eviction.
- `docs/asr-repair-design.md` §D5–D6 — prior design proposals for unsupported flags, VAD, input devices, and push-to-talk.
