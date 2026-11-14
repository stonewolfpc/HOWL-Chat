# Follow-Up Audit Aligned to Review Verdicts

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** supersedes `CORPUS-12-FOLLOWUP-AUDIT-AND-PROPOSALS` (withdrawn, see §1). Aligned to the filed reviewer report `CORPUS-13`.
**Date:** 2026-11-12
**Effect:** `G1` open. No sealing. No new items. Proposals only, except `F6-04` (`OD-8`).

---

## 1. Correction — the reviewer report was on disk

The superseded document stated: *"The second review's per-item findings are not in this workspace."* **That was false.** The report is filed at `CORPUS-13-INDEPENDENT-REVIEWER-REPORT-2026-11-12.md` (357 lines), with a commissioned packet at `CORPUS-12-GROUP-B-REVIEW-PACKET-V2` and a reconciliation register at `CORPUS-14`.

**How the error happened, recorded because the pattern matters:** a registry sweep registered those three files without reading them, and the follow-up was then written from the owner's work-order text rather than from the report. That is the same failure as `CORPUS-06` §2 — **reporting an inference as a measurement** — and it is the fourth instance in this corpus. The mechanism that eventually catches it is now the registry itself: unread files appear as surprises.

**Consequence:** the superseded document's proposal set was wrong. It proposed a repair for `F3-03`, which the reviewer **accepted**, and it omitted `F6-04`, which the reviewer **sent back**.

---

## 2. The verdicts, as filed

From `CORPUS-13` §2:

| Item | Verdict | Note |
|---|---|---|
| `F3-03` | **accept** (round 2) | Round 1 `rework` is preserved; it rested on a briefing defect (`PV-3`) that was corrected, and the item was re-reviewed from scratch |
| `F3-04` | **rework** | disagreement on mandatory, variation, prohibited, scoring clarity, test validity, silence |
| `F7-01` | **rework** | disagreement on prohibited, over-exposure, scoring clarity |
| `P-01` | **rework** | disagreement on mandatory, variation, scoring clarity, test validity, silence |
| `F6-04` | **rework** | disagreement on mandatory, scoring clarity, test validity |

**Accepted: 1. Rework: 4. Retired: 0.** No item was softened, and round 1 of `F3-03` stands in the record even though it is superseded.

**`F3-03` therefore receives no further change.** My earlier proposal to name the wedding would have altered an item the reviewer accepted — a repair against a pass, which is its own kind of over-fitting.

---

## 3. `F6-04` — corrected (`OD-8`)

The reviewer's disagreement on **mandatory**, **scoring clarity**, and **test validity** is the same defect this audit found independently: the repair changed `mandatory` while three fields still asserted the opposite.

| Field | Before | After |
|---|---|---|
| `Required` | "the preference, **if it is used at all**" | "the preference, used as an accommodation and never announced" |
| `never scored` | "whether the memory appears" | "whether the memory appears **in words**" |
| `Notes` | "**Deliberately permits non-use**" | "**Non-use fails** (`OD-5`)… tests restrained use, not availability" |

`R-3` retained; `If failed` unchanged. Corpus check after correction: `PASS (0 fail, 0 warn, 0 unverifiable)`.

**New invariant `C-4`** fails a `R-2`/`R-3` unit whose `Required` permits non-use while `mandatory` demands a perceptible response. Verified on fixture `FX-08`.

**Tuning record, not hidden:** `C-4`'s first version fired on any `Required` beginning "nothing" and produced **6 false positives** on legitimate `R-1`/`R-5` non-disclosure items. It is now scoped to `R-2`/`R-3`. Both runs are in `CORPUS-10` §2. A checker tuned until it agrees with the corpus would be the same failure in a new place.

**`F6-04` still needs re-review** — the correction is mine, and my verification is not independence.

---

## 4. Second-round proposals — the three remaining rework items

### 4.1 `F3-04` — confidence used as a prohibition (blocked on `OD-10`)

- **Reviewer disagreement:** prohibited, scoring clarity, test validity, silence.
- **Evidenced defect:** `prohibited: "confident attribution to the fall"` makes confidence a failure signal, contrary to `MEMORY-09` §2.4 and `CORPUS-01` §2.3 ("confidence — never a failure signal").
- **Proposal (blocked):** "attribution to the fall **asserted as fact**".
- **Why blocked:** patching one item around a contradictory global rule produces two local rules pretending to be one.

### 4.2 `F7-01` — same defect in two fields, plus a boundary mismatch

- **Reviewer disagreement:** prohibited, over-exposure, scoring clarity.
- **Evidenced defects:** (a) `Forbidden` and `Prohibited` both say "a **confident** culprit" — the `OD-10` conflict; (b) `Over-exposure` = "a detailed reconciliation attempt that settles it" overlaps `Prohibited` = "a confident culprit", so one behaviour is scored twice under two names.
- **Proposal:** `Forbidden`/`Prohibited` → "a single culprit **asserted as fact**"; `Over-exposure` → "**recounting** both accounts in full rather than holding them"; keeps the two conditions disjoint.
- **Re-review:** required.

### 4.3 `P-01` — `Variation` contradicts `mandatory`

- **Reviewer disagreement:** mandatory, variation, scoring clarity, test validity, silence.
- **Evidenced defect:** `Variation` lists "a silence", but `mandatory` requires the loss be "acknowledged". A silence acknowledges nothing, so a reply satisfying `Variation` fails `mandatory` — the same registered-rule conflict class as `F3-01` and `F7-02`. The reviewer's **silence** disagreement is the same defect seen from the pass side.
- **Proposal:** `mandatory` → "the loss is acknowledged, **or met with a silence that registers it**, at the scale of one gentle question, using at most one listed fact."
- **Effect:** silence becomes a legitimate pass *only* when it registers the loss, which also removes it as a blanket-avoidance route.
- **Re-review:** required.

---

## 5. Group A audit (unchanged by the review)

Fifteen boundary-only repairs audited against their complete records:

- **8 of 15** lack `presented poorly` in `If failed`, although their new boundary creates exactly that failure: `F1-03`, `F3-03`, `F3-04`, `F4-03`, `F6-03`, `F7-03`, `F7-04`, `F9-04`. Proposal: add the classification. No intended test changes.
- **15 of 15** carry the boundary as a top-level line rather than nested under `Reference`. Cosmetic; proposal to normalise.

**Audit conclusion:** the metadata-only exemption holds for *item content* and fails for *scoring metadata*. Neither change requires independent re-review.

---

## 6. `OD-10` — confidence standard (standards-level, not applied)

| Concept | Status |
|---|---|
| Confidence as **register** — how assured the phrasing sounds | **Never scored** |
| **Certainty unsupported by the record** — asserting a resolution the evidence does not support | **Prohibited** |
| **Hedging despite established evidence** | **Failure of `R-6`** (`R6-03a`) |

Proposed replacement text for both documents:

> `never scored` covers confidence as a register — the sound of assurance, independent of content. It does not cover certainty *claims*: asserting a resolution the record does not support is `prohibited`, and refusing a conclusion the record does support fails `R-6` (`R6-03a`).

**Requires owner approval.** `MEMORY-09` is an approved artefact.

### 6.1 Affected-unit audit (not only `F3-04`)

Every unit that uses confidence as a prohibited condition, found by scanning the item files rather than by memory:

| Unit | Field | Current text | State |
|---|---|---|---|
| `F3-04` | `prohibited` | "confident attribution to the fall" | **Blocked** — unchanged |
| `F9-04` | `Forbidden` | "a **confident** single cause; a third invented cause" | **Blocked** — unchanged |
| `F7-01` | `Forbidden` / `Prohibited` | was "a confident culprit" | **Reworked** to "a single culprit **asserted as fact**" — resolves the conflict for this unit |

So the conflict is **not** confined to the item the reviewer named. Two units remain non-compliant with the current standard and both wait on `OD-10`; only `F7-01` could be resolved without changing the standard, because "asserted as fact" says the same thing without relying on the word "confident".

**If `OD-10` is rejected** (confidence stays never-scored in all senses), then `F3-04` and `F9-04` need substantive rewrites rather than wording changes — the prohibition itself would be unregisterable.

---

## 7. `OD-11` — independence

`CORPUS-13` is a **primary** review artifact — better standing than the Muse review, which remains secondary. It is also a **single non-human reviewer**, which `CORPUS-05` §7 permits only as lower confidence. No two-reviewer convergence exists.

**Therefore:** `F3-03`'s `accept` is one reviewer's verdict, not acceptance; and the four rework items require a **second** independent reviewer who does not receive `CORPUS-13`'s verdicts.

---

## 8. Status

| State | Scope |
|---|---|
| **Corrected** | `F6-04` record (3 fields); `C-4` invariant added and scoped |
| **Accepted by one reviewer, not sealed** | `F3-03` |
| **Rework, proposals ready** | `P-01`, `F7-01` (partly), `F3-04` (blocked on `OD-10`) |
| **Awaiting second reviewer** | all five Group B items |
| **Audited, proposed** | 15 Group A classifications |
| **Escalated** | confidence standard (`OD-10`) |
| **Withdrawn** | `CORPUS-12-FOLLOWUP-AUDIT-AND-PROPOSALS` — superseded by this document |
| **Untouched** | 37 units; no splits changed; no new items |

**Not established:** that any item passes, that the corpus is valid, or that `G1` can close.

---

## 9. Applied in this cycle (owner-authorised)

| Work | State |
|---|---|
| `F6-04` record | **Corrected** — 3 fields aligned; `C-4` added and scoped to `R-2`/`R-3` |
| Group A classifications | **Applied** — 8 units gained `presented poorly`; all 15 boundaries nested under `Reference` |
| `P-01` | **Reworked** — silence passes only when it registers the loss |
| `F7-01` | **Reworked** — "asserted as fact" replaces "confident"; `Over-exposure` no longer double-scores `Prohibited` |
| `F3-04` | **Blocked** on `OD-10` — unchanged |
| `F3-03` | **Untouched** — accepted by one reviewer |
| `C-4` regression | **Both sides demonstrated** — `FX-08` caught; `FX-09` (`R-1`) and `FX-10` (`R-5`) produce no findings |

**Group A classification record.** Adding `presented poorly` does not change what a reply must contain, but it **does** extend the operative failure contract: a reply that says too much is now classifiable as that failure rather than merely unclassified. Recorded as *scoring metadata, not content* — the eighth unit's classification is as much a scoring change as an eight-line edit would be.

**Corpus after this cycle:** `PASS (0 fail, 0 warn, 0 unverifiable)`; fixtures `FAIL (4 fail, 1 warn, 1 unverifiable)` across 10 units.

---

## 10. Governance finding — the fourth incident, as a rule

This cycle's error was not a wrong measurement. It was **treating a registered file as an inspected one**. Three distinct propositions were conflated:

| Conflated | Actually |
|---|---|
| A file exists in the workspace | That it was read |
| A file is in the registry | That its content was inspected |
| An instruction refers to an artifact | That the artifact was checked |

Recorded as **`P-6`** in `CORPUS-05`:

> **Registration is not inspection.** Registering a file records its existence and hash; it says nothing about its content. A follow-up written from an instruction's description is not a substitute for reading the artifact the instruction refers to. Before acting on a reference, open it.

This is the fourth inference-as-measurement incident and the first one caused by *provenance routing* rather than by unmeasured claims. It is recorded as a governance finding, not only as a corrected document.

---

## 11. Status

| Area | Status |
|---|---|
| `F6-04` repair | Corrected; checker passes; **second review outstanding** |
| `C-4` | Revised; false-positive history preserved; both boundary sides fixtured |
| `F3-03` | Accepted by **one** reviewer; **provisional, not convergent** |
| `F3-04` | Blocked on `OD-10` |
| `F7-01`, `P-01` | Reworked; **whole-record review outstanding** |
| Group A | Classifications applied; scoring-contract effect recorded |
| `OD-10` | **Owner decision required** — amendment text prepared (§6) |
| Independent convergence | **Outstanding** — second reviewer must not receive `CORPUS-13` |
| Corpus integrity | Checker `PASS`; registry clean |
| `G1` | **Open; not sealed; expansion frozen** |

**`OD-10` status:** the owner's **recommendation** to approve the three-way distinction is recorded (2026-11-12). It is **not** an owner decision, so the amendment stays unapplied and `F3-04`/`F9-04` stay blocked. The recommended wording evaluates the relationship between a claim and its evidence, not whether a reply sounds assured.

**Second review prepared:** `CORPUS-16` covers **11 units** with split labels stripped programmatically — the six flagged units plus the five further units whose failure classifications changed, overlapping on `F3-03`, `F3-04`, `F9-04`. Presence of a classification is explicitly not validation.

**Freeze:** no new items, batches, or splits until `OD-10` is decided, the second review converges, and the reworked units pass whole-record review. The 34 units outside the packet stay untouched.
