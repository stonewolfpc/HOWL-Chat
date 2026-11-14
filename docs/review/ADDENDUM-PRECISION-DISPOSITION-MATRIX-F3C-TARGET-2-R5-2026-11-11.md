# ADDENDUM — Precision Corrections to the R5 Disposition Matrix

**Status:** **DOCUMENTATION-ONLY ADDENDUM — NOT R6, NOT A CONTRACT REVISION, NOT FREEZE, NOT ACCEPTANCE, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Supplements:** `docs/review/DISPOSITION-MATRIX-F3C-TARGET-2-R5-MUSE-2026-11-11.md`
**Matrix SHA-256:** `061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6`
**Matrix size:** 279 lines; 22,248 bytes.
**Independent matrix review:** `docs/review/REVIEW-INDEPENDENT-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md`
**Matrix review SHA-256:** `e28e7aa0c81f2bfd739bf396f6e28dfbbd57bc2f28a00575c6302c6200c9a9bb`
**Matrix review size:** 359 lines; 18,866 bytes.
**R5 proposal:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md` (`bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e`)
**Frozen Master:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md` (`61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`)
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md` (`7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`)

**Scope:** this addendum preserves the original matrix and all existing authority hashes. It records the six precision corrections as accepted, tightens the OD-02a/OD-02b boundary, requires an explicit frozen §4.1 ruling, states the clock/incarnation/recovery-gate questions that R6 must answer, defines abandonment as an ordered durable transaction, and requires immutable hash-bound A14 context evidence. It does not write R6, authorize R6, implement anything, run tests, adopt V18–V24, freeze F3-C, or close phase 0.

---

## 1. Disposition of the six corrections

All six corrections are accepted. None reopens the original C1–C6 or H1–H7 findings; they tighten the adjudication record that will guide a future R6.

| Correction | Disposition | Effect on the matrix |
|---|---|---|
| **1. OD-02 boundary** | **Accept** | OD-02a is reduced to per-file replacement mechanics. OD-02b exclusively owns durable projection status, repair markers, read barriers, and recovery state. A joint interface is required. |
| **2. Lease clock** | **Accept** | Monotonic time alone is not sufficient. R6 must define incarnation authority, persisted-monotonic prohibitions, suspend/snapshot behavior, one expiry boundary, and the atomic-install I/O qualification. |
| **3. Startup gate** | **Accept** | The recovery gate must be an enforceable storage-access boundary covering every supported reader and mutation path, not a startup UI convention. Ambiguous state fails closed. |
| **4. Frozen §4.1** | **Accept** | The owner must expressly rule whether frozen §4.1 is exhaustive before OD-02b or R6. No new root may be introduced by implication. |
| **5. Abandonment ordering** | **Accept** | Abandonment becomes an ordered durable transaction under `chatWriteMu`, with explicit serialization against candidate commit and explicit revision/history consequences. |
| **6. A14 evidence pointer** | **Accept** | A bare mutable candidate ID is insufficient. Durable history must embed the observed context head or bind to immutable candidate evidence by hash/version. |

No correction is rejected. The matrix remains the R5 adjudication authority as corrected by this addendum.

---

## 2. Revised OD-02a / OD-02b boundary

### 2.1 OD-02a — replacement mechanics only

OD-02a may authorize only:

- byte-verified replacement of individual current target files from an already durable, verified commit package;
- per-file same-volume replacement, with prior target bytes retained in the package;
- verification that each package-derived file matches the package manifest and target hash;
- package-derived optional-file deletion, performed idempotently;
- post-commit application of target bytes after the commit record is installed.

OD-02a must **not** authorize:

- any durable projection-status file, recovery marker, read-barrier state, repair state, or durable “repair needed” flag;
- any read blocking or unblocking decision that survives a process restart;
- any new root, journal, authority record, lease, barrier, plan, batch, attempt, candidate, or mutable operation-evidence class;
- any multi-file atomic-replacement claim;
- replacement of commit records, receipts, tombstones, packages, archives, plans, attempts, candidates, batches, authority records, leases, or barriers.

The matrix’s former OD-02a bullet “projection-status tracking, read blocking until repair” is withdrawn. If that tracking affects restart behavior, it is durable recovery state and belongs exclusively to OD-02b.

### 2.2 OD-02b — durable coordination and recovery state

OD-02b exclusively owns:

- every new durable root, including package, authority, generation, and predecessor storage;
- every durable projection-status record, repair marker, recovery barrier, or durable read-gate state;
- all mutable authority, lease, barrier, batch, attempt, and candidate lifecycle records;
- the state machine, write protocol, recovery behavior, inventory treatment, and wipe treatment for those records;
- the durable gate that makes an incomplete multi-file replacement unobservable to supported readers and mutations.

OD-02b must state, for each durable record:

- its exact path and path-hash domain;
- its owner and allowed writers;
- its permitted and forbidden states and transitions;
- its write, flush, install, and verification rules;
- its behavior when missing, torn, duplicated, stale, corrupt, or inconsistent;
- its interaction with ordinary lease barriers and wipe barriers;
- its inventory and wipe treatment;
- which bytes are hashed, which fields are excluded from identity, and which fields are audit-only.

### 2.3 Joint interface

The two decisions must define one joint interface:

- **OD-02a** defines what may be written to target paths, how package-derived replacement is verified, and the per-file ordering.
- **OD-02b** defines the durable gate that records the replacement as incomplete, repairing, verified, or released, and blocks access until verification succeeds.
- A target replacement may not be treated as complete merely because the per-file writes returned success.
- A repair marker may not be written under OD-02a authority.
- The record installation remains the single commit point. Projection completion is a post-commit recovery/apply condition, not a second commit point.

The required sequence is:

```text
verified package
  -> final lease/authority proof
  -> atomic create-only record installation (commit point)
  -> durable projection state says repairing/incomplete
  -> per-file replacement under OD-02a mechanics
  -> durable manifest verification
  -> durable projection state says verified/released
  -> supported readers and mutations may proceed
```

No durable `repairing` or `verified` state exists under OD-02a alone.

---

## 3. Required frozen §4.1 ruling

Before OD-02b can be approved, and before R6 can be commissioned, the owner must expressly rule whether frozen Master §4.1 is:

### Option A — exhaustive

The frozen layout is the complete allowed persistent root set. Any new root requires its own formal layout amendment. OD-02b is not sufficient by itself; AM-1 cannot be used as the layout amendment.

### Option B — non-exhaustive

The frozen layout permits additional owned paths under explicit contract rules. The owner must still record that ruling, and OD-02b must enumerate every new path, owner, transition, recovery rule, inventory rule, and wipe rule. “Unspecified” does not mean “authorized.”

### If the ruling is absent or ambiguous

- No new root may be created.
- OD-02b is not approved.
- R6 is not authorized.
- No implementation may infer permission from R5, the matrix, or this addendum.

Track this ruling as **OD-17 — frozen layout exhaustiveness** for review purposes. Final decision numbering remains the owner’s.

---

## 4. Clock, incarnation, and commit-boundary questions

The matrix’s monotonic-clock direction is retained, but it is not yet a proven recovery rule. R6 must answer the following explicitly.

### 4.1 Monotonic source and discontinuity

R6 must name the exact monotonic clock source and state its guarantees across:

- normal process execution;
- sleep and hibernation;
- host suspend and resume;
- virtual-machine migration;
- snapshot restore;
- container pause/resume;
- NTP or wall-clock adjustment;
- process restart.

If the source does not preserve elapsed-time authority across an event, that event must be treated as a process/incarnation change for lease authority. The design may not rely on a frozen, paused, repeated, or rewound monotonic reading.

### 4.2 Persisted monotonic values

Persisted `issuedMono`, `deadlineMono`, or equivalent values may be retained only for diagnosis. They must never:

- be compared with a new process’s monotonic clock;
- determine remaining lifetime after restart;
- override incarnation fencing;
- authorize a commit.

Wall-clock timestamps are audit evidence only. A wall-clock observation may trigger investigation; it does not become commit authority.

### 4.3 Current-incarnation record

R6 must define a durable current-incarnation authority record. It must specify:

- where the record lives and which root/chat scope owns it;
- how the current incarnation is allocated, uniqueness-checked, and written;
- the update order relative to authority/lease startup reconciliation;
- crash behavior before, during, and after the update;
- behavior when the record is missing, torn, duplicated, or inconsistent;
- inventory and wipe treatment;
- how a clean prior shutdown is distinguished from a crash.

Ambiguous incarnation state must fail closed. Recovery may not rely on process memory alone.

### 4.4 One expiry boundary

The matrix’s mixed `>=`/`>` comparison is not uniform. R6 must select one explicit boundary convention for both ordinary and maximum expiry. Recommended convention: a lease is invalid at `nowMono >= deadlineMono` and at `nowMono >= maxDeadlineMono`; equality always expires.

Golden lease vectors must cover:

- immediately before expiry;
- exact equality;
- immediately after expiry;
- exact maximum lifetime;
- restart with an old incarnation;
- suspend/snapshot discontinuity.

### 4.5 Atomic-install I/O qualification

R6 must replace the absolute “no I/O” wording with a durability-primitive rule:

- after final validation, no application or filesystem reads/writes may occur except the specified atomic create-only record installation and its required file and directory durability operations;
- define exclusive creation, file synchronization, directory synchronization, ordering, and crash behavior for that installation;
- distinguish validation reads from post-validation staging, mutation, or projection I/O;
- define whether the record installation is a same-directory atomic rename and what happens on collision, torn temporary state, or unsupported filesystem behavior.

The commit point remains the atomic installation of the record. Target projection and receipt installation remain post-commit apply work and cannot retroactively authorize an expired lease.

---

## 5. Startup/read gate

The gate is not a startup notification and not an initial UI check. It is an enforceable boundary at the storage access layer.

### 5.1 Covered paths

R6 must enumerate at least:

- ordinary chat history and target readers;
- metadata, preference, cast, continuation, and projection readers;
- plan, batch, attempt, and candidate readers;
- migration source selection;
- wipe source selection and inventory readers;
- recovery readers themselves;
- gallery, attachment, world, scenario, and imported-content paths where applicable;
- background, autosave, diagnostic, export, or copy paths capable of reading or publishing chat state.

An out-of-contract direct filesystem read cannot be physically prevented, but every supported entry point must pass through the same gate. No background path may bypass it because it is “only diagnostic.”

### 5.2 Fail-closed recovery state

The gate must treat missing, corrupt, torn, stale, duplicated, or inconsistent recovery/projection state as **unverified**:

- unverified current targets are not current;
- no supported reader may publish or copy them;
- migration and wipe may not select them as sources;
- no new commit may proceed;
- recovery must re-derive repair work from the verified package;
- successful repair must be followed by manifest verification before the gate opens;
- a crash after repair but before gate release must cause re-verification, not an assumption that repair completed.

### 5.3 Recovery-barrier states

The durable recovery barrier belongs under OD-02b and H2. R6 must define:

- creation;
- states such as open, repairing, verified, released, and blocked;
- permitted and forbidden transitions;
- operations allowed in each state;
- crash behavior for every transition;
- interaction with ordinary lease barriers and wipe barriers;
- inventory and wipe treatment.

The matrix’s illustrative names are not approved vocabulary. R6 must define the authoritative state table from scratch and demonstrate that every state has a contractual meaning.

---

## 6. Abandonment as an ordered durable transaction

`AbandonBatchSlot` must be a durable operation, not a flag that other components are expected to notice.

R6 must define:

- the operation identity, plan hash, batch ID, participant index, expected slot state, reason, and actor;
- the lock and commit order, using `chatWriteMu` as the serialization point;
- whether the operation increments `fenceEpoch`, writes a slot-scoped revocation, or does both;
- whether it advances ordinary chat revision/history or only batch-control state;
- whether it produces a commit record, a batch-operation record, or another equally durable operation record;
- how it interacts with recovery, inventory, and wipe;
- how it fails closed when the plan, batch, slot, or current state is ambiguous.

Required ordering against a candidate:

1. Both abandonment and candidate commit serialize under the same `chatWriteMu`.
2. If abandonment wins before the candidate’s final gate, the candidate fails with a slot-terminal/abandonment error and no record is installed.
3. If the candidate record is installed first, the turn is committed; abandonment may not rewrite or invalidate that committed turn.
4. If abandonment arrives while the candidate is `commit_pending` but not yet installed, the final gate must reject the candidate or the abandonment must fail; neither may silently overwrite the other.
5. Abandonment must not become a second informal commit path with weaker ordering than ordinary mutations.
6. No abandoned slot may later commit through retry, lease replay, recovery, or a new process using old holder state.

The revision/history consequence is an owner decision. If abandonment writes any durable state, it must use an explicit operation record and the same recovery discipline as other durable mutations.

---

## 7. A14 evidence must be immutable or embedded

The H7 correction is accepted with this constraint:

- **Preferred:** embed `contextHistoryHead` directly in the durable v2 `turn`.
- **Alternative:** bind a pointer to immutable candidate evidence using a candidate version, package hash, `cand3_` fingerprint, or equivalent immutable identity.

A bare `candidateId` or `attemptId` is insufficient because candidate records are mutable. A pointer must not derive its evidentiary meaning from a mutable lifecycle record.

The history package must answer, without consulting mutable generation state:

- what history head the character actually observed;
- which plan and participant slot produced the turn;
- which generation request and candidate identity were committed;
- whether the captured context diverged from persisted prefix order.

A14 remains open. The pre-amendment multi-character baseline, OD-15 choice, and V24 evidence remain missing. **V24 remains blocked.**

---

## 8. Matrix preservation and precedence

- The original matrix is preserved unchanged at `061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6`.
- This addendum corrects only the precision points above. It does not rewrite the matrix’s C1–C6, H1–H7, A14, V18–V24, or owner-decision inventory.
- Where this addendum and the matrix differ, this addendum controls the precision point; the matrix remains the adjudication record.
- Illustrative status names in the matrix are not approved vocabulary.
- No matrix disposition is converted into a freeze claim.

---

## 9. Status and authorization

**Addendum:** **DOCUMENTATION ONLY; NOT R6; NOT A CONTRACT REVISION.**
**R5:** **PROPOSAL ONLY; NOT FROZEN; NOT ACCEPTED; NOT IMPLEMENTED; NOT VERIFIED.**
**Matrix:** preserved; precision-corrected by this addendum.
**R6:** **NOT AUTHORIZED; NOT WRITTEN.**
**F3-C phase 0:** **NOT CLOSED.**
**P1a:** **UNAUTHORIZED**, downstream of P1c.
**F3-C acceptance:** **NOT GRANTED.**
**F3-D:** **BLOCKED.**
**A9/A10/A12/A13:** closed; not reopened.
**V20:** **BLOCKED** pending F3-F evidence.
**Gate v2:** adopted for V1–V17 only; V18–V24 remain proposed.

**Required before R6:**

1. Muse reviews this addendum.
2. Owner accepts OD-02a and OD-02b separately.
3. Owner issues the OD-17 frozen-layout exhaustiveness ruling.
4. Owner accepts or replaces the clock, incarnation, commit-boundary, and startup-gate requirements.
5. Owner confirms abandonment ordering and the A14 immutable-evidence rule.
6. A14 remains open and V24 remains blocked.

No source, test, build, migration, model load, destructive test, application data, authority record, freeze, acceptance, or implementation action was performed. No R6 was drafted.

---

## Appendix A — Identity

| Artifact | SHA-256 |
|---|---|
| Original matrix | `061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6` |
| Independent matrix review | `e28e7aa0c81f2bfd739bf396f6e28dfbbd57bc2f28a00575c6302c6200c9a9bb` |
| R5 proposal | `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e` |
| Frozen Master | `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31` |
| Gate v2 | `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a` |

This addendum’s own SHA-256 is recorded in the completion report rather than embedded here, to avoid a self-referential file hash.
