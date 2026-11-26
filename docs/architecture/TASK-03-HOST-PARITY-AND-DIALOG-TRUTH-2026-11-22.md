# TASK-03 — Host parity, dialog truth, and the membership gap (`TASK-03-HOST-PARITY-AND-DIALOG-TRUTH-2026-11-22`)

**Status:** CURRENT — dialog and link corrections verified in both supported hosts; one integrity repair landed; D-3 not implemented.
**Scope:** the three corrections requested after the first click-verified deletion pass, plus the permanent acceptance gate.

## HP-01 — The two supported hosts, named

`AGENTS.md` now defines them, because "the browser" had been ambiguous:

| Host | Build | What it is |
|---|---|---|
| Packaged desktop app | `-tags "desktop,production"` | Wails window, WebView2, embedded `frontend/dist` |
| Browser dev host | `-tags uiharness` | `harness_main.go`: serves the same `frontend/dist` on `http://127.0.0.1:8799`, proxies every `window.go.main.App.*` call by reflection to the same `App` methods, streams Wails events over SSE |

There is no third host. A static file server, a mocked bridge, or the `native_test_seam` build is not accepted as evidence.

## HP-02 — Dialog copy corrected

The contradiction is gone. The warning sentence is now a caller-supplied `warning` option on `DeleteConfirm.show`, defaulting to the old irreversible-deletion sentence for the five other dialogs whose deletions are not recycled.

> **Delete current chat?**
> This chat leaves HOWL immediately. Its folder is moved to the Windows Recycle Bin, where you can normally restore it. World, scenario, and character data are preserved.
>
> The chat folder is
> `%TEMP%\…\chats\2026-11-22T19-19-37Z`
> If Windows blocks the move, the files stay there. Open the folder to check.
>
> ⚠ The chat is removed from HOWL now. Windows keeps the folder in the Recycle Bin until you empty it.

## HP-03 — The whole path is the hit target

The link was an inline anchor, so only its glyphs were clickable across a wrapped two-line path. It is now a `display:block` block with padding and its own background: measured in the browser host at **354 × 52 px**, the full width of the dialog body. A click at the far right end of the block, past every glyph, opened the folder in the packaged app (`[OpenChatDirectory] opened …2026-11-22T19-19-37Z`).

## HP-04 — Same action, same result, demonstrated in both hosts

Isolated roots: `%TEMP%\howl-s1-verify` (desktop), `%TEMP%\howl-harness-verify` (browser). Same fixture shape, same clicks, results read back from disk.

| Step | Packaged app (WebView2) | Browser dev host |
|---|---|---|
| Load fixture through the UI | both turns render (`09-loaded.png`) | both turns render (`Chat session loaded: Harness Fixture`) |
| File ▸ Delete | dialog opens with the copy above | dialog opens with byte-identical text |
| Path block geometry | full-width block (`22-copy.png`, `23-final-copy.png`) | `display:block`, 354 × 52 |
| Click the path, far right edge | `[OpenChatDirectory] opened …` ; 1 Explorer window | `[OpenChatDirectory] opened …` ; `[rpc] OpenChatDirectory ok 8ms` ; Explorer at the same `file:///` path |
| Dialog ▸ Delete | folder gone, `deleted.ndjson` line written, `$RA2V2NR` recycled | folder gone, `deleted.ndjson` line written, `$RLE5UPE` recycled |
| UI after delete | conversation cleared | 0 messages, fixture text absent |
| Dialog ▸ Cancel | chat and folder untouched (`21-cancelled.png`) | — (cancel verified in the app; the browser dialog shares the same component) |

Every backend call in the browser host was the real method (`[rpc] DeleteChat ok 197ms`), not a mock, and the same code moved the second fixture to the Recycle Bin.

## HP-05 — One integrity repair: a write can no longer invent a conversation

`RewriteHistory` and `AppendHistoryEntry` used `ensureDir`, so a write mentioning an unknown identifier materialised a conversation directory. Neither is the creation route, so both now call `requireExistingDir` and return `Refusal{RefusedUnbound, "conversation directory is missing"}`. Two tests create no directory, call the mutator, assert the refusal code, and assert the directory still does not exist.

`SaveChat` still creates its directory. It is the one mutator that also serves first save of a freshly created session; removing that needs every direct-save caller and test fixture migrated onto the claim path, which is its own pass (the creation side of `V3-72a`, unchanged from `HEKI-CP-2026-11-22-77`).

## HP-06 — D-3: designed, deliberately not implemented

The deletion intent stays the transitional authority until D-3 exists; no second index was added, as instructed.

What D-3 has to answer before code: **what creates an expectation that a conversation should exist?** Listing already derives presence from directories and subtracts `deleted.ndjson`, which means a conversation whose directory vanished is simply absent — unreportable as `Missing`. Reporting `Missing` therefore needs a record of expectation. The two candidate sources are:

1. **Persisted membership** (a manifest of known conversations) — makes `Missing` reportable, but it is a second authority that can disagree with the directory set and with `deleted.ndjson`, which is the failure this project keeps removing.
2. **References only** (the active session, `config.json`'s `lastChatId`, autosave targets) — expectation is derived from state that already exists for its own reasons; `Missing` means "a live reference points at a directory that is not there", and nothing new can disagree with anything.

Option 2 is the one consistent with "one authority per state"; option 1 needs an explicit owner decision, because it re-introduces a mirror. Implementation belongs in the next slice, with the same click-through requirement.

## HP-07 — Checks this pass

```
go build ./...                                            exit 0
go test ./internal/backend/chat/                          ok (2.1s, includes the two new refusal tests)
go test .                                                 ok (9.9s)
go vet ./internal/backend/chat/                           clean
gofmt -l internal/backend/chat/*.go                       empty
node --check frontend/dist/js/chat.js                     ok
node --check frontend/dist/js/delete-confirm.js           ok
go build -tags "desktop,production"                       exit 0
go build -tags uiharness                                  exit 0
```

## HP-08 — Remaining, with evidence of absence

- `SaveChat` still creates a missing directory (see HP-05).
- Pointer race (row 571), shadow-tree writer ownership, joining autosave goroutines at shutdown: untouched.
- D-3 membership: not implemented (HP-06).
- The browser host cannot exercise native file pickers or message boxes; that difference is recorded, not papered over.
- Screenshots: `build/verify/22-copy.png`, `23-final-copy.png` (packaged app). The browser host was inspected through the preview tooling and its RPC log, not captured as an image.
