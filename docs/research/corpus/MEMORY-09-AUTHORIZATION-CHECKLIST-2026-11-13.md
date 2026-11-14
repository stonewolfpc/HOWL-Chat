# `MEMORY-09` — Consolidated Authorization Checklist (Task 5, `WO-ARCH-04`)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-AUTH-1` · **Version:** 1 · **Date:** 2026-11-13 · **Status:** **CHECKLIST — a map of what is and is not permitted today.**
**Authority:** `MEMORY-09-WO-ARCH-04` §5.

---

## 1. The distinction this document enforces

> **Design approval** means a document may be treated as the settled design of a thing.
> **Execution permission** means the thing may be run.

They are **separate signatures**, and every row below has two columns because conflating them is the failure mode this checklist exists to prevent: a design that is approved is not a licence to run it, and a thing that can be run is not thereby a thing whose design was approved.

## 2. Status board

| Item | Design approved? | Executable today? | What execution requires |
|---|---|---|---|
| **`R-1` formal definition** (`R1D-1`) | **No** — proposal | **N/A** | Owner adoption via a dated standards instrument (`AM-01` precedent). No item edit implied |
| **Counterfactual field** | **No** — pilot probe only | **No** | Adoption decision **and** a schema change (not authorised) |
| **`F6-01` pilot protocol v2** (`PILOT-2`) | **No** — final draft | **No** | Design approval **+** permission to generate replies **+** a model run (all separately authorised) |
| **Human rater admissibility** | **No** — decision brief only | **No** | Owner clarification (option A) **+** rater appointment (not authorised) |
| **α implementation** | **No** — in-session script | **No** | Filing the script as an artifact, then re-running all five gates from the filed version |
| **α thresholds (0.667 / 0.800)** | **No** | **No** | Owner decision, fixed before any result is seen (`E-7`) |
| **20-unit floor** | **No** | **No** | Owner decision |
| **NA-as-category policy** | **No** | **No** | Owner decision (changes what α means) |
| **Reproducibility study** | Design partially blocked | **No** | `E-4` decided, two independent raters, threshold adopted, explicit authorisation |
| **Corpus amendments** | **No** | **No** | Separate dated instrument per change |
| **`G2` / instrumentation / `L4` / `P-3`** | **No** | **No** | Not authorised in any order to date |

## 3. Reading the board

**Everything executable is currently unauthorised; everything authorised is not executable.** That is the intended state after a readiness phase — the work has converted unknowns into named decisions, not into permissions.

**The three items nearest to runnable:**

1. **`F6-01` pilot** — design is final (v2), control logic is specified, decision rules are pre-registered. It needs **two signatures**: design approval, then execution permission (which authorises generating ~15 replies).
2. **Human rater path** — needs **one clarification**, then rater appointment, which is a separate act.
3. **α certification** — needs **one authorisation** to file a script; the gates themselves are already specified and passing.

## 4. Design approval, per item — what exactly is being approved

| Item | If approved, what the owner is asserting |
|---|---|
| `R-1` definition | That tact is the intended construct, and that one rater judgement is acceptable |
| Counterfactual field | That additive-class items may carry a predicted-reply contrast |
| Pilot v2 | That A/B/C with E1–E6 and K1–K5 is a sound test of content-sensitivity |
| Human raters | That `OD-M09-17`'s "rater" is type-neutral |
| α implementation | That the five gates are sufficient evidence of arithmetic correctness |
| Thresholds | That a specific α value will be treated as the bar — **before** any number exists |

## 5. Execution permission, per item — what exactly is being permitted

| Item | If permitted, what may then happen |
|---|---|
| Pilot v2 | Generating replies for A, B, C; observing; recording counts. **Not** scoring, **not** `G2` |
| α on real data | Computing α over ratings that exist. **Not** reporting a claim beyond single-source unless two raters |
| Study execution | Two independent raters rate the units. **Not** promotion beyond single-source until the protocol is satisfied |
| Corpus amendment | Editing one named item under one named instrument |

## 6. What is already true without any signature

- Nine analysis artifacts exist and are registered (rows 447–457 plus 458–462) with hashes.
- `R-1` is defined to the limit the construct allows; `R6-01` is reconciled; the pilot design is final in draft.
- The α arithmetic is certified for this session, with two published reproductions on record.
- **No corpus file, standard, tag, split, score or checker line was changed** across the entire sequence.

## 7. The one-sentence summary

**Nothing here is approved and nothing here has run — and the distance between "approved design" and "permitted execution" is now exactly one signature per item, with the work behind each signature already done.**
