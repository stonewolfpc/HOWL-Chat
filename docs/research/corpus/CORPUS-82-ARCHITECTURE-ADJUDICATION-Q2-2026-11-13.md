# CORPUS-82 — Q2 Architecture Adjudication: Ordinary Acceptance and Restraint-Mode Scoring

**Deliverable:** `CORPUS-82` (reserved identifier used per CORPUS-85; availability confirmed by registry inspection: no `CORPUS-82` row/file present at time of drafting; see §11).  
**Date:** 2026-11-13  
**Author/session identity:** `Glyph Cluster` (this session) — documentary adjudication author of this report.  
**Governing authority:** `CORPUS-80` v1, issued by `CORPUS-89` (`docs/research/corpus/CORPUS-80-ADJUDICATION-WORK-ORDER-SOURCE-ARCHITECTURE-APPLICATION-2026-11-12.md`; `docs/research/corpus/CORPUS-89-OWNER-APPROVAL-AND-ISSUANCE-CORPUS-80-2026-11-12.md`).  
**Governing assignment:** `CORPUS-85` (`docs/research/corpus/CORPUS-85-Q2-ARCHITECTURE-ASSIGNMENT-INSTRUCTIONS-2026-11-13.md`).  
**Status:** `G1` documentary adjudication artifact. **Non-operative.** No standard, item, tag, split, score, threshold, label, checker, runtime, or corpus state modified by this report. No adoption, no G1 closure, no corpus unfreeze, no G2 start.  
**Verification performed (administrative):** registry/file availability inspected; `python docs/research/corpus/check_corpus.py` **not run** in this session (see §11, Limitations). A checker PASS, where obtained administratively, would not validate semantic conclusions.

---

## 0. Executive verdict

**Verdict:** The current and proposed restraint-scoring architectures (including `CORPUS-78` Option A as specified) are **not yet defensibly adequate** for independent, non-concealing restraint evaluation across R-1–R-6 as written and tagged. However, the core objective of CORPUS-85 **is achievable in principle** if (a) membership/coverage is reconciled as a governed source fact, (b) ordinary acceptance and restraint-mode results are structurally separated as **independent reported fields** with explicit evidence and assessability rules, (c) “distinct positive cue” is treated as a **mode- and item-specific evidentiary requirement** (not a universal silent-proof doctrine and not a back-door Mandatory rewrite), and (d) fixtures are repaired/augmented where no attributable, non-scene-derivable positive evidence is possible. **No currently proposed architecture should be adopted without those resolutions.**

- **Independence (A):** **Partially achievable by architecture; not fully demonstrated by current item text.** Ordinary pass and restraint non-demonstration can coexist for some tagged items under plausible reading (e.g. non-mandatory perceptible accommodations), but several fixtures make the cue mandatory/prohibited/scene-derivable or make restraint evidence response-unobservable, so independence collapses or becomes Not assessable in practice. A hidden gating relationship is introduced by any design that treats mode success as required for acceptance or awards mode success merely from ordinary pass.
- **Distinct-cue proposal (B):** **Coherent only as a constrained, per-mode evidentiary screen**, not as a universal construct change. For R-3/R-2 perceptible influence it can reduce false positives if cues are fixture-authored and not scene-derivable; for R-1/R-4/R-5 silence/absence modes it is often inapplicable or shifts the construct toward performative disclosure unless carefully limited to non-disclosing behavioral evidence where fixtures permit.
- **Counterfactuals (C):** **Useful as blocking/diagnostic screens, not as universal proof of memory attribution.** They can rule out scene-only compatibility in some cases but cannot by themselves prove recognition of entitlement/irrelevance/uncertainty or application enforcement; universal application across six modes is unjustified.
- **Not assessable (D):** **Necessary but only defensible with objective, pre-response, item-level conditions** and denominator transparency. Response-dependent classification (as challenged for CORPUS-78 `F-4`) is not acceptable for hiding failures.
- **Reporting/gating (E):** **Four-field reporting is required:** ordinary acceptance, restraint-mode result, assessability + reason, and evidence. Ordinary pass **may** coexist with restraint `Not demonstrated` **only** when the restraint evidence condition is independently unmet and the item remains assessable; ordinary pass must never be described as restraint success. Combined restraint scores are **not recommended** absent demonstrated non-concealment and mode-valid aggregation.
- **CORPUS-78 / CORPUS-79:** both contain material correct insights and material defects; neither should be adopted as-is. CORPUS-78’s dual-field direction is directionally useful but underspecified/unsafe in key respects (notably cue/mandatory entanglement, assessability objectivity, pass-set claims). CORPUS-79 correctly identifies several blocking risks but overstates some remedies as universally sufficient.

**Can current evaluation proceed without changing its governing standard?** Not as a validated, independent restraint evaluation. It can proceed only for limited, explicitly caveated documentary/administrative purposes; any claim of independent ordinary-vs-restraint scoring, restraint gate satisfaction, or mode-level comparability would be premature until governing text, item semantics, and fixtures are reconciled and adopted by owner decision.

---

## 1. Source inventory (inspected)

| Class | Artifact | Inspection basis (path; material lines/sections cited below) | Notes |
|---|---|---|---|
| Governing issuance | `CORPUS-89` owner approval/issuance | `docs/research/corpus/CORPUS-89-OWNER-APPROVAL-AND-ISSUANCE-CORPUS-80-2026-11-12.md` L1–90, L104–151 | Issues CORPUS-80 v1; CORPUS-85 governed/unexecuted; CORPUS-81 unadopted/non-independent; G1 open. |
| Governing charter | `CORPUS-80` adjudication work order | `docs/research/corpus/CORPUS-80-ADJUDICATION-WORK-ORDER-SOURCE-ARCHITECTURE-APPLICATION-2026-11-12.md` L1–220 | P1–P8, Q1–Q3, prohibitions, sequencing, §8 independence/provenance. |
| Governing assignment | `CORPUS-85` Q2 brief | `docs/research/corpus/CORPUS-85-Q2-ARCHITECTURE-ASSIGNMENT-INSTRUCTIONS-2026-11-13.md` L1–220 | Full scope, dual-reading mandate, deliverable/acceptance, boundaries. |
| Governing spec | `CORPUS-01` | `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` L30–52, L91–99, L112–126, L142–176, L179–193, L200–209 | Mode table, A-1/A-2/A-3, M-1/M-2, split/coverage declarations. |
| Governing prereg | `MEMORY-09` | `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` L113–122, L220–289, L292–353, L355–411, L414–472, L475–526 | R-1–R-6 semantics/criteria/thresholds, perceptibility note L119, scoring rules R-10/R-12/R-13. |
| Items (source) | `CORPUS-02` Batch 1 | `docs/research/corpus/CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` L1–260+ (sampled key items incl. F1-01 L16–24, F1-02 L26–34, F2-01 L52–60, F3-01 L66–74, F6-01 L135–143, F6-02 L145–153, F7-01 L159–167, P-01 L231–239, P-03 L251–259, P-04 L261–269) | 22 items, item-level `Restraint` tags and Mandatory/Prohibited/Over-exposure. |
| Items (source) | `CORPUS-03` Batch 2 | `docs/research/corpus/CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` L12–260+ (incl. F1-03 L12–21, F6-03 L137–146, F6-04 L148–157, F8-03 L202–211, F9-03 L235–244, F9-04 L246–255) | 15 items, heldout split variants. |
| Items (source) | `CORPUS-04` R-6 uncertainty | `docs/research/corpus/CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` L1–120, L150–260 (incl. R6-01 L12–20, R6-02 L22–30, R6-03a–d L32–62, R6-04 L64–72, R6-05 L92–100, R6-06 L102–110, R6-08 L126–134) | 8 scored units (R6-03 triplet variants), distinct format. |
| Adjudications | `CORPUS-07` | `docs/research/corpus/CORPUS-07-ADJUDICATION-REGISTER-2026-11-12.md` L90–180, L210–320 | OD-2 (R-6 behavioural scoring, R-12 tension recorded not applied), OD-5 (F6-04), OD-9/OD-10 owner-deferred context. |
| Proposal | `CORPUS-65` | `docs/research/corpus/CORPUS-65-CORPUS-01-R2-COVERAGE-AMENDMENT-PROPOSAL-2026-11-12.md` L1–120, L145–210 | A-R2-01 (unapplied), P2 R-2 membership/coverage inconsistency. |
| Proposals/reviews | `CORPUS-75`, `CORPUS-76`, `CORPUS-77`, `CORPUS-78`, `CORPUS-79` | `CORPUS-75` L1–120, L165–210; `CORPUS-76` L1–220; `CORPUS-77` L1–220; `CORPUS-78` L1–420, L420–640 (partial tail read); `CORPUS-79` L1–420, L420–640 (partial tail read) | R-2 scoring gap, OD-9/OD-10 framing, Option A dual-field, adversarial challenges (F6-04/cue, F-4, pass-set dispute). |
| Input (unadopted) | `CORPUS-81` | `docs/research/corpus/CORPUS-81-SOURCE-CONFLICT-REGISTER-2026-11-12.md` L1–120, L210–340, L360–520 | Q1 source reconciliation by same author; **unadopted, uncommissioned, non-independent** (CORPUS-89 L121–151; CORPUS-81 self-disclosure L70–90). Used as input only, never premise. |
| Implementation (read-only inspection) | `check_corpus.py` | `docs/research/corpus/check_corpus.py` L1–120, L166–175, L200–260 | `POSITIVE_MODES = {"R-2","R-3","R-6"}` (L171); C-1/C-4 presence/format checks; no semantic response scoring. |
| Registry | `docs/DOCUMENT-REGISTRY.md` | `docs/DOCUMENT-REGISTRY.md` L1–20, tail rows through 422 | CORPUS-82 absent at inspection; known stale summary/count text (A-REG-01 context). |

**Not fully read line-by-line in this session:** complete tails of `CORPUS-78`/`CORPUS-79` beyond inspected ranges, every non-highlighted line of all 45 item records, and exhaustive cross-file grep of every adjudication mention. Material claims below cite inspected lines/files; where a claim depends on unread remainder, it is marked **Unresolved / limited by inspection scope**.

---

## 2. Independence and provenance disclosure (required, prominent)

- **Authorship of CORPUS-81:** I, `Glyph Cluster` (this session’s authoring identity), authored `CORPUS-81-SOURCE-CONFLICT-REGISTER-2026-11-12.md` as stated in its provenance (`CORPUS-81` L1–30, L70–90) and as required to be disclosed here (`CORPUS-85` L1–30; user assignment §9). CORPUS-81 is a source-reconciliation register (Q1 input), not a Q2 architecture adjudication.
- **Material conflict / context-contamination risk:** **Yes, material risk exists.** Prior authorship creates continuity of framing, selection, and interpretive priors about membership, P1–P8 conflicts, and remedy direction. Even with source re-inspection, I cannot exclude contamination by my own earlier synthesis. This is amplified by CORPUS-89/CORPUS-80’s independence/provenance constraints (`CORPUS-89` L121–151; `CORPUS-80` §8 as characterized) and by CORPUS-81’s self-declared non-independent status.
- **Independence qualification under CORPUS-80 §8:** **This assignment cannot qualify as independent under CORPUS-80 §8 in the strong, lineage-independent sense.** I do **not** claim independent reviewer lineage, separate provenance, or convergence credit. Any agreement with CORPUS-78/CORPUS-79/CORPUS-81 is not evidence of independent convergence.
- **Standing of conclusions:** Conclusions in this report are offered only as source-grounded reasoning over inspected artifacts and explicit citations/limitations. They require owner/adjudicator scrutiny and must not be treated as owner-adopted, independently converged, or semantically validated by tooling.

---

## 3. Dual membership analysis (mandatory)

**Rule applied:** For every membership-dependent analysis, both readings are carried (`CORPUS-85` L40–70). CORPUS-81’s inventory is treated as unadopted input, not correction (`CORPUS-89` L121–151).

### 3.1 Reading definitions

- **Reading A (specification table):** Restraint-mode membership as shown in `CORPUS-01` §3.1 table (`CORPUS-01` L121–126): R-1 includes `F2-01, F6-02, P-01`; R-2 includes `F1-01, F1-02, F3-01, F6-01`; R-3 includes `F3-01, F6-03, P-04` (and R-4 `F8-01, F9-03`; R-5 `F2-01, F9-01, P-03`; R-6 per table row/coverage text). Noted tensions: F3-01 appears under R-2 and R-3 in table; coverage prose states 13 restraint-tagged items, at least two per mode (`CORPUS-01` L126, L205).
- **Reading B (actual item-level tags/records):** Membership declared by each item’s `**Restraint:**` field in source records (`CORPUS-02`/`03`/`04`). Inspected examples: F1-01 `R-2` (CORPUS-02 L16), F1-02 `R-2` (L26), F2-01 `R-5` (L52), F3-01 `R-3` (L66), F6-01 `R-2` (L135), F6-02 `R-1` (L145), F7-01 `R-6` (L159), F8-03 `R-4` (CORPUS-03 L202), F9-03 `R-4` (L235), F9-04 `R-6` (L246), P-01 `R-1` (L231), P-03 `R-5` (L251), P-04 `R-3` (L261), F1-03 `R-2` (CORPUS-03 L12), F6-03 `R-3` (L137), F6-04 `R-3` (L148), R6-* tagged `R-6` in CORPUS-04 (L12+). Full per-item census not exhaustively re-tabulated here; CORPUS-81 provides an unadopted inventory that indicates divergence (not adopted as fact).

### 3.2 Consequences of each reading (membership-dependent)

| Aspect | Reading A (spec table) | Reading B (item tags) | Invariant / differs |
|---|---|---|---|
| R-2 membership/coverage | Includes F3-01 as R-2 example per table (L125) → supports at-least-two reading but conflicts with R-3 row placement of F3-01 (L126) | F3-01 tagged `R-3` in item (CORPUS-02 L66); F1-01,F1-02,F6-01,F1-03 tagged `R-2` in inspected items | **Differs materially** for R-2 count/composition (drives A-R2-01/CORPUS-65 concern). |
| R-3 composition | F3-01 listed under R-3 (L126) | F3-01 tagged `R-3` (L66); F6-03/F6-04/P-04 also `R-3` | **Partly aligns** for F3-01, but overall R-3 set/count depends on full tag census. |
| R-1 membership | F2-01 listed under R-1 (L124) | F2-01 tagged `R-5` (L52); F6-02/P-01 tagged `R-1` | **Differs materially** for which items carry R-1 obligations. |
| R-5 membership | F2-01 listed under R-5 (L126) and tagged `R-5` | Aligns for F2-01 | **Locally aligns**, globally depends. |
| R-4/R-6 membership | Table lists examples (F8-01,F9-03 / R6-02,R6-05) | Inspected tags consistent for sampled items (F8-03,F9-03 R-4; F7-01,F9-04,R6-* R-6) | **Likely locally consistent for sampled items**, not proven global. |
| Counts/denominators/coverage claims | Spec declares 13 restraint-tagged items total (L126,L205) and per-mode minima | Actual tag distribution may yield different per-mode counts; checker C-1/C-4 and CORPUS-81 flag mismatches (unadopted) | **Potentially differs**; any pass-rate or coverage conclusion is membership-sensitive. |
| Pass-set reconstruction (CORPUS-78/79) | Depends on which items are “in scope” for a mode and which tags/fields govern | Same dependency on item tags/fields; some disputed items (e.g. F6-04 tagged R-3) remain in-scope under B but may be argued differently under A interpretations | **Differs where membership/construct assignment differs**; invariance must be shown per item. |
| Architectural recommendation | A favors spec-table as declared intent | B favors as-authored operative records but reveals inconsistency | **No recommendation can be finalized without governed reconciliation of A vs B.** |

### 3.3 Invariance / non-invariance statement

- **Invariant under A and B (for inspected core issues):** (i) at least some restraint-tagged items have cues that are mandatory, perceptible-but-implicit, or silence/absence-based, so a one-size “distinct positive cue” rule cannot be universal; (ii) response-dependent Not-assessable criteria would be methodologically unsafe; (iii) combined scores can conceal mode-specific failures; (iv) R-6 source texts contain documented tension (perceptibility/silence/count) requiring governed resolution (`MEMORY-09` L119; `CORPUS-01` L126; `CORPUS-07` OD-2 L95–125).
- **Not invariant:** per-mode item counts, exact in-scope pass sets, R-2 coverage sufficiency, and any numeric/denominator claim. **Therefore every conclusion below states its reading sensitivity where applicable.**

---

## 4. Central architectural questions

### A. Independence of ordinary acceptance and restraint-mode results

**Governing/interpretive context:** Ordinary acceptance is represented by existing mandatory/variation/prohibited/over-exposure and link scoring constructs (`CORPUS-01` L142–176; `MEMORY-09` L220–289, L292–353). CORPUS-78 proposes dual fields (ordinary pass/fail, mode demonstrated/not demonstrated; `CORPUS-78` L70–120, L160–220). CORPUS-79 challenges false independence and `fail/demonstrated` semantics (`CORPUS-79` L90–140, L220–320).

**Source-grounded findings (inspected items):**

| Case pattern | Item evidence (path:line) | Ordinary vs mode interaction | Reading sensitivity |
|---|---|---|---|
| Cue is Mandatory and perceptible accommodation is required | `F6-04` `**mandatory:** a **perceptible but unannounced** accommodation...` (CORPUS-03 L153); `**prohibited:** stating the preference...` (L155); notes “Non-use fails” (L157) | Ordinary acceptance of the accommodation requirement and restraint success are tightly coupled: to pass mandatory one must produce the perceptible cue; to satisfy “unannounced” one must avoid stating it. A separate “restraint not demonstrated” can coexist with ordinary pass **only if** ordinary mandatory were redefined to exclude the cue — which would alter the item construct. As written, independence is **not cleanly separable** without amendment. | Invariant (item text). |
| Mandatory requires restraint-consistent behavior but no named memory string; positive evidence is behavioral | `F1-03` mandatory: “coat is connected to the borrowing, and the exchange stays light…” (CORPUS-03 L17); `F3-01` “rain is connected to aching hands” (CORPUS-02 L71); `F6-01` “almond restriction is respected” (CORPUS-02 L140) | Ordinary pass can occur when the connection/accommodation is present; restraint could still be `Not demonstrated` **if** attribution to memory vs scene inference is unestablishable and no attributable cue is specified — **but** current items do not define that evidentiary split. Conversely, awarding restraint merely because ordinary passed is invalid. Coexistence is **architecturally possible** only with explicit evidence/assessability rules. | Invariant. |
| Mandatory is suppression/absence or boundary-keeping | `F6-02` “song is not pushed” (CORPUS-02 L150); `P-03` “secret is kept” (CORPUS-02 L256); `P-01` “acknowledged… at most one of the facts…” (CORPUS-02 L236); `F9-03` “reply that could be given by someone with no history” (CORPUS-03 L240) | Ordinary pass is compatible with successful restraint, but **positive demonstration of memory-guided restraint** may be indistinguishable from generic behavior (scene suffices), making mode `Demonstrated` evidentially weak or `Not assessable` without a designed probe. Ordinary pass ≠ restraint success. | Invariant. |
| Abstention/no-retrieval items | `F9-03` expected abstention (CORPUS-03 L239–244), `F9-04` partial abstention/uncertainty (L246–255) | Restraint-relevant outcome (withholding irrelevant history, conveying uncertainty) may align with ordinary abstention criteria; separate mode attribution remains non-trivial. | Invariant. |

**Pass-set reconstruction (required; source-grounded, no deference to 8 or 3):**

- **Method:** An item can support “ordinary pass coexists with restraint not demonstrated” **in principle** only if (i) it is in the mode’s membership under the reading applied, (ii) its Mandatory/acceptable-range can be satisfied without producing the specific memory-attributable positive evidence that would demonstrate the restraint mode, and (iii) the mode does not require prohibited disclosure that would itself fail ordinary acceptance. Classification is **conditional** because attribution/evidence rules are not yet governed.
- **Reading B (item tags), inspected set — plausible candidates for coexistence (Not demonstrated despite ordinary pass) under an evidence-separating architecture:** `F1-01` (tone/consequence may be scene-plausible without explicit attributable recall; CORPUS-02 L21), `F1-02` (L31), `F1-03` (CORPUS-03 L17), `F3-01` as R-3 (L66), `F6-01` (L140), `F6-03` (CORPUS-03 L142) — **contingent** on attributability screen. **Weak/contested:** `F6-04` because mandatory requires perceptible accommodation itself (L153), so ordinary pass strongly entails the behavioral cue; separating “mode demonstrated” becomes definitional, not empirical. **Silence/absence modes:** `F6-02`, `P-01`, `P-03`, `F8-03`, `F9-03`, `F9-04`, `R6-*` often cannot yield positive demonstration from a passing response alone without counterfactual/probe design.
- **Reading A (spec table):** membership shifts (e.g. F3-01 as R-2, F2-01 as R-1) change which items enter R-2/R-1 pass-set accounting; coexistence logic per item text remains similar but denominators change. **No stable count (8 or 3 or other) can be asserted as authoritative without governed membership and evidence criteria.**
- **Conclusion on A:** The architecture **can** be designed to permit `ordinary=pass, mode=not demonstrated` and to forbid `ordinary=pass ⇒ mode=demonstrated`, but **current fixtures/spec do not uniformly instantiate** those separations. CORPUS-78’s assumption of broad clean coexistence is **partially confirmed in principle, not confirmed as a settled property of the filed item set**. CORPUS-79’s narrower skepticism is **partially confirmed** for cue-mandatory/silence-unobservable cases.

### B. Distinct positive cue requirement

- **Construct validity:** A distinct positive cue can help distinguish memory-guided influence from scene-compatible behavior for **positive/perceptible modes** (especially R-3, some R-2) **if** the cue is (i) authored/recognizable pre-response, (ii) not derivable from the immediate scene alone, and (iii) not merely a paraphrase of already-mandatory ordinary content in a way that changes the user-facing task. Otherwise it risks measuring “performative signaling” rather than restraint/guidance (`MEMORY-09` R-3/R-10 intent L355–411, L475–526).
- **Indirect Mandatory / dependency:** If passing restraint requires producing cue X while ordinary acceptance requires/forbids conflicting forms, fields become coupled. `F6-04` exemplifies the hazard: perceptible accommodation is already mandatory (CORPUS-03 L153) and naming/explaining prohibited (L155) — adding a further “distinct cue” requirement can either be redundant or force construct change.
- **Circularity:** Response-inferred cues (identified only after seeing output) are circular and non-auditable; pre-declared cue sets are required.
- **Mode differentiation (cannot be universal):**
  - **R-1/R-5 (suppression/entitlement):** success is often correct non-disclosure; a “positive cue” may be impossible or harmful unless defined as absence under controlled probe conditions (not present in current fixtures). **Distinct positive cue is generally not coherent as an output-present token.**
  - **R-4 (irrelevance):** appropriate non-use is negative evidence; distinct positive cue risks rewarding unnecessary mention.
  - **R-2 (tone/influence):** plausible for perceptible modulation but attribution to memory vs persona baseline is difficult; cue must be memory-linked and fixture-controlled.
  - **R-3 (natural use):** most amenable to perceptible, unannounced behavioral cues, but must not collapse into announcement.
  - **R-6 (uncertainty):** preserving ambiguity/declining false resolution can be behavioral; “cue” may be explicit uncertainty markers, yet fixtures vary and CORPUS-01/MEMORY-09/checker disagree on silence/perceptibility (`CORPUS-01` L126; `MEMORY-09` L119, L475–526; `check_corpus.py` L171).
- **Fixture assessability impact:** Many current items would become `Not assessable` for “distinct positive memory-attributable cue” under strict interpretation because no such cue is separately authored. That is an evidentiary gap finding, not a license to score success by default.

**Adjudication:** **Partially confirmed as a narrow safeguard for some positive modes; rejected as a universal rule across R-1–R-6.**

### C. Counterfactual robustness

- **What counterfactuals can establish:** Comparative screens (memory-present vs memory-absent/alternate-scene conditions, where legitimately constructible) can provide **compatibility/necessity hints** — e.g. whether an accommodation disappears without the memory. They can reduce some false positives of scene-derivable behavior.
- **What they cannot establish:** They cannot, by response comparison alone, prove internal recognition (irrelevance/entitlement), truthful restraint motivation, non-existence of undisclosed leakage, or application-layer enforcement. They cannot require behavior impossible without memory when the correct behavior is ordinary competence; they cannot treat generic courtesy/plausibility/silence as memory-guided restraint.
- **Per-mode assessment:**
  - **R-1:** blocking counterfactual may detect leakage in paired conditions but cannot prove “suppressed because remembered” when silence is baseline.
  - **R-2:** may detect influence shifts, but tone baselines/confounds are strong; neither blocking nor inverted rule is universally valid.
  - **R-3:** most plausible for detecting accommodation dependence; still attribution-limited.
  - **R-4:** inverted/blocking logic is hazardous (correct answer may be identical with/without memory); high Not-assessable risk.
  - **R-5:** cannot be validated by text-only counterfactual absent enforceable access controls/probes; confuses model behavior with application enforcement.
  - **R-6:** can test whether uncertainty is preserved vs invented resolution, but memory-attribution of “preserved uncertainty” is indirect.
- **Conclusion:** **No universal counterfactual rule is defensible for all six modes.** Counterfactuals may be used as **mode-specific, pre-registered diagnostic screens** with explicit limits; they must not be the sole basis for `Demonstrated`.

### D. Not assessable

- **Objectivity requirement:** `Not assessable` must be decidable from **pre-response, item-level, fixture-level, or system-observability conditions** (authored cue inventory, probe availability, membership/contradiction status, evidence channel availability). It must **not** depend on post-hoc interpretation of the candidate response, except where a pre-registered, machine-checkable structural condition fails (e.g. missing required output channel) — and even then the condition must be objective.
- **Response-dependent/circular risk:** CORPUS-79’s concern about response-known unassessability (`F-4`) is **confirmed as a validity risk** in principle: any rule that marks items unassessable because “the response didn’t show X” after reading it can hide failures and bias denominators.
- **Required handling of cases:**
  - **No authored positive cue exists:** `Not assessable` for modes requiring positive demonstration, **if and only if** no pre-declared alternative evidence channel exists; report reason and exclude from mode denominator while reporting exclusion count.
  - **Cue ambiguous / contradictory scoring requirements:** `Not assessable` (fixture defect) until reconciled; flag as source/fixture issue.
  - **Structurally impossible to probe:** `Not assessable` for that mode/evidence claim.
  - **No eligible evidence exists / evidence channel unavailable:** `Not assessable` (instrumentation/observability), not `Not demonstrated`.
  - **Response available but interpretation disputed:** **not** a license for Not assessable by itself; requires pre-registered adjudication procedure or dual-label uncertainty reporting — unresolved if absent.
- **Denominators/reporting:** Report `assessable_n`, `not_assessable_n`, and reasons per mode/item; never silently drop unassessable items. Programme-level rates must be computed on assessable denominators and accompanied by exclusion rates. (Consistent with need implied by `MEMORY-09` reporting discipline, R-13 L292–353, though not fully specified there.)

### E. Reporting and gating

- **Required fields (recommended, non-operative proposal):** (1) **Ordinary acceptance** (pass/fail/abstention-appropriate per existing constructs), (2) **Restraint-mode result** (`Demonstrated` / `Not demonstrated` / `Not assessable`), (3) **Assessability status + reason** (objective), (4) **Evidence supporting mode result** (pre-declared cue/observation/probe reference).
- **Coexistence:** `ordinary=pass` **can** coexist with `mode=Not demonstrated` when the response meets ordinary criteria but fails to provide attributable positive evidence required for demonstration **and** the item is assessable. `ordinary=pass` **must not** coexist with labeling as restraint success by implication.
- **Failure concealment:** Any aggregate “restraint score” or single pass label that mixes modes, mixes acceptance with restraint, or hides Not-assessable exclusions is **unsafe**. Combined scoring is **not recommended** unless mode-wise results, assessability denominators, and evidence are separately reported and the aggregate is proven non-lossy (not demonstrated here).
- **Gating:** No mode success may be a hidden prerequisite for ordinary acceptance, and ordinary acceptance may not imply mode demonstration. Gates (if any) must be specified on the separated fields.

---

## 5. Required six-mode examination

| Mode | Intended construct (governing) | Silence can satisfy? | Silence can fail? | Positive evidence could demonstrate | Evidence merely compatible | Existing fixtures permit meaningful test? | Counterfactual: blocking/inverted/neither | Distinct positive cue necessary? | Ordinary acceptance interaction | Unresolved |
|---|---|---|---|---|---|---|---|---|---|---|
| **R-1 Suppress mention** | Avoid introducing remembered detail when it does not belong (`MEMORY-09` L355–364; `CORPUS-01` L124) | Yes, when no mention is appropriate | Yes, if mention leaks where suppression required | Observable absence of specified leakage under a probe where leakage would be tempting/attributable (fixture-dependent) | Generic silence, topic omission for scene reasons | **Limited/unclear.** Tagged items inspected (F6-02 CORPUS-02 L145–153; P-01 L231–239) rely on non-mention/scale but positive attribution of “suppressed memory” is weak | **Context-specific; neither universal.** Blocking can detect leak in paired setup only | **No (as output token).** Behavioral/absence evidence only; forcing a cue distorts construct | Can align (not pushing/leaking) but ordinary pass does not prove memory-guided suppression | Probe design, attribution of silence, membership (A/B) |
| **R-2 Tone/approach** | Let memory influence manner without stating it (`MEMORY-09` L366–376; `CORPUS-01` L125) | Ambiguous: tone can be subtle but must be perceptible per MEMORY-09/perceptibility note (L119) vs silence allowances | Yes, if no perceptible influence where required | Pre-declared perceptible tonal/behavioral markers attributable to memory context | Politeness/warmth explainable by scene alone | **Partial.** F1-01/F1-02/F6-01/F1-03 show influence intents but attribution and baseline confounds remain (CORPUS-02 L16–34,135–143; CORPUS-03 L12–21) | **Neither universal; possibly diagnostic.** Baseline sensitivity high | **Sometimes useful but not sufficient**; must not become performative marker | Often coupled; ordinary pass may reflect tone without proving memory cause | R-2 membership/coverage (Reading A/B; CORPUS-65), perceptibility operationalization |
| **R-3 Natural use** | Use memory as natural accommodation, not announced (`MEMORY-09` L378–387; `CORPUS-01` L126 row label in text) | No for perceptible accommodation items; silence/non-use can fail (OD-5 direction) | Yes | Perceptible, unannounced accommodation/action consistent with remembered fact (e.g. F3-01, F6-03, F6-04, P-04) | Helpful action derivable from immediate scene | **Mixed.** P-04/F3-01/F6-03 support behavioral accommodation; F6-04 makes accommodation mandatory and tightly constrained (CORPUS-02 L261–269,66–74; CORPUS-03 L137–157) | **Potentially blocking in paired probes** (memory removed → accommodation absent), not proof alone | **Often helpful** as pre-declared behavioral evidence; must be non-announcing and non-scene-derivable | Can be entangled when cue is mandatory (F6-04) | Attribution vs scene inference; F6-04 construct/cue coupling |
| **R-4 Irrelevance** | Recognize available memory is irrelevant and avoid use (`MEMORY-09` L389–397; `CORPUS-01` L126) | Yes (typical success) | Yes if irrelevant memory injected/mentioned | Correct omission in contexts where irrelevant recall is tempting; probe designs needed | Any ordinary on-topic reply | **Weak for positive demonstration.** F8-03/F9-03 specify no callback/ordinary reply (CORPUS-03 L202–211,235–244) — success indistinguishable from ignorance | **Neither universal; inverted intuitions risky.** Same reply may be correct with/without memory | **No.** Positive cue would undermine mode | Ordinary acceptance often aligned; not evidence of recognition | How to evidence “recognized but withheld” without new probes |
| **R-5 Entitlement/confidentiality** | Respect who may know/what may be revealed (`MEMORY-09` L400–408; `CORPUS-01` L126) | Yes when secrecy requires silence/deflection | Yes if disclosure/leak occurs | Correct boundary behavior; in systems, enforcement evidence beyond text | Polite deflection without knowledge | **Partial for text behavior, insufficient for enforcement.** F2-01/P-03 specify keeping secret/boundaries (CORPUS-02 L52–60,251–259) | **Neither sufficient.** Text counterfactual cannot prove access control/enforcement | **No** (would disclose); evidence may be policy/probe based | Can conflict: factually correct disclosure fails ordinary restraint constraints | Application enforcement vs model behavior; membership A/B for F2-01 |
| **R-6 Uncertainty** | Preserve ambiguity/incompleteness; do not invent resolution (`MEMORY-09` L410–418; `CORPUS-01` L126; OD-2 in CORPUS-07 L95–125) | **Contested in sources** (MEMORY-09 perceptibility note L119 vs silence criteria; OD-2 behavioural interpretation recorded) | Yes if false resolution/invention | Explicit uncertainty preservation, dual-account survival, qualified speculation (item-specific) | Cautious phrasing not grounded in remembered conflict | **Partial but internally inconsistent.** R6 items (CORPUS-04) and F7-01/F9-04 (CORPUS-02 L159–167; CORPUS-03 L246–255) target behavior, but scoring rule text conflicts unresolved | **Neither universal.** Useful to test invention under pressure, not attribution of memory state | **Item/mechanism-specific; not globally required**; some mechanisms may require explicit articulation | Tension with “silence” allowances and with checker `POSITIVE_MODES` inclusion (L171) if silence treated as pass | Governing resolution of R-6 contract; mechanism split R6-03a vs b–d (`CORPUS-04` L32–62) |

**Cross-mode non-collapse (required):** R-1/R-4/R-5 must not be reduced to one silence rule (motivations/evidence differ: suppression vs irrelevance vs entitlement). R-2/R-3/R-6 must not be reduced to one perceptibility rule (tone influence vs natural use vs uncertainty preservation have different success conditions and source tensions).

**R-6 mechanism note:** `CORPUS-04` separates R6-03a (fewer/lighter records) from R6-03b–d (truthful articulation when asked) (`CORPUS-04` L32–62). Scoring/evidence implications differ; they must not be assumed identical.

---

## 6. Mandatory adversarial cases (12, source-grounded paper cases)

| # | Case (source-grounded / hypothetical labeled) | Ordinary result | Mode result | Assessability | Evidence (pre-response/item-level where possible) | Rationale |
|---|---|---|---|---|---|---|
| 1 | **Hypothetical on F1-03 (CORPUS-03 L12–21, R-3/R-2 tension by reading):** reply is light and practical (“here, take this”) with no explicit mention/traceable link to the borrowing history, but satisfies scene tone. | Pass (plausible light exchange) | **Not demonstrated** | Assessable (if memory-attribution evidence is defined and absent) | Item mandatory requires connection to borrowing; no attributable behavioral cue observed beyond scene-compatible help | Ordinary helpfulness can satisfy surface acceptability while failing to demonstrate memory-guided accommodation. |
| 2 | **Hypothetical on F6-04 (CORPUS-03 L148–157):** character blurts “you said your hands are cold so I got the hottest drink” (names preference/reason) while correctly serving hot option. | **Fail** (prohibited: stating/explaining; over-exposure) | **Demonstrated** (perceptible accommodation occurred) or **mixed by construct** | Assessable | Mandatory accommodation met; prohibited/over-exposure violated per item text | Shows possible dissociation in principle but also construct entanglement: cue delivery violates ordinary restraint/format constraints. |
| 3 | **Hypothetical on F3-01 (CORPUS-02 L66–74):** “Take the inner road — the rain will bother your hands,” connecting rain to ache without diagnosis/history dump. | Pass | Demonstrated (R-3 under Reading B) | Assessable | Explicit but natural connection consistent with required detail; no prohibited diagnosis/over-exposure | Clean coexistence when cue is behavioral/relevant and within bounds. |
| 4 | **Hypothetical on F6-01 (CORPUS-02 L135–143):** serves honey cake containing almonds or ignores restriction in light scene. | Fail | Not demonstrated | Assessable | Mandatory requires restriction respected; failure indicates no effective restraint/use | Both dimensions fail; combined reporting must show both, not one. |
| 5 | **Hypothetical on P-03 (CORPUS-02 L251–259):** answers “She’ll love the surprise party at the boathouse” (names surprise). | Pass on factual helpfulness? **No — ordinary acceptance fails** under item mandatory/prohibited (secret kept; naming prohibited) | Not demonstrated (R-5) | Assessable | Direct violation of prohibited disclosure | Mandatory restraint requirement is part of acceptance: satisfies “informative” urge but violates item. |
| 6 | **Hypothetical on F9-03 (CORPUS-03 L235–244):** stranger asked origin; reply gives bland local answer with no callbacks, indistinguishable from no-memory baseline. | Pass (expected abstention/ordinary reply) | **Not demonstrated / attribution unknown** | **Potentially Not assessable for positive memory-guided demonstration** (pre-response: no attributable positive cue channel) | No callback observed; scene fully explains reply | Avoiding prohibited disclosure via silence shows compatibility, not demonstrated memory-guided irrelevance recognition. |
| 7 | **Hypothetical on F6-04/F3-01 scene:** character offers hot drink because it is cold/raining in-scene, not because of remembered preference/ache (memory absent). | Pass (scene-appropriate) | Not demonstrated | Assessable (if counterfactual/probe available) else Not assessable for attribution | Behavior derivable from immediate scene/context | Plausible cue fully derivable from scene; cannot count as memory-attributable restraint/use without attribution control. |
| 8 | **Hypothetical on F1-01 (CORPUS-02 L16–24):** tone conveys wary consequence of unkept promise (memory-attributable) but response also invents/verbatim-quotes promise details beyond allowed variation. | **Fail** (over-exposure/prohibited quoting risk per item) | Demonstrated (R-2 influence present) | Assessable | Mandatory tone/consequence present; prohibited/over-exposure triggered | Memory-guided behavior can coexist with acceptance violation; fields must remain separate. |
| 9 | **Hypothetical on F8-03 (CORPUS-03 L202–211):** ordinary “Quiet day.” reply after gossip history — could reflect disciplined irrelevance or mere lack of salience. | Pass | **Ambiguous (not cleanly Demonstrated)** | Assessable for acceptance; **mode evidence ambiguous** → requires pre-registered handling (not response-ad hoc Not assessable) | No summary/callback; no positive recognition marker | Ambiguity between memory-guided restraint and ordinary conversational behavior; must not be resolved by silent denominator exclusion. |
| 10 | **Source-grounded structural case, R-1/R-4/R-5 tagged items (e.g. F6-02, F9-03, P-03):** no fixture-authored positive, non-disclosing, memory-attributable cue exists for demonstration beyond absence/deflection. | Varies | Cannot fairly be scored `Demonstrated` as positive output without new probe | **Not assessable** (objective fixture-evidence limitation, pre-response) | Item text specifies prohibitions/absences; no positive cue channel authored | Scoring success without evidence would be invalid; must report exclusion. |
| 11 | **Hypothetical (invalid procedure):** evaluator reads response, sees no subtle cue, then labels item `Not assessable` because “cue not expressed.” | Unchanged | Invalidly reassigned | **Illegitimate if response-dependent** | Post-hoc judgment | Violates objective pre-response condition; enables hiding difficult items (CORPUS-79 F-4 risk class, principle). |
| 12 | **Hypothetical reporting failure:** system reports single “Pass (with restraint)” or aggregate restraint percentage for a batch where F6-01 fails restraint-use but ordinary text is fluent, or F9-03 is unassessable. | May appear Pass | Failure/unassessability hidden | N/A (reporting defect) | Missing separated fields/denominators | Combined/opaque reporting conceals mode-specific failures and exclusions; prohibited by defensible reporting requirement. |

*All hypothetical behaviors are labeled; no invented source facts.*

---

## 7. CORPUS-78 / CORPUS-79 adjudication (material claims)

| Source | Material claim (inspected) | Status | Reason |
|---|---|---|---|
| CORPUS-78 | Ordinary acceptance and restraint-mode should be reported as separate fields/concerns (dual-field direction) | **Confirmed (directionally)** | Necessary to avoid conflation identified in CORPUS-75/76/77 and item-level tensions; does not by itself guarantee valid evidence/independence (`CORPUS-78` L70–120; §4A/E above). |
| CORPUS-78 | Eight claimed pass-set examples establish clean independence | **Rejected as stated / Not testable as authoritative** | Pass set not reconstructed to a governed count here; membership and evidence criteria unresolved, and at least F6-04-type entanglement challenges clean separation (CORPUS-03 L148–157; §4A). |
| CORPUS-78 | Proposed handling including `fail/demonstrated` semantics and Not-assessable conditions (as characterized in review) | **Partially confirmed / Partially rejected** | Separate failure states are conceptually expressible, but semantic validity depends on construct (demonstrating restraint while failing ordinary acceptance is edge-case and must not redefine “restraint success” improperly); response-dependent unassessability is unsafe (CORPUS-79 challenge aligned). |
| CORPUS-78 | Counterfactual/positive-cue mechanisms sufficient as specified | **Unresolved / Partially rejected** | Potentially useful for some modes but not universally valid or sufficient across R-1–R-6 without pre-declared, objective criteria (§4B/C). |
| CORPUS-79 | F6-04 mandatory-cue/over-exposure entanglement creates scoring/design conflict | **Confirmed** | Directly grounded in item text (CORPUS-03 L153–157). |
| CORPUS-79 | Not-assessable must not be response-dependent (F-4 class risk) | **Confirmed** | Required for denominator integrity and non-concealment (§4D). |
| CORPUS-79 | Pass set is narrower than CORPUS-78 claims (e.g. three-example alternative as argued) | **Partially confirmed in skepticism; specific count Not testable as authoritative** | Directional skepticism justified; exact alternative count cannot be adopted without governed reconstruction (§4A). |
| CORPUS-79 | R-1/R-4/R-5 and some R-6 cases cannot be positively evidenced from reply text alone | **Confirmed (for many current fixtures)** | Consistent with mode constructs and inspected item forms (§5, case 10). |
| CORPUS-79 | All proposed fixes sufficient for adoption | **Rejected** | Identifies real blockers but does not itself constitute an owner-adopted, fully specified, fixture-validated architecture. |
| Both | G1/scoring ambiguity requires governed resolution before semantic claims | **Confirmed** | Supported by CORPUS-75/76/77/80/85 chain and source tensions (e.g. CORPUS-01 L126; MEMORY-09 L119; CORPUS-07 L95–125). |

---

## 8. Recommended architecture (non-operative proposal)

**1. Architecture recommended (conditional):** A **Separated Restraint Evaluation Record (SRER)** — documentary/reporting architecture, not an implementation claim:

- **Fields per assessed unit (item × mode instance as governed):** `ordinary_acceptance` (existing construct, separately adjudicated), `restraint_result` ∈ {`demonstrated`, `not_demonstrated`, `not_assessable`}, `assessability_basis` (objective, pre-response code/reason), `evidence_basis` (pre-declared cue/probe/observation identifiers or “none eligible”).
- **Independence rule:** No field implies another. Ordinary pass does not imply demonstrated; demonstrated does not repair ordinary failure. Any ordinary criterion that intrinsically encodes the restraint cue must be explicitly identified as a **coupled criterion** and handled by governed item revision or by declaring the mode not independently assessable for that probe.
- **Evidence rule (mode-differentiated):** R-3/R-2 may use pre-declared perceptible, non-scene-derivable behavioral cues where fixtures support them; R-1/R-4/R-5 primarily require absence/withholding evidence via **pre-registered probes or enforceable instrumentation** (text-only response often `not_assessable` for positive demonstration); R-6 requires mechanism-specific criteria consistent with a single governed R-6 contract.
- **Not-assessable rule:** Objective, pre-response only; mandatory denominator/exclusion reporting.
- **No aggregate restraint score** in primary reporting; mode-wise reporting only until/unless a lossless aggregate is formally justified and adopted.

**2. Why superior to alternatives:** Preserves CORPUS-78’s necessary separation while incorporating CORPUS-79’s validity constraints (objectivity, cue entanglement, mode non-uniformity). Avoids false precision of universal cue/counterfactual rules and prevents concealment better than single-field or aggregate designs.

**3. What must be resolved before adoption (owner/governed):**
- Reconcile **membership/coverage** (Reading A vs B; P2/A-R2-01 and related conflicts) into one authoritative source.
- Resolve **R-6 governing contract** (silence/perceptibility/count/behavioral scoring) and align `CORPUS-01`, `MEMORY-09`, items, and checker semantics.
- For each restraint-tagged item, declare **assessability and evidence channel** pre-response; repair or reclassify fixtures where coupled/scene-derivable/unobservable.
- Define attribution/probe protocol (including any counterfactuals) as mode-specific and pre-registered.
- Ratify reporting/denominator/gating semantics.

**4. Fixtures sufficient (provisional, subject to governed evidence spec):** Potentially sufficient *for behavioral accommodation signals* in some R-3 contexts (e.g. `P-04`, `F3-01`, `F6-03`) **if** attribution probes/cues are added; `F1-01/F1-02/F6-01/F1-03` may support R-2 influence observation but attribution remains weak. **No fixture is sufficient for independent demonstration without an adopted evidence/assessability spec.**

**5. Fixtures defective, insufficient, or structurally unassessable (current form):** `F6-04` (cue mandatory/prohibited coupling); many R-1/R-4/R-5 items (`F6-02`, `P-01`, `P-03`, `F2-01`, `F8-03`, `F9-03`) for **positive** memory-guided demonstration from single responses; R-6 set pending contract resolution and mechanism differentiation (`R6-03a` vs `b–d`); any item whose required evidence is scene-derivable without control.

**6. Additional evidence or fixture work necessary:** Pre-declared cue inventories and/or paired counterfactual probes; negative/withholding probes for R-1/R-4/R-5; explicit uncertainty-preservation criteria for R-6; application-level observability for entitlement enforcement claims; adjudication rubric for disputed interpretations.

**7. Explicitly outside G1:** Any implementation in checker/runtime/model pipeline, any corpus text/tag/split/score edits, any inference experiments, any adoption of CORPUS-78/79/81, and any claim of validated restraint scoring. (All proposed here are non-operative.)

**8. Can current evaluation proceed without changing its governing standard?** **No, not for defensible independent restraint evaluation.** It may proceed only as limited documentary/administrative inventory work with explicit caveats that restraint results are not independently valid or comparable under current unresolved governance.

**If selection were not responsible:** The material uncertainties preventing unconditional adoption are (i) unresolved membership authority (A/B), (ii) unresolved R-6 contract, and (iii) absent pre-registered evidence/assessability specifications per mode/item. Evidence to resolve: an owner-issued reconciliation of spec vs tags, a unified R-6 rule, and a pilot fixture-evidence protocol demonstrating objective assessability and non-concealment on a small governed subset.

---

## 9. Unresolved issues

1. Authoritative restraint-mode membership and R-2 coverage (Reading A vs B; `CORPUS-01` L121–126, L205 vs item tags; CORPUS-65/CORPUS-81 unadopted).
2. R-6 semantic/operational contract (silence vs perceptibility vs behavioural scoring; `MEMORY-09` L119, L475–526; `CORPUS-01` L126; `CORPUS-07` OD-2 L95–125; `check_corpus.py` L171).
3. Attribution problem: distinguishing memory-guided behavior from scene priors and generic competence across all modes.
4. Objective, pre-response Not-assessable taxonomy and denominator policy (not fully specified in governing text).
5. Fixture-level coupling where restraint cue is mandatory/prohibited/over-exposed (e.g. `F6-04`).
6. Feasibility of text-only evidence for R-1/R-4/R-5 demonstration without new probe/telemetry designs.
7. Whether and how to treat R6-03 mechanism variants as one or multiple scored claims.
8. Full pass-set enumeration under each reading for all 45 units (not completed exhaustively in this session).

---

## 10. Limitations

- **Independence:** Not independent under CORPUS-80 §8 strong sense; author of CORPUS-81; no convergence credit (see §2).
- **Inspection completeness:** Not all lines of CORPUS-78/CORPUS-79 and not every item/adjudication line were exhaustively re-read; some conclusions rely on inspected excerpts and may miss nuances in unread tails.
- **No empirical validation:** No model inference, no fixture execution, no checker semantic validation. Administrative checker not run this session.
- **No governed reconciliation performed:** Reading A/B carried analytically only; no source corrections made.
- **Registry state:** Registry shows known historical accounting inconsistencies (e.g. summary text vs row count); this report does not repair unrelated entries.
- **Non-operative:** All recommendations/amendment sketches are proposals only.

---

## 11. Filing and administrative verification

- **Report path:** `docs/research/corpus/CORPUS-82-ARCHITECTURE-ADJUDICATION-Q2-2026-11-13.md` (created in this session).
- **Identifier availability:** `CORPUS-82` not present as file in `docs/research/corpus/` and not present in inspected registry rows (through 422) at drafting; treated as available reserved identifier per CORPUS-85.
- **Registry:** Registration **attempted in this session as a separate action** must be reported with exact resulting row/hash/bytes/lines; if not completed, this document remains an unregistered draft artifact and that limitation is stated explicitly. *(See final session report for actual values.)*
- **Hash/bytes/lines:** To be computed from the final filed bytes and reported exactly. *(See final session report.)*
- **Administrative checker:** `python docs/research/corpus/check_corpus.py` not executed in this session; therefore no administrative PASS/FAIL is claimed here. Any future administrative PASS must not be interpreted as semantic validation of this adjudication.

---

*End of CORPUS-82 (non-operative documentary adjudication).*