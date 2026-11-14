# AMENDMENT PROPOSAL — F3C-TARGET-2, Revision 5

**Status:** **PROPOSAL ONLY — NOT FROZEN AUTHORITY, NOT ACCEPTANCE, NOT IMPLEMENTATION.**  
**Date:** 2026-11-11  
**Supersedes for review:** R4 (`docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R4-2026-11-11.md`). R1–R4 and their reviews remain historical.  
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`  
**Architecture freeze:** `docs/review/FREEZE-F3C-ARCHITECTURE-2026-11-11.md` (`2d9220e8a9cfeb15e621836892d89b5bae23bec4b7ce7c7555d5af0d62fb315a`)  
**Gate v2:** `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md` (`7fdc30d19f392f5852e7de39b4ab5dccea533df7b9336c8205b735203bfedb2a`)  
**Scope:** documentation-only contract revision and read-only source trace. No source, tests, builds, model loads, migrations, destructive tests, application data, authority records, freeze, or implementation.

R5 resolves the fourteen findings from the independent R4 review and reconciles TARGET-2 with the frozen Master contract and adopted Gate v2. It is one proposed contract, not a patch list. Every rule below is unverified until the Gate evidence exists.

---

## 1. Purpose, retained decisions, and proposal boundary

### 1.1 What R5 retains

R5 retains these R4 decisions:

- chat identity, revisions, eight-target set, receipts, and historical proof;
- v1 coexistence and immutable v1 history;
- exact target bytes and exact-byte predecessor evidence;
- independent character generation with serialized shared-history commits;
- explicit response plans rather than completion-time ordering;
- per-chat continuation preference separate from historical provenance;
- model identity separate from execution provenance;
- F3-F ownership of media bytes and V20 blocked without F3-F evidence;
- no automatic in-place rollback after a v2 commit;
- V20, V14, V16, V17, and V24 as evidence-blocked until proven.

### 1.2 What R5 changes

R5 changes the authority architecture:

1. **No parallel authority lock hierarchy.** `MutationLease` is a durable authorization record, not a lock. It is issued and validated under the frozen Master lock hierarchy.
2. **`chatWriteMu` is the commit gate.** There is no separate `chat authority lock` or `commit gate` lock outside the frozen hierarchy.
3. **The commit record identifies a durable target package.** A digest alone no longer claims recovery capability.
4. **Only one saved-preference store exists:** `meta.json.continuationModel`, updated by an ordinary chat commit with a sequence compare-and-set.
5. **Canonicalization versions are separated by domain.** Frozen `F3C-CANON-1` remains the target/receipt/record canonicalization. Exact RFC 8785 is used only for new identity envelopes. No claim is made that RFC 8785 sorts by Unicode code point.
6. **Candidate identity excludes time.** Creation metadata is stored outside the fingerprint.
7. **Response plans, attempts, and candidates have durable operation records.** They are not target files and are not inferred from scheduler order.
8. **Gate v2 remains V1–V17.** V18–V24 are proposed additions, not adopted rows.

### 1.3 Explicit formal amendment

R5 requires one formal Master amendment, **AM-1 — target projection replacement**:

- Frozen Master §6.9 says commit records, receipts, tombstones, and targets are “installed create-only by same-volume atomic move. Final collisions fail.”
- A normal chat revision must replace the current target files (`meta.json`, `history.json`, etc.).
- Literal create-only installation cannot perform that replacement.

R5 therefore defines **target projection installation** as a same-volume atomic replacement from a verified commit package, while **commit records, receipts, tombstones, and all new evidence files remain create-only**. The replaced target's prior bytes remain in the immutable commit package.

This is not a silent reinterpretation. It is the proposed amendment **AM-1**. Until the owner accepts AM-1, R5 MUST NOT claim that ordinary v2 revision updates conform to the frozen §6.9 sentence. New-chat creation, which creates target paths for the first time, remains create-only.

**Owner decision:** OD-02.

---

## 2. Authority and one global lock order

### 2.1 Frozen lock order is controlling

The Master order is retained:

```text
castApplyMu
  -> chatSwitchMu
    -> chatTreeMu.RLock/Lock
      -> stateApplyMu(chat)
      -> chatWriteMu(chat)
        -> rootJournalMu

chatStateMu is a leaf.
```

R5 adds no lock that participates in this order. It adds one **runtime-only** lock outside it:

```text
modelTransitionMu
```

`modelTransitionMu` serializes model load/unload/restore. It MUST NOT be held while acquiring any F3-C lock, and no F3-C lock may wait for a model load, stream, or process stop while holding it. A model transition finishes its runtime work before an explicit preference commit enters the F3-C order.

### 2.2 Lock roles

| Lock | Scope | Role in R5 |
|---|---|---|
| `castApplyMu` | root | Cast apply/recovery and wipe exclusion. |
| `chatSwitchMu` | root | Active-chat/session transitions and deletion/wipe switching. |
| `chatTreeMu` | root tree | Root barrier against wipe; normal `.RLock`, wipe `.Lock`. |
| `stateApplyMu` | per chat | Authoritative state mutation and publication. |
| `chatWriteMu` | per chat | **The commit gate.** Targets, packages, records, receipts, leases, plans, attempts, candidates, verified reads. |
| `rootJournalMu` | root | Root journals and authority records; acquired only under the relevant `chatWriteMu`. |
| `chatStateMu` | App | Leaf registration/snapshot/publication only. No disk I/O and no other lock acquired while held. |
| `modelTransitionMu` | runtime | Model transition serialization only; outside F3-C. |

There is no separate `root authority lock`, `chat authority lock`, or `commit gate lock`. The durable authority records are protected by `chatWriteMu`/`rootJournalMu`; they do not create a second lock order.

### 2.3 Standard durable mutation

1. Allocate the operation and register it in the in-memory operation map under `chatStateMu`; release.
2. If the operation modifies cast state, acquire `castApplyMu`.
3. If the operation changes the active chat/session, acquire `chatSwitchMu`.
4. Acquire `chatTreeMu.RLock` for an ordinary mutation, or `chatTreeMu.Lock` for wipe.
5. Acquire `stateApplyMu(chat)` for state-changing mutations.
6. Acquire `chatWriteMu(chat)`.
7. Acquire `rootJournalMu` if a root journal or authority record is needed.
8. Under `chatWriteMu`, validate the operation, session, epoch, cast epoch, base, lease, and fence.
9. Create/validate the lease record, target package, commit record, receipt, and projections.
10. Release `rootJournalMu`, then `chatWriteMu`.
11. While retaining `stateApplyMu`, publish exact state under `chatStateMu` briefly.
12. Release `stateApplyMu` and the tree lock. Release cast/switch locks in reverse order.

No target/model/stream/disk deletion waits are permitted while a higher-order lock is held. `chatStateMu` is never used for disk I/O.

### 2.4 Barrier ordering

**Chat deletion:**

1. `castApplyMu` (if cast state is involved)
2. `chatSwitchMu`
3. `chatTreeMu.RLock`
4. `stateApplyMu(chat)`
5. `chatWriteMu(chat)`
6. `rootJournalMu` if the barrier record is root-journaled

Under `chatWriteMu`, R5 sets `closing`, revokes new lease issuance, and records the barrier. It releases the locks before waiting for active leases. No new commit can enter because the barrier state rejects it at the commit gate.

**Root wipe:**

1. `castApplyMu`
2. `chatSwitchMu`
3. `chatTreeMu.Lock`
4. For each chat in stable `(rootId, chatId)` order: `stateApplyMu(chat)`, `chatWriteMu(chat)`, close and record the barrier, then release.
5. Stage every owned root; commit the new root identity.

The barrier never holds `stateApplyMu` or `chatWriteMu` while waiting for a producer or stream. The tree lock prevents new root work; the per-chat barrier state prevents new leases.

### 2.5 Authority records

The durable authority records are operation evidence, not target files:

```text
f3c-authority/<F3C-PATH-AUTH-1(rootId, chatId)>/authority.json
f3c-authority/<F3C-PATH-AUTH-1(rootId, chatId)>/leases/<leaseId>.json
f3c-authority/<F3C-PATH-AUTH-1(rootId, chatId)>/barriers/<barrierId>.json
```

They are written under `chatWriteMu` and `rootJournalMu` as applicable, are included in the wipe inventory, and are validated by exact root/chat identity. They do not grant permission by existence alone: a lease is valid only when its holder token, process incarnation, session, epoch, cast epoch, fence epoch, and commit-gate validation all match.

New operation-root path hashes use the frozen path function:

```text
pathHash(domain, parts...) =
  lowercase_hex(SHA-256(domain || 0x00 || join(parts, 0x00)))
```

| Domain | Parts |
|---|---|
| `F3C-PATH-AUTH-1` | `rootId`, `chatId` |
| `F3C-PATH-PKG-1` | `rootId`, `chatId`, `operationId` |
| `F3C-PATH-PLAN-1` | `rootId`, `chatId`, `batchId` |
| `F3C-PATH-ATTEMPT-1` | `rootId`, `chatId`, `attemptId` |
| `F3C-PATH-CAND-1` | `rootId`, `chatId`, `candidateId` |
| `F3C-PATH-BATCH-1` | `rootId`, `chatId`, `batchId` |
| `F3C-PATH-PRED-2` | `rootId`, `chatId`, `predecessorRevision` |

Exact UTF-8 parts, one `0x00` separator, no trailing separator, lowercase 64-hex result. Lease and barrier filenames remain their canonical lowercase UUIDs.

---

## 3. MutationLease lifecycle

### 3.1 Lease record

```text
leaseVersion, leaseId, holderToken, rootId, rootIncarnation,
chatId, stateSessionId, stateEpoch, castEpoch,
characterId, characterVersion,
scope: chat | preference | migration | recovery | destructive,
operationId, operationKind,
fenceEpoch, baseRevision, baseHistoryHead,
issuedAt, lastHeartbeatAt, expiresAt, maxExpiresAt,
ownerProcessIncarnation,
status: issued | active | committing | committed | aborted |
        expired | orphaned | revoked
```

`holderToken` is a 256-bit random value generated by the authority and returned only to the holder. `leaseId` alone is not authorization.

### 3.2 Issuance

- The operation is registered under `chatStateMu`; no disk I/O occurs there.
- The durable lease is created under `chatWriteMu`, after `stateApplyMu`, before the first staged write.
- The lease is created only when the chat authority state is `active`, the barrier is absent, the root/incarnation match, and the session/epoch/cast values match current state.
- The lease records the exact base revision and history head against which the operation will commit.
- A lease MUST NOT be issued for a non-active chat or a sealed barrier.
- A lease from a previous `ownerProcessIncarnation` is never adopted. It is `orphaned`.

### 3.3 Lifetime

- Default `leaseTtl` is **120 seconds** from `issuedAt` or the last successful heartbeat.
- A holder MAY renew under `chatWriteMu` no more frequently than every 30 seconds and no later than 30 seconds before expiry.
- `maxExpiresAt` is **10 minutes** from issuance. Renewal MUST NOT extend a lease past it.
- A lease that reaches `expiresAt` or `maxExpiresAt` becomes `expired`. Expiry revokes commit authority; it does not authorize deletion or release of the barrier.
- A commit MUST validate that the lease is `active`, the token matches, `now < expiresAt`, and `now <= maxExpiresAt`.
- No lease is held during model inference. A generation attempt obtains an execution record, not a write lease, and obtains a new write lease only for a durable candidate/commit transaction.

### 3.4 Revocation and fencing

A lease is revoked when any of these changes or is detected:

- root identity or incarnation;
- chat tombstone/barrier state;
- `stateSessionId` or `stateEpoch`;
- `castEpoch`;
- `characterVersion` for a character-scoped operation;
- fence epoch;
- operation identity or holder token;
- process incarnation.

For every fencing change to root, session, chat, cast, character version, or barrier state, the authority MUST increment `fenceEpoch` under `chatWriteMu`, persist the new authority record, and mark every lower-epoch lease `revoked` before releasing `chatWriteMu`. Revocation is durable and idempotent; detection alone is not a substitute for incrementing the epoch. A revoked or expired holder cannot commit by retrying, renewing, replaying, or creating a new process with the old token. A new attempt requires a new operation identity and a new lease.

### 3.5 Commit-time proof

At the commit gate, under `chatWriteMu`, R5 MUST verify all of:

- lease status/token/holder/incarnation;
- lease expiry;
- root/root incarnation;
- chat identity and non-tombstoned state;
- state session and epoch;
- cast epoch;
- character identity/version where applicable;
- fence epoch;
- operation ID and kind;
- base revision/history head;
- no conflicting commit for the same operation;
- package and record hashes.

A mismatch is `ErrMutationFenceStale` or `ErrMutationLeaseInvalid`; no write occurs. The lease is terminalized as `revoked`/`orphaned` as appropriate.

### 3.6 Crash and startup

- On startup, the authority scans all `leases` records.
- A lease whose `ownerProcessIncarnation` is not the current incarnation is marked `orphaned`.
- An orphan with a matching valid commit record is reconciled as committed; projections are completed from the package.
- An orphan without a valid record is `aborted`; staged payloads are cleaned only after proving no commit record references them.
- An expired lease is terminal. It is never renewed.
- A process crash MUST NOT leave a writer that can later commit. The holder token and process incarnation fence a resumed process.

### 3.7 Barrier and timeout

- Closing a barrier revokes issuance and rejects new leases at `chatWriteMu`.
- Existing active leases must become `committed`, `aborted`, `expired`, or `orphaned` before the barrier seals.
- `blocked_timeout` is a durable barrier state, not a release. It preserves the barrier record, keeps new leases rejected, and permits only diagnosis/read-only access.
- A retry of the same barrier ID resumes at the drain/reconcile step; it never skips drain.
- A timeout never authorizes tombstone, deletion, or wipe staging.

---

## 4. Commit package, commit point, and recovery

### 4.1 Commit package

Every authoritative commit has an immutable package:

```text
f3c-packages/<F3C-PATH-PKG-1(rootId, chatId, operationId)>/
  package.json
  targets/
    meta.json
    history.json
    history.ndjson
    characters.json
    lorebooks.json
    memories_snapshot.json
    world.json          // if present
    scenario.json       // if present
```

`package.json` contains:

```text
packageVersion, packageHash, rootId, rootIncarnation, chatId,
operationId, operationType, baseRevision, committedRevision,
targetFormat: "F3C-TARGET-2",
targetManifest,
files: [{ path, sha256, bytes }]
```

`packageHash` is:

```text
"pkg1_" + lowercase_hex(
  SHA-256("F3C-PKG-1" || 0x00 || F3C-CANON-2(package.json without packageHash))
)
```

The package contains the exact target bytes, not paths or promises. The package is written and read back before the commit record.

### 4.2 Commit protocol

The sequence is normative:

1. Acquire the frozen lock sequence in §2.
2. Validate the lease, session, epoch, cast, base, and fence.
3. Build the complete target bytes and manifest in memory.
4. Write `package.json` and the exact target bytes under `f3c-tmp`; flush file contents and directory entries.
5. Read back every file; recompute every hash; verify the manifest and package hash.
6. Install the package at its operation-derived `f3c-packages/...` path using a same-volume atomic move. The package path is create-only; a collision must match the identical package hash or fail.
7. Write the commit record under `f3c-tmp`; flush; install it create-only at the frozen commit path. The record's `targetManifest` MUST equal the package's `targetManifest`, and its `targetPackageHash` MUST equal the package hash. **The atomic installation of this record is the single commit point.**
8. Install the target projections from the package using AM-1 atomic replacement; read back and verify.
9. Install the receipt envelope from the record's exact `receiptBodyBytes`; verify body, record, and envelope hashes.
10. Publish the package's exact `meta.Revision` and `meta.lastCommitOperationId` as in-memory state under `stateApplyMu`/`chatStateMu`; do not rewrite target bytes outside the package.
11. Mark the lease `committed`; release locks.

The record is authoritative only when its package is present and verified. A record whose package is missing/corrupt is an **evidence gap**, not a commit that can be guessed.

### 4.3 Crash recovery

| Crash point | Required outcome |
|---|---|
| Before package install | No record; abort; clean staging; prior state current. |
| Package installed, before record | No commit; package is orphaned evidence; clean only if no record references it; prior state current. |
| Record installed, package valid, targets not installed | Recover by installing targets from the package; verify; reconstruct receipt; publish. |
| Record and targets installed, receipt absent | Reconstruct receipt only from record body bytes; verify; publish. |
| Record, targets, and receipt installed, not published | Verify and publish. |
| Package missing/corrupt with record | Hard block `ErrCommitPayloadMissing`/`ErrCommitPayloadCorrupt`; do not adopt, delete, or invent. |
| Target missing/corrupt but package valid | Reinstall from package; if installation fails, block. |
| Record malformed or hash mismatch | Hard block; never overwrite. |
| Repeated recovery | Idempotent: same package/record/receipt/targets; no duplicate revision or turn. |

A target or receipt without a valid record is not a commit. A valid record with incomplete projections is recoverable only through its verified package. Neither rule may be inverted by an implementation.

### 4.4 Retention and garbage collection

- Commit packages are retained until **chat deletion or wipe**; R5 performs no automatic package garbage collection.
- A package may be collected only after a separately specified retention policy proves that no current/historical recovery or export path can reference it. R5 does not claim such a policy.
- A package collision must be byte-identical; otherwise the commit fails with `ErrCommitPackageConflict`.
- A missing package is never reconstructed from a manifest digest.

### 4.5 Idempotency

`(chatId, operationId)` is unique. A repeated commit with the same operation ID and same `candidateFingerprint` and package hash returns the existing proof only after verifying the package, record, receipt, and current targets. A different fingerprint or package is `ErrCommitOperationConflict`. A different operation ID is a different commit.

---

## 5. Canonicalization and identity domains

### 5.1 Frozen canonicalization

`F3C-CANON-1` remains unchanged and controls:

- target manifests inside receipts/records;
- receipt body bytes;
- commit record bytes;
- receipt envelope bytes;
- tree hashes and tombstones;
- all v1 and v2 target-manifest hashes.

It is not redefined by R5.

### 5.2 `F3C-CANON-2` for new identity envelopes only

`F3C-CANON-2` is **exactly RFC 8785 JSON Canonicalization Scheme**, without additions:

- UTF-8, no BOM;
- object property names sorted by UTF-16 code units, as RFC 8785 requires;
- arrays preserve order;
- strings use RFC 8785 escaping; no Unicode normalization;
- finite IEEE-754 binary64 numbers serialized by the RFC 8785/ECMAScript rules;
- no `NaN`, `Infinity`, or `-Infinity`;
- duplicate keys invalid;
- no trailing whitespace.

R5 does **not** claim that RFC 8785 sorts by Unicode code points.

`F3C-CANON-2` is used only for new model, execution, generation-request, candidate, plan, package, and archive identity envelopes. It does not replace `F3C-CANON-1` for target bytes, receipts, records, or historical proof.

### 5.3 Hash domains

| Identity | Hash domain and value |
|---|---|
| `modelId` | `mdl2_` + SHA-256(`F3C-MODEL-2\0` + canonical model envelope) |
| `executionFingerprint` | `exe2_` + SHA-256(`F3C-EXEC-2\0` + canonical execution envelope) |
| `generationRequestFingerprint` | `grq2_` + SHA-256(`F3C-GENREQ-2\0` + canonical generation request) |
| `candidateFingerprint` | `cand3_` + SHA-256(`F3C-CAND-3\0` + canonical candidate identity) |
| `engineConfigDigest` | `excfg1_` + SHA-256(`F3C-EXEC-CONFIG-1\0` + canonical execution controls) |
| `targetPackageHash` | `pkg1_` + SHA-256(`F3C-PKG-1\0` + package manifest without hash) |
| `responsePlanHash` | `plan1_` + SHA-256(`F3C-PLAN-1\0` + canonical plan) |
| `predecessorArchiveHash` | `arch2_` + SHA-256(`F3C-ARCH-2\0` + archive manifest without hash) |
| target hash | `tgt2-` + SHA-256(exact target bytes) |
| frozen `requestFingerprint` | `F3C-REQUEST-1` per Master §5.2 |
| frozen `commitRecordHash` | `F3C-COMMIT-RECORD-1` per Master §5.3 |
| frozen `receiptBodyHash` | `F3C-RECEIPT-BODY-1` |
| frozen `receiptEnvelopeHash` | `F3C-RECEIPT-ENVELOPE-1` |

No new field may reuse `requestFingerprint` for generation identity.

### 5.4 Identity vs metadata

- `candidateCreatedAt`, `selectedAt`, `lastHeartbeatAt`, and operation timestamps are metadata; they are excluded from candidate/model/generation identity.
- `origin` is metadata about how an effective value was chosen. It is excluded from `generationRequestFingerprint`; two requests with the same effective values hash the same even if one value was inherited and one explicit. `seedMode` is included because it changes generation semantics.
- `engineExtensionParameters` is a canonical map. A key that duplicates a named parameter MUST be rejected; it is not silently merged.
- Optional fields are omitted only where the schema says omitted. `null` is not an alias for omission.

---

## 6. Complete durable record schemas

### 6.1 `F3C-TARGET-2` target set

The target set remains exactly eight paths:

| Path | Presence | Content |
|---|---|---|
| `meta.json` | required | Chat metadata and continuation preference. |
| `history.json` | required | Ordered history array. |
| `history.ndjson` | required | Same entries, one compact JSON line each. |
| `characters.json` | required | Chat character contexts, frozen order. |
| `lorebooks.json` | required | Chat-scoped lorebook snapshot. |
| `memories_snapshot.json` | required | Generated lore snapshot. |
| `world.json` | optional; absent if none | World snapshot. |
| `scenario.json` | optional; absent if none | Scenario snapshot. |

`characters.json`, `lorebooks.json`, `memories_snapshot.json`, `world.json`, and `scenario.json` retain the frozen `F3C-TARGET-1` field order and semantics. R5 changes only `meta.json` and history entries.

### 6.2 `meta.json`

Exact field order:

```text
id, name, createdAt, updatedAt, worldName, scenarioName,
characterNames, lorebookNames, worldId, scenarioId,
activeCharacters, removedCharacters, addedCharacters,
activeLorebooks, removedLorebooks, addedLorebooks,
rootId, revision, lastCommitOperationId,
continuationSeq, continuationModel
```

- `continuationSeq` is an unsigned integer, starting at `0`; every accepted explicit preference write increments it by exactly one.
- `continuationModel` is omitted when absent. When present, exact order:

```text
modelId, displayName, sourceHint, selectedAt
```

`modelId` and `selectedAt` are required. `displayName` and `sourceHint` are optional and omitted when empty. `selectedAt` is RFC3339. `sourceHint` is non-authoritative.

### 6.3 History entry

Exact field order:

```text
id, sender, text, role, timestamp, hasImage, characterId,
turn, generation, attachments
```

- `turn` is present for v2 user and assistant entries; omitted for legacy entries.
- `turn` order:

```text
batchId, participantIndex, commitSequence,
userCommitOperationId, characterVersion, attemptId
```

`characterVersion` is `0` for narrator/user slots with no character. `participantIndex` is from the response plan.

- `generation` is present only for v2 assistant entries; omitted for user/system and legacy entries.
- `attachments` is always `[]` for v2 entries, never `null`. `hasImage` remains a compatibility projection; `true` with `[]` is legacy-unknown.
- `history.json` and `history.ndjson` must decode to the same ordered entries and values.

### 6.4 `generation` object

Exact field order:

```text
modelId, displayName, modelSchemaVersion,
engine, engineVersion, backend,
executionFingerprint,
parameterSchemaVersion, generationRequestFingerprint,
defaultSetId,
parameters,
promptDigest, templateFingerprint, adapterConfigDigest,
seedMode, seed, effectiveSeed, seedStatus,
execution
```

- `displayName` optional; omitted if absent.
- `modelSchemaVersion` is `F3C-MODEL-2`.
- `parameterSchemaVersion` is `F3C-GENPARAMS-2`.
- `seed` appears only for `seedMode = explicit`; it is a signed 64-bit integer and `0` is valid.
- `effectiveSeed` appears only when `seedMode = random` and the engine reports it.
- `seedStatus` is `known | unknown | not_applicable`.
- `execution` order:

```text
engineName, engineVersion, engineBuildId, runtimeName,
runtimeVersion, executableHash, acceleration, engineConfigDigest,
provenanceStatus
```

  `engineBuildId`, `runtimeVersion`, and `executableHash` may be omitted only when unavailable, in which case `provenanceStatus` is `partial`/`unavailable`; no value is invented. `acceleration` order is `backend, device, offloadMode, offloadLayers`; `offloadLayers` optional.

### 6.5 `parameters` object

Exact order and effective values:

```text
temperature, topPEnabled, topP, topK, minP,
typicalPEnabled, typicalP, tfs,
repeatPenalty, repeatLastN,
frequencyPenaltyEnabled, frequencyPenalty,
presencePenaltyEnabled, presencePenalty,
mirostatEnabled, mirostat, mirostatTau, mirostatEta,
dynamicTempRangeEnabled, dynamicTempRange, dynamicTempExponent,
dryMultiplier, dryAllowedLength, dryBase,
smoothingFactor, smoothingCurve,
topAEnabled, topA, epsilonCutoff, etaCutoff,
encoderRepeatPenalty, noRepeatNGramSize,
maxTokens, nKeep, contextTokens,
stopSequences, samplerOrder,
grammar, jsonSchema, logitBias, ignoreEOS
```

All numeric fields are finite numbers; integer fields are integers. `stopSequences` and `samplerOrder` are arrays; `grammar` and `jsonSchema` are string-or-null; `logitBias` is a sorted map; `ignoreEOS` is boolean. Defaults follow the frozen engine defaults and the selected family overrides; `defaultSetId` identifies the exact base/family default set.

### 6.5.1 `F3C-DEFAULTS-2`

The base defaults are the current `llama.NewInferenceOptions()` values, before family overrides and explicit settings:

| Parameter | Base default |
|---|---|
| `temperature` | `0.7` |
| `topPEnabled` | `true` |
| `topP` | `0.9` |
| `topK` | `40` |
| `minP` | `0.05` |
| `typicalPEnabled` | `true` |
| `typicalP` | `1.0` |
| `tfs` | `1.0` |
| `repeatPenalty` | `1.1` |
| `repeatLastN` | `64` |
| `frequencyPenaltyEnabled` | `true` |
| `frequencyPenalty` | `0.0` |
| `presencePenaltyEnabled` | `true` |
| `presencePenalty` | `0.0` |
| `mirostatEnabled` | `true` |
| `mirostat` | `0` |
| `mirostatTau` | `5.0` |
| `mirostatEta` | `0.1` |
| `dynamicTempRangeEnabled` | `true` |
| `dynamicTempRange` | `0.0` |
| `dynamicTempExponent` | `1.0` |
| `dryMultiplier` | `0.0` |
| `dryAllowedLength` | `2` |
| `dryBase` | `1.0` |
| `smoothingFactor` | `0.0` |
| `smoothingCurve` | `1.0` |
| `topAEnabled` | `true` |
| `topA` | `0.0` |
| `epsilonCutoff` | `0.0` |
| `etaCutoff` | `0.0` |
| `encoderRepeatPenalty` | `1.0` |
| `noRepeatNGramSize` | `0` |
| `maxTokens` | `512`, plus any effective thinking budget |
| `nKeep` | `0` |
| `contextTokens` | effective runtime context, not a fixed UI default |
| `stopSequences` | template/engine stop set |
| `samplerOrder` | `[]` (engine order) |
| `grammar` | `null` |
| `jsonSchema` | `null` |
| `logitBias` | `{}` |
| `ignoreEOS` | `false` |
| `seedMode` | `random` |
| `seed` | `null` in random mode; explicit integer when fixed |
| `effectiveSeed` | reported random seed when available; otherwise `null` |
| `seedStatus` | `known`, `unknown`, or `not_applicable` |

Resolution order: base defaults, selected family overrides, explicit/inherited settings, runtime normalization (`topPEnabled = false` -> `topP = 1.0`; thinking budget added to `maxTokens`), then seed resolution. `defaultSetId` is:

```text
"dfl2_" + lowercase_hex(
  SHA-256("F3C-DEFAULTS-2" || 0x00 || F3C-CANON-2({
    engineName, engineVersion, family,
    baseDefaults, familyOverrides
  }))
)
```

Changing the base or family default table requires a new schema version and fresh golden vectors. The `-1` engine sentinel is never an explicit fixed seed.

### 6.6 Attachment entry

```text
mediaId, role, provider, mimeType, byteSize, contentHash
```

`contentHash` is required for `f3f-media` and optional for `legacy-indexeddb`. No bytes, base64, local path, thumbnail, or generated placeholder is a target field.

### 6.7 Commit record

The frozen commit record fields are retained. R5 adds only:

```text
targetFormat: "F3C-TARGET-2"
targetPackageHash
generationRequestFingerprint?     // assistant-producing commits
responsePlanHash?                 // user/batch commits
predecessorArchiveHash?           // v1-to-v2 migration
```

Exact record order:

```text
version, rootID, chatID, operationID, operationType,
requestFingerprint, baseRevision, committedRevision,
candidateFingerprint, targetManifest, targetFormat,
targetPackageHash,
receiptBodyBytes, receiptBodyHash,
previousCommitOperationID, previousCommitRecordHash,
previousReceiptEnvelopeHash, commitSequence, createdAt,
generationRequestFingerprint, responsePlanHash,
predecessorArchiveHash, memoryMigrationJournalHash,
commitRecordHash
```

Optional fields are omitted when not applicable; they are never `null`. `requestFingerprint` remains the frozen **operation** identity. `generationRequestFingerprint` is the new effective-generation identity. `commitRecordHash` remains the frozen domain over the record excluding that field.

`targetFormat` selects the record's fingerprint semantics. For a v2 record, `candidateFingerprint` is exactly the `cand3_` value defined in §9.2. A v1 record retains its original candidate fingerprint and is never recomputed with `F3C-CAND-3`. A verifier MUST select the algorithm from the record format, not from the field name.

### 6.8 Receipt body and envelope

The frozen receipt body fields and order are retained exactly:

```text
version, rootID, chatID, operationID, operationType,
requestFingerprint, baseRevision, committedRevision,
candidateFingerprint, targetManifest,
previousCommitOperationID, previousCommitRecordHash,
previousReceiptEnvelopeHash, commitSequence, createdAt,
result payload,
memoryMigrationJournalHash?      // optional; memory migration only
```

R5 adds the following optional fields after the frozen result payload and before any optional memory migration field:

```text
targetFormat, targetPackageHash,
generationRequestFingerprint, responsePlanHash,
predecessorArchiveHash
```

The body excludes all hash fields. The receipt envelope remains exactly:

```text
receiptBodyBytes, receiptBodyHash, commitRecordHash, receiptEnvelopeHash
```

`receiptBodyBytes` is exact Base64 of the canonical body. Old receipts have none of the optional fields and remain valid.

### 6.9 Response plan

Path:

```text
f3c-generations/plans/<F3C-PATH-PLAN-1(rootId, chatId, batchId)>/plan.json
```

Fields:

```text
planVersion, planId, batchId, rootId, rootIncarnation, chatId,
userCommitOperationId, historyBaseHead,
contextMode: captured | prefix_committed,
participantOrder: [
  { participantIndex, characterId, characterVersion }
],
deadlineAt?,
createdAt
```

- The plan header is written create-only before the user-turn commit and is immutable. It has no mutable status field.
- The user-turn commit record includes `responsePlanHash`; the plan's `userCommitOperationId` must match.
- A plan that is not referenced by a valid user commit is orphaned; it cannot start an attempt.
- `contextMode = captured` means each slot records the history head used by the existing prompt builder at attempt start; no prompt is rewritten. `prefix_committed` means a slot may start only after lower slots are terminal/committed and uses the resulting head. The default is `captured`; selecting `prefix_committed` is an owner decision.

### 6.9.1 Batch and slot state

Mutable batch/slot state is separate from the immutable plan header:

```text
f3c-generations/batches/<F3C-PATH-BATCH-1(rootId, chatId, batchId)>/batch.json
```

Fields:

```text
batchVersion, batchId, planHash, rootId, rootIncarnation, chatId,
status: prepared | active | partially_committed | completed | failed | aborted,
slots: [
  { participantIndex, characterId, characterVersion,
    status: planned | started | streaming | completed | committed |
            failed | aborted | removed | skipped | abandoned | blocked,
    attemptId?, terminalReason? }
],
updatedAt
```

- `batch.json` is written under `chatWriteMu` and is the authoritative aggregate/slot state.
- Every slot state used in §7 MUST appear in this schema. `planned` and `blocked` are nonterminal; `committed` and the no-turn terminals are terminal for that slot.
- `planHash` MUST equal `responsePlanHash` computed over the immutable plan header and stored in the user commit; slots MUST correspond exactly and in order to `participantOrder`.
- Batch status is derived from slot states and the explicit `AbandonBatchSlot` operation; it is persisted, not inferred only at read time.

### 6.10 Attempt and candidate

Attempt path:

```text
f3c-generations/attempts/<F3C-PATH-ATTEMPT-1(rootId, chatId, attemptId)>/attempt.json
```

Candidate path:

```text
f3c-generations/candidates/<F3C-PATH-CAND-1(rootId, chatId, candidateId)>/candidate.json
```

Attempt fields:

```text
attemptId, attemptOrdinal, operationId, operationType,
rootId, rootIncarnation, chatId,
stateSessionId, stateEpoch, castEpoch,
characterId, characterVersion,
batchId, participantIndex, planHash, contextHistoryHead,
userCommitOperationId, baseRevision, baseHistoryHead,
assistantMessageId, modelId, executionFingerprint,
generationRequestFingerprint, candidateFingerprint,
status, terminalReason?, partialDigest?, startedAt
```

Candidate fields:

```text
candidateId, batchId, participantIndex, planHash,
rootId, rootIncarnation, chatId,
stateSessionId, stateEpoch, castEpoch,
characterId, characterVersion, attemptId, attemptOrdinal,
userCommitOperationId, baseRevision, baseHistoryHead,
contextHistoryHead, assistantMessageId,
assistantContentDigest, targetManifestDigest,
targetPackageHash, candidateFingerprint,
generationRequestFingerprint, modelId, executionFingerprint,
status, createdAt
```

Candidate statuses:

```text
prepared | staged | eligible | commit_pending | committed |
rejected | failed | aborted | removed | skipped | abandoned |
conflict
```

`terminal_no_turn` is the set `{failed, aborted, removed, skipped, abandoned}`. `committed` and `conflict` are terminal. `prepared`/`staged`/`eligible`/`commit_pending` are nonterminal.

### 6.11 Predecessor archive

Path:

```text
f3c-predecessors/<F3C-PATH-PRED-2(rootId, chatId, predecessorRevision)>/
  archive.json
  targets/...
  evidence/...
```

Archive fields:

```text
archiveVersion, archiveHash, rootId, rootIncarnation, chatId,
predecessorRevision, predecessorManifest,
files: [{ path, sha256, bytes }],
v1RecordHash, v1ReceiptHash, createdAt
```

`archiveHash` uses `F3C-ARCH-2` and excludes itself. The archive contains exact v1 target bytes and exact v1 record/receipt evidence. It is immutable and create-only.

### 6.12 Barrier, authority, and lease records

Barrier fields:

```text
barrierVersion, barrierId, rootId, rootIncarnation, chatId,
scope: chat | root, phase, barrierEpoch, sealSequence,
createdAt, sealedAt?, completedAt?, timeoutAt?, reason
```

Authority fields:

```text
authorityVersion, rootId, rootIncarnation, chatId, state,
fenceEpoch, commitSequence, historyHead, lastCommittedOperationId,
activeLeaseCount, barrierId?, recoveryIncarnation
```

Lease fields are defined in §3.1. All three are operation evidence and are included in wipe inventory.

---

## 7. Response plan, candidate authority, and concurrent characters

### 7.1 User-turn commit and plan

A user turn MUST be durably committed before the first generation attempt. Its commit package contains the user entry and the `responsePlanHash` in the record/receipt; the plan record is already create-only. If the user commit fails, the plan remains `prepared` and no attempt may start.

### 7.2 Participant order

`participantOrder` is durable and immutable. A candidate for index `n` may commit only when every lower index is `committed` or in `terminal_no_turn`.

The coordinator never uses completion time, scheduler order, map iteration, or UI array order to determine history order.

### 7.3 Prefix commit

The expected parent is:

```text
historyBaseHead
  followed by the committed candidate hashes for lower indices
```

It must equal the current `historyHead`. A mismatch is `ErrGenerationCommitConflict`; no overwrite, no merge, no silent append.

History order is `participantIndex` order. A candidate’s own `contextHistoryHead` records what it actually saw; prefix order does not rewrite its prompt.

### 7.4 Terminal states and deadlines

Slot states are exactly those in §6.9.1:

- `planned`: no attempt has started.
- `started`/`streaming`: attempt exists and is nonterminal.
- `completed`: generation ended; candidate is pending commit.
- `committed`: durable assistant turn.
- `failed`, `aborted`, `removed`, `skipped`, `abandoned`: terminal no-turn states.
- `blocked`: waiting for a lower slot; not terminal.

Batch states are exactly `prepared | active | partially_committed | completed | failed | aborted`. A lower slot that never starts blocks later slots unless an explicit `AbandonBatchSlot` operation writes `abandoned` into `batch.json` under `chatWriteMu`. The operation records the reason and invalidates any attempt for that slot. A batch deadline MAY trigger the same durable operation only if `deadlineAt` is present and the owner has enabled deadline abandonment. There is no implicit deadline side effect.

A partially committed batch remains `partially_committed`. Committed turns are visible; uncommitted slots are not displayed as assistant turns. The batch cannot report `completed` while a slot is nonterminal.

### 7.5 Character removal and cast changes

- Removal increments `castEpoch` and increments `characterVersion` for that character.
- The slot becomes `removed`; its active attempt becomes invalid.
- The commit gate checks `castEpoch`, `characterVersion`, and character membership.
- A recreated character receives a new ID; a same-display-name character never inherits an old attempt.
- Chat deletion/wipe invalidates every plan, attempt, candidate, lease, and character slot in scope.

### 7.6 Batch recovery

- On restart, the plan header is read from `f3c-generations/plans`, and the durable batch/slot state from `f3c-generations/batches`.
- Each attempt/candidate is read by its identity path.
- The first nonterminal slot is resumed only if its authority fields still match the plan and current state.
- Committed slots are never replayed.
- A candidate that cannot prove its plan, slot, session, epoch, cast, or base is `conflict`, not adopted.
- Recovery never creates a missing assistant turn from partial stream text.

### 7.7 Multi-character preservation

R5 does not alter character prompt construction, context flushing, memory isolation, cast synchronization, lore visibility, or interaction semantics. With `contextMode = captured`, the existing prompt builder’s history head is recorded rather than rewritten. With `prefix_committed`, later slots wait for earlier commits; this is an explicit owner choice because it may delay generation. Single-character operation uses a singleton plan and does not acquire multi-character behavior.

---

## 8. Model artifact and execution provenance

### 8.1 `modelId`

Envelope:

```text
schemaVersion: "F3C-MODEL-2"
components: [
  { role, roleQualifier, artifactHash, byteSize, format, required }
]
```

`role` is one of:

```text
main-model, tokenizer, vocabulary, projector, vision-encoder,
adapter, auxiliary-model, embedding-model, scoring-model,
quantization-artifact, conversion-metadata, other-inference-artifact
```

- `roleQualifier` is required and nonempty when two components share a role; otherwise it is omitted.
- Sort components by `(role, roleQualifier or "", artifactHash, byteSize)`.
- `artifactHash` is `sha256:<hex>` over exact bytes.
- `format` is one of `gguf`, `safetensors`, `bin`, `json`, `text`, `image`, `embedded:<component>`, or `other:<label>`; the label is part of identity.
- `required` is true when absence blocks generation. Optional absent components are omitted; `required: false` and omission are not interchangeable.
- An embedded component uses the main artifact hash and `embedded:<component>`.
- Any required unhashable component blocks generation with `ErrModelArtifactIncomplete`; no placeholder ID.

`modelId` is `mdl2_` plus the domain-separated hash in §5.3. Paths, basenames, labels, and mtimes are hints, never identity.

### 8.2 Execution provenance

`F3C-EXEC-2` records:

```text
engineName, engineVersion, engineBuildId, runtimeName, runtimeVersion,
executableHash, acceleration { backend, device, offloadMode, offloadLayers },
loadControls {
  contextSize, batchSize, gpuLayers, threads, useMMap, useMLock,
  memoryFraction, flashAttention, offloadKv,
  ropeScalingType, ropeFreqBase, ropeFreqScale,
  mainGPU, splitMode
},
tokenizerControls { addBOS, addEOS, padTokenID },
engineConfigDigest,
provenanceStatus: complete | partial | unavailable
```

`engineConfigDigest` is `excfg1_` over the canonical `F3C-EXEC-CONFIG-1` envelope of these controls. It MUST NOT recurse through `F3C-REQ-2` or the generation request. Load/context/tokenizer controls that can change inference are included; machine-specific acceleration is provenance, not model identity.

### 8.3 Historical availability

Historical turns retain `modelId`, `executionFingerprint`, `generationRequestFingerprint`, and the exact `execution` snapshot in `generation`. Uninstalling or replacing a model does not rewrite history. Current artifact observation is `available | missing | relocated | replaced | corrupt | unreadable`; a mismatch is reported, not substituted.

---

## 9. Generation request, candidate identity, and fingerprints

### 9.1 Generation request

`F3C-GENREQ-2` envelope:

```text
schemaVersion, modelId, executionFingerprint, defaultSetId,
promptDigest, templateFingerprint, adapterConfigDigest,
parameters, seedMode, effectiveSeed, seedStatus,
engineExtensionParameters
```

`promptDigest` is `sha256:<hex>` over canonical prompt messages; raw prompt text is not stored. `origin`, timestamps, display labels, and requested-only settings are excluded.

`generationRequestFingerprint` is `grq2_` plus the domain hash. It is not the frozen operation `requestFingerprint`.

### 9.2 Candidate identity

`F3C-CAND-3` envelope:

```text
schemaVersion, candidateId, rootId, rootIncarnation, chatId,
stateSessionId, stateEpoch, castEpoch,
characterId, characterVersion,
batchId, participantIndex, planHash,
userCommitOperationId, userMessageId,
attemptId, attemptOrdinal,
baseRevision, baseHistoryHead, contextHistoryHead,
assistantMessageId, assistantContentDigest,
targetManifestDigest, targetPackageHash,
generationRequestFingerprint, modelId, executionFingerprint
```

Excluded: `createdAt`, display labels, raw prompt text, secrets, requested-only model labels. `candidateFingerprint` is `cand3_` plus the domain hash. The same logical candidate always produces the same fingerprint regardless of creation time.

`candidateId` is allocated at attempt start and is stable across retries of the same attempt. A retry is a new attempt/ordinal.

### 9.3 Commit identity

The commit record binds:

- frozen operation `requestFingerprint`;
- frozen `candidateFingerprint` compatibility field, set to the `cand3_` value;
- `generationRequestFingerprint`;
- `targetPackageHash`;
- `responsePlanHash` for batch/user commits;
- `predecessorArchiveHash` for migration.

No field aliases another. A verifier reads the schema version and computes each independently.

---

## 10. Continuation preference and model transition

### 10.1 One authoritative saved preference

`meta.json.continuationModel` is the sole authoritative saved preference. There is no authoritative root `preferences.json`. Any cache of the preference is rebuildable and MUST NOT be used to overwrite `meta`.

`meta.continuationSeq` is the compare-and-set sequence.

### 10.2 Operations

- `SelectContinuationModel(modelId)` loads/activates the model; it does not write preference.
- `SetContinuationPreference(modelId, expectedSeq, operationId)` writes only `meta.continuationModel` and increments `continuationSeq`; it may save an unavailable model but never claims active.
- `LoadModel`/auto-load never writes preference.
- A failed load never changes saved preference.
- An older load/selection cannot overwrite a newer explicit preference; completion checks `continuationSeq`.
- Two rapid preference writes serialize under `chatWriteMu`; `expectedSeq` mismatch returns `ErrPreferenceConflict`.
- A failed preference commit leaves the previous saved value.
- Sending uses the ready active model. If no active model exists, it loads the saved preference; if unavailable, `ErrModelUnavailable`, no fallback.
- Chat deletion/wipe invalidates pending preference operations.

### 10.3 Runtime transition

Runtime states are `preflight | stopping | starting | health_check | ready | active | restoring | failed | unavailable | cancelled`. Only `active` may generate. `LoadModel` returns only a terminal result (`ready`, `failed`, `unavailable`, `cancelled`); `starting` is never success. `ready` requires health and a generation-capable warm-up, not process/port existence.

The transition captures a full prior snapshot (model identity, projector/adapters, spawn settings, admission/eligibility evidence, resource estimate). On preflight failure, old model is untouched. On post-stop failure, kill the partial process, attempt one restore; if restore fails, `unavailable`; never false-ready, never substitute.

`modelTransitionMu` is released before any preference commit enters the F3-C lock order. No F3-C lock is held during model load/stop.

---

## 11. Migration, import, and predecessor archive

### 11.1 One commit protocol

Migration uses the same package/record/receipt protocol as ordinary v2 writes. It does not install v2 targets before the record.

Sequence:

1. Validate v1 current, receipt, chain, manifest, root identity.
2. Create and verify the immutable predecessor archive; compute `predecessorArchiveHash`.
3. Build/verify the v2 target package.
4. Install the package and write the commit record with `predecessorArchiveHash`; the record is the commit point.
5. Install target projections (AM-1), receipt, and publish v2.
6. Keep the archive until chat deletion/wipe; no automatic in-place rollback.

### 11.2 Archive binding

The archive contains exact v1 target bytes and exact v1 record/receipt evidence. `archiveHash` covers the archive manifest; the migration record/receipt carries the same value.

- If the archive is missing, corrupt, or hash-mismatched **before** the migration record is installed, migration MUST fail with `ErrPredecessorEvidenceMissing`; no v2 commit or completion marker is written.
- If a committed v2 later loses or corrupts its archive, the chat is marked `predecessorEvidence: incomplete` for export/rollback purposes. The v2 commit itself remains governed by its package/record/targets; v1 bytes are never invented or substituted.

A missing/corrupt archive is never treated as migration-complete evidence.

### 11.3 No in-place rollback

There is no automatic in-place rollback after a v2 commit. The archive supports forensic export/reconstruction, not a claim that the live chat was rolled back. A separate maintenance operation would require its own barrier and owner authorization.

### 11.4 Import and frozen Q3

Import remains the frozen Q3 operation: create-only, idempotent by identity plus content hash, journaled, reversible, and browser data untouched. R5 does not replace or weaken the import journal. If import writes v2 targets, it uses the same package/record/receipt protocol and preserves its own journal and rollback scope. Import is not silently equated with Save As.

---

## 12. Wipe and complete inventory

### 12.1 Wipe commit point

The new root identity file remains the wipe commit point. Staging and `preCommitGate=passed` are preconditions, not commit points.

### 12.2 Inventory

| Path | Wipe treatment |
|---|---|
| `chats/` | Stage, hash, remove. |
| `f3c-commits/` | Stage, hash, remove. |
| `f3c-receipts/` | Stage, hash, remove. |
| `f3c-tombstones/` | Stage, hash, remove. |
| `f3c-creation-intents/` | Stage, hash, remove. |
| `f3c-memory-migrations/` | Stage, hash, remove. |
| `f3c-predecessors/` | Stage, hash, remove. |
| `f3c-packages/` | Stage, hash, remove. |
| `f3c-authority/` | Stage, hash, remove; new root starts empty. |
| `f3c-generations/` | Stage, hash, remove. |
| `f3c-cast-journal.json` | Stage if present, hash, remove. |
| `f3c-tmp/` | Delete; no user evidence may remain. |
| `f3c-root-identity.json` | Replace at root-identity commit; old identity is not left current. |
| `f3c-wipe-intent.json` | Durable wipe journal; preserve through phases; remove after completed new root. |
| `config.json` | Preserve user settings; reset `lastChatId` after root identity commit; never treat as chat content. |
| gallery root | Stage separately, hash, remove/restore per frozen gallery rules. |
| F3-F media | Outside F3-C; never deleted by this wipe. |

The tree manifest includes every staged file/dir, including empty directories and the exact root identity/wipe-intent state required by the frozen protocol. Symlinks/reparse points/hardlinks/unsupported objects fail closed.

### 12.3 Crash behavior

- Before `root-committing`: restore old root; no new identity.
- After new root identity: finish cleanup; do not roll back.
- Ambiguous source/staged paths: `ErrWipeBoundaryAmbiguous`; block.
- Gallery cleanup failure: `wipe_cleanup_incomplete`; no rollback claim.
- `config.json` reset failure after root commit: report incomplete cleanup; no false completion.

---

## 13. V2 exact serialization

### 13.1 Target bytes

`F3C-TARGET-2` uses the frozen `F3C-TARGET-1` JSON rules:

- `json.MarshalIndent(value, "", "  ")`, two spaces;
- default HTML escaping;
- UTF-8, LF, no BOM, no trailing newline;
- frozen Go struct declaration order;
- deterministic map keys;
- required arrays `[]`, never `null`.

`history.ndjson` uses compact JSON, one LF per entry, one trailing LF, empty history zero bytes.

### 13.2 Manifest

```json
{
  "format": "f3c-target-manifest-2",
  "targets": [
    {"path":"characters.json","presence":"present","bytes":123,"hash":"tgt2-..."},
    {"path":"scenario.json","presence":"absent"}
  ]
}
```

Exactly eight paths, sorted bytewise by path. Present entries have bytes/hash; absent entries forbid the file. The manifest is canonicalized under `F3C-CANON-1` inside records/receipts.

### 13.3 v1 coexistence

- v1 targets and `tgt1-` receipts remain valid and are never rehashed.
- A v1 chat may remain v1 indefinitely.
- The first v2 write advances one revision and writes a v2 manifest.
- Mixed v1/v2 target manifests are invalid.
- Historical verification dispatches on manifest format and never compares historical targets to current files.
- New chats, Save As, and every post-amendment write are v2.

---

## 14. Gate v2 status matrix

Gate v2 adopts **V1–V17 only**. R5 does not amend that document. V18–V24 are **proposed additions**, not adopted rows.

No row is PASS. V20 remains BLOCKED. V14/V16/V17/V24 are evidence-blocked.

| Row | Status | R5 coverage / evidence requirement |
|---|---|---|
| V1 | ADOPTED — NOT RUN | CANON-1 unchanged; `F3C-CANON-2` exact RFC 8785; domain-separated identity; golden vectors. |
| V2 | ADOPTED — NOT RUN | Receipt reconstruction from exact record body; package/record/projection crash tests. |
| V3 | ADOPTED — NOT RUN | Fork/orphan/missing inventory; record-without-target and target-without-record cases. |
| V4 | ADOPTED — NOT RUN | v1/v2 exact targets, manifest, AM-1 replacement, unknown fields. |
| V5 | ADOPTED — NOT RUN | append/autosave/session/epoch; lease and session fence. |
| V6 | ADOPTED — NOT RUN | frozen lock order, barrier, commit gate, deadlock/race trace. |
| V7 | ADOPTED — NOT RUN | Save As source binding; no receipt copy; package/migration. |
| V8 | ADOPTED — NOT RUN | tombstone/absent/root mismatch; deletion recovery. |
| V9 | ADOPTED — NOT RUN | creation journal; new chats v2. |
| V10 | ADOPTED — NOT RUN | memory migration conflict/receipt/journal agreement. |
| V11 | ADOPTED — NOT RUN | history IDs/order/unknown fields; archive. |
| V12 | ADOPTED — NOT RUN | complete wipe inventory, tree hash, gallery boundary, root identity. |
| V13 | ADOPTED — NOT RUN | fault injection at package/record/target/receipt/lease/barrier boundaries. |
| V14 | ADOPTED — BLOCKED | Native Windows symlink/reparse/hardlink/partial-staging evidence absent. |
| V15 | ADOPTED — NOT RUN | race harness for save/append/switch/delete/wipe/character batch. |
| V16 | ADOPTED — BLOCKED | Native Wails/WebView2 preference/restore/unavailable evidence absent. |
| V17 | ADOPTED — BLOCKED | Three-stage Mimo challenge not run. |
| V18 | PROPOSED — NOT ADOPTED | Exact model/execution/request provenance; all assistant paths. |
| V19 | PROPOSED — NOT ADOPTED | zero, null, omitted, fixed/random/unknown seed, legacy absence. |
| V20 | PROPOSED — BLOCKED | F3-F ownership/integrity/deletion/backup evidence absent. |
| V21 | PROPOSED — NOT ADOPTED | saved/requested/loading/ready/failed/restoring/unavailable; preference CAS. |
| V22 | PROPOSED — NOT ADOPTED | plan/attempt/candidate identity, exact commit, barrier, idempotent recovery. |
| V23 | PROPOSED — NOT ADOPTED | v1/v2 coexistence, archive, migration crash matrix, no false rollback. |
| V24 | PROPOSED — BLOCKED | Pre-amendment multi-character prompt/context/memory/cast baseline absent. |

### 14.1 Proposed row contract

Each proposed row must define claim, contract section, evidence artifact/command, oracle, positive and negative cases, fault cases, native/Mimo requirement, dependencies, and PASS/FAIL/UNRESOLVED/BLOCKED. The existing Gate v2 row-record rules apply. A proposed row does not become a pass by appearing here.

### 14.2 Adoption

Adopting V18–V24 requires an owner decision and a new Gate v2 record bound to the new authority hash. R5 does not do that.

### V18 — Proposed: model, execution, and generation provenance

- **Claim:** Every assistant-producing path records the actual model artifacts, execution snapshot, effective generation request, and seed state used.
- **Contract:** §8, §9; Master §5.2/§5.3.
- **Evidence:** artifact hashes, `F3C-MODEL-2` and `F3C-EXEC-2` vectors, generation-request hashes, narrator/character/history-path fixtures.
- **Oracle:** recompute `modelId`, `executionFingerprint`, and `generationRequestFingerprint` from recorded artifacts/request.
- **Positive:** fixed/random seed, adapters/projector/tokenizer, character and narrator paths.
- **Negative:** requested-only model, changed artifact, changed load/context control, path-only identity, missing required component.
- **Fault:** model removed/replaced/corrupt; executable unavailable; provenance partial.
- **Native/Mimo:** native model run; Mimo Stage 3 evidence review.
- **Dependencies:** V1, V19, V24.
- **PASS:** every assistant turn recomputes and an unavailable value is explicit.
- **FAIL:** wrong/fabricated identity or requested-only provenance is accepted.
- **UNRESOLVED:** artifact set or execution path not exercised.
- **BLOCKING:** model/native runtime absent.
- **Evidence identity:** artifact hashes, execution/request envelopes, command logs, output digests.

### V19 — Proposed: defaults, zero, null, omission, seeds

- **Claim:** Base/family defaults, explicit zero, null, omission, fixed/random seed, unknown effective seed, and legacy absence are distinguishable and preserved.
- **Contract:** §5.4, §6.5.1, §9.1.
- **Evidence:** golden vectors for defaults, zero, null/omitted, `defaultSetId`, fixed/random/unknown seed, legacy no generation.
- **Oracle:** exact canonical bytes and recomputed fingerprint.
- **Positive:** all valid distinctions; same effective values from explicit/inherited paths.
- **Negative:** `-1` as fixed seed, guessed seed, absent treated as zero, origin changing identity, duplicate extension key.
- **Fault:** engine cannot report effective seed.
- **Native/Mimo:** Mimo Stage 3; no native requirement.
- **Dependencies:** V1, V18.
- **PASS:** all vectors match; unknown is explicit.
- **FAIL:** conflation or invented value.
- **UNRESOLVED:** vectors not reproduced.
- **BLOCKING:** no effective request capture.
- **Evidence identity:** vector file hash, `defaultSetId`, serializer version, output.

### V20 — Proposed and blocked: F3-F media ownership

- **Claim:** F3-C attachment references never claim byte durability; F3-F owns presence, integrity, recovery, deletion, and backup.
- **Contract:** §12; frozen F3-F boundary and Q5.
- **Evidence:** F3-F ownership contract and present/corrupt/missing/delete/backup tests; F3-C reference fixtures.
- **Oracle:** reference resolves to exact bytes or reports unavailable; no substitution.
- **Positive:** present matching hash, external reference, legacy IndexedDB reference preserved.
- **Negative:** false durability, re-encoded placeholder, cross-owner deletion, bytes in target.
- **Fault:** media missing, corrupt, truncated, or backup absent.
- **Native/Mimo:** native media run; Mimo Stage 3.
- **Dependencies:** F3-F acceptance.
- **PASS:** only after F3-F evidence passes.
- **FAIL:** false durability or substitution.
- **UNRESOLVED:** F3-F incomplete.
- **BLOCKING:** F3-F ownership absent; V20 remains **BLOCKED**.
- **Evidence identity:** media hash, owner record, F3-F test evidence, reference hash.

### V21 — Proposed: continuation preference and model transition

- **Claim:** saved, requested, loading, ready, failed, restoring, active, and unavailable are distinct; no stale operation overwrites a newer explicit preference.
- **Contract:** §10; frozen ChatMeta/v2 preference rules.
- **Evidence:** native fixtures for explicit save, select-without-save, failed load, rapid changes, restore, crash.
- **Oracle:** exact `meta.continuationSeq`/`continuationModel` after each transition.
- **Positive:** select without save, explicit save, failed load, unavailable saved model, restore.
- **Negative:** load completion saves, older sequence overwrites newer, false ready, fallback, split root/meta preference.
- **Fault:** kill during transition; preference commit failure.
- **Native/Mimo:** native Wails/WebView2; Mimo Stage 3.
- **Dependencies:** V16, V18.
- **PASS:** one saved preference and truthful runtime state.
- **FAIL:** conflation or silent overwrite.
- **UNRESOLVED:** no native run.
- **BLOCKING:** native runtime absent.
- **Evidence identity:** transition trace, preference sequence/hash, native artifact.

### V22 — Proposed: plans, candidates, exact commit, recovery

- **Claim:** every batch has one plan; every attempt/candidate binds plan, slot, session, epoch, cast, user commit, base, and package; recovery is exact/idempotent.
- **Contract:** §6.9–§6.11, §7, §16.
- **Evidence:** plan/attempt/candidate fixtures; crash at plan, candidate, package, record, target, receipt, publish; deleted/removed/root-mismatch cases.
- **Oracle:** exact slot/candidate/record identity; no adoption on mismatch.
- **Positive:** committed, interrupted, stale, removed, skipped, abandoned, repeated recovery.
- **Negative:** wrong plan, slot, session, cast, candidate, receipt, or base; late write; missing lower slot.
- **Fault:** kill at each durable boundary; ambiguous recovery.
- **Native/Mimo:** native restart; Mimo Stage 3.
- **Dependencies:** V3, V6, V8, V13, V18.
- **PASS:** exact/idempotent; partial batch has defined state.
- **FAIL:** adoption, duplicate, resurrection, or permanent implicit block.
- **UNRESOLVED:** untested boundary.
- **BLOCKING:** no authority/plan evidence.
- **Evidence identity:** plan/attempt/candidate/package/record hashes, recovery trace.

### V23 — Proposed: v1/v2 coexistence and migration

- **Claim:** v1 opens unchanged; old `tgt1-` receipts remain valid; v2 migration is exact, archive-bound, and recoverable; no false rollback.
- **Contract:** §11, §13; frozen target rules.
- **Evidence:** v1/v2 fixtures, archive, package, migration crash matrix, unknown fields, missing/corrupt archive.
- **Oracle:** exact v1 bytes and v2 manifest; archive hash; record/package agreement.
- **Positive:** clean/unknown-field migration, repeated recovery, v1 current, Save As.
- **Negative:** mixed manifest, rehashed v1, missing archive, claimed rollback, new v1.
- **Fault:** every boundary in migration/package/record/receipt sequence.
- **Native/Mimo:** native migration; Mimo Stage 3.
- **Dependencies:** V4, V7, V11, V13.
- **PASS:** exact evidence and no false rollback.
- **FAIL:** loss, mixed manifest, or unprovable migration.
- **UNRESOLVED:** boundary untested.
- **BLOCKING:** archive/package not implemented.
- **Evidence identity:** archive hash, v1/v2 hashes, package/record, migration record.

### V24 — Proposed and blocked: multi-character baseline parity

- **Claim:** deterministic participant ordering and `contextMode` preserve the captured pre-amendment character prompt/context/memory/cast behavior; no cross-character leak.
- **Contract:** §7.7, §7.2–§7.4.
- **Evidence:** pre-amendment baseline; two/three-character fixtures; captured and prefix modes; removal/cast change; native run.
- **Oracle:** exact baseline comparison of prompts/context/memory/cast/history order.
- **Positive:** concurrent characters, removed character, retry/restart, partial batch.
- **Negative:** scheduler changes context, cross-character leak, reordered response, stale commit, changed prompt.
- **Fault:** remove character during generation; delete during batch; crash.
- **Native/Mimo:** native multi-character; Mimo Stage 3.
- **Dependencies:** V15, V16, V18.
- **PASS:** baseline parity and deterministic ordering.
- **FAIL:** behavior change without owner-approved contract.
- **UNRESOLVED:** baseline absent.
- **BLOCKING:** no pre-amendment baseline/native run; V24 remains blocked.
- **Evidence identity:** baseline manifest/hash, run IDs, comparison hash.

---

## 15. Complete prior-finding disposition

Statuses: **Resolved by R5 rule**, **Partially resolved**, **Unresolved**, **Contradicted**, **Owner decision**, **Superseded**. Resolution is contract-level, not verification.

### 15.1 R1 review R1–R10

| ID | Disposition | R5 rule / residual |
|---|---|---|
| R1 | Resolved by R5 rule | `f3c-generations` is root/chat-owned, tombstone/wipe inventory, and subordinate to authority. |
| R2 | Resolved by R5 rule | Plan/attempt/candidate bind user commit, batch, attempt, candidate, receipt, session, epoch, cast. |
| R3 | Resolved by R5 rule | Recovery precedence in §16; ambiguous combinations block. |
| R4 | Resolved by R5 rule | Runtime transition state machine and terminal readiness. |
| R5 | Resolved by R5 rule | Full effective request, model/execution identity, attempt-start capture. |
| R6 | Resolved by R5 rule | Archive, strict v2, exact predecessor evidence. |
| R7 | Partially resolved | F3-F boundary and V20 block retained; F3-F ownership/integrity evidence is absent. |
| R8 | Resolved by R5 rule | Explicit non-goals; V24 baseline required but absent. |
| R9 | Resolved by R5 rule | Complete v2 field order/presence/null rules. |
| R10 | Resolved by R5 rule | All assistant paths use the attempt protocol. |

### 15.2 R2 review N1–N15

| ID | Disposition | R5 rule / residual |
|---|---|---|
| N1 | Resolved by R5 rule | Absent without tombstone is `ErrChatStateUnknown`; no cleanup. |
| N2 | Resolved by R5 rule | Exact candidate/attempt/record/package identity. |
| N3 | Resolved by R5 rule | Lease barrier plus frozen lock integration. |
| N4 | Resolved by R5 rule | Migration uses the ordinary package/record order. |
| N5 | Resolved by R5 rule | `F3C-MODEL-2` artifact envelope. |
| N6 | Resolved by R5 rule | Complete effective parameter schema and default set. |
| N7 | Resolved by R5 rule | Exact v1 archive + strict v2. |
| N8 | Resolved by R5 rule | Immutable predecessor archive; no automatic in-place rollback. |
| N9 | Resolved by R5 rule | Explicit preference operation; no load-time save. |
| N10 | Resolved by R5 rule | Complete V1–V24 status matrix, with V18–V24 proposed. |
| N11 | Resolved by R5 rule | Per-scope execution and serialized commit coordinator. |
| N12 | Resolved by R5 rule | Root mismatch blocks/preserves. |
| N13 | Resolved by R5 rule | `F3C-CAND-3` without time; package binding. |
| N14 | Resolved by R5 rule | New chats/Save As/post-amendment writes are v2. |
| N15 | Partially resolved | Baseline required, but no pre-amendment baseline exists; V24 blocked. |

### 15.3 R3 review M1–M10

| ID | Disposition | R5 rule / residual |
|---|---|---|
| M1 | Resolved by R5 rule | Every durable path is lease-bound; barrier drains leases/gates. |
| M2 | Resolved by R5 rule | Plan, prefix order, exact parent, candidate authority. |
| M3 | Resolved by R5 rule | Model artifacts and execution split. |
| M4 | Resolved by R5 rule | CANON-1/CANON-2 split; exact RFC 8785; no time in identity. |
| M5 | Resolved by R5 rule | Archive/package/record/receipt recovery; no in-place rollback. |
| M6 | Resolved by R5 rule | Complete candidate envelope and package binding. |
| M7 | Resolved by R5 rule | Durable `blocked_timeout`; no release. |
| M8 | Resolved by R5 rule | `stateSessionId`, `stateEpoch`, `castEpoch`, character version. |
| M9 | Partially resolved | Gate rows are specified, but V18–V24 are not adopted and no evidence exists. |
| M10 | Resolved by R5 rule | `meta` authority, `continuationSeq` CAS, no load save. |

### 15.4 Independent R4 review A1–A14

| ID | Disposition | R5 rule / residual |
|---|---|---|
| A1 | Resolved by R5 rule | Durable commit package, record references package, deterministic recovery, no GC. |
| A2 | Resolved by R5 rule | One frozen lock order; lease/commit gate mapped onto `stateApplyMu`/`chatWriteMu`; no parallel authority lock. |
| A3 | Resolved by R5 rule | Lease token, expiry, heartbeat, max lifetime, session/epoch/cast, startup reconciliation. |
| A4 | Resolved by R5 rule | Candidate/commit include `stateSessionId`, `stateEpoch`, `castEpoch`, character version. |
| A5 | Resolved by R5 rule | Durable plan header, batch/slot state, attempt/candidate records, complete statuses, explicit abandonment, partial-batch state. |
| A6 | Resolved by R5 rule | Exact RFC 8785 for CANON-2; no time in identity; origin/extension rules. |
| A7 | Resolved by R5 rule | Frozen `requestFingerprint` preserved; separate `generationRequestFingerprint`. |
| A8 | Resolved by R5 rule | `meta.continuationModel` sole authority; `continuationSeq` CAS; no root preference authority. |
| A9 | Resolved by R5 rule | One package/record order; archive domains; no in-place rollback; import Q3 retained. |
| A10 | Resolved by R5 rule | Complete wipe inventory including root identity, wipe intent, tmp, authority, packages, generations. |
| A11 | Resolved by R5 rule | Complete v2 target/record/receipt/package schemas; frozen fields retained. |
| A12 | Resolved by R5 rule | Role qualifier, formats, load/context/tokenizer controls, execution snapshot. |
| A13 | Resolved by R5 rule | Gate v2 V1–V17 adopted; V18–V24 marked proposed. |
| A14 | Partially resolved | `contextMode = captured` is the safe interim rule preserving prompt construction; `prefix_committed` requires the OD-15 owner choice and the V24 baseline. |

No prior finding is marked verified. The remaining partials are evidence or owner decisions, not omitted rules.

---

## 16. Recovery precedence

1. Root identity and incarnation.
2. Wipe intent/barrier.
3. Chat barrier state.
4. Tombstone inventory.
5. Commit package/record/receipt inventory.
6. Ordinary save/append journal.
7. Cast journal and cast epoch.
8. Response plans and generation attempts.
9. Candidate staging.
10. Terminal cleanup.

Rules:

- Root mismatch → block, never clean.
- Wipe intent → no chat recovery; complete or block wipe.
- Tombstone present → no recreation; clean only after sealed barrier.
- A commit record is authority only with its verified package; a package without a record is orphaned.
- Ordinary/cast/attempt nonterminal for the same chat → `ErrRecoveryAmbiguous`; no silent choice.
- A removed character/cast change terminalizes affected candidates and updates affected `batch.json` slots.
- Recovery holds a recovery lease and obeys the same locks/fence.
- Repeated recovery is idempotent; no duplicate turns, records, receipts, packages, or revisions.

---

## 17. A1–A14 acceptance criteria

| Finding | Required positive | Required negative | Recovery/evidence |
|---|---|---|---|
| A1 | Valid record has package, targets, receipt; crash at each step recovers. | Record without package/target blocks; no guessed bytes. | Package retention; fault matrix V2/V13. |
| A2 | Every path follows one frozen lock order; no deadlock. | Reverse order, blocked commit gate, lock inversion fails. | Native lock trace; V6. |
| A3 | Active lease only commits before expiry with token/session/epoch/cast match. | Expired/orphaned/revoked/old-incarnation cannot commit. | Startup/barrier crash tests; V5/V6/V15. |
| A4 | Removal/cast/session change invalidates candidate. | Same-name recreation/stale version cannot commit. | V8/V15/V22/V24. |
| A5 | Every plan/slot reaches a legal terminal/blocked state; partial batch visible. | Missing lower slot cannot be skipped implicitly; no duplicate. | V13/V15/V22. |
| A6 | Identical content yields identical bytes; explicit/inherited same value same fingerprint. | Time, code-point order, extension duplicate, or origin changes fingerprint incorrectly. | Golden vectors V1/V19. |
| A7 | Operation and generation fingerprints independently verifiable. | One field cannot satisfy both; v1 remains valid. | V1/V2/V18. |
| A8 | One saved preference; CAS prevents stale overwrite; restart restores it. | Autosave/load cannot overwrite explicit preference; false active fails. | V16/V21. |
| A9 | Migration uses package/record order; archive binds v1; no mixed manifest. | Target-before-record, missing archive, false rollback fail. | V4/V13/V23. |
| A10 | Every owned root staged/handled; no old identity/temp survives. | Omitted root, config loss, late writer, root mismatch fail. | V12/V14. |
| A11 | Every record has complete schema, hash domain, location, validation. | Missing field/duplicate identity/undefined path fails. | V1–V4/V22/V23. |
| A12 | Historical turn identifies artifacts and actual execution. | Same artifacts/different controls collide or path-only identity fails. | V18/V19. |
| A13 | V1–V17 adopted status preserved; V18–V24 proposed. | Any claim of adoption/PASS without owner action fails. | Gate record. |
| A14 | Same plan/input yields deterministic context/order; baseline parity. | Scheduler changes prompt/order silently; cross-character leak fails. | V24 baseline/native. |

---

## 18. Owner decision register

Frozen requirements are marked **FROZEN** and are not optional preferences.

| ID | Question | Options / consequence | Recommendation | Frozen? | Status |
|---|---|---|---|---|---|
| OD-01 | Retain every commit package until deletion/wipe, or compact? | Retain = safe recovery, more disk. Compact = requires separate proof and GC policy. | Retain until deletion/wipe; no automatic GC. | No | Proposed |
| OD-02 | Accept **AM-1** target projection replacement? | Accept = revisions can replace current targets atomically from package. Reject = only new-chat v2 writes; updates blocked. | Accept AM-1 with package pre-image. | **Yes — formal amendment** | Awaiting owner |
| OD-03 | Keep one frozen lock order with `chatWriteMu` as commit gate? | Yes = no parallel authority. No = deadlock/second authority; not recommended. | Yes. | **Yes** | Proposed |
| OD-04 | Lease TTL/renewal values | 120s/30s, max 10m; shorter = more false expiry; longer = slower recovery. | 120s/30s, max 10m. | No | Proposed |
| OD-05 | Candidate/session/cast fields | Include session, epoch, cast, character version. Exclude = stale commit risk. | Include all. | No | Proposed |
| OD-06 | Response plan and batch-state storage in `f3c-generations` | Separate root = clear evidence; extra wipe root. Embed in target = violates eight-target set. | Separate immutable plan header plus mutable batch/slot state in `f3c-generations`. | No | Proposed |
| OD-07 | Canonicalization versions | CANON-1 targets/receipts; exact RFC 8785 CANON-2 identity. Alternative = new target canonicalization (not recommended). | Split as specified. | Partly frozen | Proposed |
| OD-08 | Fingerprint names/domains | Operation vs generation vs candidate distinct. | Separate names/domains. | No | Proposed |
| OD-09 | Preference authority | `meta` only with CAS; root cache rebuildable. Root authoritative = split-brain. | `meta` only. | No | Proposed |
| OD-10 | Migration uses normal package/record protocol | Same = one recovery model. Different = two proofs. | Same. | No | Proposed |
| OD-11 | Wipe config/root behavior | Preserve settings, reset active chat; stage all other roots. Delete settings = product change. | Preserve settings/reset pointer. | Partly frozen | Proposed |
| OD-12 | Exact v2 schemas | Adopt §6 schemas; no implementation-defined fields. | Adopt. | No | Proposed |
| OD-13 | Model/execution identity | Model-2/Exec-2 complete; unavailable explicit. | Adopt. | No | Proposed |
| OD-14 | Adopt V18–V24? | Adopt via new Gate record; or keep proposed. | Keep proposed until owner/Gate action. | **Gate v2** | Awaiting owner |
| OD-15 | Multi-character context mode | `captured` preserves existing prompt builder; `prefix_committed` may serialize. | `captured`; V24 baseline decides if change needed. | No | Proposed |
| OD-16 | Import Q3 scope | Preserve frozen journal/reversibility; do not equate with Save As. | Preserve Q3; import remains separate. | **FROZEN Q3** | Proposed |

---

## 19. Unresolved risks and evidence still required

- **AM-1** is an open formal amendment. R5 MUST NOT be frozen or claimed compliant with frozen §6.9 target replacement until OD-02 is decided.
- `contextMode = captured` is the safe default, but no pre-amendment V24 baseline exists. The owner must approve the default or choose `prefix_committed`; V24 remains blocked.
- No native Wails/WebView2, Windows reparse, model-load, crash, migration, media, or destructive evidence exists.
- No Gate row passes. V14/V16/V17/V24 are evidence-blocked; V20 is F3-F-blocked.
- F3-F media ownership remains absent; F3-C does not claim attachment durability.
- Import Q3 remains a separate frozen requirement; R5 does not prove its implementation.
- The current source still has direct `DeleteChat`, multi-file `SaveChat`, async `performSave`, return-before-ready `LoadModel`, and transient image handling. These are implementation gaps, not evidence for R5.
- Package retention consumes storage; OD-01 is unresolved. No automatic GC is claimed.
- Lease timings are proposed defaults; OD-04 must be accepted or changed before implementation.
- The authority/lease records are proposed operation evidence; no runtime evidence exists.

---

## 20. Status and authorization

**R5:** **PROPOSAL ONLY; NOT FROZEN; NOT ACCEPTED; NOT IMPLEMENTED; NOT VERIFIED.**  
**R5 SHA-256:** recorded in the completion report, not embedded here to avoid a self-referential file hash.  
**Frozen authority:** unchanged; hash verified.  
**Gate v2:** adopted V1–V17 only; V18–V24 proposed.  
**F3-C phase 0:** **NOT CLOSED.**  
**P1a:** **UNAUTHORIZED**, downstream of P1c.  
**F3-C acceptance:** **NOT GRANTED.**  
**F3-D:** **BLOCKED.**  
**F3-B:** closed; not reopened.  
**A9/A10/A12/A13:** closed; not reopened.  
**V20:** **BLOCKED** pending F3-F evidence.

No implementation authorization exists. No source, test, build, migration, model load, destructive test, application data, or authority record was changed.

R5 is ready for independent adversarial review as a documentation proposal. It is **not ready for freeze** while OD-02 (AM-1), OD-14 (Gate additions), OD-15 (multi-character context), and the missing native/Mimo/F3-F/baseline evidence remain unresolved.
