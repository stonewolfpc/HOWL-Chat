# TASK-03 — OQ-7c-ii: active-chat pointer consistency

| Field | Value |
|---|---|
| Purpose | Whether `LastChatID` is a preference or participates in conversation authority, recovery, or mutation targeting — and what contract it therefore needs |
| Status | **Read-only investigation. Design only. V3.1 unchanged. Nothing authorized.** |
| Method | Source inspection of every read and write of the pointer, the startup restore path, and the wipe sequence |
| Recovery token | `TASK-03-OQ7CII-POINTER-CONSISTENCY-2026-11-22` |

Continues OQ-7's numbering (`O7-130`…). Answers OQ-7c-ii as opened by `O7-106` (row 568) and `O7-87` (row 567).

---

## 1. What the pointer actually is

**O7-130 — There are three selection records, and only one of them holds an identifier:**

| Record | Holds | Written by | Read by |
|---|---|---|---|
| `AppConfig.LastChatID` (`config.json`) | a chat ID | five sites (§1.2) | every write path (§2) |
| `HOWLStorage` `howl-chat-current-chat` | **a full content copy plus an ID** | frontend, on load and on sync | frontend startup, delete service, state restore |
| recency (no stored value) | nothing | — | `discoverMostRecentChatSession`, mtime ordering |

The frontend key is already recorded as a competing *content* authority (row 523, row 524). It is also a second *selector*. The backend pointer is the only selection record the Go side consults.

**O7-131 — The decisive structural fact: the application holds no active-conversation identifier of its own.** `App` carries `chatMessages`, `world`, `scenario`, `characters`, `lorebooks`, `triggeredLore` and `chatPersistence` (`app.go:50-68`), and no chat ID field. So the process's knowledge of *which conversation it is currently holding* exists **only** as `pm.config.LastChatID` — an in-memory field of a structure that also mirrors that value to a file shared with every other process on the profile.

That single fact explains most of what follows: a field carrying process-local state is being carried across call boundaries and process boundaries at the same time.

**O7-132 — Five write sites, four of them side effects:**

| Site | Trigger | Sets to |
|---|---|---|
| `persistence.go:183-186` | `NewChat` | the chat just created |
| `persistence.go:245-247` | **`SaveChat`** | the chat just saved |
| `persistence.go:478-481` | `DeleteChat`, only if it matches | `""`, **error discarded** (`_ = pm.saveConfig()`) |
| `persistence.go:545-547` | `SetLastChatID` (explicit selection) | the argument, unvalidated |
| `persistence.go:105-106` | `loadConfig` when the file is absent | `""` |

Only one of the five is an explicit selection. `SaveChat` — the hottest write in the package — repoints the file as a side effect of persisting content (already recorded as `O7-87`).

**O7-133 — `saveConfig` is the only non-atomic document write in the package.** Every other JSON document goes through `saveJSON`, which writes `name.tmp` and renames. `saveConfig` calls `os.WriteFile` directly on `config.json` (`persistence.go:123-137`): no temp file, no rename, no fsync. So the one file whose value coordinates selection has the weakest durability of any file in the store — a torn write leaves unparseable JSON, `loadConfig` returns an error, `initialize` fails, and `NewApp` logs the failure and leaves `chatPersistence` nil (`app.go:159-165`). No conversation is harmed; the application loses persistence entirely for that run.

---

## 2. What the pointer decides

**O7-134 — Ten consumers resolve a *content mutation target* from it:**

| Consumer | Effect |
|---|---|
| `performSave` → `SaveChat` (`app-chat.go:314-325`) | rewrites the whole document set for the named chat |
| `performFinalSave` → `RewriteHistory` (`app-chat.go:326-352`) | **rewrites `history.ndjson` from in-memory state at shutdown** |
| `appendChatMessage` (`app-prompt.go:289-310`) | appends one log line |
| `ReplaceChatHistory` (`app-prompt.go:383-392`) | rewrites the log from supplied entries |
| `ClearChatHistory` (`app-prompt.go:400-402`) | truncates the log |
| `SaveChatAs` (`app-chat.go:76-80`) | reads the source conversation |
| `AddWorldToChat` / `AddScenarioToChat` (`app-characters.go:241,256`) | `performSave` |
| `loadChatDataFromStorage` (`app-characters.go:273`) | reads the active conversation |
| `DeleteChat` (`app-chat.go:121-129`) | clears in-memory state for the named chat |
| `GetActiveChatMeta` (`persistence.go:550-565`) | reads `meta.json` for the named chat |

**O7-135 — So the pointer is load-bearing for targeting, and it is still not a representation.** No conversation's content is derived from the value; nothing reads it to decide what a conversation *contains*. Holding a target selector is not the same as being an authority over the thing selected.

---

## 3. Does it participate in authority or recovery?

**O7-136 — Not in content authority.** Under V3.1 the six states are decided by per-conversation content. The pointer holds no content, so:

| Proposition | Holds? |
|---|---|
| its loss can lose conversation content | no |
| its corruption can fabricate content | no |
| it can make a conversation `Healthy` | no |
| it can make a conversation `Unavailable` | no — a dangling pointer names nothing, and nothing is the absence of a directory, not a pointer state |

The pointer is therefore **outside the reachability domain entirely** (`O7-149`).

**O7-137 — One narrow authority edge exists.** `GetActiveChatMeta` backfills an empty `meta.ID` from the pointer (`persistence.go:561-563`). The pointer can *supply* an identity that a `meta.json` failed to carry. The edge is bounded: it completes an identity for display, it never alters content, and it does not make such a conversation Healthy — a scanner reading the same `meta.json` still sees no ID and must report identity unestablished.

**O7-138 — It does participate in recovery, as a selector — and the recovery path rewrites it.** Startup is `GetActiveChatMeta` first; on any error the frontend falls back to `discoverMostRecentChatSession`, then `loadChatSessionById` → `App.LoadChat` → `SetLastChatID` (`sidebar-session.js:101-127`; `app-chat.go:107`). So an unreadable `meta.json` does not merely fail to restore one conversation: it **silently re-selects the newest conversation and repoints the file**. That is an automatic repair of the pointer, performed by the application on a damaged conversation's behalf.

Two obligations follow. Reconciliation must not imitate it — re-pointing is a selection change, not a recovery, and must never be reported as one. And a repointed value must never be read as evidence about the damaged conversation, whose state remains `Unavailable` for its own reasons.

---

## 4. Concurrency

**O7-139 — The field carries two meanings at once.** Every save writes the ID it just persisted (`O7-132`), so `LastChatID` means both *what the user last selected* and *what this process last persisted*. The two agree only in single-instance use, which is why the conflation has gone unnoticed.

**O7-140 — Stale writes are real and they revert selection, not content.** Last-writer-wins with no CAS and no version. Instance A selects Y; instance B, still holding X loaded at launch, saves or shuts down and writes X. The pointer now names the older selection even though nothing about Y changed. The next launch opens X. Cost: one re-selection. **No conversation content is affected by the regression itself**, which is the fact that decides the contract class (§6).

**O7-141 — The hazard is not the value; it is the shutdown overwrite the value selects.** `performFinalSave` rewrites `history.ndjson` from *that process's* in-memory state, and the pointer is only how it finds the target. So a stale pointer at shutdown means instance B rewrites conversation X from B's possibly-stale memory while the user is working in Y. That hazard survives any pointer contract, because it is a writer-authority question; it is decided by the per-conversation lock and generation (OQ-7b, OQ-7c), not by the pointer. **Removing the pointer from the conversation lock's critical section (`O7-106`) removes an ordering edge; it does not make two instances agree on which conversation is active, and it does not make the shutdown overwrite safe.** This is the explicit answer to the question OQ-7c-ii was opened with.

**O7-142 — Selection during save is two config writes with no atomicity across them.** `LoadChat` first `performSave`s the previous conversation — which writes the pointer to the *old* ID — then loads, then `SetLastChatID(new)` writes again (`app-chat.go:95-108`). A crash between the two leaves the pointer naming the previous conversation while the caller's `LoadChat` had not returned. Preference loss only; no content is retired, and the loader on the next launch restores a conversation that still exists.

**O7-143 — `DeleteChat` clears the pointer only on a match, and discards the failure.** `_ = pm.saveConfig()` (`persistence.go:480`). A failed clear leaves a dangling pointer with no diagnostic anywhere — the exact opposite of the wipe path, which treats a pointer-write failure as a failed wipe (§5). The consequence of a dangling pointer is one failed selection at startup, not data loss.

---

## 5. Wipe interactions

**O7-144 — Staging creates the pointer's only window of naming a conversation that still exists elsewhere.** `WipeAllData` renames `chats/` to a sibling path, then resets the pointer, then clears inventory, then deletes the stages (`app-wipe.go:47-72`). During the window, and permanently if cleanup fails, **the conversations' data is outside the history root while the pointer may still name one of them.** A scan at that moment sees them absent.

Reconciliation obligations:

| Situation | Must report | Must not report |
|---|---|---|
| scan lands inside the staging window | `Unavailable` for the named conversation, with the reason that resolution found no directory — plus the scan's own uncertainty if it cannot distinguish stages from absence | `Missing`; a scan error promoted to absence |
| cleanup failed, stages remain | the staged remnants as **unreconciled**, named by path, outside the history boundary | `Missing`; `Deleted`; a clean wipe |

This is a scan-boundary obligation rather than a pointer rule, but the pointer is what makes the window visible: it names a conversation the root no longer contains.

**O7-145 — Pointer restoration precedes root restoration, so a failed restore can leave it dangling.** `app-wipe.go:62` clears, `:71` restores the old value if inventory clearing fails, and `restoreStagedRoots` runs after. If both root restores fail, the pointer names a conversation that no longer sits under `chats/`. That is correctly reported as `Unavailable` (no directory); the dangling pointer must not be reported as conversation damage, and must not be "repaired" by the scanner.

---

## 6. Verdict: which class of record is it?

**O7-146 — Three classes, and the pointer is being written like the wrong one:**

| Class | Written how | Recoverable | Coordination |
|---|---|---|---|
| authoritative representation (content) | atomic, CAS on `(state, generation)` | no | per-conversation lock |
| derived registration (cache, index) | atomic, rebuildable | yes, by recomputation | none |
| **selection hint** | atomic | yes, by re-selection | none |

`LastChatID` is currently a hint written like an authoritative record: on the hot path, as a side effect of content persistence, with the weakest write in the package (`O7-133`).

**O7-147 — Recommended classification: selection hint.** Exactly one writer — an explicit selection operation. Atomic replace by temp-and-rename. **No CAS, no lock, no generation, loss tolerated.**

The justification is not convenience, it is that the damage model has no room for coordination: the worst outcome of a lost, torn or raced value is that the next launch opens a different conversation than the user expected, and the user re-selects. A lock or CAS would serialise a decision whose failure cost is a click. This is also why the pointer must not be added to the gate: a gate exists to prevent content damage, and nothing here can damage content.

**O7-148 — The classification is unsound until the pointer stops being a proxy for process-local state.** Two dependencies, stated as dependencies rather than deferred:

1. `App` must hold its own active conversation identifier, bound at creation and at load, and every write path in §2 must use **that**, not a re-read of a shared file.
2. `SaveChat` must stop writing the pointer (`O7-87`, `O7-132`); the file is written on selection and on creation only.

Point 1 is what removes the concurrent-instance and shutdown-targeting classes of failure: the process already knows which conversation it holds, and reading it back from a contested register is a round trip that buys nothing and loses correctness. Point 2 is what stops selection and persistence sharing a field. **Neither is authorized here**; the classification is recorded as conditional on them, not as valid without them.

**O7-149 — The pointer must never enter the six-state reachability.** No conversation's state may depend on it, in either direction: not `Healthy` because it is named, not `Missing` because it is not, not `Unavailable` because it dangles, not `Pending` because it is ambiguous. Its absence from a scan is not evidence about any conversation, and a scanner that reported one would be converting a selection defect into a content verdict — the same class of error as promoting a scan error to `Missing`.

**O7-150 — Two selectors must not become two authorities.** The frontend copy is a presentation concern and a content-competition problem already recorded (row 523), not a second voice on selection for reconciliation. Until the ledger exists (row 524), selection across the boundary is not reconciled, and the contract must not imply that it is.

**O7-151 — What this does not answer.** Which conversation *should* be selected after its predecessor is deleted remains a product rule (FC-1's open item, row 523), not a consequence of anything here. Roaming across machines is OQ-7c-i's limit. Nothing here authorizes a pointer change, an `activeChatID` field, a `SaveChat` signature change, a `config.json` write change, a wipe change, a lock, a gate, CAS, a migration, or deletion enablement.

---

## 7. Tests this would require

- **O7-T37** Two instances, each selecting a different conversation, leave the file naming the later writer — and neither instance's *content* is altered by the disagreement.
- **O7-T38** A save in an instance holding a stale selection repoints the file to the stale conversation; the next startup restores that one; the other conversation is byte-identical before and after.
- **O7-T39** A torn `config.json` makes `initialize` fail cleanly and leaves every conversation directory untouched and loadable by a fresh manager.
- **O7-T40** An unreadable `meta.json` for the pointer's target produces a startup restore that re-selects and repoints — and a scan of that conversation still reports `Unavailable`, not `Healthy`, after the repoint.
- **O7-T41** A scan inside the wipe staging window reports the staged conversations as unreconciled, never `Missing`.
- **O7-T42** A failed pointer clear after `DeleteChat` leaves a dangling pointer that surfaces as a failed selection, and no conversation is reported damaged as a result.
- **O7-T43** No reachability state changes when `config.json` is deleted, emptied, or set to a nonexistent ID, for every conversation in the root.
- **O7-T44** `GetActiveChatMeta` backfills an empty `meta.ID` from the pointer, and a scan of the same conversation still reports identity unestablished.
- **O7-T45** The pointer path appears in no lock path and no gate transition — asserted structurally, not by convention.
