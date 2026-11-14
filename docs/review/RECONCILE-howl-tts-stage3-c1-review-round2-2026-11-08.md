# C1 — Reconciliation, Round 2: W1–W4, X-1, X-2, Terminal Projection

**Date:** 2026-11-08
**Status:** AUTHOR RECONCILIATION. **Not closure, not acceptance.**
**Inputs:** review record (owner-supplied, provenance stated there); source tree at HEAD `9c7d2a1`.
**Frozen, untouched:** attestation `598d40c5…`, ledger `e854e9f8…`, contract rev1 `7e0b61b1…`, packet `a1c26a1b…`, analysis `bdb68695…`, addendum `62e4005e…`.

Reviewer disposition is separated from author response in every row. No row is closed.

## 1. W1 — admission rule accepted, API incomplete

**Reviewer:** the rule is right; `Process` cannot carry identity; client has its own check.

**Source (**F**):** `Supervisor.Start` returns `s.proc` whenever non-nil (`:121–125`), with no state or admission check. `Client.Start` short-circuits on `started && transport != nil && sup.ProcessAlive()` (`:197–201`) and never consults the supervisor's admission decision.

**Author response (P):** concede both. `Start` must return `StartResult{proc, attach}`, where `attach` is a supervisor-minted, single-use capability bound to token and epoch, invalidated at retirement. It is not a bearer token handed to callers for reuse; the client presents it back to the supervisor to attach.

**Residual:** the client's early return must remain valid for an already-attached live session. It proves liveness, not admission; the contract must state that it grants no attachment right.

## 2. W2 — bounded reaping and retained ownership

**Reviewer:** four outcomes; quarantine; no indefinite `startMu` hold; ordering contract.

**Source (**F**):** factory creation and install are inside `s.startMu` for the whole `Start` (`:113–163`); no refusal path exists, so a created process is always installed. `observeExit` calls `Wait` → `Stderr` → `Close` (`:166–168`) and `osProcess.Close` closes `stdin` and `stderr` only (`:480–491`); `Kill` kills the tree and the process handle (`:512–519`); `Wait` is guarded by `waitOnce` (`:502–510`).

**Author response (P):** concede. Admission decision is made under `s.mu`; `startMu` is released before any reaping. A per-token reaper runs single-flight outside both locks, bounded by `reapTimeout`, and records one of four outcomes:

| Outcome | Meaning | Handling |
|---|---|---|
| O1 | killed/exited, `Wait` returned, handles closed | retired; token released |
| O2 | exited, but `Wait` error or cleanup error | retired with failure retained; not `absent` |
| O3 | termination not confirmed within `reapTimeout` | **quarantine**, ownership retained |
| O4 | handle/pipe close failed | **quarantine**, ownership retained |

Quarantine is a set under `s.mu` keyed by token, holding the process and its `Wait` ownership; it is never dropped silently, and it never blocks a new admission.

**Process contract (P):** exactly one `Wait` owns reaping; `Kill` is optional and precedes it; `Close` is idempotent but is **not** a substitute for `Wait`; `stderr` is not closed while its drain goroutine still reads it. The current order (`Wait` then `Close`) is compatible with that, but the drain's EOF dependency on `Wait` must be stated explicitly.

**Residual (**F**/**P**):** Go's `Wait` closes the pipes it created, so the correct join point for `drainStderr` is **after** `Wait` returns, not before. This needs a platform-specific check on Windows before the contract is fixed.

## 3. W3 — registration versus execution

**Reviewer:** `s.mu` check plus `c.mu` attach is not atomic; registration and execution are distinct; identity before publish; `Kill` targets the original.

**Source (**F**):** `Client.Start` installs the transport and starts the reader under `c.mu` (`:222–231`), then handshakes (`:233`). The supervisor never learns that a reader exists. `readLoop(p)` closes over `p`, so the reader's identity is correct today, but nothing invalidates it at retirement.

**Author response (P):** concede. `Attach(cap)` is a supervisor operation: under `s.mu` it validates token, epoch, and state, and returns a registration `R{token, state: registered}`. The reader goroutine must re-validate `R` under `s.mu` before its first delivery (transition to `running`), and check it again before every publish. Retirement sets `R.state = retired`; a retired registration may not publish. `KillIf(R.token, intent)` binds to the registration's original instance.

**Residual (**F**):** `waitForReader` currently returns nothing and is silent on timeout (`:249`). It must return `joined | exited | timeout`, and timeout must be a typed failure, not a pass.

## 4. W4 — atomic admission plus transport lease

**Reviewer:** `BeginShutdown` must be one transaction; transport ownership during the permitted send; unambiguous consumption on send failure.

**Source (**F**):** the current order is inverted. `Shutdown` sends the frame first (`:651`), then terminates (`:664`/`:674`), and only then sets `closing` (`:669–673`). `call` routes through `Start`, so the shutdown frame can itself spawn a process.

**Author response (P):** concede. `BeginShutdown()` performs, in one `s.mu` transaction: latch `primaryIntent=shutdown`, set `shutdownRequested`, enter `retiring`, mint the one-shot permit, and acquire the transport lease. The permit is the only writer allowed during the send; the observer may not `Close` the transport until the lease is released, or the lease expires, whichever comes first. Order:

```text
BeginShutdown() → permit+lease → send shutdown (permit, no Start)
→ ack | timeout | send failure → consume permit → release lease
→ TerminateIf(token, shutdown) → retire
```

A failed send consumes the permit and records `sendFailed`; it does not return the permit and does not reopen admission.

**Residual:** lease expiry while the observer waits must have one defined winner; the contract proposes the observer wins at expiry and records `transportClosedDuringShutdown`.

## 5. Terminal projection — open

**Reviewer:** operation and outcome are not interchangeable; ordering must be defined; handshake mapping exact; `stopping` proves nothing.

**Source (**F**):** `observeExit` classifies on `providerState == stopping` alone (`:176–187`): under `stopping`, any nonzero exit becomes `Exited` + `Absent`, with the error retained only in `waitErr`. So an unexpected crash during an orderly stop is reported as a normal disappearance.

**Author response (P):** partially concede the mechanism; the design stays open. Proposed rules:

1. `crashObserved` is set only by the observer, only from `Wait` returning a nonzero/abnormal result or a decode failure, and only once, ordered by a monotonic sequence assigned under `s.mu`.
2. `terminationAction` is recorded with its own sequence at `BeginShutdown`/`TerminateIf`.
3. If `crashObserved.seq < terminationAction.seq`, the crash is the fact; the action is context.
4. If `crashObserved.seq > terminationAction.seq`, the action requested the stop and the exit is **not** automatically a crash: the exit status decides (zero → expected; nonzero → `crashed` unless a prior crash was observed).
5. `stopping` is not evidence of expectation and must not suppress a nonzero exit.

Handshake table (pre-handshake vs post-handshake):

| Condition | Handshake phase? | Projection |
|---|---|---|
| incompatible/minor unsupported hello | pre | `incompatible` (handshake-class) |
| protocol violation before hello completes | pre | `incompatible` + protocol detail |
| protocol violation after successful hello | post | `crashed` + `protocol_violation` |
| exit 0 during shutdown | post | `absent` |
| nonzero during shutdown, no prior crash | post | `crashed` + exit |

**Residual:** whether a crash observed *before* a termination request but delivered *after* it can be ordered correctly by sequence alone; the reviewer should test that interleaving.

## 6. X-1, X-2

| ID | Source evidence | Author response |
|---|---|---|
| X-1 | `supervisor.go:121–125` returns any non-nil proc | concede; W1 admission fixes it |
| X-2 | `client.go:206–212` clears transport without joining reader | concede; rotation must join, and timeout is an error |

## 7. Ledger

| Row | Reviewer disposition | Author response | Closed |
|---|---|---|---|
| W1 | accept rule; API incomplete | concede, `StartResult` proposed | no |
| W2 | revise | concede, four outcomes + quarantine | no |
| W3 | revise | concede, attach transaction | no |
| W4 | revise | concede, transaction + lease | no |
| Terminal projection | open | partial mechanism concession | no |
| X-1 | confirmed open | concede | no |
| X-2 | confirmed open | concede | no |

Attestation remains unsigned and unpopulated. No amendment is accepted. **C1 remains unauthorized.**
