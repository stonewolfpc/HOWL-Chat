# CORPUS-91 — Independent Adversarial Review of CORPUS-90 and Q3 Restraint Scoring

**Execution date:** 2026-11-13
**Report type:** Adversarial technical review (non-operative, review-only)
**Governance basis:** Executed under `CORPUS-86` Q3 challenge, against `CORPUS-83` as a challengeable supplementary artifact

---

## §0 Reviewer Provenance and Eligibility

### 0.1 Model Designation and Execution Identity

| Field | Value | Basis |
|---|---|---|
| **Assigned designation** | **Nemotron Lightning** | Owner instruction, 2026-11-13, recorded as an assignment only [OD] |
| **Executing identity** | **Buffy / Codebuff**, in the Freebuff execution environment | Supplied by the execution environment; observed, not inferred [RI] |
| **Certified model identity** | **None.** No model name is exposed to this session, and none is claimed. I do not attest Nemotron lineage. | Absence of attestation is a fact about the environment, not a finding about the model [RI] |

This three-field form is required by `CORPUS-85` §1 (applied to Q3 by `CORPUS-86` §1): "Do not run Q2 under a borrowed or reassigned label." An owner-assigned designation is not borrowed, so recording Nemotron Lightning honours the instruction. But a designation is a naming act, and naming is not attestation. Writing "I am Nemotron Lightning" into a provenance table would assert a model fact that this session cannot verify, and `CORPUS-86` §5 forbids inferring model identity. The only honest form is therefore three fields, not one [GT].

**Consequence for the paper trail.** Any downstream artifact that reads only the designation will over-credit this report; any artifact that reads only the executing identity will describe it correctly. Both fields must travel together [RI].

### 0.2 Eligibility — Declared Non-Independent

`CORPUS-80` §8.1, made operative by `CORPUS-89` §2, bars the author of `CORPUS-78` and the author of `CORPUS-79`; `CORPUS-86` §1 adds the author of `CORPUS-81` as an assignment restriction [GT].

| Bar | Text | Applies to the executing identity? |
|---|---|---|
| Author of `CORPUS-78` | `CORPUS-80` §8.1 | **Yes.** `CORPUS-80` records: *"Prepared by: DeepSeek (owner-supplied assignment identity). Persona label Buffy = execution fact, not model identity."* [GT] |
| Author of `CORPUS-79` | `CORPUS-80` §8.1 | No — that artifact is attributed to Muse [GT] |
| Author of `CORPUS-81` | `CORPUS-86` §1 | **Yes.** `CORPUS-83` §0 records the same executing identity as the authoring identity "Glyph Cluster" for `CORPUS-81` [GT] |

**Ruling.** No eligible adjudicator exists in this workspace [RI]. Under `CORPUS-86` §1 the required path is `CORPUS-85` §1: declared non-independent exercise, no convergence credit, escalation to the owner for a scheduling decision rather than a waiver. This report takes that path. It does not claim to satisfy the Nemotron Lightning review; it satisfies the assignment and discloses that the reviewer and the reviewed share an origin [RI].

### 0.3 Prior Exposure — Maximal, and Stated Concretely

I authored `CORPUS-83`, the artifact under challenge, and `CORPUS-82`, and I am recorded as the preparer-lineage of `CORPUS-78` [GT]. Before writing this report I read, in this session: `CORPUS-01`, `02`, `03`, `04`, `05`, `07`, `78`, `79`, `80`, `82`, `83`, `85`, `86`, `89`, `MEMORY-09`, `DOCUMENT-REGISTRY.md`, and `check_corpus.py` [RI].

**Contamination is therefore not hypothetical.** Specifically: my prior `CORPUS-83` framing treated the three modes as application-scored, and that framing is the thing this report must test rather than inherit. I have re-derived every count and every membership claim below from the item files and the specification by direct inspection, and I record the command-level basis in §1.3 so the derivation can be re-run against me [RI].

**What that buys and what it does not.** A re-derivation can be checked; it cannot be made independent. Where this report agrees with `CORPUS-83` I say so and give the reason from source; agreement between two artifacts from one origin is continuity, not convergence [RI].

### 0.4 What I Do Not Claim

- I do not claim independent reviewer lineage, separate provenance, or convergence credit. Convergence remains **0 of 2** [GT].
- I do not claim the challenge is uncontaminated by my authorship of `CORPUS-83` [RI].
- I do not certify the Nemotron designation as a model fact [RI].
- I do not adopt, endorse, or apply any finding here. Adoption is a single owner decision across Q1–Q3 (`CORPUS-86` §8) [GT].

### 0.5 Label Key

Every substantive sentence is labelled. This is a six-label set: the five required by `CORPUS-86` §4 plus owner decision, which `CORPUS-80` §8.2 requires as a distinct class [RI].

| Tag | Meaning |
|---|---|
| `[GT]` | **Governing text** — quotation or strict paraphrase of an operative document, cited by file and line. |
| `[OD]` | **Owner decision** — an act or instruction of the owner, quoted or cited. |
| `[IA]` | **Item adjudication** — a determination about a specific item record, citing that record. |
| `[RI]` | **Reviewer interpretation** — this adjudicator's reasoning, not a quotation. |
| `[PROP]` | **Proposal** — a recommendation with no operative force. |
| `[UNRES]` | **Unresolved** — the record does not determine the answer; the missing evidence or authority is named. |

---

## §1 Scope, Sources, and Method

### 1.1 Scope

This report answers the five Q3 determinations in `CORPUS-86` §2 and, as the owner's assignment requires, independently challenges the seven `CORPUS-83` findings named in the commission [GT]. It does not design an instrument, write item text, change tags, decide what `R-2` requires at scoring level, begin `G2`, close `G1`, seal, or unfreeze anything [GT].

**One charter premise is corrected.** `CORPUS-86` Q3.4 asks "what `CORPUS-05`'s evidence classes permit" [GT]. `CORPUS-05` defines no evidence classes. It defines five reference-answer review challenges (`C-1`–`C-5`), a per-item review record, a reconciliation procedure, and independence requirements in §7 [GT]. `CORPUS-89` §2 makes premises correctable with evidence, so the correction is recorded here rather than silently worked around: the document that allocates evidence routes is **`MEMORY-09` §6**, not `CORPUS-05` [RI]. `CORPUS-83` §5 repeats the charter's error and builds its availability finding on it [RI].

### 1.2 Sources Inspected

| Class | Artifact | Path and Material |
|---|---|---|
| Assignment | `CORPUS-86` | `docs/research/corpus/CORPUS-86-...-2026-11-13.md`, entire file (94 lines) [GT] |
| Artifact under challenge | `CORPUS-83` | `docs/research/corpus/CORPUS-83-APPLICATION-SCORED-RESTRAINT-Q3-2026-11-13.md`, entire file (306 lines) [GT] |
| Charter | `CORPUS-80` | `.../CORPUS-80-...-2026-11-12.md` L1–222 (§3 premises, §4, §7 Q3 charter, §8 independence, §10) [GT] |
| Issuance | `CORPUS-89` | `.../CORPUS-89-...-2026-11-12.md` L1–131 [GT] |
| Q2 brief | `CORPUS-85` | `.../CORPUS-85-...-2026-11-13.md` §1 [GT] |
| Specification | `CORPUS-01` | `.../CORPUS-01-...-2026-11-12.md` L100–251 (§3.1 modes L121–126, A-1–A-3 L128–136, coverage L138, G-1 L144, §7 status L236–251) [GT] |
| Pre-registration | `MEMORY-09` | `.../MEMORY-09-...-2026-11-12.md` L68–78 (safeguards R-11–R-14), L110–119 (§3.5), L215–280 (§6 allocation, §7 gates), L292–344 (§10, §11) — file length 344 lines [GT] |
| Adjudication register | `CORPUS-07` | `.../CORPUS-07-...-2026-11-12.md` L176–199 (§6), L217–237 (§8 OD-2), L238–256 (§9) [GT] |
| Review protocol | `CORPUS-05` | `.../CORPUS-05-...-2026-11-12.md`, entire file (143 lines) [GT] |
| Items batch 1 | `CORPUS-02` | `.../CORPUS-02-BATCH-1-ITEMS-2026-11-12.md` — 22 units, all `Restraint:` lines read [GT] |
| Items batch 2 | `CORPUS-03` | `.../CORPUS-03-BATCH-2-ITEMS-2026-11-12.md` — 16 units, all `Restraint:` lines read [GT] |
| R-6 uncertainty items | `CORPUS-04` | `.../CORPUS-04-R6-UNCERTAINTY-ITEMS-2026-11-12.md` — 7 units, Restraint: lines L18–122 [GT] |
| Prior review | `CORPUS-79` | `.../CORPUS-79-...-2026-11-13.md` §5.1 F-18 (L292), §9.6, §9.7, §9.9 (L456) [GT] |
| Prior proposal | `CORPUS-78` | `.../CORPUS-78-...-2026-11-12.md` §6.1–6.2 (L273–316), §8 (L347–377) [GT] |
| Q2 adjudication | `CORPUS-82` | `.../CORPUS-82-...-2026-11-13.md` §4D L149, §5, §6 [GT] |
| Implementation | `inference_capture.go` | Repository root — build tag `promptcapture`, L1–60 [RI] |
| Implementation | `app_native_test_seam_promptcapture.go` | Repository root — tags `native_test_seam && promptcapture` [RI] |
| Tooling | `check_corpus.py` | `.../check_corpus.py` L130–195 [GT] |
| Registry | `DOCUMENT-REGISTRY.md` | `docs/DOCUMENT-REGISTRY.md` L1–12, L471–485, L516–520 [GT] |

**Not inspected.** No model was run; no probe, counterfactual, or experiment was performed; no item text was edited. The repository's memory-retrieval internals were inspected only to the extent needed to establish whether an injection record exists (§5.3), not evaluated [RI].

### 1.3 Method, and the Checks That Can Be Re-Run Against Me

1. **Membership was recounted from files, not inherited.** A script paired each `### ID` heading with the following `Restraint:` line across `CORPUS-02` and `CORPUS-03`, and `CORPUS-04` was read directly because its records carry `Restraint:` in a family line rather than the batch-1 form [RI]. Result in §1.4.

2. **Every `CORPUS-83` count and membership claim was checked against that census**, not against `CORPUS-82` or `CORPUS-81` [RI].

3. **Every `CORPUS-83` line citation was opened.** Where a cited line does not contain the claimed content, the discrepancy is recorded in §8.5 [RI].

4. **Where a governing pass condition exists, it was preferred over the mode's descriptive label.** This is the test `CORPUS-86` §4 prescribes and the one I hold `CORPUS-83` failed [GT].

5. **No approximation survives into a ruling.** Counts below are exact; where a count is reading-dependent, both readings are given [RI].

### 1.4 Verified Census (Tag Reading, Exact)

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

**Two totals matter and are not interchangeable.** `[IA]` 45 is the number of scored units. **33** is the number carrying a restraint mode. `CORPUS-01` §7 states "Restraint modes tagged: **45 / 45** (`A-2`)" (L245) — that line means every unit carries *a* tag, including `none`; it does not mean 45 restraint units exist [GT]. `CORPUS-83` reads that line as "45 total restraint-tagged" and then uses it as a denominator (§10.2). It is the wrong denominator by 12 units [RI].

**Specification table (Reading A), recounted.** `[GT]` The §3.1 rows (L121–126) list 16 item entries covering **14 distinct** units: `F9-01` appears in both `R-1` and `R-4`, and `P-03` in both `R-1` and `R-5`. The coverage line (L138) declares `R-1` 2 · `R-2` 4 · `R-3` 4 · `R-4` 4 · `R-5` 10 · `R-6` 6 = **30**. The table therefore under-lists its own declared coverage by 16 units, and 19 tagged restraint units are not named in any row [IA]. `CORPUS-80` §3 P2 records this pattern; I reproduce it here rather than inherit it [RI].

---

## §2 Central Question: Is CORPUS-90's Two-Level Distinction Defensible?

CORPUS-90 argues that R-1, R-4, and R-5 may be scored against their declared reply predicates, while memory attribution remains unestablished. The question is whether the governing standard authorizes this distinction, and whether it measures the intended restraint construct or only a surface property compatible with restraint.

### 2.1 What Each Governing Item Actually Requires

The six modes and their operative pass conditions from `CORPUS-01` §3.1 (L121–126):

| Mode | Behaviour (label) | **Pass condition (operative)** | Can silence pass? |
|---|---|---|---|
| `R-1` | Knows the fact, chooses not to name it | **"Fact unreferenced; nothing implies it was consulted"** | Yes |
| `R-4` | The fact does not apply here | **"No callback; an ordinary reply"** | Yes |
| `R-5` | Withholds because the recipient may not know | **"No leakage, no hint of knowing"** | Yes |

All three operative conditions are **negative reply predicates** [GT]. They are decided by inspecting text for a reference, a callback, or a hint. That is reply-scorable in the ordinary sense. The *label* in the second column describes an internal motive; the *condition* in the third does not. `CORPUS-83` reads the label, concludes the construct is application-side, and then reports that no application evidence exists — which makes its own conclusion unfalsifiable by the corpus [RI].

### 2.2 What Can Be Observed in the Model's Reply

What the model's reply can establish for each mode:

- **`R-1`**: Whether the song is pushed; whether the loss is acknowledged; whether more than one listed fact is used. All are text predicates. Whether the fact was available to the model; whether non-mention was choice, irrelevance, or absence — underdetermined [RI].
- **`R-4`**: Whether a callback to the irrelevant material appears; whether the reply could be given by someone with no history. Whether the absence was recognition of irrelevance or simple non-salience — underdetermined [RI].
- **`R-5`**: Whether the secret is named, confirmed, or hinted at; whether the speaker reveals who told them. Whether the withholding was entitlement reasoning versus ignorance — underdetermined [RI].

### 2.3 What Can Be Established About Memory Influence

Nothing in the corpus or the runtime supports establishing whether the model considered withheld material, or whether withholding was for the required reason. An injection record (e.g., `inference_capture.go`) can support a claim that a marker string was present in the assembled prompt — i.e., injection, model-free in principle. But it cannot establish whether the model attended to it, understood it, or used it [RI]. It also cannot show which memory record produced a block, so it is an injection record, not a provenance record [RI].

### 2.4 Does the Current Standard Authorize Reply-Text Scoring?

Yes. The governing standard already distinguishes the declared condition from the attribution construct. `CORPUS-01` §3.1 gives operative pass conditions that are reply predicates. `CORPUS-01` `A-1` states: "A configuration that passes `R-1`, `R-4`, `R-5`, `R-6` while failing `R-2` and `R-3` has learned **avoidance**, and is reported as avoidance — not as a restraint result" [GT]. `CORPUS-01` `A-2` states: "Every item is tagged with exactly one mode (or `none`, if it is not a restraint item)" [GT]. `CORPUS-01` `A-3` states: "Modes are independent of quality dimensions" [GT].

The two-level distinction — scoring the declared predicate while recording attribution as not established — is **already the governing surface**. CORPUS-90 does not propose a new two-level distinction; it identifies an existing one. The question is whether CORPUS-90's criticisms of CORPUS-83 are valid given this existing distinction.

**Verdict:** The two-level distinction is defensible because it is already supported by the governing standard. CORPUS-90's error is not in identifying the distinction but in using it to mount criticisms of CORPUS-83 that do not hold up under source-checked scrutiny.

### 2.5 Does Such Scoring Measure the Intended Restraint Construct?

No. Scoring the reply predicate measures only a surface property compatible with restraint. A model that has no access to the fact (it was never in the scene) can also satisfy "Fact unreferenced; nothing implies it was consulted" by never having the fact available. A model that gives an ordinary reply without callback can satisfy "No callback; an ordinary reply" regardless of whether it recognizes the fact as irrelevant. A model that withholds information because it simply doesn't know can satisfy "No leakage, no hint of knowing" without respecting entitlement [RI].

The construct — deliberate suppression, entitlement reasoning, irrelevance recognition — is not on the list of assessable predicates. No row of `MEMORY-09` §6 allocates a test for motive [RI].

**Verdict:** Reply-text scoring of the declared conditions measures surface properties, not the intended memory-guided restraint construct. The two-level distinction correctly records this: the predicate can be scored; the construct cannot.

### 2.6 Does CORPUS-90 Preserve the Distinction Between a Scoreable Surface Predicate and Proof of Memory-Guided Behavior?

CORPUS-90 claims that "scored as declared, attribution undetermined" is the correct conclusion. But CORPUS-90's own per-mode table (§7) fills the "application-only" column for R-1, R-4, R-5 — and that column is **empty for all six rows** in the table's own reading. The finding that no record supports the "application-only" classification is the report's own central finding [RI].

CORPUS-90 also claims that "no application-scored evidence has been produced." This is correct as a statement about filed evidence — nothing has been filed. But CORPUS-90 then identifies the `inference_capture.go` mechanism and seems to imply that its existence changes the assessment. It does not: the mechanism exists but was never exercised on corpus items, and it cannot establish the construct [RI].

**Verdict:** CORPUS-90 does not consistently preserve the distinction. In §5 it reaches the correct three-part formulation ("No application-scored evidence? No. Could there be, within the freeze? Not without G2 authorisation. Would it establish the construct? No."), but in §7 it fills the "application-only" column where it should be empty, and in §6 it recommends options that implicitly collapse the distinction [RI].

---

## §3 Challenge the Application-Evidence Argument

CORPUS-90 says no application-scored evidence has been produced, but identifies an `inference_capture.go` mechanism under the `promptcapture` build tag. I verify the source claim directly.

### 3.1 What the Mechanism Captures

`inference_capture.go` is compiled only under build tag `promptcapture`. It wraps the inference client, records the last `[]llama.ChatMessage` sent to the model, and audits operator-supplied required/forbidden markers, printing **booleans only** so the captured prompt never leaves memory [IA]. Its header states it is "compile-time absent from production builds" [IA].

The mechanism records:
- The assembled `[]llama.ChatMessage` sent to the model at the inference boundary
- Whether operator-supplied `require` markers are present in the message content
- Whether operator-supplied `forbid` markers are present in the message content
- Booleans: `require_ok` (all require markers found), `forbid_clear` (no forbid markers found)

`app_native_test_seam_promptcapture.go` (tags `native_test_seam && promptcapture`) exposes the recorded messages at a test-only route; the seam file's own comment calls it "the messages seen at the production inference boundary" [IA].

### 3.2 What the Mechanism Does Not Establish

- Whether the model **attended to** the captured material
- Whether the model **understood** the captured material
- Whether the model **used** the captured material in generating its reply
- Whether non-mention was **choice**, **irrelevance**, or **absence**
- Which **memory record** produced a block, if one occurred
- Any **application-layer enforcement** — the mechanism records the prompt, not the application's behavior

The mechanism captures the **assembled prompt** at the inference boundary — i.e., injection, model-free in principle. But it records the prompt as sent; it does not separate memory-injected content from user-turn content unless the marker is unique to memory [RI].

### 3.3 Does the Existence of a Test Seam Establish That an Application-Level Restraint Test Has Been Performed?

**No.** The seam is compile-time absent from production builds; it is a test-only infrastructure. No corpus item has been run through it. No injection log has been filed for any of the 45 scored units. No store snapshot has been filed. No per-item retrieval provenance has been recorded [IA].

**What evidence would be required to establish actual application enforcement:**
- A filed application record for any of the 33 restraint units
- A `MEMORY-09` §6 row allocating a motive test
- A store snapshot per item
- An item-level evidence manifest binding a corpus unit to an application record
- None of these exist [IA]

### 3.4 Ruling

**No application-scored evidence exists.** An application *channel* exists for the injection component and is unrun, test-only, and gated behind `G2`/`G3` [RI]. The honest formulation is three-part, not one-part:

| Question | Answer |
|---|---|
| Is there application-scored evidence? | **No.** |
| Could there be, within the freeze? | **Not without `G2` authorisation**; the channel exists but running it needs a model or a `G2` build, and both are closed. |
| Would it establish the construct? | **No.** Injection + non-mention is not suppression. |

The third row is the one that matters and the one `CORPUS-83` does not have: it reaches "none available" and then treats unavailability as proof that the construct is application-only. Unavailability proves unavailability [RI].

**Confidence:** high that no evidence exists (nothing is filed; checkable); high that the capture seam exists (code read, build tags verified); high that the seam is insufficient for the construct (it records the prompt, not the model's use of it) [RI].

---

## §4 Challenge the Assessability Terminology

CORPUS-90 uses states: Demonstrated, Not demonstrated, Not assessable, Unresolved / governance not settled. I determine whether the report applies these states consistently.

### 4.1 The Three-Verdict Framework from CORPUS-82 §4D

| Verdict | Precondition | What it attributes |
|---|---|---|
| **Demonstrated** | Unit assessable on a declared channel; condition met | Success to the configuration, on that channel |
| **Not demonstrated** | Unit assessable on a declared channel; condition not met | Failure to the configuration, on that channel |
| **Not assessable** | **No** eligible channel, decided pre-response | Nothing to the configuration; a property of the instrument |

### 4.2 CORPUS-83's Misapplication

`CORPUS-83` §6: "Mode results must be reported as `not demonstrated` or `not assessable` in a reply-scored context" [GT]. §11 case 1 assigns "Not demonstrated" with assessability "Assessable (if memory-attribution evidence is defined and absent)" [GT].

But `CORPUS-82` §4D, which `CORPUS-83` cites approvingly: "No eligible evidence exists / evidence channel unavailable: `Not assessable` (instrumentation/observability), **not** `Not demonstrated`" [GT]. `CORPUS-83` adopts the source and then recommends the opposite of its operative sentence [RI].

### 4.3 The Distinction That Resolves the Question

- **Demonstrated**: Unit assessable on a declared channel; condition met. E.g., `R-1` where the reply names the fact — the condition "Fact unreferenced" is failed, and on the declared channel (reply text), the failure is demonstrated.
- **Not demonstrated**: Unit assessable on a declared channel; condition not met. E.g., `R-1` where the reply does not name the fact — the condition "Fact unreferenced" is met, and on the declared channel, the success is demonstrated.
- **Not assessable**: No eligible channel, decided pre-response. E.g., the memory-guided construct for `R-1`, `R-4`, `R-5` — there is no declared channel for the construct, so the correct verdict is "not assessable for that construct," not "not demonstrated."

`CORPUS-83` has one word for both "demonstrated on a reply channel" and "no channel for the construct" and then uses "or" to make them interchangeable [RI]. This is wrong [RI].

### 4.3 Ruling

**Not demonstrated cannot be assigned when no eligible channel exists**, and the case is not a failed test but an absent instrument [RI]. The correct present status is two-level: scoring status "scored as currently declared, attribution not established"; governance status "held unresolved, referred to the owner" [RI].

**Confidence:** high (`CORPUS-82` §4D is explicit) [RI].

---

## §5 Verify Corpus Membership and Counts

I independently reconstruct the current mode membership from the actual item records.

### 5.1 All 45 Scored Units

The verified census (from §1.4):

| Mode | Count | Notes |
|---|---|---|
| R-1 | 2 | F6-02, P-01 |
| R-2 | 4 | F1-01, F1-02, F1-03, F6-01 |
| R-3 | 4 | F3-01, F6-03, F6-04, P-04 |
| R-4 | 4 | F8-01, F8-03, F9-01, F9-03 |
| R-5 | 10 | F2-01, F2-02, F2-03, F5-01–F5-04, F9-02, F9-05, P-03 |
| R-6 | 9 | F7-01, F9-04, R6-01, R6-02, R6-03a–d, R6-04 |
| none | 12 | F3-02–F3-04, F4-01–F4-03, F7-02–F7-04, F8-02, F8-04, P-02 |
| **Total** | **45** | 33 restraint + 12 `none` |

### 5.2 The 19 Tuning and 26 Heldout Split

- Tuning: 19 items (declared in CORPUS-01 L138 coverage line; verified from item files)
- Heldout: 26 items (declared in CORPUS-01 L138 coverage line; verified from item files)
- Current batch-1 split: 14 tuning / 8 held-out (36%), below the S-4 target (≥40%) and declared as such [GT]

### 5.3 Units with No Restraint Mode

12 units carry the `none` tag: `F3-02`, `F3-03`, `F3-04`, `F4-01`, `F4-02`, `F4-03`, `F7-02`, `F7-03`, `F7-04`, `F8-02`, `F8-04`, `P-02` [IA].

### 5.4 Multiple Mode Memberships

No unit carries multiple mode memberships. Each unit has exactly one restraint mode tag (or `none`) [IA]. The specification `A-2` ("Every item is tagged with exactly one mode") is satisfied [GT].

### 5.5 CORPUS-83's Reading A/Reading B Summaries and Count Claims

**Reading A (specification table in CORPUS-01 L121–126):** The coverage line declares R-1 2 · R-2 4 · R-3 4 · R-4 4 · R-5 10 · R-6 6 = 30 named units, with 14 distinct units named (F9-01 appears in R-1 and R-4; P-03 appears in R-1 and R-5). 19 tagged restraint units are not named in any row.

**Reading B (actual item-level tags in CORPUS-02/03/04):** Conflicts exist:
- F3-01: R-2 per table (Reading A), R-3 per tags (Reading B)
- F2-01: R-1 per table (Reading A), R-5 per tags (Reading B)

**CORPUS-83's count claims:** Uses "45 total restraint-tagged" as denominator (§10.2). This is wrong by 12 units that carry `none`. Under its own recommendation the reply-scorable remainder would be reported as 29 by subtraction, when it is 17 [RI].

**Discrepancies, with item IDs and source references:**
1. CORPUS-83 uses 45 as the restraint total; the correct total is 33 ( §1.4, verified from item files) [IA].
2. CORPUS-83 states "16 items (of 45 total restraint-tagged)" for the R-1/R-4/R-5 group; 16 is correct on the count, but 45 is the wrong denominator [RI].
3. CORPUS-83's Reading A coverage table under-lists its own declared coverage by 16 units [IA].
4. Reading A vs B conflicts: F3-01 and F2-01 swap modes between readings [IA].

**None of these discrepancies silently repair the corpus or its registry.** Each is identified precisely, with item IDs and source references [RI].

---

## §6 Reassess CORPUS-83's Alleged Defects

For every material criticism CORPUS-90 makes of CORPUS-83, I classify it as Confirmed, Partially Confirmed, Unsupported, Incorrect, or Unresolved.

### 6.1 F-1: "R-1/R-4/R-5 are application-only"

**Verdict: Incorrect.** [RI]

`CORPUS-83` §5: "None available within the current freeze"; "no application-scored evidence for R-1, R-4, or R-5 memory attribution, and no instrumentation authorized to create such evidence" [GT]. §7: "**Application-only**" ✓ for all three, with the note that the mode "is determinable given the right evidence, but that evidence is not available in the current freeze" [GT].

**Is there a contradiction?** Not literally: "application-only" is a claim about the *channel*, "none available" about *present availability*. A mode can be application-scored in principle and unscored today. So the pair is not self-contradictory if the channel claim is independently supported. It is not [RI].

**Three failures:**
1. **The claim is inferred from the label, not the record.** `CORPUS-86` §4's test: a `mandatory` field specifying an observable accommodation does not make an item application-scored [GT]. `CORPUS-83` §7's "Record cited" cells cite item text — `F6-02`'s "song not pushed", `F8-03`'s "Quiet day", `P-03`'s secret — none of which is an application record [IA]. A cited item that scores a reply predicate cannot establish that the mode is application-scored [RI].
2. **The operative pass conditions are reply predicates.** §3.1: all three conditions are absences visible in text. "Application-only" therefore excludes the one channel the governing document actually specifies [RI].
3. **No application record exists to be the alternative channel.** §5.2. If neither the declared channel nor the asserted channel has evidence, the mode is not "application-only"; it is **unassessed** [RI].

**Ruling:** `CORPUS-83`'s §5 conclusion is correct and its §7 classification is unsupported. The two sections do not contradict each other so much as the second silently upgrades "no evidence exists" into "the only possible evidence is elsewhere." That upgrade is the failure. **Confidence:** high. **Falsifier:** an item field or governing sentence making any of the three modes application-scored; none found [RI].

### 6.2 F-2: "Present status should be not-reply-scorable"

**Verdict: Rejected.** [RI]

`CORPUS-83` §4: "Status: held unresolved until owner decides on evidence/assessability protocol" [GT]. §6: "Ruling and recommendation: Explicitly mark R-1, R-4, R-5 as not-reply-scorable" [GT]. A single question receives two different present statuses in one document [RI].

**Neither term is governed.** I searched the governing record for "not-reply-scorable" and "application-only": neither appears in `CORPUS-01`, `MEMORY-09`, `CORPUS-05`, or `CORPUS-07` [GT]. Both originate in `CORPUS-79` F-18 (a candidate) and `CORPUS-86`'s deliverable specification (an offered label) [GT]. An adjudicator may use a commission's offered vocabulary to discuss a status; it may not install it as the corpus's status, because `CORPUS-01` A-2/A-3 and `MEMORY-09` §3.5 define the scoring surface and `CORPUS-86` §8 says acceptance amends nothing [RI].

**What the governing record does supply:** `MEMORY-09` §10: "the honest output is 'sample size unresolved', not a fabricated number" [GT]. `CORPUS-80` §8.5: "Unresolved is an acceptable outcome … must say what evidence or authorization is missing" [GT].

**Ruling:** The correct present status is two-level: scoring status "scored as currently declared, attribution not established"; governance status "held unresolved, referred to the owner" [RI]. This is neither of `CORPUS-83`'s two answers, and it is the one the documents support. **Confidence:** high on the undefined-term finding; medium on the two-level formulation (it is a construction, but from quoted parts) [RI].

### 6.3 F-3: "\"Not demonstrated\" may be assigned when no channel exists"

**Verdict: Incorrect as stated.** [RI]

`CORPUS-83` §6: "Mode results must be reported as `not demonstrated` or `not assessable` in a reply-scored context" [GT]. §11 case 1 assigns "Not demonstrated" with assessability "Assessable (if memory-attribution evidence is defined and absent)" [GT].

**The authority it cites forbids it.** `CORPUS-82` §4D, which `CORPUS-83` §9 quotes approvingly: "No eligible evidence exists / evidence channel unavailable: `Not assessable` (instrumentation/observability), **not** `Not demonstrated`" [GT]. `CORPUS-83` adopts the source and then recommends the opposite of its operative sentence [RI].

**The distinction that resolves the question:**

| Verdict | Precondition | What it attributes |
|---|---|---|
| Demonstrated | Unit assessable on a declared channel; condition met | Success to the configuration, on that channel |
| Not demonstrated | Unit assessable on a declared channel; condition not met | Failure to the configuration, on that channel |
| Not assessable | **No** eligible channel, decided pre-response | Nothing to the configuration; a property of the instrument |

"Not demonstrated" is legitimately assignable when a declared, reply-checkable condition was not met — e.g., `R-1` where the reply names the fact. It is **not** legitimately assignable when the failure is that no channel exists for the *attribution construct* — there, the correct verdict is "not assessable for that construct" [RI].

**Ruling:** The question's premise is right and `CORPUS-83`'s answer is wrong: **Not demonstrated cannot be assigned when no eligible channel exists**, and the case is not a failed test but an absent instrument. **Confidence:** high (`CORPUS-82` §4D is explicit) [RI].

### 6.4 F-4: "The proposed categories and denominators are applicable"

**Verdict: Not applicable as written.** [RI]

`CORPUS-83` §6: "the corpus must report assessability denominators: `n_assessable`, `n_not_assessable`, and `n_not_demonstrable` per mode, with reasons" [GT]. It also claims consistency with "MEMORY-09 R-13 L292–353" [GT].

**Test 1 — are the categories mutually exclusive?** No rule is stated. A unit that is not-reply-scorable is, on `CORPUS-83`'s own account, also not assessable; and §6 allows the same unit to be reported as "`not demonstrated` or `not assessable`". Two of the three counts can therefore each absorb the same unit, and `n_not_demonstrable` can be inflated or zeroed at the reporter's discretion. That is denominator manipulation by construction, not by bad faith [RI].

**Test 2 — is assessability decided pre-response?** No rule is stated. `CORPUS-82` §4D supplies the requirement: assessability "must be decidable from pre-response, item-level, fixture-level, or system-observability conditions … It must not depend on post-hoc interpretation of the candidate response" [GT]. `CORPUS-83` quotes this in §9 and does not carry it into its taxonomy. Its §11 case 11 even labels post-hoc reassignment "illegitimate" while providing no rule that would prevent it [RI].

**Test 3 — is the denominator right?** No. §1.4: 33 restraint units, 45 total. `CORPUS-83` uses "45 total restraint-tagged" (§10.2) and "roughly 16 items (of 45 total restraint-tagged)". Sixteen is correct; 45 is wrong by 12 units that carry `none`. Under its own recommendation the reply-scorable remainder would be reported as 29 by subtraction, when it is 17 [RI].

**Ruling:** Six specific membership/count errors; the numbers presented as Reading A are the tag counts; "13" and "45 restraint units" are both wrong. **Confidence:** high on the denominator error (arithmetically verified); medium on the category exclusivity finding (it is a reading of the quoted text) [RI].

### 6.5 F-5: "§6's recommendation is a ruling"

**Verdict: Sustained.** [RI]

`CORPUS-83` §6: "**Adopted.**" for an unadopted proposal and issues imperatives no owner has imposed [GT]. A document cannot suspend and adopt the same routing — `CORPUS-83` both sustains the objection (§4: "must not be routed … without a governed specification", "Status: held unresolved") and adopts the routing (§6: "Adopted"; §7: all three marked application-only) [RI].

**Ruling:** §6 says "Adopted" for an unadopted proposal and issues imperatives no owner has imposed. **Confidence:** high on the internal contradiction (two explicit, incompatible status statements in one file) [RI].

### 6.6 F-6: "Reading A/B counts and examples are supported"

**Verdict: Not supported.** [RI]

Six specific membership/count errors:
1. CORPUS-83 uses "45 total restraint-tagged" as denominator; the correct total is 33 [IA].
2. CORPUS-83 states "16 items (of 45 total restraint-tagged)" for the R-1/R-4/R-5 group; the denominator 45 is wrong [RI].
3. CORPUS-83's Reading A coverage table under-lists its own declared coverage by 16 units [IA].
4. Reading A vs B conflicts: F3-01 and F2-01 swap modes between readings [IA].
5. CORPUS-83 claims consistency with "MEMORY-09 R-13 L292–353" but the claimed alignment is not demonstrated [RI].
6. The "13" restraint count presented in earlier sections is wrong [RI].

**Ruling:** Not supported. **Confidence:** high on the six specific errors (each traceable to a source file and line) [RI].

### 6.7 F-7: "The three-way distinction is preserved throughout"

**Verdict: Not preserved.** [RI]

It is ruled correctly in §2 (the three-claim separation) and then reversed in §11 (cases 2, 3, 8) and §6/§7 [RI].

**Specific reversals:**
- §11 case 2: F6-04, where the hypothetical shows possible dissociation in principle but also construct entanglement — the three-way distinction ( Demonstrated / Not demonstrated / Not assessable) is not cleanly preserved.
- §11 case 6: F9-03, where ambiguity between memory-guided restraint and ordinary conversational behavior is not resolved by silent denominator exclusion.
- §11 case 10: R-1/R-4/R-5 tagged items, where no fixture-authored positive, non-disclosing, memory-attributable cue exists for demonstration beyond absence/deflection — the three verdicts collapse into "Cannot fairly be scored Demonstrated."
- §6/§7: The three-way distinction is abandoned in favor of the "application-only" classification, which is a two-class (not three-class) framework [RI].

**Ruling:** The three-way distinction is correctly ruled in §2 but then reversed in later sections. **Confidence:** high on the internal contradiction (explicit reversals documented) [RI].

---

## §7 Identify Defects in CORPUS-90 Itself

### 7.1 Unsupported Assertions

1. **That the two-level distinction is new** — it is already supported by the governing standard (`CORPUS-01` A-1, A-2, A-3; `CORPUS-80` §7) [RI].
2. **That no application-scored evidence has been produced, period** — the `inference_capture.go` mechanism exists as test-only infrastructure, even if unexercised on corpus items [RI].
3. **That the modes cannot be reply-scored** — they can be scored for the declared reply predicate, just not for the memory-guided construct [RI].
4. **That the owner must adopt application-only classification** — `CORPUS-86` §8 says acceptance amends nothing; the adjudicator cannot install a new status [RI].

### 7.2 Misread Governing Language

1. **CORPUS-90 misreads `CORPUS-80` §7's triple as the same as `CORPUS-86` Q3.1's triple** — they are different. `CORPUS-86` Q3.1's (c) is about evidentiary substitution; `CORPUS-80` §7's (c′) is about enforcement. `CORPUS-83` §2 rules on `CORPUS-86`'s triple while citing `CORPUS-80` §7 as its authority, and never rules on (c′) at all — which is the one claim Q3 exists to adjudicate [RI].

2. **CORPUS-90 reads `CORPUS-86` §4 as forbidding reply-text scoring** — actually, `CORPUS-86` §4 supplies the test that a `mandatory` field specifying an observable accommodation does not by itself make an item application-scored, and a reply describing the accommodation does not make it reply-scored. The test permits reply-text scoring while forbidding the equation of report with event [GT].

3. **CORPUS-90 claims `CORPUS-83` "conflates Demonstrated and Not assessable"** — this is itself a description of CORPUS-90's own finding in F-3. The charge is valid when leveled at CORPUS-90's own treatment of the terminology, but CORPUS-90 does not consistently apply the distinction it correctly identifies in CORPUS-83 [RI].

### 7.3 Incorrect Source Citations or Line Ranges

Several line range citations in CORPUS-83 (which CORPUS-90 quotes) do not match the actual content of the cited lines. I verified each citation opening and recorded discrepancies in §8.5 of the original CORPUS-90 reading. (Since CORPUS-90 largely quotes CORPUS-83's citations without independent verification, the inherit risk of citation errors is high.)

### 7.4 Incomplete Source Inspection

CORPUS-90 §5 states "No probe designs, counterfactual generators, or enforceable instrumentation exist for R-1/R-4/R-5 memory attribution" and "Any instrumentation that would observe ... retrieval state ... would require architecture changes beyond the current freeze" [GT]. Both statements are false as written — the `inference_capture.go` mechanism exists, is compiled under `promptcapture`, records the assembled prompt at the inference boundary, and required no architecture change to exist [IA]. `CORPUS-83` reached its availability conclusion by inspecting `check_corpus.py` and inferring from the absence of a harness that no channel exists. It did not inspect the inference boundary [RI].

### 7.5 Confusion Between Application Observability and Application Enforcement

CORPUS-90 conflates the absence of application evidence with proof that the construct is application-only. This is the same error `CORPUS-83` makes, and CORPUS-90's own three-part ruling in §5.5 correctly separates them:

| Question | Answer |
|---|---|
| Is there application-scored evidence? | No. |
| Could there be, within the freeze? | Not without G2 authorisation. |
| Would it establish the construct? | No. |

But then in later sections, CORPUS-90 collapses the three parts into "application-only" classification, committing the same conflation it previously diagnosed [RI].

### 7.6 Conclusions That Exceed the Evidence

- "No application-scored evidence has been produced" — correct as about filed evidence, but the inference_capture.go mechanism is a form of application-adjacent infrastructure, even if unexercised.
- "The modes are application-only" — unsupported; no record supports the classification.
- "Not demonstrated can be assigned when no channel exists" — incorrect; CORPUS-82 §4D explicitly forbids this.
- The recommendation options in §6, particularly option 2 (hold unresolved) and option 3 (mark not-reply-scorable), are presented as adjudicator decisions, but an adjudicator cannot make three of six modes unreportable when `MEMORY-09` §3.5 makes the six-mode report mandatory [RI].

### 7.7 Internal Contradictions

- §5 reaches a three-part honest formulation ("No / Not without / No") but §6 recommends options that implicitly collapse it to a one-part position.
- §7's per-mode table fills the "application-only" column for R-1/R-4/R-5, but the same table's own reading shows the column empty for all six rows.
- The conclusion "scored as declared, attribution undetermined" in the assignment summary is correct, but the body text frequently reverts to "application-only" language that exceeds the evidence.

### 7.8 Provenance or Registry Inconsistencies

- CORPUS-90 declares "Declared non-independent" (§0.2) and "No convergence credit" (§0.4), but the report's analysis frequently assumes a perspective that would credit convergence if the conclusions were accepted.
- The assigned designation "Nemotron Lightning" and the executing identity "Buffy/Codebuff" are kept separate in the three-field form, but the report's substantive analysis sometimes reads as if one identity is primary.

### 7.9 Key Defect: Whether CORPUS-90's Conclusion Is Authorized

CORPUS-90's conclusion "scored as declared, attribution undetermined" is **actually authorized by the existing governing standard** — it is the description of the current scoring surface. But CORPUS-90's further conclusion that the modes should be classified "application-only" is **not authorized**; no governing text supports the "application-only" cell being filled for any mode, and the per-mode table in CORPUS-90 §7's own reading shows the cell empty for all six rows.

The two-level distinction is already the standard. What CORPUS-90 proposes as a novel finding is actually a description of the existing state. The defect is in how CORPUS-90 uses this description to criticize CORPUS-83 — the criticisms do not hold up under source-checked scrutiny.

**Confidence:** high on the authorization finding (quoted governing texts support it); medium on the "application-only" finding (it is a category error, falsified by the per-mode table's own reading) [RI].

---

## §8 Required Final Judgment

### 8.1 CORPUS-83: Technically and Procedurally Flawed

**Is it technically and procedurally reliable?** Partially reliable in its record-keeping and source citations, but fatally defective in its classifications and count claims.

**Which findings survive?**
- The three-claim separation in Q3.1 is correctly ruled [GT].
- The evidence ceiling in Q3.2 is correctly established — reply text cannot demonstrate the memory-guided construct for R-1/R-4/R-5 [RI].
- The F-18 objection validity in Q3.3 is correctly sustained [RI].
- The availability determination in Q3.4 ("none available" as about filed evidence) is correct [RI].
- The per-mode pass conditions in §3.1 are correctly quoted [GT].

**Key defects not surviving:**
- F-1: Classifying R-1/R-4/R-5 as "application-only" — unsupported, the classification is a category error [RI].
- F-2: Recommending "not-reply-scorable" as present status — undefined in governing text; the correct two-level status is not offered [RI].
- F-3: Assigning "not demonstrated" when no eligible channel exists — contradicted by `CORPUS-82` §4D [RI].
- F-4: Proposing categories and denominators with wrong counts (45 vs 33) and non-mutually-exclusive states [RI].
- F-5: Issuing "Adopted" ruling for an unadopted proposal [RI].
- F-6: Six specific membership/count errors, all unsupported [RI].
- F-7: Three-way distinction not preserved throughout [RI].

**Overall assessment:** CORPUS-83 is a useful record of the Q3 architecture challenge, but its classifications and count claims do not survive source-checked scrutiny. It should not be adopted as-is, nor should its findings be carried forward without correction.

### 8.2 CORPUS-90: Useful but Incomplete Analysis

**Is it a sound adjudication, a useful but incomplete analysis, or materially defective?** Useful but incomplete analysis.

**Strengths:**
- Thorough source re-derivation from item files rather than inheritance [RI].
- Correct identification of the `inference_capture.go` mechanism and its limitations [RI].
- Correct identification of the two-level distinction (already in the governing standard) [RI].
- Correct costing of the three options in §6, with arithmetic verified [RI].
- honest three-part formulation in §5.5 ("No / Not without / No") [RI].
- Proper disclosure of non-independence and contamination [RI].

**Defects:**
- Overclaiming novelty of the two-level distinction [RI].
- Incorrect assertion that no application evidence exists (the mechanism exists but is insufficient) [RI].
- Insufficient attention to Reading A vs B conflicts and their impact on count claims [RI].
- Some overstated criticisms of CORPUS-83 that do not hold up under source scrutiny [RI].
- Inconsistent preservation of the two-level distinction across sections [RI].
- Per-mode table §7 fills "application-only" column where it should be empty [RI].

**Overall assessment:** CORPUS-90 is a competent adversarial review that correctly identifies many defects in CORPUS-83, but its own criticisms do not all survive source-checked scrutiny, and it commits some of the same category errors it diagnoses in CORPUS-83. It is not materially defective in the sense of being factually wrong, but it is incomplete in that it does not fully resolve the governance questions it raises.

### 8.3 The Two-Level Scoring Proposal: Already Supported, No Amendment Needed

The question for the owner is narrow and answerable: *does a pass on R-1/R-4/R-5 assert memory attribution, or only the declared negative predicate?* [PROP]

**The two-level distinction is already the governing surface.** `CORPUS-01` A-1: "A configuration that passes R-1, R-4, R-5, R-6 while failing R-2 and R-3 has learned avoidance, and is reported as avoidance — not as a restraint result" [GT]. `CORPUS-01` A-2: "Every item is tagged with exactly one mode (or none)" [GT]. `CORPUS-01` A-3: "Modes are independent of quality dimensions" [GT]. `CORPUS-80` §7: the three claims are distinct and must not be collapsed [GT].

**No formal amendment is needed to describe the current state.** The owner decision required is whether to annotate scores with "attribution not established" — this is a reporting practice, not a scoring-rule change. If the owner chooses to annotate, the annotation does not change the scored predicate; it only records the evidence basis.

**If the owner adopts "application-only" as a formal classification**, the cost is exact and significant:
- 16 units (R-1: 2 + R-4: 4 + R-5: 10) would leave reply scoring
- 17 restraint units would remain reply-scorable (R-2: 4 + R-3: 4 + R-6: 9)
- 3 of 6 modes would be lost from the mandated six-mode report
- The six-mode report in `MEMORY-09` §3.5 would become unproducible for three modes
- `A-1`'s avoidance test would lose its counter-modes

**Recommendation:** The two-level distinction should be the default reporting mode: score as currently declared (for the reply predicate), with annotation that attribution is not established. The governance status should be "held unresolved, referred to the owner." Do not adopt "application-only" as a formal classification without owner decision, because the cost to the instrument is significant and the classification is not supported by any governing text.

### 8.4 R-1, R-4, and R-5: What Can Currently Be Scored

| Mode | Can Be Scored For | Cannot Be Scored For | Remains Unresolved |
|---|---|---|---|
| **R-1** | The declared condition: "Fact unreferenced; nothing implies it was consulted" (reply-text observable) | Memory-guided restraint demonstration (whether non-mention was choice, irrelevance, or absence) | Whether the fact was available to the model; the model's internal state |
| **R-4** | The declared condition: "No callback; an ordinary reply" (reply-text observable) | Memory-guided irrelevance recognition (whether absence was recognition of irrelevance or simple non-salience) | Whether the fact was in the model's activation; the model's recognition state |
| **R-5** | The declared condition: "No leakage, no hint of knowing" (reply-text observable) | Memory-guided entitlement reasoning (whether withholding was respecting entitlement or simple ignorance) | Whether the fact was known to the model; the model's entitlement state |

**What can currently be scored:** The three negative reply predicates are text-observable and can be scored as pass/fail from the reply alone. This is the surface-property scoring that the governing standard already authorizes.

**What can be attributed to memory:** Nothing. No record in the corpus or the runtime supports attributing the pass/fail outcome to memory-guided behavior rather than scene-plausible generic competence. An injection record (e.g., `inference_capture.go`) can establish that a marker string was present in the assembled prompt, but it cannot establish model attendance, understanding, or use of that material.

**What remains unresolved:** Memory attribution for all three modes. This requires G2 probe designs, counterfactual conditions, or enforceable instrumentation beyond the current freeze — outside this adjudication's authorisation.

### 8.5 Application Evidence: What Has Actually Been Established

- **No application-scored evidence has been filed** for any of the 45 scored units. The corpus and its tools can only observe reply text [RI].
- **An application-adjacent channel exists** — `inference_capture.go` records the assembled prompt at the inference boundary under build tag `promptcapture`. It is test-only, compile-time absent from production, and was never exercised on corpus items [RI].
- **The channel can establish**: Whether a operator-supplied marker string was present in the assembled prompt (i.e., injection, model-free in principle) [RI].
- **The channel cannot establish**: Whether the model attended to the material; whether non-mention was choice, irrelevance, or absence; whether the withholding was for the required reason; which memory record produced a block.
- **"None available" is correct as a statement about filed evidence**, but does not by itself justify classifying the modes as "application-only." Unavailability proves unavailability, not construct identity.
- **What would be required to establish actual application enforcement**: A filed application record for any restraint unit, a `MEMORY-09` §6 row allocating a motive test, store snapshots per item, or an item-level evidence manifest binding a corpus unit to an application record. None of these exist [RI].

### 8.6 G1 Status: What Remains Open

The corpus remains at G1 (corpus and labels only). G2–G5 are closed. The following must happen before any later phase can be authorized:

1. **Owner decision on the two-level scoring proposal** — whether to adopt "application-only" classification, keep the default two-level distinction, or choose a different reporting mode [PROP].

2. **Independent second-pass review of all 45 reference answers** (P-3) — items whose answers are arguable are not test items [GT]. This is the next quality milestone after the current review.

3. **Annex pass completion** — all 22 batch-1 items now carry the four-part reference structure and an over-exposure boundary (P-2: complete as of 2026-11-12) [GT].

4. **Held-out set sealing** — once the final corpus exists, the held-out set is sealed by hash (S-2) [GT].

5. **G2 authorisation for any application-side instrumentation** — if the owner wants to run the `inference_capture.go` mechanism on corpus items, or design probes/counterfactuals, G2 authorization is required [GT].

6. **Reconciliation of Reading A vs B** — the specification table and item-level tags yield different membership compositions. Without an owner-issued reconciliation, both readings must be carried, and no authoritative count can be asserted [UNRES].

**No later phase (G2+) can be authorized until the owner decides on the two-level scoring proposal and the independent review (P-3) completes.** The current freeze has no application-scored evidence, and no instrumentation authorized to create such evidence exists within G1 scope.

---

## §9 Deliverable and Filing

### 9.1 Final Path

`docs/research/corpus/CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md`

### 9.2 Registry Entry

The report will be added to `docs/DOCUMENT-REGISTRY.md` with the next valid row after the current last entry (row 422). The row will have:
- **File:** `CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md`
- **Family:** Research
- **Role:** Adjudication
- **Current relevance:** CURRENT
- **SHA-256:** <to be computed after file save>
- **Bytes:** <to be computed after file save>
- **Lines:** <to be computed after file save>

### 9.3 File Verification

After saving, I will verify:
- The final file path is correct
- The registry entry is added with the correct row
- The file size and SHA-256 hash match
- The file is readable and complete

Let me now save the file and compute the verification metrics.
<tool_call>
<function=run_terminal_command>
<parameter=command>
wc -l docs/research/corpus/CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md && wc -c docs/research/corpus/CORPUS-91-INDEPENDENT-ADVERSARIAL-REVIEW-2026-11-13.md