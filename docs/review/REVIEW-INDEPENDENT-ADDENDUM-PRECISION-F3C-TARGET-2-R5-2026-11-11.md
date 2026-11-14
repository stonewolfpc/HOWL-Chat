# REVIEW — Independent Review of the Precision Addendum to the R5 Disposition Matrix

**Status:** **INDEPENDENT REVIEW ONLY — NOT APPROVAL, NOT FREEZE, NOT R6, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Reviewed addendum:** `docs/review/ADDENDUM-PRECISION-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md`
**Addendum SHA-256:** `d059054b678f95d8697b8987bf856035a20a6415bd61828a940ebbd99d7549b0`
**Addendum size:** 340 lines; 18,478 bytes.
**Original matrix:** `docs/review/DISPOSITION-MATRIX-F3C-TARGET-2-R5-MUSE-2026-11-11.md`
**Original matrix SHA-256:** `061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6`
**Independent matrix review:** `docs/review/REVIEW-INDEPENDENT-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md`
**Independent matrix review SHA-256:** `e28e7aa0c81f2bfd739bf396f6e28dfbbd57bc2f28a00575c6302c6200c9a9bb`
**R5 proposal:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md`
**R5 SHA-256:** `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e`
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md`
**Gate v2 SHA-256:** `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`
**Scope:** documentation-only independent review. No source, test, build, migration, model-load, destructive-test, application-data, authority-record, freeze, acceptance, P1a, R6, or F3-D action was performed or authorized. No existing file was modified.

---

## 1. Verdict

**Preserve the addendum as the precision authority for the matrix, subject to the R6-proof requirements below. Do not commission R6 yet.**

The addendum correctly tightens the matrix in the three areas emphasized for this review:

- OD-02a cannot create durable recovery authority;
- OD-17 is a hard prerequisite to new storage roots;
- abandonment and candidate-commit ordering has a defined winner rather than an ambiguous middle state.

Its proposed recovery sequence is also directionally sound:

```text
verified package
  -> final lease/authority proof
  -> atomic create-only record installation
  -> durable repairing/incomplete state
  -> per-file replacement
  -> manifest verification
  -> verified/released state
```

But the sequence is not yet crash-proof as written. In particular, the addendum does not fully close the window between record installation and durable `repairing` state, does not separate atomic visibility from power-loss durability with enough precision, and leaves several enforcement and retention linkages for R6 to prove.

Accordingly:

- **OD-17:** preserve as a hard prerequisite.
- **OD-02a/OD-02b boundary:** preserve, with an activation condition added.
- **Commit-to-repairing recovery:** preserve the direction, but require an explicit package/record/current-target reconciliation rule.
- **Atomic visibility versus durability:** require separate concurrent-linearization and power-loss-commit rules.
- **Lease expiry/incarnation:** preserve the direction, but require the missing incarnation, heartbeat, suspend, and process-model records.
- **Startup gate:** preserve the direction, but require privileged-recovery access, atomic gate checks, and explicit inclusion/exclusion treatment.
- **Abandonment ordering:** preserve the winner rule, but make the winner deterministic and cover lease issuance as well as final commit.
- **A14:** preserve the immutable-evidence rule, but require retention linkage and enough durable fields to establish divergence.
- **R6:** remains unauthorized.
- **Implementation:** remains unauthorized.

---

## 2. Artifact-integrity check

The addendum’s Appendix A correctly records:

- original matrix hash `061b...`;
- independent matrix-review hash `e28e...`;
- R5 hash `bbe9...`;
- frozen-contract hash `6104...`;
- Gate v2 hash `7fdc...`.

Those values match the current bytes of the corresponding files. The addendum follows the same self-hash convention as the matrix: it does not embed its own hash. Its current computed hash is recorded above.

This verifies present bytes and stated identities. It does not establish historical authorship or repository provenance.

---

## 3. OD-17 and the exhaustive-layout question

The addendum’s OD-17 treatment is correct and should be preserved as a hard prerequisite.

It properly requires the owner to choose between:

- **Option A:** frozen §4.1 is exhaustive, so every new root needs its own formal layout amendment; or
- **Option B:** frozen §4.1 is non-exhaustive, but every new path still needs explicit OD-02b enumeration and cannot be inferred from silence.

It also correctly states the closed default:

- no ruling means no new root;
- no ruling means OD-02b is not approved;
- no ruling means R6 is not authorized.

That is the right control. Ambiguity about the layout must not become implied storage authority.

For R6, one further requirement is needed. The owner ruling should not be a bare selection of “A” or “B.” It should cite the controlling frozen language and address whether any other frozen section independently authorizes operational, evidentiary, temporary, gallery, imported-content, or recovery storage. Otherwise, R6 may resolve §4.1 while leaving another layout-adjacent provision open to expansive interpretation.

OD-17 numbering itself can remain provisional. The substance—no new persistent root by implication—is what must be preserved.

---

## 4. OD-02a/OD-02b boundary

The addendum correctly withdraws durable projection-status tracking from OD-02a and places all durable recovery coordination under OD-02b. That prevents a narrow replacement permission from becoming accidental authorization for persistent recovery infrastructure.

One activation condition is still missing: **OD-02a must be inert without OD-02b.**

The addendum defines a joint interface, but does not expressly state whether approval of OD-02a alone would permit any target replacement. It should not. Per-file replacement after a commit is unsafe without:

- a durable repairing/incomplete marker owned by OD-02b;
- a shared gate blocking supported readers;
- package-derived repair and manifest verification;
- a durable verified/released transition.

**Concrete boundary failure:**

1. The owner approves OD-02a but defers OD-02b.
2. An implementation replaces target files from a verified package.
3. It has no approved place to record repairing state.
4. A crash leaves mixed targets with no durable gate.
5. Restart cannot distinguish safe completion from interrupted replacement except by re-scanning package/record/target state—functionality that itself belongs to the deferred OD-02b design.

**Required R6-proof rule:**

- OD-02a authorizes replacement mechanics only when invoked through the approved OD-02b recovery gate.
- Approval of OD-02a alone authorizes no target write.
- The same operation may exercise both authorities only through their defined joint transaction; it may not write an OD-02b marker under OD-02a authority.

The lock discipline also needs to be explicit. The joint transaction should state whether package verification, final proof, record installation, marker installation, per-file replacement, manifest verification, and gate release all occur under the same `chatWriteMu` acquisition or under a defined sequence of acquisitions. If multiple acquisitions are used, R6 must prove that every intermediate state is either unobservable or recoverable.

---

## 5. Commit-to-`repairing` crash window

The user’s technical concern is valid: there is a crash window between:

```text
atomic create-only record installation
```

and:

```text
durable projection state says repairing/incomplete
```

The addendum mitigates this by requiring ambiguous recovery state to fail closed and requiring recovery to re-derive repair work from the verified package. That is the correct principle, but R6 still needs an explicit discovery rule.

A missing `repairing` marker cannot mean “no repair needed.” After a commit record exists, startup must compare three things:

1. the latest authoritative commit record;
2. the package and manifest identified by that record;
3. the hashes of the currently projected target files.

If the live files do not match the committed package manifest, the committed operation is authoritative and its targets are incomplete. Recovery must reinstall them from the package. The absent marker is therefore evidence of interruption, not evidence of completion.

**Concrete sequence that R6 must survive:**

1. The package is durable and verified.
2. The final lease/authority proof passes.
3. The commit record is installed.
4. The process crashes before durable `repairing` state exists.
5. Restart finds a valid record and package but no repair marker.
6. Current targets still contain the prior revision.
7. Startup must not publish the old targets as current for the new record.
8. It must not treat the missing marker as proof that projection finished.
9. It must reinstall all package-derived files and verify the manifest before opening the read gate.

The package makes this recoverable, but only if R6 defines the above reconciliation explicitly. “Re-derive repair work from the verified package” is a requirement, not yet a recovery algorithm.

R6 must also cover the converse case:

- durable `repairing` exists;
- no corresponding commit record exists;
- package may be orphaned or incomplete.

That state must resolve to “no commit,” with orphan cleanup governed by explicit retention rules. It must not be repaired as if it were a commit.

Finally, R6 must define retention linkage: a committed package must remain available until its targets are verified/released and any applicable retention rule permits cleanup. A package-cleanup rule that can delete the recovery source before projection completes would defeat the entire package-authority design.

---

## 6. Atomic visibility versus power-loss durability

The addendum correctly asks R6 to replace absolute “no I/O” language with a durability-primitive rule, including:

- exclusive creation;
- file synchronization;
- directory synchronization;
- ordering;
- crash behavior;
- unsupported-filesystem behavior.

That must go one step further and distinguish two different commit properties.

### 6.1 Concurrent visibility

For threads and processes that remain alive, the commit point can be the atomic appearance of the create-only record name. Exclusive creation or an atomic same-directory rename can establish a single linearization point.

### 6.2 Power-loss durability

Atomic visibility alone does not establish that the commit survives power loss. Without the required file and directory synchronization, a crash or power failure may lose:

- the new record;
- the directory entry that makes it visible;
- the package data on which recovery depends;
- the ordering between package, record, marker, and projection state.

R6 must therefore define:

- when a commit is visible to live participants;
- when it is durable against process crash;
- when it is durable against host crash or power loss;
- what recovery may infer in each window.

The most demanding window is visibility without proven durability:

1. Final validation passes.
2. The record becomes visible.
3. Required directory synchronization has not completed.
4. Power is lost.
5. On reboot, the record may be present or absent.
6. If present, recovery must be able to verify the referenced package before treating the record as authoritative.
7. If absent, recovery must treat the package as orphaned and make no commit.

That behavior is achievable, but only if the package is made crash/power-loss durable before the final gate and the record’s durability primitive is fully specified. The addendum requires package flushing/read-back before final validation, but “flush” must be expanded into exact file and directory durability operations.

Platform behavior also matters. R6 must not assume uniform POSIX rename/sync semantics across all supported filesystems and operating environments. It must define the required same-volume behavior, prohibited filesystem features, and failure modes for unsupported behavior.

---

## 7. Lease expiry and process-incarnation recovery

The addendum preserves the correct recovery principle: old monotonic readings cannot be compared with a new process’s clock, and incarnation fencing—not timestamp arithmetic—must defeat a resumed stale holder.

The following records and behaviors are still required.

### 7.1 Current-incarnation authority

R6 must define the durable current-incarnation record, including:

- path/root/scope and owner;
- allocation and uniqueness discipline;
- write order relative to lease startup reconciliation;
- crash behavior before, during, and after update;
- missing/torn/duplicated/inconsistent behavior;
- inventory and wipe treatment.

Ambiguous incarnation state must fail closed. Recovery cannot rely on process memory alone.

The phrase “uniqueness-checked” also needs exact meaning. A 128-bit random value has a negligible collision probability, but that is not the same as a checked uniqueness property. R6 must state whether uniqueness means:

- probabilistic uniqueness from a named CSPRNG;
- comparison against retained incarnation history;
- allocation from a durable monotonic sequence; or
- another explicit mechanism.

If incarnation history is used, its retention, bounds, inventory, and wipe treatment must also be defined.

### 7.2 Clean shutdown versus crash

R6 must not assume that a missing shutdown marker proves a crash or that a present shutdown marker proves a clean shutdown. Either marker state can be interrupted.

The safe default is:

- every restart allocates a new incarnation;
- every old incarnation is unusable;
- a clean-shutdown marker may optimize recovery but cannot confer authority.

If the design distinguishes clean shutdown from crash, the marker protocol needs the same crash-state analysis as leases and barriers.

### 7.3 Heartbeats as durable mutations

A validating heartbeat is itself a durable state change. R6 must define:

- heartbeat record and storage;
- cadence relative to TTL and maximum lifetime;
- missed-heartbeat threshold;
- validation performed before lifetime extension;
- crash behavior during heartbeat update;
- OD-02b ownership and H2 transition-table treatment.

A heartbeat must never revive a revoked, expired, orphaned, terminal, or barrier-drained lease. It also must not be required inside the final no-I/O commit section. Those two rules are compatible only if staging is kept outside the final section or oversized work fails closed before entering it.

### 7.4 Process model

Incarnation fencing assumes that only one authority process can hold the relevant lock and update the current-incarnation record. R6 must state the supported process model and define fail-closed behavior for:

- two live processes;
- split storage mounts;
- failed exclusive locking;
- snapshot duplication;
- container duplication.

This does not require solving arbitrary distributed consensus. It requires stating the deployment assumption and making violations observable and fail-closed rather than silently permitting two current incarnations.

---

## 8. Startup-gate coverage

The addendum correctly moves the gate from a startup convention to an enforceable storage-access boundary and enumerates ordinary, operational, migration, wipe, recovery, gallery, attachment, background, autosave, diagnostic, export, and copy paths.

Two structural refinements remain.

### 8.1 Privileged recovery access must be distinguished from ordinary access

The addendum includes “recovery readers themselves” in the covered set. Read literally, that would prevent recovery from inspecting the unverified state it must repair.

R6 must define two access classes:

- **privileged recovery/repair access**, permitted to inspect package, record, marker, and current-target state for reconciliation; and
- **ordinary supported access**, prohibited from publishing, copying, migrating, wiping from, or committing on unverified state.

Diagnostic or export tools must not be placed in the privileged class merely for convenience.

### 8.2 Gate checks must be atomic with the protected action

A gate is not enforceable if a reader checks “released,” then acts after a new commit begins projection. R6 must define the shared lock or sequence discipline that binds:

- gate-state observation;
- target reads;
- migration/wipe-source selection;
- new commit admission;
- projection and repair;
- gate release.

At minimum, gate release and the start of the next mutation’s projection must be serialized so an ordinary reader cannot observe a released gate and then read files being replaced by a successor commit.

Post-release corruption or out-of-band modification is a separate integrity issue, but R6 should state whether ordinary readers verify target hashes on read, only recovery verifies them, or verification occurs at another defined boundary.

---

## 9. Abandonment versus candidate-commit ordering

The addendum’s winner rule is correct:

- abandonment before the candidate final gate defeats the candidate;
- an installed candidate record defeats later abandonment of that turn;
- no middle state may resurrect an abandoned candidate.

To make that rule deterministic, R6 must add three refinements.

First, the winner must be determined by serialization order under the shared lock, not by an implementation choice at the race point. The addendum’s “final gate must reject the candidate or abandonment must fail” should therefore be rewritten as:

- whichever operation reaches the serialization point first wins;
- the loser fails with an explicit abandonment/serialization error;
- neither operation may overwrite the other silently.

Second, slot-state validation must occur at both:

- candidate/commit lease issuance, for fail-fast behavior; and
- the final pre-record commit gate, for authority.

Final-gate validation alone is necessary but not ideal: it permits an abandoned slot to acquire lease state that must later be rejected. Issuance-time validation reduces, but does not replace, final-gate authority.

Third, `commit_pending` must have an exact durability meaning. If it is only an in-memory phase between final validation and record installation, R6 should say so. If it is durable, R6 must define:

- its writer;
- its transition into and out of the state;
- its crash behavior;
- its interaction with abandonment;
- its inventory and wipe treatment.

A durable `commit_pending` state without recovery semantics could recreate the same ambiguity the winner rule is meant to eliminate.

The abandonment operation itself also needs a settled mechanism. R6 may use a fence increment, slot-scoped revocation, or both, but it must select one interoperable rule rather than preserve alternatives. Likewise, if abandonment writes durable state, that state belongs in the OD-02b inventory and H2 transition tables.

---

## 10. A14 immutable context evidence

The addendum correctly requires either embedded context evidence or a hash-bound immutable pointer. A bare mutable candidate or attempt identifier is properly rejected.

The remaining R6-proof requirement is retention and structural linkage.

If history points to immutable candidate evidence, R6 must ensure that the evidence remains available for as long as the history entry that cites it. Otherwise, a hash-bound pointer may be cryptographically precise but operationally dangling.

R6 must therefore define:

- whether immutable candidate evidence is a separate create-only class from mutable candidate lifecycle state;
- the exact hash/version identity used by the history pointer;
- retention linkage between history entries and cited candidate evidence;
- inventory treatment for both sides of the link;
- wipe behavior ensuring that linked history and evidence are removed consistently;
- behavior when cited evidence is missing, corrupt, or outside the wipe/source boundary.

The durable history entry should also contain enough direct fields to establish divergence without first retrieving candidate evidence. At minimum, R6 should preserve both:

- the history/parent/base ordering context in which the turn was persisted; and
- the context head actually observed during generation.

A pointer can supply the latter, but embedding it is simpler and less dependent on long-term cross-root retention. Either approach remains acceptable if the immutability and retention proofs are complete.

A14 remains open. The baseline, OD-15 decision, and V24 evidence remain missing, and V24 remains blocked.

---

## 11. Final authorization boundary

This review:

- does not approve R5;
- does not rewrite the original matrix;
- does not approve the addendum as a contract revision;
- does not authorize R6;
- does not adopt V18–V24;
- does not freeze or accept F3-C;
- does not authorize P1a or implementation;
- does not unblock F3-D;
- does not modify any reviewed artifact.

The controlled record should now be treated as:

- original matrix: preserved adjudication authority;
- precision addendum: preserved precision authority, subject to §§3–10 above;
- R5: proposal only;
- R6: unauthorized and unwritten;
- A14: open;
- V24: blocked;
- F3-C phase 0: open;
- implementation: unauthorized.

The next step, when authorized, is not another broad rewrite. It is an owner ruling on OD-02a, OD-02b, and OD-17, followed by R6 contract language that proves the crash, clock, gate, abandonment, and evidence-linkage cases above.
