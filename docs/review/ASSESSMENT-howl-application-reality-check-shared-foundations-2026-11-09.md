# HOWL Chat — Application-Wide Reality Check and Shared-Foundation Assessment

**Date:** 2026-11-09  
**Scope:** Whole application, read-only source and documentation inspection  
**Purpose:** Establish the engineering baseline before subsystem completion  
**TTS governance:** P17 remains UNAPPROVED; contract rev3 remains frozen; OD-9 remains unmet; C1 remains UNAUTHORIZED. This report does not change those decisions.

---

## 1. Executive summary

Howl Chat is a substantial Wails desktop application, not a prototype shell. It contains real implementations for chat, character/world/scenario editing, lorebooks and memories, image generation, ASR/TTS, model discovery and loading, GPU monitoring, settings, persistence, and error logging. The application also contains several parallel implementations of the same responsibility, and the boundaries between them are not consistently owned.

The audit identified **15 major subsystems**:

| Primary status | Count |
|---|---:|
| WORKING / VERIFIED | 0 |
| IMPLEMENTED / UNVERIFIED | 3 |
| INCOMPLETE | 8 |
| STUBBED | 0 |
| BROKEN / DEFECT IDENTIFIED | 3 |
| UNKNOWN | 1 |
| **Total** | **15** |

No subsystem can be called WORKING / VERIFIED from this audit: tests and builds were not run, the application was not launched, and no runtime experiment was performed. “IMPLEMENTED / UNVERIFIED” therefore means real source-level functionality exists without current end-to-end proof.

The most consequential findings are shared-foundation defects, not isolated feature gaps:

1. **The registered shutdown path is missing.** `App.OnShutdown` exists at `app.go:283`, but `main.go` registers only `OnStartup` (`main.go:22-29`). Normal window shutdown therefore does not run the final chat save, model-server shutdown, image-engine shutdown, audio-manager drain, or `proc.CloseJob`.
2. **Persistence has no single root or contract.** Chats live under `%APPDATA%/HOWLChat` (`internal/backend/chat/persistence.go:41-56`), while characters, worlds, scenarios, and images are resolved through the working directory or executable directory (`app-characters.go:50-90`, `app-image-generate.go:165-207`). The full-wipe path removes `<exe>/chats` (`app-wipe.go:20-28`), which is not where the chat persistence manager stores chats.
3. **Chat has two architectures.** The live chat path is `App.SendMessageStreamWithImage` → `buildPromptMessages()` → `HTTPClient` (`app-prompt.go:25-147`). `internal/backend/chat.Service` has its own prompt builder, model manager, lore routing, and RAG integration (`internal/backend/chat/service.go:35-180`, `:648-745`), but is not used for generation. The RAG toggle is connected to that unused service (`app-audio.go:565-576`).
4. **Generation cancellation is a single shared slot.** `App.streamCancel` is one field (`app.go:48`), overwritten by every stream (`app-prompt.go:61`, `:205`) and cleared by whichever stream finishes (`app-prompt.go:115`, `:245-255`). Concurrent character replies cannot be independently cancelled or safely identified.
5. **The resource manager is not the authority for loaded models.** Registration happens before readiness, failed loads and model switches do not reconcile entries, image auto-load is not registered, and TTS has no admission kind (`app-model.go:504-531`, `app-image-model.go:143-150`, `app-vram-arbitration.go:47-120`).

The TTS findings are **not wholly unique**. The protocol and TTS-helper-specific lifecycle problems are subsystem-specific, but the underlying ownership failures — stale generation ownership, incomplete shutdown ordering, and state that is projected rather than owned — recur in chat streaming and across model/image/audio lifecycles.

Recommended order: **lifecycle and process ownership → persistence and path contract → generation/request ownership → resource and model-state authority → settings, error, and UI consistency → subsystem completion in dependency order.** TTS implementation remains outside this order until its governance changes.

---

## 2. Scope and evidence limitations

### 2.1 What was inspected

- Go entry point, application state, startup, shutdown, and Wails bindings.
- Chat persistence, live streaming, prompt assembly, character/world/scenario context, lore and memory services.
- Image engine, job manager, gallery, model loading, and image settings.
- Audio manager, ASR/TTS pipeline boundaries, and audio settings bindings.
- Model manager, llama HTTP client, model registry, process helpers, hardware/resource management.
- Settings stores, paths, error logger/reporter, wipe behavior.
- Frontend pages, page scripts, navigation, theme handling, and documented UI defects.
- Existing assessments, `ISSUES.md`, `TECH_DEBT.md`, `STAGE_STATUS.md`, and TTS decision records as documentary evidence.

### 2.2 What was not done

- No tests, builds, benchmarks, or runtime experiments.
- No application launch, process start/stop, model load, image generation, audio playback, or UI interaction.
- No source, test, configuration, contract, proposal, assessment, review, or decision-register modification.
- No dependency installation or repository cleanup.

### 2.3 Evidence classes used

| Class | Meaning in this report |
|---|---|
| Source fact | Directly visible in current source. |
| Documented test evidence | A saved audit, issue record, or stage record reports an observed result. Not independently reproduced here. |
| Inference | A consequence reasoned from source structure, not observed at runtime. |
| Unknown | Evidence is insufficient. |

A potential risk is not presented as a confirmed defect. A source-level implementation is not presented as working behavior.

---

## 3. Application subsystem inventory

### 3.1 Application shape

- Wails v2 desktop application with an embedded `frontend/dist`.
- Go backend bound as one `App` object (`main.go:29-31`).
- Nine HTML pages: chat (`index.html`), character, world, scenario, lorebooks, image generator, audio settings, story tools, and the application index.
- Persistent data is split across `%APPDATA%/HOWLChat`, the project root, the executable directory, `localStorage`, and IndexedDB.
- Long-running helpers: `llama-server.exe` on fixed port 8080, `sd-server.exe` on default port 8081, and optional ASR/TTS helper processes.

### 3.2 Subsystem table

| ID | Subsystem | Primary responsibility | Main source locations | Principal UI entry points | Primary status |
|---|---|---|---|---|---|
| S1 | Application lifecycle and process supervision | Startup, shutdown, child ownership, orphan handling | `main.go`, `app.go:283-315`, `internal/proc/*`, `app-llama.go` | Window lifecycle | BROKEN / DEFECT IDENTIFIED |
| S2 | Chat, prompt assembly, and streaming | Conversation state, prompt construction, token streaming, abort | `app-chat.go`, `app-prompt.go`, `frontend/dist/js/chat-*.js` | `index.html` | INCOMPLETE |
| S3 | Character management | Character CRUD, sheets, chat membership, character images | `app-characters.go`, `internal/backend/chat/character*.go`, `character*.js` | `character.html` | IMPLEMENTED / UNVERIFIED |
| S4 | World and scenario context | World/scenario CRUD and chat context | `app-characters.go`, `internal/backend/chat/world*.go`, `scenario*.go`, `world*.js` | `world.html`, `scenario.html` | IMPLEMENTED / UNVERIFIED |
| S5 | Lorebooks, memory, and RAG | Lore CRUD, trigger resolution, memory snapshots, document retrieval | `internal/backend/lorebook/*`, `internal/backend/memory/*`, `internal/backend/rag/*`, `app-rag-upload.go` | `lorebooks.html`, chat sidebar | INCOMPLETE |
| S6 | Image generation and gallery | sd-server lifecycle, jobs, txt2img/img2img/inpaint/video, gallery | `internal/image/*`, `app-image*.go`, `image-generator*.js` | `image-generator.html` | INCOMPLETE |
| S7 | Audio, ASR, and TTS | Audio generations, leases, transcription, synthesis | `internal/audio/*`, `app-audio.go`, `audio-settings.js` | `audio-settings.html`, chat attachments | BROKEN / DEFECT IDENTIFIED |
| S8 | Model discovery, registry, and loading | Model scan, family/profile admission, llama-server spawn, unload | `internal/backend/gguf/*`, `internal/backend/loader/*`, `internal/image/registry/*`, `app-model.go` | Chat/image model menus | IMPLEMENTED / UNVERIFIED |
| S9 | Hardware monitoring and resource arbitration | GPU/RAM sampling, model estimates, fit checks, eviction advice | `internal/resource/*`, `app-resource.go`, `app-vram-arbitration.go` | Model/image dialogs | INCOMPLETE |
| S10 | Settings and configuration | Unified settings, sampler/model/image/audio sections | `internal/settings/*`, `internal/image/settings/*`, `app-settings*.go` | Settings panels across pages | INCOMPLETE |
| S11 | Persistence, storage, and import/export | Chat/entity files, settings files, gallery, backup/wipe | `internal/backend/chat/persistence.go`, `app-wipe.go`, `app-characters.go`, frontend storage modules | Chat list, save/export, wipe | BROKEN / DEFECT IDENTIFIED |
| S12 | UI navigation and shared UI | Page navigation, menus, themes, dialogs, notifications | `frontend/dist/*.html`, `frontend/dist/js/menu*.js`, `navigation.js`, `theme.js` | All pages | INCOMPLETE |
| S13 | Backend communication, events, and cancellation | Wails bindings, event emission, streaming, cancellation | `app-*.go`, `frontend/wailsjs/*`, `frontend/dist/js/chat-stream.js` | All pages | INCOMPLETE |
| S14 | Error handling and recovery | Structured logging, user-facing errors, retry and recovery | `internal/errors/*`, `app-errors.go`, error paths in each subsystem | Error notifications/logs | INCOMPLETE |
| S15 | Story tools and session utilities | Auxiliary writing/session actions | `story-tools.html`, `chat-actions.js` | `story-tools.html` | UNKNOWN |

### 3.3 Additional observations

- **Model registry** is part of S8, not a separate user-facing subsystem. It records scan directories and model kinds and is consumed by both chat and image menus.
- **Gallery** is part of S6. It stores images and metadata, but image persistence is not integrated with the chat/entity persistence contract.
- **Microphone capture** is part of S7. The audio binding reports microphone availability separately from ASR readiness (`app-audio.go:263-337`).
- **Character identity persistence** and **T2.9.5B multimodal input** are cross-cutting concerns inside S2/S3, not separate subsystems.
- No package-level import cycle was identified. Runtime coupling is bidirectional: `App` owns services and services call back into application state indirectly through the Wails layer.

---

## 4. Implementation-status matrix

Status definitions are used exactly as requested. The matrix records a **primary** status; a subsystem can contain both working and broken paths.

| ID | Subsystem | Primary status | Evidence and verification scope |
|---|---|---|---|
| S1 | Lifecycle and process supervision | **BROKEN / DEFECT IDENTIFIED** | `main.go` registers only `OnStartup`; `App.OnShutdown` (`app.go:283`) is never wired. `proc.ReapOrphans` kills by executable path, not ownership (`ISSUES.md` T0.5). Crash-time job close exists, but normal shutdown cleanup does not. |
| S2 | Chat and streaming | **INCOMPLETE** | Live generation exists (`app-prompt.go:25-147`), but documented T0.1/T0.2 remain open; vision is current-turn only; cancellation is a single shared slot; history and prompt state are duplicated between `App.chatMessages`, `chat.Service`, and the frontend DOM. |
| S3 | Character management | **IMPLEMENTED / UNVERIFIED** | CRUD and context parsing exist (`app-characters.go:41-153`); path resolution is split between cwd and executable directory; no end-to-end CRUD/load verification was performed. |
| S4 | World and scenario context | **IMPLEMENTED / UNVERIFIED** | Add/load/save paths and prompt injection exist (`app-characters.go:209-228`); persistence follows the same split path rules as characters; no runtime verification. |
| S5 | Lorebooks, memory, and RAG | **INCOMPLETE** | Lore resolution and memory services exist, but the live prompt path does not call the RAG-enabled `chat.Service`; frontend lorebook state is also read from `localStorage` (`chat-stream.js:31-42`). RAG cannot be shown to affect replies. |
| S6 | Image generation and gallery | **INCOMPLETE** | Real engine/job/gallery code exists, but `ISSUES.md` records T0.3 (image model load/path failure) and T0.4 (preview broken) as open release blockers. No runtime proof they are closed. |
| S7 | Audio, ASR, and TTS | **BROKEN / DEFECT IDENTIFIED** | The manager/lease architecture is substantial, but the TTS lifecycle findings and P17/OD-9/C1 governance remain unresolved. The TTS helper path is also not governed by a currently approved implementation. |
| S8 | Model discovery and loading | **IMPLEMENTED / UNVERIFIED** | Real GGUF admission, projector pairing, spawn validation, and monitor code exist (`app-model.go:21-180`, `:250-370`). A second `internal/backend/model.Manager` is not the live owner. No load/unload verification was run. |
| S9 | Hardware and resource management | **INCOMPLETE** | `resource.Manager` tracks LLM/image/TTS types, but arbitration handles only `llm` and `image`; image auto-load does not register; failed/switch loads do not reconcile; admission is advisory. |
| S10 | Settings | **INCOMPLETE** | Unified and image-specific stores exist, with mixed atomicity. Not every persisted setting is consumed; audio settings are restored without activation by design. No restart round-trip was verified. |
| S11 | Persistence and storage | **BROKEN / DEFECT IDENTIFIED** | Chat root is `%APPDATA%/HOWLChat`; entities use cwd/exe; wipe removes `<exe>/chats` and does not remove the actual chat root or all settings. Save paths have mixed atomicity. |
| S12 | UI navigation and shared UI | **INCOMPLETE** | Nine pages and shared scripts exist, but documented T2.1 theme defects and T2.2 ghost-ID inventory remain open (18 never-created IDs over 30 reference sites). No page was rendered. |
| S13 | Backend communication and cancellation | **INCOMPLETE** | Wails bindings and event emission work by structure, but `streamCancel` is global, stream events are not request-scoped, and UI/backend state is duplicated. |
| S14 | Error handling and recovery | **INCOMPLETE** | A structured logger and optional reporter exist, but many paths swallow errors or only `fmt.Printf`; user-visible recovery is inconsistent across subsystems. |
| S15 | Story tools and session utilities | **UNKNOWN** | The page and scripts exist, but intended scope and backend completion were not established in this audit. |

### 4.1 Status totals

| Status | Count | Subsystems |
|---|---:|---|
| WORKING / VERIFIED | 0 | — |
| IMPLEMENTED / UNVERIFIED | 3 | S3, S4, S8 |
| INCOMPLETE | 8 | S2, S5, S6, S9, S10, S12, S13, S14 |
| STUBBED | 0 | — |
| BROKEN / DEFECT IDENTIFIED | 3 | S1, S7, S11 |
| UNKNOWN | 1 | S15 |
| **Total** | **15** | — |

S6 is classified INCOMPLETE because substantial image functionality exists, while its documented T0.3/T0.4 blockers remain confirmed risks. S1, S7, and S11 are BROKEN because their defining shared behavior is defective: normal shutdown does not run cleanup, TTS remains governed by unresolved lifecycle findings, and persistence cannot reliably locate or wipe its own data.

---

## 5. Shared-foundation architecture map

### 5.1 A. Application and process lifecycle

**Current structure**

- `main.go` creates `NewApp()`, embeds assets, binds `app`, and registers `OnStartup` only (`main.go:14-35`).
- `NewApp` starts `proc.ReapOrphans(paths.ProjectRoot())` in a goroutine (`app.go:119-121`).
- `App.startup` sets the Wails context, initializes error handling, and starts background model scans (`app.go:264-280`).
- `App.OnShutdown` is a complete-looking cleanup routine (`app.go:283-315`), but it is never registered.

**Ownership models**

| Child type | Start path | Ownership mechanism | Shutdown path |
|---|---|---|---|
| `llama-server.exe` | `spawnLlamaServerWithModel` in `app-llama.go` | `proc.Start` (suspended start + job assignment) | `OnShutdown` kill + `proc.CloseJob` (not wired) |
| `sd-server.exe` | `exec.Command` + `cmd.Start()`, then `proc.Attach(cmd)` (`internal/image/engine/sdcpp/engine.go:178-191`) | Job attachment after start, single `Wait` reaper | `Engine.Shutdown` (`engine.go:413-427`); `OnShutdown` calls it (not wired) |
| TTS helper | `proc.Start` in the supervisor design (`internal/audio/tts/helper/supervisor.go`) | Job ownership in the approved design, not the current C1 implementation | `audio.Manager.Close` (not wired) |
| Video | `exec.Command` in `internal/image/engine/sdcpp/video.go` | Process tree kill in the video path | Engine shutdown/cancel |

**Finding:** process ownership is not one foundation. It is three implementations with different start semantics. The Windows job object is the strongest mechanism, but `sd-server` is attached after creation and `ReapOrphans` still uses path matching rather than recorded ownership.

### 5.2 B. Concurrency and thread safety

**Real synchronization that exists**

- `resource.Manager` guards hardware and model maps.
- `image/job.Manager` owns job state.
- `sdcpp.Engine` has `procMu` and `loadMu`, with a documented single `Wait()` owner (`supervisor.go:28-122`).
- `audio.Manager` owns desired/applied generations, leases, build accounting, and bounded drain (`internal/audio/manager.go:40-220`).
- `chat.Service` uses `sync.RWMutex` around its context and lore state (`internal/backend/chat/service.go:35-50`).

**Ownership gaps supported by source**

- `App.streamCancel` is a single mutable field with no mutex (`app.go:48`). Every generic or character stream overwrites it (`app-prompt.go:61`, `:205`), and completion clears it without checking whether the clearing stream still owns it (`app-prompt.go:115`, `:245-255`).
- `App.chatMessages`, `a.characters`, `a.world`, `a.scenario`, and `a.lorebooks` are written from multiple bindings without an application-level lock. `SendMessage` and `SendMessageStreamWithImage` can run concurrently.
- `App.performSave` loads the active chat, mutates it from app state, and rewrites multiple files. Autosave runs in a goroutine (`app-prompt.go:141-147`); a final save can overlap it.
- `chat.Service.SendMessage*` holds `s.mu` for the entire model generation (`service.go:454-520`, `:535-620`). If that service were used concurrently, cancellation and reads would block behind generation. It is currently not the live path.

These are **ownership findings**, not proven data races. No race detector or runtime scheduling test was run.

### 5.3 C. Hardware and resource management

`resource.Manager` is a singleton that stores hardware snapshots and loaded model records (`internal/resource/manager.go`). `resource.Allocate` computes eviction advice (`internal/resource/allocator.go`). `App.CheckModelFit` and `App.FreeVRAMFor` implement chat/image arbitration only (`app-vram-arbitration.go:47-120`).

The foundation is not authoritative:

- Chat registers a model immediately after spawn, before the asynchronous health monitor establishes readiness (`app-model.go:504-514`). A later load failure does not unregister it.
- `UnloadModel` unregisters all LLM entries, not the specific model (`app-model.go:530-532`).
- Switching chat models kills the old process but does not reconcile old entries before registering a new name.
- Image auto-load through `StartImageJob` does not register at all; only `SetActiveImageModel` registers, with fixed `4000 MB` GPU / `2000 MB` RAM (`app-image-model.go:143-150`).
- `CheckModelFit` has no TTS kind. TTS resources are tracked as a type but do not participate in admission.
- Frontend preflight is advisory and can proceed when `can_fit` is false with no blocking entry (`INVESTIGATION-system-resource-management-architecture-2026-11-06.md`, findings 6-7).

**Consequence:** the manager is useful telemetry, but it cannot be used as a correctness boundary for “may I load this now?” or “what is actually loaded?”

### 5.4 D. Backend/UI communication

- All Go methods on `App` are bound at once (`main.go:29-31`).
- Backend-to-UI uses named Wails events through `emitEvent` (`app.go:314-336`).
- `emitEvent` guards a nil Wails context and drops events when no dispatcher exists. This prevents `runtime.EventsEmit` from terminating the process, but it also means events can silently disappear during teardown.
- Streaming uses `chat:start`, `chat:chunk`, `chat:reasoning`, `chat:complete`, `chat:error`, and `chat:aborted`; character streams add an ID suffix (`app-prompt.go:70-135`, `:210-250`).
- There is no request or generation ID in the event payload. The UI cannot distinguish a stale event from the current turn when two streams overlap.
- The frontend keeps its own chat state and DOM, and sends state back through `SaveChatData` (`app-chat.go:291-306`). The backend also builds state from `App` fields. This is a two-way projection, not a single source of truth.

### 5.5 E. Persistence and data integrity

**Storage roots**

| Data | Root | Writer | Atomicity |
|---|---|---|---|
| Chat metadata/history/world/scenario/characters/lore/memories | `%APPDATA%/HOWLChat/chats/<id>/` | `chat.PersistenceManager` | JSON files via temp+rename; `config.json` direct write; NDJSON append |
| Characters | `<cwd>/Characters` or `<exe>/Characters` | `app-characters.go` | direct `os.WriteFile` |
| Worlds | `<cwd>/Worlds` or `<exe>/Worlds` | `app-characters.go` | direct writes in surrounding code |
| Scenarios | `<cwd>/Scenarios` or `<exe>/Scenarios` | `app-characters.go` | direct writes |
| Unified settings | settings directory | `settings.Store.Save` | direct `os.WriteFile`; non-atomic by design/comment |
| Image settings | settings directory | `internal/image/settings.Store` | temp+rename |
| Image registry | settings directory | `internal/image/registry` | temp+rename |
| Gallery | image output directory | `gallery.Service` | direct writes of image and metadata |
| Lorebook/memory UI state | `localStorage` | frontend | browser-managed |

**Confirmed persistence defects**

1. **Wipe targets the wrong chat directory.** `WipeAllData` removes `filepath.Join(exeDir, "chats")` (`app-wipe.go:20-28`), while `NewPersistenceManager` uses `home/AppData/Roaming/HOWLChat/chats` (`persistence.go:41-56`). The wipe reports completion while leaving chats and `config.json` behind.
2. **Wipe does not clear all settings.** It resets in-memory fields and entity folders, but does not remove the unified settings file, image settings, model registry, or audio state.
3. **Entity paths have three resolvers.** `CreateCharacter` checks cwd then exe (`app-characters.go:50-60`); `getCharactersDir` checks cwd then exe (`:80-87`); `chat.Service.LoadChatData` uses `paths.ProjectRoot()` (`service.go:88-125`). These can disagree in development, installed, and portable layouts.
4. **`config.json` is not written atomically.** `saveConfig` calls `os.WriteFile` directly (`persistence.go:107-117`). A crash during that write can leave an unreadable active-chat pointer. `SaveChat` also writes each JSON file separately, so a multi-file save is not transactional.
5. **History append is not locked or fsynced.** `AppendHistoryEntry` opens `history.ndjson` with `O_APPEND` and writes one line (`persistence.go:678-710`). Concurrent autosaves/final saves can interleave, and a power loss can lose the tail. The final-save rewrite mitigates duplicates, but only when `OnShutdown` or `ExitApp` runs.
6. **Frontend and backend both own lorebooks.** `chat-stream.js` reads `howl-chat-lorebooks` from `localStorage` and sends it to `SetLorebooks` (`chat-stream.js:31-42`), while chat persistence also stores `lorebooks.json`. There is no single authoritative copy.

### 5.6 F. Error handling and recovery

- `internal/errors` provides a structured JSONL logger with rotation (`logger.go:1-190`) and an optional privacy-safe reporter (`reporter.go:1-130`).
- `App.initErrorSystem` and `bindErrorEmitter` connect logging to the UI (`app-errors.go`).
- Many paths still use `fmt.Printf` and return raw errors; some discard errors entirely (`_ = a.performSave()`, `_ = a.chatPersistence.SaveChat(...)`, `_ = pm.loadJSON(...)`).
- `LoadChat` is deliberately best-effort: missing/corrupt optional files are logged and replaced with empty values (`persistence.go:240-330`). That protects startup but can silently drop world/scenario/lore data.
- Model loading has meaningful failure records and user messages (`app-model.go:50-180`), but image and audio failures do not share one error contract.
- Recovery after a failed model load is inconsistent: chat may keep a stale registration, image may leave `currentModel` cleared but a process dead, and audio relies on manager generations.

### 5.7 G. Shared UI architecture

- Navigation is page-based (`window.location.href`) through `navigation.js` and menu scripts. There is no router or shared application shell.
- Shared scripts include `theme.js`, `menu.js`, `navigation.js`, `delete-confirm.js`, `delete-service.js`, `filter-engine.js`, `storage-*.js`, and tutorial data.
- `theme.js` cycles a decorative accent palette and separately handles light/dark mode (`theme.js:1-160`). T2.1 documents unresolved light-theme overrides across pages.
- T2.2 documents 18 never-created element IDs over 30 reference sites, plus duplicate/orphaned markup. This is a functional-consistency issue, not cosmetic.
- Dialogs, notifications, loading states, and error presentation are page-local rather than shared components.
- Several pages load different storage stacks (`storage-manager`, `storage-indexeddb`, `storage-scenarios`, `storage-chat`), which reinforces the split-brain persistence model.

---

## 6. Cross-subsystem dependency map

### 6.1 Direct dependencies

| Consumer | Depends on | Evidence | Nature |
|---|---|---|---|
| S2 Chat | S8 model loading + llama HTTP | `app-prompt.go:25-147`; `app-chat.go:165-188` | Hard runtime dependency |
| S2 Chat | S3 characters, S4 world/scenario, S5 lore/memory | `app-prompt.go:648-745` | Prompt-content dependency |
| S2 Chat | S11 chat persistence | `app-chat.go:291-350`; `persistence.go` | Durable-state dependency |
| S2 Chat | S13 events/cancellation | `app-prompt.go:70-135`; `app.go:314-336` | UI liveness dependency |
| S5 Lore/RAG | S2 prompt path | `chat.Service.SetRAGService` (`app.go:154-156`) but live path bypasses it | Broken integration |
| S6 Image | S8 model registry/loading | `app-image-model.go`; `internal/image/engine/sdcpp` | Model discovery dependency |
| S6 Image | S9 resource arbitration | `app-vram-arbitration.go`; `image-generator-core.js:594-607` | Advisory dependency |
| S6 Image | S10 image settings, S11 gallery | `app-image.go:30-120`; `gallery.Service` | Persistence dependency |
| S7 Audio/TTS | S1 process ownership | `internal/audio/tts/helper/supervisor.go` | Lifecycle dependency |
| S7 Audio/TTS | S9 resource model | `resource.ResourceTTS` exists; no arbitration kind | Incomplete dependency |
| S7 Audio/TTS | S10 audio settings | `app-audio.go:338-505` | Configuration dependency |
| S8 Model loading | S1 process lifecycle, S9 hardware | `app-model.go:250-370`; `app.go:283-315` | Ownership/resource dependency |
| S9 Resource | S8 model registration, S6 image load | `app-model.go:504`; `app-image-model.go:143` | State dependency |
| S10 Settings | S11 settings files | `internal/settings/store.go`; image settings store | Persistence dependency |
| S12 UI | S13 bindings/events | `frontend/wailsjs`; all page scripts | Transport dependency |
| S13 Events | S1 lifecycle | `emitEvent` context guard | Teardown dependency |
| S14 Errors | All subsystems | `logError` calls throughout | Cross-cutting |

### 6.2 Shared state and contention

| Shared state | Owners | Contention / risk |
|---|---|---|
| Loaded model records | `resource.Manager`, `App.serverProc`, `App.imageEngine`, `chat.model.Manager` | Multiple truths; no reconciliation on failure/switch |
| Chat history | `App.chatMessages`, `chat.Service.context`, `history.ndjson`, frontend DOM | Duplicate writes and reload divergence |
| Character/world/scenario state | `App` fields, `chat.Service` fields, frontend state, disk | Two-way projection |
| Lorebook state | Backend `lorebooks.json`, frontend `localStorage`, `chat.Service.manualLore` | No canonical copy |
| Stream cancellation | `App.streamCancel` | One slot for all streams |
| Audio generation | `audio.Manager` active/retired slots | Well-owned internally; not closed on normal shutdown |
| GPU/VRAM | `resource.Manager`, `sd-server`, `llama-server`, driver | Advisory checks, no reservation |
| Settings | `settings.Store`, image store, frontend state | Mixed atomicity and ownership |

### 6.3 Indirect and circular coupling

- **Runtime cycle:** `App` owns `chat.Service`, while `chat.Service` calls back into application-level concepts only through data; there is no Go import cycle, but state is duplicated.
- **Persistence cycle:** chat persistence stores characters/lorebooks, while character/lorebook screens also read/write project folders and `localStorage`.
- **Lifecycle cycle:** model loading needs resource admission; resource admission needs to know what is loaded; model loading registers itself before it is loaded. This is a state-ordering cycle, not a package cycle.

### 6.4 Independent components

- GGUF admission and projector pairing are largely self-contained and testable without the UI.
- Image metadata extraction and gallery file operations are mostly independent of chat.
- Hardware sampling is independent of model state.
- Error logging is independent of feature state, though its coverage is uneven.

---

## 7. TTS findings applicability

The TTS investigation is a diagnostic reference, not a mandate to re-audit TTS. The table below maps each TTS concern to the rest of the application.

| TTS concern | Analogous mechanism elsewhere | Applicability |
|---|---|---|
| Generation isolation | Image jobs have a job ID and owner-aware load abort (`supervisor.go:50-95`); audio manager has generations and leases (`manager.go:40-220`); chat has one global `streamCancel` | **Shared gap in chat.** Image/audio are better isolated. TTS is not the only place generation ownership is weak. |
| Process ownership | llama uses `proc.Start`; sd-server uses `exec.Command` + `proc.Attach`; TTS helper design uses `proc.Start` | **Partially shared.** Ownership is inconsistent across child types. The TTS design’s job model is stronger than the current sd-server attach path. |
| Shutdown re-entry | `OnShutdown` exists but is unregistered; audio `Close` has terminal state; image `Shutdown` is idempotent-ish but not proven | **Shared defect at the application boundary.** TTS would be shut down never, exactly like image and llama, on normal window close. |
| Terminal-state immutability | Audio manager has terminal close result; image engine resets generation on close; chat has no terminal request state | **Mostly TTS/audio-specific.** Chat lacks the concept entirely. |
| Action and exit causality | Image engine distinguishes `errProcessExited` from timeout and cancellation (`supervisor.go:16-24`, `:150-170`); llama monitor uses progress stages | **Image is the model to reuse.** TTS’s causality problem is not universal, but chat’s abort/error classification is weaker. |
| Reader lifecycle | Image has one `Wait` owner and a single stdout pump; llama uses an HTTP scanner; TTS uses NDJSON reader | **TTS-specific protocol risk.** The single-reader discipline is reusable. |
| Protocol compatibility | TTS helper protocol v1 has explicit capability/version negotiation; image/llama use HTTP endpoints | **TTS-specific.** Do not generalize protocol machinery to HTTP engines. |
| State projection | Audio `GetAudioStatus` distinguishes desired, applied, and capabilities (`app-audio.go:230-337`); chat projects UI state from multiple stores | **Shared lesson:** publish owned state, not inferred state. Audio already does this better than chat. |

**Conclusion:** the TTS defects are a mix of (a) TTS-specific protocol/helper lifecycle, and (b) shared ownership failures that also affect chat, image, and application shutdown. Fixing only TTS would leave the shared lifecycle and cancellation defects in place. Conversely, a shared lifecycle/cancellation foundation would remove some of the conditions that make TTS unsafe without deciding TTS’s unapproved contract.

---

## 8. Confirmed architectural risks

Risk identifiers are for engineering tracking only. “Confirmed” means source or documented evidence supports the condition; “inference” means the runtime consequence is reasoned, not observed.

| ID | Classification | Finding | Source | Affected | Consequence | Handle before subsystem completion? |
|---|---|---|---|---|---|---|
| R1 | Source fact | `OnShutdown` is not registered | `main.go:22-29`; `app.go:283-315` | S1, S2, S6, S7, S11 | Final save/history rewrite, model kill, image shutdown, audio drain, and `proc.CloseJob` do not run on normal close | **Yes — shared foundation** |
| R2 | Source fact | Chat root and entity roots differ | `persistence.go:41-56`; `app-characters.go:50-90`; `service.go:88` | S3, S4, S5, S11 | Data appears missing after install/portable move; backups/wipe miss data | **Yes — shared foundation** |
| R3 | Source fact | `WipeAllData` deletes `<exe>/chats`, not the actual chat root | `app-wipe.go:20-28`; `persistence.go:41-56` | S11, S2 | “Wipe all data” leaves chat history and config behind | **Yes — persistence** |
| R4 | Source fact | One `streamCancel` for all streams | `app.go:48`; `app-prompt.go:61`, `:115`, `:205`, `:245-255` | S2, S13 | Abort can cancel the wrong reply; concurrent character replies cannot be independently stopped | **Yes — generation ownership** |
| R5 | Source fact | Live chat bypasses `chat.Service`; RAG is attached to the unused service | `app-prompt.go:25-147`; `app.go:139-156`; `app-audio.go:565-576` | S2, S5 | RAG/lore state can be configured without affecting the prompt; duplicate prompt logic drifts | **Yes — chat foundation** |
| R6 | Source fact | Resource registration is not authoritative | `app-model.go:504-532`; `app-image-model.go:143-150`; `app-vram-arbitration.go:47-120` | S6, S8, S9 | Fit checks and “loaded models” can be wrong; TTS cannot participate | **Yes — resource foundation** |
| R7 | Source fact | Settings writes are mixed atomicity | `internal/settings/store.go` (`Save`); `persistence.go:107-117`; `app-characters.go:68-75`; `internal/image/settings/settings.go:44-55` | S3, S4, S10, S11 | Crash during direct write can corrupt config or entity JSON | **Yes — persistence** |
| R8 | Source fact | Frontend and backend both own lorebook state | `chat-stream.js:31-42`; `persistence.go` lorebooks; `chat.Service.manualLore` | S5, S11 | Stale or divergent lore between UI and prompt | **Yes — lore foundation** |
| R9 | Documented test evidence | T0.3 image model load/path failure and T0.4 preview remain open | `ISSUES.md` Tier 0 | S6 | Image page cannot be considered complete or judged on layout | Within image subsystem, after shared lifecycle |
| R10 | Documented test evidence | T0.5 path-based orphan reaping kills concurrent instances’ children | `ISSUES.md` T0.5; `internal/proc/job_windows.go:110-136` | S1, S6, S8 | Two app instances or a harness can kill each other’s model processes | **Yes — process ownership** |
| R11 | Documented evidence | T2.1 theme and T2.2 ghost IDs remain open | `ISSUES.md`; `STAGE_STATUS.md` T2.2 | S12 | Light theme unreadable; handlers may target missing elements | UI foundation, after functional boundaries |
| R12 | Inference | Concurrent autosave/final-save and append can interleave | `app-prompt.go:141-147`; `app-chat.go:310-350`; `persistence.go:678-710` | S2, S11 | Duplicate or lost tail entries under overlap or power loss | Persistence foundation; verify with concurrency test later |
| R13 | Inference | Unlocked `App` chat/character fields can be mutated concurrently | `app.go:40-60`; `app-prompt.go:25-205`; `app-characters.go:41-92` | S2, S3, S13 | Stale prompt context, inconsistent character list, abort confusion | Generation/state foundation |
| R14 | Unknown | Runtime extent of image/audio/model failures | No runtime evidence | S6, S7, S8 | Cannot claim working end-to-end | Resolve during subsystem verification |

### 8.1 Prioritization

**Tier A — shared, broad reach, must precede feature completion**

- R1 lifecycle registration and shutdown order.
- R2/R3/R7 persistence root, wipe, and atomic write contract.
- R4/R13 generation and cancellation ownership.
- R5 single chat/prompt/RAG architecture.
- R6 resource authority.
- R10 process ownership boundary.

**Tier B — shared reliability, can follow Tier A**

- R8 lorebook canonical state.
- R12 save serialization.
- R11 shared UI/theme/ghost-ID cleanup.

**Tier C — subsystem-specific**

- R9 image blockers.
- TTS protocol/lifecycle details under the frozen governance.
- Story tools scope.

---

## 9. Shared-foundation engineering recommendations

These are work areas, not implementation authorization. Each area has a narrow purpose and an explicit non-goal.

### F1. Lifecycle and process ownership

**Why shared:** every long-running child (llama, sd-server, TTS/ASR helpers, video) depends on one shutdown and ownership contract.

**Dependents:** S1, S2, S6, S7, S8, S11.

**What to establish:**

- Register `OnShutdown` in `main.go` with the same handler `ExitApp` uses.
- Define one shutdown order and make it idempotent: stop accepting new work → cancel active generations → save durable chat state → stop model servers → stop image engine → drain audio → close job.
- Make every child start under recorded ownership. `sd-server` should not rely on attach-after-start if suspended start is available. `ReapOrphans` must require proven ownership, not executable-path prefix.
- Ensure window close and explicit exit use the same path.

**Verification later:** two-instance test (no cross-kill), window-close test (no surviving children), final-save test (history rewritten), audio-drain test.

**Outside scope:** redesigning the process model, adding a service manager, or changing TTS governance.

### F2. Persistence and path contract

**Why shared:** chat, characters, worlds, scenarios, lore, settings, gallery, and wipe all depend on where data lives and how it is written.

**Dependents:** S2, S3, S4, S5, S6, S10, S11.

**What to establish:**

- One canonical data root, or an explicit portable-vs-installed policy with one resolver used by every subsystem.
- One migration path for existing split data.
- One atomic-write helper for JSON and settings; no direct `os.WriteFile` for durable state.
- A real wipe contract that removes the actual chat root and all settings, or explicitly preserves selected assets.
- A save serialization boundary for chat history/autosave/final save.

**Verification later:** path matrix for dev/installed/portable; crash-during-write test; wipe-and-relaunch test; import/export round-trip.

**Outside scope:** changing the data schema, adding a database, or designing cloud sync.

### F3. Generation, request, and cancellation ownership

**Why shared:** chat streams, character streams, image jobs, and audio operations all need to know which request owns which cancellation and which completion event.

**Dependents:** S2, S6, S7, S13.

**What to establish:**

- A request/generation ID carried through Wails events and status queries.
- Per-request cancel handles, not one `streamCancel`.
- A rule that a completion from generation N cannot clear or overwrite generation N+1.
- Stale-event rejection in the UI.
- Shared vocabulary for cancelled vs failed vs completed.

**Verification later:** two concurrent character streams, abort one, confirm the other completes; image job cancel isolated from a second job; late event from an old stream ignored.

**Outside scope:** redesigning the chat UI or adding multi-user sessions.

### F4. Resource admission and model-state authority

**Why shared:** chat, image, and audio all consume VRAM/RAM and all need one answer to “what is loaded?”

**Dependents:** S6, S7, S8, S9.

**What to establish:**

- Register only after a model is ready; unregister on failure, switch, unload, and process exit.
- Make image auto-load register the same way manual load does.
- Add TTS as an admission participant, or explicitly document why TTS is outside VRAM arbitration.
- Make backend admission mandatory, not only a frontend advisory check.
- Reconcile `resource.Manager`, `serverProc`, `imageEngine`, and actual HTTP/process liveness.

**Verification later:** failed load leaves no ghost; switch leaves one entry; image job appears in summary; TTS allocation is accounted for; forced low-VRAM path behaves as specified.

**Outside scope:** rewriting the allocator algorithm or adding a scheduler.

### F5. Settings, error, and shared UI consistency

**Why shared:** every subsystem reads settings and reports failures through the same pages and event surface.

**Dependents:** all subsystems, especially S10, S12, S14.

**What to establish:**

- One atomic settings writer and one ownership rule for each settings section.
- A single user-facing error contract: code, operation, recovery action, and whether retry is safe.
- Replace silent `_ =` persistence failures with visible state.
- Shared notification/dialog components where behavior is currently page-local.
- Close T2.1/T2.2 as part of functional UI consistency, not a cosmetic redesign.

**Verification later:** restart round-trip for every setting section; induced write failure shows a user-visible error; theme screenshots in both modes; ghost-ID inventory reaches zero write-to-missing-ID.

**Outside scope:** visual redesign, new design system, accessibility overhaul beyond functional readability.

---

## 10. Recommended engineering order

### 10.1 Shared-foundation order

| Order | Work area | Depends on | Why now |
|---:|---|---|---|
| 0 | Freeze TTS governance; no C1 implementation | — | Prevents an unapproved subsystem from driving foundation changes |
| 1 | F1 Lifecycle and process ownership | — | Every other foundation assumes a reliable shutdown and child owner |
| 2 | F2 Persistence and path contract | F1 for final-save behavior | Data loss and wipe correctness outrank feature work |
| 3 | F3 Generation/request/cancellation ownership | F1 | Chat/image/audio concurrency cannot be completed safely without it |
| 4 | F4 Resource admission and model-state authority | F1, F3, S8/S6 behavior | Prevents ghost models and contradictory fit decisions |
| 5 | F5 Settings, error, and UI consistency | F1-F4 | Shared surfaces can now reflect owned state instead of projections |

### 10.2 Subsystem completion order after foundations

1. **S8 Model discovery/loading** — foundational for chat and image; verify admission, load, unload, failure recovery.
2. **S3 Characters and S4 Worlds/scenarios** — CRUD and path contract; low concurrency, high data-integrity value.
3. **S2 Chat** — prompt assembly, streaming, history, vision-input current-turn behavior, cancel, save/reload.
4. **S5 Lore/memory/RAG** — one canonical lore state and a prompt path that actually consumes retrieval.
5. **S6 Image** — close T0.3/T0.4, verify job ownership, gallery, and resource registration.
6. **S7 Audio/ASR/TTS** — only after governance changes; reuse F1/F3/F4 contracts rather than creating a parallel lifecycle.
7. **S12/S15 UI and story tools** — finish remaining controls and unknown scope after functional subsystems are stable.

**Do not** start with TTS because its audit is most detailed. Detail is not dependency order.

---

## 11. Proposed subsystem-completion workflow

For each subsystem, once its foundations are stable:

1. **Establish intended functionality with the user.** Capture the complete product requirement, including advanced functions not currently in the UI. Do not infer “never intended” from absence.
2. **Inventory every control.** Buttons, menus, dials, sliders, toggles, dialogs, context menus, keyboard actions, drag/drop, and empty states.
3. **Map each control to a backend binding and a durable effect.** Mark unconnected, no-op, duplicate, or stale controls.
4. **Define the subsystem’s data contract.** Ownership, canonical path, schema, migration, deletion, export.
5. **Complete missing functionality.** No placeholders, no silent no-ops, no “coming soon” controls.
6. **Correct confirmed defects.** Fix causes, not symptoms; preserve working behavior.
7. **Audit lifecycle and concurrency.** Start/stop, cancel, retry, re-entry, stale completion, shutdown.
8. **Audit persistence and failure recovery.** Partial write, corrupt file, missing asset, permission denied, disk full.
9. **Audit resource behavior.** Registration, release, contention, failure when monitoring is unavailable.
10. **Polish UI and interaction.** Consistent states, disabled states, loading, errors, readability, theme.
11. **Verify end to end.** Use the subsystem’s own evidence standard: logs, files, counts, and screenshots where visual.
12. **Confirm integration with completed subsystems.** Re-run the relevant cross-subsystem path, not only the subsystem’s isolated test.
13. **Declare completion only with evidence.** Record what was verified, what was not, and what remains.

The workflow is deliberately front-to-back. “The headline feature works” is not completion.

---

## 12. Unresolved questions requiring user input

These are genuine product/ownership decisions that source inspection cannot settle.

1. **Canonical data location.** Should Howl be portable (all data beside the executable), installed per user (`%APPDATA%`), or support both with an explicit mode? This determines the migration, wipe, backup, and path contract.
2. **Wipe semantics.** Should “Wipe All Data” remove settings, model registry entries, logs, and downloaded models, or preserve assets and only remove creative work? The current implementation does neither consistently.
3. **Concurrent character replies.** Are multiple character responses intended to run simultaneously, or strictly sequential turns? This determines whether F3 needs multi-generation cancellation or only one active turn.
4. **Story tools scope.** What is `story-tools.html` intended to do, and which of its controls are committed versus exploratory? Current evidence is insufficient to classify it.
5. **TTS governance.** P17, contract rev3, OD-9, and C1 remain frozen. No implementation question should be answered by this audit; the owner must separately decide when that changes.

---

## 13. Final recommendation

Do not begin subsystem completion with TTS, image polish, or chat features.

First make the application’s ownership boundaries true:

1. Register and exercise one shutdown path.
2. Establish one persistence root, one atomic write contract, and a wipe that matches reality.
3. Give every generation an owner and every cancellation a target.
4. Make the resource manager the authoritative model-state record, or explicitly demote it to telemetry and build a real admission boundary.
5. Collapse the duplicate chat/RAG architecture before adding chat features.

After those foundations, complete subsystems in dependency order: models, characters/worlds, chat, lore/memory, image, audio/TTS, then UI/story tools. This order follows actual dependencies and data risk, not audit recency.

The application is not a blank slate and does not need a rewrite. It needs its existing foundations to stop disagreeing with each other.

---

*End of assessment.*
