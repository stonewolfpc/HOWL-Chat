# REVIEW — F3-C Decisions and Compatibility

**Date:** 2026-11-10
**Scope:** validate the six owner recommendations against the actual world/scenario records, every
reference consumer, chat snapshot behavior, and the F3-F media contract. No implementation.
**Method:** source reading of the current tree. No runtime test, build, or migration was performed.

## 1. Outcome

All six recommendations are accepted. One qualification applies to Q2. The compatibility pass
produced nine findings — seven conflicts (C1, C2, C3, C4, C6, C7, C9) and two gaps (C5, C8). None
invalidates a decision, but three change what "done" means: the chat service is a second writer
keyed by name (C3), the world/scenario editor already drops unknown fields on save (C6), and
folder-per-item assets conflict with F3-F media ownership (C7).

## 2. Decision record

| ID | Decision | Binding detail |
| --- | --- | --- |
| F3-C-Q1 | **Folder per item** | `<root>/<id>/world.json` and `<root>/<id>/scenario.json`; assets are referenced, not required to live in the folder (C7) |
| F3-C-Q2 | **Stable identity, preserve usable legacy IDs** | Existing unique IDs are kept verbatim; new items get `wld_`/`scn_`; identity is opaque and never parsed; duplicates are reassigned with a journaled mapping; every mapping is recorded |
| F3-C-Q3 | **Explicit journaled folder import** | Create-only, idempotent by identity plus content hash, journaled, reversible; browser data untouched |
| F3-C-Q4 | **Pinned snapshots with visible staleness** | Chat records canonical revision IDs; snapshot persists until explicit refresh; staleness is visible; refresh is an explicit, explained operation |
| F3-C-Q5 | **Preserve and display unavailable** | Missing references keep their ID, display as unavailable, are excluded from prompts, and are repairable; never substituted |
| F3-C-Q6 | **Stable IDs for world/scenario lore now** | World/scenario lore resolves by world/scenario ID; renames do not change ownership; deleted owners stay unresolved; same-name recreation inherits nothing |

## 3. Compatibility findings

### C1 — World and scenario IDs are not one format

- World IDs are bare millisecond strings: `Date.now().toString()` in `world-form.js:200`,
  `world-grid.js:104`, `world-manager.js:348`.
- Scenario IDs have **two shapes** in the same build: `scenario-<ms>-<random>` in
  `scenario-form.js:163` and `scenario-form.js:166`, and bare `Date.now().toString()` in
  `scenario-manager.js:345`.
- Same-millisecond creation can collide, and world/scenario/lorebook IDs share the timestamp space.

**Consequence:** Q2's "preserve where possible" produces a mixed namespace with no common prefix.
That is acceptable if identity is opaque, but no code may branch on the shape.
**Safe migration:** preserve unique IDs verbatim; mint prefixed IDs only for new records; detect
within-category duplicates during preflight and reassign with a journaled `legacy_id → new_id` map;
never rewrite references by parsing an ID.

### C2 — Chat metadata can hold a name in an "ID" field

`internal/backend/chat/persistence.go:432-436` backfills `meta.WorldID = meta.WorldName` and
`meta.ScenarioID = meta.ScenarioName` whenever the ID is empty. Legacy chats therefore have
display names in fields named `worldId`/`scenarioId`.

**Consequence:** a migration that trusts those fields as IDs would mint identity from a name — the
defect W27 removed. **Safe migration:** treat a value as an ID only when it matches a candidate
record's real ID; otherwise resolve through the explicit name→ID mapping and mark unmapped values
unresolved. Never auto-create or auto-relink from a name.

### C3 — The chat service is a second writer, keyed by name

`internal/backend/chat/service.go:91-99` loads `Worlds/<name>/world.json` and
`Scenarios/<name>/scenario.json` from the state root. Lines 118-120 and 149-150 load
`world_lore.json`, `scenario_lore.json`, `automated_world_lore.json`,
`automated_scenario_lore.json` from the same name folders. Lines 412-417 **write**
`automated_world_lore.json` / `automated_scenario_lore.json` under those name paths.

**Consequence:** after a cutover to `<root>/<id>/`, chat reads fail, and the chat service keeps
creating name-keyed lore files outside the canonical store. That breaks the single-writer rule Q3
and Q6 rely on. **Safe migration:** resolve world/scenario documents and lore through the store by
ID; move generated lore writes into the canonical item folder (or explicitly declare them
chat-scoped generated data, which is a separate decision); keep a read-only legacy fallback for the
existing name folders during the transition, never as a write target.

### C4 — Reference matching falls back to display name

Both resolvers match ID first and then name/title: `context-data.js:44-56` (`refMatches`,
`resolveLinkedItems`) and the duplicate `sidebar-state.js:130-141`. `context-selectors.js:195-235`
uses them to attach scenarios, characters, and lorebooks from references.

**Consequence:** a recreated world or scenario with the same title will satisfy a stale reference —
the silent substitution Q5 forbids. **Safe migration:** typed references; resolve by ID only after
migration; name-only legacy references are marked `legacyName` and displayed unresolved unless an
explicit mapping exists. Do not keep the name fallback in the canonical path.

### C5 — Snapshots are already pinned; the revision is what is missing

`ChatData` embeds `World`/`Scenario` (`internal/backend/chat/types.go:78-86`), and
`persistence.go:215-235` writes `world.json`/`scenario.json` into the chat folder. Loading reads
those copies. So pinning exists; there is no revision field, no staleness indication, and no refresh
operation.

**Consequence:** Q4 is additive, but refresh is not a pure document swap: `scenario-form.js` derives
characters from the scenario, and `chat-state-bridge.js` emits a system message when a scenario
change removes characters. **Safe migration:** store `worldRevision`/`scenarioRevision` plus the
source IDs in chat metadata; define refresh as "replace snapshot, recompute characters/lorebooks,
report the diff"; never mutate the snapshot implicitly.

### C6 — The editor is lossy today

`world-form.js:199-229` rebuilds the world document from form fields (`id`, `title`, `details`,
`images`, `thumbnailIndex`, `scenarios`, `characters`, `createdAt`, `updatedAt`, `howl`).
`scenario-form.js:130-166` does the same and preserves only `metadata.created` when editing.

**Consequence:** migration can preserve every field and the next ordinary save will delete whatever
the editor does not know about. Requirement 3 (nothing disappears) fails at the editor, not the
store. **Safe migration:** the repository save must merge unknown top-level and `howl`/`metadata`
keys from the existing document; the editor sends its known fields, the store retains the rest.

### C7 — Folder-per-item conflicts with F3-F media ownership

World/scenario images are IndexedDB blobs (`storage-indexeddb.js:28-30`: `world-images`,
`scenario-images`), stripped from localStorage by `_stripWorldBlobs`/`_stripScenarioBlobs` and
restored with `restoreWorldImages`/`restoreScenarioImages`. F2-C §5.6 defines generated media as a
paired `<media-root>/<id>.png` + `<id>.json` in a registered Images root with atomic commit.

**Consequence:** copying blobs into `<world>/media/` would create a second media owner, duplicate
identity, and bypass the F3-F commit rules. **Safe migration:** the item folder holds the document
plus a `media.json` manifest of typed references (`mediaId`, role, legacy provider); binaries stay
in the F3-F media roots. Until F3-F, references keep `legacyProvider: "indexeddb"` and the existing
rehydration path. Q1's "room for assets" is satisfied by references and a stable folder boundary,
not by relocating binaries.

### C8 — No native world/scenario writer exists to inherit

`storeWorld` (`storage-scenarios.js:330-391`) writes localStorage and IndexedDB only; there is no
frontend call to a native world save. `ExportContentItem` (`app-libraries.go`) can write
`howl-content-export-v1`, is explicitly `authoritative: false`, and has no frontend caller. The
`Worlds/<name>/world.json` files, where present, are legacy export copies with no read path.

**Consequence:** there is no canonical file source to adopt; import must read browser records. The
old name folders must stay untouched and read-only, or importing them would reinstate name identity.

### C9 — World/scenario lore shares the item's folder and name

Manual lore lives beside the document (`world_lore.json`, `scenario_lore.json`) and automated lore
carries `OwnerID` equal to the display name; `service.go:210-212` filters on that equality.
**Safe migration:** move the lore documents with the item, rewrite `OwnerID` to the stable ID with
a journaled mapping, and stamp generated entries with the ID. Unmapped legacy owners stay unresolved
and are excluded from prompts, mirroring W28.

## 4. Reference consumer inventory

| Consumer | Reads | Impact of an ID cutover |
| --- | --- | --- |
| `world-grid.js`, `world-form.js`, `world-manager.js`, `world-state.js` | `howl-worlds` | Must read the repository; IDs preserved or mapped |
| `scenario-grid.js`, `scenario-form.js`, `scenario-manager.js`, `scenario-state.js`, `storage-scenarios.js` | `howl-chat-scenarios` | Same, plus two legacy ID shapes (C1) |
| `context-data.js`, `context-selectors.js`, `sidebar-state.js` | both catalogs; name fallback | Name fallback must be removed from the canonical path (C4) |
| `chat-state-bridge.js` | both catalogs; builds chat meta and embedded snapshot | Must carry revision and source IDs (C5) |
| `data-sync.js` | embedded snapshot → meta IDs | Meta follows the snapshot; must not overwrite revisions |
| `character-manager.js`, `character-form.js` | world/scenario records for attachment counts | Attribute access unchanged; reference shape must stay readable |
| `delete-service.js` | both catalogs; detaches by ID | Detach must target the store, not localStorage |
| `lorebook-manager.js`, `world-dropdowns.js`, `scenario-dropdowns.js` | attachment refs | Typed reference reader required |
| `sidebar-context-world.js`, `sidebar-context-scenario.js`, `sidebar-panels.js`, `sidebar-session.js` | records, then `AddWorldToChat`/`AddScenarioToChat` | Backend receives full objects; must receive IDs and resolve canonically |
| `app-wipe.js` | localStorage key list | Registry-aware wipe is F3-H; F3-C must not add a second list |
| `internal/backend/chat/service.go` | name-keyed files | Must resolve through the store (C3) |
| `internal/backend/chat/persistence.go` | embedded snapshots, name→ID backfill | Backfill must stop (C2) |

## 5. F3-F media validation

| Requirement | Current state | Verdict |
| --- | --- | --- |
| Stable media identity | `img_<ms>_<random>` in IndexedDB; gallery uses its own IDs | F3-C references must use F3-F identity, not image IDs invented by the form |
| Atomic binary + sidecar commit | F3-F contract; not implemented for world/scenario images | Out of F3-C scope; reference-only until F3-F |
| No second media store | IndexedDB today; folder-per-item could create one | Conflict C7; resolved by `media.json` references |
| Migration of existing images | Blobs stripped from localStorage, live in IndexedDB | F3-F owns extraction; F3-C must not lose the `{id, stored:true}` reference |

## 6. Additional acceptance requirements

1. **Typed cross-category references.** One schema (`{kind, id, title?, legacyName?}`) replaces
   untyped `[]interface{}` and duplicated `lorebooks`/`howl.lorebooks` arrays. Validation on save;
   unresolved references are visible and repairable.
2. **Distinguishable failure states.** Unavailable root, duplicate identity, incomplete import,
   invalid reference, permission denial, collision, and failed write each have a distinct user-facing
   state. A cache write is never presented as a canonical save.
3. **Lossless preservation.** The repository merges unknown fields on save (C6); migration accounts
   for document fields, `howl`/`metadata` keys, lore files, and image references; nothing disappears
   because the new store does not recognize a legacy field.

## 7. Revised implementation order

1. Store skeleton with identity, roles, atomic write, revision hash, and the lossless merge (C6).
2. Reference schema and normalizer, with ID-only resolution in the canonical path (C4).
3. Explicit import with preflight (C1, C2) and journaled rollback (C8 legacy folders untouched).
4. Chat service cutover to store-by-ID reads plus relocation of generated lore writes (C3, C9).
5. Snapshot revision recording, staleness indicator, explicit refresh (C5).
6. Media reference manifest, leaving binaries to F3-F (C7).
7. Native isolated acceptance with prompt-boundary capture and regression suites.

## 8. Not decided here

- Refresh semantics detail: whether refresh also recreates the chat's character list or only reports
  the diff. This needs a UI decision before step 5.
- Whether generated world/scenario lore is canonical item content or chat-scoped generated data.
  Recommendation: canonical, stamped with the stable ID, because the chat currently writes it beside
  the world and treats it as world lore.
- Whether any assets may ever live inside the item folder. Current recommendation: no binaries until
  F3-F amends the media contract.

## 9. Status

Decisions recorded; implementation not started. Two conflicts (C3, C6) are prerequisites, not
polish: without them the store can be correct and the product still lose data or keep writing
outside ownership.
