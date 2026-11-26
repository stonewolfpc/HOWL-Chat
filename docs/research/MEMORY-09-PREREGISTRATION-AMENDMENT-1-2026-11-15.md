# MEMORY-09 — Preregistration Amendment 1: Metric Symmetry

**Status:** proposed amendment, awaiting approval. Documentation and design only. No
implementation, no run, no new results, no change to frozen artifacts.

**Supersedes for inference:** plan v3 §4.2/§9, preregistration §6, in the specific respects
listed in §7. The plan and preregistration files are **not edited**, so their hashes stay valid
as authorized; this amendment is the record of the change.

**Finding that motivates it:** the success criterion as written compares A1 against arms on a
metric that those arms **cannot attain by construction**. That is a genuine experimental-design
defect, not a documentation nuance: a difference that is guaranteed by construction looks like
evidence of capability. It is corrected below.

**`N` = 16 and the stopping rule are unchanged.** The defect is in outcome definition, not in
sample size, so sample size is not the remedy.

---

## 1. Why the current comparison is unsound

| Arm | Can it select a supporting record? | Can it produce a memory-supported answer? |
|---|---|---|
| A1 | yes | yes, by design |
| D | yes, but records are decoys by construction | **no, except by instrument failure** |
| A0 | **no records exist** | **no — undefined, not zero** |

Consequences of ignoring this:

1. **A1 − A0 on memory-supported** is a comparison between "possible" and "impossible". Any
   positive value is guaranteed; it measures the definition, not the architecture.
2. **A1 − D on memory-supported** has the same flaw in a subtler form: D's value is pinned near
   zero *by construction of the decoys*, so the contrast is near-structural too.
3. **A0 scored 0** on that metric would then be averaged into a "memory score", presenting a
   structural zero as a performance failure.

The fix is not to drop the memory-specific measure. It is to separate **a common outcome all
arms can be scored on** from **a memory-specific quantity whose applicability differs by arm**,
and to state which comparison licenses which claim.

---

## 2. M1 — answer quality (common outcome, identical scoring in every arm)

**Definition.** For every probe, the response is classified into exactly one category by a
scorer that sees only the response and the probe's key, never the arm:

| Category | Condition |
|---|---|
| `CORRECT` | matches the probe's key (exact/normalised match, or key-listed acceptable variants) |
| `INCORRECT` | contradicts the key, or answers a different question |
| `ABSTAIN_APPROPRIATE` | the probe is unanswerable from eligible evidence, and the response declines or says the information is not established |
| `ABSTAIN_INAPPROPRIATE` | declines although eligible evidence exists |
| `MALFORMED` | unusable output (empty, truncated mid-claim, refusal to engage) |

**Scoring rule (symmetric by design):**

| Probe kind | Credit |
|---|---|
| answerable (12 of 13 per conversation) | `CORRECT` = 1; every other category = 0 |
| unanswerable (1 of 13) | `ABSTAIN_APPROPRIATE` = 1; every other category = 0 |

**M1 per conversation** = mean credit over its 13 probes. **M1 per arm** = mean of the 16
conversation values (equal weight per conversation, never pooled per probe).

**Why this is scoreable identically everywhere:** it depends only on the response and the key.
A0, A1, D and even C are scored by the same rule, so their values are comparable without any
arm-specific caveat. This is the metric the inferential comparisons use.

**Scorer identity:** deterministic matcher by default. If any response needs judgement, that is a
declared second instrument — arm-blind, with its own identity and an agreement check against the
deterministic matcher on a fixed subset. Scorer identity goes in the run receipt.

---

## 3. M2 — evidence-supported answers (memory-specific, arm-dependent applicability)

**Definition.** For a probe, M2 credit requires **all** of:

1. the response is `CORRECT` (M1 credit = 1 on an answerable probe);
2. at least one record selected at that turn **supports** the keyed answer (a declared
   support check, same instrument for every arm);
3. the probe's supporting span was **absent from the window** at that turn
   (`Window.rollovers[].dropped_spans`).

**Applicability by arm — the part that must never be averaged across:**

| Arm | M2 status | Interpretation |
|---|---|---|
| A1 | **applicable** | answers produced from retrieved, window-absent evidence |
| D | **applicable but expected ≈0** | non-zero means decoys carry relevance → **instrument failure**, investigate before any comparison |
| A0 | **not applicable (N/A)** | no records exist; reported as `N/A`, never as `0`, never averaged |
| C | not applicable | injected records are plumbing, not formation |

**M2 is reported as an absolute rate for A1** (per conversation, with interval), and for D as a
**manipulation check**. It is **not** used as the inferential contrast, precisely because its
attainability differs by arm.

---

## 4. What each contrast establishes

| Contrast | Metric | Establishes | Does **not** establish |
|---|---|---|---|
| **A1 − D** | **M1** | that *relevant selection* at matched token cost and structure improves answers | that memory alone caused it — retrieval quality also varies |
| **A1 − A0** | **M1** | the **total** capability effect of the memory system versus no memory | anything memory-specific: extra tokens, formatting and retrieval are all inside this difference |
| A1 vs D | M2 | manipulation check: relevance, not tokens, drives support | a capability claim on its own |
| A1 absolute | M2 | how often correct answers rest on retrieved, window-absent evidence | causality of use (§4.2 of the plan's causal caution still applies) |

**The capability claim requires A1 > D on M1**, because M1 is attainable by every arm and D holds
tokens, shape and placement constant. **A1 − A0 is reported as the total effect** and is never
described as proof that memory works — it is the headline number a reader would misread, so it is
labelled as total, not memory-specific.

**Prohibited presentations:** a composite "memory score"; any aggregate mixing M1 and M2; any
table where A0's M2 appears as a value rather than `N/A`.

---

## 5. Seed matching across arms

**Yes — seeds are matched per `(conversation, rep)`**, and matching is enforced, not assumed.

- **Seed derivation:** `seed(c, r) = H(conversation_id ‖ rep_index)` computed **before** the run
  from a committed table, so all three arms use the identical value for the same cell.
- **Enforcement (mechanical):**
  1. the seed table is written to the run directory and hashed before the first call;
  2. every call's receipt records `seed` and `seed_honoured`;
  3. after the run, a check asserts the multiset of seeds per `(c, r)` is identical across arms
     and fails the run otherwise;
  4. any cell with `seed_honoured` false/unavailable is **invalid and re-run under identical
     conditions** — never replaced by another conversation.
- **Ordering:** cells are run in interleaved `(c, r)` order with the arm order rotated across
  conversations, so drift in the runtime cannot systematically favour one arm.

---

## 6. Changes this amendment makes to the preregistration

| Preregistration item | Change |
|---|---|
| §2 parameter table | unchanged: `N` = 16, 3 reps, 13 probes |
| §3 rationale | unchanged; the 13/16 resolution now applies to **M1** |
| §4 repetition structure | §5 of this amendment adds the seed table and its enforcement |
| §6 primary statistic | **replaced**: primary = paired **M1**, A1 − D; secondary inferential = paired M1, A1 − A0 |
| §6 secondary reporting | M2 becomes an absolute rate for A1 + manipulation check for D; A0 shown as `N/A` |
| §7 stopping rule | unchanged |
| §9 deviation protocol | unchanged; this amendment is itself the record of the change |

**Plan v3 §9 success criterion** ("memory-supported answers rise above **both** A0 and the
token-matched placebo D") is **superseded for inference** by: *M1 improves for A1 over D at
matched cost and structure, with M2 reported absolutely for A1 and ≈0 for D.* The plan's
falsification clause is amended to match: a null or reversed **M1** result at A1 vs D is a
failure of the memory claim **regardless of A1's M2 rate**.

---

## 7. Scope of the defect and what it does not change

- **G0 stands.** The gate verified plan identity, prerequisites, provenance and internal
  consistency; this defect is in the inferential use of a metric, discovered by review before any
  data existed. No G0 finding is withdrawn.
- **Plan v3's hash is preserved.** The change is recorded here, as the plan's own deviation
  protocol requires, rather than by editing an authorized artifact.
- **No sample-size change.** `N` = 16, three reps, fixed-sample stopping all stand. Making the
  metric symmetric is what allows a *smaller* effect to be detected honestly; enlarging `N` would
  be tuning, and is not done.
- **The known insensitivity remains.** With 16 pairs, only ≥13/16 concordant pairs reach
  two-sided significance. This amendment does not improve power; it stops an asymmetric metric
  from manufacturing a difference that power has nothing to do with.

---

## 8. Proposal

1. Adopt M1 as the common outcome and the sole basis for inferential comparisons.
2. Adopt M2 with per-arm applicability, A0 reported as `N/A`, and D used as a manipulation check.
3. Report A1 − D as the memory-relevant contrast and A1 − A0 as total effect, labelled as such.
4. Adopt the seed table and its mechanical enforcement.
5. Apply the §6 changes to the preregistration on approval; keep `N` = 16 and the stopping rule.

**Awaiting approval.** No implementation, no run, no frozen artifact touched, nothing committed.
