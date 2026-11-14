# ASSESSMENT — F3-C Phase 0 Final Requirements Impact and Clarification Closure

**Date:** 2026-11-11  
**Status:** **PHASE 0 NOT CLOSED — CONTROLLED TARGET-SET AMENDMENT REQUIRED.**  
**Authority:** `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md`  
**Authority SHA-256:** `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31`  
**Verification protocol:** Gate v2, `docs/review/GATE-F3C-IMPLEMENTATION-VERIFICATION-V2-2026-11-11.md`.  
**Mode:** read-only source and document review; no source changes, tests, builds, model loads, or implementation.

The four remaining Phase-0 clarifications are resolved below as decisions. Phase 0 still cannot close because the frozen target set cannot preserve four product behaviors: per-turn model provenance, explicit-versus-random seed state, image/attachment references, and the per-chat continuation model. The smallest necessary amendment is proposed in §6; it is not applied to the frozen contract by this report.

Architecture freeze remains granted. P1a implementation remains unauthorized and downstream of P1c. F3-C acceptance is not granted. F3-D remains blocked. A9/A10/A12/A13 remain closed and are not reopened by this assessment.

---

## 1. Method and evidence boundaries

This assessment is static. It reads the frozen contract and the current working tree. It does not execute the application, load a model, run a test, build a binary, or verify a native Wails/WebView2 path.

The current tree is not a clean review packet: `git status` shows many modified files. This report therefore describes the current working tree, not a frozen implementation baseline, and it cannot serve as Gate v2 evidence.

Conclusions use three labels:

- **Covered:** the contract has an explicit requirement and the current source has a credible implementation path.
- **Partial:** either the contract or the source supports part of the requirement, but a material field or behavior is missing.
- **Gap:** the requirement is absent or the source contradicts it.

No conclusion below claims runtime correctness. A source path is evidence of design intent, not proof that the behavior works.

---

## 2. Executive conclusion

The frozen F3-C contract is strong for durable chat identity, revisions, receipts, target hashing, recovery, wipe, and text/context persistence. It is not sufficient for the user's stated trust requirement, **“If it's in your chat screen, it should come back,”** because the visible session contains data the target set does not represent.

The current source also does not implement the contract's durability protocol. That is expected: §12 records the implementation gaps. The important Phase-0 issue is narrower and more serious: even a correct implementation of the frozen target set would still lose per-turn model provenance, seed mode, image/attachment references, and the per-chat continuation model.

**Verdict:** Phase 0 is **not ready to close**. A controlled amendment to the target serialization, followed by a new freeze and Gate v2 authority binding, is required before implementation planning can treat the architecture as complete.

---

## 3. Requirements-impact assessment

### 3.1 Messages: text, identity, order, and metadata

| Requirement | Contract | Current source | Impact |
|---|---|---|---|
| Message text and order | §5.5 requires `history.json` and `history.ndjson`; §5.6 fixes their order and exact bytes. | `ChatHistory` stores text and role (`internal/backend/chat/types.go:72`); `history.ndjson` is appended and rewritten. | **Covered in contract; partial in source.** |
| Stable message identity | §8.9 requires order and IDs preserved; §12 records the current index-derived IDs. | `buildChatHistory` assigns `fmt.Sprintf("%d", len(history))` (`app-chat.go:373`); `appendChatMessage` assigns `len(a.chatMessages)-1` (`app-prompt.go:281`); the frontend falls back to `dataset.id || index.toString()` (`frontend/dist/js/data-sync.js:238`). | **Contract covered; implementation gap.** Index IDs are not stable across reloads or edits. |
| Per-message metadata | §3.3 lists only `chatMessages`, world, scenario, characters, and lorebooks; §5.6 fixes `ChatHistory` to `id, sender, text, role, timestamp, hasImage, characterId`. | `types.Message.Metadata` has model, temperature, top-p, top-k, images, audio, and documents (`internal/backend/types/message.go:31-43`), but the persisted `ChatHistory` does not carry them. | **Gap.** The richer type exists in a separate/legacy message path; it is not the persisted chat shape. |

**Consequence:** the target set can preserve message text, role, sender, timestamp, character owner, and an image-presence boolean. It cannot preserve model identity, generation settings, seed, token count, or attachment references. Stable IDs are a contract requirement but not a current implementation fact.

### 3.2 Character, world, scenario, and lorebook context

| Requirement | Contract | Current source | Impact |
|---|---|---|---|
| Active character/world/scenario identity | §3.3 and §5.5 require `characters.json`, optional `world.json`, optional `scenario.json`; `ChatMeta` carries stable IDs and tracking arrays. | `ChatMeta` and `ChatData` carry these fields (`internal/backend/chat/types.go:50`, `:86`); persistence writes and loads them (`internal/backend/chat/persistence.go:192`, `:256`). | **Covered in contract; source path exists.** |
| Chat-scoped lorebook selection | §5.5 requires `lorebooks.json`; §5.8 records it as required. | `applyChatDataToState` parses and applies `lorebooks` (`app-characters.go:355`, `:398`); `PersistenceManager` writes and loads `lorebooks.json`. | **Covered in contract; source path exists.** |
| Chat-owned generated memory snapshot | §5.5 and §5.8 require `memories_snapshot.json`; §8.8 defines migration and receipt binding. | `chat.Service` reads `memories_snapshot.json` (`internal/backend/chat/service.go:200`); no production writer was found; `PersistenceManager` writes `memories.json`, and the memory service also writes `memories.json` (`internal/backend/memory/service.go:345`, `:363`, `:452`). | **Contract covered; implementation gap.** A correct target-set implementation must add the snapshot writer; the current writer does not substitute it. |
| Generated memory completeness | §8.8 requires exact source/destination evidence and source preservation. | The legacy `memories.json` and the chat service's snapshot are separate artifacts; no migration evidence path exists in the current source. | **Gap; amendment not required for the target set, but implementation and migration are required.** |

The context and lorebook portions are structurally representable. The generated-memory portion is not currently written by the production persistence path, which is a known §12 gap rather than a new target-set defect.

### 3.3 Images and attachments

| Requirement | Contract | Current source | Impact |
|---|---|---|---|
| Image presence | §5.6 preserves `hasImage` as a boolean in `ChatHistory`. | `ChatHistory.HasImage` exists (`internal/backend/chat/types.go:72`); the frontend records only `msgDiv.dataset.hasImage` (`frontend/dist/js/chat-actions.js:62`; `frontend/dist/js/chat-session.js:146`). | **Partial only.** A boolean is not the image or a durable reference. |
| Image bytes or stable reference | The target set has no image/attachment target or field. | `chat-input.js` holds the selected image in a transient `currentImage` variable and sends base64 to the backend (`frontend/dist/js/chat-input.js:20`, `:50`, `:240`); it is cleared after send (`:80`). `types.Message.Metadata.Images` exists in the separate message type (`internal/backend/types/message.go:38-42`) but is not part of `ChatHistory`. | **Contract gap plus implementation gap.** Reloading a chat cannot restore the image, only the flag. |
| Audio/document attachments | Not represented in the target set. | Audio and file attachment handlers are explicit stubs (`frontend/dist/js/chat-input.js:86`, `:133`). | **Gap, but currently not an implemented visible feature.** Do not claim it as working. |
| Missing attachment behavior | The contract does not define unavailable attachment semantics. | No durable reference exists to test. | **Gap.** Any amendment must preserve a missing reference as unavailable, never substitute. |

This is the clearest violation of the user's trust requirement. The chat screen can show an image while the persisted chat retains only `hasImage: true`.

### 3.4 Per-generation model provenance

| Requirement | Contract | Current source | Impact |
|---|---|---|---|
| Actual model used for each response | No target field. §3.3 does not list model provenance; §5.6 fixes the history shape without it. | `types.Message.Metadata.Model` exists (`internal/backend/types/message.go:38`), but the persisted `ChatHistory` does not use it. The actual streaming path builds a prompt and calls the client without writing model identity (`app-prompt.go:25`, `:932`). | **Contract gap; implementation gap.** |
| Relevant generation configuration | No target field beyond message text. | Sampler and model settings are global maps (`app.go:54-55`), read into inference options from `a.samplerSettings` (`app-prompt.go:932`). | **Contract gap.** A later change to global settings retroactively changes the apparent configuration of old turns. |
| Historical model removal | No rule in the persistence contract. | No durable model identity exists. | **Gap.** The requirement that historical provenance never be rewritten or substituted is not representable. |

The separate `types.Message` type is not a solution: the user-facing chat persistence path uses `ChatHistory` and frontend maps. A type that exists only in an unused or legacy path does not preserve the visible chat.

### 3.5 Seed semantics

| Requirement | Contract | Current source | Impact |
|---|---|---|---|
| Explicit seed value | No target field. | `InferenceOptions.Seed` exists and defaults to `-1` (`internal/backend/llama/options.go:76`, `:156`); `buildInferenceOptions` reads a global `seed` (`app-prompt.go:932`). | **Contract gap; implementation gap.** |
| Distinguish random/unset from explicit | No target field. | `chat.Service.SaveChatData` writes `"seed": s.currentSeed()` with the comment “-1 means random, otherwise the pinned seed” (`internal/backend/chat/service.go:521`); `currentSeed` only returns the global inference option (`internal/backend/chat/seed.go:6`). There is no per-turn mode. | **Contract gap.** A persisted `-1` cannot distinguish “random” from “unset/legacy unknown,” and nothing records the actual seed used for a random generation. |
| Preserve explicit zero | No target field. | `setInt` can set `Seed` to zero from settings; no per-chat persistence. | **Gap.** Zero is a valid explicit seed and must not be treated as unset. |

The legacy `chat.json` seed is not a target file and is not per-turn. It cannot satisfy the requirement.

### 3.6 Continuation-model preference

| Requirement | Contract | Current source | Impact |
|---|---|---|---|
| Restore the chat's last-selected model | No field in `ChatMeta`; §5.6 fixes `meta.json` without a model field. | `ChatMeta` has no model field (`internal/backend/chat/types.go:50`); `LoadChat` restores messages/context, not model selection (`app-chat.go:96-126`). | **Contract gap; implementation gap.** |
| Do not lock the conversation to one model | No contract rule, but there is no model field to lock. | Model selection is global: `howl-last-model-path` is stored in `localStorage` (`frontend/dist/js/chat-model.js:169`, `:393`); `autoLoadDefaultModel` uses that global path (`:166`); model settings live in the unified global settings store. | **Partial behavior, wrong scope.** The current selection is not per chat. |
| Historical provenance survives a model switch | No target field. | No durable historical model identity. | **Gap.** |
| Model unavailable after removal | No target rule. | `autoLoadDefaultModel` shows a load error and “No model loaded”; it does not report a historical model as unavailable because it did not store one. | **Gap.** The requirement to preserve identity and report unavailability is not implemented. |

The user's intended experience — load chat, restore last-used selection, type, orchestrate loading — is not present. The current code has a manual model picker and a global last-used auto-load, not a per-chat continuation preference.

### 3.7 Hot-swapping and model orchestration

This is adjacent to F3-C, not part of its durable target set. It is still a product dependency because the continuation model and per-turn provenance must survive it.

| Requirement | Current source | Assessment |
|---|---|---|
| Switch models during an active conversation | `App.LoadModel` validates admission and projector, then stops an existing server and spawns a replacement (`app-model.go:28`); `model.Manager.LoadModel` unloads the current model before loading (`internal/backend/model/manager.go:39`). | **Partial path exists.** No runtime execution was performed, so this is not evidence that it works. |
| Preserve the old model on replacement failure | Admission and projector checks occur before the stop (`app-model.go:28-330`). | **Partial safety.** Rejection before stop preserves the old model. After the stop, a spawn or startup failure returns an error with no rollback; the old model is gone (`app-model.go:330-500`). |
| In-flight generation safety | `App.LoadModel` checks shutdown, not an active stream (`app-model.go:28`); there is no `streamMu` acquisition. | **Gap.** A model swap can stop the server while a stream is active. |
| Resource-aware orchestration | `CheckModelFit` and `FreeVRAMFor` handle chat versus image model contention (`app-vram-arbitration.go:47`, `:84`); the frontend calls them before some loads (`frontend/dist/js/chat-model.js:330`). | **Partial.** There is no multi-LLM residency or general model orchestration; the current design stops the existing chat server. |
| Subsequent generation uses the newly selected model | The direct app path uses the HTTP client and the newly spawned server. | **Not verified.** No native run was performed. |
| No conversation loss/reset | `LoadModel` does not clear chat state. | **No model-induced reset seen**, but stream interruption and save failures can still lose visible turns. |
| One model authority | The app has a direct server path (`App.LoadModel`) and a `model.Manager` path used by `chat.Service`; `App.LoadModel` registers with the resource manager, not necessarily the `modelManager`. | **Design gap.** The two lifecycle representations can disagree about “loaded.” |

The intended automatic experience is absent: `LoadChat` does not select or load a model. `autoLoadDefaultModel` runs on DOM ready using a global `localStorage` path. The user's manual Load button remains the effective path.

### 3.8 Interruption and recovery

| Requirement | Contract | Current source | Impact |
|---|---|---|---|
| Durable committed state | §3.4, §5.7, §6, §7, and §8 define revision, target hashes, commit records, receipts, and publication. | `PersistenceManager.AppendHistoryEntry` writes `history.ndjson` with `O_APPEND` and no sync (`internal/backend/chat/persistence.go:668`); `SaveChat` writes JSON via temp-and-rename but does not create a commit record or receipt (`:192`). | **Contract covered; implementation gap.** |
| Streaming is visible before durable | §8.5 defines token events as dirty state with no revision; a later complete append/autosave commits the full snapshot. | Streaming writes chunks to the frontend; the assistant turn is appended only after stream completion (`app-prompt.go:25-130`). The partial response is visible but not in history. | **Partial.** The visible partial can be lost on crash; the contract does not promise token-level durability. |
| Save success is not an in-memory claim | §8.3 and §8.5 require candidate/readback/commit/publication ordering. | `SaveChatData` applies state before persistence and then saves (`app-chat.go:298`); `SendMessage`/stream use a goroutine for `performSave` and do not await it (`app-prompt.go:138-146`); the frontend post-turn `SaveChat` is not awaited (`frontend/dist/js/chat-input.js:257-259`). | **Contract covered; implementation contradicts it.** |
| Failed outgoing save blocks a switch | §7.5 requires a failed outgoing save to abort. | `LoadChat` prints a warning and continues after a failed outgoing save (`app-chat.go:108-117`). | **Contract covered; implementation violates it.** |
| Corrupt history is reported | §5.7 and §8.9 require malformed/mismatched history to fail. | `loadHistoryFile` stops at a corrupt NDJSON line and returns the partial prefix (`internal/backend/chat/persistence.go:585`); `LoadChat` skips corrupt optional world/scenario files (`:256`). | **Contract covered; implementation silently degrades.** |
| Browser/WebView cache clearing survives | The durable root is the backend store; the frontend cache is not authority. | `LoadChatData` can fall back to persisted storage (`app-chat.go:206`; `app-characters.go:267`); the frontend also has a cache path (`frontend/dist/js/sidebar-session.js:246`). | **Partial.** Backend persistence can restore text/context if the save succeeded; a failed unacknowledged save can still lose the visible state. |
| Power-loss guarantee | The contract defines durable commit boundaries, not a blanket power-loss guarantee. | No fsync/durable commit protocol in the current writer. | **Not established.** Do not claim power-loss safety. |

The current code has a real “visible but not durable” window: a streamed response can be on screen while only the user turn is on disk, and the assistant turn exists in memory until completion. That is compatible with the contract's dirty-state concept, but the current implementation does not report or recover it.

---

## 4. Coverage table for the four clarified product behaviors

| Behavior | Contract coverage | Current source support | Classification | Phase-0 result |
|---|---|---|---|---|
| **Persistence of the visible chat** | Text, context, characters, world/scenario, lorebooks, and memory snapshot are required (§3.3, §5.5, §8.4, §8.8). Model/seed/attachments are absent. | Text/context paths exist; memory snapshot writer is missing; image is only `hasImage`; message IDs are index-derived. | **Contract defect for missing fields; implementation gap for the existing fields.** | **Not closed.** |
| **Model provenance and continuation selection** | No target field; §3.3 does not include model; §5.6 fixes the history shape. | `types.Message.Metadata.Model` exists in a separate path; persisted `ChatHistory` omits it; global `localStorage`/settings hold the last model. | **Contract defect; implementation gap.** | **Not closed.** |
| **Hot-swapping** | Outside the F3-C durable target set; F3-C only needs to preserve the resulting provenance/selection. | Manual load/unload path exists; no per-chat restore; no stream coordination; no rollback after stopping the old server. | **Implementation/design gap; no contract expansion for runtime state.** | **Not closed as a product behavior; not an F3-C target-set requirement by itself.** |
| **Interruption and recovery** | §3.4, §5.7, §6, §7, §8.3–8.5 define the intended commit/dirty/recovery model. | Current writer is best-effort NDJSON/JSON; no commit records/receipts; `LoadChat` continues after failed save; corrupt history truncates silently. | **Contract covered; implementation and evidence gaps.** | **Not closed.** |

---

## 5. Authoritative versus ephemeral state

The following are authoritative chat state and must survive a reopen:

- message text, order, role, sender, timestamp, stable message ID, character owner;
- active character/world/scenario identities and their pinned snapshots;
- chat-scoped lorebook selection and snapshot;
- chat-owned generated-memory snapshot;
- per-turn model identity and generation metadata;
- explicit seed value and seed mode;
- per-chat continuation-model preference;
- durable image/attachment references.

The following may remain ephemeral unless separately promoted:

- current model loading progress and load stage;
- the model dropdown's current in-memory selection before a commit;
- streaming token chunks and the visible partial assistant bubble before a complete append/autosave;
- reasoning/thinking text that is display-only and excluded from history;
- transient image previews and file-picker state before send;
- hover/focus/scroll state and other presentation controls.

The current target set preserves the first group only partially.

---

## 6. Proposed controlled amendment (not applied)

The smallest amendment is to extend the existing eight target files, not add a ninth target. This keeps the target set, path hashing, lock order, and receipt structure intact while versioning the target bytes.

### 6.1 Target format

- Introduce `F3C-TARGET-2` for the changed `meta.json` and history serialization.
- The target manifest must identify the target format; `f3c-target-manifest-2` is the cleanest version marker.
- Do not reuse `tgt1-` for bytes produced under different field rules; use a new hash/format marker or a versioned target manifest that makes the difference explicit.
- The eight target paths remain: `meta.json`, `characters.json`, `history.json`, `history.ndjson`, `lorebooks.json`, `memories_snapshot.json`, optional `world.json`, optional `scenario.json`.

### 6.2 Per-turn history fields

Add one optional, versioned generation object per `ChatHistory` entry. The field set is semantic, not a premature implementation schema:

- stable model identity for the model that produced the turn;
- optional display name; display name is not authority;
- generation parameters needed to explain the turn, at minimum the values the product exposes as relevant (for example temperature, top-p, top-k, and a template/config fingerprint);
- seed mode: `explicit`, `random`, or `legacy_unknown`;
- seed value only when mode is `explicit`; zero is valid; no `-1` sentinel;
- typed attachment references: stable IDs, kind, path/thumbnail if applicable, MIME type, and a content hash or equivalent integrity marker;
- no runtime state (`loading`, `loaded`, `unavailable`) in the durable target; runtime state is derived.

Historical entries with no generation object are `legacy_unknown`; no model, seed, or attachment may be invented.

### 6.3 Per-chat continuation preference

Add to `meta.json` a versioned `continuationModel` object:

- stable model identity;
- optional display name and source hint;
- last explicit selection time or equivalent ordering token;
- no lock and no promise that the model is installed.

The preference is restored as a selection only. If the model is unavailable, the chat reports it and does not substitute another model. Historical provenance is never rewritten because the preference changes.

### 6.4 Serialization and receipt impact

- `meta.json` and `history.json` field order changes; `history.ndjson` gains the generation/attachment fields.
- Exact target bytes and target hashes change; every new receipt's target manifest reflects the new bytes.
- The receipt shape remains; the target manifest content and hash domain/version change.
- Legacy chats remain readable as `legacy_unknown`; migration writes new fields only on a new commit. Old receipts remain historical and are not compared with current files.
- `memories_snapshot.json`, `lorebooks.json`, characters, world, scenario, and path hashes are otherwise unaffected.
- The consolidated contract hash changes; Gate v2 must be re-bound to the new hash. A verification against the old hash is invalid.

### 6.5 Acceptance impact

Add explicit Gate v2 coverage for:

- per-turn model provenance survives save, reopen, and model removal;
- continuation selection is per chat and is not a lock;
- explicit seed, zero seed, random seed, and legacy-unknown seed are distinguishable;
- image/attachment references survive save/reopen and missing media is reported as unavailable, never substituted;
- legacy chat migration does not invent provenance;
- target bytes/manifests/receipts for `F3C-TARGET-2`;
- no historical provenance rewrite after a model switch.

No A9/A10/A12/A13 result is reopened by the amendment. Those gates predate this target-set change and do not prove it.

---

## 7. Four Phase-0 clarification results

These are **closed as decisions**. They are not yet inserted into the frozen contract, so they remain pending controlled amendment.

### 7.1 Pristine-state determination includes a comprehensive canonical lore scan

**Decision:** A chat state is not pristine merely because its target files are absent or empty. Pristine determination must perform a read-only scan of every canonical lore source relevant to the chat — world, scenario, character, lorebook, and chat-owned generated-memory sources — and must distinguish `empty` from `unavailable`.

Rules:

- If any canonical root is unavailable, unreadable, or authority-indeterminate, the result is `unknown`, not pristine.
- If a canonical source is malformed, the result is `blocked`, not pristine.
- If a valid source contains user-generated lore, the state is not pristine.
- The scan is read-only; it creates no target, journal, receipt, or migration.
- No wipe, delete, import, or “fresh chat” decision may treat `unknown` or `blocked` as pristine.

Current source has no `pristine` term and no comprehensive scan of this kind. `worldstore` exposes unavailable/indeterminate outcomes (`internal/worldstore/store.go`), and `chat.Service.LoadChatData` loads a chat snapshot and memory service entries, but that is not a pristine-state decision. This is a **contract clarification plus implementation gap**, not an existing behavior.

**Status:** decision closed; contract insertion required.

### 7.2 Cast confirmation tokens use deterministic, versioned canonical hashing and same-process apply-time consistency

**Decision:** A cast confirmation token is a versioned, canonical object, not an opaque frontend string.

Minimum semantics:

- token version and operation ID;
- root ID, chat ID, base revision;
- candidate fingerprint and the target-manifest/candidate binding;
- cast scope;
- issue time and expiry;
- a canonical token hash under a domain-separated, versioned rule (for example `F3C-CAST-TOKEN-1`, with a `ct1-`-style marker);
- same-process issuance and application; no cross-process token acceptance;
- apply-time revalidation of the exact candidate, base, chat, scope, and token hash immediately before the write;
- a changed candidate, changed base, changed scope, malformed token, or expired uncommitted token fails without writing;
- a token that resolves to an already committed receipt reports that receipt and is never reapplied;
- single-use operation identity; a replay with a different candidate conflicts.

The frozen contract §8.7 currently says the token binds chat, base, candidate, scope, and expiry, but it does not define canonical hashing or apply-time consistency. Current source has no cast confirmation implementation. This is a **contract ambiguity/defect plus implementation gap**.

**Status:** decision closed; contract insertion required.

### 7.3 `AuthorityState.unavailable` semantics

**Decision:** `AuthorityState.unavailable` means the application could not complete an authoritative read or search of the relevant source. It does not mean “not found,” “empty,” or “deleted.”

Required behavior:

- do not inject the item into a prompt;
- do not render it as an ordinary absent item;
- preserve and display the last-known state as explicitly stale/unavailable where applicable;
- do not write, overwrite, delete, or migrate based on the unavailable result;
- do not substitute a name match, a different root, or a different model;
- block destructive or identity-changing operations until an explicit successful refresh establishes authority;
- fail closed on ambiguous multiple roots/copies.

This is distinct from:

- `not_set`: no reference was configured;
- `not_found`: a complete authoritative search found no item;
- `authority_indeterminate`: multiple/unavailable roots or copies prevent a single authority;
- `revision_mismatch`: the item exists but the expected revision is not the current one.

Current source has `ContextStatus` values `not_set`, `available`, `not_found`, `unavailable`, `authority_indeterminate`, and `revision_mismatch` (`internal/backend/chat/service.go:24-35`), but there is no `AuthorityState` type. The source's `ContextUnavailable` is close to the required behavior, but it is not the frozen contract's name or full semantics.

**Status:** decision closed; contract enum/insertion required. No target-set change.

### 7.4 Reference-expectation structure must validate before use

**Decision:** A reference-expectation structure must be parsed and validated as a closed, typed structure before any read/write decision. Malformed input fails the operation without writing.

Minimum validation:

- kind is one of the typed reference kinds;
- stable ID is required for an expectation; a display name or `legacyName` is not identity;
- required fields are present and correctly typed;
- empty IDs, duplicate references, invalid kinds/scopes, unknown fields, and mixed legacy/identity records are rejected;
- no silent normalization, drop, or “best effort” behavior;
- validation occurs before state mutation, target creation, journal advancement, receipt creation, or cast application;
- an unresolved but well-formed expectation follows the existing unavailable rules; a malformed expectation is a hard input error.

Current source defines `Reference` with `Kind`, `ID`, `Title`, and `LegacyName` (`internal/worldstore/refs.go:10-21`). `NormalizeReferences` silently skips malformed entries and converts name-only entries into legacy references (`:40`, `:76`); `ResolveReferences` reports unresolved references (`:121`). That is appropriate for a permissive legacy normalizer, but it does not satisfy the required fail-closed validation for a reference-expectation input.

**Status:** decision closed; contract insertion and implementation correction required. No target-set change.

---

## 8. Unresolved risks and limitations

1. **No runtime evidence.** This assessment did not load a model, run the app, exercise a stream, clear a WebView cache, or test a crash. All source conclusions are static.
2. **Working tree is not a frozen packet.** The source under review has uncommitted modifications. It cannot support Gate v2 verification or a reproducible acceptance claim.
3. **The amendment is not applied.** The proposed `F3C-TARGET-2` fields are a design recommendation; the frozen contract remains unchanged.
4. **Memory writer gap.** No production writer for `memories_snapshot.json` was found. The current `PersistenceManager` writes `memories.json`; the memory service also writes `memories.json`. This is a known implementation gap, not a target-set decision.
5. **ID instability.** Message IDs are derived from array length or DOM order. The contract requires stable IDs and explicit migration; current code does not provide them.
6. **Best-effort load behavior.** Corrupt history and optional context files can be skipped or truncated without surfacing a hard failure. That contradicts a trust-oriented persistence contract.
7. **Hot-swap is unverified.** The code has a manual replacement path, but no runtime evidence and no rollback after the old server is stopped. It must not be described as working.
8. **Two model authorities.** The app's direct server path and the `model.Manager` path can disagree about whether a model is loaded. This is outside F3-C's target set but affects the continuation/hot-swap experience.
9. **No power-loss claim.** The contract does not establish a blanket power-loss guarantee, and the current writer does not implement durable commit semantics.
10. **A9/A10/A12/A13 remain closed.** Nothing here reopens them; the evidence limitation in the freeze record still applies.

---

## 9. Final Phase-0 status

**Phase 0: NOT READY TO CLOSE.**

The four clarifications are closed as decisions, but the frozen target set cannot represent model provenance, explicit-versus-random seed state, image/attachment references, or the per-chat continuation model. Implementing the current contract faithfully would still fail the user's trust requirement for those visible parts of the chat.

Required before Phase 0 can close:

1. Accept or reject the proposed `F3C-TARGET-2` amendment in §6.
2. If accepted, publish a new contract revision and freeze record with a new exact SHA-256.
3. Re-bind Gate v2 to the new authority hash.
4. Add the four clarification decisions to the contract as explicit requirements.
5. Re-run the requirements-impact assessment against the amended contract.
6. Do not begin P1a; P1a remains downstream of P1c.

No source changes, tests, builds, or implementation were performed. P1a is not authorized. F3-C is not accepted. F3-D remains blocked.
