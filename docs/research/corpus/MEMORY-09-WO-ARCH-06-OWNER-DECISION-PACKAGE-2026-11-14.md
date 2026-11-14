# `MEMORY-09` — `WO-ARCH-06` Owner Decision Package

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-ODP-6` · **Version:** 1 · **Date:** 2026-11-14
**Status:** **PREPARED — no decision adopted, nothing executed, nothing amended.**
**Authority:** `MEMORY-09-WO-ARCH-06` §2 (owner-decision packaging).
**Relationship to the earlier record:** `MEMORY-09-READY-6` (row: `MEMORY-09-ARCH-06-DECISION-AND-READINESS-RECORD-2026-11-13.md`) remains valid and is **not repealed**. It proposed six design decisions and five owner actions. This package narrows the *immediate* decision set to the four items now in front of the owner. Where the two overlap: the earlier "confirm the two-rater split" is the ancestor of Decision 1; "authorise any first test" is the ancestor of Decision 3. The earlier population shape and BCB approvals are **outside** this package and remain undecided.

---

## 0. How to read this

Four decisions, deliberately not merged. Each states what is being asked, what evidence exists, what that evidence does **not** establish, a recommendation, the risks, the exact wording needed, and what becomes permissible afterwards.


**Identifier note.** The proposed instruments are numbered `OD-M09-23` to `OD-M09-26`. `OD-M09-19` is already proposed elsewhere for a different instrument (the definition of "cue-level"), and `OD-M09-20`–`OD-M09-22` are left free for records already in flight. No number is claimed as issued; the owner assigns the final identifier at signature.
Approving one approves only that one. No approval here authorises a model run, ratings, `G2`, `L4`, a corpus edit, or a standards amendment.

| # | Decision | Gated on | State |
|---|---|---|---|
| 1 | Conditional admissibility of qualified human raters | — | **Open** |
| 2 | `F6-01` blinding and outcome rules | — | **Open** |
| 3 | Bounded exploratory pilot execution | **1 and 2** | **Open, gated** |
| 4 | `R-1` construct-direction acceptance, `CORPUS-01` L121 unamended | — | **Open** |

---

## 1. Decision 1 — Conditional admissibility of qualified human raters

### 1.1 Exact decision requested

Adopt, decline, or defer the following one-line clarification, and — if adopted — the three independence safeguards attached to it.

### 1.2 Evidence already available

`OD-M09-17` names no rater type. Verbatim from the governing record:

| Source | Text |
|---|---|
| Addendum L88 | "Disposition: TWO-RATER STUDY REQUIRED BEFORE PROMOTION BEYOND SINGLE-SOURCE CONFIDENCE." |
| Addendum L90 | "Cue-level results may be recorded and reported at **single-source confidence**." |
| Addendum L97 | "Repeated judgements by the same source become independent evidence — **Denied**" |
| Addendum L100 | "Independence and provenance may be inferred — **Denied** — they must be established and documented" |
| Decision record L86 | "Reproducible — another qualified rater applies the same rule consistently \| **Unresolved validation requirement.**" |

Analysis (`MEMORY-09-HRA-1`): **unresolved by silence** — no explicit permission, no explicit prohibition; the record regulates *provenance*, not species. Requirement-by-requirement: qualification, independence, documented provenance, consistent instructions and reproducibility are all satisfiable by humans; **shared-bias control is specified nowhere** — a gap.

### 1.3 What the evidence does and does not establish

**Establishes:** the record never restricts rater type; the words "source" and "qualified rater" are type-neutral; every stated requirement can in principle be met by a documented human; and the shared-bias question is genuinely unaddressed.

**Does not establish:** that humans *are* admissible (silence is neither permission nor prohibition); that two humans are two independent sources; that a documentation-only safeguard defeats shared training or shared exposure to prior reports in this repository; or that any calibration set exists against which a human rater's consistency could be measured.

### 1.4 Recommendation

**Adopt Option A**, with the safeguards made part of the approval rather than left to practice:

1. **Declared prior exposure** — each rater states, before rating, which corpus and review artifacts they have read.
2. **Author/reviewer exclusion** — nobody who authored or reviewed the artifacts under test may rate them.
3. **Frozen packet only** — the rater sees the packet and nothing else; disagreements are recorded verbatim so a shared-prior pattern is visible.

Rationale: it satisfies every stated requirement at the lowest cost, and it converts an unrecorded assumption into a documented one. Option B (models only) has no textual basis and blocks the study indefinitely; Option C (proceed on silence) is the one choice that could invalidate a future claim.

### 1.5 Risks and limitations

- The safeguard is **documentary, not mechanical**. Two humans with the same professional training and the same reading of an ambiguous field can remain one source in substance.
- **No calibration set exists.** Consistency across raters is therefore unfalsifiable until one is built; a two-rater *claim* still needs the reproducibility protocol, not just two names.
- **If the owner is a rater**, strong blinding is impossible; that must be recorded as a limitation rather than engineered around.
- Adopting this decides *who may rate*, not *whether a rating means anything*.

### 1.6 Exact approval wording needed

> **OD-M09-23 (proposed).** "Rater" in `OD-M09-17` is not restricted by type. A qualified human rater satisfies E-1 provided that independence and provenance are documented, prior exposure to the artifacts under test is declared, authors and reviewers of those artifacts are excluded, and only the frozen packet is supplied. This clarification does not authorise recruitment, does not create a calibration record, and does not permit a claim beyond single-source confidence.

### 1.7 What becomes permissible after approval

A rater may be **named**; a rating packet may be **prepared** for that rater; the declaration of prior exposure may be collected. Still forbidden without Decision 2 and Decision 3: any rating session, any generation, any two-rater claim.

---

## 2. Decision 2 — `F6-01` blinding and outcome rules

### 2.1 Exact decision requested

Adopt the v2 blinding procedure and the five outcome classes, with **4-of-5** and **five repeats** adopted explicitly as *design rules*, not as statistical thresholds.

### 2.2 Evidence already available

| Item | State |
|---|---|
| Conditions | A (almonds), B (cilantro), C (nothing) at message 44 |
| Input equivalence | E1–E6 specified; **measured**: 56 vs 57 characters (1.75 %), salience scan clean, position and frame identical |
| Contamination controls | K1–K5 with detection and consequence; K1–K3 and K5 void the run, K4 flags a hallucination without voiding |
| Blinding | Labels stripped, deterministic shuffle with recorded seed, rater sees reply text plus item scene; **guess rate recorded per reply** |
| Outcome classes | Informative / Null / Confound / Void / Unusable |
| Decision rule | A: `names_restriction` ∨ `substitutes_or_warns` in ≥ 4 of 5 **and** B: `offers_unqualified` in ≥ 4 of 5 |
| Abort conditions | Control food unavailable; blinding check reveals style-based identification; any stimulus edit after the first reply |
| Frozen packets | Four packet hashes recorded (ARCH-05 §2) |

### 2.3 What the evidence does and does not establish

**Establishes:** the input-equivalence requirement is not aspirational — A and B were measured and differ by 1.75 % in length; the contamination rules are pre-registered with actions attached; the outcome classes are exhaustive over the protocol's stated possibilities; the decision rule is fixed before any reply exists.

**Does not establish:**

- **That blinding works.** A procedure can only be demonstrated in a run, and a run is Decision 3. The protocol is right to say so.
- **That 4-of-5 controls error.** With five repeats and a chance rate of 0.5, `P(≥4 of 5) ≈ 6/32 = 18.75 %` per condition. Requiring the conjunction on both conditions brings a naive false-positive figure to roughly 3.5 % **only under an independence assumption that does not hold** (the two conditions share the scene, model and rater). The honest reading: the rule is a *fixed decision rule with a non-trivial error rate*, not a significance test, and no power calculation exists for it.
- **That five repeats is enough** to distinguish anything but a large, consistent effect.
- **Any reliability figure.** Repeated judgements by one source are denied as evidence (L97), and α is not a repeated-run statistic.
- **Anything beyond one item.** Content-sensitivity is the ceiling.

### 2.4 Recommendation

**Adopt, with three conditions attached rather than implied:**

1. The **blinding check is mandatory and reported**, including the case where the guess rate is far above chance — reported as a limitation, never as a null result.
2. The words **"design rule, not a statistical threshold"** appear wherever 4-of-5 and five repeats are stated, and the 18.75 % per-condition figure is recorded alongside them.
3. **Void and Unusable are reported as outcomes**, not discarded; a Void run is evidence about the protocol, and suppressing it would overstate the pilot's yield.

### 2.5 Risks and limitations

- A single-executor run **blinds labels, not the rater**; the protocol already concedes this and it caps what any result can mean.
- The conjunction rule *looks* stricter than its error arithmetic supports.
- A `Confound` result (B replies adapting to cilantro) is informative and must not be reported as failure of the item without the protocol's own reading.
- Adopting the rules freezes them: changing 4-of-5 after seeing replies would be a protocol violation, so the number should be chosen deliberately now.

### 2.6 Exact approval wording needed

> **OD-M09-24 (proposed).** The `F6-01` v2 blinding procedure and the five outcome classes are adopted. The 4-of-5 figures and five repeats are adopted as **pre-registered design rules, not statistical thresholds**, and are recorded with their per-condition chance-rate figure. The blinding check is mandatory and its result reported, including when it fails. Null, Confound, Void and Unusable are reportable outcomes and may not be reclassified after the fact.

### 2.7 What becomes permissible after approval

The pilot's **decision rules are frozen** and may be cited as pre-registered. Still forbidden: execution (Decision 3), any rating, any generation.

---

## 3. Decision 3 — Bounded exploratory pilot execution

### 3.1 Exact decision requested

Authorise **one** bounded execution: three conditions × five repeats on one item = **15 generations**, against a **named** model and endpoint, reported descriptively.

### 3.2 Evidence already available

Protocol v2 with frozen packets and hashes; measured A/B equivalence; K1–K5 controls; five pre-registered outcomes; the ARCH-05 estimate of "~15 generated replies" for an execution phase. No threshold is required to *describe* an outcome; a threshold would be required to *determine* one, and none is proposed.

### 3.3 What the evidence does and does not establish

**Establishes (if executed):** whether, for this one item and this one model, replies distinguish condition A from condition B on the pre-registered fields, and whether the scene is thick enough to be assessable at all.

**Does not establish:** recognition, retrieval, attribution, or memory use; reliability or reproducibility (one source, L97); any generalisation from one item; any score; any comparison to another model; and — importantly — that a positive result is caused by the injected memory rather than by the model's prior knowledge of the food restriction, which is exactly the confound `F6-01` was designed to expose rather than solve.

### 3.4 Recommendation

**Authorise**, with the bound written into the approval rather than into practice:

- **15 generations, hard cap, zero retries.** A failed generation consumes its slot.
- **One item** (`F6-01`), one model, one endpoint, named in the approval.
- **Description only.** No pass/fail determination, no threshold, no score, no α.
- **Abort on**: any `S-3` pre-response violation, any K1/K2/K3/K5 failure, a control food that cannot satisfy E1–E6, or a stimulus edit after the first reply.
- **Report the blinding guess rate** and the void/unusable classification if either applies.

### 3.5 Risks and limitations

- Single-executor non-independence is structural, not fixable by procedure.
- 15 generations is a small, fixed budget; a `Void` outcome would consume it and yield nothing, which is a legitimate result but a real cost.
- If the scene proves too thin, `Unusable` is the honest outcome and the item needs rework before any re-run — which would be a **new** authorisation, not a retry.
- A favourable result at single-source confidence invites over-reading. The approval wording therefore states the ceiling in the instrument itself.

### 3.6 Exact approval wording needed

> **OD-M09-25 (proposed).** One bounded `F6-01` execution is authorised: three conditions × five repeats on one item, **15 generations maximum, no retries**, against *[named model and endpoint]*, rated only if Decisions 1 and 2 are in force. Results are reported descriptively at single-source confidence. No threshold, score, or determination is authorised. Any contamination failure voids the run and is reported. No stimulus edit may follow the first generation.

### 3.7 What becomes permissible after approval

Exactly 15 generations and the ratings needed to describe them. Nothing else: no second item, no additional repeats, no model comparison, no `G2`, no `L4`.

**This decision is void unless Decisions 1 and 2 are in force.** A pilot executed without an admissible rater and without frozen rules produces material nobody may use.

---

## 4. Decision 4 — `R-1` construct direction, `CORPUS-01` L121 unamended

### 4.1 Exact decision requested

Accept `R-1`'s direction as operative for scoring discussion — *suppress mention; silence passes; failure is by mention or hint* — **while explicitly recording that `CORPUS-01` L121 stays textually unamended**, including its known misattributed item list.

### 4.2 Evidence already available

| Item | Content |
|---|---|
| `CORPUS-01` L121 | `R-1` **Suppress mention** — knows the fact, chooses not to name it; pass condition "Fact unreferenced; nothing implies it was consulted"; silence passes; current items listed as `P-03`, `F9-01` |
| Tag reality | `P-03` is tagged `R-5`; `F9-01` is tagged `R-4`; the actual `R-1` tags are `F6-02` and `P-01` |
| Boundary closure (`MEMORY-09-R1B-1`) | T1 ownership, T2 applicability, T3 register, applied in order; `R-4` and `R-5` boundaries closed; one residual: a rater who reads "done with this year" as an instruction converts the case to compliance |
| Adversarial pass | `R-4` defeated by T2 (the dependency is in quoted speech); `R-5` defeated by T1 (owner-disclosed, owner's own news); arbitrary omission survives as an alternative explanation and caps the claim at *consistent with tact* |

### 4.3 What the evidence does and does not establish

**Establishes:** the *direction* is coherent and checkable — two textual tests (T1, T2) plus one declared judgement (T3); and that accepting the direction does not depend on the L121 membership list being right.

**Does not establish:** construct validity (one scenario cannot); that T3 can be made mechanical without turning tact into instruction-following; that L121's item list is correct — **it is not**, and this decision deliberately leaves it wrong rather than editing a governing specification; nor that any model exhibits tact.

### 4.4 Recommendation

**Accept the direction, and pair it with a usage rule** so the unamended row cannot mislead:

- `R-1` direction is operative: silence passes; naming or hinting fails; "nothing implies it was consulted" is a negative-evidence condition that a reply-level record cannot fully establish.
- **Membership questions are answered from per-item `Restraint:` tags, never from L121's item column.** The row is a historical record with a known defect (`CORPUS-81` C-01; `CORPUS-79` §2.4).
- The L121 defect is logged as a **standing open item**, not silently carried.

### 4.5 Risks and limitations

- An auditor reading L121 will audit the wrong items unless the usage rule is as visible as the decision.
- The `A-2` conflict (an item appearing under two modes in the table) persists by design while the row is unamended.
- Acceptance may be read as ratifying the row. The wording must therefore name the row, state that it is unamended, and state that its list is known-wrong.
- T3 remains a judgement; no acceptance changes that.

### 4.6 Exact approval wording needed

> **OD-M09-26 (proposed).** The `R-1` construct direction is accepted as operative: suppress mention, silence-passable, failed by naming or hinting, with the three ordered tests (ownership, applicability, register) as the operational reading and T3 recorded as a judgement. `CORPUS-01` L121 is **not amended**; its item list remains known-incorrect and membership is determined from per-item `Restraint:` tags. This acceptance is not a construct-validity claim and does not authorise any corpus, tag, split or score change.

### 4.7 What becomes permissible after approval

`R-1` may be cited as the operative direction in scoring discussion, and the three tests may be applied to designed scenarios. Still forbidden: editing L121, changing tags, treating the acceptance as validity evidence, or scoring an item that does not exist.

---

## 5. What this package does not decide

Unchanged and unauthorised: model population shape and slot substitutions; the Balanced Chat Baseline as a frozen profile; the presentation contract and override record; alpha thresholds, sample floor, NA policy; `G2`, `P-3`, `L4`; any corpus item, tag, split or score; any standard; any model run, rating session or pilot execution.

## 6. Filing

Written to `docs/research/corpus/`. No registry file exists in this repository that this record could update; row assignment is an owner/registry action, and no row number is claimed here.

## 7. Claim ceiling

Packaging only. No decision adopted, no rater named, no rating performed, no reply generated, no model loaded, no corpus or standard amended. Nothing here indicates that `MEMORY-09` improves behaviour, or that any of the four decisions is correct — only that each is now stated precisely enough to be decided.
