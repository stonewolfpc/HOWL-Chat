# AUDIT — Howl thread/hardware/performance requirements, revision r4

**Reviewer:** Mimo (independent of the r4 author)
**Date:** 2026-11-07
**Mode:** Read-only adversarial review. No subject file was edited, repaired, reconciled, staged, committed, or executed. No build, test, replay, or measurement was run.
**Subject:** `docs/review/REQUIREMENTS-howl-thread-and-hardware-safety-r4-2026-11-07.md`
**Subject SHA-256:** `52a550853291e6d4b0db8473e7f4c78185f1a7414d06ae207c5363b12e754429`
**Subject size:** 603 lines / 68,813 bytes

## Verdict

**Not adoptable as written.** One blocking defect, five major, five minor, seven observations.

This is the strongest revision in the sequence and the first whose central architecture is right rather than merely quieter. The reviewer role is now bounded by construction, the authority matrix resolves the identifier collision without touching the original criteria, the C1/C7 allocation matches the design's own statements, and the adoption instrument no longer requires the owner decisions to be answered first. The blocking item is a single clause, not a design failure. The majors are consequences of the corrections, not recurrences of the original errors.

The blocking clause is worth stating plainly: r4 says a helper start is permitted "only from `absent`" (r4 L469). The contract defines a start from `crashed` through `backoff` into `starting` on the next explicit request while the restart budget permits (contract L1796–1797), and r4's own C5 row demands restart/backoff evidence (r4 L473). A contract-conforming implementation therefore fails C1, and a C1-conforming implementation must suppress the contract's crash-restart path. That is the one place where r4 would freeze a wrong interpretation into stage acceptance.

## Severity scale

- **Blocking** — would freeze a wrong interpretation into stage acceptance, or let a stage pass while a contract-settled behavior is unimplemented or suppressed.
- **Major** — authority, scope, or allocation defect that does not by itself produce a wrong pass.
- **Minor** — citation, ambiguity, or drafting defect.
- **Observation** — not a defect; recorded because it affected the review.

Observations are numbered `OB-n` rather than `O-n` to avoid collision with the contract's owner decision O-1.

---

## Blocking

### R4-1 — The start map contradicts the contract's crash-restart path, and C5 depends on that path

**Target:** state transitions; C1/C7 separation; adoption without contradicting the contract.

**Evidence.** r4 §8 C1 states: "A start is permitted only from `absent`, for an explicit TTS/synthesis, `ProviderCatalog`, or `ProviderStatus` request under contract §17.7.2's lazy-start row" (r4 L469). The same row makes a `GetTTSVoiceCatalog` start a C1 failure. §8 C5 then requires "restart count and backoff are bounded" and "restart/backoff bound tests" (r4 L473).

The contract's transition table contains starts that are not from `absent`:

- `crashed` + next explicit request + budget permits → `backoff` (contract L1796)
- `backoff` + delay expires → `starting`, "New process, same key" (contract L1797)
- `stopping` + process exits 0 → `absent` (contract L1793)
- `incompatible`/`unavailable` + helper/assets/config updated → `absent` (contract L1799)

The contract's own note says `backoff` is a supervisor scheduling state, not a provider-loaded state (contract L1801). The design does not restrict starts to `absent`: it records that "explicit status from `absent` remains the §17.7.2 C7 lazy start" and that "no restart from health/status/catalog observation" occurs (design L348), that the status start is "already settled by §17.7.2" (design L36), and that "status-from-`absent` lazy start is verified at C7 as §17.7.2 behavior" (design L374). Its C5 acceptance requires "three restarts within 10 minutes; a fourth attempt exhausts the budget" (design L374), with "Three restarts then unavailable; 1/2/5 sequence with fake clock" as the test (design L362). Those restarts are starts from `crashed`, which r4's rule forbids. So the exclusivity is new in r4, not inherited, and it contradicts the design's own C5 acceptance as well as the contract.

**Why it matters.** The crash-restart budget is a contract-settled behavior: a second or third attempt after a crash is a start, and it does not originate in `absent`. Under r4 as written:

1. A C1 report that records the real trigger map — including `crashed` → `backoff` → `starting` — fails C1, because it documents a start that is not from `absent`.
2. A C1 report that passes can only do so by omitting the restart edge, which then cannot support C5's required restart/backoff evidence without contradicting C1.
3. An implementer optimizing for acceptance is pushed toward suppressing the contract's restart path, which is a contract violation.

This is a wrong-pass and wrong-fail mechanism at the same time, which is why it is blocking rather than editorial.

**Fix.** Replace the exclusive rule with the contract's actual set of start sources, in the same state-scoped form already used for `unavailable`:

- From `absent`: explicit TTS/synthesis, `ProviderCatalog`, or `ProviderStatus`, per L1773.
- From `crashed`: the next explicit synthesis or load request — not a health, status, or catalog observation (design L348) — through `backoff` into `starting`, only while the restart budget permits (L1796–1797); budget exhaustion yields `unavailable` (L1798).
- From `unavailable`: the A2-D12 in-place `load` rule as already written.
- From `incompatible`/`stopping`: no start; the update-gated or exit-gated return to `absent` governs first (L1793, L1799).

C1's failure example should name the two prohibited cases it already intends — a start from `unavailable` and `GetTTSVoiceCatalog` as a trigger — rather than asserting an exclusivity the contract does not have.

---

## Major

### R4-2 — A-11 is structurally closed, but withholding remains an unbounded veto and contract interpretation is forbidden yet unavoidable

**Target:** A-11 closure; reviewer authority boundary.

**What is fixed.** §10 now denies the reviewer requirement, scope, stage-order, contract-value, architecture, and acceptance authority, and states that the owner remains the only change authority (r4 L531–550). The two-process example is correctly framed as a contract conflict rather than a budget choice (r4 L529). §7.1 repeats the boundary for every row. This is a real fix and it is implemented as recommended.

**Residual, first limb — withholding.** §10 lists "withhold conformance approval" as a permitted action with no stated precondition beyond "records that unresolved item" (r4 L537). There is no duty to certify a row that conforms, no requirement that a withholding name the specific non-conforming element and its evidence, no time bound, and no appeal or owner override short of changing the limit. Because R4-AC-4 is part of every stage's cumulative acceptance (§7 L410) and "a stage with no declared resource budget is not accepted" (§7 L425), an indefinite withholding blocks acceptance. §10 then says the reviewer "may not ... accept or reject a stage, or substitute for the stage's acceptance authority" (r4 L546). The reviewer cannot reject, but it can prevent acceptance indefinitely — a functional veto over the one act the text says it lacks. It can also create pressure to lower a conforming limit, which is a change the reviewer is forbidden to require.

**Residual, second limb — interpretation.** The reviewer may not "interpret a requirement, scope, stage order, contract value, or architecture" (r4 L542), yet must "verify that each limit ... is within contract-fixed values" (r4 L535). Determining whether a proposed limit is within a contract value *is* an interpretation of that value. The only escalation triggers are missing, conflicting, or decision-dependent limits (r4 L529), and §12's conflict rule covers only a requirement that cannot be met without violating the contract (r4 L564). A genuinely ambiguous contract clause is neither. So the reviewer must either interpret (forbidden) or stall with no recorded route.

**Why it matters.** These are not owner-authority leaks in the sense A-11 feared: the reviewer cannot change a requirement, widen a limit, or resolve a decision. They are accountability gaps in the gate the reviewer holds. The document now says who may not decide; it does not say how a conforming row gets certified, how a withholding is challenged, or who resolves a contract-reading dispute.

**Fix.** Three sentences: (a) the reviewer must certify a conforming row and may withhold only by naming the specific element that is missing, conflicting, or decision-dependent; (b) a withholding is an unresolved R4-AC-4 item that the owner may disposition by recorded decision, including a decision that the row conforms; (c) a dispute over what a contract value means is recorded as a contract-interpretation question and escalated to the owner, not resolved by the reviewer.

### R4-3 — No instrument assigns stage acceptance authority, so the author is the only actor not excluded from accepting

**Target:** matrix completeness; authority allocation.

**Evidence.** §7 states acceptance outcomes passively: "A stage that has functional tests but no applicable TH/HW evidence is not accepted. A stage with no declared resource budget is not accepted" (r4 L425). §10 bars the reviewer from accepting or rejecting and refers to "the stage's acceptance authority" as an existing thing (r4 L546). §7.1's matrix has columns for originating authority, evidence owner, conformance assessor, and change authority (r4 L433–459) — but no acceptance column and no acceptance actor. §2.4 refers to "stage acceptance" as an act (r4 L135). §13 leaves C1–C8 unauthorized, which is authorization, not acceptance (r4 L581).

The design imposes an independent review per stage and gives each stage an "Acceptance" section (design L179–181, L481), but names no acceptor either. The correction record does not either.

**Why it matters.** §10 excludes the reviewer and §12 excludes self-approval of exceptions and budget review. The change author is not excluded from stage acceptance anywhere. Because the same document requires the author to produce the evidence, the only actor left eligible to declare the stage accepted is the party whose work is being accepted. The alternative reading — that no one may accept — is also unsupported and would make C1–C8 uncompletable. Either way the gate is unassigned.

**Fix.** Add an "Acceptance authority" statement — most naturally in §7 and as a matrix column — naming the owner as the acceptor acting on the independent review, and stating that the reviewer's conformance finding and the owner's acceptance are separate acts. If acceptance is meant to be the reviewer's finding plus the owner's authorization, say so; the current text supports neither reading.

### R4-4 — The design proposal is simultaneously non-controlling and load-bearing for acceptance and change authority

**Target:** matrix completeness; C1/C7 separation; authority.

**Evidence.** §2.1 says the design and the correction record "are review inputs, not controlling authority; r4 does not adopt either one or make either binding" (r4 L79). §6.4 repeats that the design's C8 scope "is a review input, not controlling authority" and that "the absence of that authority, not the proposal's status, controls" (r4 L379).

§8 then says the overlay "supplements—and does not replace—each stage's recorded design acceptance criteria" (r4 L465), which makes those criteria part of what a stage must satisfy. §7.1 gives every stage row a change authority of "Owner; a change to stage allocation requires owner approval and a design-proposal amendment" (r4 L450–459). So the proposal's criteria gate acceptance, and its stage allocation is protected by a change gate that requires formally amending it.

**Why it matters.** If the design has no authority, it cannot supply acceptance criteria or require amendment. If it does supply them, §2.1's denial is false and the adoption record's "not amended" confirmation (r4 L90) becomes a substantive control on a document that was never adopted. r4 needs one coherent role for the proposal — a plan of record the owner executes, subordinate to the requirements but not itself a source of requirements — and it currently asserts both roles in different sections.

This is the residue of R3-3. The word "frozen" is gone and the C8 rationale is now correct; what remains is an allocation question the wording change did not settle.

**Fix.** State the role once, in §2.1: the design is the owner's recorded execution plan for C1–C8, not a source of requirements; r4's acceptance overlay applies *in addition to* the plan's own acceptance criteria because the owner has adopted the plan as the stage plan, and any change to stage allocation is an owner decision that must be reflected in the plan. Then §7.1's amendment condition follows from the plan's status instead of contradicting §2.1.

### R4-5 — Restating the original criteria as R4-AC rows can narrow them, and §2.1's supersession rule does not say which text governs

**Target:** matrix completeness; adoption without contradicting the original.

**Evidence.** Six of the ten overlay criteria restate originals (R4-AC-1, -2, -3, -5, -6, -10; r4 L413–421), and the matrix labels them "restates original AC-n" (r4 L440–449). §7 says the originals "remain unchanged and independently operative" and that "Where an overlay restates an original rule for readability, the original remains separately enforceable" (r4 L406–410). §2.1 says adoption "supersedes the original only where the two texts differ; where r4 is silent, the original remains in force" (r4 L77).

The restatements are not textually identical, and the differences are material:

- Original AC-3 requires that "Evidence must be repeatable, not a single happy-path run" (original L250). R4-AC-3 omits repeatability entirely (r4 L416).
- Original AC-2 requires "Static review is required where tests cannot prove the property" (original L249). R4-AC-2 omits it (r4 L415); §4.6 restores the idea for lock-order claims only.
- R4-AC-6 adds "except the contract-defined idle grace" (r4 L419) where original AC-6 has no exception (original L253). The exception is defensible as a clarification — r4 §4.3 and §5.3 give the idle timer an owner, bound, and trigger — but as written it reads as an exemption from AC-6's three-part requirement.

**Why it matters.** Clause-by-clause application of §2.1 makes r4's text supersede the original wherever it speaks, so AC-3's repeatability and AC-2's static-review requirement would be replaced by versions that omit them. Document-level application keeps both texts and preserves them. A stage report assessed under the first reading can pass AC-3 with a single run — a wrong pass on a settled requirement. §7's "for readability" language signals that no normative change was intended, which makes the omission a drafting error with acceptance consequences rather than a deliberate narrowing.

**Fix.** Make the restatements verbatim or mark them non-normative ("R4-AC-1 is a pointer to AC-1; it adds nothing"), and add to §2.1 that a restated overlay criterion does not supersede its original and that the original text governs on any difference. This also removes the double-counting described in OB-5.

### R4-6 — If A-1 Option B is selected, the adoption record cannot be completed as specified

**Target:** adoption mechanics; A-1.

**Evidence.** §2.2 requires the adoption record to contain "confirmation that the governing contract, the recorded design proposal, and the correction record match the table and were not amended" (r4 L94) — an unqualified non-amendment confirmation for the design. §6.5 Option B requires the owner to authorize "an amendment to C8's scope through the appropriate recorded instrument and records that amendment in the adoption record," and adds that Option B "cannot be paired with an unqualified statement that the recorded design proposal was not amended" (r4 L400). §2.3's adoption block has a "**Contract-amendment confirmation:**" field only (r4 L125); there is no field for a design-proposal amendment.

**Why it matters.** A-1 is open and either option is permitted. Selecting B requires the adoption record both to confirm the design was not amended and to record its amendment. §6.5 acknowledges the conflict in the abstract but §2.2's checklist item is unconditional, and the adoption template cannot express the qualified statement §6.5 demands. Adoption of r4 under Option B is therefore impossible without editing the instrument during adoption — exactly the kind of ad hoc reconciliation the sequence has been avoiding.

**Fix.** Make §2.2's confirmation conditional: "…and were not amended, except an amendment recorded under §6.5 Option B, which must be identified by scope and instrument." Add a design-amendment field to §2.3.

---

## Minor

### R4-7 — C1's required evidence says "likely N/A"

C1's I/O evidence cell reads "I/O impact record (likely N/A)" (r4 L469). A stage cannot record "likely": §9.1 requires an explicit N/A with a reviewer (r4 L482), and R4-AC-7 requires explicit per-category N/A against §3.2 (r4 L420). Replace with "I/O impact record, with explicit per-category N/A under §3.2."

### R4-8 — "C1 does not execute supervisor start tests" conflicts with the design acceptance r4 says it preserves

C1's evidence cell states "C1 does not execute supervisor start tests; C7 owns them" (r4 L469). The design's C1 acceptance requires real-process fake-helper handshake tests and states they "are sufficient for C1" (design L179, L181). §8 says the overlay does not replace the design's acceptance criteria (r4 L465), so C1 is required to run those tests and told it does not execute start tests. The two can coexist only if "supervisor start test" means something narrower than "a test that starts the fake helper," which r4 never defines. Say "C1 owns no *trigger-map* start evidence; C7 owns executed supervisor start and reachability," and keep the design's handshake tests explicitly in scope.

### R4-9 — C3's "no spawn from status" reads as forbidding the lazy status start that C1 and C7 permit

C3's focus reads "No eager spawn or model load from inspection/status/startup" (r4 L471). An explicit `ProviderStatus` request from `absent` is a permitted start (r4 L469, L475; contract L1773, L2266). "Eager" qualifies "spawn," but the phrase "from inspection/status/startup" without a qualifier is the same construction r4 uses elsewhere for prohibitions. Add "unsolicited" or a cross-reference: "no unsolicited spawn from inspection/status/startup; the explicit status start is C7's executed-evidence row."

### R4-10 — The decision ledger omits O-1, which C1 and C7 depend on

§13 lists the eight Stage 3 decisions plus A-1 (r4 L577–583), and §3.3 lists the same eight plus A-1 (r4 L167–172). The contract defers "whether the Character picker may start a helper" to owner decision O-1 (contract L2270), and its freeze table records "UI startup remains O-1" (contract L2365). r4 fixes the adapter-side answer — `GetTTSVoiceCatalog` must not start, `ProviderCatalog` may (r4 L469, L475, L501) — which is compatible with either O-1 answer provided the picker uses `ProviderCatalog`, but §13's rule that an answered decision "does not incorporate" would force r4 to ignore an O-1 answer it structurally depends on. Add O-1 to §13 as a contract-deferred owner decision, with the statement that an O-1 answer authorizing picker startup is satisfied through `ProviderCatalog` and does not make `GetTTSVoiceCatalog` a trigger.

### R4-11 — "More permissive than the authorized source" is undefined for a provisional value

§10 requires that a limit "may not be more permissive than the authorized source" (r4 L520) and that a limit may use a contract-provisional value with its source named (r4 L520; §5.4 L299). For a provisional source there is no fixed number, so a longer timeout or a larger buffer cannot be compared to it. The contract labels timeout values and the 120 s idle grace provisional (contract L2366 ff.). State the rule for provisional sources — either the provisional value is the ceiling until measured evidence is recorded, or the source must be named and the deviation justified in Basis — rather than leaving the comparison undefined.

---

## Observations

**OB-1 — All nine governing-input rows verified exact.** I recomputed every hash, line count, and byte count in §2.2: original `ee38f031…` 344/25,933; r2 `c9bddd10…` 447/36,025; r2 audit `567bc86c…` 221/29,130; r3 `e45f1d89…` 513/47,710; r3 audit `c9804763…` 226/30,143; contract `ec3c789e…` 2,586/182,904; design `802ae811…` 711/69,168; correction record `8fe9356c…` 172/16,190; instrument `24d33445…` 610/45,891. Every row matches. r4 omits its own hash, which is correct — it cannot self-hash.

**OB-2 — The two-process negative test is correct and well chosen.** The contract freezes "a hard cap of two helper processes" (contract L2366 ff.), and §10's refusal of a proposed three is the right illustration of the difference between conformance review and discretion (r4 L529). This is the clearest single demonstration in the document that A-11 was understood.

**OB-3 — The disjunctive §9.2 rule and the `ProviderCatalog`/`GetTTSVoiceCatalog` distinction are correct.** Original §9.2 was conjunctive (original L294); r4 makes the owner decision necessary only where the contract is silent, and distinguishes the adapter method that may start from the no-context method that may not (r4 L501), matching contract L2266–2267 and L2270. R3-2 and R3-10 are genuinely closed.

**OB-4 — The C1/C7 allocation matches the design.** The design states that the status-from-`absent` lazy start is settled by §17.7.2 (design L36) and "verified at C7" (design L374). r4 assigns the documented trigger map to C1 and executed start/reachability to C7 (r4 L469, L475). R3-6 is closed as intended, and the one remaining boundary problem is the exclusivity error in R4-1, not the allocation.

**OB-5 — Six of ten overlay criteria are restatements.** Beyond the narrowing risk in R4-5, the duplication lets a stage report ten satisfied overlay criteria from four unique requirements. Not a defect in itself; it inflates the apparent coverage of the overlay.

**OB-6 — r4's self-description of its evidence is honest.** §1.2 states that the r3 audit's reported hashes were rechecked but its repository-inspection claims were not independently verified (r4 L54). That is accurate. For this review I re-verified the specific contract and design passages r4 relies on, so its material claims are independently grounded even though its inherited provenance is not.

**OB-7 — Adoption is possible with all nine owner questions open.** I checked each dependency. §2.4 permits adoption with A-1 open if PR-1–PR-4 are inactive (r4 L133). §13 states that adoption neither requires nor resolves any decision (r4 L577). C7's D1a/D1b-dependent evidence is a stage gate, not an adoption gate, and §8 preserves the pending decision gates (r4 L465). Nothing in the adoption checklist requires an answered decision except A-1 for PR effect. Subject to R4-6, the adoption mechanism works.

---

## Claimed closure versus actual closure

### A-1–A-11 (r2 findings)

| Finding | r4 claim | Review result |
|---|---|---|
| A-1 | Open owner scope decision; PR inactive | **Closed as handling** — correctly left open, no acceptance effect, both options recorded without selection |
| A-2 | Adoption no longer requires the eight pending | **Closed** — §13 and §14 verified |
| A-3 | Original in force until adoption | **Closed** — §2.1; original hash verified |
| A-4 | §2.2 hash table | **Closed** — all nine rows exact (OB-1) |
| A-5 | State-scoped trigger map | **Partially closed** — the map is correct for `absent` and `unavailable`, wrong for `crashed`/`backoff` (R4-1) |
| A-6 | PR timing | **Closed** — §6.4 assigns no stage or C8 effect |
| A-7 | Scope restored | **Closed** — §3.1 |
| A-8 | Exception cannot fail PR-3 | **Closed** — §12 |
| A-9 | PR-3 N/A and signal | **Closed** — §6.4 table |
| A-10 | R4-AC-7 N/A by §3.2 categories | **Closed** — r4 L420 |
| A-11 | Bounded conformance review | **Partially closed** — the role is bounded against change authority, but withholding and interpretation are unbounded (R4-2) |

### R3-1–R3-12 (r3 findings)

| Finding | Review result |
|---|---|
| R3-1 unbounded approver | **Closed in substance** — role limited to conformance, owner retains change authority; residual veto in R4-2 |
| R3-2 conjunctive status gate | **Closed** — §9.2 disjunctive, verified against original L294 |
| R3-3 "frozen/controlling" | **Closed as wording; replaced by an unresolved role question** — R4-4 |
| R3-4 AC collision | **Closed** — bare `AC-n` / `R4-AC-n` separation verified; original titles match exactly (original L248–254) |
| R3-5 state-blind triggers | **Closed as scoping, but its replacement introduced the blocking defect** — R4-1 |
| R3-6 C1/C7 reallocation | **Closed** — matches design L36, L374 (OB-4) |
| R3-7 citation | **Closed** — both contract rows cited at r4 L469 |
| R3-8 AC-7 path categories | **Closed** — §3.2 |
| R3-9 C8 misdescription | **Closed** — F10 now distinguishes design content from r4's additive checks |
| R3-10 "catalog" ambiguity | **Closed** — every start rule distinguishes the two methods |
| R3-11 unparseable approval sentence | **Closed** — §10 conformance paragraph |
| R3-12 conditional native-provider authority | **Closed** — PR-1 N/A with review, no authorization class |

Net: 10 of 12 r3 findings fully closed, 2 partially. The claim that all twelve are closed is true at the level of "each finding has a corresponding correction," not at the level of behavior.

---

## Answers to the five review targets

**1. Can the reviewer still exercise owner authority by withholding?** Not owner authority in the A-11 sense — no requirement, scope, contract value, or architecture can be changed, and the two-process example is correctly treated as a conflict rather than a decision. But withholding is unbounded: no duty to certify, no evidentiary precondition beyond recording an unresolved item, no appeal, no time bound. It blocks acceptance without rejecting the stage, which is the distinction §10 relies on and does not fully earn (R4-2).

**2. Matrix completeness.** All 17 criterion rows are present with the right originating authority, evidence owner, and change authority; the original titles are exact; bare and prefixed identifiers are cleanly separated. Two gaps: no acceptance authority anywhere (R4-3), and a change-authority condition that depends on amending a document §2.1 calls non-controlling (R4-4). The six restatement rows also carry the narrowing risk in R4-5.

**3. C1/C7 separation.** The allocation is now correct and matches the design: C1 owns the documented map and structural checks, C7 owns executed start and reachability. Two residual problems, both textual: C1's exclusive "only from `absent`" rule reaches into C5's restart behavior and contradicts the contract (R4-1), and C1's "does not execute supervisor start tests" sits awkwardly against the design's C1 runtime evidence (R4-8).

**4. State transitions.** `absent`, `unloaded`, and `unavailable` are handled consistently and correctly — the `unavailable` no-recovery list matches A2-D12 exactly (contract L2513), and the in-place `load` condition matches. `crashed`, `backoff`, `stopping`, and `incompatible` are missing from the start map, and `crashed`/`backoff` is the blocking error.

**5. Adoption mechanics.** Adoption is possible with all eight Stage 3 decisions and A-1 still open, provided PR-1–PR-4 are recorded inactive; I verified each dependency. Two obstacles: Option B cannot be recorded under §2.2's unqualified non-amendment confirmation and §2.3's template (R4-6), and the supersession rule does not say whether a restated overlay criterion replaces its original (R4-5).

---

## What adoption requires

Textual changes only; no redesign:

1. Replace C1's exclusive start rule with the contract's actual start sources, including the `crashed` → `backoff` → `starting` restart (R4-1).
2. Bound withholding and add an ambiguity-escalation route (R4-2).
3. Name the stage acceptance authority (R4-3).
4. State the design proposal's role once, coherently (R4-4).
5. Make the restatements non-normative or verbatim, and state that they do not supersede their originals (R4-5).
6. Make §2.2's non-amendment confirmation conditional on Option B and add a design-amendment field to §2.3 (R4-6).
7. The five minor corrections (R4-7–R4-11).

Until then, r4 is **PROPOSED — NOT BINDING**, the original requirements remain in force, r2 and r3 remain non-adopted, A-1 remains open with PR-1–PR-4 inactive, and the eight Stage 3 decisions remain untouched.

---

## Limitations

- **L1 — Provenance.** F1–F10 still exist only as r4's own summary rows. I verified that each has a correction and that the corrections are internally consistent; I could not verify that the findings themselves were recorded as stated, because the source document was not provided to this review.
- **L2 — Design status.** The corrective sequencing document is a proposal, not an adopted artifact. Where r4 and the design differ I record the difference; I do not treat the design as authority, and I do not resolve the role question in R4-4.
- **L3 — Reviewer position.** I authored the r2 and r3 audits this revision responds to. This review is independent of the author but not of the review program.
- **L4 — No execution.** No build, test, replay, measurement, or benchmark was run, and no stage is authorized. All findings are document-internal or document-versus-contract.
- **L5 — Hash scope.** I verified the r4 hash and all nine §2.2 inputs. I did not verify the r4 audit's claims about files outside `docs/review/`.

## Review record

- Subjects read in full: r4 (603 lines), original AC table and §9.2, contract §17.7.2 transition table, §17.12.1 adapter table, §17.12.1 `GetTTSVoiceCatalog` note, process-cap freeze row, A2-D12; design C1 scope/prerequisites/runtime evidence, C1/C7 acceptance, process-cap block, C8 prerequisites; owner-decision instrument D1a, D1b, D3 rows and pending count (19 occurrences).
- Subjects hashed and compared: all nine §2.2 inputs plus r4.
- Files changed by this review: this file only. No subject was modified.
