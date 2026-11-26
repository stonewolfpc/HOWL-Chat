# TASK-03 — clarification 1c: acquisition reader, abandoned entries, directory identity

| Field | Value |
|---|---|
| Purpose | Close the final three deletion-authority gaps |
| Status | **Contract text. Nothing adopted. No implementation authorized** |
| Relationship | **Additive.** None of rows 528–531 is edited |
| Supersedes | L-06 (reader count), L-03 (row 1 and the `pending`+absent row), M-14 (exclusive release rule); adds the identity boundary |
| Recovery token | `TASK-03-CONTRACT-CLARIFY-1C-2026-11-22` |

Prefixes: `N-` normative, `NT-` acceptance test.

---

## 1. The acquisition path is a third reader — and it must be

**Finding, accepted.** L-06 says "exactly two readers"; M-11 requires every deletion request to inspect the journal under lock. That is a third reader, and without it M-10 and M-12 cannot be enforced at all. The count was wrong, not the mechanism.

**N-01 — Replaces L-06. Exactly three readers.**

| Reader | Permission | Prohibited |
|---|---|---|
| **Startup discovery** | Read-only. Reports pending `operation_id`s | Reconciling, writing, removing, exposing paths |
| **Deletion acquisition** | Under the journal lock: inspect ownership for the candidate directory; insert exactly one `prepared` entry; reject with a status. May write one terminal `absent_prior` receipt when the directory is absent | Reconciling, cleaning stale entries, removing entries, revealing any other entry's contents, touching any conversation directory |
| **Explicit reconciliation** | Read-only inspection, then the authorised writes of C-36/M-06/M-07/M-19 | Any other write |

**N-02 — What acquisition may return.** Only a status and, when ownership is held elsewhere, the owning `operation_id`. It never returns another entry's path, counts, timestamps or state beyond the fact that an entry exists and who owns it. `busy_pending_deletion` conveys exactly that.

**N-03 — Acquisition writes are bounded to two.** One `prepared` insertion, or one `absent_prior` receipt. It never writes a terminal record for an operation it did not perform, and never removes an entry — including one it just inserted, because a path that inserted `prepared` then failed proceeds to the intent, where reconciliation's authority takes over.

**N-04 — Reporting abandoned entries.** Startup may additionally report entries in `prepared` with no corresponding intent as abandoned candidates. Reporting is read-only and names `operation_id`s only.

**Alternatives considered.**
- *Keep two readers and let acquisition ask reconciliation.* Rejected: it puts a write-capable authority on the hot path of every deletion and serialises all deletions behind one reconciler.
- *Enumerate readers by behaviour rather than by name.* Rejected: the earlier wording was too coarse precisely because it counted readers instead of stating scopes; the table states scopes.

**NT-01** — A deletion request on an unowned path performs one `prepared` insertion and nothing else.
**NT-02** — A deletion request on an owned path returns `busy_pending_deletion` with the owner id and touches nothing.
**NT-03** — Acquisition never removes an entry and never writes a terminal receipt for an operation it did not perform.
**NT-04** — Acquisition's response contains no path, count or timestamp (asserted on the serialised status).

---

## 2. `abandoned_no_intent` had no release path — and the state machine closes it

**Finding, accepted.** M-10 makes a `prepared` entry an ownership record; L-03 refuses to write a receipt for it; M-14 releases ownership only after a terminal receipt. A crash between `prepared` and the intent therefore blocks that directory **forever**. That is a deadlock, not a conservative choice.

**N-05 — The entry state is a progress barrier.** Because each transition is durable and ordered, the state proves what could not yet have happened:

```
insert {prepared} → intent → {pending} → {removing} → remove call → verify → {removed} → terminal
```

| State | Removal call | Inference |
|---|---|---|
| `prepared` | **cannot have been made** | Any absence is external |
| `pending` | **cannot have been made** | Any absence is external |
| `removing` | may have been made, result unknown | Effect attributable only if verified |
| `removed` | made and verified absent | `attribution: operation` permitted |

**N-06 — Replaces L-03's row 1. The abandon action.**

> Reconciliation may remove a `prepared` entry **without writing any receipt**, and only when, under the journal lock, it verifies all three: the entry is located by exact `operation_id`; its state is exactly `prepared`; no intent record and no terminal receipt exist for that `operation_id`. Result: `abandoned_no_intent`. It touches no conversation directory and manufactures no deletion record.

This is the single case in which reconciliation removes an entry with no receipt. It exists because no deletion was ever requested — the intent *is* the request — so there is nothing to record and nothing to claim.

**N-07 — Replaces the `pending` + absent row.** A `pending` entry with an absent directory and a present intent yields **`absent_prior`, `attribution: prior_absence`**, not `complete`/`unknown`. The barrier proves the removal call was never made, so the absence cannot be ours. This is a correction of L-03 in the direction of claiming less.

**N-08 — `prepared` + intent present.** Yields **`aborted`** (removal never called, nothing removed) and a terminal receipt, then normal cleanup. A crash between the intent and the `pending` update therefore resolves without ambiguity.

**N-09 — The `receipt_missing` case.** A `pending`/`removing`/`removed` entry with no intent is not abandonment: the intent precedes those states, so its absence means the receipt store lost a record. Result: `indeterminate`, `receipt_missing`; entry preserved; nothing claimed.

**N-10 — Supersedes M-14's exclusivity.** Ownership ends in exactly two ways: removal of the entry after a terminal receipt, or the abandon action of N-06. No other path releases it.

**N-11 — Complete inference table.** This replaces L-03 wholesale; rows marked † are corrections.

| Entry | Intent | Terminal | Directory | Outcome | Attribution |
|---|---|---|---|---|---|
| `prepared` | absent | absent | any | `abandoned_no_intent` — entry removed, **no receipt** | — |
| `prepared` | present | absent | any | `aborted` (removal never called) | — |
| `pending` | present | absent | present | `aborted` | — |
| `pending` | present | absent | absent | **`absent_prior`** † | `prior_absence` |
| `pending`/`removing`/`removed` | absent | absent | any | `indeterminate`, `receipt_missing`; entry preserved | `unknown` |
| `removing` | present | absent | absent | `complete` | `unknown` |
| `removing` | present | absent | intact | `aborted` | — |
| `removing` | present | absent | partly emptied | `partial` | `unknown` |
| `removed` | present | absent | absent | `complete` | `operation` |
| `removed` | present | absent | present | `indeterminate`, `state_inconsistent`; entry preserved | `unknown` |
| any | any | present | any | return receipt; cleanup per M-06/M-07 | — |

**Alternatives considered.**
- *Write a `aborted` receipt for an abandoned entry.* Rejected: no deletion was requested, so a receipt would describe an operation that never existed — exactly what C-36's "a receipt cannot be invented" forbids.
- *Release ownership by age (a timeout on `prepared`).* Rejected: a timer is a sweep, and it would silently release a directory whose owner is merely slow.
- *Treat `prepared` + absent directory as `absent_prior`.* Rejected: it would write a receipt for an operation with no intent, and N-06 already resolves this case without one.

**NT-05** — Crash between `prepared` and the intent: reconciliation returns `abandoned_no_intent`, removes only that entry, writes no receipt, and the directory is byte-identical and still deletable afterwards.
**NT-06** — The abandon action refuses when an intent exists (`aborted` path), when the state is not `prepared`, or when a terminal receipt exists.
**NT-07** — `pending` + absent directory → `absent_prior`/`prior_absence`, never `complete`.
**NT-08** — `pending` with no intent → `indeterminate`/`receipt_missing`; the entry survives.
**NT-09** — Ownership of an abandoned entry blocks new acquisitions until reconciled, then releases (both halves asserted).

---

## 3. Directory identity, not path strings

**Finding, accepted.** M-10 keys ownership by `directory_path`. Two different strings can name one directory — case variants, junctions, `..` segments, trailing separators, 8.3 short names, symlinked parents. String comparison would let two operations acquire separate entries for one physical directory, defeating M-10, M-12 and M-15 at once.

**N-12 — Ownership is keyed by canonical directory identity.**

> The ownership key is the canonical directory identity established by **PathService**, not the caller's string. Two paths that PathService resolves to the same identity are the same key. No second path-resolution implementation may exist.

PathService is the existing authority (TASK-01): it resolves final paths by handle and is the component that already knows about junctions — `EvalSymlinks` does not follow them, which is why the service exists.

**N-13 — Resolution happens inside the acquisition critical section,** after the lock is taken and before the presence check, so the identity used for the ownership comparison is the identity used for the removal.

**N-14 — Fail closed.** If identity cannot be established — resolution error, insufficient access, unsupported filesystem — the request returns **`identity_unresolved`**: no entry, no intent, no removal, no receipt. Ownership is never acquired on an unproven identity.

**N-15 — Absence is a separate outcome from failed identity.** If resolution fails because the target is absent, the request may write `absent_prior` (a statement about the requested path, not an identity claim). Any other resolution failure is `identity_unresolved`.

**N-16 — Canonical form is what is stored.** The journal holds the canonical path only; the caller's original string never enters durable state, a log, or an error message.

**N-17 — Re-check before removal.** Immediately before the removal call, the operation re-establishes identity and compares it with the identity it acquired. A mismatch returns `aborted` with `identity_changed`, and **no removal is performed**. This closes the window where a directory is removed and re-created at the same path while an operation is pending — the case where a pending deletion would otherwise destroy someone else's conversation.

**Alternatives considered.**
- *Normalise strings (lowercase, trim separators) and compare.* Rejected: it is a second resolution implementation, and it cannot see junctions or hard links at all — precisely the aliases that matter.
- *Let the filesystem decide by attempting the removal.* Rejected: it destroys the evidence of which directory was acquired and makes attribution unverifiable.
- *Acquire on the caller's string and reconcile duplicates later.* Rejected: two ownership records for one directory are the defect; reconciling them afterwards means two receipts or a retraction the vocabulary forbids.

**NT-10** — Case variant, trailing separator, `..` segment, junction and symlinked-parent forms of one directory all resolve to one key: the second request returns `busy_pending_deletion`.
**NT-11** — Unresolvable identity → `identity_unresolved`, with no entry, intent, receipt or removal.
**NT-12** — Absent target → `absent_prior`; a permission error → `identity_unresolved`. The two are distinguished.
**NT-13** — The journal contains no non-canonical path and no caller string appears in any durable record, log or error.
**NT-14** — Identity change between acquisition and removal → `aborted`/`identity_changed`, removal never called (asserted on filesystem instrumentation).

---

## 4. Delta consistency matrix

| Clause | 1c disposition | Coherent with |
|---|---|---|
| **L-06** readers | **superseded** by N-01/N-02 | M-11, L-13 |
| L-03 row 1 | **superseded** by N-06 | M-14/N-10 |
| L-03 `pending`+absent | **superseded** by N-07 | M-01, M-02 |
| L-03 `removed` row | retained; N-11 adds the inconsistent-state case | M-02 |
| L-01 write order | unchanged; N-05 states the barrier the order creates | N-06 |
| L-02 entry schema | unchanged | N-03, N-16 |
| L-04 degradation | unchanged; `receipt_missing` added as a distinct class | N-09 |
| L-13 terminal-first | unchanged | N-11 last row |
| L-14 / M-06 / M-07 cleanup | unchanged | N-06 |
| **M-10** ownership key | retained; key defined by N-12 | N-13 |
| M-11 acquisition | refined: resolution precedes the presence check | N-13 |
| M-12 second request | unchanged | N-02 |
| M-13 stale entry | unchanged | N-11 |
| **M-14** release | **superseded** by N-10 | N-06 |
| M-15 crash retention | unchanged; now releasable via N-06 | N-09 |
| C-34 / C-35 / C-36 | unchanged | N-06 |
| C-31 byte preservation | untouched | ✓ |

**Event-kind count: unchanged — five ledger kinds plus two evidence kinds.** `abandoned_no_intent`, `identity_unresolved`, `identity_changed`, `receipt_missing`, `state_inconsistent` and `journal_cleanup_*` are **return values and result codes**, not receipt states or events.

### Residual items

1. `selection_resolved` and D-3's ancestor rule remain deferred.
2. A permanently corrupt journal leaves an irremovable stale entry (reported, never forced).
3. Identity is established at acquisition and re-checked before removal; a rename or mount change between those two points is detected, not prevented.
4. Partial deletions carry no authorship claim.
5. Power-loss durability `UN`; workload `Unestablished`.

### Authorization

No implementation, migration, ledger creation, persistence adoption, or FC-1 approval. The contract is consistent as text only; no test of any prefix has been executed, because no implementation exists.
