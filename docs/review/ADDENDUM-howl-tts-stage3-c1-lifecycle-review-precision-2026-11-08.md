# C1 — Review Precision Addendum: Retirement, Shutdown Barrier, Cleanup Ownership

**Date:** 2026-11-08
**Status:** REVIEW INPUT ADDENDUM. E1 correction is a **review item, not an accepted amendment**. **C1 REMAINS UNAUTHORIZED.**

**Frozen, not edited:**

| Artifact | Hash |
|---|---|
| lifecycle analysis | `bdb68695…` |
| review packet | `a1c26a1b…` |
| disposition | `31137b7e…` |
| reconcile | `479ac26c…` |
| writer audit | `1541332e…` |
| A1-rev1 / A2-rev1 | `c061e9d5…` / `26342c5e…` |

This addendum adds precision only. It closes nothing and approves nothing.

## 1. E1 refinement — retirement must be a checked state

The refinement (advance on install, retire on exit, increment on next install) is directionally right but insufficient alone: after retirement the epoch number is unchanged, so a late completion carrying that number still matches numerically.

**Required representation (P):**

```text
instance := { epoch, state: active | retired, intent, retiredAt }
```

`epoch` is the install counter. `state` is what rejects stale work after exit. `intent` is the termination cause (L-8). `retiredAt` exists only for diagnosis.

**Acceptance rule (P):** an event is publishable only if it satisfies **both** conditions:

1. `event.epoch == instance.epoch`;
2. the event's role is permitted for `instance.state`.

| Event role | `active` | `retired` |
|---|---|---|
| exit observer for this instance | allowed — performs the retirement | allowed, idempotent |
| ordinary request terminal (`ready`/`degraded`/`crashed`) | allowed | **rejected** |
| `SetSelected` / identity publish | allowed | **rejected** |
| cleanup / resource release | allowed | allowed — **no publication** |
| new install | not an event on this instance; requires `state == retired` and a cleared handle | — |

The exit observer is the one role allowed to act on a retired instance, because retiring it is the act itself. Every other late event is rejected on `state`, not on the number.

**Review question:** should a *second* exit observation of the same epoch be an error or a no-op? Recommended: no-op, because `observeExit` already guards on `s.proc == p` (`supervisor.go:172`).

**Falsifier:** any acceptance path that compares only the epoch integer.

## 2. L-9 — shutdown admission barrier

The defect is ordering: `closing` is set after `Terminate` returns (`client.go:669–673`), while `Cancel` → `call` → `Start` can run throughout.

**Required ordering (P):**

| Step | Action | Lock | Invariant |
|---|---|---|---|
| 1 | `beginShutdown()` sets `closing=true` and a monotonic `shutdownReason` | `Client.mu` | no new work admitted |
| 2 | reject `begin()` for any new request with `shutdown` | `Client.mu` | no new pending entries |
| 3 | cancel/cleanup paths observe `closing` before `call`/`Start` | `Client.mu` | L-9 path closed |
| 4 | `Terminate(intent=shutdown)` | `Supervisor.mu` | stopping published once |
| 5 | clear transport, reader, session | `Client.mu` | no reuse |

`Start` already refuses once `closing` is set (`:213–216`); the fix is that step 1 happens **before** any call or termination, not after. An in-flight `Cancel` may still send on a transport that is closing; that is a typed failure, not a lifecycle change, and must not respawn.

**Idempotency:** a second `Shutdown` returns the first result, does not re-terminate, and does not reopen admission. Reopening is out of C1 scope.

**Falsifier:** any spawn path reachable after step 1.

## 3. L-7 — cleanup must prove ownership of the specific resource

Epoch equality is not enough; the resource must be proven owned. Four encounters, four required outcomes:

| Encounter | Required behaviour | Proof needed |
|---|---|---|
| owned by the original request, instance still current and active | release, then request termination through the lifecycle owner | request nonce + `s.proc == expected` |
| transferred to another operation (transport, reader, output path) | **do nothing**; the new owner is responsible | transfer record written atomically before handoff |
| belongs to a replacement process | **do nothing**; never terminate the replacement | `s.proc == expected` fails → stop |
| partially created, ownership never committed | treat creator as owner, release idempotently | creation record keyed by nonce |

**Rule (P):** cleanup proves ownership of the particular resource it is about to destroy; it never infers ownership from "the request is mine" or "the epoch matches".

**Current exposure (F):** `cleanupRejected` (`client.go:859–864`) has no instance parameter and terminates whichever process is current, so the replacement case is reachable. Output deletion by exact path is safe today only because `workMu` serializes synthesis; that is an accident of the current caller, not an ownership rule.

## 4. Q6 — bounded conclusion

Replace any stronger phrasing with:

> No lock-order cycle was identified in the audited acquisition paths.

Explicitly not claimed: deadlock freedom. The reviewer should include deferred callbacks and goroutine handoffs, which the static trace cannot fully order:

- `observeExit` goroutine (`supervisor.go:162`) publishes while holding `s.mu` and closes `done`;
- `readLoop` goroutine (`client.go:236`) calls `failAll` then `Kill`, sequentially;
- `watchCancellation` goroutine (`client.go:870–872`) calls `Cancel` after the request may have retired;
- `Terminate` waits on `done` while holding no lock, which is why the close-under-`s.mu` does not deadlock.

## 5. Priority for the independent reviewer

| Item | Why first |
|---|---|
| L-7 | cross-instance termination of a replacement |
| L-8 | intent cannot be reconstructed; blocks the five-intent table |
| L-9 | shutdown admits a new process |
| E1 retirement (§1) | numeric match is not acceptance |
| Q6 | acyclic observed, full-path verification pending |
| ticket survival | restart persistence must not cross identity/configuration |

**Disposition fields (reviewer, intentionally empty):**

| Item | Disposition | Rationale |
|---|---|---|
| L-7 | | |
| L-8 | | |
| L-9 | | |
| E1 refinement | | |
| Q6 conclusion | | |

## 6. Status

E1 is a review item, not an accepted amendment. The analysis and packet remain frozen exactly as reported. No matrix regeneration, no r9 revision, no epoch implementation, nothing built or run. **C1 remains unauthorized.**
