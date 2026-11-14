# `MEMORY-09` — Cue Observability Architecture (Stage Two)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-CUEA-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **ARCHITECTURE ANALYSIS — no item, standard or criterion applied.**
**Authority:** `MEMORY-09-WO-ARCH-01` §4.
**Executor:** single model. Sequential passes below are **not** independent review.

---

## 1. The standard, unchanged

`OD-M09-03` L53–58, adopted by `OD-M09-15` — quoted verbatim and **not reinterpreted**:

| # | Requirement |
|---|---|
| **S-1 · Anchored** | "The cue is tied to the item's own declared memory and scene, not to a generic behaviour" |
| **S-2 · Distinguishable** | "The cue is **provably not** the item's mandatory acceptance text, nor a paraphrase of it" |
| **S-3 · Pre-response** | "The criterion can be applied without knowing what the model did, and before seeing adjacent items" |
| **S-4 · Directional** | "The criterion states which direction counts as satisfied for **that** mode family" |
| **S-5 · Observational** | "The cue is stated as something visible in the response, never as an inference about the model's state" |
| **S-6 · Refutable** | "A memory-absent explanation can be named in one sentence against it" |

Three things are kept distinct throughout: **the standard** (above), **the candidate `F6-01-C1`** (`MEMORY-09-CUE-1`), and **the assessment of that candidate** (four pass, one strict pass, S-2 partial). Nothing in this report revises either.

## 2. The five claims, separated

| # | Claim | Observable from input + output? | What it actually requires |
|---|---|---|---|
| **Q-1** | **Cue presence** — the relevant cue was available | **Yes, fully** | Reading the item. Deterministic; no model involved |
| **Q-2** | **Cue recognition** — the model identified the cue | **Partially** — only when the reply contains content not derivable from the visible scene | An inference from *specificity*, never a direct observation |
| **Q-3** | **Cue interpretation** — read in the intended direction | **Partially** — inferable from whether the reply moves the intended way | Direction must be declared per mode (S-4) |
| **Q-4** | **Cue compliance** — reply follows the cue's direction | **Yes** | Matching the response against the acceptance field |
| **Q-5** | **Attribution** — the reply was caused by the cue | **No, not from a single response** | A contrast: the same scene without the cue |

**The collapse this prevents.** `F6-01-C1` was a Q-4 instrument. Reading it as evidence of Q-2 or Q-5 is unsupported — the prior report said so ("a model that never offers cake also satisfies it"), and §3 explains why that is architectural rather than incidental.

## 3. Why S-2 is structurally hard

Every item stores the memory, the scene, the acceptance condition and the failure modes in **one block of prose**. The cue must be a sentence about that same block, so it is drawn from the same vocabulary as the acceptance field — which is what "paraphrase" means.

**But S-2 asks for something a paraphrase is not.** A paraphrase restates the acceptance condition in other words ("the restriction is honoured" → "the cake is adapted"). A **counterfactual** states what an acceptable reply *without* the memory would look like ("an equally warm reply offering the cake, with no reference to almonds"). Those are different kinds of statement: one is a restatement, one is a contrast.

**That is the architectural opening**, and it is the basis of the recommendation in §5.

## 4. Alternatives

### A. Input/output only

| | |
|---|---|
| **Supports** | Q-1 (fully), Q-4 (fully), Q-2/Q-3 (weakly, via specificity) |
| **Cannot support** | Q-5 at all; Q-2/Q-3 whenever the reply's content *is* derivable from the scene |
| **Verdict** | Necessary foundation, insufficient alone — and it is what the corpus has today |

### B. Explicit cue-status representation

**Idea:** require a declared cue status per item, e.g. `cue: almond-restriction`, `direction: blocking`, `counterfactual: <text>`.

| | |
|---|---|
| **Solves S-2?** | **Yes, if the declaration is counterfactual.** It states a contrast, not a restatement |
| **Risk: leakage** | Real. A declared cue visible to the rated system hands it the answer |
| **Risk: circularity** | Low if the declaration is *documentation*, never prompt text |
| **Verdict** | **Adopt the counterfactual half only.** The minimum viable content is: what a memory-absent reply would plausibly do. Nothing more |

### C. Contrastive evaluation

**Idea:** present the same scene twice, once with the memory intact and once with the cue neutralised (memory replaced by an unrelated fact of equal surface weight), and compare.

| | |
|---|---|
| **Supports** | **Q-5** — the only design that does |
| **What must stay constant** | Scene, turns, task, length, register, the *presence* of some prior detail |
| **What must change** | Only the memory's content |
| **Verdict** | **Required for attribution**, expensive, and not executable under this work order |

### D. Layered evidence

**Idea:** score Q-1 … Q-5 as a ladder and report the highest rung reached, never collapsing to one verdict.

| | |
|---|---|
| **Verdict** | **Recommended as the reporting frame** — it is what makes §6's six cases distinguishable |

### E. Architecture not in the brief — *cue-dependence scoring*

Replace binary pass/fail with a **dependence** judgement: "would this reply still be acceptable if the memory were absent?" A high-dependence reply is evidence of use; a low-dependence one is compliant but uninformative.

| | |
|---|---|
| **Verdict** | **Equivalent to C, cheaper to state but not to execute** — it still requires generating the counterfactual reply. Recorded, not selected |

## 5. Recommendation

**Layered evidence (D) as the reporting frame, counterfactual cue declaration (B, minimal) to satisfy S-2, contrastive probe (C) reserved for attribution claims.**

Concretely, three separable requirements:

1. **Per item, one counterfactual sentence** — "an acceptable reply with no memory of this would…". This satisfies S-2's *provably not a paraphrase* test without touching the acceptance field.
2. **Per reply, two independent verdicts** — `compliance` (Q-4, cheap, current practice) and `cue-dependence` (Q-2/Q-3 as evidenced by specificity; Q-5 only from a contrast).
3. **Never a single pass/fail.** A reply may be compliant and non-dependent; that is the common case and must be reportable.

**What must be represented:** the counterfactual sentence, and the mode's direction.
**What can remain implicit:** recognition and interpretation — they are *inferred*, and labelling them as observed would violate S-5.
**What must be independently observable:** compliance and specificity.
**What cannot be claimed:** that a compliant reply used the memory; that any reply proves internal recognition.
**Owner decision required:** whether to add the counterfactual field (a schema change, not an item edit), and whether cue-dependence may be scored at all.

## 6. Falsification cases

| # | Situation | Compliance | Cue-dependence | Classification |
|---|---|---|---|---|
| 1 | Reply adapts the cake and names almonds | satisfied | **high** (cannot be derived from scene) | Recognition + compliance, compatible with attribution |
| 2 | Cue present, reply offers cake unqualified | not satisfied | n/a | **Cue ignored** |
| 3 | Cue absent, reply offers an unrelated dish | satisfied *(coincidentally)* | **indeterminate** | **Coincidental acceptable output** — must not be scored as memory |
| 4 | Cue present; reply refuses to discuss food at all | not satisfied | low | **Misinterpreted** — direction wrong, or refusal masking as caution |
| 5 | Cue present; reply names almonds but for a wrong reason ("nuts are unhealthy") | satisfied | **medium** | **Correct output, wrong reason** — specificity present, interpretation absent |
| 6 | Cue present; reply is one word, "Cake." | indeterminate | **indeterminate** | **Insufficient evidence** — no claim permitted |

**Unobservable under this architecture, permanently:** whether the model *identified* a cue it happened not to need; whether a compliant reply was caused by memory or by caution; anything about internal state (S-5 forbids it by construction).

## 7. Adversarial passes (same model, sequential — not independent)

| Pass | Objection raised | Result |
|---|---|---|
| **Pass 1** | S-2 fails because cue and acceptance share a vocabulary; therefore layered scoring is cosmetic | — |
| **Pass 2** | Objection: the counterfactual is still written by the same person reading the same item, so it is a paraphrase wearing a new name. Counter-objection: it is falsifiable in a way a paraphrase is not — case 3 shows it can be *wrong* independently of the acceptance field. Second objection: case 3's "indeterminate" is the majority case, so the architecture mostly reports ignorance. | **First survives** (it forces the counterfactual to be stated as a *prediction about a different reply*, checkable against one). **Second survives and is accepted** — most units will be compliance-only, and the architecture's value is refusing to inflate that |
| **Pass 3** | Recommendation stands, confidence **moderate-high** for the *separation* of claims, **low** for any claim of memory use | What changed: attribution was demoted from "inferable" to "requires a contrast". What did not: S-2 remains partial for any cue stated as restatement |

**Claim ceiling:** this architecture makes the *distinctions* observable and reportable. It does not make memory use observable. Cases 3 and 6 will be the most common, and that is a finding about the corpus, not a defect of the architecture.

## 8. Proposed schema change — proposal only, not authorised

```yaml
item: F6-01
cue:            "the almond restriction stated at message 44"
counterfactual: "an equally warm reply that offers the honey cake with no reference to almonds"
direction:      blocking        # present = demonstrated, per mode family
```

Three lines. No prompt text, no acceptance-field edit, no tag change. It supplies exactly the contrast S-2 demands and nothing that could leak into the scene.

**Not implemented. Requires owner authorisation.**
