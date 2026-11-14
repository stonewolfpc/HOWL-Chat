# C1 / P17 — P17-to-OD Reconciliation Matrix (Read-Only)

**Date:** 2026-11-09

**Status:** Read-only reconciliation. **P17 remains UNAPPROVED.** OD-1 through OD-9 are accepted; U1–U8 remain OPEN; contract rev3 remains frozen; C1 implementation remains unauthorized.

**Baseline:**

- Decision register: `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md`, SHA-256 `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99`, 434 lines.
- P17: `01d887c1…`, 166 lines, frozen author proposal, not approved.
- Contract rev3: `ed0da1a7…`, 183 lines, proposed revision, not accepted.
- Counts: **9 ACCEPTED** (OD-1 through OD-9); **0 DEFERRED**; **0 UNRESOLVED**.

This document does not edit P17, amend the contract, update the register, change source, run tests, or authorize implementation.

## 1. Method and disposition vocabulary

P17 mixes three kinds of statements:

1. **Design requirements** — what the proposed evidence model is intended to do.
2. **Current-handling descriptions** — what the existing source does today.
3. **Unresolved-case declarations** — cases P17 explicitly leaves open.

The reconciliation treats these separately. A P17 rule that is sound as a future requirement may still conflict with an accepted decision if it states a projection or authority the decision does not permit.

| Disposition | Meaning |
|---|---|
| **Aligned** | Consistent with the accepted decision; no conflict identified. |
| **Qualify** | Consistent only with an explicit limitation, definition, or additional rule. |
| **Conflict** | P17’s claim cannot stand as written against an accepted decision. |
| **Open** | Correctly left unresolved; no outcome may be assigned now. |
| **Unsupported** | The current source does not implement the claim, or the record/field does not exist. This is not necessarily a policy conflict. |

The matrix is finding-level, not a P17 rewrite. No row authorizes a code change.

## 2. Accepted decision controls

| Decision | Control relevant to P17 |
|---|---|
| **OD-1** | Evidence record is projection-supporting as well as diagnostic; supervisor is the designated authority for process-lifecycle evidence and provider-state publication; client/session observations must be represented and retained; records bind to process generation; event order and observation order are distinct; missing evidence is not proof of non-occurrence; only justified evidence supports public projection; internal evidence does not create public states. |
| **OD-2** | Composite session + supervisor generation + PID; generation is the process-lifetime key; delayed observations from an old instance must not modify a replacement. |
| **OD-3** | One immutable first terminal decision per process instance; one linearization point; later observations retained but cannot overwrite the terminal; late evidence must not silently disappear. |
| **OD-4** | Confirmed exit without orderly shutdown → `unavailable` unless a more specific prior failure is independently justified; orderly shutdown → `absent`; confirmed exit with justified prior failure preserves that classification; no exit established → do not pretend absent; deliberate-action causation is internal under OD-8, with no public state; no new public state. |
| **OD-5** | `crashed` means observed process exit where orderly shutdown is not established; no physical mechanism, cause, intent, or deliberate-action attribution; a live process must not be represented as `crashed`; `ProcessAlive` is independent of provider state, but the five-state mapping is not accepted. |
| **OD-6** | Only a structurally valid, semantically identified failure of a mandatory compatibility requirement justifies `incompatible`; structural ambiguity, malformed content, and unestablished compatibility meaning remain `unavailable`; helper-declared and client-detected instances of the same established mandatory failure share the public state while retaining distinct internal evidence; missing mandatory service capability, unknown required core capability, and unsupported required version → `incompatible`; duplicate identities, invalid semantic descriptors, malformed version/limits → `unavailable`; misplaced descriptors and unsupported extension versions remain unclassified; no new public state. |
| **OD-7** | `quarantined` is internal supervisor ownership/reaping uncertainty for the exact generation; not a provider state; no public projection; termination timeout alone is not quarantine; later exit does not automatically establish reaping; P17 §10's public claim must be corrected. |
| **OD-8** | Temporal association is not causation; A–E evidence is required for causal attribution; A-O is association only; A-U remains unresolved; identity, action outcome, exit fact, event order, mechanism operation, and competing explanations matter; U5 remains open. |
| **OD-9** | P17 is not approved while U1–U8 remain open without determined outcomes or explicit unknown rules; acceptance of OD-9 does not approve P17. |

## 3. Executive reconciliation

P17 is not uniformly inconsistent with the accepted decisions. Its broad constraints align; its specific projection rules do not.

### 3.1 Direct conflicts requiring correction before approval

1. **`crashed` from fatal protocol or transport without an observed exit.** P17 §7 rules 2 and 3, and §7's combination table, project fatal protocol/transport to `crashed`. OD-5 permits `crashed` only for an observed process exit where orderly shutdown is not established. This also conflicts with P17's own S3 and S6.
2. **`absent` for a nonzero exit observed while stopping.** P17 O5 says the provider becomes `absent` despite a nonzero exit. OD-4 requires `unavailable` unless a more specific prior failure is independently justified.
3. **Public/`outcome_unknown` projection for confirmed exits with no established orderly shutdown.** P17 §7 maps zero/no-action and ambiguous nonzero exits to `outcome_unknown`. OD-4 requires `unavailable` for a confirmed exit without orderly shutdown unless a prior justified failure applies; `outcome_unknown` is not a public provider state. P17 must define whether this label is internal only.
4. **Handshake `protocol_violation` → `unavailable` for a missing mandatory capability.** P17 O16 and §7 rule 1 preserve the current split. OD-6 classifies missing mandatory service capability as `incompatible`, subject to the accepted mandatory-compatibility principle.
5. **`quarantined` as a required projection.** P17 §10 challenge 5 lists `quarantined` among five required projections. OD-7 makes it internal-only with no public projection and requires P17's claim to be corrected.
6. **Diagnostics selecting provider states.** P17 S4 says diagnostics never select a projection, but O1/O2 report diagnostic discoveries as `setUnavailable`, and O19 reports request diagnostics as `ready` or `degraded`. Either the observations are misclassified as diagnostics or the current handling conflicts with S4.
7. **Missing generation binding.** P17's `Observation` and `Action` records have `id` and `seq` but no generation/process-instance field. OD-1 and OD-2 require every relevant record to bind to the process generation. This is a direct model omission, not merely an implementation gap.

### 3.2 Qualifications, not necessarily conflicts

- **Terminal immutability:** P17 says precedence selects a projection but does not define a terminal latch or linearization point. OD-3 requires one. P17 needs an explicit terminal rule.
- **Late evidence:** P17 says precedence never deletes evidence, but it does not define retention, buffering, or reconciliation after terminal publication. OD-3 requires later evidence to be retained.
- **Action causation:** P17 uses “known action cause” without OD-8's A–E standard. The term must be replaced by the accepted evidence test.
- **Missing evidence:** P17 does not say what an absent record means. OD-1 says missing evidence must not automatically mean the event did not occur.
- **Authority:** P17 names records and a log but no accountable owner or publisher. OD-1 designates the supervisor for process-lifecycle evidence and publication; client/session observations must still be retained.
- **Retention:** P17 defines an append-only log but no retention/loss rule. OD-1 and OD-3 require one before projection can be authoritative.
- **`ProcessAlive`:** P17 does not resolve the five-state mapping. OD-5 accepts independence but not the current mapping.

### 3.3 Open, not to be resolved by this matrix

U1–U8 remain open. P17 correctly declares them unresolved, but several of its §7 rules nevertheless state provisional winners. Those provisions require qualification after the U cases are decided or given explicit unknown rules.

## 4. Section 1 — Constraints

| P17 claim | OD control | Disposition | Required qualification |
|---|---|---|---|
| §1.1 All observations are preserved in an append-only log. | OD-1, OD-2, OD-3 | **Aligned / Unsupported** | The proposal is aligned, but the current source has no log and loses evidence (E1–E6). The record must bind to generation and retain client/session observations. |
| §1.2 Precedence selects a projection; it never deletes evidence. | OD-1, OD-3 | **Qualify** | Preservation is not sufficient; the record must be authoritative, retained, and ordered. Precedence cannot overwrite or erase a terminal or late evidence. |
| §1.3 Severity is assigned by the observing site, not looked up from the failure code. | OD-1, OD-5, OD-6, OD-8 | **Aligned** | Site-assigned severity does not itself establish public projection or compatibility classification. It must remain internal evidence. |
| §1.4 An action is not evidence, and an exit fact is not an action. | OD-8 | **Aligned** | Action occurrence/outcome is evidence of an action, not proof of causation; exit fact is separate. This supports the required A–E standard. |
| §1.5 Where the source cannot establish a classification, the case stays unresolved. | OD-9, OD-8 | **Aligned** | P17 §7 must not assign a projection merely because a U case is listed as unresolved. A declared `outcome_unknown` rule must be explicit and internal/public-scope defined. |

## 5. Section 2 — Model

| P17 claim | OD control | Disposition | Required qualification |
|---|---|---|---|
| `Observation { id, seq, site, code, class, severity, detail }` | OD-1, OD-2 | **Conflict / omission** | No process generation or instance identity. `seq` has no defined allocator or ordering domain. Supervisor authority must be represented. |
| `Action { id, seq, site, intent, grace, completion }` | OD-1, OD-8 | **Qualify** | Intent/completion are useful, but the model does not define requested/accepted/attempted/completed/rejected/timed-out/failed status or action outcome. Caller identity and publisher authority are absent. |
| `EvidenceLog = append-only [Observation | Action | ExitFact]` | OD-1, OD-3 | **Qualify / Unsupported** | No owner, retention, loss policy, terminal latch, late-evidence buffer, or publication rule. The current source has no such log. |
| “The proposal adds [fields] at the observing site, not inside `Failure`.” | OD-1 | **Qualify** | The record must still be submitted to the supervisor authority; a site-local field does not establish record integrity or publication. |

## 6. Section 3 — Evidence classes

| P17 class | OD control | Disposition | Required qualification |
|---|---|---|---|
| handshake (fatal capable) | OD-6 | **Qualify** | “Handshake fatal” is not a compatibility classification. It must be split into mandatory compatibility failures, structural ambiguity/malformed cases, and unclassified cases under OD-6. |
| protocol (fatal capable) | OD-5 | **Conflict at projection** | The class can be fatal, but fatal protocol does not establish process exit. §7's `crashed` projection is not supported by OD-5. |
| transport (fatal capable) | OD-5 | **Conflict at projection** | A transport failure can be fatal, but it does not establish process exit. §7's `crashed` projection is not supported by OD-5. |
| exit (result of `Wait`) | OD-4, OD-5 | **Aligned** | An exit fact is not proof of orderly shutdown or causation; the actual status/exit evidence must be preserved. |
| timeout | OD-8 | **Qualify** | A timeout is an action outcome, not a crash or terminal state; the action, deadline, and completion/failure must be recorded. |
| request | OD-1, OD-4 | **Aligned** | Request diagnostics do not establish process exit or public provider state. |
| output | OD-1, OD-4 | **Aligned** | Output validation is not a lifecycle outcome; it may trigger cleanup, but not a public crash claim. |
| discovery | OD-1, OD-4 | **Qualify** | Discovery failure can be diagnostic, but current `setUnavailable` means it selects a provider state; either reclassify the observation or record it as a projection rule. |

## 7. Section 4 — Severity rules

| P17 rule | OD control | Disposition | Required qualification |
|---|---|---|---|
| S1 severity bound at site; timeout is diagnostic at request site and action outcome at supervisor | OD-1, OD-8 | **Aligned** | The same code may have different meanings by site; the record must retain site and action context. |
| S2 one code may appear at more than one severity | OD-1, OD-5, OD-6 | **Aligned** | Severity is not the public state; code names do not establish meaning. |
| S3 fatal is not synonymous with crash; handshake fatal ≠ crash | OD-5 | **Aligned** | This is correct, but P17 §7 rules 2–3 contradict it by mapping fatal protocol/transport directly to `crashed`. |
| S4 diagnostics never select a projection | OD-1, OD-4 | **Conflict with current handling** | O1/O2 (`setUnavailable`) and O19 (`ready`/`degraded`) are classified as diagnostics but select states. The classification or handling must change. |
| S5 an action never overrides a pre-existing fatal observation | OD-3, OD-4, OD-8 | **Qualify** | Directionally consistent, but P17 has no terminal latch, no publisher rule, and no A–E causal test. It cannot be applied as written. |
| S6 provider state is not derived from disposition or an unobserved outcome | OD-4, OD-5 | **Conflict with §7** | P17 §7 derives `crashed` from fatal protocol/transport without an observed exit, and derives `absent` from a stop-intent branch. S6 is correct; §7 must be reconciled to it. |

## 8. Section 5 — Observation inventory

| ID | P17 claim/current handling | OD control | Disposition | Required qualification |
|---|---|---|---|---|
| **O1** | `helper_missing` is discovery/diagnostic; current handling calls `setUnavailable`. | OD-1, OD-4 | **Conflict with S4** | A diagnostic discovery is selecting `unavailable`. Either classify it as a projection-bearing discovery rule or stop deriving state from it. |
| **O2** | `helper_start_failed` is discovery/diagnostic; current handling calls `setUnavailable`, no process. | OD-1, OD-4 | **Conflict with S4** | Same issue as O1; no process does not by itself establish a provider-state projection. |
| **O3** | `helper_crashed` is an exit fact; current handling publishes `crashed`, process exited. | OD-5 | **Aligned** | `crashed` is defensible only if the exit was observed and orderly shutdown was not established; the code name alone is not evidence. |
| **O4** | Exit fact; current handling publishes `absent`. | OD-4 | **Qualify** | `absent` requires independently established orderly shutdown; the current `stopping` branch is intent-based and is not proof. |
| **O5** | Nonzero exit while stopping; current handling publishes `absent`. | OD-4 | **Conflict** | OD-4 requires `unavailable` unless a more specific prior failure is independently justified. |
| **O6** | `internal_error` → `helper_crashed`, transport/fatal; `failAll`, no kill. | OD-5, OD-8 | **Conflict at projection** | A transport failure does not establish process exit; `crashed` cannot be published merely from this observation. |
| **O7** | `invalid_frame`, `frame_too_large`, `protocol_violation` fatal; `failAll + Kill`. | OD-5, OD-8 | **Conflict at projection** | Fatal protocol/transport evidence may justify termination, but not `crashed` without an observed exit; action causation needs OD-8. |
| **O8** | Decode errors protocol/fatal; `failAll + Kill`. | OD-5, OD-8 | **Conflict at projection** | Same as O7. |
| **O9** | `protocol_violation` fatal; `protocolViolation + Kill`. | OD-5, OD-8 | **Conflict at projection** | Same as O7. |
| **O10** | `protocol_violation` fatal; `protocolViolation + Kill`. | OD-5, OD-8 | **Conflict at projection** | Same as O7. |
| **O11** | `protocol_violation` fatal; `protocolViolation + Kill`. | OD-5, OD-8 | **Conflict at projection** | Same as O7. |
| **O12** | `helper_crashed` transport/fatal; returned only, no termination. | OD-5 | **Conflict at projection** | No observed exit; `crashed` is not established. The lack of termination also leaves the process/action outcome open. |
| **O13** | Request `timeout` diagnostic; request fails. | OD-8 | **Aligned** | Must not be published as process exit or crash; action/deadline context should be retained. |
| **O14** | `cancelled` request diagnostic; request fails. | OD-1 | **Aligned** | No lifecycle projection follows from cancellation alone. |
| **O15** | Handshake `incompatible`, `minor_unsupported`, others; Terminate then `incompatible`/`unavailable`. | OD-6 | **Qualify** | The two compatibility codes may justify `incompatible`; “others” must be classified under OD-6, not lumped into a generic default. |
| **O16** | All `validateHello` failures → `protocol_violation`; falls to `unavailable`. | OD-6 | **Conflict** | A missing mandatory capability is `incompatible` under OD-6; duplicate identities, invalid descriptors, malformed version/limits remain `unavailable`; misplaced/unsupported extension cases remain unclassified. |
| **O17** | Wire failures handshake/load diagnostic; provider unavailable. | OD-6 | **Qualify** | `unavailable` is acceptable for structural ambiguity/malformed cases; a mandatory compatibility failure must be separated. |
| **O18** | `Synthesize` send failure publishes `crashed` while process may be alive. | OD-5 | **Conflict** | A live process must not be `crashed`. This is the U6/R13 conflict directly. |
| **O19** | Request `timeout`, `cancelled`, `degraded`; current handling publishes `ready` or `degraded`. | OD-1, S4 | **Conflict with S4** | Request diagnostics are selecting provider states. Reclassify or treat the handling as a projection rule, not a diagnostic. |
| **O20** | Supervisor `timeout` is action outcome; process state unchanged. | OD-8 | **Aligned** | Action outcome must be recorded with intent/result; a timeout is not a terminal state. |

## 9. Section 6 — Termination actions

| P17 claim | OD control | Disposition | Required qualification |
|---|---|---|---|
| T1 handshakeFailure at `client.go:248`; intent proposed, not recorded. | OD-1, OD-8 | **Qualify** | Supervisor is the designated authority for process-lifecycle evidence; the client action must be represented and retained. Intent alone does not establish outcome. |
| T2–T4 protocol termination at client sites; intent proposed, not recorded. | OD-1, OD-8 | **Qualify** | Action identity and outcome are missing; `Kill` aliasing `Terminate` makes the action indistinguishable. |
| T5–T9 shutdown actions; intent/grace proposed, not recorded. | OD-1, OD-3, OD-8 | **Qualify** | Outcome and completion are missing; later evidence must be retained without overwriting a terminal. |
| T10 shutdown-grace-expiry at supervisor; action outcome. | OD-3, OD-7, OD-8 | **Qualify** | A timeout is not quarantine; the exact ownership/reaping condition and action outcome must be recorded. |
| `Kill` aliases `Terminate`; T2–T4 indistinguishable from T10 in evidence. | OD-8 | **Conflict / omission** | The action identity and intent cannot be reconstructed. This must be fixed before action causation can be evaluated. |

## 10. Section 7 — Precedence rules

| Rule | P17 projection | OD control | Disposition | Required qualification |
|---|---|---|---|---|
| 1 | Fatal handshake → `incompatible` for two helper codes; `unavailable` for handshake `protocol_violation`. | OD-6 | **Conflict** | The blanket `unavailable` path conflates missing mandatory capability with malformed/ambiguous cases. Apply OD-6's categories; missing mandatory capability → `incompatible`. |
| 2 | Fatal protocol → `crashed + protocol_violation`. | OD-5, OD-4 | **Conflict** | Fatal protocol does not establish process exit; `crashed` requires observed exit and no established orderly shutdown. |
| 3 | Fatal transport → `crashed + helper_crashed`. | OD-5 | **Conflict** | Same defect as rule 2; a transport failure is not an observed process exit. |
| 4a | Exit zero + completed shutdown → `absent`. | OD-4 | **Qualify** | `completed` must mean independently established orderly shutdown, not merely a `Terminate` return. |
| 4b | Exit zero + no action → `outcome_unknown`. | OD-4, OD-9 | **Conflict / qualify** | If the exit is confirmed and no orderly shutdown is established, OD-4's public fallback is `unavailable`. `outcome_unknown` must be internal-only or explicitly defined. |
| 4c | Nonzero + known action cause → internal `terminated_by_action`. | OD-4, OD-8 | **Aligned / Qualify** | OD-4 keeps it internal; “known action cause” must satisfy OD-8's A–E standard, not temporal association. |
| 4d | Nonzero + ambiguous cause → `outcome_unknown`. | OD-4, OD-8 | **Conflict / qualify** | A confirmed exit without orderly shutdown projects `unavailable` unless a prior justified failure applies; `outcome_unknown` is not a public state. |
| 5 | Diagnostics never change projection. | OD-1 | **Aligned** | The current O1/O2/O19 handling contradicts this rule. |

### 10.1 Combination table

| Combination | P17 projection | OD control | Disposition | Required qualification |
|---|---|---|---|---|
| Fatal handshake + any exit | `incompatible` or `unavailable` | OD-3, OD-6 | **Qualify** | “Any exit” does not define terminal publication order; OD-3 requires a latch. OD-6 requires classification by compatibility category. |
| Fatal protocol + any exit | `crashed + protocol_violation` | OD-4, OD-5 | **Conflict** | A fatal observation is not an observed process exit; orderly shutdown and prior-failure precedence must be applied. |
| Fatal transport + any exit | `crashed + helper_crashed` | OD-4, OD-5 | **Conflict** | Same as fatal protocol. |
| Exit zero, shutdown completed | `absent` | OD-4 | **Aligned** | Only if orderly shutdown is independently established. |
| Exit zero, no action | `outcome_unknown` | OD-4, OD-9 | **Conflict / qualify** | Confirmed exit with no orderly shutdown → `unavailable`; `outcome_unknown` needs explicit internal scope. |
| Exit nonzero, ambiguous | `outcome_unknown` | OD-4, OD-8 | **Conflict / qualify** | OD-4 fallback is `unavailable`; causation remains unresolved under OD-8. |
| Exit nonzero, action caused, no fatal | `terminated_by_action` (internal) | OD-4, OD-8 | **Aligned / Qualify** | No public state; A–E evidence is required for “action caused.” |
| Diagnostics only | unchanged | OD-1 | **Aligned** | Current O1/O2/O19 handling contradicts this. |

## 11. Section 8 — Evidence-loss findings

| Finding | Claim | OD control | Disposition | Required qualification |
|---|---|---|---|---|
| E1 | `failAll` keeps only the first read failure. | OD-1, OD-3 | **Conflict / Unsupported** | Violates preservation and late-evidence retention; projection-supporting authority cannot rely on a first-only failure. |
| E2 | `SetState` overwrites `lastFailure`. | OD-1, OD-3 | **Conflict / Unsupported** | A terminal or later evidence may be erased; a terminal latch and append-only record are required. |
| E3 | Stopping branch discards nonzero exit. | OD-4, OD-5 | **Conflict** | OD-4 requires `unavailable` unless a justified prior failure applies; a nonzero exit cannot simply become `absent`. |
| E4 | Termination intent is not a parameter; `Kill` aliases `Terminate`. | OD-1, OD-8 | **Conflict / omission** | Action identity, intent, and outcome cannot be established. |
| E5 | `waitForReader` returns nothing on timeout. | OD-1, OD-3 | **Conflict / omission** | Join versus timeout is lost; late evidence and terminal publication cannot be evaluated. |
| E6 | Stream closure is marked crashed but does not terminate for that code. | OD-5 | **Conflict** | A stream closure is not necessarily an observed process exit; no `crashed` projection is justified without that fact. |

P17's own statement that these losses must be fixed before precedence can be trusted is consistent with OD-1. The findings are not resolved merely by declaring the log authoritative.

## 12. Section 9 — Unresolved cases

| Case | P17 status | OD control | Disposition | Required qualification |
|---|---|---|---|---|
| U1 | Open | OD-1, OD-4, OD-8 | **Open** | Do not assign `absent` or `crashed` from the stop-intent branch; resolve stream-close provenance and action/exit evidence. |
| U2 | Open | OD-3, OD-4, OD-7, OD-8 | **Open** | Define fatal-versus-action arbitration; a timeout is not quarantine and not a crash. |
| U3 | Open | OD-3 | **Open** | Apply terminal latch and linearization; a later crash cannot overwrite the first terminal. |
| U4 | Open | OD-4, OD-5, OD-8 | **Open** | Fatal protocol plus zero exit does not equal `crashed`; apply orderly-shutdown and prior-failure rules. |
| U5 | Open | OD-4, OD-8 | **Open** | A nonzero exit after a deliberate action needs A–E causation; public projection remains governed by OD-4. |
| U6 | Open | OD-5 | **Open** | A live process cannot be `crashed`; the `client.go:756` write is non-conforming as written. |
| U7 | Open | OD-3, OD-7 | **Open** | Terminal immutability and internal quarantine/ownership must be defined; a timeout does not create a public state. |
| U8 | Open | OD-6 | **Open** | Missing mandatory capability is `incompatible`; malformed/ambiguous cases remain `unavailable`; unclassified categories stay unclassified. |

P17 correctly leaves these cases open. The problem is not the list; it is that §7's provisional projections often decide them by precedence.

## 13. Section 10 — Reviewer challenges

| Challenge | P17 request | OD control | Disposition | Required qualification |
|---|---|---|---|---|
| 1 | Verify every line and code claim. | OD-1, OD-9 | **Aligned** | Review request only; no test or implementation authorization. |
| 2 | Challenge each fatal assignment, especially O6, O12, O15, O16. | OD-5, OD-6 | **Aligned** | This is the correct challenge; O6/O12 and O18 conflict with OD-5, O15/O16 need OD-6 classification. |
| 3 | Test U1–U8, including shutdown timeout plus protocol violation and handshake incompatibility plus crash. | OD-9 | **Future / not authorized** | Testing may be needed later, but OD-9 and the register do not authorize tests now; this matrix does not run them. |
| 4 | Check precedence never deletes evidence and no projection is unreachable. | OD-1, OD-3, OD-7 | **Qualify** | Current E1/E2/E5 and unreachable `quarantined` fail this check; P17's log is not implemented. |
| 5 | Confirm five projections (`absent`, `crashed`, `incompatible`, `unknown`, `quarantined`) are reachable or declare them not. | OD-4, OD-5, OD-6, OD-7 | **Conflict** | `quarantined` is internal-only under OD-7; `unknown` is not a public provider state; `crashed` must be exit-based. The challenge's premise is not compatible with the accepted decisions. |
| 6 | Return finding-level dispositions; do not accept the design on the author's behalf. | OD-9 | **Aligned** | This matrix is finding-level and does not approve P17. |

## 14. Section 11 — Status

| P17 claim | OD control | Disposition | Required qualification |
|---|---|---|---|
| Frozen proposal; no contract edit; no rev4; no implementation; no tests; C1 unauthorized. | OD-9 | **Aligned** | Correct status. P17 remains unapproved and this matrix does not change it. |

## 15. Direct conflict register

The following are the highest-priority P17 corrections. They are not implementation tasks and do not authorize a revision.

| # | Conflict | P17 location | Accepted control |
|---|---|---|---|
| C-01 | Fatal protocol/transport → `crashed` without observed exit | §4 S6; §5 O6/O7/O8/O12/O18; §7 rules 2–3 and table | OD-5 |
| C-02 | Nonzero exit while stopping → `absent` | §5 O5; §7 table | OD-4 |
| C-03 | Confirmed exit, no orderly shutdown → `outcome_unknown` instead of `unavailable` | §7 rules 4b/4d; table | OD-4 |
| C-04 | Handshake `protocol_violation` → `unavailable` for missing mandatory capability | §5 O16; §7 rule 1 | OD-6 |
| C-05 | `quarantined` listed as a required public projection | §10 challenge 5 | OD-7 |
| C-06 | Diagnostics select projections | §4 S4 vs O1/O2/O19 | OD-1/OD-4 |
| C-07 | No generation binding in records | §2 Observation/Action | OD-1/OD-2 |
| C-08 | No terminal latch or late-evidence rule | §2, §7 | OD-3 |
| C-09 | “Known action cause” lacks A–E test | §7 rule 4c | OD-8 |
| C-10 | `seq` has no allocator or event/observation distinction | §2, §7 | OD-1/OD-8 |

## 16. Conditions before a P17 revision

A future P17 revision should address these conditions in order. This matrix does not perform the revision.

1. **Authority:** state that the supervisor is the designated authority for process-lifecycle evidence and provider-state publication; require client/session observations to be retained.
2. **Generation:** add the OD-2 generation/session/PID binding to every observation, action, and exit fact.
3. **Ordering:** define `seq` allocation and distinguish event order from observation/processing order.
4. **Terminal:** add OD-3's immutable first-terminal latch and late-evidence retention rule.
5. **Projection:** replace fatal-wins rules with OD-4/OD-5-compliant exit and orderly-shutdown tests; never publish `crashed` without an observed exit.
6. **Compatibility:** replace O15/O16's generic split with OD-6's categories.
7. **Quarantine:** remove `quarantined` from the public projection challenge; make it internal under OD-7.
8. **Causation:** replace “known action cause” with OD-8's A–E standard; keep `terminated_by_action` internal.
9. **Unknown:** define whether `outcome_unknown` is internal only; if public, it requires a separate contract decision. OD-4 currently supplies `unavailable` for confirmed exit without orderly shutdown.
10. **Evidence loss:** map every E1–E6 loss to a retention or ordering requirement; do not claim precedence can be trusted while the losses remain.
11. **U cases:** resolve each U case or apply an explicit unknown rule; do not let §7 decide them by precedence.
12. **No new state:** preserve the frozen provider-state set; no contract amendment is implied by this matrix.

## 17. Conclusion

P17's general direction is compatible with the accepted evidence model, but its projection table is not. The largest defects are:

- `crashed` is derived from fatal protocol/transport without an observed exit;
- `absent` is used for a nonzero exit while stopping;
- `outcome_unknown` is used where OD-4 requires `unavailable`;
- `quarantined` remains a claimed public projection despite OD-7;
- the handshake split does not apply OD-6;
- diagnostics currently select provider states;
- the record lacks generation, terminal, retention, and allocator semantics.

P17 remains **unapproved**, U1–U8 remain **OPEN**, contract rev3 remains **frozen**, and no implementation is authorized. The next step is a revision proposal after the authority and model gaps are resolved—not a source change now.

## 18. Verification and limitations

This is a read-only reconciliation. No existing file was modified. The register remains at `3c64797b…`, 434 lines; P17 remains `01d887c1…`; contract rev3 remains `ed0da1a7…`; no source, test, configuration, or generated file was changed.

**Limitations:** this matrix relies on P17's own inventory and the existing source assessment for current-handling claims. It did not re-run the source audit, execute U1–U8, observe interleavings, or verify platform behavior. It does not decide the U cases, approve P17, amend contract rev3, or authorize implementation. A P17 revision remains separate work and is not performed here.

**No implementation was performed or authorized.**
