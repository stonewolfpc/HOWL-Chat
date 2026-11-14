# HOWL Chat — Automated Lorebook Research Package

**Status:** Research and architecture-evaluation only. **No architecture selected. No implementation authorization.** No source, schema, or runtime change.
**Date:** 2026-11-12
**Revision:** 2 — expanded from the 327-line automated-lore assessment into the full research package required by the automated-lorebook work order: generation (Track 1), maintenance (Track 2), retrieval and activation (Track 3), scope and knowledge boundaries (Track 4), local-model suitability (Track 5), existing systems (Track 6), orchestration (Track 7), deletion and cache integrity (Track 8), plus a lore-specific ranking, cooperation matrix, evaluation scenarios, and unresolved questions.
**Status change (revision 3):** Per the owner decision of 2026-11-12 (`docs/architecture/OWNER-DECISION-MEMORY-FORMATION-2026-11-12.md`), automated lore-entry **generation** is downgraded from an independent subsystem to one possible destination of a general memory-formation pipeline. Track 1 is retained as research, not as a subsystem proposal. Tracks 2–8 remain valid in substance and are re-homed to the memory store (maintenance, retrieval, scope, deletion, orchestration, local-model suitability, existing mechanics). The document's conclusions about atomicity, provenance, candidate-first review, checkpoint timing, and budget discipline are unchanged and now apply to every memory record, not only to lore.

The owner names the resulting capability **general conversational knowledge formation** and keeps the decision **provisional** until `Q-MF-01` and `Q-MF-10` close. Three owner amplifications in that record constrain this research directly: one authoritative home with many references (records are referenced, not copied), one pipeline with several processing passes (chunk → extract → reconcile → dedupe/contradict → validate), and model latitude to notice significance that no keyword or predefined type would capture (including an `unclassified` proposal path).

**Why this document exists:** Automated lore is the one part of Howl's memory surface that already exists in some form, has known noise problems, and has never been measured. It is also the part most likely to be over-built. This document is the evidence record for deciding what automated lore is genuinely good at, where it fails, and how it should cooperate with transcript retrieval, episodic memory, structured state, and hybrid retrieval. It is separate from `MEMORY-01` because lore has its own mechanism evidence (eight tracks), its own failure taxonomy, and its own value ledger, and folding it into the general synthesis would bury both.
**Companions:** `MEMORY-01` (rev 2), `MEMORY-02` (sources S1–S50), `MEMORY-03` (protocol rev 2, Stage E), `MEMORY-05` (flush), `MEMORY-06` (prompt formats), `MEMORY-07` (hardware), `MEMORY-08` (ranking).
**Controlling records:** `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md`, `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md`, `docs/review/F3C-LEGACY-LORE-CLASSIFICATION-2026-11-10.md`, `docs/review/F3C-MASTER-ARCHITECTURE-AND-IMPLEMENTATION-CONTRACT-2026-11-11.md` (frozen), `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md` (rev 9), `docs/architecture/OWNER-DECISION-USER-PERSONAS-2026-11-12.md`.

**Authority:** Research evidence. It alters no frozen contract, owner decision, or amendment, and authorizes no work.

---

## 0. Evidence labels

| Label | Meaning |
|---|---|
| **Observed** | Reproduced in a controlled environment by this investigation. **Not used** — no model was loaded, no index built, no benchmark run. |
| **Source-derived** | Stated by a primary paper, official documentation, or identified implementation. |
| **Calculated** | Arithmetic from stated values. |
| **Inferred** | Reasoning about Howl applicability, not directly demonstrated. |
| **Unestablished** | No adequate evidence found. |
| **Superseded** | Replaced by a later record. |

No number in this document is a Howl measurement. No vendor feature description is treated as proof of effectiveness.

---

## 1. Established owner constraints

These are constraints the design must satisfy, not findings. They are restated here so the research is read against them.

| ID | Constraint | Consequence for lore |
|---|---|---|
| `C-A` | Hybrid memory: multiple complementary methods, not one architecture | Lore is one contributor among transcript retrieval, episodic memory, structured state, lexical/dense/hybrid retrieval, summaries, temporal validity |
| `C-B` | Application orchestrates; models interpret | Identity, scope, storage, indexing, budgets, cache, deletion, provenance, and assembly are deterministic application work. Model calls are reserved for interpretation, extraction, contradiction, and ambiguity |
| `C-C` | Smallest useful selection; preserve creative voice | A lore system that makes dialogue read like a database report has failed, even if its facts are correct |
| `C-D` | Objects are independent; attachment is temporary | Lorebook and lore-entry identity ≠ attachment ≠ relevance ≠ prompt inclusion |
| `C-E` | Historical transcripts are evidence | "Remember when?" should be answerable from the original messages, not only from a derived entry |
| `C-F` | Segmented persistence, continuous chat | Segment boundaries are not memory boundaries for extraction or retrieval |
| `C-G` | Personas are first-class; Role is UI-only | Active persona is a scope dimension; other personas and the Role field never reach the AI |

---

## 2. The question, stated precisely

"Automated lorebook" names at least five mechanisms that fail differently:

| Mechanism | What it does | Where it runs |
|---|---|---|
| **Keyword-activated world info** | Injects entries whose trigger keys appear in recent text | Prompt assembly |
| **Semantic/vector activation** | Injects entries near the recent context in embedding space | Prompt assembly + index |
| **Recursive/chain activation** | An injected entry's keys can activate further entries | Prompt assembly |
| **Automatically generated lore** | A model writes entries from conversation | Background or checkpoint |
| **Discovered/chat-owned entries** | Entries created inside a conversation, scoped to it | Conversation store |

Howl's current surface is a mixture: legacy automated files that are prompt-eligible but unattributed, legacy files inert in old name-keyed folders, and newly generated conversation-owned discoveries (F3-C classification). The research question is therefore not "is lore good" but:

> **Which population, generated by which method, maintained by which rules, activated by which policy, at which token cost, produces a continuity benefit that transcript retrieval plus structured state would not produce on its own?**

No answer exists in the project record, and none exists in the published literature for creative-fiction lore specifically.

---

## 3. Track 1 — Generation

### 3.1 Generation methods

Each row is assessed against the required reporting dimensions. Cost is per 1,000 tokens of source conversation, at the model tier Howl would plausibly use.

| Method | What it does | Solves | Does not solve | Appropriate when | Cost / risk | Evidence |
|---|---|---|---|---|---|---|
| **Window summarization** | Summarizes a conversation window into prose | Cheap overview | Exact names, numbers, promises; provenance | Never as sole memory; useful as an index key | 1 call / window; summary drift on repeats | Source-derived (S2 summary detail loss; S1 facts-as-keys) |
| **Event extraction** | Emits discrete events with participants, time, place | Episodic recall; "when did X happen" | Current state; static world facts | Scene or session boundaries | 1 structured call per segment; schema missed fields | Source-derived (S5 event storage; S3 episodes) |
| **Entity extraction** | Identifies people, places, items, organizations | Indexing, alias discovery, dedupe keys | Meaning; relationships | Always, as an index layer | Low if deterministic NER; model cost if generative | Source-derived (S5 graph entities; S13 lexical strength) |
| **Relationship extraction** | Emits typed edges (kin, ally, rival, debt) | Relationship continuity; multi-hop questions | Change over time unless versioned | On relationship change or introduction | Higher error rate; needs stable entity IDs | Source-derived (S3/S5 graphs) |
| **Fact extraction (subject–predicate–object)** | Atomic claims with validity | Structured state; contradiction detection | Narrative nuance, tone, voice | Most durable-lore candidates | Model cost; over-extraction risk | Source-derived (S1 extraction-as-keys; S6 FAMA) |
| **Character-development tracking** | Tracks traits, goals, injuries, skills over time | Believability; voice consistency | Objective facts | Arc-level checkpoints | Needs supersession or it accumulates contradictions | Inferred; no direct study |
| **World-building extraction** | Rules, factions, geography, customs, technology | World consistency; reduced re-explanation | Personal/emotional continuity | When new world material appears | Over-generation; duplication | Source-derived (S40 lorebook practice) |
| **Scenario extraction** | Scene-specific setup, stakes, constraints | Scenario-local correctness | Cross-scenario knowledge | Scenario attach or change | Must be scoped or it leaks across scenarios | Source-derived (S46 chat/persona binding precedent) |
| **Preference / recurrence detection** | Repeated details (drink, nickname, habit) | Personalization; natural callbacks | One-off facts | After ≥2–3 occurrences | Cheap signal; false positives | Inferred |
| **Emotional/relational beat detection** | Marks turns with relational weight | Relationship continuity; callback salience | Factual recall | Significant beats only | Subjective; needs review | Inferred; `Unestablished` |
| **Narrative-arc tracking** | Open/closed threads, promises, foreshadowing | Plot continuity; payoff | Detail retrieval | Arc checkpoints | Model-heavy; drift | Inferred; `Unestablished` |

### 3.2 Automatic creation versus proposing candidates

| Mode | Advantage | Disadvantage | Evidence |
|---|---|---|---|
| **Auto-create (canon)** | Zero user effort; immediate availability | No provenance discipline; errors become canon; deletion cannot distinguish origin; floods the store | `Unestablished`; contradicted by observed deletion/provenance gaps (S10) |
| **Propose candidate (review)** | Provenance, correction, user control; supports `H5`/`H6` | Review burden; may never be reviewed | **Promising**; matches the memory design's evidence-backed proposal model |
| **Auto-create with trust tiers** | Low-risk facts immediate; risky facts reviewed | Needs a reliable risk classifier | `Promising`, unproven |
| **Silent store, no prompt** | Keeps extraction for later mining | Adds cost with no visible benefit | `Weak` |

**Judgement:** the evidence supports candidate-first for anything that will be treated as an assertion about the world or a character. Auto-commit is defensible only for objective, reversible, source-linked facts with a visible audit trail. This is an inference, not a measured result.

### 3.3 Batch versus incremental versus checkpoint

| Strategy | Latency | Consistency | Risk |
|---|---|---|---|
| **Batch (whole chat at once)** | High; deferred | Good global view; easy dedupe | Long context degradation (S1, S15); stale between runs; expensive |
| **Incremental (per turn)** | Low per run | Poor global view; duplicate-heavy | Constant model activity; cache thrash; cost per turn |
| **Checkpoint (scene/session/flush)** | Moderate | Good; bounded input; natural dedupe point | Requires a checkpoint definition (`MEMORY-05`) |
| **User-triggered** | None until asked | Depends on user discipline | Predictable cost; easy to forget |

**Judgement:** checkpoint-based is the only strategy that matches both the flush boundary (Track 7, `MEMORY-05`) and the cost constraint (`C-B`). Incremental per-turn extraction is the worst cost/benefit in the table. `Inferred`.

### 3.4 What belongs where

This is the placement decision that prevents lore from becoming a duplicate truth source.

| Information type | Home | Why |
|---|---|---|
| Stable world fact (geography, factions, rules) | **Lore** | Reusable across chats; rarely changes |
| Stable character fact (role, appearance, ties) | **Character card + structured state** | Character card is authored; state is derived |
| Current character state (injury, location, mood) | **Structured state** | Changes constantly; supersession required |
| Relationship state (trust, debt, promise) | **Structured state** | Versioned edges, not lore |
| Event ("the bridge collapsed") | **Episodic memory** | Time-anchored; retrieved as evidence |
| Original dialogue | **Transcript** | Authoritative evidence; never replaced by lore |
| Temporary fact (weather, a passing stranger) | **Nowhere** (or expiry-tagged) | Not worth storing; pollutes retrieval |
| Scenario-local setup | **Scenario object** | Scoped to the scenario |
| User preference | **Persona / preference store** | Person-scoped, not world-scoped |
| Model interpretation ("he seemed afraid") | **Nothing, or low-confidence note** | Speculation must not become canon |
| Correction / retcon | **State + provenance event** | Supersession, not a new lore entry |

**Anti-pattern:** promoting every interesting sentence. The placement table is the primary defence — if a piece of information has no home in the left column, the correct action is usually to store nothing.

---

## 4. Track 2 — Quality and maintenance

### 4.1 Deduplication

| Method | Catches | Misses | Cost |
|---|---|---|---|
| Exact match | Identical text | Paraphrase | Negligible |
| Normalized match (case, whitespace, unicode, punctuation) | Trivial variants | Rewording | Negligible |
| Fuzzy string (edit distance / token overlap) | Near-identical entries | Semantic paraphrase | Low |
| MinHash / SimHash | Near-duplicate documents | Short entries with different wording | Low |
| Embedding clustering + threshold | Semantic duplicates | Opposite claims (dangerously merged) | Embedding cost per write |
| Entity-key + predicate key | Duplicate facts about the same subject | Different predicates on the same subject | Low, needs entity resolution |

**Key hazard:** semantic dedupe can merge *contradictory* entries ("the detective retired" vs "the detective is active"). Dedupe must run only on same-subject, same-predicate candidates and must never merge across differing validity intervals.

### 4.2 Conflict taxonomy

| Conflict | Example | Correct handling |
|---|---|---|
| **Direct contradiction** | "She is alive" vs "she died" | Supersede with validity interval; keep both; state the change |
| **Refinement** | "A detective" → "a detective in the 14th precinct" | Merge as a more specific fact |
| **Temporal change** | "Lives in Metropolis" → "moved to Gotham" | Validity intervals; never overwrite |
| **Perspective difference** | Character A believes X; character B believes Y | Two asserted-by facts, not one contradiction |
| **Error** | Model misattributed an action | Retract with provenance; do not silently delete |
| **Scope difference** | World A fact vs World B fact | Separate scope keys; never "conflict" |

A system that cannot tell these apart will either freeze on false conflicts or silently overwrite real ones. `Source-derived` for the failure mode (S6 FAMA: invalid memories reused, evolution not reconciled); the taxonomy is `Inferred`.

### 4.3 Entry structure

| Structure | Advantage | Disadvantage | Verdict |
|---|---|---|---|
| **One giant document** | Simple | Unaddressable; cannot scope, cite, version, or delete a part | **Incompatible** with provenance and deletion |
| **Atomic entries** | Addressable, citable, scoped, deletable | More records; needs dedupe | **Strongly supported** as the default |
| **Grouped (entry + children)** | Balances addressability and coherence | Group deletion semantics needed | **Promising** |
| **Linked graph** | Multi-hop, relationship queries | Extraction cost, entity resolution | **Conditional** |
| **Hierarchical (entry → summary → world)** | Cheap navigation | Summary drift; two truth sources | **Conditional** |

`Source-derived` basis: atomic records are what makes deletion, provenance, and scoping possible at all (S10, S5, S3). The recommendation for atomic-by-default is `Inferred`.

### 4.4 Maintenance rules that follow

1. Every entry carries a source span and an entry ID.
2. Every assertion carries a validity interval and an assertion status (`current`, `superseded`, `disputed`, `deleted`).
3. User edits outrank automated edits; an automated edit never silently overwrites a user edit.
4. Version history is retained; supersession is a new version, not a mutation.
5. Confidence is stored, not implied.
6. Expiry is a field, not a cleanup job.
7. A merge produces a new entry plus a merge record naming its parents.

---

## 5. Track 3 — Retrieval and activation

### 5.1 Activation methods

| Method | Catches | Misses | Failure mode | Evidence |
|---|---|---|---|---|
| **Exact key** | Named entity, term, quoted phrase | Paraphrase | Common-word triggers | Source-derived (S46, S47, S40) |
| **Alias list** | Nicknames, titles, epithets | Unlisted variants | Alias explosion | Source-derived (S46) |
| **Regex** | Structured patterns, variants | Semantics | Catastrophic false positives | Source-derived (S46) |
| **Phrase match** | Multi-word terms | Word-order variance | Low recall | Source-derived (S46) |
| **Semantic/dense** | Meaning-level match | Exact strings | Thematic-but-irrelevant entries; silent threshold sensitivity | Source-derived (S46 vector storage; S13) |
| **Hybrid + rerank** | Both coverage types | Still threshold-dependent | Latency | Source-derived (S13) |
| **Contextual/state-conditioned** | Correct entry for the scene | Requires trustworthy state | Wrong state → wrong entry | `Inferred` |
| **Relationship expansion** | Related entries via links | Over-expansion | Budget exhaustion | `Inferred` |
| **Recursive** | Chained entries | Depth control | **Recursion storms** | Source-derived (S46) |
| **Always-on / constant** | Core world rules | Budget cost | Crowds out everything else | Source-derived (S46, S48) |

### 5.2 Documented product mechanics worth copying or avoiding

| Mechanic | Observed behaviour | Assessment |
|---|---|---|
| Per-entry token budget + global budget | Entries dropped when the budget fills | Necessary; the drop must be **logged**, not silent |
| Activation probability | Probabilistic insertion after activation | Adds non-determinism; useful for variety, harmful for evaluation |
| Inclusion groups | Only one entry from a group may activate | Useful to prevent variant pile-up |
| Sticky / cooldown / delay | Keeps an entry active for N turns | Useful for scene continuity; needs a determinate rule |
| Scan depth (how far back keys are matched) | Bounded lookback | Critical for cost; too deep = false activation |
| `max entries` for vector matches | Caps semantic activations | Necessary for budget control |
| Keyword rules vs embedding rules | Rule-based is "more rigid but predictable"; embedding is flexible but less controllable | **The trade-off is explicit in vendor documentation** (S48) |
| Entry text must be self-contained | The model never sees the key or title | **Important**: entries must read as complete statements (`S47`) |
| Dynamic prompt sources cause cache misses | Documented for world info, summarization, and chat vectorization | Directly supports the flush/cache analysis in `MEMORY-05` (S46) |

### 5.3 Budget control

Activation without a budget is the single most common failure in every documented implementation. The minimum control set:

1. A per-turn lore token ceiling, declared as a fraction of the context window.
2. A per-entry ceiling.
3. A deterministic priority order (explicit pin > state-conditioned > keyword exact > alias > semantic).
4. A drop log: which entries were considered, which were admitted, which were dropped, and why.
5. No unbounded recursion; depth and total expansion are capped.
6. A visible count to the user (entries admitted / tokens used).

---

## 6. Track 4 — Scope, identity, and knowledge boundaries

### 6.1 The four-layer ladder

The work order's central distinction, and the one most systems collapse:

| Layer | Question | Controlled by |
|---|---|---|
| **1. Attached** | Is this object connected to this context? | Attachment record |
| **2. Available** | May it be retrieved here? | Scope key on the object + attachment |
| **3. Relevant** | Should it be injected this turn? | Retrieval + budget |
| **4. Known** | Does this character actually know it? | Knowledge record (witnessed / told / inferred / read) |

**Attachment is not knowledge.** A lorebook attached to a chat containing "the butler did it" does not mean the detective knows it. Availability is not knowledge either. Only layer 4 answers "what does this character know", and it must be a separate record from the entry.

### 6.2 Knowledge acquisition modes

| Mode | Example | Who knows |
|---|---|---|
| **Witnessed** | Character was present at the event | That character |
| **Told** | Another character told them | The listener (and the teller as "told") |
| **Inferred** | Deduced from evidence | The deducer, with lower confidence |
| **Read** | Read a document/note in-world | The reader |
| **Author-granted** | Card or world states they know it | That character only |
| **Player-known** | The user knows it, no character does | No character |

A character's answer may only use facts whose knowledge record includes that character. This is the same zero-tolerance boundary as `H3`; persona adds a user-side dimension (`C-G`).

### 6.3 Flexible cross-category attachment without leakage

The owner requires that a world lorebook may attach to a character, a character may use several lorebooks, and one resource may attach to multiple scenarios. That is compatible with strict isolation **only if** the scope key is on the *attachment and the entry*, not inferred from the object's category:

- A lorebook attached to three scenarios is still one object; each attachment is its own scope.
- Retrieval for scenario A may not return scenario B's attachment rows.
- "Shared" does not mean "pooled": sharing an object must never mean sharing a candidate set.

`Source-derived` precedent: SillyTavern binds lorebooks to a character, a persona, **or an individual chat**, and documents that chat-bound entries do not carry to other chats (S46). That is the attachment model Howl needs, minus the multi-context reuse requirement.

---

## 7. Track 5 — Local model suitability

### 7.1 What the model should and should not do

| Task | Deterministic alternative? | Verdict |
|---|---|---|
| Chunking, IDs, timestamps, scope keys | Yes | **Application** |
| Turning text into structured records | Partly (constrained decoding guarantees shape) | **Application constrains, model fills** |
| Deciding significance | No | **Model**, with review |
| Detecting contradiction between two entries | Partly (same subject+predicate+differing object) | **Application detects, model adjudicates** |
| Entity resolution | Partly (exact/alias) | **Application for exact; model for ambiguous** |
| Writing an entry's prose | No | **Model** |
| Deciding what a character knows | No | **Application from event participation; model only for inference** |
| Budgeting and ordering | No | **Application** |
| Deletion and invalidation | No | **Application** |

### 7.2 What constrained decoding actually guarantees

Verified from primary documentation:

- llama.cpp converts a subset of JSON Schema to GBNF; tokens that would violate the grammar get logits set to −∞, so output **parses** (S49).
- The schema is **not injected into the prompt**; the model has no visibility into it. If the structure must be understood, it must be described in the prompt (S49).
- Grammar constraints can still produce **incomplete** output if generation is truncated by a token limit (S49 supporting discussion).
- Token-boundary artifacts can force a model onto awkward continuations (S49).
- A real, reported defect exists where JSON-Schema-constrained decoding fails on certain chat-template control tokens while equivalent GBNF succeeds (S50).

**Consequence for Howl:** structured output is an application-enforced property, not a model capability. That removes "does the small model emit valid JSON" as a legitimate reason to prefer a larger model — but it does **not** remove the semantic problem (wrong entity, wrong predicate, invented fact).

### 7.3 Evidence on small-model extraction

| Finding | Evidence | Implication |
|---|---|---|
| 8B models can perform extraction usefully, but extraction quality caps gains | S1 (LongMemEval uses Llama-3.1-8B extraction) | Do not assume frontier extraction |
| Smaller models often miss temporal conflicts; larger models detect but fail to act | S18 | Put update semantics in the store |
| Memory agents reuse invalid memories and fail to reconcile changes | S6 | Supersession must be structural |
| Long-context extraction degrades with input length | S1, S15 | Segment input; do not extract from a 100k-token window |
| Generic benchmark skill does not predict creative extraction | `Unestablished` | No creative-memory extraction benchmark exists (`Q-M3`) |
| Quantization effect on extraction specifically | `Unestablished` | Must be measured, not assumed |

### 7.4 Multi-stage pipelines

The evidence supports splitting the work rather than asking one call to do everything:

```text
stage 1  deterministic segmentation + IDs + scope keys        (no model)
stage 2  candidate spans (entity/event/fact)                  (model, constrained)
stage 3  schema validation + dedupe + conflict detection      (no model)
stage 4  significance scoring                                 (application rules + model)
stage 5  entry prose                                          (model, only for admitted candidates)
stage 6  provenance attach + budget admission                 (no model)
```

Stages 2 and 5 are the only places where a model is genuinely required, and stage 5 runs far less often than stage 2. `Inferred` as a design, but every element is `Source-derived`.

---

## 8. Track 6 — Existing systems

| System | What it implements | Solves | Weakness | Transfers to Howl? |
|---|---|---|---|---|
| **SillyTavern** | World Info: keys, secondary logic, regex, vectorized entries, recursion, inclusion groups, probability, position/depth, sticky/cooldown, budgets; bind to character, persona, or chat | Flexible, author-controlled context injection | No provenance, no deletion model, budget drops are silent, dynamic sources cause cache misses (S46) | **Mechanics yes; governance no** |
| **NovelAI Lorebook** | Keys, regex, order/position, search range, force activation, probability, budget | Simple, predictable | No quality measurement | **Field set yes** |
| **KoboldAI / KoboldCpp** | Memory + Author's Note + World Info; Context Shift; separate TextDB | Keyword-only insertion conserves tokens; entry text self-contained (S47) | No maintenance layer; context shift is not memory | **Entry prose guidance yes** |
| **Layla AI** | Imports ST lorebooks (keyword rules, same behaviour); its own lorebooks use embedding similarity; local model (S48) | Two activation modes with an explicit predictability/flexibility trade-off | No published effectiveness data | **The comparison yes** |
| **LM Studio** | No lorebook; overflow policies | — | — | **Overflow comparison only** |
| **Ollama** | Template-driven; no cross-request memory | — | — | **No** |
| **MemGPT / Letta** | Tool-driven memory editing | Model-managed memory | Tool-call discipline failures (S11) | **No — model should not manage the store** |
| **Mem0** | LLM extraction + ADD/UPDATE/DELETE/NOOP over vectors | Consolidation | Loses to full context in its own table (S4) | **The operation taxonomy yes** |
| **Zep / Graphiti** | Bi-temporal graph + hybrid retrieval | Temporal validity | Extraction/graph cost (S3) | **Bi-temporal model yes** |
| **LangMem** | Store-first primitives | Low infrastructure | No benchmark (S25) | **Primitive separation yes** |
| **RAPTOR-class** | Hierarchical summaries | Global questions | Detail loss, reader confound (S14) | **Only as an index** |
| **GraphRAG** | Entity graph + community summaries | Thematic queries | LLM-heavy indexing (S27) | **No as default** |

**What transfers:** the field set, the activation mechanics, the budget discipline, bi-temporal validity, and the ADD/UPDATE/DELETE/NOOP operation taxonomy. **What does not:** model-managed stores, unbounded recursion, unbudgeted injection, and any architecture that cannot name a source span for an entry.

---

## 9. Track 7 — Operational timing and orchestration

| Trigger | Operation | Sync/async | Cost | Risk |
|---|---|---|---|---|
| **Every turn** | Nothing beyond scope keys and budget accounting | sync, negligible | ~0 | None |
| **Turn end** | Append event to the log (deterministic); no model | sync, cheap | Low | None |
| **Scene boundary** | Event extraction, state update, candidate lore | async preferred | Moderate | Mid-scene flush loses context (`FL-E`) |
| **Flush / rollover** | Capture → extract → consolidate → reconstruct (`MEMORY-05`) | async, with a visible completion state | Moderate–high | Interruption (`FL-D`) |
| **Chat save / close** | Consolidation and index flush | async | Moderate | User may close before completion |
| **User request** | Rebuild, re-index, "remember this", "forget this" | sync for the command, async for the work | Controlled | None if receipted |
| **World/scenario switch** | Re-scope, re-assemble, do not re-extract | sync | Low | Cross-scope contamination if scope keys are wrong |
| **Attach/detach** | Update attachments; invalidate cached package | sync | Low | Stale cache if invalidation is missed |
| **"Remember when?"** | Transcript retrieval + evidence package | sync (on the turn) | Retrieval only | Must respect knowledge boundaries |
| **Contradiction detected** | Propose correction; do not silently rewrite | async + review | Low | Silent overwrite is worse than a visible conflict |
| **Continuous background** | — | — | High, unpredictable | **Not recommended** |

**Default posture (inferred):** no continuous background extraction. Extract at scene/flush/save checkpoints and on explicit request. The model work should be bounded, deferred, and cancelable, and the conversation must never wait on it.

**Why:** continuous extraction violates the cost constraint (`C-B`), thrashes the prompt cache (S46 documents cache misses from dynamic prompt sources), and produces a stream of low-value entries that no one reviews — the exact mechanism by which automated lorebooks degrade.

---

## 10. Track 8 — Deletion, cache, and integrity

### 10.1 States (six, not two)

| State | Meaning | Visible to retrieval? |
|---|---|---|
| **Not loaded** | Exists in storage, not in memory | No |
| **Loaded and current** | Resident; valid | Yes |
| **Modified** | Changed but not committed | No (or the old version only) |
| **Superseded** | Replaced by a newer version; history retained | Only for as-of queries |
| **Invalidated** | Its source was deleted or changed; unusable | No |
| **Explicitly deleted** | Removed | No, and must not reappear |

**A missing cache entry is not proof of deletion.** Absence and deletion are different states and must be reported differently.

### 10.2 Propagation matrix

| Deleted object | Must propagate to |
|---|---|
| Source message | Events, facts, entries, summaries, embeddings, cached packages, receipts |
| Lore entry | Index rows, embeddings, cached package, parent/child links |
| Character | Character-scoped entries, knowledge records, relationship edges |
| World | World-scoped entries, attached-resource links |
| Scenario | Scenario-scoped entries, attachments |
| Persona | Persona-scoped facts, attachments, cached package (never other personas) |
| Chat | Everything scoped to it, including segments |

### 10.3 Resurrection vectors

| Vector | Mechanism | Control |
|---|---|---|
| Stale index | Index not rebuilt after delete | Tombstone check at query time |
| Embedding neighbour | Deleted text reachable via a retained vector | Delete vectors with the record |
| Cached prompt | Package built before deletion | Invalidate on any scope change |
| Summary | Summary contains the deleted fact | Rebuild or invalidate derived summary |
| Backup / restore | Old file restored | Restore must replay the deletion log |
| Second population | Legacy or imported copy | One authoritative store; imports are candidates |

`Source-derived` basis: no evaluated system simultaneously achieves utility, access control, and reliable forgetting (S10); memory agents reuse invalid memories (S6).

---

## 11. Lore-specific ranking

Ranks are for **the lore subsystem only**; `MEMORY-08` ranks the whole approach space. Rank is by fit to `C-A`–`C-G` and evidence quality, not by popularity.

| Rank | Approach | Class | Rationale |
|---:|---|---|---|
| 1 | **Atomic, source-linked entries** with validity intervals | Strongly supported | Only form that supports provenance, scoping, supersession, and per-entry deletion |
| 2 | **Exact key + alias activation** with a hard budget and a drop log | Established | Predictable; matches exact-name strength; cheap |
| 3 | **Candidate-first generation** at scene/flush checkpoints | Promising | Provenance and review; matches `C-B` cost discipline |
| 4 | **Detail assistance** as the entry's purpose | Strongly supported | Targets the documented summary weakness (S2) |
| 5 | **Hybrid retrieval** (keyword ∪ semantic) + rerank for admission | Strongly supported | Coverage without abandoning exactness |
| 6 | **State-conditioned activation** via structured state | Conditional | Most principled; depends on `C7`-class state |
| 7 | **Relationship expansion** with caps | Conditional | Useful for multi-hop; over-expansion risk |
| 8 | **Semantic-only activation** | Conditional | Flexible but threshold-fragile and less controllable (S48) |
| 9 | **Recursive activation** with hard depth+token caps | Weakly supported | High storm risk; little unique benefit over links |
| 10 | **Hierarchical lore summaries** | Conditional | Navigation only; two truth sources |
| 11 | **One monolithic lorebook document** | Incompatible | Cannot scope, cite, version, or delete a part |
| 12 | **Auto-canon generation, no review** | Incompatible | Violates `H5`/`H6`; floods the store |
| 13 | **Continuous background extraction** | Weakly supported | Cost, cache thrash, no review path |
| 14 | **Lore as the primary memory** | Incompatible | Duplicates state and transcript; loses evidence |

**Nothing is eliminated except where it conflicts with a stated owner constraint** (`C-D`, `C-E`, `H5`, `H6`) or cannot support deletion.

---

## 12. Cooperation matrix

| Partner | Lore's role | Boundary rule |
|---|---|---|
| **Transcript retrieval** | Index and shortcut to the right region | Transcript is authoritative; a lore entry never overrides it |
| **Episode memory** | Names, context, and stable framing for an event | Events live in episodic memory; lore does not store them |
| **Structured state** | Stable background; not current state | State wins on any conflict about the present |
| **Summaries** | Entry text may seed a summary; summary is not an entry | One consolidation pipeline, declared precedence |
| **Hybrid retrieval** | Lore entries are an additional candidate source | Same budget, same reranker, same dedupe |
| **User persona** | Persona-scoped details never enter lore | Persona scope is its own key (`C-G`) |
| **Character knowledge** | Entries are knowledge candidates, not knowledge | Layer 4 decides who knows |
| **Prompt assembly** | Inserted at declared template-safe points | Never mid-turn (`MEMORY-06`) |

---

## 13. Evaluation scenarios for later validation

These are the concrete fixtures for Stage E (`MEMORY-03`). They exist so a future validation is pre-specified, not improvised.

| ID | Scenario | Fixture | Primary measure |
|---|---|---|---|
| `E-01` | Lore on/off ablation | One world, 12 sessions, 20 lore-eligible queries | Correctness, detail retention, token overhead |
| `E-02` | False activation | Common-word keys, 20 near-miss turns | False-activation rate |
| `E-03` | Budget exhaustion | 60 entries, 4 budget levels | Admitted vs dropped; answer impact of drops |
| `E-04` | Recursion storm | Chained entries, depth 0/1/2/3 | Tokens consumed; relevance of expansion |
| `E-05` | Duplicate flood | Three paraphrases of one fact | Dedupe accuracy; duplicate rate at prompt |
| `E-06` | Contradiction | Character retires mid-arc | Supersession honored; stale reuse rate |
| `E-07` | Scope leak | Same entry attached to two scenarios | Cross-scenario leakage (must be zero) |
| `E-08` | Knowledge boundary | Entry known to one character only | Spoiler rate (must be zero) |
| `E-09` | Deletion | Delete a source message owning 3 entries | Resurrection count (must be zero) |
| `E-10` | Segment boundary | Retrieval across two storage segments | Recall vs single-segment control |
| `E-11` | Checkpoint timing | End-of-scene vs per-turn extraction | Cost, duplicate rate, answer quality |
| `E-12` | Competitive prose | Blinded rater: lore-on vs transcript-only replies | "Database report" failure rate |

`E-12` is the constraint-`C-C` test. A lore system that wins `E-01` and loses `E-12` has failed the product requirement.

---

## 14. Unresolved questions

| ID | Question | Priority |
|---|---|---|
| `Q-L1` | What is the current false-activation rate of Howl's lore? | P1 |
| `Q-L2` | Does lore improve correctness, or only style? | P1 |
| `Q-L3` | What token budget is appropriate for the smallest sampled reader? | P1 |
| `Q-L4` | Keyword, semantic, hybrid, or state-conditioned activation? | P1 |
| `Q-L5` | Candidate or canon for generated lore? | P1 |
| `Q-L6` | Can lore be compressed without losing detail? | P1 |
| `Q-L7` | Inside the memory pipeline or beside it? | P2 |
| `Q-L8` | How does a user see why an entry activated? | P2 |
| `Q-L9` | What happens to legacy populations on migration? | P2 |
| `Q-L10` | Is deterministic ordering achievable across restarts? | P1 |
| `Q-L11` | What is the right consolidation unit: turn, scene, or session? | P1 |
| `Q-L12` | Does a lore entry ever outrank a transcript passage? | P2 |
| `Q-L13` | How should `Role`-style UI-only metadata be prevented from entering entry text? | P2 |
| `Q-L14` | Can entry prose be generated with voice constraints (no database-report tone)? | P1 |
| `Q-L15` | What is the cost of entry prose generation at the chosen checkpoint cadence? | P1 |
| `Q-L16` | Does a second population (legacy/imported) need its own trust tier? | P2 |

---

## 15. What this document does not do

- It does not select a lorebook architecture, storage model, or activation policy.
- It does not authorize implementation, schema changes, or prompt changes.
- It does not measure anything: no model, no index, no benchmark.
- It does not modify the frozen F3-C records or the legacy-lore classification.
- It does not decide the disposition of any existing lore population.
- It does not treat any vendor feature list as evidence of effectiveness.

**The evidence supports one narrow conclusion and refuses the rest: automated lore is defensible only as atomic, source-linked, budgeted, candidate-first detail assistance that cooperates with transcript retrieval and structured state. Whether that is worth building, and in what form, is an owner decision after validation.**
