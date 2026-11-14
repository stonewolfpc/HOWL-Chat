# HOWL Chat — Owner Decision: User Personas as a First-Class Object

**Status:** **Owner decision record.** Requirements are established. **Implementation is not authorized by this document.** No source, schema, UI, or data change was made in producing it.
**Date:** 2026-11-12
**Decision owner:** Project owner (Stone Wolf Systems)
**Why this document exists:** The owner has established user personas as a first-class Howl concept with a defined product requirement set, including a UI surface, a UI-only Role field, and a chat-bar identity indicator. Until now the only persona code in Howl was a localStorage side-door with no reader, no identity, and no prompt presence. That gap was recorded in the audit below. A decision of this scope needs a controlling record so it is not re-litigated from chat history, and so later design and implementation work has one authority to point at.
**Connections:** extends the object/attachment/relevance framing recorded in `MEMORY-01` §9 and `MEMORY-08` §3; interacts with `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md` (identity, provenance, deletion), `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md` rev 9 (audience and character knowledge boundaries), and `docs/research/ROOT-01-PRODUCTION-ROOT-OWNERSHIP-AND-CONVERSATION-IDENTITY-2026-11-11.md` (canonical identity prerequisites).
**Authority:** Owner requirements. It does not freeze a design, select an architecture, or grant implementation authorization. Where it conflicts with an existing record, this record states the conflict and does not silently overwrite the other document.

---

## 1. The decision

**User personas are a first-class part of Howl, on the same architectural footing as characters, worlds, scenarios, and lorebooks.**

A persona is not a settings preset, not a display name, and not a property of a chat. It is an independently saved object the user creates, edits, reuses, and deliberately inhabits. The user is not only talking to characters — they are deciding who they are in each creative context, and Howl must support that as a continuous, customizable, controlled capability.

Three words the owner used to frame it, and which this record treats as the design intent:

**Continuation. Customization. Control.**

---

## 2. The three separate concepts

This distinction is the core of the decision. Conflating any two of them produces the exact failure the owner described: characters seeing the user as every persona at once.

| Concept | Definition | What it controls |
|---|---|---|
| **Persona library** | Every persona the user has created. Independently saved, reusable, discoverable, editable, deletable | What exists |
| **Context attachments** | Which personas are *available* to a particular chat, world, or scenario | What may be chosen here |
| **Active user identity** | The one persona the user has actually selected for the current interaction | What the AI and the characters actually receive |

**Rules that follow directly:**

1. A persona being attached to a scenario does **not** make every character in that scenario aware of every attached persona. Attachment is availability, not knowledge.
2. Characters receive **only the active persona**, never the library and never the attachment list.
3. The active persona is selected per interaction context, not globally by side effect.
4. Having ten personas available across a project must be indistinguishable, to the characters, from having one — until the user switches.

### 2.1 Why the library must be wider than the prompt

The owner's own example is the test case. The same person plays a loving father at home, a reporter at the Daily Planet, and a photographer in another universe. These are not three names for one identity; they are three identities with different relationships, responsibilities, and circumstances. Each is correct in its scenario and wrong in the others. If the library leaked into the prompt, every character would meet a person who is simultaneously all three, which destroys the point of having personas at all.

---

## 3. Established requirements

These are owner requirements, not design proposals. Each is stated so it can be verified later.

| ID | Requirement | Acceptance |
|---|---|---|
| `UP-01` | User personas are stored as first-class objects with stable IDs | A persona can be referenced by ID independently of its name or position in a list |
| `UP-02` | Multiple personas can be saved independently | Creating a second persona does not modify, replace, or reorder the first |
| `UP-03` | A dedicated **User Personas** column exists in the Character editor tab | The column is a persistent, first-class surface alongside the character grid, not a modal or a settings page |
| `UP-04` | Consistent creation, editing, browsing, and management | Create, open, edit, save, duplicate, and delete all behave the same way and are reachable from that column |
| `UP-05` | Personas are attachable to chats, worlds, and scenarios | An attachment can be created and removed without changing the persona object |
| `UP-06` | An explicit active persona is selected per context | The active persona is visible and changeable by the user; it is never inferred silently |
| `UP-07` | Characters receive only the active persona | The assembled prompt contains the active persona's AI-directed fields and no other persona's |
| `UP-08` | Persona identity reaches prompt assembly **and** chat presentation | Both the prompt and the visible transcript reflect the active persona consistently |
| `UP-09` | Switching personas preserves conversation continuity | No message is deleted, rewritten, reordered, or re-attributed by a switch |
| `UP-10` | Persona changes are scoped | Changing the active persona in one context does not change it in any other context |
| `UP-11` | Persistence is durable and typed | Personas survive restart, are not stored as one global blob, and are not lost by an unrelated write |
| `UP-12` | Deletion is complete and receipted | Deleting a persona removes it from the library, its attachments, and any derived surface; the receipt states what happened to messages attributed to it |
| `UP-13` | Full-wipe coverage | The application wipe removes personas. Today it does not |
| `UP-14` | The persona's `Role` field is UI-only metadata | `Role` never enters prompt assembly, character knowledge, retrieval, memory, or any AI-visible artifact |

### 3.1 Requirement `UP-14` is a hard boundary, not a style note

`UP-14` is the one requirement with a security-shaped failure mode. A user-facing label that quietly leaks into the prompt would tell characters how the *player* categorises their own persona — organizational shorthand the user never intended the AI to see. It is therefore specified as a **field-level boundary**, not as "be careful in the prompt builder":

- The persona record separates **AI-directed fields** (Name, Description, Additional notes) from **UI-only fields** (`Role`).
- Prompt assembly reads only the AI-directed set.
- Retrieval, memory extraction, and character knowledge read only the AI-directed set.
- A test asserts that no `Role` value appears anywhere in an assembled prompt, in any memory record, or in any retrieval candidate set.

If a future field is added that should also be UI-only, it joins the UI-only set by declaration and inherits that test.

---

## 4. UI specification (owner decisions)

### 4.1 Persona editor fields

| Field | Purpose | AI-visible? |
|---|---|---|
| **Name** | Who you are | Yes — the identity the characters meet |
| **Role** | How you identify this persona to yourself | **No — UI only** |
| **Description** | The actual identity/background information intended for the AI | Yes |
| **Additional notes** | Other persona details intended for the AI | Yes |

`Role` is free text (example: `Reporter`). It has no enumerated list, no validation beyond length, and no semantic role in the system.

### 4.2 Role tooltip (exact wording)

On hover, the `Role` field explains itself:

> This isn't included in the user prompt for the AI. It's for your organization, so you can keep your personas straight and know who you inhabit in each chat.

The tooltip is part of the requirement. It is what stops `Role` from being read as another AI-facing field.

### 4.3 Chat bar identity indicator

When a persona is active, the chat bar shows a permanent identity indicator:

| Property | Specification |
|---|---|
| Content | `Name - Role` (for example `Clark Kent - Reporter`) |
| Position | In the chat bar, adjacent to the message input |
| Size | Small font |
| Colour | White |
| Persistence | Always visible while that persona is active — not a transient badge, not a toast |
| Update | Updates immediately when the active persona changes, without a reload |
| Source fields | `Name` and `Role` only — **never** the description or notes |
| No persona selected | Show nothing. Do not invent a name, do not fall back to "User" as if it were an identity, do not display a placeholder persona |

**Why this exists:** it is the user's own orientation device. It answers "who am I in this chat right now?" without opening a panel, and it is the visible proof that the active identity is scoped rather than global.

### 4.4 Relationship between the indicator and the transcript

The chat bar indicator shows the **current** identity. The transcript shows what was actually said, under the identity in force when it was said. These are different questions and must not be collapsed: changing the active persona must not relabel past messages. See `UP-09` and §6.

---

## 5. What the AI receives

| Persona element | Reaches the prompt? | Notes |
|---|---|---|
| `Name` | Yes | The identity the characters address |
| `Description` | Yes | The AI-directed identity content |
| `Additional notes` | Yes | The AI-directed supplementary content |
| `Role` | **No** | UI only (`UP-14`) |
| Other personas in the library | **No** | `UP-07` |
| Attachment list | **No** | Availability is not knowledge |
| Persona ID | Not as text | May be used for bookkeeping and prompt caching; never rendered as content |
| Timestamps / audit fields | No | Never rendered |

**Open design item:** whether `Description` and `Additional notes` are assembled as one block or two labelled blocks is unresolved (§8). What is decided is that both are AI-directed and `Role` is not.

---

## 6. Continuity rules

The owner's requirement is explicit: switching personas changes who the user is presenting as. It must not change the conversation's history.

| Rule | Statement |
|---|---|
| No rewriting | A persona switch never edits, deletes, reorders, or regenerates existing messages |
| Attribution is historical | Each message records the persona in force when it was written; later switches do not alter that attribution |
| No contamination | A persona switch in one chat, world, or scenario has no effect on any other |
| No silent identity | The prompt never receives a persona the user did not select for that context |
| Reversible | The user can switch back and continue as before, with history intact |
| Visible | The active identity is always inspectable via the chat bar indicator (§4.3) |

**Consequence requiring a design decision:** if the user speaks as Clark Kent and then switches to the photographer, the earlier messages were said by Clark Kent. The system must be able to say that truthfully. That implies per-message persona attribution, which is a schema requirement, not a UI one (recorded as `Q-UP-04`).

---

## 7. Persistence, deletion, and wipe

| Requirement | Statement |
|---|---|
| Durable store | Personas live in a typed, durable store with stable IDs — not in a single localStorage blob |
| Attachment integrity | Deleting a persona removes its attachments; it does not orphan them |
| Message attribution | Deleting a persona must define what happens to messages attributed to it (label, tombstone, or retention with a "deleted persona" marker) — never silent rewriting |
| Derived surfaces | Persona-derived summaries, memory entries, or cached prompt packages are invalidated or re-attributed on deletion |
| Wipe | The application wipe removes personas and their attachments (`UP-13`) |
| Receipt | Deletion produces a truthful receipt of what was removed, retained, or unresolved |

**Current-state gap (evidence):** the wipe key list at `frontend/dist/js/app-wipe.js:15–25` does not include the persona key, so a full data wipe currently leaves the persona behind. That is a live `H6` deletion-gap instance in miniature.

---

## 8. Current state (audit evidence)

Recorded so the requirement set is grounded in what actually exists, not in assumption. Source inspection only; no runtime test.

| Aspect | Current state |
|---|---|
| Entry point | One "User Details" button on the Character Creator page (`character.html:118`) plus a `Ctrl+U` shortcut (`character.js:43`) |
| Fields | Name, Description, Additional notes — three free-text strings |
| Save | Writes one JSON object to `localStorage['howl-chat-user-details']` (`character-form.js:241`) |
| Read | Only the same modal reads that key. The chat renderer reads a **different** key, `howl-user-profile` (`chat-render.js:36`), which no code writes |
| Dead path | `chat-render.js:33` checks `currentChatData.playerCharacter.name`; the Go source contains no `player` reference at all |
| Chat label | Messages are appended with sender `'User'` (`chat-input.js:247`, `chat-actions.js:172`), and history falls back to `'User'` (`sidebar-session.js:303`) |
| Multiple personas | None. One global object, last-write-wins, no ID, no list, no active pointer |
| Prompt | No persona block in `app-prompt.go`; no `{user}` / `{{char}}` substitution anywhere in Go or frontend |
| Backend | Zero persona code, zero schema, zero migration |
| Deletion | Not covered by the wipe list |
| Tests | Zero |

**Conclusion:** there is no persona system today, only a form with a save button and no consumer. The redesign is therefore unconstrained by existing usage — nothing depends on the current shape.

---

## 9. Consequences for the memory architecture

The owner's examples force a conclusion that reaches beyond UI: **persona becomes a fourth scope dimension.**

Character knowledge is already a retrieval-authorization boundary (`H3`, `S2`). Worlds and conversations are already scope keys. A persona adds a user-side audience dimension with the same zero-tolerance property: the active persona's facts must reach the characters in that context, and no other persona's facts may, including through memory, summaries, lore, RAG, or cache.

| Consequence | Why |
|---|---|
| Persona key on memory records | A fact learned about the user under one persona must not resurface under another |
| Persona key on attachments | Availability must not become a candidate set |
| Persona-aware prompt caching | A cached package built for one persona must not be reused for another |
| Persona in deletion scope | Deleting a persona must invalidate persona-scoped derived data |
| Persona-aware evaluation | Persona leakage is a zero-tolerance failure class, sibling to character leakage |

**This is a genuine addition to the memory design's scope and is not covered by Amendment 02 rev 9**, which addresses character/conversation audience routing, not user-side identity. The conflict is recorded here rather than resolved by implication.

---

## 10. Open questions (investigation required before implementation)

| ID | Question | Why it blocks |
|---|---|---|
| `Q-UP-01` | Where does the persona store live (which store adapter, which file/database), given `DA-M1` is unresolved? | Determines `UP-11` |
| `Q-UP-02` | What is the persona ID scheme, and how does it interact with the unresolved canonical root/chat identity (`ROOT-01`)? | Stable IDs cannot be built on unstable identity |
| `Q-UP-03` | What happens to the existing localStorage persona on first run — migrate, adopt, or discard? | Silent data loss risk |
| `Q-UP-04` | Are messages attributed to a persona ID, and how is a deleted persona's attribution displayed? | `UP-09`, `UP-12` |
| `Q-UP-05` | Is attachment many-to-many across chats, worlds, and scenarios, and does a chat inherit attachments from its world/scenario? | Defines availability semantics |
| `Q-UP-06` | Is the active persona stored per chat, per scenario, or per chat-and-scenario pair? | `UP-06`, `UP-10` |
| `Q-UP-07` | Is there a default persona, and if so does it auto-apply to new chats? | Risk of an implicit, unselected identity |
| `Q-UP-08` | How are Description and Notes assembled — one block or two labelled blocks — and in what order relative to world/scenario/character? | Prompt quality and caching |
| `Q-UP-09` | Do persona facts enter memory and structured state, and under what scope key? | `H3`/`H4` interaction (§9) |
| `Q-UP-10` | Do persona switches need a conversation marker (an event in history) so the transcript is explainable? | Continuity truthfulness |
| `Q-UP-11` | What does the chat bar indicator do in a chat with no persona, and in a chat where the persona was deleted? | `UP-06` edge |
| `Q-UP-12` | Is persona export/import required, and in what format? | Portability, not blocking |
| `Q-UP-13` | How is `UP-14` enforced structurally so a future field cannot leak? | Boundary durability |
| `Q-UP-14` | Does the persona affect the "user" sender label in stored history, and does that change existing transcripts? | Backward compatibility |

---

## 11. What this document does not do

- It does not authorize implementation, schema changes, or UI changes.
- It does not select a storage adapter, an ID scheme, or an architecture.
- It does not modify Amendment 02 rev 9, the F3-C frozen contracts, `DESIGN-MEMORY-CONTINUITY`, or any frozen artifact.
- It does not resolve the memory-scope implications in §9; it records them as a requirement on the memory design.
- It does not claim the existing localStorage persona is safe to migrate without a decision (`Q-UP-03`).

**The next step is owner review of this requirement set and a separate implementation authorization that names the scope (store, schema, UI, prompt, wipe).**
