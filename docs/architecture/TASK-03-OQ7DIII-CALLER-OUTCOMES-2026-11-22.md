# TASK-03 — OQ-7d-iii: caller-visible outcomes and the disposition of unsaved work

| Field | Value |
|---|---|
| Purpose | Define what a caller must receive for every mutation outcome, what happens to displayed-but-unstored work on chat switch and shutdown, which existing callers discard a result, and how a retry stays bound to its original conversation |
| Status | **Read-only investigation and proposal. V3.1 unedited. Nothing adopted, nothing authorized.** |
| Method | Source inspection of every mutation entry point and its callers (Go bindings, autosave goroutines, shutdown stages, Wails bindings and frontend call sites), read against OQ-7d-ii's outcome table |
| Recovery token | `TASK-03-OQ7DIII-CALLER-OUTCOMES-2026-11-22` |

Continues OQ-7's numbering (`O7-244`…). Closes OQ-7d-iii, opened by the owner's review of row 574. Consumes `O7-221`…`O7-243` (row 574) and `O7-157` (row 572). It does **not** select either retry policy; §9 records a recommendation and no adoption.

---

## 1. Scope, and the severity this record is obliged to state

**O7-244 — This record answers four questions and no others:** what the caller receives for each outcome; what happens to captured work on a refusal; what happens to displayed-but-unstored work on chat switch and on shutdown; and what binds a retry. It does not design a UI, a journal, a lock, a gate or a repair.

**O7-245 — The two retry policies remain unselected.** Row 574 left (a) refusal-only retry and (b) durable operation identity as the only honest options (`O7-224`). This record separates what each needs (§8) but adopts neither; the choice is the owner's and belongs in the contract.

**O7-246 — Severity, stated plainly: the discarded append error is release-blocking silence, not a missing diagnostic.** Today three different user-visible situations are indistinguishable to the person using the product:

| What the user sees | What is actually true | What the product says |
|---|---|---|
| the message on screen | stored in the authoritative log | nothing; the call returned nil |
| the message on screen | present only in memory, absent from the log | nothing; the call returned an error that was discarded |
| the message on screen | write attempted, outcome unknown | nothing; no vocabulary exists to say it |

A missing log line is recoverable and diagnosable. A user who is told nothing cannot know that anything needs recovery, so the second and third rows persist until a reload silently removes the turn. That is why this belongs in the contract's failure-visibility clause (`V3-27`) and not in a later polish pass.

---

## 2. The result a caller must receive

**O7-247 — A mutation result is a statement over independent domains, not a boolean and not a bare `error`.** The minimum shape is seven fields:

| Field | Values | Why it cannot be omitted |
|---|---|---|
| `Outcome` | `applied` · `no_op` · `not_attempted` · `not_applied` · `unknown` | the only field a caller can branch on |
| `ConversationID` | the identity captured at the boundary | a retry must not re-resolve it (§8) |
| `Stored` | `observed` · `not_stored` · `unknown` | the authoritative store, stated separately |
| `Captured` | `retained` · `not_retained` | the in-memory copy, which a refusal must not destroy |
| `Displayed` | `shown` · `not_shown` · `unknown` | the frontend, which is neither of the above |
| `Retry` | `permitted` · `prohibited` · `not_applicable` | derived from `Outcome`, never guessed by the caller |
| `Reason` | a closed vocabulary, including `identity_unbound` | distinguishes refusals that look identical |

**O7-248 — The five outcome values extend row 574's three, and the extension is deliberate.** `O7-239` named `not_attempted`, `not_applied` and `unknown` because the question was *failure*. A result type that can only report failure cannot report success either, so `applied` (the content write returned success) and `no_op` (the request was validly empty) are the missing counterparts. `applied` means **observed write success**, never durability: `AppendHistoryEntry` discards the written count and never syncs (`persistence.go:702`, `O7-204`), so `Stored = observed` is the strongest honest claim.

**O7-249 — The three domain statements may not be collapsed, and only one combination may be presented as saved:**

| `Stored` | `Captured` | `Displayed` | May be presented as saved? | The only permitted wording |
|---|---|---|---|---|
| `observed` | `retained` | `shown` | **yes** | "saved" (write observed; durability not claimed) |
| `not_stored` | `retained` | `shown` | **no** | "not saved — kept in this window" |
| `unknown` | `retained` | `shown` | **no** | "save outcome unknown" |
| `not_stored` | `not_retained` | `shown` | **no** | "not saved" |
| `unknown` | `not_retained` | `shown` | **no** | "outcome unknown; content no longer held" |

**O7-250 — `unknown` is a terminal label, not a spinner.** It may never be rendered, logged, summarised, retried or promoted as saved. The rule is the writer-side form of the scanner's existing "a scan error must never become `Missing`" (`V3-60`, `V3-47`): a missing observation is not a negative observation.

**O7-251 — Every mutation entry point must return this result; a bare `error` is insufficient.** All conversation writers currently return `error` (`O7-203`), so the vehicle for a typed result does not exist yet. `error` can carry "something failed" and cannot carry *which domain failed*, *whether content was retained*, or *whether retry is permitted*.

**O7-252 — Calculated: today's `nil` is not even a claim about the write.** `appendChatMessage` returns `nil` in two places before any disk write is attempted (§3), and `AppendHistoryEntry` returns `nil` after a `Write` whose count is discarded. So "no error" is compatible with *applied*, *no-op*, *never bound* and *partially written*. That is the precise reason the result type is a contract object rather than a refactor convenience.

---

## 3. The positions the existing code reaches but row 574's table does not name

**O7-253 — Observed: an append with no active conversation mutates memory and reports success.** `appendChatMessage` appends to `a.chatMessages` at `app-prompt.go:294`, reads `GetLastChatID()` at `:295-298`, and returns `nil` at `:299-301` when it is empty. The message is captured and displayed; nothing is written; no conversation is identified; the caller is told nothing went wrong.

**O7-254 — Observed: an empty-content append is a different `nil`.** `app-prompt.go:282-284` returns `nil` before appending when `TrimSpace(content) == ""`. That is a legitimate no-op, but with a bare `error` it is byte-identical to the success of §O7-253 and to a completed append.

**O7-255 — Calculated: therefore "success" currently conflates four distinct positions.** They are *applied*, *valid no-op*, *captured but unbound*, and *refused* — and a fifth, *torn*, is reachable through the same return path. Row 574's five-row outcome table (`O7-221`) is correct about durable evidence and incomplete about what the caller is told, because two of these positions never reach the store at all. The contract must add the unbound position rather than infer it from `chatID == ""` at each call site.

---

## 4. Disposition of captured work on a refusal

**O7-256 — A refusal states what happens to the captured copy, in the result, not in a comment.** When `Outcome` is `not_attempted`, `Captured` must be `retained`, and the copy must remain associated with the `ConversationID` captured at the boundary. Discarding it is a data loss performed by the component that declined to write (`O7-238`).

**O7-257 — "Retained" is not "retryable".** The retained copy permits an explicit user retry only where the outcome proves nothing was applied. The mapping is fixed:

| `Outcome` | `Retry` | Rationale |
|---|---|---|
| `not_attempted` | `permitted` if the identity is still bound | nothing was written (`O7-222`) |
| `not_applied` | `permitted` if the identity is still bound | attempted, provably no effect |
| `unknown` | `prohibited` | row 574: retry duplicates or extends a tear (`O7-242`) |
| `applied` / `no_op` | `not_applicable` | there is nothing to retry |

**O7-258 — For `unknown`, the only honest affordances are explicit ones that do not write again:** keep the captured copy visible, mark it unknown, and offer the user a choice between an explicit discard (reported as data loss, never as cleanup) and an export of the captured content. A "Retry" button on an unknown outcome would implement the one action row 574 prohibits.

**O7-259 — A refusal may not be reported as a state of the conversation.** It is a result of an operation. No refusal may move a reconciliation verdict toward `Healthy`, and no scanner may read a refusal record as evidence of content (`O7-243`).

---

## 5. Displayed-but-unstored work on chat switch

**O7-260 — Observed: switching conversations attempts a save, prints a warning on failure, and continues.** `App.LoadChat` calls `performSave()` when the pointer names a different conversation (`app-chat.go:96-99`); on failure it writes `Warning: Failed to save current chat: …` to stdout and proceeds to load the requested conversation. The save failure is not returned, so the frontend cannot distinguish a clean switch from a lossy one.

**O7-261 — Observed: the switch then replaces the in-memory state.** `a.updateAppStateFromChat(chatData)` at `app-chat.go:106` overwrites `a.chatMessages` with the newly loaded conversation. The captured turn of the outgoing conversation therefore has no in-process holder after the switch.

**O7-262 — Calculated: after a switch the unsaved turn is not merely unsaved — it is unrecoverable in-process.** This qualifies row 574's wording. A later rewrite *may persist* a captured turn only while that conversation remains loaded and bound; `performFinalSave` rewrites the conversation the pointer now names (`app-chat.go:327-350`), which after a switch is the new one. So the final rewrite cannot restore a turn abandoned at a switch, and the accurate statement is: *a later rewrite may persist captured work only before the conversation is switched away; after the switch the possibility is gone.*

**O7-263 — Therefore a switch is a disposition point and must be specified as one.** One of three things must happen before the in-memory state is replaced: the save is confirmed (`Stored = observed`); or the outgoing conversation is reported as having unsaved captured work, with the user able to cancel the switch; or the switch is refused (`not_attempted`, reason `unsaved_content`). What may not happen is the current behaviour — discard the copy and continue silently.

**O7-264 — A stdout warning is not a caller-visible result.** It is invisible in a GUI, absent from the RPC result, and not recorded where a user could find it. `LoadChat` returns only the load error (`app-chat.go:101-105`); the save error never reaches its caller.

---

## 6. Displayed-but-unstored work on shutdown

**O7-265 — Observed: shutdown makes one save attempt and reports failure only to the console.** `runShutdown` sets `shuttingDown`, cancels streams and jobs, then runs `performFinalSave` as a single stage (`app-lifecycle.go:96-122`); stage failures are accumulated and printed by `OnShutdown` (`app.go:345-349`) or `ExitApp` (`app-chat.go:155-161`). There is no user decision and no second attempt.

**O7-266 — Calculated: a failed shutdown save and a clean one look identical to the user.** The process exits either way; the only channel is stdout. For the active conversation, the honest report when the stage fails is `Stored = unknown` (or `not_stored` when the failure is known to precede the write), never silence and never `saved`.

**O7-267 — And a successful shutdown does not make the turn authoritative until both writes complete.** `performFinalSave` runs `SaveChat` (snapshot from app state) and then `RewriteHistory` (log and snapshot) — row 570's two-representation window. A death between them leaves the log behind the snapshot; the loader prefers the non-empty log, so the displayed turn is absent on the next launch. This is the precise reason the permitted wording is *may persist*, not *restores*: the rewrite is an intention to persist, executed in an order that has its own inconsistency window until row 570 is repaired and tested.

**O7-268 — Therefore shutdown must report disposition per conversation, or report `unknown`.** It may not report success from the fact that the process shut down cleanly. If the shutdown cannot state which conversations had captured work, it must say so for the active conversation rather than omit the statement.

**O7-269 — Inferred: shutdown cannot compensate for a switch.** Since the outgoing conversation's captured content is gone from memory at `app-chat.go:106` (`O7-261`), no shutdown-time mechanism can save it. The switch disposition point is the only place that work can be preserved, which is why §5 is not a UI preference.

---

## 7. Every existing caller, and the discard sites

**O7-270 — Observed: the product has no conforming caller today.** Every path that can receive a mutation failure either discards it or reduces it to a console message:

| Site | Operation | Disposition |
|---|---|---|
| `app-prompt.go:138` | `_ = a.AppendChatMessage("assistant", …)` | Go-side discard; the reply is displayed and may be absent from the log |
| `frontend/dist/js/chat-input.js:247` | `AppendChatMessage('user', …).catch(() => {})` | user turn; refusal invisible |
| `frontend/dist/js/chat-actions.js:172` | `AppendChatMessage('user', newText, …).catch(() => {})` | edited turn; refusal invisible |
| `frontend/dist/js/chat-actions.js:77` | fallback re-append loop `.catch(() => {})` | surviving entries may be partially written |
| `frontend/dist/js/chat-input.js:258` | post-turn `SaveChat().catch(err => console.warn(…))` | console only |
| `frontend/dist/js/chat-actions.js:69` | `await app.ReplaceChatHistory(entries)` inside `try` | outer `catch` → `console.warn` only |
| `app-chat.go:137,148` | `_ = a.performSave()` after world/scenario clear | Go-side discard |
| `app-characters.go:242,257` | `_ = a.performSave()` after character edit | Go-side discard |
| `app-prompt.go:142-148`, `app-chat.go:187-193` | autosave goroutines, `fmt.Printf` on failure | console only, unjoined |
| `app-lifecycle.go:121` | final save via `runStage` | accumulated, printed at exit |
| `app-chat.go:97-99` | switch save | warning, then continues |

**O7-271 — Calculated: the three frontend `AppendChatMessage` sites are exactly the user-turn path.** The one operation whose loss is most visible to the user is the one whose result is discarded at every entry point. The Go-side assistant path (`:138`) is the same shape on the other side of the turn.

**O7-272 — The contract rule must be a negative assertion, not a checklist.** "The callers we know about handle the result" establishes nothing about the next caller. The assertable form is the same shape as the write-containment assertion (`O7-38`): *the set of call sites that discard or downgrade a mutation result is a subset of an explicitly enumerated allowlist*, and the allowlist may contain only adapters that cannot lose content. Any new ignore site fails the check.

**O7-273 — The Wails boundary cannot be the enforcement point.** A JavaScript promise rejection is optional to handle, so returning a typed result to the frontend does not force the frontend to read it. Enforcement therefore needs two independent parts: the backend returns a result the caller must consume, and the displayed state is derived from that result rather than from the optimistic DOM update. A UI that adds the bubble before the result arrives has already displayed an unverified state; it may do so only with the `Displayed = shown, Stored = unknown` label of `O7-249`.

**O7-274 — The four `_ = a.performSave()` sites are not history mutations but they carry the same failure.** They are settings/selection writes that also persist chat state. Converting them is not optional under this rule: a caller that receives `nil` from `ClearWorldFromSession` today is told the clear succeeded even when the save failed. Either the operation returns a result, or the binding is explicitly classified as best-effort in the contract's allowlist.

---

## 8. Retry binding, separated from retry identity

**O7-275 — Every retry carries the conversation identity captured when the work was created.** It never re-reads `GetLastChatID()`. This is `O7-157` restated as a caller rule: a deferred mutation that finds its conversation no longer bound refuses rather than retargeting (`O7-216`, `O7-225`).

**O7-276 — Calculated: the current API cannot express a bound retry.** `performSave` and `appendChatMessage` both resolve the target at execution time (`app-chat.go:314-316`, `app-prompt.go:295-298`), so a retry affordance built on today's bindings would write into whichever conversation is selected when the user clicks. That is the row 572 retargeting defect with a user gesture attached.

**O7-277 — The separation that lets the safe policy proceed without a journal:**

| Need | Required for refusal-only retry? | Required for ambiguous retry? |
|---|---|---|
| **Conversation binding** — the captured `ConversationID` | **yes** | yes |
| **Operation identity** — a durable key recorded before content | no | **yes** |
| **Outcome vocabulary** — `not_attempted` / `not_applied` / `unknown` | **yes** | yes |
| **Terminal settlement record** | no | yes |

Conversation binding is a property of the caller and the boundary; operation identity is a new durable record. Conflating them is what made row 574's two options look like one decision. They are not: the refusal-only policy is implementable once binding and the result type exist, with no journal.

**O7-278 — If the bound conversation is no longer available, retry refuses.** Deleted, switched away, or unreadable: the result is `not_attempted`, reason `identity_unbound`, `Stored = not_stored`, `Captured = retained`, and the content is offered for export or explicit discard. It is never written into another conversation and never silently dropped. A retry that cannot find its destination has not been refused by contention; it has been refused by identity, and the two must not share a diagnosis.

---

## 9. Recommendation, and why it is not an adoption

**O7-279 — Recommended baseline, carried as a proposal:** refusal-only retry, the typed result of §2, the switch and shutdown disposition points of §5–§6, the negative discard assertion of §7, and conversation binding as a precondition. Durable operation identity is a separately specified capability, with row 574's seven properties, required before any ambiguous retry may be enabled.

**O7-280 — Rationale: this ordering does not make honest failure reporting depend on an unbuilt journal.** The product can stop lying about saves — and stop losing turns at switches — using only a result type, a captured identity and the existing writers. Ambiguous retry is the narrower, later capability, and it is the only part that needs a new durable record.

**O7-281 — Neither policy is adopted, and no clause here amends V3.1.** This record supplies the caller-facing clauses a later consolidation may adopt; it edits nothing and authorizes nothing.

---

## 10. Tests this would require

- **O7-T73** `AppendChatMessage` with no active conversation returns a result whose `Outcome` is not `applied`, `Stored = not_stored`, `Captured = retained`, and `ConversationID` empty — the unbound position of §3.
- **O7-T74** A whitespace-only append returns `no_op`, distinct from `applied` and from `not_attempted`.
- **O7-T75** An injected append failure returns `not_applied` or `unknown` with all three domain statements, never `applied`; the conversation directory is byte-identical for `not_applied`.
- **O7-T76** A refusal retains the captured message in memory, still associated with the originally captured conversation, and does not clear it.
- **O7-T77** `LoadChat` with a failing save does not silently continue: it either refuses the switch or returns a result naming the conversation with unsaved captured work.
- **O7-T78** After a switch, the outgoing conversation's captured turn is absent from memory, and a subsequent `performFinalSave` does not write it — the `O7-262` limit, asserted rather than assumed.
- **O7-T79** A failing final save reports `Stored = unknown` or `not_stored` for the active conversation, and the shutdown path's result is not reduced to a console line only.
- **O7-T80** A retry after selection moved refuses with `identity_unbound`, writes nothing to either conversation, and retains the captured copy (extends `O7-157`).
- **O7-T81** The negative assertion: no call site outside the enumerated allowlist discards a mutation result — assembled over Go and frontend call sites, failing on a new ignore pattern.
- **O7-T82** `unknown` is never rendered, logged or summarised as saved — a string/result assertion over every label the UI can produce.

---

## 11. Nothing here is authorized

**No implementation follows.** This record authorizes no result type, no change to `AppendChatMessage`, `performSave`, `LoadChat`, `performFinalSave`, the autosave goroutines or any frontend call site; no switch refusal, no shutdown reporting, no retry affordance, no binding field on `App`, no operation identity, no journal, no marker change, no lock, no gate, no repair of row 570's window, and no deletion enablement. It states what a caller must be told and what may not be discarded; the two retry policies remain unselected.
