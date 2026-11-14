# C1 Readiness Review — Revision 2

**Date:** 2026-11-08  
**Status:** UPDATED READINESS REVIEW — C1 UNAUTHORIZED; READY FOR INDEPENDENT REVIEW  
**Supersedes:** `docs/review/C1-READINESS-howl-tts-stage3-2026-11-08.md` (historical review; do not edit)  
**Scope:** Review of the revised C1 design and acceptance criteria. Design/documentation only; no source change, implementation, build, test, runtime execution, benchmark, or deployment.

## 1. Artifacts

| Artifact | SHA-256 | Dimensions |
|---|---|---|
| Revised C1 design | `003827f12658d1e556657220754795bff2094b39dca960321a71c3fb302cec6d` | 313 / 22,680 |
| Original C1 design / plan | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151` | 711 / 69,168 |
| Prior C1 readiness review | `e8b077485eea21af5312417d55e33a2d00803d7b56b6013b4929e2c9998da9eb` | 158 / 17,126 |
| Owner decisions, post-D10 | `677f884926de7e58c3dc7ab33f01340f1fbb5966c1b5c089b1aafd057bcf0fba` | 612 / 49,191 |
| D10 finalization note | `84fa78dd214f16f40ab125b544727fce2239e66bb1bb54dffc3c8777df9beb02` | 22 / 1,887 |
| Adopted r9 requirements | `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55` | 769 / 112,270 |
| r9 adoption record (frozen) | `22e79a3173cc6b9b4358d80eaa82029189e0fb0663a5ba4504af562153c7f10d` | 121 / 15,398 |
| Governing helper contract | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 / 182,904 |

The revised design supersedes only the C1 portions of the original plan. The original plan remains historical and unmodified. The owner decision, r9 adoption record, adoption finalization note, and governing contract are unchanged.

## 2. Verdict

**READY — for independent review, not for implementation.**

The revised design closes all seven prior readiness gaps at the design and acceptance level:

| Prior finding | Closure |
|---|---|
| C1-1 — authoritative snapshot validity | Closed by §3: supervisor-owned snapshot, generation, invalidation rules, and private validity predicate. |
| C1-2 — stale pre-D10 C1 text | Closed by §4 and §10: direct unknown results now assert A10a-2/A10b-2, including no-frame/no-spawn, no false success, bounded detail, and no new public taxonomy. |
| C1-3 — supervisor path calls exported `Cancel` | Closed by §5: separate `cancelForSupervisor` operation; watcher must not call exported `Cancel`. |
| C1-4 — capability policy/classification | Closed by §6: required/optional floor, scope validation, and `protocol_incompatible` semantic classification. |
| C1-5 — private C1 versus public C2 status | Closed by §7 and §10: C1 adds only private snapshot/validity; C2 owns public status fields and projection; C1-AC-16 excludes a new public taxonomy. |
| C1-6 — missing start/restart map | Closed by §8: state-scoped map with C1/C5/C7 ownership boundaries. |
| C1-7 — missing C1 evidence obligations | Closed by §9 and C1-AC-15: resource bounds, I/O impact, performance statement, and TH-1/HW-1/HW-2 applicability. |

No further owner decision is identified by this review as a C1 design prerequisite. The six other Stage 3 decisions remain pending and unselected. C1 remains unauthorized; this verdict is readiness for independent review, not authorization, implementation, or stage acceptance.

## 3. Requirements coverage

| Requirement | Revised design location | Acceptance / evidence |
|---|---|---|
| A2-D1 — mandatory `load`/`synthesize`; optional `catalog`/`unload`/`cancel` | §6.1 | C1-AC-1; missing-mandatory tests |
| A2-D2 — wire/application and service/provider scope | §6.2 | C1-AC-2; scope and duplicate tests |
| A2-D3 — semantic `hello` → `protocol_incompatible` | §6.3 | C1-AC-3; negative handshake matrix |
| A2-D4 — incompatible terminal/no restart | §3.2, §6.3 | C1-AC-3/C1-AC-4; terminate/Wait evidence |
| A2-D5/A2-D6 — `Loaded`, service status, known flag | §3, §7 | C1-AC-4; C2 boundary; no public status in C1 |
| A2-D7 — known-absent direct `unload` | §4.1, §6.1 | C1-AC-7 |
| A2-D8 — idle close uses `shutdown` | §8 | C1-AC-13 |
| A2-D9 — supervisor cancellation without `cancel` | §5 | C1-AC-10/C1-AC-11 |
| A2-D10 — no silent downgrade/no-op | §4, §8 | C1-AC-5/C1-AC-7/C1-AC-9 |
| D10(a) — A10a-2 | §4.1 | C1-AC-5/C1-AC-9 |
| D10(b) — A10b-2 | §4.2 | C1-AC-6/C1-AC-9 |
| r9 C1-CAPABILITY | §3, §9.1 | C1-AC-4/C1-AC-14 |
| r9 C1-START-MAP | §8 | C1-AC-12/C1-AC-13 |
| AC-1/AC-2/AC-3 — concurrency, locks, race evidence | §9.1 | C1-AC-14; race/overlap tests |
| AC-4/R8-AC-4 — resource budget | §9.2 | C1-AC-15; budget record |
| AC-5 — safety under load | §4, §9.2 | C1-AC-12/C1-AC-14 |
| AC-6 — no hidden background work | §9.2 | C1-AC-15; ownership review |
| AC-7/R8-AC-7 — I/O evidence | §9.3 | C1-AC-15 |
| R8-AC-8 — performance preservation | §9.4 | C1-AC-15; no performance claim |
| R8-AC-9 — applicability review | §9.4 | C1-AC-15; explicit N/A entries |

## 4. Acceptance case coverage

The revised design includes all required categories:

- **Positive:** optional capability present; known-present direct operations; service/provider scope with shared `synthesize`; valid handshake.
- **Negative:** missing mandatory; unknown/duplicate/wrong-version capabilities; scope violations; malformed descriptors; invalid cancel target; known-absent operations.
- **Stale:** after exit, stop, crash, incompatibility, backoff, and replacement; no frame, no start/spawn, no false `unloaded`/cancelled claim.
- **Concurrency:** handshake publication versus direct gate; process exit/transport close versus direct send; generation change; no torn snapshot; no double start.
- **Boundary:** direct `Cancel` versus supervisor path on the same helper; public status deferred to C2; no production import.

## 5. Explicit C1/C2 boundary

C1 owns only:

- the private supervisor-owned service snapshot;
- the generation-based validity predicate;
- private invalidation at process boundaries;
- direct-operation gating and the no-start current-instance send;
- the path-scoped supervisor cancellation operation.

C2 owns:

- `Loaded bool`;
- `ServiceCapabilities []TTSCapability`;
- `ServiceCapabilitiesKnown bool`;
- public `TTSProviderStatus` projection and its single-source matrix.

C1 acceptance must not claim any C2 field. C2 must not claim C1’s lifecycle invalidation or direct-operation behavior. C4–C6 own the later full exit/backoff/replacement transitions; C1’s generation invalidation is the prerequisite that prevents stale authorization.

## 6. Dependencies on later stages

- **C2:** public status projection of the private snapshot/known state.
- **C3:** ordinary load and ensure-load behavior.
- **C4:** full exit intent, shutdown, and crash classification; C1’s generation invalidation is already required.
- **C5:** restart budget/backoff and executed restart evidence.
- **C6:** identity replacement; C1 increments generation at replacement begin.
- **C7:** production adapter and explicit load reachability.

None is implemented by this revision.

## 7. Independent review checklist

1. Confirm the supervisor-owned snapshot is the only source for direct service-capability gating.
2. Confirm no direct `Unload`/`Cancel` path calls `Start` or `Supervisor.Start`.
3. Confirm stale snapshots after handshake, exit, stop, crash, incompatibility, backoff, and replacement are treated as unknown.
4. Confirm D10(a) and D10(b) results are asserted exactly as A10a-2/A10b-2, including target-first validation and the existing `AudioError` projection.
5. Confirm the supervisor cancellation path uses a separate internal operation and A2-D9 is unchanged.
6. Confirm A2-D1/A2-D2/A2-D3 validation and `protocol_incompatible` classification.
7. Confirm C1 adds no public status fields and no production import.
8. Confirm the start/restart map and resource/I/O/performance/applicability evidence are in the C1 report.
9. Confirm no source, test, build, runtime, benchmark, or deployment work was performed by this design pass.

## 8. Unresolved technical questions

This review identifies no unresolved owner decision as a C1 design prerequisite; the six other Stage 3 decisions remain pending and unselected. The following implementation choices remain for engineering and independent review; none changes the required behavior:

- exact naming and placement of the generation field;
- whether the private snapshot read is a small accessor or a supervisor method;
- exact invalidation hook placement in the existing `Terminate`/`observeExit` paths;
- whether the `snapshot=unknown` detail is retained only on `Failure` or also in internal diagnostics.

## 9. Authorization boundary

C1 remains unauthorized. This report does not implement C1, authorize C2–C8, modify source, run tests or builds, benchmark, deploy, or alter any frozen artifact. The next step is independent review of the revised design, not implementation.
