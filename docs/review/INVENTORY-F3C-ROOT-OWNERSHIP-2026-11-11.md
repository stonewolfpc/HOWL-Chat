# F3-C ROOT OWNERSHIP EVIDENCE INVENTORY

**Status:** **READ-ONLY SOURCE INVENTORY — NOT R6, NOT A DECISION RECORD, NOT FREEZE, NOT ACCEPTANCE, NOT IMPLEMENTATION.**
**Date:** 2026-11-11
**Decision request:** `docs/review/OWNER-DECISION-REQUEST-F3C-FOUNDATIONS-2026-11-11.md`
**Adjudication package:** `docs/review/ADJUDICATION-PACKAGE-F3C-OD18-INTEGRATION-2026-11-11.md`

**Purpose:** supply evidence for **D-ROOT** only. This document does not choose the canonical F3-C root, authorize D-STORAGE, authorize lock testing, or draft R6.

**Scope:** F3-C chat identity, history, context, memory, candidate, continuity, active-pointer, and wipe/reset paths. Model, gallery, RAG, audio, and unrelated content roots are mentioned only where they explain a resolver or a bypass.

---

## 1. Confidence labels

| Label | Meaning |
|---|---|
| **SC** | Source-confirmed by direct inspection of the cited file and line range. |
| **RC** | Runtime-confirmed by an executed application or filesystem test. **No RC rows exist in this inventory.** |
| **INF** | Inferred from source structure, but not executed or not proven to be the packaged-runtime behavior. |
| **UNK** | Not established by the inspected source or available evidence. |

A source test definition is **SC** as code, not **RC** as behavior. No application launch, packaged-build inspection, Windows alias test, or two-process test was performed.

---

## 2. Central finding

**There is no single canonical F3-C root.**

The inspected source contains at least four distinct root surfaces relevant to F3-C chat state:

1. **AppData legacy chat store** — `PersistenceManager` uses `%USERPROFILE%\AppData\Roaming\HOWLChat` in normal mode. This holds chat metadata, history, context snapshots, and a `LastChatID` config.
2. **State-root memory store** — `memory.Service` and `MemoryStorage` use `paths.StateRoot()` and write `Chats/<chatID>` beneath it.
3. **Working-directory candidate/continuity store** — `NewCandidateStorage` and `NewCharacterContinuity` use `os.Getwd()` and write `Chats/<chatID>` beneath the current working directory.
4. **State-root/native-service context store** — `chat.Service` uses `paths.StateRoot()` and can write `Chats/<chatName>`; its save path is reached only from the `native_test_seam` build in the inspected callers.

In addition, isolated mode substitutes `StatePath("persisted-chats")` for the legacy AppData store, producing a fifth surface:

5. **`<state>/persisted-chats`** — the legacy store when isolation is active.

And there is a non-F3-C content-library path:

6. **`Documents\Howl\Chats`** — the default content-library destination for the `chats` category, not the chat-state writer.

The application therefore is not one logical store with two aliases. It is a set of independently rooted stores, some sharing the same chat ID or display name.

**Consequence:** a lock on any one root protects only that root. A lock on AppData protects chat metadata/history but not memory; a lock on `StateRoot` protects memory but not history; a lock on `cwd` protects neither consistently.

---

## 3. Resolver inventory

| Resolver | Source | Resolution | Relevance | Confidence |
|---|---|---|---|---|
| `paths.ProjectRoot` | `internal/paths/paths.go:29-73` | Searches upward from the executable for `wails.json` or `go.mod`, max depth 6; if not found, searches upward from cwd; if cwd is available but has no marker, returns cwd; if cwd is unavailable, returns the executable directory; result cached process-wide. | Base for `StateRoot`; not a canonical file identity. | **SC** |
| `paths.StateRoot` | `internal/paths/isolation.go:185-191` | Isolated `<root>/data` when isolation is active; otherwise `ProjectRoot()`. | Base for memory and native chat-service paths. | **SC** |
| `paths.StatePath` | `internal/paths/isolation.go:194-196` | Joins elements onto `StateRoot()`. | Memory, native context, settings, isolation legacy root. | **SC** |
| `paths.InternalConfigDir` | `internal/paths/content.go:74-76` | `StateRoot()/config`. | Settings, logs, registries, test seam token; not chat identity. | **SC** |
| `paths.DefaultContentDir` | `internal/paths/content.go:25-64` | Known folder for category; `chats` maps to `Documents/Howl/Chats`; isolated mode maps to `<state>/content/...`. | Explicit content-library destination only; no chat-state writer found. | **SC** |
| `chat.NewPersistenceManager` | `internal/backend/chat/persistence.go:41-58` | Calls `getAppDataDir()`, initializes that root. | Production chat metadata/history store. | **SC** |
| `getAppDataDir` | `internal/backend/chat/persistence.go:62-75` | Isolation: `StatePath("persisted-chats")`. Normal: `os.UserHomeDir()` plus literal `AppData/Roaming/HOWLChat`. | Legacy AppData root; hardcoded, not `%APPDATA%` or a known-folder API. | **SC** |
| `NewPersistenceManagerWithDir` | `internal/backend/chat/persistence.go:32-38` | Accepts an arbitrary root. | Test-oriented; inspected callers are tests and benches, not production. | **SC** |
| `memory.NewService` / `NewMemoryStorage` | `internal/backend/memory/service.go:40-58` | `projectRoot: paths.StateRoot()`. | Production memory records via app bindings. | **SC** |
| `memory.NewCandidateStorage` | `internal/backend/memory/candidates.go:30-35` | `root, _ := os.Getwd()`. | Candidate queue under `cwd/Chats/<chatID>`. | **SC** |
| `memory.NewCharacterContinuity` | `internal/backend/memory/continuity.go:59-62` | `root, _ := os.Getwd()`. | Continuity config under `cwd/Chats/<chatID>`. | **SC** |
| `chat.NewService` | `internal/backend/chat/service.go:99-110` | Constructs a service and calls `LoadChatData("", "", "", "", []string{}, "")`; `LoadChatData` uses `paths.StateRoot()` at line 148. | Native/service context path; reads state-root data at construction. | **SC** for construction and root; **UNK** for packaged effect. |
| `settings.NewStore` | `internal/settings/store.go:15-18` | Takes `configDir`, stores `settings.json`. | App settings; `app.getSettingsDir()` passes `InternalConfigDir()`. | **SC** |
| `app.getSettingsDir` | `app-model.go:655-660` | `paths.InternalConfigDir()`. | `StateRoot/config`; not chat identity. | **SC** |
| `config/settings.json` `chatStorageLocation` | `config/settings.json:123`; `frontend/dist/js/settings-state.js:39` | Literal `AppData/Chats` default in config/frontend. | No Go backend reader found in inspected source. | **SC** for no reader found; **UNK** for any external consumer. |

---

## 4. Caller and entry-point inventory

### 4.1 Production AppData chat paths

| Entry point | Source | Operation | Root | Confidence |
|---|---|---|---|---|
| `NewApp` | `app.go:159-165` | `chat.NewPersistenceManager()`; on error logs and sets `chatPersistence = nil`. | AppData legacy root. | **SC** |
| `GetChatList` | `app-chat.go:28-33` | Lists chat directories via manager. | AppData. | **SC** |
| `GetActiveChatMeta` | `app-chat.go:36-41` | Reads active `LastChatID` and `meta.json`. | AppData. | **SC** |
| `NewChat` | `app-chat.go:44-61` | Creates AppData chat; sets memory service to new chat ID. | AppData for chat; `StateRoot/Chats` for later memory writes. | **SC** |
| `SaveChat` | `app-chat.go:64-69` | `performSave()`; writes AppData chat. | AppData. | **SC** |
| `SaveChatAs` | `app-chat.go:72-89` | Copies through manager; sets memory service to new ID. | AppData for chat; `StateRoot/Chats` for later memory. | **SC** |
| `LoadChat` | `app-chat.go:92-114` | Loads AppData chat; sets memory service to selected ID. | AppData for chat; `StateRoot/Chats` for memory. | **SC** |
| `DeleteChat` | `app-chat.go:117-129` | Calls `PersistenceManager.DeleteChat`; clears in-memory state only when active. | AppData chat directory only. | **SC** |
| `LoadChatData` fallback | `app-characters.go:269-329` | Reads active chat through manager. | AppData. | **SC** |
| `SaveChatData` | `app-chat.go:298-307` | Converts frontend map and saves through manager. | AppData. | **SC** |
| `performSave` | `app-chat.go:314-325` | Loads current chat and saves. | AppData. | **SC** |
| `performFinalSave` | `app-chat.go:327-350` | Saves and rewrites `history.ndjson`. | AppData. | **SC** |

### 4.2 Production memory paths

| Entry point | Source | Operation | Root | Confidence |
|---|---|---|---|---|
| Memory bindings (`GetMemories`, `AddMemory`, `EditMemory`, `PinMemory`, `DeleteMemory`, `ExportMemories`, `GetAutoCaptureConfig`, `SetAutoCaptureConfig`) | `app-memory.go:66-251`; `internal/backend/memory/service.go:63-226`, `:265-310` | Reads/writes chat-scoped memory records. | `StateRoot/Chats/<chatID>`. | **SC** |
| Memory record save | `internal/backend/memory/service.go:363-389` | `memories.json` and `memories_backup.json`. | `StateRoot/Chats/<chatID>`. | **SC** |
| Memory active-scope switch | `internal/backend/memory/service.go:33-38` | Empty chat ID is ignored. | No clearing of the previous ID. | **SC** |

### 4.3 Candidate/continuity paths

| Entry point | Source | Operation | Root | Confidence |
|---|---|---|---|---|
| `NewCandidateStorage` | `internal/backend/memory/candidates.go:30-35` | Resolves `os.Getwd()`; writes `candidates.json`. | `cwd/Chats/<chatID>`. | **SC** |
| `NewCharacterContinuity` | `internal/backend/memory/continuity.go:59-62` | Resolves `os.Getwd()`; reads/writes continuity config. | `cwd/Chats/<chatID>`. | **SC** |
| `ProcessMessageForEvents` / `AddCandidates` | `internal/backend/memory/service.go:315-340`; `candidates.go:85-95` | Candidate pipeline; no App binding found. | `cwd/Chats/<chatID>` if called. | **SC** resolver; **UNK** production reachability |
| Continuity export/import | `internal/backend/memory/service.go:421-448` | Uses continuity manager; no App binding found. | `cwd/Chats/<chatID>` if called. | **SC** resolver; **UNK** production reachability |

### 4.4 Native/service context path

| Entry point | Source | Operation | Root | Confidence |
|---|---|---|---|---|
| `chat.NewService` | `internal/backend/chat/service.go:99-110` | Constructor calls `LoadChatData` with empty arguments. | `StateRoot` via line 148. | **SC** |
| `Service.LoadChatData` | `internal/backend/chat/service.go:147-215` | Reads canonical world/scenario, character files, and `Chats/<chatName>/memories_snapshot.json`. | `StateRoot`. | **SC** |
| `Service.SaveChatData` | `internal/backend/chat/service.go:514-552` | Writes `Chats/<chatName>/chat.json`, `history.json`, memory data; optional generated-lore sync. | `StateRoot`. | **SC** |
| Native test save | `internal/backend/chat/snapshot_probe_native_test_seam.go:19-24` | Calls `Service.SaveChatData(chatName, true)`. | `StateRoot`. | **SC**, build-tagged `native_test_seam` only |
| Native snapshot | `app_native_test_seam.go:315-339` | Calls `LoadChatData`. | `StateRoot`. | **SC**, build-tagged `native_test_seam` only |
| `syncAppStateFromChatService` | `app-characters.go:343-352` | Copies service getters into App state. | `StateRoot` state, if called. | **SC** definition; **UNK** caller |

**Production classification:** `Service.SaveChatData` has no production App binding in the inspected callers; `NativeTestSaveChatData` is built only with `native_test_seam`. `NewService` still runs at production startup and calls `LoadChatData` with empty arguments, so the service is instantiated, but this inventory does not claim it writes chat state in production.

### 4.5 Wipe, reset, and content-library paths

| Entry point | Source | Operation | Root | Confidence |
|---|---|---|---|---|
| `WipeAllData` | `app-wipe.go:31-110` | Stages/removes `chatPersistence.RootDir()/chats` and the gallery; resets `LastChatID`; calls `memoryService.SetChatID("")`. | AppData (normal) or `StateRoot/persisted-chats` (isolated); **not** `StateRoot/Chats`. | **SC** |
| Empty memory clear | `internal/backend/memory/service.go:33-38` | `SetChatID("")` returns without changing state. | Prior memory chat ID remains. | **SC** |
| `DeleteChat` | `internal/backend/chat/persistence.go:471-484` | `os.RemoveAll(<root>/chats/<chatID>)`. | AppData/legacy root only. | **SC** |
| `ResetApplicationSettings` | `app-settings-reset.go:17-56` | Resets settings stores. | `StateRoot/config`; no chat deletion. | **SC** |
| `AddDefaultContentLibrary("chats")` | `app-libraries.go:49-65` | Creates/registers `Documents/Howl/Chats`. | Content library only; `ExportContentItem` refuses chat export at `app-libraries.go:222-231`. | **SC** |

---

## 5. Logical store map

| Logical data | Normal-mode location | Isolated-mode location | Writer / reader | Confidence |
|---|---|---|---|---|
| Chat identity, metadata, history, context snapshots | `%USERPROFILE%\AppData\Roaming\HOWLChat\chats\<chatID>\` | `<state>\persisted-chats\chats\<chatID>\` | `PersistenceManager` | **SC** |
| Active chat pointer | `%USERPROFILE%\AppData\Roaming\HOWLChat\config.json` | `<state>\persisted-chats\config.json` | `PersistenceManager` | **SC** |
| Memory records | `<StateRoot>\Chats\<chatID>\memories.json` | Same under isolated state root | `memory.Service` | **SC** |
| Memory backup | `<StateRoot>\Chats\<chatID>\memories_backup.json` | Same | `memory.Service` | **SC** |
| Auto-capture config | `<StateRoot>\Chats\<chatID>\auto_capture_config.json` | Same | `memory.Service` | **SC** |
| Candidate queue | `<cwd>\Chats\<chatID>\candidates.json` | `<cwd>\Chats\<chatID>\candidates.json` | `CandidateStorage` if called | **SC** resolver; **UNK** reachability |
| Continuity config | `<cwd>\Chats\<chatID>\continuity_config.json` | Same | `CharacterContinuity` if called | **SC** resolver; **UNK** reachability |
| Native/service chat context | `<StateRoot>\Chats\<chatName>\chat.json`, `history.json`, `memories_snapshot.json` | Same | `chat.Service`; save only native-test path | **SC** |
| App settings | `<StateRoot>\config\settings.json` | `<state>\config\settings.json` | `settings.Store` | **SC** |
| Content-library `chats` destination | `Documents\Howl\Chats` | `<state>\content\Chats` | Explicit registration only; no chat-state writer | **SC** |

The same chat ID can appear in the AppData store and the `StateRoot/Chats` memory store. The `chat.Service` context store is keyed by **display name**, not necessarily by chat ID, and the candidate/continuity store is keyed by chat ID under **cwd**.

---

## 6. Lifecycle trace

### 6.1 Creation

- `NewApp` constructs `PersistenceManager` before the Wails run loop. `initialize()` creates the legacy root, its `chats` directory, and loads/creates its `config.json`.
- `NewChat` creates the AppData chat directory and saves its initial JSON; then `App.SetChatID`-equivalent calls point the memory service at the new chat ID.
- No memory directory is created by `NewChat` itself. It appears when a memory write occurs.
- `NewApp` also constructs `chat.Service`, which immediately calls `LoadChatData` with empty arguments against `StateRoot`. This reads but does not establish a canonical chat root.
- Failure to initialize `PersistenceManager` is logged and leaves `chatPersistence` nil; the application continues, and the memory service remains separately rooted.

### 6.2 Migration

- `PersistenceManager.migrateLegacyMeta` migrates legacy **fields inside a loaded chat**. It does not migrate a chat from one root to another.
- `getAppDataDir` has no migration branch: normal mode always resolves the AppData path; isolated mode always resolves `StatePath("persisted-chats")`.
- `NewPersistenceManagerWithDir` can point anywhere, but inspected callers are tests/benches; no production root-selection or migration path was found.
- No code moves or reconciles `cwd/Chats`, `StateRoot/Chats`, `StateRoot/persisted-chats`, and AppData chats.
- The `chatStorageLocation` setting is not consumed by the backend in the inspected source and cannot serve as a migration switch.

### 6.3 Recovery

- `PersistenceManager.initialize()` creates the legacy root and `chats`; it does not scan for a missing or alternate root.
- `LoadChat` is best-effort for optional files; `meta.json` is required. Missing or corrupt optional files are logged and ignored.
- `GetChatList` skips directories whose `meta.json` cannot be loaded. This can hide a chat without reporting a root-level inconsistency.
- Memory loading treats a missing `memories.json` as empty. A wrong root can therefore look like an empty chat rather than a root error.
- No startup reconciliation checks whether the same chat ID exists in AppData, `StateRoot/Chats`, or `cwd/Chats`.
- No root-level journal, orphan package, tombstone, or authority marker was found for F3-C chat state.

### 6.4 Deletion and wipe

- `DeleteChat` removes only `<legacy-root>/chats/<chatID>`. It does not remove `StateRoot/Chats/<chatID>`, `cwd/Chats/<chatID>`, or any native-service snapshot.
- `WipeAllData` stages and removes only `chatPersistence.RootDir()/chats` and the gallery. In normal mode that is AppData chats; in isolation it is `StateRoot/persisted-chats/chats`.
- `WipeAllData` does not remove `StateRoot/Chats`.
- `WipeAllData` calls `a.memoryService.SetChatID("")`, but `SetChatID` returns immediately for an empty string. The previous memory chat ID remains selected.
- `ResetApplicationSettings` resets settings and image settings, not chat or memory content.
- The inspected wipe tests construct a temporary `PersistenceManager` and assert only the AppData-style `chats` root and gallery. They do not create or assert a `StateRoot/Chats` memory tree.

### 6.5 Configuration

- The active chat pointer is `PersistenceManager.config.LastChatID`, stored in the legacy root's `config.json`.
- No equivalent active-chat pointer exists in `StateRoot/Chats` or `cwd/Chats`.
- `config/settings.json` contains `chatStorageLocation: "AppData/Chats"`, but no Go reader was found. It is not proof that the backend uses that setting.
- `paths.InternalConfigDir()` is `StateRoot/config`; it is used for application settings, logs, registries, and the native test token, not for the legacy chat identity.

---

## 7. Canonical identity and alias analysis

### 7.1 What is source-confirmed

- `ProjectRoot` is a **cached path string**. It is not a volume/file identity.
- `StateRoot` is either that path or a path beneath the isolation root.
- `getAppDataDir` is a **literal profile-relative path**, not `%APPDATA%` expansion or a Windows known-folder API.
- `CandidateStorage` and `CharacterContinuity` use `os.Getwd()` directly.
- No inspected F3-C root resolver calls `filepath.EvalSymlinks`, a Windows final-path API, volume serial/file-ID lookup, or a reparse-point check.
- No inspected F3-C path writer verifies that two paths resolve to the same physical directory before writing.

### 7.2 What follows

- A junction, symlink, 8.3 name, case variant, mapped drive, UNC alias, or changed drive letter can make two path strings refer to the same physical directory without detection.
- A changed working directory can move the candidate/continuity tree even when `StateRoot` is unchanged.
- A packaged build without `wails.json` or `go.mod` falls back to the launch working directory, not necessarily the executable directory. This is source behavior; packaged behavior is **UNK**.
- `ProjectRoot` caches the first resolution. A later working-directory change does not re-resolve it.
- `getAppDataDir` ignores Windows known-folder redirection. Whether the packaged app writes where the operator expects is **UNK**.

### 7.3 Same data through different routes

Normal mode resolves AppData chats and `StateRoot/Chats` to different path strings and, absent an alias, different directories. Isolation mode resolves the legacy store to `<state>/persisted-chats` and memory to `<state>/Chats`, also different strings.

**No proof exists that any two routes currently point at the same physical data, and no proof exists that they cannot.** The source does not perform the identity check needed to answer the question. This is **UNK**, not “no.”

### 7.4 Inference register

| Inference | Confidence | What would confirm it |
|---|---|---|
| A packaged build with neither `wails.json` nor `go.mod` will resolve `ProjectRoot` from its launch cwd. | **INF** from `paths.go` fallback; packaged build not run. | Launch a clean packaged build from two different working directories and inspect the resulting roots. |
| Normal-mode AppData chats and `StateRoot/Chats` are physically distinct directories. | **INF** for ordinary layouts; no file-identity check exists. | Compare volume/file identity and final paths at runtime. |
| `chatStorageLocation` is presentation metadata with no backend effect. | **INF** from “no Go reader found”; not proof that no other consumer exists. | Inspect release symbols and all settings consumers. |
| Candidate/continuity are not production-reachable chat-state writers. | **UNK/INF**; resolvers exist, but no App binding was found. | Confirm a production caller or declare the paths dormant. |
| The same chat’s AppData identity and StateRoot memory are independently writable. | **INF** from separate service/manager roots; not executed. | Two-root write test after D-ROOT selection. |

---

## 8. Answers to the D-ROOT questions

| Question | Answer | Basis |
|---|---|---|
| Which root owns F3-C chat state? | **Undecided. There is no single owner root.** AppData owns production metadata/history; `StateRoot/Chats` owns production memory; `cwd/Chats` can own candidates/continuity; the native service surface owns `StateRoot/Chats/<chatName>` context. | **SC** for the split; decision is **UNK**. |
| Does the legacy AppData path remain active for F3-C reads or writes? | **Yes.** It is the production `PersistenceManager` root for chat CRUD, active pointer, and wipe. | **SC** |
| Can both paths refer to the same data through different resolution routes? | **Not established, and not excluded.** No canonical file-identity check exists; aliases and fallback paths are unchecked. | **UNK** |
| Does every relevant caller use the same canonical root identity? | **No.** AppData, `StateRoot/Chats`, `cwd/Chats`, `StateRoot/persisted-chats`, and a content-library path are distinct surfaces. | **SC** |

---

## 9. Ownership conclusion for OD-18

A storage-level lock cannot be scoped from the current source because the source does not define one F3-C root.

- A lock on the AppData legacy root protects chat identity/history and the active pointer, but not memory.
- A lock on `StateRoot/Chats` protects production memory, but not chat identity/history.
- A lock on `cwd/Chats` would not consistently protect candidates/continuity because `cwd` can change and may differ from `StateRoot`.
- A lock on `StateRoot/persisted-chats` covers only the isolated legacy store.
- A lock on `Documents\Howl\Chats` would protect a content-library destination that is not the chat-state store.

**This inventory does not decide which root should win.** It establishes that any answer other than a single canonical root plus an explicit migration/quarantine rule would leave an unprotected F3-C surface.

---

## 10. Unknowns and evidence required

| Unknown | Why it matters | Required evidence |
|---|---|---|
| Packaged application's actual `ProjectRoot` | `StateRoot` depends on executable/cwd markers and fallback behavior. | Packaged-build launch evidence from a clean install and a portable run. |
| Windows profile/AppData redirection | `getAppDataDir` hardcodes the profile path. | Known-folder comparison and redirected-profile run. |
| Physical identity of each root | Path strings do not prove distinct or identical storage. | Volume/file identity and reparse-point handling. |
| Alias behavior | Junction/symlink/8.3/UNC can bypass a root lock. | Windows alias matrix with two-process checks. |
| Full caller set | Wails bindings and frontend paths may add writers. | Complete source/API inventory and runtime call trace. |
| Candidate/continuity production reachability | The resolvers exist, but no App binding was found. | Confirm callers or declare the paths dormant/non-production. |
| Native service save reachability | `SaveChatData` is build-tagged; only native-test callers were found. | Release-build symbol and endpoint inspection. |
| Cross-root recovery and migration | No reconciliation exists in source. | Owner decision and source-backed migration design; not implementation here. |
| Wipe completeness | Current wipe leaves memory and candidate/continuity roots. | Root inventory accepted before wipe semantics are finalized. |
| Windows lock behavior | No two-process, stale-handle, or deletion test was run. | D-STORAGE authorization and T1–T9 evidence. |

### Existing test artifacts

- `internal/paths/isolation_test.go` asserts `StateRoot`, `InternalConfigDir`, `StatePath`, and isolation containment.
- `app-wipe_test.go` constructs a temporary `PersistenceManager` and tests only the legacy `chats` root and gallery; it does not create a `StateRoot/Chats` memory tree.
- `internal/backend/chat/persistence_identity_test.go` exercises `NewPersistenceManagerWithDir`, not production root selection.

These are source-test definitions. **None is runtime-confirmed for this inventory, and none closes D-ROOT.**

---

## 11. Status

- **D-ROOT:** **undecided; evidence supplied, not decided.**
- **D-STORAGE:** **undecided; cannot be scoped until D-ROOT is recorded.**
- **D-LAYOUT:** undecided.
- **D-COMPAT:** undecided.
- **D-LOCK:** undecided.
- **OD-18:** **unresolved; no single protected root exists.**
- **State machine:** still not total.
- **Lock proof:** still conditional and incomplete.
- **R6:** **not authorized; not written.**
- **F3-C phase 0:** not closed.
- **F3-C acceptance:** not granted.
- **F3-D:** blocked.
- **V20:** blocked pending F3-F evidence.
- **A14 / V24:** open / blocked.

No source, tests, builds, migrations, model loads, destructive tests, application data, authority records, freeze, acceptance, or implementation action was performed.

---

## Appendix A — Source index

| Area | Source |
|---|---|
| Project and state root resolution | `internal/paths/paths.go:29-73`; `internal/paths/isolation.go:185-196` |
| Settings/config root | `internal/paths/content.go:74-76`; `internal/settings/store.go:15-18`; `app-model.go:655-660` |
| Legacy chat root and persistence | `internal/backend/chat/persistence.go:32-75`, `:78-136`, `:147-254`, `:256-435`, `:471-484` |
| App chat bindings and saves | `app-chat.go:28-129`, `:207-350` |
| App chat-state fallback | `app-characters.go:269-352` |
| Memory service and storage | `internal/backend/memory/service.go:33-58`, `:63-226`, `:265-310`, `:363-389`, `:421-448` |
| Candidate/continuity cwd roots | `internal/backend/memory/candidates.go:30-35`; `internal/backend/memory/continuity.go:59-62` |
| Native service context | `internal/backend/chat/service.go:99-110`, `:147-215`, `:514-552` |
| Native test seam build gate | `app_native_test_seam.go:1-20`; `app_native_test_seam_off.go:1-10`; `internal/backend/chat/snapshot_probe_native_test_seam.go:1-24` |
| Wipe and reset | `app-wipe.go:31-110`; `app-settings-reset.go:17-56`; `app-wipe_test.go:19-170` |
| Content-library path | `app-libraries.go:49-65`, `:222-231`; `internal/paths/content.go:25-64` |
| Settings metadata | `config/settings.json:123`; `frontend/dist/js/settings-state.js:39` |
