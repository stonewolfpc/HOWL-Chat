# IMPLEMENTATION — F3-C Chat Consumer Conversion

**Date:** 2026-11-10
**Status:** converted; focused, regression, and native A10 gates pass
**Scope:** `internal/backend/chat` world/scenario reads and writes, the native test seam, and the
chat persistence identity backfill
**Authorization:** explicit user approval after the root-authority gate passed; chat conversion was
the next authorized implementation
**Contracts:** `FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md` (R1, R2, A10, A11, A13),
`DECISION-F3C-ROOT-AUTHORITY-2026-11-10.md`

## 1. What changed

`chat.Service` is now a consumer of canonical world/scenario content:

- World and scenario documents are loaded by **stable ID** through `worldstore`; the display name is
  never used to locate a folder.
- Manual and generated item lore are read from the canonical item folder through
  `worldstore.LoadLore` by the same ID.
- Automated world/scenario lore is classified as `legacy_generated_unattributed` and remains
  prompt-eligible, with no chat attribution invented.
- Name-keyed ownership filters (generated lore relevance, world/scenario memories) now match stable
  IDs. Empty or foreign owners are unresolved and excluded from prompts; files are not rewritten.
- `SaveChatData` no longer writes `automated_world_lore.json` or `automated_scenario_lore.json`.
  Character generated lore keeps its existing write until F3-D, behind one allowlist function.
- Chat persistence no longer copies `WorldName`/`ScenarioName` into `WorldID`/`ScenarioID`. Names
  stay display-only; legacy chats with no canonical ID remain unresolved.
- `LoadChatData` reports per-source status through `GetContextStatus` so a consumer can display
  `not_set`, `not_found`, `unavailable`, `authority_indeterminate`, or `revision_mismatch` instead
  of inferring from nil data.
- The native test seam passes IDs and returns the context status. The app attaches the content
  library registry to the chat service after registry initialization.

## 2. Before/after inventory

Every world/scenario surface found in the audit is listed. “Item folder” means a canonical
`<root>/<id>/` directory; “chat-owned” means a file under `Chats/<chatId>/`.

### Reads

| Surface | Before | After | Class |
| --- | --- | --- | --- |
| World document | `Worlds/<worldName>/world.json` | `worldstore.Load(worldId)` | Converted |
| Scenario document | `Scenarios/<scenarioName>/scenario.json` | `worldstore.Load(scenarioId)` | Converted |
| Manual world lore | `Worlds/<worldName>/world_lore.json` | `LoadLore(worldId, "world_lore.json")` | Converted |
| Manual scenario lore | `Scenarios/<scenarioName>/scenario_lore.json` | `LoadLore(scenarioId, "scenario_lore.json")` | Converted |
| Generated world lore | `Worlds/<worldName>/automated_world_lore.json` | canonical item folder, classified | Converted |
| Generated scenario lore | `Scenarios/<scenarioName>/automated_scenario_lore.json` | canonical item folder, classified | Converted |
| Character document | `Characters/<name>/character.json` | unchanged | Retained |
| Character manual lore | `Characters/<name>/character_lore.json` | unchanged | Retained |
| Character generated lore | `Characters/<name>/automated_character_lore.json` | unchanged | Retained — F3-D |
| Automated lore relevance | `OwnerID == worldName/scenarioName` | `OwnerID == worldID/scenarioID`; empty owner stamped from item folder | Converted |
| World/scenario memories | name match or empty owner accepted | ID-only match; empty/name owners excluded | Converted |
| Chat snapshot document | `Chats/<id>/world.json`, `scenario.json` | unchanged; never refreshed on load | Retained — chat-owned |
| Chat meta identity | names copied into `WorldID`/`ScenarioID` | names stay display-only | Removed |
| RAG query labels | `WorldName`/`ScenarioName` | unchanged; search labels only, never resolution | Retained |
| Loaded context accessors | `GetWorldData`/`GetScenarioData` | unchanged; nil when unavailable | Retained |
| Prompt world/scenario blocks | from name-folder documents | only from authoritative canonical load | Converted |
| Dead path loaders | `loadWorldData`/`loadScenarioData` | deleted | Removed |

### Writes

| Surface | Before | After | Class |
| --- | --- | --- | --- |
| Generated world lore | written on `SaveChatData(sync=true)` | no write path | Removed |
| Generated scenario lore | written on `SaveChatData(sync=true)` | no write path | Removed |
| Generated character lore | written on `SaveChatData(sync=true)` | unchanged, one allowlist function | Retained — F3-D |
| Shared update helper | `updateAutomatedLoreFile` reachable for any scope | reachable only for character entries | Converted |
| Chat documents/history/memories | `chat.json`, `history.json`, memories | unchanged | Retained — chat-owned |
| Chat snapshot documents | `Chats/<id>/world.json`, `scenario.json` | unchanged | Retained — chat-owned |
| Chat meta backfill | wrote names into ID fields on load | removed | Removed |

### Explicitly deferred

Editor/repository wiring, import and rollback, lore promotion, snapshot revision fields and refresh,
cast sync, media manifests, the unavailable-root index, and all F3-D character generated-lore work.
None was started.

## 3. Authority behavior in the consumer

- `LoadChatData` receives IDs, not names, for resolution. An empty ID is `not_set`, not a lookup by
  name.
- A complete search that proves absence is `not_found`. An unavailable root, unreadable file, or
  decode failure is `unavailable`. Multiple copies or any unavailable root is
  `authority_indeterminate`.
- Indeterminate or unavailable content contributes **nothing** to `worldData`, `scenarioData`,
  `manualLore`, or `automatedLore`. No name-keyed fallback, recreation, import, or reference
  reassignment occurs.
- A stored revision mismatch is reported as `revision_mismatch`; the document remains readable
  because the file itself is present, but its revision is not trusted for writes.
- Missing lore files are empty, not missing items. A foreign `OwnerID` is excluded from the prompt
  pool while the file stays untouched; an empty `OwnerID` is stamped from the canonical folder that
  holds the entry.

## 4. Legacy compatibility

There are three distinct populations; see
`F3C-LEGACY-LORE-CLASSIFICATION-2026-11-10.md` for the full table.

1. **Legacy automated file already in a canonical item folder.** It is read through
   `worldstore.LoadLore`; `chat.Service.loadItemLore` classifies the entries in memory as
   `originKind: "legacy_generated_unattributed"`, `Automated: true`, and keeps them prompt-eligible.
   Empty owners are stamped from the canonical item ID; foreign owners are excluded. Nothing is
   written, migrated, or promoted. Tests assert the file bytes and folder contents are unchanged.
2. **Legacy automated file in an old name-keyed folder.** No production reader handles it. It is not
   read, not classified, and not prompt-eligible, and it is not a read-path fallback. LF-1 now
   authorizes a single explicit, journaled migration of these folders by verified identity or user
   resolution; unresolved content is quarantined. Migration code is not implemented, so the folders
   remain inert at runtime.
3. **New conversation discovery.** It stays chat-owned in memory storage and is not written into
   canonical world/scenario lore. Promotion is explicit and deferred.

No import was performed. The earlier statement that entries were “imported from `automated_*` files”
was too broad and has been corrected.

## 5. Verification

| Check | Result |
| --- | --- |
| `gofmt` on changed files | clean |
| `go vet ./internal/worldstore ./internal/backend/chat ./internal/backend/lorebook ./internal/library .` | exit 0 |
| `go test ./internal/worldstore` | `ok`; 38 tests |
| `go test ./internal/backend/chat` | `ok`; 23 tests |
| New canonical lore tests | 5 pass |
| New consumer tests | 6 pass |
| New persistence identity test | 1 pass |
| `go test -race ./internal/worldstore` | pass |
| A13 (`lorebook`, `chat`, `memory`, root) | all `ok`, exit 0 |
| Frontend identity + memory tests | exit 0, 4/4 each |
| Source audit | no `"Worlds"`/`"Scenarios"` paths remain in `internal/backend/chat` |
| Repository audit | `AUDIT-F3C-SINGLE-WRITER-2026-11-10.md`: dead name readers removed; export bypass closed |
| Write audit | only chat-owned writes and the character generated-lore allowlist remain |

## 6. Limits of this evidence

- `chat.Service` is not on the desktop inference path: `App` builds prompts from its own state and
  the chat snapshot. This conversion removes a real second writer and fixes the service's canonical
  reads, but it does not by itself change what the desktop prompt builder reads. That builder is a
  chat-owned snapshot consumer and is outside this slice.
- No canonical import exists yet, so a world/scenario that still lives only in the browser or in a
  legacy name folder has no canonical item. Legacy chats resolve to `not_set` until an explicit
  import assigns IDs. That is the intended fail-closed cutover consequence, not a silent fallback.
- Legacy name folders are authorized only as LF-1 migration input, not as a runtime source. Their
  automated lore is therefore unread and unclassified until that explicit migration exists. This
  consumer may not repair it at read time.
- Unit and filesystem tests only; no native run, no UI, no prompt-boundary capture.
- The store write lock is per process, as documented for the desktop runtime.

## 7. Gate status

A6 (no chat writes under item roots) has source and unit evidence. A10 has a source-level pass after
the repository audit: the store is the only canonical writer, the generic export bypass is refused,
and dead name readers are removed. A10's native runtime check also passes: `LoadChat`, `SaveChat`,
`SaveChatData`, canonical chat load, and generated-lore synchronization left canonical world and
scenario item trees byte-identical. Evidence:
`VERIFICATION-F3C-NATIVE-A9-A10-2026-11-10.md`. A11 is enforced by excluding
unavailable/indeterminate content from prompts. A13 passes. A9's native import evidence is recorded
in the same verification document. Snapshot refresh remains a later gate. A12 passed separately in
`VERIFICATION-F3C-NATIVE-A12-2026-11-10.md`.
