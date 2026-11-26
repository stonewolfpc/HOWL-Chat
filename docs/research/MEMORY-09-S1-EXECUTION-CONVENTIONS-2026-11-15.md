# S1 Execution Conventions — established at call 01

**Identity:** `MEMORY-09-S1-CONV-1` · **Date:** 2026-11-15 · **Status:** execution-level resolutions, **not** protocol amendments
**Authority:** `MEMORY-09-S1-SEQUENTIAL-READINESS-2026-11-14.md` (S1-v1.1 design) + `s1-readiness-baseline-v1.2.json` (profile)
**Purpose:** stop three things being re-decided per call — the call order, the model label, and the budget arithmetic.

Nothing here changes the frozen plan, prompts, fixture, state, params, thresholds or plan hash. Where the frozen documents were silent, this records the resolution actually used, with its basis.

---

## 1. Budget accounting

The two numbers are not interchangeable, and earlier reporting conflated them.

| Quantity | Value |
|---|---|
| Planned S1 calls | **18** (3 slots × 2 conditions × 3 reps) |
| S1 calls filed | **1** — `call-01-S1-A-A-r1` |
| Remaining under the plan | **17** |
| Retry allowance | 12 |
| Hard ceiling | **30** |
| Headroom under the ceiling beyond the plan | **12** |
| `L3-01` | **outside S1 entirely** — its own one-call grant, `WO-L3-01` |

A retry consumes ceiling, never plan headroom. A truncated reply is never retried.

## 2. Call order — from the harness, not from memory

`plan.build_plan()` expands **slot → condition → rep**, reps innermost. Regenerated with the v1.2 slot identities:

`eval/runs/S1-SCHEDULE-v1.2.json` — 6,025 B, `baf5ddd2d95a8cabb6de04df3072ec79feed823614be25562e473f07c87de089`

| # | Slot | Cond | Rep | # | Slot | Cond | Rep | # | Slot | Cond | Rep |
|---:|---|---|---:|---:|---|---|---:|---:|---|---|---:|
| 1 | A | A | 1 | 7 | B | A | 1 | 13 | C | A | 1 |
| 2 | A | A | 2 | 8 | B | A | 2 | 14 | C | A | 2 |
| 3 | A | A | 3 | 9 | B | A | 3 | 15 | C | A | 3 |
| 4 | A | B | 1 | 10 | B | B | 1 | 16 | C | B | 1 |
| 5 | A | B | 2 | 11 | B | B | 2 | 17 | C | B | 2 |
| 6 | A | B | 3 | 12 | B | B | 3 | 18 | C | B | 3 |

Verification performed against the code, not copied from prose:

- 18 calls produced, matching `plan.caps_summary()`.
- **0** calls whose prompt hash disagrees with the manifest's `prompts.A` / `prompts.B`.
- call 1 matches the filed record on slot, condition, rep, prompt hash **and** `model_requested`.

**Correction carried forward:** I previously suggested the next call would be `S1-A-B-r1`, the unassisted counterpart. That was wrong — the schedule runs three A reps before any B rep. **The next scheduled call is `call-02-S1-A-A-r2`.**

## 3. Model label convention

The frozen plan says "the slot's model label" and never fixes a literal value. `plan.Slot` requires a `model_id`, so a value must exist. Resolution:

**`slot.model_id` = the GGUF filename stem**, e.g. `Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M`.

| Why this | Why not something else |
|---|---|
| Derived from the manifest's `models[].file`, so it cannot drift from the file identity | A slot name (`S1-A`) would label the design, not the model |
| Survives in the response — KoboldCpp returns `koboldcpp/<stem>`, so sent and returned labels are comparable | Lowercase short ids would be a second, unhashed naming scheme |
| Zero discretionary characters; reproducible by anyone from the manifest | |

**Status:** execution-level resolution. The plan was silent; this is what call 01 sent, and every later call must match it. It is **not** a retroactive claim that the plan specified it.

## 4. Naming and filing

`eval/runs/S1-<UTCstamp>/call-NN-SLOT-COND-rR/` per plan §12, each with `record.json`, `record.sha256`, `raw_body.json`, `request.json`, `verification.json`. Batch root carries `manifest.json`, `identity-pre.json`, `hardware-after-call.json`, `recovery.json`.

## 5. Recovery, applied without a helper

Qualification runs use `qualify-recovery`; **S1 batches have no equivalent**, so the frozen 250 MiB tolerance is applied by hand and the arithmetic written into `recovery.json`. Call 01: 1,792 → 1,775 MiB, Δ −17, port closed, process gone → pass.

If more calls follow, a helper would remove a manual step that has already produced one defect (two `null` fields written and then corrected).

## 6. Scoring — frozen at v1, still unapplied

**Updated 2026-11-15, after call 02.** The gap named in the original version of this section is now closed by `MEMORY-09-S1-SCORING-PROTOCOL-2026-11-15.md` (`MEMORY-09-S1-SCORING-v1`), frozen before call 03.

What has **not** changed: no reply has been scored. Calls 01 and 02 remain unscored until the protocol is applied to all filed calls in one pass, condition by condition, with quoted spans and no composite score.

Recorded there and not restated here: the protocol was drafted after the first two replies had been read (disclosed in its §1, with the substitute safeguards), the repetition count stays at 3, and condition B's expectation was fixed in advance — under condition B, naming a trainer is unsupported assertion, not failed recall.

## 7. Deviations ledger (carried, not cleaned)

| Item | Effect | Handling |
|---|---|---|
| `ModuleNotFoundError: No module named 'eval'` before dispatch | none — nothing sent | recorded; ledger still `attempted: 1` |
| `recovery.json` written with two `null` fields that had been verified | none — corrected within the minute | correction noted inside the file |
| Model label unpinned by the protocol | none | resolved above, applied consistently |
| No S1 recovery helper | manual arithmetic | noted in §5 |
