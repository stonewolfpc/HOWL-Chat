# TASK-03 — S-1: implementation and audit (mutation boundary, best-effort deletion)

| Field | Value |
|---|---|
| Purpose | Record what S-1 now does, the audit that followed it, and what remains unproven |
| Status | **Implemented and tested in this workspace. Not committed. No deployment. D-2 remains open until its acceptance criteria are demonstrated on a real build** |
| Basis | Rows 586–587; the code as changed; test runs listed in §6 |
| Scope | S-1 as authorized: constructor, one mutation entry, per-conversation lock, tail rule, build-failing analysis, App-layer surfacing; plus the owner's deletion clarification |
| Recovery token | `TASK-03-S1-IMPLEMENTATION-AND-AUDIT-2026-11-22` |

---

## §1 The two changes, in one line each

**S1-01 — The conversation write boundary now exists.** Every conversation mutation goes through `withChatLock`, every conversation path through `conversationDir`, and every filesystem write through four helpers in `writers.go`. A package-level analysis test fails the build if a fifth write appears. **[Observed]**

**S1-02 — Deletion is now the operation the owner specified.** Confirmation removes the conversation from the visible collection; physical removal is attempted and not verified; a blocked delete neither restores the item nor fails the call; the confirmation names the directory and links to it. **[Observed]**

---

## §2 What was built

| Item | Where | Property |
|---|---|---|
| `conversationDir(id)` / `chatsRoot()` | `boundary.go` | One path constructor; refuses a non-component identifier; containment decided on physical paths via `paths.CheckContained`, so a junction planted inside `chats/` is refused. Replaces nine inline joins. |
| `withChatLock(op, id, fn)` | `boundary.go` | Per-conversation exclusive lock, sibling namespace `<root>/.chat-locks`, refusal (not waiting) on contention; same-process double entry is diagnosed as `refused_reentrant`, never reported as contention. |
| `Refusal` / `RefusalCode` | `boundary.go` | Typed "nothing was written" result: `refused_busy`, `refused_reentrant`, `refused_unbound`, `refused_containment`, `refused_not_written`. |
| `ensureDir`, `writeFileAtomic`, `removeTree`, `appendLine` | `writers.go` | The only filesystem writes in the package. `appendLine` terminates an unterminated tail before appending; it never concatenates onto it. |
| `LockFile` | `durable/lock_api.go` | Reuses the store's proven Windows byte-range lock rather than inventing a second locking mechanism. |
| Deletion intent | `deleted.ndjson` at the root | Append-only, one record per deletion, written **before** removal. A leftover directory and a late autosave both refuse; the conversation cannot reappear. |
| Analysis test | `boundary_analysis_test.go` | Fails the build if any function outside the four helpers references `os.WriteFile/OpenFile/Create/Truncate/Rename/MkdirAll/Mkdir/Remove/RemoveAll/Chmod/Chtimes`. Includes a self-test proving it rejects a synthetic bypass. |
| Deletion notice | `delete-confirm.js`, `chat.js`, `App.GetChatDirectory`, `App.OpenChatDirectory` | The confirmation states that Windows sometimes blocks deletion, shows the exact directory as a link, and opens it on click. |

**S1-03 — One write of a finished copy.** `SaveChatAs` previously published the destination twice, the first time empty (row 572's observation). It now builds the copy and writes once. **[Observed]**

---

## §3 The audit the owner asked for

**S1-04 — Coverage of every mutation path, not just the fourteen writes.** The eight exported mutators — `NewChat`, `SaveChat`, `SaveChatAs`, `DeleteChat`, `RewriteHistory`, `AppendHistoryEntry`, `RenameChat`, `SetLastChatID` — were walked one by one. Seven acquire the conversation lock and write only through the helpers. `SetLastChatID` writes the root `config.json` and takes no conversation lock: it is the active-chat pointer, not conversation content, and its cross-process last-writer-wins behaviour is row 571's separately owned defect. Calling that covered would be false. **[Observed]**

**S1-05 — Concurrent operations.** Two processes contending on one conversation: the second is refused with `refused_busy`, verified by taking the lock directly and attempting a save. Nesting inside one operation: refused with `refused_reentrant`, verified. Autosave goroutines now contend on the same lock instead of interleaving, which is the point. **[Observed]**

**S1-06 — Failure handling, and the regression it exposed.** The lock introduced a new outcome: an append can now be refused where previously it always returned. `app-prompt.go` discarded that result with `_ = a.AppendChatMessage(...)`, so the assistant reply could be displayed, absent from the log, and silent — the exact class row 574 recorded. It now logs at High and emits `chat:unsaved` with the reason. Four `_ = a.performSave()` sites were the same shape and now log. **[Observed]**

**S1-07 — A race the first draft had, found in audit and removed.** Deletion intent was a rewritten JSON document, so two concurrent deletions could lose one another's record. It is now an append-only log (one line per deletion), which cannot lose an entry and tolerates a torn tail by ignoring it. **[Observed]**

---

## §4 Deliberately not covered

**S1-08 — What this does not claim.** Creation identity collision (two creations in one second adopting a directory) is unchanged: `V3-72a` owns it. The shadow-tree writer in `service.go` uses the helpers but writes a different root with no conversation lock, and is build-excluded from production (row 573). `WipeAllData` and its rename staging live outside the boundary and are unchanged. The active-chat pointer's cross-process race is unchanged. Repair, the final-save window, retry identity and memories are untouched. No marker, generation or settlement field was added.

**S1-09 — What "no resurrection" now means, and does not.** An explicitly deleted conversation cannot be listed, loaded, saved into or appended to, even if its directory survives — that is enforced and tested. It does **not** follow that leftover bytes were destroyed: the dialog says so, and the folder is the user's to inspect.

---

## §5 Product semantics recorded from the owner's clarification

**S1-10 — Deletion is best-effort by decision, not by omission.** Removal from the visible collection is the contract; physical deletion is an attempt. A failed or unverifiable removal does not restore the item, does not block the UI, and does not raise a deletion failure to the user. The confirmation carries the disclosure, including the possibility that Windows blocks the delete, and the directory path as a clickable link.

**S1-11 — The distinction that matters for reconciliation.** Unexpected disappearance and explicit deletion are different facts. The intent log is the record of the second; reconciliation must read it as *intent*, never as proof of destruction, and must never treat a deleted identifier as `Missing` or as damage. This binds the index foundation, not this slice.

---

## §6 What was run

```
go build ./...                              exit 0
go test . ./internal/backend/chat/ ./internal/durable/   ok (all packages)
go vet ./internal/backend/chat ./internal/durable       clean
gofmt -l (changed Go files)                 clean
node --check frontend/dist/js/chat.js
node --check frontend/dist/js/delete-confirm.js         syntax OK
```

New tests, all passing: `TestWritesOnlyThroughBoundaryHelpers`, `TestAnalysisDetectsABypass`, `TestConversationDirRefusesNonComponentIdentifiers`, `TestAppendTerminatesUnterminatedTail`, `TestSaveRefusedWhileAnotherWriterHoldsTheConversation`, `TestReentrantSaveIsRefusedWithItsOwnCode`, `TestDeletedConversationStaysHiddenAndCannotBeRecreated`, `TestDeletionIsIdempotent`.

---

## §7 Not verified

**S1-12 — Stated plainly.** No Wails build was run and the application was not launched, so the dialog copy, the clickable directory link and `OpenChatDirectory` are **syntax-checked only** — not seen working. No two-process contention test was run against a real second process; the lock's cross-process behaviour is proven at the primitive level by the existing durable tests and here by a second handle in one process. No crash-between-intent-and-removal test was run. No performance measurement was taken on the save or append path, which now perform a containment resolution and an intent-log read per call.

---

**Scope of this record.** It describes code that now exists in the workspace. It is not adoption of v3.2, not a D-2 closure, and not a deployment. D-2's acceptance criteria remain criteria until the unverified items in §7 are exercised on a running build.
