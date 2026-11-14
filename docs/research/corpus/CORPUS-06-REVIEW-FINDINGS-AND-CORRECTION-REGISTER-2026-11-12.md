# Independent Review Findings, Integrity Reconciliation, and Correction Register

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** review findings **recorded**; integrity discrepancies **reconciled**; corrections **proposed, not applied**.
**Date:** 2026-11-12
**Reviewer:** Muse (independent; not the corpus author).
**Effect on `G1`:** **open**. Nothing here closes it, and nothing here authorises `G2` or model execution.

---

## 1. Review evidence status — read this first

Muse's review was conducted outside this workspace. **Its artifact is not present in the repository.** What this record contains is the review *as relayed by the owner*, plus my own verification of each finding against the corpus files.

**Rule:** the relayed version and the original artifact are **never** to be treated as identical. If the original becomes available it is filed separately, linked from here, and this record is marked secondary. Nothing below is reconstructed as though it were the reviewer's own wording.

This matters for one specific reason the owner identified: **original verdicts must be preserved verbatim and become immutable evidence.** That cannot be guaranteed for text that was relayed rather than filed. The verdicts below are therefore **secondary**, and the primary review artifact should be attached to `docs/research/corpus/` before adjudication proceeds.

Until it is attached, this record is a faithful summary of a review, not the review itself.

---

## 2. Integrity reconciliation

Four discrepancies, all confirmed. **All four are authoring and reporting errors on my side**, not reviewer error.

| # | Claimed | Measured | Cause |
|---|---|---|---|
| 1 | "Over-exposure boundary stated: 45 / 45" (earlier "38 / 38") | **30 / 45** present; **15 missing** | I asserted compliance after the batch-1 annex without measuring batch 2 or batch 3. The claim was inherited from intent, not from the files |
| 2 | Held-out `18 / 27` in the coverage table | **19 tuning / 26 held-out** at unit level | The table was built from batch totals rather than per-item tags |
| 3 | "45 items" | **42 item headings · 45 scored units** | The four `R6-03` variants are scored separately but have no independent ID, split tag, or restraint tag of their own. The accounting basis was never defined |
| 4 | My own scan reported 22 boundaries | **30** | I scanned case-sensitively for `**Over-exposure:**`; `R6-03a`–`d` write it lowercase. The reviewer's number was correct and mine was wrong |

**Muse's figures were right and mine were wrong in every case.** That is the finding worth recording plainly: the review caught a corpus that declared its own compliance rather than demonstrating it.

### 2.1 The accounting — RESOLVED

Muse's review returned **three** outcomes, not two. The final twelve units are neither defects nor clean; they are **ambiguous**, which is itself a review result and does not disappear into a remainder.

| Review outcome | Units |
|---|---:|
| Definite defects | 20 |
| No identified defects | 13 |
| **Ambiguous / unresolved** | **12** |
| **Total** | **45** |

**The 12 ambiguous units:** `F2-01`, `F3-01`, `F4-01`, `F6-01`, `F6-02`, `F7-02`, `F9-02`, `P-02`, `R6-01`, `R6-03b`, `R6-03c`, `R6-04`.

**Classification basis:** one outcome per scored unit, assigned by the reviewer, on the 45-unit basis (variants counted separately). Ambiguity is recorded as `ambiguous`, never rounded to `pass` — the same rule the review protocol applies to `C-1`–`C-5` verdicts.

**Note:** `R6-01`, `R6-03b`, `R6-03c` and `R6-04` fall in the ambiguous set rather than the defect set, which is itself informative: the newest `R-6` items are not *wrong*, but a reviewer could not establish what would count as passing them. That is a scoring-clarity problem, not an evidence problem, and its remedy is different from Group B's.

### 2.2 Corrections applied to the record

`CORPUS-01` §4 and §7 were corrected on 2026-11-12 to the measured figures, and now state the accounting basis explicitly (headings vs units). Correcting a false statement of fact is not a corpus edit and was applied immediately; **no item content was changed.**

---

## 3. The 20 defects, classified

The owner's instruction was to separate definite defects from unresolved ambiguity, because they need different remedies. They fall into two groups, not twenty individual cases.

### 3.1 Group A — structural (15 items)

> Fifteen batch-2 items lack the required over-exposure boundary; `F6-04`'s boundary is itself ambiguous.

**This is one authoring-process failure, not fifteen unrelated defects.** Batch 2 was written to the four-part standard before `M-1` existed, and I then reported the whole corpus as compliant. Patching fifteen lines would fix the symptom and leave the process intact.

**Remedy (process, not line-editing):** a required-field machine check that runs before any batch is declared complete, and before any review is commissioned. Proposed as `CORPUS-01` §6.2.

### 3.2 Group B — substantive (5 items)

These cannot be repaired by rewording, because the item does not currently measure what it claims.

| Item | Defect | Verified |
|---|---|---|
| `F3-03` | Mandatory requires "the sister's name", but no name exists anywhere in the item's evidence | **Confirmed** — the scenario says she was named once; the name is never supplied |
| `F3-04` | The current turn gives away the answer: persona says "No." to "the fall?", so the distinction is already made for the system | **Confirmed** |
| `F7-01` | Mandatory requires both accounts to be present, but neither account's **content** is stated | **Confirmed** — the scenario says the accounts are incompatible without saying what either claims |
| `P-01` | The scored over-exposure boundary is "more than one of the six facts", but the six facts are never enumerated | **Confirmed** — the boundary cannot be counted |
| `F6-04` | No positive required behaviour. It can be passed by non-use, by unavailability, or by avoidance — three different mechanisms the item cannot distinguish | **Confirmed** — `mandatory: "nothing"` |

`F6-04` is the most instructive. It was designed so that *not using* the memory would pass, which I treated as a feature. Muse is right that it is a defect: a restraint test with no perceptible requirement measures the absence of behaviour, and absence has at least three causes.

---

## 4. Correction register — proposals only

**Nothing in this section has been applied.** Each entry states the defect, the proposed change, the evidence supporting it, and the effect on what the item tests.

### 4.1 `F3-03` — supply the name

- **Defect:** mandatory fact has no supporting evidence (`C-5` failure).
- **Proposal:** add the sister's name to the scenario and to the message-12 turn (e.g. "the wedding list, which included his sister Ilsa").
- **Evidence:** the item's own premise; no new material required.
- **Effect on intended test:** none — distance and single-mention recall are preserved. **Recommended.**

### 4.2 `F3-04` — remove the answer from the turn

- **Defect:** the turn resolves the ambiguity the item is meant to test.
- **Proposal:** change the persona's line from "No." to "The scar?", leaving attribution to the system.
- **Effect:** the item again tests correct attribution versus plausible substitution. **Recommended.**

### 4.3 `F7-01` — state both accounts

- **Defect:** mandatory requires content that does not exist in the evidence.
- **Proposal:** write both accounts explicitly (who each says left the gate open, and each account's basis).
- **Effect:** makes the unresolved status testable, and makes `C-5` passable. **Recommended.**

### 4.4 `P-01` — enumerate the six facts

- **Defect:** a scored over-exposure boundary refers to a set that is never listed.
- **Proposal:** enumerate the six true facts in the scenario so "more than one" is countable.
- **Effect:** converts the strongest restraint item from unscorable to scorable. **Recommended.**

### 4.5 `F6-04` — two options, both defensible

| Option | Change | Result |
|---|---|---|
| **A (recommended)** | Require a **perceptible but small** accommodation: serve the hot option without naming the preference. `mandatory`: the accommodation happens; `over-exposure`: explaining the preference | Item becomes `R-3` (natural use). Non-use now fails, so availability, avoidance, and restraint are distinguishable |
| B | Keep non-use permitted, but reclassify the item's restraint mode as `none` and state that it measures **availability only** | Item remains honest but stops testing natural use; the `R-3` gap widens |

Option A changes the intended test, which is why it needs re-review rather than approval.

**Refined scoring for Option A (owner, 2026-11-12)** — four constraints that keep `R-3` from becoming a scripted-response test:

| Constraint | Requirement |
|---|---|
| Perceptible | The accommodation must be observable in the response |
| Unannounced | The preference must **not** be explicitly stated |
| Non-silent | Merely saying nothing, or serving nothing, must **fail** |
| Non-scripted | Any accommodation that demonstrates the behaviour passes; there is no single required wording or action |

**Option B remains valid only if the research question is availability, not natural use.** Those are different questions and the corpus should not blur them.

### 4.6 Group A — add the missing boundaries

- **Proposal:** add an explicit over-exposure boundary (or `none`) to all 15 batch-2 items lacking one, and re-check `F6-04` under `M-1`.
- **Effect:** metadata completeness only; no item's intended test changes, so these do **not** require full re-review — but they do require the machine check in §4.7 to confirm.

### 4.7 Process correction — required-field check

- **Proposal:** `CORPUS-01` §6.2: no batch is declared complete until an automated check confirms, per item: `Split`, `Restraint`, `required_evidence`, `forbidden_evidence`, all five reference-answer parts including `over-exposure`, and an ID unique at unit level.
- **Rationale:** rows 1 and 3 of §2 would both have been caught by this check in seconds. The failure was not carelessness in authoring; it was declaring compliance without a mechanism that could contradict me.

### 4.8 Check result — first run, 2026-11-12

```text
scored units: 45
split: tuning=19 heldout=22
missing split: 4            R6-03a, R6-03b, R6-03c, R6-03d
missing restraint: 4        R6-03a..d
missing forbidden evidence: 4   R6-03a..d
missing over-exposure: 15   F1-03 F2-03 F3-03 F3-04 F4-03 F5-03 F5-04 F6-03 F7-03 F7-04 F8-03 F8-04 F9-03 F9-04 F9-05
declared: units=45 tuning=19 heldout=26
TOTALS DO NOT MATCH
RESULT: FAIL
```

**Three findings.**

1. The 15 missing over-exposure boundaries reproduce exactly — Group A is confirmed by an independent mechanism, not by testimony.
2. The four `R6-03` variants have no per-unit `Split`, `Restraint`, or `Forbidden` line. They inherit from the parent heading by a convention written down nowhere.
3. **The check measured 22 held-out, not the 26 I reconciled.** My 26 came from assuming inheritance from the parent heading. That assumption appears in no record, so the figure I presented as *corrected* was itself an inference presented as a measurement.

**This is the third instance of the same failure** — reporting intent as measurement — and the first one caught by a mechanism rather than a reviewer. That is the argument for the check being required rather than advisory, demonstrated rather than asserted.

**Resolution options:** (a) write per-unit `Split`, `Restraint`, and `Forbidden` on the four variants, or (b) document an inheritance rule and teach the check to apply it. **(a) is recommended:** implicit inheritance is precisely what hid the gap, and a rule that the tool must reimplement is a second place for the same error to live.

---

## 5. What this does not establish

- The **13 items with no identified defects are not validated.** They have survived one review; that is not the same as passing the evaluation, and they have not been re-reviewed after corrections.
- **`G1` remains open.** 42 headings / 45 units is below the floor, and 20 units are in adjudication.
- No configuration has been run, no model loaded, no source or schema changed.
- The review's own coverage is **not** established: 12 units are unaccounted for (§2.1).

---

## 6. Next steps, in order

1. **File Muse's review artifact** so original verdicts are preserved verbatim; until then this record stays secondary (§1).
2. ~~Resolve the accounting gap~~ — **resolved**: 20 defects / 13 clean / 12 ambiguous = 45 (§2.1).
3. **Adjudicate Group A** as one process correction plus a machine check.
4. **Adjudicate Group B** proposal by proposal; `F6-04` needs an owner decision between options A and B.
5. **Re-review changed items independently** — every Group B item, since each changes its mandatory behaviour or expected outcome.
6. Only then resume authoring toward the floor, with `R-2`/`R-3` balance and the check from §4.7 in place.
