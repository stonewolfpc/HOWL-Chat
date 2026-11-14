# `MEMORY-09` — `OD-M09-17` Source Audit

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-OD17A-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **AUDIT — the earlier reconciliation is corrected here.**
**Authority:** owner instruction of 2026-11-13: audit `L88` and `L90` in context before treating the split as settled.
**Source:** `MEMORY-09-OWNER-DECISION-ADDENDUM-1-OD-M09-15-18-2026-11-13.md`, 163 lines.

---

## 1. Why this audit exists

Row 470 proposed a split: *judgement reproducibility* needs two raters; *behavioural observation* does not. That was presented as resolving `OD-M09-17`. The owner required the underlying text be examined directly, and was right to require it — **reading it in context changes the conclusion.**

## 2. The text, in context

| L | Text (verbatim) |
|---|---|
| 5 | "**Status:** **RECORDED — documentation direction only.**" |
| 9 | "**Effect:** this addendum records decisions. It authorises **no** item repair, scoring, evaluation, evidence collection, corpus change or implementation." |
| 15 | "These four decisions narrow and conclude four questions that `OD-M09-03` Version 2 left open in its §11, **and nothing more.**" |
| 25–27 | Register lists `OD-M09-17` as "Reproducibility — **TWO-RATER STUDY REQUIRED** before promotion beyond single-source confidence; not authorised for execution" |
| 88 | "**Disposition: TWO-RATER STUDY REQUIRED BEFORE PROMOTION BEYOND SINGLE-SOURCE CONFIDENCE.**" |
| 90 | "**Cue-level results** may be recorded and reported at **single-source confidence**. They may not support claims beyond it until a properly designed two-rater reproducibility study exists." |
| 118 | "It does not prohibit bounded **design** work on the protocol." |

## 3. What the context establishes

**L88 is a heading; L90 is its scope.** The disposition is stated absolutely, but the operative sentence restricts it to **cue-level results** — and §1 (L15) bounds the entire addendum to questions `OD-M09-03` §11 left open, i.e. **cue-criterion design**, not system behaviour.

**Three readings tested:**

| Reading | Claim | Verdict |
|---|---|---|
| **R1 · narrow** | The decision governs cue-level judgement reproducibility only; behavioural engineering observation is outside its subject matter | **Supported** — L90 names cue-level; L15 bounds the record to cue-design questions; L118 confirms design work is not prohibited |
| **R2 · broad** | Any claim beyond single-source confidence, of any kind, requires two raters | **Not supported** — rests on the L88 heading alone, and is contradicted by L90's explicit subject |
| **R3 · prohibitive** | The record forbids evaluation | **Not supported** — L9 records that the addendum *authorises no* evaluation. An absence of authorisation is not a prohibition |

## 4. The correction to row 470

Row 470 said engineering results "are single-source activities that must be labelled as such". **That implied the record permits them. It does not.**

Correct statement:

> **`OD-M09-17` does not reach behavioural engineering evaluation. It grants no permission for it, and imposes no two-rater requirement on it. Engineering evaluation therefore requires separate owner authorisation, and its results are neither cue-level nor covered by this decision.**

This is weaker in permission and stronger in accuracy. The distinction matters because "the standard allows it" and "the standard doesn't apply" lead to different records: the second means an engineering result can never be cited as satisfying `OD-M09-17`, and must never be described as single-source-confidence *cue* evidence.

## 5. A weakness in the boundary

**"Cue-level" is never defined in the addendum.** The scope boundary therefore rests on an undefined term. Its direction is not in doubt — L15 and L9 both point to cue-design — but a future reader could stretch "cue-level" to cover any observational claim.

**Minimum remedy:** an owner note stating the boundary, or a new instrument (`OD-M09-19`) defining "cue-level" and naming behavioural observation as out of scope. Neither is drafted here.

## 6. Consequence for planning

Until the owner rules, the Python-first plan must not rely on the split for authority. It should record engineering evaluation as **requiring its own authorisation**, and state plainly that `OD-M09-17` is silent on it rather than permissive.

## 7. Owner decision requested

| Option | Effect | Cost |
|---|---|---|
| **A (recommended)** | Adopt R1 as interpretation **and** authorise engineering evaluation by a separate decision | Keeps the record honest; one extra authorisation |
| **B** | Adopt R2 (broad) — every reported result requires two raters | Blocks all engineering reporting; no textual basis |
| **C** | Issue `OD-M09-19` defining "cue-level" and excluding behavioural observation | Cleanest long-term; one more instrument |

## 8. Claim ceiling

This audit establishes what the addendum's text does and does not cover. It does **not** decide the owner's interpretation, and it does not authorise any evaluation. It corrects one overstated inference in row 470 — an error of my own, found by reading the source rather than my summary of it.
