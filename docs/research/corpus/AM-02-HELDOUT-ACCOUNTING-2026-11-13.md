# `AM-02` — Held-Out Accounting and the `F-9` Row

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `AM-02` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **EFFECTIVE — documentation correction only.**
**Authority:** owner authorization, 2026-11-13, on review of the read-only investigation of the `19/22` versus `26` held-out question.
**Amends:** `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` — §4 family table (`F-9` row) and the superseded "Open discrepancy" note beneath it.
**Effect on items:** **none.** No item text, split value, tag, score, threshold, checker line or runtime behaviour changed.

---

## 1. Two separate findings, two separate corrections

The investigation distinguished an aggregate that was already resolved from a transcription error that was not.

| # | Finding | Kind | Correction |
|---|---|---|---|
| 1 | The checker measures 19 tuning / 26 held-out, matching the declared totals | Documentation **staleness** — the `22` figure predates the variants gaining their own split lines | Replace the note, preserve it verbatim as superseded |
| 2 | `F-9`'s family row reads `2 tuning / 3 held-out`; the units are `3 / 2` | **Transposed pair** in one row | Correct the row |

**26 − 22 = 4**, exactly the four `R6-03` variants, each of which now carries `**Split:** heldout`.

---

## 2. Evidence

**Checker output (2026-11-13):**

```
scored units: 45
split: tuning=19 heldout=26
declared in CORPUS-01: {'units': 45, 'tuning': 19, 'heldout': 26}
TOTALS MATCH
```

**Unit rule (read from `check_corpus.py` L76–110):** a unit is a heading containing a `**Reference:**` line; split is read from that heading's own `**Split:**` line, case-insensitively matched against `held-?out|tuning`. Declared totals are parsed from the `Total` row and compared at L199–202 — which is why a transposed *family* row could never have been caught by the checker.

**`F-9` units, individually:**

| Unit | Split |
|---|---|
| `F9-01` | tuning |
| `F9-02` | tuning |
| `F9-03` | held-out |
| `F9-04` | tuning |
| `F9-05` | held-out |
| **Total** | **3 tuning / 2 held-out** |

**Family consistency after the correction:** every family row matches its units exactly, and the rows sum to `45 / 19 / 26` — the independently verified totals.

---

## 3. Basis disclosure

Per-family figures were re-derived by **assigning units to families from their item-ID prefixes** (`F7-*` and `R6-*` → `F-7`), because batch-1 items carry no `Family:` field. This is an accounting convention, not a metadata read of every item, and it is stated in `CORPUS-01` alongside the correction rather than left implicit.

The `Total` row did not depend on that convention: it was already correct and is **unchanged**.

---

## 4. The superseded note, preserved verbatim

> **Open discrepancy:** `check_corpus.py` measures **19 tuning / 22 held-out** — the four `R6-03` variants carry no per-unit split line. The `26` above assumes inheritance from the parent heading, which is written nowhere. Not reproducible → not compliant. See `CORPUS-06` §4.8.

Its diagnosis was accurate when written: at that time the variants had no per-unit split line, and the `26` did depend on an unwritten inheritance assumption. The corpus later gained explicit per-unit split lines, which removed the assumption and made the figure 26. The note was not updated. It is retained in `CORPUS-01` as a quoted superseded record, and here in full, so the earlier claim remains citable and its resolution traceable.

---

## 5. What did not change

- No split value was altered. `F9-03` and `F9-05` remain held-out; `F9-01`, `F9-02`, `F9-04` remain tuning.
- No item, tag, membership, score, threshold, split declaration or checker line changed.
- `Held-out: 26 / 45 (58%)` remains correct as stated.
- The "≥2 held-out items per family" claim remains satisfied — `F-9` has 2.
- No other §7 ratio was verified or altered; that is a separate investigation and was not authorised here.

---

## 6. Verification

| Check | Result |
|---|---|
| `check_corpus.py` | `PASS`, `TOTALS MATCH`, exit 0 — administrative only |
| Family rows sum | 45 / 19 / 26 after the correction |
| Registry metadata | `CORPUS-01` row refreshed to the amended file's bytes, lines and SHA-256 |
| Items, tags, splits, thresholds, checker | unchanged |

---

## 7. Provenance

| Layer | Fact |
|---|---|
| Authorization | Owner, 2026-11-13 |
| Discovery | Read-only investigation of the `19/22` versus `26` question; the `F-9` transposition was found while verifying family rows |
| Filed by | DeepSeek (per work order), repository worker — artifact-reported identity; labels are provenance data, not lineage proof |
| Model executed | No · `P-3` not started · no inference, benchmark or `G2` |
| Registry | Append-only; see `docs/DOCUMENT-REGISTRY.md` |
