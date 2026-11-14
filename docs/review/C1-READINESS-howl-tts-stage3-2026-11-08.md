# C1 Readiness Review — Howl TTS Stage 3

**Date:** 2026-11-08  
**Status:** READINESS REVIEW ONLY — C1 UNAUTHORIZED; NO SOURCE CHANGES  
**Scope:** Adopted r9, governing helper contract, Stage 3 corrective design/plan, the completed D10 decisions, and the current helper source. No build, test, replay, benchmark, runtime execution, or deployment was performed.

## 1. Artifact identities

| Artifact | SHA-256 | Dimensions |
|---|---|---|
| Adopted r9 requirements | `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55` | 769 / 112,270 |
| r9 adoption record (frozen) | `22e79a3173cc6b9b4358d80eaa82029189e0fb0663a5ba4504af562153c7f10d` | 121 / 15,398 |
| Adoption finalization note | `cf1fdfc3eebbcfd3b20ca1ac0b746ae19f3d0716a211bba79924d8beca6b488c` | 21 / 1,348 |
| Governing helper contract | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 / 182,904 |
| Owner-decision instrument, post-D10 | `677f884926de7e58c3dc7ab33f01340f1fbb5966c1b5c089b1aafd057bcf0fba` | 612 / 49,191 |
| Owner-decision adoption-time identity | `24d3344554e7ed6a465c8b7bdbca4461cda23f78353aa7626e3cf4ba7df122ba` | 610 / 45,891 |
| C1 design / corrective plan | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151` | 711 / 69,168 |

The owner-decision instrument was completed in this assignment; the C1 plan was reviewed unchanged. The r9 adoption record and its finalization note were not altered.

## 2. Readiness verdict

**READY WITH CONDITIONS — not authorization-ready as written.**

The C1 design direction is sound and the D10 owner gates are now recorded, but C1 should not be authorized until the conditions in §5 are incorporated into the C1 plan and its acceptance tests. No new owner decision is required for those corrections; they are engineering-level design and test corrections.

C1 remains unauthorized. This review does not implement or accept C1, and it does not authorize C2 or any later stage.

## 3. Requirements-to-implementation coverage

| Requirement | C1 design coverage | Current source / evidence | Determination |
|---|---|---|---|
| A2-D1 — mandatory `load`/`synthesize`; optional `catalog`/`unload`/`cancel` | C1 plan says replace the five-capability mandatory set and test independent optional absence | `protocol.go:53-56` still lists all five as mandatory | **Gap C1-4**; C1 must change |
| A2-D2 — service versus provider capability scope | C1 plan says scope validation and separate storage | `validateCapabilities` does not enforce service-only versus provider-only scope | **Gap C1-4**; C1 must change |
| A2-D3 — semantic `hello` failures are `protocol_incompatible` | C1 plan says distinguish semantic failures from frame violations and use `detail.missing` | `protocol.go:527-531` returns `protocol_violation` for missing mandatory capabilities | **Gap C1-4**; C1 must change |
| A2-D4 — `incompatible` terminal, terminate/observe, no restart | C1 plan says terminate/observe and set `incompatible` | `Start` sets `incompatible` only for `protocol_incompatible`/`protocol_minor_unsupported`; semantic failures currently produce `protocol_violation` and `unavailable` | **Gap C1-4**; follows the A2-D3 correction |
| A2-D5/A2-D6 — `Loaded`, service capability status, `ServiceCapabilitiesKnown` | C2 adds the public fields; C1 plan says no status expansion | `types/provider.go` `TTSProviderStatus` lacks `Loaded`, `ServiceCapabilities`, and `ServiceCapabilitiesKnown`; `Supervisor.Status` cannot express unknown | **Gap C1-5 / C2-owned**; C1 may use a private predicate, but must not claim the public status surface |
| A2-D7 — known-absent direct `unload` | C1 plan says `unsupported_capability`, no frame | Current method returns `unsupported_capability` when the snapshot is empty; however, the all-five mandatory set prevents a conforming helper from reaching known-absent | **Gap C1-4**; behavior exists but is not reachable until the floor is corrected |
| A2-D8 — idle close uses `shutdown`, not `unload` | No D10 interaction | Current `Shutdown` uses `shutdown`; no C1 change required | **Pass/unchanged** |
| A2-D9 — supervisor cancellation without `cancel` | C1 plan requires a path-scoped internal operation and says no `unsupported_capability` | `client.go:866-873` calls exported `Cancel`; the direct method synthesizes the error and the watcher discards it | **Gap C1-3**; C1 must separate the path |
| A2-D10 — no silent downgrade or no-op | C1 plan says absent is never supported; direct operations fail closed | Unknown pre-handshake fails closed; a stale `c.hello` after exit can pass the gate | **Gap C1-1**; stale snapshot can defeat the rule |
| D10(a) A10a-2 — unknown direct `unload` | Plan still says C1 freezes no result; must now assert `unsupported_capability` with no frame/no spawn | A never-started client reaches `unsupported_capability`; stale-after-exit does not reliably remain unknown | **Gap C1-1/C1-2** |
| D10(b) A10b-2 — unknown direct `Cancel` | Plan still says C1 freezes no result; must now assert `unsupported_capability` after target validation, no frame/no spawn | A never-started client reaches `unsupported_capability`; target validation exists; stale-after-exit does not reliably remain unknown | **Gap C1-1/C1-2** |
| r9 C1-CAPABILITY — atomic capability publication, single writer, no torn reads | Plan requires lock-protected snapshot and focused race tests | `requireServiceCapability` reads `c.hello` under `c.mu`; no explicit publication/validity generation | **Partial**; C1 must define the predicate and lock discipline |
| r9 C1-START-MAP — documented state-scoped start/restart map | Plan has no dedicated C1 start/restart map; the r9 overlay allocates it to C1 | Current `Start` has a crash retry path; no documented map at C1 acceptance | **Gap C1-6** |
| TH-1/AC-1–AC-3 — concurrency model, lock ownership, race/overlap evidence | C1 plan has race/concurrency tests and lock-protected snapshot reads | No executable evidence; source is not a committed baseline | **Condition**; C1 report must supply |
| AC-4/R8-AC-4 — bounded resource budget | C1 plan does not state a C1 budget record | No C1 budget record | **Gap C1-7** |
| AC-5 — safety under load | No model load or writes from negotiation; direct operations are control-only | Source supports the intent, but stale snapshot can start a helper | **Condition**; test with overlap/failure injection |
| AC-6 — no hidden background work | C1 plan controls starts; watcher remains bounded | Current watcher is a goroutine per synthesis; no new C1 background loop is planned | **Pass with existing lifecycle caveat** |
| R8-AC-7/AC-8/AC-9 — I/O, performance, applicability records | C1 plan has no explicit C1-stage records | No C1 I/O, performance, or N/A record | **Gap C1-7** |
| No production caller / no model load / no write | C1 plan says no production import; helper is internal | Found callers only in tests and `watchCancellation`; no production adapter wiring | **Pass**; preserve as a C1 report assertion |

## 4. Confirmed gaps and required design corrections

### C1-1 — Authoritative snapshot validity is not implemented

This is the primary blocking gap.

- `Client.Start` stores `c.hello` only after a successful handshake (`client.go:262`).
- `Start` clears `started` and `transport` after a failed handshake, but does not clear an older `c.hello`.
- `Supervisor.observeExit` (`supervisor.go:166-197`) does not clear `s.capabilities`.
- `TTSProviderStatus` has no `ServiceCapabilitiesKnown` field.
- Consequently, after a successful handshake and later exit, a direct `Unload` or `Cancel` can pass the capability gate. The known-present path calls `c.call`, which calls `Start`; that can start a replacement helper. This violates the D10 no-spawn rule.

Required correction: define an authoritative `serviceCapabilitiesKnown` predicate in C1, consistent with A2-D6. It is true only for the current helper instance after a compatible handshake and while the instance has not entered `stopping`, `exited`, `crashed`, `incompatible`, or `backoff`, and no replacement has begun. An asset-specific `unavailable` on a live, handshaken instance retains the snapshot; an exited or no-instance `unavailable` does not.

C1 must gate direct `Unload` and `Cancel` on this predicate before `c.call`, and must clear or ignore a stale `c.hello` after failed handshake, exit, stop, incompatibility, or replacement. It may use `Supervisor.ProcessAlive()` plus the lifecycle state as the conservative C1 predicate; C4 later owns the transition writers. C1 must not start a helper to discover the snapshot.

### C1-2 — C1 plan text is stale after the D10 decisions

The C1 plan still says the unknown-snapshot result is unresolved, that C1 tests only the structural boundary, and that C1 must not assert a result. That contradicts the recorded D10(a) A10a-2 and D10(b) A10b-2 decisions.

Required correction: update the C1 design, tests, and acceptance text to assert:

- direct `Unload` unknown → `unsupported_capability`; no frame; no start/spawn/restart; no `unloaded`;
- direct `Cancel` unknown, after target validation → `unsupported_capability`; no frame; no start/spawn/restart; no cancellation claim;
- bounded detail `capability=howl.tts.unload|cancel` and `snapshot=unknown`.

### C1-3 — The supervisor cancellation path still uses the direct method

`client.go:866-873` calls exported `Cancel`. A2-D9 requires a path-scoped internal operation that records `cancel_requested`, skips the frame when `cancel` is absent or unknown, synthesizes no `unsupported_capability`, and follows grace/`KillTree`/`Wait()` and cancellation precedence. The current discarded error does not satisfy the “no synthesis in the path” rule.

Required correction: add the internal path-scoped operation and change `watchCancellation` to call it. Keep exported `Cancel` as the direct operation with the D10 result.

### C1-4 — Capability floor, scope validation, and semantic classification

The current source is still pre-Policy-B in three respects:

- all five service capabilities are mandatory (`protocol.go:53-56`);
- `validateCapabilities` does not enforce service-only versus provider-only scope;
- semantic `hello` failures are `protocol_violation` (`protocol.go:527-531`), not `protocol_incompatible`.

Required correction: implement A2-D1/A2-D2/A2-D3, including `detail.missing` containing only required IDs, duplicate detection, invalid descriptors, service/provider scope, and terminal `incompatible` behavior with terminate/observe and no restart.

### C1-5 — Public status cannot represent unknown, but C2 owns that surface

`TTSProviderStatus` lacks `Loaded`, `ServiceCapabilities`, and `ServiceCapabilitiesKnown`; `Supervisor.Status` reports only provider capabilities. A2-D6 requires the service-scoped fields and the known flag.

C1 should use a private authoritative predicate and must not claim the public status surface. C2 owns the public projection. C1 acceptance should state this boundary explicitly so a private predicate is not mistaken for the C2 status contract.

### C1-6 — The C1 start/restart map is missing

r9 allocates `C1-START-MAP` to C1: a documented state-scoped trigger/start/restart map, with no optional-operation spawn. The C1 plan has no dedicated map and no acceptance for it.

Required correction: add the map to C1, covering `absent`, `crashed`, `backoff`, `stopping`, `incompatible`, and `unavailable`; no direct `unload`/`Cancel` spawn; no status/catalog restart from `crashed`/`backoff`; no start from `stopping` or `incompatible`; no recovery from `unavailable` except the explicit same-identity `load` path under D1a/D1b.

### C1-7 — C1 evidence obligations are not yet planned

r9 requires C1 to produce stage-level evidence for AC-1 through AC-7 and the applicable R8-AC overlays. The C1 plan has race tests but no C1 resource-budget record, I/O impact record, performance-preservation statement, or applicability determination. These are documentation/evidence additions, not owner decisions.

Required correction: add a C1 evidence plan with a numeric or contract-provisional budget, explicit I/O impact, a no-performance-claim statement, and TH-1/HW-1 applicability entries rather than blanket N/A.

### C1-8 — Minor precision items

- Preserve target-ID validation before the direct `Cancel` snapshot gate; malformed targets remain `invalid_request`.
- Keep the direct and supervisor cancellation branches physically separate; do not let the internal path call exported `Cancel`.
- State that the D10 detail is internal/bounded and does not create a wire message or public error code.
- Assert no model load and no disk write from negotiation or inspection.

## 5. Contradictions and missing requirements

| Item | Finding |
|---|---|
| C1 plan vs D10(a)/D10(b) | **Contradiction.** The plan says the unknown results are unresolved and tests must not assert a code; the owner-decision record now fixes A10a-2/A10b-2. Correct the plan before authorization. |
| Current mandatory set vs A2-D1 | **Contradiction.** All five capabilities are mandatory in `validateHello`; A2-D1 makes three optional. |
| Current semantic classification vs A2-D3 | **Contradiction.** Semantic `hello` failures return `protocol_violation`; A2-D3 requires `protocol_incompatible`. |
| Current watcher vs A2-D9 | **Contradiction.** The watcher calls exported `Cancel`, which can synthesize the direct-operation error; A2-D9 forbids that on the supervisor path. |
| Stale snapshot vs A2-D6/D10 | **Missing requirement in the plan.** The plan does not define snapshot invalidation or test stale-after-exit. |
| C1 start/restart map vs r9 C1 | **Missing deliverable.** The map is allocated to C1 but absent from the C1 plan. |
| Public status vs A2-D6 | **Deferred to C2, not a C1 defect.** C1 must not claim it; C1 needs only the private predicate. |
| AC-4/R8-AC evidence | **Missing from the plan.** C1 acceptance cannot be complete without the required stage records. |

## 6. Required test coverage

1. **Unknown, pre-handshake:** direct `Unload` and direct `Cancel` return `unsupported_capability`, with `snapshot=unknown`; zero operation frames; zero starts/spawns; no `unloaded` or cancellation claim.
2. **Known absent:** optional capability omitted from a valid handshake; direct operation returns `unsupported_capability`, zero operation frames, no additional spawn.
3. **Known present:** optional capability advertised; exactly one operation frame; normal response; no additional spawn.
4. **Stale after successful handshake:** helper exits; direct operations must treat the snapshot as unknown, return `unsupported_capability`, send no frame, and start no replacement. Cover exit, stop, crash, incompatibility, backoff, and replacement where constructible.
5. **Target validation:** malformed `Cancel` target returns `invalid_request` before the snapshot gate, with no frame/start.
6. **Semantic handshake failures:** missing mandatory, unknown capability, wrong version, duplicate, invalid descriptor, and scope violations all return `protocol_incompatible`; process terminated and `Wait()` observed; state `incompatible`; no restart.
7. **Optional combinations:** each optional capability independently present/absent; `detail.missing` contains only missing mandatory IDs, in order, with no spaces/duplicates.
8. **Supervisor path:** absent or unknown `cancel` produces no frame and no `unsupported_capability`; grace/`KillTree`/`Wait()` and precedence remain authoritative.
9. **Two-sided regression:** same helper, supervisor cancellation path and direct `Cancel` produce distinct behavior.
10. **Race/overlap:** handshake completion versus direct gate; snapshot read under lock; no torn fields; no frame after gate failure; no double start.
11. **Evidence assertions:** one start maximum for handshake; no model load; no disk write; no production import; process/goroutine ownership bounded.

## 7. Unresolved technical questions

These are implementation choices for engineering, not owner decisions:

- Whether C1 clears `c.hello` immediately on exit or uses a generation/validity token checked under lock. Either is acceptable if the D10 no-spawn and no-stale-authorization tests pass.
- The exact C1 predicate boundary versus C4’s transition ownership. C1 can conservatively require `ProcessAlive()` and a non-invalid state; C4 later writes the full transition set.
- Whether the `snapshot=unknown` detail is stored on `Failure` only or also surfaced through the internal diagnostic status. It must not become a wire message or new public error code.
- Whether C1 adds a small private snapshot-known method or derives the predicate in `requireServiceCapability`; the choice must preserve single-writer publication and lock ordering.

## 8. Authorization boundary

C1 remains unauthorized. This review recommends design corrections and test coverage only. It does not modify source, implement C1, run tests or builds, authorize C2–C8, or change A2-D9, the contract, the adopted r9 baseline, or the frozen r9 adoption record.
