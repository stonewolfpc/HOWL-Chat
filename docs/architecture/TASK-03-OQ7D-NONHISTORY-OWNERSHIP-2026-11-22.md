# TASK-03 — OQ-7d: non-history ownership and lifecycle

| Field | Value |
|---|---|
| Purpose | Who owns each non-history file in a conversation directory, and what lifecycle it has |
| Status | **Read-only inventory. Design only. V3.1 unchanged. Nothing authorized.** |
| Method | Source inspection: writers, readers, path bases, wipe and duplication paths, and reachability |
| Recovery token | `TASK-03-OQ7D-NONHISTORY-OWNERSHIP-2026-11-22` |

Continues OQ-7's numbering (`O7-64`…). Contains three corrections to earlier rows (§1).

---

## 1. Corrections to what this branch already published

**O7-64 — Withdrawn: the case-collision claim.** Row 565 §O7-54 said `chats` and `Chats` could collapse into one directory on a case-insensitive volume. **They cannot, because the roots differ** (Source-derived):

| Writer group | Root | Subdirectory |
|---|---|---|
| `PersistenceManager` | `%APPDATA%\HOWLChat` (isolation: `<StateRoot>/persisted-chats`) | `chats` |
| `chat.Service.SaveChatData` | `ProjectRoot()` — the installation root | `Chats` |
| memory package (service, continuity, candidates) | `ProjectRoot()` | `Chats` |

These are **two separate trees**, not one directory under two identities. The identity difference (ID versus display name) stands; the collision consequence does not.

**O7-65 — Corrected: `memories.json` is not one contested file.** Row 556 said `SaveChat` and the memory service have "two uncoordinated writers" of `memories.json` with last-writer-wins. In the default configuration they write **different files** — `<appdata>/chats/<id>/memories.json` from `ChatData.Memories`, and `<install>/Chats/<id>/memories.json` from `[]types.MemoryRecord`. The defect is not a race; it is **two divergent memories stores with different schemas**, one of which reconciliation never sees.

**O7-66 — Corrected: the second history writer is dormant, not live.** Row 565 §O7-51 described `Service.SaveChatData` as a live second writer. It is reachable **only through `NativeTestSaveChatData`**, a native test seam (`snapshot_probe_native_test_seam.go:18`); no production path calls it. The finding stands — a test seam can still produce a shadow `history.json` — but its severity is materially lower than row 565 implied.

---

## 2. The inventory

**O7-67 — Twelve non-history files, four writer packages, one of which is unreachable except through a test seam.**

| File | Writer(s) | Reader(s) | Path basis | Reachability | Owner |
|---|---|---|---|---|---|
| `meta.json` | `SaveChat` (atomic) | `LoadChat`, active-chat lookup | `chats/<id>` | live | conversation identity |
| `world.json` | `SaveChat`, only when non-nil | `LoadChat` | `chats/<id>` | live | conversation projection of `worldstore` content |
| `scenario.json` | `SaveChat`, only when non-nil | `LoadChat` | `chats/<id>` | live | conversation projection |
| `characters.json` | `SaveChat` | `LoadChat` | `chats/<id>` | live | conversation projection |
| `lorebooks.json` | `SaveChat` | `LoadChat` | `chats/<id>` | live | conversation projection |
| `memories.json` (app-data) | `SaveChat` (atomic) | `LoadChat` | `chats/<id>` | live | **user-authored content** |
| `memories.json` (install) | `MemoryStorage.saveMemories` (plain write) | `loadMemories` | `Chats/<id>` | enclosing methods' production callers **Unestablished** | user-authored content, second store |
| `memories_backup.json` | `saveMemories` — writes `{}` then content | none found | `Chats/<id>` | same | application-owned, misnamed |
| `memories_snapshot.json` | `Service.SaveChatData` | `loadMemoriesSnapshot` | `Chats/<name>` | **test seam only** | application cache |
| `chat.json` | `Service.SaveChatData` | export path (`contentItemFile`, `library.Chats`) | `Chats/<name>` | **test seam only** | application session metadata, doubles as a library item |
| `auto_capture_config.json` | `SetAutoCaptureConfig` | `GetAutoCaptureConfig` | `Chats/<id>` | **live** (`App.SetAutoCaptureConfig`) | application setting |
| `continuity_config.json` | `SaveContinuityConfig` | `LoadContinuityConfig` | `Chats/<id>` | **writer: no caller found**; read live | application setting |
| `candidates.json` | `saveCandidates` | candidate storage | `Chats/<id>` | reachable via `ProcessMessageForEvents`; production caller **Unestablished** | application cache |

**O7-68 — Every one of these uses plain writes or `saveJSON`.** None uses `DurableStore`. Two (`saveMemories`, `SetAutoCaptureConfig`, `SaveChatData`) are non-atomic `os.WriteFile` calls, so a crash mid-write truncates the file.

**O7-69 — Two identity bases remain in force**, now across trees: ID for persistence and most memory files, display name for `chat.json` and `memories_snapshot.json`. A rename changes the name-keyed files' address; nothing migrates them.

---

## 3. Ownership classification

**O7-70 — Four classes, and only one of them is the user's data:**

| Class | Files | Consequence |
|---|---|---|
| conversation identity | `meta.json` | required for the directory to mean anything |
| conversation projection | `world`, `scenario`, `characters`, `lorebooks` | reconstructible from `worldstore` and selections; losing them degrades but does not destroy content |
| **user-authored** | `memories.json` (both trees) | **irreplaceable**; the only non-history file whose loss is data loss |
| application-owned | `chat.json`, `memories_snapshot`, `memories_backup`, `auto_capture_config`, `continuity_config`, `candidates` | settings and caches; regenerable or re-settable |

**O7-71 — `memories_backup.json` is not a backup.** It is overwritten with `{}` immediately before each save, so at the moment a write fails the "backup" holds an empty object. Classified application-owned; whether its behaviour is intentional is **Unestablished** and out of scope here.

---

## 4. Lifecycle

**O7-72 — Creation.** `SaveChat` creates the directory and writes seven files (`meta`, `history.json`, `world`, `scenario`, `characters`, `lorebooks`, `memories`). Every install-tree file is created on demand by its own package, at whatever moment that package first writes.

**O7-73 — Update.** `SaveChat` overwrites its seven from in-memory `ChatData`. The install-tree files are overwritten by their owners independently, with no coordination.

**O7-74 — Duplication copies seven, drops five.** `SaveChatAs` → `NewChat` → copy fields → `SaveChat` in a new directory. Copied: `meta`, `history.json`, `world`, `scenario`, `characters`, `lorebooks`, `memories`. **Not copied:** `chat.json`, `memories_snapshot.json`, `auto_capture_config.json`, `continuity_config.json`, `candidates.json` — and `history.ndjson`, which is not created (row 550's collision aside, the duplicate legitimately starts with the snapshot only).

**O7-75 — Wipe removes one tree.** `WipeAllData` stages `<RootDir>/chats` only (`app-wipe.go:46`). Every install-tree file survives the wipe: memories, settings, caches, and any shadow `history.json`.

**O7-76 — Deletion removes one directory.** `DeleteChat` calls `RemoveAll` on `<RootDir>/chats/<id>`. The install-tree counterparts are left in place.

---

## 5. Adjacency does not establish shared ownership

**O7-77 — Verdict: seven of the twelve are application-owned, four are projections, one is user content.** None of them inherits the history contract's authority, and none may be pulled inside the history boundary because it happens to sit beneath the same directory.

**O7-78 — The converse must also hold: the history boundary must not claim to have protected them.** A boundary that reports "the conversation is intact" while `memories.json` in the install tree was truncated mid-write would be describing a different directory from the one the user means.

**O7-79 — OQ-7d's own finding, outside reconciliation's scope:** a **shadow conversation tree** exists at `ProjectRoot()/Chats/<id|name>/`. V3.1 scans `<rootDir>/chats` only. If the test seam ever runs, a `history.json` exists there that row 1 of the decision table reports as `Missing` — because it is absent from the scanned tree while present in another. This is **not** to be absorbed into the reconciliation contract; it is recorded as a live-scope question (**OQ-7d-i**).

---

## 6. Partial success, restated as an API requirement

**O7-80 — A result must be three-valued, not a boolean:**

| Outcome | Meaning |
|---|---|
| refused before content changed | the store is exactly as it was |
| content changed, marker not recorded | **unresolved on-disk transition awaits reconciliation** |
| both content and marker recorded | completed |

**O7-81 — The middle outcome must never render as a no-op**, and compensation by rollback is forbidden (row 565 §O7-47). Idempotent retry is the direction; its precise retry semantics — which operation identity is retried, and how a retry distinguishes "already applied" from "not applied" — still need proof (**OQ-7d-ii**).

---

## 7. Standing

**O7-82 — Nothing here is authorized.** No path consolidation, no second-writer repair or removal, no memories-store unification, no wipe or deletion change, no boundary extension, no marker change, no migration, no backfill. This inventory establishes the deficit and constrains the design; it does not begin the work.

**OQ-7d-i** — is the install-tree conversation store live, vestigial, or a migration target? **OQ-7d-ii** — what proves retry idempotency at the operation level? **OQ-7b** (lock scope) and **OQ-7c** (re-entrancy) remain open and are now better informed: the writer population spans two roots and four packages.
