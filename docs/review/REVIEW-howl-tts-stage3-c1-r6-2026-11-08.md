# Independent Review — Howl TTS Stage 3 C1 Revision 6

**Date:** 2026-11-08  
**Review type:** Independent, read-only design review  
**Reviewer:** Muse  
**Artifact set reviewed:** C1 design r6, readiness r6, integrity r6  
**Verdict:** **FAIL — DESIGN NOT READY; C1 REMAINS UNAUTHORIZED**

## 1. Artifact identity and independent hash verification

The following were computed directly from the workspace, not taken from the integrity report:

| Artifact | SHA-256 | Lines / bytes | Match |
|---|---|---|---|
| `docs/review/DESIGN-howl-tts-stage3-c1-revision-2026-11-08-r6.md` | `dae4b371e901c11144a80dd1e3b99239ed25d1341088546ce058447e232a6fb4` | 853 / 127,414 | Yes |
| `docs/review/C1-READINESS-howl-tts-stage3-2026-11-08-r6.md` | `d5d21bbd9055b51cba8aac58cd0ea5b46ed83ea3d634e29700310b6afb8823e7` | 141 / 13,300 | Yes |
| `docs/review/C1-INTEGRITY-howl-tts-stage3-2026-11-08-r6.md` | `6ce92cedb19937f187dc62ebb05194badce14a90278f39ec5d921078347375b4` | 78 / 6,437 | Yes |

Structural checks also pass:

- 23 distinct `C1-R` identifiers.
- 26 distinct `C1-AC` identifiers.
- 26 acceptance-criteria definitions and 26 §10.1 verification rows.
- 15 state/trigger rows and 14 health-state rows.
- No duplicate `### N.M` section numbers.

**Preserved inputs independently hashed:**

| Artifact | SHA-256 | Lines / bytes |
|---|---|---|
| Revision 5 design | `6469203d9b8ae869988cf862b293f464e3e788c47badbc964f2534a56e04c387` | 656 / 89,192 |
| Revision 5 readiness | `67ce40b5759aa3ea85b647148ff62cb40c71922048f02d68aa084ded33275179` | 146 / 13,514 |
| Revision 5 integrity | `b1829b4eaa31ac826a4bbc5da9703dd5e31d0b202dba2bb43ac230f0be3cbf06` | 73 / 5,736 |
| Revision 4 design | `b099afe0302672bf0459ae9ec1b493e64bd11c021c6568d2a6a5ab3e08c78e3f` | 539 / 65,849 |
| Revision 4 readiness | `e291967f56312bb07c177b2d832a1a69683a2410165def6d77b7d1c255ec7e5b` | 128 / 11,076 |
| Revision 3 design | `e9b30c6d9072375750851d13a05c2c53bb5cb983ec69498418a404d0291c754e` | 455 / 46,369 |
| Revision 3 readiness | `aca2d59921709c7e0239547d1b172572a790140f185389ffdc00a04e2b714bd7` | 138 / 13,476 |
| Prior C1 design | `003827f12658d1e556657220754795bff2094b39dca960321a71c3fb302cec6d` | 313 / 22,680 |
| Prior readiness r2 | `45457897844a5fe90c247bf97225420df0700ea431d2a7ac65cf14617867692d` | 129 / 8,910 |
| Original C1 plan | `802ae8114993b7889140f5f1c216eca042f74e84f7dbcb1665f05434e355b151` | 711 / 69,168 |
| Owner decisions | `677f884926de7e58c3dc7ab33f01340f1fbb5966c1b5c089b1aafd057bcf0fba` | 612 / 49,191 |
| D10 finalization | `84fa78dd214f16f40ab125b544727fce2239e66bb1bb54dffc3c8777df9beb02` | 22 / 1,887 |
| Adopted r9 requirements | `0f71cee5e0e5d25b58d1b3956ff549df8cc625ad7f9e8c82189f2f2ef9366c55` | 769 / 112,270 |
| r9 adoption record | `22e79a3173cc6b9b4358d80eaa82029189e0fb0663a5ba4504af562153c7f10d` | 121 / 15,398 |
| r9 adoption finalization | `cf1fdfc3eebbcfd3b20ca1ac0b746ae19f3d0716a211bba79924d8beca6b488c` | 21 / 1,348 |
| Governing contract | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | 2,586 / 182,904 |

All preserved identities match. No source or prior artifact was modified.

## 2. Reconciliation of C1-R-1 through C1-R-23

“Design-closed” below means the design text addresses the cited design defect. It does **not** mean acceptance-closed, independently verified, or evidence-satisfied.

| Finding | Review disposition | Basis / remaining work |
|---|---|---|
| C1-R-1 — Start/restart map | **design-closed for design; evidence pending** | §8 map covers the required states; C5/C7 execution remains pending. |
| C1-R-2 — Resource/I/O evidence | **design-closed; acceptance-evidence open** | Records exist, but independent determinations do not. C1-AC-15 remains **OPEN**. |
| C1-R-3 — Asset-specific `unavailable` | **partially closed** | Live-retain versus exited-invalidate is stated, but the unexpected-exit path still misclassifies to `crashed`; see C1-R-25. |
| C1-R-4 — Semantic handshake termination | **design-closed for `incompatible`; acceptance pending** | Intent-aware incompatible termination is described; no executed evidence. |
| C1-R-5 — D10 caller projection | **design-closed; tests pending** | Shape is explicit; no executed projection test. |
| C1-R-6 — Traceability | **design-closed** | Original AC-7 and R8-AC-7 are separated. |
| C1-R-7 — Supervisor known-absent | **design-closed; test pending** | Known-absent branch remains distinct. |
| C1-R-8 — Concurrency | **partially closed** | The matrix is broader, but the writer/event race model is incomplete; see C1-R-26/C1-R-27. |
| C1-R-9 — Protocol/evidence details | **design-closed** | Exact `missing` key/order/omission rule is present. |
| C1-R-10 — Readiness overstatement | **author-reconciled; not independent acceptance** | The r6 readiness correctly says the r6 review is pending; r6 is now reviewed, but evidence closure remains open. |
| C1-R-11 — Direct `Unload` clears `unavailable` | **design-closed; tests pending** | The no-aggregate-transition rule is explicit. |
| C1-R-12 — Incompatible start bypass | **design-closed; acceptance pending** | Both start gates are specified; implementation/evidence absent. |
| C1-R-13 — Process-cap citation | **design-closed** | The §5.5 citation is correct. |
| C1-R-14 — AC/HW-1.6 precision | **design-closed; review pending** | 26 criteria/rows exist; HW-1.6 assumptions still require independent review. |
| C1-R-15 — Health/Catalog contradiction | **design-closed; tests pending** | Live wire health and advertised catalog reads are specified; no execution. |
| C1-R-16 — Exact error semantics | **partially closed** | `provider_not_loaded` scope is fixed, but the no-frame `request_not_found` cancellation result lacks clear authority; see C1-R-31. |
| C1-R-17 — `MarkHandshake` bypass | **partially closed** | Supervisor publication is process-owned, but client-side hello/identity copies remain outside that ownership; see C1-R-26. |
| C1-R-18 — Aggregate `incompatible` during termination | **partially closed** | Incompatible intent is described, but derived `ProcessAlive` and the writer-test contradiction undermine the assertion; see C1-R-24/C1-R-29. |
| C1-R-19 — Independence wording | **closed for wording** | The report no longer claims the r6 author reconciliation is independent. This review is the separate r6 review. |
| C1-R-20 — Aggregate-state mutation surface | **partially closed** | The inventory is incomplete for `activeRequestID`, derived status predicates, and some terminal synthesis outcomes; see C1-R-24/C1-R-27. |
| C1-R-21 — Handshake provenance | **partially closed** | The supervisor attempt is stronger, but `Client.hello`/provider identity and the old `helloHandshake` path remain; see C1-R-26. |
| C1-R-22 — Catalog start boundary | **partially closed** | The table distinguishes direct and C7 paths, but the existing `ProviderCatalog` still delegates to direct `Client.Catalog`, and catalog-driven selection is removed without a replacement; see C1-R-28. |
| C1-R-23 — Complete projection matrix | **partially closed** | The matrix is expanded, but its “six classes × seven operations” count is inconsistent with the direct/adapter Catalog rows, and the no-frame cancel origin is unresolved; see C1-R-30/C1-R-31. |

No prior finding is **acceptance-closed**. No C1-R-1 through C1-R-23 item is fully independent-closed for r6 except the wording item; the rest are design-level dispositions with tests or independent evidence outstanding.

## 3. New findings

### C1-R-24 — Derived `ProcessAlive` and `Available` contradict the contract for live `unavailable`

**Severity:** High  
**Artifacts:** design §3.4 (`DESIGN...:176-184`), §6.4 (`:385-404`), §8.1 (`:531-532`), C1-AC-18 (`:688`), C1-AC-24 (`:694`); source `supervisor.go:324-326, 351-366`; contract A2-D4 (line 2505), A2-D5/A2-D6 (lines 2506-2507).

**Observed evidence.** The current `processAlive` function returns false whenever `providerState` is `crashed` or `unavailable`, even when `processState` is `alive`, `starting`, or `stopping`:

```go
func processAlive(process, provider) bool {
    if process == alive || process == starting || process == stopping {
        return provider != crashed && provider != unavailable
    }
    return false
}
```

`Status()` calls this at `supervisor.go:324`. The current `providerAvailable` function returns false for `unavailable` at `:358-366`.

**Contract comparison.** A2-D4 says `ProcessAlive` is derived solely from `TTSProcessState`: true for `starting`, `alive`, and `stopping`; false for `absent` and `exited`, and “never derived from aggregate `State`.” A2-D5 says `Available` is true for asset-specific `unavailable` on a live compatible helper.

**Inference.** The r6 design makes live asset-specific `unavailable` a central state but only addresses `ProcessAlive` in the incompatible Case B barrier and after exit. It does not correct these derived predicates. A live `unavailable` helper will therefore report `ProcessAlive=false`, and its `Available` projection will be false, contrary to A2-D4/A2-D5. C1-AC-18 requires the live snapshot to be retained but does not assert the status predicates; C1-AC-24 only tests incompatible termination.

**Status:** **OPEN — high-severity contract contradiction.** C1 cannot claim the live-`unavailable` state is correctly modeled until the derived predicates are fixed or explicitly assigned to C2 with a binding transition and test.

**Required correction:** Add the derived-status contract to the design and acceptance criteria: `ProcessAlive` must be computed only from `processState`; `Available` must be true for live asset-specific `unavailable` and false for exited/no-instance failure where the contract requires. Add status assertions to the live-retain and exited-invalidate tests. If C2 owns the projection, state that C1's state contract makes the C2 result mandatory and do not mark C1-AC-18 complete without it.

### C1-R-25 — Unexpected exit from live `unavailable` is still classified `crashed`

**Severity:** High  
**Artifacts:** design §3.4 (`:182`), §6.5 (`:445-449`), §6.7 `beginTermination`/`observeExit` rows (`:496-497`), §8.2 (`:537-559`), C1-AC-18 (`:688`), matrix X (`:311-318`); source `supervisor.go:166-195`; contract A2-D6 (line 2507), A2-D12 (line 2513).

**Observed evidence.** Current `observeExit` classifies an exit as normal only when `providerState == stopping`; otherwise it writes `providerState=crashed` and sets a crash failure (`supervisor.go:176-191`). The r6 design defines a private termination intent only for `incompatible` (`§6.4`). `beginTermination` is allowed to preserve `unavailable`, but `observeExit` has no `unavailable` intent and the event table says provider state is chosen “per intent” without defining that case.

**Inference.** A live asset-specific `unavailable` helper that exits unexpectedly—or is terminated by cleanup that does not first establish an `unavailable` exit intent—will be observed as `crashed`. That contradicts the design's own live/exited `unavailable` distinction, matrix X (“remains `unavailable`”), and A2-D6's requirement that an unavailable with no authoritative live instance clears the snapshot rather than becoming a crash. C1-R-3 is therefore not fully closed.

**Status:** **OPEN — high-severity design gap.**

**Required correction:** Define an explicit `unavailable` exit/termination intent keyed to the process and generation, or an equivalent observer rule. On exit from a live asset-specific `unavailable`, preserve `State=unavailable`, invalidate the snapshot, set `ProcessState=exited`/`ProcessAlive=false`, retain the failed attempt, and do not enter `crashed`, backoff, budget spend, replacement, or restart. Add separate tests for spontaneous exit and supervisor-driven termination. C1-AC-24 should not be the only barrier test.

### C1-R-26 — Client-side handshake copies remain outside the process-owned attempt

**Severity:** High  
**Artifacts:** design §3.1 (`:134-146`), §6.5 (`:405-433`), C1-AC-23 (`:693`); source `client.go:191-272`, especially `:241, 262-270, 274-292`, and `:877-920`; `protocol.go:123-129`.

**Observed evidence.** `Client.Start` still calls the ordinary `helloHandshake` path (`client.go:241`), obtains a `helloResult`, and writes `c.hello`, `c.providerID`, and `c.engine` directly (`:262-270`) before/around publication. `helloHandshake` (`:274-292`) sends through the ordinary current-transport `callOn`, not through a process-pinned attempt. The design's attempt stores the validated hello and forbids a caller-supplied `helloResult` **for supervisor publication**, but it explicitly permits `Client` to retain `hello` for provider descriptors and provider-scoped work (`§3.1`). It does not say how that copy is generation-bound, how it is obtained only from the successful attempt, or that `helloHandshake` is removed.

**Inference.** A package-internal caller can still obtain a hello through `helloHandshake`/a second `callOn("hello")` path, or `Client.Start` can expose `c.hello`/provider identity before the attempt publishes. `Client.requireProviderCapability` and `Client.selectedProvider` read that copy (`client.go:887-920`), so stale or foreign provider data can influence provider-scoped work even though supervisor capability publication is gated. The design's statement that “the attempt's exchange is the only method that sends hello” is not enforced or tested by C1-AC-23, which only covers publication.

**Status:** **OPEN — R21 is only partially closed.**

**Required correction:** Make the attempt the sole owner of the hello exchange and client copy. Remove or rewrite `helloHandshake`; `Client.Start` must obtain the validated hello from the current attempt, not from a separate call. Publish/mutate `c.hello`, `c.providerID`, and `c.engine` only after successful attempt publication and bind them to the attempt generation. Invalidate the client copy on exit/replacement. Add negative tests for a foreign transport/attempt, a stale client hello, and a provider capability check after generation change. C1-AC-23 needs an explicit “no second hello send” assertion.

### C1-R-27 — The closed event set cannot represent the existing synthesis terminal states, and omits `activeRequestID`

**Severity:** High  
**Artifacts:** design §6.5 event set (`:435-442`), §6.6 inventory (`:454-479`), C1-AC-25 (`:695`); source `client.go:737, 748, 771, 783, 793, 796, 803, 806, 814, 820, 828, 839, 847`, and `supervisor.go:233-242`; contract A2-D4/A2-D6.

**Observed evidence.** The current client writes a `ready` terminal with a failure for timeout, cancellation, and cancellation-before-acceptance (`client.go:771, 783, 793, 803, 839`), and a `degraded` terminal for ordinary failures/decode/validation (`:796, 806, 814, 820, 828`). The r6 event set contains only `synthesisSucceeded` and `synthesisFailed`; `synthesisFailed` is defined as “`degraded` for an ordinary non-fatal failure; preserve failure.” It does not define timeout, cancellation, encode failure, decode/validation failure, or the current `ready`-with-failure outcomes. The same `SetState` method writes `activeRequestID` (`supervisor.go:236`), while §6.6 inventories only `providerState`, `lastFailure`, and `processState`, and C1-AC-25 does not require an `activeRequestID` writer inventory.

**Inference.** A single `synthesisFailed` event cannot preserve the current state distinctions, and a mechanical “all writes through apply” migration would change timeout/cancel/validation outcomes or leave the active request ID unwritten/stuck. This is a C1-R-20 completeness defect, not merely a test omission.

**Status:** **OPEN — high-severity design incompleteness.**

**Required correction:** Define the exact terminal event set (or a parameterized terminal event with a closed outcome enum) for success, ordinary failure, timeout, cancellation, cancellation-before-acceptance, encode failure, decode failure, validation failure, and write failure. Map each current `SetState` line in §6.6. Include `activeRequestID` in the ownership/inventory and acceptance criteria. Require a table-driven test for every terminal disposition, not only writer rejection after incompatible.

### C1-R-28 — Removing `SetSelected` removes the catalog selection path, and the C7 adapter is still delegated to the direct method

**Severity:** High  
**Artifacts:** design §4.7 (`:335-347`), §6.5 (`:433`), §6.6 rows (`:469, 477`), C1-AC-7/C1-AC-12/C1-AC-26 (`:677, 682, 696`); source `client.go:517-548`, especially `:546`, and `:961-963`; contract §17.12.1 (lines 2266-2267).

**Observed evidence.** Current `Client.Catalog` calls `SetSelected(result.ProviderID, result.Engine, "", nil)` after a successful catalog (`client.go:546`). Current `ProviderCatalog` is a `Client` method that delegates to `Catalog` (`:961-963`). The r6 design removes `SetSelected`, says identity is set only by handshake publication or `loadSucceeded`, and distinguishes a “C7 `ProviderCatalog`” method that may start from ordinary `absent`—but does not define that method or its selection behavior. The existing `ProviderCatalog` cannot both delegate to direct `Catalog` (now no-start) and retain its §17.12.1 start/handshake authority.

**Inference.** The removal of `SetSelected` loses the current catalog-driven provider/engine selection path. If a catalog request selects an engine different from the handshake's first provider, subsequent `Synthesize`/provider-scoped work has no replacement selection rule. Separately, the claimed C7 adapter start path has no concrete entry point in the reviewed design; the only current entry point delegates to the no-start direct path.

**Status:** **OPEN — R22 is only partially closed.**

**Required correction:** Decide and specify provider-selection semantics after removing `SetSelected`: either preserve a read-only catalog-selection event, or make selection explicit through `load`/C7 and state how catalog results map to later synthesis. Define a distinct C7 adapter entry point that owns the absent-only start/handshake and does not call direct `Client.Catalog`. Add direct-versus-adapter selection tests and a test that catalog cannot silently change a loaded/current provider without explicit selection.

### C1-R-29 — C1-AC-24 and the writer test contradict `beginTermination`

**Severity:** Medium  
**Artifacts:** design C1-AC-24 (`:694`), §6.6 `beginTermination` row (`:496`), test matrix “Writer rejection after incompatible” (`:779`), barrier test (`:748`).

**Observed evidence.** C1-AC-24 says every writer inventory entry after `incompatible` “either rejects without mutation or is the authorized exit/update-reset path,” while also saying `Terminate` may set `ProcessState=stopping`. The `beginTermination` row is allowed to set `processState=stopping`. The test row says each writer/event after `incompatible` leaves “process state unchanged,” which contradicts the required barrier observation `ProcessState=alive`/`stopping` and the explicit `Terminate` permission.

**Inference.** At least one acceptance test cannot pass as written. The design needs to distinguish “aggregate state and `LastError` unchanged” from “only `processState` may transition to `stopping`” for `beginTermination`. The current wording is an internal contradiction, not a test interpretation issue.

**Status:** **OPEN — medium-severity acceptance-test contradiction.**

**Required correction:** Split the writer-rejection test into (a) ordinary events that must mutate nothing and (b) the authorized `beginTermination` path, which may change only `processState` and must preserve aggregate `incompatible` and `LastError`. Rewrite C1-AC-24 so it does not claim that every writer rejects or is exit/update-reset. Keep the barrier assertion separate.

### C1-R-30 — The “six-state × seven-operation” matrix is not uniform

**Severity:** Medium  
**Artifacts:** design §4.6 (`:245-334`), §4.7 (`:335-347`), C1-AC-26 (`:696`), verification row (`:729`).

**Observed evidence.** The matrix contains eight rows in `A`, `X`, and `U` because direct `Catalog` and `ProviderCatalog` are separate, but only seven rows in `K0`, `K1`, and `L`; those three classes omit a distinct `ProviderCatalog` row. The design/AC calls the matrix “six state classes × seven operations” and also says direct and C7 Catalog are distinct. The direct/adapter distinction therefore cannot be both a seven-operation matrix and an eight-row matrix without explanation.

**Inference.** The matrix is broader than r5, but its stated completeness claim is inaccurate. A reviewer cannot verify “no blank cell” when the operation set changes by state class, and the C7 adapter rows in the live known-absent/present/unavailable classes are absent from the primary matrix.

**Status:** **OPEN — medium-severity traceability/completeness defect.**

**Required correction:** Define a fixed operation set, with `Catalog` as one logical operation containing two explicitly named paths, or make it six classes × eight operations. Fill adapter `ProviderCatalog` behavior in every state class in which it is permitted to be distinct, or state that §4.7 is the sole adapter table and cross-reference every cell. Update C1-AC-26 and its verification row to match the actual count.

### C1-R-31 — No-frame `request_not_found` is a wire-origin code without explicit adapter authority

**Severity:** Medium  
**Artifacts:** design §4.2 (`:201-209`), §4.5 (`:239-244`), §4.6 L row (`:302`), C1-AC-18/C1-AC-19 (`:688-689`), test matrix live-unavailable cancel (`:769`); contract §17.9 registry (`DESIGN...:1902`), §17.4.6 (`:1396-1402`), A2-D1/A2-D10.

**Observed evidence.** The design returns `request_not_found` for live-`unavailable` `Cancel` with **no frame**. The contract registry identifies `request_not_found` as origin **wire** and defines it as “Cancel target is not active.” The contract explicitly authorizes a no-frame `unsupported_capability` for an absent optional capability (A2-D1), but no equivalent express authority is cited for an adapter-synthesized wire-origin `request_not_found`.

**Inference.** The design may be semantically right that no synthesis is active, but it has not established that a no-frame `request_not_found` is the contract-authorized representation. It is either an unsupported app-side fabrication of a wire error, or it needs one current-transport cancel frame (which would also produce the exact wire result). The current text states a conclusion without the required authority.

**Status:** **OPEN — medium-severity error-provenance ambiguity.**

**Required correction:** Either (a) send exactly one cancel frame on the current transport and use the helper's `request_not_found` terminal, with no state/recovery change, or (b) cite an explicit contract/owner authority for adapter-synthesized no-frame `request_not_found` and add it to C1-AC-18/C1-AC-19. Do not describe it as both “exact contract result” and “no frame” without that authority.

### C1-R-32 — The C6/C7 reset exception is asserted to be token-gated but the token authority is undefined

**Severity:** High  
**Artifacts:** design §6.5 `updateGatedReset` (`DESIGN...:443-447`), §6.7 reset row (`:501`), C1-AC-24/C1-AC-25 (`:694-695`); source `supervisor.go:233-242` (current exported setter being removed); contract A2-D4/A2-D12 and r9 §10.2/§12.

**Observed evidence.** The design says `updateGatedReset` requires `(token, updateEpoch)` from the C6/C7 helper/assets/config update path, a terminal state, and no live process. It does not define the token type, issuer, creation/consumption store, or binding to the actual helper/assets/config change. It says “C1 defines this guard shape; C6/C7 must supply the actual update evidence.” The event is private, but any package-internal caller can construct the event; C1-AC-24/C1-AC-25 do not test a forged/reused/foreign token or an ordinary package-internal call. The design therefore relies on a future stage's unspecified token rather than a C1-enforceable narrow exception.

**Inference.** The reset exception is the only ordinary path that can clear `incompatible` or an `unavailable` terminal failure. Calling the event with a plausible token and epoch would clear the failure and advance the generation; there is no demonstrated type-level or ownership barrier that prevents a forged token or an accidental caller. Because C1-AC-25 says “only `observeExit` and a token-gated C6/C7 `updateGatedReset` may change `incompatible`,” this is a direct weakness in the R20 closure claim. It is high severity even though C6/C7 is a later stage: C1 is defining the guard shape and acceptance criteria now.

**Status:** **OPEN — high-severity R20 boundary gap.**

**Required correction:** Define the update token as an opaque, typed value issued only by the C6/C7 update path, with an explicit one-use/epoch store and an ownership check that the caller cannot fabricate. Add negative tests for a forged token, a reused token, a stale epoch, a token from another process/generation, and a package-internal ordinary caller. If C1 cannot implement the issuer, C1 must make the reset event unreachable until C6/C7 supplies a verified token and must not claim R20 is closed before that gate exists.

## 4. Acceptance criteria and test audit

The design has 26 criteria and 26 verification rows, but the rows do not fully prove the invariants:

- **C1-AC-18** does not assert `ProcessAlive` or `Available` for live `unavailable`; see C1-R-24.
- **C1-AC-18** does not test unexpected exit from live `unavailable`; see C1-R-25.
- **C1-AC-23** tests publication rejection but not the surviving `Client.hello`/`helloHandshake` path; see C1-R-26.
- **C1-AC-24/C1-AC-25** have the `beginTermination` contradiction and do not cover every synthesis terminal; see C1-R-27/C1-R-29.
- **C1-AC-26** overstates matrix uniformity; see C1-R-30.
- **C1-AC-18/C1-AC-19** do not resolve the no-frame cancel code origin; see C1-R-31.

The test matrix is stronger than r5, but several tests assert only final state and omit the forbidden intermediate action or the exact derived projection. A passing final-state assertion cannot close C1-R-24, C1-R-25, or C1-R-27.

## 5. Resource, I/O, and acceptance evidence

- **C1-AC-15 remains OPEN.** No independent written budget conformance determinations or named/dated I/O N/A reviews were found.
- `MaxPending=32` remains unresolved; the design correctly does not claim an authorized source.
- All 13 resource rows remain pending independent reviewer determination.
- All eight I/O categories are documented, but reviewed N/A entries still lack the required independent reviewer/date.
- No tests, runtime checks, benchmarks, failure logs, or conformance results were executed in this review.
- The r6 design does not establish `ProcessAlive`/`Available` status behavior for live `unavailable`; this is a design defect independent of resource acceptance.

## 6. Overall verdict

**FAIL.**

The artifact identities and preserved inputs are correct, and r6 substantially improves the prior design: it correctly separates direct and adapter behavior, adds a process-owned publication attempt, and expands the operation matrix. However, the following unresolved high-severity issues prevent design approval:

1. `ProcessAlive`/`Available` still contradict A2-D4/A2-D5 for live `unavailable` (C1-R-24).
2. Unexpected exit from live `unavailable` still classifies as `crashed` (C1-R-25).
3. Client-local hello and provider identity remain outside the process-owned attempt, and the old `helloHandshake` path remains (C1-R-26).
4. The closed event set cannot represent existing synthesis terminal states, and `activeRequestID` is omitted from the inventory (C1-R-27).
5. Catalog selection and the C7 adapter start path are not actually designed after `SetSelected` removal (C1-R-28).

C1-R-29 through C1-R-32 are additional correctness/traceability defects that must be fixed before the next review. No finding is acceptance-closed. C1 remains unauthorized; this review does not authorize implementation, testing, runtime execution, or deployment.

## 7. Review boundary

This was a read-only independent review. No r6 artifact, source file, governing input, owner decision, or prior revision was modified. No test, build, runtime check, benchmark, or deployment was performed. The review is submitted as the separate independent r6 review; the author must produce a new revision rather than treating this report as acceptance.
