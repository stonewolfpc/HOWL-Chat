# REVIEW — Independent Adversarial Review of F3C-TARGET-2 R3

**Status:** **INDEPENDENT ADVERSARIAL REVIEW — NOT APPROVAL, NOT FREEZE, NOT IMPLEMENTATION.**  
**Date:** 2026-11-11  
**Reviewed artifact:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R3-2026-11-11.md`  
**R3 SHA-256:** `5013eb6fad7e0ec6f067d0fdbd350fd7f7c68ff9e263a7967b37f39e6daaf42c`  
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`  
**Prior reviews:** R1 review and R2 fresh review are used for finding history, not as proof.  
**Scope:** read-only review. No source, tests, builds, model loads, migrations, destructive tests, authority edits, implementation, or freeze.

---

## 1. Overall verdict

**No — R3 does not yet define a safe, internally consistent, testable contract suitable for a formal freeze decision.**

R3 is a substantial improvement. It closes the *policy* gaps that the previous reviews identified:

- absent chat now fails closed;
- root-incarnation mismatch blocks;
- exact assistant commit identity is required;
- deletion/wipe have a generation barrier;
- v1 predecessor bytes are retained;
- new chats are v2;
- model identity and request fingerprint are no longer purely deferred;
- per-character execution scopes replace chat-wide serialization;
- Gate V1–V24 is mapped.

But R3 still leaves two **critical concurrency defects**:

1. the destructive barrier fences stream producers, but not necessarily all downstream persistence continuations that can still commit;
2. per-character scopes allow concurrent attempts to commit against the same base without a defined merge/order policy.

Those defects can still produce lost updates, duplicate assistant turns, or a late write after deletion. R3 also leaves model identity, request canonicalization, and predecessor-archive integrity partially defined, so two compliant implementations can still diverge.

**Evidence note:** R3 is a proposal; all R3 rules are level 5 (unverified) unless restating a frozen rule. No tests or runtime/native evidence were produced. No Gate row passes.

---

## 2. Evidence classification

| Class | Meaning | Present in this review |
|---|---|---|
| 1. Contractually established | Frozen Master contract rule. | Yes. |
| 2. Source-inspected | Read-only current-source behavior. | Yes, for stream/delete/wipe/model/character paths. |
| 3. Test-proven | Reproducible positive/negative tests. | None. |
| 4. Runtime/native-verified | Real Wails/WebView2/media run. | None. |
| 5. Unverified/inferred | Proposed rule or design inference. | R3 and all Gate rows. |

No R3 proposal is promoted to established fact in this review.

---

## 3. Severity-ranked findings

### M1 — The destructive barrier does not fence queued persistence continuations

**Severity: Critical. Defect class: contract/implementation.**

**Exact R3 rule:** §2.6 barrier phases: close, cancel, await producer, reconcile, destructive step, release. §2.6 says a write after the barrier epoch is rejected. §2.4 says recovery verifies exact commit identity.

**Invariant:** Once the barrier is closed, no write path capable of creating a target, receipt, directory, or assistant turn may commit until the destructive operation has a durable outcome.

**Adversarial counterexample:** A stream producer is canceled and terminates. Before cancellation, it queued an autosave/commit continuation, an `AppendChatMessage` continuation, or an asynchronous save goroutine. The barrier sees the producer terminal and proceeds to tombstone/delete. The queued continuation does not check the barrier epoch and commits an assistant turn or autosave into the deleted chat. Producer termination was true; downstream work was not fenced.

A second variant: a commit is already in progress when the barrier starts. R3 says “reconcile a commit that completed before the barrier,” but does not define whether an in-flight commit that has not yet published is adopted, aborted, or blocked. Different implementations will choose different outcomes.

**Does R3 survive?** No. It fences the producer, not every durable mutation continuation.

**Required correction:** Define a single persistent-authority fence token/epoch checked by every write path:

- generation append;
- autosave and final save;
- history rewrite/clear;
- context/cast writes;
- commit record/receipt installation;
- target installation;
- migration/predecessor writes.

The barrier must close all of them, not just streams. Any in-flight commit must have a defined pre-barrier/post-barrier outcome. A queued continuation must fail closed, not commit.

**Gate rows:** V5, V6, V8, V12, V13, V15, V22.

---

### M2 — Per-character scopes allow concurrent commits against the same base

**Severity: Critical. Defect class: contract.**

**Exact R3 rule:** §2.3 replaces chat-wide exclusivity with one active attempt per `(chatID, characterId, stateSessionID)`. §2.4 requires a base revision; §7 says different characters may generate concurrently.

**Invariant:** Execution concurrency must not create lost updates, conflicting revisions, duplicate/reordered assistant turns, or stale-session writes.

**Adversarial counterexample:** Character A and character B start concurrently at base revision 5. Both produce assistant candidates. A commits revision 6. B still holds base 5. If B commits, it either overwrites A, fails stale, or appends a second assistant turn against the wrong base. Nothing in R3 defines merge order, expected revision, per-attempt ordering, or whether concurrent character results must be serialized at commit time.

A second variant: A and B both append to the global history. The final ordering depends on scheduler timing, not on the interaction order. A third: one character is removed while its attempt is active; R3’s scope does not include a character version/epoch, so a stale character result can commit after removal.

**Does R3 survive?** No. It distinguishes execution concurrency from chat-wide serialization, but does not define the persistence serialization that concurrent execution requires.

**Required correction:** Define:

- execution concurrency per character scope;
- a single commit serialization point for the chat;
- expected base/merge policy for concurrent character results;
- deterministic assistant ordering (for example, commit sequence plus attempt sequence);
- stale-base behavior that does not overwrite a newer commit;
- character-version/removal fence for a character attempt;
- duplicate suppression per `(characterId, userMessageId, attempt)`.

If concurrent characters must produce one combined assistant response, define the merge. If each produces a separate turn, define ordering and race resolution. Silence is not a contract.

**Gate rows:** V5, V6, V15, V22, V24.

---

### M3 — Model identity omits inference-relevant artifacts

**Severity: High. Defect class: contract.**

**Exact R3 rule:** §4.2 defines `modelId` as SHA-256 over `{model_sha256, projector_sha256}`. It does not include tokenizer/vocabulary, adapters/LoRAs, auxiliary models, quantization/conversion metadata, backend executable/version, or engine/runtime version.

**Invariant:** The identity must distinguish artifacts that can materially change inference, while leaving machine-specific data out of the model artifact identity where irrelevant.

**Adversarial counterexample:** The same base GGUF and projector are used with different LoRAs, tokenizer files, or auxiliary vision components. R3 assigns the same `modelId`; two turns with different behavior look identical. Alternatively, a backend executable upgrade changes sampling or template behavior while `modelId` and `engine/backend` strings remain unchanged.

**Does R3 survive?** No. It captures two artifacts, not the full inference identity.

**Required correction:** Separate:

- **model artifact identity:** model, tokenizer/vocabulary, projector/vision, adapters/LoRAs, auxiliary models, quantization/conversion identity;
- **execution provenance:** backend/engine name and version, runtime/executable identity, template/schema version;
- **request fingerprint:** effective controls.

Do not put machine-specific paths into identity; do include artifact-affecting components.

**Gate rows:** V18, V21.

---

### M4 — `requestFingerprint` canonicalization is undefined

**Severity: High. Defect class: contract.**

**Exact R3 rule:** §4.3 says `requestFingerprint` is SHA-256 over “the canonical effective inference request.” It does not define canonical JSON, number/float encoding, map ordering, array ordering, default resolution, omitted vs null vs zero, or schema-version transitions.

**Invariant:** Semantically identical effective requests must hash identically; materially different requests must not collide.

**Adversarial counterexample:** Two implementations serialize floats differently (`0.7` vs `0.7000000000000001`), default a missing key differently, order `samplerOrder` differently, or encode a stop sequence as `null` versus omitted. Same effective request, different fingerprint; or different requests, same fingerprint after normalization.

**Does R3 survive?** No. A fingerprint without a canonicalization domain is not testable.

**Required correction:** Freeze a canonicalization specification:

- JSON canonical form and scalar encoding;
- resolved defaults before hashing;
- omitted/null/zero rules;
- map key ordering and array semantics;
- numeric precision/rounding;
- schema-version transition rules;
- whether prompt text is excluded and template/adapter identity included.

**Gate rows:** V18, V19.

---

### M5 — Predecessor archive integrity, rollback, and post-commit recovery are incomplete

**Severity: High. Defect class: contract.**

**Exact R3 rule:** §5.2–§5.5 retain `f3c-predecessors/...`, add optional `predecessorArchiveHash` to record/receipt, and define crash recovery through publish. It says predecessor archives are immutable evidence and retained until deletion/wipe.

**Invariant:** The predecessor archive must be cryptographically and structurally bound to the exact v1 target and v2 successor; a missing/corrupt archive must not be mistaken for a complete migration.

**Adversarial counterexample:** The archive is written and verified, v2 commits, then the archive is corrupted or deleted. R3 has no post-commit recovery rule. A rollback claim cannot be honored. A second implementation treats the optional receipt field as sufficient proof even when the archive is absent. A third restores v1 from a partially staged archive without verifying its manifest.

**Does R3 survive?** No. It defines the happy path and crash windows, but not the archive integrity verification/lifecycle after commit.

**Required correction:** Define:

- the exact archive manifest format and canonical hash;
- archive verification on every use and on recovery;
- missing/corrupt archive behavior (`ErrPredecessorArchiveFailed`/evidence gap, no silent completion);
- rollback procedure and whether rollback is supported after v2 commit;
- retention/deletion/wipe behavior;
- binding of archive hash to v1 target manifest and v2 receipt/record.

**Gate rows:** V4, V13, V23.

---

### M6 — `F3C-CAND-1` is named but not fully specified

**Severity: High. Defect class: contract.**

**Exact R3 rule:** §2.5 says the candidate fingerprint is SHA-256 over a canonical envelope containing chat ID, operation ID, base revision, assistant message ID, ordered target manifest, “canonical candidate history/state bytes,” and root/chat identity. It does not define the canonical bytes, the target ordering format, or how state bytes are encoded.

**Invariant:** Recovery must be able to prove the exact candidate the receipt commits; the fingerprint must be stable across implementations.

**Adversarial counterexample:** One implementation hashes the in-memory history, another hashes `history.json` only. One includes `attachments: []`, another omits it. Both claim `F3C-CAND-1`; recovery cannot compare fingerprints, or two different candidates collide.

**Does R3 survive?** No. The domain name is not a specification.

**Required correction:** Freeze the candidate envelope schema and canonicalization, including exact target manifest ordering, history bytes, optional fields, and state encoding. Reuse the same canonicalization used by the receipt/target proof where possible.

**Gate rows:** V13, V18, V22.

---

### M7 — Barrier timeout has no resolution or recovery state

**Severity: High. Defect class: contract.**

**Exact R3 rule:** §2.6 says timeout/uncertain termination blocks deletion/wipe with `ErrGenerationBarrierTimeout`; no tombstone/staging occurs. It does not say whether the barrier remains closed, how it is cleared, or how retry/recovery works.

**Invariant:** A blocked destructive operation must leave a defined, recoverable state; it must not permanently deadlock the chat or permit a later late write.

**Adversarial counterexample:** A producer ignores cancellation. Delete returns `ErrGenerationBarrierTimeout`. The barrier remains runtime-closed but no durable state records why. A later delete retry cannot acquire, or clears the barrier while the old producer is still alive. A crash/restart loses the barrier and the producer’s intent may be adopted.

**Does R3 survive?** No. The timeout branch is a dead end.

**Required correction:** Define barrier lifecycle:

- barrier ID/status and whether it is durable or runtime-only;
- whether timeout leaves the barrier closed;
- how the user resolves/retries;
- how crash/restart reconciles a closed barrier;
- guarantee that a late producer cannot commit after a timeout/retry.

**Gate rows:** V6, V12, V13, V22.

---

### M8 — Character removal has no attempt fence

**Severity: High. Defect class: contract.**

**Exact R3 rule:** §7 allows per-character scopes and says no cast changes. Character removal/removal-during-generation is not addressed.

**Invariant:** A character result must not commit after that character is removed from the chat or its context is invalidated.

**Adversarial counterexample:** Character A is generating. The user removes A from the cast. A’s stream completes and commits an assistant turn with A’s `characterId`, reintroducing a removed character into history or associating context that is no longer valid. The scope key includes `characterId`, but not character version/removal epoch.

**Does R3 survive?** No. The execution scope lacks the character lifecycle fence.

**Required correction:** Add a character/version/epoch fence to the scope or attempt; removal invalidates the attempt; recovery marks it stale/aborted; no assistant commit after removal.

**Gate rows:** V15, V22, V24.

---

### M9 — Gate mapping is present, but pass/fail and evidence definitions are incomplete

**Severity: Medium. Defect class: evidence/test.**

**Exact R3 rule:** §9 maps V1–V24 with claims and negative evidence, but does not define row-level PASS/FAIL/UNRESOLVED/BLOCKED criteria, positive/negative case matrix, native prerequisites, or row dependencies beyond a short phrase.

**Invariant:** Each Gate row must be independently testable and cannot pass from proposal text.

**Adversarial counterexample:** V13 lists “fault injection at barrier, intent, migration…” but not the exact fault points or expected recovery per point. V16 says native preference/restore but not the account/model setup. V17 says Mimo required but not the challenge protocol. A row can be checked off by a test name.

**Does R3 survive?** It is a useful map, not a complete Gate v2 audit.

**Required correction:** For each V1–V24, add exact claim, positive/negative cases, evidence artifact/command, result status, prerequisites/dependencies. Preserve V17. Do not freeze until all rows have dispositions.

**Gate rows:** V1–V24.

---

### M10 — Preference operation races with transition and commit

**Severity: Medium. Defect class: contract.**

**Exact R3 rule:** §3.2 defines `SelectContinuationModel` and `SetContinuationPreference`; §3.3 says one transition at a time; it does not define what happens when a preference operation races a transition, a chat switch, or a preference commit failure.

**Invariant:** Saved/requested/active cannot become inconsistent under concurrent user actions.

**Adversarial counterexample:** A load of B is in progress. The user explicitly sets preference C. B becomes ready and is published as active; its selection path also saves B, overwriting C. Or a preference commit fails while active B is ready; the UI reports B saved when it is not.

**Does R3 survive?** Not fully. It needs an operation ordering/ownership rule.

**Required correction:** Serialize preference/selection/transition operations per chat; define last-writer and conflict behavior; make preference commit failure explicit; no state is reported saved before receipt.

**Gate rows:** V15, V21.

---

## 4. Disposition of the fifteen previous findings

| Prior finding | Disposition | Residual issue |
|---|---|---|
| N1 absent chat | **Resolved in principle.** Fail-closed rule is explicit. | Barrier/root recovery still needs M1/M7 evidence. |
| N2 exact commit | **Partially resolved.** Exact match list exists. | M6 candidate envelope is not fully specified. |
| N3 delete/wipe barrier | **Partially resolved.** Barrier concept exists. | M1 queued continuations; M7 timeout. |
| N4 migration order | **Partially resolved.** Staged sequence and crash table exist. | M5 archive integrity/post-commit. |
| N5 model identity | **Partially resolved.** Concrete digest exists. | M3 omits artifacts/backend version. |
| N6 parameter completeness | **Partially resolved.** Request fingerprint added. | M4 canonicalization undefined. |
| N7 unknown fields | **Resolved in principle.** Predecessor archive preserves exact v1 bytes. | M5 archive lifecycle. |
| N8 predecessor bytes | **Resolved in principle.** Bytes retained. | M5 rollback/verification missing. |
| N9 preference semantics | **Resolved in principle.** Explicit operations. | M10 races. |
| N10 Gate mapping | **Partially resolved.** V1–V24 present. | M9 pass/fail/evidence detail. |
| N11 per-character scopes | **Partially resolved.** Scope replaces chat-wide serialization. | M2 commit conflicts; M8 removal. |
| N12 root mismatch | **Resolved in principle.** Blocks. | Needs evidence. |
| N13 fingerprint semantics | **Partially resolved.** Operation rules added. | M6 domain undefined. |
| N14 new-chat v2 | **Resolved in principle.** Explicit. | Needs Gate V9 evidence. |
| N15 V24 baseline | **Partially resolved.** Baseline required. | No baseline evidence yet. |

No prior finding is fully closed merely because R3 names it.

---

## 5. Independent challenge results

### 5.1 Destructive-operation barrier

**Result: not safe.** R3 fences producers and awaits producer termination, but does not prove that queued persistence continuations, in-flight commits, or async save paths cannot write. The barrier needs a single fence authority checked by every durable mutation (M1). Timeout/recovery is undefined (M7). “Producer terminated” is not “all write authority closed.”

### 5.2 Per-character execution scopes

**Result: not safe.** Execution concurrency is defined; persistence commit serialization is not. Concurrent character attempts can collide on base revision, order, and history (M2). Character removal has no fence (M8). Different character context can remain isolated in prompt construction, but shared history/persistence still needs a merge/order contract.

### 5.3 Model identity and provenance

**Result: incomplete.** SHA-256 over model+projector is a real start, but it omits tokenizer/vocabulary, adapters, auxiliary models, backend executable/version, and runtime configuration. Those can change inference. Model identity and execution provenance must be separated (M3).

### 5.4 Migration and predecessor archive

**Result: incomplete.** The archive-before-install shape is correct. R3 does not fully bind archive content to a canonical manifest, verify it after commit, define rollback after v2, or specify missing/corrupt-archive recovery (M5). Archive existence is not integrity.

### 5.5 Request fingerprint and parameter schema

**Result: incomplete.** `requestFingerprint` is the right idea, but its canonicalization is not specified. Defaults, float encoding, map/array order, null/omitted/zero, and schema transitions are unresolved (M4). A testable hash domain is required.

### 5.6 Model switching and preference

**Result: improved, still racy.** Readiness and fail-closed restore are coherent. Preference/transition races and save-failure reporting are not fully ordered (M10). No native evidence exists.

### 5.7 V1/V2 compatibility and unknown fields

**Result: improved, not complete.** Predecessor bytes preserve unknown fields. Archive lifecycle and post-commit integrity are not fully defined (M5). Strict v2 validation is coherent, but archive retention/deletion/wipe must be explicit.

### 5.8 Generation intent and recovery

**Result: improved, not fully testable.** Exact commit checks are strong. Candidate fingerprint domain is not fully specified (M6). Retry/operation rules are present. Recovery cannot resurrect deleted chats by design, but barrier/queued-write and timeout gaps remain (M1/M7).

### 5.9 Media and F3-F

**Result: boundary holds.** F3-C core can proceed independently; F3-C full and V20 remain blocked. References are not durability. No F3-F contract/evidence exists.

### 5.10 Multi-character isolation and V24

**Result: boundary holds, verification does not.** R3 states no cast/prompt/memory changes. V24 is the right regression shape. Baseline/evidence is absent (N15); concurrency and removal races remain (M2/M8).

### 5.11 Gate V1–V24

**Result: mapped, not audited.** All 24 rows are named. Row-level pass/fail criteria and evidence are incomplete (M9). No row passes; V20 is blocked; V16/V17/V24 lack native/evidence.

---

## 6. New defects and divergent-compliance cases

- **M1:** queued persistence continuation after barrier → late commit.
- **M2:** two character attempts → same base, lost update or duplicate order.
- **M3:** two model identities for different effective artifacts.
- **M4:** same request, different fingerprint due to serialization.
- **M5:** archive missing/corrupt after v2 commit; rollback impossible.
- **M6:** candidate fingerprint implementations differ.
- **M7:** barrier timeout has no recovery.
- **M8:** character removal during generation.
- **M9:** Gate row checked by label, not evidence.
- **M10:** preference/selection race overwrites explicit preference.

Additional divergence:

- Whether an in-flight commit that began before the barrier is adopted or aborted.
- Whether `requestFingerprint` includes prompt text (privacy) or only template identity.
- Whether predecessor archive retention is permanent or deletion-scoped.
- Whether concurrent character responses commit as separate turns or one merged turn.
- Whether `SetContinuationPreference` can save unavailable B while A is active and what the next send does.

---

## 7. Owner decisions still required

1. **Barrier authority:** which component owns the fence epoch, and which write paths must check it.
2. **Barrier timeout:** does it leave a durable closed barrier, block until manual resolution, or have a safe recovery?
3. **Concurrent character commits:** serialize at commit, merge, or deterministic order; define stale-base outcome.
4. **Character removal fence:** version/epoch or explicit invalidation.
5. **Model identity components:** tokenizer/vocabulary, adapters/LoRA, auxiliary models, backend/executable version.
6. **Request fingerprint canonicalization:** exact canonical form and schema transitions.
7. **Predecessor archive lifecycle:** integrity verification, rollback after v2, retention, wipe, missing/corrupt recovery.
8. **Candidate fingerprint:** exact `F3C-CAND-1` domain.
9. **Preference races:** ordering of selection/preference/transition operations.
10. **Gate evidence:** V1–V24 pass/fail criteria and prerequisites, including V16/V17/V24.

These are not optional implementation details; each changes the contract’s observable behavior.

---

## 8. Complete Gate V1–V24 audit

No row is PASS. This is a contract review, not a test run.

| Row | Disposition | Reason |
|---|---|---|
| V1 | UNCHANGED / NOT TESTED | Canonicalization/hash chain; R3 does not alter; no tests. |
| V2 | UNCHANGED / NOT TESTED | Receipt inventory; no tests. |
| V3 | UNCHANGED / NOT TESTED | Commit fork/orphan; no tests. |
| V4 | UNRESOLVED | M5/M6 target/migration bytes. |
| V5 | UNRESOLVED | M1 queued writes, M2 concurrent commits. |
| V6 | UNRESOLVED | M1/M2/M7/M10 lock/barrier ordering. |
| V7 | UNRESOLVED | M5 predecessor archive and Save As. |
| V8 | UNRESOLVED | M1 barrier, M8 character removal. |
| V9 | UNRESOLVED | Creation v2; no evidence. |
| V10 | UNCHANGED / NOT TESTED | Memory migration; no tests. |
| V11 | UNRESOLVED | M5 archive/unknown-field lifecycle. |
| V12 | UNRESOLVED | M1/M7 wipe barrier/staging. |
| V13 | UNRESOLVED | M5/M7/M8 fault boundaries. |
| V14 | UNMAPPED | Filesystem recovery; no row-level disposition. |
| V15 | UNRESOLVED | M1/M2/M8/M10 races. |
| V16 | UNVERIFIED | No native evidence. |
| V17 | UNRESOLVED | Mimo challenge requirement; no run. |
| V18 | UNRESOLVED | M3/M4 identity/fingerprint. |
| V19 | UNRESOLVED | Seed semantics defined, no evidence. |
| V20 | BLOCKED | F3-F ownership evidence absent. |
| V21 | UNRESOLVED | M10 preference, M3 identity, no runtime. |
| V22 | UNRESOLVED | M1/M2/M6/M7/M8 recovery races. |
| V23 | UNRESOLVED | M5 archive/rollback. |
| V24 | UNRESOLVED | M2/M8 concurrency/removal; no baseline. |

V20 is blocked, not failed. Gate v2 itself is unchanged and unfrozen.

---

## 9. Evidence and integrity

- Frozen authority SHA-256 verified: `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`.
- R3 SHA-256: `5013eb6fad7e0ec6f067d0fdbd350fd7f7c68ff9e263a7967b37f39e6daaf42c`.
- R3 line count: 520.
- R3 fenced-block count: 22; balanced.
- No source, test, build, migration, model load, destructive test, or application-data change.
- No R3, frozen-contract, prior-review, or source edit.

---

## 10. Status

**Question:** Does R3 define a safe, internally consistent, testable contract suitable for formal freeze?  
**Answer: No. It is ready for another revision and independent review, not freeze.**

**R3:** not frozen, not accepted.  
**Frozen contract:** unchanged, hash verified.  
**F3-C phase 0:** **NOT CLOSED**.  
**P1a:** **UNAUTHORIZED**, downstream of P1c.  
**F3-C acceptance:** **NOT GRANTED**.  
**F3-D:** **BLOCKED**.  
**F3-B:** closed.  
**A9/A10/A12/A13:** closed; not reopened.  
**V20:** **BLOCKED** pending F3-F evidence.

No implementation authorization exists.
