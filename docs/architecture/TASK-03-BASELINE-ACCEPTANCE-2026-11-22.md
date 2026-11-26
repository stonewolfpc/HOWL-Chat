# TASK-03 — V3.1 accepted as the contract baseline

| Field | Value |
|---|---|
| Status | **Baseline acceptance record. OQ-6 and OQ-7 open. No implementation authorized.** |
| Subject | `TASK-03-RECONCILIATION-CONTRACT-V3-1-2026-11-22.md` (row 560), hash `b7421bd051b3…` |
| Recovery token | `TASK-03-BASELINE-ACCEPTANCE-2026-11-22` |

---

## 1. What is accepted

**A-01 — V3.1 is the current contract baseline.** It is the single authoritative ordered decision procedure for TASK-03. Rows 547–560 remain historical evidence and are not required to execute it.

**A-02 — Acceptance is of the *contract*, not of any mechanism in it.** The gate, CAS, cross-process serialization, `gate_enforced`, the settlement record and the enforcement record are all specified and none is implemented.

**A-03 — No implementation authorization follows.** No code, no migration, no marker creation, no gate activation, no deletion enablement, no duplication change, no persistence change.

**A-04 — Accepted with two questions explicitly open**, not as defects to reopen the table: **OQ-6** (settlement evidence) and **OQ-7** (enforcement evidence). Their absence is a known limitation the table already pays for in rows 18–19.

## 2. Two records, two claims, no cross-certification

**A-05 — OQ-6 and OQ-7 stay separate.** A settlement record must not certify enforcement, and an enforcement record must not certify that earlier generations were cleanly settled:

| Record | Establishes | Cannot establish |
|---|---|---|
| **OQ-6 — settlement** | a generation transition completed | that the writer serialised mutations |
| **OQ-7 — enforcement** | the mutation boundary is enforced | that any past transition was settled |

**A-06 — Neither may be derived from the other**, and neither may be inferred from the absence of observed contamination.

**A-07 — Each needs its own crash-consistency and lifecycle definition**, or the ambiguity simply moves from the conversation marker into the new evidence mechanism.

## 3. Standing

**A-08 — Next: OQ-6**, because it addresses the historical-availability limitation. OQ-7 remains open and is not investigated here.

**A-09 — Rows 547–560 remain byte-identical.** This record changes no clause of V3.1.
