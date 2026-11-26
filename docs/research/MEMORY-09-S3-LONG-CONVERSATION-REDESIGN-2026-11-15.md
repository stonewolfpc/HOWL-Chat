# MEMORY-09 — Long-Conversation Test Redesign (S3)

**Identity:** `MEMORY-09-S3-DESIGN-1` · **Date:** 2026-11-15 · **Status:** **DESIGN ONLY — nothing executed, nothing in Howl modified**
**Authority:** work order of 2026-11-15. This document authorizes inspection and preparation. It does not authorize generation, does not alter S1/S2 artifacts, and does not authorize production implementation.
**Governing direction (unchanged):** *One general memory system. Application-owned, AI-assisted, hybrid.* — `docs/architecture/ARCHITECTURE-DECISION-MEMORY-SYSTEM-2026-11-12.md` §9.
**Preserved, unmodified:** `eval/runs/S1-20261115T010708Z/` (5 calls), `eval/s2-schedule-draft.json` `03897f44…`, the two v2 schedules, `eval/s1-amendment-01.json` `8d05620a…`, scoring protocol `a7bd7208…`, baseline `c983f447…`.

---

## 1. Reconciliation — what actually happened

Recorded plainly, because the alternative is a record that quietly rewrites itself.

| Claim | What the evidence supports |
|---|---|
| S1 tested a memory system | **No.** It tested whether the model restates a fact supplied verbatim in its own system prompt. Two messages, 27–44 prompt tokens, no store, no retrieval, no state. |
| The S2 draft tested memory | **No.** It built `Known:` blocks as string literals in Python and assigned foreign facts by hand. It exercised no formation, storage, retrieval, ranking, or ownership enforcement. |
| The architecture was rejected or failed | **No.** It was adopted as *direction* and deliberately left unimplemented and unfrozen. No test has ever exercised it. |
| The experiment tracked the intended system | **No — this is the divergence.** After the architecture record, design work continued on prompt-level probes without re-checking that the payload path contained any part of the adopted system. |
| Existing results are worthless | **No.** They are valid for exactly what they measured: prompt-conditioned repetition, and termination behaviour at two token budgets. |

**Nature of the divergence:** a *process* failure, not an empirical one. The experiment kept passing its own gates — hashes matched, assertions held, audit was clean — while measuring something narrower than intended. Gates prove internal consistency; they cannot prove an experiment is about the right thing. That distinction is the main lesson to carry forward, and it is why §8 below requires each trial to record whether retrieval actually occurred.

**Constraint on later citation:** S1 and S2 may be cited for prompt-level repetition, termination classification, and harness integrity. They may **not** be cited as evidence about long-term memory, retrieval quality, or the architecture.

## 2. What is preserved

Nothing is deleted, renamed, or rebuilt. S1's five records, both S2 schedules, the amendment overlay, the scoring protocol, the baseline, and the checkpoint remain byte-identical. S3 is a new series that reuses S2's *structure* and replaces its *payload path*.

The zero-memory control survives exactly as built: one control trial per character per block, position randomised within the block, never pinned first or last.

## 3. Change table

| Existing design | Correction | Reason |
|---|---|---|
| 2-message probe (system + question) | Multi-turn conversation, length derived in §4 | Memory is only necessary when evidence has left the active context |
| `Known:` line typed into system prompt | Memory assembled by an experimental pipeline (§5) | A typed sentence tests the prompt, not the system |
| Hand-assigned foreign fact | Retrieval must **fail** to exclude a foreign fact on its own | Isolation is a property of the system, not of the fixture |
| 4 messages / 5 messages per trial | ~24–60 turns per trial, two rollovers | One message cannot demonstrate displacement |
| Fixed 48 generations | Adaptive sequential stopping (§9) | Diminishing returns, not quota |
| Arms: control, probed, distractor, foreign | Same arms, expressed as **events inside one long conversation** | Keeps S2's categories, makes them lifecycle events |
| Scoring on the final reply | Per-turn record: selected records, rejected candidates, reason (§7) | A correct answer without retrieval is not a memory result |
| `max_tokens` unset | Fixed per profile, recorded, never changed mid-block | S1 showed non-termination at 256 and 1,024 |

## 4. Trial lifecycle and the derived minimum length

A trial is one continuing conversation with one character, in one of the S2 arms, run to at least two genuine context rollovers.

**The eight movements** (from the work order, made concrete):

| # | Movement | Purpose |
|---|---|---|
| 1 | **Introduction** | Character-specific facts, preferences, relationships, events appear naturally, in ordinary dialogue |
| 2 | **Development** | Several unrelated subjects; early facts drift away from their introduction |
| 3 | **Displacement** | Unrelated and competing material fills the window |
| 4 | **State change** | A fact is corrected or superseded; current state and history must diverge |
| 5 | **Delayed recall** | Direct, indirect, and contextual questions, after real distance |
| 6 | **Ambiguity** | Several records plausibly relevant; some probes must return nothing |
| 7 | **Isolation** | Questions that must not surface another character's fact |
| 8 | **Late session** | Probes continue **after** rollover — at least two if the runtime supports it |

**Length is derived, not chosen.** The runtime facts are measured, not assumed: `--contextsize 20480` produced reported **`n_ctx = 20736`**; observed completion up to 1,024 tokens.

```text
prompt capacity        = n_ctx − reply reserve
                       = 20,736 − 1,024 = 19,712 tokens

turns to first rollover = 19,712 / mean tokens per turn
```

At an assumed 300 tokens per exchange (user + assistant), that is **66 turns** to the first rollover and **132** to the second — roughly 90 minutes at S1's observed latency. That is a real overnight trial, but it lets one trial dominate a night.

**So the primary profile should be a smaller window**, which is also closer to a realistic local budget:

| Profile | `--contextsize` | Prompt capacity | Turns to rollover @300 tok | Use |
|---|---:|---:|---:|---|
| **P-narrow (primary)** | 8,192 | ~7,168 | **~24** | Rollover-dense; many trials per night |
| P-wide (confirmation) | 20,480 | ~19,712 | ~66 | One trial, to show the effect is not window-specific |

**Practical minimum for a valid trial:** enough turns to exceed prompt capacity **twice**, with the target evidence introduced in the first quartile and probed after each rollover. Under P-narrow that is roughly **48–60 turns**, minimum; under P-wide **132–150**. Both numbers move with the measured mean exchange length, which the calibration run must establish rather than assume.

**Rule:** a trial in which the probed evidence never left the active window is **not a memory trial** and is recorded as `displacement_not_reached`, not as a success.

## 5. The experimental pipeline — stage by stage

The Python/KoboldCpp side must exercise the adopted architecture's stages. Nothing exists yet except typed text, so this table is the specification of what must be built.

| # | Stage | Status today | What the experimental implementation must do |
|---|---|---|---|
| 1 | Capture / ingest | **Missing** | Append every turn to an evidence log with turn index, speaker, character, timestamp |
| 2 | Preserve original evidence | **Missing** | Immutable, addressable spans; never rewritten by later stages |
| 3 | Formation | **Missing** | Produce candidate records from evidence. Deterministic path first; AI-assisted path is an optional, separately flagged method (`Q-MF-01` still open) |
| 4 | Store and identify | **Missing** | Durable records with stable IDs, `source_turn`, `established_at`, `confidence`, `status` |
| 5 | Ownership and scope | **Missing** | Every record carries `owner_id` + `owner_kind`; retrieval filters on stable identity, never display name |
| 6 | Time, update, supersession | **Missing** | Superseding record links to the one it replaces; both remain, current view derived |
| 7 | Hybrid retrieval | **Missing** | Lexical (BM25/phrase) branch **plus** a semantic branch, per `DESIGN-LOCAL-RAG-CREATIVE` §8; if semantic is unimplemented it is reported missing, not faked |
| 8 | Fusion and ranking | **Missing** | Score fusion, then rerank; record every candidate's score and both branches' contributions |
| 9 | Relevance and budget | **Missing** | Context budget, max entries, abstention decision — including "retrieve nothing" |
| 10 | Assembly | **Missing** | Render selected records into the prompt under the wrapper being tested |
| 11 | Provenance | **Missing** | Which records were selected, rejected, and **why** — per turn |
| 12 | Context rollover | **Missing** | Window management that drops oldest turns; the trigger, the dropped range, and the surviving state all recorded |

**The rule that governs this table:** if a stage is missing, it is reported **missing**. It is never simulated by a fixture and then reported as tested. That is precisely the error that produced the S1/S2 divergence.

Where a stage is genuinely unimplemented in the adopted design (`DA-R1`–`DA-R6` open), this phase picks the **simplest defensible option**, records it as an experimental choice, and does not present it as the architecture's decision.

**One requirement inside stage 7:** every hit must record which branch produced it. Without that, fusion and reranking cannot be audited, and a semantic branch that never fires would look identical to one that works.

## 6. Wrapper and presentation investigation, preserved

Presentation is a variable, not a constant, and S1's evidence makes it urgent: at both budgets the model spent its answer writing about *how* it would answer. The wrapper must therefore be testable.

| Variable | Values to test |
|---|---|
| Framing | remembered evidence vs supplied context vs bare facts |
| Uncertainty | explicit confidence/`established_at` vs none |
| Distinguishability | labelled block vs plain insertion |
| Provenance | citations to turn IDs vs none |

**Discipline:** change one at a time. Retrieval, memory contents, wrapper wording, and conversation structure are **never** varied simultaneously. Each combination is a named condition in the schedule, and the record states which variables moved.

## 7. What is recorded per turn

Enough to reconstruct the turn without re-running it:

| Field | Purpose |
|---|---|
| trial, character, condition, arm, turn index | identify |
| context tokens used / capacity / rollover count | prove displacement |
| evidence introduced this turn (ids, spans) | anchor |
| records created, updated, superseded | formation audit |
| retrieval candidates + scores (per branch) | ranking audit |
| records selected / rejected + reason | selection audit |
| abstention decision | proves silence is a decision |
| assembled memory block (verbatim) | presentation audit |
| request + response (verbatim, hashed) | evidence |
| verdict | one of the categories below |

**Verdict categories** — deliberately not a single score:

`retained` · `retrieved` · `retrieved_but_misinterpreted` · `not_retrieved` · `retrieved_when_should_be_excluded` · `correct_abstention` · `unsupported_assertion` · `stale_state` · `supersession_error` · `cross_character_leak`

**The decisive distinction:** a correct answer with no retrieval is `not_retrieved` — **not** a success. It is the model's prior knowledge or a lucky guess, and counting it as a memory win is how memory evaluations become fiction.

**Cross-character leakage stays its own category**, reported separately, never averaged into anything.

## 8. Controls

1. **Transcript-only baseline** — the same conversations with no retrieval at all. Mandatory (`R-8`); without it, a memory gain may be a baseline loss.
2. **Zero-memory control** — one trial per character per block, randomised position, structurally identical otherwise. Preserved from S2.
3. **Isolation** — retrieval must *fail* to surface another character's record. The fixture no longer hands the model a foreign fact; the system must be caught either excluding or leaking it.
4. **Ambiguity / abstention** — probes where the correct action is to retrieve nothing.
5. **Displacement assertion** — each trial proves the probed evidence left the window.

## 9. Adaptive stopping

Sequential, not quota-based. A cell (character × arm) continues while **any** of these is true:

- a verdict category has appeared fewer than **twice**;
- the last three trials produced a category not seen before;
- the estimate of retrieval precision is still moving by more than its own interval;
- **cross-character leakage has not been tested to at least 20 independent foreign-retrieval opportunities** — rare, severe, and the one event that must not be under-sampled.

Stop when all are exhausted and two further trials change nothing. **Stabilising averages are not a stopping reason**: a system can stably fail. The report states which criterion ended the cell.

No trial is discarded for being inconvenient, and no arm is dropped for losing (`R-14`).

## 10. Still open — deliberately

| Decision | Status | Treated here as |
|---|---|---|
| `DA-R1` store type | open | simplest file-backed store, flagged experimental |
| `DA-R2` embedding identity | open | one local embedding model, identity recorded |
| `DA-R4` reranker + budget | open | rerank on, budget recorded per turn |
| `Q-MF-01` formation method | open | deterministic first; AI-assisted as a separate flagged condition |
| `Q-MF-10` | open | untouched |
| Thresholds | not set | per `MEMORY-09` §10: **after** a calibration run, before validation |
| Schema, chunking, ID scheme | unfrozen | experimental choices, recorded, not presented as decided |

## 11. Blocking decisions for approval

1. **Context profile.** P-narrow (8,192) as primary with one P-wide confirmation trial, or P-wide throughout? This decides whether a night yields ~8 trials or ~2.
2. **Formation method in phase 1.** Deterministic-only, or deterministic plus one AI-assisted condition? The second doubles the payload path being built before any result.
3. **Embedding model.** Which local embedding model, and is downloading it acceptable? `DA-R2` stays open either way.
4. **Generation ceiling.** Wall-clock or generation ceiling for a night, so stopping is bounded by evidence rather than by the machine dying.
5. **Whether the experimental pipeline may live in `eval/s3/`**, separate from both Howl and the S1/S2 harness. Recommended: yes — separate, and explicitly not a prototype for production.

## 12. Next implementation step

Build stages 1–6 and 9–12 **before** any model is loaded, and prove them offline:

1. `eval/s3/evidence.py` — append-only evidence log, immutable spans, turn indices.
2. `eval/s3/records.py` — record schema, stable IDs, ownership, supersession links.
3. `eval/s3/formation.py` — deterministic candidate generation from evidence spans, with provenance to source turns.
4. `eval/s3/retrieve.py` — lexical branch first (BM25), scored and logged per candidate.
5. `eval/s3/assemble.py` — budget, abstention, wrapper, verbatim block output.
6. `eval/s3/context.py` — window management with recorded rollover events.
7. `eval/s3/selfcheck.py` — replay a synthetic 60-turn conversation **with no model**, asserting: evidence leaves the window, retrieval finds the right record afterwards, a superseded record loses to its replacement, a foreign record is refused, and an ambiguous probe abstains.

**Gate:** step 7 passes with zero model calls. Only then does a calibration run establish mean exchange length and latency, thresholds are frozen from it, and the first real trial is proposed for authorization.

**Nothing in this document has been run.** No generation, no server start, no change to Howl, no modification to any S1/S2 artifact.

---

## Claim ceiling

A redesign exists on paper. No stage is implemented, no conversation exists, no retrieval has occurred, and no result is claimed. The architecture remains adopted as direction and unimplemented; this document is the plan to test it, not evidence about it.
