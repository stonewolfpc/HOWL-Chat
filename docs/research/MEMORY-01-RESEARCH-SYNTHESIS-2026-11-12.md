# HOWL Chat — Memory Systems Research and Evidence Reconciliation

**Status:** Research evidence only. No architecture selection, implementation authorization, migration, runtime experiment, benchmark execution, or owner decision is created here.
**Date:** 2026-11-12
**Scope:** Howl's existing memory/continuity/RAG/automated-lore work, external evidence on memory-system approaches, small/local-model constraints, failure modes, and remaining research gaps.
**Revision:** 2 — expanded for the owner research work order of 2026-11-11. New material: §9 (expansion findings and pointers), plus corrections to the model envelope and the automated-lore and context-flush treatment. Sections 1–8 are retained as revision 1 wrote them; where revision 2 supersedes a statement, it is called out explicitly.
**Companions:** `docs/research/MEMORY-02-SOURCE-REGISTER-2026-11-12.md` (sources S1–S45), `MEMORY-04-AUTOMATED-LORE-RESEARCH-ASSESSMENT-2026-11-12.md`, `MEMORY-05-CONTEXT-FLUSH-AND-RECONSTRUCTION-RESEARCH-2026-11-12.md`, `MEMORY-06-PROMPT-FORMAT-COMPATIBILITY-MATRIX-2026-11-12.md`, `MEMORY-07-HARDWARE-AND-INFERENCE-OPTIMIZATION-2026-11-12.md`, `MEMORY-08-COMPARATIVE-RANKING-AND-UNRESOLVED-QUESTIONS-2026-11-12.md`
**Controlling records:** `docs/README.md`, `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md`, `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md`, `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md`, and `docs/review/REQUIREMENTS-RECONCILIATION-HOWL-2026-11-11.md`.

---

## 0. Evidence discipline

This report uses the project's established evidence labels:

| Label | Meaning in this report |
|---|---|
| **Observed** | Reproduced in a controlled environment by this investigation. **Not used**: no runtime experiment was run. |
| **Source-derived** | Stated by a primary paper, benchmark, or identified implementation. |
| **Calculated** | Derived by arithmetic from stated values. |
| **Inferred** | A reasoned applicability judgment, not directly demonstrated. |
| **Unestablished** | No adequate evidence found. |
| **Superseded** | Replaced by a later record. |

Every external result below is **source-derived**. Every statement about Howl's fitness is **inferred** unless it restates a project record. No external benchmark is treated as a Howl benchmark.

**No cross-benchmark ranking is made.** LoCoMo, LongMemEval, DMR, BEIR, RULER, and vendor leaderboards measure different tasks, use different judges, and often use different reader models. A higher number on one is not evidence of superiority on another.

---

## 1. Prior Research Reconciliation

### 1.1 What the project actually has

The project has substantial **design and reconciliation** work. It does not have a memory experiment.

| Artifact | Purpose | What was actually tested | Result / limit | Current use | Open question |
|---|---|---|---|---|---|
| `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md` | Proposed canonical state, event history, memory items, retrieval, deletion, migration | **Nothing.** The document says "Acceptance design — not run"; no model was loaded | A coherent proposed model: append-oriented events, `(continuity_id, character_id)` state, evidence-backed proposals, user locks, as-of queries, tombstones | Design input only; no implementation authorization | `DA-M1`–`DA-M6`, especially store, extraction contract, automation threshold, deletion policy |
| `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md` | Proposed local creative RAG: corpus/rights, structure, lexical/vector indexes, hybrid retrieval, citations | **Nothing.** "No index was built and no model was loaded." | Proposed hybrid retrieval with citation and rights model; acceptance design only | Design input only; no runtime authorization | `DA-R1`–`DA-R6`, especially index/store, embedding identity, reranker, citation granularity |
| `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md` | Requirements for continuity, audience routing, memory/lore deletion provenance | **Nothing.** The research register says "No question in this register has been researched yet." | Requirements and research questions `R-10`, `R-11`, `DA-C6`, `DA-C8`; no answers | Requirements baseline; research questions remain open | Character knowledge boundaries, audience propagation through summaries/memory/RAG, deletion receipts |
| `docs/review/REQUIREMENTS-RECONCILIATION-HOWL-2026-11-11.md` | Reconcile owner requirements against existing capability | Static/source reconciliation only; "No model was loaded, no benchmark was run" | `REC-CH-09` structured memory **partially satisfied**; `REC-CH-10` automated capture **partially satisfied**; `REC-CH-11` rollover **unresolved**; `D7` memory schema/retrieval/compression/reinjection is a dependency | Current requirements evidence | Typed state, conflict behavior, rollover semantics, retrieval/reinjection proof |
| `docs/research/ROOT-01-PRODUCTION-ROOT-OWNERSHIP-AND-CONVERSATION-IDENTITY-2026-11-11.md` | Source investigation of production root/chat identity | Source inspection only; no runtime claim | No root ID; chat IDs collide; memory records have no root/chat field; `SetChatID("")` is a no-op; initial scope is `"default"`; delete/wipe gaps | Current source evidence; defects open | Canonical identity must exist before durable memory can be designed safely |
| `docs/review/HOWL-COMPLETE-PRODUCT-ARCHITECTURE-REQUIREMENTS-DEPENDENCY-MASTER-PLAN-2026-11-11.md` | Whole-product inventory | Source/evidence classification, no new runtime test | S6 memory **source-confirmed incomplete**; S7 RAG **source-confirmed integration gap**; live `App` prompt path does not demonstrably call RAG | Planning inventory | One live prompt path, durable identity, retrieval effect proof |
| `docs/review/CORRECTED-MASTER-PLAN-HOWL-2026-11-11.md` | Proposed corrected plan | Consistency verification only | D-ROOT open; chat/memory roots split; "per-chat" schema unsettled; memory work is premature before one prompt path and identity | Proposed, not adopted | Sequencing and owner decisions |
| `docs/review/F3C-LEGACY-LORE-CLASSIFICATION-2026-11-10.md` | Clarify legacy generated lore populations | Source/test-confirmed classification; no import run | Three populations: canonical-folder legacy lore (prompt-eligible), old name-folder lore (inert), new chat discoveries (chat-owned) | Current classification | No import/migration; no quality measurement |
| `docs/review/ENGINEERING-F2-STORAGE-MODEL-INTELLIGENCE-2026-11-09.md` and `ENGINEERING-F2B-STORAGE-MODEL-INTEGRATION-2026-11-09.md` | Storage/model foundations | Source/unit tests for storage, registry, model intelligence | Content-library and model-intelligence foundations; no memory retrieval or consolidation test | Engineering evidence | Storage foundation is not a memory architecture |
| `docs/review/ASSESSMENT-howl-application-reality-check-shared-foundations-2026-11-09.md` | 15-subsystem reality check | Source inspection | S5 lorebooks/memory/RAG **incomplete**; RAG cannot be shown to affect replies | Historical assessment | Same integration gap persists in later records |

### 1.2 What the project has not established

- No measured memory retrieval quality.
- No measured consolidation/summarization fidelity.
- No measured character-knowledge-boundary behavior.
- No measured cross-character or cross-conversation leakage.
- No measured deletion/forget behavior across derived data.
- No measured automated-lore relevance, noise, token cost, or over-expansion.
- No Howl-specific small-model extraction benchmark.
- No storage/index benchmark on Howl's target hardware.
- No runtime proof that memory or RAG affects the live reply path.

**Source presence is not runtime evidence.** The existence of `memories.json`, `candidates.json`, `memory-store.js`, RAG files, or a prompt string does not establish correct retrieval, correct scope, or correct deletion.

### 1.3 The four existing memory-relevant designs are not experiments

1. **Canonical state + events + memory items.** The design is conceptually strong and directly matches the provenance/correction requirements. It is unbuilt and untested.
2. **Local creative RAG.** The design is careful about rights, citations, index identity, and corpus separation. It is unbuilt and untested.
3. **Audience/character isolation.** Amendment 02 defines the requirement and the leakage criteria. It explicitly has not researched the questions.
4. **Automated lore.** F3-C separates canonical legacy lore from chat-owned generated discoveries. No quality evaluation exists.

### 1.4 Reuse map

This research is evidence for the existing decisions, not a replacement for them:

| Existing open item | What this report supplies | What it does not supply |
|---|---|---|
| `DA-M1` / `DA-01` physical store | Evidence on event/projection and index trade-offs | A store selection |
| `DA-M2` / `DA-02` extraction contract | Small-model extraction/conflict limits | A model choice |
| `DA-M3` / `DA-06` automation threshold | Evidence that explicit store semantics beat model judgment for updates | A field-by-field threshold |
| `DA-M4` continuity branching | Temporal/event evidence; no branch benchmark | A merge semantics decision |
| `DA-M5` deletion/forget | Evidence that deletion and derived-data resurrection are unsolved | A deletion policy |
| `DA-M6` conflict policy | Evidence that last-write-wins and model-only conflict handling fail | A conflict taxonomy |
| `DA-R1`–`DA-R6` RAG | Retrieval, reranking, graph, and citation evidence | An index/product choice |
| `R-10` audience isolation | Direct evidence that post-retrieval filtering leaks | A complete Howl audience model |
| `R-11` provenance/deletion | Evidence on provenance laundering, stale reuse, and deletion gaps | A final provenance schema |

---

## 2. External Evidence Review

### 2.1 The benchmark landscape is not a single ladder

| Benchmark | What it measures | What it does not measure | Key limitation |
|---|---|---|---|
| **LoCoMo** (S2) | QA over long multi-session personal conversations; single-hop, multi-hop, temporal, open, adversarial | Creative continuity, character knowledge boundaries, deletion | Generated/edited dialogues; token-overlap F1 and LLM judges differ across papers |
| **LongMemEval** (S1) | Information extraction, multi-session reasoning, knowledge update, temporal reasoning, abstention | Creative writing, multi-character isolation, provenance UI | 500 questions; LLM judge; histories still finite (115k / 1.5M tokens) |
| **DMR** (used by MemGPT/Zep) | Short conversational memory retrieval | Long-horizon behavior | 60-message conversations; full context is already near-perfect |
| **BEIR** (S13) | Zero-shot retrieval across 18 datasets | Chat memory, personalization, temporal state | Heterogeneous tasks; qrels bias; no single winner |
| **RULER** (S15) | Effective long-context retrieval/reasoning | Memory architecture quality | Synthetic tasks; claimed context ≠ usable context |
| **GateMem** (S10) | Utility + access control + forgetting in shared memory | Creative writing, local desktop constraints | Synthetic institutional episodes; no method solved all three |
| **KBF-QA** (S7) | Character knowledge-boundary fidelity over novels | Product UX, local model performance | New preprint benchmark; no independent replication |

**Consequence:** A vendor number such as "94% on LoCoMo" cannot be compared to "71% on LongMemEval" or to Howl's needs. The correct comparison is always *same task, same reader, same judge, same retrieval budget*.

### 2.2 Structured, entity-based, event-sourced, and temporal memory

**What the approach is.** Store events, entities, facts, and relationships explicitly; derive current state; preserve history; resolve time and conflict at query time.

**Evidence:**

- **APEX-MEM** (S5, ACL 2026) combines a domain-agnostic property graph, append-only event storage, and a multi-tool retrieval agent that resolves conflicting/evolving information. It reports **88.88% on LoCoMo** and **86.2% on LongMemEval**, using frontier QA agents. The architecture directly matches Howl's need for temporal truth and auditable state, but the reported results do not establish small-model performance or that a graph is necessary. **Source-derived; no independent replication.**
- **Zep/Graphiti** (S3) uses a bi-temporal knowledge graph: episodes, semantic facts with validity intervals, and community summaries. On LongMemEval-S with GPT-4o it reports **71.2% vs 60.2% full-context**, at **~1.6k tokens vs ~115k** and **2.58s vs 28.9s**; with GPT-4o-mini, **63.8% vs 55.4%**. On DMR it reports **94.8% vs 94.4% full-context**, a near-tie on a saturation benchmark. Later vendor numbers on its benchmark page (**90.2% LongMemEval, 94.7% LoCoMo**, S23) are not reconciled with the paper's protocol and must be treated as vendor claims.
- **Memora** (S6, ACL 2026 Findings) introduces Forgetting-Aware Memory Accuracy (FAMA), which penalizes reliance on obsolete or invalidated memory. Across four LLMs and six memory agents, it reports **frequent reuse of invalid memories and failures to reconcile evolving memories**. This is the closest peer-reviewed evidence for Howl's "dead/retconned fact keeps resurfacing" risk.

**What this establishes:**

- Explicit event/fact history is a credible way to preserve temporal truth and correction provenance.
- Temporal invalidation is a real, measurable failure mode, not a theoretical concern.
- Structured memory can substantially improve QA over full-context on long histories.

**What it does not establish:**

- That a graph is necessary. APEX-MEM's own ablations report that graph-only retrieval is weaker than the full hybrid tool system.
- That small local models can perform the extraction/ontology work.
- That these systems handle deletion, character perspective, or creative continuity.
- That vendor claims are reproducible.

### 2.3 Episodic, semantic, and consolidation memory

**What the approach is.** Keep experiences (episodic), derive reusable facts/preferences (semantic), and compress/consolidate over time.

**Evidence:**

- **Generative Agents** (S12) uses a memory stream, recency/importance/relevance retrieval, and reflection. It demonstrated believable multi-agent behavior in a 25-agent sandbox. It did **not** measure factual fidelity, cross-character leakage, correction, or deletion. **Source-derived demonstration, not a memory benchmark.**
- **LoCoMo** (S2) reports human QA F1 **87.9**, GPT-4-turbo **32.1**, and Mistral-7B **13.9**. Its summary-based RAG reached high evidence recall but worse answer F1 than observation-level retrieval: summaries retrieved the right session yet lost answer-bearing detail. **Source-derived; exact table values vary by paper version.**
- **LongMemEval** (S1) reports a **30–60%** drop when long-context models read the full history instead of oracle evidence sessions. With Llama-3.1-8B, full-history QA was **0.454 vs 0.710 oracle**; with GPT-4o, **0.606 vs 0.870**. It also reports that **40–50% of errors** were "correct retrieval, wrong generation," higher for weaker readers. **Source-derived.**
- **LightMem** (S16, ICLR 2026) uses sensory/short-term/long-term stages with offline consolidation. It reports up to **10.9% accuracy gain**, **117× token reduction**, **159× fewer API calls**, and **>12× runtime reduction** on LongMemEval with GPT/Qwen backbones. This is efficiency evidence, not proof that a small local model can do the same consolidation.

**What this establishes:**

- Long context alone is not memory; performance degrades even for strong models.
- Summaries are useful as indexes/routing, not as the sole evidence store.
- Consolidation can save tokens but can also lose the exact detail needed to answer.

**What it does not establish:**

- That hierarchical summaries improve creative continuity.
- That a particular consolidation schedule is safe.
- That small models can summarize without introducing contradictions.

### 2.4 Retrieval: lexical, dense, hybrid, reranking, graph, hierarchy

**Evidence:**

- **BEIR** (S13) found no single retrieval architecture wins across 18 zero-shot datasets. BM25 is a strong, cheap zero-shot baseline; a cross-encoder reranker improved average nDCG by about **11%** over BM25 while costing roughly **450 ms CPU/query** in the reported setup. Dense retrievers are cheaper at query time but domain-fragile.
- **LongMemEval** (S1) found that **index design** mattered more than a simple dense-vs-lexical choice on chat memory: appending extracted user facts to keys improved recall by **9.4%** and final accuracy by **5.4%** on average. It also found that rank-merging multiple retrieval pathways underperformed merging keys at index time. **This is a direct warning against "just add RRF."**
- **RAPTOR** (S14) reports gains from recursive abstractive trees, including a large QuALITY gain with GPT-4. The comparison changes the reader, so the retrieval-effect size is not clean. **Source-derived; contested.**
- **GraphRAG** (S27, Microsoft technical report) uses LLM-generated entity graphs and community summaries for query-focused summarization. It reports wins on comprehensiveness/diversity with LLM judges, but no retrieval recall/nDCG; indexing is LLM-heavy. For a local app with small models, the construction cost and extraction fragility are material. **Vendor technical report; not independent.**
- **Zep** and **APEX-MEM** show that hybrid graph + lexical + semantic retrieval can work well when the graph is high quality. They do not show that a graph is the best default for a small personal corpus.

**What this establishes:**

- Lexical retrieval is a mandatory baseline, especially for names, places, invented terms, and quoted phrases.
- Hybrid retrieval is credible for coverage, but its ranking benefit is corpus-dependent.
- Reranking is the highest-leverage optional compute in the evidence, but it must be budgeted.
- Graph/hierarchical structures are not free wins; they add extraction and token costs.

### 2.5 Agent-memory frameworks

| System | Architecture | Published evidence | Independent evidence / caveat |
|---|---|---|---|
| **MemGPT / Letta** (S11) | OS-style virtual context; main/recall/archival memory; model edits memory via tools | DMR **93.4%** (GPT-4-turbo) in the MemGPT paper; full-context **94.4%** in Zep's comparison | DMR is short and saturated; tool-call discipline is a failure mode |
| **Mem0** (S4) | LLM extraction + ADD/UPDATE/DELETE/NOOP consolidation over a vector store; optional graph | LoCoMo **66.88%** J vs **72.90%** full-context; graph **68.44%**; 1,764 vs 26,031 tokens | In its own table, Mem0 loses to the trivial full-context baseline. Later managed-platform claims (S24: **94.4% LongMemEval, 92.5% LoCoMo**) are a different product/protocol. |
| **Zep / Graphiti** (S3) | Bi-temporal graph + hybrid retrieval | LongMemEval **71.2% vs 60.2%** (GPT-4o); DMR near-tie | Later vendor claims (S23) are unreconciled; graph DB + strong extraction model required |
| **LangMem** (S25) | Store-first primitives; hot-path/background memory tools | **No vendor benchmark** | All circulating numbers are third-party or competitor re-runs; lowest infrastructure footprint |
| **Cognee** (S26) | ECL pipeline: relational + vector + graph + summaries | Vendor reports **~91% LoCoMo**, **0.79 BEAM-100K**; vendor says "not neutral" | 10M run not reproducible; local GLiNER path exists but its accuracy is not publicly established |
| **A-Mem** | Zettelkasten-style notes, links, memory evolution | LoCoMo gains in the paper; Mem0's re-run reports **48.38%** J, below full context | Metric mismatch (F1/BLEU vs J); three LLM decisions per memory compound error |

**What this establishes:**

- Framework headline numbers are not portable. The same system can score ~32–83% on "LongMemEval" across harnesses.
- Full-context and plain retrieval are strong baselines that many memory products omit or lose to.
- Operational complexity varies enormously: LangMem is a library; Zep/Graphiti is a graph database plus extraction pipeline.

**What it does not establish:**

- That any framework is suitable for Howl's local, single-user, multi-character, offline constraints.
- That vendor-managed numbers apply to the open-source/local build.
- That a framework solves deletion, provenance, or character boundaries.

### 2.6 Character, multi-character, and multi-principal isolation

This is the most directly relevant evidence for Howl's multi-character requirement.

- **ReverieMem** (S7, arXiv preprint) names **Factual Overreach**: shared retrieval or parametric memory lets a character use facts outside its perspective. It proposes first-person episodic memory, **visibility-tagged semantic facts**, and a presence roster. On its KBF-QA benchmark (4,386 questions over eight novels), it reports **+34.6 percentage points** knowledge-boundary fidelity over the strongest prior method and **~79%** win rate on BookWorld's narrative protocol. **Preprint; no code; no independent replication.**
- **Character-LLM** (S20, EMNLP 2023) fine-tunes a separate agent per role to avoid cross-role knowledge collision. It supports the isolation principle but is impractical for many dynamic characters and does not handle entry/exit/death.
- **MemLeak** (S8, NeurIPS 2026 workshop) reports **70–100% incidental cross-user leakage** under pooled vector retrieval, with **90–100%** top-k placement for crafted memories. It found **metadata filtering alone insufficient**; only hard post-retrieval ownership gating consistently restored the clean baseline. **Synthetic, small-n; direction robust, magnitudes indicative.**
- **Authorization-First Retrieval** (S9, TrustNLP 2026) reports **86.1%** structural exposure when retrieval happens before filtering, **0%** structural exposure when authorization constrains the candidate set first. Behavioral guardrails still leaked at **41.3%/29.5%** depending on model. **Peer-reviewed workshop paper; synthetic corpus.**
- **GateMem** (S10, arXiv preprint) evaluates utility, access control, and active forgetting together. It reports that **no method simultaneously achieved strong utility, robust access control, and reliable forgetting**; long-context prompting had the best governance at high token cost, while retrieval/external-memory methods still leaked unauthorized or deleted information.
- **Memora** (S6) reports frequent reuse of invalidated memories, which maps directly to dead/retconned characters and changed world facts.

**What this establishes:**

- Character knowledge must be a **retrieval authorization boundary**, not a prompt instruction.
- Scope must be applied **before** similarity search or generation.
- Metadata-only filtering is not a reliable control.
- Deletion and forgetting are not solved by current memory systems.
- Presence/lifecycle is a real modeling problem, not just a UI state.

**What it does not establish:**

- That these results transfer quantitatively to a single-user creative-writing app.
- That any product implements a character-presence model comparable to Howl's.
- That hard partitioning alone handles spoilers from the model's parametric knowledge or author notes.

### 2.7 Small and local models

Detailed in §4. The headline evidence:

- **LongMemEval** (S1): Llama-3.1-8B's performance drops sharply beyond **3k retrieved tokens**; GPT-4o continues improving past **20k**. At top-10 on LongMemEval-M, 8B scored **0.572** vs GPT-4o **0.720**.
- **RULER** (S15): effective context is often far below the claimed window; aggregation and multi-hop tracing degrade most.
- **When Facts Change** (S18): smaller models often fail to detect a temporal conflict; larger models may detect it but fail to act on mutability.
- **Memory-R1** (S17): a learned memory manager (ADD/UPDATE/DELETE/NOOP) can improve memory behavior at 3B–14B scale, but training/evaluation is benchmark-specific and not replicated.
- **LightMem** (S16): consolidation efficiency is real, but reported with GPT/Qwen backbones; it does not prove a 7B local model can do the same.

### 2.8 Automated lore is a separate research question

Existing project evidence:

- F3-C's `F3C-LEGACY-LORE-CLASSIFICATION` separates three populations:
  1. legacy automated files already inside a canonical item folder — read-only, `legacy_generated_unattributed`, prompt-eligible;
  2. legacy automated files in old name-keyed folders — inert, only LF-1 migration may read them;
  3. newly generated conversation discovery — chat-owned, no canonical origin.
- No import has run. No quality measurement exists.
- The existing automated-lore behavior is reported to have relevance/over-expansion concerns, but the project has not quantified them.

External evidence relevant to automated lore:

- Summary-based retrieval can retrieve the right session yet lose the answer-bearing detail (S2).
- LongMemEval shows that summaries/facts can help as index keys but are not a substitute for raw values (S1).
- Graph/consolidation pipelines add extraction cost and can propagate errors (S5, S3, S26).
- No external benchmark measures Howl-style automated lore.

**Disposition options — none selected:**

| Option | Evidence for | Evidence against | Current judgment |
|---|---|---|---|
| Retain as-is | It already exists and is prompt-eligible in one population | No quality/token measurement; known noise risk | **Unestablished** |
| Rework as source-linked candidates | Aligns with memory design's evidence-backed proposals | No Howl test; adds review burden | **Plausible, unproven** |
| Integrate with memory candidate pipeline | Avoids two extraction/review systems | Risks conflating lore with character/world state | **Plausible, needs schema** |
| Make optional | User control; reduces noise when off | No measured benefit/harm | **Reasonable, unproven** |
| Replace with retrieval over raw text | Summary evidence favors raw evidence | More tokens; may be slower | **Plausible, needs budget test** |
| Retire | Removes a noisy subsystem | No evidence it is net-negative; may remove useful recall | **Not supported** |

The evidence supports treating automated lore as a **separate, low-trust, source-linked retrieval/index layer**, not as canonical state. It does not establish which disposition is best.

---

## 3. Comparative Evaluation

**No numerical ranking.** The table compares evidence strength and fit, not benchmark scores.

| Approach | Evidence strength | Strengths | Weaknesses | Howl fit | Cost / risk |
|---|---|---|---|---|---|
| **Append-only event log + derived state** | Moderate (APEX, Zep, Memora) | Provenance, correction, as-of queries, deletion receipts, temporal truth | Requires schema, migration, write discipline; storage growth | **High** for canon/provenance | Medium–high design cost; low runtime cost |
| **Entity/fact graph** | Moderate | Relationships, multi-hop, temporal invalidation | Extraction errors, entity resolution, graph maintenance | **Medium** | High extraction/index cost; small-model risk |
| **Episodic raw retrieval** | Moderate (LongMemEval, LoCoMo) | Exact evidence, quotes, source spans; no summary loss | Noise, token cost, ranking difficulty | **High** as evidence layer | Medium storage; retrieval tuning |
| **Semantic facts/profiles** | Moderate | Compact personalization, fast lookup | Staleness, overgeneralization, provenance loss | **Medium–high** as derived layer | Extraction cost; must remain derived |
| **Hierarchical summaries** | Limited/contested (RAPTOR, LoCoMo, LongMemEval) | Cheap routing, multi-session overview | Detail loss, drift, hallucination amplification | **Low as sole memory; medium as index** | Summarizer cost; fidelity risk |
| **Lexical/BM25** | Strong (BEIR, LongMemEval) | Cheap, exact names/terms, no model | Vocabulary mismatch, no semantics | **High baseline** | Low |
| **Dense embeddings** | Moderate | Semantic recall, multilingual potential | Domain fragility, embedding identity, local model cost | **Medium** | Medium; index versioning |
| **Hybrid + rerank** | Moderate | Coverage + precision; exact + semantic | Latency, tuning, corpus-dependent gains | **Medium–high** | Medium; reranker CPU/GPU |
| **GraphRAG / community summaries** | Limited for local | Global/thematic queries | LLM-heavy indexing, huge query tokens, extraction fragility | **Low as default** | High |
| **Agent-memory frameworks** | Mixed | Fast start, built-in extraction/retrieval | Vendor numbers not portable; operational complexity; deletion/scope gaps | **Variable; not a shortcut** | High integration risk |
| **RL-trained memory manager** | Emerging (Memory-R1) | Learns ADD/UPDATE/DELETE/NOOP policy | Training data, benchmark-specific, no independent replication | **Unknown** | High research cost |

### 3.1 What the evidence supports as constraints

1. **Raw evidence must remain retrievable.** Summaries alone lose detail (S1, S2).
2. **Current state must be a projection, not an overwrite.** Temporal/correction evidence (S3, S5, S6).
3. **Scope must be enforced before retrieval.** Character/multi-principal evidence (S7–S10).
4. **Deletion must cover derived data and indexes.** GateMem (S10), Memora (S6); no complete solution.
5. **Retrieval must be small and precise for weak readers.** LongMemEval (S1); 3k budget for 8B.
6. **Lexical retrieval must remain a baseline.** BEIR (S13); exact names/terms.
7. **Full-context and plain retrieval must be baselines in any Howl evaluation.** Mem0 (S4), Zep (S3), vendor divergence.
8. **No framework should be adopted from headline numbers.** Independent harness divergence (S4, S3, S24).

### 3.2 What the evidence does not support

- Selecting a graph database, vector database, or agent framework as the architecture.
- Assuming a 7B/8B local model can do frontier-style extraction, conflict resolution, or temporal reasoning.
- Assuming summaries are safe because they are shorter.
- Assuming metadata filtering is a security boundary.
- Assuming a benchmark win on LoCoMo/LongMemEval means better creative continuity.
- Assuming deletion is solved because a delete API exists.

---

## 4. Small/Local Model Assessment

### 4.1 Demonstrated limits

| Capability | Evidence | Implication for Howl |
|---|---|---|
| **Long-context reading** | LongMemEval (S1): 30–60% drop vs oracle; 1B/3B near-collapse on full history; 8B cliff at ~3k retrieved tokens | Do not rely on a large context window; retrieve a small, precise package |
| **Answer synthesis over retrieved evidence** | LongMemEval (S1): 40–50% of errors are retrieval-correct/generation-wrong, worse for 8B | Retrieval improvements cannot fix reader errors; evaluate end-to-end |
| **Structured extraction** | LongMemEval uses Llama-3.1-8B for extraction successfully, but extraction quality caps gains; other small-model extraction evidence is mixed | Use narrow schemas, deterministic validation, and review; do not assume open-ended extraction |
| **Entity resolution** | Reported 14B near-saturation vs 8B false positives (small-model brief); cross-script/ambiguous pairs remain hard | Stable IDs and conservative merges; quarantine ambiguity |
| **Conflict detection** | When Facts Change (S18): smaller models often miss conflicts; larger models detect but fail to act | Put update semantics in the store; do not ask the model to be the conflict authority |
| **Temporal query inference** | LongMemEval (S1): 8B often invents time ranges and prunes the answer | Do not hard-filter on model-inferred time for ambiguous queries |
| **Summarization fidelity** | LoCoMo (S2): summaries retrieve but lose answer-bearing detail | Keep raw turns; summaries as index keys |
| **Abstention** | LongMemEval/conflict evidence: weaker models over-answer | Explicit "insufficient evidence" path and UI |
| **Tool-call discipline** | MemGPT/Letta: memory management depends on tool calls | Avoid requiring a weak model to manage its own memory via tools |

### 4.2 What the architecture can offload

| Function | Offloadable? | Residual model need |
|---|---|---|
| Storage, IDs, timestamps, turn/session granularity | **Yes, fully** | None |
| Scope/authorization filters | **Yes, fully** | None if enforced in the store |
| Schema validation / JSON well-formedness | **Yes** | Semantic correctness still required |
| Candidate generation / blocking / lexical+dense retrieval | **Yes** | Ranking judgment remains |
| Deduplication by stable identity | **Yes** | Ambiguous merges need review |
| Update semantics (supersede/retract/tombstone) | **Yes, structurally** | Detecting implicit contradiction remains model/human work |
| Summarization | **Partly** | Fidelity judgment; source re-injection |
| Conflict detection | **No, not reliably** | Stronger model, verifier, or human review |
| Final answer synthesis | **No** | Reader capability |
| Abstention calibration | **No** | Model or explicit policy |

### 4.3 Howl-specific unknowns

- No Howl model/hardware baseline exists. The project cannot yet say whether "small" means 3B, 7B, 8B, 14B, or a quantized variant on the target GPU/CPU.
- No Howl prompt/token budget is validated against a tokenizer.
- No Howl retrieval corpus statistics exist (turns per chat, sessions, tokens, entity density).
- No Howl character-boundary test set exists.
- No Howl deletion/forget test exists.

**Inference:** The safest architecture is one that makes the model's job small: deterministic scope, compact state cards, a few high-precision evidence passages, explicit "unknown", and no requirement that the model maintain its own memory bookkeeping. This is an inference, not a measured Howl result.

---

## 5. Failure-Mode Analysis

| Failure | Evidence | Howl relevance | Mitigation evidence | Residual risk |
|---|---|---|---|---|
| **Irrelevant retrieval / false association** | BEIR (S13); LongMemEval (S1) | Character/world facts can be pulled by semantic similarity alone | Lexical baseline + rerank + scope filters | Medium; no Howl measurement |
| **Summary omission / drift** | LoCoMo (S2); LongMemEval (S1) | A summary can lose the exact promise, injury, or clothing detail | Keep raw evidence; summaries as index only | High if summaries are sole memory |
| **Recursive summarization hallucination** | Hierarchical summarization evidence (S14 and later work) | Long arcs can accumulate false canon | Re-inject source text; review diffs | High; no Howl test |
| **Stale fact reuse after change** | Memora (S6); LongMemEval knowledge-update | Dead characters, retcons, changed relationships | Validity intervals, supersession, tombstones | High; no complete solution |
| **Cross-character knowledge leakage** | ReverieMem (S7); MemLeak (S8); AFR (S9) | A character can use another character's private knowledge or future plot | Visibility tags + presence roster + retrieval-before-generation gate | Medium–high; direct evidence, no Howl replication |
| **Cross-conversation contamination** | MemLeak (S8); GateMem (S10) | A fact from one story can appear in another | Hard partition by world/chat/character; no pooled store | Medium; design choice |
| **Deleted/forgotten data resurrection** | GateMem (S10); Memora (S6) | Deleting a message may leave summaries/embeddings/caches | Tombstones + derived-data invalidation + cache guard | High; unsolved |
| **Provenance laundering** | No direct benchmark; related provenance/poisoning evidence (S3, S5, S21, S22) | A model summary can turn an external claim into "user history" | Immutable source links outside the model's rewrite path | High; no complete solution |
| **Memory poisoning** | AgentPoison (S21); MINJA (S22) | Imported documents or character cards can inject persistent instructions | Trust labels, quarantine, provenance, write review | High; local single-user still imports untrusted content |
| **Index staleness / embedding drift** | LongMemEval (S1); RAG design | Reindexing can mix models/dimensions | Index manifest + embedding identity + rebuild | Medium; engineering-solvable |
| **Context budget overflow** | LongMemEval (S1) | A weak model degrades sharply beyond ~3k retrieved tokens | Hard budget + reserved response + small candidate set | Medium; needs Howl calibration |
| **Over-trust in benchmark numbers** | Mem0 (S4), Zep (S3), vendor divergence | A framework can look good on a benchmark and fail on Howl | Same-task, same-reader, same-judge evaluation | High; evaluation discipline required |
| **Human correction burden** | No direct benchmark | Automatic capture can flood review | Safe auto-commit only for objective reversible facts; review queue | Unknown; UX research needed |

---

## 6. Research Gap Register

| ID | Gap | Why it matters | Evidence needed | Priority |
|---|---|---|---|---|
| `G-01` | No Howl creative-continuity benchmark | LoCoMo/LongMemEval are QA benchmarks, not long-form creative continuity | Howl-specific scenarios: callbacks, tone, character voice, retcons, off-page events | **P0** |
| `G-02` | No Howl character-boundary test | ReverieMem is direct but preprint; no Howl data | Should-know/must-not-know pairs per character; spoiler/future-plot probes | **P0** |
| `G-03` | No Howl deletion/forget test | GateMem shows no system solves utility + access + forgetting | Delete a message/character/fact; probe raw, summary, embedding, cache, prompt | **P0** |
| `G-04` | No canonical identity prerequisite | ROOT-01: no root ID, colliding chat IDs, memory scope gaps | F3-C root/chat identity resolution | **P0 blocker** |
| `G-05` | No small-model extraction benchmark on Howl schemas | 8B vs 14B extraction quality varies; no Howl schema | Fixed schema, target models, error taxonomy, review cost | **P1** |
| `G-06` | No retrieval evaluation on personal chat scale | BEIR/LongMemEval are not Howl corpora | Recall/nDCG on labelled Howl evidence; lexical vs dense vs hybrid | **P1** |
| `G-07` | No consolidation fidelity measurement | Summaries can lose detail | Atomic-fact retention, contradiction rate, token savings, downstream QA | **P1** |
| `G-08` | No automated-lore quality measurement | Existing behavior has known noise/over-expansion concerns | Precision/relevance/noise/token cost with lore on/off | **P1** |
| `G-09` | No local storage/index benchmark | SQLite FTS/vector choices affect latency/cost | Build/read/update/delete latency; index size; crash recovery | **P1** |
| `G-10` | No temporal/conflict test for creative canon | QA benchmarks do not capture retcons/death/relationship changes | Scripted changes; as-of queries; stale-reuse probes | **P1** |
| `G-11` | No cross-character leakage test in one process | MemLeak/AFR are multi-tenant; Howl is single-user but multi-character | Adversarial character probes; derived-data propagation | **P1** |
| `G-12` | No provenance UI/usability study | Users must inspect/correct memory | Task-based study; trust/calibration; review burden | **P2** |
| `G-13` | No graph/RAG cost-benefit on Howl corpus | GraphRAG is expensive; no local evidence | Indexing time, query tokens, quality gain vs hybrid baseline | **P2** |
| `G-14` | No poisoning threat model for local imports | Character cards/documents are untrusted | Threat model; quarantine; write-review evidence | **P2** |
| `G-15` | No independent replication of most 2026 results | Vendor/preprint claims dominate | Reproduction or independent benchmark | **P2** |

---

## 7. Executive Synthesis

### 7.1 What the evidence establishes

1. **Long context is not memory.** Even frontier models lose 30–60% when reading full histories; small models collapse. Retrieval is necessary, not optional. (S1, S15, S19)
2. **Summary-only memory is unsafe.** Summaries can retrieve the right session and still lose the answer-bearing detail. Raw evidence must remain retrievable. (S2, S1)
3. **Temporal truth requires explicit store semantics.** Validity intervals, supersession, retraction, and tombstones are the only credible way to handle changing facts; model judgment alone is insufficient. (S3, S5, S6, S18)
4. **Character scope is an authorization boundary.** Shared retrieval leaks; post-retrieval metadata filtering is insufficient; scope must constrain the candidate set before generation. (S7–S10)
5. **Deletion/forgetting is unsolved.** No evaluated system simultaneously achieves strong utility, access control, and reliable forgetting. (S10)
6. **Small models need small, precise contexts.** The 8B cliff is around 3k retrieved tokens in LongMemEval; 1B/3B models are near-useless on full history. (S1)
7. **Framework headline numbers are not portable.** Same-system scores vary wildly across harnesses; full-context and plain retrieval are strong baselines. (S3, S4, S23, S24)
8. **Automated lore is not the same problem as memory.** It is a separate, low-trust generated-content population with no Howl quality measurement. (F3C classification; S1, S2)

### 7.2 What the evidence does not establish

- That any specific architecture (event store, graph, vector-only, hybrid, framework) is the right choice for Howl.
- That a particular local model size can perform the required extraction/conflict work.
- That LoCoMo/LongMemEval gains transfer to creative writing.
- That any framework solves character boundaries, deletion, or provenance.
- That automated lore should be retained, reworked, made optional, replaced, or retired.
- That the 2026 preprint/vendor results are independently reproducible.

### 7.3 Questions now ready for owner consideration

These are questions, not decisions:

1. **Capture policy:** Which memory writes may be automatic, and which require review? The design proposes safe reversible facts auto-commit; the evidence supports caution but does not set the threshold.
2. **Character epistemology:** Should characters be limited to what they witnessed/were told, or may the author grant omniscient knowledge? ReverieMem's evidence favors visibility tags and presence rosters.
3. **Deletion semantics:** What must "forget" remove: raw event, derived fact, summary, embedding, cache, citation, and prompt history? GateMem shows this is not solved by ordinary retrieval.
4. **Automated lore:** Is it a optional recall aid, a candidate generator, or a retired feature? No quality data exists.
5. **Evaluation envelope:** What model sizes and hardware define "small/local" for Howl? Without this, external small-model evidence cannot be mapped to a target.

### 7.4 What should be researched next — not implemented

A bounded research protocol, if separately authorized, should compare:

- **E1:** raw event/turn retrieval + lexical baseline vs hybrid lexical+dense vs hybrid+rerank.
- **E2:** hard character/chat scope before retrieval vs pooled retrieval + metadata filter.
- **E3:** summary-as-index + raw evidence vs summary-only memory.
- **E4:** explicit validity/supersession vs model-judged last-write-wins.
- **E5:** automated lore on/off/reworked as candidates, measured for relevance, noise, tokens, and user correction burden.
- **E6:** the same tasks across the actual target local models, not a frontier proxy.

No runtime experiment is authorized by this report. The protocol design is recorded in `docs/research/MEMORY-03-EVALUATION-PROTOCOL-2026-11-12.md`; it is a research proposal only and does not authorize execution.

---

## 8. Placement and Traceability

| Record | Role |
|---|---|
| `docs/README.md` | Master navigation; current memory-research entry point |
| `docs/DOCUMENT-REGISTRY.md` | Registry row and hash for this report and the source register |
| `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md` | Proposed memory/continuity design; open decisions `DA-M1`–`DA-M6` |
| `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md` | Proposed local RAG design; open decisions `DA-R1`–`DA-R6` |
| `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md` | Requirements `R-10`, `R-11`; audience and deletion provenance |
| `docs/review/REQUIREMENTS-RECONCILIATION-HOWL-2026-11-11.md` | `REC-CH-09`, `REC-CH-10`, `REC-CH-11`; dependency `D7` |
| `docs/research/ROOT-01-PRODUCTION-ROOT-OWNERSHIP-AND-CONVERSATION-IDENTITY-2026-11-11.md` | Identity/root defects that block durable memory |
| `docs/review/F3C-LEGACY-LORE-CLASSIFICATION-2026-11-10.md` | Automated-lore populations |
| `docs/review/HOWL-COMPLETE-PRODUCT-ARCHITECTURE-REQUIREMENTS-DEPENDENCY-MASTER-PLAN-2026-11-11.md` | S6/S7 current-state classification |
| `docs/research/MEMORY-02-SOURCE-REGISTER-2026-11-12.md` | Full source metadata and limitations |
| `docs/research/MEMORY-03-EVALUATION-PROTOCOL-2026-11-12.md` | Pre-registered evaluation protocol; design only, no execution authorization |

**Authority:** This is research evidence. It does not alter a frozen contract, owner decision, amendment, or governance authority. If an existing requirement conflicts with this evidence, the conflict is recorded here; the requirement is not silently changed.

**No frozen material was edited. No source code, schema, runtime configuration, database, or user data was changed. No model was loaded. No benchmark was run.**

---

## 9. Revision 2 — expansion findings and pointers

Revision 1 answered the synthesis question for memory and retrieval. Revision 2 adds five domains the work order requires, and corrects two framings that revision 1 stated too narrowly. Detailed work lives in the companion documents; this section records the conclusions and the pointer.

### 9.1 Correction: the model envelope

Revision 1 §4 discussed "small/local models" and the stage-gated matrix sampled `T-S`/`T-M`/`T-L` (3–4B, 7–8B, 12–14B). Read alone, that could imply Howl's supported envelope is 3–14B. **It is not, and revision 1 did not intend it to be.**

- Howl's objective is **broad local-model support, including large models where hardware permits**.
- Tiers are **experimental sampling categories**, not a support boundary, a default, or a recommendation.
- Users choose the model and may override resource policy.
- No model is excluded by parameter count.

`MEMORY-03` §4 now states this explicitly and adds a fourth optional tier (`T-D`) for the largest model the owner's hardware can host, so the protocol cannot be read as capping the envelope.

### 9.2 Correction: automated lore is not a memory feature

Revision 1 §2.8 placed automated lore inside the memory discussion. Revision 2 treats it as a **separate, low-trust generated-content population** with its own mechanism evidence, failure taxonomy, and value ledger (`MEMORY-04`). Three findings change the shape of the question:

1. Every mature local creative tool that has lorebooks uses keyword activation with token budgets, and most add recursion and probabilistic activation (S38–S40). **None publishes a precision/recall measurement.**
2. Budget exhaustion is a **silent** failure in all of them: entries vanish without an error.
3. The one documented strength of lore is **exact detail**, which is exactly what summaries lose (S2). That points at a scope reduction (detail assistance with provenance and bounds), not at retention-as-is or retirement.

**Revision 2 of `MEMORY-04`** now carries the full research package required by the automated-lorebook phase: eight tracks (generation, maintenance, activation, scope and knowledge boundaries, local-model suitability, existing systems, orchestration, deletion and cache integrity), a lore-specific ranking, a cooperation matrix, and twelve pre-specified evaluation fixtures (`E-01`–`E-12`). Two of its findings matter to this synthesis: structured output is an **application-enforced** property (grammar-constrained decoding guarantees parseability but not meaning, S49/S50), and the documented product trade-off between rule-based and embedding-based activation is explicitly one of predictability versus flexibility (S48).

**Revision 3 supersedes one part of that framing.** Per `docs/architecture/OWNER-DECISION-MEMORY-FORMATION-2026-11-12.md`, automated lore-entry **generation** is no longer an independent subsystem: Howl has one memory-formation pipeline with multiple destinations, and a lorebook is one destination (an organization layer), not the mechanism that creates knowledge. The activation, maintenance, scope, and deletion findings in `MEMORY-04` survive unchanged and now apply to every memory record. See `MEMORY-08` for the adjusted `A13` entry.

No disposition is selected. `H8`'s four options remain open, and retirement is not supported by the current evidence.

### 9.3 New: context flushing as a maintenance boundary

Revision 1 treated rollover as a context-management topic. Revision 2 treats the flush as the **predictable commit point** for memory maintenance (`MEMORY-05`). Six stages: pre-flush capture, event extraction, state update, consolidation, reconstruction, verification. The key measures are `silent_loss_count` and `reconstruction_coverage`, not token savings.

Two verified facts constrain any design:

- Context shift is **off by default** upstream since PR #15416 (S28, S32), so a design that assumes shift is measuring nothing.
- Current Howl behaviour is character-based truncation at `contextSize × 4`, keeping the system prompt and newest turns. It is not the `memory-guided` rollover the settings expose.

`P9` (flush-aware consolidation beats plain truncation) is **unestablished**: no external benchmark isolates it.

### 9.4 New: prompt-format fidelity is a confound, not a detail

Revision 1 assumed a single prompt scaffold. `MEMORY-06` shows that format is a first-class variable:

- ChatML, Llama 2/3, Mistral, Qwen2/3, DeepSeek, Gemma, Command-R, Phi, GLM, InternLM, and Yi do not share role markers, turn terminators, or system-role rules.
- **Gemma IT supports only `user` and `model` roles; a system role is unsupported** (S33). **DeepSeek has no system token**; system text is folded into a flat prefix (S34). Command-R prescribes a structured preamble (S35). Qwen3's thinking toggle is **not exposed in llama.cpp** (S36).
- Howl has family-aware **sampling** but generic **formatting**: the legacy builder hardcodes ChatML, and `--chat-template` is passed **without `--jinja`**, which upstream documents as the condition for accepting a full Jinja template (S28). That is a probable defect requiring verification, not a confirmed runtime failure.

`P10` (format fidelity changes answer quality enough to invalidate cross-format comparisons) is **unestablished** but is now a protocol axis (`F0`–`F2`, Stage G).

### 9.5 New: hardware is a precondition, not an afterthought

Revision 1 had no hardware research. `MEMORY-07` adds it, with four verified corrections to common assumptions:

| Assumption | Verified fact | Source |
|---|---|---|
| `--fit-target` is a percentage of VRAM | Fixed **1024 MiB** default | S32 |
| Context shift happens by default | **Off** by default since PR #15416 | S28, S32 |
| KV-cache quantization is independent of attention | **Requires Flash Attention** | S32 |
| A 10% reserve is an engine behaviour | It is a **Howl policy**, not an engine constant | S32 + policy |

**Source-level defect found:** Howl's current settings pair `kv_cache_type: flash-attention` with `flash_attention: false`, and leave `kv_cache_size: 512` without a unit. Any benchmark run on that configuration cannot be attributed to a known memory policy. No universal OOM-prevention claim is admissible.

### 9.6 New: the ranking, and what it is for

`MEMORY-08` ranks eighteen approaches on thirteen dimensions with seven evidence classes. Summary of the order:

| Band | Approaches |
|---|---|
| Top | `A1` event log + derived state with validity/provenance; `A11` structural scope-before-retrieval; `A3` lexical baseline; `A5` hybrid + rerank; `A6` summary-as-index with raw evidence |
| Middle | `A12` flush consolidation; `A14` format fidelity; `A15` residency/reserve; `A4` dense retrieval; `A13` automated lore (conditional) |
| Lower | `A2` graph; `A8` hierarchical summaries; `A16` tool-mediated memory; `A9` frameworks; `A18` full context alone; `A17` GraphRAG default; `A10` RL memory manager |
| Eliminated | `A7` summary-only as a complete solution; post-filtering as a security control |

**The ranking is not a selection.** It is a test order that the protocol is designed to be able to overturn (§7 of `MEMORY-08`). Neither a framework nor a custom build is preferred by it: no existing framework has been shown to satisfy `H1`–`H8` together, and none has been shown incapable.

### 9.7 What revision 2 does not change

- No architecture is selected.
- No owner decision is closed.
- No threshold is frozen; `MEMORY-03` §9.4 adds a threshold-review step before freeze.
- No implementation, benchmark, model load, or runtime change occurred.
- No frozen artifact was edited.

### 9.8 Document map after revision 2

| Document | Answers |
|---|---|
| `MEMORY-01` (this file) | What the evidence says, and what the project has not established |
| `MEMORY-02` | Which sources support each claim, and their limitations |
| `MEMORY-03` | How the competing approaches would be tested |
| `MEMORY-04` | Whether automated lore has a defensible role, the eight research tracks, and the lore-specific ranking |
| `MEMORY-05` | Whether a flush can be turned into safe maintenance |
| `MEMORY-06` | Whether Howl constructs each model's prompt correctly |
| `MEMORY-07` | What the hardware can host and how it must be measured |
| `MEMORY-08` | Which approaches rank where, and what remains unknown |
