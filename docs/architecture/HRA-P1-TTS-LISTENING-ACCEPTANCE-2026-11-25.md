# HRA-P1 — TTS Listening Acceptance (human operator)

## Status: **ACCEPTED — human listening confirmed by the operator (2026-11-25)**

| Field | Value |
|---|---|
| Date opened | 2026-11-25 |
| Operator | The user (human listening operator) |
| Agent role | Install and verify the local engine, report the artifact path/format/size/status, and state how to play it. Never start playback, never change the output device. |
| Automated evidence | `internal/process`, `internal/audio`, `internal/audio/tts/melotts`, `internal/audio/tts_ownership_test.go` pass; the installed engine was exercised through the production provider (see §2). None of this is a quality judgment. |
| Listening result | **Accepted by the operator** — character read-aloud "works fine" with acceptable latency; composer draft playback "it worked!" (the operator's own words, not the agent's judgment) |

## 1. Installed assets (verified 2026-11-25)

| File | Bytes | SHA-256 | Source |
|---|---|---|---|
| `bin/piper.exe` (plus `espeak-ng.dll`, `piper_phonemize.dll`, `onnxruntime.dll`, `onnxruntime_providers_shared.dll`, `espeak-ng-data/`, `libtashkeel_model.ort`) | 509,952 (exe) | `96f3da3811151580073e40bb4dd20eb0fb8115f5f5f76e2fb54282b3edfa5c1f` | `https://github.com/rhasspy/piper/releases/download/2023.11.14-2/piper_windows_amd64.zip` (22,477,236 B, SHA-256 `f3c58906402b24f3a96d92145f58acba6d86c9b5db896d207f78dc80811efcea`) |
| `models/en_US-lessac-medium.onnx` | 63,201,294 | `5efe09e69902187827af646e1a6e9d269dee769f9877d17b16b1b46eeaaf019f` (equals the published Hugging Face LFS digest) | `https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx` |
| `models/en_US-lessac-medium.onnx.json` | 4,885 | `efe19c417bed055f2d69908248c6ba650fa135bc868b0e6abb3da181dab690a0` | same voice directory, `.../en_US-lessac-medium.onnx.json` |

- The model and config agree: `espeak.voice = en-us`, `dataset = lessac`, `num_symbols = 256`, `audio.sample_rate = 22050`, `language.code = en_US`. MODEL_CARD: single speaker, medium quality, 22,050 Hz, “Trained from scratch.”
- Executable format: `PE32+ executable for MS Windows 6.00 (console), x86-64`.
- Distribution: the Piper project moved to `OHF-Voice/piper1-gpl`; that repository now publishes Python wheels only (v1.8.0), so the last official standalone Windows CLI is the archived `rhasspy/piper` release `2023.11.14-2`, which is what this CLI adapter (`bin/piper.exe`, `--model/--output_file/--length_scale`) requires. No mirror was used. GitHub publishes no checksum for that older asset; the archive hash above matches the hash independently quoted for the official zip in OHF-Voice issue #337.
- Licensing: engine GPL-3.0-or-later. Voice rights come from the Lessac Blizzard 2013 dataset (`https://www.cstr.ed.ac.uk/projects/blizzard/2013/lessac_blizzard2013/license.html`); the Piper project makes no separate grant over voice weights.
- The previous project-root `models` entry was a dangling link (`models -> /e/HOWL_Chat/models`, target absent). It was removed without following it and replaced by a real directory; the dead target was not written to.

## 2. Verification performed (no audio played)

| Check | Result |
|---|---|
| Standalone synthesis, text on stdin, temp destination | exit 0 in 1,831 ms; log “Loaded voice in 0.28 s / Initialized piper / Terminated piper” |
| Artifact | `build/tmp/piper-install/verify.wav`, 297,748 B, RIFF/WAVE, mono, 22,050 Hz, 16-bit, 6.75 s |
| Process cleanup | `tasklist` shows no `piper.exe` after exit |
| Production provider (`audio.TTSModelPathFor`, `piper.New`, `Initialize`, `Synthesize`) | resolved `models/en_US-lessac-medium.onnx`; capability `piper` selectable with model state available; synthesis before `Initialize` refused |
| Ownership during synthesis | max registered `tts-child` slots = 1; 0 after completion |
| Cache artifact | `cache\tts\tts_1795633353870418600.wav`, 196,884 B, RIFF, 22,050 Hz |
| Cancellation | 150 ms deadline aborted synthesis; 0 registered children afterwards; no orphan process |
| Capability catalog | `piper` available; `melotts` unavailable (not a Howl-owned service); `kokoro` unavailable (no assets) |

## 3. Host readiness

- Rebuilt `build/bin/howl-chat-s1.exe` and `build/bin/howl-ui-harness.exe`.
- **Browser dev host (verified visually):** at `http://127.0.0.1:8799/audio-settings.html` with an isolated root, General tab shows **Enable Text-to-Speech (TTS)** enabled and **Status: ready**; TTS tab shows engine **Piper (Local, fast)** and voice **en_US-lessac-medium (built-in)**. `ApplyAudioSettings` returned in 255 ms and the isolated `settings.json` records `ttsEnabled: true`, `ttsEngine: piper`, `ttsVoice: en_US-lessac-medium`.
- **Packaged desktop host (launched, not clicked):** `HOWL_WEBVIEW_DATA_DIR=<isolated root> build/bin/howl-chat-s1.exe` reached `OnStartup: webview reached backend` with no error-log entries and no `piper.exe` child. Its Settings panel was **not** clicked — this environment has no native mouse/screenshot control for a WebView2 window.
- No synthesis was triggered in either host: the isolated cache directory stayed empty, so nothing was spoken and no autoplay occurred.

## 4. Your listening test

1. Launch from the project root: `start-howl-chat.bat` (it `cd`s to the project root, so `bin\` and `models\` resolve) — or use the already-running browser host at `http://127.0.0.1:8799`.
2. Settings → Audio (in the browser host: the gear menu → **TTS Settings**, or `audio-settings.html`) → confirm **Status: ready**, engine **Piper (Local, fast)**, voice **en_US-lessac-medium**.
3. Type or paste this into the chat input and click the speaker icon (**Speak the draft — does not send**). It speaks only the draft: it never sends, saves or creates a message, and needs no chat model. Every AI reply also carries its own **Speak this message** button:
   > Howl, this is the local speech engine. Read this sentence at a normal pace so I can judge clarity, volume and timing. The next sentence checks that the end of a longer passage is still audible and not clipped.
4. Expected artifact: `tts_<unix-nanoseconds>.wav` in `cache\tts` (isolated run: `<isolation root>\data\cache\tts`), RIFF/WAVE, 22,050 Hz mono 16-bit, tens to hundreds of KB.
5. Playback happens only when you click the speaker icon; volume follows the app's TTS volume setting.

## 5. Operator result (fill in)

| Question | Answer |
|---|---|
| Was speech audible? | Yes — operator confirmed both the per-message read-aloud and the composer draft playback. |
| Clarity / pronunciation acceptable? | Yes — operator: "works fine"; no pronunciation defect reported. |
| Volume acceptable? | Not separately reported; no volume defect reported at the default 0.8. |
| Timing/speed acceptable? | Yes — operator: "it isn't that much lag either." |
| Artifact path and size | Composer run in the browser host: `tts_1795635727385868000.wav`, 71,840 B, RIFF/WAVE 22,050 Hz mono 16-bit (~1.6 s) under the isolated root's `data/cache/tts`. |
| `piper.exe` observed during synthesis, absent after? | No `piper.exe` remained after synthesis or after the run; each synthesis is a job-owned child reaped on completion. |
| Errors shown anywhere in the UI or logs | None for TTS. Chat send failed separately with `dial tcp [::1]:8080: connectex … refused` because no llama-server was listening; that is the chat backend, not TTS. |
| Host used | The operator's running app (they reported the result from it directly; the exact host was not separately confirmed); browser dev host (`howl-ui-harness`, isolated root) for the driven composer and message-button verification. |
| Accepted by | The user (human listening operator) |
| Date | 2026-11-25 |

## 6. Acceptance rule

Speech quality is **not** accepted until §5 is filled by the operator. Automated lifecycle and readiness evidence does not substitute for listening, and this record must not be marked accepted on the operator's behalf.

**Resolution (2026-11-25).** The operator supplied the listening result directly, so §5 records their assessment rather than an agent inference. Gate 8's voice-quality question is therefore closed; the gate itself remains open on integration and lifecycle coverage (cross-workload resource/lifecycle testing, packaged-host audio integration evidence, Kokoro/MeloTTS coverage, ASR admission policy/test discrepancy), not on further proof that Piper can speak.
