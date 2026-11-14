# `CORPUS-97` — Completion Report: Registry and Identifier Reconciliation

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `CORPUS-97` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **COMPLETION REPORT — documentation only.**
**Work order:** `MEMORY-09-WO-REG-02` (2026-11-13) and its owner numbering amendment of the same date.
**Registry row:** **442**.
**This report establishes no substantive corpus validity, no reviewer independence and no convergence of any kind.**

---

## 1. Executive summary

Three previously unregistered documents are now registered (`CORPUS-26`, `CORPUS-32`, `CORPUS-93`), two pre-existing registry row collisions are documented in a separate cross-reference artifact (`CORPUS-96`), and the legacy identifiers `OD-9[T]` and `OD-10[T]` were investigated. **The legacy mapping could not be established and remains ambiguous**; no identifier was rewritten and no citation altered.

No corpus item, tag, split, score, threshold, checker line, standard or historical report was modified.

---

## 2. Starting state and deviations

| Item | Reported | Found | Outcome |
|---|---|---|---|
| Registry entries | 439 | **439** | match |
| Highest assigned row | (implied 439) | **437** | **deviation** |
| Next available row | 440 | **438** | **deviation — stop condition triggered** |

**Why the stop was correct:** 439 is the number of *entries*; 437 is the highest *number*. Rows 423 and 424 each carry two documents, so entries exceed the maximum by two. Assigning 440–444 as originally ordered would have permanently stranded rows 438 and 439.

**Resolution:** owner numbering amendment of 2026-11-13 assigned **438–442**. Verified free before use.

**Second deviation, from Deliverable C's premise:** the order stated that `OD-9[T]` and `OD-10[T]` each "appears twice in corpus documentation". They do appear twice each — but **only inside the owner decision record** (`MEMORY-09-ODR-1` L42 and L166). **No corpus document uses either bracket form.** Recorded in `CORPUS-96` §5.1 rather than converted into a mapping.

---

## 3. Registry assignments

| Row | Artifact | Bytes | Lines | SHA-256 |
|---:|---|---:|---:|---|
| 438 | `CORPUS-26` | 19,351 | 163 | `6d0242ef35c94dfcf57ef96902f91a2a9b1442afbe05bc43de67966b8a656081` |
| 439 | `CORPUS-32` | 24,287 | 223 | `28986a59fa43df5840177f77be9d7830a34bf37a164de1d586e630298286b750` |
| 440 | `CORPUS-93` | 46,151 | 539 | `9a135e54ad610865b312027c86b907e0136ef270e0a7f3d773d29092c060c91e` |
| 441 | `CORPUS-96` | 7,095 | 102 | `92ff437c3b1ba5813549d1ace16e0048f4b52105e2aefcce3743d37328183dc2` |
| 442 | `CORPUS-97` (this report) | see registry row 442 | see registry row 442 | see registry row 442 |

Metadata was recomputed from each file at write time, not copied from the earlier inspection.

---

## 4. Duplicate-row reconciliation

| Collision | Documents | Cross-reference |
|---|---|---|
| Row 423 | `CORPUS-82`, `CORPUS-83` | `CORPUS-96` §2 |
| Row 424 | `CORPUS-90`, `CORPUS-91` | `CORPUS-96` §3 |

- Both collisions are **row-number collisions between distinct documents**, not duplicate entries for one file.
- **No file was deleted, merged, overwritten or renumbered.** Rows 423 and 424 remain exactly as encountered.
- Future citations must name the **path and SHA-256**, not the row number alone.
- Unrelated registry defects: **none discovered** in the scope inspected; none corrected.

---

## 5. Legacy identifier results

| Citation | Occurrences | Mapping | Status |
|---|---:|---|---|
| `OD-9[T]` | 2 (both in `MEMORY-09-ODR-1`) | candidacy: `OD-M09-02` (gating) **and** `OD-M09-06` (vocabulary) | **Ambiguous** |
| `OD-10[T]` | 2 (both in `MEMORY-09-ODR-1`) | candidacy: `OD-M09-03` (cue standard) | **Ambiguous** |

**Unresolved, and deliberately so:** the `[T]` notation is defined nowhere in the corpus; `OD-9` spans two separate decisions; `OD-10` spans the cue standard and a broader five-layer framework; and `OD-M09-13`'s conditional permission to assign `OD-M09-nn` identifiers has not been exercised. **No mapping was recorded as confirmed and no citation was rewritten.**

---

## 6. Checker

```
scored units: 45
split: tuning=19 heldout=26
declared in CORPUS-01: {'units': 45, 'tuning': 19, 'heldout': 26}
TOTALS MATCH

UNVERIFIABLE (whole-corpus):
  semantic Required<->Variation / Over-exposure conflicts â€” not decidable from field text
  evidence sufficiency (C-5) â€” scenario prose is not machine-parseable
  proper-noun check (C-2) is WARN only; no warning is NOT a pass

RESULT: PASS (0 fail, 0 warn, 0 unverifiable)
This check reports only; no file was modified.
```

**Command:** `python docs/research/corpus/check_corpus.py` · **exit code:** 0

The checker reads corpus item files only. It cannot see registry rows, so it does not validate this work; it confirms that nothing in the item corpus changed. `TOTALS MATCH` reflects the held-out accounting corrected in `AM-02`.

---

## 7. Exclusions honoured

No change to: `CORPUS-01` or any corpus item; restraint tags, membership, splits, scores, thresholds or acceptance criteria; `MEMORY-09`, owner decisions, amendments or governing standards; `CORPUS-82`, `CORPUS-83`, `CORPUS-90` or `CORPUS-91`; checker code, validators, runtime, model or application; and no new review, adjudication, experiment, probe or benchmark was performed. No claim of reviewer independence or convergence is made anywhere in this work.

---

## 8. Files created or modified

| Action | File |
|---|---|
| Created | `docs/research/corpus/CORPUS-96-REGISTRY-ROW-COLLISION-CROSS-REFERENCE-2026-11-13.md` |
| Created | `docs/research/corpus/CORPUS-97-COMPLETION-REPORT-WO-REG-02-2026-11-13.md` |
| Modified | `docs/DOCUMENT-REGISTRY.md` — rows 438–442 appended, header count updated |

The three orphan documents, the four collided documents and every standard were **read only**.

---

## 9. Provenance

| Layer | Fact |
|---|---|
| Authority | Owner work order `MEMORY-09-WO-REG-02` + numbering amendment of 2026-11-13 |
| Executor | DeepSeek (per work order), repository worker — artifact-reported identity; labels are provenance data, not lineage proof |
| Model executed | No · `P-3` not started · no inference, benchmark or `G2` |
| Limit | This report validates bookkeeping, not substance. Registration does not make any registered document correct, independent or authoritative |
