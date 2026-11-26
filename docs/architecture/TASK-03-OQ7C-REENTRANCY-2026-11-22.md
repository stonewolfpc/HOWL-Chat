# TASK-03 — OQ-7c: re-entrancy

| Field | Value |
|---|---|
| Purpose | Can one non-reentrant per-conversation hold contain every history transition? |
| Status | **Read-only investigation. Design only. V3.1 unchanged. Nothing authorized.** |
| Method | Call-graph inspection of every mutation entry point, plus event and seam channels |
| Recovery token | `TASK-03-OQ7C-REENTRANCY-2026-11-22` |

Continues OQ-7's numbering (`O7-103`…).

---

## 1. The two qualifications, now decided

**O7-103 — The seam exclusion is real, not asserted.** `app_native_test_seam.go` and `chat/snapshot_probe_native_test_seam.go` both carry `//go:build native_test_seam`, with `app_native_test_seam_off.go` (`!native_test_seam`) as the counterpart. The seam is therefore **absent from a default build** — the exclusion the owner demanded is a build fact.

**O7-104 — But the writer it calls is compiled into production.** `Service.SaveChatData` lives in `service.go`, untagged. Only the seam method that invokes it is tagged. So the correct statement is: *the seam is build-excluded; the writer remains compiled and uncalled.* A static assertion must still flag it, because nothing but a build tag prevents someone from adding a caller.

**O7-105 — Root resolution is not yet stable or authoritative.** `getAppDataDir` returns `os.UserHomeDir()/AppData/Roaming/HOWLChat`; `paths.IsolationActive()` switches to `StatePath("persisted-chats")`; the memory and service packages use `ProjectRoot()`. Consequences: two Windows users resolve different roots (correctly — but "two instances" therefore means two same-user instances); a roaming-profile redirection moves the root between sessions; an isolated run and a production run never share a lock. **A lock in a sibling namespace is only as good as the root derivation that names it.** Recorded as **OQ-7c-i**.

**O7-106 — Pointer consistency is not implied by removing a deadlock edge.** `saveConfig` writes `<root>/config.json` unguarded, last-writer-wins, no CAS. Separating it from the conversation transition removes the upward edge; it does not make two instances agree on which conversation is active. Recorded as **OQ-7c-ii**.

---

## 2. The actual call graph

**O7-107 — History mutations do call history mutations.** Source-derived, non-test only:

| Entry | Nested persistence calls | Conversations touched |
|---|---|---|
| `App.AppendChatMessage` | `AppendHistoryEntry` | 1 |
| `App.ReplaceChatHistory` | `RewriteHistory` (log + snapshot) | 1 |
| `App.ClearChatHistory` | `RewriteHistory([])` | 1 |
| `App.SaveChatData` | `SaveChat` (7 files) | 1 |
| `App.performSave` | `LoadChat` → `updateChatDataFromAppState` → `SaveChat` | 1 |
| **`App.performFinalSave`** | `performSave` (`LoadChat`+`SaveChat`) → `LoadChat` → **`RewriteHistory`** | 1, **twice rewritten** |
| `PersistenceManager.RenameChat` | `LoadChat` → `SaveChat` | 1 |
| `PersistenceManager.SaveChatAs` | `LoadChat`(source) → `NewChat`→`SaveChat`(dest) → `SaveChat`(dest) | **2** |
| `PersistenceManager.NewChat` | `SaveChat` | 1 |
| **`App.LoadChat`** | **`performSave` of the previous chat**, then `LoadChat` of the next | **2** |
| `App.DeleteChat` | `RemoveAll` | 1 |

**O7-108 — Three paths break the "one lock, one operation" assumption if the boundary is naively wrapped around the caller:**

1. **`performFinalSave`** performs two complete history writes for one conversation — `SaveChat` (which writes `history.json`) followed by `RewriteHistory` (which rewrites both). Wrapped in one hold, this is sequential and safe; but the *intermediate* state has the snapshot newer than the log.
2. **`SaveChatAs`** reads one conversation and writes another. Two locks are needed; taking them in ID order is required, or the copy is torn.
3. **`App.LoadChat`** mutates the conversation being left and reads the one being entered. A single hold spanning both is an upward violation; the correct decomposition is save-then-release, then read.

**O7-109 — Every clean shutdown passes through an incoherent pair.** `SaveChat` writes `history.json` from app state, then `RewriteHistory` writes `history.ndjson` and `history.json` again. Between them the log is older than the snapshot — **log ⊂ snapshot, row 6** — and a crash in that window persists it. The transient state is not a bug in the write order; it is the consequence of writing both representations twice in one operation.

---

## 3. Callbacks and event dispatch

**O7-110 — No persistence call emits an event inside itself.** All `emitEvent` sites are in the streaming loops (`chat:chunk`, `chat:reasoning`) or after the mutation returns (`chat:complete` after `AppendChatMessage` at `app-prompt.go:138-139`). The boundary can therefore stay **callback-free without changing product behaviour**, provided emission remains outside the hold.

**O7-111 — Two re-entry channels still exist and must be excluded by rule, not assumed absent:**

| Channel | Behaviour | Constraint |
|---|---|---|
| `runtime.EventsEmit` → frontend | asynchronous; frontend may call a binding in response | no binding may be invoked from inside a hold; a `chat:complete` handler calling `AppendChatMessage` must find the hold already released |
| `devEventSink` | **synchronous** call into UI-harness code when set (`app.go:369-372`) | a harness callback that mutates would re-enter the same conversation; the boundary must forbid mutation from a sink callback |

**O7-112 — A non-reentrant lock converts re-entry into a deadlock, not into corruption.** That is the correct failure direction — but only if the lock has a bounded wait and reports `Busy` (row 567 §O7-99). A re-entrant acquisition attempt that blocks forever would be an availability defect, and one that silently succeeds would be worse.

---

## 4. Answer

**O7-113 — One non-reentrant per-conversation hold can contain every history transition, but only with three decompositions and one behavioural decision:**

| # | Change | Why |
|---|---|---|
| 1 | `performFinalSave` becomes **one** boundary operation writing both representations once | removes the double write and the row-6 window (O7-109) |
| 2 | `SaveChatAs` becomes a two-conversation operation with a stated lock order | prevents a torn copy (O7-108) |
| 3 | `App.LoadChat`'s implicit save becomes an explicit boundary call that **returns before** the load | no upward hold, no re-entry (O7-108) |
| 4 | Emission and sink callbacks stay strictly outside the hold | keeps the boundary callback-free (O7-110, O7-111) |

**O7-114 — No product behaviour must change to keep the boundary callback-free.** Progress and completion events already sit outside persistence calls; only discipline is required, and one rule: no mutation from an event handler while a hold is held — which the asynchronous dispatch already guarantees for the frontend, and a rule must guarantee for the synchronous sink.

**O7-115 — Nothing here is authorized.** No decomposition, no merged final-save, no lock ordering, no pointer change, no emission rule, no gate, CAS, marker, migration or deletion change.

---

## 5. Open on this branch

| # | Question |
|---|---|
| OQ-7c-i | Root derivation: what makes the lock root stable and authoritative across profiles, isolation and redirection? (O7-105) |
| OQ-7c-ii | Pointer consistency: `config.json` is unguarded, last-writer-wins; what replaces that? (O7-106) |
| OQ-7c-iii | Does anything require true re-entrancy (an operation that must call another while holding the hold), or are all four decompositions sufficient? |
| OQ-7d-ii | Retry identity — unchanged, and now more pressing: a merged final-save and a two-conversation save-as both need stable operation identity |

---

## 6. Tests this would require

- **O7-T23** Re-entrant acquisition on the same conversation reports `Busy` (or a typed re-entry refusal) rather than deadlocking indefinitely.
- **O7-T24** `performFinalSave` as one operation leaves rows 14–17 and never a persisted row 6.
- **O7-T25** `SaveChatAs` with a defined lock order produces a copy that is either the source at a single instant or an explicit failure.
- **O7-T26** `App.LoadChat` saves the previous conversation before acquiring the next, and never holds two conversation locks.
- **O7-T27** No `emitEvent` call executes while a conversation hold is held (instrumented).
- **O7-T28** A `devEventSink` callback that attempts a mutation is refused, not deadlocked.
- **O7-T29** The default build contains no `native_test_seam` symbols, and the static assertion still lists `Service.SaveChatData` as a compiled writer.
- **O7-T30** Two instances with different resolved roots are shown not to share a lock — the finding, asserted as a limitation rather than a bug.
