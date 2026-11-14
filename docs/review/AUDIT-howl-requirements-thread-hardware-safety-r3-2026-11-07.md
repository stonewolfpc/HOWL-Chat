# AUDIT — Requirements revision r3 (thread safety, hardware safety, performance preservation)

**Date:** 2026-11-07
**Mode:** READ-ONLY ADVERSARIAL REVIEW — no edits to any subject
**Subject:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r3-2026-11-07.md`
**Subject SHA-256:** `e45f1d898ab3d66016d2714472fd1023bcc93926a5e250a94bbb9072d9082d16` (513 lines / 47,710 bytes)
**Reviewer:** Mimo
**Status of subject:** PROPOSED CORRECTIVE REVISION — NOT ADOPTED — NOT BINDING

## 0. Scope, method, and what this audit does not do

This audit examines r3 against the authority it claims to serve: the original recorded requirements, the audited r2 and its findings A-1–A-11, the frozen TTS contract, the corrective design proposal, the correction record, and the owner-decision instrument. It asks two separate questions:

1. Does each prior finding now have a resolution that actually holds?
2. Do the resolutions introduce new contradictions with the inputs, with the r3 text itself, or with the authority structure?

Question 2 is the reason this review exists. A revision can close every previous finding and still be unadoptable because of what the closure cost elsewhere.

This audit does not adopt, approve, authorize, or reject anything. It creates no adoption record, modifies no subject, and authorizes no stage, implementation, measurement, or benchmark. Nothing here is a decision on any owner question.

Every hash, line count, and byte count below was recomputed from the working tree during this review.

### 0.1 Severity taxonomy

- **Blocking** — would freeze a wrong interpretation into a stage's acceptance tests, or would let a stage pass while a contract-settled behavior is unimplemented or suppressed.
- **Major** — an authority, scope, or allocation defect that changes what a stage must prove, without by itself producing a wrong pass.
- **Minor** — a citation, terminology, or drafting defect with a clear correct reading.
- **Observation** — a property worth recording that is not a defect.

### 0.2 Preservation checks performed

| Check | Result |
|---|---|
| Governing contract `DESIGN-howl-tts-helper-architecture-2026-11-06.md` | `ec3c789e…c25b50f`, unchanged from the r2 review |
| Corrective design `DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` | `802ae811…e355b151`, unchanged |
| Correction record `CORRECTIVE-howl-tts-stage3-corrective-design-2026-11-07.md` | `8fe9356c…f6355e076`, unchanged |
| Owner-decision instrument `OWNER-DECISIONS-howl-tts-stage3-2026-11-07.md` | `24d33445…df122ba`, unchanged; all eight decisions still `PENDING — OWNER DECISION REQUIRED` |
| Original requirements `REQUIREMENTS-howl-thread-and-hardware-safety-2026-11-07.md` | `ee38f031…439ac0e5`, unchanged |
| Audited r2 `REQUIREMENTS-howl-thread-and-hardware-safety-r2-2026-11-07.md` | `c9bddd10…deb8ab`, unchanged |
| Mimo r2 audit `AUDIT-howl-requirements-thread-hardware-safety-r2-2026-11-07.md` | `567bc86c…48994b`, 221 lines / 29,130 bytes, matches r3 §1.1 and §2.2 exactly |
| r3 §2.2 governing-input table | All seven rows match recomputed hashes, line counts, and byte counts; no mismatch |

No subject was edited. No adoption record was created. No performance measurement, benchmark, build, test, or execution was performed.

## 1. Verdict

**Not adoptable as written.** r3 is a substantial and mostly successful correction, and it resolves the r2 conflicts it set out to resolve. Two blocking findings and four major findings remain. Both blocking findings are authority defects, not safety contradictions, and each is fixable with a bounded clause rather than a redesign.

The blocking findings matter because r3's own §14 checklist requires that "the review found no unresolved authority, contract, or safety contradiction" before adoption. Under r3's own criteria, adoption is blocked until R3-1 and R3-2 are closed.

Specific disposition of the question the owner asked this review to target — the A-11 budget-approval separation — is R3-1. The separation is directionally right, and r3 implemented exactly what the r2 audit recommended. The recommendation itself was under-specified, and r3 did not add the missing bound, so the reviewer currently acquires stage-gate authority over numeric limits and pass/fail thresholds without any textual limit tying that authority to already-authorized requirements.

## 2. Blocking findings

### R3-1 — The budget approver holds stage-gate requirement authority with no bound to authorized requirements

**Where:** §7 AC-4 (line 390), §9.1 (line 430), §10 Approval field and prose (lines 461, 465), §12 (lines 477–486), §2.1 (line 53).

r3 implements the A-11 fix by naming "the independent reviewer" as the approver of the resource budget at the stage acceptance gate, with the owner approving only exceptions. The approval covers the two fields that are themselves requirements:

- **Limit** — "A numeric value, or an explicit contract-provisional value with its source";
- **Acceptance threshold** — "The pass/fail boundary for this stage."

Three gaps follow.

**(a) The role is defined nowhere outside r3.** "Independent reviewer" appears in the frozen contract, the corrective design, the correction record, and the owner-decision instrument zero times. The design does anticipate an *independent review* function — C8's prerequisites require C1–C7 to be "independently reviewed" — but it grants that function no approval power, no threshold authority, and no gate. r3 does not merely name an approver; it creates a role and hands it a precondition of stage acceptance. §2.4 further requires this same undefined function to review r3 itself. No governing artifact assigns either duty.

**(b) Approval is an acceptance precondition but is not separated from acceptance.** §7 states "A stage with no declared resource budget is not accepted." Because approval happens "at the stage acceptance gate … before acceptance," a reviewer who declines to approve blocks acceptance. r3 nowhere says that budget approval is a conformance check rather than an acceptance act, nor that withholding approval records an unresolved AC-4 item instead of rejecting the stage. An approval power that gates acceptance is acceptance authority in substance.

**(c) Nothing binds an approved limit to previously authorized requirements.** This is the owner's specific concern, and the text does not contain the guardrail. §10 constrains the *provenance* of a limit (numeric, or contract-provisional with a named source) and forbids an arbitrary normal-operation CPU/GPU ceiling. It does not say that:

- a limit must lie within contract-fixed values, adopted requirements, and recorded owner decisions;
- a limit may not be more permissive than a frozen value;
- a limit that depends on an unresolved owner decision must be recorded as pending rather than approved;
- approving a budget cannot be conditioned on changes to requirement scope, stage order, or architecture;
- the approver may not introduce a requirement through the threshold field.

The concrete failure: the contract fixes the process pool cap at two helper processes (`DESIGN-howl-tts-helper-architecture-2026-11-06.md` line 2366: "Frozen v1 hard cap of two helper processes"; design line 467 repeats it as frozen). Nothing in §10 prevents a stage from proposing, and the reviewer from approving, a process-count limit of three; §8 C6 and C7 then test "process-count assertions" and "no path bypasses the limit" against three. The stage passes while a contract-settled behavior — the two-process cap — is unimplemented. That is the blocking failure mode.

**Mitigation, stated fairly:** §2.1 says the contract remains controlling and a conflict "must be recorded and escalated." That clause catches a direct conflict if someone notices it. It does not sit at the approval act, does not bind the threshold the reviewer signs, and does not stop the reviewer from approving a limit the author proposed. The general precedence clause is not a substitute for a bounded approval authority.

**Suggested fix (textual, small):** add to §10 that budget approval is a conformance act and is limited to verifying that each limit and threshold is within (i) contract-fixed values, (ii) contract-provisional values with a named source, (iii) adopted requirements, and (iv) recorded owner decisions; that the approver may not approve a limit more permissive than (i) or (iii), may not resolve, pre-empt, or substitute for a pending owner decision, may not condition approval on a change to requirement scope, stage order, or architecture, and that withholding approval records an unresolved AC-4 item rather than rejecting a stage. Define the role in the same section: who appoints it, its independence from author and owner, and that r3 creates it rather than the contract. Add one line to §7 making approval ≠ acceptance.

### R3-2 — §9.2's conjunctive gate forbids the contract-settled status/catalog lazy start that §8 C1 requires as evidence

**Where:** §9.2 line 438; §8 C1 row line 406; design line 36 and 362; instrument D9 row; contract line 1773.

§9.2 reads: "**Renderer/catalog:** do not start a helper or load a model from an inspection path unless the contract **and owner decision** permit it."

The conjunction is unsatisfiable for the ordinary case. The contract settles the start: §17.7.2's `absent` row is "explicit TTS/catalog/status request and helper resolved → `starting`, Lazy start." The corrective design says so explicitly in its own words: "the separate question whether an explicit status request may start an `absent` helper is already settled by §17.7.2 and is not part of D9" (design line 36), and "status-from-`absent` lazy start is contract-owned and verified at C7" (design line 362). No owner decision governs it — D9 is only whether observation resets the 120 s idle deadline, and the design excludes the start question from D9. So "the contract and owner decision permit it" can never be true for a status or catalog start: the contract permits, no owner decision exists, and the conjunction fails.

Meanwhile §8's C1 row requires, as C1 acceptance evidence, a "trigger-to-start mapping for TTS/synthesis, catalog, and status" and "authorized-trigger start evidence." A C1 stage cannot satisfy both §8 C1 and §9.2. One of the two readings must be wrong at acceptance time, which is exactly how a wrong interpretation gets frozen.

The sentence is inherited (original requirements line 294, r2 line 379), so it is not a new sentence. It becomes a live contradiction only because r3 adds the explicit status/catalog start evidence in §8 C1 without reconciling §9.2.

**Suggested fix:** replace the conjunction with a disjunction that matches the authority: "unless the contract permits that start; an owner decision is required only where the contract is silent." Then state in §8 C1 that the permitted catalog trigger is the adapter's `ProviderCatalog` method (contract §17.12.1), not the renderer-facing `GetTTSVoiceCatalog` projection, which the design keeps static with no helper start while O-1/D5 remain open.

## 3. Major findings

### R3-3 — r3 calls an unadopted design proposal "frozen" and asks the owner to certify a status the design does not have

**Where:** §1 F10 row (line 28), §2.2 bullet (line 70), §6.4 (line 354), §6.5 (line 371), §8 header (line 402), plus the §2.2 table row at line 84 which labels the same file "Corrective proposal."

The design's own header says: "**Mode:** DESIGN ONLY", "**Proposal status:** Not implementation authorization", and it opens as "the Stage 3 corrective design and sequencing proposal." There is no design adoption record anywhere in `docs/review/`, and the instrument does not adopt it. It is a proposal.

r3 nonetheless refers to "the governing contract and frozen design inputs" (§2.2), "the frozen design's C8 exclusion remains controlling" (§6.4), "the frozen design's evidence matrix" (§6.4, §8), and "each stage's frozen functional acceptance criteria" (§8).

This matters for three reasons. First, §2.2 requires the adoption record to certify that "frozen design inputs … were not amended" — designating as frozen an artifact whose own status is "proposal," so the certification would assert an authority the design has never been granted. Second, "remains controlling" is authority language: an unadopted proposal does not control anything; what controls is the contract, and what limits C8 is C8's own proposed scope plus the absence of any authorization for performance claims. Third, the document is internally inconsistent about the same file's status in the space of twelve lines — the table says "Corrective proposal," the bullet above it says "frozen design inputs."

**Suggested fix:** use "recorded design proposal" or "the design proposal's C8 scope" throughout, and state in §6.4 why the exclusion binds: not because a frozen document says so, but because no performance claim has been authorized and r3 assigns PR evidence no acceptance effect. If the intent is that the design is fixed as a review reference, say "fixed for review purposes," not "frozen" or "controlling."

### R3-4 — Restoring the original requirements to operative force creates an AC identifier collision

**Where:** §2.1 (line 53–57), §7 (lines 387–396); original requirements lines 248–254.

r3's A-3 fix is correct in substance: the original recorded requirements remain in force until r3 adoption. But the original defines **AC-1 … AC-7**, where **AC-7 is "No stage substitution."** r3 defines **AC-1 … AC-10**, where **AC-7 is "Application I/O evidence"** and no-substitution is **AC-10**. AC-4 and AC-6 also differ in content: the original's AC-4 requires bounded limits and supporting measurements, r3's AC-4 adds a measurement method and an acceptance threshold; the original's AC-6 has no idle-grace exception, r3's AC-6 does.

While both documents are operative, the same identifier denotes two different rules. A stage report that says "AC-7 satisfied" is ambiguous between I/O evidence and no-stage-substitution, and §2.1's "where r3 is silent, the original remains in force" does not resolve it, because r3 is not silent on AC-7 — it simply says something else there. A stage could satisfy r3's AC-7 and miss the original's no-substitution rule, or vice versa, with both documents cited in support.

**Suggested fix:** either keep the original's numbering for AC-1–AC-7 and append the new rules as AC-8–AC-10 with a mapping note, or state explicitly that r3's AC set replaces the original's AC set *in whole* on adoption and that until adoption the original's numbering governs with r3's AC-8–AC-10 having no analogue. Do not leave two operative numbering schemes sharing seven identifiers.

### R3-5 — The C1 trigger list is state-blind and contradicts A2-D12 for the `unavailable` state

**Where:** §8 C1 row (line 406); contract A2-D12 (line 2513) and §17.7.2.

§8 C1 states a start "is permitted only for an explicit TTS/synthesis, catalog, or status request authorized by contract §17.7.2." That citation does scope the permission to the `absent` row, which is the correct location — but the C1 evidence list then asks for a "trigger-to-start mapping for TTS/synthesis, catalog, and status" and "authorized-trigger start evidence" without restating that scope.

A2-D12 says the opposite for the failed state: once `State=unavailable`, "`synthesize`, `catalog`, voice-list, `health`, and status requests do not trigger that recovery"; recovery is the explicit `load` only, and if no live helper remains, the update-gated `unavailable` → `absent` row governs first. So a stage could demonstrate "catalog request starts a helper" from `unavailable`, cite r3's C1 trigger map as authorization, and pass C1 while implementing a contract violation. The later C8 matrix row "unavailable no-recovery" would eventually catch it, so this is not a program-level escape — but it is a stage-level false pass, which is what a stage overlay is supposed to prevent.

**Suggested fix:** state the permission as "from `absent`" explicitly, and add to the C1 row: "A2-D12 governs `unavailable`: no status, catalog, or synthesis request starts or recovers a process from that state; only an explicit `load` does, and only on a live helper whose helper-side state is `unloaded`."

### R3-6 — C1's new start evidence reallocates verification the design assigns to C7

**Where:** §8 C1 evidence (line 406); design lines 362 and 473; design C1 tests (lines 133–183).

The design says plainly: "status-from-`absent` lazy start is contract-owned and verified at C7, not tested as a D9 choice" (line 362). Its C7 tests include "status/catalog no load" and the adapter-behavior split between `ProviderStatus`/`ProviderCatalog` (may start) and `GetTTSVoiceCatalog` (no start), with the C7 acceptance covering reachability (line 473–481). Its C1 test list is protocol/client capability classification and does not include a supervisor start-trigger test.

r3 makes "trigger-to-start mapping" and "authorized-trigger start evidence" part of **C1's** required evidence. C1's scope in the design does not own the supervisor start path; its only supervisor touch is storing service and provider capabilities separately. So either C1 must prove behavior its code does not own (and can be failed for evidence outside its stage), or C1 expands into C7's reachability work, violating the design's stage boundaries and the "do not combine" rule.

**Suggested fix:** keep the trigger map as a *design assertion* checked at C1 (documented mapping, no test), and place the executed start evidence at C7 where the design puts it. Alternatively, state explicitly that r3 amends the stage allocation and that the design's C1/C7 split is superseded — but then say so, because the design is an input the review must reconcile.

## 4. Minor findings

### R3-7 — The ProviderStatus/ProviderCatalog rows are cited to the wrong section

§8 C1 (line 406) cites "contract §17.7.2 (including the `ProviderStatus` and `ProviderCatalog` rows)." §17.7.2 is "Supervisor transitions" (contract line 1769); its lazy-start row is "explicit TTS/catalog/status request and helper resolved." The `ProviderStatus` and `ProviderCatalog` rows are in **§17.12.1 Adapter behavior** (contract line 2255; rows at 2266–2267). The substantive rule r3 states is correct in both places — but the trigger map it makes checkable is only partly checkable, because one half of the citation cannot be found where it is cited. Fix the pointer: "§17.7.2 lazy-start row and §17.12.1 `ProviderStatus`/`ProviderCatalog` rows."

### R3-8 — AC-7's N/A rule points at a section that defines no path categories

AC-7 (line 393) requires "explicit N/A for each in-scope path category it did not touch, per §9.1." §9.1 (lines 417–431) lists the fields of an I/O impact record and contains no path categories. The categories are defined in §3.2: settings, caches, model/asset preparation, status and health observation, renderer and picker paths, logging and diagnostics, helper lifecycle, temporary/output files. Point AC-7 at §3.2 (or move the list into §9.1), otherwise "each in-scope path category" has no enumerable referent and the A-10 fix is only nominally closed.

### R3-9 — Two C8 statements misdescribe what the design requires

The design contains **no occurrences** of "hardware" or "workload" (verified by search across all 711 lines). Yet §1's F10 row (line 28) says C8 retains the evidence matrix "including hardware/workload details where the design requires them" — a clause with no referent, made vacuous by its own qualifier. Separately, §8's C8 row (line 413) says "Resource, handle, and process observations are limited to the safety evidence the design already requires," but the design's C8 work list and acceptance rule require none of those observations; the leak checks appear in r3's own TH column, which is an additive overlay, not existing design evidence.

The substantive position — no performance claim in C8 — is correct and is preserved. The problem is that these two sentences are the reassurance that C8 gained nothing, and both are inaccurate about the design's content. Say plainly that the leak/handle/process checks are overlay additions of the same kind as the rest of TH-1, and delete the hardware/workload clause or replace it with the design's actual evidence rules (environment and process boundary per row).

### R3-10 — "Catalog" as a permitted trigger is ambiguous against the renderer projection

§8 C1's trigger list says "catalog" without distinguishing the adapter's `ProviderCatalog` (contract §17.12.1: "may start/handshake and call `catalog`") from the renderer-facing `GetTTSVoiceCatalog`, which the design explicitly keeps as a static fallback with "no helper start" while O-1 remains open (design lines 456, 473, 481). r3 names the `ProviderCatalog` row in the parenthetical, which is the right disambiguation, but the evidence line reverts to the bare word. Given that D5/O-1 are pending and §11 forbids renderer changes, spell it out once: the permitted trigger is `ProviderCatalog`, and `GetTTSVoiceCatalog` must not start a helper.

### R3-11 — "The same approval does not create its own input" is unparseable

§10 (line 465) reads: "Budget approval is recorded at the stage acceptance gate by the independent reviewer, after the budget is prepared but before acceptance; the same approval does not create its own input." The intent is presumably that the approval is not itself evidence of the budget's correctness and that the approver does not supply the measurements. As written it does not say that. Replace with: "Approval certifies that the recorded limits and thresholds meet §10's conformance conditions; it is not a measurement, does not substitute for the Measurement field, and does not authorize work beyond the approved limits."

### R3-12 — PR-1 conditions a measurement dimension on an authority that does not exist

PR-1 (line 360) says GPU/accelerator dimensions "apply only if an authorized native-provider scope includes them; otherwise record N/A with review." No native-provider scope exists or is authorized; the design's C8 exclusion is "No native provider," and native validation is reported as unestablished. The conditional therefore always takes the N/A branch, which is fine, but it names an authorization class the program has deliberately not created. Either say "native-provider scope is not authorized; record N/A" or drop the conditional.

## 5. Verified resolutions of A-1–A-11

Each prior finding was checked against the text, not against r3's summary table.

| Finding | Verified? | Evidence and residual |
|---|---|---|
| **A-1 — C8 performance conflict** | Yes, textually | §6.4 and §8's C8 row assign PR-1–PR-4 no C8 or stage effect; "no performance claim" and "PR results are not C8 evidence" are stated. §6.5 records the two options without selecting one; Option A is "recommended for investigation, not selected." The design's exclusion is no longer contradicted. Residual: R3-3's "frozen/controlling" wording. |
| **A-2 — all-eight-pending adoption condition** | Yes | §13 explicitly makes the pending status "a status report, not an adoption condition"; §14's checklist requires that adoption "does not require or imply that any decision remains pending." The r2 defect is gone. |
| **A-3 — predecessor status** | Yes in substance | §2.1 keeps the original in force until adoption and states r3 supersedes only where the texts differ. Residual: R3-4's numbering collision, which is created by this fix. |
| **A-4 — no governing-input hashes** | Yes, exactly | §2.2's table was recomputed row by row: all seven hashes, line counts, and byte counts match, including the r2 audit's 221 lines / 29,130 bytes. |
| **A-5 — handshake trigger unnamed** | Yes in substance | §8 C1 names the permitted triggers and excludes direct `unload`, direct `Cancel`, unknown-capability inspection, startup, and health, and excludes a start merely to complete a handshake. Residuals: R3-5 (state scoping), R3-6 (stage allocation), R3-7 (citation), R3-10 (catalog ambiguity). |
| **A-6 — inconsistent PR timing** | Yes | §6.4, §7 AC-8, and §8's C8 row now agree: stages record an impact statement only; no stage runs the suite; no PR result is C8 evidence. |
| **A-7 — scope narrowed** | Yes | §3.1 restores "relies on … or directly integrated with" and keeps the anti-refactor guard. |
| **A-8 — exception can fail PR-3** | Yes | §12 adds "No exception may fail HW-2" and treats HW-1.8's throttle prohibition as part of HW-2.1 for that purpose. |
| **A-9 — PR-3 has no N/A or signal** | Yes | §6.4's PR-3 row adds the reviewed N/A rule, names acceptable contention signals (queue depth, rejection count, measured latency), and requires a stated detection latency. |
| **A-10 — AC-7 N/A over-breadth** | Yes, nominally | AC-7 now scopes N/A to "each in-scope path category." Residual: R3-8, the categories are not in §9.1. |
| **A-11 — circular approval field** | Implemented as recommended; guardrail missing | The approver is named and the owner is limited to exceptions, as the r2 audit proposed. The recommendation did not bound the role, and r3 did not add the bound: R3-1. |

## 6. What r3 gets right, for the record

- The authority separation is real: r3 is explicitly non-binding, cannot accept or reject a stage, and names adoption as a separate owner act with its own record and hash.
- The original requirements are preserved as operative, which is the correct reading of their recorded status.
- No owner decision is resolved. r3 §13 matches the instrument exactly: D10(a), D10(b), D1a, D9, D1b, D2, D5, D3 all remain pending, and A-1 is correctly kept outside those eight.
- A-1 is preserved as an owner decision, with neither option selected and PR-1–PR-4 expressly inactive. "Recommended for investigation" is the correct strength of language.
- The C8 protection is textually sound: no performance claim, no PR evidence, no throughput or GPU claim.
- F3 (replacement sequencing), F4 (timeout/exit as first-class events), F5 (120 s grace), F6 (app-wide I/O record), F7 (measurable budgets), F8 (exception closure), and F9 (observability) remain correctly stated, and F2 is now concrete enough to test once R3-5/R3-6 are fixed.

## 7. Observations

- **O-1 — The A-11 residual is mine as much as the author's.** My r2 audit recommended naming the independent reviewer without defining the role or bounding its authority. r3 implemented that recommendation faithfully. The gap in R3-1 is a defect in the recommendation that the correction inherited.
- **O-2 — HW-1.8 has no measurement vehicle while A-1 is open.** AC-8 records an impact statement only, so "no arbitrary ceiling" is enforced through §10 limits and design review, not measurement. This is the correct consequence of leaving A-1 open, but it should be said out loud in the adoption record: adopting r3 as a non-performance revision defers HW-1.8's measurable half.
- **O-3 — The §14 checklist is now a real gate, not a formality.** Because it requires an independent review that finds no unresolved authority, contract, or safety contradiction, R3-1 and R3-2 block adoption by r3's own rule. That is the checklist working as intended.
- **O-4 — The r3 §1.1 disposition table is accurate about what was done**, including where a finding was closed only in prose. It does not overstate.
- **O-5 — Same-program review.** I authored the r2 audit and the three design audits r3's lineage responds to. This review is independent of the r3 author, but it is not independent of the program. A different reviewer would be a stronger check on R3-3 and R3-6 in particular.
- **O-6 — No execution was possible or permitted.** Every finding here is textual. Nothing about runtime behavior, measurement, or acceptance was tested, and this audit establishes none of it.

## 8. Limitations

- F1–F10 still exist only in r3's own §1 table. No recorded audit of the original requirements document exists anywhere in `docs/review/`, so the findings' provenance remains unverifiable and their resolutions were checked against r3's description of them.
- The corrective design is a proposal, so "the design says" is evidence of intent, not of authority. Where r3 and the design differ (R3-6), this audit states the difference rather than resolving it.
- This review read the design's C1 and C7 blocks, C8's scope and acceptance, the contract's §17.4.1, §17.7.2, §17.12.1, §17.8 cap line, and A2-D12, plus the instrument's decision table. It did not re-derive the full 2,586-line contract.
- All eight Stage 3 decisions remain separate and untouched by this requirements review.

## 9. What adoption would require

1. Close R3-1 by bounding the approver's authority and defining the role (§10, §7).
2. Close R3-2 by making §9.2's condition disjunctive and consistent with §8 C1.
3. Correct R3-3's authority language about the design.
4. Resolve R3-4's identifier collision.
5. Decide R3-5 and R3-6 deliberately: state the `unavailable` exclusion and either keep start evidence at C1 with an explicit allocation note or move it to C7.
6. Fix the citation and drafting defects R3-7 through R3-12, which are small and do not require a review round of their own.

None of these requires changing the safety requirements, resolving an owner decision, or amending the contract. All are r3-internal text changes.

---

**Audit conclusion:** r3 is the strongest revision of this document so far and is close to adoptable. It is not adoptable as written, because it grants an undefined reviewer stage-gate authority over requirements, and because its I/O rule still forbids a start its own C1 row demands as evidence. Both are closing-text defects, not design defects. No subject was modified; no adoption record exists; no stage, implementation, measurement, or benchmark is authorized.
