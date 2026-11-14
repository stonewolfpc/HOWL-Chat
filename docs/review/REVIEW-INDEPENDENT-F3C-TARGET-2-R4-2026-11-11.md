# REVIEW — Independent Adversarial Review of F3C-TARGET-2 R4

**Status:** **INDEPENDENT ADVERSARIAL REVIEW — NOT APPROVAL, NOT FREEZE, NOT IMPLEMENTATION.**  
**Date:** 2026-11-11  
**Reviewed artifact:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R4-2026-11-11.md`  
**R4 SHA-256:** `7a2c221147621706e299924c20de1fe29080b7d649fac25264ce212d11e7417a`  
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`  
**R3 SHA-256:** `5013eb6fad7e0ec6f067d0fdbd350fd7f7c68ff9e263a7967b37f39e6daaf42c`  
**R1 review:** `docs/review/REVIEW-ADVERSARIAL-F3C-TARGET-2-2026-11-11.md` (`61ae9d934fd10e526ecb10fabca994274e989f02139dd5d7cac81cce05de0209`)  
**R2 review:** `docs/review/REVIEW-ADVERSARIAL-F3C-TARGET-2-R2-2026-11-11.md` (`70e987297e53f6ba87a84190eef2366eef4e9cf42449ad6b5680551db93bfa73`)  
**R3 review:** `docs/review/REVIEW-INDEPENDENT-F3C-TARGET-2-R3-2026-11-11.md` (`c9e6d25029bc87a44703c16a6c6c4f917af1c0107cf2ddf01c5b083522756f62`)  
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md`  
**Scope:** read-only review. No source, test, build, migration, model load, destructive test, authority change, implementation, acceptance, or freeze.

---

## 1. Overall verdict

**No. R4 is not sufficiently precise or internally consistent to support a formal architecture freeze.**

R4 is materially stronger than R3. It adds a durable mutation authority, a barrier that drains leases and commit gates, a response plan, a separate model-artifact identity, a canonical request envelope, predecessor evidence, and explicit Gate criteria. Those are real improvements, not wording changes.

The proposal still fails the freeze question for six reasons:

1. **The commit record does not retain or identify the target bytes it claims to have committed.** R4 makes the record the linearization point, then installs targets afterward. A crash in that window leaves a valid record with no reconstructable target set. The record has a target-manifest digest, not the target bytes or a durable target-set package.
2. **The lease is not a complete lifecycle.** There is no lease expiry, renewal, heartbeat, holder token, or bounded acquisition policy. `stateSessionID` and `stateEpoch` are absent from the lease and commit record, so a stale autosave or state-only write can survive a switch.
3. **Character removal is not actually bound into candidate identity.** `characterVersion` is present, but `castEpoch` is not in `F3C-CAND-2`, the commit record, or the exact-adoption proof. A removal that advances only `castEpoch` does not necessarily invalidate an in-flight candidate.
4. **The response plan and candidate/attempt records have no durable location or complete state machine.** `abandoned` is described as a terminal state but is not in the candidate status set. A missing lower slot can block every later slot indefinitely, and a partially committed batch has no defined visible state.
5. **`F3C-CANON-2` and `F3C-CAND-2` are not deterministic as written.** The key ordering conflicts with RFC 8785, and `candidateCreatedAt` is inside a fingerprint that is required to be identical for the same candidate.
6. **R4 adds a second authority/lock order without reconciling it with the frozen Master order.** The frozen contract requires `stateApplyMu` and `chatWriteMu`; R4 instead declares `root authority -> chat authority -> commit gate`. R4 does not say how both orders coexist. A dual-order implementation can deadlock or bypass one source of permission.

The result is a better proposal, but not a freeze-ready contract. The most dangerous remaining failure is not “the barrier is missing”; it is “the barrier and commit record exist, but the durable target state is not recoverable from the evidence the proposal calls authoritative.”

**Evidence note:** R4 is a proposal and therefore level 5 (unverified). The frozen contract is level 1. Current source references are level 2 and were read only. There are no level 3 test results, level 4 native/runtime results, or destructive-test results in this review.

---

## 2. Hashes, integrity, and review basis

| Artifact | SHA-256 | Lines | Result |
|---|---|---:|---|
| R4 proposal | `7a2c221147621706e299924c20de1fe29080b7d649fac25264ce212d11e7417a` | 1496 | Matches the supplied R4 hash. |
| Frozen Master | `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31` | 1091 | Matches the supplied authority hash. |
| R3 proposal | `5013eb6fad7e0ec6f067d0fdbd350fd7f7c68ff9e263a7967b37f39e6daaf42c` | 520 | Matches the supplied R3 hash. |
| R1 review | `61ae9d934fd10e526ecb10fabca994274e989f02139dd5d7cac81cce05de0209` | 543 | Prior finding inventory read. |
| R2 review | `70e987297e53f6ba87a84190eef2366eef4e9cf42449ad6b5680551db93bfa73` | 616 | Prior finding inventory read. |
| R3 review | `c9e6d25029bc87a44703c16a6c6c4f917af1c0107cf2ddf01c5b083522756f62` | 461 | Prior finding inventory read. |

R4 has 60 fenced markers, balanced. All 35 prior finding IDs are present in its resolution matrix, and all 24 proposed Gate rows are present. Presence is not resolution: this review checks the actual rules and finds multiple rows that cannot yet be verified.

---

## 3. Severity-ranked findings

### A1 — Critical: commit record can exist without the target bytes needed to recover

**Exact R4 rule:** §2.4 steps 4–8 acquire the lease, stage, verify, acquire the authority and commit gate, install the durable commit record, install the receipt and publish the target, then mark the lease committed. §2.6 lists the commit record fields as `recordVersion, operationId, commitSequence, ... targetManifestDigest, candidateFingerprint, requestFingerprint, receiptId, ...`. §2.4 then states: “The commit record is authoritative for both target bytes and receipt identity.”

**Defect:** The record contains a **digest** of the target manifest, not the manifest bytes and not the target bytes. The staged target set is lease-owned temporary data. The proposal does not require that staged set to remain after the record is installed, and §2.8 step 5 only says a committed record is “included in the destructive outcome”; it does not say the staged target payload is retained until reconstruction succeeds.

**Counterexample:** A save commits a new history, meta, and context target set. The process crashes immediately after the durable commit record is installed and before the target installation/publish step. On restart:

- the commit record is valid;
- `targetManifestDigest` identifies a target set;
- the current target files are still old or partial;
- the staged target set may have been cleaned, because the lease is now terminal;
- no field in the record contains the target bytes or a durable pointer to them.

There is no specified outcome. Recovery can:

1. serve the old target while `meta.revision`/record say the new commit is current — a current-state mismatch;
2. treat the record as absent — losing a committed turn;
3. reconstruct from the digest — impossible without the bytes; or
4. adopt an unverified target — violating A1/A7.

§7.5 has no row for “valid record, missing target.” V2 tests “record without receipt” but not “record without target.” V3 lists “crash after record/target/receipt” without defining the missing-target result. The proposal therefore fails the explicit requirement that a valid record with incomplete projections have a deterministic recovery outcome.

**Required correction:** Choose one enforceable protocol:

- retain the exact staged target-set package, with a manifest and all target bytes, until the commit record, receipt, and target installation are verified; or
- write target bytes into the record (or a durable content-addressed store) before the record is authoritative; or
- define a recovery store that can reconstruct the exact target set from record evidence and prove it by hash.

Whichever is chosen, add a recovery table row for every combination of record, target, receipt, and staging. A record without reconstructable target evidence must be a hard block, not an implicit adoption.

**Gate rows:** V2, V3, V4, V5, V13, V22, V23.

---

### A2 — Critical: R4’s authority/lock order conflicts with the frozen Master locks

**Exact R4 rule:** §2.7 says the only permitted order is:

```text
root authority lock
  -> chat authority lock
    -> chat commit gate
      -> lease-owned staging I/O
```

**Frozen rule:** Master §7.1 defines `castApplyMu`, `chatSwitchMu`, `chatTreeMu`, `stateApplyMu`, `chatWriteMu`, `rootJournalMu`, and `chatStateMu`. Master §7.2 fixes:

```text
castApplyMu
  -> chatSwitchMu
    -> chatTreeMu.RLock/Lock
      -> stateApplyMu(chat)
      -> chatWriteMu(chat)
        -> rootJournalMu
```

`chatStateMu` is a leaf lock. Master §7.3 requires `stateApplyMu` and `chatWriteMu` for the standard durable mutation; §7.4 requires autosave to revalidate `stateSessionID` and `stateEpoch`; §7.5 requires the switch path to hold the outgoing/incoming state and write locks.

**Defect:** R4 does not map `MutationAuthority`, `chat authority`, or the commit gate onto the frozen lock set. It declares a different “only permitted” order and never says whether `stateApplyMu` and `chatWriteMu` are acquired before, inside, or after the commit gate.

**Deadlock counterexample:** Preserve Master §7.2 and add R4:

- a normal save acquires `stateApplyMu`, then waits for the commit gate to obtain a `MutationLease`;
- a character commit acquires the chat authority and commit gate, then waits for `stateApplyMu` to publish the new state.

The two paths hold each other’s required lock. One implementation can deadlock; another can drop `stateApplyMu` and violate the frozen serialization rule. Both can claim compliance with one document or the other.

**Bypass counterexample:** A state-only autosave acquires a lease but not `stateApplyMu`/`chatWriteMu`; it commits a target while an ordinary save is publishing a different candidate. R4’s `baseHistoryDigest` does not fence state-only fields such as world, scenario, cast, or lorebook state. The lease is not the frozen state-session fence.

R4 says the frozen authority is unchanged. It cannot simultaneously add a new “only permitted” lock order and leave the old lock order normative without an explicit bridge. This is a direct contract conflict, not an implementation detail.

**Required correction:** Either:

- amend the frozen lock order explicitly and give one combined order, including `castApplyMu`, `chatSwitchMu`, `chatTreeMu`, `stateApplyMu`, `chatWriteMu`, `rootJournalMu`, `chatStateMu`, the authority lock, and the commit gate; or
- define the authority as an additional lease inside the existing Master order, with no lock acquired in the opposite direction.

The bridge must identify where leases are acquired and released and how publication remains after durable commit.

**Gate rows:** V5, V6, V8, V12, V15, V22.

---

### A3 — Critical: the lease has no expiry, renewal, holder token, or state-session fence

**Exact R4 rule:** §2.3 defines the lease fields as `leaseId, rootId, rootIncarnation, chatId, scope, operationId, operationKind, fenceEpoch, issuedAt, ownerProcessIncarnation, baseRevision, baseHistoryDigest, status`. Issuance rules say the authority records the lease before the first durable write; a prior-incarnation lease is `orphaned`, never adopted.

**Defects:**

- No `expiresAt`, `renewedAt`, heartbeat, maximum lifetime, or authority-side timeout exists.
- No holder credential or token is defined. `leaseId` is not an authorization token; any code in the process can present it unless the authority also authenticates the holder.
- No rule says how a live but abandoned goroutine is distinguished from a still-working owner.
- No rule says how an active lease is swept on ordinary restart, outside a barrier. §2.9 discusses orphan marking on barrier restart, but not a general startup sweep.
- The lease has `baseRevision` but no `stateSessionId` or `stateEpoch`. Frozen Master §7.4 requires autosave to revalidate both. A state-only world/cast/lorebook write can be accepted by the new authority while its state session is already invalid.

**Counterexample:** An autosave obtains a lease in chat A. The user switches to chat B. The autosave goroutine blocks. Its lease remains `active`; the switch does not invalidate it because the lease has no state-session field. When it resumes, its base revision may still be valid for a state-only target, and it commits the discarded outgoing state. Even if the commit fails a history-base check, the contract has no state-session check that proves it must fail.

**Second counterexample:** A process holds a lease for an autosave and hangs. No barrier is open. A later delete waits for the lease, times out, and enters `blocked_timeout`; the chat cannot be deleted until an operator acts. The proposal has no lease expiry or safe abandonment, so an ordinary failure becomes a permanent destructive-operation block. The timeout is safe for data but not a complete lifecycle.

**Third counterexample:** A process restarts after a normal save crash. The authority sees a durable `active` lease whose owner incarnation is old. R4 says the lease is “orphaned, never adopted silently,” but does not define the general startup reconciliation that changes it to `orphaned`; the marking rule is stated in the barrier path. An implementation can leave it `active` forever or adopt it.

**Required correction:** Define:

- a bounded lease lifetime and renewal/heartbeat rule;
- holder identity/token validation;
- authority-side expiry and safe terminalization;
- general startup reconciliation for all active leases, not only barrier-scoped ones;
- `stateSessionId` and `stateEpoch` in the lease and commit record;
- a stale-session outcome that does not depend on history-base equality.

**Gate rows:** V5, V6, V8, V13, V15, V22.

---

### A4 — Critical: `castEpoch` is not bound into candidate or commit identity

**Exact R4 rule:** §3.1 defines `ExecutionScope = (chatId, characterId, characterVersion, stateSessionId)`. §3.4 says removal “terminalizes the slot” and the candidate cannot commit. §2.2 puts `castEpoch` in the chat authority; §10.1 puts `castEpoch` in `meta.json`.

**Exact `F3C-CAND-2` envelope:** §5.6 includes `characterId` and `characterVersion`, but **not** `castEpoch` and **not** `stateSessionId`. The §2.6 commit record also omits both. The §7.4 exact adoption list includes character ID/version and state session in the attempt, but the candidate/record fields needed to prove those values are absent.

**Counterexample:** Character A has `characterVersion = 4` and is generating. The user removes A. The removal updates `castEpoch` from 8 to 9 and terminalizes the slot, but the implementation does not increment `characterVersion` because the version identifies the character artifact, not cast membership. A’s candidate still has `characterVersion = 4`. Its exact parent and history base are valid. The commit coordinator sees a candidate with the same character version and accepts it, reintroducing A’s turn after removal.

The proposal says removal increments version/epoch, but it does not require both, and the candidate does not carry the cast epoch that would make the check enforceable. “A recreated character gets a new ID or version” is also not a rule for the removal path.

**Required correction:** Add `castEpoch` (and `stateSessionId`) to the candidate envelope, commit record, exact-adoption proof, and lease/commit check. Define that removal increments both the slot generation and `characterVersion` where required, or make `castEpoch` the sole binding and carry it everywhere.

**Gate rows:** V8, V15, V22, V24.

---

### A5 — Critical: response plan, attempt, and candidate records have no durable location or complete terminal-state model

**Exact R4 rule:** §3.2 defines a `ResponsePlan` with `planVersion, batchId, userCommitOperationId, historyBaseDigest, participantOrder`. It says the plan is committed by the user-turn mutation. §7.2 says the user-turn record contains `responsePlan`. §10.1 lists additions to `meta.json` but does **not** list the response plan. §7.3 defines an attempt record and statuses, but no path or schema location. §5.6 defines candidate fields, but no candidate store.

**Defects:**

- There is no exact durable path or target for the response plan. “Committed by the user turn” does not say whether it is inside `history`, `meta`, a new record, or a commit-record payload.
- There is no exact durable path or record schema for attempts or candidates.
- Recovery precedence names “response plan and generation attempts” and “candidate staging” but cannot locate them by a defined path.
- Candidate status is `staged | eligible | committed | rejected | conflict`. §3.4 says a lower slot may be `terminal_no_turn` (`failed`, `aborted`, `removed`, `skipped`). `abandoned` is used for a batch deadline but is not in either set. `terminal_no_turn` is not a candidate status.
- No rule defines a lower slot that never starts, never produces a candidate, or is deliberately not invoked. Prefix commit rule 2 requires every lower slot to be `committed` or `terminal_no_turn`; if no such terminal marker is written, later slots wait forever.
- No rule defines whether a partially committed batch is visible as a complete turn, a partial batch, or a transaction. The proposal says prefix commits but not the read/UI visibility contract.

**Counterexample:** Three characters are planned. Slot 2 is never started because the UI cancels it before invocation. Slots 1 and 3 complete. Slot 3 cannot commit because slot 2 is neither committed nor terminal. A deadline is said to mark slot 2 `abandoned`, but `abandoned` is not a valid candidate status and is not in the durable plan. The batch remains indefinitely pending. A restart does not know whether slot 2 was skipped, failed, or is still running.

**Second counterexample:** Slot 1 commits; slot 2 is in flight; the process crashes. On restart, the user sees slot 1’s assistant turn, but no rule says whether the batch is visible, whether slot 2 is retried, or whether the user turn is considered partially answered.

**Required correction:** Define exact durable paths and schemas for the response plan, attempt records, and candidates; add all terminal statuses including `skipped`, `abandoned`, and `terminal_no_turn`; define the deadline operation, the partial-batch visibility state, and recovery for a missing lower slot.

**Gate rows:** V5, V13, V15, V22, V24.

---

### A6 — Critical: `F3C-CANON-2` is not one canonical form

**Exact R4 rule:** §5.4 says “`F3C-CANON-2` uses RFC 8785 JSON Canonicalization Scheme” and then says “Object keys are sorted lexicographically by Unicode code point.”

**Defect:** RFC 8785 sorts object property names by **UTF-16 code units**, not Unicode code points. The distinction is observable for non-BMP keys, escaped surrogate pairs, and strings whose code-point order differs from UTF-16 order. Two serializers can both claim RFC 8785 and R4 and produce different bytes. That changes `modelId`, `requestFingerprint`, `candidateFingerprint`, and every hash that depends on them.

The section also leaves these ambiguities:

- `origin` is described as part of the request identity, but the envelope has only `parameters` values. If `origin` is hashed, the rule that explicit and inherited defaults with the same value produce the same fingerprint is false. If it is not hashed, the stated provenance field is missing.
- `engineExtensionParameters` may duplicate keys already in `parameters`. No precedence or rejection rule exists.
- The “same candidate” rule in §5.6 is contradicted by including `candidateCreatedAt` in the candidate envelope. Time is not a stable property of candidate content. Two identical candidates generated at different times produce different fingerprints, violating the requirement that the same candidate produce the same fingerprint on every implementation.
- `timestamp` formats and timezone normalization are not defined for the candidate/record fields.
- There is no explicit domain separation for `F3C-MODEL-1`, `F3C-REQ-1`, `F3C-CAND-2`, and `F3C-ARCH-1` hashes; the schema version is inside the envelope, but the hash input lacks a stated domain tag.

**Counterexample:** Two implementations canonicalize `{"\uE000":1,"\u{10000}":2}` differently under code-point versus UTF-16 ordering. They produce different `requestFingerprint` values for the same request. The reviewer cannot tell which one is correct from R4.

**Second counterexample:** A candidate is staged twice after a retry. The assistant content, target bytes, model, and request are identical, but `candidateCreatedAt` differs. The fingerprints differ, so duplicate suppression and “same candidate” recovery fail.

**Required correction:** Choose RFC 8785’s UTF-16 ordering or explicitly define a different canonicalization and stop calling it RFC 8785. Remove time from candidate identity or define candidate identity as time-dependent and revise duplicate suppression. Define `origin`, extension precedence, timestamp encoding, and hash domain tags.

**Gate rows:** V1, V13, V18, V19, V22.

---

### A7 — High: `requestFingerprint` collides with the frozen `F3C-REQUEST-1` field

**Exact R4 rule:** §5.1 defines `requestFingerprint = "req_" + SHA-256(F3C-CANON-2(request_envelope))`, where the envelope contains `modelId, executionFingerprint, defaultSetId, promptDigest, templateFingerprint, adapterConfigDigest, parameters, seedMode, effectiveSeed, engineExtensionParameters`.

**Frozen rule:** Master §5.2 defines `requestFingerprint = SHA-256("F3C-REQUEST-1" || 0x00 || canonical request bytes)`, binding operation type, root/chat identity, operation ID, base revision, operation-specific input, Save As source binding, and cast token binding. The frozen commit record and receipt store that value.

**Defect:** R4 uses the same field name for a different semantic object. A v2 commit record cannot store both the operation/request identity required by the frozen receipt chain and the generation request identity required by provenance without one overwriting the other. The two hash domains also differ (`F3C-REQUEST-1` versus a bare `F3C-CANON-2` envelope, with a `req_` prefix).

This is not a harmless naming choice. Recovery, idempotency, and Gate V2/V3 compare receipts by `requestFingerprint`. A generation fingerprint that omits operation identity cannot prove the exact operation; an operation fingerprint that omits model/parameters cannot prove the generation. R4 needs two fields with distinct names and domains.

**Required correction:** Rename the generation value (for example `generationRequestFingerprint`) and preserve the frozen operation `requestFingerprint`, or explicitly amend the frozen record/receipt schema and explain how both identities are retained.

**Gate rows:** V1, V2, V3, V18, V22.

---

### A8 — High: saved preference is assigned to two incompatible authorities

**Exact R4 rule:** §8.1 says the saved continuation preference is the chat’s `continuationModel` in `meta.json`. §8.2 says the root `PreferenceAuthority` serializes explicit writes. §2.5 says: “Model preference | preference | `SetContinuationPreference` writes through the root preference authority; it is not a chat target write.” §10.1 adds `continuationModel` to `meta.json`.

**Defect:** If the preference is stored in `meta.json`, it is a chat target write and must go through the chat commit protocol, including the chat target manifest, revision, receipt, and the frozen `stateApplyMu`/`chatWriteMu` order. If it is stored only in the root `preferences.json`, the chat’s `meta.json.continuationModel` is either stale or must be updated by a second write. R4 does not define the relationship, ordering, or recovery between the two locations.

**Counterexample:** `SetContinuationPreference(C)` updates `preferences.json` seq 9. A concurrent autosave commits a `meta.json` snapshot taken before the preference operation, containing `continuationModel = B`. The autosave has a valid chat commit record and overwrites the durable chat metadata to B, while the root preference says C. On next load, one reader uses meta and another uses the root authority. The user’s explicit choice is lost or becomes implementation-dependent.

**Second defect:** `preferenceSeq` is described as strictly increasing and “last accepted sequence wins,” but there is no `expectedPreferenceSeq` compare-and-set or operation ordering proof. The authority lock may serialize writes, but an older UI operation can acquire the lock after a newer one and receive the higher sequence. “Last accepted” is not necessarily “last user action.”

**Required correction:** Choose one durable source of truth. If `meta.json` is authoritative, make preference an ordinary chat commit with the frozen locks and a `preferenceSeq` CAS. If the root file is authoritative, remove `meta.continuationModel` or define it as a verified projection with a deterministic repair rule. Define the expected-sequence conflict result.

**Gate rows:** V5, V15, V21, V22.

---

### A9 — High: migration order contradicts the ordinary commit protocol, and the archive schema is incomplete

**Exact R4 rule:** §2.4 installs the commit record at step 6, then installs the receipt and publishes the target at step 7. §6.4 instead installs staged v2 targets at step 6, installs the migration commit record at step 7, and installs the receipt at step 8. §6.5 says “All v2 targets installed, before record → restore v1 from verified archive.”

**Defects:**

- The two protocols define different linearization orders. A reader cannot know whether a target installed before a record is uncommitted staging, a recoverable migration, or an invalid state.
- “Install staged v2 targets atomically” is asserted for a multi-file target set with no atomic multi-file primitive, no per-file commit order, no fsync/rename protocol, and no crash matrix for partial installation.
- `F3C-ARCH-1` names `F3C-PATH-PRED-1` but never defines that path-hash domain.
- `archiveHash` is defined over a manifest with itself omitted, but `predecessorManifestDigest`, `v1RecordDigest`, and `v1ReceiptDigest` have no exact hash domains or canonical byte rules.
- The predecessor archive is defined only for v1→v2. The ordinary v2 commit has no equivalent exact byte archive, which is the A1 defect.
- The explicit no-in-place-rollback policy is coherent, but it must be distinguished from frozen Q3, which requires import to be journaled and reversible. R4 says imports are v2 but gives no import journal, identity, content hash, reversibility, or source preservation rule.

**Counterexample:** A migration installs six of eight v2 targets, then crashes before the record. §6.5 says restore v1 from the archive. If the archive is verified, restoration may work, but the contract has not defined the restoration atomicity or how a partially installed optional target is removed. If the archive is missing, the recovery table blocks; if the archive is present but corrupt, it reports evidence gap. Those outcomes are distinguishable, but the target-install step remains non-atomic.

**Required correction:** Use one commit order for ordinary saves and migrations. Define a target-set package and its durable installation protocol. Define every path-hash and digest domain. Add an import section satisfying Q3, or explicitly state that import is outside R4 and remains blocked.

**Gate rows:** V4, V7, V11, V13, V23.

---

### A10 — High: wipe inventory omits frozen roots and temporary evidence

**Exact R4 rule:** §6.7 stages `f3c-predecessors/`, `chats/`, `f3c-commits/`, `f3c-receipts/`, `f3c-tombstones/`, `f3c-creation-intents/`, `f3c-generation-intents/`, `f3c-memory-migrations/`, `f3c-authority/`, the cast journal if present, and the gallery root.

**Frozen rule:** Master §4.1 includes `f3c-root-identity.json`, `f3c-wipe-intent.json`, `config.json`, and `f3c-tmp/` in the root layout. Master §8.10 requires a root-identity commit and a complete staging gate.

**Defects:**

- R4 does not list `f3c-root-identity.json` in the wipe plan. The wipe’s own commit point is the new root identity; without staging/replacing the old identity, the old root can remain current or the new identity can be installed while old identity evidence survives.
- `f3c-wipe-intent.json` is not in the list. The wipe’s durable state machine is then outside its own staged tree.
- `f3c-tmp/` is not in the list. Temporary target/receipt/commit artifacts are exactly the partial state that must not survive a wipe.
- `f3c-generation-intents/` is listed without a path schema or hash domain. It cannot be verified as a complete root.
- `config.json` is not addressed. If it is intentionally preserved, that is a decision that must be stated; if it is owned user content, the wipe is incomplete.

**Counterexample:** A wipe stages the listed roots and installs a new root identity. An old `f3c-tmp` directory contains a staged commit record or target bytes. No inventory scans it, so the broken prior state survives under the new incarnation. Alternatively, a crash after the new identity leaves `f3c-wipe-intent.json` from the old root and no defined reconciliation.

**Required correction:** Define the complete wipe inventory against Master §4.1 and every R4-added root. State which roots are staged, replaced, preserved, or excluded and why. Add `f3c-root-identity.json`, `f3c-wipe-intent.json`, `f3c-tmp`, `f3c-authority` subpaths, and the generation-intent schema. Add exact crash outcomes for each.

**Gate rows:** V8, V12, V13, V14, V22.

---

### A11 — High: the v2 target/record schema is not complete enough to verify exact bytes

**Exact R4 rule:** §10 says R4 “retains the R3 v2 target set and serialization rules” and adds fields. §10.3 says new fields are inserted “at the contract-declared order before the final hash,” but no final target hash is defined there. §10.4 names `f3c-target-manifest-2` but does not define its exact hash prefix, target entry schema, or canonicalization. §2.6 defines a commit record but omits the frozen `receiptBodyBytes`, `receiptBodyHash`, `commitRecordHash`, `receiptEnvelopeHash`, and operation-type fields.

**Defects:**

- The exact v2 target byte format is not in R4; it is deferred to R3. A freeze review of R4 must still verify that the referenced R3 bytes satisfy the frozen target requirements and the new fields. R4 does not state that check.
- The v2 commit record is neither the frozen record nor a complete replacement. It has no record hash or receipt body, so `F3C-CANON-1` record verification cannot apply.
- `meta.historyHead`, `meta.castEpoch`, and `meta.lastCommitOperationId` are added without exact presence/omission rules, field order relative to the existing frozen meta order, or hash treatment.
- `history.generation` omits `defaultSetId`, the execution snapshot, `castEpoch`, and `stateSessionId`; V18 claims a turn’s identity can be recomputed from the target, but required inputs are absent.
- `F3C-CAND-2` uses `targetManifestDigest` but not the exact target manifest bytes. That is a digest-only identity, which is insufficient for A1 reconstruction.
- The “final hash” reference has no defined hash domain; it can be read as a receipt hash, a manifest hash, or a target hash.

**Counterexample:** Two compliant writers serialize a v2 `meta.json` with the same fields in different order because R4 says “contract-declared order” without declaring the new order. Their `tgt2-` hashes differ. V4 cannot establish exact target bytes, and V18 cannot recompute provenance.

**Required correction:** Freeze the complete v2 target schema and byte order, the v2 manifest and hash domain, and the complete v2 commit/receipt record schema. Do not defer these to an unstated R3 rule. Add `defaultSetId`, execution snapshot, `castEpoch`, and state session where the verification claim requires them.

**Gate rows:** V1, V4, V11, V13, V18, V19, V22, V23.

---

### A12 — High: `F3C-MODEL-1` and `F3C-EXEC-1` still leave inference-affecting configuration and component qualifiers undefined

**Exact R4 rule:** §4.1 defines `modelId` over `components: [{ role, artifactHash, byteSize, format, required }]`. It says duplicate roles are permitted only “when each has a distinct role qualifier,” but the envelope has no qualifier field. `format` values are not enumerated.

§4.2 defines `engineConfigDigest` as “`F3C-REQ-1` over effective engine/runtime controls that can change output.” It lists `engineName, engineVersion, engineBuildId, runtimeName, runtimeVersion, executableHash, acceleration, engineConfigDigest, adapterLoadingMode, templateFingerprint, provenanceStatus`.

**Defects:**

- No `roleQualifier` field exists to make the duplicate-role rule enforceable.
- `format` is unconstrained, so two implementations can label the same bytes differently and produce different canonical envelopes.
- `required: false` and omission are not fully distinguished in the identity rule.
- The load/context/tokenizer controls that can affect output are not enumerated: context size, GPU layers, threads, mmap/mlock, memory fraction, RoPE scaling, batch size, tokenizer add-BOS/add-EOS/pad behavior, and offload settings. `engineConfigDigest` is a blanket reference, not a schema.
- `engineConfigDigest = F3C-REQ-1(...)` risks recursion because `F3C-REQ-1` includes `executionFingerprint` and `F3C-EXEC-1` contains the digest.
- The target `history.generation` has `executionFingerprint` but no location for the `F3C-EXEC-1` snapshot that §4.2 says historical turns retain.

**Counterexample:** The same GGUF and projector are loaded with different context size or RoPE scaling. R4’s `modelId` is identical, and the unspecified `engineConfigDigest` may or may not differ. Two turns with materially different inference conditions are indistinguishable, or two implementations produce different execution fingerprints for the same conditions.

**Required correction:** Add `roleQualifier`, close or version the `format` enum, define optional/absent semantics, enumerate load/context/tokenizer controls, break the `engineConfigDigest` recursion, and store the exact execution snapshot with the turn (or in the commit record) with a defined schema.

**Gate rows:** V1, V18, V19, V21, V22.

---

### A13 — High: the freeze boundary expands Gate v2 from V1–V17 to V1–V24 without adopting that expansion

**Exact R4 rule:** §12 is titled “Gate v2 V1–V24” and provides rows V18–V24. §12.2 says “Gate v2 itself is unchanged and unfrozen.”

**Frozen rule:** Gate v2 §6.1 has an adopted row inventory of **V1–V17**. V18–V24 are not Gate v2 rows. The contract’s §11 evidence program has 17 items; R4’s V18–V24 are new verification categories.

**Defect:** R4 cannot simultaneously claim the authority of “Gate v2” and add seven rows without a Gate v2 amendment or explicit proposal status. A reader can treat V18–V24 as adopted rows, while an implementation can claim Gate v2 only requires V1–V17. V24 is especially important because the frozen contract has no baseline row.

**Required correction:** Label V18–V24 as **proposed Gate v2 additions**, not adopted Gate v2 rows. State the owner decision required to adopt them. Keep V1–V17 as the only adopted inventory until that decision.

**Gate rows:** V18–V24.

---

### A14 — Medium: the response-batch design may change multi-character context and interaction behavior

**Exact R4 rule:** §3.4 commits in `participantIndex` prefix order. §11.1 says R4 does not change character-specific prompt construction, context flushing, memory isolation, cast synchronization, or interaction ordering at the prompt/context layer. §11.2 requires a pre-amendment baseline, but no baseline exists.

**Defect:** Commit order is not generation order. If character B generates from the user base while character A is still generating, B’s prompt cannot include A’s newly committed assistant turn. If the existing product generates sequentially and B sees A’s turn, R4 changes prompt context without changing prompt construction. If the existing product also generates concurrently, V24 must prove that. R4 does not state which behavior is required.

The prefix commit rule can also delay B’s visible result until A finishes, even though B finished earlier. That is a user-visible interaction change. R4 acknowledges the delay as an owner decision but does not make it a requirement with a deadline/abandonment rule.

**Counterexample:** Two characters respond to one user turn. A is slow; B is fast. B’s candidate is complete but cannot commit. A fails. Under R4, B’s commit order is after A’s terminal slot; the user sees no B response until A’s failure is resolved. No rule defines the maximum wait or the visible state while waiting.

**Required correction:** Capture the pre-amendment generation/context behavior and decide explicitly whether later slots see earlier assistant turns. Define the batch deadline and terminal states. V24 must compare prompts/context as well as persisted order.

**Gate rows:** V15, V24.

---

## 4. Independent assessment of the required areas

### 4.1 MutationLease authority

**Result: not sufficient.**

R4 correctly makes staging, autosave, append, replacement, commit records, receipts, migration, tombstones, wipe, and recovery lease-bound. That addresses the policy shape of M1. It does not make the lease a complete authority:

- No expiry, renewal, heartbeat, or holder token (A3).
- No `stateSessionId`/`stateEpoch`, so state-only writes survive switches (A3).
- No reconciliation of the authority and commit gate with frozen `stateApplyMu`/`chatWriteMu` (A2).
- No durable target payload tied to the commit record (A1).
- No exact target/record path for ordinary v2 commits (A11).

A lease that can be held forever, lacks a state-session fence, and is not ordered against the frozen locks is not one enforceable durable-write authority. It is a new authority beside the old one.

### 4.2 Commit-record authority and crash consistency

**Result: not crash-recoverable as specified.**

The record is declared the linearization point, but the target publication is a later projection. The record has a manifest digest and a receipt ID, not the target bytes or a durable reconstruction pointer. A crash after the record and before publication cannot be recovered from the record alone. The proposal has recovery rows for migration, but not for an ordinary valid record with a missing target. A repeated commit can return the existing proof without verifying that the projection exists. A malformed or absent target has no defined recovery outcome.

The v2 record is also not the frozen record: it lacks receipt body bytes, record/envelope hashes, operation type, and committed revision. R4 does not provide a replacement hash chain. The record is therefore neither fully self-verifying nor tied to the frozen inventory.

### 4.3 Destructive barrier

**Result: improved but not complete.**

The barrier now closes new leases, cancels producers, drains leases and the commit gate, reconciles, seals the epoch, and only then tombstones or stages. It correctly treats producer termination as insufficient and keeps `blocked_timeout` from authorizing deletion. That is a genuine improvement.

Remaining failures:

- A stuck lease has no expiry/renewal, so a timeout can be permanent without a defined resolution (A3).
- The barrier does not reconcile with `castApplyMu`, `chatSwitchMu`, `chatTreeMu`, `stateApplyMu`, and `chatWriteMu` (A2).
- The wipe inventory omits root identity, wipe intent, and temporary evidence (A10).
- `f3c-generation-intents` has no path schema or hash domain (A10).
- A crash after seal before tombstone resumes destruction, but the exact target/record projection state is not defined (A1).
- A crash after record but before target publication is not covered by the barrier reconciliation table.

No stale operation is *shown* able to write after seal in the narrative, but the lock integration and lease lifecycle are not precise enough to prove it.

### 4.4 Concurrent character execution and response batches

**Result: improved but not deterministic or visibly defined.**

The separation between independent generation and a per-chat commit coordinator is the right architecture. Prefix order and exact-parent checks prevent the specific R3 lost-update case. Duplicate suppression and stable assistant IDs are useful.

The proposal still lacks:

- a durable response-plan location and schema (A5);
- durable attempt/candidate locations (A5);
- a complete terminal-state set, including `abandoned`/`skipped`/`terminal_no_turn` (A5);
- a deadline operation and partial-batch visibility rule (A5);
- `castEpoch` binding to candidate and commit (A4);
- a rule for whether later characters see earlier assistant turns (A14);
- recovery for a lower slot that never starts (A5).

Independent character thinking is preserved in principle. Shared-history commits are serialized in principle. The contract still does not define a deterministic, visible outcome for all batch states.

### 4.5 Character-removal and chat-wide fencing

**Result: partially resolved, not enforceable.**

`characterVersion` is a useful attempt fence, and the removal slot is terminal in the narrative. But the candidate, commit record, and adoption proof do not carry `castEpoch` or `stateSessionId`. Removal can advance `castEpoch` without changing `characterVersion`, leaving a stale candidate eligible. Recreation with the same display name is prohibited in prose but not made a candidate/record invariant. Chat deletion and wipe dominate only if the frozen lock order and the barrier are actually integrated; A2/A3 leave that unproven.

### 4.6 Model identity and execution provenance

**Result: improved, still not exact.**

The separation between model artifact identity and execution provenance is correct. The component roles cover the major artifacts. The remaining defects are concrete:

- no `roleQualifier` field despite a duplicate-role rule (A12);
- undefined `format` values;
- unenumerated load/context/tokenizer controls;
- circular/blanket `engineConfigDigest`;
- no target field for the exact execution snapshot;
- no exact model identity for a component whose absence is a distinct mode.

Historical provenance after removal, corruption, replacement, and backend upgrade is directionally correct, but the record needed to reproduce a turn is not fully stored.

### 4.7 Request canonicalization

**Result: not stable across implementations.**

`F3C-REQ-1` has a useful full parameter table and default-resolution order. It still fails on the canonical form:

- RFC 8785 key ordering is misstated (A6).
- `candidateCreatedAt` makes candidate identity time-dependent (A6).
- `origin` is described but not present in the envelope (A6).
- extension keys can duplicate named parameters (A6).
- `requestFingerprint` collides with the frozen operation fingerprint (A7).
- `defaultSetId` is not carried into the target generation record, so V18 cannot recompute it from the target (A11).

Two compliant implementations can produce different request/candidate hashes. That is a freeze blocker for provenance and duplicate suppression.

### 4.8 V1/V2 predecessor archive

**Result: policy improved, transaction not exact.**

The archive-before-v2 policy, exact-byte copy, manifest hash, required `predecessorArchiveHash`, strict v2 unknown-field handling, and explicit no-in-place-rollback choice are coherent. The archive is correctly not treated as a rollback promise.

Remaining defects:

- `F3C-PATH-PRED-1` is undefined (A9).
- `predecessorManifestDigest`, `v1RecordDigest`, and `v1ReceiptDigest` lack exact domains (A9).
- Migration target installation is ordered differently from ordinary commit (A9).
- Multi-file “atomic” installation has no primitive or crash matrix (A9).
- Ordinary v2 commits have no target-byte archive equivalent (A1).
- Import is named as v2 but not defined as journaled/reversible under frozen Q3 (A9).

No direct conflict with the no-rollback choice is established; the choice is an owner decision. The *transaction* that makes the archive usable is not complete.

### 4.9 Continuation preference and model switching

**Result: policy improved, durable storage is contradictory.**

`SelectContinuationModel`/`SetContinuationPreference`, terminal-only `LoadModel`, and the separation of requested/active/saved are the right shape. The storage contradiction is fatal to the claim:

- saved preference is in `meta.json` (A8);
- the write is said not to be a chat target write (A8);
- no synchronization or projection rule exists;
- no expected-sequence CAS exists.

A newer explicit preference can be overwritten by a stale autosave, or the root preference and meta can disagree. The native UI evidence is absent, so V16/V21 remain unverified.

### 4.10 Media and F3-F boundary

**Result: boundary holds.**

F3-C core can proceed without F3-F; F3-C full acceptance and V20 remain blocked. References are not durability. Missing, corrupt, unavailable, or externally owned media is reported rather than substituted. F3-C does not write or delete F3-F bytes. This area is properly scoped, but V20 is blocked and no F3-F evidence exists.

### 4.11 Multi-character preservation

**Result: not proven; one contract ambiguity remains.**

R4 states the right non-goals: prompt construction, context flushing, memory isolation, cast semantics, lore visibility, and interaction ordering are not changed. The persistence fields are additive. However, concurrent generation plus prefix commit can change what a later character sees (A14), and no pre-amendment baseline exists. V24 cannot pass without that baseline and without a rule for ordering/visibility.

### 4.12 Gate v2 V1–V24

**Result: not auditable as a Gate v2 record.**

R4 names 24 rows, but Gate v2 only adopts V1–V17 (A13). Of the proposed rows, V4 itself says the target schema is not frozen; V6 has a conflicting lock order; V13 lacks key boundaries; V18–V24 are not adopted Gate rows; V20 is blocked. No row can be PASS from prose. The full row audit is in §7.

---

## 5. Complete prior-finding disposition

The statuses below are this review’s dispositions, not R4’s self-report. “Resolved in principle” means the policy direction is defensible but evidence and exactness are still required.

### 5.1 R1-review findings R1–R10

| ID | R4 claim | Independent disposition | Why |
|---|---|---|---|
| R1 | Intents and operation roots under authority and wipe | **Partially resolved** | Intent root has no exact path/schema; wipe list is incomplete (A10). |
| R2 | Attempts bind user/assistant commit, batch, candidate, receipt, session | **Partially resolved** | Fields exist, but attempt/candidate/plan have no durable location and no exact state machine (A5). |
| R3 | One recovery precedence; ambiguity blocks | **Partially resolved** | Precedence exists, but plan/attempt/candidate records cannot be found by a defined path (A5). |
| R4 | Transition state machine, readiness, failure matrix | **Partially resolved** | Preference storage contradiction and no native evidence (A8). |
| R5 | Effective request captured; model/execution split | **Partially resolved** | Canonicalization and execution-config gaps (A6, A12). |
| R6 | Exact archive; strict v2/lossless predecessor | **Partially resolved** | Archive path/domains and migration order incomplete (A9). |
| R7 | Core/full split; V20 blocked | **Resolved in principle** | Boundary is correct; evidence absent and V20 blocked. |
| R8 | Non-goal and V24 baseline | **Partially resolved** | Non-goal stated; baseline absent and concurrent context unresolved (A14). |
| R9 | Presence, ordering, null/omission rules | **Partially resolved** | v2 target order and hash domains are not complete (A11). |
| R10 | All assistant paths use attempt protocol | **Resolved in principle** | Paths are named; exact protocol and persistence are not (A5). |

### 5.2 R2-review findings N1–N15

| ID | R4 claim | Independent disposition | Why |
|---|---|---|---|
| N1 | Absent chat without tombstone is unknown | **Resolved in principle** | Rule is explicit; still needs V8/V22 evidence. |
| N2 | Commit existence tied to attempt | **Partially resolved** | Candidate/attempt records are not durably specified; candidate identity is nondeterministic (A5, A6). |
| N3 | Delete/wipe quiesce stream | **Partially resolved** | Barrier is improved; lease lifecycle and frozen locks are incomplete (A2, A3). |
| N4 | Migration record/target order resolved | **Contradicted by a new rule** | §2.4 and §6.4 define opposite target/record orders (A9). |
| N5 | Model identity deferred | **Partially resolved** | More components, but qualifier/format/config gaps remain (A12). |
| N6 | Parameter set incomplete | **Partially resolved** | Full table exists, but origin/extension/canonicalization defects remain (A6). |
| N7 | Strict unknown fields brick chats | **Resolved in principle** | Archive plus strict v2 is a coherent policy; archive exactness is still open (A9). |
| N8 | Predecessor bytes not retained | **Partially resolved / owner decision** | Bytes are retained, but no automatic in-place rollback is chosen; accepted as a policy only if the owner approves (A9). |
| N9 | Preference-update ambiguity | **Partially resolved** | Explicit operations exist, but storage is split and CAS is absent (A8). |
| N10 | Gate V1–V24 only partly mapped | **Partially resolved** | Rows exist, but Gate v2 adoption and row verifiability are incomplete (A13). |
| N11 | One attempt per chat changes multi-character | **Partially resolved** | Concurrency is allowed, but context/order/deadline behavior is unresolved (A5, A14). |
| N12 | Root mismatch cleans evidence | **Resolved in principle** | Blocking is explicit; evidence still needed. |
| N13 | Candidate fingerprint undefined | **Contradicted by a new rule** | `candidateCreatedAt` destroys deterministic candidate identity (A6). |
| N14 | New-chat format unspecified | **Resolved in principle** | v2 is explicit; exact v2 bytes are not (A11). |
| N15 | V24 has no baseline | **Partially resolved** | Baseline is required but not captured; V24 is blocked (A14, §7). |

### 5.3 R3-review findings M1–M10

| ID | R4 claim | Independent disposition | Why |
|---|---|---|---|
| M1 | Every durable continuation leases; barrier drains leases/gates | **Partially resolved** | Policy is right; no expiry/token/state-session and no frozen-lock bridge (A2, A3). |
| M2 | Concurrent thought, serialized commit, response plan | **Partially resolved** | Coordinator exists; plan/candidate locations and terminal states are missing (A5). |
| M3 | Full model artifacts; execution separate | **Partially resolved** | Roles are broad but qualifier/format/config exactness is missing (A12). |
| M4 | Request canonicalization defined | **Contradicted by a new rule** | RFC 8785 ordering conflict and `origin` ambiguity (A6). |
| M5 | Archive integrity, rollback, post-commit | **Partially resolved** | Archive policy exists; transaction and exact domains are incomplete (A9). |
| M6 | `F3C-CAND-2` full specification | **Contradicted by a new rule** | Envelope includes `candidateCreatedAt`; same candidate cannot have same fingerprint (A6). |
| M7 | Durable `blocked_timeout` recovery | **Partially resolved** | State exists; lease expiry/reconciliation and exact restart behavior are incomplete (A3). |
| M8 | Character removal/version fence | **Partially resolved** | `characterVersion` exists; `castEpoch`/state session are not bound (A4). |
| M9 | Gate claim/evidence/oracle/PASS/FAIL | **Partially resolved** | Rows are detailed, but V18–V24 are not adopted and several oracles are contradictory (A13). |
| M10 | Preference sequencing and explicit operations | **Partially resolved** | Sequence concept exists; storage split and CAS are missing (A8). |

No prior finding is closed by this review. R4’s “Resolved in contract” labels are too strong for M4, M6, N4, and N13, and overstate the evidence status of the rest.

---

## 6. New counterexamples and divergent-compliance cases

1. **Record before target:** crash after record, before target publication; valid record, old target, no target bytes (A1).
2. **Stuck lease:** autosave lease never expires; delete barrier times out forever; no safe abandonment (A3).
3. **Stale session:** user switches chats; blocked autosave resumes with a valid base and commits old state because the lease has no `stateSessionId` (A3).
4. **Removed character:** removal advances `castEpoch` but not `characterVersion`; candidate still matches and commits (A4).
5. **Missing lower batch slot:** slot 2 never starts; slot 3 waits forever because there is no `terminal_no_turn` record; `abandoned` is not a valid status (A5).
6. **Partial batch:** slot 1 commits, slot 2 is in flight, crash; no rule for visible state or retry (A5).
7. **RFC key order:** non-BMP property order differs between code-point and UTF-16 canonicalization (A6).
8. **Candidate timestamp:** identical content, different time, different candidate fingerprint (A6).
9. **Fingerprint namespace:** generation fingerprint overwrites/aliases frozen operation fingerprint (A7).
10. **Preference overwrite:** autosave writes old `meta.continuationModel` after explicit preference save; meta and root disagree (A8).
11. **Migration partial install:** v2 targets installed before record; crash; restore depends on archive and multi-file atomicity that is not defined (A9).
12. **Wipe omission:** old root identity or `f3c-tmp` survives the staged wipe (A10).
13. **Target schema divergence:** two writers insert new `meta`/`history` fields in different declared orders (A11).
14. **Model identity collision:** same artifacts, different context/rope/tokenizer settings, same or unspecified execution fingerprint (A12).
15. **Gate scope:** V18–V24 treated as adopted Gate rows although Gate v2 only adopts V1–V17 (A13).
16. **Character context change:** fast character B generates before slow A commits; B’s prompt omits A’s turn while commit order puts B after A (A14).

These are not theoretical polish items. Each can produce a different durable history, a different provenance record, or a different visible result under two implementations that both claim R4 compliance.

---

## 7. Gate V1–V24 audit

Gate v2’s adopted inventory is V1–V17. V18–V24 are R4 proposals, not adopted Gate v2 rows. No row is PASS. This is a contract review, not a test run.

| Row | Status | Independent reason |
|---|---|---|
| V1 | **UNRESOLVED** | `F3C-CANON-2` contradicts RFC 8785 key ordering; `origin` and extension precedence are undefined (A6). |
| V2 | **UNRESOLVED** | Valid record without target has no reconstruction path; receipt-only reconstruction is not enough (A1). |
| V3 | **UNRESOLVED** | No exact candidate/record target-set schema; fork/orphan oracle cannot distinguish record-without-target (A1, A11). |
| V4 | **BLOCKING** | R4 itself says the v2 target schema is not frozen; exact bytes/order/domains are absent (A11). |
| V5 | **UNRESOLVED** | Lease lacks `stateSessionId`/`stateEpoch`; stale autosave not proven to fail (A3). |
| V6 | **BLOCKING** | R4 and Master define incompatible lock orders; no combined order or native stress (A2). |
| V7 | **UNRESOLVED** | Save As/import scope is not defined; Q3 import reversibility is absent (A9). |
| V8 | **UNRESOLVED** | Barrier improved, but lease expiry and wipe/tombstone evidence are incomplete (A3, A10). |
| V9 | **UNRESOLVED** | Creation journal is referenced but not specified in R4; no exact creation schema (A11). |
| V10 | **UNRESOLVED** | No new memory-migration evidence; inherited contract only. |
| V11 | **UNRESOLVED** | Archive and history schema/domains are incomplete (A9, A11). |
| V12 | **UNRESOLVED** | Wipe inventory omits root identity, wipe intent, tmp, and exact generation-intent schema (A10). |
| V13 | **UNRESOLVED** | Fault matrix omits record-without-target, lease expiry, response-plan/candidate boundaries (A1, A3, A5). |
| V14 | **BLOCKED** | Native Windows reparse/hardlink evidence is absent. |
| V15 | **UNRESOLVED** | No deterministic race harness; batch/context/lease states are not complete (A3–A5, A14). |
| V16 | **BLOCKED** | Native Wails/WebView2 evidence is absent. |
| V17 | **BLOCKED** | No three-stage Mimo challenge evidence. |
| V18 | **UNRESOLVED** | Model/execution identity and target provenance are incomplete (A11, A12). |
| V19 | **UNRESOLVED** | Defaults/canonical bytes/seed status are contradictory or incomplete (A6). |
| V20 | **BLOCKED** | F3-F ownership/integrity/deletion/backup evidence is absent; V20 remains blocked. |
| V21 | **UNRESOLVED** | Saved preference has two stores and no CAS; no native evidence (A8). |
| V22 | **UNRESOLVED** | Candidate/plan/record recovery is not exact or reconstructable (A1, A4, A5). |
| V23 | **UNRESOLVED** | Migration order, archive domains, and rollback/import policy are incomplete (A9). |
| V24 | **BLOCKED** | No pre-amendment baseline; concurrent context behavior unresolved (A14). |

**V20 remains BLOCKED, not failed.** V14, V16, V17, and V24 are blocked by missing native/Mimo/baseline evidence. The remaining rows are unresolved contract or evidence defects, not passes.

---

## 8. Owner decisions still required

These decisions cannot be inferred from R4, and this review does not make them on the owner’s behalf.

1. **Commit evidence model:** retain the exact target-set package until projection, or make the record contain reconstructable target bytes. A record that only identifies a digest is insufficient.
2. **Authority integration:** amend the frozen lock order to include the mutation authority and commit gate, or define the lease as a subordinate of the existing `stateApplyMu`/`chatWriteMu` order.
3. **Lease lifecycle:** choose bounded lifetime, renewal, holder token, expiry, and general startup reconciliation.
4. **State-session fence:** decide whether `stateSessionId`/`stateEpoch` are mandatory in every lease and commit record.
5. **Character removal:** choose `characterVersion`, `castEpoch`, or both as the authoritative candidate fence.
6. **Response-batch semantics:** define plan storage, candidate/attempt storage, every terminal status, deadline behavior, and partial-batch visibility.
7. **Multi-character context:** decide whether concurrent later slots see earlier assistant turns; capture the pre-amendment baseline.
8. **Canonicalization:** choose RFC 8785’s UTF-16 ordering or a different explicitly named canonical form; remove `candidateCreatedAt` from stable candidate identity.
9. **Fingerprint namespace:** keep the frozen operation `requestFingerprint` and add a separate generation-request fingerprint.
10. **Preference storage:** choose `meta.json` as the sole durable source, or a root projection with a defined repair/ordering rule and CAS.
11. **Migration/import:** choose one target/record order, define archive path/domains, and decide whether import is in TARGET-2 or remains blocked under Q3.
12. **Wipe inventory:** include or explicitly exclude every frozen root, especially root identity, wipe intent, `f3c-tmp`, and config.
13. **Model/execution identity:** define role qualifiers, format values, load/context/tokenizer controls, and execution snapshot storage.
14. **Gate scope:** formally adopt V18–V24 as Gate v2 additions or keep Gate v2 at V1–V17 and treat the rest as proposed evidence categories.

---

## 9. Evidence classification and integrity

| Class | Present | Note |
|---|---|---|
| 1. Contractually established | Yes | Frozen Master, Gate v2, prior review records. |
| 2. Source-inspected | Yes | Read-only current paths: `LoadModel` returns before readiness (`app-model.go:460` region); `SaveChat` writes targets sequentially (`persistence.go:190` region); `DeleteChat` removes the directory directly (`persistence.go:471`); `performSave` reloads and writes (`app-chat.go:314`); stream and character paths append then autosave (`app-prompt.go:138`, `:230`). |
| 3. Test-proven | No | No tests run or supplied. |
| 4. Runtime/native-verified | No | No Wails/WebView2, native delete, model load, or media run. |
| 5. Unverified/inferred | Yes | R4 and all proposed Gate rows. |

**Integrity results:**

- R4 SHA-256 matches: `7a2c221147621706e299924c20de1fe29080b7d649fac25264ce212d11e7417a`.
- Frozen authority SHA-256 matches: `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`.
- R3 SHA-256 matches: `5013eb6fad7e0ec6f067d0fdbd350fd7f7c68ff9e263a7967b37f39e6daaf42c`.
- R4 line count: 1496; fenced markers: 60; balanced.
- Prior reviews: 543, 616, and 461 lines; hashes recorded in the header.
- No R4, frozen-contract, prior-review, source, test, runtime, or application-data file was modified.

The current source confirms the gaps R4 names as future work: `LoadModel` returns `nil` after spawning while the monitor later emits `model:loaded`; `SaveChat` writes `meta.json`, `history.json`, world/scenario, characters, lorebooks, memories, and config sequentially; `DeleteChat` removes the chat directory without a tombstone or barrier; `performSave` reloads completed chat data and writes it asynchronously; streaming and character paths append and autosave after stream completion. These are implementation gaps, not evidence that R4 is correct.

---

## 10. Status and authorization

**Question:** Does R4 establish one enforceable durable-write authority, safe concurrent character execution, and crash-recoverable persistence sufficiently precisely to support a formal architecture freeze?

**Answer: No.**

R4 establishes the right *shape* for a single authority, but not an enforceable one. It establishes concurrent thinking plus serialized commits, but not deterministic batch lifecycle or context behavior. It establishes record-based recovery in prose, but not a record from which the committed target set can be reconstructed. Those are freeze-level defects, not implementation polish.

**Status:**

- **R4 proposal:** **NOT FROZEN; NOT ACCEPTED; NOT IMPLEMENTED; NOT VERIFIED.**
- **Frozen authority:** unchanged; SHA-256 verified.
- **F3-C phase 0:** **NOT CLOSED.**
- **P1a:** **UNAUTHORIZED**, downstream of P1c.
- **F3-C acceptance:** **NOT GRANTED.**
- **F3-D:** **BLOCKED.**
- **F3-B:** closed; not reopened.
- **A9/A10/A12/A13:** closed; not reopened.
- **V20:** **BLOCKED** pending F3-F evidence.
- **Gate v2:** adopted V1–V17; V18–V24 remain proposed additions.

No implementation authorization exists. No source, test, build, migration, model load, destructive test, application data, or authority record was changed by this review.
