# TASK-03 — D-3: one membership authority, demonstrated in both hosts (`TASK-03-D3-MEMBERSHIP-AUTHORITY-2026-11-22`)

**Status:** CURRENT — implemented, tests green, Missing and Orphaned demonstrated in the packaged app and the browser host. Uncommitted.
**Scope:** conversation membership and lifecycle; the integrity defects that share that boundary; the acceptance evidence the previous pass was missing.

## D3-01 — The decision: one lifecycle log, directories are the payload

A new authority in `internal/backend/chat/membership.go` owns **which conversations exist and what happened to them**. It is an append-only log, `lifecycle.ndjson`, one complete JSON record per line, written through the same `appendLine` helper as every other write in the package.

```
{"event":"created","id":"…","name":"…","at":"…"}
{"event":"deleted","id":"…","at":"…","directory":"…"}
```

Reading it is a fold: `created` adds a member, `deleted` removes it, and **deletion wins permanently** — a later `created` for a deleted identifier is ignored, so no file, restore or cached client can put a deleted conversation back.

Everything else is derived, never stored twice:

| Question | Answer | Written by |
|---|---|---|
| Does this conversation belong to the collection? | fold of the log | create/delete lifecycle only |
| Is it there? | `Reconcile` stats its directory | nobody (read-only) |
| Where is its folder? | `conversationDir(id)` | same constructor every mutation uses |
| May I write into it? | `Present(id)` — claimed, not deleted, directory exists | nobody (read-only) |

## D3-02 — `deleted.ndjson`: absorbed, then retired

Evaluated three ways and decided: **absorbed**. `deleted.ndjson` was the only record of intent on existing installs, so on the first load with no `lifecycle.ndjson` its records are migrated once into the log as `deleted` events. After that the file is never written or read again; it is left on disk as user data. It is not a second authority, and it is not kept alive because it already existed.

The same first load **bootstraps** the existing directories: a folder that identifies itself with readable `meta.json` and a valid identifier becomes a `created` record. Without this every existing install would have reported its whole collection as orphaned. Bootstrap runs exactly once, only when the log is absent; a directory that appears later is never claimed, it is reported as orphaned.

## D3-03 — What the authority reports

`GetChatList` is now a read-only reconciliation: **present** (claimed + folder there), **missing** (claimed, folder gone), **orphaned** (folder nobody claims, including the leftover of a deletion whose files survived). Each summary carries its state and its directory.

Nothing is repaired, adopted or removed by listing. The three states are computed, not stored.

## D3-04 — The gates that make resurrection impossible

| Path | Rule |
|---|---|
| `NewChat` / `SaveChatAs` | claim the directory, then record `created`, then write |
| `SaveChat` | a **known** conversation whose folder is gone is refused (`refused_unbound`); an unrecorded id is recorded as a member before the write, so it can never become an orphan |
| `AppendHistoryEntry` / `RewriteHistory` | refuse a missing directory |
| `CommitChat` | refuses any id the authority does not consider **present** |
| `App.SaveChatData` | refuses a frontend payload whose id is not present — a cached client cannot mint a conversation |
| `DeleteChat` | record `deleted` **before** removal; recycle; surface leftovers |

## D3-05 — The shadow writer that clicking found

While demonstrating D-3, the browser host created a conversation in a fresh isolated root that the authority had never recorded. Two causes, both in the frontend's cached session (`HOWLStorage` / `localStorage`):

1. the validator accepted any id present in the list — including missing and orphaned ones;
2. when the backend bridge was not ready, validation was skipped entirely and the cache was pushed anyway.

Fixed fail-closed: a cached session is restored only when the authority reports it **present**, and a cache without a conversation id is discarded instead of synced. This is the shadow-tree ownership defect: the cache is now a view of the authority, never a writer of it.

## D3-06 — The other integrity items

| Item | State |
|---|---|
| `SaveChat` recreating a missing conversation | fixed (D3-04) |
| Direct-save fixture migration | not needed under the rule chosen: fixtures that save unrecorded ids are recorded first, and the rule that matters (no recreation of a known conversation) has its own test. The stricter "only the claim path may create" was rejected because it forces every fixture through a different API without changing any user-visible outcome |
| `LastChatID` cross-process race | fixed: `saveConfig` holds `config.lock` through `durable.LockFile` across the read-modify-write |
| Shadow-tree writer ownership | fixed for the startup restore paths (D3-05); one more consumer was removed by the `Present` gate on `SaveChatData` |
| Autosave cancellation / shutdown joining | startup background work (`rescan`, `model discovery`) now runs through `startBackground` and is joined by `waitForBackground(10s)` before the final save. Frontend autosave timers are outside Go's lifetime and are not joined |
| Pointer race (row 571) | not addressed; the config lock covers the same race for `LastChatID`, which is the field row 571 describes |

## D3-07 — Demonstrated through the application, in both hosts

Packaged app (`-tags "desktop,production"`), isolated root, real mouse events:

1. File ▸ New Chat created a conversation — the log itself recorded `{"event":"created","id":"2026-11-22T19-52-05Z",…}`.
2. The folder was removed outside the app, and an unclaimed folder with its own `meta.json` was placed in `chats/`.
3. Restart, File ▸ Load — `build/verify/25-membership-states.png` shows **“Chat 11/22/2026, 1:52:05 … — Missing — its folder is gone”** and **“Folder nobody claimed — Orphaned folder — not in the collection”**, neither selectable.

Browser host (`-tags uiharness`, same backend over `/__rpc`), separate isolated root seeded with one claimed-but-absent id and one unclaimed folder: the Load dialog renders exactly the same two cards with the same states, both non-selectable, `clickable:false`, and **nothing is written** — the chats directory and `lifecycle.ndjson` are byte-identical before and after a full page load.

Cancel parity, the item left open by the previous pass: File ▸ Delete ▸ Cancel in the browser host leaves `DeleteChat` uncalled (`grep -c` = 0 in the RPC log), the `deleted.ndjson` line count unchanged, and `meta.json`/`history.ndjson` with unchanged hashes and mtimes.

## D3-08 — Tests

New, all in `internal/backend/chat`:

- present / missing / orphaned reconciliation, with the orphan's label read from its folder but never trusted as identity;
- deletion survives a reload, refuses recreation, and a leftover stays orphaned;
- a torn final log line is ignored while a complete record is kept;
- legacy `deleted.ndjson` is absorbed once and never duplicated;
- bootstrap claims existing directories once, and a later directory stays orphaned;
- `AppendHistoryEntry` / `RewriteHistory` / `SaveChat` (known conversation) refuse and create nothing when the folder is gone.

## D3-09 — Checks

```
go build ./...                        exit 0
go test . ./internal/...              all ok
go vet ./internal/backend/chat/ .     clean
gofmt -l …                            empty
node --check chat.js chat-session.js sidebar-session.js delete-confirm.js   ok
go build -tags "desktop,production"   exit 0
go build -tags uiharness              exit 0
```

## D3-10 — Remaining

- Selecting a missing or orphaned card is impossible, but there is no explicit "this conversation is missing" recovery action beyond showing its folder; what the user should be offered is a product decision.
- Frontend autosave timers are not joined at shutdown.
- Row 571 pointer race beyond `LastChatID`; the two hosts still differ where native surfaces are required (file pickers, message boxes).
