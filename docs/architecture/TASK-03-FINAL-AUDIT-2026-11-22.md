# TASK-03 — v2 final audit: contradictions that the addendum stack created

| Field | Value |
|---|---|
| Purpose | Final contract-level audit of v2 (row 554) plus amendments 555, 556, 557 |
| Status | **Adds V2-90…V2-104 and VT-40…VT-47. Contract work only. No implementation authorized.** |
| Method | Four documents read end to end; each amendment clause tested for a home in the decision table |
| Recovery token | `TASK-03-FINAL-AUDIT-2026-11-22` |

---

## 1. The result the audit actually produced

The call-graph and pair-consistency questions are settled and I am not reopening them. What remains is **not another missing clause**. It is that the decision table (§7 of row 554) is the only executable part of the contract — "applied in order; the first matching row wins" — and every amendment since has imposed **preconditions on existing rows** while living in a different document.

A scanner implemented from the table alone would reproduce every contradiction below. So the contract cannot be closed by adding rows 555–557 to row 554's file list; **v2 needs a consolidation in which the amendment conditions are inline in the table**.

That is the answer to "are there other contract-level contradictions": yes, eleven, and they share one cause.

---

## 2. Contradictions

**CA-01 — Content with marker `none` has no row, and VT-18 asserts it is Healthy.** `SaveChatAs` creates a directory and writes only `history.json` (row 549 P6); `NewChat` calls `SaveChat` before any append. The copy therefore holds ≥1 record while the marker is `(none,0)`. Rows 5, 6, 7, 8, 10, 11 all require a marker value that excludes `none`-with-content, so **no row matches** — yet VT-18 states "a conversation whose log is duplicated scans Healthy".

Two ways to resolve it, and they are not equivalent:

| Resolution | Consequence |
|---|---|
| Add a row: `none` + content → Healthy | This makes the **probe 12 mixture Healthy**. Unattributed content would be legitimised on sight. Rejected. |
| `none` + content → **Pending**, reason `unattributed`; duplication writes `(pending,1)` → `(committed,1)` as part of the copy | Fail-closed, and it forces §4.2 to describe duplication. Adopted below. |

**VT-18 as written is wrong** and must be replaced, not reinterpreted.

**CA-02 — V2-09 and V2-10 contradict each other on the empty log.** V2-09: valid iff every non-blank line parses **and the file ends with a newline**. V2-10: an empty `history.ndjson` is 0 bytes. A 0-byte file does not end with a newline, so under V2-09 the post-clear log is **invalid**, while V2-10 calls it empty — and row 5 requires a representation to be "structurally valid **and** empty". `RewriteHistory([])` produces exactly this file (row 548 measurement), so the contradiction lands on the most common clear outcome.

**CA-03 — Row 8 can return Healthy on a marker the amendments say is not authoritative.** Row 8's evidence is presence, validity, ≥1 record and coherent coexistence — it has **no marker-currency condition**. Row 557 V2-86 says a marker read returning `Recovered` never counts as evidence that a transition is settled. A recovered `(committed,G)` with valid content therefore satisfies row 8 and violates V2-86. Row 8 must carry the condition inline.

**CA-04 — Row 6 names a reason the contract cannot justify.** Row 6 reports content under a `cleared` marker as `retracted-residue`. But V2-16 states the generation does not identify records, and probe 12 showed an append landing in a log after an interrupted clear. Content present under `cleared` is therefore **unattributable** — it may be G's residue or G+1 records written without a marker. Calling it retracted asserts provenance the contract spent three documents denying.

**CA-05 — V2-19 says no state is stored; §4 stores one.** The marker's `state` field holds `pending`/`committed`/`cleared`. V2-19's "No state is stored" is false as written. It is meant to say the six *derived* states are not stored, which is a different claim.

**CA-06 — V2-55 forbids the scanner from resolving a transition while requiring it to promote.** `pending → committed` is a resolution of the marker. The amendment needs the distinction it left implicit: **marker resolution** (the scanner may do this, under V2-56's provenance) versus **content-operation settlement** (the scanner may never do this).

**CA-07 — Deletion is both exempt and gated.** V2-53 (row 555) lists "append, rewrite, clear, save or duplication" — delete is not among them. V2-75 (row 556) says every mutation entry point must pass through the gate, and `DeleteChat` is one of the twelve. One of the two is wrong. Deleting destroys both representations and is a legitimate way out of an unresolved transition, so the correct resolution is an **explicit exemption**, not silence.

**CA-08 — v2 is silent on concurrent instances.** Row 547 RC-27 permitted them without coordination; V2-72 superseded it because the gate needs serialization. But v2, which supersedes row 547 wholesale, contains no clause about instances at all. Silence is not permission here: **unspecified**, and the gate depends on it.

**CA-09 — §4.2 has no ordering for two operations that happen today.** The table covers first append, subsequent append, correction, clear and append-after-clear. It has no row for **duplication** (P6), and none for **content appearing without an append** — which is not hypothetical: that is what probe 12 did, and what any writer bypassing the gate does.

**CA-10 — V2-13 is stronger than the store can honour.** "The reloaded marker must be the previous or the new value, never a third state and never unreadable." On the **first** marker write there is no previous document and no pre-image to recover, so a death at the replacement point leaves no marker at all (`Absent`, `no_document`). That is not "unreadable", but it is not "previous or new" either. The clause must read: previous, new, or absent — and absent is only benign when no content claims it.

**CA-11 — Row 5 and V2-32 disagree about creation.** Row 5 requires **all** representations valid and empty; V2-32 says a missing history representation is normal during creation. A `(none,0)` directory holding `history.json = []` and no `history.ndjson` satisfies neither row 5 (absent is not valid-and-empty), nor row 8 (needs a record), nor row 10 (nothing invalid), nor row 11 (needs `committed`). **Unclassified.**

---

## 3. Evidence classification correction

**V2-90 — Row 557's "both crash windows are safe" is correct only for a document that already existed.** For the first write the outcome is `Absent`, and the classification stays **Source-derived** from the store's source and its crash test — not executed in this session, and not a claim that every marker death is recoverable. VT-39's wording ("previous whole pair or new whole pair") narrows to "previous, new, or absent-with-no-pre-image".

---

## 4. Proposed clauses

**V2-91 — The decision table is the contract's only executable surface.** Any amendment that adds a precondition to a row must be written into that row. Conditions that cannot be expressed in the table are not in force.

**V2-92 — `none` with content is `Pending`, reason `unattributed`.** Never Healthy; never promoted. Content cannot create its own generation.

**V2-93 — Duplication is a marked operation.** `SaveChatAs` writes `(pending,1)` before the copy's content and `(committed,1)` after; the copy starts its own generation at 1, and never inherits the source's.

**V2-94 — An empty representation is vacuously valid.** V2-09's trailing-newline rule applies to non-empty files only; the 0-byte log of a completed clear is empty and valid, not malformed.

**V2-95 — Row 8 requires a current marker.** Presence, validity, ≥1 record, coherent coexistence **and** a marker read whose outcome is `Saved` (not `Recovered`) **and** a settled transition (V2-52, V2-57).

**V2-96 — Two distinct reasons replace one conflated one.** `retracted-residue` means content **proven** to belong to a closed generation (attributable by evidence). `unattributed` means content whose generation cannot be established. Row 6 uses `unattributed` unless provenance is proven.

**V2-97 — V2-19 is restated**: the six observed states are never stored; the marker stores intent, which is a different thing.

**V2-98 — Marker resolution and content settlement are separate authorities.** The scanner may resolve the marker under V2-56; it may never settle a content operation.

**V2-99 — Deletion is exempt from the gate, deliberately.** It destroys both representations and is an admissible exit from an unresolved transition. It is the one mutation that cannot corrupt the pair, and the exemption is stated rather than implied.

**V2-100 — Concurrent instances are specified.** The gate's critical section is cross-process (V2-76); an instance that cannot acquire it reports `Busy` and refuses, never waits unboundedly with a stale observation.

**V2-101 — Content appearing without a recorded operation is `Pending`**, reason `unattributed`, whether it arrived through a bypassing writer, a lost marker, or an interrupted clear. This is the state probe 12 produced; it is reported, never promoted.

**V2-102 — V2-13 reads: previous, new, or absent.** Absent is benign only when no content is present; with content it is V2-101.

**V2-103 — Row 5 is expressed as "every present representation is valid and empty, and no representation holds content"**, so a missing representation at `(none,0)` is a match rather than an unclassified input.

**V2-104 — This audit authorizes nothing.** No gate, CAS, marker write, duplication change, or table edit is implemented; the proposed clauses are proposals until the owner accepts them, and the table remains as row 554 wrote it.

---

## 5. Tests

- **VT-40** `(none,0)` + content → Pending (`unattributed`), never Healthy, never promoted. Replaces VT-18.
- **VT-41** A duplicated conversation is `(committed,1)` and Healthy only after its own marker write.
- **VT-42** A 0-byte log is empty and valid, not `malformed` (V2-94).
- **VT-43** `Recovered` marker + valid content → Pending, not Healthy (V2-95).
- **VT-44** Content under `cleared` without proven provenance reports `unattributed`, not `retracted-residue`.
- **VT-45** Delete is permitted during an unresolved transition and leaves no marker requiring promotion.
- **VT-46** Two instances: the second reports `Busy`, and neither writes across an unresolved transition.
- **VT-47** `(none,0)` with `history.json = []` and no log is Healthy (V2-103).

---

## 6. Standing

| Item | Status after this audit |
|---|---|
| Mutation call-graph coverage | specified, unimplemented |
| Clear retry idempotency | specified |
| Cleared-state validation | specified |
| Inventory maintenance | specified |
| Pair atomicity | Source-derived, settled under the single-document constraint |
| CAS in one lock hold | required design, unimplemented — retained, unchanged |
| Cross-generation mixture | reproduced defect |
| Detection after mixture | not reliably possible |
| Per-representation stamps | Unestablished alternative |
| **Addendum-stack consistency** | **eleven contradictions, one cause (V2-91)** |

**Verdict:** the contract is not closable as v2 plus rows 555–557. The next artefact should be a **v3 consolidation** that folds the amendments into the table inline, replaces VT-18, and adopts or rejects V2-91…V2-103 one by one. Until then the amendments are proposals, and nothing is authorized.
