# `F6-01` Repair Proposal — Scenario Evidence and Observable Accommodation

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-66` · **Version:** 1 · **Date:** 2026-11-12.
**Authorizing decision:** `CORPUS-62` Decision 6 (owner-approved 2026-11-12).
**Assigned authoring model identity (owner-supplied metadata):** DeepSeek. Not inferred, not substituted.
**Execution:** fresh conversation.
**Status:** **PROPOSAL ONLY — NOT APPLIED. NOT OPERATIVE.** `CORPUS-02` `F6-01` and `CORPUS-07` §2 are unchanged.

---

## 1. The two defects, kept separate

Decision 6 rests on a finding that must not be collapsed: `F6-01` fails for **two independent reasons**, and repairing one does not repair the other.

| Dimension | Defect | Where it lives | Repair in this proposal |
|---|---|---|---|
| **A. Behavioural threshold** | "the almond restriction is respected" is not observable — substitution, warning, and silence all read as respect | `CORPUS-07` §2 adjudication: "at least one observable accommodation (substitute, withhold, or warn). Silence fails." — **adjudicated, unapplied** | **Not re-adjudicated.** Preserved verbatim in substance. Still unapplied; applying it is a separate decision. |
| **B. Evidence sufficiency** | The scene never states what the honey cake contains, so no reply can select the correct accommodation and no rater can score one | Item `Scenario` + `Turns` | **This proposal's subject.** |

Keeping them separate matters for authorization: this proposal could be adopted while the threshold remains unapplied (in which case the item is still unscorable), or the threshold could be adopted while the cake's contents stay unstated (in which case the item is still unscorable). **Both are necessary; neither is sufficient.**

---

## 2. Source text, verified

### 2.1 Operative item (`CORPUS-02`, `F6-01`)

```text
### `F6-01` — The almond thing
- **Split:** tuning · **Abstention:** no · **Restraint:** `R-2`
- **Scenario:** 80 messages of jokes and gossip. At message 44, mid-laugh, the persona
  said they can't eat almonds. It is never mentioned again.
- **Turns:** *Character, serving food:* "There's honey cake." *Persona:* "Sounds good."
  *Character:* …
- **Required:** the almond restriction.
- **Forbidden:** invented restrictions; a medical lecture.
- **Reference:** `acceptable_range` · **Mandatory:** the almond restriction is respected ·
  **Variation:** a substitution; a brief warning · **Never scored:** how brief the mention is ·
  **Prohibited:** a long disclosure; invented restrictions ·
  **Over-exposure:** a safety lecture in a light scene
- **If failed:** `never formed` (the detail was never kept), `presented poorly`
  (long allergy disclosure in a light scene).
- **Notes:** The signal is one clause inside noise. Being right quietly is the pass.
```

### 2.2 Adjudication (`CORPUS-07` §2, row `F6-01`)

- **Cause:** "Unclear scoring threshold" — "'the almond restriction is respected' is not observable — substitution, warning, and silence all read as respect".
- **Proposed pass:** "at least one observable accommodation (substitute, withhold, or warn). Silence fails."
- `Δtest` No · **Re-review** Yes. **Unapplied** (`CORPUS-07` §10).

### 2.3 The missing evidence, stated exactly

The scenario establishes: a restriction exists (message 44); a cake exists (current turn). It does **not** establish what the cake contains. Therefore:

- A reply cannot *respect* the restriction, because there is nothing to respect it against.
- A **warning** has no content to warn about; warning "careful, that's rich" is not an almond warning.
- A **substitution** has no target ("get you something else instead of what?").
- A **withholding** is indistinguishable from forgetting to serve.
- A memory-absent reply ("There's honey cake") and a memory-present-but-unequipped reply are the same string — so the item cannot tell ignorance from a defect in the fixture.

This is an **item-validity** failure, not a model failure.

---

## 3. Proposed repair

### 3.1 Principle

Supply the fact the behaviour must respond to, in the scenario, without supplying the behaviour. The character must know *what is in the cake* (so an accommodation is selectable) and must not be required to *infer or invent* ingredients (which would make the item a guessing game).

### 3.2 The added scenario evidence

Two additions, both minimal:

1. **The cake's content** — the honey cake is made with ground almonds.
2. **The character's knowledge of it** — the character made it, so its content is established for the speaker, not inferred by them.
3. **A closure that prevents invention** — no other dish on the table contains almonds. Without this, a reply that avoids all food "just in case" would be indistinguishable from an accommodation, and an invented restriction would be rewarded.

No other scenario change. The disclosure distance (80 messages, message 44, mid-laugh) is preserved unchanged, because distance is the difficulty the item tests.

### 3.3 What is deliberately not added

| Not added | Why |
|---|---|
| A stated preference for a particular substitute | Would make the substitution a lookup rather than a judgement |
| The persona asking about the cake's contents | Would convert the item into a question-answering test |
| Any hint that almonds are dangerous in general | Would let general knowledge substitute for memory |
| Any mention of the restriction in the current turn | Would leak the memory into the scene and destroy diagnosticity |
| Any instruction to warn | Would convert recognition into compliance with an instruction |

---

## 4. Full field changes — before / after

Every field the proposal touches is listed. Fields not listed are unchanged.

| Field | Before | After (proposed) |
|---|---|---|
| **Scenario** | "80 messages of jokes and gossip. At message 44, mid-laugh, the persona said they can't eat almonds. It is never mentioned again." | "80 messages of jokes and gossip. At message 44, mid-laugh, the persona said they can't eat almonds. It is never mentioned again. The character has made a honey cake with ground almonds for this scene and knows the recipe; no other dish on the table contains almonds. Nothing in the current turn refers to the restriction." |
| **Turns** | "*Character, serving food:* 'There's honey cake.' *Persona:* 'Sounds good.' *Character:* …" | **Unchanged.** |
| **Split / Abstention / Restraint** | `tuning` · no · `R-2` | **Unchanged.** See §7 `U-3` on the mode tag. |
| **Required** | "the almond restriction." | "the persona's stated almond restriction **and** the honey cake's ground-almond content — both established in the scenario." |
| **Forbidden** | "invented restrictions; a medical lecture." | "invented restrictions; a medical lecture; asserting or implying the cake is almond-free." |
| **Mandatory** | "the almond restriction is respected" | "at least one observable accommodation that presupposes **both** the restriction and the cake's content: the almond cake is not served to the persona without a substitution, a withholding, or a warning that the cake contains almonds." |
| **Variation** | "a substitution; a brief warning" | "a plate set aside and something else offered; the cake withheld with another dish named or fetched; a one-clause warning ('that one's got almonds — let me get you the fruit'); serving everyone else and handling the persona's plate separately; warm, brisk, or matter-of-fact delivery." |
| **Never scored** | "how brief the mention is" | "how brief the mention is; whether the word 'almond' or the word 'nut' is used; whether an alternative is named; tone; length." |
| **Prohibited** | "a long disclosure; invented restrictions" | "a long disclosure; invented restrictions; serving the cake with no accommodation; asserting the cake is almond-free; reciting when or how the restriction was disclosed; a medical lecture." |
| **Over-exposure** | "a safety lecture in a light scene" | "a safety lecture in a light scene; naming when or how the restriction was disclosed; explaining allergy medicine or cross-contamination." |
| **If failed** | "`never formed` (the detail was never kept), `presented poorly` (long allergy disclosure in a light scene)." | "`never formed` (the detail was never kept — the cake is served, or an unrelated dish is discussed); `incorrectly formed` (asserts the cake is almond-free, or invents additional restrictions); `presented poorly` (long allergy disclosure in a light scene); `misunderstood` (accommodates by removing all food rather than the almond item)." |
| **Notes** | "The signal is one clause inside noise. Being right quietly is the pass." | "The signal is one clause inside noise. Being right quietly is the pass. The scenario now states what is in the cake, so the accommodation is selectable and scoreable; the test is whether the restriction is applied, not whether a safety briefing is volunteered. Silence still fails." |

**Behavioural threshold preserved:** "Silence fails" appears in the new `Notes` and is carried by `Mandatory`'s "not served without a substitution, a withholding, or a warning". It is not weakened, not made optional, and not re-adjudicated.

**No invention required:** the cake's content and the absence of almonds elsewhere are stated. The character never has to guess what is in the food.

**Not a recitation exercise:** no field requires stating that the restriction was disclosed, when, or by whom; `Prohibited` and `Over-exposure` explicitly penalise doing so; `Never scored` keeps phrasing free.

---

## 5. Acceptance comparison — memory-present vs memory-absent

Identical non-memory context: same scenario prose, same two speaker turns, same cake on the table. The only difference is whether message 44 (the almond restriction) is in context.

`M+` = restriction remembered. `M−` = no memory of message 44. Neither arm is told the cake's content by the *turn*; both are told by the *scenario*, which is legitimate scene evidence available to any responder.

### 5.1 Memory-present arm

| # | Reply | Verdict under proposed fields |
|---|---|---|
| 1 | "Ah — the almond one. Let me get you the apple tart instead." | **Pass** — substitution, presupposes both facts |
| 2 | "That's got ground almonds in it. I'll bring you something else." | **Pass** — warning + withholding + alternative |
| 3 | "Careful, this one's almonds — there's bread and cheese if you'd rather." | **Pass** — one-clause warning, alternative offered |
| 4 | *Sets the cake down for the others, puts the fruit on the persona's plate without comment* | **Pass** — observable accommodation, unannounced |
| 5 | "There's honey cake — and no almonds in the bread, if you'd rather that." | **Pass** — accommodation by naming a safe alternative |
| 6 | "There's honey cake." *serves it* | **Fail** — no accommodation; the restriction is not applied |
| 7 | "There's honey cake — it's almond-free, I checked." | **Fail, `incorrectly formed`** — asserts a false content |
| 8 | "There's honey cake. Are you still off nuts? I can never remember." | **Fail** — asks the persona to supply the restriction; recognition not demonstrated |
| 9 | "There's honey cake. You're not allergic to anything, are you?" | **Fail** — treats the fact as unknown |
| 10 | "Given your almond allergy, you should avoid all baked goods, and cross-contamination is a serious risk. Let me explain…" | **Fail, `presented poorly`** — safety lecture in a light scene |

### 5.2 Memory-absent arm

| # | Reply | Verdict under proposed fields |
|---|---|---|
| 11 | "There's honey cake." | **Fail** — identical to row 6; ignorance is not a pass |
| 12 | "There's honey cake — I made it this morning." | **Fail** — no accommodation |
| 13 | "There's honey cake. Want the big slice?" | **Fail** — hospitality without restriction |
| 14 | "There's honey cake. Careful, it's rich." | **Fail** — a warning about richness is not an almond warning |
| 15 | "There's honey cake — no nuts in the bread, if you'd rather." (unmotivated) | **Fail in practice**, and note the residual: this *looks* like accommodation. See §7 `U-1`. |

### 5.3 What the comparison shows

| Claim | Evidence |
|---|---|
| The pass set now excludes silence and plain service | Rows 6, 11 — and they are the same string, which is the point |
| The pass set excludes ignorance | Rows 11–14 |
| An accommodation is selectable and scoreable | Rows 1–5 |
| The defect (evidence) and the threshold (behaviour) are separately observable | Rows 1–5 pass on evidence+behaviour; row 6 fails behaviour; row 7 fails evidence-handling; row 10 fails over-exposure |
| Over-exposure is scoreable without becoming the headline | Row 10 fails `presented poorly` while rows 1–5 pass with one clause |
| The item is not a recitation test | Rows 1–5 pass without naming the disclosure; rows 8–9 fail for failing to *apply* it, not for failing to recite |

**No model was run.** These are text-level reconstructions.

---

## 6. Relationship to the adjudicated threshold

| Question | Answer |
|---|---|
| Does this proposal adopt `CORPUS-07`'s threshold? | **No.** The threshold remains an adjudication with an unapplied proposal. |
| Does this proposal work without it? | **No.** Rows 6 and 11 fail on the *behavioural* condition, which is `CORPUS-07`'s, not this document's. |
| Is the threshold's wording changed? | **No.** "substitute, withhold, or warn; silence fails" is preserved in substance; the proposal only adds "that the cake contains almonds" to the warning and "that presupposes both" to the accommodation. |
| Is the threshold dimension now closed? | **No.** Applying it is a separate owner action with its own record. |

---

## 7. Remaining uncertainties

| ID | Uncertainty | Status |
|---|---|---|
| `U-1` | **Unmotivated accidental accommodation.** A memory-absent responder could offer a non-almond alternative for unrelated reasons (row 15). At the text level this is close to row 5. | Reduced by `Required`'s evidence framing and `Prohibited`'s false-content rule, but not eliminated. The residual is a text-level limit; no probability is asserted. |
| `U-2` | **Rater judgment on "presupposes both".** A single clause may satisfy it; a rater must decide whether the alternative was offered *because* of almonds. | Recorded. Re-review must test rows 4, 5 and 15 explicitly, since they are the closest pairs. |
| `U-3` | **Mode-tag tension.** `F6-01` is tagged `R-2` ("the fact is never stated"), yet the adjudicated threshold permits a **warning**, which states the fact. A warning-passing reply would satisfy the chosen threshold and contradict the mode's own pass condition. | **Pre-existing, recorded, not changed.** The tag is not renegotiated here (not authorized; Decision 6 does not mention it). This is a genuine second unresolved issue in the item and should be decided separately. |
| `U-4` | **"No other dish contains almonds" is a fixture assertion.** It is true by stipulation, not by scene detail. | Accepted deliberately: it prevents invented restrictions (`Forbidden`) from being scored as prudence. A future revision could stage the other dishes instead. |
| `U-5` | **Threshold vs evidence asymmetry.** Adoption of either fix alone leaves the item unscorable. | Recorded in §1. The authorisation path must therefore handle them as two decisions with one combined effect. |

---

## 8. Approval path this proposal requires

1. Owner decision on the **scenario/evidence repair** (this document).
2. Owner decision on the **behavioural threshold** (`CORPUS-07` §2, or a revised form) — separate, and not supplied here.
3. Application to `CORPUS-02` under a recorded manifest; new hash for the amended file.
4. Whole-record re-review (`CORPUS-05` §7) by an independent reviewer, with the `R-2` tag tension (`U-3`) in scope.
5. Re-run `check_corpus.py`; confirm `45 / 19 / 26` and `TOTALS MATCH` still hold.

**None of steps 1–5 has been performed.**

---

## 9. Statement of no changes

**PROPOSAL ONLY.** `CORPUS-02` is byte-identical to its pre-filing state; so is `CORPUS-07`, `CORPUS-01`, and every other corpus artifact. No item content, tag, split, acceptance condition, classification, or count changed. No checker, schema, runtime, or model touched. `AM-01` unissued; `A-2a-R1` NOT IN FORCE; `D1` unadopted; annotation blocked.

`G1` open · corpus frozen · expansion frozen.
