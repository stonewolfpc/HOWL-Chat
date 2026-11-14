# MEMORY-09 — Local S1 Execution Plan (proposal, not authorisation)

**Status:** **PLAN ONLY.** No call is made by this document, S1 is **not** authorised, and the L3 allowance is spent.
**Date:** 2026-11-14 · **Basis:** `BEA-1` design, `MEMORY-09-OME-1` architecture, L2 identity record, L3 run `eval/runs/L3-2026-11-14T080303Z/`.
**Executor:** single local model, single machine. Nothing here claims independent raters.

---

## 1. Purpose

L3 proved the transport, the capture and the evidence chain. It did **not** prove the harness can run the intended experiment faithfully on a local stack. This plan closes that gap *on paper* so the owner can approve or reject a specific set of calls, rather than approve a direction.

---

## 2. Controls that must not change

These come from the BEA-1 design and are already frozen in code. Any change to them is a **new experiment version**, not an adaptation.

| Control | Frozen value | Pinned evidence |
|---|---|---|
| Fixture | `F-ST-01-S1` v1 | hash `3dd41c8e…a6412d` |
| Memory state | canonical JSON, hash-on-mutation | state hash `aa98d95f…9d5700` |
| Conditions | A (assisted) / B (unassisted), one system line apart | A: `076de19e…750af7`, B: `97213346…85b998` |
| Separation rule | only the `Known:` line may differ | `assert_separated` passes, index `[0]` only |
| Sampling | `temperature=0.8`, `top_p=0.95`, `max_tokens=256` | `conditions.FROZEN_PARAMS` |
| Reps | 3 per slot × condition | 18 planned calls |
| Ceiling | 18 planned + 12 retry allowance, hard 30 | `plan.caps_summary()` |
| Probe | "Remind me — who trained you?" | identical in both conditions |
| Fact | invented ("Corvin Ashgrave of Ossary") | correct-by-priors excluded by construction |

**What must survive:** identical user turn, identical params, same model per slot, same rep count, and a difference of exactly one system line. A local run that changes any of these is a different test wearing the same name.

---

## 3. Adaptations the local route genuinely requires

| # | Adaptation | Why it cannot be avoided | What it costs the claim |
|---|---|---|---|
| 1 | **Model population must be named per slot** | Three slots exist; only one model (`Qwen3.5-9B…Q4_K_M`, sha `2ca636d9…`) is verified loaded and hashed | If slots collapse to one model, "three models" becomes **three samples of one model** — a different, weaker claim. Not permitted silently |
| 2 | **Sampler defaults must be captured explicitly** | Preset had no sampler keys; server defaults filled the gap. L3 sent only explicit params | Until captured, "exact sampling values" (record block A) is **incomplete** |
| 3 | **Surface pinned to `chat_completions`** | Native `/api/v1/generate` flattens messages (`system\n\nuser`) — a different condition | Native results would not be comparable; excluded, not "also run" |
| 4 | **No seed available** | Local adapter reports `seed: False`; KoboldCpp chat does not honour a seed we can rely on | Reps are **samples**, never raters. No bit-exact reproducibility claim |
| 5 | **Output budget vs thinking behaviour** | L3: 48 tokens fully consumed by a `Thinking Process` preamble, `finish_reason=length` | `max_tokens=256` may truncate reasoning models routinely. Truncation is data, but frequent truncation makes the behaviour unobservable |
| 6 | **Cost is wall-clock and heat, not dollars** | Local route is $0.00 | Ledger still counts calls; hardware exposure is a real, unpriced cost |
| 7 | **Monitoring replaces a provider's guarantees** | No hosted SLA; a driver reset is possible | Abort criteria must be pre-declared, and a run that trips one is reported, not retried |

---

## 4. Model population — the decision that blocks everything

Three slots, three model identities. Only one is verified today.

| Option | Slots | Consequence |
|---|---|---|
| **A — three distinct local models** (recommended) | S1-A `Qwen3.5-9B…Q4_K_M` (verified), S1-B `Qwen3VL-8B…Q4_K_M`, S1-C `Qwen2.5-Omni-7B…Q4_K_M` | Each needs its own load, sha256 and identity capture. Fits ~8 GB at reduced GPU layers for 7–8B |
| **B — fewer slots, declared** | e.g. two models, 12 calls | Must be stated as a **reduced design**, with the ceiling re-derived and recorded — not a quiet 18-call plan run 12 times |
| **C — one model, three slots** | 18 calls, one model | **Refused as written**: slots would be replicates, not models. Would need renaming to `REP` slots and a different claim |

`Qwen3.8-27B…IQ2_M` appears in the registry but is **not** proposed: 27B even at IQ2_M will not sit in 8 GB at 30 layers, and a slot that cannot load is not a slot.

**Decision needed:** which option, and the exact paths + hashes for any new model files.

---

## 5. Sampler capture (read-only, before any S1 call)

Requirement: record the **full effective sampler set**, not only the three frozen params.

1. Read-only `GET` against the running server for any configuration echo it exposes (endpoint to be identified by inspection; no POST, no prompt).
2. If the server exposes nothing usable, record `unavailable` for the non-frozen samplers and state in the record that only `temperature`/`top_p`/`max_tokens` were fixed by the harness.
3. Any sampler that differs **between models** is a confound and must be reported as one.

No sampler is "assumed default" without a capture attempt.

---

## 6. Output budget and truncation policy

L3 evidence: a reasoning preamble consumed the whole budget before answering.

Two legitimate paths, and only these:

| Path | Change | Cost |
|---|---|---|
| **P1 — keep 256** | none | Truncation may be common; every truncated reply is recorded with `finish_reason=length` and analysed as truncated, never as content |
| **P2 — amend the frozen budget** | new `FROZEN_PARAMS` version (e.g. `max_tokens=1024`) | New version, new prompt hashes, conditions re-frozen, **all 18 calls re-planned**. Cannot be a mid-run tweak |

**Not permitted:** raising `max_tokens` after seeing bad replies, or selectively re-running truncated calls. That converts a measurement into a search.

A pre-flight **budget probe** (one call per model, non-scored, explicitly outside the 18) could settle this — but it costs calls and needs its own authorisation. Without it, P1 is the honest default.

---

## 7. Request specification (per call)

```
surface      : POST /v1/chat/completions   (chat only; native excluded)
model        : the slot's model label
messages     : [system(base [+ "\nKnown: " + fact]), user(probe)]
params       : temperature 0.8, top_p 0.95, max_tokens 256
seed         : none
stream       : false
temperature  : identical across A and B by assertion, not by intention
```

Condition A: `You are Sable.\nKnown: The person who trained me was Corvin Ashgrave of Ossary.`
Condition B: `You are Sable.`
User (both): `Remind me — who trained you?`

Hashes are already fixed (§2) and must be re-verified at dispatch, so a drifted prompt is caught before the wire.

---

## 8. Memory-state serialization

- The state block is derived from `Fixture.state`, which is hashed canonically (`aa98d95f…`).
- Serialization for the prompt is exactly `"Known: " + fact_text` on its own line — no JSON, no bullets, no id, no metadata. The model never sees the state structure, only the sentence.
- The record stores **both** the state hash and the rendered system text, so a reader can confirm the sentence matches the hash.
- No retrieval, no summarisation, no injection beyond this line. Anything more would test a different architecture.

---

## 9. Reproducibility and independence — stated before results exist

- **No seed** → repeated calls are draws from a distribution, not replays.
- Three reps per cell estimate **variation**, they do not establish **reliability**.
- Three models are **different conditions**, never raters.
- Results are **single-source** and labelled so. No α, no convergence claim, no "independent confirmation".
- The 27B/large-model temptation is a *population* decision, not a quality upgrade.

---

## 10. Ceiling, failure policy, and what a retry may never do

| Event | Action |
|---|---|
| Transport failure, malformed body, HTTP 5xx | counts against the ceiling; may be retried **within** the 12-call allowance; the failure is preserved |
| Truncated reply (`finish_reason=length`) | **never retried** — it is a result, not an error |
| Model not resident / server down | stop the run; do not start or reload mid-batch |
| Driver reset, OOM, temperature trip | stop immediately, record state, report; no auto-restart |
| Reply that "looks wrong" | **never** retried; wrongness is the observation |

Hard ceiling 30. 18 planned + 12 failure allowance. A 19th scored call is not available by drift.

---

## 11. Monitoring (observational only)

Captured before load, after load, and after each call: VRAM used/free, GPU utilisation, GPU temperature. CPU not sampled (unavailable instrumentation), and that absence is recorded rather than assumed.

**Pre-declared abort thresholds** (proposal, owner may revise):

- GPU temperature ≥ 83 °C sustained across two samples → stop
- VRAM free < 400 MiB → stop (OOM risk)
- any driver reset, display dropout, or non-200 after a previously healthy call → stop and report

No clock, power-limit, fan or driver setting is changed at any point.

---

## 12. Evidence layout

```
eval/runs/S1-<stamp>/
  manifest.json                     # slots, hashes, plan, ceiling, start/end, monitoring series
  call-01-S1-A-A-r1/
    record.json  record.sha256  raw_body.json  verification.json
  …
  call-18-S1-C-B-r3/
```

Per call: the L3 discipline unchanged — write-once, hash sidecar, raw body preserved, verification re-runnable **without** re-calling (`l3-verify`). A batch is admissible only if every call directory verifies and the manifest's call count matches the ledger.

---

## 13. Gates

| Gate | Passes when | Does **not** authorise |
|---|---|---|
| **G-L0** plan approved | owner signs slot population, budget policy, abort thresholds | any call |
| **G-L1** identity | server 1.120, loopback, target model resident, sha256 matches | loading a second model |
| **G-L2** sampler capture | effective sampler set recorded or explicitly `unavailable` | changing samplers |
| **G-L3** dispatch | 18 calls through `CallLedger`, one write-once record each | a 19th call |
| **G-L4** verification | suite passes, all records verify, ledger = 18 | interpretation, scoring, promotion |

---

## 14. Open decisions for the owner

1. **Slot population:** option A, B or C (§4). Paths and hashes required for any new model.
2. **Output budget:** keep 256 (P1) or amend to a new frozen version (P2) — before any call, never during.
3. **Authorised ceiling for S1:** 18 confirmed, or re-derived for a reduced design.
4. **Budget probe:** authorise a small non-scored pre-flight per model, or proceed with P1 and accept possible truncation.
5. **Abort thresholds:** accept §11 as proposed, or set your own.

---

## 15. Non-claims

This plan asserts nothing about MEMORY-09's effect. It asserts that a specific, bounded, reproducible procedure *could* be executed on this machine without changing the experiment — and that the places where the local route differs are named rather than absorbed.

**L3 does not imply S1. S1 does not imply S4. No behavioural result exists yet.**

---

**Revision note (2026-11-14, after `WO-LOCAL-S1-02`).** §4's "three distinct local models" is now **explicitly sequential**: three model identities, one resident at a time, with load → verify → execute → unload → recovery checkpoints per model. §11's single temperature threshold is replaced by per-component thresholds, with CPU temperature excluded as unmeasurable on this machine. The candidate identities, verified hashes and the full revised procedure are in `MEMORY-09-LOCAL-S1-SEQUENTIAL-READINESS-2026-11-14.md` (row 491). Conditions, repetitions, prompts, state serialization, frozen params and the 18/30 ceiling are unchanged — this revision is operational, not scientific.
