# C1 — Concurrency Contract rev3: Enforceable Boundaries

**Date:** 2026-11-09
**Status:** PROPOSED REVISION FOR INDEPENDENT REVIEW. Not an accepted amendment.
**Supersedes for review only:** rev2 `501e58a4…`, rev1 `7e0b61b1…` (both preserved, unedited).
**Inputs:** rev2 review record; owner's six revision areas.
**C1 REMAINS UNAUTHORIZED.**

Rev3 withdraws four rev2 mechanisms that could be bypassed, and replaces them with boundaries a caller cannot cross by convention.

## 0. Withdrawals from rev2

| Rev2 rule | Why withdrawn | Replaced by |
|---|---|---|
| §3 rule 4: `Wait` is the EOF join point | contradicts `os/exec`: reads must finish before `Wait`, which closes the pipes | A2 drain-before-reap order |
| `StartResult.Proc` documented "I/O only" | a comment is not enforcement; the handle still exposes `Wait`/`Kill`/`Close` | A1 interface split |
| §5 transport lease held by the supervisor | client can write and close the transport independently | A5 arbiter for every writer |
| §6 ordering by sequence alone | sequence orders observations, not events | A6 causality record |

## A1 — Process ownership: split interfaces

```text
Transport  interface { Send([]byte) error; Receive() ([]byte, error); Stderr() string }
Lifecycle  interface { Kill(ctx) error; Wait(ctx) (Exit, error); Close() error }  // supervisor only
```

The client receives `Transport`; only the supervisor holds `Lifecycle`. The client cannot call `Wait`, `Kill`, or `Close`.

**Kill-after-reap prevention.** Instance state machine: `pending → live → reaping → reaped`, with `→ quarantined` from `reaping`.

| State | `Kill` allowed | `Wait` allowed | Notes |
|---|---|---|---|
| pending | no | no | not yet installed |
| live | yes, once | no | normal operation |
| reaping | no (second kill refused) | yes, single owner | kill issued at entry |
| reaped | **no** | no | handle release only; PID kills are forbidden |
| quarantined | escalation only | retained owner | see A3 |

The supervisor checks the token and state under its lock before forwarding any kill. A stale PID is never signalled: on Windows a recycled PID would target an unrelated process.

## A2 — Pipe lifecycle: four responsibilities, one order

| Responsibility | Owner | Completion signal |
|---|---|---|
| stdin closure | transport owner | idempotent; not a reap signal |
| stdout read | reader goroutine | EOF or read error |
| stderr drain | drain goroutine | EOF or read error |
| reap | one `Wait` owner | `Wait` returned |
| final release | reaper | handles closed, result recorded |

Order, corrected against `os/exec`:

```text
1. Kill (optional, once, state live)
2. stdout reader and stderr drain reach EOF and report done
3. Wait — exactly once, only after both drains are done
4. Close remaining handles (idempotent)
5. publish projection (A6)
```

`Wait` is **not** used to force EOF. Because `Wait` closes the pipes, calling it before the drains finish is incorrect; the drains must be joined first.

**Hang handling.** If a drain does not finish within `drainTimeout`, close the read ends to force EOF (no further reads are expected), then `Wait` within `reapTimeout`. Failure of either step moves the instance to quarantine (A3). `stdin` may be closed early to unblock the child; `stderr` is never closed while its drain runs.

## A3 — Quarantine: bounded ownership, no admission lock

```text
quarantine[token] = { proc, lifecycle, waitOwner: true, lastErr, attempts }
```

Escalation ladder, each step bounded:

| Step | Action | On success |
|---|---|---|
| E1 | kill, drains, `Wait` | reaped |
| E2 | close read ends, `Wait` again | reaped |
| E3 | second kill is **not** attempted; record stuck state | quarantined |
| E4 | retain `Wait` ownership; surface typed diagnostic | quarantined |

Quarantine never blocks admission, holds no lock during any blocking call, is never silently dropped, and does not permit installing over its token. A quarantined process is reported distinctly from `absent` and from `crashed`.

## A4 — Publication: authorize, release, deliver

Every publication — handshake result, request delivery, failure delivery, supervisor state, exit projection — uses one protocol:

```text
1. authorize under the lifecycle lock (linearization point)
2. register in-flight publication {token, seq, state: authorized}
3. release the lock
4. deliver outside the lock
5. complete or cancel; record outcome
```

| Event | Rule |
|---|---|
| retirement | closes authorization for new publications |
| publication authorized before retirement | may complete; tagged `preRetirement` |
| registered but never running by `abandonDeadline` | `abandoned`; may not publish |
| reader rotation | old registration retired and joined before new attach |
| delivery after retirement | refused; recorded as `postRetirement` violation |

No callback executes while the lifecycle lock is held, so blocking delivery cannot stall retirement.

## A5 — Transport arbitration: every writer

One arbiter, states `idle → leased → closing → closed`.

| Writer | Must hold lease | On refusal |
|---|---|---|
| ordinary request | yes | typed `transportBusy` |
| shutdown frame | yes (permit + lease) | `sendFailed`, permit consumed |
| transport close | yes | defer until release |
| acknowledgment processing | yes | dropped, recorded |
| expiry / cancellation | yes | observer wins, records outcome |
| late completion | no (already closed) | recorded `lateAckDiscarded` |

No-transport shutdown is a first-class typed outcome (`shutdownNoTransport`); it never falls back to `Start`. The permit and the lease are distinct: the permit authorises one frame, the lease arbitrates every access.

## A6 — Terminal evidence: events, observations, ambiguity

```text
evidence {
  exitStatus        // from Wait
  observedAt        // when the observer recorded it
  actionKind        // shutdown | cleanup | protocol | none
  actionStartedAt
  actionCompletedAt
  failureObserved[] // {seq, source, kind} — reader, handshake, observer
  causality         // known | ambiguous
}
```

Rules:

1. No inference of liveness from action timing. `actionStartedAt` before `observedAt` does not prove the process was alive at the action.
2. `failureObserved` is written by the reader on a fatal protocol violation and by handshake on a fatal handshake failure; a later clean exit does not erase it.
3. A fatal `failureObserved` projects `crashed` with that kind even when `exitStatus == 0`.
4. When causality cannot be established, publish conservatively and retain `causality: ambiguous`; never publish `absent` as "expected" without evidence.
5. `providerState == stopping` is not evidence and is not consulted.

| Exit | failureObserved | Causality | Projection |
|---|---|---|---|
| 0 | none | known | `absent` |
| 0 | fatal | any | `crashed` + kind |
| nonzero | none | known | `crashed` |
| nonzero | fatal | any | `crashed` + kind |
| any | none | ambiguous | `crashed` if nonzero, else `absent` + ambiguous |
| any | fatal, action after failure | any | `crashed` + kind; action recorded as context |

## A7 — Finding map

| Finding | Rev3 treatment | Status |
|---|---|---|
| W1 | A1 capability plus admission at commit | open |
| W2 | A1 split, A2 order (rev2 rule withdrawn) | open |
| W3 | A4 authorize-release-deliver, abandonment, rotation | open |
| W4 | A5 arbiter for every writer, no-transport outcome | open |
| Terminal projection | A6 evidence and ambiguity | open |
| X-1 | admission at commit | open |
| X-2 | A4 rotation join | open |
| R2-01…R2-05 | **reserved**, text not in workspace | open |
| Rev2 open item 1 (pipe/`Wait`) | resolved in A2 by withdrawal | review required |
| Rev2 open item 2 (crash-before-request) | A6 ambiguity, not a sequence fix | open |
| Rev2 open item 3 (lease expiry winner) | A5 observer wins, late ack recorded | review required |

No finding is merged or closed.

## A8 — Invariants

| # | Invariant |
|---|---|
| I1 | one `Wait` owner per process for its whole life |
| I2 | no `Kill` after `reaped`; no PID-based signalling |
| I3 | no `Wait` before both drains complete |
| I4 | no publication without authorization; none after retirement |
| I5 | no transport access without the lease |
| I6 | no admission after `shutdownRequested` |
| I7 | no quarantine entry without retained `Wait` ownership |
| I8 | no "expected exit" classification without evidence |

## A9 — Status

No r9 edit, no matrix regeneration, no implementation, no tests, nothing built or run. Attestation unsigned. Frozen evidence unchanged. **C1 remains unauthorized.**
