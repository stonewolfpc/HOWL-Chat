# HOWL Chat — Memory Evaluation Pre-Registration Package

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**
>
> Every arm, rubric, gate, and fixture below exists to test that sentence against a transcript-only baseline. A configuration that wins on volume has won nothing.

**Status:** **APPROVED — governing pre-registration authority** (owner, 2026-11-12; see `docs/architecture/OWNER-APPROVAL-MEMORY-09-AND-G1-2026-11-12.md`). **`G1` (corpus and labels) is authorized; `G2`–`G5` remain closed.** No model has been loaded and no run executed. This document fixes what would be tested and how, so that results cannot be reinterpreted after the fact.
**Revision:** 2 — approval recorded; safeguards `R-11`–`R-14` locked; principle-test requirement and tuning/held-out split added (owner, 2026-11-12).
**Date:** 2026-11-12
**Why this document exists:** Three records now describe the same programme from three angles: `ARCHITECTURE-DECISION-MEMORY-SYSTEM` (acceptance philosophy and gates), `MEMORY-03` (execution protocol), and `MEMORY-01`–`MEMORY-08` (evidence and fixtures). They overlap in phase numbering, test naming, and threshold policy. This package is the **single reconciliation**: it resolves those conflicts, fixes the corpus and arm definitions, assigns each test to application-side or model execution, and states exactly what evidence is required before any unfrozen architectural choice is made.
**Authority:** Pre-registration design. It does not authorize execution (`G1`–`G5` remain closed), select an architecture, change a schema, or modify a frozen artifact.
**Sequence:** protocol (this) → corpus → controlled experiments → architecture decisions.

---

## 1. Reconciliation of the three records

| Record | Owns | Does not own |
|---|---|---|
| `ARCHITECTURE-DECISION-MEMORY-SYSTEM` | The governing principle, the four layers, capability dispositions, acceptance philosophy, acceptance gates, failure-class taxonomy | Execution detail, corpus, arm definitions |
| `MEMORY-03` (protocol rev 2) | Stages, conditions, measures, configuration axes, authorization gates `G0`–`G5` | Acceptance philosophy, corpus representativeness, pre-registration rules |
| `MEMORY-09` (this) | Pre-registration decisions: rules lock, corpus spec, arms, variable control, rubrics, gate table, execution allocation, freeze-evidence checklist | New gates, new principles, implementation |

### 1.1 Conflicts resolved

| Conflict | Resolution |
|---|---|
| Three phase numbering schemes (architecture §8, `MEMORY-03` §9, owner programme Phases 1–5) | **One canonical sequence** in §9. `MEMORY-03`'s internal phases map onto it; the architecture record's §8 is a summary of it |
| "Test groups A–F" vs "Stages A–H" | Groups are **acceptance** categories; stages are **execution** units. The mapping table in `MEMORY-03` §6.7 is authoritative |
| Thresholds "provisional" vs "not invented" | Thresholds are set **after** the pilot establishes variance and **before** final validation, by a stated procedure (§10). No number is set in this document |
| Architecture record says "no schema frozen"; `MEMORY-03` defines conditions that presume a schema | The conditions in `MEMORY-03` operate on **abstract record shapes**; the concrete schema remains unfrozen until the evidence checklist in §8 is satisfied |
| `MEMORY-04` §13 fixtures `E-01`–`E-12` vs the programme's nine test families | The twelve fixtures are **retained as required members** of the corpus, extended by the families in §3 |

### 1.2 Determinism, stated per layer

A reproducible retrieval decision does not imply reproducible model interpretation. Reproducibility is therefore defined separately:

| Layer | Reproducibility requirement |
|---|---|
| Application decisions (scope, budget, routing, deletion, ID assignment, assembly order) | **Deterministic and replayable**; byte-identical given the same inputs |
| Retrieval ranking (deterministic parts) | Reproducible given fixed index and query |
| Model outputs (extraction, prose, adjudication, replies) | **Sampled**; reported with variance, never as a single number |

A run that cannot replay its application decisions is invalid. A run that shows model variance is normal and must be reported, not smoothed.

---

## 2. Phase 1 — Rules that cannot change

These are locked before any run and appear in the run manifest. Changing one invalidates every run collected under the old rule.

| Rule | Statement | Why it is locked |
|---|---|---|
| `R-1` | Competing configurations use the **same source conversations** | Otherwise the comparison measures corpus, not configuration |
| `R-2` | **Same model file and version** within a comparison; hash pinned | Model substitution changes the reader |
| `R-3` | **Same relevant system prompts and generation settings** within a comparison | Format and sampler are known confounds (`MEMORY-06`) |
| `R-4` | Context, reasoning, and answer budgets recorded **separately** | Conflating them invalidates the thinking-budget comparison |
| `R-5` | Fixed sampling settings; seeds recorded; deterministic mode where the backend permits | Enables replay and variance reporting |
| `R-6` | Test cases, scoring criteria, and failure classes **pre-registered** | Prevents post-hoc success definitions |
| `R-7` | **No tuning a configuration against the evaluation results and then reporting those results as validation** | The single most common way an evaluation becomes fiction |
| `R-8` | Every arm reported with the **transcript-only baseline** run in the same session | Without it, a memory "win" may be a baseline loss |
| `R-9` | Every failed item assigned to exactly one **failure class** | Aggregate scores cannot be acted on |
| `R-10` | Application-side decisions must be **replayable from the manifest** | Distinguishes a bug from model variance |

**A rule change is a dated amendment that starts a new collection.** It is never a note added mid-analysis.

### 2.1 Safeguards locked at approval

| Rule | Requirement |
|---|---|
| `R-11` | **Corpus contamination.** Evaluation material is isolated from configuration tuning. The corpus is split `tuning` / `heldout`; the held-out set is hash-sealed before tuning begins and is not read by whoever tunes. A result measured on material that shaped the configuration is not a result. |
| `R-12` | **Reviewer independence.** Creative reviewers are blind to the arm that produced a conversation. |
| `R-13` | **Repeated runs, two levels.** Variance is reported across seeds *within* a scenario and across scenarios. Scenario-level variance is larger and is never pooled away. |
| `R-14` | **Negative outcomes preserved.** Every arm is retained and reported, including configurations that perform worse than transcript-only. A deleted losing arm makes later comparison unverifiable. |

**Mechanism for `R-11`:** every item carries `split: tuning | heldout`. Minimum split for the final corpus: **≥40% held out**, and **≥2 held-out items per family**. Moving an item from held-out to tuning after tuning began is a contamination event and is reported as one.

---

## 3. Phase 2 — Corpus specification

### 3.1 Principles

1. **Synthetic by default.** No user chat, character card, world, or lorebook is read without separate written authorization (`MEMORY-03` §5.6).
2. **Different kinds of conversation, not many variations of one.** Coverage beats volume.
3. **Every item has a reference answer or an explicitly documented acceptable range.** An item without one is not a test item.
4. **Originals are preserved as evidence.** The corpus stores the source conversation untouched, alongside labels.
5. **Include the cases where the right answer is "no memory".** Abstention is a capability, not a gap.

### 3.2 Test families

| # | Family | What it challenges | Minimum items |
|---:|---|---|---:|
| `F-1` | Long-running relationships | Emotional and interpersonal continuity | 10 |
| `F-2` | Multi-character conversations | Individual knowledge boundaries | 10 |
| `F-3` | Long conversations | Important details separated by hundreds of messages | 10 |
| `F-4` | Changing circumstances | Updates, supersession, temporal validity | 10 |
| `F-5` | World and scenario changes | Context isolation and reusable knowledge | 8 |
| `F-6` | Casual conversations | One unexpectedly important detail among noise | 8 |
| `F-7` | Contradictory information | Uncertainty and conflict handling | 8 |
| `F-8` | Irrelevant information | Resisting unnecessary formation | 8 |
| `F-9` | No-retrieval cases | Knowing when to leave memory out | 8 |

**Total minimum: 80 scored items**, plus the twelve `MEMORY-04` fixtures which are retained as required members (`E-01`–`E-12`), not replaced.

**This floor is not proof of adequacy.** Whether the final corpus is sufficient is decided by the pilot and the CI procedure in §10, not by reaching a count.

### 3.5 Principle tests — required members

The corpus must test the **principle**, not the proposed architecture. Four case types are mandatory, because an evaluation built only from "did it recall X" items rewards the exact behaviour the governing principle forbids.

| Case | What it proves | Content requirement |
|---|---|---|
| **Less is better** | Over-formation has a real cost | A correct-but-verbose presentation produces a worse reply than a sparse one on the same knowledge |
| **Nothing is correct** | Abstention is a capability | Items where the right action is silence: insufficient, unauthorized, or irrelevant |
| **Accurate but contextually wrong** | Link-1 correctness does not excuse links 2–4 | A true memory that must not reach this character, persona, time, or moment |
| **Natural use without recitation** | Guidance over instruction | A character behaves as though they remember, without stating the memory |

**Restraint is not silence (owner, 2026-11-12).** Six behaviours are scored separately and never aggregated: suppress mention, influence tone, natural use, recognise irrelevance, respect entitlement, preserve uncertainty. A configuration that suppresses everything passes four of the six and fails the two that require the memory to be perceptible (`R-2`, `R-3`); that result is reported as **avoidance**, not as restraint. The six-mode report is mandatory output for every arm.

**Provenance rule (owner, 2026-11-12).** Corroboration requires **known independence**. Three records that descend from one origin are one source; three records whose relationship is unknown are **not** three confirmations. An incomplete provenance record is preserved as incomplete and never completed by assumption (`R6-03a`–`R6-03d`).

### 3.6 Split discipline

### 3.3 Per-item record

```text
item_id
family
source_conversation_id
turn_range
reference_answer | acceptable_range[]
required_evidence_ids[]
forbidden_evidence_ids[]        # leakage set
expected_abstention: true|false
expected_character_perspective
as_of_time (optional)
expected_failure_classes_if_wrong[]   # pre-registered diagnostics
notes
```

### 3.4 Reference-answer discipline

| Type | When | Requirement |
|---|---|---|
| Exact answer | Factual, name, number, event | One canonical string plus allowed variants |
| Acceptable range | Open creative response | Explicit criteria (must mention X, must not claim Y) |
| Expected abstention | Insufficient, unauthorized, or irrelevant | The reason must be stated for each item |
| Expected correction | Superseded information | Must state the current truth *and* that it changed |

An item whose "correct" answer cannot be stated before the run is removed from the corpus, not judged afterwards.

**Sufficiency rule (owner, at G1 batch-1 acceptance):** every creative reference answer separates four parts — `mandatory` (scored), `variation` (equally correct, all pass), `never scored` (length, tone, dialect, humour, confidence — never a failure signal), and `prohibited` (scored). A creative item that scores a stylistic difference as a factual failure is a **defective item** and is reworked rather than reinterpreted.

**Gate rule `G-1`:** a blocking-gate failure (scope leakage, deletion failure, fabricated provenance) is reported before any score and is not offset by performance elsewhere. An impressive aggregate never overrides a gate.

**Minimality rule `M-1` (owner review, 2026-11-12):** `mandatory` is a floor, never a checklist to be recited, and every item states its **over-exposure boundary** — the additional true material that would fail it. Without this, a verbose configuration is never worse than a concise one on required facts, and the corpus rewards the behaviour it exists to catch. For restraint items the boundary is **scored**.

---

## 4. Phase 3 — Configurations and variable control

### 4.1 Arm A — Formation

| Arm | Configuration | Purpose |
|---|---|---|
| `A0` | Transcript only, no derived memory | Baseline (mandatory) |
| `A1` | Deterministic extraction only | Availability without model dependence |
| `A2` | Chunked AI-assisted formation | The proposed direction |
| `A3` | Chunked formation + cross-chunk reconciliation | Tests the value of stages 3–4 |
| `A4` | Chunked formation + reconciliation + validation | Full pipeline |

### 4.2 Arm B — Retrieval

| Arm | Configuration | Purpose |
|---|---|---|
| `B0` | Lexical only | Mandatory baseline (`D-RET-1`) |
| `B1` | Dense only | Semantic coverage |
| `B2` | Hybrid + fusion | Coverage |
| `B3` | Hybrid + rerank | Precision |
| `B4` | Hybrid + rerank + state-conditioned filters | Full selection |

### 4.3 Arm C — Presentation (critical acceptance)

| Arm | Configuration | Constraint |
|---|---|---|
| `C1` | Verbose memory block | **Constant token budget** |
| `C2` | Compact guided block | **Constant token budget** |
| `C3` | Minimal cue set | **Constant token budget** |
| `C0` | No memory (control) | Separate control, not one of the three |

**Constant budget is the whole experiment.** If the verbose arm is permitted more tokens, the test measures budget, not presentation.

### 4.4 Arm D — Reasoning budget

| Arm | Configuration |
|---|---|
| `D0` | Reasoning off |
| `D1` | Small budget |
| `D2` | Proposed 2048-token budget, where the backend permits |
| `D3` | A non-reasoning model of comparable size |

**Record three quantities independently for every arm:** reasoning tokens, supplied-context tokens, reserved answer tokens.

### 4.5 Variable-control matrix

| Held constant | Varies by design | Recorded always |
|---|---|---|
| Source conversations | The arm | Model file hash, template hash, quant |
| Model file within a comparison | Presentation form (`C`) | Context / reasoning / answer budgets |
| System prompt and template | Retrieval strategy (`B`) | Sampler chain and seed |
| Token budget within `C` | Formation method (`A`) | Hardware, offload, free VRAM |
| Scoring rubric and raters | Reasoning budget (`D`) | Corpus and label hashes |

---

## 5. Scoring rubrics

Six dimensions. Each is scored per item on a 0–4 scale with anchors, then aggregated with the failure-class distribution reported alongside.

| Dimension | 0 | 2 | 4 |
|---|---|---|---|
| **Formation quality** | Invents or misstates | Correct but noisy or bloated | Accurate, meaningful, durable, no invention |
| **Retrieval precision** | Wrong or irrelevant material | Mostly relevant, some noise | Relevant, minimal, correctly scoped |
| **Contextual correctness** | Wrong character/persona/world/time | Partly correct scope | Fully correct scope and as-of time |
| **Presentation quality** | Recitation or over-exposition | Serviceable but intrusive | Natural guidance, no recitation |
| **Creative continuity** | Breaks voice or relationship | Adequate | Voice, relationship, and history all coherent |
| **Abstention** | Injects what it should withhold | Inconsistent | Correctly withholds irrelevant, uncertain, or unauthorized material |

**Aggregation rule:** dimensions are reported separately. No weighted total is published as a headline, because a total invites trading correctness for volume.

**Raters:** two per item, blinded to arm; disagreement adjudicated; κ reported; automated judge only as a secondary, calibrated signal (`MEMORY-03` `D-HUMAN-1`).

---

## 6. Execution allocation

Which tests need a model, and which can be settled before any model is loaded. This is what makes the programme affordable.

| Test | Needs a model? | Why | Can run before `G3`? |
|---|---|---|---|
| Corpus construction, labels, reference answers | **No** | Authoring work | **Yes** (needs `G1`) |
| Scope enforcement (structural) | **No** | Store-level constraint | **Yes** (`G1`+`G2`) |
| Deletion propagation and receipts | **No** | Store-level mechanics | **Yes** |
| Provenance completeness | **No** | Link integrity | **Yes** |
| Token accounting by section | **No** | Tokenizer only | **Yes** |
| Index build / read / update / delete latency | **No** | Store benchmark | **Yes** |
| Cache invalidation correctness | **No** | Application state | **Yes** |
| Deterministic extraction | **No** | Rules only | **Yes** |
| Retrieval recall / nDCG / precision | **No** | Ranking over labelled evidence | **Yes** |
| Schema validity and template fidelity of the assembly | **No** | Renderer and validator | **Yes** (`G2`) |
| **Formation significance and prose** | **Yes** | Model interpretation | **No** |
| **Conflict adjudication** | **Yes** | Model judgment | **No** |
| **Presentation effect** | **Yes** | Reply behavior | **No** |
| **Creative continuity** | **Yes** | Reply behavior | **No** |
| **Abstention judgment** | **Partly** | Application decides; model may still inject | **Partly** |
| **"Misunderstood" failure class** | **Yes** | Requires a reply | **No** |

**Consequence:** roughly half the programme — scope, deletion, provenance, budget, indexing, retrieval ranking, schema and template validity — can be completed **before** a single model is loaded. A large share of the risk is retired without spending inference.

---

## 7. Gates

| Gate | Required result | Blocking? |
|---|---|---|
| **Integrity** | No fabricated provenance, identity leakage, or silent deletion failure | **Absolute** |
| **Formation** | Meaningful information retained without uncontrolled accumulation | Absolute on accumulation |
| **Retrieval** | Strong relevance and appropriate abstention | Absolute on leakage and abstention-into-injection |
| **Presentation** | No systematic degradation vs the no-memory baseline | Absolute on degradation |
| **Creative acceptance** | Demonstrable continuity improvement without loss of naturalness | Absolute |
| **Operations** | Practical performance on supported local hardware | Conditional on a stated target |

**No trading rule:** a configuration that improves recall while leaking knowledge fails. A configuration that improves retrieval while making characters robotic fails. An aggregate score never overrides a gate.

### 7.1 Failure classes (required output)

Per `ARCHITECTURE-DECISION-MEMORY-SYSTEM` §6.9: never formed, incorrectly formed, mis-scoped, stored incorrectly, not retrieved, presented poorly, misunderstood, **correctly withheld**. The last is success and must never be counted with *not retrieved*.

---

## 8. Evidence required before freezing anything

Each unfrozen architectural choice names the evidence that would justify committing to it. No choice is frozen because it looks clean.

| Unfrozen choice | Evidence required |
|---|---|
| **Record schema** | Every class exercised in the corpus; zero duplicate-truth violations across all fixtures; 100% provenance on admitted records; deletion propagation verified per class; `Q-MF-10` answered |
| **Retrieval algorithm** | `D-RET-1`/`D-RET-2` met against lexical; abstention gate passed; no scope failure |
| **Chunk size / overlap** | Family `F-3` result plus a stated cost curve; no cross-chunk reconciliation loss |
| **Reasoning budget** | Arm `D` shows a decision-quality gain that justifies its processing cost, or it is rejected |
| **Presentation form** | Arm `C` shows a difference (or shows none, which is itself a finding) with diagnostics |
| **Store adapter** | Latency, growth, and delete cost measured on `HW-1`; `DA-M1`/`DA-R1` resolved |
| **Model selection** | Per-tier results with no tier extrapolation; `D-SMALL-1` |
| **Formation method** | `A2`–`A4` compared; `Q-MF-01` closed |

**A choice may remain unfrozen indefinitely.** Nothing forces a decision that the evidence does not support.

---

## 9. Canonical sequence

| Phase | Content | Exit condition |
|---|---|---|
| **P0 — Rules lock** | §2 locked; manifests defined | Rules recorded; no run yet |
| **P1 — Corpus** | §3 built; labels and reference answers | Corpus hash frozen |
| **P2 — Application-side tests** | §6 rows marked "Yes"; no model | Deterministic results recorded |
| **P3 — Model comparisons** | Arms A–D | Comparative results with variance |
| **P4 — Scorecard and Stage H** | §5 rubrics, §7 gates, blinded creative acceptance | Gate decisions recorded |
| **P5 — Analysis and replication** | Effect sizes, CIs, replication on a second model file or `HW-3` | Report; no silent averaging |

`MEMORY-03`'s phases map onto this: its Phase 0 = P0, Phase 1 (pilot) = inside P3, Phases 2–4 = P3–P5, Phase 5 = P5.

**Stop rules** are `MEMORY-03` §9.8, unchanged: real user data, a safety-gate failure, model/template drift, >5% crashes, store corruption, reserve breach, or budget exhaustion halts the affected condition.

---

## 10. Sample size and thresholds

**No sample size or threshold is invented here.** Both are set by procedure, after the pilot, before final validation.

| Quantity | Procedure |
|---|---|
| Pilot | 2 families, 16 items, two tiers; purpose is variance estimation and fixture validation only |
| Sample size | Chosen from pilot variance to reach a target confidence-interval width on the primary outcome; practical caps (time, hardware) recorded |
| Thresholds | Set from the pilot baseline and the smallest effect the product would act on; frozen before final collection |
| Adjustments | Any adjustment after final collection invalidates the run; a dated amendment starts a new collection |
| Reporting | Effect sizes with 95% CIs, raw counts, per-class failure distribution, and variance across seeds |

**If the pilot cannot estimate variance usefully, the honest output is "sample size unresolved", not a fabricated number.**

---

## 11. What this document does not authorize

- No model load, no corpus build, no run (`G1`–`G5` closed).
- No schema, source, or runtime change.
- No architecture selection; §8 lists what evidence *would* justify one, and none is met.
- No threshold, sample size, or scoring weight invented to appear rigorous.
- No modification of a frozen artifact or an owner decision.

**Current state:** this package is **approved** as the governing pre-registration authority, and **`G1` is authorized** for corpus and labels only. `G2`–`G5` remain closed, and no model may be loaded until a separate authorization.

**Not established by this document:** that the memory architecture works, that any corpus is adequate, or that any architectural choice is justified. Nothing has been measured.
