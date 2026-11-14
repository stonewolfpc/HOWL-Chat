# HOWL Chat — Architecture Decision: The Memory System

> ## Governing principle
>
> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**
>
> **Memory guides the model's understanding; it does not dictate its behavior.**
>
> Every layer, record type, retrieval method, and acceptance test below is justified against these two sentences. A capability that does not serve them is not part of the memory system, however impressive it is in isolation.

**Status:** **Owner architecture-direction decision.** The direction is adopted. **Implementation is not authorized.** The schema, retrieval algorithm, chunking strategy, reasoning budget, and model selection are deliberately **not** frozen.
**Revision:** 2 — adds the **guidance-not-command** principle, the five governing principles, and the four-link evaluation chain (owner, 2026-11-12). These govern implementation, not only research.
**Revision:** 3 — adds presentation **diagnostics** (why, not only which), a **mandatory transcript-only baseline**, **corpus representativeness** requirements, and the **failure-classification taxonomy** (owner, 2026-11-12).
**Date:** 2026-11-12
**Decision owner:** Project owner (Stone Wolf Systems)
**Why this document exists:** The research programme reached a decision point. `MEMORY-01`–`MEMORY-08` established what is known; `OWNER-DECISION-MEMORY-FORMATION` established that automated lore generation is not a subsystem; this record establishes **what the memory system is for, what shape it takes, and what evidence would justify committing to it**. It exists so that further work is measured against a stated direction instead of continuing to enumerate technologies.
**Connections:** `docs/research/MEMORY-01`–`MEMORY-08`; `docs/architecture/OWNER-DECISION-MEMORY-FORMATION-2026-11-12.md`; `docs/architecture/OWNER-DECISION-USER-PERSONAS-2026-11-12.md`; `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md`; `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md`; `docs/architecture/AMENDMENT-02-...-2026-11-11.md` rev 9; `docs/research/ROOT-01-...-2026-11-11.md`.
**Authority:** Owner decision on direction. It does not freeze a design, close a research question by assertion, modify a frozen artifact, or authorize code.

---

## 1. The governing principle, unpacked

| Clause | Technical meaning |
|---|---|
| **Right memories** | Accuracy, significance, durability, and usefulness matter more than record count |
| **Right times** | Retrieved when relevant, updated when circumstances change, withheld when not needed |
| **Right places** | Respect character knowledge, persona scope, world and scenario boundaries, and the context a thing belongs to |

**A fourth implication, and it is the one most systems get wrong:**

> **Sometimes the correct memory action is to retrieve nothing.**

A system that injects an irrelevant or misleading memory is worse than one that stays silent. A character must not know something merely because Howl has it stored.

**Consequence for evaluation:** raw recall is not a success metric. A system recalling 95% of everything while constantly injecting irrelevant material can perform worse in real roleplay than one recalling less with high precision. This is the difference between a memory database and a memory system.

**Corollary for formation:** the smallest useful representation is the target. Extracting 500 records is not better than extracting 40; it is usually worse — more duplicates, more contradictions, more noise, more prompt consumption, more chances to break continuity.

---

### 1.1 Guidance, not command

| Approach | What it does | Why it fails or succeeds |
|---|---|---|
| **Commanding** | Injects every remembered fact, relationship, event, and detail; expects the model to obey the accumulated record | Produces recitation, over-exposition, and brittle replies; a local model under context pressure may follow the wrong part of it |
| **Guiding** | Supplies carefully selected, relevant information in a format that helps the model understand the situation and keep continuity | Preserves interpretation, voice, and improvisation; the model reasons over guidance instead of obeying a script |

**The local-model reason is decisive.** Frontier systems compensate with retrieval, ranking, summarization, and orchestration infrastructure behind the API. Howl cannot assume any of it. Local models differ substantially in context utilization, instruction following, reasoning, prompt-structure sensitivity, structured-output reliability, and ability to reconcile conflicting information. A smaller model may accept a large context and still perform worse when the important part is buried in it.

**Therefore:** Howl compensates on the application side. It does not ask more of the model — it asks less of it, more precisely.

### 1.2 The four-link chain

Every memory decision is evaluated across the whole chain, not at a single point:

```text
1. AVAILABLE KNOWLEDGE        everything Howl has legitimately retained
        |
2. CONTEXTUAL SELECTION       what matters to this conversation, character, and moment
        |
3. MODEL-APPROPRIATE          the clearest, least intrusive way to present it
   PRESENTATION
        |
4. NATURAL INTERACTION        the outcome we actually care about
```

**The rule this chain exists to enforce:** a memory can be factually correct and still be harmful — presented at the wrong time, to the wrong character, or in the wrong context. Correctness at link 1 does not excuse failure at links 2–4.

**Testing consequence:** each link is measured separately, and a good result at one link never substitutes for another. A retrieval score is not a presentation score is not a creative outcome score.

---

## 2. What the research changed

Seven things were being treated as separate memory systems. They are related responsibilities of one system:

| Problem | Actual home |
|---|---|
| Remembering what happened | Original evidence + episodic records |
| Knowing what matters | Formation (significance), not retrieval |
| Tracking changes over time | Temporal records + supersession |
| Preserving character knowledge boundaries | Scope enforcement on every record |
| Organizing reusable world knowledge | Knowledge collections |
| Finding relevant historical information | Retrieval over evidence and records |
| Supplying the right information at the right time | Prompt assembly + budget + abstention |

---

## 3. The four layers

```text
1. ORIGINAL EVIDENCE          authoritative, immutable, segmented, one logical conversation
        |
2. MEMORY FORMATION           application-controlled interpretation pipeline (AI-assisted, optional method)
        |
3. KNOWLEDGE AND STATE        authoritative homes for retained information; typed records + references
        |
4. RETRIEVAL AND ASSEMBLY     identity-aware selection, ranking, budget, assembly, cache
```

| Layer | Responsibility | Owner | Explicit non-goal |
|---|---|---|---|
| **1. Original evidence** | Conversations, messages, timestamps, speakers, attachments, source references; segmented storage that never fragments identity | Application | Never rewritten by derived layers |
| **2. Memory formation** | Chunked reasoning, candidate generation, reconciliation, validation, provenance; produces candidates, not canon | Application orchestrates; model interprets | Not a mandate that everything be extracted |
| **3. Knowledge and state** | Episodic records, character state, user/persona state, world knowledge collections, decisions, commitments, relationships, validity | Application | Not a universal record that flattens all kinds |
| **4. Retrieval and assembly** | Identity-aware retrieval, lexical/semantic/historical matching, ranking, budget allocation, assembly, cache | Application | Not "inject the most" |

**The layers cooperate; they do not become four agents.** They are four responsibilities of one system, not four subsystems competing to maintain the same information. Formation does not retrieve; retrieval does not write; evidence is never authored by either.

---

## 4. Capability disposition

| Capability | Disposition | Note |
|---|---|---|
| Original conversation history | **Essential** | Authoritative evidence; never replaced by summaries |
| Chunked AI-assisted memory formation | **Strong candidate — pursue** | Optional method, not a dependency |
| Separate thinking-token budget | **Investigate and test** | Optimization, not a requirement |
| Deterministic extraction | **Keep, complementary** | Provides availability without successful AI extraction |
| Automated lore generation | **No separate generator** | `OWNER-DECISION-MEMORY-FORMATION` |
| Lorebooks | **Keep** | Reusable knowledge collections; organization, not authorship |
| Exact keyword retrieval | **Keep** | Names, terms, quoted phrases |
| Semantic / vector retrieval | **Investigate and test** | Flexible, less controllable |
| Hybrid retrieval + reranking | **Strong candidate** | Coverage plus precision |
| Character and persona knowledge boundaries | **Non-negotiable** | Zero-tolerance gate |
| Temporal state and supersession | **Non-negotiable** | Current truth and history both preserved |
| Continuous background extraction | **Avoid** unless evidence shows a real need | Cost, cache thrash, unreadable volume |
| Model tool calling | **Never a dependency** | The model must never have to manage its own memory |
| Monolithic conversation summaries | **Not authoritative memory** | Index at most; detail loss is documented |

### 4.1 Three architectural guarantees

These are guarantees, not quality improvements:

1. **AI-assisted formation is never the only path to information.** Howl must be able to retrieve original messages, find exact names and terms, and answer historical questions without any successful extraction.
2. **Formation must not erase uncertainty.** If the source does not support a claim, the application must not silently promote it to fact. Low confidence is stored as low confidence.
3. **Retrieval must respect knowledge boundaries.** Presence in the database is not entitlement to know.

---

## 5. Information classes and their validation rules

The pipeline can be universal; the records must not be interchangeable.

| Class | What is preserved | Validation rule |
|---|---|---|
| **Episodic** | What happened, when, who was involved | Source span required; immutable once committed |
| **Character state** | Current beliefs, feelings, relationships, condition, changes | Versioned; supersession required; current view derived |
| **User / persona state** | The user's relevant identity and circumstances in context | Persona-scoped; never merged across personas |
| **World knowledge** | Reusable facts about places, people, objects, rules, history | Reusable; scope-keyed; may be edited by the user |
| **Decisions** | What was decided, by whom, whether still in effect | Supersedable; the historical existence remains visible |
| **Commitments / open threads** | Promises, unresolved events, plans, obligations | Status-tracked until fulfilled, broken, or abandoned |
| **Historical evidence** | Original material retrievable when a summary is insufficient | Never regenerated; always addressable |

### 5.1 Memory represents current truth *and* how we got there

Every class keeps both: the current state and the version history that produced it. A character's opinion changing must not rewrite the event it is about; a decision being superseded must not erase that it was once made.

### 5.2 Three things that must not be conflated

| Concept | Question |
|---|---|
| **The information** | What is asserted |
| **Its authoritative storage location** | Where the one true record lives |
| **The contexts in which it is relevant** | Which contexts may reference it, and in what role |

This preserves the independent-object and temporary-attachment philosophy: one home, many references, relevance evaluated per context.

---

## 6. What testing must constitute

Not a generic "memory quality" benchmark. Six groups, each tied to a real Howl failure.

### 6.1 Group A — Formation quality (highest priority)

**Corpus:** conversations containing important events among irrelevant dialogue; casual remarks that become significant later; explicit decisions; implicit character development; relationship changes; promises and unresolved commitments; repeated facts expressed differently; contradictory statements; corrections; details that must not become permanent.

**Measures:** identification of important information; rejection of unimportant material; meaning preserved; fact versus speculation separated; who said / witnessed / knew preserved; no invented detail; records concise and readable.

**Comparison arms:** deterministic extraction alone; chunked AI-assisted formation; chunked formation plus reconciliation; transcript-only baseline.

**Purpose:** establish whether formation adds real value, not merely more records.

### 6.2 Group B — Long-conversation handling

**Corpus:** short, medium, long, very long conversations; repeated topics; important information spread across distant sections; multiple chunk sizes and overlaps.

**Key test:** a fact introduced early, modified mid-conversation, resolved near the end. **If the final memory reflects the original fact instead of the resolution, formation has failed.**

### 6.3 Group C — Thinking-budget evaluation

**Compare:** supported configurations at different thinking budgets, including the proposed 2048 where the backend permits; thinking-capable versus non-thinking models.

**Measures:** memory-selection accuracy; important-detail retention; false-memory rate; reasoning and generation time; output consistency; resource consumption; whether extra thinking improves decisions.

**Distinguish three budgets:** the model's reasoning budget, the context supplied to it, and the space reserved for the final answer. Conflating them invalidates the comparison.

### 6.4 Group D — Retrieval and character knowledge

**Method:** ask the same event from several perspectives — what happened; what does Mara know; what does the user know; what does the king believe; what was promised; is it fulfilled.

**Deliberate inclusion:** information present in the database but unavailable to a given character.

**A knowledge leak is a serious failure even when the underlying fact is accurate.**

### 6.5 Group E — Maintenance and lifecycle

**Tests:** duplicate detection; contradiction handling; superseded facts; changed character states; deleted messages; deleted memories; source removal; cache invalidation; reprocessing; interrupted formation; restart after incomplete work; detach and reattach.

**Central test:** resurrection. A deleted memory must not reappear via an old summary, vector index, cache, or stale candidate.

### 6.6 Group F — Actual creative performance (critical acceptance)

**Arms:** transcript only; formation + retrieval; formation + retrieval + relevant original excerpts. Same scenario, model, character definitions, and equivalent budgets.

**Measures:** character consistency; relationship continuity; recall of important events; naturalness; emotional continuity; avoidance of repetitive exposition; ability to continue an unfinished storyline; characters knowing only what they should.

**Blinded evaluators.** This group carries the `E-12` concern from `MEMORY-04`: a system that improves isolated fact recall while degrading the conversation has failed.

**Presentation variants (new, and required by §1.1):** the same memory set is presented in at least three forms — a verbose block, a compact guided block, and a minimal cue set — with the budget held constant. If outcomes differ materially across presentations, then presentation is a first-class variable and no single algorithm result is meaningful without naming the presentation used. This is the direct test of whether Howl is guiding or commanding.

### 6.7 Mapping to the execution protocol

`MEMORY-03` remains the authority for execution detail. The two taxonomies map as follows:

| Test group | Protocol stage |
|---|---|
| A — formation quality | Stage B (reader) + Stage A (selection), plus a formation-specific fixture set |
| B — long conversation | Stage F (flush) + a long-context fixture |
| C — thinking budget | Configuration axis §4.5, run across A and F |
| D — retrieval and knowledge | Stage C (isolation) |
| E — maintenance and lifecycle | Stage D (deletion) + interruption fixtures |
| F — creative performance | **Stage H (new)** — creative acceptance |

**No test group may be reported using another group's measures.** A formation score is not a retrieval score is not a creative score.

### 6.8 Diagnostics: why, not only which

Reporting which presentation won is not enough. Every result is accompanied by the candidate explanation it supports, and the explanation is stated **before** the outcome is inspected.

| Observed outcome | Candidate explanation (hypothesis, not conclusion) |
|---|---|
| Verbose block performs worse | Excessive detail or distraction |
| Compact guided block performs better | Better balance of context and relevance |
| Minimal cues perform worse | Insufficient information to reconstruct relationships |
| Minimal cues perform better | The model already understands the surrounding context |
| All three perform similarly | Presentation may not be the limiting factor |

**Rule:** a presentation result is reported with its diagnostic, and a diagnostic that was not pre-registered is labelled post-hoc. Rank order alone is not a finding.

### 6.9 Failure classification

A single aggregate score hides the only information needed to fix the system. Every failed item is assigned to **exactly one** class, and each class maps to one link of the chain in §1.2.

| Class | Meaning | Detected at | Fix locus |
|---|---|---|---|
| **Never formed** | The information was present but no candidate was produced | Formation coverage against the ground-truth ledger | Formation extraction and salience |
| **Incorrectly formed** | A candidate exists but asserts something the source does not support | Source-span verification | Formation validation; model interpretation |
| **Mis-scoped** | Correct content, wrong character/persona/world/conversation scope | Scope probe: contains it / must not contain it | Scope keys and routing policy |
| **Stored incorrectly** | Correct candidate, but provenance, validity, or kind is lost, or a second truth was created | Schema and provenance audit | Storage and routing policy |
| **Not retrieved** | In scope and retrievable, but not selected for a turn that needed it | Recall against required evidence | Retrieval, ranking, budget |
| **Presented poorly** | Retrieved correctly, but format, position, or budget made it unhelpful | Presentation-variant comparison | Presentation layer |
| **Misunderstood** | Presented correctly, but the reply contradicts or ignores the supplied information | Reply compared to supplied context | Model capability; prompt design |
| **Correctly withheld** | Not presented because staying silent was the right action | Abstention audit | **None — this is success** |

**The last row is not a failure.** Distinguishing *correctly withheld* from *not retrieved* is the entire purpose of the abstention test; a system that conflates them will be tuned to inject more, which is the opposite of the governing principle.

**Required output of every evaluation run:** counts per class, and the class distribution for each arm. A run that reports only an aggregate cannot be acted on.

---

## 7. Acceptance criteria

| Area | Required result |
|---|---|
| Fidelity | No unsupported fact silently promoted to an established memory |
| Coverage | Important information consistently survives formation |
| Conciseness | Meaning retained without reproducing whole conversations |
| Provenance | Every retained record traces to its evidence |
| Identity | No unauthorized character or persona knowledge leakage |
| Temporal accuracy | Changes and superseded information represented correctly |
| Retrieval | Relevant information found without flooding the prompt |
| Deletion | Deleted information cannot reappear through derived data or caches |
| Reliability | Interrupted operations do not corrupt or silently lose committed results |
| Creative quality | Memory-assisted conversations at least as coherent and natural as baseline |
| Guidance | Memory informs interpretation without producing recitation or over-exposition |
| Diagnosis | Every failed item is assigned to exactly one failure class, and every presentation result carries a pre-registered diagnostic |
| Performance | Processing time and resource use reasonable for local use |

**No arbitrary numerical thresholds are set today.** Initial runs establish a baseline and failure severity; explicit thresholds are then set **before** final validation (`MEMORY-03` §9.4).

### 7.1 Absolute blockers

These fail the system regardless of aggregate score:

1. Cross-character or cross-persona knowledge leakage.
2. Memory resurrection after deletion.
3. Unsupported facts presented as established truth.
4. Loss of original evidence.
5. Corruption of conversation history.

**A high average cannot compensate for any of these.** This is the same rule as `MEMORY-03` `D-SAFE-1`–`D-SAFE-5`, extended to the whole architecture.

---

## 8. The testing process

| Phase | Purpose | Exit condition |
|---|---|---|
| **1. Architecture validation** | Establish minimum record requirements, provenance contract, scope rules, temporal behaviour, application responsibilities | These are specified and reviewed; no large implementation commitment |
| **2. Comparative evaluation** | Compare formation methods, thinking budgets, chunking, retrieval, reconciliation on a controlled corpus; include the twelve fixtures from `MEMORY-04` §13, extended as needed | Comparative evidence exists for each open method question |
| **3. End-to-end creative validation** | Real conversations, continuity, long-term retrieval, deletion, recovery; real model behaviour | Acceptance criteria applied; blockers tested |

**Scope discipline:** start with a small number of representative local models and expand only when an unanswered question requires it. No hardware benchmarking for its own sake — latency and resource use are measured inside these tests, and additional performance work happens only if a problem appears.

**This section is a summary.** The canonical pre-registered sequence (`P0`–`P5`), the corpus specification, arm definitions, variable control, rubrics, and the freeze-evidence checklist are in `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md`.

---

## 9. The recommended direction

**One general memory system. Application-owned, AI-assisted, hybrid.**

- Original conversations remain authoritative evidence.
- One formation pipeline produces candidate knowledge records.
- AI-assisted chunked reasoning is an **optional** formation method.
- Separate thinking budgets are investigated as an optimization, not a dependency.
- Deterministic extraction and retrieval remain available.
- Application policy controls validation, classification, scope, and storage.
- Knowledge collections organize reusable information; they do not generate it.
- Retrieval combines lexical, semantic, and historical methods as evidence allows.
- Temporal state, provenance, identity, and deletion are foundational.
- Models never call tools or manage their own memories.

### 9.1 Decisions made

| # | Decision | Status |
|---|---|---|
| 1 | Retire automated lore generation as an independent subsystem; preserve its research and findings | **Made** (`OWNER-DECISION-MEMORY-FORMATION`) |
| 2 | Adopt general memory formation as the intended direction; keep the processing method open to evaluation | **Made** |
| 3 | Retain lorebooks as independent, reusable knowledge collections | **Made** |
| 4 | Make the application responsible for the complete memory lifecycle; AI contributes interpretation, not operational authority | **Made** |

### 9.2 The five governing principles

| # | Principle | Operational meaning |
|---|---|---|
| 1 | **Relevance over volume** | Minimize unnecessary context; every injected token must justify itself |
| 2 | **Guidance over instruction** | Memories inform interpretation; they never become an inflexible script |
| 3 | **Application intelligence over model autonomy** | Howl owns retrieval, identity, organization, validation, and delivery |
| 4 | **Model-appropriate presentation** | Memory formatting and context allocation adapt to the actual model's capabilities (`MEMORY-06`, `MEMORY-07`) |
| 5 | **Natural continuity over mechanical recall** | The test is whether the conversation feels coherent, not whether the model can recite its database |

**Corollary:** no retrieval algorithm and no memory schema is frozen before the full chain in §1.2 has been tested. The research phase ends here; the evidence phase begins.

### 9.3 Deliberately not decided

The final schema; the retrieval algorithm; chunk size and overlap; the reasoning budget; model selection; the store adapter; the ID scheme. **Each must earn its place through testing.** `Q-MF-01` and `Q-MF-10` remain open and are the two questions that would confirm or overturn the formation decision.

---

## 10. The biggest risk

**Confusing more memories with better memory.** It is called out as its own section because it can make or break the project.

A system with 500 records is not better than one with 40. It is usually worse: more duplicates, more contradictions, more irrelevant material, more retrieval noise, more prompt consumption, more opportunities for incorrect continuity.

The standard: **retain the smallest useful representation of the information necessary to preserve continuity, while keeping the original evidence available when more detail is needed.**

Operationally that means record count is never a success metric, and every proposed record must justify its retrieval value.

**The test must not become a contest to maximise recall, memory count, or retrieval volume.** The winning configuration is the one that produces the most accurate, natural, and contextually appropriate continuity **with the least unnecessary interference**. Any ranking that rewards volume is measuring the wrong thing.

### 10.1 Corpus representativeness

Convenient memory examples will flatter any system. The corpus must reflect actual Howl use:

1. Long-running character relationships.
2. Multi-character conversations with differing knowledge boundaries.
3. Scenario and world changes mid-history.
4. Contradictory and superseded information.
5. Important events separated by hundreds of messages.
6. Casual conversations containing one unexpectedly important detail.
7. **Conversations where the correct response is to retrieve nothing.**

**A mandatory baseline:** no configuration is reported without a **transcript-only** arm run on the same material, prompts, and model. Without it, three memory presentations can all be worse than simply letting the model read the conversation, and the result would be misread as a memory improvement.

---

## 11. What this document does not do

- It does not authorize implementation, schema changes, or a runtime change.
- It does not freeze the schema, algorithm, chunk size, budget, store, or model.
- It does not close `Q-MF-01` or `Q-MF-10` by assertion.
- It does not modify any frozen artifact, Amendment 02, or the F3-C contracts.
- It does not replace `MEMORY-03` as the execution protocol; it defines what that protocol must prove.

**The three pillars on which the system ultimately stands: creative continuity, knowledge-boundary correctness, and preservation of historical truth.**
