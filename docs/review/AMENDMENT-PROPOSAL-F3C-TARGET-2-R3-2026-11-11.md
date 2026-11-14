# AMENDMENT PROPOSAL — F3C-TARGET-2, Revision 3

**Status:** **PROPOSAL ONLY — NOT FROZEN AUTHORITY, NOT ACCEPTANCE.**  
**Date:** 2026-11-11  
**Supersedes for review:** R2 (`docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-R2-2026-11-11.md`). R1/R2 remain historical.  
**Fresh review addressed:** `docs/review/REVIEW-ADVERSARIAL-F3C-TARGET-2-R2-2026-11-11.md`  
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`  
**Scope:** contract-resolution proposal only. No implementation, tests, builds, model loads, migrations, destructive tests, source changes, or freeze.

R3 resolves all fifteen fresh-review findings. It makes the safety-critical decisions explicit rather than leaving them to implementation interpretation. It remains a proposal; independent adversarial review is still required.

---

## 1. Resolution summary

| Finding | R3 resolution |
|---|---|
| N1 absent chat | Absent without tombstone is `ErrChatStateUnknown`; only an authoritative tombstone permits cleanup. |
| N2 commit identity | Recovery verifies exact attempt/message/commit/candidate/session/base/receipt identity; a merely existing receipt is insufficient. |
| N3 delete/wipe race | A destructive-operation generation barrier rejects new attempts, cancels active attempts, waits for confirmed producer termination, and blocks on timeout. |
| N4 v2 migration order | A predecessor archive and a staged, journaled migration define install/record/receipt order and crash recovery. |
| N5 model identity | `modelId` is a concrete digest of model + projector bytes; path/display are hints. |
| N6 parameter completeness | Add `parameterSchemaVersion`, a complete canonical `requestFingerprint`, and an exact effective parameter set. |
| N7 unknown fields | Exact v1 target bytes are archived before v2; unknown fields are preserved losslessly; target v2 remains strict. |
| N8 predecessor bytes | V1 predecessor bytes are retained immutably in `f3c-predecessors` with a receipt/record hash; rollback/reconstruction is possible. |
| N9 preference semantics | Separate explicit selection and explicit preference operations; generic loading never rewrites the saved preference. |
| N10 Gate mapping | Full V1–V24 mapping, including previously omitted rows. |
| N11 multi-character concurrency | Exclusivity is per execution scope `(chatID, characterId, stateSessionID)`, not the whole chat; different characters may run concurrently. |
| N12 root mismatch | Root-incarnation mismatch blocks with `ErrRootIncarnationChanged`; no evidence cleanup. |
| N13 fingerprint semantics | `F3C-CAND-1` candidate fingerprint and exact commit-operation allocation/retry rules are defined. |
| N14 new-chat format | All new chats, Save As, and post-amendment writes are v2; v1 is legacy-only. |
| N15 V24 baseline | V24 requires a captured pre-amendment multi-character baseline and native evidence. |

---

## 2. Generation attempt lifecycle

### 2.1 Root, chat, and deletion authority

Generation intents live under a root-owned operation area:

```text
f3c-generation-intents/<F3C-PATH-GEN-1 hash>/<operationID>.json
```

They are not targets and not receipts. They are subordinate to root identity and chat deletion authority.

**Absent-chat rule:**

| State | Result |
|---|---|
| Tombstone present, chat directory present or absent | Deletion authoritative; clean intents; complete deletion; no recreation. |
| No tombstone, chat directory absent | `ErrChatStateUnknown`; no cleanup; no target; preserve intent. |
| No tombstone, chat directory present | Live chat; normal recovery. |
| Root identity/incarnation mismatch | `ErrRootIncarnationChanged`; block; no cleanup. |
| Wipe intent present | Wipe precedence; no generation recovery. |

An absent directory is not deletion evidence. Only a valid tombstone authorizes cleanup.

`DeleteChat` and wipe must also acquire the destructive generation barrier in §2.6 before tombstone/staging.

### 2.2 Intent fields

```text
version, rootId, rootIncarnation, chatId, operationId,
operationType, executionScope, stateSessionId, baseEpoch,
baseRevision, userMessageId, userCommitOperationId,
userMessageHash, assistantMessageId,
assistantCommitOperationId, assistantCandidateFingerprint,
modelId, startedAt, status, terminalReason
```

Rules:

- `executionScope` is `(chatId, characterId, stateSessionId)`; the generic user turn uses `characterId = ""`.
- `userCommitOperationId` must identify the durable user-turn commit.
- `assistantMessageId` is the candidate message ID.
- `assistantCommitOperationId` is allocated before the assistant commit attempt.
- `assistantCandidateFingerprint` is `F3C-CAND-1`, defined in §2.5.
- `status` is `started`, `streaming`, `completion_pending`, `committed`, `interrupted`, `aborted`, `failed`, or `stale`.
- No generation starts without a durable user-turn commit.

### 2.3 Attempt concurrency

R3 replaces “one active attempt per chat” with **one active attempt per execution scope**:

- `(chatID, "", stateSessionID)` for a generic user turn;
- `(chatID, characterId, stateSessionID)` for each character response.

Different characters may generate concurrently. Two attempts for the same character/session are serialized. The destructive barrier in §2.6 is chat-wide or root-wide, not scope-wide.

This preserves legitimate multi-character operation while preventing duplicate attempts for the same character/session.

### 2.4 Exact commit verification

The intent is adopted only when all of the following match:

- root ID and root incarnation;
- chat ID and execution scope;
- `stateSessionId`, `baseEpoch`, `baseRevision`;
- `assistantCommitOperationId`;
- `assistantMessageId`;
- `assistantCandidateFingerprint`;
- the commit record and receipt body;
- the receipt envelope and target manifest;
- no conflicting commit for the same attempt.

A receipt that merely exists at the same chat/chat revision is **not** evidence. A mismatch is `ErrGenerationCommitConflict`; no adoption, no cleanup as committed.

If no assistant commit exists:

- base current, no tombstone → `interrupted`;
- base stale → `stale`;
- absent chat without tombstone → `ErrChatStateUnknown`;
- tombstone → cleanup only;
- root mismatch → `ErrRootIncarnationChanged`.

### 2.5 `F3C-CAND-1` candidate fingerprint

`assistantCandidateFingerprint` is:

```text
sha256:<64 lowercase hex>
```

It is SHA-256 over the canonical `F3C-CAND-1` candidate envelope containing:

- chat ID, operation ID, base revision;
- assistant message ID;
- ordered target manifest for the candidate;
- canonical candidate history/state bytes;
- root/chat identity fields.

Prompt text and secrets are excluded. The same candidate must produce the same fingerprint on every implementation.

`assistantCommitOperationId` is allocated once for the candidate. A failed commit may be retried only with the same operation ID and same fingerprint while the intent is `completion_pending`; a different candidate requires a new intent/operation. A completed operation is never reused.

### 2.6 Destructive-operation generation barrier

Deletion and wipe must establish a generation barrier before removing or staging content.

Barrier phases:

1. **acquire**: take the chat/root write authority;
2. **close**: set the barrier epoch so no new generation attempt can start;
3. **cancel**: request cancellation of every active attempt in scope;
4. **await**: wait for every producer to confirm termination and every intent to reach terminal or recoverable state;
5. **reconcile**: adopt a commit that completed before the barrier; mark later output unsaved;
6. **destructive step**: tombstone/delete or stage/wipe;
7. **release**: clear the barrier only after the destructive step has a durable outcome.

Rules:

- Cancellation **requested** is not cancellation **completed**. The barrier waits for producer acknowledgement/terminal intent.
- If an attempt reaches a commit before the barrier, reconcile it.
- If an attempt attempts a write after the barrier epoch, the commit path rejects it; the output is unsaved, not completed.
- Timeout, uncertain termination, or an unresolved intent blocks deletion/wipe with `ErrGenerationBarrierTimeout`; no tombstone/staging occurs.
- A late stream completion must not create a target, directory, receipt, or assistant turn.
- Wipe uses a root-wide barrier; chat deletion uses a chat-wide barrier.
- Barrier state is runtime; durability comes from the tombstone/wipe intent. A crash ends producers; restart never resumes a pre-crash stream.

### 2.7 Recovery precedence

1. Root identity/incarnation.
2. Wipe intent.
3. Tombstone inventory.
4. Ordinary create/save/append journal.
5. Cast journal.
6. Generation intents.
7. Terminal cleanup.

Rules:

- Root mismatch → block, never clean.
- Wipe intent → no chat recovery.
- Tombstone present → clean intents only; no recreation.
- No tombstone and absent chat → `ErrChatStateUnknown`; preserve evidence.
- Ordinary/cast journal nonterminal with a generation intent for the same chat → `ErrGenerationJournalAmbiguous`; no silent choice.
- Generation recovery never creates a chat revision.
- Recovery is idempotent: repeated runs produce the same terminal intent status and no duplicate user/assistant turns or receipts.

### 2.8 All assistant-producing paths

The same protocol applies to `SendMessage`, `SendMessageStreamWithImage`/`SendMessageStream`, and `GenerateCharacterResponse`. Character generation uses its existing character-specific prompt/context behavior; only persistence and attempt identity are added.

---

## 3. Hot-swap, requested vs active vs saved

### 3.1 Explicit concepts

| Concept | Durability | Owner |
|---|---|---|
| `savedContinuation` | Durable chat preference | `meta.continuationModel` |
| `requestedModel` | Ephemeral | UI/runtime request |
| `loadingModel` | Ephemeral | Loader |
| `activeModel` | Ephemeral, healthy | Runtime |
| `failedModel` | Ephemeral error state | Runtime |
| `restoringModel` | Ephemeral | Loader |

Only `activeModel` may generate. Only an explicit preference operation may change `savedContinuation`.

### 3.2 Preference operations

- **`SelectContinuationModel(modelId)`** — explicit user selection; load; saved preference changes only after `ready`.
- **`SetContinuationPreference(modelId)`** — explicit preference-only operation; may save an unavailable model; does not load or claim active.
- **Generic `LoadModel`/auto-load** — changes active/requested only; never rewrites `savedContinuation`.

A failed load never rewrites the saved preference. If A is saved/active and B is requested and fails, saved remains A; if A is restored, active is A; if not, active is `unavailable`. If B is ready after an explicit selection, saved may become B. No automatic fallback.

When saved ≠ active, sending must load the saved preference. If it is unavailable, generation fails; it must not silently use another model.

### 3.3 Transition state machine

`idle` → `preflight` → `wait_for_stream` (if needed) → `capture_previous` → `stopping` → `starting` → `health_check` → `publishing` → `ready`; failures go to `rollback` then `ready`/`unavailable`.

Rules:

- One transition at a time.
- `LoadModel` returns only after `ready` or a terminal failure. There is no “returned without error but not ready” state.
- UI can display `requested`, `loading`, `ready`, `failed`, `restoring`, `unavailable` independently.
- Active stream policy is **wait**; a transition never terminates an active stream by itself. The user may explicitly abort.
- `capture_previous` stores model path/identity, projector, effective spawn settings, admission/eligibility evidence, and resource estimate.
- Stop failure: do not start replacement; health-check old.
- Spawn/health failure: kill failed process; one restore with captured snapshot; otherwise `unavailable`.
- Restore failure: `unavailable`; saved preference retained; no fallback.
- Crash: no durable active model; saved preference remains; no auto-load; no false ready.
- Preference commit failure after ready: active model is ready but saved preference is not durable; report both.

`App.serverProc`, HTTP client, `resource.Manager`, and `model.Manager` must be reconciled by one transition owner before `ready`/`unavailable`.

---

## 4. Per-turn provenance

### 4.1 Capture point

Capture the effective configuration after model defaults, clamping, template injection, and backend normalization, immediately before the first token. Do not reconstruct from current settings at save time.

### 4.2 Model identity

`modelId` is:

```text
mdl_<64 lowercase hex>
```

SHA-256 over canonical JSON:

```json
{"model_sha256":"<64 hex>","projector_sha256":"<64 hex or empty>"}
```

- `model_sha256` is the digest of the exact model artifact bytes.
- `projector_sha256` is the digest of the paired projector, or empty if none.
- A moved file keeps identity.
- A changed quant/content/projector changes identity.
- If either required digest cannot be computed, generation does not start; no placeholder identity.
- Registry IDs, paths, and display names are aliases/hints, never `modelId`.

### 4.3 Provenance object

For a v2 assistant turn:

```text
modelId, displayName, engine, backend, parameterSchemaVersion,
parameters, requestFingerprint, seedMode, seed, effectiveSeed,
templateFingerprint
```

Rules:

- `modelId`, `engine`, `backend`, `parameterSchemaVersion`, `parameters`, `requestFingerprint`, `seedMode`, `templateFingerprint` are required.
- `parameterSchemaVersion` is `f3c-genparams-1`.
- `parameters` contains the exact required effective controls:
  `temperature`, `topP`, `topK`, `repeatPenalty`, `presencePenalty`, `frequencyPenalty`, `maxTokens`, `contextTokens`, `stopSequences`, `samplerOrder`.
  No unknown keys. Each numeric value is finite; arrays preserve order.
- `requestFingerprint` is `sha256:<64 hex>` over the canonical effective inference request, excluding prompt text and secrets, including all backend controls, modelId, parameters, seed, and template. If a backend uses an unlisted control, the fingerprint still changes, and a schema-version bump is required to expose it.
- `seedMode` is `explicit` or `random`; `seed` is present only for explicit, signed 64-bit, zero valid; `effectiveSeed` is optional for known random.
- `templateFingerprint` is `sha256:<64 hex>` or `none`.
- `displayName` is optional/display-only.

Absence of `generation` is legacy-unknown. It is never fabricated.

### 4.4 Retries and partials

- A partial attempt has an intent, not a committed assistant turn.
- No partial response becomes a completed entry.
- A retry is a new operation, new candidate, new fingerprint, and new generation object.
- The prior attempt remains terminal evidence.
- Mid-stream model selection does not affect the attempt.

### 4.5 All paths

All assistant-producing paths use the same provenance helper. `characterId` remains the owner field; no character prompt/context semantics are changed.

---

## 5. V1-to-V2 migration and predecessor retention

### 5.1 New chats

All new chats, `Save As`, and every post-amendment write use `F3C-TARGET-2`/`f3c-target-manifest-2`/`tgt2-`. V1 is legacy-only.

### 5.2 Predecessor archive

Before the first v2 write, retain the exact v1 target bytes immutably:

```text
f3c-predecessors/<F3C-PATH-PRED-1 hash>/<chatID>/<v1Revision>/
```

The archive contains exact copied v1 target files plus the v1 manifest/receipt evidence. It is not a target, not a read path, and never a second source. Its manifest hash is `predecessorArchiveHash`.

The commit record/receipt gains an optional `predecessorArchiveHash` field for the migration. V1 receipts are unchanged; old receipts remain valid.

### 5.3 Unknown fields

- Exact v1 target bytes, including unknown fields, are preserved in the predecessor archive.
- V2 target objects remain strict: unknown fields in v2 are `ErrTargetMalformed`.
- A v1 chat with unknown fields can upgrade because the unknown data is proven preserved in the archive; it is not silently dropped.
- If the predecessor archive cannot be written/verified, the v2 write blocks with `ErrPredecessorArchiveFailed`; no v1 loss.
- `Save As` from v1 also archives before writing v2.

This resolves the deadlock and the unknown-field preservation obligation without adding arbitrary fields to the v2 target schema.

### 5.4 Migration transaction

1. Verify v1 current, receipt, chain, and target manifest.
2. Write `f3c-predecessors/...` create-only; read back and verify exact v1 bytes; compute `predecessorArchiveHash`.
3. Stage v2 targets; read back; verify.
4. Under `chatWriteMu`, install v2 targets create-only/atomic; verify installed bytes.
5. Write commit record with v1 predecessor, v2 manifest, and `predecessorArchiveHash`.
6. Write receipt.
7. Publish v2.
8. Retain the predecessor archive until chat deletion/wipe or a separate explicit retention policy.

Recovery:

| Crash point | Recovery |
|---|---|
| Before predecessor archive verified | v1 current; remove temps; no v2. |
| Archive verified, before v2 install | v1 current; no v2; archive may remain as evidence. |
| V2 targets installed, before record | Restore v1 from archive; verify; block if restore fails. |
| Record written, before receipt | Reconstruct receipt from record; verify v2 targets. |
| Receipt written, before publish | Publish/verify v2. |
| After publish | Normal current/commit verification. |

No mixed v1/v2 manifest is valid. `predecessorArchiveHash` must match the archive and receipt/record.

### 5.5 Predecessor retention and wipe

- Predecessor archives are immutable evidence; they are not edited.
- Chat deletion removes its predecessor archives only after the tombstone is authoritative.
- Wipe stages `f3c-predecessors/` with the other operation roots.
- A predecessor archive is not a substitute for the current v2 target; it is rollback/evidence.

### 5.6 Legacy behavior

- v1 receipt/record/manifest chain remains valid.
- v1 current chats open unchanged.
- Migrated entries get `attachments: []`, no `generation`; `hasImage` remains a legacy-unknown projection.
- No global last-model path is written into `meta`.
- IDs are not renumbered.
- Old receipts are never rehashed.

---

## 6. Media and F3-F

- F3-C core work may proceed without F3-F.
- F3-C full acceptance and V20 remain blocked until F3-F proves byte ownership, integrity, deletion, and backup.
- `attachments` refs use `mediaId`, `role`, `provider`, `mimeType`, `byteSize`, `contentHash`.
- `f3f-media` refs require `contentHash`; `legacy-indexeddb` may omit it.
- Missing/corrupt/unavailable media is reported, never substituted.
- F3-C never deletes F3-F bytes; F3-F owns lifecycle.
- No base64, local path, or thumbnails in targets.
- `F3-C-CORE` excludes attachment durability; `F3-C-FULL` includes it.

---

## 7. Multi-character isolation and scopes

No changes to:

- character-specific context flushing;
- character memory isolation;
- character-specific prompt construction;
- cast synchronization semantics;
- multi-character interaction behavior.

Persistence rules:

- `characterId` remains the message owner.
- `generation` is per assistant turn.
- `attachments` is per message.
- `continuationModel` is per chat, not per character.
- Attempt exclusivity is per `(chatID, characterId, stateSessionID)`.
- Different characters may generate concurrently; the destructive barrier is chat/root-wide.
- No intent/recovery path may alter cast state.

V24 requires a captured pre-amendment baseline for character context, memory, prompts, and cast behavior, plus native evidence where applicable. It is a regression requirement, not a new cast contract.

---

## 8. V2 serialization

### 8.1 `meta.json`

```text
id, name, createdAt, updatedAt, worldName, scenarioName,
characterNames, lorebookNames, worldId, scenarioId,
activeCharacters, removedCharacters, addedCharacters,
activeLorebooks, removedLorebooks, addedLorebooks,
rootId, revision, lastCommitOperationId, continuationModel
```

`continuationModel` is omitted when absent. When present: `modelId, displayName, sourceHint, selectedAt`; `modelId` and `selectedAt` required; display/source optional.

### 8.2 History

```text
id, sender, text, role, timestamp, hasImage, characterId,
generation, attachments
```

`generation` absent for non-assistant/legacy; present for v2 assistant. `attachments` always `[]` for v2, never null. `hasImage` is a compatibility projection; true with `[]` is legacy-unknown.

### 8.3 `generation` order

```text
modelId, displayName, engine, backend, parameterSchemaVersion,
parameters, requestFingerprint, seedMode, seed, effectiveSeed,
templateFingerprint
```

Rules as §4.3. Unknown keys are malformed.

### 8.4 Attachment order

```text
mediaId, role, provider, mimeType, byteSize, contentHash
```

Required as in §6. No bytes or paths.

### 8.5 Exact bytes and manifest

- v2 JSON/NDJSON formatting unchanged from v1.
- `f3c-target-manifest-2`, `tgt2-`.
- Exactly eight targets; no ninth.
- Receipt/record/envelope shapes unchanged except optional `predecessorArchiveHash` on migration.
- `F3C-CANON-1` unchanged.
- Unknown v2 fields fail closed.
- `null` is not an omitted optional field.
- `attachments` never omitted for v2; `generation` omitted only for legacy/non-assistant.

---

## 9. Gate v2 V1–V24 mapping

The Gate v2 document is not amended or frozen here. This is the required mapping for a later freeze.

| Row | Claim | Required negative / evidence |
|---|---|---|
| V1 | Canonicalization, fingerprints, Base64, hash chain. | Mismatch fails; unchanged from Gate v2. |
| V2 | Receipt inventory/reconstruction. | Missing/malformed receipt blocks; unchanged. |
| V3 | Commit fork/orphan/missing inventory. | Fork/orphan blocks; unchanged. |
| V4 | v1/v2 bytes, manifests, `tgt1`/`tgt2`, predecessor archive hash, new-chat v2. | Mixed manifest/unknown fields/omission fail. |
| V5 | User turn durable before attempt; stream-end not commit. | Late append/complete-before-commit fails. |
| V6 | Locks: intent/barrier/delete/wipe/order. | Deadlock or late write fails. |
| V7 | Save As v1→v2 archives predecessor and verifies source. | No receipt copy; archive failure blocks. |
| V8 | Tombstone vs absent vs root mismatch. | Absent without tombstone -> `ErrChatStateUnknown`; mismatch blocks. |
| V9 | Creation journal and new-chat v2. | Creation never v1 post-amendment. |
| V10 | Memory migration. | Unchanged; conflict/receipt agreement. |
| V11 | History migration IDs/order/unknown fields archived. | Duplicate/empty ID/unknown-loss fails. |
| V12 | Wipe stages chats, commits, receipts, tombstones, intents, predecessors, gallery. | Barrier/staging/intent failure blocks. |
| V13 | Fault injection at barrier, intent, migration, install, record, receipt, restore. | No ambiguous adoption or duplicate. |
| V14 | Filesystem symlink/reparse/hardlink/partial staging. | Escape blocks; unchanged. |
| V15 | Races: stream/switch/delete/wipe/multi-character scopes. | Late write/duplicate attempt fails. |
| V16 | Native Wails/WebView2 preference/restore/unavailable. | No false ready; no substitution. |
| V17 | Independent Mimo Flash Fast challenge. | Required before approval; not omitted. |
| V18 | Exact effective provenance, identity, request fingerprint, all paths. | Requested-only/provenance rewrite fails. |
| V19 | Explicit zero, random known/unknown, legacy absence. | `-1`/guessed seed fails. |
| V20 | F3-F media ownership, present/corrupt/missing. | **BLOCKED** until F3-F; reference alone fails. |
| V21 | Saved/requested/loading/ready/failed/restoring/unavailable; preference rules. | Failed load rewrites preference or false ready fails. |
| V22 | Intent identity, exact commit, absent/tombstone/root, idempotent recovery, barrier. | Wrong commit/late write/resurrection fails. |
| V23 | v1/v2 coexistence, predecessor archive, unknown fields, interrupted migration. | Loss/no rollback/mixed manifest fails. |
| V24 | Multi-character baseline parity and per-character scopes. | Cross-character leak/behavior change fails. |

All rows remain `UNRESOLVED`/`NOT TESTED` until evidence exists; V20 is `BLOCKED`. No row passes by prose.

---

## 10. Status

**R3:** proposal only; not frozen; not accepted. Independent adversarial review is required before any freeze decision.  
**Frozen contract:** unchanged; hash verified.  
**F3-C phase 0:** **NOT CLOSED**.  
**P1a:** **UNAUTHORIZED**, downstream of P1c.  
**F3-C acceptance:** **NOT GRANTED**.  
**F3-D:** **BLOCKED**.  
**F3-B:** closed.  
**A9/A10/A12/A13:** closed; not reopened.  
**V20:** **BLOCKED** pending F3-F evidence.

No implementation, tests, builds, model loads, migrations, destructive tests, source changes, or contract changes were performed.

---

## Appendix — evidence anchors

- Intent/deletion/recovery: `F3C-MASTER...:690-722`, `:771-805`, `:825-875`; `persistence.go:471`; `app-wipe.go:1-125`.
- Stream/commit: `app-prompt.go:25-146`; character path `:160-244`.
- Hot-swap: `app-model.go:28-174`, `:317-512`, `:537`; `app-lifecycle.go:221`; `chat-model.js:347-393`.
- Provenance: `app-prompt.go:932`; `types/message.go:31-43`.
- V1/v2: `F3C-MASTER...:275-350`, `:777-805`; `persistence.go:585`.
- Media/F3-F: `FREEZE-F3C-DATA...:56-62`; `REVIEW-F3C-DEPENDENCY...:74`; `gallery.go:94-128`.
- Gate: `GATE-F3C...:201-245`.
