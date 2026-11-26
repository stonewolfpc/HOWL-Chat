# TASK-03 — the memories limitation

| Field | Value |
|---|---|
| Purpose | Establish what memory data exists, who owns it, where it lives, whether reconciliation may assess it, what delete/wipe/move does to it, and what the product must disclose when conversation content is healthy but memories are missing, unreadable or unassessed |
| Status | **Read-only investigation. Outcome: Clarify. V3.1 unedited. No memory state, lifecycle, repair or implementation authorized.** |
| Method | Source inspection of both memory stores, their schemas, writers and readers; delete, wipe, rename, copy and launch-root paths; row 573's shadow-tree scope; row 581's verdict meaning |
| Recovery token | `TASK-03-MEMORIES-LIMITATION-2026-11-22` |

The sixth and final baseline question from row 576, after OQ-1 (row 577), OQ-2 (row 578), OQ-3b-b (row 579), row 550 (row 580) and `V3-44` (row 581).

---

## 1. Answer

**ML-01 — Scope.** This record answers: which memory data is conversation-owned versus global or application-owned; which memory files sit inside the conversation candidate directory and which in the excluded shadow namespace; whether the scanner can assess memory integrity without becoming a memory authority; whether delete, wipe or move can leave memory data behind or make it unreachable; and what the product must communicate when conversation content is healthy but associated memories are missing, inaccessible or unassessed.

**ML-02 — Answer: there are two different memory stores, and the conversation verdict belongs to neither.** A conversation-owned companion `memories.json` sits inside the candidate directory but is not assessed and is effectively dormant. The live, user-authored memory records sit in a conversation-keyed shadow namespace outside the scan, survive deletion and wipe, and are not copied with a conversation. **Conversation health must never imply memory health.** The product may disclose memory association facts, but those facts are not conversation states and must not enter the six-state machine.

**ML-03 — Outcome: Clarify.** No `V3-39` row changes. V3.2 gains three things: an explicit scope clause that memories are outside the verdict, a non-implication clause, and a separate memory-disclosure clause whose vocabulary is observations, not states. `V3-73` closes the memories item.

---

## 2. What memory data exists

**ML-04 — Two files share the name `memories.json` and share nothing else.**

| Artifact | Path | Schema and writer | Owner and lifetime | In the candidate directory? |
|---|---|---|---|---|
| **Conversation companion** | `<pm.RootDir()>/chats/<id>/memories.json` | `[]interface{}` from `ChatData.Memories`; written by `NewChat` (as `[]`), `SaveChat`, `SaveChatAs` | The conversation; deleted with it | **yes**, but excluded from assessment |
| **Live memory records** | `StateRoot()/Chats/<id>/memories.json` | `[]types.MemoryRecord`; written by `AddMemory`, `EditMemory`, `DeleteMemory`, continuity import and the memory service's `SaveChatData` | The user; **outlives the conversation** | **no** |
| **"Backup"** | `StateRoot()/Chats/<id>/memories_backup.json` | literal `{}`, rewritten before every save; **no reader** | none — not a backup | no |
| **Prompt snapshot** | `StateRoot()/Chats/<id>/memories_snapshot.json` | read at session init; **no writer found** in this tree | legacy/test artifact | no |
| **Capture settings** | `StateRoot()/Chats/<id>/auto_capture_config.json` | `AutoCaptureConfig`; read/written by `App` bindings | application-owned setting | no |
| **Continuity settings** | `StateRoot()/Chats/<id>/continuity_config.json` | read; no writer found | application-owned setting | no |
| **Candidates** | `StateRoot()/Chats/<id>/candidates.json` | derived cache | application-owned cache | no |
| **Generated lore** | content folders (`Characters/.../automated_character_lore.json`, world/scenario lore) | lore pipeline | global/character content, not chat memory | no |

**ML-05 — The conversation companion is conversation-owned and dormant.** `NewChat` initialises `Memories: []interface{}{}` (`persistence.go:174`); `SaveChat` writes it (`:239-241`); `SaveChatAs` copies `currentChat.Memories` (`:460`). No production writer populates it from the memory service: the frontend's `collectForSave` payload has no `memories` key (`chat-state-bridge.js:525-560`), and `App.SaveChatData` only assigns `chatData.Memories` if such a key is present (`app-chat.go:609-610`). The file is therefore preserved and copied, but not a live view of user memories.

**ML-06 — The live store is conversation-keyed, not conversation-owned.** `MemoryRecord` has no chat ID (`types/memory.go:22-38`); the directory is the entire association. Its `Scope` field (`global`/`world`/`scenario`/`character`) is a retrieval and ownership label, not a storage location, and `OwnerKind` may say `chat` without changing the path. The records are user-authored and live in `StateRoot()/Chats/<chatID>/` (`memory/service.go:33-55`, `:362-380`).

**ML-07 — The "backup" cannot restore anything.** `saveMemories` writes the literal bytes `{}` to `memories_backup.json` and then overwrites `memories.json` (`memory/service.go:455-472`). Nothing reads the backup (no reader in the tree). It is a placeholder wearing a recovery name; treating it as a recovery source would be a memory-domain bug.

**ML-08 — The prompt snapshot is read but not written.** `LoadChat`'s service reads `StateRoot()/Chats/<chatName>/memories_snapshot.json` at session init (`chat/service.go:200`); no production writer was found. Row 573 classifies it as test-only; this record does not rely on that classification beyond "not an authoritative memory store".

**ML-09 — There is no global memory store.** A record with `Scope: global` still lives inside one chat's shadow directory. It is not reachable from another chat, and it is not deleted when the chat is deleted, because the shadow directory is not deleted. "Global" is a label on a record, not a location.

**ML-10 — The two files must never be conflated.** Same filename, different schema, different owner, different lifetime, different root, one parsed by the scanner never and one parsed only by the memory service. A report that says "the conversation's memories" without naming which file is ambiguous by construction.

---

## 3. Inside the candidate, outside the scan

**ML-11 — The companion is inside the candidate directory and outside the assessed scope.** Row 581 defines the assessed scope as the two history representations plus the marker; companions are excluded. `memories.json` is therefore physically present in the conversation directory and evidentially absent from the verdict.

**ML-12 — The live store is outside the boundary and cannot become a candidate.** Row 573 fixes the boundary at `<rootDir>/chats` and excludes `StateRoot()/Chats`. No shadow artifact carries `meta.json`, so candidacy by identity document keeps it out structurally, not by a path rule (row 573, `O7-189`).

**ML-13 — Association is a string match and nothing stronger.** The conversation's ID names the shadow directory. A copied conversation with a new ID has no association; a collision that reuses an ID shares one. There is no record-level chat binding and no durable link to verify.

---

## 4. What delete, wipe and move do

**ML-14 — Delete removes the conversation and the companion; the live records survive.** `DeleteChat` removes `<rootDir>/chats/<id>` only (`persistence.go:471-483`). The shadow directory, its `memories.json` and its settings persist indefinitely, keyed by an ID that no longer names a conversation. That is an orphan, not a conversation (row 573, `O7-185`).

**ML-15 — The memory service keeps the deleted chat's ID.** `App.DeleteChat` resets in-memory chat state but never calls `SetChatID` (`app-chat.go:117-129`). The service therefore remains scoped to the deleted ID, and the next memory write recreates its shadow directory.

**ML-16 — Wipe removes the conversations and leaves every memory.** `WipeAllData` stages `<RootDir>/chats` and the gallery (`app-wipe.go:45-52`); it does not touch `StateRoot()/Chats`. It calls `SetChatID("")` (`:85`), which returns early on an empty value (`memory/service.go:33-38`), so the last chat ID is retained and the next write recreates the shadow directory (row 573, `O7-183`/`O7-184`).

**ML-17 — Rename is safe for association.** `RenameChat` changes the display name and preserves the ID (`persistence.go:719-729`), so the shadow path is unchanged.

**ML-18 — Copy, move and launch-root change break association.** A manual copy of the conversation directory does not carry the shadow store; moving the conversation to a different history root leaves memories behind. Because `StateRoot()` resolves from the executable location or working directory (row 573, `O7-195`), launching the same binary from another directory points at a different shadow namespace while conversations stay in the same app-data root. The memories are not deleted; they are unreachable from the current scope.

**ML-19 — `SaveChatAs` creates a conversation with no memories, or shares them.** The copy gets a new ID and copies only the dormant companion (`persistence.go:437-468`); the shadow store is not copied. If creation collides and the ID is reused (row 580), the "copy" shares the source's shadow memory store, and both conversations read and write the same records.

**ML-20 — Absence at one scope is not deletion.** A missing file under the current `StateRoot()` may mean never created, deleted, moved, or a different launch root. The memory service itself conflates the first three: `loadMemories` returns an empty slice when the file is absent (`memory/service.go:362-371`). No evidence in the file can distinguish them.

---

## 5. Can the scanner assess memory integrity without becoming an authority?

**ML-21 — It can observe presence, readability and parseability at a named root; it cannot assess integrity.** Integrity would include completeness, authorship, ownership, correct association, backup validity and recoverability. None is observable: there is no record count to compare, no authoritative snapshot, no per-record chat binding, no valid backup, and absence is indistinguishable from emptiness (ML-20).

**ML-22 — It cannot say "there are no memories".** The service returns `[]` for a missing file, so a zero-length list is not evidence of zero memories. A report must distinguish "not found at this scope" from "empty".

**ML-23 — It cannot repair, and must not choose.** No write, recreate, delete, merge, dedup, restore-from-backup, precedence between `memories.json` and `memories_backup.json`, or reassignment of records. The backup is `{}` (ML-07); choosing it would destroy data.

**ML-24 — Reading the shadow store would extend the boundary, so the conversation verdict must not do it.** Row 573 keeps the scan at the history root. If the product needs association facts, they belong to a **separate memory-disclosure report** that does not feed any row and does not make the conversation verdict wider.

**ML-25 — That separation is what prevents authority.** The scanner may report what it saw at a named scope; it may not decide whether memory data is valid, complete, correctly owned, recoverable, or safe to delete. The moment it decides those, it owns memory lifecycle, which this contract does not.

---

## 6. What the product must communicate

**ML-26 — Disclosure uses observations, not states.** Five findings, none of which is a lifecycle:

| Finding | Meaning | Must not be read as |
|---|---|---|
| `unassessed` | no memory observation was attempted | "no memories" |
| `not_found_at_scope` | no memory file at the named root and ID | "deleted", "empty", "lost" |
| `unreadable` | a file exists but cannot be read or parsed | "empty", "damaged beyond recovery" |
| `observed` | a file exists and parses as a memory array | "complete", "intact", "correct" |
| `orphaned` | memory data exists for an ID with no conversation | a conversation state, or a deletion prompt |

**ML-27 — `observed` is deliberately weak.** It says the store was seen and parsed. It does not say the records are complete, that the backup is usable, that the records belong to this conversation, or that the user authored all of them.

**ML-28 — `not_found_at_scope` must name its scope.** The report must carry the root or scope it looked under, and the product must phrase it as "not found at the current memory location", never "memories deleted" or "no memories". Launch-root changes make this distinction load-bearing (ML-18, ML-20).

**ML-29 — `unreadable` is not `empty`.** A parse or I/O failure must be shown as inaccessible, never as a zero count. The service's silent `[]` for absence must not be reused for failure.

**ML-30 — `orphaned` is reported separately and never as a conversation.** It is memory data without a conversation: not `Missing`, not `Damaged`, not `Pending`, no verdict, no delete suggestion. The same rule row 573 gives shadow directories (`O7-190`).

**ML-31 — The conversation verdict always carries the scope statement.** Every `Healthy`, `Pending`, `Damaged`, `Unavailable`, `Missing` or `Orphaned` verdict says memories were **not assessed**. Row 581 already requires this for `Healthy`; the same sentence is mandatory for the others, because silence is read as coverage.

**ML-32 — No aggregation.** A product must not combine a `Healthy` conversation with an `observed` memory file into "all data healthy". There is no joint verdict, no combined badge, and no API field that means "conversation and memories are fine".

**ML-33 — Memory lifecycle and recovery are a separate ownership problem.** If memories need their own states, repair, backup, migration or recovery policy, that is a memory-domain design with its own evidence and authorization. It is not obtained by adding states to this machine.

---

## 7. The amendment

**ML-34 — Scope clause:** a conversation verdict is a statement about the candidate directory's two history representations and its marker; the conversation companion `memories.json` and the entire `StateRoot()/Chats` namespace are outside it, whatever they contain.

**ML-35 — Non-implication clause:** conversation health never implies memory health, memory presence, memory completeness, correct memory association, or that any memory data was assessed.

**ML-36 — Disclosure clause:** memory association facts are reported by a separate, read-only disclosure that names its scope and uses the five observations of ML-26; it feeds no row, adds no state, and performs no write, repair, restore or deletion.

**ML-37 — `V3-73` closes the memories item.** No `V3-39` row, `V3-40` state, reason or transition changes; `V3-38`'s ownership boundary already assigns creation, loader and rewrite defects elsewhere, and memories join that list.

**ML-38 — Nothing is authorized, and the limitation survives a rewrite.** No memory state, lifecycle, repair, backup restoration, migration, relocation, wipe/delete change, `SetChatID` change, dedup, orphan cleanup or path consolidation. If persistence is replaced, conversation health still may not imply memory health; the two stores' separation is a statement about evidence and ownership, not about the current two-file layout.

---

## 8. Outcome

**ML-39 — Outcome: Clarify.** Conversation health never implies memory health. Memories are conversation-keyed but not conversation-owned; the live store is outside the scan and outlives deletion; the companion inside the directory is dormant and unassessed; absence is not deletion; and the product discloses memory association as observations, never as conversation states. Three clauses plus the `V3-73` closure; no state machine change.

---

## 9. Tests this would require

- **ML-T1** A conversation with `Healthy` content and a shadow memory file reports `Healthy` plus `memories: unassessed` on the verdict, and a separate `observed` disclosure that names its scope.
- **ML-T2** `DeleteChat` removes the candidate directory and leaves `StateRoot()/Chats/<id>/memories.json` intact; the disclosure reports `orphaned` for that ID, never a conversation verdict.
- **ML-T3** After `WipeAllData`, the shadow namespace is unchanged and the retained chat ID is shown to recreate its directory on the next memory write.
- **ML-T4** A missing memory file reports `not_found_at_scope`, never `empty`, and the report names the root it looked under.
- **ML-T5** A malformed memory file reports `unreadable`, never a zero count and never a repair attempt.
- **ML-T6** A same-second `SaveChatAs` collision makes both conversations read the same shadow store; the disclosure does not claim separate memory ownership.
- **ML-T7** The conversation state machine still has exactly six states, and no memory finding can produce or change one.
- **ML-T8** The negative assertion: no disclosure path writes, recreates, deletes, dedups, restores from `memories_backup.json`, or chooses between memory files.

---

## 10. Nothing here is authorized

**No implementation follows.** This record authorizes no memory state, lifecycle, recovery, backup restoration, dedup, migration, orphan cleanup, path change, wipe/delete change, `SetChatID` fix, disclosure UI, API field or scanner extension. It states what the conversation verdict may not imply and what a separate disclosure may observe.
