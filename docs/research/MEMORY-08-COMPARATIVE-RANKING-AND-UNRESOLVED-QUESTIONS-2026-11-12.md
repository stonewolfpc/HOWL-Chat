# HOWL Chat — Comparative Ranking and Unresolved Questions

**Status:** Evidence-based comparison only. **No architecture selected. No implementation authorization.** Nothing here is an approval, a design freeze, or a decision.
**Date:** 2026-11-12
**Why this document exists:** The work order requires an evidence-based ranking of the available approaches, a classification of how strong the evidence is for each, and a consolidated register of unresolved questions. `MEMORY-01` holds the evidence narrative and `MEMORY-03` holds the test design; this document holds the **comparison and the open-question ledger**, so that the owner and engineering can see the whole option space in one place without reading five documents.
**Companions:** `MEMORY-01` (rev 2), `MEMORY-02` (sources), `MEMORY-03` (protocol rev 2), `MEMORY-04` (automated lore), `MEMORY-05` (context flush), `MEMORY-06` (prompt formats), `MEMORY-07` (hardware).

**Authority:** Research evidence. It does not alter a frozen contract, an owner decision, or an amendment.

---

## 1. How to read the ranking

### 1.1 What is ranked

Eighteen **approaches** — techniques and technique-groups, not products. A product is not ranked because product quality changes with version and because vendor numbers are not portable (`MEMORY-01` §2.5).

### 1.2 Ranking dimensions

Per the work order: functional suitability, evidence quality, compatibility, interoperability, resource efficiency, scalability, reliability, privacy, deletion guarantees, user control, implementation complexity, maintenance burden, creative continuity benefit.

Each is graded **Strong / Adequate / Weak / Unknown**. Grading rules:

- **Strong** requires either independent replication, a peer-reviewed result, or a documented mechanism plus a direct owner-requirement match.
- **Adequate** requires official documentation or a single credible source plus a plausible mechanism.
- **Weak** means the mechanism is plausible but unsupported, or the evidence is vendor/preprint-only and contradicted elsewhere.
- **Unknown** means no basis to grade.

**Grades are about evidence and fit, not about benchmark scores.** No cross-benchmark number is used, because they are not comparable (`MEMORY-01` §2.1).

### 1.3 Classification

| Class | Meaning |
|---|---|
| **Established** | Mechanism is well documented and directly usable; the residual question is tuning, not viability |
| **Strongly supported** | Multiple credible sources or a peer-reviewed result, with a direct match to an owner requirement |
| **Promising** | Credible mechanism and partial evidence; needs Howl-specific testing |
| **Conditional** | Works only under stated conditions (a dependency, a budget, a data property); otherwise unsuitable |
| **Weakly supported** | Evidence is thin, vendor-only, or contradicted; may still be worth testing |
| **Incompatible** | Conflicts with a stated owner requirement or a hard constraint |
| **Unestablished** | No adequate evidence either way |

---

## 2. Ranked comparison

Rank is by **weighted fit to Howl's stated requirements**, with evidence quality as a tie-breaker. Rank is not a recommendation to implement.

| Rank | Approach | Class | Functional | Evidence | Compat. | Interop. | Resource | Scale | Reliable | Privacy | Deletion | User control | Complexity | Maint. | Continuity |
|---:|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| 1 | **A1** Append-only event log + derived state, with validity/supersession and provenance | **Strongly supported** | Strong | Adequate | Strong | Strong | Adequate | Strong | Adequate | Strong | Adequate | Adequate | Weak (high) | Weak (high) | Strong |
| 2 | **A11** Character/conversation scope enforced before retrieval, structurally at the store | **Strongly supported** | Strong | Strong | Strong | Strong | Strong | Strong | Strong | Strong | Adequate | Adequate | Adequate | Adequate | Strong |
| 3 | **A3** Episodic raw evidence retrieval (lexical/FTS baseline) | **Established** | Strong | Strong | Strong | Strong | Strong | Strong | Strong | Strong | Adequate | Strong | Strong | Strong | Adequate |
| 4 | **A5** Hybrid lexical + dense + rerank | **Strongly supported** | Strong | Adequate | Adequate | Strong | Adequate | Adequate | Adequate | Strong | Adequate | Strong | Adequate | Adequate | Adequate |
| 5 | **A6** Summary-as-index with raw-evidence re-injection | **Strongly supported** | Strong | Adequate | Strong | Strong | Strong | Strong | Adequate | Adequate | Adequate | Strong | Adequate | Adequate | Strong |
| 6 | **A12** Flush-aware consolidation and reconstruction | **Promising** | Strong | Weak | Adequate | Strong | Adequate | Strong | Weak | Strong | Adequate | Adequate | Weak (high) | Weak (high) | Strong |
| 7 | **A14** Prompt-format fidelity layer | **Established** (mechanism) | Strong | Strong | Strong | Strong | Strong | Strong | Strong | Strong | Strong | Strong | Adequate | Adequate | Adequate |
| 8 | **A15** Hardware-aware residency, reserve policy, pre-flight | **Established** (mechanism) | Adequate | Adequate | Strong | Adequate | Strong | Adequate | Adequate | Strong | Strong | Strong | Adequate | Adequate | Adequate |
| 9 | **A4** Dense embedding retrieval | **Promising** | Adequate | Adequate | Adequate | Strong | Adequate | Adequate | Adequate | Adequate | Weak | Strong | Adequate | Adequate | Weak |
| 10 | **A13** Automated lore (keyword, semantic, or state-conditioned) | **Conditional** | Adequate | Weak | Strong | Weak | Adequate | Adequate | Weak | Adequate | Weak | Adequate | Strong | Weak (high) | Adequate |
| 11 | **A2** Entity/fact graph with temporal edges | **Conditional** | Strong | Adequate | Weak | Adequate | Weak | Adequate | Adequate | Adequate | Weak | Adequate | Weak (high) | Weak (high) | Adequate |
| 12 | **A8** Hierarchical summaries / tree retrieval (RAPTOR-class) | **Conditional** | Adequate | Weak | Adequate | Adequate | Adequate | Strong | Weak | Adequate | Weak | Adequate | Adequate | Adequate | Adequate |
| 13 | **A16** Tool-mediated memory (model edits its own memory) | **Conditional** | Adequate | Weak | Adequate | Weak | Strong | Strong | Weak | Adequate | Weak | Adequate | Adequate | Adequate | Adequate |
| 14 | **A9** Agent-memory framework (MemGPT/Mem0/Zep/LangMem-class) | **Weakly supported** | Adequate | Weak | Weak | Weak | Adequate | Strong | Adequate | Weak | Weak | Adequate | Weak (high) | Weak (high) | Adequate |
| 15 | **A18** Full-context window management alone (`C0-W`/`C0-WS`) | **Established as insufficient** | Weak | Strong (against) | Strong | Strong | Adequate | Weak | Adequate | Strong | Adequate | Strong | Strong | Strong | Weak |
| 16 | **A17** GraphRAG / community summaries as a default | **Weakly supported** | Weak | Weak | Weak | Weak | Weak | Adequate | Adequate | Adequate | Weak | Adequate | Weak (high) | Weak (high) | Weak |
| 17 | **A7** Summary-only memory | **Incompatible** | Weak | Strong (against) | Strong | Weak | Strong | Strong | Weak | Adequate | Weak | Adequate | Strong | Adequate | Weak |
| 18 | **A10** RL-trained memory manager (ADD/UPDATE/DELETE/NOOP policy) | **Unestablished** | Unknown | Weak | Unknown | Unknown | Unknown | Unknown | Unknown | Unknown | Unknown | Unknown | Weak (high) | Weak (high) | Unknown |

### 2.1 Why the top of the list is where it is

- **A1** is first because it is the only approach that simultaneously serves `H4` (temporal validity), `H5` (provenance), and `H6` (deletion receipts) at the storage layer rather than in a prompt. Its weakness is design and maintenance cost, not capability. Evidence is adequate, not strong: APEX-MEM, Zep, and Memora support the mechanism, but none is independently replicated and none is a local-model result.
- **A11** is second and is the only approach graded **Strong** on privacy and reliability together. Authorization-first retrieval has structural evidence (0% vs 86.1% structural exposure, S9) and is a hard owner requirement (`H3`). It is a *gate*, not a feature; it belongs in every combination regardless of rank.
- **A3** is third because it is the mandatory baseline that many memory products lose to (S4). Its grade is `Established` and its cost is the lowest in the table.
- **A5** and **A6** are the two credible quality improvements over A3, and they compose with A1 and A11 rather than competing with them.
- **A12** is sixth despite weak direct evidence because flushing is already Howl's design decision and the alternative is unmanaged loss. It is `Promising`, not `Strongly supported`: no benchmark isolates it.

### 2.2 Why the bottom of the list is where it is

- **A7** (summary-only) is **Incompatible**, not merely weak: `H2` states summary-only is insufficient, and S1/S2 show detail loss. That is a requirement conflict, which is the only legitimate reason to place a mechanism last.
- **A17** (GraphRAG as a default) is not impossible — it is `Weakly supported` for Howl's scale and local-model constraints. It remains admissible as an optional extension (`G-13`).
- **A10** is `Unestablished`: no independent replication, benchmark-specific training. It is not eliminated; it is simply not a basis for a decision.
- **A9** is `Weakly supported`, not eliminated: framework numbers are not portable, but a framework may still be a useful implementation starting point. The class reflects evidence, not a ban.
- **A18** is `Established as insufficient`: long context is necessary headroom, not memory. It remains a baseline condition, not a candidate architecture.

**Reframe after the 2026-11-12 owner decision (`OWNER-DECISION-MEMORY-FORMATION`):** `A13` now means the **organization/presentation layer** (contextual collections of retained knowledge), not a formation method. Automated lore *generation* is subsumed by `A1`/`A6` and is no longer separately ranked. The rank of `A13` therefore reflects the value of collections for organizing and presenting knowledge, not the value of a second extractor — which this ranking does not endorse.

---

## 3. Interoperability and conflicts

### 3.1 Complements

| Pair | Why they compose |
|---|---|
| A1 + A6 | State projection answers "what is true now"; summary-as-index routes to the raw evidence that proves it |
| A1 + A3 | Provenance needs the raw span; the state needs the event |
| A3 + A4 → A5 | Lexical and dense cover different failure modes; reranking restores precision |
| A1 + A11 | Scope keys live on events, facts, and derived surfaces, so the gate can be structural |
| A6 + A11 | A summary must carry its own audience key or it becomes a leak channel |
| A12 + A1 | A flush is the natural commit point for consolidation |
| A12 + A6 | Reconstruction needs both state and a route to evidence |
| A14 + everything | Format fidelity is a precondition; a wrong template confounds every comparison |
| A15 + everything | Residency and reserve determine whether any of it runs at all |
| A13 + A6 | Lore as a detail index beside summaries, if Stage E shows value |
| A2 + A1 | A graph is a projection of the event log; it is not a separate truth source |

### 3.2 Conflicts

| Pair | Conflict | Resolution |
|---|---|---|
| A4/A5 vs A3-only | Dense retrieval adds a model, an index, and versioning | Keep A3 as the baseline; adopt A4/A5 only if `D-RET-1` passes |
| A6 vs A7 | Summary-as-index and summary-only are often conflated | Never merge them; `D-SUM-1/2` exists to keep them apart |
| A8 vs A6 | Two summary hierarchies compete for the same budget | Pick one summary structure; do not run both |
| A13 vs A1 | Two truth sources can contradict | Lore must be derived detail, never canon |
| A13 vs A11 | Global lore entries leak across worlds | Scope key on every entry |
| A12 vs prompt caching | Reconstruction changes the prefix | Declare section order; measure hit rate |
| A16 vs A11 | A model that writes memory can write out of scope | Never let the model choose the scope key |
| A17 vs local budget | Graph construction is LLM-heavy | Optional extension only |
| A9 vs A1 | A framework's internal store may not expose provenance/deletion | Framework adoption would need a provenance/deletion audit first |

### 3.3 Dependency order (not a build order)

```text
A14 (format fidelity)  ─┐
A15 (residency/reserve) ─┤ preconditions for any measurement
A11 (scope gate)        ─┤ hard requirement, applies to every layer
                        ▼
A3 (lexical baseline) → A5 (hybrid+rerank) → A6 (summary-as-index)
                        ▼
A1 (event/state/provenance) → A12 (flush consolidation)
                        ▼
A13 (lore, only if Stage E shows value)
```

---

## 4. Elimination log

Only approaches that are genuinely impossible, incompatible, unsafe, or demonstrably unsuitable are eliminated. Complexity is not grounds for elimination.

| Approach | Status | Grounds |
|---|---|---|
| Summary-only memory as the *complete* solution | **Eliminated as a complete solution** | Conflicts with `H2`; detail loss is documented (S1, S2). Remains valid as an index |
| Pooled retrieval with post-hoc filtering as a *security boundary* | **Eliminated as a security control** | Structural exposure 86.1% vs 0% (S9); metadata filtering insufficient (S8). May remain as a ranking heuristic |
| Long context alone | **Eliminated as a memory architecture** | Degrades with length (S1, S15, S19); headroom is not memory |
| Unbounded recursive lore activation | **Eliminated as a configuration** | Budget exhaustion and non-determinism (S38 mechanics); admissible only with hard caps |
| Universal OOM prevention claim | **Eliminated as a claim** | Not technically achievable (S28, S31) |
| Vendor benchmark numbers as a selection basis | **Eliminated as evidence** | Same system scores vary widely across harnesses (S3, S4, S23, S24) |

**Everything else remains on the table**, including the custom-built option. No existing framework has been shown to satisfy `H1`–`H8` together, and none has been shown to be incapable either.

---

## 5. Fact, inference, recommendation

| Category | Content |
|---|---|
| **Fact (source-derived)** | Every claim in §2's Evidence column, §3's conflicts, and §4's grounds, traced to S1–S45 |
| **Inference** | The rank order, the grades, the dependency order, and the claim that A1/A11/A3/A5/A6 compose |
| **Not established** | That this combination works on Howl's tasks, models, or hardware; that any rank predicts Howl performance |
| **Recommendation (narrow)** | Test the top-ranked set first, in the dependency order, with A11 as a gate rather than a variable |
| **Not recommended** | Any architecture selection, any implementation, any runtime change, any threshold change after results |

---

## 6. Consolidated unresolved-questions register

Priority: `P0` blocks any evaluation; `P1` blocks a credible result; `P2` blocks a product decision.

### 6.1 Identity, storage, and product path

| ID | Question | Priority | Blocks |
|---|---|---|---|
| `Q-01` | What is the canonical root/chat identity? | **P0** | Layer 2, scope keys, deletion, migration (`ROOT-01`) |
| `Q-02` | Which storage adapter holds events, state, and indexes? | **P0** | Every Layer 2 result (`DA-M1`, `DA-R1`) |
| `Q-03` | Is memory/RAG reachable from the live reply path at all? | **P0** | Any claim that a technique affects replies |
| `Q-04` | What is the deletion/retention policy across derived surfaces? | **P0** | `H6` verification (`DA-M5`) |
| `Q-05` | What is the conflict/automation policy for memory writes? | **P1** | `DA-M6` |
| `Q-06` | Can conversation branching be represented without corrupting provenance? | **P2** | `DA-M4` |

### 6.2 Evaluation design

| ID | Question | Priority | Blocks |
|---|---|---|---|
| `Q-07` | Are the provisional thresholds in `MEMORY-03` §8 acceptable to the owner and engineering? | **P0** | Protocol freeze (`G0`) |
| `Q-08` | Is the synthetic corpus realistic enough to stand in for creative continuity? | **P1** | External validity |
| `Q-09` | What is the minimum sample that yields a usable confidence interval? | **P1** | Analysis plan |
| `Q-10` | How should the automated judge be calibrated, and against how many human labels? | **P1** | Secondary measures |
| `Q-11` | Does the `C1`/`C8` gap reflect owner knowledge not in the record? | **P2** | Label integrity (`MEMORY-03` §6.2) |

### 6.3 Prompt format

| ID | Question | Priority | Blocks |
|---|---|---|---|
| `Q-P1` | Does `--chat-template` without `--jinja` apply a custom template in the pinned backend? | **P1** | `F1` validity |
| `Q-P2` | Which families accept a system message after the first turn? | **P1** | Memory injection contract |
| `Q-P3` | How should system content be folded for Gemma, Llama 2, DeepSeek? | **P1** | Format fidelity |
| `Q-P4` | Is the legacy hardcoded-ChatML builder reachable? | **P1** | Confound control |
| `Q-P5` | What replaces character-based truncation? | **P1** | Token accounting |
| `Q-P6` | Which reasoning format for non-DeepSeek reasoning models? | **P2** | Budget accounting |

### 6.4 Context flush

| ID | Question | Priority | Blocks |
|---|---|---|---|
| `Q-F1` | Does flush-aware consolidation beat plain truncation? | **P1** | `H1` |
| `Q-F2` | Maximum flush frequency before compounding loss dominates? | **P1** | Design |
| `Q-F3` | Is per-audience reconstruction affordable at small-model budgets? | **P1** | `H3` + cost |
| `Q-F4` | Can reconstruction avoid invalidating the prompt cache? | **P2** | Cost |
| `Q-F6` | What is the consolidation unit: turn, scene, or session? | **P2** | Schema |

### 6.5 Automated lore

| ID | Question | Priority | Blocks |
|---|---|---|---|
| `Q-L1` | What is the current false-activation rate? | **P1** | `H7` |
| `Q-L2` | Does lore improve correctness or only style? | **P1** | Retention vs retirement |
| `Q-L3` | Appropriate lore budget for the smallest sampled reader? | **P1** | `P7` |
| `Q-L5` | Candidate or canon for generated lore? | **P2** | `H5` |
| `Q-L7` | Lore inside the memory pipeline or beside it? | **P2** | Two-truth-source risk |

### 6.6 Hardware and inference

| ID | Question | Priority | Blocks |
|---|---|---|---|
| `Q-HW1` | Actual `HW-1` VRAM and post-load free VRAM per quant? | **P0** | Every cost measure |
| `Q-HW2` | What does `kv_cache_size: 512` mean? | **P0** | Configuration validity |
| `Q-HW3` | Is `kv_cache_type: flash-attention` with `flash_attention: false` coherent? | **P0** | Configuration validity |
| `Q-HW4` | Is 10% the right reserve, or a MiB floor? | **P1** | `D-RES-1` |
| `Q-HW8` | Can embedding and chat models be co-resident? | **P1** | Retrieval feasibility |

### 6.7 Missing research

| ID | Question | Priority | Blocks |
|---|---|---|---|
| `Q-M1` | No peer-reviewed measurement of lore precision/recall in any product | **P1** | External benchmark reference |
| `Q-M2` | No replication of 2026 preprint/vendor memory results | **P2** | Evidence confidence |
| `Q-M3` | No creative-continuity benchmark exists; must be built for Stage B/F | **P1** | `G-01` |
| `Q-M4` | Poisoning of memory via imported cards/documents is unmodelled | **P2** | Security |

---

## 7. What would change this ranking

| Observation | Effect |
|---|---|
| `D-RET-1` fails (hybrid ≈ lexical) | A4/A5 drop; A3 rises; the memory design simplifies |
| `D-SUM-1/2` pass | A6 rises to the top tier of quality techniques |
| `D-ANS-1` fails (explicit state no better) | A1's temporal rationale weakens sharply |
| `D-FLUSH-1` fails | A12 drops to `Weakly supported`; flushing reverts to a pure cost |
| `D-LORE-1`/`D-LORE-2` fail | A13 moves to `Incompatible` for the current design; retirement becomes evidenced |
| `D-FMT-2` shows no penalty | The format hypothesis weakens and `MEMORY-06` loses priority |
| A framework passes a provenance/deletion audit | A9 rises from `Weakly supported` |
| `HW-1` cannot host the retrieval stack under the reserve | A4/A5 become conditional on hardware, and A3/A1 gain |

**No ranking here is final.** It is a starting order for testing, and the tests are designed to be able to overturn it.

---

## 8. What this document does not do

- It does not select an architecture or a framework.
- It does not authorize implementation, benchmarking, or model loading.
- It does not close any open decision (`DA-M*`, `DA-R*`, `R-10`, `R-11`).
- It does not claim a rank predicts Howl performance.
- It does not prefer custom-built or off-the-shelf; both remain open on the evidence.
