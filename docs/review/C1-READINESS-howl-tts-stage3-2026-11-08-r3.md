# C1 Readiness Review — Revision 3

**Date:** 2026-11-08  
**Status:** UPDATED READINESS REVIEW — C1 UNAUTHORIZED; READY FOR INDEPENDENT REVIEW  
**Supersedes:** `docs/review/C1-READINESS-howl-tts-stage3-2026-11-08-r2.md` (historical; do not edit)  
**Scope:** Reconciliation of Muse findings C1-R-1 through C1-R-9 and review of the corrected C1 design. Design/documentation only; no source, implementation, build, test, runtime execution, benchmark, or deployment.

## 1. Artifacts

| Artifact | SHA-256 | Dimensions |
|---|---|---|
| Corrected C1 design | `e9b30c6d9072375750851d13a05c2c53bb5cb983ec69498418a404d0291c754e` | 455 / 46,369 |
| Prior C1 design revision | `003827f12658d1e556657220754795bff2094b39dca960321a71c3fb302cec6d` | 313 / 22,680 |
| Prior C1 readiness report | `45457897844a5fe90c247bf97225420df0700ea431d2a7ac65cf14617867692d` | 129 / 8,910 |
| Original C1 design / plan | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151` | 711 / 69,168 |
| Owner decisions, post-D10 | `677f884926de7e58c3dc7ab33f01340f1fbb5966c1b5c089b1aafd057bcf0fba` | 612 / 49,191 |
| D10 finalization note | `84fa78dd214f16f40ab125b544727fce2239e66bb1bb54dffc3c8777df9beb02` | 22 / 1,887 |
| Adopted r9 requirements | `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55` | 769 / 112,270 |
| r9 adoption record | `22e79a3173cc6b9b4358d80eaa82029189e0fb0663a5ba4504af562153c7f10d` | 121 / 15,398 |
| r9 adoption finalization | `cf1fdfc3eebbcfd3b20ca1ac0b746ae19f3d0716a211bba79924d8beca6b488c` | 21 / 1,348 |
| Governing helper contract | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 / 182,904 |

The corrected design supersedes only the proposed C1 design/acceptance portions of the original plan and the prior revision. It does not modify the frozen plan, contract, r9 adoption, owner decisions, or prior review artifacts.

## 2. Verdict

**READY — for independent review, not for implementation.**

All nine Muse findings are substantiated and closed at the design/acceptance level. No finding is rejected or silently omitted. The revised design adds the missing state/health map, complete evidence-record requirements, live-versus-exited `unavailable` split, intent-aware incompatible termination, D10 caller projection, traceability correction, known-absent supervisor branch, concurrency coverage, and exact protocol details.

C1 remains unauthorized. The verdict is readiness for another independent review by Muse, not authorization, implementation, or stage acceptance.

## 3. Finding-by-finding reconciliation

| Finding | Disposition | Closure in corrected design | Independent evidence checked |
|---|---|---|---|
| **C1-R-1 — Start/restart map** | **Accepted** | §8.1 covers every state and includes `crashed → backoff → starting`, fourth-attempt exhaustion → `unavailable`, exit-zero/nonzero/grace classification, and the A2-D4 exception. §8.2 covers every §8.1 health edge. C1-AC-13 and tests enforce completeness. | r9 §8 C1 and §8.1; contract A2-D4; §17.11.5/§17.7.2. |
| **C1-R-2 — Resource and I/O evidence** | **Accepted** | §9.2 has all seven §10 fields for every resource row. §9.3 has all required fields for all eight §3.2 I/O categories, including reviewed N/A. §9.4 has R8-AC-8 and R8-AC-9. C1-AC-15 requires a complete package and written independent determinations. | r9 §3.2, §7, §9.1, §9.2, §10; R8-AC-4/7/8/9. |
| **C1-R-3 — Asset-specific `unavailable`** | **Accepted** | §3.2 rule 7 and §3.4 distinguish live-retain from exited/no-instance-invalidate. §8 has separate rows; §11 has separate tests; C1-AC-18 requires both. | Contract A2-D5/A2-D6. |
| **C1-R-4 — Semantic handshake termination** | **Accepted** | §6.4 latches `terminationIntent=incompatible` before termination. Timeout and later exit preserve `incompatible`; no `crashed`, backoff, budget spend, or replacement. C1-AC-3 adds timeout/late-exit tests. | Contract A2-D4; current `Client.Start`/`Supervisor.observeExit` behavior. |
| **C1-R-5 — D10 caller-facing projection** | **Accepted** | §4.3 specifies `Failure.Code=unsupported_capability`, bounded detail, `AudioError.Code=ErrCodeSynthesisFailed`, and the same `Failure` in the cause. C1-AC-5/6/19 and tests require both layers. | Owner D10 A10a-2/A10b-2; `errors.go` mapping. |
| **C1-R-6 — Traceability corrections** | **Accepted** | §12 separates original AC-7 (C8-only) from R8-AC-7 (I/O), removes the A2-D8 coverage claim, and records A2-D8 as a non-substitution boundary. §4.1 says `Unload` is not shutdown. | r9 §7.1/§7.3; contract A2-D8. |
| **C1-R-7 — Supervisor known-absent capability** | **Accepted** | §5 separates unknown and known-absent; C1-AC-20 and a named test require the known-absent branch; both skip the frame and synthesize no capability error. | Contract A2-D9 and D10(b) boundary. |
| **C1-R-8 — Concurrency coverage** | **Accepted** | §9.1 documents existing serialization, `workMu`/`startMu`, lock discipline, and winner rules. C1-AC-21 and §11 cover concurrent unload/cancel, startup/loading/stopping/recovery, supervisor cancellation during startup/shutdown/recovery, and exit/replacement overlaps. | Original AC-1/2/3; r9 §8 C1 TH-1 focus; current `client.go`/`supervisor.go`. |
| **C1-R-9 — Exact protocol and evidence details** | **Accepted** | §6.1 states exact `"missing"` key, table order, and omission. §9.4 states PR impact and reviewed N/A. C1-AC-22 requires static review of C1-AC-16. | Contract A2-D3; R8-AC-8/9; C1-AC-16. |

No disposition is `Rejected`; no correction is left for a later design pass.

## 4. Requirements coverage

| Requirement | Corrected design | Acceptance / evidence |
|---|---|---|
| A2-D1 — mandatory/optional floor | §6.1 | C1-AC-1; exact missing tests |
| A2-D2 — service/provider scope | §6.2 | C1-AC-2 |
| A2-D3 — semantic classification and `"missing"` | §6.1, §6.3 | C1-AC-1/3 |
| A2-D4 — incompatible terminal/no restart | §6.4 | C1-AC-3; timeout/late-exit |
| A2-D5/A2-D6 — status and known snapshot | §3, §7 | C1-AC-4/18 |
| A2-D7 — known-absent unload | §4.1 | C1-AC-7 |
| A2-D8 — idle shutdown, not unload | §4.1 boundary only | No C1 implementation/coverage claim; C4/C7 own shutdown; C1-AC-22 |
| A2-D9 — supervisor cancellation | §5 | C1-AC-10/11/20 |
| A2-D10 — no silent no-op | §4, §8 | C1-AC-5/7/9 |
| D10(a)/D10(b) — A10a-2/A10b-2 | §4.1–§4.3 | C1-AC-5/6/19 |
| r9 C1-CAPABILITY | §3, §9.1 | C1-AC-4/14 |
| r9 C1-START-MAP | §8 | C1-AC-13 |
| Original AC-1/2/3 | §9.1 | C1-AC-14/21 |
| Original AC-4 / R8-AC-4 | §9.2 | C1-AC-15 |
| Original AC-5 | §4, §9.1 | C1-AC-12/14/21 |
| Original AC-6 | §9.1, §9.2 | C1-AC-15/17 |
| R8-AC-7 — I/O | §9.3 | C1-AC-15 |
| R8-AC-8 — performance | §9.4 | C1-AC-15; no claim |
| R8-AC-9 — applicability | §9.4 | C1-AC-15; reviewed N/A |
| **Original AC-7 — no stage substitution** | **No C1 design content; C8-only** | **Not a C1 claim** |
| R8-AC-10 — pointer | No C1 design content; C8-only | Not a C1 claim |

Bare `AC-7` is not used for I/O anywhere in revision 3. The I/O rule is R8-AC-7.

## 5. State, health, and evidence coverage

The complete map covers `absent`, `starting`, `handshaking`, `unloaded`, `loading`, `ready`, `degraded`, `busy`, `stopping`, `crashed`, `backoff`, `incompatible`, live asset-specific `unavailable`, exited/no-instance `unavailable`, and process `exited`.

- `crashed` uses only the next explicit TTS request, enters `backoff`, then `starting` within budget; a fourth attempt in ten minutes exhausts to `unavailable`.
- Status, catalog, health, and direct optional operations do not restart a crashed/backoff helper.
- Normal `stopping`: exit 0 → `absent`; nonzero or grace expiry → `crashed`.
- Incompatible termination is the A2-D4 exception: expected/late exit preserves `incompatible`, and there is no crash/backoff/budget/replacement.
- Health is unchanged in every applicable state; `handshaking` and `incompatible` are the only wire-health exclusions; live and exited `unavailable` are separately recorded.

Resource evidence includes every §10 field, including the independent conformance determination. I/O evidence includes all eight §3.2 categories, path/component, operation, trigger/cadence, necessity, bound, category, measurement/design proof, and reviewed N/A rationale. The performance statement names PR-1 and PR-4 dimensions, records PR-2/PR-3 as unaffected/N/A, and makes no measurement or claim. The applicability table does not use blanket N/A.

## 6. D10 and protocol exactness

- Unknown direct `Unload`: `unsupported_capability`, `capability=howl.tts.unload`, `snapshot=unknown`; no frame/start/spawn/restart; no `unloaded` claim.
- Unknown direct `Cancel`: target validated first; `unsupported_capability`, `capability=howl.tts.cancel`, `snapshot=unknown`; no frame/start/spawn/restart; no cancellation claim.
- Both project through the existing `AudioError` mapping with `ErrCodeSynthesisFailed`, preserving the precise `Failure` and bounded detail in the cause.
- No new public error code, taxonomy, wire message, or Wails binding is added.
- `detail.missing` uses the exact key `"missing"`; order is `howl.tts.load,howl.tts.synthesize`; no spaces or duplicates; the key is omitted when none are missing.
- A2-D9 remains separate: unknown and known-absent supervisor cancellation both skip the frame and synthesize no `unsupported_capability`, but are separately implemented and tested.

## 7. Concurrency boundaries and winner rules

C1 reuses existing synchronization rather than adding a new architecture:

- `Client.workMu` serializes `Load`, `Unload`, `Catalog`, and the whole `Synthesize` operation. Concurrent `Unload` calls serialize; an `Unload` issued during loading/busy waits and re-gates.
- `Client.startMu` and `Supervisor.startMu` preserve one start owner. Direct optional operations never enter either path.
- `Client.mu` guards transport/pending state; `Supervisor.mu` guards snapshot publication/invalidation and termination intent.
- `Process.writeMu` prevents interleaved frame bytes. Direct `Cancel` remains outside `workMu`, preserving cancellation responsiveness.
- No lock is held across `Start`, `Send`, `Terminate`, or `Wait`; `workMu` may span a request wait because it is the request-serialization boundary.
- Direct optional operations cannot start, spawn, or restart a helper. A post-gate exit yields a typed failure, not a replacement or false state.
- Concurrent cancels retain independent request IDs; a late acknowledgement cannot overwrite the target's terminal.

## 8. Residual conditions and limitations

This is a design package, not C1 evidence. The following are explicitly not claimed as complete:

1. **No implementation or test execution.** All acceptance criteria, tests, measurements, failure logs, and conformance determinations are prospective.
2. **Resource conformance is not self-certified.** The C1 report must obtain the independent budget reviewer's written determination per row; the author cannot fill it.
3. **The `MaxPending=32` row is an inherited implementation cap, not a new C1 limit.** The design requires the reviewer to accept an authorized source or record an unresolved R8-AC-4 item. C1-AC-15 cannot close while that item remains unresolved.
4. **I/O N/A entries require named reviewer/date.** The design supplies the required categories and rationale; the C1 report must supply the reviewed entries.
5. **Executed restart/backoff and start/reachability are C5/C7 evidence.** C1 owns the map and structural checks, not their execution.
6. **Independent review is still required.** Muse has not reviewed revision 3; this report is the author's reconciliation, not the independent verdict.

These conditions do not reopen a design finding. They identify the evidence that a C1 implementation and its independent reviewers must produce.

## 9. Independent review checklist

1. Confirm all nine findings are substantively closed, not merely renamed.
2. Confirm the map includes `crashed → backoff → starting`, exhaustion, exit classification, the A2-D4 exception, and every §8.1 health edge.
3. Confirm live asset-specific `unavailable` retains and exited/no-instance `unavailable` invalidates, with separate tests.
4. Confirm timeout and late exit after semantic incompatibility cannot produce `crashed`, backoff, budget spend, or replacement.
5. Confirm A10a-2/A10b-2 require both the helper `Failure` and `AudioError.Code=ErrCodeSynthesisFailed`, with preserved detail.
6. Confirm original AC-7 is C8-only and R8-AC-7 is the only I/O rule; confirm A2-D8 is not claimed as C1 coverage.
7. Confirm unknown and known-absent supervisor cancellation are separate branches and tests.
8. Confirm concurrent unload/cancel, lifecycle-state overlaps, and supervisor cancellation during startup/shutdown/recovery are covered.
9. Confirm the exact `"missing"` key/order/omission, R8-AC-8 dimension statement, R8-AC-9 reviewed N/A, and C1-AC-16 static assertion.
10. Confirm no source/test/build/runtime/benchmark/deployment was performed by this design pass.

## 10. Authorization boundary

C1 remains unauthorized. This report does not implement C1, authorize C2–C8, modify source, run tests/builds/runtime/benchmarks, deploy, amend the contract, amend adopted r9, amend the frozen plan, or alter any owner decision. The next step is independent review by Muse of revision 3, not implementation.
