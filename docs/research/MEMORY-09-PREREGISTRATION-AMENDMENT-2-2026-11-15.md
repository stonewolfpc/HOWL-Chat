# MEMORY-09 — Preregistration Amendment 2: Capability Acceptance Rule

**Status:** proposed amendment, awaiting approval. Documentation and design only. No
implementation, no subject run, no new results, no frozen artifact changed, no commit.

**Append-only.** Amendment 1 (`97503340…`) and every earlier document are preserved unchanged.
This amendment adds the decision rule; it does not rewrite Amendment 1's metric definitions.

**Preserved without change:** `N` = 16 paired conversations, 3 generation reps per cell,
13 probes per conversation, seed matching and its enforcement, fixed-sample stopping with
outcome-blind halts only. **No genuine defect in those parameters was found**, so none is
altered.

---

## 1. Claim structure

Two distinct claims, each reported on its own, never fused into one word like "works":

| Claim | Requires | What it means |
|---|---|---|
| **C1 — net capability improvement** | **A1 > A0 on M1** | memory-equipped conversation beats no memory at all |
| **C2 — selection-specific benefit** | **A1 > D on M1** | the gain is attributable to *relevance*, at matched tokens and shape |

**Full capability claim = C1 ∧ C2.** Either alone is a weaker, separately named statement.

| C1 | C2 | Permitted statement |
|---|---|---|
| ✓ | ✓ | full capability claim: memory improved capability, and relevance drove it |
| ✓ | ✗ | capability improved, but **not** shown to be selection-specific (tokens/format could explain it) |
| ✗ | ✓ | relevance helps **relative to decoys**, but not enough to beat no-memory — report as such |
| ✗ | ✗ | no capability claim; if either comparison is significantly reversed, report a regression |

Both C1 and C2 are computed on **M1** only (Amendment 1 §2), never on M2.

---

## 2. Paired units, ties and zero differences

For conversation *i*, let `d_i = M1(A1, i) − M1(X, i)` where X ∈ {A0, D}. `d_i` ranges over
[−1, +1], since M1 is a proportion of 13 probes.

- **Tie:** `d_i = 0` exactly. **Ties are counted, reported, and excluded from the sign count.**
  They are never treated as favourable, never as unfavourable, and never discarded silently.
- **No tolerance band in the primary test.** `d_i = 0` is the only tie definition; a small
  positive difference is a positive sign even if it is tiny (magnitude is handled by the
  minimum-meaningful-difference rule in §3.3, not by redefining the sign).
- **Effective sample size:** `n_eff` = number of pairs where `d_i ≠ 0`.
- **Reported for every comparison:** all 16 `d_i` values, the tie count `16 − n_eff`, and `n_eff`.

**Floor for interpretability:** if `n_eff < 12`, the comparison is declared **indeterminate**
regardless of its p-value, because at that point a "significant" result would rest on fewer than
twelve independent conversations. This floor is fixed now, not after seeing ties.

---

## 3. Test, error budget and magnitude

### 3.1 Two claims, one family

C1 and C2 are two tests on the same 16 conversations. Their errors are correlated, but the
family of two is declared in advance:

- Family-wise α = **.05**, split across the two claims.
- Each claim is evaluated at **α = .025**, two-sided, exact sign test on `n_eff`.
- No third inferential test is added to the family. Secondary analyses are labelled descriptive.

### 3.2 Significance thresholds by `n_eff`

Exact two-sided binomial, α = .025, minimum concordant pairs required:

| `n_eff` | minimum favourable pairs | two-sided p at that point |
|---|---|---|
| 12 | 11 | .0063 |
| 13 | 11 | .0225 |
| 14 | 12 | .0129 |
| 15 | 13 | .0074 |
| 16 | 13 | .0213 |

The threshold is **recomputed from `n_eff`**, not fixed at 13/16: fixing it would silently make
tied conversations count against a claim.

### 3.3 Minimum meaningful difference (preregistered)

A significant sign test can rest on trivially small gains, so each claim additionally requires

**median `d_i` ≥ 0.10** (10 percentage points of M1, i.e. ≥1.3 probes per conversation).

Both conditions must hold. Report the median, the quartiles and all 16 values. This threshold is
fixed before data; a claim that is significant but below it is reported as **significant but
negligible**, which is not a capability claim.

### 3.4 Sensitivity analysis (secondary, disclosed)

As a declared **secondary** view, the same sign test is repeated counting only `|d_i| ≥ 0.10` as
sign-bearing pairs. This shows whether the result survives dropping near-zero pairs. It is
reported alongside and **cannot** be substituted for the primary test — choosing whichever view
looks better is exactly the freedom this document removes.

---

## 4. Invalid pairs, vetoes and exclusions

Two different things, with different consequences:

| Event | Consequence |
|---|---|
| **Veto fires** (cross-character leakage, erased history, fabricated provenance, temporal correctness — plan §4.3) | **the run fails**; no claim is computed and no aggregate is offered |
| **Instrument-invalid pair** (missing receipt; probe's supporting span still present in the window; `seed_honoured` false/unavailable; token accounting absent) | the pair is **excluded** from that comparison and counted |

Rules for exclusion, fixed now:

1. **Outcome-blind determination.** Invalidity is decided by the criteria above, applied without
   reference to answer correctness. A pair is never excluded because its difference was
   inconvenient.
2. **Reported, not hidden.** Excluded pairs are listed by id with the criterion that invalidated
   them, in the run receipt.
3. **Tolerance:** at most **2 of 16** pairs may be excluded per comparison. At 3 or more, that
   comparison is **indeterminate** — the run is not rescued by adding conversations
   (stopping rule forbids it).
4. **No replacement.** An invalid pair is not re-run with a different conversation; if the cause
   is instrument-level, the whole fixed sample restarts per the stopping rule.

If a pair is excluded, `n_eff` is computed from the remaining pairs; the §2 floor and the §3.2
threshold table then apply unchanged.

---

## 5. What is reported, in order

1. **Veto status.** If any veto fired: "run failed", stop.
2. **C1 and C2 separately**, each with: `n_eff`, tie count, all `d_i`, median/quartiles, p-value,
   threshold met or not, MMD met or not.
3. **Combined disposition** per the §1 table, stated in words, not a score.
4. **M2** as an absolute rate for A1, and for D as a manipulation check (expected ≈0; non-zero is
   an instrument failure). A0's M2 shown as `N/A`.
5. **Descriptives:** per-class ladder counts, token accounting, retrieval diagnostics.

**Prohibited:** a composite score; a single "capability improvement" number without both claims
named; reporting C1 alone as though it were the full claim.

---

## 6. Preserved parameters and one clarification

| Item | Status |
|---|---|
| `N` = 16 paired conversations | **unchanged** |
| 3 generation reps per cell, matched seeds | **unchanged** |
| 13 probes per conversation, quotas | **unchanged** |
| Fixed-sample stopping, outcome-blind halts only | **unchanged** |
| M1 / M2 definitions (Amendment 1) | **unchanged** |

**Clarification, not a change:** Amendment 1's phrasing "the capability claim requires A1 > D on
M1" is hereby refined to the two-claim structure above. A1 > D alone is a **selection-specific**
claim; the full capability claim additionally requires A1 > A0. No metric, parameter or
definition is altered.

**No defect found in `N`, reps, probes, seeds or stopping.** The only genuine issue this review
surfaced is the undeclared multiplicity of two tests, corrected by §3.1 without touching sample
size or resolution thresholds' meaning.

---

## 7. Proposal

1. Adopt C1 and C2 as separate, independently reported claims; full claim requires both.
2. Adopt tie handling: exact-zero ties excluded from the sign count, counted and reported,
   `n_eff` floor of 12.
3. Adopt per-claim α = .025 (family-wise .05) with thresholds recomputed from `n_eff`.
4. Adopt median `d_i` ≥ 0.10 as the minimum meaningful difference.
5. Adopt the exclusion rules: outcome-blind, ≤2 pairs, no replacement, veto fails the run.

**Awaiting approval.** G1 and G4 remain unauthorized.

---

## 8. Provenance

| Document | Status |
|---|---|
| `MEMORY-09-EXPERIMENT-PLAN-2026-11-15.md` | unchanged, `6c4078bd…` |
| `MEMORY-09-G0-PREFLIGHT-2026-11-15.md` | unchanged, `5c245283…` |
| `MEMORY-09-G0-RECONCILIATION-2026-11-15.md` | unchanged, `2b9c55b6…` |
| `MEMORY-09-PREREGISTRATION-2026-11-15.md` | unchanged, `8c788a96…` |
| `MEMORY-09-PREREGISTRATION-AMENDMENT-1-2026-11-15.md` | unchanged, `97503340…` |

This amendment is append-only and does not edit any of the above.
