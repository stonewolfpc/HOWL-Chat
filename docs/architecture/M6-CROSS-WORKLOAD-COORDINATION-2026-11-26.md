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

## Owners

Unchanged. Admission decides, the lifecycle gate serialises transitions, the
registry is the ledger, the reclaim path is the only place another workload may
be stopped (and only with the user's confirmation). No new owner was added; the
two fixes are message-selection and a guard on an existing rescue path.

## Verified limitations

- **Packaged host not driven** for this scenario; all evidence is the browser
  dev host over the same bindings.
- **No ASR run.** ASR needs a microphone and a human speaker; the audio leg is
  Piper synthesis only.
- **No generation in this scenario.** The image model was loaded, not sampled;
  generation and residency measurement for image were covered by M3–M5.
- **The explicit image→chat RPC was not invoked in this run.** The return to
  chat was exercised through the primary-workload path in `LoadModel`, which calls
  the same `FreeVRAMFor → reclaimOther → waitForAdmissionFit` sequence. The
  release, the fresh admission, the restored placement and the residency record
  were all observed; only the separate RPC entry point was not clicked.
