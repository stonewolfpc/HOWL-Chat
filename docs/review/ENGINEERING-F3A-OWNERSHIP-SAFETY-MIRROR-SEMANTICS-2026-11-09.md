# F3-A — Content Ownership Safety and Mirror Semantics

Date: 2026-11-09

Scope: bounded safety slice before any content ownership migration

Status: implemented and source/test verified; runtime UI acceptance not performed

## 1. F2-C findings and dispositions

| F2-C finding | F3-A disposition | Evidence class |
|---|---|---|
| F2-B automatic JSON copies were written after browser saves and never read back. | Removed from all automatic page-save paths. The F2-B helper and its script tags are gone. | Source-confirmed |
| World/scenario/lorebook mirror writes could fail while the UI still reported success. | Removed the optional mirror write, so the remaining save result describes the authoritative browser store only. World's success alert now follows only `storeWorld`. | Source-confirmed |
| Character saves could finish while the library copy failed silently. | Removed the mirror call. Character save still reports its browser-store failure through its existing catch path. | Source-confirmed |
| Library writes used title-based paths and `RoleSave` could replace files. | Export paths now include a stable identity hash; replacement requires `RoleWrite` and an explicit `replace=true` request. | Source-confirmed / Test-confirmed |
| Deletes and renames did not propagate to library copies. | No automatic propagation was added. External registration removal remains association-only; lifecycle synchronization is deferred to the ownership migration. | Source-confirmed / Deferred |
| Gallery images and sidecars were written non-atomically, and gallery failures only printed a warning. | Image and sidecar writes use a temporary file plus rename; a sidecar failure removes the just-written image; generation now returns `gallery_error` as a partial result. | Source-confirmed / Test-confirmed |
| Gallery reads accepted arbitrary paths; deletes swallowed errors. | Reads resolve through `ResolveFile` with lexical and symlink containment; deletes validate names and return real errors. | Source-confirmed / Test-confirmed |
| `WipeAllData` targeted `exeDir/chats`, not the AppData chat root, and swallowed failures. | Wipe now uses `PersistenceManager.RootDir()`, preflights both roots, reports failures, and only clears the chat and gallery roots. | Source-confirmed / Test-confirmed |
| Unavailable external character roots made their cards disappear. | Character inventory now persists the last successful scan and returns stale items marked unavailable. | Source-confirmed / Test-confirmed |

## 2. Exact files and packages changed

- `app-libraries.go`: removed `SaveContentItem`; added `ExportContentItem`, stable export identity, safe segments, structured results, and inventory.
- `app-libraries_test.go`: export, collision, replacement, unavailable-root, and registration-removal coverage.
- `app-wipe.go`, `app-wipe_test.go`: real chat-root resolution, scoped deletion, failure reporting.
- `app-image-generate.go`, `app-image-gallery_test.go`: gallery-bound reads and partial gallery-save reporting.
- `internal/library/registry.go`, `internal/library/registry_test.go`: typed errors; save-only create versus write-role replacement; collision and failed-write preservation tests.
- `internal/library/result.go`: operation status/result contract.
- `internal/library/inventory.go`: persisted last-known inventory cache.
- `internal/image/gallery/gallery.go`, `internal/image/gallery/gallery_test.go`: gallery containment, atomic writes, truthful delete errors.
- `internal/image/types/types.go`: `ImageResult.GalleryError`.
- `internal/backend/chat/persistence.go`: `PersistenceManager.RootDir()`.
- `app.go`: inventory store initialization.
- `frontend/dist/js/character-form.js`, `world-form.js`, `scenario-form.js`, `lorebooks.js`: deleted automatic mirror calls; lorebook storage failures now roll back the in-memory change and report an error.
- `frontend/dist/js/content-libraries.js`: unavailable/stale inventory display.
- `frontend/dist/js/image-generator-core.js`: displays gallery partial-save error.
- `frontend/dist/js/app-wipe.js`: stops frontend erasure when backend wipe fails; confirmation states that external libraries are untouched.
- `frontend/dist/character.html`, `world.html`, `story-tools.html`, `scenario.html`, `lorebooks.html`: removed `content-save.js` tags.
- `frontend/dist/js/content-save.js`: deleted.
- `PROJECT_ISSUES_AND_EXPANSION.md`: F3-A decisions appended.
- This report.

## 3. Ownership and operation-result contracts

Normal page saves now have one browser-store write path. They do not create library copies. Existing F2-B JSON files are left untouched on disk; they are not read automatically and are not authoritative.

`internal/library.OperationResult` applies to content-library export operations, not normal browser saves. It distinguishes:

- `save_succeeded`
- `export_succeeded`
- `partial`
- `failed`
- `destination_unavailable`
- `permission_denied`
- `name_collision`
- `replaced`
- `cancelled`

`ExportContentItem` is the only content-library write binding. It writes a `howl-content-export-v1` envelope containing `authoritative: false`, the category, item identity, timestamp, and original content. A missing save destination returns `destination_unavailable` without writing. The binding is explicit; no page-save workflow calls it.

Gallery generation uses a separate partial result: `ImageResult.GalleryError` is set when the generated bytes exist but the gallery copy failed. The image remains returnable, and the UI shows the warning instead of claiming a successful gallery save.

## 4. Collision and overwrite protections

- New exports use `<first-12-hex-of-identity>/<category>.json`; the display name is stored in the export envelope rather than used as identity.
- `WriteNewFile` returns `ErrCollision` when the target exists.
- `WriteFile` returns `ErrReadOnly` for a save-only root, even when the target exists.
- Replacement is possible only on `RoleWrite` roots and only with `replace=true`; it returns `replaced`.
- Invalid Windows characters, control characters, leading/trailing dots and spaces, and reserved device names are neutralized.
- Path traversal is rejected by the registry before directory creation.
- The registry writes to a temporary file and restores the prior file if commit fails. A focused test blocks the staging path and confirms the old content remains.

## 5. Gallery security boundary

`gallery.Service.ResolveFile` requires a non-empty path, resolves the gallery root and target absolutely, evaluates symlinks/reparse points where they exist, and checks containment after resolution. `Read`, `List`, and export image-resolution use it. `GetGalleryImage` delegates to it, so arbitrary paths are rejected.

`Delete` accepts a base name only, rejects separators, absolute paths, `.` and `..`, and normalizes a supplied `.png`, `.json`, or `.webp` suffix. It removes the image, sidecar, and thumbnail; missing companions are tolerated, but a real error or a wholly missing item is reported.

Gallery images and metadata use `writeFileAtomic` (temporary file, sync, rename). If the metadata commit fails, the new image is removed so no image is left without its sidecar.

## 6. Reset and wipe behavior

Settings reset remains separate: `ResetApplicationSettings` writes settings and image settings and does not call `WipeAllData`; its existing test verifies that creative files are untouched.

`WipeAllData` now:

1. Requires both services before deleting anything.
2. Resolves the chat root through `PersistenceManager.RootDir()` and removes only `<root>/chats`.
3. Removes and recreates only the gallery service's base directory.
4. Does not guess executable-relative entity folders or image subdirectories.
5. Does not traverse or delete external library roots.
6. Stages both roots to same-filesystem siblings before resetting state, restores the first if the second cannot be staged, resets `LastChatID`, clears the inventory cache, and only then removes the staged data.
7. Reports recoverable staged data if cleanup fails rather than calling the wipe complete.

The frontend calls the backend first. If it fails, the frontend does not clear IndexedDB/localStorage and does not reload as though the wipe succeeded. A fully previewed, separately confirmed comprehensive cleanup remains future work.

## 7. Tests and exact results

Focused commands:

```text
go test ./internal/library ./internal/image/gallery .
ok  howl-chat/internal/library
ok  howl-chat/internal/image/gallery
ok  howl-chat

go build ./...
exit 0

go test -race ./internal/library ./internal/image/gallery .
ok  howl-chat/internal/library
ok  howl-chat/internal/image/gallery
ok  howl-chat

node --check (seven changed frontend files)
exit 0

git diff --check (all changed tracked files and pages)
exit 0
```

New coverage includes save-success/export-unavailable separation, collision without overwrite, write-role replacement, failed-write preservation, stale inventory with unavailable items, external registration removal without file deletion, gallery traversal and symlink escape (skipped if symlinks unavailable), gallery delete errors, actual chat-root wipe, and gallery reads through the app binding.

An earlier full `go test ./...` run was **not clean**: `internal/audio` failed `TestConcurrentDifferentConfigurationsRemainIsolated` with `hits 5 misses 2, want 0/2`. The test passed in 10 consecutive isolated runs, and the final F3-A-R1 full-suite run passed. No audio code was changed; the intermittent failure remains unresolved even though it did not reproduce in the final run.

## 8. Runtime checks performed and not performed

Performed: source inspection, focused unit tests, race tests, `go build`, JavaScript syntax checks, and diff whitespace checks.

Not performed: launching the Wails application, clicking through reset/wipe confirmations, observing native folder pickers, testing a real unavailable external drive, or exercising gallery rendering. No runtime-confirmed claim is made.

## 9. Remaining risks and blockers

- Legacy F2-B JSON files remain on disk and may be mistaken for current content by a user. They are no longer written or read automatically, but no cleanup or migration was authorized.
- Delete/rename propagation and Howl-owned index synchronization remain deferred to the ownership migration.
- Inventory has a live scanner for characters only. Other categories show registrations and explicitly state that no live scanner exists.
- Chat, image, and video export are not part of this slice. Videos still lack a sidecar/catalog.
- Wipe deletion is scoped but not transactional across two roots; a second-root failure returns an error after the first root was cleared. The frontend stops and reports it.
- Overlapping registrations are still a registry-level concern; F3-A prevents path traversal and filename collisions but does not deduplicate physical content across overlapping roots.
- No explicit export UI was added. The explicit backend operation and structured contract exist; export is currently backend-only.
- Export identity and replacement safety were tightened in F3-A-R1; the old `.bak` cleanup path and parent-directory symlink traversal remain outside this slice.

## 10. Ledger changes and next slice

`PROJECT_ISSUES_AND_EXPANSION.md` records F3-A1 through F3-A8. F2-B future requirements remain preserved and deferred.

The next ownership slice should be **F3-B: lorebook ownership cutover**, after this gate. It should use stable IDs, write through the authoritative lorebook store, keep exports explicit, and prove unavailable-root and rollback behavior before touching worlds, scenarios, characters, chats, images, or videos.
