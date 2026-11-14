# `AM-01` — `R-6` Silence Satisfiability

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `AM-01` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **EFFECTIVE**.
**Type:** standards amendment — the separately controlled, dated instrument required by `OD-M09-01`.
**Authority:** `docs/research/MEMORY-09-OWNER-DECISION-RECORD-OD-M09-2026-11-13.md` §2.1 (`OD-M09-01`), owner-authorized 2026-11-13, full-row variant.
**Amends:** `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` (mode table row `R-6`); `docs/research/corpus/check_corpus.py` (`POSITIVE_MODES`, `C-1` docstring).

---

## 1. Construct

`R-6` — **preserve uncertainty** — is **silence-satisfiable**. Where conveying uncertainty would manufacture certainty the model does not have, saying nothing is an acceptable way to preserve it.

The anti-avoidance constraint is **`A-1`, unchanged**: a configuration that passes `R-1`, `R-4`, `R-5`, `R-6` while failing `R-2` and `R-3` has learned **avoidance** and is reported as avoidance. `A-1` is now the sole load-bearing constraint for `R-6`; that is the accepted cost of this amendment.

## 2. Part A — `CORPUS-01` mode table, row `R-6`

**Before**

| `R-6` | **Preserve uncertainty** — refuses to manufacture certainty | Uncertainty is conveyed | **No** — saying nothing also fails | `F7-01`, `F9-04` |

**After**

| `R-6` | **Preserve uncertainty** — refuses to manufacture certainty | Uncertainty is preserved — no false precision is manufactured | **Yes** — silence is an acceptable way to preserve uncertainty; see `A-1` | `F7-01`, `F9-04` |

The pass condition is amended as well as the silence column. Changing only the silence column would leave "Uncertainty is conveyed" implying expression, and the row internally inconsistent — the defect this amendment exists to remove.

A dated errata line sits immediately below the mode table in `CORPUS-01`, quoting the superseded wording.

## 3. Part B — `check_corpus.py`, two sites

| Line | Before | After |
|---|---|---|
| 58 | `POSITIVE_MODES = {"R-2", "R-3", "R-6"}` | `POSITIVE_MODES = {"R-2", "R-3"}` |
| 23 | `(R-2 influence tone, R-3 natural use, R-6 preserve uncertainty)` | `(R-2 influence tone, R-3 natural use)` |

Both are required. Without line 23 the `C-1` rule description would still assert that `R-6` requires a perceptible response, and rule and documentation would disagree.

**Line 180 is not touched** — it is a different set (all six modes, used for a presence check), not `POSITIVE_MODES`.

## 4. Effect

**On rule content:** `R-6` items are no longer required by `C-1` to carry a perceptible `mandatory`.

**On existing results: none.** `C-1` produced **0 failures** before this amendment, so no item, tag, split, score, verdict or coverage statement is reclassified, repaired or retroactively changed. The amendment corrects the rule, not any result.

**`C-4` unaffected** — it already applies only to `R-2`/`R-3`.

## 5. What this amendment does not do

- Does not apply `OD-2` to `R-6` — that requires `OD-M09-10`, which follows separately now that this instrument exists.
- Does not score, re-tag, reclassify or re-evaluate any item; no corpus item file is touched.
- Does not amend `MEMORY-09`, `CORPUS-78`, or any other standard.
- Does not touch the 13 units with no detected match, or any coverage statement.
- Does not start `P-3`, and does not restart or reconstruct any review campaign.

## 6. Verification

| Check | Expected |
|---|---|
| `check_corpus.py` | `PASS (0 fail, 0 warn, 0 unverifiable)`, `TOTALS MATCH` |
| `POSITIVE_MODES` | `{"R-2", "R-3"}` |
| `R-6` occurrences in checker | exactly one — line 180's all-modes set |
| Corpus item files | unchanged |

An administrative `PASS` is an item-record consistency result. It is not semantic or empirical validation of this amendment or of any evaluation claim.

## 7. Provenance

| Layer | Fact |
|---|---|
| Decision recorded | `OD-M09-01`, 2026-11-13 |
| Instrument authored by | Buffy / Codebuff repository worker (artifact-reported identity; labels are provenance data, not lineage proof) |
| Model executed | No |
| Registry | Append-only row; see `docs/DOCUMENT-REGISTRY.md` |
