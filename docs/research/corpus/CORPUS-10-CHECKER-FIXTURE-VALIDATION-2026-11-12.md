# Checker Validation — Defective Fixtures

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** 10 of 10 fixtures behaved as predicted, **including two the checker cannot detect**, and `C-4` now has regression coverage on **both** sides of its boundary.
**Date:** 2026-11-12
**Fixtures:** `docs/research/corpus/fixtures/FIXTURE-SET-2026-11-12.md` — isolated, never part of a scored run.

---

## 1. Why this exists

A checker that passes the corpus proves only that the corpus matches the checker. To know whether the checker can **discriminate**, it must be run against material built to fail it.

This matters more than usual here: the checker was written by the same author as the corpus, and it has already been wrong once (a case-sensitive scan that missed four `over-exposure` fields).

---

## 2. Results

Run: `python docs/research/corpus/check_corpus.py docs/research/corpus/fixtures/FIXTURE-SET-2026-11-12.md`

| Fixture | Deliberate condition | Expected | Actual | Distinguishes? |
|---|---|---|---|---|
| `FX-01` | over-exposure boundary absent | `FAIL missing over-exposure: 1 [FX-01]` | identical | **Yes — FAIL** |
| `FX-02` | `R-3` with negative-only `mandatory` | `FAIL C-1 FX-02` | identical | **Yes — FAIL** |
| `FX-03` | generic evidence reference ("the name") | `UNVERIFIABLE U-1 FX-03` | identical | **Yes — UNVERIFIABLE, not pass** |
| `FX-04` | no per-unit `Split` | `FAIL missing split: 1 [FX-04]` | identical | **Yes — FAIL** |
| `FX-05` | `Required` "must name the harbour" vs `Variation` "no answer may name it" | **not detected** | **not detected** | **No — known limitation** |
| `FX-06` | restraint item declaring `over-exposure: none` | `WARN C-3 FX-06` | identical | **Yes — WARN, not FAIL** |
| `FX-07` | `mandatory` = "what happened at the mill", no supporting scenario | **not detected** | **not detected** | **No — known limitation** |
| `FX-08` | `R-3` whose `Required` permits non-use while `mandatory` demands accommodation | `FAIL C-4 FX-08` | identical | **Yes — FAIL, record contradiction** |
| `FX-09` | legitimate `R-1` non-disclosure with `Required: nothing about the debt` | **no finding for FX-09** | no finding | **Yes — correctly not flagged** |
| `FX-10` | legitimate `R-5` non-disclosure | **no finding for FX-10** | no finding | **Yes — correctly not flagged** |

Aggregate: `RESULT: FAIL (4 fail, 1 warn, 1 unverifiable)` across 10 units, with `mode: explicit paths` and declared-total comparison correctly skipped.

**`C-4` was tuned during validation, recorded rather than hidden.** Its first version fired on any `Required` beginning "nothing", which produced **6 false positives** across the repaired corpus — legitimate `R-1`/`R-5` items whose requirement genuinely *is* non-disclosure. It was scoped to `R-2`/`R-3`, the two modes where non-use contradicts the mode. Both runs are recorded: the naive version over-fired, and the scoped version fires only on `FX-08`.

---

## 3. What the fixtures establish

1. **The three outcomes are genuinely distinct.** `FX-06` produced a WARN and did **not** escalate to FAIL — a legitimate warning condition stays advisory. `FX-03` produced UNVERIFIABLE and did **not** appear as a pass.
2. **A graded restraint mode cannot be scored by a negative requirement.** `FX-02` reproduces the `F6-04` class in isolation, so that invariant holds independently of the repaired item.
3. **Absence of a finding is not a pass.** `FX-05` and `FX-07` both pass the checker silently. The corpus could contain either defect today and the tool would not know.

---

## 4. Limitations, stated as limitations

| Gap | Consequence | Mitigation |
|---|---|---|
| Semantic `Required`/`Variation` contradiction (`FX-05`) | A self-contradictory item can pass mechanically | Human adjudication; the three known conflicts are registered in `CORPUS-07` §2 |
| Undeterminable evidence (`FX-07`) | Unanswerable items can pass mechanically | `U-1` covers only explicit generic phrasing; the rest is review |
| `C-2` proper-noun check | Cannot find a missing entity that is never named (`F3-03` originally) | Review; and repairs now name entities explicitly |
| `C-5` evidence sufficiency | Scenario prose is not machine-parseable | Reviewer challenge `C-5` in `CORPUS-05` |

**`C-4` boundary coverage.** A rule that only ever fires proves as little as one that never fires. `FX-08` shows `R-2`/`R-3` non-use is caught; `FX-09` and `FX-10` show that genuine `R-1`/`R-5` non-disclosure is **not** flagged. Both directions are required, because the first implementation of `C-4` failed the second.

**No claim of complete coverage.** Five of seven defect classes are mechanically detectable; two are not. `CORPUS-09` §3 records a PASS on the repaired corpus, and that PASS means "no mechanically detectable defect" — not "no defect".

---

## 5. Boundaries

Fixtures are outside the scored corpus, were never counted in the 45 units, and changed no corpus content. No rule was altered to make a fixture pass.
