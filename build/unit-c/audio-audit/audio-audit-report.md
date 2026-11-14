# Chat page — Audio audit (read-only)

Date: 2026-11-05. Isolated profile: `build/tmp/unit-c-rerun` (test settings, not the real
`%APPDATA%\HOWLChat`). No product code was changed for this audit.

Evidence directory: `build/unit-c/audio-audit/`

| Artifact | What it shows |
| --- | --- |
| `audio-audit.log` | persistence pass: seed → save → relaunch restore, 1920×1080 resize, light theme |
| `audio-audit2.log` | 1280×800 layout pass, tab geometry, offscreen counts |
| `general-tree.txt`, `restore-tts-tree.txt`, `1280-*.txt` | UIA trees with control ids and rects |
| `tab-general.png`, `tab-vad.png`, `tab-hardware.png`, `tab-textto-speech.png` | rendered tabs at 1400×900 |
| `tts-1920.png`, `tts-1920-light.png` | 1920×1080 and light-theme rendering |
| `1280-general.png`, `1280-asr.png`, `1280-hardware.png`, `1280-vad.png` | 1280×800 rendering |
| `build/tmp/unit-c-rerun/config/settings.json` | persisted audio block before/after (logs record each state) |

## 1. Component inventory

Wiring shorthand: **B** = Wails binding, **S** = settings.json via `GetSettings`/`SetSettings`,
**L** = localStorage mirror, **F** = frontend-only (no backend read).

### 1.1 Menus and entry points

| # | Component | Intended | Implementation | Status |
| --- | --- | --- | --- | --- |
| 1 | Settings → Audio Settings | open panel on General | `openAudioSettings()` → iframe `audio-settings.html#general` | VERIFIED (opened in every pass) |
| 2 | Settings → TTS Settings | open panel on TTS tab | `openTTSSettings()` → same iframe with `#tts` | PARTIAL — same panel, hash works by code; not separately exercised |
| 3 | Mic button (input bar) | record → transcribe → fill draft | `AudioService.toggleRecording()` → `getUserMedia` → `TranscribeAudioData` (B) | PARTIAL — capture path exists; not exercised (no active mic on this machine) |
| 4 | TTS button (input bar) | speak composer text | `AudioService.playTTS(value)` → `SynthesizeSpeech` (B) | PARTIAL — silently returns when TTS is off; no feedback |

### 1.2 General tab

| # | Control | Intended | Wiring | Persist | Applied | Status |
| --- | --- | --- | --- | --- | --- | --- |
| 5 | Enable ASR `asr-enabled` | request ASR capability | `ApplyAudioSettings` (B) | S | yes | VERIFIED (initialized `qwen3-asr-1.7b-local`) |
| 6 | Enable TTS `tts-enabled` | request TTS capability | `ApplyAudioSettings` (B) | S | yes | PARTIAL — apply path verified by code; no installed Piper binary to prove synthesis |
| 7 | Auto-play TTS `tts-autoplay-message` | speak AI replies | `chat:complete` hook in `audio-service.js` | S | yes (frontend) | VERIFIED by code; not run end-to-end |
| 8 | Push-to-Talk `push-to-talk` | hold vs toggle | **no reader anywhere** | S | no | UNWIRED — mic button tooltip always says "Hold to record"; click always toggles |
| 9 | TTS volume `tts-volume` | playback gain | `new Audio().volume`, browser `utterance.volume` | S | frontend only | PARTIAL — backend audio file gain not applied |
| 10 | ASR status text | readiness/error | `_asrStatusText(GetAudioStatus)` | — | — | VERIFIED ("ready (local engine initialized)") |
| 11 | TTS status text | readiness/error | `status.tts_ready` | — | — | VERIFIED ("not initialized") |
| 12 | Save & Close | persist + apply | `flushSettings` → `SetSettings` (B) → `ApplyAudioSettings` (B) | S | yes | VERIFIED (ttsSpeed 1.35 written, modal closed, toast shown) |
| 13 | Reset Audio Defaults | restore defaults | Object.assign + flush; `confirm()` | S | no (does not apply) | PARTIAL — not exercised; by code it saves only |
| 14 | Retry activation | re-apply | `ApplyAudioSettings` | — | yes | PARTIAL — not exercised |
| 15 | Outcome pills | saved/applied/ready | `_setOutcome` | — | — | VERIFIED (Saved/Applied shown; separate states rendered) |

### 1.3 Speech Recognition tab

| # | Control | Intended | Wiring | Status |
| --- | --- | --- | --- | --- |
| 16 | ASR Engine | explicit provider | `ApplyAudioSettings` (B), capability contract | VERIFIED — explicit engine honored; failures surface with reason (error log `ASR-MODEL-LOAD`) |
| 17 | Spoken Language | whisper `--language` | **not read**; `initASR` hardcodes `Language: "auto"` | UNWIRED |
| 18 | Compute Device | cuda/cpu | **no reader** | UNWIRED |
| 19 | CPU Threads | `--threads` | **not read**; `initASR` hardcodes `NumThreads: 6` | UNWIRED |
| 20 | Beam Size | `--beam-size` | supports the flag, but value never reaches config | UNWIRED |
| 21 | Best-of | `--best-of` | same | UNWIRED |
| 22 | Temperature | decode temperature | **no reader** | UNWIRED |
| 23 | Translate to English | translate output | **no reader** | UNWIRED |
| 24 | Suppress Timestamps | strip timestamps | **no reader**; whisper JSON timestamps are dropped regardless | UNWIRED |
| 25 | Flash Attention | runtime flag | **no reader** | UNWIRED |
| 26 | Engine availability warning | blockers | `_renderCapabilities` from `GetAudioStatus.capabilities` | VERIFIED |

### 1.4 Text-to-Speech tab

| # | Control | Intended | Wiring | Status |
| --- | --- | --- | --- | --- |
| 27 | TTS Engine | explicit provider | `ApplyAudioSettings` (B) | VERIFIED by code path; no local Piper/Kokoro install to prove synthesis |
| 28 | TTS Voice Model | pick Piper voice | backend only via **character** override; global value ignored by `SynthesizeSpeech` | PARTIAL/UNWIRED globally — list itself is VERIFIED (built-in + registered voices) |
| 29 | + Dir (scan voices) | register folder | `OpenDirectoryPicker` + `ScanTTSDir` (B) | VERIFIED in discovery work (TestVoiceA appeared) |
| 30 | Speech Speed | rate | applied only in browser TTS or character override; `synthesizeWithCache` sets `Speed: 1.0` | PARTIAL/UNWIRED globally |
| 31 | Noise Scale | Piper `--noise_scale` | `VoiceOverride` carries it; pipeline ignores it | UNWIRED |
| 32 | Noise Width | Piper `--noise_width` | same | UNWIRED |
| 33 | Sentence Silence | inter-sentence pause | **no reader** | UNWIRED |
| 34 | Enable TTS Cache | cache switch | `cache.NewManager(100)` always; `CacheEnabled: true` hardcoded | UNWIRED |
| 35 | Cache Size Limit | MB cap | `MaxCacheSize: 100` hardcoded | UNWIRED |

### 1.5 VAD tab

| # | Control | Intended | Wiring | Status |
| --- | --- | --- | --- | --- |
| 36 | Enable VAD | segmentation | only a `VADEngine` interface in `types`; `docs/audio-architecture/README.md` lists `realtime/vad.go` as **Stage 5 PENDING**; never called by the pipeline | UNWIRED |
| 37 | Speech Threshold | 0–1 | same | UNWIRED |
| 38 | Min Speech | ms | same | UNWIRED |
| 39 | Min Silence | ms | same | UNWIRED |
| 40 | Speech Padding | ms | same | UNWIRED |

### 1.6 Hardware tab

| # | Control | Intended | Wiring | Status |
| --- | --- | --- | --- | --- |
| 41 | Input Device | choose mic | `enumerateDevices` populates the list; capture calls `getUserMedia({audio:{...}})` **without deviceId** | UNWIRED |
| 42 | Output Device | choose speaker | never read; no `setSinkId` anywhere | UNWIRED |
| 43 | Refresh (×2) | re-enumerate | `navigator.mediaDevices.enumerateDevices` | PARTIAL — not exercised visually; labels require mic permission |
| 44 | Sample Rate | 16k/22k/44.1k/48k | **no reader**; capture forces 16 kHz context | UNWIRED |
| 45 | Audio Channels | mono/stereo | **no reader**; `channelCount: 1` hardcoded | UNWIRED |
| 46 | Bit Depth | 16/24/32 | **no reader**; WAV encoder is 16-bit | UNWIRED |
| 47 | Mic Sensitivity | input gain | **no reader** | UNWIRED |

### 1.7 Workflows and states

| Workflow | Expected | Observed |
| --- | --- | --- |
| Open panel | hydrate from settings, poll status | VERIFIED — controls hydrate, status text renders |
| Save & Close | write settings, apply, report saved/applied/ready | VERIFIED — `ttsSpeed 1.0 → 1.35` written and applied; toast shown; panel closed |
| Close without saving | no write | PARTIAL — `close()` auto-saves by design (documented in code) |
| Relaunch | restore persisted selections | VERIFIED — 1.35 and all UI-only values survived relaunch |
| Explicit engine unavailable | refuse, keep prior generation, explain | VERIFIED by code + error log (`ASR-MODEL-LOAD` reason) |
| Concurrent activation | lock, retry, no double apply | PARTIAL — code paths exist; not exercised |
| Mic unavailable | block ASR allocation with specific message | **BROKEN** — see §2 |
| Device disappears mid-init | recover, release capture | PARTIAL — `releaseCapture()` on error, generic message only |
| Shutdown | drain, no orphan processes | VERIFIED — `LEFTOVER=0` after each pass |

## 2. Microphone and ASR safeguards — the central failure

1. **No Windows microphone probe exists.** No WASAPI/MMDevice enumeration in Go; the only device
   list is the browser's `enumerateDevices()` in the settings iframe, which is never consulted
   before loading ASR.
2. **No gate before ASR allocation.** `ApplyAudioSettings → Manager.Apply → Pipeline.Initialize`
   loads the recognizer with no reference to capture state.
3. **Direct evidence of the violation.** This machine's Windows endpoints currently have **no
   active capture device** (all `AudioEndpoint` capture entries report `Unknown`; only render
   devices are `OK`). The app still reported
   `Status: ready (local engine initialized; transcription not verified)` and toasted
   *“Audio settings saved and applied. ASR is initialized.”* — i.e. the ASR model was allocated
   and a llama-server started with no usable microphone.
4. **Required message is absent.** `"Windows shows a mic isn't selected. Please connect a
   microphone to use this option"` does not exist in the repository.
5. **No state distinction.** Capture failure produces one generic toast
   (`"Microphone access denied or unavailable"`); no separate no-device, disconnected, access
   denied, init-failure, or initialized states.
6. **No recheck on load request.** Availability is never checked at startup or at apply time.
7. **No device-disappearance handling** beyond `releaseCapture()` in the catch block.
8. **Auto-load contradicts the stated rule.** `restorePersistedAudio()` calls
   `ApplyAudioSettings` whenever `asrEnabled`/`ttsEnabled` was persisted, so a saved selection
   silently loads a model (and reserves VRAM) on launch. Persisted *selection* should be restored
   without activation.

## 3. Output routing

- No `setSinkId` and no use of `audioOutputDevice`; playback is `new Audio(path)` on the system
  default endpoint.
- Therefore the requirement “audio playback always routes to headphones, never speakers” is not
  implemented at all. Nothing enforces or even prefers a headphone endpoint.

## 4. UI, theme, resolution

| Finding | Evidence |
| --- | --- |
| Panel is a fixed 900×600 iframe | `RECT 190,31,900,600` at 1280×800; `510,31,900,600` at 1920×1080 |
| Vertically top-aligned, large unused area | top edge y=31 at both sizes; ~170 px empty at 800, ~450 px at 1080 |
| Content clipped behind the sticky footer | cards cut mid-description in `1280-general.png` (Auto-play), `1280-vad.png` (Min Speech), `1280-hardware.png` (Sample Rate), `tts-1920.png` (Speech Speed) |
| Tabs fit at 1280 | General 215→1090; last tab ends at 819 |
| Light theme not applied to the panel | `tts-1920-light.png`: app is light, audio modal stays dark |
| Emoji tab icons degrade in some renderers | UIA name `??? Speech Recognition`; screenshots render them correctly |
| Two-column grid never reflows | `minmax(380px, 1fr)` inside 850 px content; single column below ~800 px — at 1280/1920 it stays one column, wasting width |
| Disabled/unavailable options look identical to working ones | no disabled styling for unwired settings |

## 5. Hover help

- Exactly **one** `title` attribute exists in the whole panel (`+ Dir`).
- 33 of 34 settings have **no hover help**. Visible descriptions are one-line summaries and mostly
  omit defaults, ranges, and consequences (e.g. Beam Size has no default, VAD has no valid ranges,
  Cache Size has no upper-bound rationale).
- Mic button tooltip is static (“Hold to record”) and contradicts the Push-to-Talk setting.

## 6. Defects, ordered by severity

1. **BROKEN — mic safeguard absent** (no probe, no gate, no message, no state distinction,
   no recheck, no hidden VRAM allocation prevention).
2. **BROKEN — silent auto-load** of persisted ASR/TTS on launch.
3. **ABSENT — headphone-only routing / output-device application.**
4. **UNWIRED — 28 persisted settings** have no backend effect (all 9 ASR decode options, all 5
   VAD controls, all 6 device/format controls, 7 TTS voice/speed/noise/cache controls, and
   Push-to-Talk).
5. **BROKEN — Push-to-Talk** does nothing; mic tooltip misleading.
6. **UI — content clipped behind footer** at both tested sizes.
7. **UI — panel not vertically centred**; fixed 900×600 ignores available height.
8. **UI — light theme not propagated** into the audio panel.
9. **UI — no hover help** for 33 settings; no disabled state for unavailable options.
10. **UX — silent no-ops** (TTS icon when disabled; Reset only saves; close auto-saves).

## 7. Proposed implementation sequence (bounded)

Each step is independently testable and keeps backend-first ordering.

1. **Settings contract (backend)** — extend `audio.Settings`/`PipelineConfig` with exactly the
   fields the providers support (`ASRLanguage`, `ASRThreads`, `ASRBeamSize`, `ASRBestOf`,
   `ASRNoTimestamps`, `TTSVoice`, `TTSSpeed`, `TTSNoiseScale`, `TTSNoiseW`, `TTSCacheEnabled`,
   `TTSCacheSizeMB`, `VAD*` only if VAD ships). Unknown/unsupported values must be refused with a
   reason, never silently ignored. Update the frontend to send them.
2. **Mic gate (backend)** — add a Windows capture-endpoint probe with explicit states
   (`no-device`, `disconnected`, `denied`, `init-failed`, `ready`); check it in
   `ApplyAudioSettings` before ASR allocation and re-check on every apply; expose the state in
   `GetAudioStatus`; return a typed error the UI maps to the exact required string. Handle removal
   during initialization by aborting the build and cleaning up.
3. **Remove silent auto-load** — restore persisted selections into the UI and desired state only;
   load on explicit user action (Save/Retry or first use), never at launch.
4. **Output routing** — apply `audioOutputDevice` via `setSinkId` on both backend audio and browser
   TTS; enforce a headphone policy (explicit endpoint or endpoint-type check) with a clear refusal
   instead of falling back to speakers.
5. **VAD decision** — either implement a VAD engine and wire it into the capture→ASR path, or mark
   the VAD tab unavailable with a reason. Do not leave five controls that promise behaviour.
6. **UI completion** — responsive dialog (max-height with a scroll region above the footer,
   vertically centred), theme propagation into the iframe, two-column reflow, disabled styling for
   unavailable settings, and a tooltip for every setting with what/when/range/default/consequence.
7. **Tests** — settings round-trip, apply-field mapping, mic-gate states with a fake probe,
   output-device selection, and UI geometry assertions at 1280×800 and 1920×1080 in light/dark.

## 8. Limits of this audit

- No microphone exists on this machine, so capture success, volume, and transcription were not
  exercised; the no-device path is what was observed.
- No Piper/Kokoro binaries are installed here, so TTS synthesis was not performed; engine apply,
  voice enumeration, and settings→backend paths were verified by code and status only.
- The ASR tab was not visually captured (harness click landed outside the tab at 1280×800);
  its markup, ids, and wiring were verified statically.
- Reset, Retry, and concurrent-activation locks were not exercised.
- The real `%APPDATA%\HOWLChat` profile was not touched; the test profile was restored afterwards.
