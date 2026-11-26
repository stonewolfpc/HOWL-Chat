# TASK-03 — OQ-7c-iii: does anything require true re-entrancy?

| Field | Value |
|---|---|
| Purpose | Whether any product requirement needs recursive lock acquisition once OQ-7c's four decompositions are applied |
| Status | **Read-only investigation. Design only. V3.1 unchanged. Nothing authorized.** |
| Method | Call-graph inspection of every persistence entry point, every goroutine that reaches one, and every callback between writes |
| Recovery token | `TASK-03-OQ7CIII-REENTRANCY-2026-11-22` |

Continues OQ-7's numbering (`O7-152`…). Narrows OQ-7c-iii as scoped by the owner: synchronous nesting, async overlap, frontend selection change, native callbacks, and cross-conversation operations, kept apart.

---

## 1. Five situations that are routinely called "re-entrancy"

**O7-152 — Only one of the five is re-entrancy; the others need different guarantees:**

| Situation | Occurs today? | Is it re-entrancy? | Actual requirement |
|---|---|---|---|
| synchronous nested persistence | **yes** (§2) | nominally, but only through composition | one operation, one critical section |
| asynchronous overlap | **yes** (§3) | no | concurrency control **and** identity binding |
| frontend selection change mid-flight | **yes** (§3) | no | identity binding |
| native callbacks between writes | **no** (§4) | no | structural prohibition |
| cross-conversation operation | **yes** (§5) | no | deterministic order over distinct locks |

Conflating them is how a lock design acquires recursion it does not need. Each has its own answer below.

---

## 2. Synchronous nesting: real, and it already collapses

**O7-153 — Two composition sites nest same-conversation writes:**

| Operation | Inner calls | Writes to the same conversation |
|---|---|---|
| `NewChat` (`persistence.go:166-190`) | `SaveChat` | 1, then the pointer write |
| `SaveChatAs` (`persistence.go:437-466`) | `LoadChat`, `NewChat` → `SaveChat`, then `SaveChat` | **2 to the destination** — an empty document set, then the copied one |

`SaveChatAs` therefore writes the destination twice in one operation, with no lock and no intermediate consistency. If the process dies between them, the destination exists as an empty conversation carrying the new name and no history — indistinguishable, to a scanner, from a real conversation that lost its content. That is a distinct and narrower defect than row 570 (which is about the *source* representation order at shutdown), and it is recorded here as an observation, not as a repaired behaviour.

**O7-154 — Neither site requires recursive acquisition, because neither is a lock boundary in the target design.** OQ-7's containment shape is an unexported writer plus a single `Mutate` operation over a closed command set. Under that shape:

| Today | Under the boundary |
|---|---|
| `NewChat` → exported `SaveChat` | `Mutate(create)` → unexported writer, once |
| `SaveChatAs` → `NewChat` → `SaveChat`, `SaveChat` | `Mutate(copy)` → unexported writer, twice inside **one** critical section |

The nesting is a property of *composition within one operation*, and one operation is one acquisition by definition. A re-entrant lock would be the mechanism needed only if the exported writer stayed exposed and callable from inside itself — which the shape removes. **This is the decisive point: re-entrancy is not a requirement, it is a symptom of leaving multi-write operations as separately-lockable steps.**

---

## 3. Asynchronous overlap: the real problem, and it is not recursion

**O7-155 — There are two unjoined autosave goroutines, and they are the strongest evidence against a recursion design:**

| Site | Spawned where | Joined? |
|---|---|---|
| `app-chat.go:187-193` | end of `SendMessage` | **no** |
| `app-prompt.go:142-148` | after `chat:complete` emit | **no** |

Both call `a.performSave()` fire-and-forget, discarding the error into a log line. Nothing joins them, including at shutdown: `ensureAcceptingWork` guards new *entry points*, and a goroutine already in flight is not one. So **`performFinalSave` can run concurrently with an autosave of the same conversation**, which is the row 570 ordering defect with a second writer added.

**O7-156 — And they expose the identity defect OQ-7c-ii predicted, in shipped code.** `performSave` reads `GetLastChatID()` *at execution time* (`app-chat.go:315-319`), not at launch:

| Step | State |
|---|---|
| user is in X, sends a message | autosave goroutine spawned for X's content in `a.chatMessages` |
| user switches to Y before the goroutine runs | pointer is now Y |
| goroutine executes | loads **Y**, overwrites it with `a.chatMessages` — X's content |
| X's in-memory content | written into Y's documents |

This is the general rule the owner asked to make explicit, and it belongs in the contract as a requirement on the boundary, not on the pointer:

**O7-157 — Every mutation must use the conversation identity bound to the content it is mutating, not the identity current when the operation executes.** A deferred or asynchronous operation carries its conversation ID from the moment the content was captured; if that conversation is no longer the one bound to the state being written, the operation **refuses** rather than retargeting. Retargeting is the failure in O7-156, and no lock prevents it — the second writer acquires the wrong conversation's lock perfectly legitimately.

**O7-158 — Overlap is therefore answered by two things, neither of which is recursion:** a single writer per critical section (the lock), and a bound identity per operation (O7-157). Where two operations of the same conversation overlap, the later one waits or is refused under the bounded-refusal rule (O7-165); it never proceeds on the assumption that it holds the lock already.

---

## 4. Native callbacks: no persistence edge exists, and none may be added

**O7-159 — No callback in the codebase reaches a write.** The streaming token callbacks write to a `thinkSplitter` and emit display events (`app-prompt.go:104-118`, `:219-239`). The history append happens **after** the stream returns (`:139`), and the autosave goroutine is spawned **after** that (`:142`). `emitEvent` calls `devEventSink` and `runtime.EventsEmit` and nothing else (`app.go:355-377`); `devEventSink` is a package-level var set only by the harness (`harness_main.go:97`).

**O7-160 — The exclusion is architectural and must stay a contract rule.** `devEventSink` is a display sink with no persistence edge; a callback that reaches a write would put a persistence operation inside another operation's execution, which is the only construction that could create genuine same-thread re-entrancy. Prohibited by shape: no persistence call from an event sink, a token callback, or an emitted-event handler.

---

## 5. Cross-conversation operations: distinct locks, ordered

**O7-161 — Three operations touch more than one conversation in a single call:**

| Operation | Conversations touched | Shape |
|---|---|---|
| `App.LoadChat` (`app-chat.go:92-108`) | save **X**, then load **Y** | two sequential acquisitions |
| `pm.SaveChatAs` (`persistence.go:437-466`) | read **X**, write **Y** | two simultaneous |
| `WipeAllData` (`app-wipe.go:31`) | **all** | root-scoped, deletion-driven, outside enforcement |

The first is two acquisitions in sequence and needs no coordination beyond the identity binding: a crash between them is OQ-7c-ii's `O7-142` window and remains a preference loss. The second is the only true multi-lock operation in the content path, and it needs the **deterministic total order by stable conversation identity** already required (`OQ-7b`). The third is not a content mutation and must never acquire per-conversation locks in a loop — a root-scoped wipe that takes N locks in directory order is a deadlock generator, and it is also deletion, which is unauthorized.

**O7-162 — So the maximum simultaneous acquisitions in any authorized content operation is two, and the order is fixed.** That is the whole concurrency-complexity budget; no design needs to hold more.

---

## 6. Shutdown

**O7-163 — `performFinalSave` must be one critical section, and the in-flight autosave must be joined first.** Two requirements, both ordering rather than recursion:

1. `SaveChat` and `RewriteHistory` for the same conversation execute under **one** acquisition, not two — that is what makes the final save a coherent transition instead of the row 570 window.
2. The autosave goroutine (O7-155) must be awaited before the final save begins, or the final save races a writer that holds the same conversation legitimately.

Neither is re-entrancy, and neither is fixed by a lock that permits recursion.

---

## 7. Verdict

**O7-164 — No product requirement needs true re-entrancy.** Every situation that resembles it is either (a) composition inside one operation, which the `Mutate` shape makes a single acquisition (O7-154), or (b) overlap between operations, which needs units and identity, not recursion (O7-158). The contract may therefore **prohibit** re-entrancy rather than accommodate it.

**O7-165 — Bounded refusal, with four distinct outcomes the caller must be able to tell apart:**

| Outcome | Meaning | Caller obligation |
|---|---|---|
| acquired | the operation proceeds | — |
| `busy` | another writer holds this conversation (in-process or not) | report; do not silently retry |
| `refused-reentrant` | the same operation tried to acquire a conversation it already holds | **programming error** — must never be reported as contention |
| `closed` | shutdown has begun | discard the operation, do not defer it |

**O7-166 — `refused-reentrant` must be diagnosed above the OS lock.** On Windows `LockFileEx` locks are per-handle, so a second acquisition in the same process fails the same way a foreign process's would; without a separate in-process held-set, an internal composition bug would surface as "another instance is writing" and send every investigation in the wrong direction. The classification must be made before the OS call, not inferred from its failure.

**O7-167 — Refusal is bounded and non-blocking on the refusal paths.** No waiting inside the critical section, no retry loop inside the lock, no queue that grows during shutdown. Waiting is permitted only on the initial acquisition, bounded by the existing 5 s default with fail-fast posture.

**O7-168 — Refusal must never be converted into a content verdict.** A refused write is *not written*: reconciliation sees the previous state, which is exactly correct, and must not be told anything about it. An autosave that loses a race leaves a history that is short by one turn and **internally consistent** — Healthy, not Damaged, not Pending.

**O7-169 — Retry identity stays open.** Whether a refused deferred write is retried, and under what identity, is OQ-7d-ii. Nothing here decides it; O7-157 constrains any retry to the originally bound conversation.

**O7-170 — The `SaveChatAs` double write is recorded, not repaired.** Two writes to the destination with an empty intermediate state is the one place in the current code where a crash can leave a *shape* that looks like a real conversation with lost content. It is not fixed by re-entrancy, by a lock, or by ordering inside the lock: it is fixed by making the copy one write. Ownership: writer authority, alongside row 570, and unauthorized.

**O7-171 — Nothing here is authorized.** No `activeChatID`, no pointer change, no save-signature change, no shutdown repair, no goroutine join, no lock, no gate, no CAS, no `SaveChatAs` change, no migration, no deletion enablement.

---

## 8. Tests this would require

- **O7-T46** `SaveChatAs` completes under one acquisition, and no code path can attempt a second acquisition of the same conversation from inside it.
- **O7-T47** A deliberately nested acquisition returns `refused-reentrant`, and never `busy`.
- **O7-T48** A nested acquisition is refused before the OS lock call is attempted — asserted by observing that no lock file is opened twice.
- **O7-T49** Two overlapping autosaves of the same conversation: one proceeds, one is refused as `busy`, and the resulting history is internally consistent.
- **O7-T50** An autosave captured in X and executed after selection moved to Y writes to **X**, or refuses; it never writes to Y.
- **O7-T51** A final save racing an in-flight autosave is serialised, and both writes target the same conversation.
- **O7-T52** A refused autosave leaves a conversation that scans `Healthy`, short by one turn, and never `Damaged`.
- **O7-T53** `SaveChatAs` reads X and writes Y under the declared total order; a reversed order attempt is refused, not deadlocked.
- **O7-T54** No event sink, token callback or `devEventSink` implementation can reach a persistence write — asserted structurally (no reachability from the sink types to the writer).
