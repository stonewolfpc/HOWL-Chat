# AUDIT — Audio model loading and ASR settings

**Date:** 2026-11-05
**Repository:** `C:\Users\StoneWolfPC\HOWL_Chat`, branch `master`, HEAD `4f5defe`
**Runtime tested:** `runtime/llama/llama-server.exe` version `0.8.0-dev (build 12015, commit a34de8618)`, Windows x86_64
**Mode:** investigation only. No application source or settings were changed by this audit. **Implementation is not authorized.**
**Review:** owner review incorporated 2026-11-05; see §9 and §11.

This document exists so the loader and ASR settings failures can be reviewed independently, including by an external reviewer. It separates what was reproduced from what is still a hypothesis.

> **Review disposition:** The findings were accepted on 2026-11-05. This revision adds the required metadata-based architecture classification, the shared projector-pairing contract, the Saved/Applied/Ready distinction, the failure-boundary tests, and the A/B/C execution boundaries. No implementation has begun.

---

## 1. Bottom line

1. **The `vocoder-LFM2.5-Audio-1.5B-Q4_0.gguf` crash is expected.** The file is a vocoder component, not a chat LLM. Its GGUF `general.architecture` value is literally the refusal string:
   `this model cannot be used as LLM, use it via --model-vocoder in TTS examples`.
   llama-server exits with status 1 before loading weights. The app exposes it as a selectable LLM and reports only `LLM-SRV-CRASH … exit status 1`, so the actual reason never reaches the user or the error log.
2. **The Qwen3-ASR model itself is healthy in the current runtime.** The local pair
   `Qwen3-ASR-1.7B-Q8_0.gguf` + `mmproj-Qwen3-ASR-1.7B-Q8_0.gguf`
   loaded in about 4 seconds, returned `{"status":"ok"}` from `/health`, and transcribed `build/tmp/asr-probe.wav` correctly. The ASR failure is therefore in **settings persistence and apply lifecycle**, not in the model files or llama-server compatibility.
3. **All six audio-capable base models in `G:\AI GGUF\voice` have a matching projector.** Two projectors are hybrid audio+vision; four are audio-only. The inventory is in §3.
4. **The app does not implement direct audio understanding in chat.** The capability contract explicitly blocks `llama.cpp/audio-understanding` because no application code constructs an `input_audio` request. These models can be inventoried and loaded by llama-server, but the chat pipeline cannot consume audio with them today.
5. **The chat projector finder is not model-safe.** It falls back to the first `mmproj-*.gguf` in the directory. In the voice directory that is the 0.6B ASR projector, so loading the 1.7B ASR model, Ultravox, Voxtral, or Gemma through the chat loader would attach the wrong projector. `Qwen2.5-Omni`'s projector is not found at all because its filename does not start with `mmproj-`.
6. **The ASR settings are not reaching disk.** The save path schedules a 150 ms debounced `SetSettings` call inside the iframe, then immediately removes the iframe. The timer dies with the iframe, so the backend write never happens. `config/settings.json` still contains the defaults and has not been modified since 2026-11-01. The backend-loaded defaults then override the localStorage copy on the next open.

---

## 2. Confirmed vs. unproven

| ID | Status | Statement |
|---|---|---|
| C1 | Confirmed | The LFM vocoder is not loadable as an LLM; llama-server rejects its architecture string. |
| C2 | Confirmed | The local Qwen3-ASR 1.7B model + projector load, become healthy, and transcribe a real WAV. |
| C3 | Confirmed | The ASR save path removes the iframe before the debounced backend save can fire. |
| C4 | Confirmed | `config/settings.json` still has `asrEngine: "auto"`, `asrEnabled: false`, and an old mtime. |
| C5 | Confirmed | The chat projector fallback is directory-order based and can select a different model's projector. |
| C6 | Confirmed | The capability contract blocks direct chat audio understanding; no `input_audio` request is built. |
| C7 | Confirmed | The model picker only filters `mmproj`; the vocoder remains selectable. |
| H1 | **Hypothesis** | The app's `stderrBuf` is empty because `io.MultiWriter(os.Stderr, stderrW)` stops at `os.Stderr` when the packaged GUI has no usable stderr handle. The manual console run has a valid stderr and captures 1007 bytes. This has not been proven against the packaged app. |
| H2 | **Unknown** | Which ASR engine the user selected before the settings were lost. The backend file does not retain it. |
| H3 | **Unknown** | Whether the user's last save applied the engine in memory before the iframe closed. The call is fire-and-forget and is not logged. |

---

## 3. Audio model inventory

Scanned with GGUF metadata, not filenames. All files are under `G:\AI GGUF\voice` unless noted.

| Base model | `general.architecture` | Size | Paired projector | Projector kind | Pair status | App use today |
|---|---|---:|---|---|---|---|
| `Gemma-4-E2B-it-Q8_0.gguf` | `gemma4` | 3864.1 MB | `mmproj-gemma-4-E2B-it-Q8_0.gguf` | hybrid (audio+vision) | Complete | Chat loader only; chat audio is blocked. |
| `Qwen2.5-Omni-3B-Q4_K_M.gguf` | `qwen2vl` | 819.6 MB | `Qwen2.5-Omni-3B-mmproj-Q8_0.gguf` | hybrid (audio+vision) | Complete on disk | Chat loader only; its projector is **not found** by `findMMProj`. |
| `Qwen3-ASR-0.6B-Q8_0.gguf` | `qwen3vl` | 767.5 MB | `mmproj-Qwen3-ASR-0.6B-Q8_0.gguf` | audio-only | Complete | No local ASR registry entry; only the remote 0.6B entry exists. |
| `Qwen3-ASR-1.7B-Q8_0.gguf` | `qwen3vl` | 2064.7 MB | `mmproj-Qwen3-ASR-1.7B-Q8_0.gguf` | audio-only | Complete | Local ASR engine `qwen3-asr-1.7b-local`; manually verified. |
| `Ultravox-Llama-3.2-1B-Q4_K_M.gguf` | `llama` | 668.2 MB | `mmproj-Ultravox-v0_5-llama-3_2-1b-f16.gguf` | audio-only | Complete | Chat loader only; chat audio is blocked. |
| `Voxtral-Mini-3B-2507-Q4_K_M.gguf` | `llama` | 2358.4 MB | `mmproj-Voxtral-Mini-3B-2507-Q8_0.gguf` | audio-only | Complete | Chat loader only; chat audio is blocked. |

Other relevant files:

| File | Architecture / meaning | Size | Projector | Conclusion |
|---|---|---:|---|---|
| `G:\AI GGUF\GGUF\LFM2.5-8B-A1B-heretic.Q4_K_M.gguf` | `lfm2moe` | 4916.7 MB | none | Text model; not LFM2.5-Audio. |
| `G:\AI GGUF\GGUF\vocoder-LFM2.5-Audio-1.5B-Q4_0.gguf` | refusal string (see C1) | 103.9 MB | none | Component-only; no base LFM2.5-Audio model or matching projector exists under `G:\AI GGUF`. |

**Answer to the original inventory question:** six base models accept audio, and all six have their matching projector on disk. The LFM vocoder is not a seventh model; it has no counterpart base model in the collection.

### Projector metadata evidence

- Qwen2.5-Omni and Gemma 4 are `hybrid`:
  - `clip.has_vision_encoder = true`
  - `clip.has_audio_encoder = true`
- Qwen3-ASR, Ultravox, and Voxtral are `audio`:
  - `clip.has_audio_encoder = true`
  - `clip.has_vision_encoder` absent/false
- The projector parser (`internal/backend/gguf/projector.go`) classified all six correctly through `clip.has_audio_encoder` and `clip.audio.*` keys.

---

## 4. Chat model loading pipeline

### 4.1 Flow

```
model picker
  frontend/dist/js/chat-model.js
    _isChatModel()                 // only filters ^mmproj
    _loadModelByPath()
        -> App.LoadModelMetadata()
        -> App.LoadModel(modelPath)
              app-model.go:27
              - stat path
              - port 8080 check
              - ValidateLlamaServer()
              - gguf.ExtractModelProfile()
              - applySafeDefaults()
              - findMMProj()       app-llama.go:43
              - profile.HasAudioInput / HasVision from projector metadata
              - ValidateSpawn()
              - spawnLlamaServerWithModel()  app-llama.go:181
                    - calls findMMProj() AGAIN
                    - builds -m, -c, --threads, -b, -ngl, --mmproj
                    - cmd.Stdout = os.Stdout
                    - cmd.Stderr = io.MultiWriter(os.Stderr, stderrBuf)
              - waitMonitor(cmd)   app-llama.go:302
              - monitor loop       app-model.go:228+
                    - every 500 ms checks process alive and HTTP /health
                    - on exit: read stderrBuf, keep last 4 lines, emit model:error,
                      log LLM-SRV-CRASH
```

### 4.2 Where the vocoder fails

Manual command:

```bash
cd /c/Users/StoneWolfPC/HOWL_Chat
./runtime/llama/llama-server.exe \
  -m "G:\AI GGUF\GGUF\vocoder-LFM2.5-Audio-1.5B-Q4_0.gguf" \
  --host 127.0.0.1 --port 8099 -c 512 --threads 4 -b 128 -ngl 0 --fit off
```

Observed stderr (`build/bin/vocoder-stderr.log`, 1007 bytes):

```text
0.00.131.058 I srv    load_model: loading model 'G:\AI GGUF\GGUF\vocoder-LFM2.5-Audio-1.5B-Q4_0.gguf'
0.00.131.661 E llama_model_load: error loading model: unknown model architecture: 'this model cannot be used as LLM, use it via --model-vocoder in TTS examples'
0.00.131.668 E llama_model_load_from_file_impl: failed to load model
0.00.131.671 E cmn  common_init_: failed to load model '...'
0.00.131.675 E srv    load_model: failed to load model, '...'
0.00.131.677 I srv  operator(): operator(): cleaning up before exit...
0.00.132.455 E srv  server_run_i: exiting due to model loading error
```

Exit code: `1`. stdout was 0 bytes; all diagnostics were on stderr.

The app's error log (`config/logs/howl-errors-20261105.jsonl`) contains only:

```json
{"t":"2026-11-05T00:39:04.5811911Z","c":"LLM-SRV-CRASH","s":"critical","op":"llm.load","loc":"app-model.go:146","msg":"llama-server exited unexpectedly: exit status 1"}
{"t":"2026-11-05T00:41:58.5767413Z","c":"LLM-SRV-CRASH","s":"critical","op":"llm.load","loc":"app-model.go:146","msg":"llama-server exited unexpectedly: exit status 1"}
```

The real reason is absent from both records.

### 4.3 Why the real reason is missing

`spawnLlamaServerWithModel` sets:

```go
cmd.Stdout = os.Stdout                                   // app-llama.go:245
cmd.Stderr = io.MultiWriter(os.Stderr, stderrW)          // app-llama.go:247
```

`io.MultiWriter` writes to each writer in order and returns on the first error. In a packaged Wails GUI, `os.Stderr` may be an unusable handle. If it fails, `stderrW` (the in-memory `stderrBuf`) is never written. In the console reproduction `os.Stderr` is valid, so the buffer would work; the app's persisted record shows the buffer was empty.

The same ordering exists in the image engine (`internal/image/engine/sdcpp/engine.go:184`), so this is a shared pattern, not unique to the LLM loader.

Additional diagnostics gaps:

- No standalone llama-server log file is created. `find` found only `.freebuff/llama-b12015.log` and `.freebuff/llama-test.log` from earlier manual experiments, plus `server-test-err.log` from an older build. None is a current app-produced log.
- `build/logs/app-run*.log` contain image-generation output, not `[LoadModel]` lines.
- The monitor keeps only the last four stderr lines in the in-memory message (`app-model.go:261-291`).
- `LLM-SRV-CRASH` is sent through the error system. The file logger records `msg`; the UI event payload carries code/severity/operation/location/time, not the message. That is why the toast shows a code/operation label rather than the architecture error.
- The location string is hardcoded as `app-model.go:146`; the actual log call has moved in the file. The stale location is cosmetic but makes source correlation harder.

### 4.4 Projector pairing defects

`findMMProj` (`app-llama.go:43`) has a long list of exact names, then:

```go
// Fallback: scan directory for any file starting with "mmproj" and ending ".gguf"
...
if strings.HasPrefix(name, "mmproj") && strings.HasSuffix(name, ".gguf") {
    return filepath.Join(dir, entry.Name())
}
```

`os.ReadDir` returns entries sorted by filename, so the first match in `G:\AI GGUF\voice` is:

```text
mmproj-Qwen3-ASR-0.6B-Q8_0.gguf
mmproj-Qwen3-ASR-1.7B-Q8_0.gguf
mmproj-Ultravox-v0_5-llama-3_2-1b-f16.gguf
mmproj-Voxtral-Mini-3B-2507-Q8_0.gguf
mmproj-gemma-4-E2B-it-Q8_0.gguf
```

Consequences if the voice directory is used as a chat model directory:

- `Qwen3-ASR-1.7B` gets the 0.6B projector.
- `Ultravox`, `Voxtral`, and `Gemma` get the 0.6B projector.
- `Qwen2.5-Omni` gets no projector because its file is `Qwen2.5-Omni-3B-mmproj-Q8_0.gguf`, which does not start with `mmproj`, and the model-specific patterns expect `Qwen2.5-Omni-3B-Q4_K_M-mmproj.gguf` or `.mmproj.gguf`, neither of which exists.

The ASR resolver does **not** use `findMMProj`. It resolves the exact `ProjectorFile` from the ASR registry, so the ASR pair is correct. The defect is in the generic chat loader, and it is latent while `G:\AI GGUF\voice` is only registered as an ASR directory.

### 4.5 Picker filtering

`_isChatModel` (`frontend/dist/js/chat-model.js:237`) only rejects projectors:

```js
return !/^mmproj(?:[-_.]|$)/.test(base) && !/^mmproj(?:[-_.]|$)/.test(file);
```

The vocoder passes. It is in `config/model-registry.json` with `type: "llm"` and the refusal string as its `family`, so it appears in the picker and can be selected again.

### 4.6 Direct audio chat is not implemented

`internal/audio/capability/capability.go:438-452` declares:

```go
ID:            "llama.cpp/audio-understanding",
Task:          TaskAudioUnderstanding,
Provider:      "llama.cpp chat completions (input_audio)",
ProviderState: Blocked("no application code constructs a multimodal audio request"),
ModelState:    UnknownDim("audio model compatibility is not negotiated"),
Projector:     UnknownDim("projector kind is read from GGUF metadata; no audio projector is enumerated"),
TaskSupport:   Blocked("no adapter implements direct audio understanding"),
```

Therefore the six audio models can be detected and loaded, but audio attachments cannot be sent to them through chat. Only the separate ASR pipeline can consume audio, and only the Qwen3-ASR 1.7B local entry is registered for it.

---

## 5. ASR settings and loading pipeline

### 5.1 Flow

```
audio-settings.html (iframe)
  user selects ASR engine / enables ASR
  _wireSelect() updates window.settingsState.current.audio.asrEngine
  _wireToggle() updates window.settingsState.current.audio.asrEnabled

Save & Close
  audio-settings.js:340  saveAudioSettings()
    -> settingsState.saveSettings()          // localStorage now; backend setTimeout 150ms
    -> App.ApplyAudioSettings(...)           // fire-and-forget
    -> close()                               // audio-settings.js:57
         -> settingsState.saveSettings()     // schedules another 150ms timer
         -> parent audio-settings-modal.remove()
    -> show "Audio settings saved"
```

Backend path when `ApplyAudioSettings` is actually invoked:

```
app-audio.go:287
  -> applyAudioSettings
    -> audio.Manager.Apply
      -> Pipeline.Initialize
        -> asrCandidatesFor(plan, explicit engine)
          -> capability.Select(engine)
        -> ASRModelPathFor / ASRProjectorPathFor
          -> internal/audio/model_resolver.go
            -> exact LocalFile and ProjectorFile
        -> internal/audio/asr/llama.Initialize
          -> llama-server -m <model> --mmproj <projector>
          -> waits for /health
```

### 5.2 The persistence defect

`settings-state.js:164-190`:

```js
async function saveSettings() {
    // localStorage immediate
    ...
    if (saveDebounceTimer) clearTimeout(saveDebounceTimer);
    saveDebounceTimer = setTimeout(async () => {
        saveDebounceTimer = null;
        if (window.go?.main?.App?.SetSettings) {
            await window.go.main.App.SetSettings(currentSettings);
        }
    }, 150);
}
```

`audio-settings.js:340-358`:

```js
function saveAudioSettings() {
    window.settingsState.saveSettings();
    ...
    app.ApplyAudioSettings(...).catch(...);   // not awaited
    close();
    _notify('Audio settings saved', 'success');
}
```

`close()` (`audio-settings.js:57-67`) removes the parent modal:

```js
function close() {
    ...
    window.settingsState.saveSettings();
    const parentModal = window.parent.document.getElementById('audio-settings-modal');
    if (parentModal) parentModal.remove();
}
```

`menu-functions.js:28-52` creates the iframe and deletes the whole modal, not just hides it:

```js
modal.id = 'audio-settings-modal';
...
const onEsc = ... modal.remove() ...
```

Removing the iframe destroys its document and its pending timers. The two `setTimeout(..., 150)` callbacks therefore never call `SetSettings`.

The immediate localStorage write does happen, but it cannot win on the next open:

1. `loadSettings` calls `GetSettings()` first.
2. `config/settings.json` exists and contains a non-empty `audio` section, so `loaded = true`.
3. Because `loaded` is true, the localStorage fallback is skipped.
4. The backend's old `audio.asrEngine: "auto"` and `audio.asrEnabled: false` are merged over the defaults.

`config/settings.json` confirms the result:

```json
"audio": {
  "asrEnabled": false,
  "asrEngine": "auto",
  ...
}
```

Its mtime is `2026-11-01 10:16:25 -0600`; the current investigation is on 2026-11-05. There is no sign that an ASR save reached the file.

### 5.3 The apply defect

`saveAudioSettings` calls `app.ApplyAudioSettings` but does not `await` it. `close()` runs immediately. If the model load fails, the promise rejection goes to `_notify(...)`, but the modal is already gone and the success notification has already been shown. There is no persistent ASR error log; `applyAudioSettings` returns the error to JS and does not call `a.logError`.

`audio-service.js:62-78` re-applies persisted settings only on `howl:settings-loaded` and only when `asrEnabled` or `ttsEnabled` is true:

```js
async function restorePersistedAudio() {
    if (_restoreApplied) return;
    const s = window.settingsState?.current?.audio;
    ...
    if (!s || (!s.asrEnabled && !s.ttsEnabled)) return;
    ...
    await app.ApplyAudioSettings(!!s.asrEnabled, !!s.ttsEnabled, s.asrEngine || 'auto', s.ttsEngine || 'auto');
}
```

If the backend file says `asrEnabled: false`, startup performs no ASR initialization at all. That is correct behavior for a disabled feature, but it means a lost `asrEnabled: true` also prevents the selected engine from ever being applied.

### 5.4 Cross-frame state is stale

The settings UI runs in an iframe with its own `settings-state.js` instance. The parent `audio-service.js` has its own `_asrEnabled`/`_ttsEnabled` variables, initialized at `DOMContentLoaded` and updated only by `syncFromSettings()` on the parent's `howl:settings-loaded` event.

When the iframe saves:

- The iframe's `settingsState.current.audio` changes.
- The parent's `settingsState.current.audio` does not change.
- The parent's `AudioService._asrEnabled` does not change.
- No `howl:settings-loaded` event is dispatched by the iframe.

`startRecording()` checks the parent variable (`audio-service.js:101`: `if (!_asrEnabled) { showNotification('Speech recognition is disabled'...) }`). So even if the backend apply succeeds, the microphone can still refuse to record until the app is reloaded.

### 5.5 The ASR model/projector itself is not the problem

Manual load command (same model/projector pair the ASR registry names):

```bash
cd /c/Users/StoneWolfPC/HOWL_Chat
./runtime/llama/llama-server.exe \
  -m "G:\AI GGUF\voice\Qwen3-ASR-1.7B-Q8_0.gguf" \
  --mmproj "G:\AI GGUF\voice\mmproj-Qwen3-ASR-1.7B-Q8_0.gguf" \
  --host 127.0.0.1 --port 8098 -c 2048 -ngl 99 --no-warmup -t 6
```

Observed:

```text
0.00.136.552 I srv load_model: loading model 'G:\AI GGUF\voice\Qwen3-ASR-1.7B-Q8_0.gguf'
0.02.112.492 I srv load_model: loaded multimodal model, 'G:\AI GGUF\voice\mmproj-Qwen3-ASR-1.7B-Q8_0.gguf'
0.02.228.443 I srv server_run_i: model loaded
0.02.228.452 I srv server_on_li: listening on http://127.0.0.1:8098
```

Health response:

```json
{"status":"ok"}
```

Transcription of `build/tmp/asr-probe.wav`:

```json
{"type":"transcript.text.done","text":"language English<asr_text>The quick brown fox jumps over the lazy dog. This is an audio transcription gate test, November fourth.","usage":{"type":"tokens","input_tokens":119,"output_tokens":25,"total_tokens":144,"input_tokens_details":{"cached_tokens":0}}}
```

This proves:

- The model and projector files are complete.
- The bundled llama-server supports the model's `qwen3vl` architecture and audio input.
- `-ngl 99` works on this machine.
- The ASR endpoint works with the correct multipart request.
- The remaining ASR failure is above this layer: settings persistence, apply sequencing, and parent UI state.

### 5.6 Resolver evidence

`internal/audio/paths_api.go`:

```go
func ASRModelPathFor(model types.ASRModelType) string {
    if installed := resolveInstalledASRModel(model); installed != "" {
        return installed
    }
    return defaultASRModelPath(model)
}

func ASRProjectorPathFor(model types.ASRModelType) string {
    return resolveInstalledASRProjector(model)
}
```

`internal/audio/model_resolver.go` resolves `Qwen3-ASR-1.7B-Q8_0.gguf` from `config/asr-model-dirs.json` (`G:\AI GGUF\voice`) and then `ProjectorFile` `mmproj-Qwen3-ASR-1.7B-Q8_0.gguf` from the same directory. The ASR registry marks `ProjectorKind: audio`, so the projector is mandatory. This path is correct and is not the generic `findMMProj` path.

The local ASR registry has only one local entry:

- `qwen3-asr-1.7b-local` → `Qwen3-ASR-1.7B-Q8_0.gguf` + `mmproj-Qwen3-ASR-1.7B-Q8_0.gguf`
- The 0.6B local pair exists on disk but has no local registry entry. The UI's 0.6B option is the DashScope remote entry.

---

## 6. Findings

### F-01 — Component GGUF is exposed as a chat model (critical, confirmed)

**Where:** `frontend/dist/js/chat-model.js:237-243`, `config/model-registry.json` vocoder entry, `app-model.go:27`.

**Evidence:** The vocoder's `general.architecture` is a refusal string. It is registered as `type: "llm"` and passes `_isChatModel` because it is not named `mmproj`.

**Impact:** The user can select it repeatedly; every attempt exits with status 1 and a generic error.

**Fix direction:** Classify the asset from GGUF metadata into supported chat, supported multimodal chat, component-only, or unknown. Unknown fails closed; component-only is a wrong-asset-type error, not merely an unsupported model. The registry must not override metadata. The picker must use the same classifier. Do not attempt to load a vocoder as an LLM.

### F-02 — Loader has no architecture gate (critical, confirmed)

**Where:** `app-model.go:95-160`, `app-llama.go:181-250`.

**Evidence:** `ExtractModelProfile` accepts any parseable GGUF; `spawnLlamaServerWithModel` passes it directly to `llama-server`. The server, not the app, rejects the architecture.

**Impact:** Any non-LLM GGUF reaches llama-server and produces an opaque crash.

**Fix direction:** Implement an architecture capability classifier rather than relying on a bare hand-maintained allowlist. Distinguish supported chat, supported multimodal chat, component-only, and unknown/unsupported; fail closed with a precise user-facing reason. Registration cannot make an incompatible GGUF loadable.

### F-03 — Crash diagnostics are incomplete and not persisted (critical, confirmed gap; H1 for cause)

**Where:** `app-llama.go:245-249`, `app-model.go:261-302`, `app-errors.go:176`, `internal/errors/wails.go:120`.

**Evidence:** Manual stderr is 1007 bytes; the app error record has only `exit status 1`; no standalone llama-server log exists. `io.MultiWriter(os.Stderr, stderrBuf)` is ordered with the potentially invalid GUI stderr first.

**Impact:** The user and maintainers cannot see why a model failed; every loader failure becomes a generic code.

**Fix direction:** Tee stderr to the buffer first (or use a custom writer that never short-circuits); capture stdout as well; write a per-load `config/logs/llama-server-<timestamp>.log` with args, model path, projector path, exit code, and full output; include a bounded tail in the error record and UI.

### F-04 — Projector fallback can attach the wrong model's mmproj (high, confirmed by code; latent in current UI)

**Where:** `app-llama.go:43-112`, especially the wildcard at 99-109.

**Evidence:** First `mmproj-*` in `G:\AI GGUF\voice` is the 0.6B ASR projector. The loader does not compare model stem, projection dimensions, or architecture.

**Impact:** If the voice directory is added to chat scanning, Qwen3-ASR 1.7B, Ultravox, Voxtral, and Gemma receive the 0.6B projector; llama-server may fail or produce invalid multimodal behavior.

**Fix direction:** Use a shared projector-pairing contract, not directory order and not a simple reuse of the ASR resolver. The contract takes explicit model identity and expected companion metadata, resolves the exact intended companion, validates compatibility where metadata allows, and rejects missing, ambiguous, or mismatched pairs. It must never silently substitute another model's projector.

### F-05 — Qwen2.5-Omni projector is undiscoverable (high, confirmed)

**Where:** `app-llama.go:43-97`, `G:\AI GGUF\voice\Qwen2.5-Omni-3B-mmproj-Q8_0.gguf`.

**Evidence:** The file does not start with `mmproj`, and the model-specific patterns expect a different suffix. `findMMProj` returns `""`.

**Impact:** Even if Qwen2.5-Omni is selected as a chat model, its audio/vision projector is not attached.

**Fix direction:** Treat the Omni naming variation as a test case for the shared pairing contract. Support an exact `<model-stem>-mmproj*.gguf` form only when identity is unambiguous; do not add another arbitrary fallback that can select an unrelated file.

### F-06 — Chat audio understanding is unimplemented (high, confirmed)

**Where:** `internal/audio/capability/capability.go:438-452`.

**Evidence:** `ProviderState` and `TaskSupport` are both `Blocked`; no `input_audio` request exists.

**Impact:** The six audio-capable GGUFs cannot be used for audio chat; only the ASR pipeline can consume audio, and only Qwen3-ASR 1.7B is registered locally.

**Fix direction:** Keep direct audio understanding out of scope for the loader and ASR repair units. Until an `input_audio` adapter exists, keep `llama.cpp/audio-understanding` blocked and do not advertise these models as chat-audio capable.

### F-07 — ASR save is destroyed with the iframe (critical, confirmed)

**Where:** `frontend/dist/js/settings-state.js:164-190`, `frontend/dist/js/audio-settings.js:57-67, 340-358`, `frontend/dist/js/menu-functions.js:28-52`.

**Evidence:** Both `saveSettings()` calls schedule a 150 ms timer. `close()` immediately removes `#audio-settings-modal` from the parent document. The iframe document and its timers are destroyed. `config/settings.json` still has the default ASR values and an old mtime.

**Impact:** Every ASR engine/toggle change is lost on close/reopen and app restart. The UI appears to save, then resets.

**Fix direction:** Add a `flushSettings()`/`saveSettingsNow()` that awaits `SetSettings`; make the save handler async; await persistence before removing the iframe; or keep the iframe alive until the write resolves. Complete hydration before enabling Save. Do not rely on a debounce timer inside a modal that is about to be destroyed.

### F-08 — Apply is fire-and-forget and errors are not persisted (high, confirmed)

**Where:** `frontend/dist/js/audio-settings.js:340-358`, `app-audio.go:287-303`.

**Evidence:** `ApplyAudioSettings` is called without `await`; `close()` and the success toast happen immediately. There is no `logError` in `applyAudioSettings`; failures only reach the JS `.catch`.

**Impact:** A failed ASR load looks like a successful save; the user has no actionable reason.

**Fix direction:** Model three separate outcomes: **Saved** (backend persisted the choice), **Applied** (audio manager accepted the configuration), and **Ready** (the ASR engine initialized and can transcribe). Persist the selected preference even if activation fails, but never present "saved" as a synonym for "ready." Await `ApplyAudioSettings`, show visible progress, keep the failure recoverable, and display the backend error. Log ASR load failures with model, projector, server path, and stderr tail.

### F-09 — Parent audio state is not synchronized after iframe save (high, confirmed)

**Where:** `frontend/dist/js/audio-service.js:46, 53-58, 101`, `frontend/dist/js/audio-settings.js:340-358`.

**Evidence:** The iframe has its own settings state; it does not update `window.parent.settingsState` or dispatch a parent event. Parent `_asrEnabled` remains its startup value.

**Impact:** Even when the backend applies ASR, `startRecording()` can still show "Speech recognition is disabled" until reload.

**Fix direction:** On save, update the parent settings state or dispatch a parent event that calls `AudioService.syncFromSettings()`; or move audio settings into the parent document instead of an iframe.

### F-10 — Settings load/hydrate is racy (medium, confirmed)

**Where:** `frontend/dist/js/audio-settings.js:45-54, 77-80`.

**Evidence:** `show()` calls `loadSettings()` without awaiting, then immediately calls `_hydrateAllControls()`. The `DOMContentLoaded` handler does the same. The controls can hydrate from defaults before the async backend response merges the saved values.

**Impact:** Visual flicker or a stale value on open; combined with F-07 it makes the reset look intentional.

**Fix direction:** `await loadSettings()` before hydrating; use an initialization promise; disable Save until hydration completes.

### F-11 — Backend settings override localStorage after a lost write (medium, confirmed)

**Where:** `frontend/dist/js/settings-state.js:128-160`.

**Evidence:** `loadSettings` sets `loaded = true` whenever a non-empty backend object is returned and skips localStorage. The backend file always has a non-empty `audio` section, even when it is stale.

**Impact:** The localStorage copy of the new ASR choice is ignored on the next open.

**Fix direction:** Fix F-07 first. Do not independently redesign precedence in this repair unit; keep the current backend-first behavior until the write path is proven. Revisit precedence only if the evidence after the write fix requires it.

### F-12 — Picker only filters projectors, not component files (medium, confirmed)

**Where:** `frontend/dist/js/chat-model.js:237-243`.

**Evidence:** The vocoder passes the `_isChatModel` check; it is registered as `type: "llm"`.

**Impact:** The same crash is easy to reproduce.

**Fix direction:** Derive picker admission from the same metadata-based architecture classifier used by the loader. Component-only assets are hidden as the wrong asset type; do not maintain a second filename-only filter.

### F-13 — Loader location strings are stale (low, confirmed)

**Where:** `app-model.go:302` logs `"app-model.go:146"`.

**Evidence:** The error records show `loc: "app-model.go:146"` while the current call is at line 302. The line number is hardcoded.

**Impact:** Source correlation is misleading.

**Fix direction:** Use a stable operation/location identifier or `runtime.Caller` instead of a hardcoded line.

---

## 7. What is not the problem

These were checked and should not be the first suspect:

- **The ASR model files are not corrupt.** Both files are present and load.
- **The bundled llama-server is not too old for Qwen3-ASR.** Build 12015 loaded the `qwen3vl` model and audio projector and produced a transcript.
- **The ASR projector pairing is not the bug.** The ASR resolver uses the registry's exact `ProjectorFile`; the generic `findMMProj` path is not used by ASR.
- **`config/asr-model-dirs.json` is correct.** It contains `G:\AI GGUF\voice`.
- **The settings store can write.** `internal/settings/store.go` uses a direct `os.WriteFile`; the file exists and is readable. The issue is that the frontend never calls `SetSettings` before destroying the iframe.
- **`GetSettings` is not returning an empty map.** `config/settings.json` is non-empty; that is why the stale backend wins over localStorage.
- **Port 8080 was not the reported failure.** The error records are model-load exits, not port conflicts; the manual tests used free alternate ports.
- **The ASR model does not need the chat `findMMProj` fallback.** It has its own resolver and exact projector path.
- **The old `server-test-err.log` is not current evidence.** It is from an older build (`b9049`) and a different router experiment.

---

## 8. Reproduction cookbook

### 8.1 Reproduce the vocoder rejection

```bash
cd /c/Users/StoneWolfPC/HOWL_Chat
rm -f build/bin/vocoder-stdout.log build/bin/vocoder-stderr.log
./runtime/llama/llama-server.exe \
  -m "G:\AI GGUF\GGUF\vocoder-LFM2.5-Audio-1.5B-Q4_0.gguf" \
  --host 127.0.0.1 --port 8099 -c 512 --threads 4 -b 128 -ngl 0 --fit off \
  > build/bin/vocoder-stdout.log 2> build/bin/vocoder-stderr.log
echo "exit=$?"   # 1
cat build/bin/vocoder-stderr.log
```

Expected: `unknown model architecture: 'this model cannot be used as LLM, use it via --model-vocoder in TTS examples'`.

### 8.2 Reproduce the ASR load and transcription

Start the server in one terminal:

```bash
cd /c/Users/StoneWolfPC/HOWL_Chat
./runtime/llama/llama-server.exe \
  -m "G:\AI GGUF\voice\Qwen3-ASR-1.7B-Q8_0.gguf" \
  --mmproj "G:\AI GGUF\voice\mmproj-Qwen3-ASR-1.7B-Q8_0.gguf" \
  --host 127.0.0.1 --port 8098 -c 2048 -ngl 99 --no-warmup -t 6
```

Wait for `listening on http://127.0.0.1:8098`, then in another terminal:

```bash
curl -sS \
  -F "file=@build/tmp/asr-probe.wav" \
  -F "model=qwen3-asr" \
  -F "response_format=json" \
  http://127.0.0.1:8098/v1/audio/transcriptions
```

Expected transcript starts with:

```text
language English<asr_text>The quick brown fox jumps over the lazy dog...
```

### 8.3 Reproduce the settings loss

1. Open the audio settings iframe, enable ASR, select `Qwen3 ASR 1.7B (local llama.cpp)`, and click **Save & Close**.
2. Immediately inspect `config/settings.json`.
3. Expected current behavior: it still contains `"asrEnabled": false`, `"asrEngine": "auto"` and its old mtime.
4. Reopen the iframe: the controls show the defaults.

This is a deterministic consequence of the 150 ms timer being destroyed by `parentModal.remove()`.

### 8.4 Check whether the packaged app has a usable stderr

This is the unresolved H1. A controlled test is needed:

- Launch the packaged app with an explicit stdout/stderr file handle, or add a temporary diagnostic that writes the result of a small `os.Stderr.Write` before spawning llama-server.
- If stderr writes fail in the GUI, the `io.MultiWriter` order is the reason the buffer is empty.
- If stderr works, the buffer is empty for another reason and the loader path needs more instrumentation.

Do not treat H1 as proven until this test is run.

---

## 9. Execution plan after review (not authorized)

**Implementation is not authorized.** The review accepted the findings but required the following boundaries. The work is split into three independently reviewable units.

### Unit A — Loader integrity (first)

1. **Classify assets from GGUF metadata, not filenames and not the registry.**
   - Supported chat architecture.
   - Supported chat architecture with a multimodal projector.
   - Component-only asset (vocoder, encoder/decoder-only).
   - Unknown or unsupported architecture.
   - Unknown fails closed with an informative message.
   - Component-only is a wrong-asset-type error, not merely an unsupported model.
   - The registry cannot override what the GGUF metadata says.
   - Invariant: a model's registration cannot make an incompatible GGUF loadable as a chat model.
2. **Persist useful loader diagnostics.**
   - Capture stdout and stderr to a per-load file.
   - Buffer first, or use a writer that does not short-circuit; H1 still requires proof.
   - Include architecture classification, model path, projector path, args, exit code, and a bounded output tail in the error record and UI.
3. **Introduce the shared projector-pairing contract.**
   - Do not simply reuse the ASR resolver; chat discovery has different inputs and may support more families.
   - Inputs: explicit model identity and expected companion metadata.
   - Resolve the exact intended companion; validate compatibility where metadata allows; reject missing, ambiguous, or mismatched pairs; never silently substitute.
   - Qwen2.5-Omni naming is a test case, not a reason for another arbitrary fallback.
4. **Use the same classifier for picker admission.** Component-only assets are hidden as the wrong asset type.

### Unit B — ASR settings lifecycle

1. **Persist reliably.** Add `saveSettingsNow()`/flush that returns the `SetSettings` promise; await it before removing the iframe. Complete hydration before enabling Save so F-10 cannot race with user interaction.
2. **Separate the outcomes.** Model **Saved**, **Applied**, and **Ready** independently. Persist the preference even if activation fails, report the activation failure separately, and never display "saved" as "ready."
3. **No blind 120-second modal block.** Show visible progress and make the operation cancellable or otherwise recoverable. The user should know whether saving, loading, or initialization is underway.
4. **Synchronize the parent.** Update the parent `settingsState` or dispatch a parent event; call `AudioService.syncFromSettings()` so recording immediately after a successful activation uses the new `asrEnabled` without restarting.
5. **Log activation failures.** Record model, projector, server path, and stderr tail to `config/logs`.

### Unit C — Integrated verification

- Real ASR model load, `/health`, and a real WAV transcription.
- Save, close, and reopen; then save, restart the app, and reopen.
- Immediate microphone use after successful activation.
- Ordinary text-chat regression.

Required failure and lifecycle tests:

**Loader**
- Unknown architecture.
- Known component-only architecture.
- Missing projector.
- Wrong projector.
- Multiple plausible projectors.
- Child process exits immediately.
- Child process produces stderr when the GUI stderr destination is unavailable.
- Child process produces large diagnostic output without exhausting memory.

**ASR settings**
- Save, close, reopen.
- Save, restart app, reopen.
- Save succeeds but ASR initialization fails.
- Save fails.
- User changes settings while initialization is underway.
- Parent `AudioService` state updates without restarting.
- Recording immediately after successful activation.
- Repeated open/close cycles leave no stale timers or duplicate event handlers.

### Out of scope for these units

- Direct audio understanding / `input_audio` chat.
- Registering the local Qwen3-ASR 0.6B pair.
- Expanding the audio model registry.
- Redesigning settings precedence beyond ensuring the write occurs.

---

## 10. Evidence index

| Artifact | Purpose |
|---|---|
| `build/bin/vocoder-stderr.log` | Exact 1007-byte stderr from the vocoder load; architecture refusal. |
| `build/bin/vocoder-stdout.log` | 0 bytes; confirms the diagnostics were on stderr. |
| `build/bin/vocoder-load.log` | Earlier combined-stream reproduction; same failure. |
| `build/bin/asr-load.log` | Qwen3-ASR 1.7B + projector load; model loaded and listening. |
| `build/bin/asr-transcribe.log` | ASR server timing and successful request. |
| `build/bin/asr-transcribe.json` | Successful transcript response. |
| `build/tmp/asr-probe.wav` | Real WAV used for the ASR gate. |
| `config/logs/howl-errors-20261105.jsonl` | Two generic `LLM-SRV-CRASH` records; no stderr. |
| `config/settings.json` | Stale ASR defaults; mtime 2026-11-01. |
| `config/asr-model-dirs.json` | Points at `G:\AI GGUF\voice`. |
| `config/model-registry.json` | Vocoder registered as `llm` with the refusal family. |
| `app-model.go` | `LoadModel` and monitor; no architecture gate. |
| `app-llama.go` | `findMMProj`, spawn args, stdout/stderr tee, `waitMonitor`. |
| `frontend/dist/js/chat-model.js` | Picker filtering; only `mmproj` is excluded. |
| `frontend/dist/js/settings-state.js` | Debounced backend save; localStorage fallback. |
| `frontend/dist/js/audio-settings.js` | Save handler and iframe close. |
| `frontend/dist/js/menu-functions.js` | Iframe modal creation/removal. |
| `frontend/dist/js/audio-service.js` | Parent audio state and restore path. |
| `internal/audio/model_resolver.go` | Exact ASR model/projector resolution. |
| `internal/audio/asr/llama/llama.go` | ASR llama-server startup and `/health` wait. |
| `internal/audio/capability/capability.go` | Direct chat audio understanding is blocked. |

---

## 11. Review decisions and remaining unknowns

**Decisions recorded 2026-11-05:**

| Item | Decision |
|---|---|
| GUI stderr | H1 remains open. Prove it with the controlled test before attributing the empty buffer to `MultiWriter` ordering. |
| Audio models in chat | Only expose them if they are genuinely usable for the advertised chat task. ASR-only models belong in Speech Recognition, not the chat model picker. |
| Local 0.6B ASR | Do not register yet. Inventory is not authorization to expand the supported model list. |
| Qwen2.5-Omni audio | Not verified. Metadata establishes capability indicators, not working application-level audio input. |
| Persist selection if apply fails | Yes, generally. Report persistence and activation separately; preserve the preference. |
| 120-second modal block | No blind wait. Use visible progress and a cancellable or recoverable operation. |
| Settings precedence | Do not independently redesign it in this repair. Fix the write path first; revisit only if evidence requires it. |
| Implementation | Not authorized by the review. |

**Remaining unknowns:** H1 (packaged-app stderr); Qwen2.5-Omni audio compatibility; physical microphone behavior; packaged WebView behavior.

---

## 12. Testing charter (preserved)

These are acceptance constraints for when implementation is authorized. They are not implementation instructions and do not change the current authorization boundary.

- **Order:** Loader integrity (Unit A), then ASR settings lifecycle (Unit B), then integrated verification (Unit C).
- **Manual ASR transcription is not packaged-app proof.** A successful direct llama-server transcription only proves the model, projector, and runtime path. It does not prove the packaged Wails application wires the same path correctly.
- **Microphone testing must exercise the actual chat interface.** Verify the packaged app's microphone control, recording lifecycle, transcription insertion, visible result, and error states — not merely that the backend receives an audio payload.
- **TTS testing must include character-specific voices in real multi-character conversations.** Verify the correct voice is selected per character, survives turn changes, and behaves correctly when multiple characters speak in sequence.
- **Loader/ASR repair does not close the larger audio-and-TTS chat integration phase.** Unit C completion is a gate for the repair, not completion of the full audio product.
- **The final test must be user-observable.** The end state should be something the user can hear and use in the chat interface, not only a green terminal check.
- **Implementation remains unauthorized.** These notes are preserved for the eventual authorized test phase.
