# C1 Design Revision — Capability Negotiation and Classification

**Date:** 2026-11-08  
**Status:** PROPOSED C1 DESIGN REVISION — NOT ADOPTED; C1 UNAUTHORIZED  
**Mode:** Design and documentation only. No source code, test, build, runtime, benchmark, or deployment change.  
**Supersedes:** Only the C1 design and C1 acceptance portions of `docs/review/DESIGN-howl-tts-stage3-corrective-sequencing-2026-11-07.md`. The frozen plan remains unchanged and historical; this revision does not amend the governing contract, the r9 adoption record, the owner-decision record, or any other stage.

## 1. Inputs and authority

| Input | Identity |
|---|---|
| Adopted r9 requirements | `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55`; 769 / 112,270 |
| r9 adoption record (frozen) | `22e79a3173cc6b9b4358d80eaa82029189e0fb0663a5ba4504af562153c7f10d`; 121 / 15,398 |
| Governing TTS helper contract | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f`; 2,586 / 182,904 |
| Owner decisions, post-D10 | `677f884926de7e58c3dc7ab33f01340f1fbb5966c1b5c089b1aafd057bcf0fba`; 612 / 49,191 |
| D10 finalization note | `84fa78dd214f16f40ab125b544727fce2239e66bb1bb54dffc3c8777df9beb02`; 22 / 1,887 |
| Original C1 design / plan | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151`; 711 / 69,168 |
| Prior C1 readiness review | `e8b077485eea21af5312417d55e33a2d00803d7b56b6013b4929e2c9998da9eb`; 158 / 17,126 |

The contract controls. The owner decisions record D10(a) = A10a-2 and D10(b) = A10b-2. This revision records no new owner decision.

## 2. Purpose and scope

This revision closes C1-1 through C1-7 from the prior readiness review:

- C1-1: authoritative snapshot validity;
- C1-2: stale pre-D10 C1 text;
- C1-3: supervisor cancellation path still calling exported `Cancel`;
- C1-4: mandatory/optional capability policy, scope validation, and semantic classification;
- C1-5: private C1 validity versus public C2 status;
- C1-6: missing C1 start/restart map;
- C1-7: missing C1 evidence obligations.

It does not implement C1. It defines the design, acceptance criteria, tests, and evidence required for an independent review and a later stage-authorization decision.

## 3. Design decision — supervisor-owned service snapshot with a generation

### 3.1 Decision

The authoritative service-capability snapshot is owned by `Supervisor`, not by the exported `Client` method or by a copied `hello` value. The supervisor publishes a private, immutable snapshot object:

```text
serviceSnapshot {
    generation uint64
    known      bool
    caps       map[string]capability   // all five service IDs when known
}
```

The `Client` retains `hello` for provider descriptors and provider-scoped work, but direct service-capability gating reads only `Supervisor`’s snapshot. This is the recommended technical solution because the supervisor owns process start, stop, exit, incompatibility, and replacement; a client copy cannot be invalidated atomically at those boundaries.

### 3.2 Publication and invalidation

`Supervisor` must observe these rules under its existing mutex:

1. On entering a new process instance, invalidate any previous snapshot and increment `generation`; then start the process and handshake.
2. On a compatible handshake, publish a new snapshot with all five service IDs present or absent, `known=true`, and the current `generation`.
3. On entering `stopping`, invalidate immediately (`known=false`, clear `caps`, increment generation).
4. On observed exit or crash, invalidate immediately and increment generation.
5. On a semantic handshake failure that enters `incompatible`, invalidate; the snapshot must never survive into `incompatible`.
6. On entering `backoff` or beginning a replacement, invalidate; a replacement must have a new generation.
7. On an asset-specific `unavailable` on a live, handshaken instance, retain the snapshot; on an exited/no-instance `unavailable`, invalidate. C2/C4 own the public and lifecycle projections, but C1’s private predicate must already fail closed for the no-instance case.
8. No code may set `known=true` except a successful compatible handshake for the current generation.

A snapshot read returns a copy plus generation; no caller receives the mutable map.

### 3.3 Validity predicate

`Client` must use this predicate before any direct optional operation:

```text
known == true
AND snapshot.generation == current supervisor generation
AND current process is alive
AND aggregate state is not stopping, exited, crashed, incompatible, or backoff
AND no replacement has begun
```

The predicate is private to C1. It is not the public `ServiceCapabilitiesKnown` field. A live asset-specific `unavailable` may remain valid when the helper instance is alive and the handshake is retained; C1 does not need to create that state, but the predicate must not falsely invalidate it if it is constructed by a test.

### 3.4 Why not a client-only `hello` check

`c.hello` is a value copied at handshake and is not tied to the current process instance. `Start` currently clears `started` and `transport` after a failed handshake but does not clear an older `c.hello`, and `observeExit` does not clear supervisor capabilities. A client-only check can therefore authorize a direct operation after exit and then call `c.call`, which calls `Start` and may spawn a replacement. The generation-checked supervisor snapshot prevents that.

## 4. Direct operations

### 4.1 Direct `Unload`

`Client.Unload` must:

1. Read the authoritative snapshot; if invalid, return the precise `Failure` code `unsupported_capability` with bounded detail `capability=howl.tts.unload`, `snapshot=unknown`; write no frame; call no `Start`; set no state.
2. If valid but `howl.tts.unload` is absent, return `unsupported_capability` with `capability=howl.tts.unload`, `snapshot=known`, `advertised=false`; write no frame; call no `Start`; set no state.
3. If valid and advertised, use a no-start current-instance send path (`callOnCurrent`) rather than `c.call`. If the instance changes before the frame is written, return the existing `helper_crashed`/not-started failure; do not start a replacement and do not report `unloaded`.
4. Set `TTSStateUnloaded` only after a valid successful response with `unloaded=true`.

### 4.2 Direct `Cancel`

`Client.Cancel` must:

1. Validate the target ID first; a malformed target returns `invalid_request` before the snapshot gate.
2. Read the authoritative snapshot; if invalid, return `unsupported_capability` with `capability=howl.tts.cancel`, `snapshot=unknown`; write no cancel frame; call no `Start`; make no cancellation claim.
3. If valid but `howl.tts.cancel` is absent, return `unsupported_capability` with `capability=howl.tts.cancel`, `snapshot=known`, `advertised=false`; write no frame; call no `Start`.
4. If valid and advertised, use `callOnCurrent`; if the instance changes, return the existing typed failure and do not start a replacement.
5. Do not set `unloaded` or report `cancelled` from the direct method itself; the synthesis request owns its terminal state.

### 4.3 Direct failure surface

Both direct failures use the existing `unsupported_capability` code and the existing `AudioError` projection. `capability=howl.tts.unload|cancel` and `snapshot=unknown|known` are bounded internal `Failure.Detail` entries; they add no wire message, no new public error code, and no new public taxonomy. A malformed `Cancel` target remains `invalid_request` before this gate.

### 4.4 No-start current-instance send

`callOnCurrent(ctx, op, params, timeout, expectedGeneration)` must:

- lock the client transport state and verify the transport belongs to `expectedGeneration`;
- return an existing typed failure when the transport is absent or the generation changed;
- never call `Start` or any `Supervisor.Start`;
- write at most one operation frame after the check;
- preserve the existing request ID, timeout, and cancellation behavior.

This is the technical guarantee that an unknown or stale snapshot cannot spawn a helper.

## 5. Supervisor cancellation path

The supervisor path is a separate internal operation, `cancelForSupervisor(ctx, target, acceptance)`, and must not call exported `Client.Cancel`.

On context cancellation:

1. Record `cancel_requested`.
2. Read the authoritative snapshot privately.
3. If the snapshot is known and advertises `howl.tts.cancel`, send the cancel frame on the current instance using the no-start path. If the instance changes, skip the frame and proceed.
4. If the snapshot is absent or unknown, skip the frame without synthesizing `unsupported_capability`.
5. Observe the normal cancel grace. If no terminal arrives, `KillTree`, observe `Wait()`, and delete only the request’s exact output paths. The request ends `cancelled` or `timeout` under the existing precedence rules.
6. If a terminal arrives, existing precedence still governs; a cooperative terminal is not a missing-terminal timeout.
7. The direct `Cancel` method and the supervisor path remain separately implemented and separately tested on the same helper.

## 6. Capability policy and handshake classification

### 6.1 Service floor

- Mandatory: `howl.tts.load`, `howl.tts.synthesize`.
- Optional: `howl.tts.catalog`, `howl.tts.unload`, `howl.tts.cancel`.
- `health` and `shutdown` remain mandatory control operations.
- A missing mandatory service capability is terminal `protocol_incompatible`.
- `detail.missing` contains only missing mandatory IDs, in table order, comma-separated with no spaces or duplicates.

### 6.2 Scope validation

- Service array: the five service IDs; provider-only feature capabilities are invalid.
- Provider descriptor: provider capabilities; service-only `catalog`, `load`, `unload`, and `cancel` are invalid.
- `howl.tts.synthesize` may appear in both arrays; this overlap is intentional.
- Unknown capability IDs, wrong versions, duplicates, duplicate providers, duplicate extensions, and invalid descriptors are semantic `hello` failures.

### 6.3 Classification

After a frame has been parsed and its envelope validated, a well-framed `hello` result that fails semantic validation is `protocol_incompatible`, not `protocol_violation` and not `invalid_frame`.

The adapter must terminate the helper, observe the single `Wait()`, set aggregate state `incompatible`, and enter no `backoff`, spend no restart budget, and start no replacement. Raw JSON/frame failures remain `invalid_frame`/`frame_too_large`; request-schema failures remain `invalid_request`.

## 7. C1/C2 boundary

C1 adds only private supervisor/client state and accessors needed for validity and gating. It does not add public status fields.

C2 owns:

- `Loaded bool`;
- `ServiceCapabilities []TTSCapability`;
- `ServiceCapabilitiesKnown bool`;
- public status projection and the single-source matrix.

C1 must not claim those fields in its acceptance. C1 may expose a private snapshot read used by C2 later. C2 must not claim transition ownership or invent its own validity rule; C1’s generation invalidation is the prerequisite for safe gating, while C4–C6 own the full lifecycle transitions.

## 8. C1 start/restart map

This map is the C1 deliverable; C5 and C7 own executed restart and reachability evidence.

| State | Permitted start / restart from C1 work | Direct `Unload` / `Cancel` | Notes |
|---|---|---|---|
| `absent` | Explicit TTS/synthesis, `ProviderCatalog`, or `ProviderStatus` may start/handshake (r9/§17.7.2); `GetTTSVoiceCatalog` may not | No | No direct optional spawn; no model load |
| `starting` / `handshaking` | No second start; wait for the single handshake outcome | No | One start owner; no parallel helper |
| `unloaded` | No process start; ordinary synthesis may ensure-load under C3 | Known-present only; unknown/absent fail without a frame | C1 does not implement ordinary load |
| `loading` | No process start; one active load under C3 | No direct optional operation | C3 owns load |
| `ready` / `degraded` | No process start | Known-present only; unknown/absent fail without a frame | C1 gates only |
| `busy` | No process start | Direct `Cancel` may be sent when known-present; `Unload` is not an idle operation | Precedence remains A2-D9/§17.11.6 |
| `stopping` | No direct start | No | Snapshot already invalid |
| `crashed` | Only the next explicit TTS request under §17.11.5/C5, within budget | No | Status/catalog/health do not restart |
| `backoff` | Delay expiry enters `starting` under C5, same key | No | C5 owns timer; C1 owns no spawn |
| `incompatible` | No direct start; only the update-gated return to `absent` under C6/C7 | No | No restart/budget |
| `unavailable` | No synthesis/catalog/voice-list/health/status recovery; only explicit same-identity `load` on a live helper under D1a/D1b | No | C1 documents; C3/C7 implement |
| `exited` | No process start from C1; next explicit TTS/C5 | No | No stale snapshot |

## 9. Concurrency, resource, and I/O design

### 9.1 Lock and publication rules

- Supervisor snapshot publication and invalidation occur under `Supervisor.mu`.
- Client transport and generation capture occur under `Client.mu`.
- Do not hold `Client.mu` while calling `Supervisor.Start` or `Supervisor.Terminate`.
- Direct gate reads the supervisor snapshot, releases it, then captures the client transport/generation; never hold both locks across a start, send, or wait.
- Snapshot reads return a copy; no mutable map escapes.
- Generation increments on new process start, invalidation, exit, stop, incompatibility, backoff, and replacement begin.

### 9.2 Resource budget

C1 must record and preserve these existing bounds; it adds no unbounded worker, queue, timer, process, or buffer:

| Resource | Bound / source |
|---|---|
| One NDJSON frame | 4,194,304 UTF-8 bytes (`MaxFrameBytes`) |
| Synthesis text | 65,536 UTF-8 bytes (`MaxTextBytes`) |
| Request/correlation ID | 128 bytes (`MaxIDBytes`) |
| Error detail | 16 entries (`MaxErrorDetail`) |
| In-flight requests | 32 (`ClientOptions.MaxPending`) |
| Helper processes | 2 app-session cap; C1 starts none from optional operations |
| Protocol reader | One per helper process; owned by client |
| Model memory / ORT | None in C1; no model load |
| Disk writes | None in C1; stdout/stderr are inherited pipes |
| Network | None; parent-owned stdio only |

### 9.3 I/O impact

C1 adds no file, network, renderer, settings, logging, or status I/O. Its only I/O is the existing NDJSON request/response traffic on the inherited helper stdio. The C1 evidence record must state explicit N/A for each untouched I/O category.

### 9.4 TH-1 / HW-1 / performance

- TH-1: capability publication is atomic and single-writer; gate/read races are tested.
- HW-1: no native model, ONNX, GPU, or asset work in C1; N/A with reason.
- HW-2: no model load or resource residency in C1; N/A with reason.
- Performance: no normal-operation performance claim; PR-1–PR-4 remain inactive under A-1; C1 records no benchmark claim.

## 10. Revised C1 acceptance criteria

| ID | Criterion |
|---|---|
| **C1-AC-1** | Mandatory service floor is exactly `load` + `synthesize`; `catalog`, `unload`, and `cancel` are optional; missing mandatory yields `protocol_incompatible`. |
| **C1-AC-2** | Service/provider scope, unknown IDs, versions, duplicates, and invalid descriptors are validated; `synthesize` overlap is accepted. |
| **C1-AC-3** | Semantic `hello` failures are terminal `protocol_incompatible`; process is terminated, `Wait()` observed, state `incompatible`, no restart/backoff. |
| **C1-AC-4** | The supervisor-owned snapshot is published only after a compatible handshake and invalidated on stopping, exit, crash, incompatibility, backoff, and replacement. |
| **C1-AC-5** | Direct `Unload` with unknown snapshot returns `unsupported_capability` with `snapshot=unknown`; no frame, no start/spawn/restart, no `unloaded` claim. |
| **C1-AC-6** | Direct `Cancel` validates target first; with unknown snapshot it returns `unsupported_capability` with `snapshot=unknown`; no frame, no start/spawn/restart, no cancellation claim. |
| **C1-AC-7** | Known-absent direct `Unload`/`Cancel` return `unsupported_capability`; no frame; no start; no false success. |
| **C1-AC-8** | Known-present direct `Unload`/`Cancel` send exactly one operation frame on the current instance; no additional start. |
| **C1-AC-9** | Stale after a successful handshake is treated as unknown: no frame, no start/spawn, no unloaded/cancelled claim. |
| **C1-AC-10** | Supervisor cancellation path with absent/unknown `cancel` sends no frame, synthesizes no `unsupported_capability`, and follows grace/`KillTree`/`Wait()` and precedence. |
| **C1-AC-11** | Direct `Cancel` and the supervisor path are separately implemented; the watcher never calls exported `Cancel`; two-sided test passes on one helper. |
| **C1-AC-12** | No direct optional operation calls `Start`; known-present uses a no-start current-instance send; a generation change yields a typed failure, never a replacement. |
| **C1-AC-13** | C1 documents the state-scoped start/restart map and demonstrates no optional-operation spawn. |
| **C1-AC-14** | Capability publication/read is atomic; no torn snapshot; lock order is documented; race/overlap tests are included. |
| **C1-AC-15** | C1 records resource bounds, I/O impact, performance-preservation statement, and TH-1/HW-1/HW-2 applicability. |
| **C1-AC-16** | No production import, no model load, no disk write, no public status field, and no new public error taxonomy are added by C1. |
| **C1-AC-17** | C1 delivery includes the capability matrix, exact tests, failure logs, and a statement of the C1/C2 boundary. |

## 11. Required test matrix

| Case | Setup | Required result |
|---|---|---|
| Positive: optional present | Valid handshake with `unload` and `cancel` | Direct operations send one frame each; no extra start; normal responses |
| Positive: service/provider scope | Shared `synthesize`, provider features in provider array | Accepted; no semantic failure |
| Negative: missing mandatory | Omit `load`, `synthesize`, or both | `protocol_incompatible`; terminate/Wait; `incompatible`; no restart |
| Negative: unknown/duplicate/wrong version | Invalid capability/descriptor cases | `protocol_incompatible`; no `protocol_violation` for semantic cases |
| Negative: scope violation | Service-only capability in provider array or provider-only in service array | `protocol_incompatible` |
| Known absent: unload | Valid handshake without `unload` | `unsupported_capability`; no frame; no start |
| Known absent: cancel | Valid handshake without `cancel` | `unsupported_capability`; no frame; no start |
| Unknown: unload | New client, no handshake | `unsupported_capability` with `snapshot=unknown`; no frame/start |
| Unknown: cancel | New client, valid target | `unsupported_capability` with `snapshot=unknown`; no frame/start |
| Invalid cancel target | Malformed target, no snapshot | `invalid_request`; no frame/start |
| Stale: exit | Handshake, then process exit | Unknown; no frame/start; no `unloaded`/cancelled |
| Stale: stop | Handshake, then stopping | Unknown; no frame/start |
| Stale: crash | Handshake, then unexpected exit | Unknown; no frame/start; no restart |
| Stale: incompatible | Failed handshake after prior snapshot | Unknown; no frame/start; terminal |
| Stale: backoff/replacement | Constructed generation change | Unknown; no frame/start; no replacement |
| Supervisor absent/unknown cancel | Active synthesis, capability omitted/unknown | No cancel frame; no `unsupported_capability`; grace/kill/wait |
| Two-sided regression | Same helper direct vs supervisor | Direct returns D10 result; supervisor follows A2-D9 |
| Concurrency: gate vs exit | Exit/transport close races with direct call | No spawn; typed failure; no false success |
| Concurrency: handshake vs gate | Snapshot publication and read overlap | No torn snapshot; one terminal handshake |
| Resource/I/O | C1 evidence review | Bounded values recorded; no model/disk/network; N/A entries explicit |

## 12. Requirements-to-design coverage

| Requirement | Design section | Acceptance / tests |
|---|---|---|
| A2-D1 — capability floor | §6.1 | C1-AC-1; missing mandatory tests |
| A2-D2 — service/provider scope | §6.2 | C1-AC-2; scope tests |
| A2-D3 — semantic `hello` classification | §6.3 | C1-AC-3; negative matrix |
| A2-D4 — incompatible terminal/no restart | §3.2, §6.3 | C1-AC-3/C1-AC-4; exit/Wait tests |
| A2-D5/A2-D6 — status and known snapshot | §3, §7 | C1-AC-4; C2 boundary; no public status in C1 |
| A2-D7 — known-absent unload | §4.1, §6.1 | C1-AC-7 |
| A2-D8 — shutdown, not unload | §8 | C1-AC-13; no unload substitution |
| A2-D9 — supervisor cancellation | §5 | C1-AC-10/C1-AC-11 |
| A2-D10 — no silent no-op | §4, §8 | C1-AC-5/C1-AC-7/C1-AC-9 |
| D10(a) A10a-2 | §4.1 | C1-AC-5/C1-AC-9 |
| D10(b) A10b-2 | §4.2 | C1-AC-6/C1-AC-9 |
| r9 C1-CAPABILITY | §3, §9.1 | C1-AC-4/C1-AC-14 |
| r9 C1-START-MAP | §8 | C1-AC-12/C1-AC-13 |
| AC-1/AC-2/AC-3 — concurrency/locks/race evidence | §9.1 | C1-AC-14; race/overlap matrix |
| AC-4/R8-AC-4 — resource budget | §9.2 | C1-AC-15; budget record |
| AC-5 — safety under load | §4, §9.2 | C1-AC-12/C1-AC-14 |
| AC-6 — no hidden background work | §9.2 | C1-AC-15; ownership review |
| AC-7/R8-AC-7 — I/O evidence | §9.3 | C1-AC-15 |
| R8-AC-8 — performance preservation | §9.4 | C1-AC-15; no claim |
| R8-AC-9 — applicability review | §9.4 | C1-AC-15; explicit N/A |

## 13. Later-stage dependencies

- **C2:** public `Loaded`, `ServiceCapabilities`, and `ServiceCapabilitiesKnown`; status projection. C1 supplies only the private snapshot/validity.
- **C3:** ordinary load/ensure-load behavior. C1 does not add load semantics.
- **C4:** full exit/intent/stop transitions. C1 adds the generation invalidation needed to prevent stale authorization; C4 refines intent and terminal classification.
- **C5:** restart budget/backoff. C1 documents the map but does not implement timers or restart accounting.
- **C6:** identity/replacement. C1 increments generation on new process/replacement begin; C6 owns the replacement sequence.
- **C7:** production adapter and explicit load reachability. C1 adds no production wiring.

## 14. Authorization boundary

C1 remains unauthorized. This revision is design and documentation only. It does not modify source, implement C1, run tests or builds, authorize C2–C8, benchmark, deploy, or change A2-D9, the contract, the adopted r9 baseline, the frozen plan, or the owner decisions. Independent review is required before a stage-authorization request.
