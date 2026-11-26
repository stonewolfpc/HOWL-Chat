# MEMORY-09 Research Checkpoint — Specialized Memory Capabilities, Retrieval, and Perspective-Aware Context Selection

| Field | Value |
|---|---|
| Record | `MEMORY-09-SPECIALIZED-MEMORY-CHECKPOINT-2026-11-19` |
| Date | 2026-11-19 |
| Status | **Research checkpoint only.** Not an architectural decision, implementation authorization, experiment approval, or binding requirement. |
| Authority created | **None.** Preservation of a research direction and its open questions. |
| Provenance | Checkpoint text supplied by the project owner in session and recorded here without substantive edit. Third-person references to "the user" are the text's own. No claim was added, resolved, promoted or reclassified. |
| Scope | Future research direction for memory capabilities. Production Howl, MEMORY-09 contracts, ratified architecture, frozen requirements and existing decision records are untouched. |
| Registry | No `docs/DOCUMENT-REGISTRY.md` row was added. Registry edits are not made in this workspace without authorization. |

**This file records intent to investigate, not a plan to build.** Naming a candidate
component does not select it. Recording a hypothesis does not establish it. Section
9 governs how the rest is to be read.

---

## 1. Why this checkpoint exists

Following the Round 3 memory-extraction evaluation, we identified a potentially important next research direction: investigating whether specialized, relatively small AI components can improve Howl's memory capabilities when used alongside its existing deterministic memory architecture.

This direction grew partly out of the earlier BGE investigation. BGE should not be considered conclusively rejected. The user believes it was not adequately evaluated and wants to understand whether retrieval components can provide meaningful quality improvements worth their resource costs.

The immediate Round 3 evidence review remains a separate activity. This checkpoint preserves the ideas for the later research stage; it does not replace or interrupt the current evidence review.

## 2. Central architectural observation

Howl's memory system is not simply a mechanism for storing and retrieving facts. It must preserve distinctions about:

- What happened.
- Who experienced it.
- Who directly knows about it.
- Who merely heard about it.
- Who is permitted to know it.
- Who owns the relevant information.
- Whether information belongs to a world, scenario, character, or another scope.
- Whether information is an established fact, an observation, a belief, an interpretation, a rumor, or speculation.
- When information is relevant.
- When information should deliberately be excluded.

This is an epistemic and perspective-management problem, not merely a search problem.

Much of this responsibility is already handled by Howl's deterministic architecture. Research should identify where deterministic mechanisms are sufficient and where AI could provide a measurable benefit. Do not assume AI is needed just because a capability sounds intelligent.

## 3. Research hypotheses worth preserving

These are hypotheses to investigate, not conclusions.

**H1 — Specialized components may outperform a single general-purpose memory model.**

Different capabilities may benefit from different tools or models:

- Retrieval and semantic matching.
- Classification.
- Extraction of events and information.
- Identification of attribution, uncertainty, and temporal context.
- Interpretation of character perspectives.
- Relevance ranking.
- Selection of context for the current conversation.

Potential candidates mentioned include BGE, lexical retrieval, embedding models, rerankers, LFM-family models, Qwen, Gemma, and other small local models. These are candidate examples, not a fixed shortlist or endorsement.

**H2 — Retrieval must be perspective-aware.**

Finding a semantically relevant memory is not enough. The system must retrieve the correct memory for the current character, scenario, world, and knowledge boundary.

A memory can be highly relevant in subject matter but still be inappropriate to expose to a particular character.

Authorization, ownership, visibility, and scope must not be overridden by an AI relevance score.

**H3 — Interpretive memory may deserve separate investigation from factual memory.**

Examples:

- An earthquake altered a mountain range: an event or world fact.
- A character witnessed the earthquake: an observation or experience.
- A character privately thought the earthquake was frightening: a personal interpretation.
- Another character secretly enjoyed it: a different personal interpretation.

Preserving what an event meant to a character may be as important to long-form roleplay continuity as remembering that the event occurred.

However, distinguish information explicitly expressed by a character from interpretations inferred by a model. A model-generated interpretation must not silently become an established character belief or objective fact.

**H4 — A memory editor or context selector may be useful.**

One possible component could sit between the deterministic memory system and the storyteller.

Its job would be to help select which eligible memories are useful for the current turn. It would not own the memory store, rewrite authoritative records, change permissions, or independently establish facts.

This idea is provisional. Investigate whether it provides value beyond deterministic retrieval and ranking before adopting the concept.

## 4. Important distinctions to preserve

Do not collapse these separate research questions:

1. **Extraction:** Can a model identify useful information in conversation?
2. **Interpretation:** Can it identify meaningful character perspectives, beliefs, and emotional significance without inventing them?
3. **Retrieval:** Can a system find the correct stored information?
4. **Perspective-aware retrieval:** Can it find the correct information for the current character and context?
5. **Relevance selection:** Can it choose which eligible memories matter for the current turn?
6. **Authorization:** Is the information actually permitted to be used or disclosed?
7. **Validation:** Can deterministic checks ensure that model suggestions obey the system's contracts?

Success in one capability does not establish success in another. In particular, good semantic retrieval does not prove correct ownership, attribution, scope, or disclosure behavior.

## 5. Preferred research strategy

The user's preferred approach is to investigate capabilities separately first, then investigate combinations.

**Stage A — Isolated capability trials**

Potential trial areas:

- Retrieval: BGE versus lexical retrieval and other embedding/retrieval approaches.
- Reranking: determine whether a specialized reranker improves results beyond initial retrieval.
- Classification: small models identifying information categories and distinctions.
- Extraction: identifying events, observations, attribution, uncertainty, and temporal information.
- Interpretation: identifying character-specific meaning without promoting inference into fact.
- Context selection: selecting relevant memories for a given character and turn.
- Deterministic validation: enforcing scope, ownership, authorization, and state-transition constraints.

Each trial should have an explicit baseline and measure the capability it claims to improve.

**Stage B — Combined capability trials**

Only after isolated capabilities have been evaluated, investigate combinations such as:

- Retriever + reranker.
- Retriever + classifier.
- Retriever + extraction model.
- Classifier + deterministic validation.
- Retrieval + perspective-aware filtering + relevance selection.
- A more complete pipeline combining multiple specialized components.

The purpose is to identify whether combinations produce genuine additional value, not merely more complexity.

**Stage C — System-level evaluation**

Evaluate the complete approach against the deterministic baseline for:

- Memory quality and continuity.
- Correct perspective and attribution.
- False positives and false negatives.
- Unauthorized information exposure.
- Scope contamination.
- Retrieval relevance.
- Latency and throughput.
- RAM and VRAM requirements.
- Resource contention with Howl's other functions.
- Operational complexity and reliability.

A component should earn its place through demonstrated benefit.

## 6. BGE must remain an open candidate

Do not record BGE as rejected or selected.

The user's concern is that it was not adequately tested before being set aside. The intended research direction is to establish a fair comparison:

- Baseline without BGE.
- Matched experiment with BGE.
- Same prompts, corpus, and evaluation conditions wherever applicable.
- Compare retrieval quality and downstream memory quality.
- Measure performance and resource costs.
- Determine whether any quality improvement justifies the additional overhead.

Do not dismiss BGE solely because it consumes resources. Equally, do not assume that better retrieval automatically improves memory quality.

## 7. Architectural boundaries that research must respect

Any proposed AI assistance must be evaluated against Howl's established ownership and authority boundaries.

In particular:

- Howl's chat owns its memory and the state machine operating on that chat's state.
- MEMORY-09 operates within that ownership boundary; it does not own the chat or its lifecycle.
- Authoritative state must not be mutated through unauthorized direct or indirect paths.
- AI-generated suggestions must not silently become authoritative memory.
- Character beliefs must not become world facts.
- Private knowledge must not leak to other characters.
- Scenario-specific information must not become world-wide information.
- Rumor and speculation must not become confirmed truth.
- The conversation layer establishes disclosure events and their participants; MEMORY-09 must not independently infer exposure.

These are existing architectural boundaries, not hypotheses being reopened by this research checkpoint.

## 8. Open questions for the eventual research program

Preserve these questions without prematurely answering them:

1. Which capabilities genuinely need AI, and which are better handled deterministically?
2. Does BGE improve retrieval enough to justify its cost?
3. Does a specialized reranker outperform simpler retrieval methods?
4. Can a small model reliably distinguish facts, observations, beliefs, interpretations, rumors, and speculation?
5. How can inferred character interpretations be represented without turning them into established facts?
6. Can context selection improve continuity without introducing perspective or authorization errors?
7. Should retrieval eligibility be enforced entirely before AI ranking, and what validation is needed afterward?
8. Does a memory editor provide measurable value over a deterministic context-selection baseline?
9. What is the smallest useful combination of components?
10. What are the actual resource and latency costs on the hardware Howl is intended to support?

## 9. Research governance

The user explicitly wants this treated as exploration, not law.

- Do not convert suggestions into requirements without explicit approval.
- Do not assume the proposed components are necessary.
- Challenge ideas when evidence or reasoning warrants it.
- Investigate uncertain claims rather than presenting them as established.
- Preserve competing explanations and negative results.
- Separate observed evidence, source-derived facts, calculations, inferences, unestablished claims, and superseded claims.
- Do not let a scorer or automated evaluator become the authority over human understanding.
- Do not begin implementation merely because a research direction has been recorded.

The goal is to discover what actually improves Howl, not to justify a predetermined architecture.

## 10. Resume point

When this checkpoint is reopened:

1. Confirm the current Round 3 raw-evidence review status and its conclusions.
2. Revisit this checkpoint without treating it as an approved plan.
3. Review the original BGE investigation and why it was set aside.
4. Define isolated capability trials and their baselines.
5. Establish matched evaluation conditions and resource measurements.
6. Only after individual capability results are available, design combination trials.

**Checkpoint conclusion:** The promising direction is not necessarily to build one larger or more intelligent memory model. It may be to combine Howl's existing deterministic ownership and state architecture with a small number of specialized components that each demonstrate a distinct, measurable benefit. That possibility remains open for investigation.

---

## Artifact pointers

Navigation only. Paths are listed because §10 names them as resume inputs; nothing
here characterises their conclusions, and none of them is amended by this file.

| §10 step | Artifact |
|---|---|
| 1 — Round 3 evidence review | `docs/research/MEM-EXTRACTION-ROUND3-2026-11-19.md` |
| 1 — instrument corrections | `docs/research/MEM-EXTRACTION-ROUND3-INSTRUMENT-CORRECTIONS-2026-11-19.md` |
| 1 — independent evidence series | `docs/research/MEM-R3-INDEPENDENT-EVIDENCE/REPORT.md` |
| 3 — BGE evaluation | `docs/research/MEMORY-09-S3-BGE-EVALUATION-2026-11-15.md` (its own header records `Gate B: 6/7 — stays open`) |
| 3 — embedder evaluation | `docs/research/MEMORY-09-S3-EMBEDDER-EVALUATION-2026-11-15.md` |
| 3 — reranker candidates | `docs/research/MEMORY-09-S3-RERANKER-CANDIDATES-2026-11-15.md` |

No experiment was authorised, scheduled or run in producing this record. No model
was loaded, no component was downloaded, and no production file was changed.
