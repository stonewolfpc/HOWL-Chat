# MEMORY-09 — Preregistration Amendment 3: MMD Interpretation and Calculation

**Status:** proposed append-only clarification, awaiting approval. Documentation and design only.
No implementation, no subject run, no new results, no frozen artifact changed, no commit.

**Append-only.** Amendments 1 (`97503340…`) and 2 (`c3ffd7ed…`) and every earlier document are
preserved unchanged. This amendment adds interpretation and calculation rules only.

**Preserved:** `N` = 16, 3 matched reps, 13 probes, seed enforcement, fixed-sample stopping, the
α = .025 per-claim test, the C1/C2 structure, tie handling, and the MMD **value** of 0.10.

---

## 1. What 0.10 is — and is not

**0.10 is an operational threshold selected for this experiment.** It is **not** a predeclared
practical-significance standard, and no real-world significance is established by it.

The distinction matters, so it is stated twice:

- There is **no external anchor** for 0.10 — no prior result, no user study, no field benchmark —
  that would license the phrase "practically significant".
- 0.10 exists for one purpose: to stop a **statistically** significant but trivial difference from
  being reported as a capability gain.

**Therefore, in any report:** "the median paired improvement exceeded the preregistered minimum"
is correct; "the improvement is practically meaningful" is **not** — the latter remains
unestablished and must not be inferred from this threshold.

### 1.1 Why 0.10 rather than a round one-probe value

M1 is a mean over 13 probes, so **one probe-equivalent = 1/13 ≈ 0.0769**.

| Candidate | Probe-equivalents | Disposition |
|---|---|---|
| 1/13 ≈ 0.077 | 1.0 | rejected: one probe in the median is too thin to call an improvement |
| **0.10** | **≈ 1.3** | **adopted**: strictly above the one-probe floor, fixed before data |
| 2/13 ≈ 0.154 | 2.0 | rejected: too demanding to be a *minimum* meaningful difference |

0.10 is therefore a round number chosen to sit **above** one probe-equivalent, not a value
derived from an effect-size estimate or from any result.

### 1.2 The threshold's exact attainable value

Because the median is computed on averaged repetitions (§2), it cannot take every real value:

| Quantity | Granularity | Consequence |
|---|---|---|
| one rep's M1 | `k/13` | 13 possible values |
| cell value (3 reps averaged) | `m/39` ≈ 0.0256 | — |
| `d_i` | multiples of `1/39` | — |
| **median of 16 `d_i`** | multiples of `1/78` ≈ 0.0128 | **0.10 is unattainable** |

The smallest attainable median that satisfies the threshold is therefore **4/39 ≈ 0.1026**.
In practice the rule reads: *the median paired improvement must reach at least four
probe–repetition units*. This is recorded so the criterion is reproducible to the last decimal
rather than "about 0.10".

---

## 2. How the median is calculated

Stated precisely, so two analysts cannot differ:

1. **Per rep:** `M1(arm, conv, rep)` = mean credit over that rep's 13 probes
   (Amendment 1 §2 scoring: `CORRECT` = 1 on answerable probes, `ABSTAIN_APPROPRIATE` = 1 on the
   unanswerable probe, every other category 0).
2. **Per cell:** `M1(arm, conv)` = arithmetic mean of the conversation's **3 matched reps**,
   equal weights of 1/3, using the same seeds across arms (Amendment 1 §5).
   **A cell with any invalid rep is invalid** — it is never averaged over the two remaining reps,
   because unequal rep counts would make cells incommensurable.
3. **Per pair:** `d_i = M1(A1, conv_i) − M1(X, conv_i)` for X ∈ {A0, D}.
4. **Median:** over **all valid pairs, including exact-zero ties** — the median uses the full
   `n_valid`, *not* the `n_eff` of the sign test. These two denominators are deliberately
   different: the sign test drops ties because a tie carries no direction; the median must
   include them because dropping them would inflate the reported magnitude.
5. **Even counts:** arithmetic mean of the two central order statistics (standard median, no
   interpolation or smoothing).

**Reported with each claim:** `n_valid`, tie count, all `d_i` sorted, median, quartiles, and the
number of pairs meeting `|d_i| ≥ 0.10` (the declared secondary sensitivity view, Amendment 2
§3.4).

**Excluded pairs** (Amendment 2 §4: missing receipt, span still present, unhonoured seed, absent
token accounting) are removed from `n_valid` **and** from the median — they are not zero-filled,
and they are not counted as ties.

---

## 3. Worked illustration (no data — arithmetic only)

To show the rule is unambiguous, take a hypothetical pair set of three conversations with
per-rep probe credits summing to 9, 8, 9 of 13 for A1, and 7, 8, 7 for D, with one rep invalid in
conversation 2:

- Conversation 1: A1 = 9/13 = 0.6923; D = 7/13 = 0.5385; `d = 0.1538`.
- Conversation 2: **excluded** (invalid rep) — not averaged over its two good reps.
- Conversation 3: A1 = 9/13 = 0.6923; D = 7/13 = 0.5385; `d = 0.1538`.

`n_valid` = 2, median = 0.1538, no ties. If a fourth conversation gave `d = 0`, `n_valid` = 3,
median still 0.1538, tie count 1, while `n_eff` = 2 for the sign test. This is the intended
behaviour: the tie lowers the sign test's sample size but not the reported magnitude.

*This illustration is arithmetic, not a result; no model was called.*

---

## 4. What this amendment changes

| Item | Change |
|---|---|
| MMD value 0.10 | **unchanged** |
| MMD interpretation | now explicitly **operational**, not practical significance |
| Attainable threshold | recorded: smallest satisfying median is 4/39 ≈ 0.1026 |
| Median calculation | now specified: reps → cell → pair → median over all valid pairs incl. ties |
| Invalid rep handling | cell invalid, never partially averaged |
| `N`, test, α, stopping | **unchanged** |

---

## 5. Freeze procedure (on approval)

Once approved, the preregistration set is frozen as:

1. `MEMORY-09-PREREGISTRATION-2026-11-15.md` (`8c788a96…`) as modified by
   Amendment 1 (`97503340…`) and Amendment 2 (`c3ffd7ed…`) and this Amendment 3.
2. Their four hashes recorded in the G4 receipt before the first model call.
3. Any later change is a **dated, append-only amendment**, never an edit — including changes a
   reader might consider cosmetic.

---

## 6. Provenance

| Document | Status |
|---|---|
| `MEMORY-09-EXPERIMENT-PLAN-2026-11-15.md` | unchanged, `6c4078bd…` |
| `MEMORY-09-G0-PREFLIGHT-2026-11-15.md` | unchanged, `5c245283…` |
| `MEMORY-09-G0-RECONCILIATION-2026-11-15.md` | unchanged, `2b9c55b6…` |
| `MEMORY-09-PREREGISTRATION-2026-11-15.md` | unchanged, `8c788a96…` |
| `MEMORY-09-PREREGISTRATION-AMENDMENT-1-2026-11-15.md` | unchanged, `97503340…` |
| `MEMORY-09-PREREGISTRATION-AMENDMENT-2-2026-11-15.md` | unchanged, `c3ffd7ed…` |

G1 and G4 remain unauthorized.
