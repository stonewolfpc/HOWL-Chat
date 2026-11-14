# F3-A-R1 — Safety Repairs

Date: 2026-11-09
Scope: close the concrete F3-A acceptance defects without starting F3-B
Status: implemented and source/test verified; no runtime Wails verification

## 1. Defects and root causes

| Defect | Root cause | Repair |
|---|---|---|
| Explicit replacement could replace a directory or symlink. | `WriteFile` used `os.Stat` and treated any existing target as replaceable. | `Lstat` inspects the target without following it. Only a regular file may be replaced. Directories and symlinks return `ErrUnsafeTarget` and are untouched. |
| Wipe could leave a misleading partial state. | The two roots were deleted sequentially; `LastChatID` and the inventory cache were not reset; cleanup errors were not represented. | Both roots are staged to sibling paths first. If the second stage fails, the first is restored. Only after both are staged are `LastChatID`, inventory, and memory reset. Cleanup failure reports recoverable staged data. |
| Gallery listing could expose an image path outside the root. | `resolveImagePath` used `os.Stat` without canonical containment. | It now delegates to `ResolveFile` for every candidate, so listing, reading, and deleting use the same boundary. |

## 2. Source files changed

- `internal/library/registry.go`
  - `ErrUnsafeTarget`
  - `Lstat` target inspection
  - regular-file-only replacement
  - injected `renameFile` seam
  - unique backup path
  - explicit restore-failure message
- `internal/library/registry_test.go`
  - directory target
  - symlink target
  - commit failure with successful restore
  - commit and restore failure with recoverable backup
- `internal/library/inventory.go`
  - `Clear()` for completed wipes
- `app-libraries.go`
  - identity-only export directory
  - display name moved into the envelope
  - missing identity/data rejection
  - `ErrUnsafeTarget` mapping
- `app-libraries_test.go`
  - renamed-title identity collision
  - missing identity/data rejection
- `app-wipe.go`
  - same-filesystem staging
  - rollback of the first stage
  - `LastChatID` reset
  - inventory invalidation
  - staged-cleanup reporting
  - failure-injection seams
- `app-wipe_test.go`
  - successful staging/reset
  - first-stage failure
  - second-stage failure with restore
  - cleanup failure with recoverable stage
  - missing-service failure
- `internal/image/gallery/gallery.go`
  - `resolveImagePath` containment
- `internal/image/gallery/gallery_test.go`
  - legitimate in-root image
  - outside-resolving image symlink
  - outside-resolving sidecar symlink
- `docs/review/ENGINEERING-F3A-OWNERSHIP-SAFETY-MIRROR-SEMANTICS-2026-11-09.md`
- `PROJECT_ISSUES_AND_EXPANSION.md`
- This report.

## 3. Replacement safety and identity

`WriteFile` now:

1. Checks the lexical target and root.
2. Uses `os.Lstat`, not `os.Stat`.
3. Returns `ErrCollision` for any existing target when replacement was not requested.
4. Returns `ErrUnsafeTarget` for a directory, symlink, or other non-regular target when replacement was requested.
5. Creates a temporary file only after the target check.
6. Moves the existing regular file to a unique `.bak-<timestamp>` sibling.
7. Commits the temporary file.
8. Restores the backup if commit fails; a failed restore returns an error naming the recoverable backup path.

`RoleSave` alone remains create-only. Replacement requires `RoleWrite` and explicit `replace=true`.

The export directory is now the identity hash only. Renaming an item does not create a second directory. The display name is stored in the envelope, not used for identity. Empty item identity with nil data is rejected.

The `.bak` cleanup after a successful replacement is still best-effort. A cleanup failure can leave a `.bak-*` copy; this is documented rather than silently treated as a rollback guarantee. Parent-directory symlink traversal is also outside this narrow repair.

## 4. Wipe semantics

`WipeAllData` preflights `chatPersistence`, `gallerySvc`, both roots, and current `LastChatID`.

1. Stage `<chatRoot>/chats` to `.<name>.howl-wipe-<timestamp>` in the chat root.
2. Stage the gallery root to its same-filesystem sibling.
3. If the second stage fails, restore the first; state remains untouched.
4. Reset `LastChatID`.
5. Clear the persisted inventory cache.
6. Reset in-memory content state.
7. Remove the staged roots.
8. Recreate the empty gallery root.
9. Return an error naming any staged data that could not be cleaned up.

If `SetLastChatID` fails after both stages, both roots are restored. If inventory clearing fails, the old chat pointer and both roots are restored. Cleanup failure is reported as recoverable staged data, not as a clean wipe.

External registered roots are never staged, removed, or traversed.

## 5. Gallery containment

`List` validates the sidecar with `ResolveFile`. `resolveImagePath` now calls `ResolveFile` for each image candidate, so an image symlink or reparse point resolving outside the gallery is rejected. `Read`, `List`, and `Delete` therefore agree on the root boundary.

## 6. Tests and exact results

```text
go test ./internal/library ./internal/image/gallery .
ok  howl-chat/internal/library
ok  howl-chat/internal/image/gallery
ok  howl-chat

go test . -run TestWipeAllData -count=1 -v
PASS: TestWipeAllDataStagesBothRootsAndResetsState
PASS: TestWipeAllDataFirstStageFailureLeavesStateUntouched
PASS: TestWipeAllDataSecondStageFailureRestoresFirstRoot
PASS: TestWipeAllDataCleanupFailureLeavesRecoverableStage
PASS: TestWipeAllDataReportsMissingRoots

go test -race ./internal/library ./internal/image/gallery .
ok  howl-chat/internal/library
ok  howl-chat/internal/image/gallery
ok  howl-chat

go build ./...
exit 0

node --check (changed frontend files)
exit 0

git diff --check (changed tracked files)
exit 0
```

Full suite:

```text
go test ./...
all packages passed
```

The previously reported `internal/audio` failure did not reproduce in 10 isolated runs or in the final full-suite run. No audio code was changed. The intermittent failure is still treated as unresolved.

## 7. Remaining limitations

- Export is backend-only; there is no export button or page action.
- Only character roots have a live inventory scanner.
- Cached/stale character inventory is not category-wide inventory support.
- Video sidecar/catalog parity is incomplete.
- Delete/rename propagation remains deferred to ownership migration.
- Backup cleanup after successful replacement is best-effort; `.bak-*` may remain.
- Parent-directory symlink traversal is not resolved by this slice.
- No runtime Wails UI or packaged application test was performed.

## 8. Content statement

No user content was migrated, moved, renamed, deleted, or rewritten by F3-A-R1. All destructive-operation tests used isolated temporary directories and injected failures.
