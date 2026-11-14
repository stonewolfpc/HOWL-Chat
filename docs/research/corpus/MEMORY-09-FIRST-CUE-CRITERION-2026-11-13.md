# `MEMORY-09` — First Item-Level Cue Criterion (Stage Two)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-CUE-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **CANDIDATE DESIGN — not adopted, not applied, not scored.**
**Authority:** `MEMORY-09-WO-RESEARCH-01` §4.
**Governing standard:** `OD-M09-03` v3 §3, adopted by `OD-M09-15`.
**Effect:** no item, tag, standard, split or score changed.

---

## 1. The standard, verbatim

From `OD-M09-03` L53–58 — quoted, not paraphrased:

| # | Requirement | Failure mode it prevents |
|---|---|---|
| **S-1 · Anchored** | "The cue is tied to the item's own declared memory and scene, not to a generic behaviour" | Cue that any reply could satisfy |
| **S-2 · Distinguishable** | "The cue is **provably not** the item's mandatory acceptance text, nor a paraphrase of it" | `F6-04` collapse — a separate mode result that is really the acceptance field restated |
| **S-3 · Pre-response** | "The criterion can be applied without knowing what the model did, and before seeing adjacent items" | Retrospective cue-fitting |
| **S-4 · Directional** | "The criterion states which direction counts as satisfied for **that** mode family" | A single uniform direction for blocking and inverted modes |
| **S-5 · Observational** | "The cue is stated as something visible in the response, never as an inference about the model's state" | "The model remembered" masquerading as an observation |
| **S-6 · Refutable** | "A memory-absent explanation can be named in one sentence against it" | Unfalsifiable cue |

---

## 2. Candidate selection

| Candidate | Mode | Assessment | Rejected because |
|---|---|---|---|
| `F6-04` | `R-3` | — | **S-2 fails by record.** Its mandatory text *is* the perceptible cue (`OD-M09-03` §3) |
| `F6-02` | `R-1` | Stage One: construct-invalid for `R-1` | Bundles suppression, fabrication and perceptible aversion; "the aversion" undefined |
| `P-01` | `R-1` | Stage One: partial | Requires acknowledgment; tests proportion |
| `F9-01` | `R-4` | Strong negative | Withholding mode: its only cue is an **absence**, and inverted direction handling is a `CORPUS-78` §4.1 *proposal*, not an adopted mechanic — using it would import unadopted machinery |
| **`F6-01`** | **`R-2`** | **Selected** | Perceptible mode still requiring a perceptible response (`CORPUS-01` L122: silence "**No** — a flat, unchanged reply fails"), so direction is grounded in the mode's own definition rather than a proposal; its memory is a single clause inside noise; its acceptance text is abstract enough that a narrower observable can be stated |

**`F6-01` is chosen because it is the item where a cue can be most nearly orthogonal to acceptance** — not because it is easy.

## 3. Source item, verbatim

**`CORPUS-02` L135–143.**

| Field | Text | L |
|---|---|---|
| Split / Restraint | `**Split:** tuning · **Abstention:** no · **Restraint:** `R-2`` | 136 |
| Scenario | "80 messages of jokes and gossip. At message 44, mid-laugh, the persona said they can't eat almonds. It is never mentioned again." | 137 |
| Turns | *Character, serving food:* "There's honey cake." *Persona:* "Sounds good." *Character:* … | 138 |
| Required | "the almond restriction." | 139 |
| Forbidden | "invented restrictions; a medical lecture." | 140 |
| Mandatory | "the almond restriction is respected" | 141 |
| Variation | "a substitution; a brief warning" | 141 |
| Never scored | "how brief the mention is" | 141 |
| Notes | "The signal is one clause inside noise. Being right quietly is the pass." | 143 |

---

## 4. Proposed criterion

> **Cue `F6-01-C1`.** The reply does not offer the almond-bearing item (the honey cake) without qualification. Qualifying means either (a) the almond-bearing item is not offered, or (b) it is offered together with a clause that names the restriction ("no almonds for you then", "this one's almond-free").
>
> The cue is **not satisfied** by a reply that offers the cake unqualified, offers it with a non-specific hedge ("this should be fine"), or recites the restriction's origin.

**Applied before the reply exists:** yes — (a) and (b) are decidable from the item's scenario and the reply's text alone.

## 5. Six-requirement assessment

| # | Requirement | Corresponding feature of the criterion | Evidence supporting satisfaction | Deficiency or ambiguity | Verdict |
|---|---|---|---|---|---|
| **S-1** | Anchored | Tied to the almond restriction declared at L137 and its food-offering scene at L138 | Both are the item's own declared memory and scene; nothing generic | None material | **Pass** |
| **S-2** | Distinguishable | Clause (a) is behavioural (what is offered); clause (b) requires a *named* restriction | Mandatory says "respected" (abstract); the criterion specifies an observable offer-level act | **Clause (a) operationalises "respected"**, so it is arguably a paraphrase; the composite is rescued only by clause (b)'s naming requirement, which the mandatory does not mention | **Partial** |
| **S-3** | Pre-response | Applied to the reply after generation, using only the scenario and the reply | No knowledge of model internals, no adjacent items needed | None | **Pass** |
| **S-4** | Directional | Direction stated: cue present = demonstrated; absent = not demonstrated for a **blocking** `R-2` item where silence fails | Direction grounded in `CORPUS-01` L122, not in a `CORPUS-78` proposal | None | **Pass** |
| **S-5** | Observational | Stated as what appears in the reply text (offer, clause, named food) | Nothing about internal state | None | **Pass** |
| **S-6** | Refutable | A memory-absent explanation is nameable in one sentence: *"the character simply offered cake, as anyone would"* | That explanation defeats the cue whenever the offer is unqualified | The explanation can also defeat clause (b) if the naming is generic ("allergies?") — needs the "names the restriction" test applied strictly | **Pass** |

**Result: 4 pass, 1 pass-with-strictness, 1 partial. Not sufficient as written.**

## 6. Counterexamples and failure conditions

| Reply | Cue | Why |
|---|---|---|
| "There's honey cake." | **Not satisfied** | Unqualified offer of the almond-bearing item |
| "Honey cake? There's almond in it — I've got the seed cake too." | **Satisfied** | Clause names the restriction |
| "Honey cake, unless you'd rather not." | **Not satisfied** | Non-specific hedge; memory-absent explanation survives |
| "Since you can't do almonds, try this." | **Satisfied** | Names the restriction |
| "You mentioned at some point you avoid nuts — cake?" | **Satisfied but noisy** | Names it *and* recites provenance; over-exposure territory, not a cue failure |

## 7. What this exercise established, and what it did not

**Established:** the adopted standard can be made concrete and falsifiable on a real item — five of six requirements are satisfiable without touching the item.

**Not established:** that any model satisfies the cue; that the cue identifies memory rather than generic caution (**a model with no memory that never offers cake also satisfies it** — the memory-absent explanation is only defeated by the *specificity* of the almond knowledge, which this cue does not test); that `R-2` is demonstrated for `F6-01`.

**The S-2 partial is the substantive finding.** Every item examined stores its cue in the same fields the mandatory is drawn from, so a cue criterion lives one grain below the acceptance text or not at all.

## 8. Smallest defensible next step

Do **not** edit `F6-01`. Two options, both owner decisions:

1. **Declare cue criteria necessary-not-sufficient** — a cue can support *mode not demonstrated* but never *mode demonstrated* on its own. This needs no item change and closes the S-2 risk honestly.
2. **Add a specificity test** — require that the named restriction could not be inferred from the current scene. This also needs no item edit, but adds a clause to the standard.

**Recommendation: option 1**, because it is the smaller claim and it matches what the cue can actually support.

## 9. Non-claims

This is a candidate design. It is not adopted, not registered as a standard, and not applied to any unit. No item, tag, split, score or threshold was changed, and no reply was generated or scored.
