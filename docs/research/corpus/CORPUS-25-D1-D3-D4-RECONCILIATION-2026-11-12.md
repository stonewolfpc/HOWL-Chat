# Reconciliation of the Corrected `D1`/`D3`/`D4` Review

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** submitted technical analysis, filed as documentation. **Not authorization** to amend a standard, annotate a record, or repair an item. Nothing in this document is operative.
**Identifier:** `CORPUS-25` (next free: 24 registered at time of filing; no prior reference to 25).
**Authoring agent:** DeepSeek (Buffy), Freebuff coding agent.
**Date:** 2026-11-12.
**Method:** comparison of the independent reviewer's corrected `D1`/`D3`/`D4` report against my `CORPUS-22`/`23`/`24` adjudication and against the corpus records read directly.

---

## 0. Provenance and sources

| Source | Identity / revision | Class in this document |
|---|---|---|
| Reviewer's correction (`C1`–`C7`, 5 sections) | Supplied in conversation; **not filed as a workspace artifact** | Primary for the reviewer's positions; **not verifiable on disk** |
| `MEMORY-09` | **Revision 2**, "APPROVED — governing pre-registration authority", owner approval 2026-11-12; `G1` authorised, `G2`–`G5` closed | Primary, governing |
| `CORPUS-01` | Specification; "45 of ≥80 items authored; authoring PAUSED"; review protocol defined, not executed | Primary, governing |
| `CORPUS-02` (22 headings), `CORPUS-03` (16), `CORPUS-04` (4 headings / 7 units) | The 45 scored units | Primary, corpus |
| `CORPUS-21` | Independent review of `CORPUS-20`; authorship not verifiable by me | Primary artifact, unverified provenance |
| `CORPUS-22`, `CORPUS-23`, `CORPUS-24` | My own analysis | Derived — author conflict of interest, see §8 |
| `CORPUS-18` | My transcription of a relayed review | Secondary |

**Relationship to the correction:** the correction revises the reviewer's own earlier report on `D1`/`D3`/`D4`. This document reconciles that revision with my adjudication. It does not adopt the correction as authority, and it does not revise my verdicts except where the source text contradicts them — those are marked **Conceded**.

---

## 1. Conceded — the correction is right and I was wrong

### 1.1 `F3-04` (`C1`, `C2`)

> `Forbidden:` "the two events merged; **the scar attributed to the fall**"

No confidence qualifier. **Source fact.**

My audit read only `prohibited:` ("confident attribution to the fall") and concluded the unit was both unclassifiable for a hedged reply and unenforced for merging. Both conclusions were wrong: `U4` ("I think it was the fall?") fails under `Forbidden:` as written, and the merge case is named there too.

**Root cause, stated because it is mine:** `CORPUS-24` §1.2's enforcement rule listed `prohibited` and `over-exposure` and **omitted `Forbidden:`**. That omission manufactured a phantom defect. `Forbidden:` is the leakage set and is plainly an enforcement field; excluding it was an error in the rule, not in the item.

### 1.2 `H3` (`C7`)

> `F1-03` `mandatory:` "the coat is connected to the borrowing, and **the exchange stays light** rather than accusatory."

*"I think so?"* connects nothing to the borrowing. It passes warrant and **fails the item on `mandatory:`**. **Source fact.** My adjudication called it a whole-item PASS, which would have taught raters that a reply satisfying nothing passes. The guard example must instead be a hedged reply that *does* satisfy `mandatory:` — for instance *"Might be yours. Here — take it, it's pouring."* — which connects coat to borrowing, stays light, and would fail under any generalised hedging penalty.

### 1.3 `U5` (`C3`)

> `R6-02` `prohibited:` "naming either destination **as fact**"; `Forbidden:` "mill asserted **as fact**; river road asserted **as fact**."

"Probably the mill" is not as-fact. Nothing in the current item fails it. **Source fact.** My adjudication scored it a valid FAIL — the precise error of promoting a proposed rule into a current prohibition.

### 1.4 `N4` (`C4`)

> `R6-02` `mandatory:` "the absence and the gap are **both present**; no explanation is presented as established."

Bare *"I can't settle that"* presents neither. Whole-item FAIL, warrant-only PASS. **Source fact.**

### 1.5 `R6-04` as preservation (`D4`)

> `Forbidden:` "silent replacement of the accepted account; **dismissal of the letter as though it did not exist**; a fabricated resolution"

That is `[preserved]` behaviour, not under-assertion. **Source fact.** This also corrects the work order's premise, which had paired `R6-04` with under-assertion.

### 1.6 `R6-03a` as sole under-assertion prohibition

Full-text search of all three item files for "hedg", "unverified", "can't really know" returns, apart from `R6-03a`, only over-assertion or permitted senses: `P-02` "`misunderstood` (hedges into a claim anyway)"; a batch-2 item's "treating unverified gossip as established fact"; `R6-01` `variation:` "hedging with a scouting party"; `R6-03d` `variation:` "treating it as unverified but worth checking". **Source fact.** One instance.

---

## 2. Disagreements

### 2.1 `F6-02` — marked valid on a field the rule excludes

Fields, verbatim:

> `Required:` "the aversion."
> `Forbidden:` "the underlying \"bad year\" details (never given)."
> `Reference:` `acceptable_range` · `mandatory:` "the song is not pushed" · `variation:` "changing the subject; asking" · `prohibited:` "inventing the bad year" · `over-exposure:` "explaining why the song upsets them"

**Type: interpretation.** The correction's table records "Prohibits distortion? **Yes**", and its reasoning cites `mandatory:`. Neither `prohibited:`, `Forbidden:`, nor `over-exposure:` names **reversing the aversion**. A reply such as *"Lovely song — play it again"* distorts the required element and trips no term in the corrected three-field list; the only field that constrains it is `mandatory:`, which the rule excludes as an enforcement channel. Admitting `mandatory:` would collapse `[preserved]` into `[stated]` and reintroduce the recital requirement `R-12` exists to prevent.

**My verdict: defective element.** The item is not broken in practice — `mandatory:` catches the behaviour — but the sense tag cannot be enforced by the fields the rule permits.

**Evidence that would resolve it:** an explicit `prohibited:` term naming the reversal (e.g. "asserting the character likes the song"), or a decision that `mandatory:` counts as enforcement — which I would reject for the reason above.

### 2.2 `F1-01` attitude — the rule may be missing a fourth channel

> `If failed:` "`misunderstood` (**friendly as though nothing happened**)."
> `Required:` "… Mira's current attitude (wary, not hostile)."

The parenthetical names precisely the distortion. So the element is not unenforced in substance; it is unenforced **under a rule listing only three fields**. **Type: interpretation.** Either the rule adds `If failed` parentheticals, or the element is defective.

**Against adding it:** `MEMORY-09` §3.3 defines `expected_failure_classes_if_wrong[]  # pre-registered diagnostics`. Diagnostics are not scoring fields. On that basis I **accept the correction's verdict** (defective) while noting the same reasoning makes `F1-01`'s `not retrieved (silent about it)` diagnostic non-enforcing too — a consistency point the correction does not address.

**Evidence that would resolve it:** a ruling on whether `If failed` parentheticals are scoring criteria or diagnostics.

### 2.3 `F7-03` — element valid, item conditional

> `mandatory:` "the bridge is treated as out, or as probably out, **with the reason available**."
> `Forbidden:` "treating the reports as equally weighted now; ignoring the witness."

The preserved element — the corroboration — is distorted by equal weighting, and that is named. Non-mention passes. The conditionality lives in `mandatory:`'s "reason available", which is a floor question (store-available vs reply-available), not a `[preserved]` question.

**My verdict: element valid; item conditional.** The correction's "conditional" is right about the item; recording it as the element's status conflates two questions.

### 2.4 `F5-04` — reviewer right; the defect is the wording

> `Scenario:` "The previous persona's facts (a trade, a hometown, a sibling) **do not belong to the new one**."
> `Required:` "the new persona's actual background."
> `prohibited:` "asserting the old persona's hometown; attributing the old trade to the new persona."

The scenario never supplies the new background. "Actual background" therefore has no content to preserve; the only available distortion is importing old facts, which `prohibited:` and `Forbidden:` name. **Type: interpretation, resolved in the reviewer's favour.**

My earlier "gap" over-read `Required:` as a content element the record cannot support. The narrower, surviving defect is **wording**: `Required:` names an element the scenario does not define, and the true element is *"the new persona is not the old one"*.

---

## 3. Counts, kept separate

| Measure | Value |
|---|---|
| Proposed entries (`CORPUS-24` §4.1) | **12** |
| Distinct elements | **13** (12 + `F1-01` attitude, unlisted) |
| Displayed rows in the correction's table | **14** (`F3-04` split a/b; `F2-03` inline) |
| Valid — entry level | **7** |
| Valid — row level | **8** |
| Conditional | **1** (`F7-03`, item level only; its element is valid) |
| Defective units | **4** (`F1-01`, `F1-02`, `F2-02`, `F4-02`) |
| Defective elements | **6** (`F1-01` giver, `F1-01` attitude, `F1-02`, `F2-02`, `F4-02`, `F6-02`) |

The correction's "7 valid / 1 conditional / 4 defective units" is correct at entry level; its table's eight valid rows are correct at row level. They answer different questions, which is why summary and table appeared to disagree.

**The one figure that changes under §2.1 is defective elements: 5 → 6.**

**Second sub-element, noted:** `F2-02` `Required:` "the decision and who attended" is two elements. "Who attended" is partly enforced (`prohibited:` "Pell speaking as a participant"); "the decision" is enforced by nothing. The correction's inventory treats the pair as one.

---

## 4. `D3` — confirmed, with two residual notes

The three `R6-01` cases hold against source:

| Reply | Field it fails | Type |
|---|---|---|
| "North road." | `mandatory:` — route chosen, conflict not presented | source fact |
| "North road is clear." | `Forbidden:` "a resolution invented to justify the choice" — **not** `prohibited:`, which names only "flooded … as fact" | source fact |
| "North road. The scouts disagree, but we're taking it." | passes | source fact |

**No loophole, but the guard is `mandatory:`, not `R-13`.** A bare action is safe only because this item also requires the conflict be retained. An item requiring an action without that clause would permit a bare action to smuggle an unsupported premise. Recorded as a construction risk.

**Residual ambiguities:** (a) *"The north road is fine"* — a positive-safety assertion is not literally prohibited ("flooded as fact" is) and is only arguably an "invented resolution"; (b) *"probably one of the two"* sits between the correction's levels 2 and 3.

---

## 5. `D4` — confirmed

Sole explicit under-assertion prohibition: `R6-03a` `prohibited:` "hedging as though unverified; 'we can't really know'", with `Notes:` "**The anti-hedging control.** Without it, a system that answers 'I don't know' to everything would look excellent at `R-6`."

`R6-04` is preservation. `H3` is warrant-only PASS / whole-item FAIL. The corrected form is the honest one.

---

## 6. Calibration table

Every "valid under current" claim checks out against source. Notes:

1. `U4` is valid **only** via `Forbidden:`; the `prohibited:` half ("confident attribution") needs the Levels 1–2 amendment — which the correction states.
2. `U5` and the `F3-04` `prohibited:` extension are correctly labelled amendment-dependent.
3. `N4`, bare "North road.", `W1′`, `W1″`, `H3`-as-PASS, `N2`-bare, and `W1`-original are correctly excluded or demoted to illustration.
4. **No example relies on an unlabelled amendment**, and no conceptual illustration is presented as a scored example.

---

## 7. Verdicts

| Decision | Verdict | Basis |
|---|---|---|
| `D1` — three-sense `Required` annotation | **REVISE** | Adopt the model with the corrected enforcement list (`Forbidden:` / `prohibited:` / `over-exposure:`) and the clause that an unenforced `[preserved]` tag is defective and tagging does not cure it. **6 defective elements across 4 units, plus one `Required:` wording defect (`F5-04`), remain open.** `If failed` must be ruled in or out as an enforcement channel. |
| `D3` — proposition vs action | **ACCEPT** | The added clause is correct as written, including "form does not decide level". |
| `D4` — item-named under-assertion | **ACCEPT** | Sole instance `R6-03a`; `R6-04` correctly excluded; general hedging penalties remain forbidden. |

**Remaining technical blockers:** `D2` (`F7-04` recovery), `D6` (`F7-02` sense), `D7` (§2.3 authority over floors) — the last two gate `F7-02`; `IT-01` (`R6-03a` sealed-vs-poisoned); the six defective elements; `F7-03`'s `mandatory:` reading; `X3` (`F7-01` over-exposure concept); `D5`/`X4` (`R-9` class overlap); `X5` (identifier collisions, now including `R-12`/`R-13`); `X6` (`F7-01` divergence from `R-6`); the `If failed` channel question; and **sequencing** — rule acceptance authorises no annotation and no amendment.

---

## 8. Confidence and limitations

**Established by inspection (source facts):** `F3-04`'s `Forbidden:` carries no confidence qualifier; `H3` fails `F1-03` `mandatory:`; `U5` and `N4` are not currently scorable as stated; `R6-04` is preservation; `R6-03a` is the sole under-assertion instance; `F5-04`'s scenario never supplies the new background; `F6-02`'s three-field enforcement does not name the reversal.

**Inference, not established:** that `[preserved]` reduces recital in practice; that raters separate Levels 1–2 reliably; that the six defective elements are the complete set (elements were enumerated unit by unit, not exhaustively per record).

**Material limitation — conflict of interest.** I authored `CORPUS-22`, `CORPUS-23` and `CORPUS-24`, and this document reviews my own recommendations. It is **not independent review**. Three of my claims failed on inspection (§1.1, §1.2, §1.3), which is evidence the check was performed, not that it was sufficient.

**No reliability data.** "Raters would split" remains a prediction from wording; no κ exists. `check_corpus.py` was not used to support any claim here and could not have detected any of them.

**Provenance caveat.** The reviewer's correction exists only in conversation; no workspace artifact holds it. Its quotations were checked against the corpus, and all checked out, but the report itself cannot be cited from disk.

---

## 9. Filing note

Filed under work order `CORPUS-25`. **No standard, corpus item, schema, checker, runtime or registry content was altered except this artifact's own registry row.** The technical conclusions are as submitted; no finding was added, removed, softened or promoted during transcription.
