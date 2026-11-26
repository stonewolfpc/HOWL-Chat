# Decision — reader sharing during atomic replacement (ADR-001 §8.1)

| Field | Value |
|---|---|
| Task | `TASK-02` — DurableStore, step 1 |
| Resolves | ADR-001 §8.1: *"Compatible reader sharing is an acceptance requirement of the replacement protocol"* — currently **`UN`** |
| Date | 2026-11-22 |
| Status | **Decided** (architecture decision for TASK-02's scope) |
| ADR amendment required | **None.** This supplies a value for an `UN` field; it changes no normative statement in ADR-001 |
| Recovery token | `TASK-02-READER-SHARING-2026-11-22` |

---

## 1. The question

During atomic replacement, what may a concurrent reader do, what does it observe, and who is responsible when
replacement cannot proceed?

## 2. Requirements and evidence

### 2.1 From ADR-001 (authoritative)

- **§8 D4:** lock contention → **`Busy`, no bypass**.
- **§8 rotation ordering:** `acquire lock → write temp → flush → pre-image copy → rename → drop oldest beyond N
  (only after rename succeeded) → release`. A failed commit never rotates.
- **§8 two timers:** lock acquisition (default **5000 ms**) and rename retry (**3 × 50 ms**) *inside* the acquired
  lock.
- **§8.1:** replacement requires delete access to the target; "a reader that does not share delete blocks it". The
  current settings reader is a short-lived `os.ReadFile`; that "does not establish that overlap cannot occur, and
  does not establish delete sharing". Two protocols are offered: **short-lived read + bounded writer retry**, or a
  **delete-sharing reader** whose exact flag is `UN`.
- **§8:** `os.WriteFile` is "a direct, non-atomic write — no temp, no flush, no rename, no backup" and DurableStore
  "must replace this; it must not be assumed to exist".

### 2.2 Measured on this build (go1.27.3 windows/amd64, NTFS, `C:`)

A probe held the live document open with `os.Open` — the exact handle `os.ReadFile` uses mid-read — and attempted
both replacement and in-place write:

| Operation | Result |
|---|---|
| `os.Rename(next, live)` while a reader holds `live` | **`Access is denied`** — replacement is blocked |
| `os.WriteFile(live, …)` while a reader holds `live` | **succeeds** |
| `os.Rename` after the reader closes | succeeds |

**What this establishes.** Go's `os.Open` on Windows shares read and write but **not delete**, so a held reader
blocks replacement and therefore cannot coexist with a rename. It also establishes something worse about the
current implementation: a **direct write truncates a document a reader may be holding**, so the reader can observe a
partially written file. Non-atomicity is not merely theoretical here; it is observable.

### 2.3 What existing readers actually do

| Reader | Pattern | Holds the handle across other work? |
|---|---|---|
| `internal/settings.Store.Load` | `os.ReadFile` — open, read, close | No |
| `internal/backend/chat/persistence.go` loaders | `os.Open` + scanner, deferred close within the call | No |
| `internal/library/registry.go` | read within the call | No |
| `internal/backend/memory/service.go` | read within the call | No |

No reader memory-maps a document, keeps a handle in a struct field, or opens a document with a custom share mode.
No code in the repository calls `LockFileEx`; the only lock is an in-process `TryLock` (`internal/image/engine/sdcpp/video.go:154`).

### 2.4 The falsified premise in the current code

`internal/settings.Store.Save` justifies its direct write with: *"a direct `os.WriteFile` is instantaneous and
eliminates the complexity of temp-file rename dances that cause Windows file-locking errors and leave stale
`.tmp`/`.bak` files behind."*

The measurement shows the causation is inverted. Rename does not fail because rename is fragile; it fails because
**a reader holds the file**. The fix is shorter reader holds plus a bounded writer retry — not abandoning atomicity.
Recorded as finding **R-1**; not changed by this task.

---

## 3. Alternatives

| # | Model | Reader change | Writer change | Assessment |
|---|---|---|---|---|
| **A** | **Short-lived reader + bounded writer retry** | none; readers must keep holds short | retry rename 3 × 50 ms inside the lock, then `Busy` | **Selected.** Uses no capability that is `UN`, works with every existing reader, and is the ADR's first listed protocol |
| B | Delete-sharing readers (`FILE_SHARE_DELETE`) | every reader must open with a custom handle | none beyond A | **Rejected for now.** `os.Open` cannot request it; requires `NtCreateFile` (`UN`, capability gate per ADR §7) and rewrites every reader call site |
| C | Readers take shared locks | readers become lock clients; a writer blocks all reads | writer unchanged | **Rejected.** Turns reads into failures during writes and adds a second lock discipline |
| D | Readers read pre-images / snapshots instead of live | readers must know which file is current | writer maintains a pointer | **Rejected.** Introduces a second authority for "which document is live" — the exact defect the rebuild removes |

## 4. Selected behaviour

**A — short-lived reader with a bounded writer retry.**

### 4.1 Reader contract

1. A reader opens the document, reads it to completion, and closes it **within one call**. No handle is retained in
   a struct, cache or goroutine.
2. A reader never memory-maps a document, never opens it with a custom share mode, and never re-reads a handle it
   already holds.
3. A reader opens the **live** path only. It never opens a temp file and never opens a pre-image — those are
   recovery artifacts, not readable state.
4. If a reader's open or read fails because replacement is in progress, it retries at the caller's discretion or
   reports the failure. **The store does not paper over reader failures** and does not serve a stale in-memory copy.

### 4.2 Writer contract

1. Acquire the exclusive lock on `<doc>.lock` (byte range `[0,1)`). Not acquired within **5000 ms** → **`Busy`**.
   No bypass, no "overwrite anyway".
2. Write the temporary document beside the live one; flush; close.
3. Copy the current live document to `v1` pre-image, shifting `v1 → v2`; **never** discard the newest pre-image on a
   failed write.
4. Commit the manifest describing the intended new state (document hash, schema, pre-image list).
5. **Rename** temp → live with **3 attempts × 50 ms**, all inside the held lock. If replacement is still denied, the
   writer returns **`NotSaved(rename_denied)`** — not `Busy`, because the block is a reader, not another writer —
   removes the temp file, and **does not rotate pre-images**.
6. Only after a successful rename, drop pre-images beyond `N = 2`.
7. Release the lock.

### 4.3 What a reader may observe

| Situation | Reader observes |
|---|---|
| No write in progress | the committed document |
| Writer between temp write and rename | the **previous** committed document (complete) |
| Writer during the rename | the previous or the new document; never a mixture |
| Rename denied (reader held too long) | previous document; the writer reports `NotSaved` |
| Interrupted process (crash before rename) | previous document; a stale temp remains and is swept on next open |
| Crash after rename | the new document |
| Corrupt live document | recovery from the newest pre-image whose hash matches the manifest, reported as `Recovered(from)` |
| Live corrupt and no valid pre-image | **`NotSaved(corrupt_no_valid_preimage)`** — never a silent empty document |

**A reader never observes a partially written document**, because no path writes into the live file: the only
mutation of the live path is a rename.

### 4.4 Lock ownership and release

`LockFileEx` on `<doc>.lock`, exclusive, range `[0,1)`, via `golang.org/x/sys/windows` (already a dependency).
Windows releases the lock when the owning process terminates, so an abandoned lock cannot wedge the application. A
leftover lock **file** is harmless and is not "stale state" to delete — deleting it would defeat mutual exclusion.

### 4.5 Compatibility with the existing single-writer and `Busy` rules

- Readers take **no lock**, so `Busy` can arise only from writer-writer contention (lock timeout) or from a rename
  denied after retries (reported as `NotSaved(rename_denied)`, distinguished from `Busy`).
- Every writer takes the lock — save, reset, migration importer, recovery, rotation. No exemptions.
- The single-writer guarantee is unchanged; this decision adds no second authority.

## 5. Consequences and limitations

- **Requires reader discipline that is not mechanically enforced.** The store cannot detect a reader that holds a
  handle for a long time; it surfaces the consequence as `NotSaved(rename_denied)`. Enforcing short holds is a
  review obligation for every call site that adopts the store.
- **Protocol B remains available and unselected.** If measurement later shows long reader holds, the upgrade path is
  delete-sharing readers behind the `UN` capability gate — a call-site change, not a contract change.
- **Power-loss durability is untouched and remains `UN`.** Flush ordering is implemented; NTFS metadata durability
  is not claimed.
- **Rename is not free of platform behaviour.** Replacement needs delete access to the target, which is why the
  reader contract is load-bearing rather than advisory.

## 6. ADR amendment

**None required.** ADR-001 §8.1 already states the requirement, lists this protocol first, and leaves the value
`UN`. This record supplies it:

> §8.1, "Reader-open contract": the selected reader protocol is **short-lived read with a bounded writer retry**
> (3 × 50 ms inside the acquired lock, then `NotSaved(rename_denied)`), measured on go1.27.3 windows/amd64: a reader
> holding the document with `os.Open` blocks `os.Rename` while permitting `os.WriteFile`. No reader may retain a
> handle, memory-map a document, or use a custom share mode.

**ADR-001 is not rewritten.** The paragraph above is the exact wording proposed *if* the owner later chooses to fold
this decision into the ADR. Until then, this record is the decision of record for TASK-02.
