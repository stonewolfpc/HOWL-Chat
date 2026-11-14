# C1 Design Revision 3 — Capability Negotiation and Classification

**Date:** 2026-11-08  
**Revision identity:** C1-R3  
**Status:** PROPOSED C1 DESIGN REVISION — NOT ADOPTED; C1 UNAUTHORIZED  
**Mode:** Design and documentation only. No source, test, build, runtime, benchmark, or deployment change.  
**Supersedes:** Only the C1 design and C1 acceptance portions of `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md` and the proposed design in `docs/review/DESIGN-howl-tts-stage3-c1-revision-2026-11-08.md`. It does not modify, overwrite, or supersede any historical artifact. The frozen plan remains unchanged and historical.

## 1. Inputs and authority

| Input | Identity |
|---|---|
| Adopted r9 requirements | `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55`; 769 / 112,270 |
| r9 adoption record (frozen) | `22e79a3173cc6b9b4358d80eaa82029189e0fb0663a5ba4504af562153c7f10d`; 121 / 15,398 |
| r9 adoption finalization note (frozen) | `cf1fdfc3eebbcfd3b20ca1ac0b746ae19f3d0716a211bba79924d8beca6b488c`; 21 / 1,348 |
| Governing TTS helper contract | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f`; 2,586 / 182,904 |
| Owner decisions, post-D10 | `677f884926de7e58c3dc7ab33f01340f1fbb5966c1b5c089b1aafd057bcf0fba`; 612 / 49,191 |
| D10 finalization note (frozen) | `84fa78dd214f16f40ab125b544727fce2239e66bb1bb54dffc3c8777df9beb02`; 22 / 1,887 |
| Original C1 design / plan (frozen) | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151`; 711 / 69,168 |
| Prior C1 design revision | `003827f12658d1e556657220754795bff2094b39dca960321a71c3fb302cec6d`; 313 / 22,680 |
| Prior C1 readiness report | `45457897844a5fe90c247bf97225420df0700ea431d2a7ac65cf14617867692d`; 129 / 8,910 |
| Independent review disposition | Muse, `CONDITIONAL — SPECIFIC CORRECTIONS REQUIRED`; findings C1-R-1 through C1-R-9 as supplied for this reconciliation. No separate report file was present in the workspace; each cited requirement was reopened and checked directly. |

The contract controls. The owner decisions record D10(a) = A10a-2 and D10(b) = A10b-2. This revision records no new owner decision and reopens none.

## 2. Finding-by-finding reconciliation

All nine findings are accepted. No finding is rejected or silently narrowed. The technical counterexample was checked against the frozen text and current source before correction.

| Finding | Governing requirement checked | Counterexample verified | Disposition and correction |
|---|---|---|---|
| **C1-R-1 — Start/restart map** | r9 §8 C1 (line 589) and §8.1; contract A2-D4; §17.11.5 / §17.7.2 | Valid. Revision 2 omitted `crashed → backoff → starting`, fourth-attempt exhaustion, normal exit-zero/nonzero/grace classification, and a complete §8.1 health-edge table. | **Accepted.** §8 now supplies the complete map and §8.2 the state-scoped health table. C1-AC-13 and the matrix assert completeness. Executed restart remains C5; executed start/reachability remains C7, as r9 allocates. |
| **C1-R-2 — Resource and I/O evidence** | r9 §9.1, §9.2, §10, §3.2, R8-AC-4/7/8/9 | Valid. Revision 2 listed resource values but not basis, measurement, acceptance threshold, or written conformance determination; its I/O evidence was a single assertion, not per-category records. | **Accepted.** §9.2 supplies every §10 field; §9.3 supplies every §9.1 I/O field for all eight §3.2 categories; §9.4 supplies R8-AC-8 and R8-AC-9. C1-AC-15 requires complete, non-blank records and independent determinations. |
| **C1-R-3 — Asset-specific `unavailable`** | Contract A2-D5 and A2-D6 (lines 2506) | Valid. A live handshaken asset-specific `unavailable` retains the authoritative snapshot; an exited/no-instance `unavailable` clears it. Revision 2 did not make the two states separately testable. | **Accepted.** §3.2 rule 7 and §3.4 distinguish them. C1-AC-18 requires separate live-retain and exited-invalidate tests. |
| **C1-R-4 — Semantic handshake termination** | Contract A2-D4 (line 2505) | Valid. A2-D4 requires an intent-aware exit observer. Current `Client.Start` calls `Terminate`, then `SetState(incompatible)`; if termination times out and the exit arrives later, `Supervisor.observeExit` treats the exit as unexpected and writes `crashed`. That would permit backoff, budget spend, or replacement contrary to A2-D4. | **Accepted.** §6.4 specifies an expected-termination intent latched before termination. Late exit preserves `incompatible`; termination timeout does not reclassify. C1-AC-3 adds timeout and late-exit cases. |
| **C1-R-5 — D10 caller-facing projection** | Owner D10 A10a-2/A10b-2; `errors.go` `Failure.AudioError` | Valid. The owner text requires `ErrCodeSynthesisFailed`, preserved precise code, and bounded detail. Revision 2 named the projection but did not make it an explicit acceptance/test assertion. | **Accepted.** §4.3 gives the exact internal and caller-facing shapes. C1-AC-5, C1-AC-6, and C1-AC-19 require both the helper-boundary failure and the caller-facing `AudioError`. |
| **C1-R-6 — Traceability corrections** | r9 §7.3 original AC-7; §7.1 matrix; R8-AC-7; contract A2-D8 | Valid. Original AC-7 is C8-only no-substitution, not I/O evidence. R8-AC-7 is the I/O rule. Revision 2 wrote `AC-7/R8-AC-7` and claimed A2-D8 coverage without shutdown design content. | **Accepted.** §12 separates original AC-7 (C8-only, no C1 claim) from R8-AC-7. A2-D8 is now an explicit non-substitution boundary: `Unload` is not idle shutdown and cannot terminate the process. |
| **C1-R-7 — Supervisor known-absent capability** | Contract A2-D9 (line 2510); owner D10(b) boundary | Valid. Both unknown and known-absent skip the frame and synthesize no `unsupported_capability`, but they are distinct states and revision 2 combined their test. | **Accepted.** §5 separates the branches; C1-AC-20 and a separately named test cover known-absent. |
| **C1-R-8 — Concurrency coverage** | Original AC-1, AC-2, AC-3; r9 §8 C1 TH-1 focus | Valid. Revision 2 covered gate-versus-exit and handshake-versus-gate but not concurrent direct `Unload`/`Cancel`, optional operations during lifecycle states, or supervisor cancellation during startup/shutdown/recovery. | **Accepted.** §9.1 documents existing serialization and winner rules; C1-AC-21 and the concurrency matrix require the missing cases. No new architecture is invented. |
| **C1-R-9 — Exact protocol and evidence details** | Contract A2-D3; R8-AC-8, R8-AC-9; C1-AC-16 | Valid. The exact `"missing"` key, table order, omission rule, PR-dimension statement, reviewed N/A entries, and static-review assertion were not all explicit. | **Accepted.** §6.1 states the exact key/order/omission. §9.4 states PR impact and reviewed N/A. C1-AC-22 requires the C1-AC-16 static review. |

## 3. Authoritative service snapshot

### 3.1 Decision

`Supervisor` owns a private immutable service snapshot. `Client` retains `hello` for provider descriptors and provider-scoped work, but every direct service-capability gate reads only this snapshot:

```text
serviceSnapshot {
    generation uint64
    known      bool
    caps       map[string]capability   // exactly the five service IDs when known
}
```

A snapshot read returns the generation, `known`, and a copied capability view. No mutable map escapes.

### 3.2 Publication and invalidation

Under the existing `Supervisor.mu`:

1. A new process instance begins a new generation and clears any prior snapshot before start.
2. Only a successful compatible handshake for the current instance publishes `known=true` with all five service IDs.
3. Entering `stopping` invalidates immediately.
4. Observing an unexpected exit/crash invalidates immediately.
5. Entering `incompatible` invalidates immediately and never publishes `known=true`.
6. Entering `backoff` or beginning replacement invalidates; replacement receives a new generation.
7. A live, handshaken, asset-specific `unavailable` **retains** the snapshot. An `unavailable` with no live handshaken instance — discovery failure, factory failure, or exit — **invalidates** it.
8. No path sets `known=true` without a compatible handshake for the current generation.

### 3.3 Validity predicate

The private predicate requires all of:

```text
known
snapshot.generation == current generation
the captured transport belongs to that generation
the process is live
state is not stopping, exited, crashed, incompatible, or backoff
no replacement has begun
```

`unavailable` is not excluded when the helper is live, handshaken, and asset-specific. It is excluded when there is no instance.

### 3.4 No-instance versus live `unavailable`

The snapshot decision is instance-scoped, not label-scoped:

- `SetState(unavailable)` after a typed load failure on a live handshaken helper retains the snapshot.
- `setUnavailable` after discovery/factory failure invalidates it.
- `observeExit` during `unavailable` invalidates it and classifies the unexpected exit under the existing lifecycle rules.
- Tests construct both states separately; a single `unavailable` fixture is not sufficient.

### 3.5 Why a client copy is insufficient

`c.hello` is copied at handshake and is not tied to the current process. `Start` can clear `started` and `transport` after a failed handshake without clearing the old `hello`; `observeExit` does not clear supervisor capabilities. A client-only check can therefore authorize `c.call`, which starts a replacement. `callOnCurrent` and the generation check prevent that.

## 4. Direct operations

### 4.1 Direct `Unload`

1. Read the authoritative snapshot. If invalid, return `unsupported_capability` with `capability=howl.tts.unload`, `snapshot=unknown`; no frame; no `Start`/spawn/restart; no `unloaded` claim.
2. If valid and absent, return `unsupported_capability` with `capability=howl.tts.unload`, `snapshot=known`, `advertised=false`; no frame; no start; no state change.
3. If valid and advertised, use `callOnCurrent`; no `Start`; no replacement. A generation/transport change yields the existing typed failure.
4. Set `TTSStateUnloaded` only after a valid successful `unloaded=true` response.

`Unload` is not idle shutdown. It MUST NOT send `shutdown`, terminate the process, or be treated as satisfying A2-D8. A2-D8 remains C4/C7 work.

### 4.2 Direct `Cancel`

1. Validate the target ID first; malformed target returns `invalid_request` before the snapshot gate.
2. If the snapshot is invalid, return `unsupported_capability` with `capability=howl.tts.cancel`, `snapshot=unknown`; no frame; no start/spawn/restart; no cancellation claim.
3. If valid and absent, return `unsupported_capability` with `capability=howl.tts.cancel`, `snapshot=known`, `advertised=false`; no frame; no start.
4. If valid and advertised, use `callOnCurrent`; no `Start`; a generation/transport change yields the existing typed failure.
5. The direct method does not own the synthesis terminal and does not set `unloaded` or claim `cancelled`.

### 4.3 D10 caller-facing failure

For both unknown direct branches, the helper boundary returns a `Failure` with:

```text
Code   = "unsupported_capability"
Detail = {"capability": "howl.tts.unload" | "howl.tts.cancel",
          "snapshot": "unknown"}
```

`Failure.AudioError()` projects that exact `Failure` through the existing mapping:

```text
AudioError.Code = types.ErrCodeSynthesisFailed
AudioError.Cause = the same *Failure
```

The precise `Failure.Code`, `Message`, `Field`, `Retryable`, and bounded `Detail` remain available through the cause/`AsFailure`; they are not collapsed into the public code. No wire message, new public error code, Wails binding, or new taxonomy is added. The known-absent branch uses the same existing surface with `snapshot=known`.

### 4.4 No-start current-instance send

`callOnCurrent(ctx, op, params, timeout, expectedGeneration)`:

- verifies the captured transport still belongs to `expectedGeneration`;
- returns an existing typed failure if the transport is absent or has changed;
- never calls `Start` or `Supervisor.Start`;
- writes at most one operation frame after the check;
- preserves request IDs, timeout, and cancellation behavior.

## 5. Supervisor cancellation path

The supervisor path is `cancelForSupervisor(ctx, target, acceptance)`, a separate internal operation. It MUST NOT call exported `Client.Cancel`.

1. Record `cancel_requested`.
2. Read the snapshot privately.
3. If known and advertised, send exactly one cancel frame on the current instance through `callOnCurrent`. If the instance changes, skip the frame and continue.
4. **Unknown snapshot:** skip the frame; synthesize no `unsupported_capability`; continue the normal path.
5. **Known snapshot with `cancel` absent:** skip the frame; synthesize no `unsupported_capability`; continue the normal path. This is a distinct branch and test from unknown.
6. Observe normal cancel grace. On no terminal, `KillTree`, observe the single `Wait()`, delete only the exact request output paths, and end under the existing `cancelled`/`timeout` precedence.
7. A cooperative terminal remains authoritative; it is not a missing-terminal timeout.
8. Direct `Cancel` and this path remain separately implemented and separately tested on one helper.
9. Cancellation during startup before the watcher exists sends no frame. During shutdown, snapshot invalidation prevents a cancel frame. During live `unavailable`, an existing watcher continues under the normal path; no recovery is triggered. With no live target, the path sends no frame and performs no new termination beyond the established precedence.

## 6. Capability policy and handshake classification

### 6.1 Service floor and exact missing detail

- Mandatory service capabilities: exactly `howl.tts.load` and `howl.tts.synthesize`.
- Optional service capabilities: `howl.tts.catalog`, `howl.tts.unload`, `howl.tts.cancel`.
- Mandatory control operations: `health`, `shutdown`.
- A missing mandatory service capability is terminal `protocol_incompatible`.
- `detail` is the existing `map[string]string`. The absent mandatory IDs are encoded only under the exact key `"missing"`, comma-separated in table order, with no spaces and no duplicates. The order is `howl.tts.load` then `howl.tts.synthesize`; when both are absent the value is `howl.tts.load,howl.tts.synthesize`.
- The `"missing"` key is omitted when no mandatory ID is absent. No array or new container is added.

### 6.2 Scope validation

- Service array: the five service IDs; provider-only feature capabilities are invalid.
- Provider descriptor: provider capabilities; service-only `catalog`, `load`, `unload`, and `cancel` are invalid.
- `howl.tts.synthesize` may appear in both arrays; the overlap is intentional.
- Unknown IDs, wrong versions, duplicates, duplicate providers, duplicate extensions, and invalid descriptors are semantic `hello` failures.

### 6.3 Classification

After frame/envelope validation, a well-framed `hello` result failing semantic validation is `protocol_incompatible`, not `protocol_violation` and not `invalid_frame`. Raw frame failures remain frame errors; request-schema failures remain `invalid_request`.

### 6.4 Intent-aware incompatible termination

A semantic incompatibility sets a private expected-termination intent under `Supervisor.mu` **before** calling `Terminate`:

```text
terminationIntent = incompatible
```

`observeExit` consumes that intent:

- The exit is expected; it MUST NOT write `crashed`, enter `backoff`, spend restart budget, or trigger replacement.
- After exit, `ProcessState=exited`, `ProcessAlive=false`, and aggregate `State` remains `incompatible`.
- If termination times out, `State` remains `incompatible`; a later exit is still consumed as expected. The timeout is an internal termination diagnostic, not a crash classification.
- A live-but-unreaped process remains tracked; no new process starts. The single `Wait()` owner remains authoritative.
- Only the update-gated return to `absent` under C6/C7 clears `incompatible`.

Normal shutdown classification is separate: exit 0 → `absent`; nonzero exit or grace expiry → `crashed`. The incompatible intent is the explicit exception required by A2-D4.

## 7. C1/C2 boundary

C1 adds only private snapshot/generation state, the validity predicate, the no-start send, intent-aware incompatible termination, and the path-scoped cancellation operation. It adds no public status field.

C2 owns `Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown`, and public `TTSProviderStatus` projection. C2 must not claim C1's lifecycle invalidation or invent a second validity rule. C4–C6 own the full lifecycle transitions; C1's generation invalidation is their prerequisite.

## 8. Complete C1 start/restart map

### 8.1 State and trigger map

This is a documented map. C1 owns the map and structural no-spawn checks; C5 owns executed restart/backoff; C7 owns executed start/reachability.

| State | Start/restart trigger | Direct `Unload` / `Cancel` | Snapshot | Required classification |
|---|---|---|---|---|
| `absent` | Explicit TTS/synthesis, `ProviderCatalog`, or `ProviderStatus` may start/handshake. `GetTTSVoiceCatalog` never starts. Health does not start. | No | Unknown/cleared | No optional spawn. |
| `starting` | No second start; one start owner; await handshake. | No | Unknown | Health observation does not complete handshake or start. |
| `handshaking` | No second start; one terminal handshake outcome. | No | Unknown until compatible publication | Wire health excluded; no transition/retry caused by observation. |
| `unloaded` | No process start; ordinary synthesis may ensure-load under C3. | Known-present only; unknown/absent fail without a frame. | Known if handshaken/live | `unload` does not start a process. |
| `loading` | No process start; one active load under C3. | Serialized behind `workMu`; re-gate after load. | Known, live | Health stays `loading`; no load completion attributable to health. |
| `ready` | No process start. | Known-present only. | Known, live | No optional spawn. |
| `degraded` | No process start. | Known-present only. | Known, live | No repair/unload transition from health. |
| `busy` | No process start. | `Cancel` may send when known-present; `Unload` waits behind `workMu`. | Known, live | Health stays `busy`; no cancellation from health. |
| `stopping` | No direct start. | No | Invalid immediately | Normal exit 0 → `absent`; nonzero exit or grace expiry → `crashed`; health stays `stopping`. |
| `crashed` | Only the next explicit TTS request under §17.11.5/C5: within budget, `backoff` then `starting`; fourth attempt in 10 minutes exhausts to `unavailable`. Status/catalog/health do not restart. | No | Invalid | No restart from status, catalog, health, or direct optional operation. |
| `backoff` | C5 delay expiry (1/2/5 s) enters `starting` with the same key while budget permits; exhaustion enters `unavailable`. | No | Invalid | Health stays `backoff`; no early start or budget change. |
| `incompatible` | No direct start. Only update-gated return to `absent` under C6/C7. | No | Invalid | Expected termination/late exit remains `incompatible`; no `crashed`, backoff, budget, or replacement. |
| `unavailable` — live asset-specific | No synthesis/catalog/voice-list/health/status recovery. Only explicit same-identity `load` on a live helper whose helper-side state is `unloaded`, subject to D1a/D1b. | Known-present only; does not trigger recovery. | **Retained** | Health stays `unavailable`; no recovery/load/start. |
| `unavailable` — exited/no instance | No start; only the applicable explicit recovery path after C6/C7. | No | **Invalidated** | Health is supervisor-level; no start/restart/recovery. |
| `exited` | No start by C1. C5/C6/C7 own the next permitted explicit trigger. | No | Invalid | Exit 0 → `absent`; nonzero/grace expiry → `crashed`, except the A2-D4 incompatible intent. |

No start is permitted before `absent` or `crashed` is observed, or from any other live state. A `load` transition on an already-live helper is not a process start. Direct optional operations never start, spawn, or restart.

### 8.2 State-scoped health no-transition edges

The table records all applicable r9 §8.1 provider states. "Wire" means a completed handshake exists; "supervisor" means no handshaken helper exists.

| State | Health surface | Required no-transition assertion |
|---|---|---|
| `absent` | Supervisor | Remains `absent`; no start, process creation, or handshake. |
| `starting` | Supervisor | Remains `starting`; no second start or handshake completion. |
| `handshaking` | **Wire-health exclusion** | No wire health accepted; no transition or retry attributable to the read. |
| `unloaded` | Wire | Remains `unloaded`; no load, synthesis, or transition. |
| `loading` | Wire | Remains `loading`; no load completion, cancellation, or transition. |
| `ready` | Wire | Remains `ready`; no load, synthesis, or transition. |
| `busy` | Wire | Remains `busy`; no cancellation or synthesis completion. |
| `degraded` | Wire | Remains `degraded`; no repair, reset, unload, or transition. |
| `stopping` | Wire while alive | Remains `stopping`; no reversal, acceleration, or restart. |
| `crashed` | Supervisor | Remains `crashed`; no restart, backoff, budget change, or start. |
| `backoff` | Supervisor | Remains `backoff`; no skipped delay, early start, or budget change. |
| `incompatible` | **Wire-health exclusion** | No wire health accepted; no retry, restart, or state change. |
| `unavailable` — live | Wire | Remains `unavailable`; no recovery, load, start, or transition. |
| `unavailable` — exited/no instance | Supervisor | No start, restart, or recovery. |

If exit occurs during a health observation, the single `Wait()` observer remains authoritative; the exit is not attributed to health. `handshaking` and `incompatible` are the only wire-health exclusions, each grounded in the contract.

## 9. Concurrency, resource, and I/O evidence

### 9.1 Existing serialization and winner rules

C1 adds no new synchronization architecture. It uses the existing boundaries:

| Boundary | Role | C1 consequence |
|---|---|---|
| `Client.workMu` | Serializes `Load`, `Unload`, `Catalog`, and the entire `Synthesize` operation. | Concurrent `Unload` calls serialize. `Unload` issued during `loading` or `busy` waits and re-gates. |
| `Client.startMu` | Serializes one client `Start`/handshake; `Supervisor.startMu` serializes process start. | No second start; direct optional operations never enter either path. |
| `Client.mu` | Guards transport, pending requests, `hello`, and client state. | Transport/generation capture and pending registration are atomic; it is not held across `Send`. |
| `Supervisor.mu` | Guards snapshot publication/invalidation, process state, and termination intent. | Publication and invalidation are single-writer and atomic; readers receive copies. |
| `Process.writeMu` | Serializes frame writes on one transport. | Concurrent direct `Cancel` frames cannot interleave bytes. |
| `Supervisor.Terminate` | Stops/kills and waits on the single observer. | Snapshot invalidation happens before close/wait; cancellation never waits on `workMu`. |

Lock discipline: do not hold `Client.mu` or `Supervisor.mu` across `Start`, `Send`, `Terminate`, or `Wait`; `workMu` may span a request wait because it is the request-serialization boundary, not a state lock. Lock order and any exception are documented in the C1 report. This is a design constraint on the existing code, not a new global lock.

Winner/terminal rules:

- **Unload vs Unload:** serialized by `workMu`; each call re-gates after acquiring it. Each valid call sends at most one frame. The helper's idempotent unload response governs; no `unloaded` state is claimed without a response.
- **Unload vs Load/Synthesize:** serialized by `workMu`; no interleaving. Cancellation precedence remains A2-D9/§17.11.6.
- **Cancel vs Cancel:** not serialized by `workMu`, so cancellation stays responsive. Each has its own request ID; `Process.writeMu` serializes writes. A late acknowledgement cannot reopen or overwrite the target's terminal; a helper terminal error is returned as that typed failure, not converted to success.
- **Direct optional vs start:** if the snapshot is not yet published, the direct call fails closed; it cannot enter `Start`. If a start completes after the gate, `callOnCurrent` uses the captured generation/transport and never starts.
- **Supervisor cancel vs shutdown:** stopping/shutdown invalidates the snapshot; the supervisor path skips the frame and the shutdown/terminal precedence governs.
- **Direct optional vs exit/replacement:** invalidation is atomic; a post-gate exit yields a typed not-started/send failure; no replacement and no false state.
- **Supervisor cancel vs live asset-specific `unavailable`:** the watcher continues under the established grace/termination path; no recovery and no capability error.
- **Supervisor cancel during startup:** no watcher/target exists before request registration; no cancel frame. During recovery with no active synthesis, no new cancel target is created.

### 9.2 Resource-budget record

Every row requires the §10 fields. The conformance column is a **written independent budget-reviewer determination**; the author cannot fill or self-approve it. A blank or unresolved authorized source prevents C1 acceptance, not reviewer discretion.

| Resource | Limit | Authorized source | Basis | Measurement / design proof | Acceptance threshold | Conformance determination |
|---|---|---|---|---|---|---|
| Service snapshot | One snapshot; exactly 5 service IDs when known | Contract A2-D6; r9 §8 C1 | Atomic publication prevents torn/stale capability use | Static count; copied-read race test | 5 IDs; no mutable map escape; one generation | Pending independent reviewer |
| Pending requests | Existing `ClientOptions.MaxPending = 32`; C1 adds no new queue or class | Existing implementation bound, preserved by C1; no C1 increase. Treated as contract-provisional and must be dispositioned under §10 | Prevents unbounded correlation state; C1 direct operations use the existing map | Existing admission path; pending-count test; static diff | Never exceeds 32; no new queue | Pending independent reviewer; unresolved if no authorized source is accepted |
| NDJSON frame | 4,194,304 UTF-8 bytes | Contract §17.2.5 and §4.5 | Bounds parser/allocator exposure | Frame-limit test; max-frame fixture | ≤ 4,194,304 bytes | Pending independent reviewer |
| Synthesis text | 65,536 UTF-8 bytes | Contract §17.2.5 | Existing bound; C1 does not change it | Existing validation test | ≤ 65,536 bytes | Pending independent reviewer |
| Request/correlation ID | 128 bytes | Contract §17.2.5 | Prevents unbounded IDs | Pattern/length test | ≤ 128 bytes | Pending independent reviewer |
| Error detail | 16 entries; key 64; value 256 bytes | Contract §17.2.5 | Bounds D10 diagnostics and wire errors | Detail-bound test | ≤ 16 entries; key/value bounds | Pending independent reviewer |
| Helper processes | 2 app-session cap; C1 optional operations start 0 | Contract §17.11.5 hard cap; r9 §8 C1 | No optional-operation spawn | Process-count and start-counter test | 0 starts from C1 optional paths; global cap ≤ 2 | Pending independent reviewer |
| Goroutines | 0 new; one existing protocol reader per helper | r9 §8 C1; existing supervisor ownership | No new background work | Goroutine count/diff during direct operations | 0 new C1 goroutines | Pending independent reviewer |
| Timers | 0 new; existing request/cancel timers only | r9 §8 C1; contract timeout policy | No new polling or retry loop | Static review and timer-count trace | 0 new C1 timers | Pending independent reviewer |
| Pipes/handles/job objects | 0 new; one existing process transport | Contract single-process ownership; r9 §8 C1 | C1 reuses the current instance | Handle/process review | 0 new handles; no second process | Pending independent reviewer |
| Model/asset memory and ORT | 0 | r9 §8 C1 no model load; contract C1 boundary | Capability negotiation loads no model | No-load instrumentation; static review | 0 model/ORT initialization | Pending independent reviewer |
| Disk writes / temporary output | 0 | r9 §8 C1 no-write evidence; HW-1.1 | C1 writes no files or output | File-write trace; static diff | 0 C1 writes | Pending independent reviewer |
| Log volume | 0 new disk logs; existing stderr ring only | HW-1.1; contract bounded stderr | No per-frame/per-operation logging | Static review; log-call diff | 0 new log sites/cadence | Pending independent reviewer |

No row may be accepted on a bare "bounded" claim. The C1 report must record the reviewer's written conforms/does-not-conform determination and any unresolved R8-AC-4 item.

### 9.3 Application I/O impact record

All eight r9 §3.2 categories are recorded. "Reviewed N/A" must name the independent C1 reviewer and date before acceptance; "N/A" alone is not evidence.

| Path / component | Operation | Trigger / cadence | Necessity | Bound | Category | Measurement / design proof | Reviewed N/A rationale |
|---|---|---|---|---|---|---|---|
| `Client.Start` / `callOnCurrent` / `Supervisor.MarkHandshake` | Read and write NDJSON frames on inherited stdio | Handshake once per instance; one frame per valid direct operation | Capability negotiation and optional operation dispatch | 4,194,304-byte frame; ≤ 1 operation frame per call | Helper lifecycle | Frame counters, send-count tests, max-frame fixture | Not N/A; actual protocol I/O |
| `Client.Unload` / `Client.Cancel` / `cancelForSupervisor` snapshot gates | In-memory read only | Each gate | Authoritative capability decision | Copy of 5-entry snapshot; no file/network | Status and health observation | Static proof of no I/O call; race test | Reviewed N/A for storage I/O; in-memory only |
| Capability/snapshot storage | In-memory read/write | Handshake, invalidation, gate | Single authoritative view | 5 IDs and one generation | Caches | Static bound; no mutable map escape | Not a disk cache; no rescan |
| Settings paths | None | None | C1 changes no setting | 0 | Settings | Static diff | Reviewed N/A; no settings read/write |
| Model/asset preparation | None | None | C1 loads no model/asset | 0 | Model/asset preparation | No-load instrumentation | Reviewed N/A; out of C1 scope |
| Renderer / picker paths | None | None | C1 adds no production import or adapter | 0 | Renderer and picker | Import/static diff | Reviewed N/A; C7 owns production adapter |
| Logging / diagnostics | No new disk log | None added | Existing bounded diagnostics suffice | 0 new log sites; existing stderr ring | Logging and diagnostics | Static log-call diff | Reviewed N/A for new logging; D10 detail stays in memory/status |
| Temporary/output files | None | None | C1 produces no output | 0 | Temporary/output files | No file-path write in C1 | Reviewed N/A; C1 does not synthesize |

No network I/O is added; parent-owned stdio is the only transport. The record is completed by the C1 author and reviewed by the independent stage reviewer. It is not a performance claim.

### 9.4 Performance and applicability

**R8-AC-8 impact statement.** C1 affects PR-1 only through a bounded in-memory snapshot read and generation comparison on direct optional operations; it adds no work to the normal synthesis path. PR-4 is affected because new cancellation/exit/generation races require safety coverage. PR-2 is not affected (no contention policy is added). PR-3 is N/A with review (no adaptive control). No PR measurement, benchmark, pass/fail result, or performance claim is made; PR-1–PR-4 remain inactive under A-1.

**R8-AC-9 reviewed applicability.** No blanket N/A is used.

| Area | C1 applicability | Required reviewed entry |
|---|---|---|
| TH-1.1–TH-1.6 | Applicable | Snapshot/gate atomicity, lock ownership, race/overlap tests, liveness, integrity, shutdown discipline. |
| TH-1.7 | Applicable | C1 report maps each TH-1 requirement to evidence; no unstated exclusion. |
| HW-1.1 / HW-1.2 | N/A for new I/O/writes | Reviewed N/A: no settings/cache/model/output/log write is added; NDJSON is recorded above. |
| HW-1.3 | Applicable | Five-entry snapshot and copied read; no polling or unbounded allocation. |
| HW-1.4 | Applicable | No new queue/process/goroutine/timer; process start count 0 from C1 optional paths. |
| HW-1.5 | N/A | Reviewed N/A: C1 adds no scheduling, priority, or workload control. |
| HW-1.6 | Applicable to lock correctness; no hardware assumption | No model/GPU/accelerator dependency; synchronization behavior is covered by TH-1. |
| HW-1.7 | N/A | Reviewed N/A: no new contention signal or degradation behavior; existing bounded diagnostics unchanged. |
| HW-1.8 | Applicable | No throttle, sleep, delay, or concurrency reduction is added. |
| HW-2.1 | Applicable | No safety behavior is traded for performance. |
| HW-2.2 | N/A | Reviewed N/A: C1 adds no deliberate degradation path. |
| HW-2.3 | N/A | Reviewed N/A: C1 makes no performance claim and runs no performance measurement. |
| HW-2.4 / PR-1–PR-4 | Inactive under A-1; impact statement only | PR-3 N/A with review; no PR suite execution or acceptance result. |

## 10. Revised C1 acceptance criteria

| ID | Criterion |
|---|---|
| **C1-AC-1** | Mandatory service floor is exactly `load` + `synthesize`; `catalog`, `unload`, and `cancel` are optional. Missing mandatory IDs produce `protocol_incompatible`; `detail` uses the exact key `"missing"`, table order `howl.tts.load,howl.tts.synthesize`, no spaces/duplicates, omitted when none are missing. |
| **C1-AC-2** | Service/provider scope, unknown IDs, wrong versions, duplicates, duplicate providers/extensions, and invalid descriptors are validated; shared `synthesize` is accepted; service-only/provider-only misplacement is rejected. |
| **C1-AC-3** | Semantic `hello` failures are terminal `protocol_incompatible`; the process is terminated and `Wait()` observed; state is `incompatible`; no backoff/restart. A termination timeout and a later exit preserve `incompatible`; no `crashed`, budget spend, or replacement. |
| **C1-AC-4** | The supervisor-owned snapshot is published only after a compatible handshake and invalidated on stopping, exit, crash, incompatibility, backoff, and replacement; generation is current and copies do not expose a mutable map. |
| **C1-AC-5** | Direct `Unload` with unknown snapshot returns `unsupported_capability` with `capability=howl.tts.unload`, `snapshot=unknown`; no frame, start/spawn/restart, or `unloaded` claim. |
| **C1-AC-6** | Direct `Cancel` validates target first; unknown snapshot returns `unsupported_capability` with `capability=howl.tts.cancel`, `snapshot=unknown`; no frame, start/spawn/restart, or cancellation claim. |
| **C1-AC-7** | Known-absent direct `Unload` and `Cancel` return `unsupported_capability` with `snapshot=known`, no frame, no start, and no false success. |
| **C1-AC-8** | Known-present direct `Unload` and `Cancel` send exactly one operation frame on the current instance; no additional start; a generation change yields a typed failure. |
| **C1-AC-9** | Stale after a successful handshake is unknown for stopping, exit, crash, incompatibility, backoff, replacement, and no-instance `unavailable`: no frame, no start/spawn, no false state. |
| **C1-AC-10** | Supervisor cancellation with unknown or absent `cancel` sends no frame, synthesizes no `unsupported_capability`, and follows grace/`KillTree`/`Wait()` and existing precedence. |
| **C1-AC-11** | Direct `Cancel` and the supervisor path are separately implemented; the watcher never calls exported `Cancel`; the two-sided test passes on one helper. |
| **C1-AC-12** | No direct optional operation calls `Start`; known-present uses `callOnCurrent`; a generation change cannot start a replacement. |
| **C1-AC-13** | The documented map covers every provider state, live/exited `unavailable`, and process `exited`; it includes `crashed → backoff → starting`, fourth-attempt exhaustion → `unavailable`, exit-zero/nonzero/grace classification, the A2-D4 exception, and every §8.2 health no-transition edge. A structural test asserts coverage and no optional-operation spawn. |
| **C1-AC-14** | Snapshot publication/read is atomic and single-writer; no torn/stale snapshot; lock order and blocking exclusions are documented; race/overlap tests pass. |
| **C1-AC-15** | The stage report contains the complete §10 resource record (all seven fields, with a written independent conformance determination per row), complete §9 I/O records for all eight §3.2 categories (all fields and reviewed N/A entries), the R8-AC-8 PR-dimension/no-claim statement, and the R8-AC-9 reviewed applicability entries. |
| **C1-AC-16** | No production import, model load, disk write, public status field, or new public error code/taxonomy is added by C1. |
| **C1-AC-17** | C1 delivery includes the capability matrix, exact tests, failure logs, resource/I/O evidence, and the C1/C2 boundary statement. |
| **C1-AC-18** | Asset-specific `unavailable` on a live handshaken helper retains the authoritative snapshot; exited/no-instance `unavailable` invalidates it. Separate tests prove both. |
| **C1-AC-19** | A10a-2/A10b-2 caller projection is exact: helper `Failure.Code=unsupported_capability`, bounded detail, and `AudioError.Code=ErrCodeSynthesisFailed` with the same `Failure` in the cause; no frame/start/claim and no new taxonomy. |
| **C1-AC-20** | Supervisor known-absent `cancel` is a separately named branch and test: no frame, no `unsupported_capability`, normal grace/termination/precedence; it is not collapsed with unknown. |
| **C1-AC-21** | Concurrency cases cover concurrent direct `Unload`, concurrent direct `Cancel`, direct optional operations during startup/loading/stopping/recovery, supervisor cancellation during startup/shutdown/recovery, and exit/replacement/snapshot-invalidation overlaps; `workMu`/`startMu` and winner rules are documented. |
| **C1-AC-22** | A static-review assertion verifies C1-AC-16: no production import, no public field/error taxonomy, no model/disk path, and no A2-D8 unload-as-shutdown substitution. |

## 11. Required test matrix

| Case | Setup | Required result |
|---|---|---|
| Positive: optional present | Compatible handshake with `unload` and `cancel` | Direct operations each send one frame; no extra start |
| Positive: scope | Shared `synthesize`; provider features in provider array | Accepted |
| Missing mandatory: one/both | Omit `load`, `synthesize`, or both | `protocol_incompatible`; exact `detail["missing"]`; terminate/`Wait`; `incompatible`; no restart |
| Missing mandatory: none | Both mandatory present | No `"missing"` key |
| Negative: semantic registry | Unknown/duplicate/wrong-version capability; duplicate provider/extension; invalid descriptor | `protocol_incompatible`, not `protocol_violation` |
| Negative: scope violation | Service-only ID in provider array; provider-only feature in service array | `protocol_incompatible` |
| Negative: frame/request separation | Duplicate JSON keys; unparseable/oversized frame; malformed request schema | `invalid_frame`/`frame_too_large`/`invalid_request` as applicable |
| Incompatible termination timeout | Semantic failure; termination exceeds grace; process exits later | `incompatible` retained; no `crashed`, backoff, budget spend, or replacement; `Wait` observed once |
| Incompatible late exit | Semantic failure; delayed process exit after timeout | Expected exit consumed; `ProcessState=exited`, state `incompatible`; no crash/restart |
| Known absent: unload | Valid handshake without `unload` | `unsupported_capability`, `snapshot=known`; no frame/start |
| Known absent: cancel | Valid handshake without `cancel` | `unsupported_capability`, `snapshot=known`; no frame/start |
| Unknown: unload | New client, no handshake | `unsupported_capability`, `snapshot=unknown`; no frame/start/`unloaded` claim |
| Unknown: cancel | New client, valid target | `unsupported_capability`, `snapshot=unknown`; no frame/start/cancellation claim |
| Invalid cancel target | Malformed target, no snapshot | `invalid_request`; no frame/start; gate not reached |
| D10 caller projection | Unknown direct unload and cancel | `AudioError.Code=ErrCodeSynthesisFailed`; cause preserves `unsupported_capability` and exact detail; no new code |
| Stale: exit | Handshake, then process exit | Unknown; no frame/start; no false state |
| Stale: stop | Handshake, then stopping | Unknown; no frame/start |
| Stale: crash | Handshake, then unexpected exit | Unknown; no frame/start; no restart |
| Stale: incompatible | Failed handshake after prior snapshot | Unknown; no frame/start; terminal |
| Stale: backoff/replacement | Constructed generation change | Unknown; no frame/start/replacement |
| Asset-specific unavailable | Live handshaken helper, typed load failure | Snapshot retained; direct known-present gate remains valid; no recovery |
| Exited/no-instance unavailable | Discovery/factory failure or exit | Snapshot invalidated; direct gate fails unknown; no start |
| Supervisor unknown cancel | Active synthesis, unknown snapshot | No frame; no `unsupported_capability`; grace/kill/wait |
| Supervisor known-absent cancel | Active synthesis, snapshot known but `cancel` absent | Separately named: no frame; no capability error; grace/kill/wait |
| Two-sided regression | Same helper, direct vs supervisor | Direct D10 result; supervisor A2-D9 path; no cross-leak |
| Concurrency: concurrent unload | Two direct `Unload` calls | `workMu` serializes; each re-gates; no interleaved frames; no start; no contradictory state |
| Concurrency: concurrent cancel | Two direct `Cancel` calls, same/different targets | Each own request ID; writes serialized; typed ack/error; no fabricated success or terminal overwrite |
| Concurrency: startup | Direct optional during `starting`/`handshaking` | Unknown; no frame/start; handshake unaffected |
| Concurrency: loading | Direct `Unload`/`Cancel` issued during load | `Unload` waits on `workMu` and re-gates; `Cancel` only if valid target; no load/start corruption |
| Concurrency: stopping | Direct optional during stopping | Unknown/invalid; no frame/start |
| Concurrency: recovery | Direct optional during live asset-specific `unavailable` | Known-present behavior; no recovery; exited variant fails unknown |
| Concurrency: supervisor cancellation | Cancellation during startup, shutdown, recovery | No frame when no watcher/snapshot; no exported `Cancel`; existing precedence |
| Concurrency: exit/replacement | Exit or generation change races gate/send | No spawn; typed failure; no false `unloaded`/cancelled; snapshot invalidated |
| Map completeness | Static map check | All states, crash/backoff/exhaustion/exit/health rows present; no optional spawn |
| Evidence completeness | Static evidence check | All resource/I/O fields present; N/A reviewer entries; PR impact; no blank determination |
| Static AC-16 | Production import/static scan | No public field/taxonomy/model/disk path; no unload-as-shutdown |

## 12. Requirements-to-design coverage

| Requirement | Design section | Acceptance / tests |
|---|---|---|
| A2-D1 — capability floor | §6.1 | C1-AC-1; missing-mandatory tests |
| A2-D2 — service/provider scope | §6.2 | C1-AC-2; scope/registry tests |
| A2-D3 — semantic classification and `detail.missing` | §6.1, §6.3 | C1-AC-1/C1-AC-3; exact-detail tests |
| A2-D4 — incompatible terminal, no restart, intent-aware exit | §6.4 | C1-AC-3; timeout/late-exit tests |
| A2-D5/A2-D6 — status and known snapshot | §3, §7 | C1-AC-4/C1-AC-18; live/exited split |
| A2-D7 — known-absent direct `unload`; unload is not termination | §4.1, §6.1 | C1-AC-7; known-absent test |
| A2-D8 — idle shutdown, not unload | §4.1 boundary only; no C1 implementation | No C1 coverage claim; C4/C7 own shutdown; static AC-22 |
| A2-D9 — supervisor cancellation | §5 | C1-AC-10/C1-AC-11/C1-AC-20; three-path tests |
| A2-D10 — no silent no-op | §4, §8 | C1-AC-5/C1-AC-7/C1-AC-9 |
| D10(a) A10a-2 | §4.1, §4.3 | C1-AC-5/C1-AC-19 |
| D10(b) A10b-2 | §4.2, §4.3 | C1-AC-6/C1-AC-19 |
| r9 C1-CAPABILITY | §3, §9.1 | C1-AC-4/C1-AC-14 |
| r9 C1-START-MAP | §8 | C1-AC-13 |
| Original AC-1/AC-2/AC-3 | §9.1 | C1-AC-14/C1-AC-21; concurrency matrix |
| Original AC-4 / R8-AC-4 — budget record | §9.2 | C1-AC-15; seven-field record |
| Original AC-5 — safety under load | §4, §9.1 | C1-AC-12/C1-AC-14/C1-AC-21 |
| Original AC-6 — no hidden background work | §9.1, §9.2 | C1-AC-15/C1-AC-17 |
| R8-AC-7 — application I/O | §9.3 | C1-AC-15; eight-category record |
| R8-AC-8 — performance impact | §9.4 | C1-AC-15; PR-dimension statement, no claim |
| R8-AC-9 — applicability review | §9.4 | C1-AC-15; reviewed N/A table |
| **Original AC-7 — no stage substitution** | **No C1 design content; C8-only** | **Not a C1 acceptance claim** |
| R8-AC-10 — pointer to original AC-7 | No C1 design content; C8-only | Not a C1 acceptance claim |

Bare `AC-7` is never used here for I/O evidence. The I/O rule is R8-AC-7.

## 13. Later-stage dependencies

- **C2:** public `Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown`, and status projection.
- **C3:** ordinary load/ensure-load behavior; C1 adds no load semantics.
- **C4:** full exit intent, normal shutdown, timeout, grace, and stop classification. C1 supplies the incompatible-intent exception and generation invalidation.
- **C5:** executed crash → backoff → starting, 1/2/5 s delays, three-restart/10-minute budget, fourth-attempt exhaustion, and idle close.
- **C6:** identity matching and replacement; C1 increments generation at replacement begin.
- **C7:** production adapter, explicit load reachability, and executed start evidence.

No downstream stage is implemented or accepted here.

## 14. Authorization boundary

C1 remains unauthorized. This revision is design and documentation only. It does not modify source, implement C1, run tests/builds/runtime/benchmarks, deploy, authorize C2–C8, change A2-D9, amend the contract, amend the adopted r9 baseline, amend the frozen plan, or alter any owner decision. Independent review by Muse is required before any C1 authorization request.
