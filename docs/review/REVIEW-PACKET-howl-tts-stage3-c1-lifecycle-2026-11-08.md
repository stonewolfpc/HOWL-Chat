# C1 — Independent Review Packet: Lifecycle Serialization

**Date:** 2026-11-08
**Status:** REVIEW INPUT. Not a design, not an approval. **C1 REMAINS UNAUTHORIZED.**
**Artifact under review:** `ANALYSIS-howl-tts-stage3-c1-lifecycle-serialization-2026-11-08.md`, `bdb68695…`, 193 lines — **frozen**; this packet does not edit it.
**Other frozen artifacts:** disposition `31137b7e…`, reconcile `479ac26c…`, audit `1541332e…`, A1-rev1 `c061e9d5…`, A2-rev1 `26342c5e…`.

The author cannot close any item in this packet. Reviewer dispositions belong in the empty fields; none are filled here.

## 1. What the reviewer is being asked

Six questions, each with the evidence the claim rests on, what would falsify it, and what remains unproven. Verification is static — no execution, no race-detector run, no stress test. Any claim that requires execution is marked so.

## 2. Q1 — Can install and exit race, double-advancing the epoch?

**Evidence (F).** `Start` holds `startMu` across discovery and spawn (`supervisor.go:117–160`); install writes `s.proc` under `s.mu` (`:153–160`); `observeExit` clears `s.proc` only in its own locked section (`:171–173`); `Start` returns early while `s.proc != nil` (`:120–125`).

**Finding (F).** A new install cannot begin before the old exit is observed: `s.proc` is the gate and only `observeExit` clears it. Two installs cannot race. Exit cannot publish for a process it does not own (`:172`).

**But E1 as written is wrong (P).** E1 advances the epoch on installation **and** exit completion. That makes every in-flight terminal event stale at the moment its own instance exits — the request's own `SetState(ready)` would be rejected as stale, which is not the intent. Proposed refinement:

- epoch advances **only on install** (monotonic install counter);
- exit completion **retires** the current epoch (marks it closed, clears `proc`) without incrementing;
- the next install increments.

**Falsifier:** any path that clears `s.proc` without a new install, or installs while `s.proc != nil`.
**Unproven:** whether a clone or restart policy outside this package could install twice. No such path was found.

## 3. Q2 — Can two termination requests race and lose intent?

**Evidence (F).** `Terminate` stores no intent. It writes two enum values (`stopping`/`stopping`, `:256–259`) and waits on the same `done`. A second `Terminate` on the same instance repeats the writes and waits on the same channel; after exit, `proc == nil` returns immediately (`:261–263`).

**Finding (F).** Intent is a boolean marker, not a record. Shutdown, cleanup, and protocol-violation intents are indistinguishable once termination begins, and whichever observer runs sets the terminal state for all of them. This is independent of T1: even with no concurrent writer, the observer cannot know which intent it is completing.

**Consequence:** the five-intent table (disposition §B.5) cannot be implemented against the current supervisor without adding per-instance intent storage.

**Falsifier:** a per-instance intent field, or a termination call that carries intent. Neither exists (**F**, grep).
**Open:** whether intent should be a single value or a set (shutdown racing cleanup). Reviewer question.

## 4. Q3 — Cancellation, shutdown, and synthesis completion concurrently

Three traces were checked:

| Order | Result |
|---|---|
| completion → shutdown | completion publishes `ready`; shutdown sets stopping; exit → `absent`. Consistent |
| shutdown → completion | stopping is overwritten by `ready`/`alive` (`:240`); observer publishes `crashed` (T1) |
| cancel watcher during shutdown | watcher calls `Cancel` → `call` → `Start` (`:389`). If exit was already observed, `Start` sees `proc == nil` and **spawns a replacement** while `closing` is still false (`client.go:669–673` sets `closing` only after `Terminate` returns; `Start` refuses at `:213–216` once it is set) |

**New finding L-9 (F):** shutdown can be defeated by a cancel-triggered respawn. The same is true of `Shutdown`'s own `call`, which may start a helper if the transport exists but the process does not.

**Falsifier:** `closing` set before termination begins, or `Start` refusing while a shutdown is in progress.
**Unproven:** frequency; requires execution to observe.

## 5. Q4 — Resource ownership for stale cleanup

**Evidence (F).** `cleanupRejected` (`client.go:859–864`) takes only `(dir, filename, force)`. It calls `c.sup.ProcessAlive()` and `Terminate` on **whatever process is current**, not on the instance the request used.

**New finding L-7 (F):** if the request's instance has already exited and a replacement was started by `Health`/`Catalog`, the stale request's cleanup terminates the **replacement**. There is no instance token in the cleanup path.

**Rule needed (P):** a cleanup must carry the resource generation it owns — `(instance pointer, epoch, request nonce)` — and may act only when the current generation matches. Proposed guards:

| Resource | Owner | Guard |
|---|---|---|
| process/transport | instance | `s.proc == expected` before close/kill |
| output file | request | exact path **and** request nonce still owns it |
| pending/retired entry | request | request id, already unique per session |
| stderr ring | instance | read-only, labelled by epoch |

`retire(id)` is already safe (unique ids, `:467`). Output deletion is safe while `workMu` serializes synthesis, but not if a future caller cleans up outside `workMu`.

## 6. Q5 — What a ticket may authorize across a restart

The binding must state what survives a restart. Proposed (P):

| Field | Survives restart | Advances on |
|---|---|---|
| subject | yes (identity of caller) | — |
| `targetID` (provider + engine) | yes, must match | identity change |
| `generation` (hash of provider, model identity, config) | yes, must match | any of those changing |
| `epoch` | **no** — restart invalidates | process install |
| operation class | yes | — |
| request nonce | no | per request |

| Operation | Binding checked | After restart |
|---|---|---|
| status / catalog | targetID + generation | valid |
| load / unload | targetID + generation | valid |
| synthesize (new) | targetID + generation | valid |
| cancel of request R | nonce + epoch of R | **rejected** — the request died |
| any "resume"/"continue" | nonce + epoch | **rejected** |

Key point: surviving a restart must never mean authorizing against a **different** identity or configuration. `targetID` and `generation` are compared together; a matching generation under a different target is a mismatch, not a match.

**Falsifier:** a ticket field that is neither identity-bound nor instance-bound, or a consume path that compares only one of the two.

## 7. Q6 — Lock discipline and deadlock

Edges observed in the audited source (F):

```text
workMu  → startMu        (Synthesize/Catalog/Load/Unload → Start)
workMu  → Client.mu      (capability checks, begin, retire)
workMu  → Supervisor.mu  (SetState, SetSelected)
startMu → Client.mu      (Start, :201)
startMu → Supervisor.mu  (sup.Start)
Client.mu → (no other lock held across the call)   (:201–224 release before sup.Start; dispatch unlocks before protocolViolation)
Supervisor.mu → (no other lock)                    (:234–241, :253–260)
```

No reverse edge was found: nothing acquires `workMu` while holding `Client.mu` or `Supervisor.mu`; nothing acquires `Client.mu` while holding `Supervisor.mu`. `Terminate` waits on `done` **without** holding `s.mu`, and `observeExit` closes `done` **while** holding it — safe, because the waiter holds nothing.

**Finding (F):** the graph is acyclic, so no deadlock from lock inversion exists in the audited paths. `Terminate` can still block for the full grace while a caller holds `workMu`; that is latency, not deadlock.

**Falsifier:** a path that calls a `Client` method from inside `s.mu`, or a `Supervisor` method from inside `c.mu`. Neither was found; the reviewer should re-check `Status`/`ProviderStatus` and any future adapter.

## 8. New findings (open)

| # | Finding |
|---|---|
| L-7 | Stale cleanup terminates the replacement instance; no instance token in `cleanupRejected` |
| L-8 | Termination intent is a boolean marker; concurrent or sequential intents collapse |
| L-9 | Shutdown can spawn a replacement: `Cancel`/`Shutdown` → `call` → `Start` before `closing` is set |

L-1..L-6 and W-1..W-14 remain open. No item is closed by this packet.

### 8.1 Citation errata (frozen artifacts not edited)

The reconcile artifact's R-27 row cites `cleanupRejected:856–864`. The comment block begins at 855; the function body spans 859–864, with `Terminate` at 861. Substance unchanged, and the frozen file is left as written so the correction is visible rather than silent.

## 9. Reviewer checklist — what would overturn each conclusion

1. Q1: find a path clearing `proc` without install, or two installs with the gate open.
2. Q1 refinement: show a legitimate terminal that must be rejected after its own instance exits.
3. Q2: show intent is recoverable from existing state (it is not, but prove it).
4. Q3/L-9: find a guard that prevents `Start` after exit while `closing == false`.
5. Q4/L-7: find an instance binding in `cleanupRejected` that this reading missed.
6. Q5: identify a ticket operation that must survive an identity change.
7. Q6: produce a lock-inversion path.

## 10. Status

Independent review requested; the author cannot perform or claim it. No matrix regenerated, no r9 design revision, no epoch implementation. **C1 remains unauthorized.**
