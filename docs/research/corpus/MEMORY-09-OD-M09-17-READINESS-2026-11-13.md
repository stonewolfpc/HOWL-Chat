# `MEMORY-09` — `OD-M09-17` Reproducibility Readiness (Stage Three)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-REPR-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **READINESS ANALYSIS — study not executed, no threshold adopted.**
**Authority:** `MEMORY-09-WO-ARCH-01` §5.
**Relationship to prior work:** conceptually updates `MEMORY-09-PROTO-1` (row 445) without overwriting it. That report stands as written.

---

## 1. The eight elements, verbatim from the governing record

`MEMORY-09-OWNER-DECISION-ADDENDUM-1-OD-M09-15-18-2026-11-13.md` **L104–111** — read from the record, not from the prior report's summary:

| # | Requirement |
|---|---|
| **E-1** | "Rater qualification and identity/provenance requirements." |
| **E-2** | "Independent assessment procedures." |
| **E-3** | "Blinding and the information available to each rater." |
| **E-4** | "The exact judgement being compared." |
| **E-5** | "How disagreements are recorded and handled." |
| **E-6** | "The reliability measure and its justification." |
| **E-7** | "Any acceptance threshold, specified **before** results are observed." |
| **E-8** | "How insufficient sample size, missing judgements and unresolved disagreements are reported." |

Governing constraints from the same record: L90 (results "may not support claims beyond" single-source confidence until a two-rater study exists), L97–L100 (repeated judgements by one source are **not** independent; different sessions or conflicting labels do **not** establish lineage; independence "must be established and documented"), L116 (the record "invents **no** reliability threshold and **no** statistic"), L118 (bounded **design** work is not prohibited).

## 2. Dependency classification

| # | Classification | Why |
|---|---|---|
| **E-1** | **Finalizable with an explicit assumption** | The record says "rater", never "model rater" (L104). Human raters are therefore **not excluded** — but independence must be *documented*, not assumed (L100). Assumption to be recorded: raters may be qualified humans or models, with provenance established per §5.1. **Owner should confirm** the human-rater reading |
| **E-2** | **Ready to finalize** | Procedures are writable now: separate submissions, no cross-access, hashed on receipt |
| **E-3** | **Ready to finalize** | Packet contents are decidable now (item, scene, reply, mode definitions, scale; **not** expected verdicts, prior reviews, or aggregates) |
| **E-4** | **Blocked by a measurement decision** | Whether the compared judgement is **compliance** alone or **compliance + cue-dependence** depends on Stage Two. This is the single clearest cross-stage dependency |
| **E-5** | **Ready to finalize** | Raw matrix preserved; adjudication separate and never substituted |
| **E-6** | **Finalizable with an explicit assumption** | Measure justified below; one sub-decision (below) needs owner approval |
| **E-7** | **Requires owner decision** | Thresholds remain unapproved proposals. The record invents none (L116), and none may be chosen after results |
| **E-8** | **Finalizable with an explicit assumption** | Reporting rules writable; the 20-unit floor remains an unapproved proposal |

**Nothing is classified "blocked by unavailable resources" except execution itself**, which this order does not authorise anyway. That is the practical answer to the single-model constraint: **the constraint blocks the study, not the preparation.**

## 3. The single-model constraint, handled honestly

**Prohibited and not done:** running the same model twice and calling it two raters; prompting it with a different persona or "reviewer" framing to manufacture a second source; describing separate passes as agreement. L97–L99 denies all three, and `WO-ARCH-01` §2.1 forbids the claim.

**What the constraint does *not* block:**

| Work | Feasible with one model | Why |
|---|---|---|
| Rating instructions and operational definitions | Yes | Documentation |
| Source packets and blind presentation rules | Yes | Deterministic file assembly |
| Disagreement and missing-data forms | Yes | Templates |
| Data schema for the rating matrix | Yes | Design |
| Reproducibility calculation, validated | **Yes** | Deterministic arithmetic; validated against known cases (§4) |
| Audit and provenance procedure | Yes | Design |
| The rating itself, twice, independently | **No** | Requires a genuinely separate source |

**Human raters as the practical route.** The record's language does not restrict raters to models. A pair of qualified humans, with documented and non-overlapping provenance, would satisfy E-1's letter — subject to (a) owner confirmation that human raters are admissible, and (b) the corpus being ratable by a human reading text, which the Stage Two architecture actually improves, because it asks for observable compliance rather than inferred cognition.

## 4. Statistical design

### 4.1 Measure

**Krippendorff's α (nominal)** remains the recommendation, on grounds that are structural rather than preferential:

- The scale is nominal (four unordered labels).
- Missing cells are expected (E-8), and α handles them without imputation.
- α generalises beyond two raters, so a third adjudicator later does not force a measure change.
- Cohen's κ assumes two raters and no missing data, and is prevalence-sensitive; percent agreement ignores chance.

**Reported alongside:** observed agreement, category distribution, and the unit count — because α alone hides the two failure modes that matter here (see 4.3).

### 4.2 Implementation must be validated before use — evidence from this session

While preparing this report, a hand-written α implementation was checked against a six-unit worked example (five usable units, one missing cell; categories `D`, `ND`, `NA`). **The first three attempts were wrong, in three different ways:**

| Attempt | Defect | Yielded |
|---|---|---|
| 1 | Counted all observed pairs as agreement | α = **1.00** (should be 0.45) |
| 2 | Included a unit that α requires dropping (m < 2) | marginals 4/4/3 instead of 4/3/3 |
| 3 | Inverted the agreement/disagreement convention — used `(Do−De)/(1−De)` instead of `1 − Do/De` | α = **−1.25** (impossible: α ≥ −1 only at maximal disagreement) |

Correct value, confirmed by two independent computation paths (pair enumeration and coincidence matrix): **Do = 2/5, De = 11/15, α = 5/11 ≈ 0.455**. Sanity checks then passed: perfect agreement → α = 1.000; random raters (n = 40) → α = −0.060.

**This is the argument for a validation requirement, not an anecdote about carelessness.** Three of four implementations produced plausible-looking numbers. A protocol that reports α without validating its implementation will report a wrong number with full confidence.

**Rule proposed:** before any real α is computed, the implementation must reproduce (a) α = 1 on perfect agreement, (b) α ≈ 0 on independent random ratings, and (c) a published worked example. Results of all three recorded.

### 4.3 The one open statistical sub-decision

`not assessable` is a **judgement**, not a missing value. Treating it as a category inflates expected disagreement (it is common and clustered); treating it as missing discards real information and can make the sample look cleaner than it is.

| Option | Effect | Recommendation |
|---|---|---|
| `not assessable` as a **category** | De inflated; α conservative | **Primary** — it is a judgement the rater made |
| `not assessable` as **missing** | De lower; α can rise misleadingly | **Sensitivity only**, reported beside |

Reporting both is cheap and makes the choice visible. **Requires owner approval** because it changes what the number means.

### 4.4 Thresholds

**Unchanged and unapproved.** α ≥ 0.667 / ≥ 0.800 and the 20-unit floor remain proposals from row 445. They are not adopted here, not used to interpret anything, and must be fixed before results if adopted at all (E-7). No result exists, so nothing can have been chosen after seeing one.

## 5. What is ready, blocked, and prohibited

| Status | Items |
|---|---|
| **Ready to finalize now** | E-2, E-3, E-5; instructions, packets, forms, schema, audit procedure, α validation |
| **Finalizable with a recorded assumption** | E-1 (human raters admissible), E-6 (α nominal; NA-as-category), E-8 (20-unit floor as proposal) |
| **Blocked by a construct decision** | E-4 in part — which modes are ratable at all depends on Stage One's `R-1` outcome |
| **Blocked by a measurement decision** | E-4 in part — compliance vs compliance + cue-dependence |
| **Requires owner decision** | E-7 (threshold), NA-as-category, human-rater admissibility |
| **Prohibited until separately authorised** | Any rating session, any model execution, any rater recruitment, any published α |

## 6. Adversarial passes (same model, sequential — not independent)

| Pass | Objection | Result |
|---|---|---|
| **Pass 1** | Eight elements can be finalized on paper now | — |
| **Pass 2** | Objection: finalizing protocols for a study that may never run is documentation for its own sake. Counter: §2 shows only execution is blocked; the artifacts are the precondition for any future run, and E-4's block is *informative* — it says the corpus is not yet ratable in a defined sense. Second objection: the human-rater route is an escape hatch that quietly weakens "two-rater" into "two people", who may share every bias a single model has. | **First survives weakly** — it lowers the value of finalizing E-2/E-3/E-5 but does not make them wrong. **Second survives and is material**: human raters with shared training are not automatically independent; L100's "documented" requirement must do real work, and §5.1's denial of inferred independence applies to humans too |
| **Pass 3** | Classification stands; confidence **high** for the dependency map, **moderate** for the measurement choices | What changed: the human-rater route is now stated as *conditional* on documented non-overlapping provenance, not as a solution. What did not: α remains justified over κ on missing-data and multi-rater grounds |

**Claim ceiling:** a complete dependency map and validated-measurement procedure. **Not** evidence about reproducibility, and not a claim that the corpus can be rated by anyone.
