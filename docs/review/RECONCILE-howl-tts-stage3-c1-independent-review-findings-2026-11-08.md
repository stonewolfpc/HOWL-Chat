# C1 — Author Reconciliation of the Independent Review Findings

**Date:** 2026-11-08
**Status:** AUTHOR RECONCILIATION. **Not independent closure.** **C1 REMAINS UNAUTHORIZED.**
**Review text:** **not present in the workspace.** Searched `docs/**` for `E1-01`…`E1-05`, `L-7-01`…`L-7-06`, `NOT READY`, `instance authentication` — no review artifact contains those IDs. Rows below reconcile only the claims quoted by the owner; IDs without quoted text are marked **TEXT REQUIRED** and are not reconstructed.

**Frozen, untouched:** packet `a1c26a1b…`, analysis `bdb68695…`, precision addendum `62e4005e…`, disposition `31137b7e…`, reconcile `479ac26c…`, audit `1541332e…`, A1-rev1 `c061e9d5…`, A2-rev1 `26342c5e…`.

Individual findings stay separate. No row is merged into a theme, and no row is closed by this document.

## 1. E1 — acceptance predicate and event model

### E1-05 — epoch equality is not instance authentication (CONCEDE)

**Reported claim:** an epoch match does not prove the event belongs to the instance being retired.

**Evidence (F):** the current observer authenticates by pointer, not by number: `if s.proc == p` (`supervisor.go:172`). That comparison is the only thing preventing an exit from publishing for a process it does not own. The proposed epoch-only rule would discard it.

**Correction (P):** an instance token captured at install — pointer identity plus a monotonically unique instance id — is the authentication primitive. The epoch remains a label and an ordering device; it never authenticates by itself.

**Residual question:** whether a recycled pointer could ever compare equal. With `proc` cleared at exit (`:173`) and no reuse before reinstall, no; the reviewer should confirm no allocator/pool path violates that.

### E1 (reported gap) — `active + intent != none` permits resurrection (CONCEDE)

**Reported claim:** if termination intent is latched without atomically changing the acceptance state, ordinary completions still publish.

**Evidence (F):** `SetState` writes the state, the active request id, and the failure under `s.mu`, and touches nothing else (`:233–242`). `Terminate` latches only two enum values (`:256–259`). There is no `intent` field, so any `active` instance accepts `SetState(ready)` regardless of a pending termination.

**Correction (P):** termination must move the instance to a distinct `retiring` state **in the same critical section** that records the intent. Acceptance then reads state, not intent separately:

```text
publishable(event) :=
  event.instance == instance.token
  && event.epoch == instance.epoch
  && ( (instance.state == active  && event.role in {ordinary, cleanup})
     || (instance.state == retiring && event.role == exitObserver)
     || (instance.state == retired  && event.role == exitObserver && idempotent) )
```

`ordinary` is never publishable once `retiring`, which closes the gap without relying on intent inspection.

### E1 (reported conflict) — handshake failure path violates the role rules (CONCEDE, PARTIAL)

**Reported claim:** the proposed role rules conflict with an existing production path.

**Evidence (F):** on handshake failure the client terminates, waits for the reader, then publishes the typed state **after** termination (`client.go:248–253`), with a comment stating the intent: preserve `incompatible` rather than report `absent`. Under the role matrix, that publication is an ordinary event against a retiring/retired instance and must be rejected.

**Correction options (P):**

| Option | Change | Consequence |
|---|---|---|
| (a) observer owns the typed terminal | `Terminate(intent=handshakeFailure{code})`; observer writes `incompatible`/`unavailable` itself | removes the post-hoc patch and the T1 exposure on this path; preferred |
| (b) explicit `handshakeFailure` role | allow one post-termination publish for this role | keeps two publishers; preserves the pattern the owner rejected |

Recommended: (a). It makes the handshake result a lifecycle outcome rather than a client override.

### E1-01 … E1-04 — TEXT REQUIRED

The reported review describes five E1 findings; only E1-05 and the two gaps above were quoted. These IDs are reserved and not reconstructed. Each needs its reviewer text, evidence line, and falsifier before any response is recorded.

## 2. L-7 — termination and resource ownership

### L-7 (reported) — check-then-act termination (CONCEDE)

**Evidence (F):** `cleanupRejected` tests `ProcessAlive()` and then calls `Terminate()` (`client.go:860–861`). Between the two, the process can exit and a replacement can be installed, so the termination lands on the replacement.

**Correction (P):** replace both calls with one synchronized operation on the supervisor:

```text
TerminateIf(expectedInstance, intent) bool   // compare token + epoch, set retiring, then close
```

The expected-instance comparison happens inside `s.mu`; the caller never checks liveness first. Return `false` when the instance no longer matches, and the caller must not retry against the current process.

### L-7 (reported) — reader and protocol-violation kill paths (CONCEDE)

**Evidence (F):** `readLoop` kills the current process on frame errors (`client.go:305`, `:312`), and `protocolViolation` does the same (`:366`). Neither passes the process it was reading from.

**Correction (P):** `KillIf(p, intent)` with the same token comparison. `readLoop` already holds `p`, so the token is available at the call site.

### L-7-03 — process ownership and output ownership are different (CONCEDE)

**Reported claim:** a stale request must not kill a replacement, but may still own an output file that needs cleanup.

**Evidence (F):** `cleanupRejected` couples both actions behind one liveness test (`:860–863`).

**Correction (P):** split by resource class; ownership is proven per resource, not per request:

| Resource | Ownership proof | On proof failure |
|---|---|---|
| process / transport | instance token + epoch | never terminate; return |
| output file | request nonce (path alone insufficient) | do not delete |
| pending/retired entry | request id | no-op |
| reader | process instance that started it | never close |

A stale request therefore may delete its own output while being forbidden to touch the process.

### L-7-01, L-7-02, L-7-04 … L-7-06 — TEXT REQUIRED

Reserved; not reconstructed.

## 3. L-8 — concurrent intent semantics

**Reported state:** open. **Evidence (F):** no intent is stored at all (`Terminate` writes two enums, `:256–259`), so concurrent intents cannot be ordered, merged, or reported.

**Question the contract must answer (P):**

1. Is the first intent authoritative, or is there a precedence order?
2. Can a later intent upgrade the terminal result (e.g., shutdown arriving during cleanup)?
3. What does the second caller observe — success, no-op, or a typed "already terminating"?

Recommended: first intent wins, recorded once; later intents are recorded as `supersededBy` for diagnosis and return a typed `terminating` result; precedence applies only to the exit code interpretation, not to the recorded cause. Owner decision, not author's.

## 4. L-9 — serialization, not ordering language

### L-9 (reported) — `Start` already past its check (CONCEDE)

**Evidence (F):** `Start` checks `closing` under `Client.mu` (`:213–216`), releases it, then calls `sup.Start` (`:226`). A shutdown that sets `closing` after that check but before install does not stop the spawn; the check and the install are in different critical sections and different objects.

**Correction (P):** admission must be enforced where installation happens — a `shutdownRequested` flag on `Supervisor`, checked under `s.mu` inside `Start` immediately before `s.proc = p` (`:153–160`). `Client.closing` remains the client-side mirror.

### L-9 (reported) — the shutdown control frame (CONCEDE)

A barrier that refuses all new work also refuses the `shutdown` frame, because `Shutdown` sends it through `call` → `Start` (`client.go:389`, `:657`).

**Correction (P):** `beginShutdown()` issues a **one-shot control permit**:

| Step | Action | Effect |
|---|---|---|
| 1 | acquire control permit | exactly one frame may be sent; `Start` refused for everything else |
| 2 | send `shutdown` using the permit, no `Start` | frame cannot respawn |
| 3 | set `retiring` with intent `shutdown` | observer owns the terminal |
| 4 | consume permit, clear transport | no further frames |

The permit is the exception the reviewer is asking for, and it is bounded to one use.

## 5. Q6 — accept the bounded conclusion, correct the inventory

The conclusion stands as: **no lock-order cycle was identified in the audited acquisition paths.** Not a deadlock-freedom guarantee.

One inventory correction is required: the earlier table listed `startMu → Client.mu` and `startMu → Supervisor.mu` as if they were one acquisition path. They are two separate critical sections with a release between them, and that gap **is** L-9. The corrected inventory must show the release points, not just the edges.

## 6. Disposition ledger (author, not closure)

| Item | Author response | Closed? |
|---|---|---|
| E1-05 | concede | no |
| E1 intent gap | concede | no |
| E1 handshake conflict | concede, option (a) proposed | no |
| E1-01..E1-04 | text required | no |
| L-7 check-then-act | concede | no |
| L-7 kill paths | concede | no |
| L-7-03 | concede | no |
| L-7-01..02, 04..06 | text required | no |
| L-8 | open, questions listed | no |
| L-9 Start race | concede | no |
| L-9 shutdown frame | concede | no |

## 7. Next gate

Author reconciliation is complete only for the quoted findings. The next gate is another independent review of the revised contract, after the reviewer's full text is supplied for the reserved IDs. No r9 edit, no matrix regeneration, no implementation, no tests, nothing built or run. **C1 remains unauthorized.**
