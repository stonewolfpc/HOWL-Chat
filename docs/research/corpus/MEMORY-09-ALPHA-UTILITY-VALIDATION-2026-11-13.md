# `MEMORY-09` — Alpha Utility Validation Record

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-AUV-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **VALIDATED — v1.0.1 passes the full gate suite, from the filed source.**
**Authority:** `MEMORY-09-WO-ARCH-05` §3.
**Files under test:** `MEMORY-09-ALPHA-UTILITY-2026-11-13.py` (v1.0.1), `MEMORY-09-ALPHA-UTILITY-FIXTURES-2026-11-13.py`.

---

## 1. Reconstruction disclosure (required before any result)

The implementation used earlier in session ran from an **unnamed heredoc and was never written to disk**. It **cannot be recovered exactly.**

This file set is therefore a **reconstruction**, built to the same specification and validated from scratch. No claim is made that it is byte- or logic-identical to the in-session version. What is claimed: **the filed version passes every gate**, and gates 1–5 of the earlier record are reproduced (§4).

## 2. Specification as filed

| Property | Value |
|---|---|
| Statistic | Nominal Krippendorff's α = 1 − Do/De |
| Missing data | `None` = missing; units with fewer than two values **dropped**, as the published method requires |
| n | Pairable values only (dropped units excluded from marginals) |
| Coincidence | o_ck = Σ_u n_uc·n_uk/(m−1), diagonal o_cc = Σ_u n_uc(n_uc−1)/(m−1) |
| Arithmetic | exact `Fraction`, so no float drift |
| Paths | `coincidence` (matrix) and `pairs` (unordered pair enumeration) |
| NA policy | **not invented**: `na_policy="missing"` (default) or `"category"` with an explicit label; an undeclared value raises |
| Degenerate De = 0 | returns `alpha = None` with a `reason`; never coerces |
| I/O | no network, no corpus access, deterministic |

## 3. Gate results — final run, filed source

Full suite executed via `python docs/research/corpus/MEMORY-09-ALPHA-UTILITY-FIXTURES-2026-11-13.py`, exit 0.

| Gate | Result | Evidence |
|---|---|---|
| **A** Perfect agreement | **PASS** | α = 1 exactly (45 units) |
| **B** Random simulations | **PASS** | 400 reps × 3 conditions; means 0.0039 / 0.0052 / 0.0070; p05–p95 −0.121…0.146, −0.124…0.155, −0.140…0.173 |
| **C1** Krippendorff canonical | **PASS** | α = **113/152 = 0.74342** (published .743); n = **40**; marginals **9, 13, 10, 5, 3** |
| **C2** ReCal matrix | **PASS** | α = **108/127 = 0.85039** (published .85); n = **20** |
| **D** Two paths agree | **PASS** | canonical 113/152 both paths; perfect 1 both paths |
| **E** Intermediates inspectable | **PASS** | Do = 1/5, De = 152/195, dropped unit index 11, diagonal 7/10/8/4/3 |
| **F** Negative controls | **PASS** | insufficient data, undeclared category, undefined NA policy all rejected; De = 0 flagged not invented; valid 3-value unit accepted |

## 4. Failures preserved (not concealed)

### Gate D — a real defect in the second path

**First run: FAIL.** Coincidence gave `113/152`; pairs gave **`1087/1672` ≈ 0.6501**.

**Cause:** the pair path computed Do as *disagreeing units ÷ units* — correct only when every unit has exactly two coders. With four coders a unit holds six unordered pairs, so each disagreeing pair must carry weight **2/(m−1)**, summed over n padded values.

**Fix:** weighting corrected; version 1.0.0 → **1.0.1**; both paths now `113/152`.

**Why this matters:** two implementations agreeing is not evidence if one is wrong in the same direction — but here they disagreed, which is exactly what gate D exists to surface. A single-path implementation would have shipped the correct answer for the wrong reason and failed the first four-coder fixture silently.

### Gate F — two faulty expectations (fixture defect, not utility defect)

| Expectation as first written | Truth |
|---|---|
| a single unit holding three values should be rejected | **False** — it is a valid α input with n = 3. The insufficient-data case needs units that all drop out (`[[1],[2]]`) |
| De = 0 should raise | **Wrong design** — it is a defined degeneracy; returning `alpha = None` with a reason is what "flag, never coerce" requires |

Both were corrected in the fixtures; the utility was not changed for them.

## 5. What is certified, precisely

**Certified:** the filed v1.0.1 computes nominal Krippendorff's α correctly on two published fixtures of different shape, reproduces perfect agreement, behaves as expected across 1,200 random samples, agrees across two code paths, exposes every intermediate, and refuses invalid input.

**Not certified:** that α is the right statistic for any particular study design; that any threshold is defensible; any NA policy beyond the published method; **and nothing at all about `MEMORY-09` ratings**, none of which exist.

## 6. Standing restrictions

No real corpus ratings were used. No α was computed on study data. No threshold, sample floor, or NA policy was adopted. The utility must not be pointed at live corpus ratings without separate authorisation.
