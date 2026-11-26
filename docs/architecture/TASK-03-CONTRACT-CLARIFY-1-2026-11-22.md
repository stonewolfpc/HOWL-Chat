# TASK-03 — contract 1.0-rc, clarification 1

| Field | Value |
|---|---|
| Purpose | Dispose of five consistency findings against the 53-clause contract |
| Status | **Contract text for acceptance. Nothing adopted. No implementation authorized** |
| Relationship | **Additive.** Every clause here either adds `K-` clauses or replaces exactly one named clause's wording. `TASK-03-CONTRACT-ACCEPTANCE-2026-11-22.md` is **not edited** |
| Supersedes | C-09 (actor enum), C-11, C-18, C-19 (`selection_resolved` row, `turn_committed` row), C-26, C-33, C-34, C-35, C-36, C-44 (unchanged), C-52 — each named at its clause |
| Recovery token | `TASK-03-CONTRACT-CLARIFY-1-2026-11-22` |

Clause prefixes: `K-` normative, `KT-` acceptance test.

---

## 1. Deletion reconciliation identity without a conversation ID

**Finding.** C-36 lets a reconciliation call name an `operation_id`, but C-33 forbids the conversation ID in the receipt. After a restart nothing maps the two. As written, reconciliation cannot locate the directory it is reconciling.

**K-01 — The operation locator.** A transient app-state journal, `evidence/pending_deletions`, holds one entry per deletion in flight:

| Field | Type | Purpose |
|---|---|---|
| `operation_id` | UUIDv4 | joins to the receipt |
| `directory_path` | string | locates the conversation |
| `prepared_at` | RFC3339 UTC | ordering evidence |
| `record_count`, `snapshot_count` | int | copied from the intent for a later accounting |

**K-02 — Authority.** The deletion path writes the entry; **only** the reconciliation call reads it. It is never exposed through the UI, any API, diagnostics bundles or exports, and it is never written into a ledger. An entry is deleted when the terminal receipt record is written.

**K-03 — Lifetime.** From the intent record until the terminal record — i.e. exactly the window in which the outcome is unknown. If the process dies, the entry persists; that persistence is the mechanism by which a later reconciliation can work, not an oversight.

**K-04 — Ordering, and what its absence proves.**

```
observe presence  →  write intent (receipt)  →  write journal entry  →  remove  →  verify absence  →  write terminal  →  delete journal entry
```

Because the journal entry is written **before** removal, a missing entry with a present intent proves removal **never started**, which yields `aborted` rather than `indeterminate`.

**K-05 — Behaviour after restart.** Startup reads the journal **read-only** and reports pending `operation_id`s. Reconciliation remains an explicit call (C-36). It reads the path from the journal, observes the filesystem, writes one terminal record, and deletes the entry.

**K-06 — Locator loss.** If an intent has no journal entry and no terminal record, the outcome is `indeterminate` with reason `locator_missing` and attribution `unknown`. The directory is **never** guessed by scanning for candidates, and never matched by counts or timestamp — that would reconstruct the identity the receipt deliberately omits.

**K-07 — Privacy.** The journal is the only place a conversation path is stored for a deletion. Its entry is removed on completion. **Residual risk, stated:** after a crash inside the window, a path persists in app-state until reconciliation runs; that file is therefore treated as sensitive as the conversation itself, is excluded from every export, and is not backed up by this contract's own mechanisms.

**Alternatives considered.**
- *Put a directory reference in the receipt.* Rejected: C-33, and timestamp-named directories are enumerable, so any derivable reference is a lookup path.
- *Require the caller to supply the path.* Works only in-process; after restart the caller has nothing, so the outcome would stay unknown forever, contradicting C-36.
- *Store a salted hash of the path.* Rejected again for the same reason as the salted conversation hash: salt beside hash is linkage.
- *A separate durable mapping file with long retention.* Rejected as a second identity store; the journal is scoped to the window and deleted at the terminal record.

**KT-01** — A deletion completes: journal entry exists during the operation and is gone after the terminal record; the receipt contains no path or identifier.
**KT-02** — Kill the process inside the window, restart: startup reports exactly one pending `operation_id` and writes nothing; the explicit call writes one terminal record and removes the entry.
**KT-03** — An intent without a journal entry yields `indeterminate`/`locator_missing`, and no directory is scanned or guessed.
**KT-04** — The journal file appears in no export, diagnostic bundle or API response.

---

## 2. Prior absence versus completed removal

**Finding.** C-35 permits `complete` when the directory is absent, without distinguishing "we removed it" from "it was already gone". Later absence alone must not prove this operation deleted anything.

**K-08 — Presence is observed before removal and recorded.** The deletion path observes presence first and writes it durably (`observed_present_at_start`) before any removal.

**K-09 — `absent_prior` is its own terminal state.**

> If the directory is absent when the operation begins, the operation performs no removal and writes a single terminal record `absent_prior`, `attribution: prior_absence`, `removed_record_count: 0`. No intent and no journal entry are written, because no deletion was started.

**K-10 — Permitted outcomes, by evidentiary case.** These are all the permitted combinations; any other claim is a violation.

| Presence observed first | Removal result | Directory now | Terminal state | Attribution |
|---|---|---|---|---|
| absent | not attempted | absent | `absent_prior` | `prior_absence` |
| present | success (in this process) | absent | `complete` | `operation` |
| present | no result (crash in `W`) | absent | `complete` | `unknown` |
| present | failure | present, intact | `aborted` | — |
| present | failure | present, partly emptied | `partial` | `operation` |
| present | not recorded | absent or present | `indeterminate` | `unknown` |

**K-11 — The attribution rule, stated once.**

> `attribution: operation` requires **all three**: observed presence before removal, a successful removal result in this process, and observed absence afterwards. Any missing element downgrades the attribution to `unknown`; an operation that began with an absent directory may only claim `prior_absence`.

**K-12 — Honest absence.** A receipt whose directory is absent does not prove this installation removed it. Reporting "deleted by this operation" from absence alone — including after a crash in `W` — is a contract violation, not an inference.

**Alternatives considered.**
- *Treat ENOENT from `rm` as success.* Rejected: it collapses "we removed it" into "it was already gone", which is the exact overclaim this finding names.
- *Omit attribution entirely.* Rejected: it hides the distinction rather than resolving it, and makes `absent_prior` indistinguishable from a completed deletion.
- *Re-check presence only at reconciliation.* Rejected: after a crash the earlier state is unrecoverable, so presence must be recorded before removal.

**KT-05** — Deleting an already-absent directory writes exactly one `absent_prior` record, no intent, no journal entry, and never `complete`.
**KT-06** — A crash in `W` yields `complete` with `attribution: unknown`, never `operation`.
**KT-07** — Every terminal record's state/attribution pair appears in K-10's table; an unlisted pair fails the test.
**KT-08** — A generated receipt never states that this installation performed a deletion unless K-11's three elements hold.

---

## 3. Selection determinism

**Finding.** C-19 says `turn_committed` means "selection follows", while C-26 says selection is the last selection event, else the newest state. A new state-producing event after an explicit selection is undefined by one clause and pinned by the other.

**K-13 — Replaces C-26.**

> Let `j` be the sequence of the last selection-affecting event (`selection_changed`) at or before `k`.
> - If a state-producing event exists in `(j, k]`, then `selection(k)` = the **newest state** at `k`.
> - Otherwise, if `j` exists, `selection(k)` = the state named by that event.
> - Otherwise, if any state exists at `k`, `selection(k)` = the newest state.
> - Otherwise `selection(k)` = `empty`.

An explicit selection therefore **persists until the next state-producing event**, which re-selects the head. "Selection follows" in C-19 is this head-follow, nothing more.

**K-14 — Consequence, stated.** A user viewing an earlier state who receives a new committed turn is moved to the new head. This is deliberate: the alternative leaves new content invisible. Phase 1 has no "hold position" mode.

**K-15 — Worked example.** Ledger `e1..e9`; states in bold.

| Seq | Event | Newest state | Selection after |
|---|---|---|---|
| 1 | `conversation_created` | — | `empty` |
| 2 | `turn_committed` (user) | **S2 = e2** | S2 |
| 3 | `turn_committed` (assistant) | **S3 = e3** | S3 |
| 4 | `selection_changed` → S2 | S3 | **S2** (pinned) |
| 5 | `turn_committed` (user) | **S5 = e5** | S5 (pin cleared) |
| 6 | `turn_corrected` (e2) | **S6 = e6** | S6 |
| 7 | `selection_changed` → S3 | S6 | **S3** (pinned) |
| 8 | `turn_retracted` (e3) | **S8 = e8** | S8 |
| 9 | `selection_changed` → S5 | S8 | **S5** (pinned) |

**Alternatives considered.**
- *Pins persist until explicitly changed.* Rejected: a new turn would not become visible, which contradicts the intent of C-19 and produces a conversation that appears frozen.
- *Pins survive corrections but not new turns.* Rejected as unprincipled: both produce states, and special-casing kinds re-creates the ambiguity a deterministic rule is meant to remove.
- *Selection stored as a field on the projection.* Rejected: C-21 requires the projection to be a pure function of the sequence; a stored pointer is an input that is not an event.

**KT-09** — Replaying the K-15 ledger produces exactly the listed selections at every prefix.
**KT-10** — A `selection_changed` followed by a state-producing event ends at the new head, not the pinned state.
**KT-11** — Two `selection_changed` events with no state-producing event between them resolve to the later target.

---

## 4. `selection_resolved` has no Phase 1 trigger

**Finding.** C-19 defines `selection_resolved` around `deleted_state_id` and an `ancestor` rule. C-22 makes every state permanently addressable, and C-03 excludes state deletion. **Nothing in Phase 1 deletes a state**, so the event has no valid trigger; retaining it means specifying a branch that can never run.

**K-16 — Finding confirmed.** No Phase 1 operation removes a state, a turn, or a state's addressability. `selection_resolved` is therefore **deferred out of the Phase 1 vocabulary**: five ledger kinds, not six.

**K-17 — Vocabulary, replaces C-18.**

> Phase 1 `event_type` is one of: `conversation_created`, `turn_committed`, `turn_corrected`, `turn_retracted`, `selection_changed`. A record carrying any other `event_type` — including `selection_resolved` — is **refused at commit** (`OUT_OF_VOCABULARY`) and **`unresolved` at read**.

**K-18 — C-26 simplification.** With no resolver, selection changes only via `selection_changed`, which cannot express `empty`; `empty` is reachable only as the initial state of a conversation with no state-producing events.

**K-19 — D-3 is re-homed, not dropped.** D-3's nearest-surviving-ancestor selection concerns a **conversation that may no longer exist**. It therefore cannot be an event in that conversation's ledger, which is why it was misfiled under `selection_resolved`. It belongs to conversation-level selection in app-state, outside this contract, and is recorded here as an **open item for a later contract**, not as a Phase 1 ledger event.

**K-20 — Effect on acceptance tests.** C-19's `selection_resolved` row is removed; T-07's `rule`/`chosen_state_id` cases are removed; the remaining five kinds' tests are unchanged; `KT-12` covers vocabulary rejection. C-52's second owner decision is restated as **five ledger kinds plus two evidence kinds**.

**Alternatives considered.**
- *Keep `selection_resolved` unused.* Rejected: an event with no trigger accumulates a specification nobody can test, and invites an implementation to invent a deletion path to justify it.
- *Keep it and add a state-deletion operation.* Rejected: contradicts C-22's permanent addressability and C-03, and would destroy the historical-state guarantee the owner just secured.
- *Reduce vocabulary to four kinds and drop selection entirely.* Rejected: selection is needed once a conversation has more than one state; that need is real, unlike the resolver's.

**KT-12** — A `selection_resolved` record is refused at commit with `OUT_OF_VOCABULARY` and leaves the file unchanged; at read it marks the ledger `unresolved`.
**KT-13** — The enum accepts exactly the five kinds and rejects a sixth.
**KT-14** — `empty` selection is reachable only before the first state-producing event.

---

## 5. Actor/role mapping for `turn_committed`

**Finding.** C-09 allows four `actor` values; C-19 requires `role` on turn events and states `actor` must "agree" without defining agreement. The combination matrix is missing, so invalid input is not deterministically testable.

**K-21 — Replaces C-19's `turn_committed` row note with an explicit matrix.**

| `actor` | `role` | Result |
|---|---|---|
| `user` | `user` | **accept** |
| `assistant` | `assistant` | **accept** |
| `user` | `assistant` | refuse `ROLE_ACTOR_MISMATCH` |
| `assistant` | `user` | refuse `ROLE_ACTOR_MISMATCH` |
| `system` | any | refuse `ACTOR_NOT_PERMITTED` |
| `recovery` | any | refuse `ACTOR_NOT_PERMITTED` |
| any other | any | refuse `ACTOR_INVALID` |
| any | any other | refuse `ROLE_INVALID` |

**K-22 — Authority for `system` and `recovery`.**

> In Phase 1, `system` may author non-turn ledger events (`conversation_created`, `selection_changed`) and may not author turns. `recovery` authors no ledger record at all; its value exists for evidence records only. No Phase 1 path may create a user- or assistant-role turn under another actor's authority.

**K-23 — Per-kind actor allowances.**

| Kind | Permitted `actor` |
|---|---|
| `conversation_created` | `user`, `system` |
| `turn_committed`, `turn_corrected`, `turn_retracted` | `user`, `assistant` (matching `role` on commit) |
| `selection_changed` | `user`, `system` |

**Alternatives considered.**
- *Allow `system` to create assistant turns.* Rejected for Phase 1: it would let automation write content indistinguishable from a model reply, and there is no automation in scope yet.
- *Allow `recovery` to write turns when repairing.* Rejected: repair would become an authored-content path with no authorship boundary; recovery reports, it does not write (C-41, C-45).
- *Drop `role` and derive it from `actor`.* Rejected: `role` is what the user sees; deriving it removes the ability to detect a mismatch, which is the point of the test.

**KT-15** — All eight rows of K-21 produce the stated result, each with a distinct error where refused.
**KT-16** — A `system`-authored `turn_committed` is refused with `ACTOR_NOT_PERMITTED` and leaves the file unchanged.
**KT-17** — A `recovery`-authored ledger record of any kind is refused.

---

## 6. Consistency matrix

`✓` coherent unchanged · `A` amended by this document · `—` not involved.

| Clause | 1 locator | 2 absence | 3 selection | 4 resolver | 5 actor |
|---|---|---|---|---|---|
| C-09 `actor` enum | — | — | — | A (K-22) | **A** (K-21, K-23) |
| C-10 conditional fields | — | — | — | — | ✓ |
| C-11 absent fields | ✓ | ✓ | — | — | — |
| C-18 kind set | — | — | — | **A** (K-17) | — |
| C-19 schemas | — | — | **A** (K-13) | **A** (row removed) | **A** (K-21) |
| C-20 invalid events | ✓ | — | — | ✓ | ✓ |
| C-21 determinism | ✓ | — | **A** (K-13, K-14) | — | — |
| C-22 addressability | ✓ | — | — | **A** (K-16 confirmed) | — |
| C-26 selection function | — | — | **A** (replaced) | **A** (K-18) | — |
| C-31 byte preservation | ✓ | — | — | ✓ | — |
| C-32 Option A | ✓ | ✓ | — | — | — |
| C-33 receipt fields | ✓ | ✓ | — | — | — |
| C-34 lifecycle | ✓ | **A** (K-08) | — | — | — |
| C-35 crash window | ✓ | **A** (K-10, K-11) | — | — | — |
| C-36 reconciliation | **A** (K-01–K-06) | **A** (K-10) | — | — | — |
| C-37 correlation | ✓ | ✓ | — | — | — |
| C-41 / C-45 evidence | ✓ | — | — | — | **A** (K-22 cross-ref) |
| C-49 turn identity | — | — | ✓ | — | **A** (K-21) |
| C-52 decisions | — | — | — | **A** (K-20) | — |
| D-1 (substrate) | ✓ | ✓ | — | ✓ | ✓ |
| D-2 (lock/sequence) | ✓ | ✓ | — | ✓ | ✓ |
| D-3 (ancestor rule) | — | — | — | **A** (re-homed, K-19) | — |
| D-4 (no journal) | ✓ | ✓ | — | ✓ | ✓ |
| `FC1-14b` linear | — | — | ✓ | ✓ | — |
| `FC1-14c` one/turn | — | — | ✓ | — | ✓ |

**No remaining contradiction found.** The two places where the contract previously contradicted itself — selection precedence and the resolver's missing trigger — are removed rather than papered over; the two places where it could not be executed — reconciliation identity and absence attribution — are given mechanisms. Nothing in this document weakens C-21, C-22 or C-31.

### Residual items carried forward

1. `selection_resolved` and D-3's ancestor rule need a later contract (§4, K-19).
2. The pending-deletions journal holds a path during the window (§1, K-07).
3. Power-loss durability remains `UN`; workload remains `Unestablished`.
4. Nine owner decisions, with decision 2 restated as five ledger plus two evidence kinds.

### Authorization

No implementation, migration, ledger creation, persistence adoption, or FC-1 approval by implication. This is contract review only.
