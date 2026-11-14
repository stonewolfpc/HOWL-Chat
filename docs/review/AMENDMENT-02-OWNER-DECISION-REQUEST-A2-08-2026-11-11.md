# Owner decision record — `A2-08` / `O-11`

**Status:** Resolved by `OD-02`, `OD-03`, and `OD-04` on 2026-11-11.  
**Decision owner:** The owner of Howl.  
**Architecture approval:** Still not granted.  
**Not authorized by this record:** implementation, migration, runtime testing, or research execution.

The original filename is retained for continuity; this is now the decision record, not an open request.

## 1. The decision

World and scenario transitions must preserve the conversation's continuity. State is layered, not a single pile over which the scenario and chat compete:

| Layer | Owns |
|---|---|
| World state | Scenarios, world rules, locations, and world-level lore |
| Scenario state | The selected scenario within its world; only that scenario is active |
| Character state | Identity, memories, experiences, and evolving character facts |
| Chat state | This conversation's history, participants, active world/scenario, and chat-specific lifecycle |

A world or scenario switch is a transition, not a reset. Conversation history survives the transition. Characters may follow, remain behind, join, step out, die, return, or appear as user-authorized ghosts. Participation lifecycle belongs to the individual chat and does not globally alter the character in other chats.

An explicit user change is preserved at conversation scope unless the user explicitly promotes it to character or scenario scope. Howl must never infer the broader scope from the fact that the user typed something, and must never silently overwrite the explicit change during a scenario reload or transition.

## 2. `OD-02` — layered state and continuity

> Howl must support persistent, long-form conversations across world and scenario transitions without requiring users to create new chats to maintain continuity. World and scenario state, character identity and memory, and chat-specific state must remain distinct. Transitions must preserve conversation history, change active world/scenario state appropriately, and manage participant continuity according to who actually accompanies the chat. Character participation lifecycle events, including departure, return, death, and user-authorized ghost participation, belong to the individual chat and must not globally alter the character across other chats. Historical events remain distinguishable from active state. Memory must preserve relevant continuity for both the group and individual characters without inventing experiences for characters who did not participate.

The existing Step Out and Death mechanics are preserved and must be reconciled with this policy, not redesigned from scratch.

## 3. `OD-03` — transition confirmation threshold

**Automatic with a visible record:**

- selecting the active world and its scenario;
- loading declared scenario defaults;
- updating the active participant list according to who follows;
- authorized departures and returns;
- ordinary non-destructive changes with an unambiguous owner.

**Requires confirmation:**

- removing a character from the chat entirely;
- replacing or discarding an explicit user change;
- changing or exposing private character information;
- permanently deleting or rewriting conversation history or associated records;
- resolving conflicting changes whose ownership is unclear;
- any change that would silently destroy meaningful continuity.

Undo is not a substitute for confirmation when an operation is destructive or privacy-sensitive. A failed transition must recover to a consistent state rather than leave a partially transitioned conversation.

## 4. `OD-04` — explicit-change scope

An explicit user change:

- is a distinct operation, not a scenario default, migration, or background refresh;
- is preserved at the narrowest scope—the current conversation—unless the user explicitly selects a broader scope;
- may be promoted to a character-persistent change or a scenario edit only through an explicit user choice;
- must show its scope to the user;
- must retain a historical record of what changed and why;
- must not be silently overwritten when a scenario reloads or a transition applies a different default;
- requires confirmation under `OD-03` before a transition replaces or discards it.

The exact field taxonomy and the UI mechanics for promotion remain design work under `R-09` and `DA-C4`. They may not be used to narrow this policy.

## 5. Field ownership matrix

This matrix states the owning layer and the default treatment of an explicit change. It does not claim that every field has been enumerated.

| Field | Owning layer | Explicit-change treatment | Active/history behavior | Confirmation |
|---|---|---|---|---|
| Active world | Chat selects; world defines | Conversation-scoped selection; scenario definition unchanged | Transition record | Automatic selection |
| Active scenario | Chat selects; world owns definition | Conversation-scoped selection; broader edit requires explicit scenario operation | Transition record | Automatic selection |
| Active participant roster | Chat | Conversation-scoped follow/stay/join; global roster change requires explicit operation | Participation event | Confirmation on removal |
| Scene, location, and situation | Scenario default; chat may override | Conversation-scoped override; scenario edit requires explicit operation | Historical values remain | Confirmation if an override is replaced or discarded |
| Character identity and traits | Character | Conversation-scoped by default; character-persistent change requires explicit promotion | Character state | Confirmation on global change |
| Character memories and relationships | Character | Conversation-scoped by default; persistence across chats requires explicit promotion | Character memory | No invented experience |
| Step Out / return | Chat | Conversation-scoped lifecycle action | Chat lifecycle event | User-guided event |
| Death / ghost | Chat | Conversation-scoped lifecycle action | Chat lifecycle and graveyard | Explicit user authorization |
| Private audiences and material | Chat / character boundary | Conversation-scoped; must not broaden silently | Must not broaden | Confirmation on exposure |
| Derived memory, lore, summary, RAG | Source-scoped | Source-scoped; audience must be preserved | Historical/source-linked | Confirmation on destructive change |

## 6. What remains design work

The following are not open owner questions:

- the exact field taxonomy within each layer;
- the UI mechanics for promoting a conversation-local change to character or scenario scope;
- the data representation of historical versus active values;
- the precise transition and rollback protocol;
- the mapping of existing Step Out and Death data and behavior to the layered model.

These belong to `R-09`, `DA-C4`, the memory design, and the implementation-planning gates. They must remain consistent with `OD-02`, `OD-03`, and `OD-04`.

### 6.1 Enforcement criteria added in revision 7

The owner policy is testable through `AC-C67`–`AC-C72`:

- explicit-change preservation and scope visibility;
- confirmation gating and decline behavior;
- conversation-scoped memory and explicit promotion;
- chat-lifecycle isolation across chats;
- durable actor/reason/scope/previous/resulting/confirmation records.

`R-09` may specify the field taxonomy and promotion mechanics, but it may not relax this policy. Any broader scope requires a new owner decision.

Revision 8 adds the revision-7 delta corrections without changing `OD-02`, `OD-03`, or `OD-04`: the applied-marker commit boundary and restart behavior (`AC-C73`), atomic promotion with explicit demotion/revocation semantics (`AC-C74`), and derived-surface/new-conversation lifecycle isolation (`AC-C71`).

Revision 9 adds the mechanism-boundary corrections without changing those owner decisions: commit state and ownership in the canonical model (`AC-C75`), separation of atomic authorization from asynchronous materialization and read-time enforcement (`AC-C76`), four explicit revocation operations (`AC-C77`), and an observed non-receipt test in `AC-C71`.

## 7. Gate status

| Question | Answer |
|---|---|
| `O-11` closed? | **Yes, by `OD-04`** |
| `A2-08` resolved? | **Yes at the owner-policy level** |
| Architecture approved? | **No** |
| Research authorized? | **No** |
| Implementation authorized? | **No** |
| Runtime testing authorized? | **No** |

## 8. Related records

- Amendment 02 §1.2 — `OD-02`, layered state and continuity.
- Amendment 02 §1.3 — `OD-03`, transition confirmation threshold.
- Amendment 02 §1.4 — `OD-04`, explicit-change scope.
- Amendment 02 §7 — scenario, world, character, and lore transitions.
- Amendment 02 §14 — `DA-C4`.
- Amendment 02 §17 — `A2-08`.
- Muse review `M2-02` — `docs/review/REVIEW-AMENDMENT-02-REVISION-3-MUSE-2026-11-11.md`.
