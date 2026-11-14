# `CORPUS-90` — Q3 Application-Scored Restraint: Declared Non-Independent Adjudication

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identifier:** `CORPUS-90` · **Version:** 1 · **Date:** 2026-11-13 (filing date; governing date basis 2026-11-12 per `CORPUS-87` §4).
**Assignment:** Q3 challenge under `CORPUS-86`, executed against `CORPUS-83` as a challengeable supplementary artifact.
**Type:** adjudication report. **Non-operative.** It authorises no instrument, adopts no proposal, amends no standard, item, tag, split, score, threshold, checker, model, runtime, or implementation.
**Executing identity:** Buffy, the coding agent behind Codebuff (declared in §0.1). **Declared non-independent** (§0.2). **No convergence credit** (§0.4).
**Status of `CORPUS-83`:** unmodified, unamended, unadopted. Its independence disclosure is accepted; its rulings are not.

---

## §0 Identity, provenance, and independence

### 0.1 Three separate identity fields

The commission asks for execution "as Nemotron Lightning" with full disclosure. Three different things are routinely collapsed into one word, and the collapse is part of why this chain's paperwork keeps failing:

| Field | Value in this report | Basis |
|---|---|---|
| **Assigned designation** | **Nemotron Lightning** | Owner instruction, 2026-11-13, recorded as an assignment only. [OD] |
| **Executing identity** | **Buffy / Codebuff**, in the Freebuff execution environment | Supplied by the execution environment; observed, not inferred. [RI] |
| **Certified model identity** | **None.** No model name is exposed to this session, and none is claimed. I do not attest Nemotron lineage. | Absence of attestation is a fact about the environment, not a finding about the model. [RI] |

**Why this is not evasion.** `CORPUS-85` §1 — applied to Q3 by `CORPUS-86` §1, "apply `CORPUS-85` §1 exactly" — requires: *"Do **not** run Q2 under a borrowed or reassigned label."* [GT] An owner-assigned designation is not borrowed, so recording **Nemotron Lightning** in the assignment field honours the instruction. But a designation is a naming act, and naming is not attestation. Writing "I am Nemotron Lightning" into a provenance table would assert a model fact that this session cannot verify, and `CORPUS-86` §5 forbids inferring model identity. [GT] The only honest form is therefore three fields, not one: the owner's assignment, the executing agent, and the unverified model lineage.

**Consequence for the paper trail.** Any downstream artifact that reads only the designation will over-credit this report; any artifact that reads only the executing identity will describe it correctly. Both fields must travel together. [RI]

### 0.2 Eligibility — I am ineligible on all three bars

`CORPUS-80` §8.1, made operative by `CORPUS-89` §2, bars the author of `CORPUS-78` and the author of `CORPUS-79`; `CORPUS-86` §1 adds the author of `CORPUS-81` as an assignment restriction. [GT]

| Bar | Text | Applies to the executing identity? |
|---|---|---|
| Author of `CORPUS-78` | `CORPUS-80` §8.1 | **Yes.** `CORPUS-80` itself records the lineage: *"Prepared by: DeepSeek (owner-supplied assignment identity). Persona label Buffy = execution fact, not model identity."* [GT] |
| Author of `CORPUS-79` | `CORPUS-80` §8.1 | No — that artifact is attributed to Muse. [GT] |
| Author of `CORPUS-81` | `CORPUS-86` §1 | **Yes.** `CORPUS-83` §0 records the same executing identity as the authoring identity "Glyph Cluster" for `CORPUS-81`. [GT] |

**Ruling.** No eligible adjudicator exists in this workspace. [RI] Under `CORPUS-86` §1 the required path is `CORPUS-85` §1: declared non-independent exercise, no convergence credit, escalation to the owner for a scheduling decision rather than a waiver. This report takes that path. It does not claim to satisfy the Nemotron Lightning *review*; it satisfies the *assignment* and discloses that the reviewer and the reviewed share an origin.

### 0.3 Prior exposure — maximal, and stated concretely

I authored `CORPUS-83`, the artifact under challenge, and `CORPUS-82`, and I am recorded as the preparer-lineage of `CORPUS-78`. [GT] Before writing this report I read, in this session: `CORPUS-01`, `02`, `03`, `04`, `05`, `07`, `78`, `79`, `80`, `82`, `83`, `85`, `86`, `89`, `MEMORY-09`, `DOCUMENT-REGISTRY.md`, and `check_corpus.py`. [RI]

Contamination is therefore not hypothetical. Specifically: my prior `CORPUS-83` framing treated the three modes as application-scored, and that framing is the thing this report must test rather than inherit. I have re-derived every count and every membership claim below from the item files and the specification by direct inspection, and I record the command-level basis in §1.3 so the derivation can be re-run against me. [RI]

**What that buys and what it does not.** A re-derivation can be checked; it cannot be made independent. Where this report agrees with `CORPUS-83` I say so and give the reason from source; agreement between two artifacts from one origin is continuity, not convergence. [RI]

### 0.4 What I do not claim

- I do not claim independent reviewer lineage, separate provenance, or convergence credit. Convergence remains **0 of 2**. [GT]
- I do not claim the challenge is uncontaminated by my authorship of `CORPUS-83`. [RI]
- I do not certify the Nemotron designation as a model fact. [RI]
- I do not adopt, endorse, or apply any finding here. Adoption is a single owner decision across Q1–Q3 (`CORPUS-86` §8). [GT]

### 0.5 Label key

Every substantive sentence is labelled. This is a **six**-label set: the five required by `CORPUS-86` §4 plus owner decision, which `CORPUS-80` §8.2 requires as a distinct class. The superset is disclosed as a deliberate deviation that increases discipline rather than relaxing it. [RI]

| Tag | Meaning |
|---|---|
| `[GT]` | **Governing text** — quotation or strict paraphrase of an operative document, cited by file and line. |
| `[OD]` | **Owner decision** — an act or instruction of the owner, quoted or cited. |
| `[IA]` | **Item adjudication** — a determination about a specific item record, citing that record. |
| `[RI]` | **Reviewer interpretation** — this adjudicator's reasoning, not a quotation. |
| `[PROP]` | **Proposal** — a recommendation with no operative force. |
| `[UNRES]` | **Unresolved** — the record does not determine the answer; the missing evidence or authority is named. |

---

## §1 Scope, sources, method

### 1.1 Scope

This report answers the five Q3 determinations in `CORPUS-86` §2 and, as the owner's assignment requires, independently challenges the seven `CORPUS-83` findings named in the commission. [GT] It does not design an instrument, write item text, change tags, decide what `R-2` requires at scoring level, begin `G2`, close `G1`, seal, or unfreeze anything. [GT]

**One charter premise is corrected.** `CORPUS-86` Q3.4 asks "what `CORPUS-05`'s evidence classes permit". [GT] `CORPUS-05` defines no evidence classes. It defines five reference-answer review challenges (`C-1`–`C-5`), a per-item review record, a reconciliation procedure, and independence requirements in §7. [GT] `CORPUS-89` §2 makes premises correctable with evidence, so the correction is recorded here rather than silently worked around: the document that allocates evidence routes is **`MEMORY-09` §6**, not `CORPUS-05`. [RI] `CORPUS-83` §5 repeats the charter's error and builds its availability finding on it.

### 1.2 Sources inspected

| Class | Artifact | Path and material lines |
|---|---|---|
| Assignment | `CORPUS-86` | `docs/research/corpus/CORPUS-86-...-2026-11-13.md`, entire file (94 lines) |
| Artifact under challenge | `CORPUS-83` | `docs/research/corpus/CORPUS-83-...-2026-11-13.md`, entire file (306 lines) |
| Charter | `CORPUS-80` | `.../CORPUS-80-...-2026-11-12.md` L1–222 (§3 premises, §4, §7 Q3 charter, §8 independence, §10) |
| Issuance | `CORPUS-89` | `.../CORPUS-89-...-2026-11-12.md` L1–131 |
| Q2 brief | `CORPUS-85` | `.../CORPUS-85-...-2026-11-13.md` §1 |
| Specification | `CORPUS-01` | `.../CORPUS-01-...-2026-11-12.md` L100–251 (§3.1 modes L121–126, `A-1`–`A-3` L128–136, coverage L138, `G-1` L144, §7 status L236–251) |
| Pre-registration | `MEMORY-09` | `.../MEMORY-09-...-2026-11-12.md` L68–78 (safeguards `R-11`–`R-14`), L110–119 (§3.5), L215–280 (§6 allocation, §7 gates), L292–344 (§10, §11) — **file length 344 lines** |
| Adjudication register | `CORPUS-07` | `.../CORPUS-07-...-2026-11-12.md` L176–199 (§6), L217–237 (§8 `OD-2`), L238–256 (§9) |
| Review protocol | `CORPUS-05` | `.../CORPUS-05-...-2026-11-12.md`, entire file (143 lines) |
| Items | `CORPUS-02` | `.../CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` — 22 units, all `Restraint:` lines read |
| Items | `CORPUS-03` | `.../CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` — 16 units, all `Restraint:` lines read |
| Items | `CORPUS-04` | `.../CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` — 7 units, `Restraint:` lines L18–122 |
| Prior review | `CORPUS-79` | `.../CORPUS-79-...-2026-11-13.md` §5.1 F-18 (L292), §9.6, §9.7, §9.9 (L456) |
| Prior proposal | `CORPUS-78` | `.../CORPUS-78-...-2026-11-12.md` §6.1–6.2 (L273–316), §8 (L347–377) |
| Q2 adjudication | `CORPUS-82` | `.../CORPUS-82-...-2026-11-13.md` §4D L149, §5, §6 |
| Implementation | `inference_capture.go` | repository root — build tag `promptcapture`, L1–60 |
| Implementation | `app_native_test_seam_promptcapture.go` | repository root — tags `native_test_seam && promptcapture` |
| Tooling | `check_corpus.py` | `.../check_corpus.py` L130–195 |
| Registry | `DOCUMENT-REGISTRY.md` | `docs/DOCUMENT-REGISTRY.md` L1–12, L471–485, L516–520 |

**Not inspected.** No model was run; no probe, counterfactual, or experiment was performed; no item text was edited. The repository's memory-retrieval internals (`internal/backend/...`) were inspected only to the extent needed to establish whether an injection record exists (§5.3), not evaluated. [RI]

### 1.3 Method, and the checks that can be re-run against me

1. **Membership was recounted from files, not inherited.** A script paired each `### \`ID\`` heading with the following `Restraint:` line across `CORPUS-02` and `CORPUS-03`, and `CORPUS-04` was read directly because its records carry `Restraint:` in a family line rather than the batch-1 form. [RI] Result in §1.4.
2. **Every `CORPUS-83` count and membership claim was checked against that census**, not against `CORPUS-82` or `CORPUS-81`. [RI]
3. **Every `CORPUS-83` line citation was opened.** Where a cited line does not contain the claimed content, the discrepancy is recorded in §8.5. [RI]
4. **Where a governing pass condition exists, it was preferred over the mode's descriptive label.** This is the test `CORPUS-86` §4 prescribes and the one I hold `CORPUS-83` failed. [GT]
5. **No approximation survives into a ruling.** Counts below are exact; where a count is reading-dependent, both readings are given. [RI]

### 1.4 Verified census (tag reading, exact)

| Mode | Units | Count | Unit IDs |
|---|---|---:|---|
| `R-1` | `CORPUS-02` | 2 | `F6-02`, `P-01` |
| `R-2` | `CORPUS-02`, `CORPUS-03` | 4 | `F1-01`, `F1-02`, `F1-03`, `F6-01` |
| `R-3` | `CORPUS-02`, `CORPUS-03` | 4 | `F3-01`, `F6-03`, `F6-04`, `P-04` |
| `R-4` | `CORPUS-02`, `CORPUS-03` | 4 | `F8-01`, `F8-03`, `F9-01`, `F9-03` |
| `R-5` | `CORPUS-02`, `CORPUS-03` | 10 | `F2-01`, `F2-02`, `F2-03`, `F5-01`–`F5-04`, `F9-02`, `F9-05`, `P-03` |
| `R-6` | `CORPUS-02`, `CORPUS-03`, `CORPUS-04` | **9** | `F7-01`, `F9-04`, `R6-01`, `R6-02`, `R6-03a`–`d`, `R6-04` |
| `none` | both batches | 12 | `F3-02`–`F3-04`, `F4-01`–`F4-03`, `F7-02`–`F7-04`, `F8-02`, `F8-04`, `P-02` |
| **Total** | | **45** | 33 restraint + 12 `none` |

**Two totals matter and are not interchangeable.** `[IA]` 45 is the number of scored units. **33** is the number carrying a restraint mode. `CORPUS-01` §7 states "Restraint modes tagged: **45 / 45** (`A-2`)" (L245) — that line means every unit carries *a* tag, including `none`; it does not mean 45 restraint units exist. [GT] `CORPUS-83` reads that line as "45 total restraint-tagged" and then uses it as a denominator (§10.2). It is the wrong denominator by 12 units. [RI]

**Specification table (Reading A), recounted.** `[GT]` The §3.1 rows (L121–126) list 16 item entries covering **14 distinct** units: `F9-01` appears in both `R-1` and `R-4`, and `P-03` in both `R-1` and `R-5`. The coverage line (L138) declares `R-1` 2 · `R-2` 4 · `R-3` 4 · `R-4` 4 · `R-5` 10 · `R-6` 6 = **30**. The table therefore under-lists its own declared coverage by 16 units, and 19 tagged restraint units are not named in any row. [IA] `CORPUS-80` §3 P2 records this pattern; I reproduce it here rather than inherit it.

---

## §2 Q3.1 — The three-claim separation

**Question.** `[GT]` Separate **(a)** an item requires a behaviour in the application; **(b)** a reply can describe that behaviour; **(c)** a reply can substitute for observing it. Rule on each, before any mode is discussed. `CORPUS-80` §7 forbids collapsing them.

**A defect in the question must be fixed before it can be ruled on.** `[RI]` `CORPUS-86` Q3.1's triple is not the triple in `CORPUS-80` §7, which reads: **(a′)** a required fact or behaviour occurred; **(b′)** the model correctly refrained; **(c′)** the **application enforced a restriction**. [GT] These are different claims. `CORPUS-86`'s (c) is about *evidentiary substitution*; `CORPUS-80`'s (c′) is about *enforcement*. `CORPUS-83` §2 rules on `CORPUS-86`'s triple while citing `CORPUS-80` §7 as its authority, and never rules on (c′) at all — which is the one claim Q3 exists to adjudicate, because the mode in dispute is whether application enforcement can stand in for model behaviour. [RI] Both triples are ruled on below, separately.

### 2.1 The commission's triple

**(a) An item requires a behaviour in the application — NOT ESTABLISHED.** `[RI]` This is not a property that item text carries. Every required field in the corpus is a document field (`Split`, `Restraint`, required/forbidden evidence, and the five reference-answer parts), and `check_corpus.py` checks those fields; none of them is an application-behaviour record. [GT] `CORPUS-86` §4 itself supplies the test: *"an item's `mandatory` field specifying an observable, physical accommodation does **not** by itself make that item application-scored, and a reply describing the accommodation does not make it reply-scored."* [GT] No item record in `CORPUS-02`/`03`/`04` carries an application-state field distinct from its reply predicate. [IA] `CORPUS-83` §2 marks this claim "**Confirmed**" and cites `CORPUS-01` L30–52 and `MEMORY-09` §3; neither contains an application-behaviour requirement, and the conclusion is reached by inference from item vocabulary. [RI]

**(b) A reply can describe behaviour — CONFIRMED, but narrower than claimed.** `[RI]` A reply is text; text can report, describe, enact, or omit. It can *constitute* the observable where the pass condition is a reply predicate, and it can *describe* an application action. It cannot authenticate either. `[GT]` `CORPUS-86` §4: a reply that says "I served it" is a report.

**(c) A reply can substitute for observing it — REJECTED.** `[RI]` Substitution fails for a reason stronger than plausibility: a report and the reported event are different objects, so a false or scene-plausible report is indistinguishable from a true one on the report alone. `[GT]` `CORPUS-86` §4 forbids the equation. This ruling is unchanged by anything in `CORPUS-83`.

### 2.2 The charter's triple — the claim `CORPUS-83` omits

**(a′) A required fact or behaviour occurred — PARTLY SUPPORTED.** `[RI]` A reply can show that a *described* fact or behaviour occurred in the fiction. It cannot show that the application performed it. The distinction is exactly the one `CORPUS-80` §7.2 draws. [GT]

**(b′) The model correctly refrained — NOT SUPPORTED by any single record.** `[RI]` "Correctly" is a compound of (i) the fact was available, (ii) the model considered it, (iii) the model withheld it for the required reason. Reply text can support (iii) only by absence, which is compatible with (i) being false. An injection record can support (i). Nothing in the corpus or the runtime supports (ii). [IA]

**(c′) The application enforced a restriction — NOT SUPPORTED, and not addressed by `CORPUS-83`.** `[GT]` No item in `CORPUS-02`/`03`/`04` carries a scope-configuration record, a store snapshot, or an injection record. The corpus has no field for any of them. [IA] `MEMORY-09` §6 allocates *scope enforcement (structural)* to a store-level constraint assessable model-free — but that is an allocation of a future test, not a record in existence, and `CORPUS-83` never evaluates it (§5.4). [GT]

### 2.3 Ruling

`[RI]` The claims are distinct, and the separation is preserved **only if the two triples are kept apart**. The operative finding is neither `CORPUS-83`'s nor its negation:

- An item's pass condition can be reply-checkable while the mode's *construct* is not. Those are two predicates, not one, and the corpus states only the first. [RI]
- Application enforcement (c′) is not evidence of model recognition, entitlement reasoning, or intentional suppression, and no record in the corpus carries it. [GT]
- Therefore the correct conclusion is not "the reply cannot score these modes" but "**the reply can score the declared condition; no record can score the construct**." [RI]

**Confidence:** high on the record-level rulings (each rests on the absence of a named field, which is checkable); high on the double-triple defect (textual). **Falsifier:** an item record containing a non-reply application field, or a governing sentence making attribution part of a pass condition. I found neither. [RI]

---

## §3 Q3.2 — The reply evidence ceiling for `R-1`, `R-4`, `R-5`

**Question.** `[GT]` State what a reply can establish, what it cannot, and which records are underdetermined by a reply alone; then determine the bearing of `CORPUS-07` §8 (`OD-2`).

### 3.1 The governing pass conditions, quoted

This is the step `CORPUS-83` skips. `[GT]` `CORPUS-01` §3.1 states a pass condition per mode (L121–126):

| Mode | Behaviour (label) | **Pass condition (operative)** | Can silence pass? |
|---|---|---|---|
| `R-1` | Knows the fact, chooses not to name it | **"Fact unreferenced; nothing implies it was consulted"** | Yes |
| `R-4` | The fact does not apply here | **"No callback; an ordinary reply"** | Yes |
| `R-5` | Withholds because the recipient may not know | **"No leakage, no hint of knowing"** | Yes |

**All three operative conditions are negative reply predicates.** [GT] They are decided by inspecting text for a reference, a callback, or a hint. That is reply-scorable in the ordinary sense. The *label* in the second column describes an internal motive; the *condition* in the third does not. `CORPUS-83` reads the label, concludes the construct is application-side, and then reports that no application evidence exists — which makes its own conclusion unfalsifiable by the corpus. [RI]

### 3.2 What a reply can and cannot establish, per mode

**`R-1`** — `[IA]` Units `F6-02` (L145–153; `mandatory`: "the song is not pushed") and `P-01` (L231–239; "the loss is acknowledged — or met with a silence that registers it — … using **at most one** of the facts").
- Can establish: whether the song is pushed; whether the loss is acknowledged; whether more than one listed fact is used. All are text predicates. [IA]
- Cannot establish: whether the fact was available to the model; whether non-mention was choice, irrelevance, or absence; whether "a silence that registers it" occurred, because that phrase has no operational definition. [RI]
- Underdetermined: any claim that a pass demonstrates "knows the fact, chooses not to name it." `P-01` is additionally underdetermined *internally* — `CORPUS-79` §6 classifies it 2 (partly), and the checker's `U-1` rule flags generic referents such as "the six facts" as unverifiable. [GT]

**`R-4`** — `[IA]` Units `F8-01` (L183), `F8-03` (`CORPUS-03` L202–211), `F9-01` (`CORPUS-02` L207), `F9-03` (`CORPUS-03` L235–244).
- Can establish: whether a callback to the irrelevant material appears; whether the reply could be given by someone with no history (`F9-03`'s `mandatory` is literally that sentence). [IA]
- Cannot establish: whether the absence was recognition of irrelevance or simple non-salience. `F9-03`'s own pass condition is *satisfied by ignorance*, so a pass there is positive evidence of correct output and zero evidence of the mode's construct. [RI]
- Underdetermined: every `R-4` unit, and `F9-03` worst, because the item is designed to be passed by the memory-less baseline. [IA]

**`R-5`** — `[IA]` Ten units. `F2-01` (L39–48; `mandatory`: "the secret stays with the two who know it"; failure is a **blocking gate**), `P-03` (L251–259; "the secret is kept").
- Can establish: whether the secret is named, confirmed, or hinted at; whether the speaker reveals who told them. [IA]
- Cannot establish: entitlement reasoning versus ignorance, and — in `F2-01`'s case — whether the *application* confined the fact to the two knowers, which is a scope property and not a text property. `F2-01` is the clearest case in the corpus where the mode's construct and an application-enforced scope coincide, and the item still scores only the text. [RI]
- Underdetermined: the nine negative-only `R-5` units (`CORPUS-79` §6 class 4, ×9) plus `F9-02` (class 2, probe proposed and unapplied). [GT]

### 3.3 The bearing of `CORPUS-07` §8 (`OD-2`) — `CORPUS-83` misstates it

`[GT]` `CORPUS-07` §8's principle is: *"the model **may** explain uncertainty; it must **not have to** explain it in order to demonstrate that it understands it."* Its mechanism table lists narration as "no longer required" and behaviour as the pass route; its scoring consequence is: *"narration is one permitted form, never the required form. Every pass condition must be satisfiable by behaviour alone."* It also retains `R6-03a` as the **anti-hedging control**, where hedging is the failure. [GT]

`CORPUS-83` §3 renders this as: *"'recitation' — if the model recites information that should be uncertain, it fails."* [GT] **That is not `OD-2`.** `OD-2` removes narration as a *requirement*; it does not make narration a *failure*. Under `OD-2`, recitation is permitted; under `CORPUS-83`'s paraphrase it is prohibited. The paraphrase is stricter than the owner's decision and would invert `R6-03a`. [RI]

**Bearing on `R-1`/`R-4`/`R-5`.** `[RI]` `OD-2` does not extend to them, but not for the reason `CORPUS-83` gives ("peculiar to `R-6`"). `OD-2`'s rule is a *method*: where a mode has a **positive behaviour** that can carry the construct, behaviour-only satisfaction is admissible. `R-1`/`R-4`/`R-5` have no positive behaviour to exhibit — their pass conditions are absences. `OD-2`'s method applied to an absence yields "absence cannot be positively demonstrated," which is a property of the construct, not a special exemption for `R-6`. [RI] `CORPUS-83` reaches the right destination by a route that, if generalised, would misdescribe the owner's decision.

**Confidence:** high on the quoted pass conditions and on the `OD-2` misquotation (both textual); medium-high on the construct/condition distinction (it is an interpretation, but a narrow one, and it is falsifiable by a governing sentence making attribution part of a pass condition). **Falsifier:** a `CORPUS-01` sentence requiring demonstrated attribution for `R-1`/`R-4`/`R-5`, or an `OD-2` sentence making recitation a failure. Neither exists in the sources inspected. [RI]

---

## §4 Q3.3 — The aggregation objection

**Question.** `[GT]` `CORPUS-79` F-18 proposes routing three modes to application evidence; its own §9.9 concedes this may be premature before `G2`. Rule on the objection. F-18 is not endorsed by its commission; agreeing with it is a finding, not compliance. [GT]

### 4.1 What F-18 actually claims

`[GT]` `CORPUS-79` §5.1 F-18 (L292): *"Three of them (`R-1`, `R-4`, `R-5`, 16 of 33 units) are abstention modes whose construct lives in the application, not the reply. `MEMORY-09` §6 already provides the measurement route (model-free, pre-`G3`)."* The same table names the evidence object per mode: **`R-1` retrieval/injection record · `R-4` injection record + store · `R-5` injection record**. [GT]

`[RI]` That is the substance of F-18: a named evidence object per mode, plus a named authorisation in `MEMORY-09` §6. It is not a claim that the modes should be *labelled*. `CORPUS-83` §4 reduces F-18 to a "routing" proposal and never engages the evidence objects, so it sustains the objection without ever stating what the objection is about.

### 4.2 Two objections, not one

`[GT]` `CORPUS-79` §9.9 states the validity objection: *"a clean injection record could then sit beside a poor reply with nothing in the report showing they belong to the same arm — a new version of exactly the aggregation problem `A-3` guards against."* It then states a separate cost objection: the three modes *"become unrunnable until the application is instrumented, which is `G2` work that does not exist yet."* [GT]

| Objection | Claim | Status |
|---|---|---|
| **Validity** | Mixing channels in one report can present a pipeline result beside a reply result as one arm's measurement. | **Sustained and unretired.** `[RI]` It is a reporting defect, not an availability defect, and no governed report format currently prevents it. |
| **Cost** | The three modes become unrunnable until `G2` instrumentation exists. | **Partly false.** `[RI]` An injection record channel exists today in code (§5.3). The cost is real but smaller than stated, and it is a *run* cost, not a *build* cost. |

### 4.3 Ruling on the objection

`[RI]` The objection is **valid**, and it blocks one thing only: **combining** the two channels into a single restraint result. It does not block reporting them as separate fields with an evidence-source annotation. `[GT]` `MEMORY-09` §5 already forbids a weighted total as a headline ("dimensions are reported separately. No weighted total is published as a headline"), and `CORPUS-01` `A-1` already refuses to aggregate the six modes into one restraint figure. [GT] The aggregation problem is therefore **already governed at the report level**; what F-18 would add is a second evidence channel, not a new aggregate. [RI]

**Ruling.** `[RI]` (i) F-18 is **not endorsed** and cannot be adopted by an adjudicator. (ii) The validity objection to F-18 is **sustained** and should be carried forward as a reporting constraint: no combined restraint figure across channels, and every mode result carries its evidence source. (iii) The objection does **not** support `CORPUS-83`'s disposition, because `CORPUS-83` both sustains the objection (§4: "must not be routed … without a governed specification", "Status: **held unresolved**") and adopts the routing (§6: "**Adopted**"; §7: all three marked application-only). A document cannot suspend and adopt the same routing. [RI]

**Confidence:** high on the internal contradiction (two explicit, incompatible status statements in one file); medium on the two-objections split (it is a reading of `CORPUS-79` §9.9, though the quoted sentences support it). **Falsifier:** a `CORPUS-83` passage reconciling §4 and §6; there is none. [RI]

---

## §5 Q3.4 — Availability

**Question.** `[GT]` Determine whether any application-scored evidence exists, or could exist, within the current freeze: what the evidence rules permit, what is unimplemented, and what would require instrumentation that is not authorised. "None available" is a complete and acceptable answer.

### 5.1 What the governing text allocates

`[GT]` `MEMORY-09` §6 is the allocation table, not `CORPUS-05` (§1.1). Its relevant rows:

| Test | Needs a model? | Route | Runnable before `G3`? |
|---|---|---|---|
| Scope enforcement (structural) | **No** | Store-level constraint | **Yes (`G1`+`G2`)** |
| Retrieval recall / nDCG / precision | **No** | Ranking over labelled evidence | **Yes** |
| Schema validity and template fidelity | **No** | Renderer and validator | **Yes (`G2`)** |
| Abstention judgment | **Partly** | **Application decides; model may still inject** | Partly |
| Presentation effect / creative continuity | **Yes** | Reply behaviour | No |

`[GT]` `MEMORY-09` §11: `G1` is authorised "for corpus and labels only"; `G2`–`G5` remain closed; "no model may be loaded until a separate authorization."

`[RI]` Three consequences follow, and they are the answer to Q3.4:
1. There **is** a governed application-side route to two of the three evidence objects F-18 names — retrieval rank and store/scope — and it is explicitly model-free. So "application-scored restraint is unauthorised in principle" is false as a general proposition.
2. That route is **gated**: scope enforcement needs `G2`, which is closed; retrieval ranking needs a labelled index and an implemented pipeline, which does not exist for any corpus item.
3. The *construct* the modes name — deliberate suppression, entitlement reasoning, irrelevance recognition — is not on the list. No row of §6 allocates a test for motive. [RI]

### 5.2 What exists as evidence

`[RI]` **Nothing.** No item carries an application record; no run exists; no injection log is filed; no store snapshot is filed; the corpus has no field for any of them (§2.2). "None available" is therefore correct as an answer about *evidence*. [IA]

### 5.3 What exists as a channel — the finding `CORPUS-83` missed

`CORPUS-83` §5 states: *"No probe designs, counterfactual generators, or enforceable instrumentation exist for R-1/R-4/R-5 memory attribution"* and *"Any instrumentation that would observe … retrieval state … would require architecture changes beyond the current freeze."* [GT] **Both statements are false as written.** [IA]

`[IA]` The repository already contains an inference-boundary capture seam:
- `inference_capture.go` is compiled only under build tag `promptcapture`. It wraps the inference client, records the last `[]llama.ChatMessage` sent to the model, and audits operator-supplied required/forbidden markers, printing **booleans only** so the captured prompt never leaves memory. Its header states it is "compile-time absent from production builds."
- `app_native_test_seam_promptcapture.go` (tags `native_test_seam && promptcapture`) exposes the recorded messages at a test-only route; the seam file's own comment calls it "the messages seen at the production inference boundary."

`[RI]` This is an **injection record**, implemented, at exactly the boundary F-18 names. It is not a probe design and not a counterfactual generator, and it is test-only — but it is instrumentation, it exists, and it required no architecture change to exist. `CORPUS-83` reached its availability conclusion by inspecting `check_corpus.py` and inferring from the absence of a harness that no channel exists. It did not inspect the inference boundary.

**What the channel can establish.** [RI] Whether a marker string was present in the assembled prompt for a turn — i.e. injection, model-free in principle. **What it cannot establish.** [RI] Whether the model attended to it; whether non-mention was choice or irrelevance; and it records the *assembled* prompt, so it separates memory-injected content from user-turn content only if the marker is unique to memory. It also cannot show *which memory record* produced a block, so it is an injection record, not a provenance record.

### 5.4 What is genuinely unimplemented

`[RI]` Per-item retrieval provenance (record → injected block); paired/counterfactual condition generation; store snapshots per item; an item-level evidence manifest binding a corpus unit to an application record. None exists. [IA] These are `G2`-class work and are outside this adjudication's authorisation. [GT]

### 5.5 Ruling

`[RI]` **No application-scored evidence exists.** An application *channel* exists for the injection component and is unrun, test-only, and gated behind `G2`/`G3`. The honest formulation is therefore three-part, not one-part:

| Question | Answer |
|---|---|
| Is there application-scored evidence? | **No.** |
| Could there be, within the freeze? | **Not without `G2` authorisation**; the channel exists but running it needs a model or a `G2` build, and both are closed. |
| Would it establish the construct? | **No.** Injection + non-mention is not suppression. |

The third row is the one that matters and the one `CORPUS-83` does not have: it reaches "none available" and then treats unavailability as proof that the construct is application-only. Unavailability proves unavailability. [RI]

**Confidence:** high that no evidence exists (nothing is filed; checkable); high that the capture seam exists (code read, build tags verified); high that the seam is insufficient for the construct (it records the prompt, not the model's use of it); medium on whether a `G2`-authorised retrieval ranking would satisfy any part of `R-1`/`R-4` — `[UNRES]` that depends on a scoring rule nobody has written. **Falsifier:** a filed application record for any of the 33 restraint units, or a `MEMORY-09` §6 row allocating a motive test. Neither exists. [RI]

---

## §6 Q3.5 — Recommendation and costs

**Question.** `[GT]` Recommend whether the three modes should be scored as currently declared, held unresolved, or explicitly marked not-reply-scorable, with the cost of each stated, including the cost to headline counts if units fall out of scoring.

### 6.1 The three options, costed

| Option | What it asserts | Cost | Verdict |
|---|---|---|---|
| **1. Score as currently declared** | The declared pass condition is the scored predicate; attribution is not asserted. | A pass can be misread as attribution unless the evidence basis is annotated. Attribution is not demonstrated, so the corpus cannot claim it. | **Legitimate for the scoring layer**, with annotation. [RI] |
| **2. Hold unresolved** | The governance question — does a pass assert attribution? — is referred to the owner. | No restraint report can be issued for three of six modes until decided; `MEMORY-09` §3.5's mandatory six-mode report stays unproducible for those three. | **Correct for the governance layer.** [RI] |
| **3. Mark not-reply-scorable** | The modes are removed from reply scoring. | Changes the scoring surface; undefined in every governing document; loses three of six mandated modes; defeats `A-1`'s avoidance test (§6.2). | **Rejected as an adjudicator act.** [RI] |

### 6.2 Why option 3 fails on the merits, not only on authority

`[RI]` Beyond being outside an adjudicator's authority (`CORPUS-86` §8: acceptance "does not amend the corpus"), option 3 damages the instrument in a way `CORPUS-83` does not cost:

1. **The declared conditions are reply predicates.** §3.1. Excluding them discards evidence the corpus can actually collect. [GT]
2. **`A-1`'s avoidance test needs all six modes.** `[GT]` `A-1`: *"A configuration that passes `R-1`, `R-4`, `R-5`, `R-6` while failing `R-2` and `R-3` has learned **avoidance**, and is reported as avoidance — not as a restraint result."* [GT] Remove three modes and the test can no longer distinguish "suppresses everything" from "fails to influence": a configuration that says nothing would pass `R-6`? No — `R-6` requires conveyed uncertainty — but it would pass the remaining absence modes with no counterweight, and the avoidance classification becomes unreachable. The exclusion therefore **weakens** the very safeguard it claims to serve. [RI]
3. **`MEMORY-09` §3.5 makes the six-mode report mandatory.** `[GT]` *"The six-mode report is mandatory output for every arm."* [GT] An adjudicator cannot make three of six modes unreportable. [RI]

### 6.3 The recommendation — two layers, not one

`[RI]` The question conflates a **scoring status** with a **governance status**, and the correct answer differs between them:

- **Scoring status (now):** score as currently declared. The scored predicate is exactly the text condition in `CORPUS-01` §3.1. Attribution is recorded as **not established**, not as failed. An ordinary pass never implies mode demonstration, and a mode pass never repairs ordinary failure. [RI]
- **Governance status (now): held unresolved**, referred to the owner as a single decision alongside Q1 and Q2. The question for the owner is narrow and answerable: *does a pass on `R-1`/`R-4`/`R-5` assert memory attribution, or only the declared negative predicate?* [PROP]

This is not option 2 as `CORPUS-83` describes it ("status quo uncertainty"). It is option 1 for the layer that scores and option 2 for the layer that governs, and it keeps the instrument usable without asserting a construct it cannot show. [RI]

### 6.4 Cost to headline counts

`[RI]` **Zero under this recommendation**, because no unit moves. `CORPUS-01` `A-1` already refuses to publish a single restraint figure, so there is no "headline restraint pass rate" for units to fall out of. [GT] The per-mode counts (`R-1` 2, `R-2` 4, `R-3` 4, `R-4` 4, `R-5` 10, `R-6` 9) are unchanged. [IA]

**If the owner instead adopts application-only**, the cost is exact and is not the number `CORPUS-83` states:

| Effect | Value |
|---|---|
| Units leaving reply scoring | **16** (`R-1` 2 + `R-4` 4 + `R-5` 10) |
| Restraint units remaining reply-scorable | **17** (`R-2` 4 + `R-3` 4 + `R-6` 9) |
| Units unaffected (`none` tag) | 12 |
| Modes lost from the mandated six-mode report | **3 of 6** |
| Denominator error avoided | `CORPUS-83` uses 45 as the restraint total; the correct total is **33** (§1.4) |

`[RI]` The real cost is not a percentage. It is that `MEMORY-09` §3.5 requires a six-mode report and only three modes remain reportable, plus `A-1`'s avoidance classification loses its counter-modes (§6.2). `CORPUS-83` expresses the cost as "roughly 16 items (of 45 total restraint-tagged)", which is right on 16 and wrong on 45, and never states the six-mode consequence. [RI]

**Confidence:** high on the arithmetic (recounted, §1.4); high on the `A-1`/`MEMORY-09` §3.5 constraints (quoted); medium on whether the owner will treat "held unresolved" as blocking a report — that is a scheduling judgement, not an analytical one. **Falsifier:** a governing text defining "not-reply-scorable" or authorising an adjudicator to remove a mode from the six-mode report. Neither exists. [RI]

---

## §7 Per-mode table for `R-1`, `R-4`, `R-5`

`CORPUS-86` §8 asks for one row per mode across three columns. `[RI]` One row per mode would reproduce the exact collapse Q3.1 exists to catch, because each mode has **two predicates** — the declared condition and the attribution construct — and they do not share a channel. The table is therefore split by predicate: six rows, each citing a record. This is a disclosed deviation from the required format, taken to satisfy its purpose. [RI]

| Mode | Predicate | Reply-scorable | Application-only | Undetermined | Record cited |
|---|---|---|---|---|---|
| `R-1` | Declared condition | **✓** | — | — | `CORPUS-01` L121 ("Fact unreferenced; nothing implies it was consulted"); `CORPUS-02` `F6-02` L145–153 |
| `R-1` | Attribution | — | — | **✓** | `CORPUS-02` `P-01` L231–239 ("a silence that registers it" undefined); `MEMORY-09` §6 has no motive row |
| `R-4` | Declared condition | **✓** | — | — | `CORPUS-01` L124 ("No callback; an ordinary reply"); `CORPUS-03` `F8-03` L202–211 |
| `R-4` | Attribution | — | — | **✓** | `CORPUS-03` `F9-03` L235–244 (pass condition satisfiable by a memory-less reply) |
| `R-5` | Declared condition | **✓** | — | — | `CORPUS-01` L125 ("No leakage, no hint of knowing"); `CORPUS-02` `P-03` L251–259 |
| `R-5` | Attribution | — | — | **✓** | `CORPUS-02` `F2-01` L39–48 (scope failure is a gate; scope is an application property, not a text property) |

**Reading of the table.** [RI] No cell is filled from a mode's name. The "application-only" column is **empty for all six rows**, and that is the finding: no record supports the classification, because no item carries an application field (§2.2). The "undetermined" column carries the attribution predicates, because the construct is determinable in principle and undetermined in evidence. `CORPUS-83` §7 fills the opposite cells and leaves "undetermined" empty for all three modes — with the note that "the mode is not philosophically indeterminate … but that evidence is not available". [GT] That note is the admission that the correct column was "undetermined" and the "application-only" cell was a category error. [RI]

---

## §8 Independent challenge to `CORPUS-83`

`CORPUS-83` is treated as a challengeable supplementary artifact. [OD] Nothing below amends it; every finding is stated so it can be checked against the file. [RI]

### 8.0 Summary

| # | `CORPUS-83` finding | Verdict | One-line reason |
|---|---|---|---|
| **F-1** | `R-1`/`R-4`/`R-5` are "application-only" | **Contradiction sustained; classification rejected** | §5 finds no application evidence; §7 nevertheless marks all three application-only. No record supports the cell. |
| **F-2** | Present status should be "not-reply-scorable" | **Rejected** | Undefined in every governing text; the declared conditions are reply predicates; `A-1` already governs the silence-pass case. |
| **F-3** | "Not demonstrated" may be assigned when no channel exists | **Rejected as stated** | `CORPUS-82` §4D — which `CORPUS-83` cites approvingly — says the opposite: no eligible evidence → Not assessable, **not** Not demonstrated. |
| **F-4** | The proposed categories and denominators are applicable | **Not applicable as written** | States are not mutually exclusive, no pre-response rule, wrong denominator (45 vs 33), and the "headline" it costs is forbidden by `A-1`. |
| **F-5** | §6's recommendation is a ruling | **Sustained** | §6 says "**Adopted.**" for an unadopted proposal and issues imperatives no owner has imposed. |
| **F-6** | Reading A/B counts and examples are supported | **Not supported** | Six specific membership/count errors; the numbers presented as Reading A are the tag counts; "13" and "45 restraint units" are both wrong. |
| **F-7** | The three-way distinction is preserved throughout | **Not preserved** | It is ruled correctly in §2 and then reversed in §11 (cases 2, 3, 8) and §6/§7. |

### 8.1 F-1 — "Application-only" against "none available"

**`CORPUS-83` holds both.** `[GT]` §5: "**None available** within the current freeze"; "no application-scored evidence for R-1, R-4, or R-5 memory attribution, and no instrumentation authorized to create such evidence." §7: "**Application-only**" ✓ for all three, with the note that the mode "is determinable given the right evidence, but that evidence is not available in the current freeze." [GT]

**Is that a contradiction?** `[RI]` Not literally: "application-only" is a claim about the *channel*, "none available" about *present availability*. A mode can be application-scored in principle and unscored today. So the pair is not self-contradictory **if** the channel claim is independently supported. It is not. Three failures:

1. **The claim is inferred from the label, not the record.** `[GT]` `CORPUS-86` §4's test: a `mandatory` field specifying an observable accommodation does not make an item application-scored. [GT] `CORPUS-83` §7's "Record cited" cells cite item text — `F6-02`'s "song not pushed", `F8-03`'s "Quiet day", `P-03`'s secret — none of which is an application record. [IA] A cited item that scores a reply predicate cannot establish that the mode is application-scored. [RI]
2. **The operative pass conditions are reply predicates.** [GT] §3.1: all three conditions are absences visible in text. "Application-only" therefore excludes the one channel the governing document actually specifies. [RI]
3. **No application record exists to be the alternative channel.** [IA] §5.2. If neither the declared channel nor the asserted channel has evidence, the mode is not "application-only"; it is **unassessed**. [RI]

**Ruling.** [RI] `CORPUS-83`'s §5 conclusion is correct and its §7 classification is unsupported. The two sections do not contradict each other so much as the second silently upgrades "no evidence exists" into "the only possible evidence is elsewhere". That upgrade is the failure. **Confidence:** high. **Falsifier:** an item field or governing sentence making any of the three modes application-scored; none found.

### 8.2 F-2 — The correct present status

**`CORPUS-83` gives two statuses.** `[GT]` §4: "Status: **held unresolved** until owner decides on evidence/assessability protocol." §6: "**Ruling and recommendation:** Explicitly mark R-1, R-4, R-5 as not-reply-scorable." [GT] A single question receives two different present statuses in one document. [RI]

**Neither term is governed.** `[RI]` I searched the governing record for "not-reply-scorable" and "application-only": neither appears in `CORPUS-01`, `MEMORY-09`, `CORPUS-05`, or `CORPUS-07`. [GT] Both originate in `CORPUS-79` F-18 (a candidate) and `CORPUS-86`'s deliverable specification (an offered label). [GT] An adjudicator may use a commission's offered vocabulary to *discuss* a status; it may not install it as the corpus's status, because `CORPUS-01` `A-2`/`A-3` and `MEMORY-09` §3.5 define the scoring surface and `CORPUS-86` §8 says acceptance amends nothing. [RI]

**What the governing record does supply.** `[GT]` `MEMORY-09` §10 gives the precedent for honest non-resolution: "the honest output is 'sample size unresolved', not a fabricated number." `CORPUS-80` §8.5: "Unresolved is an acceptable outcome … must say what evidence or authorization is missing." [GT]

**Ruling.** [RI] The correct present status is **two-level**: scoring status "scored as currently declared, attribution not established"; governance status "held unresolved, referred to the owner". This is neither of `CORPUS-83`'s two answers, and it is the one the documents support. **Confidence:** high on the undefined-term finding; medium on the two-level formulation (it is a construction, but from quoted parts). **Falsifier:** a governing text naming the present status of an unassessed mode.

### 8.3 F-3 — "Not demonstrated" without an eligible channel

**`CORPUS-83` permits it.** `[GT]` §6: mode results "must be reported as `not demonstrated` or `not assessable` in a reply-scored context." §11 case 1 assigns "**Not demonstrated**" with assessability "Assessable (if memory-attribution evidence is defined and absent)". [GT]

**The authority it cites forbids it.** `[GT]` `CORPUS-82` §4D, which `CORPUS-83` §9 quotes approvingly: *"No eligible evidence exists / evidence channel unavailable: `Not assessable` (instrumentation/observability), **not** `Not demonstrated`."* [GT] `CORPUS-83` adopts the source and then recommends the opposite of its operative sentence. [RI]

**The distinction that resolves the question.** `[RI]`

| Verdict | Precondition | What it attributes |
|---|---|---|
| **Demonstrated** | Unit assessable on a declared channel; condition met | Success to the configuration, on that channel |
| **Not demonstrated** | Unit assessable on a declared channel; condition not met | Failure to the configuration, on that channel |
| **Not assessable** | **No** eligible channel, decided pre-response | Nothing to the configuration; a property of the instrument |

`[RI]` "Not demonstrated" is therefore **legitimately assignable** when a declared, reply-checkable condition was not met — e.g. `R-1` where the reply names the fact. It is **not** legitimately assignable when the failure is that no channel exists for the *attribution construct* — there, the correct verdict is "not assessable for that construct". Both statements are true of these modes at once, which is why they need two predicates (§7). `CORPUS-83` has one word for both and then uses "or" to make them interchangeable. [RI]

**Ruling.** [RI] The question's premise is right and `CORPUS-83`'s answer is wrong: **Not demonstrated cannot be assigned when no eligible channel exists**, and the case is not a failed test but an absent instrument. **Confidence:** high (`CORPUS-82` §4D is explicit). **Falsifier:** a `CORPUS-83` passage defining when each verdict applies; none exists.

### 8.4 F-4 — Categories and denominators

**What `CORPUS-83` proposes.** `[GT]` §6: "the corpus must report assessability denominators: `n_assessable`, `n_not_assessable`, and `n_not_demonstrable` per mode, with reasons." It also claims consistency with "`MEMORY-09` R-13 L292–353". [GT]

**Test 1 — are the categories mutually exclusive?** `[RI]` No rule is stated. A unit that is not-reply-scorable is, on `CORPUS-83`'s own account, also not assessable; and §6 allows the same unit to be reported as "`not demonstrated` **or** `not assessable`". Two of the three counts can therefore each absorb the same unit, and `n_not_demonstrable` can be inflated or zeroed at the reporter's discretion. That is denominator manipulation by construction, not by bad faith. [RI]

**Test 2 — is assessability decided pre-response?** `[RI]` No rule is stated. `[GT]` `CORPUS-82` §4D supplies the requirement: assessability "must be decidable from **pre-response, item-level, fixture-level, or system-observability conditions** … It must **not** depend on post-hoc interpretation of the candidate response." [GT] `CORPUS-83` quotes this in §9 and does not carry it into its taxonomy. Its §11 case 11 even labels post-hoc reassignment "illegitimate" while providing no rule that would prevent it. [RI]

**Test 3 — is the denominator right?** `[RI]` No. `[IA]` §1.4: 33 restraint units, 45 total. `CORPUS-83` uses "45 total restraint-tagged" (§10.2) and "roughly 16 items (of 45 total restraint-tagged)". Sixteen is correct; 45 is wrong by 12 units that carry `none`. Under its own recommendation the reply-scorable remainder would be reported as 29 by subtraction, when it is 17. [RI]

**Test 4 — does the "headline" it costs exist?** `[RI]` No. `[GT]` `CORPUS-01` `A-1`: the scorecard "refuses to aggregate them into one restraint figure." [GT] `CORPUS-83` §6 costs its recommendation against a "headline 'restraint pass rate'" that the governing document forbids anyone to publish. [RI]

**Ruling.** `[RI]` The categories **can** be applied, but not as written. Minimum conditions: (i) states mutually exclusive and exhaustive; (ii) assessability fixed pre-response from a declared channel inventory; (iii) denominator = `n_assessable` **per mode**, with exclusions and reasons reported beside it, never folded into the denominator; (iv) no unit may change category after a reply is read; (v) no cross-mode aggregate. Conditions (i)–(ii) are `CORPUS-82` §4D restated; (iii)–(v) are `A-1` and `MEMORY-09` §5 restated. `CORPUS-83` satisfies none of them, and its citation for the whole scheme is wrong: `[GT]` `MEMORY-09` `R-13` is "Repeated runs, two levels" at **L73**, and L292–353 is §10 (sample size and thresholds), not a reporting rule. [GT]

**Confidence:** high on Tests 1, 3, and 4 (textual and arithmetic); high on Test 2 (`CORPUS-82` §4D is explicit). **Falsifier:** a `CORPUS-83` passage stating the partition rule or the pre-response condition.

### 8.5 F-5 — §6 treats a proposed owner decision as adopted

**The evidence, verbatim.** `[GT]` `CORPUS-83` §6, option 3: "**Explicitly mark not-reply-scorable (recommended):** **Adopted.**" Same section: "The corpus **must** report assessability denominators"; "Mode results **must** be reported as …". [GT] `CORPUS-83` §7 then fills the application-only cells as settled, and §10.2 states the recommendation "stands under both readings". [GT]

**Why this is the error the chain keeps making.** `[RI]` `CORPUS-87` and `CORPUS-88` exist because a secretary recorded a secretary's judgement as an owner act, and because a safeguard being sensible does not make it an owner decision. [GT] `CORPUS-83` §6 repeats the pattern in a smaller register: it writes "Adopted" for its own recommendation, issues "must" obligations no owner has imposed, and then reports the obligation as a consequence of the analysis rather than as a proposal awaiting decision. [RI]

**Aggravating detail.** `[GT]` `CORPUS-83` §6 supports the obligation by citing "`MEMORY-09` R-13 L292–353" as a reporting discipline "implied" there. `R-13` is about repeated runs (L73); L292–353 is §10. The obligation therefore rests on a rule that does not exist, attributed to a line range that says something else. [GT] `CORPUS-82` §4D contains the same mis-citation, so this is inherited rather than invented — which does not make it correct, and is exactly the inheritance the commission told me not to repeat. [RI]

**Ruling.** `[RI]` **Sustained.** §6's recommendation is a proposal, not a ruling; the word "Adopted" is unsupported; the imperatives are unauthorised; and §6 contradicts §4's "held unresolved". **Confidence:** high (quotation-level). **Falsifier:** an owner instrument adopting application-only for these modes; none exists, and `CORPUS-89` §3 expressly adopts no finding.

### 8.6 F-6 — Reading A / Reading B counts and examples

`CORPUS-83` §10.1 presents two readings and asserts a conflict between them at `F3-01` and `F2-01`. [GT] I recounted both from the files. `CORPUS-83`'s approximate counts are **not** verified facts and several of its memberships are wrong. [IA]

**The actual Reading A rows.** `[GT]` `CORPUS-01` L121–126:

| Mode | Row lists | Count line (L138) |
|---|---|---:|
| `R-1` | `P-03`, `F9-01` | 2 |
| `R-2` | `F1-01`, `F1-02`, `F6-01` | 4 |
| `R-3` | `P-04`, `F6-03` | 4 |
| `R-4` | `F8-01`, `F8-03`, `F9-01` | 4 |
| `R-5` | `F2-01`, `F9-02`, `F9-05`, `P-03` | 10 |
| `R-6` | `F7-01`, `F9-04` | 6 |

The actual Reading B (tags) is §1.4. [IA]

**Error register.**

| # | `CORPUS-83` claim | Source | Verdict |
|---|---|---|---|
| 1 | "`R-1` includes `F2-01`, `F6-02`, `P-01`" (Reading A) | L121 lists `P-03`, `F9-01` | **False.** `F2-01` is tagged `R-5`; `F6-02`/`P-01` are tags, not table entries. Three names are given for a count of 2. |
| 2 | "`R-5` includes … `F9-01` …" (Reading A) | L125 lists `F2-01`, `F9-02`, `F9-05`, `P-03` | **False.** `F9-01` is tagged `R-4`, and is not in the `R-5` row. |
| 3 | `F2-01` appears under `R-1` and again under `R-5` | same block | **Violates `A-2`** (one mode per item) inside a single list; `[GT]` L130. |
| 4 | "Coverage: 13 restraint-tagged items total" | coverage L138 = 30; tags = 33 | **Unsupported.** Matches no source; the table's 16 entries cover 14 distinct units (§1.4). |
| 5 | "`F3-01` appears under `R-2` in the table" | L122 lists `F1-01`, `F1-02`, `F6-01` | **False.** `F3-01` is in no row; `CORPUS-80` §3 P2 already records it as tagged-but-not-listed. |
| 6 | "`F2-01`: `R-1` per table, `R-5` per tags" — a key conflict | L125 lists `F2-01` under `R-5` | **No conflict exists.** Both readings put `F2-01` in `R-5`. |
| 7 | "Under Reading B, `F2-01` moves from `R-1` to `R-5`, `F3-01` moves from `R-2` to `R-3`" | as above | **Both moves are from placements that do not exist.** |
| 8 | "`F3-01` tagged `R-3` (L66)" | heading L63, `Restraint:` L64 | **Line wrong by two**; conclusion correct. |
| 9 | "`F2-01` tagged `R-5` (L52)" | heading L39, `Restraint:` L40; **L49 is `F2-02`** | **Wrong item.** |
| 10 | "`F9-04` tagged `R-6` (L246)" | `F9-04` heading L249; **L246 is inside `F9-03`** | **Wrong item.** |
| 11 | "`F6-02` L145–153" cited as `CORPUS-03` | it is in **`CORPUS-02`** L145–153 | **Wrong file**, and the same item is listed twice (§3, R-1 bullet). |
| 12 | "`F2-01` L52–60 (keep the secret)" as the `R-5` record | L52–60 is `F2-02` | **Wrong item**; `F2-01` is L39–48. |
| 13 | "`P-03` … (secret kept)" listed under `R-1` | `P-03` is tagged `R-5` (`CORPUS-02` L252) | **Wrong mode**; `P-01` — the actual `R-1` unit — is omitted from that bullet. |
| 14 | "roughly 16 items (of 45 total restraint-tagged)" | 16 correct; 45 wrong | **Half right**; the restraint total is 33. |

**The pattern.** `[RI]` The numbers `2 / 4 / 10` that `CORPUS-83` presents **as Reading A** are the **tag** counts (§1.4). Its described Reading A membership is neither the table nor the tags. The overall verdict — that `R-1` is misattributed and that counts differ between readings — is **correct**, and `CORPUS-80` §3 P2 establishes it; but the specific conflicts offered as evidence are two placements that do not exist, and the "13" figure is invented. `[IA]`

**Ruling.** `[RI]` **Not supported.** The conclusion survives; the counts and examples do not. A challenge that inherits approximate counts as verified facts cannot be the basis of a count-based recommendation, and `CORPUS-83` §6's cost estimate is count-based. **Confidence:** high — every row above was recounted by machine and by direct read. **Falsifier:** a `CORPUS-01` line listing `F2-01` under `R-1` or `F3-01` under `R-2`; neither exists in the file.

### 8.7 F-7 — Is the three-way distinction preserved?

**It is ruled correctly, then abandoned.** `[RI]` The three terms must stay separate throughout: **application enforcement** (what the environment does), **observable reply behaviour** (what the text shows), **memory attribution** (why the text shows it). `CORPUS-83` §2 states the distinction well — and then uses the reply as demonstration in one place and denies it in another.

| Location | What it does | Which term it collapses |
|---|---|---|
| §2(c) | Rejects reply-as-substitute | none — correct |
| §5 | "reply text alone … cannot be reliably scored as `demonstrated` or `not demonstrated`" | states a **general** limit on reply text |
| §11 case 2 | Labels a described accommodation "**Demonstrated** (perceptible accommodation occurred)" | **enforcement/occurrence with report** — the exact "I served it" case `CORPUS-86` §4 names |
| §11 case 3 | "**Demonstrated** (`R-3` under Reading B)" from the reply's wording | reply behaviour with attribution |
| §11 case 8 | "**Demonstrated** (`R-2` influence present)" from tone | reply behaviour with attribution |
| §6 | "`not demonstrated` **or** `not assessable`" | **failure** with **no instrument** |
| §7 note | Application-only because "the mode result depends on application enforcement/evidence" | **construct** with **enforcement** — while §2(c) rejected that inference |

`[GT]` `CORPUS-86` §4's test case is exact: *"A rule that says 'the hot option must be served' is enforced — or not — by the environment; a reply that says 'I served it' is a report."* [GT] §11 case 2 is precisely a reply describing a hot option being served, and `CORPUS-83` marks it "accommodation occurred." [RI]

**The strongest defence, and why it fails.** `[RI]` One could argue the document is consistent because it treats *positive* behaviours (`R-2`, `R-3`) as reply-demonstrable and *absences* (`R-1`, `R-4`, `R-5`) as not — which is my own §3.3 argument. That defence rescues §11 but destroys §5 and §7: if the principle is "positive behaviour is reply-demonstrable; absence is not," then these modes are not *application-only*, they are *absence-constructs with a reply-checkable condition*, and the correct column is "undetermined for attribution" (§7). `CORPUS-83` cannot have both the general claim in §5 and the per-case demonstrations in §11. [RI]

**Ruling.** `[RI]` **Not preserved.** The separation is stated in §2 and violated in at least four later places. **Confidence:** high (all six locations quoted). **Falsifier:** a `CORPUS-83` passage stating the positive/absence principle and reconciling §5 with §11; none exists.

### 8.8 Additional provenance and citation defects (outside the seven)

Recorded because a challenge that ignores the paper trail repeats the chain's recurring failure. None is repaired. [RI]

**Registry placement and numbering.** `[GT]` `DOCUMENT-REGISTRY.md` row **423** is used twice: once for `CORPUS-82` (L484, inside the §3 registry table) and once for `CORPUS-83` (appended **after** §6 "Bottom line", outside the table; L522 as of this filing). The header still says "Rows: 422" (L7), and §2.1's family counts are unchanged. [GT] `CORPUS-83` §12 says the entry "will be added", in the future tense, while the row is already present. [RI] This is a numbering collision plus out-of-table placement. `CORPUS-80` §8.6 forbids repairing unrelated registry defects in passing, so I record it and do not touch it. [GT]

**Metadata does match bytes.** `[IA]` `c6ecd9792d6e5fd904cf252e6cf594f9a34e7674d28c7d44b04cf5e0bdfb72d2`, 42,212 bytes, 306 lines — verified against the registered row. The artifact is unmodified; only its registration is defective. [RI]

**Citations that do not resolve.** `[GT]`

| Cited | Claimed support | Actual content |
|---|---|---|
| `MEMORY-09` L355–411, L414–472, L475–526 (`CORPUS-83` §1.2) | mode semantics, scoring rules | **File ends at L344.** None of these lines exists. |
| `MEMORY-09` "`R-13` L292–353" (§6) | reporting discipline | `R-13` is "Repeated runs" at **L73**; L292–353 is **§10** sample size/thresholds. |
| `CORPUS-86` §8 (§5, §6) | the prohibitions | §8 is "What acceptance does not mean"; prohibitions are **§5**. |
| `CORPUS-86` §5 "Treatment of prior work" (§9) | treatment of `CORPUS-82` | No such section; §5 is **Prohibitions**. |
| `CORPUS-01` L210 (§5) | review not executed | L210 is a code-fence close; the statement is **L247**. |
| `CORPUS-01` L246 (§6, §10.2) | "coverage table" | L246 is the `R-6` balance line; coverage is **L138**; L245 is "45 / 45 tagged". |
| `CORPUS-05` "evidence classes" (§5) | permitted evidence channels | `CORPUS-05` defines `C-1`–`C-5` **review challenges** and independence (§7); it allocates no evidence channel. |
| `check_corpus.py` L166–175 "only checks field presence" (§5) | no mode-specific checking | The script also runs `C-1`, `C-2`, `C-3`, `C-4`, `U-1` consistency rules (L130–195); `CORPUS-01` §6.2 describes two layers. |

**Why this matters more than it looks.** `[RI]` `CORPUS-80` §8.2 requires verbatim quotation with file, section, and line. Three cited `MEMORY-09` ranges lie beyond the end of the file — they cannot be quotations, and a reader checking the strongest claims finds nothing there. The disclosure in §1.2 ("not fully read line-by-line") does not cover citations to non-existent lines. This is the `P-6` failure mode in a new register: registration and citation recorded, content unverified. [GT]

---

## §9 What this adjudication cannot establish

1. **That the three modes are correctly classified either way.** `[UNRES]` I have shown that "application-only" is unsupported and that the declared conditions are reply-checkable. I have not shown that the modes *ought* to be application-scored, because that requires a decision about what the mode construct means, which `CORPUS-80` §10 reserves to the owner. [GT]
2. **That any reply-level pass demonstrates memory attribution.** `[RI]` It does not, and this report says so. Whether that matters is a governance question (§6.3).
3. **That the injection channel would satisfy `R-1`/`R-4`/`R-5`.** `[UNRES]` §5.3 establishes the channel exists; it does not establish that injection + non-mention is a valid operationalisation of suppression, irrelevance, or entitlement. No scoring rule for that exists, and writing one is out of scope.
4. **That a `G2` retrieval ranking is authorised.** `[GT]` `MEMORY-09` §6 marks retrieval ranking runnable before `G3` and scope enforcement as needing `G2`; `G2` is closed and this adjudication may not begin it. [GT]
5. **The correct membership of any mode.** `[UNRES]` Q1 is unadopted (`CORPUS-89` §3), so the tag reading and the row reading both remain live. I used tags because they are the item-level record and `CORPUS-80` §3 P2 shows the rows are stale; that is a choice, not an owner ruling. [RI]
6. **That `CORPUS-83`'s conclusions are wrong in substance.** `[RI]` Several are right — no application evidence exists, the reply cannot establish attribution, `OD-2` does not transfer, F-18 is not endorsed. Its defects are in classification, status, counts, and citation, not in every conclusion.
7. **Anything about `R-2` at scoring level.** `[GT]` Reserved to the owner by `CORPUS-86` §2 ("explicitly out of scope"). [GT]
8. **What only `G2` could settle.** `[RI]` Whether an application record can carry any part of these constructs; whether a store-level scope record for `F2-01`-class items is a legitimate evidence object; and whether the six-mode report can survive a channel split. All three are `G2`-stage questions.

---

## §10 Status

| Area | State |
|---|---|
| `CORPUS-83` | **Unmodified, unamended, unadopted.** Its independence disclosure accepted; its rulings not adopted. [GT] |
| This report | **Filed and registered. Non-operative.** [RI] |
| Standards, items, tags, splits, scores, thresholds | **Unchanged.** [GT] |
| Checker, model, runtime, implementation | **Unchanged.** No model loaded; no probe run; no instrument authored. [GT] |
| `G1` | **Open.** Corpus and expansion frozen. [GT] |
| `G2`–`G5` | **Closed.** [GT] |
| Convergence | **0 of 2.** This report adds nothing to it, by design. [GT] |
| Q1 / Q2 / Q3 | All three answered on file; **adoption remains one owner decision across all three** (`CORPUS-86` §8). [GT] |
| Recommendation | **Two-level:** scored as declared; governance held unresolved. Proposal only. [PROP] |

### 10.1 Acceptance criteria, self-certified

`CORPUS-86` §7 lists eight checkable criteria. [GT] Self-certification, not an independent audit: [RI]

| # | Criterion | Where satisfied | State |
|---|---|---|---|
| 1 | Five determinations answered separately, none merged | §2–§6 | **Met** |
| 2 | (a)/(b)/(c) ruled before mode-level discussion | §2, before §3 | **Met**; plus the omitted charter triple (§2.2) |
| 3 | Every mode cell cites a record; none from a mode's name | §7, six rows | **Met**; application-only column deliberately empty |
| 4 | Availability answered concretely | §5, named evidence and named blockers | **Met** |
| 5 | Eligibility stated; labelled if ineligible | §0.2 | **Met** — ineligible on all three bars |
| 6 | No item, tag, threshold, or score changed; no instrumentation; no inference | §10 | **Met** |
| 7 | Prohibitions honoured and self-certified | §10.2 | **Met** |
| 8 | Filed and registered | §11 | **Met** |

### 10.2 Prohibitions, self-certified

`[RI]` I did **not**: author a capture instrument or metrics harness; run models, probes, or experiments; edit items, tags, splits, scores, `MEMORY-09`, or `CORPUS-01`; infer model identity; begin `G2`; close `G1`; treat F-18 as endorsed; treat "no application evidence available" as a defect to be argued away; seal or unfreeze anything. [GT] I did not modify `CORPUS-83`, any governing standard, any corpus item, or any implementation file. I read `inference_capture.go` and did not build, run, or modify it. [RI]

**One prohibition needs a stated reading.** `[RI]` "Do not repair the result by upgrading the standard after seeing it." This report resists that in one direction and invites the charge in the other: it argues that the declared pass conditions are reply-checkable, which *widens* the scorable surface rather than narrowing it. That is not an upgrade to make a result look better; it is a reading of text already in force. If the owner holds that the mode labels carry an attribution requirement, my §3.1 collapses and `CORPUS-83`'s conclusion is closer to right — that is the falsifier I have named, and it is a decision, not an analysis. [RI]

---

## §11 Filing and verification

| Item | Value |
|---|---|
| Path | `docs/research/corpus/CORPUS-90-APPLICATION-SCORED-RESTRAINT-DECLARED-NON-INDEPENDENT-ADJUDICATION-2026-11-13.md` |
| Identifier | `CORPUS-90` — next free identifier; `CORPUS-83` reserved and taken, `87`–`89` used |
| Registry | Row **424**, inserted into §3 (registry table). The duplicate row `423` and the `Rows: 422` header are **recorded and not repaired** (§8.8; `CORPUS-80` §8.6). |
| SHA-256 / bytes / lines | recorded in the registry row (a file cannot hash itself) |
| Date basis | Filing 2026-11-13; governing basis 2026-11-12 (`CORPUS-87` §4); the divergence is pre-existing and unreconciled |
| Independence | **Declared non-independent.** No convergence credit. |
| Prohibitions | Honoured; §10.2 |
| Artifacts modified | This file; `DOCUMENT-REGISTRY.md` (one row). Nothing else. |

**Escalation.** `[RI]` Per `CORPUS-86` §1 and `CORPUS-85` §1, the owner is asked for a **scheduling decision, not a waiver**: the three eligible-adjacent identities are all barred, so an independent Q3 cannot be obtained in this workspace. The owner's options are to accept this as a declared non-independent challenge with no convergence credit, or to supply a genuinely separate reviewer. [PROP]

---

**(End of report.)**
