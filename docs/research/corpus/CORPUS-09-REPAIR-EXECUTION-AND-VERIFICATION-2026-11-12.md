# Repair Execution and Verification

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** 20 of 20 authorised proposals **applied exactly**; repaired corpus **mechanically verified** (`PASS`); **awaiting independent re-review** of the five Group B items.
**Date:** 2026-11-12
**Authority:** owner repair authorisation, 2026-11-12.
**Effect:** `G1` remains **open**. A checker PASS is evidence of mechanical integrity, not of corpus validity or evaluation readiness.

---

## 1. Pre-repair snapshot

Held in `docs/research/corpus/snapshots/` — copies, not reconstructions.

| File | SHA-256 (pre-repair) |
|---|---|
| `CORPUS-02-BATCH-1-ITEMS` | `a18d02626c48…` |
| `CORPUS-03-BATCH-2-ITEMS` | `83b896257875…` |
| `CORPUS-04-R6-UNCERTAINTY-ITEMS` | `3d34f0530699…` |
| `check_corpus.py` | `791159844c27…` |

Full hashes: `snapshots/pre-repair-hashes.json`. Baseline output: `snapshots/BASELINE-CHECK-2026-11-12.txt`. `CORPUS-08` is unchanged and retains the proposals as approved.

**Pre-repair accounting** (`CORPUS-07`): 18 defect · 15 clean · 12 ambiguous = 45 units. Preserved unchanged.

---

## 2. Repair manifest — 20 changes

| # | Unit | Change | Applied |
|---:|---|---|---|
| 1–15 | `F1-03` `F2-03` `F3-03` `F3-04` `F4-03` `F5-03` `F5-04` `F6-03` `F7-03` `F7-04` `F8-03` `F8-04` `F9-03` `F9-04` `F9-05` | over-exposure boundary added (Group A) | **15 / 15** |
| 16 | `F3-03` | sister named **Ilsa** in scenario and mandatory | ✓ |
| 17 | `F3-04` | turn changed to "The scar?"; scenario notes nothing is confirmed | ✓ |
| 18 | `F7-01` | both accounts stated (Halla / Bren) | ✓ |
| 19 | `P-01` | six facts enumerated; boundary made countable | ✓ |
| 20 | `F6-04` | Option A: perceptible, unannounced, non-silent, non-scripted; `R-3` retained | ✓ |

**Traceability incident, recorded:** the first Group B application run **aborted** on a placeholder assertion after writing nothing; the second run applied all five. The corpus never contained a partial Group B state — verified by the checker, which still reported `C-1 F6-04` and both `U-1` findings between the two runs.

**Boundary-only check:** no Group A change altered an item's ask, evidence, mandatory behaviour, restraint mode, or difficulty. No item required individual stop-and-report.

---

## 3. Prediction versus result

Prediction was recorded in `CORPUS-08` §3 **before** application:

| Check | Predicted | Actual |
|---|---|---|
| missing over-exposure | 15 → **0** | **0** ✓ |
| `C-1 F6-04` | cleared | **cleared** ✓ |
| `U-1 P-01` | cleared | **cleared** ✓ |
| `U-1 F3-03` | cleared | **cleared** ✓ |
| totals | match | **TOTALS MATCH** ✓ |
| split | unchanged | **19 tuning / 26 held-out** ✓ |
| `RESULT` | `PASS (0, 0, 0)` | **`PASS (0 fail, 0 warn, 0 unverifiable)`** ✓ |

**Post-repair output:**

```text
scored units: 45
split: tuning=19 heldout=26
declared in CORPUS-01: {'units': 45, 'tuning': 19, 'heldout': 26}
TOTALS MATCH

UNVERIFIABLE (whole-corpus):
  semantic Required<->Variation / Over-exposure conflicts — not decidable from field text
  evidence sufficiency (C-5) — scenario prose is not machine-parseable
  proper-noun check (C-2) is WARN only; no warning is NOT a pass

RESULT: PASS (0 fail, 0 warn, 0 unverifiable)
```

Nothing was forced. No rule was weakened, no finding suppressed, no threshold moved.

---

## 4. Status classification

| State | Units / scope |
|---|---|
| **Repaired and applied** | 20 findings across 18 units |
| **Mechanically verified** | 45 / 45 units — `PASS` from `check_corpus.py`, plus 7-of-7 fixture discrimination (`CORPUS-10`) |
| **Awaiting independent re-review** | `F3-03`, `F3-04`, `F7-01`, `P-01`, `F6-04` — packet in `CORPUS-11` |
| **Unresolved** | 3 rule conflicts (`F3-01`, `F7-02`, `R6-01`); 9 other ambiguities with proposed pass conditions (`CORPUS-07` §2) |
| **Not applied** | `R-6` behavioural scoring model (`CORPUS-07` §8) — 5 `R-6` units still require narration-shaped pass conditions |
| **Requiring owner decision** | sealing; application of the 12 ambiguity proposals; `R-6` model adoption |

**Explicitly not claimed:** that the corpus is valid, sufficient, or ready. 42 headings / 45 units remains below the 80 floor; 13 units have survived exactly one review; and my own verification is **not** independent validation.

---

## 5. Boundaries observed

No new items. No split changes. No changes to `MEMORY-09`. No source, schema, runtime, or data change. No model loaded. No record deleted or silently retired. No sealing. `G1` open.
