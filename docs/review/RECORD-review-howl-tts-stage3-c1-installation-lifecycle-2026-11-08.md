# C1 — Review Record: Installation Lifecycle (W1–W4, X-1, X-2)

**Date recorded:** 2026-11-08
**Status:** RECORD OF A REVIEW SUPPLIED BY THE OWNER. **Not a signed attestation.**
**Provenance:** findings and dispositions received in the owner's message of 2026-11-08. Reviewer identity, date, and method were not stated. This file records them; it does not certify them.
**Attestation:** `ATTESTATION-…-independent-review-2026-11-08.md` remains unsigned and was **not** populated from these findings.

## 1. Verdict

**NOT READY.** The ownership model is accepted in direction; the synchronization and failure semantics that make it enforceable are missing. Not a reason to abandon the architecture; a reason not to promote it.

## 2. Findings as received

### W1 — admission rule accepted, API incomplete

- Requiring an active process and open admission is correct.
- Returning only `Process` cannot convey the identity capability W3 needs.
- The client has its own admission check and early-return path.

**Disposition:** accept rule; API incomplete. A `StartResult` or equivalent opaque attachment capability is required, and it must not be a freely reusable token — lifecycle and consumption stay under supervisor control.

### W2 — ownership problem, not merely error handling

- Synchronous reaping is incomplete without a defined failure path.
- The reaper must not hold `startMu` indefinitely; a non-terminating process must not block future lifecycle operations.

Required outcomes: (1) terminated and reaped; (2) exited, but exit status or cleanup errored; (3) termination not confirmed; (4) handle or resources not closable.

Quarantine is appropriate for (3) and possibly (4). **Ownership survives the failed attempt.** The `Kill`/`Wait`/`Close` ordering needs a defined process-abstraction contract; a deferred `Close` is not automatically safe while the process or its I/O goroutines still depend on those resources.

**Disposition:** revise.

### W3 — the most important concurrency finding

- Checking under `s.mu` and attaching under `c.mu` is not atomic; they are separate critical sections.
- Supervisor-controlled attachment is accepted, with one addition: **reader registration and reader execution are separate lifecycle events.**
- Registration must be validated against the current instance; retirement must invalidate it; the reader must validate identity before delivering messages or publishing state; `Kill` targets the reader's original process instance, never the current one.
- `waitForReader` silently returning after a timeout does not establish exit.

**Disposition:** revise.

### W4 — atomicity must include the transport

`BeginShutdown` must atomically latch intent, close ordinary admission, enter the retiring state, and issue the one-shot permit. After that transaction the shutdown frame still travels through a transport retirement may otherwise close: the contract must say who owns the transport during the permitted send and what stops the observer closing it early. The permit must not reopen admission, and consumption must be unambiguous if the send fails.

**Disposition:** revise.

### Terminal projection — open

- The requesting operation and the observed process outcome are different facts and are not interchangeable.
- `crashObserved` and `terminationAction` are useful, but a boolean does not resolve the race: the contract must define when a crash becomes an observed fact, how that is ordered against termination, and what happens when both race.
- Handshake mapping needs an exact table: `protocol_violation` during handshake must not automatically classify as a post-handshake violation.
- `stopping` cannot serve as proof that a nonzero exit was expected.

**Disposition:** open.

### X-1, X-2

Confirmed and open. X-1: `Start` returns a retiring process (`supervisor.go:121–125`). X-2: rotation clears the transport without joining the reader (`client.go:206–212`).

## 3. Disposition summary

| Area | Disposition |
|---|---|
| W1 | accept rule; API incomplete |
| W2 | revise; bounded reaping, retained ownership |
| W3 | revise; atomic registration, reader identity |
| W4 | revise; atomic admission, transport ownership |
| Terminal projection | open; event ordering, exact classification |
| X-1 | confirmed, open |
| X-2 | confirmed, open |
| Overall | not ready for acceptance |

## 4. Governance sequence retained

1. Record the reviewer's findings and provenance — this file.
2. Reconcile the findings individually — `RECONCILE-…-round2`.
3. Revise the proposed contract — `CONTRACT-…-rev2`.
4. Independently review the revised contract.
5. Owner acceptance before any design amendment.

This record does not populate the attestation and does not convert author concessions into reviewer-confirmed findings.

## 5. Standing constraints

No r9 edit, no matrix regeneration, no implementation, no tests, no builds. Frozen evidence unchanged. **C1 remains unauthorized.**
