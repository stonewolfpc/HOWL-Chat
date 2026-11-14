# HOWL Chat — Owner Decision: Memory Formation as the General Pipeline; Lorebooks as an Organization Layer

**Status:** **Owner decision record.** Establishes a scope decision. **Implementation is not authorized.** No source, schema, UI, or data change was made in producing it.
**Revision:** 2 — owner endorsed the direction on 2026-11-12 and **kept it provisional** until `Q-MF-01` and `Q-MF-10` are closed. Revision 2 adds three owner amplifications (`MF-05` one home ≠ one relationship; `MF-06` one pipeline ≠ one pass; `MF-07` model latitude), renames the capability to **general conversational knowledge formation**, and adds four open questions.
**Date:** 2026-11-12
**Decision owner:** Project owner (Stone Wolf Systems)
**Why this document exists:** The automated-lorebook research phase produced a finding that changes a proposal rather than confirming it. Automated lore *generation* was scoped as an independent subsystem while a general memory-formation pipeline was separately described. Those two do substantially the same job, and keeping both would mean two extractors scanning the same conversations, producing overlapping records, and then reconciling each other's output. The owner has decided to collapse them. That decision needs a controlling record so later design work does not rebuild the second system by accident, and so the substantial research already done is preserved rather than discarded.
**Connections:** reframes `docs/research/MEMORY-04-AUTOMATED-LORE-RESEARCH-ASSESSMENT-2026-11-12.md` (Tracks 2–8 re-homed; Track 1 downgraded); adjusts the `A13` entry in `docs/research/MEMORY-08-COMPARATIVE-RANKING-AND-UNRESOLVED-QUESTIONS-2026-11-12.md`; relates to `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md` (the evidence-backed proposal pipeline) and `docs/architecture/OWNER-DECISION-USER-PERSONAS-2026-11-12.md` (independent objects and attachment).
**Authority:** Owner requirement and scope decision. It does not freeze a design, select an architecture, or grant implementation authorization. It does not modify any frozen artifact.

---

## 0. Terminology note

The owner referred to a **"Form Memories"** proposal. **No file, symbol, document, or UI element with that name exists in this workspace**, and it does not appear in this conversation's history. The closest existing artifacts are:

- the **evidence-backed candidate/proposal pipeline** described in `DESIGN-MEMORY-CONTINUITY` (events → proposed memory items → review → canonical state), and
- the **Track 1 generation methods** in `MEMORY-04`.

This record treats "Form Memories" as the owner's name for that concept and does not attribute a specific document to it. If a separate proposal exists elsewhere, it should be linked here by whoever holds it; nothing in this decision depends on its exact wording.

---

## 1. The decision

Three statements, in order of consequence.

**`MF-01` — One pipeline, not two.** Howl has **one** memory-formation pipeline. It is not an "automated lore generator" plus an "automated memory extractor." There is one pass over conversation material, one provenance model, one candidate/review path, one deletion path, and one reconciliation discipline.

**`MF-02` — Automated lore-entry generation is downgraded.** It is no longer an independent proposed subsystem. It becomes:
1. an **unresolved research question** (is there any generation behaviour that the general pipeline cannot express?), and
2. at most, **one destination** that the general pipeline may write to.

**`MF-03` — Lorebooks remain, as organization.** Lorebooks survive as a first-class **organization and presentation layer** over retained knowledge: contextual collections of reusable knowledge that the user assembles, attaches, and reuses. They are a *way of organizing knowledge*, not *the mechanism responsible for creating it*.

The distinction in one line: **lorebooks are a shelf, not an author.**

**`MF-04` — Name the capability correctly.** This is not "automated lore creation by another route." It is **general-purpose conversational knowledge formation**: identifying what is worth preserving from a conversation, across kinds — events, facts, relationships, state changes, decisions, commitments, discoveries — and routing each to its proper home. The owner's framing is deliberate: the capability is larger than the subsystem it replaces, and the naming should reflect the larger thing.

**Status: provisional.** The owner supports this as a strong candidate for the permanent architecture, not as a settled one. `MF-01`–`MF-04` become permanent only when `Q-MF-01` (is a separate generator ever justified?) and `Q-MF-10` (can one schema carry every kind without losing their distinctions?) are answered by evidence.

---

## 2. What dissolves, what survives

| Existing concept | Disposition | Note |
|---|---|---|
| Automated lore-entry generation as an independent subsystem | **Dissolved** | Replaced by the general formation pipeline (`MF-01`) |
| Lorebook entries as a *kind* of retained knowledge | **Survives** | One destination among several |
| Character / world / scenario lorebooks as contextual collections | **Survives** | Collections of reusable knowledge; no ownership over their contents |
| Lorebook attachment to chats, worlds, characters, scenarios | **Survives** | Same attachment model as all other objects |
| Retrieval of lore entries | **Survives**, generalized | Retrieval finds relevant knowledge regardless of destination (`MEMORY-04` Track 3) |
| Curation UI for entries (create, edit, merge, delete) | **Survives** | Becomes curation of a knowledge collection |
| Second extractor with its own schema | **Dissolved** | The specific thing this decision exists to prevent |
| Research on activation, conflicts, scope, maintenance, deletion | **Survives in full** | Re-homed to the memory store (`§5`) |

**Nothing valuable is thrown away.** What is thrown away is the *duplicate mechanism*, not the knowledge about how knowledge should be organized and retrieved.

---

## 3. Destination taxonomy

Every formed record has exactly **one authoritative destination**. A second destination may hold a *view or link*, never a second copy.

| Record kind | Authoritative destination | Presentation | Owner's example |
|---|---|---|---|
| Reusable world knowledge | Knowledge collection (lorebook) | Entry in a collection | "The kingdom's capital is Veyr" |
| Event and its consequences | Episodic memory | Timeline / event record | "Mara promised to protect the village" |
| Character state development | Character state projection | State card, versioned | "Mara now distrusts the king" |
| Relationship state | Relationship graph | Relationship view, versioned | — |
| User decision or preference | User-scope store (persona-scoped) | Persona / preference view | "The user abandoned the northern expedition" |
| Open commitment or promise | Open-thread record | Thread list with status | "Mara promised…" (the obligation, not the event) |
| Detail that would otherwise be lost (name, object, injury, number) | Detail index with source span | Searchable detail | — |
| Scenario-local setup | Scenario object | Scenario view | — |
| Model interpretation or speculation | **Nothing**, or a clearly-labelled low-confidence note | — | "He seemed afraid" |
| Conversational noise | **Nothing** | — | Most turns |

### 3.1 The rule that prevents duplicate truth

> **One record, one authoritative home.** If the same statement could plausibly live in two places, the application chooses one home by declared policy, and the other place holds a reference — never a copy that can drift.

This is the same discipline the memory design already needs for events versus state, extended to collections. It is the specific defence against rebuilding two systems with one problem.

### 3.2 Destination is chosen by the application, not the model

The model proposes a record and its content. **The application decides the destination** from the record's kind, its scope keys, and its validity — deterministic rules, not model judgment (`C-B`). A model that chooses its own destination will put things in the wrong store, and the failure is invisible until someone notices a character "knowing" something that was only ever a user decision.

**`MF-05` — One authoritative home is not one possible relationship.** Ownership and reference are different axes, and the decision must not collapse them.

| Axis | Meaning | Cardinality |
|---|---|---|
| **Authoritative home** | The single record that holds the content | Exactly one |
| **References** | Contexts that may use that record | Many |
| **Relationship kind** | *Why* the context references it (relevant-to, known-by, attached-to, derived-from) | Per reference |

A character-state record may be referenced by a scenario, a world, and several conversations without being copied into any of them. A reference is an edge with a kind and a scope, not a duplicate. This preserves reuse — the whole point of independent objects — while keeping one truth per record.

**Consequence for deletion and scope:** a reference never carries the content, so it cannot resurrect it; and scope is evaluated on the reference, so the same record can be legitimately visible in one context and not another.

**`MF-06` — One pipeline is not one processing pass.** A single formation system may require several stages, and quality must not be traded for superficial simplicity:

| Stage | What it does | Needs global view? |
|---|---|---|
| 1. Chunk examination | Read a bounded segment; identify candidate material | No |
| 2. Candidate extraction | Emit structured candidates with source spans | No |
| 3. Cross-chunk reconciliation | Merge candidates that describe the same thing across segments | **Yes** |
| 4. Duplicate and contradiction analysis | Deduplicate; detect direct contradiction, refinement, temporal change, perspective difference | **Yes** |
| 5. Validation and routing | Apply schema, provenance, policy, destination, admission | No (per record) |

Stages 3 and 4 require a global view, which means buffering and ordering across chunks — an engineering constraint with a real cost (`MEMORY-05` covers the flush boundary as the natural commit point). "One pipeline" is a statement about *who owns the work*, not about *how many passes it takes*.

**`MF-07` — The model's judgment must remain useful.** Determinism should constrain what happens to a proposal, not what the model is capable of noticing.

| Decision | Owner | Why |
|---|---|---|
| Storage, IDs, scope keys, validity representation | **Application** | Must be stable and testable |
| Destination and admission | **Application** | Policy, not judgment |
| Whether a candidate is well-formed and fully sourced | **Application** | Mechanical |
| Salience: what is significant, recurring, or load-bearing | **Model** | Not expressible as a keyword or a fixed type |
| Meaning: what a passage actually asserts | **Model** | Requires interpretation |
| Relationship hypotheses and ambiguity resolution | **Model**, with application validation | Requires judgment |
| Entry prose, where prose is required | **Model** | Voice and readability |

**Rule:** the pipeline must accept an **`unclassified` proposal** — a candidate the model judges significant but cannot map to a known kind. The application routes it to review or a holding area rather than discarding it, and never silently coerces it into the nearest predefined type. Without this, the pipeline can only find what the taxonomy already anticipates, which defeats the reason for model-assisted formation.

---

## 4. Why this is the better architecture

| Problem with two systems | Consequence | Effect of one pipeline |
|---|---|---|
| Two extractors over the same text | Duplicate records, wasted tokens | One pass, one cost |
| Two schemas | Two dedupe rules, two conflict taxonomies | One schema with a `kind` and a destination |
| Two provenance models | Two deletion behaviours, one of which is wrong | One source-span contract |
| Two review queues | Review fatigue; one queue is ignored | One queue, one trust policy |
| Two activation systems | Two budget fights for the same context | One retrieval layer with destination filters |
| Two "what is canon" answers | Contradiction between systems | One authority per record |

**Evidence basis:** the failure of independently-developed memory and lore populations is already visible in the project record — the F3-C classification required three populations precisely because lore generation and memory storage grew separately and disagreed about ownership and provenance. That is a `Source-derived` project observation, not a hypothesis.

---

## 5. What the research still buys us

`MEMORY-04` is **not** invalidated. Its eight tracks are re-homed:

| Track | Original framing | New framing |
|---|---|---|
| Track 1 — Generation | "How to generate lore entries" | **"How to form memory records"** — one pipeline, many destinations. Lore generation is one candidate behaviour, not a subsystem |
| Track 2 — Quality and maintenance | "How to maintain lorebooks" | **How to maintain memory records** — dedupe, conflicts, supersession, validity, versioning. Unchanged in substance; applies to every destination |
| Track 3 — Retrieval and activation | "How to activate lore entries" | **How retrieval selects any knowledge record** — keyword/semantic/hybrid/budget. Applies to collections *and* to episodic and state records |
| Track 4 — Scope and identity | "How to scope lore" | **How to scope every record** — the four-layer ladder (attached → available → relevant → known) is a memory-wide requirement |
| Track 5 — Local-model suitability | "What local models can do for lore" | **What local models can do for formation** — constrained decoding, per-task split. Unchanged |
| Track 6 — Existing systems | "How products do lorebooks" | **What mechanics transfer** — field sets, budgets, attachment binding, bi-temporal validity. Unchanged |
| Track 7 — Orchestration | "When to run lore extraction" | **When to run formation** — checkpoint-based, no continuous background. Unchanged |
| Track 8 — Deletion and integrity | "How to delete lore" | **How to delete any record and its derived forms.** Unchanged; and now the single deletion path |

**The research conclusion that survives intact:** formation should be atomic, source-linked, candidate-first, checkpoint-driven, and budgeted — and it should never be continuous background extraction. That conclusion was reached about lore and applies more strongly to a general pipeline.

---

## 6. Effect on the ranking

In `MEMORY-08`, `A13` ("Automated lore") changes meaning:

| Before | After |
|---|---|
| A ranked approach: automated lore as a memory method | **A13 = the organization/presentation layer**, not a formation method. Its rank reflects that collections are useful for *organizing and presenting* retained knowledge |
| Generation treated as A13's mechanism | Generation is **subsumed by `A1`/`A6`** (event log + derived state; summary-as-index with raw evidence) and is no longer separately ranked |
| Lore "as primary memory" ranked low | Still incompatible, for the same reason: it duplicates state and transcript |

The general pipeline inherits the top-ranked position that `A1` already held. This decision does not promote a new approach; it removes a duplicate one.

---

## 7. Constraints carried forward (unchanged)

1. One extractor, one schema, one provenance contract, one deletion path.
2. Provenance is mandatory: every record names its source span.
3. Candidate-first for anything asserted about the world or a character.
4. Atomic records; collections are groupings, not documents.
5. Retrieval is budgeted and logged; drops are visible.
6. Scope is enforced before retrieval, structurally, on every record kind.
7. Model proposes; application decides destination and scope.
8. No continuous background extraction by default.
9. A collection never owns its contents; detach is not delete.
10. UI-only metadata (such as a persona `Role`) never enters any formed record.

---

## 8. What this decision does not do

- It does not authorize implementing the formation pipeline, altering the lore schema, or migrating legacy lore populations.
- It does not delete or supersede any research document. `MEMORY-04` remains valid with a changed status.
- It does not select a storage model, an ID scheme, or a review design.
- It does not decide whether automated lore generation survives as a distinct *behaviour* (`Q-MF-01`).
- It does not reopen frozen governance, F3-C contracts, or Amendment 02. It changes a proposal, not a frozen record.
- It does not resolve the legacy-population migration question (`Q-L9`), which remains open and remains constrained by the frozen F3-C classification.

---

## 9. Open questions

| ID | Question | Priority |
|---|---|---|
| `Q-MF-01` | Is there any generation behaviour that a general formation pipeline cannot express, and that justifies a separate lore generator? | **P0** — this is the question that decides whether `MF-02` is permanent |
| `Q-MF-02` | Does a record need a `kind` discriminator with a shared schema, or separate stores per kind? | P1 |
| `Q-MF-03` | Who owns destination policy — a declarative table, per-kind rules, or user configuration? | P1 |
| `Q-MF-04` | What does a "collection" mean structurally: a set of references, a tag, or an ordered document? | P1 |
| `Q-MF-05` | Can one record appear in two collections without becoming two records? | P1 |
| `Q-MF-06` | How does attachment apply to a collection versus an individual record? | P1 |
| `Q-MF-07` | How does review differ per destination (entry prose vs state change vs user decision)? | P2 |
| `Q-MF-08` | Does the user edit a destination-specific view, or the underlying record? | P2 |
| `Q-MF-09` | How are legacy lore populations mapped onto the taxonomy without violating the frozen classification? | P2 |
| `Q-MF-10` | What is the minimum viable formation schema that covers all destinations without becoming generic mush? | **P1** |
| `Q-MF-11` | Does the pipeline form records at checkpoints only, or also on explicit "remember this"? | P1 |
| `Q-MF-12` | How is the "one authoritative home" rule enforced and tested? | P1 |
| `Q-MF-13` | Does a collection ever become a retrieval unit distinct from its entries? | P2 |
| `Q-MF-14` | How does a deleted record behave in a collection that still references it? | P1 |
| `Q-MF-15` | What is the reference schema — kind, scope, and provenance of the edge itself? | P1 |
| `Q-MF-16` | How long may an `unclassified` proposal live, and who reviews it? | P1 |
| `Q-MF-17` | Where is the global-view buffer for stages 3–4 held, and how is it invalidated? | P1 |
| `Q-MF-18` | Can model-proposed salience override a routing policy, or only rank within it? | P2 |
| `Q-MF-19` | Does a reference ever need its own validity interval (attached-then-detached)? | P2 |

---

## 10. Summary

Automated lore generation was a proposal that solved a general problem with a specific tool. The owner has replaced it with the general solution: **one memory-formation pipeline, many stages, multiple destinations, one of which is a lorebook collection** — and has named the capability for what it is: **general conversational knowledge formation**.

The three amplifications in revision 2 are the conditions under which that general solution stays general: one home with many references (`MF-05`), one owner with many passes (`MF-06`), and application control over routing without application control over noticing (`MF-07`).

The research is kept. The duplicate subsystem is not. The decision remains **provisional** until `Q-MF-01` and `Q-MF-10` are answered by evidence.
