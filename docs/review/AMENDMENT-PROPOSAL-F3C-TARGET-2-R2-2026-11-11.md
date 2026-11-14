# AMENDMENT PROPOSAL — F3C-TARGET-2, Revision 2

**Status:** **PROPOSAL ONLY — NOT FROZEN AUTHORITY, NOT ACCEPTANCE.**  
**Date:** 2026-11-11  
**Supersedes:** `docs/review/AMENDMENT-PROPOSAL-F3C-TARGET-2-2026-11-11.md` for review only. The original proposal remains historical.  
**Reviewed against:** `docs/review/REVIEW-ADVERSARIAL-F3C-TARGET-2-2026-11-11.md`  
**Frozen authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`  
**Scope:** bounded contract-resolution pass. No implementation, tests, builds, model loads, source changes, or freeze.

This revision resolves the ten adversarial findings and the two owner refinements. It remains a proposal. The adversarial review is accepted as a review record, not as acceptance of TARGET-2.

---

## 1. Resolution summary

| Finding | Resolution |
|---|---|
| R1 — intent outside deletion/wipe ownership | Generation intents move to a root-owned, chat-scoped, tombstone-aware operation area; wipe stages the area; deletion terminalizes/cleans intents; recovery refuses tombstoned/absent chats. |
| R2 — no commit linkage or switch policy | Intent binds root incarnation, chat, session, epoch, base revision, user commit, assistant commit, and candidate message; switch terminalizes; no cross-chat adoption. |
| R3 — third journal without precedence | A single recovery precedence is defined: root/wipe → tombstone → ordinary journal → cast journal → generation intents. Ambiguous nonterminal combinations block. |
| R4 — rollback aspiration | A runtime transition record, explicit state machine, snapshot of the prior model, one restore attempt, crash behavior, and one stream policy are defined. |
| R5 — provenance not actual runtime config | Effective configuration is captured at attempt start; model identity, parameter schema, seed encoding, partial/retry rules, and all generation paths are frozen. |
| R6 — v1 predecessor/unknown fields | Unknown v1 fields block a v2 write rather than being silently dropped; v1 target bytes are explicitly not retained, only the receipt/manifest chain; first v2 write is staged and recoverable. |
| R7 — media scope ambiguity | F3-F is required for attachment durability and full-workspace acceptance, not for core text/context/provenance work. Core and full acceptance are separated. |
| R8 — multi-character boundary | Explicit non-goal and V24 regression requirement; no cast, prompt, memory, or context-flush changes. |
| R9 — optional-byte rules | Exact optional-field, omission, numeric, array, and unknown-field rules are defined. |
| R10 — non-streaming/character paths | One generation-attempt protocol covers `SendMessage`, streaming, and character responses. |

The two owner refinements are resolved in §4 and §5:

- a failed load never rewrites the saved continuation preference;
- interrupted recovery is idempotent, deletion-safe, and never creates a completed assistant turn from a partial response.

---

## 2. Generation attempt: ownership, linkage, and recovery

### 2.1 Intent is subordinate to chat identity and deletion

The intent path becomes a root-owned operation area:

```text
f3c-generation-intents/<F3C-PATH-GEN-1 hash>/<operationID>.json
```

It is not a target and not a receipt. It is subordinate to:

- root identity and root incarnation;
- chat existence;
- tombstone/deletion;
- ordinary journal precedence;
- wipe intent.

An intent must never outlive or outvote a tombstone. Recovery of a tombstoned or absent chat cleans the intent; it does not recreate the chat, create a target, create a receipt, or revive the operation.

`DeleteChat` must, under the same chat/root locks as the tombstone:

1. verify the last commit evidence;
2. write the tombstone before directory removal;
3. terminalize or remove every intent for that chat;
4. remove the chat directory;
5. leave no live intent that can be adopted later.

If deletion is interrupted after the tombstone, recovery cleans the intents and completes deletion. A tombstone is authoritative over an intent.

Wipe staging must include `f3c-generation-intents/` explicitly. A root wipe removes intents; no intent from an earlier root incarnation may be adopted.

### 2.2 Intent fields

The intent must carry at least:

```text
version, rootId, rootIncarnation, chatId, operationId,
operationType, stateSessionId, baseEpoch, baseRevision,
userMessageId, userCommitOperationId, userMessageHash,
assistantMessageId, assistantCommitOperationId,
assistantCandidateFingerprint, modelId,
startedAt, status, terminalReason
```

Rules:

- `rootIncarnation` must match the current root incarnation. A mismatch is an orphan, not a recoverable generation.
- `chatId` must exist and must not be tombstoned.
- `stateSessionId` and `baseEpoch` bind the attempt to the state that started it.
- `baseRevision` is the revision after the durable user-turn commit.
- `userCommitOperationId` is the commit that made the user turn durable.
- `assistantMessageId` is the candidate assistant message ID.
- `assistantCommitOperationId` is allocated and written to the intent before the assistant commit attempt.
- `assistantCandidateFingerprint` binds the candidate bytes/state that the commit would prove.
- `status` is one of `started`, `streaming`, `completion_pending`, `committed`, `interrupted`, `aborted`, `failed`, or `stale`.
- `terminalReason` records switch, delete, abort, error, stale base, or completion.

No generation starts without a durable user-turn commit. If the user turn cannot be committed, the generation does not start.

### 2.3 One active attempt per chat

Generation attempts are serialized per chat/session. A second attempt is rejected or explicitly queued; it never writes a second intent for the same chat/session while one is active. The existing single-stream token model is preserved.

A retry is a **new** operation ID and a **new** assistant candidate. It does not reuse, overwrite, or merge the prior attempt. The prior attempt must be terminal before the retry starts.

### 2.4 Switch behavior

A chat switch:

- binds the outgoing chat;
- cancels an active stream and records `aborted` with a switch reason;
- waits until the outgoing intent is terminal or clean;
- does not adopt the outgoing intent into the incoming chat;
- does not reuse the outgoing `stateSessionId`.

A completed response whose commit is pending must not be silently discarded by a switch. The switch either waits for the commit/rollback decision or explicitly reports that the visible response is unsaved. It never reports completion.

### 2.5 Recovery precedence

The single recovery order is:

1. **Root identity/incarnation.** Mismatch → orphan; no adoption.
2. **Wipe intent.** If present, no chat/generation recovery; complete or block the wipe.
3. **Tombstone inventory.** If the chat is tombstoned, clean intents and do not recreate.
4. **Ordinary create/save/append journal.** A nonterminal ordinary journal takes precedence; do not adopt a generation intent for the same chat.
5. **Cast journal.** A nonterminal cast journal blocks generation adoption; no independent cast/generation race.
6. **Generation intents.** Only for a live, non-tombstoned chat with no blocking ordinary/cast journal.
7. **Terminal cleanup.** Record terminal status and clean artifacts; no chat revision is created by recovery.

If an ordinary or cast journal and a generation intent are simultaneously nonterminal for the same chat, recovery returns `ErrGenerationJournalAmbiguous`; it does not choose one silently.

### 2.6 Idempotent recovery

Recovery is a deterministic decision over the intent, current chat evidence, and receipts. It is repeatable without creating duplicate user turns, assistant turns, attempts, or receipts.

| Evidence | Recovery result |
|---|---|
| Intent missing | Ordinary load only. |
| User commit missing | `ErrGenerationAnchorMissing`; no target; block/clean. |
| Assistant commit/receipt exists | Mark `committed`; clean; never append again. |
| No assistant commit, base current, no tombstone | Mark `interrupted`; no target; no partial assistant. |
| No assistant commit, base stale | Mark `stale`; no target; do not merge. |
| `completion_pending`, no receipt | Mark `interrupted` or `failed`; no target. |
| Tombstoned/absent chat | Delete/terminalize intent; no target, no recreation. |
| Root incarnation mismatch | Orphan; clean; no target. |
| Malformed/conflicting intent | Block with `ErrGenerationIntentMalformed`; no adoption. |
| Wipe intent present | No generation recovery; wipe precedence. |

A partial assistant response is **not** persisted as a completed assistant turn. This revision does not retain partial tokens as a recoverable draft. A retry starts a new generation. If the owner later wants an explicit interrupted-draft sidecar, that is a separate amendment; it is not silently introduced here.

### 2.7 All assistant-producing paths

The same attempt protocol applies to:

- `SendMessage` (synchronous);
- `SendMessageStreamWithImage` / `SendMessageStream`;
- `GenerateCharacterResponse`.

No path may append an assistant turn without a user-turn commit, intent, and assistant commit/receipt. Character generation uses the same intent but keeps its existing character context/prompt behavior.

### 2.8 Durable events

- `chat:stream_end` means the token stream ended; it is not durable proof.
- `chat:committed` means the assistant commit/receipt is verified.
- `chat:save_failed` means the response is visible but not durably committed.
- `chat:interrupted` means an intent exists without a committed assistant turn.

`chat:complete` may remain for frontend compatibility only if it means stream-end, never durable completion.

---

## 3. Hot-swap and continuation state

### 3.1 Three separate concepts

| Concept | Owner | Durability |
|---|---|---|
| `savedContinuation` | `meta.continuationModel` | Durable chat preference. |
| `requestedModel` | Runtime/UI selection | Ephemeral until a successful load. |
| `activeModel` | Runtime loader/process | Ephemeral; must be healthy. |

Only `activeModel` may generate. Only a successful selection may update `savedContinuation`.

### 3.2 Preference rule

- A failed load never rewrites `savedContinuation`.
- If A is saved and loaded, and the user requests B, B fails, and A is restored, saved remains A and active is A.
- If A cannot be restored, saved remains A and active is `unavailable`. The UI must say A is unavailable and B failed; it must not claim A or B is ready.
- If B becomes healthy and the user’s selection is accepted, B becomes active and may become saved.
- An explicit preference-only action may save an unavailable model, but must mark it unavailable and must not claim it is loaded. Normal selection does not do this.
- On explicit unload, saved remains; active becomes unavailable.
- On crash/restart, saved remains; active is absent; no automatic model load and no substitution.

### 3.3 Transition state machine

| State | Meaning | Allowed next |
|---|---|---|
| `idle` | No transition. | `preflight` |
| `preflight` | Validate model and resources before stop. | `blocked`, `wait_for_stream`, `capture_previous` |
| `wait_for_stream` | Active generation must finish; no automatic cancel. | `capture_previous` after stream end, or `idle` after explicit abort |
| `capture_previous` | Snapshot old path/projector/settings/evidence/estimate. | `stopping` |
| `stopping` | Stop old server; verify exit. | `starting`, `unavailable` |
| `starting` | Spawn replacement. | `health_check`, `rollback` |
| `health_check` | Verify health, not process existence. | `publishing`, `rollback` |
| `publishing` | Publish active model; optionally commit preference. | `ready` |
| `ready` | Healthy active model. | `idle` |
| `rollback` | Restore previous model once. | `ready`, `unavailable` |
| `unavailable` | No healthy model; no false active selection. | `preflight` |
| `blocked` | Preflight failed; old state unchanged. | `idle` |

Rules:

- Only one transition may run at a time.
- `LoadModel` returns only after `ready` or a terminal failure; if the product needs an accepted operation, a single explicit `model:ready` event is the only publication signal.
- The frontend does not update `howl-last-model-path` or `savedContinuation` before `ready`.
- Active stream policy is **wait**, not cancel. If the user explicitly aborts the stream, transition may proceed. A transition must not terminate an active stream by itself.
- `App.serverProc`, the HTTP client, `resource.Manager`, and `model.Manager` are updated by one transition owner and reconciled before `ready` or `unavailable`.

### 3.4 Failure matrix

| Failure | Required result |
|---|---|
| Path/admission/eligibility/projector rejection | Old model and selection unchanged; no stop. |
| Resource-fit rejection | Old model unchanged; no stop. |
| Stop failure | Do not start replacement; health-check old; report stop failure; no false selection. |
| Spawn failure | Kill partial process; one restore with the captured prior snapshot; otherwise `unavailable`. |
| Health timeout | Stop failed process; one restore; otherwise `unavailable`. |
| Restore failure | `unavailable`; saved preference retained; no fallback model. |
| Concurrent transition | Reject/serialize; never two transitions. |
| Crash mid-transition | No durable active model; startup reports no active model; saved preference remains; no auto-load. |
| Preference commit failure after ready | Model is active but preference is not durable; report that distinction. |

No failure branch may publish a model as loaded merely because a process was spawned or a preference was saved.

---

## 4. Per-turn provenance

### 4.1 Capture point

Capture the effective configuration at attempt start, after model defaults, clamping, template injection, and backend normalization, but before the first token. Never reconstruct it from current settings at save time.

### 4.2 Model identity

`modelId` is a stable opaque identity:

- use the model registry's persisted identity when available;
- otherwise use a documented digest identity over the model artifact, not a path or display name;
- a moved model keeps identity only if its identity evidence is unchanged;
- a quant or content change creates a different identity;
- a removed model remains an identity with an `unavailable` availability state.

The exact registry/digest implementation must be frozen in the model identity record before TARGET-2 freeze. A path, filename, or display name is never authority.

### 4.3 Exact provenance object

For a v2 assistant turn, `generation` is present with:

```text
modelId, displayName, engine, backend, parameters,
seedMode, seed, effectiveSeed, templateFingerprint
```

- `modelId`, `engine`, `backend`, `parameters`, `seedMode`, and `templateFingerprint` are required.
- `displayName` is optional and display-only.
- `parameters` has the exact required keys `temperature`, `topP`, `topK`, `repeatPenalty`, `maxTokens`, and `contextTokens`; no unknown keys; numbers are finite; values are the effective runtime values.
- `seedMode` is `explicit` or `random`.
- `seed` is present only for `explicit`, is a signed 64-bit integer, and `0` is valid.
- `effectiveSeed` is optional and only for a random seed actually reported by the engine.
- `templateFingerprint` is `sha256:<64 lowercase hex>` or `none`.

`legacy_unknown` is represented by absence of `generation`, never by an invented model or a `-1` seed.

### 4.4 Partial attempts and retries

- A partial/interrupted attempt has an intent, not a committed assistant turn.
- No partial response is recorded as a completed assistant entry.
- A retry is a new operation, a new assistant candidate, and a new generation object.
- The previous attempt remains terminal evidence; it is not overwritten or merged.
- A mid-stream model change does not affect the attempt; the model is frozen at attempt start.

### 4.5 All generation paths

The same provenance helper applies to synchronous, streaming, and character responses. `characterId` remains the existing owner field; provenance does not redefine character prompt/context behavior.

---

## 5. V1-to-V2 compatibility

### 5.1 What is preserved

- v1 receipts, records, and `tgt1-` target manifests remain valid;
- old receipts are never rehashed;
- v1 current chats remain openable at `F3C-TARGET-1`;
- the first new write creates a v2 revision with a v1 predecessor;
- no historical receipt is rewritten or copied as v2 proof.

### 5.2 Unknown fields

Unknown fields in v1 target files are **not silently dropped**.

- A v1 file with unknown target fields is read-only/openable but cannot be upgraded to v2.
- A v2 write attempt returns `ErrLegacyUnknownField` and leaves the v1 source unchanged.
- A separate repair/export path may be defined later; this amendment does not invent a `legacyExtras` field.
- V2 target files with unknown fields are malformed and fail closed.

This is a deliberate strict policy: preserve the source by refusing to rewrite it, rather than claiming lossless migration the target cannot represent.

### 5.3 Predecessor bytes

Historical proof is receipt/manifest/hash-level, not byte-level retention of the old target files. After a successful v2 commit, the old v1 target bytes are not retained as part of the target set; the v1 receipt/record/manifest chain remains.

If the owner requires byte-level predecessor retention, that is a separate archive requirement, not an implied promise of `F3C-TARGET-2`.

### 5.4 First v2 write

The first v2 write is a staged ordinary commit:

1. v1 remains current;
2. build and stage v2 targets;
3. read back and verify;
4. verify v1 base, receipt, and chain;
5. create record/receipt with a v1 predecessor and v2 manifest;
6. install and publish v2.

Pre-commit failure leaves v1 current. Post-commit publication failure follows ordinary recovery. No mixed v1/v2 target manifest is valid.

### 5.5 Migrated entries

- `meta.continuationModel` is absent unless the user explicitly selects a model.
- Migrated history entries have `attachments: []` and no `generation`.
- `hasImage: true` with no attachment remains a legacy-unknown image.
- No global last-model path is written into an existing chat.
- Duplicate/empty/malformed IDs block; IDs are not renumbered.

### 5.6 Save As

`Save As` from v1 to v2:

- source binding and source revision are verified;
- unknown v1 fields block the save;
- old receipts are not copied;
- new revision is v2;
- provenance is not invented.

---

## 6. Media dependency and acceptance split

### 6.1 F3-F is not a prerequisite for all F3-C work

Core F3-C work may proceed without F3-F:

- text/context persistence;
- provenance;
- continuation preference;
- generation intent/recovery;
- v1/v2 compatibility.

Attachments remain `reference-only` or `unavailable` until F3-F proves durability. F3-C must not claim attachment durability before then.

### 6.2 Full-workspace acceptance requires F3-F

`F3-C-FULL` includes attachment durability and requires F3-F ownership evidence. `F3-C-CORE` excludes attachment durability and may be accepted independently.

| Acceptance | Includes | F3-F required |
|---|---|---|
| `F3-C-CORE` | Text, context, provenance, continuation, interruption, v1/v2 | No |
| `F3-C-FULL` | Core plus durable attachments | Yes |

V20 remains **BLOCKED**, not PASS, until F3-F exists.

### 6.3 Media reference rules

- `attachments` entries use `mediaId`, `role`, `provider`, `mimeType`, `byteSize`, and `contentHash`.
- `provider` is `f3f-media` or `legacy-indexeddb`.
- F3-F owns bytes, presence, corruption, deletion, and backup.
- F3-C deletes references, not arbitrary F3-F bytes.
- Missing, corrupt, or legacy-unknown media is unavailable; no substitution.
- No base64, local path, or thumbnail bytes in the target.

---

## 7. Multi-character boundary

F3C-TARGET-2 does not change:

- character-specific context flushing;
- character memory isolation;
- character-specific prompt construction;
- cast synchronization semantics;
- multi-character interaction behavior.

The persistence layer preserves the information those systems need:

- `characterId` remains the existing message owner;
- `generation` is per assistant turn;
- `attachments` is per message;
- `continuationModel` is per chat, not per character;
- no generation intent or recovery path may alter cast state.

### V24 — multi-character isolation regression

**Claim:** existing character-specific context, memory, prompt, and cast behavior is unchanged for equivalent inputs.

**Required evidence:**

- existing character-context and cast suites;
- two characters with distinct contexts produce the same prompts/context as before the amendment;
- persistence/recovery does not move, duplicate, or cross-assign character context;
- a per-turn generation record carries the correct `characterId`;
- attachments/provenance do not enter another character’s prompt.

**Negative case:** an implementation that shares one character’s memory, prompt, or cast state with another fails V24.

V24 is a regression requirement, not a new cast contract.

---

## 8. V2 exact serialization

### 8.1 `meta.json`

```text
id, name, createdAt, updatedAt, worldName, scenarioName,
characterNames, lorebookNames, worldId, scenarioId,
activeCharacters, removedCharacters, addedCharacters,
activeLorebooks, removedLorebooks, addedLorebooks,
rootId, revision, lastCommitOperationId, continuationModel
```

`continuationModel` is omitted when absent. When present:

```text
modelId, displayName, sourceHint, selectedAt
```

`modelId` and `selectedAt` are required; `displayName` and `sourceHint` are optional and omitted when empty.

### 8.2 History

```text
id, sender, text, role, timestamp, hasImage, characterId, generation, attachments
```

- `generation` is omitted for non-assistant and legacy entries; present for v2 assistant entries.
- `attachments` is always `[]` for v2 entries; never `null`.
- `hasImage` is a compatibility projection.
- A v2 entry with `hasImage: true` and `attachments: []` is a legacy-unknown image.
- `history.json` and `history.ndjson` use identical presence rules and values.

### 8.3 Attachment entry

```text
mediaId, role, provider, mimeType, byteSize, contentHash
```

`mediaId`, `role`, `provider`, `mimeType`, `byteSize` are required. `contentHash` is required for `f3f-media` and optional for `legacy-indexeddb`. No bytes or paths.

### 8.4 Unknown fields and optional rules

- Unknown fields in v2 target objects are `ErrTargetMalformed`.
- Missing required fields are malformed.
- `null` is not a substitute for an omitted optional field.
- `attachments` is never omitted for v2; `generation` is omitted only for legacy/non-assistant.
- Maps use deterministic sorted keys; arrays preserve order.
- `history.ndjson` uses compact marshal, one LF per entry, one trailing LF; `history.json` uses two-space indent, no trailing newline.

### 8.5 Manifest and hashes

- `f3c-target-manifest-2` for v2 targets.
- `tgt2-` for v2 present targets.
- `tgt1-` remains valid for v1.
- `F3C-CANON-1` remains unchanged.
- Receipt body/record/envelope shapes do not change; the embedded manifest identifies the format.
- Exactly eight target paths; no ninth target.

---

## 9. Gate v2 mapping update

The Gate v2 document is not amended by this proposal. This table is the required mapping if TARGET-2 is later frozen.

| Row | Testable claim | Required negative/evidence |
|---|---|---|
| V4 | v1 and v2 exact bytes, manifests, `tgt1`/`tgt2`, optional fields. | Unknown field, wrong omission, wrong array, or mixed manifest fails. |
| V5 | User turn durable before attempt; autosave; switch; stream-end is not commit. | Switch with active intent does not adopt; `chat:complete` cannot prove commit. |
| V7 | Save As v1→v2 preserves source binding and blocks unknown fields. | No receipt copy; source mismatch fails. |
| V8 | Delete tombstone terminalizes/cleans intents; no resurrection. | Tombstoned chat with intent cleans, never recreates. |
| V11 | v1/v2 history IDs/order; no invented provenance. | Unknown fields block; duplicate/empty IDs fail. |
| V12 | Wipe stages chat, receipt/commit/tombstone/intent roots, gallery boundary. | Generation intent or gallery staging failure blocks; no partial wipe. |
| V13 | Fault injection at intent start, stream, commit, migration, rollback. | No ambiguous adoption or duplicate append. |
| V15 | One transition/one active attempt; switch/stream races. | Concurrent generation/transition fails closed. |
| V16 | Native Wails/WebView2 path, preference restore, unavailable state. | No false loaded; no model substitution. |
| V18 | Effective provenance captured at attempt start; all paths. | Requested settings changed later do not alter old turn. |
| V19 | Explicit zero, explicit nonzero, random known/unknown, legacy absent. | `-1` or a guessed seed fails. |
| V20 | F3-F media ownership, present/corrupt/missing. | **BLOCKED** until F3-F; reference alone cannot pass. |
| V21 | Saved/requested/active separation; rollback; crash. | Failed load does not rewrite preference; no false active. |
| V22 | Intent recovery idempotent, deletion-safe, all paths. | Repeat recovery no duplicates; tombstone wins. |
| V23 | v1/v2 coexistence, predecessor chain, unknown-field policy, interrupted first v2. | Old receipts valid; no half-migration. |
| V24 | Multi-character context/prompt/cast parity. | Cross-character context leak or altered prompt fails. |

The Gate v2 row inventory becomes V1–V24; until frozen, all new rows are `UNRESOLVED`, and V20 is `BLOCKED`. No row passes by test name, positive case alone, or prose.

---

## 10. Decision status and boundaries

**This revision:** proposal only; not frozen; not acceptance.  
**F3C-TARGET-2:** not frozen.  
**Adversarial review:** accepted as a review record, not as acceptance.  
**Phase 0:** remains **NOT CLOSED**.  
**P1a:** not authorized; remains downstream of P1c.  
**F3-C:** not accepted.  
**F3-B:** closed; not reopened.  
**F3-D:** blocked.  
**A9/A10/A12/A13:** remain closed and are not reopened.

No implementation, tests, builds, model loads, source changes, or contract changes were performed. The next step is owner review of this R2 revision, then a fresh adversarial review against the frozen authority hash. Only after that can a freeze decision be considered.

---

## Appendix — resolution-to-evidence map

| Resolution | Contract/source basis |
|---|---|
| Intent deletion/wipe | `F3C-MASTER...:694-716` tombstone; `:771-805` wipe staging; `persistence.go:471` current delete. |
| Recovery precedence | `F3C-MASTER...:825-875` journal recovery; `app-wipe.go:1-125`. |
| Stream/commit split | `F3C-MASTER...:686`; `app-prompt.go:138-146`. |
| Hot-swap readiness | `app-model.go:28-174`, `:317`, `:345-512`, `:537`; `app-lifecycle.go:221`. |
| Preference vs active | `chat-model.js:347-393`; `types.go:50`; `types.go:72`. |
| Effective provenance | `app-prompt.go:932`; `types/message.go:31-43`; `llama/options.go:76`, `:156`. |
| V1 unknown/predecessor | `F3C-MASTER...:275-350` target rules; `:777-805` history migration; `persistence.go:585`. |
| Media/F3-F | `FREEZE-F3C-DATA...:56-62`; `REVIEW-F3C-DEPENDENCY...:74`; `gallery.go:94-128`. |
| Multi-character boundary | `app-prompt.go:160`; existing character context/cast suites; no target changes. |
| Gate v2 | `GATE-F3C...:201-245`; V18–V24 proposed. |
