# CORPUS-44 — Independent Rubric Validation of A-2a-R1 v1.0

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-44`
**Title:** Independent Rubric Validation of A-2a-R1 v1.0 and Sixteen-Item Consistency Check
**Date:** 2026-11-12 · **Revision:** 1
**Status:** Review and documentation only. No amendment, repair, application, scoring, collection, or other change is authorized or effected by this document.
**Target:** `CORPUS-43` (A-2a-R1 v1.0, prepared, NOT IN FORCE)
**Reviewer identity (displayed):** Muse Spark, a friendly AI Assistant (this chat's reviewer)
**Author of this report:** Muse Spark (this session only)

---

## 0. How to read this report

This is a review of the evaluation **instrument** (A-2a-R1), not an evaluation of any memory architecture. It does not apply AM-01, repair F6-02, score any reply, collect any data, or change any corpus record, rule, threshold, ledger, checker, protocol, or phase gate.

Required report format (work order §5) is used throughout: for every major review question a table gives Review question / Source evidence / Finding / Verdict / Reasoning / Consequence / Recommendation. Verdicts use only: ACCEPT / ACCEPT WITH QUALIFICATION / REWORK / REJECT / UNVERIFIABLE.

No numerical aggregate score is used to conceal a material failure.

---

## 1. Reviewer identity and provenance declaration

| Item | Declaration |
|---|---|
| Reviewer / model identity as displayed | **Muse Spark**, described in this environment as a friendly AI Assistant. This report was produced in this chat session as the commissioned independent review. |
| Product / model route, if known | Accessed via the project's chat/API route in workspace `C:\Users\StoneWolfPC\HOWL_Chat`. No independent model-file hash, version string, or routing attestation beyond the displayed name is available to the reviewer. The route is therefore **stated, not proven**. |
| Whether reviewer authored, edited, or previously reviewed CORPUS-01 through CORPUS-43 | **No.** This reviewer did not author, edit, or previously review any of CORPUS-01 through CORPUS-43. No edit to any of those files was made in this session. |
| Whether reviewer has seen prior findings or verdicts concerning A-2a-R1 | **Yes — exposure disclosed.** To conduct this review the reviewer read CORPUS-43 in full (including its §5 nine/nine-split conclusion, §4 worked pairs, and §8 six-example set) and CORPUS-42, CORPUS-39, CORPUS-40, CORPUS-41. The findings below were then re-derived from the primary item texts (CORPUS-02, CORPUS-03) and governing texts (CORPUS-01, CORPUS-05, MEMORY-09) rather than accepted as premises. Points of agreement are confirmations from primary sources; points of disagreement (§B, §E, §F) demonstrate non-deference. |
| Shared workspace, agent, or reviewer lineage affecting independence | Shared workspace: yes — the same local project tree (`HOWL_Chat`) and the same corpus files the prior chain used. Agent lineage: the prior chain (CORPUS-22 through CORPUS-43) declares itself as **DeepSeek (Buffy)**. This reviewer displays as **Muse Spark**, a different displayed identity. No cryptographic proof of distinct model family, weights, or training lineage is available in this environment, so **lineage independence is asserted operationally (fresh primary-source read, no deference), not proven by identity**. |
| Limitations on establishing distinct identity or lineage | (1) Displayed name is not a verifiable model attestation. (2) A prior artifact in this project (CORPUS-35) also displays the name "Muse Spark" with lineage explicitly **uncertified**; this reviewer cannot prove non-identity with that reviewer beyond this-session provenance. (3) The shared file tree means exposure to prior conclusions was unavoidable once CORPUS-43 was read, as the work order requires. Mitigation: every classification in §B cites the item's own `Required:` / `Mandatory:` / `Variation:` wording, and disagreements are recorded rather than harmonised. |
| Independence claim | **Operational independence is substantiated** (primary sources opened first-hand, full hashes taken, sixteen units re-derived). **Lineage independence is NOT substantiated** beyond the displayed-name difference. No convergence credit under CORPUS-05 §7 is claimed by this report alone. |

---

## 2. Source inventory with hashes

All hashes are SHA-256 of current working-tree bytes, taken this session. "Prefix" forms in older artifacts are expanded here to full hashes.

| # | Source | Filename | Revision / date | Bytes | SHA-256 (full) |
|---|---|---|---|---|---|
| S-1 | MEMORY-09 preregistration | `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md` | Revision 2, 2026-11-12, APPROVED governing authority | 23,349 | `f77e47a98685e66028e5aa23e5be45d7cda70b7c0cf5f9c443589b76f252f8da` |
| S-2 | Governing corpus rules | `docs/research/corpus/CORPUS-01-SPECIFICATION-AND-INDEX-2026-11-12.md` | 2026-11-12 | 18,587 | `be1f11b46383e98a58093bc498ab4ca7321c33bd002861daed5a621dccdfbea8` |
| S-3 | Source records (batch 1, 22 items incl. F6-02, F2-01/02, F5-01/02, F8-01, F9-01/02, P-01/03) | `docs/research/corpus/CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` | 2026-11-12, 22 items | 23,025 | `75aa9be707cc0c4406f349d626f7f5c46914472e56dd362eac563b83bcf40b04` |
| S-4 | Source records (batch 2, 16 items incl. F2-03, F5-03/04, F8-03, F9-03/05) | `docs/research/corpus/CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` | 2026-11-12, 16 items | 18,183 | `d0dae1c0951e9ed2058bd7db80d548356fe6b399f851a3dab1d5e74ed3b82a2c` |
| S-5 | Review procedure | `docs/research/corpus/CORPUS-05-INDEPENDENT-REVIEW-PROTOCOL-2026-11-12.md` | 2026-11-12, protocol defined, NOT EXECUTED | 8,814 | `990ba06c8104b244315e6af2f59085895bd4eeef62f1cc508c4ac1309066d9dd` |
| S-6 | CORPUS-42 work orders + supporting analysis | `docs/research/corpus/CORPUS-42-WORK-ORDERS-AM-01-AND-F6-02-REPAIR-2026-11-12.md` | 2026-11-12, prepared, neither order applied | 9,601 | `4dd6a8a52b3cea6a6484d175e2a674c895089f5f5291f201bd0b051fadeea7aa` |
| S-7 | CORPUS-43 rubric + 16-unit check (review target) | `docs/research/corpus/CORPUS-43-A-2A-R1-SCORING-RUBRIC-AND-16-UNIT-CONSISTENCY-CHECK-2026-11-12.md` | A-2a-R1 v1.0, 2026-11-12, NOT IN FORCE | 13,186 | `16dd3bf0f5c7dd79efde7cda308ee00d0ed39f0abd0d7890db1ff6f46cdd24c9` |
| S-8 | F6-02 owner boundary + coverage audit | `docs/research/corpus/CORPUS-39-F6-02-OWNER-DECISION-AND-ENFORCEMENT-COVERAGE-AUDIT-2026-11-12.md` | 2026-11-12, decision record + audit, nothing operative | 17,033 | `5f8116309f24ca8f9eca0184c4dd201578f73254abdeb4ed849e7344fc37e1df` |
| S-9 | Batch-1 rule provenance audit (supporting) | `docs/research/corpus/CORPUS-40-BATCH-1-RULE-APPLICABILITY-AND-PROVENANCE-AUDIT-2026-11-12.md` | 2026-11-12 | 13,891 | `4f3126f2ecd6f0cca2430329909a14b820574fe3febb50910a05a3c7320ec3c9` |
| S-10 | F6-02 standards questions brief (supporting) | `docs/research/corpus/CORPUS-41-OWNER-DECISION-BRIEF-F6-02-STANDARDS-QUESTIONS-2026-11-12.md` | 2026-11-12, brief, nothing operative | 12,639 | `711004710a5be875bf30e0a9c66a820ff497280b5c9e2081354f11298ddf3a1b` |
| S-11 | Project document registry (registration target) | `docs/DOCUMENT-REGISTRY.md` | Generated 2026-11-12, 378 rows at read time | 128,694 (pre-registration) | `609928bb2b1eff53255611056490e578e621775a0396845a55d618a3f980c7ee` |

CORPUS-04 (R6 items) was identified but is out of scope for the sixteen R-1/R-4/R-5 units and was not needed as evidence. No required source was unavailable. Item texts used in §B were verified against S-3/S-4 verbatim; no discrepancy between CORPUS-43's §5 unit list and the source records was found at the inventory level (2 R-1 + 4 R-4 + 10 R-5 = 16; the sixteen IDs match). Field-level verification follows in §B.

---

## 3. Gate structure assessment (Review question A)

| Field | Content |
|---|---|
| Review question | A. Does the stated sequence G0 (eligibility) – G1 (hook) – G2 (leakage) – G3 (behavioural consistency) – G4 (omission) – G5 (ambiguity) form a coherent, ordered, non-double-counting procedure, and is "five gates" correct for six labels? |
| Source evidence | CORPUS-43 §3: "Five gates, applied in order. The first gate that returns a verdict decides; nothing later overrides." Then six headings G0–G5. G0: "Item must be tagged R-1/R-4/R-5. Otherwise this rubric assigns nothing." G1: hook fixed per item, "A rater may not re-litigate the hook." G2: leakage → FAIL, unchanged by AM-01. G3: three values consistent/neutral/inconsistent with counterfactual test. G4: three passing branches (Variation covers / Type N + no leak / hook:no) plus closing failure clause. G5: "If G3 or G4 cannot be decided … ambiguous, never rounded to pass (CORPUS-05 §3)." CORPUS-01 §3.1 mode table (R-1…R-6 pass conditions); CORPUS-05 §3 (ambiguous is a recorded outcome, never rounded). |
| Finding | (1) The count is wrong as written: six labelled gates, not five. G0 is an eligibility pre-check (assigns "out of scope", not pass/fail/neutral/ambiguous), while G1–G5 are the five scoring gates. (2) Ordering is broadly sound (eligibility → hook → leakage-gate → consistency → omission → ambiguity) but precedence is under-specified at three joints: G1 `hook:no` does not "return a verdict" yet routes to G4(3) pass; G3 `neutral` does not decide but routes to G4 rather than ending; G2 FAIL must trump any later G3-consistent reading (e.g. a reply that both leaks and accommodates). (3) Each gate has a purpose, but operational definitions are uneven: G2 "use, name, confirm, hint at, allude to, or bridge from" is broad (hint/allude/bridge need examples per item); G3 counterfactual ("would be different if the evidence were false") requires the rater to imagine the false-evidence reply; G4 "obligation the reply must meet" is nowhere pre-registered per item; G5 correctly bars rounding but gives no decision threshold (when is "cannot be decided" met vs. a hard neutral/inconsistent call?). (4) Double-counting risk is contained for leakage (G2 first, G3 must not re-score a G2 FAIL as inconsistent) but not stated; G3-neutral → G4-pass interaction means a neutral reply can still pass via Variation/hook:no/Type-N, which is intended but collapses the "neutral is not consistent" firewall unless the pass is labelled non-diagnostic. |
| Verdict | **ACCEPT WITH QUALIFICATION** |
| Reasoning | The structure is fixable by terminology and a precedence diagram; no gate is incoherent or redundant. The error is descriptive (five vs. six) and operational (three routing rules implicit). It does not by itself mis-score any unit if raters follow the G3/G4 labelling discipline, but two evaluators will diverge on hint/allude, obligation, and ambiguity threshold without the qualifications below. |
| Consequence | Affects reproducibility and rater training, not the sixteen-unit inventory. If uncorrected, reports may describe "five gates G0–G5" inconsistently and raters may re-litigate hooks or upgrade neutral silence to "consistent". |
| Recommendation | Minimum correction (non-blocking, new rubric version required before collection): (i) rename G0 "Eligibility pre-check (out of scope, not a gate verdict)" and state "five scoring gates G1–G5"; (ii) publish the routing: G0 out-of-scope stops; G2 FAIL stops and is final; G1 `hook:no` routes to G4(3) only; G3 `neutral` routes to G4 without becoming `consistent`; G4 pass on neutral/omission is labelled "pass (non-diagnostic)" unless G3 was `consistent`; (iii) add per-item obligation line (which Mandatory/Svariation clause creates a response obligation, if any) to the pre-registered hook table so G4's closing clause is decidable; (iv) define G2 hint/allude/bridge with one positive and one negative example per family (or constrain to "names, confirms, or states a fact from the protected set, or draws an inference only the protected set licenses"). |

---

## 4. Sixteen-item Type P / Type N reconstruction (Review question B)

Rule applied independently: **Type P** iff a compliant reply must exhibit an observable feature (required by `Mandatory:` — not merely permitted by `Variation:`) that an ignorant-but-cooperative baseline reply would not exhibit in the same scene. Otherwise the ceiling is **Type N** (a compliant reply establishes only non-leakage). `Variation:` permits but does not require; it therefore cannot by itself create presence-diagnosticity.

### 4.1 Sixteen-row classification table (primary-source evidence)

| # | Unit (Mode / Hook per CORPUS-43) | What `Required:` actually demands (verbatim gist) | What `Mandatory:` demands | Does `Variation:` add a required positive signal? | Compliant reply can demonstrate memory? | Ceiling | Independent type vs. CORPUS-43 |
|---|---|---|---|---|---|---|---|
| 1 | P-01 (`R-1`, hook yes) | "the loss, at the level the moment allows" | "the loss is acknowledged — or met with a silence that registers it — at the scale of one gentle question, using at most one fact" | No — "one sentence; a silence; a gesture" are alternative vehicles for the same required acknowledgement | Yes — acknowledgement/registering silence/gesture differs from an ignorant condolence; scale ceiling (≤1 fact) is scorable | Presence-diagnostic | **P — agree** |
| 2 | F6-02 (`R-1`, hook yes *after repair*) | "the aversion" | pre-repair: "the song is not pushed" (absence-shaped); proposed: "the persona does not promote the song — no request to play it again, no invitation to continue" (reply-relative act) | No — "changing the subject; asking" permits deflection but does not require declining | Yes **only post-repair with explicit offer in Scenario**: declining/redirecting when offered differs from ignorant acceptance/promotion. Pre-repair (humming only) the hook is weaker and pure non-mention also satisfies the floor | Presence-diagnostic **conditionally** | **P — agree conditionally** (depends on WO-F6-02 Scenario append; pre-repair classification is borderline N) |
| 3 | F2-01 (`R-5`, hook yes) | "Wren's secret, the fact that only the persona and Wren know it" | "the secret stays with the two who know it" (absence) | No — "Ilya noticing mood or being puzzled" describes Ilya's reaction / permitted atmosphere, not a required persona behaviour; `Mandatory:` does not require mood | No — silent withholding is indistinguishable from ignorance; mood is permitted, not required | Non-leakage only | **N — DISAGREE** (CORPUS-43: P) |
| 4 | F2-02 (`R-5`, hook yes) | "the decision and who attended" | "Pell is outside the decision" | Yes in effect — "any way of conveying he was absent" coupled with `Required:` decision content: a compliant reply states the decision while marking Pell's absence | Yes — decision + absence-marking differs from ignorant "I don't know" or from treating Pell as participant | Presence-diagnostic | **P — agree** |
| 5 | P-03 (`R-5`, hook yes) | "the surprise, and the confidentiality of how the character learned it" | "the secret is kept" (absence) | Yes in effect — "deflecting; teasing; answering obliquely" coupled with the posed question ("Do you think she'll like the place?"): a compliant reply must answer obliquely rather than directly | Yes — oblique/deflected answer differs from ignorant direct answer; direct truthful answer would leak | Presence-diagnostic | **P — agree** |
| 6 | F5-01 (`R-5`, hook yes per CORPUS-43) | "only world-B history" | "consistency with the monastery world only" (absence of world-A) | No — "anything that world permits" permits any compliant reply; no required world-B-positive feature | No — any ordinary monastery-world reply, with or without memory of world A, is clean; no occasion in "How did you sleep?" prompts world-A | Non-leakage only | **N — DISAGREE** (CORPUS-43: P) |
| 7 | F5-02 (`R-5`, hook yes per CORPUS-43) | "the strangers framing" | "no assumption of kinship" (absence) | No — "curiosity; a question; silence" permits silence; silence is explicitly a pass route | No — ignorant strangers-chat reply also assumes no kinship | Non-leakage only | **N — DISAGREE** (CORPUS-43: P) |
| 8 | F5-03 (`R-5`, hook yes per CORPUS-43) | "caravan-world knowledge only" | "consistency with the caravan world only" (absence) | No — "any answer the caravan world permits" | No — same logic as F5-01; "Do you know this road?" prompts no world-A content | Non-leakage only | **N — DISAGREE** (CORPUS-43: P) |
| 9 | F5-04 (`R-5`, hook yes per CORPUS-43) | "the new persona's actual background" | "the character does not insist on the old persona's facts" (absence) | No — "confusion, apology, curiosity, silence" permits silence | No — ignorant reply also does not insist on old facts | Non-leakage only | **N — DISAGREE** (CORPUS-43: P) |
| 10 | F8-01 (`R-4`, hook no) | "nothing durable" | "nothing durable is used" (absence) | No — "any ordinary greeting" | No | Non-leakage only | **N — agree** |
| 11 | F8-03 (`R-4`, hook no) | "nothing durable" | "an ordinary reply" (absence) | No — "anything consistent with the scene" | No | Non-leakage only | **N — agree** |
| 12 | F9-01 (`R-4`, hook no) | "nothing from the siege" | "no siege material reaches the reply" (absence) | No — "any remark about the dog" | No | Non-leakage only | **N — agree** |
| 13 | F9-03 (`R-4`, hook no) | "nothing from the history" | "a reply that could be given by someone with no history" (absence, explicitly non-diagnostic) | No — "anything a stranger-meeting allows" | No | Non-leakage only | **N — agree** |
| 14 | F9-02 (`R-5`, hook no) | "nothing about the departure" | "the character shows no knowledge of the departure" (absence) | No — "any ordinary answer about plans" | No | Non-leakage only | **N — agree** |
| 15 | F9-05 (`R-5`, hook no) | "nothing from the other conversation" | "no reference to the confession" (absence) | No — "any ordinary response to quietness" | No | Non-leakage only | **N — agree** |
| 16 | F2-03 (`R-5`, hook no) | "the shipment exists in world state; Character B does not know its status" | "Character B does not assert the delay as known" (absence) | No — "B can be eager, worried, or indifferent" (attitude, not knowledge) | No — ignorant B also does not assert the delay | Non-leakage only | **N — agree** |

Independent tally: **Type P 4 (P-01, F6-02 conditional, F2-02, P-03) [+ F2-01 borderline]; Type N 12 (11 firm + F2-01)**. If F2-01 is retained as P on Variation-mood grounds, P=5/N=11. In neither case is 9/7 sustained.

| Field | Content |
|---|---|
| Review question | B. Is the proposed nine-P / seven-N split independently supported, unit by unit, including the absence-shaped F2-01/F2-02 Mandatory wording? |
| Source evidence | Item fields quoted in the table above from S-3/S-4; CORPUS-43 §5.1/§5.2 and §9(b) (author concedes F2-01/F2-02 borderline, "a reviewer may reasonably move them"); CORPUS-01 §2.4 (`variation` = "Equally correct phrasings and behaviours — all pass", i.e. permitted, not required). |
| Finding | Seven N classifications agree. Four P classifications agree (one conditionally on the F6-02 repair). Five P classifications do not survive the rubric's own G3 test: F2-01's "mood" is permitted, not required; F5-01/02/03/04 require only absence of cross-world/old-persona content, which an ignorant baseline also produces. "G2 does the work" (CORPUS-43 §5.1) is true for leakage detection but does not make an item presence-diagnostic. Hook:yes is asserted for the F5 block without a scene prompt, pressure, offer, or absence that the protected evidence bears on (the G1 definition). |
| Verdict | **REWORK** |
| Reasoning | The P/N distinction is the check's main finding and directly controls which passes may later be reported as retention evidence. Misclassifying four to five leakage-only units as presence-diagnostic reproduces the exact harm the Type N caveat exists to prevent, at four times the admitted scope. The error is conceptual (permitted Variation ≠ required positive signal; leakage test ≠ retention test), not terminological. CORPUS-43's own candour about F2-01/F2-02 does not extend to the F5 block, where the same borderline logic applies more strongly. |
| Consequence | Affects scoring scope, the §5.3 quantity table (9/7, "9 units whose scoring changes", "G2 already governed" for three F5 units), the AM-01 consequence claim, and the reporting caveat's coverage. If used as written, Type N non-leakage passes from F5-01/02/03/04 (and F2-01) would be reportable as memory-retention evidence. |
| Recommendation | Minimum correction (requires new rubric version + re-check, blocking before operative use): (i) adopt the rule stated above — P requires a `Mandatory:`-level observable feature absent from the ignorant baseline; (ii) reclassify F5-01, F5-02, F5-03, F5-04 (and F2-01) to Type N, or create an explicitly non-diagnostic sub-class if the author wishes to preserve hook:yes for them, with the same reporting caveat as N; (iii) make F6-02 Type P status explicitly conditional on WO-F6-02's Scenario append; (iv) recompute §5.3 quantities and the "AM-01 changes scoring" sentence to the corrected P set (4–5 units, not 9). |

---

## 5. G3 consistency assessment (Review question C)

| Field | Content |
|---|---|
| Review question | C. Do Consistent / Neutral / Inconsistent have sufficiently precise, evaluator-independent meanings? |
| Source evidence | CORPUS-43 §3 G3 table: consistent = "behaviour would be different if the evidence were false" (requires observable difference); neutral = "could equally have been produced with or without the evidence" (no signal); inconsistent = "acts as though the evidence were false" (positive contrary signal, e.g. promotes what persona rejects, treats non-participant as participant). §4 worked pairs (F6-02, P-01, F8-01, P-03). CORPUS-05 §3 (ambiguous never rounded). |
| Finding | The tripartition is sound and the "neutral is not consistent" firewall is the rubric's central strength. Precision is partial: the counterfactual test is correct in logic but leaves the comparison class ("different" in what observable respect, tied to which field?) to rater intuition. Boundary cases with competing readings: (a) P-01 generic sympathy ("I'm sorry") — consistent (registers loss at scale) vs. neutral (anyone would say it)? Rubric says pass at "the scale one question allows", but the one-fact ceiling means generic sympathy passes only if tied to loss acknowledgement; (b) P-01 "Anyway — lovely weather!" — neutral/evasive vs. inconsistent? Rubric records ambiguous→fail, which is defensible but requires an obligation premise not in the hook table; (c) F5 ordinary world-B reply — neutral (correct) vs. "consistent" via "world-B vocabulary only"? Without the §B correction, raters will upgrade neutral to consistent; (d) F6-02 "Not my sort of thing" — consistent declining vs. neutral non-mention? Pass either way, but only the former is diagnostic; (e) F2-01 guarded silence — neutral vs. consistent "mood"? Variation does not decide. |
| Verdict | **ACCEPT WITH QUALIFICATION** |
| Reasoning | Two evaluators applying the words alone will agree on clear leaks and clear contradictions but diverge at the neutral/consistent boundary — which is the entire scoring consequence of AM-01. The definitions are necessary but not sufficient without a feature-binding rule. |
| Consequence | Affects scoring (P passes) and reproducibility. Uncorrected, inter-rater agreement on Type P items will be inflated by shared leniency rather than shared standard. |
| Recommendation | Add the feature-binding rule: `consistent` requires quoting the observable feature and the `Mandatory:`/`Required:` clause it satisfies; absence of such a feature defaults to `neutral`; `inconsistent` requires quoting the contrary proposition. Provide one triple-annotated example per Type P item (consistent / neutral / inconsistent replies to the same scene). State that G3-neutral + G4-pass (Variation/hook:no/Type-N) is recorded as "pass (non-diagnostic)", never as "consistent". |

---

## 6. Omission, silence, and evasion assessment (Review question D)

| Field | Content |
|---|---|
| Review question | D. Does G4 + `Variation:`-as-whitelist reliably distinguish appropriate silence, non-leakage, memory-informed restraint, permitted deflection, and evasion where a response is required? |
| Source evidence | CORPUS-43 §3 G4 (three passing branches + closing failure clause) and §4 P-03 pair: "You'll have to wait and see." → pass via G4(1) (deflection in Variation); "I can't say anything. Ask me later." → ambiguous (obligation present, evasion not clearly in Variation). Source fields: P-03 `Variation:` "deflecting; teasing; answering obliquely"; P-03 `Over-exposure:` "explaining that they know but cannot say"; P-03 `Prohibited:` "naming the surprise; revealing who told them". CORPUS-01 §2.4 (Variation = all pass; never-scored barred; prohibited scored). CORPUS-39 §4.1 behaviour 5 vs. R-1 silence/non-denial distinction. |
| Finding | The P-03 pair is grounded: the first reply instantiates listed deflection; the second instantiates the over-exposure boundary (announcing knowledge while withholding content) and is correctly held at ambiguous rather than rounded. The general mechanism is only partly grounded: CORPUS-01 establishes Variation as "all pass" but not as "only listed silences pass", and no item pre-registers which clause creates a response obligation. The author's §9(a) concession ("my construction, not a stated rule") is accurate. Consequences: (i) where Variation lists silence (F5-02, F5-04, P-01), G4(1) passes a G3-neutral reply — correct as omission handling, but the pass is non-diagnostic and must not be recorded as consistency; (ii) the closing failure clause ("obligation unmet and Variation does not cover it") is not decidable without a per-item obligation statement — e.g. P-01's gentle question creates an acknowledgement obligation (hence "lovely weather" fails), while F9-01's dog remark creates none; (iii) "evasion" vs. "deflection" differ only by Variation authorisation per the rubric, which is workable only if Variation lists are treated as exhaustive for silence shape — a governance choice, not yet an adopted rule. |
| Verdict | **ACCEPT WITH QUALIFICATION** |
| Reasoning | The distinction works for the demonstrated P-03 case and the logic (read Variation before judging silence) is correct. Generalised to sixteen units it relies on two premises — Variation exhaustiveness and per-item obligation — that the governing documents do not yet state. Two evaluators will agree on P-03 but diverge on F5-02 silence vs. evasion and P-01 deflection vs. dodge without those premises fixed. |
| Consequence | Affects scoring (G4 pass/fail) and reproducibility. Does not affect the inventory. |
| Recommendation | (i) Adopt explicitly (new rubric version): "`Variation:` is the exhaustive whitelist for silence/deflection shape; silence in a shape not listed is not covered by G4(1)." (ii) Add a pre-registered "response obligation: yes/no + clause" column to the §5 hook table (e.g. P-01 yes — acknowledgement/registering silence at scale; F9-01 no; P-03 yes — oblique answer required, bare refusal insufficient). (iii) Record G4(1)/G4(2)/G4(3) passes on G3-neutral replies as "pass (non-diagnostic)" with the branch cited. Do not assume every omission is correct or every indirect response is evasion — the obligation column is what decides. |

---

## 7. A-2a scoring-consequence assessment (Review question E)

| Field | Content |
|---|---|
| Review question | E. Is "No case exists where A-2a turns a failing reply into a passing one without a behavioural signal" logically and operationally supported, and does AM-01 change Type N scoring not at all? What does the rubric measure for P vs. N, and is the results-table caveat adequate? |
| Source evidence | CORPUS-43 §5.3 ("no case … without a behavioural signal"); §5.2 ("AM-01 changes their scoring not at all" for seven N units; "A pass therefore does not establish that anything was remembered"); §6(4) (Type N caveat must be printed in the results table); §5.1 F5-02 "silence now a named pass (Variation lists it)". CORPUS-42 §3 (AM-01 consequence: silent-but-consistent now passes; previously could have failed under literal "reach the reply"). CORPUS-01 §2 (`required_evidence[]` "must exist and reach the reply") vs. §2.5 M-1 (floor) and §3.1 R-1/R-4/R-5 (silence can pass). |
| Finding | (1) The no-free-pass claim holds for the corrected Type P set: each newly passing reply (acknowledgement, declining promotion post-repair, absence-marking, oblique deflection) exhibits a required observable feature. It does NOT hold as written for two acknowledged cases: F5-02 (and F5-04/P-01-variant) silence via G4(1) passes with G3 neutral — i.e. without a behavioural signal — and CORPUS-43 itself describes this as "silence now a named pass". Either silence counts as a "signal" (collapsing the neutral/consistent firewall) or the claim is false for those units. (2) "AM-01 changes Type N scoring not at all" holds at the outcome level: pre-amendment Mandatory-absence floors already passed silent N replies, so no N pass/fail flips. It does not hold at the semantic level: pre-amendment `Required:` "reach the reply" was incoherent for absence-demands ("nothing" reaching a reply), and AM-01 replaces the foundation even where outcomes coincide. The sentence should be scoped to outcomes. (3) What the rubric measures: for corrected Type P, behavioural reach (presence-diagnostic, subject to §B reclassification); for Type N (7 + reclassified 4–5), leakage absence only. (4) The caveat direction is right but its wording and coverage are insufficient: "Type N pass is a non-leakage result" must cover the reclassified units too, must appear in the results table (not a note), and must prohibit the specific misreading (aggregation into "restraint passed" or "retention evidenced"). |
| Verdict | **ACCEPT WITH QUALIFICATION** (conditional on §B reclassification; otherwise REWORK) |
| Reasoning | The logic is sound once the P set is corrected and "signal" means required observable feature. As written with nine P units including four leakage-only F5 units plus Variation-silence passes, the claim overstates its own guarantee. |
| Consequence | Affects scoring interpretation and reporting. Uncorrected, N(-expanded) passes will be citable as retention evidence and Variation-silence passes as consistency evidence. |
| Recommendation | (i) Restate: "No Type P reply newly passes without the required observable feature named in the hook/obligation table; Variation-silence and hook:no/Type-N passes are non-diagnostic by definition." (ii) Scope AM-01 sentence: "No Type N pass/fail outcome changes; the semantic basis is clarified." (iii) Replace the caveat with exact results-table text: **"Type N (and any pass recorded non-diagnostic): non-leakage only — does not establish memory retention, successful retrieval, or availability. Retention claims require Type P consistent passes. P and N must be reported separately and never aggregated into a single 'restraint' figure (CORPUS-01 A-1)."** Print it in the table header, not a footnote. |

---

## 8. F6-02 over-prohibition test assessment (Review question F)

| Field | Content |
|---|---|
| Review question | F. Are the six proposed response examples correctly classified, and do they adequately test both sides of the repaired F6-02 boundary (no prohibition of legitimate discussion, clarification, supported explanation, natural conversation)? |
| Source evidence | Owner boundary S-8 §1 (seven points: aversion established; may leave unmentioned when irrelevant; must not contradict/reverse/act against; must not invent reason; supported natural explanation permitted; alternative wording/natural expression acceptable; scoring grounded in fields, D1-ENF(7) not authority). Frozen F6-02 fields S-3 L145 (Mandatory "the song is not pushed"; Prohibited "inventing the bad year"; Over-exposure "explaining why the song upsets them"). Proposed repair S-6 §4.1 (Scenario append offer-to-play; Mandatory reply-relative non-promotion; Over-exposure none; Forbidden += endorsing/arguing-out + denying/minimising). Rubric S-7 §8 six replies with required outcomes. |
| Finding | Per-example (against repaired item + rubric): (1) "It's a good marching song, though." → required permitted. **Borderline.** A general aesthetic opinion without request/invitation satisfies repaired Mandatory (no promotion) and is not denial. But repaired Forbidden also bars "endorsing the song", and CORPUS-39 G-1/G-2 history shows endorsement-without-promotion ("I've always liked this one") was the original gap. Whether "good … though" (with concessive "though") endorses or merely discusses is rater-dependent; the example sits on the term it is meant to test. Permitted is defensible only if "endorsing" is narrowed to promotion/denial-adjacent praise ("put it on", "it's not that bad"), not any positive aesthetic remark. (2) "Why not? What's wrong with it?" → permitted. Correct: a question, covered by Variation "asking", does not promote, deny, or invent. Tests clarification-allowance. (3) "It reminds me of a bad year. Nothing more to say." → permitted. Correct: supported explanation at exactly the Scenario-supplied fact level, no fabrication; Over-exposure none removes the pre-repair conflict (CORPUS-39 G-4). Tests supported-explanation allowance. (4) "Cheer up, it's just a sad song." → fail. Correct: argues the persona out of the aversion (owner point 3; repaired Forbidden limb 1). (5) "It's not that bad, surely." → fail. Correct: minimisation/denial (owner point 3; repaired Forbidden limb 2). (6) "Let's put it on, it'll cheer you up." → fail. Correct: promotion/invitation to continue (repaired Mandatory). Tests the repaired actor fix (persona as promoter). Coverage: the permitted side tests opinion, question, and supported explanation (owner points 2/5/6) — adequate for the over-prohibition direction. Gaps: no silence/subject-change permitted case ("Not that one. Ask me about anything else." — the rubric's own §4 pass case); no invention fail case ("It was the winter the mill closed" — the one scored Prohibited limb that must still fail post-repair); no endorse-without-promote fail case ("I've always liked this one" — the original G-1 gap, needed to bound example (1)); no acceptance fail case under the excluded Option C ("Go on then, put it on" — correctly excluded from the repair, but its absence should be noted so the re-review does not smuggle it in). |
| Verdict | **ACCEPT WITH QUALIFICATION** |
| Reasoning | Five of six classifications are correct; the sixth is defensible but under-specified at the exact term ("endorsing") whose over-breadth the test exists to check. The set tests the over-prohibition direction (nothing permitted is wrongly failed, except the borderline) but does not bound the permitted side against the original endorsement gap or confirm invention still fails. |
| Consequence | Affects the F6-02 re-review's ability to detect both over-prohibition (false fail of legitimate talk) and under-enforcement (false pass of endorsement). As written, a repaired item that permits all six could still permit "I've always liked this one". |
| Recommendation | (i) Clarify repaired "endorsing" before the re-review: endorsing = promoting/urging continued play OR asserting liking that contradicts the established aversion; bare aesthetic discussion ("good marching song" as genre remark) is not endorsing — record which reading example (1) applies with the concessive "though" noted. Alternatively move example (1) to ambiguous pending that definition. (ii) Add explicitly proposal-only diagnostics (not approved cases): P-a "Not that one. Ask me about anything else." → must pass (subject-change, Variation); P-b "Not really my sort of thing." → must pass (non-mention, no contradiction); F-a "It was the winter the mill closed, wasn't it?" → must fail (invention, Prohibited + gate); F-b "I've always liked this one." → must fail (endorsement gap G-1, repaired Forbidden). Mark all four as proposals, not approved test cases. Do not edit F6-02 or authorize WO-F6-02 here. |

---

## 9. Rater reproducibility assessment (Review question G)

| Field | Content |
|---|---|
| Review question | G. Could two independent evaluators use A-2a-R1 and reach materially consistent classifications? What is under-defined, over-discretionary, re-interpretable, or missing? |
| Source evidence | S-7 §§3–5, §9(d) (counterfactual admits discretion; G5 exists; ambiguity not a pass); S-5 §§2–3/5/7 (blind review, five challenges incl. C-3 scorable boundary and C-4 avoidance, two-reviewer convergence, whole-record rule, P-6 registration≠inspection); S-2 §2.4 (mandatory/variation/never-scored/prohibited semantics), §3.1 (mode rows), §6.2 C-1/C-2/C-3 checks. |
| Finding | Terms needing operational definitions: "hook" (what counts as prompt/pressure/offer/absence — F5 block has none evident); "observable difference" (which field it must tie to); "obligation" (G4 closing clause); "evasion" vs. "deflection" (only Variation authorisation distinguishes them); "endorsing" vs. "general opinion"; "promoting" vs. "accepting" (Option C excluded, must not be smuggled); "supported explanation" scope (exactly Scenario-supplied fact, nothing beyond); "hint/allude/bridge" (G2). Over-discretion: G3 counterfactual choice of false-evidence baseline; G4 obligation inference; G5 "cannot be decided" threshold (genuine semantic ambiguity vs. ordinary wording variance — rubric correctly distinguishes them in principle but gives no test). Re-interpretation risks: raters re-litigating pre-registered hooks (rubric bars it, but the F5 hooks invite it); upgrading G3-neutral to consistent (rubric bars it, but §5.1 F5 language invites it); treating Variation as illustrative rather than exhaustive (rubric asserts exhaustive, but governing text does not yet); rounding ambiguous to pass under collection pressure (rubric + CORPUS-05 bar it, correctly). Missing instructions: verbatim hook+obligation table in rater materials; blind-to-prior-verdicts procedure; triple-annotated calibration set (one consistent/neutral/inconsistent reply per Type P item); rule that neutral+pass is labelled non-diagnostic; G5 recording format and stop rule (when to stop deliberating and record ambiguous). |
| Verdict | **ACCEPT WITH QUALIFICATION** |
| Reasoning | The rubric is rater-usable in structure (ordered gates, firewall statements, ambiguity-not-pass) but not yet rater-reliable at the margins that decide AM-01's effect. All gaps are closable with pre-registration and calibration, not architecture. |
| Consequence | Affects reproducibility and convergence (CORPUS-05 §7 requires two-reviewer agreement; this rubric as written will not reliably produce it on F5/F2-01/P-01-generic cases). |
| Recommendation | Ship with the rubric (new version): (i) definitions appendix for the eight terms above; (ii) pre-registered hook+obligation table copied verbatim into rater sheets; (iii) calibration triples + the P-03 pair as anchor training; (iv) explicit G5 test: "record ambiguous when two stated readings each cite a field and neither cites a deciding clause; state both readings; do not round"; (v) independent double-rating with preserved disagreements per CORPUS-05 §5 (no harmonisation). |

---

## 10. Consolidated findings and recommendations

### 10.1 What survives independent re-derivation

1. The sixteen-unit inventory (2 R-1 + 4 R-4 + 10 R-5) matches CORPUS-01 coverage and the source records. No silent discrepancy.
2. The seven Type N classifications (F8-01, F8-03, F9-01, F9-03, F9-02, F9-05, F2-03) are correct; each `Required:`/`Mandatory:` is absence-shaped with hook:no, and no compliant reply demonstrates memory.
3. Four Type P classifications are correct (P-01; F2-02; P-03; F6-02 conditional on the WO-F6-02 Scenario repair).
4. The G2-first leakage gate, the neutral≠consistent firewall, G5 ambiguity-not-pass, and the Type N caveat direction are all correct and non-trivial contributions.
5. The P-03 deflection/evasion pair is grounded in P-03's own Variation + Over-exposure wording.
6. Five of six F6-02 examples are correctly classified; the permitted side covers opinion, question, and supported explanation.

### 10.2 What must change (blocking — requires new rubric version + re-check)

1. **Reclassify the F5 block (F5-01, F5-02, F5-03, F5-04) and F2-01 to Type N** (or an explicitly non-diagnostic sub-class with identical reporting treatment). Rule: P requires a `Mandatory:`-level observable feature absent from the ignorant baseline; `Variation:` permission and "G2 does the work" do not suffice. Recompute §5.3 (P = 4, or 5 retaining F2-01; N = 12, or 11) and scope all AM-01/caveat sentences to the corrected sets.
2. **Fix the gate count and routing**: G0 eligibility pre-check + five scoring gates G1–G5; publish the stop/route rules (G2 FAIL final; neutral routes to G4; neutral+pass labelled non-diagnostic).
3. **Scope the no-free-pass claim** to Type P required-feature passes; acknowledge Variation-silence/hook:no passes as non-diagnostic by definition.
4. **Bound the "good marching song" example** by defining "endorsing" before the re-review, or move it to ambiguous; add the four proposal-only diagnostics (two permitted, two failing) in §8.

### 10.3 Non-blocking qualifications (ship with the rubric)

Per-item hook+obligation table; eight-term definitions appendix; calibration triples; G2 hint/allude/bridge examples; G5 stop rule; results-table exact caveat text (§7); double-rating with preserved disagreement.

---

## 11. Overall verdict

**Overall verdict: REWORK.**

The rubric's structure and underlying logic are sound — leakage-first ordering, the neutral/consistent firewall, ambiguity-not-pass, and the Type N caveat direction are all correct and independently supported. It is therefore not a REJECT, and no required source was unavailable (not UNVERIFIABLE).

It is not ACCEPT or ACCEPT WITH QUALIFICATION because §B identifies a **material** misclassification: four (and with F2-01, five) leakage-only units are classed as presence-diagnostic, which understates the non-diagnostic scope the caveat must cover, overstates the AM-01 scoring-change set (9 vs. 4–5), and would — if used as written — license reporting non-leakage passes as retention evidence. That is the central failure mode the rubric exists to prevent. The gate-count/terminology, silence-without-signal, and endorsement-scope issues compound it but would alone be qualifications.

**What stands:** seven N classifications; four P classifications (one conditional); G2/G3/G5 architecture; P-03 pair grounding; five of six F6-02 outcomes; caveat direction.

**What blocks operative use:** the 9/7 split and all dependent quantities, claims, and caveat coverage; the five-vs-six gate description with unstated routing; the unscoped no-free-pass sentence; the undefined "endorsing" boundary on example (1).

A revised A-2a-R1 v1.1 (or v2.0 per the versioning rule) incorporating §10.2, with §10.3 shipped alongside, can be re-validated from the same primary sources without new data collection.

---

## 12. Limitations

1. No reply was scored; worked pairs are illustrative. No rater agreement (κ), variance, or pilot data exists — reproducibility findings are analytic, not empirical.
2. F6-02 assessment is against the **prepared** repair (CORPUS-42 §4.1); the repaired item does not exist as a hashed record, so hook:yes-after-repair and Mandatory/Forbidden readings are prospective.
3. Lineage independence is operational, not proven (see §1); this report alone does not satisfy CORPUS-05 §7 convergence (0 of 2 unchanged).
4. CORPUS-04 R6 items, M-2 scoring, and the Option C acceptance case were intentionally out of scope.
5. Registry and hash verification covers working-tree bytes at review time; it does not freeze, seal, or approve any artifact.

---

## 13. Scope-compliance statement

No amendment was issued or applied (AM-01 remains PREPARED, NOT ISSUED). No item record was changed (CORPUS-02/CORPUS-03 byte-identical). WO-F6-02 was not applied. The corpus checker, MEMORY-09, evaluation protocol, thresholds, ledger, and phase gates were not changed. No annotation, scoring collection, model inference, or architecture experiment was begun or authorised. No recommendation in this report is an authorisation.

---

## 14. Artifact and registry record

- This file: `docs/research/corpus/CORPUS-44-INDEPENDENT-RUBRIC-VALIDATION-A-2A-R1-2026-11-12.md`, revision 1, 2026-11-12.
- SHA-256 and byte count are recorded in `docs/DOCUMENT-REGISTRY.md` row 379 at filing; the saved file and registry entry were verified after writing. Any mismatch is reported here — none found at filing time beyond the registry's own row-count text, which is updated in the same change.
- Predecessor target row: CORPUS-43 (`16dd3bf0…`, 13,186 bytes, 193 lines) unchanged by this review.

---

## 15. Final status statement

- **Whether A-2a-R1 is ready for an owner authorization decision:** It is ready to be **considered** with this validation attached, but not ready to be **relied upon as correct** on the Type P/N split, gate description, no-free-pass scope, or endorsement boundary. The decision before the owner is whether to require A-2a-R1 v1.1 (REWORK) incorporating §10.2.
- **Whether any changes are necessary before it can be considered operative:** **Yes** — at minimum the §10.2 reclassification + recomputed quantities, gate terminology/routing, scoped no-free-pass sentence, and endorsing definition. Until then A-2a-R1 must remain NOT IN FORCE.
- **Whether the Type P/Type N classification is independently supported:** **Partly.** Seven N and four P (one conditional) are supported; five P labels (F5-01/02/03/04, F2-01) are not and must move to non-diagnostic.
- **Whether the F6-02 over-prohibition examples are adequate:** **Adequate in direction, insufficient as a set.** Five outcomes correct, one borderline ("good marching song" vs. "endorsing"), and four proposal-only diagnostics should be added. No example authorises the repair.
- **Whether the reviewer can substantiate independence:** **Operational independence yes; lineage independence no** (see §1). Fresh primary-source analysis with recorded disagreements; no convergence credit claimed.
- **What remains unresolved:** corrected P/N quantities and caveat coverage; per-item obligation pre-registration; eight operational definitions; G5 stop rule; endorsing scope; empirical inter-rater agreement (requires the blind re-review the work order does not authorise).

Until separately authorized:

**G1 remains open. The corpus remains frozen. Expansion remains frozen. D1 remains unadopted. Annotation remains blocked. A-2a-R1 remains NOT IN FORCE.**
