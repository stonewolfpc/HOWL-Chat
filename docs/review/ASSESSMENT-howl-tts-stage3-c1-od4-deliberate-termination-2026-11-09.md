# C1 / P17 — OD-4 Deliberate-Termination Projection Assessment (Read-Only)

**Date:** 2026-11-09

**Status:** Assessment only. **OD-4 remains UNRESOLVED.** No register update, disposition, contract amendment, P17 edit, source change, or implementation is made or recommended by this document.

**Baseline:**

- Decision register: `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md`, SHA-256 `7aef96a0…`, 371 lines.
- Six decisions are accepted: OD-2, OD-3, OD-5, OD-6, OD-7, OD-8. OD-1 is deferred; OD-4 and OD-9 are unresolved.
- P17 remains unapproved. Contract rev3 is a frozen proposal, not an accepted amendment.
- OD-9 remains untouched.

## 1. Exact OD-4 question and current status

The register asks:

> **When a known deliberate action explains a nonzero exit, what public provider state, if any, should be published?**

The recorded status is **UNRESOLVED**. The alternatives are:

- **A. Publish `absent` (current).** Treat the action as explaining the exit.
- **B. Keep `terminated_by_action` internal; publish no new provider state.** Preserve the last operationally justified state.
- **C. Publish `outcome_unknown` as a public state.** Requires a contract change; `unknown` is not in the frozen set.
- **D. Publish a new public state.** Requires a separately approved amendment.

The non-binding default is B. This assessment tests B against the accepted decisions and the source.

## 2. Source inventory

| Location | Source behavior | Relevance |
|---|---|---|
| `supervisor.go:166–195` | `observeExit` waits for the process, closes it, and classifies the result. | The only path that observes an actual exit. |
| `supervisor.go:176–182` | If `providerState == stopping` at observation, a **zero or nonzero** exit becomes provider `absent`; process state becomes `absent` for `err == nil`, otherwise `exited`. | The current public projection for a nonzero exit observed while stopping is `absent`, not `crashed`. |
| `supervisor.go:183–190` | Otherwise the process state becomes `exited`, provider state becomes `crashed`, and `helper_crashed` is recorded. | A crash is selected from the observer's view of pre-exit provider state, not from the action identity. |
| `supervisor.go:246–288` | `Terminate` sets process/provider state to `stopping`, closes the process, waits, kills after grace, and can return `timeout` after a second grace. | The action writes `stopping`; it does not record intent, completion, consequence, or a causal link. |
| `supervisor.go:291–293` | `Kill` delegates to `Terminate` with the default grace. | Kill and shutdown are not distinguishable in the exit classification. |
| `supervisor.go:233–242` | `SetState` can overwrite provider state; when `proc != nil` and the new state is not `stopping`, it sets process state `alive`. | Provider state is mutable and not a terminal record. |
| `supervisor.go:306–343` | `Status` copies both dimensions; `ProcessAlive` is derived. | The public status surface has no `terminated_by_action` field. |
| `client.go:645–674` | `Shutdown` sends the shutdown request, calls `Terminate`, then clears `closing`, `started`, and `transport` regardless of the termination result. | A shutdown action is not recorded as a causally linked exit, and a timeout can leave a retained process. |
| `client.go:859–864` | `cleanupRejected` calls `Terminate` after a rejected synthesis when the process is alive. | Cleanup termination is another deliberate action whose exit can be classified as `absent`. |
| `client.go:248`, `:654`, `:661`, `:664`, `:861` | Termination results are often ignored or only returned. | No action-result evidence reaches the projection path. |
| `types/provider.go:8–16` | Process states: `absent`, `starting`, `alive`, `stopping`, `exited`. | No `terminated_by_action` or `outcome_unknown` process state. |
| `types/provider.go:18–36` | Provider states are the frozen thirteen; no `terminated_by_action`. | There is no public field for the action-causation distinction. |

### 2.1 Source-verified behavior

1. **The current source maps a nonzero exit observed while `stopping` to `absent`.** The mapping is based on provider state at observation time, not on an action record or an orderly-shutdown finding.
2. **`stopping` is set by the action.** `Terminate`/`Kill` set it before waiting; the observer therefore cannot distinguish a shutdown, a cleanup termination, a protocol-triggered kill, or another caller's termination.
3. **The source does not record action intent or completion in the exit path.** `Terminate` returns an error to its caller, but the supervisor's exit observer does not consult an action record, action ID, or completion result.
4. **`Terminate` can return `timeout` while a process remains owned.** That path is an OD-7 internal ownership question, not a public action projection.
5. **`observeExit` can publish `crashed` when the observer does not see `stopping`.** The same process exit can therefore become `absent` or `crashed` depending on interleaving, even when a deliberate action exists.
6. **The source does not establish orderly shutdown.** Neither `stopping`, a successful `Terminate` return, nor a zero/nonzero exit code establishes that the shutdown was orderly.

## 3. Accepted decisions that constrain OD-4

### 3.1 OD-5 — `crashed`

OD-5 accepts:

> `crashed` describes an observed process exit for which orderly shutdown is not established.

OD-5 does not assert a physical mechanism, cause, intent, or action attribution. It also does not make `crashed` a synonym for an unexpected exit or a provider/runtime failure.

Consequences for OD-4:

- A nonzero exit after a known action is not automatically `crashed`.
- A zero exit after a known action is not automatically `absent`.
- `crashed` may be justified only when the observed-exit/no-orderly-shutdown condition is independently established; the action label does not establish it.
- OD-4 must not rewrite OD-5's meaning.

### 3.2 OD-3 — terminal decision

OD-3 accepts an immutable first terminal decision. It does not make the provider state immutable, and it does not define the public projection for an action-caused exit. OD-4 must not use a provider-state write to overwrite a terminal decision, and it must not treat `absent` as a terminal decision by implication.

### 3.3 OD-8 — deliberate-action causation

OD-8 requires A–E evidence for causal attribution: process-instance identity, action identity and outcome, event ordering, mechanism operation, and competing-explanation assessment.

A known action or a temporal sequence is **association**, not causation. `terminated_by_action` is an action-causation claim; it cannot be established merely because `stopping` preceded the exit, because the action completed, or because the exit was nonzero.

### 3.4 OD-7 — quarantine

OD-7 accepts internal-only quarantine for an explicitly recorded unresolved ownership/reaping condition. A termination timeout is not automatically quarantine, and quarantine is not a public provider state. OD-4 must not use `quarantined` as the public outcome for an action-caused nonzero exit.

### 3.5 Lifecycle decision

The lifecycle decision (`DECISION…`, §1.3–1.4, §3–4) is directly relevant:

- `unknown` and `quarantined` are not provider states.
- Silent `unknown → absent` is forbidden.
- Until public exposure is separately approved, no terminal provider state is written from a disposition; the provider state remains at its last operationally justified value.

This supports internal `terminated_by_action` and argues against publishing `absent` merely because the process was observed while `stopping`.

## 4. What P17 and contract rev3 propose

### 4.1 P17 §7

P17's proposed precedence table says:

- zero exit plus completed shutdown → `absent`;
- zero exit with no action → `outcome_unknown`;
- nonzero plus known action cause → internal `terminated_by_action`, public projection undecided;
- nonzero plus ambiguous cause → `outcome_unknown`.

P17 is unapproved. Its model is useful because it separates the action-caused case from the ambiguous case, but:

- it assumes a completed shutdown action and a known action cause without defining the evidence for either;
- it treats `outcome_unknown` as a projection while the frozen provider-state set does not contain it;
- it does not specify what public state should accompany internal `terminated_by_action`.

### 4.2 Contract rev3 A6

Contract rev3 A6 proposes:

- nonzero exit with no `failureObserved` and known causality → `crashed`;
- zero exit with no failure and known causality → `absent`;
- any fatal failure → `crashed` with kind;
- ambiguous cases publish conservatively, never `absent` without evidence.

Contract rev3 is a frozen proposal, not an accepted amendment. Its A6 table cannot be adopted as authority by OD-4. In particular, its “nonzero → crashed” rule is not the same as OD-5's process-exit/no-orderly-shutdown definition, and its “known causality” field depends on evidence the source does not retain.

## 5. The three facts OD-4 must not conflate

The following are distinct:

1. **An action was requested or attempted.** This establishes action identity, not execution or effect.
2. **The action or an exit was observed.** `Terminate` returning, `stopping` being set, or a nonzero status appearing does not establish which event caused which.
3. **A nonzero exit after an action-adjacent event.** This is a sequence. Without OD-8's A–E evidence, it is not a causal finding and does not establish an orderly shutdown.

The current source collapses these into the `stopping` branch and publishes `absent`. That is an intent-based classification, not a demonstrated clean outcome.

## 5.1 Last justified state versus current availability

A confirmed process exit creates a distinct problem for B: the last provider state may be historically justified but no longer operationally true.

The `ready` case is the clearest:

- `ready` asserts that the provider is operational and usable.
- A confirmed process exit means the provider is no longer operational.
- Retaining `ready` after the exit would preserve historical information while making a false current-availability claim.

The same problem applies to other active states that imply a live or usable provider: `busy`, `degraded`, `starting`, `handshaking`, `loading`, `stopping`, and `backoff`. These states may describe the pre-exit condition, but they must not be retained as the current public state after the process is confirmed gone.

A prior failure state is different. `incompatible`, `unavailable`, and `crashed` describe a provider condition or process outcome, not operational readiness. If one remains independently justified, retaining it does not assert that the provider is currently usable. The evidence for that condition must still be retained internally; it must not be silently overwritten by the termination action.

The frozen enum may not be able to express both “the last failure was X” and “the provider is currently not available” at once. OD-4 must choose a current public state and preserve the other information internally, rather than retaining an active state as a substitute for current availability.

This refines B: retain the failure evidence and the action classification internally; do not retain `ready` or another active state as a current public claim.

### 5.2 Answers to the four questions

| Question | Answer |
|---|---|
| Does retaining `ready` after a confirmed exit create a stale operational claim? | **Yes.** `ready` asserts current usability. Once the process is confirmed gone, retaining it is a false current-availability claim, even if it was historically justified. |
| Is `unavailable` appropriate when the exit is confirmed but orderly shutdown is not? | **Yes, as the preferred fallback.** It is an existing state that can truthfully describe current non-availability without asserting an orderly shutdown. It must not be treated as a clean-outcome claim. Its post-exit use requires explicit contract wording; it is not a new enum value. |
| Can `absent` be reserved for independently established orderly shutdown? | **Yes.** `absent` should not be inferred from `stopping`, a nonzero exit, or an action request. Reserving it for established orderly shutdown or another established absence narrows its terminal use; it does not add a new state. |
| Do these choices introduce a new terminal meaning into the frozen enum? | **No new state is introduced.** `unavailable` and `absent` remain the same thirteen values. Their post-exit/terminal use is a clarification that must be explicit in contract text; a new terminal meaning must not be inferred from the state name alone. |

This resolves the stale-state distinction: the public state describes current availability, while the prior failure and the action classification remain internal evidence. Exact orderly-shutdown evidence and the contract wording for `unavailable`/`absent` remain open.

## 6. Alternatives assessed

| Alternative | Source fit | Public consequence | Assessment |
|---|---|---|---|
| **A. Publish `absent` (current)** | Matches `observeExit` for an exit observed while `stopping`. | Asserts a clean/absent outcome for both zero and nonzero exits, even when the action was a kill, cleanup termination, or protocol-triggered termination. | Rejected. It treats `stopping` as sufficient and can erase a prior `incompatible`, `unavailable`, or `crashed` condition. |
| **B. Keep `terminated_by_action` internal; no new provider state; preserve the last operationally justified state** | Does not require a new enum value; aligns with the lifecycle decision and the absence of action evidence in the exit path. | Must not retain `ready` or another active state after a confirmed exit; preserve failure evidence internally and use a current-availability projection. | **Recommended, refined by §5.1.** |
| **C. Publish `outcome_unknown` as a public state** | `outcome_unknown` is an internal/projection concept, not one of the frozen provider states. | Requires a contract change; risks turning uncertainty into a provider condition. | Rejected for now. It conflicts with the frozen set and the lifecycle decision's separation. |
| **D. Add a new public state** | No source value; requires a separately approved amendment. | Largest contract change; no demonstrated necessity. | Rejected. The internal distinction is sufficient unless a future exposure decision justifies a new state. |

### 6.1 Why B is the narrowest defensible decision

B preserves the distinction without adding a public state:

- `terminated_by_action` remains internal action-causation evidence;
- a nonzero exit does not automatically become `absent` or `crashed`;
- prior failure evidence is retained internally, but active states such as `ready` are not retained as current public claims;
- a confirmed exit without orderly shutdown is projected as `unavailable` unless an independently justified failure state remains applicable;
- an established orderly shutdown can still support `absent` through a separate rule;
- a process-exit fact can still support `crashed` under OD-5 when its independent condition is met.

B does not claim that the prior state is fresh or that the exit is orderly. The refinement accepts that the public state cannot carry both historical failure and current availability; it chooses current availability and preserves the failure internally. That is preferable to a false clean outcome or a stale `ready` claim.

## 7. Recommended OD-4 decision formulation

The following is a recommendation for owner consideration, not a recorded decision.

> **OD-4 — Deliberate-termination projection (proposed):**
>
> 1. `terminated_by_action` remains an internal action-causation classification. It is not a provider state and has no public projection.
> 2. No new public provider state is added. The thirteen-value provider-state set remains frozen.
> 3. A nonzero exit after a known action does **not** by itself justify `absent`. `absent` is reserved for an independently established orderly shutdown or another established absence, not for `stopping`, a stop request, a successful `Terminate` return, or a nonzero exit.
> 4. A nonzero exit after a known action does **not** by itself justify `crashed`. `crashed` remains governed by OD-5's observed-exit/no-orderly-shutdown meaning. The action label does not create a crash.
> 5. After a confirmed process exit, active provider states — including `ready`, `busy`, `degraded`, `starting`, `handshaking`, `loading`, `stopping`, and `backoff` — must not be retained as the current public state. `ready` is specifically false once the provider is no longer operational.
> 6. If an independently established orderly shutdown exists, `absent` may be published. If the process exit is confirmed but orderly shutdown is not established, `unavailable` is the preferred existing current-availability projection, unless a prior failure state (`incompatible`, `unavailable`, or `crashed`) remains independently justified and does not assert operational readiness.
> 7. A prior failure state is not silently erased. If it cannot remain the current public state, its evidence is retained internally alongside the action classification; `unavailable` must not be used to discard it.
> 8. `outcome_unknown` remains an internal/epistemic distinction unless a separate approved amendment gives it a public projection. It is not added to the provider-state enum by OD-4.
> 9. The post-exit use of `unavailable` and the orderly-shutdown restriction on `absent` must be explicit in contract text; they are clarifications of existing states, not new enum values. OD-4 does not infer a new terminal meaning from the state names.
> 10. A public action-cause field, diagnostic exposure, or any new public state requires a separately approved contract amendment.
> 11. OD-4 does not decide the exact fallback when no prior provider state remains operationally justified; that is a separate projection rule.
> 12. OD-4 does not decide U2, U4, U5, U7, or U3, and does not reopen OD-3, OD-5, OD-7, or OD-8.

### 7.1 What it would make non-conforming

If accepted, the current `observeExit` mapping of a nonzero exit observed while `stopping` to provider `absent` would be non-conforming as a general rule. It would not establish an orderly shutdown merely because the action set `stopping`. A future implementation would need a separate, evidence-based orderly-shutdown rule before publishing `absent`, or use `unavailable` as the current-availability fallback.

Retaining `ready` after a confirmed exit would also be non-conforming under the refinement. The source's `crashed` branch is not automatically conforming either: it must meet OD-5's independent condition. No source change is proposed by this assessment.

### 7.2 What it would not decide

- The exact evidence that establishes an orderly shutdown for `absent`.
- Whether a completed shutdown action can ever support `absent` without a separate orderly-exit finding.
- The fallback when no prior provider state remains justified.
- Whether `outcome_unknown` should later become a public projection.
- U2's shutdown-timeout/protocol-violation winner.
- U4's zero-exit-after-fatal-observation case.
- U5's action-caused nonzero-exit classification.
- U7's terminal immutability across a later ownership condition.
- Any implementation or test.

## 8. Consequences

### 8.1 U4

U4 is a zero exit after a fatal protocol violation. OD-4-B does not decide whether the fatal observation or the exit controls; it requires the prior justified state to survive the action unless new evidence displaces it. U4 remains OPEN.

### 8.2 U5

U5 is a nonzero exit after a deliberate kill. OD-4-B keeps `terminated_by_action` internal, refuses automatic `absent`, and refuses automatic `crashed`. U5 remains OPEN until the orderly-shutdown and action-causation evidence rules are settled.

### 8.3 U2

U2 is a shutdown timeout racing a protocol violation. OD-4-B does not choose a public winner and does not turn an action timeout into a state. U2 remains OPEN.

### 8.4 U7

U7 is a `Terminate` timeout after a terminal projection is published. OD-4-B does not make the provider state immutable or decide whether the first terminal survives. U7 remains OPEN.

### 8.5 U1 and U3

U1's orderly EOF during shutdown still needs stream-close provenance and an action-completion rule. U3's handshake-incompatible-then-crash sequence remains governed by OD-3 and is untouched.

### 8.6 OD-6 and OD-7

A prior `incompatible` state is not silently erased by a later cleanup termination merely because `stopping` was set. Under the refinement, it may remain the current public state only if the incompatibility remains independently justified and the state does not assert operational readiness; otherwise `unavailable` becomes the current-availability projection and the compatibility evidence remains internal. OD-6's public rule is unchanged.

OD-7's internal quarantine condition is separate. OD-4-B does not turn a retained process after a timeout into a public state and does not reopen OD-7.

### 8.7 OD-9 and P17

OD-9 remains untouched. P17 remains unapproved; its internal `terminated_by_action` row is consistent with B only if its public projection remains undecided and its causal assumptions are not treated as established. Contract rev3 A6 is not adopted.

## 9. Remaining uncertainty

1. **Orderly-shutdown evidence.** The source does not define what establishes an orderly shutdown. `stopping`, a completed `Terminate`, and an exit code are insufficient by themselves.
2. **Fallback wording.** The refinement selects `unavailable` as the preferred current-availability fallback when exit is confirmed and orderly shutdown is not, but its contract wording and its interaction with an independently justified `incompatible`/`crashed` state remain open.
3. **Single-state limitation.** The public state cannot express both the last failure and current operational availability at once. The refinement chooses current availability and preserves the failure internally; whether that is sufficient for consumers is a contract-exposure question, not a new enum question.
4. **Action causation.** A known action is not a causal finding. OD-8's A–E standard still applies; without it, `terminated_by_action` itself may be unjustified.
5. **`crashed` boundary.** A deliberate exit without orderly shutdown may meet OD-5's process-exit condition, but the action label does not make it a crash. The interaction remains U5.
6. **`absent` boundary.** A zero exit with a completed shutdown is proposed by P17 as `absent`; the exact evidence for “completed” and “orderly” remains open.
7. **`outcome_unknown` exposure.** Whether it ever becomes a public projection is a separate amendment question.
8. **P17 and contract rev3.** P17's `terminated_by_action` row and contract A6's nonzero-exit rules are proposals, not authority. No edit is made here.

## 10. What this assessment does not decide

This assessment does not:

- record an OD-4 disposition, acceptance, rejection, or deferral;
- change OD-1–OD-9 or U1–U8;
- close U2, U4, U5, U7, or U3;
- add, rename, or expose a provider state;
- define the exact orderly-shutdown, fallback, or action-causation evidence;
- edit P17, the decision register, contract rev3, or the lifecycle decision;
- authorize implementation, source changes, tests, builds, or runtime experiments.

## 11. Verification note

This is a read-only assessment. No existing artifact was modified. The decision register remains at SHA-256 `7aef96a0…`, 371 lines; OD-7 remains `e878f561…`; OD-6 remains `cd9d5c65…`; OD-5 remains `47d19bf4…`; P17 remains `01d887c1…`; contract rev3 remains `ed0da1a7…`; the lifecycle decision remains `aa136bc9…`. No source, test, configuration, or generated file was modified; no tests, builds, benchmarks, or runtime experiments were run.

**OD-4 remains UNRESOLVED.**
