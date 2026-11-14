# `OD-M09-05` — Membership Divergence Record

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `OD-M09-05-DIV-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **RECORDED — documentation only.**
**Authority:** `OD-M09-05` in `docs/research/MEMORY-09-OWNER-DECISION-RECORD-OD-M09-2026-11-13.md` §2.2 — tags adopted as **preferred membership authority**, all divergence from the tables and coverage statements to be documented, existing memberships preserved.
**Effect on sources:** **none.** No item, tag, table, coverage statement or count was changed, added, corrected or re-tagged in producing this record.

---

## 1. What "preferred authority" does and does not mean

| | |
|---|---|
| **Preferred** | Where the two sources disagree, an item's `**Restraint:**` tag is the operative membership for future work |
| **Not deleted** | The specification table's memberships remain on the record, unchanged, as a divergent historical source |
| **Not retroactive** | This record does not make past table-based statements correct or incorrect; it makes their divergence visible |
| **Not a decision about correctness** | Where the two disagree about *which* mode an item exercises, this record does not decide who is right (§4) |

---

## 2. Method

Enumerated 2026-11-13 from the three authored item files only:

- `docs/research/corpus/CORPUS-02-BATCH-1-ITEMS-2026-11-12.md`
- `docs/research/corpus/CORPUS-03-BATCH-2-ITEMS-2026-11-12.md`
- `docs/research/corpus/CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md`

**Units** = level-2/3 headings matching ``^#{2,3} `([A-Za-z0-9-]+)` `` — **46 headings, 46 unique IDs, zero duplicates**.
**Tags** = the first `**Restraint:** \`R-n\`` / `` `none` `` line under each heading — **46 tags, one per unit**.

**Table memberships** = the item lists in the `CORPUS-01` §5 mode table (`R-1` … `R-6`) — **16 citations across 14 distinct items**.

Matrix compared: table → tags, and tags → table. Both directions are reported because each reveals a different defect: a table member tagged elsewhere is a **contradiction**; a tagged unit absent from the table is an **omission**.

---

## 3. Measured inventory

| Mode | Table (distinct items) | Coverage line (after batch 3) | Tags (units) |
|---|---|---|---|
| `R-1` | 2 — `P-03`, `F9-01` | 2 | 2 — `F6-02`, `P-01` |
| `R-2` | 3 — `F1-01`, `F1-02`, `F6-01` | 4 | 4 — `F1-01`, `F1-02`, `F1-03`, `F6-01` |
| `R-3` | 2 — `P-04`, `F6-03` | 4 | 4 — `F3-01`, `F6-03`, `F6-04`, `P-04` |
| `R-4` | 3 — `F8-01`, `F8-03`, `F9-01` | 4 | 4 — `F8-01`, `F8-03`, `F9-01`, `F9-03` |
| `R-5` | 4 — `F2-01`, `F9-02`, `F9-05`, `P-03` | 10 | 10 — `F2-01`, `F2-02`, `F2-03`, `F5-01`, `F5-02`, `F5-03`, `F5-04`, `F9-02`, `F9-05`, `P-03` |
| `R-6` | 2 — `F7-01`, `F9-04` | 6 | 10 — `F7-01`, `F9-04`, `R6-01`, `R6-02`, `R6-03`, `R6-03a`…`R6-03d`, `R6-04` |
| `none` | — (no `none` row exists) | 12 | 12 — `F3-02`, `F3-03`, `F3-04`, `F4-01`, `F4-02`, `F4-03`, `F7-02`, `F7-03`, `F7-04`, `F8-02`, `F8-04`, `P-02` |
| **Total** | 16 citations / 14 distinct | 42 | **46** |

Three observations that hold across the whole table:

1. The table is **frozen at an earlier batch state**. It names 14 distinct items; 46 units are tagged.
2. The coverage line and the tags **agree exactly** for `R-2`, `R-3`, `R-4`, `R-5` and `none`. The two places they disagree are `R-1` and `R-6`.
3. `R-6` and `none` are the only modes whose dedicated item series (`R6-*`) exists at all; the table was never updated for the `R6` series or for `none`.

---

## 4. Divergences

### D-1 · `R-1` — total conflict, not staleness

| Direction | Items |
|---|---|
| Table says `R-1` | `P-03`, `F9-01` |
| Tags say `R-1` | `F6-02`, `P-01` |
| `P-03` is actually tagged | `R-5` |
| `F9-01` is actually tagged | `R-4` |

**Neither table member is tagged `R-1`, and neither tagged `R-1` member appears in the table.** The count agrees (2 and 2) while the *sets are disjoint* — the one divergence a count-based check would never catch, and the reason this record compares sets rather than totals.

**A hypothesis, labelled as such:** the table predates a later re-tagging in which `P-03` moved to `R-5` and `F9-01` settled at `R-4`, with `F6-02` and `P-01` taking their place. Supporting this is the table's own structure: it lists `P-03` and `F9-01` **twice** (see D-4), which is what a pre-`A-2` table looks like. This record does not adopt the hypothesis and does not change the tags.

### D-2 · `R-6` — the table omits the entire dedicated series

The table names `F7-01` and `F9-04` only. Eight `R6-01`…`R6-04` units — four mechanisms plus four provenance variants — are tagged `R-6` and absent from the table. The table therefore describes `R-6` as a two-item mode when it is the corpus's most heavily authored mode.

### D-3 · Systematic omission across all modes

`R-2` omits `F1-03`; `R-3` omits `F3-01` and `F6-04`; `R-4` omits `F9-03`; `R-5` omits six units (`F2-02`, `F2-03`, `F5-01`…`F5-04`). In every case the tag is present and the table entry absent. Sixteen table citations against 46 tagged units is a 30-unit gap.

### D-4 · The table permits multi-membership; `A-2` forbids it

`P-03` appears under `R-1` **and** `R-5`; `F9-01` appears under `R-1` **and** `R-4`. `CORPUS-01` `A-2` requires **exactly one** mode per item, and the single tag per unit satisfies it — 46 units, 46 tags.

**Structural consequence:** the table cannot be read as a membership list under `A-2` without collapsing its duplicates. `A-2` is the operative rule; the table records a superseded, multi-membership view of mode assignment.

### D-5 · Four different totals for the same corpus

| Source | Claim |
|---|---|
| `CORPUS-01` §5 coverage line | 42 units (`2+4+4+4+10+6+12`) |
| `CORPUS-01` §7 status | "Item headings: **42** · scored units: **45**" |
| Tags measured here | **46** units |
| Headings measured here | **46** headings, 46 unique |

The coverage line's 42 equals the tagged total **minus four** — precisely the `R6-03a`…`R6-03d` variants, which the coverage line does not count as separate units. The §7 "45" is one short of the 46 headings measured, and this record does not resolve which unit the figure excludes. **No figure is corrected here**; all four are recorded as they stand.

### D-6 · `R-6` has five different counts

| Source | Count |
|---|---|
| `CORPUS-01` §5 table | 2 |
| `CORPUS-01` §5 coverage line | 6 |
| `CORPUS-01` §7 "`R-6` resolved (9 items, 4 mechanisms + 4 provenance variants)" | 9 |
| Tags | 10 |
| Dedicated `R6-*` units | 8 |

The §7 line's own arithmetic (4 + 4) yields 8, not 9, so the "9" is internally unexplained as well as inconsistent with the tags.

---

## 5. What this record deliberately does not resolve

1. **Whether `R-1` has any members at all.** Under the tags it has two (`F6-02`, `P-01`); under the table it has two different ones. Deciding which pair is right requires reading those four items and is not a documentation act.
2. **Whether `P-03`'s and `F9-01`'s dual table membership was ever intended.** `A-2` says no item carries two modes, so either the table was wrong at authoring time or `A-2` post-dates it. This record states the conflict, not the cause.
3. **The 45-versus-46 scored-unit arithmetic.** One unit's status is unresolved; no attempt is made to identify it here beyond noting the gap.
4. **Whether `none` should have a table row.** It has twelve tagged units and no table representation, which is an omission of shape rather than of item.
5. **Any re-tagging.** No tag was changed, added or removed.

---

## 6. Consequence under the adopted authority

- **Operative membership = tags**, 46 units, one mode each, satisfying `A-2`.
- **The table is retained unchanged** as a divergent record, and its 16 citations remain citable with the divergence disclosed — as this record discloses it.
- **Coverage statements are flagged, not corrected.** The §5 coverage line and §7 counts remain as written; anyone relying on them must note D-5 and D-6.
- **`R-6` membership for future work is the 10 tagged units**, not the table's 2 and not §7's 9.
- **`R-1` is the one open membership question** requiring item-level reading before any future work depends on it.

---

## 7. Reproduction

| Step | Rule |
|---|---|
| Units | `^#{2,3} ` + backtick + `[A-Za-z0-9-]+` + backtick, in the three item files |
| Tags | `**Restraint:**` followed by a backticked `R-1`…`R-6` or `none`, first occurrence under each heading |
| Table | `^\| ` + backticked `R-n`, item IDs backticked in the final cell of `CORPUS-01` §5 |
| Expected | 46 units · 46 tags · 16 table citations across 14 distinct items · no duplicate unit IDs |

A parser-level count is reproducible; it does not adjudicate whether any individual assignment is **correct**.

---

## 8. Provenance and limits

| Layer | Fact |
|---|---|
| Decision recorded | `OD-M09-05`, 2026-11-13 |
| Measured | 2026-11-13, read-only, from the three item files and `CORPUS-01` |
| Authored by | Buffy / Codebuff repository worker — artifact-reported identity; labels are provenance data, not lineage proof |
| Sources changed | **None** — this record is the only artifact created |
| Model executed | No · `P-3` not started · no inference, benchmark or `G2` |
| Limit | Heading-based extraction; a unit whose tag sits under a differently-styled heading would be missed. Zero duplicates and a 1:1 heading-to-tag ratio make that unlikely but not excluded |
| Limit | Table memberships are read from `CORPUS-01` §5 only; other documents may carry their own item lists, unexamined here |
