# F3-C D-ROOT CLOSURE INVESTIGATION

**Status:** **READ-ONLY SOURCE INVESTIGATION — NOT R6, NOT A DECISION RECORD, NOT FREEZE, NOT ACCEPTANCE, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Companion inventory:** `docs/review/INVENTORY-F3C-ROOT-OWNERSHIP-2026-11-11.md` (`cd722ef0628bf5d3f2517d41f6276520a55552193d05976be30f88e67f43fdb8`)
**Decision request:** `docs/review/OWNER-DECISION-REQUEST-F3C-FOUNDATIONS-2026-11-11.md`

**Purpose:** close the remaining D-ROOT evidence gaps identified after the ownership inventory. This document does not select a canonical root, authorize cleanup, authorize migration, authorize lock testing, or draft R6.

**Scope:** production reachability of candidate and continuity storage; `chat.Service` startup read effects; every `SetChatID` path and its empty-ID behavior; deletion, wipe, and recovery coverage across the discovered storage trees; residual unknowns; defect entries.

---

## 1. Verdict

| Question | Verdict |
|---|---|
| Are candidate and continuity storage production-reachable? | **No production binding, generated Wails method, indirect service caller, or startup registration was found in the inspected source.** They remain latent and must not be assumed production-dead without the residual external-caller caveat. |
| Can `chat.Service` startup reads influence production behavior? | **Yes, narrowly.** `NewService` always calls `LoadChatData` at startup against `StateRoot`. A corrupt or unreadable `StateRoot/Chats/memories_snapshot.json` can fail construction, leaving `chatService` nil and skipping RAG/content-library wiring. Startup does not write chat data. |
| Does any path clear the memory service's active chat identity? | **No.** `SetChatID("")` returns without changing state. `DeleteChat` and `WipeAllData` do not clear it; after either, the previous memory ID remains selected. |
| Does deletion/wipe cover all discovered storage trees? | **No.** `DeleteChat` covers the legacy chat tree. `WipeAllData` covers the legacy chat tree and gallery. Memory, candidate, continuity, and native-service context trees are not covered. |
| Is D-ROOT ready to close? | **No.** The evidence now supports rejecting the assumption of one existing canonical store, but the defect register and residual unknowns remain open. |

---

## 2. Reachability of candidate and continuity storage

### 2.1 Candidate queue

Source path:

- `memory.Service.ProcessMessageForEvents` builds candidates and calls `s.AddCandidates` (`internal/backend/memory/service.go:315-340`).
- `AddCandidates` constructs `NewCandidateStorage(s.chatID)` (`internal/backend/memory/candidates.go:85-95`).
- `NewCandidateStorage` resolves `os.Getwd()` and writes `Chats/<chatID>/candidates.json` (`internal/backend/memory/candidates.go:30-35`, `:49-80`).

Reachability checks:

| Route | Result | Evidence |
|---|---|---|
| Direct `App` binding | **Not found.** No `func (a *App) ProcessMessageForEvents` or candidate method exists. | Source search. |
| Generated Wails binding | **Not found.** `frontend/wailsjs/go/main/App.d.ts` and `App.js` contain memory CRUD, but no `ProcessMessageForEvents`, `AddCandidates`, or candidate method. | Generated binding source. |
| Frontend call | **Guarded and non-functional.** `frontend/dist/js/memory-candidates-menu.js:100-107` checks `window.go?.main?.App?.ProcessMessageForEvents` before calling it; without a generated binding the check is false and the local fallback runs. | Source. |
| Indirect service caller | **Not found.** `ProcessMessageForEvents` has no caller other than the frontend reference. `AddCandidates` has no caller other than `ProcessMessageForEvents`. | Source search. |
| Startup registration | **Not found.** No constructor or startup path registers candidate processing. | Source. |
| Native/test build | **No candidate-specific native seam found.** The native test seam exposes memory CRUD through the App, not this method. | Source. |

**Verdict:** **No production reachability found in the inspected source.** The path is latent, not proven dead against future or external callers.

### 2.2 Continuity storage

Source path:

- `memory.Service.ExportCharacterMemoriesForContinuity` and `ImportCharacterMemoriesForContinuity` call `NewCharacterContinuity` (`internal/backend/memory/service.go:421-448`).
- `NewCharacterContinuity` resolves `os.Getwd()` (`internal/backend/memory/continuity.go:59-62`).
- `SaveContinuityConfig` and `LoadContinuityConfig` write/read `cwd/Chats/<chatID>/continuity_config.json` (`continuity.go:154-185`).

Reachability checks:

| Route | Result | Evidence |
|---|---|---|
| Direct `App` binding | **Not found.** No App method for continuity export or import. | Source search. |
| Generated Wails binding | **Not found.** No continuity method in the generated App bindings. | Generated binding source. |
| Frontend call | **Not found.** No continuity call appears in the inspected frontend. | Source search. |
| Indirect service caller | **Not found.** The two continuity methods have no caller except tests. | Source search. |
| Startup registration | **Not found.** No startup path calls continuity. | Source. |

**Verdict:** **No production reachability found in the inspected source.** The continuity tree is latent. A missing direct binding alone is not proof of unreachability, but the absence of generated binding, indirect caller, and startup registration closes the inspected routes.

### 2.3 Why this matters despite no production caller

A latent path is still a future ownership hazard:

- it resolves a third root (`cwd/Chats`);
- it is keyed by chat ID;
- it is not included in `DeleteChat` or `WipeAllData`;
- if a future binding or internal caller enables it, the existing canonical-root decision would not cover it.

**Disposition:** record as an open reachability/ownership defect, not as “safe because no caller was found.”

---

## 3. `chat.Service` startup reads

### 3.1 Initialization chain

1. `NewApp` constructs `chat.NewService(...)` at `app.go:167-172`.
2. `chat.NewService` calls `s.LoadChatData("", "", "", "", []string{}, "")` at `internal/backend/chat/service.go:99-110`.
3. `LoadChatData` sets `projectRoot := paths.StateRoot()` at line 148.
4. With empty world/scenario IDs, canonical world and scenario return `ContextNotSet`; no content is loaded.
5. With no character names, no character files are read.
6. It reads `filepath.Join(projectRoot, "Chats", chatName, "memories_snapshot.json")` at line 200. With `chatName == ""`, this is `<StateRoot>/Chats/memories_snapshot.json`.
7. `loadMemoriesSnapshot` returns an empty slice when the file is absent, but returns an error for an unreadable file or malformed JSON (`internal/backend/chat/loader.go:41-63`).
8. User memories are loaded through `loadMemoriesFromService`; its error is logged and downgraded to an empty slice, not a construction failure (`service.go:204-218`).

### 3.2 Production effect

| Effect | Evidence | Confidence |
|---|---|---|
| Startup read is always attempted. | `NewService` calls `LoadChatData`; `NewApp` constructs the service. | **SC** |
| Startup does not write chat data. | No write call in `LoadChatData`; `SaveChatData` is separate. | **SC** |
| A corrupt/unreadable `<StateRoot>/Chats/memories_snapshot.json` can make `NewService` return an error. | `loadMemoriesSnapshot` returns the error; `LoadChatData` returns it. | **SC** |
| `NewApp` logs that error and sets `chatService = nil`; it does not fail startup. | `app.go:167-172`. | **SC** |
| With `chatService == nil`, RAG attachment and `SetContentLibraries` are skipped. | `app.go:184-190`, `:248-253`. | **SC** |
| RAG toggles are nil-safe. | `app-audio.go:570-581`. | **SC** |
| Service getters are not used by a production caller. | `syncAppStateFromChatService` has no caller; only native seam calls `LoadChatData`/`GetContextStatus`. | **SC** for no caller found; **INF** for no runtime effect. |

**Verdict:** startup reads are production-reachable and can affect **service availability and RAG/content-library wiring**. They do not write chat data and do not appear to influence production chat content, because the service's data getters have no production caller and its save path is test-seam-only.

### 3.3 What is still unknown

- Whether a packaged profile contains a corrupt `<StateRoot>/Chats/memories_snapshot.json`; not tested.
- Whether any external release build calls `syncAppStateFromChatService`; no caller found in source.
- Whether a future consumer will use the service getters. That is outside this read-only check.

---

## 4. Every `SetChatID` path and empty-ID behavior

### 4.1 Contract

`memory.Service.SetChatID` (`internal/backend/memory/service.go:33-38`):

```go
if strings.TrimSpace(chatID) == "" {
    return
}
s.chatID = chatID
```

There is no clear operation. An empty ID is ignored, not applied.

### 4.2 Production callers

| Caller | Argument | Empty-ID behavior | Resulting scope |
|---|---|---|---|
| `App.NewChat` (`app-chat.go:56`) | New chat ID | Non-empty. | Memory scope switches to new ID. |
| `App.SaveChatAs` (`app-chat.go:86`) | New chat ID | Non-empty. | Memory scope switches to copy ID. |
| `App.LoadChat` (`app-chat.go:111`) | Selected chat ID | Non-empty. | Memory scope switches to selected ID. |
| `App.GetMemories` (`app-memory.go:72`) | Caller chat ID | Empty retains prior scope. | Prior memory chat, or empty initial scope. |
| `App.GetMemoryDetail` (`app-memory.go:97`) | Caller chat ID | Empty retains prior scope. | Prior memory chat. |
| `App.EditMemory` (`app-memory.go:114`) | Caller chat ID | Empty retains prior scope. | Prior memory chat. |
| `App.AddMemory` (`app-memory.go:137`) | Caller chat ID | Empty retains prior scope. | Prior memory chat. |
| `App.PinMemory` (`app-memory.go:170`) | Caller chat ID | Empty retains prior scope. | Prior memory chat. |
| `App.DeleteMemory` (`app-memory.go:185`) | Caller chat ID | Empty retains prior scope. | Prior memory chat. |
| `App.ExportMemories` (`app-memory.go:200`) | Caller chat ID | Empty retains prior scope. | Prior memory chat. |
| `App.GetAutoCaptureConfig` (`app-memory.go:225`) | Caller chat ID | Empty retains prior scope. | Prior memory chat. |
| `App.SetAutoCaptureConfig` (`app-memory.go:247`) | Caller chat ID | Empty retains prior scope. | Prior memory chat. |
| `App.GetMemoryIdentityDiagnostics` (`app-memory.go:269`) | Caller chat ID | Empty retains prior scope. | Prior memory chat. |
| `App.WipeAllData` (`app-wipe.go:85`) | `""` | **Ignored; no clear.** | Prior memory chat remains selected. |
| `chat.Service.NativeTestSaveChatData` (`snapshot_probe_native_test_seam.go:20`) | Chat name | Non-empty test path. | Test-only, `native_test_seam` build. |

### 4.3 Non-callers that matter

- `App.DeleteChat` does not call `SetChatID` at all (`app-chat.go:117-129`).
- `App.ResetApplicationSettings` does not call `SetChatID` (`app-settings-reset.go:17-56`).
- `chat.NewService` constructs `memory.NewService(chatName)`; at startup `chatName` is empty, so the memory service begins with an empty scope (`service.go:99-110`). An empty-scope memory read/write resolves to `<StateRoot>/Chats` rather than a per-chat directory.
- `chat.Service.SaveChatData` calls `s.memoryService.SaveChatData(chatName, ...)` without first calling `SetChatID`; the native test helper calls `SetChatID(chatName)` first (`service.go:546`; `snapshot_probe_native_test_seam.go:19-23`).

### 4.4 Verdict

There is **no supported clear operation** for memory scope. Empty ID is silently ignored. Delete and wipe leave the prior chat selected. A subsequent memory read/write can recreate or continue writing under the deleted chat ID.

---

## 5. Deletion, wipe, and recovery coverage

### 5.1 Storage trees

| Tree | Location | Primary content |
|---|---|---|
| **A** | `%USERPROFILE%\AppData\Roaming\HOWLChat` normal; `<StateRoot>\persisted-chats` isolated | Chat metadata, history, context snapshots, `LastChatID` |
| **B** | `<StateRoot>\Chats\<chatID>` | Memory records, backup, auto-capture config |
| **C** | `<cwd>\Chats\<chatID>` | Candidate queue, continuity config, if reached |
| **D** | `<StateRoot>\Chats\<chatName>` | Native/service `chat.json`, `history.json`, `memories_snapshot.json`; save is test-seam-only, startup read occurs |

### 5.2 Coverage matrix

| Operation | A legacy chat | B memory | C candidate/continuity | D native context |
|---|---|---|---|---|
| Create | `PersistenceManager` creates root/chats/config. | Created on first memory write. | Created if memory methods are reached. | Created only by test-seam save; startup does not create. |
| Read | `LoadChat`, `GetChatList`, active meta. | App memory bindings. | Memory methods if reached. | Startup `LoadChatData` with empty name; native seam. |
| Update | `SaveChat`, `SaveChatAs`, `performSave`. | Memory add/edit/pin/config. | Candidate/continuity methods if reached. | `Service.SaveChatData` test-seam only. |
| Delete | `PersistenceManager.DeleteChat` removes `chats/<id>`. | **No delete path found.** | **No delete path found.** | **No delete path found.** |
| Wipe | `WipeAllData` removes `RootDir()/chats` and gallery. | **Not removed.** | **Not removed.** | **Not removed.** |
| Recovery | `initialize` creates missing root/chats; `LoadChat` best-effort; no orphan scan. | Missing file reads as empty; no reconciliation. | No reconciliation. | No reconciliation; corrupt file can fail service construction. |
| Active pointer | `LastChatID` in A. | `memoryService.chatID` in memory. | `cwd`-derived path; no pointer. | `chatName` supplied by caller; no pointer. |

### 5.3 Consequence

Deleting or wiping a chat can leave:

- B memory records for the same chat ID;
- C candidates/continuity under the same chat ID if those paths become reachable;
- D native context under the display name if a test-seam save occurred.

The active memory ID can remain selected after both delete and wipe because empty-ID clearing is ignored.

---

## 6. Defect register

Each defect is open. No disposition here authorizes a fix. Acceptance requires a recorded owner decision or formal amendment, not merely a code change.

### DEF-F3C-ROOT-01 — Split chat and memory ownership

| Field | Entry |
|---|---|
| Severity | **Critical** |
| Evidence | AppData `PersistenceManager` for chat identity/history (`persistence.go:41-75`, `:192-254`); `memory.Service` for `StateRoot/Chats` (`service.go:40-58`, `:363-389`). Same chat ID used by both after `NewChat`/`LoadChat`. |
| Impact | A chat can be deleted or wiped from A while B remains; no transaction spans the stores; one lock cannot cover both. |
| Disposition | **OPEN — blocks D-ROOT and OD-18.** |
| Owner / authority | D-ROOT decision owner; OD-18 for enforcement scope. |
| Acceptance criteria | A recorded canonical root, or explicit per-surface ownership with a shared identity and cross-root consistency rule; creation, update, delete, wipe, recovery, and lock scope all cover every retained surface; source and runtime evidence. |

### DEF-F3C-WIPE-02 — Incomplete content wipe

| Field | Entry |
|---|---|
| Severity | **Critical** |
| Evidence | `WipeAllData` removes only `chatPersistence.RootDir()/chats` and gallery (`app-wipe.go:31-110`). It does not remove `StateRoot/Chats`, `cwd/Chats`, or native context. `WipeAllData`'s `memoryService.SetChatID("")` is a no-op (`app-wipe.go:85`; `memory/service.go:33-38`). |
| Impact | “Wipe all data” can report success while memory and latent chat context survive. |
| Disposition | **OPEN — blocks wipe/retention acceptance and OD-02b.** |
| Owner / authority | D-ROOT for surface list; OD-02b for wipe/retention authority. |
| Acceptance criteria | Every F3-C surface is enumerated; each is wiped, explicitly excluded with authority, or proven absent; wipe is staged and failure-atomic across the retained set; post-wipe tests assert no B/C/D residue and no active memory ID; test covers a previously deleted/selected chat. |

### DEF-F3C-IDENTITY-03 — Empty `SetChatID` cannot clear memory identity

| Field | Entry |
|---|---|
| Severity | **High** |
| Evidence | `SetChatID` returns on empty input (`memory/service.go:33-38`); `WipeAllData` calls it with `""`; `DeleteChat` does not call it. |
| Impact | After delete/wipe, memory operations can continue against the old chat ID; an empty initial ID can resolve to the `Chats` root rather than a per-chat directory. |
| Disposition | **OPEN — blocks lifecycle correctness.** |
| Owner / authority | D-ROOT and OD-02b. |
| Acceptance criteria | Explicit clear/reset operation with defined behavior; delete and wipe invoke it; empty ID cannot address the shared `Chats` root; tests assert no read/write after clear and no orphan recreation. |

### DEF-F3C-REACH-04 — Latent `cwd/Chats` candidate and continuity roots

| Field | Entry |
|---|---|
| Severity | **Medium** |
| Evidence | `NewCandidateStorage` and `NewCharacterContinuity` use `os.Getwd()` (`candidates.go:30-35`; `continuity.go:59-62`). No direct App binding, generated Wails method, indirect caller, or startup registration was found. |
| Impact | If a future binding or service call enables them, they create a third root outside every current cleanup and lock decision. |
| Disposition | **OPEN — reachability and ownership unproven; no production writer observed.** |
| Owner / authority | D-ROOT. |
| Acceptance criteria | Recorded production reachability matrix; either proven dormant and excluded from production, or included in the canonical root/wipe/lock surface; no `cwd`-dependent root allowed in a supported path. |

### DEF-F3C-STARTUP-05 — `chat.Service` startup read can fail and alter wiring

| Field | Entry |
|---|---|
| Severity | **Medium** |
| Evidence | `NewService` always calls `LoadChatData` (`service.go:99-110`); reads `<StateRoot>/Chats/memories_snapshot.json` (`service.go:200`); corrupt/unreadable file returns an error (`loader.go:41-63`); `NewApp` sets `chatService = nil` and skips RAG/content-library wiring (`app.go:167-190`, `:248-253`). |
| Impact | A malformed root-level snapshot can silently disable service wiring while startup continues; service data getters have no production caller, but the failure is not reported as a root-integrity error. |
| Disposition | **OPEN — startup behavior and root-integrity reporting unverified.** |
| Owner / authority | D-ROOT and OD-18. |
| Acceptance criteria | Startup behavior is explicit: fail closed or report a root-integrity error; no silent nil-service degradation; corrupt/unreadable snapshot test; no startup write; defined RAG/content-library outcome. |

### DEF-F3C-ALIAS-06 — No canonical path identity or alias rejection

| Field | Entry |
|---|---|
| Severity | **High** |
| Evidence | `ProjectRoot`, `StateRoot`, and `getAppDataDir` produce path strings; no file-identity or reparse-point check exists in the inspected F3-C resolvers. |
| Impact | Junctions, symlinks, 8.3 names, case variants, mapped drives, or a changed cwd can create an unprotected alternate route or a false “different root” conclusion. |
| Disposition | **OPEN — blocks a storage-level lock claim.** |
| Owner / authority | D-ROOT and OD-18. |
| Acceptance criteria | Canonical identity defined by volume/file identity or equivalent; alias tests for junction/symlink/8.3/case/UNC/mapped path; any unprovable alias fails closed. |

---

## 7. Residual unknowns

| Unknown | Status | Why it remains |
|---|---|---|
| External/dynamic caller enables candidate or continuity | **UNK** | No source caller found; future bindings or external code cannot be excluded. |
| Packaged build contains native test seam | **UNK** | Build-tagged source exists; release binary not inspected. |
| Corrupt `<StateRoot>/Chats/memories_snapshot.json` in a real profile | **UNK** | No packaged runtime test. |
| Actual packaged `ProjectRoot` and launch cwd | **UNK** | No clean-install or portable launch was run. |
| Windows profile/AppData redirection | **UNK** | `getAppDataDir` hardcodes the profile path; no redirected-profile run. |
| Physical identity of A/B/C/D | **UNK** | No volume/file identity check. |
| `chatStorageLocation` consumed outside inspected Go | **UNK** | No backend reader found; not proof of all consumers. |
| Wipe/delete behavior under permission failures or missing roots | **UNK** | No failure-injection tests for the split roots. |
| `NewPersistenceManager` failure behavior in production | **SC** for logging/nil; **UNK** for packaged UX | App continues with nil persistence; no fail-closed test. |
| Memory operations before the first `SetChatID` | **INF** | Empty scope resolves to `Chats` root; no runtime test. |
| Concurrent A/B/D writes | **UNK** | No concurrency test; no shared lock. |

---

## 8. D-ROOT status

- **D-ROOT:** **not ready to close.** The assumption of one existing canonical store is rejected; the complete production ownership boundary is not yet established.
- **D-STORAGE:** **blocked by D-ROOT.** No lock scope can be selected while A/B/C/D remain independently addressable.
- **D-LAYOUT, D-COMPAT, D-LOCK:** undecided.
- **OD-18:** unresolved.
- **State machine:** not total.
- **Lock proof:** conditional and incomplete.
- **R6:** **not authorized; not written.**
- **F3-C phase 0:** not closed.
- **F3-C acceptance:** not granted.
- **F3-D:** blocked.
- **V20:** blocked pending F3-F evidence.
- **A14 / V24:** open / blocked.

**Next D-ROOT input:** a recorded disposition for DEF-F3C-ROOT-01, DEF-F3C-WIPE-02, and DEF-F3C-IDENTITY-03; then acceptance or rejection of the latent-path and startup-read findings. No cleanup, migration, root unification, lock test, or implementation is authorized by this document.

---

## Appendix A — Source index

| Area | Source |
|---|---|
| Candidate reachability | `internal/backend/memory/service.go:315-340`; `internal/backend/memory/candidates.go:30-35`, `:49-95` |
| Continuity reachability | `internal/backend/memory/service.go:421-448`; `internal/backend/memory/continuity.go:59-62`, `:154-185` |
| Generated App bindings | `frontend/wailsjs/go/main/App.d.ts`; `frontend/wailsjs/go/main/App.js` |
| Frontend candidate call | `frontend/dist/js/memory-candidates-menu.js:100-107` |
| `chat.Service` startup | `internal/backend/chat/service.go:99-110`, `:147-215`; `internal/backend/chat/loader.go:41-63` |
| App startup wiring | `app.go:159-190`, `:248-253`; `app-audio.go:570-581` |
| `SetChatID` | `internal/backend/memory/service.go:33-38`; `app-chat.go:56`, `:86`, `:111`; `app-memory.go:72-269`; `app-wipe.go:85`; `internal/backend/chat/snapshot_probe_native_test_seam.go:20` |
| Delete/wipe | `app-chat.go:117-129`; `internal/backend/chat/persistence.go:471-484`; `app-wipe.go:31-110` |
| Inventory baseline | `docs/review/INVENTORY-F3C-ROOT-OWNERSHIP-2026-11-11.md` |
