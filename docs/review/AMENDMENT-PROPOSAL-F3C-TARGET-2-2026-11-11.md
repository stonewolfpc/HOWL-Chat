# AMENDMENT PROPOSAL — F3C-TARGET-2

**Status:** **PROPOSAL ONLY — NOT FROZEN AUTHORITY.**  
**Date:** 2026-11-11  
**Basis:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Frozen authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`  
**Related assessment:** `docs/review/ASSESSMENT-F3C-PHASE-0-FINAL-REQUIREMENTS-IMPACT-2026-11-11.md`  
**Scope:** contract amendment proposal and read-only source trace. No implementation, tests, builds, model loads, or source changes.

This document does not amend the frozen contract. It does not authorize P1a. It does not accept F3-C. It does not reopen A9/A10/A12/A13. The owner must accept or reject this proposal before a new freeze and Gate v2 binding can exist.

---

## 1. Executive proposal

Adopt a versioned target format, `F3C-TARGET-2`, that preserves the eight-target set but adds:

1. **Per-turn generation provenance** — the model, engine, backend, relevant generation parameters, template fingerprint, and seed state that actually produced each assistant turn.
2. **A per-chat continuation preference** — the model the user last selected for this chat, separate from historical provenance and separate from runtime loaded state.
3. **Media attachment references** — typed references owned by F3-F, never bytes copied into the chat target.
4. **A durable generation-intent journal** — so an interrupted stream is recoverable as “interrupted,” not falsely completed or silently discarded.
5. **A compatibility boundary** — `F3C-TARGET-1` remains valid for historical and current legacy revisions; old receipts are never invalidated or rehashed.

The amendment does **not** put runtime model state, partial tokens, or media bytes into the target set. F3-C owns durable chat identity and references. F3-F owns media bytes. The hot-swap loader remains runtime machinery, but its publication and failure rules must be corrected so a chat preference is never reported as loaded before it is ready.

**The media precondition is explicit:** F3-C cannot claim that a chat image is durable merely because a reference exists. Until F3-F owns the bytes and reports them present with matching integrity, an attachment is `reference-only` or `unavailable`. Current chat image data is transient base64 and a `hasImage` boolean; it has no durable media owner.

---

## 2. Current-state evidence

### 2.1 The target set cannot represent the visible creative session

The frozen target set is fixed at eight paths (`F3C-MASTER...:260`): `meta.json`, `characters.json`, `history.json`, `history.ndjson`, `lorebooks.json`, `memories_snapshot.json`, optional `world.json`, optional `scenario.json`.

`F3C-TARGET-1` fixes history to exactly `id, sender, text, role, timestamp, hasImage, characterId` (`:275`). `ChatHistory` matches that shape (`internal/backend/chat/types.go:72`). The richer `types.Message.Metadata` has `Model`, temperature, top-p, top-k, images, audio, and documents (`internal/backend/types/message.go:31-43`), but it is not the persisted chat shape.

Consequences:

- Only an image **presence boolean** survives (`HasImage`); image bytes and references do not.
- The actual generating model is not persisted with the turn.
- The seed mode and value are not persisted per turn.
- The chat has no continuation-model field in `ChatMeta` (`internal/backend/chat/types.go:50`).
- Stable IDs are required by the contract, but current IDs are array-length/DOM-derived (`app-chat.go:373`; `app-prompt.go:281`; `frontend/dist/js/data-sync.js:238`).

### 2.2 Current image ownership is not durable chat media

The F3-C ownership freeze already states the intended boundary:

- image bytes stay under F3-F media roots;
- F3-C never copies, moves, or writes image bytes;
- legacy IndexedDB references are preserved as `provider: "legacy-indexeddb"` until F3-F migrates them (`docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md:56-62`).

The dependency review reaches the same conclusion: F3-C may define references now, but canonical image ownership and atomic PNG+sidecar commit belong to F3-F (`docs/review/REVIEW-F3C-DEPENDENCY-AND-ACCEPTANCE-2026-11-10.md:74`). The compatibility review requires F3-F identity, not form-invented IDs, and says F3-C must not lose `{id, stored:true}` references (`docs/review/REVIEW-F3C-DECISIONS-AND-COMPATIBILITY-2026-11-10.md:154-162`).

What exists today is weaker than that boundary:

- Chat image input is transient base64 in `currentImage` (`frontend/dist/js/chat-input.js:20`, `:50`, `:240`), cleared after send.
- Chat persistence records only `hasImage` (`frontend/dist/js/chat-actions.js:62`; `frontend/dist/js/chat-session.js:146`).
- World/scenario images are separate IndexedDB images (`storage-manager.js:40-84`; `storage-indexeddb.js` stores `world-images`, `scenario-images`), not chat message attachments.
- Generated gallery images use a separate `img_<timestamp>_<nano>` identity and a PNG plus JSON sidecar under `paths.StatePath("images", "gallery")` (`app-image.go:56`; `internal/image/gallery/gallery.go:94-128`). The sidecar write removes the PNG if it fails, but this is not a registered F3-F media transaction with recovery, ownership, or a reference index.
- No dedicated `docs/**/*F3-F*` contract file exists in the tree. F3-F is a planned boundary, not an accepted implementation.

**Therefore:** a chat can reopen with `hasImage: true` and no image at all. That is not a cosmetic gap; it is the exact trust failure the amendment must stop.

### 2.3 Model selection is global, not per chat

`ChatMeta` has no model field. `LoadChat` restores chat state but not model selection (`app-chat.go:96-126`). The frontend stores `howl-last-model-path` globally in `localStorage` (`frontend/dist/js/chat-model.js:169`, `:393`); `autoLoadDefaultModel` uses that global path (`:166`).

The picker shows a model after `await app.LoadModel(modelPath)` returns, then writes the global last path (`chat-model.js:347-393`). The backend `LoadModel` has already registered the model and returned before the health monitor emits `model:loaded` (`app-model.go:345-512`; monitor at `:360-500`). This is a false-success window: the UI can publish “selected” before readiness.

### 2.4 Hot-swap is a partial manual path

`App.LoadModel` validates path, GGUF admission, chat eligibility, and projector pairing before stopping the current server (`app-model.go:28-174`). That preflight is sound and should be preserved.

After preflight it:

1. stops the existing server (`:174`);
2. waits for port 8080;
3. validates the server binary;
4. extracts the GGUF profile;
5. validates spawn parameters;
6. spawns and adopts the process (`:317`, `:345`);
7. starts a monitor;
8. registers resources and returns before health (`:512`).

There is no rollback after the stop. A spawn or startup failure leaves the old model gone and the new model not ready. There is no `streamMu` acquisition in `LoadModel`, so a model stop can race an active stream. `stopLlamaServer` kills the owned tree and waits (`app-lifecycle.go:221`), but no transaction restores the previous model. `UnloadModel` simply stops and unregisters (`app-model.go:537`).

The app also has a separate `model.Manager` lifecycle used by the chat service. The direct `App.LoadModel` path and the manager can disagree about what is loaded. The amendment must require one reconciled runtime truth; it must not create another model authority.

### 2.5 Interrupted generation is currently invisible to persistence

A stream:

- records the user turn in memory (`app-prompt.go:58`; `recordUserTurn` at `app-chat.go:355-381`);
- starts a stream token (`app-prompt.go:64`);
- sends chunks to the UI;
- appends the assistant message only after the stream completes (`app-prompt.go:138`);
- emits `chat:complete` before the asynchronous save (`app-prompt.go:139`);
- calls `performSave` in a goroutine that only logs failure (`app-prompt.go:138-146`).

`performSave` loads the last chat, updates it from app state, and calls `SaveChat` (`app-chat.go:314-327`). `recordUserTurn` itself only mutates memory. `LoadChat` prints a warning and continues after a failed outgoing save (`app-chat.go:96-111`). `loadHistoryFile` stops at a corrupt NDJSON line and returns the partial prefix (`internal/backend/chat/persistence.go:585`).

There is no generation-intent record. The UI can show a partial response that is not in history, and `chat:complete` can be emitted before a durable save. A crash therefore cannot distinguish “completed but unsaved” from “interrupted mid-stream.”

### 2.6 Legacy behavior

`F3C-TARGET-1` and its `tgt1-` hashes are already accepted for existing revisions. Any amendment that makes a legacy `tgt1-` receipt invalid, rehashes historical receipts, or treats a v1 chat as corrupt without a write is rejected by this proposal. Compatibility is a first-class requirement, not a migration afterthought.

---

## 3. Scope boundaries

### F3-C owns

- durable chat identity, revisions, targets, receipts, and commit proof;
- message identity/order and durable text;
- per-turn provenance and seed state;
- per-chat continuation preference;
- typed media references and their availability state;
- generation-intent/recovery records that are not target revisions;
- migration/version dispatch for chat persistence.

### F3-F owns

- media bytes, media roots, stable media identity, atomic PNG+sidecar commit;
- media presence, corruption detection, deletion, and backup/recovery;
- the authority that says a referenced attachment is actually durable.

### F3-C does not own

- loading, stopping, or monitoring `llama-server`;
- runtime model state (`loading`, `ready`, `unavailable` as active-process state);
- copying image bytes into the chat directory;
- deleting F3-F media because a chat was deleted;
- claiming that an IndexedDB cache is canonical media.

### Hot-swap is outside the target set

The target stores only the continuation preference and historical provenance. Runtime phases are ephemeral. The amendment still defines the hot-swap publication/failure rules because a chat preference that can be falsely published is not a usable feature.

---

## 4. `F3C-TARGET-2` target model

### 4.1 Target set remains eight

No ninth target file. The eight paths and path-hash rules remain. The format version is carried by the target manifest, and target hashes use the `tgt2-` domain for v2 bytes. `F3C-CANON-1` remains the receipt/record canonicalization rule.

### 4.2 `meta.json`

Append one optional field after `lastCommitOperationId`:

```text
id, name, createdAt, updatedAt, worldName, scenarioName,
characterNames, lorebookNames, worldId, scenarioId,
activeCharacters, removedCharacters, addedCharacters,
activeLorebooks, removedLorebooks, addedLorebooks,
rootId, revision, lastCommitOperationId, continuationModel
```

`continuationModel` is omitted when no explicit per-chat selection exists. It is not populated from the global last-model path.

### 4.3 `history.json` and `history.ndjson`

Append two fields after `characterId`:

```text
id, sender, text, role, timestamp, hasImage, characterId, generation, attachments
```

- `generation` is omitted for entries whose provenance is legacy-unknown. It is never fabricated.
- `attachments` is always `[]` for v2 entries with no attachment; it is never `null`.
- `hasImage` remains a compatibility projection. It is `true` if a legacy entry claims an image or a current entry has an attachment reference. It is not proof that bytes exist.
- Both files must decode to the same ordered entries and exact same values.

### 4.4 `generation` object

Semantic field order:

```text
modelId, displayName, engine, backend, parameters,
seedMode, seed, effectiveSeed, templateFingerprint
```

Rules:

- `modelId` is a stable, opaque model identity. A path or filename is never the authority. A path may be retained as a non-authoritative `sourceHint` only if a later field is added explicitly.
- `displayName` is display-only and never resolves identity.
- `engine` and `backend` record what actually executed the generation (for example the engine family and server/API adapter), not what the picker displayed.
- `parameters` is a deterministic object of relevant generation controls. At minimum it must be able to carry the values the product treats as generation-relevant (for example temperature, top-p, top-k, repeat penalty, maximum tokens, and context size). It must not contain secrets.
- `seedMode` is `explicit` or `random` for a recorded v2 generation.
- `seed` is present only for `explicit`; `0` is a valid explicit seed and must not be treated as unset.
- `effectiveSeed` is optional; it is present only when the engine reports the actual seed used for a random generation. Its absence does not convert a random generation into an explicit one.
- `templateFingerprint` identifies the prompt/template version used, without storing the template text.
- A v2 entry with `generation` present must have a complete object. Partial or malformed provenance is a target error, not an “unknown” fallback.
- A v1 entry has no `generation`; it is legacy-unknown by absence.

`legacy_unknown` is therefore represented by **absence**, not by an invented model or a `-1` seed. This keeps old history honest.

### 4.5 `attachments` array

Each entry has:

```text
mediaId, role, provider, mimeType, byteSize, contentHash
```

Rules:

- `mediaId` is the F3-F media identity, not a chat-local ID.
- `role` is `attachment` for chat message media; other F3-F roles are not silently reused.
- `provider` is `f3f-media` for canonical F3-F media, or `legacy-indexeddb` only for preserved legacy references.
- `mimeType`, `byteSize`, and `contentHash` are required for a new `f3f-media` reference.
- No bytes, base64, thumbnails, or local absolute paths are stored in the target.
- A missing or corrupt F3-F object is an availability state, not a reason to drop the reference.
- `hasImage: true` with `attachments: []` is a legacy-unknown image, not an empty image.

### 4.6 `continuationModel`

Semantic field order:

```text
modelId, displayName, sourceHint, selectedAt
```

Rules:

- `modelId` is the same stable identity used by provenance.
- It is a preference for the next generation, not a lock, not a claim that the model is installed, and not a rewrite of history.
- `sourceHint` is optional and non-authoritative.
- `selectedAt` is the time of the explicit selection, not the time a load finished.
- If the model is unavailable, the preference is retained and marked unavailable; another model is never substituted automatically.
- A new chat may use the application default. An existing chat must not have its preference overwritten by the global last-model value.

### 4.7 Exact bytes, manifest, and hashes

`F3C-TARGET-2` uses the same JSON rules as v1: `MarshalIndent` with two spaces, UTF-8, LF, no BOM, no trailing newline, frozen field order, deterministic map-key order, and `[]`/`{}` normalization. `history.ndjson` uses compact `Marshal`, one LF per entry, one trailing LF, and no blank lines.

The manifest becomes:

```json
{
  "format": "f3c-target-manifest-2",
  "targets": [
    {"path":"characters.json","presence":"present","bytes":123,"hash":"tgt2-..."},
    {"path":"scenario.json","presence":"absent"}
  ]
}
```

- Exactly the eight target paths, sorted ascending by path.
- `tgt2-` for every present v2 target.
- `tgt1-` remains valid for v1 targets and their historical records.
- The receipt body and commit-record field order do not change. Their embedded `targetManifest` identifies the format, so a separate receipt field is unnecessary.
- `commitRecordHash`, `receiptBodyHash`, and `receiptEnvelopeHash` still cover the same structures with the embedded manifest.
- `F3C-CANON-1` is unchanged.

### 4.8 Revisions and version coexistence

The target format is a property of each committed revision, not a new revision kind.

- A v1 chat may remain at `F3C-TARGET-1` indefinitely; opening it does not rewrite it.
- The first ordinary v2 write advances the revision and writes v2 targets with a v2 manifest.
- A v2 chat may contain history entries without `generation`; these are legacy-unknown.
- Historical manifests are never compared with current files, consistent with §5.7.
- Historical receipts remain valid and are never rehashed into a different target format.
- A v1 predecessor followed by a v2 commit is valid.
- `Save As` from a v1 source creates a new v2 revision, preserving source binding and unknown provenance.
- No receipt is “upgraded” in place.

---

## 5. Per-turn model provenance

### 5.1 What must be recorded

For every assistant turn produced after the amendment, record what actually generated it:

- stable model identity;
- engine and backend adapter;
- generation parameters relevant to reproducibility/explanation;
- template fingerprint;
- seed mode;
- explicit seed, or effective random seed when known.

The record is written from the actual request/generation result, not from the current picker value at save time. A later selection cannot alter it.

### 5.2 Identity

`modelId` must survive model removal and path changes. The proposal leaves the exact registry/digest scheme to the implementation decision, with these constraints:

- a filename or directory path is not identity;
- a display name is not identity;
- the same model moved to another path retains identity if its identity evidence is unchanged;
- two different models with the same display name are distinct;
- an unresolvable historical identity remains an identity with an `unavailable` availability state, not a reason to substitute.

### 5.3 Seed semantics

| Mode | Meaning | Stored state |
|---|---|---|
| `explicit` | The user or request supplied `seed`. | `seedMode: explicit`, `seed: <integer>`; zero is valid. |
| `random` | The engine chose the seed. | `seedMode: random`; `effectiveSeed` only if the engine reports it. |
| legacy-unknown | The turn predates provenance. | No `generation` object; never converted to a guessed seed. |

No `-1` sentinel is used as the only marker. The existing inference default `-1` (`internal/backend/llama/options.go:76`, `:156`) is a runtime default, not a durable seed mode.

### 5.4 Historical model removal

When a model is no longer installed:

- the turn still shows its recorded `modelId` and display metadata;
- the model is marked unavailable for future continuation;
- no historical field is rewritten;
- no different model is presented as the author;
- selecting a different model affects the next generation and the continuation preference only.

---

## 6. Per-chat continuation selection

### 6.1 Restore, do not lock

On chat open:

1. read `meta.continuationModel`, if present;
2. restore the selection in the picker;
3. do not load a different model because a global last-model path exists;
4. do not rewrite the preference merely because the chat opened.

If the model is unavailable, show that state. The user may choose another model. The choice is allowed; it does not change prior provenance.

### 6.2 When the preference is committed

The preference is committed only after an explicit user selection is accepted as the chat's selection. A failed load does not publish a new preference. A successful load may publish it. The exact save operation is an ordinary chat mutation, not a new receipt type.

### 6.3 On send

If the selected continuation model is not ready, Howl loads it before generation. If it cannot be loaded:

- do not generate with a different model;
- do not silently fall back to the global default;
- report the selected model as unavailable;
- leave the preference and history intact.

This is the “start typing and Howl handles the machinery” behavior, without model substitution.

### 6.4 New chats

A new chat may inherit the application default. That default is not written into an existing chat's `continuationModel` merely because it was used for a new chat.

---

## 7. Hot-swap and model transition

### 7.1 Preserve the sound preflight

Keep the current path validation, GGUF admission, chat-task eligibility, and projector identity checks before stopping any running server (`app-model.go:28-174`). These checks are the existing partial capability and should not be rewritten.

### 7.2 Required state machine

Add one serialized transition authority around the existing loader. Proposed states:

| State | Meaning | Durable effect |
|---|---|---|
| `idle` | No transition. | None. |
| `preflight` | Validate target and resources before stopping. | None. |
| `blocked` | Preflight rejected; old state unchanged. | None. |
| `draining` | No new generation; active stream is cancelled or completed. | Generation intent records abort if cancelled. |
| `stopping` | Old server is being stopped. | None to target. |
| `starting` | New server spawned. | None to target. |
| `health_check` | Waiting for readiness, not mere process existence. | None to target. |
| `publishing` | Publish active selection and, when explicit, the chat preference. | Ordinary chat mutation only after readiness. |
| `ready` | New server is healthy and authoritative. | Active runtime state. |
| `rollback` | Restore the prior model after a post-stop failure. | No new preference. |
| `unavailable` | No usable backend; do not claim loaded. | No false active selection. |

Only one transition may run at a time. `LoadModel` must not return success while the state is still `starting` or `health_check`; or it must return an accepted operation whose readiness is separately published. The frontend must not set `howl-last-model-path` or the chat preference until `ready`.

### 7.3 Failure recovery

| Failure point | Required result |
|---|---|
| Path/admission/eligibility/projector rejection | Old server and old selection unchanged. No stop. |
| Resource-fit rejection | Old server unchanged. No stop. |
| Stop failure | Do not start a new server; report failure; do not claim a new selection. |
| Spawn failure after stop | Attempt one restore of the previous model; if restore succeeds, old selection remains active. If restore fails, enter `unavailable`; never claim a loaded model. |
| Health timeout after spawn | Stop the failed process; attempt one restore; otherwise `unavailable`. |
| Concurrent transition | Reject or serialize; never run two load/unload sequences. |
| Active stream | Do not stop the server during generation. Cancel and report an aborted generation, or wait for completion. Do not claim the stream completed. |

The existing GGML crash retry may remain as an implementation detail, but it must not bypass the transition authority or publish readiness before health.

### 7.4 One runtime truth

`App.serverProc`, the HTTP client, `resource.Manager`, and `model.Manager` must agree after every transition. A failed transition must not leave a stale “loaded” entry in one authority and no process in another. The amendment requires reconciliation, not a second manager.

### 7.5 Resource scope

Keep the existing chat-versus-image arbitration (`app-vram-arbitration.go:47`, `:84`; `chat-model.js:330`). This amendment does not introduce multi-LLM residency or a general model orchestrator. If the target cannot fit, fail before stopping the old server.

### 7.6 No target runtime state

`loading`, `starting`, `ready`, `failed`, and `unavailable` are runtime/UI states. They are not written to `meta.json` or history. Only the explicit continuation preference is durable.

---

## 8. Media attachments and the F3-F boundary

### 8.1 The reference is not the bytes

An attachment reference proves only that the chat knows which media it used. It does not prove that the bytes exist. Durability requires:

1. a stable F3-F media identity;
2. bytes under an F3-F-owned media root;
3. an atomic commit and integrity record owned by F3-F;
4. a presence/corruption check when resolving the reference;
5. backup/recovery and deletion rules owned by F3-F.

F3-C must not copy image bytes into the chat directory, create a second media store, or use IndexedDB/localStorage as the canonical media authority.

### 8.2 F3-F prerequisite

No dedicated F3-F contract was found in the tree. The current gallery writer is not a substitute: it writes a PNG and a JSON sidecar with `writeFileAtomic`, and removes the PNG if the sidecar write fails (`internal/image/gallery/gallery.go:94-128`), but it has no registered media root, reference index, recovery journal, or F3-F ownership proof.

Therefore:

- F3-C may define the reference shape now;
- F3-C may not claim that chat attachments are durable until F3-F is accepted and proves ownership;
- Gate v2 attachment rows are **BLOCKED**, not PASS, until F3-F evidence exists;
- if the owner wants attachments to work before F3-F, that would create a second media owner and is rejected by this proposal.

### 8.3 Missing, corrupt, and legacy media

| State | Display/behavior |
|---|---|
| F3-F object present and hash matches | Attachment available. |
| F3-F object missing | Attachment unavailable; chat still opens; no substitution. |
| F3-F object corrupt/hash mismatch | Attachment unavailable/corrupt; no substitution; no silent re-encode. |
| Legacy `hasImage: true` with no reference | Legacy-unknown image; show unavailable; do not fabricate a reference. |
| `legacy-indexeddb` reference | Preserve provider; resolve only through the legacy path; migrate only under F3-F. |

No failed attachment may be replaced with a different image, a thumbnail, or a generated placeholder that claims to be the original.

### 8.4 Deletion, backup, and ownership

F3-C deletes chat targets and chat references; it does not delete F3-F bytes. F3-F owns media lifecycle:

- deletion of a chat does not automatically delete shared media;
- F3-F must define reference counting or ownership before deleting bytes;
- F3-C must not report media deletion success it did not perform;
- F3-F owns backup and recovery; a chat receipt does not prove media backup;
- wipe must respect the F3-F boundary; chat-root wipe is not an instruction to delete arbitrary media.

These are F3-F decisions, not F3-C target fields.

---

## 9. Interrupted generation and persistence

### 9.1 Required distinctions

The UI and persistence layer must distinguish:

- **committed** — target/receipt proves the message;
- **visible streaming** — text is displayed but not committed;
- **interrupted** — generation intent exists; no committed assistant message;
- **persistence failed** — generation completed but the commit failed; visible text is unsaved;
- **pending operation** — an ordinary journaled save/creation is recoverable;
- **aborted by user** — normal cancellation; no assistant message is invented.

### 9.2 User turn durability

Before the first model token, the user turn must be durably committed (or otherwise proven by an existing valid revision). A crash after sending must not silently lose the user's visible message. If the user turn cannot be committed, do not start generation.

### 9.3 Generation-intent journal

Add an operation artifact outside the eight targets, for example under the chat's F3-C operation area:

```text
f3c-generations/<F3C-PATH-GEN-1 hash>/<operationID>.json
```

Semantic fields:

```text
version, rootId, chatId, operationId, operationType,
baseRevision, userMessageId, userMessageHash,
modelId, startedAt, status
```

`status` is one of `started`, `streaming`, `completed`, `aborted`, `failed`, or `completion_pending`.

Rules:

- written create-only and flushed before the first token;
- not a target, not a revision, and not a receipt;
- not a substitute for commit proof;
- no partial tokens or secret prompt content are stored by default;
- a terminal record is retained until the next successful commit/recovery decision, then cleaned according to the operation inventory;
- malformed or conflicting intent blocks recovery rather than being guessed.

### 9.4 Completion and event order

The durable sequence becomes:

1. commit the user turn;
2. write the generation intent;
3. stream;
4. on normal completion, create the assistant candidate and commit a new revision;
5. update the intent to `completed`;
6. publish `chat:committed` (or an equivalent durable-success event).

A stream-complete event must not be the same signal as “committed.” If `chat:complete` is retained for compatibility, it means only that the stream ended; it must not imply durable success. A failed or unconfirmed commit publishes `chat:save_failed` or an equivalent unsaved state.

### 9.5 Crash recovery

On load, the recovery scan checks the intent:

- no intent → ordinary load;
- intent with committed assistant revision → mark completed/clean up;
- intent `started`/`streaming` with no assistant commit → mark **interrupted**; keep the committed user turn; do not restore partial tokens as a completed answer;
- intent `completion_pending` with no commit → mark **unsaved/persistence failed**; do not claim completion;
- `aborted` → no assistant message;
- malformed/conflicting intent → block the chat with a precise error.

The user may retry a new generation or discard the incomplete attempt. “Retry” is a new generation, not a claim to resume invisible tokens.

### 9.6 Power-loss boundary

This does not create a general power-loss guarantee. It creates an honest boundary: the last durable revision is authoritative; an interrupted or uncommitted response is clearly marked and never silently presented as complete. The current writer has no fsync/commit receipts, so this amendment is a specification requirement, not an implemented behavior.

---

## 10. Legacy compatibility and migration

### 10.1 v1 remains valid

`F3C-TARGET-1` remains a supported historical and current format:

- a v1 chat opens normally;
- `hasImage: true` without an attachment is legacy-unknown image;
- missing `generation` is legacy-unknown provenance;
- missing `continuationModel` means no explicit per-chat preference;
- the global last-model path is not written into the chat;
- a v1 target manifest and `tgt1-` hashes remain valid;
- old receipts are never rehashed or invalidated;
- no read-time migration rewrites files.

### 10.2 Upgrade on write

- The first ordinary write after the amendment creates a v2 target set and a v2 manifest.
- The new revision has a v1 predecessor; the chain remains valid.
- Legacy history entries are preserved; no model/seed/anchor is invented.
- `attachments: []` is written for v2 history entries that have no attachments; `hasImage` is preserved.
- `generation` is omitted for legacy-unknown entries.
- `continuationModel` is omitted unless the user explicitly selected one.

### 10.3 Save As

`Save As` from a v1 source creates a new v2 revision:

- source binding and baseRevision remain explicit;
- `legacy_unknown` history is preserved;
- no historical receipt is copied as v2 proof;
- the new source's target manifest is v2;
- no provenance is inferred from the current picker.

### 10.4 Receipts and revision coexistence

- A v1 current revision may coexist with v2 historical/current revisions in the same root.
- Verification dispatches by the manifest format on that revision.
- Historical proof checks hashes and links; it does not compare historical targets with current files.
- The receipt body/commit-record shapes do not change; the embedded target manifest carries the format.
- `F3C-CANON-1` remains unchanged.
- No old receipt becomes corrupt because a newer v2 revision exists.

---

## 11. Affected surfaces

| Surface | Current | Amendment impact |
|---|---|---|
| `meta.json` serializer | 19 fields; no model. | Append optional `continuationModel`; v2 field order. |
| `history.json` / `history.ndjson` | 7 fields; `hasImage` only. | Append `generation` and `attachments`; v2 exact bytes. |
| `ChatHistory` | No provenance/attachments. | New optional generation/attachments types; stable ID requirement remains. |
| `ChatMeta` | No continuation model. | Optional continuation preference; no runtime state. |
| Target manifest | `f3c-target-manifest-1`, `tgt1-`. | `f3c-target-manifest-2`, `tgt2-` for v2; v1 retained. |
| Receipt body/record/envelope | Fixed shape. | Shape unchanged; embedded manifest identifies v2. |
| Revision model | One target format. | Format per revision; v1→v2 ordinary write; no new revision kind. |
| Commit inventory | `f3c-commits`. | Unchanged. |
| Generation journal | None. | New `f3c-generations/.../<operationID>.json` operation artifact. |
| Operation map | Ordinary save/create. | Add generation operations; duplicate/crash recovery. |
| Locks | `chatWriteMu`, state locks, model runtime. | Add one model-transition serialization; do not hold target locks during model load. |
| Wipe | Chat root, gallery separate. | Include generation intents in chat cleanup; do not delete arbitrary F3-F media. |
| Media/F3-F | IndexedDB/gallery split; chat images transient. | References only; F3-F owns bytes; F3-F prerequisite for durability. |
| `LoadChat` | Restores context; no model. | Restore continuation preference; scan generation intents. |
| `LoadModel` | Returns before readiness; no rollback. | Readiness gate; rollback/unavailable; stream coordination; one runtime truth. |
| Frontend picker | Global last path; selection after return. | Per-chat restore; no global overwrite; commit only after ready. |
| Frontend stream | `chat:complete` before save. | Separate stream-end from committed; interrupted/unsaved states. |
| Test matrix | 17 Gate v2 rows. | Update V4/V5/V7/V11/V13/V15/V16; add V18–V23. |
| Freeze/authority | Exact contract hash. | New freeze and new Gate v2 binding if accepted. |

---

## 12. Gate v2 impact

### 12.1 Existing rows to revise

- **V4** — target bytes and manifest must cover `tgt1-`/`tgt2-`, v1/v2 dispatch, legacy history, and exact field order.
- **V5** — append, dirty state, autosave, session invalidation, stream-end vs committed, and interrupted generation.
- **V7** — Save As from v1 to v2, source binding, unknown provenance, no receipt reuse.
- **V10** — memory migration remains valid; media attachments are not a memory snapshot.
- **V11** — history migration preserves order/IDs and does not invent provenance.
- **V12** — wipe covers chat targets and generation intents without deleting F3-F media arbitrarily.
- **V13** — fault injection at generation-intent creation, assistant commit, target install, receipt install, and recovery.
- **V15** — model swap versus active generation/autosave; per-chat selection races.
- **V16** — native Wails/WebView2 path: load chat, restore selection, model unavailable, interrupted recovery.

### 12.2 New rows

| ID | Focus | Required evidence |
|---|---|---|
| V18 | Per-turn provenance | Two models in one chat; each turn records the actual model/engine/backend/parameters; removal leaves identity intact; no historical rewrite. |
| V19 | Seed semantics | Explicit zero, explicit nonzero, random with known effective seed, random without one, legacy-unknown; no `-1` conflation. |
| V20 | Attachments and F3-F | Reference survives reopen; F3-F present/hash matches; missing/corrupt is unavailable; no substitution; no bytes in target; old `hasImage` is legacy-unknown. **BLOCKED until F3-F proof.** |
| V21 | Continuation and hot swap | Three chats restore three selections; unavailable model is not substituted; failed swap preserves old or reports unavailable; no task-lock; active stream is not silently completed. |
| V22 | Interrupted generation | Kill after intent, during stream, and before commit; user turn survives; no fake assistant; retry/discard; malformed intent blocks. |
| V23 | v1/v2 compatibility | v1 opens unchanged; old `tgt1-` receipts valid; first write becomes v2; v1 predecessor chain valid; Save As v2; no rehash. |

Each row needs positive and negative evidence. V20 cannot be marked PASS by a chat reference test alone; it needs F3-F media ownership evidence. A missing F3-F contract is **BLOCKED**, not “passed by reference.”

### 12.3 Binding

If the owner accepts this amendment:

1. publish a new contract revision containing the v2 rules;
2. record its exact SHA-256;
3. update the freeze record;
4. re-bind Gate v2 to the new hash;
5. update the §11 matrix and row inventory;
6. rerun the Phase-0 assessment against the amended authority.

---

## 13. Decisions required from the owner

| ID | Decision | Proposal |
|---|---|---|
| OD1 | Adopt `F3C-TARGET-2`? | Yes, as a versioned target with v1 coexistence. |
| OD2 | Is F3-F a prerequisite for durable attachments? | Yes. Define references now; do not claim durable media without F3-F. |
| OD3 | Model identity scheme | Stable registry/digest identity; path is only a hint. |
| OD4 | Seed representation | Explicit/random by mode; zero valid; random effective seed optional; legacy by absence. |
| OD5 | Hot-swap failure policy | Preserve old if preflight fails; attempt one restore after post-stop failure; otherwise `unavailable`, never false loaded. |
| OD6 | Stream policy during swap | Do not stop an active stream; abort explicitly and mark aborted, or wait. No false completion. |
| OD7 | Interrupted generation | Durable user turn, generation intent, no partial authority; retry is a new generation. |
| OD8 | v1 lifetime | Allow v1 indefinitely; upgrade only on new write; never rehash old receipts. |
| OD9 | Save As | New v2 revision; preserve source binding and unknown provenance. |
| OD10 | Event semantics | Separate `chat:stream_end` from `chat:committed`; no durable success from a stream-end event. |

---

## 14. Status and authorization boundary

**Proposal status:** not authority; not accepted; not frozen.  
**Phase 0:** remains **NOT CLOSED** until the owner accepts or rejects this proposal, a new contract/freeze is published, and Gate v2 is re-bound.  
**P1a:** not authorized; remains downstream of P1c.  
**F3-C:** not accepted.  
**F3-D:** blocked.  
**A9/A10/A12/A13:** remain closed and are not reopened by this proposal.

No source changes, tests, builds, model loads, or implementation were performed. This document is the only artifact created by this task.

---

## Appendix A — exact references used

### Frozen authority

- `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`
  - §3.3 authoritative state, `:99`
  - §5.5 complete target set, `:260`
  - §5.6 exact target serialization, `:275`
  - §5.7 target hashes and manifest, `:330`
  - §5.8 compatibility decisions, `:352`
  - §6.1 receipt body, `:388`
  - §6.2 commit record, `:414`
  - §6.4 current commit proof, `:444`
  - §6.5 historical proof, `:455`
  - §6.9 temporary installation, `:500`
  - §7.3 standard durable mutation, `:575`
  - §8.1 creation, `:599`
  - §8.7 cast confirmation, `:728`
  - §8.9 history migration, `:777`
  - §12 known gaps, `:1018`

### F3-C / F3-F boundary

- `docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md:56-62` — image bytes under F3-F; F3-C never copies/moves/writes bytes; `media.json` and `legacy-indexeddb`.
- `docs/review/REVIEW-F3C-DECISIONS-AND-COMPATIBILITY-2026-11-10.md:105-116` — folder-per-item conflict; F3-F media identity and atomic commit.
- `docs/review/REVIEW-F3C-DECISIONS-AND-COMPATIBILITY-2026-11-10.md:154-162` — F3-F media validation; stable identity; legacy migration.
- `docs/review/REVIEW-F3C-DEPENDENCY-AND-ACCEPTANCE-2026-11-10.md:74` — F3-F partly blocks F3-C; references now, ownership later.
- `docs/review/ARCHITECTURE-F2C-CONTENT-OWNERSHIP-CUTOVER-2026-11-09.md:529` — F3-F gallery/PNG+sidecar scope.
- No dedicated `docs/**/*F3-F*` contract file was found.

### Chat persistence and model selection

- `internal/backend/chat/types.go:50` — `ChatMeta` has no continuation model.
- `internal/backend/chat/types.go:72` — `ChatHistory` has only `hasImage`.
- `internal/backend/chat/types.go:86` — `ChatData`.
- `internal/backend/types/message.go:31-43` — richer `Metadata` exists separately.
- `internal/backend/chat/persistence.go:192` — `SaveChat`.
- `internal/backend/chat/persistence.go:256` — `LoadChat`.
- `internal/backend/chat/persistence.go:585` — corrupt NDJSON truncation.
- `internal/backend/chat/persistence.go:668` — `AppendHistoryEntry`.
- `app-chat.go:96-111` — `LoadChat` warns and continues after failed save.
- `app-chat.go:298` — `SaveChatData` applies state before persistence.
- `app-chat.go:314-327` — `performSave`.
- `app-chat.go:355-381` — `recordUserTurn` only mutates memory.
- `app-chat.go:373` — index-derived IDs.
- `app-prompt.go:25` — streaming path.
- `app-prompt.go:58` — user turn recorded in memory.
- `app-prompt.go:138-146` — assistant append, `chat:complete`, async save.
- `app-prompt.go:281` — index-derived append IDs.
- `app-prompt.go:419` — `AbortStream`.
- `app-prompt.go:932` — global sampler settings.
- `frontend/dist/js/chat-model.js:166` — global auto-load.
- `frontend/dist/js/chat-model.js:169`, `:393` — global `howl-last-model-path`.
- `frontend/dist/js/chat-model.js:347-393` — selection after `LoadModel` returns.
- `frontend/dist/js/chat-input.js:20`, `:50`, `:240` — transient image/base64.
- `frontend/dist/js/chat-actions.js:62` — only `hasImage`.
- `frontend/dist/js/chat-session.js:146` — only `hasImage`.
- `frontend/dist/js/data-sync.js:238` — DOM index fallback ID.

### Hot-swap lifecycle

- `app-model.go:28-174` — preflight then stop.
- `app-model.go:174` — existing server stop.
- `app-model.go:317` — spawn.
- `app-model.go:345` — adopt.
- `app-model.go:345-512` — monitor/register/return before readiness.
- `app-model.go:537` — `UnloadModel`.
- `app-lifecycle.go:166` — `beginStream`.
- `app-lifecycle.go:193` — `adoptLlamaServer`.
- `app-lifecycle.go:221` — `stopLlamaServer`.
- `app-llama.go:158-205` — `spawnLlamaServerWithModel`.
- `app-vram-arbitration.go:47`, `:84` — chat/image fit and release.
- `internal/backend/model/manager.go:39` — manager unload-before-load path.

### Media

- `internal/image/gallery/gallery.go:38-128` — gallery save, PNG+sidecar, atomic file helper.
- `frontend/dist/js/storage-manager.js:40-84` — IndexedDB/localStorage image methods.
- `frontend/dist/js/storage-indexeddb.js` — `HOWLChatStorage`, `world-images`/`scenario-images`.
- `frontend/dist/js/chat-state-bridge.js:43-89` — legacy image resolution.

### Gate v2

- `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md:201-219` — V1–V17.
- `:220-245` — required row evidence and isolated destructive testing.
