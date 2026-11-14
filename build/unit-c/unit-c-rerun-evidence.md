# Unit C rerun (post-D3) — evidence record

Date: 2026-11-05. Status: **EVIDENCE ONLY — NOT ACCEPTED, NOT CLOSED.**
Owner review required. No application code, menus, frontend, settings
architecture, D1 records or D3 records were modified. No commit, push, release
or deployment.

This rerun was authorized as verification only. It stopped at the first
unresolved blockers rather than repairing them: no usable chat LLM could be
loaded, and the physical-microphone capture produced an empty transcript.

## 0. Isolation (established before any write-capable test)

The accepted D3 mistake was not repeated. Nothing was launched against the
user's real AppData chat data.

- Chat root is `os.UserHomeDir()/AppData/Roaming/HOWLChat` (verified in
  `internal/backend/chat/persistence.go`), so `USERPROFILE` — not `APPDATA`
  alone — is the isolation lever.
- `paths.ProjectRoot()` searches upward for `wails.json`/`go.mod`, so a copied
  executable with its own markers isolates `config/`, `logs/` and runtime
  resolution.

Isolation target: `build/tmp/unit-c-rerun/`

```
howl-chat-unit-c.exe   copy of build/bin/howl-chat-d3-tags.exe (29,233,152 B)
config/                copy of project config (own settings.json, logs/)
wails.json, go.mod     markers so ProjectRoot resolves here
runtime -> <project>/runtime   directory junction (read-only use)
sd-server.exe          copy (1,582,080 B)
profile/               USERPROFILE / APPDATA / LOCALAPPDATA for the child
```

Launch environment (per process, never system-wide):

```
USERPROFILE=C:\Users\StoneWolfPC\HOWL_Chat\build\tmp\unit-c-rerun\profile
APPDATA=%USERPROFILE%\AppData\Roaming
LOCALAPPDATA=%USERPROFILE%\AppData\Local
```

Before/after proof (`real-before.txt`, `real-after.txt`):
`find %APPDATA%\HOWLChat -type f | sha256sum` is **IDENTICAL** before and after
both launches. `config/settings.json` stayed
`e786f53f2fd7f91b794936dad1930e2b09594de8a81167179a857ecb011757bf`
(mtime 2026-11-04 21:39:55). The isolated profile contains only its own
`HOWLChat\config.json`; no chat directory was created because no message was
sent. **Isolation verdict: PASS.**

## 1. Build and model identity

Binary under test: `build/bin/howl-chat-d3-tags.exe` (the D3-verified build,
tags `desktop,production`), copied byte-for-byte into the isolated root and run
as `howl-chat-unit-c.exe` so the existing UIA driver could address it.

| Artifact | SHA-256 | Size |
|---|---|---|
| `Qwen3-ASR-1.7B-Q8_0.gguf` | `58e22d0532d4eacaf034cfac17a6fed159f37c41390c710186783be439d1fc57` | 2,165,034,944 |
| `mmproj-Qwen3-ASR-1.7B-Q8_0.gguf` | `46c1d533af3f354ceb37ce855dbceff7da7fa7cf1e6a523df3b13440bd164c0d` | — |
| `Gemma-4-E2B-it-Q8_0.gguf` | `c25a0ce5f49cec613242a6d328a0d417a35a5f77c2c77a97efe942386ad6c5da` | 4,051,836,928 |
| `Qwen2.5-Omni-3B-Q4_K_M.gguf` | `67a70bc7a3b6cc2d5306389472705da1f1323206629b8741d2a454f77f9185c9` | 859,443,200 |

ASR settings in the isolated `config/settings.json` (unchanged copy):
`asrEnabled: true`, `asrEngine: "qwen3-asr-1.7b-local"`.

Hardware at test time: NVIDIA GTX 1080, 8,192 MiB total.

## 2. Verdicts

| # | Check | Verdict | Evidence |
|---|---|---|---|
| 1 | Physical microphone capture | **FAIL** | capture state entered (`mic-button recording`), stimulus played 14.80 s, composer never populated |
| 2 | ASR initialize + transcribe | **PARTIAL — init PASS, transcription FAIL** | `/health` → `{"status":"ok"}`; transcript empty |
| 3 | Transcript insertion / user-controlled send / assistant reply | **FAIL (blocked)** | no chat LLM could be loaded, so no reply path existed |
| 4 | Ordinary typed message with ASR enabled | **BLOCKED** | not run: chat LLM absent; would only append a user turn |
| 5 | Persistence | **PARTIAL — ASR settings PASS, messages NOT TESTED** | ASR applied and ready on both launches from unchanged settings file |
| 6 | ASR + chat coexistence | **NOT DEMONSTRATED** | only one `llama-server` ever existed (the ASR child) |
| 7 | Normal shutdown | **PASS** | app + ASR child gone after close, zero orphans, both runs |

### 2.1 Run 1 (chat-load attempts) — PID 26292

- `llama-server` **4732**, parent 26292:
  `-m "G:\AI GGUF\voice\Qwen3-ASR-1.7B-Q8_0.gguf" --mmproj "…mmproj-Qwen3-ASR-1.7B-Q8_0.gguf" --host 127.0.0.1 --port 64537 -c 2048 -ngl 99 --no-warmup -t 6`
- Picker selected `Gemma-4-E2B-it-Q8_0 (3.8 GB)`; Load produced, in 0.39 s:

```
E llama_model_load: error loading model: tensor 'blk.21.ffn_up.weight'
  data is not within the file bounds, model is corrupted or incomplete
E llama_model_load_from_file_impl: failed to load model
```

  Recorded as `LLM-SRV-CRASH` / `llm.load` at 16:28:37 and 16:35:03 and
  16:38:02 in `config/logs/howl-errors-20261105.jsonl`; full invocation is in
  `logs/model-load-failures.log`.
- `Qwen2.5-Omni-3B-Q4_K_M (0.8 GB)` was the only model small enough to
  plausibly coexist with ASR on 8 GB, but its row rendered below the window
  (`rect=1801,1241 … OFFSCR`). Wheel scrolling over the flyout did not move the
  list; after maximizing (window 2576×1456) the row remained off-screen and a
  coordinate click did nothing, so the trigger stayed on Gemma and Load simply
  repeated the corrupt-model failure. **No chat server was ever started.**
- VRAM while only ASR was loaded: 3,854–4,082 MiB used of 8,192 MiB.
- Normal close: `CLOSE=True`, `EXITED=True`, then **no** `howl-chat` or
  `llama-server` process; ASR 4732 gone.

### 2.2 Run 2 (relaunch persistence + microphone attempt) — PID 25660

- ASR child **10372**, parent 25660, same model and projector, port 56448;
  `curl http://127.0.0.1:56448/health` → `{"status":"ok"}`.
  This is the saved→applied→ready chain: the unchanged isolated settings file
  selected `qwen3-asr-1.7b-local`, the backend launched exactly that pair, and
  the server answered health.
- Capture: `mic-button` class changed to `input-icon mic-button recording`
  after the first invoke, so capture start was accepted. Default render was
  switched from the headset to `Speakers (Realtek(R) Audio)` at 100 % (same
  method as the previous Unit C), `stimulus2.wav` played for 14.80 s, capture
  stopped. Composer value stayed empty across eight polls (24 s).
- Error recorded: `ASR-TRANS-FAIL` / `asr.transcribe` /
  `all recognizers failed, last error: llama-server returned an empty transcript`
  (16:41:31).
- Audio default was restored to `Headphones (WH-1000XM4)` at 54 % immediately
  after the attempt.
- Normal close: `CLOSE=True`, `EXITED=True`, no leftover processes.

No transcript was produced, so there is no observed transcription and no
accuracy claim of any kind.

## 3. Direct UI evidence

Under `build/tmp/unit-c-rerun/` (copied to
`build/unit-c/unit-c-rerun-evidence/`):

- `ui-1.png` — isolated app, ASR active, "NO MODEL LOADED".
- `ui-2.png` — model picker open with the inventory (no `ui-3.png` was
  captured).
- `ui-5.png` — maximized window with picker still unable to expose the small
  models.
- `ui-6.png` — composer empty after the stimulus attempt.
- `ui-tree-1..9.txt` — UIA trees with element rectangles and off-screen flags.
- `run1.out`, `run1.err`, `run2.out`, `run2.err`, `logs/model-load-failures.log`,
  `config/logs/howl-errors-20261105.jsonl`, `real-before.txt`, `real-after.txt`.

Representative commands:

```
Start-Process howl-chat-unit-c.exe -WorkingDirectory <isolated root>   (with env above)
uia.ps1 -Action invoke -Class "mic-button"
AudioCtl.ps1 -Action setdefault -Match "Speakers (Realtek" -Roles 0,1,2
PlayWav.ps1 -Path "build/unit-c/stimulus2.wav"
uia.ps1 -Action value -Class "input-field"
Get-CimInstance Win32_Process -Filter "Name='llama-server.exe'"
nvidia-smi --query-gpu=memory.used,memory.total --format=csv
$p.CloseMainWindow()
```

## 4. Deviations and limitations

1. **No chat LLM available.** The smallest installed chat-eligible model is a
   corrupt file; the next candidate could not be reached through the picker.
   Coexistence, assistant response, transcript send and ordinary text send are
   therefore untested, not passed.
2. **Picker overflow.** Model rows below the viewport were reported `OFFSCR`
   and were not clickable or scrollable through the existing UIA driver. This
   is a tooling limitation of this rerun; it is not established whether a human
   can select them.
3. **Microphone result is inconclusive as to cause.** Capture initialized and
   the ASR server was healthy, but the captured audio yielded an empty
   transcript. Possible causes (headset not worn / room audio not reaching its
   mic, muted capture endpoint, capture-path defect) were not distinguished.
   The model and projector loaded, so this is not a model-load failure.
4. **Stimulus method.** Audible playback through the Realtek speakers replaced
   the headset as default render for ~16 s; the previous default and volume
   were restored. The capture endpoint itself was not changed.
5. **Test surface.** Verification used a copy of the accepted packaged build
   with a junctioned `runtime/`; the in-tree install was never launched.
6. No messages were persisted, so message-level persistence and history
   round-trip remain unverified.

## 5. Comparison with the previous Unit C

| Aspect | Previous Unit C (2026-11-04/05) | This rerun |
|---|---|---|
| Chat server | Never healthy; `llm.infer.character` failed with `localhost:8080` refused at 03:27:48 and 03:35:55 | Never started; corrupt Gemma file, picker blocked the 0.8 GB model |
| Assistant response | None observed | None; not attempted because no model |
| ASR transcription | Succeeded once from speaker→headset-mic stimulus (`check1c`) | Empty transcript from the same stimulus method |
| ASR ownership / leftover | `llama-server` 17740 orphaned after app exit (check6/6b) | ASR child 4732 and 10372 both terminated with the app; zero orphans |
| Protected data | Rewrote `history.json` and sibling chat files | Real AppData byte-identical; all writes confined to the isolated profile |
| Settings | `asrEnabled`/`asrEngine` changed on disk | `settings.json` hash unchanged |

The D3 ownership fix is observable here: two independent normal-shutdown
cycles left no ASR descendant, including a relaunch where the settings-driven
ASR child started automatically.

## 6. Status

Unit C is **not accepted and not closed.** Checks 3, 4 and 6 are blocked by the
absence of a loadable chat LLM; check 1/2 transcription failed with an empty
transcript; checks 5 (settings) and 7 passed. No repair was attempted, per the
verification-only authorization. Awaiting owner review and a decision on:

- whether to supply/repair a chat model that fits 8 GB alongside ASR, and
- whether the empty microphone transcript should be investigated as a capture
  issue or retried with the headset confirmed active.
