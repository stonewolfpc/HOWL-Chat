# TASK-03 — ledger schema and event definitions (proposal for review)

| Field | Value |
|---|---|
| Purpose | Present the exact `FC1-08` envelope, the exact `FC1-11` event kinds, and reconcile `FC1-14b`/`FC1-14c` |
| Status | **Proposal. Not adopted. No owner decision is resolved by this document.** |
| Supersedes | Only the *phrasing* of the earlier counts in `TASK-03-FC1-DISPOSITIONS-2026-11-22.md`; two of its summaries were incomplete and are corrected in §2.2 and §3.1 |
| Authority | FC-1 is `PROPOSED`; ADR-001/002 accepted. Nothing here approves a clause |
| Recovery token | `TASK-03-LEDGER-SCHEMA-2026-11-22` |

**Correction up front.** The earlier report said "eight core fields, plus `v` and `actor`" and "eight event kinds in the
ledger". Both counts were ambiguous or wrong. The exact envelope is **nine required fields** (§2.1), and the exact
event model is **six ledger kinds plus two evidence kinds** (§3.1) — because an event that records a deletion cannot
be stored in the ledger it deletes.

---

## 1. Container

One JSON object per line, per conversation, in that conversation's existing `history.ndjson`. Commits replace the
whole file through `DurableStore` (decision record D-1 / E2). Prior lines are carried as bytes; they are never
re-serialised.

A record that cannot be parsed does not invalidate the file: the valid prefix stands, and the remainder is
`unresolved` (A-12). **Line order is the commit order**; `sequence` states it.

---

## 2. `FC1-08` — the exact envelope

### 2.1 Required fields (nine, all events)

| # | JSON name | Type | Required | Invariant |
|---|---|---|---|---|
| 1 | `v` | integer | **always** | Format version. `1` for Phase 1. A reader that sees a higher value **refuses** the record (ADR-001 D6) and never rewrites it. Per-record rather than per-file so a truncated tail cannot make the whole ledger unreadable |
| 2 | `event_id` | string | **always** | UUIDv4, lowercase, 36 chars. Unique within the ledger. Never reused. **Never derived from sequence, position or content** (FC1-06) |
| 3 | `conversation_id` | string | **always** | Must equal the conversation directory's id. A record whose value disagrees with its location is **refused**, not relocated |
| 4 | `sequence` | integer | **always** | ≥ 1, strictly increasing by exactly 1, assigned under the lock (D-2). A gap is an error that marks the ledger `unresolved` from the gap; sequences are never renumbered |
| 5 | `event_type` | string | **always** | One of §3's six ledger kinds. Unknown value → record is `unresolved`, not skipped |
| 6 | `timestamp` | string | **always** | RFC 3339 UTC with nanoseconds. The **event** time, distinct from commit evidence |
| 7 | `idempotency_key` | string | **always** | Unique per conversation. Repeat with identical payload → existing result returned, no new record. Repeat with different payload → **refused** (FC1-10) |
| 8 | `actor` | string | **always** | `user` \| `assistant` \| `system` \| `recovery`. Phase 1 does not invent finer provenance |
| 9 | `durability` | string | **always** | `requested` \| `attempted` \| `observed` \| `verified` \| `unknown`. Phase 1 writes **`observed`** — the file was replaced. `verified` is **not** claimed: power-loss durability remains `UN` |

### 2.2 Conditional fields

| JSON name | Type | Required when | Invariant |
|---|---|---|---|
| `turn_id` | string | the event is turn-scoped (`turn_committed`, `turn_corrected`, `turn_retracted`) | UUIDv4 assigned at commit; **immutable across corrections and retractions**; never derived from sequence |
| `predecessor_event_id` | string | ~~recommended on turn events; optional otherwise~~ **superseded: required for every event with `sequence ≥ 2`, absent only at `sequence 1`** — `TASK-03-CONTRACT-FINAL-2026-11-22.md` §2.5 | Must equal the previous record's `event_id` in the **same** ledger. Records an ordered predecessor relation only — it does **not** imply branching or a state change (`FC1-14b`) |
| `payload` | object | **always** | Shape fixed per event kind in §3. Unknown fields inside `payload` are preserved on read and never rewritten |

### 2.3 Deliberately absent in Phase 1

| FC1-08 field | Disposition | Reason |
|---|---|---|
| `state_id` / `resulting_state_id` | **Derived, not stored**: `state_id` is the `event_id` of the turn event that produced the state | Storing it duplicates a value that is already an identity; FC1-14a forbids manufacturing identities |
| `affected_ids` | Not present | Its only Phase 1 use would be a bulk clear, which is expressed as one `turn_retracted` per turn instead |
| `operation` | Not present as a separate field | `event_type` plus `payload` carries the same information without a second vocabulary |
| `provenance` | Reduced to `actor` | Finer provenance has no Phase 1 producer |

**Envelope total: nine required, three conditional, four deliberately absent.** That is the exact list to approve.

---

## 3. `FC1-11` — the exact event kinds

### 3.1 The split, and why the earlier count was wrong

FC-1's semantic distinction "deletion recorded" cannot live in the ledger: Option A deletes the ledger itself. And
`A-12` forbids writing to the source while reading it, so a recovery record cannot live in the ledger either. The
eight semantic distinctions are preserved; **two of them live in the evidence store**, not in the conversation
ledger.

| Kind | Home | Replaces today's operation |
|---|---|---|
| `conversation_created` | **ledger** | `NewChat` |
| `turn_committed` | **ledger** | `AppendHistoryEntry` |
| `turn_corrected` | **ledger** | the edit half of `RewriteHistory` |
| `turn_retracted` | **ledger** | the delete half of `RewriteHistory`, and clear-history |
| `selection_changed` | **ledger** | `SetLastChatID` / `localStorage` `howl-chat-current-chat` |
| `selection_resolved` | **ledger** | nothing — new (D-3) |
| `conversation_deleted` | **evidence store** | `DeleteChat` (which today records nothing) |
| `recovery_recorded` | **evidence store** | nothing — new (A-12 reports, never repairs) |

### 3.2 Required payloads and invariants

| Kind | Required payload | Invariants |
|---|---|---|
| `conversation_created` | `{}` (optional `name`) | Must be `sequence` 1 and the first record. A ledger whose first record is not this kind is `unresolved` |
| `turn_committed` | `role` (`user`\|`assistant`), `content` (string), optional `character_id`, optional `has_image` | Requires `turn_id`. `role` must agree with `actor` for user/assistant turns. Appending produces a new state |
| `turn_corrected` | `target_turn_id`, `replacement_content` | **Never modifies the target record.** `target_turn_id` must already exist. The correcting event carries its own `turn_id` = the target's, so the turn keeps its identity |
| `turn_retracted` | `target_turn_id`, optional `reason` | **Never modifies or removes the target record.** The projection excludes it; the bytes remain |
| `selection_changed` | `state_id` (the `event_id` of a `turn_committed` or `turn_corrected`) | The target must exist and must not be retracted. Never a dangling pointer |
| `selection_resolved` | `deleted_state_id`, optional `chosen_state_id`, `rule` (`ancestor`\|`empty`) | Written by the deletion path. `rule=ancestor` requires `chosen_state_id`; `rule=empty` forbids it |
| `conversation_deleted` | operation id, timestamp, authority, record count, snapshot count, optional last-record hash (§5) | Written to the evidence store **before** the conversation directory is removed, so an interrupted deletion leaves a receipt that says `partial` |
| `recovery_recorded` | observed boundary, record count seen, ledger hash before and after, outcome | Written to the evidence store, never the ledger. Read paths must not write it themselves — they report, and the caller records |

### 3.3 Anti-erasure: correction and retraction cannot destroy history

Four properties, each testable:

1. **No in-place mutation.** A correction or retraction is a new line with a new `sequence`, `event_id` and
   `idempotency_key`. The target line is not opened for writing.
2. **Byte preservation.** After a correction, the target's line is **byte-identical** to its prior bytes and is still
   present. This is the D-1 prefix test applied at the level of a single record.
3. **Projection, not deletion.** "Erased" messages exist only in the projection, which marks the turn `retracted` or
   shows the replacement. The ledger still contains the original.
4. **No compaction escape hatch.** Phase 1 has no compaction, so there is no code path that could rewrite history "to
   save space". Introducing one later requires its own authorization and receipt.

**Documented cost:** clearing a 5,000-turn history writes 5,000 `turn_retracted` records. That is deliberate — a bulk
"clear" event that names a set is cheaper but leaves the ledger unable to say, per turn, what happened to it.

### 3.4 What is deliberately not an event kind

`alternative`/`branch` (`FC1-14b` excluded), `snapshot_captured` (`FC1-19a` unresolved), `memory_copy_committed`
(M-4 / FC1-22e–m out of scope). Adding them later is additive; no required field changes.

---

## 4. `FC1-14c` then `FC1-14b` — reconciled against the schema

### 4.1 Turn and message in Phase 1

| Concept | Phase 1 representation | Restriction, stated |
|---|---|---|
| Turn | `turn_id` (UUIDv4), assigned by `turn_committed`, immutable thereafter | One turn cannot contain a second message |
| Message | `message_id` **equals** `turn_id` | Two FC-1 identities collapse to one value by restriction, not by redesign |
| State | Derived: `state_id` = the `event_id` of the event that produced it | Never stored; never manufactured |
| Chain | `sequence` order is the chain; `predecessor_event_id` states it explicitly on turn events | Exactly one parent, always the immediately preceding state |
| Correction effect | New event → **new state** (`state_id` = correction's `event_id`), same `turn_id`, replacement content | A correction changes conversation state, so it produces a state — but it does not create a new turn |
| Retraction effect | New event → **new state**, `turn_id` retained, projection excludes the turn | The original `turn_committed` remains in the ledger |

**One consequence to accept explicitly:** because a correction produces a new state, repeatedly correcting a turn
produces a chain of states. Selection can point at any of them; the projection of a state is the turn set as of that
event.

### 4.2 Branching (`FC1-14b`) — excluded from Phase 1

- Phase 1 has **no branch event**, no multi-parent state, and no merge. `predecessor_event_id` is a single ordered
  pointer, and FC1-14b's own wording says it does not settle branch semantics.
- **Fork-as-new-conversation is recorded as a possibility, not adopted.** If it is ever chosen, a fork would create a
  *new* `conversation_id` with its own ledger whose first event references `origin_conversation_id` and
  `origin_state_id` — both **optional** fields, so adding them later does not change any required field or rewrite
  any existing record.
- Because it is only a possibility, `conversation_created`'s payload is left open to those optional fields and nothing
  in Phase 1 populates them.

### 4.3 What remains undecided after this reconciliation

Whether fork-as-new-conversation is the intended shape (recorded, not decided) · whether a turn may ever contain more
than one message · whether a selection occurrence needs its own identity (`FC1-14d`) · `FC1-19a` snapshot contents.

---

## 5. Deletion receipt — revised after owner review

The earlier field list is tightened. **Removing the conversation ID did not make the rest non-identifying.**

| Field | Status | Note |
|---|---|---|
| Operation id | **required** | Random, unrelated to the conversation |
| Timestamp | **required** | Exact; proves when. Coarsening is possible but weakens the evidence, so it is left exact and the residual correlation risk is stated rather than hidden |
| Authority | **required** | Who requested the deletion |
| Record count | **required** | Low-entropy, but proves the deletion set size |
| Snapshot count | **required** | Same |
| **Last-record hash** | **optional — proposed default: absent** | **It is a fingerprint**: anyone holding a surviving copy of the ledger can match it and identify the deleted conversation. It must earn its place operationally or be dropped |
| Salted conversation-ID hash | **withdrawn** | Enumeration of timestamp-named ids defeats it (decision record §10.2.1) |
| Conversation title | **prohibited** | — |
| Plaintext conversation ID | **prohibited** | — |

**Residual correlation risk, stated:** timestamp + counts narrow a receipt to a particular deletion. That is inherent
to a receipt that says when and how much; it is not a fingerprint of content, and this is the difference the optional
hash removes.

**Question for the owner, not answered here:** does the last-record hash provide enough operational value — proving
*which* history ended here — to justify retaining a fingerprint of deleted content? Default proposed: **no**.

---

## 6. Two counts that must stay distinct

The preparation package has **seven disposition topics** (`FC1-08`, `FC1-11`, `FC1-02`, `FC1-07`, `FC1-14b`,
`FC1-14c`, `FC1-PRIVACY-DELETE`) and **eight owner decisions** — the same seven plus **whether per-conversation
deletion proof is required at all**. Deletion-proof is a separate question because it is about *whether any
identifier-bearing proof is wanted*, not about how a receipt is shaped. They are consistent only when counted
separately, and the record keeps them separate.

## 7. Authorization

Documentation only. No implementation, no call-site migration, no ledger creation, no persistence adoption, no
contract approval. `FC1-08`, `FC1-11`, `FC1-02`, `FC1-07`, `FC1-14b`, `FC1-14c`, `FC1-PRIVACY-DELETE` all remain
undecided until the owner says otherwise.
