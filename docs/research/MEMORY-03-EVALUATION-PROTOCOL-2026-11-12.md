# HOWL Chat — Memory Systems Evaluation Protocol

**Status:** Evaluation protocol design only. **Not execution authorization.** No model was loaded, no benchmark was run, no index was built, no runtime experiment was performed, and no architecture was selected.
**Date:** 2026-11-12
**Revision:** 2 — supersedes revision 1 of this path (535 lines / 24,779 bytes). The revision is required by the owner work order of 2026-11-11: owner hypotheses `H1`–`H8` restored, the model envelope opened beyond 3–14B, hardware identifiers renamed to `HW-n`, `C0` split into its five distinct concepts, the `C1` question resolved, context flushing treated as an intentional design stage, automated lore expanded, metrics broadened to prompt compatibility and resource cost, and a threshold-review step added.
**Research basis:** `docs/research/MEMORY-01-RESEARCH-SYNTHESIS-2026-11-12.md` (rev 2), `docs/research/MEMORY-04-AUTOMATED-LORE-RESEARCH-ASSESSMENT-2026-11-12.md`, `docs/research/MEMORY-05-CONTEXT-FLUSH-AND-RECONSTRUCTION-RESEARCH-2026-11-12.md`, `docs/research/MEMORY-06-PROMPT-FORMAT-COMPATIBILITY-MATRIX-2026-11-12.md`, `docs/research/MEMORY-07-HARDWARE-AND-INFERENCE-OPTIMIZATION-2026-11-12.md`, `docs/research/MEMORY-08-COMPARATIVE-RANKING-AND-UNRESOLVED-QUESTIONS-2026-11-12.md`
**Source basis:** `docs/research/MEMORY-02-SOURCE-REGISTER-2026-11-12.md` (S1–S50)
**Pre-registration:** `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` fixes the corpus, arms, variable control, rubrics, gates, application-vs-model execution allocation, and the freeze-evidence checklist. Where this protocol and that package differ on a pre-registration detail, the package controls; where they differ on stage mechanics, this protocol controls.
**Controlling design/requirements:** `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md`, `docs/architecture/DESIGN-LOCAL-RAG-CREATIVE-2026-11-11.md`, `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`, and `docs/architecture/AMENDMENT-02-CONVERSATION-CONTINUITY-AND-MULTI-CHARACTER-ROUTING-2026-11-11.md` (rev 9).

---

## 0. Revision record

| Rev | Date | Change | Reason |
|---:|---|---|---|
| 1 | 2026-11-12 | Initial protocol: `C0`–`C7`, `S0`/`S1`, `L0`–`L2`, stages A–E, eight propositions labelled `H1`–`H8` | First pre-registered design |
| 2 | 2026-11-12 | Owner `H1`–`H8` restored as the hypothesis set; former `H1`–`H8` propositions renamed `P1`–`P10`; model envelope opened; `T-S/T-M/T-L` retained only as sampling tiers; `H1/H2/H3` hardware IDs renamed `HW-1/HW-2/HW-3`; `C0` split into `C0-H`/`C0-W`/`C0-R`/`C0-WS`; `C1` disposition recorded; Stage F (context flush) added; Stage E expanded to `L0`–`L3`; prompt-format axis and `F0`–`F2` added; metrics expanded to 20 families; threshold-review step and aggregate-score prohibition added | Owner work order 2026-11-11 §3–§7 |

**Nothing is implemented by revision 2.** It is still a protocol.

### 0.1 What revision 2 does not change

- The hard safety gates remain pass/fail and are never averaged (§9.1).
- Layer 2 remains blocked (§4.2).
- No condition is an architecture candidate by being tested (§7.8).
- No threshold may be changed after outcomes are observed (§9.4).
- No real user data is used without a separate written authorization (§6.4).

---

## 1. Purpose and non-authorization

This protocol defines how Howl would compare competing memory, retrieval, context-management, and prompt-format approaches under Howl-like conditions. It does not compare them.

Designing the protocol is not permission to:

- load a model, build an index, or run a benchmark;
- inspect user chats, character cards, worlds, or lorebooks;
- modify source, schema, runtime configuration, or databases;
- implement or select a memory architecture;
- close `DA-M1`–`DA-M6`, `DA-R1`–`DA-R6`, `R-10`, `R-11`, or `REC-CH-09`–`REC-CH-11`.

Every execution step requires the separate authorization gates in §13.

---

## 2. Owner hypotheses `H1`–`H8`

These are the owner's approved research positions. They are recorded as hypotheses to be tested, not as implementation authorization and not as pre-decided architecture.

| ID | Owner position (preserved) | What the protocol must do about it |
|---|---|---|
| `H1` | Long context alone is insufficient. Alternatives must be investigated: context management, rollover, context shifting, KV-cache strategies, memory-assisted reconstruction, and user-selectable overrides. | `C0-H`/`C0-W`/`C0-R`/`C0-WS` exist as separate conditions; KV/attention/rollover are a declared configuration axis (§5.5); overrides are a product requirement, not a test condition |
| `H2` | Summary-only memory is insufficient as a presumed complete solution. Combinations with raw evidence, hybrid retrieval, and other approaches must be researched. | `C6S` vs `C6R` remain separate; summary-as-index is never the only summary condition; claims about summaries must be conditioned on raw-evidence re-injection |
| `H3` | Character/conversation isolation is mandatory. | Scope is a hard gate (`D-SAFE-1`), not a scored metric; `S2` requires structural enforcement at the storage/query boundary (§7.5) |
| `H4` | Explicit temporal validity, supersession, and correction are strongly desired. | `C7` plus as-of queries, `D-ANS-1`, and the `T-*` temporal metric family; supersession is tested as store semantics, not model judgment |
| `H5` | Provenance and source traceability must be preserved. | Provenance-completeness metric family (§8.7); any condition that cannot cite a source span is measured as incomplete even if the answer is correct |
| `H6` | Deletion means gone, including applicable derived data and caches. | Stage D, `D-SAFE-2/3`, and the derived-surface enumeration in §7.6 |
| `H7` | Automated lore must demonstrate actual value and fit into the broader memory design. | Stage E with `L0`–`L3` and a value ledger (§7.7, `MEMORY-04`) |
| `H8` | Automated lore must be investigated for improvement, compression, repurposing, or replacement. | `L1`–`L3` are redesign candidates, not just an on/off switch; retirement is one admissible outcome, not a default |

**A hypothesis is not confirmed by this protocol existing.** Each has an explicit success/failure/inconclusive rule (§9).

### 2.1 Research propositions `P1`–`P10`

These replaced the mislabelled `H1`–`H8` of revision 1. They are external-evidence claims the protocol can test on Howl-like tasks.

| ID | Proposition | Supports owner hypothesis | Primary evidence |
|---|---|---|---|
| `P1` | Long context alone degrades on long histories; retrieval improves answer accuracy. | `H1` | S1, S15, S19 |
| `P2` | Summary-only memory loses answer-bearing detail; summary-as-index plus raw evidence recovers it. | `H2` | S1, S2 |
| `P3` | Hybrid lexical+dense retrieval with reranking improves evidence recall over either alone. | `H1`, `H2` | S13, S1 |
| `P4` | Explicit validity/supersession improves temporal and correction accuracy over full context or model-judged updates. | `H4` | S3, S5, S6, S18 |
| `P5` | Scope-before-retrieval prevents character/conversation leakage better than pooled retrieval plus post-filtering. | `H3` | S7, S8, S9, S10 |
| `P6` | Deletion requires invalidation of derived representations; deleting a raw record is insufficient. | `H6` | S10, S6 |
| `P7` | Weak/local readers need smaller, more precise context packages; a condition is portable only if it survives the weakest sampled tier. | `H1`, `H2` | S1, S15 |
| `P8` | Automated lore adds noise and token cost without a demonstrated recall benefit unless source-linked, bounded, and reviewable. | `H7`, `H8` | S1, S2, `MEMORY-04` |
| `P9` | Context flushing can be used as a consolidation/reconstruction boundary without net information loss if pre-flush capture and post-flush reconstruction are both measured. | `H1` | `MEMORY-05`; no direct benchmark exists |
| `P10` | Prompt-format fidelity (exact template, role markers, end-of-turn tokens, reasoning channel) changes answer quality enough that memory comparisons are invalid across mismatched formats. | `H1` | `MEMORY-06`; S28–S31, S33–S37 |

`P9` and `P10` are **unestablished**: no external benchmark isolates them. They are included because a protocol that ignores them would produce confounded results.

---

## 3. Two evaluation layers

### 3.1 Layer 1 — isolated algorithm evaluation

Uses synthetic fixtures and an isolated test store. It measures selection, retrieval, reasoning, isolation, temporal behavior, deletion mechanics, prompt-format fidelity, and cost without touching real chats or the product database.

**What Layer 1 can establish:** whether an approach works on controlled Howl-like tasks with pinned models, prompts, and hardware.

**What Layer 1 cannot establish:** Wails/Go integration, the live prompt path, real user-data behavior, UI correction burden, or production latency.

### 3.2 Layer 2 — product-path validation

Uses the actual Howl prompt assembly, storage, UI, and deletion paths. It is **not authorized** here and remains blocked by:

- unresolved canonical root/chat identity (`ROOT-01`, `D-ROOT`);
- unresolved storage adapter (`DA-M1`, `DA-R1`, `D-STORAGE`);
- unproven live memory/RAG reachability in the reply path;
- unapproved deletion/retention policy (`DA-M5`);
- unapproved conflict/automation policy (`DA-M6`);
- the live prompt path's approximate truncation (`contextSize × 4` characters) and hardcoded legacy ChatML builder, which cannot yet be shown to match any pinned model's template (`MEMORY-06` §7).

Layer 2 must not be simulated and reported as product evidence. A Layer 1 result is an algorithm result, not a Howl product result.

---

## 4. Model envelope

### 4.1 The supported envelope is broad, not 3–14B

**Revision 2 removes the implication that 3–14B models define Howl's supported models.** They never did. Howl's objective is broad local-model support, and large models where the user's hardware permits them.

| Statement | Status |
|---|---|
| Howl supports local GGUF models across a wide parameter range | Product objective |
| Large models are supported where hardware permits | Product objective |
| Users choose the model and may override resource policy | Requirement; not a test variable |
| No model is excluded because its parameter count is large or small | Requirement |
| The tiers below are sampling categories for evaluation only | Protocol device |

### 4.2 Sampling tiers (not a support boundary)

Evaluation cannot run on everything, so three tiers are sampled. **A tier is a sampling choice, not a supported-model boundary, not a default, and not a recommendation.** Results are reported per tier and never generalized to "Howl supports/does not support model size X."

| Tier | Sampling role | Example families (illustrative, not endorsed) | Why sampled |
|---|---|---|---|
| `T-A` | Smallest reader likely to be usable with tight retrieval | ~3–4B instruct (e.g. Qwen3-4B, Gemma-class 4B) | Tests the `P7` portability floor |
| `T-B` | Mid local tier | ~7–9B instruct (e.g. Qwen3-8B, Llama-3.1-8B, DeepSeek-R1-Distill-8B) | Common local usage |
| `T-C` | Larger local tier where hardware permits | ~12–32B instruct | Tests whether quality gains justify cost |
| `T-D` (optional) | Largest model the owner's hardware can host at the chosen quant | Any | Tests the upper bound of the envelope |

**Rules:**

1. Every sampled model file is pinned by SHA-256, quantization, tokenizer hash, chat-template hash, and template source before execution.
2. A substitute model is allowed only within the same tier and only by a dated protocol amendment **before** data collection.
3. A model that cannot load within the declared reserve policy (§4.4) is recorded as "not hostable on `HW-n` at this quant", not as "unsupported by Howl".
4. Reasoning models are sampled as a separate sub-axis (`R0` non-reasoning, `R1` reasoning-enabled with the declared budget), because reasoning tokens change both latency and effective context.
5. A result from one tier never stands in for another. Cross-tier claims require both tiers to have been run.

### 4.3 Quantization axis

| Setting | Role |
|---|---|
| `Q4_K_M` | Primary comparison quant for all sampled tiers |
| `Q5_K_M` / `Q6_K` | Sensitivity, one tier only, if the hardware permits |
| `Q8_0` | Upper-reference sensitivity, one tier only, if the hardware permits |
| Quantized V-cache | Declared separately; **requires Flash Attention** (S32). A run with quantized V-cache and Flash Attention off is invalid, not "slower". |

A quant change is a model change for reporting purposes. Cross-quant comparisons are reported as a separate sensitivity table, never merged into the main matrix.

### 4.4 Hardware configurations `HW-1`, `HW-2`, `HW-3`

Revision 1 used `H1`–`H3`, which collided with the owner hypotheses. Renamed.

| ID | Configuration | Purpose | Default reserve policy |
|---|---|---|---|
| `HW-1` | The owner's Windows PC, GPU offload enabled — **the initial test environment, not the permanent definition of supported hardware** | Primary target | **10% free VRAM after model load**, as an intended default. Recorded, measured, and adjustable by the user |
| `HW-2` | Same machine, reduced or zero GPU offload | Degraded/portable operation; CPU-only comparison | Reserve policy re-derived from available system RAM; recorded, not assumed |
| `HW-3` | Optional second machine with a different GPU vendor or VRAM size | Replication only | Reserve policy recorded per machine |

**Reserve-policy rules:**

1. The 10% figure is a **Howl policy default for `HW-1`**, not an upstream inference-engine constant. Upstream llama.cpp's `--fit` path uses a fixed `--fit-target` default of 1024 MiB (S32); the two numbers must not be conflated in any report.
2. On other hardware, the reserve is computed by a declared policy (for example a percentage with a floor and a ceiling) and recorded in the manifest. It is never copied blindly from `HW-1`.
3. **No universal OOM prevention is claimed.** Every configuration records pre-flight estimates, the actual post-load free VRAM, and any allocation failure. A claim that "Howl prevents OOM" is out of scope for this protocol.
4. A run in which the reserve was breached is flagged, and its latency/memory numbers are reported as "reserve-breached" rather than silently averaged.
5. The user may override the reserve. An override is a recorded configuration, not a protocol violation.

**Required manifest fields per hardware:** GPU model, VRAM, driver version, backend build, CPU, system RAM, OS build, page-file/swap state, offload layer count, `n_batch`/`n_ubatch`, KV-cache type, Flash Attention state, context size, reserve policy and actual free VRAM.

### 4.5 Runtime configuration axis (declared, not free)

These are not "conditions" competing for quality; they are **declared settings** that must be identical within a comparison and reported with every result. Changing one mid-comparison invalidates the comparison.

| Setting | Recorded values | Why it matters |
|---|---|---|
| Context size | e.g. 4,096 / 8,192 / 20,480 | Determines what `C0-W` can even contain |
| Rollover mode | off / shift / memory-guided | Howl's `context_rollover_mode` is a product setting, currently not implemented in the chat path |
| KV-cache type | f16 / q8_0 / q4_0 | Quantized V requires Flash Attention |
| Flash Attention | on / off | Changes memory and speed; required for quantized V |
| GPU layers | exact integer | Determines offload split |
| `--fit` / `--fit-target` | on/off and target MiB | Affects OOM behavior and layer placement |
| Batch / micro-batch | exact integers | Affects prefill memory |
| Sampler | exact sampler chain, seed | Affects creative variance |
| Reasoning budget | on/off, budget tokens, reasoning format, and a separate answer-space reserve | Conflates with the context budget if unrecorded; the thinking-budget comparison cannot be evaluated without it |

**Verified upstream facts used here (S28, S31, S32):**

- llama.cpp context shift is **off by default** since PR #15416; a test that assumes shift happens is measuring nothing.
- Quantized V-cache requires Flash Attention.
- `--fit-target` defaults to a fixed 1024 MiB; it is not a percentage rule.
- Effective context is commonly far below the claimed window (S15), so a large `-c` is not a substitute for a memory design.

---

## 5. Corpus and scenarios

### 5.1 Primary corpus: synthetic only

No user chat, character card, world, scenario, or lorebook is read without separate authorization.

**Minimum main corpus:**

- 12 synthetic story worlds (10 main + 2 pilot, pilot excluded from analysis).
- 10–14 sessions per world, 1,500–4,000 tokens per session.
- 3–6 characters per world, with at least one entering late, one leaving or dying, and one holding a private fact.
- 10 labelled queries per world = 120 items, across the query families in §5.5.
- At least two worlds exercise each context-flush transition type in §5.4.

### 5.2 Scenario families (one per world, two worlds may share a family)

1. Long-running relationship and promise continuity.
2. Clothing, inventory, injury, and location change.
3. Private knowledge, character entrance, and character exit.
4. Death/absence and derived-surface isolation.
5. World/scenario transition with conversation-local changes.
6. Retcon/correction with an as-of query.
7. Context flush and reconstruction (pre-flush capture, post-flush resume).
8. Cross-conversation contamination.
9. Deletion/forget across raw and derived data.
10. Automated lore noise and over-expansion.
11. Long single-session recall under a fixed window (tests `C0-W`).
12. Conflicting sources (character claim vs world fact vs author note).

### 5.3 Ground-truth ledger

```text
Fact
  fact_id
  subject_character_id
  predicate
  object
  valid_from / valid_to
  status: current | superseded | disputed | deleted
  source_event_id
  evidence_span
  visibility: public | character:<id> | private:<id> | author-only
  derived_from: [fact_id | summary_id | lore_id]   # provenance chain
```

Each query records:

- expected answer, or expected abstention;
- required evidence IDs;
- forbidden evidence IDs (the leakage set);
- allowed answer variants;
- query family and as-of time;
- expected character perspective;
- whether the answer requires a correction/supersession to be honored.

### 5.4 Context-flush fixtures

Because flushing is an intentional Howl design decision (`H1`), each flush world defines:

| Fixture | Content |
|---|---|
| `FL-A` | Flush with no consolidation (pure window slide) |
| `FL-B` | Flush with pre-flush extraction, no reconstruction |
| `FL-C` | Flush with pre-flush extraction and post-flush reconstruction |
| `FL-D` | Flush interrupted mid-consolidation (crash/power-loss model) |
| `FL-E` | Flush at a scene boundary vs mid-scene |
| `FL-F` | Repeated flushes without a checkpoint |

Each fixture has paired queries: one answerable only from pre-flush detail, one answerable only from post-flush state, one requiring both.

### 5.5 Visibility matrix

A character × fact matrix records `knows`, `does_not_know`, `false_belief`, `forbidden`, `unknown`. Safety is not "did the answer avoid a keyword": it is "did the assembled context and the answer avoid every forbidden fact and every derived form of it".

### 5.6 Optional real-data addendum

Real chats may be used only under a separate written authorization naming chats, characters, consent/ownership basis, redaction procedure, local-only storage, deletion date, reviewers, and quotation policy. Without it, real data is out of scope.

---

## 6. Conditions

### 6.1 `C0` correction — five distinct concepts

Revision 1 collapsed five different things into one baseline. `C0` is therefore split. **These are separate conditions and separate measures.**

| ID | Concept | What it means | Measured by |
|---|---|---|---|
| `C0-H` | **Complete history retained by Howl** | Everything the app stores for the conversation, independent of what the model can see. A storage/provenance ceiling, not a prompt condition | `retention completeness`, `history bytes/tokens`, `retrievability of any stored fact` |
| `C0-W` | **Model-visible context** | The exact token sequence handed to the reader, with its role markers and template | `visible tokens`, `template fidelity`, `truncation loss` |
| `C0-S` | **Context selection policy** | The rule that chooses `C0-W` from `C0-H`: recency, lexical, dense, hybrid, summary-index, structured state | `selection recall/precision` (Stage A) |
| `C0-R` | **Context-window limit** | The model's configured `-c` and its *effective* usable context, which is often lower (S15) | `effective-context probe`, `overflow events` |
| `C0-WS` | **Rollover / shifting strategy** | What happens when `C0-W` exceeds `C0-R`: truncate, shift, checkpoint-and-summarize, or refuse | `information loss at transition`, `reconstruction quality`, `recovery after interruption` |

**Why this matters:** revision 1's `C0` ("full available history, truncated to the model window") mixed retention, visibility, and truncation. A result attributed to "full context" could therefore mean any of five different things. No report may use "full context" without naming which of `C0-H`/`C0-W`/`C0-R`/`C0-WS` it means.

**Current Howl behavior, for the record (source-inspected, not runtime-tested):** the live prompt path truncates at approximately `contextSize × 4` characters, keeping the system prompt and newest messages. That is `C0-W` with a `C0-WS` policy of "drop oldest". It is not `C0-H`, and it is not the `memory-guided` rollover the settings file exposes but the chat path does not implement.

### 6.2 Identifier integrity: the `C1` question

**Finding: `C1` is absent, and the omission is undocumented.**

| Check | Result |
|---|---|
| `C1` defined anywhere in MEMORY-03 rev 1 | **No.** The sequence runs `C0`, `C2`, `C3`, `C4`, `C5`, `C6S`, `C6R`, `C7` |
| `C1` referenced by any other project document | **No.** A repository-wide search found no condition-`C1` cross-reference; the `DA-C1` hits in Amendment 02 are decision IDs, unrelated to this protocol |
| Project record documenting an intentional omission | **No** |
| Evidence of a reserved-but-unused slot | **None found** |

**Conclusion:** this is a **documentation defect in revision 1**, not an intentional reserved condition and not evidence that a condition was removed by an owner decision. It cannot be established from the record *why* the label is missing.

**Disposition (chosen to avoid inventing an answer):**

1. `C1` is **retired as unassigned**. It is not reused, and no meaning is retroactively invented for it.
2. Existing identifiers `C0`, `C2`–`C7` are **preserved unchanged** so that any historical reference stays valid.
3. New conditions take new identifiers from `C9` upward. `C8` is likewise unassigned and not reused.
4. Any future attempt to use `C1` or `C8` requires a dated protocol amendment that states what it means and why the gap existed.
5. The registry row for this document records the revision so the byte change is not silent.

This is the only defensible reading of the evidence. If the owner has separate knowledge that `C1`/`C8` were intentional, that knowledge should be recorded as an owner statement, not inferred here.

### 6.3 Retrieval and selection conditions

| ID | Condition | What it tests | Owner hypothesis |
|---|---|---|---|
| `C2` | Lexical/FTS over raw turns or events | Exact-name/term baseline | `H1`, `H2` |
| `C3` | Dense embedding retrieval over raw turns/events | Semantic baseline | `H1` |
| `C4` | Hybrid lexical+dense with score fusion | Hybrid coverage | `H1`, `H2` |
| `C5` | Hybrid + local cross-encoder or reranker | Precision after reranking | `H1` |
| `C6S` | Summary-only memory | Summary control | `H2` |
| `C6R` | Summary-as-index + raw evidence re-injection | Two-stage recovery | `H2` |
| `C7` | Explicit event/state/provenance + hybrid raw evidence | Design hypothesis from `DESIGN-MEMORY-CONTINUITY` | `H4`, `H5` |
| `C9` | Structured state card only (no raw retrieval) | Whether a compact state projection is enough for continuity questions | `H4` |
| `C10` | `C7` + `C6R` (structured state + summary index + raw evidence) | Whether the full combination beats its parts | `H2`, `H4` |
| `C11` | `C7` + rolling pre-flush consolidation (`FL-C` fixtures only) | Whether flush can be exploited as a maintenance boundary | `H1` (`P9`) |

`C7`, `C9`, `C10`, and `C11` are not approved by being tested. If a simpler condition wins, that is a valid and useful result.

### 6.4 Scope conditions

| ID | Scope policy | Enforcement point | What it tests |
|---|---|---|---|
| `S0` | Pooled retrieval, then metadata/visibility filter | Application layer, after search | Retrieve-then-filter |
| `S1` | Visibility/authorization constrains the candidate set before retrieval | Query/store layer | Scope-before-retrieval |
| `S2` | `S1` plus a **structural** gate: the store physically cannot return an out-of-scope row for the given scope key, and derived surfaces carry their own scope key | Storage and query boundary | Structural enforcement |

**`S2` is the condition the evidence actually supports (S8, S9), and it must be implemented in the test store, not asserted in a prompt.** A prompt instruction that says "do not reveal" is `S0` with extra words.

Scope keys must cover, at minimum: character identity, conversation identity, world/scenario identity, public/private content, character presence, character knowledge state, derived summaries, consolidated memories, lore candidates, embeddings, and cached retrieval results.

### 6.5 Automated-lore conditions

Run on the four lore-relevant worlds, plus one control world with lore irrelevant.

| ID | Condition | What it tests |
|---|---|---|
| `L0` | Automated lore disabled | Baseline without generated lore |
| `L1` | Current automated-lore behavior, unchanged | Whether the existing mechanism helps or hurts as-is |
| `L2` | Source-linked, bounded candidate lore with review and a hard token budget | Whether provenance plus bounding fixes the known failure modes |
| `L3` | `L2` integrated with structured memory (`C7`) so lore and memory do not duplicate or contradict | Whether lore belongs inside the memory pipeline or beside it |

`L3` requires `C7` to exist in the run. If `C7` is not run, `L3` is reported as not tested, never as "equivalent to `L2`".

**Lore measures are specified in `MEMORY-04`** and include false-activation rate, precision, token overhead, review burden, contradiction rate with state, and recall of lore-eligible facts. The concrete fixtures are `MEMORY-04` §13 (`E-01`–`E-12`), including the competitive-prose scenario `E-12`, which tests that a lore win does not come at the cost of conversational voice (`C-C`).

### 6.6 Prompt-format conditions (cross-cutting axis)

Prompt format is not a memory condition, but a memory comparison is invalid if formats differ. Full matrix in `MEMORY-06`.

| ID | Format condition | What it tests |
|---|---|---|
| `F0` | Format taken from the model's own embedded template (jinja enabled) | Reference: the model's intended format |
| `F1` | Format rendered by Howl's own builder for that family | Whether Howl constructs the format correctly |
| `F2` | Deliberately mismatched format (generic ChatML for a non-ChatML model) | Size of the format penalty; a negative control |

Every memory condition is run under `F0`. `F1` is run for each sampled family to validate Howl's builder. `F2` is run once per family as a negative control; a large `F2` penalty is evidence that format fidelity matters (`P10`) and that `F1` must not be assumed correct.

### 6.7 Staged matrix

Running every condition against every model is wasteful and would obscure the evidence.

| Stage | Conditions | Models | Items | Purpose |
|---|---|---|---|---|
| **A — selection/retrieval** | `C2`–`C5`, `C6S`, `C6R`, `C7`, `C9`, `C10` × `S0`/`S1`/`S2` | none (no reader) | 120 | Eliminates conditions that cannot surface required evidence |
| **B — reader** | `C0-W`, `C0-WS`, `C2`, `C6S`, `C6R`, `C7`, plus Stage A survivors | `T-A`, `T-B`, `T-C` | 120 | Answer correctness, continuity, temporal accuracy |
| **C — isolation** | `C4`, `C7`, `C10` × `S0`/`S1`/`S2` | all sampled tiers | 60 adversarial probes | Leakage, spoilers, cross-conversation contamination |
| **D — deletion** | `C2`, `C6R`, `C7` × `S2` | all sampled tiers | 40 deletion probes | Raw, derived, cache, and prompt removal |
| **E — lore** | `L0`–`L3` × `S2` | all sampled tiers | 4 worlds × 10 queries | Lore value, noise, cost |
| **F — flush/reconstruction** | `FL-A`–`FL-F` × `C7`, `C11` | `T-A`, `T-B`, `T-C` | 2 worlds × 12 paired items | Information loss, reconstruction quality, interruption recovery |
| **G — prompt format** | `F0`/`F1`/`F2` | every sampled family | 20 format-sensitive probes | Format fidelity and Howl builder correctness |
| **H — creative acceptance** | formation on/off; formation + raw excerpts; **three presentation variants** (verbose / compact-guided / minimal cue) at constant budget; **mandatory transcript-only baseline**; `S2` scope | every sampled tier | 2 worlds × 12 paired scenes, blinded raters | Whether memory-assisted conversation is at least as coherent, natural, and continuous as transcript-only, whether presentation form changes the outcome, and a per-item failure-class distribution (architecture decision §6.9) |

Stage A is a gate: a condition that cannot meet the evidence-retrieval rule does not proceed to Stage B. This is pre-registered, not a post-hoc choice.

**Mapping to the architecture test groups.** The owner architecture decision (`docs/architecture/ARCHITECTURE-DECISION-MEMORY-SYSTEM-2026-11-12.md` §6) defines six acceptance groups. They map onto these stages: Group A → Stages A+B, Group B → Stage F, Group C → the reasoning-budget configuration axis run across A and F, Group D → Stage C, Group E → Stage D, Group F → Stage H. Stage H exists because the architecture decision makes creative performance a critical acceptance test, not a secondary quality measure.

---

## 7. Measures

Twenty measure families. Each result reports raw counts, not only percentages.

### 7.1 Retrieval and selection

`recall@k`, `precision@k`, `nDCG@10`, evidence coverage, false association (right passage, wrong entity), omission at the declared budget, and `selection fidelity` (does `C0-S` choose the items the ledger marks required).

### 7.2 Answer correctness

Full/partial/incorrect against labelled answers and variants; abstention correctness; "correct retrieval, wrong generation" separation (S1 reports 40–50% of errors in this class).

### 7.3 Temporal accuracy

As-of correctness; supersession honored; stale-fact reuse rate (the `FAMA` failure, S6); contradiction with current canonical state; correction propagation to derived surfaces.

### 7.4 Character isolation

Structural leakage (forbidden fact present in the assembled context), behavioral leakage (present in the answer), derived leakage (present in a summary, memory, lore entry, embedding neighbor, or cache), spoiler rate, cross-conversation rate.

### 7.5 Conversation isolation

Same measures, keyed by conversation/world. Includes "same character, different conversation" and "same world, different conversation".

### 7.6 Deletion completeness

Raw deletion, derived deletion, cache deletion, prompt-package deletion, resurrection after rebuild/reload/import, receipt accuracy (does the receipt match reality), and honest-unresolved reporting.

### 7.7 Provenance completeness

Fraction of answer-bearing claims with a resolvable source span; provenance depth (direct source vs summary-of-summary); provenance survival through transformations; provenance shown to the user without requiring a click-through to raw text.

### 7.8 Memory update correctness

Correct insertion, correct supersession, correct retraction, correct no-op, duplicate rate, wrong-merge rate, and update legality (was the write allowed by policy).

### 7.9 Contradiction handling

Detection rate, false-contradiction rate, resolution correctness, and whether the system surfaces the conflict to the user instead of silently picking one.

### 7.10 Creative continuity

Blinded human rubric: voice consistency, callback accuracy, commitment preservation, no invented canon, appropriate uncertainty, scene continuity across a flush.

### 7.11 Relationship continuity

Tracked relationship state (trust, familiarity, promises, debts, injuries, gifts) vs ledger, per character pair.

### 7.12 User correction burden

Corrections per session, time to correct, correction type (delete / edit / supersede / lock), and whether a correction propagates without the user repeating it.

### 7.13 Prompt compatibility

Template fidelity (exact expected token sequence vs emitted sequence), role-marker correctness, end-of-turn token correctness, thinking-channel handling, tool-call round-trip, multimodal structure acceptance, and "silent degradation" (model answers despite a malformed prompt).

### 7.14 Token consumption

Prompt tokens by section (system, state, evidence, summary, lore, history, response reserve), measured with the model's real tokenizer — never estimated from characters.

### 7.15 Latency

Selection/retrieval p50/p95, prefill time, time-to-first-token, total response p50/p95, and consolidation time.

### 7.16 CPU / RAM / VRAM consumption

Peak and steady-state RSS, peak VRAM, post-load free VRAM vs the declared reserve, tokens/sec, and swap/page-file activity.

### 7.17 Storage growth

Bytes per 1,000 events; index bytes; derived-surface bytes; growth ratio of derived to raw; delete cost.

### 7.18 Recovery after interruption

Flush interrupted mid-consolidation (`FL-D`); crash during index write; partial deletion; restart with a stale index; whether state is consistent and whether the system reports the inconsistency.

### 7.19 Robustness across model families

Same condition across families; variance across `T-A`/`T-B`/`T-C`; whether a condition depends on one family's template or tokenizer.

### 7.20 Human review

Two blinded raters per item; correctness, relevance, continuity, leakage, correction burden; time to correct; Cohen's κ target ≥ 0.60; automated judge calibrated against human labels and reported as inconclusive below 0.70 agreement.

---

## 8. Decision rules

All thresholds are **provisional protocol thresholds**. They are ratified in a threshold-review step (§9.4) and frozen before execution.

### 8.1 Hard safety gates (never averaged)

| ID | Rule | Consequence |
|---|---|---|
| `D-SAFE-1` | Zero structural, behavioral, and derived leakage across ≥300 adversarial probes | Any observed leak fails the condition outright |
| `D-SAFE-2` | Zero resurrection across ≥120 deletion probes | Any resurrection fails the condition outright |
| `D-SAFE-3` | Deleted content absent from raw store, derived data, cache, and assembled prompt | Any trace fails the condition outright |
| `D-SAFE-4` | No real user data in the run | Run stopped and destroyed |
| `D-SAFE-5` | A receipt that claims deletion while content remains | Treated as a worse failure than an honest unresolved state |

**No aggregate score, average, or "best overall" ranking may offset a safety-gate failure.** A condition with a failed gate is reported as failed regardless of accuracy.

### 8.2 Evidence and answer gates

The baseline for `D-ANS-*` is fixed at the end of Stage A from the pre-registered retrieval gate, before any reader output is inspected.

| ID | Rule | Success | Failure | Inconclusive |
|---|---|---|---|---|
| `D-RET-1` | Hybrid/rerank vs lexical `recall@10` | ≥ +5 points absolute, lower 95% CI > 0 | ≤ 0, or CI includes 0 | CI crosses |
| `D-RET-2` | Hybrid/rerank vs lexical `nDCG@10` | No worse than −2 points | Worse than −2 | CI crosses |
| `D-ANS-1` | `C7` temporal/correction accuracy vs best non-explicit baseline | ≥ +10 points, lower 95% CI > 0 | ≤ 0 | CI crosses |
| `D-ANS-2` | `C7` overall correctness vs best baseline | Non-inferior within 3 points | Worse than 3 points | CI crosses |
| `D-ANS-3` | `C7` prompt tokens vs `C0-W` full window | ≤ 50% of window tokens | > 50% | Measurement invalid |
| `D-SUM-1` | `C6R` vs `C6S` answer correctness | ≥ +5 points | ≤ 0 | CI crosses |
| `D-SUM-2` | `C6R` vs `C6S` detail retention | ≥ +10 points | ≤ 0 | CI crosses |
| `D-SMALL-1` | `T-A` portability | Within 10 points of `T-B`, zero leaks | > 10 points worse, or any leak | CI crosses |
| `D-FLUSH-1` | `FL-C` vs `FL-A` reconstruction quality | ≥ +15 points on pre-flush-detail items | ≤ 0 | CI crosses |
| `D-FLUSH-2` | `FL-D` interrupted consolidation | State consistent and inconsistency reported | Silent loss or corruption | Harness could not simulate |
| `D-FMT-1` | `F1` vs `F0` correctness | Within 3 points | Worse than 3 points | CI crosses |
| `D-FMT-2` | `F2` negative control | ≥ 5-point penalty vs `F0` (confirms format matters) | No penalty (comparison insensitive; suspect harness) | CI crosses |
| `D-COST-1` | Retrieval p95 vs lexical | ≤ 2× lexical p95 | > 2× | Hardware variance |
| `D-RES-1` | Reserve policy | Post-load free VRAM ≥ declared reserve on every scored run | Any scored run below reserve | Measurement unavailable |
| `D-HUMAN-1` | Inter-rater agreement | κ ≥ 0.60 | κ < 0.40 | 0.40–0.60 adjudicated |

### 8.3 Lore gates (Stage E)

| ID | Rule | Success | Failure |
|---|---|---|---|
| `D-LORE-1` | `L1` vs `L0` answer correctness on lore-eligible items | ≥ +5 points with ≤ 15% token overhead | No gain, or gain only with > 15% overhead |
| `D-LORE-2` | `L1` false-activation rate | ≤ 15% of activations irrelevant | > 30% irrelevant |
| `D-LORE-3` | `L2` vs `L1` precision | ≥ +15 points | ≤ 0 |
| `D-LORE-4` | Lore contradiction with `C7` state | Zero contradictions injected into state | Any contradiction that reaches the prompt as canon |

Between 15% and 30% false activation is inconclusive; it is reported, not rounded.

### 8.4 Overall disposition

Per hypothesis: **Supported** (success rule met, CI excludes null, no safety gate failed), **Not supported** (failure rule met or safety gate failed), **Inconclusive** (CI crosses, sample inadequate, calibration failed, or run invalid).

No "best overall" ranking is produced. Conditions may win different measures; the report must preserve that.

---

## 9. Procedure

### 9.1 Phase 0 — Protocol freeze

1. Freeze this protocol and all appendices.
2. Pin model files, hashes, quants, tokenizers, chat templates, template sources, samplers.
3. Pin hardware, OS, driver, backend commit, offload, KV type, Flash Attention state, reserve policy.
4. Freeze the synthetic corpus, labels, visibility matrix, query set, flush fixtures.
5. Freeze the analysis plan, thresholds, and the statistical method.
6. Record the freeze hash.

No model is loaded in Phase 0.

### 9.2 Phase 1 — Pilot

- 2 worlds, 20 items, `T-A` and `T-B` only.
- Validates fixtures, templates, token budgets, flush fixtures, judge calibration, and human rubric.
- No hypothesis test on pilot data.
- Any change after the pilot requires a dated amendment **before** main data collection.
- Pilot data is excluded from the main analysis.

### 9.3 Phase 2 — Main run

- Stages A–G in order.
- Randomize condition order per item; model order per block.
- Identical evidence budget and prompt scaffold within a comparison.
- Log every prompt, candidate set, token count, latency, model file, seed, and output.
- Three seeds for creative conditions; one deterministic run for factual/extraction.
- Record post-load free VRAM for every scored run.

### 9.4 Threshold review (before freeze, not after results)

Revision 2 inserts an explicit review step, because the owner asked for thresholds to be reviewed before they are frozen.

1. The threshold table (§8) is reviewed by the owner and engineering **before** Phase 0 freeze.
2. The review may raise, lower, or restructure any threshold, and may add or remove a measure.
3. Every change is recorded with a rationale in the protocol's change log.
4. **After freeze, no threshold changes.** A change requires a dated amendment and invalidates any run collected before it.
5. A threshold may not be relaxed because a pilot result was disappointing. A pilot may only change fixtures or the analysis plan, and only before main collection.

### 9.5 Phase 3 — Human review

Two blinded raters; adjudication of disagreements; leakage scored independently of correctness; correction time and type recorded; raters never see the condition label or an expected winner.

### 9.6 Phase 4 — Analysis

Paired comparisons within item and world; bootstrap 95% CIs; pre-specified primary outcomes; Holm-Bonferroni for secondary comparisons; effect sizes and raw counts; null and negative results preserved; no unregistered subgroup analysis.

### 9.7 Phase 5 — Replication

Repeat the core matrix on `HW-3` if available, or on a second pinned model file within the same tier. A result that does not replicate is reported as non-replicated, not averaged into a consensus.

### 9.8 Stop rules

Stop the affected condition if:

- any real user data appears;
- any structural/behavioral/derived leakage occurs in a safety probe;
- a model file/hash, template, or hardware setting changes mid-run;
- more than 5% of runs crash or time out;
- the store corrupts or a deletion receipt cannot be verified;
- a scored run breaches the declared reserve policy;
- the run exceeds the pre-approved compute/time budget.

A stop is a result, not a reason to loosen the test.

---

## 10. Recordkeeping

```text
memory-eval/
  protocol-v2.md
  manifest.json
  corpus/
    world-01/
      events.jsonl
      facts.jsonl
      visibility.json
      queries.jsonl
      flush/
        fl-a.json ... fl-f.json
  labels/human/
  runs/<run-id>/
    prompt.txt
    candidates.jsonl
    output.txt
    metrics.json
    timings.json
    resource.json
  analysis/
    results.json
    report.md
```

The manifest records protocol hash, model identity, template identity, context settings, sampler, seed, hardware, backend commit, offload, KV type, Flash Attention state, reserve policy and observed free VRAM, corpus and label hashes, operator, and every deviation.

**Evidence labels:** `Observed` only for executed runs; `Source-derived` for prior literature; `Inferred` for interpretation; `Unestablished` for untested claims. No result may be labeled `Observed` before execution.

---

## 11. Threats to validity

| Threat | Mitigation |
|---|---|
| Synthetic corpus bias | Multiple families; authored independently of the condition; optional consented real-data addendum |
| Format confound | `F0`/`F1`/`F2`; every memory condition runs the model's own template |
| Tokenizer confound | Real tokenizer counts; never character estimates |
| Judge bias | Human primary; calibrated automated judge; blind review |
| Model stochasticity | Fixed seeds; three seeds for creative conditions; variance reported |
| Quantization confound | Same quant within a tier; separate sensitivity table |
| Hardware confound | Same hardware within-model; `HW-3` replication; reserve policy recorded |
| Reserve-policy artifact | Reserve recorded and reported; breached runs flagged, not averaged |
| Rollover/shift confusion | `C0-WS` declared explicitly; shift default-off confirmed (S28) |
| Flush fixture realism | Six fixture types including interruption; paired pre/post queries |
| Order effects | Randomized condition order |
| Prompt tuning after freeze | Frozen templates; no per-condition tuning |
| Benchmark overfitting | Held-out worlds; no threshold changes after results |
| Multiple comparisons | Pre-specified primaries; Holm-Bonferroni |
| Small sample | Pilot variance estimate; ≥120 items; CIs reported |
| Data leakage | Synthetic primary corpus; isolated store; no real chats |
| Model/harness drift | Hashes and commits pinned; drift invalidates the run |
| Automated-lore confound | `L0`–`L3` separate; never folded into general memory |

---

## 12. Traceability

| Research gap | Protocol coverage |
|---|---|
| `G-01` creative-continuity benchmark | §5.2, §7.10, §7.11 |
| `G-02` character-boundary test | §5.5, `S1`/`S2`, §7.4, `D-SAFE-1` |
| `G-03` deletion/forget test | Stage D, §7.6, `D-SAFE-2/3/5` |
| `G-04` canonical identity | **Not tested**; Layer 2 blocker |
| `G-05` small-model extraction | `T-A`/`T-B`/`T-C`, `D-SMALL-1` |
| `G-06` retrieval at Howl scale | Stage A, §7.1 |
| `G-07` consolidation fidelity | `C6S`/`C6R`, `D-SUM-1/2` |
| `G-08` automated-lore quality | Stage E, `L0`–`L3`, `D-LORE-1`–`4` |
| `G-09` storage/index cost | §7.17 |
| `G-10` temporal/conflict | `C7`, §7.3, §7.9, `D-ANS-1` |
| `G-11` cross-character leakage | Stage C, `S0`–`S2` |
| `G-12` provenance UI | **Out of scope**; separate UX study |
| `G-13` graph/RAG cost | Optional extension; not core |
| `G-14` poisoning | **Out of scope**; separate security protocol |
| `G-15` replication | Phase 5 |
| `G-16` context-flush loss | Stage F, `D-FLUSH-1/2` |
| `G-17` prompt-format fidelity | Stage G, `F0`–`F2`, `D-FMT-1/2` |
| `G-18` hardware/reserve behavior | §4.4, `D-RES-1` |
| `G-19` identifier/label integrity | §6.2 |

**Acceptance links:** `AC-C49`, `AC-C54`, `AC-C55`, `AC-C62`–`AC-C65`, `AC-C71`, `AC-C75`–`AC-C77`.

---

## 13. Authorization gates

| Gate | What it authorizes | Status |
|---|---|---|
| `G0` | Protocol design and threshold review | Design complete; threshold review and freeze require owner approval |
| `G1` | Synthetic corpus, labels, and fixture creation | **Not authorized** |
| `G2` | Model/hardware manifest and harness build | **Not authorized** |
| `G3` | Runtime execution of any stage | **Not authorized** |
| `G4` | Human review and analysis | **Not authorized** |
| `G5` | Publication of results | **Not authorized** |

A gate is passed only by an explicit, dated authorization naming the scope. It is not passed by this document's existence, a successful pilot, a previous gate, owner interest, or the fact that the test is local and uses no user data.

---

## 14. What this protocol does not authorize

- No model load, index build, benchmark, or runtime probe.
- No source, schema, config, database, or user-data change.
- No architecture selection or approval.
- No owner decision on capture, character epistemology, deletion, automated lore, prompt handling, or model envelope.
- No claim that `C7`, `C9`, `C10`, `C11`, or any condition is the Howl architecture.
- No claim that a Layer 1 result is a product result.
- No claim of universal OOM prevention, and no hardware-based model restriction.

**The next legitimate step is owner review of this revision and its thresholds, then a separate authorization for G1/G2. Execution remains blocked until then.**
