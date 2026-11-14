# `MEMORY-09` — α Validation Closure (Task 4, `WO-ARCH-04`)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-AVC-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **GATE 3 CLOSED — arithmetic only; no corpus data used.**
**Authority:** `MEMORY-09-WO-ARCH-04` §4; gate introduced in row 454 §3.

---

## 1. Gate status

| # | Gate | Status | Result |
|---|---|---|---|
| 1 | Perfect agreement → α = 1 | **Pass** | α = 1 exactly, 45 units |
| 2 | Random ratings, expected behaviour across repeated samples | **Pass** | 400 samples/condition: means 0.0021 / −0.0054 / 0.0054; 90% intervals ≈ ±0.13–0.18 at n = 45 |
| 3 | **Reproduction of a published worked example** | **CLOSED — pass** | Two independent published examples reproduced exactly (§2, §3) |
| 4 | Two computation paths agree | **Pass** | Pair enumeration and coincidence matrix both give identical Do, De, α |
| 5 | Retained coincidence matrix, marginals, Do, De, α | **Required per run** | Procedure specified; nothing computed on corpus data |

## 2. Published example 1 — Krippendorff's canonical data

**Source:** Krippendorff, *Computing Krippendorff's Alpha-Reliability* (UPenn, `asc.upenn.edu/.../Computing Krippendorff's Alpha-Reliability.pdf`). The same table is reproduced in the K-Alpha Calculator paper (`pmc.ncbi.nlm.nih.gov/articles/PMC11636850`), which states **α = 0.743**, **40 pairable ratings**, 4 raters, 12 units, 7 missing, 5 nominal categories.

**Data as published** (rows = observers A–D, columns = units 1–12; `·` = missing):

| | u1 | u2 | u3 | u4 | u5 | u6 | u7 | u8 | u9 | u10 | u11 | u12 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| A | 1 | 2 | 3 | 3 | 2 | 1 | 4 | 1 | 2 | · | · | · |
| B | 1 | 2 | 3 | 3 | 2 | 2 | 4 | 1 | 2 | 5 | · | 3 |
| C | · | 3 | 3 | 3 | 2 | 3 | 4 | 2 | 2 | 5 | 1 | · |
| D | 1 | 2 | 3 | 3 | 2 | 4 | 4 | 1 | 2 | 5 | 1 | · |

**Computed:**

| Quantity | Value | Published |
|---|---|---|
| **n** (pairable values) | **40** | 40 ✓ |
| Marginals n₁…n₅ | **9, 13, 10, 5, 3** | 9, 13, 10, 5, 3 ✓ |
| Do | 1/5 | — |
| De | 152/195 | — |
| **α** | **113/152 = 0.74342…** | **0.743** ✓ |

Unit 12 (B's lone "3") is dropped, exactly as the source requires — the same rule that produced my earlier marginal error.

## 3. Published example 2 — ReCal worked example

**Source:** Freelon, *Worked Examples for Nominal Intercoder Reliability* (`dfreelon.org/recal/recal-worked-examples.pdf`), which publishes the coincidence matrix and **α = .85**.

| | Cat 0 | Cat 1 | Cat 2 | total |
|---|---|---|---|---|
| Cat 0 | 8 | 1 | 0 | 9 |
| Cat 1 | 1 | 6 | 0 | 7 |
| Cat 2 | 0 | 0 | 4 | 4 |
| total | 9 | 7 | 4 | **20** |

**Computed from the published matrix:** Do = 2/20, De = 8/9 · α = **108/127 = 0.85039…** vs published **.85** ✓.

## 4. What closing gate 3 does and does not do

**Does:** the implementation reproduces published α values exactly, on two datasets of different shape (4 raters × 12 units with 5 categories and heavy missingness; 2 raters × 10 units with 3 categories and one disagreement). Combined with gates 1, 2 and 4, the arithmetic is now **certified for use on corpus ratings**.

**Does not:** validate any study, item, rater, or construct. **Gate 3 certifies a calculator, not a result.** It says nothing about whether `MEMORY-09` units are ratable or whether α on real data would be high or low.

## 5. Residual requirement before any real α

The implementation used here exists only as an in-session script — **no file was written**. Before any real α is reported:

1. **File the implementation** as a reviewable artifact (separate authorisation; creating a script is outside a documentation-only order).
2. **Re-run all five gates from the filed version**, and record the outputs with the run.
3. **Retain** the coincidence matrix, marginals, Do, De, α, n, and the unit-drop list per run (gate 5).
4. **Report which published examples were used**, with URLs, so a reviewer can repeat the reproduction.

Until steps 1–2, the certification belongs to *this session's script*, not to a version anyone can inspect.

## 6. Why this gate mattered

Three earlier implementations produced α = 1.00, a marginal error, and α = −1.25 — each plausible-looking. Gate 3 is the only one that would have caught all three: it fails loudly when the arithmetic is wrong.

## 7. Non-claims

**No α was computed on corpus ratings.** No rating, session, model run or study occurred. Pasted above are synthetic and published data only; the certification is arithmetic.
