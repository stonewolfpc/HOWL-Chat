# TASK-03 — FC-1 disposition proposals

| Field | Value |
|---|---|
| Purpose | Dispositions for the six questions that block freezing the Phase 1B ledger contract |
| Status | **Proposals for review. None is adopted. No owner decision is silently resolved.** |
| Authority | FC-1 is `PROPOSED`; ADR-001/002 accepted. A proposal here does not approve a clause |
| Evidence labels | **observed** (read in this tree) · **source-derived** (follows from code) · **calculated** · **inferred** · **unestablished** |
| Recovery token | `TASK-03-FC1-DISPOSITIONS-2026-11-22` |

**One sentence.** Six questions stand between the decided design and a contract that can be frozen; each is answered
here with the clause, the evidence, the options, a recommendation, and whether it blocks.

## Summary

| Question | Recommendation | Blocks freezing Phase 1B? | Needs owner sign-off |
|---|---|---|---|
| `FC1-08` envelope fields | Mandatory core of 8 + a format version; state/branch fields deferred | **Yes** | Yes |
| `FC1-11` Phase 1 event kinds | Eight kinds, including corrections and retractions so `RewriteHistory` can retire | **Yes** | Yes |
| `FC1-02` graph scope | Linear single-parent chain now, graph-shaped API | **Yes** (state record shape) | Yes |
| `FC1-07` legacy handling | Read-only `unresolved`, no invented IDs; migration is a later task | No for Phase 1; yes for any migration | Yes |
| `FC1-14b` branch predecessor | Defer; fork-as-new-conversation recorded as the likely later shape | **Partially** — needs an explicit Phase 1 exclusion | Yes |
| `FC1-14c` turn/message | One turn = one message in Phase 1, `turn_id` present | **Yes** | Yes |
| `FC1-PRIVACY-DELETE` | Receipt with **no conversation identifier**; no UI labels or counts | **Yes** for the receipt shape | Yes |

---

## 1. `FC1-08` — mandatory envelope fields

**1. Exact contract language.** FC-1 §3 defines the envelope as: `event_id`, `conversation_id`, `sequence`,
`predecessor_event_id`, `event_type`, `operation`, `affected_ids`, `resulting_state_id`, `idempotency_key`,
`durability`, `timestamp`, `provenance` — each with an obligation ("Stable unique identity", "Monotonic order within
the conversation", …). Separately, **ADR-001 D6** requires that *every* document carry `schemaVersion`, which FC-1's
list does not include.

**2. Observed evidence.** `ChatHistory` (`internal/backend/chat/types.go:72–83`) carries `id`, `sender`, `text`,
`role`, `timestamp`, `hasImage`, `characterId`. It has **no** sequence, event type, idempotency key, durability,
provenance, affected ids, resulting state or predecessor. `conversation_id` is **implicit** — it is the directory
name (`persistence.go:472`); the record never states which conversation it belongs to. `ChatMeta`
(`types.go:50–58`) carries `id`, `name`, `createdAt`, `updatedAt`, world/scenario/character/lorebook names.

**3. The unresolved decision.** Which of the twelve fields are mandatory in Phase 1.

**4. Alternatives.**

| Option | Content | Consequence |
|---|---|---|
| M1 minimal | `event_id`, `sequence`, `event_type`, `timestamp`, `idempotency_key` | Smallest honest set, but no format version and no actor |
| M2 full FC1-08 set now | All twelve, populated from day one | Fields with no Phase 1 producer (`resulting_state_id`, `affected_ids`) become decoration that no test can verify; FC1-14b/14c not decided |
| **M3 core + version** | Mandatory: **`v`** (format version), `event_id`, `conversation_id`, `sequence`, `event_type`, `timestamp`, `idempotency_key`, `actor`, `durability`. Optional/deferred: `predecessor_event_id`, `affected_ids`, `resulting_state_id`, `operation` detail | Every mandatory field has a Phase 1 verifier; deferred fields arrive with the decisions that need them |

**5. Recommendation: M3.** Rationale: a mandatory field that nothing sets is worse than an absent one — it invites
readers to trust a value that was never produced. `v` is proposed **beyond FC-1's list** because ADR-001 D6 requires
it; that addition is flagged rather than smuggled in. `actor` is proposed because deletion and correction are
authority-bearing operations and the receipt has an authority field.

**6. Dependencies.** `FC1-11` (event type vocabulary) · `FC1-14c` (whether `turn_id` joins the envelope) · D-2
(sequence rule) · D-4 (idempotency semantics).

**7. Blocks freezing?** **Yes.**

---

## 2. `FC1-11` — Phase 1 event kinds

**1. Exact contract language.** FC-1 §3 lists at minimum: conversation created; user turn committed; assistant turn
committed; alternative/branch produced; current state selected; snapshot captured; memory copy committed (transfer
event); deletion recorded; recovery/reconciliation recorded — with the note that "names may change; the semantic
distinctions may not".

**2. Observed evidence.** Operations that exist today: append (`AppendHistoryEntry`, `persistence.go:668`) · whole-file
rewrite (`RewriteHistory`, `:635`) called from three places — final save (`app-chat.go:348`), the delete/edit path
(`app-prompt.go:389`), and clear-history (`app-prompt.go:406`) · save-as (`SaveChatAs`, `:437`) · rename
(`RenameChat`, `:719`) · delete (`DeleteChat`, `:471`).

**3. The unresolved decision.** The minimum event-kind set for Phase 1.

**4. Alternatives.**

| Option | Kinds | Consequence |
|---|---|---|
| K1 six | created, turn committed (role in payload), selection changed, conversation deleted, recovery recorded | Clean, but **cannot express edit, delete or clear** — so `RewriteHistory` cannot be retired, which blocks acceptance A-10 |
| **K2 eight** | `conversation_created` · `turn_committed` · `turn_corrected` · `turn_retracted` · `selection_changed` · `selection_resolved` · `conversation_deleted` · `recovery_recorded` | Covers every operation the app performs today; clear-history becomes one `turn_retracted` per turn, which keeps the ledger self-describing |
| K3 nine (full FC-1 list) | K2 + branch, snapshot, memory copy | Includes kinds with no Phase 1 producer; branch and memory are out of scope |

**5. Recommendation: K2.** The decisive argument is backwards: the branch must retire `RewriteHistory` (acceptance
A-10), and it cannot do that unless correction and retraction are expressible as events. Branch/snapshot/memory kinds
arrive with the work that produces them.

**6. Dependencies.** `FC1-08` · `FC1-09` (no in-place edit) · D-3 (selection events) · `FC1-24` (deletion receipt).

**7. Blocks freezing?** **Yes.**

---

## 3. `FC1-02` — graph scope now versus later

**1. Exact contract language.** FC1-02: "The materialized domain view is a rooted state/branch graph (D1 option B).
Option B is the representation chosen for the materialized view … The graph is reconstructible from the event ledger;
it is not a second source of truth."

**2. Observed evidence.** The current model is a flat `[]ChatHistory` (`types.go:88`) with no parent pointers, no
branch identity and no state view at all. No branching/regeneration code path was located in this pass —
**unestablished**, not absent.

**3. The unresolved decision.** Whether the materialized view is a full graph in Phase 1 or a linear chain that the
graph representation can later absorb.

**4. Alternatives.**

| Option | Content | Consequence |
|---|---|---|
| **G1 linear now, graph-shaped** | Each state has exactly one parent; `parent_state_id` present; traversal is a chain walk | A chain is a strict special case of the graph, so nothing is thrown away; no fork semantics are invented |
| G2 full graph now | Multi-parent states, forks, merges | Requires merge semantics that `FC1-14b` has not decided; risks a graph model nobody can validate |
| G3 no graph at all | Only a message list | Repeats today's model and cannot host `state_id` |

**5. Recommendation: G1.** Rationale: FC1-02's chosen representation is a graph, and a single-parent chain satisfies
it without deciding fork/merge. The record must state the restriction explicitly rather than implying the graph is
complete.

**6. Dependencies.** `FC1-14b` (what a branch is).

**7. Blocks freezing?** **Yes** — it fixes the shape of state records.

---

## 4. `FC1-07` — legacy handling

**1. Exact contract language.** FC1-07: existing records without stable IDs "remain unresolved until an explicit
migration decision. This contract does not invent IDs for legacy records or silently promote a display name, file
path or browser-local record to canonical identity."

**2. Observed evidence. This one is not theoretical.** Message identity today **is** the array position:

- `app-prompt.go:303` — `ID: fmt.Sprintf("%d", len(a.chatMessages)-1)`
- `app-prompt.go:357` — `ID: fmt.Sprintf("%d", i)`

So existing `id` values are positions, and any whole-file rewrite (`RewriteHistory`) renumbers them. That is a direct
violation of **FC1-06** ("Array index, display order … is never the identity"), and it means legacy records cannot be
promoted without inventing exactly the identity FC1-07 forbids. Separately, `CharacterID`'s own comment
(`types.go:79–82`) states: "Empty means ownership is not established (legacy or non-character messages)" — the code
already distinguishes legacy from unowned.

**3. The unresolved decision.** How existing records are represented in the ledger.

**4. Alternatives.**

| Option | Content | Consequence |
|---|---|---|
| L1 synthetic IDs at import | Generate identities for legacy rows | **Violates FC1-07**; manufactures history |
| **L2 read-only `unresolved` projection** | Legacy records are read, displayed and preserved, never assigned identity, never rewritten; the ledger has no events for them | Honours FC1-07 and FC1-06; costs a projection path that must not become a writer |
| L3 user-confirmed migration | Per-conversation, explicit consent, recorded | Honest, but needs a UI, a receipt and a decision that does not exist yet |

**5. Recommendation: L2 for Phase 1; L3 recorded as a later, separately authorised migration.** Positional `id`
values are never promoted, and the code path that generates them (`fmt.Sprintf("%d", …)`) is removed when the
in-memory model retires — not "reused as a seed".

**6. Dependencies.** None to freeze Phase 1. Migration depends on an owner decision that does not exist.

**7. Blocks freezing?** **No** for Phase 1, **provided** the contract states legacy is outside the ledger until L3.

---

## 5. `FC1-14b` — branch predecessor relationships

**1. Exact contract language.** FC1-14b: how a branch relates to its predecessor is undecided — "whether a branch is
a state sequence, whether a state may have more than one parent, and how a fork or merge is represented".
`predecessor_event_id` "records an ordered predecessor relation; it does not settle branch semantics".

**2. Observed evidence.** No branch structure exists (`types.go:88`); `predecessor_event_id` has no producer; no
regenerate/fork path was located — **unestablished**.

**3. The unresolved decision.** The representation of branch relationships.

**4. Alternatives.**

| Option | Content | Consequence |
|---|---|---|
| **B1 defer** | Phase 1 permits no branching; one parent per state; no branch events | Smallest commitment; the envelope makes no claim it cannot keep |
| B2 fork = new conversation | A fork creates a new `conversation_id` with its own ledger and a recorded origin reference | Simple, immutable, no merge; but identity semantics for "the same conversation forked" need the owner's view |
| B3 multi-parent graph | Full fork/merge in the ledger | Requires merge semantics nobody has decided; high risk |

**5. Recommendation: B1 for Phase 1, with B2 recorded as the likely later shape — not decided here.** Rationale: B3
cannot be specified without a decision FC-1 explicitly defers, and a Phase 1 that pretends to support branching would
put unverified semantics into the envelope.

**6. Dependencies.** `FC1-02` (state shape).

**7. Blocks freezing?** **Partially.** Freezing can proceed if the owner explicitly excludes branching from Phase 1
and the envelope is documented as single-predecessor.

---

## 6. `FC1-14c` — turn/message relationships

**1. Exact contract language.** FC1-14c: whether a turn contains one message or several, whether `message_id` is
turn-scoped, and how regeneration or reorder relates retained messages to their turns are undecided.

**2. Observed evidence.** There is **no turn concept** in the current model: one `ChatHistory` row per message
(`types.go:72–83`), no `turn_id` anywhere in the chat types. Regeneration as a retained alternative does not exist in
the model — **unestablished** in code.

**3. The unresolved decision.** The turn↔message relationship.

**4. Alternatives.**

| Option | Content | Consequence |
|---|---|---|
| **T1 one turn = one message** | `turn_id` present and equal to the committed message's identity; documented as a Phase 1 restriction | No renumbering later; multi-message turns deferred |
| T2 turn contains N messages | Turn identity distinct from message identity | Requires ordering inside a turn and a decision nobody has made |
| T3 no turn identity | Only messages | Leaves `turn_id` in FC1-05 unimplemented, and every later reorder becomes a migration |

**5. Recommendation: T1.** Rationale: it costs one field now and avoids a renumbering later. The restriction must be
written down, not implied.

**6. Dependencies.** `FC1-11` (event naming uses turns).

**7. Blocks freezing?** **Yes.**

---

## 7. `FC1-PRIVACY-DELETE` — receipt privacy and UI disclosure

**1. Exact contract language.** FC1-25: privacy deletion "is not a second deletion authority"; a privacy deletion of
authored content must be an explicit user deletion naming its set and recording a receipt (FC1-24); the policy
"remains unresolved … and must not be hidden inside `orphaned`". FC1-31 requires consequential transitions to record
evidence; FC1-32 requires privacy projection and deletion receipts to preserve the distinction between canonical
content, conversation copies and snapshots.

**2. Observed evidence.** `DeleteChat` (`persistence.go:471–479`) is `os.RemoveAll(chatDir)` — no receipt, no
enumeration of what was removed, no report of partial deletion, and it also removes `memories.json` inside the
directory (a lifecycle FC1-22e governs, not this branch). No privacy-specific deletion path was found. Conversation
directories are **timestamp-named** (`Chats/2026-05-27T16-46-11Z`), so conversation identifiers are enumerable.

**3. The unresolved decision.** What the deletion receipt may contain about the deleted conversation, and what the UI
may disclose.

**4. Alternatives.**

| Option | Receipt content | Consequence |
|---|---|---|
| P1 stored global salt | hash(id) + global salt retained | One salt links all receipts; enumeration of timestamp IDs defeats it |
| P2 stored per-operation salt | hash(id) + salt in the receipt | Scoped, but the salt sits beside the hash, so anyone with the receipt can test candidates; **pseudonym, not anonymity** |
| P3 raw conversation id | plaintext | Maximum linkability |
| **P4 no identifier** | operation id, timestamp, authority, counts, last-record hash | Proves a deletion occurred without any lookup path; per-conversation proof is unavailable |

**Revised after owner review:** the last-record hash is now **optional with a proposed default of absent**, because it
is a fingerprint that a surviving ledger copy can match. Exact fields: `TASK-03-LEDGER-SCHEMA-PROPOSAL-2026-11-22.md`
§5. Whether *any* identifier-bearing proof is wanted is a separate owner question (§6 of that document).

UI disclosure: **D1** nothing beyond "deleted"; **D2** a count of deleted conversations; **D3** titles/labels.

**5. Recommendation: P4 with D1.** Rationale, stated plainly: the earlier salted-hash proposal **does not survive
this analysis** — with timestamp-shaped ids and a stored salt, the hash functions as a verifiable reference to a
guessable identifier. P4 creates no persistent identity and no reconstruction path. If the owner requires
per-conversation proof, P2 is the fallback **with the enumeration risk documented**, and D2/D3 must wait for a
privacy decision.

**6. Dependencies.** Retention A (decided, §10) · FC1-24 (receipt fields) · a future salt-ownership decision only if
P2 is chosen.

**7. Blocks freezing?** **Yes** for the receipt shape; **no** for the ledger core.

---

## 8. What remains undecided after this document

Nothing here is adopted. Owner sign-off is required for: the `FC1-08` core field set (including `v` and `actor`) ·
the `FC1-11` eight kinds · linear-not-graph for Phase 1 · legacy read-only · branching excluded from Phase 1 with B2
noted as the likely later shape · one turn = one message · receipt P4 with D1.

**Still unresolved and out of scope for this package:** `FC1-14d` (selection identity) · `FC1-19a`
(`FC1-SNAPSHOT-CONTENT`) · `FC1-22c`'s memory-copy clauses · `FC1-STATE-RETENTION` (decided for conversation-owned
data; nothing else) · voice revision binding (`FC1-VOICE-VERSION`).

## 9. Authorization

Documentation and decision preparation only. No Go implementation, no call-site migration, no ledger creation, no
persistence adoption, no contract approval, no change to unrelated architecture or historical records.
