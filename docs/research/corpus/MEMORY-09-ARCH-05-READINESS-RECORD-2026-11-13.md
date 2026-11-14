# `MEMORY-09` — `WO-ARCH-05` Readiness and Authorization Record

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-READY-5` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **READINESS RECORD — preparation complete; execution unauthorised.**
**Authority:** `MEMORY-09-WO-ARCH-05` §6.

---

## 1. Four states

| State | Meaning |
|---|---|
| **A · Completed and verified** | Done, and checked mechanically |
| **B · Ready for owner approval** | Finished work awaiting a signature, not a fix |
| **C · Blocked** | Waiting on an authority decision or missing capability |
| **D · Unauthorised** | Executable in principle, permitted by nobody |

## 2. Readiness matrix

| Item | State | Evidence / blocker |
|---|---|---|
| **Alpha utility** | **A** | v1.0.1 filed; all six gates pass from the filed source |
| **Alpha validation** | **A** | Gates A–F, two published examples reproduced exactly, failures preserved |
| **F6-01 pilot packets** | **A** (frozen) | Rows 466–467; four packet hashes recorded |
| **Input equivalence** | **A** | Measured: 56 vs 57 chars (1.75%), salience scan clean, position and frame identical |
| **Contamination controls** | **A** | K1–K5 specified with detection and consequence; no control weakened |
| **Blinding** | **B** | Procedure final; whether it *works* is unknown until a run (and cannot be claimed in advance) |
| **Decision rules** | **B** | Five classes fixed; the 4-of-5 figures and repeat counts await approval |
| **Human-rater admissibility** | **C** | Unresolved by silence; clarification drafted, **not issued** |
| **`R-1` architecture** | **C** | Definition proposed (row 458); adoption requires a dated instrument |
| **Study execution** | **D** | No raters, no threshold, no authorisation |

## 3. The four separations this record exists to enforce

1. **Passing the α gate does not authorise study analysis.** The utility is certified to compute; it has never touched a `MEMORY-09` rating, and none exist.
2. **A frozen pilot does not authorise a pilot run.** Packets and hashes are *inputs awaiting permission*, not a run.
3. **A human-rater recommendation does not establish admissibility.** Row 468 requests a decision; the record is unamended.
4. **`R-1` architecture approval would not amend the standard.** Row 458 defines; `CORPUS-01` L121 remains operative until an instrument changes it.

## 4. What changed in this work order

| Before | After |
|---|---|
| α implementation existed only in session | Filed, versioned, independently gated, with a real defect found and fixed |
| Pilot had conditions but no frozen text | Frozen packets with measured equivalence and hashes |
| Human-rater question analysed | Decision request with exact wording and consequences |
| Readiness checklist (row 462) | Readiness matrix separating verified work from permission |

## 5. What did not change

No corpus item, tag, split, score, standard, threshold, sample floor, NA policy, checker line, schema, runtime or owner decision was altered. `check_corpus.py` was executed unchanged as an administrative check only.

## 6. The one open engineering item

**Gate D found a genuine defect** — the second computational path weighted disagreeing *units* rather than *pairs*, yielding `1087/1672` where the canonical fixture requires `113/152`. It was corrected, re-gated, and preserved in the validation record (row 465 §4). A single-path implementation would have shipped the right number for the wrong reason and failed silently on any four-coder fixture.

## 7. What a future execution phase would need

1. **Design approval** of the pilot package (§9 of row 467).
2. **Execution permission**, covering ~15 generated replies.
3. **A named rater** and agreement to report at single-source confidence.
4. For anything beyond compliance: **two independent raters** — which returns to the human-rater decision.
5. **A threshold**, only if a determination is to be made rather than described.

## 8. Summary

**Preparation is complete and verified; authority is entirely outstanding.** The α calculator is certified, the pilot is frozen, the decision request is written — and nothing has been executed, rated, generated or amended.
