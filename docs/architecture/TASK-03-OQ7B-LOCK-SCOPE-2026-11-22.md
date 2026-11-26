# TASK-03 — OQ-7b: lock scope

| Field | Value |
|---|---|
| Purpose | Which lock scope can make one conversation's history transition atomic against everything that touches it |
| Status | **Read-only investigation. Design only. V3.1 unchanged. Nothing authorized.** |
| Method | Source inspection of the store's lock placement, the writer population, and the creation/delete/wipe paths |
| Recovery token | `TASK-03-OQ7B-LOCK-SCOPE-2026-11-22` |

Continues OQ-7's numbering (`O7-83`…).

---

## 1. Two guardrails carried forward

**O7-83 — Reachability is four-valued, and a production claim may not drift.** Every writer statement must name its class:

| Class | Meaning | Example |
|---|---|---|
| production | reached by a shipped path | `SaveChat`, `AppendHistoryEntry`, `RewriteHistory` |
| test-seam | reached only through a native test hook | `Service.SaveChatData` via `NativeTestSaveChatData` |
| dormant | no caller found in the repository | `SaveContinuityConfig` |
| unestablished | callers exist but production reachability was not proven | `saveMemories`' enclosing methods |

**O7-84 — "`memories.json` is the only irreplaceable non-history data" is an inventory conclusion, not a guarantee.** Row 566 derived it from the twelve files it examined; preservation requirements must be established per file before any wipe, duplication or deletion work.

---

## 2. What the current lock actually covers

**O7-85 — The store's lock is per document, beside the document, and held only for that document's write.** `LockPath()` is `s.path + ".lock"`; `Save` acquires it, writes, and releases. Defaults: `LockTimeout` 5 s, `RenameAttempts` 3, `RenameWait` 50 ms, `PreImages` 2, 25 ms poll to `ErrBusy`.

**O7-86 — Nothing else in a transition is guarded.** The artefacts and their current protection:

| Artefact | Location | Writer | Lock today |
|---|---|---|---|
| marker (planned) | `<chat>/expectation.json` | `DurableStore` | its own `.lock`, during `Save` only |
| log | `<chat>/history.ndjson` | `AppendHistoryEntry`, `RewriteHistory` | **none** |
| snapshot | `<chat>/history.json` | `SaveChat`, `RewriteHistory`, `saveJSON` | **none** |
| companions | `<chat>/meta|world|scenario|characters|lorebooks|memories.json` | `SaveChat` | **none** |
| active-chat pointer | `<root>/config.json` | `saveConfig`, plain `WriteFile` | **none** |

**O7-87 — `SaveChat` couples the pointer to the transition.** Every save ends with `pm.config.LastChatID = …; pm.saveConfig()` (`persistence.go:244-247`), an unguarded root-level write on the hot path. The pointer is what `App.LoadChat` reads to decide which conversation to save next (`app-chat.go:95-99`), so it is both contended and load-bearing.

---

## 3. Four candidate scopes against the writer population

**O7-88 — Per-document is insufficient by construction.** It serialises marker saves and nothing else. An append and a save can still interleave on the log and snapshot, producing rows 6 and 9 exactly as probe 12 did. The current scope cannot satisfy V3-26 even if the gate is implemented on top of it.

**O7-89 — Per-conversation is the minimum sufficient scope.** All four transition artefacts live in one directory, so one lock over that conversation serialises every operation that changes its representations or its marker: append, rewrite, clear, save, duplication-target, deletion. It also keeps contention O(1) — a slow conversation cannot block another.

**O7-90 — Per-root is correct but needlessly contended.** It serialises all conversations, so a 100 k-record rewrite (measured ~105 ms) blocks every other chat for its duration, and a slow or disconnected volume converts that into a 5 s stall per operation before `ErrBusy`. Correctness is not improved over per-conversation for any single-conversation transition.

**O7-91 — Global across both trees is unnecessary and would overstate the claim.** No operation spans the two trees: persistence writes `%APPDATA%\HOWLChat`, the memory and service packages write `ProjectRoot()`. Claiming a global lock would describe a guarantee no code path needs.

---

## 4. The constraint that decides the design: a lock cannot live inside what it protects

**O7-92 — Three paths break if the lock file sits in the conversation directory:**

| Path | Failure |
|---|---|
| **creation** | `NewChat` writes into a directory that does not yet exist; there is nowhere to put the lock before `MkdirAll` |
| **deletion** | `DeleteChat` calls `RemoveAll(<chat>)` while the deleter holds `<chat>/…lock` open; on Windows an open file cannot be removed |
| **wipe** | `WipeAllData` renames the whole `chats` tree; an open handle inside prevents a directory rename on Windows |

**O7-93 — Therefore the lock belongs outside the protected conversation directory**, in a sibling namespace such as `<root>/locks/<id>.lock`, and outside `chats/` so a wipe's rename does not move it. This is a placement requirement, not a preference: the naive `path + ".lock"` placement is correct for a single document and wrong for a directory scope.

**O7-94 — Creation still needs an exclusive name.** `generateChatID` is second-resolution and `MkdirAll` is non-exclusive (row 550), so creation must acquire a lock keyed by the candidate ID *before* the directory exists — which O7-93's namespace allows, and the current placement does not.

---

## 5. Two locks, and the ordering rule that avoids deadlock

**O7-95 — A root-level lock is required for wipe, and only for wipe.** Renaming `chats/` (and staging the gallery root) is a root-scoped operation. Per-conversation locks cannot make it atomic.

**O7-96 — Transitions must not take the root lock, and wipe must not take conversation locks in the opposite order.** If a transition holds a conversation lock and then needs the root lock, while wipe holds the root lock and then visits conversations, the two orders admit deadlock.

**O7-97 — Resolution: the active-chat pointer leaves the conversation transition.** `config.json` is not a representation of any conversation; it is an application pointer. It must have its own document and its own lock, acquired **after** the conversation lock is released — never while held. Then the hierarchy has one upward edge only (root → conversation), held by wipe alone.

**O7-98 — Stated order, for the eventual implementation:** acquire root only when the operation is root-scoped; acquire at most one conversation lock per operation; never acquire a conversation lock while holding the root lock except in wipe's single-threaded enumeration, in which each conversation lock is taken and released before the next.

---

## 6. Contention and refusal

**O7-99 — Refusal is the observable behaviour to specify, not the timeout value.** With a 5 s timeout and a 25 ms poll, contention surfaces as `ErrBusy`. That must reach the caller as a typed refusal (V3-36, row 564 §O7-36); today four sites discard it. Per-conversation scope keeps the *frequency* of refusal low; it does not remove the obligation to surface it.

**O7-100 — The test seam never contends with production.** `NativeTestSaveChatData` writes the installation tree, while the production lock lives in the appdata tree. Per-conversation locking therefore does not cover the seam, and the enforcement claim must exclude it by build configuration rather than by lock placement.

---

## 7. Answer to the key question

**O7-101 — Per-conversation scope is the only one that can guarantee a conversation's complete transition cannot interleave with any other operation changing its representations or marker — and only under all five:**

1. the lock lives outside the conversation directory (O7-93);
2. log, snapshot, companions and marker are written under **one** hold (O7-89);
3. creation acquires the lock before the directory exists (O7-94);
4. the active-chat pointer is not part of the transition, and wipe is the only upward acquisition (O7-97, O7-98);
5. `Busy` and every other failure surface as refusals (O7-99).

Per-document fails (2). Per-root satisfies (2) but not (1)–(3) economically and raises contention to O(all conversations). Global is unnecessary and would misdescribe the system.

**O7-102 — Nothing here is authorized.** No lock placement, no lock namespace, no pointer extraction, no wipe ordering, no gate, no CAS, no marker, no migration, no deletion change. This establishes which scope is sufficient and what it costs; it does not build it.

---

## 8. Tests this would require

- **O7-T15** Two concurrent appends to one conversation serialise; writes to different conversations do not block each other.
- **O7-T16** An append interleaved with a save under the conversation lock leaves rows 14–17, never row 6 or 9.
- **O7-T17** Creation acquires the lock before the directory exists; two same-second creations do not collide (the row 550 shape is refused, not merged).
- **O7-T18** Deletion succeeds while holding the conversation lock — proving the lock is not inside the removed directory.
- **O7-T19** Wipe during an active conversation transition leaves neither a half-renamed tree nor an orphaned lock.
- **O7-T20** The pointer write never occurs while a conversation lock is held (instrumented ordering assertion).
- **O7-T21** `Busy` surfaces as a typed refusal at every entry point.
- **O7-T22** The test seam is excluded from the production enforcement claim by configuration, and the claim says so.
