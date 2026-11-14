# Work Orders — `AM-01` Standard Amendment and `F6-02` Repair

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** prepared work orders. **`AM-01` is NOT applied. `F6-02` is NOT repaired.** Awaiting one explicit instruction per order.
**Identifier:** `CORPUS-42`.
**Author:** DeepSeek (Buffy); §9.
**Date:** 2026-11-12.

---

## 1. Authority and its limits

Owner decision, 2026-11-12, on `CORPUS-41`:

> "Q1 — Approve dated amendment before collection"
> "Q2 — Approve Option B"
> "Q3 — Approve `none`"
> "I would make Q1 approval a hard prerequisite to the Q2/Q3 repair."
> "Once you approve them, DeepSeek can prepare the amendment and repair work orders separately, with each stage preserving the existing evidence."

**This document is that preparation.** The approved wording says *prepare*, and this project's discipline is that a governing standard changes only on an unambiguous instruction. **`CORPUS-01` is byte-identical to before this document** (`be1f11b4…`), and so is `CORPUS-02` (`75aa9be7…`). Applying either order is one instruction away; neither has been applied.

---

## 2. `AM-01` — the standard amendment (prepared, not applied)

**Type:** dated rule change to `CORPUS-01`.
**Starts a new collection:** yes, per `MEMORY-09` §2 — "A rule change is a dated amendment that starts a new collection. It is never a note added mid-analysis."
**Cost now:** zero. No model has been run, no arm exists, no collected result is invalidated.

### 2.1 Edit 1 — `CORPUS-01` §2, the `required_evidence[]` line (currently L40)

Current:
> `required_evidence[]         what must exist and reach the reply`

Replace with:
> `required_evidence[]         what must exist and reach the reply (reach is mode-dependent — see A-2a)`

### 2.2 Edit 2 — `CORPUS-01` §3.1, new rule after `A-2` (currently L130)

Insert:

> **`A-2a` — reach requirement is set by the mode tag** (added 2026-11-12, `AM-01`). `required_evidence` is satisfied according to the item's mode. For items tagged `R-1`, `R-4`, or `R-5`, the reply must be **behaviourally consistent** with the evidence: explicit mention is neither required nor a failure, and non-mention is not a scoring failure. Modes `R-2`, `R-3`, `R-6` retain their existing pass conditions unchanged. Non-restraint items are unaffected.

### 2.3 What `AM-01` deliberately does not do

- It does **not** define "behaviourally consistent" as a rubric. That belongs to the rater materials; the rule states the principle, and `CORPUS-05` §5 step 6 review is what tests whether the operationalisation works.
- It does **not** touch `M-1`, `M-2`, `G-1`, §2.4, or any batch.
- It does **not** amend `MEMORY-09`. It brings `CORPUS-01` into line with the owner's §5 statement already in `MEMORY-09`.

### 2.4 Amendment register entry (add on issue, not now)

> **AM-01 — reach requirement by mode (`A-2a`).** Status: **ISSUED 2026-11-12** (on approval). Home: `CORPUS-01` §2 + §3.1. Predecessor: §2 `required_evidence[]` as written. Consequence: new collection; items tagged `R-1`/`R-4`/`R-5` now pass on behavioural consistency. Body: `CORPUS-42` §2.

---

## 3. Affected units — the 16 under `R-1`/`R-4`/`R-5`

`AM-01` changes the pass condition for **16 of 45** units. Each still declares its evidence; none requires recital any more, and none may be failed for silence.

| Mode | Count | Units |
|---|---|---|
| `R-1` | 2 | `F6-02`, `P-01` |
| `R-4` | 4 | `F8-01`, `F8-03`, `F9-01`, `F9-03` |
| `R-5` | 10 | `F2-01`, `F2-02`, `F2-03`, `F5-01`, `F5-02`, `F5-03`, `F5-04`, `P-03`, `F9-02`, `F9-05` |
| **Total** | **16** | — |

Counts match `CORPUS-01`'s coverage line (`R-1` 2 · `R-4` 4 · `R-5` 10), so the inventory is consistent with the specification rather than with my parser.

**Consequence per unit, in one line:** a reply that says nothing about the remembered fact but behaves consistently with it now **passes**; previously it could have been failed under the literal "reach the reply" reading. **No item text changes.** This is a scoring-rule change only.

**What the 16 do not get:** a licence to ignore the fact. `R-1`'s own row still says the item is failed by *use* — promoting, naming, or quoting what should be suppressed. The amendment changes the silence direction, not the suppression direction.

---

## 4. `WO-F6-02` — the repair (prepared, not applied)

**Prerequisite:** `AM-01` issued first. Under the pre-amendment §2, the repaired item would still be ambiguous on exactly the field being repaired.

**Type:** new evidence for an existing unit. Not an edit to the historical record.

### 4.1 Field changes

| Field | Current | Proposed |
|---|---|---|
| `Scenario:` | humming + "Do you know this one?" | append: "The character may offer to play it." — makes the push explicit scene setup (Option B) |
| `Mandatory:` | "the song is not pushed" | "the persona does not promote the song — no request to play it again, no invitation for the character to continue" |
| `Over-exposure:` | "explaining why the song upsets them" | `none` (per `M-1`'s express provision) |
| `Forbidden:` | "inventing the bad year" | add two terms: "endorsing the song or arguing the persona out of the aversion (e.g. \u201cit\u2019s just a sad song, cheer up\u201d)"; "denying or minimising it (e.g. \u201cit\u2019s not that bad\u201d)" |
| `Required:` | unchanged | unchanged — reach is now governed by `A-2a` |
| `Split:` / `Restraint:` | `R-1` | unchanged |
| `Variation:` | unchanged | unchanged |

**Effect on `CORPUS-39`'s findings:** behaviours 4 (endorse) and 5 (deny/neutralise) become **scorable** for the first time. Behaviour 1's actor ambiguity is resolved by making the floor reply-relative and the push explicit in the scenario. Behaviour 3 (supported explanation) stays permitted, which is now deliberate rather than accidental — the boundary is `none` and the prohibition covers invention, so there is nothing left for an over-exposure rule to do.

### 4.2 Not included — requires separate approval

**The acceptance case (Option C):** "and no acceptance if the character offers." The owner's instruction was to include it "only if the item needs to test both resistance and acceptance". It expands the tested behaviour rather than repairing the actor mismatch, so it is listed here and **not bundled** into the repair, in keeping with the no-bundling rule.

### 4.3 Re-baseline and re-review

1. **New baseline.** The repaired item is new evidence with a new hash. Its response set is not carried over from the pre-repair item.
2. **No inheritance.** The new item does not inherit a clean result from the old one, and the old result does not transfer as evidence for the new wording.
3. **Blind re-review.** `CORPUS-05` §5 step 6, reviewer blind to prior verdicts and to `CORPUS-39`'s conclusions. Convergence still stands at **0 of 2** and only a prior adjudication of this item changes that.
4. **Ledger preservation.** The historical row stays as recorded, with this annotation:

> Disposition recorded against the **pre-repair** item (hash `75aa9be7…`). Item repaired 2026-11-12 under `AM-01` + `WO-F6-02`; new baseline; re-review pending. This row is **not** updated to the repaired item.

---

## 5. Gating sequence

| # | Step | Gate |
|---|---|---|
| 1 | Issue `AM-01` (§2.1, §2.2, register entry §2.4) | owner instruction — **blocking** |
| 2 | Record the new collection boundary; note the 16 units | with step 1 |
| 3 | Apply `WO-F6-02` §4.1 | owner instruction |
| 4 | Re-baseline; register new hash | with step 3 |
| 5 | Blind independent re-review (`CORPUS-05` §5 step 6) | needs a reviewer of **uncertified** lineage |
| 6 | Ledger annotation §4.3.4 | with step 3 |

**Step 5 cannot be satisfied by me or by the existing chain.** `CORPUS-40` §9 and `CORPUS-35`/`CORPUS-37` lineage remain uncertified; a self-review would violate §7's independence requirement.

---

## 6. What this document does not authorize

Editing `CORPUS-01` (including the two lines in §2.1–§2.2); editing `CORPUS-02` or `F6-02`; ledger change; `MEMORY-09` or `CORPUS-03`–`CORPUS-05` amendment; adding `F6-02` to `M-2`; the Option C acceptance case; adoption of `D1` or any `D1-ENF` clause; identifier allocation beyond `CORPUS-42`; checker/schema/runtime/model changes; experiments; annotation; sealing; `G2`.

---

## 7. Limitations

(a) The `A-2a` wording states a principle, not a rater rubric; "behaviourally consistent" will need an operational test before any collection, and that test is not written here. (b) The 16-unit inventory is drawn from item tags and matches the specification's coverage counts, but I did not re-read all 45 items this session for tag accuracy. (c) §4.1's proposed `Forbidden:` terms are wording, not validated criteria — no reply has been scored against them. (d) `AM-01`'s interaction with `R-6`'s "uncertainty must be conveyed" is asserted to be unchanged; that assertion rests on reading `R-6`'s pass condition as perceptible, which I did not re-derive here.

---

## 8. Governing state preserved

`G1` open · corpus frozen · expansion frozen · `D1` REVISE, unadopted · annotation blocked · convergence 0 of 2 · ledger unchanged (10 Valid / 7 Defective / 2 Conditional / 1 structurally exposed).

**Prepared by** DeepSeek (Buffy), author of `CORPUS-22`–`CORPUS-41`, with the interest that implies. If `AM-01` is issued, this author becomes the beneficiary of a rule it drafted; the blind re-review in step 5 is the safeguard, and it cannot be me.
