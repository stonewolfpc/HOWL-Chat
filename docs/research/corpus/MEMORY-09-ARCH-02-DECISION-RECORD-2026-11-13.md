# `MEMORY-09` — `WO-ARCH-02` Decision Record and Synthesis

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-DEC-2` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **RECOMMENDATIONS ONLY — nothing below is owner-approved.**
**Authority:** `MEMORY-09-WO-ARCH-02` §6–§7.
**Executor:** single model. Passes are sequential and self-challenging; **no independence credit**.

---

## 1. Decision table

| Decision | Recommendation | Evidence | Remaining uncertainty |
|---|---|---|---|
| **Counterfactual field for S-2** | **Pilot, additive-cue items only** | Three drafts (row 451): `F6-01` distinct, `P-01` partially distinct, `R6-01` a restatement | Falsifiability holds only where the memory supplies content absent from the scene |
| **`R-1` architecture** | **Tact (A)**, now operationalized | Scenario `R1-PROP-01` (row 452) with four required elements, `R-5` excluded, `R-4` narrowed not closed | The relevance-domain question ("subject" vs "moment") |
| **Cue-dependence** | **Separate verdict, with `indeterminate` permitted** | Attribution unobservable from one reply (row 448) | Attribution needs a contrast that does not exist yet |
| **Human-rater admissibility** | **Admissible under clarification (A)** | Record never names a rater type; L100 requires documented independence (row 453) | Shared-bias control is unspecified in the authority |
| **Statistical policy** | **Defer thresholds; adopt the validation gate** | α incident (row 449) + 400-sample simulations (§3) | Thresholds, NA handling, sample floor all unapproved |

**Unchanged and unapproved:** α ≥ 0.667, α ≥ 0.800, the 20-unit floor, NA-as-category, any study execution.

## 2. What each stage actually established

| Stage | Established | Not established |
|---|---|---|
| One | A counterfactual is **independent** for additive cues and **collapses** for quantitative and relational ones (1 of 3 passes) | That a pilot would improve reliability |
| Two | Tact can be operationalized; `R-5` is excludable, `R-4` is narrowed | Construct validity from one scenario |
| Three | The record is **silent** on rater type; every requirement is satisfiable by humans under documentation | That humans are *permitted* |
| Four | — | — |

**The one-of-three result is the most useful finding in this order.** A universal counterfactual field would have looked successful while doing nothing on two of three item classes.

## 3. Statistical validation gate (carried forward)

The `WO-ARCH-01` incident — three successive wrong α implementations yielding 1.00, a marginal error, then −1.25 before two paths agreed on 5/11 — motivates a **mandatory gate before any real α is reported**:

| # | Gate | Pass criterion | Status of the implementation used here |
|---|---|---|---|
| 1 | **Perfect agreement** | α = 1 exactly | **Pass** — 45 units, α = 1 |
| 2 | **Random-rating simulations**, expected behaviour across repeated samples | Mean ≈ 0; individual samples need **not** be 0 | **Pass** — 400 samples per condition: mean α = 0.0021 (0% missing), −0.0054 (10%), 0.0054 (25%); 90% intervals −0.129…0.138, −0.141…0.142, −0.157…0.178 |
| 3 | **Published worked example** | Exact reproduction | **Not performed** — no published example was available offline; this gate is **open** |
| 4 | **Two independent computation paths agree** | Identical α, Do, De | **Pass** — pair enumeration and coincidence matrix both give Do = 2/5, De = 11/15, α = 5/11 |
| 5 | **Retained artifacts** | Coincidence matrix, marginals, Do, De, α, unit counts all inspectable | **Required by protocol**; retained per run when a study is authorised |

**On gate 2, explicitly:** a single random sample withholding α ≈ 0 is *not* a failure. The correct expectation is distributional — individual samples scatter (90% interval ≈ ±0.14 at n = 45), and only the mean should sit near zero. My earlier single-sample check (α = −0.060) was **insufficient evidence** on its own; the 400-sample runs above replace it.

**Gate 3 remains open and blocks any real α report.** No α was computed on corpus ratings here, and none may be.

## 4. Dependency map

| Can proceed now (one model) | Requires owner decision | Requires independent raters | Requires contrastive evaluation / instrumentation | Deferred |
|---|---|---|---|---|
| Counterfactual drafts for additive-cue items | Counterfactual field adoption | Any rating session | Attribution (Q-5) claims | Study execution |
| `R-1` scenario refinement (`R-4` closure) | `R-1` architecture A | Human-rater appointment | Internal-state claims (S-5 forbids regardless) | Threshold adoption |
| Validation gates 3–4 on any α implementation | Human-rater clarification | Reproducibility claims | — | `G2`, `P-3`, `L4` |
| Packet, forms, schema design | Threshold / NA / floor | — | — | — |

**No path requires a second model to make progress on design.** Only the rating itself does.

## 5. Surviving objections

1. **`R-4` boundary is narrowed, not closed** (row 452 §6). A reader who treats "the moment" as the relevance domain can still classify the `R-1` scenario as irrelevance.
2. **The counterfactual is authored by the same reader of the same item** (row 451). Its independence is structural, not authorial — it can be wrong, which is the only defence offered.
3. **Human raters may be one source in two bodies** (row 453 §5). Disclosure mitigates; it does not establish independence.
4. **`R6-01` is stricter than its mode.** `AM-01` made `R-6` silence-satisfiable while the item's mandatory requires a chosen route — a source-level inconsistency recorded, not repaired.
5. **Gate 3 is unmet**, so no implementation is yet certified for real use.

## 6. Owner decisions required (5)

| # | Question | Options | Recommended | Consequence |
|---|---|---|---|---|
| 1 | Counterfactual field | Pilot additive-only / universal / reject | **Pilot additive-only** | Universal use would launder S-2 failures |
| 2 | `R-1` architecture | A tact / C redundant / defer | **A** | Without it, `R-1` stays unscoreable |
| 3 | Cue-dependence verdict | Separate verdict / compliance-only | **Separate** | Compliance-only cannot distinguish memory from caution |
| 4 | Human raters | Clarify admissible / restrict to models / defer | **Clarify admissible** | Without clarity, any human study rests on an unrecorded assumption |
| 5 | Statistical policy | Adopt gate only / gate + thresholds / defer | **Gate only** | Thresholds without a validated implementation are meaningless |

## 7. Claim ceiling

Established: one of three counterfactual drafts is independently falsifiable; tact is operationalizable with four scenario elements; the authority is silent on rater type; a validated α implementation is achievable but not yet certified. **Not established:** construct validity, reliability, reproducibility, recognition, attribution, or any model behaviour. Nothing was scored, executed, recruited, or amended.
