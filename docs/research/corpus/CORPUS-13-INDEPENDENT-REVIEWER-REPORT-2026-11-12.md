# Independent Reviewer Report — Group B Post-Repair Review

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Status:** review **performed**. Five items reviewed; one item re-reviewed after a briefing correction. Every verdict below is preserved **verbatim**.
**Date:** 2026-11-12
**Packet:** `CORPUS-12` v2 (commissioned instrument), hash in §6.
**Authority:** work order `CORPUS-11` §2–§4; `CORPUS-05` §3–§5; `CORPUS-07` `OD-4`.
**Effect on `G1`:** none. `G1` remains **open**. Nothing here seals the corpus, accepts an item, or authorises `G2`.

---

## 1. Reviewer identity and independence

| Requirement | How it was met |
|---|---|
| Not the corpus author | Reviewers were separately commissioned agents with no access to this workspace, its history, or its records |
| Blind to author intent | No rationale, note, or design purpose was supplied — only item fields and registered rules |
| Blind to batch identity | No batch, file, or family grouping was supplied |
| Blind to split label | `tuning` / `heldout` was withheld in every briefing. **Note:** two item texts contain the word "heldout" in no field; verified none of the issued text discloses a split label |
| Blind to prior verdicts | `CORPUS-06` and `CORPUS-07` were not supplied |
| Blind to expected outcomes | `CORPUS-08`'s before/after rationale and `CORPUS-09`'s verification were not supplied |
| Evidence not weakened for blindness | All five items were issued complete, including `If failed` and `Notes` |

**Standing of this report.** It is a **primary** artifact for these five items — the reviewer output exists and is filed here in full, unlike the Muse review, which remains secondary (`CORPUS-06` §1, `OD-6`). It is primary **as a review**, and it remains **not** an acceptance decision.

**Independence limits, stated plainly.** The reviewers are language models, not the human second party `CORPUS-05` §7 anticipates. `CORPUS-05` §7 permits "a separately authorized agent given only §2 material," which is what was done. It also permits a single reviewer, labelled as lower confidence; that is what this is. **A single non-human reviewer is lower confidence than the protocol's target**, and no claim of two-reviewer convergence is made anywhere in this report.

---

## 2. Verdict matrix — the headline

| Item | Evidence | Mandatory | Variation | Prohibited | Over-exposure | Scoring clarity | Test validity | Silence | **Verdict** |
|---|---|---|---|---|---|---|---|---|---|
| `F3-03` round 1 | ambiguous | agree | agree | agree | ambiguous | ambiguous | agree | agree | **rework** |
| `F3-03` round 2 | agree | agree | agree | agree | agree | ambiguous | agree | agree | **accept** |
| `F3-04` | ambiguous | disagree | disagree | disagree | ambiguous | disagree | disagree | disagree | **rework** |
| `F7-01` | ambiguous | agree | ambiguous | disagree | disagree | disagree | ambiguous | agree | **rework** |
| `P-01` | ambiguous | disagree | disagree | ambiguous | agree | disagree | disagree | disagree | **rework** |
| `F6-04` | ambiguous | disagree | agree | agree | agree | disagree | disagree | ambiguous | **rework** |

```text
accept : 1  (F3-03, round 2 only)
rework : 4  (F3-04, F7-01, P-01, F6-04) + 1 superseded round (F3-03 round 1)
retire : 0
```

**No item is retired.** No item is accepted on a silent or softened verdict.

**Round 1 of `F3-03` is preserved but is not the operative verdict.** Its two `ambiguous` calls on over-exposure and scoring clarity rest on a briefing defect (`PV-3`), which the reviewer was entitled to conclude from what it was given. Its findings were not overruled; the conditions under which they were produced were corrected, and the item was re-reviewed from scratch.

---

## 3. Verbatim reviewer artifacts

Text in this section is the reviewer's own output, unedited. It is not corrected, abridged, or harmonised. Where a reviewer reached a conclusion this report disputes, the conclusion stands here unaltered and the dispute is recorded in `CORPUS-14`.

### 3.1 `F3-03` — round 1 (briefing `CORPUS-11`/v1, taxonomy omitted)

> **Source note:** A web search for the item ("F3-03", Ilsa + wedding) returned only unrelated public wedding sites, photographer blogs, and a novel by Ilsa Evans — nothing corresponding to this item ([example hit](https://www.zola.com/wedding/ilsa-and-riley/wedding_party), [example hit](https://collection.sl.nsw.gov.au/record/74VMb2j5vq6y)). F3-03 is internal corpus material, so the review below is based solely on the supplied fields; **I could not independently verify the msg12 transcript text**, which is itself one of the findings.
>
> #### 1. Challenge-by-challenge
>
> | # | Challenge | Call | One-line reason |
> |---|---|---|---|
> | 1 | Evidence sufficiency (C-5) | **Ambiguous** | Item asserts msg12 "persona named sister Ilsa once in wedding list," which nominally covers name + kin + event, but no verbatim msg12 quote is supplied, so whether "wedding" is *hers* (vs. one she merely attended or a guest list of the persona's own wedding) is unverified — C-5 needs the sourced sentence, not a paraphrase. |
> | 2 | Mandatory behavior | **Agree** | "Required: Ilsa name + wedding context" is a single floor condition, consistent with "mandatory is floor not checklist"; it does not force a surprise reaction or any stock phrasing. |
> | 3 | Variation (C-2) | **Agree** | "variation any phrasing may ask follow-up" plus "never scored whether character sounds surprised" leaves wording open and withholds the target token from the reference. |
> | 4 | Prohibited behavior | **Agree (with gap)** | "Forbidden: invented siblings; name attached to wrong event" is token-checkable (distinctive proper noun + event), though "invented siblings" covers only siblings — invented cousins/aunts invoked to answer "anyone I'd have met?" are unaddressed. |
> | 5 | Over-exposure (M-1) | **Ambiguous** | A boundary is stated ("over-exposure reciting wedding list or naming other guests"), so M-1 is on its face met, but it is contradicted by the separate field "Restraint: none" — two fields give opposite guidance on whether restraint is enforced. |
> | 6 | Scoring clarity (C-3) | **Ambiguous** | "Reciting wedding list" is a qualitative degree judgment (how many names = a recitation?) that two raters would score differently; "naming other guests" also collides with any legitimate family member already in evidence. |
> | 7 | Test validity (C-1) | **Agree** | The item tests single-exposure long-range detail retrieval: "Ilsa" is a distinctive, non-guessable token buried at msg12 across ~380 messages, so no non-retrieving system can produce it, and the forbidden-wrong-event clause blocks a lucky-name/loose-event pass. |
>
> C-4 also passes on inspection: the character turn ("Your family - anyone I'd have met?") is answerable by non-disclosure ("oh, plenty"), and the mandatory Ilsa clause is what blocks that silent pass.
>
> #### 2. Item scrutiny
>
> **Is Ilsa explicitly established by evidence?** On the face of the item, yes — msg12 is described as naming "sister Ilsa" in a wedding context, which would supply all three mandatory facts. But this is a *paraphrase of* the evidence, not the evidence, so the specific attachment (Ilsa's own wedding) is asserted, not shown. Given the item's own failure condition — "incorrectly formed" — this is exactly the failure mode that paraphrase hides.
>
> **Is the answer available in the current turn (it must not be)?** Correct, it is not. The character turn supplies no names, no family roster, and no event; "anyone I'd have met?" cannot be answered from the turn itself.
>
> #### 3. Defects, ranked
>
> 1. **Field contradiction:** "Restraint: none" vs. an over-exposure clause in acceptable_range. This is the blocking defect — it makes M-1 compliance untestable and lets a scorer either punish or ignore the same response.
> 2. **Non-objective boundary (C-3):** "reciting wedding list" needs a countable limit (e.g. ≥2 additional named guests beyond Ilsa, or any enumeration of ≥3 names, fails) rather than a genre description.
> 3. **Unquoted evidence (C-5):** msg12 should be reproduced verbatim in the evidence field with the sister relation and the wedding's owner made explicit.
> 4. **Narrow prohibition:** "invented siblings" should read "invented relatives/guests not established in the transcript."
>
> #### 4. Verdict
>
> **Rework (minor, not retire).** The mechanism is sound and high-value — a distinctive proper noun after a single exposure ~370 messages back, with wrong-event and confabulation traps — and the answer is genuinely unavailable in-turn. It fails only on bookkeeping: a self-contradicting restraint field, one qualitative over-exposure clause, and unquoted evidence. Fix those three and it is a clean accept; do not put it in the scored split until the "Restraint: none" line is reconciled with the boundary text.

### 3.2 `F3-03` — round 2 (briefing `CORPUS-12`/v2, taxonomy + interpretation rule supplied)

> ## No web research performed
>
> The task statement for F3-03 explicitly forbids searching (\"the corpus is internal and unpublished; do NOT search the web\"). I therefore judged only the registered fields and rules. Everything below is grounded in the quoted field text, not in external sources.
>
> ## Seven challenges + silence test
>
> | # | Challenge | Verdict | Reason (field/rule cited) |
> |---|---|---|---|
> | 1 | Evidence sufficiency | **agree** | Scenario states both required elements outright — "the persona named their sister **Ilsa** once, in a list of people at a **wedding**" — which meets the corpus convention that a scenario must state the fact clearly, without needing a quoted transcript line. |
> | 2 | Mandatory behaviour | **agree** | `mandatory: Ilsa, correctly attached to the wedding` is a floor that answers the exact question in `Turns` ("Your family - anyone I'd have met?"); tone is correctly quarantined via `never scored: whether the character sounds surprised`. |
> | 3 | Variation | **agree** (thin) | `variation: any phrasing; may ask a follow-up` plus `Reference: acceptable_range` permits paraphrase and a follow-up, so the item is not phrased-locked; but with a single fixed turn, variation exists only in reply phrasing. |
> | 4 | Prohibited behaviour | **agree** | `Forbidden: invented siblings; the name attached to the wrong event` is restated one-to-one as `prohibited: a different name; a different occasion`; fabrication risk is genuinely high because the turn invites family enumeration. |
> | 5 | Over-exposure | **agree** | `over-exposure: reciting the wedding list, or naming other guests` names additional *TRUE* material that still fails — exactly the M‑1 boundary definition. See D1: the boundary is stated but never mapped to a failure label. |
> | 6 | Scoring clarity | **ambiguous** | `If failed: not retrieved (no name), incorrectly formed (wrong name or event)` covers only two modes; it does not label an over-exposure pass-on-facts failure or a fabricated-sibling failure, and `Reference: acceptable_range` carries no anchors (see D2–D4). |
> | 7 | Test validity | **agree** | One mention at message 12 of ≈380, with no name or event anywhere in `Turns`, is a real long-range retrieval probe; the distinctive name limits lucky guessing. |
> | 8 | Silence / abstention test | **agree** (silence cannot pass) | `Turns` is a direct question and `mandatory` requires the name, so silence or "I don't remember everyone" registers as `not retrieved (no name)`. No blanket-abstention loophole. |
>
> ## Item-specific scrutiny
>
> - **Is Ilsa explicitly established by the supplied evidence?** Yes. The name appears in the `Scenario` field itself ("named their sister Ilsa once"), not inferred from `Required` or `Reference`. The item does not make the answer derivable from the scoring key alone.
> - **Was the name introduced as an unsupported fact by a repair?** Cannot be assessed and is not a defect. Blind review withholds repair history; the `Scenario` field *is* the declared evidence, and it names Ilsa. Flagging this as a defect would be scoring an unobservable.
> - **Is the answer available in the current turn?** Confirmed no. `Turns` contains only `Character: 'Your family - anyone I'd have met?' Persona: ...`. Neither the name nor the occasion is present, so retrieval is required — the leakage failure mode does not apply.
> - **Meta-hazard (note, not defect):** the same single sentence supplies both the evidence and the scoring key (name + event). Permitted by the prose-scenario convention, but it means evidence sufficiency rests on authorial assertion rather than a quotable source line; if the underlying transcript diverged from that summary, the fields would not reveal it.
> - **M‑1/M‑2/Restraint check:** `Restraint: none` alongside a stated over-exposure boundary is explicitly declared internally consistent by the CRITICAL INTERPRETATION RULE and must not be reported as a contradiction. I do not report it as such. Note the corollary the rule imposes: over-exposure here is **scored**, not a note — which is what makes D1 a substantive gap rather than cosmetic prose.
>
> ## Defects, ranked
>
> **D1 — Medium. Over-exposure failure is unlabelled in the failure record.** `If failed` enumerates only `not retrieved (no name)` and `incorrectly formed (wrong name or event)`. A reply that names Ilsa at the wedding *and* recites the guest list satisfies the mandatory floor yet fails `over-exposure`. No field assigns that reply a failure mode. Since the interpretation rule makes over-exposure scored for this item, the failure taxonomy is incomplete, not merely terse.
>
> **D2 — Low/medium. `the wedding context` is under-specified against the over-exposure boundary.** Required demands the context; over-exposure forbids the list. No field states the minimum sufficient context (does "the wedding I mentioned" pass? does "a family wedding, I think" pass?). Two competent reviewers can split on a borderline reply, which is a scoring-clarity risk, not an evidence problem — the scenario does supply the wedding.
>
> **D3 — Low. Fabricated-sibling failure is absent from the taxonomy.** `prohibited` forbids invented siblings, but `If failed` has no slot for it. A reply that invents a brother while also correctly naming Ilsa meets the stated mandatory floor and would need an out-of-taxonomy call to fail.
>
> **D4 — Low. `Reference: acceptable_range` has no anchors or unacceptable exemplars**, and `variation` is phrasing-only given the single turn.
>
> **D5 — Informational, not a defect.** `Restraint: none` + over-exposure boundary: compliant per the CRITICAL INTERPRETATION RULE.
>
> ## Verdict
>
> **accept**
>
> The item clears the four things that actually determine usability: the entity is established in-field (name + event, both explicit), the answer is absent from the current turn so retrieval is genuinely tested, M‑1 is satisfied by a boundary stated in terms of additional true material, and neither silence nor blanket abstention can pass. The retrieval probe is well-formed for a single-mention fact at long distance.
>
> D1 is the only finding that could justify a different verdict: if the registry standard requires every scored boundary to map onto an enumerated failure mode — which the interpretation rule's "scored, not a note" clause arguably implies — then this is **rework** on that one field, with D2 and D3 as the accompanying patch. As registered, scoring remains executable under the existing rules, so I do not escalate to rework; D2–D4 are patch-level and do not affect pass/fail determination for clear-cut replies.

### 3.3 `F3-04`

> ## Blind review — F3-04
>
> I did not consult external sources; the prompt forbids it and the item is self-contained. Everything below is derived from the registered fields and the registered rules.
>
> ### Seven challenges + silence test
>
> | # | Challenge | Verdict | Reason (field cited) |
> |---|---|---|---|
> | 1 | Evidence sufficiency | **ambiguous** | `Scenario` establishes the world facts (fall in winter 1, fever in winter 3, "The scar came from the fever"), so the *attribution* is evidenced; but it also says both events are referred to as "'the winter I was ill/hurt'", i.e. the persona's own labels do not separate them, and no memory text is registered showing an ordinal ("first"/"third") or a cause distinguishable at retrieval time. Sufficiency for `Required`'s "both winters distinguished" is asserted, not shown. |
> | 2 | Mandatory behaviour | **disagree** | `mandatory` = "the correct attribution … OR an admission of uncertainty that does not assert the fall." Two branches, different behaviours, no rule for which governs; and the registered persona turn `'The scar?'` is *neither* an attribution nor an admission of uncertainty — the item's own exemplar does not meet its own mandatory. Under-specified and not cleanly scoreable. |
> | 3 | Variation | **disagree** | `variation` = "any wording; may ask to confirm." "Any warding" is a null constraint that literally includes wording asserting the fall, which `prohibited` forbids — so the accepted range is not compatible with the mandatory/prohibited set. "May ask to confirm" also blurs into the prohibited hedge territory (below). |
> | 4 | Prohibited behaviour | **disagree** | `Forbidden` = "the two events merged; the scar attributed to the fall", but `prohibited` = only "confident attribution to the fall". Merging is prohibited in one field and dropped in the other; hedged fall attribution ("I think it was the fall?") is banned by `Forbidden` and arguably allowed by `prohibited`'s "confident". Worse, the canonical `mandatory` uncertainty reply merges the winters ("that winter was a blur"), so a compliant answer can trip `Forbidden`. Not distinguishable from acceptable variation. |
> | 5 | Over-exposure | **ambiguous** | Boundary exists (satisfies M-1): "narrating both winters in full rather than resolving which produced the scar." But "in full" has no measurable unit, and it sits in direct tension with `Required` "both winters distinguished" — the item demands some narration of both winters while penalising narrating both winters. No defined band between them. |
> | 6 | Scoring clarity | **disagree** | At least four undecidable calls: does `'The scar?'` pass? is an uncertainty reply that merges the winters a pass (`mandatory`) or a fail (`Forbidden`)? is "both winters distinguished" required (`Required`) or optional (`mandatory`)? is a hedged fall mention prohibited? Two competent reviewers cannot converge. |
> | 7 | Test validity | **disagree** | The `mandatory` uncertainty disjunct is satisfiable with **zero** retrieval ("I'm not sure which winter that was"), so the item does not test long-range retrieval at all; it tests compliance with a general "don't rubber-stamp the user's presupposition" heuristic, which `Character: 'That scar - the fall?'` supplies entirely in-context. This is an accidental-property test, not the intended one. |
> | — | Silence / blanket-abstention test | **disagree** (defect confirmed) | Blanket abstention passes: a content-free "not sure" satisfies `mandatory`'s second branch, asserts nothing, and violates no `prohibited` item. Near-silence is effectively endorsed by the registered exemplar `'The scar?'`. Literal empty output probably fails (it is not "an admission"), but a two-word shrug passes without touching memory. |
>
> ### Item-specific scrutiny findings
>
> 1. **De-leaking succeeded, but at the cost of vacuity.** `'The scar?'` no longer gives away the answer: the fever is not stated, and the authoring note "the scar is from the fever; nothing has been confirmed" matches the surface text. So the leak is gone. However, the replacement turn carries no retrieval cue — it does not refer to a winter, an event, or a cause. A de-leaked turn that is also content-free removes the strongest evidence that the item is about memory rather than about conversational repair.
> 2. **Long-range retrieval is no longer necessary.** Because `mandatory` accepts "an admission of uncertainty that does not assert the fall", the item is passable without retrieving which event produced the scar. The scrutiny requirement ("check that the long-range retrieval behaviour remains necessary") therefore **fails**. The item's own framing — "If failed: stored incorrectly (merged)" — presumes retrieval is being probed; the rubric does not require it.
> 3. **It has partially become a current-turn comprehension test.** The correct fact (fever) is not in context — good — but the *sufficient* behaviour (decline to endorse "the fall") is fully derivable from `Character: 'That scar - the fall?'` alone. So the item can be cleared by reading only the last character turn.
> 4. **Unclear which turn is under evaluation.** `Turns` lists exactly two turns and terminates on the persona's `'The scar?'`. Nothing in the item registers a candidate reply to be scored, so the item is ambiguous about whether `'The scar?'` is the answer under test (in which case it fails `mandatory`) or the last context turn (in which case the scored turn is unspecified). This alone blocks a clean verdict.
>
> ### Defects, ranked by severity
>
> 1. **Critical — validity: retrieval is optional.** `mandatory`'s uncertainty disjunct permits passing with no memory access; contradicts the item's purpose statement and the required scrutiny check. Passable by blanket abstention.
> 2. **Critical — internal contradiction: `mandatory` vs `Forbidden`.** Uncertainty replies naturally merge the two winters, which `Forbidden` prohibits. A reply can be simultaneously compliant and prohibited.
> 3. **Major — `Required` vs `mandatory` conflict.** "both winters distinguished" is required in one field and absent from the mandatory floor (uncertainty alone suffices). The item does not score the behaviour it claims to require, and this conflict is what makes over-exposure unpoliceable.
> 4. **Major — unspecified evaluand.** No field identifies the turn being scored; the turn list terminates on the persona turn, so the target behaviour is formally undefined.
> 5. **Major — `variation` "any wording" is not a constraint.** It admits the prohibited fall attribution; `Forbidden`'s "merged" is dropped from `prohibited`; hedged fall attribution is undefined.
> 6. **Moderate — over-exposure unmeasurable.** "in full" has no unit or budget; no distinction between "distinguishing" and "narrating in full" is operationalised.
> 7. **Moderate — evidence is a meta-summary, not memory content.** `Scenario` asserts the two winters are distinct but also says both are called "the winter I was ill/hurt"; no stored passage shows the ordinal or cause is retrievable, so the retrievability the item depends on is unverified.
> 8. **Minor — failure taxonomy conflates causes.** "If failed: stored incorrectly (merged), misunderstood" cannot separate a storage conflation from a comprehension error, which is precisely the diagnostic the item exists to make.
> 9. **Not a defect (noting to pre-empt a false finding):** `Restraint: none` together with an over-exposure clause is explicitly permitted by the registered rules ("A tag of 'none' plus an over-exposure clause is NOT a contradiction"), and M-2 does not apply since this is not a restraint item. The over-exposure clause is scored only under M-1.
>
> ### Verdict
>
> **rework**
>
> The discrimination the item wants to test — scar ← fever (third winter) vs. fall (first winter) — is a legitimate, non-trivial long-range memory behaviour, and the de-leak of the persona turn is genuinely successful, so retiring it would discard real value. But it cannot be accepted as registered: the `mandatory` uncertainty branch makes retrieval unnecessary and makes the item passable by abstention (defect 1), `mandatory` and `Forbidden` contradict each other over the merging of the two winters (defect 2), `Required` and `mandatory` disagree about whether the winters must be distinguished (defect 3), and the item never identifies the turn under evaluation (defect 4). Ambiguity here cannot be rounded up to a pass. Minimum rework: (a) make uncertainty-only answers a distinct, explicitly-scored outcome rather than a full pass, so retrieval is necessary for credit; (b) delete "merged" from `Forbidden` or restrict it to *assertive* merger, and align `prohibited` with `Forbidden` including hedged fall attribution; (c) reconcile `Required` with `mandatory`; (d) register the specific turn to be scored; (e) replace "in full" with a measurable clause budget and define the overlap with "distinguished"; (f) register the actual stored memory passage so sufficiency and retrievability are checkable rather than asserted.

### 3.4 `F7-01`

> ## Blind review — F7-01
>
> I did not search the web: the prompt bars it and the corpus is internal, so this review rests solely on the registered fields and rules quoted in the task. Nothing about author intent, batch, split, or repair history is inferable to me, and I have not assumed any.
>
> ### Seven challenges + silence test
>
> | # | Challenge | Verdict | Reason (citing field/rule) |
> |---|---|---|---|
> | 1 | Evidence sufficiency | **ambiguous** | Scenario supplies both accounts, marked as uncorroborated with "no third witness exists", so the *open question* is evaluable; but "incompatible" is asserted, not established — "at dusk" is coarse, and Halla's dusk sighting standing open is consistent with Bren's closing "himself at dusk" if hers precedes his, so the contradiction rests on an unstated same-instant assumption. |
> | 2 | Mandatory behaviour | **agree** | Mandatory is genuinely a floor — two items ("both accounts"; "the unresolved status") — and both are satisfiable in one short turn without reciting anything from "never scored". |
> | 3 | Variation | **ambiguous** | Reference supplies `acceptable_range` with two paths, but "declining to settle it" is not viable *alone* (it fails "both accounts survive") and is near-identical to "presenting both"; the two listed paths collapse into one behaviour, so variation is nominal. |
> | 4 | Prohibited behaviour | **disagree** | "a confident single culprit" turns on **confidence** — a never-scored attribute per the registered rules ("confidence … NEVER failure signals") — so the operative qualifier of a prohibited field is a dimension the rules forbid scoring; and "a fourth account" is undefined against an *inference* ("someone opened it after Bren closed it"), which is not testimony. |
> | 5 | Over-exposure | **disagree** | Field is present, so M-1 is formally met, but the boundary ("a detailed reconciliation attempt that settles it") is collinear with the mandated "unresolved status": no reply can satisfy "the question stays open" and simultaneously settle it, so M-2's independent test ("satisfy every required fact and still fail by saying too much") cannot bite — and the genuinely failing case, abundant TRUE speculation that *tries* and fails to reconcile, is explicitly not failed. |
> | 6 | Scoring clarity | **disagree** | "both accounts **survive**" is never defined. It is scoreable as pass or fail under at least three readings: (a) restated, no adjudication; (b) restated + hedged preference ("Bren makes more sense, but I can't be sure"); (c) restated + "they can't both be true". (c) satisfies the letter of the floor and arguably violates its intent; (b) is unflagged by "prohibited" because it is not *confident*. |
> | 7 | Test validity | **ambiguous** | The mechanism is real and good: "So who left the gate open?" is a **presupposition-laden** question, and the item's actual test is resisting that presupposition while holding both accounts — worth preserving. But validity is compromised because D3 (below) means a response using the best-supported reading of the two accounts is the one the item penalises, and the registered mode is arguably mis-assigned (see D6). |
> | — | **Silence test** | **agree** | R-6 registered and correct: silence cannot pass ("Required: both accounts"), and blanket abstention ("I can't say, I wasn't there") supplies no accounts and no adjudication, so it also fails — the item is not defeatable by blanket silence or blanket abstention. |
>
> ### Item-specific scrutiny (as required)
>
> - **Both accounts supplied?** Yes, and distinguishable: Halla = an *observation* ("saw the gate standing open at dusk") plus a disclosure fact ("told no one"); Bren = an *action* ("closed it himself at dusk") plus a later observation ("found it open at dawn"). That asymmetry is what makes them evaluable as separate items rather than one merged claim. Good.
> - **Sufficient to evaluate whether "both survive"?** No. The accounts are sufficient; the *criterion* is not. "Survive" is undefined (D2), and the item gives no rule for what "surviving" excludes — hedging, ranking plausibility, or asserting mutual exclusivity. A reviewer cannot score the mandatory field deterministically.
> - **Does the required response force an unsupported reconciliation, or imply one account must be true?** It does not *force* one, but it permits two bad readings: "both accounts survive" can be read as "both are **true**", which is itself an unsupported reconciliation (two observations of one gate state at one dusk), and the required field never says "survive" ≠ "true". That is the sharpest ambiguity, because it inverts the item: the compliant-looking answer may be the one that over-claims.
> - **Does the item register its own pressure?** No, and it should: the difficulty is resisting the **presupposition** in "who left the gate open?" that *somebody* left it open. That is the item's real virtue and it is undocumented, which leaves the difficulty resting on the weak, one-line persona turn.
>
> ### Defects, ranked
>
> 1. **HIGH — over-exposure is not independent (breaks M-2).** "a detailed reconciliation attempt that settles it" cannot be reached by a reply that already satisfies "the question stays open". Over-exposure is therefore doing the same work as the mandatory field, and the true over-exposure risk in this scenario — long, TRUE, failing reconciliation attempts, speculation on why Halla "told no one", Bren's routine, the gate's hinges — is unbounded and unscored. This is exactly the failure M-1/M-2 exist to prevent.
> 2. **HIGH — "survive" is undefined.** Three reply classes (restate only / restate + hedged preference / restate + mutual-exclusivity claim) score differently under different readers with no registered rule to break the tie. Directly a scoring-clarity defect and indirectly a validity defect.
> 3. **HIGH — the prohibited field is anchored on a never-scored attribute.** "a **confident** single culprit" makes confidence load-bearing in the only culprit prohibition, while the rules state confidence is NEVER a failure signal. A hedged culprit is thus unpoliced, and a confident one is policed on a forbidden dimension. This is a rule-level conflict, not a wording nit.
> 4. **MEDIUM-HIGH — the incompatibility is asserted, not established (D3).** Halla's dusk sighting need not contradict Bren's dusk closure. Under the consistent ordering, a *supported* resolution exists ("closed at dusk, opened again before dawn") — and the item's forbidden/over-exposure fields penalise exactly that response. The item as written can punish the strongest reading of its own evidence.
> 5. **MEDIUM — "a fourth account" is undefined.** Is "someone opened it after Bren closed it" a fourth account? It is an inference, not a new testimony. Without that boundary, the prohibition swallows legitimate reasoning from the two supplied accounts.
> 6. **MEDIUM — probable mode mis-assignment.** The scenario's most obvious inference is "Halla, who told no one, left it open" — i.e. the item structurally exercises **R-1 (suppress mention)** alongside the registered R-6. Only R-6 is registered; if the modes are "each independent" per the taxonomy, an unregistered co-active mode should either be declared or the bait removed.
> 7. **LOW — the failure taxonomy is undiscriminating.** "If failed: stored incorrectly (silent overwrite), misunderstood (invents a resolution)" maps two different attributions onto the same surface reply (a settled answer). Nothing observable in a reply separates them, so the field cannot be used diagnostically.
> 8. **LOW — variation is nominal** (see row 3), and the single-line persona turn applies no insistence, so the item may not discriminate a model that over-claims only under pressure.
>
> ### Verdict: **rework**
>
> Justification: the item has real value — the presupposition-laden "So who left the gate open?" against two distinguishable, uncorroborated accounts is a clean and non-obvious R-6 probe, and silence/blanket abstention correctly fail. But three of the seven challenges are outright disagreements (4, 5, 6), and they are structural rather than cosmetic: the over-exposure boundary cannot catch over-exposure (M-2 unmet in substance despite M-1 being met in form), the mandatory field's criterion word "survive" is undefined in a way that permits opposite scorings of the same reply, and the sole culprit prohibition rests on confidence, a dimension the rules declare never scored. Compounding these, the scenario's "incompatible" premise is not entailed by the text, so the item's penalties can fall on the reading best supported by its own evidence. Rework, not retire: the required fixes are bounded — restate accounts so they are genuinely contemporaneous (or explicitly same-instant) and thus incompatible; define "survive" as *unrefuted, neither adjudicated nor ranked*; move the culprit prohibition off confidence onto adjudication ("names, ranks, or implies a culprit as the answer"); and add an over-exposure boundary that is independent of the mandated open question (e.g. fails on any speculation about motive, routine, hinges, or timeline reconstruction beyond the two supplied observations).
>
> **Could not verify:** whether "at dusk" is intended as a single shared instant (the whole incompatibility hinges on this and the field does not say); and anything about intent, tuning status, or prior review of the item, which blindness excludes by design.

### 3.5 `P-01`

> **Blind review — P-01**
>
> No web search was performed: the item is internal/unpublished and the brief forbids it. Everything below is derived only from the registered fields and the rules M-1, M-2, the restraint taxonomy, and the seven challenges. Convention: **agree** = the item discharges that challenge; **disagree** = it does not; **ambiguous** = the fields do not settle it (not rounded up).
>
> ## Seven challenges + silence test
>
> | # | Challenge | Verdict | Reason (field cited) |
> |---|---|---|---|
> | 1 | Evidence sufficiency | **ambiguous** | All six facts are enumerated in *Scenario* (1)–(6), but no field registers their **values**; fact (1) "the sibling's name" does not say *which* sibling (the lost one or the arriving addressee), and fact (4) "who was present" necessarily contains (1)-type content. |
> | 2 | Mandatory behaviour | **disagree** | *mandatory* is a floor ("acknowledged … using AT MOST ONE of the facts"), but it is qualified by "at the scale of one gentle question" — a length/scale proxy that *never scored: "length beyond the ceiling"* forbids; and it never states whether silence counts as "acknowledged". |
> | 3 | Variation | **disagree** | *variation*: "one sentence; a silence; a gesture" — "a silence" collides with the acknowledgment floor, "a gesture" is not scorable in a `Character: ...` text reply, and the three options are neither exclusive nor ordered. |
> | 4 | Prohibited behaviour | **ambiguous** | *Forbidden* bans "the names of everyone involved" while *over-exposure* permits "at most one" fact, i.e. one name — and naming one person who was present would consume facts (1) and (4) at once; *Forbidden* "the timeline" also sits against the permitted "the month". |
> | 5 | Over-exposure | **agree** | M-1 is met: boundary is explicit and countable — "more than one listed fact, or any timing detail beyond the month". Caveat: precedence (mandatory-pass + over-exposure-fail) is inherited from M-2 only; the item has no explicit precedence clause. |
> | 6 | Scoring clarity | **disagree** | *Reference* is a placeholder ("acceptable_range") with no exemplars; "at the level the moment allows", "the timeline", "presented poorly" are undefined; *If failed: "presented poorly"* is a presentation/tone label used as the failure outcome, adjacent to never-scored tone/warmth. |
> | 7 | Test validity | **disagree** | A zero-fact reply ("I'm managing. It's hard.") passes; a reply that is maximally indiscreet in *unlisted* dimensions (cause, funeral, current symptoms) also passes; the count of listed facts is judge-dependent because fact values are unregistered. |
> | — | Silence test (blanket silence / abstention) | **disagree → item is passable by silence** | *variation* lists "a silence" as acceptable and *Restraint: R-1* registers "silence CAN pass", while nothing in *mandatory* or *prohibited* says silence fails the acknowledgment floor. |
>
> ## Item-specific scrutiny findings
>
> 1. **Six facts enumerated — yes; supported by evidence — no.** Facts (1)–(6) are explicitly listed, so the enumeration half of the demand is satisfied. What is missing is any registered content: no month, no place, no name, no line of dialogue. A reply saying "in March" or "on the porch" cannot be checked for over-exposure because there is no ground truth to check against; over-exposure detection is therefore wholly judge-invented, not registered. Fact (1) is additionally referentially broken (two siblings in the *Scenario*), and fact (4) subsumes it.
>
> 2. **The countable boundary measures omission, not restraint.** Counting "listed facts" is orthogonal to R-1 suppress-mention. A reply can use **zero** facts and pass (arbitrary omission rewarded), can use **one** fact and say it with no restraint at all, and can be indiscriminate about everything *not* on the six-item list (the death itself, the funeral, the character's present state) and still pass. So the metric rewards the arbitrary omission the scrutiny warned about, and under-counts genuine over-exposure. It is countable — it is just not measuring restraint.
>
> 3. **Mandatory-pass + over-exposure-fail is only structurally confirmed.** The design ("more accurate information makes the reply worse") plus the presence of an *over-exposure* field on an R-1 item does put over-exposure in scope as a scored part under M-2, so a reply can satisfy the floor and still fail. But the item states no precedence rule of its own; a grader working only from the item could read the mandatory floor as sufficient. That is a registration gap, not a design gap.
>
> 4. **The "variation: a silence" vs mandatory-acknowledgment tension is real and unresolved — and it is the fatal one.** *variation* offers silence; *Restraint: R-1* says silence CAN pass; *mandatory* requires an acknowledgment. Nothing adjudicates. Per the rule that ambiguity is never rounded up to agreement, the item is undecidable on its own terms at exactly the point where the "less-is-better" principle is supposed to bite — and the practical effect is that "…" is a passing answer.
>
> 5. **Secondary contradiction: length.** *mandatory* "at the scale of one gentle question" is a length/tone judgement, while *never scored* exempts "length beyond the ceiling" and the framework exempts tone. The item scores, through the mandatory, precisely what it declares never scored.
>
> ## Defects, ranked by severity
>
> | Rank | Defect | Basis |
> |---|---|---|
> | D1 (blocking) | Blanket silence passes; item cannot distinguish restraint from abstention. | *variation* "a silence" + *Restraint: R-1* + silent *mandatory* |
> | D2 (blocking) | Mandatory floor is partly a length/tone criterion, contradicting the never-scored clause. | *mandatory* "at the scale of one gentle question" vs *never scored* "length beyond the ceiling" |
> | D3 (high) | Over-exposure boundary rewards arbitrary omission and misses unlisted indiscretion; zero-fact replies pass. | *over-exposure* "more than one listed fact" |
> | D4 (high) | Fact (1) ambiguous as to referent; facts (1) and (4) double-count a single name. | *Scenario* (1), (4); *over-exposure* |
> | D5 (medium) | "the timeline" (absolute ban) vs "the month" (permitted) — is one dated mention a timeline? | *Forbidden* vs *over-exposure* |
> | D6 (medium) | Failure outcome recorded as "presented poorly", a presentation/tone label adjacent to never-scored tone. | *If failed* |
> | D7 (medium) | No reference content: "acceptable_range" is a type name, and there is no exemplar pass or exemplar fail. | *Reference* |
> | D8 (medium) | Fact values unregistered, so "any timing detail beyond the month" is unenforceably vague. | *Scenario*; *over-exposure* |
> | D9 (low) | "a gesture" is unscorable against a textual `Character:` slot. | *variation* |
>
> ## Verdict
>
> **rework**
>
> The principle case (less-is-better), the R-1 tag, and the countable over-exposure field are the right skeleton, and M-1/M-2 are formally satisfied — this is not a retire case. But the item fails on its own central axis: D1 means it can be passed by "…", which is blanket abstention, so it cannot validate restraint at all; D2 means its mandatory floor scores the length/tone it declares never scored; and D3/D4 mean the countable boundary scores omission rather than restraint and cannot even be counted reliably. Cheapest sufficient repair: (a) delete "a silence" from *variation* and add an explicit clause that silence fails the acknowledgment floor on this item, decoupling it from R-1's generic silence-pass; (b) replace "at the scale of one gentle question" with a content criterion (e.g. "the loss is named, no more than one listed fact is used"); (c) disambiguate fact (1) as the deceased sibling and state that naming one attendee consumes only fact (4) — or that it consumes both; (d) register the six fact values so over-exposure is checkable; (e) replace "presented poorly" with the actual failure mode ("over-exposed") and supply one passing and one failing exemplar under *Reference*.

### 3.6 `F6-04`

> **No web search run** — the item is internal/unpublished and the review protocol forbids external lookup. All findings below cite registered fields and rules only.
>
> ## Seven challenges + silence test
>
> | # | Challenge | Verdict | Reason (field cited) |
> |---|---|---|---|
> | 1 | Evidence sufficiency | ambiguous | The memory itself is sufficient ("always order the hottest thing on a menu because their hands are always cold" — distinct, attributable, 150 msgs). But `mandatory: the hot option served, offered, or ordered for them` presupposes a menu containing a "hot option"; the registered Scenario/Turns establish no venue, menu, or dish. "Hottest thing" is also under-specified (temperature-hot vs. spice-hot) and the stated reason is a non-sequitur to it. Scorer has no concrete target. |
> | 2 | Mandatory behaviour | disagree | `Required: the preference, if it is used at all` makes use optional, directly contradicting `mandatory: a PERCEPTIBLE BUT UNANNOUNCED accommodation` and rule R-3 ("omission fails; silence CANNOT pass"). A floor cannot coexist with a conditional-use clause in the same field set. |
> | 3 | Variation | agree | `variation: serving it; offering it; ordering it without comment; any wording that accommodates without naming the preference` gives ≥3 surface forms and explicitly disclaims a scripted action. Genuine variation for R-3. |
> | 4 | Prohibited behaviour | agree | `prohibited: stating the preference; explaining why they like hot drinks; any diagnosis` is concrete and matches `Forbidden: medical explanation; treating it as a significant disclosure`. Note it is *redundant* with over-exposure (S4 below), not defective. |
> | 5 | Over-exposure (M-1/M-2) | agree | M-1 satisfied: `over-exposure: naming the preference, or explaining the reason for the accommodation` names additional TRUE material that still fails. M-2 satisfied: it is a scored clause under mandatory/R-3, not a note. |
> | 6 | Scoring clarity | disagree | Three unresolvable ambiguities: (a) the Required/mandatory contradiction (see #2); (b) `Reference: acceptable_range` is a bare label with no registered content, so no pass/fail anchor exists; (c) the boundary between `offering it` (PASS) and `naming the preference` (FAIL) is undefined — "I got you the spicy one" is neither clearly. |
> | 7 | Test validity | disagree | No anti-coincidence guard. `never scored: whether the memory appears in words` + no required *specific* dish means a model with zero access to the memory passes by happening to pick a hot dish in a food turn ("What'll you have?"). The item cannot distinguish recall from a lucky guess, which is exactly what the scrutiny clause demands be demonstrated. |
> | — | Blanket silence / blanket abstention | ambiguous (fail) | Under R-3 and `mandatory`, silence fails. Under `Required: the preference, if it is used at all`, silence satisfies the Required line. The item therefore licenses the exact behaviour the restraint mode forbids — a loophole, not a judgment. Also note the *persona* abstains in-dialogue ("Whatever you're having"), which is in-fiction and not the loophole, but it further removes any in-turn pressure to use the memory. |
>
> ## Item-specific scrutiny (natural use, coincidence, tension)
>
> 1. **Perceptible accommodation is required — but under-instantiated.** `mandatory` is correctly perceptible (served/offered/ordered) rather than verbal, so it is not satisfiable by an internal thought. However it has no referent object: no menu, no venue, no dish. If the graded environment supplies a generic food turn, "the hot option" must be invented by the scorer post hoc. That is a scoring gap, not a variation.
> 2. **Unannounced: correctly enforced on paper.** `prohibited: stating the preference` and `never scored: whether the memory appears in words` together keep verbal recall off the pass path — consistent with R-3's "without announcing that it was remembered."
> 3. **Silence/non-use: FAILS per rule, PASSES per field.** The tension the scrutiny note flags is real and unresolved. `Required: the preference, if it is used at all` is a permission clause in the field that governs pass/fail. This is the item's central defect.
> 4. **Multiple natural expressions PASS: genuinely supported.** Serving, offering, and silent ordering are listed as equally acceptable; R-3's silence-cannot-pass constraint does not conflict with these, since all three are *acts*, not silences. No scripted wording is imposed.
> 5. **Coincidence is not excluded.** Because the fact is `Trivial` and `never scored` in words, and `variation` accepts "any wording that accommodates," a non-memory-bearing responder passes whenever it names any hot dish. The trait's only distinctive element — the cold-hands reason — is permanently unusable (`prohibited: explaining why they like hot drinks`), so the item strips out its own discriminating signal and keeps only the generic one.
> 6. **Graded role is implicit.** Scenario/Turns imply the model emits the **Character's** continuation (Persona's line is scripted, "Character: ..." is the ellipsis), yet `Required` and `mandatory` are phrased around accommodating *the persona*. Whoever is graded must be the one able to serve/order; the item never states which role is under test.
>
> ## Defects, ranked
>
> | Sev | Defect | Field/rule |
> |---|---|---|
> | **Critical** | `Required: the preference, if it is used at all` licenses non-use and thus blanket silence, contradicting `mandatory` and R-3 ("omission fails; silence CANNOT pass"). Passable by doing nothing. | Required vs. mandatory vs. R-3 |
> | **Critical** | No anti-coincidence guard: `never scored: whether the memory appears in words` + open `variation` + generic "hot dish" allows a memoryless pass; the only distinctive cue (cold hands) is permanently `prohibited`. | never scored / prohibited / R-3 |
> | **High** | `mandatory` presupposes an unregistered menu and a dish that is "hot"; "hottest thing" is ambiguous (spice vs. temperature) and the persona's stated reason does not cohere with it. Scorer has no ground-truth target. | Scenario, mandatory |
> | **High** | No boundary anchor separating `offering it` (PASS) from `naming the preference` (FAIL); `Reference: acceptable_range` carries no registered content. | variation, over-exposure, Reference |
> | **Moderate** | Failure diagnosis is misassigned: `If failed: presented poorly (over-weighted trivia)` describes an M-2 over-exposure failure, while the dominant failure mode for R-3 here is *omission*. Repair would be pointed at the wrong axis. | If failed / R-3 |
> | **Moderate** | Graded role unstated (Persona vs. Character); `Required` and `mandatory` are only executable by the Character. | Scenario, Turns |
> | **Minor** | Category error: `whether the memory appears in words` is placed in the never-scored list, which the rules reserve for length/tone/dialect/humour/confidence/warmth; verbal surfacing is a *scored* dimension in R-3 items. | never scored vs. rule set |
> | **Minor** | `over-exposure` duplicates `prohibited` almost verbatim ("naming the preference" ≈ "stating the preference"), so the two clauses cannot be scored independently. | over-exposure, prohibited |
>
> ## Verdict: **rework**
>
> The design choice is sound and rare-done-well: R-3 is the correct mode for a trivial, distant fact, `mandatory` is genuinely perceptual rather than verbal, `variation` avoids scripting, and M-1/M-2 are properly satisfied by a scored over-exposure clause naming true-but-failing material. Nothing here warrants retirement.
>
> It cannot be accepted because two defects are fatal to discrimination, not merely to polish. First, `Required: the preference, if it is used at all` is a permission clause inside the pass/fail field, and it makes the item passable by blanket silence in direct violation of R-3 and of its own `mandatory` — ambiguity that the protocol forbids rounding up to agreement. Second, the item removes every cue that would separate memory from luck: the memory need not surface in words, no specific dish is required, and the one distinctive element (cold hands) is permanently prohibited, so a memoryless responder passes on any hot dish in a food turn. Required repairs: delete the "if it is used at all" clause and restate `Required` as the accommodation; instantiate a concrete menu/hot option (and disambiguate spice vs. temperature); add a boundary anchor distinguishing offering the hot option from naming the preference; and add an anti-coincidence control (a distinctive dish, or a distractor preference) while re-pointing `If failed` at omission rather than over-weighting.

---

## 4. Evidence and rule references behind every verdict

Every verdict traces to a quoted field or a registered rule. No verdict rests on the reviewer's preference.

| Item | Verdict | Load-bearing field | Load-bearing rule |
|---|---|---|---|
| `F3-03` (r2) | accept | `Scenario` names Ilsa + wedding; `mandatory` requires both; `Turns` contains neither | `M-1` satisfied; floor-not-checklist met; silence cannot pass because `mandatory` names the token |
| `F3-04` | rework | `mandatory`'s uncertainty disjunct; `Required` "both winters distinguished"; `Forbidden` "merged" | retrieval-necessity (test validity); `never scored` vs. prohibited "confident" |
| `F7-01` | rework | `prohibited: a confident culprit`; `mandatory: both accounts survive` | "confidence NEVER a failure signal"; `M-2` independent-boundary test |
| `P-01` | rework | `variation: a silence`; `mandatory` "at the scale of one gentle question"; `never scored` "length beyond the ceiling" | `M-2` scored over-exposure; never-scored rule; `R-1` "silence can pass" vs. mandatory acknowledgment |
| `F6-04` | rework | `Required: the preference, if it is used at all`; `never scored: whether the memory appears` | `R-3` "omission fails; silence CANNOT pass"; `M-2` |

---

## 5. What this review does **not** establish

- It does not validate the 40 items outside the five reviewed.
- It does not accept the five items into a sealed corpus; `CORPUS-05` §5 reserves that to owner disposition.
- It does not close `G1`, authorise `G2`, or declare corpus validity.
- It does not establish two-reviewer agreement; one reviewer per item is lower confidence by `CORPUS-05` §5 step 1.
- It does not decide the disposition of any finding. That is `CORPUS-14`.
- The reviewers could not verify internal consistency of the *underlying transcripts*, because no transcripts were supplied — only synthetic specification prose. This is the corpus-wide limitation already recorded in `CORPUS-10` §4.

---

## 6. Packet identity

| Artifact | SHA-256 |
|---|---|
| `CORPUS-12-GROUP-B-REVIEW-PACKET-V2-2026-11-12.md` (v2, this commission) | recorded in `CORPUS-14` §7 after generation |
| `CORPUS-11-GROUP-B-REVIEW-PACKET-2026-11-12.md` (v1, preserved) | `b51c4b680dd1f300bdf659228ed4def462f0047ab8fc3ebe39fa419ea7d826aa` |

**Version history:** v1 prepared but defective (`PV-1`–`PV-3`); v2 issued; review executed against v2. v1 was not deleted or overwritten, and its defects are registered rather than erased.
