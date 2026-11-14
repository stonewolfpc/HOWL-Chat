# F3-B — Lorebook Ownership Cutover

Date: 2026-11-10

Scope: canonical lorebook ownership, explicit import, read-back through the library adapter,
legacy store retained

Status: implementation complete; automated tests passed; integration verified through the UI harness; native-window acceptance pending

## 1. Starting point

F3-A-R1 was accepted with documented limitations. It removed the unlabelled automatic mirror
writes, made replacement and wipe safe, and left lorebooks browser-owned with no file identity.
The F3-B scope gate recorded five unanswered contract questions. This stage resolves them and
implements the resulting architecture.

## 2. Decisions

| ID | Question | Decision |
|---|---|---|
| F3-B-Q1 | File granularity | One JSON document per entry: `<root>/<id>.json`, envelope `howl-lorebook-v1` |
| F3-B-Q2 | Read path | One canonical store; a frontend repository derives the browser cache, and chat refresh reads canonical before sending |
| F3-B-Q3 | Identity | Preserve existing identities on import; random `lb_` identities for new entries; duplicates reassigned and reported |
| F3-B-Q4 | Import and recovery | Explicit, create-only, idempotent by identity + content hash, journaled, reversible |
| F3-B-Q5 | Delete/rename propagation | No propagation subsystem; canonical delete removes the file, rename retains identity and path |

### Alternatives and why they were rejected

**Granularity.** A single versioned collection document was the alternative. It forces a
read-modify-write of every entry for any edit, creates a larger atomic write, and invents a
collection concept the application does not have — worlds, scenarios, and characters attach
individual entries. Per-entry files match the attachment model and the registry's existing
single-item write protections.

**Read path.** Keeping the page-only read was the alternative. Because chat prompt assembly and
every attachment dropdown read the browser store directly, a page-only read would leave the
browser store authoritative for the rest of the application — the exact unlabelled dual-source
state F3-A removed. Instead, the canonical store is the only writer, the repository hydrates the
cache from it, and chat sync refreshes before prompt assembly. Synchronous readers keep working
against a rebuildable cache.

**Identity.** Minting new identities at import would have required rewriting every
`howl.lorebooks` reference across worlds, scenarios, and characters — a user-content rewrite
that F3-B is not authorized to perform. Preserving identities keeps REG-1 attachments intact.
New entries get random identities from the start, so persistence after import is stable.
Duplicates that already share an identity cannot be distinguished by reference; the second is
imported under a new identity and reported for review rather than silently re-linked.

**Import.** A silent or automatic migration was rejected outright. Import is an explicit page
action over all browser entries, writes create-only files, skips identical content on re-run,
reports conflicts without overwriting, and journals what it created. Rollback removes only
those journaled files and clears the cutover marker. The browser store is never modified, so it
remains the recovery source even after rollback.

**Propagation.** With identity-based paths, a rename does not change the file and a delete
targets exactly one file. A separate propagation subsystem would be a second owner. Existing
cross-reference cleanup stays where it already lives.

## 3. Architecture

### Canonical owner

`internal/lorebookstore` is the only writer of lorebook content.

- Reads are bounded: one directory level per root, regular `*.json` files only, 4 MiB per entry,
  5,000 entries per root, no recursion, no background scanning, no parallel hashing.
- Only `howl-lorebook-v1` envelopes are listed. Old F2-B export envelopes and unrelated JSON are
  ignored and counted, so nothing is read back through a stale export.
- Duplicate identities across roots are reported as conflicts; the first root wins
  deterministically and nothing is merged or rewritten.
- A disconnected root keeps its registration and its last-known entries from a bounded index
  cache, returned as unavailable and non-editable. A reachable but empty root is reported empty,
  not stale.
- Writes go through the registry: create requires a save root, replace requires the write role
  on the owning root, delete requires a writable root and re-derives the path from the identity.
  The F3-A-R1 `Lstat`, regular-file-only, unique-backup, and restore-reporting protections are
  inherited unchanged.
- The ownership transition is a marker plus a journal in the settings directory:
  `lorebook-cutover.json` and `lorebook-import-journal.json`.

### Bindings

`app-lorebooks-store.go` exposes `ListLorebooks`, `SaveLorebook`, `DeleteLorebook`,
`ImportLorebooks`, `RollbackLorebookImport`, `GetLorebookStoreStatus`,
`GetLorebookStoreDiagnostics`, and `EnsureLorebookLibrary`. Mutations return structured
results; the UI never infers success from the absence of an error.

### Frontend data path

`frontend/dist/js/lorebook-repository.js` is the single frontend path.

- Before import, `cutover` is false and the browser store remains the source; saves continue to
  the browser store exactly as before, so no file is created behind the user's back.
- After import, hydration rebuilds the browser cache from canonical entries, and every existing
  synchronous reader keeps working against that cache.
- Save and delete go canonical-first and refresh the cache only after success. Failures surface
  as an alert and leave the cache untouched.
- The lorebooks page shows the storage state, the entry count, and any identity conflicts, and
  offers `Import to Library` and `Undo Import`.
- Chat lorebook sync refreshes from the canonical store before prompt assembly.

### Ownership before and after

| Concern | Before F3-B | After F3-B |
|---|---|---|
| Lorebook content | `localStorage['howl-chat-lorebooks']` | Registered JSON files; browser store is a derived cache |
| Write path | Page writes the browser array | Store writes files; page writes through the repository |
| Identity | `Date.now()` string, no file identity | Preserved on import, random `lb_` for new entries |
| Unavailable root | Content vanished with the browser store | Last-known entries remain visible as unavailable |
| Recovery | Browser store only | Legacy browser store retained; import rollback removes only imported files |

## 4. Files changed

- `internal/lorebookstore/store.go` — new canonical store: bounded listing, index cache, save,
  delete, journaled import, rollback.
- `internal/lorebookstore/store_test.go` — new ownership, failure, conflict, rollback, and
  unavailable-root coverage.
- `app-lorebooks-store.go` — new bindings and diagnostics.
- `app-lorebooks-store_test.go` — binding-level cutover gating, save/list/delete, rollback.
- `app.go` — `lorebookStore` field and initialisation.
- `frontend/dist/js/lorebook-repository.js` — new repository and derived cache.
- `frontend/dist/js/lorebooks.js` — repository-backed load, save, delete; status and
  import/undo handlers.
- `frontend/dist/lorebooks.html` — repository script and store bar.
- `frontend/dist/css/lorebooks.css` — store bar styling.
- `frontend/dist/js/chat-stream.js` — canonical refresh before lorebook sync.
- `PROJECT_ISSUES_AND_EXPANSION.md` — F3-B decisions and status.
- This report.

## 5. Verification

```text
go build ./...
BUILD_EXIT=0

go test ./internal/lorebookstore -count=1 -v
15 passed, 1 skipped (symlink creation requires a privilege this host does not grant)
ok  howl-chat/internal/lorebookstore

go test . -run 'TestSaveLorebook|TestImportThenSaveListAndDelete|TestRollbackLorebookImport|TestLorebookBindings'
ok  howl-chat  1.720s

go test -race ./internal/lorebookstore -count=1
ok  howl-chat/internal/lorebookstore
RACE_EXIT=0

go test ./... -count=1
FULL_EXIT=0

node --check frontend/dist/js/lorebooks.js
node --check frontend/dist/js/lorebook-repository.js
node --check frontend/dist/js/chat-stream.js
JS_SYNTAX=0

git diff --check (changed files)
DIFF_CHECK=0
```

The symlink listing test skips on this host; registry-level symlink and directory rejection is
already covered by the accepted F3-A-R1 tests, which still pass in the full suite. Two earlier
test runs hung on a self-deadlock introduced during this stage (an index write taking the store
mutex from inside a locked section); the lock split fixed it, and the reruns above are the
post-fix results. No audio or unrelated package failure appeared in the full suite.

## 6. Regression evidence

- No world, scenario, character, chat, gallery, audio, model, or prompt-assembly source file was
  modified. `REG-1` behavior is therefore untouched by construction: attachment resolution,
  scopes, activation, `Attached to —` labels, and hot-swapping read the same data as before.
- Chat prompt assembly still consumes the same browser-array shape; only the refresh point was
  added before sync.
- The lorebooks page keeps its existing normalization, validation, and editor; only the storage
  target changed behind it.
- The existing F3-A and F3-A-R1 suites pass in the full run.

## 7. Limitations

- Before the explicit import, canonical files do not exist and page saves remain browser-only.
  That is deliberate: writing files while the browser store is still the reader would recreate
  the dual source F3-A removed.
- The cache is refreshed on page load, save, delete, and chat send. External edits made while a
  page is open are not watched; there is no background scanning by design.
- Import reassigns invalid or duplicate identities and reports them. Existing cross-references
  cannot distinguish entries that already collided, and no reference rewriting was performed.
- Rollback covers the most recent import only, and leaves adopted files in place after an interrupted-import recovery.
- Conflicts are reported, never merged or overwritten; resolving them is a user action.
- No runtime Wails UI, native folder picker, or packaged-application test was performed. The
  lorebooks page was not exercised in a browser.
- Worlds, scenarios, characters, chats, images, and videos keep their current ownership.
- `F3-A8` propagation is now inherent for lorebooks (identity-based delete/rename) but remains
  deferred for other categories.

## 8. Deferred

The next ownership slice remains worlds and scenarios (F3-C) in the recorded order. No deferred
`FUT-*` requirement was implemented. `FUT-LORE1` (automated lore intelligence) is untouched; this
stage changed storage and retrieval plumbing, not trigger or memory behavior.

## 9. Acceptance review additions (2026-11-10)

### Crash-window recovery

An import can be interrupted between three states: files written, journal written, marker
written. Previously a retry after the first state reported `nothing_to_do` and ownership was
never completed. The store now treats an identical canonical file found while the marker is
absent as evidence of an interrupted import: it verifies the content hash, journals the item as
`adopted`, writes the journal, and completes the marker.

Rollback removes only entries journaled as `created` or `reassigned`; `adopted` files are left in
place because this store did not create them in that journal. A retry after rollback can adopt
them again and re-complete ownership.

Covered by `TestInterruptedImportAfterFileWriteIsRecovered`.

### Duplicate identity resolution

When two roots contain the same identity, listing returns one entry and reports a conflict. The
winner is the first registered root in the registry's deterministic ordering, and repeated reads
return the same entry. An attachment that references the identity therefore resolves to that
winner consistently; the conflict remains visible so the duplication can be resolved
deliberately. Covered by
`TestDuplicateIdentitiesAcrossRootsAreReportedAndDeterministic`.

### Remaining acceptance work

Integration and runtime acceptance remain open: duplicate-identity behavior in a real
attachment, interruption during a real process, the Lorebooks page in the running application,
and `REG-1` behavior under the repository and cache arrangement. Source tests are not a
substitute for those checks.

## 10. Consumer resolution audit (source)

The goal was to confirm that no consumer resolves lorebook entries by scanning roots on its own.

| Consumer | Path | Result |
|---|---|---|
| Lorebooks page (`lorebooks.js`) | repository `hydrate`/`save`/`remove`, canonical-first | Single path |
| Chat prompt assembly (`chat-stream.js`) | hydrate, then send the cache to `SetLorebooks` | Canonical-backed |
| Context assembly (`context-data.js`, `context-selectors.js`) | read the cache by ID | Canonical-backed |
| World/scenario/character attachments (`world-state.js`, `scenario-state.js`, `character-state.js`, dropdowns, forms) | read the cache by ID, save reference arrays into their own records | Canonical-backed |
| Chat lorebook panel (`ui/lorebook-manager.js` in practice `lorebook-manager.js`) | reads the cache; chat-scoped toggle/remove writes chat data only | No catalog write reachable |
| Go backend | receives entry maps through `SetLorebooks`; chat persistence keeps a per-chat `lorebooks.json` snapshot | No root scan |

No consumer other than `internal/lorebookstore` opens lorebook entry files. `contentItemFile`
names an export copy used by the share/export flow, not the canonical store, and the store
ignores those envelopes.

Two findings are recorded rather than changed:

- `lorebook-manager.js`, loaded on the chat page, still contains `saveLorebook`,
  `deleteLorebook`, and `openLorebookCreationPopup`, which write `howl-chat-lorebooks`
  directly. They have no callers in the current markup or scripts, so no live path bypasses the
  repository, but a future caller would write a cache-only change that the next hydration
  discards.
- Per-chat `lorebooks.json` snapshots keep chat-scoped copies of entries. Prompt assembly
  refreshes from canonical, but chat restoration enriches from those copies. This is the
  documented chat-scope duplicate, not a second root reader; reconciling it belongs to the chat
  stage, not F3-B.

Neither finding changes the ownership architecture. Both are listed in the acceptance gate.

## 11. Runtime acceptance through the UI harness (2026-11-10)

The project's `uiharness` build served the real `frontend/dist` to a browser and proxied every
`window.go.main.App.*` call by reflection to the real backend. The harness ran from a throwaway
folder whose own `wails.json` pinned `ProjectRoot` there, so config, library registrations,
journals, and entries were isolated from the user's data.

| Check | Result |
|---|---|
| Pre-import state | Browser-storage message, import enabled, undo disabled, no files |
| Create, three scopes | Three entries saved through the repository; no dialogs, no console errors |
| Import without a save root | `Import failed. no lorebook save destination is configured`, nothing written |
| Import | Three files with original identities, marker, journal (`created`), index |
| Rename | Same file name, updated title, no new file |
| Delete | Exactly one file removed; page and cache followed |
| Undo import | Files removed, marker cleared, journal `rolled_back`, browser store kept |
| Re-import | Two files recreated, cutover restored |
| Restart + cache wipe | Entries rebuilt from canonical files after process restart |
| Interruption recovery | Marker removed, reload, retry: journal `adopted`, mtimes unchanged, marker restored |
| Duplicate identity | Two roots, one conflict reported, deterministic winner displayed |
| REG-1 attachment | World saved with lorebook reference to the canonical identity; survived reload |
| Console / network | Only the expected no-destination 500 and the reset from the deliberate restart |

Limitations: the harness replaces Wails IPC with HTTP reflection, so it is the real backend and
real page code but not the native WebView2 window; chat prompt assembly needs a model and was not
exercised; character and scenario attachment forms, activation, and hot-swapping were not
exercised; the Windows symlink case remains skipped.

### Assessment (2026-11-10)

Harness integration acceptance passed. Runtime acceptance is still partial, and the remaining
work is bounded:

- **Stale Import button (low severity).** Seen only when the marker was removed under a live
  page. A real restart hydrates the marker-less state and Import is enabled, so recovery does
  not normally need a reload; the store bar simply does not self-refresh on external marker
  changes during a session.
- **`REG-1` partial.** World attachment verified end to end; scenario and character forms,
  activation, private/public lore, and hot-swapping are not.
- **Native runtime outstanding.** HTTP-reflected bindings are not a packaged WebView2 session.

Remaining: native window smoke test, character and scenario attachments, chat prompt
synchronization with a model, creative hot-swapping, and the Windows symlink case. F3-B remains
open for final acceptance; no architecture change follows from this record.

## 12. Final acceptance pass (2026-11-10)

Verdict: **ACCEPTED WITH DOCUMENTED LIMITATIONS**. F3-C is not started.

| Check | Status | Evidence |
|---|---|---|
| Native desktop build | Passed | `go build -tags desktop` exit 0 |
| Native WebView2 runtime | Blocked | WebView2 user-data folder cannot be isolated from the real profile; not run |
| Scenario attachment | Passed | `Black River Vigil` → `1794283090983` (Tide Oath), survived reload |
| Character attachment | Passed | `Rill Final` → `howl.lorebooks` contains `1794283090983` |
| Reusable attachment | Passed | Same identity attached to scenario and character |
| Legacy catalog writes | Passed | No external caller of the three dormant methods |
| Symlink listing | Blocked | Windows privilege error from `mklink` |
| Chat prompt synchronization | Unverified | No prompt-preview binding; no inference server |
| Activation / hot-swapping | Unverified | No active chat workspace exercised |
| Private character lore injection | Unverified | Field stored canonically; injection not exercised |
| Disconnected-root display | Unverified | Not exercised |

Pre-existing character-form observation: stale selected-lorebook chips can remain after closing
and reopening the form, because the chip list is not re-rendered on open. Saved data reflects the
actual selection. This is a character-surface display issue, outside F3-B ownership.Remaining limitations: native WebView2 launch (data isolation unavailable), chat prompt
  delivery, activation/hot-swapping, private-lore injection, disconnected-root display, and the
  Windows symlink test. None contradict the verified ownership, persistence, recovery, identity,
  or attachment-resolution behaviour.

## 13. Addendum — native gate cleared (2026-11-10)

Written after §12; it corrects status, not history. The isolation work it depended on landed, and
the W27 native scenario ran in the real Wails window against a marked isolated root:

- **Native WebView2 runtime:** no longer blocked. The run used `HOWL_WEBVIEW_DATA_DIR` with the
  required marker; the real profile and config were not used.
- **Private character lore injection:** exercised natively for two stable character IDs, including
  the cross-owner conflict defect found and fixed during that run (§27 of the isolated-runtime
  report; W27a in the ledger).
- **Chat prompt delivery:** observed at the inference boundary (character composer → `inferenceClient`),
  not with a loaded model. Generation quality remains unverified.
- **Still open:** activation/hot-swapping in a live creative workflow, disconnected-root display in
  the native window, and the Windows symlink listing test.

Verdict unchanged: **ACCEPTED WITH DOCUMENTED LIMITATIONS**, now with the native gate satisfied.
F3-C is not started.
