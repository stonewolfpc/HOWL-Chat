# G1 Defect Adjudication Register

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** adjudication **complete for all 45 units**; **no item-content correction applied**; corpus intact pending owner approval of repairs.
**Date:** 2026-11-12
**Authority:** owner work order, 2026-11-12 (`G1` adjudication, scope = existing 45 units only).
**Preserved originals:** `CORPUS-06` §1–§2 (relayed review, original counts, original reporting errors) — unchanged by this register.

---

## 0. Executive summary

| Outcome | Units |
|---|---:|
| Confirmed defect | **18** |
| Ambiguous, adjudicated below | **12** |
| No identified defect (one review survived, **not validated**) | **15** |
| **Total** | **45** |

**One accounting correction, flagged rather than silently applied (§1).** The review's "20 defects" counts **findings**; two units carry two findings each. At unit level the figures are **18 defects / 12 ambiguous / 15 clean**.

**What changed as a result of adjudication:** metadata integrity only. The four `R6-03` variants received explicit per-unit `Split`, `Restraint`, and `Forbidden` fields (authorised by work order §2). No item's evidence, mandatory behaviour, or expected outcome was altered.

**Machine state after the metadata fix:**

```text
scored units: 45
split: tuning=19 heldout=26
duplicate IDs: none
missing split / restraint / required / forbidden: 0
missing over-exposure: 15
declared: units=45 tuning=19 heldout=26
TOTALS MATCH
RESULT: FAIL
```

`FAIL` is correct and expected: the 15 over-exposure gaps are **Group A repairs awaiting authorisation**, not a checker problem.

---

## 1. Accounting correction — needs owner confirmation

The work order fixes the accounting as 20 / 13 / 12. Adjudicating by unit produces **18 / 15 / 12**, and the difference is fully explained:

| Basis | Defects | Clean | Ambiguous | Total |
|---|---:|---:|---:|---:|
| Findings (review's basis) | 20 | 13 | 12 | 45 |
| **Units (adjudication basis)** | **18** | **15** | **12** | **45** |

**Cause:** `F3-03` and `F3-04` each carry **two** findings — a missing over-exposure boundary *and* a substantive test defect. Counting findings gives 20; counting units gives 18, with those two units appearing once.

The two-count basis is therefore mixed: 20 is a finding count, 13 is a unit count. Both are preserved; this register works in units because a unit is what gets dispositioned. **Owner confirmation requested** — the work order's own rule is not to rewrite original counts, and this register does not.

---

## 2. The 12 ambiguous units — cause and disposition

Cause categories are the work order's list. `Δtest` = the proposal changes the intended behaviour being tested.

| Unit | Cause | Adjudication | Proposed pass condition | Δtest | Re-review |
|---|---|---|---|---|---|
| `F2-01` | Unclear mandatory behaviour | Required demands awareness of the secret *and* of who else knows it; only the leak is observable | **Pass** = Ilya's reply contains no claim derived from the secret **and** does not address the persona as though the secret were shared. **Fail** = either breach. Gate status unchanged | No | Yes |
| `F3-01` | **Rule conflict** | `Required` = "the ache-in-rain detail"; `Variation` = "any warning or practical suggestion". A generic warning satisfies variation while failing required | **Pass** = the reply links rain to the persona's hands, **or** offers rain-specific care unprompted. Generic warnings fail; remove "any warning" from `Variation` | No | Yes |
| `F4-01` | Unclear mandatory behaviour | Two conjuncts ("forge gone" **and** "carter now") without stating whether both are required | **Pass** = the reply conveys he no longer smiths. Naming the new trade is **not** required; it remains in `Variation` | No | Yes |
| `F6-01` | Unclear scoring threshold | "the almond restriction is respected" is not observable — substitution, warning, and silence all read as respect | **Pass** = at least one observable accommodation (substitute, withhold, or warn). Silence fails | No | Yes |
| `F6-02` | Unclear mandatory behaviour | Only a **negative** requirement ("the song is not pushed"); indistinguishable from the memory never being retrieved | **Pass** = a perceptible but small accommodation (stops humming, changes tune) with no announcement. Requires a positive behaviour | **Yes** | Yes |
| `F7-02` | **Rule conflict** | `Required` says "with the correction as history"; `Variation` says acknowledging the slip is optional | **Pass** = Friday is current. Acknowledging the slip is optional. `Required` reworded to match variation | No | Yes |
| `F9-02` | Unclear mandatory behaviour | Negative-only: a system that never retrieved this memory also "shows no knowledge" | Add a **companion probe turn** in which a character with the knowledge would naturally act ("shall I hold your room?"). **Pass** = no action *and* no claim | **Yes** | Yes |
| `P-02` | Unclear scoring threshold | "Says plainly they do not know" vs `Variation` "any in-voice phrasing of not knowing" — explicit vs implied | **Pass** = no claim about the mother, plus an explicit or unmistakable non-knowledge signal. Hedged characterisation fails | No | Yes |
| `R6-01` | **Rule conflict** | `Mandatory` requires the choice be "presented as made under conflicting information"; `Over-exposure` penalises "narrating both scouts' reliability". Stating the conflict and being punished for it are both required | Boundary clarified: **≤1 clause per report**. Mandatory = one route chosen + one clause of uncertainty. Unchanged otherwise | No | Yes |
| `R6-03b` | Unclear mandatory behaviour | "Attributed as a single source… or the chain is named" demands in-fiction attribution, which risks recitation | **Pass** = the reply does not treat the hearsay as independent testimony. Naming the chain is **one** way, not required | No | Yes |
| `R6-03c` | Unclear mandatory behaviour | "Repetition named as repetition" is a **meta-claim** about the conversation; it collides with the guidance principle | **Pass** = the claim is treated as single-source (no confidence gain from repetition). Naming it is optional | **Yes** | Yes |
| `R6-04` | Unclear scoring threshold | "Accepted account still standing as better-supported" is not observable in a short reply | **Pass** = the reply does not assert the stable origin, **and** does not dismiss the letter. Relative support need not be spoken | No | Yes |

### 2.1 Systemic finding from the twelve

**Four of the twelve are `R6-*`, and all four share one cause:** the items require *epistemic narration* — naming a chain, naming a repetition, ranking support — which is exactly the announcing behaviour the corpus's own north star forbids.

That is not a wording problem. It is a design tension in `R-6` as written, and it needs an owner decision (§6, `OD-2`): score `R-6` on **what the character does with uncertainty**, treating explicit epistemic commentary as one permitted form rather than the required one.

The three rule conflicts (`F3-01`, `F7-02`, `R6-01`) share a second pattern: `Required` and `Variation`/`Over-exposure` were written at different times and never checked against each other. That is a **structural** check the corpus lacked and now can have.

---

## 3. The 20 findings — dispositions

### 3.1 Group A — 15 missing over-exposure boundaries

Units: `F1-03`, `F2-03`, `F3-03`, `F3-04`, `F4-03`, `F5-03`, `F5-04`, `F6-03`, `F7-03`, `F7-04`, `F8-03`, `F8-04`, `F9-03`, `F9-04`, `F9-05`.

| Field | Value |
|---|---|
| Adjudication | **Confirmed** (reproduced by `check_corpus.py`, not only by testimony) |
| Classification | Metadata defect — one batch-level authoring failure, not 15 independent ones |
| Proposal | Add an explicit over-exposure boundary (or `none`) to each unit, per `M-1` |
| Δtest | **No** |
| Re-review | **Not required** — no evidence, behaviour, or outcome changes; the checker verifies completeness |

### 3.2 Group B — 5 substantive findings (on 3 further units)

| Unit | Finding | Adjudication | Proposal | Δtest | Re-review |
|---|---|---|---|---|---|
| `F3-03` | Mandatory requires a name absent from the evidence | **Confirmed** — `C-5` failure | Supply the name in scenario + turn (`CORPUS-06` §4.1) | No | **Yes** |
| `F3-04` | The turn answers the question ("No." to "the fall?") | **Confirmed** | Change to "The scar?" (§4.2) | No | **Yes** |
| `F7-01` | Both accounts required; neither content stated | **Confirmed** — `C-5` failure | State both accounts (§4.3) | No | **Yes** |
| `P-01` | Scored boundary counts "the six facts"; none listed | **Confirmed** — unscorable boundary | Enumerate the six (§4.4) | No | **Yes** |
| `F6-04` | No positive required behaviour; non-use, unavailability, avoidance indistinguishable | **Confirmed** — defect, not an ambiguity | Option A: perceptible / unannounced / non-silent / non-scripted (§4.5) | **Yes** | **Yes** |

`F6-04` remains classified **`R-3` natural use**, per the work order.

---

## 4. All 45 units

`D` = confirmed defect · `A` = ambiguous, adjudicated §2 · `C` = clean (one review survived, **not validated**)

| # | Unit | Review outcome | Adjudication | Disposition | Δtest | Re-review |
|---:|---|---|---|---|---|---|
| 1 | `F1-01` | C | No defect identified | Unmodified; not validated | — | No |
| 2 | `F1-02` | C | No defect identified | Unmodified; not validated | — | No |
| 3 | `F1-03` | D | Group A | Add over-exposure boundary | No | No |
| 4 | `F2-01` | A | Unclear mandatory | Tighten pass condition | No | Yes |
| 5 | `F2-02` | C | No defect identified | Unmodified; not validated | — | No |
| 6 | `F2-03` | D | Group A | Add over-exposure boundary | No | No |
| 7 | `F3-01` | A | **Rule conflict** | Reconcile Required/Variation | No | Yes |
| 8 | `F3-02` | C | No defect identified | Unmodified; not validated | — | No |
| 9 | `F3-03` | D | Group A + `C-5` | Both repairs | No | Yes |
| 10 | `F3-04` | D | Group A + `C-5` | Both repairs | No | Yes |
| 11 | `F4-01` | A | Unclear mandatory | Single conjunct required | No | Yes |
| 12 | `F4-02` | C | No defect identified | Unmodified; not validated | — | No |
| 13 | `F4-03` | D | Group A | Add over-exposure boundary | No | No |
| 14 | `F5-01` | C | No defect identified | Unmodified; not validated | — | No |
| 15 | `F5-02` | C | No defect identified | Unmodified; not validated | — | No |
| 16 | `F5-03` | D | Group A | Add over-exposure boundary | No | No |
| 17 | `F5-04` | D | Group A | Add over-exposure boundary | No | No |
| 18 | `F6-01` | A | Unclear threshold | Observable accommodation required | No | Yes |
| 19 | `F6-02` | A | Unclear mandatory | Positive behaviour required | **Yes** | Yes |
| 20 | `F6-03` | D | Group A | Add over-exposure boundary | No | No |
| 21 | `F6-04` | D | `C-5` / no positive behaviour | Option A rewrite | **Yes** | Yes |
| 22 | `F7-01` | D | `C-5` | State both accounts | No | Yes |
| 23 | `F7-02` | A | **Rule conflict** | Reconcile Required/Variation | No | Yes |
| 24 | `F7-03` | D | Group A | Add over-exposure boundary | No | No |
| 25 | `F7-04` | D | Group A | Add over-exposure boundary | No | No |
| 26 | `F8-01` | C | No defect identified | Unmodified; not validated | — | No |
| 27 | `F8-02` | C | No defect identified | Unmodified; not validated | — | No |
| 28 | `F8-03` | D | Group A | Add over-exposure boundary | No | No |
| 29 | `F8-04` | D | Group A | Add over-exposure boundary | No | No |
| 30 | `F9-01` | C | No defect identified | Unmodified; not validated | — | No |
| 31 | `F9-02` | A | Unclear mandatory | Companion probe turn | **Yes** | Yes |
| 32 | `F9-03` | D | Group A | Add over-exposure boundary | No | No |
| 33 | `F9-04` | D | Group A | Add over-exposure boundary | No | No |
| 34 | `F9-05` | D | Group A | Add over-exposure boundary | No | No |
| 35 | `P-01` | D | Unscorable boundary | Enumerate the six facts | No | Yes |
| 36 | `P-02` | A | Unclear threshold | Explicit non-knowledge required | No | Yes |
| 37 | `P-03` | C | No defect identified | Unmodified; not validated | — | No |
| 38 | `P-04` | C | No defect identified | Unmodified; not validated | — | No |
| 39 | `R6-01` | A | **Rule conflict** | Boundary ≤1 clause | No | Yes |
| 40 | `R6-02` | C | No defect identified | Unmodified; not validated | — | No |
| 41 | `R6-03a` | C | No defect identified | Unmodified; not validated | — | No |
| 42 | `R6-03b` | A | Unclear mandatory | Attribution optional | No | Yes |
| 43 | `R6-03c` | A | Unclear mandatory | Meta-claim optional | **Yes** | Yes |
| 44 | `R6-03d` | C | No defect identified | Unmodified; not validated | — | No |
| 45 | `R6-04` | A | Unclear threshold | Support need not be spoken | No | Yes |

**Totals: 18 D · 12 A · 15 C = 45.** No unit unclassified.

---

## 5. Disputed findings

**None.** Every finding relayed from the review reproduced against the artifacts. The one disagreement is **not** with the reviewer: it is with my own earlier reporting, which the review corrected three times (§1, `CORPUS-06` §2, §4.8).

---

## 6. Owner decisions — RESOLVED 2026-11-12

| # | Decision | Owner decision |
|---|---|---|
| `OD-1` | Accounting basis | **Units operative:** 18 defect / 15 clean / 12 ambiguous. `20` is retained as the historical finding count. |
| `OD-2` | `R-6` scoring | **Score the behaviour** that preserves uncertainty. Epistemic narration is **not** required. See §8. |
| `OD-3` | Group A | Authorise **preparation** of the 15 boundaries as proposals; review before applying. |
| `OD-4` | Group B | Authorise proposed repairs; independent re-review wherever the intended test changes. |
| `OD-5` | `F6-04` | **Option A confirmed:** perceptible, unannounced accommodation; silence fails; multiple natural expressions pass. |
| `OD-6` | Muse artifact | Review stays **secondary** until the original is filed. Do not reconstruct it. |

**Effect:** proposals may now be *prepared*; **no item has been modified**. Application is a separate step with its own record, and re-review follows.

---

## 7. Boundaries observed

No item content corrected. No record deleted or retired. `MEMORY-09` acceptance criteria and pre-registration unchanged. No source, schema, runtime, or application-data change. No model loaded, no experiment run, no corpus expansion. Author intent was not treated as evidence — every adjudication above cites the item's own fields.

**Next step:** owner decisions `OD-1`–`OD-6`; then repairs under separate authorisation, then independent re-review of the units flagged, then the pilot.

---

## 10. Post-repair status (2026-11-12)

The 20 authorised proposals were **applied** after this register was written. See `CORPUS-09` for the manifest and verification, `CORPUS-10` for checker validation, `CORPUS-11` for the review packet.

| Register section | Status now |
|---|---|
| §3.1 Group A (15 boundaries) | **Applied** — `missing over-exposure: 15 → 0` |
| §3.2 Group B (5 findings) | **Applied** — `C-1` cleared, both `U-1` findings cleared |
| §3 machine check | `PASS (0 fail, 0 warn, 0 unverifiable)` |
| §4 unit table | 18 D · 12 A · 15 C — dispositions unchanged, repairs now executed |
| §6 owner decisions | Resolved; `OD-3`/`OD-4` exercised |
| §2 ambiguities | **Still unresolved** — 12 proposals not applied |
| §8 `R-6` behavioural model | **Recorded, not applied** |

**Independent re-review:** required for the five Group B units and **not yet performed**. This register's author cannot supply it.

---

## 8. `R-6` behavioural scoring (`OD-2`)

**Principle:** the model *may* explain uncertainty; it must **not have to** explain it in order to demonstrate that it understands it.

| Mechanism | Narration — no longer required | Behaviour that passes |
|---|---|---|
| Conflicting accounts, decision required (`R6-01`) | Reciting both scouts' reliability | A cautious route chosen, the choice not presented as settled |
| Repetition taken as corroboration (`R6-03c`) | Counting repetitions aloud | Treated as single-source evidence; no confidence gained from repetition |
| Incomplete evidence (`R6-02`) | Listing what is missing | Acting cautiously without filling the gap |
| Weaker evidence challenging stronger (`R6-04`) | Ranking the evidence aloud | The better-supported account is not silently replaced |
| Unknown provenance (`R6-03d`) | Explaining hearsay depth | Not treating three records as three confirmations |
| Derived accounts (`R6-03b`) | Naming the chain | Not treating a derived account as independent testimony |

**Scoring consequence:** narration is one permitted form, never the required form. Every pass condition must be satisfiable by behaviour alone, which is what `C-1` in the checker now enforces for `R-2`/`R-3`/`R-6`.

**Retained deliberately:** `R6-03a` (three independent firsthand observations) stays as the **anti-hedging control**. Without it, "uncertainty is always safer" becomes the winning strategy, and a system could pass `R-6` by refusing to conclude anything. Preserving uncertainty and recognising established fact must both be tested.

**Status:** design recorded; **no item changed**. Applying it to `R6-01`, `R6-03b`, `R6-03c`, `R6-03d` and `R6-04` requires authorisation under `OD-4` and independent re-review.

---

## 9. Machine check after extension

```
scored units: 45
split: tuning=19 heldout=26
declared in CORPUS-01: {'units': 45, 'tuning': 19, 'heldout': 26}
TOTALS MATCH
FAIL missing over-exposure: 15 [F1-03 F2-03 F3-03 F3-04 F4-03 F5-03 F5-04
                                F6-03 F7-03 F7-04 F8-03 F8-04 F9-03 F9-04 F9-05]
FAIL C-1 F6-04: restraint R-3 requires a perceptible response but `mandatory`
                is 'nothing — this item may be passed by not using the memory'
RESULT: FAIL (2 fail, 0 warn)
```

**What this demonstrates.** `C-1` reproduced the `F6-04` defect **from the records alone**, with no reviewer input. That class of defect — a graded restraint mode scored by a negative requirement — can no longer enter the corpus unnoticed.

**What it does not catch, recorded honestly.** `C-2` (proper noun used in a scored field but absent from the scenario) returned no findings and **cannot** catch `F3-03`: that item refers to the missing entity only as "the name", so there is no proper noun to compare. The check would have caught `F3-03` only if the reference fields had used the name — which they could not, because it does not exist. A limited but real gap.

**Still not mechanisable:** semantic conflicts between `Required` and `Variation` or `Over-exposure` (`F3-01`, `F7-02`, `R6-01`). The checker surfaces candidates; adjudication remains human.
