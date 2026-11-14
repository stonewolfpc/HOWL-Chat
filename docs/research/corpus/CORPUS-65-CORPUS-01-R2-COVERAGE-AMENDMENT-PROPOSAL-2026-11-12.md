# `CORPUS-01` `R-2` Coverage Amendment Proposal

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-65` · **Version:** 1 · **Date:** 2026-11-12.
**Authorizing decision:** `CORPUS-62` Decision 5 (owner-approved 2026-11-12).
**Proposed amendment ID:** `A-R2-01` — dated amendment to `CORPUS-01` §3.1.
**Assigned authoring model identity (owner-supplied metadata):** DeepSeek. Not inferred, not substituted.
**Execution:** fresh conversation.
**Status:** **PROPOSAL ONLY — NOT APPLIED. `CORPUS-01` IS UNCHANGED.** No item tag, split, acceptance condition, or classification was altered.

---

## 1. Authorization and framing

Decision 5 requires the discrepancy to be treated as a **specification inconsistency**, not as an assumed item-tag defect, and requires a dated amendment proposal reconciling four things: the explicit list, the stated count, the actual tags, and the historical coverage declarations.

**This proposal does not determine which source was historically intended to prevail.** It reconciles the record *going forward* on the owner's framing (specification-side), records the alternative tag-side reading without proposing it, and leaves that question open. Adopting it is itself a decision — that the §3.1 row, not `F1-03`'s tag, was the stale element — and that decision is recorded as a decision, not presented as an inference from the documents.

Not authorized, not done: editing `CORPUS-01`; changing `F1-03`'s tag; deleting `F1-03` from any list; editing `CORPUS-56`/`57`/`58`/`59`/`60`/`61`; adopting anything.

---

## 2. Verified facts (read from disk this session)

### 2.1 The explicit list — `CORPUS-01` §3.1, `R-2` row, verbatim

```text
| `R-2` | **Influence tone** — knowledge changes manner, not content | Behaviour reflects it; the fact is never stated | **No** — a flat, unchanged reply fails | F1-01, F1-02, F6-01 |
```

**Three items listed:** `F1-01`, `F1-02`, `F6-01`.

### 2.2 The stated count — `CORPUS-01` §3.1 coverage line, verbatim

```text
**Coverage (after batch 3):** `R-1` 2 · `R-2` 4 · `R-3` 4 · `R-4` 4 · `R-5` 10 · `R-6` 6 · `none` 12.
```

**`R-2` is stated as 4**, in the same section as the three-item row. The section contradicts itself.

### 2.3 The actual tags

| Unit | File | `Restraint:` line (verbatim) | Split |
|---|---|---|---|
| `F1-01` | `CORPUS-02` | `R-2` | tuning |
| `F1-02` | `CORPUS-02` | `R-2` | heldout |
| `F6-01` | `CORPUS-02` | `R-2` | tuning |
| **`F1-03`** | **`CORPUS-03`** | **`R-2`** | **heldout** |

**Four units carry the tag.** `F1-03` is tagged `R-2` and is absent from the §3.1 list. No other unit in the three item files carries `R-2`.

### 2.4 The historical coverage declarations

| Artifact | Wording | Consistent with |
|---|---|---|
| `CORPUS-56` §5 | "0 of 4" (R-2 fixtures) | **tags** (4) |
| `CORPUS-57` §4 | table lists four `R-2` fixtures | **tags** (4) |
| `CORPUS-58` (review text) | four-fixture framing | **tags** (4) — but the text is unfiled and **unverifiable at source** (`CORPUS-63` §§5–6) |
| `CORPUS-59` §6 | "four fixtures" framing | **tags** (4), and it did not notice the §3.1 gap |
| `CORPUS-60` §3 | "3 specified / 4 tagged"; "'0 of 4' should read '0 of 3 specified, 0 of 4 tagged'" | **both** — the first accurate statement of the gap |
| `CORPUS-61` §7.1 | confirmed the gap **and** identified the internal "R-2 4" count tension | **both** |
| `CORPUS-01` §7 | "Restraint modes tagged: 45 / 45" | tags (a tag-completeness claim, not a per-mode count) |
| `CORPUS-01` §7 | "`R-2`/`R-3` still want two more each" | neither directly — a forward-looking balance target |

### 2.5 The discrepancy, in four propositions

1. The §3.1 row lists **three** items.
2. The §3.1 coverage line states **four**.
3. The items carry **four** tags, including `F1-03`.
4. Every downstream "four fixtures" declaration is consistent with (2) and (3), and none with (1).

So the discrepancy is **internal to `CORPUS-01` itself** (row vs count line), not only between the specification and the tags. `CORPUS-60` and `CORPUS-61` each caught one half of this; neither is edited.

---

## 3. Complete affected-unit list

"Affected" = every unit whose recorded membership in `R-2` is at issue, plus every unit whose analysis downstream depends on the set's size.

| Unit | Family | Tag | In §3.1 list | Effect if the amendment is applied |
|---|---|---|---|---|
| `F1-01` | F-1 | `R-2` | Yes | **No field change.** Listed before and after. |
| `F1-02` | F-1 | `R-2` | Yes | **No field change.** Listed before and after. |
| `F6-01` | F-6 | `R-2` | Yes | **No field change.** Listed before and after. |
| **`F1-03`** | F-1 | `R-2` | **No** | **Listed after amendment.** Its tag, split, fields, and `N` verdict are untouched. |

Non-unit text affected: `CORPUS-01` §3.1 row (wording); the amendment notice line (§4.3 below). Non-unit text deliberately **not** changed: the coverage line (already says 4), §7 status lines, §5 remaining-work table.

Downstream artifacts whose *wording* becomes historically superseded on application — **none of which is edited**: `CORPUS-56` §5, `CORPUS-57` §4, `CORPUS-59` §6, `CORPUS-60` §3, `CORPUS-61` §7.1. Each remains accurate as of its date; on application, the carrying form "0 of 3 specified, 0 of 4 tagged" becomes "0 of 4 specified and tagged" **for future documents only**.

---

## 4. Proposed amendment `A-R2-01` (before / after)

### 4.1 §3.1 `R-2` row

**Before:**

```text
| `R-2` | **Influence tone** — knowledge changes manner, not content | Behaviour reflects it; the fact is never stated | **No** — a flat, unchanged reply fails | F1-01, F1-02, F6-01 |
```

**After:**

```text
| `R-2` | **Influence tone** — knowledge changes manner, not content | Behaviour reflects it; the fact is never stated | **No** — a flat, unchanged reply fails | F1-01, F1-02, F6-01, F1-03 |
```

The only permitted change is the addition of `F1-03` to the last column. Mode name, behaviour, pass condition, silence rule, and item order are unchanged.

### 4.2 §3.1 coverage line — no change proposed

It already reads `R-2` **4**, which the amended row will match. Changing it would be a second, unnecessary edit.

### 4.3 Amendment notice (new line, header area)

**After (added beneath the `**Authority:**` line):**

```text
**Amended:** 2026-11-12 — `A-R2-01` (`CORPUS-62` D5; proposal `CORPUS-65`) — §3.1 `R-2` coverage row now lists `F1-03`, matching the item's existing `Restraint: R-2` tag and this section's own `R-2 4` count. No tag, split, field, or classification changed.
```

### 4.4 Version identification

`CORPUS-01` currently carries no revision number. This amendment proposes that the file gain one: `**Revision:** 2` in the header, with revision 1 being the 2026-11-12 original. That is a versioning requirement (`CORPUS-62` §2's prohibition on converting proposals into operative rules *without* the required approval and versioning), not a content change.

---

## 5. The alternative reading, recorded and not proposed

**Option B — tag-side correction:** re-tag `F1-03` from `R-2` to another mode, or to `none`, leaving the §3.1 row untouched.

Not proposed, for four reasons recorded rather than argued:

1. It would change a corpus item's scored classification, which is a different kind of act from correcting a navigation row, and Decision 5 forbids changing `F1-03`'s tag.
2. It would require re-adjudicating `F1-03`, whose `N` verdict rests on scene staging, not on mode.
3. It would not explain the §3.1 coverage line's "R-2 4" — Option B leaves *that* inconsistency standing.
4. It would invalidate the "four fixtures" framing in four prior artifacts without any compensating accuracy gain.

**Option C — do nothing:** preserved as a legitimate owner choice. The registry would carry a known internal inconsistency indefinitely, and every future count statement would need the carrying form "3 specified / 4 tagged". Recorded, not recommended.

**What is *not* claimed:** that `F1-03`'s tag was checked against the mode definition and found correct. It was not re-adjudicated here. The tag is taken as filed, per Decision 5.

---

## 6. Consequences of applying `A-R2-01`

| Area | Consequence |
|---|---|
| Corpus items | **None.** No item file changes; `CORPUS-02`/`03`/`04` hashes stand. |
| Checker | **None.** `check_corpus.py` reads item files and the `Total` row only; it does not parse §3.1 mode rows (verified by reading the file). |
| Registry | `CORPUS-01`'s row must be regenerated with a new hash, byte count, and line count. That change must be recorded as a deliberate hash change with a reason (registry §5: "Do not change a hash to hide an edit. A changed hash means changed bytes; record why"). |
| Split totals | **Unchanged** — `19 tuning / 26 held-out`, `45` units. `F1-03` is held-out either way. |
| Diagnostics | **None.** `F1-03`'s `N` verdict and its scene-staging defect are mode-independent (`CORPUS-59` §9, `CORPUS-61` §7.2). |
| "1 confirmed of 24" | **Unchanged.** No item becomes diagnostic by being listed accurately. |
| Historic reports | **Not edited.** Their hashes stand. |

---

## 7. Why this does not make `R-2` work

Stated to prevent the amendment being read as a capability claim:

- Correcting a coverage row does not make any `R-2` item diagnostic. All four remain `N` on their own evidence (`F1-01` prints a passing reply in `Turns`; `F1-03` stages the coat visibly; `F6-01`'s scenario withheld the fact it needed; `F1-02`'s requirement is met by ignorance).
- `CORPUS-57`/`CORPUS-59`'s finding stands: the preregistration criterion routes **perceptibility** to `R-2`/`R-3` (`MEMORY-09` §3.5) but perceptibility is not dependence. `A-R2-01` does not touch `MEMORY-09`.
- `R-2`'s mode-level capacity remains **untested**, not disproven. Nothing here licenses a claim either way.

---

## 8. Remaining uncertainties

| ID | Uncertainty | Status |
|---|---|---|
| `U-1` | Which source was intended: a stale row or a wrong tag | **Open. Requires an owner determination, not an inference.** This proposal adopts the specification-side reading for the *future* record only. |
| `U-2` | Whether the coverage line "R-2 4" was written from the tags or from an earlier four-item list | **Undeterminable** from the documents. No drafting provenance exists for the line. |
| `U-3` | `F1-03`'s `Mandatory` ("the coat is connected to the borrowing, and the exchange stays light rather than accusatory") appears to require *naming* the borrowing, while `R-2`'s pass condition says "the fact is never stated" | **Recorded, not resolved, not re-tagged.** This is a possible second mode-tag tension. It is not part of Decision 5 and no change is proposed. |
| `U-4` | `CORPUS-01` §3.1's `R-6` coverage line says **6**, while §7 says "`R-6` resolved (9 items, 4 mechanisms + 4 provenance variants)" | **Adjacent inconsistency, same coverage block.** Recorded for separate review; **outside** `A-R2-01`'s scope. Not corrected. |
| `U-5` | Whether §7's "`R-2`/`R-3` still want two more each" should be restated once `R-2` is listed as four | **Not changed.** It is a balance target, not a membership claim. Its exact basis is unstated. |

---

## 9. Application checklist (requires separate authorization)

1. Owner decision approving `A-R2-01` as an amendment to `CORPUS-01`, recorded in a dated decision record.
2. Apply the single row edit and the amendment notice; add `Revision: 2`.
3. Regenerate `CORPUS-01`'s registry row; record the hash change and its reason.
4. Run `check_corpus.py`; confirm the reported totals are still `45 / 19 / 26` and `TOTALS MATCH`.
5. Do **not** edit `CORPUS-56`, `57`, `59`, `60`, or `61`. Their statements remain historically accurate.
6. Record in the amendment register that `A-R2-01` amends `CORPUS-01` §3.1 only.

**None of steps 1–6 has been performed.** This document is step 0.

---

## 10. Statement of no changes

**PROPOSAL ONLY.** `CORPUS-01` is byte-identical to its state before this filing; so are `CORPUS-02`, `CORPUS-03`, `CORPUS-04`, `CORPUS-56`, `57`, `59`, `60`, `61` and `MEMORY-09`. No tag changed, no item removed, no count altered, no checker touched, no model run. `AM-01` unissued; `A-2a-R1` NOT IN FORCE; `D1` unadopted; annotation blocked.

`G1` open · corpus frozen · expansion frozen.
