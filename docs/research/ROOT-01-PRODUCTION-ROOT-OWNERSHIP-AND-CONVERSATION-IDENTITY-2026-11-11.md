# ROOT-01 — Production Root Ownership and Conversation Identity Investigation

**Project:** Howl Chat — Stone Wolf Systems  
**Work item:** ROOT-01  
**Priority:** 1  
**Date:** 2026-11-11  
**Authorization:** Read-only source and working-tree investigation. No implementation, migration, runtime test, user-data access, architecture selection, or Git mutation was authorized.

---

## 1. Executive findings

**Central answer: on the inspected source, production chat-owned data is not consistently identified or isolated by the established `(rootID, chatID)` contract.**

The decisive finding is not only that the application has several storage roots. It is that the production chat and memory types contain **no root identity at all**. `ChatMeta` has `ID` and display fields, but no `RootID`; `MemoryRecord` has `Scope`, `OwnerID`, and `OwnerKind`, but no root or conversation field; `AppConfig` stores only `LastChatID`. The persistence manager's `rootDir` is a process-local path string, not a root identity. A `(rootID, chatID)` tuple therefore cannot be formed from the production data model, and no code path passes one.

The second decisive finding is that chat identity is not unique. `generateChatID` formats UTC time to **one-second resolution** and includes no random or monotonic component. Two `NewChat` calls in the same second produce the same ID; `SaveChatAs` calls `NewChat` while the source chat may have been created in that same second, so "Save As" can target the source directory instead of creating a copy. `NewChat` does not reject an existing directory.

The third is that chat IDs are used as filesystem path components without validation. `SaveChat`, `LoadChat`, `DeleteChat`, `AppendHistoryEntry`, `RewriteHistory`, and `GetActiveChatMeta` all join a caller-supplied ID under `chats`. `SaveChatData` accepts a frontend-supplied `meta.id` and writes it. Relative traversal is therefore source-reachable; absolute-path and Windows alias behavior remain unestablished.

The fourth is that the memory service's conversation scope is a mutable, process-global field. `SetChatID("")` is a no-op, the initial scope is `"default"`, and every memory binding that receives an empty chat ID leaves the previous scope selected. Missing conversation identity therefore does not fail closed. Delete and wipe do not clear the scope, and backend wipe does not touch `StateRoot/Chats` at all.

Other source-confirmed findings include:

- `GetChatList` reads each chat's directory name but reports `meta.ID`; a missing or different metadata ID produces an unloadable or mispointed summary.
- `LoadChat` backfills a missing `meta.ID` from the directory name, while `GetChatList` and `GetActiveChatMeta` use different fallback behavior. Identity is reconstructed inconsistently.
- `SaveChatAs` copies `ChatData.Memories` but does not copy the memory-service store; the two stores both use the filename `memories.json` under different roots.
- `chat.Service` is keyed by `chatName`, not chat ID. Its startup call uses an empty name, so it reads root-level `StateRoot/Chats/memories_snapshot.json` and `StateRoot/Chats/memories.json`.
- Backend memory IDs are also generated at one-second resolution (`mem_YYYYMMDDHHMMSS`), so multiple memories created in one second share an ID.
- The frontend memory store is a single global `localStorage` array, `howl-chat-memories`, with no chat key in the record or lookup. Chat memories therefore mix across chats in the frontend store.
- `LoadChatData()` can return `meta.id = ""` when in-memory context exists; the bridge then falls back to `Date.now().toString()` on save, which the backend accepts as a new directory identity.
- Candidate and continuity storage resolve `os.Getwd()` rather than `paths.StateRoot()`. No production caller was found in this pass, so production reachability is unestablished, not disproved.

**What is conformant or partially conformant:**

- `MemoryRecord` has an explicit `OwnerKind`; `IsStableCharacterOwnership` requires `Scope == "character"`, `OwnerKind == "character_id"`, and a non-empty stable `OwnerID`, so legacy or name-owned character records fail closed at that filter.
- `history.ndjson` is treated as the canonical history log and `history.json` as a fallback snapshot, which is a coherent snapshot/canonical distinction inside the chat directory.
- World and scenario canonical loading explicitly uses stable IDs and refuses name-keyed folder fallback; characters and chat metadata do not apply the same rule.
- The isolation switch validates an operator-supplied root, requires a marker, checks containment, and rejects overlap with the real project. That protects the isolated test layout, but it is not a production root identity and does not apply to normal mode.

**What this report does not establish:** runtime behavior, packaged root resolution, AppData redirection, actual filesystem aliasing, actual Save As collision in a live session, actual cross-chat display leakage, candidate/continuity reachability through external callers, or concurrency outcomes. Those are listed in sections 13–15.

---

## 2. Scope and explicit exclusions

### 2.1 Included

- Production Go source for root resolution, chat persistence, chat service, memory service, candidate/continuity storage, character context loading, RAG, wipe, and Wails chat/memory bindings.
- Production frontend source under `frontend/dist/js` for session loading, save payload construction, memory storage, and fallback IDs.
- The authority and evidence artifacts named in section 3, read for requirements and prior findings.
- Current working-tree state and hashes, including tracked modifications and untracked files.

### 2.2 Excluded

- Any implementation, fix, refactor, schema change, migration, or architecture selection.
- Running the application, loading a model, starting a server, or executing a runtime test.
- Reading user chat data, databases, profiles, or persisted content.
- Import/export, recovery, deletion, wipe, reset, or move operations.
- Benchmarks, performance experiments, network requests, or external service calls.
- Git reset, stash, checkout, restore, stage, commit, or any other Git-state mutation.
- The full IDENTITY-03 investigation. Clear/reset behavior is examined only for its effect on root/chat identity and ownership; the full lifecycle contract is deferred.
- A14 execution, audio/TTS, and other queue items not implicated by a concrete chat identity path.

### 2.3 Search boundary

The source pass covered Go files under the repository, frontend JavaScript under `frontend/dist/js`, and the named documentation artifacts. It did not inspect compiled binaries, `node_modules`, external callers, Windows registry/known-folder state, or live profiles. **No absence claim in this report is proof of absence beyond that boundary.**

---

## 3. Input artifact identities and hashes

All hashes are SHA-256 of the current files. The two root-ownership reports are under `docs/review/`; the paths named in the request under `docs/research/` do not exist in this tree.

| Artifact | Role / authority | SHA-256 |
|---|---|---|
| `docs/review/FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md` | Frozen contracts; requirements | `aa11be0047c9410a72a37334d758dbd4d05dfb60cd318477cd5b103157df40ca` |
| `docs/review/IMPLEMENTATION-F3C-CHAT-CONSUMER-2026-11-10.md` | Earlier implementation assessment; observational, not current proof | `7ea9381817d5f8ae358395b4d9e4130adc2ed4ccfcc603669154faef3de31bbc` |
| `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md` | Architecture and implementation contract; intended design | `61042f9543ad5b44c22aacf22c94c1ef0157772471702c671977bb65f2725b31` |
| `docs/review/INVENTORY-F3C-ROOT-OWNERSHIP-2026-11-11.md` | Prior root inventory; source-derived evidence inventory | `343104b528818df699d752c419770fdc1bfd7fe71bc1f249e0dc6847f308c6b8` |
| `docs/review/INVESTIGATION-F3C-DROOT-CLOSURE-2026-11-11.md` | Prior D-ROOT investigation; source-derived evidence | `494539b9dde27a78c9c4c091131a93dce99c1debd81d11fea911e6d411e2e6a6` |
| `docs/review/DISPOSITION-REQUEST-F3C-DEFECTS-2026-11-11.md` | Open-defect disposition request | `ce32ac0bedbe216b55c065be073b2165209e0ebad4cb435897b829336d29d2c5` |
| `docs/research/GATE-B-SOURCE-INVESTIGATION-2026-11-11.md` | Gate B current-tree findings within its scope | `3ad665497a9f16fd75bcefaee3e2734f7495addc43c84c6767c5dd3ea7aca17f` |
| `docs/research/GATE-B-RECONCILIATION-AGAINST-EXISTING-ARTIFACTS-2026-11-11.md` | Gate B reconciliation | `c75483fab28b743349e1025bda03351795eb82cd714bf41d96ca2972914568db` |
| `docs/research/GATE-B-TO-F3C-DEFECT-WORK-QUEUE-2026-11-11.md` | Queue and dependency order | `d55aa7282059a48ab63679770357aade7db6fd1308578c34d45cf15e2ad5145c` |
| `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md` | Amendment 02 revision 9; architecture/research planning, not implementation authority | `073d0d4be768918841763a4ed714877ea8db2412053a1671fef2007a1327a640` |

The amendment's own header records revision 9 and states that it is architecture and research planning only, with no implementation, migration, runtime, or pairing authorization. That status is preserved here.

---

## 4. Current working-tree baseline and limitations

### 4.1 Baseline

- `HEAD`: `9c7d2a17dd7f702d7db6d95d9316e7640e34bf6f`
- `git status --porcelain=v1`: 258 entries
- Tracked modifications: 107
- Untracked entries: 151
- The tree is dirty. `HEAD` is **not** the reviewed implementation. Every source finding below is against the working tree identified by the hashes in section 4.3.

### 4.2 Relevant status

| File | Status |
|---|---|
| `app.go` | modified |
| `app-chat.go` | modified |
| `app-memory.go` | modified |
| `app-characters.go` | modified |
| `app-wipe.go` | modified |
| `internal/paths/paths.go` | unmodified (HEAD content) |
| `internal/paths/isolation.go` | untracked |
| `internal/backend/chat/persistence.go` | modified |
| `internal/backend/chat/types.go` | modified |
| `internal/backend/chat/service.go` | modified |
| `internal/backend/chat/loader.go` | modified |
| `internal/backend/memory/service.go` | modified |
| `internal/backend/memory/candidates.go` | modified |
| `internal/backend/memory/continuity.go` | modified |
| `internal/backend/memory/classify.go` | modified |
| `internal/backend/memory/identity.go` | untracked |
| `internal/backend/types/memory.go` | modified |
| `frontend/dist/js/chat-state-bridge.js` | modified |
| `frontend/dist/js/data-sync.js` | modified |
| `frontend/dist/js/sidebar-session.js` | modified |
| `frontend/dist/js/memory-store.js` | modified |
| `internal/backend/rag/*` | unmodified |

### 4.3 Principal source hashes (working tree)

| File | SHA-256 |
|---|---|
| `app.go` | `96d2e43ff908dde89e37d5fbac9912121d562412a76ffa1d2ea3af93dd2923b1` |
| `app-chat.go` | `f08037a8e3aaf32d5e70cad8fc25ecedefeddc3bd1a34a07e6a33d10256fea05` |
| `app-memory.go` | `7e7a8e83724f3c06e11a22193a1cb246ac89608669b51d15acdaaf1cbb427f60` |
| `app-characters.go` | `9d832100862b78c22bbe06c777e7eb911387d8a6988840b440ef0c24b5165a04` |
| `app-wipe.go` | `beb0b5b566e32433a91d4c4011812f2daf505837a7ba91ed24db0ae94a238d3e` |
| `internal/paths/paths.go` | `8d9bafbbe445e8a64e39a04a646606978f382515bbaf3a808f59f31b776716d3` |
| `internal/paths/isolation.go` | `c41dfde1e0c3b695d94081a477ad2c2d99faf9e80cf8f0a602752ecf9406a607` |
| `internal/backend/chat/persistence.go` | `ad9663a48278cff0ee35eb48d28799702cdbfa7754f9f71b7541a3a6762a87c9` |
| `internal/backend/chat/types.go` | `aa9f1f71adcc849bb77d8dd6c25b6be1ba2e1362d939206d00255f9b66fb1d46` |
| `internal/backend/chat/service.go` | `1f068fc15f0c89c6ee46a26a202e84f5dda505608c5acfe454f75c7585d4ae8c` |
| `internal/backend/memory/service.go` | `5f6c7d7824099e6999997afdd6b793bda422b360e73529c17a6f33e4f8837ed2` |
| `internal/backend/memory/candidates.go` | `fe72182aef09498b08f2aec33382599e78ae614b3b49ad21e750930a4a5750f7` |
| `internal/backend/memory/continuity.go` | `c4c9ce17e3e3d2fee7228d7f7c3c4253e224be961cdbd8cd8aede54d51a425e7` |
| `internal/backend/memory/identity.go` | `15c7c9ec37f839e708793e22a16f448d8d884f171ff058e69ac3987e20848ea3` |
| `internal/backend/types/memory.go` | `a76413f89a46d1249f6c705a05edd0d77c136d7357005ce4b87ea586f9053f50` |
| `frontend/dist/js/chat-state-bridge.js` | `3267e269c7d699b5a3f9c3b8d710a78aabb2986d27b5fa1060dc01e031ad2b2b` |
| `frontend/dist/js/data-sync.js` | `5708493201c9125873891749ad6c631614cd9fff700d6f8286bd272b6d0f497b` |
| `frontend/dist/js/sidebar-session.js` | `6f3547521a77e944f219b44207b16a72dabe1bf90be15c16d30217f3acb9011c` |
| `frontend/dist/js/memory-store.js` | `78a8f6b8c76d435ccdf4b58614b5692e7bad19d8c9cc6505436973d24176a4cd` |

### 4.4 Limitations

- Source-level observation does not prove runtime correctness. No application, model, database, or profile was started.
- The working tree is dirty and includes untracked identity/isolation files. Untracked code is not automatically production-reachable; its callers are traced separately.
- Prior reports do not name a source commit or working-tree hash binding. Their source-version applicability is therefore **not established**, even where their artifact hashes match the recorded values.
- Filesystem identity (volume/file ID, reparse points, 8.3 names, case aliases, mapped drives, UNC) was not examined at runtime. Path-string equality is not physical identity.

---

## 5. Root identity source trace

### 5.1 `paths.ProjectRoot` is a cached path string

`internal/paths/paths.go:29-35` returns `cachedRoot`, computing it once through `resolve()`.

`internal/paths/paths.go:37-59`:

1. Takes the executable directory.
2. Searches upward for `wails.json` or `go.mod`, up to `maxDepth = 6` (`:14-19`, `:62-76`).
3. If no marker is found from the executable, searches upward from the current working directory.
4. If cwd is available but has no marker, **returns cwd** (`:45-54`).
5. If cwd is unavailable, returns the executable directory, then `"."`.

This is a path resolver with a launch-directory fallback. It does not produce an opaque root ID and does not fail closed when no marker is found. A packaged launch from an arbitrary working directory can therefore use that directory as the state root. Whether the packaged shortcut sets such a working directory is a runtime question (unestablished).

### 5.2 `paths.StateRoot` selects one of two path bases

`internal/paths/isolation.go:185-190`:

```go
func StateRoot() string {
    if cfg := mustIsolation(); cfg != nil {
        return cfg.StateRoot
    }
    return ProjectRoot()
}
```

`StatePath` joins elements onto that base (`:192-195`). The isolated layout is `<operator-root>/data`, with WebView data under `<operator-root>/webview2` and temp under `<operator-root>/temp` (`:29-41`, `:122-127`). Isolation validates that the root is absolute, has the marker `.howl-isolated-root` with content `howl-isolated-test`, does not overlap the real project, and is writable (`:78-146`). That is a strong containment contract for isolated runs, but it is not a production root identity.

### 5.3 `PersistenceManager` uses a path, not a root ID

`internal/backend/chat/persistence.go:21-28` exposes `RootDir()` as a string. `NewPersistenceManager` (`:40-59`) calls `getAppDataDir()` and stores the result in `rootDir`. `getAppDataDir` (`:61-75`):

- If isolation is active: `paths.StatePath("persisted-chats")`.
- Otherwise: `filepath.Join(homeDir, "AppData", "Roaming", "HOWLChat")`.

The path is literal and profile-relative. It does not use `%APPDATA%`, a Windows known-folder API, or a stable volume/file identity. `NewPersistenceManagerWithDir` (`:30-38`) accepts an arbitrary directory and is documented as test-oriented.

### 5.4 No production type carries `RootID`

- `ChatMeta` (`internal/backend/chat/types.go:41-69`) has `ID`, `Name`, timestamps, world/scenario display names, character/lorebook display names, and tracking arrays. No `RootID`, root path, or root generation.
- `ChatData` (`:85-94`) contains meta, history, world, scenario, characters, lorebooks, and memories. No root.
- `AppConfig` (`:96-99`) contains only `LastChatID`.
- `ChatSummary` (`:101-112`) contains only `ID` and display/count fields.
- `MemoryRecord` (`internal/backend/types/memory.go:12-36`) has `ID`, `Scope`, `OwnerID`, `OwnerKind`, tags, status, and timestamps. No root and no conversation field.
- `MemoryCandidate` (`:78-95` in the same file) has scope/owner fields but no root or conversation field.

The only `RootID` hits in the repository are in `internal/lorebookstore/store.go` (`:116`, `:128`, `:142`, `:512`, `:520`, `:602`, `:644`, `:649`). Those belong to the content-library registry, not to chat or memory. They do not satisfy the chat contract.

### 5.5 Root selection lifecycle

`app.go:159-170` constructs the chat persistence manager and the memory service once at startup. The chat manager resolves the AppData (or isolated) path. The memory service resolves `paths.StateRoot()` independently. `chat.NewService` also resolves `paths.StateRoot()` at `internal/backend/chat/service.go:148`. There is no production API that switches roots during a chat's lifetime. Root change requires a process/environment change.

That does not make the root an identity. Two process instances, two launches from different working directories, a copied profile, or a path alias can present the same path string differently or the same location under different strings. No inspected resolver calls `filepath.EvalSymlinks`, a Windows final-path API, volume serial/file-ID lookup, or a reparse-point check. The prior inventory's alias finding remains open.

### 5.6 Root conformance result

| Question | Source-derived answer |
|---|---|
| Where does a root ID originate? | No root ID exists. A root is a path string returned by `ProjectRoot`, `StateRoot`, or `getAppDataDir`. |
| How is it validated? | Isolation roots are validated; normal roots are not identity-validated. `ProjectRoot` may fall back to cwd. |
| Is it passed explicitly? | No. It is captured inside `PersistenceManager`, `memory.Service`, `MemoryStorage`, `chat.Service`, candidates/continuity, and RAG. |
| Does missing/invalid root fail closed? | `NewPersistenceManager` returns an error; `NewApp` logs it and sets `chatPersistence = nil`. Memory and `chat.Service` can still be constructed. `ProjectRoot` does not fail closed on a missing marker. |
| Can root identity change during a chat? | Not through a production switch API. It can differ between processes, launches, or cwd changes for cwd-based stores. |
| Can separate roots resolve to the same persisted location? | Not established and not excluded. No canonical file-identity check exists. |

---

## 6. Chat identity source trace

### 6.1 Generation

`internal/backend/chat/persistence.go:138-144`:

```go
func (pm *PersistenceManager) generateChatID() string {
    now := time.Now().UTC()
    return fmt.Sprintf("%d-%02d-%02dT%02d-%02d-%02dZ",
        now.Year(), now.Month(), now.Day(),
        now.Hour(), now.Minute(), now.Second())
}
```

Resolution is one second. There is no random suffix, counter, UUID, or collision check.

### 6.2 Creation

`NewChat` (`:146-188`):

1. Generates the timestamp ID.
2. Joins `pm.rootDir/chats/<chatID>` (`:149`).
3. Calls `os.MkdirAll(chatDir, 0755)` (`:152`). `MkdirAll` succeeds if the directory already exists.
4. Builds an empty `ChatData` with that ID (`:157-175`).
5. Calls `SaveChat`.
6. Sets `config.LastChatID` and saves config (`:182-186`).

Two calls in one second collide. The second call does not receive an error from `MkdirAll`; `SaveChat` then writes empty meta/history over the first chat. This is a calculated consequence of the source, not a runtime observation. `NewChat` also does not validate the display name or the generated ID before using it as a path.

### 6.3 Save

`SaveChat` (`:191-251`):

- `chatDir := filepath.Join(pm.rootDir, "chats", chatData.Meta.ID)` (`:193`).
- `MkdirAll` (`:196`).
- Writes `meta.json`, `history.json`, `world.json`, `scenario.json`, `characters.json`, `lorebooks.json`, and `memories.json` (`:203-242`).
- Sets `config.LastChatID = chatData.Meta.ID` and saves config (`:244-248`).

There is no validation that `Meta.ID` is non-empty, is a single path segment, is not `.`/`..`, or belongs to this root. `filepath.Join` cleans the result. A relative ID containing `..` can escape `chats`; a crafted ID can write outside the intended directory. Absolute-path behavior differs by platform and was not tested. `SaveChat` is the common writer for `NewChat`, `SaveChatAs`, `performSave`, `RenameChat`, `SaveChatData`, and `saveChatDataToStorage`.

### 6.4 Load

`LoadChat` (`:253-272`):

- Joins `pm.rootDir/chats/<chatID>` (`:257`).
- Returns `"chat not found"` if the directory does not exist (`:260-262`).
- Loads `meta.json`; if `meta.ID` is empty, **sets it to the requested directory name** (`:270-272`).

This is an identity inference. A chat directory without an ID is made loadable by trusting the path. It does not validate that a non-empty `meta.ID` equals the directory name. A mismatched meta ID is accepted and then used by `SaveChat` as the write location on the next save, which can move the chat to a different directory or overwrite another chat.

### 6.5 List

`GetChatList` (`:485-536`):

- Reads each directory under `chats` (`:490-503`).
- Loads `meta.json` and skips invalid directories (`:505-511`).
- Builds `ChatSummary{ID: meta.ID, ...}` (`:521-531`).

The directory name is used to find the files, but the summary ID is the metadata ID. If `meta.ID` is empty, the summary ID is empty even though the directory name is a usable identity. If `meta.ID` differs from the directory name, the summary points at a different identity than the one just read. `GetChatList` has no fallback; `LoadChat` does. `discoverMostRecentChatSession` in `frontend/dist/js/sidebar-session.js` selects by `updatedAt`/`createdAt`, then loads by the returned `id`/`ID` (`:29-52`). A list entry with an empty or mismatched ID is therefore not reliably loadable.

### 6.6 Save As

`SaveChatAs` (`:436-468`):

1. Loads the source chat by `currentChatID` (`:439`).
2. Calls `NewChat(newName)` (`:445`).
3. Copies world/scenario names, history, world, scenario, characters, lorebooks, and `ChatData.Memories` (`:451-460`).
4. Saves the new chat (`:463`).

If the source chat was created in the same UTC second as the Save As call, `NewChat` generates the source ID. `NewChat` has already overwritten the source directory with an empty chat; `SaveChatAs` then copies the in-memory source data back into the same directory and saves. The result is not a copy: it is the source directory rewritten under the new name, with a new `CreatedAt` and the same ID. This is a calculated collision path. It is more likely immediately after chat creation or in rapid duplicate operations; no runtime test was run.

Even when the ID differs, `SaveChatAs` copies only `ChatData.Memories`, which is written to the AppData chat directory's `memories.json`. It does not copy the memory service's records under `StateRoot/Chats/<newID>/memories.json`. The backend wrapper then calls `memoryService.SetChatID(newID)` (`app-chat.go:80-88`), selecting an empty memory store. The copied chat therefore has one "memories" representation and an empty service store.

### 6.7 Delete

`DeleteChat` (`:470-483`):

- Joins `pm.rootDir/chats/<chatID>` (`:472`).
- Calls `os.RemoveAll` (`:474`).
- If the ID was `LastChatID`, clears the pointer and saves config (`:478-481`).

It does not validate the ID, does not remove `StateRoot/Chats/<chatID>` memory, does not remove candidate/continuity data, and does not clear or change the memory service scope. The App wrapper (`app-chat.go:116-130`) clears in-memory world/scenario/characters/lorebooks/messages only when the deleted chat is the active one; it does not call `memoryService.SetChatID`, `DeleteMemory`, or any cleanup on the other stores.

### 6.8 Active pointer

`GetLastChatID` and `SetLastChatID` (`:539-548`) read/write `AppConfig.LastChatID`. `GetActiveChatMeta` (`:550-565`) builds `pm.rootDir/chats/<LastChatID>/meta.json`, loads it, and backfills `meta.ID` from `LastChatID` if empty (`:561-563`). The active pointer is a global string with no root identity. It can be stale after a failed delete, wipe, or root change. `SetLastChatID` does not validate that the chat exists.

### 6.9 Wails and frontend identity

- `app-chat.go:28-33` passes `GetChatList` through.
- `NewChat` (`:44-58`) calls persistence, clears chat messages/triggered lore, and sets the memory scope to `chatData.Meta.ID`.
- `SaveChatAs` (`:72-89`) uses `LastChatID`, then sets memory scope to the returned ID.
- `LoadChat` (`:92-113`) optionally saves the current chat, loads the requested ID, updates app state, sets `LastChatID`, and sets memory scope.
- `SaveChatData` (`:298-308`) calls `applyChatDataToState`, converts the frontend map, requires non-empty `meta.id` (`app-chat.go:612-614`), then calls `SaveChat`.
- `performSave` (`:314-325`) saves `LastChatID` by loading it, updating from app state, and saving. It does not receive a chat ID argument.

`loadChatSessionById` in `frontend/dist/js/sidebar-session.js:156-235` calls `App.LoadChat(sessionId)` and uses the requested `sessionId` for `sidebarState.currentChatData`, `dataSync.currentChatId`, and the current-chat cache. It does not replace it with the returned meta ID. `hydrateFromBackend` in `chat-state-bridge.js:113` uses `meta.ID || meta.id || chatData.sessionId || null`. `collectForSave` (`:505-549`) uses `_state.sessionId || Date.now().toString()` at `:526`. `data-sync.js:170` uses `this.currentChatId || Date.now().toString()` in its own collector.

`LoadChatData` in `app-chat.go:207-295` can return `meta.id = ""` (`:261-272`) whenever in-memory world/scenario/characters/lorebooks exist. The frontend calls this as a fallback in `sidebar-session.js:128-143` and `context-data.js:119-124`. If the bridge receives that payload and the user saves, `collectForSave` can substitute `Date.now()` and `SaveChatData` will write a new timestamp directory. This is a source-reachable identity invention path. The exact user flow depends on frontend state and was not run.

---

## 7. Production data-root inventory

Paths below are logical, not proof of physical identity. "ID key" is the identifier used by the writer/reader, not necessarily an enforced ownership contract.

| Data category | Physical root | Logical owner | Identity key | Writer / reader | Load / legacy behavior | Scope risks | Evidence |
|---|---|---|---|---|---|---|---|
| Chat metadata, history, world/scenario/character/lorebook snapshots | Normal: `%USERPROFILE%\AppData\Roaming\HOWLChat\chats\<id>`; isolated: `StateRoot/persisted-chats/chats/<id>` | Chat | `ChatMeta.ID` (not root-scoped) | `SaveChat`, `AppendHistoryEntry`, `RewriteHistory`, `RenameChat` / `LoadChat`, `GetChatList`, `GetActiveChatMeta` | `LoadChat` backfills empty meta ID from directory; `GetChatList` does not; `history.json` fallback | No root ID; timestamp collisions; unvalidated path; metadata/directory mismatch; no transaction with memory | `persistence.go:138-251`, `:253-272`, `:436-536`, `:550-565` |
| Canonical history append log | Same chat directory; `history.ndjson` | Chat | Chat directory + message record | `AppendHistoryEntry`, `RewriteHistory` / `LoadChat` via `loadHistoryFile` | NDJSON is canonical; JSON is snapshot fallback | Message IDs are index-based in `app-chat.go:405`; not stable across rebuilds; no root | `persistence.go:651-706`; `app-chat.go:373-415` |
| Chat `memories.json` (frontend payload) | Same chat directory; `memories.json` | Chat (by directory) | `ChatMeta.ID` | `SaveChat` / `LoadChat` | Copied by `SaveChatAs`; unrelated to memory service records | Duplicate name with memory-service file; no root; content shape is `[]interface{}` | `persistence.go:239-242`; `app-chat.go:609-611` |
| Active chat pointer | AppData/isolated `config.json` (`AppConfig`) | Application/global | `LastChatID` | `NewChat`, `SaveChat`, `SetLastChatID` / `GetLastChatID`, `GetActiveChatMeta` | Empty if deleted/wipe; no root | Global pointer, no root or generation; stale pointer can fail or select another root's same ID | `types.go:96-99`; `persistence.go:539-565` |
| Memory records | `StateRoot/Chats/<memoryService.chatID>/memories.json` | Chat, by directory | `memoryService.chatID` + `MemoryRecord.ID` | `AddMemory`, `EditMemory`, `PinMemory`, `DeleteMemory`, `SaveChatData` / `GetMemories`, `GetMemoryDetail`, `GetMemoryIdentityDiagnostics` | Initial scope `"default"`; empty `SetChatID` ignored; no root/chat field in record | `default` shared scope; last-set scope; no record-level conversation ID; memory ID collision; no root | `memory/service.go:27-60`, `:98-125`, `:344-379`; `identity.go:73-80` |
| Memory snapshot | `StateRoot/Chats/<chatName>/memories_snapshot.json` | Chat, by name | `chatName` (not chat ID) | `chat.Service.SaveChatData` / `chat.Service.LoadChatData` | Startup uses `chatName=""`, yielding `StateRoot/Chats/memories_snapshot.json`; missing file -> empty, malformed -> error | Name-keyed; empty name maps to shared root; `SaveChatData` ignores its `chatName` argument and uses `s.chatID` | `chat/service.go:147-225`, `:545-548`; `app.go:170` |
| Continuity config | `cwd/Chats/<chatID>/continuity_config.json` | Chat (latent) | `chatID` | `CharacterContinuity.SaveContinuityConfig` / `LoadContinuityConfig` | No production caller found | cwd-dependent third root; no cleanup; path unvalidated; production reachability unestablished | `continuity.go:59-62`, `:154-187` |
| Candidate queue | `cwd/Chats/<chatID>/candidates.json` | Chat (latent) | `chatID` | `AddCandidates` / `GetCandidates` | No production caller found | cwd-dependent third root; no cleanup; path unvalidated; reachability unestablished | `candidates.go:22-35`, `:38-80`, `:84-120` |
| RAG documents/chunks | `StateRoot/RAG/<docID>.json`, `.txt` | Document; optional `ScopeChat` + `ScopeID` | `Document.ID`; `Document.Scope`/`ScopeID` | `rag.Storage.SaveDocument`, `SaveChunks` / `LoadDocument`, `ListDocuments`, `LoadIndex` | `Query.ChatID` is declared but not used elsewhere in Go; index loading takes scope/scopeID | No root ID; physical storage is a global doc-ID namespace; chat scope is a field, not a directory or enforced ownership | `rag/document.go:14-56`; `rag/types.go:3-33`, `:57-66` |
| Frontend memory store | Browser `localStorage` key `howl-chat-memories` | Global frontend array; records may claim `scope:chat` | Memory `id`; optional `ownerId`/`ownerKind` | `createMemoryEntry`, `updateMemoryEntry`, `deleteMemoryEntry` / `getMemoryEntries`, `getCharacterMemories` | No chat key in storage or lookup; backend sync uses `sessionId || 'default'` | Cross-chat display/leak; legacy records have empty `ownerKind`; no root | `memory-store.js:14`, `:57-118`, `:203-208` |
| Frontend current-chat cache | `HOWLStorage` / `localStorage` current-chat records | Chat | `currentChatId`/sessionId | `sidebar-session.js` save/load; `data-sync.js` | Falls back to discovery or generic `LoadChatData` | Cache can outlive backend chat; no root; fallback can create a new ID | `sidebar-session.js:128-143`, `:196-213`; `data-sync.js:48-78` |
| Character persistent data | `StateRoot/Characters/<charName>/character.json`, `character_lore.json`, `automated_character_lore.json` | Character | Folder name for loading; `Character.ID` inside; memory ownership uses stable ID | Content loaders / `chat.Service.LoadChatData` | Names are used to locate character files; stable ID is used for memory ownership | Display-name path; mixed identity model; no continuity ID; no root | `chat/service.go:160-197`; `entities.go:32-49`; `identity.go:20-25` |
| Gallery / generated images | `gallerySvc.BaseDir()` (app-owned) | App/gallery; chat association not established | Gallery record IDs | Gallery service / `WipeAllData` | Wipe treats gallery as app content; no chat ID coverage in backend wipe | Chat linkage not established; separate root; not cleared by chat delete | `app-wipe.go:41-70` |

---

## 8. Character-versus-chat ownership analysis

### 8.1 Stable character identity exists in the type and in the memory filter

`entities.Character` has a stable `ID` and a separate `Name` (`internal/backend/types/entities/entities.go:32-49`). `CharacterContext` embeds it (`internal/backend/chat/types.go:35-39`).

`IsStableCharacterOwnership` (`internal/backend/memory/identity.go:20-25`) requires:

- `characterID != ""`
- `Scope == "character"`
- `OwnerKind == types.OwnerKindCharacterID` (`"character_id"`)
- `OwnerID == characterID`

Legacy records with an empty `OwnerKind` or a display-name `OwnerID` fail this filter. `summarizeMemoryIdentities` counts those as unresolved and lists owners without reassigning them (`identity.go:38-70`). This is a conformant fail-closed rule at the memory prompt boundary.

### 8.2 Character loading still uses display names

`chat.Service.LoadChatData` builds `Characters/<charName>/character.json`, `character_lore.json`, and `automated_character_lore.json` from the supplied `characterNames` (`internal/backend/chat/service.go:160-197`). The comment at `:142-146` says world and scenario are addressed by stable ID and "the name parameters are display labels"; characters are not covered by that rule. The memory filter then uses `character.ID` from the loaded file (`:208-214`).

This is a mixed model:

- Character **memory** ownership is ID-based.
- Character **file resolution** is name-based.
- Chat **character selection** is stored both as display names (`ChatMeta.CharacterNames`) and stable tracking arrays (`ActiveCharacters`, `RemovedCharacters`, `AddedCharacters`) (`types.go:50-69`).
- `collectForSave` writes `characterNames` from display names and separate tracking arrays (`chat-state-bridge.js:529-542`).

A display-name change can therefore affect which character file is loaded even though the stable ID is unchanged. The exact rename behavior was not run.

### 8.3 No `(continuity_id, character_id)` or `(rootID, chatID, character_id)` contract exists

`CharacterContinuity` has only `projectRoot string` (`continuity.go:54-62`). It has no continuity ID field. Its export/import functions take `characterID` and a config, and filter memories by stable character ownership (`:64-151`). `ContinuityConfig` (`:32-43`) has no identity. There is no `ContinuityID` in the inspected chat, memory, or entity types.

`MemoryRecord` has no chat ID or root ID. The chat association is implied by the directory (`StateRoot/Chats/<memoryService.chatID>/`). Therefore the actual available keys are:

- `character_id` for character-owned memory content.
- The current `memoryService.chatID` for the directory.
- No `continuity_id`.
- No `rootID`.

`(continuity_id, character_id)` is **not implemented**. `(rootID, chatID, character_id)` is **not implemented**; only `chatID` and `characterID` exist, and the chat ID is not stored in the record.

### 8.4 Same character in multiple chats

Chat-specific data is separated by the chat directory and by the `memoryService.chatID` selection. Persistent character files are shared. That part is structurally plausible. However:

- `memoryService.chatID` is mutable and process-global; separate bindings can leave it at a previous chat.
- The frontend memory store is global, so character memories are visible across chats independent of the backend directory.
- Continuity export/import has no continuity ID and no production caller found.
- `SaveChatAs` changes the memory scope to the new ID without copying the service store.

Cross-chat isolation for the same character is therefore **not established**, and the global frontend store contradicts isolation by default.

### 8.5 Legacy character memories

The stable-ID filter is a positive control: legacy name-owned or unmarked character records are retained but excluded from prompts (`identity.go:20-35`; `chat/service.go:220-223`). This is the clearest conformant ownership rule in the inspected code. It does not solve chat/root identity, and it does not migrate or assign provenance.

---

## 9. Alternate and fallback path analysis

| Path | Source | Reachability | Ownership consequence |
|---|---|---|---|
| `SaveChatData` accepts frontend `meta.id` | `app-chat.go:298-308`; `:538-616` | Production Wails binding; frontend save paths call it (`data-sync.js:147-150`, `sidebar-session.js:449-455`) | Any non-empty frontend string becomes the chat directory ID. `NewChat` is bypassed. Missing ID fails (`:612-614`), but no path/segment validation. |
| `saveChatDataToStorage` duplicate writer | `app-characters.go:332-342` | Production helper; same `SaveChat` path | Same unvalidated ID and no root identity. |
| `LoadChatData` returns empty meta ID | `app-chat.go:207-278` | Production binding; called from `sidebar-session.js:132` and `context-data.js:121` | Frontend bridge sees `sessionId = null`; next save can use `Date.now()`. |
| Frontend `Date.now()` save ID | `chat-state-bridge.js:526`; `data-sync.js:170` | Production fallback | Creates a timestamp directory if session ID is absent; can establish a new active chat instead of updating the old one. |
| `GetChatList` summary ID | `persistence.go:521-531` | Production binding | Reports `meta.ID`, not the directory ID used to read; empty/mismatched metadata makes the list entry unusable or misdirected. |
| `LoadChat` meta backfill | `persistence.go:270-272` | Production load path | Infers ID from the directory. `GetChatList` and `GetActiveChatMeta` do not use the same fallback. |
| Memory service initial/default/last scope | `app.go:167`; `memory/service.go:33-38`; `app-memory.go:71-73`, `:96-98`, `:113-115`, `:136-138` | Production bindings | Empty chat ID is ignored; operations use `"default"` or the last-selected chat. |
| `chat.Service` empty-name startup | `app.go:170`; `chat/service.go:92-107`, `:200`; `:546` | Production startup | Reads root-level `Chats/memories_snapshot.json` and `Chats/memories.json`; not per-chat. |
| Candidate queue `cwd/Chats` | `candidates.go:22-35`, `:38-80` | No production App/Wails caller found in this pass; prior D-ROOT found none | Latent third root. If enabled, not covered by delete/wipe and not root-consistent. |
| Continuity `cwd/Chats` | `continuity.go:59-62`, `:154-187` | No production App/Wails caller found; prior D-ROOT found none | Latent third root; continuity config has no continuity ID. |
| RAG `StateRoot/RAG` | `rag/document.go:14-56`; `rag/types.go:3-33`; `app.go:184-187` | Production RAG service is constructed and attached | Global document namespace; `ScopeChat`/`ScopeID` are fields, not directories. `Query.ChatID` is declared but not used elsewhere in Go. Chat scoping is not enforced by storage identity. |
| Frontend global memory | `memory-store.js:14`, `:57-118`, `:203-208` | Production scripts loaded by `index.html:81` and other pages | All chats share `howl-chat-memories`; lookup has no chat filter. Backend sync chooses `sessionId || 'default'`. |
| Frontend current-chat cache | `sidebar-session.js:210-213`; `data-sync.js:48-78` | Production | Can restore a cached chat view without a backend identity; no root. |
| `WipeAllData` | `app-wipe.go:31-110`; `memory/service.go:33-38` | Production binding | Removes AppData `chats` and gallery only; memory/candidate/continuity/native roots remain; `SetChatID("")` is ignored. |

No alternate path was classified as test-only except `NewPersistenceManagerWithDir` and build-tagged native seams. No path was declared dead. Absence of a caller is a search-boundary result, not proof.

---

## 10. Clear/reset dependencies referred to IDENTITY-03

This section records only dependencies that affect root/chat identity or persistence. The full clear/reset contract belongs to IDENTITY-03.

### 10.1 Delete

`DeleteChat` removes the AppData/isolated chat directory and clears `LastChatID` only if it matches. It does not:

- clear or change `memoryService.chatID`;
- delete `StateRoot/Chats/<chatID>/memories.json`;
- delete `StateRoot/Chats/<chatID>/memories_snapshot.json`;
- delete `cwd/Chats/<chatID>/candidates.json` or `continuity_config.json`;
- clear RAG documents scoped to the chat;
- clear frontend `howl-chat-memories`;
- establish a tombstone or generation that prevents ID reuse.

### 10.2 Wipe

`WipeAllData` stages and removes `chatPersistence.RootDir()/chats` and the gallery (`app-wipe.go:31-110`). It resets `LastChatID` and clears in-memory chat state. At `:85` it calls `a.memoryService.SetChatID("")`, which returns immediately because empty is ignored (`memory/service.go:33-38`). It does not remove `StateRoot/Chats` or `cwd/Chats`.

### 10.3 Reactivation path

Because chat IDs are second-resolution timestamps, an ID can be reused. If a chat is deleted or wiped from the legacy store while `StateRoot/Chats/<id>` remains, a new chat created in the same second with the same ID will select the surviving memory directory. Whether that exact sequence occurs at runtime is not established, but the source permits it. The same applies to a frontend `Date.now()` ID that happens to match an existing memory directory.

### 10.4 Scope

- **Do not treat this section as the IDENTITY-03 investigation.**
- The blocking dependency is explicit: no clear/reset contract can be correct while `(rootID, chatID)` is absent and `SetChatID("")` cannot clear scope.
- The prior D-ROOT defect `DEF-F3C-IDENTITY-03` remains open.

---

## 11. Contract conformance matrix

Statuses are based only on inspected source. "Contradicted" means an inspected production path conflicts with the contract; it does not assert a runtime failure was observed.

| Contract | Status | Evidence | Why |
|---|---|---|---|
| Stable root identity | **Not implemented / contradicted** | No `RootID` in `ChatMeta`, `ChatData`, `AppConfig`, `ChatSummary`, or `MemoryRecord`; roots are path strings (`persistence.go:21-75`; `paths.go:29-59`; `isolation.go:185-190`) | There is no root identity to be stable. A path can change or alias without detection. |
| Stable chat identity | **Contradicted** | `generateChatID` at one-second resolution (`persistence.go:138-144`); no collision check (`:146-154`); frontend `Date.now()` fallback (`chat-state-bridge.js:526`; `data-sync.js:170`) | IDs are neither unique nor guaranteed stable; save can invent a new one. |
| `(rootID, chatID)` is authoritative | **Not implemented** | All persistence funcs take only `chatID`; `PersistenceManager` holds one root path; no call passes a tuple (`persistence.go:192-272`, `:436-483`) | The contract's key cannot be formed or enforced. |
| Stable character IDs | **Partially conformant** | `entities.Character.ID` (`entities.go:32-49`); `IsStableCharacterOwnership` (`identity.go:20-25`) | Memory ownership uses stable IDs; character file loading uses display-name paths (`chat/service.go:160-197`). |
| Character-persistent vs chat-specific separation | **Partially conformant** | Persistent `Characters/<name>` vs chat `characters.json`; memory service directory | Separation exists structurally, but memory records lack chat/root fields and the frontend store is global. |
| Conversation-local ownership by default | **Contradicted** | `"default"` initial scope (`app.go:167`); empty `SetChatID` no-op (`memory/service.go:33-38`); global `howl-chat-memories` (`memory-store.js:14`); `'default'` sync fallback (`:88`) | Missing identity does not fail closed; it reuses a shared/default/last scope. |
| Explicit promotion of chat-generated content | **Not established** | No promotion API traced; generated-lore sync has a character allowlist (`chat/service.go:550-570`) | The inspected paths do not establish a general promotion contract. This is not proof that none exists outside the boundary. |
| No invented provenance for legacy data | **Mixed / partially contradicted** | `LoadChat` infers ID from directory (`persistence.go:270-272`); `GetChatList` reports `meta.ID` (`:521-531`); character memory filter fails closed (`identity.go:20-25`) | Character memory handles legacy conservatively; chat identity invents/backfills an ID and applies it inconsistently. |
| Snapshot vs canonical store distinction | **Partially conformant** | `history.ndjson` canonical, `history.json` fallback (`persistence.go:208-213`, `:651-661`); memory snapshot separate (`chat/service.go:199-225`) | History distinction is explicit. Two different `memories.json` stores and name-keyed snapshots blur ownership. |
| Fail-closed on missing/foreign ownership | **Contradicted** | Empty `SetChatID` no-op; initial `"default"`; frontend `Date.now()`; `meta.ID` backfill; no ID validation | Multiple production paths continue with an inferred, default, or previous identity. |
| No accidental cross-chat/cross-root exposure | **Contradicted by source paths** | Global frontend memory; default/last memory scope; cwd roots; unvalidated IDs; delete/wipe gaps | Source permits cross-chat display and cross-root persistence. Runtime occurrence was not measured. |

---

## 12. Findings register

Each finding states what the evidence establishes and what it does not. Source-version applicability of prior artifacts is addressed in section 13.

### R1-F01 — No production root identity

- **Classification:** Observed / Source-derived.
- **Evidence:** `internal/backend/chat/types.go:41-99` (`ChatMeta`, `ChatData`, `AppConfig`, `ChatSummary`); `internal/backend/types/memory.go:12-36` (`MemoryRecord`); `internal/backend/chat/persistence.go:21-75` (`RootDir`, `NewPersistenceManager`, `getAppDataDir`); `internal/paths/paths.go:29-59`; `internal/paths/isolation.go:185-190`.
- **Caller/callee chain:** `app.go:159-170` → `chat.NewPersistenceManager` / `memory.NewService` / `chat.NewService` → path-based stores. No caller constructs a root ID.
- **Establishes:** The production types and persistence APIs have no root identity; root is a path captured inside each service.
- **Does not establish:** That two path strings cannot refer to the same physical data, or that a packaged build resolves a particular root.
- **Prior artifact:** Confirms the root-inventory central finding; extends it from "multiple roots" to "no root identity type."

### R1-F02 — Root selection can fall back to cwd and does not fail closed

- **Classification:** Observed / Source-derived.
- **Evidence:** `internal/paths/paths.go:37-59`; `internal/paths/paths.go:62-76`; `app.go:159-164`.
- **Caller/callee chain:** `paths.StateRoot` → `ProjectRoot` → `resolve` → `searchUp`; `NewApp` logs persistence failure and continues with `nil`.
- **Establishes:** A missing marker can yield cwd; a persistence failure does not stop memory/chat-service construction.
- **Does not establish:** Which cwd a packaged shortcut uses, or whether a real installation has markers.
- **Prior artifact:** Consistent with the inventory's packaged-root unknown.

### R1-F03 — Chat IDs are not unique

- **Classification:** Observed / Calculated.
- **Evidence:** `internal/backend/chat/persistence.go:138-144`; `:146-154`; `:157-175`.
- **Caller/callee chain:** `NewChat` → `generateChatID` → `MkdirAll` → `SaveChat`.
- **Establishes:** Two creations in the same UTC second produce the same ID; `MkdirAll` does not reject the existing directory; the second `SaveChat` overwrites.
- **Does not establish:** How often the frontend calls `NewChat` twice in one second in practice.
- **Prior artifact:** Not found in the prior inventory; this is a new source-derived finding.

### R1-F04 — `SaveChatAs` can overwrite the source instead of copying

- **Classification:** Calculated / Source-derived.
- **Evidence:** `internal/backend/chat/persistence.go:436-468`; `:138-154`; `app-chat.go:72-89`.
- **Caller/callee chain:** `App.SaveChatAs` → `PersistenceManager.SaveChatAs` → `LoadChat(current)` → `NewChat` → `SaveChat`.
- **Establishes:** If the source ID was generated in the same second, `NewChat` returns that ID and writes an empty chat to the source directory before the copy. The result is not a separate chat.
- **Does not establish:** Runtime frequency; the exact behavior if the source was created in a previous second.
- **Prior artifact:** The prior inventory did not identify this collision.

### R1-F05 — Chat IDs are unvalidated path components

- **Classification:** Observed / Source-derived.
- **Evidence:** `persistence.go:193`, `:257`, `:472`, `:669`; `:545-548`, `:556`; `app-chat.go:538-560`, `:612-614`.
- **Caller/callee chain:** Wails `LoadChat`/`DeleteChat`/`SaveChatData` → persistence join → `os.Stat`/`ReadFile`/`RemoveAll`/`MkdirAll`.
- **Establishes:** No function validates non-empty, single-segment, `.`/`..`, separator, or containment rules. Relative traversal can escape `chats`; `SaveChat` writes meta/history under the resulting path.
- **Does not establish:** Absolute-path handling on Windows, symlink/reparse behavior, or whether a caller can supply a crafted ID in the packaged UI.
- **Prior artifact:** The inventory found no canonical path identity; this adds the concrete join sites and the `SaveChatData` entry point.

### R1-F06 — Frontend can choose the chat directory identity

- **Classification:** Observed / Source-derived.
- **Evidence:** `app-chat.go:298-308`; `:538-616`; `frontend/dist/js/chat-state-bridge.js:505-549`; `frontend/dist/js/data-sync.js:147-150`, `:170`.
- **Caller/callee chain:** `collectForSave` → `App.SaveChatData` → `convertFrontendDataToChatData` → `SaveChat`.
- **Establishes:** `SaveChatData` requires only a non-empty `meta.id`, then writes it as the directory. It does not require the ID to be the active chat or to have been generated by `NewChat`.
- **Does not establish:** Whether the production UI exposes a direct payload editor; the path is reachable through the normal save bridge.
- **Prior artifact:** Not identified as an identity bypass in the prior inventory.

### R1-F07 — `GetChatList` reports metadata ID, not directory ID

- **Classification:** Observed.
- **Evidence:** `persistence.go:497-531`, especially `:502`, `:521-531`.
- **Caller/callee chain:** `App.GetChatList` → `GetChatList` → directory read → summary.
- **Establishes:** Empty `meta.ID` yields an empty summary ID; mismatched `meta.ID` yields a summary that points away from the directory just read.
- **Does not establish:** How often legacy metadata lacks an ID in real profiles.
- **Prior artifact:** Consistent with the prior catalog's directory-vs-metadata concern; this pins the exact fields.

### R1-F08 — Missing-ID fallback is inconsistent

- **Classification:** Observed / Source-derived.
- **Evidence:** `persistence.go:270-272` (`LoadChat` backfills); `:521-531` (`GetChatList` does not); `:561-563` (`GetActiveChatMeta` backfills from `LastChatID`).
- **Caller/callee chain:** Three readers of the same metadata apply different identity reconstruction rules.
- **Establishes:** The same on-disk chat can have different IDs depending on which reader loads it.
- **Does not establish:** Which reader the frontend uses in each flow.
- **Prior artifact:** Extends the D-ROOT identity finding.

### R1-F09 — Memory records have no root or conversation field

- **Classification:** Observed.
- **Evidence:** `internal/backend/types/memory.go:12-36`; `internal/backend/memory/service.go:344-379`; `internal/backend/memory/identity.go:73-80`.
- **Caller/callee chain:** `AddMemory`/`GetMemories` → `NewMemoryStorage(s.chatID)` → `Chats/<chatID>/memories.json`.
- **Establishes:** Chat ownership is represented only by the directory selected at call time. A record does not carry the chat or root it belongs to.
- **Does not establish:** That records are never copied between directories, or that no external index adds ownership.
- **Prior artifact:** Confirms the split-store finding; adds the record-level absence.

### R1-F10 — Empty chat ID does not clear or fail closed

- **Classification:** Observed / Source-derived.
- **Evidence:** `memory/service.go:33-38`; `app.go:167`; `app-memory.go:71-73`, `:96-98`, `:113-115`, `:136-138`; `memory/classify.go:91-95`; `memory-store.js:88`.
- **Caller/callee chain:** Wails memory calls → `SetChatID(chatID)` only if non-empty → `NewMemoryStorage(s.chatID)`.
- **Establishes:** Initial scope is `"default"`; empty input leaves the previous scope; classify falls back to `"default"`; frontend sync uses `'default'`. Missing identity is not an error.
- **Does not establish:** The exact interleaving of concurrent Wails calls; the mutable field is unsynchronized, but no race was run.
- **Prior artifact:** Confirms `DEF-F3C-IDENTITY-03`.

### R1-F11 — Delete and wipe leave memory and latent roots

- **Classification:** Observed / Source-derived.
- **Evidence:** `persistence.go:470-483`; `app-chat.go:116-130`; `app-wipe.go:31-110`, `:85`; `memory/service.go:33-38`.
- **Caller/callee chain:** `DeleteChat`/`WipeAllData` → legacy `chats` removal; memory service remains selected.
- **Establishes:** Backend delete/wipe covers AppData/isolated chat data and gallery, not `StateRoot/Chats`, `cwd/Chats`, or RAG. `SetChatID("")` is a no-op.
- **Does not establish:** Whether a user-visible "wipe" is expected to remove those stores; that is an owner/lifecycle decision.
- **Prior artifact:** Confirms `DEF-F3C-WIPE-02` and `DEF-F3C-IDENTITY-03`; adds the reuse scenario in section 10.3.

### R1-F12 — `SaveChatAs` copies the wrong memory representation

- **Classification:** Observed.
- **Evidence:** `persistence.go:451-460`; `:239-242`; `memory/service.go:344-357`; `app-chat.go:80-88`.
- **Caller/callee chain:** `SaveChatAs` copies `ChatData.Memories` into the AppData chat directory; `App.SaveChatAs` then sets the memory service to the new ID.
- **Establishes:** The copied `memories.json` is the chat-data field, not the memory-service store under `StateRoot/Chats/<id>`. The service store for the new ID is empty unless separately written.
- **Does not establish:** Whether the frontend later rehydrates the service store from the copied field.
- **Prior artifact:** Not identified in the prior inventory.

### R1-F13 — `chat.Service` is name-keyed and starts empty

- **Classification:** Observed / Source-derived.
- **Evidence:** `chat/service.go:92-107`; `:147-225`; `:545-548`; `app.go:170`.
- **Caller/callee chain:** `chat.NewService` → `memory.NewService(chatName)` → `LoadChatData(..., "")` → `Chats/memories_snapshot.json`; `SaveChatData` ignores `chatName` and uses `s.chatID`.
- **Establishes:** Startup uses empty `chatName`; snapshot path is root-level; `SaveChatData`'s `chatName` parameter is unused; memory service scope is empty at startup.
- **Does not establish:** Whether the service save path is production-reachable; prior D-ROOT found it test-seam-only.
- **Prior artifact:** Confirms `DEF-F3C-STARTUP-05` and the startup read path.

### R1-F14 — Backend memory IDs collide at one-second resolution

- **Classification:** Observed / Calculated.
- **Evidence:** `internal/backend/types/memory.go:107-110`; `memory/service.go:98-125`.
- **Caller/callee chain:** `AddMemory` with empty `record.ID` → `NewMemoryID()` → append → save.
- **Establishes:** Two memories created in one second receive the same `mem_...` ID. `EditMemory`/`DeleteMemory` search by ID and act on the first match.
- **Does not establish:** Whether the frontend always supplies its own random memory ID; the frontend generator uses `Date.now()` plus random (`memory-store.js:31-33`).
- **Prior artifact:** Not identified in the prior inventory.

### R1-F15 — Frontend memory store is global, not chat-scoped

- **Classification:** Observed / Source-derived.
- **Evidence:** `frontend/dist/js/memory-store.js:14`, `:57-118`, `:203-208`; `frontend/dist/index.html:81`; `data-sync.js:538-550`.
- **Caller/callee chain:** `createMemoryEntry` → `_saveMemories` (global key) → `getMemoryEntries` (no chat filter) → UI.
- **Establishes:** All chats share one array. `scope:'chat'` and `ownerId` are fields, not storage partitions. `getMemoryEntries` filters only scope/query/status. Backend sync uses `sessionId || 'default'`.
- **Does not establish:** Whether every UI surface displays all scopes; the store itself does not partition by chat.
- **Prior artifact:** Not identified in the prior inventory.

### R1-F16 — Frontend save can invent a chat ID

- **Classification:** Observed / Source-derived.
- **Evidence:** `app-chat.go:207-278` (`LoadChatData` meta ID `""`); `chat-state-bridge.js:113`, `:526`; `data-sync.js:170`; `app-chat.go:612-614`.
- **Caller/callee chain:** `LoadChatData` → `hydrateFromBackend` (sessionId null) → `collectForSave` (`Date.now`) → `SaveChatData` (non-empty requirement passes) → `SaveChat`.
- **Establishes:** A source-reachable path creates a timestamp directory instead of updating the active chat.
- **Does not establish:** Whether the fallback is taken in the normal startup flow; `LoadChatData` is used as a fallback in `sidebar-session.js:132` and `context-data.js:121`.
- **Prior artifact:** The prior D-ROOT report noted empty ID behavior; this report traces the frontend consequence.

### R1-F17 — Candidate/continuity roots depend on cwd

- **Classification:** Observed / Unestablished reachability.
- **Evidence:** `candidates.go:22-35`, `:38-80`; `continuity.go:59-62`, `:154-187`.
- **Caller/callee chain:** `ProcessMessageForEvents` → `AddCandidates` → `NewCandidateStorage`; continuity service methods → `NewCharacterContinuity`.
- **Establishes:** Both use `os.Getwd()` and write `Chats/<chatID>` beneath it. That root is independent of `StateRoot` and AppData.
- **Does not establish:** Production reachability. No App binding, generated Wails method, indirect caller, or startup registration was found in this pass or by prior D-ROOT.
- **Prior artifact:** Confirms `DEF-F3C-REACH-04`; does not close it.

### R1-F18 — Character identity is mixed and continuity identity is absent

- **Classification:** Observed / Source-derived.
- **Evidence:** `entities.go:32-49`; `chat/types.go:35-39`, `:50-69`; `chat/service.go:160-197`; `memory/identity.go:20-25`; `continuity.go:54-62`.
- **Caller/callee chain:** `LoadChatData` loads character files by name; memory filters by stable character ID; continuity exports/imports by character ID with no continuity ID.
- **Establishes:** Character IDs are stable in the entity and memory filter; file resolution and chat selection remain name-based; no `continuity_id` exists.
- **Does not establish:** That renaming a character breaks loading in practice, or that the same character's state leaks across chats at runtime.
- **Prior artifact:** Consistent with the freeze's stable-character requirement; adds the name/ID split and the missing continuity key.

### R1-F19 — RAG has chat scope fields but no root identity

- **Classification:** Observed / Unestablished enforcement.
- **Evidence:** `rag/types.go:3-33`, `:57-66`; `rag/document.go:14-56`; `app.go:184-187`.
- **Caller/callee chain:** `rag.NewService(StateRoot())` → `Storage.SaveDocument` → `RAG/<docID>.json|.txt`; `LoadIndex(scope, scopeID)`.
- **Establishes:** Documents can declare `ScopeChat` and `ScopeID` (chat id), but physical storage is a global `RAG` namespace keyed by document ID. `Query.ChatID` is declared but has no other Go use in this package/search.
- **Does not establish:** Whether a caller filters documents by chat scope before retrieval; storage-level chat isolation is not present.
- **Prior artifact:** The prior inventory listed RAG as a non-F3-C path; this records the field/namespace nuance.

---

## 13. Prior-evidence reconciliation and source-version applicability

### 13.1 Artifact hashes

All input hashes in section 3 match the recorded values. The two root reports are at `docs/review/`, not `docs/research/`; their content is still the named evidence. The Amendment 02 revision-9 file hash is `073d0d4b...`, matching the recorded value.

### 13.2 What the prior evidence established

The prior root inventory established:

- There is no single canonical F3-C root.
- AppData chat state, `StateRoot/Chats` memory, `cwd/Chats` candidate/continuity, and `StateRoot/Chats/<name>` service context are distinct surfaces.
- `ProjectRoot` is a cached path string, not file identity.
- `getAppDataDir` is a literal profile path, not a known-folder API.
- Candidate/continuity had no production caller found.
- Alias and packaged-root behavior were unestablished.

This report confirms those points. It adds the stronger type-level finding that **no root ID exists at all**, plus the timestamp collision, path-validation, frontend-ID, memory-ID, and global-frontend-store findings.

The prior D-ROOT report established:

- `NewCandidateStorage` and `NewCharacterContinuity` use `os.Getwd()`.
- No production binding/generated method/startup registration was found for them.
- `chat.Service` startup reads `StateRoot/Chats/memories_snapshot.json` with an empty name and can fail construction.
- `SetChatID("")` is a no-op.
- Delete covers the legacy chat tree; wipe covers the legacy chat tree and gallery; memory/candidate/continuity/native context are not covered.
- `DEF-F3C-ROOT-01`, `WIPE-02`, `IDENTITY-03`, `REACH-04`, `STARTUP-05`, and `ALIAS-06` remain open.

This report confirms those findings at the cited source locations. It does not close any of them.

### 13.3 Source-version applicability

| Artifact | Can it be treated as current behavior? | Reason |
|---|---|---|
| Freeze F3-C | Requirements only | It is a frozen contract, not an implementation observation. |
| Implementation F3-C chat consumer | No, without a source-version check | It does not bind its observations to `HEAD` or the current dirty-tree hashes. It may describe an earlier/different tree. |
| F3-C master contract | Intended design only | It specifies `(rootID, chatID)`, target sets, and receipts; it is not evidence that the working tree implements them. |
| Root inventory | Evidence inventory for its inspected tree; not current proof | It records source paths and unknowns, but no source commit/hash binding. Its central finding is consistent with current source. |
| D-ROOT closure | Evidence for its inspected tree; not current proof | Same limitation. Its verdicts are consistent with current source. |
| Gate B source/reconciliation/queue | Current-tree findings within Gate B's stated scope | Use as prior findings, not as a substitute for this trace. |
| Amendment 02 rev 9 | Planning only | Its own header disclaims implementation/runtime authority. |

**Superseded claims:** Any prior statement that treats `(rootID, chatID)` as implemented, that treats a root path as a stable root identity, that treats chat IDs as unique, that says `SetChatID("")` clears memory scope, or that says delete/wipe covers all chat-owned data is superseded by the current source evidence. This report does not assert that any named prior artifact made those claims; it records the current conflict.

### 13.4 No architecture decision

The freeze/master artifacts define an intended contract. This report compares the current source to that contract and finds non-conformance. It does not select a canonical root, a migration, a lock design, or an implementation approach. Those remain owner decisions outside this authorization.

---

## 14. Unestablished claims and investigation limitations

| Claim | Status | Minimum evidence needed | Executed? |
|---|---|---|---|
| Packaged app resolves a specific `ProjectRoot` | **Unestablished** | Launch a clean packaged build from two working directories; inspect resolved roots | No |
| AppData redirection/known-folder behavior | **Unestablished** | Redirected-profile run; compare `%APPDATA%` and actual path | No |
| Two roots are physically distinct/identical | **Unestablished** | Volume/file identity, final path, reparse-point checks | No |
| A junction/symlink/8.3/case/UNC alias reaches the same store | **Unestablished** | Windows alias matrix with two processes | No |
| `SaveChatAs` collision occurs in a real session | **Calculated, not observed** | Create a chat and immediately Save As within one second; inspect IDs/directories | No |
| Frontend save actually invents a `Date.now()` ID | **Source-reachable, not observed** | Drive the `LoadChatData` empty-ID flow and save; inspect created directory | No |
| Cross-chat memory display occurs | **Source-derived, not observed** | Create memories in two chats; inspect `howl-chat-memories` and UI | No |
| Candidate/continuity are production-reachable | **Unestablished** | Find an App binding, generated Wails method, indirect caller, or external caller | No |
| `chat.Service.SaveChatData` is production-reachable | **Unestablished** | Release symbol/endpoint trace; prior D-ROOT found test seam only | No |
| Concurrent memory bindings interleave `SetChatID` | **Inferred risk** | Concurrent two-chat test with controlled ordering | No |
| A deleted chat's memory is reactivated by ID reuse | **Source-permitted, not observed** | Delete, recreate in same second, inspect memory reads | No |
| `SaveChatData` path traversal is exploitable through the UI | **Source-reachable, not observed** | Controlled payload with a traversal ID in a disposable profile | No |

**Search-boundary caveat:** The absence of a production caller for candidates/continuity is a finding about the inspected Go/frontend/Wails surface. It is not proof that no compiled or external caller exists. The prior D-ROOT report reached the same conclusion and preserved the same caveat.

No runtime test, model load, profile read, database operation, network request, or destructive action was performed.

---

## 15. Dependencies and handoff questions for subsequent queue items

### 15.1 Open defects this report informs but does not close

- `DEF-F3C-ROOT-01` — split chat and memory ownership. **Open.**
- `DEF-F3C-WIPE-02` — incomplete content wipe. **Open.**
- `DEF-F3C-IDENTITY-03` — empty `SetChatID` cannot clear memory identity. **Open.**
- `DEF-F3C-REACH-04` — latent `cwd/Chats` candidate/continuity roots. **Open; reachability still unestablished.**
- `DEF-F3C-STARTUP-05` — `chat.Service` startup read can fail and alter wiring. **Open.**
- `DEF-F3C-ALIAS-06` — no canonical path identity or alias rejection. **Open.**

### 15.2 IDENTITY-03 dependency

IDENTITY-03 must define:

- what "clear" means for the active chat pointer, memory scope, candidate queue, continuity config, RAG chat-scoped documents, and frontend caches;
- whether `SetChatID` accepts an explicit clear operation or requires a separate reset;
- whether deleting a chat creates a tombstone/generation that prevents ID reuse from reactivating surviving data;
- how the frontend current-chat cache is invalidated when the backend chat is deleted or wiped.

This report supplies the source dependencies but does not choose those semantics.

### 15.3 Owner decision brief R1-D01

**Question:** Which identity is authoritative for chat-owned data across the legacy chat store, memory service, candidate/continuity store, RAG, and frontend cache?

**Status:** No decision is made here. The current source has no `(rootID, chatID)` implementation; the freeze/master artifacts define an intended contract. Selecting a canonical root, a migration rule, or a lock scope is an owner/architecture decision and remains unauthorized.

**Why it blocks:** Without an authoritative identity, a fix could make one store internally consistent while leaving another store addressable by a different path or timestamp. The prior inventory's OD-18 conclusion still applies: a storage-level lock cannot be scoped from the current source.

### 15.4 Runtime evidence requests (not authorized)

If separately authorized, the minimum tests are:

1. **R1-T01 Root identity:** Launch a packaged build from two different working directories; record `ProjectRoot`, `StateRoot`, AppData, and final paths.
2. **R1-T02 Alias:** Create a junction/symlink/8.3/case variant to the same directory; verify whether writes collide and whether IDs remain distinguishable.
3. **R1-T03 Chat collision:** Call `NewChat` twice within one second; record IDs and whether the first chat survives.
4. **R1-T04 Save As collision:** Create a chat and immediately Save As within one second; record whether a second directory exists.
5. **R1-T05 Path validation:** Use a disposable profile and a controlled chatID containing `..`; verify whether `Save`/`Load`/`Delete` stay under `chats`.
6. **R1-T06 Memory scope:** Exercise `AddMemory`/`GetMemories` with empty, `default`, and two chat IDs; verify no cross-chat read/write.
7. **R1-T07 Delete/wipe:** Delete and wipe with memory present; verify `StateRoot/Chats`, candidate/continuity, and RAG residues.
8. **R1-T08 Frontend fallback:** Drive `LoadChatData` with active context and no meta ID; save; verify whether a new directory is created.
9. **R1-T09 ID reuse:** Delete a chat, create a new one in the same second, and verify whether surviving memory is selected.
10. **R1-T10 Reachability:** Inspect release symbols/endpoints for candidate, continuity, and `chat.Service.SaveChatData`; confirm or exclude production callers.

None of these tests was run.

### 15.5 Handoff to other queue items

- **F3-D:** The character name/ID split and absent continuity ID are inputs; do not infer that character persistence is already conversation-scoped.
- **A14 / V24:** Not implicated by this read-only investigation; no execution authorization.
- **Audio/TTS:** Separate; no concrete chat identity dependency was found in the inspected paths.
- **R5 proposal:** Remains proposal-only; this report does not revive or alter it.

---

## 16. Confirmation of prohibited actions avoided

This investigation was read-only. The following were not performed:

- No source, schema, configuration, build, test, or existing documentation change.
- No implementation, fix, refactor, migration, or architecture selection.
- No application launch, model load, server start, benchmark, or runtime experiment.
- No user chat data, database, profile, gallery, or memory-content access.
- No import/export, recovery, delete, wipe, clear, reset, or move operation.
- No network request.
- No Git reset, stash, checkout, restore, stage, commit, or other Git-state mutation.
- No modification of any frozen or evidence artifact.
- No claim that absence of a search result proves absence beyond the stated boundary.

The only file created is this report. All prior artifacts named in section 3 retain their recorded hashes.

---

## Conclusion

### Established

1. Production chat and memory data have **no root identity**. `ChatMeta`, `ChatData`, `AppConfig`, `ChatSummary`, and `MemoryRecord` contain no `RootID`; roots are process-local path strings.
2. Chat IDs are **second-resolution timestamps** with no collision check. `NewChat` can overwrite an existing chat, and `SaveChatAs` can rewrite the source instead of copying it.
3. Chat IDs are used as **unvalidated filesystem path components** by save, load, delete, history, and active-meta paths. `SaveChatData` accepts a frontend-supplied `meta.id`.
4. The memory service has **no per-record conversation identity**. `SetChatID("")` is a no-op; the initial scope is `"default"`; empty bindings reuse the last scope.
5. Delete and backend wipe do **not** cover `StateRoot/Chats`, candidate/continuity, or RAG; `WipeAllData`'s clear call is ineffective.
6. `GetChatList`, `LoadChat`, and `GetActiveChatMeta` apply **different ID fallbacks**; metadata and directory identity can disagree.
7. `SaveChatAs` copies the chat-data `memories.json`, not the memory-service store; the two stores share a filename but not ownership.
8. `chat.Service` is **name-keyed** and starts with an empty name, reading root-level snapshot/memory paths.
9. Backend memory IDs also collide at one-second resolution.
10. The frontend memory store is a **global** `localStorage` array, and the save fallback can invent a `Date.now()` chat ID.

### Contradicted

- The contract's `(rootID, chatID)` authoritative identity is not implemented.
- Chat identity is not stable across save, load, list, Save As, and frontend fallback.
- Missing or foreign ownership does not consistently fail closed.
- Conversation-local ownership by default is contradicted by `"default"`/last-scope memory, the global frontend store, and cwd-based latent roots.
- Delete/wipe do not establish that chat-owned data is gone.

### Unknown

- Packaged `ProjectRoot`, AppData redirection, and physical path identity.
- Actual frequency and user-visible effect of timestamp collisions.
- Actual cross-chat display/leak in the running UI.
- Production reachability of candidate/continuity and the native service save path.
- Concurrency outcomes of the mutable memory scope.
- Whether a runtime `SaveChatData` path traversal is reachable through the UI.

### Findings that block reliable ownership conclusions

`R1-F01`, `R1-F03`, `R1-F04`, `R1-F05`, `R1-F06`, `R1-F09`, `R1-F10`, `R1-F11`, `R1-F13`, `R1-F15`, and `R1-F16` block a claim that production ownership is reliable. The absence of a root ID alone is sufficient to block the `(rootID, chatID)` conclusion; the collision and path-validation findings make the chat-ID half unreliable as well.

### Questions requiring a separate owner decision or runtime evidence

- **Owner decision R1-D01:** Which identity is authoritative across the split stores, and what migration/quarantine rule applies? No selection is made here.
- **Runtime evidence R1-T01–R1-T10:** The minimum tests listed in section 15.4. None was executed.
- **IDENTITY-03:** Full clear/reset semantics remain a separate work item. This report records dependencies only.

**Completion statement:** ROOT-01 has mapped the evidence problem. It has not proven the ownership contract satisfied, selected an architecture, authorized a fix, or established runtime behavior. The open defects and the missing root identity remain blocking.

---

## Report identity

- Report body SHA-256 (excluding this block): `bd2564e6f8708fc6dbafdb8b3c53bd4fe1f6bb43e31ba3ea9f54e86fafb12068`
- Report body bytes (excluding this block): `74160`
- Report body lines (excluding this block): `859`

This block is excluded from the body digest so the digest is not self-referential. The full-file SHA-256, byte count, and line count are recorded in the delivery message after the file is finalized.
