# INTEGRATION REVIEW — F3C-TARGET-2 R5 and Precision Layers

**Status:** **INTEGRATION REVIEW ONLY — NOT R6, NOT FREEZE, NOT ACCEPTANCE, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**R5 proposal:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md` (`bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e`)
**Original matrix:** `docs/review/DISPOSITION-MATRIX-F3C-TARGET-2-R5-MUSE-2026-11-11.md` (`061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6`)
**Precision addendum 1:** `docs/review/ADDENDUM-PRECISION-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md` (`d059054b678f95d8697b8987bf856035a20a6415bd61828a940ebbd99d7549b0`)
**Precision addendum 2:** `docs/review/ADDENDUM-2-PRECISION-RECOVERY-F3C-TARGET-2-R5-2026-11-11.md` (`be399e5bfda3fec62e4356ad8e0dadde11a643881bd74fc314ba2bf9b116e403`)
**Latest independent review:** `docs/review/REVIEW-INDEPENDENT-ADDENDUM-2-PRECISION-RECOVERY-F3C-TARGET-2-R5-2026-11-11.md` (`d97870088f74e8c9b2e79580225f5bbb56a08cca08e1d9cf9f66e9cfdab9361d`)
**Frozen Master:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md` (`61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`)
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md` (`7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`)

**Scope:** documentation-only integration review. It consolidates the accepted findings and decisions, tests whether they compose without contradicting the frozen Master or adopted Gate v2, and identifies the remaining proof obligations. It does not modify those artifacts, write R6, adopt V18–V24, authorize implementation, or claim any Gate row passes.

**Limits:** this review is contract-level. It does not establish source-level completeness of every storage entry point, native filesystem behavior, or runtime evidence. The existing `AUDIT-F3C-SINGLE-WRITER-2026-11-10.md` is an adjacent world/scenario audit, not proof of F3-C chat-store single-writer authority.

---

## 1. Integration verdict

**The accepted architecture is coherent in direction but not yet integrated enough to authorize R6.**

The five integration findings are valid:

1. commit-record ambiguity needs a complete taxonomy;
2. incarnation identity is not storage-level exclusion;
3. gate coverage needs a closed entry-point boundary and privileged recovery class;
4. cleanup needs reference-aware retention;
5. the frozen lock hierarchy needs a whole-system proof covering recovery and destructive paths.

No accepted correction is rejected. The integration gap is that several rules are individually correct but their interaction is not yet a single authority model.

**Verdict:** keep R5 and all three precision layers preserved; keep R6 unauthorized; require the owner decisions and proof obligations below.

---

## 2. Consolidated finding and decision inventory

| Layer | IDs | Count | Current state |
|---|---|---:|---|
| R4 review | A1–A14 | 14 | Contract-level dispositions; A14 open. |
| Matrix blocking | C1–C6 | 6 | Accepted; future-R6 requirements. |
| Matrix exactness | H1–H7 | 7 | Accepted; future-R6 requirements. |
| Precision addendum 1 | P1–P6 | 6 | Accepted; OD-02a/OD-02b/OD-17 controls. |
| Precision addendum 2 | S1–S7 | 7 | Accepted; recovery/durability/enforcement controls. |
| Integration review | I1–I5 | 5 | Accepted as integration requirements; proof still pending. |

### 2.1 Open or gating items

| Item | State | Why it remains open |
|---|---|---|
| A14 / V24 | Open / blocked | Pre-amendment multi-character baseline and OD-15 decision absent. |
| C1 / OD-02 | Accepted, not resolved | OD-02a and OD-02b must be separately decided; AM-1 cannot absorb OD-02b. |
| C2 / S2 / I1 | Accepted, not resolved | Commit taxonomy and durability windows require the single state machine in §3. |
| C3 / S6 | Accepted, not resolved | Abandonment needs operation-map integration and exact winner/durability rules. |
| C4 | Accepted, not resolved | Seed golden vectors and canonical omission rules remain R6 work. |
| C5 / S5 | Accepted, not resolved | Entry-point gate must be enforced across every supported path. |
| C6 | Accepted, not resolved | One fencing inventory and recovery-incarnation rules remain. |
| H1–H7 | Accepted, not resolved | Schemas, transitions, migration, preference, sampler, archive, and A14 linkage remain R6 work. |
| I2 / OD-18 | New owner decision | No storage-level single-writer authority has been selected. |
| I4 | Accepted, not resolved | Reference-aware retention and cleanup lifecycle require a positive rule. |
| I5 | Accepted, not resolved | Whole-system lock proof is conditional pending exact recovery acquisitions. |
| OD-17 | Awaiting owner | Frozen §4.1 exhaustiveness and other layout-adjacent authority must be ruled. |
| OD-02a / OD-02b | Awaiting owner | Joint activation and durable recovery ownership unresolved. |
| V18–V24 | Proposed | Not adopted by Gate v2. |
| V20 | Blocked | F3-F evidence absent. |

### 2.2 Closed prior items

A9, A10, A12, and A13 remain closed and are not reopened by this integration review. The world/scenario single-writer audit does not become evidence for the chat-store authority model.

---

## 3. Single authoritative commit/recovery/release state machine

This is the required contract-level state machine. It is not R6 text and not an implementation.

### 3.1 States

| State | Meaning | Ordinary access |
|---|---|---|
| `S0_STABLE` | No pending recovery; current target manifest matches the last verified state. | Allowed through gate. |
| `S1_ORPHAN_PACKAGE` | Durable package exists with no authoritative record referencing it. | Not treated as a commit; cleanup only under retention. |
| `S2_RECORD_INDETERMINATE` | Record is visible but durability is unconfirmed, or record/package state is ambiguous. | Closed. Recovery or explicit block. |
| `S3_REPAIR_REQUIRED` | Authoritative record and valid package exist; live targets do not match the committed manifest or verified state is absent. | Closed. |
| `S4_REPAIRING` | Durable repair state exists; package-derived replacement is in progress. | Closed except privileged repair. |
| `S5_VERIFIED_PENDING_RELEASE` | Targets match the committed manifest; durable verified state exists; gate release not yet complete. | Closed until release. |
| `S6_RELEASED` | Manifest verified and release gate durable; ordinary access may proceed. | Allowed. |
| `S7_BLOCKED` | Torn, truncated, hash-invalid, duplicate, conflicting, missing-package, or otherwise irreducible evidence. | Hard block; no repair by guess. |
| `S8_ORPHAN_MARKER` | Durable repair/marker state exists with no authoritative record. | No commit; cleanup only under retention. |
| `D1_DELETING` | Tombstone/deletion state for the chat. | No ordinary chat access. |
| `W1_WIPING` | Root wipe in progress. | No ordinary access; destructive path only. |

### 3.2 Commit points

Three distinct points must not be conflated:

1. **Visibility / linearization:** the atomic create-only record name becomes visible to live participants.
2. **Durability:** the record, its directory entry, and the package needed for recovery are durable against process crash and host power loss.
3. **Release:** package/record/live-target verification is complete and durable verified state allows ordinary access.

The record's visibility is not the release point. A record's presence is not proof of durability. A matching target byte set is not proof of commitment.

### 3.3 Required transitions

| From | Event | To | Rule |
|---|---|---|---|
| `S0_STABLE` | New commit, package durable, record visible | `S2_RECORD_INDETERMINATE` | No ordinary reads until record durability and record/package agreement are established. |
| `S2_RECORD_INDETERMINATE` | Record and directory synchronization complete; package verified | `S3_REPAIR_REQUIRED` | Even exact target bytes do not skip verification. |
| `S2_RECORD_INDETERMINATE` | Record absent after power loss; package verified but unreferenced | `S1_ORPHAN_PACKAGE` | No commit. |
| `S2_RECORD_INDETERMINATE` | Torn, hash-invalid, duplicate, or conflicting record | `S7_BLOCKED` | Torn bytes cannot identify a package. |
| `S3_REPAIR_REQUIRED` | Durable repair marker written | `S4_REPAIRING` | Marker belongs to OD-02b. |
| `S4_REPAIRING` | Package-derived targets installed and manifest verified | `S5_VERIFIED_PENDING_RELEASE` | Do not open gate yet. |
| `S5_VERIFIED_PENDING_RELEASE` | Durable release state written | `S6_RELEASED` | Release must be serialized with successor projection. |
| `S6_RELEASED` | New commit begins projection | `S2_RECORD_INDETERMINATE` | Ordinary reads must not observe mixed successor state. |
| Any | Wipe wins exclusion | `W1_WIPING` | Recovery and wipe cannot run concurrently. |
| Any | Tombstone/delete wins | `D1_DELETING` | No ordinary chat access. |
| `S1_ORPHAN_PACKAGE` | No reference and retention permits cleanup | `S0_STABLE` or absent | Cleanup only after reference proof. |
| `S8_ORPHAN_MARKER` | No record and retention permits cleanup | `S0_STABLE` or absent | Must not be repaired as a commit. |

### 3.4 Taxonomy of record outcomes

| Record/package observation | Classification | Required outcome |
|---|---|---|
| No record; no marker; package absent | `S0_STABLE` | Prior verified state remains current. |
| No record; orphan package | `S1_ORPHAN_PACKAGE` | No commit; cleanup only under retention. |
| Record visible; durability not confirmed | `S2_RECORD_INDETERMINATE` | Closed; re-check after synchronization/recovery. |
| Record durable and verified; package valid; targets mismatch | `S3_REPAIR_REQUIRED` | Repair from package. |
| Record durable; torn/truncated/unreadable/hash-invalid | `S7_BLOCKED` | Do not infer package from torn bytes. |
| Duplicate/conflicting records for one operation | `S7_BLOCKED` | Resolve via exclusive-create and operation map; never last-writer inference. |
| Record references unreadable/missing package | `S7_BLOCKED` | Hard evidence block; no guessed bytes. |
| Record + valid package + matching targets + no durable verified state | `S5_VERIFIED_PENDING_RELEASE` | Verify manifest and install release state; do not open merely because bytes match. |
| Durable marker with no record | `S8_ORPHAN_MARKER` | No commit; retention-governed cleanup. |
| Wipe intent/root identity conflict | `W1_WIPING` or block | Wipe exclusion wins over ordinary repair. |

### 3.5 Invariants

- No ordinary read before `S6_RELEASED`.
- No state transition to `S6_RELEASED` without `S5_VERIFIED_PENDING_RELEASE`.
- No adoption of a package without a record except as an independent repair source.
- No commit inference from matching target bytes alone.
- No repair from torn record bytes.
- No package cleanup while a record, marker, history link, or pending repair references it.
- No recovery/wipe overlap.
- No state is called “released” while any required durability or verification step is unresolved.

---

## 4. Single-writer authority decision

### 4.1 The problem

An incarnation value is an identity, not an exclusion mechanism. Two processes with different incarnations can both believe they own the same root if they have different storage views or locks.

The frozen Master's locks are in-process synchronization. They do not establish cross-process or cross-mount exclusion.

### 4.2 Owner decision required

Track as **OD-18 — storage-level single-writer authority** (provisional numbering; final numbering remains the owner’s).

| Option | Meaning | Consequence |
|---|---|---|
| **A. One authority process per root, enforced by a storage-level exclusive lock** | Acquire an OS/filesystem exclusion primitive before any F3-C storage access; fail closed if not acquired. | Requires lock path, acquisition order, crash behavior, inventory/wipe treatment, and proof that the lock is not defeated by split mounts. |
| **B. Shared atomic fencing primitive** | A storage-level compare/fence operation makes the second writer fail before it acts. | Requires a concrete primitive and proof across supported filesystems; a process-local mutex is insufficient. |
| **C. Explicit bounded deployment assumption** | Only one authority process per root is supported; split mounts, snapshots, and duplicated containers are outside the supported deployment. | Requires owner acceptance of residual risk; must not be described as enforced single-writer safety. |

**Recommendation:** Option A plus a bounded deployment assumption for unsupported split-view cases. The storage-level lock path and ownership must be enumerated under OD-02b/OD-17; the lock must be acquired before any F3-C lock and must not be nested inside `chatWriteMu`.

### 4.3 Required startup order

R6 must choose and prove one order:

1. acquire storage-level exclusion;
2. allocate and durably install the new incarnation;
3. reconcile authority, leases, barriers, and old incarnation state;
4. reconcile packages, records, markers, and targets;
5. release the recovery gate.

No partial order may leave old and new authority both usable. Ambiguous incarnation state fails closed.

### 4.4 Status

**OD-18 is not decided.** Without it, the contract may prove in-process locking but cannot claim one authoritative storage writer.

---

## 5. Storage-access entry-point inventory

This is the contract-level inventory required before R6. Source-level exhaustive enumeration of every caller remains unverified.

| Class | Entry point / path | Gate class | Required lock / boundary |
|---|---|---|---|
| Chat read | `LoadChat`, `LoadChatData`, hydrated history/targets | Ordinary | Gate + `chatWriteMu` or defined read sequence; no unverified targets. |
| Chat write | `SaveChat`, `SaveChatData`, append, replace, clear | Ordinary | `chatTreeMu.R` → `stateApplyMu` → `chatWriteMu`. |
| Creation / Save As | `NewChat`, `SaveChatAs` | Ordinary mutation | Creation/journal + package/record protocol. |
| Cast | preview, confirm | Ordinary mutation | Same commit protocol; preview non-durable. |
| Delete | `DeleteChat` | Destructive | `chatSwitchMu` + tree/state/write order; tombstone. |
| Wipe | `WipeAllData` | Destructive root | `castApplyMu` → `chatSwitchMu` → `chatTreeMu.Lock` → per-chat state/write; new root identity. |
| Memory migration | canonical snapshot migration | Ordinary mutation | Same package/record/receipt/revision protocol. |
| History migration | JSON/NDJSON reconciliation | Ordinary mutation | Same package/record/receipt protocol. |
| Plan/batch/attempt/candidate | generation readers/writers | Ordinary or privileged as appropriate | `chatWriteMu`; no mutable-state bypass. |
| Startup recovery | package/record/marker/target reconciliation | Privileged recovery | Frozen lock order; tree/stable chat order; gate remains closed. |
| Gate release | verified → released | Privileged recovery/publication | Must serialize with successor projection and ordinary readers. |
| Gallery | staging, cleanup, copy, display | Separate boundary | No bypass of F3-C gate when it reads/copies chat state. |
| Attachments/media | F3-F-owned bytes, target references | F3-F / ordinary reference | F3-C does not own bytes; reference metadata still follows target rules. |
| Import | frozen Q3 browser/legacy import | Ordinary mutation | Create-only/journaled; same v2 package protocol; no silent Save As. |
| Export | chat/content copy | Ordinary or explicit privileged | Must not publish unverified state; no bypass. |
| Diagnostics | inspection, logs, support copy | Explicit privileged or read-only | No publication/copy of unverified state. |
| Background/autosave | scheduled chat writes | Ordinary | Same gate and commit proof; no “background” exemption. |
| World/scenario snapshot | `AddWorldToChat`, `AddScenarioToChat`, clear | Ordinary mutation | Part of chat target, not canonical item store. |
| Content library/worldstore | canonical world/scenario items | Outside F3-C chat store | Separate single-writer audit; not evidence for chat-store authority. |
| Browser storage | `localStorage` world/scenario records | Outside F3-C / import source | Never treated as canonical chat target. |
| Operation map/waiters | duplicate/historical/cancellation | Control path | `chatStateMu` brief leaf; no disk I/O; no upward lock. |

**Inventory finding:** every supported path needs a named gate check, lock, access class, error behavior, and source-level owner. “Normal reads” alone is not an entry-point inventory.

---

## 6. Evidence-reference and cleanup lifecycle

### 6.1 Reference graph

| Object | Created by | Referenced by | Cleanup rule |
|---|---|---|---|
| Commit record | Commit point | Package, previous record, receipt, history, recovery, current proof | Retain while chain/current proof or recovery requires; no deletion by age alone. |
| Receipt | Commit | Record, proof, historical verification | Retain with record; reconstruct only create-only. |
| Commit package | Commit staging | Record, repair, migration/export | Retain until targets verified/released and no reference/retention requirement remains. |
| Target projection | Post-commit apply | Current proof, readers | Never “verified” without manifest match and durable release. |
| Recovery marker/state | OD-02b | Record/package, gate, recovery | Retain until release and no pending repair; missing marker is interruption, not completion. |
| Plan header | User-turn plan | Batch, attempt, candidate | Immutable; clean only when no batch/attempt/candidate/history reference remains. |
| Batch | Coordinator | Attempts, candidates, operation record | Retain while nonterminal or referenced; cleanup only after reference proof. |
| Attempt | Generation | Candidate, audit | Mutable lifecycle; clean after no candidate/history/audit reference. |
| Candidate lifecycle | Generation | Commit eligibility, audit | Mutable; immutable evidence is separate. |
| Immutable A14 evidence | Generation | Durable history turn | Retain as long as cited history; hash/version-bound; wipe with linked history. |
| Predecessor archive | Migration | Migration record, export | Immutable; retain until deletion/wipe/retention rule. |
| Incarnation history | Startup | Lease authority/diagnostics | Retain only under an explicit bounded rule; no authority from stale entries. |
| Operation record | Operation map | Waiters, historical results | Retain per operation-map lifecycle; no eviction with waiters. |
| Wipe intent/journal | Wipe | Wipe recovery | Retain through phases; remove after completed new root. |

### 6.2 Positive cleanup rule

Cleanup may delete an object only after proving, under the relevant lock and authority:

- no authoritative commit record references it;
- no recovery marker, pending repair, or gate state references it;
- no history entry or A14 immutable pointer references it;
- no migration, import, export, or operation-map waiter references it;
- no retention policy requires it;
- no wipe/source boundary requires it;
- the deletion itself is recorded and recoverable.

No TTL-only deletion. No prefix-based deletion. No “tmp” deletion of a referenced object. No package cleanup before projection verification.

### 6.3 Destruction sequences that must fail closed

1. Orphan cleanup deletes a package still referenced by an authoritative record.
2. Temporary cleanup deletes a durable marker/package because it lives under a tmp-like path.
3. Partial recovery deletes a package while targets remain unverified.
4. Migration discards candidate/context evidence still cited by history.
5. Wipe omits a new OD-02b root or leaves linked evidence behind.
6. Candidate lifecycle cleanup deletes immutable A14 evidence sharing an identifier.
7. Storage pressure deletes a package needed to repair an interrupted projection.

### 6.4 Retention bound

A14 evidence must live as long as the history cites it. If history is append-only, R6 must define the owner-approved retention, export, summarization, or wipe rule that bounds growth. “Retain forever” is not a contract unless explicitly approved.

---

## 7. Lock-order proof

### 7.1 Controlling order

```text
castApplyMu
  -> chatSwitchMu
    -> chatTreeMu
      -> stateApplyMu(chat)
      -> chatWriteMu(chat)
        -> rootJournalMu

chatStateMu is a brief leaf; no disk I/O under it.
```

No new F3-C lock may be inserted. A storage-level OD-18 lock, if adopted, is an admission lock outside this order and must be acquired before any F3-C lock.

### 7.2 Path proof table

| Path | Required acquisition | Release / publication | Status |
|---|---|---|---|
| Normal save/append | tree.R → stateApply → chatWrite → rootJournal if needed | publish after durability; chatState brief | **Conditional proof:** matches frozen order. |
| Creation/Save As | tree.R → stateApply → chatWrite → rootJournal if needed | creation journal/record; publish after durability | **Conditional.** |
| Delete | cast/switch as needed → tree.R → stateApply → chatWrite | tombstone then removal; no ID reuse | **Conditional.** |
| Wipe | cast → switch → tree.Lock → per-chat state/write → rootJournal; new identity | cleanup after root commit | **Conditional:** tree.Lock excludes recovery tree.R. |
| Startup recovery | tree.R → stable per-chat order → stateApply → chatWrite → rootJournal as needed | gate release only after verification | **Proof gap:** exact acquisition sequence not yet fixed. |
| Multi-chat recovery | stable `(rootId, chatId)` order; no reverse wait | release each chat in stable order | **Proof gap:** must not hold lower while waiting higher. |
| Root journal/incarnation | rootJournal only under relevant chatWrite | durable state before publication | **Conditional.** |
| Verification helper | pure functions; no upward lock | none | **Required rule.** |
| Gate release/read | gate observation + target read under shared lock/sequence; release serialized with successor projection | ordinary access after release | **Proof gap:** exact reader lock sequence unresolved. |
| Operation map | chatState brief for registration/snapshot; no disk I/O | historical/waiter lifecycle | **Conditional.** |
| Publication | after durability; chatState leaf only | no disk under chatState | **Required rule.** |

### 7.3 No-inversion obligations

- Recovery may not hold `chatWriteMu` and wait for `stateApplyMu`, `chatTreeMu`, `chatSwitchMu`, or `castApplyMu`.
- Wipe may not hold `chatTreeMu.Lock` and wait for a recovery path holding `chatWriteMu`.
- Root journaling may not occur outside the relevant `chatWriteMu`.
- Verification may not acquire an upward lock under `chatWriteMu`.
- Multi-chat operations must use stable `(rootId, chatId)` order.
- `chatStateMu` must never be held across disk I/O or another F3-C lock.
- Wipe and startup recovery must be mutually exclusive.
- Gate release must not race successor projection.

### 7.4 Proof status

The frozen order is preserved conditionally. The lock proof is **not complete** until R6 fixes:

- startup incarnation/authority acquisition and its relation to root journals;
- exact multi-chat recovery acquisition/release;
- exact ordinary-read/gate-release serialization with successor projection;
- OD-18 storage-level lock placement and fail-closed behavior;
- wipe/recovery exclusion under crash.

No lock inversion is proved absent by this review; the obligations above are the required proof target.

---

## 8. Contradiction register against frozen Master and Gate v2

| ID | Frozen / adopted provision | Integration requirement | Status | Required resolution |
|---|---|---|---|---|
| CR-01 | Master §4.1 layout is the frozen path set. | New `f3c-packages`, `f3c-authority`, `f3c-generations`, `f3c-predecessors`, projection/repair records. | **OPEN** | OD-17 ruling; OD-02b enumeration; formal amendment if exhaustive. |
| CR-02 | Master §6.9: commit records, receipts, tombstones, targets installed create-only. | Repeated target replacement after each commit. | **OPEN** | AM-1/OD-02a; no silent reinterpretation. |
| CR-03 | Master §6.9: same-volume atomic move. | Eight target paths plus optional-file deletion. | **OPEN** | OD-02a/OD-02b per-file protocol; no multi-file atomic claim. |
| CR-04 | Master §7 locks are in-process synchronization. | Storage-level single writer across processes/mounts. | **OPEN** | OD-18 choice; bounded deployment assumption if no shared primitive. |
| CR-05 | Master §5.5/§5.6 target set and target field order. | V2 meta/history additions and new generation fields. | **OPEN** | Formal v2 compatibility/version ruling; no claim of frozen compliance yet. |
| CR-06 | Master §6.2 commit-record field list. | TARGET-2 optional fields and package reference. | **OPEN** | Formal v2 compatibility/version ruling. |
| CR-07 | Master §6.9 retained evidence while a chain depends on it. | Reference-aware cleanup, retained packages, A14 evidence. | **OPEN** | Positive retention rule; no age/prefix cleanup. |
| CR-08 | Gate v2 adopts V1–V17 only. | V18–V24 proposed. | **COMPATIBLE — HOLD** | Do not adopt or claim PASS without a new Gate record. |
| CR-09 | V24 is not adopted; A14 open. | Multi-character baseline/context mode. | **COMPATIBLE — OPEN** | Keep V24 blocked; capture baseline. |
| CR-10 | `AUDIT-F3C-SINGLE-WRITER` covers world/scenario items, not chat storage. | Chat-store single-writer authority. | **SCOPE MISMATCH** | Do not cite the audit as F3-C chat-store proof. |
| CR-11 | Master §8.10 wipe inventory. | New OD-02b roots, markers, projection state. | **OPEN** | Add all new paths to wipe inventory; no stale evidence. |
| CR-12 | Master §6.9 temp leftovers cleaned without inventory visibility. | Durable markers/packages may share tmp-like storage. | **OPEN** | Never delete referenced evidence by path prefix or age. |

No contradiction is closed by this review. Items marked compatible are boundaries to preserve, not approvals.

---

## 9. Required owner decisions

| Decision | Status | Blocks |
|---|---|---|
| OD-02a — per-file target replacement mechanics | Awaiting owner | Any target write; inert without OD-02b. |
| OD-02b — durable recovery/layout/mutable evidence | Awaiting owner | Any target replacement and recovery state. |
| OD-17 — frozen §4.1 exhaustiveness and related layout authority | Awaiting owner | New roots, OD-02b, R6. |
| OD-18 — storage-level single-writer authority | Awaiting owner | Any claim of cross-process authority. |
| OD-04 — lease timing | Proposed | Time-bound lease behavior. |
| OD-06 — plan/batch storage | Proposed | Batch/plan lifecycle. |
| OD-12 — exact schemas | Proposed | Record/state machine implementation. |
| OD-15 — multi-character context mode | Proposed | A14/V24. |
| OD-14 — V18–V24 adoption | Awaiting owner | Gate additions. |
| OD-01 — package retention/GC | Proposed | Package cleanup. |

No decision is inferred from this review.

---

## 10. Status and authorization

**Integration review:** **DOCUMENTATION ONLY; NOT R6; NOT A CONTRACT REVISION.**
**R5:** **PROPOSAL ONLY; NOT FROZEN; NOT ACCEPTED; NOT IMPLEMENTED; NOT VERIFIED.**
**Original matrix:** preserved.
**Precision addendum 1:** preserved.
**Precision addendum 2:** preserved.
**Latest independent review:** accepted as findings; not approval.
**R6:** **NOT AUTHORIZED; NOT WRITTEN.**
**F3-C phase 0:** **NOT CLOSED.**
**P1a:** **UNAUTHORIZED**, downstream of P1c.
**F3-C acceptance:** **NOT GRANTED.**
**F3-D:** **BLOCKED.**
**A9/A10/A12/A13:** closed; not reopened.
**V20:** **BLOCKED** pending F3-F evidence.
**Gate v2:** adopted for V1–V17 only; V18–V24 remain proposed.
**A14/V24:** open / blocked.

**Required before R6:**

1. Muse reviews this integration review.
2. Owner rules on OD-17 and OD-18.
3. Owner decides OD-02a and OD-02b separately and confirms their joint activation.
4. Owner accepts or replaces the commit/recovery/release state machine.
5. Owner accepts the entry-point inventory and evidence-retention lifecycle.
6. R6 proves the lock-order obligations and closes the contradiction register.

No source, test, build, migration, model load, destructive test, application data, authority record, freeze, acceptance, or implementation action was performed. No R6 was drafted.

---

## Appendix A — Reviewed artifacts

| Artifact | SHA-256 |
|---|---|
| R5 proposal | `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e` |
| Original matrix | `061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6` |
| Precision addendum 1 | `d059054b678f95d8697b8987bf856035a20a6415bd61828a940ebbd99d7549b0` |
| Precision addendum 2 | `be399e5bfda3fec62e4356ad8e0dadde11a643881bd74fc314ba2bf9b116e403` |
| Latest independent review | `d97870088f74e8c9b2e79580225f5bbb56a08cca08e1d9cf9f66e9cfdab9361d` |
| Frozen Master | `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31` |
| Gate v2 | `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a` |

This review's own SHA-256 is recorded in the completion report rather than embedded here, to avoid a self-referential file hash.
