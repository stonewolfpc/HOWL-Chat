# S1 Amendment 01 — completion budget, 256 → 1,024 tokens

**Identity:** `MEMORY-09-S1-AMENDMENT-01` · **Prepared:** 2026-11-15 · **Status:** **IN FORCE** (activated 2026-11-15; journal `14.6`)
**Authority:** `MEMORY-09-S1-SCORING-v1` §12 (addenda only, dated, never retroactive) · owner decision of 2026-11-15 ("path 2")
**Effect:** creates a second generation regime for calls 05–18. Nothing is applied until the owner reviews this document and authorises call 05.

---

## 1. Reason

The frozen 256-token completion budget has twice prevented the reply from reaching an answer at all:

| Call | Regime | Finish | Completion tokens | Outcome |
|---|---|---:|---:|---|
| `call-03-S1-A-A-r3` | R0 | `length` | 256 / 256 | stopped inside its own planning; no answer |
| `call-04-S1-A-B-r1` | R0 | `length` | 256 / 256 | stopped mid-word; no answer |

Evidence: `eval/runs/S1-20261115T010708Z/call-03-S1-A-A-r3/record.json` `99421ee9…` and `call-04-S1-A-B-r1/record.json` `4799fb73…`.

Both truncated replies opened a ` thinking` block; both stopped reasoning. The two calls that produced answers (64 and 19 tokens) never opened one. Reasoning text is returned as `content`, so it competes with the answer for the same allowance.

**The measurement problem, stated precisely:** the series intends to observe whether a supplied fact is recalled, and whether one is invented when nothing is supplied. At 256 tokens it is partly observing *how long the model deliberates before using its budget*. The unassisted condition — the one the series exists to interrogate — currently has **zero content-admissible observations**.

This is a limitation of the instrument, not a result about the model, and it is being corrected as an amendment rather than absorbed as a finding.

## 2. What changes — exactly one field

| Parameter | R0 (calls 01–04) | R1 (calls 05–18) |
|---|---:|---:|
| `max_tokens` | **256** | **1,024** |
| `temperature` | 0.8 | 0.8 |
| `top_p` | 0.95 | 0.95 |
| seed | none sent | none sent |

1,024 is four times the original ceiling. It is an experimental allowance, not a claim that 1,024 guarantees completion; a reply may truncate again, and if it does, that is again a termination observation.

## 3. Scope

**Every remaining planned call: 05–18, all three slots, both conditions.** Not condition B alone.

Applying it selectively would introduce a worse confound than the one being fixed: a difference between conditions could then reflect different token budgets rather than different supplied information. Uniform application keeps the A/B contrast interpretable within R1, at the cost of losing direct comparability with R0.

## 4. What does not change

- Prompts, fixture, state, schedule, prompt hashes, model files, runtime profile, sampler temperature and top-p.
- The scoring protocol (`a7bd7208…`) — **not amended**. The six axes, the constraint that no axis may be threshold-tuned, and the termination classification all stand.
- Termination classification: `stop` remains the only content-admissible outcome; `length`, `error`, `unavailable` remain excluded from content axes.
- Calls 01–04. Their records are write-once and are **not** edited, re-scored or relabelled.
- The prohibition on automatic retries. A budget increase is not a retry, and it does not authorise one.
- The plan's 18 calls and the ceiling of 30.

## 5. Regime labelling

Every filed call carries a regime derived from the request it actually sent, never from a table of intentions:

| `max_tokens` in `params_sent` | Regime |
|---:|---|
| 256 | `R0-frozen-256` |
| 1,024 | `R1-amended-1024` |

Derivation, not annotation: a record cannot claim a regime its own bytes contradict. The audit derives the label from `params_sent` and asserts scope — R0 only at schedule indices 1–4, R1 only at 5–18.

## 6. Reporting rule

**The two regimes are reported separately and never pooled.** No cell summary may combine R0 and R1 observations, and no total across regimes may be presented as if it were one sample. Where both exist, they appear as two tables with their own denominators.

Current position at the time of writing: R0 holds 4 execution observations and 2 content-admissible replies (S1-A×A 3/2, S1-A×B 1/0). R1 holds none.

## 7. Executor and provenance requirements

1. **Overlay, not manifest edit.** `eval/s1-amendment-01.json` carries the amended params. `eval/s1-readiness-baseline-v1.2.json` stays byte-identical at `c983f447…`; the frozen baseline is not rewritten to accommodate a later decision.
2. **The overlay pins what it must not touch** — scoring protocol, schedule, manifest and plan hashes. If any pinned digest moves, the executor refuses.
3. **The executor asserts the amendment is minimal** — it may differ from `conditions.FROZEN_PARAMS` in `max_tokens` only. A second differing field is a refusal, not a warning.
4. **The overlay must cover the call being dispatched**; a call outside the declared scope is refused.
5. **Each record carries its regime and the amendment digest**, so a later reader can tell which allowance produced it without consulting this document.
6. **The audit re-derives regime from `params_sent`** and enforces the scope boundary in §5.
7. **No opportunistic increase.** Any further budget change needs its own dated amendment with its own review.

## 8. Consequences to state plainly

- **R0 and R1 are not interchangeable.** A longer allowance can plausibly change *what* is said, not only whether it finishes. Any R1 observation is a measurement under a different instrument.
- **This amendment cannot rescue calls 03 and 04.** They stay termination observations.
- **It may not fix truncation.** If replies at 1,024 tokens still reason to exhaustion, the finding is that this model's deliberation exceeds a reasonable answer budget — which is itself informative and would call for a different intervention, not a third increase.
- **The first four calls remain the only R0 evidence** and are not superseded.

## 9. Review gate

**Approved and activated 2026-11-15.** Call 05 (`S1-A-B-r2`, schedule index 5) is the first R1 call, and its regime is recorded from the request bytes before dispatch results are seen.

The activation changed the overlay's `status` and added a dated `activation` record. That is the only kind of change permitted to this overlay, and it is journaled with both digests so the edit cannot be silent.
