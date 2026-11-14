# ADDENDUM 2 — Recovery, Durability, and Enforcement Corrections

**Status:** **DOCUMENTATION-ONLY ADDENDUM — NOT R6, NOT A CONTRACT REVISION, NOT FREEZE, NOT ACCEPTANCE, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Supplements:** `docs/review/ADDENDUM-PRECISION-DISPOSITION-MATRIX-F3C-TARGET-2-R5-2026-11-11.md`
**First addendum SHA-256:** `d059054b678f95d8697b8987bf856035a20a6415bd61828a940ebbd99d7549b0`
**First addendum size:** 340 lines; 18,478 bytes.
**Original matrix:** `docs/review/DISPOSITION-MATRIX-F3C-TARGET-2-R5-MUSE-2026-11-11.md` (`061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6`)
**Independent second review:** `docs/review/REVIEW-INDEPENDENT-ADDENDUM-PRECISION-F3C-TARGET-2-R5-2026-11-11.md`
**Second review SHA-256:** `46033f073db1739a4f6c9b1e92295ab632e081b727af2c89dce1eba3fa7783b5`
**Second review size:** 421 lines; 21,795 bytes.
**R5 proposal:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R5-2026-11-11.md` (`bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e`)
**Frozen Master:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md` (`61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`)
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md` (`7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`)

**Scope:** this addendum preserves the original matrix, the first precision addendum, R5, frozen Master, Gate v2, and all prior reviews. It accepts the second review's seven findings, adds the missing crash-recovery protocol, makes OD-02a explicitly inert without OD-02b, separates visibility from power-loss durability, defines the storage-access boundary, fixes deterministic abandonment ordering, and requires immutable retention-linked A14 evidence. It does not write or authorize R6.

---

## 1. Disposition of the seven findings

All seven findings are accepted. None reopens C1–C6 or H1–H7. They tighten the adjudication record that will guide R6.

| # | Finding | Disposition | Required effect |
|---|---|---|---|
| **1** | OD-02a must be inert by itself. | **Accept** | OD-02a approval defines replacement mechanics but authorizes no target write. Replacement requires an approved OD-02b gate and the joint transaction. |
| **2** | A missing `repairing` marker is a recovery problem. | **Accept** | Startup must treat record/package/live-target mismatch as interruption, not completion, and reconstruct repair from the authoritative record and package. |
| **3** | Atomic visibility is not power-loss durability. | **Accept** | Separate live linearization from crash/power-loss durability, with exact package, record, directory, synchronization, and ambiguous-interruption rules. |
| **4** | Incarnation recovery needs one authoritative process model. | **Accept** | Define the durable current-incarnation record, uniqueness meaning, clean-shutdown limits, durable heartbeats, and multi-process fail-closed behavior. |
| **5** | Startup gate must separate privileged recovery from ordinary access. | **Accept** | Define two access classes and bind gate checks atomically to the protected read, source selection, commit admission, projection, and release. |
| **6** | Abandonment ordering must be deterministic and cover both validation gates. | **Accept** | Serialization order under the shared lock determines the winner; validate slot state at lease issuance and at final commit; define `commit_pending`. |
| **7** | A14 immutable evidence must remain available while history references it. | **Accept** | Define immutable evidence retention, exact hash/version linkage, inventory/wipe treatment, and direct durable fields sufficient to establish divergence. |

The second review's **OD-17 precision requirement** is also accepted: the layout ruling must cite the controlling frozen language and address any other frozen provision that might independently authorize storage. It is a control requirement, not an eighth finding.

No finding is rejected. The first addendum remains the precision authority except where this second addendum tightens or supersedes it.

---

## 2. Revised OD-02a / OD-02b activation rule

### 2.1 OD-02a is not a write authorization

Approval of OD-02a alone authorizes **no target write**.

OD-02a defines only the mechanics by which an already authorized replacement may occur:

- per-file replacement from a verified, durable package;
- package-manifest and target-hash verification;
- package-derived optional-file deletion;
- post-commit application of target bytes.

It does not authorize:

- the first target write;
- creation of a durable repair/status marker;
- a read gate;
- recovery coordination;
- a new root;
- a mutable operation-evidence record.

### 2.2 Both approvals are required

The permission model is:

| Approved authority | Permitted effect |
|---|---|
| OD-02a only | Replacement mechanics defined; **no target write authority**. |
| OD-02b only | Durable recovery infrastructure defined; **no target replacement authority**. |
| OD-02a and OD-02b | The joint transaction may be authorized, subject to the remaining gates. |

A target replacement is authorized only when invoked through the approved OD-02b recovery gate and only within the joint transaction. An OD-02b marker may not be written under OD-02a authority. An OD-02a replacement may not bypass the OD-02b gate.

### 2.3 Lock discipline

R6 must state whether the joint transaction uses:

- one continuous `chatWriteMu` acquisition for package proof, record installation, marker installation, per-file replacement, manifest verification, and gate release; or
- a defined sequence of acquisitions.

If multiple acquisitions are used, R6 must prove that every intermediate state is either unobservable to ordinary access or recoverable from the record and package. No intermediate state may be treated as verified current state merely because a prior step returned success.

The commit record remains the single commit point. The recovery gate is a post-commit authority boundary, not a second commit point.

---

## 3. Missing-marker crash recovery

### 3.1 Record installed, marker absent

A missing `repairing` marker does not mean “no repair needed.” After a commit record exists, startup must reconcile:

1. the latest authoritative commit record;
2. the package and target manifest identified by that record;
3. the hashes of the currently projected target files;
4. the presence or absence of any durable projection/recovery state.

If the live targets do not match the committed package manifest, the commit is authoritative and the targets are incomplete. Recovery must reinstall all package-derived files and verify the manifest before ordinary access.

The required rule is:

```text
record present
  -> verify referenced package and package hash
  -> verify record/manifest agreement
  -> compare current targets to committed manifest by hash
  -> mismatch, missing, or ambiguous current state
       => recovery_required
       => reinstall all package-derived targets
       => verify manifest
       => release gate
  -> exact match plus durable verified state
       => ordinary access may open
```

A missing marker must be treated as evidence of interruption, not evidence of completion.

### 3.2 Marker present, record absent

If a durable `repairing` state exists without a corresponding authoritative commit record:

- no commit exists;
- the state must not be “repaired” as if a commit existed;
- the package is orphaned or incomplete;
- cleanup follows explicit retention rules;
- ordinary access must not treat the orphan as current.

### 3.3 Target/package retention linkage

A committed package must remain available until:

- its targets are verified and released; and
- any applicable retention rule permits cleanup.

No package-cleanup rule may delete the recovery source before projection completes. A package retained only by a manifest digest is not sufficient.

### 3.4 Crash sequence that must be covered

1. Package is durable and verified.
2. Final lease/authority proof passes.
3. Commit record is installed.
4. Process crashes before durable `repairing` state exists.
5. Restart finds a valid record and package but no repair marker.
6. Current targets still contain the prior revision.
7. Startup must not publish the old targets as current for the new record.
8. Startup must not treat the missing marker as proof that projection finished.
9. Startup must reinstall package-derived files and verify the manifest before releasing ordinary access.

---

## 4. Atomic visibility versus power-loss durability

### 4.1 Two distinct properties

R6 must separately define:

- **Concurrent visibility / linearization:** the moment live participants can observe the create-only record name.
- **Power-loss durability:** the moment the record, its directory entry, and the package needed for recovery are guaranteed to survive host crash or power loss.

Atomic visibility alone is not durability. A record may be visible to live threads and still be lost, or become ambiguous, after power loss.

### 4.2 Package durability

Before the final authority gate, R6 must require the package to be durable:

- package files written and flushed;
- package manifest and hashes verified;
- file contents and required directory entries synchronized;
- read-back verification completed;
- package path installed create-only with collision behavior defined.

A failure in any required package durability step aborts the operation before the commit record is installed.

### 4.3 Record visibility and durability

R6 must define:

- the exact exclusive-create or atomic create-only primitive;
- why it cannot overwrite an existing record;
- how the record becomes visible;
- when its file contents are synchronized;
- when its parent directory entry is synchronized;
- what happens if visibility succeeds but synchronization fails;
- what recovery may infer in each window.

The commit linearization point may be the record’s atomic appearance, but the operation may not be reported as power-loss durable until the required file and directory synchronization completes.

### 4.4 Ambiguous interruption

If power loss leaves the record present:

- recovery must verify the package before treating the record as authoritative;
- a record without a valid package is a hard evidence block, not a guessed commit.

If power loss leaves the record absent:

- the package is orphaned evidence;
- no commit exists;
- cleanup follows retention rules;
- no target projection may be adopted from the orphan alone.

### 4.5 Platform rules

R6 must define the supported same-volume behavior, required filesystem guarantees, prohibited filesystem features, and failure modes for unsupported behavior. It must not assume uniform POSIX rename/sync semantics across platforms.

It must also state which synchronization failures invalidate the operation, which leave a visible-but-not-durable record for recovery, and which require a hard block.

---

## 5. Incarnation and process model

### 5.1 Current-incarnation authority

R6 must define a durable current-incarnation record with:

- exact path, root/chat scope, and owner;
- allocation and write order relative to lease startup reconciliation;
- crash behavior before, during, and after update;
- missing, torn, duplicated, stale, or inconsistent behavior;
- inventory and wipe treatment;
- how a clean shutdown is distinguished from a crash, if at all.

Ambiguous incarnation state fails closed. Recovery may not rely on process memory.

### 5.2 Uniqueness meaning

“Uniqueness-checked” must have an exact meaning. R6 must select one:

- probabilistic uniqueness from a named CSPRNG;
- comparison against retained incarnation history;
- allocation from a durable monotonic sequence;
- another explicit mechanism.

If incarnation history is retained, its bounds, inventory, and wipe treatment must be defined.

### 5.3 Clean shutdown is not authority

A missing shutdown marker does not prove a crash; a present marker does not prove a clean shutdown. The safe default is:

- every restart allocates a new incarnation;
- every old incarnation is unusable for authority;
- a clean-shutdown marker may optimize recovery but cannot confer authority;
- any shutdown marker protocol has the same crash analysis as leases and barriers.

### 5.4 Durable heartbeats

A validating heartbeat is a durable state change. R6 must define:

- the heartbeat record and storage;
- cadence relative to TTL and maximum lifetime;
- missed-heartbeat threshold;
- validation before lifetime extension;
- crash behavior during heartbeat update;
- OD-02b ownership and H2 transition-table treatment.

A heartbeat must never revive a revoked, expired, orphaned, terminal, or barrier-drained lease. It must not be required inside the final no-I/O commit section.

### 5.5 Supported process model

R6 must state the deployment assumption and fail closed for:

- two live authority processes;
- split storage mounts;
- failed exclusive locking;
- snapshot duplication;
- container duplication.

The design does not need arbitrary distributed consensus. It must make multiple-current-incarnation states observable and closed, not silently accepted.

---

## 6. Storage-access boundary

### 6.1 Two access classes

R6 must define:

- **privileged recovery/repair access:** permitted to inspect package, record, marker, and current-target state for reconciliation, under the recovery locks and authority.
- **ordinary supported access:** prohibited from publishing, copying, migrating, wiping from, or committing on unverified state.

Diagnostic, export, gallery, attachment, imported-content, or background tools do not become privileged merely because they need access. Each must have a defined class and a defined gate.

### 6.2 Covered paths

At minimum, R6 must enumerate:

- ordinary chat history and target readers;
- metadata, preference, cast, continuation, and projection readers;
- plan, batch, attempt, and candidate readers;
- migration source selection;
- wipe source selection and inventory readers;
- recovery readers themselves;
- gallery, attachment, world, scenario, and imported-content paths where applicable;
- autosave, background, diagnostic, export, and copy paths capable of reading or publishing chat state.

An out-of-contract direct filesystem read cannot be physically prevented, but every supported entry point must pass through the same gate. No subsystem may bypass the gate because it is “separate” or “only diagnostic.”

### 6.3 Atomic gate checks

A gate check that is separated from the protected action is not enforceable. R6 must bind, under a shared lock or defined sequence:

- gate-state observation;
- target reads;
- migration/wipe-source selection;
- new commit admission;
- projection and repair;
- gate release.

At minimum, gate release and the next mutation’s projection start must be serialized so an ordinary reader cannot observe a released gate and then read files being replaced by a successor commit.

R6 must state whether ordinary readers verify target hashes on read, only recovery verifies them, or verification occurs at another defined boundary. Post-release out-of-band modification is a separate integrity issue but must not be silently treated as verified.

### 6.4 Fail-closed state

Missing, corrupt, torn, stale, duplicated, or inconsistent recovery/projection state means **unverified**:

- unverified targets are not current;
- no supported reader may publish or copy them;
- migration and wipe may not select them as sources;
- no new commit may proceed;
- recovery re-derives repair from the verified package;
- manifest verification is required before gate release;
- a crash after repair but before release requires re-verification.

---

## 7. Deterministic abandonment ordering

### 7.1 One durable operation

`AbandonBatchSlot` is a durable operation, not a flag. R6 must define:

- operation identity, plan hash, batch ID, participant index, expected state, reason, actor;
- exact operation record and path;
- fence/revocation mechanism;
- revision/history consequence;
- recovery, inventory, and wipe treatment;
- fail-closed behavior for ambiguous plan, batch, slot, or current state.

### 7.2 One mechanism

R6 must select one interoperable abandonment mechanism. Recommended single rule:

- increment `fenceEpoch`; and
- revoke affected slot-scoped leases.

No implementation may choose fence increment or slot revocation independently. If both are used, they are one defined rule with one ordering and one durable record.

### 7.3 Winner rule

The winner is determined by serialization order under the shared `chatWriteMu`, not by an implementation choice at the race point:

1. Whichever operation reaches the serialization point first wins.
2. The loser fails with an explicit abandonment/serialization error.
3. Neither operation may silently overwrite the other.
4. If abandonment wins before the candidate’s final gate, no record is installed.
5. If the candidate record is installed first, the turn is committed; abandonment cannot rewrite or invalidate it.
6. No abandoned slot may commit through retry, lease replay, recovery, or old holder state.

### 7.4 Two validation gates

Slot state must be validated at:

- candidate/commit lease issuance, for fail-fast behavior; and
- the final pre-record commit gate, for authority.

Final-gate validation alone is necessary but not sufficient. Issuance-time validation does not replace final-gate authority.

### 7.5 `commit_pending`

R6 must give `commit_pending` one exact meaning. Recommended: it is an in-memory phase under `chatWriteMu` between final validation and record installation. It is not durable, cannot be cited by recovery, and cannot authorize a commit after restart.

If `commit_pending` is persisted instead, R6 must define its writer, transitions, crash behavior, abandonment interaction, inventory, and wipe treatment. A durable `commit_pending` without recovery semantics is not acceptable.

### 7.6 Revision/history consequence

Recommended rule: abandonment writes a durable batch-control operation record and does **not** advance ordinary chat history or revision unless the owner separately authorizes a history mutation.

The operation record must be as durable and ordered as a commit record. It belongs in the OD-02b inventory and the H2 transition tables.

---

## 8. A14 immutable evidence retention

### 8.1 Evidence class

If A14 history points to candidate evidence, that evidence must be a separate immutable, create-only class from mutable candidate lifecycle state.

A bare `candidateId` or `attemptId` is insufficient. The pointer must bind to:

- candidate version plus immutable hash;
- package hash;
- `cand3_` fingerprint;
- another explicitly immutable identity.

### 8.2 Retention linkage

R6 must define:

- exact hash/version identity used by the history pointer;
- retention linkage between the history entry and cited evidence;
- inventory treatment for both sides;
- wipe behavior that removes linked history and evidence consistently;
- behavior when cited evidence is missing, corrupt, or outside the wipe/source boundary.

A hash-bound pointer that can dangle is not sufficient.

### 8.3 Direct fields

The durable history entry should preserve both:

- the history/parent/base ordering context in which the turn was persisted; and
- the context head actually observed during generation.

Embedding `contextHistoryHead` directly is preferred. A pointer is acceptable only with the immutability and retention proofs above.

### 8.4 A14 remains open

A14 remains open. The pre-amendment multi-character baseline, OD-15 decision, and V24 evidence remain missing. **V24 remains blocked.**

---

## 9. OD-17 precision

The second review's OD-17 refinement is accepted:

- the owner ruling must cite the controlling frozen language;
- it must address whether any other frozen section independently authorizes operational, evidentiary, temporary, gallery, imported-content, or recovery storage;
- a bare “§4.1 is non-exhaustive” selection is insufficient without that analysis;
- no new persistent root may be created until the ruling is recorded and OD-02b enumerates the path.

The ruling remains tracked as **OD-17 — frozen layout exhaustiveness**; final decision numbering remains the owner’s.

---

## 10. Precedence and preservation

- The original matrix remains preserved at `061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6`.
- The first addendum remains preserved at `d059054b678f95d8697b8987bf856035a20a6415bd61828a940ebbd99d7549b0`.
- This second addendum does not rewrite either document.
- Where this addendum and the first addendum differ, this addendum controls the precision point.
- The matrix remains the adjudication authority; this addendum is its second precision layer.
- Illustrative status names remain unapproved vocabulary.
- No disposition is converted into a freeze claim.

---

## 11. Status and authorization

**Addendum 2:** **DOCUMENTATION ONLY; NOT R6; NOT A CONTRACT REVISION.**
**R5:** **PROPOSAL ONLY; NOT FROZEN; NOT ACCEPTED; NOT IMPLEMENTED; NOT VERIFIED.**
**Original matrix:** preserved.
**First addendum:** preserved and qualified.
**Second review:** accepted as findings; not approval.
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

1. Muse reviews this second addendum.
2. Owner rules on OD-17 using the controlling frozen language.
3. Owner decides OD-02a and OD-02b separately and confirms the joint activation rule.
4. Owner accepts or replaces the crash-recovery, visibility/durability, incarnation, and gate rules.
5. Owner confirms deterministic abandonment and immutable A14 retention.
6. A14 remains open and V24 remains blocked.

No source, test, build, migration, model load, destructive test, application data, authority record, freeze, acceptance, or implementation action was performed. No R6 was drafted.

---

## Appendix A — Identity

| Artifact | SHA-256 |
|---|---|
| Original matrix | `061bbfd2daeeab2ccb398b96a41a76ba39390d77b87f6868fea573e6b19678f6` |
| First precision addendum | `d059054b678f95d8697b8987bf856035a20a6415bd61828a940ebbd99d7549b0` |
| Second precision review | `46033f073db1739a4f6c9b1e92295ab632e081b727af2c89dce1eba3fa7783b5` |
| R5 proposal | `bbe9263a3251168e4ba22d8c327e665b7dbaf50cf4692f2ddf93e62de206275e` |
| Frozen Master | `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31` |
| Gate v2 | `7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a` |

This addendum’s own SHA-256 is recorded in the completion report rather than embedded here, to avoid a self-referential file hash.
