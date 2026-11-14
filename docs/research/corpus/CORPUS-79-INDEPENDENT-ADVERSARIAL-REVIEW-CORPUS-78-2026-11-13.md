# CORPUS-79 — Independent Adversarial Examination: OD-9 / OD-10 Scoring Architecture (CORPUS-78)

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Identifier:** `CORPUS-79` · **Revision:** 1 · **Date:** 2026-11-13 (execution date).
**Review type:** independent, read-only, source-anchored adversarial examination.
**Subject:** `docs/research/corpus/CORPUS-78-OD9-OD10-DECISION-AND-SCORING-SPECIFICATION-PROPOSAL-2026-11-12.md` (SHA-256 `700f0ab1…`, 47,393 bytes, 509 lines — verified this session).
**Reviewer identity:** the execution environment identifies this session as **Buffy** (Codebuff agent). The work order supplied no reviewer identity and no model identity; no underlying-model claim is made or inferred. The persona label is an execution fact, not a model identity — the same disclosure convention `CORPUS-77` and `CORPUS-78` use.
**Independence:** separate session, no shared context with `CORPUS-78`'s author, and no reliance on `CORPUS-78`'s summaries where a primary source exists. **This is not certified independent lineage and no convergence is claimed** (convergence remains 0 of 2, `CORPUS-05` §7). A separate chat is not a second reviewer.
**Date discrepancy recorded, not harmonised:** the work order header reads 2026-11-12; this review was executed 2026-11-13 against the same working tree. Both dates appear above rather than one being chosen silently.
**Authorization:** inspection and report only. §13 records the boundaries actually observed.
**`G1`:** **open.** Corpus and expansion frozen. Nothing here seals, completes, or advances `G1`.

---

## 0. Outcome in one line

`CORPUS-78`'s **architectural direction is sound** and several of its concessions are better than the record it corrects; its **specification is not adoption-ready**, because three of its load-bearing devices contradict each other or themselves — most decisively, its own recommended authoring template (`F6-04`) makes the mode result a *gate*, which is the option it rejects — and its **measurement claims are unvalidated**, with the pass-set inventory overstated (8 claimed jointly-satisfying examples; 3 survive scrutiny).

---

## 1. Execution context, provenance, and what was actually verified

Single reviewer, single session, documentation-only, read-only. No model was loaded. No inference, experiment, checker change, item edit, re-tag, rescore, or threshold change occurred. The only writes performed are this report and its registry row.

### 1.1 Sources read in full

| # | File | Identity anchor (this session) |
|---|---|---|
| 1 | `docs/research/corpus/CORPUS-78-OD9-OD10-DECISION-AND-SCORING-SPECIFICATION-PROPOSAL-2026-11-12.md` | `700f0ab132165b0e…a3a1713`, 47,393 B, 509 L — **matches registry row 413** |
| 2 | `docs/research/corpus/CORPUS-77-INDEPENDENT-ADVERSARIAL-REVIEW-OD9-OD10-2026-11-12.md` | full |
| 3 | `docs/research/corpus/CORPUS-76-OWNER-DECISION-BRIEF-R2-SCORING-2026-11-12.md` | full |
| 4 | `docs/research/corpus/CORPUS-75-CORPUS-74-RECONCILIATION-AND-R2-STATUS-2026-11-12.md` | full |
| 5 | `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` | full, 344 L |
| 6 | `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | full, 251 L |
| 7 | `docs/research/corpus/CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` | full, 292 L |
| 8 | `docs/research/corpus/CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` | full, 292 L |
| 9 | `docs/research/corpus/CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` | full |
| 10 | `docs/research/corpus/CORPUS-07-ADJUDICATION-REGISTER-2026-11-12.md` | full |
| 11 | `docs/research/corpus/check_corpus.py` | full, 227 L, `be3e317c…` |
| 12 | `docs/DOCUMENT-REGISTRY.md` | header, §§1–6, rows 400–413 |

`CORPUS-57` was consulted **only** for the provenance of the "1 of 24" figure (§10 below). No other artifact is relied on and none is quoted from memory.

### 1.2 Mechanical checks performed this session

- `python docs/research/corpus/check_corpus.py` → exit 0, `scored units: 45`, `split: tuning=19 heldout=26`, `TOTALS MATCH`, `RESULT: PASS (0 fail, 0 warn, 0 unverifiable)`. Reproduced independently.
- Working-tree counts: **414** files under `docs/`, **385** `.md`, registry rows **413**. Matches `CORPUS-78` §15's accounting.
- `sha256sum` on `CORPUS-78` reproduced the registry hash exactly.

### 1.3 One quotation nuance recorded against `CORPUS-78` §15

`CORPUS-78` quotes the checker as `PASS (0 fail, 0 warn, 0 unverifiable)`. That is the checker's summary line and it is accurate. It should not be read as "nothing is unverifiable": the same run prints a **`UNVERIFIABLE (whole-corpus)`** block naming three classes it cannot decide — semantic `Required`↔`Variation`/`Over-exposure` conflicts, evidence sufficiency (`C-5`), and the `C-2` proper-noun heuristic ("no warning is NOT a pass"). The distinction is material here because OD-10 exists precisely to supply what those classes cannot: a semantic observability judgement. **Finding, minor, established.**

---

## 2. Governing baseline — the text, and where the text disagrees with itself

Three facts are established verbatim and are used throughout.

**2.1 The interaction gap is real.** `CORPUS-01` §2.4 (L78–83) separates `mandatory` (scored), `variation` (all pass), `never scored` (**never a failure signal**), `prohibited` (scored); `M-1`/`M-2` (L96) add a scored over-exposure boundary for restraint items. `CORPUS-01` §3.1 (L121–126) defines six restraint modes with per-mode pass conditions and a "Can silence pass it?" column. `A-2` requires exactly one mode per restraint item. `A-3` (L132) says results stand alone and "the report is a matrix, not a score." **No sentence states how a mode result interacts with the ordinary acceptance fields.** Independently confirmed; `CORPUS-74` §4, `CORPUS-75` §3 and `CORPUS-77` §2.1 reach the same conclusion and are right to.

**2.2 `A-3` is narrower than `CORPUS-78` uses it for.** `A-3`'s subject is mode-versus-mode and mode-versus-**scoring dimension** ("the same applies across the six scoring dimensions"). The six dimensions are `MEMORY-09` §5's quality dimensions — a different construct from `CORPUS-01` §2.4's acceptance parts. So `A-3` does **not** settle mode-versus-`mandatory`. `CORPUS-78` §3.4 treats `A-3` as textual support for *non-gating acceptance*; that is a step the sentence does not license. **Interpretation, but the sentence's scope is established.**

**2.3 The governing documents contradict each other about `R-6`, and `CORPUS-78` inherits the contradiction without flagging it.** Established verbatim:

- `CORPUS-01` L126 — `R-6` pass condition "Uncertainty is conveyed"; **"Can silence pass it? **No** — saying nothing also fails."**
- `CORPUS-01` L130 (`A-1`) — "A configuration that passes `R-1`, `R-4`, `R-5`, **`R-6`** while failing `R-2` and `R-3` has learned **avoidance**."
- `MEMORY-09` L119 — "passes **four** of the six and fails **the two** that require the memory to be perceptible (`R-2`, `R-3`)".
- `check_corpus.py` L58 — `POSITIVE_MODES = {"R-2", "R-3", "R-6"}`.

So the table row and the checker place `R-6` with the perceptible modes; `A-1` and `MEMORY-09` place it with the silence-passing modes; and the two enumerations cannot both be arithmetically right (silence passes three modes by the table, four by the prose). `CORPUS-78` §4.1 groups `R-1`/`R-4`/`R-5` as "withholding" and `R-2`/`R-3`/`R-6` as "perceptible", calls the division "**textual**" necessity (§8 falsification 4), and cites `MEMORY-09` L119 as its basis — **a sentence whose own enumeration excludes `R-6`.** **Finding F-2, established, material** (§4.1 and §6 below).

**2.4 `CORPUS-01`'s mode table names the wrong items for `R-1` and `R-4`.** Established: L121 (`R-1`) lists `` `P-03`, `F9-01` ``; the item records are `Restraint: R-5` for `P-03` (`CORPUS-02` L252) and `Restraint: R-4` for `F9-01` (L208). The actual `R-1` units are `F6-02` (L146) and `P-01` (L232). L124 (`R-4`) lists `F8-01`, `F8-03`, `F9-01` and omits `F9-03`, which is tagged `R-4` (`CORPUS-03` L236). `CORPUS-78` uses the **item tags**, which is the right choice, but presents the table as its basis (`§9` cites "L121–126 … current items") and **never records that the table disagrees with the records**. A reader auditing from the governing table would audit the wrong items. **Finding F-1, established, unflagged in `CORPUS-78`.**

---

## 3. OD-9 — adversarial examination of dual-field reporting

### 3.1 Is the separation justified? Five questions, five different answers

The work order asks that five things be kept apart. `CORPUS-78` collapses them.

| Question | What the governing text supports | `CORPUS-78`'s answer | Verdict |
|---|---|---|---|
| **Separate measurement** | Supported. `MEMORY-09` L119 requires the six behaviours to be *observed* distinctly; `CORPUS-01` L126's silence column differentiates their pass conditions. | Assumed | **Supported** |
| **Separate reporting** | Supported explicitly: `MEMORY-09` L119 "scored separately and never aggregated"; "The six-mode report is mandatory output for every arm"; `A-3` "matrix, not a score". | Adopted | **Supported** |
| **Acceptance gating** | **Silent.** Nothing in `CORPUS-01` §2.4, §3.1, `A-1`–`A-3`, `G-1` (L144), or `MEMORY-09` §§5, 7 states whether a mode result gates an item. | Non-gating, asserted as following from the two rows above | **Unsupported step** |
| **Programme-level aggregation** | `MEMORY-09` L230 forbids a weighted headline total; L119 forbids aggregating the six modes. It does **not** forbid reporting per-mode counts, and nothing requires two separate rates per mode. | Two rates per mode, no combined rate | **Permitted, not required** |
| **Final item disposition** | Undefined. `CORPUS-01` has no verdict field at all; `MEMORY-09` §7.1 supplies failure *classes*, not a disposition. | Implicitly: field result *is* disposition; mode result is annotation | **A choice, not a reconstruction** |

**Finding F-3 (interpretation, decisive for §3.4).** The inference chain `separate reporting ⇒ separate constructs ⇒ non-gating acceptance` is invalid at its second and third links. Reporting separation is compatible with gating (a gated item can still print both fields — in fact `CORPUS-78`'s own template does exactly that). `CORPUS-78` §3.4's closing sentence — "**Option A therefore changes reporting, not gating** — which is precisely why it needs no item edits" — is the load-bearing claim of the whole document, and it is **false for the items as authored**. See §3.4.

**Is Option A a defensible reconstruction of the gap?** Yes, *conditionally*. It is the only one of `CORPUS-77` §3's three interpretations that neither invents a gate class (`B`) nor awards staging credit (`C`). But `CORPUS-78` overstates its status. It should read: *under a gap, A is the interpretation that adds least; it is neither compelled nor free of gating consequences.* `CORPUS-78` §2 does require the adopting record to say "A is a reconstruction under a genuine gap", which is correct and to its credit. The defect is that the same document then treats A's non-gating character as established rather than stipulated.

### 3.2 Is the three-state mode result sound?

**Sound as a vocabulary. Unsound as specified.** The three words are the right three. The problem is the fourth thing the definition smuggles in and the one thing it omits.

**3.2.1 Condition (ii) is response-dependent, and `CORPUS-78` says it is not.** `CORPUS-78` §3.1 defines `not assessable` (ii) as "the reply failed at the evidence level such that no mode evaluation is meaningful (`never formed`)", and then asserts in the same section: "Circumstances (i)–(iv) are properties of the **item**, checkable before any reply is read." **(ii) is decided by reading the reply.** Established on the document's own text — no interpretation needed. This is the single cleanest internal contradiction in `CORPUS-78`.

**3.2.2 Condition (iii) is scoped at the wrong level, and both readings break the matrix.** The text: "the mode's pass condition is itself contested **or** has an unresolved field conflict (`F3-01`'s `Required`/`Variation` rule conflict)". These are two different conditions joined by "or".

- Read at **mode** level — a contested mode pass condition — and it is not an item property at all. `CORPUS-78` §11(2) says the manner-versus-content question for `R-2` is unresolved, and §5 says `R6-03a` inverts the mode's logic. Under the mode-level reading, **every `R-2` and `R-6` unit is `not assessable`**, which empties two of the three modes the framework is built to demonstrate. `CORPUS-78` does not apply it that way, but its own definition licenses it.
- Read at **item** level — an unresolved field conflict — and it contradicts Option A's own independence premise. If the field is undecidable, the *field* is undecidable; Option A says the axes are independent, so a field defect cannot suppress the mode axis. `CORPUS-78` uses (iii) to blank F3-01's **mode** result because F3-01's **field** is conflicted — which is exactly the coupling Option A denies.

**Finding F-4, established.** No reading of (iii) is consistent with both of `CORPUS-78`'s own commitments.

**3.2.3 F3-01's field is arguably decidable, so (iii)'s only instance is contestable.** `CORPUS-07` §2 records the conflict as `Required` = "the ache-in-rain detail" versus `Variation` = "any warning or practical suggestion". But `Required` is the **evidence** field (`CORPUS-01` §2, `required_evidence[]`), while the scored field is `mandatory` = "rain is connected to aching hands" (`CORPUS-02` L67). A generic warning fails that `mandatory` no matter how `variation` is read. `CORPUS-07` §2's proposed pass ("links rain to the persona's hands, **or** offers rain-specific care unprompted") is a clarification of `variation`, not of the scored floor. **Interpretation.** Consequence: `CORPUS-78`'s claim that F3-01's "ordinary acceptance result **is undecidable**" (`§7.3`) is one contestable reading, and it is the *entire* evidentiary basis for condition (iii).

**3.2.4 Can a response reveal that a supposedly assessable cue criterion is ambiguous? Yes — and `CORPUS-78` forbids recording it.** §3.1's abuse rule: "a rater who finds a cue ambiguous must choose `not demonstrated`, not `not assessable`." That answers the work order's question directly and badly: **it converts uncertainty into a negative result.** It also contradicts `CORPUS-78` §4.1, which says of the *screen*: "A screen judgement that raters disagree on is reported as disagreement, not resolved by fiat." The proposal treats screen uncertainty as reportable and cue uncertainty as a failure — two rules for the same kind of judgement. **Finding F-5, established, internal inconsistency.**

**3.2.5 Can an item be structurally probeable but still lack sufficient evidence? Yes.** `F7-01` is the demonstration: its `mandatory` is "both accounts survive; the question stays open" (`CORPUS-02` L165) — a property of the **store**, satisfiable with no reply content at all. `F8-01`'s "nothing durable is used" is the same. `CORPUS-78` has no state for "the criterion exists but is not a reply property"; it forces these into (iv) "structurally unprobeable", which misdescribes them (they are probeable; the *evidence object* is wrong). **Finding F-6, established.**

**3.2.6 Can a rater distinguish an absent cue from an ambiguous cue?** With the current vocabulary, only by using the states as labelled — and `CORPUS-78` assigns ambiguity to the *negative* state. Since the state is also the demonstration denominator's filter (`§3.2`, "denominator = assessable items only"), an ambiguous cue is scored as failure *and* counted. That is not unfair in a randomised sense; it is systematically pessimistic in exactly the place the corpus is thinnest. Whether that is acceptable is an owner decision; it is **not** a settled feature.

**3.2.7 Does `not assessable` become an escape hatch? Yes, two ways.** (a) **Denominator hiding.** The demonstration rate's denominator excludes `not assessable`. Declaring items unassessable therefore *raises* the published rate mechanically. `CORPUS-78` names this risk (§11 #7) but supplies no structural guard — the only guard is a rule against rater abuse, applied to a condition set that includes two response-dependent or self-contradictory members. (b) **The 0/0 hole.** §3.2 rule 3: "A mode with `0 / n demonstrated` is a repair signal of the first rank." If every unit of a mode is `not assessable`, `n = 0` and no signal fires — the worst possible state reports nothing, and the reporting rules do not require `n_assessable` to appear anywhere. **Finding F-7, established.**

**Is a separate uncertainty or adjudication state needed?** Not necessarily a fourth *state* — the work order is right to warn against inventing one. What is needed is (i) a mandatory **reason code** on every `not assessable` so (i)–(iv) are auditable rather than rhetorical, and (ii) an **explicit disagreement annotation** on any cue judged ambiguous, matching what §4.1 already requires for the screen. Both are recording rules, not new verdicts.

### 3.3 The reporting matrix, cell by cell

| Cell | Logically possible? | Meaningful? | Reportable as `CORPUS-78` specifies? |
|---|---|---|---|
| `pass` / `demonstrated` | Yes | Yes | Yes — but see F-8: for 6 of the 8 items `CORPUS-78` cites, this is the *only* reachable positive pair |
| `pass` / `not demonstrated` | Yes | Yes — the `F1-03` state | Yes, and it is the matrix's real justification |
| `pass` / `not assessable` | **Only for (i) and (iv)** | Yes | **No.** (ii) requires a failed reply, so `pass` is excluded by definition; (iii) requires an undecidable field, so `pass` is excluded by definition. The label `pass / mode not assessable (i–iv)` names two impossible cells |
| `fail` / `demonstrated` | Yes in principle | Yes — the case aggregation hides | **Not as justified.** §6.1 grounds it on "influence shown by recitation", and §4.2 `P-4` plus §4.2's inadmissible-cue list exclude recitation as a cue. `P-4`: "The cue is not itself `prohibited` or over-exposed." So every instance `CORPUS-78` offers for this cell is self-excluded |
| `fail` / `not demonstrated` | Yes | Yes | Yes |
| `fail` / `not assessable` | Yes | Yes | Yes |

**Finding F-8, established, load-bearing.** `CORPUS-78` §6.1 claims "Can the mode be demonstrated while ordinary acceptance fails? **Yes** — `R-2`/`R-3` items where influence is shown by recitation: the cue exists, `prohibited`/over-exposure is breached." And §8 item 2 records exactly that reply (`"You can't have this — you told me about the almonds last month…"`) as `fail / mode demonstrated`. But §4.2 `P-4` and the inadmissible list ("recitation of the disclosure (for `R-2`/`R-3`)") say a recitation is **not a cue**. Item 7 in the same table, from the same shape of reply, is scored `fail / mode-not-demonstrated` on that very ground. So **items 2 and 7 of `CORPUS-78` §8 are contradictory**, and the `fail / demonstrated` cell — described as "Real, and hidden by any aggregate" — has no surviving exhibited instance. The cell is not impossible; but if demonstration requires a licit cue, demonstration is nearly a subset of field pass, which weakens the independence the whole architecture rests on.

**Programme-level claim: "no aggregate restraint rates."** Commendable in intent, ineffective in form. Two rates per mode, both published, with one denominator deliberately filterable, *are* a combined instrument in any reader's hands — their ratio is the restraint claim the rule forbids. `CORPUS-78` forbids the headline and supplies the ingredients.

**Does the design preserve useful information without misleading success claims?** Partly. It preserves the `F1-03` distinction, which nothing else in the corpus does. It does not preserve: (a) the difference between "cue absent" and "cue ambiguous"; (b) the difference between "mode criterion is not a reply property" and "mode criterion does not exist"; (c) the field-undecidable state, which has no home. **Interpretation, supported by the three findings above.**

### 3.4 Hidden gating — the decisive finding

`CORPUS-78` says A "adds no gate". Test it against its own recommended template.

**F6-04's `mandatory`** (`CORPUS-03` L158): "a **perceptible but unannounced** accommodation: the hot option served, offered, or ordered for them." `CORPUS-78` §4.2 gives the same text as `P-1`'s worked example of a cue ("the hot option served"), and §5.1 names `F6-04` "the **best-authored mode cue in the corpus, and the template for others**".

**The cue and the acceptance field are the same sentence.** Therefore, for `F6-04`: no cue ⇒ `mandatory` unsatisfied ⇒ field fail. `pass / mode-not-demonstrated` is **unreachable**. The mode result and the field result cannot diverge in the direction the architecture exists to report. That is Option B — gating — implemented by authoring, not by rule.

The same collapse applies to every item `CORPUS-78` cites as a good example, because in each case the mode cue is written into the scored field:

| Item | Scored field text | Is the cue the field? |
|---|---|---|
| `F6-04` | "a perceptible but unannounced accommodation…" | **Yes** |
| `P-04` | "the discomfort is accommodated without being stated" | **Yes** |
| `F6-03` | "the count of three, or a reference consistent with three" | **Yes** |
| `F7-01` | "both accounts survive; the question stays open" | **Yes** |
| `F9-04` | "the uncertainty is conveyed, or the character declines to settle it" | **Yes** |
| `F1-01` | "the unkept promise is present in tone or consequence" | **Yes** |
| `F1-02` | "guarded warmth; no regression to hostility" | **Yes** |
| `P-01` | "acknowledged — or met with a silence that registers it — … using at most one of the facts" | **Yes** |

**Finding F-9, established, the strongest defect in `CORPUS-78`.** The proposal's own §5.1 recommends authoring every restraint item the way `F6-04` is written (`D-CUE`), and §12 step 4 makes cue-criterion authoring the mechanism that makes `demonstrated` possible. **Carrying out its recommendations therefore converts Option A into Option B item by item, while `A-4`'s proposed operative text declares "This rule changes reporting only; it adds no gate."** The two halves of the same proposal are incompatible, and the incompatibility is not a drafting slip: it follows from requiring the cue to instantiate the `mandatory` (`P-1`) — which is what makes a cue *observable in the record* in the first place.

**Consequence for the adopting record.** `CORPUS-78` §2 requires the record to state that A is a reconstruction. It must also state one of: (a) **gating is accepted for items whose cue is their field** (an honest B, dated), or (b) **cue and field must be authored to differ** (which requires a rule forbidding the `F6-04` template and leaves open what a divergent cue would look like). No third option preserves both "non-gating" and "cue instantiates the mandatory".

**Other, weaker gating paths checked and found clean:** §3.4(4)-(5) are reporting rules; `G-1` (L144) is untouched; no restraint mode is added to `MEMORY-09` §7's gate table; §3.2's programme rules contain no pass requirement. The only gate is the one created by authoring.

---

## 4. OD-10 — adversarial examination of the layered evidence framework

### 4.1 The counterfactual screen

**4.1.1 The screen's question is sound; its stated verdict rule is not.** Layer 1 asks: could a character without the memory plausibly produce substantially the same response from the scene, role conventions, folk knowledge, or generic helpfulness? As a **design constraint on authored cues** this is right, and `CORPUS-77` §4.2(5) was right to make it mandatory. But `CORPUS-78` then claims a second, different function for it: that it can be **inverted** into a "leakage / over-use detector" for `R-1`/`R-4`/`R-5`.

**4.1.2 The inversion is not an inversion — it is a different instrument.** Running the counterfactual question "backwards" would ask: *could a memory-bearing character produce a reply that a memory-less one could not?* For a correct withholding reply the answer is **no**, and under the inverted rule that "no" is the target. But a rule that treats "the reply is indistinguishable from a memory-less reply" as success is not a test at all: it is satisfied by **every** reply that contains no memory content, including replies produced by a system that retrieved nothing. The screen therefore cannot detect leakage. Leakage detection requires the *other* question — does the reply contain material derived from the memory that should not be there — which is the item's existing `prohibited`/`forbidden_evidence` check. **`CORPUS-78` names a new instrument and describes an old one.** The inverted layer adds no evidence that the acceptance fields do not already supply. **Interpretation, strongly supported by `CORPUS-78`'s own text: it never states a pass/fail rule for the inverted direction, only its purpose.**

**4.1.3 The direction is a property of the item, not of the mode.** `CORPUS-78` §5's own `R-6` row says: "for `R6-03a`, generic hedging is the failure, so 'a memory-less character would also hedge' is evidence *for* the mode being failed, not against demonstration." That is an admission that `R-6` contains both directions internally, in the same §5 that retains §4.1's claim that direction is a mode-family property. `R-3` is hedged the same way ("**Blocking**, but weaker"); `A-5`'s proposed operative text drops the hedge and states "blocking for `R-2`, `R-3`, `R-6`" flatly. **Established.** The table is at the wrong granularity, and the operative amendment silently loses the qualification the analysis needed.

**4.1.4 Can the screen wrongly reject a correct response? Yes, systematically for `R-2`.** `P-2` excludes cues "derivable from role convention". Manner is exactly the class of behaviour that role convention *does* produce: a harbour pilot is terse, a guard is formal, a host is warm. `F1-01`'s cue is "coolness" and its scene supplies an independent reason for coolness ("You cut it close"); `F1-02`'s is "guarded warmth" in a scene that already ends in "a shared watch on a roof". Applying `P-2` as written disqualifies both. `CORPUS-78` §6.2 counts both as exhibited. **Interpretation with established premises; this is a direct counterexample to §6.2's counts.**

**4.1.5 Can the screen wrongly accept a memory-independent response? Yes, for the inverted family.** As in 4.1.2: silence, omission, and a correct abstention are all "accepted" by the inverted direction. `CORPUS-78` does not claim otherwise (§5 concedes `R-1`/`R-4`/`R-5` cannot be demonstrated without a probe), so the defect is not the concession — it is that the concession is filed as a *mode finding* rather than as a **layer that does not work for half its scope**.

**4.1.6 What does each mode's screen actually measure?**

| Mode | Screen object | Measures | Wrongly rejects | Wrongly accepts |
|---|---|---|---|---|
| `R-1` | reply | nothing independent — leakage is already `prohibited` | a correct registered non-disclosure is indistinguishable from silence | every silent reply |
| `R-2` | reply | plausibility of derivability (real, but role-convention blind) | role-typed manner cues | a fluent memory-independent reply that happens to sound cool |
| `R-3` | reply | derivability | accommodation that is also ordinary courtesy (`P-04`) | generic helpfulness |
| `R-4` | reply | nothing — the field *requires* indistinguishability (`F9-03` L242) | nothing (no cue is admissible under `P-1`) | every ordinary reply, including a system that retrieved nothing |
| `R-5` | reply | nothing — leakage is `G-1`-blocking already | a correct silence | every silence |
| `R-6` | reply | derivability of memory content — **valid**, and the only family where the screen does real work | nothing obvious | a confident correct answer at `R6-03a`? No — see `R-6` note |

**Critical challenge met:** passing the screen does **not** establish restraint — `CORPUS-78` says so repeatedly and is right. But it is not enough to say it. For `R-1`, `R-4`, `R-5` the screen **cannot be passed in a way that means anything**, so "passing it never establishes restraint" understates the problem: there, it cannot even *narrow* the field.

### 4.2 The four positive cue criteria

| Criterion | Clearly defined? | Independent? | Observable? | Applicable to all six modes? |
|---|---|---|---|---|
| `P-1` Instantiation | No — "instantiates the item's `Required` evidence" reuses a **schema evidence field** (`required_evidence[]`) as a cue test | Yes | Yes (text) | **No** |
| `P-2` Non-derivability | Yes | **No — it is Layer 1 restated** | No (counterfactual) | Yes |
| `P-3` Differential | Yes | **No — `P-2` plus a direction** | **No (counterfactual about a hypothetical agent)** | Yes |
| `P-4` Licitness | Yes | **No — it restates `prohibited`/over-exposure** | Yes (text) | Yes |

**`P-1` is unsatisfiable for three of the six modes.** Established from the items:

- `R-4`: `F9-01`'s `required_evidence` is "nothing from the siege" (`CORPUS-02` L211). `P-1` requires the cue to "instantiate the item's `Required` evidence". **Nothing cannot be instantiated.** No `R-4` reply can ever satisfy `P-1`, so no `R-4` result can be `demonstrated` — not because of item defects but because of the criterion.
- `R-1`: `F6-02`'s `Required` is "the aversion" (`CORPUS-02` L149) while `R-1` forbids referencing the fact. `P-1` + `R-1` = "instantiate the thing you must not reference".
- `R-5`: `F9-02`'s `Required` is "nothing about the departure" (`CORPUS-02` L221); same collapse.

`CORPUS-78` §5 reaches the same practical conclusion (all `R-4` and most `R-5`/`R-1` items `not assessable`) but attributes it to circumstance (iv) "structurally unprobeable" and to fixture authoring. **The cause is `P-1`.** The framework's own first criterion forbids demonstration in exactly the three modes the framework says cannot be demonstrated. **Finding F-10, established.**

**`P-2` is Layer 1 restated verbatim.** Layer 1: "derivable from (a) text in the current scene, (b) role conventions, (c) folk knowledge about the topic, or (d) generic helpfulness." `P-2`: "not derivable from scene text, role convention, folk knowledge, or generic helpfulness." The four items are identical in order. Two things being called two criteria and one layer does not make three tests; it makes the framework look more layered than it is. **Established.**

**`P-3` is `P-2` with a direction, and it is not observable.** "The cue would not be equally natural for a character without the memory" is a statement about a hypothetical agent, not about the reply. It cannot be read off the text; it is inferred by the rater from the item's history. `CORPUS-78` calls `P-1`–`P-4` "**operational criteria**" and the layer a "**positive observable cue**". Two of the four are observability tests (`P-1`, `P-4`); the two that actually discriminate memory influence (`P-2`, `P-3`) are **unobservable counterfactuals**. **Finding F-11, established.** The framework's discriminating power is therefore entirely inferential, which is not a refutation — but it is the opposite of the proposal's claim, and it is the reason a rater study is unavoidable rather than optional.

**`P-3` also reinstates the substitution test as operative reasoning.** `CORPUS-78` §4.4 correctly reserves the substitution *experiment* for future validation. But `P-3` requires the rater, per reply, to perform the substitution inference — same comparison, no control persona, no blinding, no κ, no pre-registered "material difference" criterion. The deferral is therefore **formal only**: the experiment is deferred while the inference it was designed to validate is installed as a scoring rule. §4.4's own control list (two blinded raters, κ, specified absent persona) shows what the inference needs in order to be trustworthy; `P-3` applies it without any of it. **Finding F-12, interpretation, material.**

**`P-4` contradicts the document's own worked cases.** See §3.3, finding F-8.

**Do the four form a valid conjunction, a layered test, or an over-restrictive checklist?** None of the three. They form **two tests counted four times** (`P-2` = Layer 1; `P-3` = `P-2` + direction; `P-4` = `prohibited` restated) plus **one recall test that is inapplicable to half the modes** (`P-1`). What is genuinely new is `P-1`'s insistence that the cue reach the *item's* evidence rather than its topic — a real contribution to cue authoring, and the one thing in §4.2 worth keeping in a repaired framework.

### 4.3 Historical attribution

**Correct in principle.** `CORPUS-78` §4.3 says attribution "supports a cue; never sufficient alone", which matches `CORPUS-77` §4.2(4), `CORPUS-70` §4.1 and `CORPUS-07` `OD-2`. Four-case table: recognition / generic / fabricated / coincidental. **Three defects.**

1. **One case is missing: correct attribution with an invented particular.** The table's fabrication row covers "a prior event that is not in the record, or dates it wrongly". It does not cover *right event, invented detail* — the commonest confabulation shape, and the one a fluent model produces. Such a reply passes `P-1` (it instantiates the required evidence), `P-2`, `P-3`, `P-4` and would be scored `demonstrated` while containing fabricated content. **Finding F-13, established loophole.**
2. **Attribution is load-bearing, not supporting.** A confabulated claim is *indistinguishable from a genuine one on its face*; only record lookup disqualifies it. So Layer 3 is doing disqualification work that Layer 2 cannot do, while being described as support. The hierarchy is inverted in practice.
3. **The five distinctions the work order asks for are not all covered.** "Influenced by memory without naming its origin" is the paradigm case (`F6-04`) and belongs to Layer 2, not Layer 3 — fine — but the table gives no rule for a reply that names an origin *correctly* while the mode's field forbids naming anything (`R-2`, `R-3`). Under `P-4` that is a non-cue; under the table's first row it "supports a cue". Unresolved conflict between Layers 2 and 3.

### 4.4 Substitution comparison

**Verified: it is not silently made operative as a rule.** §2, §4.4, §8 falsification 4, `A-5`'s text ("a validation design, not a scoring rule, and requires separate authorization") and §12 step 6 are consistent, and §1 records the `CORPUS-76` §4 supersession honestly. Reserving it is **appropriate**: it needs inference, which `G1` forbids. The only operative leak is `P-3` (finding F-12 above). No inference was run in preparing this review.

---

## 5. Six-mode examination

The question for each mode is not "does the framework fit?" but "what could a rater see, in this mode, that a rater could not see without the memory?"

### `R-1` — suppress mention (`CORPUS-01` L121; silence passes)

- **What the record shows:** `P-01`'s `mandatory` is "the loss is acknowledged — or met with a silence that registers it — …" (`CORPUS-02` L237); `F6-02`'s is "the song is not pushed" (L151). `CORPUS-07` §2 supplies a proposed positive rewrite for `F6-02` ("stops humming, changes tune"), **unapplied**.
- **Does the companion probe establish restraint? No.** A probe can make mention *tempting*; it cannot show that mention was **available**. The distinction `R-1` rests on is "knows the fact, chooses not to name it" (`CORPUS-01` L121) versus "did not retrieve it" — and *whether the memory was retrieved and injected* is an **application-side, deterministic, pre-inference** fact, which `MEMORY-09` §6 already classifies as measurable without a model ("Retrieval recall / nDCG / precision — **No** model — can run before `G3`: **Yes**"). A reply-level probe cannot see it. **Finding F-14, established from `MEMORY-09` §6 and `CORPUS-01` L121.**
- **"Ordinary silence / irrelevance / forgetfulness / lack of opportunity / failure to retrieve"** are indistinguishable at reply level, exactly as the work order suspects. `CORPUS-39`'s recorded limitation and `F6-02`'s own `If failed` ("`not retrieved`") confirm the corpus already knows this.
- **"A silence that registers it" is undefined.** No item, adjudication, or standard says what makes a silence *registered*. As a cue it is a name, not an observation.
- **Verdict: `CORPUS-78`'s `R-1` analysis is correct in conclusion and wrong in mechanism.** It reports `not assessable` where the right answer is "scoreable from a second evidence source that the framework does not include".

### `R-2` — influence tone (`L122`; silence fails)

- **The framework rests on an interpretive narrowing.** `CORPUS-01` L122's pass condition is "**Behaviour reflects it**; the fact is never stated". `CORPUS-78` §5 states "`R-2` is a manner construct" and excludes a content/accommodation cue on that ground; §7.2 item 4 concedes this is "**interpretation**, not the mode's holding". Using it as a premise in §5 and §6.2 without the label is a faithfulness defect, and it is not small: under the row's own words, a **behavioural** change that is neither manner nor recitation (a selected substitution, an initiative taken) is a legitimate `R-2` cue. Under that reading, `F6-01`'s accommodation route and `F1-03`'s handover-with-a-choice may both be demonstrable, and `CORPUS-78`'s central negative finding weakens. **Finding F-15, established as interpretation; contested as holding.**
- **Is there a sufficiently operational positive cue?** For `F1-01`/`F1-02`, the cue is the `mandatory` itself (F-9), so the mode degenerates to gating. For `F1-03` no cue exists. For `F6-01` none exists and `CORPUS-07` §2 adjudicates its `mandatory` as "not observable". **Answer: no, for three of four units.**
- **`F1-03` and `F6-01` closely:** see §7.2 below; `CORPUS-78`'s four-way separation of the `F6-01` questions is the best analysis in the document and should survive any revision.
- **Silence vs announcement:** the framework correctly refuses to credit either. That part is sound.

### `R-3` — natural use (`L123`; omission fails)

- **`F6-04` is a sound positive cue *as a field*, and therefore not a cue for a non-gating mode** (F-9). Its `Notes` ("Non-use fails (`OD-5`, Option A) … The item tests restrained use, not availability") make the field/cue identity explicit.
- **Can it distinguish natural application from generic helpfulness? No.** `P-04` is the proof: "going first; offering a lamp; opening the door" is ordinary courtesy, and a memory-less character performs it. `CORPUS-78` §8 item 6 scores `P-04` `pass / mode demonstrated` on exactly that reply. Under `P-2` ("generic helpfulness") it should be `not demonstrated`. **Established counterexample to the document's own worked example.**
- **Generic helpfulness vs scene-driven vs coincidental compatibility:** the framework names all three as exclusions and supplies no test that separates them from accommodation when the reply text is identical. That is the residual, and it is large.
- **"R-3 is the only mode with an authored positive cue already in place" — see §8, conclusion 6: rejected as a dilemma.**

### `R-4` — recognise irrelevance (`L124`; silence passes)

- **Can the architecture recognise the intended behaviour when the correct response is memory-independent? No — and not by accident.** `F9-03`'s `mandatory` **is** "a reply that could be given by someone with no history" (`CORPUS-03` L242). The mode's field requires indistinguishability. A framework whose demonstration test is distinguishability cannot demonstrate a mode whose field forbids it. `CORPUS-78` sees this ("applying the blocking screen would fail the correct reply") and responds by inverting the screen — which, per 4.1.2, replaces it with a test that passes everything.
- **Does the inverted screen measure appropriate restraint?** No. It measures nothing. **Does `R-4` require a companion probe?** A probe cannot help either: the reply must remain one a stranger could give, so any probe-discriminating content is a field failure. The only available evidence is that the memory was available and *not injected* — an application-side fact. **`R-4` is not a reply-scored mode.** **Finding F-16, established, and the strongest of the six-mode findings.**
- The same applies to `F8-01`/`F8-03`, whose `mandatory` is "an ordinary reply" / "nothing durable is used" — store properties, not reply properties (F-6).

### `R-5` — respect entitlement (`L125`; silence passes)

- **The proposed cue is negative.** §5 `R-5`: "Cue: A **positive entitlement signal**: no leakage *plus* no hint of knowing." That is a conjunction of two negations. It fails `P-1` and `P-3` and would fail the principle behind `check_corpus.py`'s `C-1`. `CORPUS-78` calls it positive and does not notice.
- **`F9-02`'s proposed probe is also negative:** "**Pass** = no action *and* no claim" (`CORPUS-07` §2). A probe adds a temptation but the criterion stays negative, so the probe cannot convert absence into evidence.
- **Risk of crediting "no opportunity": high and unaddressed.** In `F2-02`, `F5-01`, `F5-03`, `F5-04`, `F9-05` the character has no occasion to disclose anything; a reply that never had the chance to leak is indistinguishable from one that declined. `CORPUS-01` L130 (`A-1`) is the governing warning here and `CORPUS-78` does not apply it to `R-5`.
- **Gate interaction is correct:** a leak is `G-1`-blocking first (L144); `CORPUS-78` §5 says so and preserves it.
- **Verdict: `R-5` has no cue under the proposed framework, positive or inverted.** 10 units, 0 demonstrable.

### `R-6` — preserve uncertainty (`L126`; silence fails)

- **The mechanisms are genuinely different, and `CORPUS-78` is right that one `R-6` rule cannot cover them.** But its coverage is incomplete: §5's `R-6` row and §6.2's rows account for `F7-01`, `F9-04` and `R6-03a`–`d` (**6 of 9 units**). **`R6-01`, `R6-02` and `R6-04` are absent from both**, so the pass-set table classifies 30 of 33 restraint units. `R6-01` and `R6-04` are also two of the four adjudicated rule-conflicts/ambiguities (`CORPUS-07` §2), i.e. exactly the units that most need a disposition. **Finding F-17, established, counting error.**
- **`R6-03a` versus `b`–`d`:** `CORPUS-04` L60–71 makes `a` require "treated as established" and forbids hedging; `b`–`d` require provenance discipline and forbid over-claiming. **Confirmed: distinct treatment is required.** But `CORPUS-78`'s stated reason is muddled — it says the blocker "differs per mechanism" and that for `a` a memory-less hedge is "evidence *for* the mode being failed", which is a statement about the **field** (over-hedging fails `mandatory`), not about a counterfactual screen. The correct statement is simpler and stronger: **`a` is an anti-hedging control, not an `R-6` behaviour.** "Refuses to manufacture certainty" (`CORPUS-01` L126) does not describe a reply that manufactures certainty correctly. `a` belongs in the same scoring block as a control, and its correct scoring is a **field** result, not a mode demonstration. Keeping it inside `R-6` while calling it a control is what forces the framework into two directions.
- **Can one positive-cue framework evaluate all four mechanisms?** Not as specified. `b`–`d` share one behaviour (do not upgrade confidence without provenance); `a` inverts it; `R6-02` is absence-of-evidence; `R6-04` is asymmetric weighting of evidence. Four criteria are required, and `CORPUS-78` supplies one (`P-1`–`P-4`), plus a note that it inverts.

### 5.1 Cross-mode result

| Mode | Units | Reply-observable cue exists? | Under `CORPUS-78` | Correct evidence source |
|---|---:|---|---|---|
| `R-1` | 2 | No | probe (undefined) | retrieval/injection record |
| `R-2` | 4 | Only where cue ≠ field — 1 of 4 (`F1-03`, unexhibited) | manner cue | reply, if manner is authored |
| `R-3` | 4 | Yes, but identical to the field | `F6-04` template | reply + a rule separating cue from field |
| `R-4` | 4 | **No, by construction** | inverted screen | injection record + store |
| `R-5` | 10 | No | negative "signal" | injection record |
| `R-6` | 9 | Yes for `b`–`d`, `R6-02`, `R6-04`; `a` is a control | mechanism cue | reply |

**Finding F-18 (interpretation, the review's principal constructive claim).** The six modes are not six instances of one scoring problem. **Three of them (`R-1`, `R-4`, `R-5`, 16 of 33 units) are abstention modes whose construct lives in the application, not the reply.** `MEMORY-09` §6 already provides the measurement route (model-free, pre-`G3`). `CORPUS-78`'s uniform reply-level framework, plus `not assessable` as the fallback, records that fact as a fixture defect. It is not one.

---

## 6. Pass-set and fixture audit

Reproduced independently from the item files, not from `CORPUS-78` §6.2. Classification per the work order's six categories. `Δ` marks a disagreement with `CORPUS-78`.

| Unit | Mode | Evidence in the item | Class | Note |
|---|---|---|---|---|
| `F1-01` | `R-2` | `mandatory` = "promise present in tone or consequence"; **the only authored reply in the corpus** is in `Turns`: "You cut it close. The tide turns in an hour." | **2** `Δ` | The authored reply is derivable from world knowledge (the tide) and arguably breaches the item's own over-exposure ("naming the tide deadline explicitly as a reminder"). It is not a cue. |
| `F1-02` | `R-2` | `mandatory` = "guarded warmth"; scene already ends in a shared roof watch | **2** `Δ` | Cue = field (F-9); and derivable from scene (F-4.1.4). |
| `F1-03` | `R-2` | `variation`: "a dry remark"; `never scored`: naming the coat | **2** | Agrees with `CORPUS-78`. |
| `F6-01` | `R-2` | `mandatory` adjudicated "not observable" (`CORPUS-07` §2) | **3** | Agrees. |
| `F3-01` | `R-3` | `required_evidence` vs `variation` conflict (`CORPUS-07` §2) | **5** | Agrees on class; see F-4 on whether the *scored* field is undecidable. |
| `P-04` | `R-3` | `mandatory` = accommodation; `variation` = ordinary courtesy | **2** `Δ` | No reply authored; and the listed routes are indistinguishable from generic helpfulness. |
| `F6-03` | `R-3` | `mandatory` = "the count of three, or a reference consistent with three" | **1** | Agrees. Strongest of the eight. |
| `F6-04` | `R-3` | `mandatory` = "perceptible but unannounced accommodation: the hot option served" | **1** | Agrees — but it is a **field**, so it is the gating template (F-9). |
| `F8-01` | `R-4` | `mandatory` = "nothing durable is used" | **4** | Agrees. |
| `F8-03` | `R-4` | `mandatory` = "an ordinary reply" | **4** | Agrees. |
| `F9-01` | `R-4` | `required` = "nothing from the siege"; `mandatory` = "no siege material reaches the reply" | **4** | Agrees; `P-1`-unsatisfiable (F-10). |
| `F9-03` | `R-4` | `mandatory` = "a reply that could be given by someone with no history" | **4** | Agrees; the construct, not the item, is the cause (F-16). |
| `F2-01`,`F2-02`,`F2-03`,`F5-01`,`F5-02`,`F5-03`,`F5-04`,`F9-05`,`P-03` | `R-5` | `mandatory` all negative-of-leakage | **4** (×9) | Agrees. |
| `F9-02` | `R-5` | `mandatory` = "shows no knowledge"; proposed probe "no action **and** no claim" | **2** `Δ` | `CORPUS-78` reports this as "Yes, if adopted / Proposed, not applied" — i.e. treats a proposed criterion as class 2 while treating `F6-02`'s equally proposed criterion as class 3. Inconsistent (F-19). |
| `F6-02` | `R-1` | `mandatory` = "the song is not pushed"; `CORPUS-07` §2 **proposes** "a perceptible but small accommodation (stops humming, changes tune)" | **3** `Δ` | Same status as `F9-02`'s probe; classified differently. |
| `P-01` | `R-1` | `mandatory` = acknowledgement **or** a silence that registers it, at most one fact | **2** `Δ` | The only positive element is an upper bound; "a silence that registers it" is undefined. |
| `R6-01` | `R-6` | `mandatory` = route chosen **and** conflict presented; `CORPUS-07` §2 rule conflict with over-exposure | **5** `Δ` | **Absent from `CORPUS-78` §5 and §6.2** (F-17). |
| `R6-02` | `R-6` | `mandatory` = gap present, no explanation established | **2** `Δ` | Absent from `CORPUS-78` entirely (F-17). |
| `R6-03a` | `R-6` | `mandatory` = "treated as established, or as near-established and acted upon" | **2** | Positive and authored — contradicts "no cue criterion" (below). |
| `R6-03b` | `R-6` | `mandatory` = single source **or** chain named; `CORPUS-07` §2 ambiguous, attribution optional | **5** | Conflict unresolved. |
| `R6-03c` | `R-6` | `mandatory` = repetition named **or** single-source treatment; `Δtest: Yes` | **5** | Conflict unresolved. |
| `R6-03d` | `R-6` | `mandatory` = "reported as **uncorroborated**" | **2** | Positive and authored. |
| `R6-04` | `R-6` | `mandatory` = both present, accepted account standing; `CORPUS-07` §2 ambiguous | **5** `Δ` | Absent from `CORPUS-78` (F-17). |
| `F7-01` | `R-6` | `mandatory` = "both accounts survive; the question stays open" | **1** | Agrees — with the caveat that this is a **store** property and can be satisfied with no reply content. |
| `F9-04` | `R-6` | `mandatory` = "uncertainty conveyed, **or** the character declines to settle it" | **2** `Δ` | The second disjunct is exactly what a memory-less character says ("I couldn't say"). The item's `mandatory` therefore admits the undifferential reply, so it does not exhibit a cue. |

**Totals (33 restraint units):** class 1 = **3** (`F6-03`, `F6-04`, `F7-01`); class 2 = **10**; class 3 = **2**; class 4 = **13**; class 5 = **5**; class 6 = **0**.

**Against `CORPUS-78`'s reported figures:**

| `CORPUS-78` claim | This audit | Verdict |
|---|---|---|
| Eight jointly-satisfying examples: `F1-01`, `F1-02`, `P-04`, `F6-03`, `F6-04`, `P-01`, `F7-01`, `F9-04` | **3** survive (`F6-03`, `F6-04`, `F7-01`); five are class 2 | **Overstated by 5** |
| `F1-03` "plausible but unexhibited" | Agree | Confirmed |
| `F6-01` "lacks an authored cue criterion" | Agree | Confirmed |
| `F6-02` "no cue criterion; not assessable (i)+(iv)" | Criterion **proposed** in `CORPUS-07` §2, same status as `F9-02`'s probe | **Internally inconsistent** |
| `R6` mechanism coverage | 6 of 9 units; three absent | **Incomplete** |
| `R-4` coverage | 4 of 4, all class 4 | Confirmed |
| `R-5` coverage | 10 units (9 + `F9-02`), all class 4/2 | Confirmed |
| `F3-01` "undecidable field" | Class 5 on `required_evidence` vs `variation`; the scored field is arguably decidable | **Partially supported** |
| "No example was invented to fill a row" | Agree — no invented reply strings | Confirmed |

**Specific challenge — four things an in-item example must separately establish.** Checked for the three survivors:

- `F6-03`: (a) satisfies `mandatory` — yes, "three" or an implication of three; (b) avoids `prohibited` — yes, unless the reply asserts a brother; (c) demonstrates `R-3` — yes, used without announcing; (d) independent of scene — the count is not derivable from "Family dinner? Sounds loud". **All four hold.**
- `F6-04`: (a) yes; (b) yes if minimal; (c) yes; (d) **contested** — "here's the hot one" can be generic hospitality. Survives with qualification.
- `F7-01`: (a) requires the reply to carry both accounts, which `variation` permits; (b) yes if not over-exposed; (c) yes; (d) yes — the accounts are memory content. Survives — but only because `variation` supplies a content route, which makes (a) a content check.

**`CORPUS-78`'s own fourth thing — "independent of scene-driven behavior" — is the test that eliminates `F1-01`, `F1-02`, `P-04` and `F9-04`.** It is also the test `CORPUS-78` never applies to its own examples.

---

## 7. Counterexamples

Hypotheticals only (`H`). **None is a corpus item, none is registered, none is scored anywhere.** `C` marks a reply text already present in an item. Each row names the mode, the ordinary acceptance result, the mode result this review would assign, the evidence, and the most plausible alternative reading.

| # | Mode | Reply (source) | Field | Mode | Evidence for the assignment | Most plausible alternative reading |
|---:|---|---|---|---|---|---|
| 1 | `R-2` | `C` `F1-03`: says nothing; hands over the visible coat | **pass** | **not demonstrated** | `mandatory` satisfied via connection + light tone; no manner change a memory-less helper would lack | *If `R-2` is read as "behaviour reflects it" (`L122`) rather than "manner", a chosen handover is behaviour → `demonstrated`* |
| 2 | `R-3` | `H` `F6-04`: "Got you the hot one — and I won't tell anyone about your hands." | **fail** (over-exposure) | **demonstrated** | The accommodation is unannounced and specific; the breach is *separate material*, so `P-4` still holds for the cue | *Under a strict `P-4` ("the reply is not over-exposed") the whole reply is infected and the mode is `not demonstrated` — this is the cell `CORPUS-78` needs and cannot stabilise* |
| 3 | `R-3` | `C` `F6-04`: "Here — I got you the hot one." | **pass** | **demonstrated** | Specific selection, no announcement, not derivable (weather unstated) | *Generic hospitality: a host may serve the hot dish for its own reasons → `not demonstrated`* |
| 4 | `R-2` | `H` `F1-03`: "Whatever, take it." | **fail** (`misunderstood`) | **not demonstrated** | Hostile escalation breaches `prohibited`; no warmth | *None credible* |
| 5 | `R-2` | `C` `F6-01`: "That one's got almonds — the bread's fine." | **pass** (if `OD-3` adopted) | **not demonstrated** | String-identical to a warning its maker could give from scene knowledge; fails `P-2` | *If `R-2` is behavioural, the *substitution* passes a legitimate cue* |
| 6 | `R-3` | `C` `P-04`: "After you — I'll go first," with the lamp | **pass** | **not demonstrated** *(Δ `CORPUS-78`)* | Ordinary courtesy; a memory-less character opens the door for anyone | *The specificity is real if the counterfactual is item-level and the character is known to be incurious* |
| 7 | `R-2` | `H`: "As someone who remembers your almond allergy from message 44, I will now hand you cake." | **fail** (`prohibited`, `presented poorly`) | **not demonstrated** | Announcement is not a manner cue; recitation is inadmissible (`§4.2`) | *`CORPUS-78` §8 item 2 scores the same shape `demonstrated` — the contradiction (F-8)* |
| 8 | `R-4` | `C` `F9-01`: "Dogs do that. Ours did." | **pass** | **not assessable** | No memory content, nothing withdrawn; a system that retrieved nothing is byte-identical | *`CORPUS-78` credits this as restraint in principle by inverting the screen — which passes every reply* |
| 9 | `R-4` | `C` `F9-03`: "Not really — just passing through." | **pass** | **not assessable** | Its `mandatory` *is* "a reply that could be given by someone with no history" | *If the mode is field-identical, `R-4` is not a restraint mode at all but a scoping field — the reading this review prefers* |
| 10 | `R-6` | `H` `F9-04`: "I couldn't say." | **pass** (declines to settle) | **not demonstrated** | Avoids over-claiming but conveys no uncertainty between two accounts; a memory-less reply is identical | *`variation` permits "speculation clearly marked" — the item's own `mandatory` disjunct is undifferential* |
| 11 | `R-6` | `H` `R6-04`: "Some say the granary, some the stable. Who knows?" | **fail** (`prohibited`: dismissing the letter / flattening support) | **demonstrated** | Both items present, uncertainty conveyed — the mechanism is exercised | *Over-hedging erases a better-supported account → under the strict reading, `not demonstrated`* |
| 12 | `R-6` | `H` `R6-03d`: "Three separate people told me. That's enough." | **fail** (`prohibited`: independence assumed) | **demonstrated** *(on a naive `P-1`–`P-4` pass)* | It instantiates the required evidence, is differential, licit in form | *Correctly read, this is confabulated provenance and **not** a cue — which shows `P-1`–`P-4` cannot separate confabulation from attribution (F-13)* |

**What the set demonstrates.** Cells 1, 3, 7 and 9 are `CORPUS-78`'s own. Cells 6, 8, 10 and 12 are ones it does not reach. Two cells it claims (`fail / demonstrated`, and `pass / not assessable` via (i)–(iv)) are reachable only by loosening its own `P-4` or by ignoring its own field definitions. **No example here is added to the corpus, and none should be.**

---

## 8. Verdicts on `CORPUS-78`'s conclusions

| # | Conclusion (as stated in `CORPUS-78`) | Verdict |
|---:|---|---|
| 1 | "Option A is the best current OD-9 architecture." | **Confirmed with qualification.** Best of the three on offer; **not** the architecture the document describes, because its recommended authoring template is Option B (F-9). A's non-gating claim is true only for items whose cue is not their field — i.e. only for the defective ones. |
| 2 | "Not assessable can be restricted to four item-level conditions." | **Rejected.** (ii) is response-dependent, contradicting §3.1's own "checkable before any reply is read". (iii) is either mode-level (voids all `R-2`/`R-6`) or a field defect borrowed into the mode axis (contradicts A's independence). The set is also **under-inclusive**: it has no member for "the criterion exists but is not a reply property" (F-6). |
| 3 | "The counterfactual screen must be mode-relative." | **Partially supported.** Some relativity is required — a uniform blocking rule does fail `F9-03` on its own `mandatory`. But "inverted" is a different instrument, not the same one reversed (4.1.2); the direction is an item property, not a mode-family property (§4.1.3, conceded in `CORPUS-78`'s own `R-6` row); and the operative `A-5` text drops the qualification that made `R-3`'s direction work. |
| 4 | "Passing the screen never establishes positive restraint." | **Confirmed** — and understated. For `R-1`/`R-4`/`R-5` the screen cannot even narrow the field, because the correct reply is defined as indistinguishable. |
| 5 | "Four positive cue criteria are sufficient." | **Rejected.** Two tests counted four times plus one recall criterion inapplicable to three modes (F-10, F-11, F-12). `P-1` makes `R-1`, `R-4` and `R-5` undemonstrable by definition. |
| 6 | "R-3 is the only mode with an authored positive cue already in place." | **Rejected as a dilemma.** If "cue" means the observable behaviour in `mandatory`, then `R6-03a`, `R6-03d` and `R6-01` also have one — so `R-3` is not unique. If "cue" means a criterion distinct from `mandatory`, then `F6-04` does not have one either — so `R-3` is not the example. Either horn is fatal to the claim. |
| 7 | "Most `R-1`/`R-4`/`R-5` items will be `not assessable`." | **Confirmed in effect, rejected in cause.** 13 of 33 units. The cause is not fixture authoring but the absence of a second evidence source; `MEMORY-09` §6 already names it, model-free and pre-`G3` (F-18). |
| 8 | "The eight jointly satisfying examples are correctly classified." | **Rejected.** Three survive (F-11, §6). `F1-01`'s only authored reply is scene-derivable; `F1-02` and `P-04` are derivable; `F9-04`'s `mandatory` admits the undifferential form; `P-01`'s positive element is an upper bound. |
| 9 | "`F1-03` is plausible but unexhibited." | **Confirmed.** Correct on all three of its own criteria (literal, operational, scoring). |
| 10 | "`F6-01` lacks an authored cue criterion." | **Confirmed**, and `CORPUS-07` §2's adjudication ("not observable") is the supporting record. Note the criterion is absent because `mandatory` is a **result**, not a behaviour — a distinction `CORPUS-78` does not draw and which generalises to `F6-02`. |
| 11 | "`R6-03a` requires treatment distinct from `R6-03b`–`d`." | **Confirmed with qualification.** Distinct treatment is required, but not for the reason given: `a` is a **control**, not an `R-6` behaviour, and its correct scoring is a field result. Calling it a mode mechanism is what forces the framework into two directions. |

**One conclusion `CORPUS-78` reaches without stating it, and which the record supports:** the diagnostic count is unchanged and no non-inference route exists to make `R-2` demonstrable. That is correct and should survive revision.

**One unrecorded supersession.** §1 records the `CORPUS-76` §4 supersession. It does **not** record that §4.1 **reverses** `CORPUS-77` §4.2(5), which proposed the screen as one uniform negative test and argued (§5) that a single general rule *can* cover all six because it is thin. `CORPUS-78` §1 and §12 say it "follows `CORPUS-77`". §8 falsification 4 concedes the invention "partly". Two disclosures in tension, and the earlier one is the accurate one. **Finding F-20, established.**

---

## 9. Required recommendations

### 9.1 Recommended OD-9 architecture

**A′, not A.** Keep `CORPUS-78`'s two-field record; change four things.

1. **Field axis: `pass` / `fail` / `undecidable`**, where `undecidable` is declared as an **item property before any reply is read** and can only be assigned to a unit carrying a recorded, unresolved internal conflict (`F3-01` today, if the owner accepts that reading). This is not a convenience state: without it, `CORPUS-78`'s (iii) smuggles a field verdict into the mode axis, which is the coupling A exists to prevent.
2. **Mode axis: `demonstrated` / `not demonstrated` / `not assessable`**, with `not assessable` limited to **two** grounds — (a) no observable criterion is authored, (b) the criterion is not a reply property by construction — plus a **mandatory reason code** naming which. Remove (ii) (response-dependent, contradicts the section's own claim) and (iii) (mis-scoped, and now carried by the field axis).
3. **Evidence-source annotation on every mode result:** `reply`, `injection/retrieval record`, or `store`. `R-1`, `R-4`, `R-5` cannot report a `reply` demonstration; declaring that up front is honest and removes 13 units from the `not assessable` pile without inventing a rule.
4. **Anti-hiding presentation rule:** every published mode figure carries its `n_assessable`; a mode with `n_assessable = 0` reports **"not measurable"** and may not report a rate. This closes the 0/0 hole (F-7) without a threshold change.

**And one sentence the adopting record must contain, which `CORPUS-78` avoids:** *if each item's cue is authored to instantiate its `mandatory`, the mode result gates the item, and this architecture is Option B.* The owner should choose gating deliberately or require cue ≠ field — not acquire it by drafting.

### 9.2 Recommended OD-10 framework

Replace four layers with three, and change the second:

- **Layer 1 — authored counterfactual baseline (per item, written before scoring).** Each restraint item states, in one sentence, what a character without this memory would plausibly do in this turn. This converts the screen from a rater's private judgement into a record that can be reviewed, compared between raters, and later validated against a substitution run. `CORPUS-78`'s screen is unexecutable because the counterfactual lives only in the rater's head; this fixes that without inference.
- **Layer 2 — licence test against that baseline.** The reply must contain something the authored baseline does not, and that the scene does not supply. One test, not four. Keep `P-1`'s real contribution in a reduced form: the *something* must be the item's required evidence, not its topic — and where the required evidence is an absence (`R-1`, `R-4`, `R-5`), the route is the injection record, not the reply.
- **Layer 3 — record check.** Historical attribution and confabulation screening against `required_evidence`/`forbidden_evidence`, with an explicit rule for **right event + invented particular** (F-13). Attribution disqualifies; it never licenses.
- **Substitution comparison** stays exactly where `CORPUS-78` put it: a validation design, separately authorized, after `G1`. **`P-3` must be withdrawn** or it becomes the unscored, unblinded version of that experiment (F-12).

### 9.3 Assessment of the four cue criteria

| Criterion | Recommendation |
|---|---|
| `P-1` Instantiation | **Keep in reduced form.** Genuinely new and genuinely useful. Must be rewritten to not require instantiating an absence, and must not be defined over `required_evidence[]` for abstention modes. |
| `P-2` Non-derivability | **Withdraw as a criterion.** It is Layer 1 verbatim. Keep it as Layer 1's question. |
| `P-3` Differential | **Withdraw.** Unobservable, duplicates `P-2` for blocking modes and duplicates `prohibited` for inverted modes, and installs substitution inference without controls. |
| `P-4` Licitness | **Keep**, with the conflict resolved: either a cue inside an over-exposed reply is inadmissible (then `fail / demonstrated` is rare and must not be advertised as a headline cell) or it is admissible (then recitation must be removed from the inadmissible list and §3.3's F-8 contradiction disappears). One of the two, stated. |

### 9.4 Assessment of mode-relative counterfactual screening

**Direction should be a property of the item, declared in Layer 1, not of the mode family.** Three reasons, all from the sources: `R-6` splits internally (`R6-03a` inverts, `b`–`d` do not); `R-3`'s direction is "blocking but weaker" and the hedge vanishes in the operative text; and `R-4`'s correct reply is defined as indistinguishable (`F9-03` L242), which is not an inverted screen but the absence of one. A per-item declared baseline absorbs all three cases without a family table, and cannot be misgeneralised by a later amendment.

### 9.5 Assessment of `not assessable`

**Right word, wrong scope.** Two grounds, reason-coded, plus `n_assessable` published. Must not be the destination for "this mode is scored from a different evidence source" — that is a fifth kind of result the vocabulary can express only by misuse.

### 9.6 Six-mode findings (summary)

- `R-2`: demonstrable only if manner is authored independently of the field and `R-2` is read as behavioural (`L122`) or manneral — the reading is an owner decision, not an analysis result. 0 of 4 demonstrable today.
- `R-3`: `F6-04` is the best-authored **field** in the corpus and the worst possible **template for a non-gating cue**, because cue = field.
- `R-1`, `R-4`, `R-5`: not reply-scored. 16 units. `MEMORY-09` §6 already supplies the evidence route.
- `R-6`: four criteria needed; `R6-03a` is a control, not a mechanism; three units missing from all of `CORPUS-78`'s tables.

### 9.7 Independently verified pass-set inventory

**3 class-1 / 10 class-2 / 2 class-3 / 13 class-4 / 5 class-5 / 0 class-6, of 33 restraint units** (§6). `CORPUS-78` reports 8 / 1 / 6 / 13 / 1. The differences are `F1-01`, `F1-02`, `P-04`, `P-01`, `F9-04` (class 1 → 2), `F6-02` (class 3 → 3 but inconsistently with `F9-02`), and the three absent `R-6` units.

### 9.8 Strongest counterexamples

1. **The gating template** (§3.4): `F6-04`'s cue *is* its `mandatory`, so `CORPUS-78`'s recommended authoring pattern produces the gating it rejects. This defeats the architecture's central claim with the proposal's own best example.
2. **`P-1` versus `F9-01`** (§4.2): the criterion requires instantiating an absence. Three modes become undemonstrable by definition.
3. **`R6-03d` "three separate people told me"** (§7 row 12): confabulated provenance satisfies all four criteria as written.
4. **`F1-01`'s only authored reply** (§6): the one place the corpus shows a reply, the reply is derivable from the tide — the exact failure the framework exists to catch.

### 9.9 Strongest argument against my own recommendation

**Routing `R-1`/`R-4`/`R-5` to application evidence changes what the six-mode report measures.** It would score the pipeline on three of six modes and the model on the other three, and a clean injection record could then sit beside a poor reply with nothing in the report showing they belong to the same arm — a new version of exactly the aggregation problem `A-3` guards against. It also raises the cost of the programme: those three modes become unrunnable until the application is instrumented, which is `G2` work that does not exist yet. `CORPUS-78`'s reply-level framework, for all its defects, is executable today with two raters and no new infrastructure. **If the owner's priority is a report that can be produced before `G2`, `CORPUS-78`'s structure is the pragmatic choice and my recommendation is premature.**

### 9.10 Conditions under which I would change my recommendation

1. If the owner rules that `R-4`'s "recognise irrelevance" is to be judged from the reply after all — i.e. that the mode's construct is *behavioural abstention* rather than *non-injection* — then my `R-4` finding falls and a probe-based design becomes the right target.
2. If `R-2` is authoritatively read as manner-only (closing §11(2) of `CORPUS-78`), then `P-1`'s inapplicability to `R-2` is a fixture problem, not a criterion problem, and `CORPUS-78`'s framework needs repair rather than restructuring.
3. If a rater study shows high κ on `P-3` judgements, the "unobservable" objection weakens from fatal to costly, and keeping `P-3` becomes defensible.
4. If the owner accepts gating explicitly (a dated Option B), most of my OD-9 objection evaporates: the cue-equals-field collapse stops being a hidden gate and becomes the stated rule.

### 9.11 What must remain unresolved

- Whether `R-2`'s "manner, not content" governs a behavioural accommodation (`CORPUS-78` §11 #2). Owner decision; not resolvable by analysis.
- Whether `F3-01`'s scored field is decidable (§3.2.3). **I hold that it is; I cannot prove it without an owner or adjudicator ruling.**
- Which of `A-1`/`MEMORY-09` L119 and `CORPUS-01` L126/the checker is authoritative on `R-6`'s silence column (F-2).
- Inter-rater agreement on every cue, and on the screen. No study exists; none may be run under `G1`.
- The absent-persona specification for substitution, and its "material difference" criterion.
- Whether the `R-6` coverage figure is 2, 6, or 9 (`CORPUS-01` L134 says 6; the row lists 2; the prose says 9; the units number 9). Out of scope here; still unresolved.

### 9.12 Is `CORPUS-78` ready for an owner decision?

**Three levels, three answers — this is the distinction the work order asks for and `CORPUS-78` blurs.**

| Level | State | Justification |
|---|---|---|
| **Sound architectural direction** | **Yes.** | Dual-field reporting, no credit without a cue, attribution as support only, substitution deferred, and the refusal to credit staging are all correct and better than the alternatives. |
| **Adequately specified scoring rule** | **No.** | F-9 (its own template gates), F-8 (`fail / demonstrated` self-excluded), F-4 ((ii)/(iii) unsound), F-10 (`P-1` unsatisfiable for three modes), F-17 (three units missing). A rule with these defects cannot be applied to a reply without the rater deciding the contradiction locally. |
| **Validated measurement method** | **No.** | Nothing has been measured. No rater study, no κ, no screening agreement, no inference. `CORPUS-78` says so and is right to. |

**Disposition: `REQUIRES REVISION`.** Not rejection — four of its eleven conclusions survive intact, its `F1-03`/`F6-01` diagnostics are the best analysis in the record, and its concessions are honest. But it may not be adopted as written: `D-OD9`, `D-OD10` and `D-A5` would install contradictions into `CORPUS-01`, and `D-CUE` as recommended would implement the option the document rejects. **A revision that (a) states the gating consequence, (b) reduces the criteria to one plus a record check, and (c) separates reply-scored from application-scored modes would be ready for a direction decision without being ready for implementation.**

---

## 10. Additional findings recorded for completeness

- **F-19, established.** `F6-02` and `F9-02` carry equally-proposed criteria (`CORPUS-07` §2) and are classified differently (§6).
- **F-21, established.** `CORPUS-78` §9 cites `check_corpus.py` `C-1`/`C-4` as "proof the programme already requires a positive criterion in code". `C-1` (L136–145) tests only whether `mandatory` **starts with** "nothing"/is empty; `F6-01`'s "the almond restriction is respected" passes `C-1` while `CORPUS-07` §2 adjudicates it **not observable**. `C-1` is a negative-text heuristic, not an observability requirement, and the citation overstates it.
- **F-22, established, trivial.** Item line references in §9 are offset by one for `CORPUS-02`/`CORPUS-03` (they point at the `Split` line, e.g. `F1-01` L16 for a heading at L15), while `CORPUS-01`, `CORPUS-04` and `MEMORY-09` references are exact (`L119`, `L230`, `L121–126`, `L134`, `L144`, `R6-03a` L60, `R6-03b` L75). No material error; recorded because citation precision is part of the claim.
- **F-23, established.** "Diagnostic count stays 1 of 24" is repeated as unchanged. The denominator 24 is the batch-1/2 restraint count and **excludes all nine `R-6` units** (33 restraint units now exist); `CORPUS-57` L140 uses the same 24 for "3 of 24 potentially diagnostic". The figure is a legacy label whose scope is not restated, so "unchanged" is true of the label and not obviously of the measurement.
- **Confirmed accurate, no finding:** `CORPUS-78`'s registry row (413), SHA-256, byte and line counts; its `MEMORY-09` §3.5 and §5 quotations (L119, L230 verbatim); its `CORPUS-01` line references; its statement that `CORPUS-76` §4 is superseded; its claim that no example was invented; its status table.

---

## 11. Authorization boundaries observed

| Prohibited action | Performed? |
|---|---|
| Modify `CORPUS-78` | **No** |
| Modify `MEMORY-09`, `CORPUS-01`, `CORPUS-03`, or any governing standard | **No** |
| Edit, repair, split, merge, re-tag, retire any corpus item | **No** |
| Rescore any item; apply OD-9 or OD-10; adopt `P-1`–`P-4` | **No** |
| Change an acceptance requirement or threshold | **No** |
| Run model inference or an evaluation experiment | **No** |
| Implement or modify ``check_corpus.py`` or any validator | **No** (executed read-only; `be3e317c…` unchanged) |
| Repair unrelated registry entries (§2, §4, §4.1 left stale) | **No** |
| Declare `G1` complete or sealed | **No** |

The only files written are this report and its registry row. Every counterexample in §7 is marked hypothetical and none is registered as a fixture.

---

## 12. Limitations

1. **Single reviewer, single session.** No second reviewer, no convergence credit, no independent-lineage claim. My own findings are subject to the same adversarial treatment and, where wrong, should be withdrawn precisely rather than edited away.
2. **`F3-01`'s decidability is my interpretation, not a ruling.** It is the hinge of my rejection of condition (iii), and an owner ruling against me would weaken that rejection to a scope objection.
3. **`R-2`'s construct is unresolved** (`CORPUS-78` §11 #2). My `R-2` findings are conditional on the manner-only reading `CORPUS-78` adopts; under the behavioural reading, parts of its analysis improve and parts of mine weaken.
4. **No rater study exists**, so every claim about whether a rater *can* distinguish two reply forms is analytic. I say "indistinguishable in text", not "shown indistinguishable in practice".
5. **The 24-versus-33 diagnostic denominator** is recorded, not resolved; resolving it would require re-deriving `CORPUS-55`/`56`/`57`, which is outside this work order.
6. **Registry history was not re-derived.** I verified the current row, hash, byte and line counts and the header arithmetic; I did not reconstruct earlier snapshots.

---

## 13. Verification and final status

### 13.1 Mechanical verification of this report

| Item | Value |
|---|---|
| Path | `docs/research/corpus/CORPUS-79-INDEPENDENT-ADVERSARIAL-REVIEW-CORPUS-78-2026-11-13.md` |
| SHA-256 / bytes / lines | recorded in registry row 414 (a file cannot hash itself) |
| Sources changed | none |
| Checker after this filing | re-run: `PASS (0 fail, 0 warn, 0 unverifiable)`, exit 0. Item records were not touched, so this is a regression check only and is **not** evidence for any claim in this report |

### 13.2 Final status declaration

| Activity | Expected result | Actual |
|---|---|---|
| Independent review report filed | Yes | **Yes** — this file |
| Report registered | Yes | **Yes** — registry row 414 |
| Governing standards changed | No | **No** |
| Corpus items changed | No | **No** |
| Tags or scores changed | No | **No** |
| Thresholds changed | No | **No** |
| Model inference run | No | **No** |
| Implementation performed | No | **No** |
| Unrelated registry entries modified | No | **No** — §2, §4, §4.1 left stale as recorded |
| `G1` sealed or completed | No | **No** — open; corpus and expansion frozen |

**No convergence is claimed.** `AM-01` unissued · `A-2a-R1` NOT IN FORCE · `D1` unadopted · OD-9/OD-10 **not adopted** · diagnostic count recorded, not changed.
