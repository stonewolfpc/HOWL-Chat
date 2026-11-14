# FROZEN CONTRACTS — F3-C Worlds and Scenarios

**Date:** 2026-11-10
**Status:** decisions resolved; contracts frozen; implementation not started
**Supersedes:** the open items in `REVIEW-F3C-DECISIONS-AND-COMPATIBILITY-2026-11-10.md` §8
**Change control:** any change to a frozen item requires a new dated decision record. Silent
adjustment during implementation is not permitted.

## 1. Resolved decisions

### R1 — Snapshot refresh

**Document refresh and scenario cast sync are separate, explicit operations.**

- Editing a world or scenario never mutates a chat.
- Opening a chat never refreshes a snapshot.
- Refresh replaces the embedded snapshot and its recorded revision. It keeps the chat's character
  list and lorebook selection.
- Cast sync is a second, explicit operation. It recomputes the chat's characters from the refreshed
  scenario, returns the diff (added, removed, changed), and requires confirmation. It appends a
  system message naming the removed characters, matching existing behavior.
- A missing source is `source_unavailable`: the snapshot is kept, refresh is refused, and nothing is
  substituted.

Rejected: combining both into one "refresh" action. It would silently change who is in a scene,
which is a creative decision the user must see.

### R2 — Generated lore ownership

**Canonical lore belongs to the item; generated lore belongs to the chat until explicitly promoted.**

- Manual world/scenario lore (`world_lore.json`, `scenario_lore.json`) is canonical item content,
  owned by the world/scenario ID, and moves with the item folder.
- Lore produced from a conversation is chat-owned and lives in that chat's
  `memories_snapshot.json`. The chat service stops writing `automated_world_lore.json` and
  `automated_scenario_lore.json` entirely.
- Promotion is explicit: a promoted entry is copied into canonical item lore with a new `lb_`
  identity and `origin {kind: "chat_promotion", chatId, promotedAt}`. The chat copy is retained.
- Legacy `automated_*_lore.json` files cannot be attributed to a chat: any chat could append to the
  same world file, and the file carries no chat ID. They are imported once as canonical
  `origin.kind = "legacy_generated_unattributed"`, `Automated: true`, flagged in the UI, and remain
  prompt-eligible for that item so existing worlds keep working. No attribution is invented.
- Consequence, stated plainly: newly generated lore no longer becomes shared world knowledge
  automatically. Promotion replaces that implicit behavior with a visible action.

Rejected: auto-promoting generated entries to world canon (keeps the cross-chat leak), and assigning
legacy entries to the most recent chat (invents provenance the file does not contain).

### R3 — Binary assets

**Item folders contain the document and references only. No binaries.**

- The item folder holds `world.json` or `scenario.json` plus an optional `media.json`.
- `media.json` is `{format: "howl-media-refs-v1", items: [{mediaId, role, provider, legacyRef?}]}`
  with roles `thumbnail`, `gallery`, `attachment`.
- Image bytes stay under F3-F media roots; F3-C never copies, moves, or writes image bytes.
- Legacy IndexedDB references (`{id, name, type, size, stored: true}`) are preserved verbatim in the
  document and carried as `provider: "legacy-indexeddb"` until F3-F migrates them.
- `thumbnailIndex` remains a display concern in the document.

Rejected: `media/` subfolders inside the item. They would create a second media owner and bypass the
F3-F atomic commit and gallery identity.

## 2. Frozen document contract

### 2.1 Envelopes

```json
{ "format": "howl-world-v1",    "id": "...", "revision": "rv1-...", "created_at": "...",
  "updated_at": "...", "origin": {"kind": "howl|imported|external", "legacy_ids": [], "legacy_name": ""},
  "document": { "...": "verbatim world record" } }
```

`howl-scenario-v1` is identical with a scenario document. `document` is the existing record shape and
is preserved verbatim, including unknown fields.

### 2.2 Identity

- Identity is opaque: `[A-Za-z0-9_-]{1,64}`, unique within a category, path-safe, never parsed.
- Existing unique IDs are preserved verbatim, including bare timestamps and both scenario shapes.
- New items get `wld_<16 hex>` or `scn_<16 hex>`.
- Unsafe or duplicate IDs are reassigned through a journaled `legacy_id → new_id` map. References
  are rewritten only through that map.
- An ID is never derived from a title, folder name, or display name.

### 2.3 Paths

```text
<root>/<id>/world.json
<root>/<id>/scenario.json
<root>/<id>/media.json          # optional, references only
```

Renaming an item changes no path. Moving a root changes no identity.

### 2.4 Content and revision

- `revision` is computed by the store, not supplied by a client:
  `rv1-` + first 16 hex of SHA-256 over canonical JSON (sorted keys, no whitespace) of `document`,
  excluding `updated_at` and media path fields.
- The revision changes when canonical content changes, not when metadata timestamps do.
- A stored revision that does not match recomputation is reported as `revision_mismatch`; the item is
  readable but flagged.

## 3. Frozen ownership contract

- **Single writer:** only the world/scenario store writes item folders. `chat.Service` becomes a
  consumer; it must not write anywhere under `Worlds/` or `Scenarios/`.
- **Read path:** one repository, one canonical read path, one rebuildable cache. A canonical
  consumer must not fall back to localStorage when the store is unavailable; it must display the
  failure.
- **Roles:** create requires a save root. Replacing requires the write role on the root that holds
  the item. Read-only roots can be browsed, attached, and referenced; never written.
- **Lore files:** `world_lore.json` / `scenario_lore.json` move with the item and are written only by
  the store. `OwnerID` on world/scenario-scope entries is the stable item ID. Unmapped legacy owners
  stay unresolved and are excluded from prompts.
- **Availability:** a disconnected root keeps its registration and last-known index entries, marked
  unavailable, excluded from prompts, and never treated as deletion.

## 4. Frozen reference contract

```json
{ "kind": "world|scenario|character|lorebook", "id": "...", "title": "optional display cache" }
```

- Resolution is by `kind` + `id` only. Display names never authorize or resolve a reference.
- Location: world `scenarios`, `characters`, `howl.lorebooks`; scenario `characters`, `lorebooks`,
  `howl.world`; chat `active/removed/added*` arrays keep IDs.
- Legacy name-only references become `{kind, id: "", legacyName: "..."}` during import. They are
  displayed as unresolved and can be repaired explicitly. They never match by title at read time.
- Canonical writes use `howl.lorebooks` and `lorebooks` as one normalized list; the repository
  normalizes reads. The top-level duplicate is not written going forward.
- A dangling reference is preserved, displayed unavailable, excluded from prompts, and repairable.

## 5. Frozen update contract

- `Save` is create-only. `Update` carries `baseRevision`; a mismatch returns `revision_conflict`
  with the current revision and requires an explicit choice.
- Update merges: incoming known fields are applied over the canonical document; unknown top-level
  keys and unknown keys inside `howl` and `metadata` are preserved; arrays are replaced wholesale
  when present; explicit `null` clears only allowlisted fields.
- `id`, `created_at`, `origin.legacy_ids`, and `revision` are immutable to clients.
- Writes are atomic: temp file in the item folder, flush, rename over the target. A failed write
  leaves the prior document intact.
- Failure states are distinct and user-facing: `destination_unavailable`, `collision`,
  `permission_denied`, `unsafe_target`, `invalid_document`, `invalid_reference`,
  `revision_conflict`, `duplicate_identity`, `partial_import`, `write_failed`,
  `legacy_unmapped_owner`. A cache write is never reported as a canonical save.

## 6. Frozen import and rollback contract

- Explicit page-level import of browser worlds and scenarios, plus the narrowly scoped legacy-folder
  migration authorized by LF-1; no automatic migration or cutover.
- Preflight: enumerate records, validate identity shape and uniqueness, normalize references,
  preserve unknown fields, inventory lore and image references, report conflicts before writing.
- Write order per item: stage folder → document → lore → `media.json` → journal → marker.
- Journal: `<config>/worlds-import-journal.json` and `scenarios-import-journal.json`; every created
  or adopted item recorded. Rollback removes only journaled creations and clears the marker;
  browser data is never modified.
- Idempotent by identity plus content hash; identical files are adopted and journaled as adopted.
- Interrupted import is recoverable: marker absent plus identical files yields adoption and
  completion. ~~Legacy `Worlds/<name>` and `Scenarios/<name>` folders are never imported from and never
  written to.~~ **Superseded by LF-1 (2026-11-10):** legacy name folders are never read at runtime and
  never written to. One explicit, journaled migration may read them as migration input under
  `DECISION-F3C-LEGACY-FOLDER-MIGRATION-2026-11-10.md`; no other import path may. Mapping requires
  verified identity evidence or explicit user resolution; name matching alone is insufficient.
  Unresolved content is quarantined, not discarded or guessed. Source folders remain untouched.

**Amendment LF-1 — 2026-11-10.** The change above is a controlled amendment, not a silent edit. It
selects the explicit journaled migration option in
`F3C-LEGACY-LORE-CLASSIFICATION-2026-11-10.md` §4. All other §6 rules and all other frozen clauses
remain in force.

## 7. Frozen snapshot contract

- Chat meta gains `worldRevision`, `scenarioRevision`, `worldSnapshotAt`, `scenarioSnapshotAt`.
  `worldId`/`scenarioId` carry identity; names stay display-only.
- Snapshot files remain `world.json` / `scenario.json` inside the chat folder.
- `GetChatSnapshotStatus(chatId)` returns per source: id, snapshot revision, current revision, stale,
  source unavailable.
- `RefreshChatWorld` / `RefreshChatScenario` replace the snapshot and revision only.
- `SyncChatScenarioCast` returns a diff, requires confirmation, and appends the existing removal
  system message. It is never implicit.
- A deleted or disconnected source keeps the snapshot; refresh is refused with `source_unavailable`.

## 8. Acceptance matrix

| ID | Contract | Evidence required |
| --- | --- | --- |
| A1 | Identity and paths | Two worlds and two scenarios; rename keeps ID and path; restart read-back |
| A2 | Lossless update | Regression: unknown top-level and nested fields survive an editor save (C6) |
| A3 | Roles | Create-only; replace requires write role; read-only root refuses |
| A4 | Revision | Content change changes revision; timestamp-only change does not; conflict refuses silently losing a write |
| A5 | References | ID-only resolution; same-title recreation does not resolve (C4); dangling refs unavailable and repairable |
| A6 | Generated lore | No writes under `Worlds/`/`Scenarios/` by chat; promotion is explicit; legacy generated lore flagged and preserved |
| A7 | Snapshots | No implicit mutation on load or edit; document refresh and cast sync separated; staleness visible |
| A8 | Media | No binaries in item folders; legacy IndexedDB refs preserved; `media.json` reference-only |
| A9 | Import | Preflight conflict report; verified identity mapping; unresolved content quarantined; idempotent retry; rollback removes only journaled creations and never modifies source folders (LF-1) |
| A10 | Single writer | Source audit plus runtime check that chat cannot write item folders |
| A11 | Availability | Disconnected root keeps items unavailable, excluded from prompts |
| A12 | Native evidence | Isolated native run; prompt-boundary capture; production binary free of test seams |
| A13 | Regression | Lorebook identity/privacy suites, frontend identity tests, conflict scope-boundary test |

## 9. Consequences to keep visible

- Newly generated lore stops being shared world knowledge automatically; promotion is the replacement.
- Existing `automated_world_lore.json` content becomes canonical but visibly unattributed, not silently
  merged with user-authored lore.
- Character generated lore (`automated_character_lore.json`) keeps its current behavior until F3-D,
  where the same second-writer pattern is removed. F3-C must not extend it.
- The chat service's name-keyed reads and writes are replaced in this slice; leaving either would
  break the single-writer contract.

## 10. Status

Contracts frozen. No code changed. Implementation begins with the store skeleton, the lossless
update path, and the reference normalizer, in that order.
