# MEMORY-09 — A1 Preregistration (proposed)

**Status:** proposal, not approved. Fixes the parameters plan v3 §4.4 left open (`N`,
repetitions) **before** any run, so G4 can be scored without post-hoc choices. No implementation
exists; nothing here has been executed.

**Binds to:** plan v3 (`6c4078bd…`), G0 receipt (`5c245283…`), reconciliation record.

**Rule of this document:** every number below is fixed now. If a number later proves wrong,
that is recorded as a **deviation**, not silently adjusted. Choosing `N`, seeds or probe quotas
after seeing outcomes is the tuning the plan forbids (§2.6).

---

## 1. Units: what is independent and what is repeated

This distinction governs the whole analysis, so it is settled first.

| Level | Unit | Counted as evidence? |
|---|---|---|
| **Conversation** | one scenario + one evidence log, distinct content | **Yes — the independent experimental unit** |
| **Probe** | one question within a conversation | No — repeated measure within a unit |
| **Generation rep** | one sampling run of the same conversation under the same arm | No — repeated measurement, averaged per conversation |

**Consequence:** the headline statistic is computed per **conversation**, then compared across
conversations. Turn-level and probe-level numbers are descriptive only. Treating probes as
independent would inflate significance by roughly the number of probes per conversation.

---

## 2. Fixed parameters

| Parameter | Value | Note |
|---|---|---|
| Independent conversations `N` | **16** | distinct scenario seeds, fixed before the run |
| Arms analysed for capability | **A0, A1, D** | A2 inert; C excluded from capability scoring |
| Repetitions per conversation × arm | **3** generation seeds | repeated measures, averaged |
| Design | **paired** | every conversation appears in every analysed arm |
| Probes per conversation | **13** | composition in §5 |
| Post-rollover probe instances per arm | 16 × 13 = **208** | descriptive base |
| Model calls (upper bound) | 16 × 3 arms × 3 reps × (turns + probes) | budgeted before the run |

Arms A0 and A1/D use the **same** conversation, queries, window capacity and sampling; only the
memory block differs (plan §2.2). Arm D's block is matched to A1 per turn (plan §2.5).

**A0 cannot produce a memory-supported answer by construction** — it has no memory block. It is
the floor: what the model achieves from the window alone. This is stated so that nobody later
reads A0's zero as a defect.

---

## 3. Why `N` = 16

No pilot data exists, so a power calculation against an assumed effect size would be invented
precision. Instead:

1. **`N` is set by budget, not by a target result.** Sixteen conversations is what the apparatus
   can run across three arms and three reps without the run becoming unattributable to any one
   scenario.
2. **The design is paired**, so each conversation is its own control and between-scenario
   variance does not enter the primary comparison.
3. **The exact test's resolution is stated in advance.** With 16 pairs, a two-sided sign test at
   α = .05 requires **≥ 13 of 16** concordant pairs (P(X≥13) = 697/65,536 ≈ .0106, two-sided
   ≈ .021); at 12/16 it is ≈ .077 and **not** significant. So this design detects **large,
   consistent** effects only.
4. **Anything smaller is reported as indeterminate**, not as a null result and not as a trend.

**This is a limitation accepted now, not discovered later.** An `N` chosen to make a smaller
effect significant would have to be much larger, and choosing it after seeing the data is
exactly what preregistration prevents.

---

## 4. Repetition structure

- **Three generation seeds per conversation × arm**, fixed and recorded before the run.
- Reps exist to reduce sampling noise, not to multiply sample size: the three are **averaged into
  one per-conversation value** before the paired comparison, and their spread is reported
  separately as measurement noise.
- **Independent content, not replayed scenarios, for new units.** A seventeenth conversation
  would need new evidence; re-running the same sixteen with new seeds adds repetitions, never
  units.
- If the runtime reports `seed_honoured` false or unavailable, that rep is **invalid and re-run
  under identical conditions**; it is not replaced by a different conversation.
- Conversation seeds are committed before the first call. No conversation is added, dropped or
  re-seeded after results are visible.

---

## 5. Probe composition and rollover coverage

Fixed per conversation, identical across arms:

| Probe class | Count | Requirement |
|---|---|---|
| character-specific (owner A only) | 4 | evidence outside the window at probe time |
| ambient / world-scope | 2 | same |
| restricted (authored, narrow audience) | 2 | same |
| experiential (learned by a character in-conversation) | 2 | same |
| correction (value superseded after a rollover) | 1 | correction occurs post-rollover |
| retraction (asserted fact later retracted) | 1 | retraction occurs post-rollover |
| unanswerable (correct response is abstention) | 1 | no eligible evidence exists |

**Rollover coverage per conversation:** at least **3** window rollovers after the earliest
planted fact, and **every** probe asserts `supporting_span ∉ window.present_spans()` at the
moment it is asked (`Window.rollovers[].dropped_spans` is the evidence). A probe whose span is
still in the window is an **instrument failure**, not a data point.

**Lifecycle placement:** correction at ≈40% of turns, retraction at ≈65%, each followed by at
least one rollover before its probe. Placement is fixed before the run.

**Knowledge asymmetry:** two characters in every conversation, with at least two facts known to
only one. The foreign-owner probe is part of the character-specific class.

**P5** (`Who trained the horse?`) is run as its own case in every arm and reported standalone,
not folded into any aggregate.

---

## 6. Reporting

**Primary (inferential):**

- Per-conversation proportion of **memory-supported answers** among probes with confirmed window
  absence.
- Paired differences **A1 − A0** and **A1 − D**, reported per conversation, all 16 pairs listed.
- Exact paired sign test, two-sided, at the resolution in §3; a 95% interval for the paired
  difference.
- **Vetoes reported first.** Per plan §4.3, any veto fails the run: no aggregate is offered.

**Secondary (descriptive, never substituted for primary):**

- Per-class breakdown of the four ladder levels (plan §4.2).
- Token accounting: `tokens.total` delta per turn with variation (plan §2.4).
- Retrieval diagnostics: selected ids, abstention reasons, `not_compared` where applicable.
- A1-vs-D structural fidelity: per-turn token parity and record-shape match (plan §2.5).

**Explicitly not reported as success:** answer-correct without support; any A0 or D
"memory-supported" figure (undefined by construction for A0); any composite "memory score".

---

## 7. Stopping rule

**Fixed sample, no peeking.**

- No interim analysis of outcomes. The run ends when the preregistered cells are complete.
- **Outcome-blind halt conditions only:**
  1. more than **10%** of probes across the run lack a receipt;
  2. token accounting missing for any arm or turn;
  3. an instrument assertion fails (probe span still present; record cites a nonexistent span);
  4. decoy pool cannot meet per-turn token parity within the declared tolerance.
- A halt is **not** a result. The instrument is repaired and the **entire** fixed sample is
  restarted; no partial data is reused, because a repaired instrument changes the measurement.
- **No extension after seeing outcomes.** If the result is indeterminate, the response is a
  **new preregistration** with a larger `N`, disclosed as a replication — never "a few more
  conversations" on the existing sample.

---

## 8. What this preregistration does not decide

- It does not approve any architecture rule; the contract's proposals remain proposals.
- It does not resolve `RQ-MF-01`; A1 still uses deterministic formation.
- It does not set retrieval thresholds or gates; those come from the existing calibration and are
  recorded with the run, per plan §2.6.
- It does not authorize G1 or G4. Both require separate authorization.

## 9. Deviation protocol

Any departure from §2, §4, §5 or §6 is recorded in the G4 receipt with:

- what changed, when, and why;
- whether it was decided before or after outcomes were visible;
- its effect on inference.

A deviation decided after outcomes were visible **invalidates the affected comparison**, and the
run is reported as such rather than adjusted.

---

## 10. Approval record

| Item | Status |
|---|---|
| `N` = 16 paired conversations | proposed |
| 3 generation reps per cell | proposed |
| 13 probes per conversation, quota in §5 | proposed |
| ≥3 rollovers, all probe spans absent | proposed |
| Fixed-sample, outcome-blind stopping (§7) | proposed |
| Primary statistic: paired memory-supported proportion | proposed |

**Pending owner approval.** Once approved, this document's SHA-256 is recorded in the G4
receipt, and it cannot be edited afterwards — corrections go in a dated addendum.
