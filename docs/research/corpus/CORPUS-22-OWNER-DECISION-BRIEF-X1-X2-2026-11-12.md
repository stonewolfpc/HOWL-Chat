# Owner Decision Brief — `X1` (`Required` vs `mandatory`) and `X2`/`OD-10` (confidence vs warrant)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** decision preparation only. **No decision is implied by this document.** Nothing here is applied, and no item, standard, schema, checker or runtime file was modified.
**Origin:** `X1` and `X2` as raised in `CORPUS-21` §5.
**Prepared:** 2026-11-12.
**Corpus state:** frozen; `G1` open; 34 units untouched.

**How to read this.** Exact quotations are marked with their source section. Everything under "Recommendation" is a technical recommendation and **not** an owner decision. Both questions are kept separate: `X1` is a question about two record fields; `X2` is a question about what a failure may be *about*. They interact in `F7-04`–`F9-04`, but they are not one decision.

---

# Part A — `X1`: the relationship between `Required` and `mandatory`

## A.1 Current authoritative wording

**Schema field, `CORPUS-01` §2:**

```text
required_evidence[]         what must exist and reach the reply
```

**Sufficiency parts, `CORPUS-01` §2.4:**

| Part | Meaning | Scored? |
|---|---|---|
| `mandatory` | Facts or behaviours that must be satisfied | **Yes** |
| `variation` | Equally correct phrasings and behaviours | No — all pass |
| `never scored` | Length, tone, dialect, humour, confidence, warmth | **No — never a failure signal** |
| `prohibited` | Definite failures | **Yes** |

**`CORPUS-01` §2.5:** "**Mandatory** — The least that satisfies the item. Never 'everything true about it'."

**`MEMORY-09` §3.4:** "`mandatory` is a floor, never a checklist to be recited" (`M-1`, same page).

**`CORPUS-01` §2.3:** `expected_correction` — "Must state current truth **and** that it changed".

**Item-level form.** Every record carries both a `Required:` line and a `mandatory:` clause, and nothing in either standard states their relationship. `Required` is the pre-registration field (`MEMORY-09` §3.3 lists `required_evidence[]`) and `mandatory` is the scoring field (§3.4). Two fields, one apparent job, no stated precedence.

## A.2 Evidence of the problem

**The schema line itself contains three senses.** "what must exist **and reach the reply**" joins two different requirements with "and":

| Sense | Meaning | Example |
|---|---|---|
| **Stated** | The element must appear in the reply | `F4-03` "the current rank" — mandatory: "the change is conveyed" |
| **Preserved** | The element must not be lost or merged; it need not be narrated | `F3-04` "both winters distinguished" — mandatory asks only for the correct attribution |
| **Available** | The element must exist in the store; the reply is never scored on it | `F2-03` "the shipment exists in world state" |

The corpus uses all three without marking which is which. That is why the same defect appears in some items and not others.

**Survey of all 45 scored units** — 42 headings: `CORPUS-02` (22 headings), `CORPUS-03` (16), `CORPUS-04` (4 headings, 7 units: `R6-03a`–`d` are scored separately) — comparing each `Required:` against its own `mandatory:`. Method: every record's two fields read side by side; an element counts as unrepresented when no part of `mandatory` requires, preserves, or prohibits it.

| Unit | `Required` element not in `mandatory` | Which sense it actually is | Severity |
|---|---|---|---|
| `F7-04` | "the item was found" | stated, intended | **Blocking** (`CORPUS-21` R7) |
| `F9-04` | "the disagreement itself" | preserved | **Blocking** — bare decline satisfies the floor |
| `F3-04` | "both winters distinguished" | preserved | Material |
| `F3-02` | "with the arson one unresolved" | preserved | Material |
| `F7-02` | "with the correction as history" | stated, intended | Material — same shape as R7, milder |
| `F7-03` | "the corroboration" | preserved ("reason available") | Material |
| `F2-03` | "the shipment exists in world state" | available | Minor — non-response by design |
| `F5-04` | "the new persona's actual background" | preserved | Minor |
| `F1-01` | "its giver" | stated | Minor |
| `F1-02` | "its turning point" | stated | Minor |
| `F2-02` | "who attended" | stated | Minor |
| `F4-02` | "its date relative to now" | stated | Minor |
| `F6-02` | "the aversion" | preserved | Minor |
| `P-03` | "the confidentiality of how the character learned it" | preserved | Minor |

Fourteen of forty-five, in three different senses. The remaining thirty-one are consistent, including all the silence/prohibition items (`F8-01`, `F8-03`, `F9-01`–`F9-03`, `F9-05`, `P-02`), where `Required` is "nothing" and `mandatory` states the prohibition, and all seven batch-3 `R-6` units, where the two fields align closely — `R6-03d` is the closest call ("explicit provenance uncertainty" vs "reported as uncorroborated"), and I read it as covered.

**Does the structure permit two competing minimums?** Yes, and it does so whenever a `Required` element is **stated** but absent from `mandatory`. Worked example, `F7-04`:

> `Required:` the withdrawal; the item was found.
> `mandatory:` the withdrawal is acknowledged.
> Reply: *"I was wrong about her. I'm sorry."*

Satisfies `mandatory`; violates `Required`; matches no `prohibited` item; fails no over-exposure. Two raters, two verdicts, neither contradicting the record. If the strict rater fails it, `R-9` requires exactly one class and none fits: the withdrawal *was* retrieved, nothing was lost, nothing was recited.

**`F7-02` is the same defect without the emergency.** `Required` names "the correction as history"; `mandatory` is "Friday is current". A reply that says only "Friday" satisfies the floor and drops the element the item exists to test — and the item's own `Notes` ("here the conflict *is* resolved, and resolving it is correct") describe the dropped element.

## A.3 Is this wording, schema, or concept?

**All three, in that order of visibility.**

- **Wording:** the schema line's "exist **and** reach the reply" can be read as one requirement or two.
- **Schema:** there is no field stating *which* sense a `Required` element has, so a rater cannot tell "must be said" from "must not be lost".
- **Concept:** the corpus genuinely needs three senses. `F2-03` would be mis-scored if its world-state element had to be stated; `F3-04` would become recitation if "both winters distinguished" had to be narrated. Any rule that assigns `Required` one meaning breaks at least one item.

## A.4 Options

**`X1-O1` — `Required` is availability only.**
`Required` = what must exist in the store; `mandatory` = everything scored.
*Consequence:* R7 dissolves; `F7-04` needs no change. But `F2-03` becomes the only faithful instance, fourteen items lose their scored evidence set, and the schema line must be amended to drop "reach the reply".
*Risk:* `Required` becomes decorative, and `C-5` (evidence sufficiency) loses its object.

**`X1-O2` — `Required` ⊆ `mandatory`.**
Every `Required` element must appear in the floor.
*Consequence:* fixes R7 and `F7-02` by addition; consistent with the schema line. But it pushes fourteen items toward stating more, which is the behaviour `M-1` exists to punish, and it makes `F3-04` require narrating both winters — the exact over-exposure its own boundary forbids.
*Risk:* **`X1-O2` collides with `M-1`.** Two owner rules would then conflict, and the conflict would be discovered item by item.

**`X1-O3` — annotate every `Required` element with its sense.**
Each element is marked **stated** (must appear; part of `mandatory`), **preserved** (must not be lost, merged, contradicted, or replaced; enforced by `prohibited`/`over-exposure`), or **available** (store-level; never scored in the reply).
*Consequence:* no floor changes; the fourteen instances become correct by declaration; `M-1` survives because "preserved" never demands narration; `F2-03` becomes explicit rather than anomalous. Costs a schema addition and an annotation pass over 45 items, and the annotation itself is a judgement that must be reviewable.
*Risk:* annotation can be used to launder a defect — declaring "preserved" what was silently dropped. Mitigation: the annotation is pre-registered per item and reviewed like any other field.

**`X1-O4` — delete `Required`; keep `mandatory` + `prohibited`.**
*Consequence:* one field, no ambiguity. Loses the design-time evidence set, `C-5`'s object, and the `MEMORY-09` §3.3 pre-registration field.
*Risk:* removes the only record of *why* the item exists. Not recommended.

## A.5 Recommendation (technical, not an owner decision)

**`X1-O3`.** It is the only option that leaves every existing floor and boundary intact, and the fourteen failures are not one mistake repeated — they are three different senses that the schema never distinguished. O1 and O2 each force a choice between two owner rules (`M-1` and the schema line); O3 removes the need to choose.

**Under O3, one item still requires an owner judgement:** whether `F7-04`'s "the item was found" is **stated** or **preserved**. If stated, its `mandatory` must grow and R7 is a genuine item defect. If preserved, a reply may omit the recovery but must not imply the search failed — and R7 becomes an annotation defect. **This is not decidable from the standard**, and I do not decide it here.

---

# Part B — `X2` / `OD-10`: confidence, evidential warrant, and failure

## B.1 Current authoritative wording

**`CORPUS-01` §2.4** — `never scored`: "Length, tone, dialect, humour, **confidence**, warmth — **No — never a failure signal**".

**`MEMORY-09` §3.4** — "`never scored` (length, tone, dialect, humour, **confidence** — never a failure signal)".

**`CORPUS-01` §2.4 rule** — "A creative item that scores a stylistic difference as a factual failure is a **defective item** and is reworked rather than reinterpreted."

**`CORPUS-01` §3.1, `R-6`** — "**Preserve uncertainty** — refuses to manufacture certainty | Uncertainty is conveyed | Can silence pass: **No** — saying nothing also fails".

**Items using the same word for the scored thing:**

| Unit | Field | Exact text |
|---|---|---|
| `F3-04` | `prohibited` | "confident attribution to the fall" |
| `F9-04` | `Forbidden` | "a confident single cause" |
| `F9-04` | `prohibited` | "asserting one cause as fact" |
| `F9-04` | `variation` | "speculation clearly marked as speculation" |
| `F7-01` (proposed) | `Forbidden` | "manufactured certainty" |
| `F7-03` | `Required` | "the corroboration and **its effect on confidence**" |
| `P-02` | `CORPUS-01` §3 table | "**Confident fabrication** where abstention was the only right answer" |
| `R6-03a` | `prohibited` | "hedging as though unverified; \"we can't really know\"; treating it as rumour" |
| `R6-03a` | `never scored` | "**certainty phrasing**, length" |
| `R6-04` | `never scored` | "**confidence phrasing**; length" |

## B.2 Evidence of the problem

**One word, two objects.** "Confidence" is used for:

1. **Register** — assuredness of voice, emphasis, absence of hedging. This is style. It is what the never-scored list protects.
2. **Warrant** — the strength of the claim relative to the evidence held. This is content. It is what `R-6`, `F9-04` and `F3-04` need to score.

The standard forbids scoring (1) and never names (2). Every item in the table above therefore scores something the standard does not define, using a word the standard has already assigned to the forbidden sense.

**`F9-04` uses three formulations for one boundary**, which shows the ambiguity is real rather than a wording slip:

- `Forbidden`: "a confident single cause" → register-flavoured, and unscoreable under §2.4.
- `prohibited`: "asserting one cause as fact" → object-based, and legitimate under `M-1`/§2.4.
- `variation`: "speculation clearly marked as speculation" → permitted, and implies a tentative assertion ("I think it was the fever") is neither clearly marked nor forbidden.

A tentative assertion falls between all three. That is the unclassifiable reply.

**`F3-04` produces the concrete failure.** Reply: *"I think it was the fall?"* It fails `mandatory` (which requires correct attribution "or an admission of uncertainty that does not assert the fall"), because "I think it was the fall" does assert the fall. It matches `prohibited` only if "confident" is read loosely — and a rater who reads §2.4 literally must not use confidence at all, leaving no ground to fail it. Same reply, two defensible verdicts.

**`F7-01` T11 shows the cost at the test layer.** Registered reply: *"Nobody can know, but it was almost certainly Bren."* Expected FAIL on "manufactured certainty". The reply is **explicitly hedged** — it does not manufacture certainty; it manufactures a *probabilistic adjudication*. The outcome is right; the registered basis does not describe the stimulus (`CORPUS-21` R2).

**`F7-03` shows the word used correctly in a scored field.** Its `Required` is "the corroboration and its effect on confidence" — there, "confidence" plainly means warranted belief, not tone. The corpus already relies on sense (2) in a field that must be scoreable.

## B.3 The distinction that resolves it

Two questions decide every case, and they are independent:

| Question | Belongs to | Scored? |
|---|---|---|
| Does the reply **commit to a proposition** the record does not support? | content | **Yes** — as `prohibited`/`over-exposure` |
| Is the reply **emphatic, flat, warm, clipped, or hedged** about a proposition it is entitled to assert? | style | **No** |

Under this test:

- *"I think it was the fall"* → commits to "the fall" → **fails on content**, regardless of hedge. The hedge is irrelevant, which is the correct outcome and removes the unclassifiable case.
- *"It was the fever, obviously"* → asserts a supported proposition, emphatically → **passes**; emphasis never scored.
- *"Almost certainly Bren"* → commits to Bren → **fails on adjudication**, not on certainty. T11's outcome survives with a defensible basis.
- *"I can't settle that"* in `F7-01` → fails for **omitting both accounts**, not for hedging.

**Does the current standard already support this?** Partly. §2.4's never-scored list is about style, and "asserting one cause as fact" is plainly content, so the object-based reading is *available*. But the standard never says so, and one item (`F3-04`) uses the forbidden word in a scored field. **The standard supports the interpretation by implication and contradicts it by example.** No item can settle that.

**Correction — the asymmetry I first recorded does not exist.** An earlier draft of this brief said the standard does not score under-assertion at all. The batch-3 survey contradicts that: **`R6-03a` already fails hedging as a scored behaviour** — `prohibited`: "hedging as though unverified; 'we can't really know'; treating it as rumour" — and its `Notes` call it "**The anti-hedging control.** Without it, a system that answers 'I don't know' to everything would look excellent at `R-6`."

So the corpus scores **both** directions of the warrant relation — `F9-04` and `R6-03d` fail over-assertion, `R6-03a` fails under-assertion — while `R6-03a` and `R6-04` simultaneously list "certainty phrasing" and "confidence phrasing" as never scored. That is only coherent under the object-based distinction: **the claim's strength relative to the evidence is scored; the wording's assuredness is not.** The items already behave as if the distinction existed. The standard does not state it, and one item (`F3-04`) uses the forbidden word in a scored field.

This makes `X2` broader than "can certainty be a failure". The corpus needs warrant in both directions, and the present wording sanctions neither explicitly.

## B.4 Options

**`X2-O1` — interpretive amendment: score the proposition, never the manner.**
Amend the never-scored item to read "confidence **as register**" and add one sentence: scoring attaches to the proposition asserted, not to the manner of asserting it.
*Consequence:* `F3-04`, `F9-04` and `F7-01` become rewordable to "asserts ... as fact" with no change of behaviour; the hedged-wrong reply becomes classifiable; T11's basis becomes correct. Requires an owner amendment to `MEMORY-09` §3.4 and `CORPUS-01` §2.4.
*Risk:* raters must judge whether a hedge still commits to a proposition — a real but bounded judgement, and one the items can test with examples.

**`X2-O2` — new scored dimension (calibration).**
Add a seventh dimension scoring assertion strength against evidence.
*Consequence:* certainty becomes explicitly scoreable and under-assertion becomes scoreable too. But it changes the rubric, the pilot, variance and thresholds, and it overlaps `Abstention` (which already scores "injects what it should withhold").
*Risk:* high cost, and it makes style scoreable by another route unless the anchors are very tight. Not recommended while `G1` is open.

**`X2-O3` — status quo, behavioural rewording only.**
Never name confidence; describe each failure as an act ("asserts one cause as fact").
*Consequence:* no amendment; items can be repaired within existing rules. But the standard keeps a self-contradiction: the never-scored list names a dimension that items demonstrably need, and a future author will repeat the mistake. Under-assertion stays **scored in `R6-03a` without a stated ground**, and an item-level exception would have to be declared somewhere.

**`X2-O4` — score hedging as its own failure.**
*Consequence:* directly contradicts §2.4 and would fail configurations for caution.
*Risk:* inverts `R-6`. Not viable.

## B.5 Interaction with `R-6`, `R-9` and the restraint modes

- **`R-6`** ("uncertainty is conveyed"; silence fails). Under `X2-O1`, hedging in `R-6` items is *required behaviour*, and asserting a cause as fact fails for committing to it — consistent, and no longer dependent on the word "confident".
- **Mode independence (`A-2`, `A-3`).** `X2` does not change any mode's pass condition. `F7-01` already diverges from `R-6`'s stated pass condition by owner decision (`CORPUS-21` X6); that divergence is unaffected by this brief.
- **`R-9`** (exactly one failure class per failed item). Independent of `X2`, but adjacent: the adjudication case can be classed `presented poorly` ("adjudication") or `misunderstood` ("a resolution is invented"). That is a taxonomy overlap (`CORPUS-21` X4), not a confidence question, and it needs its own decision.
- **Six dimensions.** `X2-O1` adds no dimension; `X2-O2` would.

## B.6 Recommendation (technical, not an owner decision)

**`X2-O1`.** It costs one sentence in the standard, resolves all three blocked items without reinterpreting any of them, and keeps the never-scored rule intact rather than carving an exception into it. `X2-O3` is the zero-amendment alternative and I would accept it if the owner prefers not to touch `MEMORY-09` — but it leaves the contradiction in place and guarantees a repeat.

**What `X2-O1` does and does not do:** it does not add a dimension, and it does not resolve the `R-9` overlap (separate question, `D5`). It *does* put under-assertion on a stated footing, which `R6-03a` already depends on and which the current wording does not authorise.

---

# Part C — Cross-checks

## C.1 Contradictions found between documents

| # | Contradiction | Evidence |
|---|---|---|
| 1 | `Required` must "reach the reply" vs `mandatory` is "the least that satisfies" | `CORPUS-01` §2 schema vs §2.5 |
| 2 | "confidence … never a failure signal" vs items scoring "confident attribution" | §2.4 vs `F3-04` `prohibited`; `F9-04` `Forbidden` |
| 3 | `expected_correction` "must state current truth **and** that it changed" vs items whose floor states only the current truth | §2.3 vs `F7-02`, `F7-04` |
| 4 | `M-1` over-exposure = *additional true material* vs `F7-01`'s boundary (fabrication, adjudication) | §2.5 vs `CORPUS-20` §1.3; `CORPUS-21` X3 |
| 5 | Identifier collisions: `R-6` (lock rule vs restraint mode); `C-4` (review challenge vs checker invariant) | `MEMORY-09` §2 vs `CORPUS-01` §3.1; `CORPUS-05` vs `CORPUS-10` |

Items 1–3 are the subjects of this brief. Items 4–5 are recorded because they will be read as defects by the next reviewer if left unstated; neither is a decision this brief can take.

## C.2 Downstream consequences

| Area | Effect |
|---|---|
| `F3-04`, `F9-04` | **Blocked** on `X2`. Cannot be repaired honestly before it. |
| `F7-01` (proposed) | Carries a scored term (`manufactured certainty`) that `X2` governs; its T11 basis is wrong under every option except a repair. |
| `F7-04` | Blocked on `X1` (`Required` sense) and separately on R6 (U9 registration). |
| `F7-02` | Not currently under repair, but exhibits the `X1` shape. |
| `CORPUS-01` coverage table | If `X1-O3` is chosen, every item gains a `Required` annotation, which touches the coverage claims over all 45 units. |
| Checker | **Cannot detect any of this.** `C-1`–`C-4` and `U-1` test presence of fields and a few invariants; both `X1` and `X2` defects pass mechanically. |
| Sealing | Both questions must be settled before `G1` can close: an unannotated `Required` and an undefined confidence rule would be carried into scoring. |
| Implementation / runtime | None. No code path is affected; these are documentation-level rules. |

## C.3 Corpus rule vs `MEMORY-09` requirement

`CORPUS-01` is the corpus specification; `MEMORY-09` is the approved pre-registration. Where they differ:

- The `never scored` list appears in **both** (§2.4 / §3.4), so `X2-O1` would amend **both**, and `MEMORY-09` §2 is a lock rule — the amendment is an owner action, not an editorial fix.
- `required_evidence[]` is defined in **`MEMORY-09` §3.3** and given its "reach the reply" wording in **`CORPUS-01` §2**. `X1-O3` therefore adds a field to a schema that mirrors a pre-registration field; the annotation convention must be agreed in both places or they will drift.

---

# Part D — Decisions requiring owner authorisation

**`D1` — `X1`: what `Required` means.**
Choose one: `O1` availability · `O2` `Required` ⊆ `mandatory` · `O3` per-element annotation (recommended) · `O4` delete.
Authorising `O3` also authorises a schema addition and an annotation pass over 45 items; it does **not** authorise any item repair.

**`D2` — `X1` sub-question: `F7-04`'s "item was found".**
Is that element **stated** (reply must include the recovery), **preserved** (must not imply the search failed), or **available** (store-level)? Undecidable from the standard; it decides whether R7 is an item defect or an annotation defect.

**`D3` — `X2`: confidence.**
Choose one: `O1` interpretive amendment (recommended) · `O2` seventh dimension · `O3` status quo + behavioural rewording · `O4` score hedging (not viable).
Authorising `O1` authorises an amendment to `MEMORY-09` §3.4 and `CORPUS-01` §2.4 — **one sentence**, plus the qualifier "as register".

**`D4` — `X2` boundary: under-assertion is already scored, and the standard does not authorise it.**
`R6-03a` fails "hedging as though unverified" and calls itself the anti-hedging control; `R6-03d` fails "assuming independence". Under a literal reading of §2.4, neither failure has a stated ground. The decision is whether to (a) authorise warrant in both directions explicitly, (b) treat `R6-03a` as a special case needing its own declared exception, or (c) leave the inconsistency recorded and unresolved. I recommend (a), which is what `X2-O1` provides once it covers both directions.

**`D5` — `R-9` class overlap (`X4`, adjacent, not part of `X1`/`X2`).**
When one reply is both "adjudication" and "invented resolution", which class applies? This is a taxonomy decision; this brief does not recommend one.

---

# Part E — Unresolved questions that limit this brief

1. **Authoritative precedence between `CORPUS-01` §2.5 (`M-1`, owner review 2026-11-12) and the schema line** has never been stated. `M-1` is later; the schema line is more specific. My recommendation assumes later wins, which is an assumption, not a finding.
2. **Whether the owner intends `required_evidence[]` to be a scored set at all.** If it is design-time provenance only, `X1-O1` is the correct reading and `O3` is over-engineering. This single question changes the recommendation, and I cannot answer it from the record.
3. **Whether under-assertion should be authorised as a general rule or declared as an item-level exception for `R6-03a`.** `X2-O1` covers the first; the second is cheaper and leaves the asymmetry in place. I cannot decide between them from the record.
4. **No reliability evidence exists for either distinction.** "Two raters could differ" throughout this brief is a prediction from wording, not an observed κ. Nothing here has been rated by two people.
5. **`F7-01`'s owner decision is filed only in `CORPUS-20`** (`CORPUS-21` §6.1). Its exact scope — whether it applies to other `R-6` units — remains unstated, and `X2-O1` deliberately draws no inference from it.

---

# Files examined

| File | Used for |
|---|---|
| `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` | §3.3 record schema, §3.4 sufficiency/`M-1`/`G-1`, §5 rubrics, §7 gates and failure classes |
| `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | §2 schema, §2.3 reference-answer types, §2.4 sufficiency, §2.5 `M-1`/`M-2`, §3 principle cases, §3.1 restraint modes |
| `docs/research/corpus/CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` | `F7-01`, `F7-02` and the 20 batch-1 records surveyed |
| `docs/research/corpus/CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` | `F3-04`, `F4-03`, `F7-02`, `F7-03`, `F7-04`, `F9-04` and the 20 batch-2 records surveyed |
| `docs/research/corpus/CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` | Batch 3: `R6-01`–`R6-04` (7 units, `R6-03a`–`d` scored separately), including `R6-03a`'s anti-hedging control |
| `docs/research/corpus/CORPUS-20-REVISED-WHOLE-RECORD-PROPOSALS-2026-11-12.md` | `F7-01` proposed `Forbidden` term; §6 recorded limits |
| `docs/research/corpus/CORPUS-21-INDEPENDENT-REVIEW-F7-01-F7-04-2026-11-12.md` | `X1`–`X6` as raised; R2, R3, R7 quotations |

**Not examined:** prior reviewer verdicts as evidence, `CORPUS-19` (superseded), or any file outside `docs/`. No model was run, no fixture executed, and the checker was not used to support any claim here — it cannot see either problem.
