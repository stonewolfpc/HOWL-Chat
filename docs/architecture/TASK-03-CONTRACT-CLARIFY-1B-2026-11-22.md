# TASK-03 — clarification 1b: partial attribution, cleanup boundary, deletion ownership

| Field | Value |
|---|---|
| Purpose | Close three remaining inconsistencies in the deletion contract |
| Status | **Contract text. Nothing adopted. No implementation authorized** |
| Relationship | **Additive.** None of `TASK-03-CONTRACT-ACCEPTANCE-2026-11-22.md`, `-CLARIFY-1-`, or `-CLARIFY-1A-` is edited |
| Supersedes | L-03 (partial rows), L-14 (unreadable-entry row), L-21 (scope narrowed), K-10 (partial row), and adds ownership rules the earlier text did not state |
| Recovery token | `TASK-03-CONTRACT-CLARIFY-1B-2026-11-22` |

Prefixes: `M-` normative, `MT-` acceptance test.

---

## 1. Partial deletion cannot claim operation attribution

**Finding, accepted — and this contradiction is mine.** L-03 assigns `partial, attribution: operation` to a partly emptied directory, while L-21 requires either an in-process successful removal or a `removed` marker written after verifying absence. A partly emptied directory is **not** absent, and a process that died mid-removal produced no successful result. Both L-03 rows are unsatisfiable under L-21.

The deeper error: observing a directory before starting proves it existed, not that every subsequent change was ours. Attribution is a claim about **authorship**, and a partial state is only a claim about **effect**.

**M-01 — Replaces K-10's and L-03's partial rows.**

| Entry state | Directory | Terminal state | Attribution |
|---|---|---|---|
| `pending` | present, intact | `aborted` | — (no removal claim) |
| `pending` / `removing` | present, partly emptied | `partial` | **`unknown`** |
| `removing` | absent | `complete` | `unknown` |
| `removed` | absent | `complete` | `operation` |

**M-02 — Narrowed `attribution: operation`.** It is permitted **only** for `state: complete`, and only with observed presence before removal **and** either a successful in-process removal result or a durable `removed` marker. No other state may carry it.

**M-03 — What a `partial` receipt records instead.** Facts about the operation, not claims about the world:

| Field | Meaning |
|---|---|
| `state: partial` | the directory was observed present and is now partly emptied |
| `effect: present_partial` | a state observation, not an attribution |
| `removal_result` | `success` \| `failure` \| `no_result` — the outcome of **our** call, which is a fact |
| `removed_record_count` | observed remaining/removed counts, whichever the schema names |
| `attribution: unknown` | no authorship claim |

This keeps the information a partial deletion exists to convey — something was removed and the operation did not finish — without asserting who removed it.

**M-04 — Why no exception is granted for in-process failure.** A removal that returns failure in-process establishes that *our* call ran and failed. It does not establish that the observed partial emptiness resulted from it: between our presence observation and our failure return, another actor could have removed files. "Our call failed" and "we removed these files" are different propositions, and only the first is evidenced.

**Alternatives considered.**
- *`attribution: operation` for in-process partial with a documented condition.* Rejected: the condition would have to exclude external interference, which no filesystem evidence can do without an exclusive lock on the directory — and the directory is precisely what is being destroyed.
- *A new attribution value such as `probable`.* Rejected: it names a confidence rather than an actor, and invites consumers to treat it as a claim.
- *Drop `partial` and report only `indeterminate`.* Rejected: the effect is genuinely observed and is worth recording; only the authorship is unknown.

**MT-01** — Every partial case yields `attribution: unknown`, including in-process failure, crash, and external interference.
**MT-02** — `attribution: operation` appears only alongside `state: complete` (table-driven, exhaustive over the state/attribution product).
**MT-03** — A partial receipt records `removal_result` truthfully and never states that this installation removed files.
**MT-04** — An `aborted` record carries no attribution value at all.

---

## 2. Cleanup must distinguish an invalid entry from an unreadable journal

**Finding, accepted.** L-14 permits removing an unreadable entry; L-18 says an unparseable journal is corrupt for **every** pending intent; L-04 routes corruption to `indeterminate`. Removing a stale entry requires rewriting the journal — and a rewrite of a file we cannot fully parse can silently destroy other pending operations.

Two classes, with different powers:

| Class | Definition | Cleanup permitted? |
|---|---|---|
| **Entry-level invalid** | The journal parses as a whole; this entry fails validation (missing field, bad type, malformed value) but its `operation_id` is readable | **Yes, conditionally** (M-06) |
| **Journal-level unreadable** | The file does not parse; is truncated mid-record; carries an unsupported `journal_v`; or cannot be read at all | **No** (M-07) |

**M-05 — Cleanup never changes an outcome.** A terminal receipt is the outcome of a deletion. It does not depend on the journal, so journal corruption cannot alter it. The receipt is returned unchanged in every cleanup case.

**M-06 — Conditional entry-level cleanup.** Removal of a stale entry is permitted only when **all** hold:

1. the journal parses as a whole;
2. the entry is located by exact `operation_id` match, not by position or inference;
3. every other entry validates and re-emits byte-equivalently under the rewrite;
4. the rewrite is atomic (temp, flush, rename) under the journal lock.

If any condition fails, the entry is preserved and the call reports `journal_cleanup_refused`.

**M-07 — Journal-level unreadable: preserve and report.**

> If the journal cannot be parsed, fully validated and safely rewritten, reconciliation **does not write it at all**. The file is preserved byte-identically, the existing terminal receipt is returned, and the call reports `journal_cleanup_unavailable`. No automatic repair, no partial rewrite, no truncation to a parseable prefix.

**M-08 — No silent discard.** No code path may drop an entry it cannot re-emit byte-equivalently. Preservation is always preferred to cleanliness; a stale entry that outlives its terminal receipt is a reported condition, not a fault to be cleared by force.

**M-09 — Consequence, stated.** A permanently corrupt journal leaves a stale entry that no reconciliation can remove. Diagnostics report it once per call and stop there. Recovery requires an out-of-band decision by the owner, which this contract does not automate.

**Alternatives considered.**
- *Rewrite a truncated journal from its parseable prefix.* Rejected: the prefix is exactly the part that may already be missing entries; writing it back makes the loss permanent and silent.
- *Delete the whole journal when it is corrupt.* Rejected: it would discard other pending operations' recovery information to tidy one entry.
- *Quarantine by renaming the journal.* Rejected for Phase 1: a rename is a write to a file we do not understand, and it creates a second file whose lifetime nobody governs.

**MT-05** — Entry-level invalid + terminal → the entry is removed, every other entry byte-identical, receipt returned.
**MT-06** — Journal-level corrupt + terminal → the file is byte-identical afterwards, the receipt is returned, and the result is `journal_cleanup_unavailable`.
**MT-07** — A rewrite is refused when any other entry fails validation or carries an unknown version.
**MT-08** — No path deletes or truncates the journal file.

---

## 3. Concurrent deletions need an owner per directory

**Finding, accepted.** The journal lock serializes individual read-modify-write operations; it does not serialize a deletion's *lifecycle*. Two operations can both observe presence, both create intents, and then both race to remove — producing two receipts for one physical deletion and an unresolvable attribution question.

**M-10 — Ownership record.** A non-terminal journal entry (`prepared`, `pending`, `removing`, `removed`) keyed by `directory_path` **is** the ownership record for that directory. At most one such entry may exist per path.

**M-11 — Acquisition is atomic with the presence observation.**

```
under journal lock:
  existing non-terminal entry for path?  → release; return busy (§M-12)
  terminal receipt for existing entry?   → release; return stale_entry_reconcile_required (§M-13)
  stat directory:
      absent  → release; write terminal `absent_prior`; no entry, no intent
      present → insert entry{prepared, observed_present_at_start: true}; release
→ intent → journal{pending} → remove → journal{removing} → verify → journal{removed} → terminal → cleanup
```

Inserting the entry while still holding the lock, after the presence check, is what makes ownership precede any removal. No process can reach `remove` without an entry it owns.

**M-12 — A second request never claims the removal.**

> While another operation owns the directory, a second deletion request performs **no** presence observation for deletion purposes, writes **no** intent, removes **nothing**, and writes **no** receipt. It returns `busy_pending_deletion` with the owner's `operation_id`.

The second caller learns the outcome from the owner's receipt when one exists. Two receipts for one physical deletion are a contract violation.

**M-13 — Stale entry at acquisition.** If an entry exists whose terminal receipt is already written, the second request returns `stale_entry_reconcile_required` and does nothing else. Cleanup remains reconciliation's exclusive job (M-06/M-07); acquisition is not a cleanup path. A new attempt may proceed only after that cleanup succeeds.

**M-14 — Scope.** Ownership is per directory: deletions of different directories proceed concurrently. Cross-process serialization is the journal lock; within a process the same rule applies, so two goroutines cannot both own one path.

**M-15 — Release.** Ownership ends when the entry is removed, which happens only after the terminal receipt is durable. A crashed owner keeps ownership until an explicit reconciliation ends it — which is the correct conservative behaviour and is why startup reports pending owners.

**Alternatives considered.**
- *Per-directory lock file inside the conversation directory.* Rejected: the directory is the thing being removed, so the lock would vanish mid-operation.
- *A separate lock directory keyed by a path hash.* Rejected: a second locator store with its own lifetime and privacy surface, when the journal already holds the mapping.
- *Let the second operation proceed and resolve the duplicate at reconciliation.* Rejected: the duplicate receipts are already written by then, and one of them would have to be retracted — which the receipt vocabulary does not permit.
- *Return an error and let the caller retry.* Rejected as the only rule: `busy_pending_deletion` carries the owner's identifier, which is more useful than a bare failure and does not require polling.

**MT-09** — A and B race: exactly one intent, exactly one terminal receipt, B receives `busy_pending_deletion` with A's `operation_id`, and B performs no filesystem change.
**MT-10** — Only the owner's receipt may carry `attribution: operation`; a second `complete` for the same path fails the test.
**MT-11** — Deletions of two different directories both complete concurrently.
**MT-12** — Two processes, same path: one acquires, one receives busy; the journal is not corrupted.
**MT-13** — A stale entry makes a new request return `stale_entry_reconcile_required` and change nothing; after successful reconciliation the request proceeds.
**MT-14** — A crashed owner blocks new deletions of that path until reconciliation, and startup reports it.

---

## 4. Delta consistency matrix

| Clause | 1b disposition | Coherent with |
|---|---|---|
| K-10 partial row | **superseded** by M-01 | L-21/M-02 |
| L-03 partial rows | **superseded** by M-01 | M-02, M-03 |
| L-14 unreadable-entry row | **superseded** by M-06/M-07 | L-18, L-04 |
| **L-21** attribution | **narrowed** to `complete` only | M-02 |
| L-01 write order | refined by M-11: acquisition + observation are one locked section; the inversion is preserved | L-03 |
| L-02 entry schema | unchanged | M-10 |
| L-04 degradation | unchanged; now explicitly scoped to *outcome* conclusions, not cleanup | M-05 |
| L-06 two readers | unchanged | M-12 |
| L-13 terminal-first | unchanged | M-05, M-13 |
| L-15 no directory contact | reinforced by M-12, M-13 | C-36 |
| L-18 integrity | unchanged; M-06 adds rewrite preconditions | M-07, M-08 |
| L-19 retention | unchanged | M-09 |
| K-09 `absent_prior` | unchanged; now reached inside the acquisition lock | M-11 |
| C-34/C-35/C-36 | unchanged | M-01, M-05 |
| C-37 correlation | unchanged | — |
| C-31 byte preservation | untouched | ✓ |

**Event-kind count: unchanged — five ledger kinds plus two evidence kinds.** Nothing here adds a kind; `busy_pending_deletion`, `stale_entry_reconcile_required` and the `journal_cleanup_*` results are **return values**, not receipt states or events.

### Residual items

1. `selection_resolved` and D-3's ancestor rule remain deferred.
2. A permanently corrupt journal leaves an irremovable stale entry (M-09) — reported, not repaired.
3. Partial deletions carry no authorship claim, so a user-facing summary can say "deletion did not complete" but never "some files were deleted by the app".
4. Power-loss durability `UN`; workload `Unestablished`.

### Authorization

No implementation, migration, ledger creation, persistence adoption, or FC-1 approval. The contract can be internally consistent without being implemented or empirically validated; no `T-`, `KT-`, `LT-` or `MT-` test has been executed.
