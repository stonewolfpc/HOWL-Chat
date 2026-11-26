# TASK-03 — Deletion architecture and S-1 correction pass

| Field | Value |
|---|---|
| Purpose | Record the deletion rework, the S-1 corrections, and the ownership decisions behind them |
| Status | **Implemented and tested in this workspace. Uncommitted. User-visible behaviour not demonstrated on a running build** |
| Basis | Rows 586–588; the directive of 2026-11-22; code and tests in this workspace |
| Recovery token | `TASK-03-DELETION-AND-S1-CORRECTION-2026-11-22` |

---

## §1 Windows-native deletion

**DE-01 — What the platform offers, and what it refuses to promise.** `SHFileOperationW` is superseded by `IFileOperation` (Vista+), which reports accurate HRESULTs; `SHFileOperation` returns codes that can be misleading. Both recycle only **"if possible"**: `FOF_ALLOWUNDO` preserves undo information *if possible*, and a volume without a Recycle Bin, or a full one, still ends in permanent deletion with no error. There is therefore **no Windows API that guarantees recycling**. **[Source-derived]** Microsoft documentation for `SHFILEOPSTRUCTW`, `IFileOperation::SetOperationFlags`, and the `SHFileOperationW` remarks.

**DE-02 — Decision: recycle through the shell, and hold no permanent-delete path.** `recycleTree` (`recycle_windows.go`) calls `SHFileOperationW` with `FO_DELETE | FOF_ALLOWUNDO | FOF_NOCONFIRMATION | FOF_SILENT | FOF_NOERRORUI`, on a COM-initialised, OS-locked thread, with the double-NUL-terminated source the API requires. If the shell refuses or aborts, the error is returned and the files stay. **This program never falls back to deleting permanently itself.** What Windows does inside a successful call is the platform's decision, and the deletion notice tells the user where the files were so they can verify.

**DE-03 — Why not `IFileOperation` now.** It needs COM interface interop (vtable calls through `shobjidl`) for a single delete path that `SHFileOperationW` performs correctly on supported Windows. Choosing it would add the most error-prone code in the package in exchange for better HRESULT text. Revisit if per-file diagnostics are ever needed; the call site is one function.

**DE-04 — The obsolete path is removed, not wrapped.** The direct `os.RemoveAll` deletion is gone from Windows builds: `removeTree` was deleted from `writers.go`, and the only remaining `os.RemoveAll` is in `recycle_other.go`, compiled only off Windows, so a non-Windows test run still builds. There is no runtime switch between the two.

**DE-05 — Deletion semantics, unchanged by the swap.** Confirmation removes the conversation from the visible collection; the filesystem operation is attempted and not verified; a refusal leaves files and does not restore the item or fail the call; the dialog names the directory and links to it. `DeleteChat` reports an incomplete removal to the console and returns success, because the item is gone from the user's view either way.

---

## §2 The tombstone: kept, with a single owner

**DE-06 — Necessity.** A blocked delete leaves the directory in `chats/`. Without a record of intent, the next startup scan lists the conversation again — the user's decision silently reverted. The Recycle Bin does not change this: a refused recycle leaves exactly that leftover.

**DE-07 — What it is, and is not.** `deleted.ndjson` is an append-only log of *intent*: one line per deletion, written **before** any removal. It is not proof of destruction, not a verification subsystem, and not a second membership authority. Its single owner is `PersistenceManager`, which writes it inside `DeleteChat` and reads it in `List`/`Load`/`Save`/`RewriteHistory`/`Append`.

**DE-08 — Lifecycle and growth.** Appended on each deletion, never rewritten or compacted in this slice; a duplicate line is harmless because it is read as a set. Growth is one short line per deleted conversation, and the file is torn-tail tolerant. Compaction belongs with the index, which will supersede the read path; compacting now would create a second writer for the same fact.

**DE-09 — The index must consume it, not duplicate it.** When membership becomes authoritative, the tombstone remains the only record of explicit deletion: the index records membership, the tombstone records the user's decision to remove it. Neither may resurrect the other, and reconciliation may never report a deleted identifier as `Missing` or as damage.

---

## §3 S-1 corrections found by review

**S1C-01 — The final-save window is closed.** `performFinalSave` used to save the snapshot, reload, and then rewrite both representations — two writes of one history, with a window in which the snapshot was newer than the log (row 570). It is now one load and one `CommitChat`: companions, then the log, then the snapshot, under one lock hold. An interruption can only leave the log newer, which is the state that loses nothing, because the loader prefers the log. **[Observed]** in code and covered by the existing shutdown tests.

**S1C-02 — Creation no longer adopts an existing conversation.** `NewChat` and `SaveChatAs` claim their directory with an exclusive `Mkdir` and retry with a suffixed identifier when the same-second identifier is taken. The row 580 data loss — a new chat silently inheriting an existing conversation's metadata, companions and `memories.json` — is no longer reachable through either path. **[Observed]** and tested by creating three conversations in one second and loading each back under its own name.

**S1C-03 — A refused append is no longer silent.** The lock introduced refusals where previously every append returned. `app-prompt.go` discarded that result for the assistant reply; it now logs at High and emits `chat:unsaved`. Four `_ = a.performSave()` sites now log the same way.

**S1C-04 — What review did not change, and why.** `SaveChat`, `RewriteHistory` and `AppendHistoryEntry` still create their directory with `MkdirAll`, so a write can still materialise a conversation that does not exist. That is a real remaining defect (the creation side of `V3-72a`), but changing it invalidates callers that build a `ChatData` and save it directly, including tests that construct their fixture that way; it needs its own pass with those call sites migrated. Creation itself no longer adopts, which is the data-loss half.

---

## §4 Ownership, as it now stands

| State | Authority | Competing owners still present |
|---|---|---|
| Conversation directory path | `conversationDir` (single constructor, containment-checked) | none inside the package |
| Conversation content (meta, companions, log) | `PersistenceManager`, one lock per conversation, writes only in `writers.go` | `service.go` shadow tree writes a different root (build-excluded from production); `WipeAllData` at root level |
| History bytes | `history.ndjson` is read first; `history.json` is the fallback and is written last | `SaveChat` writes the snapshot from app state; a turn captured but never appended exists only in memory (row 574's third domain) |
| Deletion intent | `deleted.ndjson`, written by `DeleteChat` only | none; the index will consume it |
| Active-chat pointer | `AppConfig.LastChatID` in `config.json`, atomic write, last-writer-wins | each running instance; row 571's defect, unchanged |
| Membership | **nobody** — derived from directory contents, so nothing can report a missing or orphaned conversation | this is D-3 |

---

## §5 Run in this pass

```
go build ./...                                        exit 0
go test . ./internal/backend/chat/ ./internal/durable/ ok
go vet ./internal/backend/chat ./internal/durable .    clean
gofmt -l                                               clean
```

New: `TestRecycleTreeMovesDisposableFixtureToRecycleBin` (executes the real shell call on a disposable fixture and observes the directory removed), `TestRecycleTreeIsIdempotentForAMissingPath`, `TestConcurrentCreationsClaimDistinctDirectories`. All twelve boundary tests pass.

---

## §6 Not demonstrated

**DE-10 — Stated plainly.** The application was not launched, so the dialog copy, the clickable folder link and Explorer opening are **syntax-checked only**. Whether the Recycle Bin received the fixture, rather than Windows destroying it, was not verified — the test asserts only that the call succeeded and the directory is gone. No second process was run. No measurement of the added containment resolution and tombstone read on the save/append path.
