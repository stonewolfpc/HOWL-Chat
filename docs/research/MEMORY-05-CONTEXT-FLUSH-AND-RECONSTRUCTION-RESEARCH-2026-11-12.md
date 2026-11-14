# HOWL Chat — Context Flushing and Reconstruction: Research Assessment

**Status:** Research evidence and protocol design only. **No implementation authorization. No architecture selection.**
**Date:** 2026-11-12
**Why this document exists:** Context flushing is an intentional Howl design decision, not an accident to be eliminated. Owner hypothesis `H1` requires investigating rollover, context shifting, KV-cache strategies, and memory-assisted reconstruction as continuity mechanisms. This document records what flushing costs, what it can be made to do, and how to measure it. It is separate from `MEMORY-01` because flushing is a **transition event** with its own failure modes (interruption, partial state, cache invalidation) that a general memory synthesis would treat only in passing.
**Companions:** `MEMORY-01-RESEARCH-SYNTHESIS-2026-11-12.md` (rev 2), `MEMORY-03-EVALUATION-PROTOCOL-2026-11-12.md` (Stage F, `FL-A`–`FL-F`), `MEMORY-07-HARDWARE-AND-INFERENCE-OPTIMIZATION-2026-11-12.md`.
**Controlling records:** `docs/architecture/DESIGN-MEMORY-CONTINUITY-2026-11-11.md`, `docs/review/INVESTIGATION-system-resource-management-architecture-2026-11-06.md`, `docs/architecture/DESIGN-HARDWARE-RESOURCE-ORCHESTRATION-2026-11-11.md`.

**Authority:** Research evidence. It does not alter a frozen contract or an owner decision.

---

## 1. What flushing is in Howl today

Source-inspected facts (not runtime-tested):

| Fact | Where | Consequence |
|---|---|---|
| `context_rollover_mode` is set to `memory-guided` in the app settings | `config/settings.json` | The UI exposes an intent to roll over with memory assistance |
| The live chat path does not implement rollover | `docs/review/INVESTIGATION-system-resource-management-architecture-2026-11-06.md` | The setting is not honoured by the reply path |
| The live prompt path truncates at approximately `contextSize × 4` characters, keeping the system prompt and the newest messages | `app-prompt.go` | The effective behaviour is "drop oldest" with a character-based estimate, not a tokenizer count |
| `context_size` is 20,480 | `config/model-settings.json` | A large window, but effective context is often lower (S15) |
| Context shift is off by default upstream since PR #15416 | S28, S31 | Unless Howl passes an explicit flag, no shift occurs; a design that assumes shift is measuring nothing |
| Quantized V-cache requires Flash Attention | S32 | A KV-cache strategy and an attention setting are coupled |

**Therefore:** "flushing" in Howl currently means *the model's oldest turns fall out of the prompt because the window filled*. It is not a managed transition, has no capture step, no consolidation step, and no reconstruction step, and its size is estimated in characters.

That gap is the research target.

---

## 2. Flushing is a maintenance boundary, not a failure

Three framings are possible, and only one is useful:

| Framing | Implication | Verdict |
|---|---|---|
| **Flushing as loss** | Minimize it; keep as much history as possible | Partially true but unachievable: the window is finite and effective context is smaller than claimed (S15) |
| **Flushing as purely a cost** | Accept overhead; optimize tokens | Incomplete: it treats the transition as an engineering tax |
| **Flushing as a maintenance boundary** | A scheduled point at which the system consolidates state, commits evidence to durable memory, and rebuilds context | **The productive framing** |

A flush is the only reliable moment when the system *knows* history is about to leave the window. That makes it the natural place to:

1. **Capture** what is about to be lost, before it is lost.
2. **Extract** events and state changes from the segment being retired.
3. **Consolidate** those into durable memory with provenance.
4. **Reconstruct** a compact context for the next segment.
5. **Verify** that nothing answer-bearing was lost without a trace.

Framing flushing as a boundary converts an unavoidable cost into the maintenance schedule the memory design needs anyway. Framing it as pure loss produces an endless fight against a finite window.

**This is an inference, not a measured result.** No published benchmark evaluates "flush-aware consolidation" as a technique. It is `Unestablished` and is exactly what Stage F exists to test.

---

## 3. The six stages of a managed flush

| Stage | What happens | Failure if absent | What must be durable |
|---|---|---|---|
| `F-1` Pre-flush capture | Persist the exact retiring segment with IDs, roles, timestamps, template form | Nothing to consolidate from; irreversible loss | Raw segment, byte-exact |
| `F-2` Event extraction | Identify events, state changes, promises, injuries, item/location changes, relationship shifts | State silently stale; retcons lost | Event records with source spans |
| `F-3` State update | Apply supersession/retraction to the projection, with validity intervals | Contradictory state; stale facts resurface (S6) | Projection plus change log |
| `F-4` Consolidation | Merge duplicates, summarize where useful, keep raw evidence reachable | Token bloat or detail loss (S2) | Summary plus pointers to raw |
| `F-5` Reconstruction | Build the next context: state card, recent turns, retrieved evidence, character-specific view | Cold restart; character forgets relationship | Assembled package plus its scope keys |
| `F-6` Verification | Check that every ledger fact needed for continuity is either in state or retrievable | Silent loss discovered only by the user | Loss report / receipt |

### 3.1 Ordering constraint

`F-1` must complete **before** the window is truncated. `F-3` must complete before `F-5` reads state, or reconstruction reads a stale projection. `F-6` must run after `F-5` and must be able to fail loudly.

**A flush that truncates first and captures later is not a managed flush.** It is data loss with extra steps.

---

## 4. Reconstruction: what has to come back

Continuity after a flush is not "the summary of the old segment". It is a set of distinct things, and each has a different source:

| Restoration target | Source | Why a summary alone fails |
|---|---|---|
| Current scene state (place, time, participants, injuries, items) | `F-3` projection | Summaries describe events, not current state |
| Relationship state per character pair | `F-3` projection | Trust/debt/promise state is rarely narrated after the fact |
| Active commitments and promises | Event records with open status | The exact wording matters; paraphrase loses the obligation |
| Private knowledge per character | Visibility-tagged facts | A pooled summary leaks private facts (S7–S9) |
| Character voice and recent emotional trajectory | Recent raw turns + style anchors | Voice is lost in paraphrase |
| Open plot threads | Event records with `open` status | Summaries close threads that are still open |
| Exact details on demand (names, numbers, appearances) | Raw evidence retrieval | Summaries lose answer-bearing detail (S2) |
| Conversation-local conventions (nicknames, inside jokes, established facts) | Conversation-scoped memory | Global summaries cannot carry per-conversation convention |

**Character-specific restoration is the hardest part.** A single reconstruction package shared by all characters cannot satisfy `H3`. Reconstruction must be assembled **per audience**, from a scope-constrained candidate set (`S2`), not filtered afterward.

---

## 5. Cache reuse across a flush

A flush changes the prompt prefix, which interacts with inference-engine caching.

| Mechanism | Behaviour | Constraint |
|---|---|---|
| **Prompt-prefix caching** (server-side) | Reuses computed KV for an unchanged prefix; a changed early block invalidates it | Inserting memory or lore early in the prompt destroys reuse |
| **Slot save/restore** (server-side) | Serializes a slot's state to disk and restores it | Restoring a slot is not the same as restoring conversation state; it is a performance device |
| **`--cache-reuse` chunking** | Reuses matching chunks to salvage part of a shifted prompt | Salvage is partial and implementation-defined; must be verified against the pinned backend |
| **Context shift** | Discards the oldest tokens and keeps the rest of the cache | Off by default upstream (S28); changes positions and can contradict the memory design's own record |
| **Howl-level assembly cache** | Caches a rendered context package | Must be invalidated on state change, deletion, or scope change |

**Design consequence:** the *order* of prompt sections is a performance decision, not only a formatting one. If state and evidence blocks are placed before history, every turn's change invalidates the cache. If they are placed after the stable history prefix, caching helps but position effects (S19) must be checked. `MEMORY-06` owns the template-safe placement rules.

**Unestablished:** the actual cache-hit rate Howl achieves. It must be measured, not assumed. A cache hit rate claim without a measurement is a guess.

---

## 6. Recovery from interrupted consolidation

A flush is a multi-step write. Any step can be interrupted by a crash, power loss, forced quit, or a model call that never returns.

| Interruption point | Resulting state | What the system must do |
|---|---|---|
| During `F-1` capture | Partial segment persisted | Detect the incomplete segment, re-capture from the source of truth, mark the run |
| During `F-2` extraction | Some events extracted, others not | Idempotent extraction keyed by segment and event ID; re-run without duplicating |
| During `F-3` state update | Projection partially updated | Transactional apply per event batch; replay from the change log |
| During `F-4` consolidation | Orphan summary or half-merged duplicate | Candidate state, not canon, until committed |
| During `F-5` reconstruction | Context assembled without state | Refuse to answer from an unverified assembly, or answer with an explicit degraded flag |
| After `F-6` | Receipt written but content inconsistent | Verification must re-read, not trust the write |

**Rules:**

1. Every step is idempotent and keyed by stable IDs.
2. Nothing is deleted until its replacement is durably committed.
3. A partial flush is a **reported** state, never a silent one.
4. The model must never be asked to answer from a context whose assembly is unverified, because a confident answer from a half-built context is worse than an error.

**Evidence status:** the general pattern (write-ahead, idempotent steps, commit-then-delete) is standard engineering practice, not a memory-research result. Its application to flush consolidation is `Inferred`.

---

## 7. Information-loss accounting

The only way to know whether a flush is safe is to account for what left the window.

| Category | Question | Measure |
|---|---|---|
| Captured | Was every retiring turn persisted byte-exactly? | `capture_completeness` = captured tokens ÷ retiring tokens |
| Extracted | Did extraction cover the events the ledger requires? | `extraction_recall` against the fact ledger |
| Projected | Did every current fact survive into state? | `state_completeness` |
| Retrievable | Can a flushed fact still be retrieved when asked? | `post_flush_recall@k` |
| Reconstructed | Does the assembled context contain what the next turn needs? | `reconstruction_coverage` |
| Silently lost | Which ledger facts are in no durable surface? | `silent_loss_count` — must be zero for scored conditions |
| Recoverable | Can a lost fact be recovered from raw source? | `recoverability` |
| Explainable | Can the system say *why* a fact is absent? | `loss_receipt_accuracy` |

**`silent_loss_count` is the key number.** A flush that loses facts but reports the loss is a bounded failure. A flush that loses facts silently violates the continuity requirement and, in the deletion direction, violates `H6`.

---

## 8. Failure modes

| Failure | Cause | Evidence | Mitigation |
|---|---|---|---|
| **Summary drift** | Repeated consolidation with no raw anchor | S2, S14 | Raw evidence reachable from every summary |
| **State staleness** | `F-3` skipped or partially applied | S6 | Transactional projection + change log |
| **Cold restart** | No reconstruction; model starts from system prompt only | Direct reasoning; no benchmark | Reconstruction package per audience |
| **Character amnesia** | Shared reconstruction drops private facts | S7–S9 | Per-audience assembly (`S2`) |
| **Cache thrash** | Early-prompt changes invalidate prefix | S28 mechanical | Section ordering + measured hit rate |
| **Double-consolidation** | Interrupted then re-run extraction | Inherent | Idempotent keys |
| **Orphaned derivation** | Summary survives deletion of its source | S10 | Deletion propagation |
| **Over-compression** | Budget pressure pushes summary-only | S1, S2 | Hard rule: summaries index, raw evidence answers |
| **Position artifact** | Reconstructed block placed where the model ignores it | S19 | Format/position tests (`MEMORY-06`) |
| **Non-determinism** | Map iteration / timing-dependent ordering | Inherent | Deterministic assembly contract |

---

## 9. Cost model

Flushing as a maintenance boundary has a cost, and it must be measured rather than asserted.

| Cost | Driver | How to measure |
|---|---|---|
| Extraction calls | Segment size, schema breadth | Calls per 1k tokens; total tokens |
| Consolidation calls | Duplicate rate, summary cadence | Tokens per flush |
| Added latency | Synchronous vs background flush | Time-to-first-token impact |
| Prompt tokens | State card + evidence + summaries | Tokenizer counts by section |
| Cache loss | Prefix invalidation | Measured hit rate before/after |
| Storage | Raw + derived growth | Bytes per 1k events |
| User-visible stall | Blocking flush | Seconds of blocked UI |

**A flush that is correct but blocks the reply for several seconds on `HW-1` is a product problem even if its quality is good.** The measurement must include the stall, not only the quality.

---

## 10. What to test (Stage F)

| Fixture | Purpose | Primary measure |
|---|---|---|
| `FL-A` no consolidation | Baseline: pure window slide | `post_flush_recall`, `silent_loss_count` |
| `FL-B` capture + extraction, no reconstruction | Isolates the value of capture | `state_completeness` |
| `FL-C` capture + extraction + reconstruction | Full managed flush | `reconstruction_coverage`, answer correctness |
| `FL-D` interrupted mid-consolidation | Recovery | State consistency, honest reporting |
| `FL-E` scene-boundary vs mid-scene flush | Placement sensitivity | Continuity rubric |
| `FL-F` repeated flushes without checkpoint | Compounding loss | Detail retention across ≥5 flushes |

Paired items per fixture: one answerable only from pre-flush detail, one only from post-flush state, one requiring both.

**Thresholds** (`D-FLUSH-1`, `D-FLUSH-2`) are in `MEMORY-03` §8.2 and are provisional until ratified.

---

## 11. Unresolved questions

| ID | Question | Status |
|---|---|---|
| `Q-F1` | Does flush-aware consolidation beat plain truncation on continuity? | Unestablished; Stage F |
| `Q-F2` | What is the maximum flush frequency before compounding loss dominates? | Unestablished (`FL-F`) |
| `Q-F3` | Is per-audience reconstruction affordable at `T-A` token budgets? | Unknown; cost model §9 |
| `Q-F4` | Can reconstruction be assembled without invalidating the prompt cache? | Unknown; needs measurement |
| `Q-F5` | Should flushing be synchronous or background? | Product decision; depends on `Q-F1` and stall tolerance |
| `Q-F6` | What is the correct unit to consolidate: turn, scene, or session? | Unknown; affects schema |
| `Q-F7` | How should a flush interact with a mid-generation abort? | Unknown |
| `Q-F8` | Does context shift belong in Howl at all, given upstream default-off and the memory design's own record? | Open design question, not this document's decision |
| `Q-F9` | What is the user-visible signal that a flush happened? | UX; out of scope here |
| `Q-F10` | Can a flush receipt be verified independently of the code that wrote it? | Needs design |

---

## 12. What this document does not do

- It does not authorize implementing a flush pipeline.
- It does not change `context_rollover_mode`, the context size, or the truncation behaviour.
- It does not claim any measured benefit for flush-aware consolidation.
- It does not decide between shift, truncate, and checkpoint.
- It does not touch user chats or run a model.

**The claim being made is narrow and testable: a flush is the predictable moment at which Howl can maintain its memory, and it should be measured as such rather than treated only as overhead.**
