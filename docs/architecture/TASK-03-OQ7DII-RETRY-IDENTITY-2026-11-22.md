# TASK-03 — OQ-7d-ii: retry identity

| Field | Value |
|---|---|
| Purpose | Whether a deferred or partially applied conversation mutation can be retried against its originally bound conversation — and whether the existing data contains enough information to recognise such a retry safely |
| Status | **Read-only investigation. Design only. V3.1 unchanged. Nothing authorized.** |
| Method | Source inspection of every conversation writer and reader, the append and rewrite paths, the message identity fields, the pointer resolution sites, the OQ-6 marker shape, and the two in-repository precedents (`contentimport`, `durable`) |
| Recovery token | `TASK-03-OQ7DII-RETRY-IDENTITY-2026-11-22` |

Continues OQ-7's numbering (`O7-196`…). Resolves OQ-7d-ii as opened by `O7-79` (row 566) and scoped by `O7-193` (row 573). Consumes `O7-157` (row 572) as the binding requirement it extends, and keeps OQ-6 and OQ-7 separate (§6).

---

## 1. The question is five questions

**O7-196 — "Retry identity" is asked as one question and is actually five, and they are answerable at different prices:**

| # | Question | Answered by |
|---|---|---|
| Q1 | What operation identity is retried, and what proves idempotency for it? | operation identity (§2, §8) |
| Q2 | What binds the content to a conversation? | binding (§5) |
| Q3 | What is the outcome if the operation was refused before any mutation? | decidable today (§7) |
| Q4 | What is the outcome if content mutated but the marker was not settled? | undecidable today (§7) |
| Q5 | What is the outcome after the caller or the process restarts? | undecidable today (§3, §7) |

**O7-197 — Four of the five are not answerable from the existing data, and one of them (Q3) is answerable for a reason that does not generalise.** That asymmetry, not the absence of a retry API, is the finding. It is stated exactly in §7.

---

## 2. There is no operation identity in the conversation writer path

**O7-198 — Observed: no operation identity exists anywhere in the conversation store.** `ChatMeta` and `ChatHistory` carry no operation field (`internal/backend/chat/types.go`); `AppConfig` holds only `LastChatID`; `AppendHistoryEntry(chatID string, entry ChatHistory) error` has no key parameter; and `internal/backend/chat` contains no journal, operation table, receipt, intent record or sequence counter. Nothing in the package can be asked "has operation X been applied".

**O7-199 — Source-derived: message identity is positional and derived, not assigned.** Three sites produce it, and each computes it from length:

| Site | Expression | Consequence |
|---|---|---|
| `app-prompt.go:303` | `fmt.Sprintf("%d", len(a.chatMessages)-1)` | id = index in the in-memory list at append time |
| `app-chat.go:405` | `fmt.Sprintf("%d", len(history))` | id = index while rebuilding from app state |
| `app-prompt.go:357` | `fmt.Sprintf("%d", i)` | id = index while rebuilding from frontend entries |

**O7-200 — Calculated: a message id therefore cannot serve as an idempotency key in either direction.** It is not unique across time (a deleted message frees its number for the next one), not stable across a rewrite (every id after an insertion or deletion shifts), and not assigned at the mutation boundary (it is derived from a list that may already have changed). Matching a retry by id would match the wrong message; not matching it by id means there is nothing to match. This is the same property `FC1-06` prohibits and row 525 already recorded as the reason legacy history cannot be promoted.

**O7-201 — Source-derived: the append has no precondition of any kind.** `AppendHistoryEntry` (`persistence.go:668-706`) computes the directory, optionally migrates `history.json` on first use, marshals the entry, opens `O_APPEND|O_CREATE|O_WRONLY` and writes one line. It does not read the marker, compare a generation, inspect existing content, or count anything. Every call is by construction "apply again".

**O7-202 — Observed by absence: no writer records that an operation occurred.** `SaveChat`, `AppendHistoryEntry`, `RewriteHistory`, `DeleteChat` and `saveConfig` each either succeed or return an error; none leaves a durable statement that this attempt existed. A retry therefore has no counterpart to compare itself against, which is the precise sense in which "retry identity" is missing rather than merely unimplemented.

**O7-203 — Source-derived: failures carry no classification.** Every conversation writer returns a bare `error`. The package has no equivalent of `durable.Outcome` (`internal/durable/store.go:33-76`: `Saved`, `NotSaved(reason)`, `Busy`, `Recovered`, `RefusedUnknownSchema`), so a caller cannot distinguish *refused*, *not applied* and *unknown* — which is also why OQ-7c-iii's refusal taxonomy (`acquired`/`busy`/`refused-reentrant`/`closed`) has nothing to attach to today. The refusal vocabulary exists in the contract and nowhere in the code.

**O7-204 — Source-derived, with a consequence the contract must not overstate: the append does not fsync and discards the written count.** `persistence.go:702` is `if _, err := f.Write(line); err != nil` — `n` is discarded and no `Sync` is called; `saveJSON` (`:571-582`) likewise writes a temp file and renames without flushing. So a reported success is an *observed successful write call*, which is exactly the `durability: observed` meaning V3.1 already fixed, and explicitly not a power-loss claim. Any retry design that assumes a returned success survives a crash would be asserting durability the writer never claimed.

---

## 3. The one existing idempotency attempt is a content heuristic, and it must not be generalised

**O7-205 — Observed: the only idempotency logic in the conversation path is `recordUserTurn` (`app-chat.go:360-370`).** It suppresses a user turn when the immediately preceding in-memory message has the same role and byte-identical content.

**O7-206 — What it actually is, stated precisely, because each property disqualifies it as a retry mechanism:**

| Property | Value | Why it cannot serve |
|---|---|---|
| medium | in-memory `a.chatMessages` only | gone after a restart (Q5 unanswerable) |
| scope | the last element only | a duplicate at any other position is invisible |
| predicate | role + exact `content` equality | identity is inferred from content, which the contract forbids |
| coverage | ignores `sender`, `characterId`, `hasImage`, `timestamp` | two distinct turns can collide, one turn can evade |

**O7-207 — Calculated: the guard is sound only under an unstated caller protocol.** It cannot distinguish "this turn was already recorded by the frontend" from "the user sent the same text twice and the frontend has not recorded the second yet". Its correctness depends on `AppendChatMessage` having been called first for every turn — a protocol assumption held in two places and stated in a comment, not an invariant. That is the same defect class the governing principle excludes: a state that depends on a condition held in a comment or a caller convention.

**O7-208 — Inferred: it is evidence that the product wanted idempotency here, and evidence that the substitute chosen was content comparison.** The correct conclusion is narrow: this is a handling of one specific double-append, useful where it is, and unusable as precedent. Any retry rule that infers "already applied" from matching content would reproduce this heuristic at a larger scope with less information.

---

## 4. What the durable content can and cannot prove

**O7-209 — Source-derived: the durable content unit has no integrity information.** A `history.ndjson` line is the marshalled `ChatHistory`: id, sender, text, role, timestamp, hasImage, characterId. No checksum, no sequence, no operation id, no turn id, no record terminator beyond `\n`.

**O7-210 — Calculated: a torn append is a permanent truncation point, not a recoverable ambiguity.** Two source facts compose into this: the append writes at EOF (`O_APPEND`), and the reader stops at the first line it cannot unmarshal and returns the prefix (`persistence.go:609`, `break // return what we have so far`). A partially written line therefore has no newline; the retry's line is appended directly after it; the two form one invalid line; the reader stops there. Every subsequently successful append also sits beyond the tear and is unreadable. Reachability of the partial write itself is **Unestablished** — not reproduced in this session — but the consequence is not a guess: it follows from the two sites.

**O7-211 — Calculated: the snapshot fallback does not rescue that case, so the loss is silent.** `loadHistoryFile` prefers `history.ndjson` whenever it exists and is non-empty (`persistence.go:589-591`); a torn log is non-empty, so `history.json` is never consulted and the conversation presents the pre-tear prefix as its history. The function returns no error to its caller — by design, since a load must never fail — and prints only a line to stdout, so neither the loaded state nor the caller's result carries the truncation.

**O7-212 — Source-derived: `RewriteHistory` is the only writer with retry-safe semantics, and only by substitution.** `persistence.go:630-661` builds the complete file in memory and replaces it by temp-and-rename, then rewrites the snapshot. Re-running it with the same list produces the same bytes: the operation is idempotent because it *replaces* the state, not because it *recognises* a previous attempt. It cannot express "append message M", and it destroys evidence of any write that intervened between the caller capturing its list and the rewrite executing — which is why the same mechanism is behind the renumbering in O7-200.

**O7-213 — Inferred: two incompatible retry semantics coexist in one store with nothing recording which one a caller used.** An append is neither idempotent nor detectable; a rewrite is idempotent but requires the caller to still hold the whole intended state. A retry designed for one is wrong for the other, and no durable field distinguishes them.

---

## 5. The write is not bound to a conversation at capture

**O7-214 — Source-derived: the target is resolved at execution time, from the shared pointer.** `performSave` reads `GetLastChatID()` when it runs (`app-chat.go:314-316`); so does `appendChatMessage` (`app-prompt.go:296-298`). No conversation identifier is passed by the caller who created the work.

**O7-215 — Observed: no write carries the conversation it was created for.** The autosave goroutines capture `a.chatMessages` implicitly and the pointer explicitly-at-execution, so a save created while conversation X was selected can execute after selection moved to Y and write X's content into Y (row 572). A retry is the same shape: deferred work re-reading a mutable selector.

**O7-216 — O7-157 already names the requirement and supplies no mechanism.** Row 572 requires a deferred mutation to use the identity bound to its captured content and to refuse rather than retarget. A retry *is* a deferred mutation, so it is already covered in principle — but the identity it must use does not exist (O7-198), which is why this record and O7-157 close on the same missing object.

**O7-217 — And the pointer cannot be promoted into that identity.** Row 571 establishes `LastChatID` as a selection hint with no content authority in either direction; it cannot prove which conversation a write belonged to, only which one was selected when something read it. Binding must be captured at the boundary, never re-derived from it.

---

## 6. Settlement evidence is not a retry ledger

**O7-218 — Source-derived: the OQ-6 marker records a pair and a settlement, never an operation.** Its fields are `(state, generation)` with `settled_generation` and `settled_at` (row 562, `S-06`/`S-07`, corrected by `O7-20`). It holds no operation identity, no operation count and no per-operation outcome.

**O7-219 — Calculated: therefore `settled_generation = G` answers "generation G closed", and can never answer "operation X applied".** Two appends inside one generation are indistinguishable in the marker: closing the generation is one fact, and how many mutations composed it is not recorded.

**O7-220 — Inferred: the marker must not be extended into a retry ledger.** Adding per-operation entries would give one document two lifetimes (settlement of a generation, completion of an operation) and two authorities, and would recreate the defect row 573 was written to avoid — one artifact answering two questions that can be answered independently. If retry identity is adopted, it is its own record with its own crash and lifecycle definition, and the two must never certify each other, exactly as the two-record rule already requires of settlement and enforcement.

---

## 7. The five outcomes, enumerated against the actual writers

**O7-221 — The table is the answer to Q1–Q5.** Each row is a position the current code can actually reach, with what survives it and whether "already applied" is decidable:

| Position | Durable evidence | Decidable? | Permitted action |
|---|---|---|---|
| refused before mutation (lock busy, `refused-reentrant`, `closed`, identity invalid) | none — nothing written | **yes**: not applied | safe to retry, provided the identity is still bound (O7-216) |
| applied in memory, nothing written (killed before the append) | none | **no**: indistinguishable from never attempted | retry is the only forward action and is safe *because* nothing was written — not because the attempt was recognised |
| append partially written (torn line) | invalid tail | **no** — and not fixable by retry: the retry extends the tear (§4) | report unknown; do not retry; the log needs its own repair question, which this record does not open |
| append written, marker not settled | a complete line | **no**: the content is present, but that it is *this* operation's content is unprovable | report unknown; the content may not be presented as this operation's completion |
| restart after a complete append | a complete line | **no**: identical to the row above, and the in-memory heuristic is gone | report unknown |

**O7-222 — Only the refusal row is decidable, and only because "nothing was written" is a negative that no missing evidence can hide.** That is why OQ-7c-iii's refusal taxonomy is load-bearing for this question too: it is the only class of outcome the store can answer honestly today.

**O7-223 — The central asymmetry: the existing data is sufficient to prove non-application and insufficient to prove application.** Every row with a positive outcome leaves a record that is compatible with both "this operation applied" and "some other operation applied, or this one applied twice". No amount of reading closes that gap, because the distinction was never written down.

**O7-224 — Calculated: therefore the contract has exactly two honest options.** Either retry support is restricted to the refusal row, or it requires a durable operation identity. There is no third option that stays truthful — content matching, id counting, timestamp proximity and pointer inspection are all guesses, and each has a demonstrated failure mode above (O7-200, O7-207, O7-210).

---

## 8. What a durable operation identity would have to satisfy

These are necessary properties derived from the writer path and the gate work, not a design and not an authorization.

**O7-225 — Bound at capture.** It carries the conversation identity captured when the work was created, never one re-read at execution (O7-216, O7-157).

**O7-226 — Written before the content, where the content can be recognised against it.** If the record is written after the content, the window between them is exactly the ambiguous window it exists to remove. The ordering is therefore a contract clause, not an implementation detail — the same reasoning that made OQ-6's settlement precondition the *hold* rather than the observation (`O7-20`).

**O7-227 — Compared inside the same critical section as the content write.** A check followed by a separate write races by construction; the gate audit already found that caller-side checks cannot be enforcement (row 556). Retry identity has the identical shape: the "has this applied" test and the apply must be one atomic boundary.

**O7-228 — Terminal settlement distinct from the content.** "Applied but not settled" must be a stored state, not an inference, or Q4 remains unanswerable in the same way it is today.

**O7-229 — Never derived from content.** Two identical turns are two turns; one turn recorded twice is one turn; content equality cannot tell them apart (O7-200, O7-207).

**O7-230 — Per-conversation scope.** The store, the lock and the transition are all per-conversation; an identity scoped wider would need the deterministic total lock order for an operation that touches one directory.

**O7-231 — Terminal and re-execution-refusing.** A second execute of a terminal identity is refused with a typed result rather than re-applied. This is the shape `contentimport` already implements (§9), which establishes that the pattern is buildable in this repository — not that it exists in this path.

---

## 9. The in-repository precedent, and what it does not license

**O7-232 — Source-derived: `contentimport` already implements this exact contract shape for a different domain.** `Execute` accepts a caller-supplied or generated operation id, refuses an existing operation as `already_imported` (terminal) or `operation_active`, compares a plan hash and refuses a mismatch (`plan_mismatch`), journals per-item states including `StateAmbiguous` for "cannot be proven", and exposes `Resume` as the only recovery entry point, explicitly never called at startup (`internal/contentimport/service.go`, operation-identity block at `:299-313`; `worldstore.StageImport` likewise refuses a missing or invalid operation identity, `import.go:166-172`).

**O7-233 — Inferred: what that precedent does not license.** Its idempotency rests on whole-item ownership evidence — staging owner files, document and lore hashes, a journaled revision per item — and its unit is a whole canonical document, not a line in an append log. Nothing in it supplies message-level identity, torn-line handling, or a per-generation interaction with the OQ-6 marker. It is evidence that the *shape* is idiomatic here, and it is not a component that can be reused by reference.

**O7-234 — Source-derived: `durable` supplies the outcome vocabulary the conversation path is missing.** Its `Result` distinguishes `Saved` from `NotSaved(reason)` from `Busy` from `Recovered` from `RefusedUnknownSchema`, and it proves the negative case concretely: `NotSaved(rename_denied_by_open_reader)` with the previous document intact (`store.go:33-76`, tests at `store_test.go:189-195`). The conversation path reports none of this — which means the retry question cannot even be *stated* to a caller until the writer can say which of those five things happened.

---

## 10. Disposition of unsaved work on refusal

**O7-235 — A refusal is a statement about one state domain, and the contract must name which.** Three domains hold the turn and they fail independently: `stored` (the conversation directory, unchanged by a refusal), `captured` (`a.chatMessages`, in memory, never verified against disk), and `displayed` (the frontend). "Stored state is intact" is true and says nothing about the other two; a refusal that reports only the first invites the reader to conclude the turn is saved.

**O7-236 — Observed: the failure path already loses that distinction, in shipped code.** `appendChatMessage` appends to `a.chatMessages` *before* attempting the disk write (`app-prompt.go:281-305`), and the streaming path discards the result outright: `_ = a.AppendChatMessage(...)` at `app-prompt.go:138`. So a failed append leaves the message in memory and on screen, absent from the authoritative log, and reports nothing to the user.

**O7-237 — Calculated: the divergence persists for the rest of the session and is repaired only at exit.** The next `performSave` writes `history.json` from app state, while the log — which the loader prefers whenever non-empty — still lacks the entry, so a reload would show the conversation without it. The only path that restores it is `performFinalSave`'s rewrite, which row 570 already records as two representations written in one operation with its own window. The turn is therefore *not* lost by refusal; it is dependent on a clean shutdown that the shutdown path does not currently guarantee.

**O7-238 — Required by the owner, carried as a clause: a refusal on invalid identity must not silently discard captured content.** It must return an explicit failure result and a defined disposition — and the disposition may not be "clean up the in-memory copy", because discarding captured work is a data loss performed by the component that declined to write.

**O7-239 — Inferred: the honest result set is three names, not a boolean.** `not_attempted` (refused before any mutation), `not_applied` (attempted, provably no effect), and `unknown` (attempted, effect unprovable) — each accompanied by the three domain statements of O7-235. `unknown` may never be rendered, logged or summarised as saved; that is the same rule the scanner already carries in "a scan error must never become `Missing`", applied to a writer.

---

## 11. Verdict: what the contract may and may not require

**O7-240 — Stated plainly: retry identity is not derivable from the existing data.** No operation id, no message identity that survives a rewrite, no failure classification, no settlement per operation, and no binding carried with the work. The store can prove that nothing happened; it cannot prove that something did.

**O7-241 — Therefore only two retry semantics may be specified today:** retry of a *refused* operation, which is safe because nothing was written (O7-222), and retry of a whole-state *rewrite*, which is content substitution and must be described as such — never as "append exactly once" (O7-212).

**O7-242 — And retry of an ambiguous append is prohibited, not merely discouraged.** It must be reported as `unknown`. Retrying duplicates the record in the ordinary case and extends a tear in the torn case (O7-210), so the "safe" retry is the one that makes the state worse and unreadable.

**O7-243 — Reconciliation is unaffected in both directions.** The scanner never retries, never repairs and never chooses a winner; a retry executed by a writer cannot promote a state, and no retry rule may be used to move `Unavailable` toward `Healthy`. Likewise a settlement record cannot certify enforcement and an enforcement record cannot certify settlement — the two-record rule stands, and retry identity, if adopted, becomes a third record that certifies neither.

---

## 12. Tests this would require

- **O7-T64** A message id is shown to be positional: deleting an earlier entry renumbers the survivors, and two different messages can hold the same id at different times.
- **O7-T65** `AppendHistoryEntry` is called twice with one entry and produces two lines; no path detects the duplicate.
- **O7-T66** A partially written final line makes the reader return the prefix, and a subsequent successful append is not visible — the tear is a truncation point, asserted structurally.
- **O7-T67** A torn log is non-empty, so `history.json` is not consulted; the fallback does not engage.
- **O7-T68** `RewriteHistory` run twice with the same list produces byte-identical files — content idempotency by substitution, asserted separately from append.
- **O7-T69** `recordUserTurn` suppresses an identical adjacent turn and cannot suppress one at any other position, documenting the heuristic's boundary rather than its correctness.
- **O7-T70** `_ = a.AppendChatMessage(...)` on an injected append failure leaves the message in `a.chatMessages`, absent from `history.ndjson`, and present in `history.json` after the next save — the O7-236 divergence, instrumented.
- **O7-T71** An operation refused before mutation leaves the conversation directory byte-identical, and reports `not_attempted` with all three domain statements.
- **O7-T72** No code path answers "was this applied?" from content, id, count, timestamp or pointer — asserted as the absence of any such predicate, mirroring the closure test the gate work requires.

---

## 13. Nothing here is authorized

**No implementation follows.** This record authorizes no code, no migration, no repair of the torn-line behaviour, no operation journal, no change to `AppendHistoryEntry`, `RewriteHistory`, `SaveChat`, `recordUserTurn` or their callers, no message-identity change (positional ids stay positional), no binding field on `App`, no pointer change, no lock, no gate, no shutdown repair, no marker field, no new record, no backfill, and no deletion enablement. It establishes what any retry claim would have to rest on, and that two of the three retry shapes the product might want are unsafe in the current store. Reconciliation keeps reporting and never retries.
