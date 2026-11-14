# REVIEW — Independent Review of the R5 Disposition Matrix

**Status:** **INDEPENDENT REVIEW ONLY — NOT APPROVAL, NOT FREEZE, NOT R6, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Reviewed matrix:** `docs/review/DISPOSITION-MATRIX-F3C-TARGET-2-R5-MUSE-2026-11-11.md`
**Matrix SHA-256:** `061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6`
**Matrix size:** 279 lines; 22,248 bytes.
**R5 proposal:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md`
**R5 SHA-256:** `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e`
**Prior R5 review:** `docs/review/REVIEW-INDEPENDENT-F3C-TARGET-2-R5-2026-11-11.md`
**Prior R5 review SHA-256:** `a8441f9dc9b30eddfc9493f11d4dca2d0252da3f89bab735bedef957d4499c2d`
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md`
**Gate v2 SHA-256:** `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`
**Scope:** documentation-only independent review. No source, test, build, migration, model-load, destructive-test, application-data, authority-record, freeze, acceptance, P1a, R6, or F3-D action was performed or authorized. No existing file was modified.

---

## 1. Verdict

**Preserve the matrix as the R5 adjudication authority, subject to the precision corrections below. Do not commission R6 yet.**

The matrix correctly:

- accepts C1–C6 as freeze-blocking;
- accepts H1–H7 as future-R6 exactness requirements;
- keeps A14 open;
- keeps V18–V24 proposed and V24 blocked;
- separates authorization from post-commit application;
- connects batch-slot abandonment to fencing;
- refuses multi-file atomicity;
- refuses to infer an R6, freeze, acceptance, implementation, or Gate adoption.

The two precision concerns identified for this review are valid and require tightening. The matrix also contains three additional boundary defects concerning OD-02a/OD-02b, monotonic-clock semantics, and the startup recovery barrier.

None of these is a reason to reject the matrix. They are reasons not to convert its otherwise sound dispositions directly into contract language without correction.

---

## 2. Artifact-integrity check

The matrix’s Appendix A correctly records the four related artifacts as they presently exist:

- R5 proposal hash matches the computed R5 bytes.
- Prior independent R5 review hash matches the computed review bytes.
- Frozen Master hash matches the computed frozen-contract bytes.
- Gate v2 hash matches the computed Gate bytes.

The matrix does not embed its own hash, stating that the value belongs in the completion report to avoid a self-referential hash. Its current computed hash is recorded above.

This verifies present bytes and stated identities. It does not establish historical authorship, repository provenance, or that no earlier matrix revision existed.

---

## 3. OD-02a/OD-02b: the split is correct, but the boundary leaks

The matrix is right to separate narrow target-projection replacement from new roots and mutable operation evidence. That separation prevents AM-1 from becoming a broad exception to frozen §4.1 and §6.9.

Its present OD-02a/OD-02b boundary nevertheless assigns projection-status responsibility to both sides.

Matrix §5.1 permits OD-02a to authorize:

- target-manifest verification;
- projection-status tracking;
- read blocking until repair;
- recovery from the package.

Matrix §5.2 and §7 then require OD-02b to authorize:

- any durable projection-status or repair marker;
- every mutable record’s path, transitions, write protocol, recovery behavior, inventory treatment, and wipe treatment.

A read barrier that survives a crash cannot be implemented through behavior alone. It needs durable shared state. Therefore, OD-02a cannot coherently authorize “read blocking until repair” while OD-02b separately owns “any projection-status or repair marker.”

**Concrete boundary failure:**

1. A future R6 implements status tracking only in memory under OD-02a.
2. It replaces targets after a durable commit.
3. The process crashes before projection completes.
4. Restart loses the in-memory status.
5. An ordinary reader observes a mixture of old and new target files.
6. Alternatively, R6 writes a durable marker under OD-02a, thereby exercising OD-02b authority without an OD-02b decision.

Both outcomes defeat the intended narrow scope.

**Required correction:**

- OD-02a may authorize only the mechanics of replacing individual target files from an already verified package.
- Every durable projection-status record, recovery marker, barrier, or repair state belongs exclusively to OD-02b.
- The two decisions must define one joint interface: OD-02a defines what may be written and how package-derived replacement is verified; OD-02b defines the durable gate that makes incomplete replacement unobservable.
- Matrix §5.1 must not use “projection-status tracking” without saying whether that tracking is ephemeral or durable. If it affects restart behavior, it is durable and belongs under OD-02b.

---

## 4. Lease clock: the matrix improves R5 but remains imprecise across restarts

The matrix correctly moves lease authority away from blind cross-restart timestamp comparison. It states that:

- monotonic time governs lifetime while the process is alive;
- persisted wall-clock timestamps are audit-only;
- the prior process incarnation is not reusable;
- existing leases become orphaned rather than trusted across a process boundary.

That direction is sound. It is not yet a complete recovery rule.

### 4.1 Persisted monotonic values must never become authority

The matrix says no persisted monotonic value is trusted across a process boundary. A future R6 must enforce that as a structural rule, not merely a startup convention:

- persisted `issuedMono`, `deadlineMono`, or equivalent values may be retained for diagnosis;
- they must not be compared with a new process’s monotonic clock;
- recovery must not calculate remaining lifetime from them;
- a mismatch between old wall-clock audit values and new process state must not override incarnation fencing.

**Concrete sequence:**

1. Process P1 starts, allocates incarnation I1, and issues a lease.
2. P1 persists wall-clock audit timestamps.
3. P1 crashes.
4. Process P2 starts with unrelated monotonic origin and incarnation I2.
5. P2 reads the old durable lease.
6. If P2 compares any old monotonic deadline with its own clock, the comparison is meaningless.
7. The only safe outcome is to orphan the I1 lease because I1 is not the current incarnation.

The matrix states the correct outcome but does not yet require R6 to define where the current incarnation is durably recorded and how an ambiguous incarnation record fails closed.

### 4.2 The current-incarnation record is missing

Orphaning “all leases from a previous incarnation” presupposes that startup can reliably identify the current incarnation and distinguish it from:

- a clean prior shutdown;
- a crashed prior incarnation;
- a corrupt incarnation record;
- two ostensibly current processes caused by a lock, filesystem, or deployment failure.

The matrix requires a fresh 128-bit random value recorded under `chatWriteMu`, but does not specify:

- the durable record containing the current incarnation;
- its update order relative to lease startup reconciliation;
- its crash behavior;
- its inventory and wipe treatment;
- the behavior when it is missing, torn, duplicated, or otherwise ambiguous.

**Required correction:** define the current-incarnation authority record and require ambiguous incarnation state to fail closed. Recovery cannot rely on memory alone.

### 4.3 Suspend, virtualization, and clock discontinuity remain open

A monotonic clock is not automatically safe across suspend, hibernation, virtual-machine migration, snapshot restore, or container pause/resume. The matrix does not name an exact clock source or define these cases.

**Concrete concern:**

1. A holder obtains a lease with a monotonic deadline.
2. The host suspends or restores the process from a snapshot.
3. The monotonic source pauses, repeats, jumps, or otherwise loses its elapsed-time guarantee.
4. Wall-clock audit values show that the nominal maximum lifetime has passed.
5. The resumed process treats its frozen or rewound monotonic reading as authority.
6. A stale holder commits after the intended maximum lifetime.

The matrix should require R6 to:

- name the exact monotonic source;
- state its guarantees across sleep, suspend, snapshot restore, virtualization, and NTP adjustment;
- treat loss of those guarantees as equivalent to a process restart or incarnation change where necessary;
- define whether wall-clock audit evidence may trigger investigation without itself becoming commit authority.

### 4.4 The “uniform” expiry comparison is not uniform

Matrix §6.1 states expiry as:

```text
nowMono >= deadlineMono
```

or:

```text
nowMono > maxDeadlineMono
```

fails.

That is not a uniform boundary. The first deadline is exclusive at equality; the maximum deadline is inclusive at equality. A contract should not call mixed inclusive/exclusive comparisons uniform.

**Required correction:** select one explicit boundary convention for both deadlines or explain why they intentionally differ. Golden lease-lifetime vectors should cover exact-deadline and exact-maximum cases.

### 4.5 No-I/O wording needs a durability-primitive qualification

The matrix correctly separates:

- package construction, flushing, read-back, and verification before the final gate;
- atomic record installation at the commit point;
- target projection and receipt installation afterward.

It then says no I/O may occur between the final proof and record installation. Read literally, record installation itself is I/O and required durability flushes are I/O.

The intended rule is sound, but R6 needs precise language:

- after validation is complete, perform no additional application or filesystem reads/writes other than the specified atomic create-only record installation and its required file/directory durability operations;
- define exclusive creation, file synchronization, directory synchronization, ordering, and crash behavior for that installation;
- distinguish validation reads from post-validation mutation or staging I/O.

---

## 5. Startup/read gate: required, but not yet placed at an enforceable boundary

The matrix correctly requires normal reads, publication, migration, wipe-source selection, and further commits to wait for startup repair and manifest verification. The user’s second precision concern is therefore well founded: this must be enforced at the actual shared access boundary.

The matrix does not yet define that boundary.

### 5.1 Every covered reader must be enumerated

“Normal reads” is not an enforceable category unless R6 lists the covered access paths. At minimum, R6 must identify:

- ordinary chat history and target readers;
- metadata, preference, cast, continuation, and projection readers;
- plan, batch, attempt, and candidate readers;
- migration source selection;
- wipe source selection and inventory readers;
- recovery readers themselves;
- gallery, attachment, world, scenario, or imported-content paths where applicable;
- any background or diagnostic reader capable of publishing or copying chat state.

An out-of-contract direct filesystem read cannot be physically prevented, but every supported entry point must pass through the same gate. The matrix should expressly require that enumeration.

### 5.2 Ambiguous recovery state must default to closed

The matrix requires blocking until repair completes, but does not explicitly state the default when the repair marker itself is missing, torn, duplicated, stale, or inconsistent with the package/record state.

**Concrete sequence:**

1. A commit record is installed.
2. The process crashes before projection completes and before—or during—a repair-marker write.
3. Restart finds either no marker or a torn marker.
4. A reader interprets the absence of a marker as “no repair needed.”
5. It observes mixed old/new targets or selects them as a migration/wipe source.

**Required correction:**

- missing, corrupt, stale, or inconsistent recovery/projection state must mean “unverified”;
- unverified current targets must not be treated as current;
- recovery must re-derive repair work from the verified package;
- successful repair must be followed by manifest verification before the gate opens;
- a crash after repair but before gate release must cause re-verification, not an assumption that prior repair completed.

### 5.3 Recovery-barrier transitions need the same exactness as leases

Because the recovery/projection barrier must be durable, it belongs under OD-02b and H2. R6 must define:

- barrier creation;
- open, repairing, verified, released, and blocked states;
- permitted and forbidden transitions;
- operations allowed in each state;
- crash behavior for every transition;
- inventory and wipe treatment;
- interaction with ordinary lease barriers and wipe barriers.

The matrix’s illustrative barrier vocabulary must not be mistaken for an approved state machine. It is a requirement for R6 to create one.

---

## 6. Remaining matrix-level corrections

### 6.1 Decide whether frozen §4.1 is exhaustive

Matrix §8 correctly identifies the issue: if frozen §4.1 is exhaustive, the new roots require their own amendment. It leaves the interpretive question conditional.

For adjudication purposes, that condition should become an explicit owner ruling before R6. A future R6 should not be written while it remains unclear whether:

- `f3c-packages/`;
- `f3c-authority/`;
- `f3c-generations/`;
- `f3c-predecessors/`;
- projection/repair markers

are amendments to an exhaustive frozen layout or merely unspecified paths under a non-exhaustive layout.

### 6.2 Abandonment ordering needs a durable-operation rule

The matrix correctly requires R6 to decide whether abandonment is itself a durable commit and how it is ordered against an in-flight candidate.

That decision should also specify:

- whether abandonment advances ordinary revision/history or only batch-control state;
- whether an abandonment operation is subject to the same universal commit proof;
- how concurrent abandonment and candidate-commit operations are serialized under the same lock;
- how abandonment interacts with fence increments, recovery, inventory, and wipe.

Abandonment must not become a second informal commit path with weaker ordering than ordinary mutations.

### 6.3 H7’s candidate pointer must be hash-bound or embedded

The matrix permits either embedding `contextHistoryHead` in durable history or using a normative pointer to immutable candidate evidence. The second alternative needs an additional constraint.

If candidate records have mutable lifecycle state, a bare `candidateId` or `attemptId` pointer is insufficient. R6 must either:

- embed the observed context head directly in the durable v2 turn; or
- use a pointer bound to an immutable candidate version, package hash, fingerprint, or similarly immutable identity.

Otherwise, mutable generation state could alter the evidentiary meaning of an already committed history entry.

### 6.4 Illustrative new statuses must not become silent approvals

Matrix §7 uses lifecycle terms such as:

- `issued`;
- `committing`;
- `blocked-timeout`;
- `completed`;
- expanded lease/barrier transitions.

Those terms are useful as R6 requirements, but they are not part of R5 and must not be treated as already approved vocabulary. R6 must define the authoritative status tables from scratch and demonstrate that every listed matrix state has a contractual definition.

---

## 7. Disposition-by-disposition assessment

- **C1:** accept and split. The OD-02a/OD-02b separation is necessary, but §3’s projection-status boundary correction is required.
- **C2:** accept. The commit-point/final-gate direction is correct, but §4’s clock, incarnation-record, suspend, expiry-boundary, and durability-primitive corrections are required.
- **C3:** accept. The abandonment-to-fencing connection is necessary, but §6.2’s durable-operation and ordering corrections are required.
- **C4:** accept. The recommended omit-unavailable-seeds rule is sound, but R6 must also normalize terminology and add exact golden vectors.
- **C5:** accept. The refusal to overclaim atomicity is correct, but §5’s enforceable read/recovery gate corrections are required.
- **C6:** accept. A unified fencing-event inventory is necessary, but recovery-incarnation and batch invalidation need exact allocation and transition rules.
- **H1:** accept. The complete commit proof must expressly include abandonment, batch/slot linkage, and recovery scope.
- **H2:** accept. No mutable lifecycle state should remain prose-only.
- **H3:** accept. Migration needs a complete initialization mapping, not merely an archive-preservation rule.
- **H4:** accept. Every accepted preference change must use the ordinary durable-transaction protocol.
- **H5:** accept. Engine-opaque sampler order and thinking-budget behavior must narrow, not overstate, the effective-request claim.
- **H6:** accept. Archive, authority, lease, barrier, plan, batch, attempt, and candidate identity require explicit canonical-byte and hash-domain rules.
- **H7:** accept with §6.3’s pointer correction. A mutable-state-dependent A14 pointer would be inadequate.
- **A14:** remains open. The matrix correctly refuses closure and correctly blocks V24.
- **V18–V24:** correctly remain proposed. No Gate row passes.
- **A9/A10/A12/A13:** the matrix correctly does not reopen them by implication.

---

## 8. Authorization boundary

This review:

- does not approve the matrix as a contract revision;
- does not approve R5;
- does not authorize R6;
- does not adopt V18–V24;
- does not freeze or accept F3-C;
- does not authorize P1a or implementation;
- does not unblock F3-D;
- does not modify R5, the prior independent review, the matrix, frozen Master, Gate v2, or any source/test artifact.

Recommended next control remains:

1. Apply the precision corrections in §§3–6 to the adjudication record.
2. Obtain separate owner decisions for OD-02a and OD-02b.
3. Obtain an explicit ruling on whether frozen §4.1 is exhaustive.
4. Approve or replace the final commit-boundary, clock, heartbeat, incarnation, and startup-gate rules.
5. Keep A14 open and V24 blocked.
6. Only then consider whether to commission an R6.

The matrix is suitable to preserve as adjudication authority after those corrections. It is not yet suitable to use as direct R6 contract text.
