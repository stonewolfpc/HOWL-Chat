# TASK-03 — clarification 1a: deletion journal, readers, actor matrix, terminal cleanup

| Field | Value |
|---|---|
| Purpose | Close four remaining consistency findings in clarification 1 |
| Status | **Contract text. Nothing adopted. No implementation authorized** |
| Relationship | **Additive.** Neither `TASK-03-CONTRACT-ACCEPTANCE-2026-11-22.md` nor `TASK-03-CONTRACT-CLARIFY-1-2026-11-22.md` is edited. Named clauses are superseded, not rewritten |
| Supersedes | K-02, K-04, K-06, K-11 (third element), K-21 — and adds journal integrity obligations the earlier text left to implementation |
| Recovery token | `TASK-03-CONTRACT-CLARIFY-1A-2026-11-22` |

Prefixes: `L-` normative, `LT-` acceptance test.

---

## 1. Proven-never-created versus missing: K-04 and K-06 describe one condition

**Finding, accepted.** K-04 infers `aborted` from a missing entry; K-06 infers `indeterminate` from the same observation. The write order `intent → journal → remove` was meant to make absence proof that removal never started, but absence of a *durable* record cannot distinguish "never written" from "written and then lost, corrupted, or unreadable". **Absence is not proof of non-creation.** The clause is wrong, not merely ambiguous.

**L-01 — The write order is inverted, so absence acquires a meaning.**

```
observe presence  →  journal{prepared}  →  intent  →  journal{pending}
                  →  remove  →  journal{removing}  →  verify absence
                  →  journal{removed}  →  terminal receipt  →  delete entry
```

Because the entry is created **before** the intent, an intent with no entry can never mean "the entry was never created". It means the entry was created and later became unavailable. That single inversion is what makes the inference sound.

**L-02 — Entry schema, version 1.**

| Field | Type | Note |
|---|---|---|
| `journal_v` | int | `1`; higher → `journal_unversioned` |
| `operation_id` | UUIDv4 | joins to the receipt |
| `directory_path` | string | locator; sensitive |
| `observed_present_at_start` | bool | observed before any removal |
| `record_count`, `snapshot_count` | int | accounting copies |
| `state` | `prepared` \| `pending` \| `removing` \| `removed` | durable progress |
| `updated_at` | RFC3339 UTC | progression evidence |

**L-03 — The inference table. This replaces K-04 and K-06.**

| Receipt | Entry | Directory | Result | Ordinary meaning |
|---|---|---|---|---|
| none | `prepared` | any | `abandoned_no_intent` — report, remove entry, **write no receipt** | No deletion was requested; removal could not have started |
| intent | **missing** | any | **`indeterminate`, `locator_unavailable`** | Entry created then lost — progress unknowable. **Never `aborted`** |
| intent | unparseable / invalid | any | `indeterminate`, `locator_corrupt` | Entry retained, not deleted |
| intent | `journal_v` unknown | any | `indeterminate`, `journal_unversioned` | Refuse to interpret |
| intent | `pending` | absent | `complete`, `attribution: unknown` | Removal may have run; authorship unknown |
| intent | `pending` | present, intact | `aborted`, `removed_record_count: 0` | Removal had not started |
| intent | `pending` / `removing` | present, partly emptied | `partial`, `attribution: operation` | Partial effect observed |
| intent | `removing` | absent | `complete`, `attribution: unknown` | Removal began; absence unverified by us |
| intent | `removed` | absent | `complete`, `attribution: operation` | The operation itself recorded verified absence |
| terminal | entry present | any | stale entry → remove, return receipt (§4) | Normal crash tail |
| terminal | entry missing | any | return receipt, write nothing | Already clean |

**L-04 — The safety property, stated once.**

> Every failure of the pending-deletion journal — missing, unreadable, truncated, corrupt, partially written, unknown version, or lost — degrades the outcome toward **`indeterminate`**. No journal failure may produce `aborted`, `absent_prior`, `complete`, or any claim that a deletion occurred or did not occur.

Truncation is the case that matters most: a shortened journal still parses, so an entry can vanish silently. The table routes that to `locator_unavailable`, which claims nothing.

**L-05 — A truncated or missing journal never licenses a scan.** The directory is not discovered by filename pattern, count match, timestamp proximity, or directory listing. L-03's `indeterminate` is the terminal answer until the owner supplies evidence out of band.

**Alternatives considered.**
- *Keep `intent → journal` and treat absence as proof.* Rejected: it is the defect being fixed. A lost entry is indistinguishable from an uncreated one.
- *Add a `journal_committed: true` flag to the receipt.* Rejected: the receipt is immutable and written before the journal under the old order; adding a later flag creates a second mutable record with the same trust problem.
- *Reconcile by scanning for candidate directories and matching counts.* Rejected: it reconstructs the identity the receipt deliberately omits, and would attach an outcome to the wrong conversation when two deletions share a count.
- *Treat every missing entry as `indeterminate` but still allow startup to guess.* Rejected: guessing is the failure mode, wherever it happens.

**LT-01** — Intent present, entry missing → `indeterminate`/`locator_unavailable`; never `aborted` (asserted across intact, absent and partial directories).
**LT-02** — Journal truncated after the intent and before `pending`: the case parses, the entry is gone, and the result is `indeterminate`.
**LT-03** — Entry with `journal_v: 2` → `indeterminate`/`journal_unversioned`; nothing is rewritten.
**LT-04** — `prepared` entry with no intent → `abandoned_no_intent`, entry removed, receipt file byte-identical.
**LT-05** — No journal failure path yields `aborted`, `absent_prior` or `complete` (table-driven over every corruption case).
**LT-06** — Reconciliation performs no directory listing or pattern scan (asserted by instrumenting filesystem calls).

---

## 2. Two readers, named: K-02 and K-05

**Finding, accepted.** K-02 says only reconciliation reads the journal; K-05 has startup reading it. Contradiction as written.

**L-06 — Replaces K-02.**

> The pending-deletion journal has exactly **two** readers, and no others:
> 1. **Startup discovery** — read-only. Reports pending `operation_id`s. Never reconciles, never writes, never removes an entry, and never surfaces a path.
> 2. **Explicit reconciliation** — read-only inspection, followed by the authorized writes of C-36: one terminal receipt record and removal of that operation's entry.
>
> No other component — UI, API, diagnostics, export, migration, or test harness — may open the file.

**L-07 — Prohibited surfaces.** No directory path appears in the UI, any API response, diagnostics bundles, exports, or logs. Startup reports operation identifiers only. An error raised while reading the journal names the failure class, never the path.

**L-08 — Startup never reconciles.** Discovery is passive. Acting requires the explicit call; an automatic reconcile at startup is outside the contract and is a defect, not an optimization.

**Alternatives considered.**
- *Make reconciliation the only reader and let startup report nothing.* Rejected: a pending deletion would then be invisible after a restart until someone knew to ask — the journal exists precisely to survive that gap.
- *One shared reader with a write flag.* Rejected: it makes "who may write" a runtime parameter instead of a stated boundary, which is how the original contradiction arose.

**LT-07** — A startup pass over a journal with two pending entries performs zero writes and returns two `operation_id`s and no paths.
**LT-08** — The journal file is opened by no other module (import-graph and call-site assertion).
**LT-09** — No log line, API payload or export contains `directory_path` or any path substring.

---

## 3. The actor matrix must be per-kind: K-21 against K-22/K-23

**Finding, accepted.** K-21's row `system | any → ACTOR_NOT_PERMITTED` is categorical; K-22 and K-23 permit `system` on `conversation_created` and `selection_changed`. The matrix was written for turn events and then generalised by the word "any". The restriction was always about **turns**, not about ledger authorship.

**L-09 — Replaces K-21. Conditional on event kind.**

| Kind | `actor` | `role` | Result |
|---|---|---|---|
| `conversation_created` | `user` | — | accept |
| `conversation_created` | `system` | — | **accept** |
| `conversation_created` | `assistant` | — | `ACTOR_NOT_PERMITTED` |
| `conversation_created` | `recovery` | — | `ACTOR_NOT_PERMITTED` |
| `turn_committed` | `user` | `user` | accept |
| `turn_committed` | `assistant` | `assistant` | accept |
| `turn_committed` | `user` | `assistant` | `ROLE_ACTOR_MISMATCH` |
| `turn_committed` | `assistant` | `user` | `ROLE_ACTOR_MISMATCH` |
| `turn_committed` | `system` | any | `ACTOR_NOT_PERMITTED` |
| `turn_committed` | `recovery` | any | `ACTOR_NOT_PERMITTED` |
| `turn_corrected` | `user` \| `assistant` | — | accept if `actor` equals the target turn's role |
| `turn_corrected` | `system` \| `recovery` | — | `ACTOR_NOT_PERMITTED` |
| `turn_retracted` | `user` \| `assistant` | — | accept if `actor` equals the target turn's role |
| `turn_retracted` | `system` \| `recovery` | — | `ACTOR_NOT_PERMITTED` |
| `selection_changed` | `user` | — | accept |
| `selection_changed` | `system` | — | **accept** |
| `selection_changed` | `assistant` | — | `ACTOR_NOT_PERMITTED` |
| `selection_changed` | `recovery` | — | `ACTOR_NOT_PERMITTED` |
| any kind | unknown value | — | `ACTOR_INVALID` |
| `turn_committed` | any | unknown value | `ROLE_INVALID` |

**L-10 — Corrections and retractions carry no `role`.** Their authority is derived: the `actor` must equal the role of the turn named by `target_turn_id`. A mismatch is `ROLE_ACTOR_MISMATCH`; a missing target is `TARGET_NOT_FOUND`, evaluated before the actor check so the error is deterministic.

**L-11 — The rule behind the table.**

> `system` may author events that describe the application's own state — creation and selection. It may not author turns, corrections or retractions, because those carry user or model content. `recovery` authors no ledger record in any kind. `assistant` authors only assistant-side turns and their corrections and retractions.

**L-12 — Correction recorded, not erased.** K-21's categorical reading applied to turn kinds only. K-22 and K-23 were correct; the matrix was under-specified. Clarification 1's text stands as written, with this table governing.

**Alternatives considered.**
- *Amend K-22/K-23 to forbid `system` entirely.* Rejected: creation and selection legitimately belong to the application, and forbidding them would make an app-driven conversation impossible to record.
- *Allow `system` on turns for automation.* Still rejected: no automation is in scope, and it would let generated text be indistinguishable from a model reply.
- *Drop the `role`/`actor` distinction.* Rejected: `role` is what the user sees and mismatch detection is the point.

**LT-10** — All nineteen rows of L-09 produce the stated result with distinct errors.
**LT-11** — A `system`-authored `conversation_created` and a `system`-authored `selection_changed` commit successfully.
**LT-12** — A `system`-authored `turn_committed` is refused with `ACTOR_NOT_PERMITTED`; the file is unchanged.
**LT-13** — A `turn_corrected` whose `actor` differs from the target turn's role is refused with `ROLE_ACTOR_MISMATCH`; an unknown target yields `TARGET_NOT_FOUND` first.

---

## 4. Crash after the terminal receipt, before entry removal

**Finding, accepted.** The sequence ends `terminal → delete entry`. A crash between those two leaves a receipt and a stale entry, and clarification 1's lifetime claim ("exists only while the outcome is unknown") is then false unless a rule restores it.

**L-13 — Terminal first, then the entry. Replaces K-16's ordering assumption.**

> Reconciliation inspects the terminal receipt **before** the journal. If a terminal record exists for the `operation_id`, it is returned as-is and no second receipt is written.

**L-14 — Stale entry cleanup.**

| Receipt | Entry | Action | Result |
|---|---|---|---|
| terminal exists | entry present | remove entry | return the existing receipt unchanged |
| terminal exists | entry missing | none | return the existing receipt unchanged |
| terminal exists | entry unreadable | remove entry | return the existing receipt unchanged; report `entry_unreadable` |

**L-15 — Never repeat the deletion.** Reconciliation removes journal entries and writes at most one receipt. It **never** touches a conversation directory: not to delete, not to verify, not to re-create. A receipt that already says `complete` is not re-verified by reconciling — re-verification would be a new operation with its own receipt.

**L-16 — Lifetime restored.** After `L-14`'s cleanup the entry is gone, so clarification 1's claim holds: an entry exists exactly while the outcome is unknown, plus at most the interval between a terminal write and its cleanup.

**L-17 — Idempotence.** A second call after cleanup returns the same receipt, writes nothing, removes nothing. Reconciliation is safe to retry after any crash.

**Alternatives considered.**
- *Write the terminal receipt only after deleting the entry.* Rejected: it inverts the evidence order, so a crash in between would lose the receipt and leave a deletion no record describes.
- *Have reconciliation delete the directory again to confirm.* Rejected explicitly by C-36 and by the finding; it would also destroy a recreated conversation under the same path.
- *Report a stale entry as an error.* Rejected: it turns a normal crash tail into a permanent alarm, and the correct disposition is known.

**LT-14** — Crash simulated after the terminal write and before entry removal: the next call returns that receipt, writes no second receipt, removes the entry, and performs no filesystem change inside the conversation directory.
**LT-15** — Two consecutive calls return the same receipt and the second performs zero writes.
**LT-16** — A `complete` receipt is never re-verified and the directory is never re-deleted (asserted on filesystem instrumentation).

---

## 5. The journal is part of the deletion contract

Accepted: the journal is no longer an implementation convenience. It is a durable, sensitive recovery mechanism that temporarily holds a conversation path, and its failure modes are contractual.

**L-18 — Integrity.**

| Requirement | Rule |
|---|---|
| Atomicity | Every journal write replaces the file atomically (temp, flush, rename). No in-place edit |
| Serialization | A journal lock covers each read-modify-write; a second writer waits or reports `Busy` |
| Versioning | `journal_v` per entry; unknown version → `journal_unversioned`, never interpreted |
| Validation | An entry missing any required field is `locator_corrupt`; the file is never partially trusted |
| Truncation | A shortened file is read as-is and produces `locator_unavailable` for affected intents — never a fallback source |
| Corruption | Unparseable file → `journal_corrupt`; every pending intent becomes `indeterminate` |

**L-19 — Lifetime and retention.** An entry exists from `prepared` until the terminal receipt is written and the entry removed. There is no expiry, no timer, no background sweep, and no automatic repair. An entry survives restarts indefinitely until an explicit reconciliation ends it.

**L-20 — Sensitivity.** The journal contains conversation paths. It is excluded from every export, diagnostics bundle and backup the contract controls; it is never logged; its filesystem permissions are owner-only where the platform allows. Loss of the journal never causes a claim in either direction — it only removes the ability to conclude, which L-04 already routes to `indeterminate`.

**L-21 — Replaces K-11's third element.** `attribution: operation` requires: observed presence before removal, **and** either a successful removal result in this process **or** a durable `removed` journal marker written after verifying absence. Absent both, a `pending`/`removing` entry with an absent directory yields `attribution: unknown`; `prior_absence` belongs only to `absent_prior`.

**Alternatives considered.**
- *One journal file per operation.* Rejected: it multiplies sensitive files and needs its own discovery scheme.
- *Store the entry inside the receipt file.* Rejected: it puts the path next to the receipt, which is the one place it must not be.
- *Give the journal a TTL and let expiry clean it.* Rejected: a timer is a sweep, and expiry would silently convert an unknown outcome into no record.

**LT-17** — A crash between `prepared` and `pending` leaves one entry and no intent; reconciliation abandons it without writing a receipt.
**LT-18** — Concurrent journal writes do not interleave; the loser reports `Busy` and no entry is lost.
**LT-19** — No export, log or diagnostics bundle contains journal content.
**LT-20** — `attribution: operation` appears only for a successful in-process removal or a `removed` marker.

---

## 6. Delta consistency matrix

Rows are the clauses clarification 1a touches, plus the clauses those touch.

| Clause | 1a disposition | Now coherent with |
|---|---|---|
| K-01 locator | schema extended (L-02), mechanism unchanged | L-01, L-18 |
| **K-02** readers | **superseded** by L-06/L-07/L-08 | C-36, C-41 |
| **K-04** order + aborted | **superseded** by L-01/L-03 | C-34, C-35 |
| **K-06** locator loss | **superseded** by L-03 row 2 | L-04, L-05 |
| K-05 restart | unchanged, refined by L-08 | L-06 |
| K-07 privacy | unchanged, strengthened by L-20 | L-07 |
| K-09 `absent_prior` | unchanged; no entry, no intent | L-03 row 1 |
| K-10 outcome table | refined by L-03 | K-11/L-21 |
| **K-11** attribution | **third element replaced** by L-21 | L-03, L-14 |
| K-12 honest absence | unchanged, supported by L-04 | C-35 |
| **K-21** actor matrix | **superseded** by L-09/L-10 | K-22, K-23, C-09 |
| K-22 authority | unchanged; L-11 states the rule it implied | L-09 |
| K-23 per-kind actors | unchanged; now consistent with L-09 | L-09 |
| K-16 reconcile ordering | refined by L-13 (terminal first) | L-14 |
| C-34 lifecycle | unchanged; journal is evidence, receipt is the record | L-01 |
| C-35 crash window | unchanged; `W` now bounded by the `removed` marker | L-21 |
| C-36 bounded reconcile | unchanged; "deletes nothing" clarified by L-15 (entries only) | L-14, L-15 |
| C-41/C-45 evidence | unchanged | L-19, L-20 |
| C-31 byte preservation | untouched | ✓ |

**Event-kind count is unchanged by 1a: five ledger kinds plus two evidence kinds.** L-09's table has one row group per existing kind; no kind is added or removed.

### Residual items

1. `selection_resolved` and D-3's ancestor rule remain deferred to a later contract.
2. The journal holds a conversation path for the life of an unresolved deletion — now bounded, versioned and swept only by explicit reconciliation, but still sensitive.
3. Power-loss durability remains `UN`; workload remains `Unestablished`.
4. `abandoned_no_intent` is a reported status, not a receipt state; the receipt vocabulary is unchanged.

### Authorization

No implementation, migration, ledger creation, persistence adoption, or FC-1 approval. `go build` passing and documents agreeing are not implementation evidence; no `LT-` or `T-` test has been executed, because no implementation exists.
