# FC-1 — Conversation State Contract

## Contract paper `FC-1-2026-11-21-01`

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Type | Contract paper — review draft |
| Authority | D1, D2, D3 and D7 owner decisions; `HF-CONTRACT-INTAKE-2026-11-20-01` |
| Gate | M1 |
| Status | **PROPOSED — not approved** |
| Implementation authorization | **None** |
| Tests/models/process runs | **0** |
| Revision | **2026-11-21 review revisions** — first: D1 attribution, D2 retention policy, transition semantics, acceptance recheck; second: state identity versus current selection, ordering language, manual deletion boundaries; third: selected-state deletion invariant, privacy-deletion boundary, lore and image acceptance coverage. FC-1 only; FC-2–FC-6 and the register are unchanged |

**Labels:** **Owner-approved** · **Proposed** · **Unresolved** · **Observed** · **Unestablished**.

**Identifier note:** Clauses suffixed with a letter (for example `FC1-14a`) were added in this revision. Existing identifiers were not renumbered, so cross-references from FC-2–FC-6 and the register remain valid; those documents have not been updated.

---

## 1. Contract decision

**FC1-01 (Owner-approved — D1).** Howl treats an append-only, ordered conversation-event ledger as the durable source of conversation state (D1 option C). State and branch structures are derived, materialized projections. A transition that produces new retained conversation state — committing a turn or producing an alternative, for example — creates a stable state ID for that retained state. A transition that selects among retained states, or that records an operation without changing conversation state, changes the current-selection projection or records evidence; it does not manufacture a new identity for retained material.

**FC1-02 (Proposed).** The materialized domain view is a rooted state/branch graph (D1 option B). Option B is the representation chosen for the materialized view, not the selected overall D1 architecture; the event ledger in FC1-01 remains the durable source of truth. The graph is reconstructible from the event ledger; it is not a second source of truth.

**FC1-03 (Owner-approved — D3 boundary).** The conversation ledger does not own authored world, scenario, character, persona or lore content. It records immutable references to canonical revisions and may retain reproducibility snapshots. Canonical authored objects remain under the library authority defined by FC-4.

**FC1-04 (Proposed).** Storage technology, file format, database choice and transport API are outside this paper. The contract defines identities, transitions, obligations and evidence, not an implementation schema.

---

## 2. Identity model

**FC1-05 (Proposed — required identities).**

| Identity | Meaning | Stability rule |
|---|---|---|
| `conversation_id` | Durable conversation root | Survives save/load, rename, export/import and restart |
| `event_id` | One immutable ledger event | Never reused; never derived from list position |
| `state_id` | One retained materialized conversation state | Stable across restart, branch switch, selection of that state and unrelated deletion; never reassigned to a different retained state |
| `message_id` | One retained message/turn artifact | Assigned at append; never reassigned by reorder |
| `turn_id` | One user or assistant turn | Stable local ordering key; not a display index |
| `snapshot_id` | One reproducibility artifact | Content-bound; immutable |

**FC1-06 (Proposed — no positional identity).** Array index, display order, filename order or message text is never the identity of a state, message or turn. Renumbering the UI does not renumber retained material.

**FC1-07 (Unresolved — legacy identity).** Existing records without stable IDs remain unresolved until an explicit migration decision. This contract does not invent IDs for legacy records or silently promote a display name, file path or browser-local record to canonical identity.

---

## 3. Durable event ledger

**FC1-08 (Proposed — event envelope).**

| Field | Obligation |
|---|---|
| `event_id` | Stable unique identity |
| `conversation_id` | Owning conversation |
| `sequence` | Monotonic order within the conversation |
| `predecessor_event_id` | Ordered predecessor relation, where applicable |
| `event_type` | Named transition kind |
| `operation` | What was requested or committed |
| `affected_ids` | States, messages, turns, snapshots or deletion set affected |
| `resulting_state_id` | Current state after the event, where applicable |
| `idempotency_key` | Duplicate-suppression key |
| `durability` | Requested / attempted / observed / verified / unknown |
| `timestamp` | Event time, distinguishable from commit evidence |
| `provenance` | Actor, capability, model or process origin where relevant |

**FC1-09 (Proposed — append-only rule).** A committed event is not edited in place. A correction, retraction or supersession is a new event that names what it corrects. Deleting a ledger event is not a normal operation.

**FC1-10 (Proposed — idempotency).** Repeating an event with the same `idempotency_key` and identical payload returns the existing result and produces no second effect. Reusing a key with a conflicting payload is an explicit error, not a silent overwrite.

**FC1-11 (Proposed — event kinds).** At minimum: conversation created; user turn committed; assistant turn committed; alternative/branch produced; current state selected; snapshot captured; deletion recorded; recovery/reconciliation recorded. Names may change; the semantic distinctions may not.

**FC1-12 (Unresolved — ordering implementation).** D1 establishes an ordered conversation-event ledger, and FC1-01 and FC1-08 carry that requirement; the existence of event ordering is not open. What remains unresolved is how ordering is implemented: concurrent append serialization, sequence assignment and gap policy, multi-process writer authority, and conflict handling. These must be resolved consistently with FC-4 and FC-5 before implementation.

---

## 4. Materialized state and transitions

**FC1-13 (Proposed — state view).** A materialized state view contains at least:

- `state_id`;
- parent state or states;
- selected turn/message references;
- active world/scenario/character/persona references;
- snapshot reference, if captured;
- the event that produced the state;
- validity status (`valid`, `degraded`, `unresolved`).

**FC1-13a (Proposed — current selection is a projection, not an identity).** The conversation's current selection is a `current_state_id` pointer held by the materialized projection and changed by selection events. Selecting or restoring a retained state changes that pointer and is recorded as an event; the retained state keeps its own identity. `current_state_id` must never silently point to a deleted or nonexistent state. If the currently selected state is manually deleted, the resulting selection is recorded explicitly and resolves either to a valid retained state or to an explicit empty or unresolved condition. Silent selection of an unrelated state, restoration of deleted material and manufacture of a historical state identity are prohibited. The exact transition and recovery mechanism for this case is unresolved (`FC1-SELECTION-AFTER-DELETION`).

**FC1-14 (Proposed — transition table).**

| Transition | Retained material | New identity | New `state_id` | Destructive? |
|---|---|---|---|---|
| Append user turn | Prior state remains | New event and state | Yes | No |
| Append assistant turn | Prior state remains | New event and state | Yes | No |
| Regenerate/alternative | Prior assistant material remains | New branch/state | Yes | No |
| Revert/restore selection | All retained states remain | Selection event; pointer changes | No — the retained state keeps its identity | No |
| Snapshot capture | State and canonical refs remain | New snapshot | No — artifact identity only, unless a distinct state is also committed | No |
| Explicit manual deletion of a named set | Only the named deletion set removed | Deletion event and receipt | Unresolved — depends on the named set; the current-selection pointer is re-resolved under FC1-13a | Yes, bounded to the named set |
| Ordinary chat switch | All retained material remains | No identity change | No — no conversation state change | No |

**FC1-14a (Proposed — identity is allotted only to new retained state).** A new `snapshot_id`, receipt, deletion record or event identity is not a new conversation state, and neither is a selection. Snapshot capture records an event and an artifact while leaving conversation state unchanged. Ordinary chat switching changes which conversation is displayed; it changes no conversation's state. Selecting or restoring a retained state changes the current-selection pointer and is recorded as an event, without creating a new `state_id`. Only transitions that produce new retained conversation state are marked as producing a `state_id` in FC1-14.

**FC1-14b (Unresolved — branch predecessor relationships).** How a branch relates to its predecessor is not decided: whether a branch is a state sequence, whether a state may have more than one parent, and how a fork or merge is represented. `predecessor_event_id` in FC1-08 records an ordered predecessor relation; it does not settle branch semantics.

**FC1-14c (Unresolved — turn/message relationships).** Whether a turn contains one message or several, whether `message_id` is turn-scoped, and how regeneration or reorder relates retained messages to their turns are not decided here.

**FC1-14d (Unresolved — selection occurrence identity).** Whether a selection operation or a selection occurrence needs an identity distinct from the selection event is not decided (`FC1-SELECTION-IDENTITY`). This paper allots no `state_id` to a selection and proposes no separate selection-identity scheme; the recorded event's `event_id` is the only identity it currently gives a selection occurrence.

**FC1-15 (Owner-approved — D2, retention and orphanhood).** Howl retains user data by default. Authored content is not automatically deleted because a conversation, world, scenario, character or other related object is deleted. Orphanhood describes a relationship or ownership condition; it does not, by itself, grant permission for automatic deletion of retained authored content. An earlier draft of this clause stated that deletion is permitted only when an object is orphaned; that shorthand is superseded here. Manual deletion from inside the application is the user's deletion authority for retained authored content, as FC1-22 and FC1-22a record.

**FC1-15a (Proposed — orphan status is not deletion authorization).** A missing reference, lack of current use, absence from the current composition, or loss of a parent relationship is not, by itself, authorization to delete an authored object.

**FC1-16 (Proposed — no silent deletion).** No automatic deletion path is authorized. A deletion must name its deletion set and record a receipt. A view being filtered, collapsed, archived or absent from the current screen does not make authored data disposable.

---

## 5. Canonical references and snapshots

**FC1-17 (Owner-approved — D3).** A conversation stores immutable references to canonical authored-object revisions plus reproducibility snapshots. It does not become the canonical owner of those objects.

**FC1-18 (Proposed — reference shape).** A reference identifies at least the object class, canonical object ID and immutable revision ID. A mutable “latest” pointer is not a reproducibility reference.

**FC1-19 (Proposed — snapshot obligations).** A snapshot records the exact canonical revision IDs used for composition, the state ID, the prompt-plan identity where applicable, and a content hash. It may retain display copies for reproducibility, but those copies do not acquire write authority.

**FC1-19a (Unresolved — snapshot completeness).** Which snapshot contents are required for reproducibility and which display copies are optional is not decided (`FC1-SNAPSHOT-CONTENT`). A snapshot must not claim reproducibility beyond what it can verify.

**FC1-20 (Proposed — restoration).** Restoration resolves the snapshot references against the canonical library:

- all revisions present → `valid`;
- some revisions missing or changed → `degraded` with named missing references;
- identity cannot be resolved → `unresolved`.

There is no silent substitution of the latest revision, a same-name object or a browser-local copy.

**FC1-21 (Calculated).** This is the precise D1/D3 relationship: the ledger owns conversation history; the library owns authored revisions; the snapshot binds a state to exact revisions for reproducibility; none of the three may claim another's authority.

---

## 6. Retention and deletion

**FC1-22 (Owner-approved — D2, authored-object retention).** Howl retains user data by default. Authored library content is not automatically deleted because a conversation, world, scenario, character or other related object is deleted. In particular:

- worlds remain stored unless manually deleted inside Howl;
- scenarios remain stored unless manually deleted inside Howl;
- characters remain stored unless manually deleted inside Howl;
- lore remains stored unless manually deleted inside Howl;
- images remain stored unless manually deleted inside Howl;
- other authored library objects follow the same retention principle.

**FC1-22a (Owner-approved — D2, manual deletion authority).** Manual deletion from inside the application is the user's deletion authority for retained authored content. A deletion operation identifies what the user is deleting and does not silently expand into unrelated assets. Only the object or objects explicitly included in the user's deletion operation may be deleted.

**FC1-22b (Proposed — no cascading deletion of authored assets).** No cascading deletion of authored assets is authorized. Deleting a conversation, world, scenario, character or other parent object does not delete the authored objects related to it. Deleting a world does not implicitly delete its scenarios, characters, lore or images. Deleting a scenario does not implicitly delete the characters, lore or images related to it, and the same principle applies to the deletion of characters, lore, images and other authored objects. Deleting a conversation does not delete the authored objects it referenced.

**FC1-22c (Unresolved — conversation-state and event-history retention).** D2 establishes retention by default for authored data. It does not establish the retention policy for conversation events, materialized state history, snapshots or other chat-scoped records. This paper authorizes no automatic retention window for them and claims no broader retention rule than the owner established; the applicable policy is unresolved (`FC1-STATE-RETENTION`).

**FC1-23 (Unresolved — `FC1-D2-ORPHAN`, orphan definition).** What makes a world, scenario, character, persona, lore object, conversation or state “orphaned”? Deleting a chat may orphan chat-associated memories or other chat-scoped records; deleting a world, scenario or character may orphan records that depended on that object. Candidate definitions may include no reachable reference, no owning conversation, no current composition and no retained snapshot. None is selected here. Orphanhood is a condition, not a deletion authorization (FC1-15, FC1-15a).

**FC1-23a (Unresolved — orphaned records).** The precise classification, ownership and handling of chat-scoped memories and other orphaned records is not decided. This paper does not establish an automatic orphan-cleanup policy (`FC1-CHAT-SCOPED`).

**FC1-24 (Proposed — deletion receipt).** A deletion record names the deletion set, the authority, the reason, the affected references, the prior state and the resulting state. A deletion that cannot name its set is refused.

**FC1-25 (Proposed — privacy deletion).** A privacy deletion and an orphan cleanup are not interchangeable. D2's manual-deletion rule for authored content does not prohibit a separately authorized privacy-deletion mechanism. A privacy deletion must have its own explicit authority, scope and evidence requirements, and must not be conflated with orphan cleanup or ordinary deletion. Privacy deletion may require removing content that is still referenced; that policy remains unresolved (`FC1-PRIVACY-DELETE`) and must not be hidden inside `orphaned`.

---

## 7. D7 voice identity preservation

**FC1-26 (Owner-approved — D7).** Character voice assignment is persisted character configuration, not transient prompt context. A hard context flush must not erase the voice profile, voice selection or character identity needed for TTS.

**FC1-27 (Proposed).** Voice identity is referenced by stable character ID and a versioned voice-profile ID. Display name is not the key. A voice profile may be unavailable; that produces an explicit unavailable state, not a silent fallback to another character's voice.

**FC1-27a (Unresolved — voice-profile revision binding).** Whether and how a voice-profile revision is bound to a character revision is not decided (`FC1-VOICE-VERSION`; FC-3 input).

**FC1-28 (Proposed).** TTS invocation is tied to a committed assistant turn. A TTS failure or cancellation does not roll back the durable text turn. This paper does not implement TTS.

---

## 8. Recovery and error semantics

**FC1-29 (Proposed — recovery obligations).**

| Failure | Required behavior |
|---|---|
| Duplicate append | Idempotent result; one durable effect |
| Torn write / partial event | Event is uncommitted or recoverable; never half-applied silently |
| Crash during state projection | Rebuild the projection from the ledger; report reconciliation |
| Missing canonical revision | `degraded` with named missing reference; no latest fallback |
| Snapshot hash mismatch | `unresolved`; do not claim reproducibility |
| Deletion interrupted | Report partial/unknown; do not claim complete deletion |
| Selected state deleted | Selection is recorded explicitly and resolves to a valid retained state or an explicit empty/unresolved condition; no silent reselection, no restored material, no manufactured identity |
| Parent object deleted | Retained authored assets remain; any orphaned records remain and are reported as unresolved, not cleaned up |
| Legacy record | Remains unresolved; no synthetic identity |

**FC1-30 (Proposed).** Recovery is not a new source of truth. It reconciles durable records with observed runtime state and records what was recovered, what was inferred and what remains unknown.

---

## 9. Observability and privacy

**FC1-31 (Proposed).** Every consequential transition records: operation, affected IDs, resulting state, durability level, timestamp, provenance and evidence reference. “Saved,” “rendered” and “committed” are distinct claims.

**FC1-32 (Proposed).** Conversation state may contain private character context and user material. Privacy projection, export and deletion receipts must preserve the distinction between canonical content, conversation copies and snapshots. A snapshot is not automatically safe to expose merely because it is a copy.

---

## 10. Open questions

| ID | Question | Status |
|---|---|---|
| `FC1-D2-ORPHAN` | What is an orphaned object, and when, if ever, does orphan status authorize deletion? | **Unresolved; orphan status is not deletion authorization** |
| `FC1-STATE-RETENTION` | What is the retention policy for conversation events, state history, snapshots and other chat-scoped records? | Unresolved; D2 does not establish it |
| `FC1-CHAT-SCOPED` | How are chat-scoped memories and other orphaned records classified, owned and handled? | Unresolved; no automatic cleanup |
| `FC1-CONCURRENCY` | How are appends serialized, sequences assigned and gaps handled, and who holds multi-process writer authority? | Unresolved; ordering itself is required |
| `FC1-BRANCH-PREDECESSOR` | How are branch predecessor relationships represented? | Unresolved |
| `FC1-TURN-MESSAGE` | What is the relationship between a turn and its messages? | Unresolved |
| `FC1-SELECTION-IDENTITY` | Is an identity distinct from the selection event required for a selection operation or occurrence? | Unresolved; no scheme proposed |
| `FC1-SELECTION-AFTER-DELETION` | How is the current selection re-resolved when the selected state is deleted? | Unresolved; the invariant is fixed, the mechanism is not |
| `FC1-LEGACY` | How are existing records reconciled without invented IDs? | Unresolved |
| `FC1-SNAPSHOT-CONTENT` | Which snapshot contents are required versus optional display copies? | Proposed, not decided |
| `FC1-PRIVACY-DELETE` | How does privacy deletion interact with retained references? | Unresolved |
| `FC1-VOICE-VERSION` | How are voice-profile revisions bound to a character? | Unresolved; FC-3 input |

---

## 11. Acceptance design — proposed, not run

1. **Replay:** rebuilding materialized state from the ledger yields the same state/branch identities.
2. **Duplicate:** the same idempotency key produces one effect; a conflicting reuse is rejected.
3. **Branch:** regenerate/alternative retains prior material and produces a distinct state.
4. **Revert:** selecting a retained state changes the current pointer without deletion or rewrite.
5. **Restart:** all state/message/turn identities survive restart.
6. **Snapshot:** exact canonical revision refs restore; a missing revision produces `degraded`, not latest substitution.
7. **Manual deletion:** deleting only the objects explicitly included in the operation removes exactly those objects, with a receipt; related and unrelated authored assets remain.
8. **Parent deletion:** deleting a chat does not delete its worlds, scenarios, characters, lore or images; deleting a world, scenario or character does not delete the authored objects that depended on it — deleting a world leaves its scenarios, characters, lore and images. Deleting a lore object or an image does not delete other authored objects related to it, and deleting other authored objects does not delete lore or images. Resulting orphaned records remain present and are reported as unresolved.
9. **Voice:** character voice assignment survives hard flush, character switch and conversation restoration.
10. **Crash:** recovery reports verified, partial or unknown; no phantom committed state.
11. **Selected-state deletion:** deleting the currently selected state records the resulting selection explicitly; the selection resolves to a valid retained state or to an explicit empty or unresolved condition. No dangling `current_state_id`, silent fallback, restored deleted material or manufactured state identity is observed.

The complete deletion policy is not testable yet. `FC1-D2-ORPHAN`, `FC1-CHAT-SCOPED` and `FC1-STATE-RETENTION` are unresolved; tests 7 and 8 verify only the authorized manual-deletion boundary and the prohibition on cascading deletion. Test 11 verifies the FC1-13a invariant without deciding the unresolved transition and recovery mechanism.

---

## 12. Traceability and non-authorization

| Link | Artifact |
|---|---|
| Owner decisions | D1, D2, D3, D7 in `HF-CONTRACT-INTAKE-2026-11-20-01` |
| State requirements | `HF1-PLAN-2026-11-21-01` |
| Readiness | `HFR-CP-2026-11-21-01` |
| Retention and manual deletion | D2 in `HF-CONTRACT-INTAKE-2026-11-20-01` |
| Canonical ownership | FC-4 (this set) |
| Prompt composition | FC-2 (this set) |
| Multi-character voice | FC-3 (this set) |
| Evidence | FC-6 (this set) |

**FC1-33 (Observed — no authorization).** This paper is a review draft. It does not approve a schema, authorize implementation, select a storage technology, define `orphaned`, establish a conversation-state retention window, authorize cascading deletion, authorize a privacy-deletion mechanism, or change any protected record.

**End of FC-1.** The event ledger owns conversation history; the library owns authored revisions; snapshots bind them; authored content is retained until the user deletes it; orphanhood is a condition, not a deletion authorization.
