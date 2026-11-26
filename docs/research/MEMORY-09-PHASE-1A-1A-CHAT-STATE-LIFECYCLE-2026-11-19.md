# MEMORY-09 — Phase 1A-1A: Chat-State Lifecycle — Invalidation, Restoration and Applicability

**Status:** research addendum. Read-only. **No implementation, no test, no probe, no model call, no
corpus, no canonical edit, no status vocabulary.** Nothing here adopts a model, names a status value,
amends a requirement, or decides Q-D or Q-B.

| Field | Value |
|---|---|
| Created | **2026-11-19** |
| Authority | None beyond the owner's follow-up to Phase 1A-1. Creates no decision and records no new T1 entry. |
| Scope | Whether a chat-state transition (regenerate / next alternative / undo / revert) is a memory lifecycle operation, what the sources say about it, and what Howl and MEMORY-09 actually do today |
| Predecessor | `MEMORY-09-PHASE-1A-1-RETRACTION-AND-KNOWLEDGE-PERSISTENCE-2026-11-19.md` (RP-01–RP-14, RQ-01–RQ-06, RA-1–RA-3) |
| Inputs verified | Hashes in §3, read from disk this session |
| Tests executed | **0.** Every behavioural statement is a static reading of source. Nothing was run, and nothing was rendered. |
| Files changed | This file, and one cross-reference section added to the 1A-1 report. Nothing staged, committed, or registered. |
| **Scope boundary (2026-11-19)** | The state machine is recorded as a **Howl** requirement, not a MEMORY-09 one: `docs/architecture/REQUIREMENTS-HOWL-CHAT-STATE-LIFECYCLE-2026-11-19.md` (HSL-1–HSL-11). This file remains research: RPS findings, RQS questions, RS shapes. It adds nothing to the MEMORY-09 experiment plan, corpus, probes or conformance rows. |

---

## 1. Why this record exists

Phase 1A-1 separated five operations that are routinely treated as one — forgetting, source
retraction, account retraction, correction, supersession — and found that both implementations
express every one of them by writing a status on a record.

The owner's follow-up names a sixth operation that neither the sources nor Phase 1A-1 separated:

> a chat state transition (regenerate, next, revert, or another state transition) changes the
> conversation, and the memory system must follow that state history without destroying anything.

The failure mode is specific and severe: a memory system that handles correction and retraction
correctly can still lose a character's knowledge every time the user presses **Next**, or can present
knowledge from a conversation branch the user abandoned.

This record treats that as its own subject. It does **not** fold it into the retraction review, because
the operation is different in kind: retraction changes what a source asserts; a state transition
changes which conversation is current, and asserts nothing about anything.

**Boundary, settled 2026-11-19.** The state machine, branching, regeneration, retention and deletion
are **Howl's** responsibilities, now recorded as requirements HSL-1–HSL-11 in
`docs/architecture/REQUIREMENTS-HOWL-CHAT-STATE-LIFECYCLE-2026-11-19.md`. This record does not claim
them for MEMORY-09 and does not turn any finding below into an experimental variable. It investigates
what a state-aware memory representation would require, and what MEMORY-09 may and may not infer —
which is MEMORY-09's side of the boundary, together with the contract that will eventually cross it.

Three results were not visible in Phase 1A-1:

1. **No source governs this operation at all.** r32 and v9 contain no rule, invariant, or example
   about regenerate, alternatives, branch, revert or undo. Chat *deletion* is governed in detail
   (K-R35, K-R40–K-R44); chat *state* is silent. The 118 requirements likewise have no state row
   (RPS-01, RPS-10).
2. **Howl's current chat machine cannot supply the provenance the owner's principle 2 requires, and
   it is destructive in one direction and possibly duplicative in another** (RPS-02, RPS-07).
3. **`status` cannot carry this dimension without collapsing two distinctions the owner explicitly
   wants kept.** Applicability has to be a separate axis from lifecycle status, and the shape that
   already fits it is r32's projection-not-field pattern (K-R12/K-R13), not a fourth enum value
   (RPS-03).

---

## 2. Method and evidence labels

Labels and tiers are unchanged from Phase 1A-1 §2 and are not redefined here. **Observed** means read
in the working tree at the hash in §3, never run. **Inferred** is a judgement this record makes and
invites challenge on; where a runtime consequence is inferred from a static code path, it says so and
is not presented as demonstrated.

**A note on the owner's five principles.** They were stated in this session as premises *for the
investigation*. They are treated here as the owner's stated design premises, not as recorded T1 ledger
decisions, because no ledger entry contains them. Where a principle is already entailed by T1/T2/T3
text, the entailment is cited; where it is not, the principle is assessed as a premise and the gap is
recorded rather than closed.

**ID namespaces** (new, collision-free with `CO-`, `PH1A-`, `RP-`, `RQ-`, `RA-`, `RT-`):

- `RPS-nn` — findings
- `RQS-nn` — open questions requiring an owner decision
- `RS-n` — relation alternatives
- `RST-nn` — implementation traces (Appendix A)

**Terminology, stated because the word is already overloaded (§6.5).** This record calls the owner's
subject a **conversation state** (CS), and a transition between such states a **conversation-state
transition**. "Chat state" is avoided as a term of art here because `chat_state` is already a MEMORY-09
*scope* with a different meaning. The term is a proposal for discussion, not a decision.

---

## 3. Sources read, with hashes and tiers

| Source | SHA-256 (first 16) | Role here |
|---|---|---|
| `MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md` (r32) | `e88ef2d1697da1d9…` | T2; K-R12/R13/R19, K-R35, K-R40–K-R44, §8.4 |
| `MEMORY-09-ARCHITECTURE-v9-2026-11-17.md` | `698254ca262c0874…` | T3; rules 2, 12, 20, 27–31, 38, 42 |
| `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` | `da3e72398ecbaa0a…` | T1; D-01, O-1, O-3, O-5, O-6a, O-13 |
| `MEMORY-09-CONFORMANCE-MAP-2026-11-17.md` | `b4ba33ccbfe8f8c1…` | 118 ratified requirements |
| `MEMORY-09-PHASE-1A-1-…-2026-11-19.md` | `82173a6f…` (predecessor, unchanged at its own §3) | RP/RQ/RA baseline |
| `eval/mem1/m09/records.py` | `ad2d4d11501c357b` | `Record` fields, `STATUSES` |
| `eval/mem1/m09/retrieval.py` | `cd2d380bd5d9d1f7` | candidate pool construction |
| `eval/mem1/m09/engine.py` | `7505e9acca8e3229` | operations; no state parameter |
| `eval/mem1/model.py` | `813cdbff56c01ae4` | `SCOPES`, harness record shape |
| `eval/mem1/form.py` | `f472052738388bab` | `chat_state` scope use; `state_change` kind |
| `eval/mem1/conversation.py` | `fa1d9ce630462fa9` | declared scope state, not CS |
| `frontend/dist/js/chat-actions.js` | `577f7e570b7f5b61` | undo, delete, edit, regenerate, prev/next |
| `frontend/dist/js/chat-render.js` | `6bc1689c329b59a4` | bubble action wiring |
| `frontend/dist/js/tutorials-data-chat.js` | `c4defa9e4a585d5d` | documented behaviour of Prev/Next |
| `app-prompt.go` | `93c9e6b875993619` | `ReplaceChatHistory`, `ClearChatHistory` |
| `app-chat.go` | `f08037a8e3aaf32d` | `buildChatHistory`, final-save rewrite |
| `internal/backend/chat/types.go` | `aa9f1f71adcc849b` | `ChatHistory` |
| `internal/backend/chat/persistence.go` | `ad9663a48278cff0` | append log, atomic rewrite, load preference |
| `internal/backend/types/context.go` | `14cdb64c17a6ef2f` | `RemoveMessageIDs` (unused by the chat path) |
| `Chats/2026-05-29T12-28-48Z/memories.json` | (data sample, not hashed) | Howl's live memory record shape |

---

## 4. The causes, separated

### 4.1 The owner's five, stated precisely

| Cause | What changed | What must persist | What a character may see |
|---|---|---|---|
| **New information** | an earlier claim is contradicted or updated | the old claim, and the replacement, and the fact of replacement | both, with the later one governing only where authority says so |
| **Source retraction** | the source withdraws its assertion | the source record, its provenance, every knower's account | the account, unchanged; whether the withdrawal is visible is RQ-01 |
| **Mistaken presentation** | information was presented when it should not have been | the history of the presentation and the error | the correction, not a silently cleaned history |
| **Relevance change** | the information is no longer useful now | the record, unchanged | nothing changes; retrieval ranking changes |
| **Conversation-state change** | regenerate / next / undo / revert changes which conversation is current | the affected records, and which state they belong to | only what applies to the current state, plus whatever history surfaces are authorised |

The last row is the new subject. Its distinguishing property is that **nothing in the fiction changed
and no one's assertion changed**: the same claim is true or false exactly as before, and the same
character knows exactly what they knew. Only *applicability to the current conversation* changes.

### 4.2 Two operations Phase 1A-1 separated that this table does not show

- **Forgetting** removes a character's knowledge and is not one of the five causes above. It remains
  governed by K-R3/K-R19 and RP-01.
- **Account retraction** removes one character's copy while leaving the source intact (RQ-03). It
  belongs to the "source retraction" family above only in the sense that a source is involved; its
  effect on knowledge is the opposite.

Neither is a state operation, and neither becomes one here.

### 4.3 Why the state row is not any of the other four

| Test | New information | Source retraction | Mistaken presentation | Relevance change | State change |
|---|---|---|---|---|---|
| Does a record's content change? | no — a new record is added | no | no — an error is recorded | no | **no** |
| Does any record's truth status change? | supersession may change which governs | **yes** | possibly, via correction | no | **no** |
| Does any character's knowledge change? | no | no (O-3) | no | no | **no** |
| Does *what is retrievable now* change? | yes, by content | sometimes | sometimes | **yes, by ranking** | **yes, by applicability** |
| Is the change caused by the fiction? | yes | yes | yes | no — caused by the user's need | **no — caused by the interface** |

The last row is decisive. A state transition is the only cause in the set that is caused by the
interface rather than by the fiction, which is exactly why the owner's principle 5 assigns the state
machine to the chat and not to the memory system.

---

## 5. What Howl actually does today

All Observed at the hashes in §3 unless marked Inferred.

### 5.1 The current conversation is a single mutable list

`ChatHistory` is `{id, sender, text, role, timestamp, hasImage, characterId}`
(`internal/backend/chat/types.go:72`). There is no parent, no branch, no alternative set, no state
id, and no tombstone. `context.go`'s `RemoveMessageIDs` exists but the chat path does not use it; the
frontend path calls `ReplaceChatHistory` instead.

**RST-01 — message identity is positional.** `buildChatHistory` assigns
`ID: fmt.Sprintf("%d", len(history))` (`app-chat.go:405`), and `ReplaceChatHistory` assigns
`fmt.Sprintf("%d", i)` while rebuilding (`app-prompt.go:327–390`). A message's id is therefore a
function of how many messages precede it. Insert, delete, prune or reorder renumbers everything after
the change.

### 5.2 Undo destroys

`undoMessage` calls `_pruneMessagesAfter(messageDiv, true)` (`chat-actions.js:88–108`), which removes
the target and every later DOM node, then `_syncDomHistoryToBackend` (`chat-actions.js:48–86`) calls
`ReplaceChatHistory`, which calls `RewriteHistory` (`app-prompt.go:389`).
`RewriteHistory` atomically replaces `history.ndjson` and `history.json`
(`persistence.go:631–661`), and its own comment states the intent: *"so deletions and edits actually
remove messages."*

So **undo is a durable destructive operation on the conversation.** There is no branch behind it and
no way back to the state that was removed.

### 5.3 Regenerate discards the reply, but may not discard it on disk

`regenerateMessage` removes the last assistant bubble from the DOM and re-streams from the last user
text (`chat-actions.js:273–304`). Unlike undo and delete, it does **not** call
`_syncDomHistoryToBackend` before streaming.

The append log is the preferred load source (`persistence.go:583–629`), an assistant reply is appended
to it (`persistence.go:668`), and the log is only rewritten from current state at final save
(`app-chat.go` `performFinalSave`). **Inferred from that path, not executed:** the removed reply can
survive in `history.ndjson` and be loaded alongside the replacement, so regenerate may *duplicate*
where undo *deletes*. This is stated as an inference with its code path, not as a demonstrated defect;
it is recorded because the owner's model assumes regenerate retains exactly one alternative and
discards onto a branch, and the current machine does neither reliably.

### 5.4 Delete and edit are the same destructive family

`deleteMessage` removes one node and syncs (`chat-actions.js:35–46`); `editMessage` on commit prunes
everything after the edited message and, for a user turn, re-sends (`chat-actions.js:124–232`). Both
route through the same rewrite. Editing history is therefore implemented as truncating it and
regenerating forward.

### 5.5 "Prev / Next" does not browse alternatives

`previousMessage` and `nextMessage` (`chat-actions.js:329–348`) scroll the chat pane to the adjacent
assistant message. They do not select between alternatives, and no alternative set exists to select
from. The in-app documentation says the opposite: *"Prev / Next — Navigates through generation
variations if you have regenerated a message multiple times"*
(`tutorials-data-chat.js:117`). **Documented behaviour and implemented behaviour disagree**, and the
owner's step 2 ("Howl generates an alternative response") describes the documented behaviour, which
the code does not have.

### 5.6 What survives a state transition today

| Survives? | Messages | Character memories |
|---|---|---|
| Undo | **no** — rewritten away | **yes**, silently — records are not touched |
| Regenerate | possibly the discarded reply (Inferred) + the new one | yes |
| Delete | no | yes |
| Edit | truncated after the edit | yes |

### 5.7 Howl's live memory store has no state anchor

The sample store (`Chats/…/memories.json`) records `id`, `title`, `summary`, `content`, `scope`,
`owner_id` (the chat id), `tags`, `status: active`, `source: auto`, timestamps. There is no message
id, no turn number, and no state reference. Whatever produces those entries, the record shape cannot
say which conversation state produced a memory or which state it belongs to.

---

## 6. What MEMORY-09 provides today

### 6.1 The record carries containment and a turn, but no state

`Record` carries `record_id`, `record_class`, `chat_id`, `text`, `created_turn`, owner, domain,
`proposition_id`, `provenance`, `status`, `supersedes`/`superseded_by`, `promotion_id`,
`corrected_by` (`records.py:75–110`). There is **no** field for a conversation state, a branch, a
lineage, a message, or an alternative.

`status ∈ {active, superseded, revoked}` (`records.py:72`) is the entire lifecycle vocabulary.

### 6.2 No operation accepts a state

`Engine.retract(record_id, turn)`, `disclose(...)`, `acquire(...)`, `correct(...)`,
`supersede_authority(...)`: none takes a state parameter, and no method exists that could receive one.
There is no "state changed" event, no applicability field to write, and no projection to recompute.
Observed across `m09/engine.py`.

### 6.3 Retrieval has no state filter

`retrieval.perspective` builds its pool from `engine.all_records()` filtered only by
`owner_character == character and status != "revoked"` (`retrieval.py:46–52`). `world` and
`belief_facts` pass through `engine.authoritative` on domain and target. **Nothing consults a state.**
If two conversation states both produced records in the same chat, both sets are candidates in both
states.

### 6.4 `created_turn` cannot identify a state

The turn is an integer supplied by the driver. A regenerated turn is the same turn number with
different content, so `(chat_id, created_turn)` is not unique across states, and a record formed in
state B is indistinguishable by turn from one formed in state A at the same position. This also
matters for evaluation: a probe keyed `(character, turn)` becomes ambiguous the moment a state
transition exists (RPS-05).

### 6.5 "State" already has four unrelated meanings in this codebase

| Term | Where | Means |
|---|---|---|
| `chat_state` | `model.py:35` `SCOPES`; `form.py:168`; `retrieve.py` context map | a **visibility scope**: facts belonging to the chat container |
| `state_change` | `model.py:48` `KINDS`; `form.py:223–226` | a **record kind**: a value that changed and supersedes the prior value |
| `state_changes` / `active_state` | `fixtures.py:53–58`; `conversation.py:144–147`; `retrieve.py:176` | **declared world/scenario scope activation** over turns |
| conversation state (CS) | this record | **which conversation is current** |

Four meanings, one word. Any new identifier using "state" without a qualifier will be ambiguous on
sight, and `chat_state` in particular must not be reused for CS.

---

## 7. The owner's five principles, assessed against the sources

Premises are the owner's (see §2). "Entailed" means existing text already requires it; "gap" means no
text addresses it.

| Principle | r32 | v9 | 118 requirements | Today's code |
|---|---|---|---|---|
| **P1 Preservation** — a transition must not destructively delete memories of the previous state | **entailed** by K-R3/rule 20 ("nothing is erased"), but only for memory operations; nothing addresses message deletion | **entailed** by rule 20 | no row; the closest is E1a/E2 | **Howl violates it at the message layer** (undo rewrites); MEMORY-09 satisfies it **vacuously** — nothing ever deletes a record |
| **P2 State association** — provenance sufficient to establish which state produced a memory | **gap** | **gap** | no row | **gap**: no state field; `chat_id` + `created_turn` do not discriminate (§6.4); message ids are positional (§5.1) |
| **P3 Reversible applicability** — inactive ≠ nonexistent or forgotten | **shape exists** (K-R12 status-as-projection; K-R13 computed-at-read), **not applied** | **gap** — `status` is a mutable field | no row | **gap**: no dimension; `revoked` is the only exclusion and it means retraction/forgetting |
| **P4 Replacement** — a new state's memory consequences are added, not overwritten | **entailed** by rule 20 / K-R3 for records; **silent** for presentations | **entailed** by rule 20 | E2 for supersession only | records: **satisfied** (nothing overwrites); presentation: **violated** by regenerate/undo |
| **P5 Ownership** — the chat owns the state machine; MEMORY-09 does not decide abandonment | **consistent** with K-R40 ("no automatic transfer"; manual acts name source and destination) and with the owner's D-01/O-13 | rule 31 assigns cross-chat reach to the application | no row | **no interface exists** for the chat to inform memory of anything |

Two observations matter more than the table:

- **P1 and P3 pull in opposite directions if satisfaction is claimed by doing nothing.** MEMORY-09
  currently preserves everything because it has no deletion path, which satisfies P1 and P4 by
  accident and fails P3 completely. Howl currently deletes presentations, which satisfies P3's spirit
  and violates P1. Neither layer alone satisfies the pair.
- **P5 is the load-bearing one.** If the chat deletes state A's messages (it does), and MEMORY-09 is
  required to preserve A's memories (it must), then **memory becomes the only surviving record of a
  conversation the chat has erased**. That is a real consequence of P1 + P5 that no source has
  considered, and it is recorded as RPS-08 rather than resolved.

---

## 8. The owner's worked sequence, traced

Sequence: state A → *"The merchant gave us a silver key."* → Next → alternative *"…a map."* → state B
→ revert to A.

| Step | Chat layer (today) | MEMORY-09 (today) | Required by the owner's model |
|---|---|---|---|
| 1. A: silver key | message appended, id = position | account record, `chat_id` + `created_turn`, `status=active` | active in A |
| 2. Next | new reply appended; the A reply is removed from the DOM, **not** rewritten from disk (Inferred) | **no change** — nothing tells it a state changed | key record retained, inapplicable in B |
| 3. B: map | map reply appended | new account, `active`, same turn integer | active in B; key record inapplicable but present |
| 4. Revert to A | **not implemented** — there is no revert action; the nearest real operation is undo, which would delete B's messages | the map record remains `active` and retrievable; the key record was never deactivated | key applicable again; map preserved but inapplicable |

**What each surface would show, under the owner's model versus today:**

| Surface | Owner's model | Today |
|---|---|---|
| Current conversation | one branch, switchable | one list, editable only by destruction |
| Character's retrieval in A | key, not map | **key and map** (both `active`, no state filter) |
| Character's retrieval in B | map, not key | **key and map** |
| History view | both, with their states | both, with no state to show |
| After undo of B | A intact | A is gone from disk; **only memory still holds it** |

That last cell is the finding. The owner's concern was that pressing Next would lose memories. The
static reading says the current machine has both faults: it loses *messages* on undo, and it would
never lose *memories* — it would instead leak abandoned-branch memories into the restored branch,
because nothing can mark them inapplicable.

---

## 9. Findings

| ID | Finding | Label | Source |
|---|---|---|---|
| **RPS-01** | No source governs conversation-state transitions. r32 and v9 contain no rule about regenerate, alternatives, branch, revert or undo; chat *deletion* is governed in detail (K-R35, K-R40–K-R44), chat *state* is silent | Observed + Calculated (zero occurrences) | r32; v9; conformance map |
| **RPS-02** | Howl's chat machine has no state identity to publish: one mutable list, no parent, no branch, no alternative set, no tombstone (`types.go:72`) | Observed | `internal/backend/chat/types.go` |
| **RPS-03** | `status` cannot carry applicability. It is 3-valued and already means retraction/forgetting; adding "inactive" would make one field carry two axes, and a reverted state could not restore activity without a fourth value — the enum collapse the owner ruled out | Observed + Inferred | `records.py:72`; 1A-1 RP-01 |
| **RPS-04** | MEMORY-09 retrieval is state-blind: the pool is owner + `status != revoked`; records from an abandoned branch are candidates in every branch of the same chat | Observed | `retrieval.py:46–52` |
| **RPS-05** | `created_turn` cannot discriminate states: a regenerated turn reuses the integer, so `(chat_id, turn)` is not unique across branches; probe keys `(character, turn)` become ambiguous once a transition exists | Observed + Inferred | `records.py:81`; `conversation.py` turn loop |
| **RPS-06** | "State" already carries four unrelated meanings (`chat_state` scope, `state_change` kind, declared scope `active_state`, and CS). Any new identifier must be qualified | Observed | §6.5 |
| **RPS-07** | Howl's documentation and code disagree about Prev/Next: the tutorial promises variation browsing, the code scrolls. The owner's step 2 describes behaviour the code does not implement | Observed | `tutorials-data-chat.js:117`; `chat-actions.js:329–348` |
| **RPS-08** | Preservation and ownership compose into an unconsidered consequence: because the chat destroys messages on undo, memory would be the sole surviving record of an erased conversation | Inferred | §5.2; §7 P1/P5 |
| **RPS-09** | A retraction is itself an event *in* a state, so source status is state-relative. Q-B's derived notice, if adopted, must join on `(state, status)`, not on a global status | Source-derived + Inferred | r32 §8.4, K-R12/R13; §10 below |
| **RPS-10** | No requirement in the 118 has a state dimension; A2a/A2b/E1a–E1c all assume one linear history. A conformance row cannot adjudicate a behaviour that differs by state | Calculated (row inspection) | conformance map |
| **RPS-11** | The only lifecycle vocabulary MEMORY-09 has is a mutable `status` field, while r32's shape for "known but no longer asserted" is a projection over an event stream. The state subject needs the projection shape, not the field | Observed + Source-derived | K-R12, K-R13; `engine.py` |
| **RPS-12** | Regenerate does not rewrite the history log before streaming, so the discarded reply may survive and reload beside its replacement (Inferred from path; not executed) | Inferred | `chat-actions.js:273–304`; `persistence.go:583–629, 668` |

---

## 10. Why this is not retraction — and the one place they couple

A state transition is not:

- **retraction** — no source withdraws anything, and the claim's truth status is unchanged;
- **forgetting** — no character stops knowing anything;
- **correction** — no content changes;
- **supersession** — no record stops governing;
- **relevance change** — relevance is a ranking judgement about usefulness; applicability is a fact
  about which conversation is current.

The coupling is narrow and real: **a retraction event occurs at a point in a conversation state.** If
A retracts the key in state A and the user reverts to a state before the retraction, the retraction
did not happen in the current conversation. A status field that records "this source is retracted"
globally cannot express that; a projection over state-scoped lifecycle events can.

This is RPS-09, and it is the reason Q-B cannot be answered independently of Q-D. Phase 1A-1's
R-narrow reading requires a read-time join against the source's retraction status. That join is
well-defined only once the right-hand side is `(state, status)`. The state subject therefore does not
*change* Q-B's question; it changes what a correct answer to Q-B has to be computed from.

It also raises the cost of one Q-B answer. R-broad ("a knower learns only by being told") requires
acquisition events for the notice. Under state transitions those acquisition events are themselves
state-scoped — a notice acquired in an abandoned branch must not govern the restored one. So R-broad
now needs state-scoped acquisition, which no source provides. This is a cost observation, not a
ranking: R-broad is not refuted, and §7 of the predecessor's three readings remain unranked.

---

## 11. Revisiting Q-B

**What does not change.** Q-B asks whether a knower may be shown that what they hold was retracted
without having acquired that fact. The state subject does not answer it. §8.4 still says yes for any
live projection; A2b/E1c still forbid unaccounted retraction knowledge; the three readings
(R-narrow, R-broad, R-split) remain each defensible and none is selected by any source. **No part of
the predecessor's §7 is withdrawn.**

**What changes, in three ways:**

1. **The join is state-relative** (RPS-09). "Live projection" must be read as "live in the current
   conversation state". A notice computed from a global status would report a retraction that did not
   happen in the branch the user restored.
2. **R-broad's cost rises.** It needs state-scoped acquisition events, not just state-scoped status.
   R-narrow's cost also rises, but differently: it needs a state-aware projection over the memory
   stream, which is the shape r32 already has (K-R12/K-R13) and v9 does not (RPS-11).
3. **The ordering becomes forced.** Q-B presupposes an answer to Q-D, because it cannot say what the
   notice joins to without one. Answering Q-B first would silently decide the global-status model,
   which is exactly the collapse the owner ruled out.

**One thing this record deliberately does not do:** it does not convert Q-B into a question about
applicability. "This does not apply to the current conversation" and "you were told this and it was
withdrawn" are different statements to a character, with different epistemic standing, and merging
them would be the same error as merging retraction with forgetting.

---

## 12. Relationship alternatives for conversation state

Stated for the owner's decision. **No alternative is recommended, and none names a status value or
proposes an enum** — each describes where the state relation lives, not what its values are called.

**Rescoped 2026-11-19.** Because the state machine belongs to Howl (HSL-1), these are now candidate
**contract shapes** between Howl and MEMORY-09, not MEMORY-09-internal architecture choices. Howl can
guarantee stable state identity and an ordered transition history in all three; what differs on the
MEMORY-09 side is where applicability is represented. The choice is therefore joint, and it is HSL-O5
in the requirements record.

### RS-1 — State-tagged records (identity on the record)

Each record carries the conversation state (or lineage) in which it was formed. Applicability is a
comparison between the record's tag and the current state. History surfaces query without the filter.

- Satisfies P1, P2, P3, P4 directly; P5 requires a published current-state identifier from the chat.
- Cost: every record gains a field, and every existing record needs a value it does not have. The
  scope-free semantics of `world_truth`/`scenario_truth`, which are chat-scoped rather than
  turn-scoped, have to be decided explicitly — a promoted truth that outlives its branch is RQS-04.

### RS-2 — Current-state applicability as a projection (nothing on the record)

Records carry no state. The chat publishes an ordered set of state transitions; MEMORY-09 derives
which records are applicable by joining record creation against those transitions at read time.

- Mirrors K-R12/K-R13's "status is a projection, never a field" and adds no per-record field.
- Cost: the projection is only as good as the transition log the chat supplies, and the log is the
  thing Howl does not have (RPS-02). It also makes every retrieval path depend on a join, exactly the
  join v9 refuses for retraction notices.

### RS-3 — State-local namespaces (a store or partition per state)

Each state (or lineage) has its own partition; reverting selects a partition. Nothing is deleted.

- Simplest to reason about for P1/P3; strongest separation.
- Cost: a record that must span states (world truth, a character's standing knowledge formed before
  the branch point) exists in one partition and is invisible from the other unless duplicated, which
  reintroduces the identity problem RA-1/RA-3 were about. History across states becomes a union the
  system has no identity for.

### 12.1 Against the principles

| | P1 preservation | P2 state association | P3 reversible applicability | P4 replacement | P5 chat owns the machine |
|---|---|---|---|---|---|
| **RS-1** tagged records | pass | pass | pass | pass | pass, given a published state id |
| **RS-2** projection | pass | **pass, if the transition log is durable** | pass | pass | **pass, and the chat remains the only authority** |
| **RS-3** partitions | pass | pass | pass | pass | pass |

All three can satisfy the principles on paper; they differ in what they require of Howl and in what
they do to cross-state records. That difference — not the principle table — is the decision.

---

## 13. Questions and the revised brief

### 13.1 Open questions

| ID | Question | Why it cannot be settled from source |
|---|---|---|
| **RQS-01** | Does MEMORY-09 receive a conversation-state identity at all, and does applicability live on the record (RS-1), in a projection (RS-2), or in partitions (RS-3)? | No source mentions the operation; the three shapes are each consistent with the owner's principles |
| **RQS-02** | What is the unit: a state, a branch, or a lineage? Regenerate creates a sibling; undo truncates; revert re-selects a prior state. These are different relations and the owner's model uses all three loosely | No source defines any of them; the three have different persistence and query consequences |
| **RQS-03** | On revert, may anything query the abandoned branch's memories, and who? | P3 says inactive ≠ inaccessible to history; it does not say who may look |
| **RQS-04** | Do records that outlive a branch — promoted `world_truth`/`scenario_truth`, standing character knowledge formed before the branch point — remain authoritative/applicable after the branch is abandoned? | Promotion is chat-scoped (O-13, rule 31); states are inside a chat, so no text reaches this |
| **RQS-05** | Is a state transition a knowledge event for characters? If A's line exists only in an abandoned branch, does the account formed from it remain A's account? | Depends on RQS-01; no source states whether applicability affects ownership or only visibility |
| **RQS-06** | What must Howl guarantee (durable state identity, retained abandoned states) before MEMORY-09 can satisfy P1 at all? | P5 gives the state machine to the chat, but no requirement constrains what the chat must retain; today it retains nothing (§5.2) |

### 13.2 The brief, reordered

Phase 1A-1's brief put three questions (Q-A relation, Q-B notice, Q-C attribution) and held PH1A-12
downstream of Q-A. The state subject adds one question that must precede Q-B, and does not displace
the others:

**Q-D's ownership is assigned; its answer is not.** Howl owns the state machine, state identity and
the current-state determination (HSL-1–HSL-3); MEMORY-09 owns representation, retrieval and dependency
(HSL-9–HSL-11); the crossing contract is deferred to Howl's foundational architecture work. So Q-D is
no longer a question about whether MEMORY-09 should model states — it should not — but about where
applicability is represented on MEMORY-09's side of the boundary. RQS-01 is narrowed accordingly, and
RQS-02/RQS-03/RQS-06 move to the Howl requirements record as HSL-O1–HSL-O3. RQS-04 and RQS-05 remain
MEMORY-09 semantics.

**Q-D. The state relation (RQS-01/RQS-02).** Does MEMORY-09 track conversation states at all, and if
so does applicability live on the record, in a projection over a transition log the chat supplies, or
in partitions? This is RPS-03's consequence: whichever is chosen, applicability must not be a fourth
value of `status`.

**Q-A. The relation between knowledge and the thing known (unchanged).** RA-1, RA-2, RA-3.

**Q-B. The notice (unchanged in substance, now state-relative).** RQ-01, joined on `(state, status)`
per RPS-09.

**Q-C. The attribution (unchanged).** RQ-02.

**Still not put to the owner:** PH1A-12 (experimental subject). It remains downstream of Q-A, and now
also of Q-D, because a state-blind fixture cannot host a state-aware experiment (RPS-10).

---

## 14. What this record does not decide, and its limits

No architecture is adopted. No status value or enum is proposed. No requirement is added, amended or
relaxed. No interface is specified. No implementation change is proposed, and none was made. The
owner's five principles are assessed as premises, not ratified as requirements.

**Limits.** Every behavioural claim is a static reading at the hashes in §3; nothing was executed, so
RPS-12 and RPS-08 are inferences from code paths, not demonstrations. The Howl frontend was read at
`frontend/dist/js/`, which is the served copy; whether a source copy under a build step differs was
not investigated and is not claimed either way. Nothing was rendered, so no claim is made about what
the interface looks like — only about what its handlers do. `Chats/…/memories.json` was read as one
sample record set, not as a census of every stored shape.

---

## Appendix A — Traces

**RST-01** `app-chat.go:405` positional `ID`; `app-prompt.go:327–390` index-assigned `ID`.
**RST-02** `chat-actions.js:88–108` `_pruneMessagesAfter` → `:48–86` `_syncDomHistoryToBackend` →
`app-prompt.go:389` `RewriteHistory` → `persistence.go:631–661`.
**RST-03** `chat-actions.js:273–304` regenerate removes the bubble, streams, no history sync.
**RST-04** `chat-actions.js:329–348` prev/next scroll; `tutorials-data-chat.js:117` documents variations.
**RST-05** `types.go:72` `ChatHistory` — no parent/branch/state.
**RST-06** `records.py:72` `STATUSES`; `records.py:75–110` `Record` — no state field.
**RST-07** `retrieval.py:46–52` pool by owner and status only.
**RST-08** `model.py:35` `chat_state` scope; `model.py:48` `state_change` kind; `fixtures.py:53–58` and
`conversation.py:144–147` declared scope state.
**RST-09** `persistence.go:583–629` load prefers the append log; `:668` append path;
`app-chat.go` `performFinalSave` rewrites the log only at exit.
**RST-10** r32 K-R40–K-R44 govern chat deletion only; no r32/v9 rule names regenerate, branch, revert
or undo.

## Appendix B — Text locations for the load-bearing quotations

| Quotation | Location |
|---|---|
| "Deleting this chat will permanently delete its associated world, scenario and character memories…" | r32 K-R42 |
| "the append log (history.ndjson) is the true source of truth and must NOT be deleted" | `persistence.go:208–210` |
| "so deletions and edits actually remove messages" | `persistence.go:631–634` |
| "Navigates through generation variations if you have regenerated a message multiple times" | `tutorials-data-chat.js:117` |
| "Nothing is erased" (rule 20) | v9 rule 20 |

---

**End of addendum.** No canonical artifact was modified. Q-B is not withdrawn; it is re-scoped by
RPS-09 and now follows Q-D (RQS-01). The next action requires an owner response to Q-D before the
notice question can be answered on a stable footing.
