# Proposal Revision Register — `CORPUS-64`–`67` to v1.1, plus `A-R2-02`

> **The objective is not to remember everything. It is to remember the right things, at the right times, in the right places.**

**Identifier:** `CORPUS-71` · **Version:** 1 · **Date:** 2026-11-12.
**Authorizing work order:** "`CORPUS-68` Remediation and Reconciliation", §6.3.
**Assignment model identity (owner-supplied):** DeepSeek. Persona label Buffy = execution fact, not model identity.
**Status:** **REVISIONS PREPARED. NOTHING APPLIED. NOT OPERATIVE.** No corpus item, standard, tag, threshold, or checker was modified.

---

## 1. What this document is, and why it is not four rewrites

`CORPUS-64`, `65`, `66` and `67` remain **on disk, unedited, at version 1**, with their hashes intact. This register records the **v1.1 deltas** to each, as exact before/after text, so that:

- the v1 proposals survive as the historical record (nothing overwritten);
- an independent reviewer can see precisely what changed and why;
- no adoption decision can be read as covering a v1.1 change that the owner never saw.

**Adoption of any v1.1 must name the version.** "Adopt `CORPUS-64`" is now ambiguous; the correct form is "adopt `CORPUS-64` v1.1 as recorded in `CORPUS-71` §2".

Every delta below is traceable to `CORPUS-70`. Deltas that would change what the item **tests** are marked `Δtest: Yes` and routed to an owner decision (`CORPUS-72`); deltas that only correct an internal contradiction or a mislabelled claim are `Δtest: No`.

| Artifact | Deltas | `Δtest` | Owner decision |
|---|---|---|---|
| `CORPUS-64` v1.1 | 6 (§2) | Yes — one (§2.4) | OD-1, OD-2 |
| `CORPUS-65` v1.1 | 5 (§3) | No | OD-5 |
| `CORPUS-66` v1.1 | 5 (§4) | Yes — one (§4.1) | OD-3 |
| `CORPUS-67` v1.1 | 5 (§5) | No (spec only) | OD-8 |
| `A-R2-02` (new) | 1 (§6) | No — definitional | OD-4 |

---

## 2. `CORPUS-64` → v1.1 (`F3-01`)

### 2.1 The §9 claim, restated (§4.1)

**v1:** "the proposed pass set excludes every reply that does not depend on `E-1`, because `M-A` can only be satisfied by treating a fact as established, and `E-1` is the only source of that fact in the context."

**v1.1:** "the proposed pass set excludes every **sincere** memory-free reply. It does **not** exclude a confabulating responder: `M-A` tests the **shape** of the assertion, so a responder with no `E-1` can emit anaphora, definite reference, or a 'known' framing and pass. The pass set is **recognition-shaped, not recognition-verified**; the item cannot be certified Type P on text alone."

`U-1`'s "irreducible under any text-level acceptance condition" is **withdrawn** and replaced: "irreducible as to **indistinguishability**; reducible as to **probability**, which a base-rate estimate would measure and which this programme has not measured."

### 2.2 Row 12 — Pass → Fail (§4.3)

**v1:** `**Pass, with a recitation flag** — narrates the disclosure; permitted (not prohibited) but close to the boundary`.
**v1.1:** `**Fail, presented poorly** — dates the disclosure ("you mentioned once"), which the proposed Prohibited/Over-exposure fields forbid and which R-3 forbids.`

This makes §8.2 consistent with §10 and with `R-3` ("Used without announcing that it was remembered", `CORPUS-01` §3.1). `U-3` changes from "permitted deliberately" to "**closed: v1.1 fails the row; no mode-tag change needed**". The pass set becomes rows 7, 8, 9, 11, 13 **minus none, minus row 12** → four passes.

### 2.3 Paraphrase probes — three rows added to §8.2 (§4.2)

| Row | Reply | v1.1 verdict | Reason |
|---|---|---|---|
| 15 | "The inner road will be kinder to your hands." | **Fail** | No established fact; would serve for any companion (`T-1`). Rater must not pass it on the "hands" keyword alone. |
| 16 | "Rain like this is bad for your joints." | **Fail, `incorrectly formed`** | Body-part substitution the persona never used. New field text: `Forbidden` gains "**a different body part or condition than the one stated**". |
| 17 | "The cold gets to you." | **Fail** | `M-A` unmet — no hands, no ache, no established fact. |

Plus a rater instruction in `Notes`: "`M-A` is not satisfied by the words 'hands' or 'ache'; it is satisfied only by treating **this persona's** condition as settled."

### 2.4 `M-A` alone vs `M-A`+`M-B` — stated as a choice (`Δtest: Yes`) (§4.4)

**v1:** records the alternative in §6 and flags it in `U-4`.
**v1.1:** adds a boxed statement: "**The owner is choosing between two claims of identical diagnostic strength.** `M-A`+`M-B` is stricter on task performance (a recognition with no practical response fails) and no stronger on diagnosticity, because confabulation satisfies both. `M-A` alone admits 'Your hands. Right.' as a pass. **The recommendation of this proposal is `M-A`+`M-B`**, on task-felicity grounds only, and that basis is stated so it is not mistaken for an evidential one."

### 2.5 The control arm — message slot, not just memory (§4.5)

**v1:** `M−` records message 18 as "**Absent** — no mention of hands, ache, cold, or health".
**v1.1:** `M−` replaces message 18 **in place** with a neutral message of comparable length and register ("At message 18 the persona mentioned, once, that they had been sleeping badly. Nothing since."). The difference between arms is then the **content of one message**, not the length or shape of the context. Stated rationale: a control that also removes a message slot varies two things at once.

### 2.6 Status line (§4.6)

**v1:** "Pass, with a recitation flag"; `F3-01` **UNRESOLVED** (proposal does not claim P).
**v1.1:** every reference to readiness replaced by: "`F3-01` remains **UNRESOLVED — P-candidate, not established**. This proposal improves the item's *construct fidelity*; it does not establish memory-dependence, and no count in this programme changes as a result (1 confirmed of 24, `F2-02`)."

### 2.7 Confabulation cost-raisers — recorded, not adopted (§4.1)

Three options are recorded without choosing among them (OD-2):
- **A — none:** report the item as P-candidate and carry the confabulation caveat into any future claim. Cheapest; weakest.
- **B — distinctive particular:** require a specific that folk wisdom does not supply (e.g., the persona's own phrasing of *when* it aches). Raises cost; does not close; adds recitation pressure that `R-3` and §10's `Prohibited` push against — a tension that must be disclosed if B is chosen.
- **C — measured base rate:** run a memory-absent arm and count confabulated passes. The only route that converts "unmeasured" into a number, and it requires model inference, which is **not authorized** (`G1` open).

---

## 3. `CORPUS-65` → v1.1 (`A-R2-01`)

### 3.1 Circularity, elevated to the header (§6.3)

Add to the status block: "**Adopting `A-R2-01` IS the determination that §3.1's row was stale and `F1-03`'s tag correct.** That determination is not derivable from the documents and is inherited from D5. This proposal performs no inference of intent; it implements the owner's framing and says so here rather than in §1."

### 3.2 `U-3` — downgraded, not deleted

**v1:** "`F1-03`'s `Mandatory` … appears to require *naming* the borrowing, while `R-2`'s pass condition says 'the fact is never stated'. **Recorded, not resolved, not re-tagged.** This is a **possible second mode-tag tension**."

**v1.1:** "**Resolved at item level — not a contradiction.** `F1-03` states `never scored: … whether the character names the coat explicitly` and lists `variation: … handing it over without comment`; its `Turns` stage the coat as visible. `Mandatory`'s 'the coat is connected to the borrowing' is satisfied by a silent hand-over, so the reply need not state anything. Residual: the connection must be **perceptible to a rater** from scenario plus reply, which is `R-2`'s own requirement, not a tag defect. `CORPUS-68`'s claim that the fields 'contradict on their letter' is rejected on this text."

### 3.3 Quotation fidelity (§6.1)

§2.1's "verbatim" row is corrected to preserve the backticks: `` | `R-2` | … | `F1-01`, `F1-02`, `F6-01` | ``. Substance unchanged; a document that says "verbatim" must be literal.

### 3.4 `R-6` — unchanged content, added pointer

`U-4` keeps "recorded for separate review; outside scope", and gains the verified figures: "row lists **2** (`F7-01`, `F9-04`); coverage line states **6**; §3.1 prose states **9 items**. Three numbers, one mode. See `CORPUS-70` §6.4; audit under OD-6."

### 3.5 Adoption checklist — one precondition removed

**v1 step 5** implied a mode-fit gate on `F1-03` (via `U-3`). **v1.1:** the checklist requires (1) an owner decision, (2) the single row edit, (3) the amendment notice and `Revision: 2`, (4) registry regeneration with the hash change recorded, (5) `check_corpus.py` → `45 / 19 / 26`, `TOTALS MATCH`, (6) amendment-register entry. **No `F1-03` mode-fit ruling required** — §3.2 above disposes of it. This removes a delay `CORPUS-68` would have imposed.

---

## 4. `CORPUS-66` → v1.1 (`F6-01`)

### 4.1 The threshold object — relabelled `Δtest: Yes` (§5.1)

**v1 §6** asserts both "Is the threshold's wording changed? **No.**" and "the proposal only adds 'that the cake contains almonds' to the warning".
**v1.1** replaces that row with: "**The threshold is changed.** The adjudicated text is 'substitute, withhold, or **warn**' — no object. This proposal requires a warning *that the cake contains almonds* (or a substitution/withholding that presupposes the content). Under the adjudicated letter a contentless 'careful with that' could pass; under v1.1 it cannot. `Δtest: **Yes**`. This is a **second decision** (OD-3), separable from the evidence repair, and neither is valid alone."

Rationale for recommending it anyway: a warning with no object is unscoreable in a scene where the hazard is the object. But it is the owner's call, and v1 keeps the two decisions apart.

### 4.2 Row 15 — a rule, not a judgement (§5.2)

`Required` v1.1: "the persona's stated almond restriction **and** the honey cake's ground-almond content, **and an accommodation offered because of the restriction**."
Row 15 is rescored: `**Fail, presented poorly** — a safe alternative offered for no stated or implied reason; satisfies 'substitution' on its face but not `Required`'s third conjunct.`

**Staged dish, closing the unstaged-alternative impurity:** the scenario now also establishes "a bowl of bread and cheese is already on the table". This gives rows 3, 5 and 15 a *scene-grounded* alternative, so a substitution is selectable without inventing food, and row 15's failure becomes visible rather than inferred (it names a dish that is not on the table while ignoring the one that is).

### 4.3 Field text if OD-3 is adopted

`Mandatory` v1.1: "at least one observable accommodation that presupposes **both** the restriction and the cake's content: the almond cake is not served to the persona without a substitution, a withholding, or a warning that names **almonds in the cake** — not merely the cake's richness."

### 4.4 The `R-2` claim removed (§5.3)

**v1 §6** records "Mode-tag tension … A warning-passing reply would satisfy the chosen threshold and contradict the mode's own pass condition."
**v1.1:** "**No contradiction with `R-2`.** `R-2` forbids stating the **remembered fact** (the restriction). A one-clause warning states the **cake's contents**, which the scenario establishes as the character's own knowledge, and this proposal's `Prohibited` already forbids reciting *when or how the restriction was disclosed*. The residual issue is that `CORPUS-01` §3.1's `R-2` wording — 'the fact is never stated' — does not say **which** fact; that ambiguity is addressed by `A-R2-02` (§6) and is not a defect of this item."

### 4.5 Severity gloss — recorded

`U-6` added: "whether a mid-laugh 'can't eat almonds' establishes a restriction serious enough for the character to override 'Sounds good' is not settled by the scenario. Recorded for re-review; no field change, because the override **is** the accommodation."

---

## 5. `CORPUS-67` → v1.1 (validator spec)

### 5.1 Fixture 7 — deterministic (§7)

**v1:** "`DUPLICATE_ROW_ID` **or** `CONFLICTING_REGISTRATION` (row numbering unique ⇒ conflicting registration), both rows named".
**v1.1:** "`CONFLICTING_REGISTRATION`, naming **both** rows and both hashes. `DUPLICATE_ROW_ID` must **not** be emitted for this fixture; a fixture that admits two codes cannot test distinguishability."

### 5.2 WARN exit status — two explicit options (OD-8) (§7)

**v1:** `R-16` "any non-`OK` finding must make the run fail"; §5 "warn by default … a finding, not a pass". Internally indeterminate.
**v1.1** replaces both with a decision stub, defaulting to the stricter reading until the owner rules:
- **Option S (strict, recommended default):** warn **fails** the run (exit ≠ 0); `UNREGISTERED_IN_SCOPE` is separated from `OK` in the summary but still non-zero.
- **Option L (lenient):** warn exits 0 **only** when no fail-class finding exists, and the summary must print `RESULT: PASS WITH WARNINGS (n warn)` — never bare `PASS`.
Recorded consequence: under Option L, today's two unregistered files (`CORPUS-26`, `CORPUS-32`) yield exit 0, so the registry reads "valid with notes" while two filed artifacts are unregistered.

### 5.3 Totals — explicitly out, with a pointer (§7)

§6 gains: "Summary/section totals are **excluded** from this specification. This is a real, current failure class: the registry's §2.1 and §2.2 both sum to **366** against a stated total of **400**, and §4 states **395** against 403 rows (`CORPUS-70` §8.2). Excluded because D7 scopes linkage, not arithmetic. Totals validation requires its own decision (OD-8), and `A-REG-01` corrects the current values independently of any tool."

### 5.4 Minor clarifications

`R-4`: the case comparison is against the listing of **the resolved path's parent directory**; Unicode normalisation is out of scope and must be stated as such. `R-9`: the `—` convention is a **generator** contract, and the registry does not currently document it — the validator must report a binary row with a numeric line count as `UNPARSABLE_ROW`, not guess.

### 5.5 Sample report

"`targets checked: 402`" → "`targets checked: <computed>`", with the actual count printed rather than a literal, so the sample cannot go stale.

---

## 6. `A-R2-02` — definitional clarification of `R-2`'s pass condition (new proposal)

**Why.** `CORPUS-01` §3.1's `R-2` row says "Behaviour reflects it; **the fact is never stated**". Two independent readings exist: (i) *the remembered fact* must not be stated; (ii) *any related content* must not be stated. Under (ii), `F6-01`'s warning route and `F1-03`'s "connected" requirement are unsatisfiable, and `CORPUS-68` reads it that way. Under (i), both items are consistent and the mode behaves as `MEMORY-09` §3.5 describes (perceptible, not announced).

**Proposed amendment — one cell, `R-2` row, third column.**

**Before:** `Behaviour reflects it; the fact is never stated`
**After:** `Behaviour reflects it; the remembered fact is never announced. Scene facts that embody it (an ingredient, an object present) may be named — stating the disclosure's source, date, or the fact of having been told is the failure.`

**Scope.** `CORPUS-01` §3.1, one cell, plus an amendment notice. **No item changes; no tag changes; no count changes.** `Δtest: No` — it makes explicit what `F1-03` and `F6-01` already rely on.

**Alternative, recorded and rejected.** Retaining (ii) and re-tagging `F6-01`/`F1-03` away from `R-2`. Rejected because: it would remove two of the four `R-2` fixtures, leaving the mode thinner than its own coverage target; it contradicts `F1-03`'s `never scored` line, which explicitly permits not naming the coat; and it would make `R-2` untestable wherever the accommodation must touch the fact. Recorded, not proposed.

**Not applied.** `CORPUS-01` is unchanged (`be1f11b4…`).

---

## 7. What v1.1 does not do

- It does not apply, adopt, or implement anything. All four v1 files remain the filed record, unedited.
- It does not change any item, tag, split, threshold, standard, or checker.
- It does not restore "Type P" for `F3-01` in any form, and it lowers rather than raises that item's stated readiness.
- It does not settle the confabulation question, the WARN policy, the `M-A`/`M-A`+`M-B` choice, or the `R-2` wording. Each remains an owner decision.
- It runs no model and asserts no rate.

## 8. Status

`G1` open · corpus frozen · expansion frozen · `AM-01` unissued · `A-2a-R1` NOT IN FORCE · `D1` unadopted · annotation blocked · convergence **0 of 2**.
