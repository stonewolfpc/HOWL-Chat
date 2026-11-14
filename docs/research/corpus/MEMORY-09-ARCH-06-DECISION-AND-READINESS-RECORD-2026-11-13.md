# `MEMORY-09` — `WO-ARCH-06` Decision and Readiness Record

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-READY-6` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **READINESS — design complete; nothing executed, nothing amended.**
**Authority:** `MEMORY-09-WO-ARCH-06` §15.

---

## 1. Completed design work

| # | Artifact | Row | What it settles |
|---|---|---|---|
| 1 | Owner-Mediated Evaluation Architecture | 470 | Four layers with claim ceilings; the two-rater split; the A–D record; labelling discipline |
| 2 | Model Population and Hardware Plan | 471 | Eligibility rule, seven categories, ten candidate slots, bounded measurement |
| 3 | Balanced Chat Baseline | 472 | Two-layer profile; published values only; freeze before test |
| 4 | Controlled Engineering Test Matrix | 473 | Eleven tests with hypotheses, repetitions, failure conditions, ceilings |
| 5 | Owner Review and Adjudication Protocol | 474 | Presentation order, five owner actions, override record, escalation triggers |

## 2. Proposed decisions (recommendations, not approvals)

| # | Proposal | Recommended | If declined |
|---|---|---|---|
| 1 | **Two-rater split** (§470 §3) | Adopt | Engineering results stay blocked behind a requirement meant for different claims |
| 2 | Population shape (§471) | Ten slots, seven categories, non-exhaustive | Evaluation drifts into collecting similar models |
| 3 | BCB as two layers (§472) | Adopt, reusing family defaults | Baseline becomes either a creative preset or invented numbers |
| 4 | Repetition doctrine (§473 §1) | Four reasons only, no universal N | Compute spent without an evidential purpose |
| 5 | Presentation contract (§474 §2) | One sentence, three observations, then evidence | Records become unreadable and adjudication degrades to trust |
| 6 | `C2` vs `C3` classification mandatory (§473, T-11) | Mandatory | Every favourable result becomes unattributable |

## 3. Owner decisions required (narrow)

1. **Confirm the two-rater split** — that engineering observation is not a cue-level reproducibility claim. *This is the load-bearing decision; without it, nothing below is usable.*
2. **Approve the population shape** (or name substitutions). Burden: scope.
3. **Approve BCB as the frozen reference profile.** Burden: design adequacy.
4. **Approve the presentation contract and override record.** Burden: legibility, and only the owner can judge it.
5. **Authorise any first test** — separate from approving the matrix. Burden: execution.

## 4. Unresolved authority

| # | Question | State |
|---|---|---|
| 1 | Human-rater admissibility (`OD-M09-17`) | Unresolved by silence; clarification drafted, not issued |
| 2 | `R-1` adoption as construct | Proposed (row 458); no instrument |
| 3 | Whether engineering results may ever be cited alongside cue-level results | **Open** — they answer different questions; a joint claim needs a stated bridging rule |
| 4 | Alpha thresholds, sample floor, NA policy | Unapproved; and α is excluded from repeated-run reporting here |
| 5 | F6-01 pilot authority | Approved design? **No.** Executed? **No.** |

## 5. `F6-01` review — retained, unexecuted, unapproved

The pilot **remains useful** as one narrow cell of the broader matrix (it is closest to T-04 with a content-sensitivity contrast). It is **not** subsumed and **not** discarded. Three differences from the engineering matrix must stay visible:

| Dimension | F6-01 pilot | Engineering matrix |
|---|---|---|
| Unit of evidence | Content-sensitivity between two conditions | Behaviour under recorded conditions |
| Decision rules | Pre-registered 4-of-5 counts | Per-test failure conditions |
| Claim ceiling | Content-sensitive compliance for one item | Behaviour for one model/scenario |

**No conflict of authority**, but three constraints: the pilot is **not executed** by this order; its **frozen packets are unaltered**; and this architecture **does not authorise it** — approving a matrix is not approving a pilot.

## 6. Alpha boundary — restated

The filed utility (row 463) certifies **arithmetic**. It does not approve thresholds, sample floors, missing-data policy, or statistical suitability. And specifically: **α is not used for repeated-run variability** — that is a different quantity, and reporting it would be a category error the protocol now forbids.

## 7. Prepared but not executed

| Prepared | Status |
|---|---|
| Eleven test designs | `NOT RUN` |
| Ten model slots | `NOT TESTED` — no model loaded |
| BCB parameter policy | No setting applied to any run |
| Record templates | No record exists |
| Measurement plan | No measurement taken |

## 8. Still unauthorised

Model runs; pilot execution; load/performance measurement; rater recruitment; corpus, tag, split or score changes; standards amendments; alpha thresholds, floor or NA policy; `G2`, `P-3`, `L4`; any change to Howl's source or runtime.

## 9. What the owner will be able to inspect and decide

> **You will be handed, for each test, one sentence saying what happened, three observations explaining why it matters, and — underneath — the exact settings, the full unedited reply, and the stored state at that moment. You will be able to accept it, argue with the interpretation while keeping the evidence, override the conclusion and record why, or ask for a smaller test instead of starting over. Nothing an AI concludes will be recorded as your decision, and nothing you decide will delete an observed defect.**

The single most useful output the architecture is built to produce is the distinction between *"right because the system supplied the fact"* and *"right because the model happened to know it"*. Without that line, every favourable result is unattributable, and the evaluation would measure the model twice while calling one of the measurements architecture.

## 10. Claim ceiling

Design only. No model was loaded, no reply generated, no state inspected, no measurement taken, and no standard amended. Nothing here indicates that `MEMORY-09` improves behaviour.
