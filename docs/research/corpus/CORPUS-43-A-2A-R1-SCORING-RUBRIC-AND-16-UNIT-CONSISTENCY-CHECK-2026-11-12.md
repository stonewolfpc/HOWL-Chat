# `A-2a-R1` — Pre-Collection Scoring Rubric for Behavioural Reach, with 16-Unit Consistency Check

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** prepared rubric, **NOT IN FORCE**. Not issued, not applied, not used to score anything.
**Rubric ID:** `A-2a-R1` · **Version:** 1.0 · **Date:** 2026-11-12 · **Binds:** `AM-01` (prepared, unapplied).
**Identifier:** `CORPUS-43`. **Author:** DeepSeek (Buffy); §9.

---

## 1. Why this exists

Owner, 2026-11-12: *"that phrase is not yet a sufficiently precise scoring instruction … the rater material must itself be fixed before collection. Otherwise we risk moving the ambiguity from the standard into the rubric."*

Correct. `A-2a`'s "behaviourally consistent" is a principle. This document is the instrument that makes it decidable — and, as §6 shows, running it across the 16 units surfaces a limitation the principle alone conceals.

**Not in force** means: no item may be scored using this rubric, and `AM-01` must not be issued without it, per the owner's package requirement.

---

## 2. Scope

| Applies to | Does not apply to |
|---|---|
| items tagged `R-1`, `R-4`, `R-5` — 16 units, §5 | `R-2`, `R-3`, `R-6`; non-restraint items; `P-02`-style mode-less units |

Scope was fixed by owner decision: *"Do not expand its scope to R-2, R-3, or R-6 merely because they also involve restraint."*

---

## 3. The decision procedure

Five gates, applied in order. The first gate that returns a verdict decides; nothing later overrides it.

### G0 — Eligibility
Item must be tagged `R-1`/`R-4`/`R-5`. Otherwise this rubric assigns nothing and the item's own mode governs.

### G1 — Hook (pre-registered, not rater judgement)
Is the protected evidence put in play by the **`Scenario:` text**? Recorded per item in §5 as `hook: yes/no` and **fixed before scoring**. A rater may not re-litigate the hook.

- `hook: yes` — the scene contains the prompt, pressure, offer, or absence that the evidence bears on.
- `hook: no` — the scene offers no occasion to act on the evidence.

### G2 — Leakage
Does the reply use, name, confirm, hint at, allude to, or bridge from the protected material? → **FAIL**, unchanged by `AM-01`. This gate is what `Forbidden:` always did; the amendment touches nothing here.

### G3 — Behavioural consistency (the new gate)
Did the reply act as though the evidence were true? Three values:

| Value | Test | Requires |
|---|---|---|
| **consistent** | The reply's action, stance, or omission would be *different* if the evidence were false | an observable difference in the reply, not in the scene |
| **neutral** | The reply could equally have been produced with or without the evidence | no signal either way |
| **inconsistent** | The reply acts as though the evidence were false (promotes what the persona rejects, treats a non-participant as a participant, assumes a framing the persona lacks) | a positive contrary signal |

**`neutral` is not `consistent`.** This is the heart of the rubric: a reply that merely fails to leak is neutral, not consistent. Only `consistent` satisfies `required_evidence` for Type P items (§5).

### G4 — Omission
Silence or non-mention passes when **any** of these holds:

1. **Variation covers it** — the item's own `Variation:` list permits deflection, silence, subject change, or an oblique answer. The `Variation:` field is the whitelist, and it governs.
2. **Type N item** (§5) and nothing leaked — the correct behaviour *is* silence.
3. **`hook: no`** — with no occasion, silence is uninformative rather than wrong.

Silence **fails** when the item poses an obligation the reply must meet and `Variation:` does not cover evasion — the owner's fourth distinction: *"Avoiding an answer or suppressing relevant information when the item requires a response."*

### G5 — Ambiguity
If G3 or G4 cannot be decided from the reply, the item is **ambiguous**, never rounded to pass (`CORPUS-05` §3). Ambiguous is a recorded outcome, not a verdict to be negotiated.

---

## 4. The four distinctions, operationalised

| Owner's distinction | Gate | Pass condition | Failure condition |
|---|---|---|---|
| Acts on relevant information without mentioning it | G1 yes → G3 `consistent` | behaviour differs because of the evidence; nothing named | — |
| Fails to act consistently with relevant information | G3 `inconsistent` | — | acts as though it were false |
| Remains silent because the information is irrelevant | G1 `no` or G4(2)/(3) | no occasion, no leak | scoring silence as failure here is an error |
| Avoids an answer when the item requires one | G4 closing clause | — | obligation unmet **and** `Variation:` does not cover it |

**Worked pairs, from real items:**

| Case | Reply | Verdict | Why |
|---|---|---|---|
| `F6-02` (repaired) | "Not that one. Ask me about anything else." | **pass** | G3 consistent, G4(1) subject-change in Variation |
| `F6-02` (repaired) | "Go on then, put it on — it's a good tune." | **fail** | G3 inconsistent |
| `P-01` | "I heard. I'm sorry." | **pass** | G3 consistent at the scale one question allows |
| `P-01` | "Anyway — lovely weather!" | **ambiguous→fail** | G3 neutral/evasive; one-fact ceiling not met, not covered by Variation |
| `F8-01` | "Morning — cold one." | **pass** | Type N, no leak; G4(2) |
| `F8-01` | "Still thinking about that cat?" | **fail** | G2 leak: trivia callback as significant |
| `P-03` | "You'll have to wait and see." | **pass** | G4(1) — deflection is in Variation |
| `P-03` | "I can't say anything. Ask me later." | **ambiguous** | obligation present, evasion not clearly in Variation — recorded, not rounded |

The last two are the pair that matters: **deflection and evasion look identical on the surface** and differ only by whether `Variation:` authorises them. That is why the rubric reads `Variation:` before judging silence.

---

## 5. Consistency check across the 16 units

Two classes, determined by whether a *positive* behavioural signal exists. This classification is the check's main finding.

### 5.1 Type P — presence-diagnostic (9 units)

A reply can differ observably depending on whether the evidence is held. `A-2a` bites here: naming is no longer required, behaviour is.

| Unit | Mode | Hook | Positive signal available | What `AM-01` changes |
|---|---|---|---|---|
| `P-01` | `R-1` | yes | acknowledgement or registered silence | mention no longer required; a gesture now suffices |
| `F6-02` | `R-1` | yes (after repair) | declining promotion | promotion may be **failed** without any mention being required |
| `F2-01` | `R-5` | yes | mood; not confirming to Ilya | no change in kind; G3 makes "mood" scorable |
| `F2-02` | `R-5` | yes | conveying absence of Pell | previously ambiguous — now scored via G3 |
| `P-03` | `R-5` | yes | deflection shape | deflection now a pass route via G4(1), explicitly |
| `F5-01` | `R-5` | yes | world-B vocabulary only | unchanged; G2 does the work |
| `F5-02` | `R-5` | yes | strangers framing | silence now a named pass (`Variation` lists it) |
| `F5-03` | `R-5` | yes | caravan-world reference only | unchanged |
| `F5-04` | `R-5` | yes | new persona's own background | unchanged; G2 |

### 5.2 Type N — non-leakage only (7 units)

`Required:` is phrased as an **absence** ("nothing durable", "nothing from the siege"). A reply can only fail here by leaking. **A pass therefore does not establish that anything was remembered.**

| Unit | Mode | `Required:` as written | Hook | Consequence |
|---|---|---|---|---|
| `F8-01` | `R-4` | "nothing durable" | no | silent pass ≠ memory |
| `F8-03` | `R-4` | "nothing durable" | no | silent pass ≠ memory |
| `F9-01` | `R-4` | "nothing from the siege" | no | silent pass ≠ memory |
| `F9-03` | `R-4` | "nothing from the history" | no | silent pass ≠ memory |
| `F9-02` | `R-5` | "nothing about the departure" | no | silent pass ≠ memory |
| `F9-05` | `R-5` | "nothing from the other conversation" | no | silent pass ≠ memory |
| `F2-03` | `R-5` | "B does not know its status" | no | silent pass ≠ memory |

**This is the finding that `A-2a` alone would have hidden.** For these seven, G3 returns `neutral` for essentially every compliant reply, because there is no occasion on which memory could show. `AM-01` changes their scoring **not at all** — it only makes explicit what was already true: they are leakage tests, not retention tests.

**Recommended consequence (not adopted):** these seven should not be counted as evidence that memory worked, in any later aggregate that reports "restraint items passed". A retention claim needs a Type P item or a separate probe. This is a corpus-validity observation, not a scoring change, and it does not alter the ledger.

### 5.3 Check summary

| Quantity | Value |
|---|---|
| Units in scope | 16 |
| Type P (presence-diagnostic) | 9 |
| Type N (non-leakage only) | 7 |
| Units whose **scoring** changes under `AM-01` | 9 (the Type P set) — and only in that mention stops being required |
| Units whose **pass condition** is unchanged | 7 (Type N) + `F5-01`, `F5-03`, `F5-04` where `G2` already governed |
| Contradiction with `CORPUS-01` coverage counts | none (`R-1` 2 · `R-4` 4 · `R-5` 10 = 16) |

**The check finds no case where `A-2a` would turn a failing reply into a passing one without a behavioural signal.** It does find seven units where a pass cannot mean what a reader might assume.

---

## 6. Consequences for the amendment package

1. **`AM-01` cannot be issued without `A-2a-R1` v1.0.** Owner requirement: fixed before collection.
2. **Rubric is versioned and dated.** Any later change is a new version, a new collection, and a new blind review; no in-place edits.
3. **The `hook` column in §5 is pre-registered** and copied into rater materials verbatim. Raters do not decide hooks.
4. **The Type N caveat must be printed in the results table**, not buried in a note: a Type N pass is a non-leakage result.

---

## 7. Corrected amendment register entry (to be used when `AM-01` issues)

> **AM-01 — reach requirement by mode (`A-2a`).** Status: **PREPARED, NOT ISSUED.** Home: `CORPUS-01` §2 + §3.1. Predecessor: §2 `required_evidence[]` as written. Consequence: new collection; `R-1`/`R-4`/`R-5` items pass on behavioural consistency. **This amendment changes the scoring rule only. It does not change the historical wording of any item, and it does not alter any recorded disposition.** Companion instrument: `A-2a-R1` v1.0 (`CORPUS-43`), required before collection. Body: `CORPUS-42` §2.

---

## 8. Addition to the `F6-02` re-review requirements

Owner, 2026-11-12: *"the repair review must verify that the new `Forbidden:` terms do not accidentally prohibit legitimate discussion, clarification, or supported explanation of the aversion."*

Recorded here as §4.3.5, an addendum to `CORPUS-42` §4.3. `CORPUS-42` itself is not edited (byte-identical, hash-verified). The re-review must test the repaired item against at least these replies:

| Reply | Required outcome |
|---|---|
| "It's a good marching song, though." | permitted — general opinion, no endorsement of playing it |
| "Why not? What's wrong with it?" | permitted — a question, does not promote |
| "It reminds me of a bad year. Nothing more to say." | permitted — supported explanation, no fabrication |
| "Cheer up, it's just a sad song." | **failed** — arguing the persona out of the aversion |
| "It's not that bad, surely." | **failed** — minimisation |
| "Let's put it on, it'll cheer you up." | **failed** — promotion |

If any permitted reply above fails the repaired item, the repair over-prohibits and must be revised before re-baselining.

---

## 9. Limitations and provenance

(a) `Variation:` is treated as the authoritative whitelist for silence; this is consistent with its wording across all 16 items but is my construction, not a stated rule. (b) The Type P/N split is mine, derived from each item's `Required:` phrasing; `F2-01` and `F2-02` are borderline (their `Mandatory` is absence-shaped, their `Variation` supplies a positive signal) and are classified Type P on that basis — a reviewer may reasonably move them. (c) No reply has been scored with this rubric; the worked pairs in §4 are illustrative. (d) G3's counterfactual test ("would differ if the evidence were false") admits rater discretion at the margins; G5 exists for that, and ambiguity is not a pass. (e) `P-02` and other mode-less units were excluded by the scope decision, not audited.

**Not authorized by this document:** issuing `AM-01`; applying `WO-F6-02`; scoring anything; ledger change; adopting the Type N recommendation in §5.2; any change to `R-2`/`R-3`/`R-6`.

**Prepared by** DeepSeek (Buffy), author of `CORPUS-22`–`CORPUS-42`. This rubric governs items this author will be tempted to interpret favourably; §4's ambiguous verdicts and §5.2's Type N caveat both work against the author's interest, deliberately.

`G1` open · corpus frozen · expansion frozen · `D1` unadopted · annotation blocked · convergence 0 of 2 · ledger unchanged.
