# TASK-03 — gate audit against the mutation call graph

| Field | Value |
|---|---|
| Purpose | Test the row 555 gate against the actual mutation call graph, not the five named paths |
| Status | **Amendment to row 555. Adds V2-67…V2-78. Contract work only.** |
| Supersedes | row 547 RC-27 (concurrent instances); row 555's implicit "five paths" inventory |
| Evidence | source reads only; no probe, no writes |
| Recovery token | `TASK-03-GATE-CALLGRAPH-2026-11-22` |

---

## 1. What was audited

Every function that writes inside a conversation directory, every entry point that reaches one, and every serialization mechanism between them. Not the five paths of row 549 — the call graph.

---

## 2. Findings

**V2-67 — There is no serialization at the mutation boundary. (Source-derived.)** `PersistenceManager` has **no mutex, no lock file, and no atomic check-and-write** (whole-file method list, `persistence.go`). `Service.mu` in `service.go:53` is an `RWMutex` over in-memory service state — inference context, model handles, lore and memory caches — and does not cover persistence.

So even a perfect caller-side check would race: read `(committed,G)`, another operation records a clear, the first write lands. That is the exact sequence row 555 §1 warned about, and nothing in the code prevents it.

**V2-68 — The five paths are an undercount.** Four writers inside persistence are reached from **at least twelve** entry points:

| Writer | Reached from |
|---|---|
| `SaveChat` | `NewChat` (`persistence.go:178`), `SaveChatAs` (`:463`), `RenameChat` (`:728`), `performSave` (`app-chat.go:314`), `SaveChatData` (`app-chat.go:68`), character edits (`app-characters.go:340`) |
| `AppendHistoryEntry` | `AppendChatMessage` (`app-prompt.go:310`) — called for the user turn, the assistant turn (`app-prompt.go:138`) and audio transcripts |
| `RewriteHistory` | `ReplaceChatHistory` (`:389`), `ClearChatHistory` (`:406`), `performFinalSave` (`app-chat.go:348`) |
| `os.RemoveAll` | `DeleteChat` (`persistence.go:471`) |

**V2-69 — Three writers sit entirely outside `PersistenceManager`: (Source-derived.)**

| Writer | File | Path |
|---|---|---|
| `MemoryStorage.saveMemories` | `memories.json`, `memories_backup.json` | `internal/backend/memory/service.go:363,452` |
| auto-capture settings | `auto_capture_config.json` | `internal/backend/memory/service.go:298` |
| `SaveContinuityConfig` | `continuity_config.json` | `internal/backend/memory/continuity.go:161` |

`SaveChat` also writes `memories.json` (`persistence.go:240`). **Two independent writers own that file**, with no coordination and last-writer-wins. This is not a history representation, so it is outside the gate's scope — but a gate that claims to cover "the conversation directory" would be claiming something untrue.

**V2-70 — A read path mutates.** `App.LoadChat` saves the *previous* conversation before switching (`app-chat.go:95-99`). A gate applied only to explicit mutations would be bypassed by navigating the chat list.

**V2-71 — A root-level operation bypasses every per-conversation rule.** `WipeAllData` renames the entire `chats` directory (`app-wipe.go:57-70`). No per-conversation gate can observe or prevent that.

**V2-72 — Cross-process races are currently permitted by the contract.** Row 547 RC-27 said coordination between instances is not required "because every state is derived from observation". That held when nothing was written. With a marker and a gate, two instances can interleave a check and a write. **Superseded.**

---

## 3. Verdict on row 555's gate

**V2-73 — The gate as written cannot be enforced at the current boundary.** It says mutations "must complete or refuse", which a caller can satisfy by checking and then writing. That is not enforcement; it is a convention with a race.

**V2-74 — The check and the write must be inseparable.** One critical section performs all three: read and validate the `(state, generation)` pair, confirm the mutation is permitted, and perform the write — under the same hold. No preflight call, no UI-level check, no startup check.

**V2-75 — One enforcement point, or none.** Every mutation entry point in §2 must pass through the same boundary. A boundary with two entrances is two boundaries.

**V2-76 — Serialization must be visible to a second instance.** Whatever guards the critical section must guard it across processes, because two instances are permitted and the store already has the discipline to copy (`DurableStore`'s exclusive lock). An in-process mutex alone leaves the cross-instance race intact.

**V2-77 — Writers outside the boundary are declared, not absorbed.** `memories.json` / `auto_capture_config.json` / `continuity_config.json` are not history representations; the gate does not claim them, and reconciliation's guarantees do not extend to them. Their own last-writer-wins conflict is a **separate defect** and is not fixed here.

**V2-78 — The gate cannot cover a rename of its own root.** Root-level operations are excluded from reconciliation's guarantees, and must be named as such rather than silently assumed safe.

---

## 4. The three requirements preserved

**V2-79 — Clear retry is idempotent.** An interrupted clear that replaced the log but not the snapshot is completed by re-running the same clear. Retrying must not treat the empty log as a new conversation, and must not permit an append in between. `RewriteHistory` is a full replace, so retry is safe; the gate is what makes the interval closed.

**V2-80 — `cleared` alone is not resolution.** Resolution is marker `cleared` **and** both representations structurally valid **and** both empty (row 555 V2-52, restated). A populated snapshot is never dismissed as harmless residue because the log is empty — that is exactly the state row 552 §2 case 9 showed resurrecting deleted messages.

**V2-81 — Inventory maintenance rule.** The gate's sufficiency depends on the entry-point inventory staying complete. Implementation tests must **enumerate mutation entry points and assert each is gated**, and a new mutation path must add a coverage test. A path that bypasses the gate silently reintroduces the defect that probe 12 reproduced.

---

## 5. What this changes in the contract's claims

| Claim | Before | Now |
|---|---|---|
| "Mutations complete or refuse" | sufficient | **not sufficient**; needs atomicity and a single boundary |
| "Five mutation paths" | the inventory | **twelve entry points, four writers, three external writers, one root operation** |
| "Concurrent instances need no coordination" (RC-27) | permitted | **superseded** — the gate requires cross-process serialization |
| "The scanner observes the conversation directory" | true | true for history representations; **not** for memory/config files or root renames |

---

## 6. Tests added

- **VT-29** A mutation whose check succeeded is refused if a transition became unresolved before the write (interleaved, deterministic via a test hook).
- **VT-30** There is exactly one enforcement point; each of the twelve entry points is asserted to pass through it.
- **VT-31** A second process attempting a mutation during an unresolved transition is refused, not merely slow.
- **VT-32** Retried clear after an interrupted one is idempotent, and no append is accepted between the two attempts.
- **VT-33** `cleared` with a populated snapshot is unresolved, even when the log is empty.
- **VT-34** `App.LoadChat`'s save of the previous conversation is gated (V2-70).
- **VT-35** `WipeAllData` is outside the gate and the contract says so; the test asserts the documented exclusion rather than silence.
- **VT-36** A new mutation path added without a coverage test fails the inventory test (meta-test, deliberately).

---

## 7. Unestablished

| # | Item |
|---|---|
| G-1 | Whether `projectRoot/Chats` (memory, continuity) and `<rootDir>/chats` (persistence) are the same directory in practice. They differ in case; on NTFS that is likely identical, on a case-sensitive volume it is two trees. **Source-derived divergence, not resolved here.** |
| G-2 | Whether `memories_backup.json` being overwritten with `{}` before every memory save is intentional. Noticed while auditing writers; not investigated. |
| G-3 | Whether the twelve entry points are a closed inventory for *history*; the search covered the repository, but a future writer would not appear. |

**V2-82 — Contract work only.** No gate, no lock, no marker, no recovery implemented. Row 555 stands as the amendment; this document amends it in turn, and both remain drafts.
