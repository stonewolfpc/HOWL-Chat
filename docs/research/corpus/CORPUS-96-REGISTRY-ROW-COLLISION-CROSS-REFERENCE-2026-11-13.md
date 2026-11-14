# `CORPUS-96` — Registry Row Collisions and Legacy Identifier Cross-Reference

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `CORPUS-96` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **RECONCILIATION RECORD — no artifact altered.**
**Authority:** `OD-M09-14` (append-only reconciliation) as recorded in `MEMORY-09-ODR-1`; executed under work order `MEMORY-09-WO-REG-02` and its owner numbering amendment of 2026-11-13.
**Registry row:** **441**.
**Effect:** no document is deleted, merged, overwritten, renumbered or edited by this record. It is the authoritative *description* of two pre-existing collisions and of two legacy identifiers.

---

## 1. Purpose

Two numbered registry rows each carry two distinct documents, and two bracket-form identifiers appear in the record without a definition anywhere in the corpus. This record states both facts precisely, with verified metadata, so that later work can cite the truth without rewriting history.

---

## 2. Collision at registry row 423

| Identifier | Path | Title | Bytes | Lines | SHA-256 |
|---|---|---|---:|---:|---|
| `CORPUS-82` | `research/corpus/CORPUS-82-ARCHITECTURE-ADJUDICATION-Q2-2026-11-13.md` | CORPUS-82 — Q2 Architecture Adjudication: Ordinary Acceptance and Restraint-Mode Scoring | 52,403 | 281 | `959527990014dadb6381bb1473131cd6f0e3f083815e3d031ec66e87876a8f81` |
| `CORPUS-83` | `research/corpus/CORPUS-83-APPLICATION-SCORED-RESTRAINT-Q3-2026-11-13.md` | CORPUS-83 — Q3 Application-Scored Restraint Adjudication Report | 42,212 | 306 | `c6ecd9792d6e5fd904cf252e6cf594f9a34e7674d28c7d44b04cf5e0bdfb72d2` |

Both documents are distinct artifacts filed against the same registry number. **Neither is a duplicate entry for the other.**

---

## 3. Collision at registry row 424

| Identifier | Path | Title | Bytes | Lines | SHA-256 |
|---|---|---|---:|---:|---|
| `CORPUS-90` | `research/corpus/CORPUS-90-APPLICATION-SCORED-RESTRAINT-DECLARED-NON-INDEPENDENT-ADJUDICATION-2026-11-13.md` | `CORPUS-90` — Q3 Application-Scored Restraint: Declared Non-Independent Adjudication | 69,527 | 599 | `a53427d1ec6130c491ac691da606fc0ba41bbe3efbb2a040b649aabebc2a69a7` |
| `CORPUS-91` | `research/corpus/CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md` | CORPUS-91 — Independent Adversarial Review of CORPUS-90 and Q3 Restraint Scoring | 56,012 | 627 | `fda9b5e12f49fa308f24d38641c866dbd759e917cc196ee2b309753c694d542a` |

---

## 4. What a collision is, and what is done about it

| | |
|---|---|
| **What happened** | Two separate documents were registered against one row number, twice. The registry's numbers are therefore not unique as keys. |
| **What it is not** | Not a duplicated entry for a single file; not a failed write; not an indication that either document is invalid. |
| **Action taken** | **None** on the four documents. Rows 423 and 424 remain exactly as found. `OD-M09-14` forbids renumbering them, and this record does not. |
| **Consequence** | Any citation to "row 423" or "row 424" is ambiguous unless it names the file. Cite by **path and SHA-256**, as tabulated above. |

**Historical state preserved:** the assignments above are exactly the state encountered on 2026-11-13. This record does not imply they were ever corrected elsewhere.

---

## 5. Legacy identifiers `OD-9[T]` and `OD-10[T]`

### 5.1 Where they occur

An exhaustive text search of `docs/` finds the bracket form in **exactly two files** — in fact in **one** file, at two locations:

| Citation | File | Location | Status |
|---|---|---|---|
| `OD-9[T]` / `OD-10[T]` | `research/MEMORY-09-OWNER-DECISION-RECORD-OD-M09-2026-11-13.md` | L42 (decision register) | **Two occurrences of each**, same sentence |
| `OD-9[T]` / `OD-10[T]` | same | L166 (§2.12) | **Two occurrences of each**, same heading |

**They appear in no corpus document.** The work order anticipated that each identifier "appears twice in corpus documentation"; that is not what the workspace shows. They occur twice each, but both occurrences sit inside the owner decision record's own description of `OD-M09-13`.

### 5.2 What the underlying decisions are

The legacy labels do have identifiable referents in the corpus, even though the bracket form does not:

| Legacy label | Corpus referent | Evidence |
|---|---|---|
| `OD-9` | "What governs when an item passes but its mode is not demonstrated?" — gating options A/B/C | `CORPUS-76` §3 (L45–47) |
| `OD-9` | "dual-field scoring specification (proposed)" — the two result dimensions | `CORPUS-78` §3 (L41) |
| `OD-10` | "What counts as the perceptible cue?" — cue candidates | `CORPUS-76` §4 (L63–65) |
| `OD-10` | "layered evidence framework (proposed)" — the five layers | `CORPUS-78` §4 (L95) |

### 5.3 Mapping table

| Historical citation | Source and location | Referent established by source | Candidate current identifier | Mapping status | Evidence |
|---|---|---|---|---|---|
| `OD-9[T]` | `MEMORY-09-ODR-1` L42 | Legacy family covering mode-gating **and** dual-field scoring | `OD-M09-02` (gating) **and** `OD-M09-06` (vocabulary) | **Ambiguous** | `CORPUS-76` §3 asks the gating question; `CORPUS-78` §3 specifies the two dimensions. Two decisions answer two different parts |
| `OD-9[T]` | `MEMORY-09-ODR-1` L166 (§2.12) | Same | Same | **Ambiguous** | Same |
| `OD-10[T]` | `MEMORY-09-ODR-1` L42 | Legacy family covering the perceptible cue **and** the layered framework | `OD-M09-03` (cue standard) | **Ambiguous** | `CORPUS-76` §4 asks what the cue is; `CORPUS-78` §4 supplies the layers. `OD-M09-03` adopts the direction but explicitly not the mechanics |
| `OD-10[T]` | `MEMORY-09-ODR-1` L166 (§2.12) | Same | Same | **Ambiguous** | Same |

### 5.4 Why no mapping is recorded as confirmed

1. **The bracket form is undefined.** `[T]` appears nowhere in the corpus and no document explains what `T` denotes. A notation with no definition cannot yield a confirmed identifier.
2. **The relation is many-to-many.** `OD-9` spans two separate decisions (`OD-M09-02` gating, `OD-M09-06` vocabulary); `OD-10` spans the cue standard (`OD-M09-03`) and a five-layer framework broader than that decision alone.
3. **`OD-M09-13` is conditional and unexercised.** It permits assigning `OD-M09-nn` identifiers *after* `OD-M09-01`–`03`, and nothing in the record shows such an assignment being made.

**Therefore:** the ambiguity is preserved. No `OD-9[T]` or `OD-10[T]` text was rewritten, replaced or annotated in any file, and no new identifier was assigned.

---

## 6. Provenance

| Layer | Fact |
|---|---|
| Work order | `MEMORY-09-WO-REG-02` + numbering amendment (rows 438–442) |
| Metadata | Recomputed from the files at write time, not copied from the registry |
| Executor | DeepSeek (per work order), repository worker — artifact-reported identity; labels are provenance data, not lineage proof |
| Model executed | No · `P-3` not started · no inference, benchmark or `G2` |
| Files altered | **None** — this record is additive only |
