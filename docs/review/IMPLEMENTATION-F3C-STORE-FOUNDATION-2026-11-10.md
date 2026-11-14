# IMPLEMENTATION — F3-C Store Foundation

**Date:** 2026-11-10
**Scope:** canonical world/scenario store skeleton, deterministic revision, lossless update,
typed reference normalization. No consumers connected.
**Contracts:** `FREEZE-F3C-DATA-AND-OWNERSHIP-CONTRACTS-2026-11-10.md`
**State:** implemented and unit-tested; import, chat integration, snapshot refresh, and media are
deliberately absent.

## 1. What was built

| File | Purpose |
| --- | --- |
| `internal/worldstore/store.go` | Envelope, identity, create, lossless update, load, role enforcement, structured failures |
| `internal/worldstore/revision.go` | Deterministic canonical JSON and content revision |
| `internal/worldstore/refs.go` | Typed references, normalization, ID-only resolution |
| `internal/worldstore/store_test.go` | Create, collision, roles, lossless merge, revision conflict, tamper detection |
| `internal/worldstore/revision_test.go` | Determinism: key order, number spelling, arrays, null, timestamps, media, Unicode |
| `internal/worldstore/refs_test.go` | Legacy names, dedupe, ID-only resolution, kind separation |

No existing file was modified. The package is isolated so it cannot absorb import or consumer work
by accident.

## 2. Contracts as implemented

### Identity and paths

- `^[A-Za-z0-9][A-Za-z0-9._-]{0,79}$`; path traversal is rejected before any filesystem call.
- Existing IDs, including bare timestamps and both scenario shapes, are accepted unchanged.
- New items get `wld_<16 hex>` or `scn_<16 hex>`.
- Paths are `<root>/<id>/world.json` and `<root>/<id>/scenario.json`; a rename cannot move a file
  because the path never contains a title.

### Revision

- `rv1-` + 16 hex of SHA-256 over canonical JSON.
- Keys sorted; array order preserved; `null` distinct from omission.
- Integers kept as exact literals (precision beyond float64 preserved); fractional/exponent forms
  normalized through float64.
- UTF-8 emitted without Unicode normalization; composed and decomposed sequences differ.
- Excluded: `updatedAt`/`updated_at`, `metadata.updated`, and image payload keys
  (`data`, `medium`, `original`, `path`, `thumbnail`).
- `Load` recomputes the revision. A mismatch is reported in `revision_mismatch`; the document stays
  readable and is never trusted as current.

### Lossless update

- `Save` is create-only and requires a save root.
- `Update` requires the write role on the owning root and a non-empty `baseRevision`.
- A missing or stale base returns `revision_conflict` with `current_revision`, and the file is not
  touched.
- Known fields are applied; unknown top-level keys survive; `howl` and `metadata` merge key by key.
- Arrays replace wholesale when present.
- Explicit `null` clears only the allowlist (`description`, `details`, `notes`, `image`, `images`,
  `thumbnail`, `thumbnailIndex`, `world`, `worldId`, `image_style`, `avatar`); a null elsewhere is
  ignored rather than permission to erase data.
- `id`, `createdAt`, `created_at`, and `revision` are immutable; a change is `ErrInvalidDocument`.

### Roles and failures

- Create: save root only.
- Replace: write role on the root that already holds the item.
- Structured statuses: `saved`, `updated`, `collision`, `permission_denied`, `unsafe_target`,
  `invalid_path`, `write_failed`, `not_found`, `revision_conflict`, `destination_unavailable`,
  `authority_indeterminate`.
- Writes go through the library registry's atomic replace, so a failed write leaves the prior
  document intact.

### References

- Shape `{kind, id, title?, legacyName?}`; kinds are world, scenario, character, lorebook.
- Name-only legacy entries normalize to `legacyName` with an empty ID and resolve as
  `legacy_name_unresolved`.
- Resolution is by kind plus ID only. A same-title item never satisfies a reference.
- Duplicates collapse on identity, or on legacy name when no identity exists.
- An identity in one kind never resolves in another.

## 3. Verification

| Check | Result |
| --- | --- |
| `gofmt -l internal/worldstore` | clean (no output) |
| `go vet ./internal/worldstore` | exit 0 |
| `go test ./internal/worldstore ./internal/library` | both packages `ok` |
| Store contract tests | 13 tests, all pass |
| Revision determinism tests | 5 tests, all pass |
| Reference tests | 4 tests, all pass |
| Concurrency and atomic-write tests | 4 tests, all pass |
| Availability tests | 4 tests, all pass |
| Root authority tests | 3 tests, all pass |
| Canonical lore read tests | 4 tests, all pass |
| `go test -race ./internal/worldstore` | pass |
| A13 regression: `lorebook`, `chat`, `memory`, root packages | all `ok`, exit 0 |
| Frontend character identity and memory identity tests | exit 0 (4/4 each) |

The lossless test is the C6 guard: an update that touches only `title` and `howl.lorebooks` preserves
`legacyNotes`, `howl.notes`, `howl.customFlag`, and `metadata.extra`. The reference test is the C4
guard: a catalog entry with the same title but a different identity does not resolve.

## 4. Deliberately absent

- Import, preflight, journal, rollback.
- Import, preflight, journal, rollback.
- Chat service conversion to a consumer; generated-lore relocation. (Completed 2026-11-10;
  record: `IMPLEMENTATION-F3C-CHAT-CONSUMER-2026-11-10.md`.)
- Snapshot revision fields, staleness display, refresh operations.
- Media manifests and any binary handling.
- Editor/repository wiring and cache hydration.
- Availability index for disconnected roots; `Lookup` exposes search completeness, but persisting
  root availability across restarts belongs to a later step.

These belong to later steps and each has its own acceptance criteria. Connecting the editor before
the consumer cutover would create a second read path.

## 5. Limits of this evidence

- Unit and filesystem tests only; no native run, no UI, no prompt-boundary capture.
- The write lock is per process. A second OS process could still lose an update; the desktop
  runtime is single-process, and atomic replacement still prevents partial files.
- No migration has been executed against real records.
- Revision stability is proven within this implementation; it is not a claim that legacy documents
  imported later will hash identically to any previously stored value, because none exist.

## 6. Correctness gates closed — 2026-11-10

### Concurrent and interrupted writes

- Read-check-write is serialized process-wide (`writeMu`), so an update re-reads the document and
  its revision inside the critical section. A second update holding the same `baseRevision`
  observes the winner and returns `revision_conflict`.
- Eight goroutines sharing one base revision: exactly one `updated`, seven
  `revision_conflict`; the stored document equals the winner and its computed revision.
- Two independent `Store` instances over one registry produce the same result: one winner.
- Eight concurrent creates: one `saved`, seven `collision`; the winning document survives.
- Atomicity: blocking the registry temp path makes the update fail as `write_failed`; the
  canonical file is byte-identical afterward, still loads, and its revision is not mismatched.
- `go test -race ./internal/worldstore` passes on this Windows host.
- Limit: the lock is per process, not a filesystem lock. A second OS process could still lose an
  update. The desktop runtime is single-process, and atomic replacement still prevents partial
  documents.

### Missing versus unavailable

`Lookup` returns `Search{Found, Complete, AuthorityIndeterminate, Copies, UnavailableRoots, Errors}`:

- **Found** — the items in every readable root are scanned; the first copy in deterministic path
  order is returned for display. Path order is retrieval order, not authority.
- **Authority determinate** — only when every registered root was readable and searched without
  error and exactly one root holds the identity. Only then is the copy canonical and editable.
- **Found, authority indeterminate** — a root is unavailable, unreadable, or errored, or two or
  more readable roots hold the identity. The copy stays readable with
  `AuthorityIndeterminate=true` and `Editable=false`; `Update` reports
  `authority_indeterminate` and writes nothing.
- **Not found, complete** — every registered root was readable and searched without error, so
  absence is proven; `Load` returns `ErrNotFound`.
- **Not found, incomplete** — a root is missing, unreadable, or errored, so absence cannot be
  proven; `Load` returns `ErrUnavailable`. `Update` reports `destination_unavailable`.
- **No registered root** — incomplete search, `ErrUnavailable`, never `ErrNotFound`.

Root authority is decided in `DECISION-F3C-ROOT-AUTHORITY-2026-11-10.md`: roots are equivalent
locations, never ranked, and copy selection never implies precedence.

Consumer rules, binding for the chat conversion:

1. `Found && AuthorityIndeterminate` — show the item as not-current/unavailable-authority, keep it
   readable for reference, and never edit, recreate, import, overwrite, deduplicate, or reassign
   references. A disconnected sibling root is not a missing item.
2. `Found && !AuthorityIndeterminate` — use the canonical item; editing still requires `Editable`.
3. `!Found && Complete` — genuinely absent; only an explicit user action may create it.
4. `!Found && !Complete` — indeterminate; show as unavailable and never recreate, import,
   overwrite, or reassign references. Never fall back to chat-local or name-keyed data.
5. No state permits silent promotion or invented provenance.

No availability index is added here; this is the documented foundation behavior that consumers
must honor until that index exists.

## 7. Next gate

1. ~~Chat service conversion to read through this store, with generated-lore writes removed (C3).~~
   Completed 2026-11-10; record: `IMPLEMENTATION-F3C-CHAT-CONSUMER-2026-11-10.md`.
2. Reference integration: typed reads in the canonical path, no title fallback (C4).
3. Then import, snapshot revisions, media manifests, and native acceptance.
