# TASK-03 — Front-end deletion verified by clicking the real executable (`TASK-03-FRONT-END-DELETION-VERIFIED-2026-11-22`)

**Status:** CURRENT — implementation and front-of-screen evidence. Uncommitted.
**Scope:** the user-visible delete path only: menu item, confirmation dialog, folder link, removal, leftover handling across restart.

## FE-01 — The first real defect found by clicking

The Delete menu item **could never run**. `deleteActiveChatSession` guarded on `window.currentChatId` and passed it to `GetChatDirectory` / `DeleteChat`, but nothing in the shipped front end ever assigned that global: it was initialised to `null` at `frontend/dist/js/chat.js:27` and afterwards only cleared. The conversation identifier lives on `window.dataSync.currentChatId`, set by the load path (`frontend/dist/js/sidebar-session.js:207`, `data-sync.js:67/285/327/367`).

Evidence before the fix, from the real executable:

| Observation | Value |
|---|---|
| File ▸ Delete clicked (menu item highlighted) | no dialog rendered |
| Store tombstone `deleted.ndjson` | absent |
| Conversation directory | still present, unchanged |
| App log | no deletion attempt, no error |

That is the ghosting class this task exists to remove: the UI accepted the click and did nothing, and every backend test passed throughout, because the backend was never reached.

## FE-02 — One front-end owner for the active conversation

`chat.js` now reads the active conversation from its single owner and nothing else:

```js
function activeChatId() {
    if (window.dataSync && window.dataSync.currentChatId) {
        return window.dataSync.currentChatId;
    }
    return null;
}
```

`window.currentChatId` is gone; the guard, `GetChatDirectory`, `OpenChatDirectory` and `DeleteChat` all use the same value, and a successful delete clears `window.dataSync.currentChatId`. There is no second copy to fall out of step.

## FE-03 — A link that cannot open the folder now says so

`onOpenDirectory` previously reported failure to the console only, which is invisible in a shipped WebView2 window. It now surfaces `Could not open the folder: <error>` through the normal notification path. This is also what made the click decidable during verification.

## FE-04 — What was clicked, in order, in `build/bin/howl-chat-s1.exe`

Isolated store: `%TEMP%\howl-s1-verify` (`HOWL_WEBVIEW_DATA_DIR` + `.howl-isolated-root`). Build: `go build -tags "desktop,production"`. All input was real mouse events; all state was read back from disk.

| Step | Screenshot | Observed |
|---|---|---|
| Load fixture `2026-11-22T00-00-00Z` ("Deletion Fixture") | `09-loaded.png` | both turns rendered |
| File ▸ Delete | `13-dialog.png` | dialog: *Delete current chat?*, best-effort explanation, directory `…\chats\2026-11-22T00-00-00Z` as a link, *Deletion is permanent and cannot be undone.* |
| Dialog ▸ Delete | `14-after-delete.png` | conversation cleared from the window; list empty |
| Disk after confirm | — | `deleted.ndjson` = `{"id":"2026-11-22T00-00-00Z","at":"2026-11-22T19:21:14Z","directory":"…\\chats\\2026-11-22T00-00-00Z"}`; directory gone from `chats/`; `config.json` reduced to no active chat |
| Recycle Bin | — | `$RA2V2NR` present in `C:\$Recycle.Bin\S-1-5-21-…`, named `2026-11-22T00-00-00Z` — the item was **recycled, not destroyed** |
| Restore the directory by hand, restart the app, File ▸ Load | `15-leftover-hidden.png` | only the surviving chat is listed; the deleted one does not come back |
| File ▸ Delete on the surviving chat, click the directory link | `20-link.png` | `[OpenChatDirectory] opened …\chats\2026-11-22T19-19-37Z` in the app log; one Explorer window present, then closed |
| Dialog ▸ Cancel | `21-cancelled.png` | dialog closes, conversation and directory untouched |

The leftover case was produced by restoring the deleted directory with a `meta.json` and a `history.ndjson` by hand, because the shell refused nothing on this volume. That tests the stronger claim: intent, not filesystem success, decides what is listed.

## FE-05 — Verified limitations

- **The folder link has a two-line hit box.** A click aimed at the first line's pixel row missed it while the second line worked; this is a coordinate artefact of driving by screenshot, but it also means the link is only as large as its glyphs. Widening it is a UI decision, not done here.
- **UI Automation cannot see into WebView2** (`TREE` shows only `Window/Pane/Pane`), so all driving is coordinate-based with screenshot confirmation. `WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--force-renderer-accessibility` did not change that.
- **Only the packaged WebView2 host was exercised.** No browser/dev-server path was driven for this flow, so "browser and executable agree" is not claimed; the only supported host for these dialogs is the Wails window that was tested.
- **Copy conflict, owner decision:** the dialog says *Deletion is permanent and cannot be undone* while the code recycles the folder. The item is gone from the application either way, but a user can restore it from the Recycle Bin. Which sentence is true for the product is a product decision, not an engineering one.
- No second process, no locked-file refusal exercised, no timing measurement.

## FE-06 — Checks

```
go build ./...                                          exit 0
go test . ./internal/backend/chat/ ./internal/durable/  ok (8.8s / cached / cached)
go vet ./internal/backend/chat/ .                       clean
gofmt -l app-chat.go internal/backend/chat/*.go         empty
node --check frontend/dist/js/chat.js                   ok
node --check frontend/dist/js/delete-confirm.js         ok
```

## FE-07 — Owners after this pass

| State | Owner |
|---|---|
| Active conversation (front end) | `window.dataSync.currentChatId` — the only reader is `activeChatId()` |
| Deletion intent | `PersistenceManager` → `deleted.ndjson` (append-only, written before removal) |
| Directory removal | `recycleTree` → shell `SHFileOperationW` + `FOF_ALLOWUNDO`; no permanent-delete path in this program |
| Listing | `GetChatList` / `LoadChat`, which refuse tombstoned identifiers |

Competing owner found and removed: the dead `window.currentChatId` global, which made the delete action unreachable while appearing wired.
