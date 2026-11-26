# M6 — Cross-workload coordination

Date: 2026-11-26
Isolated root: `C:\HOWL-Test\hra-image-20261125` (`.howl-isolated-root`)

Two real workloads were run against one card to find out whether Howl coordinates
them or merely serialises them. It did both, and the run exposed two defects.

## Scenario (browser dev host, real engines)

1. Audio activated through the production binding; Piper synthesised a WAV
   (`data/cache/tts/tts_…wav`) — CPU-only, no VRAM claim.
2. Chat loaded (Qwen3.5-9B): `admitted: 6057 MB VRAM required`; engine placement
   32/33 layers; residency recorded `planned=6057 engine=4887 device_delta=5209
   ledger=5209 within plan`.
3. Image model requested while chat was resident.

## Defect 1 — the M5 rescue silently degraded a contended start (fixed)

With chat holding the card, admission refused the image model and named chat as
the blocker. The M5 bounded-placement rescue then re-planned inside the leftover
budget and started SDXL **on the CPU**:

```
[App] bounded placement admitted: mode=cpu-only budget=1041MB
Starting sd-server ... -m Illustrious-XL-v2.0.safetensors --max-vram 1 --backend cpu --vae-tiling
```

Running a GPU-capable model on the CPU because another workload holds the GPU is
not a placement decision; it is a handoff decision the user has not made. The
rescue now applies only when nothing resident is blocking:

```
[App] image admission: admitted=false ... blocking model: llm "Qwen3.5-9B…" holds 5209 MB VRAM
```

and no engine starts. Pinned by
`TestResidentBlockerIsNamedRatherThanDegradedToCPU` and, in the other direction,
`TestBoundedPlacementStillAppliesWithNoBlocker`.

## Defect 2 — the handoff dialog blamed a CPU resident (fixed)

The unconfirmed handoff prompt read:

```
Your chat model "piper" is loaded and this one needs about 6945 MB, with 7090 MB free of 8192 MB.

Unload it and continue?
```

`piper` is a 360 MB CPU speech model with no VRAM claim, named as the chat model
blocking a 7 GB VRAM request. The dialog now selects the first resident that
actually holds GPU memory, and labels each workload by type (chat, image, speech
recognition, speech). With only a CPU resident left it states the shortfall
instead of blaming it:

```
This model needs about 6945 MB and only 7089 MB of 8192 MB is free, so it cannot start.
```

Pinned by `TestFitMessageBlamesTheGPUHolderNotACPUResident` and
`TestFitMessageNamesEachWorkloadKind`.

## Verified coordination

| Step | Result |
|---|---|
| Audio (Piper) beside chat | Synthesised while chat resident; no VRAM claim, no admission refusal |
| Image request while chat resident | Refused, blocker named, no silent CPU degradation |
| Unconfirmed handoff | `unloading the chat model to start the image model was not confirmed; nothing was unloaded` |
| Confirmed handoff | `INFO: llama-server stopped, model memory freed`; image then admitted at `required=5084 MB` and sd-server started on the GPU |

## Defect 3 — the handoff reports a timeout after it has succeeded (fixed)

The confirmed handoff once returned:

```
VRAM reclamation for image did not complete within 10s: This model needs about 6945 MB
and only 7089 MB of 8192 MB is free, so it cannot start.
```

The unload had in fact succeeded — the image loaded on the GPU immediately
afterwards. The reclaim loop waited for the *full-precision* request to admit,
while the load path that follows falls back to q8 (5084 MB), so the wait could
never observe the condition it was waiting for.

The fix is not a longer wait and not a looser condition. The reclaim wait now
holds the same **placement ladder** the load path will run — stored precision
first, then the automatic q8_0 conversion when that policy allows it — and
admits each candidate against fresh snapshots. Success is an admission of one
candidate on live telemetry; a stopped process is still never treated as freed
memory, and when no candidate fits the handoff still fails and says so
(`TestReclaimStillTimesOutWhenNothingFits`). The result names the placement it
admitted, so the caller is not told a full-precision load is coming when a q8
conversion will run:

```
Fits: needs about 5084 MB and 7093 MB is free of 8192 MB.
; this load will use the automatic q8_0 conversion
can_fit=true weight_type=q8_0
```

Pinned by `TestReclaimSucceedsOnTheSameLadderTheLoadUses`.

## End-to-end transition: chat → image → chat, with audio present

One session, three workloads, in order:

| Step | Evidence |
|---|---|
| Piper synthesises | `data/cache/tts/tts_1795661005018019600.wav`; CPU-only, no VRAM claim |
| Chat loads (primary) | `admitted: 6057 MB VRAM required`; `actual placement layers=32/33 context=3328 matches_plan=true`; residency `planned=6057 engine=4887 device_delta=5210 ledger=5210 within plan` |
| Confirmed handoff to image | `can_fit=true weight_type=q8_0`; `llama-server stopped, model memory freed` |
| Image loads and generates | 512×512, 2 steps, `status=completed`; residency `planned=5093 engine=0 device_delta=5681 ledger=5681 above plan` |
| Return to chat | `[VRAM] unloading image model to make room for the llm model`; chat re-admitted at the same plan; `layers=32/33 context=3328 matches_plan=true`; residency again `planned=6057 engine=4887 device_delta=5210 ledger=5210 within plan` |
| Audio after both transitions | Synthesis still succeeds — the CPU workload was never evicted by either handoff |

The chat return is the primary-workload path: loading chat is itself the user's
intent, so it hands the device over through the same `reclaimOther` sequence
with confirmation already given, rather than asking again mid-load. The image
direction cannot start that way — it requires the explicit confirmation shown
above.

## Evidence closure 1 — explicit image→chat entry point (browser host)

The primary-load path had been proven; the separate public entry point had not.
With the image model resident, `FreeVRAMFor("llm", <chat path>, true)` was called
directly over the harness RPC:

```
error: none
can_fit: True required: 6182 free: 7090
message: Fits: needs about 6182 MB and 7090 MB is free of 8192 MB.
[VRAM] unloading image model to make room for the llm model
[LoadModel] actual placement: known=true layers=32/33 vram=4677MB ram=679MB context=3328 matches_plan=true
```

So the entry point behaves like the mechanism it shares: refuse without consent,
stop the resident, wait on fresh telemetry, admit, and report. The 6182 MB here
is the pre-plan request (no fitted plan exists before the load), which is why it
differs from the 6057 MB the fitted load reports.

## Evidence closure 2 — packaged host pass (WebView2, real input)

Driven with real mouse clicks and keystrokes against the maximised window of the
packaged app, on the same isolated root, with screenshots read back after every
step:

| Clicked | Visible afterwards |
|---|---|
| Chat panel, typed a message, Enter | Chat loaded: `actual placement layers=32/33 context=3328 matches_plan=true` |
| Image Generator icon | Panel with `Illustrious-XL-v2.0 (sdxl)` selected, 512×512 |
| GENERATE | Native `wails.localhost` message box: *Your chat model "Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M.gguf" is using about 5183 MB, and this one needs about 6945 MB — more than the 1694 MB free of 8192 MB. Unload it and continue?* |
| OK | `[VRAM] handoff complete: Fits: needs about 6057 MB and 6985 MB is free of 8192 MB.`; sd-server started, a new lighthouse image appeared in the gallery |
| Chat icon, typed a message | Chat reloaded; `[LoadModel] actual placement: layers=32/33 context=3328 matches_plan=true` |

Residency measured in this host: image `planned=5093 engine=0 device_delta=4501
ledger=4501 within plan`; chat `planned=6057 engine=4887 device_delta=5223
ledger=5223 within plan`.

**Host difference recorded:** the handoff confirmation is a native Windows
message box in the packaged app, not an in-page dialog. The browser host has no
such surface, so the two hosts ask the same question through different chrome.

## Defect 4 — the first message after a handoff reload can fail (open)

In the packaged host, the message typed immediately after returning to chat
produced this in the transcript while the reload was still starting:

```
[Error: failed to send request: Post "http://localhost:8080/v1/chat/completions":
 dial tcp [::1]:8080: connectex: No connection could be made because the target
 machine actively refused it.]
```

The model was listening moments later; the next message streamed normally. The
send path begins its request as soon as the user submits, while the reload the
nav transition had begun was still bringing llama-server up, so the first
request can reach a port that is not open yet.

A defensive guard was written (wait for `/health` while the engine reports a
loading stage) and then **reverted**: in this failure the send left before the
load had even entered a loading stage, so the guard never armed and the error
still appeared. The cause is ordering between the load request and the send,
not a missing wait inside the send. It is recorded open rather than patched with
something the host evidence does not support.

## Intent-driven activation (chat)

Submission is the intent to execute, so the send owns the whole sequence:
load (admission decides, a needed handoff goes through the same reclaim
route), verify the engine answers, then submit. Navigation and selection stop
moving GPU memory.

- `checkModelStatus` no longer calls `autoLoadDefaultModel` — that function is
  deleted, not bypassed. Opening the chat page shows `NO MODEL LOADED`.
- Both send paths (`streamAssistantResponseFromBackend` and
  `…WithImage`) begin with `ensureChatModelReadyForSend()`, which awaits
  `App.EnsureChatModelReady(modelPath)` before creating the bubble or the request.
- `EnsureChatModelReady` refuses an empty selection
  (`no chat model selected, so there is nothing to load for this message`),
  otherwise loads through `LoadModel` (admission, lifecycle gate, confirmed
  handoff) and then waits for `/health`. Readiness is an answering engine, never
  a started process (`TestWaitForHealthRequiresAnAnswer`).
- The explicit Load control and model selection are unchanged; they are the
  deliberate preload path.
- Image already worked this way: `Generate` → `admitImageLoad` → `EnsureReady`,
  with `SetActiveImageModel` reserved for the explicit Load/Reload button. No
  page or selection change there starts sd-server.

### Verified in the packaged host

| Action | Observed |
|---|---|
| Launch, click Chat | header `NO MODEL LOADED`; `Get-Process llama-server,sd-server` empty; zero `LoadModel` lines in the console log |
| Type "demand activation check", Enter | `[Chat] demand load: no model is loaded, loading "…Q4_K_M.gguf" for this message`; `[LoadModel] actual placement: layers=32/33 context=3328 matches_plan=true`; stream completed, 486 chunks |
| Same message before this change | failed instantly with `dial tcp [::1]:8080 … actively refused it` |

The refused-connection error did not recur: the request now leaves only after
the engine answers, because the send itself started it.

**Observed but not explained:** the assistant bubble rendered empty in the
packaged host even though 2,285 bytes of content streamed (`Stream completed,
total chunks=486, total content=2285 bytes`). A second send behaved the same. The
same empty bubble appeared in the packaged host *before* this change, so it is a
separate rendering/host issue, not evidence for or against demand activation —
and it is not claimed as a working reply.

### Safeguards kept

- A refusal (no selection, admission says no, load fails) surfaces and the send
  does not submit; nothing is retried automatically.
- One send awaits one ensure, so a failed or cancelled send cannot cause a
  second execution of the same message.
- Startup failure is reported by `LoadModel` and by the readiness wait; neither
  is converted into a success.

## Owners

Unchanged. Admission decides, the lifecycle gate serialises transitions, the
registry is the ledger, the reclaim path is the only place another workload may
be stopped (and only with the user's confirmation). No new owner was added; the
two fixes are message-selection and a guard on an existing rescue path.

## Verified limitations

- **Packaged host pass covers chat and image only.** It was not extended to a
  Piper synth or an ASR capture; audio was verified in the browser host.
- **No ASR run at all.** Capture devices are present (`CORSAIR VOID ELITE
  Wireless Gaming Headset`, `Realtek High Definition Audio`), but transcription
  needs a human speaker at the microphone. Nothing here claims ASR works beside
  a resident workload.
- **Defect 4:** the refused-connection race did not recur after Send took
  ownership of the load, but **the assistant reply rendered empty** in that same
  packaged run, so the end-to-end reply is still unproven in that host.
- **Image generation in the packaged host used the 20-step default at 512×512**;
  the browser run used 2 steps. Neither is a speed measurement.
- **Screenshot/click driving is physical-screen input**; the packaged window was
  resized to 1400×900 because the maximised window's DPI scale made coordinate
  mapping unreliable. Every reported step was read back from a screenshot.
