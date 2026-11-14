# C1 / P17 — Revision Plan (Read-Only)

**Date:** 2026-11-09

**Status:** Planning only. **P17 remains UNAPPROVED and unchanged.** OD-1 through OD-9 remain accepted; U1–U8 remain OPEN; contract rev3 remains frozen; C1 implementation remains unauthorized.

**Baseline:**

- Decision register: `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md`, SHA-256 `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99`, 434 lines.
- P17: `01d887c1…`, 166 lines, frozen proposal, not approved.
- Contract rev3: `ed0da1a7…`, 183 lines, proposed revision, not accepted.
- Reconciliation matrix: `3a06f8e3…`, 281 lines; C-01 through C-10 are the basis of this plan.

This plan does not edit P17, amend contract rev3, update the register, change source, run tests, or authorize implementation.

## 1. Purpose

The reconciliation identified ten concrete P17 corrections. This plan converts them into a reviewable revision sequence:

1. map each conflict to the exact P17 sections;
2. separate corrections dictated by accepted owner decisions from decisions still requiring technical judgment;
3. identify text that would exceed OD-1 through OD-9 authority;
4. preserve every evidence-loss finding and U case; and
5. define verification criteria for a future P17 revision.

No P17 text is revised here.

## 2. Authority boundary

The accepted decisions fix the direction of the revision, not every implementation detail.

### 2.1 Directly controlled by accepted decisions

The following outcomes are not open for P17 to choose again:

- Evidence is projection-supporting as well as diagnostic, with the supervisor as designated authority for process-lifecycle evidence and publication (OD-1).
- Client/session observations must be represented and retained (OD-1).
- Records bind to the process generation (OD-1, OD-2).
- Event order and observation order are distinct (OD-1, OD-8).
- A first terminal decision is immutable; later evidence remains retained (OD-3).
- Confirmed exit without orderly shutdown projects `unavailable` unless a prior justified failure applies; orderly shutdown projects `absent`; no exit established does not become `absent` (OD-4).
- `crashed` requires an observed process exit and no established orderly shutdown; a live process cannot be `crashed` (OD-5).
- Compatibility classification follows OD-6's accepted categories; missing mandatory capability is `incompatible`, structural ambiguity/malformed cases are `unavailable`, and unclassified cases remain unclassified.
- `quarantined` is internal-only, generation-bound ownership/reaping uncertainty, not a public provider state (OD-7).
- Causal attribution requires OD-8's A–E evidence; temporal association is not causation.
- P17 remains unapproved while U1–U8 remain open without determined outcomes or explicit unknown rules (OD-9).

### 2.2 Not authorized by OD-1 through OD-9

A P17 revision must not:

- add a public provider state or rename the thirteen-value set;
- make `outcome_unknown` a public provider state without a separate contract decision;
- make `quarantined` public;
- make `terminated_by_action` public;
- treat `crashed` as a live-process provider condition;
- replace OD-2's generation identity with PID or session alone;
- use `unavailable` as a way to discard a more specific, independently justified failure;
- declare U1–U8 resolved without evidence or an explicit unknown rule;
- amend contract rev3;
- authorize source changes, tests, builds, benchmarks, or runtime experiments.

## 3. Conflict-to-correction matrix

### C-01 — `crashed` without an observed exit

| Field | Plan |
|---|---|
| **P17 locations** | §4 S3 and S6; §5 O6, O7, O8, O12, O18; §7 rules 2–3 and combination table. |
| **Conflict** | Fatal protocol/transport observations are projected to `crashed` even when no process exit is observed. OD-5 restricts `crashed` to an observed process exit with no established orderly shutdown. |
| **Required correction** | Remove any `crashed` projection that is not conditioned on an observed exit. Keep fatal protocol/transport observations as internal evidence. Do not infer process exit and do not assign an unsupported terminal or public projection. The public-state treatment when exit is not established remains explicitly unresolved. |
| **Authority** | OD-1, OD-4, OD-5. |
| **Additional technical decision** | Which public state, if any, follows a fatal observation when process exit is not established, including whether a prior state remains valid. OD-4 and OD-5 set the boundary but do not choose the fallback. |
| **Verification** | No P17 rule or table maps fatal protocol/transport directly to `crashed`; every `crashed` row has an observed-exit condition and no orderly-shutdown condition; the no-exit public treatment is explicitly unresolved rather than assigned an unsupported terminal. |

### C-02 — Nonzero exit while stopping projected to `absent`

| Field | Plan |
|---|---|
| **P17 locations** | §5 O5; §7 rule 4 and combination table. |
| **Conflict** | P17 says a nonzero exit observed while stopping becomes `absent`. OD-4 says a confirmed exit without orderly shutdown is `unavailable` unless a more specific prior failure is independently justified. |
| **Required correction** | Replace the `absent` projection for this path with OD-4's `unavailable` fallback, preserving a prior justified failure where applicable. Do not treat a nonzero exit as proof of orderly shutdown. |
| **Authority** | OD-4, OD-5, OD-8. |
| **Additional technical decision** | How a prior justified failure is identified and retained; whether a specific failure classification survives the exit. |
| **Verification** | No nonzero-exit-while-stopping path yields `absent`; the prior-failure exception is explicit and evidence-bound. |

### C-03 — `outcome_unknown` used as a public projection

| Field | Plan |
|---|---|
| **P17 locations** | §7 rules 4b and 4d; §7 combination table; related U4/U5 references. |
| **Conflict** | P17 projects zero/no-action and ambiguous-nonzero exits to `outcome_unknown`. OD-4 requires `unavailable` for a confirmed exit without orderly shutdown unless a prior justified failure applies. `outcome_unknown` is not in the frozen provider-state set. |
| **Required correction** | Define `outcome_unknown` as internal only, or remove it from the public projection column. For a confirmed exit with no orderly shutdown, use OD-4's public fallback. For no established exit, do not assign a public terminal; the treatment remains unresolved under C-01. |
| **Authority** | OD-4, OD-9. |
| **Additional technical decision** | The exact internal unknown-case representation and the evidence required to record it; no public state is added. |
| **Verification** | No public projection uses `outcome_unknown`; every confirmed-exit/no-orderly-shutdown row uses OD-4's policy; the internal label does not leak into the provider-state set. |

### C-04 — Handshake split for missing mandatory capability

| Field | Plan |
|---|---|
| **P17 locations** | §5 O15/O16; §7 rule 1; §9 U8. |
| **Conflict** | P17 preserves a blanket `protocol_violation` → `unavailable` path. OD-6 classifies a missing mandatory service capability, unknown required core capability, and unsupported required version as `incompatible`; duplicate identities, invalid semantic descriptors, and malformed version/limits remain `unavailable`; unclassified categories remain unclassified. |
| **Required correction** | Replace the blanket handshake mapping with OD-6's category table. Do not use a failure-code lookup as the classification; state the observed condition and its category. Preserve distinct internal evidence for helper-declared and client-detected instances. |
| **Authority** | OD-6. |
| **Additional technical decision** | Classification of each of the nine `validateHello` paths and the still-unclassified descriptor/extension cases. The plan does not invent those classifications. |
| **Verification** | Every handshake path has an explicit OD-6 category or is marked unclassified; missing mandatory capability yields `incompatible`; malformed/ambiguous cases yield `unavailable`; no generic default hides an unknown category. |

### C-05 — `quarantined` public projection

| Field | Plan |
|---|---|
| **P17 locations** | §10 challenge 5; the §7 projection table's omission and §10's public claim. |
| **Conflict** | P17 lists `quarantined` among required projections. OD-7 makes it internal-only and forbids a public provider-state projection. |
| **Required correction** | Remove `quarantined` from the public projection challenge. Describe it as an internal supervisor ownership/reaping condition for the exact generation, not a provider state. Do not add a fourteenth state. |
| **Authority** | OD-7. |
| **Additional technical decision** | Internal representation, transition, retention, and release evidence; none is a public state. |
| **Verification** | P17 contains no public `quarantined` projection; §10 no longer lists it as a required public state; no enum change. |

### C-06 — Diagnostic/projection role contradiction

| Field | Plan |
|---|---|
| **P17 locations** | §4 S4; §5 O1, O2, O19 (and the discovery/output/request class boundaries). |
| **Conflict** | S4 says diagnostics never select a projection, but O1/O2 call `setUnavailable` and O19 publishes `ready` or `degraded`. |
| **Required correction** | Treat diagnostic and projection roles as separate. Evidence may be retained diagnostically and may support a provider-state projection only under an explicit, justified publication rule; a diagnostic label alone does not select a provider state. Either reclassify O1/O2/O19 as projection-bearing observations with explicit rules, or stop deriving provider state from them. Preserve the underlying evidence in either case. |
| **Authority** | OD-1, OD-4, OD-5. |
| **Additional technical decision** | Which discovery/request conditions are diagnostic, which are projection-bearing, and which may legitimately serve both roles under separate rules; exact replacement handling. |
| **Verification** | Every projecting observation has an explicit rule; no diagnostic label or diagnostic-only record selects a provider state by itself; diagnostic evidence remains retained; evidence may support both roles only through an explicit publication rule. |

### C-07 — Missing generation binding

| Field | Plan |
|---|---|
| **P17 locations** | §2 Observation/Action/EvidenceLog; §1.1 preservation; §7 references to `seq`. |
| **Conflict** | P17 records have `id` and `seq` but no process generation or instance binding. OD-1 and OD-2 require every relevant record to bind to the process generation, with session scope and PID corroboration. |
| **Required correction** | Add the accepted identity binding to every relevant observation, action, and exit fact. Use generation as the process-lifetime key; session provides outer scope; PID is metadata, not identity. |
| **Authority** | OD-1, OD-2. |
| **Additional technical decision** | Exact field/schema representation and how a generation is allocated and retained; no PID-only or session-only substitute. |
| **Verification** | Every record type carries or references the generation; delayed old-instance evidence cannot modify a replacement; the generation is unique across supervisor replacement. |

### C-08 — No terminal latch or late-evidence rule

| Field | Plan |
|---|---|
| **P17 locations** | §1.2 preservation; §2 EvidenceLog; §7 precedence; §10 challenge 4. |
| **Conflict** | P17 says precedence selects a projection but does not define a terminal latch, linearization point, or late-evidence handling. OD-3 requires an immutable first terminal decision and retained later evidence. |
| **Required correction** | Add an explicit terminal rule: one immutable first terminal per process instance, one linearization point, later evidence retained but not republished as the terminal, and old-instance evidence unable to modify a replacement. The immutability applies to the terminal decision; it does not by itself decide whether later evidence may affect diagnostics or non-terminal projections. |
| **Authority** | OD-1, OD-2, OD-3. |
| **Additional technical decision** | Exact synchronization/linearization mechanism, terminal set, late-evidence buffering, and reconciliation behavior; whether and how later evidence affects diagnostics or non-terminal projections remains unresolved. |
| **Verification** | A later exit cannot overwrite an established terminal; late evidence remains recorded; old-instance evidence cannot modify a new generation; no terminal overwrite path remains; P17 does not claim all later evidence is inert for non-terminal projections. |

### C-09 — “Known action cause” lacks OD-8’s standard

| Field | Plan |
|---|---|
| **P17 locations** | §7 rule 4c; §7 combination table; §9 U5; §4 S5. |
| **Conflict** | P17 uses “known action cause” as if an action or stop intent established causation. OD-8 requires A–E evidence; temporal association is insufficient. |
| **Required correction** | Replace “known action cause” with the OD-8 A–E test. Keep `terminated_by_action` internal and public projection governed by OD-4. Record action identity, outcome, exit fact, event order, mechanism operation, and competing explanations. |
| **Authority** | OD-4, OD-8. |
| **Additional technical decision** | Exact action/outcome fields and evidence collection needed to apply A–E; no implementation is authorized. |
| **Verification** | No P17 rule uses action-before-exit as sufficient causation; every causal claim requires A–E; U5 remains open; no public action-causation state. |

### C-10 — `seq` has no allocator or order semantics

| Field | Plan |
|---|---|
| **P17 locations** | §2 `seq` fields; §7 precedence; §9 U1–U7 dependencies. |
| **Conflict** | P17 uses `seq` without defining who allocates it, whether it records event or processing order, or how concurrent writers are ordered. OD-1 requires event order and observation order to remain distinct; OD-8 requires reliable event order. |
| **Required correction** | Define the allocator and ordering domain. A sequence allocated at observation/processing time establishes observation order only; it must not be presented as event order. State how concurrent entries are linearized and retained. |
| **Authority** | OD-1, OD-8. |
| **Additional technical decision** | Allocator design, concurrency mechanism, and clock/sequence policy. The plan does not select one. |
| **Verification** | P17 distinguishes event order from observation/processing order; no request counter is used as evidence order; concurrent writes have a defined record; causal conclusions cannot rest on processing order alone. |

### G-01 — OD-9 historical OD-1 condition (read-only consistency check)

| Field | Plan |
|---|---|
| **Issue** | OD-9 was accepted with wording tied to OD-1 being deferred. OD-1 has since been accepted. U1–U8 still independently prevent P17 approval, but a superseded condition must not look current. |
| **Required treatment** | Flag the condition for a read-only consistency review. Do not rewrite the register as part of the P17 revision. Any correction to OD-9's historical wording requires separate authorization. |
| **Verification** | The revision plan identifies any superseded OD-1 condition; no register edit is bundled into the P17 revision. |

## 4. Corrections dictated by owner decisions vs technical decisions

| Conflict | Direction fixed by OD-1–OD-9 | Technical decision still needed |
|---|---|---|
| C-01 | `crashed` requires an observed exit; no live-process `crashed`. | Public state, if any, when fatal evidence exists but exit is not established, including whether a prior state remains valid. |
| C-02 | Nonzero exit while stopping → `unavailable` unless prior justified failure. | Prior-failure identification and retention. |
| C-03 | No public `outcome_unknown`; confirmed exit/no orderly shutdown → OD-4 fallback. | Internal unknown representation and evidence rule. |
| C-04 | Apply OD-6 categories; missing mandatory capability → `incompatible`. | Classification of all nine handshake paths and unclassified categories. |
| C-05 | `quarantined` internal-only; no public state. | Internal transition, retention, and release evidence. |
| C-06 | A diagnostic label alone cannot select a projection; evidence may serve both roles under an explicit rule. | Reclassify O1/O2/O19 or change their handling. |
| C-07 | Every relevant record binds to generation. | Schema/field representation and allocator. |
| C-08 | Immutable first terminal; later evidence retained. | Linearization, terminal set, late-evidence mechanics, and effects on non-terminal projections. |
| C-09 | A–E required for causal attribution. | Action/outcome fields and evidence collection. |
| C-10 | Event order and observation order are distinct. | Sequence allocator and concurrency policy. |

A future revision may propose technical answers, but it must mark them as proposals, not as accepted owner decisions.

## 5. Text that would exceed OD-1 through OD-9 authority

The following would be out of scope or a silent policy change:

1. Treating `outcome_unknown` as a public provider state.
2. Treating `quarantined` as a public provider state or diagnostic exposure.
3. Adding `terminated_by_action` to the public projection set.
4. Allowing `crashed` to describe a live process or a fatal observation without an observed exit.
5. Replacing generation identity with PID or session alone.
6. Using `unavailable` to erase a prior, independently justified failure classification.
7. Declaring U1–U8 resolved by precedence, source assertion, or a new default.
8. Declaring `outcome_unknown` without a defined rule or evidence meaning.
9. Changing the thirteen-value provider-state set.
10. Amending contract rev3 or presenting a contract change as implied by P17.
11. Treating OD-1's authority model as implementation authority.
12. Claiming test, build, runtime, or implementation evidence that does not exist.
13. Rewriting OD-9's historical condition as part of the P17 revision instead of a separate consistency review.

## 6. Preservation requirements

A revision must preserve, not delete or relabel away:

### 6.1 Evidence-loss findings

- **E1:** `failAll` keeps only the first read failure.
- **E2:** `SetState` overwrites `lastFailure`.
- **E3:** the stopping branch discards a nonzero exit from classification.
- **E4:** termination intent is not recorded; `Kill` aliases `Terminate`.
- **E5:** `waitForReader` loses join versus timeout.
- **E6:** stream closure is marked crashed without termination.
- **E7:** late terminal evidence can be dropped by dispatch (`client.go:339–342`).
- **E8:** cancellation and cleanup paths ignore errors from `Cancel` and `Terminate`.
- **E9:** `Shutdown` calls `Terminate`, then clears lifecycle and transport state without establishing successful process reaping.

These remain findings until separately resolved by evidence or implementation. A revision may change their surrounding text, but must not omit them or claim they are fixed.

### 6.2 Unresolved cases

- **U1** orderly EOF versus reader crash.
- **U2** shutdown timeout versus protocol violation.
- **U3** handshake incompatible followed by crash.
- **U4** zero exit after fatal protocol violation.
- **U5** nonzero exit after deliberate kill.
- **U6** provider `crashed` while process alive.
- **U7** `Terminate` timeout after terminal projection.
- **U8** missing mandatory capability versus helper-declared incompatibility.

No U case may be closed by the revision plan. Each must remain OPEN or receive a separately justified determined outcome or explicit unknown rule.

### 6.3 Existing status

P17 remains unapproved; contract rev3 remains frozen; the thirteen provider states remain unchanged; no implementation is authorized.

## 7. Proposed revision order

This is a document-revision order, not an implementation sequence.

### Phase 0 — Freeze and verify

- Confirm register `3c64797b…`, P17 `01d887c1…`, contract `ed0da1a7…`.
- Preserve the reconciliation matrix `3a06f8e3…`.
- Perform the read-only OD-9 historical-condition check in G-01; flag, do not rewrite.
- Do not edit any artifact.

### Phase 1 — Authority, identity, and ordering

Address **C-07, C-08, C-10** first:

- generation-bound records;
- supervisor authority and client observation retention;
- immutable terminal latch and late-evidence rule;
- `seq` allocator and event-versus-observation ordering.

Reason: terminal projection and compatibility text cannot be evaluated on a record model that lacks identity or order.

### Phase 2 — Exit and terminal projections

Address **C-01, C-02, C-03, C-09**:

- `crashed` requires observed exit;
- nonzero exit while stopping → `unavailable` unless prior failure;
- `outcome_unknown` internal only;
- OD-4 public fallback for confirmed exits;
- OD-8 A–E standard for action causation.

Reason: these are the critical terminal-state defects and the core of P17's projection table.

### Phase 3 — Compatibility and quarantine

Address **C-04, C-05**:

- OD-6 handshake categories;
- `quarantined` internal-only;
- no new public state.

Reason: these are bounded policy corrections with direct accepted-decision authority.

### Phase 4 — Diagnostics and observation inventory

Address **C-06**:

- diagnostic labels do not select by themselves;
- every projection has an explicit publication rule;
- diagnostic evidence remains retained and may also support projection only under that rule;
- O1/O2/O19 reclassified or rehandled.

Reason: the inventory must not conflate diagnostic retention with projection authority, and must not let the existing S4 contradiction survive.

### Phase 5 — Evidence-loss and U-case preservation

- retain E1–E9;
- retain U1–U8;
- add no unsupported outcome;
- apply explicit unknown rules only where separately defined.

### Phase 6 — Cross-decision verification

- check every P17 claim against OD-1 through OD-9;
- confirm no text exceeds authority;
- confirm no contract or public-state change;
- confirm no implementation claim;
- produce a reviewable revised proposal only after separate authorization.

## 8. Verification criteria for a future P17 revision

A revision is ready for review only if all of the following hold:

1. Every C-01 through C-10 has an explicit correction or a documented reason it is not applicable.
2. No P17 rule maps fatal protocol/transport to `crashed` without an observed exit.
3. No nonzero exit observed while stopping maps to `absent`.
4. No public projection uses `outcome_unknown`; its internal meaning is explicit.
5. Confirmed exit without orderly shutdown follows OD-4's `unavailable` policy unless a prior justified failure applies.
6. Every handshake path is classified under OD-6 or marked unclassified.
7. `quarantined` is nowhere a public projection.
8. No diagnostic label or diagnostic-only record selects a provider state by itself; every projection has an explicit publication rule; evidence may serve both roles only under that rule.
9. Every relevant record is generation-bound.
10. A first terminal is immutable and later evidence is retained; the revision does not claim all later evidence is inert for non-terminal projections.
11. “Known action cause” has been replaced by the OD-8 A–E standard.
12. `seq` distinguishes event order from observation/processing order.
13. E1–E9 remain present and unresolved where applicable.
14. U1–U8 remain OPEN or have separately justified dispositions.
15. The thirteen-value provider-state set is unchanged.
16. Contract rev3 is not amended.
17. No source, test, build, runtime, or implementation work is claimed.
18. The revision does not turn an OD-1 authority model into an implementation authorization.
19. OD-9's historical OD-1 condition has been explicitly checked; no register edit is bundled into the P17 revision.

## 9. What is still not decided

Even a complete P17 revision will not decide:

- runtime interleavings for U1–U7;
- the exact compatibility classification of every handshake path;
- the internal quarantine transition and retention mechanics;
- the terminal linearization mechanism;
- the sequence allocator implementation;
- the action/exit evidence collection mechanism;
- the public projection, if any, when fatal evidence exists without an observed exit;
- whether any implementation conforms to the revised model.

Those remain separate technical or evidentiary questions.

## 10. Conclusion

The revision plan is complete as a plan. It does not revise P17.

The next permissible artifact is a **P17 revision proposal** that follows this order and passes the verification criteria, after separate authorization. P17 remains unapproved, U1–U8 remain open, contract rev3 remains frozen, and no implementation is authorized.

## 11. Verification and limitations

This is a planning document. No artifact other than this plan was modified. The register remains at `3c64797b…`, 434 lines; P17 remains `01d887c1…`; contract rev3 remains `ed0da1a7…`; the reconciliation matrix remains `3a06f8e3…`; no source, test, configuration, or generated file was changed.

**Limitations:** the plan relies on the reconciliation matrix and existing source assessment; it does not re-run the source audit, decide U1–U8, select technical mechanisms, or verify implementation. It does not approve P17 or authorize a P17 edit.

**No implementation was performed or authorized.**
