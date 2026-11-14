# C1 — Lifecycle Serialization Analysis (read-only)

**Date:** 2026-11-08
**Status:** ANALYSIS. Not r9 design text. Not an approval of `Supervisor.apply(event)`. **C1 REMAINS UNAUTHORIZED.**
**Frozen:** reconcile `479ac26c…`, audit `1541332e…`, disposition `31137b7e…`, A1-rev1 `c061e9d5…`, A2-rev1 `26342c5e…`, originals `92acde55…`/`0a5ea047…`.

Purpose: establish a concurrency contract for synthesis, shutdown, termination and completion; separate state publication from resource cleanup; give the epoch enough semantics to be evaluated rather than merely named.

No source, test, dependency or configuration changed. Nothing built or run.

## 1. Reachability audit — exact scope and bounded conclusion

Method: `find . -name '*.go' -not -path './.git/*'` for the file set, then a recursive literal search for `audio/tts/helper"` over that same set, plus markers for generated files.

| Measure | Count |
|---|---|
| `.go` files in tree | 350 |
| files with `//go:build` | 28 |
| `.go` files under `build/` | 41 |
| imports of the helper outside the package | **0** |
| `Code generated … DO NOT EDIT` under `internal/audio` | 0 |

Bounded conclusion: **no production import or call path was found in the audited source tree.** Build-tagged files were included in the search; generated production sources were checked for presence and none exist under `internal/audio`. This is a negative search result, not a universal claim about every possible runtime integration — reflection, plugin loading, or a build produced outside this tree are outside its scope.

Consequence for evidence: the ten exported operations are API surface. Every side effect in the inventory describes what the code would do if called. The 21 `SetState` sites are package-internal coupling, not 21 active application publishers.

## 2. Lock inventory

| Lock | Protects | Notes |
|---|---|---|
| `Supervisor.mu` (`:66`) | `proc`, `procDone`, `providerState`, `processState`, `lastFailure`, identity, versions | the publication lock |
| `Supervisor.startMu` (`:68`, `:117`) | discovery + spawn | held across `Discover()` and `Factory()` |
| `Client.mu` (`:89`) | `transport`, `started`, `closing`, `hello`, `pending`, `retired`, `session` | transport/request bookkeeping |
| `Client.startMu` (`:129`, `:198`) | client-level startup | held across handshake |
| `Client.workMu` (`:142`) | catalog, load, unload, synthesis | **request serialization only; not lifecycle** |
| `requestAcceptance.mu` (`:542`) | acceptance state machine per request | not shared |

No lock is shared between publication and request serialization. That is the structural fact behind the interleavings in §4.

## 3. Traced acquisition order

### 3.1 Synthesize (`client.go:679–854`)

1. request validation, no locks;
2. `workMu.Lock()` — released last by defer;
3. `Start` → `startMu` → `mu` (read transport) → release → `sup.Start` → `startMu` → `s.mu`;
4. capability checks: `mu` per check (`requireServiceCapability`, `selectedProvider`);
5. `begin` registers the pending request under `mu`;
6. `sup.SetState(busy, id)` — `s.mu`;
7. `watchCancellation` goroutine started; it may issue `Cancel` with **no** `workMu`;
8. `currentTransport()` under `mu`; `Send` outside any lock;
9. wait on `pending.ch`, timeout, or `ctx.Done()`;
10. terminal `SetState(...)` — `s.mu`, one call per exit path (`:771`, `:783`, `:793`, `:796`, `:803`, `:839`, `:847`);
11. defer: `retire(id)` under `mu`, `state.close()`, then `workMu` release.

### 3.2 Shutdown (`client.go:645–674`)

1. `currentTransport()` under `mu`;
2. `call(ctx, "shutdown", ...)` — `call` may invoke `Start` (`:389`), so **shutdown can spawn**;
3. `sup.Terminate` — `s.mu` to set stopping (`:257–258`), release, then `p.Close()` and wait on `done`;
4. `mu` to set `closing`, `started=false`, `transport=nil`.

No `workMu` at any point.

### 3.3 Termination (`supervisor.go:246–259`)

1. `s.mu`: read `proc`, `procDone`; if `proc != nil`, write `stopping` + `stopping`; release;
2. `proc == nil` → `return nil` **without any state write** (`:261–263`);
3. `p.Close()` outside the lock;
4. wait on `done` / ctx / grace; `Kill()` on the latter two.

### 3.4 Completion (`supervisor.go:166–192`)

1. `p.Wait()` and `p.Stderr()` outside locks;
2. `s.mu`;
3. if `s.proc == p`: clear `proc`, record `waitErr`/`lastStderr`, then branch on `providerState == stopping`;
4. publish `absent`/`exited` (exit 0) or `absent`/`exited` (nonzero) on the stopping branch, or `crashed`/`exited` otherwise;
5. `close(done)` while holding `s.mu`;
6. release.

`done` is closed under the publication lock, which is the only ordering guarantee currently offered to waiters.

### 3.5 Load / Catalog / Unload / Cancel / Health

- `Load` (`:551–596`): `workMu` → `Start` → `SetState(loading)` → `callOn` → `SetState(unavailable|ready)` → `SetSelected`. Publication and identity are separate critical sections.
- `Catalog` (`:517–547`): `workMu` → `Start` → `callOn` → `SetSelected`. No state write.
- `Unload` (`:599–618`): `workMu` → capability (`mu`) → `call` → `SetState(unloaded)`.
- `Cancel` (`:622–643`): no `workMu`, no state write.
- `Health` (`:501–515`): no `workMu`; `call` may `Start`.

## 4. The interleavings that make an epoch necessary

### 4.1 T1 — an orderly termination is republished as `crashed`

```text
T1 (shutdown/cleanup)                 T2 (in-flight synthesis terminal)
Terminate: s.mu, provider=stopping    ...
release s.mu                          SetState(ready): s.mu
p.Close()                               provider=ready
                                        processState=alive (proc != nil, state != stopping)
                                        release s.mu
observeExit: provider != stopping → crashed + "exited unexpectedly"
```

`SetState` unconditionally overwrites the stopping marker (`:240`), and `observeExit` branches only on that marker (`:176`). A termination that began cleanly is therefore reported as an unexpected crash. This is not a benign label error: it changes `LastError`, `Available`, and every downstream decision keyed on `crashed`.

Severity is highest for `cleanupRejected` (`:861`), which terminates from inside a request path while that same request may still be completing.

### 4.2 T2 — a request publishes after the process is gone

```text
S: workMu held, pending.ch waiting
Shutdown: no workMu → Terminate → observeExit publishes absent, proc=nil
S: timeout/terminal → SetState(ready) with proc == nil
```

Result: `providerState=ready`, `processState=absent`, `Available=true` (derived from `ready`), `ProcessAlive=false`. The status is internally contradictory, and no code path corrects it. This is the concrete form of W-11 and W-9.

### 4.3 T3 — a stale request publishes against a new instance

`Health` (`:501`) and `Catalog` (`:520`) may spawn without `workMu`. After an exit, a new instance can reach `unloaded` while an old request is still unwinding; that request's terminal `SetState` then describes an instance it never talked to. Identity happens to match, so identity comparison cannot catch it — only an instance epoch can.

### 4.4 T4 — `Terminate` with no process

`proc == nil` returns without writing state (`:261–263`). A stale `ready` with an absent process survives shutdown unchanged.

## 5. Epoch contract — proposed semantics

Six rules, each answering one owner question. All **P**; none approved.

| # | Rule | Proposal |
|---|---|---|
| E1 | **Advancement** | Monotonic `epoch` in `Supervisor`, advanced under `s.mu` when a spawned process is installed (`Start`, after `Factory` succeeds) and when `observeExit` clears it. Never advanced by the client, never by identity change |
| E2 | **Capture** | Captured once, after `workMu` and `Start` return, before `Send`; stored with the request id and passed to every publication for that request. A ticket captures the current generation at issuance for the same reason |
| E3 | **Validation** | `apply(event, epoch)` compares and publishes in **one** `s.mu` critical section. Mismatch → typed `stale_epoch` failure, no write, no state change. No check-then-publish across a lock boundary |
| E4 | **Termination** | A termination event carries the epoch it observed. It may stop only the instance it names. Exit observation of instance *N* is authoritative for *N* and publishes even if the requesting operation has been abandoned — but nothing may publish provider state on behalf of a request whose epoch is stale |
| E5 | **Cleanup** | Stale work may always release resources: `retire(id)` (so a late response is dropped, `:339`), exact output deletion, reader shutdown, stderr capture. It may not publish provider state, process state, `LastError`, `Loaded`, or `ActiveRequestID` |
| E6 | **Replacement** | A new process instance gets a new epoch even with identical provider and configuration. Identity replacement advances the **generation**, not the epoch |

### 5.1 Epoch and generation are two counters, not one

| | Lifecycle epoch | TargetID generation |
|---|---|---|
| Scope | one process instance | one provider identity/configuration |
| Advances on | spawn, exit completion | identity or configuration change |
| Answers | "is this operation still talking to the instance it started with?" | "is this ticket still valid for this target?" |
| Stale result | drop publication, keep cleanup | refuse consumption |

They share one validation point (consume/publication) but not one advancement trigger. Collapsing them would either invalidate tickets on every restart or let a stale request publish to a replacement process — neither is acceptable. This is the reconciliation the owner asked for between TargetID and epoch ownership.

## 6. Publication vs cleanup vs diagnostics

| Category | Members | Gated by epoch |
|---|---|---|
| Publication | `providerState`, `processState`, `lastFailure`, `activeRequestID`, `Loaded`, availability class, `SetSelected` identity | **yes** |
| Resource | `proc` handle, `procDone`, transport, `pending`, `retired`, reader goroutine, output file, pipes | no — release regardless |
| Diagnostic | `lastStderr`, `waitErr`, stderr ring | recorded for the instance, labelled by epoch |

Rule: a stale operation is allowed to free anything and forbidden to assert anything.

## 7. Indivisibility requirement

`Supervisor.apply(event, epoch)` must be the only writer, and validation and publication must occur in the same critical section. The current code demonstrates the failure mode: `Terminate` publishes `stopping`, releases `s.mu`, closes the process, and a concurrent `SetState` republishes `ready` and `alive` in between (§4.1). Any design that checks the epoch before taking `s.mu` reproduces T1 exactly.

`Status()` already reads the whole tuple once under `s.mu`, so it is internally consistent; `ProviderStatus` (`:952–957`) calls `Start` and then `Status` as two steps, so the returned snapshot may describe a different instance than the one that answered `Start`. That is acceptable only if the status carries the epoch.

## 8. New findings — open

| # | Finding |
|---|---|
| L-1 | A concurrent `SetState` during termination converts an orderly stop into `crashed` (T1) |
| L-2 | A late request terminal publishes `ready` after exit; status becomes `ready` + `ProcessAlive=false` (T2) |
| L-3 | A stale request can publish against a restarted instance (T3) |
| L-4 | `Terminate` with `proc == nil` leaves stale state uncorrected (T4) |
| L-5 | `cleanupRejected` publishes `stopping` from a request path, so cleanup mutates lifecycle state |
| L-6 | `watchCancellation` sends `Cancel` with no `workMu` and no epoch, possibly after termination closed the transport |

None closed here. L-1..L-6 join W-1..W-14 as independently open.

## 9. Decisions requested

| # | Decision | Recommendation |
|---|---|---|
| O-22 | accept E1–E6 as the evaluation contract for the event owner | accept, then re-evaluate `apply(event)` against it |
| O-23 | two-counter model (epoch per instance, generation per identity) | accept |
| O-24 | cleanup/publication split: stale work frees, never asserts | accept |
| O-25 | whether `Terminate` may be called from a request path at all, or only from shutdown/lifecycle | decide before r9 |

## 10. What this does not do

No matrix regeneration. No approval of `Supervisor.apply(event)`. No exit-observer specification. No design edit, no source change, no build, no run.

**C1 remains unauthorized.**
