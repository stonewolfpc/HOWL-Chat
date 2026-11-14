# C1 / P17 — OD-9 Final Approval-Gate Assessment (Read-Only)

**Date:** 2026-11-09

**Status:** Assessment only. **OD-9 remains UNRESOLVED.** No register update, P17 approval, contract amendment, source change, or implementation is made or authorized by this document.

**Baseline:**

- Decision register: `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md`, SHA-256 `6854363b…`, 396 lines.
- Seven decisions accepted: OD-2, OD-3, OD-4, OD-5, OD-6, OD-7, OD-8. OD-1 is deferred; OD-9 is unresolved.
- P17: `01d887c1…`, 166 lines, status **AUTHOR PROPOSAL FOR INDEPENDENT REVIEW. Frozen on write. Not a contract edit, not rev4.** P17 is not approved; C1 remains unauthorized (`P17 §11`).
- Contract rev3: `ed0da1a7…`, 183 lines, status **PROPOSED REVISION FOR INDEPENDENT REVIEW. Not an accepted amendment.**
- Lifecycle decision: `aa136bc9…`, 74 lines, an owner decision but not a contract amendment.
- U1–U8 remain OPEN; R11 is narrowed but not closed; R12 and R13 retain their open dispositions in the register.

## 1. Exact OD-9 wording

The register states:

> **OD-9 — P17 approval gate**
>
> **Question:** May P17 be approved while U1–U8 remain open, or must every case have a determined outcome or an explicit `outcome_unknown` rule before approval?
>
> **Status:** UNRESOLVED.
>
> **Alternatives:**
>
> - **A.** No approval until OD-1–OD-8 are dispositioned and every U case has a determined outcome or an explicit unknown rule.
> - **B.** Approve P17 as a partial model with open U cases.
> - **C.** Split P17: approve a minimal, fully specified subset and defer the rest.
> - **D.** Defer P17 entirely until rev4/P18 work is authorized.
>
> **Proposed default (non-binding):** **A.**

The exact wording is taken from the current register's OD-9 entry (§3). It is quoted for decision precision; it is not a new disposition.

## 2. What authority OD-9 asks the owner to exercise

OD-9 asks the owner to set the **evidence threshold for approving P17 as a model**, not to amend the contract, approve an implementation, or authorize C1.

The authority is narrow:

1. **Model approval:** decide whether P17 is sufficiently specified and dependable to be called approved, given the open U cases.
2. **Gate policy:** choose whether approval is blocked until every U case has a determined outcome or an explicit unknown rule, or whether some partial approval is acceptable.
3. **Scope decision:** if C is chosen, identify a minimal, fully specified subset; no such subset is currently defined.

It does **not** include:

- changing the thirteen-value provider-state contract;
- accepting contract rev3;
- amending P17;
- authorizing C1 implementation, tests, builds, or runtime work;
- closing U1–U8 by approval alone;
- making an owner decision equivalent to independent verification.

The register is explicit that owner acceptance records direction only and does not close R09–R13 or U1–U8 (§5). Approval of P17 would not be implementation evidence.

## 3. Evidence classification

The evidence supporting OD-9 is governance evidence, not runtime evidence. The categories below are kept separate.

### 3.1 Observed

Directly read from the current register and governance documents:

- The register records **7 accepted**, **1 deferred**, **1 unresolved**: OD-1 deferred, OD-9 unresolved, OD-2–OD-8 accepted.
- OD-9 is unresolved and its non-binding default is A.
- P17 is not approved; its status is a frozen author proposal for independent review, not a contract edit or rev4.
- Contract rev3 is a proposed revision for independent review, not an accepted amendment.
- P17 lists U1–U8 as unresolved and asserts no classification for them.
- P17 §10 challenge 5 names five required projections: `absent`, `crashed`, `incompatible`, `unknown`, and `quarantined`; P17 §7's projection table has no `quarantined` row.
- P17 §11 says C1 remains unauthorized and no tests, builds, or implementation were performed.
- The lifecycle decision says `unknown` and `quarantined` are not provider states and that public exposure requires a separate approved contract amendment.

### 3.2 Source-derived

The following source facts are taken from the existing assessments and register citations; this OD-9 assessment did not re-run or re-derive the full source review:

- `observeExit` maps an exit observed while `providerState == stopping` to provider `absent`, including a nonzero exit (`supervisor.go:176–182`).
- `Terminate` sets `stopping` and records no intent, action completion, or causal link (`supervisor.go:246–288`).
- The source has no authoritative evidence log, no all-writer sequence, and no stable association between action and exit (`OD-1` source entry).
- `crashed` has conflicting source writers: the observer path and the synthesis write-failure path (`OD-5` source entry).
- `quarantined` has no source value or path; the supervisor has no quarantine record (`OD-7` source entry; prior assessment R11).
- The handshake validator returns `protocol_violation` for all nine `validateHello` failures; `Start` maps only two helper-declared codes to `incompatible` (`OD-6` source entry).
- `ProcessAlive` has inconsistent derivations between `Status` and `Supervisor.ProcessAlive` (`OD-5` source entry; R13).

These facts explain why U1–U8 remain open. They are not a substitute for the underlying assessments.

### 3.3 Reported

Reported but not independently re-verified by this assessment:

- The owner accepted OD-2–OD-8 as governance directions; those acceptances are recorded in the register.
- The prior assessments classify R09–R13 as open or narrowed, and U1–U8 as open.
- The prior OD-5, OD-6, OD-7, and OD-8 assessments report their source conclusions; this assessment relies on them where cited.
- The register's hash and line count are a reported verification unless independently recomputed; this assessment recomputes the current register hash and line count in §10.

### 3.4 Inferred

The following are analytical consequences, not observed facts:

- Approving P17 under B would be a partial-model approval with unresolved cases; the model would not have a demonstrated outcome for U1–U8.
- Accepting A would leave P17 unapproved but would establish a standing gate for future approval.
- Rejecting OD-9 would not, by itself, approve P17; it would require a separate approval motion.
- A split under C cannot be evaluated until a minimal subset and its boundaries are defined.

### 3.5 Unestablished

- Who owns the authoritative evidence record and whether it is diagnostic-only or projection-supporting (OD-1).
- The outcome of each U1–U8 case.
- Closure of R11, R12, or R13.
- Whether P17's projections are reachable and supported for every case.
- Whether P17 is approved.
- Whether contract rev3 should be adopted.
- Whether any implementation satisfies the model.

## 4. Accepted decisions OD-2–OD-8 and their remaining boundaries

The accepted decisions are relevant, but none closes the U cases or approves P17.

| Decision | Accepted direction | What remains open |
|---|---|---|
| **OD-2** | Generation-based process identity; observations/actions/exit facts bind to one instance. | U cases still lack a full evidence record; OD-1 authority remains deferred. |
| **OD-3** | Immutable first terminal decision. | Publisher, synchronization, late-evidence reconciliation, and lifecycle interaction remain unresolved; U3/U7 remain open. |
| **OD-4** | Confirmed exit without orderly shutdown → `unavailable` unless a more specific prior failure applies; orderly shutdown → `absent`; active states cannot survive exit; action causation is internal under OD-8. | U2/U4/U5/U7 remain open; exact contract wording and implementation are separate. |
| **OD-5** | `crashed` = observed process exit where orderly shutdown is not established; no mechanism/cause/intent attribution. | R13, U4, U5, U6 remain open; source writers are not declared conforming. |
| **OD-6** | Only a structurally valid, semantically identified mandatory compatibility failure justifies `incompatible`; ambiguity/malformed/unestablished meaning remains `unavailable`. | R12, U8, U3 remain open; scope/extension-version categories remain unestablished. |
| **OD-7** | `quarantined` is internal-only, generation-bound ownership/reaping uncertainty; no public state. | R11 narrowed but not closed; U2/U6/U7 remain open; P17 still claims a public projection. |
| **OD-8** | Deliberate-action causation requires A–E evidence; association is not causation. | U1–U5 remain open; R10 is not resolved; exact action/exit evidence is not retained. |

No accepted decision supplies the missing U outcomes. None approves P17. None authorizes a contract amendment or implementation.

## 5. Remaining dependencies and blockers

### 5.1 OD-1 is deferred

OD-1 asks who owns the authoritative evidence record and whether it may support provider-state projection. The owner deferred it without selecting an authority model. The register states that R09 cannot be closed, R10 cannot be evaluated as a resolution, and U1–U7 cannot be assigned evidence-based outcomes (`OD-1` entry).

OD-1 is not a mere formatting gap. P17's projection model depends on evidence that is preserved, ordered, associated with the same process instance, and authorized to support a projection. The current source has no authoritative record; `failAll` keeps only the first read failure, `SetState` overwrites `lastFailure`, and a late terminal is dropped. A P17 approval that assumes a projection-supporting record would therefore precede the authority decision.

### 5.2 U1–U8 are open

P17 lists them as unresolved, and no classification is asserted:

| Case | Open issue |
|---|---|
| U1 | Orderly EOF during shutdown versus `helper_crashed` from the reader. |
| U2 | Shutdown timeout racing a protocol violation; no defined winner. |
| U3 | Handshake incompatible followed by a process crash; terminal-versus-sequence order. |
| U4 | Zero exit after a fatal protocol violation; fatal-wins may overstate certainty. |
| U5 | Nonzero exit after a deliberate kill; crash versus `terminated_by_action` is not established. |
| U6 | Provider `crashed` while the process is alive. |
| U7 | `Terminate` timeout after terminal projection; terminal immutability is undefined. |
| U8 | Missing mandatory capability yields `protocol_violation`/`unavailable`, while helper-declared incompatibility yields `incompatible`. |

These are not closed by OD-2–OD-8. The register says all remain OPEN, and each accepted decision explicitly preserves its dependent U cases.

### 5.3 R11, R12, R13

- **R11:** `quarantined` was a P17 projection with no path. OD-7 narrowed the question by making quarantine internal-only, but R11 is not closed while P17 §10 still claims it as a required public projection.
- **R12:** handshake compatibility still has an unresolved classification and P17/contract alignment issue; OD-6 does not close it.
- **R13:** `crashed` and live-process/provider-state semantics remain open; OD-5 does not close it.

### 5.4 P17 and contract rev3

P17 remains an unapproved frozen proposal. Contract rev3 is a proposed revision, not an accepted amendment. Neither has been revised to reflect the accepted OD-4, OD-5, OD-6, OD-7, or OD-8 directions. Approving P17 without a revision would leave the proposal inconsistent with the accepted decisions; accepting OD-9 cannot repair those inconsistencies and is not a contract amendment.

P17 also has an internal reachability problem independent of U outcomes: §10 requires confirmation of five projections including `quarantined`, while §7 has no `quarantined` path. OD-7 confirms the public claim should be corrected but does not edit P17.

## 6. Can OD-9 be decided independently?

The answer depends on what “decide OD-9” means.

### 6.1 As a negative gate: yes

The owner can decide **A** now. The evidence needed for A is already present: OD-1 is deferred, U1–U8 are open, and no U case has a determined outcome or explicit unknown rule. A decision to withhold P17 approval does not require resolving those cases; it records that they are conditions of approval.

This is a decision about the gate, not a finding that P17 is approved or rejected as a model. It is the narrowest decision the current evidence supports.

### 6.2 As approval of P17: no

Approval under **B** is not supported. It would require treating an incomplete, internally inconsistent model as dependable while:

- OD-1 has no authority decision;
- U1–U8 have no determined outcome or explicit unknown rule;
- R11, R12, and R13 remain open;
- P17's own public projections are unreachable/inconsistent in at least the `quarantined` case;
- P17 does not reflect the accepted OD-4–OD-8 directions.

Approving a partial model could be a deliberate scope decision, but the owner has not defined the partial scope, the disclosure, or the consequences. Existing evidence does not establish that such an approval is safe.

### 6.3 As a split: not ready

**C** cannot be evaluated without a proposed minimal subset. No document identifies which P17 rules, projections, or cases form a fully specified subset; no subset has been checked against U1–U8. The absence of a defined subset is a genuine blocker to a split approval, not merely an open case.

### 6.4 As a full deferral: possible but broader than required

**D** is available. It would defer P17 entirely until rev4/P18 work is authorized. It is safer than B but broader than A: A preserves a defined path to approval when the conditions are met, while D stops the model work. The evidence supports not approving P17 now; it does not require abandoning the model.

## 7. Consequences of accepting, deferring, or rejecting OD-9

### 7.1 Accept A

- P17 remains **unapproved** until OD-1–OD-8 are dispositioned and each U case has a determined outcome or explicit unknown rule.
- OD-1 remains deferred; U1–U8 remain OPEN; R11 is narrowed but not closed; R12/R13 remain open.
- No contract amendment, source change, test, build, or C1 implementation is authorized.
- The gate is recorded; future P17 approval requires a separate owner decision and evidence.
- This is the recommended outcome.

### 7.2 Accept B

- P17 would be approved as a partial model with open U cases.
- The approval would have to state its partial scope and known gaps explicitly; otherwise it treats unresolved cases as settled.
- It does not itself amend the contract or authorize implementation.
- The evidence does not support this choice now.

### 7.3 Accept C

- P17 would be split; only a minimal, fully specified subset would be approved.
- The subset and its boundaries must be defined and checked before approval.
- No such subset currently exists; implementation and contract amendment remain unauthorized.

### 7.4 Accept D

- P17 approval is deferred entirely until rev4/P18 work is authorized.
- No P17 revision, contract amendment, or implementation follows.
- This is broader than A and is not required by the evidence.

### 7.5 Defer OD-9

- OD-9 remains UNRESOLVED; P17 remains unapproved.
- The current boundary is unchanged; no gate is recorded.
- This is permissible but leaves the final approval question open. It should be distinguished from a decision to approve or reject P17.

### 7.6 Reject OD-9

Rejection is ambiguous in the current register. It could mean:

- reject the proposed gate and therefore allow consideration of P17 without the gate;
- reject P17 approval; or
- reject only the non-binding default A.

The register does not define which. Rejecting the gate would **not** approve P17 by itself, and rejecting A would not establish B or C. A separate explicit approval or deferral would still be required. This ambiguity must be resolved before any rejection is recorded.

## 8. Ambiguity that prevents a fully defensible approval decision

The following ambiguities do not prevent A, but they prevent a defensible approval of P17:

1. **“Dispositioned” and the deferred OD-1.** The gate requires OD-1–OD-8 to be dispositioned. OD-1 is deferred, not accepted, rejected, or resolved. Whether deferral satisfies “dispositioned” is undefined.
2. **Explicit unknown rule.** The gate allows an explicit `outcome_unknown` rule, but no document defines who records it, what evidence it requires, or how it maps to P17's projections. P17 uses `outcome_unknown` as a projection while the frozen provider-state set does not contain it.
3. **Approval scope.** “P17 approved” is not defined as a model-only approval, a contract approval, or an implementation authorization. The register says no implementation follows, but the approval's scope should be explicit.
4. **Minimal subset.** Alternative C has no candidate subset or boundary definition.
5. **Rejection meaning.** “Reject OD-9” does not state whether it rejects the gate, the default, or P17 approval.
6. **P17 consistency.** P17's §10 `quarantined` claim conflicts with its §7 table; P17 has not been revised for OD-4–OD-8. Approval without revision would approve a text inconsistent with accepted decisions.

## 9. Recommendation

**Recommendation: accept A — P17 is not approved while U1–U8 remain open and OD-1 remains deferred.**

This is the only decision the current evidence supports:

- It does not pretend the U cases are resolved.
- It does not approve an internally inconsistent P17.
- It does not convert an owner gate into a contract amendment or implementation authorization.
- It preserves a defined path: resolve OD-1, provide a determined outcome or explicit unknown rule for each U case, close or narrow the relevant R findings, and then bring P17 back for a separate approval decision.

If the owner is not prepared to record the standing gate, **defer OD-9** rather than approve P17 under B. A split under C would first require a defined, fully specified subset; the current record has none. D is available but is broader than the evidence requires.

The recommendation does **not** mean:

- P17 is correct or complete;
- U1–U8 are resolved;
- R11, R12, or R13 is closed;
- contract rev3 is accepted;
- C1 is authorized;
- implementation has been verified.

## 10. Verification and limitations

This is a read-only assessment. No existing artifact was modified. The decision register remains at SHA-256 `6854363b…`, 396 lines; P17 remains `01d887c1…`, 166 lines; contract rev3 remains `ed0da1a7…`, 183 lines; the lifecycle decision remains `aa136bc9…`, 74 lines. The OD-4 assessment remains `7c7cc86f…`; OD-7 `e878f561…`; OD-6 `cd9d5c65…`; OD-5 `47d19bf4…`. No source, test, configuration, or generated file was modified; no tests, builds, benchmarks, or runtime experiments were run.

**Inspection limitation:** this assessment directly inspected the register, P17, contract rev3, the lifecycle decision, and the existing assessment material. It did not re-inspect every source path cited by those documents; source-derived statements are attributed to the existing assessments and register entries rather than presented as fresh runtime verification. It did not test any U case, and it did not verify implementation behavior.

**OD-9 remains UNRESOLVED.**
