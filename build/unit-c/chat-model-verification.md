# Unit C — chat model selection and load verification

Status: **ACCEPTED by owner on 2026-11-05. Unit C is closed; no further testing or reopening.**
Result at acceptance: PASS for model selection, load, conversation turn, persistence and active
ASR + chat coexistence.
Date: 2026-11-05
Authorization: owner asked for an inventory audit (integrity, architecture, memory, actual picker
registration), then selection of a viable candidate and verification of loading in the isolated
environment. No code, UI, settings-architecture, D1/D3 or TTS changes were made.

## 1. Isolation

- App root: `build/tmp/unit-c-rerun/` (own `config/`, `runtime/` junction, `profile/`, `logs/`).
- Redirected: `USERPROFILE`, `APPDATA`, `LOCALAPPDATA`, `TEMP`/`TMP` into `profile/`.
- Protected data checked before/after every run: `APPDATA/HOWLChat` file hashes diff against
  `build/tmp/unit-c-rerun/real-before.txt` → `REAL_IDENTICAL`.
- `config/settings.json` (real, protected): `e786f53f2fd7f91b794936dad1930e2b09594de8a81167179a857ecb011757bf`
  — unchanged after the run.

## 2. Inventory audit

Full structural scan of `G:\AI GGUF\voice` and `G:\AI GGUF\GGUF`:
`build/unit-c/model-audit/inventory.json` (tensor table bounds per file; `data_ok` = last tensor ends
inside the file).

Corrupt / truncated (excluded):

| File | Finding |
| --- | --- |
| `Gemma-4-E2B-it-Q8_0.gguf` | `blk.21.ffn_up` overruns file by 1,080,320 bytes |
| `Qwen2.5-Omni-3B-Q4_K_M.gguf` | overrun +20,903,232 bytes |
| `Qwen2.5-Omni-3B-mmproj-Q8_0.gguf` | overrun; runtime load fails in <1 s |
| `Ultravox-Llama-3.2-1B-Q4_K_M.gguf` + `mmproj-Ultravox…` | overrun |
| `Gemma-4-E2B` runtime probe | `load_model: failed to load model` (see `failure-log-tail.txt`) |

Intact chat-eligible candidates (all `chat_template: true`, all present in the app registry):

| Model | Bytes | Arch | Layers | CPU probe |
| --- | --- | --- | --- | --- |
| LFM2.5-8B-A1B-heretic.Q4_K_M | 5.16 GB | `lfm2moe` | 24 | health ok, 1.5 s |
| Qwen3.5-9B-Uncensored…Q4_K_M | 5.63 GB | `qwen35` | — | health ok, 3 s |
| **Gemma-4-E4B-Uncensored…Q4_K_M(1)** | **5,335,285,728** | **`gemma4`** | **42** | **health ok, 3 s** |
| LongWriter-glm4-9b | 5.30 GB | `chatglm` | — | not probed |
| command-r7b | 5.06 GB | `cohere2` | — | not probed |
| 12B llama family (7.48 GB each) | — | `llama` | — | too large to coexist with ASR |

Selected file:
`G:\AI GGUF\GGUF\Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M(1).gguf`
size 5,335,285,728 B, mtime 2026-09-04 20:53:48 -0500,
sha256 `7096c0dd243da9762247386d476801714971e53e6894b4001710287a522fea1b`,
720 tensors, embd 2560, native ctx 131072, ftype 15 (Q4_K_M), `data_ok: true`, `problems: []`.

## 3. Why earlier picker attempts failed (resolved, no repair)

Two independent causes, both environmental to the test harness:

1. **Clicks landed before the model rescan finished.** The app scans `G:\AI GGUF\GGUF`
   (~60.2 s, 27 entries) then `G:\AI GGUF\voice` (~12.6 s, 6 entries) — 72.9 s total. Earlier
   attempts clicked the trigger 2–5 s after launch, when the selector had no rows.
   Evidence: `config/logs/startup.log` `rescan: done G:\AI GGUF\GGUF 27 entries`,
   `rescan: done G:\AI GGUF\voice 6 entries`.
2. **The dropdown is viewport-clipped.** Once the scan completes the trigger opens
   `.model-select-options` (tree 217 → 272 elements) and it renders under
   `cls=model-select-option` (not `model-select-options`). Rows below the clip carry `OFFSCR`;
   wheel events earlier were sent outside the list box.

Registration evidence (this run, `model-audit/picker-open-final.txt`): all 27 GGUF entries are
listed with human sizes, e.g. `Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M(1) (5.0 GB)` at
`rect=1801,889,118,16`. So candidate registration is proven through the existing supported
backend route; no UI/menu change was needed.

LFM2.5 row sits just below the dropdown viewport (`rect=1801,1087` `OFFSCR`). Wheel input inside
a visible row moved the list (272 → 276 elements) but left LFM `OFFSCR`, so LFM was not selectable
by pointer in this harness. Recorded as a harness limitation; nothing was changed in the UI.

## 4. Load verification

Route: picker trigger click → visible row click → `load-button` invoke. No CDP, no storage seeding,
no code change.

Timeline (isolated app PID 19212, ASR child 2288 on port 60667):

| t | Observation |
| --- | --- |
| 0 s | app + ASR child running, VRAM 3,924 MiB / 8,192 MiB |
| 10 s | second `llama-server` present; `:8080/health` → `503 {"message":"Loading model"}` |
| **12 s** | **`:8080/health` → `{"status":"ok"}`** |
| — | `/v1/models` reports `G:\AI GGUF\GGUF\Gemma-4-E4B-…Q4_K_M(1).gguf` |
| — | VRAM 4,124 MiB / 8,192 MiB |
| — | UI shows **READY**, `UNLOAD` enabled (`model-audit/loaded-e4b.png`) |

Coexisting chat `llama-server` (PID 7116, parent 19212):

```
llama-server.exe --host 127.0.0.1 --port 8080
  -m "G:\AI GGUF\GGUF\Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M(1).gguf"
  -c 3096 --threads 8 -b 256 -ngl 0 --fit off
  --chat-template "{%- for message in messages %}<start_of_turn>{{ message.role }} {{ message.content }}<end_of_turn> {%- endfor %}<start_of_turn>model "
  --reasoning on --reasoning-format deepseek --reasoning-budget 2048
  --reasoning-budget-message "I have thought enough. Time to answer."
```

ASR child unchanged and still resident:
`-m "G:\AI GGUF\voice\Qwen3-ASR-1.7B-Q8_0.gguf" --mmproj "…mmproj-Qwen3-ASR-1.7B-Q8_0.gguf"
--host 127.0.0.1 --port 60667 -c 2048 -ngl 99 --no-warmup -t 6`.

**Important memory finding:** the app clamped the chat model to `-c 3096 -ngl 0 --fit off`, i.e. it
declined GPU offload while ASR holds ~3.9 GB of the 8 GB card and ran the chat model entirely on
system RAM. Coexistence therefore passes on memory, but chat inference is CPU-only and its
generation speed is not yet measured.

## 5. Shutdown and isolation

`CloseMainWindow()` → app exited, `tasklist` shows no `howl-chat*` and no `llama-server`
(both children reaped), protected-data diff `REAL_IDENTICAL`, release `settings.json` hash
unchanged. `logs/model-load-failures.log` gained no new entry for this model — the only lines are
the earlier `Gemma-4-E2B` probe.

## 6. Verdict

| Check | Result |
| --- | --- |
| Model file integrity (tensor bounds) | PASS |
| GGUF version/arch admissions (`gemma4`, v3) | PASS |
| Chat template present | PASS |
| Registry + picker registration | PASS (27 entries listed) |
| Selection through existing supported route | PASS |
| Load alongside Qwen3-ASR (both servers up) | PASS (12 s) |
| `:8080/health` = ok, `/v1/models` = selected path | PASS |
| VRAM coexistence within 8,192 MiB | PASS (4,124 MiB; chat on CPU) |
| GPU offload of chat model | FAIL / degraded (`-ngl 0 --fit off`) |
| Clean shutdown, no leftover processes | PASS |
| Protected-data isolation | PASS |
| Actual assistant response (typed message) | PASS (`PONG`, 13 s round trip) |
| Typed-message regression (user bubble → backend → reply → save) | PASS |
| Persistence across restart | PASS (both turns restored from disk) |
| ASR + chat simultaneous use in one session | PASS (section 10; supersedes the earlier INCONCLUSIVE row) |
| Microphone dictation transcribed into the composer, both models loaded | PASS (section 10.2) |
| Dictated message sent through `.send-button` → assistant reply | PASS (6 s) |
| Second, typed message in the same session → assistant reply | PASS (`LEDGER`, 2 s) |
| Restart restores all six turns with the model unloaded | PASS (section 10.4) |

## 7. Conversation workflow evidence (same isolated build)

Run: app PID 26872; ASR child 18260 (`--port 58738`), chat server 21568 (`--port 8080`).
Model loaded as described above, then one typed message sent through the real UI:

1. `message-input` set to `Reply with exactly this word and nothing else: PONG`, value read back
   (`VALUE=…PONG`) before sending.
2. `.send-button` invoked. No character was loaded, so the generic-assistant path ran
   (`chat-stream.js`: "No characters loaded, using generic AI response"), which is the documented
   fallback, not an error state.
3. Reply completed at t=12 s after send. UI tree: user bubble
   `Reply with exactly this word and nothing else: PONG` at `rect=1089,565`, assistant bubble
   `PONG` at `rect=1614,668`. `model-audit/conversation.png`.
4. Both `llama-server` processes still resident at reply time; VRAM 4,161 MiB / 8,192 MiB.

Persistence (`profile/AppData/Roaming/HOWLChat/`):

```
config.json      {"lastChatId": "2026-11-05T18-01-55Z"}
chats/2026-11-05T18-01-55Z/meta.json       updatedAt 2026-11-05T18:02:08Z
chats/2026-11-05T18-01-55Z/history.ndjson  (source of truth)
  {"id":"0","sender":"User","role":"user","text":"Reply with exactly this word and nothing else: PONG","timestamp":"2026-11-05T12:01:55-06:00"}
  {"id":"1","sender":"AI","role":"assistant","text":"PONG","timestamp":"2026-11-05T12:02:08-06:00"}
chats/2026-11-05T18-01-55Z/history.json   canonical snapshot, 360 B
```

File hashes after the turn and after the restart are identical:
`history.ndjson 46fb4a2321382469a677347ffc391e2fd2dfb7590d8bafe0312c18199c6ab1c9`,
`meta.json 30880d2091be79858e94ed42e8f206f7a26c47a082d1d970cd308a36430a3f3c`,
`config.json e7adbd279d1053cc3e5f88e8786873d7118f766bd56bc01b323424ab78f7d63c`.

Restart check (new app PID 29452, same isolated profile): after the model rescan the UI restored
both bubbles from disk (`Reply with exactly … PONG` at `rect=901,398`, `PONG` at `rect=1426,504`)
with **NO MODEL LOADED** and only the ASR server running (port 52611). So the transcript is
independent of the inference server, and the chat model is *not* auto-loaded on startup even
though the selector remembers the last choice — relevant to the upcoming resource-management
design, not a defect on its own.

App closed, no `howl-chat*` or `llama-server` left, protected data `REAL_IDENTICAL`,
release `settings.json` hash unchanged.

## 8. Limitations

- LFM2.5-8B could not be selected by pointer (row clipped below dropdown viewport); Gemma-4-E4B was
  chosen instead. Both were intact; LFM remains untested in-app.
- Chat ran CPU-only, so nothing here demonstrates usable generation latency.
- Duplicate picker entries exist for the same file (`…Q4_K_M` and `…Q4_K_M(1)`); harmless for this
  verification, not investigated.
- No assistant reply was produced or inspected; the model is loaded but its output path is unproven.
- Transcription accuracy and typed-chat regression remain outside this record.

## 9. Next decision

Unit C stays open pending owner review. The conversation turn, typed-message regression, persistence
and ASR + chat *actively used in one session* are now proven (section 10), so the item that was
outstanding when this paragraph was first written is closed.
No further model audit is needed. A GPU-offloaded chat run would require freeing ASR VRAM or
lowering `gpu_layers`/`context_size`, which changes the tested configuration and needs an explicit
decision; the 13 s CPU-only round trip is adequate evidence that inference works but says nothing
about sustained throughput. Resource management (load/park/unload policy, VRAM budgeting,
auto-load-on-start, user-visible simplicity with advanced controls opt-in) is the next engineering
problem and should be designed before frontend assumptions are fixed.

## 10. Active ASR + chat coexistence — final Unit C step (2026-11-05)

Purpose: physical-microphone dictation transcribed, sent, and answered by the chat model while both
ASR and chat were loaded and operating in one session, then restored after a full restart. No
application code, UI, menu, settings, D1/D3 record or TTS path was changed; every write stayed inside
`build/tmp/unit-c-rerun/`.

### 10.1 Session

| Item | Value |
| --- | --- |
| App process | PID 21040, `howl-chat-unit-c.exe`, isolated root |
| ASR server | PID 26068, `--port 54935`, `Qwen3-ASR-1.7B-Q8_0.gguf` + `mmproj-Qwen3-ASR-1.7B-Q8_0.gguf` |
| Chat server | PID 2200, `--port 8080`, `Gemma-4-E4B-Uncensored-HauhauCS-Aggressive-Q4_K_M(1).gguf` |
| ASR readiness | 0.1 s after the child was discovered (child spawned within ~1 s of launch) |
| Chat readiness | 14.1 s after `Load` (`:8080/health` = ok) |
| VRAM | 3,931 MiB startup → 4,169 both loaded → 4,455 after dictation → 4,483 after the typed turn (8,192 MiB total) |
| GPU utilisation | 30–31 % during generation (chat is CPU-bound: `-ngl 0 --fit off`) |
| App RAM | 132 MB working set |
| Selection route | existing picker; option rows are unnamed wrappers and the model name is in child text, so the row was chosen by clicking its on-screen coordinates |

### 10.2 Acoustic path and dictation

- Capture endpoint: `Headset (WH-1000XM4)` (Bluetooth HFP), Windows default for console and
  communications, microphone consent `Allow`.
- The agent cannot speak, so the spoken stimulus was synthesized once with SAPI (`Microsoft David`,
  16 kHz mono, 4.06 s, peak 31,000 = −0.5 dBFS) and played **acoustically** through
  `Speakers (Realtek(R) Audio)` at 100 % into the headset's physical microphone. The default render
  endpoint was restored to `Headphones (WH-1000XM4)` immediately afterwards.
- The owner was physically present during this run and spoke spontaneously while the microphone was
  capturing — calling for River, then saying `Sorry, sweetheart!`. That utterance was picked up by
  the same recording and appears in the transcript ahead of the synthesized sentence.
- A first attempt routed the stimulus to the A2DP `Headphones` endpoint while the headset was
  recording on HFP: the headset produced no audible output, the recognizer returned unrelated text
  (`Yeah, that's only one of the functions. All, all close enough to one hundred percent. Everything
  else is mostly well.`), and nothing matched. Recorded as a **failed stimulus**, not an application
  failure; it is the reason the speakers endpoint was used.

Spoken text:

```
The blue lantern guild meets at seven in the copper harbour.
```

Transcript inserted into the composer (read back through ValuePattern, then sent unedited):

```
Sorry, sweetheart. The Blue Lantern Guild meets at seven in the Copper Harbor.
```

Every intended word is present. The leading `Sorry, sweetheart.` is **not** recognizer noise: the
owner was in the room speaking and said exactly that while the microphone was recording. An earlier
revision of this record attributed the prefix to a recognition artifact on the synthesized voice;
that reading was incorrect and is corrected here.

The transcript therefore contains one **uncontrolled human-speech capture** — the owner's spontaneous
remark, transcribed verbatim — followed by the intended synthesized sentence. It was sent unedited,
as recorded. This is an incidental observation, not a controlled human-speech accuracy test: there is
no reference recording, no fixed script and no word-error measurement against human speech.

### 10.3 Send, reply, and the typed regression

| Step | Evidence |
| --- | --- |
| Send | `.send-button` invoked 13:14:18.7, normal user action, no scripted bypass |
| Assistant reply | `Copper Harbor. Seven sharp. I won't be late again, darling. Promise.` saved 13:14:24 (≈6 s) |
| Typed message | `Reply with exactly this word and nothing else: LEDGER`, value read back, sent 13:16:57 |
| Assistant reply | `LEDGER` saved 13:16:59 (≈2 s) |
| Residency | `SERVERS[after-dictation] count=2` — chat `:8080` and ASR `:54935` both alive; `READY` badge, `UNLOAD` present |
| Screenshots | `coexistence/both-loaded.png`, `coexistence/dictated-reply.png`, `coexistence/typed-reply.png` |

### 10.4 Persistence and restart

`history.ndjson` after both workflows (first two records are the earlier PONG turn in the same chat):

```
{"id":"2","sender":"User","text":"Sorry, sweetheart. The Blue Lantern Guild meets at seven in the Copper Harbor.","role":"user","timestamp":"2026-11-05T13:14:18-06:00"}
{"id":"3","sender":"AI","text":"Copper Harbor. Seven sharp. I won't be late again, darling. Promise.","role":"assistant","timestamp":"2026-11-05T13:14:24-06:00"}
{"id":"4","sender":"User","text":"Reply with exactly this word and nothing else: LEDGER","role":"user","timestamp":"2026-11-05T13:16:57-06:00"}
{"id":"5","sender":"AI","text":"LEDGER","role":"assistant","timestamp":"2026-11-05T13:16:59-06:00"}
```

Hashes: `history.ndjson 5DFFB392D4DEE71D4476089D8B05AFED96BD0A06DB13912100D67BA0653321CA`,
`meta.json 0B166359A0BA375E168B036A1A45B3CCB460A98ED08A5A9F27E86E6741BF1148`,
`history.json 9A114016033B3E09617F07CE6E8CA09BCF4CD221D34E51E8B84B2F1E38A80D6E`,
`config.json E7ADBD279D1053CC3E5F88E8786873D7118F766BD56BC01B323424AB78F7D63C` (unchanged,
`lastChatId` still this chat). `meta.json.updatedAt` = `2026-11-05T19:17:25Z`.

Restart (new app PID 25592, same isolated profile): the 302-element tree contains `PONG` ×2,
`Copper Harbor.`, the full reply sentence, and `LEDGER` ×2 — all six bubbles restored from disk with
no message resent and no inference run. `NO MODEL LOADED` / `No model loaded` shown, only the ASR
server resident, `:8080` absent (`RESTORE_CHAT_SERVER_PRESENT=False`). Screenshot `restore.png`,
tree `restore-tree.txt`.

### 10.5 Shutdown and isolation

Both runs: `CloseMainWindow()` → app exit, `LEFTOVER=0`, no `howl-chat*` or `llama-server` left.
Real `%APPDATA%\HOWLChat` aggregate (sorted SHA-256 list) was
`9a405586012d7e33f77574b21dea8d5e5a5440f64ae23dd16710945be3fb3d90` before the session and identical
after the restart run. Release `config/settings.json` = `e786f53f…757bf` unchanged. Default render
endpoint and volume restored (`Headphones (WH-1000XM4)`, 62 %).

### 10.6 Verdicts for this step

| Check | Result |
| --- | --- |
| Both models loaded and healthy in one session | PASS |
| Physical microphone delivers speech-grade audio | PASS (transcript matches intended words) |
| Spontaneous human speech captured | OBSERVED (uncontrolled) — owner's `Sorry, sweetheart!` transcribed ahead of the synthesized sentence |
| Transcript inserted into the composer | PASS (ValuePattern read-back) |
| Dictated message sent by user action → assistant reply | PASS (6 s) |
| Second typed message → assistant reply | PASS (2 s) |
| Both backends operational throughout | PASS (2 servers resident before and after) |
| Conversation persisted | PASS (4 new records) |
| Restart restores conversation, model unloaded | PASS (6 bubbles, no inference) |
| Clean shutdown, no orphans | PASS |
| Protected data unchanged | PASS |
| Chat GPU offload | FAIL / degraded (unchanged, not repaired) |
| Captured-audio level measurement | NOT TESTED (temp WAV deleted before sampling) |

### 10.7 Harness limitations (not product failures)

- The app deletes its staged ASR temp WAV before it can be sampled (two attempts, zero files seen),
  so there is no captured-audio level figure for this run. Audio presence is evidenced by the
  renderer's `RMS > 0.002` gate plus a matching transcript.
- The UIA tree exposes message bubbles without a usable class/id, so the harness's class-based reply
  poll logged `REPLY_TIMEOUT` after 150 s although both replies were already on disk (6 s and 2 s).
  Those `REPLY_TIMEOUT` lines are harness noise; generation times here come from history timestamps,
  and the rendered bubbles are confirmed by screenshot.
- ASR latency: a precheck with immediate composer polling measured 3.1 s from stop to inserted text;
  the official run's composer read was delayed by a 45 s capture watcher, so its ~49 s figure is a
  harness bound, not ASR time.
- The stimulus is synthesized (the agent cannot speak), so this verifies the acoustic capture → ASR
  → send → LLM path. The same recording also contains one spontaneous human utterance from the owner;
  that is an observation only, not a controlled accuracy test of human speech.
- Server discovery by parent PID briefly reported `BOTH_RESIDENT_AFTER_DICTATION=False` because of a
  PowerShell expression-precedence bug in the harness; the server dump in the same log shows both
  processes alive, which is the evidence used above.

## 11. Acceptance and locked decisions (2026-11-05)

**Owner accepted Unit C.** The recorded results support the intended behaviour. Evidence in sections
6–10 is unchanged by acceptance; nothing was retested, corrected or reopened.

Two decisions are locked:

1. **A model staying unloaded after restart is correct.** Restart restores the conversation, not the
   inference workload. The application must not auto-load a model until system resource management
   can make that decision intelligently.
2. **GPU offload is deferred.** It belongs to the system resource-management branch, which must handle
   chat, image, video and audio together, rather than patching chat in isolation.

Consequences for the record:

- The `Chat GPU offload | FAIL / degraded` row in 10.6 remains a factual measurement, but it is not a
  Unit C defect: it is deferred work owned by the resource-management branch.
- No automatic model loading is to be added on the basis of the restart observation in 10.4; that
  behaviour is the intended product behaviour.
- Unit C is closed. Further work proceeds under new assignments.
