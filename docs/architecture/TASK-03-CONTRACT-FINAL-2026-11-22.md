# TASK-03 — final contract clarifications

| Field | Value |
|---|---|
| Purpose | Close the clarification round so FC-1 can be presented for acceptance |
| Status | **Contract text for review. Nothing adopted. No implementation authorized** |
| Supersedes | `TASK-03-LEDGER-SCHEMA-PROPOSAL-2026-11-22.md` §2.2 and §4.1 **in three specific statements**, listed in §1 |
| Authority | FC-1 `PROPOSED`; ADR-001/002 accepted. This document approves no clause |
| Recovery token | `TASK-03-CONTRACT-FINAL-2026-11-22` |

## 1. What this document corrects

Preserved, not rewritten — the earlier statements stay readable beside their replacements.

| Earlier statement | Corrected here |
|---|---|
| `predecessor_event_id` is "recommended on turn events; optional otherwise" | **Required for every event with `sequence ≥ 2`**, absent only at `sequence 1`, because the contract adopts a strict linear chain (§2.5) |
| "a correction changes conversation state, so it produces a state" — left open whether *every* event does | **Only three kinds produce a state**; selection events never do (§6.2) |
| Evidence records were described as "evidence store" without an authority or lifecycle | Authority, two homes, schemas and transitions defined in §5 |

---

## 2. `FC1-08` — envelope invariants

### 2.1 `event_id` uniqueness

| Aspect | Rule |
|---|---|
| Scope | **Per conversation ledger.** Not global |
| Value | UUIDv4, lowercase, 36 characters |
| Enforcement | At commit, under the lock, against the committed records already being read for the sequence: a collision is **refused** and a fresh id generated. Never rewrites, never reuses an id |
| Prohibited derivations | sequence, array position, timestamp, content hash, filename, display order (FC1-06) |
| Cross-conversation | A collision across two conversations is not an error, but nothing may depend on global uniqueness |

### 2.2 `sequence` allocation, commit atomicity, failed commits

**The invariant that makes gaps impossible: sequence is derived at commit time, never pre-allocated.**

```
lock → read committed records → next = last_valid_sequence + 1 → build new document
     → temp write + flush → stage pre-images → manifest → rename → promote → trim → release
```

| Situation | Result |
|---|---|
| Any failure before the rename (temp write, staging, manifest, reader denial) | The live file is **unchanged**. The next attempt re-reads the same last sequence and allocates **the same number** |
| Crash before the rename | Same as above: no gap, no consumed position |
| Crash after the rename, before the caller sees it | The event **is** committed at that sequence. A retry with the same `idempotency_key` returns it; it does not allocate a second position |
| Idempotent replay | Consumes **no** sequence |
| Two writers | Serialized by the lock; the second either waits or reports `Busy`. No interleaving, so no duplicate allocation |

**Explicitly prohibited:** an in-memory counter, a pre-reserved range, or "burn the number if the commit fails". A permanent gap or a silently consumed position is a contract violation, not a cosmetic defect.

**Gap handling on read remains:** a gap between committed records is an error that marks the ledger `unresolved` from the gap. Sequences are never renumbered and a gap is never filled by inference.

### 2.3 `idempotency_key` — scope and exact comparison

| Aspect | Rule |
|---|---|
| Scope | **Per conversation.** The same key in another conversation is unrelated |
| Replay: identical | Same key **and** identical `event_type`, `actor`, `turn_id` (when present) and `payload` → return the **existing** event as the result. No new record, no new sequence, no state change |
| Conflict: any difference | Same key with a different `event_type`, `actor`, `turn_id` or `payload` → **refused with an explicit error.** Never overwritten, never "last write wins" |
| Comparison basis | Canonical comparison of those four components — object keys sorted, no insignificant whitespace, numbers in a stated canonical form. Comparison is **not** a byte comparison of the stored line, because key order and whitespace are not semantics |
| Relationship to byte preservation | The canonical form is used **only to compare**. Committed bytes are never re-serialised for storage (§4.4) |
| Unknown fields in `payload` | Included in the comparison. A replay that adds an unknown field is a **conflict**, not a match |

### 2.4 What `durability: observed` means

| Value | Precise guarantee |
|---|---|
| `requested` | The caller asked for a commit. Nothing more |
| `attempted` | A temp document was written; the replacement was not confirmed |
| **`observed`** | **The writer observed a successful `rename` onto the live path, after a flushed temp write, under the lock.** This is the Phase 1 value |
| `verified` | A subsequent read confirmed the document by hash. **Not written in Phase 1**; reserved |
| `unknown` | Legacy or recovered records whose commit was never observed |

**What `observed` does not mean:** that the bytes survived a power cut, that the filesystem flushed directory metadata, or that any reader has since agreed. **Power-loss durability remains `UN`** and `observed` must never be presented as it.

### 2.5 `predecessor_event_id`

| Aspect | Rule |
|---|---|
| Required | **Every event with `sequence ≥ 2`** |
| Absent | **Only** at `sequence 1` (`conversation_created`) |
| Value | Must equal the `event_id` of the record at `sequence − 1` in the same ledger |
| **Meaning** | **Event ordering only.** It is not semantic state ancestry, it does not imply the predecessor is visible, and it does not imply a state change |
| Branching | A single ordered pointer only. It makes no branch claim (`FC1-14b`) |

---

## 3. The six ledger payload schemas

Common to all six: `payload` is required; unknown fields are preserved on read and never rewritten; a record failing any invariant below is **refused at commit** or **`unresolved` at read** — never silently dropped, never repaired.

### 3.1 `conversation_created`

| | |
|---|---|
| Required payload | none (empty object) |
| Optional | `name` (string, ≤ 200 chars) |
| Invariants | `sequence` must be **1**; must be the first record; exactly one per ledger |
| Target rules | none |
| Invalid event | Any other kind at `sequence 1` → ledger `unresolved`. A second `conversation_created` → **refused** |
| Projection effect | Creates the conversation with zero turns. Selection = `empty`. **No state is produced** |

### 3.2 `turn_committed`

| | |
|---|---|
| Required payload | `role` (`user` \| `assistant`), `content` (string, may be empty) |
| Optional | `character_id` (string), `has_image` (bool) |
| Invariants | Envelope `turn_id` required and **not already present** in this ledger; `actor` must agree with `role` for user/assistant turns |
| Target rules | none |
| Invalid event | Duplicate `turn_id` → refused. `role`/`actor` mismatch → refused. Unknown `role` → refused |
| Projection effect | Appends the turn. **Produces a state** (`state_id` = this `event_id`). Selection moves to that state |

### 3.3 `turn_corrected`

| | |
|---|---|
| Required payload | `target_turn_id` (string), `replacement_content` (string) |
| Optional | `reason` (string) |
| Invariants | `target_turn_id` must exist as a `turn_committed`; the target must **not** be retracted; the envelope `turn_id` must equal `target_turn_id` |
| Target rules | May target a turn that has already been corrected — **yes** (latest correction by sequence wins in the projection; every correction remains) |
| Invalid event | Unknown target → refused. Retracted target → refused (§4.2) |
| Projection effect | Replaces turn content in the projection. **Produces a state.** **Original bytes untouched** |

### 3.4 `turn_retracted`

| | |
|---|---|
| Required payload | `target_turn_id` (string) |
| Optional | `reason` (string) |
| Invariants | Target must exist; target must **not already be retracted**; envelope `turn_id` = `target_turn_id` |
| Target rules | May target a corrected turn — **yes**; the correction then has no visible effect but remains in the ledger |
| Invalid event | Unknown target → refused. Double retraction of the same turn → refused (idempotent only via identical `idempotency_key`) |
| Projection effect | Turn excluded from the projection. **Produces a state.** **Original bytes untouched; the target is never deleted** |

### 3.5 `selection_changed`

| | |
|---|---|
| Required payload | `state_id` (string) |
| Optional | none |
| Invariants | `state_id` must name an existing, non-retracted state in this ledger |
| Target rules | May select any retained state, forward or backward in the chain |
| Invalid event | Dangling or retracted target → refused. Cannot be used to express `empty` (see 3.6) |
| Projection effect | Moves the selection pointer. **Produces no state** |

### 3.6 `selection_resolved`

| | |
|---|---|
| Required payload | `deleted_state_id` (string), `rule` (`ancestor` \| `empty`) |
| Conditional | `chosen_state_id` (string) — **required** when `rule=ancestor`, **prohibited** when `rule=empty` |
| Invariants | `deleted_state_id` must be the state being deleted in this operation; `chosen_state_id` must exist, be non-retracted, and be an ancestor |
| Target rules | Written only by the deletion path |
| Invalid event | `rule`/`chosen_state_id` mismatch → refused; non-ancestor choice → refused |
| Projection effect | Sets selection to the chosen state, or to `empty`. **Produces no state** |

---

## 4. Correction and retraction semantics

### 4.1 What may be targeted

| Operation | May target | May not target |
|---|---|---|
| `turn_corrected` | A live or already-corrected `turn_committed` | A retracted turn; a non-existent turn; another correction event |
| `turn_retracted` | A live or corrected `turn_committed` | A retracted turn; a non-existent turn; a selection event |

### 4.2 Are prior states addressable?

**Yes.** Every produced state keeps its `state_id` and remains selectable. A corrected turn's **earlier** states still project the original content, because each state is the turn set as of its producing event. Retraction does not remove earlier states either.

### 4.3 Is retraction reversible?

**No — not in Phase 1.** There is no "unretract" event. Restoring a message requires committing a **new turn with a new `turn_id`**, which is a new fact, not an undo. Stated as a limitation: a user who retracts by mistake cannot recover the turn's identity, though the original bytes remain in the ledger.

### 4.4 Can corrections target corrected states?

**Yes** — corrections may be corrected, chained indefinitely. The projection of a given state applies corrections **in sequence order up to that state's producing event**. There is no "correction of a correction" special case: each is a `turn_corrected` naming the same `turn_id`.

### 4.5 Repeated or conflicting operations

| Case | Behaviour |
|---|---|
| Same operation retried with the same `idempotency_key` and identical content | Returns the existing event; one durable effect |
| Same key, different content | Refused (§2.3) |
| A correction followed by a retraction of the same turn | Both retained; the projection excludes the turn. Not an error |
| Two corrections of the same turn in one state's history | Last by `sequence` wins for that state |

### 4.6 Byte preservation, restated as a prohibition

The target record is **never opened for writing**. After any correction or retraction the target's line is **byte-identical** to its previous bytes. No in-place mutation, no re-serialisation, no compaction, no "tidy-up" pass. This is checked by the D-1 prefix test and by a per-record byte comparison.

---

## 5. Evidence store — separate authority, two homes, explicit lifecycle

### 5.1 Authority

| Question | Answer |
|---|---|
| Owner | **Evidence store** — deletion receipts and recovery observations |
| Must not own | Conversation state, conversation identity, authored content |
| Writers | The deletion path (receipts) and an explicit recording call (recovery). **Never a reader, never during a read** |
| Readers | Diagnostics and audit surfaces only |
| Lifecycle | Append-only. Records are never edited or deleted individually; no background expiry or sweep |
| Failure | A missing receipt is `Unknown`, never success |

### 5.2 Two homes, because they have different lifetimes

| Record | Home | Why |
|---|---|---|
| `conversation_deleted` | **App-state** (`evidence/deletions.ndjson`) | Must survive the deletion of the conversation |
| `recovery_recorded` | **The conversation's own directory** (`recovery.ndjson`) | Identified by location rather than by content, and dies with the conversation under Option A |

### 5.3 `conversation_deleted` — schema and the intent/outcome rule

| Field | Type | Required | Note |
|---|---|---|---|
| `operation_id` | string (UUIDv4) | yes | Unrelated to the conversation |
| `timestamp` | string (RFC3339 UTC) | yes | |
| `authority` | `user` \| `system` \| `migration` | yes | |
| `record_count` | integer | yes | Count at intent time |
| `snapshot_count` | integer | yes | |
| `state` | `intent` \| `partial` \| `complete` | yes | See below |
| `removed_record_count` | integer | only when `partial`/`complete` | What was actually removed |

**Lifecycle, and the rule that matters most:**

```
write {state: intent}  →  remove the conversation directory  →  write {state: complete} or {state: partial}
```

- A record written **before** removal is **intent**. It proves a deletion was requested and started. **It is never evidence that deletion completed.**
- Completion requires a **second** record with `state: complete`. If removal fails or is interrupted, the outcome record is `partial` with `removed_record_count`.
- **No reader may report "deleted" from a lone `intent` record.** An intent without an outcome is reported as `unknown outcome`, and the next startup treats it as an unresolved lifecycle record — reported, never swept.
- The intent record is never mutated; the outcome is a new record naming the same `operation_id`.

**Superseded in part** — the `state` enumeration is extended to `intent` \| `complete` \| `partial` \| `aborted` \|
`indeterminate` by `TASK-03-CONTRACT-ACCEPTANCE-2026-11-22.md` C-35, which also defines the crash window and the
reconciliation table. `attribution` is added there. Nothing else in this section changes.

### 5.4 `recovery_recorded` — schema

| Field | Type | Required | Note |
|---|---|---|---|
| `observation_id` | string (UUIDv4) | yes | |
| `timestamp` | RFC3339 UTC | yes | |
| `scope` | `ledger` \| `manifest` \| `projection` | yes | What was being read |
| `boundary` | integer | yes | Sequence (or byte offset for `scope=manifest`) where validity ended |
| `records_seen` | integer | yes | Records successfully read before the boundary |
| `ledger_hash_before` | string (hex sha256) | yes | Proves the read modified nothing |
| `ledger_hash_after` | string (hex sha256) | yes | Must equal `ledger_hash_before` |
| `outcome` | `prefix_valid` \| `refused` \| `recovered_generation` | yes | |

**No conversation identifier field.** The record is identified by its location. Written after the read completes, by an explicit call — never by the reader itself, so A-12's "no modification during reading" holds literally.

---

## 6. `FC1-14b` / `FC1-14c` — reconciliation

### 6.1 Turn identity and state derivation

| Concept | Rule |
|---|---|
| Turn | `turn_id` (UUIDv4) assigned by `turn_committed`; **immutable** through corrections and retraction |
| Message | `message_id` **equals** `turn_id` — two FC-1 identities collapsed by a Phase 1 restriction, not by redesign |
| State | **Derived, never stored**: `state_id` = the `event_id` of the event that produced it |
| State parent | Exactly one: the previous state in `sequence` order |
| Event predecessor | `predecessor_event_id` = previous **event**, whether or not it changed state |

### 6.2 A linear event chain is not a chain of state changes

This is the clarification the review asked for:

| Kind | Produces a state? | Changes visible conversation? |
|---|---|---|
| `conversation_created` | **No** | Creates an empty conversation |
| `turn_committed` | **Yes** | Yes |
| `turn_corrected` | **Yes** | Yes |
| `turn_retracted` | **Yes** | Yes |
| `selection_changed` | **No** | Changes *which* retained state is shown, not the state |
| `selection_resolved` | **No** | Same |

So `predecessor_event_id` is **ordering only**: it links events, and several consecutive events can share one visible state. Semantic ancestry is the **state** relation, derived from the sequence of state-producing events.

### 6.3 `FC1-14b` — branching excluded, fork recorded as a possibility

- Phase 1 has **no branch event, no multi-parent state, no merge**. `predecessor_event_id` is a single ordered pointer and makes no branch claim.
- **Fork-as-new-conversation is not adopted.** If it were ever adopted, a fork would create a new `conversation_id` with a new ledger whose creation record references `origin_conversation_id` and `origin_state_id` — both **optional** fields. The fields are **not permitted in Phase 1**; their absence now means no required field changes later.
- **Explicitly outside Phase 1:** branch selection, merging, cross-conversation shared state.

---

## 7. Deletion receipt — final default

The default receipt contains **none of**: conversation ID · salted conversation-ID hash · last-record hash · title · UI counts or labels.

| Field | Status |
|---|---|
| `operation_id` | retained |
| `timestamp` | retained |
| `authority` | retained |
| `record_count` | retained |
| `snapshot_count` | retained |
| `state`, `removed_record_count` | retained (intent/partial/complete rule, §5.3) |

### 7.1 Can any retained field be correlated with surviving conversation data?

**Yes, two of them — and this must be stated rather than denied:**

- `record_count` and `snapshot_count` can be matched against the size of a **surviving copy** of the ledger (a backup, a sync target, a pre-image that escaped deletion). That is a weak fingerprint; it is not content.
- `timestamp` narrows the candidate set and, combined with the counts, makes matching easier.

**Mitigations available if the owner wants them:** bucket the counts (`<100`, `<1000`), or coarsen the timestamp to a date. Both reduce evidence value, so neither is adopted by default; they are offered, not assumed.

### 7.2 Access, retention, exposure

| Aspect | Rule |
|---|---|
| Location | App-state, never inside a conversation directory |
| Access | Not exposed through the UI or any API; no export; no inclusion in diagnostics bundles by default |
| Retention | Retained. The file is tiny and append-only. **No background sweep** (FC1-16) |
| Deletion-of-evidence | Only an explicit user action naming the set, with its own receipt |
| Per-conversation proof | **Not added.** Available only if the owner explicitly requires it, and then only as a separate decision (§8) |

---

## 8. Consistency statement, residual risks, remaining decisions

### 8.1 Is the contract internally consistent?

**Yes, within the stated Phase 1 restrictions — no contradiction was found while assembling this document.** The three earlier defects found in my own text are corrected in §1, and the two structural findings stand: `conversation_deleted` cannot live in the ledger it deletes, and `recovery_recorded` cannot be written during a read.

### 8.2 Residual risks that acceptance does not remove

1. **Power-loss durability `UN`.** `observed` is not a durability claim.
2. **Retraction is irreversible** in Phase 1; the bytes survive, the identity does not come back.
3. **Counts + timestamp are correlatable** with a surviving ledger copy (§7.1).
4. **Legacy records are outside the ledger** and never assigned identity (`FC1-07`).
5. **The workload measurement is `Unestablished`**, so the whole-document commit cost model is arithmetic, not measured.
6. **Branching, snapshots and memory copies are absent by design**, and adding them changes no required field.

### 8.3 Remaining owner decisions

| # | Decision | Status |
|---|---|---|
| 1 | Retention Option A | **Decided** |
| 2 | `FC1-08` nine-field envelope + required `predecessor_event_id` | Open |
| 3 | `FC1-11` six ledger + two evidence kinds | Open |
| 4 | `FC1-02` linear chain, not graph | Open |
| 5 | `FC1-07` legacy read-only | Open |
| 6 | `FC1-14b` branching excluded, fork unadopted | Open |
| 7 | `FC1-14c` one message per turn | Open |
| 8 | Receipt with none of the five excluded items | Open |
| 9 | Whether per-conversation deletion proof is required at all | Open |

### 8.4 Authorization

No Go implementation, no call-site migration, no ledger creation, no persistence adoption, no FC-1 approval by implication. This document is contract text for review.
