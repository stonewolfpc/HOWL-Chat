# TASK-03 — OQ-7a: path centralisation

| Field | Value |
|---|---|
| Purpose | Can history writes be structurally confined to one package, or only hoped for? |
| Status | **Read-only investigation. Design only. V3.1 unchanged. Nothing authorized.** |
| Method | Inventory of every site that opens a history representation for writing, plus every path construction for a conversation directory |
| Recovery token | `TASK-03-OQ7A-PATH-CENTRALISATION-2026-11-22` |

Continues OQ-7's numbering (`O7-45`…). Includes one correction to row 564's ordering table (§1).

---

## 1. Correction first: a failure response is not a rollback

**O7-45 — Two refusal outcomes must not be conflated.** Row 564's ordering table listed them, but the difference is a consistency property, not a formatting one:

| Outcome | Consistency debt |
|---|---|
| refused **before** any content mutation | none — the store is exactly as it was |
| content mutated, marker recording failed | **the transition is now unresolved on disk** |

**O7-46 — A failure response does not restore consistency.** When content changed and the marker did not, the caller receives a failure *and* the conversation is left in a state reconciliation must see as unresolved (row 12, 18, 19, or 7 for a torn write). Reporting failure is the correct behaviour; it is not a repair.

**O7-47 — The boundary must not attempt compensating rollback.** A rollback is another mutation with its own failure mode, and it would rewrite content the marker does not describe. The resolution path is retry of the same idempotent operation under the gate (V3-33), not an undo.

**O7-48 — Consequence for the API:** the result of a partial transition must name the state it left behind, so the caller cannot present it as "nothing happened".

---

## 2. Inventory: every write-opening site (non-test)

**O7-49 — `history.ndjson`, three sites, all inside `internal/backend/chat/persistence.go`:**

| # | Site | Operation | Atomic? |
|---|---|---|---|
| 1 | `:653-656` `RewriteHistory` | temp write + rename over the log | yes |
| 2 | `:678` `AppendHistoryEntry` migration | `OpenFile(O_APPEND\|O_CREATE\|O_WRONLY)` backfilling the snapshot into a new log | **no**, and its error is discarded (`if openErr == nil`) |
| 3 | `:696` `AppendHistoryEntry` | `OpenFile(O_APPEND\|O_CREATE\|O_WRONLY)` | append |

**O7-50 — `history.json`, four sites, and one is in a different package:**

| # | Site | Operation | Atomic? |
|---|---|---|---|
| 4 | `:211` `SaveChat` → `saveJSON` | temp + rename | yes |
| 5 | `:661` `RewriteHistory` → `saveJSON` | temp + rename | yes |
| 6 | `service.go:536-542` **`Service.SaveChatData`** | **`os.WriteFile` directly** | **no** |

Plus the sites that reach these: `SaveChatAs` (`:463`), `RenameChat` (`:728`), `NewChat` (`:178`) all call `SaveChat`.

**O7-51 — Site 6 is not another route to the same writer; it is a different writer.** `Service.SaveChatData` serialises `s.context.GetMessages()` — the in-memory inference context — into `history.json`. Its content source is neither of the two on-disk representations. It writes non-atomically, so a crash mid-write leaves a truncated snapshot.

**O7-52 — Test-only sites also open the log directly** (`app_chat_character_test.go:112`, `history_character_test.go:80`). They must be inside the assertion's declared scope or the assertion fails on the suite itself; excluding them by accident would be the checklist error in reverse.

---

## 3. Path construction: three conventions, no constructor

**O7-53 — There is no single conversation-path constructor. Source-derived:**

| Convention | Where | Identity |
|---|---|---|
| `<rootDir>/chats/<chatID>` | `persistence.go` at nine sites (`:85, :149, :193, :257, :472, :487, :556, :636, :669`) | ID, lowercase |
| `<projectRoot>/Chats/<chatID>` | memory package at seven sites (`service.go:269,292,363,451`, `continuity.go:156,172`, `candidates.go:59,79`) | ID, **capital C** |
| `<projectRoot>/Chats/<chatName>` | `chat/service.go:200,518,536` | **name**, capital C |

**O7-54 — The third convention is the important one.** `Service.SaveChatData` addresses a conversation by *display name*, while persistence addresses it by *ID*. On a case-insensitive volume `Chats` and `chats` collapse into one directory, so the name-based writer and the ID-based writers can land in the same place while disagreeing about identity — and the project already has a documented creation-identity collision (row 550) in the same area.

**O7-55 — Therefore no assertion over paths can be a proof today.** A negative assertion is only meaningful if paths are produced in one place; three conventions and nine repeated literals mean the checker would be describing an accident of style, not a structure.

---

## 4. Verdict on the question asked

**O7-56 — Confinement to one package is currently false, in three independent ways:**

1. **a different package writes `history.json`** (O7-50 site 6), with a different content source and non-atomic writes;
2. **no central path constructor exists** (O7-53), so "the set of writing sites" is not enumerable from a single place;
3. **the identity basis differs** — ID versus name (O7-54), so even the same directory can be reached under two identities.

**O7-57 — The assertion is therefore a heuristic, and must be described as one.** Stated plainly, as requested: it would today catch a *new literal* write, and it would not catch a writer that assembles the name, nor would it establish that the known set is complete. The existing site 6 is precisely the class of thing a checklist of "known writers" would have omitted.

---

## 5. What would make it structural

**O7-58 — Three changes, none authorized:**

| Change | What it closes |
|---|---|
| one unexported `conversationDir(id)` constructor, with the two basenames as package constants | naming drift and the nine repeated literals |
| a `ConversationPath` type with an unexported field, constructible **only** by that constructor, accepted by the boundary API in place of a string | **dynamic assembly**: a runtime-built string cannot become a `ConversationPath` |
| `Service.SaveChatData`'s history write removed or routed through the boundary | the second writer, and its non-atomic path |

**O7-59 — With those, the negative assertion becomes two checkable facts**, and only then is it proof within its scope:

1. no literal `"history.ndjson"` / `"history.json"` outside the boundary package;
2. no boundary-side API accepts a raw `string` path where a `ConversationPath` is required — so a computed name has nowhere to enter.

**O7-60 — Even then, the claim stays scoped.** `memories.json`, `memories_snapshot.json`, `candidates.json`, `auto_capture_config.json`, `continuity_config.json` and the root wipe remain outside it (O7-06, OQ-7d). Their path conventions are the capital-C ones, so typing the history path does not type theirs.

---

## 6. Effect on the branch

**O7-61 — The enforcement claim cannot be certified today, at any confidence level.** Not because the lock is weak (it is adequate, O7-16) and not because the shape is wrong (it is provable, row 564), but because the *population* of writers is not enumerable. Certification requires O7-58.

**O7-62 — OQ-7d is now partly answered by evidence:** there are at least six non-history files written into conversation directories by three packages, under two case conventions. Their ownership is a separate lifecycle question and must not migrate into the history guarantee by adjacency.

**O7-63 — Nothing here is authorized.** No constructor, no path type, no constants, no removal or routing of `Service.SaveChatData`, no assertion, no lock composition, no marker change, no migration, no backfill, no deletion enablement.

---

## 7. Tests this would require

- **O7-T9** The negative assertion fails on a newly added literal history write outside the boundary package.
- **O7-T10** A raw string path cannot be passed where `ConversationPath` is required (fixture that must not compile).
- **O7-T11** Two identity conventions cannot address the same directory: a name-based and an ID-based construction either coincide demonstrably or the name path is rejected.
- **O7-T12** `Service.SaveChatData`'s history write no longer exists, or is exercised through the boundary (asserted, not assumed).
- **O7-T13** A partial transition reports the unresolved state it left behind, and no code path reports it as success (O7-46, O7-48).
- **O7-T14** The assertion's declared scope includes the test-only direct writers, so the suite does not fail the check it is meant to pass.
