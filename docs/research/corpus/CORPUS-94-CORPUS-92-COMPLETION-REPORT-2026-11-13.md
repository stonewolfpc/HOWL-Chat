# CORPUS-94 — Completion Report: `CORPUS-92` Internal-Consistency Correction

**Execution date:** 2026-11-13 (date carried by the `CORPUS-90`–`CORPUS-93` artifacts this report covers; the 2026-11-12 basis of the `CORPUS-84`–`CORPUS-89` package is untouched).
**Report type:** completion report for a narrow document-correction assignment. **Non-operative. Nothing adopted.**
**Assignment:** WORK ORDER — `CORPUS-92` FINAL INTERNAL-CONSISTENCY CORRECTION.
**Prepared by:** Buffy / Codebuff, in the Freebuff execution environment. **Certified model identity: none asserted** — no model name is exposed to this session and none is claimed. This is a repository-worker task, not a review of my own substantive work.
**Independence:** not claimed and not relevant to the assignment. No convergence credit arises from it.

---

## 1. Target and final state

**Path:** `docs/research/corpus/CORPUS-92-SOURCE-INTEGRITY-AND-CITATION-AUDIT-OF-CORPUS-91.md`

| | Before this pass | After this pass |
|---|---|---|
| Bytes | 33,724 | **38,840** |
| Lines | 469 | **511** |
| SHA-256 | `c76db062b6abe31fe02d650e69c55441d7eb64cb9bbf655d78fe52069ef98299` | `fad67830d026e080cb6dd6767f887fce812960e938c7abe737813aee84a407d9` |
| Registry row | 425 | **425** |
| Registry metadata | correct (`c76db062…`, 33,724, 469) | **correct** (`fad67830…`, 38,840, 511) |

**The work order's premise was already stale.** It reported a conflict between a "completion claim" of `f61da454…` / 33,793 / 469 and §9.2/§9.3 holding `701b710b…` / 30,328 / 444. On the saved file, **neither** set was in §9.2/§9.3: both carried `f61da454…` / 33,793 / 469, against an actual file of 33,724 bytes with digest `c76db062…`. The registry was already correct; only the document's self-description was wrong. I recomputed everything from the saved bytes rather than adopting either stated figure, as instructed.

---

## 2. Corrections applied — before / after

**A. Identity and metadata.** §9.2 previously asserted a hash, byte count and line count that were false for the file containing them. It now records the byte and line counts exactly as measured on the final save, and **defers the digest to the registry row**, with the reason stated: a file cannot contain its own final SHA-256, because inserting the digest changes the bytes it summarises. §9.3's quoted registry-row block — which reproduced the stale hash — is replaced by the row's identity fields plus a pointer to the registry for the digest. This removes the impossible claim rather than substituting a new stale one.

**B. Classification tally.** §10 already contained 22 findings classified 10 Confirmed / 6 Defective / 1 Partially supported / 5 Unable to verify. §11 asserted the same totals but as prose with "counted within #3 / #6 / #9" parentheticals that invited double counting. §11 is now an enumeration table naming the exact §10 finding numbers per class, with a stated reconciliation method (count §10 rows once each). **No classification was changed, and no finding was added or removed.** The "14 Confirmed" figure named in the work order is not present in the saved file; it appears to describe an earlier draft.

**C. Opening summary.** Added `§0.1 Correction pass 2 — what this document now claims`: what this pass corrected, what it deliberately did not do (repair `CORPUS-91`), what remains unresolved (#18–#22 plus the registry numbering defect), what was verified by execution versus merely asserted, and what is outside scope.

**D. Citations and cross-references.**

| Defect | Finding | Action |
|---|---|---|
| `CORPUS-91` cites the coverage line as **L138**; actual **L134** | #11 | **Verified again**: `CORPUS-91` L82 and L126 both say L138; `CORPUS-01` L134 is the coverage line. Recorded in §5.4. **Not corrected**, because the citation lives in `CORPUS-91` |
| `[GT]` used for the three-field-form inference in `CORPUS-85` §1 | #12 | Reported as-is. The label defect is in `CORPUS-91`'s text; `CORPUS-85` §1 was left untouched |
| "§8.5 of the original CORPUS-90 reading" ambiguous | #13 | **Resolved.** `CORPUS-90` §8.5 is "F-5 — §6 treats a proposed owner decision as adopted" (L437); its citation defects are in **§8.8** (L507). `CORPUS-90` L109 carries the same mis-point, which is the likely origin. Intended source: `CORPUS-90` **§8.8**. The defect remains in `CORPUS-91`; §8.1 item 4 now states the intended section |

**E. Registry duplicate numbering — disclosed, not repaired.** `DOCUMENT-REGISTRY.md` contains two rows numbered **423** (`CORPUS-82`, `CORPUS-83`) and two numbered **424** (`CORPUS-90`, `CORPUS-91`), from concurrent filings in other sessions. `CORPUS-92`'s row 425 is unique, but the *sequence* is not, so "row 424" names two artifacts. I did **not** renumber: that would rewrite other artifacts' cited row numbers and is outside this assignment. The defect is now disclosed in `CORPUS-92` §1.2 and §9.3, and in §3 below.

---

## 3. Unresolved items

1. **Registry numbering**: duplicate 423 and 424 pairs remain. Needs an owner decision — renumber with a correction note, or adopt non-numeric stable identifiers for future rows.
2. **`CORPUS-93` is unregistered.** `docs/research/corpus/CORPUS-93-INDEPENDENT-VERIFICATION-OF-MUSE-G1-FINDINGS-AND-REMEDIATION-PROPOSAL-2026-11-13.md` (46,151 B, 539 lines) has no registry row. Another session's artifact; not registered here because this assignment covers `CORPUS-92` only.
3. **`CORPUS-26` and `CORPUS-32` remain unregistered** (pre-existing; unchanged).
4. **Findings #18–#22 remain Unable to verify**: `CORPUS-82` §4D exact wording, `CORPUS-86` §5 exact wording, `CORPUS-78`/`CORPUS-79` line ranges, `MEMORY-09` R-13, `CORPUS-90`'s citations.
5. **Findings #7–#13 describe defects inside `CORPUS-91`** and remain uncorrected there by design. Fixing them requires `CORPUS-91`'s custodian, not this assignment.

---

## 4. Final §10 / §11 totals

| Class | Count | Findings |
|---|---|---|
| Confirmed | 10 | 1, 2, 3, 4, 5, 6, 14, 15, 16, 17 |
| Defective | 6 | 7, 8, 9, 11, 12, 13 |
| Partially supported | 1 | 10 |
| Unsupported | 0 | — |
| Unable to verify | 5 | 18, 19, 20, 21, 22 |
| **Total** | **22** | every §10 row exactly once |

---

## 5. Verification performed

1. Re-read the whole final `CORPUS-92` after editing (all nine replacements applied; no placeholders left; no stale digest asserted as current).
2. Recomputed SHA-256, bytes and lines from the saved file: `fad67830…`, 38,840, 511.
3. Refreshed registry row 425 and verified it matches those values exactly (path, hash, bytes, lines), with exactly one registry row for `CORPUS-92`.
4. Recomputed §10 counts and matched them to §11 by enumeration: 22 = 10 + 6 + 1 + 5 + 0.
5. Confirmed §0.1 describes the document as it now stands, including the unresolved list.
6. Ran `python docs/research/corpus/check_corpus.py` → `PASS (0 fail, 0 warn, 0 unverifiable)`, exit 0, `TOTALS MATCH`. **This is an administrative consistency check on item records only; it is not semantic validation of anything in `CORPUS-92`.**
7. Verified no governing standard, corpus item, tag, split, or score changed: no file outside `CORPUS-92`, the registry, and this report was written in this pass; `CORPUS-01`, `MEMORY-09`, `CORPUS-05`, `CORPUS-07`, the checker and all item files are untouched.

---

## 6. What this assignment did not do

- It did not amend `MEMORY-09`, `CORPUS-01`, or any governing standard.
- It did not change items, labels, tags, membership, splits, or scoring rules.
- It did not adopt any architecture or scoring recommendation.
- It did not reopen substantive adjudication recorded in `CORPUS-92`; no classification changed.
- It did not modify `CORPUS-91` or rewrite its history.
- It did not claim independent review, convergence, or validation. Convergence remains **0 of 2**; `G1` remains **open**.
- It did not make `CORPUS-92` an adopted standard. `CORPUS-92` remains a **non-operative, review-only** audit.

---

## 7. Status

| Area | State |
|---|---|
| `CORPUS-92` | Internally consistent; metadata verified against the saved artifact; non-operative |
| Registry row 425 | Matches the saved file exactly |
| Registry numbering | **Defective and disclosed** (duplicate 423, 424); unrepaired by design |
| `CORPUS-93` | **Unregistered** — flagged, not fixed |
| Standards · items · tags · splits · scores | **Unchanged** |
| Checker · runtime · inference | **Unchanged** · unchanged · none performed |
| Convergence | **0 of 2** |
| `G1` | **Open** |
