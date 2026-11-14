# DESIGN — F3-C Explicit Import and LF-1 Legacy Migration

**Date:** 2026-11-10
**Status:** implemented 2026-11-10; native A9/A10 evidence outstanding
**Scope:** explicit page-level import of browser worlds/scenarios, plus the LF-1 one-time legacy
name-folder migration
**Contracts:** `FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md` §3–§6; A9, A10, A13;
`DECISION-F3C-LEGACY-FOLDER-MIGRATION-2026-11-10.md` (LF-1);
`AUDIT-F3C-SINGLE-WRITER-2026-11-10.md`

## 1. Goal and non-goals

Goal: move existing user content into canonical items without losing fields, inventing ownership,
silently overwriting canonical work, or creating a runtime name-based path.

Non-goals in this slice:

- editor/repository wiring;
- automatic startup migration or cutover;
- snapshot refresh or cast sync;
- image-byte or IndexedDB migration (F3-F);
- lore promotion from chats;
- character generated-lore changes (F3-D);
- deletion or renaming of source folders;
- a permanent legacy-folder reader.

## 2. Inputs and authority

There are exactly two explicit sources:

| Source | Supplied by | Authority |
| --- | --- | --- |
| Browser records | Renderer sends a deep copy of selected `howl-worlds` / `howl-chat-scenarios` records | Legacy browser content; import input only |
| Legacy name folders | User explicitly selects a root; Go reads it under LF-1 | One-time migration input only |

The Go import operation never reads `localStorage`. The renderer reads it to build the request; that
read is observational. The renderer must not write or normalize records during preview or execute.

A legacy root is **not** added to the content-library registry. It is an ephemeral source selection
for the operation. It must be outside every registered content root, including equivalent, parent,
and child paths. A source that overlaps the destination produces `source_destination_overlap`; a
source that is or contains another registered root produces `source_registered`. Neither case
writes anything.

No import runs at startup, on chat load, on editor open, or as a side effect of another operation.

## 3. Operation state machine

```
idle
  → previewing          (pure read; returns plan; writes nothing)
  → authorized          (user reviewed plan; plan hash confirmed)
  → executing           (per-item stage/commit; journal active)
  → completed           (all ready items committed/adopted; quarantines reported)
  → completed_partial   (some quarantined/conflicted; no pending unsafe work)
  → rolled_back         (journal-created items removed; adopted/source untouched)
```

`previewing` is not persisted. A plan is returned to the renderer and resubmitted for execution; the
app does not create a plan file, directory, journal, marker, or temp file during preview.

Execution is refused unless the plan hash matches the current source fingerprints and destination
authority. Any source or destination change after preview returns `plan_stale` and requires a new
preview.

The renderer-supplied plan is not an authority. `ExecuteContentImport` receives the original request
plus the plan hash, re-runs preflight from the sources and destination, and compares the re-derived
plan hash. A plan whose write set, mappings, or dispositions cannot be reproduced is rejected. The
UI cannot add, remove, or redirect writes by editing the plan object.

## 4. Preflight — observational only

Preflight must not create directories, normalize source files, rewrite browser records, create
journals/markers/staging/temp files, write app logs containing source content, or partially migrate
anything. It may only read, validate, hash, and report. Diagnostics are returned in memory with the
plan.

For browser records:

- validate the record is an object and within the document size limit;
- record `id`, display name, timestamps, references, images, and unknown fields;
- compute a raw-record hash and a canonical-document hash;
- detect inline binary/base64 media and report `inline_media_unmigrated` rather than dropping it;
- detect duplicate IDs and duplicate names in the supplied set.

For legacy folders:

- enumerate only direct child folders of the explicitly selected root;
- require exactly one `world.json` or `scenario.json` matching the category;
- read with no symlink/reparse traversal; reject symlinked folders and files;
- inventory `world_lore.json` / `scenario_lore.json` and the automated counterparts;
- inventory optional `media.json`; report other files as `unsupported_source_file`, not import them;
- compute raw hashes for every source file and a canonical hash for the document;
- record whether the document is a raw record or an existing canonical envelope;
- record missing IDs, invalid IDs, duplicates, and destination collisions.

For every proposed item, preflight checks:

- destination registration exists, is a save root, is writable for rollback, and is available;
- `worldstore.Lookup` across all roots is complete and authoritative; any unavailable root or
  duplicate copy is `authority_indeterminate` / `destination_unavailable`;
- an existing item with the same ID: identical document, lore, and media may be adopted; differing
  content is `identity_conflict`, `lore_conflict`, or `media_conflict`; no write is planned;
- destination lore conflicts, entry-ID collisions, and unmapped owners;
- whether the source/folder mapping has verified identity evidence or needs user resolution.

The plan reports counts for `ready`, `adoptable`, `conflicted`, `quarantined`, `unsupported`, and
`writes`, plus every item's reason. The renderer must show this before enabling execution.

## 5. Identity resolution

Canonical identity is opaque and validated by the store. Import never derives an ID from a title,
folder name, timestamp, or content.

Accepted evidence:

1. **Exact stable ID:** the source record/envelope has one non-empty ID that matches exactly one
   candidate destination item.
2. **Explicit user resolution:** the user maps an unresolved source to a destination. The mapping,
   evidence shown, and actor are journaled.

Not accepted:

- name/title match, case folding, trimming, punctuation normalization, fuzzy matching;
- content similarity alone;
- recency, list position, or “only one candidate”;
- a chat snapshot carrying only a display name.

Resolution rules:

- missing source ID → `missing_document_id` quarantine;
- invalid ID shape → propose a store-generated ID; requires explicit approval; journal
  `legacy_id → new_id`;
- duplicate source IDs → all conflicting sources quarantined;
- ID matches zero destinations → new item with that ID, if shape is valid;
- ID matches one destination with identical document → adopt;
- ID matches one destination with different content → `identity_conflict`; options are skip or
  explicit new-ID import, never overwrite;
- ID matches multiple destinations → `ambiguous_identity` quarantine.

When a new ID is assigned, only references whose exact source ID is in the journaled mapping are
rewritten. Name-only references remain `legacyName` and unresolved.

A folder name and a same-named browser record are not assumed to be the same item. They need exact
ID evidence or explicit user resolution.

## 6. Canonical conversion

The canonical envelope is produced by the store, not by the renderer.

- Envelope: correct `format`, `version`, store-computed `revision`, `created_at`, `updated_at`, and
  `origin.kind = "imported"`.
- `origin.legacy_ids` records source IDs that differ from the destination ID.
- `origin.legacy_name` records the source folder name or browser display name.
- `document` is the source record with known identity fields normalized and unknown fields
  preserved. JSON is decoded with `UseNumber`; numeric precision and nested values are not lost.
- Known reference fields are normalized with the existing `worldstore.NormalizeReferences` shape.
  Duplicate/name-only reference transformations are reported in the plan; unknown fields are not
  stripped.
- Inline browser media is not copied. A record containing unextracted binary data is quarantined as
  `inline_media_unmigrated`; image-byte migration belongs to F3-F. Reference-only fields such as
  `{id,name,type,size,stored}` are preserved.
- Inline browser lore remains in the document. Only actual legacy lore files become
  `world_lore.json` / `scenario_lore.json` or the automated counterparts.
- `media.json` is copied only when the legacy source has a valid reference-only file; image bytes
  are never written into item folders.
- Source timestamps are preserved when parseable; otherwise the import time is used. Original
  timestamp fields stay in `document` as supplied.

Manual lore:

- entries are preserved as raw JSON values, including unknown fields;
- an empty owner is stamped with the mapped canonical item ID;
- a non-empty owner that does not equal the mapped item ID is `unmapped_owner`; the item is
  quarantined rather than rewriting the owner;
- manual and automated lore are never merged.

Automated lore:

- entries receive `Automated: true` and `originKind: "legacy_generated_unattributed"`;
- no chat ID or chat attribution is invented;
- entries may receive the mapped item ID as owner because item ownership is established by the
  mapping;
- no promotion occurs.

## 7. Destination and store boundary

A10 remains in force: only `worldstore` writes canonical item folders. The import orchestrator never
calls `library.Registry.WriteFile`/`WriteNewFile` directly.

Proposed store surface (design only; not implemented):

```go
type ImportItem struct {
    ID                   string
    DestinationLibraryID string
    Document             map[string]interface{}
    Origin               Origin
    CreatedAt            string
    Lore                 map[string][]byte // allowlisted filenames only
    Media                []byte            // optional reference-only media.json
}

type StagedImport struct {
    OperationID  string
    ItemID       string
    Root         string
    StagingPath  string
    FinalPath    string
    DocumentHash string
    LoreHashes   map[string]string
    Revision     string
}

func (s *Store) StageImport(item ImportItem, operationID string) (StagedImport, error)
func (s *Store) CommitImport(staged StagedImport) (SaveResult, error)
func (s *Store) AbortImport(staged StagedImport) error
func (s *Store) RemoveImported(id, expectedRevision string, expectedHashes map[string]string) (SaveResult, error)
func (s *Store) ProposeID(category library.Category) string
```

Rules:

- `StageImport` validates the ID, envelope, lore filenames, sizes, and destination role. It writes
  only under a hidden staging path, never into the final item folder.
- Staging path: `<destination root>/.howl-import-staging/<operationID>/<itemID>/`. The staging
  directory is not an item; its first segment is not a valid ID.
- `CommitImport` re-verifies staged hashes and requires the final path to be absent; it atomically
  renames staging to `<root>/<id>/`. It never replaces an existing item.
- `AbortImport` removes only the staging path.
- `RemoveImported` removes a canonical item only when its current revision equals the imported
  revision **and** every journaled document/lore/media file still matches its imported hash. Any
  unexpected file in the item folder also refuses removal. If the user edited the document, a lore
  file, media, or added content after import, removal is refused with `rollback_conflict`. Document
  revision alone is not sufficient because lore is outside the document revision.
- Adoption performs no store write: the orchestrator compares `Item.Document`/revision, lore
  semantics, and normalized reference-only media bytes, then journals `adopted`. Missing versus
  present media is a conflict, not an adoption.

Destination requirements:

- registration category matches, root exists, `RoleSave` is present;
- `RoleWrite` is required for an import destination so rollback is possible; a save-only root is
  reported as `destination_not_rollback_safe` and is not imported;
- every root is searched for the identity; any unavailable root makes the item indeterminate;
- no import writes into a read-only or external root unless it has the required roles.

## 8. Journal, interrupts, and rollback

Journal files live outside content roots:

- `<config>/worlds-import-journal.json`
- `<config>/scenarios-import-journal.json`
- completion marker: `<config>/worlds-import-marker.json` / `scenarios-import-marker.json`

Journal writes are atomic (temp + rename). The journal is app metadata, not canonical content.

Entry states:

- `prepared` — staged files exist; final path absent; journal written;
- `committed` — final path exists; revision/hashes recorded;
- `adopted` — identical destination already existed; no writes;
- `quarantined` — no writes; reason and evidence recorded;
- `rolled_back` — created paths removed, including prepared staging removed by rollback;
- `failed` — no commit; cleanup pending or completed;
- `cleanup_failed` — owned staging or its marker could not be verified or removed;
- `cleanup_complete` — owned staging of a failed or quarantined item was removed and no canonical item exists.

Per-item execute order:

1. orchestrator revalidates the plan and destination authority;
2. store stages the item (document, lore, media);
3. orchestrator appends/replaces the `prepared` journal entry;
4. store commits the staging rename;
5. orchestrator records `committed` and then writes/updates the operation marker when no pending
   safe work remains.

Crash recovery (only on explicit resume, never at startup):

| Observed state | Action |
| --- | --- |
| Staging exists, no journal entry | remove staging only; no canonical change |
| `prepared`, staging exists, final absent | verify hashes; commit or roll back |
| `prepared`, final exists (crash after rename) | verify final hashes; mark committed |
| `committed`, marker absent | verify final hashes; complete marker |
| `failed`/`quarantined`, owned staging present | verify marker and exact staged files, then remove staging |
| `failed`/`quarantined`, no staging or marker | mark `cleanup_complete`; no canonical change |
| `cleanup_failed` | retry only after ownership evidence is verified; otherwise remain visible |
| Marker present | no-op |
| Source/destination changed | conflict; do not overwrite; require new preview |

Rollback:

- iterates only journaled `committed` entries with `created_paths` and recorded file hashes;
- calls `RemoveImported` with the imported revision and the journaled document/lore/media hashes;
- refuses changed items or items with unexpected files (`rollback_conflict`), never deletes adopted
  items, never touches sources;
- inspects prepared, failed, quarantined, ambiguous, cleanup-failed, and committed item states;
- records prepared staging removal as `rolled_back`; `cleanup_complete` is reserved for cleaned
  failed/quarantined dispositions and is the only cleanup-complete state counted as quarantine;
- reports `rollback_partial`/`cleanup_failed` rather than `rolled_back` when owned staging or a
  committed item's marker evidence remains;
- removes the marker only after all requested removals succeed or are reported;
- reports cleanup failures explicitly.

## 9. Idempotency and conflicts

A plan hash covers: category, destination library ID, every source fingerprint, every mapping
decision, conversion version, and intended write set. Execution re-derives the plan and hash before
writing; the hash is a comparison token, not a capability.

- Executing the same plan against unchanged sources/destination is `already_imported` / `adopted`,
  with no writes.
- A changed source after preview is `source_changed` / `plan_stale`.
- A changed destination after preview is `destination_changed` / `plan_stale`.
- A changed source after a completed import is a conflict; the earlier canonical item is not
  overwritten.
- A changed destination after import is a conflict; no merge.
- A repeated operation with new quarantine resolutions is a new plan, not an automatic retry.

No silent merge is allowed for documents or lore. Manual/automated lore files remain separate.

## 10. Quarantine and later resolution

Quarantine is a retained, visible outcome:

- no canonical write for that item;
- source content remains untouched;
- reason and evidence are journaled and returned to the UI;
- not prompt-eligible and not counted as imported.

Minimum reason codes:

`missing_document_id`, `invalid_document_id`, `duplicate_identity`, `duplicate_name`,
`ambiguous_identity`, `identity_conflict`, `lore_conflict`, `unmapped_owner`,
`unsupported_source_file`, `unsupported_media_manifest`, `media_conflict`, `inline_media_unmigrated`,
`destination_unavailable`, `destination_not_rollback_safe`, `authority_indeterminate`,
`source_destination_overlap`, `source_registered`, `cleanup_failed`, `write_failed`.

Resolution is explicit:

- map to an existing canonical item (identity evidence recorded);
- assign a new store-generated ID (journaled mapping);
- skip the source.

Resolution requires a new preview and plan hash. It never changes a prior canonical item
automatically and never promotes lore.

## 11. Renderer/API shape

Proposed Wails bindings (design only):

```go
PreviewContentImport(req ImportPreviewRequest) (ImportPlan, error)
ExecuteContentImport(req ImportPreviewRequest, planHash string) (ImportResult, error)
ResumeContentImport(category, operationID string) (ImportResult, error)
RollbackContentImport(category, operationID string) (RollbackResult, error)
ListContentImportJournal(category string) ([]ImportOperation, error)
ResolveContentImportEntry(req ImportResolution) (ImportPlan, error)
```

`ImportPreviewRequest` contains the category, destination library ID, browser records (deep copy),
and explicitly selected legacy roots. `ImportPlan` contains the plan hash, source classifications,
per-item dispositions, conflicts, quarantines, and write counts. Execute accepts the request and the
hash, then re-derives the plan; it never trusts the plan object as a write authorization.

The import page is the only caller. It must show the plan before enabling execute, and must not
auto-execute on page load.

## 12. Native acceptance strategy

A9 and the remaining A10 runtime check run in an isolated native process with temp source/destination
roots and temporary config. Browser records are fixtures passed through the same binding shape.

A9 fixtures and assertions:

1. Preview writes nothing: hash the destination tree and source tree before/after; assert equal.
2. Browser world/scenario import preserves unknown top-level and nested fields.
3. Legacy folder exact-ID import writes canonical envelope and lore; automated lore is
   `legacy_generated_unattributed`; manual owner rules hold.
4. Duplicate names, duplicate IDs, same title with different IDs, and missing IDs quarantine; none
   merges.
5. Existing destination identical content is adopted with no writes; differing document is
   `identity_conflict`; differing lore is `lore_conflict`.
6. Unavailable/indeterminate destination root stops the affected item; no success claim.
7. Interrupt after stage, after journal, and after rename; explicit resume completes or rolls back.
8. Rollback removes only journal-created items whose document, lore, and media hashes are unchanged;
   a lore-only edit, an unexpected file, an adopted item, and a source folder all survive rollback.
9. Rerun is idempotent; changed source/destination conflicts without overwrite.
10. Unresolved owners and quarantined items remain in the report and are excluded from prompts.
11. No import occurs on startup, chat load, or editor open.

A10 runtime evidence:

- hash the canonical item tree before and after each operation;
- invoke chat load/save with world/scenario IDs, including the generated-lore sync path;
- assert the tree is unchanged;
- invoke `ExportContentItem` for worlds/scenarios and assert refusal with no target;
- retain the source audit and the new no-write assertions as the runtime half of A10.

A13 remains unchanged: lorebook identity/privacy suites, frontend identity tests, conflict
scope-boundary test, and the chat/store suites.

A12 passed separately: the production prompt path was captured in a test build, and the shipped
binary contains no import seam or prompt-capture symbols. See
`VERIFICATION-F3C-NATIVE-A12-2026-11-10.md`.

Every crash boundary in §15 is exercised by deterministic fault injection in the import service
tests: after the operation journal entry, before and after staging, after the `prepared` journal
write, after the canonical rename but before the `committed` journal write, and after the final
terminal journal write. Faults are injected at the service boundary before store commits, never by
changing production code paths; recovery is only invoked explicitly by the test or user.

## 13. Security and data safety

- Validate every ID with the store validator; reject traversal, absolute paths, reserved names, and
  symlinks/reparse points.
- Legacy source folders are read-only. No write, rename, delete, or normalization occurs there.
- Never write a journal, marker, staging folder, or temp file inside a source folder.
- Enforce document/lore size limits (8 MiB document, 4 MiB lore) and total plan limits.
- Decode JSON with `UseNumber`; do not silently drop unknown fields.
- Do not log source content, full record payloads, or user file bodies.
- Browser records are copied, not mutated; no `localStorage` write from the import path.
- The journal may contain paths and hashes; it must not contain image bytes or secrets.
- Import is local-only; no network calls.

## 14. Implementation sequence

1. Store import primitives (`StageImport`, `CommitImport`, `AbortImport`, `RemoveImported`,
   `ProposeID`) with unit tests for no-overwrite, revision safety, limits, and atomic staging.
2. Plan/preflight model and observational tests (zero writes).
3. Identity resolver, reference normalization, and quarantine reasons.
4. Browser adapter and lossless conversion tests.
5. LF-1 legacy adapter, lore classification, and source-preservation tests.
6. Journal state machine, interrupt injection, resume, and rollback.
7. Idempotency/conflict tests and plan-hash staleness.
8. Wails bindings and import page; no editor/snapshot wiring.
9. Isolated native A9/A10 evidence; A13 regression.
10. Only then consider cutover or editor integration.

## 15. Store-boundary resolutions (implementation gate)

These are the binding resolutions for the approved implementation. They refine, but do not weaken,
§3–§12.

### A. Multi-item crash consistency

Atomicity is **per item**, not transaction-wide. A multi-item operation may end with earlier items
committed and later items quarantined or failed; it never claims all-or-nothing. The journal is the
operation ledger, and the completion marker records `completed` or `completed_partial`.

Journal states per item: `planned`, `prepared`, `committed`, `adopted`, `quarantined`, `rolled_back`,
`failed`, `ambiguous`.

Crash boundaries and recovery:

| Boundary | Durable evidence | Recovery decision |
| --- | --- | --- |
| After operation journal, before staging | `planned`, no staging, no final | Mark failed/aborted; no content change |
| During staging | `planned`, partial staging, owner marker if written | Remove staging only if the ownership marker matches this operation/item; otherwise `ambiguous` and stop |
| After complete staging, before `prepared` | `planned`, complete staging with owner marker | Verify owner marker and staged hashes from the marker; either complete staging journal or abort staging. No canonical write |
| After `prepared`, before commit | `prepared`, staging present, final absent | Revalidate authority; verify staged hashes; commit or abort |
| After rename, before `committed` | `prepared`, final present, staging absent | Verify final file hashes recorded at `prepared`; record `committed`. Never assume the rename did not happen |
| After `committed`, before marker | `committed`, final present | Verify hashes; complete marker |
| After marker | marker + journal | No-op |
| Any state with final present but no verifiable hashes | journal without file evidence | `ambiguous`; no deletion or overwrite |

A rename that succeeds but is not journaled is therefore recoverable from `prepared` plus final
hashes. A final path that appeared without journal evidence is never adopted automatically.

### B. Commit-time authority

`CommitImport` runs under the same process-wide store write lock as `Save` and `Update`. Under that
lock it:

- re-reads the destination registration and requires the same root, category, save role, write role,
  availability, and rollback capability;
- calls `Lookup` across every registered root and rejects unavailable roots, duplicate copies, or
  indeterminate authority;
- requires the final item path to be absent;
- re-verifies every staged file hash;
- atomically renames only after those checks.

A final item that appears after preview is a collision. It is never overwritten. Cross-process
writers remain outside the desktop single-process contract; the process lock does not claim to
serialize an unrelated external process.

### C. Staging isolation

Staging is `<destination root>/.howl-import-staging/<operationID>/<itemID>/` with an ownership marker
`<itemID>.owner.json` beside that folder (file name `owner.json`, not inside the item folder), containing a magic value, operation ID, item ID, category, plan hash, and
staged hashes. The first path segment is not a valid item ID, and `Lookup` only reads
`<root>/<id>/world.json` or `scenario.json`, so staging is not discoverable as a canonical item.
Chat loading, startup, and ordinary item enumeration do not scan staging. Staging and final are in
the same destination root, so the rename is same-filesystem.

Cleanup is operation-scoped: a staging directory is removed only when its ownership marker matches
the requested operation and item and its document/lore/media set is exactly the journaled set.
Unexpected files, unmarked directories, another operation's staging, and final items are never
removed. An abandoned staging directory with no marker is reported as `cleanup_failed` or
`ambiguous`, not deleted.

### D. Concurrency

- Staging may run concurrently because each operation has a unique staging path.
- `CommitImport`, `AbortImport`, `RemoveImported`, `Save`, and `Update` all use the existing
  process-wide `writeMu`; there is no independent import lock that bypasses it.
- Two imports of the same identity may both stage; only the first commit succeeds, the second
  returns `collision`.
- An import commit and an ordinary create of the same identity are serialized; exactly one wins.
- An import commit and an ordinary update are serialized. If the item exists, the import cannot
  create or replace it; if the import creates first, the update sees the new revision or a missing
  base item.
- Journal writes use a separate short-lived journal mutex only to keep the JSON file consistent;
  that mutex never authorizes or performs canonical writes.

### E. Recovery and rollback proof

Recovery requires journal ownership evidence: category, operation ID, plan hash, item ID, staging
marker, and recorded hashes. Missing or conflicting evidence yields `ambiguous`; no file is changed.

Rollback is two-phase:

1. Verify the whole operation: every committed entry has matching document/lore/media hashes, no
   unexpected files, no adopted item, no ambiguous item, no leftover owned staging or marker, and
   complete journal evidence. Any failure stops before any deletion.
2. Remove only verified, unchanged, journal-created items through `RemoveImported` under `writeMu`,
   and remove verified failed/quarantined staging through `AbortImport`. Adopted items, source
   folders, and browser records are untouched. A cleanup failure leaves the operation
   `rollback_partial` with `cleanup_failed`, never `rolled_back`.

A later document or lore edit, an added file, a missing file, a revision change, an unavailable
authority, or incomplete evidence produces `rollback_conflict` or `ambiguous_state` and leaves the
item intact. Source hashes are checked before and after every recovery/rollback test.

### F. Reference representation

The import boundary now normalizes known reference fields to the frozen typed shape without
discarding unknown document or per-reference fields: `internal/contentimport/references.go` copies
maps losslessly, adds or corrects `kind`, preserves `id`/`title`/`legacyName`, collapses duplicate
identities, and reports legacy-name, missing-ID, malformed, duplicate, unknown-kind, and kind-mismatch
cases. Name-only references become `{kind, id: "", legacyName}` and remain unresolved. Unknown-kind
or mismatched-kind entries are preserved in the document but are skipped by runtime
`worldstore.NormalizeReferences`, so a title or an unrecognized kind can never resolve to an ID.
`legacyName` is retained only as migration/display evidence; resolution remains ID-only. This
implements the frozen §6 rule rather than amending it.

## 16. Status

- Design: **recorded and implementation-authorized after §15 resolutions; corrected 2026-11-10**.
- Store import primitives: **implemented in `internal/worldstore/import.go`**.
- Import/migration service: **implemented in `internal/contentimport`** (browser + LF-1, journal,
  explicit resume and rollback). Audit findings F-01 through F-05 and F-07 through F-09 are corrected;
  F-06 bindings were regenerated and the native build passed.
- UI/editor wiring: **not started** (outside this authorization).
- A9: implementation tests and matrices recorded in
  `docs/review/IMPLEMENTATION-F3C-IMPORT-2026-11-10.md`; native runtime evidence remains open.
- A10: source audit passes; native runtime evidence remains open.
- A13: existing gate intact.

## 17. Final marker and media correction — 2026-11-10

### A. Completion-marker authority

The completion marker is operation-scoped at
`<configDir>/<category>-import-marker-<operationID>.json`. It is not a category-wide singleton. The
marker records the marker version, category, operation ID, plan hash, destination library and root,
sorted item IDs, terminal status, and completion time.

- `completed` and `completed_partial` require a marker.
- A marker is written only after the journal reaches a terminal completion state.
- Creation is exclusive when the path is absent.
- An existing marker may be replaced only when category, operation, plan, destination, item set, and
  terminal status all match. A foreign or contradictory marker is never overwritten.
- `execute`, `resume`, and `rollback` verify the marker before acting. Missing, malformed, stale,
  contradictory, or foreign evidence produces `integrity_failure`; canonical content is not deleted
  or inferred into completion.
- The read-only journal listing verifies markers too and reports `integrity_failure` in its returned
  copy without changing the journal.
- `rollback_partial` removes the completion marker only after ownership is verified.
- A marker surviving in a non-completion journal state is stale evidence and blocks recovery.
- Result actions keep state and counts aligned: committed/rolled-back items are `ready`, adopted
  items are `adopt`, planned/prepared items are `pending`, and every other state is `quarantine`.

This keeps the journal as the operation ledger and the marker as terminal evidence. Neither alone
is treated as proof of completion when both should exist.

### B. Media validation classes

- `ErrMediaMalformed` covers undecodable JSON, missing/blank canonical `format`, non-object items,
  missing `mediaId`, malformed legacy entries, and non-array supported envelopes.
- `ErrMediaUnsupported` covers a structurally valid envelope whose format or role is not supported.
- Preparation maps these to `malformed_media_manifest` and `unsupported_media_manifest`.
- Inline binary remains separately classified as `inline_media_unmigrated`; no media bytes are
  migrated.
- Normalization preserves unknown top-level and per-item fields, and the same normalized bytes are
  used for staging, hashing, load, adoption equality, and reload.

### C. Reference integration boundary

The import slice preserves and classifies typed references, including unknown-kind, kind-mismatch,
missing-ID, malformed, and legacy-name-only outcomes. It does not add a runtime name/title fallback
or widen the reference contract. No production consumer of the resolver exists in this F3-C slice;
that absence is an explicit integration boundary and is not claimed as end-to-end runtime
resolution.

### D. Status

- Marker integrity, cleanup aggregation, media classification, and reference edge evidence are
  corrected and tested.
- A9 and A10 later passed in an isolated native run; see
  `VERIFICATION-F3C-NATIVE-A9-A10-2026-11-10.md`.
- A13 remains a separate regression gate and passed after this correction.
