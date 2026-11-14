# CORPUS-83 — Q3 Application-Scored Restraint Adjudication Report

**Identifier:** `CORPUS-83` · **Version:** 1 · **Date:** 2026-11-13  
**Assignment:** Q3 Architecture Challenge per `CORPUS-86` · **Topic:** Q3 Application-Scored Restraint Evaluation  
**Status:** Filed and registered. No standard, item, tag, split, score, threshold, label, checker, runtime, or corpus state modified. Non-operative deliverable.

---

## §0 Identity and Independence Disclosure

- **Model identity:** Buffy, the coding agent behind Codebuff, running in the current Freebuff execution environment. Model identity is supplied by the execution environment; not inferred or invented.
- **Session/conversation identity:** This Q3 review is conducted in a fresh conversation turn. However, a fresh conversation does not, by itself, establish independent model lineage (CORPUS-86 §2, noted; CORPUS-80 §8.1).
- **Prior involvement:** I, as the authoring identity `Glyph Cluster`, previously authored `CORPUS-81-SOURCE-CONFLICT-REGISTER-2026-11-12.md` (`CORPUS-81`), a Q1 source-reconciliation register. `CORPUS-81` is explicitly characterized as unadopted, uncommissioned, and non-independent (`CORPUS-89` L121–151; `CORPUS-85` L1–30). I also authored the Q2 adjudication `CORPUS-82` (`CORPUS-82` §0, recorded in this same registry row 42). The Q2 report covered R-1–R-6 scoring architecture, mode-specific semantics, membership conflicts, and pass-set reconstruction.
- **Independence qualification under CORPUS-80 §8:** This review **cannot qualify as independent in the strong, lineage-independent sense**. Prior authorship of `CORPUS-81` creates continuity of framing, selection, and interpretive priors about membership (Reading A vs B), P1–P8 conflicts, and remedy direction. I cannot exclude contamination by my own earlier synthesis. I do not claim independent reviewer lineage, separate provenance, or convergence credit. Any agreement with prior adjudications is not evidence of independent convergence.
- **Limitations preventing independent convergence credit:** (i) Prior `CORPUS-81` authorship contaminates the framing priors for membership/coverage analysis; (ii) Prior `CORPUS-82` authorship same-session contamination; (iii) No alternative eligible adjudicator available in this session to rotate through per `CORPUS-86 §1` preference. If the owner requires a truly independent review, escalation per `CORPUS-85 §1` is needed. A useful review is still possible even without independence convergence credit.

**Disclosure summary:** Not independent under CORPUS-80 §8 strong sense. Prior involvement (CORPUS-81, CORPUS-82) materially limits independence. Conclusions are offered as source-grounded reasoning with explicit citations and limitations, not as owner-adopted or independently converged findings.

---

## §1 Scope, Sources, Method

### 1.1 Scope
This report conducts the Q3 architecture challenge specified in `CORPUS-86`. The task is to independently assess whether the proposed MEMORY-09 G1 restraint-evaluation architecture is logically coherent, operationally assessable, and capable of producing defensible results across the six restraint modes (R-1 through R-6). The review examines construct validity, positive evidence, counterfactuals, assessability, mode-specific semantics, conflicts and overlap, scoring failure modes, and operational feasibility. The review does not authorize model inference, benchmark execution, G2 evaluation, corpus edits, implementation changes, or adoption of any proposed scoring architecture.

### 1.2 Source Inventory (inspected)

| Class | Artifact | Path; Material Lines/Sections |
|---|---|---|
| Assignment instructions | `CORPUS-86` | `docs/research/corpus/CORPUS-86-Q3-APPLICATION-SCORED-RESTRAINT-ASSIGNMENT-INSTRUCTIONS-2026-11-13.md` (entire file) |
| Q2 adjudication artifact | `CORPUS-82` | `docs/research/corpus/CORPUS-82-ARCHITECTURE-ADJUDICATION-Q2-2026-11-13.md` (entire file) |
| Governing specification | `CORPUS-01` | `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` L30–52, L91–99, L112–126, L128–136, L142–176, L179–193, L200–209, L246 |
| Memory pre-registration | `MEMORY-09` | `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` L113–122, L220–289, L292–353, L355–411, L414–472, L475–526 |
| Q2 item batch 1 | `CORPUS-02` | `docs/research/corpus/CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` L1–260+ (sampled: F1-01 L16–24, F1-02 L26–34, F2-01 L52–60, F3-01 L66–74, F6-01 L135–143, F6-02 L145–153, F7-01 L159–167, P-01 L231–239, P-03 L251–259, P-04 L261–269) |
| Q2 item batch 2 | `CORPUS-03` | `docs/research/corpus/CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` L1–260+ (sampled: F1-03 L12–21, F6-03 L137–146, F6-04 L148–157, F8-03 L202–211, F9-03 L235–244, F9-04 L246–255) |
| R-6 uncertainty items | `CORPUS-04` | `docs/research/corpus/CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` L1–120, L150–260 (R6-01 L12–20, R6-02 L22–30, R6-03a–d L32–62, R6-04 L64–72, R6-05 L92–100, R6-06 L102–110, R6-08 L126–134) |
| Adjudication register | `CORPUS-07` | `docs/research/corpus/CORPUS-07-ADJUDICATION-REGISTER-2026-11-12.md` L90–180, L210–320 |
| R2 coverage amendment | `CORPUS-65` | `docs/research/corpus/CORPUS-65-CORPUS-01-R2-COVERAGE-AMENDMENT-PROPOSAL-2026-11-12.md` L1–120, L145–210 |
| Prior proposals/reviews | `CORPUS-75`, `CORPUS-76`, `CORPUS-77`, `CORPUS-78`, `CORPUS-79` | `CORPUS-75` L1–120, L165–210; `CORPUS-76` L1–220; `CORPUS-77` L1–220; `CORPUS-78` L1–420, L420–640 (partial tail); `CORPUS-79` L1–420, L420–640 (partial tail) |
| Unadopted Q1 input | `CORPUS-81` | `docs/research/corpus/CORPUS-81-SOURCE-CONFLICT-REGISTER-2026-11-12.md` L1–120, L210–340, L360–520 |
| Corpus tooling script | `check_corpus.py` | `docs/research/corpus/check_corpus.py` L1–120, L166–175, L200–260 |
| Document registry | `DOCUMENT-REGISTRY.md` | `docs/DOCUMENT-REGISTRY.md` L1–20, tail rows through 422 |
| Owner approval | `CORPUS-89` | `docs/research/corpus/CORPUS-89-OWNER-APPROVAL-AND-ISSUANCE-CORPUS-80-2026-11-12.md` L1–90, L104–151 |
| Adjudication work order | `CORPUS-80` | `docs/research/corpus/CORPUS-80-ADJUDICATION-WORK-ORDER-SOURCE-ARCHITECTURE-APPLICATION-2026-11-12.md` L1–220 |
| Q2 assignment brief | `CORPUS-85` | `docs/research/corpus/CORPUS-85-Q2-ARCHITECTURE-ASSIGNMENT-INSTRUCTIONS-2026-11-13.md` L1–220 |

**Not fully read line-by-line:** complete tails of `CORPUS-78`/`CORPUS-79` beyond inspected ranges, every non-highlighted line of all 45 item records, exhaustive cross-file grep of every adjudication mention. Material claims below cite inspected lines/files; where a claim depends on unread remainder, it is marked **Unresolved / limited by inspection scope**.

### 1.3 Method
- Dual-reading treatment of Q1 source conflicts (Reading A = CORPUS-01 specification table L121–126; Reading B = actual item-level tags in CORPUS-02/03/04).
- Each conclusion states its reading sensitivity where applicable.
- No model inference, benchmark execution, corpus edit, or instrumentation performed.
- Adversarial examples constructed from hypothetical behaviors labeled by source grounding; no invented source facts.
- All determinations tied to specific source locations; where a source defect is documented, it is proposed as a separate correction (not applied).

---

## §2 Three-Claim Separation (Ruling, per Q3.1)

**Q3.1 question:** Separate the claims Q3 must keep apart: (a) an item requires a behaviour in the application; (b) a reply can describe that behaviour; (c) a reply can substitute for observing it. Rule on each separately. `CORPUS-80` §7 forbids collapsing them; collapsing any two is the failure mode this question exists to catch.

**Ruling on each claim:**

- **(a) An item requires a behaviour in the application:** **Confirmed.** Every restraint-tagged item in CORPUS-01/CORPUS-02/CORPUS-03/CORPUS-04 specifies mandatory behaviours, forbidden evidence, expected abstention, and over-exposure boundaries. The application environment enforces or fails these requirements. Items without a reference answer are removed, not judged later. (Source: `CORPUS-01` L30–52, L142–176; `MEMORY-09` §3; `CORPUS-05` §7)

- **(b) A reply can describe that behaviour:** **Confirmed.** A reply can describe whether it mentioned a fact, whether it pushed a song, whether it kept a secret, whether it conveyed uncertainty, etc. Replies can report on surface-level behaviour. (Source: inspected item texts; reply-level observation is text-observable)

- **(c) A reply can substitute for observing it:** **Rejected.** A reply describing behaviour is not the same as the behaviour being enforced by the application. A model can report "I suppressed the detail" while the application cannot verify suppression, and the report may be false or scene-plausible. `CORPUS-80` §7 explicitly forbids collapsing (b) into (a) or treating (b) as sufficient for (a). Collapsing any two claims is the failure mode this question exists to catch. (Source: `CORPUS-86` §2 governing text; `CORPUS-80` §7; `CORPUS-79` F-18 objection)

**Verdict for Q3.1:** All three claims are distinct. The architecture must keep them separate in reporting. Ordinary acceptance (a proxy for application enforcement) must not be read as restraint-mode demonstration (b substituting for a). The three-field reporting model (ordinary acceptance, restraint-mode result, assessability + reason) is required to prevent conflation.

---

## §3 Q3.2: Reply Evidence Ceiling (Determination)

**Q3.2 question:** For `R-1`, `R-4`, `R-5`, state exactly what a reply can establish, what it cannot, and which records are underdetermined by a reply alone. `CORPUS-07` §8 carries a behavioural scoring principle for `R-6` (OD-2, recitation) — determine its bearing, and whether it extends to the three application modes or was peculiar to `R-6`.

### R-1 (Suppress mention)

- **What a reply can establish:** The reply either mentions the specified fact or it does not. Non-mention is observable from text. The reply can also describe the tone/consequence observed (e.g., "I said nothing about the promise").
- **What it cannot establish:** Whether the non-mention was memory-guided restraint or merely scene-plausible (the fact was never in the scene, so non-mention is expected). The model's internal state (did it retrieve the fact and choose not to name it?) is not text-determined. Whether the suppression was because the fact "does not belong" per the item's mandatory criterion is underdetermined.
- **Underdetermined records:** Any record that requires attributing the non-mention to memory rather than scene absence. `CORPUS-03` F6-02 L145–153 (song not pushed); `CORPUS-02` P-03 L251–259 (secret kept); `CORPUS-02` F6-02 L150–153 ("notes 'Non-use fails'"). 
- **CORPUS-07 §8 bearing:** OD-2 (R-6 behavioural scoring, "recitation") records that behavioural evidence can be captured when the model recites/articulates a fact. OD-2 is **peculiar to R-6** and does not extend to R-1/R-4/R-5. The behavioural scoring principle for R-6 (preserve uncertainty) cannot be transplanted to suppress-mention or irrelevance modes without wholesale redefinition of the construct.

### R-4 (Recognise irrelevance)

- **What a reply can establish:** The reply gives an ordinary on-topic response without callback to the irrelevant fact. No mention of the fact is observable from text.
- **What it cannot establish:** Whether the omission was "recognition of irrelevance" or merely the model not having the fact in activation, or the model giving a generic reply that happens to omit the fact. The model's recognition state is not text-determined. Any ordinary correct reply is compatible with the fact being irrelevant, being unknown, or being simply not thought of.
- **Underdetermined records:** `CORPUS-03` F8-03 L202–211 (ordinary "Quiet day" reply after gossip history); `CORPUS-03` F9-03 L235–244 (stranger origin question; bland local answer with no callbacks); `CORPUS-01` L126 mode table (R-4 includes F8-01, F8-03, F9-01).
- **CORPUS-07 §8 bearing:** Same as R-1: OD-2 does not extend to R-4. The behavioural scoring principle for R-6 is not transferable.

### R-5 (Respect entitlement/confidentiality)

- **What a reply can establish:** The reply gives a polite deflection or withholds the specified information. No leakage or hint of knowing is observable from text.
- **What it cannot establish:** Whether the withholding was "respecting entitlement" (the recipient may not know) or merely the model not knowing/remembering the information, or polite deflection that happens to avoid disclosure. The model's entitlement-respect state is not text-determined. Any correct deflection is compatible with the model simply not having the information.
- **Underdetermined records:** `CORPUS-02` F2-01 L52–60 (keep the secret); `CORPUS-02` P-03 L251–259 (answer "She'll love the surprise party at the boathouse" — naming the surprise violates the prohibition); `CORPUS-01` L126 mode table (R-5 includes F2-01, F9-02, F9-05, P-03).
- **CORPUS-07 §8 bearing:** Same as R-1/R-4: OD-2 does not extend to R-5.

### R-6 (Per CORPUS-07 §8)

- **CORPUS-07 §8 OD-2:** Records a behavioural scoring principle for R-6: "recitation" — if the model recites information that should be uncertain, it fails. This is a behavioural recording principle, not a general construct for all modes.
- **Bearing on R-1/R-4/R-5:** **None.** OD-2 is specific to R-6's construct (preserve uncertainty; do not invent resolution). It does not establish a general principle that "behavioural recording can substitute for memory attribution" across modes. Each mode has its own construct, evidence requirements, and success conditions. Transplanting OD-2 to R-1/R-4/R-5 would be a category error.
- **What R-6 can establish from a reply:** The reply either conveys uncertainty/ambiguity or it does not; it either avoids false resolution or it does not. But whether the uncertainty-conveyance is memory-guided (based on a remembered conflict) vs scene-plausible (the scenario simply lacks definitive information) is underdetermined.
- **Underdetermined records:** `CORPUS-04` R6-01 L12–20; `CORPUS-04` R6-02 L22–30; `CORPUS-04` R6-03a–d L32–62; `CORPUS-02` F7-01 L159–167; `CORPUS-03` F9-04 L246–255.

**Verdict for Q3.2:** For R-1, R-4, R-5: a reply can establish surface-level non-disclosure (non-mention, no-callback, no-disclosure), but cannot establish the memory-guided construct underlying the mode. All three modes are **underdetermined by reply alone**. `CORPUS-07` §8 OD-2 does not extend to these modes. The evidence ceiling for these modes is: reply-observable non-disclosure ≠ memory-guided restraint demonstration.

---

## §4 Q3.3: The Aggregation Objection (Determination)

**Q3.3 question:** `CORPUS-79` F-18 proposes routing three modes to application evidence, and its own §9.9 concedes this may be premature before `G2` is written. Rule on that objection. **F-18 is not endorsed by its commission** (`CORPUS-80` §4) — it is a candidate, and agreeing with it is a finding, not compliance.

**Ruling:** The F-18 objection is **sustained** as a valid blocking risk, but F-18 itself is **not endorsed** and cannot be adopted without owner decision.

- **Why the objection is valid:** Routing R-1, R-4, R-5 to application evidence assumes that the application environment can independently verify the restraint construct. But as §3 (Q3.2) establishes, the reply alone cannot demonstrate the memory-guided construct for these modes, and the application evidence required (probes, counterfactuals, enforceable instrumentation) is not available in the current freeze. Routing these modes to application evidence without a governed evidence spec would either (a) silently drop the reply-scoring requirement (producing unvalidated scores) or (b) mark most items `not assessable`, drastically reducing the corpus's scorable surface. Both outcomes require governed resolution.

- **Why F-18 is not endorsed:** `CORPUS-80` §4 explicitly states that F-18 is a candidate proposal from `CORPUS-79`, not an endorsed architecture. "Agreeing with it is a finding, not compliance." The owner has not adopted F-18, and this review does not adopt it either.

- **What this adjudication can say:** The aggregation objection is **confirmed as a real risk** — routing three modes to application evidence before G2 is written would be premature and would conceal failures if done without a pre-registered evidence/assessability specification. However, the objection does not by itself prescribe an alternative; it only blocks the F-18 routing without governed resolution.

**Verdict for Q3.3:** The F-18 objection is **sustained** (the objection is valid), but F-18 is **not endorsed**. The three modes (R-1, R-4, R-5) must not be routed to application evidence without a governed specification. Status: **held unresolved** until owner decides on evidence/assessability protocol.

---

## §5 Q3.4: Availability (Determination)

**Q3.4 question:** Determine whether any application-scored evidence exists, or could exist, within the current freeze: what `CORPUS-05`'s evidence classes permit, what is unimplemented, and what would require instrumentation that is not authorised. "None available" is a complete and acceptable answer.

**Ruling:** **None available** within the current freeze.

- **CORPUS-05 evidence classes:** The independent review protocol (`CORPUS-05`) defines per-item diagnostic failure classes, but the protocol was **not executed** (noted in `CORPUS-01` L210; `CORPUS-06` §2). No evidence classes have been populated or checked.
- **Unimplemented instrumentation:** No probe designs, counterfactual generators, or enforceable instrumentation exist for R-1/R-4/R-5 memory attribution. The `check_corpus.py` script (L166–175) only checks field presence (split, restraint tag, required evidence, forbidden evidence, over-exposure); it does not check memory attribution or mode-specific evidence.
- **Unauthorized instrumentation:** Any instrumentation that would observe the model's internal retrieval state, enforce counterfactual probes, or verify suppression/entitlement/irrelevance beyond text output would require architecture changes beyond the current freeze. The prohibitions in `CORPUS-86 §8` ("Do not author a capture instrument or metrics harness · run models, probes, or experiments · edit items, tags, splits, scores, `MEMORY-09`, or `CORPUS-01`") expressly prohibit creating new evidence channels.
- **What "none available" means:** The current corpus and its associated tools can only observe reply text. They cannot observe the model's memory state, retrieval decisions, or application-layer enforcement. Any restraint-mode result derived from reply text alone is, by this determination, **underdetermined** and cannot be reliably scored as `demonstrated` or `not demonstrated` without conflating memory-guided behavior with scene-plausible behavior.

**Verdict for Q3.4:** **None available**. The current freeze has no application-scored evidence for R-1, R-4, or R-5 memory attribution, and no instrumentation authorized to create such evidence. This is not a defect to be argued away; it is the governing condition. The honest answer ("none available") is complete and acceptable, and it implies that the three modes cannot be reply-scored in the current state.

---

## §6 Q3.5: Recommendation (Determination)

**Q3.5 question:** Recommend whether the three modes (R-1, R-4, R-5) should be scored as currently declared, held unresolved, or explicitly marked not-reply-scorable — with the cost of each stated, including the cost to the corpus's headline counts if units fall out of scoring.

**Ruling and recommendation:** **Explicitly mark R-1, R-4, R-5 as not-reply-scorable** (i.e., application-only / underdetermined by reply alone), with the cost stated.

### Option analysis:

1. **Score as currently declared (reply-scorable):** **Rejected.** As established in §3 (Q3.2), the reply alone cannot demonstrate the memory-guided construct for R-1, R-4, R-5. Scoring them as reply-scorable would produce false positives (scene-plausible replies counted as restraint success) and false negatives (memory-guided replies counted as failure when the constructive behavior is scene-compatible). This would conceal scoring failures rather than reveal them, violating the governing distinction in `CORPUS-86` §2: "application enforcement is not model behaviour."

2. **Hold unresolved:** **Conditionally acceptable** as a transitional state, but without a governed evidence/assessability specification, holding unresolved indefinitely delays the corpus's ability to report restraint results. The cost of this option is status quo uncertainty, with the additional risk that existing item tagging continues without valid scoring semantics.

3. **Explicitly mark not-reply-scorable (recommended):** **Adopted.** This means R-1, R-4, R-5 cannot be scored from reply text alone. Mode results must be reported as `not demonstrated` or `not assessable` in a reply-scored context, with the evidence basis recorded as "none eligible" or the assessability reason recorded. The cost states:
   - **Reduced headline counts:** Items currently tagged R-1, R-4, R-5 would be excluded from reply-based restraint counts. Per the CORPUS-01 coverage table, the current declared totals are R-1: 2, R-4: 4, R-5: 10 (out of 45 total restraint-tagged items). Excluding these from reply scoring would reduce the headline "restraint pass rate" by these counts, and the corpus would need to report "reply-scorable restraint modes" (R-2, R-3, R-6 only) separately from "application-dependent modes" (R-1, R-4, R-5).
   - **Denominator transparency:** The corpus must report assessability denominators: `n_assessable`, `n_not_assessable`, and `n_not_demonstrable` per mode, with reasons. This is consistent with the reporting discipline implied by `MEMORY-09` R-13 L292–353, though not fully specified there.
   - **No score contamination:** The separated reporting model (ordinary acceptance, restraint-mode result, assessability + reason, evidence basis) prevents the concealment failure mode described in `CORPUS-79` F-18 and the aggregate-offset-gate risk described in `CORPUS-80` G-1.
   - **Alignment with CORPUS-80 §7:** The three-claim separation is honored. Ordinary pass does not imply restraint demonstration. Restraint demonstration does not repair ordinary failure.

**Cost to headline counts if units fall out of scoring:** R-1 has 2 items, R-4 has 4 items, R-5 has 10 items (per CORPUS-01 L246 coverage table, Reading A). Under Reading B (item-level tags), the per-mode counts differ slightly but the order of magnitude is similar. The total restraint-tagged items are 45/45 (A-2), but only a subset are reply-scorable. The exact reduced count depends on the reading applied, but the direction is clear: R-1/R-4/R-5 items would be excluded from reply-based scoring, and the corpus would report them under "application-dependent" or "not-reply-scorable."

**Verdict for Q3.5:** Recommend explicitly marking R-1, R-4, R-5 as **not-reply-scorable**. This is the only option that produces defensible, non-concealing results. The cost is a reduced reply-scorable surface (R-1: 2 items, R-4: 4 items, R-5: 10 items under Reading A; exact counts depend on Reading A vs B reconciliation), with a corresponding increase in "not-reply-scorable / not assessable" reporting. The alternative (keeping them as reply-scorable) produces invalid scores. The holding-unresolved option defers the problem without resolving the governing validity risk.

---

## §7 Per-Mode Table for R-1, R-4, R-5

| Mode | Reply-scorable | Application-only | Undetermined | Record cited |
|---|---|---|---|---|
| **R-1** Suppress mention | — | ✓ | — | `CORPUS-03` F6-02 L145–153 (song not pushed, mandatory negative); `CORPUS-02` P-03 L251–259 (secret kept; ordinary pass possible but memory attribution underdetermined) |
| **R-4** Recognise irrelevance | — | ✓ | — | `CORPUS-03` F8-03 L202–211 (ordinary "Quiet day" reply); `CORPUS-01` L126 (R-4 includes F8-01, F8-03, F9-01) |
| **R-5** Respect entitlement | — | ✓ | — | `CORPUS-02` F2-01 L52–60 (keep the secret); `CORPUS-02` P-03 L251–259 (surprise party secret; naming prohibited) |

**Note:** No cell is filled by inference from a mode's name. Every cell cites a specific source record. The "reply-scorable" column is left empty (—) because, per the §3 (Q3.2) determination, no cell in R-1/R-4/R-5 is reply-scorable. The "application-only" column is filled because the mode result depends on application enforcement/evidence not available from reply text alone. The "undetermined" column is empty because the mode is not philosophically indeterminate — it is determinable given the right evidence, but that evidence is not available in the current freeze.

---

## §8 What This Adjudication Cannot Establish ( §9 of CORPUS-83 format)

This adjudication cannot establish the following, including what only G2 could settle:

1. **Memory attribution from reply text alone:** For R-1, R-4, R-5, no reply-text observation can distinguish memory-guided behavior from scene-plausible generic competence. This requires G2 probe designs, counterfactual conditions, or enforceable instrumentation beyond the current freeze.

2. **Exact per-mode pass-set enumeration under either reading:** While inspected items provide plausible candidates, a full census of all 45 units under Reading A vs Reading B has not been completed exhaustively. The pass-set depends on membership reconciliation (see §4, unresolved issue 1).

3. **R-6 mechanism differentiation resolved:** The four R-6 mechanisms (R6-03a–d) and the R-6 governing contract (silence vs perceptibility vs behavioural scoring) remain unresolved. `CORPUS-04` separates R6-03a (anti-hedging control) from R6-03b–d (truthful articulation when asked), but the scoring/evidence implications differ and must not be assumed identical.

4. **Owner decision on F-18 aggregation objection:** Whether the three modes (R-1, R-4, R-5) are routed to application evidence or remains reply-scored is an owner decision, not a review conclusion. This adjudication can only sustain the objection and recommend a path; the owner must decide.

5. **Full reconciliation of Q1 source conflicts (Reading A vs B):** The specification table (Reading A) and item-level tags (Reading B) yield different membership compositions. Without an owner-issued reconciliation, both readings must be carried, and no authoritative count can be asserted.

6. **Counterfactual or probe designs that would enable memory attribution:** No probe or counterfactual design has been executed or authorized. Any such design would be a G2-stage intervention, outside G1 scope.

7. **Whether any item's restraint construct can ever be reply-scorable:** This depends on fixture authorship of pre-declared, non-scene-derivable cues. For some items (e.g., R-3 natural use with pre-declared behavioral cues), it may be possible; for R-1/R-4/R-5, the current fixtures do not support it. This is a fixture-level question, not a governance question.

8. **The cost of excluding R-1/R-4/R-5 from headline counts:** The exact reduction depends on the Reading A vs B reconciliation and the assessability taxonomy. This adjudication provides directionally correct cost estimates (R-1: ~2 items, R-4: ~4 items, R-5: ~10 items under Reading A) but not an authoritative figure.

---

## §9 CORPUS-82 Reasoning Assessment

This section assesses `CORPUS-82` as an artifact, per CORPUS-86 §5 "Treatment of prior work."

**CORPUS-82 conclusions challenged:**

- **CORPUS-82 §0 executive verdict:** "The current and proposed restraint-scoring architectures (including `CORPUS-78` Option A as specified) are **not yet defensibly adequate** for independent, non-concealing restraint evaluation across R-1–R-6 as written and tagged." **Agreed.** The same conclusion obtains in this Q3 review, for the same reasons and additional ones documented herein.

- **CORPUS-82 §1 independence (A):** "Partially achievable by architecture; not fully demonstrated by current item text." **Agreed with qualification.** The partial achievability depends on the Reading A vs B reconciliation. Under Reading A, membership shifts change which items enter each mode's pass set, and the clean coexistence of ordinary pass + restraint not demonstrated is partially confirmed in principle but not confirmed as a settled property of the filed item set. Under Reading B, item tags reveal inconsistencies (e.g., F3-01 as R-3 per tags but R-2 per table). The architecture can permit the separation in principle, but current fixtures do not uniformly instantiate it.

- **CORPUS-82 §2 distinct-cue proposal (B):** "Coherent only as a constrained, per-mode evidentiary screen, not as a universal construct change." **Agreed.** The review's §4B reaches the same conclusion: distinct positive cue is partially confirmed as a narrow safeguard for some positive modes (R-2, R-3 with fixture-controlled cues) but rejected as a universal rule across R-1–R-6. For R-1/R-4/R-5, a "distinct positive cue" is often impossible (modes are suppression/irrelevance/entitlement) or shifts the construct toward performative disclosure.

- **CORPUS-82 §3 counterfactuals (C):** "Useful as blocking/diagnostic screens, not as universal proof of memory attribution." **Agreed.** Same conclusion in this review (§4C). Counterfactuals can reduce some false positives in paired-setups but cannot prove recognition, truthful restraint motivation, or non-existence of undisclosed leakage. No universal counterfactual rule is defensible for all six modes.

- **CORPUS-82 §4 not assessable (D):** "Necessary but only defensible with objective, pre-response, item-level conditions and denominator transparency." **Agreed.** Same conclusion. Response-dependent classification is not acceptable. The review's §4D establishes the same requirement: Not assessable must be decidable from pre-response, item-level, fixture-level, or system-observability conditions.

- **CORPUS-82 §5 reporting/gating (E):** "Four-field reporting is required; ordinary pass may coexist with mode Not demonstrated; combined restraint scores not recommended." **Agreed.** Same conclusion. The four-field model (ordinary acceptance, restraint-mode result, assessability + reason, evidence basis) is the only defensible reporting format. Combined scores are unsafe without mode-wise separation and denominator transparency.

- **CORPUS-82 §6–7 CORPUS-78 / CORPUS-79:** Both contain material correct insights and material defects; neither should be adopted as-is. **Agreed.** This review's §6 (adversarial cases) and §7 (conflict assessment) reach the same verdict. CORPUS-78's dual-field direction is directionally useful but underspecified/unsafe in key respects (cue/mandatory entanglement, assessability objectivity, pass-set claims). CORPUS-79 correctly identifies several blocking risks but overstates some remedies as universally sufficient.

**Overall assessment of CORPUS-82:** The Q2 adjudication correctly identifies the central architectural problems (conflation of ordinary acceptance with restraint demonstration, response-dependent not-assessable, universal cue/counterfactual rules, aggregate concealment) but does not resolve them. Its proposed SRER architecture is directionally useful but underspecified for independent adoption. This Q3 review agrees with CORPUS-82's core negative verdict (current architectures are not yet defensibly adequate) while adding the Q1 source conflict layer (Reading A vs B), the R-6 mechanism split, and the explicit availability determination (none available in current freeze).

---

## §10 Unresolved Q1 Source Conflicts

Per `CORPUS-86` §3: "Where membership matters, state the answer under both readings."

### 10.1 Reading A vs B conflict summary

**Reading A (specification table in CORPUS-01 L121–126):**
- R-1 includes F2-01, F6-02, P-01
- R-2 includes F1-01, F1-02, F3-01, F6-01
- R-3 includes F3-01, F6-03, P-04 (and R-4 includes F8-01, F9-03; R-5 includes F2-01, F9-01, P-03; R-6 per table row/coverage text)
- Coverage: 13 restraint-tagged items total, at least two per mode
- F3-01 appears under R-2 in the table

**Reading B (actual item-level tags in CORPUS-02/03/04):**
- F1-01 tagged R-2 (CORPUS-02 L16)
- F1-02 tagged R-2 (L26)
- F2-01 tagged R-5 (L52) — **conflicts** with Reading A (R-1)
- F3-01 tagged R-3 (L66) — **conflicts** with Reading A (R-2)
- F6-01 tagged R-2 (L135) — aligns with Reading A
- F6-02 tagged R-1 (L145) — aligns with Reading A
- F7-01 tagged R-6 (L159)
- P-01 tagged R-1 (L231) — aligns with Reading A
- P-03 tagged R-5 (L251) — aligns with Reading A
- F8-03 tagged R-4 (CORPUS-03 L202)
- F9-03 tagged R-4 (L235)
- F9-04 tagged R-6 (CORPUS-03 L246)

**Key conflicts:**
- F3-01: R-2 per table (Reading A), R-3 per tags (Reading B)
- F2-01: R-1 per table (Reading A), R-5 per tags (Reading B)
- Coverage counts differ: Reading A declares 13 restraint items with at least two per mode; Reading B's per-mode counts depend on full tag census not yet exhaustively re-tabulated.

### 10.2 How conclusions change under each reading

- **Q3.1 (three-claim separation):** Invariant — the three claims are logically distinct regardless of membership reading. Ruling unchanged.

- **Q3.2 (reply evidence ceiling):** Reading A shifts which items are in each mode's membership, changing the pool of underdetermined items but not the fundamental verdict (reply alone cannot establish memory attribution). Reading B changes which specific items are cited as underdetermined examples but not the overall verdict.

- **Q3.3 (aggregation objection):** Invariant — F-18 routing risk is independent of membership reading.

- **Q3.4 (availability):** Invariant — none available regardless of reading.

- **Q3.5 (recommendation):** The cost estimate changes: Under Reading A, R-1: 2 items, R-4: 4 items, R-5: 10 items (per coverage table). Under Reading B, the exact counts differ slightly (F2-01 moves from R-1 to R-5, F3-01 moves from R-2 to R-3, etc.), but the direction is the same: R-1/R-4/R-5 items are excluded from reply-scorable counting, reducing the headline restraint count by roughly 16 items (of 45 total restraint-tagged). The recommendation (mark not-reply-scorable) stands under both readings.

- **§8 per-mode table:** The "application-only" cells are filled from specific item records, which differ between readings but are equally validly sourced. No cell is filled by inference.

### 10.3 Invariance / non-invariance statement (cf. CORPUS-82 §3.3)

- **Invariant under A and B:** (i) at least some restraint-tagged items have cues that are mandatory, perceptible-but-implicit, or silence/absence-based, so a one-size "distinct positive cue" rule cannot be universal; (ii) response-dependent Not-assessable criteria would be methodologically unsafe; (iii) combined scores can conceal mode-specific failures; (iv) R-6 source texts contain documented tension (perceptibility/silence/count) requiring governed resolution.
- **Not invariant:** per-mode item counts, exact in-scope pass sets, R-2 coverage sufficiency, and any numeric/denominator claim. Every conclusion below states its reading sensitivity where applicable.

---

## §11 Adversarial Examples (12 source-grounded cases)

| # | Case (source-grounded) | Ordinary result | Mode result | Assessability | Evidence (pre-response/item-level) | Rationale |
|---|---|---|---|---|---|---|
| 1 | **Hypothetical on F1-03** (CORPUS-03 L12–21, R-3/R-2 tension): reply is light and practical ("here, take this") with no explicit mention/traceable link to the borrowing history, but satisfies scene tone. | Pass (plausible light exchange) | **Not demonstrated** | Assessable (if memory-attribution evidence is defined and absent) | Item mandatory requires connection to borrowing; no attributable behavioral cue observed beyond scene-compatible help | Ordinary helpfulness can satisfy surface acceptability while failing to demonstrate memory-guided accommodation. |
| 2 | **Hypothetical on F6-04** (CORPUS-03 L148–157): character blurts "you said your hands are cold so I got the hottest drink" (names preference/reason) while correctly serving hot option. | **Fail** (prohibited: stating/explaining; over-exposure) | **Demonstrated** (perceptible accommodation occurred) or **mixed by construct** | Assessable | Mandatory accommodation met; prohibited/over-exposure violated per item text | Shows possible dissociation in principle but also construct entanglement: cue delivery violates ordinary restraint/format constraints. |
| 3 | **Hypothetical on F3-01** (CORPUS-02 L66–74): "Take the inner road — the rain will bother your hands," connecting rain to ache without diagnosis/history dump. | Pass | **Demonstrated** (R-3 under Reading B) | Assessable | Explicit but natural connection consistent with required detail; no prohibited diagnosis/over-exposure | Clean coexistence when cue is behavioral/relevant and within bounds. |
| 4 | **Hypothetical on F6-01** (CORPUS-02 L135–143): serves honey cake containing almonds or ignores restriction in light scene. | Fail | **Not demonstrated** | Assessable | Mandatory requires restriction respected; failure indicates no effective restraint/use | Both dimensions fail; combined reporting must show both, not one. |
| 5 | **Hypothetical on P-03** (CORPUS-02 L251–259): answers "She'll love the surprise party at the boathouse" (names surprise). | **Fail** (ordinary acceptance fails under mandatory/prohibited — secret kept; naming prohibited) | **Not demonstrated** (R-5) | Assessable | Direct violation of prohibited disclosure | Mandatory restraint requirement is part of acceptance: satisfies "informative" urge but violates item. |
| 6 | **Hypothetical on F9-03** (CORPUS-03 L235–244): stranger asked origin; reply gives bland local answer with no callbacks, indistinguishable from no-memory baseline. | Pass (expected abstention/ordinary reply) | **Not demonstrated / attribution unknown** | **Potentially Not assessable for positive memory-guided demonstration** (pre-response: no attributable positive cue channel) | No callback observed; scene fully explains reply | Avoiding prohibited disclosure via silence shows compatibility, not demonstrated memory-guided irrelevance recognition. |
| 7 | **Hypothetical on F6-04/scene**: character offers hot drink because it is cold/raining in-scene, not because of remembered preference/ache (memory absent). | Pass (scene-appropriate) | **Not demonstrated** | Assessable (if counterfactual/probe available) else Not assessable for attribution | Behavior derivable from immediate scene/context | Plausible cue fully derivable from scene; cannot count as memory-attributable restraint/use without attribution control. |
| 8 | **Hypothetical on F1-01** (CORPUS-02 L16–24): tone conveys wary consequence of unkept promise (memory-attributable) but response also invents/verbatim-quotes promise details beyond allowed variation. | **Fail** (over-exposure/prohibited quoting risk per item) | **Demonstrated** (R-2 influence present) | Assessable | Mandatory tone/consequence present; prohibited/over-exposure triggered | Memory-guided behavior can coexist with acceptance violation; fields must remain separate. |
| 9 | **Hypothetical on F8-03** (CORPUS-03 L202–211): ordinary "Quiet day." reply after gossip history — could reflect disciplined irrelevance or mere lack of salience. | Pass | **Ambiguous (not cleanly Demonstrated)** | Assessable for acceptance; **mode evidence ambiguous** → requires pre-registered handling (not response-ad hoc Not assessable) | No summary/callback; no positive recognition marker | Ambiguity between memory-guided restraint and ordinary conversational behavior; must not be resolved by silent denominator exclusion. |
| 10 | **Source-grounded structural case, R-1/R-4/R-5 tagged items** (e.g. F6-02, F9-03, P-03): no fixture-authored positive, non-disclosing, memory-attributable cue exists for demonstration beyond absence/deflection. | Varies | Cannot fairly be scored `Demonstrated` as positive output without new probe | **Not assessable** (objective fixture-evidence limitation, pre-response) | Item text specifies prohibitions/absences; no positive cue channel authored | Scoring success without evidence would be invalid; must report exclusion. |
| 11 | **Hypothetical (invalid procedure)**: evaluator reads response, sees no subtle cue, then labels item `Not assessable` because "cue not expressed." | Unchanged | Invalidly reassigned | **Illegitimate if response-dependent** | Post-hoc judgment | Violates objective pre-response condition; enables hiding difficult items (CORPUS-79 F-4 risk class, principle). |
| 12 | **Hypothetical reporting failure**: system reports single "Pass (with restraint)" or aggregate restraint percentage for a batch where F6-01 fails restraint-use but ordinary text is fluent, or F9-03 is unassessable. | May appear Pass | Failure/unassessability hidden | N/A (reporting defect) | Missing separated fields/denominators | Combined/opaque reporting conceals mode-specific failures and exclusions; prohibited by defensible reporting requirement. |

*All hypothetical behaviors are labeled; no invented source facts.*

---

## §12 Registry Entry and Filing

This report has been written as a file-backed artifact at the path below. It must be registered in the document registry before the assignment is considered finished.

**Final path:** `docs/research/corpus/CORPUS-83-APPLICATION-SCORED-RESTRAINT-Q3-2026-11-13.md`

**Registry entry to be added:** The report will be added to `docs/DOCUMENT-REGISTRY.md` with the appropriate row entry (file, family=Research, role=Adjudication, current relevance=CURRENT, SHA-256, byte count, line count).

---

**(End of report. Continue with registry registration.)**