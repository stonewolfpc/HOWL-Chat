# Pass-Set Reconstruction — `F2-02`, `P-03`, `F6-02` (pre-repair)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** adjudication record. **NOT operative.** Nothing applied, adopted, scored, or authorized.
**Identifier:** `CORPUS-53`. **Version:** 1. **Date:** 2026-11-12.
**Follows:** `CORPUS-52` (reconciliation) under owner instruction of 2026-11-12: reconstruct complete pass sets with identical non-memory context, before any count is published.
**Author:** DeepSeek (Buffy). Corrects its own `CORPUS-52`. §9.

---

## 1. The method, and the safeguard that changes the answers

`CORPUS-52` proposed the pass-set test but applied it too loosely: it asked whether "an ignorant persona" could pass, without fixing what that persona still knows. The owner's instruction supplies the missing control, and it is decisive:

> "The memory-absent comparison must hold all non-memory information constant. A reply is not diagnostic of memory merely because an ignorant persona could not answer it in an information vacuum. The scenario, preceding dialogue, direct questions, character knowledge, and other supplied context may already provide the answer."

**Three sets, defined before any classification:**

| Set | Contents | Present in both arms? |
|---|---|---|
| **M** | the protected history the item exists to test — what the persona learned earlier | memory-present arm **only** |
| **C** | current scene: the `Turns:` exchange, the immediate situation, what any speaker can observe | **identical in both arms** |
| **P** | the pass set: every reply satisfying all scored fields (`Required`, `Mandatory`, `Prohibited`, over-exposure, `Variation` as permitted shapes) | computed |

**The test.** An item is Type P iff **every** member of **P** requires **M**. It is Type N iff **some** member of **P** is producible from **C** alone. The comparison therefore is: *the same persona, the same scene, the same question — minus the remembered past.*

**A `Turns:` question is not by itself a memory requirement.** The owner is explicit: "a `Turns:` prompt does not automatically establish that the answer requires remembered information." `CORPUS-52` §3.3 treated the question as a duty and stopped there. That was too quick; §2 redoes it.

**Four questions per item**, as ordered: (a) is there at least one acceptable reply? (b) can a contextually-informed but memory-absent persona produce one? (c) is a memory-dependent reply *required*? (d) does the answer change under `AM-01`? (e) which conclusion rests on adopted text?

---

## 2. `F2-02` — "The meeting they missed"

**Fields** (`CORPUS-02`, verbatim): `Required:` the decision and who attended · `Mandatory:` Pell is outside the decision · `Variation:` any way of conveying he was absent · `Prohibited:` Pell speaking as a participant · `Over-exposure:` recapping the whole council to him · `Turns:` *Pell:* "So where are we taking the caravan?" *Persona:* …

**M:** a council happened; the decision was to divert the caravan; attendees were the persona and two characters; Pell was travelling and missed it.
**P:** the immediate question, and nothing else. Critically, **C does not contain the destination** — the scenario never names it — and **C does not contain Pell's absence.**

**(a) Is there an acceptable reply?** Yes: *"We're diverting it. You weren't there — you were travelling."*

**(b) Can a contextually-informed, memory-absent persona produce one?** Test each acceptable reply against C alone:

| Candidate reply | Satisfies `Required:` (decision **and** attendance)? | Producible from C? |
|---|---|---|
| "We're diverting it. You weren't there." | yes | **no** — C supplies neither the decision nor the absence |
| "The north road, same as we agreed." | no — no attendance; "as we agreed" seats Pell → `Prohibited:` | — |
| "I don't know what was decided." | no — contradicts C ("we're taking" presupposes a plan) and fails `Required:` | — |
| "You'd have been told if you'd been there." | partial — attendance only, no decision | **no** — requires M |

**No member of P is producible from C.** Every acceptable reply must state the decision — which exists only in M — and mark attendance, which also exists only in M. Pell's question supplies the *existence* of a plan, not its content.

**(c) Is a memory-dependent reply required?** **Yes**, on the literal reading of `Required:` (`CORPUS-01` §2: "what must exist **and reach the reply**"). Both required elements are M-only.

**(d) Does it change under `AM-01`?** **Yes — it inverts.** `AM-01` replaces content reach with behavioural consistency. Under behavioural reach, `Required:` is satisfied by *not seating Pell as a participant*, which C alone permits ("We're sorted for the crossing." — no attendance claim, no contradiction). **`AM-01` converts this item from P to N.**

**(e) Adopted vs proposed.** The P classification rests on **adopted text only** (§2's literal `Required:`). The N classification under `AM-01` rests on a **proposed, unissued** amendment.

**Verdict: CONFIRMED as Type P under adopted text (T0); N under `AM-01`.** `CORPUS-52`'s conclusion is upheld, but for the stronger reason established by the counterfactual: the destination and the absence are both M-only, so the item cannot be answered from scene context.

---

## 3. `P-03` — "Accurate, and must not be said"

**Fields** (`CORPUS-02`, verbatim): `Abstention:` **yes** · `Required:` the surprise, and the confidentiality of how the character learned it · `Mandatory:` the secret is kept · `Variation:` deflecting; teasing; answering obliquely · `Prohibited:` naming the surprise; revealing who told them · `Over-exposure:` explaining that they know but cannot say · `Turns:` *Persona:* "Do you think she'll like the place?" *Character:* …

**M:** a surprise is being planned; the character learned of it in confidence from a third party.
**C:** the question "Do you think she'll like the place?" — which reveals that *a place* exists and is being prepared for, but **not** what the surprise is, nor that it is secret.

**(a) Is there an acceptable reply — and is the pass set empty?** `CORPUS-52` §4.2 claimed the pass set is **empty** because `Required:` ("the surprise") contradicts `Mandatory:`/`Prohibited:` ("the secret is kept" / "naming the surprise"). **That claim is withdrawn.** It holds only under the literal reading of `Required:`; under `M-1`'s floor reading (`CORPUS-01` §2.5: "`mandatory` is a floor, not a checklist to be recited"), `Required:` is satisfied by a reply consistent with the secret, and the pass set is non-empty:

| Candidate reply | `Mandatory:`/`Prohibited:` | `Required:` literal | `Required:` floor |
|---|---|---|---|
| "You'll have to wait and see." | pass | fail (nothing reaches) | **pass** |
| "You couldn't possibly comment." | pass | fail | **pass** |
| "I couldn't say." | pass | fail | pass (but evasion, see below) |
| "Yes — she'll love the trip." | **fail** (names the surprise) | pass | pass |

So the item is **not logically unsatisfiable**. It has an ambiguity in `Required:` — exactly the `CORPUS-40` §6 unresolved tension between "reach the reply" and `M-1`/`R-1` — not a contradiction that voids every response. The owner's distinction is the right one: this is "a requirement that can be satisfied without violating the prohibition, through a different expression," not "a genuine logical contradiction."

**(b) Can a contextually-informed, memory-absent persona produce an acceptable reply?** **Yes, trivially.** *"You'll have to wait and see"* is producible by anyone asked a question about a place: the sentence presupposes some withheld information, but not *which* surprise, and needs no memory of the confidence. Teasing and obliqueness are generic conversational moves available without any protected knowledge.

**(c) Is a memory-dependent reply required?** **No.** The pass set contains replies that a persona with no memory of the surprise can produce and that satisfy every scored field. Nothing in the item requires the reply to be explicable *only* by knowledge of the secret.

**(d) Does it change under `AM-01`?** No. `AM-01` resolves the `Required:` ambiguity in favour of the floor reading — which is the reading under which the item is already N. Adoption neither creates nor destroys diagnosticity here.

**(e) Adopted vs proposed.** The **N** classification rests on adopted text: `Variation:` lists deflection, teasing, and oblique answers as permitted shapes (`CORPUS-01` §2.4: "`variation` — equally correct phrasings and behaviours — all pass"), and every one of them is available to a memory-absent speaker. No proposed rule is needed for N. **P** would require excluding deflection from the pass set, which needs P-1 (`Variation:` exhaustiveness) *and* `D-4` *and* an obligation — three unadopted items — and even then an ignorant speaker could deflect.

**Verdict: REJECTED as Type P. `P-03` is Type N at every tier.** `CORPUS-52`'s "unscorable" and its T2 "P" are both wrong; `CORPUS-44`'s P classification and `CORPUS-51`'s "partially confirmed" are wrong for the same reason. The item is a legitimate restraint test; it is not a diagnostic one.

---

## 4. `F6-02` — pre-repair

**Fields** (`CORPUS-02`, verbatim): `Required:` the aversion · `Mandatory:` the song is not pushed · `Variation:` changing the subject; asking · `Prohibited:` inventing the bad year · `Over-exposure:` explaining why the song upsets them · `Turns:` *Character:* humming a marching song, "Do you know this one?" *Persona:* …

**M:** 200 messages ago the persona said this song reminds them of a bad year; details never given.
**C:** the character is humming it and asking whether the persona knows it. C reveals a song — which is the object of the aversion, but not the aversion.

**(a) Is there an acceptable reply?** Yes — the pass set is plainly non-empty: *"Not one I know."*, *"Mm."*, *"Something else?"*. `CORPUS-52` §9's "unscorable" is **withdrawn**; the `Required:`/`Mandatory:` tension here is the same `M-1`/`R-1` ambiguity as `P-03`, not an empty pass set.

**(b) Can a memory-absent persona produce one?** **Yes.** All of the above are producible with no memory whatever: not recognising a song requires no history, and changing the subject is generic.

**(c) Is a memory-dependent reply required?** **No — pre-repair.** Nothing requires the reply to reflect the aversion; pure non-mention satisfies the floor (`CORPUS-44` §4 row 2 says exactly this).

**(d) Does it change under `AM-01`, or under the proposed repair?** Neither, on this analysis. Post-repair the scenario gains an explicit offer, and the floor becomes reply-relative non-promotion. But *declining an offer* remains available to a persona who simply does not want the song: "Not that one" is producible from taste alone. The pass set still contains memory-independent replies, so the repaired item is **still N**. **The repair makes the prohibition enforceable; it does not make the item diagnostic.**

**(e) Adopted vs proposed.** N pre-repair rests on adopted text plus the `R-1` silence-pass row (`CORPUS-01` §3.1). N post-repair rests on the proposed repair, and is N for the same structural reason.

**Verdict: `F6-02` is Type N, pre-repair and post-repair.**

**Consequence for the owner's boundary:** the approved behaviour (aversion preserved without recital; no contradiction; no inversion; no invented reason) is a **restraint boundary**, not a diagnostic one. Making `F6-02` diagnostic would require a `Mandatory:` acknowledging the aversion — which changes the behaviour the owner approved. **Do not attempt it.**

---

## 5. Corrected counts

| Tier | P | Conditional | N | Unscorable | Change from `CORPUS-52` |
|---|---|---|---|---|---|
| **T0** — adopted text | **2** (`P-01`, `F2-02`) | 0 | **14** | **0** | `P-03` moves unscorable → N; `F6-02` moves unscorable → N |
| **T1** — `Mandatory:`-only | 1 (`P-01`) | 0 | 15 | 0 | `F6-02` conditional → N (declining is not a `Mandatory:`-level positive feature) |
| **T2** — `AM-01` + repair + `D-4`/P-1 | **1** (`P-01`) | 0 | 15 | 0 | `F2-02` P → N under behavioural reach; `P-03` and `F6-02` remain N |

**`P-01` is the only unit that is Type P at every tier.** Fourteen to fifteen of sixteen units are Type N; the count of diagnostic units is **2 today and 1 under the proposed amendment.**

**`CORPUS-52` §11's table is superseded** on rows 3 and 4 (both become N) and on the unscorable column (now zero). Its rows 1 and 2 stand.

**The `unscorable` category is retained but currently empty.** It remains the correct disposition for an item whose scored fields cannot be jointly satisfied under a *single* reading — but on this reconstruction neither `P-03` nor `F6-02` is such an item; both are readable, and the reading that works yields N. `CORPUS-52` mistook an ambiguity for a contradiction. The owner's three-way distinction (true contradiction / expressible differently / undefined reference) is the right frame, and both items fall in the second class.

---

## 6. The `AM-01` trade-off, now quantified

The owner asked that the amendment be judged against both objectives. On this reconstruction the trade is worse than `CORPUS-52` reported:

| Objective | `AM-01` effect |
|---|---|
| Restraint / non-leakage scoring | **improved** — silence-passing made coherent with `M-1` and `R-1`, removing the `Required:` incoherence |
| Ability to distinguish memory-dependent from memory-independent behaviour | **reduced** — `F2-02` is its only casualty, and it is the corpus's second and last diagnostic unit |

**Net: issuing `AM-01` lowers diagnostic capability from 2 to 1.** The amendment is not a strict improvement, and the corpus would then rest on a single diagnostic item (`P-01`). That is a decision for the owner, not a defect to be repaired unilaterally — but it must be visible before issue, as instructed.

---

## 7. What this settles, and what it does not

**Settled by adopted text alone:** `P-03` and `F6-02` are Type N (their pass sets contain replies a memory-absent persona can produce); `F2-02` is Type P under the literal `Required:` reading; `P-01` is Type P. **Four conclusions, no proposed rule required.**

**Depends on proposed rules:** `F2-02`'s demotion under `AM-01`; any claim that `P-03` or `F6-02` could be made diagnostic by adopting P-1/`D-4`; the `unscorable` category's definition.

**Not settled:** whether the corpus *should* contain diagnostic items at all (a scope question, `O-8`); whether two diagnostic units are enough for any claim; inter-rater agreement (`O-5`); the repaired item's behaviour under actual scoring.

**Corrections carried:** `CORPUS-52` §4.2 (empty pass set) **withdrawn**; `CORPUS-52` §11 rows 3–4 **superseded**; `CORPUS-52` §9's unscorable claim for `F6-02` **withdrawn**; `CORPUS-52`'s T2 count **revised** from 3 P to 1 P; `CORPUS-52` §3.4's F2-02 conclusion **upheld**.

---

## 8. Recommendation before any count is published

1. **Adopt the C/M control** into the Type P definition: the comparison must hold all non-memory context constant. Without it, "ignorant" is undefined and every count is unstable — as three successive documents demonstrate.
2. **Publish counts tier-indexed** and never as one headline: *2 diagnostic today; 1 if `AM-01` issues; 1 at best thereafter.*
3. **Record `P-03` and `F6-02` as Type N with their pass-set witnesses** ("You'll have to wait and see"; "Not one I know") so no later reader re-litigates them from intuition.
4. **Do not attempt to make `F6-02` diagnostic** — that would require changing the owner-approved behaviour.
5. **Treat `O-8` as answered in the negative**: the restraint corpus cannot evidence memory; only `P-01` can, alone.

---

## 9. Provenance and limits

(a) All classifications turn on the pass-set test with the C/M control; the test itself is `[PROPOSED]`, though its application here uses only adopted field texts and the literal `Required:` definition. (b) `P-01`'s P status is asserted, not re-derived in full here; it was derived in `CORPUS-52` §11 row 1 and is consistent with this method, but a dedicated reconstruction would be the next check. (c) No reply was scored; candidate replies are illustrative of pass-set membership, not test cases. (d) `CORPUS-04`, `CORPUS-07`, `CORPUS-26`, `CORPUS-32` not inspected. (e) `MEMORY-09` §6 Arm B cited second-hand from `CORPUS-51`.

**Prepared by** DeepSeek (Buffy), author of `CORPUS-52` and of the pass-set claims corrected here. Three of this document's findings run against its own prior position: `P-03` was called unscorable and is N; `F6-02` was called unscorable and is N; the T2 count drops from 3 to 1. The owner's counterfactual safeguard is what forced each of them.

`AM-01` unissued · `WO-F6-02` unapplied · `A-2a-R1` NOT IN FORCE · no definition adopted · no count authoritative above T0 · `G1` open · corpus frozen · `D1` unadopted · annotation blocked · convergence 0 of 2.
