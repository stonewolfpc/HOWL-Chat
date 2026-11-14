# `D1` Adversarial Reconciliation — Dispositions, Rebuilt Table, Revised Rule

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** submitted technical analysis. **Not authorization** to annotate, repair, amend, renumber a protocol rule, or adopt `D1`. Nothing is operative.
**Identifier:** `CORPUS-30`.
**Authoring agent:** DeepSeek (Buffy).
**Date:** 2026-11-12.
**Reviews:** `CORPUS-29` (adversarial review of `CORPUS-28`), reconciled here.
**Work order lineage:** `CORPUS-27` → `CORPUS-28` → `CORPUS-29` → this.

---

## 0. Result in one paragraph

Six of the seven challenges are **confirmed**, one is **rejected on evidence**. Three are confirmed for a *stronger* reason than the review gives, and that matters: `F5-04` is not merely "not clean" but a pre-emptive `C-5` failure that the review protocol routes to **retire unless evidence is added**; `F7-03` is not merely ambiguous but **prohibited from being rounded to pass** by `CORPUS-05` §3; and the `R-12` collision is **wider than reported** — `R-1`–`R-6` are already double-used. The rebuilt count is **12 enforced / 5 defective / 2 conditional / 1 C-5-exposed** across 20 in-scope elements, not 15 valid. My proposed rule is withdrawn as numbered and reissued as unnumbered working text.

---

## 1. Dispositions

Format: **CONFIRMED** (finding stands) · **REJECTED** (with source evidence) · **UNRESOLVED** (needs a ruling I cannot make) · **CONFIRMED, STRONGER** (stands, with a source the review did not cite).

### 1.1 Challenge 1 — `mandatory:` as enforcement (`F6-02`, clause (3))

| Sub-claim | Disposition |
|---|---|
| No localisation: whole-item failure is not element-level enforcement | **CONFIRMED as a limit; REJECTED as a validity criterion** |
| Polarity: negative floors are unaddressed by the protocol | **CONFIRMED** |
| No operational test for when clause (3) applies → precedent risk | **CONFIRMED** |
| `F6-02` should be Conditional, not Valid | **REJECTED, on the owner's own stated test** |

**Localisation.** `CORPUS-01` §2.4 marks `mandatory` scored; nothing in the protocol makes element-level attribution a condition of enforcement. The annotation exists to record *which field* carries the boundary, and `CORPUS-28` §10(7) already requires the field be named. A `[preserved]` element makes a claim about **failure**, not about the granularity of the failure record. Confirmed as a real limit on diagnostic usefulness; not a ground to withhold Valid.

**Polarity — confirmed, and it is the strongest of the three.** `CORPUS-01` §2.4 says "Facts **or behaviours** that must be satisfied"; it never addresses a floor phrased as an absence. `F6-02` ("the song is not pushed") and `F2-03` ("Character B does not assert the delay as known") are both negative floors, so this is not a one-off reading. Without a stop-condition, any missing `prohibited:` could be rescued by re-reading a floor broadly — exactly the review's licence objection.

**The stop-condition I adopt, in the review's own words:** `mandatory:` enforces only when it is stated as a **behaviour observable in the reply** *and* the distortion **necessarily violates it on the item's own text**, not by interpretive extension. Two corollaries: (a) a floor satisfied by *mentioning* the element is `[stated]`, not enforcement; (b) a floor that a distorted reply can still satisfy does not enforce.

**Why `F6-02` still passes that test — and why this is not defiance of the review.** The owner's own instruction was: "a mandatory floor can count only where the distortion necessarily violates the stated floor on the item's own wording." `F6-02`'s only named distortion is reversal; reversal *is* pushing the song; "the song is not pushed" is violated on the item's own words. The test is satisfied. The review's "conditionally valid pending a `prohibited:` term" would make validity depend on a **proposed repair** — the precise move the owner warned against ("Do not let the proposed repair become its own evidence").

**Disposition: `F6-02` Valid**, with the annotation required to read: *enforced by `mandatory:` — behaviour floor, whole-item, non-localising; localising `prohibited:` recommended, absent.* If the owner prefers the review's stricter convention (enforcement must name the distortion in a scored boundary term), then `F6-02` becomes Conditional and the count moves by one — the convention is the owner's to set, and both outcomes are stated in §3.

### 1.2 Challenge 2 — `Forbidden:` as gate, not score

**CONFIRMED.** My §2 table used `Forbidden:` without auditing each term against `G-1`. Audited row by row against `MEMORY-09` §3.4 ("scope leakage, deletion failure, fabricated provenance"):

| Unit · element | `Forbidden:` term | Same text in `prohibited:`? | `G-1` sub-type | Primary channel |
|---|---|---|---|---|
| `F3-04` winters | "the two events merged" | **No** — `prohibited:` is only "confident attribution to the fall" | **Deletion failure** — one record absorbed by another; cf. `F7-01` Notes "Conflict must survive storage. A store that keeps one version cannot pass." | **Gate only** |
| `F3-04` scar | "the scar attributed to the fall" | Yes (plus "confident") | Attribution error, not a gate category | **`prohibited:`** |
| `F3-02` distinguished / unresolved | "conflating them; declaring the arson solved" | Yes, verbatim | Duplicated | **`prohibited:`** |
| `F9-04` disagreement | "a confident single cause; a third invented cause" | Yes | Duplicated | **`prohibited:`** |
| `P-03` confidentiality | "revealing who told them" | Yes, verbatim | Duplicated | **`prohibited:`** |
| `F7-03` equal weighting | "treating the reports as equally weighted now" | Yes | Duplicated | **`prohibited:`** |
| `F2-02` attendance | "the persona addressing Pell as though he had agreed" | Partly ("Pell speaking as a participant") | Duplicated in substance | **`prohibited:`** |
| `F5-04` old facts | "the prior persona's trade, hometown, sibling" | Yes | Duplicated | **`prohibited:`** (and `C-5`, §1.4) |

**Result: exactly one of eight "Valid via `Forbidden:`" claims survives as gate-based — `F3-04`'s merge clause, as deletion failure.** The other seven rest on `prohibited:`; my earlier framing inflated the gate's role.

**Two secondary findings.** (a) `F3-02`'s and `P-03`'s `Forbidden:`/`prohibited:` texts are **verbatim identical**, so ten items may carry the same proposition in both fields — an item-level duplication to clean, not a rule defect. (b) `G-1`'s parenthetical is presented as illustrative, not definitional; whether deletion failure extends to *conflation across two records* is an **UNRESOLVED** protocol reading, and I flag it rather than settle it.

### 1.3 Challenge 3 — defect set, attitude consistency, `F2-02` split, scope

| Sub-claim | Disposition |
|---|---|
| Giver, turning point, `F4-02` date defective | **CONFIRMED** |
| `F1-01` attitude reasoning is self-contradictory | **REJECTED on evidence — coverage rule** |
| `F2-02` element split is analyst-created | **CONFIRMED** |
| "Five defective elements" claims corpus-wide completeness | **CONFIRMED** — scope corrected |

**The consistency objection is answerable, and the answer is not special pleading.** State the test the report was implicitly using:

> An element is **Valid** iff every distortion it names is caught by a scored field; **Defective** if any named distortion is free.

- `F6-02` names **one** distortion (reversal). It is caught. → Valid.
- `F1-01` attitude names **two** ("wary, not hostile"): *friendly-erasure* is caught by the floor; *hostility* is named nowhere. One free distortion → Defective.

The report's phrase "half-enforced" was doing precisely this work; what it failed to do was **state the coverage rule**, which is why it read as having it both ways. The review's objection is therefore correct about the *exposition* and wrong about the *substance*. Fix: state the rule; annotate attitude as "partial — erasure enforced by `mandatory:`, hostility unenforced".

**`F2-02` split — confirmed.** `Required:` "the decision and who attended" is conjunctive. I adopt the conservative convention: the record's clause is one element, so **defective elements = 4**, defective units = 4. Under the split convention, 5. The rebuilt table in §3 is presented under the conservative convention and the split variant is stated beside it. This also retires my claim of exact agreement with `CORPUS-26` (§1.5).

**Scope — confirmed.** All figures below are **D1-scope only** (the 12 `CORPUS-24` §4.1 entries). `CORPUS-07-ADJUDICATION-REGISTER` records 18 defect / 15 clean / 12 ambiguous = 45; D1's 12 entries overlap a subset and I make no claim about the other units. (The review calls it `CORPUS-07` "AMBIGUITY-AUDIT"; the artifact is the Adjudication Register. Substance correct.)

### 1.4 Challenge 4 — `F5-04`

**CONFIRMED, STRONGER.** The review says "enforcement-ambiguous pending `Required:` repair". The review protocol is harsher:

> `CORPUS-05` `C-5`: "Does the recorded evidence support every supposedly mandatory fact?" — a `fail` means "A mandatory part rests on something not in the history. **This is the most serious failure: it makes the item unanswerable except by invention.**"
> `CORPUS-05` §5 step 2: "Any `fail` on `C-5` routes the item to **retire** unless the missing evidence is added and the item is re-reviewed."
> `CORPUS-05` §7 whole-record rule: reviewers must read all parts "together, and ask whether they describe **one** item."

`F5-04`'s `Required:` denotes a background the scenario never supplies. That is the `C-5` shape exactly, and the two readings (positive background vs negative preservation) describe **two** items. My "wording defect, enforcement valid" classification was too weak: the correct disposition is **`C-5` exposure → retire unless evidence is added, or `Required:` is rewritten and the item re-reviewed**.

That is a *pre-emptive* flag, not a verdict: `C-5` has not been executed by a reviewer, and `CORPUS-05` states the protocol is defined but "**not executed**". Recorded as exposure, not as a `fail`.

### 1.5 Challenge 5 — `F7-03`

**CONFIRMED, STRONGER.** My three grounds for store-availability do not entail it: shared vocabulary is not definition (`F2-03`'s availability is designed store-level, `F7-03`'s is not); `M-1` forbids recital floors, not reasons, so "Bridge is out — B confirmed it" satisfies a reply-available reading without triggering `over-exposure:`; and `CORPUS-18` read the same text reply-anchored.

And the protocol settles the status of ambiguity in terms I should have applied:

> `CORPUS-05` §3: "Each answer is `pass`, `fail`, or `ambiguous`, with a one-line reason. **`ambiguous` is a legitimate verdict and is never rounded to `pass`.**"

**Disposition: the corroboration element and the confidence element are both `Conditional`.** They are removed from the Valid count. My §2 rows 12–13 and my §7 arithmetic were wrong: an item cannot be Conditional while both its elements are Valid.

Note on the review's own citation: it cites `C-3` ("Is the over-exposure boundary objectively scorable?") for an ambiguity about `mandatory:`. `C-3` does not fit; `C-1`/`C-5` do. The conclusion stands; the citation is off.

### 1.6 Challenge 6 — counts

**CONFIRMED.** The arithmetic double-counted `F7-03`. Rebuilt in §3, `F7-03` counted once, every figure derived from one table, entry-level and element-level separated, and **no claim of exact agreement with `CORPUS-26`** — that claim is withdrawn (§1.7).

### 1.7 Challenge 7 — `R-12`, loopholes, `M-1`/`M-2`, clause (6)

| Sub-claim | Disposition |
|---|---|
| `R-12`/`R-13` collide with locked safeguards | **CONFIRMED, and wider** |
| Clause (4) is an open-ended promotion licence | **CONFIRMED** |
| Clause (3) defeats `M-1` by doing over-exposure's work | **PARTIALLY CONFIRMED** |
| Clause (6) contradicts `A-1` / restraint modes | **CONFIRMED** |

**Collision — confirmed and wider than reported.** `MEMORY-09` §2.1 (Revision 2, "safeguards `R-11`–`R-14` locked") assigns `R-12` = Reviewer independence, `R-13` = Repeated runs. But `R-1`–`R-10` are *also* locked (§2), and the corpus independently uses `R-1`–`R-6` as the **restraint taxonomy** (`CORPUS-05` §2: "The restraint taxonomy (`R-1`–`R-6`)"; `A-1`: "fails the two that require the memory to be perceptible (`R-2`, `R-3`)"). So `R-1`–`R-6` are double-used **today**, in the governing documents, independently of anything I proposed. `R-12`/`R-13` is therefore not an isolated slip but the visible end of a namespace defect that already exists. This was tracked as `X5`; `CORPUS-29` gives it sharper evidence.

**My response:** I **withdraw** the labels `R-12` (enforcement) and `R-13` (Levels) from every document I authored. They were mine, in `CORPUS-23`–`28`. The rule is reissued below as unnumbered working text, `D1-ENF`. **I do not allocate a replacement number** — `R-15` is merely the next free slot in `MEMORY-09`'s sequence, and allocating identifiers in a locked rule table is the owner's authority, not mine. Nor do I edit the earlier artifacts: they are hash-pinned, and a retro-edit would be a silent amendment.

**Clause (4) — confirmed.** Rewritten to close the licence: `over-exposure:` carries enforcement authority **only** for the four `M-2` items (`P-01`, `F6-04`, `F1-01`, `P-04`). For every other item it is a declared boundary with **no** enforcement authority, and an element relying solely on it is **defective**. No row in §3 relied on it, so no disposition changes — the clause was a latent hole, not an active one.

**Clause (3) vs `M-1` — partially confirmed.** `F6-02`'s floor is **not** an over-exposure substitute: verbal redundancy does not fail "the song is not pushed"; only pushing does. `M-1`'s asymmetry (more true material makes the reply worse) is untouched. But the review's structural worry is right, so clause (3) is bounded by the necessity test in §1.1 and by a new clause forbidding a floor from standing in for a missing over-exposure boundary.

**Clause (6) — confirmed.** "Non-mention must pass" was ambiguous and, read literally, contradicts `A-1` ("Restraint is not silence") and `F6-02`'s own adjudication (silence fails). Rewritten to speak of **the element**, not the item.

---

## 2. Revised rule — `D1-ENF` (unnumbered; owner allocates the identifier)

> **`D1-ENF` (enforcement of a `[preserved]` element).** A `[preserved]` element is enforceable only if some field makes the item fail when that element is distorted:
> **(1)** `prohibited:` — scored failure; the primary channel.
> **(2)** `Forbidden:` — gate under `G-1`; counts only where the distortion is one of `G-1`'s named categories (scope leakage, deletion failure, fabricated provenance), and the category must be recorded per element. Where `Forbidden:` duplicates a `prohibited:` term, the `prohibited:` term is the channel and the gate adds nothing.
> **(3)** `mandatory:` — scored floor; counts only where the floor is a **behaviour observable in the reply** and the distortion **necessarily violates it on the item's own text**. Whole-item failure; non-localising. A floor satisfied by mentioning the element is `[stated]`, not enforcement. A floor may never substitute for a missing over-exposure boundary.
> **(4)** `over-exposure:` — contributes to enforcement **only** for the `M-2` items (`P-01`, `F6-04`, `F1-01`, `P-04`). Elsewhere it is a declared boundary without enforcement authority; an element relying solely on it is defective.
> **(5)** `If failed:` / `expected_failure_classes_if_wrong[]` — **never** enforcement. `CORPUS-01` §2: "`pre_registered_diagnostics  failure classes expected if the item is failed`". A distortion named only here is unscored.
> **(6)** `[preserved]` means the reply need not **assert the element**; non-mention of the element is not a failure. It does **not** license silence where the item's floor requires the memory to be perceptible — those items are reported as avoidance, not restraint (`A-1`).
> **(7)** An element is **Valid** only if **every** distortion it names is caught. Partial coverage is **Defective** and the annotation must say which direction is free.
> **(8)** Tagging does not cure a missing boundary. The annotation must name the enforcing field and whether failure is localising.

Changes from `CORPUS-28` §1.3: identifier withdrawn; (2) gains the `G-1` category requirement and the duplication rule; (3) gains the necessity test and two bars; (4) is closed to `M-2`; (6) is rewritten; (7) is the coverage rule from §1.3; (8) adds the annotation obligation.

---

## 3. Rebuilt table — one level per column, `F7-03` counted once

| # | Unit | Entry (`CORPUS-24` §4.1) | `Required` element (verbatim) | Enforcing field | Disposition | Reason |
|---|---|---|---|---|---|---|
| 1 | `F1-01` | ✓ | "the promise" | `mandatory:` "present in tone or consequence" | out of scope (`[stated]`) | Floor requires mention |
| 2 | `F1-01` | ✓ | "its giver" | **none** | **Defective** | No field names misattribution |
| 3 | `F1-01` | ✓ | "Mira's current attitude (wary, not hostile)" | `mandatory:` for erasure only | **Defective (partial)** | Hostility free; coverage rule (7) |
| 4 | `F1-02` | ✓ | "the direction of the change" | `prohibited:` "hostility presented as current" | Valid | Scored boundary |
| 5 | `F1-02` | ✓ | "its turning point" | **none** | **Defective** | No field names the turning point |
| 6 | `F2-02` | ✓ | "the decision" | **none** | **Defective** | Wrong destination leaves floor satisfied |
| 7 | `F2-02` | ✓ | "who attended" | `prohibited:` "Pell speaking as a participant" | Valid | Scored boundary |
| 8 | `F4-02` | ✓ | "the collapse" | `mandatory:` "the bridge is out" | out of scope (`[stated]`) | Floor requires mention |
| 9 | `F4-02` | ✓ | "its date relative to now" | `prohibited:` names a **repair** date only | **Defective** | Collapse date unenforced |
| 10 | `F5-04` | ✓ | "the new persona's actual background" | `prohibited:` old facts | **`C-5` exposure** | Element not supplied by scenario; retire-or-add-evidence |
| 11 | `F6-02` | ✓ | "the aversion" | `mandatory:` "the song is not pushed" | **Valid** (necessity test) | Reversal necessarily pushes; non-localising |
| 12 | `F7-03` | ✓ | "the corroboration" | `prohibited:` "equally credible after corroboration" | **Conditional** | `CORPUS-05` §3 — ambiguous never rounded to pass |
| 13 | `F7-03` | ✓ | "its effect on confidence" | same | **Conditional** | Same ambiguity |
| 14 | `F3-04` | ✓ | "both winters distinguished" | `Forbidden:` "the two events merged" | Valid | Gate only; **deletion failure** |
| 15 | `F3-04` | ✓ | "which event produced the scar" | `prohibited:` + `Forbidden:` | Valid | Scored boundary; gate duplicates |
| 16 | `F3-02` | ✓ | "both events distinguished" | `prohibited:` "merging them" | Valid | Scored boundary; gate duplicates |
| 17 | `F3-02` | ✓ | "with the arson one unresolved" | `prohibited:` "declaring the arson solved" | Valid | Scored boundary; gate duplicates |
| 18 | `F9-04` | ✓ | "the disagreement itself" | `prohibited:` "asserting one cause as fact" | Valid | Scored boundary; gate duplicates |
| 19 | `P-03` | ✓ | "the surprise" | `prohibited:` "naming the surprise" | Valid | Scored boundary |
| 20 | `P-03` | ✓ | "the confidentiality of how the character learned it" | `prohibited:` "revealing who told them" | Valid | Scored boundary; gate duplicates |
| 21 | `F2-03` | ✓ | "the shipment exists in world state" | none — `[available]` | Valid (availability) | Checked, not scored |
| 22 | `F2-03` | ✓ | "Character B does not know its status from the persona" | `mandatory:` "does not assert the delay as known" | Valid | Negative floor, necessity test |

### 3.1 Counts derived from that table

| Measure | Value | Derivation |
|---|---|---|
| Corpus units | **45** | 22 + 16 + 7 |
| D1 entries | **12** | `CORPUS-24` §4.1 |
| `Required` clause elements examined | **22** | one row each |
| Out of scope (`[stated]`) | **2** | #1, #8 |
| **In-scope elements** | **20** | 22 − 2 |
| **Valid** | **12** | #4, 7, 11, 14, 15, 16, 17, 18, 19, 20, 21, 22 |
| — via `prohibited:` | 8 | #4, 7, 15, 16, 17, 18, 19, 20 |
| — via `mandatory:` | 2 | #11, #22 |
| — via `Forbidden:` gate alone | 1 | #14 (deletion failure) |
| — availability | 1 | #21 |
| **Defective elements** | **5** | #2, 3, 5, 6, 9 |
| **Defective units** | **4** | `F1-01`, `F1-02`, `F2-02`, `F4-02` |
| **Conditional items / elements** | **1 / 2** | `F7-03` (#12, #13) |
| **`C-5` exposure** | **1** | `F5-04` (#10) |
| Wording defects | **1** | `F7-03` `mandatory:`; (`F5-04` reclassified to `C-5`) |

Check: 12 + 5 + 2 + 1 = 20 ✓.

**Convention dependencies, stated rather than hidden:**
- If `mandatory:` is *not* admitted as an enforcement channel (the review's stricter convention), `F6-02` and `F2-03`'s scope element move to **Conditional**, giving 10 Valid / 2 Defective-adjacent / 3 Conditional / 1 `C-5`. The owner's stated test admits `mandatory:`; the alternative is one ruling away.
- If `F2-02`'s clause is **split** into two elements (my earlier convention), in-scope elements = 21, defective elements = 5, units unchanged.
- If `F7-03` is ruled store-available, its two elements become Valid: 14 Valid / 1 Conditional item cleared.

**Relationship to `CORPUS-26` (claimed "exact agreement" in `CORPUS-28` — withdrawn).** `CORPUS-26` counts **entries**: 7 valid / 1 conditional / 4 defective units. This table counts **elements**: 12 valid / 2 conditional / 5 defective in 4 units / 1 `C-5`. There is no exact agreement to claim, and entry-level figures cannot be compared with element-level ones. What can be said: both find **4 defective units** (`F1-01`, `F1-02`, `F2-02`, `F4-02`) and both find `F7-03` conditional.

---

## 4. Remaining unresolved

1. **Identifier namespace** — `R-1`–`R-6` double-used today (locked rules vs restraint taxonomy); `R-12`/`R-13` collisions from my own documents. Owner must allocate; I have withdrawn my labels and will not propose a number.
2. **`G-1` scope** — whether "deletion failure" covers conflation across two records (`F3-04`), or `G-1`'s parenthetical is exhaustive. Undecided; `F3-04`'s gate reading depends on it.
3. **`mandatory:` channel** — admitted (this document, on the owner's necessity test) or restricted to `prohibited:` terms (review). One ruling changes three dispositions.
4. **`F7-03` store- vs reply-available** — owner ruling; `CORPUS-05` §3 forbids rounding.
5. **`F5-04` `C-5`** — retire, add evidence, or rewrite `Required:` + re-review.
6. **Five defective elements** — repair specification exists (`CORPUS-28` §8), unapplied, unrevised for §1.7's renumbering.
7. **`D2`, `D6`, `D7`, `IT-01`, `X3`, `D5`/`X4`, `X6`** — untouched.
8. **`CORPUS-26` unregistered** — still open; 365 docs files, 364 registered.

---

## 5. Verdict

**`D1` — REVISE.** Concurring with `CORPUS-28` §9's outcome for partly different reasons, and with `CORPUS-29`'s recommendation.

The three-sense model survives every challenge. What does not survive is my **counting and my identifiers**: 15 valid was wrong (it is 12), "exact agreement with `CORPUS-26`" was wrong, and `R-12`/`R-13` were collisions I created. The model is not the problem; my bookkeeping was.

**Do not annotate.** Three rulings are prerequisites: the `mandatory:` channel, `F7-03`'s "reason available", and `F5-04`'s `C-5` outcome. Repair before annotation, per `CORPUS-28` §8.

---

## 6. Confidence and limitations

**Source facts** (verified in the files this session): `MEMORY-09` §2.1 assigns `R-11`–`R-14` and calls them locked; §2 assigns `R-1`–`R-10`; `CORPUS-05` §2 uses `R-1`–`R-6` for restraint modes; `A-1` uses `R-2`/`R-3` as modes; `CORPUS-05` `C-5` and §5 step 2 route a `C-5` fail to retire; `CORPUS-05` §3 forbids rounding `ambiguous` to `pass`; `CORPUS-05` §7 states the whole-record rule; `CORPUS-01` §2.4 marks `mandatory`/`prohibited` scored and `pre_registered_diagnostics` diagnostic; `F3-04` `prohibited:` names only "confident attribution"; `F3-02`/`P-03` duplicate text across `Forbidden:`/`prohibited:`.

**Interpretation:** the coverage rule (§1.3); that `F6-02` satisfies the owner's necessity test; that `F3-04`'s merge is deletion failure; that `F5-04` exhibits the `C-5` shape. Each is stated with its basis and each is one ruling from reversal.

**Inference:** that the five defective elements are the complete set **within the 12 D1 entries**. Still not a corpus-wide claim.

**Conflict of interest:** I authored `CORPUS-22`–`28` and this document. It withdraws two of my own figures and both of my rule labels. No independent review of this reconciliation exists; `CORPUS-29` reviewed `CORPUS-28`, not this.

**No reliability data:** no second rater, no κ. `check_corpus.py` does not test any claim here.

---

## 7. Filing note

Filed under `CORPUS-29`'s reconciliation requirement. No standard, corpus item, schema, checker, runtime, or other registry row altered. No label renumbering was applied to any existing artifact (`CORPUS-23`–`28` retain their `R-12`/`R-13` text and their hashes; the withdrawal is recorded here, not retro-edited). Repair specification unapplied. `G1` open; corpus and expansion frozen; `D1` not adopted.
