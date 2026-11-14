# Amendment 02 — conversation continuity, contextual media, scenario transitions, and multi-character information routing

**Date:** 2026-11-11  
**Status:** Architecture and research planning only. No implementation, migration, model, benchmark, runtime, network, or pairing authorization.  
**Amends:** `docs/architecture/INTEGRATION-AND-ACCEPTANCE-PLAN-2026-11-11.md`  
**Preserves:** the original plan, Amendment 01, all review documents, and all existing design documents  
**Original plan SHA-256:** `1e92a54434e8e3269015c47374ad57bd6ae755e64e27a07954c083ab952510cf`  
**Amendment 01 SHA-256:** `f3a0f8b9c5d19129db85d559f1880f21451cc7aa70596d37a408952c85c3f4ea`  
**Revision:** 7. Applies the OD-04 enforcement audit: adds a confirmation gate and decline path, conversation-scoped memory rule, chat-lifecycle isolation, auditable change records, and acceptance criteria `AC-C67`–`AC-C72`. Revision 6 is preserved as an archived artifact; revision 3 remains the independently reviewed identity.

## 1. Purpose and authority

This amendment makes the owner's conversation-continuity and information-routing requirements controlling for the affected Chat capabilities.

It does not replace Amendment 01. Amendment 01 remains authoritative for hardware intelligence, model residency, workload orchestration, readiness, release verification, and manual resource controls.

Amendment 02 adds requirements for:

- canonical conversation history independent of any model context;
- durable media and attachment identity;
- a protected initial conversation message;
- scenario, world, character, and lore transitions as real conversation events;
- contextual image generation from an assistant message;
- message deletion with memory and automatic-lore provenance;
- Public and Private information routing across multi-character Chat;
- persistent conversation continuity across world and scenario transitions, with distinct world, scenario, character, and chat state.

The central product rule is:

> The user controls the conversation, its information, its audience, and its history. Howl manages the complexity behind the scenes without hiding what it does.

### 1.1 Resolved owner decision `OD-01` — Go/Wails remains the technology foundation

**Status:** Decided by the owner on 2026-11-11; not an open question.  
**Scope:** Backend architecture, implementation strategy, migration, research framing, and any future platform proposal.

Go/Wails remains Howl's technology foundation. No language or backend rewrite is authorized or proposed on the basis of preference or speculative advantage. Reconsideration is permitted only if concrete engineering evidence establishes a significant limitation that cannot reasonably be resolved within Go/Wails.

AI-assisted maintainability and inspectability are explicit evaluation criteria, informed by the owner's ARES experience. That experience informs the criteria; it is not, by itself, evidence that a rewrite is warranted.

This decision prevents the amended architecture evaluation from being read as permission for a wholesale rewrite. It does not prohibit ordinary Go/Wails refactoring, module replacement, or interface changes that are supported by evidence and separately authorized.

### 1.2 Resolved owner decision `OD-02` — layered state and continuity

**Status:** Decided by the owner on 2026-11-11; not an open question.  
**Scope:** World and scenario state, character identity and memory, chat-specific state, participant continuity, and transition semantics.

Howl must support persistent, long-form conversations across world and scenario transitions without requiring users to create new chats to maintain continuity. World and scenario state, character identity and memory, and chat-specific state must remain distinct. Transitions must preserve conversation history, change active world/scenario state appropriately, and manage participant continuity according to who actually accompanies the chat. Character participation lifecycle events, including departure, return, death, and user-authorized ghost participation, belong to the individual chat and must not globally alter the character across other chats. Historical events remain distinguishable from active state. Memory must preserve relevant continuity for both the group and individual characters without inventing experiences for characters who did not participate.

The existing Step Out and Death mechanics are preserved and must be reconciled with this policy, not redesigned from scratch.

This decision establishes the layers and continuity principle. `OD-04` settles the explicit-change rule; the exact field taxonomy remains an R-09/DA-C4 design question.

### 1.3 Resolved owner decision `OD-03` — transition confirmation threshold

**Status:** Decided by the owner on 2026-11-11; not an open question.  
**Scope:** Automatic application, confirmation, undo, and recovery for transitions.

Transitions apply automatically, with a visible record, when ownership is unambiguous and no meaningful information is destroyed. This includes selecting the active world and scenario, loading declared defaults, updating the active participant list according to who follows, authorized departures and returns, and ordinary state changes that have an unambiguous owner.

Confirmation is required before removing a character from the chat entirely, replacing or discarding an explicit user change, changing or exposing private character information, permanently deleting or rewriting conversation history or associated records, resolving conflicting changes whose ownership is unclear, or making any change that would silently destroy meaningful continuity.

Undo is not a substitute for confirmation when an operation is destructive or privacy-sensitive. A failed transition must recover to a consistent state rather than leave a partially transitioned conversation. The field-level application of this threshold is a design matter under `R-09`/`DA-C4`, not an open owner question.

### 1.4 Resolved owner decision `OD-04` — explicit-change scope

**Status:** Decided by the owner on 2026-11-11; not an open question.  
**Scope:** A user's explicit change to a scenario-owned, scenario-defaulted, character, or chat field.

An explicit user change is a distinct operation. It must not be silently overwritten by a scenario reload or transition. Howl preserves it at the narrowest scope—the current conversation—unless the user explicitly selects a broader scope.

A character-persistent change and a scenario edit are separate, explicit operations. Howl must not infer either from the fact that the user typed something in Chat. The user must be able to see the scope of the change, promote it deliberately, and retain the historical record of what changed and why.

If a scenario transition would replace or discard an explicit conversation-level change, confirmation is required under `OD-03`. The change is not silently discarded merely because the selected scenario defines a different default.

## 2. Why this amendment is necessary

The current static implementation demonstrates several design assumptions that are not sufficient for the owner's clarified product:

- `ChatHistory` carries text, role, sender, timestamp, character ID, and a single `hasImage` boolean. It has no durable media identity or audience.
- Vision input is reduced to transient base64 and passed directly to one request.
- The visible history is rebuilt from text records; a `hasImage` flag does not restore or display the image.
- The backend memory DTO exposes `SourceMessages` but currently returns an empty list; the durable memory record has no source-message field.
- Message deletion removes a DOM node and rewrites history from the DOM. Message IDs are regenerated from array indexes.
- There is no protected initial-message entity. The only `FirstMessage` value found is a context timestamp.
- Scenario selection updates scenario metadata and can replace characters, but a scenario first message and a first-class transition record are not implemented. A system line is emitted only when characters are removed.
- Multi-character responses are orchestrated in frontend JavaScript. The backend receives a composed character-specific message, not a participant/audience routing decision.
- The contextual “Generate Image” action from an assistant message is absent.
- Image generation saves to the gallery, but no conversation/message association is recorded.

These are design gaps, not evidence that the requested behavior is impossible.

## 3. Three systems that must remain separate

| System | Question it answers | Owner |
|---|---|---|
| **Canonical conversation history** | What happened, what is visible, and in what order? | Conversation/event authority |
| **Model input** | What exact content did a particular model receive for a particular request? | Context assembler plus model-input receipt |
| **Persistent application state** | What durable memories, lore, scenario state, attachments, and derived information exist? | Their respective stores, linked by stable IDs |

The visible transcript must not be reconstructed solely from the current model context. Model context must not be assumed to contain everything visible. Persisted memories and lore must not be treated as substitutes for the transcript.

The canonical system must survive:

- application reload and restart;
- model changes;
- character changes;
- scenario and world changes;
- context-window truncation;
- summaries and compaction;
- attachment deletion or reintroduction;
- migration and rollback.

### 3.1 Amendment-01 cross-references

| Amendment-01 decision | Required use in Amendment-02 |
|---|---|
| `DA-H3` multi-GPU/backend capability | Chat, Image, voice, and private/public character runs must not assume one device or one shared memory pool. |
| `DA-H6` CPU/GPU placement | Each model-consuming workflow declares its supported placement; no hidden fallback. |
| `DA-H7` residency policy | A selected Chat model may remain selected while not resident; switching tabs does not unload it. |
| `DA-H8` scheduling | Public and private character runs, image jobs, and other consumers queue under the same scheduling authority. |
| `DA-H9` readiness verification | Image generation dispatches only after the image model is verified ready for that workload. |
| `DA-H10` release verification | Chat-to-Image and Image-to-Chat transitions distinguish verified, partial, and unknown release. |
| `DA-H11` overrides | User placement, residency, and scheduling choices remain visible and are never silently replaced. |
| `DA-I9` unified residency/scheduling contract | All model-consuming Chat and Image operations use one authority, not modality-specific managers. |
| `DA-I10` model-state persistence | The selected model remains selected across unload, tab change, and restart. |
| `DA-I11` manual-control audit | Private/public model runs and Chat image generation record the effective resource policy. |

No requirement in this amendment creates a second image loader, a second resource registry, or a tab-owned transition policy.

`OD-01` likewise constrains every orchestration, prompt, media, memory, RAG, and migration change described here: implementation remains within Go/Wails unless the owner reopens that decision with concrete engineering evidence.

## 4. Proposed canonical event model

The exact schema is not approved. The following shape is a design hypothesis for R-07 to test.

```text
Conversation
  id
  createdAt / updatedAt
  activeSceneStateRef
  initialAnchorEventId

ConversationEvent
  id                  stable, never reused
  conversationId
  sequence            monotonic within conversation
  kind                 message | transition | media | generation | system
  actor                user | character | system | tool
  body                 text and/or content blocks
  audience             public | private(recipientIds)
  sourceEventIds       provenance, including shared contributions
  modelInputReceiptIds requests that included this event/content
  createdAt
```

A message is one event kind, not the only event kind. Scenario transitions, generated images, media attachments, and protected anchors must not be forced into an ordinary assistant message with no type.

The visible transcript is a projection of canonical events. A model request is a separate projection assembled under an explicit policy and recorded as a receipt.

### 4.1 Required invariants

1. Every event has a stable identity that survives edits, reload, reorder, and rewrite.
2. Editing text does not silently change the identity or provenance of the event.
3. A model-input receipt identifies the exact events and media supplied.
4. A visible event has one canonical record; presentations are derived views.
5. No model input is hidden from the user's ability to inspect what was sent.
6. Deletion is an explicit operation against canonical records, not a DOM-only action.
7. Audience is a property of the information, not a UI label.

### 4.2 Proposed owner boundaries

| Concern | Proposed owner | Must not own |
|---|---|---|
| Conversation events and ordering | Conversation authority | Prompt policy, model memory |
| Media bytes and identity | Media/attachment store | Conversation order |
| Model-input receipts | Prompt/context assembler | Canonical transcript |
| Scene state and transition events | Scenario/conversation state authority | Character card content |
| Personal lorebooks | Canonical lorebook store | Automatic memory |
| Automatic lore | Automatic-lore store | Personal lorebook deletion |
| Memories | Memory service | Canonical transcript |
| RAG corpus and retrieval | RAG service | Audience policy |
| Prompt assembly | Context assembler | History ownership |

## 5. Durable conversation history and media

### 5.1 Visibility requirement

Everything sent to Chat must be visibly represented in conversation history, including:

- text;
- images;
- audio;
- vision-model image inputs;
- generated images;
- supported attachments;
- materially relevant scenario and conversation transition events.

The visible representation must survive:

- reload;
- close and reopen;
- tab changes;
- model changes;
- history navigation;
- reopening an old conversation;
- context truncation;
- summary generation.

### 5.2 Separate media lifecycle states

A single `hasImage` boolean is prohibited as the model for media. At minimum, the system must distinguish:

| State | Meaning |
|---|---|
| **Stored** | Durable bytes and identity exist. |
| **Linked** | The media is attached to a canonical event/message. |
| **Visible** | The user can see it in history. |
| **Previously supplied** | A model request previously included it. |
| **Currently included** | The present request includes it. |
| **Intentionally reintroduced** | The user explicitly asked to include it again. |
| **Deleted or unavailable** | The bytes or link are intentionally gone or cannot be resolved. |

A displayed attachment does not mean it is included in the next model request.

A vision image remains visible and persisted after the request that used it. Subsequent requests do not receive it unless the user intentionally reintroduces it.

A generated image remains visible without automatically entering subsequent model context.

### 5.3 Model-input receipts

Every model request must have a durable receipt containing, at minimum:

- request identity;
- conversation identity;
- model/backend identity;
- ordered event IDs included;
- ordered media IDs and content hashes;
- attachment reintroduction decisions;
- prompt-assembly version;
- exact payload hash;
- timestamp and outcome.

The receipt is how “everything supplied to a model is visibly represented” becomes testable. A text summary alone is not sufficient for multimodal inputs.

The receipt is diagnostic evidence, not a privacy bypass. A receipt that references private events, media, recipients, or derived content inherits the strictest audience of those references. It may be retained for owner inspection and audit, but it is never supplied to another character's context and its presence never authorizes resolving or reusing the referenced media.

### 5.4 Media presentation versus backend representation

The UI may render a thumbnail while the model receives a resized JPEG, a different content block, or a backend-specific encoding. That difference is permitted only if the same durable media identity links both representations and the model-input receipt records what was actually supplied.

The backend representation must never become the only evidence that the media existed.

### 5.5 Backend and frontend ownership

The current browser-first attachment flow is insufficient as the canonical store. The design must establish:

- which backend owns durable media;
- how media identities are created and deduplicated;
- how generated gallery assets relate to conversation media;
- how a reload resolves a media reference without a transient data URL;
- how deletion and retention work across conversation and gallery;
- how missing media is displayed as unavailable rather than silently disappearing.

`OD-01` fixes the implementation foundation for this ownership split: the durable owner remains the existing Go/Wails application, and frontend JavaScript must not become a second durable authority. This does not itself select a storage engine or schema.

## 6. Protected initial conversation message

Every conversation must have one required initial anchor.

Requirements:

- automatically displayed;
- always at the absolute top of the visible conversation;
- stored in exactly one canonical location;
- never duplicated into every scenario, summary, or model request;
- cannot be deleted;
- can be edited;
- remains protected after editing;
- remains present regardless of conversation length;
- remains navigable after thousands of messages;
- cannot be silently removed by summarization, compaction, model change, or context-window management.

The initial anchor is not the same as “the first ordinary message” and is not merely a timestamp. It is a canonical conversation event with a protected role.

### 6.1 Model-context behavior

The anchor may be omitted from a particular model request when the context policy requires it. It must still remain visible and canonical. If it is supplied to a model, the model-input receipt must record that fact.

### 6.2 Editing behavior

Editing the initial anchor must preserve:

- its identity;
- its protected status;
- its absolute-top position;
- its canonical storage;
- its provenance.

Whether an edited anchor changes earlier derived memories or summaries is an open decision (`O-06`).

### 6.3 Open schema questions

- Whether the anchor is a system, user, or special `anchor` event kind.
- Whether it can contain media.
- Whether it participates in normal transcript search.
- Whether editing it creates a revision event or replaces text in place.
- Whether a conversation can be forked without copying it.

These are decisions for `DA-C3` and `DA-C11`.

## 7. Scenario, world, character, and lore transitions

A selection change is an application and conversation-state transition, not merely a dropdown update.

Under `OD-02`, the transition model is layered rather than a contest over one shared state object:

- **World state** owns its scenarios, world rules, locations, and world-level lore.
- **Scenario state** identifies the selected scenario within its world; only that scenario is active in the chat.
- **Character state** owns identity, memories, experiences, and evolving character facts.
- **Chat state** owns this conversation's history, participants, active world/scenario selection, and chat-specific lifecycle events.

World or scenario selection changes the active context without resetting the conversation's identity or history. Character participation is governed by the chat, not by a global character-wide status.

### 7.1 Scenario selection transaction

Selecting a new scenario is a proposal until the confirmation gate resolves. As one coordinated operation:

1. Resolve the proposed selection and the affected world, scenario, character, participant, memory, and private state.
2. Determine whether `OD-03` requires confirmation.
3. If confirmation is required, present it before any active change. Decline leaves the active world/scenario, participant state, visible history, memory, and model context unchanged; it creates no canonical transition or first message. A rejected attempt may be recorded only as a non-active proposal.
4. On approval or when no confirmation is required, record the selection.
5. Persist the scenario identity and relevant revision.
6. Insert the scenario's first message into the conversation.
7. Display the first message in visible history.
8. Display the selected world, scenario, and active characters.
9. Record a transition event describing what changed.
10. Persist the transition and first message together, or report a recoverable partial failure.
11. Make the transition available to the next model interaction under the context policy.

No transition, first message, participant change, memory update, or model-context change may become active before the confirmation gate resolves. The model must not continue as if the previous scenario were still active.

### 7.2 Transition state

The transition record must carry, or reference:

- selected world and revision;
- selected scenario and revision;
- character additions, removals, and active state;
- relevant character cards;
- relevant personal lorebook entries;
- applicable automatic lore entries;
- scene state;
- the actor, reason, affected scope, previous value, resulting value, and confirmation outcome;
- reference to the first message.

### 7.3 First message and initial anchor are different

The scenario first message is a transition/opening message. The protected initial conversation message is the conversation anchor. They must not be conflated, duplicated, or allowed to overwrite each other.

Whether the scenario first message appears at the transition point or is inserted at the top is an open decision (`O-05`). It must not be silently decided by the implementation.

### 7.4 Transition versus ordinary message

Do not assume every transition is presented or transmitted identically. The event model must distinguish:

- scenario selection;
- world selection;
- character add;
- character remove;
- active-character switch;
- lore selection change;
- scene-state change;
- conversation fork or merge.

Each requires a defined visibility, ordering, and context-injection policy.

### 7.5 Active generation and queued work

Switching scenario while generation is active or messages are queued is an unresolved policy area (`O-05`). The design must prevent an old-scene generation from being silently attributed to the new scene, and must not silently discard work.

### 7.6 Participant continuity and lifecycle

A transition must distinguish characters who follow the conversation, characters who remain behind, and characters who join in the new context. A character who does not follow becomes inactive in the new scene but remains part of the chat's history; the group retains the departure event. Character memory must not invent personal experience of a world the character never entered.

Step Out, return, death, and user-authorized ghost participation are chat-scoped lifecycle states. They must not change the character's status in other chats. The existing Step Out and Death mechanics are preserved, but their data and behavior must be reconciled with the layered state and memory architecture.

### 7.7 Transition confirmation and recovery

`OD-03` governs confirmation. Routine, unambiguous, non-destructive changes apply automatically with a visible record. Removing a participant entirely, replacing or discarding an explicit user change, exposing private information, permanently deleting or rewriting history, resolving unclear ownership, or silently destroying continuity requires confirmation. Undo does not replace confirmation for destructive or privacy-sensitive operations. A partial failure must recover to a consistent state.

Explicit changes follow `OD-04`: preserve them at conversation scope, show that scope, require an explicit promotion to character or scenario scope, and never overwrite them silently during a transition.

## 8. Contextual image generation from Chat

This is a first-class Chat capability, not a shortcut to a separate inferior pipeline.

### 8.1 Required workflow

1. User opens the three-dot menu on an assistant message.
2. Selects **Generate Image**.
3. A flyover opens with the same relevant Image settings and controls.
4. Sensible defaults are provided.
5. User can adjust settings.
6. User presses Generate.
7. Howl queues the image workflow.
8. Amendment 01's unified orchestration engine performs any needed resource transition.
9. The image model is loaded and verified ready.
10. Generation begins.
11. The flyover closes on completion.
12. The result appears in the conversation and the Image board.

### 8.2 Shared implementation

The Chat action must use the existing Image subsystem's models, settings, backend, job management, and orchestration.

The flyover is a contextual interface to that implementation. It must not introduce:

- a second model loader;
- a second settings authority;
- a second gallery;
- a separate resource registry;
- a generation path that bypasses readiness verification.

### 8.3 Associations

A generated image must be associated with:

- the conversation;
- the originating assistant message;
- the generation operation;
- the durable media asset;
- the Image board/gallery entry.

The association must survive reload and navigation.

### 8.4 Progress contract

The progress display must cover the actual lifecycle:

- request queued;
- resource assessment;
- prior workload drain/unload if required;
- release verification;
- image-model load;
- readiness verification;
- generation dispatch;
- generation;
- completion or failure.

Generation progress must not begin before dispatch. A queued or loading operation must not look like sampling.

### 8.5 Prompt preparation

The owner has not confirmed whether the assistant message becomes the initial prompt automatically or another preparation method is preferred. This is `O-01`; it must not be silently assumed.

The design must support researched options such as:

- use the assistant message verbatim;
- use a selected portion;
- extract a scene description;
- combine a user-edited prompt with scene context;
- preserve the original assistant text unchanged.

### 8.6 Failure and cancellation

- A failed resource transition must not leave a completed-looking image.
- A cancelled job must be visibly cancelled or failed.
- A gallery save failure must be reported as a partial result, not success.
- A conversation association failure must be reported and recoverable.
- The originating assistant message must remain unchanged unless the user edits it.

## 9. Message deletion and persistent-information provenance

### 9.1 Three categories

1. **User-authored personal lorebook entries.**
2. **Automatic lore entries generated by Howl.**
3. **Automated memory entries generated by the memory system.**

Personal lorebook entries are never included in the message-associated deletion operation.

The required user-facing term for category 2 is **“automatic lore entries”** or **“automated system-generated lore entries.”**

### 9.2 Confirmation behavior

When a message has associated automated memories and/or automatic lore entries, the deletion flow must prompt the user.

Suggested wording:

> “This message contains associated automated memory entries and/or automatic lore entries (system-generated, not your personal lorebook entries). Would you also like to delete these associated entries?”

The prompt must identify which types are affected.

- **Confirm:** delete the message and the selected associated system-generated entries.
- **Cancel:** delete the message only; preserve its associated memory and automatic-lore entries.

Cancel does not cancel deletion of the message. This distinction must be explicit in the interface and in the backend operation.

### 9.3 Provenance model

The system must support:

- one message creating multiple entries;
- multiple messages contributing to one memory;
- multiple messages contributing to one automatic-lore entry;
- entries edited after creation;
- merged or summarized entries;
- entries referenced by other state;
- deletion of one source message;
- corrections and contradictory information.

Each derived entry requires durable source message/event IDs. A free-text source excerpt is not a provenance key.

A derived entry created from Chat input is conversation-scoped by default. It may enter character-persistent memory or automatic lore only through the explicit promotion required by `OD-04`; the entry records its source scope, prior scope, and promotion reason. Retrieval must enforce scope before merging the entry into any character context. A correction typed in one chat must not become a character-wide fact merely because the character appears in another chat.

### 9.4 Shared-source rule

Deleting one message must not destroy information still supported by other source messages.

The default behavior must preserve shared-source entries and remove only the deleted message's contribution. Whether an entry with no remaining sources is deleted automatically, offered for deletion, or retained is an open decision under `DA-C6`.

### 9.5 Stable message identity

The current delete path reconstructs history from the DOM and regenerates IDs from array positions. That is incompatible with reliable provenance.

The design must establish:

- stable message IDs;
- stable event sequence numbers;
- source references that survive edits and reorder;
- deletion receipts;
- tombstone or audit behavior where required;
- deterministic rewrite or append semantics.

### 9.6 Media and gallery deletion

Deleting a message with an image or generated image is not automatically the same as deleting the gallery asset. `O-07` must decide the relationship. The interface must not silently remove an image from the Image board because a chat message was deleted, nor silently leave a hidden attachment orphaned.

## 10. Multi-character Public and Private routing

Public and Private are backend-enforced information-routing modes, not presentation labels or alternate speaker selectors.

### 10.1 Public mode

Public is selected by a checkbox in multi-character Chat.

When enabled:

1. The user supplies content.
2. The model receives it.
3. The backend distributes the applicable original input to each current participant.
4. Each participant receives the actual input, not another character's description of it.
5. Each participant gets an opportunity to react.
6. Responses enter the shared conversation in the intended order.

The backend must construct the per-character context. A frontend concatenation of prior responses is not sufficient for the public contract.

### 10.2 Private mode

Private uses a recipient dropdown with two sections:

**Characters currently in this chat**

A visible divider.

**All other available characters, including nonparticipants**

Selecting a character means:

1. Only that character receives the original supplied information.
2. Howl loads the character card and necessary prompt context.
3. The character processes the information.
4. Its response is inserted into the shared Chat transcript.
5. Other participants can see and react to the response.
6. Other participants do not receive or see the original private input.
7. The character does not automatically become a group participant.

### 10.3 Strict privacy boundary

Private input must not leak through:

- shared transcript presentation;
- subsequent character prompts;
- conversation summaries;
- memory extraction;
- automatic lore generation;
- RAG indexing or retrieval;
- attachment references;
- generated descriptions or derived content;
- context reconstruction;
- scenario transitions;
- persistence, reload, or export.

The selected character's response is shared. The original private input is not shared with other characters.

Making a response public does not make its private source public.

### 10.4 Audience propagation

Every derived object must preserve or deliberately transform audience information:

| Derived object | Required privacy treatment |
|---|---|
| Memory | Retain source IDs and audience; do not broaden. |
| Automatic lore | Retain source IDs and audience; do not broaden. |
| Summary | Must not quote or paraphrase private content into a public summary. |
| RAG chunk/index | Private material is retrievable only under its audience. |
| Attachment reference | Must not expose bytes or a resolvable path to unauthorized characters. |
| Generated description | Must not reveal private source content. |
| Scenario transition | Must not backfill private content into shared scene state. |
| Model-input receipt | Retain for owner inspection and audit under the strictest audience of the referenced events, media, recipients, and derived content; never supply it to another character's context; a receipt never authorizes media resolution or reuse. |

### 10.5 Backend enforcement

A frontend checkbox/dropdown is a request, not an authority. The backend must know:

- the conversation participant list;
- the recipient or public audience;
- which character is being prompted;
- which events and attachments are permitted;
- which memories, lore, and RAG chunks are permitted;
- what response event may be shared.

A character must not receive private information merely because it exists in the same conversation file or because a generic history projection includes it.

## 11. Integration map

| Current seam | Current static behavior | Target | Transition / rollback | Dependencies |
|---|---|---|---|---|
| `ChatHistory` | Text + role + sender + timestamp + `hasImage` | Stable event/message identity, content blocks, audience, provenance | Add new canonical fields; keep text projection for old chats | R-07, R-11 |
| History persistence | `history.ndjson` append log, `history.json` fallback | Canonical event log with deterministic rewrite and migration | Dual-read; write new log; preserve old fallback | R-07 |
| Attachments | Base64 passed to one request; preview only | Durable media identity, link, visibility, reintroduction, receipt | Preserve immediate vision path; add media link before canonical cutover | R-07, R-08 |
| Initial message | No protected entity; `FirstMessage` is a timestamp | Protected anchor at absolute top, one canonical record | Introduce anchor at migration; preserve old first display | R-07 |
| Scenario selection | Metadata change; conditional system line | Atomic transition + first message + scene state + model context | Event recording behind old UI; no implicit bulk rewrite | R-09 |
| Contextual image | No Chat action | Menu → flyover → shared image job → Chat + board | Reuse Image job; add association and progress adapter | R-08, Amendment 01 |
| Generated image | Gallery save only; gallery error partial | Durable conversation/media/generation/board links | Add association without duplicating bytes | R-07, R-08 |
| Delete message | Immediate DOM removal; rewrite from DOM | Stable ID deletion + provenance prompt + receipts | Add prompt and IDs before switching writer | R-11, DA-M5 |
| Memory deletion | ID-based memory removal | Source-aware deletion preserving shared sources | Keep current ID delete as low-level primitive | R-11, R-01 |
| Automatic lore | Character-scoped files; no message provenance | Source IDs + audience + separate store | Introduce links before enabling bulk deletion | R-11 |
| Personal lore | Canonical lorebook store | Never included in message-associated deletion | No change to ownership | Existing lorebook cutover |
| Multi-character | Frontend sequential prompts, prior responses | Backend participant/audience routing and per-recipient context | Shadow audience resolution before enforcement | R-10, R-03 |
| Prompt assembly | App, backend service, frontend character path | One audience-aware assembly contract | Keep existing path until fixtures pass | R-10, DA-I1 |
| RAG | Scope-based retrieval; injected into generic prompt | Audience-aware retrieval and injection | Tag/deny private material before enforcement | R-10, R-02 |
| Image board | Gallery list | Same asset linked to conversation and operation | Add link table/record; no byte duplication | R-08 |
| Backend foundation | Go application and Wails bindings | Remain Go/Wails under `OD-01` | Ordinary Go/Wails refactoring only; any rewrite requires a separate evidence-backed owner decision | `OD-01` |

## 12. New research workstreams

### R-07 — Canonical conversation history and durable media

Questions:

- Event-based versus message-based persistence.
- Stable message/event identity across edits, deletes, reorder, reload.
- Durable attachment storage and content-addressed media.
- Media lifecycle and deduplication.
- Visible history versus model context.
- Model-input receipts and payload hashing.
- History reconstruction, migration, and recovery.
- Large-conversation navigation and indexing.
- Context-window independence.
- Generated-media associations.
- Attachment deletion, retention, and orphan handling.
- Multiple representations of one media asset.

Deliverable: a history/media architecture with invariants, migration risks, and acceptance fixtures.

### R-08 — Contextual image generation

Questions:

- Contextual generation workflows in comparable local applications.
- Reuse of existing image job/settings/backend interfaces.
- Flyover/modal design and accessibility.
- Prompt preparation options and user control.
- Defaults and preservation of last-used settings.
- Resource transitions and stage-accurate progress.
- Generation cancellation and partial failure.
- Conversation, message, media, operation, and board associations.
- Duplicate generation and idempotency.

Deliverable: a workflow and association design that reuses the existing Image pipeline.

### R-09 — Scenario and context transition architecture

Questions:

- World/scenario/character switching semantics.
- Layered ownership across world, scenario, character, and chat state.
- Explicit user changes, scenario defaults, and field-level conflict rules.
- Promotion of an explicit change between conversation, character, and scenario scope.
- Participant continuity when characters follow, remain, join, leave, die, return, or appear as authorized ghosts.
- Group history versus individual character memory across transitions, without invented experience.
- Transition event persistence and ordering.
- Context invalidation and reconstruction.
- Scenario first-message handling.
- Character-card selection and participant changes.
- Model awareness of transitions.
- Hot-swappable lore and scene state.
- Interaction with active generation and queued messages.
- Reload and partial-transition recovery.
- Fork/merge and historical navigation.

Deliverable: a transition event model, state machine, and recovery policy.

`R-09` may specify the field taxonomy, storage, and promotion mechanics, but it may not relax `OD-02`, `OD-03`, or `OD-04`. Any proposal that broadens change scope or weakens confirmation must be reported as a conflict requiring a new owner decision, not resolved by research or implementation.

### R-10 — Multi-character audience and private information isolation

Questions:

- Public group delivery.
- Private character-directed delivery.
- Per-message and per-input audiences.
- Character-specific prompt assembly.
- Nonparticipant character interactions.
- Information provenance and derived-data visibility.
- Memory, automatic-lore, summary, and RAG privacy.
- Private attachment handling.
- Shared responses from private interactions.
- Adversarial leakage tests.
- Backend enforcement boundaries.

Deliverable: an audience model, per-recipient context contract, and leakage test matrix.

### R-11 — Memory, automatic lore, and message-deletion provenance

Questions:

- Message-to-memory and message-to-automatic-lore associations.
- Shared-source and merged entries.
- Correction and contradiction handling.
- Deletion and retention policies.
- User inspection and correction.
- Separation from personal lorebooks.
- Stable IDs across edit/reorder/migration.
- Audit receipts and rollback.

Deliverable: a provenance and deletion design that never confuses personal lore with system-generated entries.

### 12.1 Dependencies

```text
R-07 canonical history/media
  ├── R-08 contextual image
  ├── R-09 scenario transitions
  ├── R-10 audience routing
  └── R-11 deletion provenance

R-10 audience routing
  ├── R-01 memory
  ├── R-02 RAG
  ├── R-03 orchestration (private/public model runs)
  └── prompt assembly (DA-I1)

R-08 contextual image
  └── R-03 orchestration + image subsystem

R-11 deletion provenance
  ├── R-01 memory
  ├── automatic-lore owner
  └── R-07 stable IDs
```

### 12.2 Research-question register

| ID | Research question |
|---|---|
| `RQ-07.1` | Message-based versus event-based canonical history; what is the minimum event taxonomy? |
| `RQ-07.2` | Stable event/message identity across edits, deletes, reorder, migration, and reload. |
| `RQ-07.3` | Durable media storage, content identity, deduplication, and lifecycle. |
| `RQ-07.4` | Model-input receipts sufficient to prove what each model received, including audience-scoped retention, non-disclosure to other characters, and non-authorizing media references. |
| `RQ-07.5` | History/media migration, legacy compatibility, and recovery. |
| `RQ-07.6` | Large-conversation navigation and indexing independent of model context. |
| `RQ-07.7` | Generated-media association with conversation, message, operation, and board. |
| `RQ-07.8` | Retention, deletion, orphan handling, and missing-media display. |
| `RQ-08.1` | Contextual image workflow and accessible flyover design. |
| `RQ-08.2` | Reuse of existing Image settings, jobs, backends, and gallery without a second path. |
| `RQ-08.3` | Prompt extraction, editing, and user control. |
| `RQ-08.4` | Stage-accurate progress across orchestration, load, readiness, generation, and failure. |
| `RQ-08.5` | Conversation/message/media/operation/board association and reuse. |
| `RQ-08.6` | Cancellation, partial failure, duplicate actions, and idempotency. |
| `RQ-09.1` | World/scenario/character switching semantics and visible state. |
| `RQ-09.2` | Scenario first-message placement, editing, and missing-message behavior. |
| `RQ-09.3` | Transition ordering, scene-state invalidation, and model awareness. |
| `RQ-09.4` | Active generation and queued-message interaction during transitions. |
| `RQ-09.5` | Reload, fork/merge, partial-transition recovery, and historical navigation. |
| `RQ-10.1` | Public delivery of original content to each current participant. |
| `RQ-10.2` | Private delivery to one selected character, including nonparticipants. |
| `RQ-10.3` | Per-recipient character-card and scene context for nonparticipants. |
| `RQ-10.4` | Audience propagation through memory, automatic lore, summaries, RAG, media, descriptions, and model-input receipts. |
| `RQ-10.5` | Backend enforcement boundaries and inspectable audience provenance. |
| `RQ-10.6` | Adversarial leakage tests across derived data and subsequent turns. |
| `RQ-11.1` | Durable message-to-memory and message-to-automatic-lore source links. |
| `RQ-11.2` | Multi-source, merged, summarized, and shared-contribution semantics. |
| `RQ-11.3` | Edits, corrections, contradictions, and source revision behavior. |
| `RQ-11.4` | Deletion, retention, recovery, and receipt semantics. |
| `RQ-11.5` | Separation of personal lorebooks from automated memory and automatic lore. |
| `RQ-11.6` | Auditability, rollback, and unsupported-source handling. |

Each research answer must carry an evidence class and source-quality record under Section 13. No question in this register has been researched yet.

## 13. Research standards and evidence

Each workstream must investigate:

- methods previously attempted;
- what worked and what failed;
- user complaints and praise;
- operational limitations;
- newer approaches;
- primary technical documentation;
- relevant papers;
- real implementation evidence;
- measured benchmarks where applicable;
- suitability for small and weak local models;
- resource requirements;
- privacy and persistence implications;
- integration complexity;
- failure and recovery behavior.

Evidence must be labeled:

- **Observed:** reproduced in a controlled environment with named conditions.
- **Source-derived:** stated by primary documentation or an identified implementation.
- **Calculated:** derived from stated assumptions and shown arithmetic.
- **Inferred:** reasoned from evidence but not directly demonstrated.
- **Unestablished:** no adequate evidence.
- **Superseded:** replaced by newer evidence or a later decision.

Marketing statements and theoretical designs are not proof of capability. Frontier-model results must not be assumed to transfer to small local models.

## 14. Decision register amendments

| ID | Decision | Required evidence | Blocks approval | Blocks implementation | Runtime grant |
|---|---|---|---|---|---|
| `DA-C1` | Canonical conversation/event schema and ownership | R-07 | Yes | Yes | `OA-09` if runtime |
| `DA-C2` | Durable media identity, storage, lifecycle, and reintroduction | R-07 | Yes | Yes | `OA-09` |
| `DA-C3` | Protected initial-anchor semantics | R-07; owner decision | Yes | Yes | `OA-09` |
| `DA-C4` | Scenario/world/character transition event model | R-09; `OD-02`; `OD-03`; `OD-04` | Yes | Yes | `OA-09` |
| `DA-C5` | Contextual image prompt preparation and association contract | R-08; owner decision `O-01` | Yes | Yes | `OA-09` |
| `DA-C6` | Message-associated automatic memory/lore deletion and shared-source policy | R-11; `DA-M5` | Yes | Yes | `OA-09` |
| `DA-C7` | Public/Private audience model and backend enforcement | R-10 | Yes | Yes | `OA-09` |
| `DA-C8` | Audience propagation through transformations (memory, lore, summary, RAG) | R-10, R-01, R-02 | Yes | Yes | `OA-09` |
| `DA-C9` | Nonparticipant private interaction and response visibility | R-10; owner decision | Yes | Yes | `OA-09` |
| `DA-C10` | History/media migration and rollback | R-07, R-11 | Yes | Yes | `OA-09` |
| `DA-C11` | Initial-anchor edit/revision semantics | R-07; owner decision | Yes | Yes | `OA-09` |
| `DA-C12` | Deletion, tombstone, retention, and gallery coupling | R-07, R-11; `DA-M5` | Yes | Yes | `OA-09` |

Existing related decisions remain open, including `DA-I1`, `DA-I3`, `DA-I5`, `DA-I6`, `DA-M5`, and the relevant RAG and orchestration decisions.

`OD-01` is settled, not a candidate. `DA-C1`, `DA-C10`, and the related migration decisions must be resolved within Go/Wails. A proposal to change the language or backend requires a separate owner decision supported by concrete evidence and is outside this amendment's scope.

## 15. Open owner decisions

| ID | Question | Why it must remain open |
|---|---|---|
| `O-01` | Does **Generate Image** automatically use the originating assistant message as the initial prompt? | The owner has not approved one prompt-preparation method. |
| `O-02` | How is public input delivered across different character/model configurations? | Per-character prompting, batching, and model switching can change delivery semantics. |
| `O-03` | How should private-derived memories and automatic lore behave when the same information is independently established elsewhere? | Audience and provenance can conflict; deletion must not destroy independent support. |
| `O-04` | How are private conversation events represented in shared visible history without exposing original contents? | User visibility, participant visibility, and model-input visibility are distinct. |
| `O-05` | How do scenario transitions interact with active generation and queued messages? | Ordering, cancellation, and attribution are unresolved. |
| `O-06` | How do edits to historical messages affect previously generated memories, automatic lore, and derived context? | Revision semantics are not established. |
| `O-07` | How does message deletion interact with images or attachments that also exist in the Image board? | Shared asset ownership and retention are unresolved. |
| `O-08` | What role and content model does the protected initial message have? | System/user/anchor semantics and media support are not established. |
| `O-09` | May a nonparticipant private character see relevant shared scene history? | The owner has defined the response behavior, not the full context boundary. |
| `O-10` | What should happen when a scenario has no first message? | A missing opening must not silently become an invented message. |

`O-11` is resolved by `OD-04` and is not an open owner question. The exact field taxonomy and promotion mechanics remain design work under `R-09` and `DA-C4`; they may not be used to narrow the owner policy in §1.4.

## 16. Acceptance criteria

All criteria are **unvalidated design targets**. They require appropriate evidence before any gate may pass. IDs are stable and must not be reused for different behavior.

### 16.1 Conversation history and media

| ID | Criterion |
|---|---|
| `AC-C01` | A vision image remains visible after application restart. |
| `AC-C02` | The image is not automatically resent to subsequent model requests. |
| `AC-C03` | Explicit reintroduction supplies it again and records why. |
| `AC-C04` | Stored, visible, and model-input identities are linked without collapsing them into one flag. |
| `AC-C05` | Generated images remain associated with the originating conversation and assistant message. |
| `AC-C06` | Generated images also appear in the Image board and are reusable without regeneration. |
| `AC-C07` | Every model input has a visible representation and a receipt. |
| `AC-C08` | The protected initial message remains at the absolute top and cannot be deleted. |
| `AC-C09` | Editing the initial message does not remove its protection or position. |
| `AC-C10` | Context truncation does not remove the visible initial anchor. |
| `AC-C11` | A missing media asset is shown as unavailable, not silently omitted. |

### 16.2 Contextual image generation

| ID | Criterion |
|---|---|
| `AC-C12` | **Generate Image** is available from the assistant message menu. |
| `AC-C13` | The flyover exposes the relevant Image settings. |
| `AC-C14` | The user can adjust settings before generating. |
| `AC-C15` | Resource transitions are visible and stage-accurate. |
| `AC-C16` | Generation does not begin before the image model is ready. |
| `AC-C17` | Progress distinguishes queued, unloading, loading, ready, dispatched, generating, and failed. |
| `AC-C18` | Completion registers the image in Chat and the Image board. |
| `AC-C19` | Failure and cancellation do not leave misleading partial states. |
| `AC-C20` | The originating assistant message is unchanged unless explicitly edited. |
| `AC-C21` | The workflow uses the same orchestration authority as the Image tab. |
| `AC-C22` | Duplicate or retried actions do not create contradictory associations or duplicate generations without consent. |

### 16.3 Scenario transitions

| ID | Criterion |
|---|---|
| `AC-C23` | After any required `OD-03` confirmation, selecting a new scenario records and displays the transition. A declined selection does not create a transition. |
| `AC-C24` | After the transition applies, its first message appears automatically under the approved placement policy. |
| `AC-C25` | The selected world, scenario, and characters are visible. |
| `AC-C26` | Only after the transition applies, the next model interaction receives the required transition context. |
| `AC-C27` | Reloading preserves the correct scenario state. |
| `AC-C28` | Old scenario state is not treated as current. |
| `AC-C29` | A transition cannot silently apply to an in-flight generation from the old scene. |
| `AC-C30` | Missing first-message behavior is explicit and does not invent content. |
| `AC-C31` | The initial conversation anchor and the scenario first message remain distinguishable. |

### 16.4 Memory and automatic-lore deletion

| ID | Criterion |
|---|---|
| `AC-C32` | Deleting a message with automated memory prompts the user. |
| `AC-C33` | Deleting a message with automatic lore prompts the user. |
| `AC-C34` | Wording identifies automatic lore as system-generated, not personal lorebook content. |
| `AC-C35` | Confirm removes the selected associated system-generated entries. |
| `AC-C36` | Cancel preserves those entries while still deleting the message. |
| `AC-C37` | Personal lorebook entries remain untouched. |
| `AC-C38` | Shared-source entries are not destroyed merely because one contributing message is deleted. |
| `AC-C39` | A deletion receipt records what was removed, retained, or left unresolved. |
| `AC-C40` | Stable IDs survive edit, rewrite, reload, and history navigation. |
| `AC-C41` | Deleting a message does not delete a gallery image without the approved policy. |

### 16.5 Public and Private routing

| ID | Criterion |
|---|---|
| `AC-C42` | Public input reaches every current participant. |
| `AC-C43` | Each participant receives the original applicable input, not a paraphrase. |
| `AC-C44` | Private input reaches only the selected character. |
| `AC-C45` | A nonparticipant can receive private input without automatically joining. |
| `AC-C46` | The selected character's response enters the shared conversation. |
| `AC-C47` | Other participants can react to that response. |
| `AC-C48` | Other participants do not receive the original private input. |
| `AC-C49` | Private data does not leak through memory, automatic lore, summaries, RAG, attachments, descriptions, transitions, model-input receipts, persistence, or subsequent context assembly. |
| `AC-C50` | The backend, not the frontend, decides what each character may receive. |
| `AC-C51` | A private response does not make its private source public. |
| `AC-C66` | A model-input receipt for a private request is visible only to its owner, is not supplied to another character's context, and does not by itself authorize media resolution or reuse. |

### 16.6 Adversarial and failure-oriented tests

| ID | Criterion |
|---|---|
| `AC-C52` | A public response that quotes private material does not grant other characters the private source. |
| `AC-C53` | A summary generated for public context omits private-only facts. |
| `AC-C54` | A memory derived from private input remains private when retrieved later. |
| `AC-C55` | RAG retrieval cannot cross an audience boundary. |
| `AC-C56` | A deleted source message does not erase a shared memory still supported elsewhere. |
| `AC-C57` | Reload reconstructs audiences and media links without widening access. |
| `AC-C58` | A failed scenario transition can be retried without duplicate first messages. |
| `AC-C59` | A failed image association can be repaired without regenerating the image. |
| `AC-C60` | A stale attachment cannot be reintroduced after intentional deletion. |
| `AC-C61` | A nonparticipant's private response cannot silently add them to the participant list. |

### 16.7 Required isolation scenarios

| ID | Scenario |
|---|---|
| `AC-C62` | Character D is privately sent an image; D receives the actual image and responds; D's response appears in shared history; Character A later receives D's response but not the image, a private description, or a summary revealing its contents; D does not become a participant merely from receiving the private input. |
| `AC-C63` | The audio variant of `AC-C62`: D receives the original audio, the response is shared, and A receives neither the audio nor a derived description. |
| `AC-C64` | The RAG variant of `AC-C62`: D receives the original retrieved material, the response is shared, and A cannot retrieve the private material or a derived summary. |
| `AC-C65` | A private-derived memory or automatic-lore entry is later retrieved in a public context; it remains inaccessible or is transformed only under an approved visibility change. |

### 16.8 OD-04 enforcement criteria

These criteria make the owner decisions testable at the acceptance gate. They are unvalidated design targets.

| ID | Criterion |
|---|---|
| `AC-C67` | A scenario transition does not silently replace, discard, or reset a conversation-local explicit change. If the transition would do so, it does not apply until confirmation; an approved transition preserves the change's recorded scope and value or explicitly supersedes it. |
| `AC-C68` | No transition, first message, participant change, memory update, or model-context change becomes active before required confirmation. A declined transition leaves the pre-transition state active and creates no canonical transition or first message. |
| `AC-C69` | The user can see whether an explicit change is conversation-local, character-persistent, or scenario-level; promotion to a broader scope requires a separate explicit action and is not inferred from Chat input. |
| `AC-C70` | A correction or derived memory originating from Chat is conversation-scoped by default, is not retrievable in another chat unless explicitly promoted, and records its source scope, prior scope, and promotion reason when promoted. |
| `AC-C71` | Death, Step Out, return, or ghost participation in one conversation leaves the character's state in every other conversation unchanged; the departed character's history remains visible in the originating conversation. |
| `AC-C72` | The durable change record includes actor, reason, scope, previous value, resulting value, transition identity, and confirmation outcome; it survives reload and remains retrievable. |

## 17. Conflicts discovered

| ID | Conflict | Consequence |
|---|---|---|
| `A2-01` | `ChatHistory` has only text and `hasImage`; no media identity. | Durable media and model-input provenance need a new canonical layer. |
| `A2-02` | Vision image input is transient base64; no persisted link. | “Image remains visible after restart” cannot be satisfied by the current record. |
| `A2-03` | The protected initial message is not modeled; `FirstMessage` is a timestamp. | The anchor must be introduced as a distinct event, not inferred. |
| `A2-04` | Delete rewrites history from the DOM and regenerates IDs. | Source-message provenance is unstable. |
| `A2-05` | No confirmation for associated memory/lore deletion. | Current delete behavior violates the required consent semantics. |
| `A2-06` | `MemoryRecord` has no source-message field; the DTO returns empty `SourceMessages`. | Memory-to-message deletion cannot be made reliable without provenance. |
| `A2-07` | Scenario selection is metadata-only in part; first message is absent. | Scenario changes are not a complete conversation transition. |
| `A2-08` | Scenario comment says scenarios never override chat state. | `OD-02` resolves the layering and continuity principle; `OD-03` resolves the confirmation threshold; `OD-04` resolves explicit-change scope. The remaining field taxonomy is design work under `R-09`/`DA-C4`. |
| `A2-09` | Multi-character orchestration lives in frontend JS with prior-response context. | Backend cannot enforce Public/Private routing or audience. |
| `A2-10` | Backend prompt is generic; frontend character-specific path composes a different message. | Per-recipient audience-aware assembly is not implemented. |
| `A2-11` | RAG has no per-message audience model; it retrieves by world/scenario/character scope. | Private RAG material can leak through generic retrieval. |
| `A2-12` | Contextual Generate Image is absent. | No current workflow or association exists. |
| `A2-13` | Generated images save to the gallery without a conversation/message link. | Chat and Image board are not synchronized. |
| `A2-14` | `hasImage` is used as a display/rebuild flag without image data. | Rebuilt history cannot show the image. |
| `A2-15` | Personal lore and automatic lore are separate stores, but automatic lore lacks message provenance. | Deletion cannot reliably distinguish shared and exclusive derived entries. |

## 18. Assumptions deliberately not adopted

The amendment does **not** assume:

- the assistant message is the correct image prompt;
- Public means “same model context for every character”;
- Private means only a UI label;
- the original private input should be hidden from the user;
- the original private input may be shown to other characters;
- every transition is an ordinary system message;
- the initial message is a system message;
- deleting a chat image deletes its gallery asset;
- deleting a gallery asset deletes the chat image;
- `hasImage` is sufficient to restore media;
- localStorage is the canonical history;
- summaries are automatically safe to share;
- memory or lore can be shared because its source is visible;
- RAG material is public because it is local;
- a private response can be used as a public description of its source;
- scenario changes may silently rewrite historical messages.

## 19. Sequencing and gates

```text
R-07 history/media
      │
      ├── R-08 contextual image
      ├── R-09 transitions
      ├── R-10 audience routing
      └── R-11 deletion provenance
      │
      ▼
DA-C1–DA-C12 + owner decisions O-01–O-10
      │
      ▼
Integration/acceptance design review
      │
      ▼
Owner architecture approval
      │
      ▼
Implementation authorization
```

Product staging remains:

- **Chat**, including character voice/TTS and Narration;
- **Image** next;
- future **Music** and **3D** studios;
- broader UI unification later.

The new conversation/media/audience work is cross-cutting and must be designed before the affected Chat behavior is implemented.

`OD-01` remains fixed throughout this sequence. No research workstream or migration phase may treat a language/backend rewrite as an available design option, a fallback, or a speculative optimization.

`§12.1` remains the dependency authority. The diagram is illustrative and does not override cross-workstream dependencies.

## 20. Governance and authorization

This amendment does not authorize:

- modifying application source code;
- modifying model files or runtime configuration;
- changing databases, memory stores, or lore stores;
- running migrations;
- loading models or running benchmarks;
- probing runtime behavior;
- network or Android pairing tests;
- implementing Public/Private routing;
- implementing contextual image generation;
- changing production conversation behavior;
- proposing or performing a wholesale Go/Wails language or backend rewrite.

OA-09 restrictions remain in force for any development-time runtime investigation.

This amendment does not approve the architecture, authorize implementation, or pass any acceptance criterion.

## 21. Migration, compatibility, rollback, and recovery

No migration is authorized. This section defines the required properties of any later migration proposal.

### 21.1 Backend foundation constraint

All migration and rollback work in this section is scoped to the existing Go/Wails application. `OD-01` forbids treating this migration as a language or backend transition. No migration plan may introduce a parallel runtime, a second backend, or a rewrite gate without a separate owner decision supported by concrete engineering evidence.

Ordinary Go/Wails refactoring, extraction of Go packages, and Wails binding changes remain possible only under the existing authorization gates. They are not a basis for claiming that the requested conversation, media, or audience behavior is implemented.

### 21.2 Migration phases

| ID | Phase | Required property |
|---|---|---|
| `MIG-01` | Inventory and shadow read | Read legacy history and new event records without changing visible behavior; report unknown IDs, missing media, and ambiguous ownership. |
| `MIG-02` | Additive identity | Assign durable IDs and media links without deleting or rewriting legacy records. |
| `MIG-03` | Initial anchor | Designate or create the anchor under the approved `DA-C3`/`DA-C11` policy; never invent a second copy. |
| `MIG-04` | Transition events | Add transition records prospectively; do not fabricate missing historical first messages. |
| `MIG-05` | Audience shadowing | Compute per-recipient permitted context in shadow mode; private material fails closed where audience is unknown. |
| `MIG-06` | Provenance and deletion | Attach source IDs to existing derived entries where demonstrable; mark unknown provenance and disable bulk associated deletion for those entries. |

### 21.3 Compatibility rules

- Older text-only chats remain displayable.
- `hasImage` may be preserved as a legacy hint, but it must not be interpreted as a durable media identity.
- A legacy image whose bytes are unavailable is shown as unavailable; it is not regenerated or described as if present.
- Legacy messages receive stable migration IDs where possible; array position is never the long-term identity.
- A legacy first message may be designated as the anchor only under `O-08`; no second copy is created.
- Historical scenario changes are not backfilled with invented first messages or transition events.
- Legacy audience is not inferred as public merely because no private marker exists.
- Existing personal lorebooks are never migrated into automatic lore or memory.
- Existing automatic lore and memories with unknown source messages are preserved and excluded from bulk associated deletion.
- Existing gallery images remain valid; they are not automatically linked to a conversation without evidence.

### 21.4 Rollback

A rollback must:

- disable new canonical writers before restoring the legacy projection;
- preserve the new event/media records for diagnosis rather than deleting them;
- preserve the legacy `history.ndjson`/`history.json` fallback;
- never delete personal lorebooks, original media, or the gallery as a rollback shortcut;
- report which new links and events are unreadable by the old path;
- be idempotent and leave no duplicate anchor or duplicate scenario first message.

### 21.5 Recovery

- A partial event/media write is quarantined and reported; it is not silently completed with guessed data.
- A missing audience on private material fails closed for character access.
- A missing media asset remains visible as unavailable and cannot be silently reintroduced.
- A failed transition is retryable without duplicating its first message.
- A failed image association is repairable without regenerating the image.
- A deleted message leaves an auditable tombstone or receipt as required by `DA-C12`.
- A restored backup is reconciled against media and provenance links before normal writes resume.

## 22. Source and evidence register

| ID | Source | Class | Use and limitation |
|---|---|---|---|
| `SRC-01` | Owner clarification in this request | Owner-defined | Controlling requirements; ambiguous interaction details remain open decisions. |
| `SRC-02` | Original integration plan | Source-derived / frozen | Prior planning baseline; preserved unchanged. |
| `SRC-03` | Amendment 01 | Source-derived / frozen | Hardware, residency, scheduling, and orchestration authority. |
| `SRC-04` | Existing design package | Proposed / source-derived | Design hypotheses and prior constraints; not evidence of implementation. |
| `SRC-05` | Independent review and re-review | Source-derived / review | Planning-status evidence, not architecture approval. |
| `SRC-06` | Static source inspection | Observed (static) | Seams cited in Section 17; no runtime behavior was exercised. |
| `SRC-07` | Existing tests and comments | Source-derived | Evidence of intended/covered behavior only; not proof that a requirement is met. |
| `SRC-08` | R-07–R-11 dossiers | Not produced | Research evidence is outstanding. |
| `SRC-09` | Owner decision `OD-01` on the Go/Wails foundation and ARES-informed maintainability/inspectability criteria | Owner-defined | Settled constraint; no implementation or separate evidence review is implied. |
| `SRC-10` | Muse review of Amendment 02 revision 3 | Source-derived / review | PASS for requirement capture only; advisory findings `M2-01` and `M2-02` are preserved; not architecture approval. |
| `SRC-11` | Owner direction on layered continuity, transition confirmation, and explicit-change scope (`OD-02`, `OD-03`, `OD-04`) | Owner-defined | Settles the layers, confirmation threshold, and explicit-change rule; the exact field taxonomy remains R-09/DA-C4 design work. |
| `SRC-12` | Muse addendum on OD-04 consistency at revision 6 | Source-derived / review | PASS WITH REQUIRED CHANGES; `M3-01`–`M3-05` require enforceable acceptance criteria and a confirmation gate; not architecture approval. |

No external product research, benchmark, user study, or runtime evidence was produced for this revision. Claims about comparable systems remain research questions, not findings.

## 23. Change log

| ID | Change | Status |
|---|---|---|
| `CH-01` | Captured canonical history, durable media, initial anchor, transitions, contextual image, deletion provenance, and Public/Private routing. | Added in Amendment 02 revision 1. |
| `CH-02` | Added explicit Amendment-01 cross-references for residency, scheduling, readiness, release, and shared orchestration. | Added in revision 2. |
| `CH-03` | Added stable research-question IDs `RQ-07.1`–`RQ-11.6`. | Added in revision 2. |
| `CH-04` | Added stable acceptance IDs `AC-C01`–`AC-C65`, including the D/A image, audio, and RAG isolation scenarios. | Added in revision 2. |
| `CH-05` | Added migration, compatibility, rollback, recovery, source/evidence, and change-log sections. | Added in revision 2. |
| `CH-06` | Added the independent-review handoff requirement. | Added in revision 2. |
| `CH-07` | Recorded resolved owner decision `OD-01` (Go/Wails foundation, no speculative rewrite, evidence and ARES-informed maintainability criteria) and cross-referenced it through backend, migration, change-log, and review planning. | Added in revision 3. |
| `CH-08` | Accepted Muse's Revision 3 PASS; closed `M2-01` on model-input-receipt privacy; recorded `M2-02` as `O-11` without deciding the scenario/chat authority conflict; accepted `M2-03` as informational. | Added in revision 4. |
| `CH-09` | Recorded `OD-02` layered state and continuity and `OD-03` transition confirmation; narrowed `O-11` to field ownership and the default scope of explicit user changes. | Added in revision 5. |
| `CH-10` | Recorded `OD-04` explicit-change scope, closing `O-11`: preserve at conversation scope, require explicit promotion to character or scenario scope, and never overwrite silently. | Added in revision 6. |
| `CH-11` | Applied the OD-04 consistency audit: added confirmation gating and decline behavior, conversation-scoped memory, chat-lifecycle isolation, auditable change records, and `AC-C67`–`AC-C72`; preserved revision 6 as an archive. | Added in revision 7. |

Deliberately unchanged:

- the original integration plan and Amendment 01;
- the owner's product staging;
- the distinction between selection, residency, and execution;
- the 10% VRAM reserve floor and Amendment-01 authorization boundaries;
- all open owner decisions `O-01`–`O-10`; `O-11` is resolved by `OD-04`;
- all conflict IDs `A2-01`–`A2-15`;
- the prohibition on implementation, runtime investigation, research execution, and acceptance claims.

## 24. Independent review gate

Revision 3 passed the review gate. Revisions 4–7 are controlled post-review deltas; no affected research workstream may begin without explicit owner authorization.

| ID | Review question |
|---|---|
| `RG-01` | Is every owner-defined requirement captured accurately? |
| `RG-02` | Has any requirement been weakened, silently reinterpreted, or converted into an implementation assumption? |
| `RG-03` | Is Public/Private isolation addressed across prompts, media, memory, automatic lore, summaries, RAG, persistence, and subsequent turns? |
| `RG-04` | Are transcript history and model input kept distinct and independently inspectable? |
| `RG-05` | Are personal lorebooks, automatic lore, and automated memory kept distinct? |
| `RG-06` | Are durable media identity, reuse, persistence, deletion, and Chat/Image-board association covered? |
| `RG-07` | Do Chat image generation and all multi-character runs use Amendment-01 orchestration? |
| `RG-08` | Is research scope adequate, evidence-led, and explicit about unestablished behavior? |
| `RG-09` | Are dependencies and acceptance criteria testable, including constrained systems? |
| `RG-10` | Are the original plan, Amendment 01, frozen documents, and prior digests intact? |
| `RG-11` | Is implementation or runtime authorization accidentally implied anywhere? |
| `RG-12` | Does the package support a clear PASS/FAIL review verdict without treating review as architecture approval? |
| `RG-13` | Does the package preserve `OD-01`—Go/Wails as the foundation, no speculative rewrite, and concrete evidence as the only reopening condition—without turning review into permission to reopen the language debate? |
| `RG-14` | Does revision 7 make `OD-04` enforceable through the confirmation gate, memory scope, lifecycle isolation, audit record, and `AC-C67`–`AC-C72`, without weakening `OD-02` or `OD-03`? |

### 24.1 Muse review disposition

Muse's PASS applies to revision 3 and is a requirement-capture result only. It is not architecture approval or research authorization.

| Finding | Status | Disposition |
|---|---|---|
| `M2-01` | Closed in revision 4 | Model-input receipts are audience-scoped in §5.3 and §10.4, named in `AC-C49`, given the dedicated criterion `AC-C66`, and carried into `RQ-07.4` and `RQ-10.4`. |
| `M2-02` | Closed by `OD-02`, `OD-03`, and `OD-04` | `OD-02` records the world/scenario/character/chat layering and continuity principle; `OD-03` records the confirmation threshold; `OD-04` records explicit-change scope and closes `O-11`. Field taxonomy is design work under `R-09`/`DA-C4`. |
| `M2-03` | Accepted as informational | §12.1 is recorded as the dependency authority; the §19 diagram is illustrative. |

### 24.2 M3 audit disposition

Muse's revision-6 addendum required enforceable acceptance coverage. Revision 7 applies the changes; the delta review request is `docs/review/AMENDMENT-02-REVISION-7-MUSE-DELTA-REVIEW-REQUEST-2026-11-11.md`.

| Finding | Status | Disposition |
|---|---|---|
| `M3-01` | Applied | `AC-C67`–`AC-C72` cover preservation, confirmation, scope visibility, promotion, memory scope, lifecycle isolation, and audit history. |
| `M3-02` | Applied | §7.1 now gates active changes behind confirmation and defines decline; `AC-C23` and `AC-C68` make the ordering testable. |
| `M3-03` | Applied | §9.3 makes Chat-derived entries conversation-scoped by default; `AC-C70` tests non-retrieval in other chats. |
| `M3-04` | Applied | `AC-C71` tests cross-chat lifecycle isolation and originating-chat history. |
| `M3-05` | Applied | §7.2 and `AC-C72` require actor, reason, scope, previous value, resulting value, transition, and confirmation outcome to survive reload. |
| `M3-06` | Accepted as informational | `AC-C66` remains a stable non-contiguous ID; no change is made to avoid ID reuse. |

The review request for the reviewed revision 3 remains `docs/review/AMENDMENT-02-MUSE-REVIEW-REQUEST-2026-11-11.md`. The PASS acceptance and finding disposition are recorded in `docs/review/AMENDMENT-02-MUSE-PASS-ACCEPTANCE-2026-11-11.md`. Research execution, implementation, migrations, and runtime work remain unauthorized pending owner authorization.

## 25. Verification record

The amendment's own digest is recorded in:

`docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md.sha256`

The companion verification record documents hash checks, preserved baselines, and the static source seams cited by this amendment. The review request records the package handed to Muse.

No source or runtime changes were made.
