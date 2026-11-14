# C1 / P17 — OD-8 Causality Assessment (Read-Only)

**Date:** 2026-11-09

**Subject:** What evidence is sufficient to attribute a process exit to a deliberate action, rather than merely observing an exit after an action was requested?

**Status of this document:** Assessment only. It records no OD-8 disposition, changes no OD or U status, and authorizes no implementation.

## 1. Decision boundary preserved

This assessment starts from the current register and does not reopen settled or deferred items:

| Item | Current status | Boundary used here |
|---|---|---|
| OD-1 | **DEFERRED** | The authoritative evidence-record owner is not assumed. |
| OD-2 | **ACCEPTED** | Supervisor generation is the process-lifetime identity concept; session is the outer scope; PID is corroborating metadata. The source does not yet retain that identity with actions or exits. |
| OD-3 | **ACCEPTED** | One immutable first terminal decision per process instance, with one linearization point. Evidence-record ownership, publisher authorization, synchronization, late-evidence reconciliation, and lifecycle/shutdown behavior remain unresolved. |
| OD-4 through OD-9 | **UNRESOLVED** | No disposition is selected here. |
| U1–U8 | **OPEN** | No case is closed here. |
| Provider-state contract | **FROZEN** | The thirteen values are unchanged and no public projection is proposed. |
| P17 | **NOT APPROVED** | This assessment does not approve it. |

OD-3's accepted principle is narrower than the full terminal-outcome problem. It says a first terminal decision is immutable for that process instance. It does **not** settle who owns the evidence record, who may publish, how concurrent writers are synchronized, how late evidence is reconciled, or how shutdown interacts with a terminal decision. Those exclusions are material to OD-8 and are not treated as resolved.

OD-1 remains deferred. This assessment can describe what evidence would be sufficient in the abstract, but it cannot assume that an authoritative recorder exists or that any particular component owns it.

## 2. Exact OD-8 question and current status

The register's OD-8 entry is at `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md:192–204`. Its question is:

> **What evidence is sufficient to establish that a termination action or a fatal observation caused an exit, rather than merely preceding it?**

Its status is **UNRESOLVED**. Its recorded alternatives are:

- **A.** No causal claim without an action record, an exit fact, and a recorded ordering; otherwise `outcome_unknown`.
- **B.** Treat action-before-exit as sufficient.
- **C.** Require action completion and exit status, but not a global timestamp.
- **D.** Require independent external evidence, for example an OS exit status and a helper-side record.

The register labels **A** as the non-binding proposed default. The register's own consequence column says A is the most conservative and keeps unknown cases honest. It does not say A is sufficient to prove causation. This assessment takes that distinction seriously: **A is a useful floor, but as written it is not a sufficient causal test.**

## 3. Source inventory

The inspection was static. No tests, builds, benchmarks, or runtime experiments were run. Line references identify the current source revision inspected; they are not claims about behavior outside the code paths shown.

### 3.1 Supervisor and process lifecycle

| Location | What the source does | What it records, returns, or loses |
|---|---|---|
| `internal/audio/tts/helper/supervisor.go:28–40` | `Process` interface exposes `Send`, `Receive`, `Close`, `Stderr`, `Wait`, `Kill`, and `PID`. | `Wait` returns only an `error`. There is no exit-status object, start time, instance token, or action identity. |
| `internal/audio/tts/helper/supervisor.go:64–108` | `Supervisor` holds one process handle, a `procDone` channel, provider/process state, `lastFailure`, and `waitErr`. | No supervisor generation, action list, action outcome, exit code, timestamp, sequence number, or evidence log is stored. |
| `internal/audio/tts/helper/supervisor.go:110–164` | `Start` launches the helper, installs the new `proc` and `procDone`, resets `waitErr`, and starts the exit observer. | Replacement is represented by overwriting the current handle. No incarnation token or generation is written into later evidence. |
| `internal/audio/tts/helper/supervisor.go:166–195` | `observeExit` calls `p.Wait()`, takes the supervisor lock, clears `proc` only if it still equals `p`, stores the wait error, and classifies the provider state. | The classification reads `providerState == stopping`: if true, the provider becomes `absent`; otherwise `crashed` with a failure. It does not read an action identity, action result, timestamp, sequence, exit code, or competing fatal observation. |
| `internal/audio/tts/helper/supervisor.go:197–204` | `setUnavailable` records a failure and an unavailable provider state. | No action or exit record is created. |
| `internal/audio/tts/helper/supervisor.go:233–244` | `SetState` writes provider state and `lastFailure`; when a process exists and the new state is not `stopping`, it also sets process state to alive. | State is overwritten in place. The write is not an append-only evidence record and carries no causal ordering. |
| `internal/audio/tts/helper/supervisor.go:246–289` | `Terminate` takes the current process and done channel, sets process/provider state to `stopping`, closes stdin, waits for done or context/grace expiry, kills on timeout, and returns `nil`, a context error, or a timeout error. | It records no action intent, action ID, target instance token, request time, grace start, completion time, or result. A `nil` return means the done channel was observed, not that the action caused the exit. |
| `internal/audio/tts/helper/supervisor.go:291–294` | `Kill` is an alias for `Terminate(ctx, defaultShutdownTimeout)`. | A deliberate kill and an ordinary terminate are not distinguishable in retained evidence. |
| `internal/audio/tts/helper/supervisor.go:306–333` | `Status` projects provider status from current state fields. | It does not expose a process-outcome record or an action outcome. |
| `internal/audio/tts/helper/supervisor.go:335–343` | `ProcessAlive` returns true when a process handle exists and process state is not `exited`. | Liveness is derived from stored state, not from an independently observed OS process status. |
| `internal/audio/tts/helper/supervisor.go:442–456` | `osProcess` holds a command, PID, stderr buffer, and cached wait error. | The PID is metadata on the handle; it is not attached to an action or exit evidence record. |
| `internal/audio/tts/helper/supervisor.go:502–511` | `osProcess.Wait` calls `cmd.Wait()` and stores the returned error. | The error is retained as an error. Exit code, signal, and observation time are not extracted or stored. |
| `internal/audio/tts/helper/supervisor.go:512–522` | `osProcess.Kill` kills the process tree and the process. | It returns an error to its caller but creates no durable action record. |
| `internal/audio/tts/helper/supervisor.go:523` | `PID` returns the stored PID. | PID is available on demand, not bound into exit evidence. |

### 3.2 Client calls, failure propagation, and shutdown

| Location | What the source does | What it records, returns, or loses |
|---|---|---|
| `internal/audio/tts/helper/client.go:191–259` | `Client.Start` guards against an already-started client, starts the supervisor, and publishes the transport. | Replacement is a client-level lifecycle event. No supervisor generation or action identity is carried into the new process's evidence. |
| `internal/audio/tts/helper/client.go:295–318` | `readLoop` receives and decodes frames and calls failure handling on receive/decode errors. | It does not record whether an exit was already in progress, or whether a later action was the cause. |
| `internal/audio/tts/helper/client.go:364–367` | `protocolViolation` calls `failAll` and then `Kill`. | The kill's result is not retained as an action outcome. The fatal observation and the kill are not causally linked in a record. |
| `internal/audio/tts/helper/client.go:369–386` | `failAll` retains the first read failure, clears pending requests, and delivers the failure. | There is no timestamp or sequence. A later or competing failure is not retained as a separate causal account. |
| `internal/audio/tts/helper/client.go:398–443` | `callOn` sends a request and waits for a response or failure. | Request IDs identify protocol requests, not process-termination actions. |
| `internal/audio/tts/helper/client.go:445–469` | `begin` and `retire` allocate and remove pending request IDs. | These are request lifecycle operations, not evidence of process action or exit causation. |
| `internal/audio/tts/helper/client.go:645–677` | `Shutdown` performs shutdown handling, may call `Terminate`, and clears `started`/`transport`. | The termination result is not retained as a structured action record. Clearing the transport does not establish why the process exited. |
| `internal/audio/tts/helper/client.go:679–858` | `Synthesize` sends requests, handles responses, and on send failure can set `crashed` and clean up. | `SetState(crashed, ...)` is a provider-state write, not a process-exit finding with a cause. |
| `internal/audio/tts/helper/client.go:859–864` | `cleanupRejected` conditionally calls `Terminate` and ignores its result. | A cleanup kill can occur without any retained action outcome. |
| `internal/audio/tts/helper/client.go:866–874` | `watchCancellation` calls cancellation/termination and ignores the result. | Cancellation can be an unrecorded action relative to exit evidence. |

The callers discard or fail to retain `Terminate` results at several points, including `client.go:248`, `:654`, `:661`, `:664`, and `:861`; the cancellation path at `:871` likewise does not preserve a structured outcome. These references are important because an action can be requested, and even return without an error, without the caller retaining evidence that would support attribution.

### 3.3 Identity and contract types

| Location | What the source provides | Limitation for OD-8 |
|---|---|---|
| `internal/audio/tts/helper/protocol.go:259–268` | `newSessionID` creates a random session identifier; `newRequestID` combines that session with a counter. | Neither identifies a supervisor generation or a termination action. A request ID is not an action ID. |
| `internal/audio/types/provider.go:8–16` | Five process states: `absent`, `starting`, `alive`, `stopping`, `exited`. | `crashed` is not a process state. `stopping` describes a state, not an action outcome. |
| `internal/audio/types/provider.go:18–36` | The thirteen frozen provider states. | No state encodes action identity, action outcome, exit code, causal link, or evidence completeness. |

### 3.4 What the source does not contain

A static inspection found no retained record of:

- a supervisor generation or incarnation token attached to actions or exits;
- a unique termination-action identifier, action kind, or target-instance binding;
- whether an action was requested, accepted, attempted, completed, rejected, or timed out;
- an exit code, signal, or structured process outcome;
- timestamps or a monotonic sequence covering action and exit observations;
- a causal mechanism connecting an action to an exit;
- a competing-explanation record for a pre-existing exit, an independent failure, or a shutdown-time failure;
- an authoritative evidence owner or publisher under OD-1/OD-3.

The absence is not a claim that such evidence could never exist. It is a statement about what the current implementation records.

## 4. Observed facts, inferences, and unknowns

### 4.1 Directly observed from the source

1. `Terminate` sets `providerState` to `stopping` before waiting and potentially killing.
2. `observeExit` distinguishes only whether the provider state is `stopping` at the moment it processes the wait result.
3. A nonzero `waitErr` after an exit observed while `stopping` is projected as provider `absent`, not as an action-caused termination. A nonzero exit observed outside `stopping` becomes `crashed`.
4. `Kill` is literally `Terminate` with a default grace, so the source does not distinguish a deliberate kill from an ordinary termination request.
5. `Terminate` returns `nil` when the done channel is observed; it can also return a context error or a timeout. Those returns are not a structured action record.
6. Several callers ignore the `Terminate` result, and `Shutdown` clears client lifecycle fields without retaining a causal record.
7. `waitErr` is stored as an error. No code, signal, or timestamp is stored with it.
8. There is no generation, action ID, sequence, or evidence log in the inspected source.

### 4.2 Inferences that the source does not establish

- That a call to `Terminate` was deliberate in the sense required by the question. The call site may be cleanup, cancellation, protocol handling, or shutdown; the source does not preserve that intent with the exit.
- That the action completed. A `nil` return means the done channel was observed; the process may have exited for an unrelated reason before or during the action.
- That the action preceded the exit. There is no timestamp or sequence, and the only ordering signal is the mutable `stopping` state at observation time.
- That the action caused the exit. The source has no mechanism evidence, no competing-explanation record, and no binding of action to exit.
- That a nonzero exit is a crash, or that `absent` is a clean termination. Both are projections from limited state, not causal findings.

### 4.3 Unknowns that are material to OD-8

- Whether a particular exit was already in progress when the action was requested.
- Whether an independent failure caused the exit before the action's kill could take effect.
- Whether the exit was a shutdown-time failure distinct from both the action and the original process condition.
- Whether the action was accepted, rejected, timed out, or completed.
- Whether the observed process was the instance the action targeted, especially after supervisor replacement or PID reuse.
- Whether a fatal observation and the action both explain different parts of the history.

The source cannot answer these from retained evidence. It can only report the state at observation and the wait error.

## 5. Three conclusions that must remain separate

| Conclusion | Meaning | What the current source can support |
|---|---|---|
| **A. Action observed** | A deliberate action was requested or attempted. | A call to `Terminate`/`Kill` can be observed statically at a call site; the source does not retain a durable action record. |
| **B. Temporal association observed** | An exit was observed after an action request or attempt. | At most, the source observes an exit while provider state is `stopping`. It does not record a timestamp or sequence, so even this association is inferred from mutable state rather than ordered evidence. |
| **C. Causal attribution established** | The evidence meets the proposed standard for attributing that exit to that action. | Not established by the current source. No inspected path records the required linkage. |

A stop request followed by a nonzero exit is **B**, not **C**. The process may already have been exiting, may have failed independently, or may have encountered a separate shutdown-time failure. The source's `stopping` check cannot distinguish those cases.

This assessment does not promote A or B into C. It also does not promote a fatal observation into a cause merely because an exit follows it.

## 6. Proposed minimum causal-evidence standard

The standard below is evidence-first. Terms such as "action ID," "generation," and "sequence" name evidence requirements, not mandatory implementation fields. A mechanism that satisfies the requirement without those exact names is acceptable; a mechanism that merely produces a timestamp or a log line without the required content is not.

### 6.1 Required dimensions for conclusion C

| Dimension | Minimum evidence | Why it is required | Current source status |
|---|---|---|---|
| **1. Instance identity** | The action and the exit fact are bound to the same authoritative process instance. The binding survives supervisor replacement and cannot be satisfied by a stale instance or a reused PID. | Without this, an action against one instance can be attributed to an exit of another. OD-2's generation concept is the relevant target, but the source does not retain it with actions or exits. | Partially conceptual only. A process pointer and PID exist; no generation is attached to evidence. |
| **2. Action identity** | The specific action is identifiable: at least its kind (terminate, kill, shutdown, cleanup, cancellation), its target instance, and a unique action identity or equivalent. | "The provider was stopping" is a state, not an action. `Kill` is an alias, so the source cannot even distinguish the action kind. | Absent. `Kill` aliases `Terminate`; no action ID or intent is stored. |
| **3. Action outcome** | The record distinguishes requested, accepted, attempted, completed, rejected, and timed out. A successful API return is not automatically completion or acceptance. | An action request is not an action result. A `nil` return can mean the done channel closed for an unrelated reason; a timeout can mean the action did not complete. | Absent as structured evidence. `Terminate` returns nil/error, and callers often discard it. |
| **4. Exit fact** | The actual process exit is observed, with its exit status or error and the observation context. Zero and nonzero must be distinguishable, and an error must not be treated as an exit code. | A nonzero exit alone does not identify a cause, and `waitErr` is an error rather than a structured status. | Partial. `waitErr` is stored, but there is no exit code, signal, or observation timestamp. |
| **5. Ordering** | A recorded ordering places the action and the exit observation in sequence within the same instance. A monotonic sequence or equivalent is preferable to wall-clock time; ordering alone is not causation. | Without ordering, even temporal association is conjecture. With ordering alone, the exit may still have been already in progress or independently caused. | Inferred only. The source uses the mutable `stopping` state at the moment of observation; no sequence exists. |
| **6. Causal linkage** | Evidence of a mechanism by which the action could produce the observed exit, plus evidence that the mechanism operated in this case. Examples include an OS-reported kill/signal outcome, a helper-side record of the action reaching the process, or a supervisor observation that the targeted process exited after the action's signal. The action must be distinguished from a pre-existing exit. | This is the step from B to C. Ordering and an action record can establish that the action happened before the exit, but not that it caused the exit. | Absent. No signal/cause record, no action-to-exit link, and no exclusion of a pre-existing or independent failure. |
| **7. Competing explanations** | Contradictory, incomplete, delayed, or ambiguous evidence is retained and its effect on attribution is stated. An attribution is not selected by class order, first-seen failure, last-writer state, or convenience. | A fatal observation and an action can both be present; an exit can be already in progress; a shutdown can fail separately. A standard that cannot represent competition will manufacture certainty. | Absent. Only a first failure and a current state are retained; competing accounts are not preserved. |

### 6.2 Governance dimensions for using C in a terminal decision

| Dimension | Minimum requirement | Why it is separate |
|---|---|---|
| **8. Terminal-decision interaction** | A causal attribution may support a terminal decision only when the evidence meets dimensions 1–7 and the publisher is authorized under the applicable OD-1/OD-3 rules. The standard does not itself grant publisher authority or evidence ownership. | OD-3 accepted the immutability principle, not the publisher, synchronization, or ownership contract. OD-1 remains deferred. |
| **9. Late evidence** | Later evidence can strengthen, weaken, or leave an attribution unresolved. It must remain recordable after a terminal decision; whether it can change a projection is a separate unresolved question. The immutable first terminal decision must not be silently rewritten, and late evidence must not silently disappear. | OD-3 distinguishes terminal immutability from the continuing evidence record. This assessment preserves that distinction without deciding projection mutability. |

### 6.3 How the standard maps to conclusions

- **C is established** only when dimensions 1–7 are satisfied together. Dimensions 8–9 govern whether and how C may be used in a terminal decision.
- **B only** applies when there is an action observation and an exit observation with some ordering, but dimension 6 (or another required dimension) is missing.
- **A only** applies when an action is observed but there is no adequately identified exit or ordering.
- **Unresolved** applies when evidence is contradictory, incomplete, or ambiguous. It must not be converted into C, `crashed`, or a clean `absent` by assumption.

This is a proposed minimum. A stricter standard is compatible with it. A standard that stops at dimensions 1–5 is not sufficient for C, because it cannot distinguish an action-caused exit from an exit that was already underway or independently caused.

### 6.4 Application to fatal observations

The OD-8 question covers a "termination action **or a fatal observation**." The same separation applies:

- A fatal handshake, protocol, or transport observation can establish that a fatal condition was observed.
- It does not by itself establish that the observation caused the process exit.
- A causal attribution for a fatal observation needs the same instance identity, observation identity, outcome/exit fact, ordering, causal linkage, and competing-explanation treatment. The "action identity" becomes the identity of the fatal observation and the policy that responded to it.

The source currently records fatal observations as failures and state changes; it does not record a causal chain from the observation to the process exit. This is a second reason not to treat a fatal observation plus an exit as a causal finding.

## 7. Illustrative evidence patterns

The following are **illustrative patterns, not source findings**. They show how the proposed standard would classify evidence; they do not assert that any pattern occurs in the current implementation.

### 7.1 Sufficient pattern (illustrative)

- The action and exit are bound to the same supervisor generation G.
- The action has a unique identity A and an explicit kind (for example, terminate), targeting G.
- A's record shows it was accepted and that the termination mechanism operated.
- The exit fact records that process G exited with an OS-reported outcome consistent with that mechanism.
- A's ordering precedes the exit in a monotonic sequence for G.
- No pre-existing fatal observation, independent failure, or competing exit explanation is present, or each is recorded and ruled out by evidence.

**Classification:** C — causal attribution is established for this instance. This does not automatically establish attribution for any other event or any other instance.

### 7.2 Insufficient pattern — action plus exit only (illustrative)

- A stop request is issued against G.
- G exits nonzero.
- No action outcome, exit status, sequence, or causal mechanism is retained.

**Classification:** A and B only. The exit is temporally associated with the request, but the evidence does not establish that the request caused it. The honest result is unresolved/`outcome_unknown`, not `crashed`, `absent`, or action-caused.

### 7.3 Ambiguous pattern — completed action, unexplained exit (illustrative)

- The action reports completion, and G exits nonzero.
- There is no timestamp, and the completion report could have been produced after G had already exited.
- There is no pre-existing fatal observation and no mechanism record.

**Classification:** unresolved. The action may have caused the exit, or the action may have observed an exit already in progress. Neither explanation is uniquely supported. The evidence must be retained; it must not be resolved by choosing the action because it is the only recorded event.

### 7.4 Competing-explanation pattern (illustrative)

- A fatal protocol observation occurs first.
- A termination action follows and completes.
- The process exits zero.

**Classification:** the fatal observation explains the session's fatal condition; the action may explain the exit, or the exit may be the expected consequence of shutdown. These are different claims. The evidence must distinguish "why the session became fatal" from "why the process exited." Neither class order nor a single projection table settles both.

## 8. U5 evidence-gap assessment

U5 is the separate application case: **nonzero exit after a deliberate kill**. OD-8's standard does not decide U5. This section compares what the standard would require with what the current source records.

| Required element | Present in current source? | Assessment |
|---|---|---|
| Same authoritative process instance | Partial/conceptual | A process pointer and PID exist, but no supervisor generation is attached to the action or exit. OD-2's accepted identity concept is not yet retained as evidence. |
| Specific action identity | No | `Kill` is an alias of `Terminate`; no action ID, kind, or target-instance token is recorded. |
| Action outcome | No | `Terminate` returns nil/error/timeout, but callers ignore it at several sites; no accepted/attempted/completed/rejected/timed-out record exists. |
| Exit fact with status | Partial | `waitErr` is stored, but it is an error, not an exit code plus cause. No signal or observation time is retained. |
| Ordering | Inferred only | There is no timestamp or sequence. The only ordering signal is whether `providerState == stopping` when `observeExit` processes the wait result. |
| Causal linkage | No | No record shows that the kill mechanism reached this process and produced this exit, and no record excludes an exit already in progress. |
| Competing explanations | No | A pre-existing failure, an independent exit, and a shutdown-time failure are not separately retained. `failAll` keeps a first failure; `SetState` overwrites state. |
| Terminal-decision handling | Not established for U5 | `observeExit` maps a nonzero exit observed while `stopping` to provider `absent`. OD-3's immutability principle is accepted, but publisher authorization, evidence ownership, synchronization, and late-evidence handling remain unresolved. |
| Late evidence | No | No append-only evidence record exists to strengthen, weaken, or leave the attribution unresolved. |

**Conclusion for U5:** the current source cannot give U5 a unique, source-justified causal classification. It can observe an action request and an exit; it cannot distinguish action-caused termination from an exit already in progress, an independent failure, or a separate shutdown-time failure. `outcome_unknown` is the honest internal disposition where such a distinction is required, but this assessment does not decide its public projection. U5 remains **OPEN**.

## 9. Dependencies and downstream boundaries

### 9.1 OD-3 — terminal immutability

OD-3's accepted principle gives OD-8 an important constraint: a first terminal decision for a process instance is immutable, and later observations must not silently rewrite it or disappear. But OD-3 does not settle:

- who owns the evidence record;
- who may publish the terminal decision;
- how concurrent publishers are synchronized;
- how late evidence is reconciled;
- how shutdown interacts with terminal publication.

OD-8 can define what evidence would justify an attribution, but it cannot make that evidence authoritative or make a publisher authorized. A terminal decision that does not meet the standard is not retroactively validated by OD-8.

### 9.2 OD-5 — meaning and scope of `crashed`

OD-5 asks whether `crashed` describes a process exit, a fatal helper failure, or a narrower condition. OD-8 is adjacent but not identical:

- If `crashed` is a process-outcome claim, the causal standard may inform when an exit is classified as action-caused versus crashed.
- If `crashed` is a fatal-helper-condition claim, an action-caused exit and a crash may be different dimensions.
- A fatal observation plus an exit is not automatically a crash, and an action-caused exit is not automatically not a crash.

This assessment does not resolve OD-5. It only notes that OD-8's standard cannot be applied to `crashed` until the meaning of `crashed` is settled.

### 9.3 OD-4a and OD-4b

- **OD-4a** is whether an action can be shown to have caused an exit. OD-8 supplies the evidentiary standard; OD-4a is the application question for a particular action/exit sequence. OD-8 does not decide any OD-4a case.
- **OD-4b** is the public projection after causal and lifecycle semantics are established. OD-8 does not decide any projection, and nothing here authorizes a new public state.

### 9.4 OD-1 — deferred evidence-record ownership

OD-1 remains deferred. The standard can be stated without assuming an owner, but the standard cannot be implemented or treated as authoritative until ownership is decided. This assessment therefore separates:

- **evidence requirements** (what must be recorded), from
- **evidence ownership and publication** (who records it and who may use it).

The absence of an owner is not a reason to weaken the standard; it is a reason not to claim that the standard is already operational.

### 9.5 U5 and the other U cases

U5 is an application case, not a definition. The same is true of U1–U4 and U6–U8. OD-8 defines the standard; it does not close any U case. In particular:

- U5 cannot be closed by saying "the standard exists."
- A future implementation that records the required evidence would still need a separate application decision for U5.
- A future terminal decision for U5 would still be subject to OD-1, OD-3's unresolved exclusions, and OD-5's meaning of `crashed`.

### 9.6 R09, R10, R11–R13

- **R09:** OD-3's terminal immutability principle is accepted, but the minimal record, allocator, causal link, action result, phase, retention, and publisher remain unresolved. OD-8 depends on those evidence pieces.
- **R10:** precedence remains a selection policy, not a demonstrated resolution of competing evidence. OD-8 is the missing causal standard; it must not be replaced by a class-order rule.
- **R11–R13:** unchanged. This assessment does not resolve quarantine reachability, `crashed` semantics, or any other finding.

## 10. Alternatives for OD-8

### Alternative A — No causal claim without action record + exit fact + ordering; otherwise `outcome_unknown`

**Strengths**

- It is conservative and avoids inventing a winner.
- It preserves the distinction between association and causation better than B.
- It gives unresolved cases an honest label rather than forcing `crashed` or `absent`.

**Weaknesses**

- As written, it does not define what counts as an "action record," an "exit fact," or "ordering." A request, a state flag, and an error could be read as satisfying all three while still providing only temporal association.
- It does not require a causal mechanism or competing-explanation handling. Dimensions 1–5 can be satisfied while dimension 6 is absent.
- It does not say whether the three items are a minimum floor or a sufficient test. If treated as sufficient, it repeats the R10 defect in a more careful vocabulary.
- It does not address fatal observations separately from deliberate actions.

**Unresolved questions**

- Is A a floor, a presumption, or a sufficient test?
- Does "action record" include action identity and outcome, or only the fact that a request was made?
- Does "exit fact" include an exit code/signal and observation context, or only a wait error?
- Does "ordering" require a monotonic sequence, or is a mutable state check enough?
- Is a causal mechanism required, or may attribution remain unresolved without it?

### Alternative B — Treat action-before-exit as sufficient

**Strengths**

- Simple and easy to apply.
- It would make many cases decidable without new evidence.

**Weaknesses**

- It is exactly the correlation-to-causation error R10 identifies.
- A process may already have been exiting when the action was requested.
- A separate failure during shutdown can produce the observed exit.
- It would close U5 by assumption, contrary to the evidence gap.

**Unresolved questions**

- None that make B acceptable as a causal standard. Ordering is necessary but not sufficient.

### Alternative C — Require action completion and exit status, but not a global timestamp

**Strengths**

- Stronger than B: it asks for an action outcome and an exit fact.
- It avoids relying on wall-clock synchronization alone.

**Weaknesses**

- "Action completion" can be reported after the process already exited; completion is not causation.
- "Exit status" identifies how the process ended, not why. A nonzero status is consistent with several causes.
- Without a recorded ordering, the action and exit may be concurrent or reversed.
- It still lacks instance identity, action identity, causal linkage, and competing-explanation handling.
- It may be impossible to apply consistently without some ordering mechanism, even if a global timestamp is not required.

**Unresolved questions**

- What counts as completion, and who observes it?
- Is a per-instance sequence required if a global timestamp is not?
- Does exit status include a signal or only a numeric code?
- How are contradictory completion and exit evidence reconciled?

### Alternative D — Require independent external evidence

**Strengths**

- It offers the best chance of a genuine causal link, especially if the external evidence includes an OS-reported termination cause or a helper-side record of the action reaching the process.
- It can distinguish some action-caused exits from independent failures that an internal state flag cannot.

**Weaknesses**

- The source does not currently retain the required external evidence.
- "Independent" is not self-validating: external evidence can be incomplete, delayed, contradictory, or unavailable.
- It still needs instance identity, action identity, ordering, and competing-explanation handling; external evidence alone does not supply those.
- Platform differences may make exit status/signal evidence uneven.
- It risks treating one external artifact as conclusive without checking whether the process was already exiting.

**Unresolved questions**

- Which external source is authoritative, and under OD-1 who owns it?
- What counts as independent evidence of causation rather than mere corroboration of an exit?
- How should conflicting external and internal evidence be retained and resolved?

### Assessment refinement (not a register alternative)

The alternatives above are the register's alternatives. This assessment does not add a disposition or a fifth register alternative. It does propose a clarification for the owner's consideration:

> Treat the action record, exit fact, and ordering as necessary but not sufficient. Require instance identity, action identity, action outcome, exit fact, ordering, causal linkage, and competing-explanation handling before concluding C. Use dimensions 8–9 only to govern terminal publication and late evidence. If any required dimension is missing, record the result as unresolved/`outcome_unknown` rather than attributing the exit.

This refinement preserves A's conservatism while closing its causal gap. It is still not a disposition and does not authorize implementation.

## 11. Decision-readiness recommendation

**OD-8 is not decision-ready as currently written.** It needs clarification before a disposition is recorded.

The main reason is that alternative A, the non-binding default, is stated as "action record + exit fact + ordering." That combination can establish A and B, but it does not establish C. The register's own question is about causation "rather than merely preceding it"; a standard that stops at ordering does not answer that question.

The entry also leaves several terms undefined:

- what an action record must contain (request only, or identity and outcome?);
- what an exit fact must contain (wait error only, or status/cause and context?);
- what ordering mechanism is acceptable (state check, sequence, or timestamp?);
- whether a causal mechanism is required;
- how competing explanations are retained;
- how fatal observations are treated relative to deliberate actions;
- how terminal publication and late evidence interact with the standard.

A decision-ready OD-8 would, at minimum:

1. State whether the listed evidence is a floor or a sufficient test.
2. Require causal linkage and competing-explanation handling before C.
3. Define the evidence items without mandating a particular implementation.
4. State the standard's application to fatal observations as well as deliberate actions.
5. Preserve OD-1's deferred ownership and OD-3's unresolved publisher/synchronization/late-evidence questions.
6. Keep U5 and every other U case open.

Until those clarifications are recorded, a disposition would either adopt an underspecified standard or decide causation by implication. Neither is warranted.

## 12. What this assessment does not decide

This assessment does **not**:

- record an OD-8 disposition, acceptance, rejection, or deferral;
- change OD-1 (deferred), OD-2/OD-3 (accepted), OD-4–OD-9 (unresolved), or U1–U8 (open);
- decide that U5, or any other U case, is action-caused;
- decide that a fatal observation caused any exit;
- decide the meaning of `crashed` (OD-5);
- decide public projection for `terminated_by_action` or any other internal disposition (OD-4a/OD-4b);
- decide evidence-record ownership (OD-1) or publisher authorization, synchronization, or late-evidence reconciliation (OD-3 exclusions);
- approve P17, amend the contract, create rev4 or P18, or change the thirteen-value provider-state set;
- implement an evidence logger, action record, causal-link mechanism, or lifecycle change;
- assert runtime behavior beyond the inspected source paths;
- treat a process exit, a nonzero exit, a successful API return, or an action request as interchangeable evidence.

## 13. Baseline and integrity statement

### 13.1 Baseline artifacts

| Artifact | Baseline |
|---|---|
| P17 proposal | `01d887c1…` |
| Corrected r3 response | `60ca1cb5…` |
| Lifecycle disposition decision | `aa136bc9…` |
| Contract rev3 | `ed0da1a7…` |
| Decision register | `4a2a24b908cbb382f05711d10750bc156e0abf7358eec1b7e98ef50116663a16` (244 lines) |

### 13.2 Files read

- `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md` — OD-8 entry at lines 192–204; status table and surrounding boundary material.
- `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-2026-11-09.md` — §7 precedence at lines 102–125; U5 at line 148; rule 4 and table at lines 112 and 124.
- `docs/review/DECISION-howl-tts-stage3-c1-lifecycle-disposition-2026-11-09.md` — provider-state/disposition separation, including lines 11–14, 22–24, and 42–47.
- `docs/review/ASSESSMENT-howl-tts-stage3-c1-p17-remaining-evidence-semantics-2026-11-09.md` — R10 at lines 117–137; U5 at lines 295–312; missing-evidence findings at lines 428–433.
- `internal/audio/tts/helper/supervisor.go` — `Process` interface at 28–40; `Supervisor` at 64; `Start` at 110–164; `observeExit` at 166–195; `setUnavailable` at 197–204; `SetState` at 233–244; `Terminate` at 246–289; `Kill` at 291–294; `Status` at 306–333; `ProcessAlive` at 335–343; `osProcess` at 442–456; `Wait` at 502–511; `Kill` at 512–522; `PID` at 523.
- `internal/audio/tts/helper/client.go` — `Start` at 191–259; `readLoop` at 295–318; `protocolViolation` at 364–367; `failAll` at 369–386; `callOn` at 398–443; `begin`/`retire` at 445–469; `Shutdown` at 645–677; `Synthesize` at 679–858; `cleanupRejected` at 859–864; `watchCancellation` at 866–874.
- `internal/audio/tts/helper/protocol.go` — session/request IDs at 259–268.
- `internal/audio/types/provider.go` — process states at 8–16; thirteen provider states at 18–36.

### 13.3 Files created

- `docs/review/ASSESSMENT-howl-tts-stage3-c1-od8-causality-2026-11-09.md` (this report).

### 13.4 Files modified

**None.** The decision register, P17, corrected r3 response, lifecycle decision, contract rev3, source, tests, and configuration were not modified.

### 13.5 Tests, builds, benchmarks, and experiments

**None.** This was a static, read-only assessment. No runtime experiment was performed, and no interleaving claim in this report is a runtime observation.

### 13.6 Limitations

- The helper is untracked in Git; the source assessment relies on the current working tree and the recorded helper digest, not on a tracked commit.
- The analysis is static. It does not prove how concurrent `Terminate`, `observeExit`, `failAll`, or `SetState` interleave at runtime.
- Platform-specific exit-status behavior was not tested. The source does not retain a structured exit status on any platform.
- OD-1's ownership question and OD-3's publisher/synchronization/late-evidence questions remain unresolved, so the standard described here is not an operational evidence contract.
- Any evidence requiring instrumentation, execution, or source changes was not obtained; that would require separate authorization.

### 13.7 Final status confirmation

**OD-8 remains UNRESOLVED.** The decision register is unchanged at SHA-256 `4a2a24b908cbb382f05711d10750bc156e0abf7358eec1b7e98ef50116663a16`, 244 lines. No OD or U disposition was recorded, no P17 approval occurred, and no implementation was started.
