# C1 / P17 — OD-5 `crashed` Semantics Assessment (Read-Only)

**Date:** 2026-11-09

**Status:** Assessment only. **OD-5 remains UNRESOLVED.** No register update, disposition, contract amendment, source change, or implementation is made or recommended by this document.

**Baseline:**

- Decision register: `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md`, SHA-256 `50350c44…`, 309 lines; OD-5 entry at `:155–167`.
- P17: `01d887c1…`; still unapproved.
- Corrected r3 response: `60ca1cb5…`.
- Lifecycle decision: `aa136bc9…`.
- Contract rev3: `ed0da1a7…`; the document itself says it is a proposed revision for independent review, not an accepted amendment (`:1–5`).

The thirteen-value provider-state set remains frozen. `crashed` is one of those values, but its meaning is not settled by the set.

## 1. Exact OD-5 question and current status

The register asks:

> **Does `crashed` describe an observed process termination, a fatal helper failure that may leave the process alive, or a narrower provider condition?**

The recorded status is **UNRESOLVED**. The register's alternatives are:

- **A. Process-outcome only.** `crashed` requires an observed exit or equivalent process termination; a live write failure must use another state.
- **B. Provider-condition failure.** `crashed` may describe a fatal helper condition even if the process is alive.
- **C. Split.** `crashed` means an observed process exit; live transport failure uses `degraded`/`unavailable`, or internal `outcome_unknown`.
- **D. Remove `crashed` from the write-failure path only.** Keep the broader state but declare `client.go:756` non-conforming.

The non-binding proposed default is C. This assessment does not record that default as a decision.

### 1.1 Four distinct facts and claims

OD-5 fails if any two of these are treated as interchangeable:

| Layer | What it is | What establishes it | What it does not establish |
|---|---|---|---|
| **Observed process exit** | A process-lifecycle fact: the helper process has exited. | The exit observation itself: `Wait()` returned and the supervisor's process observation completed (`supervisor.go:168–190`). | It does not establish an orderly shutdown, a provider/runtime failure, a physical crash mechanism, or that an action caused the exit. A zero or nonzero status is compatible with this fact. |
| **Orderly shutdown** | A separately established lifecycle outcome. | Separate evidence under the applicable lifecycle rules; not a stop request, `stopping`, a termination API result, temporal order, or exit code alone. | It does not follow from the process being asked to stop, from `providerState=stopping`, or from `Terminate` returning. |
| **Deliberate-action causation** | A causal claim that a deliberate action caused the exit. | The accepted OD-8 standard: A-E evidence is required for causal attribution; A-O is association only. | It does not follow from `stopping`, a stop request, observation order, a successful `Terminate`/`Kill` result, or a nonzero exit code. |
| **Provider/runtime failure** | An operational condition of the provider or helper session. | Failure observations, transport errors, provider state, and related evidence. | It does not establish that the process exited. A live process may have a failed provider/runtime condition. |

`crashed` is a classification of the first layer. It is not a synonym for orderly shutdown, deliberate-action causation, or provider/runtime failure.

## 2. Source inventory

| Location | Source behavior | Relevance to OD-5 |
|---|---|---|
| `internal/audio/tts/helper/supervisor.go:166–195` | `observeExit` waits for the process, closes it, and classifies the result. | The only path that observes an actual exit. |
| `supervisor.go:176–182` | If `providerState == stopping` at observation time, a zero or nonzero exit becomes provider `absent`; process state is `absent` or `exited`. | The source classifies a stop-intent exit as `absent`; this is an intent-based classification, not proof of an orderly shutdown. |
| `supervisor.go:183–190` | Otherwise, the process state becomes `exited`, the provider state becomes `crashed`, and a `helper_crashed` failure is recorded. | One source-verified `crashed` writer. This happens after `Wait` returns. |
| `supervisor.go:233–242` | `SetState` writes provider state and failure; when `proc != nil` and the new state is not `stopping`, it sets `processState = alive`. | A provider-state write can leave the process dimension `alive`, including a `crashed` write. |
| `supervisor.go:306–333` | `Status` copies process and provider state and computes `ProcessAlive` at `:324` using both dimensions. | The status surface contains both dimensions and a derived liveness boolean. |
| `supervisor.go:335–343` | `Supervisor.ProcessAlive()` returns true when `proc != nil` and `processState != exited`. | The supervisor's own liveness method is process-handle based. |
| `supervisor.go:351–356` | `processAlive(process, provider)` returns true for process states `alive`, `starting`, or `stopping` only when provider state is neither `crashed` nor `unavailable`. | The status boolean vetoes liveness from provider state, unlike `Supervisor.ProcessAlive()`. |
| `internal/audio/tts/helper/client.go:755–758` | On a synthesis `transport.Send` error, `Synthesize` writes `TTSStateCrashed` with `helper_crashed`, then calls `cleanupRejected(..., true)`. | The second source-verified `crashed` writer; it can run while the process is alive. |
| `client.go:859–864` | `cleanupRejected` checks `sup.ProcessAlive()` and calls `Terminate` when it is true (or output exists). | The cleanup path uses the process-handle liveness method, not `Status.ProcessAlive`. |
| `client.go:200–208` | `Start` checks `sup.ProcessAlive()`; a live owned process returns without restart. | A process that is alive but provider-`crashed` is treated as still started by this path. |
| `client.go:295–308` | `readLoop` maps a generic stream error to `helper_crashed`, calls `failAll`, and calls `Kill` only for `invalid_frame`, `frame_too_large`, or `protocol_violation`. | A generic stream closure does not itself write `crashed`, and may not terminate. |
| `client.go:369–386` | `failAll` retains the first read failure and delivers it to pending requests. | It does not write a provider state. |
| `client.go:790–797` | A pending synthesis failure publishes `degraded` (and a timeout/cancel may publish `ready`). | A fatal protocol observation delivered through `failAll` can emerge as `degraded`, not `crashed`. |
| `internal/audio/types/provider.go:8–16` | Five process states: `absent`, `starting`, `alive`, `stopping`, `exited`. | `crashed` is not a process state. |
| `types/provider.go:18–36` | Thirteen provider states, including `crashed`. | The contract vocabulary is frozen; semantics remain open. |
| `types/provider.go:57–61` | `TTSProviderStatus` carries both `ProcessState` and `ProcessAlive` plus aggregate `State`. | The status shape exposes the distinction the implementation then blurs. |

### 2.1 Frozen/proposed evidence

| Artifact | What it says relevant to OD-5 | Authority status |
|---|---|---|
| P17 `:102–125` | Fatal protocol or transport observation projects `crashed`; exit facts without a fatal observation use `absent`/`outcome_unknown`/internal `terminated_by_action`. | **Not approved.** A proposal, not accepted contract behavior. |
| P17 `:149` | U6: provider `crashed` while the process is alive; provider state conflates condition with process outcome. | **OPEN.** Consistent with the source finding. |
| Contract rev3 `:130–148` | A6 proposes `crashed` for fatal `failureObserved` even when `exitStatus == 0`, and for nonzero exits without a fatal observation. | The document says **proposed revision for independent review, not an accepted amendment** (`:1–5`). Not governing authority. |
| Lifecycle decision `:11–14`, `:31–38` | Provider state and process-lifecycle disposition are separate dimensions; `unknown → crashed` is forbidden. | **Accepted owner decision**, but not a contract amendment and not a definition of `crashed`. |
| Amendment draft `:1–5`, `:75`, `:99` | A2-D4 proposes `ProcessAlive` solely from process state and says only an unexpected exit is `crashed` (`:75`; repeated in the supersession map at `:99`). The document is **not adopted** and on conditional hold. | Prior proposal only; not authority. |
| Prior assessment `:186–205`, `:315–333` | R13/U6 identify the `Synthesize` write and the `Status` liveness conflation. | Prior reported findings; source claims are rechecked below. |

## 3. Source-verified behavior

### 3.1 There are two `crashed` writers, not one

A source search finds `TTSStateCrashed` writes at:

1. `supervisor.go:183–190` — an unexpected exit observed by `observeExit`.
2. `client.go:755–758` — a synthesis transport write failure, before `cleanupRejected`.

The first is a process-exit path: `Wait()` has returned, and `processState` is set to `exited`.

The second is not: `SetState` can set `processState = alive` at `supervisor.go:239–240` because the process handle still exists and the new provider state is not `stopping`. The provider state is `crashed` while the process dimension may remain `alive`.

There is no third source writer for `crashed`. In particular, `readLoop` does **not** call `SetState(crashed)` on a generic stream closure; it calls `failAll`, and only kills for certain protocol/frame codes (`client.go:295–308`).

### 3.2 R13 trace: `crashed` before an observed exit

The source-verified sequence for a synthesis write failure is:

1. The helper process exists and is alive; the provider is in a usable state.
2. `transport.Send` returns an error (`client.go:755`).
3. `Synthesize` calls `SetState(TTSStateCrashed, "", helper_crashed)` (`:756`). `SetState` sets provider state to `crashed` and, because `s.proc != nil` and the state is not `stopping`, sets `processState = alive` (`supervisor.go:239–240`).
4. `Synthesize` then calls `cleanupRejected(..., true)` (`client.go:757`). `cleanupRejected` calls `sup.ProcessAlive()` (`:860`). That method returns true while `proc != nil` and `processState != exited` (`supervisor.go:335–343`), so cleanup calls `Terminate` (`client.go:861`).
5. `Terminate` sets provider state to `stopping` (`supervisor.go:256–259`), closes stdin, and waits. When the exit observer sees `stopping`, it publishes `absent`, not `crashed` (`supervisor.go:176–182`). If `Terminate` times out, the provider remains `stopping`; the process may remain alive, but the provider is no longer `crashed`.

Therefore the source can publish `crashed` before the exit is observed, but in this path the publication is typically transient: it is followed by `stopping` and then `absent`. The exact contradiction is observable in the window between the `SetState(crashed)` write and `Terminate` acquiring the state lock.

**Correction to a prior inference:** the prior U6 text says a `Terminate` timeout can leave the process alive while the provider is `crashed`. That is not supported by this path. `Terminate` sets `stopping` before waiting, so a timeout leaves `stopping`, not `crashed`. The live/`crashed` contradiction is a pre-`Terminate` observation window, not a persistent post-timeout state. This is a correction of a prior inference, not a source change.

**Scope of that correction:** this is a source-path finding, not runtime verification. U6 remains **OPEN**.

The `stopping` branch is an intent-based classification, not proof of an orderly shutdown. When `providerState == stopping` at observation time, `observeExit` publishes provider `absent` for either a zero or nonzero exit (`supervisor.go:176–182`); it does not establish that the shutdown was orderly, and it does not treat the nonzero exit as a crash. That classification boundary is part of U5/OD-4 and is not decided here.

### 3.3 The two liveness answers disagree

For a live process with provider state `crashed`:

- `Supervisor.ProcessAlive()` returns **true** (`proc != nil`, `processState != exited`; `supervisor.go:335–343`).
- `Status()` copies `ProcessState = alive` and `State = crashed`, then computes `ProcessAlive = processAlive(alive, crashed)`; `processAlive` returns **false** because provider state is `crashed` (`supervisor.go:324`, `351–356`).

So the status object can simultaneously contain:

```text
ProcessState = "alive"
State        = "crashed"
ProcessAlive = false
```

while the supervisor's own `ProcessAlive()` method returns true. The process may actually be alive; the status boolean says it is not. This is the concrete R13 conflation.

The precise reading is: the `Status` surface can report `ProcessState=alive` while `ProcessAlive=false` and `State=crashed`. The separate `Supervisor.ProcessAlive()` method reports true. The source has no single consistent answer.

### 3.4 The status conflict has behavioral consequences

- `cleanupRejected` uses `sup.ProcessAlive()`, so it does attempt to terminate a live process even when `Status.ProcessAlive` would be false (`client.go:859–864`).
- `Start` uses `sup.ProcessAlive()` and returns without restart when the process is still owned (`client.go:200–208`). A live process with provider `crashed` is therefore treated as still started by this path.
- `Status.ProcessAlive=false` does not mean the process is gone. A consumer using the status boolean can reach the opposite conclusion from the supervisor's cleanup path.

This is a semantic and observational defect, not merely a display detail.

### 3.5 Fatal observations are not source-verified as `crashed`

P17 §7 and contract rev3 A6 propose `crashed` for fatal protocol/transport observations. The current source does not implement that uniformly:

- `readLoop` maps a generic stream error to `helper_crashed` but does not write provider state; only certain frame/protocol codes trigger `Kill` (`client.go:295–308`).
- `failAll` retains the failure and delivers it to pending requests but does not write provider state (`client.go:369–386`).
- A pending synthesis receiving that failure publishes `degraded` (`client.go:790–797`), not `crashed`.
- `protocolViolation` calls `failAll` then `Kill` (`client.go:364–367`); `Kill` is `Terminate`, which sets `stopping` (`supervisor.go:291–294`, `:256–259`). The observer may then publish `absent`; `crashed` is not guaranteed.

Thus "fatal observation → `crashed`" is a proposed projection rule, not a uniform source behavior. It cannot be used to define current `crashed` without an owner decision.

## 4. Prior findings and authority distinctions

### 4.1 Source-verified findings

- `Synthesize` publishes `crashed` before an observed exit on a write failure.
- `SetState` can leave `processState=alive` when publishing `crashed`.
- `Status.ProcessAlive` and `Supervisor.ProcessAlive()` use different evidence and can disagree.
- `observeExit` publishes `crashed` only after `Wait` returns, and only when provider state is not `stopping`.
- `readLoop`/`failAll` do not write `crashed` by themselves.
- `cleanupRejected` and `Start` use `Supervisor.ProcessAlive()`, not `Status.ProcessAlive`.

### 4.2 Prior reported findings

- R13 and U6 in the prior assessment are substantially confirmed. The prior assessment's statement that a `Terminate` timeout can leave a live process with provider `crashed` is corrected above: `Terminate` leaves `stopping`.
- P17 O18/O19 and §7 describe `crashed` for fatal transport/protocol and identify U6. P17 is **not approved**, so its projection rules are not current semantic authority.
- Contract rev3 A6 proposes `crashed` for fatal `failureObserved` even with zero exit. Contract rev3 is a **proposed revision for independent review**, not an accepted amendment. It is a competing proposal, not a definition.
- Amendment draft A2-D4 proposes process-only liveness and only-unexpected-exit `crashed`. It is **NOT ADOPTED** and on conditional hold. It supports the process-outcome interpretation but is not authority.
- The lifecycle decision separates provider state from disposition and forbids `unknown → crashed`. It does not define `crashed`; its exclusions remain open.

### 4.3 Inferences

- The name `crashed` and `observeExit` support a process-outcome meaning.
- The `Synthesize` write supports a provider-condition meaning.
- The status surface currently implements neither meaning consistently.
- A process-outcome-only definition would make the `Synthesize` write non-conforming and would require the liveness boolean to use process state only.
- A provider-condition definition would preserve the write but leave `ProcessAlive=false` for a live process, and would make "crashed" a condition label rather than a process outcome.

### 4.4 Unresolved questions

- May `crashed` be published before an exit is observed?
- May a live process ever be `crashed`?
- Is a fatal observation alone enough for `crashed`, or must an exit fact also be present?
- Is `crashed` a terminal decision, or only a provider-state projection? OD-3 does not answer this.
- What should a live transport/protocol failure publish: `degraded`, `unavailable`, or an internal disposition?
- How does `ProcessAlive` relate to `State` and `ProcessState`? The source currently gives two answers.
- Does `crashed` mean a process exit, a provider failure, or a session-fatal condition? No single source definition exists.

## 5. Competing interpretations

| Interpretation | What it asserts | Source fit | Main conflict / consequence |
|---|---|---|---|
| **A. Process-outcome only** | `crashed` requires an observed process exit; a live write failure must use another state. | Strong fit with `observeExit` (`supervisor.go:183–190`). | Conflicts with `Synthesize` (`client.go:756`), which publishes before exit. Would require that write to change or be declared non-conforming. |
| **B. Provider/runtime failure** | `crashed` may describe a fatal helper or transport condition while the process is alive. | Fits `Synthesize`; matches part of P17/contract rev3 proposals. | Conflates condition with process outcome; `Status.ProcessAlive` becomes false for a live process; the name overstates process death. |
| **C. Split** | `crashed` = observed process exit; live transport failure uses `degraded`/`unavailable` or internal `outcome_unknown`. | Uses the existing `observeExit` meaning and preserves `degraded`/`unavailable` for live conditions. | Requires a projection rule for the write-failure path; the choice between `degraded` and `unavailable` remains open. |
| **D. Fatal-observation projection** | Any fatal protocol/transport observation projects `crashed`, even with a zero exit. | Proposed by P17 §7 and contract rev3 A6. | Not source-verified as uniform behavior; conflicts with process-outcome meaning and with U4's zero-exit case. |
| **E. Terminal decision** | `crashed` is the immutable terminal outcome. | None in the current supervisor; state can be overwritten by `Terminate`/`observeExit`. | OD-3 governs terminal immutability; the source has no latch. `crashed` is a provider state, not automatically a terminal decision. |

A narrower “fix only `client.go:756`” option (register alternative D) is not a full definition. It can remove the live-write contradiction while leaving the meaning of `crashed` and the status liveness rule unresolved.

### 5.1 `ProcessAlive` independence

`ProcessAlive` can be defined independently of provider state. It answers whether the supervisor currently owns a live process; provider state answers the provider's operational condition. They are different questions.

A process-only predicate is:

- true for `starting`, `alive`, and `stopping`;
- false for `absent` and `exited`;
- not derived from aggregate `State`.

Current source does not apply that predicate consistently:

- `Status` computes `ProcessAlive = processAlive(ProcessState, State)` and `processAlive` vetoes `crashed` and `unavailable` (`supervisor.go:324`, `:351–356`). This is not independent of provider state.
- `Supervisor.ProcessAlive()` is process-handle based (`proc != nil && processState != exited`; `supervisor.go:335–343`), but its predicate is not identical to the strict five-state rule: it treats every non-`exited` state as alive, including `absent` if that state were paired with a non-nil handle. The source should not be assumed to enforce the strict predicate.
- `cleanupRejected` and `Start` use `Supervisor.ProcessAlive()` (`client.go:859–864`, `:200–208`); status consumers use `Status.ProcessAlive`. The two surfaces can disagree, and the current source shows both answers for the same live process.

A process-only definition would allow `ProcessAlive=true` while provider state is `crashed`, `unavailable`, or `stopping`. That is the intended separation: a process can be alive while its provider/runtime condition is failed or while termination is pending. It also allows `ProcessAlive=false` after an observed exit even if the provider state has not yet been reconciled.

Implementation implications, without proposing source changes:

- `Status.ProcessAlive` would need to stop consulting provider state.
- The `processAlive` helper would need a process-only predicate identical to the `Supervisor.ProcessAlive()` ownership rule.
- `Supervisor.ProcessAlive()` would need an explicit `absent`/`exited` false rule if the strict five-state predicate is intended.
- Consumers that use `Status.ProcessAlive=false` to infer that no process exists would need to be reassessed; the process may still be live while the provider is `crashed` or `unavailable`.
- `Start` and `cleanupRejected` already use the supervisor method; their behavior on a live provider-failed process is a separate consequence and is not changed by this assessment.
- No code change is proposed or authorized.

## 6. Consequences for OD-4 and adjacent decisions

### 6.1 OD-4

OD-4 asks what public provider state, if any, should represent a nonzero exit after a known deliberate action. OD-5 does not decide OD-4.

The non-evidence points are explicit: `stopping` or a stop request, temporal ordering of observations, a successful `Terminate`/`Kill` result, and a nonzero exit code do not independently prove that a deliberate action caused the exit. A successful termination result is an API outcome, not a causal finding; `Terminate` waits for the observer, but it does not establish which action caused the exit. OD-8's A-E standard is required for deliberate-action causation. The same signals do not independently establish that the shutdown was orderly.

The connection is:

- If `crashed` is process-outcome-only, a deliberate action that produces a nonzero exit is not automatically `crashed`; the source currently maps a nonzero exit observed while `stopping` to `absent`, and P17 proposes internal `terminated_by_action`. OD-4's choice among `absent`, internal `terminated_by_action`, no state, or a separately approved state remains open.
- If `crashed` is broad enough to cover fatal/provider conditions, it competes with `terminated_by_action` for the same event and risks encoding causation in a provider-state label.
- The clean separation is: `crashed` describes a process outcome; `terminated_by_action` describes a causal/disposition judgment; they are not synonyms and should not be used as substitutes.

This assessment does not decide OD-4, OD-4a, or OD-4b.

### 6.2 Accepted OD-2, OD-3, and OD-8

- **OD-2:** `crashed` evidence must be bound to the same process instance. The current `Status` has no generation and does not prove which instance it describes. This assessment does not reopen OD-2.
- **OD-3:** `crashed` is not automatically an immutable terminal decision. The source can overwrite it through `Terminate`/`observeExit`; OD-3's first-terminal principle and its unresolved publisher/late-evidence exclusions remain untouched.
- **OD-8:** `crashed` is not proof that an action caused an exit. A fatal observation or a `crashed` state does not establish causality; the OD-8 standard applies to any causal claim. U5 remains open.

### 6.3 OD-1 and the frozen contract

- **OD-1:** evidence-record ownership remains deferred. This assessment describes source behavior and semantic options; it does not assume an authoritative recorder.
- **Contract:** the thirteen-value set remains frozen. A process-outcome-only definition does not add a state; it changes which writes are conforming. Any public projection change still requires the separate contract gate.

### 6.4 U4, U5, U6, and R13

- **U4 (zero exit after fatal protocol):** if `crashed` requires an observed exit not established as orderly, a zero exit after a fatal protocol observation is not automatically `crashed`. The case remains open; P17's fatal-wins rule is not settled.
- **U5 (nonzero exit after deliberate kill):** `crashed` and `terminated_by_action` are different claims. A deliberate kill's nonzero exit is not causal proof and is not automatically a crash. U5 remains **OPEN**.
- **U6 (provider `crashed` while process alive):** the source-verified contradiction would be addressed by a process-outcome-only definition plus process-based `ProcessAlive`; but U6 is not closed until the owner decides and the implementation is separately verified. This assessment does not close it. The §3.2 correction is a source-level finding only; it does not establish runtime behavior.
- **R13:** the semantic question is identified; the source defect is confirmed. R13 remains open until a decision and independent verification.

## 7. Recommended decision formulation for OD-5

The following is a **recommendation for owner consideration**, not a recorded disposition or amendment.

> **OD-5 — `crashed` semantics (proposed):**
>
> 1. `crashed` describes an **observed process-lifecycle outcome**, not a physical crash mechanism, fault, intent, or causal explanation. An observed process exit is a necessary fact; provider/runtime failure, a fatal observation, or a stop request is not a substitute for it.
> 2. In this definition, **“unexpected” means that an orderly shutdown has not been established**. It does not mean that the exit was unintended, uncaused, or caused by an action, and it does not assert a provider/runtime failure.
> 3. `crashed` **must not** be published solely because a transport write failed, a protocol/transport observation was fatal, a provider/runtime condition failed, a stop was requested, `stopping` was observed, or a lifecycle record was closed. A live process must not be represented as `crashed`.
> 4. `stopping`, a stop request, temporal ordering of observations, a successful `Terminate`/`Kill` result, and a nonzero exit code are **not independently sufficient** to establish that a deliberate action caused the exit. Deliberate-action causation requires the accepted OD-8 standard (A-E), not A-O association alone.
> 5. Orderly shutdown is a **separate lifecycle outcome**. It is not established by `stopping`, a stop request, a termination API result, temporal order, or exit code alone. The source's `stopping` branch is an intent-based classification, not proof of orderly shutdown.
> 6. A live transport, protocol, or provider failure uses an existing provider state selected by the applicable projection rules (`degraded` or `unavailable`) or an internal disposition. `outcome_unknown` remains internal and is not added to the provider-state set. The choice between `degraded` and `unavailable` is not decided by OD-5.
> 7. `ProcessAlive` is derived solely from `TTSProcessState` (`starting`, `alive`, `stopping` = alive; `absent`, `exited` = not alive), not from aggregate provider state. The status fields must not contradict the supervisor's process ownership.
> 8. `crashed` is not a synonym for a terminal decision, a fatal observation, or `terminated_by_action`. Terminal publication and immutability remain OD-3; causal attribution remains OD-8; public projection for an action-caused nonzero exit remains OD-4. A deliberate action-caused exit may still require a separate projection decision; OD-5 does not infer that projection from `crashed`.
> 9. A fatal observation followed by an exit may be evidence for a later classification, but the fatal observation alone is not a process exit and does not establish `crashed`.

### 7.1 Why this formulation

- It preserves the one source path that genuinely knows an exit occurred (`observeExit`).
- It stops `Synthesize` from using a provider-state label as a substitute for an exit fact.
- It resolves the `Status`/`ProcessAlive` contradiction without adding a public state.
- It keeps the observed exit fact, orderly-shutdown classification, and OD-8 causal attribution separate; a stop request, observation order, termination API result, or nonzero status cannot substitute for any of those claims.
- It keeps `crashed` distinct from fatal observation, terminal decision, and causal attribution.
- It leaves U4, U5, U6, OD-4, OD-8, and OD-3's exclusions open.

### 7.2 What it would make non-conforming

If accepted, the current `client.go:756` write would be non-conforming because it publishes `crashed` before an observed exit. The current `Status` liveness derivation would also be non-conforming because it lets provider state veto process liveness. The current `observeExit` `stopping` branch would also not be proof of an orderly shutdown; it is an intent-based classification that publishes `absent` for a zero or nonzero exit. This is a semantic finding, not an implementation instruction; no source change is authorized.

### 7.3 What it would not decide

- Whether a nonzero exit after a deliberate action is an orderly shutdown, `absent`, internal `terminated_by_action`, or another public projection (OD-4).
- Whether a zero exit after a fatal protocol observation is `crashed`, `absent`, or `outcome_unknown` (U4).
- Whether the existing write failure should project `degraded` or `unavailable` in every case.
- Whether `crashed` is an immutable terminal decision (OD-3).
- Whether the source implementation actually satisfies the semantic rule (implementation gate).
- Whether U5 satisfies OD-8 (it does not, on current evidence).

## 8. What this assessment does not decide

This assessment does **not**:

- record an OD-5 disposition, acceptance, rejection, or deferral;
- change OD-1, OD-2, OD-3, OD-4, OD-5, OD-6, OD-7, OD-8, OD-9, or U1–U8;
- close U4, U5, U6, or R13;
- decide public projection for `terminated_by_action`;
- decide whether a fatal observation may ever be `crashed`;
- decide the meaning of a terminal decision or its publisher;
- amend the thirteen-value provider-state contract;
- approve P17, create rev4 or P18, or authorize implementation;
- modify source, tests, configuration, or any frozen/supporting document;
- run tests, builds, benchmarks, or runtime experiments;
- claim that the current source is conforming to the recommended formulation.

## 9. Verification note

No existing artifact other than this assessment was modified. The decision register remains at SHA-256 `50350c44…`, 309 lines; P17, the corrected r3 response, the lifecycle decision, contract rev3, the OD-8 assessment, clarification memo, challenge review, and decision-ready proposal remain unchanged. No source, test, configuration, or generated file was modified; no tests, builds, or runtime experiments were run.

**OD-5 remains UNRESOLVED.**
