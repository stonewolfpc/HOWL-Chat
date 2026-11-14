# REVIEW — Independent Review of Precision Addendum 2: Recovery, Durability, and Enforcement

**Status:** **INDEPENDENT REVIEW ONLY — NOT APPROVAL, NOT FREEZE, NOT R6, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Reviewed addendum:** `docs/review/ADDENDUM-2-PRECISION-RECOVERY-F3C-TARGET-2-R5-2026-11-11.md`
**Addendum 2 SHA-256:** `be399e5bfda3fec62e4356ad8e0dadde11a643881bd74fc314ba2bf9b116e403`
**Addendum 2 size:** 489 lines; 23,070 bytes.
**First precision addendum:** `docs/review/ADDENDUM-PRECISION-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md`
**First precision addendum SHA-256:** `d059054b678f95d8697b8987bf856035a20a6415bd61828a940ebbd99d7549b0`
**Independent precision review:** `docs/review/REVIEW-INDEPENDENT-ADDENDUM-PRECISION-F3C-TARGET-2-R5-2026-11-11.md`
**Independent precision review SHA-256:** `46033f073db1739a4f6c9b1e92295ab632e081b727af2c89dce1eba3fa7783b5`
**Original matrix:** `docs/review/DISPOSITION-MATRIX-F3C-TARGET-2-R5-MUSE-2026-11-11.md`
**Original matrix SHA-256:** `061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6`
**R5 proposal:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md`
**R5 SHA-256:** `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e`
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md`
**Gate v2 SHA-256:** `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`
**Scope:** documentation-only independent review. No source, test, build, migration, model-load, destructive-test, application-data, authority-record, freeze, acceptance, P1a, R6, or F3-D action was performed or authorized. No existing file was modified.

---

## 1. Verdict

**Preserve Addendum 2 as the recovery/durability/enforcement precision layer, subject to the integration requirements below. Do not commission R6 yet.**

Addendum 2 correctly converts the prior precision findings into lifecycle rules:

- OD-02a alone authorizes no target write;
- OD-02b alone authorizes no target replacement;
- only their joint transaction may replace targets;
- a missing `repairing` marker is interruption, not completion;
- atomic visibility and power-loss durability are separate guarantees;
- ordinary access cannot bypass an incomplete recovery barrier;
- abandonment has a deterministic serialization winner;
- A14 evidence needs immutability and retention linkage.

The four difficult areas identified for this review are handled in the right direction, but the rules do not yet fully prove that they work together. In particular, Addendum 2 still needs:

1. a complete committed/uncommitted/indeterminate taxonomy for interrupted records;
2. a storage-level single-authority rule capable of defeating two live processes;
3. a closed supported-entry-point boundary with privileged recovery access;
4. a reference-aware retention/cleanup rule linking records, packages, markers, history, and A14 evidence;
5. an explicit frozen-lock-hierarchy proof covering startup recovery, root journals, wipe exclusion, publication, and operation-map integration.

None of these rejects Addendum 2. They are the remaining integration proof without which R6 would still be a collection of individually reasonable rules.

---

## 2. Artifact-integrity check

Addendum 2’s Appendix A correctly records:

- original matrix hash `061b...`;
- first precision addendum hash `d059...`;
- second precision review hash `4603...`;
- R5 hash `bbe9...`;
- frozen-contract hash `6104...`;
- Gate v2 hash `7fdc...`.

Those values match the current bytes of the corresponding files. Addendum 2 does not embed its own hash. Its current computed hash is recorded above.

This verifies present bytes and stated identities. It does not establish historical authorship or repository provenance.

---

## 3. Commit-record ambiguity: committed, uncommitted, and indeterminate

Addendum 2 §3–§4 correctly separates:

- record present with valid package: authoritative after verification;
- record present without valid package: hard evidence block;
- record absent with orphan package: no commit;
- visible-but-not-yet-synchronized record: not yet power-loss durable.

R6 still needs a complete interruption taxonomy because “record present” is not a binary filesystem observation. At minimum, R6 must define:

- absent record;
- visible but not yet synchronized;
- present and synchronized;
- present but torn, truncated, unreadable, or hash-invalid;
- duplicate/conflicting records for the same operation;
- present record whose package reference is unreadable;
- present record and valid package with ambiguous marker state;
- present record and valid package with exactly matching live targets but no durable verified state.

The most important unresolved case is a torn or hash-invalid record. Addendum 2 says a record without a valid package is a hard block, but a torn record may not reliably identify any package. R6 must not infer a package reference from torn record bytes.

**Required rule:**

- unreadable, torn, truncated, or hash-invalid record bytes are not an authoritative commit;
- they are an evidence block until reconciled;
- recovery may use an orphaned package only as a repair source after its independent package-hash verification;
- recovery must not adopt targets from an orphan alone;
- duplicate records for one operation must resolve through exclusive-create/collision semantics and the operation map, not through last-writer inference.

R6 must also clarify the “exact match without durable verified state” case. If live targets already match the committed manifest but no verified/released marker exists, the safe outcome is to verify the manifest and install durable verified state. Wholesale file replacement may be unnecessary, but the gate must not open merely because the bytes happen to match. Verification and gate installation remain required.

Finally, package durability must be proved before the final authority gate, not merely described as “written and flushed.” R6 must specify exact file synchronization, directory synchronization, read-back verification, create-only installation, and collision behavior for the package itself. Otherwise, a durable record could reference a non-durable package in the power-loss window Addendum 2 is meant to close.

---

## 4. Multi-process authority: incarnation fencing needs a storage-level rule

Addendum 2 §5 correctly requires:

- a durable current-incarnation record;
- exact uniqueness semantics;
- no authority from shutdown markers;
- durable validating heartbeats;
- observable and fail-closed behavior for duplicate processes, split mounts, failed locking, snapshots, and duplication.

That is necessary but not sufficient unless R6 states the storage-level exclusion mechanism.

An incarnation number, by itself, does not prevent two processes from believing they own the same root. If two processes do not share the same lock, each can allocate an incarnation, write its own “current” record to its own view of storage, and proceed. The later restart may detect inconsistency, but detection after divergent writes is not prevention.

**Concrete failure:**

1. Storage is mounted twice or a VM/container is duplicated.
2. Process P1 allocates incarnation I1 and writes it through its storage view.
3. Process P2 allocates incarnation I2 and writes it through another storage view.
4. Both hold what they believe is the only `chatWriteMu`.
5. Both issue leases and install records.
6. Their storage views later reconcile or one overwrites the other.
7. Incarnation fencing detects the conflict only after two incompatible authorities have acted.

R6 must therefore define one of the following, explicitly:

- a supported deployment model with exactly one authority process per root storage; plus the storage/filesystem primitive used to enforce it; or
- a shared atomic fencing primitive capable of making the second writer fail before it acts; or
- a clearly bounded residual risk, owned by the owner, where split storage is outside the supported deployment.

A purely in-process mutex cannot satisfy the second alternative. If the implementation uses a filesystem lock, R6 must define its path, acquisition order, crash behavior, split-mount behavior, and interaction with the frozen lock hierarchy. If split mounts can defeat that lock, R6 must say so rather than imply that incarnation numbers alone solve multi-writer storage.

Related startup ordering also needs precision. R6 must define whether startup:

1. allocates and durably installs the new incarnation;
2. reconciles authority, leases, barriers, and old incarnation state;
3. performs package/record/target reconciliation;
4. opens the gate.

If old leases are orphaned before the new incarnation is durable, a crash can leave neither the old nor the new authority well defined. If the new incarnation is installed before reconciliation, R6 must prove that the partially reconciled state remains closed. Either order can work, but only one exact order may be approved.

---

## 5. Gate enforcement: every supported path must use the same boundary

Addendum 2 §6 correctly requires two access classes and atomic binding of gate observation to the protected action. It also correctly refuses to let diagnostic, export, gallery, attachment, background, or copy paths bypass the gate.

The remaining integration defect is that Addendum 2 enumerates paths without proving that they all pass through the same protected boundary.

R6 must provide a supported-entry-point inventory and, for each entry, state:

- which gate is checked;
- which lock protects the check and subsequent action;
- whether the accessor is ordinary or privileged recovery;
- whether target hashes are verified on read or only during recovery;
- what error is returned for unverified state;
- how migration/wipe-source selection is serialized against repair and gate release.

Two omissions are especially important.

First, privileged recovery access must itself operate under the frozen lock order. “Under the recovery locks and authority” is not a lock-order proof. R6 must state the exact acquisition sequence for startup repair, whether it acquires `stateApplyMu`, how it uses `rootJournalMu` for root-journaled markers, and how multi-chat recovery preserves stable `(rootId, chatId)` order.

Second, gate release and successor-commit admission must be atomic with respect to ordinary readers. Otherwise:

1. recovery verifies targets and releases the gate;
2. an ordinary reader checks the released gate;
3. a new commit begins projection before the reader opens the files;
4. the reader observes files from two revisions despite having passed the gate.

The addendum identifies this hazard. R6 must eliminate it through shared-lock serialization, not merely describe it as a requirement.

Out-of-contract direct filesystem reads cannot be physically prevented, but R6 must make the supported/unsupported distinction explicit. Any tool that bypasses the gate must be unsupported for correctness purposes, even if it remains useful for forensics.

---

## 6. Evidence retention: pointers must not dangle

Addendum 2 §§3, 7, and 8 correctly require:

- package availability until projection is verified/released;
- no package cleanup before completion;
- immutable A14 evidence separate from mutable lifecycle state;
- hash/version-bound history pointers;
- retention linkage and consistent wipe behavior.

The missing integration rule is reference-aware cleanup. R6 must define the authoritative reference set consulted before deleting any package, record, marker, plan, batch, attempt, candidate, incarnation-history, or A14 evidence object.

**Concrete destruction sequences that must fail closed:**

1. Orphan-package cleanup deletes a package still referenced by an authoritative commit record.
2. Temporary-file cleanup deletes a durable marker or package because it lives under a tmp-like path.
3. Partial recovery deletes an “old” package while its targets remain unverified.
4. Migration builds v2 targets but discards candidate/context evidence still needed by history.
5. Wipe inventory omits a new OD-02b root and leaves retention-linked evidence behind—or deletes history while retaining evidence it was meant to remove consistently.
6. Candidate-lifecycle cleanup deletes immutable A14 evidence because both share an identifier or storage prefix.
7. Storage pressure causes background deletion of a package needed to repair an interrupted projection.

“A package retained only by a manifest digest is not sufficient” is the right warning. R6 must replace it with a positive rule: cleanup may delete an object only after proving, under the relevant lock, that no authoritative record, recovery marker, history entry, retention link, or pending repair references it.

Retention also needs a bounded lifecycle. If A14 evidence must live as long as history cites it, and history is append-only, R6 must define the owner-approved retention, export, summarization, or wipe rule that prevents unbounded growth. “Retain forever” is not a contract unless the owner explicitly approves it.

---

## 7. Frozen lock hierarchy: no inversion may be introduced

Addendum 2 does not overtly contradict the frozen lock hierarchy, but several of its requirements touch locks closely enough that R6 needs an explicit proof. The controlling order remains:

```text
castApplyMu
  -> chatSwitchMu
    -> chatTreeMu
      -> stateApplyMu(chat)
      -> chatWriteMu(chat)
        -> rootJournalMu
```

with `chatStateMu` as a brief leaf and no disk I/O under it.

R6 must prove at least the following.

### 7.1 Startup recovery uses the same order

Startup repair is not exempt from lock ordering because it runs at startup. If it acquires per-chat locks for multiple chats, it must use stable `(rootId, chatId)` order. If it inspects or repairs root-wide state, it must define how `rootJournalMu`, incarnation records, barriers, and per-chat locks interact.

A recovery path that acquires a per-chat write lock and then waits for root-wide wipe or barrier completion would invert the intended hierarchy: wipe acquires the tree lock and then per-chat locks. R6 must ensure recovery never holds lower-order locks while waiting for higher-order completion.

### 7.2 Root journals remain under the relevant write lock

Addendum 2 introduces or relies on several potentially root-scoped records:

- current incarnation;
- recovery barriers;
- durable projection markers;
- abandonment operation records;
- immutable evidence indexes.

If any of these is root-journaled, `rootJournalMu` must be acquired only under the relevant `chatWriteMu`, under the frozen rule. R6 must identify which records are root-journaled and which are chat-scoped. It must not allow a root-wide recovery record to become a second lock hierarchy.

### 7.3 Verification helpers must not acquire upward

The final commit proof occurs under `chatWriteMu`. Frozen rules prohibit verification helpers from acquiring `stateApplyMu` under `chatWriteMu`.

R6 must therefore ensure that final-gate validation uses only:

- state protected by already-held higher-order locks;
- `chatWriteMu`-protected authority/lease/batch/marker records;
- already durable package/record bytes;
- pure verification functions that acquire no F3-C lock.

This is especially important because the final gate is now larger: operation, lease, incarnation, session, epoch, cast, fence, plan, batch, slot, abandonment, base/history, package/manifest/record, and recovery scope. Expanding validation under `chatWriteMu` is acceptable only if none of the added checks acquires an upward lock.

### 7.4 Publication remains after durability

Gate release, publication, incarnation publication, and ordinary-read availability must follow the frozen publication discipline:

- durable state first;
- brief `chatStateMu` use only for registration/snapshot/publication;
- no disk I/O under `chatStateMu`;
- no additional F3-C lock acquired while holding `chatStateMu`.

A recovery gate that performs filesystem writes while holding publication state, or that publishes availability before directory synchronization completes, would create a new ordering defect even if every individual filesystem operation were correct.

### 7.5 Operation-map integration

Abandonment is now a durable operation. It must integrate with the frozen operation map:

- registration/cancellation semantics;
- duplicate waiter behavior;
- in-flight/completed lifecycle;
- historical-result handling;
- wipe invalidation.

Duplicate abandonment or candidate waiters must not acquire `stateApplyMu` or otherwise bypass the operation-ownership rules. The deterministic winner rule governs execution ordering; the operation map governs duplicate and historical visibility.

### 7.6 Wipe exclusion

Startup recovery and wipe must be mutually exclusive. R6 must define whether:

- startup recovery runs before wipe can begin;
- wipe initiation aborts or waits for startup reconciliation;
- a crash during wipe ветра resumes wipe rather than ordinary package repair;
- wipe inventory includes every OD-02b package, marker, authority, incarnation, batch, attempt, candidate, and retention-linked evidence path.

Wipe is not an ordinary commit, so “no new commit may proceed” is insufficient. Recovery and wipe need an explicit exclusion rule.

---

## 8. Disposition-by-section assessment

- **§1:** accept. All seven prior findings are correctly preserved as requirements; none is rejected.
- **§2:** accept, with §4’s activation condition. Neither OD-02a nor OD-02b alone authorizes replacement, and their joint lock discipline must be exact.
- **§3:** accept, with §3’s missing-marker reconciliation and torn-record taxonomy added.
- **§4:** accept, with §§3 and 6’s package-durability, visibility/durability, platform, and ambiguous-window rules completed.
- **§5:** accept, with §4’s storage-level single-authority, incarnation-record ordering, heartbeat-storage, and process-model proof added.
- **§6:** accept, with §5’s entry-point inventory, privileged/ordinary separation, gate-action atomicity, and post-release integrity rule completed.
- **§7:** accept. The single abandonment mechanism, serialization winner, dual validation gates, in-memory default for `commit_pending`, and batch-control-only revision consequence are the correct direction. R6 must still integrate abandonment with the operation map and H2 tables.
- **§8:** accept, with §6’s immutable-class separation and reference-aware retention/cleanup linkage completed.
- **§9:** accept. OD-17 remains a hard prerequisite and must include controlling-language analysis.
- **§§10–11 and Appendix A:** accept as preservation, precedence, status, and identity treatment. No disposition is converted into a freeze claim.

---

## 9. Authorization boundary

This review:

- does not approve R5;
- does not rewrite Addendum 2, Addendum 1, or the original matrix;
- does not authorize R6;
- does not adopt V18–V24;
- does not freeze or accept F3-C;
- does not authorize P1a or implementation;
- does not unblock F3-D;
- does not modify any reviewed artifact.

The controlled record should now be treated as:

- original matrix: preserved adjudication authority;
- first precision addendum: preserved precision authority;
- Addendum 2: preserved recovery/durability/enforcement authority, subject to §§3–7 above;
- R5: proposal only;
- R6: unauthorized and unwritten;
- A14: open;
- V24: blocked;
- F3-C phase 0: open;
- implementation: unauthorized.

The next defensible step is not more parallel paperwork. It is proof that the recovery, durability, gate, abandonment, retention, and lock-order rules compose into one executable transaction lifecycle. That proof belongs in R6 only after OD-17, OD-02a, and OD-02b are separately decided.
