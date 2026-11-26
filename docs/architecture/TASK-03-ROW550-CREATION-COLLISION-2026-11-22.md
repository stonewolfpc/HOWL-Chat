# TASK-03 — row 550: creation collision

| Field | Value |
|---|---|
| Purpose | Establish whether creation can silently adopt, merge with, overwrite or redirect writes into an existing conversation; whether a unique ID plus exclusive creation is sufficient; and what reconciliation may say about the result |
| Status | **Read-only investigation. Outcome: Clarify (contract) plus a creation-side prerequisite. V3.1 unedited. No repair authorized.** |
| Method | Source inspection of every creation-capable path and its callers; one disposable probe reproducing cross-process creation in a forced second, a same-second `SaveChatAs`, and a creation onto a pre-existing conversation (`build/tmp/collprobe`, log `build/tmp/collprobe.log`) |
| Recovery token | `TASK-03-ROW550-CREATION-COLLISION-2026-11-22` |

Fourth of row 576's six baseline questions, after OQ-1 (row 577), OQ-2 (row 578) and OQ-3b-b (row 579). It revisits row 550's defect record and its proposed TASK-04, without authorizing either.

---

## 1. Answer in brief

**CC-01 — Scope.** This record answers five questions: can collision occur across processes; which creation paths are vulnerable; can a caller mistake an existing conversation for a new one; can collision overwrite or merge user-authored content; and is a unique logical ID plus exclusive creation sufficient, or does the workflow need a transactional design.

**CC-02 — Answer: all four damaging outcomes are reachable today, and the contract cannot classify the result as a distinct condition.** Creation can silently adopt an existing directory, overwrite its metadata and companions (including user-authored memories), leave its log in place, and hand the caller a new conversation that is really the old one. Reconciliation may still report a valid conversation — the damage is to identity and ownership, not necessarily to availability — so forcing `Damaged` would be wrong.

**CC-03 — Outcome.** The contract change is **Clarify**: add an identity-integrity limitation and close `V3-73`'s row-550 effect without adding a state or a reason. The repair is a **creation-side prerequisite** (unique ID, exclusive claim, no implicit creation elsewhere, commit ordering), which is row 550's proposed TASK-04 and remains unauthorized.

**CC-04 — A correction to row 550's own remedy.** `RC5-7` says a scan reports a collision-shaped directory as an **anomaly**. The probe shows the surviving shape is often a coherent conversation: valid metadata plus a valid log, or a renamed source with its original log. There is no content evidence that names the collision reliably. `RC5-8`–`RC5-10` stand; `RC5-7`'s anomaly is not implementable as stated and is superseded by CC-33.

---

## 2. Which paths create or adopt a conversation

**CC-05 — Observed: the ID is wall-clock to the second.** `generateChatID` formats UTC as `YYYY-MM-DDTHH-MM-SSZ` (`persistence.go:134-141`). Two creations in one second, in one process or across processes, compute the same string.

**CC-06 — Observed: `NewChat` claims nothing exclusively.** It builds `chats/<timestamp>`, calls `os.MkdirAll` — which succeeds on an existing directory — and then `SaveChat` writes metadata and companions over whatever is there (`persistence.go:143-186`). There is no existence check, no `Mkdir`-fails-if-present, no read-back of an existing `meta.json`, and no comparison with the returned ID.

**CC-07 — Observed: `SaveChatAs` inherits the same behaviour and then copies into it.** It loads the source, calls `NewChat` for the destination, copies the data and saves (`persistence.go:437-468`). If the second has not advanced, the destination ID equals the source ID, so the copy is written into the source directory.

**CC-08 — Observed: two further paths create implicitly.** `SaveChat` calls `MkdirAll(chats/<Meta.ID>)` before writing (`:190-200`), and `RewriteHistory` calls `MkdirAll(chats/<chatID>)` before rewriting (`:631-640`). Both accept an arbitrary ID and will create or adopt a directory for it. Creation is therefore not confined to `NewChat`.

**CC-09 — Observed: `AppendHistoryEntry` is not a creation path.** It opens `history.ndjson` with `O_CREATE` but does not create the conversation directory (`:668-702`); if the directory is missing, the open fails.

**CC-10 — Observed: production callers are narrow, but the vulnerable primitive is general.** The only production creation calls are `App.NewChat` → `NewChat` (`app-chat.go:48`) and `App.SaveChatAs` → `SaveChatAs` (`:80`). However, `SaveChat` is reached from `NewChat`, `SaveChatAs`, `RenameChat`, `performSave`, `SaveChatData` and character edits, so any of those carrying an arbitrary `Meta.ID` can create or adopt.

**CC-11 — Table of the paths:**

| Path | Creates the directory | Checks existence | What it writes |
|---|---|---|---|
| `NewChat` | `MkdirAll` | **no** | meta, history.json, five companions, config |
| `SaveChatAs` | via `NewChat` | **no** | copy of source data, then meta/companions |
| `SaveChat` | `MkdirAll` | **no** | meta first, then history.json, companions, config |
| `RewriteHistory` | `MkdirAll` | **no** | log, then history.json |
| `AppendHistoryEntry` | no | directory must exist | one log line |

---

## 3. Observed: cross-process collision

**CC-12 — Method.** Two child processes of the probe were released at a forced UTC second against one shared root; each called `NewPersistenceManagerWithDir`, `NewChat(name)`, then appended one message. No product file was touched.

**CC-13 — Result (Observed):**

```
CHILD A newchat_err=failed to save new chat: failed to save meta.json:
  rename ...\2026-11-22T17-15-55Z\meta.json.tmp ...\meta.json:
  The process cannot access the file because it is being used by another process.
CHILD B id=2026-11-22T17-15-55Z newchat_err=<nil> append_err=<nil>
DIRS 1
DIR 2026-11-22T17-15-55Z meta.name=B files=[characters.json history.json history.ndjson
  lorebooks.json memories.json meta.json] texts=[from-B]
CONFIG {"lastChatId": "2026-11-22T17-15-55Z"}
```

**CC-14 — Calculated: both processes targeted the same directory, and the loser learned nothing useful.** A's failure is a sharing violation on the fixed temp file `meta.json.tmp`, not a statement that the ID was taken. A retry in the same second would collide again; a retry in the next second creates a second conversation. Neither path tells the caller that B's conversation exists.

**CC-15 — Observed: the temp-file name is fixed, so the race is on the representation as well as the identity.** `saveJSON` writes `filename + ".tmp"` and renames it (`:571-582`); two processes writing the same document use one temp path. `saveConfig` writes `config.json` directly with no temp (`:118-130`). The probe produced a sharing violation on `meta.json.tmp`; a corrupted or mixed document is **Unestablished** but the mechanism is present.

**CC-16 — Calculated: cross-process collision is production-reachable.** Two instances share the same history root (row 569) and no creation-time cross-process lock exists. Row 550's probe was in-process; this record shows the same second-resolution defect across processes.

---

## 4. Observed: same-process `SaveChatAs`

**CC-17 — Result (Observed):**

```
SAVEAS source_id=2026-11-22T17-15-56Z copy_id=2026-11-22T17-15-56Z same=true
DIRS 1
DIR ... meta.name=duplicate files=[... history.ndjson ...] texts=[original]
```

**CC-18 — Calculated: the "copy" is the source.** `SaveChatAs` returned the source's own ID with no error. The source directory's metadata was renamed to "duplicate", its `history.json` was replaced by the copy's snapshot, and its `history.ndjson` survived. The caller believes a new conversation exists; in fact one conversation was renamed and its snapshot overwritten.

---

## 5. Observed: a pre-existing conversation is adopted and its memories destroyed

**CC-19 — Method.** A conversation was created at the ID the next second would generate, with `meta.name = "existing"`, `memories.json = [{"text":"precious"}]`, and a one-record log. `NewChat("intruder")` was then called inside that second.

**CC-20 — Result (Observed):**

```
OVERWRITE target=2026-11-22T17-15-57Z returned=2026-11-22T17-15-57Z err=<nil> same=true
  meta.json       {"id":"2026-11-22T17-15-57Z","name":"intruder", ...}
  memories.json   []
  history.json    []
  history.ndjson  {"id":"0",...,"text":"old",...}
  loaded name=intruder texts=[old]
```

**CC-21 — Calculated: four things happen at once.** The existing identity is replaced; the user-authored memories are replaced by an empty array; the snapshot is replaced by an empty array; and the surviving log is adopted by the new conversation, so `LoadChat` returns the old message under the new name. The companions (`characters`, `lorebooks`, and the rest) are overwritten with empty values by the same `SaveChat` call.

**CC-22 — This is the data loss that matters most.** Row 566 and row 573 establish the chat-scoped `memories.json` as user-authored and irreplaceable. A single `NewChat` in a colliding second silently replaces it with `[]`. The log survives only because `SaveChat` does not write it; the next append merges the new conversation's messages into the old log.

**CC-23 — Observed: the caller cannot tell.** `NewChat` returns the in-memory `ChatData` it just built, with the new name and the colliding ID; it never reads the directory back, never compares an existing `meta.json`, and returns `nil`. `App.NewChat` then clears `a.chatMessages` and points the memory service at that ID (`app-chat.go:41-58`), so subsequent turns are written into the adopted directory.

---

## 6. Is a unique ID plus exclusive creation sufficient?

**CC-24 — Necessary: an ID that does not depend on clock resolution.** A UUID, ULID or random suffix removes the same-second collision by construction. The timestamp string is not unique across processes, and its one-second granularity is the direct cause of every observation above.

**CC-25 — Necessary: an exclusive claim on the final directory.** Creation must use a primitive that fails when the directory exists (`Mkdir`, not `MkdirAll` on the final path), or an equally strict no-replace rename. A claim that succeeds on an existing directory is adoption, whatever else changes.

**CC-26 — Necessary: no implicit creation in any other writer.** `SaveChat` and `RewriteHistory` must require an existing directory whose `meta.ID` matches the target (row 578's logical identity) and must refuse otherwise. Without this, a stale or wrong ID can still create or adopt a conversation after the creation path is fixed.

**CC-27 — Necessary: the commit point must be ordered or atomic.** `SaveChat` writes `meta.json` **first**, then history and companions (`:200-240`). Since `meta.json` is the candidacy marker (row 573), a crash mid-creation leaves a directory that already looks like a conversation. Creation should either write `meta.json` last as the commit record, or build in a private temp directory and move it into the claimed ID.

**CC-28 — Sufficient for the stated end result: yes, these four are sufficient.** Unique ID + exclusive claim + no implicit creation + safe commit ordering means creation can never silently adopt, merge with, overwrite or redirect into an existing conversation. A full transactional store is **not** required for that end result.

**CC-29 — What a transaction would add, and what it would not fix.** Temp-directory-then-rename gives all-or-nothing creation and simpler crash cleanup; it does not by itself prevent adoption unless the final move is exclusive, and it does not give unique IDs. Conversely, unique IDs plus exclusive `Mkdir` give collision safety without a transaction engine. The two solve different problems and should not be conflated.

**CC-30 — A creation lock is not the primary guard.** A cross-process lock would serialize same-ID attempts, but it cannot make a non-unique ID safe and it is unnecessary once IDs are unique and the claim is exclusive. Locks remain necessary for concurrent writes to an existing conversation (OQ-7b), which is a different operation.

**CC-31 — The pointer is not an ownership record.** `NewChat` and `SaveChat` both set `LastChatID` (`:180-185`, `:235-241`); two creations leave it at the last writer. Creation success must be reported from the claimed ID and the returned result, never reconstructed from `config.json` (row 571).

---

## 7. What reconciliation may say

**CC-32 — A collision is not necessarily `Damaged`.** It may leave a structurally valid, readable conversation whose identity or ownership is compromised. The probe's overwrite case loads cleanly: name "intruder", one old message. Calling that Damaged would assert content loss the scanner cannot see.

**CC-33 — V3.1 decides these inputs by its existing rows, and none of them names the collision.** Merged content with no marker is row 12 `Pending / unattributed`; an overwritten empty conversation is row 10 `Healthy / empty`; an overwritten meta with a surviving log and a settled marker can be rows 14–15 `Healthy`. The table is not wrong about availability; it is silent about creation integrity, and it cannot be made to see it from content.

**CC-34 — The limitation V3.2 needs.** `Healthy` and `Pending` assert content availability under the decision procedure. They do **not** assert that the directory was created exclusively, that its logical ID is unique, that its companions belong to the same creation as its log, or that its ownership is uncompromised. A collision is not a state; no anomaly state, reason or repair is added, and the scanner never renames, splits or deletes either party (`RC5-8`, `RC5-9`).

**CC-35 — `V3-73`'s row-550 item closes with no row change.** The effect on rows 6 and 12 is: they may be reached by a collision's residue, but they do not identify it. The limitation states the gap; the creation repair is the only thing that prevents it.

**CC-36 — The creation prerequisite, stated as an end result, not a design:** creating a conversation must never silently adopt, merge with, overwrite or redirect writes into an existing conversation. CC-24–CC-27 are the minimum conditions; how they are implemented (UUID vs ULID, `Mkdir` vs no-replace rename, meta-last vs temp-directory) is a separate authorized task.

---

## 8. Outcome

**CC-37 — Outcome: Clarify for the contract, with a prerequisite repair.** No `V3-39` row changes, no state and no reason are added. V3.2 gains the identity-integrity limitation of CC-34, and `V3-73` records the row-550 item as closed. The creation workflow needs redesign before the product can claim the end result, but that is a separate task and remains unauthorized.

**CC-38 — This is not an argument for the current layout.** If creation is rebuilt around a different store, the four conditions still apply: non-clock identity, exclusive claim, no implicit creation elsewhere, and a commit point that does not make a partial creation look complete.

---

## 9. Tests this would require

- **CC-T1** Two processes released in the same second: exactly one creation succeeds; the other fails with a collision-specific result, and the winner's directory is untouched by the loser.
- **CC-T2** A pre-existing conversation at the target ID: creation fails and the directory's `meta.json`, companions and log are byte-identical afterwards — the CC-20 memories loss cannot recur.
- **CC-T3** N creations in one second receive N distinct IDs.
- **CC-T4** `SaveChatAs` in the same second as its source returns a distinct ID and leaves the source directory byte-identical.
- **CC-T5** `SaveChat` and `RewriteHistory` on a missing directory refuse; no implicit creation.
- **CC-T6** A creation interrupted before commit leaves no conversation candidate (meta absent or still in a private temp directory).
- **CC-T7** A collision residue reports an ordinary verdict (row 10/12/14/15 as its evidence dictates) and is never forced to Damaged; no anomaly state is invented.
- **CC-T8** The negative assertions: no creation path uses `MkdirAll` on the final conversation directory; no ID is derived solely from wall-clock resolution; `LastChatID` is never used to establish that a creation succeeded.

---

## 10. Nothing here is authorized

**No implementation follows.** This record authorizes no ID change, no exclusive-create primitive, no `MkdirAll` removal, no meta-ordering change, no temp-directory creation, no creation lock, no config change, no repair, rename, split or deletion of a collided directory, no cleanup of incomplete creation, no new state or reason, and no adoption of the refusal-only retry policy. It states what creation must guarantee and what a verdict may not claim.
