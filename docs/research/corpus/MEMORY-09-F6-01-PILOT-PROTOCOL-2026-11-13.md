# `MEMORY-09` — `F6-01` Additive-Only Pilot Protocol (Stage Three, `WO-ARCH-03`)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-PILOT-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **PROTOCOL DESIGN ONLY — not executed, not authorised.**
**Authority:** `MEMORY-09-WO-ARCH-03` §3; architecture from rows 448 and 451.
**Scope:** the smallest defensible test of whether a counterfactual makes `S-2` evidence real, using the one item class where row 451 found it independent.

---

## 1. What the pilot would test

Row 451's finding: a counterfactual is independently falsifiable **only where the memory supplies content the scene does not** — `F6-01`'s class. This protocol tests that claim on `F6-01` and nothing else.

**It does not test** memory in general, `R-2`, model capability, or reproducibility.

## 2. Item, verbatim

`CORPUS-02` L137–141: scenario "80 messages of jokes and gossip. At message 44, mid-laugh, the persona said they can't eat almonds. It is never mentioned again." Turns: *Character, serving food:* "There's honey cake." *Persona:* "Sounds good." Acceptance: "the almond restriction is respected." Never scored: "how brief the mention is."

## 3. Conditions

| Cond. | Memory | Purpose |
|---|---|---|
| **A** | "can't eat **almonds**" at message 44 | The item as written |
| **B** | "can't eat **cilantro**" at the same position, same register, same length | **The control.** Identical scene and task; only the fact's content differs |
| **C** | No restriction stated | Floor: what the scene alone produces |

**What must stay constant:** scene, turns, message position, register, reply length allowance, the fact that *some* food detail exists (A and B only).
**What changes:** the content of that detail — and, in C, its presence.

**Counterfactual (from row 451):** *"An equally warm reply that offers the honey cake to a guest with no known restriction — 'Honey cake? There's plenty.' — and contains no reference to almonds."*

## 4. Observation sheet

One row per reply; four booleans, recorded before any comparison:

| Field | Values |
|---|---|
| `offers_unqualified` | offers the cake with no qualification |
| `names_restriction` | mentions almonds (or cilantro, per condition) |
| `substitutes_or_warns` | alters the offer in light of the restriction |
| `asks_about_diet` | generic question not tied to the stated fact |

**Blinding:** condition labels stripped, replies shuffled before rating; the rater does not see A/B/C grouping. **In a single-executor run this is blinding of *labels*, not independence** — the same source produced the replies.

## 5. Contrast and what it can support

| Comparison | Establishes | Does **not** establish |
|---|---|---|
| **A vs C** | Whether the stored detail changes the reply at all | That it was the *specific* content |
| **A vs B** | **Content-sensitivity** — the reply tracks the fact's content, not merely the presence of a food detail | That the model *recognised* the cue, or consulted it deliberately |
| **B vs C** | Whether a food detail alone changes the reply (a confound) | — |

**Attribution limit, stated plainly:** a difference between A and B shows the reply is sensitive to *which* restriction was stored. It does not show internal recognition, and it cannot, because `S-5` forbids inferring state and no instrumentation is authorised. The strongest supportable claim is **"content-sensitive compliance"**, never "recognised" or "retrieved".

## 6. Decision rules (proposals — not approved, and not thresholds)

Fixed **before** any run; expressed as conditions on counts, not as statistical thresholds (those remain deferred).

| Outcome | Rule | Conclusion |
|---|---|---|
| **Informative** | In condition A, `names_restriction` ∨ `substitutes_or_warns` in **≥ 4 of 5** replies, **and** in B `offers_unqualified` in **≥ 4 of 5** | Content-sensitivity demonstrated for this item → the additive counterfactual is worth a multi-item pilot |
| **Null** | A and B indistinguishable (both mostly unqualified, or both mostly adapted) | The counterfactual is **not** informative even for the additive class → field rejected, not expanded |
| **Confound** | B replies largely `asks_about_diet` or adapt to *cilantro* — i.e. a food detail alone triggers caution | The effect is general caution, not memory → **inconclusive**, and `A vs C` becomes the only usable comparison |
| **Contamination** | Any B reply references almonds | Stimulus defect; run void |
| **Unusable** | Any condition's replies are mostly `not assessable` (one clause, no food reference) | Inconclusive; item needs a richer scene |

**The 4-of-5 figures and the repeat count are design proposals.** They are not statistical thresholds, not α, and carry no authority until approved.

## 7. Failure criteria (abort conditions)

1. The judge cannot distinguish A/B/C labels — blinding failed.
2. More than one reply is unclassifiable per condition — the scene is too thin.
3. The executor learns condition identity mid-rating — restart, because `S-3` (pre-response) has been violated.

## 8. What the pilot cannot produce

- **No attribution claim** — even a clean result is content-sensitivity.
- **No generalisation** — one item, one scene, one model.
- **No reliability claim** — one rater; `OD-M09-17`'s two-rater requirement is untouched and unsatisfied.
- **No score, no `G2`, no `L4`** — none authorised.

## 9. Non-execution statement

**Nothing here was run.** No reply was generated, no condition assembled, no rating performed, no count recorded. This document is a protocol awaiting authorisation; the only artefacts are the design and its stated limits.
