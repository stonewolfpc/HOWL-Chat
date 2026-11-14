# C1 — Finding-to-Mechanism Reconciliation: R3-01…R3-21

**Date:** 2026-11-09
**Status:** AUTHOR RECONCILIATION. **Not closure. Not rev4.**
**Inputs:** rev3 review record (owner-supplied); source at HEAD `9c7d2a1`; contract rev3 `ed0da1a7…`.
**Provenance limit:** the review's per-ID text is **not in the workspace**. IDs are tracked individually; mechanism scope is inherited from the four groups named by the owner. No per-ID claim is quoted or invented.

## 1. Reading rules

1. Each R3 ID is listed separately and stays open.
2. Mechanisms are numbered P1–P18; each cites source and a rev3 section.
3. Mechanism scope is group-level where the per-ID text is unavailable; this is a provenance limit, not a merge.
4. No rev4 is written here. The reviewer examines these mechanisms first.
5. Nothing below closes W1–W4, X-1, X-2, or any R3 finding.

## 2. Finding index

| ID | Group | Mechanisms in scope | Claim text | Status |
|---|---|---|---|---|
| R3-01 | process/pipe | P1–P6 | not supplied | open |
| R3-02 | process/pipe | P1–P6 | not supplied | open |
| R3-03 | process/pipe | P1–P6 | not supplied | open |
| R3-04 | process/pipe | P1–P6 | not supplied | open |
| R3-05 | process/pipe | P1–P6 | not supplied | open |
| R3-06 | process/pipe | P1–P6 | not supplied | open |
| R3-07 | process/pipe | P1–P6 | not supplied | open |
| R3-08 | publication | P7–P9 | not supplied | open |
| R3-09 | publication | P7–P9 | not supplied | open |
| R3-10 | publication | P7–P9 | not supplied | open |
| R3-11 | publication | P7–P9 | not supplied | open |
| R3-12 | transport | P10–P12 | not supplied | open |
| R3-13 | transport | P10–P12 | not supplied | open |
| R3-14 | transport | P10–P12 | not supplied | open |
| R3-15 | transport | P10–P12 | not supplied | open |
| R3-16 | terminal | P13–P18 | not supplied | open |
| R3-17 | terminal | P13–P18 | not supplied | open |
| R3-18 | terminal | P13–P18 | not supplied | open |
| R3-19 | terminal | P13–P18 | not supplied | open |
| R3-20 | terminal | P13–P18 | not supplied | open |
| R3-21 | terminal | P13–P18 | not supplied | open |

## 3. Mechanisms

### P1 — refused `pending` process has no legal reaping edge

**Reported failure:** a refused pending process cannot enter reaping. **Source (F):** `Start` installs under `s.mu` (`supervisor.go:151–163`) and has no refusal branch; the factory result is always installed. **Rev3:** A1 table gives `pending` `Kill=no, Wait=no`, yet W2 requires refusal reaping. **Contradiction:** the only state that can own a refused process forbids the operations needed to reap it. **Correction (P):** add `pending → reaping` on refusal, owned by the installer, with `Kill`/`Wait` legal in `pending-reaping` and forbidden after `reaped`. **Depends on:** P4 (who closes stdin), P5 (single Wait). **Residual:** whether refusal can occur after the observer starts; it cannot today, but the edge must forbid it.

### P2 — forced drain does not join goroutines before `Wait`

**Reported failure:** the forced-drain path does not join the drains before `Wait`. **Source (F):** `drainStderr` is a bare goroutine with no completion channel (`supervisor.go:527–539`); the stdout reader likewise has `readerDone` only in the client, not in the process abstraction. **Rev3:** A2 says close read ends, then `Wait`; I3 says no `Wait` before both drains complete. **Contradiction:** I3 is unverifiable because the abstraction exposes no drain-done signal. **Correction (P):** require `DrainResult{kind, err}` from both drains and join on those before `Wait`; forced close records `forced`. **Depends on:** P3. **Residual:** a drain blocked in `readFrame` must observe close as an error and report `forced`, not hang.

### P3 — natural EOF and forced closure are indistinguishable

**Reported failure:** the two outcomes are not distinguishable. **Source (F):** `readFrame` returns only `([]byte, error)`; `drainStderr` stops on any error (`supervisor.go:535–538`). **Rev3:** A2 lists "EOF or read error" as one signal. **Contradiction:** the evidence model needs `failureObserved` provenance, but a forced close to unblock a hang is indistinguishable from a natural exit. **Correction (P):** drain result carries `natural | forced | error`; `forced` is attributed to the escalation step, never to the child. **Depends on:** P2. **Residual:** distinguishing an EOF from a read error after the process died is ambiguous on some platforms; must be recorded as `unknown` rather than guessed.

### P4 — stdin has no consistent owner

**Reported failure:** stdin closure has no consistent owner. **Source (F):** `Send` writes stdin (`supervisor.go:458–476`); `Close` closes stdin and stderr (`:480–491`); the client holds `Process` and can call `Close` today. **Rev3:** A1 gives the client `Transport` (which includes `Send`), A2 assigns stdin to a "transport owner" that is never defined, and A5's arbiter omits stdin. **Contradiction:** three components can close stdin; none is named. **Correction (P):** stdin belongs to the transport lease holder; `Lifecycle.Close` never closes stdin while the lease is held; unblocking a child is an explicit `closeStdin` transport operation. **Depends on:** P10–P12 (arbiter policy). **Residual:** whether closing stdin is required to make a hung child exit on Windows.

### P5 — "Wait again" conflicts with one-owner

**Reported failure:** A3's second `Wait` conflicts with I1. **Source (F):** `Wait` is wrapped in `waitOnce` (`supervisor.go:502–510`); a second call returns the cached error and does not re-reap. **Rev3:** A3 E2 says "Wait again"; I1 says one owner for life. **Contradiction:** the text implies a second reaping attempt that the abstraction cannot perform. **Correction (P):** remove "Wait again". The single owner re-arms its deadline on the same `Wait`; every other participant waits on `waitDone`/`Exit` and never calls `Wait`. Timeout never transfers ownership. **Depends on:** P1. **Residual:** none technical; the wording is the defect.

### P6 — quarantine has no terminal projection

**Reported failure:** quarantine has no explicit projection. **Source (F):** the state set is 13 values and contains no `quarantined` (`provider.go:22–36`). **Rev3:** A3 promises a "typed diagnostic"; A6 has no quarantine row. **Contradiction:** A3's promise cannot be published through a taxonomy that lacks the value. **Correction (P):** add `quarantined` to the projection set with its own publish path (P8), distinct from `absent` and `crashed`. **Depends on:** P8, P18. **Residual:** whether `quarantined` is a provider state or a process-only state; today's type conflates them.

### P7 — pre-retirement and post-retirement delivery conflict

**Reported failure:** the rule both permits and prohibits delayed delivery. **Source (rev3 §A4, F):** rows "authorized before retirement may complete; tagged `preRetirement`" and "delivery after retirement refused". **Contradiction:** any delayed completion is delivered after retirement, so both rows apply. **Correction (P):** the linearization point is authorization; delivery is legal iff authorized, and the tag records when authorization occurred. "Post-retirement delivery" means delivery **without** prior authorization. **Depends on:** P8. **Residual:** how long an authorized delivery may remain outstanding before it is abandoned.

### P8 — terminal publication has no authorization path

**Reported failure:** a terminal projection needs explicit authorization and ordinary publications must not overwrite it. **Source (rev3 §A4/§A6, F):** A6 is published through the same A4 protocol as ordinary events; no latch or precedence exists. **Contradiction:** a late ordinary publication can overwrite a terminal projection. **Correction (P):** `PublishTerminal` is a distinct authorization, exactly once per instance; after it commits, ordinary publications are refused and recorded. **Depends on:** P18. **Residual:** whether a terminal projection may be revised on new fatal evidence (recommend: no, append evidence only).

### P9 — attachment capability not restated

**Reported failure:** W1/W3 cannot close because the capability is missing from rev3. **Source (rev3, F):** A1 defines `Transport`/`Lifecycle` but no `StartResult` or `Attach`; rev2 §1 held that content. **Contradiction:** rev3's interfaces have no lifecycle authority without the attachment transaction. **Correction (P):** rev4 must restate the capability (single-use, token+epoch bound, invalidated at retirement) alongside the interface split. **Depends on:** P7, P8. **Residual:** none; omission only.

### P10 — acknowledgment deadlock

**Reported failure:** sender holds the lease while awaiting an ack that needs the lease. **Source (F):** `Shutdown` calls `call` and waits for the response (`client.go:651`), then `Terminate` (`:664`, `:674`). **Rev3:** A5 requires the lease for both the shutdown frame and acknowledgment processing. **Contradiction:** the waiter holds the resource the ack needs. **Correction (P):** separate write-lease from read/inbox arbitration. The sender releases the write-lease after the frame and awaits the ack on the inbox channel with a deadline; ack processing never takes the write-lease. **Depends on:** P11. **Residual:** ordering when an ack and a close race.

### P11 — lease policy absent

**Reported failure:** eight aspects need separate treatment. **Source (rev3 §A5, F):** one table, one state machine, no precedence or expiry rule. **Contradiction:** "everyone takes the lease" is not a policy; two writers can starve each other and closure can block forever. **Correction (P):** define precedence (ack > shutdown frame > ordinary > close), bounded queue, expiry winner, and late-completion recording. **Depends on:** P10. **Residual:** fairness policy under sustained ordinary traffic.

### P12 — evidence after transport closure

**Reported failure:** evidence recording after closure is unspecified. **Source (rev3 §A5/§A6, F):** evidence is published through A4, which is refused after retirement. **Contradiction:** the final exit evidence arrives as the transport closes. **Correction (P):** terminal evidence bypasses the transport and is written to the post-closure evidence buffer under the supervisor, then published through P8. **Depends on:** P8, P10. **Residual:** durable retention if the supervisor itself is shutting down.

### P13 — exit zero is treated as expected

**Reported failure:** zero alone does not establish an expected exit. **Source (rev3 §A6 row 1, F):** `0 + none + known → absent` without requiring an action. **Contradiction:** I8 says no expected exit without evidence, yet the row supplies none. **Correction (P):** `absent` requires an action that explains the exit; zero with no action projects `unknown`. **Depends on:** P17. **Residual:** whether a clean exit after a failed handshake is `incompatible` (P15) rather than `unknown`.

### P14 — ambiguous nonzero is treated as a crash

**Reported failure:** ambiguous causality plus nonzero does not establish a crash. **Source (rev3 §A6 row 5, F):** `ambiguous + nonzero → crashed`. **Contradiction:** this manufactures the certainty the model was built to avoid. **Correction (P):** `ambiguous` projects `unknown` with the exit status retained; `crashed` requires known causality or a fatal observation. **Depends on:** P17. **Residual:** whether a nonzero exit during a deliberate kill can ever be `unknown` vs `killed`.

### P15 — handshake incompatibility classified as crash

**Reported failure:** `protocol_incompatible` is not a crash. **Source (rev3 §A6 row 2, F):** any fatal `failureObserved` projects `crashed`; the handshake path publishes `incompatible` separately today (`client.go:248–253`). **Contradiction:** A6 collapses a handshake-class outcome into a process crash. **Correction (P):** `failureObserved` carries class `handshake | protocol | transport | exit`; handshake-class fatal projects `incompatible`. **Depends on:** P17, P18. **Residual:** post-handshake incompatibility (protocol downgrade) needs its own class.

### P16 — deliberate kill treated as unexpected crash

**Reported failure:** a deliberate kill is not automatically an unexpected crash. **Source (rev3 §A6 last row, F):** "action after failure → `crashed` + kind; action recorded as context". **Contradiction:** the kill is the cause of exit but is classified as the crash. **Correction (P):** the termination action explains the exit; `crashed` derives only from evidence that predates the action. A kill after a fatal failure projects the failure's class. **Depends on:** P17. **Residual:** whether a kill with no prior failure projects `absent` or `unknown`; recommend `unknown` until the exit is observed.

### P17 — no severity, class, or precedence

**Reported failure:** fatal vs diagnostic, precedence, and action effects are undefined. **Source (rev3 §A6, F):** `failureObserved[]` has only `{seq, source, kind}`. **Contradiction:** two failures have no defined order, and a diagnostic can be read as fatal. **Correction (P):** add `severity {diagnostic, fatal}` and `class`; precedence: fatal handshake > fatal protocol > exit status; diagnostics never change the projection. **Depends on:** P18. **Residual:** whether transport failures are diagnostic or fatal for the provider state.

### P18 — required projections do not exist

**Reported failure:** `unknown` and `quarantined` are required. **Source (F):** the 13 provider states are absent, starting, handshaking, unloaded, loading, ready, busy, degraded, stopping, crashed, incompatible, unavailable, backoff (`provider.go:22–36`). **Rev3:** A6 projects only `absent`/`crashed`. **Contradiction:** four of the five required projections exist; `unknown` and `quarantined` do not. **Correction (P):** map the taxonomy explicitly: `absent`, `crashed`, `incompatible`, `unavailable` exist; add `unknown` and `quarantined` (or map them to a process-level enum if the provider state must stay frozen). **Depends on:** owner decision, because the state set is part of the frozen helper contract. **Residual:** whether adding states is a contract amendment; **not** an author decision.

## 4. Dependency order

| Order | Must settle first | Because |
|---|---|---|
| 1 | P18 taxonomy decision | P6, P13–P16 cannot publish into a set that lacks their values |
| 2 | P17 severity/class | precedence determines P13–P16 projections |
| 3 | P4 stdin owner + P5 single Wait + P2 drain join | together define P1's reaping edge and I1–I3 |
| 4 | P7/P8 publication rule | P6 and P12 need a terminal path |
| 5 | P10/P11 arbitration policy | P4's lease reference and P12's buffer depend on it |
| 6 | P9 capability restatement | gates W1/W3, not the mechanisms above |

## 5. Status

Twenty-one IDs open; eighteen mechanisms identified; one owner decision required (P18 taxonomy). No contract edit, no rev4, no matrix, no implementation, no tests, nothing built or run. Attestation unsigned. Frozen evidence unchanged. **C1 remains unauthorized.**
