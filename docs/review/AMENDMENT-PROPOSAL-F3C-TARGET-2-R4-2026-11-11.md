# AMENDMENT PROPOSAL — F3C-TARGET-2, Revision 4

**Status:** **PROPOSAL ONLY — NOT FROZEN AUTHORITY, NOT ACCEPTANCE, NOT IMPLEMENTATION.**  
**Date:** 2026-11-11  
**Supersedes for review:** R3 (`docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R3-2026-11-11.md`). R1/R2/R3 remain historical.  
**R3 SHA-256:** `5013eb6fad7e0ec6f067d0fdbd350fd7f7c68ff9e263a7967b37f39e6daaf42c`  
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`  
**Independent review addressed:** `docs/review/REVIEW-INDEPENDENT-F3C-TARGET-2-R3-2026-11-11.md`  
**Scope:** contract-resolution proposal only. No implementation, source changes, tests, builds, model loads, migrations, destructive tests, authority changes, or freeze.

R4 resolves the ten findings in the independent R3 review, the fifteen N findings, and the ten R findings. It does not treat a prose resolution as a passed Gate row. Its central change is a single durable-mutation authority that fences every writer while preserving independent character generation.

---

## 1. Resolution matrix

### 1.1 R4 resolution of the R3-review findings

| ID | Finding | R4 resolution | Status |
|---|---|---|---|
| M1 | Barrier does not fence queued persistence continuations | Every durable mutation, including staging, autosave, append, replacement, commit record, receipt, migration, intent, tombstone and wipe, must hold a `MutationLease`; the barrier drains leases and commit gates, not just producers. | Resolved in contract |
| M2 | Per-character commits collide on the same base | Execution remains concurrent; all results pass through one per-chat commit coordinator with a durable response plan, stable participant order, exact history-parent checks, duplicate suppression and prefix commits. | Resolved in contract |
| M3 | Model identity omits inference-relevant artifacts | `F3C-MODEL-1` includes all inference-affecting artifacts; execution environment is a separate `F3C-EXEC-1` record. | Resolved in contract |
| M4 | Request fingerprint canonicalization undefined | `F3C-REQ-1` with `F3C-CANON-2`, explicit defaults, type/default/omission rules, seed semantics, prompt digest and schema versions. | Resolved in contract |
| M5 | Predecessor archive integrity/rollback/post-commit incomplete | `F3C-ARCH-1` archive manifest, exact-byte binding, verification, recovery at every boundary, explicit no-in-place-rollback policy and evidence-gap rules. | Resolved in contract |
| M6 | `F3C-CAND-1` not fully specified | `F3C-CAND-2` binds exact target bytes, manifest, commit parent, batch slot, attempt and request identity under `F3C-CANON-2`. | Resolved in contract |
| M7 | Barrier timeout has no recovery state | Durable `blocked_timeout` barrier state; no release, no destructive work; restart recovery and explicit retry/resume rules. | Resolved in contract |
| M8 | Character removal has no attempt fence | `characterVersion`/`castEpoch` in execution and batch slots; removal terminalizes the slot and invalidates its attempt. | Resolved in contract |
| M9 | Gate mapping lacks pass/fail/evidence | Every V1–V24 row below defines claim, contract, evidence, oracle, positive, negative, fault/native/Mimo, dependencies, PASS/FAIL/UNRESOLVED/BLOCKING and evidence identity. | Resolved as a verification specification; no row passed |
| M10 | Preference races with transition/commit | `PreferenceAuthority` with monotonic `preferenceSeq`; explicit preference operations; load completion never writes saved preference; terminal readiness API. | Resolved in contract |

### 1.2 Disposition of the fifteen R2-review findings

| ID | Prior finding | R4 disposition | Residual evidence |
|---|---|---|---|
| N1 | Absent chat treated as tombstoned | Resolved: absent without tombstone is `ErrChatStateUnknown`; no cleanup. | V8/V22. |
| N2 | Commit existence not tied to attempt | Resolved: exact `F3C-CAND-2` and receipt/manifest identity. | V13/V22. |
| N3 | Delete/wipe do not quiesce stream | Resolved by lease barrier; producer termination alone is explicitly insufficient. | V6/V12/V15. |
| N4 | Migration installs record before targets | Resolved by one linearization record and recovery table. | V4/V13/V23. |
| N5 | Model identity deferred | Resolved by `F3C-MODEL-1`. | V18. |
| N6 | Parameter set incomplete | Resolved by complete effective-request schema and extension rule. | V18/V19. |
| N7 | Strict unknown fields brick chats | Resolved by exact predecessor archive; v2 remains strict. | V4/V11/V23. |
| N8 | Predecessor bytes not retained | Resolved by immutable `F3C-ARCH-1` archive. | V23. |
| N9 | Preference-update ambiguity | Resolved by explicit operations and `preferenceSeq`. | V21. |
| N10 | Gate V1–V24 only partly mapped | Resolved by full row specification. | All rows. |
| N11 | One attempt per chat changes multi-character | Resolved by execution concurrency plus serialized commit plan. | V15/V24. |
| N12 | Root mismatch cleans evidence | Resolved: blocks and preserves. | V8/V13/V22. |
| N13 | Candidate fingerprint undefined | Resolved by `F3C-CAND-2`. | V13/V22. |
| N14 | New-chat format unspecified | Resolved: all new chats, Save As and post-amendment writes are v2. | V4/V9/V23. |
| N15 | V24 has no baseline | Resolved as a required captured pre-amendment baseline; evidence still absent. | V24. |

### 1.3 Disposition of the ten R1-review findings

| ID | Prior finding | R4 disposition |
|---|---|---|
| R1 | Generation-intent journal outside deletion/wipe | Intents and all operation roots are under the root mutation authority and wipe inventory; deletion requires tombstone. |
| R2 | Intent lacks commit/switch policy | Every attempt binds user commit, assistant commit, batch, candidate, receipt and session; switch invalidates the execution lease. |
| R3 | Third journal lacks recovery order | One precedence order in §7.6; ordinary/cast/intent combinations block on ambiguity. |
| R4 | Hot-swap rollback not a transaction | Runtime transition state machine, terminal readiness, failure matrix and preference authority in §8. |
| R5 | Provenance not actual runtime configuration | Request identity is captured from the resolved effective request; model and execution provenance are separate. |
| R6 | V1 migration lacks unknown/predecessor policy | Exact archive and strict-v2/lossless-predecessor split. |
| R7 | Media scope ambiguous | Core/full split retained; V20 blocked. |
| R8 | Multi-character boundary not explicit | Non-goal and V24 baseline retained; batch commit is persistence-only. |
| R9 | V2 optional bytes incomplete | R4 fixes presence, ordering, null/omission and new fields in §10. |
| R10 | Non-streaming/character paths uncovered | All assistant-producing paths use one attempt protocol; `SendMessage`, stream-with-image and `GenerateCharacterResponse` are explicitly covered. |

No prior finding is treated as verified by this matrix.

---

## 2. Authority model

### 2.1 Terms and invariants

- **Durable mutation:** any operation that creates, changes, replaces, stages, renames, removes or publishes bytes whose loss or stale retention can affect a chat, a preference, an operation journal, a receipt, a tombstone, a migration or a wipe.
- **Mutation authority:** the root-owned, durable authority that issues and drains `MutationLease` records. It is not a generation scheduler and does not own prompts, character context or model inference.
- **Execution authority:** the per-attempt authority that permits a character or narrator generation to run. It does not grant write authority.
- **Commit authority:** the per-chat serialization point that turns a completed candidate into a durable history commit. It does not serialize generation.
- **Destructive authority:** the exclusive root/chat authority used by delete and wipe after the mutation barrier is sealed.

The following invariants are normative:

1. **A1 — Single permission source.** No durable write may occur without a valid `MutationLease` from the root authority.
2. **A2 — No self-issued permission.** Chat persistence, autosave, cast sync, model preference, recovery, migration and UI adapters cannot create their own write permission.
3. **A3 — Check at commit.** The authority check must be coupled to the durable commit decision. A check performed when work starts, or before a long I/O operation, is insufficient.
4. **A4 — Stale leases fail.** A lease cannot become valid again through retry, replay, delay or a new process.
5. **A5 — Execution is not serialized.** Character thinking may overlap. Only shared durable commit decisions are serialized.
6. **A6 — Destructive work is last.** Delete/wipe may begin only after the applicable mutation boundary is proven quiescent.
7. **A7 — Evidence is not success.** A missing, corrupt or mismatched proof is a conflict or evidence gap, never a commit.
8. **A8 — No resurrection.** No recovery path may recreate a deleted chat, a removed character, or a replaced root.

### 2.2 Root and chat authority records

The authority has one root record and one record per chat:

```text
f3c-authority/<F3C-PATH-AUTH-1 root hash>/root.json
f3c-authority/<F3C-PATH-AUTH-1 root hash>/chats/<chatId>.json
f3c-authority/<F3C-PATH-AUTH-1 root hash>/barriers/<barrierId>.json
f3c-authority/<F3C-PATH-AUTH-1 root hash>/preferences.json
```

These are operation evidence, not chat targets. They are included in the wipe inventory.

`chat authority` fields:

```text
version, rootId, rootIncarnation, chatId, state,
fenceEpoch, commitSequence, historyHead,
lastCommittedOperationId, activeLeaseCount, activeCommitGate,
castEpoch, barrierId, recoveryIncarnation
```

`state` is exactly one of:

```text
active
closing
draining
sealed
destructive
blocked_timeout
completed
```

`root authority` carries the same root identity, a root `fenceEpoch`, a `wipeBarrierId`, the set of chat authority IDs, and a `preferenceSeq`. A root-wide wipe closes every chat authority before staging any root.

### 2.3 Mutation lease

A lease is issued by the authority, not by the caller:

```text
leaseId, rootId, rootIncarnation, chatId,
scope: chat | preference | migration | recovery | destructive,
operationId, operationKind,
fenceEpoch, issuedAt, ownerProcessIncarnation,
baseRevision, baseHistoryDigest,
status: active | committing | committed | aborted | orphaned
```

Issuance rules:

1. `BeginMutation` takes the root/chat authority lock for a short, bounded critical section.
2. If the chat is not `active`, the request fails with `ErrMutationFenceClosed`.
3. If the request carries an older `fenceEpoch` or an operation identity already associated with a different lease, it fails with `ErrMutationFenceStale`.
4. The authority records the lease before returning. The lease is durable before the caller performs its first durable write.
5. A lease is required for staging writes, not just final publication.
6. A lease has an owner process incarnation. A lease from a prior incarnation is `orphaned`, never adopted silently.

### 2.4 Commit protocol and linearization

Every mutation follows this logical sequence:

1. **Begin:** obtain the lease (§2.3); capture the expected base revision/history digest.
2. **Stage:** write candidate bytes to a lease-owned temporary area; all staging is covered by the lease.
3. **Verify:** read back; verify exact bytes, manifest, predecessor and operation identity.
4. **Commit decision:** acquire the chat authority lock, then the chat commit gate, following the fixed order in §2.7. Do not acquire the authority after the gate.
5. **Authority re-check:** while holding the authority lock and the commit gate, verify the lease is `active`, the fence epoch is current, the chat is `active`, the barrier is absent, the base/parent matches, and no conflicting operation is committed.
6. **Linearize:** install the durable commit record with a new `commitSequence` and `historyHeadAfter`. The atomic installation of this record is the **single linearization point**.
7. **Project:** install the receipt and publish the target; both are deterministic projections of the commit record. A crash after the record is recoverable by reconstruction; a crash before it is not a commit.
8. **Finish:** mark the lease `committed`, update the authority projection, release the commit gate and the authority lock. The authority lock and commit gate are held only for steps 5–8; no model inference, UI callback, stream wait or destructive filesystem operation occurs while they are held.

The commit record is authoritative for both target bytes and receipt identity. A receipt that appears without a matching commit record is not a commit. A target that appears without a matching commit record is uncommitted staging and must be reconciled.

The commit gate is per chat. It is not a global generation lock. It may delay another commit for the same chat; it must never wait for model inference.

### 2.5 Every durable path must be leased

| Durable path | Authority scope | Required rule |
|---|---|---|
| Autosave and final save | chat | `SaveChat`/`performSave`/`performFinalSave` obtain a lease and commit a v2 target set; no direct `saveJSON` bypass. |
| Chat append | chat | `AppendHistoryEntry` obtains a lease; assistant append is only legal after a generation candidate is committed. |
| History replacement/clear | chat | `RewriteHistory`/`ClearChatHistory` obtain a lease, verify base/history head, and commit a new record. |
| Full metadata/target save | chat | World, scenario, cast, lorebook, memory, meta and history writes use the same commit protocol. |
| User-turn commit | chat | A user turn becomes durable before an attempt starts; it returns `userCommitOperationId` and the response batch plan. |
| Assistant completion | chat | The candidate is staged; commit is separated from stream completion and uses the commit coordinator. |
| Cast synchronization | chat | Cast changes are a durable mutation; they advance `castEpoch` and invalidate affected character slots. |
| Character/world/scenario target writes | chat | No independent write path; all target writes are commits or staging under a lease. |
| Import | chat/root | Import stages into a new chat under a lease; it never writes into an existing chat without a commit. |
| Migration | chat | `F3C-ARCH-1` archive, staged targets and migration record are lease-bound. |
| Generation-intent updates | chat | Intent state changes are journal writes under a lease; terminalization cannot bypass the fence. |
| Commit records and receipts | chat | Written only by `CommitMutation` or its deterministic recovery projection. |
| Tombstones | chat/root | Written only after the destructive barrier is `sealed`. |
| Delete | chat | Barrier, tombstone, intent cleanup, archive cleanup and directory removal are one destructive authority sequence. |
| Wipe staging | root | Root barrier closes every chat, stages all operation roots and the gallery, then commits root reset. |
| Root-incarnation replacement | root | A new root has a new identity/incarnation; old-root leases are never valid. |
| Recovery/reconciliation | chat/root | Recovery obtains a recovery lease and follows §7.6; it cannot adopt or delete by convention. |
| Model preference | preference | `SetContinuationPreference` writes through the root preference authority; it is not a chat target write. |
| Active-chat pointer/config | root | The pointer is changed under a root lease after the chat commit/destructive outcome; it cannot be a bypass. |
| Gallery/F3-F bytes | F3-F | F3-C never writes media bytes; F3-F is outside this authority and V20 remains blocked. |

A path not listed above is still governed by A1. “Not listed” is not permission.

### 2.6 Commit identity and idempotency

The commit record contains:

```text
recordVersion, operationId, commitSequence, rootId, rootIncarnation,
chatId, batchId, participantIndex, attemptId, attemptOrdinal,
baseRevision, baseHistoryDigest, historyHeadAfter,
targetManifestDigest, candidateFingerprint, requestFingerprint,
receiptId, predecessorArchiveHash?, timestamp
```

The pair `(chatId, operationId)` is unique. A repeated commit request with the same operation ID and same `candidateFingerprint` returns the existing proof if and only if the record matches; a different fingerprint returns `ErrCommitOperationConflict`. A different operation ID is a different mutation.

### 2.7 Lock ordering and non-deadlock rules

The only permitted order is:

```text
root authority lock
  -> chat authority lock
    -> chat commit gate
      -> lease-owned staging I/O
```

The following are forbidden:

- holding the root authority while waiting for a producer or a commit gate;
- taking a chat authority while holding a commit gate and then taking another chat authority;
- waiting for generation while holding a lease;
- taking a lease from the commit gate without releasing the gate;
- any lock whose holder can wait for a model load, stream, UI callback or filesystem deletion.

The barrier sets `closing` under the chat authority, releases the authority lock, cancels producers, and then drains leases and the commit gate. It does not hold the authority lock while waiting. A root wipe follows the same pattern in a fixed chat-ID order.

### 2.8 Mutation barrier

A delete or wipe creates a durable `barrierId` and follows this exact sequence:

1. **Acquire destructive ownership.** Take the root/chat authority; verify root incarnation and chat identity. Reserve `barrierEpoch = current fenceEpoch + 1`; it is not published until seal.
2. **Close new leases.** Set state `closing`; future `BeginMutation` calls fail with `ErrMutationFenceClosed`. Existing leases remain tracked. Release the authority lock before waiting.
3. **Cancel producers.** Request cancellation of every generation attempt and every queued producer in scope. This is not termination proof.
4. **Drain leases.** Wait until every lease in scope is `committed`, `aborted` or `orphaned`, and every commit gate is empty. Do not hold the authority lock or a commit gate while waiting. A producer that has terminated but left an unleased continuation is still not drained; the continuation must attempt a lease and fail, or be represented by an active lease.
5. **Reconcile.** For every committed record, include it in the destructive outcome. For every lease without a commit record, abort and clean its staging. There is no partial adoption.
6. **Seal.** Under the authority, set `fenceEpoch = barrierEpoch`, state `sealed`, `sealSequence = commitSequence`, and persist the barrier. No lease can be issued after this point. Because all leases were drained in step 4, no old-epoch lease can commit after seal.
7. **Destructive step.** Write the tombstone and remove the chat, or stage the root for wipe. These operations are performed only by the destructive authority.
8. **Complete.** Record `completed`; release destructive ownership only after the destructive operation has a durable outcome.

`blocked_timeout` is a terminal barrier state for the current attempt, not a release. It preserves the barrier record, rejects new leases, and permits only diagnostic/read-only access. It never authorizes tombstone, deletion or staging.

### 2.9 Timeout, crash and recovery

- A timeout while draining is `ErrMutationBarrierTimeout`; state becomes `blocked_timeout`; active leases remain listed; no destructive step occurs.
- A user or operator may retry the same barrier ID after resolving the cause. A retry starts at step 4; it never skips draining.
- A process crash while `closing` or `draining` leaves the durable barrier. On restart, the authority marks leases whose `ownerProcessIncarnation` is not the current incarnation as `orphaned`. An orphan with a matching commit record is reconciled as committed; an orphan without one is aborted and its staging cleaned.
- A crash after `sealed` but before tombstone/staging resumes only the destructive step. No ordinary lease is valid.
- A crash after tombstone but before directory removal completes deletion; no chat recreation is permitted.
- If a producer cannot be proven terminated in the current process, the barrier remains `blocked_timeout`; no destructive authority is granted.
- `ErrGenerationBarrierTimeout` and `ErrMutationBarrierTimeout` are distinct: the former is an execution/producer failure, the latter is a durable write-boundary failure. Both fail closed.

---

## 3. Execution and commit separation

### 3.1 Execution scopes

Execution authority is:

```text
ExecutionScope = (chatId, characterId, characterVersion, stateSessionId)
```

- `characterId = ""` denotes the narrator/generic path.
- `characterVersion` is a stable version incremented on removal, replacement or identity invalidation. A recreated character with the same display name has a new ID or a new version; it never inherits a stale attempt.
- At most one active attempt exists for a given execution scope.
- Different characters may execute concurrently. A narrator and character may execute concurrently only if the response plan declares both slots and the existing product permits it.

The execution authority only grants permission to generate. It does not grant permission to write. Every write from an attempt must still pass through a `MutationLease`.

### 3.2 Response batch and participant order

Before the first character-specific attempt for a user turn, the orchestrator submits a `ResponsePlan` with the user-turn commit:

```text
planVersion, batchId, userCommitOperationId,
historyBaseDigest, participantOrder: [
  { participantIndex, characterId, characterVersion }
]
```

Rules:

- `participantOrder` is the explicit, durable ordering. It must not be inferred from completion time, scheduler order, map iteration or UI array mutation.
- The plan may contain a narrator slot and character slots. A character not in the plan cannot create an assistant turn for that batch.
- If the caller supplies no plan, the user-turn commit creates a singleton plan for the invoked path. A later attempt for the same user turn that is not in the plan is rejected with `ErrResponsePlanMissing`; it is not silently appended.
- The plan is immutable for the batch. Character removal terminalizes its slot; it does not reorder remaining slots.
- The plan is committed by the user-turn mutation and is visible to recovery before any candidate can commit.

### 3.3 Candidate submission

A completed attempt produces an immutable candidate:

```text
candidateId, batchId, participantIndex, characterId, characterVersion,
attemptId, attemptOrdinal, userCommitOperationId,
baseHistoryDigest, historyHeadAfterCandidate,
candidateContentDigest, targetManifestDigest,
assistantMessageId, candidateFingerprint, requestFingerprint,
status: staged | eligible | committed | rejected | conflict
```

A candidate is staged through a lease. It is **not** appended to history, and no `chat:complete` event is emitted, until the commit coordinator commits it.

### 3.4 Serialized commit, concurrent thought

The per-chat commit coordinator applies these rules:

1. **One commit at a time per chat.** It holds the chat commit gate, not the execution authority.
2. **Prefix order.** A candidate for `participantIndex = n` may commit only when every lower `participantIndex` in the batch is `committed` or `terminal_no_turn` (`failed`, `aborted`, `removed`, `skipped`).
3. **Exact parent.** The expected history parent is the batch's `historyBaseDigest` followed by the committed candidate hashes of lower participant indices. It must equal the current `historyHead`. If it does not, return `ErrGenerationCommitConflict`; never overwrite.
4. **Duplicate suppression.** `(batchId, participantIndex, attemptOrdinal)` is unique. A retry is a new `attemptOrdinal`; a second commit for an already committed slot is rejected unless it is the exact same candidate/operation, in which case the existing proof is returned.
5. **Stable identity.** `assistantMessageId` and `candidateFingerprint` are allocated at attempt start; the commit does not regenerate them.
6. **No fabricated order.** Completion time never determines history order. A later slot waits for earlier terminal slots. A batch deadline may mark an unstarted or running lower slot `abandoned`; that is an explicit terminal state, not a fabricated assistant turn.
7. **No lost update.** The commit record carries `baseHistoryDigest` and `historyHeadAfter`; a CAS mismatch is a conflict, not a merge.
8. **Removal.** A removed character's slot becomes `removed`; its active attempt is invalidated; its candidate cannot commit. A later slot proceeds in plan order.
9. **Restart.** The plan, candidates, terminal states and commit sequence survive restart. Recovery resumes at the first nonterminal slot; it does not reorder committed turns.
10. **Chat deletion/wipe.** The destructive barrier invalidates all execution scopes and commit candidates in the chat; no character scope can bypass it.

### 3.5 What is preserved

The commit coordinator writes an assistant message with the existing `characterId` ownership and the existing sender/character semantics. It does not:

- build, rewrite or reorder character prompts;
- change what history or lore a character sees;
- change memory isolation or context flushing;
- change cast synchronization rules;
- combine two character responses into one message;
- alter character-specific context or memory visibility.

Execution context is read-only input to generation. Persistence only records the result.

---

## 4. Model identity and execution provenance

### 4.1 `F3C-MODEL-1` artifact identity

`modelId` is:

```text
"mdl_" + lowercase_hex(SHA-256(F3C-CANON-2(model_identity_envelope)))
```

The envelope contains:

```text
schemaVersion: "F3C-MODEL-1"
components: [
  { role, artifactHash, byteSize, format, required }
]
```

`role` is one of:

```text
main-model
tokenizer
vocabulary
projector
vision-encoder
adapter
auxiliary-model
embedding-model
scoring-model
quantization-artifact
conversion-metadata
other-inference-artifact
```

Rules:

- Every component hash is SHA-256 over exact bytes. No path, basename, display label, registry name or mtime is identity.
- Multi-component lists are sorted by `(role, artifactHash, byteSize)`; duplicate roles are permitted only when each has a distinct role qualifier and artifact hash.
- An artifact embedded in the main GGUF (for example an embedded tokenizer) is represented as a component with `artifactHash = main-model hash`, `format = "embedded:<component>"`, and `role` = the embedded role.
- Adapters/LoRAs, auxiliary models, tokenizers/vocabularies, projectors/vision encoders, embedding/scoring models, quantization artifacts and conversion metadata are included when they can affect the inference path.
- A required component that cannot be hashed means no generation. No placeholder ID is permitted.
- Optional absent components are omitted from the component list; their absence is part of the identity only when the engine treats absence as a distinct mode. `required: false` and omitted are different from `required: true`.
- File paths and installed-model labels may be stored as **hints** in the target, but never used to compute `modelId`.

### 4.2 `F3C-EXEC-1` execution provenance

Execution provenance is separate from model identity:

```text
engineName, engineVersion, engineBuildId?,
runtimeName, runtimeVersion?,
executableHash?, executableVersion?,
acceleration: { backend, device, offloadMode, offloadLayers? },
engineConfigDigest,
adapterLoadingMode,
templateFingerprint,
provenanceStatus: complete | partial | unavailable
```

Rules:

- `engineConfigDigest` is `F3C-REQ-1` over effective engine/runtime controls that can change output.
- `executableHash` is required when the backend binary is launchable and its bytes are available; otherwise `provenanceStatus = unavailable` for that field, not a fabricated value.
- Acceleration/offload configuration is recorded even when it is machine-specific; it is execution provenance, not model identity.
- Historical turns retain the original `modelId`, `F3C-EXEC-1` snapshot and request fingerprint. A later installed model with the same path does not rewrite them.
- Artifact status is a separate current observation:

```text
available | missing | relocated | replaced | corrupt | unreadable
```

`replaced` means the bytes at the historical path no longer match the historical `modelId`. The historical turn remains truthful; the UI reports the mismatch.

### 4.3 Model removal/replacement

- Uninstalling a model does not alter historical `modelId` or provenance.
- Replacing bytes at the same path creates a different `modelId`; the old ID remains historical.
- Corruption or unavailability is reported as `artifactStatus` on the current artifact, not substituted.
- If a tokenizer/adapter/component is missing, a new generation that needs it fails with `ErrModelArtifactIncomplete`; a historical turn remains readable.
- Backend version changes create a new `F3C-EXEC-1` record; they do not rewrite old turns.
- Incomplete historical provenance is retained as `partial`/`unavailable`; no backfill invents values.

---

## 5. Request identity and canonicalization

### 5.1 `F3C-REQ-1` effective request

`requestFingerprint` is:

```text
"req_" + lowercase_hex(SHA-256(F3C-CANON-2(request_envelope)))
```

The envelope contains:

```text
schemaVersion, modelId, executionFingerprint, defaultSetId,
promptDigest, templateFingerprint, adapterConfigDigest,
parameters: { each key below }, seedMode, effectiveSeed,
engineExtensionParameters
```

`promptDigest` is SHA-256 over canonical prompt messages. Raw prompt text is not stored in the fingerprint and is never logged. The digest identifies the actual prompt used, including history and character context, without exposing it.

### 5.2 Parameter schema

`F3C-REQ-1` includes at least the following. Values are the **effective** values after defaults, clamping, template injection and backend normalization; `origin` records whether the value was explicit, inherited, engine-defaulted or unavailable.

| Key | Type | Omission/default rule |
|---|---|---|
| `temperature` | finite number | Required; default `0.7` from `F3C-DEFAULTS-1`, then family and explicit overrides. |
| `topPEnabled` | boolean | Required; default `true`. |
| `topP` | finite number | Required; default `0.9`; if `topPEnabled = false`, effective `topP = 1.0`. |
| `topK` | integer >= 0 | Required; default `40`. |
| `minP` | finite number | Required; default `0.05`. |
| `typicalPEnabled` | boolean | Required; default `true`. |
| `typicalP` | finite number | Required; default `1.0`. |
| `tfs` | finite number | Required; default `1.0`. |
| `repeatPenalty` | finite number | Required; default `1.1`. |
| `repeatLastN` | integer >= 0 | Required; default `64`. |
| `frequencyPenaltyEnabled` | boolean | Required; default `true`. |
| `frequencyPenalty` | finite number | Required; default `0.0`. |
| `presencePenaltyEnabled` | boolean | Required; default `true`. |
| `presencePenalty` | finite number | Required; default `0.0`. |
| `mirostatEnabled` | boolean | Required; default `true`. |
| `mirostat` | integer enum | Required; default `0` (disabled). |
| `mirostatTau` | finite number | Required; default `5.0`. |
| `mirostatEta` | finite number | Required; default `0.1`. |
| `dynamicTempRangeEnabled` | boolean | Required; default `true`. |
| `dynamicTempRange` | finite number | Required; default `0.0`. |
| `dynamicTempExponent` | finite number | Required; default `1.0`. |
| `dryMultiplier` | finite number | Required; default `0.0`. |
| `dryAllowedLength` | integer >= 0 | Required; default `2`. |
| `dryBase` | finite number | Required; default `1.0`. |
| `smoothingFactor` | finite number | Required; default `0.0`. |
| `smoothingCurve` | finite number | Required; default `1.0`. |
| `topAEnabled` | boolean | Required; default `true`. |
| `topA` | finite number | Required; default `0.0`. |
| `epsilonCutoff` | finite number | Required; default `0.0`. |
| `etaCutoff` | finite number | Required; default `0.0`. |
| `encoderRepeatPenalty` | finite number | Required; default `1.0`. |
| `noRepeatNGramSize` | integer >= 0 | Required; default `0`. |
| `maxTokens` | integer >= 1 | Required; default `nPredict = 512`, plus any effective thinking budget. |
| `nKeep` | integer >= 0 | Required; default `0`. |
| `contextTokens` | integer > 0 | Required; effective runtime context, not a UI default. |
| `stopSequences` | ordered array of strings | Always present; `[]` means none; default is the template/engine stop set. |
| `samplerOrder` | ordered array of enum strings | Always present; order is significant; `[]` means engine order. |
| `grammar` | string or null | Required key; `null` means none. |
| `jsonSchema` | string or null | Required key; `null` means none. |
| `logitBias` | map<string, number> | Sorted by key; `{}` means none. |
| `ignoreEOS` | boolean | Required; default `false`. |
| `seedMode` | `random|fixed` | Always present; default `random`. |
| `seed` | integer or null | `fixed` requires integer; `random` records `null` here; the legacy `-1` sentinel is never a literal fixed seed. |
| `effectiveSeed` | integer or null | Actual seed used; `null` only when the engine cannot report it, with `seedStatus: unknown`. |
| `promptTemplateId` | string | Required; identifies the template family. |
| `templateFingerprint` | hash | Required; exact template bytes. |
| `adapterConfigDigest` | hash or null | Required; `null` when no adapters. |
| `engineExtensionParameters` | canonical map | Required, possibly empty; every engine control not named above must be included, not dropped. |

`origin` is one of `explicit`, `inherited`, `engine-default`, `normalized`, `unavailable`.

### 5.3 `F3C-DEFAULTS-1` default resolution

The base defaults are the current `llama.NewInferenceOptions()` values above. Before hashing, the resolver applies, in order:

1. the base defaults;
2. the selected family's published overrides (`familySamplingDefaults`);
3. explicit/inherited user settings;
4. runtime normalization (`topPEnabled = false` -> `topP = 1.0`; thinking budget added to `nPredict`);
5. seed resolution (`-1` -> random; capture the effective seed when the engine reports it).

`defaultSetId` is:

```text
sha256(F3C-CANON-2({
  schemaVersion: "F3C-DEFAULTS-1",
  engineName, engineVersion, family,
  baseDefaults, familyOverrides
}))
```

The concrete numeric table above is the base table. A family override changes the effective value only when selected and is recorded in `defaultSetId`; no implementation may silently substitute a different family default. If the engine cannot report a default or effective seed, the field is `unavailable` or `null` with `seedStatus: unknown`; it is not guessed. Changing the base/family defaults requires a new `schemaVersion` and new golden vectors.

### 5.4 `F3C-CANON-2`

`F3C-CANON-2` uses RFC 8785 JSON Canonicalization Scheme, with these additions:

- UTF-8, no BOM; Unicode is not normalized.
- Object keys are sorted lexicographically by Unicode code point; duplicate keys are invalid.
- Arrays preserve order; array order is semantically significant.
- Numbers are finite IEEE-754 binary64 for numeric parameters; no `NaN`, `Infinity` or `-Infinity`; `-0` canonicalizes to `0`.
- Integer-valued fields must be encoded as integers; a fractional value in an integer field is schema-invalid.
- Strings are JSON strings with RFC 8785 escaping.
- `null` is a value, never an omitted key. Every key listed by the schema is present in the canonical envelope.
- Omitted input values are resolved to the effective value before canonicalization. An explicit default and an inherited default with the same effective value produce the same fingerprint. An explicit zero and a nonzero default do not.
- Canonical envelopes are versioned. Old fingerprints remain interpretable under their original `schemaVersion`; they are never recomputed under a new schema.
- `schemaVersion` changes only when semantics change; adding a new engine extension key requires a new schema version and a migration table.

### 5.5 Examples

Semantically identical:

- `{temperature: 0.7, topP: 0.9}` and `{topP: 0.9, temperature: 0.7}` canonicalize to the same key order.
- Omitted `repeatPenalty` with effective engine default `1.1` and explicit `repeatPenalty: 1.1` with the same origin-independent effective value produce the same fingerprint.
- `stopSequences: []` and omitted stop list resolved to `[]` produce the same fingerprint.
- Two prompts differing only in raw whitespace outside semantic template boundaries but with the same canonical message array produce the same `promptDigest`.

Materially different:

- `topK: 40` versus `topK: 41`.
- `temperature: 0` versus `temperature: 0.0001`.
- `seedMode: random` with effective seeds `1` and `2`.
- `grammar: null` versus a grammar envelope.
- Two template byte sequences with different `templateFingerprint`.
- Two adapter configurations with different `adapterConfigDigest`.

### 5.6 `F3C-CAND-2`

`candidateFingerprint` is:

```text
"cand_" + lowercase_hex(SHA-256(F3C-CANON-2(candidate_envelope)))
```

The envelope contains:

```text
schemaVersion: "F3C-CAND-2"
rootId, rootIncarnation, chatId, batchId, participantIndex,
userCommitOperationId, userMessageId, attemptId, attemptOrdinal,
characterId, characterVersion, baseRevision, baseHistoryDigest,
assistantMessageId, assistantContentDigest,
targetManifestDigest, requestFingerprint, modelId,
executionFingerprint, candidateCreatedAt
```

`assistantContentDigest` is the exact canonical assistant turn bytes; `targetManifestDigest` binds the exact target set. The fingerprint excludes secrets and raw prompt text. The same candidate must produce the same fingerprint on every implementation using `F3C-CANON-2`.

---

## 6. Predecessor archive and V1-to-V2 migration

### 6.1 New format

All new chats, Save As, imports and post-amendment writes use `F3C-TARGET-2`/`f3c-target-manifest-2`/`tgt2-`. v1 is legacy-only. There is no new-v1 creation path.

### 6.2 `F3C-ARCH-1` archive

Before the first v2 write, the exact v1 target bytes and evidence are copied into an immutable archive:

```text
f3c-predecessors/<F3C-PATH-PRED-1 root hash>/<chatId>/<v1Revision>/
```

The archive manifest is:

```text
schemaVersion: "F3C-ARCH-1"
rootId, rootIncarnation, chatId, predecessorRevision,
predecessorManifestDigest, predecessorReceiptId,
files: [{ path, sha256, byteSize }],
v1RecordDigest, v1ReceiptDigest,
archiveHash
```

`archiveHash` is `F3C-CANON-2` over the manifest with `archiveHash` omitted. The archive is immutable and write-once. `predecessorArchiveHash` is required in every migration receipt/record and equals `archiveHash`; it is not optional on a v1→v2 migration.

Binding rules:

- Each archived file hash must equal the corresponding v1 target hash in `predecessorManifestDigest`.
- The v1 receipt and record digests must match the archived evidence.
- The v2 receipt/record must include the archive hash and the v1 predecessor revision.
- An archive that is present but unverified is not evidence.
- A filename or path alone is never proof; the manifest hash is.

### 6.3 Unknown fields

- Exact v1 target bytes, including unknown fields, are copied into the archive.
- v2 targets remain strict: an unknown v2 key is `ErrTargetMalformed`.
- A v1 chat with unknown fields may upgrade because the unknown bytes are preserved and provable in the archive.
- The unknown fields are not copied into arbitrary v2 keys. They remain in the immutable predecessor archive and can be exported for repair.
- If the archive cannot be written or verified, no v2 write occurs; v1 remains current.

### 6.4 Migration transaction

1. Acquire a migration lease; verify v1 current, receipt, chain, manifest and root identity.
2. Create the archive create-only; copy exact v1 bytes and evidence.
3. Read back every archived file; verify hashes, sizes and complete manifest; compute `archiveHash`.
4. Stage v2 targets under the same lease; read back and verify.
5. Enter the chat commit gate; re-check fence, base, archive and staged manifest.
6. Install staged v2 targets atomically; verify installed bytes.
7. Install the **migration commit record** containing v1 predecessor, archive hash, v2 manifest and history head.
8. Install the receipt as a projection of the record.
9. Publish v2 and mark the lease committed.
10. Retain the archive until deletion/wipe or a separately approved retention policy.

No mixed v1/v2 manifest is valid. The record is the linearization point. A receipt without the matching record is not a migration.

### 6.5 Recovery at every boundary

| Boundary | Required recovery |
|---|---|
| Before archive creation | v1 current; remove lease staging; no v2. |
| Partial archive | Delete incomplete archive staging; v1 current; no v2. |
| Archive complete, unverified | Verify; if failed, keep v1 and report `ErrPredecessorArchiveFailed`; no v2. |
| Archive verified, before v2 staging | v1 current; archive retained; retry allowed under a new lease. |
| Partial v2 staging | Delete staged v2 only; v1 current; archive retained. |
| All v2 targets installed, before record | Restore v1 from verified archive; verify; if restore fails, block with evidence gap. |
| Record installed, before receipt | Reconstruct receipt from record; verify v2 targets; no duplicate record. |
| Receipt installed, before publish | Publish/verify v2; current becomes v2. |
| Published | Normal current/commit verification; verify archive binding. |
| Repeated recovery | Idempotent: same record/receipt/targets, no duplicate migration, no rehash of v1. |
| Archive missing/corrupt after v2 commit | Current v2 may be read, but migration evidence is `incomplete`; block rollback/export claims and V23; do not rehash or invent predecessor bytes. |
| Wipe/root-incarnation change | Root mismatch blocks; wipe precedence; no cross-root restore. |

### 6.6 Rollback policy

R4 explicitly chooses **no automatic in-place rollback after a v2 commit**. This is a deliberate contract decision:

- The predecessor archive is immutable evidence and supports forensic export/reconstruction, but restoring exact v1 bytes in place would create a v1 current state after the amendment and conflict with the new-chat-v2 invariant.
- A byte-exact restoration is a separate maintenance operation requiring its own barrier, owner authorization and a new migration plan. It is not an ordinary write and is not part of TARGET-2.
- The archive can produce a **read-only export** of the exact predecessor bytes. It must not claim that the live chat was rolled back.
- If the archive is missing/corrupt, the system reports `predecessorEvidence: incomplete`; rollback/export is unavailable, and no historical bytes are fabricated.
- `Save As` from a v1 source archives first, then creates v2; it never copies a v1 receipt as v2 proof.

This policy is honest: it retains bytes but does not promise an in-place rollback that the v2 invariant cannot safely provide.

### 6.7 Retention, deletion and wipe

- Archives are immutable; ordinary cleanup never deletes them.
- Chat deletion removes archive files only after the tombstone is authoritative and the destructive barrier is sealed.
- Wipe stages `f3c-predecessors/` with `chats/`, `f3c-commits/`, `f3c-receipts/`, `f3c-tombstones/`, `f3c-creation-intents/`, `f3c-generation-intents/`, `f3c-memory-migrations/`, `f3c-authority/`, the cast journal if present, and the gallery root.
- A root-incarnation change invalidates every old-root archive reference; no old-root lease or restore is valid.

---

## 7. Generation attempts and recovery

### 7.1 Every assistant-producing path

All of these paths use the same attempt, candidate and commit protocol:

- `SendMessage` synchronous generation;
- `SendMessageStreamWithImage` and the text streaming path;
- `GenerateCharacterResponse`;
- any future assistant-producing import, replay or test seam.

A path that produces assistant text without a durable user-turn commit, attempt record, candidate and commit record is non-compliant.

### 7.2 User turn

A user turn is durable before generation:

```text
userCommitOperationId, userMessageId, contentDigest,
userTurnSequence, batchId, responsePlan, baseRevision
```

The user turn is written through the mutation authority. `recordUserTurn` may remain an in-memory convenience, but it cannot be the durable user commit. If the append fails, no generation attempt starts.

### 7.3 Attempt record

```text
attemptId, attemptOrdinal, operationId, operationType,
rootId, rootIncarnation, chatId, executionScope,
characterId, characterVersion, stateSessionId,
batchId, participantIndex, userCommitOperationId, userMessageId,
baseRevision, baseHistoryDigest,
assistantMessageId, assistantCandidateFingerprint,
modelId, executionFingerprint, requestFingerprint,
status, terminalReason, partialDigest?, startedAt
```

Statuses:

```text
started -> streaming -> completion_pending -> committed
       \-> interrupted | aborted | failed | stale | removed | conflict
```

Rules:

- No attempt starts without the user commit.
- A retry is a new `attemptId`/`attemptOrdinal`; it never reuses a committed operation ID.
- `completion_pending` is not a commit.
- A partial response gets an attempt record and a `partialDigest`; it does not get a committed assistant turn and cannot be promoted by recovery.
- A mid-stream model switch is forbidden for the active attempt; the attempt keeps the model identity captured at start.
- The stream-end event is `chat:stream_end`; `chat:complete` is emitted only after the candidate commit record exists.

### 7.4 Exact adoption proof

Recovery adopts a commit only when all of the following match:

- root ID and incarnation;
- chat ID, batch ID and participant index;
- execution scope, character ID/version, state session and attempt identity;
- user commit operation and base history digest;
- candidate fingerprint, assistant message ID and target manifest;
- commit record and receipt body;
- request fingerprint and model/execution identity;
- no conflicting commit for the same attempt or slot.

A mismatch is `ErrGenerationCommitConflict`; no adoption and no cleanup as committed.

### 7.5 Missing commit outcomes

| State | Outcome |
|---|---|
| No tombstone, chat present, base current | `interrupted`; no assistant turn; partial remains non-durable. |
| No tombstone, base stale | `stale`; no overwrite; conflict evidence retained. |
| No tombstone, chat absent | `ErrChatStateUnknown`; preserve evidence; no cleanup. |
| Tombstone present | Deletion authoritative; clean intents and staging only after barrier/seal. |
| Root mismatch | `ErrRootIncarnationChanged`; block; no cleanup. |
| Character removed | `removed`; no assistant turn; slot terminal. |
| Barrier/timeout | `blocked`; no destructive cleanup; retry/recovery only under authority. |

### 7.6 Recovery precedence

1. Root identity and root incarnation.
2. Wipe barrier and root destructive state.
3. Chat barrier state.
4. Tombstone inventory.
5. Ordinary save/append commit journal.
6. Migration/commit records.
7. Cast journal and cast epoch.
8. Response plan and generation attempts.
9. Candidate staging.
10. Terminal cleanup.

Rules:

- Root mismatch → block, never clean.
- Wipe intent/barrier → no chat recovery; finish or block the wipe.
- Tombstone present → no chat recreation; clean only after sealed barrier.
- Ordinary/cast/attempt nonterminal for the same chat → `ErrRecoveryAmbiguous`; no silent choice.
- A commit record is authority; a target/receipt without it is staging.
- Repeated recovery is idempotent and must not duplicate turns, candidates, receipts or migration records.
- Recovery itself must hold a recovery lease and obey the fence.
- Recovery may not revive a removed character; a removed slot is terminal.

---

## 8. Preference and model transition

### 8.1 Three distinct states

- **Saved continuation preference:** the chat's `continuationModel` in `meta.json`; durable, per chat.
- **Requested model:** the model the user asked to load; runtime only.
- **Active model:** the model currently healthy and permitted to generate.

The old global `howl-last-model-path` is not authoritative. It may be a UI hint only; it cannot be read as saved preference or used to publish active state.

### 8.2 `PreferenceAuthority`

The root preference authority serializes explicit preference writes:

```text
preferenceSeq, chatId, modelId, selectedAt,
operationId, status: saved | failed | superseded
```

Rules:

- `SelectContinuationModel(modelRef)` loads/activates; it does **not** write saved preference.
- `SetContinuationPreference(modelRef)` is the only operation that writes `continuationModel`; it may save an unavailable model, but it never claims active.
- A preference write is a durable mutation under a lease and carries a strictly increasing `preferenceSeq`.
- An older asynchronous load or preference operation cannot overwrite a newer explicit preference. Completion checks the latest `preferenceSeq`; a stale completion is `superseded`.
- Two rapid preference changes are serialized; the last accepted sequence wins.
- A failed preference write reports `saved: false`; the previous saved value remains.
- If active and saved differ, sending uses the active model. If no active model exists, the send returns `ErrModelUnavailable`; there is no automatic fallback.
- A chat deletion/wipe invalidates pending preference operations for that chat.

### 8.3 Transition state machine

Runtime-only transition record:

```text
transitionId, chatId, prior: ModelSnapshot?, requested: ModelSnapshot,
phase: preflight | stopping | starting | ready | active | restoring | failed | unavailable | cancelled,
startedAt, failureStage?
```

`ModelSnapshot` includes model identity, projector/adapter identity, spawn settings, admission evidence, execution provenance and resource estimate.

Order:

1. preflight and admission before stopping the old model;
2. stop old; verify termination;
3. spawn new; verify process identity;
4. health and warm-up;
5. publish `ready` only after the backend is genuinely ready;
6. mark active;
7. on failure, kill partial process, then attempt one restore of the prior snapshot;
8. if restore fails, `unavailable`; no fallback model, no false active.

A model is not active merely because `LoadModel` returned, the process exists, the HTTP port is open, or `model:loaded` was emitted. Readiness requires a successful generation-capable health check and warm-up.

### 8.4 Readiness API

`LoadModel(request)` returns only after a terminal result:

```text
ModelTransitionResult {
  transitionId,
  state: ready | failed | unavailable | cancelled,
  activeModelId?,
  failureStage?,
  message
}
```

`LoadModel` must not return `nil` for `starting` or `loading`. Progress is published through `model:state` events; those events are not success. `ready` is the only success state. A crash during transition leaves no durable active model; startup may attempt the saved preference, but only a terminal `ready` publishes it.

### 8.5 Failure matrix

| Condition | Required result |
|---|---|
| Preflight/admission fails | Old model untouched; no stop; no preference write; `failed`. |
| Stop fails | Do not start replacement; health-check old; report `failed`; no new active. |
| Spawn fails | Kill partial process; one restore attempt; otherwise `unavailable`. |
| Health/warm-up fails | Kill process; one restore; otherwise `unavailable`. |
| Restore fails | `unavailable`; saved preference retained; no substitution. |
| Newer preference arrives | Older completion `superseded`; it never writes saved. |
| Crash mid-transition | No durable active; startup reports/retries safely; no false ready. |
| Chat deleted mid-transition | Preference operation fails; no active publication for deleted chat. |

---

## 9. Media and F3-F boundary

- F3-C core text/context/provenance work may proceed without F3-F.
- F3-C full acceptance and V20 remain **BLOCKED** until F3-F proves byte ownership, integrity, recovery, deletion and backup.
- F3-C stores references only:

```text
mediaId, role, provider, mimeType, byteSize, contentHash
```

- `f3f-media` references require `contentHash`; `legacy-indexeddb` may omit it.
- A reference is not proof of durable bytes. Missing, corrupt, unavailable or externally owned media is reported, never substituted.
- F3-C never writes or deletes F3-F bytes; F3-F owns lifecycle.
- No base64, local path or thumbnail is a target field.
- `F3-C-CORE` excludes attachment durability; `F3-C-FULL` includes it.

---

## 10. V2 serialization additions

R4 retains the R3 v2 target set and serialization rules, with these required additions. The target count remains eight; no ninth target is introduced.

### 10.1 `meta.json`

Add:

```text
historyHead, castEpoch, lastCommitOperationId, continuationModel
```

`continuationModel` remains omitted when absent; when present: `modelId, displayName?, sourceHint?, selectedAt`. `modelId` is `F3C-MODEL-1`; `displayName`/`sourceHint` are hints, never identity.

### 10.2 History entry

Add:

```text
turn: {
  batchId, participantIndex, commitSequence,
  userCommitOperationId, characterVersion, attemptId
}
generation: {
  modelId, executionFingerprint, parameterSchemaVersion,
  parameters, requestFingerprint, candidateFingerprint,
  promptDigest, seedMode, seed, effectiveSeed,
  templateFingerprint
}
```

Presence rules:

- `turn` is required for v2 user and assistant entries created after the amendment; absent for legacy.
- `generation` is required for v2 assistant entries; absent for user/system and legacy.
- `attachments` is always `[]` for v2; never null.
- `hasImage` remains a compatibility projection; true with `[]` is legacy-unknown.
- Unknown v2 keys are `ErrTargetMalformed`.
- `null` is not an omitted optional field.
- `history.ndjson` and `history.json` use the same presence and ordering rules.

### 10.3 Canonical field order

R4 does not change `F3C-CANON-1` byte formatting. New fields are inserted at the contract-declared order before the final hash; the exact order is frozen with the target schema version. Implementers may not choose map iteration order.

### 10.4 Target manifest and predecessor

- `f3c-target-manifest-2` includes the eight target hashes.
- The migration record/receipt includes `predecessorArchiveHash`; no migration is complete without it.
- A v1 receipt is never rewritten or rehashed.

---

## 11. Multi-character isolation and V24

### 11.1 Non-goals

R4 does not modify:

- character-specific context flushing;
- character memory isolation;
- character-specific prompt construction;
- cast synchronization semantics;
- multi-character interaction ordering at the prompt/context layer;
- character-specific lore visibility.

Persistence records the existing `characterId`, `characterVersion`, batch slot and generation provenance. It does not decide what a character sees or how a prompt is constructed.

### 11.2 V24 baseline

Before implementation, capture a reproducible pre-amendment baseline containing:

1. the exact build/source identity and executable hash;
2. test fixtures for narrator, one character, two characters, three characters and a removed character;
3. canonical prompt/context snapshots before generation;
4. memory visibility and context-flush observations;
5. cast synchronization event traces;
6. interaction traces for simultaneous and sequential character responses;
7. expected history ordering for the baseline product behavior.

After implementation, V24 reruns the same fixtures and compares:

- canonical prompt/context bytes;
- memory visibility;
- cast events;
- interaction outcomes;
- persisted history order.

A difference requires an explicit owner-approved contract change, not a test update. A newly invented interpretation of character isolation is not a baseline.

---

## 12. Gate v2 V1–V24 verification contract

### 12.1 Row record

Every row below has:

1. unique claim;
2. contract section;
3. implementation evidence;
4. test oracle;
5. positive cases;
6. negative/adversarial cases;
7. fault-injection cases;
8. native/runtime requirement;
9. Mimo requirement;
10. dependencies/prerequisites;
11. PASS criteria;
12. FAIL criteria;
13. UNRESOLVED criteria;
14. BLOCKING criteria;
15. evidence identity/provenance.

The standard interpretation is:

- **PASS:** all positive and negative cases behave as specified, the oracle is reproducible, fault cases recover as specified, native/Mimo prerequisites are satisfied, and evidence identity is recorded.
- **FAIL:** any required case violates the oracle.
- **UNRESOLVED:** a required case is missing, ambiguous, unexecuted, or has no reproducible oracle.
- **BLOCKING:** a prerequisite (isolation, native runtime, F3-F, Mimo, baseline) is absent.

No row passes because this proposal defines it. Gate v2 itself is not amended or frozen by R4.

### V1 — Canonicalization, fingerprints, Base64, hash chain

- **Claim:** `F3C-CANON-2` produces one stable fingerprint for identical effective requests and distinct fingerprints for material changes; hash chains and Base64 boundaries are exact.
- **Contract:** §5; Master canonicalization/hash rules.
- **Evidence:** golden vectors for objects, numbers, strings, arrays, maps, null/omitted, seeds, prompt digest and target manifests; recorded hashes.
- **Oracle:** recompute each vector with two independent serializers.
- **Positive:** reordered keys, explicit-vs-default same values, empty arrays, Unicode, zero.
- **Negative:** NaN, duplicate keys, fractional integer, different seed, different template, map order, malformed Base64.
- **Fault:** truncated hash input, receipt/hash mismatch.
- **Native/Mimo:** no native requirement; Mimo Stage 1/3 sees the vectors.
- **Dependencies:** frozen schema version.
- **PASS:** all vectors match; negatives reject.
- **FAIL:** any mismatch, acceptance of malformed input, or nondeterministic output.
- **UNRESOLVED:** vectors not independently reproduced.
- **BLOCKING:** no frozen canonical schema.
- **Evidence identity:** vector file hash, serializer version, command, output digest.

### V2 — Receipt inventory and reconstruction

- **Claim:** every commit has one reconstructable receipt; missing, duplicate, malformed or orphan receipts are detected.
- **Contract:** §2.4, §2.6.
- **Evidence:** receipt inventory fixture; reconstruct from commit record; hash/identity comparison.
- **Oracle:** exact receipt bytes equal the projection of the record.
- **Positive:** normal commit, crash after record, repeated recovery.
- **Negative:** missing receipt, duplicate receipt, receipt without record, record without receipt, wrong operation.
- **Fault:** truncate receipt, corrupt digest, partial install.
- **Native/Mimo:** no native; Mimo Stage 2/3.
- **Dependencies:** V1, V3.
- **PASS:** reconstruction is unique and idempotent.
- **FAIL:** ambiguous receipt or duplicate acceptance.
- **UNRESOLVED:** no crash fixture.
- **BLOCKING:** commit record authority not implemented.
- **Evidence identity:** fixture hash, receipt/record hashes, recovery log.

### V3 — Commit fork, orphan and missing inventory

- **Claim:** a fork, orphan target or missing commit is rejected; no silent adoption.
- **Contract:** §2.4–§2.6, §7.6.
- **Evidence:** branch fixtures at each commit boundary.
- **Oracle:** only one record with the exact parent/history head is accepted.
- **Positive:** linear commits, retry same operation, recovery.
- **Negative:** two records same parent, orphan target, record missing, conflicting candidate, stale operation.
- **Fault:** crash after record/target/receipt.
- **Native/Mimo:** Mimo Stage 2.
- **Dependencies:** V1, V2.
- **PASS:** conflict detected; no duplicate turn.
- **FAIL:** fork accepted or orphan treated as commit.
- **UNRESOLVED:** no fork fixture.
- **BLOCKING:** no authority.
- **Evidence identity:** repository hash, fixture IDs, commit record hashes.

### V4 — Exact target bytes, manifests, predecessor, new-chat v2

- **Claim:** v2 target bytes/manifests are exact; new chats are v2; migration binds predecessor; mixed manifests fail.
- **Contract:** §6, §10.
- **Evidence:** byte-level golden files, manifest hashes, new-chat and migration fixtures.
- **Oracle:** exact file hashes and manifest set; no extra/missing target.
- **Positive:** new chat, Save As, v1 migration with unknown fields.
- **Negative:** extra target, missing target, mixed v1/v2, unknown v2 key, missing archive hash, mismatched archive.
- **Fault:** partial target install, interrupted publish.
- **Native/Mimo:** Mimo Stage 1/2.
- **Dependencies:** V1, V23.
- **PASS:** exact bytes and manifest; failures reject.
- **FAIL:** mixed manifest or new v1 chat.
- **UNRESOLVED:** no byte-level fixture.
- **BLOCKING:** target schema not frozen.
- **Evidence identity:** target hash list, manifest hash, archive hash.

### V5 — Append, dirty state, autosave, session invalidation

- **Claim:** user/assistant append is durable only through a commit; autosave cannot overwrite a newer state; stale sessions fail.
- **Contract:** §2.5, §3.4, §7.
- **Evidence:** append/autosave/session fixtures; commit sequence logs.
- **Oracle:** exact history head and no lost update.
- **Positive:** normal append, concurrent autosave, retry.
- **Negative:** stale base, queued autosave after barrier, duplicate append, session mismatch.
- **Fault:** crash before/after commit record.
- **Native/Mimo:** Mimo Stage 2/3.
- **Dependencies:** V3, V6, V15.
- **PASS:** one authoritative append; stale fails.
- **FAIL:** lost update or late append.
- **UNRESOLVED:** no concurrency fixture.
- **BLOCKING:** no lease integration.
- **Evidence identity:** commit sequence trace, history hashes.

### V6 — Locks, publication, autosave vs switch, wipe vs save

- **Claim:** authority/lease/commit-gate ordering is deadlock-free and excludes late writes across switch, delete and wipe.
- **Contract:** §2.7–§2.9.
- **Evidence:** lock-order trace; blocked/queued writer fixtures; timeout/restart tests.
- **Oracle:** no deadlock; no write after seal; timeout blocks.
- **Positive:** concurrent commits, switch, delete, wipe.
- **Negative:** lock inversion, writer after closing, barrier timeout, commit gate stuck.
- **Fault:** process crash while draining, orphan lease.
- **Native/Mimo:** native stress; Mimo Stage 2.
- **Dependencies:** V5, V12, V15.
- **PASS:** all paths drain or block explicitly; no late write.
- **FAIL:** deadlock, late write, timeout release.
- **UNRESOLVED:** no lock trace.
- **BLOCKING:** no native stress.
- **Evidence identity:** trace hash, thread/process IDs, barrier records.

### V7 — Save As source binding and operation-ID conflicts

- **Claim:** Save As creates a new v2 chat bound to exact source evidence; it never copies a v1 receipt as v2 proof.
- **Contract:** §6.1, §6.4, §10.
- **Evidence:** source/target manifests; operation IDs; archive binding.
- **Oracle:** new ID, new v2 record, source hash preserved.
- **Positive:** v1 source, v2 source, unknown fields.
- **Negative:** source changes mid-copy, duplicate operation ID, missing archive.
- **Fault:** interruption before/after new chat commit.
- **Native/Mimo:** Mimo Stage 2.
- **Dependencies:** V4, V23.
- **PASS:** independent new chat; no receipt copy.
- **FAIL:** source mutation or v1 receipt reuse.
- **UNRESOLVED:** no source-change fixture.
- **BLOCKING:** no v2 target schema.
- **Evidence identity:** source/target hashes, operation ID.

### V8 — Tombstones and deletion recovery

- **Claim:** deletion is tombstone-authoritative; absent without tombstone is unknown; root mismatch blocks; no recreation.
- **Contract:** §2.8–§2.9, §7.5.
- **Evidence:** tombstone/absent/root-mismatch fixtures; barrier traces.
- **Oracle:** exact state-specific outcome.
- **Positive:** delete with active stream, delete after stream, restart.
- **Negative:** absent no tombstone, wrong root, tombstone with late write, repeated delete.
- **Fault:** crash after seal/tombstone/removal.
- **Native/Mimo:** native delete test; Mimo Stage 3.
- **Dependencies:** V6, V13, V22.
- **PASS:** no resurrection; unknown preserved.
- **FAIL:** absent cleanup, recreation, late write.
- **UNRESOLVED:** no tombstone evidence.
- **BLOCKING:** isolation not proven.
- **Evidence identity:** tombstone hash, barrier hash, root ID.

### V9 — Creation journal boundaries and identity

- **Claim:** creation is a journaled v2 operation; identity is revalidated; partial creation is recoverable.
- **Contract:** §2.5, §6.1.
- **Evidence:** creation fixtures, crash boundaries, identity checks.
- **Oracle:** no visible chat before commit; no duplicate chat.
- **Positive:** normal create, crash before/after record.
- **Negative:** duplicate ID, partial directory, tombstone collision.
- **Fault:** process death during creation.
- **Native/Mimo:** Mimo Stage 2.
- **Dependencies:** V4, V13.
- **PASS:** one durable creation or none.
- **FAIL:** duplicate/partial visible chat.
- **UNRESOLVED:** no creation fixture.
- **BLOCKING:** no creation journal.
- **Evidence identity:** creation record hash, chat ID.

### V10 — Memory migration conflict, commit, receipt/journal agreement

- **Claim:** memory migration follows the same authority; conflicts and receipt/journal disagreements block.
- **Contract:** §2.4, §7.6; Master memory rules.
- **Evidence:** memory fixtures, conflict and recovery records.
- **Oracle:** exact memory digest and receipt agreement.
- **Positive:** normal migration, retry.
- **Negative:** conflicting memory, missing receipt, journal disagreement.
- **Fault:** crash between memory and commit.
- **Native/Mimo:** Mimo Stage 2.
- **Dependencies:** V2, V3, V11.
- **PASS:** agreement; no duplicate migration.
- **FAIL:** silent overwrite or disagreement accepted.
- **UNRESOLVED:** no memory fixture.
- **BLOCKING:** no memory schema.
- **Evidence identity:** memory digest, receipt hash.

### V11 — History migration order, IDs, unknown fields

- **Claim:** history IDs and order are preserved; unknown v1 fields remain in the archive; no renumbering.
- **Contract:** §6.3, §10.2, Master history rules.
- **Evidence:** history migration fixture with duplicates/empty IDs/unknown fields.
- **Oracle:** exact ID set/order; archive byte comparison.
- **Positive:** clean v1, unknown keys, repeated recovery.
- **Negative:** renumbering, dropped unknown, duplicate ID, reordered entries.
- **Fault:** crash during history staging.
- **Native/Mimo:** Mimo Stage 3.
- **Dependencies:** V4, V23.
- **PASS:** exact preservation.
- **FAIL:** any loss or renumber.
- **UNRESOLVED:** no archive comparison.
- **BLOCKING:** no predecessor archive.
- **Evidence identity:** history hash, archive hash, ID list.

### V12 — Wipe tree hashes, gallery boundary, staging failure, root identity

- **Claim:** wipe stages every owned root, excludes gallery ownership from F3-C, and fails closed on root mismatch/staging failure.
- **Contract:** §2.8–§2.9, §6.7.
- **Evidence:** disposable root with complete tree manifest; staging failure; restore.
- **Oracle:** exact pre/post tree hashes; no user data outside disposable root.
- **Positive:** full wipe, restart after seal.
- **Negative:** root mismatch, gallery failure, partial stage, late writer, symlink escape.
- **Fault:** crash after first/second stage.
- **Native/Mimo:** native wipe; Mimo Stage 3.
- **Dependencies:** V6, V8, V14.
- **PASS:** complete owned-root wipe; no escape.
- **FAIL:** missed root, late write, real-root reach.
- **UNRESOLVED:** no isolated root.
- **BLOCKING:** disposable isolation not proven.
- **Evidence identity:** root manifest hash, stage paths, cleanup proof.

### V13 — Fault injection at every journal boundary

- **Claim:** every boundary produces one of the specified outcomes; no ambiguous adoption or duplicate.
- **Contract:** §2.4, §6.5, §7.6.
- **Evidence:** fault matrix for lease, stage, record, receipt, publish, restore, barrier, archive, wipe.
- **Oracle:** exact expected state per crash point.
- **Positive:** every boundary recovers.
- **Negative:** wrong boundary, missing record, duplicate recovery.
- **Fault:** injected crash/truncation at each point.
- **Native/Mimo:** native crash harness; Mimo Stage 3.
- **Dependencies:** V3–V12, V22.
- **PASS:** all boundaries match.
- **FAIL:** any ambiguous/duplicate state.
- **UNRESOLVED:** untested boundary.
- **BLOCKING:** no crash harness.
- **Evidence identity:** fault matrix hash, per-run logs, artifact hashes.

### V14 — Filesystem recovery: symlinks, reparse points, hardlinks, partial staging

- **Claim:** recovery cannot escape the disposable root; symlinks/reparse points/hardlinks/partial staging are detected and rejected or safely handled.
- **Contract:** §2.8, §6.7, Master filesystem rules.
- **Evidence:** path-traversal fixtures; inode/link checks; root canonicalization.
- **Oracle:** no write outside root; no ambiguous target.
- **Positive:** normal rooted paths, nested dirs.
- **Negative:** symlink escape, reparse escape, hardlink alias, partial stage, non-directory root.
- **Fault:** interrupted stage, broken link, replaced root.
- **Native/Mimo:** native Windows reparse test; Mimo Stage 2/3.
- **Dependencies:** V12, V23.
- **PASS:** escape blocked; safe recovery.
- **FAIL:** any outside-root write or alias acceptance.
- **UNRESOLVED:** no Windows reparse test.
- **BLOCKING:** native isolation absent.
- **Evidence identity:** canonical path list, file IDs/link counts, root hash.

### V15 — Race testing

- **Claim:** concurrent autosave, append, character attempts, switches, deletes and wipes obey one fence and one commit order.
- **Contract:** §2–§3, §7.
- **Evidence:** deterministic scheduler/race harness; two+ character attempts; delete/wipe races.
- **Oracle:** no lost update, duplicate, reorder, late write or deadlock.
- **Positive:** concurrent commits; two characters; delete with active streams.
- **Negative:** same base conflict, character removal, timeout, root wipe.
- **Fault:** crash during race.
- **Native/Mimo:** native stress; Mimo Stage 2/3.
- **Dependencies:** V5, V6, V8, V22, V24.
- **PASS:** all races match the specified outcome.
- **FAIL:** any lost/duplicate/late result.
- **UNRESOLVED:** race not executed.
- **BLOCKING:** no deterministic harness.
- **Evidence identity:** schedule seed, trace hash, result hashes.

### V16 — Native Wails/WebView2 preference/restore/unavailable

- **Claim:** native UI shows only terminal readiness; saved/requested/active remain distinct; restore/unavailable is truthful.
- **Contract:** §8.
- **Evidence:** native Wails/WebView2 run with real load, failure, restore, crash.
- **Oracle:** UI state matches backend terminal state; no false ready; no preference overwrite.
- **Positive:** ready, restore, unavailable, preference during load.
- **Negative:** early `LoadModel` return, monitor failure, crash.
- **Fault:** killed backend, port conflict, missing artifact.
- **Native/Mimo:** required native; Mimo Stage 3 may review evidence.
- **Dependencies:** V21.
- **PASS:** truthful UI and preference.
- **FAIL:** false ready or silent rewrite.
- **UNRESOLVED:** native not run.
- **BLOCKING:** no Wails/WebView2 environment.
- **Evidence identity:** app build hash, run log, UI/event trace.

### V17 — Independent Mimo Flash Fast challenge

- **Claim:** an independent Mimo challenge is completed in three stages: contract, implementation and evidence.
- **Contract:** Master Mimo requirements.
- **Evidence:** Mimo prompt/response artifacts, finding adjudication, resolved/rejected dispositions.
- **Oracle:** every finding is classified and dispositioned; no unreviewed stage.
- **Positive:** all three stages complete.
- **Negative:** combined mega-prompt, skipped stage, unreviewed finding.
- **Fault:** Mimo unavailable or incomplete.
- **Native/Mimo:** Mimo is the requirement.
- **Dependencies:** V1–V24 draft evidence.
- **PASS:** three complete stages and adjudication.
- **FAIL:** incomplete/invalid challenge.
- **UNRESOLVED:** not run.
- **BLOCKING:** Mimo unavailable.
- **Evidence identity:** Mimo artifact hashes, reviewer identity, adjudication record.

### V18 — Exact effective provenance, identity, request fingerprint, all paths

- **Claim:** every assistant turn records exact model artifact identity, execution provenance and effective request from all assistant paths.
- **Contract:** §4–§5, §7.1.
- **Evidence:** `F3C-MODEL-1`/`F3C-EXEC-1`/`F3C-REQ-1` vectors; narrator/character/sync paths.
- **Oracle:** recorded identity recomputes from artifacts/request.
- **Positive:** fixed/random seed, adapters, projector, character path.
- **Negative:** changed adapter/tokenizer/engine, requested-only values, missing artifact, path-only identity.
- **Fault:** unavailable artifact, backend version change.
- **Native/Mimo:** native model run; Mimo Stage 3.
- **Dependencies:** V1, V19, V24.
- **PASS:** exact provenance; unavailable is explicit.
- **FAIL:** wrong identity or fabricated value.
- **UNRESOLVED:** artifact set not exercised.
- **BLOCKING:** no model/native runtime.
- **Evidence identity:** artifact hashes, request envelope hash, provenance record.

### V19 — Explicit zero, random known/unknown, legacy absence

- **Claim:** zero, null, omitted, fixed/random seed and legacy absence are distinguishable and preserved.
- **Contract:** §5.2–§5.4, §10.2.
- **Evidence:** golden vectors for zero vs default, null vs absent, seed random/known/unknown, legacy no generation.
- **Oracle:** exact canonical bytes/fingerprint.
- **Positive:** all distinctions.
- **Negative:** `-1` as unknown, guessed seed, absent treated as zero.
- **Fault:** engine cannot report effective seed.
- **Native/Mimo:** Mimo Stage 3.
- **Dependencies:** V1, V18.
- **PASS:** all distinctions exact.
- **FAIL:** conflation or invented seed.
- **UNRESOLVED:** no seed evidence.
- **BLOCKING:** no effective request capture.
- **Evidence identity:** vectors, seed records, request hash.

### V20 — F3-F media ownership, present/corrupt/missing

- **Claim:** attachment bytes are durable only under F3-F; F3-C refs do not prove durability; missing/corrupt/unavailable behavior is truthful.
- **Contract:** §9.
- **Evidence:** F3-F contract and ownership/integrity/deletion/backup tests.
- **Oracle:** reference resolves to exact bytes or reports unavailable.
- **Positive:** present, corrupt, missing, external.
- **Negative:** false durability, substitution, cross-owner deletion.
- **Fault:** byte loss, hash mismatch.
- **Native/Mimo:** native media run; Mimo Stage 3.
- **Dependencies:** F3-F acceptance.
- **PASS:** only after F3-F evidence.
- **FAIL:** false durability.
- **UNRESOLVED:** F3-F incomplete.
- **BLOCKING:** F3-F ownership absent.
- **Evidence identity:** media hash, owner record, F3-F test evidence.

### V21 — Saved/requested/loading/ready/failed/restoring/unavailable preference rules

- **Claim:** saved preference, requested model and active runtime are separate; no asynchronous operation overwrites a newer explicit preference.
- **Contract:** §8.
- **Evidence:** state-machine fixtures, rapid preference changes, crash, restore.
- **Oracle:** exact state and saved value after each transition.
- **Positive:** select without save, explicit save, failed load, crashed transition.
- **Negative:** load completion saves, older operation overwrites newer, false ready, fallback.
- **Fault:** kill process during transition.
- **Native/Mimo:** native V16 run; Mimo Stage 3.
- **Dependencies:** V16, V18.
- **PASS:** separation and ordering hold.
- **FAIL:** saved/active conflation.
- **UNRESOLVED:** no native run.
- **BLOCKING:** native unavailable.
- **Evidence identity:** transition trace, preference sequence, saved hash.

### V22 — Intent identity, exact commit, absent/tombstone/root, idempotent recovery, barrier

- **Claim:** recovery adopts only exact commits; absent/root/tombstone outcomes are correct; recovery is idempotent and cannot bypass the barrier.
- **Contract:** §2, §7.
- **Evidence:** attempt/receipt/candidate fixtures; absent/tombstone/root; crash loops.
- **Oracle:** exact candidate/record identity.
- **Positive:** committed, interrupted, stale, removed, repeated recovery.
- **Negative:** unrelated receipt, wrong candidate, absent chat, root mismatch, late barrier write.
- **Fault:** crash at every attempt boundary.
- **Native/Mimo:** native restart; Mimo Stage 3.
- **Dependencies:** V3, V6, V8, V13, V18.
- **PASS:** exact/idempotent; failures block.
- **FAIL:** adoption/duplicate/resurrection.
- **UNRESOLVED:** no identity fixture.
- **BLOCKING:** no barrier/commit evidence.
- **Evidence identity:** attempt/candidate/receipt hashes, recovery trace.

### V23 — v1/v2 coexistence, predecessor archive, unknown fields, interrupted migration

- **Claim:** v1 remains valid; v2 migration is exact and recoverable; unknown fields are preserved; no mixed manifest or false rollback.
- **Contract:** §6, §10.
- **Evidence:** v1 fixture, archive, migration crash matrix, unknown fields, archive loss/corruption.
- **Oracle:** exact v1 bytes and v2 manifest; archive binding.
- **Positive:** clean/unknown-field migration, repeated recovery.
- **Negative:** missing/corrupt archive, mixed manifest, rehashed v1, claimed rollback, new v1.
- **Fault:** all boundaries in §6.5.
- **Native/Mimo:** native migration; Mimo Stage 3.
- **Dependencies:** V4, V7, V11, V13.
- **PASS:** exact archive/recovery; no false rollback.
- **FAIL:** loss or unprovable migration.
- **UNRESOLVED:** boundary untested.
- **BLOCKING:** archive not implemented.
- **Evidence identity:** archive hash, v1/v2 hashes, migration record.

### V24 — Multi-character baseline parity and per-character scopes

- **Claim:** concurrent character execution and serialized commits preserve existing context, memory, prompts, cast behavior and interaction; V24 baseline is captured and reproducible.
- **Contract:** §3, §11.
- **Evidence:** pre-amendment baseline artifacts; post-implementation reruns; two/three-character concurrent tests.
- **Oracle:** exact baseline comparison for prompts/context/memory/cast/history order.
- **Positive:** concurrent characters, removed character, retry, restart.
- **Negative:** cross-character context leak, reordered response, stale character commit, changed prompt.
- **Fault:** remove character during generation; delete during concurrent batch.
- **Native/Mimo:** native multi-character run; Mimo Stage 3.
- **Dependencies:** V15, V16, V18.
- **PASS:** frozen baseline parity and commit ordering.
- **FAIL:** any behavior change without approved contract.
- **UNRESOLVED:** baseline absent.
- **BLOCKING:** no pre-amendment baseline/native run.
- **Evidence identity:** baseline manifest/hash, run IDs, comparison hash.

### 12.2 Gate status

No row is PASS. V20 is **BLOCKED** pending F3-F. V16, V17 and V24 have no native/Mimo/baseline evidence. V14 has no native Windows reparse evidence. Gate v2 itself is unchanged and unfrozen.

---

## 13. Internal consistency audit

| Risk | R4 control |
|---|---|
| Two authorities granting write permission | A1/A2; all paths in §2.5 use `MutationLease`; preference is a root namespace, not a separate permission. |
| TOCTOU between check and write | Authority re-check and linearization at commit record (§2.4); staging itself is leased. |
| Stale epoch accepted after retry | Lease carries epoch/incarnation; retry cannot revive it (§2.3, A4). |
| Lock-order inversion | Fixed order and forbidden waits (§2.7). |
| Barrier releases on timeout | `blocked_timeout` is not release; no destructive step (§2.8–§2.9). |
| Producer termination mistaken for quiescence | Barrier drains leases and commit gates, not just producers (§2.8). |
| Concurrent characters lose updates | Response plan, prefix order, exact parent and CAS (§3.2–§3.4). |
| Commit order depends on scheduler | Participant order is durable; completion time cannot reorder (§3.2). |
| Character removal lets stale result commit | `characterVersion`/`castEpoch` and terminal `removed` slot (§3.4, §7.5). |
| Global lock serializes thinking | Commit gate is per chat and held only for commit; execution authority is separate (§2.7, §3.4). |
| Model identity conflated with execution | `F3C-MODEL-1` and `F3C-EXEC-1` separate (§4). |
| Fingerprint nondeterministic | `F3C-CANON-2` and versioned schema (§5). |
| Archive existence treated as integrity | `F3C-ARCH-1` manifest hash and readback verification (§6). |
| Rollback claimed without bytes/authority | No automatic in-place rollback; archive export only (§6.6). |
| Preference overwritten by load completion | Preference sequence and explicit operations (§8). |
| Migration/unknown fields deadlock | Archive preserves exact v1 bytes; v2 stays strict (§6.3). |
| Recovery recreates deleted data | Tombstone/barrier/root checks; no cleanup on unknown (§7.6). |
| Partial output promoted | Attempt state and `partialDigest`; no completed turn (§7.3). |

**Note on `F3C-CANON-2`:** this is a new canonicalization version. If the owner requires `F3C-CANON-1` byte compatibility for all v2 target hashing, the canonical envelope version must be resolved before freeze. R4 does not silently alter `F3C-CANON-1`; it uses `F3C-CANON-2` for new identity/fingerprint envelopes only. This is an explicit compatibility boundary.

---

## 14. Remaining owner decisions

R4 makes the following decisions explicitly; the owner must still accept or reject them before freeze:

1. **Batch ordering:** use the durable `participantOrder` in the response plan, with prefix commits and an explicit batch deadline. This preserves concurrent thinking but may delay a later character's visible result until an earlier slot is terminal.
2. **Model identity:** include tokenizer/vocabulary, adapters/LoRAs, auxiliary models, quantization/conversion artifacts and the exact execution provenance defined in §4.
3. **Canonicalization boundary:** `F3C-CANON-2` for new request/candidate/identity envelopes, while retaining `F3C-CANON-1` for legacy target bytes. Confirm this split is acceptable.
4. **Rollback:** accept that TARGET-2 has no automatic in-place rollback after v2 commit; the predecessor archive is immutable evidence/export only.
5. **Preference:** `SelectContinuationModel` never saves; only explicit `SetContinuationPreference` writes `continuationModel`; send uses active, not saved.
6. **Readiness:** `LoadModel` returns only terminal states; `ready` requires health/warm-up; `starting` is never success.
7. **Unknown fields:** keep the exact v1 bytes in the archive and strict v2 rejection; confirm no live `legacyExtras` field is required.
8. **F3-F:** V20 stays blocked until ownership, integrity, deletion and backup evidence exists.
9. **V24 baseline:** approve the baseline categories and require exact pre-amendment capture before implementation.
10. **V14/V16/V17:** native Windows reparse, Wails/WebView2 and Mimo evidence are prerequisites, not optional rows.

---

## 15. Unresolved risks and limitations

- No implementation, tests, builds, model loads, migrations, native runs, crash tests or destructive tests were performed.
- `F3C-MODEL-1`, `F3C-EXEC-1`, `F3C-REQ-1`, `F3C-CANON-2`, `F3C-CAND-2` and `F3C-ARCH-1` are proposed contracts, not evidenced behavior.
- R4 does not prove the current UI can supply a response plan, nor that current character generation invokes a batch plan. It defines the required contract; integration evidence is absent.
- R4 does not prove a per-chat commit gate can be implemented without blocking the Wails event loop; native stress evidence is absent.
- R4 does not prove current `SaveChat` can be converted from its multiple independent `saveJSON` calls into one atomic target-set commit; V4/V13 must test this.
- R4 does not prove F3-F media durability; V20 is blocked.
- R4 does not prove current model readiness; V16/V21 must run.
- R4 does not prove V14 filesystem isolation on Windows reparse points.
- R4 does not prove V24 baseline parity; no baseline exists.
- R4 cannot resolve the absence of a frozen owner decision on `F3C-CANON-2` versus `F3C-CANON-1` without an explicit acceptance.
- The current source has direct `DeleteChat`, multi-file `SaveChat`, async `performSave`, `LoadModel` return-before-ready, and transient image handling; R4 is a proposed target, not a claim that those paths are already compliant.

---

## 16. Status

**R4 proposal:** **PROPOSAL ONLY; NOT FROZEN; NOT ACCEPTED; NOT IMPLEMENTED; NOT VERIFIED.**  
**Frozen authority:** unchanged; SHA-256 verified.  
**F3-C phase 0:** **NOT CLOSED**.  
**P1a:** **UNAUTHORIZED**, downstream of P1c.  
**F3-C acceptance:** **NOT GRANTED**.  
**F3-D:** **BLOCKED**.  
**F3-B:** closed.  
**A9/A10/A12/A13:** closed; not reopened.  
**V20:** **BLOCKED** pending F3-F evidence.

No implementation authorization exists. No source, test, build, model load, migration, destructive operation, application data or authority record was changed.
