# C1 — P17 U1–U8 Scenario Findings

**Date:** 2026-11-09

**Status:** READ-ONLY EVIDENCE ASSESSMENT. **P17 NOT APPROVED. CONTRACT REV3 FROZEN. C1 IMPLEMENTATION UNAUTHORIZED.** B-2 is accepted only at the independent-review level and does not resolve any U case.

**Subject:** the eight open scenarios U1–U8 in P17 rev2 and the helper/client source that bears on them.

**Method:** static document and source reading, hash verification, and bounded reasoning only. No test, build, benchmark, application execution, runtime experiment, source change, contract change, register change, P17 edit, or implementation was performed.

---

## 1. Baseline and owner provenance determination

### 1.1 Inspected inputs

| Input | Lines | SHA-256 | Status |
|---|---:|---|---|
| P17 rev2 candidate, `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-rev2-2026-11-09.md` | 538 | `651d38d3e2d97f981a606a39ea4a3117c2e748ff9f5941eb8bd9671b6e022f82` | Identity verified |
| Independent adversarial review, `docs/review/REVIEW-howl-tts-stage3-c1-p17-rev2-adversarial-651-2026-11-09.md` | 224 | `cbed00a4f9fb25153d2ce7555c5aff33cccf57dc103b5e016b15f3a5271e9196` | Identity verified; no blocking finding; not an approval |
| Decision register, `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md` | 434 | `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99` | Identity verified |
| Accepted remaining-evidence assessment, `docs/review/ASSESSMENT-howl-tts-stage3-c1-p17-remaining-evidence-semantics-2026-11-09.md` | 477 | `3c4ce2beb5681d8608e6e9e807b1b76562142b9f8adf2d8ea8847f6f7f2bbd93` | On-disk working-copy identity accepted by owner |
| Contract rev3, `docs/review/CONTRACT-howl-tts-stage3-c1-installation-lifecycle-rev3-2026-11-09.md` | 183 | `ed0da1a75d4f62874f873979e6667f66f19541a7523ff96369f2b44e8391195e` | Frozen input |
| Original P17, `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-2026-11-09.md` | 166 | `01d887c18700d63070f1a2649659281ec1e875a21b49375adea12f8f11d917bf` | Frozen original |

### 1.2 Accepted assessment-provenance determination

The owner has accepted the existing on-disk assessment as the authoritative current working copy for this review. The accepted working-copy identity is **477 lines / `3c4ce2beb5681d8608e6e9e807b1b76562142b9f8adf2d8ea8847f6f7f2bbd93`**.

The historical identity has **not** been independently authenticated. The previously carried-forward expected value — `3c4ce2beb5681d8608e6e9e807b1b76562142b9f8adf2d8ea8847f6f2bbd93` — is only 62 hexadecimal characters and is not a valid SHA-256. No contemporaneous full digest was recovered.

This report does not claim that the historical digest was recovered, that the expected and on-disk values match, or that historical provenance was proven. It uses the owner-accepted on-disk file as an input only. The accepted assessment was not modified.

### 1.3 Source identity

| Item | Value |
|---|---|
| Repository HEAD | `9c7d2a17dd7f702d7db6d95d9316e7640e34bf6f` (`9c7d2a1`) |
| Helper tracked by Git | **0 files**; the helper tree is untracked |
| Helper Go files | 19 |
| Helper Go lines | 4,992 |
| Helper content digest | `d0df8b2823389b305698fc06eda28e0e01f06259fdc69174cf7fb3851e271a51` |

The helper content digest binds the relative paths and contents of the 19 `*.go` files. The Git HEAD pin does not establish helper identity because the helper directory is untracked. The digest does not bind non-Go files, runtime state, platform behavior, or a particular execution.

---

## 2. Scope, labels, and limitations

### 2.1 Evidence labels

- `[Observed]` — directly visible in the inspected source or document text.
- `[Source-derived]` — a control-flow consequence of the visible source, not a runtime observation.
- `[Calculated]` — line counts, hashes, and other deterministic calculations.
- `[Inferred]` — a bounded interpretation that follows from the evidence but is not itself directly observed.
- `[Unestablished]` — the available evidence does not determine the proposition.
- `[Superseded]` — a prior claim displaced by a later accepted decision or source correction.

`[Observed]` is not used for runtime behavior. No runtime execution was performed.

### 2.2 Limits

1. Static source reading can establish reachable paths, not which goroutine wins a concurrent interleaving.
2. No terminal, action, stream-close, or exit-cause record exists in the source.
3. No runtime evidence, platform evidence, or helper-side diagnostic trace was collected.
4. `observeExit` and `Terminate` are concurrent writers to supervisor state; line-level reachability does not prove their order in a real run.
5. Missing evidence is not treated as proof that an event did or did not occur.
6. P17's current wording is assessed as an unapproved proposal, not as current semantic authority.

---

## 3. Source trace

### 3.1 Control-flow references

| ID | Source location | Symbol/function | What the reference establishes | Label |
|---|---|---|---|---|
| S-01 | `internal/audio/tts/helper/client.go:191–259` | `Client.Start` | Starts a process and reader, performs hello handshake, and on failure calls `Terminate`, discards the `waitForReader` result, publishes `incompatible`/`unavailable`, then clears `started` and `transport`. | `[Observed]` |
| S-02 | `client.go:295–318` | `Client.readLoop` | A receive error becomes `helper_crashed` when the original code is `internal_error`; `failAll` is called; only `invalid_frame`, `frame_too_large`, and `protocol_violation` trigger `Kill`. No stream-close provenance is retained. | `[Observed]` |
| S-03 | `client.go:321–347` | `Client.dispatch` | Empty ID, duplicate terminal, and unknown ID become `protocol_violation`; a retired ID is dropped without a record. | `[Observed]` |
| S-04 | `client.go:349–362` | `waitForReader` | Returns nothing whether the reader joins or the timeout expires. | `[Observed]` |
| S-05 | `client.go:364–386` | `protocolViolation`, `failAll` | `protocolViolation` calls `failAll` then `Kill`; `failAll` stores only the first read failure and delivers it to pending requests. | `[Observed]` |
| S-06 | `client.go:398–443` | `Client.callOn` | Generic timeout and cancellation outcomes are returned without operation, phase, or deadline-origin evidence; a send failure returns `helper_crashed` without terminating the process. | `[Observed]` |
| S-07 | `client.go:645–674` | `Client.Shutdown` | Sends `shutdown`; may return early after `Terminate`; otherwise calls `Terminate`, then clears `started` and `transport` regardless of whether termination returned a timeout. | `[Observed]` |
| S-08 | `client.go:679–856` | `Client.Synthesize` | On a synthesis send failure at `:755–758`, publishes `TTSStateCrashed` before `cleanupRejected`; a pending failure at `:790–797` publishes `TTSStateDegraded`. | `[Observed]` |
| S-09 | `client.go:859–864` | `cleanupRejected` | Calls `Terminate` only if `ProcessAlive()` and either force or output exists; ignores the termination result. | `[Observed]` |
| S-10 | `internal/audio/tts/helper/supervisor.go:110–164` | `Supervisor.Start` | Creates the process, installs the exit observer, and sets the initial handshaking state. | `[Observed]` |
| S-11 | `supervisor.go:166–195` | `Supervisor.observeExit` | On exit, if provider state is `stopping`, publishes provider `absent` for either zero or nonzero exit; otherwise publishes provider `crashed`; stores `waitErr`. | `[Observed]` |
| S-12 | `supervisor.go:233–242` | `Supervisor.SetState` | Overwrites provider state and `lastFailure`; sets process state to `alive` whenever a process is owned and the new state is not `stopping`. | `[Observed]` |
| S-13 | `supervisor.go:246–288` | `Supervisor.Terminate` | Sets `stopping`, closes stdin, waits for the observer, kills on context cancellation or grace expiry, and can return `timeout` after a second grace while a process remains owned. No action intent, start, completion, or result is recorded. | `[Observed]` |
| S-14 | `supervisor.go:291–293` | `Supervisor.Kill` | `Kill` aliases `Terminate` with the default grace; it is not an independently recorded action. | `[Observed]` |
| S-15 | `supervisor.go:305–342`, `:351–356` | `Status`, `processAlive` | `Status` derives `ProcessAlive` from both process and provider state; the helper returns false for `crashed` or `unavailable`. | `[Observed]` |
| S-16 | `internal/audio/tts/helper/protocol.go:502–545` | `validateHello` | All nine validation failures return `protocol_violation`; a missing mandatory service capability is one, with a `missing` detail. | `[Observed]` |
| S-17 | `protocol.go:566–578` | `missingCapabilities` | Computes missing required capability IDs from the hello capability list. | `[Observed]` |
| S-18 | `internal/audio/tts/helper/errors.go:15–21` | `Failure` | Carries code, message, retryable flag, field, and detail; carries no process generation, sequence, timestamp, phase, action, or exit-cause field. | `[Observed]` |
| S-19 | `internal/audio/types/provider.go:8–36` | state sets | Five process states and the frozen thirteen provider states; neither contains `quarantined`. | `[Observed]` |

### 3.2 Existing source-loss findings

The P17 E1–E9 findings remain relevant and unresolved:

| ID | Loss | Source reference |
|---|---|---|
| E1 | `failAll` keeps only the first read failure; later failures are discarded. | `client.go:369–386` |
| E2 | `SetState` overwrites `lastFailure`; only the last published failure survives. | `supervisor.go:233–241` |
| E3 | The stopping branch discards a nonzero exit from classification. | `supervisor.go:176–187` |
| E4 | Termination intent is not a parameter, and `Kill` aliases `Terminate`. | `supervisor.go:246,291–293` |
| E5 | `waitForReader` returns nothing on timeout. | `client.go:349–362` |
| E6 | Stream closure is marked crashed but does not terminate for that code. | `client.go:299–308` |
| E7 | Late terminal evidence can be dropped by dispatch. | `client.go:339–342` |
| E8 | Cancellation and cleanup paths ignore errors from `Cancel` and `Terminate`. | `client.go:248,654,661,664,861`; `supervisor.go:246–288` |
| E9 | `Shutdown` calls `Terminate`, then clears lifecycle and transport state without establishing successful reaping. | `client.go:654–674`; `supervisor.go:246–288` |

These are evidence losses, not new public states.

---

## 4. U1–U8 matrix

| ID | Scenario | Source-established core | Current P17 disposition | Evidence conclusion |
|---|---|---|---|---|
| U1 | Orderly EOF racing reader crash | Reader can classify an EOF as `helper_crashed`; `Shutdown` later sets `stopping`; observer maps `stopping` to `absent`, otherwise `crashed`; no stream-close provenance or ordering record. | OPEN | **OPEN** |
| U2 | Shutdown timeout versus protocol violation | `Shutdown` and `protocolViolation`/`Kill` can call overlapping `Terminate`; timeout, exit, and failure records are not ordered. | OPEN | **OPEN** |
| U3 | Incompatibility followed by crash | `Start` publishes handshake state after `Terminate`; observer can publish `absent` or `crashed` before or after; no terminal latch. | OPEN | **OPEN** |
| U4 | Zero exit after fatal protocol violation | Violation calls `failAll` and `Kill`; observer publishes `absent` while `stopping`; `protocol_violation` remains only in the first read failure; pending request can publish `degraded`. | OPEN | **OPEN** |
| U5 | Nonzero exit after deliberate kill | `Terminate` sets `stopping`; nonzero exit observed in that state publishes provider `absent`; action intent/result/causation is not recorded. | OPEN | **OPEN** |
| U6 | Provider `crashed` while process alive | Synthesis send failure publishes `TTSStateCrashed` before any exit; `SetState` can keep process state `alive`; `Status` can report not alive from provider state. | OPEN | **PARTIALLY RESOLVED** |
| U7 | Terminal published before termination timeout | No terminal latch; `Terminate` overwrites with `stopping`; later exit can publish `absent` or `crashed`; timeout may leave ownership. | OPEN | **OPEN** |
| U8 | Two handshake projections | Helper-declared incompatibility maps to `incompatible`; client-detected missing mandatory capability maps through `protocol_violation` to `unavailable`. | OPEN | **PARTIALLY RESOLVED** |

No U case is fully resolved at the evidence level.

---

## 5. U1 — Orderly EOF racing reader crash

### 5.1 Scenario definition

A live, handshaken helper is being shut down. `Shutdown` sends the `shutdown` request while the reader can observe the helper closing its protocol stream. The question is whether the resulting EOF is an orderly consequence of shutdown or an independent reader/transport failure, and which public provider state, if any, is justified.

### 5.2 Established facts

- `[Source-derived]` `Shutdown` calls `call("shutdown")`, then, when the error is `helper_crashed` or `timeout`, falls through to `Terminate`; it clears `started` and `transport` afterward (`client.go:645–674`).
- `[Source-derived]` A reader receive error whose code is `internal_error` is reclassified to `helper_crashed` (`client.go:297–303`). The reader calls `failAll`; it kills only for `invalid_frame`, `frame_too_large`, or `protocol_violation` (`client.go:304–306`).
- `[Source-derived]` `Terminate` sets `providerState = stopping` when a process is still owned (`supervisor.go:253–259`).
- `[Source-derived]` `observeExit` publishes provider `absent` when it sees `stopping`, including for a nonzero exit; otherwise it publishes `crashed` (`supervisor.go:176–190`).
- `[Source-derived]` No stream-close provenance distinguishes natural EOF, shutdown-induced EOF, forced close, or read error. `waitForReader` discards whether the reader joined or timed out (`client.go:349–362`).
- `[Source-derived]` `failAll` records only the first read failure (`client.go:369–386`). `SetState` can overwrite `lastFailure` (`supervisor.go:233–242`).

### 5.3 Unestablished facts

- Which writer — the reader, `Shutdown`, or the exit observer — acquires the supervisor state first.
- Whether the EOF was caused by the shutdown action, by an independent helper failure, or by a transport/pipe condition.
- Whether orderly shutdown was established. `stopping` is an intent/state label, not proof of orderly shutdown.
- The actual exit status and the identity of the process instance associated with the EOF.
- Whether the reader was still active when `Shutdown` cleared the client transport.
- The final public state in a real execution.

### 5.4 Applicable decisions

- **OD-4:** confirmed exit without orderly shutdown → `unavailable` unless a more specific justified prior failure applies; orderly shutdown → `absent`; action causation is internal under OD-8.
- **OD-5:** `crashed` requires an observed process exit with no established orderly shutdown.
- **OD-3:** any first terminal decision is immutable and late evidence is retained.
- **OD-2:** observations must be bound to the process generation.
- **OD-8:** action-before-exit association is insufficient for causal attribution.

### 5.5 Current P17 handling

P17 lists U1 as OPEN at line 479. Its S7/S8 and P4 rules distinguish orderly shutdown from non-orderly exit, but they do not say whether a shutdown-induced EOF counts as an orderly consequence, nor do they define stream-close provenance. P17 supplies no unique outcome for the race.

### 5.6 Counterexamples and unresolved orders

1. Observer runs before `Terminate`: it sees a non-`stopping` state and publishes `crashed`; a later handshake/terminal write may then publish something else.
2. Observer runs after `Terminate`: it sees `stopping` and publishes `absent`, including for a nonzero exit.
3. `Terminate` times out and no exit is observed: provider can remain `stopping` while the client clears `started`/`transport`.
4. The reader’s `helper_crashed` becomes the first read failure while the provider is later published `absent`; the two evidence records describe different conditions.
5. The shutdown response is received before the EOF: the client can complete the action and still have no evidence that the EOF was action-caused.

### 5.7 Evidence sufficiency and classification

The source establishes a bounded race, not a unique outcome. The decisive questions — whether a shutdown-induced EOF is orderly and whether the fatal observation or the exit controls — require an owner rule, and the actual winner requires runtime evidence. Static source reading cannot settle either.

- **Classification:** **Owner decision required** (with runtime interleaving evidence also necessary).
- **Conclusion:** **OPEN**.
- **Reason:** the source has both outcomes and no provenance; P17 has no rule for the specific EOF case.

### 5.8 Smallest next step

A bounded read-only trace of the helper-side shutdown/stream-close behavior, if available, plus an owner rule for whether a shutdown-induced EOF is an orderly consequence. Runtime observation would still be required to establish the actual interleaving; no runtime work is authorized here.

---

## 6. U2 — Shutdown timeout plus protocol violation

### 6.1 Scenario definition

A shutdown action is in flight while the reader observes a protocol violation. Two termination paths and two failure/action outcomes can overlap, with no defined winner.

### 6.2 Established facts

- `[Source-derived]` `Shutdown` sends `shutdown` and waits for a terminal (`client.go:645–652`).
- `[Source-derived]` A dispatch/read violation calls `protocolViolation`, which calls `failAll` and `Kill` (`client.go:364–367`; `client.go:321–347`).
- `[Source-derived]` `Kill` is `Terminate` with the default grace (`supervisor.go:291–293`); `Terminate` sets `stopping`, closes stdin, and waits (`supervisor.go:253–286`).
- `[Source-derived]` A pending shutdown request can receive the protocol failure; `Shutdown` can call `Terminate` again and return (`client.go:653–655`), or on a timeout path it can call `Terminate` and then clear `started` and `transport` (`client.go:668–674`).
- `[Source-derived]` Two `Terminate` calls can overlap; one can return `nil` and the other can return `timeout` after the second grace. The timeout is not stored by the supervisor.
- `[Source-derived]` If an exit is observed while `stopping`, provider state becomes `absent`, including a nonzero exit (`supervisor.go:176–182`).
- `[Source-derived]` `readFailure` keeps the first failure; `lastFailure` can remain stale; no action record identifies the winner.

### 6.3 Unestablished facts

- Whether the protocol violation or the shutdown action was observed first.
- Whether either `Terminate` completed or timed out.
- Whether an exit occurred, and if so, its status and relationship to the kill.
- Which failure is authoritative when the first read failure, the shutdown outcome, and the exit fact disagree.
- Whether the process remains owned after a timeout; no quarantine record exists.
- The correct public state when no exit is observed.

### 6.4 Applicable decisions

- **OD-3:** first terminal decision is immutable; later evidence is retained.
- **OD-4:** no orderly shutdown → `unavailable` unless a justified prior failure applies; timeout is not a public state.
- **OD-5:** `crashed` requires an observed exit with no established orderly shutdown.
- **OD-7:** quarantine is internal-only ownership/reaping uncertainty, not a provider state; timeout alone is not quarantine.
- **OD-8:** ordering and mechanism evidence are required for causation.

### 6.5 Current P17 handling

P17 lists U2 as OPEN at line 480 and describes “two fatal/action outcomes with no defined winner.” P17 has no arbitration rule for a pre-existing protocol failure versus a shutdown action, no action-result record, and no public projection for a no-exit timeout. `quarantined` is internal-only, so it cannot resolve the public gap.

### 6.6 Counterexamples and unresolved orders

1. Protocol violation first, then `Shutdown` calls a second `Terminate`; one returns `nil`, the other `timeout`.
2. Shutdown first, then the protocol violation fails all pending requests; the shutdown terminal may never arrive.
3. Exit arrives while `stopping`; provider `absent` is published despite the nonzero/fatal history.
4. No exit arrives; provider remains `stopping` while the client has cleared transport.
5. Process remains owned after a second grace; no quarantine marker records ownership uncertainty.

### 6.7 Evidence sufficiency and classification

The source establishes overlapping paths and an absence of arbitration evidence. No source-only rule selects a public winner. A policy choice is required for fatal-versus-action precedence and for the no-exit timeout; runtime evidence is needed for the actual interleaving.

- **Classification:** **Owner decision required** (with runtime evidence).
- **Conclusion:** **OPEN**.
- **Reason:** the source neither records action order/result nor defines a timeout projection; P17 explicitly leaves the winner open.

### 6.8 Smallest next step

Define the owner rule for a pre-existing fatal observation versus a shutdown action and for a no-exit timeout. If a runtime path is later authorized, capture action start/completion, exit fact, and observation order for the same process instance.

---

## 7. U3 — Incompatibility followed by crash

### 7.1 Scenario definition

A handshake fails with a compatibility result such as `incompatible` or `unavailable`; the helper process later exits or crashes. The question is whether the handshake result or the later exit is the first terminal, and whether the later exit may replace it.

### 7.2 Established facts

- `[Source-derived]` On handshake failure, `Start` calls `Terminate`, waits for the reader without preserving the result, publishes `incompatible` for `protocol_incompatible`/`protocol_minor_unsupported` or `unavailable` otherwise, then clears `started`/`transport` (`client.go:248–258`).
- `[Source-derived]` The exit observer can run before, during, or after that `SetState`. If it sees `stopping`, it publishes `absent`; if it sees a non-`stopping` state, it publishes `crashed` (`supervisor.go:176–190`).
- `[Source-derived]` The later `SetState` can overwrite an earlier `absent`/`crashed`; a later observer can overwrite the handshake result. There is no terminal latch (`supervisor.go:233–242`).
- `[Source-derived]` `Terminate` result and reader-join result are discarded (`client.go:248–249`; `waitForReader` at `client.go:349–362`).
- `[Source-derived]` The typed handshake failure is retained in `lastFailure` when `SetState` publishes it; a later observer may overwrite `lastFailure` only in its unexpected branch (`supervisor.go:183–190`).

### 7.3 Unestablished facts

- Whether the handshake result or the exit was observed first.
- Whether `Terminate` completed, timed out, or left a live process.
- Whether the handshake terminal was published before the exit and therefore should be immutable under OD-3.
- Whether a later exit is new terminal evidence, late evidence about the same instance, or evidence of a different process generation.
- The final public state in a real execution.

### 7.4 Applicable decisions

- **OD-3:** one immutable first terminal decision per process instance; later observations retained but not republished.
- **OD-5:** `crashed` requires an observed exit without established orderly shutdown.
- **OD-6:** compatibility classification boundaries; U8 remains open.
- **OD-2:** generation binding is required for a replacement process.
- **OD-8:** later exit evidence does not by itself establish that the handshake failure caused the exit.

### 7.5 Current P17 handling

P17 lists U3 as OPEN at line 481. Its handshake rule says “fatal handshake + any exit → `incompatible`/`unavailable`,” which reads like a first-terminal policy, but P17 defines no latch, publisher, linearization point, or late-evidence rule. The current source has no terminal latch and can overwrite the handshake result.

### 7.6 Counterexamples and unresolved orders

1. Observer sees `stopping` and publishes `absent`; `Start` later publishes `incompatible`; the handshake result wins by being later, not by being first.
2. `Start` publishes `incompatible`; process later exits while not stopping; observer publishes `crashed`; the crash overwrites the compatibility result.
3. `Terminate` times out; `Start` publishes `incompatible` while the process remains alive; a later exit publishes `crashed`/`absent`.
4. A replacement process starts for the same client session; a late observer/exit from the old instance could modify the replacement if generation identity is absent.

### 7.7 Evidence sufficiency and classification

The source establishes the overwrite race and the absence of a latch. OD-3 supplies the governing principle, but the source does not implement it and no runtime order establishes which terminal is first. This is not fully source-resolvable from static text.

- **Classification:** **Partially source-resolvable; additional bounded evidence needed**.
- **Conclusion:** **OPEN**.
- **Reason:** OD-3's principle is accepted, but terminal identity, linearization, and the actual interleaving are not established in the source.

### 7.8 Smallest next step

A read-only design trace of the proposed terminal latch and generation boundary, followed by a runtime observation of the handshake-failure/exit interleaving if separately authorized. No owner policy choice is needed for the immutability principle itself; the missing element is implementation and order evidence.

---

## 8. U4 — Zero exit after a fatal protocol violation

### 8.1 Scenario definition

A fatal protocol violation is observed. The process is killed or exits with status zero. The question is whether the fatal observation, the zero exit, or the orderly/non-orderly shutdown boundary controls the public result.

### 8.2 Established facts

- `[Source-derived]` `dispatch` and `readLoop` turn frame, duplicate-terminal, unknown-ID, and decode violations into fatal `protocol_violation` evidence and call `protocolViolation`/`Kill` (`client.go:321–347`, `client.go:309–313`, `client.go:364–367`).
- `[Source-derived]` `Terminate` sets `stopping` and closes stdin (`supervisor.go:253–264`).
- `[Source-derived]` A zero exit observed while `stopping` publishes `processState = absent`, `providerState = absent` (`supervisor.go:176–182`). The `protocol_violation` remains in `readFailure` but is not published by `failAll`.
- `[Source-derived]` A pending synthesis that receives the fatal failure can publish `degraded` (`client.go:790–797`).
- `[Source-derived]` The source never publishes `crashed + protocol_violation` on this path.
- `[Source-derived]` `waitErr` retains `nil` for a zero exit; no action intent/result or event order is stored.

### 8.3 Unestablished facts

- Whether the zero exit means the process shut down cleanly, was killed, or exited for an independent reason.
- Whether the fatal observation preceded the exit as an event, as opposed to merely preceding it in observation order.
- Whether the fatal observation is a failure classification, a public failure value, a terminal, or only diagnostic evidence.
- Which public projection controls when a pending request publishes `degraded` and the provider publishes `absent`.
- Whether a zero exit after a fatal observation may ever be `crashed`.

### 8.4 Applicable decisions

- **OD-3:** first terminal immutability and late-evidence retention.
- **OD-4:** orderly versus non-orderly exit; confirmed exit without orderly shutdown → `unavailable` unless a prior failure is justified; orderly shutdown → `absent`.
- **OD-5:** `crashed` requires an observed exit with no established orderly shutdown; a fatal observation alone is not a crash.
- **OD-8:** fatal-before-exit association is not causal proof.

### 8.5 Current P17 handling

P17 lists U4 as OPEN at line 482. Its rule 2 maps a fatal protocol violation to `crashed + protocol_violation` regardless of exit status, while S3 says fatal is not synonymous with crash and rule 4 maps zero exit plus completed shutdown to `absent`. The proposal therefore contains competing routes and no unique result.

### 8.6 Counterexamples and unresolved orders

1. Violation first, `Kill` sets `stopping`, zero exit arrives: provider `absent`; pending request may publish `degraded`.
2. Violation fails a pending request before termination: `degraded` is published while the lifecycle path later publishes `absent`.
3. Exit arrives before `Kill` sets `stopping`: observer publishes `crashed`; no orderly shutdown is established.
4. `Terminate` times out; provider remains `stopping` while the fatal evidence remains in `readFailure`.
5. A zero exit may be action-caused, independent, or a clean acknowledgment; the source does not distinguish them.

### 8.7 Evidence sufficiency and classification

The source establishes the competing paths and the `absent`/`degraded` possibilities, but not a unique outcome. P17's own rules conflict. The choice between fatal-wins, exit-wins, and an explicit unknown rule is a governance decision.

- **Classification:** **Owner decision required** (with runtime/action evidence).
- **Conclusion:** **OPEN**.
- **Reason:** the model conflict is visible, but the event order, exit cause, and fatal-versus-exit policy are not established.

### 8.8 Smallest next step

Resolve the S3/rule-2 contradiction explicitly: decide whether a zero exit after a fatal observation can be `crashed`, `absent`, or an internal unknown. If the owner requires a deterministic outcome, obtain action-result and exit-order evidence for the same process instance.

---

## 9. U5 — Nonzero exit after a deliberate kill

### 9.1 Scenario definition

A deliberate `Kill` or `Terminate` is issued. The process exits with a nonzero status. The question is whether the exit is a crash, an action-caused termination, or an unresolved causal case, and what public state is justified.

### 9.2 Established facts

- `[Source-derived]` `Terminate` sets `stopping`, closes stdin, waits, kills on grace expiry, and can return `timeout` after a second grace (`supervisor.go:246–288`).
- `[Source-derived]` A nonzero exit observed while `stopping` sets `processState = exited` and publishes provider `absent` (`supervisor.go:176–182`).
- `[Source-derived]` `waitErr` retains the exit error, but the error is not separated into exit code and cause (`supervisor.go:174`).
- `[Source-derived]` The action's intent, grace, start, completion, and result are absent. `Terminate` results are ignored at multiple call sites (`client.go:248,654,661,664,861`).
- `[Source-derived]` The supervisor does not record `terminated_by_action`.
- `[Source-derived]` If the exit is observed before `stopping` is set, the observer publishes `crashed` (`supervisor.go:183–190`).

### 9.3 Unestablished facts

- Whether the deliberate action caused the exit, merely preceded it, or raced an independent exit.
- Whether the nonzero status is kill-consistent, helper-generated, or unrelated.
- Whether termination completed or timed out.
- Whether the process instance is the one the action targeted.
- The public projection for an action-caused nonzero exit when orderly shutdown is not established.
- Whether `absent` is justified or is only the current source's chosen value.

### 9.4 Applicable decisions

- **OD-4:** confirmed exit without orderly shutdown → `unavailable` unless a more specific justified prior failure applies; `terminated_by_action` is internal and has no public projection.
- **OD-5:** `crashed` requires an observed exit with no established orderly shutdown; action attribution is separate.
- **OD-8:** A–E evidence is required for deliberate-action causation; a successful action result or temporal order is insufficient.
- **OD-3:** first terminal immutability.

### 9.5 Current P17 handling

P17 lists U5 as OPEN at line 483. It retains `terminated_by_action` as an internal classification and leaves its public projection undecided. The source silently publishes `absent` for a nonzero exit observed while `stopping`; that is not a demonstrated application of OD-4 or OD-8.

### 9.6 Counterexamples and unresolved orders

1. Exit occurs before `Kill` sets `stopping`: observer publishes `crashed`, even if the action was requested.
2. `Kill` sets `stopping`; independent nonzero exit arrives: provider `absent`, no causal attribution.
3. `Kill` causes the exit after grace expiry: provider `absent`; the action result is ignored.
4. `Terminate` times out; a later exit arrives: `absent` or `crashed` depending on whether `stopping` is still set.
5. Action completed successfully but the exit was already in progress: `absent` falsely suggests a clean outcome.

### 9.7 Evidence sufficiency and classification

The source establishes current behavior but not causation. OD-8 supplies the required evidence standard and OD-4 supplies the public policy boundary, but the source lacks same-instance action/exit evidence, event order, mechanism operation, and competing-explanation assessment.

- **Classification:** **Owner decision required** for the public projection of action-caused/unknown nonzero exit; runtime evidence is also required for causation.
- **Conclusion:** **OPEN**.
- **Reason:** the current `absent` is an undocumented choice, not a source-justified application of OD-4/OD-8.

### 9.8 Smallest next step

Decide the public projection for an action-caused nonzero exit, or adopt an explicit internal unknown disposition. If deterministic attribution is required, capture action identity/outcome, exit fact/status, event order, and mechanism evidence for the same process instance.

---

## 10. U6 — Provider `crashed` while process remains alive

### 10.1 Scenario definition

A synthesis write fails. The client publishes `crashed` while the process may still be alive. The question is whether `crashed` may describe a provider condition without an observed process exit, and what must replace that publication if not.

### 10.2 Established facts

- `[Source-derived]` On `transport.Send` failure, `Synthesize` publishes `TTSStateCrashed` with `helper_crashed` before `cleanupRejected` (`client.go:755–758`).
- `[Source-derived]` `cleanupRejected` checks `ProcessAlive()` and may call `Terminate`, but ignores its result (`client.go:859–864`).
- `[Source-derived]` `SetState` sets `processState = alive` when a process is owned and the new state is not `stopping` (`supervisor.go:239–240`). The `crashed` publication therefore does not itself establish a process exit.
- `[Source-derived]` `Status` derives `ProcessAlive` from both process and provider state (`supervisor.go:324`), and `processAlive` returns false for `crashed`/`unavailable` (`supervisor.go:351–356`); a live OS process can be reported as not alive.
- `[Source-derived]` A later exit can publish `absent` if `stopping` is set, or `crashed` again (`supervisor.go:176–190`).
- `[Source-derived]` P17 O18 and §7 rule 3 preserve the fatal-transport-to-`crashed` mapping; P17 lists U6 as OPEN at line 484.

### 10.3 Unestablished facts

- Whether the process was alive at the exact publication point in a real execution.
- Whether `Terminate` after the publication succeeded, timed out, or left the process owned.
- Whether the write failure is a process condition, a transport condition, or both.
- Which public state is required for a live process whose transport failed.

### 10.4 Applicable decisions

- **OD-5:** accepted `crashed` meaning is an observed process exit with no established orderly shutdown; a live process must not be represented as `crashed`; `ProcessAlive` must be process-based, not provider-state-derived.
- **OD-4:** process exit is distinct from provider/runtime failure; no public `crashed` merely from a fatal observation.
- **OD-7:** quarantine is internal-only; it does not turn a live process into `crashed`.
- **OD-3:** any terminal decision must be distinct from a transport observation.
- **OD-9:** U6 remains an approval case until resolved or covered by an explicit unknown rule.

### 10.5 Current P17 handling

P17 lists U6 as OPEN at line 484. P17 O18 says fatal transport may publish `crashed + helper_crashed` while the process may be alive, and §7 rule 3 does the same. That conflicts with accepted OD-5. The proposal has not reconciled the source write, the `Status` liveness derivation, or the public state.

### 10.6 Counterexamples and unresolved possibilities

1. Send fails while process is alive; `SetState(crashed)` occurs; `cleanupRejected` terminates successfully; observer later publishes `absent`; the transient `crashed` was never an exit.
2. Send fails; `Terminate` times out; process remains alive while provider is `crashed`; `ProcessAlive` reports false.
3. Send fails because the process already exited; observer later publishes `crashed`; the initial publication and the observed exit are distinct facts.
4. `Status` is read between publication and cleanup; the response conflates a live process with a `crashed` provider state.
5. `cleanupRejected` does not terminate because `ProcessAlive()` is false due to provider state, leaving a live process unowned by the cleanup rule.

### 10.7 Evidence sufficiency and classification

The source-level contradiction is established: `client.go:756` publishes `crashed` before any exit, and `Status` can report the process not alive. OD-5 establishes that this is non-conforming. What is not resolved is the required replacement projection and the implementation/source correction; those are not authorized here.

- **Classification:** **Source-resolvable** for the existence and non-conformance of the source write.
- **Conclusion:** **PARTIALLY RESOLVED**.
- **Reason:** the evidence establishes a concrete violation of OD-5; it does not implement or approve a replacement, and U6 remains OPEN in the register.

### 10.8 Smallest next step

A finding-level P17/source reconciliation: remove or narrow the fatal-transport-to-`crashed` rule and define the replacement projection consistent with OD-5 and OD-4. Any source change or runtime verification requires separate authorization.

---

## 11. U7 — Terminal published before termination timeout

### 11.1 Scenario definition

A terminal provider state has been published. A termination action is still outstanding and may time out, after which the process may remain owned or later exit. The question is whether the terminal is immutable and how a timeout/ownership uncertainty is represented.

### 11.2 Established facts

- `[Source-derived]` `SetState` writes a provider state and `lastFailure` without a terminal latch or revision record (`supervisor.go:233–242`).
- `[Source-derived]` A later `Terminate` sees the owned process and sets `providerState = stopping` (`supervisor.go:253–259`).
- `[Source-derived]` If termination times out after the second grace, it returns `timeout` while the process may remain owned (`supervisor.go:277–286`).
- `[Source-derived]` A later exit publishes `absent` if `stopping`, otherwise `crashed` (`supervisor.go:176–190`).
- `[Source-derived]` `Start` publishes `incompatible`/`unavailable` after `Terminate` returns, including after a timeout (`client.go:248–254`). `Shutdown` clears `started`/`transport` after a timeout (`client.go:668–674`).
- `[Source-derived]` No terminal timestamp, terminal flag, revision record, or quarantine state exists. The client `closing` flag is not a terminal latch.

### 11.3 Unestablished facts

- Which state is the first terminal and at what linearization point.
- Whether the terminal should survive the later `stopping`, `absent`, or `crashed` write.
- Whether a timeout means the process is still owned, safely retained, or quarantined.
- Whether a later exit is late evidence or a new terminal decision.
- Whether the client's cleanup after timeout is consistent with the supervisor's ownership.

### 11.4 Applicable decisions

- **OD-3:** first terminal is immutable; later evidence is retained; publisher, synchronization, and late-evidence reconciliation remain open.
- **OD-7:** quarantine is internal-only ownership/reaping uncertainty; timeout alone is not quarantine.
- **OD-4:** process exit and orderly shutdown are distinct; timeout is not a public state.
- **OD-5:** a timeout/ownership condition is not an observed exit and cannot establish `crashed`.

### 11.5 Current P17 handling

P17 lists U7 as OPEN at line 485. It has no terminal latch or revision rule. Its Case A is proposed and assumes an established terminal, but the source has no such latch; P17 does not define how a termination timeout relates to the terminal or to internal quarantine.

### 11.6 Counterexamples and unresolved orders

1. `SetState(incompatible)` then `Terminate` sets `stopping`; later timeout leaves `stopping`, replacing the terminal.
2. Terminal published; `Terminate` times out; later exit publishes `absent`; terminal is gone.
3. Terminal published; observer runs before `Terminate`; later `Terminate` overwrites with `stopping`.
4. Timeout leaves process owned; client clears transport; no quarantine record captures ownership.
5. A replacement process starts; late old-generation evidence can affect the replacement if identity is not enforced.

### 11.7 Evidence sufficiency and classification

The source establishes a mutable terminal and the absence of a latch. OD-3 supplies the immutability principle, but the source does not implement it, and no runtime timeout/exit order is observed. The exact owner decision is already supplied in principle; implementation and evidence remain.

- **Classification:** **Partially source-resolvable; additional bounded evidence needed**.
- **Conclusion:** **OPEN**.
- **Reason:** the conflict is source-established, but the terminal linearization and timeout ownership behavior are not.

### 11.8 Smallest next step

A read-only design trace of the proposed terminal latch and timeout/quarantine relationship, followed by runtime evidence if separately authorized. No new public state is needed; quarantine must remain internal under OD-7.

---

## 12. U8 — Two handshake projections

### 12.1 Scenario definition

A helper fails a handshake compatibility requirement. Depending on whether the helper declares incompatibility or the client detects a missing mandatory capability, the source publishes `incompatible` or `unavailable` for the same compatibility concept.

### 12.2 Established facts

- `[Source-derived]` A helper wire error with `protocol_incompatible` or `protocol_minor_unsupported` is preserved by `failureFromWire`; `Start` publishes `incompatible` (`client.go:250–251`).
- `[Source-derived]` `validateHello` returns `protocol_violation` for a missing mandatory service capability, with `Detail["missing"]` (`protocol.go:527–531`).
- `[Source-derived]` `Start` maps every handshake failure other than the two incompatibility codes to `unavailable` (`client.go:252–254`).
- `[Source-derived]` Both paths are deterministic; this scenario does not require a race.
- `[Source-derived]` The failure code and `missing` detail are retained in `lastFailure` and the returned error; there is no compatibility classification field.

### 12.3 Unestablished facts

- Whether the missing mandatory capability is a structurally valid, semantically identified mandatory compatibility failure or a protocol/schema violation under OD-6's accepted categories.
- Whether the current client-detected path is non-conforming or merely a different internal code that may still map to `incompatible`.
- Whether P17's §7 rule 1 should be corrected to match OD-6.
- Whether helper-declared and client-detected instances of the same mandatory failure must share the public state in every case.

### 12.4 Applicable decisions

- **OD-6:** missing mandatory service capability, unknown required core capability, and known required capability at an unsupported required version → `incompatible`; duplicate identities, invalid semantic descriptors, and malformed helper version/limits mismatch → `unavailable`; helper-declared and client-detected instances of the same established mandatory failure share the public `incompatible` state.
- **OD-6 limitations:** misplaced descriptors, unsupported extension versions, unknown optional/provider feature requiredness, and unlisted classifications remain unclassified. Exact contract definitions remain open.
- **OD-3:** a later exit can affect the handshake terminal as in U3.
- **OD-9:** U8 remains an approval case; the register explicitly says OD-6 acceptance does not close U8.

### 12.5 Current P17 handling

P17 lists U8 as OPEN at line 486 and §7 rule 1 preserves the split: helper-declared incompatibility → `incompatible`; handshake `protocol_violation` → `unavailable`. P17 does not apply OD-6's accepted missing-mandatory-capability category to the client-detected path.

### 12.6 Counterexamples and unresolved possibilities

1. Helper declares `protocol_incompatible`: source publishes `incompatible`; conforms to the helper-declared path.
2. Helper omits a mandatory service capability: `validateHello` returns `protocol_violation`; `Start` publishes `unavailable`; OD-6 says the same mandatory failure should be `incompatible`.
3. Helper omits a mandatory capability but the helper also declares an incompatibility code: the wire code controls, potentially masking the missing capability.
4. Helper omits a capability that is outside the mandatory list: remains unclassified under OD-6; source may still return `unavailable`.
5. A later exit overwrites either handshake projection as in U3.

### 12.7 Evidence sufficiency and classification

The two paths and the policy conflict are source-established. OD-6 supplies the governing categories for the listed mandatory failures, but the register expressly leaves U8 open and the source still implements the split. The remaining work is P17 reconciliation and source/implementation alignment, not a new source discovery.

- **Classification:** **Partially source-resolvable; additional bounded evidence needed**.
- **Conclusion:** **PARTIALLY RESOLVED**.
- **Reason:** the source split and the OD-6 conflict are established; U8 remains OPEN because the source and P17 have not been reconciled or verified.

### 12.8 Smallest next step

Apply OD-6's categories to the nine `validateHello` exits in a finding-level mapping, then reconcile P17 §7 rule 1. No new public state is required. Source changes and runtime verification remain separate and unauthorized.

---

## 13. Consolidated dependencies

| Dependency | Affected cases | Why it blocks |
|---|---|---|
| No stable process-generation identity with every failure/action/exit (OD-2) | U1–U7 | The source cannot prove that an observation, action, and exit belong to the same process instance. |
| No immutable first-terminal latch or linearization (OD-3) | U1, U2, U3, U7 | Later writes can overwrite an earlier terminal; no evidence establishes which state was first. |
| No stream-close provenance or action-completion record | U1, U2, U4, U5, U7 | The source cannot distinguish natural EOF, shutdown-induced EOF, forced close, kill, and independent exit. |
| No action intent/result record (E4, E8) | U2, U5, U6, U7 | The source cannot establish what action occurred, whether it completed, or whether it caused the exit. |
| No exit code/cause separation from `waitErr` | U4, U5 | The source cannot distinguish zero/nonzero, kill-induced exit, and independent crash. |
| No terminal/quarantine representation (OD-7) | U2, U6, U7 | Timeout ownership and reaping uncertainty have no internal record; `quarantined` is not public. |
| `observeExit` maps `stopping` to `absent` even on nonzero exit | U1, U2, U4, U5, U7 | The source silently selects `absent` without proving orderly shutdown. |
| `Synthesize` publishes `crashed` before an observed exit | U6 | The source violates OD-5's process-exit meaning and `ProcessAlive` independence. |
| `validateHello` maps missing mandatory capability to `protocol_violation` | U8 | The source path conflicts with OD-6's accepted mandatory-capability classification. |
| P17 has no explicit unknown-case rule | U1–U8 | OD-9 requires a determined outcome or an explicitly defined unknown rule; none is adopted in P17. |
| E1–E9 evidence losses remain | U1–U7 | The source record cannot preserve action order, failure history, terminal evidence, or reaping outcomes. |

## 14. What remains OPEN

- **U1:** OPEN. Race and missing EOF provenance; owner rule and runtime order unresolved.
- **U2:** OPEN. No arbitration between shutdown action and protocol failure; no timeout projection.
- **U3:** OPEN. No terminal latch or established first-terminal order.
- **U4:** OPEN. Fatal-versus-exit and zero-exit policy unresolved; P17 rules conflict.
- **U5:** OPEN. Action causation and public projection unresolved; source `absent` is undocumented.
- **U6:** PARTIALLY RESOLVED. Source publishes `crashed` before exit, contrary to OD-5; replacement/remediation remains open.
- **U7:** OPEN. Terminal immutability is not implemented; timeout/ownership relation undefined.
- **U8:** PARTIALLY RESOLVED. Source split and OD-6 conflict established; P17/source reconciliation remains open.
- **CP-1:** OPEN; prior public `crashed` preservation is not decided.
- **Terminal-set definition:** OPEN.
- **N-5:** OPEN; no multi-classification selection rule.
- **N-6:** OPEN as an unresolved governance scope question; the C-B fallback contradiction is corrected.
- **S5-scope:** one OPEN tracking row; not a gate.
- **E1–E9:** present and unresolved.
- **U1–U8 approval:** all remain without a determined outcome or an explicitly defined unknown rule.

## 15. Recommended next governance step

Do not revise P17 or implement C1 from this assessment.

The smallest defensible next step is a finding-level governance decision that does three things in order:

1. Define the explicit unknown-case rule permitted by OD-9, including who records it and how it remains internal rather than becoming a public provider state.
2. Reconcile P17's U1–U8 table and projection rules with the accepted OD-3 through OD-8 boundaries, using this report's per-case evidence and not a forced unique outcome.
3. Authorize only the narrow source/implementation work that is actually required for U6 and U8, and the terminal/identity/action evidence required for U1–U5 and U7, under a separate authorization.

Source reading cannot decide the policy choices, and no runtime work is authorized here. P17 remains unapproved; contract rev3 remains frozen; C1 implementation remains unauthorized.

## 16. Integrity and limitations

- No existing artifact, source file, test, contract, manifest, or decision register was modified.
- This report is the only file created by this assessment.
- No tests, builds, benchmarks, application binaries, or runtime experiments were run.
- No missing evidence was regenerated or reconstructed.
- The source identity is the recorded helper content digest; the helper tree is untracked, so Git HEAD alone does not bind it.
- The U1–U8 outcomes are source-admissible paths, not observed executions. No claim is made about a particular runtime interleaving.
- The accepted assessment is used under the owner's working-copy determination, not as proof of historical provenance.
