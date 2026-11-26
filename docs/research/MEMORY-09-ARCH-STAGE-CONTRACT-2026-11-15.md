# MEMORY-09 — Stage contract: scope, ownership, knowledge and lifecycle

**Status:** design record, not implementation. Nothing in `internal/`, `app-*.go` or
`frontend/` is touched by this document or by this workstream.
**Date:** 2026-11-15
**Revision:** **r33** — owner decision **O-14** (2026-11-19): §8.4 is classified as **illustrative and
non-binding** (note at §8.4). No normative rule was changed; K-R7 and the K-R13 consequence stand
unamended. Previous revision — r32 — **labelling correction only.** A rule in R1d had been written as
`**I14**`, colliding with the invariant of that number and inflating invariant counts in
automated checks. It is relabelled `K-R90`; its wording is unchanged and no normative
content was added. r31 (`fc3081f3…`) extended detector identity to its dependency closure;
r30 (`a60cdba0…`) gave detectors their own implementation identity. r22 (`a924e767…`) added
the commit-time check; r21 (`a3fe5b1e…`) gave comparisons a lifecycle. r21 (`a3fe5b1e…`) gave
comparisons a lifecycle; r20 (`c6793b44…`) fixed the interval test and made comparison
history first-class. r19 (`432923e7…`) added the three-valued comparison; r18
(`f53e2490…`) fixed the conflict/temporal division. r18 (`f53e2490…`) fixed the conflict/temporal division and the temporal
states; r17 (`161c522e…`) scoped conflict by applicability. r17 (`161c522e…`) scoped conflict by temporal
applicability; r16 (`e78f47c9…`) bounded the transfer hash and split Q14. r16 (`e78f47c9…`) bounded the transfer hash and split Q14
into equivalence, incompatibility and `unresolved_identity`; r15 (`40628199…`) recorded
transfer provenance; r14 (`fca8b7c3…`) added the deletion contract. r14 (`fca8b7c3…`) added the deletion contract; r13 (`47f718c8…`) bounded audit
claims and fixed retention ordering. r13 (`47f718c8…`) bounded honest audit claims and fixed retention
ordering; r12 (`755396af…`) introduced replayability states. r12 (`755396af…`) introduced replayability states;
r11 (`83722e0c…`) extended evaluation provenance. r11 (`83722e0c…`) extended evaluation provenance and made mismatched
replay fail closed; r10 (`7be0faf0…`) added I15 and versioned declarations.
r10 (`7be0faf0…`) added I15 and the versioned-declaration proposal; r9 (`5a2971d3…`)
introduced the precondition tri-state; r8 (`e6ca9c0d…`) separated materiality from
exposure. r9 (`5a2971d3…`) introduced the precondition tri-state; r8
(`e6ca9c0d…`) formalised precedence and separated materiality from exposure; r7
(`af12f533…`) established summary-vs-account, continuation and recall. r8 (`e6ca9c0d…`) formalised precedence and separated materiality from
exposure; r7 (`af12f533…`) established summary-vs-account, continuation and recall; r6
(`f20014ac…`) specified the fold; r5 (`bf3bca39…`) closed projection purity; r4
(`eb93c01e…`) resolved acquisition-vs-edge authority, contexts and classes.
r7 (`af12f533…`) established summary-vs-account, continuation and recall; r6
(`f20014ac…`) specified the fold; r5 (`bf3bca39…`) closed projection purity; r4
(`eb93c01e…`) resolved acquisition-vs-edge authority, contexts and classes; r3
(`f8bbe6aa…`) introduced the edge model; r2 (`ade6281f…`) holds D1–D5. r6 (`f20014ac…`) specified the fold and removed `revoked_turn`;
r5 (`bf3bca39…`) closed projection purity; r4 (`eb93c01e…`) resolved acquisition-vs-edge
authority, contexts and classes; r3 (`f8bbe6aa…`) introduced the edge model; r2
(`ade6281f…`) holds D1–D5; r1 (`2a550d8b…`) is superseded on ownership.

Tags: **[E]** measured here with a receipt · **[L]** literature, not verified ·
**[X]** owner decision, authoritative design input · **[P]** proposed, awaiting owner
approval · **[D]** design detail proposed here, not yet reviewed.

**No claim here is evidence that the architecture works.** Three standing rules:

1. **Document validation is not behavioural validation.** Suites passing, probes defined and hashes matching prove the record is consistent — not that the contract works. No probe in §12 has been run.
2. **Accepting a review correction is not approving an architecture decision.** Folding a correction into a revision means the record is more accurate; it does not ratify the design it now describes.
3. **Proposals stay provisional** until examined with their consequences (§3.5, §14).

---

## 1. Why this workstream opens now

1. **[E]** A world-scope fact was unreachable under a character-scoped query (T10).
2. **[E]** A reranker preferred a retracted record over the current one even when the text said `(superseded …, no longer true)`.
3. **[E]** The single apparent temporal win did not replicate; two regressions on unseen cases.

Eligibility precedes ranking; truth state travels out-of-band; and, since r3, **memory
existence is not character knowledge**. r4 adds: **knowledge history is event-sourced,
retrieval is context-explicit, and ordinary world facts are not all knowledge edges.**

---

## 2. Owner decisions (unchanged)

| | Decision |
|---|---|
| **D1** | Experiential memories belong to the originating **chat**; `world`/`scenario`/`character_public`/`character_private` are **scopes**, not owners. Chat deletion removes memories, not reusable assets. No automatic cross-chat carryover. |
| **D2** | **Memory existence ≠ character knowledge.** Knowledge follows experience, exposure and authorised transfer — never presence, never mere existence. |
| **D3** | **Reusable assets ≠ experiential memories.** Assets keep their identity; chats reference them; asset content is context, never experience. |
| **D4** | **Provisional records** are retrievable for characters holding them, attributed, never promoted by retrieval. |
| **D5** | **Decay** affects surfacing only — never truth, ownership, scope, eligibility or history. |
| **D6** | **Unresolved conflict**: five-condition deterministic materiality predicate, application-owned detection, model instructs the user to clarify, never arbitrates; unequal authority is not a conflict. |
| **D7** | **Belief never mutates world truth**; character knowledge controls what the character can use; authoritative facts must not leak into a character's perspective. |

---

## 3. Review resolutions (r4)

### R1 — Acquisition history is authoritative; the edge is its projection [X]

Review point: *"AcquisitionEvent and KnowledgeEdge must not become competing sources of truth."*

| Object | Responsibility | Mutability |
|---|---|---|
| `AcquisitionEvent` | **Immutable historical evidence** that a character witnessed, was told, discovered, verified or lost something | append-only, never edited |
| `KnowledgeEdge` | **Current derived relationship** between character and memory | materialised projection, always reconstructable from events |

Rules:

- **K-R1** The edge for (character, memory) is the fold of that pair's acquisition events in turn order. No edge exists without at least one event.
- **K-R2** An edge carries `derived_from_event_ids`. A projection that cannot be reconstructed from its events is an **integrity failure**, not a policy question.
- **K-R3** Revocation is a new event (`via=revoked`), never a deleted history row. Forgetting changes current knowledge; it never rewrites what happened.
- **K-R4** Where provenance matters (assembly, audit), the *event chain* is rendered, not the edge alone: "told by B at turn 210, unverified" is history; "knows it" is state.

### R1b — Projection purity: what may enter the edge fold [X]

Review concern: *whether all changes affecting an edge's projection are represented in the
event history, including revocation and changes to the underlying memory's status.*

**That concern was justified against r4.** The edge carried `revoked_turn`, and assembly
applied the memory's `status` to decide how the character's knowledge should read. A
rebuild from acquisition events alone therefore could not reproduce the edge's meaning:
part of it depended on the memory record — exactly the outside-authority failure R1
exists to prevent.

Resolution: **the edge encodes knowledge only; truth status is never copied onto it.**

- **K-R11** The edge fold depends **only** on acquisition-family events for (character, memory): `witnessed | told | read | inferred | verified | narrator_grant | revoked`. Nothing else may influence it.
- **K-R12** Memory lifecycle is its **own authoritative stream** — `established | superseded | corrected | retracted`, plus `scope` and `knowledge_class` changes — and a memory's current status, class and scope are a **projection of that stream**, not mutable authoritative fields.
- **K-R13** "This character's knowledge includes a retracted claim" is a **read-time join** of edge × memory projection. It is never stored on the edge.
- **K-R14** Each stream is rebuilt and checked independently; a mismatch in either is an integrity failure. The cross-stream join has no stored form, so it cannot drift.

**Consequence for R4:** a retraction changes the memory stream only. The edge's
`derived_from_event_ids` is unchanged, the character still knows the (now retracted)
claim, and whether they have *learned of the retraction* is itself an acquisition event.

**Alternative considered and rejected:** folding status transitions into the acquisition
stream. It would restore rebuildability, but it puts one object's lifecycle inside
another object's history and duplicates memory identifiers across streams. Two streams
with a read-time join keeps each history about one thing.

### R1c — The acquisition fold, specified [X]

Review concern: *`revoked_turn` must be entirely reproducible from the acquisition
stream, including ordering and multiple learn/forget cycles; and revocation must not
erase the fact that the character once knew something.*

**Ordering.**

- **K-R15** Acquisition events for a (character, memory) pair are totally ordered by `(turn, seq)`, where `seq` is a monotonic per-chat counter. **No wall clock and no insertion order participates.** Events sharing a turn are resolved by `seq` alone.
- **K-R16** Event ids are unique; a duplicate id is rejected at write.

**Event validity.**

- **K-R17** `revoked` is valid only when the pair's ordered log currently ends in a learn event. A revoke with no open episode is **rejected at write**, not silently ignored, because a silent no-op hides formation bugs.
- **K-R18** A learn event is valid at any time. A learn while an episode is open is a *continuation* (the character learns again / confirms), not a new episode.

**Episode model.** A learn opens an episode; a revoke closes it. The fold is:

```
known(pair)  ⇔  ordered_log(pair)[-1] is a learn event
edge exists  ⇔  known(pair)                  # no edge when forgotten
```

- **K-R19** A revoked pair has **no edge**, not an edge with a "revoked" flag. "Forgot" is the absence of the current relation; the fact that it was once known lives in the log, not in a half-state row.
- **K-R20** `revoked_turn` is therefore **not a field**. It is derived on demand — the timestamp of the last revoke event — and appears only in the `knowledge_history` inspection view.

**Projection fields, defined from the current episode only:**

| Field | Definition |
|---|---|
| `acquired_turn` | `(turn, seq)` of the learn event that opened the current episode |
| `via` | strongest evidence kind **in the current episode**, precedence `witnessed > verified > read > told > inferred > narrator_grant` |
| `verified` | true iff the current episode contains a `witnessed` or `verified` event |
| `derived_from_event_ids` | the events of the current episode, in order |

**Multi-cycle example:** told at turn 10, revoked at turn 20, witnessed at turn 40.

- Edge: `acquired_turn=(40, seq)`, `via=witnessed`, `verified=true`.
- History: three events, all intact; "did A ever know this?" answers **yes**, from events.
- `revoked_turn` derived for inspection: turn 20 — it is history, not current state.

**Relearning after forgetting cannot shortcut:** the earlier `told` event does not count
as verification for the new episode, and the earlier revocation does not suppress the
later knowledge.

### R1d — Summary versus account; continuation; recall [X]

Review concerns: *`via` precedence compresses history; continuation needs an explicit
upgrade rule; forgetting and recall difficulty must not collapse into one boolean.*

**1. The edge is a summary, never an account.**

- **K-R21** `via` reports the strongest evidence in the current episode. It does **not** describe the full history: a character told a rumour and later witnessing the same event has `via=witnessed`, while both events remain in the log.
- **K-R90** Any question of the form "how did they learn this?" is answered from **acquisition events**, never from `via`. Rendering and audit read the log; eligibility reads the edge.

**2. Continuation may upgrade, never re-date.**

- **K-R22** A learn event while an episode is open may raise `via` along the precedence order and may set `verified`, but **must not change `acquired_turn`**; the episode's start is fixed by its opening event.
- **K-R23** Within an episode, `via` never downgrades. A weaker later event (e.g. `told` after `witnessed`) is appended to the log and does not lower the edge's summary.
- **K-R26** Precedence is a **fixed total order over kinds**: `witnessed > verified > read > told > inferred > narrator_grant`. Comparison is on kind only, so two events can never tie on `via` — no alphabetical, insertion, locale or hashing order participates in the result.

**3. Three states, not one boolean.**

| State | Representation |
|---|---|
| knows it now | edge present |
| knew it, forgot | no edge; learn+revoke events in the log |
| knows it but struggles to recall | edge present, **surfacing priority** reduced (D5) |

- **K-R24** Recall difficulty is a **surfacing** property, never an eligibility or truth property. It must not hide a record from a direct question, and it must never be represented by removing an edge.
- **K-R25** Forgetting (no edge) and difficulty recalling (edge, low surfacing) are distinguishable in every view that reports either.

### R2 — Retrieval is context-explicit [X]

Review point: *"Retrievable only for characters holding it" applies to character-perspective retrieval, not every application operation.*

| Context | May retrieve | Must never do | Use |
|---|---|---|---|
| `character_perspective` | memories with an active edge for that character + ambient-class records in scope | receive C1 authoritative truth as knowledge; receive other characters' private scope | dialogue, decisions, first-person narration |
| `world_state` | every chat memory and every asset, regardless of character knowledge | be rendered as a character's knowledge | consistency checks, validation, contradiction detection |
| `narrator` | world state + all edges + unresolved conflicts + restricted lore | state a secret as something the character knows; inject it into the character channel | scene construction, discovery setup, dramatic irony |
| `inspection` | world state + knowledge map + conflict map | feed the prompt | user-facing audit of who knows what |

- **K-R5** Every retrieval call names its context. **Default is `character_perspective`** (fail closed).
- **K-R6** A scene in which a character discovers something requires the `narrator` context to know what the character does not — i.e. R2 exists so that knowledge gating does not make discovery impossible.
- **K-R7** Moving information from a broader context into the character channel requires an `AcquisitionEvent`, never a context switch alone.

### R3 — Knowledge classes, not a single flag [X]

Review point: *"A newcomer shouldn't need to be told a city has a north gate, but shouldn't automatically know the secret passage or the haunted house's history."*

| Class | Examples | Access without an edge? |
|---|---|---|
| `ambient_authored` | the city has a north gate; the harbour is tidal | **Yes** — ordinary public world knowledge, available to any character in scope by default. Must be authored as ambient, never auto-promoted from an event. |
| `restricted_authored` | a secret passage; hidden lineage; the house's concealed history | **No** — requires an explicit grant edge (`via=read`/`told`/`witnessed`) or in-world discovery. |
| `experiential` | the earthquake happened; a character was wounded | **No** — requires an acquisition event. |
| `character_specific` | beliefs, suspicions, private observations | **No** — edge to that character only, ever. |

- **K-R8** `knowledge_class` is a field on memories **and** assets. Ambient applies only to authored world/scenario definitions.
- **K-R9** No record becomes ambient because many characters know it. Popularity is not authorship. **[D]** (needs review)
- **K-R10** An experiential event never downgrades to ambient by repetition or by age; decay does not change class (D5).

### R4 — Retraction does not correct minds [X]

Owner direction, recorded verbatim in substance: if A was told the king is guilty and the
accusation is later disproven, **the fact that A was told is not erased**.

- the acquisition event remains (append-only);
- the source memory becomes `retracted` or `corrected`;
- A's knowledge state reflects that the information it holds **was retracted**;
- A may still believe the accusation until they learn of the correction;
- the correction is itself information: applying it to A requires an acquisition event.

"The application correcting a fact must not magically correct every character's mind."

### 3.5 Seven directions — proposed, not approved [P]

Your provisional answers to the r3 questions, recorded as proposals. None is a decision
yet; each is carried into §14 pending approval.

| # | Direction | Effect if approved |
|---|---|---|
| Q1 | Materiality = deterministic predicate over whether the query/action depends on the disputed fact | unblocks S24–S26 (still unrunnable until the predicate is mechanical) |
| Q2 | Edge creation only from explicit acquisition events; never from presence | fixes K-R1 as the only creation path |
| Q3 | Group-level knowledge first, individualise when a character matters | keeps write volume sane; adds a group audience concept |
| Q4 | Preserve acquisition history; allow explicit revocation without deleting the memory | matches K-R3 |
| Q5 | A chat keeps the asset version it started with unless the user explicitly updates | pins `asset_ref` to a version, not a mutable pointer |
| Q6 | Preserve the edge; mark acquired information retracted/corrected | matches R4 |
| Q7 | Ordinary authored facts ambient by default; secrets and events require access/acquisition | matches R3 |

---

### 3.6 Responsibility rubric — a pattern for future requirements [D]

One vocabulary recurs: **events preserve what happened; projections represent current
state; joins combine independent truths; contexts decide what may be surfaced.**

Before any new requirement enters this contract, it must name its responsibility:

| If the requirement is about… | It belongs to | Test it by asking |
|---|---|---|
| what happened, and when | an append-only event stream | can it be replayed? |
| what is true now | a projection with a rebuild check | does it reconstruct from its stream? |
| how two independent truths combine | a read-time join | is there any stored copy that could drift? |
| who may see it now | a retrieval context | is the context explicit at the call site? |
| what an action requires to be true | an **action-schema declaration** | is it declared absent, or merely unwritten? |
| how something is phrased or emphasised | assembly | does it change eligibility? (it must not) |

**Rule [D]:** a requirement that cannot name its responsibility is under-specified, and
should not be added as a field. This rubric is itself proposed, not adopted.

## 4. Stage contract

| Stage | May do | May **never** do | Input | Output |
|---|---|---|---|---|
| **Formation** | Extract records; assign owner chat, scope, knowledge class, attribution, time, status; append acquisition events | Invent scope, class or attribution; infer knowledge from presence; auto-promote private observation to world truth or event to ambient | Evidence log | Records + events (+ edge projection) |
| **Eligibility filter** | Exclude by ownership, scope, privacy, status, time, character knowledge, conflict state — **for the named context** | Rank by relevance; widen the pool; treat presence as knowledge; apply character gating to world-state/narrator contexts | Records + request context + retrieval context | Eligible set |
| **Retrieval** | Score and rank from the eligible set | Return outside it; change status, class or knowledge | Eligible set | Ranked eligible candidates |
| **Reranking** (optional) | Reorder eligible candidates | Add a record; resurrect a retracted one; override ownership, knowledge or context; decide truth | Ranked candidates | Reordered candidates |
| **Assembly** | Build channels C1–C5 for the named context; emit clarification instruction | Grant knowledge by supplying a record; present provisional as verified; leak C1 into character perspective; arbitrate | Candidates + knowledge state + conflict state | Prompt evidence |

### Invariants

- **I1** Eligibility precedes ranking. [D]
- **I2** Ranking cannot alter truth, membership or knowledge. [D, E §1.2]
- **I3** Status, ownership, scope, attribution, class, audience and knowledge are structured and out-of-band. [E §1.2]
- **I4** History is never destroyed; supersession and retraction are states. [D]
- **I5** Ownership, scope and knowledge are enforced at retrieval time. [E for owner; D for scope/knowledge]
- **I6** Empty eligible set → abstain; never widen the pool. [D]
- **I7** Retrievability does not confer authority. [X]
- **I8** Application owns conflict detection; the model owns the clarification request. [X]
- **I9** Memory existence does not imply character knowledge; presence does not imply observation. [X]
- **I10** Asset content may be context, never automatically experience. [X]
- **I11** Acquisition history is authoritative; knowledge edges are its projection; any edge must be reconstructable from events. [X, R1]
- **I12** Retrieval context is explicit and never inferred; the character channel is the default. [X, R2]
- **I13** Knowledge class is authored, not inferred from popularity or age. [X, R3]
- **I14** The edge is a summary of current knowledge; the event log is the account of how it was acquired. Questions about the account are never answered from the summary. [X, R1d]
- **I15** **Absence of information is never treated as information about absence.** An unpopulated field, an empty list, a missing record and an unwritten schema each mean "not established", never "established as none". A system that cannot safely determine an answer says so instead of manufacturing certainty. [X]

  **Applies uniformly to:** missing acquisition history, empty knowledge sets, unspecified permissions, undeclared action requirements, missing asset versions, and incomplete retrieval evidence. "Known to be absent" and "not established" are distinct states in every one of them.
- **I16** **Conflict determines whether claims disagree; temporal applicability determines which claims apply to a given point or interval.** Neither layer answers the other's question: a system that conflates them reads ordinary world evolution as contradiction. [X, R1d]
- **I17** **A result earns authority from the conditions under which it was established, not from when or where it was written.** Applies beyond temporal comparison to retrieval, formation and conflict evaluation: never "it was computed most recently", always "it was computed under these conditions". [X]
- **I18** **The system may not use absence of evidence as evidence that nothing happened.** Every subsystem must name the **independent evidence** that would reveal its own silent omission — evidence it does not produce itself. A subsystem that cannot answer "what would show I failed to act?" is incomplete, not merely untested. [X]

---

## 5. Dimensions and where each lives

| Dimension | Mechanism | Notes |
|---|---|---|
| Ownership (chat) | field `owner_chat_id` | lifecycle/deletion boundary |
| Scope | field `scope` | applicability, not audience |
| Knowledge class | field `knowledge_class` | ambient / restricted / experiential / character_specific |
| Attribution | field `attribution{kind, by}` | observed / told / believed / asserted / authoritative |
| Audience | **derived** from edges + class + scope + context | never stored |
| Knowledge state | projection `KnowledgeEdge` | reconstructable from acquisition events only (I11, K-R11) |
| Knowledge history | append-only `AcquisitionEvent` | authoritative for knowledge |
| Memory lifecycle history | append-only memory events | authoritative for status/scope/class; status is projected, never a mutable field (K-R12) |
| Truth/lifecycle status | field `status` | provisional / active / superseded / retracted / corrected |
| Temporal provenance | fields `established_turn`, chain | |
| Confidence | field, separate from importance | |
| Importance | field | narrative salience |
| Surfacing priority | derived, decay-adjusted | never affects eligibility |
| Conflict state | field `conflict_group` | unresolved marker |
| Retrieval context | **call parameter** | character / world_state / narrator / inspection |

**Deliberate non-fields:** audience, `known_by`, and "current knowledge" are all derived.
Storing any of them alongside the event log creates the competing-source-of-truth failure
R1 exists to prevent.

---

## 6. Record model (A′, revised)

### 6.1 MemoryRecord (chat-owned)

`id`, `owner_chat_id`, `scope`, `knowledge_class`, `subject`, `predicate`, `value`,
`attribution{kind, by}`, `status`, `established_turn`, `supersedes`/`superseded_by`,
`correction_of`, `conflict_group`, `confidence`, `importance`, `surfacing_priority`,
`source_span`, `formation_method`, `asset_ref{asset_id, version}` (nullable).

### 6.2 AcquisitionEvent (authoritative, append-only)

`id`, `character_id`, `memory_id`, `turn`, `via` ∈ `witnessed | told | read | inferred |
verified | narrator_grant | revoked`, `from_character_id` (nullable), `source_memory_id`
(nullable), `confidence_in_source`, `evidence_span`.

### 6.3 KnowledgeEdge (projection)

`character_id`, `memory_id`, `acquired_turn`, `via`, `verified`,
`derived_from_event_ids[]`. Rebuilt by folding acquisition events for the pair (R1c);
never edited directly.

**Exists only while the relation is current** (K-R19): a forgotten pair has no row.
`revoked_turn` was removed in r6 — it is derived history, not projection state (K-R20).

**Carries no truth status, class or scope.** Those belong to the memory's own stream and
are joined at read time (R1b). An edge saying "A knows X" is compatible with X being
active, retracted or corrected: the edge describes the knowing, not the truth.

### 6.4 Consistency rules

- **Projection check** is deterministic and cheap enough to run on write: fold events → projection; compare to stored projection; mismatch is an integrity error surfaced to `inspection`, never silently resolved in favour of either side. **[D]**
- **Two streams, two checks.** Acquisition events project to edges; memory events project to the memory's status/scope/class. Neither projection may consult the other.
- **Truth of history** is always the relevant event log. *How someone came to know* is answered only from acquisition events; *what is currently true* only from memory events.
- **Truth of current access** is the join: edge (knowledge) × memory projection (status, scope, class) × retrieval context. The join is computed per request and never stored.
- **No status on the edge** (K-R13). If a field describes truth rather than knowing, it does not belong on `KnowledgeEdge`.
- **Fold determinism** (K-R15, K-R20): replaying the same ordered log produces the same edge, on any machine, in any insertion order; and a missing edge is not "unknown state" to be guessed — it is a determinate consequence of the last event.

### 6.5 Derived views

`known_by(memory)` = characters with active projection · `character_knowledge(c)` ·
`audience(memory, context)` = known_by ∩ scope ∩ class ∩ context rules.

**Belief** remains a MemoryRecord with `attribution.kind=believes`,
`knowledge_class=character_specific`, edge to the believer only. A character may hold a
belief and an edge to the contradicting fact simultaneously; assembly keeps them separate.

---

## 7. Data-model alternatives

| | A′ Event-sourced edges (recommended) | B Stored audience set | C Per-character copies |
|---|---|---|---|
| Complexity | Highest: events + projection + integrity check | Lowest | Medium |
| Consistency | Strong: one authority (events), projection verified | Weak: two mutable sources | Weak: N copies drift |
| Provenance | Strongest: full transfer chain | Poor | Medium |
| Isolation | Strong | Adequate for presence only | Strong |
| Retraction | Preserves that A was told, and that the fact was later retracted | Cannot represent | Copies diverge |
| Verdict | **Recommended** | **Rejected:** conflates audience with knowledge | **Rejected as authority;** viable as a read cache |

**Recommendation: A′.** R1 sharpened it: events own history, edges own access, and the two
are reconciled by a deterministic check rather than by convention.

---

## 8. Worked cases

### 8.1 Earthquake, new character at turn 300

1. Formation writes `scope=world`, `knowledge_class=experiential`, `status=active`, `established_turn=40`.
2. Events are appended only for witnesses (`via=witnessed`) and recipients (`via=told`). Presence is never an input.
3. At turn 300 the newcomer has no events, so no projection, so the `character_perspective` filter excludes it.
4. "The quake took the east wall" appends `AcquisitionEvent(told, from=B)` → projection; the memory becomes eligible for that character, marked unverified.
5. Seeing the collapsed wall appends `witnessed` instead.

### 8.2 North gate vs secret passage

- *"The city has a north gate"* → `ambient_authored`: a newcomer arriving at turn 300 may be given it as ordinary public knowledge. No event required.
- *"A passage runs behind the north wall"* → `restricted_authored`: requires discovery or a grant. A newcomer does not receive it, and the narrator context may hold it while the character channel does not.

### 8.3 Haunted house history

The house's concealed history is `restricted_authored`; entering the house produces no
knowledge. A portrait, a diary or an explanation appends an event, and only then does the
character know. The narrator context may use the full history to stage the discovery.

### 8.4 Retraction while A believes the accusation

`memory(retracted)` + A's event unchanged + projection live ⇒ A's channel: *"you were told
X; it has since been retracted; you still believe it"*. Another character who knows the
correction behaves differently. No mind is corrected by a data update.

> **Owner decision O-14 (2026-11-19) — this display is illustrative and non-binding.** The sentence
> above is a worked illustration, **not a requirement**, and it is not to be implemented as
> character-usable knowledge. A retraction reaches a character's channel only through an
> `AcquisitionEvent` (K-R7), like any other information. Acquisition is read **broadly**:
> `witnessed`, `told`, `read`, `inferred`, `verified`, `narrator_grant` (K-R11). An inference or a
> verification is an acquisition event **without being a disclosure event**; the system records the
> acquisition, it never derives one from status. The character's knowledge state — inspected
> by the **human** operator — is unaffected by this classification; only the channel display is
> withdrawn. The rule text of §8.4 is preserved above unedited; this note classifies it.

---

## 9. Prompt assembly: channels and guards

| Channel | Source | Rule |
|---|---|---|
| C1 Authoritative constraints | world state, resolved facts | for adjudication/validation and scene consistency; **never** presented as character knowledge |
| C2 Character knowledge | active projections + ambient-class records in scope | may be used and spoken |
| C3 Beliefs / suspicions | `attribution.kind=believes` | drive behaviour; never labelled fact |
| C4 Unverified reports | provisional + `via=told` | carry source attribution |
| C5 Unknown | absence of any event | absence is not denial |

Guards: **G1** no C1→C2 promotion · **G2** no C4→C2 promotion by retrieval, reranking or
repetition · **G3** absence is not a negative record · **G4** dramatic irony is a narrator
affordance, not a character licence.

---

## 10. Conflict semantics

Detection: same subject + predicate + scope + applicable time window, incompatible values,
unresolved by correction/retraction → `conflict_group`. Materiality: D6's five conditions.
Unequal authority (provisional vs established) is never material. Response: competing
claims + attribution + unresolved marker + instruction to ask. Resolution: only an
authorised correction clears the group.

### 10.1 Candidate mechanical predicate, revision 2 — proposed, awaiting approval [P]

Review of revision 1 found two false-negative risks:

1. **Condition 2 was too restrictive.** A question can depend on a disputed fact without naming it — whether to enter a room may depend on whether the room is haunted.
2. **Condition 3 conflated exposure with relevance.** A material fact omitted by retrieval would be judged immaterial, letting the system excuse its own miss.

Revision 2 separates **dependency** from **exposure**.

```
material(Q, G) ⇔ (1) at least one member of G is eligible for the character and context
             ∧ (2) Q's answer or intended action depends on G, where dependency is:
                     (a) Q or the action references G's subject/entity, or
                     (b) Q or the action references an attribute G's predicate determines, or
                     (c) G's predicate is a declared precondition of the action
             ∧ (3) G is unresolved after memory projection, authority and temporal rules
```

Exposure becomes an **assembly obligation**, not a materiality condition:

- **E1** If `material(Q,G)`, assembly must include at least one member of G, with attribution, under the same eligibility and context rules.
- **E2** If a material conflict exists, retrieval must consider all members of G regardless of initial rank. An omitted member is a retrieval failure to fix, never evidence against materiality.
- **E3** A negative materiality decision may never be justified by absence of exposure.

**Operationalisation.** (a) and (b) use a declared term/attribute map. (c) depends on an
action-schema **precondition state**, which is tri-state — not a boolean:

| `preconditions` | Meaning | Materiality behaviour |
|---|---|---|
| `declared_some: [...]` | the action has enumerated preconditions | evaluate them under (c) |
| `declared_none` | the action is **explicitly** declared to require nothing | valid; (c) simply never fires; (a)/(b) still apply |
| `undeclared` | preconditions not yet specified | **fail closed** — materiality is undetermined for this action |

- **K-R27** Absence of a precondition list is not a claim that no preconditions exist. Only an explicit `declared_none` may be treated as "nothing required"; `undeclared` must raise an evaluability error rather than silently yielding "immaterial".
- **K-R28** `undeclared` is a property of the action schema, not of any memory record. It is declared once, reviewed once, and does not migrate into per-query state.

**Why tri-state:** a boolean would make every schema without prerequisites look defective,
while a default-to-valid would let unwritten schemas quietly disable (c) for exactly the
actions most likely to need it.

**What this buys:** a material conflict is recognised even when retrieval missed it, and
the haunted-room case triggers clarification without the query ever mentioning haunting.

**Status: revision 2 is an accepted candidate (owner review, r9); formal approval is
pending, and nothing is adopted.** S24–S26 remain registered and **not runnable** until
that approval. On approval: S24 tests an (a)/(b) failure, S14n tests absent or irrelevant
(c), S25 satisfies all three, S42 tests dependency without literal naming, S43 tests
materiality surviving an initial retrieval miss, and S44 tests the tri-state behaviour
(`declared_none` evaluates normally, `undeclared` fails closed).

---

### 10.2 Schema versioning for reproducible evaluations — proposed, not adopted [P]

Review recommendation: *action-schema declarations should be versioned with the schema, so
historical materiality evaluations remain interpretable against the version they used.*

**Why it matters.** A materiality decision is an evaluation over three inputs: the query or
action, the conflict group, and the action's declared preconditions. If the declaration can
change in place, an earlier decision cannot be reproduced: the same query replayed against
a corrected schema may yield a different verdict, and nothing in the record shows whether
the difference is a system bug or a legitimate schema correction.

**Proposed rules:**

- **K-R29** Action-schema declarations are **immutable once published**; any change creates a new version. This mirrors D3: a schema is a reusable asset, not chat memory.
- **K-R30** Every materiality evaluation records `schema_version` alongside its inputs and verdict. Replay uses the recorded version, never the current one.
- **K-R31** Re-evaluating the same query under a newer schema is a **new** decision with its own record; it never overwrites or reinterprets the earlier one.

**Consequence:** a schema correction is auditable as a change of behaviour under a changed
declaration, not as a silent rewriting of past verdicts. It also makes S44's `undeclared`
state meaningful in hindsight: a schema that was `undeclared` at the time of a decision
stays `undeclared` in that decision's record, even after the schema is completed.

### 10.2.1 What an evaluation record must carry — proposed, not adopted [P]

A schema version alone does not make a decision reproducible. Replaying faithfully against
the right schema while unknowingly using a different memory projection or retrieval policy
would produce a verdict that looks valid and is not.

Proposed record contents:

| Field | Purpose |
|---|---|
| query or action text, as issued | the question actually asked |
| conflict group: member ids + `conflict_group` id | which records were in dispute |
| eligible-set snapshot: ids + status/class at evaluation time | what the system believed was in scope |
| `schema_version` | the declaration governing the action |
| memory-projection watermark | which projection state the join used |
| retrieval-policy version (k, thresholds, context) | how candidates were assembled |
| predicate revision (`10.1 rev 2`) | which rules produced the verdict |
| verdict + reason code | the conclusion and why |

- **K-R32** Replay must verify **every** recorded governing version. If any cannot be matched, the replay **fails loudly**; it may not emit a verdict from partial agreement.
- **K-R33** The evaluation log is its own append-only artifact. It adds no fields to memory records and no new state to the knowledge model — it is a record *of* decisions, not part of them.

**Minimisation [D]:** record identifiers, versions and verdicts rather than full content
snapshots; content can be resolved from the streams at the recorded watermark. This keeps
the log auditable without duplicating the corpus.

### 10.2.2 Replayability states — proposed, not adopted [P]

Identifiers and watermarks are only durably auditable while the artifacts they reference
exist. Retention therefore needs an explicit relationship with replay, and three outcomes
must never be reported as one.

| State | Artifacts | What can still be claimed |
|---|---|---|
| `reproducible` | every governing version + source content present | the decision can be re-executed and must reproduce exactly |
| `auditable_only` | decision record and content **hashes** present; source content archived | the verdict is historically verifiable (hashes match), but re-execution is not possible without restore |
| `expired` | required artifacts legitimately deleted under a recorded retention action | the verdict stands as history; replay is impossible and must say so |

- **K-R34** A replay result carries exactly one of these states. `auditable_only` and `expired` **must never** be reported as a successful reproduction.
- **K-R35** Deleting an artifact referenced by a live evaluation record requires an **explicit, recorded retention action**; a tombstone names the artefact and the policy decision that authorised removal. Silent garbage collection may not achieve it.
- **K-R36** Expiry affects **replayability only**. A verdict does not become false, withdrawn or revised because its evidence was lawfully deleted — and an `expired` record may never be quietly dropped from the log.
- **K-R37** **A hash proves integrity, not truth.** `auditable_only` claims exactly two things: the archived content matches what was hashed, and the chain of identifiers is traceable. It never claims the original evidence was accurate, complete or sufficient, and no downstream decision may cite it as evidence of correctness.
- **K-R38** **Retention is atomic from the audit's point of view.** For any observer at any instant, the state "artifact absent and no tombstone" must be unreachable. Ordering is therefore tombstone-first:

```
1. write tombstone (state=pending, artefact id, policy decision)
2. delete artifact
3. mark tombstone (state=committed)
```

  A crash may therefore leave a `pending` tombstone — visible, attributable and recoverable — but never a silent gap. Recovery resolves `pending` by confirming absence or restoring the artifact; it never deletes the tombstone to tidy up. **The storage transaction mechanism is not prescribed; the ordering invariant is.**

**Why three and not two:** a binary "replay works / replay fails" would report a lawful
retention decision as a provenance failure, and would make archival look like corruption.
Hashes preserve auditability across archival without keeping full content forever; the
`expired` state lets the system be honest instead of approximating.

**Status: proposed, not adopted.** Retention policy, archival and the deletion-authorisation
path are Q12; this subsection only fixes what the outcomes *mean*.

**Status: §10.2 as a whole — proposed, not adopted.** It also depends on Q5 (asset version
pinning), which is still open — the two should be decided together.

## 11. Lifecycle and retraction

```
evidence turn
   └─> provisional (status, class, attribution set)
          ├─> active        (episodic, preference, validated fact, belief)
          ├─> superseded    (replaced; history kept)
          ├─> corrected     (earlier statement wrong; correction chain)
          └─> retracted     (asserted, now known false; acquisitions retained)
```

Promotion per kind: episodic — on witnessing; preference — immediate for holder; belief —
never promotable to fact; rumour — stays provisional, repetition is not corroboration;
world fact — authority or corroboration; absence — attributed. Rejected: blanket
three-source rule, promotion by retrieval, decay altering truth, deleting history.

**Retraction applies to records, not to minds** (R4). A `retracted` record keeps its
acquisition events; characters keep their beliefs until they acquire the correction.

### 11.1 Chat deletion contract — proposed, not adopted [P]

Deletion is a **lifecycle** rule, not a UI confirmation. A chat must never be deleted while
leaving its experiential memories orphaned, and no memory may survive by accident.

| Item | On chat deletion |
|---|---|
| World-scope memories acquired in this chat | deleted unless manually transferred |
| Scenario-scope memories acquired in this chat | deleted unless manually transferred |
| Character-public and character-private memories from this chat | deleted unless manually transferred |
| Reusable world / scenario definitions | **retained** |
| Character definitions and personas | **retained** |
| Memories manually transferred elsewhere | **retained at their destination** |

- **K-R39** Ownership is scope-independent: every experiential memory of the chat is deleted with it, whatever its scope. `scope=world` does not exempt a memory from the chat's lifecycle. [links D1]
- **K-R40** **No automatic transfer.** Memories never silently migrate to another chat, world, scenario or character. Manual transfer is the only preservation mechanism, and it names both source and destination.
- **K-R41** **Warnings are computed, never estimated.** Any deletion prompt reports counts and categories enumerated from the actual deletion set at confirmation time; a stale or approximated count is a defect, not a rounding issue.
- **K-R42** **Explicit confirmation is required**, and the text states what is deleted, what is retained, and that the action cannot be undone. Example: *"Deleting this chat will permanently delete its associated world, scenario and character memories that have not been manually transferred. Reusable character, world and scenario assets will remain. This action cannot be undone."*
- **K-R43** Deletion removes the experiential records **and their dependent knowledge relations and acquisition history**. Dangling edges and orphaned events are forbidden. Deletion is recorded as a retention action under K-R35, so evaluations referencing those artifacts transition to `expired` by the same tombstone protocol — never to a silent gap.
- **K-R44** Reusable assets and transferred destinations are untouched. Deleting a chat must not delete a lorebook, a persona, or a record that now lives in another chat.

**D8 (proposed):** the deletion contract above, connected to D1 (ownership), D3 (assets),
R1 (acquisition history) and Q12 (retention).

### 11.2 Manual transfer — Q13 approved in principle [X]

Transfer moves **content and provenance**, never knowledge. A record arriving in a
destination chat carries no knowledge edges; characters there acquire it through play (D2).

| Carried | Not carried |
|---|---|
| content, scope, knowledge class, original attribution | knowledge edges |
| original provenance (`source_span`, `formation_method`, origin turn) | acquisition events of the source chat |
| `source_memory_id`, `source_chat_id`, transfer event id, content hash | any implication that a destination character experienced it |

- **K-R45** Transfer is an explicit event naming actor, source chat, source memory and destination. The destination record is marked `transferred_from` and its origin fields are **copied references, never rewritten**.
- **K-R46** No acquisition events or edges are copied. Any later knowledge in the destination is a **new** event with its own `via` (`read | told | narrator_grant`), attributable to that chat.
- **K-R47** The chain `original memory → transfer event → destination memory` must remain auditable even if the source chat is later deleted: the transfer record carries a **content hash and provenance stub**, so the destination's history does not dissolve into an unverifiable claim. Under **K-R37**, this hash proves *which* content was transferred and *what provenance was recorded* — never that the original memory was accurate, complete or still available.

**Why the chain, not the knowledge:** a destination character who "remembers" an event they
never experienced is a continuity error, and transferring information must not rewrite
anyone's history. The chain gives auditability without granting experience.

**Open — Q14, split into two relations.** These are not interchangeable and must not share
one mechanism:

| Relation | Condition | Permitted handling |
|---|---|---|
| **equivalence** (Q14a) | same subject, predicate and scope; values equal under a **declared normalisation** — not model judgement | additive `possible_duplicate` link; **both** records and provenance chains retained |
| **incompatibility** (Q14b) | same subject, predicate and scope; values cannot both be true under the applicable time and context (D6) | join forms a conflict group; both records retained; clarification only if material (§10.1) |

- **K-R48** Both relations are **additive only**. No transfer may delete, merge, overwrite or re-date either record, and both provenance chains survive association.
- **K-R49** Equivalence that cannot be mechanically decided is recorded as **`unresolved_identity`** — an explicit state requiring a normalisable value or inspection. It may never default silently to "distinct", consistent with **I15**.
- **K-R50** Incompatibility is evaluated over **temporal applicability**, not over a bare subject/predicate/scope key. Two values conflict only where their applicability intervals **intersect**; an empty intersection is temporal progression, never contradiction.
- **K-R51** A **partial** overlap scopes the conflict to the intersecting interval. Outside that window the records are independent facts; the conflict may be neither widened to the whole span nor narrowed to nothing.
- **K-R52** Disjoint intervals are recorded as **succession** — ordinary story progression. Succession is not contradiction and not `unresolved_identity`. Which value applies *now* is the temporal layer's question (D5), not the conflict layer's.

| Records | Intervals | Relation |
|---|---|---|
| gate open / gate closed | yesterday / today | succession — no conflict |
| gate open / gate closed | same instant | incompatibility — conflict group |
| gate open 08:00–14:00 / gate closed 12:00–20:00 | partial overlap | conflict scoped to 12:00–14:00 only |
| gate open 08:00–10:00 / gate closed 12:00–20:00 | disjoint | succession — no conflict |

**Undecided:** the equivalence normalisation itself (Q14a) and the applicability window used
for incompatibility (Q14b). Note that "same fact" with different provenance is equivalence,
not a contradiction — treating it as conflict would manufacture disputes out of agreement.

### 11.3 Temporal applicability conventions — proposed, not adopted [P]

These are conventions for Q14b. They are **not** folded into K-R50–K-R52 until approved;
if adopted they refine how intersection is computed, not what it means.

- **K-R53** Intervals are **half-open**, `[start, end)`. The instant a segment begins belongs to that segment, so adjacent intervals never overlap and never leave a gap.

Worked boundary set for `open [08:00, 14:00)` and `closed [12:00, 20:00)`:

| Query time | Expected |
|---|---|
| 07:59 | neither applies |
| 08:00 | open applies |
| 11:59 | open applies |
| 12:00 | conflict begins |
| 13:59 | conflict remains |
| 14:00 | closed applies |
| 20:00 | neither applies |

- **K-R54** An interval is one of four tiers, and the tiers are not interchangeable:

| Tier | Basis |
|---|---|
| `explicit` | the interval was established directly in the record or its source |
| `derived` | inferred from trustworthy event ordering or lifecycle history, with the derivation basis recorded |
| `partial` | one endpoint known, the other not |
| `none` | no defensible interval exists |

- **K-R55** An unknown endpoint is **not** infinity. Unknown start does not mean "since the beginning of time"; unknown end does not mean "still current". A `none` interval yields **`unresolved_temporal_applicability`** — distinct from `unresolved_identity` (K-R49), because the uncertainty is *when*, not *which record this is*. Reusing the identity state would conflate two different questions in one field.
- **K-R56** Comparing two intervals yields exactly one of three results: **`disjoint`** (evidence establishes no overlap), **`overlapping`** (evidence establishes an intersection), or **`indeterminate`** (evidence establishes neither). `partial` and `none` tiers can produce any of the three depending on what the known bounds prove.
- **K-R57** The result is computed from **provable** interval relations over known bounds, never from the tier label. For `A = [a₁, a₂)` and `B = [b₁, b₂)`:

```
disjoint     ⇔  a₂ ≤ b₁  ∨  b₂ ≤ a₁          (non-strict — equality is disjoint)
overlapping  ⇔  max(a₁, b₁) < min(a₂, b₂)     (strict)
otherwise    ⇒  indeterminate
```

  The asymmetry in strictness is deliberate. Under half-open intervals the shared instant belongs only to the interval that starts there, so `[08:00, 12:00)` and `[12:00, 14:00)` are **disjoint**. A strict comparison in the disjoint test is a boundary defect: it would manufacture spurious indeterminacy at every touching pair.
- **K-R58** **Indeterminate is never coerced.** It is not conflict and not succession: it opens no conflict group, it also permits no "no conflict" conclusion, and a material query over an indeterminate pair must gather temporal evidence or abstain — never guess a value.

Worked case from review — *A: gate open, `[08:00, ?)`; B: gate closed, `[12:00, 14:00)`*:

| Disjunct | Evaluation |
|---|---|
| `a₁ ≥ b₂` | `08:00 ≥ 14:00` — false |
| `b₁ ≥ a₂` | `12:00 ≥ ?` — unprovable |
| overlap `max(a₁,b₁) < min(a₂,b₂)` | `12:00 < ?` — unprovable |
| **Result** | **`indeterminate`** — A may have ended before B, during B, or after B |

- **K-R59** Every evaluated pair records its outcome, and four states are distinct:

| `comparison_state` | Meaning |
|---|---|
| `compared: disjoint` | compared, separation proven |
| `compared: overlapping` | compared, intersection proven |
| `compared: indeterminate` | compared, neither proven |
| `not_compared` | pair eligible but never evaluated |
| `discarded` | recomputation finished against evidence that had already changed; never authoritative |

- **K-R60** An absent comparison may **never** be read as established non-conflict. `not_compared` is recorded explicitly; a consumer needing the relation must either perform the comparison or abstain. This is E3's rule applied to comparison: a retrieval stage must not mistake its own omission for proven agreement.
- **K-R61** `compared: indeterminate` is persisted, not discarded when no conflict group forms. Dropping it would make the pair indistinguishable from one never examined — the exact collapse K-R58 forbids.
- **K-R62** A comparison is bound to **the evidence it was computed from**: interval tiers, known bounds, and derivation basis. That input fingerprint is recorded with the result. A comparison is a judgement about specific evidence, not a permanent property of the pair.
- **K-R63** When the underlying evidence changes — a missing endpoint established, a `partial` becoming `explicit` — the existing comparison becomes **`stale`** and is never silently authoritative. A consumer must recompute before relying on the relation.
- **K-R64** Re-evaluation **supersedes without erasing**: the new comparison is authoritative, the previous one is retained as history of what was knowable at the time. `superseded` and `stale` are distinct — stale means no replacement exists yet and the relation cannot be answered from old results.
- **K-R65** **Commit is a compare-and-set on the input fingerprint.** A recomputed result may become authoritative only if the evidence fingerprint still matches the current inputs at commit time. If evidence changed while the comparison was being computed, the result is **discarded**, never published as current.
- **K-R66** A discarded attempt is logged with its own fingerprint and is distinguishable from `stale` and `superseded`; it never answers a query. Discarding is not a failure of the computation — the computation was correct for evidence that no longer applies.
- **K-R67** **No last-writer-wins.** Authority derives from fingerprint match, never from commit recency; a later commit against older evidence cannot displace an authoritative result computed against current evidence.
- **K-R68** The fingerprint is **content-addressed over the complete relevant input state**, not a version counter: interval bounds and tiers for both records, derivation bases, the governing policy versions, and the **implementation identity** of the comparison itself. A version number that returns to an earlier value must produce the same fingerprint only if the content genuinely matches.
- **K-R69** A prior result is **reusable** exactly when the full fingerprint matches: same evidence content, same governing policies, same implementation identity. Then the commit check may accept it without recomputation. Match on evidence alone is insufficient.
- **K-R70** Same evidence under a **different implementation identity** or policy version is **not reusable** — the earlier result is retained as history under its own fingerprint and the pair is recomputed. Reuse never crosses a change in the rules that produced the answer.
- **K-R71** Fingerprint inputs are serialised **canonically**: fixed field order, fixed encoding, normalised numeric representation, sets sorted, and no locale-, map-iteration- or serialiser-dependent ordering. Equivalent inputs must hash identically regardless of how the fields happened to be assembled.
- **K-R72** Canonicalisation preserves **semantic distinctions**: absent, empty and null are three different states and must encode differently. Collapsing them would make `not established` and `established as none` hash alike — the I15 failure at the representation layer.
- **K-R73** Canonicalisation is part of **implementation identity** (K-R70). Changing the canonical form invalidates reuse exactly as changing the comparison does, even if the computation is untouched. Input that cannot be encoded canonically is an **error**, never a silent fallback to a looser encoding.
- **K-R74** **Numeric domain is declared, not inferred.** Within the declared domain, numerically equal values encode identically (`1`, `1.0` and `1e0` are one value); values outside it — NaN, ±Infinity, and any undeclared sentinel — are **rejected**, not encoded.
- **K-R75** **Text is UTF-8 and Unicode-normalised to a declared form** (proposed: NFC). Case is preserved: no implicit case folding, no locale-dependent collation. Normalisation is applied explicitly by the protocol, never left to the runtime.
- **K-R76** **Ordered sequences and unordered sets are distinct types** and encode differently. Sets sort by their **canonically encoded bytes**, never by locale-sensitive text comparison.
- **K-R77** The fingerprint protocol records its **digest algorithm and canonicalisation version**. A hash migration or canonical-form change is a new protocol version and **invalidates reuse** — old fingerprints are history, not authority.
- **K-R78** **No partial fingerprints.** If serialisation fails partway, no fingerprint is produced and no comparison result may become authoritative; a truncated digest is never a usable key.
- **K-R79** **Implementation identity is a declared identifier with its own version**, not a source-code hash. Irrelevant edits must not invalidate reuse, and meaningful changes must not silently preserve it — the identity is asserted deliberately and recorded with the result.
- **K-R80** **Three identities must all match for reuse**, and they answer different questions:

| Identity | Answers |
|---|---|
| **protocol identity** | which canonicalisation and digest rules produced the fingerprint |
| **implementation identity** | which deliberately versioned comparison behaviour produced the result |
| **input fingerprint** | which evidence and governing conditions were evaluated |

  A source hash answers only "did these files change?", which is neither question and cannot substitute for either.
- **K-R81** **Identity change is a release-governance act, not a refactor side effect.** A material change in comparison behaviour requires a reviewed, recorded identity increment with its rationale; the identity registry is the record. An identifier that no process is required to move becomes dangerously stable.
- **K-R82** Behaviour-affecting change without a corresponding identity increment is a **defect**, not a judgement call. It is detected and reported, never silently accepted — otherwise reuse authority would survive a change nobody noticed making.
- **K-R83** **Identity completeness is established outside the registry.** Detecting an unincremented change requires independent evidence that behaviour changed — controlled behaviour-change fixtures and reviewed change declarations — because a registry cannot prove its own completeness. A missing entry is not evidence that nothing changed (I15, applied to provenance).
- **K-R84** Every subsystem specification states an **omission detector**: the independent evidence that would expose a failure to act. Omitting that statement is an incompleteness defect at review, not a documentation gap. [links I18]
- **K-R85** **Detector coverage is bound to the subsystem contract it was written against.** A change to a subsystem's responsibilities triggers re-evaluation of its detector; coverage that no longer matches the current contract is **stale coverage** and is treated as no detector at all until re-established. This reuses K-R63's staleness model rather than inventing a second one.
- **K-R86** Detector validity is **demonstrated, not asserted**: it must have fired on a seeded omission for the **current** contract, and must not fire on a no-op. A detector validated against a superseded contract is stale coverage, however impressive its history.
- **K-R87** **A detector has its own implementation identity** (K-R79). Changing how a detector works — even with the subsystem contract and its coverage label untouched — invalidates that detector's demonstrated validity until it is re-seeded. A current label over a changed mechanism is stale coverage by another route.

The three-identity discipline therefore applies twice over, to results and to the detectors
that guard them:

| Question | Result | Detector |
|---|---|---|
| which rules produced this? | protocol identity | contract identity |
| which behaviour produced this? | implementation identity | detector implementation identity |
| what was it established from? | input fingerprint | seeded-failure evidence |

All three must match in both columns; in neither column may a label outlive the mechanism
it describes.

- **K-R88** Detector implementation identity includes its **material dependency closure** — the parsers, fixture runners, evidence sources and downstream libraries that can affect either detection behaviour or the validity of the evidence it consumes. The criterion is capability, not change: a dependency belongs in the closure if it *can* influence the outcome, regardless of whether it moved.
- **K-R89** The material set is **declared and reviewed** as part of the detector specification. A change to any declared material dependency invalidates demonstrated validity until re-seeded; updates outside the declared set do not. Discovering an undeclared material dependency is itself a defect (K-R83's rule, applied to detectors): the list cannot be certified by the list.

I18 is not new policy so much as the pattern the contract already follows in five places:

| Subsystem | Silent omission | Independent evidence |
|---|---|---|
| materiality | retrieval dropped a material conflict | E1–E3 assembly obligations |
| temporal comparison | pair never examined | `not_compared`, persisted (K-R59) |
| knowledge model | field never established | I15's three states |
| identity governance | change shipped unincremented | behaviour fixtures (K-R83) |
| audit | report asserts what was not checked | verification before reporting (§13) |

---

## 12. Evaluation contract — `E-SCOPE-CONTRACT`

| Version | Change | Reason |
|---|---|---|
| v1 | S1–S8 | initial |
| v1.1 | +S9–S15 | ownership, provisional, decay, conflict, belief |
| v1.2 | +S16–S29 | knowledge boundary, assets, leakage, conflict controls |
| **v1.3** | **+S30–S34** | event/edge authority, retrieval contexts, knowledge classes, retraction |
| **v1.4** | **+S35–S36** | projection purity: status never on the edge; two streams rebuilt independently (R1b) |
| **v1.5** | **+S37–S38** | acquisition fold: multi-cycle learn/forget/relearn; ordering and invalid-revoke rejection (R1c) |
| **v1.6** | **+S39–S41** | summary vs account; continuation upgrades; recall vs forgetting (R1d) |
| **v1.7** | **+S42–S43** | materiality revision 2: dependency without naming; exposure separated from relevance |
| **v1.8** | **+S44** | action-schema precondition state: `declared_none` valid, `undeclared` fails closed |
| **v1.9** | **+S45** | schema-versioned evaluations: immutable declarations, replay against recorded version (10.2) |
| **v1.10** | **+S46** | full evaluation provenance: mismatch in any governing version fails loudly (10.2.1) |
| **v1.11** | **+S47–S48** | replayability states and retention-linked deletion (10.2.2) |
| **v1.12** | **+S49–S50** | hash-integrity scope; retention ordering and crash recovery (K-R37/K-R38) |
| **v1.13** | **+S51–S55** | chat deletion contract (D8/11.1): no orphans, computed warnings, no knowledge inheritance |
| **v1.14** | **+S56–S58** | transfer provenance chain, no inheritance, duplicate handling (11.2) |
| **v1.15** | **+S59–S61** | equivalence vs incompatibility as distinct relations; explicit `unresolved_identity` |
| **v1.16** | **+S62–S64** | temporal applicability: succession vs conflict vs scoped overlap (K-R50–K-R52) |
| **v1.17** | **+S65–S67** | half-open boundaries, unknown endpoints, temporal-undecidability state (11.3) |
| **v1.18** | **+S68–S70** | three-valued temporal comparison; indeterminate never coerced (K-R56–K-R58) |
| **v1.19** | **+S71–S72** | boundary equality is disjoint; comparison history is first-class (K-R59–K-R61) |
| **v1.20** | **+S73–S74** | comparison lifecycle: evidence-bound, stale on change, superseded without erasure (K-R62–K-R64) |
| **v1.21** | **+S75–S76** | commit-time fingerprint check; discarded results never authoritative (K-R65–K-R67) |
| **v1.22** | **+S77–S78** | content-addressed fingerprint, ABA reuse, no reuse across implementation change (K-R68–K-R70) |
| **v1.23** | **+S79–S80** | canonical serialisation: stable ordering, absent≠empty, canonicalisation as identity (K-R71–K-R73) |
| **v1.24** | **+S81–S84** | canonicalisation protocol: numeric/text/collection domains, digest versioning, declared identity (K-R74–K-R79) |
| **v1.25** | **+S85–S86** | three-identity model for reuse; identity change as governed release act (K-R80–K-R82) |
| **v1.26** | **+S87** | identity completeness verified externally by behaviour fixtures (K-R83) |
| **v1.27** | **+S88** | every subsystem declares an omission detector; absence is an incompleteness defect (I18/K-R84) |
| **v1.28** | **+S89–S90** | detector coverage bound to contract; stale coverage counted as none (K-R85–K-R86) |
| **v1.29** | **+S91–S92** | detector implementation identity; re-seeding required after mechanism change (K-R87) |
| **v1.30** | **+S93–S94** | material dependency closure for detectors; declared set and undeclared-dependency defect (K-R88–K-R89) |

S1–S29 retained verbatim; v1.3 only adds. **Not frozen, not run.** A suite is frozen only
after independent review of its expectations.

### New in v1.3

| Probe | Expectation | Failure means |
|---|---|---|
| S30 | event/edge integrity: deleting an event leaves the projection unreconstructable and raises an integrity error | competing sources of truth (I11) |
| S31 | context separation: the same query retrieves the unknown earthquake in `narrator` and not in `character_perspective` | gating blocks scene construction, or leaks |
| S32 | ambient availability: a newcomer receives "the city has a north gate" without acquisition | ambient class missing |
| S33 | restricted not ambient: the secret passage is withheld from the newcomer, released after discovery | secrets leak as ambient |
| S34 | retraction does not correct minds: A still believes, the acquisition is intact, another character behaves differently | a data update rewrote a character's mind |
| S35 | retracting a memory leaves A's edge and its `derived_from_event_ids` unchanged; only the read-time join reports "retracted" | status was copied onto the edge (K-R13) |
| S36 | rebuilding each stream from its own events reproduces its projection exactly; neither rebuild reads the other stream | projection depends on outside authority (R1b) |
| S37 | told at t10, revoked at t20, witnessed at t40 ⇒ edge `acquired_turn=40`, `via=witnessed`, `verified=true`, no `revoked_turn` field; "did A ever know" still answers yes | multi-cycle state ambiguous, or forgetting erased history (R1c) |
| S38 | a revoke with no open episode is rejected at write; same-turn learn+revoke resolves deterministically by `seq`, never by insertion order | silent no-op hides formation bugs, or ordering is non-deterministic (K-R15/K-R17) |
| S39 | told a rumour then witnessing the same event ⇒ edge `via=witnessed`, while "how did they learn it?" from the log still reports the `told` event | the edge is being treated as a complete account (I14) |
| S40 | a continuation learn upgrades `via`/`verified` but leaves `acquired_turn` at the episode's opening event; a weaker later event does not downgrade | continuation re-dates the episode, or precedence is not monotonic (K-R22/K-R23) |
| S41 | forgetting (no edge) and struggling to recall (edge with reduced surfacing) remain distinguishable; the latter still answers a direct question | recall difficulty removes eligibility, collapsing two states (K-R24/K-R25) |
| S42 | a query that never names the disputed fact (should we enter the room?) still triggers clarification because an action precondition depends on it | dependency requires literal mention, so the trigger is missed (10.1 (b)/(c)) |
| S43 | a material conflict absent from the initial evidence still triggers, and assembly pulls a member of the group with attribution | materiality was decided by exposure, excusing a retrieval miss (E2/E3) |
| S44 | an action declared `declared_none` evaluates materiality normally (no false alarm); an action whose preconditions are `undeclared` fails closed with an evaluability error, never "immaterial" | valid schemas are flagged as defective, or unwritten schemas silently disable (c) (K-R27) |
| S45 | a materiality decision replayed against its recorded `schema_version` reproduces exactly; re-evaluating under a newer schema produces a **new** record while the earlier verdict is unchanged | historical verdicts are silently reinterpreted, or a schema correction looks like a system bug (K-R29–K-R31) |
| S46 | replay with a mismatched projection watermark, retrieval-policy version or predicate revision emits **no verdict** and reports the specific mismatch; only a complete match reproduces the decision | partial agreement manufactures a verdict that merely looks reproducible (K-R32) |
| S47 | with source content archived but hashes intact, replay reports `auditable_only`: verdict verifiable, re-execution refused | archival is reported as a successful reproduction (K-R34) |
| S48 | after a recorded retention action removes a required artifact, replay reports `expired` with the tombstone reference; the verdict remains in the log and is not marked false | lawful deletion is reported as mismatch/corruption, or the record is quietly dropped (K-R35/K-R36) |
| S49 | `auditable_only` output states integrity and traceability only; no surface reports it as correctness, completeness or sufficiency | hash match is presented as evidence of truth (K-R37) |
| S50 | a crash injected between artifact deletion and tombstone commit leaves a visible `pending` tombstone; no observer ever sees "artifact absent, no tombstone", and recovery never deletes the tombstone to tidy up | silent provenance gap (K-R38) |
| S51 | deleting a chat removes its experiential memories together with their knowledge edges and acquisition events; no dangling edge or orphaned event remains | orphaned knowledge state (K-R43) |
| S52 | the same deletion leaves reusable world, scenario, character definitions and lorebooks intact | assets deleted with chat (K-R44) |
| S53 | the confirmation reports counts enumerated from the actual deletion set, and a memory added after the prompt is invented makes the confirmation stale rather than proceeding | warning estimated or stale (K-R41) |
| S54 | a manually transferred memory arrives with provenance and **no** knowledge edges; characters in the destination chat do not know it until they acquire it | transfer silently carries knowledge (Q13/D2) |
| S55 | evaluations referencing deleted artifacts become `expired` through the retention protocol; none exhibits a silent gap and none is silently dropped | deletion bypasses retention bookkeeping (K-R43/K-R35) |
| S56 | a transferred record arrives with `transferred_from`, copied origin fields and **no** edges or acquisition events; the chain original → transfer → destination is auditable | transfer carries knowledge, or rewrites origin provenance (K-R45/K-R46) |
| S57 | after the source chat is deleted, the destination still resolves content and provenance through the transfer record's hash and stub | transfer chain dissolves when the source is deleted (K-R47) |
| S58 | transferring into a destination that already holds an equivalent memory merges and overwrites nothing; distinct identities are preserved and any incompatibility opens a conflict group | silent merge or silent overwrite (Q14) |
| S59 | an equivalent arrival receives a `possible_duplicate` link while both records and both provenance chains remain intact and separately auditable | association collapses one chain (K-R48) |
| S60 | an incompatible arrival opens a conflict group, retains both values and both provenance chains, and clarification is offered only when material | incompatibility is resolved by deletion or silent preference (K-R48/D6) |
| S61 | when equivalence cannot be decided mechanically, the relation is recorded as `unresolved_identity` and surfaced for inspection | undecided equivalence silently becomes "distinct" (K-R49/I15) |
| S62 | the gate open yesterday / closed today produces **succession**, no conflict group and no `unresolved_identity` | ordinary progression is reported as contradiction (K-R52) |
| S63 | gate open and gate closed asserted for the **same instant** open a conflict group and retain both records | simultaneous contradiction goes undetected (K-R50) |
| S64 | partially overlapping intervals conflict only within the intersection; a query outside it is not material, and the conflict is neither widened nor narrowed away | overlap scoped to the whole span, or to nothing (K-R51) |
| S65 | the boundary set evaluates exactly as tabled under half-open intervals: 07:59 neither, 08:00 open, 11:59 open, 12:00 conflict, 13:59 conflict, 14:00 closed, 20:00 neither | a boundary instant belongs to both claims or to neither (K-R53) |
| S66 | a `partial` interval with unknown start is not widened to the beginning of time, unknown end is not treated as current, and an unknowable intersection causes abstention rather than a guessed relation; tier and derivation basis are recorded | unknown endpoint silently becomes infinity, or a relation is guessed (K-R54/K-R55) |
| S67 | a `none` interval reports `unresolved_temporal_applicability`, distinct from `unresolved_identity`; the two states never share a field | temporal uncertainty is mislabelled as identity uncertainty (K-R55) |
| S68 | partial A `[08:00, ?)` vs explicit B `[12:00, 14:00)` returns `indeterminate`: no conflict group opens, succession is not asserted, and a material query abstains | indeterminate is read as conflict, or as no conflict (K-R56/K-R58) |
| S69 | a `partial` interval that is nonetheless provably disjoint (A `[08:00, ?)`, B `[–, 07:00)`) returns `disjoint`, not `indeterminate` | tier label overrides provable evidence (K-R57) |
| S70 | `indeterminate` is a distinct third result: it never appears as `overlapping` or `disjoint` in any output, and material queries over it are never answered by inference | three-valued result collapses to two (K-R56/K-R58) |
| S71 | touching half-open intervals `[08:00,12:00)` and `[12:00,14:00)` return `disjoint`, not `indeterminate` and not conflict; **at 12:00 only the second interval applies** — temporal applicability, not resolution choosing a winner | strict comparison manufactures indeterminacy, or the boundary is described as conflict resolution (K-R57) |
| S72 | `not_compared` and `compared: indeterminate` are distinct persisted states; neither is ever read as proven non-conflict, and an unevaluated eligible pair is visible as such | omission is read as agreement (K-R59–K-R61) |
| S73 | an `indeterminate` comparison becomes `stale` when a newly acquired event establishes the missing endpoint; no consumer uses the old result, and recomputation yields the now-provable relation | stale comparison silently remains authoritative (K-R62/K-R63) |
| S74 | after recomputation the prior comparison is retained as history of what was knowable, alongside the authoritative new one; `superseded` ≠ `stale`, and neither erases the record | supersession deletes the earlier comparison, or stale is confused with superseded (K-R64) |
| S75 | evidence advances from v12 to v13 during a v12 recomputation: the commit fails the fingerprint check, the result is `discarded`, and no query is answered from it until recomputation completes | stale-write race publishes a superseded result as current (K-R65) |
| S76 | a `discarded` attempt is logged distinctly from `stale` and `superseded`, and a later commit against older evidence never displaces an authoritative result computed against current evidence | recency is treated as authority (K-R66/K-R67) |
| S77 | evidence moving v12 → v13 → content identical to v12 produces a matching fingerprint, and the earlier result is reused as authoritative without recomputation | identical content needlessly recomputed, or treated as stale (K-R68/K-R69) |
| S78 | identical evidence under a changed implementation identity or policy version is **not** reused; the pair is recomputed and the earlier result retained as history under its own fingerprint | reuse crosses a change in the rules (K-R70) |
| S79 | the same inputs assembled in different field orders, set orders or numeric spellings produce an identical fingerprint and are reused | equivalent inputs hash differently, causing needless recomputation (K-R71) |
| S80 | absent, empty and null inputs encode to three distinct fingerprints; canonically distinct inputs never collide; unencodable input raises an error instead of falling back | absent and empty collide, or encoding silently loosens (K-R72/K-R73) |
| S81 | within the declared numeric domain `1`, `1.0` and `1e0` encode identically; NaN and ±Infinity are rejected rather than hashed | numeric spellings diverge, or non-finite values enter a fingerprint (K-R74) |
| S82 | NFD and NFC spellings of the same text encode identically; case-variant text does not fold, and no locale collation affects ordering | runtime decides normalisation or case folding (K-R75) |
| S83 | two independent implementations produce **byte-identical** canonical output for the same input before digesting, and the protocol records digest algorithm and canonicalisation version; a migration invalidates reuse | agreement assumed from one runtime, or a hash migration inherits old authority (K-R76/K-R77) |
| S84 | a serialisation failure yields **no** fingerprint and no authoritative result, and the implementation identity is a declared versioned identifier unaffected by irrelevant edits | partial digest used as a key, or identity inferred from source state (K-R78/K-R79) |
| S85 | reuse requires protocol identity, implementation identity **and** input fingerprint to match: a protocol change alone blocks reuse even with identical evidence and implementation | any two of three treated as sufficient (K-R80) |
| S86 | a behaviour-affecting change shipped without an identity increment is reported as a defect; the registry records who approved each increment and why | identity drift passes silently, or the registry is written after the fact (K-R81/K-R82) |
| S87 | a controlled behaviour-change fixture produces a material behavioural difference under an unchanged identity, and the defect is detected **from the fixture**, not from registry absence; the same check on a no-op edit produces no defect | completeness is asserted by the registry itself, so unincremented changes go undetected (K-R83) |
| S88 | every subsystem in the registry declares an omission detector naming evidence it does not produce; a new subsystem without one fails the review gate, and a declared detector that cannot fire on a seeded omission is rejected | silent omission is possible in a subsystem that no independent evidence would expose (I18/K-R84) |
| S89 | changing a subsystem's responsibilities without re-evaluating its detector marks coverage `stale`, and the subsystem then fails the gate as if it had no detector | outdated coverage counts as coverage (K-R85) |
| S90 | a re-established detector fires on a seeded omission for the **current** contract and stays silent on a no-op; history against a superseded contract is not accepted as evidence | prior success is accepted in place of current validity (K-R86) |
| S91 | changing a detector's mechanism while its contract and coverage label stay identical invalidates its validity until re-seeded; the stale label is not accepted | a current label outlives the mechanism it describes (K-R87) |
| S92 | revalidation requires **fresh** seeded-failure evidence against the new mechanism; evidence produced by the previous implementation is not carried forward | old demonstrations are inherited by new code (K-R87) |
| S93 | a parser or evidence-source change inside the declared material closure invalidates detection validity, while a formatting-only update outside it does not | all dependency movement invalidates, or none does (K-R88/K-R89) |
| S94 | a material dependency absent from the declared set is reported as a defect when discovered, and the detector's validity is reassessed against the corrected closure | the declared list is treated as proof of completeness (K-R89/I18) |

S24–S26 remain **registered and blocked** on the materiality predicate (Q1).

### v1–v1.2 probes (retained verbatim)

Restored as a table in r4: revision r3 listed these as definitions and r4 initially
referred to them only by number, which made "retained verbatim" untrue. Expectations are
unchanged from their registering versions.

| Probe | Expectation |
|---|---|
| S1 | a character-scope query reaches the world fact |
| S2 | another character cannot retrieve a private record |
| S3 | a current-state query excludes a retracted record |
| S4 | an explicit-history query returns it, flagged retracted |
| S5 | "post before the correction" returns the earlier value, marked corrected |
| S6 | competing records both survive; neither is overwritten |
| S7 | after decay, historical answers are unchanged |
| S8 | a mock ranking an ineligible record first still returns only eligible records |
| S9 | a memory from chat A is not retrievable in chat B |
| S10 | reusing character/world/scenario in a new chat imports no memories |
| S11 | chat deletion removes its memories; reusable assets survive |
| S12 | a rumour returns as provisional and attributed, never as established |
| S13 | after decay, status, scope and ownership are unchanged |
| S14 | a material unresolved conflict emits the clarification instruction |
| S14n | non-conflicting differences, and irrelevant conflicts, do not trigger |
| S15 | a belief is retrievable as belief; no world record changes |
| S16 | a character introduced after the event receives no memory of it |
| S17 | a witnessing character does receive it |
| S18 | a told character retrieves it, marked told and unverified |
| S19 | a newcomer inside the house does not know its concealed history |
| S20 | after a valid transfer, the same character retrieves it |
| S21 | a rumour stays attributed to its source, not shared cast-wide |
| S22 | the authoritative fact is absent from the unaware character's channel |
| S23 | a belief changes no world record; the character's prompt keeps the belief |
| S24 | a conflict irrelevant to the query triggers no clarification |
| S25 | an equal-authority, answer-relevant conflict emits the instruction |
| S26 | established fact vs provisional rumour triggers no clarification; both preserved |
| S27 | an asset definition is supplied as context; no character gains experience from it |
| S28 | experiential memories are gone after chat deletion; assets and lorebooks intact |
| S29 | with no edge to the only relevant record, the character abstains; the pool is not widened |

Failure on any probe means the corresponding invariant (I1–I13) is unenforced, not that
the system regressed: these expectations come from the contract, not from observed
behaviour.

### Coverage

I1–I14 are each covered: I11 by S30, S35, S36 and S37–S38 (rebuildability, status never
on the edge, stream independence, multi-cycle determinism); I12 by S31; I13 by S32–S33;
I14 by S39; I4 and R4 by S34, S35 and S37; I9–I10 by S16–S20, S27. R1d's remaining
rules are covered by S40–S41; materiality revision 2 by S42–S43; the precondition
tri-state by S44; schema-versioned replay by S45; complete evaluation provenance by S46;
replayability states and retention-linked deletion by S47–S48; integrity-scope bounding by
S49; retention ordering and crash recovery by S50; the deletion contract by S51–S55;
transfer provenance by S56–S58; additive relations and explicit identity states by
S59–S61; temporal applicability by S62–S64; boundaries and temporal states by S65–S67;
three-valued comparison by S68–S70; boundary equality and comparison history by S71–S72;
comparison lifecycle by S73–S74; commit-time fingerprint checks by S75–S76; fingerprint identity and reuse by S77–S78; canonical serialisation by S79–S80; the canonicalisation protocol by S81–S84; the three-identity model and identity governance by S85–S86; external verification of
identity completeness by S87; omission detectors by S88; detector coverage currency by
S89–S90; detector identity survival by S91–S92; dependency closure and its declaration by S93–S94.

**Still not behavioural evidence.** These are registered expectations derived from the
contract. S24–S26 remain blocked on the materiality predicate, and no probe has run.

---

## 13. What this document is not

Not an implementation plan, schema migration, UI design, or authorisation to build
knowledge propagation, conflict detection or prompt assembly. No component added, no
runtime changed, no suite frozen, no reranker fetched.

### 13.1 Document structure plan — proposed, not executed [P]

The record is now ~51 KB across thirteen revisions. Coherent, but a reader looking for the
contract has to read the deliberation. Proposed split, to be executed only on request and
with full traceability preserved:

| Document | Contains |
|---|---|
| **Normative contract** | current rules, definitions, invariants, unresolved decisions — the thing that gets frozen |
| **Decision rationale** | why each rule exists, alternatives rejected, worked examples, review findings |
| **Evaluation registry** | probe definitions, statuses, dependencies, results |
| **Revision history** | immutable references to prior revisions and their SHA-256 |

**Rule [P]:** the normative contract **points to** the rationale and evidence; it never
duplicates them. A split that copies rationale into two places would recreate the exact
two-sources-of-truth problem this contract spends its time preventing.

**Not executed.** Splitting mid-review would risk losing the reasoning trail while the
questions in §14 are still open.

## 14. Open questions

**All items below remain open and provisional.** Accepting the r4 review corrections into
this revision does not approve the architecture those corrections now describe.

Awaiting approval (the seven directions in §3.5 are proposals):
1. **Q1 materiality** — approve a mechanism, or keep S24–S26 blocked?
2. **Q2 edge creation** — explicit acquisition events only; is narrator-granted knowledge an acceptable `via`, and how is it marked?
3. **Q3 crowds** — group-level knowledge first: who owns a group, and when does a member individualise?
4. **Q4 forgetting** — allow explicit revocation events (representable today), or is forgetting surfacing-only?
5. **Q5 asset versioning** — pin each chat to the asset version it started with?
6. **Q6 retraction** — R4 recorded as direction; confirm the "A may still believe" clause.
7. **Q7 ambient facts** — approve the four classes, and decide who may author an `ambient_authored` record.

New in r4:
8. **Group knowledge identity** (Q3 follow-on): is a group a pseudo-character for edges, or a scope-like audience? Affects whether S-probes need a group actor.
9. **Ambient and beliefs:** may a character hold a belief that contradicts an ambient authored fact (a flat-earther in a spherical world)? Proposed: yes, as `character_specific`, never mutating the ambient record.
10. **Inspection exposure:** does `inspection` show the full knowledge map to the player, or is some of it withheld to preserve mystery?
11. **Schema versioning (§10.2) and asset pinning (Q5):** decide these together — immutable published schemas with replay against the recorded version, or mutable schemas with evaluations treated as point-in-time only?
12. **Evaluation-log retention and minimisation (§10.2.1–10.2.2):** keep full input snapshots for exactness, or identifiers plus watermarks at lower storage cost? And crucially: what is the authorised path for deleting an artifact a live evaluation record references, so that lawful retention never gets reported as provenance failure? Outcomes must be distinguished as `reproducible`, `auditable_only` or `expired` — never collapsed into "replay worked" or "replay broke". Also: approve the retention-action ordering protocol (K-R38) and decide who may authorise deletion of artifacts referenced by live evaluation records.
13. **Deletion contract (D8, §11.1):** approve the chat-deletion lifecycle. **Q13 is approved in principle** (§11.2): content and provenance move, knowledge does not.
14. **Q14a equivalence on transfer:** what declared normalisation decides that two records assert the same fact (making them `possible_duplicate` rather than in conflict)?
15. **Q14b incompatibility on transfer:** what applicability window (time, context) decides that competing values cannot both be true — and therefore open a conflict group? The interval algebra is now fixed (K-R50–K-R52): empty intersection means succession, identical instant means conflict, partial overlap is scoped to the intersection. Open: how an interval is derived for a record that was never explicitly dated. §11.3 proposes four tiers (`explicit`, `derived`, `partial`, `none`), half-open boundaries and a distinct `unresolved_temporal_applicability` state, all pending approval — including which event-ordering sources qualify as trustworthy enough to licence a `derived` interval. Neither merge nor overwrite is acceptable in any case.

## 15. Reranker status

BGE is **not shipped**; artifacts and history retained. It remains a registered,
replaceable selection stage permitted only to reorder eligible candidates. Reconsideration
requires out-of-band status/turn reaching the selection stage **and** S8 passing.
