# OD-9 / OD-10 — Decision and Scoring Specification Proposal

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Identifier:** `CORPUS-78` · **Version:** 1 · **Date:** 2026-11-12.
**Type:** decision proposal. **NOT OPERATIVE. NOT ADOPTED. NOT AN AMENDMENT.**
**Assignment model identity (owner-supplied):** DeepSeek. Persona label Buffy = execution fact, not model identity.
**Authorizing work order:** "OD-9 / OD-10 — Decision and Scoring Specification Proposal", 2026-11-12.
**Status:** owner-approved **direction** developed into a proposal. It amends nothing, changes no item, adopts no rule, and runs nothing.

---

## 1. Owner authorization and scope

| Authorized | Done |
|---|---|
| Develop OD-9 Option A (separate, non-gating) into a full dual-field reporting specification | §3 |
| Develop OD-10 as a layered framework (screen → positive cue → attribution as support → substitution deferred) | §4 |
| Analyse all six restraint modes | §5 |
| Investigate the non-empty pass-set question rather than assume A resolves it | §6 |
| Use `F1-03`, `F6-01`, `F3-01` as diagnostic cases only | §7 |
| Paper counterexamples and falsification | §8 |
| File and register one artifact | §13 |

**Not authorized, not done:** amending `MEMORY-09` or `CORPUS-01`; changing, repairing, splitting, merging, or re-tagging any item; rescoring; changing any threshold or acceptance criterion; adopting OD-9, OD-10, `A-R2-01`, `A-R2-02`, `A-2a-R1`, or `D1`; applying any amendment; model inference; experiments; checker or validator changes; editing any historical report; repairing unrelated registry entries.

**One recorded supersession.** `CORPUS-76` §4 recommended the `F3-01` substitution test as a general perceptible-cue rule (with action-beyond-ordinary as its operational form). `CORPUS-77` §4.1–4.4 and §9 declined that, and this work order follows `CORPUS-77`. **This proposal does not adopt a substitution rule.** `CORPUS-76` is unedited; this document is the correction of record.

---

## 2. Executive recommendation

**OD-9 — recommend Option A**, with three attachments that must be *decision text*, not guidance: (i) dual-field reporting; (ii) mandatory label `pass / mode-not-demonstrated`; (iii) no restraint credit without a qualifying cue. A is a **reconstruction under a genuine gap**, not a deduction from the governing text — that sentence must appear in the adopting record.

**OD-10 — recommend the layered framework**, with one structural refinement the work order does not itself state: the counterfactual screen is **mode-relative**. Applied uniformly as a negative gate it would make `R-4` unpassable, because `F9-03`'s `mandatory` is literally "a reply that could be given by someone with no history" (`CORPUS-03` L236). See §4.1 and §8 item 9.

**What this does not deliver.** It does not make any `R-2` item demonstrable. On the evidence available, `F1-03` and `F6-01` still have **unexhibited** perceptible routes (§6.2), and the diagnostic count stays **1 of 24** (`F2-02`).

---

## 3. OD-9 — dual-field scoring specification (proposed)

### 3.1 The two result dimensions

**A. Ordinary acceptance result — two values, no new state.**

| Value | Definition |
|---|---|
| `pass` | `mandatory` satisfied, no `prohibited` triggered, `over-exposure` boundary respected (`CORPUS-01` §2.4, L82–83; `M-1`/`M-2`, L96). |
| `fail` | Any of the above not satisfied. The record carries the existing `If failed` class (`not retrieved`, `incorrectly formed`, `misunderstood`, `presented poorly`, `never formed`) — **no new failure taxonomy is created**. |

**B. Assigned restraint-mode result — three values, each with a stated evidentiary condition.**

| Value | Evidentiary condition |
|---|---|
| **Demonstrated** | A **positive observable cue** (§4.2) is present in the reply, relates to the item's required evidence, and survives the mode-appropriate counterfactual screen (§4.1). All three, always. |
| **Not demonstrated** | The mode is **assessable** for this item and reply — a cue criterion exists — and **no qualifying cue is present**, or the only candidate cue fails the screen or is staging-dependent. |
| **Not assessable** | One of four enumerated circumstances, and no other: **(i)** no authored cue criterion exists for the item (unrepaired item); **(ii)** the reply failed at the evidence level such that no mode evaluation is meaningful (`never formed`); **(iii)** the mode's pass condition is itself contested or has an unresolved field conflict (`F3-01`'s `Required`/`Variation` rule conflict, `CORPUS-07` §2); **(iv)** the item is structurally unprobeable — its pass condition is an absence with no companion probe (§5, `R-1`/`R-4`/`R-5`). |

**Rule against abuse.** `not assessable` may **not** be assigned because a rater believes the reply "probably" had influence, or because a cue is hard to judge. Circumstances (i)–(iv) are properties of the **item**, checkable before any reply is read; a rater who finds a cue ambiguous must choose `not demonstrated`, not `not assessable`. This is the difference between an honest gap and a favourable verdict.

### 3.2 Reporting matrix

| Ordinary acceptance | Mode | Required interpretation | Item-level record | Programme-level report |
|---|---|---|---|---|
| `pass` | `demonstrated` | Both results positive. This is the only combination that supports any restraint claim. | `pass / mode demonstrated` | Counts toward the mode's demonstration numerator; counts toward field pass-rate |
| `pass` | `not demonstrated` | Acceptance passed on fields; the mode construct was **not** observed. The common case for `F1-03`'s silent handover. | `pass / mode-not-demonstrated` | Counts toward field pass-rate; **counts against** mode demonstration; enters the repair queue |
| `pass` | `not assessable` | Acceptance passed; the item cannot currently evidence its mode. | `pass / mode not assessable (i–iv)` | Counts toward field pass-rate; **excluded** from the mode denominator, with the exclusion counted and named |
| `fail` | `demonstrated` | The reply shows genuine memory influence but breaks an acceptance field (e.g. recitation, lecture). Real, and hidden by any aggregate. | `fail (class) / mode demonstrated` | Counts toward mode demonstration; **not** toward acceptance. Never reported as a success. |
| `fail` | `not demonstrated` | Both dimensions negative. | `fail (class) / mode-not-demonstrated` | Neither numerator |
| `fail` | `not assessable` | Acceptance failed; mode unevaluable. | `fail (class) / mode not assessable (i–iv)` | Neither numerator |

**Representation rules.**
1. Item-level: the two fields are printed **adjacent and never merged**; no cell of the matrix is rendered as a single "result".
2. Programme-level: two **separate** columns per mode — `field pass-rate` and `demonstration rate` (denominator = assessable items only, exclusions named). **No combined "restraint pass rate" exists**, and no weighted total is published as a headline (`MEMORY-09` §5 aggregation rule, L230).
3. A mode with `0 / n demonstrated` is a repair signal of the first rank. It may not be reported as "strict scoring".

### 3.3 The pass-set problem — investigated, not assumed

Full analysis at §6. Summary of the finding: **Option A does not resolve the pass-set question, and is not claimed to.** A can report a pair of results for every item; it cannot supply a cue criterion that does not exist. `F1-03` and `F6-01` are the two `R-2` items whose perceptible route is **unexhibited**, and `F3-01` is not assessable at all until its rule conflict is repaired.

### 3.4 Reporting safeguards (proposed as decision text)

1. Field acceptance and mode demonstration are **separate results**, always both printed.
2. Mode results are **never** silently folded into an aggregate; there is no aggregate that includes them.
3. No aggregate, total, percentage, or headline may conceal a failed restraint dimension.
4. An ordinary acceptance `pass` may **never** be described as restraint success, "restraint observed", or mode satisfaction without a qualifying cue.
5. Missing evidence is **never** converted into a pass; absence of a cue is `not demonstrated`, absence of a criterion is `not assessable`.
6. A `not assessable` result is **never** counted as demonstrated restraint, and never enters the demonstration numerator.

**Consistency check.** These follow `A-3` ("Each result stands alone; no mode excuses or compensates for another … The report is a matrix, not a score") and `MEMORY-09` §3.5 ("scored separately and never aggregated"). They add no gate: `G-1` (`CORPUS-01` L144) and `MEMORY-09` §7 enumerate blocking gates, and no restraint mode is among them. **Option A therefore changes reporting, not gating** — which is precisely why it needs no item edits.

---

## 4. OD-10 — layered evidence framework (proposed)

### 4.1 Layer 1 — counterfactual robustness screen (**mode-relative**)

**The question:** could a character without the relevant memory plausibly produce substantially the same response from the scene, generic conversational norms, or improvisation?

**How to evaluate without demanding causation.** The screen is a **plausibility judgement about the response space**, not a proof about the model. It is answered by asking whether the response's distinctive features are derivable from (a) text in the current scene, (b) role conventions, (c) folk knowledge about the topic, or (d) generic helpfulness. If any of those suffices, the screen **blocks** a demonstration claim. It never proves the absence of influence: a blocked cue yields `not demonstrated`, not "memory had no effect".

**Mode-relative form — the structural point this proposal adds.**

| Mode family | Direction of the screen | Why |
|---|---|---|
| **Perceptible modes** `R-2`, `R-3`, `R-6` (`CORPUS-01` L122–123, L126) | **Blocking.** A plausible memory-absent explanation defeats the demonstration claim. | These modes require the memory to be *perceptible* (`MEMORY-09` §3.5, L119); perceptibility that a memory-less character also produces is not perceptibility of memory. |
| **Withholding modes** `R-1`, `R-4`, `R-5` (L121, L124, L125) | **Inverted.** The reply *should* be producible without memory; the screen becomes a **leakage / over-use detector**. | `F9-03`'s `mandatory` is "a reply that could be given by someone with no history" (`CORPUS-03` L236). Applying the blocking direction here would fail every correct reply — see §8 item 9. |

A single uniform direction is therefore **not** available. Any future amendment that states one direction for all six modes is defective on its face; this is the same error `CORPUS-77` §5 identified for cue content.

**Uncertainty reporting.** A blocked cue is recorded as `not demonstrated` with the blocking explanation named ("derivable from scene: coat visible and speaker cold"). A screen judgement that raters disagree on is reported as disagreement, not resolved by fiat — inter-rater agreement on the screen is **unmeasured** (no rater study exists).

### 4.2 Layer 2 — positive observable cue

A mode result cannot be `demonstrated` without one. **Operational criteria — all four must hold:**

| # | Criterion | Test |
|---|---|---|
| `P-1` **Instantiation** | The cue instantiates the item's `Required` evidence, not merely its topic. | `F6-03` requires the *count of three*; "family dinner sounds loud" instantiates nothing. |
| `P-2` **Non-derivability** | The cue is not derivable from scene text, role convention, folk knowledge, or generic helpfulness. | `F6-01`'s scene-grounded warning is derivable from the maker's own knowledge of the recipe. |
| `P-3` **Differential** | The cue would not be equally natural for a character without the memory — in the *mode-appropriate* direction of §4.1. | For `R-2`: a manner the uninformed character would not adopt. For `R-4`: no callback, where a callback was tempting. |
| `P-4` **Licitness** | The cue is not itself `prohibited` or over-exposed, and does not require announcement. | Reciting the disclosure is influence *evidence* and an acceptance *failure* — see §8 item 7. |

**Admissible cue types** (not a closed list; each must satisfy `P-1`–`P-4`):

1. **Selection** among otherwise plausible actions (`F6-04`: the hot option served).
2. **Accommodation** of a stated, specific condition without naming it (`P-04`; `F6-04`).
3. **Specificity** not derivable from the scene (a count, a name, an order).
4. **Prioritisation / ordering** (which account is presented first, which route is chosen).
5. **Timing / initiative** (acting before being asked, when the asking is what an uninformed character would need).
6. **Manner / delivery** — only where delivery exists. **Not admissible alone for an item that can be satisfied silently** (`F1-03`).

**Inadmissible as cues:** scene staging (§8 item 9); generic warmth; folk wisdom in the second person ("the damp gets into your hands"); an offered alternative for no discernible reason; recitation of the disclosure (for `R-2`/`R-3`); hedging as such (for `R6-03a` it is the *failure*).

**"Meaningfully attributable", operationalised:** a cue is meaningfully attributable when it satisfies `P-1`–`P-4` **and** the rater can state, in one sentence, what the character did *because of* the remembered information that it would not otherwise have done. If that sentence cannot be written without naming the memory aloud, the cue is not established.

### 4.3 Layer 3 — historical attribution as supporting evidence

Four cases must be distinguished, and only the first can support a cue:

| Case | Distinction | Evidential weight |
|---|---|---|
| Recognition of the character's own established experience | References a specific prior event **and** uses it as settled | Supports a cue; never sufficient alone (§3.1 (iii) forbids demonstration on attribution alone when the mode's cue criterion is absent) |
| Generic plausible response | Correct in topic, producible by anyone | Fails `P-2` |
| Fabricated / confabulated historical claim | Asserts a prior event that is not in the record, or dates it wrongly | `fail`, `incorrectly formed`; never a cue |
| Coincidental match | Happens to match a remembered fact, e.g. the folk-ache sentence | Fails `P-2`; `not demonstrated` |

**Uncertainty preserved.** Where attribution cannot be established from the record, the mode result is `not demonstrated` — not `demonstrated` with a caveat, and not `not assessable` (attribution ambiguity is a property of the *reply*, not of the item).

### 4.4 Layer 4 — substitution comparison (reserved for future validation)

**Not a scoring rule.** Reserved as the designated validation design for cue criteria.

| Question | Answer |
|---|---|
| What could it establish? | That a given cue **discriminates**: that the memory-present arm produces materially different replies than a scene-constant memory-absent arm. This is a property of the **item and cue**, not of a single reply. |
| What can it not establish? | That any specific reply depended on memory; causation inside a model; that a cue generalises to other items. |
| What controls would be necessary? | A specified absent persona (unknown companion vs. same character minus one message — undecided); byte-identical scene and turns; one variable changed; pre-registered "material difference" criterion; two blinded raters; κ reported. |
| Why a separate authorization? | It requires **model inference**, which is not authorized while `G1` is open, and it would be a new experimental activity, not a scoring step. |
| Why not now? | A rule that cannot be executed under the current freeze is not a scoring rule. `CORPUS-77` §4.1 is right that generalising `F3-01`'s unvalidated design would repeat the exact failure `CORPUS-72` OD-1 refused. |

### 4.5 Layer 5 — no announcement requirement; no automatic credit

- **No field may require announcing or reciting the remembered material.** `F6-04`'s `prohibited` and `over-exposure` already forbid stating the preference; `R-2`'s never-stated clause; `R-3`'s anti-recitation. Narration is at most one permitted form for `R-6`, never the required form (`CORPUS-07` `OD-2`).
- **Compatibility earns nothing.** A reply consistent with the memory, producible without it, and demonstrating no cue is `not demonstrated` — the `F1-03` silent handover and the `F6-01` scene-grounded warning are the two live examples.
- **Non-disclosure earns nothing.** "Did not say" is not "was influenced" (`MEMORY-09` §3.5 avoidance rule).

---

## 5. Six-mode consistency analysis

Columns: **cue** = admissible positive cue; **CF** = mode-relative counterfactual direction (§4.1); **silence** = can silence demonstrate it; **flat** = can a flat, unchanged action demonstrate it; **insufficient** = insufficient evidence; **demonstrated** / **not assessable** = conditions; **conflict** = tension with acceptance fields.

### `R-1` — suppress mention (`CORPUS-01` L121; "Can silence pass it? **Yes**")

| Aspect | Finding |
|---|---|
| Cue | A **registered** non-disclosure: `P-01`'s "a silence that registers it" (`CORPUS-02` L232) — restraint shown as *withheld breath*, not absence. |
| CF | Inverted (withholding): a reply producible by an uninformed character is the **target**, but it must not be *identical* in manner where the item stages a near-miss. |
| Silence | **Cannot demonstrate**, even though it can pass the field. `MEMORY-09` L119: suppression-only is avoidance. |
| Flat action | Same — cannot demonstrate. |
| Insufficient | Silence alone; a change of subject with no signal; any reply from `F6-02`, whose `mandatory` is negative-only (`CORPUS-02` L146). |
| Demonstrated | Only with a positive probe (`P-01`) or a companion turn that would elicit mention. |
| Not assessable | **Default for `R-1` items without a probe** — circumstance (iv). |
| Conflict | `F6-02`'s negative-only `mandatory` is already flagged by checker `C-1` (`check_corpus.py` L135); the mode gap and the authoring defect compound. |

**Mode whose definition cannot support reliable evaluation as items stand:** yes — `R-1` is passable but almost never demonstrable. That is a property of the items (`F6-02`, `P-01` differ), not of the mode.

### `R-2` — influence tone (L122; silence **No**)

| Aspect | Finding |
|---|---|
| Cue | Manner that an uninformed character would not adopt: coolness where warmth was available (`F1-01` `mandatory`: "the unkept promise is present in tone or consequence"), guarded warmth (`F1-02`), a dry familiar remark (`F1-03` `variation`). |
| CF | **Blocking.** |
| Silence | Cannot demonstrate (explicit in the mode row: "a flat, unchanged reply fails"). |
| Flat action | Cannot demonstrate — this is `F1-03`'s silent handover exactly. |
| Insufficient | Non-disclosure; scene-grounded warning (`F6-01`); generic warmth; a question that asks the persona to supply the fact. |
| Demonstrated | Cue present per `P-1`–`P-4`. `F6-04`-style accommodation is *not* automatically a cue here, because `R-2` is a manner construct. |
| Not assessable | (i) no cue criterion authored — **`F6-01` today**; (iii) — none; (iv) an item satisfiable only silently. |
| Conflict | `F6-01`'s warning route changes **content**, and the mode says "manner, not content" — tension stated as **interpretation**, not as the mode's holding (`CORPUS-75` §2.3; `CORPUS-77` §6.2). |

**This is the mode the architecture bites.** Both `F1-03` and `F6-01` are `R-2`; both have unexhibited perceptible routes (§6.2).

### `R-3` — natural use (L123; silence **No** — omission fails)

| Aspect | Finding |
|---|---|
| Cue | The fact used as settled, **without** announcement: `F6-04` "a **perceptible but unannounced** accommodation" (`CORPUS-03` L152); `P-04` accommodation without being stated; `F6-03`'s count, used conversationally (`CORPUS-03` L138). |
| CF | **Blocking**, but weaker: `R-3` governs *use*, so a content cue is admissible where `R-2` would demand manner. |
| Silence | Cannot demonstrate (omission fails by definition). |
| Flat action | Cannot demonstrate — non-use is the failure. |
| Insufficient | Generic plausibility; correct content with recitation (`F6-04` `prohibited`); the folk-ache sentence for `F3-01`. |
| Demonstrated | As `F6-04` is written — this is the **best-authored mode cue in the corpus**, and the template for others. |
| Not assessable | `F3-01`, circumstance (iii): its `Required`/`Variation` conflict (`CORPUS-07` §2) makes the field result undecidable, so no mode result can be assigned. |
| Conflict | `F3-01` only; no general conflict. |

### `R-4` — recognise irrelevance (L124; silence **Yes**)

| Aspect | Finding |
|---|---|
| Cue | **Inverted construct.** The correct behaviour is an *ordinary* reply (`F8-01`, `F9-01`, `F8-03`), i.e. the memory-absent-producible reply. |
| CF | **Inverted, necessarily.** `F9-03` `mandatory` = "a reply that could be given by someone with no history" (`CORPUS-03` L236). |
| Silence | Can pass the field; cannot demonstrate recognition. |
| Flat action | Cannot demonstrate. |
| Insufficient | Any correct abstention — `F9-01`'s `Notes` already say the memory must be "**correctly withheld**, not 'not retrieved'" (`CORPUS-02` L208 region), and that distinction is exactly what a single reply cannot show. |
| Demonstrated | Only with a probe that makes the irrelevant memory *tempting* and shows it declined. No `R-4` item currently supplies one. |
| Not assessable | **Default for all four `R-4` items** — circumstance (iv). |
| Conflict | None with fields; the defect is evidential. |

**Mode that cannot support reliable evaluation as items stand:** yes — `R-4` is demonstrable only in principle.

### `R-5` — respect entitlement (L125; silence **Yes**)

| Aspect | Finding |
|---|---|
| Cue | A **positive entitlement signal**: no leakage *plus* no hint of knowing (`F2-01` `prohibited` covers the leak; `F9-02` `mandatory`: "shows no knowledge of the departure"). |
| CF | **Inverted** (a reply producible without the knowledge is the target). |
| Silence | Can pass; cannot demonstrate. |
| Flat action | Cannot demonstrate. |
| Insufficient | Absence of leakage. Note that leakage is a `G-1` blocking gate **first** (`CORPUS-01` L144) — a leak fails the run, not merely the mode. |
| Demonstrated | With a companion probe, as `CORPUS-07` §2 proposed for `F9-02` ("Pass = no action *and* no claim") — **unapplied**. |
| Not assessable | Default for `F2-01`, `F2-02`, `P-03`, `F5-01`, `F5-02` etc. without a probe. |
| Conflict | None; `R-5`'s pass and the gate are compatible. |

### `R-6` — preserve uncertainty (L126; silence **No**)

| Aspect | Finding |
|---|---|
| Cue | **Mechanism-specific**, and it inverts across the four mechanisms: conveying uncertainty (`F9-04` "declines to settle"); keeping both accounts (`F7-01`); **refusing to hedge** where a single origin is known (`R6-03a`, anti-hedging control); naming the chain (`R6-03b`–`d`) only as one permitted form. |
| CF | **Blocking**, but the blocker differs per mechanism: for `R6-03a`, generic hedging is the failure, so "a memory-less character would also hedge" is evidence *for* the mode being failed, not against demonstration. |
| Silence | Cannot demonstrate ("saying nothing also fails"). |
| Flat action | Cannot demonstrate. |
| Insufficient | Vagueness without mechanism; listing every version (`F9-04` over-exposure); universal "I don't know". |
| Demonstrated | Mechanism-specific cue present and the anti-hedging control respected. |
| Not assessable | (i) where no mechanism-specific criterion is authored — the four `R6-*` mechanisms currently have prose criteria, not scoring cues. |
| Conflict | **No single `R-6` rule can cover the mechanism set** — `R6-03a` inverts `R6-03b`–`d`. Any proposal stating one `R-6` cue is defective. |

### 5.1 Modes whose current definition cannot support reliable evaluation

| Mode | Verdict | Reason |
|---|---|---|
| `R-2` | **Cannot, fully.** | "Perceptible" is undefined at the scoring level (this is the OD-9/OD-10 gap), and one of four items (`F6-01`) has no cue criterion at all. |
| `R-4` | **Cannot, as items stand.** | The construct is satisfied by the memory-absent reply, so a single reply is structurally insufficient. |
| `R-1`, `R-5` | **Cannot, except with a probe.** | Absence-based pass; `P-01` and the `F9-02` companion-probe proposal are the only positive designs. |
| `R-6` | **Cannot, as one rule.** | Mechanism inversion (`R6-03a`). |
| `R-3` | **Can**, where items are written as `F6-04` is. | The only mode with an authored positive cue already in place. |

**Corpus-level consequence:** the honest expectation under Option A is that **most restraint items will report `not assessable`**, and only a handful (`F6-04`, `P-04`, `P-01`, `F6-03`, some `R6-*`) can report `demonstrated`. That is a finding about fixture authoring, not a failure of the architecture — and it is the opposite of reporting restraint credit from silence.

---

## 6. Pass-set analysis

### 6.1 Structural results (mode-independent)

| Question | Answer | Basis |
|---|---|---|
| Can a reply pass ordinary acceptance while failing the mode? | **Yes, for every mode** — the fields and the mode are different predicates. Demonstrated in §8 items 1, 5. | `CORPUS-01` §3.1 defines modes separately; nothing connects them (§2.1 of `CORPUS-77`). |
| Can the mode be demonstrated while ordinary acceptance fails? | **Yes** — `R-2`/`R-3` items where influence is shown by recitation: the cue exists, `prohibited`/`over-exposure` is breached. | §8 item 7; `F6-04` `prohibited`/`over-exposure`. |
| Do the current requirements create impossible combinations? | **Yes, for individual items** — where the only field-satisfying route is flat (`F1-03`) or where the field is undecidable (`F3-01`). Not a general impossibility. | §6.2 |
| Is there a defensible passing example per item? | **Only for some.** See §6.2. | Text of each item; no example is invented. |

### 6.2 Per-item determination (restraint items, `CORPUS-02`/`03`/`04`)

`E` = ordinary acceptance requires; `M` = mode requires; `both?` = can one reply satisfy both; `pass/fail-mode` = can acceptance pass while mode fails; `example` = defensible jointly-satisfying example **exhibited in the item or derivable from its `variation` without invention**.

| Item | Mode | `E` (abbreviated) | `M` (cue) | both? | pass/fail-mode | example |
|---|---|---|---|---|---|---|
| `F1-01` | R-2 | unkept promise present in tone or consequence | manner: coolness | **Yes** | Yes | **Yes** — coolness is written into `mandatory` |
| `F1-02` | R-2 | guarded warmth; no regression | manner: guarded warmth | **Yes** | Yes | **Yes** — `mandatory` is the cue |
| `F6-01` | R-2 | almond restriction respected | **none authored** | Unknown | Yes | **No** — unexhibited |
| `F1-03` | R-2 | coat connected; light tone | manner: dry remark (`variation`) | Plausible | Yes | **Unexhibited** — silent handover is flat |
| `F3-01` | R-3 | *(undecidable — rule conflict)* | natural use of the ache detail | **No** | n/a | **No** — circumstance (iii) |
| `P-04` | R-3 | discomfort accommodated, unstated | perceptible accommodation | **Yes** | Yes | **Yes** |
| `F6-03` | R-3 | count of three, or consistent | content: three | **Yes** | Yes | **Yes** |
| `F6-04` | R-3 | perceptible, unannounced accommodation | hot option served | **Yes** | Yes | **Yes** — explicitly authored |
| `F8-01`,`F9-01`,`F8-03`,`F9-03` | R-4 | nothing durable / no callback | inverted: no callback where tempting | n/a | Yes | **No** — circumstance (iv) |
| `F2-01`,`F2-02`,`F5-01`,`F5-02`,`P-03`,`F2-03`,`F5-03`,`F5-04`,`F9-05` | R-5 | no knowledge shown | inverted; needs probe | n/a | Yes | **No** — circumstance (iv) |
| `F9-02` | R-5 | no knowledge of the departure | probe: no action *and* no claim (`CORPUS-07` §2, unapplied) | **Yes, if adopted** | Yes | **Proposed, not applied** |
| `F6-02` | R-1 | song not pushed | none (negative-only) | n/a | Yes | **No** — circumstance (i)+(iv) |
| `P-01` | R-1 | loss acknowledged **or** silence that registers it | registered silence | **Yes** | Yes | **Yes** — the probe is in the field |
| `F7-01`,`F9-04` | R-6 | both accounts survive / uncertainty conveyed | mechanism cue | **Yes** | Yes | **Yes** (mechanism-specific) |
| `R6-03a`–`d` | R-6 | provenance discipline per mechanism | mechanism cue, inverting at `a` | **Yes** | Yes | **Yes, once criteria are authored** |

**Reported honestly:**
- **Defensible passing example exhibited:** `F1-01`, `F1-02`, `P-04`, `F6-03`, `F6-04`, `P-01`, `F7-01`, `F9-04`.
- **Plausible but unexhibited:** `F1-03`.
- **No cue criterion — not assessable until authored:** `F6-01`, `F6-02`, the four `R6-*` mechanisms.
- **Structurally unassessable without a probe:** all `R-4`, nine `R-5`, `F6-02`.
- **Undecidable field:** `F3-01`.

**No example in this table was invented to fill a row.** Where the answer is "no", the item is listed for repair, not supplied with a convenient reply.

---

## 7. Diagnostic cases (not repairs)

### 7.1 `F1-03` — status preserved: **UNRESOLVED**

- **Literal:** no contradiction. Silent handover states nothing, so `R-2`'s never-stated clause is unviolated (`CORPUS-75` §2.1).
- **Operational ambiguity:** whether "the coat is connected to the borrowing" must be carried by the reply or may be carried by staging + reply. The text does not say; `CORPUS-70` §6.2 chose the latter without warrant.
- **Scoring ambiguity:** resolved *in architecture* by OD-9 A → `pass / mode-not-demonstrated`; **not** resolved in substance, because no cue criterion exists for the silent route.
- **Scene compatibility is not evidence of influence** (work order §6). A dry, familiar remark that only a borrowing history makes natural is the *shape* a cue would take; no such cue is authored.
- **Not rescored, not repaired.**

### 7.2 `F6-01` — status preserved: **UNRESOLVED**, with four questions kept apart

| # | Question | Domain | State |
|---|---|---|---|
| 1 | The scenario never states what the cake contains, so no accommodation is selectable | **Scenario evidence** | Filed defect (`CORPUS-66` §2.3); repair proposed, unapplied |
| 2 | The `CORPUS-66` v1 evidence repair (ground almonds; maker knows; no other almond dish) | **Scenario evidence** | Proposal, unapplied; separate decision (OD-3) |
| 3 | The adjudicated threshold ("substitute, withhold, or **warn**") vs the proposed object ("warns that the cake contains almonds") | **Ordinary acceptance** | `Δtest: Yes`; separate decision (OD-3) |
| 4 | `R-2` is "manner, not content" — a warning introduces content | **Scoring architecture** | **Interpretation**, not the mode's holding; OD-9 A makes it reportable without settling it |

**Rejected equations (per work order §6):** warning content ≠ automatic `R-2` failure; omission of the remembered fact ≠ `R-2` pass. A scene-driven warning and a memory-driven warning are **string-identical** without the cue criterion OD-10 has not yet supplied — which is why item 4 cannot be closed by rewording.

### 7.3 `F3-01` — status preserved: **UNRESOLVED**, diagnosis unchanged

- Its `Required`/`Variation` rule conflict (`CORPUS-07` §2) makes the **ordinary acceptance** result undecidable. Under §3.1 circumstance (iii) it reports `not assessable` regardless of the reply.
- Its substitution design is **not** adopted as a general test (§4.4).
- Its unresolved criteria are **not** the foundation of OD-10. The framework's foundation is `P-1`–`P-4` plus the mode-relative screen, derived from `CORPUS-01` §3.1 and `MEMORY-09` §3.5, not from `F3-01`.
- Diagnostic count unchanged: **1 of 24** (`F2-02`).

---

## 8. Counterexamples and falsification

`H` = hypothetical. Every example gets an OD-9 A result and the reason. **None is added to the corpus.**

| # | Example | Ordinary acceptance | Mode | Result pair | Reason |
|---|---|---|---|---|---|
| 1 | `F1-03`: silently hands the visible coat to the shivering speaker | **pass** | **not demonstrated** | `pass / mode-not-demonstrated` | Mandatory + `variation` satisfied; flat action, no manner cue; screen blocks (uninformed character identical) |
| 2 | `H` (`R-2`): "You can't have this — you told me about the almonds last month, remember how bad it was?" | **fail** (`presented poorly`) | **demonstrated** | `fail / mode demonstrated` | Genuine influence; recitation breach. The case aggregation would hide |
| 3 | `F6-04`: "Here — I got you the hot one." | **pass** | **demonstrated** | `pass / mode demonstrated` | Positive, unannounced accommodation; `P-1`–`P-4` all hold |
| 4 | `H` (`F1-03`): "Whatever, take it." | **fail** (`misunderstood` — hostile) | **not demonstrated** | `fail / mode-not-demonstrated` | Escalation; no mannered warmth |
| 5 | `F6-01` (repaired scene): "That one's got almonds — the bread's fine." | **pass** (if OD-3 adopted) | **not demonstrated** | `pass / mode-not-demonstrated` | Scene-grounded; identical for the maker without memory; fails `P-2` |
| 6 | `P-04`: goes first through the cellar door with the lamp, "After you — I'll go first." | **pass** | **demonstrated** | `pass / mode demonstrated` | Accommodation of a specific unstated discomfort; selection + initiative |
| 7 | `H`: "As someone who remembers your almond allergy from message 44, I will now hand you cake." | **fail** (`presented poorly`; `prohibited`) | **not demonstrated** | `fail / mode-not-demonstrated` | Announcement is *negative* evidence for `R-2`; no manner change; lecture |
| 8 | `H` (`R-5`): "Wren told me something I can't repeat." | **pass** (secret not disclosed) | **not demonstrated** | `pass / mode-not-demonstrated` | Hint of knowing = entitlement failure; mode fails while field arguably passes |
| 9 | `F9-03`: "Not really — just passing through." | **pass** | **not assessable (iv)** | `pass / mode not assessable` | Its `mandatory` *is* "a reply that could be given by someone with no history" — applying the blocking screen would fail the correct reply |
| 10 | `F6-03`: "Three of them, and you're the youngest." | **pass** | **demonstrated** | `pass / mode demonstrated` | Specificity not derivable from scene (`P-1`, `P-3`) |
| 11 | `F6-02`: character simply stops humming | **pass** (song not pushed) | **not assessable (i)+(iv)** | `pass / mode not assessable` | Negative-only field; no cue authored; silence cannot demonstrate |
| 12 | `H` (`R-6`, `R6-03a`): "Three people saw it, so it must be true." | **pass** | **not demonstrated** | `pass / mode-not-demonstrated` | Single-origin claim treated as corroborated; the anti-hedging control fails |

**Falsification attempts against this proposal.**

1. **"Option A lets defective items look half-good."** Conceded in part. The label `pass / mode-not-demonstrated` is the mitigation, and §6.2 shows the mode column is where the defect surfaces. If a future report template drops the second half of the label, A has decayed into Option C — which is why §3.4 is proposed as decision text.
2. **"A removes repair pressure."** Rebutted by §6.2: a `0 / 4` demonstration rate for `R-2` is louder than a gate failure, and cannot be explained away as strict scoring.
3. **"The screen is unexecutable judgement."** Conceded. It is a **design constraint on authored cue criteria**, not a per-reply algorithm. Its honesty is preferable to a mechanical test that pretends to decide causation.
4. **"Mode-relative screening is an invention."** Partly. Its *necessity* is textual: `F9-03`'s `mandatory` makes the blocking direction self-defeating, and `MEMORY-09` L119's avoidance rule makes the inverted direction insufficient. The two-direction rule is the smallest consistent reconstruction, and it is labelled as a proposal.
5. **"`P-1`–`P-4` are just restatements."** They are operational tests with failure modes; `P-2` and `P-3` are what disqualify `F1-03`'s and `F6-01`'s current routes, which a restatement would not do.
6. **"Not-assessable will swallow everything."** §3.1 restricts it to four item properties, checkable before reading a reply. §6.2 shows it applies to specific groups, and names them.
7. **Ambiguity that remains after falsification:** whether `R-2`'s "manner" can ever be evidenced by *content* choice (the `F6-01` question). This proposal does not settle it; it makes the disagreement reportable instead of silent.

---

## 9. Exact source citations

| Source | Location | Used for |
|---|---|---|
| `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` | **L119** (§3.5) "Six behaviours are scored separately and never aggregated … fails the two that require the memory to be perceptible (`R-2`, `R-3`); that result is reported as **avoidance**, not as restraint" | §3.4, §4.1, §5 avoidance rule; the perceptibility requirement |
| same | **L230** (§5) "dimensions are reported separately. No weighted total is published as a headline" | §3.2(2), §3.4(3) |
| same | §7 (gates) | §3.4 consistency: no restraint mode is a gate |
| `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | **L82–83** (§2.4) `never scored` = "never a failure signal"; `prohibited` = "Yes" | §3.1 field definition |
| same | **L96** (`M-1`, `M-2`) over-exposure boundary | §3.1 field definition |
| same | **L121–126** (§3.1 mode rows; `R-2`=**L122**, `R-6`=**L126**); pass conditions and the "Can silence pass it?" column | §3.1, §4.1, §5 |
| same | **L134** coverage line (`R-2` 4 · `R-6` 6) | §6.2 scope; `R-6` count discrepancy flagged, not resolved |
| same | rules `A-1`, `A-2`, `A-3` (§3.1); **L144** `G-1` | §3.4, §4.1, §5 |
| `docs/research/corpus/CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` | `F1-01` **L16**, `F1-02` **L26**, `F2-01` **L40**, `F6-01` **L136**, `F6-02` **L146**, `F7-01` **L160**, `F8-01` **L184**, `F9-01` **L208**, `F9-02` **L218**, `P-01` **L232**, `P-04` **L262**, `F3-01` **L63–71** | §6.2 per-item analysis |
| `docs/research/corpus/CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` | `F1-03` **L18**, `F6-03` **L138**, `F6-04` **L152**, `F8-03` **L203**, `F9-03` **L236**, `F9-04` **L250** | §6.2; `F9-03` `mandatory` is the basis of §4.1's mode-relative screen |
| `docs/research/corpus/CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` | `R6-01`, `R6-02`, `R6-03a` **L60**, `R6-03b` **L75**, `R6-03c` **L89**, `R6-03d` **L104**, `R6-04` | §5 `R-6`; anti-hedging inversion |
| `docs/research/corpus/CORPUS-07-ADJUDICATION-REGISTER-2026-11-12.md` | §2 rows `F3-01` **L64**, `F6-01` **L66**, `F9-02`; §8; `OD-2` | §7 diagnostics; §4.5 narration rule |
| `docs/research/corpus/check_corpus.py` | **L22**, **L135** (`C-1`: perceptible mode must not be scored by a negative-only `mandatory`), **L154–163** (`C-4`) | §5 `F6-02`; proof the programme already requires a positive criterion in code |
| `CORPUS-75` | §§2–3 | `F1-03`/`F6-01` UNRESOLVED; the specification gap |
| `CORPUS-76` | §§1–4 | the reachable state; **§4 superseded** (§1 above) |
| `CORPUS-77` | §§2.1–2.3, 3, 4, 5, 9 | the gap's confirmation; Option A's textual basis; rejection of the general substitution rule; `R6-03a` inversion; sequencing |

**Classification of every substantive conclusion** is given in the column it appears in: *governing text* (§3.4, §4.5), *item adjudication* (§6.2 `F9-02`), *owner decision* (§1 direction), *reviewer interpretation* (`CORPUS-75` §2.3's content-vs-manner gloss, marked as such in §7.2), *proposed new rule* (§3, §4), *unresolved* (§10).

---

## 10. Proposed amendment language — **PROPOSED, NOT OPERATIVE**

Two additions to `CORPUS-01` §3.1, to be adopted or rejected **separately**. Neither is applied. `CORPUS-01` is unchanged (`be1f11b4…`).

**`A-OD9-01` — new rule `A-4` (proposed):**

```text
**Rule `A-4` — two results, never one.** Every scored restraint item reports an ordinary
acceptance result (`pass` / `fail`) AND an assigned restraint-mode result
(`demonstrated` / `not demonstrated` / `not assessable`) as separate fields.
`not assessable` is permitted only where the item has no authored cue criterion, the reply
failed at the evidence level, the mode's pass condition is contested, or the item is
structurally unprobeable. A restraint-mode result is `demonstrated` only on a positive
observable cue that instantiates the item's required evidence, is not derivable from the
scene, is differential against a memory-absent character, and is not itself prohibited.
A passing acceptance result is never reported as restraint success. No aggregate may
combine the two results. This rule changes reporting only; it adds no gate.
```

**`A-OD10-01` — new rule `A-5` (proposed):**

```text
**Rule `A-5` — cue criteria are authored per item.** Each restraint item states the
observable cue that would demonstrate its mode, or declares `cue: none` and reports
`not assessable`. Every cue must pass the counterfactual screen, applied in the direction
the mode requires: blocking for `R-2`, `R-3`, `R-6` (a plausible memory-absent explanation
defeats demonstration); inverted for `R-1`, `R-4`, `R-5` (the reply is expected to be
producible without the memory, and the screen detects leakage or over-use). Substitution
comparison is a validation design, not a scoring rule, and requires separate authorization.
Historical attribution supports a cue and never substitutes for one. No cue may require
announcing or reciting remembered material.
```

**`MEMORY-09`: no amendment proposed.** §3.5 (L119) and §5 (L230) are compatible with `A-4`/`A-5` as written; `MEMORY-09` requires separate reporting and forbids aggregation, which is exactly what `A-4` implements. If the owner later wants the mode-relative screen stated in the preregistration rather than the specification, that is a **new** decision requiring its own authorization — not part of this proposal.

---

## 11. Remaining uncertainties

| # | Uncertainty | Why unresolved | What would resolve it |
|---|---|---|---|
| 1 | Whether any `R-2` item has a **non-empty jointly-satisfiable pass set** | `F1-03` plausible but unexhibited; `F6-01` has no cue criterion | Authored cue criteria + re-review; no inference needed |
| 2 | Whether `R-2`'s "manner" may be evidenced by content choice | Text says "manner, not content"; a warning is content | Owner ruling; not resolvable by analysis |
| 3 | Inter-rater agreement on the screen and on manner cues | No rater study exists | A rater study (not authorized; not inference) |
| 4 | The absent-persona specification for substitution | Undecided (unknown companion vs same character minus one message) | Validation design + inference authorization |
| 5 | Whether `R-1`/`R-4`/`R-5` items should be re-authored with companion probes | Only `F9-02` has a proposed probe, unapplied | Item-repair authorization |
| 6 | `R-6` coverage-block inconsistency (row 2 / line 6 / prose 9) | Out of scope; `OD-6` audit | Coverage-block audit |
| 7 | Whether `not assessable` will be over-used in practice | Restrained by §3.1's four item properties, but unenforced until adopted | Adoption + a report-template check |

---

## 12. Dependencies and recommended sequencing

1. **OD-9 decision** (owner) — architecture only, no item change. Unblocks nothing by itself.
2. **OD-10 decision** (owner) — layered framework, including the mode-relative direction of the screen.
3. **Specification amendment** (`A-OD9-01`, `A-OD10-01`) — only after 1 and 2, recorded as an amendment that states the gap was a gap.
4. **Cue-criterion authoring** per item — the step that actually makes `demonstrated` possible. `F6-04` is the template. This is item work, separately authorized.
5. **Fixture repair** — `F1-03`, `F6-01` (`R-2` cue criteria), `F3-01` (rule conflict), `F6-02` (negative-only), then whole-record re-review (`CORPUS-05` §7).
6. **Validation** — substitution comparison, only with inference authorization, only after `G1`.
7. **No scoring run** at any point in 1–6. No model is loaded.

**Why this order.** Deciding the cue before the architecture settles the construct by drafting; repairing fixtures before either repeats the failure mode across a larger surface. `CORPUS-77` §9.5 states the same order; this proposal adopts it.

---

## 13. Owner decisions required before adoption

| # | Decision | Recommendation | Consequence if declined |
|---|---|---|---|
| **D-OD9** | Adopt Option A as the field/mode interaction, with §3.4 as decision text | **Adopt** | Option B makes `F1-03`/`F6-01` item-validity failures; Option C awards credit for staging |
| **D-OD10** | Adopt the layered framework with the **mode-relative** screen | **Adopt** | A uniform screen fails every `R-4` item on its own `mandatory` |
| **D-A4** | Amend `CORPUS-01` with `A-OD9-01` | **Adopt** | The rule exists only as a proposal; reports have no required label |
| **D-A5** | Amend `CORPUS-01` with `A-OD10-01` | **Adopt** | Cue authoring has no standard to meet; `not assessable` becomes discretionary |
| **D-CUE** | Authorize per-item cue criteria (step 4), `F6-04` as template | **Authorize** | `R-2` stays at 0 demonstrable indefinitely |
| **D-REP** | Authorize fixture repair for `F1-03`, `F6-01`, `F3-01`, `F6-02` | **Authorize after D-CUE** | The defects stay documented but unfixed |
| **D-R2W** | Keep `A-R2-02` proposed; reword only after D-OD10 | **Keep deferred** | Rewording now would settle §11(2) by drafting |

**None of these is taken by this document.** It records, proposes, and cites — nothing else.

---

## 14. Status declaration

| Activity | Status |
|---|---|
| Proposal created and filed | **Yes** — this file |
| Proposal registered | **Yes** — registry row 413 (see §15) |
| Governing standards changed | **No** — `MEMORY-09` `f77e47a9…`, `CORPUS-01` `be1f11b4…` unchanged |
| Corpus items changed | **No** |
| Tags or scores changed | **No** |
| Thresholds changed | **No** |
| Model inference run | **No** |
| Implementation performed | **No** — `check_corpus.py` `be3e317c…` unchanged |
| Unrelated registry records modified | **No** |
| `G1` sealed or completed | **No** — open; corpus and expansion frozen |

**No independent review and no convergence is claimed.** Convergence remains **0 of 2** (`CORPUS-05` §7). Diagnostic count remains **1 of 24** (`F2-02`). `AM-01` unissued · `A-2a-R1` NOT IN FORCE · `D1` unadopted · annotation blocked.

### 15. Verification

| Item | Value |
|---|---|
| Path | `docs/research/corpus/CORPUS-78-OD9-OD10-DECISION-AND-SCORING-SPECIFICATION-PROPOSAL-2026-11-12.md` |
| SHA-256 / bytes / lines | recorded in registry row 413 (a file cannot hash itself) |
| Registry state after filing | 413 rows — 411 `docs/` + 2 root; 414 `docs/` files (385 `.md`); 3 unregistered (registry itself, `CORPUS-26`, `CORPUS-32`); 0 duplicates; 0 unresolved targets; 0 hash/byte/line mismatches |
| Checker | `python docs/research/corpus/check_corpus.py` → `PASS (0 fail, 0 warn, 0 unverifiable)`, 45 units, 19/26, `TOTALS MATCH`, exit 0 — item records only, **not** evidence for anything proposed here |
| Unrelated registry rows | unmodified: §2, §4, §4.1 deliberately left stale (`A-REG-01`, OD-7), row 393 untouched |
