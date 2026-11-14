# FREEZE — F3-C Architecture and Implementation Contract

**Date:** 2026-11-11  
**Status:** **FROZEN — ARCHITECTURE ONLY.**  
**Authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Supporting authority:** frozen F3-C decisions, LF-1, A1–A13, and accepted native evidence.  
**Implementation-verification protocol:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md` (Gate v2, adopted).

This record freezes the architecture so implementation can be measured against a stable contract. It does **not** accept F3-C. F3-C acceptance requires implementation and evidence. It does not authorize P1a implementation, source changes, or tests.

---

## 1. Freeze declaration

The consolidated F3-C contract is frozen as the candidate implementation authority. The isolated Master-1 through Master-3c amendment sequence is retained as decision history and superseded wording, but the consolidated contract controls where they differ.

Frozen scope:

- chat identity, root incarnation, revisions;
- commit records, receipts, successor links, fork inventories;
- canonicalization, request fingerprints, target bytes, `tgt1-` hashes;
- target set and ownership;
- locks, publication, operation map;
- creation, Save As, append, autosave, switch, delete, cast;
- ordinary/cast journal recovery;
- tombstone integrity;
- memory and history migration;
- wipe staging, tree hashing, gallery boundary, root identity commit;
- implementation-verification protocol, defect taxonomy, and status ladder (Gate v2).

Not frozen as accepted behavior:

- F3-C implementation completeness;
- P1a authorization;
- F3-D;
- cast confirmation implementation;
- A9/A10/A12 evidence for the new persistence contract.

---

## 2. Consistency audit

### 2.1 Frozen six decisions

| Decision | Result | Note |
|---|---|---|
| F3-C-Q1 folder per item | PASS | Contract does not alter worldstore item paths. |
| F3-C-Q2 stable identity | PASS | Opaque identity, no title authorization. |
| F3-C-Q3 explicit journaled import | PASS | No automatic import or fallback added. |
| F3-C-Q4 pinned snapshots | PASS | No implicit refresh; chat snapshots remain explicit. |
| F3-C-Q5 preserve unavailable | PASS | No substitution or name fallback. |
| F3-C-Q6 stable lore ownership | PASS | World/scenario lore remains item-owned; chat copies are not canonical lore. |

### 2.2 LF-1

| Requirement | Result |
|---|---|
| Legacy name folders never read at runtime | PASS — no runtime reader added. |
| Never a fallback | PASS — canonical store by ID only. |
| Never written or deleted | PASS — no write target added. |
| One explicit journaled migration | PASS — outside this chat persistence contract. |
| Unresolved content quarantined | PASS — no guessing or title matching. |

### 2.3 A1–A13

This table checks design consistency against each gate; it does not re-accept evidence.

| Gate | Result | Note |
|---|---|---|
| A1 identity/paths | PASS | No conflict. |
| A2 lossless update | PASS | No conflict. |
| A3 roles | PASS | No conflict. |
| A4 revision | PASS | Chat revision is distinct from `rv1-` document revision. |
| A5 references | PASS | ID-only; no title fallback. |
| A6 generated lore | PASS | Chat-owned snapshot; no item-folder writes. |
| A7 snapshots | PASS | No implicit refresh. |
| A8 media | PASS | No media ownership change. |
| A9 import | PASS | LF-1 evidence accepted; no contract conflict. |
| A10 single writer | PASS | Contract forbids chat item-folder writes. |
| A11 availability | PASS | No persistence exception introduced. |
| A12 native prompt boundary | PASS | Accepted native evidence. |
| A13 regression | PASS | Passing regression gate. |

### 2.4 Accepted native evidence

- A9 PASS: explicit import, LF-1 migration, quarantine, rollback, idempotency, source-tree preservation.
- A10 PASS: native chat load/save, canonical item trees unchanged, generated-lore snapshot behavior.
- A12 PASS: production prompt path at inference boundary.
- A13: passing regression gate.

Evidence limitation: A9/A10/A12 predate the consolidated persistence contract. They do **not** prove the new target set, `memories_snapshot.json` writer, revisions, receipts, or migration code. They prove the accepted import, single-writer, and prompt boundaries.

### 2.5 Final consistency points

| Check | Result |
|---|---|
| Memory migration receipt and journal agreement | PASS — receipt carries the pre-commit `memoryMigrationJournalHash`; final journal carries that pre-commit hash plus commit record/receipt hashes; missing/corrupt journal is an evidence gap, not committed. |
| Wipe root-committing boundary | PASS — staging states and `preCommitGate` are preconditions; the new root identity file remains the only commit point. |
| Target ownership | PASS — `memories_snapshot.json` required; `memories.json` outside; `lorebooks.json` chat-scoped; canonical lore item-owned. |
| Creation identity | PASS — intent, journal, chat ID, root, operation, fingerprint, and source binding must agree before advancement. |
| History identity | PASS — order and IDs preserved; no index-derived IDs. |
| Gallery boundary | PASS — separate staging state; ambiguous crash states block; no hidden partial move. |

---

## 3. Decision history

| Revision | Accepted correction |
|---|---|
| Master-1 | Initial consolidated identity, revisions, receipts, locks, recovery, cast boundary. |
| Master-2 | Receipt inventory, append/dirty, publication ordering, Save As, operation identity, deletion, journal recovery. |
| Master-3 | Prepared-journal non-adoption, exact target hashing, autosave session revalidation, domain-separated paths, wipe/tombstone integrity. |
| Master-3a | `memories_snapshot.json` correction, creation ordering, wipe staging completeness, target evidence. |
| Master-3b | Memory conflict policy, deterministic tree hash, gallery boundary, ownership matrix, creation identity. |
| Master-3c | Memory migration commit boundary, gallery crash reconciliation, history identity preservation. |

---

## 4. Supersession register

| Superseded wording | Current authority |
|---|---|
| `memories.json` as canonical chat target | `memories_snapshot.json` required; `memories.json` legacy/separate. |
| `prepared + all-candidate` may adopt | `prepared` never adopts; conflict. |
| Creation directory may exist while `prepared` | `writing` durable before directory creation. |
| Autosave checks only root/base | Also session, epoch, poison, discard generation. |
| Plain `sha256(rootID||chatID)` paths | Domain-separated path hashes. |
| Wipe without staging gate | Tree-hash staging gate before `root-committing`. |
| Gallery covered by app root | Separate gallery boundary and crash state. |
| History migration may renumber IDs | Preserve order and IDs; reject missing/duplicate identities. |
| Receipt alone proves memory migration | Receipt and migration journal must agree. |
| Target manifest without exact target bytes | `F3C-TARGET-1` and `tgt1-` exact-byte contract. |

Superseded wording remains in the review history. The consolidated contract controls.

---

## 5. Authorization boundary

Frozen:

- architecture and implementation contract;
- decision inheritance;
- target ownership;
- recovery and failure outcomes;
- evidence requirements;
- Gate v2 verification method and status separation.

Not authorized:

- P1a implementation;
- source changes;
- tests;
- deployment;
- F3-D;
- claiming F3-C acceptance.

The next permissible step is a separate implementation-planning authorization. Implementation must follow the frozen dependency order and Gate v2, including the fault-injection, race, filesystem recovery, native Wails/WebView2, and staged Mimo challenge gates.

---

## 6. Freeze verdict

**Architecture freeze: GRANTED.**  
**Gate v2: ADOPTED — does not grant verification or acceptance.**  
**F3-C acceptance: NOT GRANTED.**  
**P1a implementation: NOT AUTHORIZED.**  
**F3-D: BLOCKED.**  
**A9/A10/A12/A13: CLOSED.**

Gate v2 keeps four statuses separate: implementation complete, verification passed, F3-C accepted, and F3-D authorized. This record grants none of them.
