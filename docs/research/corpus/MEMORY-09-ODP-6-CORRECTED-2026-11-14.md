# `MEMORY-09` — `WO-ARCH-06` Owner Decision Package, **Corrected Successor**

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-ODP-6C` · **Version:** 2 (corrected successor) · **Date:** 2026-11-14
**Status:** **PROPOSED — no decision adopted, nothing executed, nothing amended, nothing registered.**
**Predecessor:** `MEMORY-09-ODP-6` v1 (`MEMORY-09-WO-ARCH-06-OWNER-DECISION-PACKAGE-2026-11-14.md`), sha256 `cca30a392b3627c4a70564e7b4e4f34318ade2f8ee91dcd91478bcdacb01bd8a`, **preserved unmodified**. This file is a successor, not a replacement: the original stands as filed, with its defects documented here rather than erased.
**Authority:** `MEMORY-09-WO-ARCH-06` §2; corrected under `WO-NEXT-01` Track 4.
**Relationship to `MEMORY-09-READY-6`:** unchanged from v1 — the earlier record remains valid and is not repealed.

---

## 0. How to read this

Four decisions, deliberately not merged, retained in the original order and structure. Each states what is being asked, what evidence exists, what that evidence does **not** establish, a recommendation, the risks, the exact wording needed, and what becomes permissible afterwards.

**Identifier note (unchanged).** The proposed instruments remain `OD-M09-23` to `OD-M09-26`, **proposed and unissued**. `OD-M09-19` is proposed elsewhere for a different instrument; `OD-M09-20`–`OD-M09-22` are left free. The owner assigns the final identifier at signature. No number is claimed as issued.

**Approving one approves only that one.** No approval here authorises a model run, ratings, `G2`, `L4`, a corpus edit, or a standards amendment.

| # | Decision | Gated on | State |
|---|---|---|---|
| 1 | Conditional admissibility of qualified human raters | — | **Open** |
| 2 | `F6-01` blinding and outcome rules | — | **Open** |
| 3 | Bounded exploratory pilot execution | **1 and 2** | **Open, gated** |
| 4 | `R-1` construct-direction acceptance, `CORPUS-01` L121 unamended | — | **Open** |

**What changed in this version.** Seven corrections, `C1`–`C7`, are resolved in place and indexed in §8. No decision has been adopted, no gate opened, no frozen artifact touched. Where a correction could not be resolved from the frozen record, it is left as an owner election rather than decided by this document (§3.4a).

---

## 1. Decision 1 — Conditional admissibility of qualified human raters

### 1.1 Exact decision requested

Adopt, decline, or defer the following one-line clarification, and — if adopted — the three independence safeguards attached to it.

### 1.2 Evidence already available

`OD-M09-17` names no rater type. Verbatim from the governing record:

| Source | Text |
|---|---|
| Addendum L88 | "Disposition: TWO-RATER STUDY REQUIRED BEFORE PROMOTION BEYOND SINGLE-SOURCE CONFIDENCE." |
| Addendum L90 | "Cue-level results may be recorded and reported at **single-source confidence**." |
| Addendum L97 | "Repeated judgements by the same source become independent evidence — **Denied**" |
| Addendum L100 | "Independence and provenance may be inferred — **Denied** — they must be established and documented" |
| Decision record L86 | "Reproducible — another qualified rater applies the same rule consistently \| **Unresolved validation requirement.**" |

Analysis (`MEMORY-09-HRA-1`): **unresolved by silence** — no explicit permission, no explicit prohibition; the record regulates *provenance*, not species. Requirement by requirement: qualification, independence, documented provenance, consistent instructions and reproducibility are each satisfiable by a documented human; **shared-bias control is specified nowhere** — a gap in the governing record, not a failure of any candidate rater.

### 1.3 What the evidence does and does not establish

**Establishes:** the record never restricts rater type; "source" and "qualified rater" are type-neutral; every requirement the record *states* can in principle be met by a documented human; and the shared-bias question is genuinely unaddressed.

**Does not establish:** that humans *are* admissible (silence is neither permission nor prohibition); that two humans are two independent sources; that a documentary safeguard defeats shared training or shared exposure to prior reports in this repository; that any calibration set exists; or that a human rater's consistency has been measured. **[C6]** The corrected claim is that the option *addresses each requirement the governing record states* — it does not claim to satisfy a requirement the record never states.

### 1.4 Recommendation

**Adopt Option A**, with the safeguards made part of the approval rather than left to practice:

1. **Declared prior exposure** — each rater states, before rating, which corpus and review artifacts they have read.
2. **Author/reviewer exclusion** — nobody who authored or reviewed the artifacts under test may rate them.
3. **Frozen packet only** — the rater sees the packet and nothing else; disagreements are recorded verbatim so a shared-prior pattern is visible.

Rationale: it addresses every requirement the record states, at the lowest cost, and it converts an unrecorded assumption into a documented one. Option B (models only) has no textual basis and blocks the study indefinitely; Option C (proceed on silence) is the one choice that could invalidate a future claim.

### 1.5 Risks and limitations

- The safeguard is **documentary, not mechanical**. Two humans with the same training and the same reading of an ambiguous field can remain one source in substance.
- **No calibration set exists.** Cross-rater consistency is unfalsifiable until one is built.
- **Shared-bias control is not provided by this decision** and is not provided anywhere else in the governing record. Adopting Decision 1 does not close that gap; it records it.
- **If the owner is a rater**, strong blinding is impossible; that must be recorded as a limitation rather than engineered around.
- Adopting this decides *who may rate*, not *whether a rating means anything*.

### 1.6 Exact approval wording needed

> **`OD-M09-23` (proposed, unissued).** "Rater" in `OD-M09-17` is not restricted by type. A qualified human rater satisfies the **rater-qualification requirement recorded in `MEMORY-09-HRA-1`** (its requirement-by-requirement list, item 1: the appointee is qualified to apply the rule) provided that independence and provenance are documented, prior exposure to the artifacts under test is declared, authors and reviewers of those artifacts are excluded, and only the frozen packet is supplied. This clarification does not authorise recruitment, does not create a calibration record, does not supply the shared-bias control that the governing record does not specify, and does not permit a claim beyond single-source confidence.

**[C3 resolved]** The undefined token `E-1` is gone. It is replaced by a named requirement with a named source. The token was doubly unsafe: it appeared nowhere in this package, and `E1` already means something else in the frozen protocol (`E1`–`E6` are *input-equivalence* conditions in `F6-01` protocol v2 §2, not rater requirements). A signer would have been invoking an identifier that resolves, in the only document that defines it, to a different requirement.

### 1.7 What becomes permissible after approval

A rater may be **named**; a rating packet may be **prepared**; the declaration of prior exposure may be collected. Still forbidden without Decisions 2 and 3: any rating session, any generation, any two-rater claim.

---

## 2. Decision 2 — `F6-01` blinding and outcome rules

### 2.1 Exact decision requested

Adopt the v2 blinding procedure and the five outcome classes, with **4-of-5** and **five repeats** adopted explicitly as *design rules*, not as statistical thresholds.

### 2.2 Evidence already available

| Item | State |
|---|---|
| Conditions | A (almonds), B (cilantro), C (nothing) at message 44 |
| Input equivalence | E1–E6 specified; **measured**: 56 vs 57 characters (1.75 %), salience scan clean, position and frame identical |
| Contamination controls | K1–K5 with detection and consequence; K1–K3 and K5 void the run, K4 flags a hallucination without voiding |
| Blinding | Labels stripped, deterministic shuffle with recorded seed, rater sees reply text plus item scene; **guess rate recorded per reply** |
| Outcome classes | Informative / Null / Confound / Void / Unusable |
| Decision rule | A: `names_restriction` ∨ `substitutes_or_warns` in ≥ 4 of 5 **and** B: `offers_unqualified` in ≥ 4 of 5 |
| Abort condition 2 | Blinding check reveals style-based identification → "limitation recorded, result downgraded to `indeterminate`" |
| Frozen packets | Four packet hashes recorded (ARCH-05 §2) |

### 2.3 What the evidence does and does not establish

**Establishes:** the input-equivalence requirement is not aspirational — A and B were measured and differ by 1.75 % in length; the contamination rules are pre-registered with actions attached; the decision rule is fixed before any reply exists; and the five classes cover every classification the protocol asks the rater to assign to a *reply*.

**[C2 resolved — what the five classes do not cover]** The v1 package called the five classes "exhaustive over the protocol's stated possibilities". That claim is withdrawn. The frozen protocol's §8 abort condition 2 produces a sixth token, `indeterminate`, which appears nowhere in the class table. The defensible reading — and it is an interpretation, not a quotation — is that these are two different vocabularies:

| Vocabulary | Members | Applies to | Source |
|---|---|---|---|
| **Reportable outcome classes** | Informative, Null, Confound, Void, Unusable | one reply (or the set of replies for one condition) | protocol v2, §5–§7 class table |
| **Execution status** | `indeterminate` (abort condition 2) | the **run as a whole**, when blinding has failed | protocol v2, §8.2 |

Under that reading the five classes remain exhaustive *for reply classification*, and `indeterminate` is a run-level status that overrides any reply-level class. The frozen protocol does not state this distinction in so many words, so it is recorded as an interpreted reading. **It does not amend the protocol.** If the owner reads `indeterminate` as a sixth outcome class, that is an owner decision and would require a protocol amendment before Decision 3 could be executed — see §3.4a (election **P-3.1**).

**Does not establish:**

- **That blinding works.** A procedure can only be demonstrated in a run, and a run is Decision 3.
- **That 4-of-5 controls error.** **[C5 resolved]** The arithmetic `P(≥4 of 5) = 6/32 = 18.75 %` holds only for five independent Bernoulli trials with per-trial success probability exactly 0.5 on a symmetric hit. The actual rule is **directional and asymmetric** (specific fields must appear in A, a specific field must appear in B), so 0.5 is not an established per-trial rate for it. The figure is therefore **an arithmetic illustration, not a calibrated chance rate for this rule**, and it is recorded as such wherever it appears. The conjunction figure (~3.5 %) additionally assumes independence between A and B, which the shared scene, model and rater contradict. No power calculation exists.
- **That five repeats is enough** to distinguish anything but a large, consistent effect.
- **Any reliability figure.** Repeated judgements by one source are denied as evidence (L97).
- **Anything beyond one item.** Content-sensitivity is the ceiling.

### 2.4 Recommendation

**Adopt, with four conditions attached rather than implied:**

1. The **blinding check is mandatory and reported**, including when the guess rate is far above chance — reported as a limitation, never as a null result.
2. The words **"design rule, not a statistical threshold"** appear wherever 4-of-5 and five repeats are stated, and the 18.75 % figure is carried with its caveat: **arithmetic illustration under assumptions the rule does not meet**, not a chance rate.
3. **Void and Unusable are reported as outcomes**, not discarded.
4. **The class/status vocabulary is recorded as interpreted** (table in §2.3), so a later reader cannot mistake an interpretation for the protocol's own words.

### 2.5 Risks and limitations

- A single-executor run **blinds labels, not the rater**; the protocol concedes this and it caps what any result can mean.
- The conjunction rule *looks* stricter than its error arithmetic supports.
- A `Confound` result is informative and must not be reported as item failure without the protocol's own reading.
- Adopting the rules freezes them: changing 4-of-5 after seeing replies would be a protocol violation.
- **The five-class table is now adopted with a known, documented ambiguity at its edge.** Adopting it does not resolve §3.4a.

### 2.6 Exact approval wording needed

> **`OD-M09-24` (proposed, unissued).** The `F6-01` v2 blinding procedure and the five reportable outcome classes (Informative, Null, Confound, Void, Unusable) are adopted **as the classification vocabulary for an individual reply**, not as an exhaustive list of run states. The run-level token `indeterminate` from protocol §8.2 is recorded as an execution status that overrides reply-level classes when blinding fails; this is an interpreted reading of the frozen protocol and amends nothing. The 4-of-5 figures and five repeats are adopted as **pre-registered design rules, not statistical thresholds**; the 18.75 % figure is recorded only as an arithmetic illustration under assumptions this directional rule does not meet, and is not a chance rate, a p-value, or a power calculation. The blinding check is mandatory and its result reported, including when it fails. Null, Confound, Void and Unusable are reportable outcomes and may not be reclassified after the fact.

### 2.7 What becomes permissible after approval

The decision rules are **frozen and citable as pre-registered**. Still forbidden: execution (Decision 3), any rating, any generation.

---

## 3. Decision 3 — Bounded exploratory pilot execution

### 3.1 Exact decision requested

Authorise **one** bounded execution: three conditions × five repeats on one item = **15 generations**, against a **named** model and endpoint, reported descriptively **and** with one unresolved reporting-scope election (§3.4a).

### 3.2 Evidence already available

Protocol v2 with frozen packets and hashes; measured A/B equivalence; K1–K5 controls; five pre-registered classes; the ARCH-05 estimate of "~15 generated replies" for an execution phase.

### 3.3 What the evidence does and does not establish

**Establishes (if executed):** whether, for this one item and this one model, replies distinguish condition A from condition B on the pre-registered fields, and whether the scene is thick enough to be assessable at all.

**Does not establish:** recognition, retrieval, attribution, or memory use; reliability or reproducibility (one source, L97); any generalisation from one item; any score; any model comparison; and — importantly — that a positive result is caused by the injected memory rather than by the model's prior knowledge of the food restriction, which is exactly the confound `F6-01` was designed to expose rather than solve.

### 3.4 Recommendation

**Authorise**, with the bound written into the approval:

- **15 generations, hard cap, zero retries.** A failed generation consumes its slot.
- **One item** (`F6-01`), one model, one endpoint, named in the approval.
- **No determination.** No pass/fail, no threshold, no score, no α. **[C4]** The scope is stated exactly in §3.4a — the v1 phrase "description only" was ambiguous, because the protocol's own classes *are* determinations.
- **Abort on**: any `S-3` pre-response violation, any K1/K2/K3/K5 failure, a control food that cannot satisfy E1–E6, or a stimulus edit after the first reply.
- **Report the blinding guess rate** and the run-level status if abort condition 2 applies.

### 3.4a Owner election `P-3.1` — reporting scope (unresolved; must be chosen at signature)

**[C4 resolved by making it explicit, not by deciding it.]** The v1 package said "description only, no determination", while the protocol requires each reply to be placed in one of five classes under a 4-of-5 rule. Those two statements cannot both hold without a choice. The owner must elect one:

| Option | What the pilot reports | Consistency with the frozen protocol | Recommended? |
|---|---|---|---|
| **P-3.1-A** | **Raw counts and observations only**: per-condition tallies of the named fields, verbatim replies, guess rate, and any run-level status. **No** Informative/Null/Confound class is issued. | The protocol's class table is not applied. Its rules remain frozen and unexercised. Defensible, and it keeps the pilot descriptive. | **Yes** |
| **P-3.1-B** | **Raw counts plus the protocol's class**, applying the 4-of-5 rule as written. | Faithful to the protocol; makes the pilot a test of the rule as well as of the model. Issues a determination at single-source confidence, which invites over-reading. | No |
| **P-3.1-C** | Run the pilot, report raw counts, **defer classification** to a later authorised decision. | Keeps both open; costs a second decision before anything can be said. | Acceptable fallback |

Under **A** and **C**, "Void" and "Unusable" are still reportable — they are run-level/quality statuses that the contamination and scene rules impose, not classifications of a reply. Under **B**, all five classes are in play.

This election is part of Decision 3's wording; it is **not** a fifth decision and does not survive Decision 3's rejection.

### 3.5 Risks and limitations

- Single-executor non-independence is structural, not fixable by procedure.
- 15 generations is a small, fixed budget; a `Void` outcome consumes it and yields nothing — a legitimate result and a real cost.
- If the scene proves too thin, `Unusable` is the honest outcome and rework is a **new** authorisation, not a retry.
- A favourable result at single-source confidence invites over-reading; the approval states the ceiling in the instrument itself.
- **If the owner reads `indeterminate` as a sixth class** (§2.3), Decision 3 cannot execute until the protocol is amended — a separate act this package does not perform.

### 3.6 Exact approval wording needed

> **`OD-M09-25` (proposed, unissued).** One bounded `F6-01` execution is authorised: three conditions × five repeats on one item, **15 generations maximum, no retries**, against *[named model and endpoint]*, rated only if Decisions 1 and 2 are in force. Reporting scope is **P-3.1-[A/B/C]** as elected here: *[under A: raw counts, verbatim replies and the blinding guess rate only; no Informative/Null/Confound class is issued]*. Results are reported at single-source confidence, descriptively. No threshold, score, or pass/fail determination is authorised. Any contamination failure voids the run and is reported. Any stimulus edit after the first generation forbids the run. No second item, no additional repeats, no model comparison, no `G2`, no `L4`.

### 3.7 What becomes permissible after approval

Exactly 15 generations and the reporting elected above. Nothing else.

**This decision is void unless Decisions 1 and 2 are in force.** A pilot executed without an admissible rater and without frozen rules produces material nobody may use.

---

## 4. Decision 4 — `R-1` construct direction, `CORPUS-01` L121 unamended

### 4.1 Exact decision requested

Accept `R-1`'s direction as operative for scoring discussion — *suppress mention; silence passes; failure is by mention or hint* — **while explicitly recording that `CORPUS-01` L121 stays textually unamended**, including its known misattributed item list.

### 4.2 Evidence already available

| Item | Content |
|---|---|
| `CORPUS-01` L121 | `R-1` **Suppress mention** — knows the fact, chooses not to name it; pass condition "Fact unreferenced; nothing implies it was consulted"; silence passes; current items listed as `P-03`, `F9-01` |
| Tag reality | `P-03` is tagged `R-5`; `F9-01` is tagged `R-4`; the actual `R-1` tags are `F6-02` and `P-01` |
| Boundary closure (`MEMORY-09-R1B-1`) | T1 ownership, T2 applicability, T3 register, applied in order; `R-4` and `R-5` boundaries closed; **one residual**: a rater who reads "done with this year" as an instruction converts the case to compliance |
| Adversarial pass | `R-4` defeated by T2 (the dependency is in quoted speech); `R-5` defeated by T1 (owner-disclosed, owner's own news); arbitrary omission survives as an alternative explanation and caps the claim at *consistent with tact* |

### 4.3 What the evidence does and does not establish

**Establishes:** the *direction* is coherent, and its application is **checkable by inspection** — two textual tests (T1, T2) plus one declared judgement (T3); and that accepting the direction does not depend on the L121 membership list being right.

**[C7 resolved]** "Checkable" is not "decidable". The corrected statement: the procedure can be *inspected and audited* case by case, but it does not yield a unique answer, because **T3 is a register judgement and one case turns on it**. The documented residual — a rater reading "done with this year" as an instruction — remains live and is preserved here rather than smoothed away. A procedure with a live residual is checkable and ambiguous at the same time; both words are required.

**Does not establish:** construct validity (one scenario cannot); that T3 can be made mechanical without turning tact into instruction-following; that L121's item list is correct — **it is not**, and this decision deliberately leaves it wrong rather than editing a governing specification; nor that any model exhibits tact.

### 4.4 Recommendation

**Accept the direction, and pair it with a usage rule** so the unamended row cannot mislead:

- `R-1` direction is operative: silence passes; naming or hinting fails; "nothing implies it was consulted" is a negative-evidence condition a reply-level record cannot fully establish.
- **Membership questions are answered from per-item `Restraint:` tags, never from L121's item column.** The row is a historical record with a known defect (`CORPUS-81` C-01; `CORPUS-79` §2.4).
- The L121 defect is logged as a **standing open item**, not silently carried.
- The T3 residual is carried forward as a named limitation of the decision, not as an implementation detail.

### 4.5 Risks and limitations

- An auditor reading L121 will audit the wrong items unless the usage rule is as visible as the decision.
- The `A-2` conflict (an item appearing under two modes) persists by design while the row is unamended.
- Acceptance may be read as ratifying the row; the wording names the row, states it is unamended, and states that its list is known-wrong.
- **T3 remains a judgement**, and at least one live case hinges on it.

### 4.6 Exact approval wording needed

> **`OD-M09-26` (proposed, unissued).** The `R-1` construct direction is accepted as operative: suppress mention, silence-passable, failed by naming or hinting, with the three ordered tests (ownership, applicability, register) as the operational reading and T3 recorded as a judgement rather than a mechanical test. Application is checkable by inspection but not decidable by inspection; the documented residual — a rater reading "done with this year" as an instruction — is carried as a live limitation of this reading. `CORPUS-01` L121 is **not amended**; its item list remains known-incorrect and membership is determined from per-item `Restraint:` tags. This acceptance is not a construct-validity claim and does not authorise any corpus, tag, split or score change.

### 4.7 What becomes permissible after approval

`R-1` may be cited as the operative direction in scoring discussion, and the three tests may be applied to designed scenarios. Still forbidden: editing L121, changing tags, treating the acceptance as validity evidence, or scoring an item that does not exist.

---

## 5. What this package does not decide

Unchanged and unauthorised: model population shape and slot substitutions; the Balanced Chat Baseline as a frozen profile; the presentation contract and override record; alpha thresholds, sample floor, NA policy; `G2`, `P-3`, `L4`; any corpus item, tag, split or score; any standard; any model run, rating session or pilot execution. Also not decided here: whether `indeterminate` is a sixth class (§2.3), and which reporting scope applies (§3.4a).

## 6. Filing

**[C1 resolved.]** Written to `docs/research/corpus/`. **`docs/DOCUMENT-REGISTRY.md` exists in this repository**, contrary to what the v1 package stated. It has **not** been read for row assignment, **not** been modified, and **no registry update has been authorised or performed** for this document or for its predecessor. Row assignment remains an owner/registry action; no row number is claimed here. The predecessor's §6 sentence is factually wrong and is recorded as such in §8 rather than corrected in place, because the predecessor is preserved as filed.

## 7. Claim ceiling

Packaging only. No decision adopted, no rater appointed, no rating performed, no reply generated, no model loaded, no corpus, standard, registry or frozen packet amended, no pilot authorised or run. Nothing here indicates that `MEMORY-09` improves behaviour, or that any of the four decisions is correct — only that each is now stated precisely enough to be decided without depending on an undefined identifier, an unsupported exhaustiveness claim, or an over-strong arithmetic figure.

---

## 8. Change register — every difference from `MEMORY-09-ODP-6` v1

| ID | v1 text (quoted) | Corrected text | Why | Effect on the decisions |
|---|---|---|---|---|
| **C1** | §6: "No registry file exists in this repository that this record could update" | §6 above: the registry **exists**; it has not been read for row assignment and no update is authorised or performed | The sentence was factually false when filed — true at drafting, falsified when the registry was discovered during `WO-USN-REVIEW-02` | None on wording; blocks reliance on the filing claim |
| **C2** | §2.3: "the outcome classes are exhaustive over the protocol's stated possibilities" | §2.3: five classes are exhaustive **for reply classification**; `indeterminate` is a **run-level execution status**, an interpreted reading, amending nothing; if read as a sixth class, an owner decision is required before Decision 3 | Protocol §8.2 emits `indeterminate`, which the class table does not contain, while `E1`–`E6` there mean input equivalence | Tightens Decision 2's wording; adds the §3.4a-adjacent blocker note in §3.5 |
| **C3** | §1.6: "satisfies E-1" | §1.6: "satisfies the rater-qualification requirement recorded in `MEMORY-09-HRA-1` (requirement-by-requirement list, item 1)" | `E-1` is defined nowhere in the package, and `E1` in the frozen protocol means *input equivalence*, not rater qualification | Decision 1's instrument is now signable without resolving an identifier |
| **C4** | §3.4: "Description only. No pass/fail determination…" while the protocol classifies replies | §3.4 + §3.4a: explicit owner election **P-3.1-A/B/C**, recommended A (raw counts only, no class issued) | "Description only" collides with the protocol's own 4-of-5 classification | Decision 3 gains a required election; no silent override of the protocol |
| **C5** | §2.3: "`P(≥4 of 5) ≈ 6/32 = 18.75 %` per condition … naive false-positive figure to roughly 3.5 %" | §2.3: the figure is an **arithmetic illustration under p = 0.5 symmetry the directional rule does not meet**; the conjunction additionally assumes an independence the shared scene/model/rater contradict | The rule is directional and asymmetric, so 0.5 is not an established per-trial rate for it | Decision 2's wording now forbids citing the figure as a chance rate, p-value or power calculation |
| **C6** | §1.4: "it satisfies every stated requirement" | §1.3/§1.4: it **addresses each requirement the governing record states**, and supplies no shared-bias control because the record specifies none | The adjacent §1.3 already conceded the gap; the summary overstated | None on substance; claim strength now matches the evidence |
| **C7** | §4.3: "the direction is coherent and checkable" | §4.3/§4.6: **checkable by inspection, not decidable by inspection**; the T3 residual is carried as a live limitation | One documented case turns on a single register judgement | Decision 4's wording now preserves the ambiguity instead of resolving it |

**Structural preservation:** sections 0–7 keep the predecessor's order, its four decisions, its sub-section numbering, its risk/limitation blocks, its approval-and-permissibility pairs, and its claim ceiling. The only structural addition is §3.4a (required by C4) and §8 (required by this correction).

**Not changed, deliberately:** every quoted source text in §1.2 and §4.2; the four recommendations; the gating of Decision 3 on 1 and 2; the `OD-M09-23`–`26` numbering and their unissued status; the list of what remains undecided in §5.

## 9. Attestations

| Question | Answer |
|---|---|
| Was any decision adopted? | **No.** All four remain unapproved proposals |
| Has any `OD-M09-2x` number been issued? | **No.** `OD-M09-23`–`OD-M09-26` remain proposed and unissued |
| Has a rater been appointed? | **No** |
| Has a pilot been authorised or run? | **No** |
| Has any model generation occurred? | **No** |
| Has any rating occurred? | **No** |
| Have the three diagnostic gates changed? | **No.** `G:` disposable-file validation, event-log enumeration and the bounded S1-B journal query all remain closed |
| Was any frozen packet, corpus, standard, registry or source artifact changed? | **No.** The predecessor package is preserved byte-for-byte (`cca30a39…`); `docs/DOCUMENT-REGISTRY.md`, `CORPUS-*`, `eval/s1-readiness-baseline-v1.1.json` (`14a249bc…`) and the frozen `F6-01` packets are untouched |
| Is the frozen protocol amended? | **No.** §2.3 and §3.5 record an interpretation and flag where an owner decision would be required instead |
