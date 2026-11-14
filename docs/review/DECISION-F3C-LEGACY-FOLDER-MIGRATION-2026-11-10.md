# DECISION — F3-C Legacy Name-Folder Migration Amendment (LF-1)

**Date:** 2026-11-10
**Status:** amendment recorded; import not implemented
**Amends:** `FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md` §6 (import), §8 A9 (import evidence)
**Related:** R2 (generated lore), §3 (single writer), §5 (update), A6/A9/A10/A13;
`F3C-LEGACY-LORE-CLASSIFICATION-2026-11-10.md`
**Authorization:** owner directed a formal amendment before import implementation; this record
authorizes migration design, not import code

## 1. Decision

Frozen §6 said legacy `Worlds/<name>` and `Scenarios/<name>` folders are never imported from and
never written to. R2 says legacy automated lore is imported once and remains prompt-eligible. Those
statements conflict for content that still lives only in a name folder.

**LF-1 selects option 2: one narrowly scoped, explicit, journaled, reversible migration of legacy
name-folder content into canonical items.**

The old folders remain:

- **never read at runtime**;
- **never a fallback** when a canonical item is missing or unavailable;
- **never a permanent import source**;
- **never written to or deleted**.

They are read only by the one-time migration operation defined here, and only after preflight and
explicit user authorization.

## 2. Amendment text

The final bullet of §6 is superseded as follows:

> Interrupted import is recoverable: marker absent plus identical files yields adoption and
> completion. Legacy `<name>` folders are never read at runtime and never written to. A single
> explicit, journaled legacy-folder migration authorized by `DECISION-F3C-LEGACY-FOLDER-MIGRATION-2026-11-10.md`
> may read them as migration input under that record's identity, conflict, quarantine, journal, and
> rollback rules. No other import path may read them. Mapping requires verified identity evidence or
> explicit user resolution; a display-name match alone is never sufficient. Unresolved content is
> quarantined, not discarded, merged, or guessed. Source folders remain untouched.

The first bullet of §6 is amended only to name the second input:

> Explicit page-level import of browser worlds and scenarios, plus the narrowly scoped legacy-folder
> migration authorized by LF-1; no automatic migration or cutover.

§8 A9 is amended to require:

> Preflight conflict report; verified identity mapping; unresolved content quarantined; idempotent
> retry; rollback removes only journaled creations and never modifies source folders.

No other frozen clause changes. In particular, A6, A10, A11, R1, R2, and R3 remain in force.

## 3. Scope

In scope:

- legacy `Worlds/<name>/world.json` and `Scenarios/<name>/scenario.json`;
- their `world_lore.json` / `scenario_lore.json`;
- their `automated_world_lore.json` / `automated_scenario_lore.json`;
- optional `media.json` references, preserved but not migrated to media roots;
- unknown fields in documents and lore entries.

Out of scope:

- runtime reads from name folders;
- automatic startup migration;
- title/fuzzy/recency matching;
- deleting, renaming, or rewriting source folders;
- image bytes or IndexedDB migration (F3-F);
- lore promotion from chats;
- character generated lore (F3-D);
- editor wiring and snapshot refresh.

## 4. Identity evidence

A legacy folder is mapped to a canonical or browser item only by one of:

1. **Exact stable identity:** the legacy document contains a non-empty `id` that exactly matches one
   candidate item's stable ID. The match must be unique. The mapping records the source ID, the
   destination ID, and both content hashes. A differing display name is recorded as a rename, not
   used as evidence.
2. **Explicit user resolution:** the user selects the destination for an otherwise unresolved
   folder. The selection, the evidence shown, and the actor are journaled.

Never sufficient on their own:

- matching folder name to item name;
- case folding, trimming, punctuation normalization, or fuzzy title comparison;
- matching lore content;
- position in a list, recency, or “only one candidate” heuristics;
- a chat snapshot that carries only a display name.

If the legacy document has no stable ID, or the ID matches zero or multiple candidates, the folder
is **unresolved**. If two folders claim the same identity, both are unresolved. If the same name
appears more than once, none is automatically mapped. If a browser record and a legacy folder share
a title but carry different IDs, they are different items and must not be merged.

A mapping decision is per folder. It does not establish that the folder and a same-named browser
record are the same creative object; only identity evidence or user resolution can do that.

## 5. Preflight

Preflight performs no writes. It must:

- enumerate every candidate name folder under each explicitly selected legacy input root; that root
  is not added to the content-library registry;
- read and validate the document envelope or raw record without normalizing away fields;
- record whether a stable `id` exists, its shape, and all matching candidates;
- record duplicate folders, duplicate identities, and duplicate names;
- inventory manual lore, automated lore, and `media.json` references;
- compute a canonical content hash for the document and each lore file;
- identify destination conflicts: existing item, differing content, existing lore, unmapped owners;
- identify unavailable or authority-indeterminate destination roots;
- produce a conflict/quarantine report before any write.

Preflight reads source files only. It does not create canonical folders, journals, markers, or
temporary files inside source folders.

## 6. Write, journal, and rollback

Per item, after a resolved mapping and a clean destination check:

1. stage the canonical item folder in the destination root;
2. write the document envelope;
3. write manual lore, if any;
4. write automated lore with `legacy_generated_unattributed` classification;
5. write `media.json` references, if any;
6. append the journal entry;
7. write or update the completion marker.

Journal entries must record:

- source root, folder name, and source file paths;
- source stable ID, if any, and the mapping evidence;
- destination item ID, root, and paths;
- content hashes for document and lore;
- created paths versus adopted existing paths;
- prior destination state for rollback;
- classification and owner decisions for lore;
- conflict/quarantine result, if the item did not commit.

Rollback:

- removes only canonical files and folders created by this migration;
- never deletes or rewrites a source folder;
- never deletes an adopted pre-existing canonical item;
- restores the marker and journal to the pre-operation state;
- reports any cleanup failure rather than claiming success.

The source folders remain byte-identical after a successful migration and after rollback.

## 7. Conflict, quarantine, and partial failure

No destination item is overwritten. If the destination item exists:

- identical document and no lore conflict may be recorded as **adopted**, with no writes;
- differing document content is `identity_conflict`; the item is not written;
- existing destination lore that would be replaced or silently merged is `lore_conflict`; the item
  is not written;
- a destination with an unavailable or indeterminate root is `destination_unavailable`; no write.

Quarantine is a first-class outcome, not deletion. At minimum these reasons are distinct:

- `ambiguous_identity`;
- `duplicate_name`;
- `duplicate_identity`;
- `missing_document_id`;
- `identity_conflict`;
- `lore_conflict`;
- `unmapped_owner`;
- `destination_unavailable`;
- `authority_indeterminate`;
- `write_failed`.

A quarantined folder keeps its source files and appears in the report/journal with its reason. It is
not prompt-eligible. It is not silently skipped.

A partial failure leaves completed items journaled and recoverable. No item is reported committed
until its document and lore writes are complete. An interrupted item is cleaned up by the same
rollback path or adopted on retry if the staged files are byte-identical.

## 8. Lore ownership and preservation

- Manual lore is preserved verbatim, including unknown fields and entry IDs.
- An empty manual owner may be stamped with the mapped item ID because the migration has established
  item ownership. A non-empty owner that does not resolve to the mapped item is `unmapped_owner`:
  the entry is preserved but remains unresolved and is excluded from prompts. It is not rewritten.
- Automated lore is imported with `Automated: true` and
  `originKind: "legacy_generated_unattributed"`. No chat ID or chat attribution is invented.
- Automated entries may receive the mapped item ID as owner because item ownership is established by
  the migration; this is not chat attribution.
- Manual and automated lore are never silently merged with each other or with existing destination
  lore.
- Entry IDs are preserved; a collision with an existing destination entry is `lore_conflict`, not a
  silent overwrite. Any future append/merge operation must assign new `lb_` IDs explicitly.

## 9. Idempotency and repeated execution

A repeated migration is a no-op when:

- the mapping is unchanged;
- the destination document and lore hashes match the journaled values;
- the source hashes match.

If source content changed since the journal entry, the result is `identity_conflict` or
`lore_conflict` as applicable; the migration does not overwrite the earlier canonical result. If the
destination changed independently, the same conflict rule applies. A completed marker does not
authorize a second automatic pass.

## 10. Runtime prohibition

After migration, production code still has no name-folder reader. The only canonical read path is
`worldstore` by stable ID. The migration package may read legacy folders only when invoked
explicitly; it is not reachable from chat load, prompt assembly, editor load, or startup.

The old folders are archival after the migration. No code may delete them, rename them, or treat them
as current content.

## 11. Acceptance criteria

**LF1 — Preflight.** Every candidate folder is inventoried with document/lore/media hashes, IDs,
duplicates, conflicts, and unavailable roots; preflight writes nothing.

**LF2 — Identity.** No folder is auto-mapped by name alone; automatic mapping requires one exact
stable-ID match; ambiguous, duplicate, and missing-ID cases are quarantined; explicit user mappings
are journaled with evidence.

**LF3 — Losslessness.** Unknown document and lore fields survive migration byte-for-byte at the JSON
value level; arrays and nested objects are not flattened or dropped.

**LF4 — Lore classification.** Automated lore is `legacy_generated_unattributed` with no chat ID;
manual lore preserves its owner; unresolved owners are preserved and excluded, not rewritten.

**LF5 — No overwrite.** Existing destination items are never silently replaced; identical content
may be adopted; differing document or lore content is a reported conflict with no write.

**LF6 — Authority.** Unavailable or authority-indeterminate destination roots stop the affected
item; no guess, no fallback, no partial success claim.

**LF7 — Journal and rollback.** Every mapping and write is journaled; rollback removes only created
canonical paths, leaves adopted items and all source folders untouched, and reports cleanup failure.

**LF8 — Recovery and idempotency.** Interrupted migration is recoverable; an identical rerun adopts
and completes; a changed source or destination conflicts without overwrite.

**LF9 — Quarantine.** Unresolved content is retained, reported with a distinct reason, and never
prompt-eligible or silently discarded.

**LF10 — Runtime boundary.** Source audit shows no runtime name-folder reader or writer outside the
explicit migration package; chat and prompt assembly remain stable-ID only.

**LF11 — Test fixtures.** Tests cover duplicate names, same title with different IDs, renamed folders,
deleted/recreated IDs, conflicting lore, interrupted writes, rollback, rerun idempotency, unavailable
roots, and unresolved owners.

**LF12 — Regression.** A9, A10 source audit, A13, frontend identity, and the existing chat/store
suites remain green.

## 12. Status and gates

- §6 amendment: **recorded**.
- Migration design: authorized by this record.
- Migration code: **not started**.
- Import code: **not started**.
- A9: open until LF1–LF12 are implemented and evidenced.
- A10: source audit passes; native runtime evidence remains open.
- Runtime name-based reads/writes: still prohibited.

This record does not make old folders a permanent source, a fallback, or an automatic migration.
It authorizes one explicit, auditable transition path so user content is not abandoned when the old
storage mechanism is retired.
