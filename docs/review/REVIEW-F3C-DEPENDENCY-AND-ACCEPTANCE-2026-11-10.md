# REVIEW — F3-C Worlds and Scenarios: Dependency and Acceptance

**Date:** 2026-11-10
**Scope:** dependency and acceptance review only, before any F3-C implementation
**Not in scope:** code changes, migrations, schema finalization, F3-D/E/F/G/H, mobile work
**Method:** source reading of the current tree plus the accepted F3-B/W27 records. No runtime or
build test was executed for this review, and none is claimed.

## 1. Verdict

F3-C is **unblocked and not started**. Its dependency chain is sound: the ownership contracts,
registry roles, isolation, and native evidence model all exist and were exercised by F3-B/W27.
What does not exist is any canonical world or scenario store. Worlds and scenarios are still
browser-owned, with write-only native copies, name-keyed chat reads, and snapshots that carry no
revision. That is F3-C's actual work.

Six product decisions must be made before implementation begins. They are §6.

## 2. What F3-C inherits from F3-B

These are reusable contracts, not just precedent. Each was exercised in F3-B or W27.

| Inherited contract | Where it lives | Reuse in F3-C |
| --- | --- | --- |
| One canonical store owns the files; browser storage becomes a rebuildable cache | `internal/lorebookstore`, `frontend/dist/js/lorebook-repository.js` | Same shape for worlds/scenarios: a store package plus one frontend data path |
| Identity names the path; a rename never moves the file | `lorebookstore`: `entryID`, `encodeDocument`, ID-based filenames | World/scenario IDs must be independent of title from the first write |
| Create requires a save root; replace requires the write role on the root that already holds the item | `Store.Save`, `library.Registry.WriteFile` | Reuse unchanged; no new permission model |
| Explicit, create-only, idempotent, journaled, reversible import | `Store.Import`, `RollbackImport`, journal under the settings dir | Same contract, including "legacy store untouched" |
| Cutover marker, not a silent switch | `lorebook-cutover.json` | A world/scenario cutover marker per category |
| Bounded per-root scan, conflict reporting, deterministic winner | `Store.scanRoot`, list conflicts | Same rules; a duplicated world ID resolves to one winner and is reported |
| Disconnected roots keep registration and last-known items marked unavailable | Availability flags plus the index cache | Required for F3-C acceptance, not optional |
| Atomic file writes, unsafe-target rejection, structured save results | `library` registry errors, `writeAtomic` | Reuse; no new error vocabulary |
| Native acceptance against a marked isolated root; test seam compile-time absent from production | `HOWL_WEBVIEW_DATA_DIR` + `native_test_seam` + `promptcapture` | Same method and the same evidence-level separation |
| Prompt-boundary evidence is not model evidence | W27 §27.4 | F3-C must state its own evidence level |

**REG-1 is unchanged:** world and scenario attachments already reference lorebook identities. F3-C
must preserve those references, not rewrite them. The conflict-grouping fix is neutral here and is
pinned by `TestConflictGroupingKeepsNonCharacterScopesShared`.

## 3. Shared foundations already complete

| Foundation | State | Note for F3-C |
| --- | --- | --- |
| Library registry with `worlds` and `scenarios` categories, Read/Write/Save roles, default-save designation | Complete | No registry work needed for these categories |
| Approved default destinations (`Documents\Howl\Worlds`, `...\Scenarios`) and isolation equivalents | Complete | `paths.DefaultContentDir` already covers both |
| Path containment, unsafe-target rejection, collision and read-only errors | Complete | Reused by the world/scenario store |
| Isolated native runtime with marker root, fail-loud validation | Complete | Native acceptance can run without touching user data |
| Frontend asset injection for every page | Complete | World/scenario pages can use bindings like the Lorebooks page does |
| Prompt-capture seam and character prompt composer | Complete | Usable for the prompt-boundary part of F3-C acceptance |
| Lorebook canonical identities and owner privacy | Complete | World/scenario lorebook references can resolve by ID today |
| Stable character IDs | Complete | Scenario character references can resolve by ID, with dangling references reported |

## 4. Current state F3-C must replace

| Area | Current reality (source-confirmed) | Gap |
| --- | --- | --- |
| Primary ownership | `howl-worlds` and `howl-chat-scenarios` localStorage; images in IndexedDB | No canonical owner; browser data is still the source |
| Native copies | `ExportContentItem` can write `howl-content-export-v1` with `authoritative: false`; no frontend caller was found | The export path exists but is unused and explicitly non-authoritative |
| Native directories | `getWorldsDir`/`getScenariosDir` resolve from cwd or the executable directory | Not a registered root; not user-owned; same defect F2-C flagged |
| Chat reads | `chat.Service` loads `Worlds/<name>/world.json`, `world_lore.json`, `automated_world_lore.json` (and scenario equivalents) by name | Name is the de facto key; renames break resolution |
| Chat snapshots | Each chat embeds `world.json`, `scenario.json`, `characters.json`, `lorebooks.json`; `ChatMeta` has `WorldID`/`ScenarioID` but no revision | Staleness is invisible; no visible revision policy |
| References | `ScenarioContext` carries character IDs and `lorebooks []interface{}` | Dangling references have no defined display or rejection rule |
| World/scenario lore | Entries use `ScopeWorld`/`ScopeScenario` with `OwnerID` equal to the name | Name-keyed ownership survives in a system that otherwise uses stable IDs |

The last row is the direct analogue of the W27 defect class: correct privacy logic sitting on a name
key. F3-C must not repeat the circularity of "the world is identified by the name that the file is
stored under."

## 5. Dependencies and ordering

| Dependency | Blocks F3-C? | Why |
| --- | --- | --- |
| F3-E chats | No | The snapshot/revision policy can be defined and surfaced in the existing chat folder format; moving the chat root is separate |
| F3-F images/media | Partly | Media references can be defined now, but canonical image ownership and atomic PNG+sidecar commit belong to F3-F. F3-C should define the reference shape and keep today's image resolution |
| F3-D characters | No | Scenario character references use existing stable IDs; unresolved IDs must be displayed as unavailable |
| Lorebook canonical store | No | Already complete; F3-C consumes it |
| Mobile | No | Explicitly desktop-independent; see the ledger constraint |

Recommended order once authorized: decisions (§6) → world/scenario store with identity and roles →
page read path and cache → explicit import plus rollback → reference resolution → snapshot revision
policy → native acceptance.

## 6. Decisions required before implementation

| ID | Question | Consequence of deferring |
| --- | --- | --- |
| F3-C-Q1 | Document shape: `<root>/<id>/world.json` plus assets, or one flat `<root>/<id>.json` with media referenced elsewhere | Determines whether the media contract is inside or outside F3-C |
| F3-C-Q2 | Identity: mint `wld_`/`scn_` IDs on import, or preserve the existing localStorage IDs | Affects REG-1 attachments, chat `WorldID`/`ScenarioID`, and rollback |
| F3-C-Q3 | Import unit and authority: page-level explicit import of all browser worlds/scenarios, with the same journal/rollback contract as lorebooks | Without it, cutover cannot be reversed or verified |
| F3-C-Q4 | Snapshot revision policy: what the chat records (source ID + content revision hash), and how "this chat uses an older world" is shown | Staleness stays invisible, which is the current defect |
| F3-C-Q5 | Dangling references: reject on save, or display as unavailable and never substitute | Silent substitution would violate the ownership rule |
| F3-C-Q6 | World/scenario lore ownership: migrate `OwnerID` from name to stable ID in this slice, or record it as a named limitation | Name-keyed lore would survive the slice that exists to remove name keys |

## 7. Acceptance for the slice

Acceptance requires all of the following, with evidence levels kept separate.

1. **Identity:** two worlds and two scenarios created, each with a stable ID; rename keeps ID and
   path; read-back after restart returns the same IDs.
2. **Single writer:** only the world/scenario store writes item files; create is create-only;
   replace requires the write role on the owning root; collision, read-only, and unsafe-target
   failures are surfaced, never reported as success.
3. **Explicit import and rollback:** browser data imports through an explicit action, is idempotent
   by identity plus content hash, is journaled, and rolls back to an untouched legacy store.
4. **Canonical read path:** the world/scenario pages read through the canonical store; the browser
   cache is rebuildable; no second source after cutover.
5. **References:** lorebook references resolve by canonical identity; missing or foreign references
   are displayed as unavailable and never substituted; scenario character references behave the same.
6. **Snapshot policy:** a chat records the source ID and a visible revision; a chat using an older
   world or scenario says so.
7. **Disconnected roots:** registration and last-known entries survive, are marked unavailable, and
   are excluded from prompt assembly.
8. **Native evidence:** the scenario runs in the real Wails window against a marked isolated root,
   with the character/world prompt boundary captured and the production binary free of test seams.
9. **Regression:** the lorebook identity/privacy suites, the frontend identity tests, and the
   conflict scope-boundary test all still pass.
10. **Limits stated:** no model inference is claimed unless a model is actually loaded; image
    ownership, chat-root moves, and character-card canonicalization remain with their own slices.

## 8. Risks and non-goals

- **Risk — snapshot divergence.** The chat embeds full documents. A policy that records revisions
  but does not surface them would satisfy nothing; the visible indicator is part of acceptance.
- **Risk — identity migration breadth.** Preserving localStorage IDs keeps attachments but carries
  whatever shape those IDs already have; minting IDs requires rewriting references. This is a
  decision, not an implementation detail.
- **Risk — media scope creep.** Image ownership is F3-F. F3-C should resolve images through a
  defined reference and defer the media store.
- **Non-goals:** no character, chat, image, video, or wipe work; no removal of legacy browser data;
  no automatic migration; no mobile work.

## 9. Recommended next action

Authorize F3-C-Q1 through F3-C-Q6, then implement in the order in §5. Do not begin with the page
rewrite: the store and identity decisions determine whether that rewrite has to be done twice.
