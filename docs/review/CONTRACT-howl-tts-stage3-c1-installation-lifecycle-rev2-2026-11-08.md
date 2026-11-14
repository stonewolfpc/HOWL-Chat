# C1 — Concurrency Contract rev2: Enforceable Installation Ownership

**Date:** 2026-11-08
**Status:** PROPOSED REVISION FOR INDEPENDENT REVIEW. Not an accepted amendment.
**Supersedes for review only:** contract rev1 `7e0b61b1…` (preserved, unedited).
**Inputs:** review record (W1–W4, X-1, X-2, terminal projection); round-2 reconciliation.
**C1 REMAINS UNAUTHORIZED.**

Rev2 changes ownership statements into enforceable transitions. Terminology is not added unless it maps to a check, a sequence number, or an owned resource.

## 1. StartResult and attachment capability

```text
type StartResult struct {
    Proc   Process   // handle for I/O only; grants no attachment right
    Attach attachCap // opaque, supervisor-minted, single-use
}
```

| Property | Rule |
|---|---|
| minted | only by `Supervisor` at commit |
| bound | instance token + epoch + registration nonce |
| use | presented back to `Supervisor.Attach`; never accepted by another instance |
| reuse | single use; a consumed capability is invalid |
| invalidation | retirement invalidates unconsumed capabilities for that token |
| exposure | not a bearer token; cannot be persisted, copied, or used to bypass admission |

`Start` returns `StartResult, error`. The client's early return for an already-attached live session remains, but it proves liveness only and grants no attachment.

## 2. W1 — admission at the commit point

Under `s.mu`, admission requires `state == active && !shutdownRequested`. The existing-process path returns a result **only** when that holds; otherwise it returns typed `retiring`/`shutting_down` and never the handle (`supervisor.go:121–125` today returns any non-nil process).

`Client.Start` keeps its `closing` check as a mirror, not as the authority (`client.go:213–216`); the authority is inside the install critical section, because the two checks are in different critical sections and different objects.

## 3. W2 — bounded reaping with retained ownership

Admission and install happen under `s.mu`; `startMu` is released **before** reaping. A refused process is never installed, but it is never orphaned either.

```text
Pending{token, proc, createdAt, owner: installer}
  → admit?  yes → install, start observer
            no  → Reap(Pending) outside startMu
```

`Reap` is single-flight per token, bounded by `reapTimeout`, and returns exactly one outcome:

| Outcome | Condition | Record | Ownership after |
|---|---|---|---|
| O1 | `Kill`/exit, `Wait` returned, handles closed | retired | released |
| O2 | exited, but `Wait` or cleanup errored | retired + failure | released |
| O3 | exit not confirmed within `reapTimeout` | quarantined | **retained** |
| O4 | pipe/handle close failed | quarantined | **retained** |

Quarantine: `quarantine[token] = {proc, lastErr, owner}` under `s.mu`. It never blocks a new admission, is never silently dropped, and retains the single `Wait` owner. A later `TerminateIf` for the quarantined token may close it, but cannot install over it.

**Process contract:**

1. Exactly one `Wait` call owns reaping; it is not called twice.
2. `Kill` is optional, may precede `Wait`, and never replaces it.
3. `Close` is idempotent and is not a substitute for `Wait`.
4. `Wait` is the join point for pipe EOF; the `stderr` drain is joined **after** `Wait` returns, then handles are closed.
5. `Close` before `Wait` is permitted only for `stdin` shutdown, never for `stderr` while its drain runs.

The platform check required by W2 (Windows `Wait`/pipe semantics) must confirm rule 4 before acceptance.

## 4. W3 — atomic registration, validated execution

Attachment is a supervisor transaction, not a client assignment:

```text
Attach(cap) under s.mu:
  validate cap.token == current.token
  validate cap.epoch == current.epoch
  validate state == active && !shutdownRequested
  validate cap not consumed
  consume cap
  R := registration{token, nonce, state: registered}
  return R, nil
```

Reader lifecycle is explicit: `registered → running → retired`.

| Event | Rule |
|---|---|
| first delivery | reader must re-validate `R` under `s.mu` and move to `running` |
| every publish | reader checks `R.state == running` and token match |
| retirement | supervisor sets `R.state = retired`; publish is refused |
| kill | `KillIf(R.token, intent)` targets the registration's original instance |
| rotation | old reader must be joined before a new attach; timeout is a typed error |

`waitForReader` returns `joined | exited | timeout`; the current silent timeout (`client.go:249`) does not satisfy this.

## 5. W4 — one transaction, one lease

`BeginShutdown()` performs atomically under `s.mu`:

1. latch `primaryIntent = shutdown` (first intent wins);
2. set `shutdownRequested`;
3. enter `retiring`;
4. mint the one-shot permit;
5. acquire the transport lease.

| Artefact | Scope | Lifetime | On failure |
|---|---|---|---|
| permit | exactly one `shutdown` frame | consumed at send completion | records `sendFailed`, remains consumed |
| lease | transport writes during that send | released after ack/timeout/send error | observer may close only after release or expiry |

The observer must not `Close` the transport while the lease is held. At lease expiry the observer wins and records `transportClosedDuringShutdown`. The permit never reopens admission, never authorises `Start`, and cannot publish state.

Order (the reverse of today's `Shutdown`):

```text
BeginShutdown → send(shutdown, permit) → consume permit → release lease
→ TerminateIf(token, shutdown) → retire → observer projection
```

## 6. Terminal projection — ordering and evidence

Three fields, distinct authority:

| Field | Written by | Meaning |
|---|---|---|
| `primaryIntent` | transition owner | authorises the transition |
| `crashObserved{seq, status}` | observer only | a fact about the process |
| `terminalProjection` | observer | what is published |

Rules:

1. `crashObserved` is set once, from `Wait`, with a monotonic `seq` assigned under `s.mu`.
2. `terminationAction` carries its own `seq`.
3. `crashObserved.seq < terminationAction.seq` → the crash is the fact; the action is context.
4. `crashObserved.seq > terminationAction.seq` → the action requested the stop; a zero exit is expected, a nonzero exit is `crashed`.
5. `providerState == stopping` is not evidence and must not map a nonzero exit to `absent` (today `supervisor.go:176–187`).

Handshake classification:

| Condition | Phase | Projection |
|---|---|---|
| incompatible / minor unsupported | pre | `incompatible` |
| protocol violation before hello | pre | `incompatible` + protocol detail |
| protocol violation after hello | post | `crashed` + `protocol_violation` |
| exit 0, shutdown | post | `absent` |
| nonzero, shutdown, no prior crash | post | `crashed` |

## 7. Invariants to review

| # | Invariant |
|---|---|
| I1 | every created process has exactly one `Wait` owner at all times |
| I2 | no attachment without a consumed capability matching token and epoch |
| I3 | no publish after retirement, from any reader or exit observer |
| I4 | no transport close while the shutdown lease is held |
| I5 | no admission after `shutdownRequested`, including after quarantine |
| I6 | no `stopping`-based expectation of a nonzero exit |

## 8. Open items

- Windows `Wait`/pipe ordering (W2 rule 4) requires platform confirmation.
- Crash-before-request delivered-after ordering (sequence-only sufficiency) needs a reviewer interleaving test.
- Lease-expiry winner needs confirmation.

## 9. Status

No r9 edit, no matrix regeneration, no implementation, no tests, nothing built or run. Frozen evidence unchanged. **C1 remains unauthorized.**
