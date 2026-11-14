# C1 Design Revision 6 — Capability Negotiation and Classification

**Date:** 2026-11-08  
**Revision identity:** C1-R6  
**Status:** PROPOSED C1 DESIGN REVISION — NOT ADOPTED; C1 UNAUTHORIZED; ACCEPTANCE EVIDENCE OPEN  
**Mode:** Design and documentation only. No source, test, build, runtime, benchmark, or deployment change.  
**Current proposal:** Revision 6 is the current proposed C1 design. It supersedes only the proposed C1 design and acceptance positions in revision 5 for review purposes. Revision 5, revision 4, revision 3, and every earlier artifact remain preserved and immutable. The frozen plan remains unchanged and historical.  
**Correction scope:** C1-R-20 through C1-R-23 from the independent review of revision 5. Valid revision-5 corrections are retained.

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
| Revision 3 C1 design (historical) | `e9b30c6d9072375750851d13a05c2c53bb5cb983ec69498418a404d0291c754e`; 455 / 46,369 |
| Revision 3 C1 readiness (historical) | `aca2d59921709c7e0239547d1b172572a790140f185389ffdc00a04e2b714bd7`; 138 / 13,476 |
| Prior C1 design revision | `003827f12658d1e556657220754795bff2094b39dca960321a71c3fb302cec6d`; 313 / 22,680 |
| Prior C1 readiness report | `45457897844a5fe90c247bf97225420df0700ea431d2a7ac65cf14617867692d`; 129 / 8,910 |
| Revision 3 review disposition | Muse, `CONDITIONAL — SPECIFIC CORRECTIONS REQUIRED`; findings C1-R-1 through C1-R-9. No separate report file was present in the workspace; each cited requirement was reopened and checked directly. |
| Revision 4 review disposition | Findings C1-R-10 through C1-R-14 as supplied by the owner. No separate report file was present; each finding was author-reconciled against the governing contract, r9, D10, revision 3, and the current source. |
| Revision 4 design (historical) | `b099afe0302672bf0459ae9ec1b493e64bd11c021c6568d2a6a5ab3e08c78e3f`; 539 / 65,849 |
| Revision 4 readiness (historical) | `e291967f56312bb07c177b2d832a1a69683a2410165def6d77b7d1c255ec7e5b`; 128 / 11,076 |
| Revision 5 design (preserved) | `6469203d9b8ae869988cf862b293f464e3e788c47badbc964f2534a56e04c387`; 656 / 89,192 |
| Revision 5 readiness (preserved) | `67ce40b5759aa3ea85b647148ff62cb40c71922048f02d68aa084ded33275179`; 146 / 13,514 |
| Revision 5 integrity report (preserved) | `b1829b4eaa31ac826a4bbc5da9703dd5e31d0b202dba2bb43ac230f0be3cbf06`; 73 / 5,736 |
| Revision 5 independent review disposition | Muse, `CONDITIONAL — SPECIFIC CORRECTIONS REQUIRED`; findings C1-R-20 through C1-R-23. The review is independent of the r5 author reconciliation. |

The contract controls. Findings are **author-reconciled**, not independently accepted. The evidence ladder is: design content exists; author reconciliation against governing text; independent review where explicitly evidenced; acceptance evidence and closure. These levels are not collapsed. The owner decisions record D10(a) = A10a-2 and D10(b) = A10b-2. This revision records no new owner decision and reopens none.

## 2. Finding-by-finding reconciliation

All findings C1-R-1 through C1-R-23 are author-reconciled against the cited governing text. C1-R-1 through C1-R-19 retain their r5 design dispositions; C1-R-20 through C1-R-23 are the r6 corrections. None is independently accepted, and none is acceptance-closed. The authoritative per-finding disposition is §2.4. Acceptance evidence remains open wherever noted, especially C1-AC-15. No finding is rejected or silently narrowed.

| Finding | Governing requirement checked | Counterexample verified | Disposition and correction |
|---|---|---|---|
| **C1-R-1 — Start/restart map** | r9 §8 C1 (line 589) and §8.1; contract A2-D4; §17.11.5 / §17.7.2 | Valid. Revision 2 omitted `crashed → backoff → starting`, fourth-attempt exhaustion, normal exit-zero/nonzero/grace classification, and a complete §8.1 health-edge table. | **Design-corrected; author-reconciled.** §8 now supplies the complete map and §8.2 the state-scoped health table. C1-AC-13 and the matrix assert completeness. Executed restart remains C5; executed start/reachability remains C7, as r9 allocates. |
| **C1-R-2 — Resource and I/O evidence** | r9 §9.1, §9.2, §10, §3.2, R8-AC-4/7/8/9 | Valid. Revision 2 listed resource values but not basis, measurement, acceptance threshold, or written conformance determination; its I/O evidence was a single assertion, not per-category records. | **Design-corrected; author-reconciled.** §9.2 supplies every §10 field; §9.3 supplies every §9.1 I/O field for all eight §3.2 categories; §9.4 supplies R8-AC-8 and R8-AC-9. C1-AC-15 requires complete, non-blank records and independent determinations. |
| **C1-R-3 — Asset-specific `unavailable`** | Contract A2-D5 and A2-D6 (lines 2506) | Valid. A live handshaken asset-specific `unavailable` retains the authoritative snapshot; an exited/no-instance `unavailable` clears it. Revision 2 did not make the two states separately testable. | **Design-corrected; author-reconciled.** §3.2 rule 7 and §3.4 distinguish them. C1-AC-18 requires separate live-retain and exited-invalidate tests. |
| **C1-R-4 — Semantic handshake termination** | Contract A2-D4 (line 2505) | Valid. A2-D4 requires an intent-aware exit observer. Current `Client.Start` calls `Terminate`, then `SetState(incompatible)`; if termination times out and the exit arrives later, `Supervisor.observeExit` treats the exit as unexpected and writes `crashed`. That would permit backoff, budget spend, or replacement contrary to A2-D4. | **Design-corrected; author-reconciled.** §6.4 specifies an expected-termination intent latched before termination. Late exit preserves `incompatible`; termination timeout does not reclassify. C1-AC-3 adds timeout and late-exit cases. |
| **C1-R-5 — D10 caller-facing projection** | Owner D10 A10a-2/A10b-2; `errors.go` `Failure.AudioError` | Valid. The owner text requires `ErrCodeSynthesisFailed`, preserved precise code, and bounded detail. Revision 2 named the projection but did not make it an explicit acceptance/test assertion. | **Design-corrected; author-reconciled.** §4.3 gives the exact internal and caller-facing shapes. C1-AC-5, C1-AC-6, and C1-AC-19 require both the helper-boundary failure and the caller-facing `AudioError`. |
| **C1-R-6 — Traceability corrections** | r9 §7.3 original AC-7; §7.1 matrix; R8-AC-7; contract A2-D8 | Valid. Original AC-7 is C8-only no-substitution, not I/O evidence. R8-AC-7 is the I/O rule. Revision 2 wrote `AC-7/R8-AC-7` and claimed A2-D8 coverage without shutdown design content. | **Design-corrected; author-reconciled.** §12 separates original AC-7 (C8-only, no C1 claim) from R8-AC-7. A2-D8 is now an explicit non-substitution boundary: `Unload` is not idle shutdown and cannot terminate the process. |
| **C1-R-7 — Supervisor known-absent capability** | Contract A2-D9 (line 2510); owner D10(b) boundary | Valid. Both unknown and known-absent skip the frame and synthesize no `unsupported_capability`, but they are distinct states and revision 2 combined their test. | **Design-corrected; author-reconciled.** §5 separates the branches; C1-AC-20 and a separately named test cover known-absent. |
| **C1-R-8 — Concurrency coverage** | Original AC-1, AC-2, AC-3; r9 §8 C1 TH-1 focus | Valid. Revision 2 covered gate-versus-exit and handshake-versus-gate but not concurrent direct `Unload`/`Cancel`, optional operations during lifecycle states, or supervisor cancellation during startup/shutdown/recovery. | **Design-corrected; author-reconciled.** §9.1 documents existing serialization and winner rules; C1-AC-21 and the concurrency matrix require the missing cases. No new architecture is invented. |
| **C1-R-9 — Exact protocol and evidence details** | Contract A2-D3; R8-AC-8, R8-AC-9; C1-AC-16 | Valid. The exact `"missing"` key, table order, omission rule, PR-dimension statement, reviewed N/A entries, and static-review assertion were not all explicit. | **Design-corrected; author-reconciled.** §6.1 states the exact key/order/omission. §9.4 states PR impact and reviewed N/A. C1-AC-22 requires the C1-AC-16 static review. |

### 2.1 Revision 4 findings (C1-R-10 through C1-R-14)

| Finding | Requirement checked | Author verification against governing text | Disposition and correction |
|---|---|---|---|
| **C1-R-10 — Readiness overstates evidence closure** | r9 §7.1/§7.2, §10, R8-AC-4/7/9; C1-AC-15 | Valid. A complete design table is not an independent determination, and the revision-3 readiness used “closed at the design/acceptance level.” | **Author-reconciled; acceptance evidence pending.** Readiness 5 separates design content, author reconciliation, independent review, and acceptance closure. C1-AC-15 is expressly **OPEN**. C1-R-2 remains an evidence gap, not acceptance closure. C1-R-4 remains a design correction; its executed termination evidence is prospective. |
| **C1-R-11 — Live `unavailable` can be cleared by direct `Unload`** | Contract A2-D12 (line 2513); A2-D6; A2-D7 | Valid. Current `Client.Unload` sets `TTSStateUnloaded` after a successful response. On a live handshaken helper in asset-specific `unavailable`, that clears the failed-attempt state and lets ordinary synthesis follow the ordinary unloaded ensure-load path. | **Design-corrected; author-reconciled.** §4.1 now permits only a current-transport unload frame and forbids the aggregate `unloaded` transition; §4.2 resolves `Cancel` separately as `request_not_found`; §4.5 permits read-only `Health` and `Catalog` without recovery. The r4 `provider_not_loaded` guard for direct optional operations is withdrawn because that code is synthesis-specific. Snapshot retention/invalidation is unchanged; only explicit same-identity `load` is the recovery bypass, owned by C3/C6. |
| **C1-R-12 — Incompatible timeout still permits a later start/handshake** | Contract A2-D4 (line 2505); A2-D5; A2-D7 | Valid. The revision-3 intent-aware exit observer prevents a late exit from becoming `crashed`, but it does not prevent a later `Client.Start` or `Supervisor.Start` from returning the still-live incompatible process and running a second handshake. | **Design-corrected; author-reconciled.** §6.4 separates termination success from termination timeout and adds an explicit `incompatible` gate in both `Client.Start` and `Supervisor.Start`. The live-but-unreaped process remains tracked; no second start, handshake, transport reuse, frame, backoff, budget spend, or replacement. C1-AC-3, C1-AC-12, C1-AC-13, and C1-AC-21 and the test matrix cover repeated observation, ordinary TTS, direct start, later exit, and concurrency. |
| **C1-R-13 — Wrong helper-process cap source** | Contract §5.5 (line 520); frozen §17 material; r9 §8 C7 | Valid. Revision 3 cited §17.11.5, which governs crash/restart/shutdown, not the two-process hard cap. | **Design-corrected; author-reconciled.** §9.2 now cites contract §5.5 (line 520), repeated in the frozen §17 material, with r9 §8 C7 only as the stage overlay. The C1 no-start assertion remains separate from the global cap. |
| **C1-R-14 — Acceptance and HW-1.6 precision** | Contract A2-D1; r9 §5.6 HW-1.6; C1-AC-17 | Valid. C1-AC-1 did not sharply distinguish mandatory control operations from optional service IDs; C1-AC-17 lacked a verification disposition; HW-1.6 did not state CPU/memory/disk/GPU assumptions. | **Design-corrected; author-reconciled.** C1-AC-1 now separates mandatory service floor from mandatory control operations. C1-AC-17 is a delivery-only criterion verified by a static package checklist. §9.5 records CPU-core, memory, disk, and GPU/accelerator assumptions and their review status. All 26 criteria have a verification and traceability row in §10.1. |

### 2.2 Revision 5 findings (C1-R-15 through C1-R-19)

| Finding | Governing text checked | Author verification and source fact | Design correction |
|---|---|---|---|
| **C1-R-15 — Live `unavailable` Health/Catalog contradiction** | Contract §17.4.1 operation preconditions; §17.4.4 `unload`; §17.4.5 health; §17.12.1 `ProviderCatalog`; A2-D12 (line 2513); r9 §8.1 live-`unavailable` row | Valid. Revision 4 §4.5 returned `provider_not_loaded` and suppressed both `Health` and `Catalog`, while §8.1/§8.2 required a live wire health observation that remains `unavailable`. Current `Client.Health` and `Client.Catalog` can send on an existing transport, but both call `Start`; `Catalog` also calls `SetSelected`. The contract permits `health` after handshake and `catalog` after handshake with no load required. | **Design-corrected; author-reconciled; acceptance evidence pending.** `Health` is a mandatory control read and MUST send one current-transport `health`; `Catalog` is permitted when advertised and MUST send one current-transport `catalog`. Neither may call `Start`, spawn, handshake, load, recover, change aggregate state, clear `LastError`, or mutate the identity/snapshot; `Catalog` MUST NOT call `SetSelected` in this state. `Health` remains wire-based for live `unavailable`; exited/no-instance uses supervisor status. |
| **C1-R-16 — Exact error semantics** | Contract §4.6 typed errors; §17.4.4 unload; §17.4.6 cancel; §17.9 registry; A2-D12; A2-D1/A2-D7 | Valid. `provider_not_loaded` is defined as “Synthesis before successful load” and A2-D12 assigns it only to `synthesize` in `unavailable`. Revision 4 extended it to direct `Unload`, direct `Cancel`, `Catalog`, and `Health` without authority. | **Design-corrected; author-reconciled; acceptance evidence pending.** `Synthesize` alone returns `provider_not_loaded`, no frame/start/load. `Health` and advertised `Catalog` proceed as read-only current-transport operations. Advertised `Unload` proceeds on the current transport and returns its own result, but MUST NOT set aggregate `unloaded` or clear the failed attempt. Advertised `Cancel` returns the contract's exact non-active-target result `request_not_found` without a frame or recovery; an absent optional capability returns `unsupported_capability` with no frame. Exact projections are in §4.6. |
| **C1-R-17 — `MarkHandshake` publication bypass** | Contract A2-D3, A2-D4, A2-D6; §17.11.3 step 7; current `client.go:270`, `supervisor.go:205-220` | Valid. `Supervisor.MarkHandshake` is exported and has no process-instance, generation, lifecycle, or terminal-state gate. It writes `processState=alive`, `providerState=unloaded`, helper/protocol versions, capabilities, and selected provider data. A direct or stale call can therefore publish after `incompatible`/`stopping` or publish for a process the supervisor no longer owns. | **Design-corrected; author-reconciled; acceptance evidence pending.** Revision 6 goes beyond a method rename: publication is owned by a private `handshakeAttempt` that captures `Process`/generation and stores the validated hello read through that process. There is no caller-supplied `helloResult`; the old exported method is removed. `incompatible`, `stopping`, exited/crashed/backoff, `unavailable`, replacement, stale generation, foreign attempt, and mismatched selected-provider attempts are refused without mutation. C1-AC-23 and the expanded negative matrix cover it. |
| **C1-R-18 — Aggregate `incompatible` during termination** | Contract A2-D4 (lines 2505); §17.7.1 process/provider separation; current `supervisor.go:246-280`, `169-192` | Valid. Current `Supervisor.Terminate` sets `providerState=stopping` whenever a process exists, including after an incompatible handshake. The current `observeExit` classifies any exit not seen as `stopping` as `crashed`; it has no incompatible intent. This contradicts A2-D4's requirement that aggregate `State` remain `incompatible` while termination is pending and after the expected exit. | **Design-corrected; author-reconciled; acceptance evidence pending.** `Terminate` owns process lifecycle only: it sets `ProcessState=stopping` but MUST NOT overwrite aggregate `incompatible`. The intent-aware observer sets `ProcessState=exited`, leaves `State=incompatible`, and suppresses crash/backoff/budget/replacement. A barrier test observes `State=incompatible` while termination is still pending. C1-AC-24 covers the projection. |
| **C1-R-19 — Independence wording** | r9 §7.1/§7.2/§7.3; adoption record | Valid. Revision 4 readiness described C1-R-10 through C1-R-14 as “independently checked,” but the producing agent performed the reconciliation and no independent finding-specific review is recorded. | **Author-reconciled.** Revision 5 uses “Author-reconciled against governing text; independent review completed separately only where explicitly evidenced.” No resource, I/O, or stage evidence is called independently accepted. The four evidence levels remain separate. |

### 2.3 Revision 6 findings (C1-R-20 through C1-R-23)

| Finding | Governing requirement checked | Counterexample verified | Design correction |
|---|---|---|---|
| **C1-R-20 — Aggregate-state mutation surface** | Contract A2-D4 (line 2505), A2-D5/A2-D6, §17.7.2, §17.7.4; C1-AC-24; current `supervisor.go:233-242` | Valid. `Supervisor.SetState` is exported and writes `providerState`, `lastFailure`, and `processState` without a terminal-state guard. A call after `incompatible` can write `unloaded`, clear the terminal failure, and set `processState=alive`; the `Client.Start`/`Supervisor.Start` gates then see no `incompatible` and can return the still-live process. The same method is the normal writer for every load, unload, synthesis, and startup-failure path. | **Design-corrected; author-reconciled.** §6.6 inventories every current writer. §6.5 removes the exported setter and makes a private lifecycle transition the sole writer of `providerState`, `lastFailure`, `processState`, and `activeRequestID`. `incompatible` admits only the expected-exit observer and a C6/C7 update-gated reset token; no ordinary event, handshake, termination, or alternate exported method can overwrite it. §10.1/§11 add writer-specific negative tests. C1-AC-24 and new C1-AC-25 verify it. |
| **C1-R-21 — Handshake provenance** | Contract A2-D3, A2-D4, A2-D6; C1-AC-23; `helloResult` at `protocol.go:123-129`; `supervisor.go:205-219`; `client.go:270` | Valid. The r5 gate checks the current `Process`, generation, state, and `validateHello(h)`, but `helloResult` contains no process/transport token. A package-internal caller with the current process could publish a valid hello obtained from another process. The test row covers only `incompatible`, `stopping`, and old generation. The phrase “old method removed or made unexported” could leave a second ungated method. | **Design-corrected; author-reconciled.** §6.5 introduces a private `handshakeAttempt` that captures `p`, generation, and the validated hello read through a pinned exchange on `p`; publication uses only the attempt's stored hello and validates `selected` against it. No caller-supplied `helloResult` exists. The old exported method is removed, not renamed. §11 expands the negative matrix to every terminal/lifecycle state, foreign process, selected mismatch, replacement, stale generation, concurrent termination, plus a positive control. C1-AC-23 is rewritten. |
| **C1-R-22 — `ProviderCatalog` versus direct `Client.Catalog`** | Contract §17.12.1 (lines 2266-2267), §17.4.1 `catalog`, A2-D12; C1-AC-12; current `client.go:517-548`, `959-963` | Valid as a boundary ambiguity. `Client.Catalog` calls `Start`; `ProviderCatalog` delegates to it. C1-AC-12 says no direct optional operation calls `Start`, while §17.12.1 permits `ProviderCatalog` to start/handshake from `absent`. The design did not say which component owns that start. | **Design-corrected; author-reconciled.** §4.7 distinguishes the C7 adapter method `ProviderCatalog` (may start/handshake only from ordinary `absent`) from the direct C1 `Client.Catalog` (never starts; current instance only). From `unavailable`, neither starts or recovers; live reads use the current transport, exited/no-instance fails closed. The same split applies to `ProviderStatus` versus direct `Health`. C1-AC-7, C1-AC-12, C1-AC-18, C1-AC-19 and C1-AC-26 are updated. |
| **C1-R-23 — Complete operation/error projection** | A2-D1, A2-D6, A2-D7, A2-D10, A2-D12; contract §4.6, §17.4.1, §17.4.4-§17.4.6, §17.9 | Valid. The r5 matrix covered only live, handshaken, asset-specific `unavailable`. It did not give exact rows for ordinary `absent`, unknown capability, known-absent optional, exited/no-instance `unavailable`, or direct versus adapter `Catalog`. It also omitted known-absent `Catalog` from C1-AC-7. | **Design-corrected; author-reconciled.** §4.6 is replaced by a six-state-class by seven-operation matrix with exact frame, start, recovery, helper result, `Failure`, `AudioError`, state/failure, and delegation columns. §4.7 adds the adapter start table. `Synthesize` remains `provider_not_loaded` under A2-D12 for any `State=unavailable`; exited recovery is expressly delegated to C3/C5/C6 and is not authorized by C1. C1-AC-7, C1-AC-9, C1-AC-18, C1-AC-19 and new C1-AC-26 are updated. |

### 2.4 Consolidated finding disposition (C1-R-1 through C1-R-23)

| Finding | Disposition | Evidence/acceptance qualifier |
|---|---|---|
| C1-R-1 — Start/restart map | **design-corrected** | C5/C7 execution pending |
| C1-R-2 — Resource/I/O evidence | **design-corrected** | **acceptance evidence pending**; C1-AC-15 OPEN |
| C1-R-3 — Asset-specific `unavailable` | **design-corrected** | tests prospective |
| C1-R-4 — Semantic handshake termination | **design-corrected** | **acceptance evidence pending**; termination evidence prospective |
| C1-R-5 — D10 projection | **design-corrected** | projection tests prospective |
| C1-R-6 — Traceability | **design-corrected** | design traceability only |
| C1-R-7 — Supervisor known-absent | **design-corrected** | test prospective |
| C1-R-8 — Concurrency | **design-corrected** | race/stress evidence pending |
| C1-R-9 — Protocol/evidence details | **design-corrected** | static design only |
| C1-R-10 — Readiness overstatement | **author-reconciled** | **acceptance evidence pending**; C1-AC-15 OPEN |
| C1-R-11 — Direct `Unload` clears `unavailable` | **design-corrected** | tests prospective |
| C1-R-12 — Incompatible start bypass | **design-corrected** | **acceptance evidence pending** |
| C1-R-13 — Process-cap citation | **design-corrected** | citation verified; budget conformance pending |
| C1-R-14 — AC/HW-1.6 precision | **design-corrected** | **acceptance evidence pending**; HW-1.6 review pending |
| C1-R-15 — Health/Catalog contradiction | **design-corrected** | **acceptance evidence pending**; tests prospective |
| C1-R-16 — Error semantics | **design-corrected** | **acceptance evidence pending**; tests prospective |
| C1-R-17 — `MarkHandshake` bypass | **design-corrected** | **acceptance evidence pending**; C1-AC-23 provenance/tests expanded in r6 |
| C1-R-18 — Incompatible aggregate during termination | **design-corrected** | **acceptance evidence pending**; C1-AC-24 barrier and writer tests expanded in r6 |
| C1-R-19 — Independence wording | **author-reconciled** | independent review of r5 completed; r6 independent review pending |
| C1-R-20 — Aggregate-state mutation surface | **design-corrected** | **acceptance evidence pending**; C1-AC-25 prospective |
| C1-R-21 — Handshake provenance | **design-corrected** | **acceptance evidence pending**; C1-AC-23 prospective |
| C1-R-22 — Catalog start boundary | **design-corrected** | **acceptance evidence pending**; C1-AC-26 prospective |
| C1-R-23 — Complete projection matrix | **design-corrected** | **acceptance evidence pending**; C1-AC-26 prospective |

**independently verified:** none in this revision. The r5 independent review supplied findings C1-R-20 through C1-R-23; the r6 author reconciliation is a design correction, not acceptance. **unresolved:** C1-AC-15 and the `MaxPending=32` authorized-source/R8-AC-4 determination; neither is a finding disposition. No finding is marked accepted or closed.

### 2.5 Observed source facts versus proposed behavior

These are current-source observations, not claims that C1 is implemented, and not execution evidence. Source locations are exact for revision 6 and are evidence for the corrections only. The complete writer inventory is §6.6.

| Area | Observed current source fact | Proposed C1 behavior |
|---|---|---|
| Direct `Unload` | `client.go:599-618` calls `call` (which calls `Start`), sends `unload`, and sets `TTSStateUnloaded` after success. | Current-transport unload only; no aggregate `unloaded` transition in live `unavailable`; no `Start`/recovery. |
| Direct `Cancel` | `client.go:622-647` calls `Start` through `call`, sends a cancel frame, and has no no-active-target precheck. | In live `unavailable`, return the contract's `request_not_found` with no frame/start; absent `cancel` returns `unsupported_capability`. |
| `Health` | `client.go:501-515` calls `call`, so it can enter `Start`; it is otherwise read-only. | Current-transport wire health only for a live handshaken helper; no `Start`, handshake, load, or transition. No-instance uses the supervisor observation, not a wire frame. |
| `Catalog` | `client.go:517-548` calls `Start`, requires `catalog`, sends, then calls `SetSelected`; `ProviderCatalog` at `client.go:959-963` delegates to it. | Direct `Client.Catalog` never starts; `ProviderCatalog` is the C7 adapter start path only from ordinary `absent`. No `SetSelected`; live `unavailable` reads current transport; exited/no-instance fails closed. |
| `MarkHandshake` | `supervisor.go:205-220` is exported and ungated; `client.go:270` is its only current call. | Removed. Publication uses a private process-owned `handshakeAttempt`; no caller-supplied hello result; no stale/terminal mutation. |
| `SetState` | `supervisor.go:233-242` is exported and writes `providerState`, `lastFailure`, and `processState`; `client.go` calls it from every load, unload, synthesis, and startup-failure path. | Removed. A private lifecycle transition is the only writer; terminal rejection mutates nothing. |
| `SetSelected` | `supervisor.go:223-230` is exported; `client.go:546` calls it after catalog and `client.go:593` after load. | Identity is set only by handshake publication and successful load; `Catalog` is read-only; no standalone exported setter. |
| `Terminate` | `supervisor.go:246-288` sets both `ProcessState=stopping` and `ProviderState=stopping` whenever a process exists. | Process state may become `stopping`; aggregate `incompatible` and live `unavailable` are preserved; no general aggregate write. |
| `observeExit` | `supervisor.go:166-195` treats an exit as normal only when `providerState == stopping`; otherwise it writes `crashed`. | Consume the private termination intent for the tracked instance; preserve `incompatible`; no crash/backoff. |
| Failed handshake | `client.go:244-253` calls `Terminate` and only afterward calls `SetState(incompatible)`; a late observer can race the state write. | Latch private intent before termination; aggregate `incompatible` is never overwritten by `stopping`; late exit preserved. |
| Synthesize send failure | `client.go:756` writes `TTSStateCrashed` directly after a transport write failure, before the exit observer. | The caller requests termination; only `observeExit` may classify `crashed`. No direct crashed write. |
| Start gate | `supervisor.go:110-164` returns an existing `proc` before checking terminal state; `client.go:191-272` similarly returns an existing transport. | Reject `incompatible` before returning or starting; R12 remains in force. |

Future rows in §10.1 and §11 are prospective; they are not executed evidence.

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
2. Only a successful compatible handshake for the current instance publishes `known=true` with all five service IDs. The supervisor opens a private `handshakeAttempt` bound to the captured `Process` and generation. The attempt performs the pinned `hello` exchange and stores the validated `helloResult` itself; see §6.5. Publication through `attempt.publish(selected)` requires the attempt still to be the current attempt, the captured `Process` to equal `Supervisor.proc`, the generation to be current, the provider state to be `handshaking`, and no terminal intent. The caller cannot supply a `helloResult`. `incompatible`, `stopping`, exited/crashed/backoff, replacement, and stale-generation attempts mutate nothing and publish no capabilities.
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

- A typed load-failure transition to `unavailable` on a live handshaken helper retains the snapshot.
- A no-instance discovery/factory-failure transition to `unavailable` invalidates it.
- `observeExit` during `unavailable` invalidates it and classifies the exit under the private termination intent.
- Tests construct both states separately; a single `unavailable` fixture is not sufficient.

### 3.5 Why a client copy is insufficient

`c.hello` is copied at handshake and is not tied to the current process. `Start` can clear `started` and `transport` after a failed handshake without clearing the old `hello`; `observeExit` does not clear supervisor capabilities. A client-only check can therefore authorize `c.call`, which starts a replacement. `callOnCurrent` and the generation check prevent that.

## 4. Direct operations

### 4.1 Direct `Unload`

1. Read the authoritative snapshot. If invalid, return `unsupported_capability` with `capability=howl.tts.unload`, `snapshot=unknown`; no frame; no `Start`/spawn/restart; no `unloaded` claim.
2. If valid and absent, return `unsupported_capability` with `capability=howl.tts.unload`, `snapshot=known`, `advertised=false`; no frame; no start; no state change.
3. If valid and advertised, in live asset-specific `unavailable`: use `callOnCurrent`, never `Start`; send at most one `unload` frame on the current instance; validate the `{"unloaded":true}` result. A successful helper response means the helper is unloaded, **not** that the aggregate failed attempt recovered. Do not submit `unloadSucceeded`, do not clear `LastError`, do not change the snapshot, and do not claim aggregate `unloaded`; the aggregate remains `unavailable`. A wire error is returned unchanged and changes no aggregate state. If `unload` is not advertised, return `unsupported_capability` with `snapshot=known` and no frame.
4. If valid and advertised in a normal live state, use `callOnCurrent`; no `Start`; no replacement. A generation/transport change yields the existing typed failure.
5. After a successful `unloaded=true` response in a normal state, re-check the serialization boundary: set `TTSStateUnloaded` only if the supervisor is still in a normal live state. If `unavailable`, `incompatible`, `stopping`, or an invalid snapshot won meanwhile, preserve that state; never clear a failed attempt or claim `unloaded`.

`Unload` is not idle shutdown. It MUST NOT send `shutdown`, terminate the process, or be treated as satisfying A2-D8. A2-D8 remains C4/C7 work. Unload in live `unavailable` is not recovery and cannot be the A2-D12 explicit `load` bypass.

### 4.2 Direct `Cancel`

1. Validate the target ID first; malformed target returns `invalid_request` before the snapshot gate.
2. If the snapshot is invalid, return `unsupported_capability` with `capability=howl.tts.cancel`, `snapshot=unknown`; no frame; no start/spawn/restart; no cancellation claim.
3. If valid and absent, return `unsupported_capability` with `capability=howl.tts.cancel`, `snapshot=known`, `advertised=false`; no frame; no start.
4. If valid and advertised, in live asset-specific `unavailable`, there is no active synthesis. Return the contract's exact non-active-target failure `request_not_found` with no frame, no `Start`, no state change, and no cancellation claim. This is not `provider_not_loaded`: that code is synthesis-specific. The precise `Failure.Code=request_not_found` maps through the existing `audioCodeFor` to `AudioError.Code=ErrCodeSynthesisFailed`, with the same `Failure` preserved in the cause. The snapshot and `LastError` remain intact.
5. If valid and advertised in a normal live state, use `callOnCurrent`; no `Start`; a generation/transport change yields the existing typed failure.
6. The direct method does not own the synthesis terminal and does not set `unloaded` or claim `cancelled`.

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

The precise `Failure.Code`, `Message`, `Field`, `Retryable`, and bounded `Detail` remain available through the cause/`AsFailure`; they are not collapsed into the public code. No wire message, new public error code, Wails binding, or new taxonomy is added. The known-absent branch uses the same existing surface with `snapshot=known`. The live-`unavailable` cancel result uses the existing contract code `request_not_found` and the same `AudioError` projection; it is not a D10 unknown/absent branch.

### 4.4 No-start current-instance send

`callOnCurrent(ctx, op, params, timeout, expectedGeneration)`:

- verifies the captured transport still belongs to `expectedGeneration`;
- returns an existing typed failure if the transport is absent or has changed;
- never calls `Start` or `Supervisor.Start`;
- writes at most one operation frame after the check;
- preserves request IDs, timeout, and cancellation behavior.

### 4.5 A2-D12 no-recovery gate for non-load paths

A2-D12 makes the failed-attempt state sticky: once aggregate `State=unavailable`, only a new explicit same-identity `load` on a live helper whose helper-side state is `unloaded` may recover it. No non-load path may start, handshake, ensure-load, or load. This applies to live and exited/no-instance `unavailable` alike. The explicit `Load` path is the sole permitted recovery bypass. C1 does not implement load semantics; C3/C6 must enforce live-helper, same `provider_id`/`engine`/`model_identity`, helper-side `unloaded`/`accepting_work=true`, and the D1a/D1b gates. No ordinary path may call `Load`, ensure-load, or a process start while `unavailable`. Re-check the state after serialization because an exit or failure can race the entry gate.

`provider_not_loaded` is used only for `Synthesize` in `State=unavailable`; it is defined by contract §4.6 as synthesis before successful load and is the code A2-D12 assigns to `synthesize`. It is not a generic rejection for `Health`, `Catalog`, `Unload`, or `Cancel`. `request_not_found` is the exact §17.4.6 non-active-target result for `Cancel`; `unsupported_capability` remains scoped to an absent/unknown optional capability. No new public error code or taxonomy is introduced.

### 4.6 Complete operation/error projection matrix

This matrix is the single source of truth for C1. The state classes are:

- **A — ordinary `absent`:** no helper process, transport, or authoritative snapshot; not a failed attempt.
- **K0 — live known-absent optional:** compatible handshake, current live instance, snapshot known; the named optional capability is absent.
- **K1 — live known-present optional:** compatible handshake, current live instance, snapshot known; the named optional capability is present.
- **L — live asset-specific `unavailable`:** compatible live handshake retained; typed load failure; aggregate `State=unavailable`.
- **X — exited/no-instance `unavailable`:** failed attempt with no live handshaken helper; snapshot invalid.
- **U — other unknown/no-current-handshake lifecycle:** `starting`, `handshaking`, `stopping` with no live handshake for the named capability, `exited`, `crashed`, `backoff`, or `incompatible`; no authoritative snapshot. `stopping` retains its live wire-health exception in §8.2.

“No frame” means no operation frame. “No start” means no `Client.Start`, `Supervisor.Start`, spawn, handshake, replacement, or recovery. “No state” means no aggregate, identity, snapshot, or `LastError` mutation attributable to the operation. Health is a mandatory control operation, not a service-capability-gated operation.

#### A — ordinary `absent`

| Operation | Helper / transport | Frame | Start / handshake / load | Result / `Failure` | Caller-facing `AudioError` | Aggregate / `LastError` | Owner |
|---|---|---|---|---|---|---|---|
| `Health` | None | No wire frame | No start | Supervisor observation; no helper health result | None | Remains `absent` | C1 document / C7 projection |
| `Catalog` (direct) | None | No frame | No start | `unsupported_capability`, `capability=howl.tts.catalog`, `snapshot=unknown` | `ErrCodeSynthesisFailed`; cause preserves failure | No state change | C1 |
| `Catalog` (`ProviderCatalog`) | None | Allowed after start/handshake | May start/handshake from ordinary `absent`; no load | Catalog result after handshake, or start/catalog error | Existing helper failure mapping | C7-owned start transition; no load | C7 adapter |
| `Unload` | None | No frame | No start | `unsupported_capability`, `capability=howl.tts.unload`, `snapshot=unknown` | `ErrCodeSynthesisFailed`; cause preserves failure | No state change | C1 |
| `Cancel` | None | No frame | No start | Validate target first; then `unsupported_capability`, `capability=howl.tts.cancel`, `snapshot=unknown` | `ErrCodeSynthesisFailed`; cause preserves failure | No state change | C1 |
| `Synthesize` | None at entry | No C1 frame | C3/C7 ordinary explicit-TTS ensure path may start/handshake/ensure-load from `absent`; C1 adds no start or recovery | Delegated C3/C7 result | Existing mapped result | Delegated C3/C7 | C3/C7 |
| `ProviderStatus` | None | No wire frame | May start/handshake from ordinary `absent`; no load | Supervisor status after start, or start/status error | Existing mapped result | C7-owned start transition; no load | C7 adapter |
| Voice-list | None | No frame | No start | Static/last snapshot | None | No state change | C7/C1 static |

#### K0 — live known-absent optional

| Operation | Helper / transport | Frame | Start / handshake / load | Result / `Failure` | Caller-facing `AudioError` | Aggregate / `LastError` | Owner |
|---|---|---|---|---|---|---|---|
| `Health` | Current live transport | One `health` frame | No start/handshake/load | Validated health result | None | No state change; `LastError` retained | C1 |
| `Catalog` | Current live transport | No frame when `catalog` absent | No start | `unsupported_capability`, `capability=howl.tts.catalog`, `snapshot=known` | `ErrCodeSynthesisFailed`; cause preserves failure | No state change | C1 |
| `Unload` | Current live transport | No frame when `unload` absent | No start | `unsupported_capability`, `capability=howl.tts.unload`, `snapshot=known` | `ErrCodeSynthesisFailed`; cause preserves failure | No state change | C1 |
| `Cancel` | Current live transport | No frame when `cancel` absent | No start | `unsupported_capability`, `capability=howl.tts.cancel`, `snapshot=known` | `ErrCodeSynthesisFailed`; cause preserves failure | No state change | C1 |
| `Synthesize` | Current live transport | Determined by normal C3 state | No process start; C3 may ensure-load from ordinary `unloaded` only | Normal C3 synthesis/load result | Existing mapping | C3 normal transitions only | C3 |
| `ProviderStatus` | Current live transport | No wire frame | No start/load | Supervisor/current status | None | No state change | C7/C2 |
| Voice-list | Static/last snapshot | No frame | No start | Existing list | None | No state change | C7/C1 static |

#### K1 — live known-present optional

| Operation | Helper / transport | Frame | Start / handshake / load | Result / `Failure` | Caller-facing `AudioError` | Aggregate / `LastError` | Owner |
|---|---|---|---|---|---|---|---|
| `Health` | Current live transport | One `health` frame | No start/handshake/load | Validated health result | None | No state change | C1 |
| `Catalog` (direct) | Current live transport | One `catalog` frame | No start/handshake/load | Validated catalog result | None on success; helper failure preserved | No state or identity change; no `SetSelected` | C1 |
| `Unload` (normal state) | Current live transport | One `unload` frame | No start; no shutdown/termination | `{"unloaded":true}` validated | None on success; helper failure preserved | Normal live state may become `unloaded`; no terminal write | C1/C3 |
| `Cancel` (normal state) | Current live transport | One `cancel` frame | No start | Ack/error; `request_not_found` if target not active | `ErrCodeSynthesisFailed` for wire failure; success returns nil | No aggregate terminal write | C1 |
| `Synthesize` | Current live transport | Normal C3 frames | No process start; C3 may ensure-load only from ordinary `unloaded` | Normal C3 synthesis/load result | Existing mapping | C3 normal transitions | C3 |
| `ProviderStatus` | Current live transport | No wire frame | No start/load | Supervisor/current status | None | No state change | C7/C2 |
| Voice-list | Static/last snapshot | No frame | No start | Existing list | None | No state change | C7/C1 static |

#### L — live asset-specific `unavailable`

| Operation | Helper / transport | Frame | Start / handshake / load | Result / `Failure` | Caller-facing `AudioError` | Aggregate / `LastError` | Owner |
|---|---|---|---|---|---|---|---|
| `Health` | Current live transport | One `health` frame | No start/handshake/load/recovery | Validated health result | None | Remains `unavailable`; `LastError` retained | C1 |
| `Catalog` | Current live transport | One `catalog` frame only if advertised | No start/handshake/load/recovery; no `SetSelected` | Validated catalog result; absent capability is `unsupported_capability` `snapshot=known` | None on success; `ErrCodeSynthesisFailed` for absent/failure | Remains `unavailable`; identity and `LastError` retained | C1 |
| `Unload` | Current live transport | One `unload` frame only if advertised | No start; no shutdown/termination | `{"unloaded":true}` from helper | None on success; helper failure preserved | Remains `unavailable`; no aggregate `unloaded`; `LastError` retained | C1 |
| `Cancel` | Current live transport | No frame | No start | `request_not_found`; no synthesis is active | `ErrCodeSynthesisFailed`; cause preserves failure | Remains `unavailable`; no cancellation claim | C1 |
| `Synthesize` | Current live transport | No frame | No start/handshake/load/recovery | `provider_not_loaded` | `ErrCodeModelLoadFailed`; cause preserves failure | Remains `unavailable`; `LastError` retained | C1 |
| `ProviderStatus` | Current live transport | No wire frame | No start/load | Supervisor status including `State=unavailable` | None | No transition | C7/C2 |
| Voice-list | Static/last snapshot | No frame | No start | Existing list | None | No transition | C7/C1 static |

#### X — exited/no-instance `unavailable`

| Operation | Helper / transport | Frame | Start / handshake / load | Result / `Failure` | Caller-facing `AudioError` | Aggregate / `LastError` | Owner |
|---|---|---|---|---|---|---|---|
| `Health` | None | No wire frame | No start/recovery | Supervisor observation; no helper health result | None | Remains `unavailable`; no transition | C1 document / C7 projection |
| `Catalog` (direct) | None | No frame | No start/recovery | `unsupported_capability`, `capability=howl.tts.catalog`, `snapshot=unknown` | `ErrCodeSynthesisFailed`; cause preserved | No state change | C1 |
| `Catalog` (`ProviderCatalog`) | None | No frame | No start/recovery | Same fail-closed result; no start from `unavailable` | `ErrCodeSynthesisFailed`; cause preserved | No state change | C7 adapter |
| `Unload` | None | No frame | No start/recovery | `unsupported_capability`, `capability=howl.tts.unload`, `snapshot=unknown` | `ErrCodeSynthesisFailed`; cause preserved | No state change | C1 |
| `Cancel` | None | No frame | No start/recovery | Target validated; `unsupported_capability`, `capability=howl.tts.cancel`, `snapshot=unknown` | `ErrCodeSynthesisFailed`; cause preserved | No state change | C1 |
| `Synthesize` | None | No frame | No start/recovery | `provider_not_loaded` under A2-D12 | `ErrCodeModelLoadFailed`; cause preserved | Remains `unavailable`; `LastError` retained | C1 guard; C3/C5/C6 recovery |
| `ProviderStatus` | None | No wire frame | No start/recovery | Supervisor status including `State=unavailable` | None | No transition | C7/C2 |
| Voice-list | None | No frame | No start/recovery | Static/last snapshot | None | No state change | C7/C1 static |

#### U — other unknown/no-current-handshake lifecycle

`U` excludes ordinary `absent` (A), live `unavailable` (L), and exited/no-instance `unavailable` (X). It includes `starting`, `handshaking`, `stopping` when no current handshaken transport is usable, `exited`, `crashed`, `backoff`, and `incompatible`.

| Operation | Helper / transport | Frame | Start / handshake / load | Result / `Failure` | Caller-facing `AudioError` | Aggregate / `LastError` | Owner |
|---|---|---|---|---|---|---|---|
| `Health` | None for `starting`/`handshaking`/`exited`/`crashed`/`backoff`/`incompatible`; current transport for `stopping` | Supervisor observation; `stopping` uses wire `health` while alive; `handshaking`/`incompatible` wire-health excluded | No start/handshake/retry | Supervisor/status or validated health result | None | No transition | C1/C5/C7 |
| `Catalog` (direct) | None or current invalid instance | No frame | No start/recovery | `unsupported_capability`, `capability=howl.tts.catalog`, `snapshot=unknown` | `ErrCodeSynthesisFailed`; cause preserved | No state change | C1 |
| `Catalog` (`ProviderCatalog`) | None or current invalid instance | No frame | No start from `starting`/`handshaking`/`stopping`/`exited`/`crashed`/`backoff`/`incompatible` | Fail closed; no recovery | `ErrCodeSynthesisFailed`; cause preserved | No state change | C7 adapter |
| `Unload` | None | No frame | No start/recovery | `unsupported_capability`, `capability=howl.tts.unload`, `snapshot=unknown` | `ErrCodeSynthesisFailed`; cause preserved | No state change | C1 |
| `Cancel` | None | No frame | No start/recovery | Target validated; `unsupported_capability`, `capability=howl.tts.cancel`, `snapshot=unknown` | `ErrCodeSynthesisFailed`; cause preserved | No state change | C1 |
| `Synthesize` | State-dependent | No C1 frame | No C1 start; C5/C6/C7 own the explicit trigger and budget | `incompatible` returns `protocol_incompatible` through the start gate; `crashed`/`backoff`/`stopping`/`starting`/`handshaking` have no C1 start; `unavailable` is X | Existing mapped result | No C1 transition | C5/C6/C7, except C1 gate |
| `ProviderStatus` | None | No wire frame | No start from U; status observation only | Supervisor status | None | No transition | C7/C2 |
| Voice-list | None | No frame | No start | Static/last snapshot | None | No state change | C7/C1 static |

### 4.7 Direct client versus C7 adapter start authority

This table resolves C1-R-22. It is normative for the C1/C7 boundary.

| Path | `absent` | `starting`/`handshaking`/`stopping` | `crashed`/`backoff`/`incompatible` | Live `unavailable` | Exited/no-instance `unavailable` |
|---|---|---|---|---|---|
| **Direct `Client.Catalog`** | No start; fail closed `unsupported_capability`/unknown | No start; fail closed | No start; fail closed | Current-transport `catalog` if advertised; no recovery | Fail closed; no start |
| **C7 `ProviderCatalog`** | May start/handshake, then `catalog`; no load | No start; fail closed | No start; fail closed | Current-transport `catalog` if advertised; no recovery | No start; fail closed |
| **Direct `Client.Health`** | Supervisor observation; no start | Supervisor observation; `stopping` may use wire health; no start | Supervisor observation; no start | Wire health; no start | Supervisor observation; no start |
| **C7 `ProviderStatus`** | May start/handshake; no load | No start; status observation | No start; status observation | Status observation; no recovery | Status observation; no recovery |
| **C1 `Client.Synthesize`** | C3/C7 ordinary explicit-TTS path; C1 adds no start/recovery | No C1 start; C3/C5/C6 govern | `incompatible` start gate; C5/C6 govern the rest | `provider_not_loaded`; no start | `provider_not_loaded`; C3/C5/C6 recovery |

`ProviderCatalog` is a C7 adapter surface, not a second direct optional operation. It may invoke `Client.Start` only from ordinary `absent`; it never starts from `unavailable`, `crashed`, `backoff`, `stopping`, `starting`, `handshaking`, or `incompatible`. Direct `Client.Catalog` never calls `Start`. This is the only reading of C1-AC-12 and contract §17.12.1 that does not both forbid and require the same start. The C7 adapter must be tested separately; C1 does not claim C7 executed start evidence.

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

The intent is keyed to the captured process instance and generation. `Terminate` owns process lifecycle only; it may set `ProcessState=stopping` but MUST NOT overwrite aggregate `State=incompatible` with `stopping`. The aggregate projection is a separate field and remains `incompatible` until the expected exit is observed.

`observeExit` consumes that intent only for the tracked incompatible instance:

**Case A — termination succeeds.** The exit is expected. `ProcessState=exited`, `ProcessAlive=false`, and aggregate `State` remains `incompatible`. The observer MUST NOT write `crashed`, enter `backoff`, spend restart budget, or trigger replacement. The single `Wait()` owner observes the exit once.

**Case B — termination times out while the process remains live.** `State` remains `incompatible`; the live-but-unreaped process remains tracked by the supervisor and the single `Wait()` owner. `ProcessState` may remain `alive` or `stopping` while `ProcessAlive=true`, as A2-D4 permits. No second start, handshake, transport reuse, frame, backoff, restart-budget spend, replacement, or `crashed` classification is permitted. Repeated observation of the same instance does not change the terminal classification. A later exit is still consumed under the latched intent and also remains `incompatible`.

**Normal-start gate for both cases.** `Client.Start` and `Supervisor.Start` MUST check aggregate `incompatible` before doing any of the following: returning an existing process, creating a process, starting a reader, reusing a transport, or sending `hello`. If the state is `incompatible`, both return the existing `protocol_incompatible` failure and leave the process/observer untouched. This includes the live-but-unreaped Case B process. A normal TTS request, a `Catalog`, a `Health` call, or a direct `Client.Start` therefore cannot obtain a second handshake. Only the update-gated return to `absent` under C6/C7 clears `incompatible`.

Normal shutdown classification is separate: exit 0 → `absent`; nonzero exit or grace expiry → `crashed`. The incompatible intent is the explicit exception required by A2-D4.

### 6.5 Handshake provenance and aggregate-state ownership

**R17 publication is process-owned, not value-owned.** The only handshake publication entry point is a private `handshakeAttempt` returned by `Supervisor.beginHandshake(p)`:

```text
handshakeAttempt {
    p          Process
    generation uint64
    hello      helloResult   // set only by the pinned exchange
    validated  bool
    completed  bool
}
```

Under `Supervisor.mu`, `beginHandshake` requires a current `p == s.proc`, the current generation, `providerState=handshaking`, no terminal intent, and no replacement. It creates exactly one attempt and stores it as `s.handshakeAttempt`.

The attempt's `exchange(ctx, c)` is the only method that sends `hello` and receives its result on `p`. It uses a **pinned** exchange (`c.callOnProcess(p, ...)`), validates the result with `validateHello`, and stores it in the attempt. The caller cannot supply a `helloResult` to publication. `attempt.publish(selected)` then requires all of:

1. `attempt == s.handshakeAttempt`, `completed == false`, `p != nil`, and `s.proc == p`;
2. `s.generation == attempt.generation` and no replacement has begun;
3. `s.providerState == handshaking`;
4. `s.processState` is `starting` or the just-started alive state, not `stopping`, `exited`, `crashed`, `incompatible`, `unavailable`, or `backoff`;
5. no `terminationIntent` for this instance;
6. `attempt.validated` is true and `validateHello(attempt.hello)` still succeeds;
7. `selected` is nil or exactly one of `attempt.hello.Providers` and its `ProviderID`/`Engine` match.

Publication uses only `attempt.hello`; there is no `h` parameter. It marks the attempt completed and publishes versions, capabilities, selected identity, `processState=alive`, and `providerState=unloaded`. A failure writes nothing. `incompatible`, `stopping`, exited/crashed/backoff, replacement, stale generation, a foreign attempt, and a mismatched selected provider are explicit refusals. The old exported `MarkHandshake` is removed; it is not renamed into a second entry point. `Client.Start` must not retain or call any other publication method. This closes the value-provenance gap independently of the start gates.

**R18 aggregate-state ownership is single-writer.** `Supervisor` owns private `processState`, `providerState`, `lastFailure`, and `activeRequestID`. No exported method writes them. A private `apply(event)` transition under `Supervisor.mu` is the only mutation path; construction initialization is the only exception. The currently exported `SetState` is removed. The currently exported `SetSelected` is removed; identity is set by `attempt.publish` or `loadSucceeded`, not by `Catalog`.

The event set is deliberately small and closed:

```text
beginStart, processSpawned, handshakeSucceeded, handshakeFailed,
beginLoad, loadSucceeded, loadFailed, beginSynthesis, synthesisSucceeded,
synthesisFailed, unloadSucceeded, beginTermination, observeExit,
noInstanceUnavailable, updateGatedReset
```

`beginTermination` owns process lifecycle only. It may set `processState=stopping`; it sets `providerState=stopping` only for a normal termination from a normal live state. It MUST NOT overwrite `incompatible`, and it MUST NOT clear a live `unavailable` failure. The exit observer is the only writer of `processState=exited`/`absent`, and the only writer of `providerState=crashed` from an unexpected exit. A transport send failure requests termination; it does not write `crashed` directly. This removes the current `client.go:756` direct-crash write.

`incompatible` is terminal for ordinary events. The only allowed events are `observeExit` for the tracked process and `updateGatedReset` with a valid C6/C7 update token. `beginStart`, `processSpawned`, `handshakeSucceeded`, `beginLoad`, `loadSucceeded`, `loadFailed`, `beginSynthesis`, `synthesisSucceeded`, `synthesisFailed`, `unloadSucceeded`, `handshakeFailed`, and `noInstanceUnavailable` all reject without mutation. `beginTermination` may set process state only.

A live `unavailable` rejects `handshakeSucceeded`, `beginSynthesis`, `synthesisSucceeded`, `synthesisFailed`, `unloadSucceeded`, `beginStart`, and `processSpawned`. It permits `beginLoad` and `loadFailed` only through the C6/C7 explicit same-identity recovery path, `beginTermination` for process lifecycle, `observeExit`, and `updateGatedReset` with a valid update token. `unavailable` with no live instance permits `updateGatedReset` only; it does not permit a direct start or a general setter.

`updateGatedReset` is not a general-purpose setter and is not called by C1. It requires `(token, updateEpoch)` from the C6/C7 helper/assets/config update path, a terminal state `incompatible` or `unavailable`, and no live process. It rejects a reused/unknown token, a live process, or a missing epoch. On success it transitions to `absent`, clears the terminal failure, invalidates the snapshot, and advances the generation so a stale attempt cannot publish. C1 defines this guard shape; C6/C7 must supply the actual update evidence. No caller can impersonate an update with an arbitrary state string.

`LastError` is cleared only by an authorized success or `updateGatedReset`: `handshakeSucceeded`, `loadSucceeded`, `synthesisSucceeded`, and a normal `unloadSucceeded` may clear it; `loadFailed`, `synthesisFailed`, `handshakeFailed`, `noInstanceUnavailable`, and an unexpected `observeExit` set it; `beginLoad` and `beginSynthesis` do not clear a terminal failure. A rejected event writes no `LastError`. Rejection leaves the original terminal failure and the actual process state intact.

### 6.6 Aggregate-state writer and ownership inventory

This inventory covers every direct assignment to `providerState`, `lastFailure`, or `processState` in `internal/audio/tts/helper` at the revision-6 source baseline. A repository search found no writer outside `supervisor.go`; `Client` writes only through `Supervisor` methods. Line references are current-source evidence, not implementation claims.

| Current writer | Field(s) written | Current trigger | Current effect | Revision-6 owner/event |
|---|---|---|---|---|
| `NewSupervisor`, `supervisor.go:103-104` | `processState`, `providerState` | Construction | `absent`/`absent` | Constructor initialization only |
| `Supervisor.Start`, `supervisor.go:136-137` | `processState`, `providerState` | Discovery resolved | `starting`/`starting` | `beginStart` |
| `Supervisor.Start`, `supervisor.go:158-159` | `processState`, `providerState` | Factory succeeded | `starting`/`handshaking` | `processSpawned` + `beginHandshake` |
| `Supervisor.observeExit`, `supervisor.go:176-191` | `processState`, `providerState`, `lastFailure` | Single `Wait()` exit | Normal/absent or crash/failure | `observeExit`, sole exit writer |
| `Supervisor.setUnavailable`, `supervisor.go:199-201` | `processState`, `providerState`, `lastFailure` | Discovery/factory failure | `absent`/`unavailable` | `noInstanceUnavailable` |
| `Supervisor.MarkHandshake`, `supervisor.go:205-219` | `processState`, `providerState`, identity/caps | Direct handshake publication | `alive`/`unloaded` | Removed; `handshakeAttempt.publish` |
| `Supervisor.SetState`, `supervisor.go:233-242` | `providerState`, `lastFailure`, `processState` | Caller-supplied state | Arbitrary | Removed; closed `apply(event)` |
| `Supervisor.Terminate`, `supervisor.go:255-259` | `processState`, `providerState` | Termination request | `stopping`/`stopping` | `beginTermination`; never terminal overwrite |
| `Client.Start` failure, `client.go:244-253` | `providerState`, `lastFailure` | Failed handshake | `incompatible`/`unavailable` | `handshakeFailed`, then intent-aware termination |
| `Client.SetSelected` after catalog, `client.go:546` | identity only | Catalog result | Provider identity | Removed; `Catalog` is read-only |
| `Client.Load` begin, `client.go:578` | `providerState` | Load begins | `loading` | `beginLoad` |
| `Client.Load` failure, `client.go:581,586,590` | `providerState`, `lastFailure` | Typed load failure | `unavailable` | `loadFailed` |
| `Client.Load` success, `client.go:593-594` | identity only; `providerState` | Load result | Identity + `ready` | `loadSucceeded` |
| `Client.Unload` success, `client.go:616` | `providerState` | Normal unload | `unloaded` | `unloadSucceeded`; rejected from `unavailable`/terminal |
| `Client.Synthesize` begin, `client.go:737` | `providerState` | Synthesis begins | `busy` | `beginSynthesis` |
| `Client.Synthesize` terminal, `client.go:748,771,783,793,796,803,806,814,820,828,839,847` | `providerState`, `lastFailure` | Successful/failed/cancelled/timeout terminal | `ready`/`degraded` | `synthesisSucceeded`/`synthesisFailed` |
| `Client.Synthesize` write failure, `client.go:756` | `providerState`, `lastFailure` | Transport write failed | Direct `crashed` | Removed; request termination; `observeExit` classifies |
| `Supervisor.SetSelected`, `supervisor.go:223-230` | identity only | Catalog/load | Provider identity | Removed; `loadSucceeded` or `attempt.publish` only |

No `providerState`, `lastFailure`, or `processState` write in `Client` remains after the migration. The test inventory must also include test-only fixtures, but tests may not add a production writer or a general setter.

### 6.7 Allowed-event and terminal-state matrix

| Event | Allowed source state(s) | Required writes | Required refusal |
|---|---|---|---|
| `beginStart` | `absent` only | New generation; `processState=starting`; `providerState=starting`; no terminal clear | `incompatible`, `unavailable`, `stopping`, `starting`, `handshaking`, `crashed`, `backoff` |
| `processSpawned` | `starting` for the captured start | `processState=starting`; `providerState=handshaking` | No current start, wrong instance/generation, terminal intent |
| `handshakeSucceeded` | Current attempt, `handshaking`, captured `p`/generation | `processState=alive`; `providerState=unloaded`; caps/identity/versions; clear authorized failure | Terminal/no attempt/stale/foreign attempt/invalid hello/mismatched selected |
| `handshakeFailed` | Current attempt, `handshaking` | Terminal `incompatible` for semantic hello failure, otherwise `unavailable`; set exact `Failure`; latch private intent | Terminal/stale/foreign attempt; no state write on refusal |
| `beginLoad` | `unloaded`; live `unavailable` only through C6/C7 token | `providerState=loading`; active request empty; preserve previous `LastError` | `incompatible`, `stopping`, `exited`, `crashed`, `backoff`, no live helper, different identity |
| `loadSucceeded` | `loading` current live attempt | `ready`; selected identity; clear failure | Stale generation, exited process, terminal/stopping, identity mismatch |
| `loadFailed` | `loading` current live attempt | `unavailable`; exact failure; retain live snapshot if helper live | Exited process, stale/terminal, unrelated failure |
| `beginSynthesis` | `ready`/`degraded`, current live | `busy`; active ID | `unavailable`, `incompatible`, `stopping`, `crashed`, `backoff`, no live |
| `synthesisSucceeded` | `busy`, current request | `ready`; clear failure; clear active ID | Stale request, exit, terminal, timeout/cancel winner |
| `synthesisFailed` | `busy`, current request | `degraded` for an ordinary non-fatal failure; preserve failure | Stale request, terminal, exit; no `crashed` write |
| `unloadSucceeded` | `ready`/`degraded`/`unloaded` normal live | `unloaded`; clear active ID | `unavailable`, `incompatible`, `stopping`, exited, stale |
| `beginTermination` | Any tracked live process | `processState=stopping`; `providerState=stopping` only for normal intent; preserve `incompatible`/`unavailable` | No process; terminal intent already consumed; stale instance |
| `observeExit` | Single `Wait()` owner for tracked `p` | `processState=absent`/`exited`; provider `absent`/`crashed`/`incompatible`/`unavailable` per intent; failure only for unexpected crash | Any other process; second observer; post-reset instance |
| `noInstanceUnavailable` | No live process; discovery/factory failure | `processState=absent`; `providerState=unavailable`; exact failure; invalidate snapshot | Live process, `incompatible`, `stopping`, terminal intent |
| `updateGatedReset` | `incompatible`/`unavailable` with no live process and valid C6/C7 token/epoch | `absent`; clear failure; invalidate snapshot; advance generation | Live process, missing/reused token, no update evidence, ordinary caller |

Any event not listed, or any source state not listed for that event, is rejected with the existing typed failure and **no mutation**. In particular, no event other than `observeExit` and `updateGatedReset` may alter `incompatible`; no event other than `beginLoad`/`loadFailed`/`beginTermination`/`observeExit`/`updateGatedReset` may be accepted while live `unavailable`; and no event writes a terminal failure to `nil` unless it is an authorized success.

**R18 termination rule.** `beginTermination` may set `processState=stopping`; it leaves `providerState=incompatible` unchanged. `observeExit` for that instance sets `processState=exited`, `ProcessAlive=false`, and leaves `incompatible`; it never writes `crashed`, `backoff`, `absent`, or `stopping` for that exit. A later exit after a timeout is consumed by the same single `Wait()` owner under the latched intent and remains `incompatible`. No exit from this path spends restart budget, starts a process, or sends `hello`. A barrier-based test must observe `Status()` after termination has been requested but before the process exits: `State=incompatible`, `ProcessState` in `alive`/`stopping`, `ProcessAlive=true`; after exit: `State=incompatible`, `ProcessState=exited`, `ProcessAlive=false`. The test must not wait only for termination success or timeout.

## 7. C1/C2 boundary

C1 adds only private snapshot/generation state, the validity predicate, the no-start direct operation path, the A2-D12 no-recovery guard, the process-owned handshake attempt, the single-writer aggregate-state machine, the normal-start `incompatible` gate, intent-aware incompatible termination, and the path-scoped cancellation operation. It adds no public status field and no load semantics. It removes the exported `SetState`, `MarkHandshake`, and `SetSelected` surfaces rather than preserving general-purpose mutation.

C2 owns `Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown`, and public `TTSProviderStatus` projection. C2 must not claim C1's lifecycle invalidation or invent a second validity rule. C3 owns ordinary load/ensure-load and is the only place that may call the explicit-load recovery path; it must enforce A2-D12's live-helper and identity conditions. C4–C6 own the full lifecycle transitions; C1's generation invalidation and state gates are their prerequisite. C7 owns the production `ProviderCatalog`/`ProviderStatus` adapter start from ordinary `absent`; C1 owns the direct `Client.Catalog` no-start rule and the live `Health` read. Neither C1 nor C7 may start from `unavailable`, `crashed`, `backoff`, `stopping`, `starting`, `handshaking`, or `incompatible`. `Synthesize`, direct `Unload`, direct `Cancel`, and status/voice-list do not enter the explicit-load recovery path.

## 8. Complete C1 start/restart map

### 8.1 State and trigger map

This is a documented map. C1 owns the map and structural no-spawn checks; C5 owns executed restart/backoff; C7 owns executed start/reachability.

| State | Start/restart trigger | Direct `Unload` / `Cancel` | Snapshot | Required classification |
|---|---|---|---|---|
| `absent` | Explicit TTS/synthesis, C7 `ProviderCatalog`, or C7 `ProviderStatus` may start/handshake. Direct `Client.Catalog`, direct `Health`, `GetTTSVoiceCatalog`, `Unload`, and `Cancel` never start. | No direct optional start | Unknown/cleared | No direct optional spawn; C7 adapter owns the absent start. |
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
| `incompatible` | **No `Client.Start` or `Supervisor.Start`; no second handshake or transport reuse.** Only update-gated return to `absent` under C6/C7. A live-but-unreaped process remains tracked until its single `Wait()` exit. | No | Invalid | Expected termination/late exit remains `incompatible`; no `crashed`, backoff, budget, replacement, or start. |
| `unavailable` — live asset-specific | No synthesis/voice-list/status recovery. `Health` uses current-transport wire health; advertised `Catalog` uses current-transport catalog; only explicit same-identity `load` on a live helper whose helper-side state is `unloaded`, subject to D1a/D1b, recovers. C7 `ProviderCatalog`/`ProviderStatus` do not start or recover. | Advertised `Unload` may use one current-transport frame without an aggregate `unloaded` transition; `Cancel` returns `request_not_found`; absent capabilities fail without a frame. | **Retained** | Health stays `unavailable`; catalog/unload/cancel cause no recovery/load/start. |
| `unavailable` — exited/no instance | No start; only the applicable explicit recovery path after C6/C7. Direct `Catalog`/`Health` and C7 `ProviderCatalog`/`ProviderStatus` do not start or recover. | No | **Invalidated** | Health is supervisor-level; no start/restart/recovery; `Synthesize` returns `provider_not_loaded` under A2-D12. |
| `exited` | No start by C1. C5/C6/C7 own the next permitted explicit trigger. | No | Invalid | Exit 0 → `absent`; nonzero/grace expiry → `crashed`, except the A2-D4 incompatible intent. |

No start is permitted before `absent` or `crashed` is observed, or from any other live state. A `load` transition on an already-live helper is not a process start. Direct optional operations never start, spawn, or restart. The `incompatible` row gates both normal start layers, not just process creation: a still-live incompatible process cannot be adopted for a second handshake or transport reuse.

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
| `incompatible` | **Wire-health exclusion** | No wire health accepted; no retry, restart, start/handshake, or state change. A live-but-unreaped Case B process remains tracked and cannot be adopted. |
| `unavailable` — live | Wire `health` | Remains `unavailable`; no recovery, load, start, cancellation, or transition; current-transport `catalog` and helper-side `unload` also cause no transition. |
| `unavailable` — exited/no instance | Supervisor | No start, restart, or recovery. |

If exit occurs during a health observation, the single `Wait()` observer remains authoritative; the exit is not attributed to health. `handshaking` and `incompatible` are the only wire-health exclusions, each grounded in the contract.

## 9. Concurrency, resource, and I/O evidence

### 9.1 Existing serialization and winner rules

C1 adds no new synchronization architecture. It uses the existing boundaries:

| Boundary | Role | C1 consequence |
|---|---|---|
| `Client.workMu` | Serializes `Load`, `Unload`, `Catalog`, and the entire `Synthesize` operation. | Concurrent `Unload` calls serialize. `Unload` issued during `loading` or `busy` waits and re-gates. Direct `Catalog`/`Health` do not enter the load/synthesis state machine. |
| `Client.startMu` | Serializes one client `Start`/handshake; `Supervisor.startMu` serializes process start. | No second start; direct optional operations never enter either path; C7 adapter starts only from ordinary `absent`. |
| `Client.mu` | Guards transport, pending requests, `hello`, and client state. | Transport/generation capture and pending registration are atomic; it is not held across `Send`. |
| `Supervisor.mu` | Guards the one `apply(event)` transition, snapshot publication/invalidation, process state, termination intent, and the current `handshakeAttempt`. | All aggregate writes are single-writer and atomic; publication, rejection, and invalidation cannot interleave; readers receive copies. |
| `handshakeAttempt.exchange` | Pinned to the captured `Process`; the supervisor holds the attempt identity, not `Supervisor.mu`, across I/O. | No second exchange; a foreign result cannot be supplied; publication re-checks the attempt under `Supervisor.mu`. |
| `Process.writeMu` | Serializes frame writes on one transport. | Concurrent direct `Cancel` frames cannot interleave bytes. |
| `Supervisor.Terminate` | Stops/kills and waits on the single observer. | Snapshot invalidation happens before close/wait; cancellation never waits on `workMu`. |

Lock discipline: do not hold `Client.mu` or `Supervisor.mu` across `Start`, `Send`, `Terminate`, or `Wait`; `workMu` may span a request wait because it is the request-serialization boundary, not a state lock. Lock order and any exception are documented in the C1 report. This is a design constraint on the existing code, not a new global lock.

Winner/terminal rules:

- **Unload vs Unload:** serialized by `workMu`; each call re-gates after acquiring it. Each valid call sends at most one frame. The helper's idempotent unload response governs; no `unloaded` state is claimed without a response.
- **Unload vs Load/Synthesize:** serialized by `workMu`; no interleaving. Cancellation precedence remains A2-D9/§17.11.6.
- **Cancel vs Cancel:** not serialized by `workMu`, so cancellation stays responsive. Each has its own request ID; `Process.writeMu` serializes writes. A late acknowledgement cannot reopen or overwrite the target's terminal; a helper terminal error is returned as that typed failure, not converted to success.
- **Direct optional vs start:** if the snapshot is not yet published, the direct call fails closed; it cannot enter `Start`. If a start completes after the gate, `callOnCurrent` uses the captured generation/transport and never starts. C7 `ProviderCatalog`/`ProviderStatus` are the only adapter paths that may start from ordinary `absent`; they are serialized by `Client.startMu`/`Supervisor.startMu` and do not use a direct optional path.
- **Normal start vs incompatible:** both `Client.Start` and `Supervisor.Start` re-check `incompatible` under their ownership boundary; neither returns an existing process nor starts/handshakes. Case A and Case B are the same gate; the live-but-unreaped process stays tracked by the single `Wait()` owner.
- **Writer vs writer:** all lifecycle events enter the same `apply` under `Supervisor.mu`; the last accepted event sees the current state and either transitions or refuses atomically. There is no state setter outside this boundary, so a stale writer cannot overwrite a terminal state.
- **Attempt publication vs termination/exit:** `beginHandshake`, `attempt.publish`, `beginTermination`, and `observeExit` all serialize on `Supervisor.mu`; publication after the intent is refused, and an accepted publication cannot erase a later exit's classification. An attempt that loses the current-attempt pointer cannot publish.
- **A2-D12 recovery vs non-load operations:** only `Synthesize` returns `provider_not_loaded` without a frame/start. `Health` and advertised `Catalog` may read the current transport; `Unload` may send without aggregate transition; `Cancel` returns `request_not_found`. Only explicit `Load` may enter the C3/C6 recovery bypass after serialization.
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
| Helper processes | 2 app-session cap; C1 optional operations start 0 | Contract §5.5 (line 520) hard cap, repeated in frozen §17 material; r9 §8 C7 stage overlay | No optional-operation spawn; separate from C1's no-start assertion | Process-count and start-counter test | 0 starts from C1 optional paths; global cap ≤ 2 | Pending independent reviewer |
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
| `Client.Start` / `callOnCurrent` / `handshakeAttempt.exchange` | Read and write NDJSON frames on inherited stdio | Handshake once per instance; one frame per valid direct operation | Capability negotiation and optional operation dispatch | 4,194,304-byte frame; ≤ 1 operation frame per call | Helper lifecycle | Frame counters, send-count tests, max-frame fixture | Not N/A; actual protocol I/O |
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
| HW-1.6 | Applicable as an explicit assumption record; no C1 hardware-dependent behavior | CPU cores, memory, disk, and GPU/accelerator availability are recorded in §9.5. C1 adds no native compute, model load, disk I/O, or GPU work; its bounded snapshot/gate logic is hardware-independent. The independent C1 reviewer must accept the assumptions; this is not a blanket N/A. |
| HW-1.7 | N/A | Reviewed N/A: no new contention signal or degradation behavior; existing bounded diagnostics unchanged. |
| HW-1.8 | Applicable | No throttle, sleep, delay, or concurrency reduction is added. |
| HW-2.1 | Applicable | No safety behavior is traded for performance. |
| HW-2.2 | N/A | Reviewed N/A: C1 adds no deliberate degradation path. |
| HW-2.3 | N/A | Reviewed N/A: C1 makes no performance claim and runs no performance measurement. |
| HW-2.4 / PR-1–PR-4 | Inactive under A-1; impact statement only | PR-3 N/A with review; no PR suite execution or acceptance result. |

### 9.5 HW-1.6 cross-hardware assumptions

C1 is a protocol/state-gating change, not a compute path. These assumptions are design constraints, not measurements; the independent C1 reviewer must confirm them and record any disagreement. None is a blanket N/A.

| Resource | C1 assumption | Bounded behavior / scaling | Review status |
|---|---|---|---|
| CPU cores | C1 adds no worker, native call, or inference and does not require a minimum core count. It must remain correct on a single available core. | Snapshot/gate work is O(1) over five IDs; low-core scheduling can add latency but cannot change state, ordering, or winner rules. C1 does not scale work with core count and adds no busy wait or spin. | Pending independent C1 reviewer; not N/A. |
| Memory | The snapshot is one generation plus at most five service entries; reads copy that bounded view. C1 adds no model, ORT, asset, or frame buffer. | O(5) snapshot storage and per-read copy; no allocation proportional to provider/asset count. The existing pending-request map remains bounded by `MaxPending=32`, whose authorized-source determination is open separately in §9.2. | Pending independent C1 reviewer; `MaxPending` source open. |
| Disk | C1 performs no disk read or write. Helper discovery/start and model preparation remain outside C1. | Zero C1 file opens, writes, temp files, or log files; low disk can affect the helper's unrelated startup but cannot make C1 retry, spawn, or load. | Pending independent C1 reviewer; no C1 disk dependency. |
| GPU / accelerator | C1 performs no ORT, model, GPU, or accelerator initialization and assumes no accelerator availability. | No GPU-dependent branch, fallback, or retry. A missing or weaker GPU cannot turn into a C1 retry loop or correctness change; C1 behavior is identical with or without an accelerator. | Pending independent C1 reviewer; not N/A, recorded as no dependency. |

C1 therefore does not scale correctness or resource use with CPU cores, memory, disk, or GPU/accelerator capacity. The only scaling is bounded control-plane latency under host scheduling; the existing process/request bounds remain the resource constraints.

## 10. Revised C1 acceptance criteria

| ID | Criterion |
|---|---|
| **C1-AC-1** | Mandatory **service** floor is exactly `howl.tts.load` + `howl.tts.synthesize`; `howl.tts.catalog`, `howl.tts.unload`, and `howl.tts.cancel` are optional service capabilities. `health` and `shutdown` are mandatory control operations, not optional service IDs and not gated by a service capability. Missing mandatory service IDs produce `protocol_incompatible`; `detail` uses the exact key `"missing"`, table order `howl.tts.load,howl.tts.synthesize`, no spaces/duplicates, omitted when none are missing. |
| **C1-AC-2** | Service/provider scope, unknown IDs, wrong versions, duplicates, duplicate providers/extensions, and invalid descriptors are validated; shared `synthesize` is accepted; service-only/provider-only misplacement is rejected. |
| **C1-AC-3** | Semantic `hello` failures are terminal `protocol_incompatible` and latched before termination. **Case A (termination succeeds):** the single `Wait()` owner observes exit; `ProcessState=exited`, `ProcessAlive=false`, state remains `incompatible`; no backoff/restart-budget spend/replacement. **Case B (termination times out while live):** state remains `incompatible`; the live-but-unreaped process remains tracked; repeated observation does not change it; no `crashed`, second start/handshake, transport reuse, frame, backoff, budget spend, or replacement; a later exit remains `incompatible`. The C1-AC-24 barrier and aggregate projection are mandatory. |
| **C1-AC-4** | The supervisor-owned snapshot is published only through C1-AC-23's current-instance/current-generation handshake gate and invalidated on stopping, exit, crash, incompatibility, backoff, and replacement; generation is current and copies do not expose a mutable map. |
| **C1-AC-5** | Direct `Unload` with unknown snapshot returns `unsupported_capability` with `capability=howl.tts.unload`, `snapshot=unknown`; no frame, start/spawn/restart, or `unloaded` claim. |
| **C1-AC-6** | Direct `Cancel` validates target first; unknown snapshot returns `unsupported_capability` with `capability=howl.tts.cancel`, `snapshot=unknown`; no frame, start/spawn/restart, or cancellation claim. |
| **C1-AC-7** | Known-absent direct `Unload`, `Cancel`, and `Catalog` return `unsupported_capability` with `snapshot=known`, no frame, no start, and no false success. The C7 `ProviderCatalog` path is separately allowed to start only from ordinary `absent`; it is not folded into this direct criterion. |
| **C1-AC-8** | In a normal live state, known-present direct `Unload` and `Cancel` send exactly one operation frame on the current instance; no additional start; a generation change yields a typed failure. Live `unavailable` follows C1-AC-18 instead: `Unload` may send one current-transport frame without aggregate transition, and `Cancel` returns `request_not_found` with no frame. |
| **C1-AC-9** | Stale after a successful handshake is unknown for stopping, exit, crash, incompatibility, backoff, replacement, and exited/no-instance `unavailable`: no frame, no start/spawn, no false state. Direct `Catalog` follows the same unknown result; it is not silently conflated with ordinary `absent`. |
| **C1-AC-10** | Supervisor cancellation with unknown or absent `cancel` sends no frame, synthesizes no `unsupported_capability`, and follows grace/`KillTree`/`Wait()` and existing precedence. |
| **C1-AC-11** | Direct `Cancel` and the supervisor path are separately implemented; the watcher never calls exported `Cancel`; the two-sided test passes on one helper. |
| **C1-AC-12** | No direct optional operation calls `Start`: direct `Unload`, `Cancel`, `Catalog`, and `Health` use only the current instance or a fail-closed result. Known-present uses `callOnCurrent`; a generation change cannot start a replacement. C7 `ProviderCatalog`/`ProviderStatus` may start/handshake only from ordinary `absent` under contract §17.12.1, never from `unavailable` or another live/terminal state. Both `Client.Start` and `Supervisor.Start` reject `incompatible` before returning an existing process, creating a process, starting a reader, reusing a transport, or sending `hello`; a normal TTS request therefore cannot obtain a second handshake, including while a Case B process remains live. |
| **C1-AC-13** | The documented map covers every provider state, live/exited `unavailable`, and process `exited`; it includes `crashed → backoff → starting`, fourth-attempt exhaustion → `unavailable`, exit-zero/nonzero/grace classification, the A2-D4 Case A/Case B exception and start gate, the A2-D12 no-recovery guard, live wire `health`/`catalog` reads, the distinct live-`unavailable` `Unload`/`Cancel` behavior, and every §8.2 health no-transition edge. A structural test asserts coverage and no optional-operation spawn/start/handshake. |
| **C1-AC-14** | Snapshot publication/read is atomic and single-writer; no torn/stale snapshot; lock order and blocking exclusions are documented; race/overlap tests pass. |
| **C1-AC-15** | **OPEN.** The stage report contains the complete §10 resource record (all seven fields, with a written independent conformance determination per row), complete §9 I/O records for all eight §3.2 categories (all fields and reviewed N/A entries naming the independent reviewer and date), the R8-AC-8 PR-dimension/no-claim statement, and the R8-AC-9 reviewed applicability entries. A blank, pending, or self-authored determination does not satisfy this criterion. |
| **C1-AC-16** | No production import, model load, disk write, public status field, or new public error code/taxonomy is added by C1. |
| **C1-AC-17** | **Delivery-only criterion.** C1 delivery includes the capability matrix, exact tests, failure logs, resource/I/O evidence, the aggregate-state writer/ownership inventory, the state-class operation/error matrix, the reconciliation record, and the C1/C2/C7 boundary statement. Verification is a static delivery manifest/checklist reviewed by the independent stage reviewer; it is not runtime evidence and cannot substitute for C1-AC-15 or any executed test. |
| **C1-AC-18** | Asset-specific `unavailable` on a live handshaken helper retains the authoritative snapshot; exited/no-instance `unavailable` invalidates it. Separate tests prove both. While live `unavailable`: `Health` sends one current-transport wire `health`, returns its result, and leaves state/failure unchanged; advertised `Catalog` sends one current-transport `catalog`, returns its result, and does not call `SetSelected` or change state/failure; advertised `Unload` may send one current-transport frame but MUST NOT set aggregate `unloaded` or clear the failure; advertised `Cancel` returns `request_not_found` with no frame/start; absent optional capabilities return `unsupported_capability` with no frame; `Synthesize` returns `provider_not_loaded` with no start, handshake, ensure-load, load, or frame; status/voice-list do not start. Exited/no-instance `unavailable` follows the same no-recovery rules and `Synthesize` returns `provider_not_loaded`; it delegates any recovery to C3/C5/C6. Only explicit same-identity `load` on a live helper whose helper-side state is `unloaded`, subject to D1a/D1b, may recover; that path is C3/C6 work and C1 adds no other recovery path. |
| **C1-AC-19** | A10a-2/A10b-2 caller projection is exact: unknown/absent direct `Unload`, `Cancel`, and `Catalog` use `Failure.Code=unsupported_capability`, bounded detail, and `AudioError.Code=ErrCodeSynthesisFailed` with the same `Failure` in the cause. In live `unavailable` and exited/no-instance `unavailable`, `Synthesize` uses `provider_not_loaded`/`ErrCodeModelLoadFailed`; `Cancel` returns `request_not_found` in live `unavailable` and `unsupported_capability` when the capability/instance is unavailable; `Health`/`Catalog` return their results or existing failures. No frame for rejected cases; no new taxonomy. |
| **C1-AC-20** | Supervisor known-absent `cancel` is a separately named branch and test: no frame, no `unsupported_capability`, normal grace/termination/precedence; it is not collapsed with unknown. |
| **C1-AC-21** | Concurrency cases cover concurrent direct `Unload`, concurrent direct `Cancel`, direct optional operations during startup/loading/stopping/recovery, supervisor cancellation during startup/shutdown/recovery, exit/replacement/snapshot-invalidation overlaps, incompatible Case B repeated observation/ordinary-TTS/direct-start/shutdown/recovery overlaps, and handshake-attempt/publication/termination/update-reset overlaps; `workMu`/`startMu` and winner rules are documented. |
| **C1-AC-22** | A static-review assertion verifies C1-AC-16: no production import, no public field/error taxonomy, no model/disk path, and no A2-D8 unload-as-shutdown substitution. |
| **C1-AC-23** | Handshake publication is process-owned, instance-bound, and generation-bound: only the private `handshakeAttempt` returned by `beginHandshake(p)` may publish; its `exchange` must have read and validated `hello` through the pinned process/transport, the attempt must still be current, the process current, the generation current, the state `handshaking`, the selected provider must belong to the validated hello, and no terminal intent may be latched. There is no caller-supplied `helloResult` and no second publication method. `incompatible`, `stopping`, `unavailable`, exited/crashed/backoff, replacement, stale generation, foreign process/attempt, and mismatched selected-provider attempts publish no capability snapshot and change no state. Negative tests cover every listed condition plus a positive current-process control. |
| **C1-AC-24** | Incompatible termination keeps aggregate `State=incompatible` throughout the pending window: `Terminate` may set `ProcessState=stopping` but never overwrites aggregate state; a barrier test observes `State=incompatible` while the process is still alive/stopping; expected exit sets `ProcessState=exited`/`ProcessAlive=false` and leaves `incompatible`; no `crashed`, backoff, budget spend, replacement, or handshake occurs. Every writer inventory entry is exercised after `incompatible` and either rejects without mutation or is the authorized exit/update-reset path. |
| **C1-AC-25** | The aggregate-state writer inventory is complete for `providerState`, `lastFailure`, and `processState`; no exported general setter remains; all writes go through the private event transition except construction. A rejected event preserves aggregate state, `LastError`, process state, and `ProcessAlive`; only `observeExit` and a token-gated C6/C7 `updateGatedReset` may change `incompatible`. Static source inventory and negative writer tests cover every current writer identified in §6.6. |
| **C1-AC-26** | The §4.6 matrix and §4.7 adapter table cover ordinary `absent`, unknown/no-current-handshake, live known-absent optional, live known-present optional, live asset-specific `unavailable`, and exited/no-instance `unavailable` for `Health`, `Catalog`, `Unload`, `Cancel`, `Synthesize`, `ProviderStatus`, and voice-list, including exact frame/start/recovery, helper result or `Failure`, caller projection, state/failure effect, and delegation. Direct `Client.Catalog` and C7 `ProviderCatalog` are distinct. Static completeness and state-class tests verify no unknown/known-absent conflation. |

### 10.1 Acceptance-criteria verification and traceability

Every criterion has a verification method. “Executed” remains prospective; this table does not claim that a test ran. C1-AC-15 is explicitly open.

| Criterion | Design locus | Requirement trace | Verification method / disposition |
|---|---|---|---|
| C1-AC-1 | §6.1 | A2-D1, A2-D3 | Static capability table plus missing-mandatory fixtures; exact `detail["missing"]`; provisional executed evidence. |
| C1-AC-2 | §6.2 | A2-D2 | Registry/scope fixtures for unknown, duplicate, version, and misplacement cases; provisional executed evidence. |
| C1-AC-3 | §6.4, §6.5, §8, §9.1 | A2-D4 | Separate Case A success and Case B timeout/live tests; repeated observation, later exit, barrier aggregate-state observation, and no-start assertions; provisional. |
| C1-AC-4 | §3, §6.5, §9.1 | A2-D5, A2-D6 | Static generation/publication review; negative direct/stale publication test; copied-read and invalidation race tests; provisional. |
| C1-AC-5 | §4.1, §4.3 | A2-D7, A2-D10 | Unknown-snapshot direct-unload test; no frame/start/`unloaded`; conditional. |
| C1-AC-6 | §4.2, §4.3 | A2-D7, A2-D10 | Target-validation and unknown-snapshot direct-cancel test; no frame/start. |
| C1-AC-7 | §4.1, §4.2, §4.6 | A2-D7 | Known-absent fixture for `Unload`, `Cancel`, and `Catalog`; no frame/start; separately test C7 `ProviderCatalog` from ordinary `absent`. |
| C1-AC-8 | §4.1, §4.2, §4.4 | A2-D7, D10 | Known-present one-frame send on the captured generation; start counter remains zero. |
| C1-AC-9 | §3, §4.6 | A2-D6 | Stale-state matrix for stopping, exit, crash, incompatible, backoff, replacement, and exited/no-instance `unavailable`; direct `Catalog` unknown is distinct from ordinary `absent`. |
| C1-AC-10 | §5 | A2-D9 | Supervisor unknown and known-absent paths; grace/kill/wait; no capability error. |
| C1-AC-11 | §5, §9.1 | A2-D9 | Static separation review plus two-sided direct/supervisor test on one helper. |
| C1-AC-12 | §4.7, §6.4, §9.1 | A2-D4, A2-D9 | Static no-`Start` review; start-counter test; direct `Catalog` no-start; C7 `ProviderCatalog`/`ProviderStatus` start only from `absent`; incompatible normal/direct start gate, including live Case B. |
| C1-AC-13 | §4.6, §8 | r9 C1-START-MAP | Static map-completeness check plus matrix assertions for live health/catalog/unload/cancel and no optional spawn/start/handshake. |
| C1-AC-14 | §3, §9.1 | Original AC-1/AC-2 | Lock-order/ownership review and race/overlap tests; provisional executed evidence. |
| C1-AC-15 | §9.2, §9.3, §9.4 | Original AC-4, R8-AC-4/7/8/9 | **OPEN.** Static completeness may be checked now, but independent written budget determinations and named/dated I/O N/A reviews do not exist; C1 cannot be acceptance-ready until they do. |
| C1-AC-16 | §7 | Original AC-2/AC-6 | Static import/API/model/disk scan; no public field or taxonomy. |
| C1-AC-17 | §2, §4.6, §6.6, §10, §11, §12 | Original AC-6 delivery trace | **Delivery-only.** Static package manifest/checklist including the C1-R-1–R-23 reconciliation, writer inventory, operation matrix, and C1/C2/C7 boundary; reviewed by the independent stage reviewer; not runtime evidence and not a substitute for C1-AC-15. |
| C1-AC-18 | §3.4, §4.1, §4.2, §4.5, §4.6, §8 | A2-D6, A2-D12 | Live-retain/exited-invalidate tests; health/catalog current-transport reads; unload no-transition; cancel `request_not_found`; absent-capability no-frame tests; exited/no-instance `Synthesize` `provider_not_loaded`; explicit-load path mapped to C3/C6. |
| C1-AC-19 | §4.3, §4.6 | D10 A10a-2/A10b-2 | Error-chain projection tests for `unsupported_capability` (including direct `Catalog`), `provider_not_loaded`, and `request_not_found`; precise `Failure` preserved in `AudioError.Cause`; live and exited/no-instance cases. |
| C1-AC-20 | §5 | A2-D9 | Separately named supervisor known-absent test. |
| C1-AC-21 | §9.1, §11 | Original AC-1/AC-2/AC-3 | Race/overlap matrix, including incompatible Case B, live-`unavailable` no-recovery, attempt publication/termination, and update-reset races. |
| C1-AC-22 | §7, §10, §11 | C1-AC-16 | Static-review assertion for imports, public fields/taxonomy, model/disk paths, and A2-D8 boundary. |
| C1-AC-23 | §3.2, §6.5, §6.7 | A2-D3, A2-D4, A2-D6 | Static entry-point review plus the full negative attempt matrix (all lifecycle states, foreign process/attempt, selected mismatch, replacement, stale generation, concurrent termination) and a positive current-process control; no caller-supplied hello. |
| C1-AC-24 | §6.4, §6.5, §6.7, §8 | A2-D4 | Barrier-based termination test observing aggregate `incompatible` before exit; expected-exit and timeout/late-exit assertions; no crash/backoff/restart; every writer rejection after `incompatible`. |
| C1-AC-25 | §6.5, §6.6, §6.7 | A2-D4, A2-D6 | Static writer inventory and export scan; parameterized rejected-writer tests preserving state/failure/process; only exit/update-token changes `incompatible`. |
| C1-AC-26 | §4.6, §4.7 | A2-D6, A2-D7, A2-D10, A2-D12 | Static matrix-completeness check and state-class tests for all seven operations; direct `Catalog` versus C7 `ProviderCatalog`; no unknown/known-absent conflation. |

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
| Incompatible termination success | Semantic failure; termination succeeds within grace | State `incompatible`; single `Wait` observes exit; `ProcessState=exited`; no backoff/budget/replacement |
| Incompatible timeout alive | Semantic failure; termination exceeds grace; process still live | State remains `incompatible`; process tracked; `ProcessAlive=true`; no crash/backoff/budget/replacement |
| Incompatible repeated observation | Case B process remains live; repeat status/start attempts | No second start, handshake, transport reuse, frame, or state change; one `Wait` owner |
| Incompatible ordinary TTS | Case B process live; ordinary `Synthesize` call | Existing `protocol_incompatible` failure; no `Client.Start`/`Supervisor.Start`, handshake, frame, backoff, or replacement |
| Incompatible direct start | Case B process live; `Client.Start` and `Supervisor.Start` | Both reject before returning/starting; no reader, transport reuse, or `hello`; state remains `incompatible` |
| Incompatible late exit | Semantic failure; delayed process exit after timeout | Expected exit consumed; `ProcessState=exited`, state `incompatible`; no crash/restart |
| Incompatible termination barrier | Semantic failure; termination blocked before exit | While pending: `State=incompatible`, `ProcessState=alive`/`stopping`, `ProcessAlive=true`; after exit: `ProcessState=exited`, `State=incompatible`; no crash/backoff/restart |
| Positive: process-owned publication | Current `handshakeAttempt` on the captured live process, valid compatible hello, selected provider from that hello | Publishes exactly once; state `unloaded`; capabilities/identity current; no second attempt |
| Negative: publication state matrix | Attempt for `incompatible`, `stopping`, live/exited `unavailable`, exited, crashed, backoff, replacement, or stale generation | No publication, no state/identity/failure change, no valid snapshot |
| Negative: foreign process/attempt | Current process pointer but hello read by a different attempt/transport; or a stale attempt object | Rejected; no publication; caller cannot supply a hello result |
| Negative: selected mismatch | Valid attempt but `selected` not in its validated hello, or wrong provider/engine | Rejected atomically; no partial identity publication |
| Negative: concurrent termination | Publication races termination intent/Case A/Case B exit | Either publication wins before intent and then termination preserves it, or termination wins and publication rejects; never both |
| Known absent: unload | Valid handshake without `unload` | `unsupported_capability`, `snapshot=known`; no frame/start |
| Known absent: cancel | Valid handshake without `cancel` | `unsupported_capability`, `snapshot=known`; no frame/start |
| Known absent: catalog | Valid handshake without `catalog` | `unsupported_capability`, `snapshot=known`; no frame/start; direct and adapter paths both fail closed on the live instance |
| Unknown: unload | New client, no handshake | `unsupported_capability`, `snapshot=unknown`; no frame/start/`unloaded` claim |
| Unknown: cancel | New client, valid target | `unsupported_capability`, `snapshot=unknown`; no frame/start/cancellation claim |
| Unknown: catalog | New client, no handshake; direct and C7 adapter | `unsupported_capability`, `snapshot=unknown`; no frame; direct never starts; adapter starts only from ordinary `absent` |
| Invalid cancel target | Malformed target, no snapshot | `invalid_request`; no frame/start; gate not reached |
| D10 caller projection | Unknown/known-absent direct unload, cancel, and catalog | `AudioError.Code=ErrCodeSynthesisFailed`; cause preserves `unsupported_capability` and exact detail; no new code |
| Stale: exit | Handshake, then process exit | Unknown; no frame/start; no false state |
| Stale: stop | Handshake, then stopping | Unknown; no frame/start |
| Stale: crash | Handshake, then unexpected exit | Unknown; no frame/start; no restart |
| Stale: incompatible | Failed handshake after prior snapshot | Unknown; no frame/start; terminal |
| Stale: backoff/replacement | Constructed generation change | Unknown; no frame/start/replacement |
| Live unavailable: health | Live handshaken helper, typed asset-specific load failure; `Health` | One current-transport `health` frame; health result returned; no start/handshake/load; state remains `unavailable`; `LastError` retained |
| Live unavailable: catalog present | Same live state; `Catalog` advertised | One current-transport `catalog` frame; catalog result returned; no start/handshake/load/`SetSelected`; state and failure unchanged |
| Live unavailable: catalog absent | Same live state; `catalog` not advertised | `unsupported_capability`, `snapshot=known`; no frame/start; state unchanged |
| Live unavailable: unload | Same live state; `Unload` advertised | One current-transport unload frame; success does not set aggregate `unloaded`, clear `LastError`, or change snapshot; state remains `unavailable` |
| Live unavailable: cancel | Same live state; `Cancel` advertised, valid target | `request_not_found`; no frame/start/state change; `AudioError.Code=ErrCodeSynthesisFailed`; cause preserves precise code |
| Live unavailable: ordinary synthesis | Same live failed state; ordinary `Synthesize` | `provider_not_loaded`; no start, handshake, ensure-load, load, or frame; state remains `unavailable` |
| Live unavailable: status/voice-list | Same live failed state; status and voice-list | No frame/start/load; existing supervisor/static surfaces; state remains `unavailable` |
| Unavailable explicit-load mapping | Live helper, helper-side `unloaded`; explicit same-identity `load`; D1a/D1b satisfied | Only permitted recovery path is mapped to C3/C6; C1 adds no competing recovery and claims no executed load |
| Exited/no-instance unavailable | Discovery/factory failure or exit | Snapshot invalidated; direct unload/cancel/catalog fail `unsupported_capability` unknown; health supervisor observation; synthesis `provider_not_loaded`; no start/recovery; C3/C5/C6 recovery only |
| Ordinary absent vs unknown | Fresh client at `absent` versus stopped/exited state | Direct catalog is unknown/no-start in both; C7 `ProviderCatalog` may start only from ordinary `absent`; no conflation |
| Direct versus adapter catalog | Client.Catalog on absent/unknown versus ProviderCatalog on absent | Direct never starts; adapter starts/handshakes/catalogs without load; identical live read behavior |
| Writer inventory: SetState removed | Static source scan plus attempted call through test seam | No exported general setter; all writes through `apply`; test build cannot call a removed method |
| Writer rejection after incompatible | Exercise each §6.6 writer/event after `incompatible` | Aggregate/`LastError`/process state unchanged; `Client.Start`/`Supervisor.Start` still gated; no publication/crash/backoff/replacement |
| Writer rejection after live unavailable | Exercise load success, unload success, synthesis, handshake, start, ordinary terminate | Only C6/C7 load path, process termination intent, exit, and update token are accepted; no recovery or terminal overwrite |
| Update-gated reset | Valid C6/C7 token/epoch with no live process versus missing/reused token/live process | Valid only: `incompatible`/`unavailable` → `absent`, clear failure, invalidate snapshot, advance generation; all invalid attempts no-op |
| State-class matrix completeness | Static matrix check and parameterized tests | All six classes × seven operations covered; no blank cell; direct adapter distinction; exact `Failure`/`AudioError`/state |
| Supervisor unknown cancel | Active synthesis, unknown snapshot | No frame; no `unsupported_capability`; grace/kill/wait |
| Supervisor known-absent cancel | Active synthesis, snapshot known but `cancel` absent | Separately named: no frame; no capability error; grace/kill/wait |
| Two-sided regression | Same helper, direct vs supervisor | Direct D10 result; supervisor A2-D9 path; no cross-leak |
| Concurrency: concurrent unload | Two direct `Unload` calls | `workMu` serializes; each re-gates; no interleaved frames; no start; no contradictory state |
| Concurrency: concurrent cancel | Two direct `Cancel` calls, same/different targets | Each own request ID; writes serialized; typed ack/error; no fabricated success or terminal overwrite |
| Concurrency: startup | Direct optional during `starting`/`handshaking` | Unknown; no frame/start; handshake unaffected |
| Concurrency: loading | Direct `Unload`/`Cancel` issued during load | `Unload` waits on `workMu` and re-gates; `Cancel` only if valid target; no load/start corruption |
| Concurrency: stopping | Direct optional during stopping | Unknown/invalid; no frame/start |
| Concurrency: unavailable guard | Direct optional or ordinary synthesis races entry to `unavailable` | No frame/start/load; state remains `unavailable`; snapshot retained if live, invalidated if exited |
| Concurrency: incompatible | Case B exit/observation/start/shutdown/recovery races | No second start/handshake/transport reuse; late exit remains `incompatible`; no backoff/budget/replacement |
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
| A2-D4 — incompatible terminal, no restart, intent-aware exit | §6.4, §6.5 | C1-AC-3/C1-AC-24; timeout, barrier, and late-exit tests |
| A2-D5/A2-D6 — status and known snapshot | §3, §6.5, §7 | C1-AC-4/C1-AC-18/C1-AC-23; live/exited split and publication gate |
| A2-D7 — known-absent direct `unload`; unload is not termination | §4.1, §6.1 | C1-AC-7; known-absent test |
| A2-D8 — idle shutdown, not unload | §4.1 boundary only; no C1 implementation | No C1 coverage claim; C4/C7 own shutdown; static AC-22 |
| A2-D9 — supervisor cancellation | §5 | C1-AC-10/C1-AC-11/C1-AC-20; three-path tests |
| A2-D10 — no silent no-op | §4, §8 | C1-AC-5/C1-AC-7/C1-AC-9 |
| A2-D12 — `unavailable` recovery | §3.4, §4.1, §4.2, §4.5, §4.6, §8 | C1-AC-18/C1-AC-21; live/exited, health/catalog read, unload/cancel, and no-recovery tests; explicit load mapped C3/C6 |
| C1-R-17 — `MarkHandshake` publication gate | §3.2, §6.5, §6.7, §10.1 | C1-AC-23; process-owned attempt and full negative matrix |
| C1-R-18 — aggregate state during termination | §6.4, §6.5, §6.7, §10.1 | C1-AC-24; barrier, late-exit, and writer-rejection tests |
| C1-R-20 — aggregate-state mutation surface | §6.5, §6.6, §6.7, §10.1 | C1-AC-25; inventory and rejected-writer tests |
| C1-R-21 — handshake provenance | §3.2, §6.5, §6.7, §10.1 | C1-AC-23; foreign-process/attempt tests |
| C1-R-22 — Catalog start authority | §4.6, §4.7, §10.1 | C1-AC-7/C1-AC-12/C1-AC-26; direct-versus-adapter tests |
| C1-R-23 — complete projection matrix | §4.5, §4.6, §4.7, §10.1 | C1-AC-7/C1-AC-9/C1-AC-18/C1-AC-19/C1-AC-26 |
| D10(a) A10a-2 | §4.1, §4.3 | C1-AC-5/C1-AC-19 |
| D10(b) A10b-2 | §4.2, §4.3 | C1-AC-6/C1-AC-19 |
| r9 C1-CAPABILITY | §3, §9.1 | C1-AC-4/C1-AC-14 |
| r9 C1-START-MAP | §8 | C1-AC-13 |
| Original AC-1/AC-2/AC-3 | §9.1 | C1-AC-14/C1-AC-21; concurrency matrix |
| Original AC-4 / R8-AC-4 — budget record | §9.2 | C1-AC-15; seven-field record |
| Original AC-5 — safety under load | §4, §9.1 | C1-AC-12/C1-AC-14/C1-AC-21 |
| Original AC-6 — no hidden background work | §9.1, §9.2 | C1-AC-15/C1-AC-17 |
| C1-AC-17 — delivery package (delivery-only) | §2, §4.6, §6.6, §10.1 | Static manifest/checklist; not runtime evidence |
| C1-AC-25 — aggregate-state ownership | §6.5-§6.7 | Static writer inventory/export scan and rejected-writer tests |
| C1-AC-26 — complete operation matrix | §4.6-§4.7 | Six-state-class by seven-operation static/test matrix |
| R8-AC-7 — application I/O | §9.3 | C1-AC-15; eight-category record |
| R8-AC-8 — performance impact | §9.4 | C1-AC-15; PR-dimension statement, no claim |
| R8-AC-9 — applicability review | §9.4 | C1-AC-15; reviewed N/A table |
| **Original AC-7 — no stage substitution** | **No C1 design content; C8-only** | **Not a C1 acceptance claim** |
| R8-AC-10 — pointer to original AC-7 | No C1 design content; C8-only | Not a C1 acceptance claim |

Bare `AC-7` is never used here for I/O evidence. The I/O rule is R8-AC-7.

## 13. Later-stage dependencies

- **C2:** public `Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown`, and status projection.
- **C3:** ordinary load/ensure-load behavior; C1 adds no load semantics. C3 is the only place that may implement the explicit-load recovery bypass; it must enforce live helper, same identity, helper-side `unloaded`/`accepting_work=true`, and D1a/D1b.
- **C4:** full exit intent, normal shutdown, timeout, grace, and stop classification. C1 supplies the incompatible-intent exception, the live Case B gate, and generation invalidation.
- **C5:** executed crash → backoff → starting, 1/2/5 s delays, three-restart/10-minute budget, fourth-attempt exhaustion, and idle close.
- **C6:** identity matching and replacement; C1 increments generation at replacement begin. C6 owns executed same-identity recovery and supplies the update token/epoch for the C1 `updateGatedReset` guard. C6 must not bypass the guard with a general state write.
- **C7:** production adapter, explicit load reachability, executed start evidence, and the only `ProviderCatalog`/`ProviderStatus` start from ordinary `absent`. C7 must not bypass C1's direct `Catalog` no-start rule, the `incompatible` gate, or the `unavailable` no-recovery rule.

No downstream stage is implemented or accepted here.

**C4 note:** the aggregate-state ownership rule in §6.5 is a C1 constraint on the existing `Terminate`/exit-observer path. C4 owns the full exit/shutdown state machine and must preserve it; C1 does not implement a second exit observer.

## 14. Authorization boundary

C1 remains unauthorized. This revision is design and documentation only. It does not modify source, implement C1, run tests/builds/runtime/benchmarks, deploy, authorize C2–C8, change A2-D9, amend the contract, amend the adopted r9 baseline, amend the frozen plan, or alter any owner decision. C1-AC-15 is open; this revision is not acceptance-ready and does not claim acceptance closure. All C1-R-1 through C1-R-23 dispositions are author-reconciled, not independent acceptance. Independent review by Muse of revision 6 is required before any C1 authorization request.
