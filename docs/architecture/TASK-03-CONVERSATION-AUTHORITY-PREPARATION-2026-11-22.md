# TASK-03 — Conversation Authority: preparation and investigation package

| Field | Value |
|---|---|
| Task | Phase 1B — conversation identity, ledger, projections |
| Status | **Prepared. Implementation NOT authorized.** Investigations are authorized; code is not |
| Owner decision | 2026-11-22: TASK-03 recommended and approved as a preparation boundary |
| Authority | ADR-001 (accepted), ADR-002 (accepted). **FC-1 is a target to satisfy, not an approved contract** |
| Depends on | `TASK-01` closed · `TASK-02` closed (primitive, unadopted) |
| Must not | become a sixth conversation-state authority · migrate settings/chat/gallery · amend any contract |
| Recovery token | `TASK-03-PREP-2026-11-22` |

**One sentence.** Before any ledger is built, decide exactly which representation becomes authoritative, how the
five existing ones are retired, and how the ledger itself survives a crash.

---

## 1. Scope

**In scope (investigation and decision only).** Representation and replacement · ordering and writer authority ·
selection derivation and recovery · identity under deletion/rename/movement/restoration · ledger durability shape ·
FC-1 clause triage · containment disposition.

**Out of scope.** Any implementation · migration of settings, chat persistence or gallery · character write paths
(C-1) · composition/tokenizer (1C) · memory-record authority (M-4) · branching (FC1-14b) · snapshot contents
(FC1-19a) · voice revisions (FC1-27a).

## 2. Ownership contract

| Question | Answer |
|---|---|
| **Owner** | Conversation authority — the ledger, and the derived selection/fold |
| **State** | Append-only events; derived current selection; derived state/message/turn identities |
| **Writers** | Exactly one commit path per conversation. Readers never write |
| **Readers** | Composition (1C), UI, restore, diagnostics — all read derived views, never the ledger's internals |
| **Lifecycle** | Created with the conversation; deleted with it under a recorded receipt; never swept |
| **Must not own** | Authored content (library, C-1), composition, token budgeting, memory records, snapshots' *contents*, file layout authority (PathService) |
| **Failure** | Append failure = no commit; recovery reports what it verified and what it inferred; a missing receipt is `Unknown`, never success |

## 3. The five existing implementations — verified inventory

Read from this tree. These are the implementations the branch must retire, not coexist with.

| # | Implementation | Evidence | Role today |
|---|---|---|---|
| 1 | In-memory `a.chatMessages []map[string]interface{}` | `app-chat.go:53,126,184,361–363` | App-side conversation buffer; written and read as truth during a session |
| 2 | `chat.Service`'s own context | `internal/backend/chat/service.go:52–75` (`context`, `worldName`, `scenarioName`, `characterNames`, `worldData`, `scenarioData`, `characterData`) | A second in-memory conversation/composition context in the backend |
| 3 | `history.json` snapshot | `persistence.go:211` (write), `:515` (read) | Full-history JSON; written via `saveJSON` (temp + rename, **no flush**) |
| 4 | `history.ndjson` append log | `persistence.go:678,696` (`O_APPEND`), `:586` (read), `:631` `RewriteHistory` (atomic full replacement) | Declared "true source of truth" at `persistence.go:209`, yet **two** mutation paths exist |
| 5 | Frontend `localStorage` | `app-wipe.js:16–24` — `howl-chat-current-chat`, `howl-chat-memories`, `howl-chat-characters`, … | Holds **current chat selection** and content copies, in the browser |

Two facts worth carrying into the decision:

- **`RewriteHistory` already performs in-place replacement of the append log** (`persistence.go:631,651`). An
  append-only ledger cannot coexist with it; deleting or rewriting that path is part of this branch, not a follow-up.
- **`persistence.go:210` records that frontend DOM scraping can corrupt `history.json` roles/ids.** The fifth
  implementation is not hypothetical: it has already damaged truth, and the code comments say so.

## 4. What must be decided

### D-1 — Representation and replacement (the decision record's core)

| Question | Why it blocks |
|---|---|
| Which representation becomes authoritative? | FC1-01 says the ledger. That settles the *shape*, not *which file* |
| Does the ledger **replace** `history.ndjson`, or **evolve** it (framed records, IDs, checksums added in place)? | Introducing a new file beside the existing log creates the sixth authority the owner prohibited. Evolving it is the only option that does not |
| How does each of the five transition, and in what order? | The slice must retire four of them; without an order, the app has two authorities mid-slice |
| What happens to `RewriteHistory` and the DOM-scrape write path? | Both mutate truth in place; both must be removed or converted to append-only events |

### D-2 — Ordering and writer authority (`FC1-12`, `FC1-CONCURRENCY`)

Concurrent append serialization · sequence assignment · gap policy · multi-process writer authority · conflict
handling. **Method:** state the single-writer rule per conversation, whether it is enforced by `LockFileEx`
(reusing `internal/durable`'s lock module) or by process ownership, and what a second process does — `Busy`, never
a merge. **Output:** a decision plus the invariant each rule protects.

### D-3 — Selection after deletion (`FC1-SELECTION-AFTER-DELETION`)

The invariant is fixed (no dangling `current_state_id`, no silent reselection, no manufactured identity). The
mechanism is not. **Method:** enumerate candidate resolutions (nearest ancestor, explicit empty, unresolved marker),
pick one, and state what a restart reproduces. Note that today selection lives in `localStorage`
(`howl-chat-current-chat`), which is itself a competing authority and must stop being one.

### D-4 — Ledger durability shape (see §6)

Torn-tail detection · flush policy per commit · truncation rule · duplicate suppression under crash · compaction.

## 5. Assessment A — FC-1 clause triage

FC-1 is `PROPOSED`, 0 approvals. Implementing a clause does not approve it.

**Compatible as written (satisfiable, no clarification needed).** FC1-01 (ledger as durable truth) · FC1-03, FC1-17
(no authored-content ownership) · FC1-05/06 (stable, non-positional identities) · FC1-09 (immutable events) ·
FC1-10 (idempotency) · FC1-13a's invariant · FC1-15, FC1-22, FC1-22a, FC1-22b (retention, manual deletion, no
cascade) · FC1-24 (deletion receipt) · FC1-26 (voice survives flush) — the last only as a constraint on the ledger,
since TTS is not this branch.

**Needs clarification before implementation.**

| Clause | Clarification required |
|---|---|
| FC1-02 | Is the state/branch **graph** built in this branch, or is a linear state chain sufficient for Phase 1 with the graph derived later? The plan defers branching, but the envelope carries `predecessor_event_id` |
| FC1-08 | Which envelope fields are mandatory in Phase 1 — in particular `durability` (5 values) and `provenance` |
| FC1-11 | The minimum event-kind set for Phase 1 only |
| FC1-22c / `FC1-STATE-RETENTION` | Is the ledger deleted with its conversation? Implementation cannot proceed without this answer |
| FC1-25 / `FC1-PRIVACY-DELETE` | Whether privacy deletion rewrites retained conversation events or appends a redaction event |

**Unresolved — blocks or bounds this branch.** `FC1-12`/`FC1-CONCURRENCY` (D-2) ·
`FC1-SELECTION-AFTER-DELETION` (D-3) · `FC1-14b` branching (deferred, but the envelope must not pretend to settle
it) · `FC1-14c` turn/message relationship (Phase 1 needs a minimal position: one turn = one message, stated as a
Phase 1 limitation) · `FC1-07`/`FC1-LEGACY` (migration only; legacy stays unresolved and read-only).

**Out of scope for TASK-03.** FC1-22e–m (memory copies, M-4) · FC1-19/19a (snapshot contents) · FC1-27a (voice
revision binding) · FC1-14d (selection identity, only if D-3 needs it).

## 6. Assessment B — can the ledger use DurableStore as-is?

**Not assumed.** The primitive's guarantees are document-shaped: whole-document replacement, atomic visibility,
two retained pre-images, `Busy` on contention, hash-verified recovery.

| Option | What it means | Cost |
|---|---|---|
| **As-is** | One store document per conversation; every append rewrites the whole ledger | Correct and simplest. Each append is O(n) writes plus 2 pre-image copies — 3× amplification per turn, growing with conversation length |
| **Narrow extension** | Segmented ledger part-files plus a manifest, each segment a DurableStore document | Preserves the primitive, but the manifest becomes a second pointer to truth and compaction becomes undefined work |
| **Separate journaling primitive** | Framed append-only log (length + checksum per record), single-writer lock reusing the same `LockFileEx` module, torn-tail truncation on open, periodic compaction through DurableStore | O(1) append; matches ADR-001 §4, which explicitly separates *document persistence* from *workflow journaling* |

**Proposed (for review, not adopted):** a **separate, small journaling primitive** that composes with DurableStore
rather than replacing it — same `Result` vocabulary, same lock module, same refusal semantics — because an
append-only ledger is a journal, and ADR-001 §4 already forbids conflating the two.

**The decisive evidence is missing, and it is measurable:** the distribution of events per conversation (p50/p95/max)
and the interval between turns. If p95 is small and turns are seconds apart, **as-is** is defensible and cheaper.
That measurement must be taken before this question is decided — it is the only input that changes the answer.

**Regardless of shape**, four properties must be decided: record framing and checksum · flush policy per commit
(per event, or batched with an explicit durability claim) · the truncation rule on a torn tail (truncate to the last
valid record and **record a recovery event**, never silently continue past it) · duplicate suppression when a commit
landed but its acknowledgement did not.

**Power-loss durability remains `UN` here as it does for TASK-02**, and the ledger must not claim more.

## 7. Assessment C — ContainmentChecker disposition

The owner's concern is correct: a second containment implementation would repeat the defect TASK-01 removed.

**Verified.** `internal/paths` now provides `Roots`, `RootSet.Check`, `CheckContained`, `CheckGrant`,
`ExternalRoot`, with physical resolution by `GetFinalPathNameByHandle` (TASK-01). It answers *where* and *is this
inside*. It does **not** expose component-level reparse inspection, and it does not pin a directory identity across
an operation.

**Gap versus ADR-001 §7.** Grant-time inspection that reports *which* component is a reparse point · a pinned
parent handle retained for an operation · rejecting reparse points created **after** issuance (TOCTOU) ·
handle-relative creation (`NtCreateFile`; capability `UN`). None of these is needed by the ledger, which writes only
inside an app-owned root.

**Duplication that already exists.** `hasReparseComponent` is implemented **four times**:
`internal/audio/tts/helper/pathcheck_windows.go:16`, `pathcheck_other.go:14`, and the same pair under
`fakehelper` (`:13`, `:11`). ADR-001 §7 already directs that this be promoted into one shared pre-flight primitive.

**Disposition (proposed).** Do **not** create a `ContainmentChecker` package. Extend `internal/paths` additively
when ExternalDataGuard work begins, and retire the four copies then. For TASK-03 the conclusion is narrow and
actionable: **the ledger adds no containment code**; it calls `RootSet.Check` for its own documents. Recorded here so
the duplication cannot reappear silently.

## 8. Acceptance criteria for the eventual implementation

Mapped to FC-1 §11, scoped to Phase 1B. Each is observable and needs a test, not a statement.

| # | Criterion |
|---|---|
| A-1 | A conversation's identity survives restart, rename and file movement, and is never derived from list position or filename order |
| A-2 | An append is atomic: crash mid-append leaves the ledger parseable, with the event either present and valid or absent — never half-applied |
| A-3 | A duplicate commit with the same idempotency key produces one durable effect; a conflicting reuse is refused |
| A-4 | Replaying the ledger reproduces the same state/message/turn identities and the same current selection |
| A-5 | Selection is derived, never stored as truth; a restart reproduces it exactly |
| A-6 | Deleting the selected state records the resulting selection explicitly and resolves to a valid state or an explicit empty/unresolved condition — no dangling pointer, no silent fallback |
| A-7 | No committed event is edited in place; correction is a new event. `RewriteHistory`-style replacement is gone |
| A-8 | Corrupt projection + valid ledger rebuilds and reports reconciliation; an unusable ledger refuses explicitly rather than presenting an empty conversation |
| A-9 | Two processes cannot interleave appends to one conversation: the second reports contention, never a merge and never a silent success |
| A-10 | The four retired implementations are removed and the application still functions — the anti-sixth-authority criterion |
| A-11 | Deleting a conversation deletes its ledger and snapshots under a receipt; authored content is untouched (FC1-22b) |

**Not claimed by any of these:** power-loss durability · branch/fork semantics · snapshot contents · memory copies.

## 9. Handoff

| Field | Content |
|---|---|
| **Starting point** | This record; `FC-1` §1–§8, §10–§11; ADR-001 §4, §5 D, §8–§9; ADR-002 §4; `internal/durable` as the locking/result vocabulary; `internal/paths` for containment |
| **Ownership** | §2 |
| **Acceptance** | §8 |
| **Evidence** | §3 inventory (verified), §5 triage, §6 assessment + the measurement it needs, §7 disposition |
| **Next action** | Resolve D-1…D-4 into one decision record, then owner review. **No code until that review passes** |

## 10. Authorization

Prepares investigations and records three assessments. **Authorizes no implementation, no migration, no call-site
change, and no contract approval.** FC-1 remains `PROPOSED`; implementing a clause does not approve it.
