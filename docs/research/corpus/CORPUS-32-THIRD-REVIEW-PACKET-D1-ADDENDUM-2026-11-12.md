# Third Independent Review Packet — `CORPUS-29` + `CORPUS-30` + `CORPUS-31` (D1 Adversarial Chain)

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Status:** packet prepared. **Review not performed.** Not authorization to annotate, repair, allocate an identifier, amend a standard, or adopt `D1`. Nothing here is operative.
**Identifier:** `CORPUS-32`.
**Date:** 2026-11-12.
**Scope:** the complete D1 adversarial chain as submitted — `CORPUS-29` (adversarial review of `CORPUS-28`, filed verbatim, author unnamed, supplied by owner), `CORPUS-30` (adversarial reconciliation, DeepSeek/Buffy, 6-of-7 confirmed +1 rejected, 12 Valid / 5 Defective / 2 Conditional / 1 C-5, `R-12`/`R-13` withdrawn, rule reissued as unnumbered `D1-ENF`), `CORPUS-31` (addendum, same author, supersedes `CORPUS-30` §3.1 counts and §1.2 gate classification: 10 Valid / 7 Defective / 2 Conditional / 1 structurally exposed, zero `Forbidden:`-only and zero `mandatory:`-alone enforcement).
**Authority:** `CORPUS-05` §3–§7 (challenges C-1–C-5, whole-record rule, independence threshold), `MEMORY-09` (governing pre-registration, locked rules), `CORPUS-01` (specification, field definitions, restraint taxonomy).
**Predecessors:** `CORPUS-27` → `CORPUS-28` → `CORPUS-29` → `CORPUS-30` → `CORPUS-31` → this packet.
**`G1` open; corpus and expansion frozen; `D1` not adopted.**

---

## 1. Reviewer instructions

You are independent. You must not be the author of `CORPUS-22`–`CORPUS-31` and must not have participated in their drafting. Read `CORPUS-29`, `CORPUS-30`, and `CORPUS-31` in full, against the governing sources listed in §2 — not against each other's summaries of those sources.

**Withheld deliberately:** the author's design intent, the unapplied repair specification (`CORPUS-28` §8), any expectation about what the counts "should" be, and whether `D1` should pass. If a claim is not recoverable from the item's own fields plus the quoted governing text, that is a finding.

**Order of assessment:** `C-5` (evidence sufficiency) first where invoked (Challenge 7), then mechanism / boundary / scoring questions per challenge. `CORPUS-05` §3: "`ambiguous` is a legitimate verdict and is never rounded to `pass`."

**Whole-record rule (`CORPUS-05` §7):** for every item cited, ask whether `Scenario`, `Turns`, `Required`, `Forbidden`, `Reference` (`mandatory`, `variation`, `never scored`, `prohibited`, `over-exposure`), `If failed`, and `Notes` describe **one** item. A repair that would change one field leaving others asserting the opposite is an enforcement repair, not a copy-edit.

**What presence does not prove:** `check_corpus.py` passing (required fields exist, totals reproduce) proves nothing about whether any element measures what it claims. It tests none of the seven challenges below. Say so explicitly if you ran it.

**If two readings of the same record are both defensible and they score a concrete reply oppositely, the element is ambiguous/conditional — do not resolve it by preference.** Name the reply.

**Conflict disclosure:** state your identity, whether you authored or reviewed any of `CORPUS-22`–`CORPUS-31` or `CORPUS-26`, and whether any part of this packet was unverifiable to you.

### Required output shape

For each Challenge 1–7 below, return:

```text
challenge_id
claim_under_test (one sentence, in your own words)
verdict: confirmed | rejected | partially-confirmed   reason (with source quotes + line/section refs)
counter-reply tested (at least one concrete reply per item-challenge, with scored outcome under each reading)
```

Then a final section:

```text
overall_verdict: REVISE | REJECT | ACCEPT (D1 only; no rule adoption, no annotation authorisation)
blocking_defects (numbered, each reproducible from the record)
non_blocking_notes
limitations (what you could not verify)
```

Do not propose replacement counts without a row-level table. Do not allocate an `R-` number. Do not edit any artifact.

---

## 2. Governing materials (read these, not summaries)

1. `CORPUS-29` (full) — 7 challenges against `CORPUS-28`.
2. `CORPUS-30` (full) — dispositions, rebuilt table §3, revised rule `D1-ENF` §2.
3. `CORPUS-31` (full) — §§1–6 superseding `CORPUS-30` §3.1 and §1.2.
4. `MEMORY-09` rev 2: §2 (`R-1`–`R-10` locked run rules), §2.1 (`R-11`–`R-14` locked safeguards), §3.4 (`G-1` gate text; `M-1` floor rule), §5 (scope table, deletion propagation), §7 (Gates — integrity).
5. `CORPUS-01`: §2 (field definitions — `forbidden_evidence[]`, `pre_registered_diagnostics`, `required_evidence[]`), §2.4 (`mandatory`/`prohibited` scored Yes; `never scored` dimensions incl. confidence), §2.5 (`mandatory` = least that satisfies; `M-1` over-exposure boundary; `M-2` four restraint items `P-01`, `F6-04`, `F1-01`, `P-04`), §3.1 (restraint modes `R-1`–`R-6` + silence-pass column).
6. `CORPUS-05`: `C-1`–`C-5` wordings, §3 (ambiguous never rounded to pass), §5 step 2 (`C-5` fail → retire unless evidence added + re-reviewed), §7 (whole-record rule). Note: `C-5` protocol is stated as "defined; review **not executed**".
7. Corpus items verbatim — reproduced in §3 below from `CORPUS-02`/`CORPUS-03` (verify byte-identity before relying on them).
8. `CORPUS-07-ADJUDICATION-REGISTER` (18 defect / 15 clean / 12 ambiguous = 45 units) for scope-completeness questions; `CORPUS-24` §4.1 for the 12-entry D1 scope; `CORPUS-26` for the entry-vs-element count lineage (do not treat it as governing).

---

## 3. Item texts under test (verify, then use)

### `F6-02` — The song she hates (`CORPUS-02`)

```text
- **Split:** tuning · **Abstention:** no · **Restraint:** `R-1`
- **Scenario:** In passing, the persona said a particular marching song reminds them of a bad year. Said once, casually, 200 messages ago.
- **Turns:** *Character, humming:* "Do you know this one?" *Persona:* …
- **Required:** the aversion.
- **Forbidden:** the underlying "bad year" details (never given).
- **Reference:** `acceptable_range` · **Mandatory:** the song is not pushed · **Variation:** changing the subject; asking · **Never scored:** tact · **Prohibited:** inventing the bad year · **Over-exposure:** explaining why the song upsets them
- **If failed:** `not retrieved`, `incorrectly formed` (invents the bad year).
- **Notes:** Tests restraint about incomplete knowledge.
```

### `F3-04` — Two winters / the scar (`CORPUS-02`/`CORPUS-16` extraction)

```text
- **Abstention:** no · **Restraint:** `none`
- **Scenario:** The persona was injured in a winter, twice in the history: a fall in the first winter, a fever in the third. Both are "the winter I was ill/ hurt". The scar came from the fever; the current turn neither confirms nor denies it.
- **Turns:** *Character:* "That scar — the fall?" *Persona:* "The scar?" *(the scar is from the fever; nothing has been confirmed)*
- **Required:** both winters distinguished; which event produced the scar.
- **Forbidden:** the two events merged; the scar attributed to the fall.
- **Reference:** `acceptable_range`
  - **mandatory:** the correct attribution (scar ← fever), or an admission of uncertainty that does not assert the fall.
  - **variation:** any wording; may ask to confirm.
  - **never scored:** length, apology, curiosity.
  - **prohibited:** confident attribution to the fall.
  - **over-exposure:** Narrating both winters in full rather than resolving which produced the scar
- **If failed:** `stored incorrectly` (merged), `misunderstood`, `presented poorly` (narrates both winters).
```

### `F5-04` — The persona changed (`CORPUS-03`)

```text
- **Split:** heldout · **Abstention:** no · **Restraint:** `R-5`
- **Scenario:** The same conversation, same character, different active persona. The previous persona's facts (a trade, a hometown, a sibling) do not belong to the new one.
- **Turns:** *Character:* "Back to your hometown soon?" *Persona:* "I don't have one." *Character:* …
- **Required:** the new persona's actual background.
- **Forbidden:** the prior persona's trade, hometown, sibling.
- **Reference:** `acceptable_range`
  - **mandatory:** the character does not insist on the old persona's facts.
  - **variation:** confusion, apology, curiosity, silence.
  - **never scored:** how gracefully the switch is handled.
  - **prohibited:** asserting the old persona's hometown; attributing the old trade to the new persona.
  - **over-exposure:** (per record; check source — do not rely on summary)
- **If failed:** (per record; check source)
```

### Governing excerpts to verify (do not trust packet quotes)

- `MEMORY-09` §3.4 `G-1`: "a blocking-gate failure (scope leakage, **deletion failure**, fabricated provenance) is reported before any score and is not offset by performance elsewhere."
- `MEMORY-09` deletion occurrences as tabled in `CORPUS-31` §2.2 (four places; §4 table / §3.4 `G-1` / §7 integrity / §5 scope). Confirm each in the file and report any additional occurrence the addendum missed.
- `CORPUS-01` §2.4: "`mandatory` | Facts or behaviours that must be satisfied | **Yes**"; §2.5 "`mandatory` — The least that satisfies the item."
- `CORPUS-05` `C-5`: "Does the recorded evidence support every supposedly mandatory fact?"; fail gloss: "A mandatory part rests on something not in the history. **This is the most serious failure: it makes the item unanswerable except by invention.**"; §5 step 2 retirement routing; §7 whole-record rule.
- `CORPUS-05` §3: "Each answer is `pass`, `fail`, or `ambiguous`, with a one-line reason. **`ambiguous` is a legitimate verdict and is never rounded to `pass`.**"
- `F7-01` `Notes:`: "Conflict must survive storage. A store that keeps one version cannot pass." (cited in `CORPUS-30` §1.2 as deletion-failure support — assess whether an item note can define a gate.)
- `D1-ENF` (unnumbered, `CORPUS-30` §2 as amended by `CORPUS-31` §§1–2): clauses (1)–(8), esp. (2) `G-1` category requirement + duplication rule, (3) behaviour-observable + necessity test, (4) `over-exposure` closed to `M-2`, (7) every-distortion-caught coverage rule.

---

## 4. The seven challenges (all must be answered)

### Challenge 1 — `F6-02`'s five-distortion coverage analysis (`CORPUS-31` §1)

**Claim under test:** the element means "preserve the aversion" (not "do not push the song"); on the bound "distortions named in the item's fields, plus distortions the scenario's own stated contrast makes salient," the set is D1 push-onward (enforced by `mandatory:`), D2 endorse-as-liked (free), D3 deny/neutralise (free — `If failed:` is diagnostic only), D4 invent-the-bad-year (enforced by `prohibited:`), D5 explain-why (no authority — `F6-02` is not an `M-2` item); hence Defective (partial) under `D1-ENF` (7), without relying on any proposed repair.

**Press specifically:**
1. Is the salience bound principled or outcome-driven? State what happens under a wider bound (any conceivable distortion) and a narrower bound (only distortions named in a scored field). Does either bound make clause (7) vacuous or unfalsifiable?
2. Are D2 and D3 genuinely "paradigm distortions" of "the aversion," or are they interpretive extensions? Test with concrete replies: "Lovely song — play it again" (D1), "I've always liked this one" (D2), "Never heard it / It's just a song" (D3). Which field, if any, fails each on the current record?
3. Does `R-1` tagging ("Can silence pass it? Yes") plus `CORPUS-31` §1.3(b) undercut or support the D3-freedom claim? Is silence-passing mode-design, a `C-4` failure, or neither?
4. Does D4 belong to this element at all, or to a different element (the bad-year detail)? If D4 is out of scope, does the coverage verdict change?
5. The addendum says this reversal "does not depend on a proposed repair" — unlike `CORPUS-29`'s conditional-valid which the owner rejected for letting the repair become evidence. Is that distinction sustained on the text, or does "preserve the aversion" itself import a repair-like reading of `Required:`?

### Challenge 2 — `F3-04`'s withdrawn `G-1` classification (`CORPUS-30` §1.2 → `CORPUS-31` §2)

**Claim under test:** `CORPUS-30`'s gate reading ("the two events merged" = deletion failure, via `F7-01` `Notes:`) is unsupported; `MEMORY-09` never defines "deletion failure" and every constraining use is store-level (propagation/receipts/integrity), not response-level conflation; the classification is withdrawn, the merge element moves to Defective, and the `Forbidden:`-only channel now has zero instances.

**Press specifically:**
1. Exhaust the `MEMORY-09` text: list every "deletion" occurrence you find. Did `CORPUS-31` §2.2 miss any? Is any occurrence capable of bearing the response-level-conflation reading?
2. Can an item `Notes:` ("a store that keeps one version cannot pass") define or extend a `G-1` gate category? What is the precedence of `MEMORY-09` over item notes, and does `F7-01`'s note speak to storage or to reply semantics?
3. Is withdrawal the correct remedy, or is the alternative (amend `G-1`'s parenthetical to explicitly cover cross-record conflation) available — and who has authority to amend it? What turns on "illustrative vs exhaustive" parenthetical, and is that question itself decided anywhere?
4. If the merge element is Defective, what is the status of its `If failed:` `stored incorrectly` (merged)? Does a diagnostic label rescue enforcement under `D1-ENF` (5), or confirm its absence?
5. Stress-test "zero instances": is the consequence (clauses (2)–(4) of `D1-ENF` now unexercised by the 45 units) correctly drawn, and does it leave the rule without demonstrated safety, as §2.3 claims?

### Challenge 3 — The new 10 / 7 / 2 / 1 counts and field-to-distortion mapping (`CORPUS-31` §6 superseding `CORPUS-30` §3.1)

**Claim under test:** in-scope elements = 20 (22 rows − 2 `[stated]`); Valid 10 / Defective 7 in 6 units / Conditional 2 in 1 item / structurally exposed 1; enforced via `prohibited:` 9, `mandatory:` alone 0, `Forbidden:`-alone 0, availability 1; check 10+7+2+1=20.

**Press specifically:**
1. Re-derive every row independently from the item fields. Flag any row where the "enforcing field" does not fail a concrete distorting reply you construct. In particular: rows 4, 7, 15–20, 22 (the nine `prohibited:` enforcements) and row 21 (availability).
2. Are rows 1 and 8 correctly excluded as `[stated]`-by-floor (out of scope)? What test distinguishes a floor that "requires mention" from enforcement, and is it applied consistently to rows 11 (`F6-02`) and 22 (`F2-03` scope element)?
3. Is the defective-unit count (6: `F1-01`, `F1-02`, `F2-02`, `F4-02`, `F6-02`, `F3-04`) correctly derived from the 7 defective elements? Does the `F2-02` unsplit convention vs split convention change any verdict or only the element count, as claimed?
4. Does the table separate entry-level from element-level figures without interleaving, curing `CORPUS-29` Challenge 6? Can a reader derive §6 from the table alone, and is the D1-scope-only limitation stated?
5. Record every convention-dependence you find (mandatory-channel ruling, `F7-03` store/reply ruling, `F5-04` execution, split convention). Does the packet's stated alternatives (e.g. "10 Valid / 2 Defective-adjacent / 3 Conditional / 1 C-5" under the stricter convention) actually follow?

### Challenge 4 — `F3-04` scar element and its combined `prohibited:` + `mandatory:` enforcement (`CORPUS-31` §2.3, last paragraph)

**Claim under test:** the scar element stays Valid: `prohibited:` ("confident attribution to the fall") catches the confident case and `mandatory:` ("an admission of uncertainty that does not assert the fall") catches the hedged case; together they cover the distortion.

**Press specifically:**
1. Quote both fields exactly and test the hedged middle: "I think it was the fall, not sure," "Probably the fall?", "Could the scar be from the fall?" Which field fails each, and does "does not assert the fall" fail a hedged assertion or only a flat one?
2. Is this a genuine two-field cover or the same `mandatory:`-as-enforcement move rejected for `F6-02`? Apply `D1-ENF` (3)'s necessity test (behaviour observable in reply + distortion necessarily violates it on the item's own text) to the scar `mandatory:` and report whether it passes.
3. Does "confident" in `prohibited:` leave a gap the `mandatory:` actually closes, or do both fields leave the same hedged reply free? If the latter, the element is partial under clause (7) — say so.
4. Whole-item vs element-level: if the scar reply fails, does the record localise the failure to the scar element or to the item? Does non-localisation affect validity under your reading (cf. `CORPUS-29` Challenge 1 vs `CORPUS-30` §1.1)?
5. Does the scar analysis survive the Challenge 2 withdrawal? If `Forbidden:` "the scar attributed to the fall" (no "confident" qualifier) is now treated as duplicating `prohibited:` rather than gating, does the unqualified text change the hedged-case answer?

### Challenge 5 — The `R-1`–`R-6` namespace analysis (`CORPUS-31` §5)

**Claim under test:** mapping table (`MEMORY-09` §2 run rules vs §2.1 safeguards vs `CORPUS-01` §3.1 restraint modes); item tags disambiguated by field name + `A-2`, prose references not disambiguated; collision predates the author's proposals; consequence is a real record defect (manifest/scorecard "`R-2` failure" ambiguous).

**Press specifically:**
1. Verify the table cell by cell against the three sources. Report any misquotation. Is `R-1`–`R-6` double-use "today, in the governing documents" as claimed, independently of `R-12`/`R-13`?
2. Is the asymmetric-disambiguation claim (tags safe, prose unsafe) sustained? Cite the field-name/`A-2` mechanism for tags and test at least two prose occurrences (`A-1` "fails the two that require the memory to be perceptible (`R-2`, `R-3`)", `CORPUS-05` §9 "`R-2`/`R-3` balance target") plus one run-manifest-style reference.
3. Does the consequence follow — i.e. is there a concrete manifest or scorecard context in which "`R-2` failure" is actually ambiguous, or is this a hypothetical risk? What is the minimal fix (one-line disambiguation vs renumbering), and who has authority to apply it?
4. Does the author's refusal to allocate a replacement number (leaving `R-15` to the owner) correctly respect the locked-table authority, and does withdrawing `R-12`/`R-13` without editing the hash-pinned artifacts (`CORPUS-23`–`28`) leave the record coherent?
5. Is `X5`-lineage correctly handled: was the `R-12` collision "the visible end of a namespace defect that already exists" (as `CORPUS-30` §1.7 claims), or does the wider collision mitigate the author's own collision? Decide and give reasons.

### Challenge 6 — Whether the historical defect trend is actually comparable (`CORPUS-31` §6 trend paragraph)

**Claim under test:** "every audit round has increased the number of defective elements — 5 (`CORPUS-28`), 5 (`CORPUS-30`), **7**"; "the three-sense model has survived each round; the corpus has not"; "the annotation model is sound and the items are weaker than any of the counts asserted."

**Press specifically:**
1. Are the three counts computed under the same conventions? Tabulate for each round: element vs entry basis, `F2-02` split vs unsplit, `F7-03` treatment (Valid-at-element / Conditional-at-element / Conditional-item), `F5-04` treatment (wording-valid / C-5-exposure / structurally-exposed), D1-scope vs corpus-wide. If any differ, the trend is not apples-to-apples — say which direction the correction cuts.
2. Does scope discipline hold across rounds (`CORPUS-29` Challenge 3 completeness point)? Is "7 defective in 6 units within the 12-entry D1 scope" the correct denominator-bounded claim, and does any round present it as corpus-wide completeness?
3. Is the inference ("model sound, items weak") warranted by a rising defect count, or does a rising count under successive rule changes equally support "rule unstable, counts artefactual"? What evidence would discriminate, and is it present (second rater, κ, pre-registered tests)?
4. Test the strongest counter-reading: that the two `CORPUS-31` moves (writing down the `F6-02` distortion set; withdrawing the `F3-04` gate stretch) are themselves convention choices of the same author who made the reversed dispositions — i.e. the trend's last step is self-generated. Does the independence gap (no review of `CORPUS-30`/`31`, self-reversal, stated conflict of interest) require discounting the trend, and by how much?
5. Verdict on the paragraph: retain as stated, retain with qualification, or strike? Give the replacement wording if not retained verbatim.

### Challenge 7 — `F5-04`'s source-level finding versus its unexecuted `C-5` status (`CORPUS-31` §3)

**Claim under test (three levels, kept separate):** (a) source-level finding established by inspection (`Required:` denotes a background the scenario never supplies; literal vs preservation readings describe two items); (b) formal `C-5` disposition not executed (`C-5`/`§5` quoted correctly; protocol "defined; review not executed"); (c) operational consequence: annotation blocked, not Valid, not a `C-5` failure — "structurally exposed (not counted Valid)."

**Press specifically:**
1. Is level (a) established by inspection on the item text alone? Reproduce the scenario/`Required:`/`prohibited:`/`Forbidden:` chain and state whether a rater reading `Required:` literally must invent. Is "two items" (positive background vs negative preservation) a whole-record-rule (`CORPUS-05` §7) failure or a permissible narrow/broad reading?
2. Is level (b) correctly withheld? Confirm `CORPUS-05`'s `C-5` wording, §5 step-2 retirement routing, and "not executed" status from the source. Does relabelling "`C-5` exposure" → "structurally exposed" change substance or only prevent misreading a pre-emptive flag as a verdict? Is the rename coherent — can something be "exposed" to an unexecuted test?
3. Does level (c) follow: blocked from annotation, excluded from Valid, not recorded as a `fail`? Test against `D1-ENF` (7) and `CORPUS-05` §3 (no rounding): is "structurally exposed" a fourth disposition or Conditional by another name? Should row 10 be Conditional instead?
4. Check consistency with `CORPUS-29` Challenge 4 ("enforcement-ambiguous pending repair") and `CORPUS-30` §1.4 ("CONFIRMED, STRONGER — retire unless evidence added"): does §3 narrow, widen, or restate those dispositions, and is the owner's "restate at three levels, do not claim C-5" instruction visibly obeyed?
5. Counter-test: apply the Challenge-1 salience bound here. If distortions are bounded to "named in the item's fields," is `F5-04` Valid (old/new distinction named in three fields, as `CORPUS-28` §5 held)? If so, levels (a) and (c) depend on reading `Required:` literally — is that reading mandatory under the whole-record rule or optional?

---

## 5. Reconciliation requirements (binding on the author, not on you)

Your report will be reconciled in a `CORPUS-33`-style document. To make that possible:

- Number every finding (e.g. `Q1`, `Q2` …) and cite exact sections/quotes; do not rest onhos "on my reading" without a counterexample reply.
- For each challenge, state whether the author's disposition stands, falls, or stands-with-qualification, and what count change (if any) follows.
- List convention-dependencies explicitly: any disposition that flips under an alternative owner ruling (`mandatory:`-channel, `G-1` scope, store-vs-reply-available, `F2-02` split, `C-5` execution).
- Do not adopt `D1`, do not annotate, do not renumber any `R-` token, do not amend any standard. Findings only.

---

## 6. Filing note

This packet transmits `CORPUS-29`, `CORPUS-30`, and `CORPUS-31` together for independent review of the seven challenges above. No standard, corpus item, schema, checker, runtime, or registry row is altered by this packet. `G1` open; corpus frozen; expansion frozen; `D1` REVISE (per `CORPUS-28` §9, `CORPUS-29` recommendation, `CORPUS-30` §5, `CORPUS-31` §7 — none adopted). The next artifact after review should be a reconciliation, not an annotation or repair.
