# CLARIFICATION — F3-C Legacy Generated-Lore Categories

**Date:** 2026-11-10
**Status:** clarification recorded; no import or migration implemented
**Related:** `FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md` §R2, §3, §6
**Supersedes:** the ambiguous “imported from `automated_*` files” shorthand in
`IMPLEMENTATION-F3C-CHAT-CONSUMER-2026-11-10.md` and `PROJECT_ISSUES_AND_EXPANSION.md`

## 1. Why this record exists

R2 says legacy `automated_*_lore.json` files are imported once as canonical,
`legacy_generated_unattributed`, `Automated: true`, and remain prompt-eligible. The chat consumer was
implemented before import. That makes “legacy generated lore remains prompt-eligible” true only for
files that are already inside a canonical item folder. It is not true for files in old name-keyed
folders, and it is not true for newly generated chat discoveries.

These are three different populations and must not be described with one sentence.

## 2. The three categories

| # | Population | Where it lives | Current reader | Classification | Disk effect of reading |
| --- | --- | --- | --- | --- | --- |
| 1 | Legacy automated file already in a canonical item folder | `<root>/<id>/automated_world_lore.json` or `automated_scenario_lore.json` | `worldstore.LoadLore` → `chat.Service.loadItemLore` | In memory: `Automated=true`, `OriginKind="legacy_generated_unattributed"`, `Enabled=true`; empty owner stamped from the canonical item ID; foreign owner excluded from prompts | Read-only |
| 2 | Legacy automated file in an old name-keyed folder | `Worlds/<name>/automated_world_lore.json` or `Scenarios/<name>/automated_scenario_lore.json`, outside a canonical item | **None** | **Not classified, not read, not prompt-eligible** | None |
| 3 | Newly generated conversation discovery | Chat memory state and `Chats/<chatId>/memories_snapshot.json` | Chat/memory service | Chat-owned; no canonical origin assigned | Chat-owned writes only; no item-folder write |

## 3. Category 1 — canonical legacy file

`worldstore.LoadLore` only reads the allowlisted filename from `<root>/<id>/`. It:

- resolves the item by stable ID and applies the same completeness/authority rules as `Load`;
- returns raw entries and the item facts;
- returns no entries when authority is indeterminate or the item is unavailable;
- treats a missing lore file as empty, not as a missing item;
- never writes, creates, migrates, or rewrites anything.

`chat.Service.loadItemLore` performs the classification in memory after the read:

- sets `Automated=true`;
- sets `OriginKind=lorebook.OriginLegacyGeneratedUnattributed`;
- sets `Enabled=true`;
- fills an empty `OwnerID` with the canonical item ID and excludes an entry whose non-empty
  `OwnerID` differs from that ID;
- returns the entries to the service's automated prompt pool.

No chat ID is invented. No entry is promoted. No attribution is written to disk. The file's bytes
and folder contents are unchanged; `TestLoadLoreDoesNotModifyTheItemFolder` and
`TestLoadCanonicalWorldUsesStableIDAndClassifiesLegacyLore` assert this.

This is the only population for which “legacy generated lore stays prompt-eligible” is currently
true.

## 4. Category 2 — old name-keyed folder

No production reader handles this population. `chat.Service` has no `Worlds/` or `Scenarios/` path,
no name-folder loader, and no fallback. `worldstore.LoadLore` cannot address a folder by display
name; its only address is a stable ID inside a registered root. If a name folder happens to sit in a
registered root, requesting its name as an ID either proves absence or fails decode; it is never
silently accepted as a canonical item.

The original frozen contract said:

> Legacy `Worlds/<name>` and `Scenarios/<name>` folders are never imported from and never written
> to.

**Resolved by LF-1 on 2026-11-10:** the owner selected option 2, a one-time explicit, journaled,
reversible migration. See `DECISION-F3C-LEGACY-FOLDER-MIGRATION-2026-11-10.md`. The three options
below are retained as the decision history.

1. **Browser records only.** Import only page-level browser worlds/scenarios; leave old name folders
   untouched and unread.
2. **Amend §6 for a one-time explicit folder migration.** Add a separately authorized, journaled
   import of old name folders with preflight, duplicate handling, and no overwrite. **Selected.**
3. **Leave old folders permanently inert.** Document that they are archival only and provide a
   manual repair path outside import.

Category 2 remains **inert at runtime**. It is not a fallback and not a permanent import source. It
is read only by the explicit migration, only under verified identity evidence or user resolution,
and only with quarantine for unresolved content. Migration code has not been written.

## 5. Category 3 — new conversation discoveries

Generated conversation entries stay in chat-owned memory storage. `SaveChatData` can write generated
lore only for character scope; `automatedLoreSyncTarget` has no world or scenario branch. There is no
implicit promotion into `world_lore.json` or `scenario_lore.json`.

Promotion is a separate future operation: a promoted entry must receive a new `lb_` identity and an
explicit `origin {kind: "chat_promotion", chatId, promotedAt}` record, and the chat copy is retained.
Until then, a conversation discovery is not canonical world/scenario knowledge and must not be
described as such.

## 6. Corrected wording

The earlier sentence “generated world/scenario entries imported from `automated_*` files carry
`originKind` …” implied that an import had run. It did not. The accurate statement is:

> When a legacy automated file is already present in a canonical item folder, its entries are read
> as `legacy_generated_unattributed` and remain prompt-eligible. No import has run. Files in old
> name-keyed folders are not read. New conversation discoveries remain chat-owned.

## 7. Consequences for the import gate

Import design must state its source population explicitly:
- browser page records (explicit, page-level);
- legacy name folders only through the LF-1 migration, by verified identity or user resolution;
- never automatic migration, never a read-time fallback, never an ownership guess.

LF-1 is a contract amendment, not an implementation convenience. Category 2 must not be silently
“fixed” during import, matched by title, or merged into a same-named browser record. Unresolved
folders are quarantined and retained. The migration code and its acceptance tests are still future
work.
