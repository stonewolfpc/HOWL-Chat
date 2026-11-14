# C1 / P17 — Remaining Evidence-Gap Review (Read-Only)

**Date:** 2026-11-09

**Status:** Review phase opened. **P17 remains UNAPPROVED.** OD-9 is accepted as Option A; OD-1 remains DEFERRED; U1–U8 remain OPEN; contract rev3 remains frozen; C1 implementation remains unauthorized.

**Baseline:**

- Decision register: `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md`, SHA-256 `f9f0a1710776f2ca8870c87d9cfd4f6f6d4c61fc9b64f4509db9bf405e389960`, 409 lines.
- Dispositions: **8 ACCEPTED** (OD-2 through OD-9); **1 DEFERRED** (OD-1); **0 UNRESOLVED**.
- P17: `01d887c1…`, 166 lines, frozen author proposal, not approved.
- Contract rev3: `ed0da1a7…`, 183 lines, proposed revision, not accepted.
- U1–U8 remain OPEN; R11 is narrowed but not closed; R12 and R13 remain open.

This document begins the review of the remaining P17 evidence gaps. It does not revise P17, amend the contract, update the register, change source, run tests, or recommend implementation as the next step.

## 1. Purpose and boundary

The governance questions now have recorded dispositions. That does not make P17 technically sufficient. The remaining phase is to determine, gap by gap:

1. what evidence is missing;
2. whether the gap is an owner decision, a model-definition defect, a source-conformance defect, or missing runtime observation;
3. what evidence would close it; and
4. whether that evidence can be obtained without implementation.

The review is deliberately not a P17 revision and not a test plan. It identifies the minimum next review actions so that a later revision or test request can be scoped from evidence rather than assumption.

**Not authorized by this review:** P17 edits, contract amendments, source changes, tests, builds, benchmarks, runtime experiments, or C1 implementation.

## 2. Evidence classification

### 2.1 Observed

Directly read from the current governance record:

- OD-9 is accepted as Option A: P17 is not approved while OD-1 remains deferred and U1–U8 remain open without determined outcomes or explicit unknown-case rules.
- P17 is still an unapproved frozen proposal, not a contract edit or rev4.
- Contract rev3 is still proposed, not accepted.
- U1–U8 remain open.
- P17's §10 challenge names five projections, including `quarantined`, while §7's projection table has no `quarantined` row.
- The decision register records that the accepted decisions do not close R09–R13 or U1–U8.

### 2.2 Source-derived

The following are source-level facts established in the existing P17 assessment. They are not new runtime observations and were not re-tested here.

- There is no append-only evidence record, evidence sequence allocator, or stable instance identity stored with failures, actions, and exits.
- `Terminate` records no intent, start, completion, or result; its result is ignored at multiple call sites.
- `failAll` retains only the first read failure; `SetState` overwrites `lastFailure`; `waitForReader` discards join versus timeout.
- `Synthesize` publishes `crashed` on a send failure while the process may still be alive (`client.go:755–758`).
- A pending fatal protocol failure can publish `degraded` (`client.go:790–797`).
- `validateHello` returns `protocol_violation` for all nine validation failures, including a missing mandatory capability; `Start` maps only two helper-declared codes to `incompatible`.
- No terminal latch, terminal revision record, quarantine record, or ownership-uncertainty field exists.

### 2.3 Reported

- The owner accepted OD-2 through OD-9; those accepted directions are recorded in the register.
- The existing assessment reports R09–R13 and U1–U8 as open or narrowed.
- The register's dispositions are reported owner directions; this review does not independently verify the owner's intent beyond the recorded text.

### 2.4 Inferred

- Some gaps can be resolved by a finding-level review of existing source evidence; others cannot be resolved without an owner decision or runtime observation.
- A P17 revision that does not first resolve OD-1 would encode an evidence-record assumption the owner has not made.
- Approval remains blocked even if some U cases are reclassified, because OD-9 requires every case to have a determined outcome or an explicit unknown rule.

### 2.5 Unestablished

- Whether the evidence record is diagnostic-only or projection-supporting (OD-1).
- Which U1–U8 cases, if any, have outcomes that can be established from static source alone.
- Whether the current source's interleavings actually produce the paths the static review identifies.
- Whether P17's projection table can be reconciled with OD-4 through OD-8 without changing the contract.
- Whether any implementation conforms to the accepted directions.

## 3. Gap inventory

### 3.1 Decision and authority gaps

#### G-01 — OD-1 evidence-record authority

**Gap:** The owner has not decided whether the evidence record is diagnostic-only or projection-supporting.

**Why it blocks:** P17's precedence table claims to resolve projections from preserved evidence. Without an authority decision, the model cannot say whether the log may support a public state, whether it must be retained, or whether a missing entry is merely diagnostic. R09 cannot be assessed as a model requirement until this is settled.

**Closure:** A separate owner decision on OD-1. This review does not make that decision.

**Not sufficient:** Further source reading. The source has no log; reading it again cannot choose the log's authority.

#### G-02 — P17 has not been reconciled with accepted OD-4 through OD-9

**Gap:** P17 remains the frozen proposal. It still contains the pre-decision projection claims and has not been revised to reflect OD-4 (public projection), OD-5 (`crashed`), OD-6 (compatibility), OD-7 (internal quarantine), OD-8 (causality), or OD-9 (approval gate).

**Why it blocks:** Approving P17 as written would approve text that conflicts with accepted directions. For example, P17 §10 still claims `quarantined` as one of the five projections, while OD-7 says it is internal-only; P17's fatal-transport rule still maps to `crashed` while U6 shows the process may be alive.

**Closure:** A gap-closure matrix and, only after separate authorization, a P17 revision. This review does not edit P17.

#### G-03 — No explicit unknown-case rule

**Gap:** OD-9 permits an explicit `outcome_unknown` rule, but no document defines what that rule is, who records it, what evidence it requires, or how it maps to P17's projections.

**Why it blocks:** U1–U8 cannot be closed merely by saying "unknown." The rule has to distinguish an unresolved evidence question from a deliberately declared diagnostic outcome, and it must not silently become a public provider state.

**Closure:** A governance definition of the unknown-case rule, separate from P17 approval.

### 3.2 Model-definition gaps

#### G-04 — R09 record and ordering model

**Gap:** P17 names `id`, `seq`, `site`, `code`, `class`, `severity`, `detail`, `intent`, `grace`, `completion`, and an append-only log, but does not define the instance, the sequence allocator, the ordering domain, the action-to-exit linkage, phase, or retention.

**Evidence state:** Source-derived. The source has no corresponding fields, no stable instance identity, and no evidence sequence. The request counter is not an evidence sequence.

**Why it blocks:** The §7 precedence table consumes records whose semantics are undefined. An append-only log can preserve entries and still produce an ambiguous history.

**Closure:** OD-1 authority first; then a minimum record/allocator definition if the log is projection-supporting, or an explicit statement that it is diagnostic-only and cannot support §7.

**Source-only?** No. The authority and schema are design decisions.

#### G-05 — R10 causality and conflict resolution

**Gap:** P17's precedence order is a selection policy, not a demonstrated causal or ordering rule. It does not explain how two concurrent observations are ordered, how an exit is linked to an action, or when a combination should remain unresolved.

**Evidence state:** Source-derived. The source does not retain action outcomes, event order, or a causal link.

**Why it blocks:** U1–U5 depend on this rule. Assigning fatal-wins would encode a causation claim the source does not support.

**Closure:** Use OD-8's accepted evidence standard to define what would establish causation; otherwise leave the combination explicitly unresolved. Do not add a precedence level to manufacture certainty.

**Source-only?** No. The standard is a policy; the facts needed to apply it are absent.

#### G-06 — R11 `quarantined` claim

**Gap:** OD-7 makes `quarantined` internal-only, but P17 still lists it among the five required projections in §10 while §7 has no path for it.

**Evidence state:** Source-derived: no provider value, no process value, no supervisor quarantine field, and no public path.

**Why it blocks:** P17 is internally inconsistent as written. Approval of that text would approve an unreachable public projection.

**Closure:** Reconcile P17's text with OD-7: internal-only condition, no public projection, and no provider-state addition. This is a text/model correction, not an enum change.

**Source-only?** No; it requires a P17 correction decision, but the contradiction itself is already established.

#### G-07 — R12/U8 handshake taxonomy

**Gap:** A helper-declared incompatibility maps to `incompatible`, while a client-detected missing mandatory capability maps through `protocol_violation` to `unavailable`. P17 preserves both projections.

**Evidence state:** Source-derived and deterministic; no race is required to produce the split.

**Why it blocks:** OD-6 supplies a principle — only a structurally valid, semantically identified mandatory compatibility failure justifies `incompatible`; ambiguity or unestablished compatibility meaning remains `unavailable` — but the missing-capability case still needs a finding-level classification and a decision on whether the current mapping conforms.

**Closure:** Apply OD-6 to the nine `validateHello` exits and state which are compatibility failures, which are protocol violations, and which remain unavailable. Then reconcile U8.

**Source-only?** Partly: the split is established by source; the classification is a policy/design determination.

### 3.3 Source-conformance gaps

#### G-08 — R13/U6 `crashed` while the process may be alive

**Gap:** `Synthesize` publishes `crashed` on a transport send failure before establishing process exit; a later cleanup may or may not terminate the process. `Status` can report not alive while a live OS process remains.

**Evidence state:** Source-derived; the write is visible in the source. No runtime test is needed to establish the contradiction, though runtime could confirm the exact path.

**Why it blocks:** OD-5 defines `crashed` as an observed process exit where orderly shutdown is not established. A live process cannot satisfy that definition. P17 rule 3 nevertheless maps fatal transport to `crashed`.

**Closure:** A finding-level determination that the write is non-conforming, or a separately justified decision that `crashed` means a provider condition rather than process exit. The latter would require a contract-level decision; this review does not make one.

**Source-only?** The contradiction can be established from source; the disposition requires a governance decision.

#### G-09 — U3 handshake terminal versus later crash

**Gap:** A handshake failure may publish `incompatible`/`unavailable`; a later exit observer can overwrite it with `crashed` or `absent`. There is no terminal latch.

**Evidence state:** Source-derived race; no runtime interleaving was observed.

**Why it blocks:** OD-3 accepts an immutable first terminal decision. The source has no latch, and P17 has no explicit terminal publication rule. U3 cannot be assigned a unique outcome until the model says whether the handshake terminal survives.

**Closure:** Define terminal publication and immutability in the model; then obtain interleaving evidence for the actual path. This is model work first, runtime evidence second.

**Source-only?** No. The model rule and runtime behavior are both missing.

#### G-10 — U7 terminal followed by termination timeout

**Gap:** A terminal state can be published, then `Terminate` can set `stopping` and time out while ownership remains uncertain. No terminal latch or quarantine record exists.

**Evidence state:** Source-derived; no runtime timeout was observed.

**Why it blocks:** OD-3's immutability and OD-7's internal quarantine imply a terminal/ownership relationship P17 does not define. The current code can rewrite the terminal.

**Closure:** Define whether a terminal is immutable and how a timeout is represented internally; then test the path. No public state is added.

**Source-only?** No; the policy and eventual implementation are separate.

#### G-11 — U8 source mapping after OD-6

**Gap:** The source has two projections for one compatibility concept; OD-6 has a principle but no finding-level mapping for the missing-mandatory-capability path.

**Evidence state:** Source-derived; deterministic split.

**Why it blocks:** U8 is one of the eight cases OD-9 requires to be determined or explicitly unknown. The current P17 table reproduces the split rather than resolving it.

**Closure:** A finding-level compatibility disposition applying OD-6. If the missing capability is a mandatory compatibility failure, the current `protocol_violation`/`unavailable` path is non-conforming; if it is a protocol violation, the taxonomy must say why. This review does not choose.

**Source-only?** Partly; the mapping is visible, the classification is a decision.

### 3.4 Runtime and observational gaps

#### G-12 — U1 orderly EOF versus reader crash

**Gap:** The reader can classify a shutdown-related stream close as `helper_crashed`, while the exit observer may publish `absent` if it sees `stopping`; no stream-close provenance is retained.

**Evidence state:** Source-derived path; no observed interleaving.

**Why it blocks:** OD-4's orderly-shutdown rule depends on whether the EOF is an expected consequence of shutdown or an independent transport failure. The source does not record which.

**Closure:** Runtime evidence that distinguishes natural EOF, shutdown-induced EOF, forced close, and read error, plus an action/exit record. No implementation is authorized now.

**Source-only?** No. Static reading shows the race; it cannot establish the winner in a real execution.

#### G-13 — U2 shutdown timeout versus protocol violation

**Gap:** Two actions/observations can overlap: a protocol violation may call `Terminate`, while `Shutdown` also calls `Terminate`; one may time out while the other completes. No action result or ordering is retained.

**Evidence state:** Source-derived path; no observed interleaving.

**Why it blocks:** P17 has no arbitration rule and no public projection for the no-exit timeout case. The source's later `absent` or `stopping` is timing-dependent.

**Closure:** Define arbitration between a pre-existing fatal observation and a shutdown action; obtain action-result and exit evidence; decide whether an internal quarantine or explicit unknown applies.

**Source-only?** No.

#### G-14 — U4 zero exit after fatal protocol violation

**Gap:** A fatal protocol observation can be followed by a zero exit while `stopping` is set, producing `absent`; P17 rule 2 would produce `crashed + protocol_violation` regardless of exit status. The two rules conflict.

**Evidence state:** Source-derived conflict; no observed execution.

**Why it blocks:** P17's own rules cannot both hold. A zero exit after a fatal observation may mean the fatal observation was diagnostic, the action caused the exit, or the observation was wrong; the source does not establish which.

**Closure:** Resolve the S3/rule 2 contradiction; decide whether a zero exit after a fatal observation can be `crashed`; obtain runtime evidence of the actual path.

**Source-only?** No. The model conflict is visible, but the outcome is not.

#### G-15 — U5 nonzero exit after deliberate kill

**Gap:** A nonzero exit observed while `stopping` is published as `absent`; P17 proposes internal `terminated_by_action` but leaves the public projection undecided. No action intent, result, or causal link is retained.

**Evidence state:** Source-derived; no observed execution.

**Why it blocks:** OD-4 says confirmed exit without orderly shutdown → `unavailable` unless a specific justified prior failure applies; OD-8 requires A–E evidence for action causation. The current source does not establish either.

**Closure:** Define the public projection for action-caused nonzero exit; obtain action outcome and exit evidence; apply OD-8 rather than inferring causation from sequence.

**Source-only?** No.

#### G-16 — U6 live-process `crashed` (already G-08)

This is the same gap as G-08, recorded separately because U6 is an approval case. It should not be counted twice in a closure plan: resolving the `crashed` semantics and the source write addresses both R13 and U6.

#### G-17 — U7 terminal timeout (already G-10)

This is the same gap as G-10, recorded separately because U7 is an approval case. The model decision and the runtime path are both required.

## 4. R09–R13 and U1–U8 closure view

| ID | Current status | Gap class | Can existing evidence close it? | Minimum next review action |
|---|---|---|---|---|
| R09 | Open | Model/authority | No | OD-1 authority; then record/allocator definition or diagnostic-only declaration |
| R10 | Open | Model/policy | No | Apply OD-8; define causal/order rule or leave unresolved |
| R11 | Narrowed, not closed | Model/contract text | Partly: contradiction is established | Reconcile P17 with OD-7; no public state |
| R12 | Open | Policy + source mapping | Partly | Apply OD-6 to handshake failures |
| R13 | Open | Source conformance + policy | Partly: write is established | Decide whether the live-process `crashed` write is non-conforming |
| U1 | Open | Runtime + action provenance | No | Collect EOF/action/exit evidence; then apply OD-4 |
| U2 | Open | Runtime + arbitration | No | Define fatal-vs-action rule; collect timeout/exit evidence |
| U3 | Open | Model latch + runtime | No | Define terminal immutability; observe the interleaving |
| U4 | Open | Model conflict + runtime | No | Resolve S3/rule 2; observe zero-exit path |
| U5 | Open | Causality + runtime | No | Apply OD-8; collect action/exit evidence |
| U6 | Open | Source conformance | Yes for the contradiction; no for disposition | Decide whether the write conforms to OD-5 |
| U7 | Open | Model latch/ownership + runtime | No | Define terminal/quarantine relation; observe timeout |
| U8 | Open | Policy + source mapping | Partly | Apply OD-6 to the two handshake paths |
| OD-1 | Deferred | Authority | No | Separate owner decision |

## 5. What can be resolved without implementation

### 5.1 Can be resolved by a finding-level review

- **R11:** P17's `quarantined` claim contradicts OD-7 and its own §7 table. The textual correction can be specified without code.
- **R13/U6:** The source publishes `crashed` before establishing exit. A finding can state that this conflicts with OD-5; no test is needed to see the write.
- **R12/U8:** The two handshake paths are deterministic and visible. A finding can classify them under OD-6, though the classification is a policy choice.
- **U3/U7 model portion:** OD-3 already states terminal immutability. The model can state that a latch is required; the source's lack of a latch is a separate conformance gap.
- **P17 consistency:** A gap-closure matrix can identify every place P17 conflicts with accepted decisions without editing P17.

### 5.2 Cannot be resolved by source review alone

- **OD-1:** the authority choice is the owner's.
- **R09:** the record schema and allocator require a design decision; source has no implementation to inspect.
- **R10:** causal standard and conflict policy require a rule, not another source read.
- **U1, U2, U4, U5:** the actual interleaving and action/exit relationship require runtime evidence; static reading establishes possible paths, not observed outcomes.
- **U3, U7:** terminal and timeout behavior require both a model rule and runtime observation.
- **U6, U8:** the contradiction is static, but the disposition (conforming/non-conforming, incompatible/unavailable) is a decision.

## 6. Dependency order for the next review

This is a review order, not an implementation plan.

1. **OD-1:** decide whether the evidence record is diagnostic-only or projection-supporting. Nothing about R09, R10, or U1–U7 can be finalized without this.
2. **P17 reconciliation matrix:** map every P17 rule and projection to OD-4 through OD-9. Identify conflicts; do not edit P17.
3. **Static conformance findings:** resolve R11, R12/U8, and R13/U6 at the finding level, applying OD-6/OD-7/OD-5.
4. **Model rules:** specify terminal latch, action/exit linkage, unknown-case rule, and causal standard consistent with OD-3/OD-4/OD-8.
5. **Runtime evidence:** collect U1–U7 interleaving and action/exit evidence under a separately authorized test plan. No tests are authorized now.
6. **U-case disposition:** assign each case a determined outcome or an explicit unknown rule.
7. **Reconsider P17 approval:** only after the above, and only by a separate governance decision. Implementation remains separate.

## 7. Recommended next actions

1. **Record OD-1 as the next governance decision.** The question is not whether to implement a log; it is whether the evidence record may support a projection at all. A diagnostic-only decision would require P17 §7 to stop claiming projection resolution. A projection-supporting decision would require a record/allocator definition and retention rules.

2. **Create a P17-to-OD reconciliation matrix.** List every P17 rule, projection, and claim against OD-4 through OD-9. Required entries include:
   - `quarantined` vs OD-7;
   - fatal transport → `crashed` vs OD-5/U6;
   - handshake split vs OD-6/U8;
   - terminal immutability vs OD-3/U3/U7;
   - nonzero exit after action vs OD-4/OD-8/U5;
   - P17 §10's five-projection claim vs §7's actual table.

3. **Do not revise P17 yet.** A revision before OD-1 would guess at the evidence authority. The matrix is the prerequisite, not the revision.

4. **Separate static findings from runtime questions.** R11, R12/U8, and R13/U6 can be reviewed now. U1, U2, U4, U5, U3, and U7 require runtime evidence and a model rule; they should not be closed by source inspection alone.

5. **Accept that some cases may remain `outcome_unknown`.** OD-9 permits an explicit unknown rule; it does not require every case to be forced into a public state. The rule must be defined, not improvised per case.

## 8. Conclusion

The governance stage is complete, but the evidence stage is not. The remaining gaps are bounded and named:

- **OD-1** is an unresolved authority decision, not a source-inspection task.
- **R09/R10** are model-definition gaps; the source lacks the record and causal structure P17 assumes.
- **R11/R12/R13** contain static contradictions or unresolved classifications that can be narrowed by finding-level review, but not all can be closed without a policy choice.
- **U1–U8** remain open; none has a unique, source-justified outcome. Some require runtime evidence, and all require the prior model rules.
- **P17** remains stale relative to accepted decisions and cannot be approved as written.

The next phase should therefore begin with **OD-1 and a P17 reconciliation matrix**, not with implementation, a P17 rewrite, or a contract amendment. P17 remains unapproved; contract rev3 remains frozen; C1 remains unauthorized.

## 9. Verification and limitations

This is a read-only review. No existing file was modified. The register remains at `f9f0a171…`, 409 lines; P17 remains `01d887c1…`; contract rev3 remains `ed0da1a7…`; no source, test, configuration, or generated file was changed.

**Limitations:** this review relies on the existing P17 assessment and register for source-derived facts; it did not re-run the full source audit. It did not execute any U case, observe any interleaving, or verify platform-specific behavior. It does not establish that any implementation conforms to the accepted decisions. The source citations are those of the underlying assessment, not new runtime observations.

**No implementation was performed or authorized.**
