# C1 — Revised Concurrency Contract: Installation Lifecycle

**Date:** 2026-11-08
**Status:** PROPOSED CONTRACT REVISION. Review item, not an accepted amendment. **C1 REMAINS UNAUTHORIZED.**
**Separate from:** attestation (unsigned), corrected ledger, frozen analysis/packet/addendum.

## 1. The chain, not four checks

```text
factory returns p
   → admission decision
   → commit (install) or refusal
   → refusal: reap p, synchronously
   → commit: attach reader/transport under a token
   → retirement: observer owns the terminal
```

Every arrow has an owner. The defect class is an arrow without one.

## 2. Instance token

Allocated **at factory return**, before any admission check:

```text
pending := { token, proc, epochCandidate, state: pending }
```

The token is pointer identity plus a unique id. It is the authentication primitive (E1-05); the epoch is assigned only at commit. A pending instance is not addressable by requests and not returned by `Start`.

## 3. Four windows

| Window | Situation | Required admission | Owner during window | On refusal |
|---|---|---|---|---|
| W1 | existing process returned (`supervisor.go:121–125`) | return only if `state == active` **and** `!shutdownRequested` | supervisor | typed `retiring`/`shutting_down`; never return the handle |
| W2 | factory returned, not installed | commit under `s.mu` only if not shutting down | **installer owns `p`** | installer reaps: `Kill`, `Wait`, `Close`, then returns |
| W3 | installed, reader not yet attached | attach permitted only for this token, before retirement | client holds attach permit | do not attach; let the observer retire |
| W4 | shutdown control permit | one frame, `shutdown` only | holder of the permit | permit consumed; admission stays closed |

**W1 is the current defect (X-1, F):** `Supervisor.Start` returns `s.proc` whenever it is non-nil, including while `retiring`, so a caller can attach to a doomed instance.

**W2 rule (P):** refusing installation does not destroy the process. The party that created `p` reaps it synchronously before returning; a refusal path that returns without `Wait` leaks a process and an observer.

## 4. Attachment

- The client attaches only with a token matching `pending.token`; attachment is a commit, not an implicit consequence of receiving a handle.
- Rotation (`client.go:206–212`) currently clears `transport` without joining the old reader. **X-2 (F):** the previous reader may still be running when a new one attaches. Rotation must `waitForReader` before the new attach, or prove the reader exited.
- After retirement, attachment is refused by token state, not by liveness probing.

## 5. Shutdown permit

`beginShutdown()` returns a permit:

| Property | Requirement |
|---|---|
| scope | exactly one `shutdown` frame |
| start authority | none — `Start` refused while holding or after permit issuance |
| publication authority | none |
| consumption | single use; second use is a typed `shutdown` failure |
| admission | stays closed for the rest of the process, including after termination |

The permit exists because the shutdown frame must be sent through the same transport that admission is closing. It is the only exception, and it cannot reopen ordinary work.

## 6. Intent and terminal classification

Keep three separate fields:

| Field | Owner | Purpose |
|---|---|---|
| `primaryIntent` | transition owner | authorises the transition; first intent wins |
| `observedIntents` | observer | records later facts without changing authority |
| `terminalProjection` | observer | computes the published state from **observed outcome** |

**Rule:** the observed process outcome classifies the terminal; the requesting operation never explains an unrelated failure. A cleanup request cannot turn an unrelated crash into "expected".

| Requesting intent | Observed outcome | Projection |
|---|---|---|
| handshakeFailure{code} | exit 0 or nonzero | `incompatible` / `unavailable` per code |
| shutdown | exit 0 | `absent` |
| shutdown | nonzero | `crashed` + observed code |
| cleanup | exit 0 | `absent` (process gone) |
| cleanup | nonzero **and** no crash observed before termination | `absent` + failure retained |
| protocolViolation | exit 4 | `crashed` + `protocol_violation` |
| none | unexpected exit | `crashed` |

No precedence table is accepted yet; this projection table is a proposal for review.

## 7. New findings (open)

| # | Finding |
|---|---|
| X-1 | `Start` returns a retiring process (`:121–125`); admission is missing |
| X-2 | Rotation clears the transport without joining the reader (`:206–212`) |

X-1, X-2 join L-1…L-9 and W-1…W-14. None closed.

## 8. Status

No r9 edit, no matrix regeneration, no implementation, no tests, nothing built or run. **C1 remains unauthorized.**
