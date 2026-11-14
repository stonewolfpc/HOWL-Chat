# `MEMORY-09` — Source-Discrepancy Correction Proposal (`SD-1` … `SD-5`)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-SDP-1` · **Version:** 2 · **Date:** 2026-11-13 · **Status:** **PROPOSAL — NOT AN AMENDMENT. Nothing is in force.**
**Revision 2** (2026-11-13) — owner-directed: `SD-5` no longer asserts a chronology the sources do not establish; `SD-3` no longer claims `A-2`'s applicability to a non-scored container is settled; Appendix A adds the proposed erratum to `MEMORY-09-RCN-1`; §12 records the owner-recommended dispositions as **recommendations, not decisions**.
**Authorization:** owner, 2026-11-13 — proposal development only.
**Controlling source for findings:** `MEMORY-09-RCN-1` (`MEMORY-09-RECONCILIATION-1-ACCOUNTING-R1-WITHHOLDING-2026-11-13.md`), with one correction recorded in §2.
**Effect:** **no source file modified.** Every proposed passage below is proposed text, not current text.

---

## 1. Purpose and boundary

Make five documentation discrepancies precise and reviewable so the owner can decide whether to correct them. This proposal supplies exact replacement wording, an amendment map, and impact analysis. It does **not** apply anything, and it does not change membership, scoring, evaluation semantics or historical evidence.

**Evidence labels used throughout:** **Observed** (read directly from source), **Source-derived** (follows from observed text), **Calculated** (arithmetic), **Inferred** (a reading, not a measurement), **Unestablished** (no source support), **Superseded** (historical text retained, no longer operative).

---

## 2. Correction to the controlling reconciliation

Before the five items, one accuracy correction — required because it changes the *shape* of `SD-1`.

**`MEMORY-09-RCN-1` §2.1 describes `46` as "every level-2/3 item heading".** That conflicts with `CORPUS-01` L166, which already states: *"42 item **headings**, 45 scored **units** — `R6-03a`–`d` are scored separately but share a heading."*

**Observed:** the source already uses "item heading" to mean the **family-level** heading, with the four provenance variants folded under `R6-03`. Under that convention 46 is not a heading count — it is the count of **raw heading lines** (`##`/`###`), which includes the four variant sub-headings.

**Calculated:** 42 item headings + 4 variant sub-headings = **46 heading lines**; 42 headings − 1 unscored container (`R6-03`) + 4 scored variants = **45 scored units**.

**Consequence:** `SD-1` is **narrower** than the reconciliation implied. `CORPUS-01` L166 is substantively correct and should be *sharpened*, not replaced. The genuine defects are: the §5 coverage line is unlabelled, the 46 line-count and 46 tag-count appear nowhere explained, and `MEMORY-09-RCN-1` itself uses "headings" for the raw count.

**Not corrected here:** `MEMORY-09-RCN-1` is left unchanged (no authorization to edit it). A separate erratum to that record may be warranted and is listed as `OD-required` in §9.

---

## 3. `SD-1` — Unit-count conventions

### 3.1 Defect

**Observed:** `CORPUS-01` L136 (§5 coverage line) sums to **42** (`2+4+4+4+10+6+12`); §4 L164 totals **45**; L166's note explains the 42/45 relation but not the coverage line's convention. L247 reports "Restraint modes tagged: **45 / 45**" while 46 tags exist.

### 3.2 Proposed replacement — `CORPUS-01` L166

> **Accounting basis (revised 2026-11-13).** Three counts are in use and are **not interchangeable**:
>
> | Count | Value | What it counts | Excludes |
> |---|---:|---|---|
> | Item headings | **42** | One heading per family-level item | The four `R6-03a`–`d` provenance variants, which sit under the `R6-03` family heading |
> | Scored units | **45** | Units carrying a complete scored body | The `R6-03` container, which has no `Required` or `Forbidden` |
> | Heading lines / tags | **46** | Every `##`/`###` heading line, each carrying a `Restraint:` tag | Nothing — this is the raw line count, including the container and the four variants |
>
> The §4 table and §7's "scored units" and ratio lines count **units**. The §5 coverage line counts **item headings**. Tag totals exceed scored units by one because the unscored `R6-03` container is itself tagged.

### 3.3 Proposed replacement — `CORPUS-01` L136

> **Coverage (after batch 3, item-heading convention):** `R-1` 2 · `R-2` 4 · `R-3` 4 · `R-4` 4 · `R-5` 10 · `R-6` 6 · `none` 12 = **42 item headings**, corresponding to **45 scored units** (see the accounting basis in §4).

### 3.4 Terminology recommendation

| Count | Recommended term | Never call it |
|---|---:|---|
| 42 | **item headings** | "units", "items" without qualification |
| 45 | **scored units** | "items", "headings" |
| 46 | **heading lines** (or **tagged headings**) | "headings" — collides with the 42 convention |

**Inferred:** the collision is precisely why three numbers appeared irreconcilable. One word, "headings", currently means both 42 and 46 depending on the reader.

---

## 4. `SD-2` — `R-6` arithmetic

### 4.1 Defect

**Observed:** `CORPUS-01` L248 reads *"`R-6` **resolved** (9 items, 4 mechanisms + 4 provenance variants); `R-2`/`R-3` still want two more each."*
**Calculated:** 4 mechanisms + 4 variants = **8**, not 9. **Unestablished:** which ninth element the figure counts; no source identifies it.

### 4.2 Proposed replacement — `CORPUS-01` L248

> **Restraint balance — `R-6` resolved.** Eight dedicated fixtures (`R6-01`, `R6-02`, `R6-03` with variants `R6-03a`–`d`, `R6-04`), of which **seven are scored units** — the `R6-03` container is a family heading and is not scored. `R-6` carries **ten tagged headings**, the eight fixtures plus legacy members `F7-01` and `F9-04`. `R-2`/`R-3` still want two more each.
> *Superseded record:* an earlier figure read "9 items, 4 mechanisms + 4 provenance variants"; its arithmetic yields 8 and the ninth element is not identified in any source.

### 4.3 The six counts, kept apart

| Count | Value | Definition | Source |
|---|---:|---|---|
| Table membership | 2 | `F7-01`, `F9-04` — predates the `R6` series | `CORPUS-01` §5 |
| Coverage line | 6 | Item-heading convention | `CORPUS-01` L136 |
| §7 stated | 9 | Superseded figure; arithmetic yields 8 | `CORPUS-01` L248 |
| Tag count | 10 | 8 fixture headings + 2 legacy | item files |
| Dedicated fixtures | 8 | 7 scored + 1 container | `CORPUS-04` |
| Scored units | 7 | Fixtures minus container | Calculated |

**No collapsing.** Each answers a different question; only the §7 "9" is defective, and it is defective as arithmetic, not as a convention.

---

## 5. `SD-3` — The `R6-03` container tag

### 5.1 Observed

`CORPUS-04` L53 headings `R6-03` with `Family: F-7 · Split: heldout (all four) · Abstention: no · Restraint: R-6` (L55), an uncertainty-source field, and framing prose. **No `Required`, no `Forbidden`.** Variants at L60, L75, L89, L104 carry their own bodies and their own `Restraint: R-6`.

### 5.2 Option A — retain the tag, document the container

| Consequence | Effect |
|---|---|
| Tag totals | Stay **46**; the convention in §3.2 explains the +1 |
| `A-2` | **Not settled by this proposal.** `A-2` speaks of one mode per **item**; whether it applies to an unscored family container is established nowhere in the sources. Retaining the tag produces no untagged restraint heading either way |
| Source churn | **None** — `CORPUS-04` untouched |
| Historical membership | Preserved exactly |
| Cost | Requires the documented convention to be read alongside any tag count |

### 5.3 Option B — remove the tag from the container

| Consequence | Effect |
|---|---|
| Tag totals | Become **45**, aligning tags with scored units |
| `A-2` | **Would conflict** unless the container were first established as outside `A-2`'s scope, which no source does: `A-2` states *"Untagged restraint items are not accepted"* |
| Source churn | One line edited in `CORPUS-04`; a reader comparing digests sees an item-file change |
| Historical membership | Altered — the container would no longer record its family's mode |
| Additional risk | Sets a precedent that family containers need not declare a mode, which no rule authorises |

### 5.4 Recommendation

**Option A.** It is the least disruptive defensible treatment: it changes no item file, preserves history, and the only cost — a documented convention — is already required by `SD-1`. Option B would be cleaner arithmetic bought with an item-file edit and an `A-2` question left dangling.

**Stated precisely, so it cannot be over-read:** this is a **membership-accounting treatment**. It documents that the `R6-03` container is a non-scored family heading that carries its family's mode tag. It is **not** a new interpretation of `A-2`, and it does **not** settle whether `A-2` applies to non-scored containers — that question remains open as `U-7`.

---

## 6. `SD-4` — `R-1` table membership

### 6.1 Defect

**Observed:** `CORPUS-01` §5 `R-1` row lists `P-03` and `F9-01`, tagged `R-5` and `R-4` respectively. No item listed in the table is tagged `R-1`; the tagged `R-1` items are `F6-02` and `P-01`.

### 6.2 Proposed treatment — retain the table, add an operational block

The table is **not rewritten**. A note and a tag-derived block are added beneath it:

> **Operational membership (tag-derived, 2026-11-13).** The item lists in this table are retained as a **historical record** and are superseded for operation by the per-item `Restraint:` tag (`OD-M09-05`). Current tag membership:
>
> | Mode | Tagged items |
> |---|---|
> | `R-1` | `F6-02`, `P-01` |
> | `R-2` | `F1-01`, `F1-02`, `F1-03`, `F6-01` |
> | `R-3` | `F3-01`, `F6-03`, `F6-04`, `P-04` |
> | `R-4` | `F8-01`, `F8-03`, `F9-01`, `F9-03` |
> | `R-5` | `F2-01`, `F2-02`, `F2-03`, `F5-01`–`F5-04`, `F9-02`, `F9-05`, `P-03` |
> | `R-6` | `F7-01`, `F9-04`, `R6-01`–`R6-04` incl. `R6-03a`–`d` |
>
> Where the two disagree, the tag governs. **Neither representation establishes that an item substantively exercises its mode** — see `MEMORY-09-RCN-1` §3 on `R-1`.

### 6.3 Boundary

**Inferred, and not converted into an accounting change:** `F6-02` carries a perceptibility requirement and `P-01` tests proportionate disclosure. Correcting the table's membership does **not** claim either item cleanly satisfies `R-1`, and does not re-tag anything.

---

## 7. `SD-5` — Duplicate memberships versus `A-2`

### 7.1 Defect

**Observed:** the table lists `P-03` under `R-1` **and** `R-5`; `F9-01` under `R-1` **and** `R-4`. **Observed:** `A-2` requires *"exactly one mode"* per item. **Unestablished:** whether the duplicates were ever an intentional multi-mode policy; no source states one.

### 7.2 Options

| Option | Description | Assessment |
|---|---|---|
| **(a) Labelled historical table** | Add a note: the table contains duplicate memberships that conflict with the currently governing `A-2` convention; duplicates are evidence, not assignments | **Least disruptive; no source text rewritten; no chronology asserted** |
| (b) Separate operational table | Split into a historical table and a tag-derived operational table | Larger restructure; duplicates the data `SD-4`'s block already supplies |
| (c) Explanatory note only | State the conflict and leave both as-is | Leaves the table readable as current, which is the actual hazard |

### 7.3 Recommendation

**(a), implemented by the same note block as `SD-4`** — one insertion resolves both, without touching the table:

> **Note on duplicated memberships.** This table's item lists contain duplicate family memberships that conflict with the currently governing `A-2` single-tag convention. The available sources do not establish when those memberships originated, whether they reflected an intentional earlier policy, or whether they were carried forward inadvertently. The table is retained as historical evidence, not as operational membership authority. An item appearing under two modes therefore records that both were considered, **not** that it carries two modes; `A-2` governs, and exactly one tag applies per item.

---

## 8. Amendment map

| ID | File | Anchor | Change type | Lines touched |
|---|---|---|---|---|
| SD-1 | `CORPUS-01` | L136 coverage line | replace one line | 1 |
| SD-1 | `CORPUS-01` | L166 accounting note | replace with 3-row convention table + 2 sentences | 1 → ~8 |
| SD-2 | `CORPUS-01` | L248 `R-6` balance line | replace one line, add superseded note | 1 → 2 |
| SD-3 | — | none | **no edit** (Option A is documentary) | 0 |
| SD-4 | `CORPUS-01` | after §5 mode table | insert operational-membership block | +~12 |
| SD-5 | `CORPUS-01` | after §5 mode table | insert duplicate-membership note (may merge with SD-4) | +~4 or shared |

**One file affected.** No item file, checker, threshold, split, score or standard changes. **Observed:** `check_corpus.py` reads item files, not `CORPUS-01`, so none of these edits can alter a checker result.

---

## 9. Impact analysis

| Question | Answer |
|---|---|
| Does any item, tag, split or score change? | **No** |
| Does membership change operationally? | **No** — tags already govern; `SD-4` only *records* what they say |
| Does the checker output change? | **No** — it does not read `CORPUS-01` |
| Does any evaluation claim change? | **No** — no result, ratio or verdict is recomputed |
| Do the counts change? | **No** — 42/45/46 keep their values; only their labels and scope become explicit |
| Is history rewritten? | **No** — superseded figures are retained and labelled, not deleted |
| Does `R-1`'s construct finding change? | **No** — `SD-4` records membership; `RCN-1` §3 records exemplification |
| Residual risk | A reader may still treat "42" as the corpus total; the proposed labels reduce, not eliminate, that |

---

## 10. Dependencies and ordering

1. **SD-1 first** — it defines the vocabulary (`item headings`, `scored units`, `heading lines`) that `SD-2` and `SD-3` rely on.
2. **SD-2** — uses `SD-1`'s terms for the `R-6` breakdown.
3. **SD-3** — a decision, not an edit; depends on `SD-1`'s tag-count explanation.
4. **SD-4 + SD-5** — one insertion, implementable together; depend on `SD-1` only for terminology.

No item file, checker or standard edit is a prerequisite for any of them.

---

## 11. Unresolved and unsupported

| # | Item | Status |
|---|---|---|
| U-1 | Which ninth element the §7 "9 items" counted | **Unestablished** — no source identifies it |
| U-2 | Whether "item heading" should include or exclude the `R6-03` container | **Inferred** — the source's 42 implies inclusion; not stated |
| U-3 | Whether the four variants should ever count as headings | Terminology decision, owner's |
| U-4 | Whether the §5 coverage line should stay heading-level | Proposal keeps it; alternative is recomputation to 45 |
| U-5 | `CORPUS-01` L168's own open discrepancy (checker 19/22 vs stated 26 held-out) | **Outside this proposal** — pre-existing, already flagged in source |
| U-6 | Whether `MEMORY-09-RCN-1` needs an erratum for the "46 headings" framing | **Recommended**; proposed text supplied in Appendix A, not authorised for filing |
| U-7 | Whether `A-2` applies to a non-scored family container | **Unestablished** — no source addresses it; `SD-3` Option A does not decide it |

---

## 12. Owner decisions required before implementation

| # | Decision | Options | **Owner-recommended** (2026-11-13) — recommendation, **not a recorded decision** |
|---|---|---|---|
| **D-1** | Adopt the three-term vocabulary and the §3.2 accounting block | adopt / amend wording / decline | **Adopt** |
| **D-2** | Adopt the `R-6` restatement and the superseded-figure note | adopt / keep "9" with a caveat / decline | **Adopt** |
| **D-3** | `R6-03` container tag | Option A (retain + document) / Option B (remove) | **Option A** |
| **D-4** | `R-1` operational-membership block | adopt / table-only note / decline | **Adopt** |
| **D-5** | Duplicate-membership note | adopt merged with `D-4` / separate / decline | **Adopt, with the chronology correction now in §7.3** |
| **D-6** | Erratum to `MEMORY-09-RCN-1` for the "46 headings" framing | authorise / leave as recorded | **Authorise** — text in Appendix A; the original reconciliation stays unchanged |
| **D-7** | Coverage line convention | keep heading-level / recompute to units | **Keep heading-level**, provided it is labelled |

**None of these is decided by this document.** Until they are, `CORPUS-01` keeps its current wording and every count keeps its current ambiguity.

---

## 13. R-1 substantive boundary

Restated, unchanged from `MEMORY-09-RCN-1` §3, and **scoped to the items examined**:

> No clean `R-1` exemplar has been established among the identified tagged and table-listed items examined (`F6-02`, `P-01`, `P-03`, `F9-01`).

This is **not** a finding about the whole corpus, and it does **not** authorise re-tagging, redefining `R-1`, or declaring `R-1` unmeasurable. Documentation corrections to membership (SD-4, SD-5) change *which items are recorded as members*; they say nothing about *whether those items exercise the construct*. The two questions must not be merged in review.

---

## 14. Status

**No source corrections have been applied.** `CORPUS-01`, `CORPUS-04`, `check_corpus.py`, all item files, tags, splits, scores, thresholds, `MEMORY-09`, `OD-M09-03`, the owner decision record, its addendum, `AM-01`, `OD-M09-05` and `MEMORY-09-RCN-1` are **unchanged**. This proposal is not an amendment and confers no authority to edit anything.

---

## Appendix A — Proposed erratum to `MEMORY-09-RCN-1` (not filed)

Proposed as a **separate** documentation artifact adjusting `MEMORY-09-RECONCILIATION-1-ACCOUNTING-R1-WITHHOLDING-2026-11-13.md`. The original record would remain unchanged and byte-identical; only the naming of one count is corrected.

**Proposed text:**

> **Erratum `RCN-1-E1` — 2026-11-13 (authority: owner decision `D-6`, to be recorded).**
>
> §2.1 described **46** as *"every level-2/3 item heading"*. That is a **raw heading-line count** and conflicts with the convention already stated in `CORPUS-01` L166, where "item headings" means the **42 family-level headings**, with the four `R6-03a`–`d` variants folded under the `R6-03` family heading.
>
> Corrected terminology:
>
> | Count | Term | Composition |
> |---:|---|---|
> | 42 | **item headings** | One per family-level item; includes the `R6-03` container |
> | 45 | **scored units** | 42 − 1 unscored container + 4 scored variants |
> | 46 | **heading lines** (equivalently, **tagged headings**) | 42 + 4 variant sub-headings; each line carries a `Restraint:` tag |
>
> The findings in §2.2–§2.3 are **unaffected**: the container's absence of a scored body, the 45/46/42 relation, and the six separate `R-6` counts all stand as recorded. Only the label applied to 46 changes, from "headings" to "heading lines".
>
> **Not affected:** no count, item, tag, membership or finding is revised. This erratum is terminological.

**Consequence if authorised:** the two artifacts then agree on terminology, and the phrase "46 headings" no longer appears anywhere as a convention name.

---

## 15. Provenance and revision history

| Layer | Fact |
|---|---|
| Authorization | Owner, 2026-11-13 — proposal only |
| Sources read | `CORPUS-01` (§4 L150–170, §5 L136, §7 L237–248), `CORPUS-04` L53–58 and variant headings, item files for tag data, `MEMORY-09-RCN-1`, `OD-M09-05` |
| Arithmetic verified | 42 + 4 = 46 heading lines; 42 − 1 + 4 = 45 scored units; 46 tags = 45 scored + 1 container |
| Model executed | No · `P-3` not started · no inference, benchmark or `G2` |
| Performed by | DeepSeek (per work order), repository worker — artifact-reported identity; labels are provenance data, not lineage proof |

| Rev | Date | Change |
|---|---|---|
| 1 | 2026-11-13 | Initial proposal. Includes §2, a correction to `MEMORY-09-RCN-1`'s "46 headings" framing discovered while verifying `SD-1` against `CORPUS-01` L166 — the source already labels 42 headings / 45 units, so `SD-1` was narrowed from "unlabelled conventions" to "one unlabelled line plus two unexplained totals". |
| 2 | 2026-11-13 | **Owner-directed corrections, proposal only.** (a) `SD-5` no longer asserts the table pre-dates `A-2`: §7.2 and §7.3 now state only that the duplicates conflict with the current convention, that their origin, intent and transmission are unestablished, and that the table is historical evidence rather than operational authority. (b) `SD-3` no longer claims `A-2`'s applicability to a non-scored container is settled: Option A is recorded as a **membership-accounting treatment**, with `U-7` opening the applicability question explicitly. (c) Appendix A supplies the proposed `MEMORY-09-RCN-1` erratum as separate text, leaving that record unchanged. (d) §12 records owner-recommended dispositions in a column labelled as recommendation, **not** as decisions. No source file touched. |
