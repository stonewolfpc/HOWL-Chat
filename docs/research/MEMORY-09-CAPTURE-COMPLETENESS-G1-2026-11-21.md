# MEMORY-09 — G1 Capture Completeness Audit

## The app has the context; the memory service does not receive it

**Status:** **G1 audit at time of writing; the test-only fixture subsequently passed the contract (G1R-22).** The
audit itself built no fixture, schema, code, migration, model, corpus or L-2 amendment. It establishes the minimum
capture envelope and the gaps that the later fixture tested.

| Field | Value |
|---|---|
| Date | **2026-11-21** (local, CST) |
| Gate | **G1 — capture completeness** |
| Result | **not passed as-is**; requirements established, fixture pending |
| Current backend | no speaker, stable source character, world/scenario metadata or attribution |
| Current frontend | user-message call with empty history/identities; no App binding |
| Data already in app | `ChatHistory`, `ChatMeta`, `ChatData` entities |
| States kept separate | **capture-incomplete** ≠ **classification-unresolved** |
| L-2 / Decision A | **unchanged** |
| Implementation | **0** |
| Map | **118** |

**Labels:** **Observed** · **Source-derived** · **Calculated** · **Inferred** · **Unestablished** · **Narrowed**.

---

## 1. Gate definition

**CCG-01 (Observed — G1).** G1 establishes that the capture pipeline provides enough information for a sound class
determination, and that missing information is visible as **capture-incomplete**, not silently repaired by a
classifier.

**CCG-02 (Observed — G1 is not a classifier test).** G1 uses no class labels, no rules, no BGE and no model. It tests
only the capture envelope.

**CCG-03 (Calculated — G1 is a prerequisite).** A classification experiment on capture-incomplete examples measures
nothing except the classifier's ability to guess missing context. G1 must pass before G2.

---

## 2. Current capture inventory

### 2.1 The backend call

**CCG-04 (Observed).** `Service.ProcessMessageForEvents(message string, history []string, identities
[]CharacterIdentity)` receives:

| Input | Content |
|---|---|
| `message` | the text |
| `history` | plain strings, no roles, IDs or spans |
| `identities` | `CharacterIdentity{ID, Name}` only |

**CCG-05 (Observed — missing from the call).** The call does not receive speaker role, source message ID/span, stable
source character (except optional identity labels), world/scenario IDs, world/scenario descriptions, character
metadata, attribution or authored-asset status.

### 2.2 The frontend call

**CCG-06 (Observed).** The browser wrapper calls
`window.go.main.App.ProcessMessageForEvents(message, [], [])` after the user sends a message. It passes empty history
and identities, and it runs on the user's outgoing message, not the character's reply.

**CCG-07 (Observed — no binding).** `frontend/wailsjs/go/main/App.d.ts` has no `ProcessMessageForEvents` function. The
frontend check fails and falls through to its own keyword fallback. The backend path is not exposed in the inspected
build.

**CCG-08 (Observed — detector output).** `EventDetectionResult` carries `EventType`, `Confidence`, `Evidence` and
`Context`. It carries no source, role, attribution or class.

### 2.3 What the app already has

**CCG-09 (Observed — `ChatHistory`).**

| Field | Use |
|---|---|
| `ID` | source message identity |
| `Sender` | display label |
| `Text` | content |
| `Role` | `system` / `user` / `assistant` |
| `Timestamp` | time |
| `HasImage` | attachment flag |
| `CharacterID` | stable character identity for assistant messages; empty means not established |

**CCG-10 (Observed — `ChatMeta`).** `WorldID`, `ScenarioID`, `ActiveCharacters`, `RemovedCharacters`,
`AddedCharacters`, `ActiveLorebooks`, `RemovedLorebooks`, `AddedLorebooks`, names.

**CCG-11 (Observed — `ChatData`).** `World *WorldContext`, `Scenario *ScenarioContext`, `Characters
[]*CharacterContext`, `Meta`, `History`, `Lorebooks`, `Memories`.

**CCG-12 (Observed — entity metadata).**

| Entity | Metadata |
|---|---|
| `World` | ID, Name, Description, KeyFacts, References |
| `Scenario` | ID, Name, Description, KeyDetails, References |
| `Character` | ID, Name, Description, Personality, Backstory, Goals, Relationships, References, Age, Role, Background, Priority, Version |

### 2.4 What the memory types carry

**CCG-13 (Observed — candidate).** `MemoryCandidate` has ID, EventType, Title, Summary, Content, SourceMessage, Scope,
OwnerID, OwnerKind, Tags, Confidence, Status, CreatedAt, ReviewedAt. It has **no** `KnowledgeClass`, source message
ID, role, attribution, world/scenario ID or authored-asset reference.

**CCG-14 (Observed — record).** `MemoryRecord` has ID, Title, Summary, Content, Scope, OwnerID, OwnerKind, Tags,
Pinned, Status, Confidence, Source, CreatedAt, UpdatedAt. It has **no** class, source role, source character, turn or
attribution.

---

## 3. Minimum capture envelope

**CCG-15 (Calculated — the logical envelope).**

| Field | Required | Why |
|---|---|---|
| `source_message_id` / span | yes | provenance and replay |
| `speaker_role` | yes | user / assistant / narrator / system |
| `source_character_id` | yes, explicit nullable | stable character identity; never a name guess |
| `source_character_name` | optional label | display only; not ownership |
| `chat_id` | yes | isolation |
| `turn_id` / sequence | yes | local ordering |
| `event_type` | yes | detection signal |
| `evidence` / source text | yes | explicit rule evidence |
| `world_id`, `scenario_id` | yes | authored-context check |
| `attribution` | yes for experience/belief | observed / told / believed / asserted / authoritative |
| `authored_asset_ref` | nullable | ambient/restricted eligibility |
| `existing_class` | optional | legacy only; never trusted as new |
| `content` | yes | the memory text |

**CCG-16 (Calculated — source character).** A character message must use `ChatHistory.CharacterID`, not a display-name
match. An empty `CharacterID` is an explicit unresolved identity, not a license to infer from `Sender`.

**CCG-17 (Calculated — role).** The role must distinguish a character assistant message from a user message, a
narrator/system message and a legacy message. The current `history []string` cannot.

**CCG-18 (Calculated — attribution).** Without attribution, the classifier cannot distinguish an experience from a
belief. That is a **capture-incomplete** state for the experience/belief determination, not a classification
abstention.

**CCG-19 (Calculated — authored asset).** Without an authored-asset reference, the classifier must not assign
`ambient` or `restricted`. Auto-captured character speech is not an authored world/scenario definition (K-R8).

---

## 4. The two states must not collapse

**CCG-20 (Calculated — capture-incomplete).** The system lacks information needed for a sound determination. The
candidate may be retained for inspection, but it is not formed/classified. It must not become any class, and it must
not be reported as “unresolved class.”

**CCG-21 (Calculated — classification-unresolved).** Capture is complete; the classifier has enough inputs but cannot
confidently assign a class. The memory may be formed with an explicit unresolved class only if the owner later
authorizes that representation. The unresolved state is a classification outcome, not a missing-input outcome.

**CCG-22 (Calculated — distinct reporting).**

| State | Meaning | Allowed outcome |
|---|---|---|
| **capture-incomplete** | a required capture field is missing | no class; no formation/classification; visible for repair/inspection |
| **classification-unresolved** | capture complete; classifier uncertain | explicit unresolved class, if authorized; no default |
| **resolved** | capture complete; class determined | class value with provenance |

**CCG-23 (Calculated — the failure mode).** If capture-incomplete is reported as unresolved, a missing speaker can be
mistaken for an ambiguous class. That would corrupt both the experiment and the production boundary.

---

## 5. Fixture matrix (proposed, not built)

**CCG-24 (Calculated).**

| # | Case | Expected |
|---|---|---|
| F1 | character assistant message with stable `CharacterID` | source character ID preserved; role assistant |
| F2 | user message | source character explicitly none; role user |
| F3 | narrator/system message | role system/narrator; no character |
| F4 | assistant message with empty `CharacterID` | capture-incomplete; no name guess |
| F5 | two characters with the same display name | stable IDs distinguish; no ambiguous guess |
| F6 | character message without attribution | capture-incomplete for experience/belief; other uses may proceed |
| F7 | authored world/scenario definition | authored asset reference present; ambient/restricted eligible |
| F8 | auto-captured character speech about the world | not ambient; no authored asset reference |
| F9 | legacy record with no class/source | reported as unmigrated; never guessed |
| F10 | missing world/scenario ID | capture-incomplete; no default context |

**CCG-25 (Calculated — assertions).**

- Every character memory has a stable source character ID.
- Every message has a role.
- Missing attribution is visible as capture-incomplete for the specific determination.
- Missing authored-asset status is visible for ambient/restricted.
- `capture-incomplete` and `classification-unresolved` are different records/states.
- No classifier, no class label and no model are used in G1.

**CCG-26 (Calculated — no fixture code yet).** The fixture is not implemented. Building it requires a test harness
and, if it exercises the production capture API, a change to that API. That is not authorized by this audit.

---

## 6. Current G1 result

**CCG-27 (Observed — not passed).** The current backend call cannot carry the minimum envelope; the frontend passes
empty history/identities and is not exposed as an App binding; the detector has no role/source/attribution. The app
already holds the needed data, but it is not passed to the memory service.

**CCG-28 (Calculated — the encouraging part).** The needed data largely exists: `ChatHistory` has role and
`CharacterID`; `ChatMeta` has world/scenario IDs; `ChatData` has world/scenario/character metadata. The gap is
**plumbing and an explicit envelope**, not a missing source of truth.

**CCG-29 (Calculated — G1 exit criteria).** G1 passes only when:

1. the fixture proves the envelope for F1–F10;
2. capture-incomplete is distinct from classification-unresolved;
3. no source character is inferred from a display name;
4. no class is assigned in G1;
5. the production API either carries the envelope or explicitly rejects formation as capture-incomplete.

**CCG-30 (Observed — G1 status).** **Audit complete; G1 not passed.** The next action is the fixture, not a classifier
run.

---

## 7. What this audit does not do

**CCG-31 (Observed).**

| Untouched | State |
|---|---|
| L-2 / Decision A | **unchanged** |
| P-1 / S-3 | closed; not reopened |
| Decision B | untouched |
| schema / code / migration / tests | **0** |
| fixture | **not built** |
| classifier / BGE / corpus | **0** |
| requirements / map | **0**; **118** |
| canonical capture | **0** |

---

## 8. Traceability

**CCG-32 (Observed).**

| What | Where |
|---|---|
| G1 gate | ACX-15/24/49 |
| Backend inputs | `Service.ProcessMessageForEvents`; CCG-04/05 |
| Frontend call | `memory-candidates-menu.js`; CCG-06 |
| No App binding | `App.d.ts`; CCG-07 |
| Chat history fields | `ChatHistory`; CCG-09 |
| Chat meta/data | `ChatMeta`, `ChatData`; CCG-10/11 |
| Entity metadata | `entities.World/Scenario/Character`; CCG-12 |
| Candidate/record fields | `types.MemoryCandidate/MemoryRecord`; CCG-13/14 |
| Capture-incomplete vs unresolved | ACX-03/20; CCG-20…23 |
| Ambient authored-only | K-R8; CCG-19 |
| Map | **118** |

---

## 9. Record status

**CCG-33 (Observed).**

| Check | Result |
|---|---|
| Canonical files edited | **0** |
| Schema, code, registry or ledger edited | **0** |
| Requirements changed | **0** — map **118** |
| G1 | **fixture contract passed; production wiring not established** (G1R-22) |
| Fixture built | **test-only harness; passed** |
| Classifier / BGE run | **0** |
| L-2 / Decision A amended | **0** |
| Tests executed | **0** |
| r33 · v9 · ledger · map | unchanged (`dfb28b60…` · `698254ca…` · `e5008020…` · `b4ba33cc…`) |
| This file | new, untracked; **G1 audit** — no canonical capture |

**End of audit.** The app already holds the context the memory service lacks: role and stable `CharacterID` in
`ChatHistory`, world/scenario IDs in `ChatMeta`, and world/scenario/character metadata in `ChatData` (CCG-09…12). The
backend call carries only message text, history strings and `{ID, Name}` identities; the frontend passes empty arrays
and has no App binding (CCG-04…08). The minimum envelope adds source span, role, stable source character, turn,
attribution and authored-asset status (CCG-15…19). **Capture-incomplete is not classification-unresolved** and must not
be collapsed (CCG-20…23). The current path cannot satisfy G1 as-is; the fixture is the next step, not a classifier run
(CCG-27…30). **No L-2 amendment, schema, code, fixture, model, corpus or test follows** (CCG-31/33).
