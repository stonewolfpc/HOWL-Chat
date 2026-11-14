# `OD-M09-03` — Layered Cue Standard

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `OD-M09-03-CUE-1` · **Version:** 3 · **Date:** 2026-11-13 · **Status:** **DIRECTION BINDING · S-1…S-6 ADOPTED AS DESIGN STANDARD · MECHANICS PROPOSED · NOTHING IMPLEMENTED.**
**Revision 2** (2026-11-13) — narrow clarification: separates criterion-design failure from evidence unavailability, makes named explanations distinct from supported ones, states the single-source ceiling explicitly, retains `F6-04` unchanged. Adopted no mechanics.
**Revision 3** (2026-11-13, owner-authorized) — records `OD-M09-15`…`OD-M09-18`: S-1…S-6 adopted as the standing design standard; `F6-04` retained unchanged; a two-rater study required before promotion beyond single-source confidence; `P-3` kept separate with per-element design not authorised. **No** item, tag, threshold, score or runtime change. See §13.
**Authority:** `OD-M09-03` in `docs/research/MEMORY-09-OWNER-DECISION-RECORD-OD-M09-2026-11-13.md` §2.3 — the layered development direction approved; `CORPUS-78`'s operational definitions remain **proposals**; counterfactual judgements split into recorded / supported / unresolved-reproducibility; `P-3` kept separate unless a per-element rule is written.
**Effect on sources:** **none.** No item, tag, table, score, threshold, checker or standard was changed in producing this document.

---

## 1. What is binding and what is not

| Element | Status |
|---|---|
| The **layered direction** — screen, then positive cue, then historical attribution, then substitution reserved | **Approved direction** (`OD-M09-03`) |
| The requirement that a cue be **specified before** an item is scored on it | Approved direction |
| Recorded and supported as requirements on counterfactual judgements | **Proposed requirement** — proposed here in concrete form |
| Reproducibility | **Required for promotion beyond single-source confidence** (`OD-M09-17`). The two-rater study is **not designed, not conducted and not authorised**; no threshold, statistic or observed agreement is asserted |
| `CORPUS-78` §4.1 mode-relative screen directions (blocking for perceptible modes, inverted for withholding modes) | **Proposal** — a reasoned structural finding, not an adopted rule |
| `CORPUS-78` §4.2–§4.5 mechanics, its §10 amendment language, its decision IDs | **Proposals / suspended** — not revived by this document |
| `P-3` merge | **Refused, and per-element design not authorised** (`OD-M09-18`). Separate unless a concrete evaluation need is established and the design work separately authorised |
| The **§3 sufficiency test S-1…S-6** | **ADOPTED** as the standing design standard (`OD-M09-15`). Adoption makes **no** item's criterion sufficient: no criterion has been authored or assessed against it |
| Item-specific cue criteria | **None adopted, none authored, none assessed.** The standard is a test, not a set of cues |

Nothing in this standard makes any item scorable that was not scorable before. It describes what a cue criterion must contain **if** one is later proposed.

---

## 2. Layers

| Layer | Question | Can it establish? | Status |
|---|---|---|---|
| **L0** — construct and routing | What is this mode asking, and is it a reply-observable predicate at all? | No — design only | Approved: L0 work only, routing **open** (`OD-M09-07`) |
| **L1** — counterfactual screen | Could a character without the memory plausibly produce substantially the same response? | No — it can **defeat** a claim, never establish one | Direction approved; mode-relative directions proposed |
| **L2** — positive observable cue | Is a specified, memory-linked observation present in the response? | Yes, **conditionally** — at cue level, never at causal level | Direction approved; criteria proposed here |
| **L3** — historical attribution | Did the memory actually influence this response? | **No** — no validated causal-evidence channel exists | Direction approved as *supporting evidence only* |
| **L4** — substitution comparison | Does removing the memory change the response under control? | Would be the strongest available test | **Deferred** as validation work; not authorised |

**The layering exists so that no layer is silently promoted.** A passed L2 cue is not an L3 attribution; a failed L1 screen is not proof of absence; L4's deferral is not evidence either way.

---

## 3. The sufficiency test for a cue criterion

A cue criterion is **specified** only if all six hold. This is the checkable content of "adopted or sufficiently specified" in `OD-M09-06`.

**Status: adopted 2026-11-13 by `OD-M09-15`** as the standing design standard. The six requirements below are adopted **as worded**; they are not to be silently weakened, strengthened, reinterpreted or replaced. **No existing item is thereby declared to have a sufficient cue criterion** — establishing that requires a separate, authorised item-level assessment that has not been performed.

| # | Requirement | Failure mode it prevents |
|---|---|---|
| **S-1 · Anchored** | The cue is tied to the item's own declared memory and scene, not to a generic behaviour | Cue that any reply could satisfy |
| **S-2 · Distinguishable** | The cue is **provably not** the item's mandatory acceptance text, nor a paraphrase of it | `F6-04` collapse — a separate mode result that is really the acceptance field restated |
| **S-3 · Pre-response** | The criterion can be applied without knowing what the model did, and before seeing adjacent items | Retrospective cue-fitting |
| **S-4 · Directional** | The criterion states which direction counts as satisfied for **that** mode family | A single uniform direction for blocking and inverted modes |
| **S-5 · Observational** | The cue is stated as something visible in the response, never as an inference about the model's state | "The model remembered" masquerading as an observation |
| **S-6 · Refutable** | A memory-absent explanation can be named in one sentence against it | Unfalsifiable cue |

**S-2 is where `F6-04` fails today.** Its mandatory text *is* the perceptible accommodation, so a mode result derived from it would be the acceptance field in different clothes. Until the overlap is resolved, that item is **mode not assessable** — not "mode not demonstrated", and certainly not a pass.

---

## 4. Three states, and what each may be used for

| State | Condition | Permitted use |
|---|---|---|
| **Mode not assessable** | Either no criterion meeting §3 exists, **or** a criterion exists but the required evidence is unavailable pre-response — see §4.1 | **Nothing.** May not be counted as failure, and may not be excluded silently either — it is reported as unassessable, with its cause named |
| **Mode not demonstrated** | Criterion exists and is applicable; the cue is absent, or defeated by a named memory-absent explanation | Reporting only — it never gates an item on its own (`OD-M09-02` anti-collapse) |
| **Mode demonstrated** | Criterion exists; cue present; no memory-absent explanation survives | Cue-level result only |

**Explicit prohibition:** `Mode not assessable` must never be reported as `Mode not demonstrated`. The distinction is the reason the two-axis vocabulary was adopted, and collapsing it would manufacture failure rates out of missing standards.

### 4.1 Two causes of unassessability, reported separately

One state, two unrelated causes and two unrelated remedies. Collapsing them would let a design gap look like an evidence problem, or an evidence problem look like a design gap that someone is already fixing.

| Cause | Condition | Remedy | May it change with execution? |
|---|---|---|---|
| **Design** — no sufficient criterion | Nothing meeting §3 exists for the item | Author a criterion (design work) | No — running the item cannot produce a criterion |
| **Evidence** — criterion exists, evidence does not | A criterion exists, but the evidence it requires is unavailable under pre-response conditions | Obtain the evidence, or record the criterion as inapplicable to this item | Possibly — this cause is contingent on materials, not on design |

**Reporting rule:** every `Mode not assessable` entry names its cause (`design` or `evidence`). A blank cause is itself a defect in the report.

**Neither cause is a failure.** A design-gap item has not failed the mode; it has not been asked the question yet. An evidence-gap item may be assessable later without any change to the criterion — which is exactly why the two must not share a label that implies the same fix.

---

## 5. The counterfactual judgement — recorded, supported, reproducible

### 5.1 Recorded (proposed requirement)

Each judgement must state, in the record:

1. the memory-absent explanation considered;
2. the evidence in the item that makes that explanation plausible or implausible;
3. the resulting direction — blocked (perceptible modes) or leakage/over-use suspected (withholding modes);
4. the rater's one-sentence conclusion.

A judgement without (2) is an assertion, not a recorded judgement.

### 5.2 Named is not supported

Naming a memory-absent explanation satisfies the **recording** requirement and nothing more. Three states, only the last of which carries evidential force:

| State | Condition | Consequence |
|---|---|---|
| **Named** | An explanation is stated and reproduced in the record | Enters the record. Cannot, by itself, defeat a cue or block a claim |
| **Named, unsupported** | Stated, but not grounded in the item's declared evidence | Reported as recorded and unsupported. Carried as an open question, never as a refutation |
| **Supported** | Grounded per §5.3 | May defeat a cue claim, and must be answered before a cue counts as demonstrated |

A named-only explanation that blocks a claim without grounding turns rater intuition into a verdict. That is the failure `CORPUS-79` F-11/F-12 identify, and it is not repaired by the explanation being written down.

### 5.3 Supported (proposed requirement)

The explanation must be **grounded in the item's declared evidence** — scene details, stated facts, the reference answer's own reasoning — and not in the rater's model of how assistants usually behave. Support is a property of the link between the item's text and the explanation, and is checkable by a second reader without re-running anything.

### 5.4 Reproducible (unresolved)

Whether two qualified raters applying the same rule reach the same judgement is **unknown and unmeasured**. This standard does not:

- set an agreement threshold;
- name a statistic;
- claim any observed agreement;
- imply that single-rater judgements are therefore unreliable.

It records only that reproducibility is **required before cue-level results carry weight beyond single-source confidence**, and that no measurement exists.

**Owner decision `OD-M09-17`:** a two-rater reproducibility study is **required** before cue-level results support claims beyond single-source confidence. The study's protocol must specify, before execution: rater qualification and provenance; independent assessment procedures; blinding and available information; the exact judgement compared; how disagreements are recorded and handled; the reliability measure and its justification; any acceptance threshold, fixed before results are seen; and how insufficient sample size, missing judgements and unresolved disagreements are reported. **None of this is executed, scheduled or authorised here**, and no statistic or threshold is named by this standard.

### 5.5 Single-source confidence ceiling

Every cue-level result produced under this standard is capped at **single-source confidence**. It receives **no convergence credit**, and formal convergence remains **0 of 2**.

| Permitted | Not permitted |
|---|---|
| Reporting what one judge found, at single-source confidence | Aggregating repeated judgements as though they were independent |
| Using that result to describe a single item | Describing it as replicated, corroborated or convergent |
| Citing it alongside the limitation | Treating a second reading of the same material as a second source |

The cap is a ceiling on what the result can support — not a judgement that single-source cue results are unreliable. It is the same ceiling `OD-M09-04` applies to review evidence, applied here to cue judgements.

---

## 6. `P-3` — separate, and what would be needed to merge it

`P-3` remains a **separate** test. The merge is refused by default, because a merge without a per-element rule loses the distinct information `P-3` carries when a response contains one derivable and one non-derivable element.

**Owner decision `OD-M09-18`:** `P-3` is kept separate, and **per-element design work is not authorised** at this time. The content below is explanatory design guidance, not a work order. A future merge request must establish a concrete evaluation need and be separately authorised.

A per-element rule sufficient to justify merging would have to state, for **each** mode direction:

| Required content | Why |
|---|---|
| The element set (what counts as one element) | Otherwise "per-element" is undefined |
| Which elements the screen applies to | A screen that applies to "the response" is the aggregate test again |
| How a mixed response is scored (some blocked, some not) | The case `P-3` exists for |
| The direction for blocking modes **and** inverted modes separately | §4.1's asymmetry cannot be collapsed |
| What result a mixed response reports | Otherwise the merge produces an unlabelled average |

**No such rule is written.** Until it is, `P-3` stays separate. Writing it is a design task, not a documentation task, and it is not authorised here.

---

## 7. Three things that must stay apart

| Level | Question | Established? |
|---|---|---|
| **Cue observable** | Is the specified observation present? | **Not established** for any unit — no criterion adopted, no execution |
| **Response compatible with memory use** | Could memory explain this response, with no memory-absent alternative surviving? | **Not established** for any unit |
| **Memory influenced behaviour** | Did the memory cause it? | **Not established.** No current channel establishes it; prompt capture, as it exists today, shows injection only. Whether a future validated method could establish influence is **open** — this standard asserts no impossibility |

Confusing the first two inflates cue results into compatibility claims. Confusing the second and third converts a plausibility judgement into causation. Neither inference is permitted anywhere in this standard.

**Precedent already in the corpus:** `R6-03d` states that unknown independence is not established corroboration. The same discipline applies here — an unevaluated relationship is not a weak relationship, and an unrefuted explanation is not a demonstrated cause.

---

## 8. Constraints carried in from `OD-M09-05`

| Constraint | Effect on this standard |
|---|---|
| Tags are the adopted membership authority | Cue criteria attach to tagged items — 46 units, one mode each |
| `R-1` is operative but substantively unresolved | `F6-02` and `P-01` are `R-1` **by authority**; whether they exercise suppression is **open**. Any `R-1` cue criterion inherits that question and must state it |
| The re-tagging account is a hypothesis | `OD-M09-05` §4.1 offers a labelled reconstruction of how the table came to contradict `A-2`. It is **not** provenance: neither document asserts when or why the table changed, and this standard cites the table only as a divergent record |
| `R-6` = 10 tagged units, 8 dedicated `R6-*` | Ten is the membership count. Eight is the fixture count. They are not interchangeable in sample-size or coverage reasoning |
| Table omissions are not disproof of intent | An item's absence from `CORPUS-01` §5 says nothing about whether it was meant to carry a mode |
| No re-tagging by this standard | Membership is read, never written |

**Two-tier consequence for `R-1`.** Work that needs only the adopted membership authority may proceed: counting `R-1` units, sampling them, including them in coverage arithmetic. Work whose **validity** depends on those two items genuinely exercising suppression may not treat the question as settled. An `R-1` cue criterion, a suppression rate, or any claim that `R-1` is adequately covered all inherit the open question, and must carry it explicitly rather than resolving it by citation.

---

## 9. Illustrations from existing items only

No fixture is invented here. Existing documented cases show how §3 and §4 apply:

- **`F6-04`** — cue overlaps the acceptance field (S-2 fails) → **mode not assessable**, cause **design**. `CORPUS-79` F-9 already records that `A-4`'s "changes reporting only" claim does not hold for this shape. **Its wording is retained unchanged** (`OD-M09-16`): this standard does not edit its mandatory text, does not restate its acceptance field, does not change its tag, and does not reclassify it — it is neither a pass nor a failure. Its status is **not necessarily permanent**: a future item-level change is possible and requires separate owner authorization. Any change to that item is a separate item-level decision that has not been authorised.
- **`F1-03`, `F6-01`, `F3-01`** — carried forward as **UNRESOLVED** by `CORPUS-78` §7; this standard neither resolves nor re-diagnoses them.
- **`R6-03a`…`R6-03d`** — provenance variants whose rule is stated in `CORPUS-01`; they illustrate S-1 and S-6 without becoming scored examples.
- **`F7-01`, `F9-04`** — `R-6` members that predate the dedicated series; their cues, if any, must clear S-2 against their own mandatory text.

---

## 10. What this standard does not do

- Does not adopt `CORPUS-78` §4.2–§4.5, its §10 amendment language, or its decision IDs.
- **Adopts** the §3 sufficiency test as a design standard (`OD-M09-15`) — but adopts no item-specific cue criterion, assesses no existing item against it, and judges no criterion sufficient.
- Does not change any item, tag, table, threshold, score, split or runtime behaviour.
- Does not merge `P-3`, and does not authorise per-element design (`OD-M09-18`).
- Does not design, schedule or execute a reproducibility study (`OD-M09-17`).
- Does not implement gating — `OD-M09-02` is the gating decision, and its anti-collapse clause is a precondition, not a consequence, of this document.
- Does not score, tag, retag, repair or reclassify any item.
- Does not write the per-element rule, does not merge `P-3`.
- Does not authorise L4 substitution work, `OD-M09-08` source corpus, or any evidence collection.
- Does not claim any cue is observable, any response compatible, or any memory causal.

---

## 11. Dispositions of the four owner questions

| Question | Disposition | ID |
|---|---|---|
| Adopt the §3 sufficiency test as the definition of "sufficiently specified"? | **Adopted** as the standing design standard. No criterion is made sufficient by adoption | `OD-M09-15` |
| Rule on the `F6-04` overlap? | **Retained unchanged** — mode not assessable, cause design. No item edit authorised; not necessarily permanent | `OD-M09-16` |
| Reproducibility? | **Study required** before promotion beyond single-source confidence. Not designed, not executed, not authorised | `OD-M09-17` |
| The per-element rule for `P-3`? | **Kept separate**; per-element design not authorised at this time | `OD-M09-18` |

**Still open, and not decided by these four decisions:** which unit the corpus's 45-versus-46 scored-unit arithmetic excludes (`OD-M09-05` D-5); whether `F6-02` and `P-01` substantively exercise `R-1` (`OD-M09-05` D-1); and the routing of `R-1`/`R-4`/`R-5` beyond L0 design (`OD-M09-07`).

---

## 12. Provenance

| Layer | Fact |
|---|---|
| Decision recorded | `OD-M09-03`, 2026-11-13 |
| Authored by | Buffy / Codebuff repository worker — artifact-reported identity; labels are provenance data, not lineage proof |
| Model executed | No · `P-3` not started · no inference, benchmark or `G2` |
| Sources changed | None — this document is the only artifact created |
| Limit | A standard for criteria, not a set of criteria. It cannot be applied to any unit until specific cues are written under §3 and adopted |

---

## 13. Revision log

| Rev | Date | Change |
|---|---|---|
| 1 | 2026-11-13 | Initial standard: layer model, six-part sufficiency test (proposed), three states, recorded/supported/reproducible, `P-3` kept separate, three-way distinction, `OD-M09-05` constraints |
| 3 | 2026-11-13 | **Owner-authorized revision (`OD-M09-15`…`OD-M09-18`); not an independent research finding.** S-1…S-6 moved from proposed to **adopted design standard** with no wording change and no item declared sufficient; `F6-04` retained unchanged, cause design, explicitly not necessarily permanent; reproducibility upgraded from "unresolved" to "**required study**, not designed/executed/authorised"; `P-3` separation reaffirmed with per-element design not authorised; §7 narrowed so no impossibility claim is made about future validated methods; §10 and §11 rewritten to state the dispositions; no item, tag, threshold, score, split or runtime change |
| 2 | 2026-11-13 | **Clarification only, no mechanics adopted.** (a) §4.1 splits `Mode not assessable` into design-cause and evidence-cause with separate remedies and a mandatory cause label; (b) §5.2 makes *named* distinct from *named, unsupported* and *supported*, with evidential force attaching only to the third; (c) §5.4/§5.5 state the single-source ceiling and forbid independence-by-repetition; (d) `F6-04` retained unchanged and still mode-not-assessable, cause design; (e) `P-3` still separate; (f) the §3 sufficiency test explicitly **not adopted**; (g) no item, tag, threshold, score, split or runtime change |
