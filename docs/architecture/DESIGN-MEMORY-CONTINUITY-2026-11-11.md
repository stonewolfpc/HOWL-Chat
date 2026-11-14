# Design — memory and conversational continuity

**Status:** Proposed architecture design; no implementation authorization  
**Date:** 2026-11-11  
**Owner inputs:** `UD-06` (automatic, transparent, editable memory), `DPR-06`  
**Plan sources:** `REC-CH-09`, `REC-CH-10`, `ODR-04`, `CWP-06`, `CWP-02`, `CWP-03`  
**Related design:** `DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md`, `DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`  

## 1. Decision in one sentence

Howl owns a revisioned canonical character/world state and an auditable, append-oriented history; the language model may only submit evidence-backed proposals, while automatic capture, transparent diffs, user editing, and reversible commits keep memory useful with small local models without allowing invented or silent state changes.

## 2. What exists today

| Surface | Current behavior | Consequence |
|---|---|---|
| Backend memory | Per-chat `Chats/<chatID>/memories.json`; candidate queue in `candidates.json`; CRUD and pin/status fields | Persistent records exist, but records are prose-centric and are not a canonical state model |
| Candidate review | Detection, classification, review/approve/reject, deduplication | Good foundation for automatic capture; approval is not a state-transaction model |
| Character ownership | Stable `character_id` ownership with fail-closed legacy handling | Preserve this rule; do not infer identity from display names |
| Frontend memory | `localStorage` mirror in `memory-store.js` plus backend bindings | Two stores can diverge; the frontend copy cannot be authoritative |
| Continuity | Export/import of selected memories across scenarios with confidence decay | Useful migration idea, but not a general state history or conflict model |
| Prompt injection | Memory/RAG path is not yet proven to be exactly-once or authoritative | The design must define one prompt assembly path, not add another |

The current memory system is a **capture and retrieval foundation**, not yet a continuity authority.

## 3. Design principles

1. **Howl is the authority.** The model is an unreliable narrator and extractor, never the owner of state.
2. **History is auditable and append-oriented; current state is a projection.** There is no silent overwrite or hidden mutation. Corrections, supersession, retraction, deletion, and privacy removal are explicit, attributable operations.
3. **No evidence, no automatic fact.** A model-proposed change requires a source turn and an evidence span.
4. **Editable does not mean silently mutable.** User edits are authoritative and visible; later model output cannot overwrite them.
5. **Automatic means automatic capture, not invisible mutation.** Every committed change produces a visible diff and an undo path.
6. **Small models get structure, not a bigger prompt.** Deterministic context assembly and compact state cards matter more than asking a small model to remember everything.
7. **Separation of concerns.** Canonical state, episodic memory, lore, RAG references, and chat transcript are different stores with different authority.
8. **Unknown is a value.** Missing clothing, location, or relationship state is not “none” and must not be invented.

## 4. Logical model

### 4.1 Canonical state

A character-state record is scoped by a **continuity scope**, not by display name:

```text
StateKey = (continuity_id, character_id)
```

`continuity_id` identifies the world/scenario/chat continuity in which the character is being portrayed. It prevents a character's state in one story from silently becoming canon in another.

```text
CharacterState
  continuity_id
  character_id
  revision
  fields:
    location:      StateValue<string | unknown>
    activity:      StateValue<string | unknown>
    clothing:      StateValue<ClothingItem[] | unknown>
    equipment:     StateValue<EquipmentItem[] | unknown>
    condition:     StateValue<string[] | unknown>
    relationships: RelationshipFact[]
    notes:         UserFact[]
  updated_at
  last_event_id
  last_editor: user | howl-auto | import
```

`StateValue` carries:

- `value`;
- `source_event_ids`;
- `confidence`;
- `asserted_at` / `effective_at`;
- `status`: current, disputed, superseded, user-locked;
- `qualifiers` (for example “wearing”, “carrying”, “at”, “doing”);
- `visibility`: normal, private, hidden.

Clothing and equipment are inventories, not prose. Each item has a stable identity, label, optional quantity, condition, and provenance. A location is a place reference, not a free-text sentence. Activity is a bounded state with a start event and an optional end event.

Relationships are separate directed or symmetric facts:

```text
RelationshipFact
  subject_character_id
  predicate            # ally, rival, parent, owes, trusts, loves, fears, ...
  object_character_id
  qualifiers
  source_event_ids
  confidence
  valid_from / valid_to
  status
```

Relationship state is not reduced to a single affinity number. A number may be a derived UI hint, never the canonical fact.

### 4.2 Event history

```text
Event
  event_id
  continuity_id
  chat_id
  turn_id
  source_message_id
  event_type
  participants[]
  summary
  evidence_span
  observed_at
  effective_at
  origin: user | model_proposal | import | system
  confidence
  status: active | corrected | superseded | retracted | forgotten
  supersedes_event_id?
```

History is **auditable and append-oriented**, not literally never-changing. “Append-oriented” means normal story changes are recorded as new attributable events, not by silently mutating earlier records. It does not mean a user can never correct, retract, delete, or forget.

| Operation | Required semantics |
|---|---|
| **Correction** | A new event identifies the earlier event, states the corrected value, and preserves the prior value for audit |
| **Supersession** | A new event replaces the current projection while the earlier event remains historically visible |
| **Retraction** | A new event marks an earlier assertion withdrawn; the projection stops using it but retains the audit trail |
| **Retcon / branch** | A new continuity branch starts from a declared state boundary; it does not rewrite the old branch |
| **Deletion / privacy removal** | Payload and index content are removed according to the deletion policy; a non-content tombstone and deletion receipt preserve referential integrity |
| **Forget / tombstone** | The fact disappears from current state, prompts, memory, and indexes; a tombstone prevents stale-cache resurrection |
| **Invalidation** | Derived memory, summaries, embeddings, and citations that depended on a changed or forgotten event are marked stale and cannot be silently reused |
| **As-of** | A query can reconstruct the projection at effective/recorded time without pretending the current state existed then |

Corrections and retractions are explicit and attributable. Deletion is a privacy operation, not permission for a silent background mutation. The deletion/forget policy is a blocking pre-implementation requirement, not merely a UI preference.

### 4.3 Memory items

A `MemoryItem` is a distilled, retrievable statement derived from one or more events. It is not canonical state.

```text
MemoryItem
  memory_id
  continuity_id
  scope: global | world | scenario | chat | character
  owner_id
  kind: episodic | preference | relationship | promise | location | item | ...
  title
  summary
  content
  source_event_ids[]
  valid_from / valid_to
  confidence
  salience
  pinned
  status: active | archived | forgotten
  user_edited
  version
```

A memory can support a scene without changing current state. A memory can also be the evidence that a state proposal is grounded.

### 4.4 State proposal

```text
StateProposal
  proposal_id
  continuity_id
  character_id
  base_revision
  changes[]:
    field
    operation: set | add | remove | clear
    value
    evidence_span
    confidence
  conflicts[]
  rationale
  status: pending | applied | rejected | superseded
```

The model writes proposals, not state. A deterministic application layer validates and commits them.

## 5. Automatic capture pipeline

```text
turn transcript
  → event detection (small model + deterministic cues)
  → entity resolution by stable ID
  → evidence extraction
  → proposal construction
  → validation and conflict detection
  → auto-commit or review queue
  → visible diff / undo record
  → state projection + memory index refresh
```

### 5.1 What may auto-commit

Automatic commits are permitted only when all of these are true:

- the statement is explicit in the user or character turn;
- the subject resolves to exactly one stable character ID;
- the field is objective and reversible (location, activity, clothing/equipment add/remove);
- the evidence span supports the exact value;
- no user lock or higher-authority value is being overwritten;
- the change does not create a contradiction with another current fact;
- the change is within the configured automation level.

Everything else enters the review queue:

- inferred emotions, motives, or relationship changes;
- ambiguous pronouns or multiple possible subjects;
- destructive changes, permanent death, identity changes;
- contradictions;
- low-confidence extraction;
- changes to user-authored or user-locked facts.

The default product behavior is **automatic capture with visible commits** for safe reversible facts, and **review before commit** for ambiguous or irreversible facts. This satisfies “automatic” without making the model an invisible author.

### 5.2 Conflict policy

There is no last-write-wins rule.

- Competing claims are retained as `disputed`.
- The current projection keeps the last **confirmed** value and displays the conflict.
- A user resolution creates a `user` event and locks the chosen value.
- A later model proposal cannot silently clear a user lock; it can only request review.
- Repeated contradictory evidence raises a visible continuity warning.

### 5.3 Idempotency and ordering

Every proposal binds to `base_revision`. Commit is compare-and-swap. A stale proposal is rejected or rebased through a new proposal, never applied by position. Turn and event IDs are idempotency keys so replay does not duplicate state.

## 6. History and “as of” queries

`GetStateAsOf(continuity_id, character_id, timestamp)` reconstructs the projection from events up to the requested time. It distinguishes:

- **effective time** — when the event happened in the story;
- **recorded time** — when Howl learned it;
- **current time** — what is canonical now.

This supports:

- “What was she wearing when they entered the city?”
- retcons without erasing the original event;
- flashbacks and timeline branches;
- audit of who changed what and why;
- visible continuity warnings when a draft contradicts a current fact.

A retcon is an explicit event, not a database overwrite. A user may also create a new continuity branch; the branch copies a declared state boundary, not the whole hidden history.

## 7. Retrieval and context assembly for small models

The prompt assembler requests a deterministic, budgeted context package:

1. **Always:** compact current state cards for active characters.
2. **Always:** the current scene/location and recent turn window.
3. **Then:** pinned facts and unresolved conflicts.
4. **Then:** relevant events/memories selected by stable entity, recency, lexical overlap, and optional embeddings.
5. **Then:** RAG reference passages, clearly labeled as reference material, never as canon.
6. **Never:** the full event log or a raw JSON database dump.

### 7.1 State card format

A small model should receive a short labeled block, not a schema exercise:

```text
[CANON — CURRENT]
Mara: location=North Gate; activity=waiting; clothing=grey cloak, boots;
equipment=iron key; condition=wet; relationship to Ivo=distrust (disputed)
[END CANON]

[MEMORY — RELEVANT]
- Mara promised Ivo she would not enter the tower alone. (source: turn 41)
[END MEMORY]
```

The state card is generated by Howl. The model is told that `CANON` is authoritative and that it must not invent changes outside a proposal block.

### 7.2 Retrieval strategy

- Lexical/metadata retrieval is the baseline and always available.
- Dense retrieval is optional and only local.
- Entity/time filters run before similarity search.
- Retrieve a small candidate set, deduplicate by event lineage, and rerank by salience and recency.
- Keep a token reserve for the response; never consume the whole context with memory.
- Include source references so the UI can answer “why does Howl think this?”

### 7.3 Small-model safeguards

- Extraction uses a compact schema and a limited vocabulary.
- The model may return “unknown”; it may not fill gaps.
- Summaries are generated only from cited events.
- A deterministic validator rejects unsupported fields.
- Memory retrieval has a hard cap; no silent spill of the entire history into the prompt.
- If the extraction model is unavailable, Howl captures the event and queues it without a state commit.

## 8. UI contract

The memory/continuity surface must show:

- **Timeline:** events with source turn and timestamp.
- **Current state:** field values, provenance, confidence, and “why”.
- **Diff:** before → after for every automatic commit.
- **Review queue:** proposed changes with approve/reject/edit.
- **Conflicts:** disputed facts and the current confirmed value.
- **Undo:** reverse any automatic commit with a compensating event.
- **Locks:** user-locked facts that model proposals cannot overwrite.
- **Forget:** remove payload/index content according to a defined privacy policy.
- **As-of view:** state at a selected point in the story.

“Editable” means the user can add, correct, merge, split, pin, archive, or forget. It does not mean the user must manually maintain the database.

## 9. Persistence and ownership

The design requires a single writer per `continuity_id`. The logical store is independent of the physical adapter so that the protected F3-C/D-STORAGE decision is not pre-empted.

- Current JSON files remain the compatibility source until an authorized migration.
- New state/event records require an explicit schema version and transaction boundary.
- Frontend `localStorage` becomes a cache/read model, never a competing source of truth.
- Backups and restore are part of the transaction contract.
- No migration is authorized by this design; `OA-01` and F3-C Phase 0/D-ROOT/D-STORAGE remain prerequisites.

### 9.1 Migration bridge from current memory

The current implementation cannot be promoted by changing the documentation. It needs an explicit bridge from per-chat JSON and renderer caches to the event/projection model.

**Current-data inventory required before any migration proposal:**

- `Chats/<chatID>/memories.json` — current memory records;
- `Chats/<chatID>/candidates.json` — pending/approved/rejected candidates;
- `Chats/<chatID>/continuity_config.json` — cross-scenario inheritance settings;
- legacy/name-owned character records and records with empty `owner_kind`;
- frontend `howl-chat-memories` and `howl-chat-candidates` caches;
- `memories_backup.json` and any prior backups, including the current placeholder behavior;
- import/export artifacts and character/world references.

**Observed current constraints that the bridge must handle:**

- edits overwrite records in place;
- deletes remove records without tombstones;
- save is read-modify-write with no transaction, lock, revision check, or atomic replace;
- the current backup path is not a reliable restore source;
- memory IDs use second-resolution timestamps, which can collide during rapid capture;
- the renderer cache is a second surface that can diverge;
- candidate, continuity, and memory records are not unified into revisions, conflicts, undo, or as-of state.

**Bridge requirements (design, not execution):**

1. **Inventory and classify** every existing record, including legacy ownership that must fail closed.
2. **Assign durable IDs** (collision-resistant, not second-resolution timestamps) and preserve the old ID as provenance.
3. **Map current records to imported events/memory** without inventing canonical facts. Ambiguous legacy records stay quarantined.
4. **Stage atomically**: snapshot, write to a staging store, verify counts/hashes, then commit; no partial cutover.
5. **Use a single writer and compare-and-swap revision** for every write.
6. **Preserve tombstones, deletion receipts, locks, conflicts, revisions, and attribution** through export/import; fail visibly rather than dropping them.
7. **Guard stale caches** so a forgotten fact cannot be re-ingested from `localStorage`, an old candidate queue, or an old backup.
8. **Provide rollback** to the prior JSON store until the new store is verified.
9. **Do not infer** a character owner from a display name, and do not promote legacy prose into canonical state automatically.
10. **Keep migration blocked** behind `OA-01`, F3-C Phase 0/D-ROOT, and D-STORAGE.

A migration that loses, duplicates, resurrects, or misattributes a memory is a failed migration, even if the new store starts successfully.

### Proposed interface boundary

```text
AppendEvent(event) -> event_id
ProposeChange(proposal) -> proposal_id
ApplyChange(proposal_id, expected_revision) -> new_revision
GetState(continuity_id, character_id) -> CharacterState
GetStateAsOf(continuity_id, character_id, time) -> CharacterState
ListEvents(filter) -> Event[]
RetrieveContext(request) -> ContextPackage
UndoChange(change_id) -> compensating_event
Forget(scope, policy) -> deletion_receipt
```

The prompt layer consumes `ContextPackage`; it does not read storage directly.

## 10. Failure behavior

| Failure | Required behavior |
|---|---|
| Extraction model absent or timed out | Store the event, mark extraction incomplete, do not change state |
| Model returns unsupported field | Reject proposal with visible reason |
| Ambiguous subject | Queue for review; do not guess by name |
| Conflicting current fact | Mark disputed; preserve both; show warning |
| Stale revision | Reject/rebased proposal; never overwrite by position |
| Disk full or write failure | Abort transaction; preserve prior state; visible error |
| User deletes/forgets a fact | Remove content/index per policy; block re-ingestion from stale cache |
| Imported state conflicts | Preserve imported provenance; require explicit merge; no silent promotion |
| Small model produces malformed JSON | Validate, discard malformed proposal, keep event |

## 11. Acceptance design — not run

1. **No silent overwrite:** a model proposal cannot change a user-locked value.
2. **Evidence requirement:** every automatic change has a source turn and span.
3. **Conflict:** contradictory statements produce a visible disputed state, not last-write-wins.
4. **As-of:** a later retcon does not alter the earlier as-of projection.
5. **Undo:** undo creates a compensating event and restores the prior projection.
6. **Small model:** a compact fixture can extract a bounded field without receiving the full log.
7. **Budget:** state, memory, and RAG sections never exceed their allocated context budget.
8. **Import/export:** state, events, memory, provenance, and locks round-trip or fail visibly.
9. **Deletion:** forgetting a fact removes it from prompt/index while preserving referential integrity.
10. **Single writer:** concurrent proposals at the same revision cannot both commit.

These are test designs only. No runtime test or model load was performed.

## 12. Open architecture decisions

| ID | Decision still required | Why it cannot be silently chosen |
|---|---|---|
| `DA-M1` | Physical store: transactional embedded database versus versioned files | F3-C/D-STORAGE protection; migration cost |
| `DA-M2` | Extraction model and structured-output contract | Small-model quality and hardware budget |
| `DA-M3` | Exact automation threshold by field | Product safety versus friction |
| `DA-M4` | Continuity branching/merge semantics | Prevents cross-story contamination |
| `DA-M5` | **Blocking pre-implementation:** user-facing deletion/forget policy, retention, tombstones, receipts, and export obligations | Privacy, stale-cache prevention, and export integrity |
| `DA-M6` | **Blocking pre-implementation:** conflict/automation policy by field | Prevents silent overwrite and inconsistent auto-commit |

## 13. Non-goals

- Giving the model write authority over canonical state.
- Replacing lorebooks with memory.
- Treating RAG reference text as fictional canon.
- Storing full chat history in every prompt.
- Defining the F3-C schema or authorizing migration.
- Implementing any code or running a model.
