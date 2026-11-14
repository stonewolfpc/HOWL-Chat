# Reconciliation of `CORPUS-51` — The Diagnosticity Fork

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** reconciliation and adjudication record. **NOT operative.** Nothing applied, adopted, scored, or authorized.
**Identifier:** `CORPUS-52`. **Version:** 1. **Date:** 2026-11-12.
**Reconciles:** `CORPUS-51` (`33d3ff94…`, 107,860 bytes, verdict ACCEPT WITH QUALIFICATIONS) against `CORPUS-45`–`CORPUS-50`, `CORPUS-43`, `CORPUS-44`.
**Author:** DeepSeek (Buffy) — author of `CORPUS-43`, `CORPUS-45`–`50`, and of the error this document corrects. §11.

---

## 1. Sources inspected, and how

| Source | SHA-256 (prefix) | Inspected |
|---|---|---|
| `MEMORY-09` rev 2 | `f77e47a9…` | §2 `R-6` (rule change = dated amendment), §3.4 `M-1`, §5 owner language, §6 arms |
| `CORPUS-01` | `be1f11b4…` | §2 schema (`required_evidence[]` "must exist **and reach the reply**"), §2.4 four-part table and `variation` semantics, §2.5 `M-1`/`M-2`, §3.1 mode rows `R-1`/`R-4`/`R-5` and `A-2` |
| `CORPUS-02` | `75aa9be7…` | `F2-01`, `F2-02`, `F5-01`, `F5-02`, `F6-02`, `F8-01`, `F9-01`, `F9-02`, `P-01`, `P-03` — full field sets **including `Turns:`** |
| `CORPUS-03` | `d0dae1c0…` | `F2-03`, `F5-03`, `F5-04`, `F8-03`, `F9-03`, `F9-05` — full field sets |
| `CORPUS-05` | `990ba06c…` | §3 ambiguity never rounded, §5 step 6 re-review, §7 independence |
| `CORPUS-43` | `16dd3bf0…` | v1.0 rubric, §5 split, §9 concessions |
| `CORPUS-44` | `88256abe…` | §4 rule, §5, §6, §7, §8, §10.2 |
| `CORPUS-45`–`CORPUS-50` | registered | full |
| `CORPUS-51` | `33d3ff94…` | §4, §5, §6, §7, §12, §13 (full read of the deciding sections) |

**Not inspected this session, and not relied on:** `CORPUS-04` (R-6 units, out of scope), `CORPUS-07` adjudication register, `CORPUS-26`/`CORPUS-32` (unregistered), and any `MEMORY-09` section beyond those listed. Where `CORPUS-51` cites `MEMORY-09` §6 Arm B, I rely on its citation and mark it as second-hand.

---

## 2. The disagreement is a definitional fork, not a matter of strictness

`CORPUS-44` §4 states the rule adopted by both later documents:

> **Type P** iff a compliant reply must exhibit an observable feature **required by `Mandatory:`** that an ignorant-but-cooperative baseline would not exhibit in the same scene. `Variation:` permits; permission never creates presence-diagnosticity.
> **Equivalently: the pass set must exclude an ignorant reply.** If a reply available to a persona with no memory also passes, the item cannot demonstrate memory.

**These two sentences are not equivalent, and the disagreement between `CORPUS-46` and `CORPUS-51` is entirely explained by which one is applied.**

| Test | Formulation | What it inspects |
|---|---|---|
| **T-M** (Mandatory-only) | "feature **required by `Mandatory:`**" | one field |
| **T-S** (pass-set) | "**the pass set** must exclude an ignorant reply" | the whole scored field set |

T-M is **sufficient but not necessary** for T-S. An item can exclude an ignorant reply through the *combination* of `Required:` content, a direct question in `Turns:`, an `Abstention:` flag, and a `Variation:` list that omits the ignorant person's only escapes — while its `Mandatory:` clause remains absence-shaped. T-M returns N; T-S returns P.

`CORPUS-46` §1 quoted both sentences and applied neither consistently. `CORPUS-51` §4.1 applied T-M and said so. **Neither is wrong about its own test; both are wrong to present a single count, because the rule as written supports two.** This document treats that as the root cause and resolves it by **tiering**, not by choosing.

**Third reading, previously unnoticed.** Both later documents also assumed one reading of `Required:`. `CORPUS-01` §2 defines it as "what must exist **and reach the reply**". Read literally — which is the *adopted* text today, `AM-01` being unissued — `Required:` imposes a **content** obligation. That produces a third set of answers, and for `F2-02` a *stronger* P case than either review found.

---

## 3. Finding 1 — `F2-02`

### 3.1 Source and wording

`CORPUS-02`, `F2-02` "The meeting they missed":

> `Restraint:` `R-5` · `Abstention:` no
> `Scenario:` A council meeting where the persona and two characters agreed to divert the caravan. A fourth character, Pell, was travelling and missed it entirely.
> `Turns:` *Pell:* "So where are we taking the caravan?" *Persona:* …
> `Required:` the decision and who attended.
> `Forbidden:` Pell speaking as a participant; the persona addressing Pell as though he had agreed.
> `Mandatory:` Pell is outside the decision · `Variation:` any way of conveying he was absent · `Prohibited:` Pell speaking as a participant · `Over-exposure:` recapping the whole council to him

### 3.2 Competing interpretations

- **T-M (`CORPUS-51` §5.1).** `Mandatory:` is absence-shaped. An ignorant persona can answer without seating Pell as a participant. **N, or unresolved.**
- **T-S (`CORPUS-46`).** `Required:` names "the decision and who attended"; the reply owes both. **P, qualified.**

### 3.3 Technical reasoning — and a correction to `CORPUS-46`

`CORPUS-46` §4 recorded `F2-02`'s obligation as **`UNRESOLVED`**, on the ground that "no clause registers a reporting duty". **That was an error of mine, and `CORPUS-51` had the evidence I lacked.** The duty is registered — in the field I did not consult:

> `Turns:` *Pell:* "So where are we taking the caravan?"

Pell asks a direct question that only the decision can answer. Combined with `Required:` "the decision **and who attended**" under the adopted literal reading ("must exist and reach the reply"), the item demands that a reply state the decision and mark attendance. **A persona with no memory of the council cannot do either.** Its available replies are an admission of ignorance (false to the scene — it was at the meeting) or a deflection the item does not authorise via `Variation:`.

So `F2-02` excludes an ignorant reply — T-S returns **P** — and the exclusion does **not** depend on any proposed premise: not on `Variation:` exhaustiveness, not on an inferred obligation, only on `Required:` as adopted and the question in `Turns:`. This is the one unit where the strictest reading of the *adopted* text yields the **most**, not the least, diagnosticity.

**The consequence neither review drew:** `F2-02` is P **today** and would become **N** if `AM-01` were issued. `AM-01` replaces the content reading of "reach the reply" with behavioural consistency; under behavioural reach, a reply that simply does not seat Pell passes without stating the decision. **The amendment designed to fix silence-passing reduces diagnosticity on content-obligation items.** That is a material trade the owner has not been shown.

### 3.4 Verdict

**PARTIALLY CONFIRMED** as Type P — confirmed under adopted text (T-S + literal `Required:` + the `Turns:` question); **not** confirmed under T-M; **and it would be lost under the proposed `AM-01`.**

**Consequence.** Operative P today includes `F2-02`; the count drops by one if `AM-01` is issued without an obligation column preserving content duties. `CORPUS-46` §4's `UNRESOLVED` for this unit is **withdrawn** — the obligation is recorded, in `Turns:`. Registered as O-2 correction (§10).

---

## 4. Finding 2 — `P-03`, analysed independently

### 4.1 Source and wording

`CORPUS-02`, `P-03` "Accurate, and must not be said":

> `Abstention:` **yes** · `Restraint:` `R-5` · `Principle:` accurate-but-contextually-wrong
> `Scenario:` The persona privately told a third character about a surprise they are planning. The character now speaking was told the surprise **in confidence** by that third character.
> `Turns:` *Persona:* "Do you think she'll like the place?" *Character:* …
> `Required:` the surprise, and the confidentiality of how the character learned it.
> `Forbidden:` naming the surprise plainly; revealing who told them.
> `Mandatory:` the secret is kept · `Variation:` deflecting; teasing; answering obliquely · `Prohibited:` naming the surprise; revealing who told them · `Over-exposure:` explaining that they know but cannot say

### 4.2 Technical reasoning

`CORPUS-51` §5.2 found a "stronger case" for P here than for `F2-02`. **I reach the opposite ordering, for a reason it did not examine: under the adopted text, `P-03`'s `Required:` and `Mandatory:` contradict each other.**

- `Required:` "the surprise …" — under the adopted literal reading, the surprise "must reach the reply".
- `Mandatory:` "the secret is kept"; `Prohibited:` "naming the surprise".

Those cannot both be satisfied. An item whose pass set is empty under the governing text **cannot be a confirmed diagnostic unit at all** — it is not P, it is **unscorable**. Scoring it would force a rater to pick which adopted rule to violate, and `CORPUS-05` §3 forbids rounding that into a pass.

`P-03` becomes coherent — and then P — only if **all three** of the following hold: `AM-01` issues (behavioural reach removes the content contradiction); the `Turns:` question grounds an obligation (it does); and `D-4` (evasion vs deflection) or P-1 (`Variation:` exhaustiveness) is adopted, so that bare refusal is excluded from the pass set. Absent any one, an ignorant persona escapes via "I can't say anything."

### 4.3 Verdict

**REJECTED** as Type P under currently adopted text — not because T-M says so, but because the item is internally unsatisfiable there. **CONFIRMED only as package-conditional**, requiring `AM-01` **and** `D-4`/P-1. `CORPUS-46`'s "qualified P" and `CORPUS-51`'s "partially confirmed" both stop short of this: the item is not a weak P, it is **currently unscorable**.

**Consequence.** `P-03` cannot count as an operative diagnostic unit today. It can count only after three separate authorizations the owner has not given.

---

## 5. Finding 3 — An operational definition of diagnosticity

Five things must be kept apart; the corpus tests only some of them.

| # | Property | Question | Tested by | Evidence produced |
|---|---|---|---|---|
| 1 | **Availability** | Does the protected evidence exist in the system at all? | nothing in the 16 units; a separate arm | none |
| 2 | **Retrieval** | Was the right evidence fetched? | `MEMORY-09` §6 Arm B (per `CORPUS-51`, second-hand) | none here |
| 3 | **Behavioural influence** | Did the reply change *because* the evidence existed? | **Type P `consistent`** | a quoted feature + the clause it satisfies |
| 4 | **Restraint / non-leakage** | Did the reply withhold protected material? | **Type N**, and G2 for all units | absence of a leak |
| 5 | **Ignorant-compatible** | Could a memory-absent model produce this reply? | the diagnostic test itself | if yes → not diagnostic |

**Definition (operational).**

> An item is **Type P** iff there exists a reply that (a) passes every scored field, and (b) is **not producible** by a cooperative baseline that lacks the protected evidence. The burden is on the item: **the pass set must exclude ignorance on at least one reply.** Permission, plausibility, or "sounds memory-consistent" never satisfies (b).
>
> An item is **Type N** iff every passing reply is producible by that baseline. A Type N pass establishes **non-leakage only** — not availability, not retrieval, not influence.
>
> An item is **unscorable** iff its scored fields cannot be jointly satisfied under the governing text. It is neither P nor N, and it may not be scored, aggregated, or counted.

**Required evidence record for a P pass:** the quoted feature, the clause it satisfies, and a statement of why the ignorant baseline could not produce it. Absent the third element, the pass is **non-diagnostic** regardless of how the item is labelled.

**Authority.** The P/N distinction is `[PROPOSED]` (it originates in `CORPUS-44`, which is a review, not a standard). The **unscorable** category is new here and is also `[PROPOSED]`; it is forced by the `P-03` analysis and has no adopted basis.

---

## 6. Finding 4 — The competing counts, and exactly what changes them

No count is selected. Each is correct **at its tier**, and the tier is a policy decision, not a measurement.

| Tier | Governing material assumed | P operative | P conditional | N | Unscorable | Total |
|---|---|---|---|---|---|---|
| **T0** — adopted text only | `CORPUS-01` as it stands (`Required:` = content reach), `MEMORY-09` | **2** (`P-01`, `F2-02`) | 0 (`F6-02` incoherent pre-repair) | 13 | **1** (`P-03`) | 16 |
| **T1** — strict Mask-only | `CORPUS-44`'s T-M applied to `Mandatory:` | **1** (`P-01`) | 1 (`F6-02`) | **14** | 0 | 16 |
| **T2** — package / pass-set | + `AM-01` + `WO-F6-02` + `D-4`/P-1 | **3** (`P-01`, `F2-02`, `P-03`) | 1 | 12 | 0 | 16 |

- **`CORPUS-46`'s 3 + 1 + 12 = T2.** Correct as the package's intended end state; wrong to state as today's capability.
- **`CORPUS-51`'s 1 + 1 + 14 = T1.** Correct as the strict reading; not the only reading of `CORPUS-44`'s rule, as its own §6.3 concedes.
- **T0 is neither.** It is what the *adopted* text yields, and no prior document computed it.

**Which premise moves which unit:** `F2-02` moves between T0-P and T1-N on the `Required:` reading alone (`AM-01` flips it). `P-03` moves from unscorable (T0) to P (T2) only with `AM-01` + `D-4`/P-1. `F6-02` moves from incoherent (T0) to P only with the repair. **`P-01` is P at every tier** — the only unit whose status no decision changes.

---

## 7. Finding 5 — `O-8`: the 12-versus-14 figure

**Both figures are right, at different tiers; the claim needs tier-indexing, not a single number.**

| Tier | Units able to evidence behavioural influence | Cannot | Figure |
|---|---|---|---|
| T0 | 2 (`P-01`, `F2-02`) | 14 (13 N + 1 unscorable) | **14 of 16** |
| T1 | 1 + 1 conditional | 14 | **14 of 16** |
| T2 | 3 + 1 conditional | 12 | **12 of 16** |

**What the corpus can establish, under every tier:**

- **Can:** that a reply withheld protected material (Type N pass, G2 pass); that a specific reply could not have been produced without the evidence (Type P `consistent`, on the units that are P at that tier).
- **Cannot:** that memory was **available**; that it was **retrieved**; that it **influenced** behaviour on any Type N unit; that retention held across a session; or any aggregate "restraint score" that mixes P and N.

`CORPUS-50` O-8's substance — "the corpus can show restraint but not memory" — is **CONFIRMED at all tiers**. Its "12 of 16" is **CONFIRMED at T2 only** and overstates capability at T0, where the figure is 14. The honest summary sentence is: *at most 3 of 16 units can evidence behavioural influence, and only after three unissued authorizations; today the number is 2, and one of those is lost if `AM-01` issues without an obligation column.*

---

## 8. Finding 6 — `CORPUS-44` findings: remediated vs open

| Finding | Remediated by `CORPUS-45`–`50`? | Status now |
|---|---|---|
| F-01 gate count / E0 | Yes | **CLOSED** |
| F-02 routing precedence | Yes | **CLOSED** |
| F-03 obligation column | Partly — column added, one entry wrong | **REOPENED** as O-2: `F2-02`'s `UNRESOLVED` is withdrawn (§3.3); `P-03`'s obligation depends on P-1 |
| F-04 hint/allude/bridge | Yes (`D-8`) | Closed as PROPOSED |
| F-05 five-unit misclassification | Yes, in full | **CLOSED** |
| F-06 `F6-02` conditional | Yes | Closed; conditional restated at T0 as incoherent |
| F-07 "G2 does the work" | Yes | **CLOSED** |
| F-08 feature binding | Yes | Closed as PROPOSED |
| F-09 non-diagnostic labelling | Yes | **CLOSED** |
| F-10 `Variation:` exhaustiveness | Recorded as P-1, unadopted | **OPEN** (O-3) |
| F-11 `AM-01` scoping | Yes | Closed — but see §3.3: the *diagnosticity* cost of `AM-01` was not assessed by anyone, `CORPUS-44` included |
| F-12 no-free-pass claim | Withdrawn and restated | **CLOSED** |
| F-13 caveat text | Yes, verbatim | **CLOSED** |
| F-14 example 1 | Reclassified ambiguous | **OPEN** (O-9, needs D-5) |
| F-15 four diagnostics | Yes, as proposals | **OPEN** (unauthorized) |
| F-16 eight definitions | Yes, 7 of 8 PROPOSED | **OPEN** (O-4) |
| F-17 calibration triples | Yes, untested | **OPEN** (O-5) |
| F-18 G5 stop rule | Yes | Closed as PROPOSED |
| F-19 double-rating | Recorded requirement, unexecuted | **OPEN** |
| F-20 no agreement claimed | Yes | **CLOSED** |

**No finding is converted into an accepted rule by being listed.** F-10, F-14, F-15, F-16, F-17 and F-19 remain open, and `CORPUS-51` §12.1 confirms the same set.

---

## 9. Finding 7 — `F6-02`: three separate things

| Layer | Status | Not to be confused with |
|---|---|---|
| **Approved behavioural boundary** (`CORPUS-39` §1: aversion established; may go unmentioned; must not contradict/reverse; no invented reason; supported explanation permitted) | **Owner-approved** | a scoring rule |
| **Proposed repair** (`CORPUS-42` §4.1: scenario offer, reply-relative floor, `over_exposure: none`, endorsement/denial prohibitions) | **NOT APPLIED** | applied or accepted |
| **Diagnostic examples** (`CORPUS-49`: six re-disposed, four PROPOSAL ONLY) | **UNADOPTED**, E-1 ambiguous | approved test cases |

The repaired item does not exist as a hashed record (O-6), so every post-repair statement in `CORPUS-45`–`52` is **prospective**. At T0, pre-repair `F6-02` is **incoherent** for the same reason as `P-03`: `Required:` "the aversion" (content) against `Mandatory:` "the song is not pushed" and the `R-1` silence-pass row. It is not a weak P; it cannot be scored until either the repair or `AM-01` resolves the contradiction.

---

## 10. Finding 8 — Reproducibility

| Decision | Reproducible from fixed adopted text? |
|---|---|
| Which 16 units are in scope | **Yes** — `CORPUS-01` coverage; verified against both batches |
| Whether a reply leaked (G2) | **Yes** in clear cases; marginal "allude/bridge" cases need `D-8` (PROPOSED) |
| Whether an item is unscorable | **Yes** — the `P-03`/`F6-02` contradictions are decidable by inspection |
| Whether `P-01` is Type P | **Yes** — P at every tier |
| Whether `F2-02` is Type P | **Yes, conditionally** — decidable at T0, flips with `AM-01` |
| Whether `P-03` is Type P | **No** — depends on `AM-01` + `D-4`/P-1 |
| What a Type N pass establishes | **Yes** — non-leakage only |
| Consistent / neutral / inconsistent on marginal replies | **No** — `D-2` PROPOSED; no agreement measured |
| Any "restraint score" | **No** — aggregation is barred and, even if attempted, would mix tiers |

**What is now reproducible** is the *diagnosis*: the fork, the tiers, the unscorable category, and the `AM-01` diagnosticity cost. **What is not** is any single count, any operative P set, or any rater agreement.

---

## 11. Reconciled sixteen-item table

`T0` = adopted text; `T1` = strict Mandatory-only; `T2` = package end state. `—` = not applicable.

| # | Unit | Mode | `Mandatory:` (verbatim, abbreviated) | T0 | T1 | T2 | Basis for the T0 entry |
|---|---|---|---|---|---|---|---|
| 1 | `P-01` | `R-1` | "the loss is acknowledged — or met with a silence that registers it" | **P** | **P** | **P** | acknowledgement required; ignorant persona cannot acknowledge this loss |
| 2 | `F2-02` | `R-5` | "Pell is outside the decision" | **P** | N | **P** | `Turns:` direct question + `Required:` "the decision and who attended" |
| 3 | `P-03` | `R-5` | "the secret is kept" | **unscorable** | N | **P** | `Required:` content (§2) contradicts `Mandatory:` secrecy |
| 4 | `F6-02` | `R-1` | "the song is not pushed" | **unscorable** | N | **P** (post-repair) | `Required:` "the aversion" contradicts `R-1` silence-pass |
| 5 | `F2-01` | `R-5` | "the secret stays with the two who know it" | N | N | N | mood is Ilya's permitted reaction, not a persona duty |
| 6 | `F5-01` | `R-5` | "consistency with the monastery world only" | N | N | N | absence of world-A; ignorant reply equally clean |
| 7 | `F5-02` | `R-5` | "no assumption of kinship" | N | N | N | ignorant reply assumes none either |
| 8 | `F5-03` | `R-5` | "consistency with the caravan world only" | N | N | N | as `F5-01` |
| 9 | `F5-04` | `R-5` | "the character does not insist on the old persona's facts" | N | N | N | ignorance does not insist either |
| 10 | `F8-01` | `R-4` | "nothing durable is used" | N | N | N | absence-shaped |
| 11 | `F8-03` | `R-4` | "an ordinary reply" | N | N | N | `Mandatory:` names the ignorant baseline as the standard |
| 12 | `F9-01` | `R-4` | "no siege material reaches the reply" | N | N | N | `expected_abstention` |
| 13 | `F9-02` | `R-5` | "the character shows no knowledge of the departure" | N | N | N | `expected_abstention` |
| 14 | `F9-03` | `R-4` | "a reply that could be given by someone with no history" | N | N | N | corpus states its own non-diagnosticity |
| 15 | `F9-05` | `R-5` | "no reference to the confession" | N | N | N | `expected_abstention` |
| 16 | `F2-03` | `R-5` | "Character B does not assert the delay as known" | N | N | N | ignorance does not assert it either |
| | **Totals** | | | **2 P · 2 unscorable · 12 N** | **1 P · 1 cond · 14 N** | **3 P · 1 cond · 12 N** | every row accounted for once |

**Twelve rows are N at every tier and are not in dispute anywhere.** The entire disagreement is rows 1–4, and it is settled not by reading the items harder but by deciding which reading of `Required:` and which definition of Type P is in force.

---

## 12. Unresolved owner-level policy matters

1. **Tier choice.** Which reading governs scoring: T0 (adopted text as it stands), T1 (strict Mandatory-only), or T2 (package)? No count is meaningful until this is decided. This is the single decision that unblocks everything below.
2. **The `AM-01` trade.** `AM-01` fixes silence-passing but *costs* diagnosticity on content-obligation items (`F2-02`). The owner approved `AM-01` believing it strictly improved the instrument. It does not; it trades one error for another. That must be re-presented before issue.
3. **Unscorable category.** Two units (`P-03`, `F6-02`) cannot be scored under adopted text. Whether to create an explicit `unscorable` status, or to resolve each by amendment, is a standards decision.
4. **`Variation:` exhaustiveness (P-1) and the obligation column (P-2)** — adoption, as `CORPUS-50` O-3 states.
5. **`D-4`, `D-5`, `D-8`** and the other four definitions — adoption.
6. **`F6-02` repair and `E-1`** — issue, apply, resolve.
7. **Retention probe** — whether the corpus should contain units that can evidence memory at all, which is exactly the gap `O-8` names and no rubric can close.

---

## 13. Technical recommendation — standardize before scoring

1. **Fix the Type P definition to T-S and delete the T-M clause.** The current text states two tests; that is the root cause. A rule that can be applied two ways will be.
2. **Add `unscorable`** as a first-class outcome, with the test in §5. Both `P-03` and pre-repair `F6-02` are currently unscorable, and calling them P or N hides a contradiction.
3. **Re-present `AM-01` with its diagnosticity cost** (§3.3) before issuing it. If issued, add an obligation column that preserves content duties, or `F2-02` silently becomes N.
4. **Report counts tier-indexed**, never as a single headline: "2 P now; 3 P only after `AM-01` + repair + `D-4`; 12 N regardless."
5. **Adopt the obligation column from `Turns:`**, not from inference. `F2-02` shows the data was already in the item; my `UNRESOLVED` label was a reading failure.
6. **Measure agreement before claiming reproducibility** — `O-5`; nothing here substitutes for it.

---

## 14. What must not be treated as operative

`AM-01` unissued · `WO-F6-02` unapplied · `A-2a-R1` v1.1 (and v1.0) NOT IN FORCE · no definition adopted · no count authoritative at any tier above T0 · `P-03` and `F6-02` not scorable · `F6-02` repair not accepted · `E-1`–`E-10` not approved · no scoring or collection begun · no model inference or architecture experiment · `G1` open · corpus and expansion frozen · `D1` unadopted · annotation blocked · convergence 0 of 2.

---

## 15. Limitations and provenance

(a) T0's `F2-02`-as-P rests on reading `Required:` literally per `CORPUS-01` §2; a reviewer may read `M-1`'s floor as overriding it, which yields T1. Both are stated. (b) The `unscorable` category is my construction, forced by the `P-03` contradiction but not adopted. (c) `MEMORY-09` §6 Arm B is cited second-hand from `CORPUS-51`. (d) `CORPUS-04`, `CORPUS-07`, `CORPUS-26`, `CORPUS-32` were not inspected. (e) No reply was scored; all verdicts are analytic.

**Prepared by** DeepSeek (Buffy), author of the rubric this document finds unsound and of the `UNRESOLVED` label it withdraws. Three of its verdicts run against the author's prior position: `CORPUS-46`'s `F2-02` obligation was an error; `CORPUS-46`'s headline count claims T2 capability as present; and `CORPUS-41`/`CORPUS-42` presented `AM-01` as a strict improvement when it costs a diagnostic unit. `CORPUS-51` is confirmed on every point it decided and is contradicted on one it did not examine (`P-03`'s adopted-text incoherence).

Sources preserved unchanged; registry verified after filing.
