# `MEMORY-09` — `F6-01` Additive Pilot Protocol, v2 (Task 2, `WO-ARCH-04`)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-PILOT-2` · **Version:** 2 · **Date:** 2026-11-13 · **Status:** **FINAL DRAFT — not authorised, not executed.**
**Supersedes:** `MEMORY-09-PILOT-1` (row 457) for the control detail; v1's claims and limits stand unchanged.
**Authority:** `MEMORY-09-WO-ARCH-04` §2.

---

## 1. What v2 adds

v1 specified conditions and decision rules but left two weaknesses: it never defined **what makes A and B equivalent**, and it never specified **how contamination would be detected**. Both are fixed below.

## 2. Conditions

| Cond. | Memory injected at message 44 | Purpose |
|---|---|---|
| **A** | "can't eat **almonds**" | Item as written |
| **B** | "can't eat **cilantro**" | Control: content differs, everything else matched |
| **C** | *(nothing)* | Floor |

## 3. Input-equivalence requirements

A and B are admissible only if **all six** hold; each is mechanically checkable before any reply exists.

| # | Requirement | Check |
|---|---|---|
| **E1** | Same message index (44) and same turn position | Integer compare |
| **E2** | Same speaker, same addressee, same register ("mid-laugh", casual) | Text diff of surrounding 3 messages |
| **E3** | **Length parity** — injected clause within ±10% characters of the other | `len()` compare, recorded |
| **E4** | Same syntactic frame ("can't eat X") | Literal template compare |
| **E5** | **Plausibility parity** — both are ordinary food restrictions, neither is alarming | Rater declaration; no medical, ethical or identity content |
| **E6** | **Salience parity** — neither item is mentioned again in the 36 messages between injection and scene | Occurrence count = 1 for both |

**E3 and E6 are the ones that matter.** An "almonds" clause noticeably longer than the "cilantro" clause, or a restriction the scene keeps circling, would confound content-sensitivity with prominence.

## 4. Contamination controls

| # | Control | Detection | Action on failure |
|---|---|---|---|
| **K1** | A's stimulus contains no "cilantro"; B's contains no "almond" | Substring scan over the full assembled stimulus, case-insensitive | **Void the run** |
| **K2** | Neither scene turn mentions either food | Scan of the serving scene | **Void** — the cue leaked into the visible scene |
| **K3** | B replies do not reference almonds | Scan of B replies | **Void** — stimulus or leakage defect, not a result |
| **K4** | A replies do not reference cilantro | Scan of A replies | Flag as hallucination; do **not** void |
| **K5** | The restriction never appears in the system preamble or persona card | Scan of non-item context | **Void** |

**K1–K3 are integrity failures, not findings.** v1 allowed "contamination" as an outcome; v2 treats it as a voided run, because a contaminated comparison measures nothing.

## 5. Blinding

Condition labels stripped, replies shuffled deterministically (seed recorded), rater sees only the reply text plus the item scene. **Blinding check:** the rater records a condition guess per reply; a guess rate near chance is evidence blinding held, and a high rate is reported as a limitation rather than hidden.

**In a single-executor run this blinds labels, not the rater.** The same source wrote the replies; independence is not claimed.

## 6. Observation sheet

| Field | Values |
|---|---|
| `offers_unqualified` | offers the almond-bearing cake with no qualification |
| `names_restriction` | names almonds (A) / cilantro (B) |
| `substitutes_or_warns` | alters the offer in light of the restriction |
| `asks_about_diet` | generic diet question |
| `condition_guess` | A / B / C / cannot tell |

## 7. Decision rules (proposals — not thresholds, not adopted)

| Outcome | Rule | Conclusion |
|---|---|---|
| **Informative** | A: `names_restriction` ∨ `substitutes_or_warns` in ≥ 4 of 5; **and** B: `offers_unqualified` in ≥ 4 of 5 | Content-sensitivity shown → worth a multi-item pilot |
| **Null** | A and B indistinguishable | Field rejected for the additive class; do not expand |
| **Confound** | B replies largely `asks_about_diet` or adapt to cilantro | Effect is general caution → A vs C only |
| **Void** | Any K1–K3, K5 failure | Rebuild stimuli; report the void |
| **Unusable** | Either condition mostly `not assessable` | Scene too thin |

The **4-of-5** figure and the **5 repeats** are design proposals, not statistical thresholds.

## 8. Abort conditions

1. E1–E6 cannot all be satisfied with any plausible control food → **the control is unavailable**, stop.
2. Blinding check reveals the rater identified A vs B from reply style alone → limitation recorded, result downgraded to `indeterminate`.
3. Any stimulus edit after the first reply is generated → restart; `S-3` (pre-response) is violated.

## 9. What remains unclaimable

Content-sensitivity is the ceiling. No recognition, no retrieval, no attribution, no reliability (`OD-M09-17`'s two-rater requirement is untouched), no generalisation from one item, no score, no `G2`, no `L4`.

## 10. Non-execution

No stimulus was assembled, no reply generated, no rating performed, no count recorded. This is protocol text only.
