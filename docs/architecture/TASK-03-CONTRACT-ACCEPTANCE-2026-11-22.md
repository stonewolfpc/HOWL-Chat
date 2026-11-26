# TASK-03 — conversation authority contract, revision 1.0-rc

| Field | Value |
|---|---|
| Purpose | Single acceptance surface for FC-1: numbered clauses, exact language, testable criteria |
| Status | **Contract text for owner acceptance. Nothing adopted. No implementation authorized** |
| Supersedes | The **contract text** of rows 524–527 where restated here. Their documents remain as historical records (§1.2) |
| Authority | FC-1 `PROPOSED`; ADR-001/002 accepted. This document approves no clause |
| Recovery token | `TASK-03-CONTRACT-ACCEPTANCE-2026-11-22` |

Clause prefix `C-` for normative statements, `T-` for acceptance tests, `D-` (reserved; not used here).

---

## 1. Scope and authority

**C-01 — Authority.** This document is contract text, not an approval. FC-1 remains `PROPOSED` until the owner explicitly accepts these clauses. Silence, implementation, or a passing test is not acceptance.

**C-02 — Scope.** One conversation owns one ledger and, optionally, one recovery file. Across the installation there is one deletion-evidence file. Nothing else is in scope.

**C-03 — Non-goals for Phase 1.** No snapshots, no branching, no merge, no memories or media in the ledger, no persisted projections, no UI, no background maintenance, no compaction, no migration of existing data.

**C-04 — Supersession.** Where this document restates a clause from rows 524–527, this text governs. The earlier documents are not edited; later contradictions are resolved here.

**T-01** — The four earlier documents remain byte-identical except where a supersession note was appended, and every appended note names the clause that governs.

---

## 2. Storage substrate

**C-05 — File.** `history.ndjson` in the conversation directory, one JSON object per line, UTF-8, `\n`-terminated. Commits replace the whole document through DurableStore. A true append journal is **not** required on current evidence.

**C-06 — Logical append-only.** Physical replacement preserves logical append-only history: records at sequences `1..n` are carried as **bytes**, never re-serialised, never reordered, never rewritten.

**C-07 — Version.** Each record carries `v`. `v = 1` in Phase 1. A record with a higher `v` causes refusal, not best-effort parsing.

**C-08 — Lock.** One `ledger.lock` per conversation. Sequence is assigned under the lock; a second writer waits or reports `Busy`. No interleaving.

**T-02** — After committing sequence `n`, the first `n-1` lines are byte-identical to their previous bytes, and the line count is exactly `n`.
**T-03** — Two concurrent writers produce no duplicate and no skipped sequence; the loser reports `Busy` or succeeds with a distinct sequence.

---

## 3. Envelope `FC1-08`

**C-09 — Required fields (nine).**

| Field | Type | Invariant |
|---|---|---|
| `v` | int | `1`; higher → refuse |
| `event_id` | string | UUIDv4; unique per conversation |
| `conversation_id` | string | equals the containing directory id |
| `sequence` | int | ≥ 1; exactly previous + 1 |
| `event_type` | string | one of the six ledger kinds |
| `timestamp` | string | RFC3339 UTC with nanoseconds |
| `idempotency_key` | string | unique per conversation; replay rules C-15 |
| `actor` | string | `user` \| `assistant` \| `system` \| `recovery` |
| `durability` | string | `observed` in Phase 1 (C-16) |

**C-10 — Conditional fields (three).** `turn_id` — required on turn events, forbidden elsewhere, immutable across corrections and retraction. `predecessor_event_id` — required for every `sequence ≥ 2`, forbidden at `sequence 1` (C-17). `payload` — required, object, kind-specific (§4).

**C-11 — Deliberately absent.** `state_id` (derived, C-22), `message_id` (equals `turn_id`), `affected_ids`, `operation`, `provenance` (reduced to `actor`), and any content hash or title.

**C-12 — `event_id` uniqueness.** Scope: per conversation. Enforced at commit against the records already being read, under the lock; a collision is refused and a new id generated. Never derived from sequence, position, timestamp, content, filename, or display order.

**C-13 — Sequence allocation and failed commits.**

> Sequence is **derived at commit time from the committed file, under the lock**. It is never pre-allocated, never reserved, and never consumed by a failed attempt. Any failure before the successful rename leaves the live file unchanged, so the next attempt reads the same last sequence and allocates the same number.

Pre-reserved ranges, in-memory counters, and "burn the number on failure" are contract violations.

**C-14 — Read-side gaps.** A gap between committed records marks the ledger `unresolved` from the gap. Sequences are never renumbered; a gap is never filled by inference.

**C-15 — Idempotency.** Scope per conversation. Same key **and** identical `event_type`, `actor`, `turn_id`, `payload` → return the existing event; no new record, no sequence consumed. Any difference → refuse. Comparison is canonical (sorted keys, no insignificant whitespace, stated numeric form) and is used **only to compare** — committed bytes are never canonicalised into storage. An unknown field in a replay payload makes it a conflict.

**C-16 — `durability` values.** `requested` — the caller asked. `attempted` — a temp document was written, replacement unconfirmed. `observed` — **the writer observed a successful `rename` onto the live path after a flushed temp write, under the lock.** `verified` — a later read confirmed the document by hash; **not written in Phase 1**. `unknown` — legacy or recovered records. See §8 for what `observed` does not prove.

**C-17 — `predecessor_event_id`.** Required for `sequence ≥ 2`; absent only at `sequence 1`; must equal the `event_id` at `sequence − 1`. It records **event ordering only** — not semantic state ancestry, not visibility, not a state change, and not a branch.

**T-04** — A commit whose rename is forced to fail leaves the file byte-identical and the next successful commit uses the same sequence.
**T-05** — An idempotent replay produces one record and one sequence; a conflicting replay produces an error and no record.
**T-06** — A record with `v = 2` is refused with an explicit error and nothing is rewritten.

---

## 4. The six ledger kinds

**C-18 — Kind set.** `conversation_created`, `turn_committed`, `turn_corrected`, `turn_retracted`, `selection_changed`, `selection_resolved`. An unknown `event_type` marks the ledger `unresolved` from that record; it is never skipped or guessed.

**C-19 — Payload schemas.**

| Kind | Required payload | Optional | Target rules | Projection effect |
|---|---|---|---|---|
| `conversation_created` | `{}` | `name` (string ≤ 200) | — | Creates empty conversation; selection `empty`; **no state** |
| `turn_committed` | `role` (`user`\|`assistant`), `content` (string) | `character_id` (string), `has_image` (bool) | — | Appends turn; **new state**; selection follows |
| `turn_corrected` | `target_turn_id` (string), `replacement_content` (string) | `reason` (string) | Live or already-corrected turn in this ledger | Replaces content **in the projection only**; **new state** |
| `turn_retracted` | `target_turn_id` (string) | `reason` (string) | Live or corrected turn; not already retracted | Excludes turn from projection; **new state** |
| `selection_changed` | `state_id` (string) | — | Existing, non-retracted state | Moves pointer; **no state** |
| `selection_resolved` | `deleted_state_id` (string), `rule` (`ancestor`\|`empty`) | `chosen_state_id`, required iff `rule = ancestor`, forbidden if `empty` | State being deleted; chosen must be a non-retracted ancestor | Moves pointer; **no state** |

**C-20 — Invalid events.** A record violating any invariant is **refused at commit** and **`unresolved` at read**. It is never dropped silently, never repaired, never reordered, and never partially applied.

**T-07** — Each kind's invalid cases in the table (duplicate `turn_id`, unknown target, retracted target, dangling `state_id`, `rule`/`chosen_state_id` mismatch, second `conversation_created`, `role`/`actor` mismatch) are refused with distinct errors and leave the file unchanged.

---

## 5. Projection determinism and historical states

**C-21 — Projection is a pure function of the event sequence.**

> For a ledger with committed records `1..k`, the projection `P_k = fold(apply, events[1..k], empty)` is deterministic: it is a pure function of the record bytes in sequence order. It may not depend on wall-clock time, locale, filesystem iteration order, map iteration order, randomness, process state, or the order in which events were *received*.

Serialisation is canonical **for comparison only**; projections are not stored in Phase 1.

**C-22 — State identity and addressability.** A state exists only where a state-producing event exists. `state_id` is that event's `event_id`, derived and never stored. Each state has exactly one parent: the previous state in sequence order. Every produced state remains **addressable and selectable forever**, including after later corrections and retractions.

**C-23 — Corrections target `turn_id`, never a state.**

> `turn_corrected` and `turn_retracted` name a stable `turn_id`. They do not name a `state_id`, a sequence, a byte range, or a line.

**C-24 — Why an earlier state cannot be rewritten.** A correction committed at sequence `m` is not among `events[1..k]` for any `k < m`. By C-21, `P_k` is unchanged by it — not by convention, but because it is not an input. Earlier states are therefore immutable **by construction**, not by policy.

**C-25 — No historical-state editing.** Phase 1 has no clause that edits a state, rewrites a record in place, or alters the past. Independent historical-state editing requires separate authorization.

**C-26 — Selection is derived.** `selection(k)` = the last `selection_changed` or `selection_resolved` at sequence ≤ k; if none, the newest state at k; if no states, `empty`.

**T-08** — Replaying identical ledger bytes twice yields byte-identical projections.
**T-09** — Committing a correction at `m` leaves `P_{m-1}` byte-identical to its value before the commit.
**T-10** — A state produced before any correction remains selectable afterwards and projects the pre-correction content.
**T-11** — No code path accepts a `state_id` as a correction target.

---

## 6. Correction and retraction semantics

**C-27 — Targeting.**

| Operation | May target | Must not target |
|---|---|---|
| `turn_corrected` | live turn; already-corrected turn | retracted turn; non-existent turn; another correction |
| `turn_retracted` | live turn; corrected turn | already-retracted turn; non-existent turn; a selection event |

**C-28 — Correction chains.** Corrections may be corrected indefinitely. For `P_k`, corrections of a turn are applied in sequence order up to `k`; the last one wins there.

**C-29 — Retraction is irreversible.** There is no unretract event. Restoring content requires a **new `turn_id`** — a new fact, not an undo. The original bytes remain in the ledger; the identity does not return.

**C-30 — Repeated and conflicting operations.** Identical replay → one durable effect (C-15). Same key, different content → refused. Correction followed by retraction → both retained; the turn is excluded; not an error.

**C-31 — Byte preservation of targets.** The target record is never opened for writing. After any correction or retraction its line is **byte-identical**, and the number of lines is unchanged. No in-place mutation, no re-serialisation, no compaction.

**T-12** — After a correction and a retraction, the target line's hash equals its pre-operation hash and the ledger length is unchanged.
**T-13** — A retracted turn cannot be corrected, and an attempt leaves the file unchanged.
**T-14** — Retracting, then re-adding the same content, produces a different `turn_id` and two distinct turns in the ledger.

---

## 7. Deletion, interrupted deletion, reconciliation

**C-32 — Option A.** Deleting a conversation removes its directory: ledger, snapshots, and `recovery.ndjson`. The deletion receipt survives. **No retention window and no background sweep.** Authored content outside the boundary and both `memories.json` locations are untouched.

**C-33 — Receipt contents.** `operation_id`, `timestamp`, `authority`, `record_count`, `snapshot_count`, `state`, plus `removed_record_count` and `attribution` where applicable. Excluded: conversation ID, salted conversation-ID hash, last-record hash, title, UI counts or labels.

**C-34 — Lifecycle.**

```
write {state: intent}  →  remove the directory  →  write terminal record
```

A record written before removal is **intent**. **It proves a deletion was requested and started. It is never proof that deletion completed.** Completion requires a terminal record. The intent record is never mutated.

**C-35 — The crash window.** Window `W` is: the directory removal returned success, and the terminal record has not been flushed.

> During `W`, the only durable evidence is the `intent` record plus the filesystem. An `intent` with no terminal record is reported as **`unknown outcome`** — never as complete, never as failed, and never swept.

Terminal states are `complete`, `partial`, `aborted`, `indeterminate`. `complete` means the effect is established (receipt present, directory absent). `attribution` records `operation`, `prior_absence`, or `unknown`; the effect is never stated more strongly than the evidence, and after a crash in `W` the attribution is `unknown`.

**C-36 — Reconciliation is explicit and bounded.**

> `ReconcileDeletion(operation_id)` is an explicit operation over records the caller names. It is not a sweep, not a timer, not automatic cleanup, and it never runs on a read path. It deletes nothing, creates no identity, and infers nothing beyond the table below.

| Receipt | Directory now | Permitted terminal record |
|---|---|---|
| intent | absent | `complete`, `attribution: unknown` — effect established, actor not established |
| intent | present, intact | `aborted`, `removed_record_count: 0` |
| intent | present, partially emptied | `partial` with the observed remaining counts |
| terminal | any | none written; the existing record is returned (idempotent) |
| none | any | refused — a receipt cannot be invented for an operation that has none |

Startup may **report** unresolved intents read-only. Acting on them requires the explicit call; automatic invocation is outside this contract.

**C-37 — Correlation risk, stated.** `record_count` and `snapshot_count` can be matched against a **surviving copy** of a ledger, and `timestamp` narrows the candidates. That is a weak fingerprint, not content. Bucketing counts or coarsening the timestamp would reduce it and also reduce evidence value; neither is adopted by default. The evidence file is app-state, not exposed through the UI or any API, not exported in diagnostics, and retained without expiry.

**T-15** — A crash simulated between removal and terminal write yields `unknown outcome` on every read path, never `complete`.
**T-16** — `ReconcileDeletion` on that record writes exactly one `complete` with `attribution: unknown`; a second call writes nothing and returns the same record.
**T-17** — `ReconcileDeletion` with no matching `intent` is refused and writes nothing.
**T-18** — No deletion path removes a directory without first flushing an `intent` record.
**T-19** — The receipt contains none of the five excluded items (asserted on the serialised record).

---

## 8. Durability terminology

**C-38 — What `observed` proves.** That the writer's `rename` onto the live path succeeded after a flushed temp write, under the lock, in this process.

**C-39 — What `observed` does not prove.** It does **not** establish power-loss durability; **not** post-restart presence; **not** an independent read-back or hash verification; **not** directory-metadata flushing; **not** that any reader has agreed; **not** that a replica exists.

**C-40 — Power-loss durability remains `UN`.** No clause, value, log line, or API response may present `observed` as durability. `verified` is reserved and unwritten in Phase 1.

**T-20** — No code path writes `verified`.
**T-21** — Mapping `observed` to user-facing "saved"/"safe"/"durable" text is a test failure, not a style issue.

---

## 9. Evidence store authority and lifecycle

**C-41 — Not a second authority.** `recovery_recorded` is **diagnostic evidence about how a read went**, not durable conversation state, not a replica, not a fallback authority, and not reconstructable into a conversation. Nothing may be rebuilt from it, and no conversation may be served from it.

**C-42 — Two homes.**

| Record | Home | Lifetime |
|---|---|---|
| `conversation_deleted` | app-state `evidence/deletions.ndjson` | Outlives the conversation |
| `recovery_recorded` | the conversation's `recovery.ndjson` | Dies with the conversation directory under Option A |

**C-43 — Expected disappearance.** Losing `recovery.ndjson` when a conversation is deleted is **correct behaviour**, not data loss. Diagnostics for that conversation end with it. No replacement file is written elsewhere.

**C-44 — `recovery_recorded` schema.** `observation_id` (UUIDv4), `timestamp`, `scope` (`ledger`\|`manifest`\|`projection`), `boundary` (int), `records_seen` (int), `ledger_hash_before`, `ledger_hash_after` (must be equal — the read changed nothing), `outcome` (`prefix_valid`\|`refused`\|`recovered_generation`). **No conversation identifier field**; the record is identified by its location.

**C-45 — Written after the read, by an explicit call.** The reader performs no write, so the read-only guarantee holds literally.

**C-46 — Evidence is not the ledger.** The app-state deletion receipt is evidence **of the deletion operation**. It does not replace the ledger, does not represent conversation content, and its presence never implies the conversation still exists — nor does its absence imply deletion succeeded.

**T-22** — A read of an intact ledger and of a truncated ledger produces `ledger_hash_before == ledger_hash_after` in every recorded observation.
**T-23** — Deleting a conversation removes `recovery.ndjson` and leaves the app-state receipt; no recovery record is written outside the conversation directory.
**T-24** — No module exposes a code path that reads `recovery.ndjson` as conversation state.

---

## 10. Legacy, exclusions, limits

**C-47 — Legacy is read-only.** Records lacking the envelope, positional identifiers (`fmt.Sprintf("%d", index)`), and multi-line `history.json` are **never promoted**, never assigned sequences, and reported `unresolved`.

**C-48 — Linear only.** No branch, no multi-parent state, no merge, no cross-conversation shared state. Fork-as-new-conversation is **unadopted and outside Phase 1**; its fields (`origin_conversation_id`, `origin_state_id`) are **prohibited in Phase 1**, so adopting it later changes no required field.

**C-49 — Turn identity.** `turn_id` immutable; `message_id` equals `turn_id`; one message per turn.

**C-50 — Accepted limits.** Commit cost grows with record count; no compaction; retraction irreversible; `observed` is not durability; workload is `Unestablished`.

---

## 11. Acceptance

**C-51 — Acceptance tests are the criteria.** A clause is satisfied when its `T-` tests pass against the implementation. No clause is satisfied by documentation alone.

**C-52 — Owner decisions.** Nine remain: the nine-field envelope with required predecessor · six ledger plus two evidence kinds · linear not graph · legacy read-only · branching excluded and fork unadopted · one message per turn · receipt without the five excluded items · whether per-conversation deletion proof is required.

**C-53 — Authorization unchanged.** No implementation, no call-site migration, no ledger creation, no persistence adoption, no FC-1 approval by implication. Acceptance of this contract is a separate, explicit act.
