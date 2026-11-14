# R-07 — Evidence Report: EO-01, EO-02, EO-03
## Howl Chat / Amendment 02 Closure Follow-Through
### Date: 2026-11-11
### Status: Evidence gathering under read-only authorization; no architecture selected; submitted for independent review

## 0. Authority, scope, and evidence discipline

### 0.1 Immutable inputs

- Amendment 02 Revision 9: `073d0d4be768918841763a4ed714877ea8db2412053a1671fef2007a1327a640` — **Observed**
- Muse Revision 9 final review: `381fac1ce84e91f4e20d4649538598513f3f4de6e5bf0809bcb4d3574806a6e8` — **Observed**
- Amendment 02 revision-cycle closure: `f828fd0469cef0508e1e2c39482875ae90665bcc07d96821df293a9746863050` — **Observed**
- Closure verification record: `aa48d8a529ba0d6dd9d76639136e9e9b72d81ca670307598ec144676d0c3e1d3` — **Observed**

Revision 6-8 identities checked during evidence gathering:

- Revision 6: `ef2d5dd88fc30cf49424f6bcf20f50aca5c599ae6d9a8f76a33bc65673370433` — **Observed**
- Revision 7: `4cd760643d57f039ecb7fc11503608c8041b8e524b1aa5258e1ed8729fb42dd1` — **Observed**
- Revision 8: `7f97d3594f03cc88dc2ca9148345ad7b7336cdf4ce61f6831abe2478911890` — **Observed**

No Revision 10 file and no existing R-07 file were found before this report was created. The four reference files are byte-identical to the recorded hashes, but they appear as untracked `??` entries in `git status`; therefore Git status alone does not establish tracked-tree provenance. — **Observed**

### 0.2 What this report is and is not

This is an evidence artifact, not an implementation authorization, research-execution authorization, architecture approval, or Amendment 02 revision. It does not change source, schema, runtime behavior, models, benchmarks, networks, user data, or dependencies. — **Observed**

No runtime experiment, model load, inference test, benchmark, packet capture, migration, database mutation, or deployment was performed for R-07. Gaps that require those methods are recorded as unknowns. — **Observed**

### 0.3 Finding classifications

- **Observed:** directly inspected in source or existing records.
- **Source-derived:** supported by authoritative documentation or primary research.
- **Calculated:** derived transparently from observed inputs.
- **Inferred:** reasoned interpretation, not directly established.
- **Unestablished:** insufficient evidence.
- **Superseded:** previously held conclusion contradicted by newer evidence.

## 1. EO-01 — Cross-store commit and recovery

### 1.1 Current persistence topology

Howl has at least two different chat/memory persistence roots in the inspected code:

1. **Live chat persistence root**
   - `chat.NewPersistenceManager` uses `getAppDataDir()`.
   - Normal Windows production root is `%USERPROFILE%/AppData/Roaming/HOWLChat`.
   - Isolated native-test mode uses `StatePath("persisted-chats")`.
   - Chat files live under `<root>/chats/<chatID>/`.
   - Relevant code: `internal/backend/chat/persistence.go:15-75`. — **Observed**

2. **State-root service/memory persistence**
   - `memory.Service`, `MemoryStorage`, and `chat.Service` use `paths.StateRoot()`.
   - Memory files use `StateRoot/Chats/<chatID>/memories.json`.
   - `chat.Service` also reads `StateRoot/Chats/<chatName>/memories_snapshot.json`.
   - Relevant code:
     - `internal/backend/memory/service.go:27-58`
     - `internal/backend/memory/service.go:345-357`
     - `internal/backend/memory/service.go:361-380`
     - `internal/backend/memory/service.go:449-471`
     - `internal/backend/chat/service.go:79-110`
     - `internal/backend/chat/service.go:147-244`
     - `internal/backend/chat/service.go:200`
     - `internal/backend/chat/service.go:514-562` — **Observed**

These are not the same store. On Windows the names `chats` and `Chats` may be case-insensitive equivalents, but the parent roots differ: AppData legacy root versus project/state root. Treating them as one transactional store would be incorrect. — **Inferred**

### 1.2 Live chat files and write order

`PersistenceManager.SaveChat` writes, in order:

1. `meta.json`
2. `history.json`
3. optionally `world.json`
4. optionally `scenario.json`
5. `characters.json`
6. `lorebooks.json`
7. `memories.json`
8. `config.json` / active-chat pointer

Relevant code: `internal/backend/chat/persistence.go:191-251`. — **Observed**

Each JSON file is written through `saveJSON`, which writes `<file>.tmp` and renames it over the target. The comment states that this prevents a mid-write torn target file. Relevant code: `internal/backend/chat/persistence.go:569-581`. — **Observed**

There is no transaction spanning those files. If the process stops after `meta.json` but before `characters.json`, the next load sees a newer identity/metadata file with older context files. If it stops before `config.json`, the files may exist without becoming the active chat. No journal, commit marker, rollback, or cross-file recovery coordinator was found in the inspected persistence path. — **Observed**

### 1.3 History has two representations and a preferred source

The live chat store has:

- `history.ndjson`: append-only log; preferred at load when nonempty.
- `history.json`: snapshot/fallback.

`AppendHistoryEntry` appends one JSON line to `history.ndjson`, migrating existing `history.json` entries into NDJSON on first use. Relevant code: `internal/backend/chat/persistence.go:664-706`. — **Observed**

`loadHistoryFile` prefers nonempty NDJSON and falls back to `history.json`. If NDJSON has a corrupt line, it stops and returns the prefix; if scanner state reports an error, it returns a partial list. Missing/corrupt history otherwise yields an empty list rather than an error. Relevant code: `internal/backend/chat/persistence.go:583-629`. — **Observed**

`RewriteHistory` atomically replaces NDJSON through a temp file and then rewrites `history.json` to match. Relevant code: `internal/backend/chat/persistence.go:631-662`. — **Observed**

`SaveChat` rewrites `history.json` but does not rewrite NDJSON. Therefore a normal save does not reconcile the preferred NDJSON history source; only delete/edit and final-save paths rewrite NDJSON. — **Observed**

### 1.4 Normal turn/save path

The normal production flow inspected is:

- Frontend records the user turn through `AppendChatMessage`.
- `SendMessageStreamWithImage` de-duplicates/ensures the in-memory user turn.
- Backend builds the prompt and streams through `HTTPClient`.
- Assistant reply is appended through `AppendChatMessage`.
- Autosave calls `performSave`, which loads the active chat and overlays current app state.

Relevant code:

- `frontend/dist/js/chat-input.js:246-247`
- `app-chat.go:355-371`
- `app-prompt.go:25-149`
- `app-prompt.go:277-311`
- `app-chat.go:297-325` — **Observed**

`recordUserTurn` only appends to in-memory `a.chatMessages`; it does not itself persist to disk. Persistence happens through later append/save operations. — **Observed**

`applyChatDataToState` mutates backend app state before frontend conversion; if conversion fails, memory may already have changed. `performSave` loads the persisted chat first, then overlays memory and saves. There is no lock or transaction around load-modify-save. — **Observed**

### 1.5 Scenario and world ownership

Scenario selection has both frontend and backend effects:

- Frontend bridge updates local scenario state.
- Context selectors merge scenario characters/lore and call backend bindings.
- `AddScenarioToChat` sets backend scenario state.
- Full chat saves persist scenario/world snapshots through `SaveChatData`.

Relevant code:

- `frontend/dist/js/chat-state-bridge.js:346-374`
- `frontend/dist/js/chat-state-bridge.js:505-550`
- `frontend/dist/js/context-selectors.js:342-401`
- `frontend/dist/js/sidebar-context-scenario.js:172-200`
- `app-characters.go:248-260`
- `app-chat.go:297-308` — **Observed**

Scenario identity is represented by IDs/names in `ChatMeta` and scenario/world objects. No scenario revision, generation number, commit ID, operation ID, or deterministic replay marker was found in `ChatMeta`, `ChatHistory`, or the live `ChatData` persistence model. `chat.Service` can report canonical revision mismatch/status, but that status is not persisted as a chat commit marker. — **Observed**

`SaveChat` writes `world.json` and `scenario.json` only when the corresponding in-memory objects are non-nil; it does not delete stale world/scenario files when those objects are cleared. Therefore, clearing world/scenario state without deleting the files can leave last-known context on disk for a later best-effort load. This was established by code inspection, not by a deletion runtime test. — **Observed**

### 1.6 Memory-store split and no-op save

`memory.Service.SaveChatData` loads the existing `StateRoot/Chats/<chatID>/memories.json` and saves the same records back. Its comments explicitly say this only ensures the file exists and a full automated-lore sync is not implemented. It does not accept or use a new memory payload; the `chatName` and `syncAutomatedLore` parameters do not change the written records in the inspected implementation. — **Observed**

`MemoryStorage.saveMemories` overwrites `memories_backup.json` with `{}` on every save before writing `memories.json`. It is therefore not a useful prior-version backup in the inspected implementation. — **Observed**

`chat.Service.SaveChatData` writes `chat.json` and `history.json` with direct `os.WriteFile`, then calls memory save and optionally syncs character-scoped automated lore. The writes are sequential and non-transactional; the file writes in this secondary path do not use temp-file atomic replacement. — **Observed**

`chat.Service.SaveChatData` calls `memoryService.SaveChatData(chatName, ...)` without first setting the memory service's `chatID` in the inspected method. Since `MemoryStorage` uses `s.chatID`, a caller that has not set that scope may read/write an unintended or empty chat scope. Whether this path is live in production is separately unestablished. — **Observed**

### 1.7 EO-01 findings

| ID | Finding | Class |
|---|---|---|
| F-01-01 | Live chat files are written sequentially, not transactionally. | Observed |
| F-01-02 | Individual JSON files use temp-file rename, but the multi-file operation has no atomicity. | Observed |
| F-01-03 | NDJSON is the preferred history source, while normal `SaveChat` only rewrites `history.json`. | Observed |
| F-01-04 | `chat.Service`/memory state and live AppData chat state use different roots and formats. | Observed |
| F-01-05 | No journal, commit marker, event/operation ID, revision marker, or cross-store recovery coordinator was found in the inspected live persistence path. | Observed |
| F-01-06 | Best-effort load can skip corrupt/optional files and truncate NDJSON at the first bad line. | Observed |
| F-01-07 | A crash between files can leave newer metadata with older context, or persisted but inactive chat files. | Inferred |
| F-01-08 | Recovery cannot deterministically finish or reject a partial multi-file operation from current metadata alone. | Inferred |
| F-01-09 | Clearing world/scenario memory objects does not delete stale world/scenario files. | Observed |
| F-01-10 | Memory backup is overwritten with `{}` before every memory write. | Observed |

### 1.8 EO-01 candidate approaches

No candidate is selected. Each requires separate authorization for implementation or runtime validation.

#### A. Single transactional store

Use one local transactional store for chat identity, messages, scene state, lore selection, and memory metadata.

Advantages:

- One commit boundary for a turn or scene transition.
- Deterministic recovery: committed or not committed.
- Easier revision/operation identity.

Disadvantages/risks:

- Adds a new database dependency and migration path.
- Existing files, snapshots, canonical content stores, and frontend save payloads must be migrated or bridged.
- A database does not remove the need for careful crash and backup behavior.

Failure mode and detection:

- Failure: WAL/journal corruption, disk-full commit failure, lock contention, or incomplete migration leaving duplicate authorities.
- Detection: crash-injection tests, checksum/recovery tests, migration equivalence tests, startup integrity checks, and backup-restore tests.

Unresolved:

- Whether SQLite or another local engine is acceptable.
- Migration authority and rollback behavior.
- Performance on constrained local hardware.

#### B. File-generation plus atomic pointer swap

Write each new chat generation to a new directory, fsync files, then atomically switch an active-generation pointer/manifest.

Advantages:

- No new database dependency.
- Recovery can select the last complete generation or remain on the previous generation.
- Old generations can support audit/debug.

Disadvantages/risks:

- More storage churn and cleanup complexity.
- Windows file-replacement, antivirus locking, fsync, and crash semantics need explicit validation.
- Pointer/manifest itself becomes the critical single point.

Failure mode and detection:

- Failure: pointer updated to an incomplete generation, cleanup deleting the only good generation, or rename not durable across crash.
- Detection: generation-manifest hash checks, crash-injection tests, startup recovery tests, and orphan-generation audits.

#### C. Journaled two-phase intent over the existing files

Record an operation intent with IDs and expected hashes, apply file changes, then mark completion; recovery completes or discards intents.

Advantages:

- Preserves much of the existing file layout.
- Makes partial work explicit and recoverable.

Disadvantages/risks:

- More implementation complexity than either pure files or a database.
- The journal itself needs atomicity and recovery.
- Operation identity must be added everywhere.

Failure mode and detection:

- Failure: torn journal, duplicate replay, or intent marked complete before durable writes finish.
- Detection: journal checksum validation, idempotent replay tests, duplicate-operation tests, and crash-injection coverage.

#### D. Append-only event log with projections

Treat durable events as authoritative and derive snapshots for world, scenario, characters, lore selection, and memory.

Advantages:

- Clear ownership and replay history.
- Natural place for event IDs, operation IDs, and hashes.
- Supports audit and recovery.

Disadvantages/risks:

- Largest conceptual change.
- Snapshot/projection divergence and schema evolution risks.
- Existing save/load paths must be rebuilt or carefully bridged.

Failure mode and detection:

- Failure: projection bug, event-schema change, unbounded log growth, or compaction losing recoverability.
- Detection: replay-equivalence tests, projection regression tests, compaction tests, and storage-growth measurements.

## 2. EO-02 — Current authorization across retrieval paths

### 2.1 Retrieval and injection paths observed

At least these paths can supply stored information to a model:

1. Backend lore resolution for the generic prompt path.
2. Backend lore block for a character turn.
3. Frontend localStorage memory injection in multi-character composition.
4. `chat.Service` memory snapshot and RAG injection in the secondary service prompt path.
5. World, scenario, character cards, summaries, lorebooks, cached context, and frontend-composed framing.

Relevant code:

- `app-prompt.go:671-796`
- `app-prompt.go:798-924`
- `frontend/dist/js/chat-stream.js:100-180`
- `frontend/dist/js/memory-store.js:203-208`
- `internal/backend/chat/service.go:899-1025`
- `internal/backend/chat/service.go:1028-1050`
- `internal/backend/rag/service.go:119-169` — **Observed**

### 2.2 Where scope and audience are recorded

- Lore entries record `Scope`, `OwnerID`, `OwnerKind`, `Visibility`, `Enabled`, and `Stale`.
- `OwnerName` is descriptive metadata and is explicitly not authorization identity.
- Memory records record `Scope`, `OwnerID`, `OwnerKind`, status, confidence, and source.
- Backend memory records do not carry chat ID or source-message IDs in the inspected type.
- Frontend memory DTO exposes `SourceMessages`, but backend `types.MemoryRecord` does not preserve that provenance field.
- RAG documents record `Scope` and `ScopeID`, but not per-character/per-recipient audience.

Relevant code:

- `internal/backend/lorebook/lorebook.go:57-109`
- `internal/backend/types/memory.go:21-38`
- `app-memory.go:30-52`
- `internal/backend/rag/types.go:22-34` — **Observed**

### 2.3 Generic lore enforcement

`lorebook.Resolve` considers trigger matching, budgets, conflicts, and basic eligibility. `canConsider` excludes:

- disabled entries
- stale entries
- visibility `hidden`
- mismatched trigger direction
- already-triggered non-`always` entries

Relevant code: `internal/backend/lorebook/lorebook.go:134-188`, `249-268`. — **Observed**

`RouteResolvedEntries` then applies recipient routing:

- world entries are passed to all characters
- scenario entries are passed without comparing the entry owner to the active scenario name
- character entries require stable `OwnerID == characterID`
- unowned character entries pass only when explicitly shareable/public
- unknown scopes fall through and are included

Relevant code: `internal/backend/lorebook/lorebook.go:344-376`. — **Observed**

This contradicts any assumption that `RouteResolvedEntries` enforces current world/scenario authority. The function receives `scenarioName` and `worldName`, but the inspected implementation does not use them to filter world/scenario entries. It also does not use `ResolveRequest.ActiveWorld`, `ActiveScenario`, or `ActiveCharacter` for authorization. That is a material gap between parameter names/comments and implemented enforcement. — **Observed**

### 2.4 Canonical overlay and frontend-controlled pool

`SetLorebooks` parses frontend-supplied entries and overlays canonical content through `mergeCanonicalLore`. After cutover, canonical text can replace chat-snapshot text. Before cutover, or without a lorebook store, the pool is returned untouched and the browser store remains the source. Entries unknown to the canonical store retain frontend-supplied fields.

Relevant code:

- `app-prompt.go:433-449`
- `app-lorebooks-store.go:102-175` — **Observed**

`syncAppStateFromChatService` can replace the prompt pool with `chatService.GetAllLorebookEntries`, which combines manual and automated lore but not necessarily memory snapshots. Thus the authoritative prompt pool depends on which sync ran last. — **Observed**

### 2.5 Frontend memory injection

Multi-character composition calls backend `BuildCharacterLoreBlock`, then separately reads `window.MemoryStore.getCharacterMemories(character.id)` and injects up to four memories into the frontend-composed message. The frontend memory filter requires stable character ownership, but the injected material enters the backend through `GenerateCharacterResponse` as part of a user-content string.

Relevant code:

- `frontend/dist/js/chat-stream.js:110-180`
- `frontend/dist/js/memory-store.js:203-208`
- `app-prompt.go:154-266` — **Observed**

The backend cannot independently re-verify the scope, provenance, revocation state, or byte identity of already-composed frontend memory text from the user message alone. It can verify the character exists and route the request, but not reconstruct which stored memory records produced the injected prose. — **Inferred**

### 2.6 Secondary service path

`chat.Service.buildPrompt` can inject:

- resolved memory entries
- RAG chunks
- world/scenario/character context
- conversation history

Relevant code: `internal/backend/chat/service.go:899-1050`. — **Observed**

`rag.Service.Retrieve` performs embedding similarity, deduplication, and token-budget truncation, but does not perform per-recipient authorization in the inspected retrieval method. The query carries world/scenario/chat labels, but retrieval is scope-labeled rather than audience-authorized. — **Observed**

The live generic backend path inspected uses `App.buildPromptMessages`, not `chat.Service.buildPrompt`. Existing reviews also record this split. Whether the secondary `chat.Service.SendMessage` path is reachable in normal production remains **Unestablished** from the inspected application call sites alone.

### 2.7 EO-02 findings

| ID | Finding | Class |
|---|---|---|
| F-02-01 | There is no single shared authorization decision point across lore, frontend memory, service memory, and RAG paths. | Observed |
| F-02-02 | World and scenario lore are routed without comparing entry ownership to the active world/scenario. | Observed |
| F-02-03 | Character lore requires stable owner-ID match or explicit shareability. | Observed |
| F-02-04 | `hidden` lore is excluded; stale lore is excluded; private/unowned lore handling is fail-closed in the inspected router. | Observed |
| F-02-05 | Frontend memory injection bypasses backend memory authorization because it arrives as composed user text. | Observed |
| F-02-06 | Backend memory records lack durable source-message provenance. | Observed |
| F-02-07 | RAG retrieval is similarity- and budget-scoped, not audience-authorized. | Observed |
| F-02-08 | Prompt-pool authority depends on the last lore sync/merge path. | Observed |
| F-02-09 | Scope changes, promotion, and revocation effects on cached snapshots and derived artifacts are not centrally enforced. | Inferred |
| F-02-10 | Live reachability of the secondary `chat.Service` prompt/RAG path is not established. | Unestablished |

### 2.8 EO-02 candidate approaches

No candidate is selected.

#### A. Central context-authorization gate

Require every provider to emit claims containing entry ID, scope, owner, revision, visibility, source, and derivation. A central gate resolves those claims against current authoritative scope before assembly.

Advantages:

- One place to define deny-by-default behavior.
- Makes stale-artifact use explicit.
- Simplifies future auditing.

Disadvantages/risks:

- New coupling between providers and authorizer.
- Added latency and cache-invalidation complexity.
- Frontend-composed text still cannot be fully authorized unless the composer moves behind the gate or carries verifiable references.

Failure mode and detection:

- Failure: stale cache, revision lag, provider bypass, or frontend text injection outside the gate.
- Detection: cross-path matrix tests, canary entries, provenance assertions, and gate-bypass negative tests.

#### B. Capability-scoped retrieval per recipient

Each prompt path retrieves only records already scoped to the requesting chat/world/scenario/character.

Advantages:

- Less central coupling.
- Closer to data ownership.

Disadvantages/risks:

- Policy drift across providers.
- One missed path reintroduces bypass.
- Promotion/revocation must propagate to every cache.

Failure mode and detection:

- Failure: one retrieval path omits audience filtering.
- Detection: per-path authorization matrix tests and differential prompt-capture tests.

#### C. Audience partitioning at write time

Store audience-sensitive material in separate namespaces and prohibit cross-scope reads.

Advantages:

- Strong isolation by construction.
- Fewer runtime authorization decisions.

Disadvantages/risks:

- Duplication and migration complexity.
- Shared world/scenario material needs explicit sharing rules.
- Promotion/revocation requires sweeping copies, caches, snapshots, and derived artifacts.

Failure mode and detection:

- Failure: stale copies after promotion/revocation or accidental shared-namespace placement.
- Detection: revocation-sweep tests, namespace audits, and derived-artifact inventories.

#### D. Server-side verification for frontend memory

Keep localStorage as UI cache, but require backend memory IDs/owners/revisions before backend use of memory-derived material.

Advantages:

- Preserves UI responsiveness.
- Restores provenance for backend decisions.

Disadvantages/risks:

- Requires changing the multi-character composition contract.
- More backend/frontend coordination.
- Local-only threat model may not justify the cost.

Failure mode and detection:

- Failure: frontend and backend memory views diverge.
- Detection: differential frontend/backend memory assertions and unresolved-reference rejection tests.

## 3. EO-03 — Actual model-input boundary and independent verification

### 3.1 Generic request path

The generic live path is:

1. Stored chat, world, scenario, character, and lore state.
2. Backend lore retrieval/selection.
3. System-frame construction from world/characters/scenario/lore/rules.
4. Transcript assembly from in-memory chat messages.
5. Prompt-only final-user augmentation and context truncation.
6. `ChatCompletionRequest` construction.
7. JSON serialization.
8. HTTP POST to llama.cpp-compatible `/v1/chat/completions`.

Relevant code:

- `app-prompt.go:25-149`
- `app-prompt.go:671-796`
- `internal/backend/llama/http_client.go:594-653`
- `internal/backend/llama/http_client.go:144-180` — **Observed**

The request object is marshaled with `json.Marshal(req)` immediately before transport. No hash, byte log, receipt, signature, or independent copy is retained in the inspected production path. — **Observed**

### 3.2 Character request path

`GenerateCharacterResponse` accepts one frontend-composed string and sends it as a single `user` message. Backend character info is deliberately not prepended. Thus the actual model-input boundary for character turns includes:

- frontend-assembled character context
- backend lore block
- frontend memory text
- group context
- user message

all collapsed into one user-content string before backend serialization.

Relevant code:

- `frontend/dist/js/chat-stream.js:100-180`
- `frontend/dist/js/chat-stream.js:446`
- `app-prompt.go:154-266` — **Observed**

This means backend-only instrumentation cannot see the individual frontend sources except as undifferentiated user text. It can capture the composed string, but not its provenance. — **Inferred**

### 3.3 Image transformation

`GenerateChatStreamWithImage` replaces the last user message with OpenAI-style vision content containing text plus a base64 JPEG image URL. This is backend-specific request transformation after generic message assembly.

Relevant code: `internal/backend/llama/http_client.go:655-680`. — **Observed**

A receipt that only records pre-image `ChatMessage` values would not establish the transmitted vision payload. — **Inferred**

### 3.4 Existing receipt and capture seam

Production has no model-input receipt in the inspected path. The existing capture mechanisms are test-only:

- `promptcapture` builds a test-only recorder around the inference client.
- It records in-memory `[]llama.ChatMessage`.
- It does not serialize the request, hash bytes, or send the request to a model.
- Production builds replace it with a no-op wrapper.
- `GetPromptCapture` is available only in a `promptcapture` build.

Relevant code:

- `inference_capture.go:1-7`
- `inference_capture.go:47-77`
- `inference_capture.go:122-139`
- `inference_capture_off.go:1-6` — **Observed**

The unit-test `recordingInference` likewise records messages at the inference interface and returns a fixed reply without serialization or transport. Relevant code: `app_prompt_boundary_test.go:16-67`. — **Observed**

Therefore, existing tests can prove what message objects reached the inference interface, but they do not prove the transmitted HTTP bytes, image transformation, headers, URL, sampler fields, or backend-specific encoding. — **Inferred**

### 3.5 EO-03 findings

| ID | Finding | Class |
|---|---|---|
| F-03-01 | The inspected llama.cpp-compatible boundary is HTTP POST `/v1/chat/completions` with a JSON-marshaled `ChatCompletionRequest`. | Observed |
| F-03-02 | Production retains no hash, receipt, or byte copy of the transmitted request. | Observed |
| F-03-03 | Existing prompt capture records message objects, not serialized transport bytes. | Observed |
| F-03-04 | Existing capture is test-only and absent from production builds. | Observed |
| F-03-05 | Image requests transform the last user message after generic assembly. | Observed |
| F-03-06 | Character requests collapse multiple frontend sources into one user string. | Observed |
| F-03-07 | Exact byte identity at the transport boundary is not independently established by current code/tests. | Inferred |
| F-03-08 | Other production model backends and alternate invocation boundaries are not established from the inspected path. | Unestablished |

### 3.6 Independent verification method proposed, not executed

A future authorized investigation should capture the exact bytes handed to transport without relying on the same assembly object:

1. Instrument the HTTP transport layer to copy the outgoing request body before `Do`.
2. Record method, URL, headers necessary for interpretation, body length, and SHA-256 of the exact transmitted bytes.
3. Preserve the captured bytes separately from the application receipt.
4. Independently recompute the hash from the preserved bytes.
5. Compare the application-level assembled request manifest against the captured transport payload.
6. Cover both text and image paths.
7. Cover both generic and frontend-composed character paths.

This method is proposed only. It was not implemented or tested in R-07. — **Observed**

### 3.7 EO-03 candidate approaches

No candidate is selected.

#### A. Transport-layer tee

Hash and preserve exact request bytes at the HTTP transport boundary.

Advantages:

- Closest to what leaves the process.
- Independent of prompt-assembly objects.
- Covers sampler fields, headers, and image encoding.

Disadvantages/risks:

- Must cover every backend/client path.
- Adds sensitive-payload handling obligations.
- A proxy or transport wrapper could itself alter bytes if implemented incorrectly.

Failure mode and detection:

- Failure: uncovered backend path, transformed capture, or hash recorded after mutation.
- Detection: backend coverage matrix, byte-for-byte capture/receipt comparison, and no-transform proxy tests.

#### B. Loopback recording proxy

Point the client at a local recording proxy in an authorized test environment; forward unchanged bytes to llama-server and preserve the capture.

Advantages:

- External to application assembly.
- Can establish full transport bytes without changing production prompt code.
- Useful for image and multi-character cases.

Disadvantages/risks:

- Requires runtime/network/model authorization.
- Proxy behavior must be proven not to alter requests.
- Sensitive prompt data handling is required.

Failure mode and detection:

- Failure: proxy normalization, chunking, header changes, or TLS/body transformation.
- Detection: original/captured hash comparison and forwarding-equivalence tests.

#### C. Immutable request manifest plus independent encoder

Store an immutable manifest of system, messages, options, and image references; verify with a separate encoder/serializer and hash.

Advantages:

- Useful for regression and audit.
- Easier than full transport capture.

Disadvantages/risks:

- Still uses application-controlled data.
- Does not prove transport bytes unless compared against independent capture.
- Encoder differences can create false mismatches.

Failure mode and detection:

- Failure: manifest and transport diverge.
- Detection: manifest-vs-capture reconciliation tests.

## 4. Source and symbol inventory

### 4.1 EO-01 inventory

| Symbol/file | Relevance | Lines inspected |
|---|---|---|
| `chat.PersistenceManager` | Live chat persistence manager | `internal/backend/chat/persistence.go:15-75` |
| `NewPersistenceManagerWithDir` | Test-oriented arbitrary root | `internal/backend/chat/persistence.go:30-38` |
| `NewPersistenceManager` | Production root selection | `internal/backend/chat/persistence.go:40-58` |
| `getAppDataDir` | AppData versus isolated root | `internal/backend/chat/persistence.go:61-75` |
| `SaveChat` | Sequential multi-file save | `internal/backend/chat/persistence.go:191-251` |
| `LoadChat` | Best-effort load | `internal/backend/chat/persistence.go:253-324` |
| `loadHistoryFile` | NDJSON preference and truncation behavior | `internal/backend/chat/persistence.go:583-629` |
| `RewriteHistory` | NDJSON/history rewrite | `internal/backend/chat/persistence.go:631-662` |
| `AppendHistoryEntry` | Hot-path history append | `internal/backend/chat/persistence.go:664-706` |
| `saveJSON` | Per-file temp-file atomic write | `internal/backend/chat/persistence.go:569-581` |
| `ChatMeta` | IDs and tracking arrays; no revision/commit fields | `internal/backend/chat/types.go:41-69` |
| `ChatHistory` | Message/event fields; no commit marker | `internal/backend/chat/types.go:71-83` |
| `ChatData` | Saved chat aggregate | `internal/backend/chat/types.go:85-94` |
| `App.SaveChatData` | Frontend payload save | `app-chat.go:297-308` |
| `performSave` | Autosave load-overlay-save | `app-chat.go:314-325` |
| `performFinalSave` | Exit rewrite of NDJSON | `app-chat.go:327-353` |
| `recordUserTurn` | In-memory user-turn handling | `app-chat.go:355-371` |
| `updateChatDataFromAppState` | Memory-to-persisted overlay | `app-chat.go:417-498` |
| `convertFrontendDataToChatData` | Frontend-to-backend conversion | `app-chat.go:538-616` |
| `collectForSave` | Frontend save payload | `frontend/dist/js/chat-state-bridge.js:505-550` |
| `applyScenarioChange` | Scenario state transition | `frontend/dist/js/chat-state-bridge.js:346-374` |
| `ContextDisplay.prototype.selectScenario` | Scenario merge and backend sync | `frontend/dist/js/context-selectors.js:342-401` |
| `selectScenario` | Sidebar scenario path | `frontend/dist/js/sidebar-context-scenario.js:172-200` |
| `memory.Service.SaveChatData` | Existing-memory rewrite/no-op save | `internal/backend/memory/service.go:345-357` |
| `MemoryStorage.saveMemories` | Memory write and backup behavior | `internal/backend/memory/service.go:449-471` |
| `chat.Service.SaveChatData` | Secondary service save | `internal/backend/chat/service.go:514-562` |
| `chat.Service.LoadChatData` | Secondary service load | `internal/backend/chat/service.go:147-244` |

### 4.2 EO-02 inventory

| Symbol/file | Relevance | Lines inspected |
|---|---|---|
| `lorebook.Entry` | Scope, owner, visibility, stale fields | `internal/backend/lorebook/lorebook.go:57-109` |
| `lorebook.Resolve` | Trigger/budget/conflict selection | `internal/backend/lorebook/lorebook.go:134-188` |
| `canConsider` | Eligibility filter | `internal/backend/lorebook/lorebook.go:249-268` |
| `RouteResolvedEntries` | Recipient routing | `internal/backend/lorebook/lorebook.go:344-376` |
| `App.buildPromptMessages` | Generic prompt assembly | `app-prompt.go:671-796` |
| `App.resolveLore` / `resolveLoreBlock` | Lore prompt block assembly | `app-prompt.go:798-924` |
| `App.BuildCharacterLoreBlock` | Character lore resolver entry point | `app-prompt.go:853-877` |
| `App.SetLorebooks` | Frontend pool ingestion | `app-prompt.go:433-449` |
| `mergeCanonicalLore` | Canonical overlay | `app-lorebooks-store.go:102-175` |
| `GetAllLorebookEntries` | Secondary service lore pool | `internal/backend/chat/service.go:503-511` |
| `chat.Service.buildPrompt` | Secondary memory/RAG prompt path | `internal/backend/chat/service.go:899-1025` |
| `rag.Service.Retrieve` | Similarity retrieval and budgets | `internal/backend/rag/service.go:119-169` |
| `rag.Document` | Scope-labeled documents | `internal/backend/rag/types.go:22-34` |
| `MemoryStore.getCharacterMemories` | Frontend character memory filter | `frontend/dist/js/memory-store.js:203-208` |
| Character composer memory injection | Frontend memory into model text | `frontend/dist/js/chat-stream.js:100-180` |
| Backend memory DTO/type split | Provenance gap | `app-memory.go:30-52`; `internal/backend/types/memory.go:21-38` |

### 4.3 EO-03 inventory

| Symbol/file | Relevance | Lines inspected |
|---|---|---|
| `App.SendMessageStreamWithImage` | Generic inference path | `app-prompt.go:25-149` |
| `App.GenerateCharacterResponse` | Frontend-composed character path | `app-prompt.go:154-266` |
| `buildRequestFromMessages` | Request-object construction | `internal/backend/llama/http_client.go:594-639` |
| `GenerateChatStream` | Text streaming request | `internal/backend/llama/http_client.go:641-653` |
| `GenerateChatStreamWithImage` | Vision transformation | `internal/backend/llama/http_client.go:655-680` |
| `GenerateChatCompletionStreaming` | JSON marshal and HTTP POST | `internal/backend/llama/http_client.go:144-180` |
| `ChatCompletionRequest` | Transmitted request shape | `internal/backend/llama/http_client.go:65-98` |
| `promptCapture` | Test-only boundary recorder | `inference_capture.go:1-139` |
| Production no-op capture | Absence from production builds | `inference_capture_off.go:1-6` |
| `recordingInference` | Interface-level test recorder | `app_prompt_boundary_test.go:16-67` |

## 5. Data-flow diagrams

### 5.1 Persistence and recovery

```text
Frontend chat input/DOM
  -> AppendChatMessage
  -> in-memory App.chatMessages
  -> history.ndjson hot append
  -> performSave / SaveChatData
  -> SaveChat sequential JSON files
  -> config/active-chat pointer

Scenario selector
  -> local bridge/context state
  -> AddScenarioToChat backend state
  -> autosave/full save
  -> scenario/world snapshot files

Restart/load
  -> require meta.json
  -> prefer nonempty history.ndjson
  -> fallback to history.json
  -> optionally load world/scenario/characters/lore/memories
  -> apply tracking-array filter
```

There is no commit coordinator between these arrows. — **Inferred**

### 5.2 Authorization and retrieval

```text
Stored lore/memory/RAG
  -> provider-specific retrieval/selection
  -> lore: Resolve eligibility
  -> lore: RouteResolvedEntries recipient filter
  -> prompt-system block

Frontend MemoryStore
  -> character memory filter
  -> direct injection into frontend-composed user message
  -> GenerateCharacterResponse backend receives undifferentiated text

Secondary chat.Service path
  -> memory snapshot + RAG
  -> separate buildPrompt
  -> live reachability unestablished
```

There is no common authorization gate across these arrows. — **Inferred**

### 5.3 Model-input assembly and transport

```text
Generic path:
App state + lore + transcript
  -> buildPromptMessages
  -> []ChatMessage
  -> buildRequestFromMessages
  -> ChatCompletionRequest
  -> json.Marshal
  -> POST /v1/chat/completions

Character path:
Frontend character context + lore + memory + group context + user text
  -> one frontend-composed user string
  -> GenerateCharacterResponse
  -> []ChatMessage{user}
  -> buildRequestFromMessages
  -> json.Marshal
  -> POST /v1/chat/completions

Image variant:
  -> last user message replaced by text + base64 image content
```

No production receipt preserves the exact transmitted bytes. — **Observed**

## 6. Unknowns and evidence-gap register

| ID | Unknown/gap | Why it matters | Bounded investigation needed | Authorization needed |
|---|---|---|---|---|
| U-01 | Exact live reachability of `chat.Service.SendMessage`, secondary RAG, and memory-snapshot injection | Determines whether secondary authorization gaps affect production replies | Static call-graph completion plus isolated runtime path test | Runtime test authorization |
| U-02 | Full save timing after scenario/world transitions | Determines whether scene state can remain memory-only across crash | Isolated save-path trace and crash-window test | Runtime test authorization |
| U-03 | Frequency and impact of NDJSON/history.json divergence | Determines user-data loss/replay risk | Isolated history-divergence test | Runtime test authorization |
| U-04 | Behavior after partial multi-file writes, termination, restart | Required for EO-01 recovery design | Fault/crash-injection tests in isolated state root | Runtime experiment authorization |
| U-05 | Whether world/scenario clearing is intended to preserve files | Determines whether stale-file behavior is a bug | Requirements clarification plus isolated deletion test | Runtime test authorization |
| U-06 | Complete production model-backend inventory | Determines whether EO-03 has multiple transport boundaries | Static backend/client inventory plus configuration review | Read-only review; runtime invocation needs separate authorization |
| U-07 | Local hardware performance for SQLite, journaling, event log, or proxy capture | Prevents selecting an unsuitable mechanism | Benchmarks on representative hardware | Benchmark authorization |
| U-08 | RAG authorization requirements for local single-user deployment | Determines whether audience controls are product requirements or future multi-user requirements | Product/security decision plus threat-model review | Decision authorization, not code execution |
| U-09 | Memory revocation/promotion semantics | Determines whether derived artifacts must be swept | Isolated lifecycle tests with fixtures | Runtime test authorization |
| U-10 | Complete prior archive/sidecar enumeration | Needed for full closure provenance | Read-only filename/hash inventory | No new authorization if read-only |
| U-11 | Production user-data behavior | Cannot inspect user chats under current authorization | No user-data inspection; use synthetic fixtures only | User-data authorization if ever requested |

## 7. Bounded research plan

No research execution is authorized by R-07. The following phases require separate approval.

### Phase R-1 — Complete static inventory

- Finish source/symbol inventory for all retrieval, save, wipe, import, export, and backend invocation paths.
- Complete prior archive/sidecar enumeration.
- Produce a machine-checkable call graph and authority/root map.
- Authorization: read-only work already permitted; no new authorization expected.

### Phase R-2 — Isolated crash and recovery tests

- Use synthetic chats in an isolated state root.
- Test interruption between metadata/history/context/config writes.
- Test corrupt NDJSON, corrupt optional files, missing files, and stale world/scenario files.
- Measure whether recovery is deterministic.
- Authorization required: runtime experiments and isolated filesystem fault injection.

### Phase R-3 — Authorization-matrix tests

- Use synthetic lore, memory, RAG, and character fixtures.
- Test world/scenario/character/private/hidden/stale/legacy ownership cases.
- Test frontend memory injection and secondary service paths separately.
- Authorization required: isolated runtime tests; no user data.

### Phase R-4 — Transport-boundary capture tests

- Use a loopback proxy or transport tee in an authorized local test environment.
- Compare exact transmitted bytes against request manifests.
- Cover text, image, generic, and frontend-composed character paths.
- Authorization required: runtime, local model, and local network tests.

### Phase R-5 — Hardware and operational evaluation

- Measure candidate commit/capture mechanisms on representative constrained hardware.
- Evaluate dependency, backup, migration, and operational burden.
- Authorization required: benchmarks and implementation prototypes.

## 8. Automated lore — separate future research question

Automated lore is not approved as the future memory foundation by R-07. It is recorded here only as an early capture mechanism relevant to EO-02 and EO-03. — **Observed**

The separate future memory-system research question is:

> Should automated lore survive, be overhauled, be subsumed as one component of a broader memory architecture, or be removed?

The eventual research must compare alternatives and include measurable tests for at least:

1. Usefulness of currently generated entries.
2. Retrieval precision/recall with and without automated lore.
3. Continuity quality across chats/scenarios.
4. Noise, duplication, irrelevant detail, oversized entries, and cleanup burden.
5. Promotion, demotion, revocation, and audit behavior.
6. Processing, storage, maintenance, and local-hardware costs.
7. Whether equivalent or better continuity is available from manual lore, conversation history, summaries, RAG, character memory, or another capture method.
8. Whether overhaul improves extraction reliability enough to justify retention.
9. Whether removal harms measurable continuity or retrieval outcomes.

No redesign, replacement, removal, implementation, or lore experiment is authorized by R-07. — **Observed**

## 9. Completion checks

- [x] All four immutable reference hashes verified.
- [x] Revision 6-9 archive identities verified for the checked files.
- [x] No Revision 10 file found.
- [x] No R-07 file existed before this report.
- [x] R-07 work performed read-only; no prohibited runtime/model/benchmark/network/migration/dependency changes.
- [x] Every numbered finding carries an evidence classification.
- [x] Unknowns remain explicitly visible.
- [x] Proposed experiments are described as proposed, not completed.
- [x] Automated-lore research remains separate from the requirements baseline.
- [ ] Full prior archive/sidecar enumeration remains open and is recorded as U-10.
- [ ] Independent review remains pending.

## 10. Conclusion

The current implementation does not establish a cross-store commit, a shared retrieval authorization decision, or an independently verifiable model-input receipt. Those are evidence gaps, not merely documentation gaps.

Most importantly, the evidence contradicts any assumption that Revision 9’s proposed lore routing already enforces active world/scenario authority: world and scenario entries are passed through without comparing them to the active world/scenario. Likewise, the existing prompt-capture mechanism does not establish transmitted transport bytes.

No candidate architecture is selected or validated in R-07. The next step should be independent review, followed by separately authorized bounded research—not implementation.

**Requirements baseline closed. Architecture unvalidated. Evidence comes next.**
