# MEMORY-09 — Open Decisions (index)

**Identity:** `MEMORY-09-OPEN-DECISIONS-2026-11-15` · **Set:** artifact 4 of 5 created 2026-11-15 under owner authorisation

**Status: DRAFT — UNREVIEWED. Reconstructed from a review handoff, not from the source transcript. Textual review only; no execution, no implementation, no G4 work performed or authorised. Marked items were not supplied and must not be invented.**

**Implementation status: NOT STARTED / NOT AUTHORIZED**

---

## 1. What this document is

An **index** of items that were open when the durability contract review paused. It does not decide anything, does not replace any existing owner-decision record, and does not close any item.

Existing owner-decision records remain the authority for owner decisions: `MEMORY-09-OWNER-DECISION-RECORD-OD-M09-2026-11-13.md` and `MEMORY-09-OWNER-DECISION-ADDENDUM-1-OD-M09-15-18-2026-11-13.md`, both protected. Where an item below is also addressed there, that record governs and this index is a pointer only.

An item is open until the owner closes it in writing. Nothing in this document, and nothing in the five-file set, closes anything.

## 2. The rule that applies to every item

Every row carries the same prohibition:

> This item must not be closed by observing current implementation behaviour.

The reason is recorded in the principles reconstruction as P1 and P3: what the system currently does is an observation, and an observation does not authorise a decision about what the system should do. Deriving a contract answer from existing code would make the implementation the specification, and would silently convert every defect into a decision.

A second, related trap: an item must not be closed because a later revision appears to address it. Apparent resolution is not resolution.

## 3. Owner decisions

Closure condition for every row in this section: an owner decision, recorded in the owner-decision series. The specific condition for each item was **not supplied** to this reconstruction, so it is marked pending rather than invented.

| ID | Subject | Closure condition |
|---|---|---|
| `U-CAN-1` | `[PROSE PENDING]` | Owner decision; specific condition `[TRANSCRIPT RECONCILIATION REQUIRED]` |
| `U-CAN-2` | `[PROSE PENDING]` | As above |
| `U-CAN-3` | `[PROSE PENDING]` | As above |
| `U-A4` | `[PROSE PENDING]` | As above |
| `U-A5` | `[PROSE PENDING]` | As above |
| `U-B2` | `[PROSE PENDING]` | As above |
| `U-B6` | `[PROSE PENDING]` | As above |
| `U-ADM-1` | `[PROSE PENDING]` | As above |
| `U-ADM-2` | `[PROSE PENDING]` | As above |
| `U-ADM-3` | `[PROSE PENDING]` | As above |
| `U-ADM-4` | `[PROSE PENDING]` | As above |
| `U-ADM-5` | `[PROSE PENDING]` | As above |
| `U-ADM-6` | `[PROSE PENDING]` | As above |
| `U-ADM-7` | `[PROSE PENDING]` | As above |
| `U-ADM-8` | `[PROSE PENDING]` | As above |
| `U-ADM-9` | `[PROSE PENDING]` | As above |
| Record-ID stability | Whether an artifact's record identifier is stable across retries and generations | Owner decision; the operational requirement was discussed, the decision was not supplied |
| `Q14a` | `[PROSE PENDING]` | Owner decision; specific condition `[TRANSCRIPT RECONCILIATION REQUIRED]` |
| `Q14b` | `[PROSE PENDING]` | As above |

The prefixes are preserved exactly as supplied. No attempt has been made to infer what `CAN`, `A`, `B` or `ADM` abbreviate or group; that taxonomy is `[TRANSCRIPT RECONCILIATION REQUIRED]`.

`U-ADM-1` through `U-ADM-9` are the admission-boundary items. Beyond that grouping, no subject text was supplied.

## 4. G4 blockers

These are preconditions on any future G4 run, not decisions about the contract's text.

| ID | Blocker | Closure condition | Must not be closed by |
|---|---|---|---|
| `G07` | Runtime tokenizer unbound — the tokenizer used at run time is not bound to a recorded identity | A binding recorded such that the run's tokenizer identity is fixed and checkable. Exact requirement `[TRANSCRIPT RECONCILIATION REQUIRED]` | Inferring the tokenizer from what the loader happens to load, or from a log line |
| `G12` | Model revision and sampling unset — neither the model revision nor the sampling configuration is fixed for the run | A recorded model revision and sampling configuration fixed before the first call. Exact requirement `[TRANSCRIPT RECONCILIATION REQUIRED]` | Reading the values from a previous run, a default, or current server output |
| Retention | Sampled-conversation retention undecided | An owner decision on retention of sampled conversations | Whatever the harness currently persists |

The statements above say what each blocker is about. They do not claim the closure requirement, which was not supplied in a form that can be reproduced.

G4 remains blocked. Opening it requires closing these items and the separate authorisation that G4 has never had.

## 5. What is not open, and why that matters

Some things are often assumed still undecided and are not: the preregistration set is frozen (`MEMORY-09-PREREGISTRATION-FREEZE-2026-11-15.md`), and G4 execution, subject runs and model calls have never been authorised. Those are closed states, not open items, and this index does not reopen them.

Conversely, the contract itself is **not** a closed state. Revision 23.8 is unreviewed and unreconciled, and the items above keep it that way.

## 6. Cross-references

| Document | Relationship |
|---|---|
| `MEMORY-09-DURABILITY-CONTRACT-2026-11-15.md` | Artifact 1 of 5 — the clauses these decisions would modify |
| `MEMORY-09-GOVERNING-PRINCIPLES-2026-11-15.md` | Artifact 2 of 5 — P1 and P3 are the basis of §2 |
| `MEMORY-09-REVIEW-LEDGER-2026-11-15.md` | Artifact 3 of 5 — `L-02a`–`L-02c` remain open findings |
| `MEMORY-09-DOCUMENTATION-RECONSTRUCTION-RECORD-2026-11-15.md` | Artifact 5 of 5 — provenance and limits |
| `MEMORY-09-OWNER-DECISION-RECORD-OD-M09-2026-11-13.md` | Protected. Authority for owner decisions recorded there |
| `MEMORY-09-OWNER-DECISION-ADDENDUM-1-OD-M09-15-18-2026-11-13.md` | Protected. As above |
