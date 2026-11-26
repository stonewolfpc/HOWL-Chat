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
| Revision | **2026-11-21 review revisions** — first: D1 attribution, D2 retention policy, transition semantics, acceptance recheck; second: state identity versus current selection, ordering language, manual deletion boundaries; third: selected-state deletion invariant, privacy-deletion boundary, lore and image acceptance coverage; fourth: memory copy semantics (OD-D2-COPY), three data lifecycles, parent-deletion cleanup of memory originals, conversation-owned snapshots, transfer receipts and provenance. FC-1 only; FC-2–FC-6 and the register are unchanged |

**Labels:** **Owner-approved** · **Proposed** · **Unresolved** · **Observed** · **Unestablished**.

**Identifier note:** Clauses suffixed with a letter (for example `FC1-14a`, and `FC1-22d`–`FC1-22m` added by the memory-copy revision) were added in later revisions. Existing identifiers were not renumbered, so cross-references from FC-2–FC-6 and the register remain valid; those documents have not been updated.

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

**FC1-11 (Proposed — event kinds).** At minimum: conversation created; user turn committed; assistant turn committed; alternative/branch produced; current state selected; snapshot captured; memory copy committed (transfer event); deletion recorded; recovery/reconciliation recorded. Names may change; the semantic distinctions may not.

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

**FC1-15 (Owner-approved — D2, retention and orphanhood).** Howl retains user data by default. Authored content is not automatically deleted because a conversation, world, scenario, character or other related object is deleted. Orphanhood describes a relationship or ownership condition; it does not, by itself, grant permission for automatic deletion of retained authored content. An earlier draft of this clause stated that deletion is permitted only when an object is orphaned; that shorthand is superseded here. Manual deletion from inside the application is the user's deletion authority for retained authored content, as FC1-22 and FC1-22a record. This clause governs authored content only. Memory records and conversation reproducibility snapshots have their own directed lifecycles (FC1-22c, FC1-22d, FC1-22e) and are not covered by the authored-content retention rule; the retention-by-default statement here must not be read as retaining them indefinitely.

**FC1-15a (Proposed — orphan status is not deletion authorization).** A missing reference, lack of current use, absence from the current composition, or loss of a parent relationship is not, by itself, authorization to delete an authored object.

**FC1-16 (Proposed — no silent deletion).** No independent background deletion path is authorized: nothing is deleted because it is unreferenced, unused, absent from the current composition, archived, filtered out of a view, or discovered by a sweep. Deletion occurs only in three directed forms, and each names its deletion set and records a receipt (FC1-24):

- **explicit user deletion of a named set** — the named-set deletion mechanic that FC1-22a records for authored content. This clause extends that mechanic to memory records and copies as well as authored content; FC1-22a itself remains scoped to authored content and does not govern memory records. It operates on what the named set contains, including an individually selected memory copy owned by the scope the user is deleting from (FC1-22g, FC1-22j), and does not require that the copy's owning parent also be deleted;
- **parent-triggered cleanup of original memory records** (FC1-22e), which follows the explicit deletion of the parent that owns those originals;
- **conversation-triggered snapshot deletion** (FC1-22c), which follows the explicit deletion of the conversation that owns those snapshots.

A destination copy is therefore not limited to disappearing when a parent is deleted: the user can delete a selected copy on its own, leaving the original and any other copies in place. A view being filtered, collapsed, archived or absent from the current screen does not make authored data disposable.

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

**FC1-22c (Unresolved in part — snapshot and chat-scoped retention).** Conversation reproducibility snapshots are owned by the conversation that produced them and are deleted with that conversation; they are not retained as orphans when their owning conversation is deleted (FC1-22d). This lifecycle is directed by the owner. The retention policy for conversation events, materialized state history and other chat-scoped records remains unresolved (`FC1-STATE-RETENTION`); this paper authorizes no automatic retention window for them and claims no broader retention rule than the owner directed.

**FC1-22d (Proposed — three data lifecycles).** FC-1 distinguishes three data categories. No clause may conflate them, and no retention statement in this paper applies to more than one category.

| Category | Lifecycle direction |
|---|---|
| User-authored content | Deleted only through explicit user deletion (FC1-22, FC1-22a). Chat deletion and memory cleanup do not delete it. |
| Memory records | An original record belongs to its owning parent and is cleaned up when that parent is explicitly deleted (FC1-22e). Independently owned copies held elsewhere survive that deletion (FC1-22g, FC1-22i). |
| Conversation reproducibility snapshots | Owned by the conversation that produced them and deleted with that conversation (FC1-22c). They are not retained as orphans. |

**FC1-22e (Proposed — memory originals and parent-deletion cleanup).** A memory record's original belongs to the scope that owns it. When that owning parent is explicitly deleted, the parent's own original memory records are cleaned up as part of that deletion. This cleanup is a consequence of the deletion the user performed, not an independent cleanup path: it is bounded by the parent being deleted, names its deletion set and records a receipt (FC1-24). Whether a notification or confirmation must precede it, and whether copying is offered only at that point, are unresolved (`FC1-MEMORY-COPY-NOTIFY`, `FC1-MEMORY-COPY-TIMING`).

**FC1-22f (Proposed — transfer creates a copy).** The operation the source records call *transfer* is defined in this paper as creating a copy. A user-authorized transfer preserves the source memory record, which remains owned by its original parent, and creates an independent destination copy. Transfer is not an ownership move, and no clause of this paper may be read as moving ownership from a source to a destination. Earlier recommendations in this review adopted move semantics, under which the transfer would have removed the source record; OD-D2-COPY (2026-11-21) supersedes them. That recommendation is preserved here as superseded rather than erased, and it is not the current direction.

**FC1-22g (Proposed — independent ownership of copies).** A destination copy has its own owning scope and a retention lifecycle independent of the source. For a conversation destination, the owning scope is that conversation; whether a world destination becomes the owning scope of a copy is unresolved (`FC1-MEMORY-COPY-DESTINATION`), and the general statement that a copy has its own owning scope does not settle which scope a world destination is. Deleting one destination's copy does not delete the original or any other copy, and deleting the source parent does not delete copies already created elsewhere.

**FC1-22h (Proposed — provenance snapshot).** A destination copy records provenance at transfer time. At minimum it retains:

- the actor who authorized the copy;
- the source memory identity;
- the source chat identity where the source is a conversation;
- the source-parent identity;
- the destination identity;
- the source content hash;
- the transfer-event identity.

Every destination memory copy carries the literal `transferred_from` marker, and that marker identifies the copy's source through the copied origin references listed above. A destination copy that cannot carry the marker is not created. The marker and the provenance it identifies must survive deletion of the source record; like the provenance itself, the marker is a stored value and must not depend on a live source pointer.

The origin references are copied into the destination record and are never rewritten by the destination; a later edit to the copy does not alter what the copy records about where it came from. The provenance snapshot and the content hash must remain sufficient to identify the source after the source record is deleted. The provenance must not depend on a live pointer to the source record: a copied provenance reference is a stored value, not a live dependency, and the destination does not resolve, fetch or revalidate it against the source. Provenance establishes origin and transfer evidence; it does not establish that the original content was accurate, complete or still available.

**FC1-22i (Proposed — bounded deletion scope).** Deleting a parent deletes that parent's own original memory records and does not expand into independently owned copies held by other parents. The deletion set names the records being removed and refuses to extend to copies outside that parent's ownership. This is a required scope boundary, not an inference from FC1-22b: FC1-22b governs authored objects, and this clause governs memory copies.

**FC1-22j (Proposed — user-authorized selection and transfer-all).** Copying is user-authorized. The user may select individual memory records or transfer all of them, and each is an explicit operation rather than an automatic side effect. Where a copy may be sent — another conversation or a world — is part of the operation, and the destination is named in the transfer event. Switching the active world, scenario or character is not a transfer, a copy or a deletion of memories; it changes selection and composition only (FC1-14). The user interface for this operation is not specified here.

**FC1-22k (Proposed — copies do not carry knowledge relationships).** Copying a memory record does not silently copy its knowledge edges, relationships or other associated records. The destination copy is a record with provenance, not a reconstruction of the source's relational context. The Howl-side behaviour of K-R46 (no events or edges copied) and K-R48 (additivity) remains unresolved and is not settled by this clause; this paper claims no enforcement of either.

**FC1-22l (Proposed — transfer receipt and evidence).** A transfer records an event and a receipt naming the actor, the source chat identity where the source is a conversation, the source parent, the selected source memory identities, the destination identity, the created copies, and for each copy the provenance recorded and its content hash. The receipt records the `transferred_from` marker carried by each created destination copy, with that copy's source and destination identities consistent with FC1-22h. The receipt preserves the origin references as copied at transfer time: it does not rewrite them, and it does not depend on the source record still existing. The receipt is evidence that a copy was made and what it claims to derive from; it is not evidence that the source content was accurate. A transfer that cannot name its source records and its destination is refused. Duplicate suppression follows FC1-10. Destination-side merge and conflict behaviour remains unresolved (`FC1-MEMORY-COPY-MERGE`); this paper requires neither silent merge nor silent overwrite.

**FC1-22m (Proposed — Howl owns the workflow; MEMORY-09 boundary unchanged).** Howl owns the user-facing transfer workflow, the user's authorization, the transfer event, the destination record write, the parent lifecycle, the cleanup, and the associated notifications and receipts. MEMORY-09 exposes no transfer operation and no destination-write interface; that boundary (D-4, Option B) is unchanged and is not superseded by this paper. This clause does not define MEMORY-09's internal memory representation, edge behaviour, retrieval behaviour or integration, and FC-4 continues not to establish memory-record storage authority.

**FC1-23 (Unresolved — `FC1-D2-ORPHAN`, orphan definition).** Orphanhood is a relationship or ownership condition, not a deletion authorization (FC1-15, FC1-15a), and it is not applied indiscriminately across data categories. Authored objects are not orphan candidates under this paper: FC1-22 and FC1-22b keep them retained regardless of references. Ordinary memory records are not treated as orphan candidates either: an original is owned by its parent and is cleaned up with that parent's explicit deletion (FC1-22e); a copy belongs to the scope that owns its destination, which for a conversation destination is that conversation, while whether a world destination owns a copy remains unresolved (`FC1-MEMORY-COPY-DESTINATION`) and is not settled by the general destination-ownership wording in FC1-22g. Conversation reproducibility snapshots are not orphan candidates: they are conversation-owned and deleted with their conversation (FC1-22c). What remains unresolved is the status of records whose lifecycle is not established — legacy records, records of unknown origin, and records left by an interrupted lifecycle such as an interrupted deletion or an interrupted transfer. Whether `orphaned` is retained for those records, or retired in favour of a different classification, is not decided, and this paper does not declare the category empty. Candidate definitions may include no reachable reference, no owning parent, no current composition and no retained snapshot; none is selected here.

**FC1-23a (Unresolved — residual orphaned records).** The cleanup of a parent's own original memory records that follows an explicit parent deletion is directed (FC1-22e) and is not a general orphan-cleanup policy: it is bounded by the deletion the user performed and does not extend to independently owned copies (FC1-22i). The classification, ownership and handling of legacy, unknown-origin and interrupted-lifecycle records remains undecided (`FC1-CHAT-SCOPED`). This paper establishes no sweep, background pass or reference-scanning cleanup for any record type.

**FC1-24 (Proposed — deletion receipt).** A deletion record names the deletion set, the authority, the reason, the affected references, the prior state and the resulting state. A deletion that cannot name its set is refused.

**FC1-25 (Proposed — privacy deletion is not a second deletion authority).** A privacy deletion and an orphan cleanup are not interchangeable, and neither is a second authority for deleting authored content. Authored content is deleted only through explicit user deletion (FC1-22, FC1-22a); this paper authorizes no separately authorized mechanism that could delete authored content on its own initiative. A privacy deletion of authored content must therefore be an explicit user deletion that names its set and records a receipt (FC1-24), and it must not be conflated with orphan cleanup or ordinary deletion. Privacy deletion may require removing content that is still referenced; that policy remains unresolved (`FC1-PRIVACY-DELETE`) and must not be hidden inside `orphaned`.

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
| Parent object deleted | Retained authored objects remain (FC1-22, FC1-22b). The parent's own original memory records are cleaned up as part of that explicit deletion (FC1-22e). Independently owned copies held by other parents survive (FC1-22g, FC1-22i). A deleted conversation's own reproducibility snapshots are deleted with it (FC1-22c). Records of unknown or interrupted lifecycle remain and are reported as unresolved, not swept (FC1-23a) |
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
| `FC1-D2-ORPHAN` | Is `orphaned` retained for legacy, unknown-origin and interrupted-lifecycle records, or retired in favour of another classification? | **Unresolved; not decided, and the category is not declared empty** |
| `FC1-STATE-RETENTION` | What is the retention policy for conversation events, materialized state history and other chat-scoped records, excluding conversation-owned snapshots? | Unresolved; the snapshot lifecycle is directed, this policy is not |
| `FC1-CHAT-SCOPED` | How are legacy, unknown-origin and interrupted-lifecycle records classified, owned and handled? | Unresolved; directed lifecycles cover parent-owned originals and copies |
| `FC1-MEMORY-COPY-NOTIFY` | Is the pre-deletion notification blocking — requiring explicit confirmation or a choice before deletion proceeds — or advisory? | Unresolved; not decided |
| `FC1-MEMORY-COPY-DESTINATION` | Does attaching a copied memory to a world make the world the owning scope of that copy? | Unresolved; not decided |
| `FC1-MEMORY-COPY-TIMING` | Is copying available only as part of the deletion workflow, or can it be initiated at any time? | Unresolved; not decided |
| `FC1-MEMORY-NOTICE-SCOPE` | Does the memory-history deletion notice also cover conversation snapshots, or is snapshot deletion communicated separately? | Unresolved; not decided |
| `FC1-MEMORY-COPY-MERGE` | How are destination-side merge and conflict resolved when a copy arrives where an equivalent memory already exists? | Unresolved; neither silent merge nor silent overwrite is required or authorized |
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
8. **Parent deletion:** deleting a chat does not delete its worlds, scenarios, characters, lore or images; deleting a world, scenario or character does not delete the authored objects that depended on it — deleting a world leaves its scenarios, characters, lore and images. Deleting a lore object or an image does not delete other authored objects related to it, and deleting other authored objects does not delete lore or images. The parent's own original memory records are cleaned up with the explicit parent deletion, under a receipt that names that set. Independently owned copies held by other conversations or worlds remain, and their provenance still identifies the deleted source. A deleted conversation's own reproducibility snapshots are deleted with it. Records of unknown or interrupted lifecycle remain and are reported as unresolved, not swept.
9. **Voice:** character voice assignment survives hard flush, character switch and conversation restoration.
10. **Crash:** recovery reports verified, partial or unknown; no phantom committed state.
11. **Selected-state deletion:** deleting the currently selected state records the resulting selection explicitly; the selection resolves to a valid retained state or to an explicit empty or unresolved condition. No dangling `current_state_id`, silent fallback, restored deleted material or manufactured state identity is observed.

The complete deletion policy is not testable yet. `FC1-D2-ORPHAN`, `FC1-CHAT-SCOPED`, `FC1-STATE-RETENTION` and the `FC1-MEMORY-COPY-*` questions are unresolved; tests 7 and 8 verify only the authorized manual-deletion boundary, the prohibition on cascading deletion of authored objects, and the directed parent-deletion cleanup of a parent's own memory originals with the survival of independently owned copies. Test 11 verifies the FC1-13a invariant without deciding the unresolved transition and recovery mechanism.

---

## 12. Traceability and non-authorization

| Link | Artifact |
|---|---|
| Owner decisions | D1, D2, D3, D7 in `HF-CONTRACT-INTAKE-2026-11-20-01` |
| State requirements | `HF1-PLAN-2026-11-21-01` |
| Readiness | `HFR-CP-2026-11-21-01` |
| Retention and manual deletion | D2 in `HF-CONTRACT-INTAKE-2026-11-20-01` |
| Memory copy semantics | OD-D2-COPY (2026-11-21) |
| Memory transfer provenance and boundaries | `MEMORY-09-CONTRACT-NORMATIVE-2026-11-15.md` (K-R40, K-R45, K-R47); `MEMORY-09-ARCHITECTURE-v9-2026-11-17.md` (rules 31, 45); `MEMORY-09-OWNER-DECISIONS-2026-11-17.md` (Finding 1, D-1, D-4) |
| Canonical ownership | FC-4 (this set) |
| Prompt composition | FC-2 (this set) |
| Multi-character voice | FC-3 (this set) |
| Evidence | FC-6 (this set) |

**FC1-33 (Observed — no authorization).** This paper is a review draft. It does not approve a schema, authorize implementation, select a storage technology, define `orphaned` for legacy, unknown-origin or interrupted-lifecycle records, establish a conversation-event or state-history retention window, authorize cascading deletion of authored objects, authorize an independent background or reference-scanning cleanup, decide the notification, destination-ownership, timing or notice-scope questions, settle destination-side merge behaviour, establish memory-record storage authority, or authorize a privacy-deletion mechanism. It changes no protected record.

**End of FC-1.** The event ledger owns conversation history; the library owns authored revisions; snapshots bind them. Authored content is retained until the user deletes it; a conversation's own memory originals and reproducibility snapshots are cleaned up with its explicit deletion; an independently owned memory copy survives its source with recorded provenance; orphanhood is a condition, not a deletion authorization.
