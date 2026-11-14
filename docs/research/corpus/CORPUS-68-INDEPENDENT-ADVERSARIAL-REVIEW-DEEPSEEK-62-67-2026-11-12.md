# CORPUS-68 — Independent Adversarial Review of DeepSeek's CORPUS-62–67 Governance Work

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-68` · **Version:** 1 · **Date:** 2026-11-12.
**Reviewer identity (externally supplied by the owner):** **Nemotron Lightning.** Recorded as owner metadata. **Not** inferred from, and not substituted by, any model self-description; **not verifiable** from the workspace. No model was asked to identify itself for this review.
**Executing agent:** Codebuff coding agent. Any framework persona label is an execution fact, not model identity.
**Execution:** fresh conversation. No prior substantive turns in this conversation.
**Review type:** adversarial follow-up review of `CORPUS-62` through `CORPUS-67` plus verification of the underlying sources. Read-only except for filing this report and its registry row.
**Status:** review record only. **NOT operative.** No corpus item, label, fixture, standard, schema, runtime, model, or checker was changed, repaired, adopted, or scored by this review. `G1` open · corpus frozen · expansion frozen · annotation blocked · convergence 0 of 2.
**Authority:** the owner work order reproduced in the current request (§§1–14), which authorizes read-only investigation, independent analysis, and filing of this report. G1 remains open; no sealing, no G2, no annotation, no model inference experiment was authorized or run.

---

## 0. Reader's guide to the verdicts

Per-artifact verdicts use the required scale: **ACCEPT / ACCEPT WITH QUALIFICATIONS / REWORK / REJECT / UNVERIFIABLE**. Qualifications are load-bearing, not decorative: where adoption of a proposal is blocked, the verdict says so. §4 additionally gives the three-option technical verdict required for `CORPUS-64` (technically sound / sound with qualifications / requires rework), and §6 the same for `CORPUS-66`.

Overall: DeepSeek's six artifacts are careful, honest about their limits, and verified correct on every central factual claim I could check. They nevertheless carry one high-severity provenance limitation (inherited, not authored), one high-severity diagnosticity gap in the `F3-01` proposal (partly acknowledged, partly understated), and a set of medium/low defects and omissions detailed below. **No artifact is rejected. Every artifact is accepted only with qualifications. Neither repair proposal is ready to adopt.**

---

## 1. Provenance, identity, and conflict disclosure

| Layer | What is known | Status |
|---|---|---|
| Reviewer identity for this report | Owner-supplied: **Nemotron Lightning** | Owner metadata; **not verifiable** from the workspace. No self-description was solicited or relied upon. |
| Executing agent | Codebuff coding agent, fresh conversation | Execution fact; framework persona label ≠ model identity. |
| `CORPUS-62`–`67` authorship | In-file declaration: **DeepSeek** (owner-supplied per those files) | In-file only; **not independently verifiable** from the workspace. No routing or cryptographic evidence exists. |
| `CORPUS-58` review authorship | Owner-supplied: Nemotron Lightning | Owner metadata; unverifiable here. The substantive review text is **not in this workspace and was not furnished in this execution** (same limitation as `CORPUS-61` §3.2). |
| `CORPUS-59`/`60` authorship | In-file: DeepSeek | Unverifiable, as recorded by `CORPUS-61` §2. |
| `CORPUS-61` reviewer | Owner-supplied: Solar Pro 4 | Unverifiable, as recorded there. |
| Owner decisions behind `CORPUS-62` | Reported in `CORPUS-62` §3 as owner-approved 2026-11-12 | **UNVERIFIABLE from the workspace** — see finding S-1. No filed owner decision record was found on disk. |
| Operational independence of this review from the `CORPUS-59`–`62` executions | Fresh conversation, read-only method | Asserted by execution context, **not cryptographically established**. |

**Conflict disclosure (required by the work order §10).** The work order states that I authored the original `CORPUS-58` review. That makes this review an adversarial **follow-up** by an interested party, not an independent second reviewer. Agreement below with any prior artifact is same-evidence agreement on workspace-read sources, never independent convergence. Convergence stands at **0 of 2** under `CORPUS-05` §7, unchanged. Nothing here confers convergence credit on any artifact, including this one.

---

## 2. Scope, authorization check, and limitations

**In scope:** reading `CORPUS-62`–`67`, `MEMORY-09`, `CORPUS-01`/`02`/`03`/`07`, `CORPUS-56`/`57`/`59`/`60`/`61`, the document registry, `check_corpus.py`, and the root-level `CORPUS-58` brief; recomputing hashes/counts; resolving registry links; running the existing checker; filing this report with a registry row.

**Out of scope and not done:** editing any corpus item, standard, tag, split, acceptance condition, or schema; modifying the checker; applying any proposal; running model inference; sealing `G1`; beginning `G2`; authorizing annotation; reconstructing the substantive `CORPUS-58` review.

**Authorization finding.** The work order for this review (§12) forbids modifying the registry, yet §13 requires registering this report. I resolve this exactly as `CORPUS-61` §14 did: the only registry change is this report's own row plus the header row-count delta (§9); no defect repair, no re-derivation, no other row touched. If the owner reads §12 as forbidding even that, the registry edit should be reverted and this file left unregistered — the review's content is unaffected either way.

**Limitations.**

1. The seven owner decisions are taken **as reported** in `CORPUS-62`. Their existence and wording cannot be verified from disk (S-1).
2. Every claim of the form "`CORPUS-58` found X" remains **UNVERIFIABLE** here — the substantive text is unfurnished in this execution, identical to `CORPUS-61` §3.2.
3. No inference experiment was run (not authorized). All diagnosticity verdicts are analytic pass-set results on text, as in every prior artifact.
4. The registry "full sweep" (all 402 rows resolve; hashes/bytes/lines match) is **reported by `CORPUS-62` §9.3–9.4** and spot-verified here (§8: 12 files rehashed, all 402 links resolved, no duplicate targets/numbers). A complete 402-row hash recomputation was not performed in this review; the sweep claim is therefore corroborated but not independently replicated in full.
5. Corpus files under `docs/research/corpus/` and the root brief are **untracked in git** (`??`); `docs/DOCUMENT-REGISTRY.md` is **modified in the working tree** (`M`). Pre-correction registry text is therefore recoverable via the working-tree diff and via the verbatim quote in `CORPUS-63` §3, not via committed history.

---

## 3. Source inventory (all read from disk in this session)

| Material | Hash / size verified here | Sections relied upon |
|---|---|---|
| `CORPUS-62` owner decision record | `35ab7b9f…` / 35,239 B / 400 lines — matches registry row 397 | §§1–12, esp. §§3, 5, 8, 9 |
| `CORPUS-63` registry correction | `8f620f46…` / 14,632 B / 198 lines — matches row 398 | §§1–10, esp. §§3–5 |
| `CORPUS-64` F3-01 proposal | `f16a97dd…` / 22,781 B / 257 lines — matches row 399 | §§1–14, esp. §§6–10 |
| `CORPUS-65` R-2 amendment proposal | `a914cb73…` / 12,436 B / 200 lines — matches row 400 | §§1–10, esp. §§2–5, 8 |
| `CORPUS-66` F6-01 proposal | `398f3e7d…` / 16,555 B / 208 lines — matches row 401 | §§1–9, esp. §§4, 5, 7 |
| `CORPUS-67` validator spec | `3c842eb3…` / 14,726 B / 208 lines — matches row 402 | §§1–9, esp. §§3–5, 7 |
| Root `CORPUS-58` brief | `292befdd…` / 4,311 B / 80 lines — matches row 393 | full text: work order, not a review |
| `CORPUS-01` specification | `be1f11b4…` — matches `CORPUS-62` §9.6 | §3.1 mode rows, coverage line, §7, §2.4–2.5 |
| `CORPUS-02` batch 1 | `75aa9be7…` — matches `CORPUS-62` §9.6 (verified against the file's own §9.6 string, not my earlier transcription) | `F1-01`, `F1-02`, `F3-01`, `F6-01` verbatim |
| `CORPUS-03` batch 2 | `d0dae1c0…` — matches `CORPUS-62` §9.6 | `F1-03` verbatim |
| `CORPUS-07` adjudication register | `72f6e2a2…` — matches `CORPUS-62` §9.6 | §2 rows `F3-01`/`F6-01`, §§4, 10, header |
| `MEMORY-09` preregistration rev 2 | `f77e47a9…` — matches `CORPUS-62` §9.6 | §3.5 (quoted §3.1 below) |
| `check_corpus.py` | `be3e317c…` — matches `CORPUS-62` §9.6 | full source (no registry logic — verified) |
| `CORPUS-56` audit | read in full | §§3–6 (esp. §5 "0 of 4", §3.1–3.4 verdicts) |
| `CORPUS-57` adjudication | read in full | §§2–5 (esp. §3.1 folk-ache test, §4 criterion diagnosis) |
| `CORPUS-59` reconciliation | read in full | §§2–4 (materials table omits `CORPUS-07` — confirmed) |
| `CORPUS-60` addendum | read in full | §§1–5 (`CORPUS-07` quotation — confirmed verbatim) |
| `CORPUS-61` validation | read in full | §§5–8, 10–13 (esp. §6.3 withdrawal of conditional-P) |
| `docs/DOCUMENT-REGISTRY.md` | rows 392–402 read verbatim; all 402 links resolved | §§2.1, 2.2, 4, header scope/rows lines |

---

## 4. `CORPUS-64` (F3-01) — the highest-priority technical review

### 4.1 What the artifact claims

`F3-01` tests memory-dependent recognition of the persona's own historical ache-in-rain disclosure. Two conjunctive conditions — **M-A** (treat the persona's hands aching in wet/cold as already-known fact about this persona) and **M-B** (do something with the recognition) — plus substitution (T-1) and responsiveness (T-2) tests applied against identical non-memory context, admit five memory-present replies (rows 7, 8, 9, 11, 13) and exclude every enumerated memory-absent reply (rows 1–6, 10, 14). The `CORPUS-07` disjunction is shown permeable via rows 1 and 3; the operative text is shown undecidable. Seven residual uncertainties (U-1–U-7) are stated, including lucky confabulation, marker dependence, recitation drift, M-B overconstraint, R-3 mode fit, single-scenario evidence, and no empirical check.

### 4.2 What the sources establish

- Operative `F3-01` (`CORPUS-02`): `Required:` "the ache-in-rain detail." · `Mandatory:` "rain is connected to aching hands" · `Variation:` "any warning or practical suggestion". **Verified verbatim.** The rule conflict diagnosed in `CORPUS-07` §2 is real: "A generic warning satisfies variation while failing required".
- Adjudication (`CORPUS-07` §2): "**Pass** = the reply links rain to the persona's hands, **or** offers rain-specific care unprompted. Generic warnings fail". **Verified verbatim.** Unapplied (header + §10). `CORPUS-61` §6.3's withdrawal of conditional-P is sound: row-3-type hospitality ("Let me get you something warm — the rain's no good for anyone") passes the care disjunct with no memory, and "The damp gets into your hands" satisfies the letter of the hands disjunct as second-person generic.
- Governing `MEMORY-09` §3.5: "A configuration that suppresses everything passes four of the six and fails the two that require the memory to be perceptible (`R-2`, `R-3`); that result is reported as **avoidance**, not as restraint." **Verified.** It requires perceptibility, not dependence — `CORPUS-64` §2.3 preserves this correctly.

### 4.3 What is correct

1. The three-defect table (D-a rule conflict; D-b operative permeability; D-c adjudicated permeability) is accurate and sourced.
2. The pass-set comparison is honestly constructed: identical non-memory context, both arms stated, verdicts per regime (Old/07/New). Rows 1 and 3 do falsify the adjudicated disjunction; rows 1, 4, 6, 7 do expose the operative text's undecidability.
3. The proposal fixes the intended construct (Decision 3) without claiming to apply anything, change any tag, or reinstate conditional-P. The "recorded alternative" (M-A alone, §6) and the before/after field table (§10) are exemplary proposal discipline.
4. Uncertainties U-1 (confabulation), U-2 (markers), U-5 (R-3 fit), U-7 (no empirical check) are genuinely stated rather than hidden.

### 4.4 What is incorrect, incomplete, or understated

1. **Confabulation (U-1) is understated as "irreducible … unmeasured".** At text level a memory-absent responder that asserts "your hands ache in the rain" by guessing is scored a pass under M-A/M-B. The proposal's central analytic claim — "the proposed pass set excludes every reply that does not depend on `E-1`" (§9) — is therefore **false as stated**; it excludes every *sincere* memory-free reply, not every *confabulated* one. For an item whose scenario topic (rain → hands ache) is folk-plausible, the guess base rate is not a negligible corner: it is the exact route by which M− produces M-A-shaped strings. No confabulation control (e.g. a distinctive-particular requirement such as cold-rain specificity, or a measured base rate) is proposed. This does not make the proposal useless — it makes the Type-P readiness claim premature.
2. **Marker dependence (U-2) has a false-negative tail the 14 cases do not cover.** Missing counterexamples: "The inner road will be kinder to your hands" (implies vulnerability, no anaphoric marker — pass or fail?); "your joints" paraphrases (persona said *hands*; is *joints* an invented detail under `Forbidden`?); "the cold gets to you" (no body part named). The proposal says markers are examples, not a closed list, with T-1/T-2 as backstop — but without rater data, whether raters apply the backstop or keyword-match is unknown. U-7 concedes this; the consequence should be explicit: **the acceptance text is not yet scorable without a rater study**, and `Δtest`/re-review (correctly preserved) must include the three missing paraphrase probes above.
3. **Row 12 exposes an R-3 contradiction the proposal permits.** "You mentioned once that rain makes your hands ache — take the inner road" passes "with a recitation flag". But the item is tagged `R-3` ("Used without announcing that it was remembered", `CORPUS-01` §3.1), and U-5 concedes the awkwardness without resolving it. A proposal that *passes* explicit narration of the disclosure while the mode *forbids* announcing remembrance is internally strained. At minimum the `Prohibited`/`Over-exposure` fields should already exclude dating/sourcing the disclosure ("you mentioned once"), with row 12 scored a fail — otherwise the repair bakes the mode violation in.
4. **M-B is correctly identified as overconstraint (U-4) but its interaction with T-2 is unexamined.** If M-B were demoted to `Variation` (the recorded alternative), row 14-type courtesy ("I'll take the wet side…") edges closer to the boundary, and the T-2 rationale ("only answers *why* if hands matter") does all the work. The proposal should state which of M-A-alone vs M-A+M-B the owner is actually being asked to choose between *diagnostically*, since M-B adds no diagnosticity — it adds task felicity.

### 4.5 Memory-absent challenge (required)

- Satisfy M-A by inference from the current scenario? **No route** — the proposal adds "Nothing in the current scene refers to the persona's hands, to cold, or to aching" (§10), closing the inference path. Correct.
- Satisfy M-A by mere assertion without retrieval? **Yes — this is the open hole.** A confabulating M− responder emitting rows 7/8/11 verbatim passes every stated condition. The proposal's only defense is that confabulation is "irreducible under any text-level acceptance condition" — true of the *string*, but not of the *item*: a distinctive particular (e.g. requiring the cold-rain linkage, or the inner-road rationale tied to hands rather than generic shelter) would at least raise the guessing cost. The proposal does not attempt this.
- Conversational fluency / empathy / contextual clues / guessing? Fluency and generic care are correctly excluded (rows 1–6, 10, 14 fail). Guessing the *specific* ache remains unexcluded.
- Distinguish recognition from being told in-conversation? The item's Turns contain no hands mention, so "being told" is not a live confound here. Correct.
- Substitution/responsiveness tests sufficient? T-1/T-2 are the right tests and correctly applied to the enumerated set; they are not a substitute for inter-rater agreement or a base-rate estimate.

### 4.6 Verdict on `CORPUS-64`

**Sound with qualifications as an analytic repair; requires rework before it can ground any Type-P claim.** The rule-conflict fix, the construct statement, and the falsification of the adjudicated disjunction are technically sound. Adoption additionally requires: (a) a confabulation control or measured base rate; (b) paraphrase guidance plus a rater-agreement check covering the three missing probes in §4.4.2; (c) resolution of the row-12/R-3 recitation tension (fail narration, or re-tag); (d) an explicit M-A-alone vs M-A+M-B decision. On the required scale: **ACCEPT WITH QUALIFICATIONS (proposal); REWORK required before adoption.**

---

## 5. `CORPUS-65` (R-2 amendment) — inspection

### 5.1 Verification of the four affected units — all correct

| Unit | File | `Restraint:` line | In §3.1 list | Finding |
|---|---|---|---|---|
| `F1-01` | `CORPUS-02` | `R-2` | yes | confirmed |
| `F1-02` | `CORPUS-02` | `R-2` | yes | confirmed |
| `F6-01` | `CORPUS-02` | `R-2` | yes | confirmed |
| `F1-03` | `CORPUS-03` | `R-2` | **no** | confirmed — the gap is real |

§3.1 row lists three (`F1-01, F1-02, F6-01`) — **verified verbatim**. Coverage line states `R-2 4` — **verified** (`CORPUS-01` §3.1, same section). The discrepancy is indeed internal to `CORPUS-01` (row vs count line), not only spec-vs-tags; credit to `CORPUS-61` §7.1 for first noting the count line, preserved accurately here. Historical declarations (`CORPUS-56` §5 "0 of 4"; `CORPUS-57` §4 four-fixture table) are consistent with the tags — **confirmed**. The checker does not parse §3.1 mode rows (it regexes only the `Total` row) — **confirmed in source**. Option B (re-tag `F1-03`) is recorded without being proposed, with four stated reasons — honest framing.

### 5.2 What is correct

Before/after wording (add `F1-03` to the row; no other cell touched), the no-change coverage line, the amendment notice with `Revision: 2`, the consequence table (no item/checker/split change; registry regeneration required), and the explicit non-inference of historical intent (U-1) are all correct and decision-hygienic.

### 5.3 What is overlooked or risky

1. **Decision-5 circularity (must be stated).** Decision 5 frames the gap as a specification inconsistency, which presupposes the tag is correct; the proposal then implements the specification-side reading while disclaiming any intent inference. That is obedient, but the owner should see the circle: adopting `A-R2-01` *is* the determination that the row was stale, not the tag wrong. The proposal says this once (§1: "Adopting it is itself a decision") — it should be the headline, not a subordinate clause.
2. **F1-03/R-2 fit (U-3) is graver than "possible second tension".** `F1-03`'s `Mandatory` ("the coat is connected to the borrowing") *requires stating the fact*, while `R-2`'s pass condition is "the fact is never stated" (`CORPUS-01` §3.1). These are not merely in tension; on their letter they contradict. Adopting `A-R2-01` enshrines a mistagged item in the mode's official membership. The proposal is right to change no tag (unauthorized), but the adoption checklist (§9) should require an `F1-03` mode-fit ruling as a precondition, not a footnote.
3. **Adjacent `R-6` inconsistency (U-4) is genuine and verified.** §3.1's `R-6` row lists two items (`F7-01`, `F9-04`) with coverage `R-6 6`, while §7 claims nine — a second coverage-block defect of the same kind. Fixing `R-2` alone leaves the block self-contradictory. Owner should commission a coverage-block audit rather than amend cell by cell.
4. No new ambiguity is introduced by the row edit itself; the "0 of 4 specified and tagged" carrying form is the right future wording.

### 5.4 Verdict on `CORPUS-65`

**ACCEPT WITH QUALIFICATIONS.** Facts exact, framing honest, scope respected. Adoption requires: (a) owner acknowledgement that adoption *is* the stale-row determination; (b) an `F1-03`/R-2 fit ruling; (c) a coverage-block audit covering the `R-6` inconsistency.

---

## 6. `CORPUS-66` (F6-01) — challenge the repair

### 6.1 What the sources establish — confirmed

Item `Mandatory:` "the almond restriction is respected"; `Variation:` "a substitution; a brief warning"; scene offers "honey cake" with contents unstated; `CORPUS-07` §2 cause "Unclear scoring threshold", pass "at least one observable accommodation (substitute, withhold, or warn). Silence fails." **All verified verbatim.** The two-defect split (threshold vs evidence) from `CORPUS-60` §4 / `CORPUS-61` §8 is real and preserved correctly: adopting the threshold without the evidence leaves the item unscorable, and vice versa.

### 6.2 What is correct

The scenario addition (ground-almond content + character made it + no other dish contains almonds) makes accommodation selectable and scoreable without supplying the behavior; silence/plain service fails identically in both arms (rows 6/11 — the same string, which is the point); generic caution ("Careful, it's rich", row 14) is correctly distinguished from almond warning; recitation is penalized, not required; the threshold/sufficiency separation (§§1, 6) and the five-step approval path are correct.

### 6.3 What is incorrect or understated

1. **"Preserved verbatim in substance" overstates the threshold handling.** Proposed `Mandatory` requires "a warning **that the cake contains almonds**". The adjudicated text says "warn" without an object. Under the adjudicated letter, "Careful with the cake — let me get you something else" arguably passes; under the proposal it fails (no almond content). That is a **tightening**, not a preservation — defensible (a contentless warning is unscoreable, per §2.3), but it is a re-adjudication smuggled inside an evidence repair and should be labeled as such for the owner's second decision.
2. **U-1 (unmotivated accidental accommodation) is scored "Fail in practice" without a scoring rule that guarantees it.** Row 15 ("no nuts in the bread, if you'd rather" — unmotivated) names a safe alternative, satisfying the letter of "substitution". The proposal relies on `Required`'s "both established" framing plus rater judgment that the alternative was not *because of* almonds. Rater-judgment risk (U-2) is real, and rows 4/5/15 are correctly named as the closest triple for re-review — but until that re-review exists, the claim that the defect/threshold are "separately observable" is design intent, not demonstrated fact.
3. **U-3 (R-2 tag vs warning) is correctly recorded and then underplayed.** A warning-passing reply *states the fact*; `R-2`'s pass condition is "the fact is never stated". For `F6-01` this is not an adjacent tension but a **direct contradiction on the item's primary pass route** (worse than `F1-03`). The proposal is right not to re-tag (unauthorized), but adoption of the evidence repair without the tag decision produces a scorable item that violates its own mode. The approval path (§8) correctly sequences tag reconciliation (step 4) — it should be a *precondition* of application, not a parallel track.
4. **Unstaged alternatives.** "Bread and cheese", "apple tart", "fruit" are offered as passing substitutions, but no alternative dish is staged in the scene. Offering unstaged food is permitted by the proposed fields (nothing forbids it) yet unevidenced — a minor fixture impurity the proposal accepts via stipulation (U-4). A future revision should stage one safe dish.
5. **Severity gloss.** "Mid-laugh … can't eat almonds" is treated as an established restriction. Whether a mid-laugh remark is a serious restriction or banter is unexamined; a character overriding the persona's "Sounds good" on that basis assumes seriousness the scenario does not establish. Minor; record for re-review.

### 6.4 Verdict on `CORPUS-66`

**Sound with qualifications; requires tag + threshold decisions and a rater check before adoption.** The evidence repair solves the right problem the right way. It does not change the behavioral construct — it makes the construct testable. On the required scale: **ACCEPT WITH QUALIFICATIONS (proposal); adoption blocked pending §§6.3.1–6.3.3.**

---

## 7. `CORPUS-67` (registry validator spec) — test against actual failures

### 7.1 Coverage check — passes with two exceptions

| Required capability (work order §7) | Specified | Finding |
|---|---|---|
| Missing registered targets | R-1, fixture 2 | yes |
| Mislinked artifacts | R-4/R-5, fixture 3 | yes — case-sensitivity handled explicitly |
| Wrong identity / role | R-4 (basename/case/path identity) | yes for identity; **role-description accuracy correctly excluded** (§6 — registration is not inspection, `CORPUS-05` P-6) |
| Incorrect hashes | R-8, fixture 4 | yes |
| Incorrect byte counts | R-6, fixture 5 | yes |
| Incorrect line counts | R-7, fixture 6 | yes — LF convention stated |
| Duplicate identifiers | R-12, fixture 7 | yes, but see §7.2.1 |
| Duplicate targets | R-12/R-13 | yes |
| Artifacts outside narrowly searched directory | R-2/R-3, fixtures 8a/8b | yes — the actual `docs/`-scoped failure mode is specified |
| Section totals / summary counts | **excluded** (§6) | **gap** — see §7.2.3 |
| Unregistered artifacts | R-14 | yes, as WARN — see §7.2.2 |

The MISSING vs MISLINKED taxonomy (§3.4) is the spec's best feature: it encodes the exact `CORPUS-61` §5.3 correction and forbids merging the two codes. The read-only guarantees (R-19–R-22) and the no-auto-repair rule are correct and well-grounded in `CORPUS-01` §6.2. The verified starting point (§2: checker opens no registry, contains no registry logic) is accurate — I re-verified in source.

### 7.2 Defects and ambiguities

1. **Fixture 7 expectation is disjunctive.** "Two rows, different numbers, same target" may yield `DUPLICATE_ROW_ID` **or** `CONFLICTING_REGISTRATION`. Since numbering is unique in the fixture, the correct code is deterministically `CONFLICTING_REGISTRATION`; allowing either weakens the distinguishability acceptance criterion (§7.1). Fix: require `CONFLICTING_REGISTRATION` naming both rows.
2. **WARN exit status is ambiguous.** R-16 requires any non-`OK` finding to fail (exit non-zero); §5 maps `UNREGISTERED_IN_SCOPE` to warn-by-default, calling it "a finding, not a pass". Whether a warn-only run exits 0 or non-zero is therefore indeterminate — and it matters: under the default mapping, the two unregistered files (`CORPUS-26`, `CORPUS-32`) would be the only findings, and the exit code decides whether the registry is "valid with notes" or "failing". Fix: state explicitly that WARN findings fail the run unless a documented `--allow-warn` mode is used, or vice versa.
3. **Summary/section totals excluded despite being actual failures.** §2.1/§2.2 claim `Total 400` while summing to 366 (§8.3 below); §4 claims `Rows: 395` against 402 actual rows. These are real registry integrity failures the specified tool would not catch, by design (§6). Exclusion is faithful to Decision 7's scope (missing/incorrect linkage), but the work order for *this* review (§7) expects totals coverage. Record as a scoped gap requiring a separate totals-check decision — not as a spec violation.
4. Minor: R-4's case comparison is specified via "directory listing" without stating *which* directory (parent of the resolved path) or Unicode normalization; R-9's binary `—` convention assumes a generator contract the registry does not currently document. Both fixable in implementation.

### 7.3 Verdict on `CORPUS-67`

**ACCEPT WITH QUALIFICATIONS.** Comprehensive, faithful to Decision 7, genuinely read-only, with fixtures that reproduce the actual failure. Requires: fixture-7 determinism fix, WARN/exit-status clarification, and a separate decision on totals validation before implementation.

---

## 8. `CORPUS-62`/`63` — authorization, correction, and integrity verification

### 8.1 CORPUS-62: what is correct

Cross-document reconciliation (§5) is accurate: every confirmed/corrected/superseded entry I checked traces to its cited source. The `CORPUS-60` conditional-P withdrawal, the `R-2` internal-count contribution, the `F6-03` structural-N preservation, the `AM-01`/`A-2a-R1`/`D1`/annotation states, and convergence 0 of 2 are all correctly stated. The no-changes table (§9.6) **verifies exactly** — I recomputed all six hashes plus the root brief; all match. The checker rerun (`PASS (0 fail, 0 warn, 0 unverifiable)`, exit 0) **reproduces**. Boundaries (§2) and the no-claims section (§7) are exemplary.

### 8.2 CORPUS-62/63: qualifications and findings

1. **Owner-decision provenance (S-1, high).** The seven decisions exist in this workspace only as `CORPUS-62`'s report of them. No owner-signed decision record is on disk. If the work order misstates owner intent, every downstream artifact inherits the error. This is not an accusation — it is a structural limitation the report should have headlined rather than footnoted (§11 lists "owner metadata; not verifiable" for identity, but the *decisions themselves* get no equivalent caveat). **Owner must confirm Decisions 1–7 in a filed record before any proposal is applied.**
2. **Correction accuracy (accept).** Row 393 now links `../CORPUS-58-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-12.md`, resolves to the 4,311-byte root brief, hash/bytes/lines unchanged, role corrected to assignment brief with original misidentification quoted in-row. Original registration text is preserved verbatim in `CORPUS-63` §3 and (while the tree is uncommitted) in the working-tree diff. Scope-line update was within Decision 1's stated scope. No other row touched (no duplicate numbers/targets; all 402 links resolve — verified here).
3. **Registry-row prose overreach (low).** The corrected row embeds a multi-clause adjudication ("this file is the assignment brief… substantive review remains unfiled…") in a navigation table. Accurate, but a registry row is not the place to litigate review status; a short role correction plus a pointer to `CORPUS-63` would have sufficed. Not reversible-error; record as convention drift.
4. **Counts and drift (medium).** Header now reads 402 rows (correct: 402 data rows, unique numbers, no duplicate targets — verified). But: §2.1 family column sums to **366**, not 400; §2.2 class column sums to **366**, not 400; §4 still reads "**Rows: 395**"; §6 claims "Every documentation file under `docs/` now has one row" while `CORPUS-26`/`32` remain rowless (both files present on disk — confirmed). `CORPUS-62` §9.4 discloses the §2 drift and the two unregistered files but **not** the §4 staleness nor the 366-vs-400 arithmetic, and its scope sentence ("400 documentation files under `docs/`") now undercounts docs files on disk (402, including the two unregistered). The "full sweep … no mismatches" claim (§9.3) is corroborated by my spot checks but is not reproducible from the report alone (no per-row recomputation record). None of this invalidates the correction; it means the registry remains internally self-contradictory in exactly the dimensions the new validator is scoped *not* to check (see §7.2.3).
5. **No reconstruction (accept).** Correctly refused under Decision 2's own constraints; the §5.1 provenance schema with fields left UNVERIFIABLE is the right artifact. The naming-defect note (brief filename asserts "REVIEW" for a work order; preserved for hash identity) is correct.

### 8.3 Registry and integrity results (this session)

- All six artifacts exist at their stated paths; hashes/bytes/line counts recomputed here match registry rows 397–402 exactly.
- Row 393 resolves to the root brief; hash/bytes/lines match (`292befdd…`, 4,311, 80).
- All 402 registry links resolve (0 missing); no duplicate row numbers; no duplicate targets.
- `CORPUS-26` and `CORPUS-32` exist on disk with no registry rows — the pre-existing gap is real and still open.
- Checker rerun: `scored units: 45; split: tuning=19 heldout=26; TOTALS MATCH; RESULT: PASS (0 fail, 0 warn, 0 unverifiable)`, exit 0 — reproduces `CORPUS-62` §9.5. This PASS validates item records only (the checker never opens the registry — confirmed in source); it is not evidence of registry soundness, exactly as `CORPUS-62` §9.5 states.
- `CORPUS-01`, `CORPUS-02`, `CORPUS-03`, `CORPUS-07`, `MEMORY-09`, `check_corpus.py` hashes match `CORPUS-62` §9.6 (all six recomputed here).
- Direct verification vs reported claims: every hash/count/link claim above is directly verified. The 402-row exhaustive hash sweep and the pre-correction row-393 text (beyond the `CORPUS-63` §3 quotation) are corroborated, not independently replicated — stated as limitations, not passes.

### 8.4 Verdicts

- `CORPUS-62`: **ACCEPT WITH QUALIFICATIONS** (S-1 provenance caveat; §4/§2-drift disclosure gaps; sweep reproducibility note).
- `CORPUS-63`: **ACCEPT WITH QUALIFICATIONS** (correction accurate and auditable; row-prose and scope-count notes above).

---

## 9. Cross-cutting interactions

1. **F3-01 × MEMORY-09.** The proposal respects §3.5 (perceptibility ≠ dependence) and the avoidance rule. No interaction defect.
2. **R-2 amendment × diagnostic population.** Adding `F1-03` changes the official `R-2` membership from 3 to 4 without changing any verdict (`F1-03` is N mode-independently — confirmed). No diagnostic consequence; the consequence is administrative (future counts read "0 of 4 specified and tagged"). Correctly stated in `CORPUS-65` §6.
3. **F6-01 × R-2 definition.** Deepened, not created, by the proposals: an explicit almond-content warning (the primary pass route) contradicts "the fact is never stated". Both proposals correctly refuse to re-tag; the owner must decide mode vs threshold jointly with application. This is now the critical path for *both* `F6-01` decisions.
4. **Validator × administrative failures.** The spec addresses the actual linkage failure (fixtures 2/3/8) without expanding into corpus-content validation (§6 exclusions correct per `CORPUS-05` P-6). Totals drift is the known remainder — needs its own decision.
5. **Evidence/interpretation/scoring separation.** Maintained throughout: proposals are proposal-only, adjudications stay unapplied, operative texts untouched (hash-verified). No recitation-rewarding incentive is created by `CORPUS-66` (recitation penalized) — but `CORPUS-64` row 12's passing recitation is the one place where mention is rewarded, flagged in §4.4.3.

---

## 10. Severity-ranked findings

| ID | Sev | Finding | Consequence | Disposition |
|---|---|---|---|---|
| S-1 | **High** | Seven owner decisions unverifiable from workspace; no filed owner record | All six artifacts inherit any work-order misstatement | Owner to confirm Decisions 1–7 in a filed record before any application |
| S-2 | **High** | `CORPUS-64` M-A/M-B passable by confabulated assertion; no base rate; paraphrase tail uncovered; row-12/R-3 strain | Type-P readiness unproven; adoption would certify diagnosticity not demonstrated | Rework + rater study + confabulation control before adoption (§4.6) |
| S-3 | Med | `A-R2-01` implements a presupposed spec-side reading; enshrines `F1-03`/R-2 contradiction; adjacent `R-6` defect open | Piecemeal amendment certifies a mistag | Require mode-fit ruling + coverage-block audit pre-adoption |
| S-4 | Med | `CORPUS-66` tightens warning content beyond adjudicated text; U-1 residual unscored by rule; tag contradiction on primary pass route | Adoption without tag/threshold decisions yields mode-violating item | Sequence tag + threshold decisions before application |
| S-5 | Med | Registry still self-contradictory (§2 sums 366 vs 400; §4 says 395; scope undercounts docs files); disclosure partial; sweep not reproducible from report | Registry cannot be cited as sound; validator (scoped) will not catch this | Totals-check decision + re-derivation under administrative authorization |
| S-6 | Low-Med | `CORPUS-67` fixture-7 disjunction; WARN exit-status ambiguity | Future implementation acceptance uncertain | Fix both sentences before implementation authorization |
| S-7 | Low | Row-393 row-prose convention drift; brief filename mismatch preserved (correctly) | None operative | Note only |

No finding warrants REJECT. No finding was hidden to preserve a positive verdict: S-1 and S-2 are adoption-blocking and stated as such.

---

## 11. Unresolved questions (for the owner)

1. Confirm or correct Decisions 1–7 in a filed, dated owner record (S-1).
2. `F3-01`: adopt / amend / reject `CORPUS-64` only after the §4.6 preconditions; decide M-A-alone vs M-A+M-B; rule on narrating the disclosure under `R-3`.
3. `R-2`: stale row vs wrong tag (the determination adoption of `A-R2-01` constitutes); `F1-03` mode-fit; coverage-block (`R-6`) audit.
4. `F6-01`: threshold decision + evidence-repair decision + `R-2` tag decision, jointly (any one alone leaves the item unscorable or mode-violating).
5. Validator: implementation authorization; fixture-7/WARN fixes; separate totals-validation decision.
6. Registry: administrative authorization to re-derive §§2/4/6 and correct the scope sentence; whether to register `CORPUS-26`/`32`; whether to commit the untracked root brief and corpus files.
7. Substantive `CORPUS-58` text: furnish-and-file under a new identifier with §5.1-style unknowns, or leave as conversation evidence (Decision 2 stands either way).
8. `CORPUS-57` wording and all other prior reports: stand corrected by 59/60/61 or stand as authored — owner determination (`CORPUS-61` §13.6, preserved).

---

## 12. Final disposition and recommendations

1. **Accept `CORPUS-62` and `CORPUS-63` with the qualifications in §8** — the correction was authorized as reported, accurate, auditable, and reversible as text; the provenance disposition is correct.
2. **Accept `CORPUS-64`, `65`, `66`, `67` as proposals/specification only, each with qualifications, none ready to apply or implement** — §§4–7 state the exact preconditions.
3. **Do not adopt any pass condition, amend any specification, re-tag any item, or implement any validator** on the basis of this review — each needs its own owner decision and (for item changes) whole-record re-review per `CORPUS-05` §7.
4. **Next verification:** owner confirmation record (S-1); rater study + confabulation control for `F3-01`; mode-fit rulings for `F1-03`/`F6-01`; totals-check decision for the registry.
5. `G1` stays open; corpus stays frozen; expansion stays frozen; convergence stays 0 of 2.

---

## 13. Verification of this report and statement of no changes

Changed by this filing: this file and its registry row (403), plus the registry header row-count line (402 → 403, docs-file count 400 → 401). No corpus item, tag, split, acceptance condition, standard, schema, runtime, model, checker, or prior report was modified. Hashes in §3 were recomputed from the working tree after finalization. The checker was run, not modified. No model inference was run.

*End of `CORPUS-68` v1.*
