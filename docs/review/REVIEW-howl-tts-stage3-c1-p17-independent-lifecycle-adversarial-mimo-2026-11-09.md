# Independent Adversarial Lifecycle Review — C1 / P17 (MIMO)

**Review date:** 2026-11-09

**Status:** COMPLETE — **FINAL VERDICT: FAIL — NOT READY.** P17 is not approved. Contract rev3 is not accepted. C1 implementation remains unauthorized. No source, test, contract, proposal, assessment, register, or existing review was modified.

**Subject:** independent, source-bound review of the TTS Stage 3 C1 lifecycle, P17 rev2, the U1–U8 scenarios, contract rev3, and the OD-1 through OD-9 decisions.

---

## 1. Verdict

**FAIL — NOT READY.**

The lifecycle design is not yet safe to implement as an approved contract. The failure is not that every source path is wrong; several P17 directions are correct and materially better than the current source. The failure is that the accepted approval gate is unmet and several concrete lifecycle mechanisms remain contradictory or ambiguous:

1. **OD-9 is not satisfied.** U1–U8 remain OPEN, and P17 supplies no explicit unknown-case disposition for all eight. The register's accepted gate says P17 is not approved while U1–U8 remain open without determined outcomes or explicitly defined unknown-case rules.
2. **No immutable first terminal exists in the source, and the terminal set remains undefined in P17.** Later writes overwrite earlier outcomes.
3. **No process-generation identity is attached to evidence or client operations.** A reader from an old process can fail pending work or kill a replacement process.
4. **`Shutdown` can start a replacement helper after the old helper has crashed.** This is a concrete ownership defect, not a theoretical race.
5. **Contract rev3 A6 contradicts accepted OD-5 and OD-4** by projecting `crashed` from a fatal observation even when the exit was zero. P17 explicitly removed that rule.
6. **The primary U1–U8 assessment is not fully bound to P17 rev2.** Its U3, U4, U5, and U6 rationales quote the frozen original P17's precedence table, not the candidate under review. The cases remain open, but the stated reasons are partly stale.

The smallest defensible next step is not C1 implementation. It is a bounded lifecycle correction pass that first defines the terminal set and unknown-case rule, then the generation/ownership boundary, then action/exit and shutdown semantics, followed by a fresh independent review bound to the resulting exact artifact.

---

## 2. Inspected identities and provenance

### 2.1 Documents

| Artifact | Lines | SHA-256 | Use |
|---|---:|---|---|
| `docs/review/ASSESSMENT-howl-tts-stage3-c1-p17-u1-u8-scenario-findings-2026-11-09.md` | 640 | `33c776bf50b6a10467c319d61b7d0079815705b33e050991140d5f5611b71e95` | Primary U1–U8 assessment; independently checked, not treated as authority |
| `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-rev2-2026-11-09.md` | 538 | `651d38d3e2d97f981a606a39ea4a3117c2e748ff9f5941eb8bd9671b6e022f82` | P17 review candidate |
| `docs/review/REVIEW-howl-tts-stage3-c1-p17-rev2-adversarial-651-2026-11-09.md` | 224 | `cbed00a4f9fb25153d2ce7555c5aff33cccf57dc103b5e016b15f3a5271e9196` | Prior adversarial review; no blocking B-2 finding; not approval |
| `docs/review/DECISION-REGISTER-howl-tts-stage3-c1-p17-evidence-2026-11-09.md` | 434 | `3c64797bfdd48ddacca95e6a0def88d960e4104aa3fc3034a25dc539d97cbc99` | Accepted OD-1 through OD-9 decisions |
| `docs/review/CONTRACT-howl-tts-stage3-c1-installation-lifecycle-rev3-2026-11-09.md` | 183 | `ed0da1a75d4f62874f873979e6667f66f19541a7523ff96369f2b44e8391195e` | Proposed concurrency contract |
| `docs/review/PROPOSAL-howl-tts-stage3-c1-evidence-semantics-p17-2026-11-09.md` | 166 | `01d887c18700d63070f1a2649659281ec1e875a21b49375adea12f8f11d917bf` | Frozen original P17; used only to check stale citations |

The accepted remaining-evidence assessment was also treated as a current working-copy input under the owner's determination. Its on-disk working-copy identity is 477 lines / `3c4ce2beb5681d8608e6e9e807b1b76562142b9f8adf2d8ea8847f6f7f2bbd93`. Its historical digest was not independently authenticated, and this review does not claim that it was.

### 2.2 Source identity

| Item | Value |
|---|---|
| Repository HEAD | `9c7d2a17dd7f702d7db6d95d9316e7640e34bf6f` (`9c7d2a1`) |
| Helper tree tracked by Git | **No**; all 19 helper Go files are untracked |
| Helper Go files | 19 |
| Helper Go lines | 4,992 |
| Helper content digest | `d0df8b2823389b305698fc06eda28e0e01f06259fdc69174cf7fb3851e271a51` |

The digest was reproduced by hashing the sorted `sha256sum` output of the 19 `*.go` files. It binds file paths and contents, not Git history, non-Go files, runtime state, or platform behavior. Because the tree is untracked, the Git HEAD pin does not establish its identity.

### 2.3 Method and boundaries

The review read the candidate and frozen original directly, checked every source claim used in the U-case analysis, and reconstructed the state-writer and lifecycle-call graph from `client.go`, `supervisor.go`, `protocol.go`, `errors.go`, and `internal/audio/types/provider.go`.

No tests, builds, benchmarks, application executions, runtime experiments, source changes, contract changes, register changes, or P17 edits were performed. No existing artifact was changed. All execution sequences below are source-reachable schedules, not observed runs.

**Evidence labels:**

- `[Document]` — directly stated in an inspected document.
- `[Observed source]` — directly visible in the inspected source text.
- `[Source-derived]` — a control-flow consequence of that source, not an observed execution.
- `[Unestablished]` — available evidence does not decide the proposition.
- `[Stale]` — a document claim displaced by the current candidate or accepted decision.

---

## 3. Independent source-trace summary

The lifecycle has one supervisor mutex, one client mutex, one client work mutex, and no generation token. The supervisor owns a single current `proc` pointer and one `procDone` channel. The client owns a transport pointer, reader, pending map, retired map, and first-read-failure pointer.

| ID | Source location | Observed fact |
|---|---|---|
| S1 | `client.go:191–270` | `Start` starts a process, assigns the transport, starts one reader, performs the hello handshake, and on failure calls `Terminate`, calls `waitForReader` without preserving whether it joined, publishes `incompatible`/`unavailable`, then clears `started` and `transport`. |
| S2 | `client.go:295–318` | On any receive error, `internal_error` is reclassified to `helper_crashed`; `failAll` runs; only `invalid_frame`, `frame_too_large`, or `protocol_violation` triggers `sup.Kill`. EOF has no provenance record. |
| S3 | `client.go:321–347` | Empty ID, duplicate terminal, and unknown ID become `protocol_violation`; a retired ID is silently dropped, not retained as evidence. |
| S4 | `client.go:349–362` | `waitForReader` returns the same value whether the reader joined or the timeout expired. |
| S5 | `client.go:364–386` | `protocolViolation` calls `failAll` then `sup.Kill`; `failAll` stores only the first read failure and delivers it to every pending request. |
| S6 | `client.go:398–443` | `callOn` reports generic timeout/cancellation without operation, phase, deadline origin, action identity, or event order; a send failure becomes `helper_crashed` without termination. |
| S7 | `client.go:645–674` | `Shutdown` calls `c.call("shutdown")` (which can call `Start`), then `Terminate`, then sets `closing=true`, clears `started` and `transport`, and returns. |
| S8 | `client.go:679–853` | `Synthesize` writes `busy`, `crashed` on send failure, `ready` on timeout/cancellation, and `degraded` on outcomes; several of these writes can occur after `cleanupRejected` has terminated the process. |
| S9 | `client.go:859–864` | `cleanupRejected` calls `Terminate` only when `ProcessAlive()` is true and either force is set or output exists; it ignores the termination result. |
| S10 | `supervisor.go:110–164` | `Start` returns any non-nil current process, sets `starting`/`handshaking`, and starts the single `observeExit`. There is no generation, no stale-reader join, and no check that the process is actually live. |
| S11 | `supervisor.go:166–195` | `observeExit` calls `Wait`, reads stderr, closes the process, and only then checks `s.proc == p`; if the provider state is `stopping`, it publishes `absent` for either zero or nonzero exit; otherwise it publishes `crashed`. |
| S12 | `supervisor.go:206–220` | `MarkHandshake` unconditionally publishes `processState=alive` and `providerState=unloaded`; it does not check that a process is still owned or that the handshake belongs to the current generation. |
| S13 | `supervisor.go:233–242` | `SetState` unconditionally overwrites provider state and `lastFailure`, and sets `processState=alive` whenever `s.proc != nil` and the new state is not `stopping`. |
| S14 | `supervisor.go:246–288` | `Terminate` sets `stopping`, closes stdin, waits for `done`, kills on cancellation or grace expiry, and can return `timeout` while the process remains owned. It records no intent, start, completion, or action outcome, and it does not guard against a second kill after the process has been reaped. |
| S15 | `supervisor.go:305–356` | `Status` derives `ProcessAlive` from process and provider state; a live OS process can be reported not alive when the provider state is `crashed` or `unavailable`. `ProcessAlive()` is a label/pointer check, not independent OS-liveness evidence. |
| S16 | `protocol.go:502–545` | `validateHello` returns `protocol_violation` for its validation failures; missing mandatory service capabilities is one such failure with a `missing` detail. |
| S17 | `protocol.go:566–578` | `missingCapabilities` computes missing mandatory IDs from the hello capability list. |
| S18 | `errors.go:15–21` | `Failure` carries code, message, retryable, field, and detail only; no generation, sequence, timestamp, phase, action, or exit status. |
| S19 | `internal/audio/types/provider.go:8–36` | Five process states and thirteen provider states; no `quarantined`, no generation, no terminal marker. |

The important structural facts are:

- The only identity in evidence is the client `session` string and the `Process` pointer/PID; the supervisor publishes no generation.
- The client can write supervisor state directly through `SetState` and `MarkHandshake`; OD-1's designated supervisor authority is not implemented as a publication boundary.
- The supervisor mutex prevents simultaneous map/field access, but it does not create an immutable decision, event order, generation boundary, or action identity.
- The client mutex protects maps and flags, not lifecycle closure. `Start` and `Shutdown` do not share a lifecycle admission lock.

---

## 4. U1–U8 findings

### U1 — Orderly EOF racing reader-crash

**Status: OPEN; blocks the OD-9 gate.**

`[Observed source]` A live handshaken helper is being shut down. `Shutdown` sends `shutdown` and then calls `Terminate` (`client.go:645–674`). Concurrently, the reader's `Receive` can return EOF; `readLoop` converts the resulting `internal_error` to `helper_crashed` (`client.go:297–303`). If `Terminate` has set `providerState=stopping`, the exit observer publishes `absent`, including a nonzero exit (`supervisor.go:176–182`); if it has not, the observer publishes `crashed` (`supervisor.go:183–190`).

There is no stream-close provenance. `waitForReader` discards whether the reader joined (`client.go:349–362`), and `readFailure` retains only the first failure (`client.go:373–376`). The source therefore admits both `absent` and `crashed` from the same user-visible shutdown action, depending on observation order.

**Counterexample:** shutdown response is received, then the helper closes stdout before `Terminate` writes `stopping`; the reader records `helper_crashed` and the observer records `crashed`. In the reverse order, the observer records `absent`, while `readFailure` still holds `helper_crashed`.

P17 lists U1 OPEN and supplies no rule for whether a shutdown-induced EOF is an orderly consequence, a fatal transport observation, or both (`P17:479`). The current P17's S7/S8 rules distinguish orderly from non-orderly exit but do not define stream-close provenance.

**Smallest resolution:** define an explicit `shutdown-induced EOF` classification and its provenance; or adopt an explicit internal unknown disposition that is not a public state. Require the action outcome and stream-close evidence to be retained. This is a missing lifecycle decision, not merely a test.

### U2 — Shutdown timeout plus protocol violation

**Status: OPEN; blocks the OD-9 gate.**

`[Observed source]` A protocol violation calls `failAll` and `Kill` (`client.go:364–367`), where `Kill` is `Terminate` with the default grace (`supervisor.go:291–293`). A concurrent `Shutdown` can call `Terminate` a second time (`client.go:653–674`). Both set `stopping`, close stdin, and wait on the same `done`; one can return `nil`, the other can return `timeout` after the second grace (`supervisor.go:271–286`). The timeout is returned to the caller, not stored.

The source retains no action start, completion, failure-vs-action order, or quarantine condition. If an exit arrives while `stopping`, the observer publishes `absent`, including a nonzero exit. If no exit arrives, the provider can remain `stopping` while the client clears its transport.

**Counterexample:** protocol violation wins the first `failAll`; `Kill` sets `stopping`; `Shutdown` calls `Terminate` again; the first returns after exit, the second returns `timeout`; `readFailure` still holds the protocol violation, `lastFailure` may be stale, and the public outcome depends on which write last reached the supervisor.

P17 lists U2 OPEN and provides no arbitration rule or timeout projection (`P17:480`). `quarantined` remains internal-only, so it cannot answer the public gap. U2 requires a determined precedence/unknown rule and an ownership outcome for a timeout with no exit.

### U3 — Incompatibility followed by crash

**Status: OPEN; blocks the OD-9 gate.**

`[Observed source]` On handshake failure, `Start` calls `Terminate`, calls `waitForReader` without preserving whether it joined, then publishes `incompatible` for `protocol_incompatible`/`protocol_minor_unsupported` or `unavailable` otherwise (`client.go:248–254`). The exit observer can run before, during, or after that `SetState`: it can publish `absent` if it saw `stopping`, or `crashed` otherwise (`supervisor.go:176–190`). A later `SetState` or `MarkHandshake` can overwrite the handshake result; there is no terminal latch or terminal set.

**Counterexample:** observer sees `stopping` and writes `absent`; `Start` then writes `incompatible`, so the handshake result wins by being later, not by being first. Alternatively, `Start` writes `incompatible`; the process exits while not stopping; the observer writes `crashed`, overwriting the compatibility result.

P17 lists U3 OPEN, has no terminal latch, and leaves the terminal set undefined (`P17:323`). P17's Case A is conditioned on an “established immutable terminal” but does not establish which provider values are terminal.

**Assessment correction:** the primary assessment says P17 “rule 1” says “fatal handshake + any exit → incompatible/unavailable.” That rule appears in the **frozen original P17** (`P17 original:112–118`), not in the 538-line rev2 candidate. Rev2 O15/O16 says apply P2 and leaves path classification open; P5 lists U3 open. The assessment's stated rationale is therefore stale, although U3 remains genuinely open.

### U4 — Zero exit after fatal protocol violation

**Status: OPEN; blocks the OD-9 gate.**

`[Observed source]` A protocol violation triggers `failAll` and `Kill` (`client.go:309–313`, `364–367`). `Terminate` sets `stopping` and closes stdin. If the helper exits zero while `stopping`, the observer publishes `absent` (`supervisor.go:176–182`). A pending synthesis that receives the fatal failure can publish `degraded` (`client.go:790–797`). The source never publishes `crashed + protocol_violation` on this path.

P17 rev2 correctly says a fatal protocol/transport observation does not establish a process exit and must not publish `crashed` from the observation alone (`P17:268–275`). It does not decide whether the protocol violation is a prior **failure classification**, whether it has a resolved publication status, or whether the zero exit is an orderly shutdown. P4's failure-classification threshold and the no-current-state case remain OPEN (`P17:289–307`). U4 therefore remains open for a different reason than the assessment gives.

**Assessment correction:** the primary assessment says P17 “rule 2” maps fatal protocol to `crashed + protocol_violation` and conflicts with S3. That is the **frozen original** precedence table, not rev2. Rev2 C-01 and P3 explicitly remove the fatal-to-crashed projection. The assessment's stated contradiction is not present in rev2; the unresolved issue is the classification/threshold and orderly-exit determination. P17's own §10 U4 row also retains a “fatal-wins rule may overstate certainty” rationale (`P17:482`), while P3/O7 says no fatal-to-crashed projection; that row should be corrected to the actual remaining threshold question.

### U5 — Nonzero exit after deliberate kill

**Status: OPEN; blocks the OD-9 gate.**

`[Observed source]` `Terminate` sets `stopping`, closes stdin, and may kill after grace (`supervisor.go:246–288`). A nonzero exit observed while `stopping` publishes provider `absent` (`supervisor.go:176–182`). The action's intent, target, start, completion, and result are not recorded; `Kill` is only an alias for `Terminate`; `waitErr` is not separated into exit status, signal, and cause. OD-8's A–E standard cannot be satisfied from this source.

**Counterexample:** independent nonzero exit occurs just before `Terminate` sets `stopping`; the observer publishes `crashed` even though a kill was requested. In the reverse order, the same underlying exit publishes `absent`, without evidence that the action caused it.

P17 correctly keeps `terminated_by_action` internal and leaves no public projection (`P17:397–428`; OD-4). U5 remains open because the causal evidence required by OD-8 is absent, not because a public projection is undecided. The assessment's statement that an owner decision is still needed for the public projection is **stale** after OD-4 acceptance.

### U6 — Provider `crashed` while process remains alive

**Status: OPEN; partial policy correction; source still non-conforming.**

`[Observed source]` On a synthesis send failure, `Synthesize` publishes `crashed` before `cleanupRejected` (`client.go:755–758`). `SetState` sets `processState=alive` when a process is owned (`supervisor.go:239–240`). `cleanupRejected` may terminate, but ignores the result (`client.go:859–864`). `Status.ProcessAlive` is derived from both process and provider state and reports false for `crashed` or `unavailable` (`supervisor.go:324`, `351–356`). A live process can therefore be represented as `crashed` and as not alive.

P17 rev2 O18 improves on the source: “A live process must not be `crashed`; retain the evidence. No exit is established; preserve the last justified current state; no `absent`.” If no current state is justified, it selects no projection and leaves the case OPEN (`P17:222`, `268–275`). This is consistent with OD-5. But U6 is still OPEN because the source has not been reconciled, `ProcessAlive` is not independent OS liveness, and the replacement projection for the no-current-state branch is C-01 OPEN.

**Assessment correction:** the assessment says P17 O18 and §7 rule 3 preserve the fatal-transport-to-`crashed` mapping. Rev2 O18 says the opposite, and §7 no longer contains that rule. The assessment's U6 source findings are correct; its description of the candidate is stale.

### U7 — Terminal published before termination timeout

**Status: OPEN; blocks the OD-9 gate.**

`[Observed source]` `SetState` has no terminal latch (`supervisor.go:233–242`). `Terminate` unconditionally sets `stopping` when a process is still owned (`supervisor.go:253–259`). If the process later exits, the observer publishes `absent` or `crashed` (`supervisor.go:176–190`). `Start` publishes `incompatible`/`unavailable` after `Terminate` returns, including after a timeout (`client.go:248–254`). `Shutdown` clears client state after a timeout (`client.go:668–674`).

**Counterexample:** `SetState(incompatible)`; `Terminate` sets `stopping`; timeout leaves the process owned; a later exit publishes `absent`; the original terminal is gone. No timestamp, revision, terminal flag, or quarantine record exists.

P17 requires an immutable first terminal but leaves the terminal set, synchronization, and late-evidence mechanics OPEN (`P17:132–142`, `323`). U7 therefore cannot be implemented deterministically from P17 alone. OD-7's quarantine concept is internal-only, but the timeout/quarantine transition is still undefined.

### U8 — Two handshake projections

**Status: OPEN; the policy conflict is determinable but the candidate and source are inconsistent.**

`[Observed source]` A helper-declared `protocol_incompatible` or `protocol_minor_unsupported` is preserved by `failureFromWire` and mapped to `incompatible` (`client.go:250–254`). A missing mandatory service capability is detected by `validateHello` as `protocol_violation` with a `missing` detail (`protocol.go:527–531`) and is mapped by `Start` to `unavailable` (`client.go:252–253`).

OD-6 accepts the mapping: missing mandatory service capability → `incompatible`; helper-declared and client-detected instances of the same established mandatory failure share the public state. P17 P2 also states “Missing mandatory service capability → incompatible” (`P17:258–266`). Yet P17 O16 immediately says the exact path-by-path classification of the nine `validateHello` exits remains OPEN and U8 remains OPEN (`P17:220`). That is an internal inconsistency: the source condition is not ambiguous for the missing-mandatory path, and OD-6 already selects its public category.

**Counterexample:** helper omits `howl.tts.catalog`; `validateHello` returns `protocol_violation` with `missing=howl.tts.catalog`; source publishes `unavailable`; OD-6 and P17 P2 require `incompatible`. This is a direct source-to-decision contradiction, not a race.

**Smallest resolution:** classify the missing-mandatory path from `validateHello` explicitly as the established compatibility failure, correct the source mapping, and narrow the “path-by-path OPEN” language to categories actually unclassified by OD-6. U8 should not remain open for a path that OD-6 already decides.

---

## 5. Broader lifecycle findings

### F-01 — BLOCKING — OD-9 is unmet

**Evidence:** `[Document]` `P17:475–491` says all U1–U8 remain OPEN; the register's accepted OD-9 says P17 is not approved while U1–U8 remain OPEN without determined outcomes or explicitly defined unknown-case rules. P17 has no rule assigning `outcome_unknown` to each OPEN case; it only permits the internal label in P7.

**Consequence:** P17 cannot be approved, and C1 cannot be authorized, regardless of the other source defects. This is the accepted gate, not a new prerequisite invented by this review.

**Smallest resolution:** for each U case, either determine an outcome or define an explicit internal unknown rule with its evidence and non-public status; update the status only through the authorized governance process.

### F-02 — BLOCKING — No immutable terminal and no terminal set

**Evidence:** `[Observed source]` `SetState` overwrites (`supervisor.go:233–242`); `Terminate` overwrites with `stopping` (`253–259`); `observeExit` overwrites with `absent`/`crashed` (`176–190`); `MarkHandshake` overwrites with `alive`/`unloaded` (`206–220`); `Start` publishes handshake states after `Terminate` (`client.go:248–254`). P17 requires first-terminal immutability but says the complete terminal set is OPEN (`P17:132–142`, `323`).

**Consequence:** U3, U7, CP-1, and the public projection of any earlier terminal remain indeterminate. The same observable execution can end in different public states depending on lock acquisition and goroutine scheduling.

**Smallest resolution:** define the terminal set and one linearization/publication point, implement a terminal latch with generation binding, and retain later evidence without republishing it. The latch mechanism can be an implementation detail; the terminal set and late-evidence projection rule cannot.

### F-03 — BLOCKING — No generation isolation; an old reader can affect a replacement

**Evidence:** `[Observed source]` `readLoop(p)` closes over `p` but calls `c.failAll(f)` and `c.sup.Kill(...)` with no generation check (`client.go:295–318`). `failAll` drains the client's **current** pending map (`369–386`). `retired` is reset on restart (`client.go:211`, `219`). `Start` can run after a failed handshake or after `ProcessAlive` returns false, and the supervisor has no generation token (`supervisor.go:110–164`). `waitForReader` may time out without joining (`client.go:349–362`).

**Concrete sequence:**

1. Start p1; handshake fails; `Terminate` times out and p1 remains alive; `waitForReader` times out; client clears `started`/`transport` (`client.go:248–258`).
2. If p1 is still owned, a later Start returns p1 and installs a **second reader** on the same process. If p1 later exits, the observer clears `s.proc` and a later Start can create p2 while the first reader has still not been joined.
3. The old p1 reader returns an error or a late frame. Because it has no p1 generation check, it calls `failAll` against the current pending map and, for a protocol error, `sup.Kill` against the current supervisor process, which is p1 or p2, not the reader's generation.
4. The current process is terminated, or its pending request receives p1's failure, even though the reader belongs to a different generation.

This is a source-reachable ownership violation of OD-2 and P17's constraint that old evidence must never modify a replacement. A late response with an old request ID is worse: `retired` was reset on restart, so `dispatch` treats the old ID as unknown and calls `protocolViolation`, which kills the current process (`client.go:339–345`, `364–367`). The failure can occur without any user action on the current generation.

**Smallest resolution:** attach generation identity to every reader, failure, action, and exit; make old-generation `failAll`/`Kill`/state writes no-ops for current state; join/close the old reader before installing a replacement; never signal or publish through a stale pointer.

### F-04 — BLOCKING — `Shutdown` can start a replacement helper

**Evidence:** `[Observed source]` `Shutdown` checks `currentTransport()` and then calls `c.call("shutdown")` (`client.go:645–652`). `c.call` calls `c.Start` (`client.go:388–396`). `Start` sees `started && transport != nil`, checks `sup.ProcessAlive()`, and if false, clears the client state and calls `sup.Start`, which starts a new process (`client.go:201–226`; `supervisor.go:110–164`).

**Concrete sequence:**

1. p1 exits unexpectedly; `observeExit` sets `s.proc=nil`, provider `crashed` (`supervisor.go:171–190`).
2. The client still holds `transport=p1` and `started=true`.
3. Caller invokes `Shutdown`; `currentTransport()` is non-nil.
4. `c.call` calls `Start`; `ProcessAlive()` is false; `Start` clears `readFailure`/`retired` and starts p2.
5. The shutdown frame is sent to p2 and p2 is terminated; the old `crashed` evidence is overwritten by p2's `starting`/`handshaking`/`unloaded`/`absent` writes.

A shutdown action must not create a new process. This is a concrete lifecycle/ownership defect, not a hypothetical race. It also invalidates any assumption that `Shutdown` is idempotent or that a crashed generation remains the current evidence target.

**Smallest resolution:** make shutdown a lifecycle-closing operation that refuses admission before it can call `Start`; distinguish stale transport from an owned live process; never restart during shutdown; preserve the old generation's terminal evidence and create no new generation.

### F-05 — BLOCKING — Contract rev3 A6 contradicts OD-5 and P17

**Evidence:** `[Document]` Contract rev3 A6 rule 3 says: “A fatal `failureObserved` projects `crashed` with that kind even when `exitStatus == 0`” (`CONTRACT:136–137`), and its table maps `0 | fatal | any` to `crashed + kind` (`CONTRACT:141–146`). OD-5 says `crashed` is an observed process exit where orderly shutdown is not established; P17 P3 says a fatal protocol/transport observation does not establish an exit and must not publish `crashed` (`P17:268–275`).

**Consequence:** the proposed concurrency contract and the proposed evidence model select different public states for the same zero-exit-after-fatal case. An implementer cannot satisfy both. U4 is not merely undecided; the two proposed artifacts conflict.

**Smallest resolution:** reconcile A6 with OD-5/P17 before using contract rev3 as an implementation basis. A fatal observation can be retained and can influence routing, but it cannot independently prove `crashed`; the exit fact and orderly-shutdown determination must control.

### F-06 — HIGH — `MarkHandshake` can resurrect a dead process

**Evidence:** `[Observed source]` `MarkHandshake` unconditionally sets `processState=alive` and `providerState=unloaded`, with no `s.proc` or generation check (`supervisor.go:206–220`). `Client.Start` calls it after the handshake response and after releasing `c.mu` (`client.go:261–270`).

**Concrete sequence:** p1 returns a valid hello; the process exits before `MarkHandshake` executes; `observeExit` clears `s.proc` and publishes `crashed`; `MarkHandshake` then publishes `alive`/`unloaded`. `Status` reports an alive/unloaded provider with no process. Even if `observeExit` runs afterward, the intermediate public state is false; if the observer already ran, the final public state is false indefinitely.

**Smallest resolution:** bind the handshake result to the exact process generation, require an owned live process before publication, and route late handshake evidence through the terminal/late-evidence rules.

### F-07 — HIGH — Post-termination `ready` and `degraded` writes can survive exit

**Evidence:** `[Observed source]` `Synthesize` timeout calls `cleanupRejected(..., true)`, which may terminate the process, and then unconditionally writes `ready` (`client.go:766–772`). The cancellation path does the same (`773–784`). Outcome failure and response failure paths call `cleanupRejected` and then write `degraded` or `ready` (`790–806`, `811–829`). `cleanupRejected` ignores the termination result (`859–864`).

**Concrete sequence:** p1 is alive; synthesis times out; `cleanupRejected` calls `Terminate`; p1 exits and the observer publishes `absent` and clears `s.proc`; `Synthesize` then writes `ready`. `Status` reports `State=ready`, `Ready=true`, `Available=true`, `ProcessAlive=false`. If `Terminate` times out instead, `SetState(ready)` sees a non-nil process and sets `processState=alive`, effectively undoing the `stopping` state and leaving no quarantine record.

**Consequence:** an active operational state can be published after the process is gone, or after a termination timeout. This contradicts OD-4's exit boundary, P17's no-active-state-survives-exit rule, and OD-7's timeout/quarantine distinction. It is not covered by U4/U6 alone: it also applies to request timeouts and cancellation.

**Smallest resolution:** make request cleanup return a typed outcome; do not write a public state after a confirmed exit or a termination timeout; route post-termination uncertainty through the terminal latch or internal `outcome_unknown`/quarantine. P17 needs an explicit rule for request-timeout-triggered termination, not only shutdown.

### F-08 — HIGH — `ProcessAlive` is not independent process liveness

**Evidence:** `[Observed source]` `Status.ProcessAlive` is computed from `processAlive(processState, providerState)` and returns false for `crashed`/`unavailable` (`supervisor.go:324`, `351–356`). `SetState` can mark `processState=alive` whenever a process pointer exists (`233–240`). `ProcessAlive()` checks only `s.proc != nil && processState != exited` (`335–341`).

**Consequence:** the source can both report a live process as not alive (provider `crashed`/`unavailable`) and report a dead/replaced process as alive (stale pointer + state write). OD-5 requires actual process liveness independent of provider state. U6 is not resolved by a semantic sentence alone; the liveness predicate must be implemented from ownership/handle evidence.

**Smallest resolution:** derive process liveness from the owned process handle/OS state, not from provider state; make `ProcessAlive` an explicit input to cleanup and publication; add a generation-bound exit fact.

### F-09 — HIGH — Old reader is never joined; cleanup can race replacement

**Evidence:** `[Observed source]` `Start` launches the reader and does not join it before returning (`client.go:230–239`). `Terminate` waits only for the process's `done` channel, not for the reader (`supervisor.go:246–288`). `Shutdown` clears the transport without waiting for the reader (`client.go:668–674`). `waitForReader` is used only on handshake failure and returns no join result (`client.go:349–362`).

**Consequence:** even after a normal process exit, the old reader can still process a late frame or error after a replacement is installed. That is the cross-generation path in F-03. Contract rev3 A2's drain-before-reap ordering is the correct direction, but the source does not implement it, and P17 does not require a reader-join result as a public or internal precondition.

**Smallest resolution:** one reader per generation, joined before retirement; no replacement while an old reader/transport remains; retain late evidence against its own generation rather than dropping it or applying it to the current one.

### F-10 — HIGH — `Terminate` has no reaped guard or lifecycle state machine

**Evidence:** `[Observed source]` `Terminate` calls `p.Kill()` on context cancellation or timer expiry (`supervisor.go:271–286`). It does not consult a reaped state; `p.Kill` calls `proc.KillTree(pid)` and `cmd.Process.Kill` (`supervisor.go:512–519`). Multiple `Terminate` calls can overlap. If the timer/cancellation branch wins while `Wait` is completing, the kill can target a process that has exited or is being reaped.

**Consequence:** on Windows, a recycled PID could target an unrelated process; on any platform, the action evidence and ownership state are wrong. Contract rev3 A1/I2 explicitly forbids kill-after-reap, but the source has no `pending → live → reaping → reaped` state.

**Smallest resolution:** implement the contract's single-owner lifecycle state machine, and make a second kill/kill-after-reap a typed no-op; retain the action outcome.

### F-11 — HIGH — `Shutdown` and `Start` have no common admission barrier

**Evidence:** `[Observed source]` `Start` uses `startMu`, checks `closing` before `sup.Start`, and assigns `c.transport` only afterward (`client.go:198–235`). `Shutdown` does not take `startMu`; it returns immediately when `currentTransport()` is nil (`client.go:645–651`). Therefore a concurrent `Start` that has not yet assigned its transport can be missed by `Shutdown`, and can complete after `Shutdown` has returned.

**Concrete sequence:** Start enters `sup.Start`; Shutdown sees no transport and returns nil; Start receives p1, assigns transport, starts reader, and handshakes; a live helper remains after “shutdown” completed. This is a lifecycle closure defect, not a U1 race.

**Smallest resolution:** make `Start` and `Shutdown` share a lifecycle admission lock and closing state; no start may commit after shutdown begins, and no shutdown may report success while a start is in flight.

### F-12 — MEDIUM — Cleanup completion is not retained and can be indistinguishable from intentional retention

**Evidence:** `[Observed source]` `cleanupRejected` calls `Terminate` only if `c.sup.ProcessAlive()` is true and either force is set or output exists (`client.go:859–864`). `ProcessAlive()` is the state/pointer predicate from F-08, not independent OS liveness. The termination result is ignored, and no cleanup outcome is recorded. When `force=false` and output does not exist, the helper may intentionally remain alive; the lifecycle record cannot distinguish that from failed cleanup.

**Consequence:** cleanup success is not established, and a live helper may remain without a retained ownership outcome. This is primarily an ownership/evidence gap; it does not by itself prove that every retained process is orphaned.

**Smallest resolution:** use actual owned-process liveness and record a typed cleanup outcome; distinguish intentional helper retention from termination failure. No deletion may be presented as successful while the record cannot establish what happened to an owned process.

### F-13 — MEDIUM — The primary U1–U8 assessment is stale on U3, U4, U5, and U6

**Evidence:** `[Document]`

- U3 says P17 “rule 1” says “fatal handshake + any exit.” That is the frozen original's precedence table (`P17 original:112–118`), not rev2.
- U4 says P17 “rule 2” maps fatal protocol to `crashed + protocol_violation` and conflicts with S3. Rev2 C-01/P3 removes the fatal-to-crashed projection; that rule is not in the 538-line candidate.
- U5 says the public projection for `terminated_by_action` is undecided. OD-4 accepted internal-only/no public projection.
- U6 says P17 O18 and §7 rule 3 preserve fatal transport to `crashed`; rev2 O18 says a live process must not be `crashed`.

**Consequence:** the assessment's U-case rationales are partly not bound to the candidate under review. U3/U4/U5/U6 remain open, but for narrower/different reasons. A reviewer or owner relying on those rationales could fix the wrong artifact or overstate the remaining gap.

**Smallest resolution:** rebind the U assessment to the exact P17 rev2 text and accepted register; retain the source findings but correct the proposal citations before using it as decision evidence.

### F-14 — MEDIUM — P17's public `absent` rule depends on unaccepted P4-A

**Evidence:** `[Document]` P17's status says P4-A and its non-application of retained S5 are **PROPOSED**, not independently accepted (`P17:5`, `P17:308`). OD-4 accepted the broad orderly-shutdown/no-orderly policy. P17's actual `absent` and branch routing depend on P4-A gates 1–5 (`P17:289–320`).

**Consequence:** the reviewer cannot treat P17's `absent` projection as fully accepted policy. P4-A may be conservative, but it is still a proposed gate that can convert `absent` into OPEN. If implementation follows P17's text, it implements a rule the owner has not accepted; if it follows OD-4 alone, it may disagree with P17.

**Smallest resolution:** either accept P4-A explicitly or state the exact public projection rule that applies where P4-A is unresolved. This is a public-state eligibility question, not a code-detail preference.

### F-15 — MEDIUM — `outcome_unknown` is permitted but not assigned

**Evidence:** `[Document]` P17 P7 permits `outcome_unknown` internally and says it is not public (`P17:430–436`). U1–U8 are all OPEN and no U case says “record outcome_unknown with this evidence and no public state.” The no-current-state branches say OPEN without requiring an internal record.

**Consequence:** OD-9's “explicitly defined unknown-case rule” is not met. “OPEN” is a review status, not an operational rule for the supervisor.

**Smallest resolution:** define, per case or as a general predicate, when `outcome_unknown` is recorded, who records it, what evidence it retains, and that it never becomes a public state.

### F-16 — MEDIUM — No action identity, outcome, or exit status in the source

**Evidence:** `[Observed source]` `Terminate` has no action parameter, no start/completion, and no outcome record (`supervisor.go:246–288`); `Kill` is an alias (`291–293`); `waitErr` is an error, not status/signal/cause (`supervisor.go:166–195`); `Failure` has no action or generation (`errors.go:15–21`). Contract rev3 A6 proposes action fields and an exit status, but they are not implemented and A6 itself conflicts with OD-5 (F-05).

**Consequence:** OD-8's A–E standard cannot be applied. U5 cannot be resolved by source inspection, and the same visible exit can be attributed differently depending on observation order.

**Smallest resolution:** record action identity/target/outcome and exit status/signal/observation context bound to the generation; apply OD-8's A–E and F assessment before any `terminated_by_action` record.

### F-17 — MEDIUM — No quarantine transition among timeout, retained ownership, and reaping

**Evidence:** `[Observed source]` `Terminate` can return `timeout` while `s.proc` remains owned (`supervisor.go:277–286`); `Shutdown` clears client state regardless (`client.go:668–674`). `quarantined` is internal-only under OD-7, but no record or transition exists.

**Consequence:** U2/U6/U7 and F-04/F-09 can leave a live or unknown process after the client reports a terminal/closed state. OD-7's “explicitly recorded unresolved ownership/reaping condition” is not implemented, and P17's exact transition is open.

**Smallest resolution:** record timeout, retained ownership, and reaping uncertainty separately; only an explicit retained/reaping record is quarantine; a later exit is not automatically successful reaping.

### F-18 — LOW — `retired` is unbounded, and status/catalog can start a helper

**Evidence:** `[Observed source]` `retire` inserts every request ID into `retired` without eviction (`client.go:445–469`). `ProviderStatus` and `ProviderCatalog` call `Start` (`client.go:952–970`), so an inspector read can start/handshake a helper and create a new lifecycle generation.

**Consequence:** long-running clients accumulate memory; a status read can mutate lifecycle state and overwrite a prior terminal. The lifecycle effect is more significant than the map growth.

**Smallest resolution:** bound/rotate retired IDs by generation and make status/catalog non-mutating or explicitly declared start-capable with terminal/generation rules.

---

## 6. Concrete execution sequences

The following sequences are source-reachable. They are not runtime observations.

### CE-1 — Old reader attacks a replacement

1. p1 starts; handshake fails; `sup.Terminate` times out; `waitForReader` times out; client clears `started`/`transport` (`client.go:248–258`).
2. If p1 remains owned, a later Start returns p1 and installs a second reader. If p1 later exits, a later Start can create p2 while p1's reader has not been joined.
3. p1's old reader returns a protocol error or late frame; it calls `failAll` and `sup.Kill` with no generation check (`client.go:295–318`).
4. The current pending request receives p1's failure, or the current process (p1 or p2) is killed. The public state can be `stopping`/`absent` even though the current generation never violated anything.

**Expected:** old-generation evidence cannot affect the replacement. **Actual:** it can. **Consequence:** false crash/action attribution and unsafe process termination.

### CE-2 — Shutdown starts a helper

1. p1 exits unexpectedly; `s.proc` is cleared and provider becomes `crashed`.
2. Client still has `started=true`, `transport=p1`.
3. `Shutdown` sees a transport, calls `c.call`, which calls `Start`.
4. `Start` sees `ProcessAlive()==false`, clears state, starts p2, and handshakes.
5. Shutdown terminates p2 and returns; p1's terminal evidence has been overwritten by p2's lifecycle writes.

**Expected:** shutdown closes the existing generation; it does not create one. **Actual:** it creates and terminates a new one. **Consequence:** lifecycle history and public state are false.

### CE-3 — Handshake resurrects a dead process

1. p1 returns a valid hello.
2. p1 exits before `MarkHandshake` runs; `observeExit` clears `s.proc` and publishes `crashed`.
3. `MarkHandshake` publishes `alive`/`unloaded` without checking the process.
4. `Status` reports alive/unloaded; no process exists.

**Expected:** late handshake evidence is retained but cannot resurrect a dead generation. **Actual:** it does.

### CE-4 — Timeout leaves `ready` after exit

1. p1 is alive; synthesis times out.
2. `cleanupRejected` calls `Terminate`; p1 exits; observer publishes `absent` and clears `s.proc`.
3. `Synthesize` writes `ready` after cleanup.
4. `Status` reports `ready`, `Ready=true`, `Available=true`, `ProcessAlive=false`.

**Expected:** no active state after a confirmed exit; unknown/timeout is internal, not `ready`. **Actual:** `ready` is public after exit.

### CE-5 — Shutdown misses a concurrent Start

1. `Start` enters `sup.Start` but has not assigned `transport`.
2. `Shutdown` sees `currentTransport()==nil` and returns nil.
3. `Start` receives p1, assigns transport, starts the reader, and handshakes.
4. A live helper remains after `Shutdown` returned successfully.

**Expected:** shutdown is a barrier against admission. **Actual:** it is not.

### CE-6 — Kill without a reaped guard

1. p1 exits; `Wait` returns or is completing; observer is about to close `done`.
2. A concurrent `Terminate` timer/cancellation selects the kill branch before observing `done`.
3. `p.Kill` signals the old PID handle while the process is exiting or being reaped.

**Expected:** no kill after reaped; no PID-based signal. **Actual:** the kill path has no reaped guard.

---

## 7. Contract rev3 assessment

Contract rev3 is directionally useful but is not an implementation-ready contract.

| Contract area | Consistency with source/P17 | Finding |
|---|---|---|
| A1 process ownership split | Source exposes `Wait`/`Kill`/`Close` to the client through `Process`; no split | Not implemented; F-10/F-03 |
| A2 drain-before-reap | Source calls `Wait` in `observeExit` without joining stdout/stderr drains; `Terminate` waits for process, not reader | Not implemented; F-09/F-10 |
| A3 quarantine | Source has no quarantine record; timeout remains `stopping`/owned | Not implemented; F-17 |
| A4 publication authorization | Source has no lifecycle lock/linearization/registration; client writes supervisor state directly | Not implemented; F-02/F-03 |
| A5 transport arbitration | Source has a raw transport pointer; no lease; `Shutdown` can start | Not implemented; F-04/F-11 |
| A6 terminal evidence | Rule 3 and table project `crashed` from fatal failure even at exit 0; conflicts with OD-5 and P17 | **Contradiction; blocking** F-05 |
| A8 invariants I1–I8 | Several are directly violated by source paths (I2 kill-after-reap, I3 Wait-before-drain, I8 expected exit) | Contract is aspirational, not current behavior |

Contract A6's `failureObserved` model also records only fatal observations from reader/handshake/observer; it does not supply the missing generation, action identity, or exit status needed by OD-8. It cannot by itself close U1–U7.

---

## 8. OD-1 through OD-9 consistency

| Decision | Accepted direction | Source / P17 consistency | Review disposition |
|---|---|---|---|
| OD-1 | Projection-supporting evidence; supervisor authority; client observations retained | No evidence log, no supervisor publication boundary; client calls `SetState`/`MarkHandshake`; P17 P1 states the requirement | **Accepted governance, unimplemented; HIGH gap** |
| OD-2 | Session + supervisor generation + PID; generation is process-lifetime key; no old evidence on replacement | No generation field; old reader can affect replacement; P17 proposes generation but exact allocator open | **BLOCKING** |
| OD-3 | One immutable first terminal; one linearization point; later evidence retained; no old-instance mutation | No latch, no terminal set; writes overwrite; P17 terminal set open | **BLOCKING** |
| OD-4 | Confirmed exit without orderly → `unavailable` unless prior failure; orderly → `absent`; active states cannot survive; internal action attribution | Source maps nonzero-while-stopping to `absent`; P17 corrects but P4-A unaccepted | **Source nonconforming; P17 policy incomplete** |
| OD-5 | `crashed` = observed exit without orderly shutdown; live process must not be crashed; `ProcessAlive` independent | Source publishes `crashed` on send failure; Status is provider-derived; P17 O18 corrects | **Source nonconforming; U6 open** |
| OD-6 | Missing mandatory capability → `incompatible`; helper/client same established failure share state | Source maps client-detected missing capability to `unavailable`; P17 P2 says incompatible but O16 keeps paths open | **Contradiction; U8 open** |
| OD-7 | Quarantine internal-only, exact generation, explicit ownership/reaping condition | No record/transition; timeout leaves owned process; P17 mechanics open | **Unimplemented; U2/U6/U7** |
| OD-8 | A–E evidence required for causal attribution; A-O association only | No action identity/outcome, no event order, no exit status; U5 cannot be evaluated | **Unimplemented; U5 open** |
| OD-9 | No P17 approval while U1–U8 open without determined outcomes or explicit unknown rules | U1–U8 all OPEN in P17; no per-case unknown rule | **Gate unmet; FAIL** |

The accepted decisions are internally coherent at the policy level. The failure is implementation and specification completeness, plus one direct contract contradiction (A6).

---

## 9. Findings ranked by severity

| ID | Severity | Finding | P17 already addresses? | Prevents C1? | Smallest resolution |
|---|---|---|---|---|---|
| F-01 | BLOCKING | OD-9 gate unmet; U1–U8 open | No; U matrix open | Yes | Determine outcomes or explicit unknown rules |
| F-02 | BLOCKING | No terminal latch; terminal set undefined | Principle only; no set/latch | Yes | Define set + linearization + latch |
| F-03 | BLOCKING | No generation isolation; old reader can affect replacement | Requirement only; no mechanism | Yes | Generation-bound evidence/operations |
| F-04 | BLOCKING | Shutdown can start replacement after crash | No | Yes | Lifecycle-closing shutdown; no restart |
| F-05 | BLOCKING | Contract A6 contradicts OD-5/P17 | P17 correct; contract wrong | Yes | Reconcile A6 before use |
| F-06 | HIGH | MarkHandshake resurrects dead process | No | Yes, until fixed | Generation/current-process check |
| F-07 | HIGH | Ready/degraded after termination/exit | Partly; no request-timeout rule | Yes, until fixed | Typed cleanup outcome; no post-exit active projection |
| F-08 | HIGH | ProcessAlive not independent | OD-5 requirement; not implemented | Yes, until fixed | OS/handle liveness predicate |
| F-09 | HIGH | Reader never joined; replacement race | E5/E9 only; no join rule | Yes | One reader per generation; join before replace |
| F-10 | HIGH | Kill after reaping/repeated Terminate | E8 only; no state machine | Yes | A1 state machine; kill no-op after reap |
| F-11 | HIGH | Shutdown/Start no common barrier | No | Yes | Shared admission/closure lock |
| F-12 | MEDIUM | Cleanup outcome not retained; retention can be indistinguishable from failure | E9; no outcome record | No by itself; contributes to ownership uncertainty | Actual liveness; typed cleanup outcome |
| F-13 | MEDIUM | U assessment stale on U3/U4/U5/U6 | N/A | Not by itself; undermines evidence | Rebind assessment to rev2 |
| F-14 | MEDIUM | P17 absent rule depends on unaccepted P4-A | P4-A proposed | Yes, for public `absent` | Accept P4-A or state the rule |
| F-15 | MEDIUM | `outcome_unknown` not assigned to open cases | Permitted, not required | Yes, for OD-9 | Explicit unknown rule |
| F-16 | MEDIUM | No action/exit status for OD-8 | Proposed shape only | Yes, for U5 | Action/exit evidence model |
| F-17 | MEDIUM | No quarantine transition | Internal concept only | Yes, for U2/U6/U7 | Explicit timeout/ownership/reaping records |
| F-18 | LOW | `retired` unbounded; status/catalog start helper | No | Not directly | Bound retired; non-mutating status |

---

## 10. Required fixes or decisions in dependency order

These are findings-level requirements, not authorization to edit any artifact.

1. **Satisfy OD-9 first.** For each U1–U8, record either a determined outcome or an explicit internal unknown rule. The rule must state the evidence, recorder, retained record, and why no public state is selected. “OPEN” alone is not enough.
2. **Define the terminal set and linearization point.** State which public values can become the first terminal, who publishes them, how the latch is synchronized, and how later evidence is retained without republication. This is required before U3/U7 and CP-1 can be evaluated.
3. **Define the generation and ownership boundary.** Every reader, action, exit, failure, and publication must carry the OD-2 generation. Retire/join the old reader before replacement. Old-generation evidence may be recorded but must be unable to fail or kill the current generation. This addresses F-03, F-09, F-11.
4. **Make shutdown a lifecycle-closing operation.** Shutdown must not call `Start`, must not start after a crash, and must exclude concurrent starts. A stale transport is not permission to create a new process. This addresses F-04 and F-11.
5. **Define orderly shutdown and exit/action evidence.** Add stream-close provenance, action identity and outcome, exit status/signal, and event order. Apply OD-8's A–E standard. This addresses U1/U2/U4/U5 and F-16.
6. **Fix U6 and U8 source paths.** Never publish `crashed` for a live process; derive liveness from actual process ownership; map missing mandatory capability to the OD-6 `incompatible` category; narrow the P17 “path-by-path OPEN” statement to unclassified categories. This addresses F-06/F-08/F-12/U6/U8.
7. **Reconcile contract rev3 A6 with OD-5/OD-4.** Remove the fatal-observation-to-`crashed` projection. Do not use A6 as an implementation basis until it matches the accepted semantics. Then review A1/A2/A4/A5 as implementation mechanisms.
8. **Resolve the timeout/quarantine relationship.** A termination timeout is not a terminal state and not automatically quarantine; retained ownership and unresolved reaping are distinct. This addresses U2/U7/F-17.
9. **Resolve or explicitly retain CP-1, N-5, and N-6.** CP-1 affects prior `crashed` preservation; N-5 affects multi-classification selection; N-6 affects a non-failure public value. Do not let any of them be silently decided by the implementation.
10. **Re-run independent review against the exact paths, line counts, and SHA-256s.** Only after the OD-9 gate and the public projection rules are satisfied should C1 implementation be reconsidered.

Items 2–5 are not merely code details. They determine which public state can be published and whether an old process can affect a new one. Items such as the lock implementation, generation allocator, and storage schema may be decided during C1 once the semantic rules are fixed; this review does not demand a particular mechanism.

---

## 11. C1 recommendation

**Do not begin C1 implementation under P17 rev2 or contract rev3 as currently written.**

The current source is not merely missing tests. It has concrete paths that can:

- publish `crashed` or `ready` without an observed exit, or after an exit;
- let an old reader fail or kill a replacement;
- start a helper during `Shutdown`;
- overwrite a terminal with `stopping`, `absent`, or `crashed`;
- report process liveness from provider state;
- signal a process without a reaped guard;
- leave cleanup completion unproven when a helper is intentionally retained.

A partial C1 implementation might be technically possible with defensive safeguards, but it would not be an implementation of the governing model, and the OD-9 gate expressly forbids approving P17 while U1–U8 remain open without determined or explicit unknown dispositions. The safe recommendation is **no C1 implementation yet**.

---

## 12. Limitations and unverified assumptions

1. No runtime, test, build, benchmark, or application execution was performed. All interleavings are source-reachable schedules, not observed runs. The actual winner in a real execution is unestablished.
2. The helper tree is untracked. The content digest binds the inspected files, but Git cannot prove they are the only files that changed or that they match an earlier reviewed state.
3. The review did not inspect every helper test or every platform-specific implementation. It inspected the lifecycle-relevant production files and their call sites. `proc.KillTree` platform behavior was not executed or independently verified.
4. The accepted remaining-evidence assessment's historical digest was not authenticated; this review used the owner-accepted working copy only. The primary U1–U8 assessment is at its verified on-disk digest but is stale in the ways described in F-13.
5. P17 and contract rev3 are proposals, not accepted implementation contracts. The review distinguishes their defects from current source defects; it does not treat either document as proof of implementation.
6. The thirteen-value provider-state set is unchanged. `outcome_unknown`, `terminated_by_action`, and `quarantined` remain internal classifications; none is a public state.
7. CP-1, the terminal-set definition, N-5, N-6, S5-scope, and E1–E9 remain OPEN as described in the inspected documents. This review does not close or re-decide them.
8. No implementation files, tests, contracts, proposals, assessments, reviews, or registers were changed. No tests, builds, benchmarks, or runtime experiments were performed.

---

## 13. Final statement

**P17 remains UNAPPROVED. Contract rev3 remains frozen and unaccepted. C1 implementation remains UNAUTHORIZED. The final verdict is FAIL — NOT READY.**

The source establishes real lifecycle hazards, not merely missing wording. The most serious are generation cross-talk, shutdown-initiated replacement, missing terminal immutability, post-exit active projections, and the contract's fatal-observation-to-`crashed` rule. The accepted OD-9 gate independently prevents approval while U1–U8 remain open without determined or explicit unknown dispositions.

This report is an independent review artifact, not a correction pass, not an approval, and not implementation authorization.
