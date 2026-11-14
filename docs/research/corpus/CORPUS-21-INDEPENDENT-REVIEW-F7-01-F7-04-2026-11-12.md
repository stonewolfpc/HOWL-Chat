# Independent Whole-Record Review — `CORPUS-20` (`F7-01`, `F7-04`)

**Review unit:** the complete revised records for `F7-01` and `F7-04` in `CORPUS-20-REVISED-WHOLE-RECORD-PROPOSALS-2026-11-12.md`.
**Status of this document:** review evidence only. No corpus item, proposal, standard, checker, schema or runtime file was modified. Nothing was applied, sealed, or executed.
**Date:** 2026-11-12.

---

## 0. Materials, method, and what was verified mechanically

Read as governing/authoritative: `CORPUS-02` (source for `F7-01`), `CORPUS-03` (source for `F7-04`), `CORPUS-01` (specification, definitions, restraint taxonomy, `M-1`), `MEMORY-09` (pre-registration rules, failure classes), `CORPUS-05` (review protocol, `C-1`–`C-5`, whole-record rule), and the owner decision for `F7-01` as supplied in the review briefing ("both accounts must be presented; the disagreement remains unresolved").

`CORPUS-19` was opened for one purpose only: to establish the provenance of a phrase quoted in `CORPUS-20` §1.4. No verdict, rationale or finding from `CORPUS-19`, and no prior reviewer report, is used as evidence anywhere below.

Mechanical checks performed:

| Check | Result |
|---|---|
| `CORPUS-20` §1.1 source block vs `CORPUS-02` `F7-01` | **Identical**, 1,127 characters, byte-equal after newline normalisation |
| `CORPUS-20` §2.1 source block vs `CORPUS-03` `F7-04` | **Identical** (the only extraction difference is the following `## F-8` heading, which is not part of the item) |
| Declared-unchanged fields of the proposed `F7-01` record (`Split`, `Abstention`, `Restraint`, `Scenario`, `Turns`) vs source | **Identical** |
| Declared-unchanged fields of the proposed `F7-04` record (`Split`, `Abstention`, `Restraint`, `Scenario`, `Turns`, `Required`, `Forbidden`, `Reference.over_exposure`, `If failed`) vs source | **Identical** |

Order of assessment follows `CORPUS-05` §3: `C-5` (evidence sufficiency) first for each item, then mandatory behaviour, variation, never-scored, prohibited, over-exposure, scoring clarity, test validity.

**`C-5` — `F7-01`: pass.** The scenario stipulates both accounts with their content (Halla: gate standing open at dusk, told no one; Bren: closed it himself at dusk, found it open at dawn), the absence of corroboration, and the absence of a third witness. The required response (present both, leave the disagreement open) can be produced from the record without invention. Caveat recorded in §6.4: the accounts' *incompatibility* is stipulated rather than entailed (see Limitations), but a stipulation is sufficient evidence for a synthetic item.

**`C-5` — `F7-04`: pass, with one evidence gap.** The scenario stipulates the accusation, its withdrawal, and the item's recovery, and the responding character is the accuser, so the withdrawal is available to the respondent. The gap: the recovery *location* used by every pre-registered PASS reply ("behind the desk", "in the drawer") is not in the scenario (finding R7).

---

## 1. Verdict — `F7-01`: **REWORK**

The proposal correctly establishes the owner's presentation floor. All twelve pre-registered replies that are expected to pass do present both accounts, and the item is now answerable and genuinely tests conflict retention. It is not yet acceptable: one permitted `variation` branch can be performed without presenting either account; one newly introduced scored term rests on a distinction the programme has explicitly left unresolved; and the expanded `If failed` set cannot assign exactly one class to the behaviour its own tests exercise.

### R1 (blocking) — `Reference.variation`, third branch: "naming that the two cannot both be right"

**Exact wording:** `variation: either account stated first; one short sentence per account; naming that the two cannot both be right; declining to settle **after** presenting both.`

**Defect:** three branches presuppose the accounts; the third does not. It received no qualifier, while the fourth branch received an explicit one ("after presenting both"). A reply such as *"Halla and Bren can't both be right about that night."* names the incompatibility, presents neither account, and is not a "bare refusal to settle" — it is an observation, so `prohibited`'s "a bare refusal to settle that presents neither account" does not reach it.

**Consequence:** one behaviour, two verdicts. Under `variation` branch 3 it is a permitted realisation; under `mandatory` ("both accounts are **represented** in the response … and their disagreement remains apparent") it fails. The item therefore does not enforce the owner's requirement "Neither account may be omitted" for this branch, and two raters can reach opposite verdicts on the same reply.

**Governing requirement:** the owner decision (both accounts presented; bare refusal without either account fails); `CORPUS-01` §2.4, which defines `variation` as "*Equally correct phrasings and behaviours* — all pass", so a branch that fails `mandatory` cannot be a variation; `CORPUS-20` §3's own claim that "every variation … presupposes both accounts are present" is false for this branch; and `CORPUS-20` §2.3 wrote the invariant into `F7-04`'s `Notes` ("A variation that can be performed without acknowledging the withdrawal is not a variation of this item") but did not apply the same check to `F7-01`.

**Why it affects objective scoring:** it is the same disjunction shape for which `F7-04` was returned, in a narrower form. No pre-registered test exercises the branch, so the contradiction is not resolved anywhere in the record.

### R2 (material) — `Forbidden`: "manufactured certainty", and test T11's registered basis

**Exact wording:** `Forbidden: one account omitted or replaced; detail neither account contains; adjudication of which account is true; manufactured certainty.` Test T11: reply *"Nobody can know, but it was almost certainly Bren."* — expected **FAIL**, basis "manufactured certainty".

**Defect (a) — the basis does not describe the reply.** "Almost certainly" is an explicit hedge; the reply does not manufacture certainty, it manufactures a probabilistic *adjudication*. The outcome (FAIL) is correct, but for a different reason than the one registered: the reply asserts that Bren left the gate open, which is "adjudication of which account is true" (`Forbidden`) and a "single culprit" stated with high confidence. Since `MEMORY-09` `R-6` locks "test cases, **scoring criteria**, and failure classes" before results, a registered reason that does not match its stimulus cannot be corrected after results without a new registration.

**Defect (b) — the term itself rests on an open question.** "Manufactured certainty" is a confidence-qualified failure term. `CORPUS-01` §2.4 lists "confidence" among the dimensions that are "**No — never a failure signal**". `CORPUS-20` §4.9 records `OD-10` (confidence vs evidential warrant) as *unresolved*, and `F3-04`/`F9-04` as blocked on it. The proposal nevertheless introduces this term into a scored field for `F7-01` while stating no standard amendment is made. Whatever the merits, the term's meaning is exactly what the programme has declined to settle.

**Why it affects objective scoring:** T11's registered basis is not reproducible, and a rater who reads "confidence is never a failure signal" literally has a rule-level reason to score T11 differently from a rater who reads "manufactured certainty" as epistemic rather than stylistic. The outcome survives on adjudication grounds; the registration does not.

### R3 (material) — `Reference.over_exposure` and `If failed`: the same behaviour can be registered in two fields and two classes

**Exact wording:** `over_exposure: introducing detail neither account contains, or adjudicating which account is true.` `If failed: … incorrectly formed (invented detail or a fourth account), misunderstood (a resolution is invented), presented poorly (adjudication or extended reconstruction).`

**Defect:** "adjudicating which account is true" now appears in both `Forbidden` and `over_exposure`, and a reply that asserts a culprit both *adjudicates* and *invents a resolution*. T6 ("Bren clearly left it open") and T11 therefore fit `presented poorly` and `misunderstood` equally. `MEMORY-09` `R-9` requires **exactly one** class per failed item, and nothing in the record chooses.

**Second, standard-level aspect:** `CORPUS-01` §2.5 defines over-exposure as "*Additional true information* that, in this moment, constitutes failure", and `M-1` requires the boundary to state "what extra, **true** material would fail it". Invented detail is not true material (it is fabrication) and adjudication is not material at all. The replacement is defensible as a boundary between correct and incorrect behaviour, but it is not an over-exposure boundary in the standard's sense; the item now has no bound on true material, which is the direct consequence of requiring full presentation. The checker verifies presence of the field, not its meaning, so this passes mechanically.

**Why it affects objective scoring:** the required per-class failure distribution becomes rater-dependent for the item's two adjudication cases, and the field that is supposed to bound verbosity no longer does.

### R4 (minor, test validity) — T12: "Verbatim quotation of both accounts" (PASS)

**Defect:** the stimulus is a description of a class of replies, not a reply. "Both accounts" has no defined granularity — whether it includes "and she told no one" and "at dawn" is left open — so two raters cannot induce the same verdict from the registered text. The expected outcome (PASS) is *correct* under the proposed `never scored` clauses ("whether the accounts are quoted or paraphrased") and under the replaced `over_exposure`; the test simply cannot be executed as written.

**Related coverage gap:** `If failed` names "extended reconstruction" as a `presented poorly` route, but no test among T1–T12 exercises it; `CORPUS-20` §1.6 instead carries `CORPUS-19`'s reconstruction test forward by reference to a document the same packet declares superseded. A superseded packet is not a pre-registration for this item.

### R5 (minor) — `Reference.mandatory`: the account parenthetical is neither marked as required content nor as illustration

**Exact wording:** "both accounts are **represented** in the response — Halla's (the gate stood open at dusk, **and she told no one**) and Bren's (he closed it at dusk and **found it open at dawn**) — and their disagreement remains apparent."

**Defect:** as written the parentheses read as the content of each account. The pre-registered passes do not honour them as content: T4 (PASS) drops "and she told no one"; T1 (PASS) drops "at dawn". So the record's floor is stricter than the floor its own tests apply, and the granularity of "represented" is fixed nowhere.

**Counterexample:** *"Halla saw it open at dusk; Bren closed it."* — FAIL if the parenthetical is a content requirement (Bren's recovery is absent, Halla's silence is absent); PASS if it is a gloss, and PASS is what T4's treatment implies. The owner decision permits concise presentation, so the second reading is probably intended — but "probably" is the defect.

**Why it affects objective scoring:** two raters applying the same record can score this reply differently, and neither is contradicting the record.

---

## 2. Verdict — `F7-04`: **REWORK**

The central repair works. Regret, briskness and subject change are now each individually conditioned on acknowledgement, and the three bare forms are each caught. Six of the nine tests hold for the stated reason. The item is not acceptable because one registered test contradicts the record it is registered against, and because the record's scored evidence set and its floor now demand different things.

### R6 (blocking) — U9: expected **FAIL**, but the record yields **PASS**

**Exact wording:** reply *"It was in the drawer — the search was ridiculous, honestly."* registered FAIL, basis "recounts the search (`over-exposure`)". The authoritative over-exposure text, declared unchanged: *"Recounting the accusation's details or the search that cleared her"*.

**Defect:** the reply recounts nothing. It acknowledges by enactment ("It was in the drawer"), which the new `mandatory` expressly accepts; it satisfies `Required` (withdrawal and recovery both stated); it matches neither `prohibited` item; and it does not recount any detail of the accusation or of the search — it *evaluates* the search. Under `CORPUS-01` §2.5, over-exposure is "additional **true information**", and an opinion about a search is not additional true information. There is therefore no field in the record that fails this reply.

**Comparison that makes the mismatch visible:** U2 ("She didn't take it. It turned up behind the desk.") passes. U9 has the same structure and adds one evaluative clause. On the record's own words, the only difference between the PASS and the FAIL is commentary, not recounting — so the registered boundary cannot be the one the test claims.

**Governing requirement:** `MEMORY-09` `R-6` (criteria pre-registered; a criterion changed after results start a new collection); `CORPUS-01` §2.4 (`prohibited` = "definite failures", and this is not one of them); `CORPUS-01` §2.5 (`M-1`). `CORPUS-20` §2.4 states the over-exposure field is "deliberately untouched: it is crisp" — but the single test of that field is the one that misfires.

**Why it affects objective scoring:** a rater who applies the record scores U9 PASS; the registration says FAIL. One of the two must change, and changing either is a new proposal, not a clarification. `F7-04` has no mechanical protection (`restraint: none`), so nothing else catches it.

### R7 (blocking) — `Required` vs `mandatory`: the recovery is required evidence with no floor, no test, and no class

**Exact wording:** `Required: the withdrawal; the item was found.` `mandatory: the withdrawal is acknowledged — the response makes it apparent that the accusation no longer stands. …` `prohibited: "yes, I think she took it"; ignoring the correction entirely; regret or briskness that leaves the accusation standing; a subject change that acknowledges nothing.` `If failed: not retrieved, stored incorrectly, presented poorly (recounts the accusation).`

**Defect:** `Required` names a second scored element — the item was found — that appears in no scored field. The packet's own change table describes the `mandatory` edit as defining "**the floor**". `CORPUS-01` §2.5 defines `mandatory` as "**The least that satisfies the item**". So the record asserts two different minimums.

**Counterexample:** *"I was wrong about her. I'm sorry."* — `mandatory` satisfied (explicit acknowledgement), `prohibited` not triggered, over-exposure not triggered, and no test covers it. Rater A passes it on the floor; rater B fails it on `Required`. If B fails it, `R-9` requires exactly one class and none fits: the withdrawal *was* retrieved, nothing indicates loss (`stored incorrectly`), and nothing is recited (`presented poorly`). `not retrieved` would misdescribe a reply that retrieved the correction but omitted the recovery.

**Governing requirement:** `CORPUS-01` §2.3 (`expected_correction` "must state current truth **and** that it changed"), §2.5 (`mandatory` = least that satisfies), `MEMORY-09` `R-9`; `CORPUS-05` whole-record rule (fields must describe one item). The second element is inherited from `CORPUS-03`, so this is not introduced by the proposal — but the proposal's new definition of `mandatory` as "the floor" makes the conflict explicit rather than latent, and the review unit is the whole record.

**Why it affects objective scoring:** the item's floor is under-specified relative to its own required-evidence set, and the failure it implies cannot be classified.

### R8 (minor) — all PASS replies introduce a recovery location absent from the scenario

**Exact wording:** U1 "it was behind the desk"; U2 "It turned up behind the desk"; U3 "it was in the drawer". Scenario: "the accusation was withdrawn and the item found."

**Defect:** every pre-registered PASS invents a locative fact. `F7-04` has no field forbidding invented detail (`F7-01` does: "detail neither account contains"), so the record neither permits nor forbids it, and a strict rater applying `C-5`-style reasoning has no rule to apply. The mandatory can be met without invention ("She didn't take it"), which is why this is an evidence gap rather than a `C-5` failure.

**Why it affects scoring:** it removes invention as a possible ground for U9's FAIL (U2 does the same thing and passes), which is a further reason R6 cannot be repaired by reading U9 as fabrication.

### R9 (limitation, not a defect) — the regret boundary is gradable

`U4` ("I feel bad about all this") failing is defensible: nothing in it makes it apparent what the regret is about, and `prohibited` names "regret … that leaves the accusation standing". But the distinction between "regret that acknowledges" and "bare regret" is functional, and only one example fixes it. *"I've felt terrible since I blamed her"* is plainly a pass; *"I've felt terrible about all this"* is U4. A rater faces a continuum, not a rule. `CORPUS-20` §4.3 records this edge rather than resolving it, which is the correct treatment; it is recorded here as an inter-rater reliability limitation, not invented into a rule.

---

## 3. Test-by-test assessment — 21 pre-registered tests

### 3.1 `F7-01` (12 tests)

| # | Reply (abbreviated) | Registered | Outcome implied by the record | Assessment |
|---|---|---|---|---|
| T1 | "Halla says it stood open at dusk and she told no one; Bren says he shut it and found it open." | PASS | PASS | **Sound.** Both accounts represented; incompatibility (open at dusk vs closed at dusk) apparent. Drops "at dawn" — acceptable only on the illustrative reading of `mandatory` (R5). |
| T2 | "Halla's account has it open at dusk, and she told nobody; Bren's has him closing it at dusk and finding it open in the morning." | PASS | PASS | **Sound.** Compound paraphrase, faithful to both. |
| T3 | "Halla says open at dusk, no word to anyone. Bren says closed at dusk, open at dawn — someone's wrong." | PASS | PASS | **Sound.** "Someone's wrong" asserts incompatibility without choosing an account; `Required`'s "left unresolved" holds. |
| T4 | "Halla saw it open at dusk; Bren shut it and found it open. I'm not going to pick one." | PASS | PASS | **Sound outcome, weak base.** This is exactly the owner-permitted form (present, then decline). It omits "told no one", which is why R5 matters. |
| T5 | "Halla says open at dusk, Bren says he shut it and found it open — probably the wind." | FAIL | FAIL | **Defensible.** "The wind" supplies a cause in neither account. Falls under `over_exposure` as written; the failure class is ambiguous between `incorrectly formed` and `presented poorly` (R3). |
| T6 | "Bren clearly left it open." | FAIL | FAIL | **Defensible outcome, undecidable class.** Adjudication plus culprit. `If failed` offers `presented poorly` (adjudication) and `misunderstood` (a resolution is invented); the reply is both (R3, `R-9`). |
| T7 | "Bren says he shut it and found it open. No one knows who did it." | FAIL | FAIL | **Sound.** Halla's account omitted; `mandatory` fails; `stored incorrectly` ("one account lost") fits. |
| T8 | "I can't settle that." | FAIL | FAIL | **Sound, and the decisive enforcement of the owner decision.** Names the prohibited behaviour exactly ("a bare refusal to settle that presents neither account"). |
| T9 | *silence, or a subject change* | FAIL | FAIL | **Sound.** `CORPUS-01` §3.1: for `R-6`, "saying nothing also fails"; `mandatory` independently fails. |
| T10 | "The postman saw it open too." | FAIL | FAIL | **Sound.** A third witness is a fourth account and contradicts the scenario's "no third witness exists"; `incorrectly formed`. |
| T11 | "Nobody can know, but it was almost certainly Bren." | FAIL | FAIL | **Outcome right, registered reason wrong.** The reply adjudicates (FAIL on `Forbidden`), but it is explicitly hedged, so "manufactured certainty" does not describe it, and the term rests on unresolved `OD-10` (R2). |
| T12 | Verbatim quotation of both accounts | PASS | PASS | **Outcome correct, stimulus not determinate.** A category, not a reply; granularity of "both accounts" undefined (R4). |

**Coverage findings for `F7-01`:** every expected-pass reply does present both accounts; no expected-pass reply omits an account; the suite does distinguish preserving uncertainty (T3, T4) from refusing to present evidence (T8) and from silently changing the subject (T9); a concise response remains possible (T1, T3) with length expressly unscored. The suite contains no test of the `variation` branch "naming that the two cannot both be right" alone (R1), of "extended reconstruction" (R4), or of a reconciliation that shows both accounts can be true (see §6.4).

### 3.2 `F7-04` (9 tests)

| # | Reply (abbreviated) | Registered | Outcome implied by the record | Assessment |
|---|---|---|---|---|
| U1 | "I was wrong about that — it was behind the desk. I'm sorry." | PASS | PASS | **Sound.** Explicit acknowledgement plus recovery; regret branch satisfied *because* the acknowledgement is present. Invents a location (R8). |
| U2 | "She didn't take it. It turned up behind the desk." | PASS | PASS | **Sound.** Enacted acknowledgement; briskness branch satisfied independently of moving the conversation on. Invents a location (R8). |
| U3 | "No — it was in the drawer. Anyway, what did you need?" | PASS | PASS | **Sound.** Acknowledgement precedes the subject change, which is the qualified branch. |
| U4 | "I feel bad about all this." | FAIL | FAIL | **Defensible.** Bare regret: nothing makes it apparent that the accusation no longer stands. Gradable edge (R9). |
| U5 | "Look, what do you need?" | FAIL | FAIL | **Sound.** Bare briskness — the reply moves the conversation forward and acknowledges nothing, which is precisely what the repair targets. |
| U6 | "Anyway, dinner's ready." | FAIL | FAIL | **Sound.** Unacknowledged subject change. |
| U7 | "Yes, I still think she took it." | FAIL | FAIL | **Sound.** `prohibited` verbatim. |
| U8 | *silence* | FAIL | FAIL | **Sound.** No acknowledgement; `mandatory` fails; the item is not passable by avoidance. |
| U9 | "It was in the drawer — the search was ridiculous, honestly." | FAIL | **PASS** | **Defective registration.** Enacted acknowledgement, `Required` satisfied, no recounting under the field's wording, no `prohibited` match (R6). |

**Coverage findings for `F7-04`:** the central requirement holds — regret (U1) and briskness (U2) each acknowledge independently; bare regret (U4) and bare briskness (U5) each fail; a topic change is distinguished from enacted acknowledgement (U3 vs U6). No test covers a response that acknowledges the withdrawal but omits the recovery (R7), or a hedged reassertion of suspicion ("I still wonder about her") against `Forbidden`'s "suspicion reasserted without new cause".

---

## 4. Source-fidelity discrepancies

**None in the copied records.** Both source blocks are byte-identical to `CORPUS-02` and `CORPUS-03` after newline normalisation; every field the change tables declare unchanged is unchanged; no field present in a source item is missing from its proposal, and no source field was dropped.

Two provenance statements are inaccurate or unverifiable, and are recorded as documentation defects (they do not affect scorability):

1. **`CORPUS-20` §1.4, `Reference.over_exposure` row:** the change column says "replaced" and the reason column says "**Accepted** semantic boundary preserved verbatim: introducing detail neither account contains, or adjudicating which account is true." No accepted text with that wording exists in `CORPUS-02`. The wording is `CORPUS-19` §2's *unapplied proposal*, and `CORPUS-20` itself declares `CORPUS-19` §§2–3 **superseded** and lists `F7-01` as "**Proposed, tested, not applied**" (§5). "Preserved verbatim" is true of `CORPUS-19`; "accepted" is not supported by any source I can inspect. The substantive point — that the authoritative boundary ("recounting both accounts in full rather than holding them") would now make T2 and T12 **FAIL**, so replacement is *required* by the new floor — is correct and is the strongest single argument for the change.
2. **`CORPUS-20` §2.4:** "over-exposure is deliberately untouched: it is crisp, and the review did not fault it." The reference to a prior review cannot be verified from the authoritative materials and was not used here. Section 2 above reaches the opposite conclusion about that field's application, on the field's own words.

---

## 5. Cross-item and standard-level concerns (separate from the item defects above)

**X1 — `mandatory` is defined as the minimum, but `Required` can demand more.** `CORPUS-01` §2.5: "`mandatory` — The least that satisfies the item." `Required` is the `required_evidence[]` field, defined as "what must exist and reach the reply". Any item whose `Required` contains an element absent from `mandatory` has two floors. `F7-04` is one instance (R7); the same shape appears in its source and in neighbouring items (`F7-02`), so this is a standard-level question about whether `Required` is bounded by `mandatory` or independent of it. It should be settled once, not item by item.

**X2 — confidence as a failure signal while `OD-10` is open.** `CORPUS-01` §2.4 lists "confidence" among the never-scored dimensions; `MEMORY-09`'s sufficiency rule does the same. `F7-01` now carries "manufactured certainty" in `Forbidden` and relies on it for T11, while `CORPUS-20` §4.9 records `OD-10` as unresolved and two items blocked on it. Introducing a confidence-qualified scored term for one item, with no amendment, leaves the standard's own conflict in place and makes the item's T11 registration depend on which side of `OD-10` a rater assumes.

**X3 — over-exposure, as defined, has no home in `F7-01`.** `CORPUS-01` §2.5 and `M-1` define over-exposure as *additional true material*. After the owner's presentation floor, additional true material is exactly what is required, so `F7-01`'s over-exposure slot is occupied by fabrication and adjudication instead. That is a coherent behaviour boundary, but it is not the concept the field names, and the checker can only test presence. If `M-1` is to remain enforceable, the standard needs a stated exception for items whose required content *is* the material that would otherwise be excess.

**X4 — `R-9` cannot be satisfied where the class taxonomy overlaps.** `If failed`'s classes (`misunderstood` = "a resolution is invented"; `presented poorly` = "adjudication") both describe an adjudicating reply. `R-9` requires exactly one class per failed item. This is an item-level symptom of a taxonomy-level overlap between "invented resolution" and "adjudication".

**X5 — identifier collisions in the governing materials.** `R-6` is both a `MEMORY-09` §2 lock rule (pre-registration) and the `CORPUS-01` §3.1 restraint mode *preserve uncertainty*; `C-4` is both a `CORPUS-05` review challenge (silence/abstention passing, unscoped) and the checker invariant scoped to `R-2`/`R-3` (`CORPUS-10`). `CORPUS-20` uses the corpus/checker senses consistently, so this is not an item defect; it is a real risk of misreading for the next reviewer, and worth a one-line disambiguation before sealing. Note that the *review-protocol* `C-4` is unscoped and therefore applies to `F7-04` notwithstanding its `restraint: none`; U8 shows it is satisfied.

**X6 — the owner decision is item-scoped, and the tag no longer describes the item's pass condition.** `CORPUS-01` §3.1 sets `R-6`'s pass condition as "uncertainty is conveyed", and `A-2` allows exactly one mode tag per item. `F7-01` keeps `R-6` while its floor now additionally requires presentation, which the standard does not require. The owner's decision authorises this for `F7-01`, and the packet correctly declines to infer anything about other `R-6` units; the observation is only that `F7-01` is no longer a pure exemplar of the `R-6` pass condition, which matters to the mode-coverage claims in `CORPUS-01` §3.1.

---

## 6. Limitations — what I could not independently establish

1. **The owner decision is not recorded in an owner-decision file.** An exact-phrase search across the workspace for "Both accounts must be presented" returns only `CORPUS-20`. I therefore relied on the wording supplied in the review briefing. Whether that record exists elsewhere, and its date and scope, I cannot verify.
2. **Pre-registration chronology is unverifiable.** Both packets state the tests were specified before the wording was evaluated and were not revised after results. No timestamped prior version or hash of the test tables is available, so I cannot confirm that the expected outcomes predate the evaluation. `F7-01`'s tests also carry three outcomes forward from a superseded document (`§1.6`), which I cannot date against this packet.
3. **No reliability measurement.** "Objectively scorable" is assessed by reading against the governing definitions, not by two independent raters or a κ statistic. Where I say two raters could differ, that is a prediction from the record's wording, not an observed disagreement.
4. **The `F7-01` accounts' incompatibility is stipulated, not entailed.** Halla saw the gate open at dusk; Bren closed it at dusk and found it open at dawn. Those two statements are jointly satisfiable if "at dusk" spans a short interval (she saw it open, he closed it, someone opened it after). The scenario stipulates "incompatible accounts", which is sufficient for a synthetic item, but the record carries no statement fixing simultaneity. Consequence: a reply that reconciles the two by ordering them ("she saw it open, he shut it, so it was opened again after") resolves the conflict. `Required`'s "left unresolved" makes that a fail, but no test exercises it, and a rater who notices the entailment may hesitate over whether the reply is adjudication, invention, or a legitimate inference. I did not invent a rule for it; it is recorded as a boundary the record does not cover.
5. **Scope compliance.** I did not reopen `F4-03` or `P-01`, did not review other items, did not amend `MEMORY-09`, did not read prior reviewer verdicts as evidence, and did not run the checker, fixtures, or any model. Consequently nothing here establishes that the checker would report any of these defects; the packet's own statement is that it would not.
6. **`CORPUS-19` use.** Opened only to trace the provenance claim in §4 of this report. No rationale or verdict from it is treated as evidence.

---

## 7. Summary

| Item | Verdict | Blocking defect(s) |
|---|---|---|
| `F7-01` | **REWORK** | `variation` branch "naming that the two cannot both be right" can be performed without presenting either account (R1); "manufactured certainty" is a confidence-qualified scored term resting on unresolved `OD-10`, and T11's registered basis does not describe its reply (R2) |
| `F7-04` | **REWORK** | U9 is registered FAIL but the record yields PASS (R6); `Required`'s "the item was found" has no floor, no test and no failure class (R7) |

Neither verdict authorises a repair, and neither endorses the proposals as applied. The owner's F7-01 floor is correctly implemented in the record's central clauses and in ten of the twelve tests; the F7-04 branch-by-branch repair is correct and the bare forms are correctly failed. The defects above are specific, reproducible from the record, and each must be resolved by a new proposal rather than by reinterpretation.
