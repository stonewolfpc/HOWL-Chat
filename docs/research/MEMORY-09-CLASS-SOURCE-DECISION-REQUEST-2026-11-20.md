# MEMORY-09 — Class-Source Decision Request

## Where does `class_baseline` come from?

**Status:** **owner decision request.** It does not decide the source, authorize code, or change the implementation
plan. Implementation authorization remains valid, but the class-dependent path is **blocked** until this is resolved.
No schema, code, migration, test, requirement or canonical capture.

| Field | Value |
|---|---|
| Date | **2026-11-20** (local, CST) |
| Gap | `MemoryRecord` has no `KnowledgeClass`; `classify.go` has no class rule; no default is permitted |
| Decision requested | the source of `class_baseline`, including how auto-captured memories obtain a class |
| Recommendation | **author-supplied at creation/review** (Option A), with absent class rejected |
| Blocked path | automatic capture producing L-2 formation entries |
| Preserved | JSON per-dimension architecture; legacy separation; AC-11 tests |
| Tests executed | **0** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished**.

---

## 1. The gap

**CSR-01 (Observed — no field).** `internal/backend/types/memory.go` `MemoryRecord` has `Scope` and legacy `Status`, but
no `KnowledgeClass`. `MemoryCandidate` has no class field either.

**CSR-02 (Observed — no rule).** `internal/backend/memory/classify.go` assigns scope, owner, title, summary, content,
tags and `pending` status. It does not assign a knowledge class. A code search finds no
`knowledge_class`/`KnowledgeClass`.

**CSR-03 (Calculated — the class baseline has no source).** L-2 requires a `class_baseline` in the formation entry.
Adding the field is mechanical; supplying its value is not. The current code cannot do it.

**CSR-04 (Observed — the prohibitions).** The implementation must not:

- infer class from **scope**;
- infer class from **event type** without an owner-approved mapping;
- assign a convenient **default**;
- treat an **absent** class as an ordinary valid class.

**CSR-05 (Calculated — this is semantic, not a missing field).** A class is an authoring decision about what kind of
knowledge the memory is. A field with no source would either fail validation or invite one of the prohibited guesses.

---

## 2. Vocabulary

**CSR-06 (Source-derived — the existing set).** r33 §5 gives the knowledge-class values:
`ambient | restricted | experiential | character_specific`. The implementation has no field, but the model’s
vocabulary is available. This request does not invent new values.

**CSR-07 (Calculated — the decision should confirm the vocabulary).** Whichever source is chosen should supply one of
the existing values. If the owner wants a different vocabulary, that is a separate model decision, not an
implementation detail.

---

## 3. Options

### Option A — author-supplied at creation/review

**CSR-08 (Calculated).** `AddMemory` requires a class from the caller. Auto-capture still creates **candidates**, but a
candidate has no class until a reviewer supplies one; approval requires the class. No formation entry is written for an
auto-captured memory until approval.

- **Strength:** no inference, no default, no absent-as-valid; aligns with “knowledge class is authored” (I13); fits the
  existing `pending | approved | rejected` review workflow.
- **Cost:** a class input is needed in the AddMemory API and candidate review path. Without a UI field, auto-capture
  approval is blocked; that is a UI/API dependency, not a semantic guess.
- **Effect on auto-capture:** candidates still queue; L-2 formation waits for review.

### Option B — owner-approved event-type mapping

**CSR-09 (Calculated).** A fixed, owner-approved table maps `EventType` to a class (for example, micro-memories to
`character_specific`, discovery/location to `experiential`). The classifier applies only that table.

- **Strength:** auto-capture can form L-2 entries without review; the mapping is explicit and auditable.
- **Cost:** event type does not carry the `ambient`/`restricted` distinction. A discovery can be public or secret; a
  death can be world fact or character-specific. A mapping would often guess context the event type does not contain.
- **Effect on auto-capture:** allowed, but only for classes the mapping can justify.

### Option C — defer auto-capture L-2; manual only

**CSR-10 (Calculated).** Manual `AddMemory` requires a class; auto-capture continues to produce candidates but never
creates a formation entry until a separate auto-capture decision.

- **Strength:** safest; no inference; no UI dependency for manual records.
- **Cost:** auto-captured memories remain outside the L-2 lifecycle; the feature is incomplete until a later decision.

### Option D — unresolved class slot

**CSR-11 (Calculated — excluded).** A `class_slot` analogous to `status_slot` would preserve the obligation without a
value. But Decision A requires `knowledge_class` to be reconstructible, and the owner authorized an unresolved slot
only for **status** pending N-b. A class slot would defer the class obligation without that authorization. **Excluded.**

### Option E — default / scope inference / unapproved event-type inference

**CSR-12 (Calculated — excluded).** Directly violates CSR-04. No default, no scope inference, no unapproved mapping.

---

## 4. Comparison

**CSR-13 (Calculated).**

| Criterion | A — author-supplied | B — approved mapping | C — defer auto L-2 |
|---|---|---|---|
| no scope inference | **yes** | **yes** | **yes** |
| no unapproved event-type inference | **yes** | **owner-approved only** | **yes** |
| no default | **yes** | **yes** | **yes** |
| absent class rejected | **yes** | **yes** | **yes** |
| auto-capture forms L-2 | **no** (waits for review) | **yes** | **no** |
| semantic fidelity | **high** — author decides | **medium** — event type is coarse | **high** — no auto claim |
| UI/API dependency | class input at review | mapping table | manual only |
| fits existing review workflow | **yes** | no review needed | **yes** |
| risk of wrong class | low | **medium–high** | low |

**CSR-14 (Calculated — A is the least inventive).** A does not guess. It uses the existing review step to obtain the
author’s class. B is permissible only with an explicit owner-approved table, and the table cannot recover context the
event type lacks. C is safe but leaves auto-capture outside L-2.

---

## 5. Recommendation

**CSR-15 (Calculated — recommend A).** The owner should decide that `class_baseline` is **author-supplied** at
`AddMemory` and at candidate approval; an absent class is rejected. Auto-capture may continue to produce **candidates**
but must not write an L-2 formation entry until a reviewer supplies the class.

**CSR-16 (Calculated — why A over B).** The event taxonomy is about what happened, not about whether the knowledge is
ambient, restricted, experiential or character-specific. A mapping would encode a guess as a rule. If the owner prefers
B, the mapping must be explicitly approved and versioned; this request does not supply one.

**CSR-17 (Calculated — the UI/API dependency).** A requires a class input in the manual-add API and the candidate
review path. If the current UI cannot supply it, the implementation stops for a UI/API decision. That is a surface
dependency, not a reason to invent a class.

---

## 6. Preserved architecture constraints

**CSR-18 (Observed — JSON is compatible).** The separate `formations`, `scope_changes` and `class_changes` collections
can implement L-2/P-1. They must not create a universal ordering rule:

- formation supplies the baseline for each dimension;
- scope versions order only scope changes;
- class versions order only class changes;
- no ordering relationship is inferred between the two collections.

**CSR-19 (Observed — legacy handling).** Incomplete legacy records must remain distinguishable from corrupt records.
Their existing application behavior must not silently change merely because they lack the new lifecycle representation.
The legacy `Status` field remains separate from the lifecycle status slot.

**CSR-20 (Observed — AC-11 tests).** The status round trip must be tested end to end:

| Condition | Expected result |
|---|---|
| explicit unresolved | preserved as unresolved |
| missing slot | validation failure |
| unknown state | invalid |
| missing resolved value | invalid |
| corrupt lifecycle file | corruption reported |
| legacy `active` value | never interpreted as resolved lifecycle status |

**CSR-21 (Calculated — these are not reopened by the class decision).** The class source changes what supplies
`class_baseline`; it does not change the JSON ordering constraint, legacy separation or status-slot tests.

---

## 7. What the decision would and would not decide

**CSR-22 (Calculated — would decide).**

- the source of `class_baseline` for new records;
- how auto-captured memories obtain a class;
- that an absent class is rejected, not defaulted;
- whether auto-capture may form L-2 entries before review.

**CSR-23 (Calculated — would not decide).**

| Not decided | Where it belongs |
|---|---|
| physical field names/encoding | implementation plan |
| UI layout/class picker | UI work; separate |
| N-b status value | N-b |
| N-c general lifecycle order | N-c |
| formation authority | separate |
| Decision B | separate |
| legacy backfill migration | separate decision |
| requirement/canonical capture | separate gates |

**CSR-24 (Calculated — the implementation remains blocked).** The class-dependent path cannot pass AC-2/AC-3 until the
source is chosen. The rest of the plan may be structured, but no auto-capture L-2 entry may be written, and no
acceptance criterion may be claimed as passed.

---

## 8. Decision requested

**CSR-25 (Observed — the owner’s choice).**

| Option | Meaning |
|---|---|
| **A** | author-supplied at AddMemory and candidate approval; absent class rejected; auto-capture waits for review |
| **B** | owner-approved event-type mapping; no default; mapping explicit and versioned |
| **C** | manual only; auto-capture L-2 deferred |
| **Other** | specify the source; it must satisfy CSR-04 |

**CSR-26 (Observed — approval is the decision).** This request does not choose. The owner’s selection becomes the
class-source decision; implementation then follows the selected path.

---

## 9. Traceability

**CSR-27 (Observed).**

| What | Where |
|---|---|
| Implementation plan | IPL-05/06/09/33 |
| Class gap | `internal/backend/types/memory.go`; `classify.go` |
| L-2 class baseline | L2D-07/08 |
| Decision A class obligation | DAO-01/05 |
| Status safeguard / AC-11 | IAA-03; IAR-14 |
| JSON ordering constraint | IPL-13/14/24; CSR-18 |
| Legacy separation | IPL-16/17; CSR-19 |
| Class vocabulary | r33 §5 `:318` |
| No default/inference | CSR-04 |

---

## 10. Record status

**CSR-28 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — map **118** |
| Class source selected | **0** — request only |
| Implementation started | **0** |
| Tests | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **decision request** — no canonical capture |

**End of request.** The class baseline has **no source** today (CSR-01/02/03). The implementation must not infer it from
scope, infer it from event type without an approved mapping, default it, or accept an absent value (CSR-04). The
recommended source is **author-supplied at creation/review** (CSR-15), with auto-capture waiting for review; the
approved-mapping option remains available if the owner prefers it (CSR-09/13). The JSON ordering constraint, legacy
separation and AC-11 tests are preserved (CSR-18/19/20). **No code, schema, migration or test follows from this
request** (CSR-28).
