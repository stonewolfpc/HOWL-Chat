# `MEMORY-09` — Python-First Implementation and Validation Plan

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-PYFP-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **PLAN — nothing implemented, nothing executed.**
**Authority:** owner direction of 2026-11-13, following `WO-ARCH-06`.
**Inherits:** rows 470–475 (design), 476 (source audit), 477 (adapter), 478 (state).

---

## 1. Four stages, four gates

| Stage | Question | Environment | Passes when |
|---|---|---|---|
| **S1 · Concept** | Does the memory concept work with real models? | Python + hosted inference | Canonical facts are honoured in the assisted condition and not in the unassisted one, on at least two model families |
| **S2 · Architecture** | Does the system behave reliably across controlled tests? | Python + hosted, frozen configs | The matrix (T-01…T-11) runs end to end with recorded results; no unrecoverable defects in state handling |
| **S3 · Local inference** | Does it work with real local models and this hardware? | Python + `llama.cpp`, 40B-class included | The same tests run locally; behaviour is compared, not assumed equivalent |
| **S4 · Howl integration** | Does the implemented application deliver the behaviour? | Howl end-to-end | Owner judges the user-facing experience directly |

**Passing one stage passes no other.** S1 success says nothing about S3's hardware behaviour; S3 success says nothing about Howl's wiring. **A later stage may fail for reasons earlier stages never tested** — that is expected, not a regression.

## 2. What is inherited from `WO-ARCH-06` (unchanged)

| Kept | Why |
|---|---|
| Four layers and their claim ceilings (row 470 §2) | The authority structure is independent of language |
| A–D record blocks | Nothing about Python changes what must be recorded |
| `[MECH] [OBS] [AI] [OWNER]` tags | Same discipline |
| Presentation contract (one sentence, three observations, then evidence) | Owner-facing, unaffected by harness |
| Repetition doctrine (four reasons only) | Same evidential logic |
| `C2` vs `C3` mandatory classification | The crux, at every stage |
| T-01…T-11 test designs | Written as behaviour, not implementation |
| Confound table | Same hazards, one new one added (§5) |
| Owner adjudication protocol (row 474) | Unchanged |

## 3. What changes

| Was (ARCH-06) | Now |
|---|---|
| Howl runtime implied as subject | **Python harness; Howl deferred to S4** |
| Local models and hardware central | **Hosted inference first; local at S3** |
| Ten-slot population across all categories | **Small hosted population, creative-first** |
| Sampler defaults read from Howl's Go table | **Reproduced in harness config** (row 472 values, not Go code) |
| Measuring load time / VRAM | Deferred to S3 |

## 4. Hosted population (S1–S2) — small and creative-first

**Three to four models, not ten.** Availability, provider support and pricing are **confirmed at selection time**, never assumed here.

| Slot | Selection intent | Why |
|---|---|---|
| **H-1** | Creative-writing-oriented model, strongest available | Primary category; the product's surface |
| **H-2** | Creative-writing-oriented, smaller/cheaper | Does architecture compensation shrink as prose quality rises? |
| **H-3** | General instruct, different family from H-1 | Family effect, not quality effect |
| **H-4** *(optional)* | Reasoning-flavoured | Does deliberation help or hurt in-character persistence? |

**Rules:** one provider **pinned per model** (adapter HF-1); the same model never compared across two providers in one series; any model whose provider changes mid-series restarts the series. The local 40B-class test (row 471, M-08) **remains S3** and is not part of S1–S2.

## 5. New confound, and it is significant

| Confound | Why it is new | Handling |
|---|---|---|
| **Provider/quantisation drift** | Hosted routing may serve the same model id with different quantisation, context or sampler behaviour | Pin provider; record `provider_used` and `model_reported`; treat a provider change as a **new condition**, not a rerun |

This is the hosted analogue of "model re-quantised mid-series" — and it can happen invisibly, which is why the adapter records what actually served the call.

## 6. Assisted vs unassisted

Both conditions: **identical user messages, identical sampler values, identical max tokens**, same model, same provider. The assisted condition receives state-derived context; the unassisted condition does not.

**Disclosure requirement:** if the unassisted prompt differs structurally (system preamble, turn count), the record states the difference explicitly. **No claim of equivalence is made where equivalence was not achieved** — the same rule as row 470 §4.

## 7. Provenance per run

Adapter name/version, provider, model requested and reported, exact params, seed and whether honoured, prompt hash, response hash, state hashes before/after, fixture id/version, harness version, test id. **A run missing any of these is not admissible.**

## 8. Bounds

Hosted inference costs money and rate limits, so: cap runs per session; no unbounded retries (adapter max 3); a failed call is recorded, not re-billed; and **no exploratory sweep** — each run belongs to a named test.

## 9. What S1 must not be read as establishing

- That the architecture is validated (that is S2).
- That behaviour transfers to local models (S3) or to Howl (S4).
- That any `MEMORY-09` cue criterion is satisfied — the corpus stream stays separate (row 476).
- That a correct reply was caused by state (`C3` must be excluded case by case).

## 10. Decisions required

| # | Decision | Recommended | If declined |
|---|---|---|---|
| 1 | Adopt the `OD-M09-17` audit finding (row 476) **and separately authorise engineering evaluation** | Yes — the record grants no permission | Engineering stays unauthorised |
| 2 | Approve the adapter interface (row 477), including provider pinning | Yes | Hosted runs are not reproducible |
| 3 | Approve the state interface and fixtures (row 478) | Yes | `C2`/`C3` cannot be checked |
| 4 | Approve the hosted population (3–4 models) | Yes | Model choice becomes ad hoc |
| 5 | Authorise S1 execution, with a run/cost cap | Separate signature | Nothing runs |

## 11. Unauthorised

Model calls; any S1–S4 execution; corpus, tag, split or score changes; standards amendments; α thresholds or policy; `G2`, `P-3`, `L4`; Howl source changes; and any statement that a stage has passed.

## 12. Claim ceiling

A plan. No adapter, no fixture, no call, no model. It establishes what would count as passing each stage, not that any stage is close.
