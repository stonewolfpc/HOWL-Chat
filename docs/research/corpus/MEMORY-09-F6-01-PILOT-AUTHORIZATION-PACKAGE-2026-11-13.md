# `MEMORY-09` — `F6-01` Pilot Authorization Package

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-PKG-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **FROZEN AND READY FOR APPROVAL — NOT AUTHORISED, NOT EXECUTED.**
**Authority:** `MEMORY-09-WO-ARCH-05` §4.
**Protocol:** `MEMORY-09-PILOT-2` (row 459) as amended by the packet specification below.
**Packets:** `MEMORY-09-PKT-1` (row 466).

---

## 1. Frozen inputs and hashes

SHA-256 (first 32 hex characters shown; full value in the packets file header):

| Element | Length | SHA-256 (truncated) |
|---|---:|---|
| Scene block | 69 | `95bd673d84a4ba20330f8a7e757469b1…` |
| **A** memory line | 56 | `ec2a5ffa817909293180c627664d7c16…` |
| **B** memory line | 57 | `f89a69948ead3d612b1127e6e5ae71fb…` |
| **C** memory line | 56 | `533fe7b5e25d40fab790a584190c311c…` |
| Packets file | — | `3e019c84d4a5417847d9ffe2b27219d2…` |

**Any edit to any packet invalidates the hashes and the pre-registered rules.**

## 2. Equivalence results — measured, not asserted

| Check | Requirement | Measured | Status |
|---|---|---|---|
| E1 position | message 44 in all | 44 / 44 / 44 | **PASS** |
| E2 speaker & register | same, casual | identical frame | **PASS** |
| E3 length parity | A vs B within ±10% | **56 vs 57 = 1.75%** | **PASS** |
| E4 syntactic frame | "can't eat X" for A and B | literal match | **PASS** |
| E5 plausibility parity | both ordinary food restrictions | declared: neither medical, ethical nor identity-related | **PASS** (declaration) |
| E6 salience parity | fact appears exactly once in A and B | substring scan: almond only in A, cilantro only in B, neither again | **PASS** |
| C vs A length | not required; recorded | 56 vs 56 = 0.00% | *informational* |
| C structure | preference about tea, not a served food | differs by design, rationale in row 466 §6 | *known difference* |

## 3. Contamination controls

| # | Contamination | Detection | Consequence |
|---|---|---|---|
| **K1** | A's stimulus contains "cilantro"; B's contains "almond" | case-insensitive substring scan of assembled stimulus | **VOID RUN** |
| **K2** | Either served-food or restriction word appears in the serving scene | scan of scene block | **VOID RUN** |
| **K3** | A B-reply references almonds | scan of B replies | **VOID RUN** |
| **K4** | An A-reply references cilantro | scan of A replies | **FLAG**, not void |
| **K5** | Restriction appears in preamble or persona card | scan of non-item context | **VOID RUN** |

Each is assessable from the assembled text alone — **no hidden model state needed.** No control was weakened to ease execution.

## 4. Blinding

- **Given to the rater:** the reply text, the serving scene, and the item statement of the acceptance condition.
- **Withheld:** condition labels, condition order, all aggregate counts.
- **Procedure:** replies stripped of labels, shuffled with a recorded seed; the rater rates every reply, then records a `condition_guess` per reply.
- **Guess handling:** guesses are stored **beside** the ratings and never alter them. A high guess rate is reported as failed blinding and downgrades the result to `indeterminate`.
- **No claim of successful blinding is made in advance.**

## 5. Decision rules — proposals, not thresholds

| Class | Rule | Rationale |
|---|---|---|
| **Informative** | A: `names_restriction` ∨ `substitutes_or_warns` in ≥ 4 of 5; **and** B: `offers_unqualified` in ≥ 4 of 5 | Content-sensitivity requires the reply to track *which* restriction, not merely that one exists |
| **Null** | A and B indistinguishable | The counterfactual is uninformative even for the additive class |
| **Confound** | B replies largely `asks_about_diet` or adapt to cilantro | Effect is generic caution → only A vs C usable |
| **Void** | Any K1, K2, K3 or K5 failure | A contaminated comparison measures nothing |
| **Unusable** | Either condition mostly `not assessable` | Scene too thin; rebuild |

**The 4-of-5 rule and 5 repeats per condition are proposals.** They are not statistical thresholds, not α, and carry no authority until approved.

## 6. Execution records required

Per reply: condition (sealed until rating), reply text, four booleans, `condition_guess`, rater identifier, timestamp. Per run: packet hashes, seed, K1–K5 scan outputs, blinding guess distribution, counts per class, and the disposition.

## 7. Conditions that invalidate the run

1. Any K1/K2/K3/K5 hit.
2. A stimulus edit after the first reply is generated (`S-3` violated).
3. Blinding broken and undetectable (guess rate at ceiling with no recording).
4. Fewer than 5 replies in any condition.

## 8. Claim ceiling

**Permitted:** "content-sensitive compliance under the tested conditions", for `F6-01`, in one model, at single-source confidence.

**Not permitted:** internal recognition; retrieval; hidden state; unrestricted causal attribution; general construct validity; reliability; reproducibility; any `R-2` verdict for the corpus; any statement about `MEMORY-09` overall.

**Alternatives that remain live even on an informative result:** generic caution; accidental compliance; a model that asks about diets by habit; a reply that adapts for reasons unrelated to the stored fact.

## 9. Approvals required before any reply is generated

| # | Approval | Status |
|---|---|---|
| 1 | Approve this package as the pilot design | **NOT GIVEN** |
| 2 | Authorise execution: generating ~15 replies across A/B/C | **NOT GIVEN** |
| 3 | Authorise a model run of any kind | **NOT GIVEN** |
| 4 | Confirm the rater (and that single-rater results are reported at single-source confidence) | **NOT GIVEN** |

**No model run, reply generation, rating or condition assembly has occurred.**
