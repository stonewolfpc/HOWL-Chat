# AUDIT — F3-C Repository-Wide Single Writer and Legacy Readers

**Date:** 2026-11-10
**Status:** source audit complete; two legacy surfaces removed/blocked; A10 source-level and native
runtime acceptance pass
**Scope:** production Go and frontend paths that read or write world/scenario content, documents, or
item folders
**Related:** `FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md` §3, §6; A6, A10, A11, A13;
`F3C-LEGACY-LORE-CLASSIFICATION-2026-11-10.md`

## 1. Method

Searches covered:

- path literals `Worlds`, `Scenarios`, `world.json`, `scenario.json`, `world_lore.json`,
  `scenario_lore.json`, `automated_world_lore.json`, `automated_scenario_lore.json`;
- directory helpers and existence checks (`getWorldsDir`, `getScenariosDir`, `worldExistsOnDisk`,
  `scenarioExistsOnDisk`);
- category writers (`worldstore`, `library.Registry.WriteFile/WriteNewFile`,
  `ExportContentItem`, `os.WriteFile`, `os.RemoveAll`);
- app/service consumers (`LoadChatData`, `AddWorldToChat`, `performSave`, `LoadChat`,
  `buildPromptMessages`);
- frontend storage (`howl-worlds`, `howl-chat-scenarios`, world/scenario managers);
- adjacent packages (`lorebook`, `rag`, `memory`, `lorebookstore`) for file reads or writes of
  world/scenario folders.

Tests, fixtures, and documentation are excluded from the production inventory. Test-only writes are
noted where they exercise a production path.

## 2. Read paths

| # | Path | What it reads | Classification |
| --- | --- | --- | --- |
| 1 | `worldstore.Lookup` / `Load` | Canonical envelope by stable ID across every registered root | Canonical reader |
| 2 | `worldstore.LoadLore` | Allowlisted lore file inside a canonical item folder | Canonical reader, read-only |
| 3 | `chat.Service.loadCanonicalWorld` / `loadCanonicalScenario` | Documents and lore through the store | Canonical consumer |
| 4 | `chat.Service.loadItemLore` | Entries from canonical lore files; classifies in memory | Canonical consumer |
| 5 | `App.syncAppStateFromChatService` | Service world/scenario data (test seam) | Canonical via consumer |
| 6 | `App.buildPromptMessages` | `a.world` / `a.scenario` from the active chat snapshot or `AddWorldToChat` | Chat-owned snapshot; not a direct canonical read |
| 7 | `PersistenceManager.LoadChat` | `Chats/<id>/world.json` and `scenario.json` | Chat-owned snapshot |
| 8 | `App.loadChatDataFromStorage` / `applyChatDataToState` | Chat snapshot | Chat-owned snapshot |
| 9 | Frontend `localStorage` (`howl-worlds`, `howl-chat-scenarios`) | Browser world/scenario records and attachments | Browser-owned legacy records; import source |
| 10 | `App.GetContentLibraryCharacters` | Character roots only | Not a world/scenario reader |
| 11 | `lorebook`, `rag`, `memory` packages | Entries, labels, and chat memory; no item-folder file reads found | Not world/scenario item readers |
| 12 | Removed: `worldExistsOnDisk`, `scenarioExistsOnDisk`, `getWorldsDir`, `getScenariosDir` | Old `Worlds/<name>` / `Scenarios/<name>` existence checks | **Removed** — no production caller remained |

`worldExistsOnDisk` and `scenarioExistsOnDisk` were dead code with no callers, but they were still
legacy name-keyed readers in the repository. They are removed, together with their directory
helpers.

## 3. Write paths

| # | Path | What it writes | Classification |
| --- | --- | --- | --- |
| 1 | `worldstore.Save` | New canonical `<root>/<id>/world.json` or `scenario.json` | **Only canonical document writer** |
| 2 | `worldstore.Update` | Atomic replacement of an existing canonical document | **Only canonical document writer** |
| 3 | `library.Registry.WriteFile` / `WriteNewFile` | Low-level atomic file writer used by the store | Infrastructure; no direct world/scenario caller remains |
| 4 | `chat.Service.SaveChatData` | `chat.json`, `history.json`, memory snapshot; character automated lore only | Chat-owned; generated world/scenario branches removed |
| 5 | `updateAutomatedLoreFile` | Character automated lore, reachable only through `automatedLoreSyncTarget` | Character-owned; F3-D |
| 6 | `PersistenceManager.SaveChat` | Chat snapshot documents inside `Chats/<id>/` | Chat-owned |
| 7 | `App.performSave` / `saveChatDataToStorage` | Chat persistence | Chat-owned |
| 8 | Frontend world/scenario managers and grids | `localStorage` browser records and attachments | Browser-owned legacy records; not item folders |
| 9 | `App.ExportContentItem` | Non-authoritative export copies for characters/lorebooks | Worlds/scenarios now **refused** |
| 10 | `WipeAllData` | Chat root and gallery root only; no `Worlds`/`Scenarios` deletion | Not a world/scenario writer |
| 11 | Memory and RAG services | Memory/config and RAG storage | Not item folders |

### Closed bypass

`ExportContentItem` previously accepted `worlds` and `scenarios` and wrote a
`howl-content-export-v1` envelope to `<hash>/world.json` or `<hash>/scenario.json` under a
registered root. It was explicitly non-authoritative, but it was still a second writer of a
world/scenario filename into an item root. A10 says only the store writes item folders. The binding
now refuses both categories before any destination lookup:

> `worlds export is not available: canonical items are written only by the world/scenario store`

`TestExportContentItemRefusesCanonicalWorldAndScenarioRoots` asserts both categories are refused and
no target is produced. Generic export tests now use characters. World/scenario export is a
deliberate capability gap until a store-backed export exists.

## 4. What is not a writer

- `AddWorldToChat` accepts a frontend map and updates `a.world`; the only persistence it triggers is
  `performSave`, which writes the chat snapshot. It does not write item folders.
- `buildPromptMessages` reads app state; it does not write content.
- `WipeAllData` stages and removes `chats` and the gallery only. It never deletes canonical
  world/scenario roots.
- `lorebookstore` owns lorebook content only; its writes do not target world/scenario item folders.
- `internal/backend/memory` writes memory state, not `Worlds` or `Scenarios`.
- No production `os.Remove`/`os.RemoveAll` call targets a world/scenario root.

## 5. Findings

1. **Dead legacy readers existed.** Removed.
2. **Generic export was a real bypass.** Refused for worlds/scenarios.
3. **The desktop editor still uses browser records.** Worlds and scenarios are created and edited in
   `localStorage` (`howl-worlds`, `howl-chat-scenarios`). That is the active desktop data path and
   the intended import source. It is not a canonical item-folder writer, but it is a separate
   content authority until import and cutover.
4. **The prompt builder is snapshot-based.** `chat.Service` is canonical, but the desktop prompt
   path consumes `a.world`/`a.scenario` from the chat snapshot. The service conversion does not
   change that path. This remains a deliberate separation, not an accidental canonical read.
5. **Old name-keyed lore now has one authorized migration path.** At audit time it had none. LF-1
   (`DECISION-F3C-LEGACY-FOLDER-MIGRATION-2026-11-10.md`) amends §6 to allow a single explicit,
   journaled, reversible migration by verified identity or user resolution. Runtime name reads and
   writes remain prohibited; unresolved content is quarantined. Migration code is not implemented,
   so category 2 remains unread in production.

## 6. A10 verdict

**Source-level canonical single-writer: PASS.**

- `worldstore.Save`/`Update` are the only writers of canonical envelopes.
- `chat.Service` cannot write under `Worlds/` or `Scenarios/`; its generated-lore target is
  character-only.
- The generic export bypass is refused.
- No other production writer of canonical item-folder `world.json`, `scenario.json`, or
  world/scenario lore files remains. Chat snapshot files under `Chats/<id>/` are a separate,
  chat-owned contract.
- No production deletion path targets those roots.

**A10 acceptance: CLOSED — PASS.**

The frozen evidence requires a source audit **and** a runtime check that chat cannot write item
folders. The native A10 run created canonical world and scenario items through the real import
bindings, loaded chat data by stable IDs, called `LoadChat`, `SaveChat`, and `SaveChatData`, then ran
`chat.Service.SaveChatData(..., true)` with world and scenario generated-lore entries present. The
canonical world and scenario trees were byte-identical before and after every step; chat-owned files
were written only under `Chats/A10Chat`. Evidence:
`VERIFICATION-F3C-NATIVE-A9-A10-2026-11-10.md`.

A6 (no chat writes under item roots) has source, unit, and native evidence. A9's native import run is
recorded separately. A11 (unavailable content excluded) has consumer tests. A13 passes. A10 no longer
stays open.

## 7. Verification for this audit

Commands run after the removals and export guard:

| Check | Result |
| --- | --- |
| `gofmt` on changed files | clean |
| `go vet ./internal/worldstore ./internal/backend/chat .` | exit 0 |
| `go test ./internal/worldstore ./internal/backend/chat . -count=1` | all `ok`, exit 0 |
| Worldstore tests | 38 pass |
| Chat tests | 23 pass |
| New no-mutation lore test | pass |
| New name-folder non-fallback test | pass |
| New export-bypass test | pass |
| Source grep for `"Worlds"` / `"Scenarios"` in `internal/backend/chat` | none |
| Source grep for canonical item-folder writers outside `worldstore` | none |

## 8. Residual risk

- `library.Registry` remains a general atomic writer. A future caller could use it without the
  store; the guard is code review, tests, and the export restriction, not a type-level barrier.
- `worldstore` has no lore write API yet. Import will need one, with the same authority and atomic
  rules as documents.
- Browser records remain writable and are not yet mirrored. Until import/cutover, a desktop edit is
  not canonical and a canonical edit is not reflected in the browser editor.
- The old name-folder migration is authorized by LF-1 but not implemented. Until it exists,
  category 2 remains inert and unread. Runtime name reads and writes remain prohibited.
