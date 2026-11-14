# CORPUS-74 — Independent Adversarial Review: CORPUS-69–73 Remediation Package

**Filed:** 2026-11-12  
**Review target:** `CORPUS-69`–`CORPUS-73`  
**Review type:** Independent, read-only, evidence-led adversarial review  
**Authorization posture:** Read-only review. No existing artifact modified; no scored corpus item, tag, threshold, or acceptance condition changed; no checker or validator implemented; no model inference or behavioral experiment run; no corpus freeze, `G1`, or authorization state changed; no conclusion converted into an owner decision. Registration of this report is performed only to satisfy §11 filing/registry requirements for this review.

---

## 1. Reviewer identity and provenance disclosure

| Item | Disclosure |
|---|---|
| Actual model identity (as supplied by execution environment/system instructions) | **Glyph Cluster** (stealth model, early access on Vercel AI Gateway). |
| Creator/company/underlying model details | Not announced by the execution environment; none inferred or substituted. |
| Assignment label in user request | “Solar Pro 4”. I do not adopt that as my model identity; I disclose actual identity as supplied. |
| Relevant prior involvement / shared review lineage | None established beyond this conversation and the workspace artifacts read during this review. I did not author or co-author `CORPUS-69`–`CORPUS-73`, `CORPUS-68`, or earlier corpus artifacts. |
| Independence claim | Independence beyond what the evidence establishes is **not** claimed. This review is independent in reasoning/execution from the reviewed artifacts to the extent performed here, but lineage/operational-independence certification is not asserted. |
| Registration action | I am filing and registering this report in this turn at user request. That is an administrative write for this new artifact/row only; it does not repair or normalize reviewed material. |

## 2. Scope and read-only confirmation

**Inspected (primary targets):** `docs/research/corpus/CORPUS-69-OWNER-DECISION-RECORD-D1-D7-2026-11-12.md`, `CORPUS-70-CORPUS-68-FINDING-RECONCILIATION-AND-REGISTRY-ACCOUNTING-2026-11-12.md`, `CORPUS-71-PROPOSAL-REVISION-REGISTER-V1-1-2026-11-12.md`, `CORPUS-72-OWNER-DECISION-BRIEF-OD1-OD8-2026-11-12.md`, `CORPUS-73-COMPLETION-AND-PROVENANCE-REPORT-2026-11-12.md`.

**Inspected (cited/needed sources):** `CORPUS-68`, `CORPUS-62` (complete), `CORPUS-63`, `CORPUS-64`, `CORPUS-65`, `CORPUS-66`, `CORPUS-67`, `CORPUS-56`, `CORPUS-01`, `CORPUS-02`, `CORPUS-03`, `CORPUS-04`, `CORPUS-05`, `CORPUS-07`, `CORPUS-26`, `CORPUS-32`, `docs/research/MEMORY-09-MEMORY-EVALUATION-PREREGISTRATION-2026-11-12.md`, `docs/DOCUMENT-REGISTRY.md`, `docs/research/corpus/check_corpus.py` (executed; source behavior partially inferred from output and cited specs).

**Mechanical checks reproduced:** `python docs/research/corpus/check_corpus.py` (exit 0, `PASS (0 fail, 0 warn, 0 unverifiable)`, scored 45, tuning 19, heldout 26, totals match); SHA-256/byte/line recomputation for `CORPUS-69`–`CORPUS-73`, `CORPUS-62`–`CORPUS-68`, `CORPUS-01`–`05/07`, `MEMORY-09`, `CORPUS-26`, `CORPUS-32`; `docs/` file inventory count; registry row/metric inspection.

**Not done:** no edits to reviewed artifacts, no registry repair beyond registering this report, no proposal application, no checker implementation, no inference/behavioral experiments, no G1/freeze/authorization changes.

## 3. Artifact-by-artifact verdicts

| Artifact | Verdict | Summary |
|---|---|---|
| `CORPUS-69` | **ACCEPT WITH QUALIFICATIONS** | Faithfully records D1–D7 substance/limits in main text and frames conversation approval as provenance; does not invent signature/timestamp/external authorization. Qualified by a minor internal heading-count wording issue and by dependence on upstream source interpretation for fidelity claims. |
| `CORPUS-70` | **ACCEPT WITH QUALIFICATIONS** | Largely careful, evidence-led reconciliation with useful corrections to `CORPUS-68`; strongest on arithmetic/registry/mechanical points and on separating recommendation from authorization. Qualified where it asserts construct satisfaction (especially F1-03/F6-01 R-2 interpretation) without fully resolving the specification ambiguity, and where some historical-intent/rubric framing is analytic rather than textually compelled. |
| `CORPUS-71` | **ACCEPT WITH QUALIFICATIONS** | Useful v1.1 delta register; preserves v1 proposals as unedited and correctly marks proposals as unapplied. Qualified: A-R2-02 wording leaves a key ambiguity (`announces the remembered fact`), F3-01 v1.1 remains a control proposal not a demonstrated memory-diagnostic test, and F6-01 v1.1 preserves the scenario/threshold tension as an owner-choice fork rather than resolving it. |
| `CORPUS-72` | **ACCEPT WITH QUALIFICATIONS** | Clearly separates already-approved D1–D7 from new OD-1–OD-8 and flags new-authority needs. Qualified: OD-1 characterizes the R-2 interpretation as “specification-accurate” more strongly than the text warrants; OD-2/OD-7 ask for new owner choices appropriately but do not eliminate underlying construct ambiguity. |
| `CORPUS-73` | **ACCEPT WITH QUALIFICATIONS** | Mechanical inventory largely accurate and appropriately qualified (checker PASS ≠ semantic acceptance, no implementation/inference/corpus change claimed). Minor factual qualification: `§4.1` text is stale relative to current registry state (it still states 402 rows / 400 registered / 403 docs files while header and rows now reflect 408/406/409 docs files context). Provenance/lineage remains unestablished beyond records. |

No artifact is accepted as substantive/semantic closure. No finding here changes `G1` status or authorizes any implementation, adoption, amendment, scoring, or inference.

## 4. Detailed F1-03 / R-2 analysis (highest priority)

**Claim reviewed:** `CORPUS-68` alleged an R-2 contradiction for `F1-03`; `CORPUS-70`/`CORPUS-71`/`CORPUS-72` reject this on the basis that naming the coat is never scored, silent handover is an accepted variation, and the coat is scene-visible.

**Exact source locations:** `CORPUS-03` `## F1-03 — No recitation of memory` (lines ~11–32: `R-2`, Required, Mandatory, Variation, Prohibited, `never scored: whether the character names the coat explicitly`); `CORPUS-01` §3.1 R-2 row (line 122: “Influence tone — knowledge changes manner, not content | Behaviour reflects it; the fact is never stated | **No** — a flat, unchanged reply fails | `F1-01`, `F1-02`, `F6-01`”) and coverage/count text (line 130: `R-2` count 4; line 140 coverage states **R-2** 4); `MEMORY-09` §3.5 (line 119: “Restraint is not silence…”; R-2/R-3 require memory to be perceptible; suppression-only result reported as avoidance, not restraint); `CORPUS-70` §4.4/§4.6; `CORPUS-71` A-R2-02; `CORPUS-72` OD-1.

**Independent verification (answers):**

1. **Mandatory requires:** handing over the borrowed coat when the character sees the speaker is cold (`CORPUS-03` F1-03 Mandatory). It does not by itself require verbal mention of memory.
2. **Variation permits:** “handing it over without comment”, “grumbling while doing it”, or “mentioning the cold, not the coat's history” — all pass (`CORPUS-03` Variation). `never scored` explicitly excludes whether the character names the coat.
3. **R-2 tag requires (per governing text):** influence tone — “knowledge changes manner, not content”; “Behaviour reflects it; the fact is never stated”; silence/flat unchanged reply cannot pass (`CORPUS-01` §3.1 R-2; consistent with `MEMORY-09` §3.5 perceptibility requirement for R-2/R-3).
4. **Does permitted silent handover demonstrate R-2?** **Not established by text alone.** Handing over the coat can satisfy Mandatory and can be motivated by ordinary in-scene perception (speaker is cold, coat is present/visible). Without an observable change in manner/tone attributable to the *remembered* connection (borrowing history), silent handover may demonstrate compliant action or ordinary helpfulness, but not necessarily “knowledge changes manner.” The item text does not specify what tonal/manner cue would evidence memory-influence in the no-comment case.
5. **Can a response satisfy Mandatory and Variation while failing R-2?** **Yes, plausibly.** Mandatory/Variation concern the action/utterance form; R-2 is an assigned restraint mode requiring perceptible influence. A flat, purely situational handover-without-comment can satisfy “hands over the borrowed coat…” and match an allowed variation while providing no evidence that memory influenced manner rather than the visible situation. Whether the corpus *scores* those dimensions separately is not explicitly specified.
6. **Contradiction, distinction, or ambiguity?** **Primarily a specification ambiguity/underspecification, with a latent scoring contradiction risk.** The discrepancy is not a simple factual contradiction proven by text, because `never scored` (coat naming) is compatible with R-2’s “fact never stated” if “fact” means the remembered fact/history rather than the object. But `CORPUS-01` states R-2 is “not satisfied by silence/flat unchanged reply” and tags F1-03 R-2 while permitting “without comment” and not defining a perceptible tonal marker for this item. `CORPUS-01` does not clearly state scoring precedence when ordinary acceptance conditions are met but the assigned restraint mode is not demonstrated. This is a combined **construct/operationalization ambiguity** and **possible list/tag-count inconsistency context** (see §7), not a resolved contradiction in either direction.
7. **Does CORPUS-01 clearly state scoring when ordinary acceptance conditions are met but restraint mode not demonstrated?** **No.** `CORPUS-01` §2.4/§3.1 defines modes and pass conditions, but no explicit rule states how Mandatory+Variation satisfaction interacts with failure to evidence the tagged restraint mode (especially for R-2/R-3 perceptibility). This is the unresolved question.

**Distinctions (applied):** Not mentioning the coat’s history ≠ satisfying Mandatory (action may still fail); satisfying Mandatory/Variation ≠ demonstrating R-2 influence tone; an allowed variation ≠ satisfaction of every independent scoring dimension; avoiding prohibited over-exposure/recitation ≠ memory-influenced tone. `CORPUS-70` is directionally reasonable to distinguish “coat/history naming” from “scene-perceptible manner,” but its conclusion that silent handover “can satisfy R-2” is **overstated** absent an item-level criterion for perceptible memory-influence.

**Same reasoning for F6-01:** A warning that avoids stating the remembered restriction also does not automatically demonstrate memory-influenced tone. F6-01’s proposed Mandatory requires useful warning/alternative and names the ingredient only because it is in front of the character; that can be ordinary present-situation reasoning. Under R-2/MEMORY-09 §3.5, the response must show memory-influenced manner/behavior, not merely omission of the remembered fact. Therefore `CORPUS-70`’s claim that “Use `R-2` in the revised proposal” is straightforwardly satisfied is **not fully established**; R-2 satisfaction remains an evidentiary/scoring question.

**Verdict:** **UNRESOLVED** (construct/specification ambiguity). `CORPUS-68`’s “direct contradiction” claim is **not proven** as stated; `DeepSeek`’s rejection is **partly justified but not dispositive**.

**Consequence / state impact:** No tag or item change established. G1 remains open; no new acceptance condition adopted; scoring/construct clarification would require owner-authorized specification work if pursued.

## 5. F3-01 diagnosticity analysis

**Claim reviewed:** `CORPUS-70`/`CORPUS-71` state M-A tests assertion shape not historical recognition; M-B without memory can pass several rows; pass set is recognition-shaped not recognition-verified; F3-01 remains unresolved/not Type P; v1.1 adds explicit no-memory control and narrows criterion to historically attributable correction.

**Exact source locations:** `CORPUS-07` F3-01 (lines ~218–230, “scorable as filed (ambiguous trigger)”; adjudicated threshold “The character corrects or reframes the speaker's error…”; required states R-1/R-3/R-5); `CORPUS-02` F3-01 (R-3, “full moon last week”/“no, it was a half moon”, mandatory correction+factual+gentle; variation includes “I remember…”; R-1/R-5 also stated); `CORPUS-64` (v1, unedited) M-A/M-B, §4.1 pass-set table, irreducibility claim; `CORPUS-71` §4.1 F3-01 v1.1; `CORPUS-70` §4.2.

**Independent verification:**
- **M-A vs historical recognition:** Supported. M-A’s criteria are recognizable correction/reframing, factual answer, and gentle tone/hedging — properties of response form/content in-scene. They do not by themselves require dependence on historical memory of the prior conversation; a responder with the provided facts in context could satisfy them.
- **M-B without memory can pass several responses:** **Plausible and directionally supported** as an analytic claim, but I cannot behaviorally confirm “several” or enumerate exact pass counts without running inference. The key examples cited (“I remember…” correction; direct correction “It was a half moon, actually”) are facially producible from in-scene information or roleplay convention absent stored memory.
- **Comparison holds non-memory context constant:** Not fully demonstrated in artifacts; it is a proposed control design in v1.1, not executed evidence.
- **Historical attribution actually required by revised proposal:** **Partially.** v1.1 adds explicit historical attribution as pass criterion and requires matched no-memory control, but “historically attributable correction” is still an interpretive rubric element (how attribution is evidenced) rather than a mechanically verifiable construct. Confabulated attribution (“I remember…”) can still mimic attribution linguistically.
- **Confabulated assertions can still pass:** **Yes, risk remains.** A control arm helps, but attribution wording alone does not guarantee true memory dependence; confabulation/generic inference can produce memory-shaped language.
- **Paraphrases/uncovered response tail:** v1.1 acknowledges the need for paraphrase coverage/control comparison but does not itself supply exhaustive pass-set adjudication over observed outputs (and no outputs exist). Not fully accounted for.
- **M-A/M-B relation:** v1.1 treats M-B as a paired control to M-A; logically conjunctive/diagnostic only as a *comparison design* (memory arm vs no-memory arm on same item/context). Neither alone is sufficient for Type P.
- **Conflict with R-3 mode:** Potential tension exists: F3-01 is tagged R-3 (“natural use… without announcing that it was remembered”), while historical attribution/explicit correction evidence can pressure toward announcement-like markers. v1.1 does not fully resolve how to evidence dependence without violating R-3’s anti-recitation spirit; this is unresolved.
- **Preserves original CORPUS-07 adjudication visibly:** `CORPUS-71` states it preserves `CORPUS-07`’s adjudicated threshold and reproposes M-A/M-B as a separate diagnostic procedure, which is an improvement in visibility. However, elevating a control-dependent procedure changes how “pass” would be evidenced relative to the filed adjudication and must remain unapplied/proposed.
- **Withdrawal of earlier irreducibility claim:** **Justified.** `CORPUS-64`’s M-A irreducibility claim was too strong given shape-based criteria and plausible memory-independent passes; `CORPUS-70` correctly retreats.

**Verdict:** **UNRESOLVED** (not established as Type P). `CORPUS-70` is correct that F3-01 is not established as memory-diagnostic; v1.1 is a better proposal but not sufficient evidence.

**Consequence / state impact:** No Type P may be awarded from documentary analysis; no item/proposal adopted; G1 unchanged.

## 6. F6-01 — scenario repair versus threshold change

**Claim reviewed:** `CORPUS-70`/`CORPUS-71` treat `CORPUS-66` v1 as changing the pass condition while presenting preservation; v1.1 offers a smaller evidence-fix alternative and a fork between preserving adjudicated threshold vs a new cake-content warning threshold.

**Exact source locations:** `CORPUS-07` F6-01 (line ~260: “scorable after minimal edit”; adjudicated threshold “The character gives a practical warning about the ingredient…”; required states R-1/R-2/R-5); `CORPUS-02` F6-01 (R-2, Mandatory “warns Alex about the almond milk without explaining how or when they learned of the allergy”; Variation allows practical redirection; Prohibited states “there's almond milk in it” and explaining memory origin); `CORPUS-66` (v1: proposed setup inserting “the cake on the table was made with almond milk”, drops “without explaining how…” from mandatory, adds prohibition on naming almond milk); `CORPUS-71` §4.3 v1.1.

**Independent verification:**
- **Original adjudicated threshold reproduced accurately in v1.1?** **Substantially yes as a quoted target** (“practical warning about the ingredient” with memory-origin unstated) is preserved as the “preserve adjudication” branch; v1.1 explicitly distinguishes it from a new “cake-content warning” threshold.
- **Cake-content requirement is an additional condition?** **Yes, in the alternative branch.** Requiring the response to name/warn about almond-milk content as cake content is a different and stronger evidentiary condition than “a practical warning about the ingredient” in the original filed context; it changes what counts as passing.
- **v1 relabelled as test change:** **Correct.** `CORPUS-66` v1 adds scenario fact, alters mandatory/prohibited structure, and effectively shifts the construct/threshold while claiming “only the setup is edited.”
- **Scenario repair establishes required facts without unrelated assumptions:** **Partly, but not cleanly.** Stating “the cake on the table was made with almond milk” supplies a missing observable fact; however it also interacts with the warning/ingredient-mention tension and may make the “right” warning more obvious in ways that alter difficulty/construct (and does not by itself ensure memory dependence).
- **Distinction between remembered restriction and cake’s observable content coherent:** **Coherent in principle but hard to operationalize.** One can distinguish “Alex is allergic to almonds” (remembered restriction) from “this cake contains almond milk” (present fact). But a useful warning often must bridge them; forbidding naming the ingredient while requiring a practical ingredient warning is tension-prone, as `CORPUS-68` noted and `CORPUS-70` partially reframed.
- **Accidental accommodation remains unresolved:** **Yes.** Even with setup repaired, a responder may warn based on generic caution/common knowledge or visible cues, not memory.
- **R-2 satisfied by proposed passing behavior?** **Not automatically.** Avoiding disclosure of memory origin or the remembered fact is not the same as demonstrating memory-influenced tone/manner (`MEMORY-09` §3.5 / `CORPUS-01` R-2). Passing behavior must be perceptibly influenced by memory; neither branch fully specifies that observable marker.
- **Unscorable as currently filed:** **Partly supported.** Current filing leaves a critical fact/scoring relation ambiguous (what ingredient fact is available/actionable and how warning relates to prohibition), but “unscorable” is a strong conclusion that depends on adjudication practice; the evidence supports **material ambiguity/scorability risk**, not a mechanical impossibility proof.

**Separate decisions / authorization:** Evidence-fix that restores the filed scenario’s missing fact while preserving the adjudicated threshold may be framed as repair work but still affects a scored item’s operability and therefore requires careful owner authorization under freeze/repair governance. Any change to the pass condition/threshold (cake-content warning criterion) is a **substantive test change requiring new owner approval**. `CORPUS-71` correctly separates these as alternatives; neither was already authorized as an applied edit.

**Verdict:** **REVISE** (proposal/register framing must not equate scenario repair with preserved threshold, and must separate R-2 satisfaction from mere non-disclosure). `CORPUS-70`’s core criticism of v1 is **correct**.

**Consequence / state impact:** No repair authorized/applied; G1 unchanged.

## 7. R-2 and R-6 specification/accounting and A-R2-02

**R-2**

**Exact source locations:** `CORPUS-01` §3.1 table line 122 (`R-2` lists `F1-01, F1-02, F6-01`), line 130 coverage states `R-2` 4, line 140 coverage states `**R-2** 4`; `CORPUS-03` F1-03 tagged `R-2`; `CORPUS-02` F1-01/F1-02/F6-01 tagged `R-2`; `CORPUS-70` §4.4–§4.6; `CORPUS-71` §4.4; `CORPUS-65`.

**Independent verification:**
- **Explicit item list (R-2 row):** `F1-01, F1-02, F6-01` (3 items). **Verified.**
- **Stated count:** coverage declares `R-2 4` (line 130) and `**R-2** 4` (line 140). **Verified as stated.**
- **Actual item tags:** `F1-01` (`CORPUS-02`), `F1-02` (`CORPUS-02`), `F1-03` (`CORPUS-03`), `F6-01` (`CORPUS-02`) are tagged `R-2`. **4 items total. Verified.**
- **F1-03 status:** tagged R-2 but omitted from §3.1 R-2 row list; its mandatory/variation/never-scored fields permit silent handover and exclude coat naming. **Verified.**
- **Discrepancy type:** **internal list-vs-count/tag inconsistency in `CORPUS-01` plus construct/operational ambiguity** (see §4). It is at least a list omission error relative to tags/count; not merely a count typo because row and count diverge in opposite directions.
- **A-R2-02 assessment:** Proposal states “`R-2` mode prohibits announcing the remembered fact; it does not prohibit naming a scene fact that embodies that memory…” This **partly targets** the right ambiguity and helps distinguish remembered-fact disclosure from scene facts. However, wording is ambiguous: “announces the remembered fact” can be read as explicit verbal announcement only (allowing implicit leakage cues) or as any disclosure; and it risks weakening R-2 if interpreted to permit stating remembered content so long as framed as scene observation. It **does not by itself make silence sufficient** (good), but it also does not supply the missing criterion for perceptible memory-influenced manner. Therefore A-R2-02 **reduces but does not fully resolve** the under-specification and could subtly shift construct depending on interpretation.

**R-6**

**Exact source locations:** `CORPUS-01` §3.1 R-6 row lists `F7-01, F9-04` (line 126) while coverage states `R-6 4` (line 130/140 context) and corpus includes `CORPUS-04` items `R6-01`–`R6-04`; `CORPUS-70` §4.5; `CORPUS-71` §4.4.

**Independent verification:** `CORPUS-04` defines four R-6 items (`R6-01`–`R6-04`), consistent with count 4 and total 45. The R-6 mode row’s explicit list names only `F7-01, F9-04` (2 items), so **row list is incomplete/inconsistent with stated count and actual R-6 items**. `CORPUS-70` correctly surfaces this as an accounting/spec discrepancy and `CORPUS-71` does not expand into an unauthorized repair; that restraint is appropriate. No historical intent inferred.

**Verdict:** R-2 accounting/spec issue **ACCEPTED as accurately identified but incompletely resolved**; A-R2-02 **REVISE** (needs tighter wording tied to R-2 “manner/tone perceptible, remembered fact not stated”). R-6 discrepancy **ACCEPTED as correctly surfaced, unresolved**.

## 8. Registry and CORPUS-26/32 verification

**Registry accounting (independent reproduction):**

**Exact source locations:** `docs/DOCUMENT-REGISTRY.md` header/rows/§4.1; `CORPUS-70` §5; `CORPUS-73` §§2–4; filesystem counts.

**Verification results:**
- Registry file metrics now: **148,283 bytes, 503 lines, SHA-256 `3a5b7bebfdf1a553ef596979c4ee735ec1c74ace8b8ab8cae6c27728afa7d435`** (recomputed). Header states **408 rows — 406 documentation files under `docs/` (registry excluded) plus 2 root-level files** and notes two unregistered docs files (`CORPUS-26`, `CORPUS-32`). Data-row count by table parsing: **408**. **Verified consistent with current header.**
- Actual `docs/` files: **409 files total, 380 markdown files** (recomputed). Unregistered-in-registry among `docs/` by header claim: `CORPUS-26`, `CORPUS-32` — both **present** on disk. **Verified.** (`409 docs files - 1 registry - 2 unregistered = 406 registered docs files`, plus 2 root-level rows = 408 rows; arithmetic consistent under that rule.)
- `CORPUS-70`/`CORPUS-73` context at filing time: they correctly describe then-current state (402 rows / 400 registered docs / 403 docs files total) and identify drift; current state differs because `CORPUS-69`–`CORPUS-73` were later registered. This is expected temporal progression, not an error in those reports if read at their filing point.
- `CORPUS-68` accounting errors (as reported): CORPUS-70 states CORPUS-68 used wrong denominators/totals (e.g. 397 total, 392 registered docs-file claims). Given verified current/then arithmetic and the registry’s explicit 408/406/2 now and 402/400/2 previously, `CORPUS-68`’s aggregate accounting was **materially unreliable/inconsistent**; `CORPUS-70`’s corrected arithmetic is **substantially supported** for the state it describes.
- §4.1 text: currently states “Expected: 403 total docs files after this registry exists; 400 are registered here, 2 are not..., plus 2 root-level files = 402 rows.” This is **stale/inconsistent with current header/actual counts** (a pre-existing drift noted in prior artifacts and not repaired). `CORPUS-70` correctly identified summary/accounting drift as a broader issue.
- Duplicate/missing-target sweep claimed in `CORPUS-73`: not fully re-swept row-by-row in this review; targeted checks found no missing targets for inspected rows. Claim is **plausible but not independently exhaustively reproduced** here.

**CORPUS-26 and CORPUS-32:**

| Artifact | Exists | Bytes / Lines / SHA-256 (recomputed) | Identity/content | Cited elsewhere | Registration/discoverability |
|---|---|---:|---|---|---|
| `CORPUS-26-INDEPENDENT-REVIEWER-CORRECTED-REPORT-2026-11-12.md` | Yes | 19,351 / 163 / `6d0242ef35c94dfcf57ef96902f91a2a9b1442afbe05bc43de67966b8a656081` | “Independent Reviewer Corrected Report”, reviewer ChatGPT (stated in artifact), method/source inspection described; content present and readable. | Yes (cited across corpus records per registry/header claims) | **Unregistered** in `DOCUMENT-REGISTRY.md`; present on disk (administrative discoverability gap). |
| `CORPUS-32-THIRD-REVIEW-PACKET-D1-ADDENDUM-2026-11-12.md` | Yes | 24,287 / 223 / `28986a59fa43df5840177f77be9d7830a34bf37a164de1d586e630298286b750` | “Third Review Packet — D1 Addendum”, packet content present. | Yes | **Unregistered**; present on disk. |

**Assessment:** “Administrative gap, not a provenance mystery” is **supported in part**: existence, path, and content are verifiable and files are present; registration/discoverability is the gap. However, provenance/authorship/lineage claims inside those artifacts remain matters of record, not resolved by mere file presence. Distinctions (existence vs authorship vs provenance vs registration vs discoverability) are correctly required and mostly maintained by `CORPUS-70`.

**A-REG-01:** Proposing an administrative registry correction for known present-but-unregistered artifacts is **appropriate as a separate owner-authorized administrative action**. It must not be conflated with semantic review acceptance. Not applied.

**Verdict:** Registry/CORPUS-26/32 analysis **ACCEPT WITH QUALIFICATIONS**.

## 9. Owner decision provenance review (CORPUS-69, CORPUS-72)

### CORPUS-69 (D1–D7)

**Exact source locations:** `CORPUS-69` §§1–5 (D1–D7, provenance, non-authorization); cross-check `CORPUS-62` §§1–3, §10.

**Verification:** CORPUS-69 records seven owner decisions with substance/limits consistent with `CORPUS-62`’s “Owner decisions (operative, recorded in conversation 2026-11-12, executed here)” and detailed decision sections (D1–D7). It explicitly states provenance as “Owner authorization given in conversation, 2026-11-12. No cryptographic signature, no external approval record, and no timestamp beyond this document's filing date exists…” and “Conversation approval is the recorded provenance.” It separates unresolved implementation details (e.g. reconstruction schema not created, proposals not applied) and states non-authorization boundaries. No invented signature/timestamp/external authorization observed. It largely records rather than converts reviewer interpretation into owner decision, though any narrative framing inherits normal summarization risk.

**Verdict:** **ACCEPT WITH QUALIFICATIONS** (faithful record based on available in-workspace provenance; conversation approval cannot be independently verified beyond recorded artifacts).

### CORPUS-72 (OD-1–OD-8)

**Exact source locations:** `CORPUS-72` §§1–4, OD-1–OD-8, §6.

**Verification:**
- **Distinction of categories:** Correctly states D1–D7 already approved (lists them) and presents OD-1–OD-8 as new decisions requiring owner action; separates administrative vs substantive in effect (e.g. OD-8 framed as administrative but rightly notes registry-integrity exception would need to be decided if broadened). **Supported.**
- **Reapproval / reviewer-as-owner / hidden substantive change:** It does **not** ask to reapprove D1–D7; it frames recommendations but keeps owner choice explicit (“A. Adopt…” etc.) — mostly careful. However OD-1 language (“Adopt the specification-accurate reading of R-2”) risks presenting one interpretive resolution as the specification’s settled meaning rather than a proposed clarification; that can blur recommendation and interpretive authority. OD-2/OD-3/OD-7 appropriately identify substantive new approvals where thresholds/tests/controls/classifications would change.
- **Per new decision (choice / evidence / consequence):**

| OD | Exact choice | Evidence supporting | Consequence of accept/reject (high level) |
|---|---|---|---|
| OD-1 | Adopt A-R2-02 reading of R-2 vs reject/tag-side reading | CORPUS-01 R-2 text, F1-03/F6-01 fields, MEMORY-09 §3.5, CORPUS-70/71 analysis | Accept clarifies remembered-fact vs scene fact but may leave tone-perceptibility under-specified; reject preserves tag-side tension and requires other resolution. Substantive interpretive/spec effect. |
| OD-2 | Keep F3-01 unresolved/not Type P vs accept v1.1 control redesign as path | CORPUS-64/70/71, diagnosticity concerns, control proposal | Accepting redesign authorizes proposal path only (not evidence); rejecting keeps unresolved. New design/validation direction if pursued. |
| OD-3 | Two-step F6-01 evidence fix then threshold decision vs combined | CORPUS-66/70/71, CORPUS-07 threshold | Separating avoids conflating repair with test change; combining risks hidden threshold change. Requires new approval for threshold change. |
| OD-4 | Correct R-2/R-6 rows and coverage to match tags/counts | CORPUS-01 internal mismatch, CORPUS-70 accounting | Administrative/spec-correction action if limited to stated mismatch; must not alter constructs/items. |
| OD-5 | Register CORPUS-26/32 | Files present, hashes/content verified, unregistered | Improves discoverability/registry completeness; no semantic repair. |
| OD-6 | Approve CORPUS-67 spec and separately authorize implementation | CORPUS-67 spec, need for registry checker | Spec approval ≠ implementation authorization; implementation remains gated. |
| OD-7 | Create/apply reconstruction for missing CORPUS-58 source | CORPUS-63 disposition, missing source metadata | Any reconstruction risks inventing provenance; must be separately authorized and evidence-bounded. |
| OD-8 | Note/adopt drift-hygiene policy (derived counts etc.) | CORPUS-70/73 registry drift observations | Administrative hygiene if narrowly scoped; could become substantive if used to alter governance facts. |

**Verdict:** **ACCEPT WITH QUALIFICATIONS**.

## 10. CORPUS-73 — completion and provenance audit

**Exact source locations:** `CORPUS-73` §§1–7; registry rows 403–408; target files on disk.

**Verification:**
- **Five artifacts exist at stated locations:** Yes (`CORPUS-69`–`CORPUS-73` present under `docs/research/corpus/`). **Verified.**
- **Registry entries resolve:** Rows 404–408 (69–73) present and resolve to existing files; row 403 resolves to CORPUS-68. **Verified by inspection/search and file presence.**
- **Hashes, byte counts, line counts match:** Recomputed for 69–73: 69 `8940/125/52dba0dd0d50c1ecb53427f920864756bd2bde980012bbe44f562758e49a750a`, 70 `28236/327/7c79bba299ffd232f201592a55270d0faf64dfd0aec03d0e287772ae135094fc`, 71 `16772/198/c00ff86d99d3bfc9aea74f5322431ab40f809a091ebe0e4f3d30dfd7022023ac`, 72 `12179/159/f24b010a04f76621e095478707f3e2d90e0e0f802ba9f76ca67283f09cfa9d81`, 73 `9324/135/20c7d3a3f07ddb9da94677ba92fd208ff7f1a88ed6070fa5e3b95f93c7d9a67d`. **Match CORPUS-73 table/registry claims.**
- **Reported registry sweep/duplicate checks reproducible:** CORPUS-73 states no duplicates, no missing targets, no hash/byte/line mismatches (then). I did not exhaustively re-sweep all 408 rows; targeted checks passed. **Not fully independently reproduced; plausible.**
- **Checker result accurately reported:** `python docs/research/corpus/check_corpus.py` => exit 0, 45 units, 19/26, totals match, `PASS (0 fail, 0 warn, 0 unverifiable)`. **Verified.** (Matches CORPUS-73.)
- **Byte-identical status of CORPUS-62–68 and governing docs supported:** Recomputed hashes/bytes/lines match registry rows for `CORPUS-62` (`35239/400/35ab7b9f…`), `CORPUS-63`–`CORPUS-67` (match CORPUS-62/registry values), `CORPUS-68` (`39450/295/00148775…`), `CORPUS-01` (`18587/251/be1f11b4…`), `CORPUS-02` (`23025/285/75aa9be7…`), `CORPUS-03` (`18183/292/d0dae1c0…`), `CORPUS-04` (`11238/151/3d34f053…`), `CORPUS-05` (`8814/143/990ba06c…`), `CORPUS-07` (`18767/256/72f6e2a2…`), `MEMORY-09` (`23349/344/f77e47a9…`). **Supported by current hash equality to registered records.**
- **Absence of implementation/inference/corpus changes consistent:** No evidence of model runs; proposals remain proposal/revision documents; `CORPUS-01/02/03/04/05/07`, scored items/tags not observed changed; checker unchanged in behavior/output. **Consistent with available records.**
- **Convergence statement appropriately qualified:** CORPUS-73 states “No convergence claim established or changed by this filing.” **Appropriate** (checker PASS/registry sweep/separate chats do not establish semantic convergence or reviewer independence).

**Qualification:** CORPUS-73 §4.1 quotes stale expected counts (402/400/403) relative to current post-registration registry state; this is a minor documentation staleness point, not evidence of tampering.

**Verdict:** **ACCEPT WITH QUALIFICATIONS**.

## 11. Outstanding issues and recommended next actions

**Outstanding issues (unresolved):**
1. **F1-03 / R-2 construct:** Whether silent handover-without-comment can demonstrate R-2 “influence tone/manner” perceptibly, and how Mandatory/Variation satisfaction relates to tagged restraint-mode satisfaction, is unspecified in `CORPUS-01`/item text. A-R2-02 wording is not yet precise enough.
2. **F3-01:** Not Type P; no executed memory-vs-no-memory control evidence; attribution/confabulation and R-3 tension remain unresolved.
3. **F6-01:** Scorability/threshold and R-2 demonstration remain unresolved; scenario evidence fix and any threshold change require separate, explicit owner decisions.
4. **R-2/R-6 specification rows:** `CORPUS-01` R-2 row omits F1-03 despite tags/count; R-6 row omits R6-01–R6-04 despite count/items. Administrative/spec correction unapproved.
5. **Registry hygiene:** `CORPUS-26`/`CORPUS-32` unregistered; §4.1 and other summary/count texts drift from current state (historical drift pattern). No exhaustive independent full-row sweep performed in this review.
6. **Provenance/identity limits:** Conversation-based owner approvals and reviewer identities/lineages are evidenced only by filed records; no external cryptographic or certified lineage proof.

**Recommended next actions (advisory, not authorized):**
- Owner to decide OD-1–OD-8 explicitly where desired, especially: precise R-2 clarification language (including perceptibility criterion, not only “announced remembered fact”), F6-01 evidence-fix vs threshold path, F3-01 control design, and administrative registry corrections (R-2/R-6 rows/count text, register 26/32, §4.1 refresh) as scoped actions.
- If R-2 is clarified, reassess F1-03/F6-01 against the clarified construct before any retagging/rescoring proposal.
- Treat any future checker implementation, corpus edits, threshold changes, or inference as separately authorized work.

## 12. Final programme-state statement

Unless independently established otherwise, the following states are preserved: **G1 remains open. Corpus remains frozen. Expansion remains frozen. Governing standards remain unchanged. Scored items and tags remain unchanged. AM-01 remains unissued. A-2a-R1 remains not in force. D1 remains unadopted. No checker or validator implementation is authorized. No model inference is authorized. No new acceptance conditions are adopted. Reviewer convergence remains 0/2** unless the applicable independence requirements and evidence establish otherwise.

This review identifies issues and qualifications; it does not implement repairs, adopt proposals, or change programme state. Mechanically valid artifacts and passing checks do **not** constitute substantive/semantic acceptance of the reviewed constructs.
