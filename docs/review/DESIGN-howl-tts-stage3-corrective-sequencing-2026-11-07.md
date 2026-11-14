# DESIGN — Howl TTS Stage 3 corrective design and sequencing proposal

**Date:** 2026-11-07  
**Mode:** DESIGN ONLY  
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md`  
**Proposal status:** Not implementation authorization.  
**Revision:** Corrected 2026-11-07 in response to the adversarial design audits `71a17369e2389fae26e86783fc4df802ab91bad3daef512c0e7f743067e94ea2`, `a7384da9f3f38661fc28b3de48badb82245cfd41efb8e49964f17bcfaaf5a7e2`, and `6681e11b7eb770d59f3c79294206f05c8d3a3d9d4cfbbd038ef103fc8c48fd15`. This revision remains design-only and authorizes no implementation.

This proposal converts the Stage 3 readiness investigation into a bounded, dependency-ordered corrective program. It preserves the accepted Stage 1/Stage 2 process-ownership and protocol machinery wherever the governing contract permits a narrow correction. It does not authorize or describe an implementation already performed.

## 1. Executive recommendation

Proceed in eight reviewable stages, in this order:

1. **C1 — Capability negotiation and semantic classification.**
2. **C2 — Status and capability projection.**
3. **C3 — Lifecycle state, load gating, and ordinary-load seam.**
4. **C4 — Exit intent, shutdown, timeout, and cancellation consistency.**
5. **C5 — Restart budget, backoff, and idle close.**
6. **C6 — Identity matching, unavailable recovery, and replacement.**
7. **C7 — Production adapter and explicit load reachability.**
8. **C8 — Adopted-contract test closure and separately authorized runtime verification.**

The critical path is C1 → C2 → C3 → C4 → C5 → C6 → C7 → C8. C2 and C4 both touch status and process semantics, so parallelizing them would increase merge and state-model risk more than it saves. C3 should not be folded into C2 merely to show lifecycle progress: taxonomy/projection and actual load behavior have separate failure modes and separate rollback boundaries. C4 and C5 are separate because exit intent must be correct before restart budgeting and idle timers can be trusted.

The accepted foundation should be preserved:

- parent-owned stdio process and one `Wait()` observer;
- Windows suspended startup and job-object ownership;
- strict NDJSON framing, request correlation, terminal validation, acceptance barrier, and cancellation precedence;
- exact-path output cleanup, WAV validation, hash validation, and no-overwrite publication;
- `CGO_ENABLED=0` for the main module;
- the frozen wire `Capability` shape `{id, version integer}`;
- unchanged `types.TTSSynthesizer`.

The design does **not** find a new contradiction in the adopted contract beyond the already-recorded deferred B12a/B12b and O-1–O-7 matters. It identifies three contract-surface gaps (D1a, D3, D10) and five owner decisions (D1b, D2, D5, D9, D10); D10 is both an owner decision and a contract-surface gap because the contract is silent. A2-D12 requires recovery from `unavailable` to be triggered by an explicit `load`, but §17.12.1 defines no `load`, `unload`, or `cancel` entry point. **D1a** asks whether a narrow internal loader seam outside that table is permitted or whether an adoption-record clarification is required; **D1b** asks which existing explicit product action invokes it and whether the adapter survives that action for in-place recovery. **D2** and **D5** remain genuine product/architecture choices. **D9** asks only whether status/catalog/health observation resets the 120 s idle deadline; the separate question whether an explicit status request may start an `absent` helper is already settled by §17.7.2 and is not part of D9. **D10** now explicitly records **two** unknown-snapshot direct-operation questions—**(a)** direct `unload` and **(b)** direct `Cancel` outside the supervisor path—because the contract does not settle either result. The supervisor cancellation path is not part of D10: A2-D9 settles it. D10 originated as the unload branch; the direct-`Cancel` branch is separately answerable and is not assumed to receive the same disposition. No public API or new error code is assumed. These decisions are not resolved by this proposal and require separate authorization.

## 2. Verified input hashes

Recomputed before this corrected proposal:

| Input | SHA-256 |
|---|---|
| Governing architecture | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` |
| Stage 3 readiness investigation | `a63f9f32f4d1bd047746a7bbbf5f509588e13ebd4272cffeb23f61fcfa891c6c` |
| Stage 2 implementation report | `ec85b2b2c3760542699bcd406f1f3b6b3eb2cdc0165b9d6a3b21977a06741227` |
| Stage 2 capability correction | `c9f498f61e0444ab9a66b347f466da7ac456e5afb0a30534621ab453de83857f` |
| Stage 2 adversarial audit | `e73cef36491d2d77235c31cc051c3e756ac73d78389ca2d0f90967b9dee18755` |
| Stage 2 correction audit | `d6b2349e030ac8867d236d571a5d3e309eb6c35dfde2a7e5b6586a51cc8df969` |
| Stage 3 corrective-design audit | `71a17369e2389fae26e86783fc4df802ab91bad3daef512c0e7f743067e94ea2` |

The Stage 1 report, correction, and audit records remain the historical inputs indexed in `docs/review/README.md`. Their actual hashes are recorded in the readiness investigation and were not reinterpreted here.

## 3. Readiness findings mapped to corrective stages

| Finding | Corrective stage(s) | Why there |
|---|---|---|
| S3R-01 five mandatory capabilities | C1 | Capability floor and handshake validation are one acceptance boundary. |
| S3R-02 semantic hello is not `protocol_incompatible` | C1 | Classification and terminal handshake behavior belong with negotiation. |
| S3R-03 service/provider scope conflation | C1, C2 | C1 validates scopes; C2 projects them separately. |
| S3R-04 `detail.missing` wrong; no typed detail | C1, C2 | C1 encodes wire detail; C2 exposes application detail. |
| S3R-05 status model incomplete | C2 | Pure status/taxonomy stage. |
| S3R-06 process/available/ready contradictions | C2, C4 | C2 fixes projection types/invariants; C4 fixes exit/lifecycle truth. |
| S3R-07 no ordinary unloaded ensure-load | C3, C7 | C3 supplies state/load seam; C7 applies it in the production adapter. |
| S3R-08 unavailable recovery ungated/unreachable | C3, C6, C7 | C3 blocks synthesis recovery; C6 implements the identity-gated path; C7 exposes it. |
| S3R-09 no identity tracking | C3, C6 | C3 captures identity; C6 enforces matching/replacement. |
| S3R-10 no restart budget/backoff | C5 | Requires correct exit intent from C4. |
| S3R-11 no idle close | C5 | Requires shutdown classification from C4. |
| S3R-12 shutdown misclassifies nonzero exit | C4 | Exit intent and state truth. |
| S3R-13 timeout/cancel can report ready after kill | C4 | Post-kill status consistency. |
| S3R-14 unload/cancel gate before handshake | C1 | Validated-handshake invariant and no-frame rule. |
| S3R-15 no production adapter/caller | C7 | Adapter and reachability. |
| S3R-16 tests encode superseded policy | C1–C7, C8 | Each stage owns its tests; C8 closes the matrix. |
| S3R-17 OS-boundary/native evidence incomplete | C4, C5, C8 | Windows process tests and runtime verification. |
| S3R-18 historical audit superseded | Not corrective | Provenance only; cite correction audit. |

## 4. Current architecture and preserved invariants

### 4.1 Current layers

```text
App / audio manager
  -> Pipeline (native TTS engine map, cache, character request resolution)
  -> [no helper production path]
  -> helper.Supervisor (process, one Wait, stderr, status)
  -> helper.Client (handshake, requests, terminals, cancellation)
  -> parent-owned stdio NDJSON
  -> fakehelper (test child) or future native helper
```

The helper packages are currently inert from production. `cmd/howl-tts-helper-fake/main.go` is the only importer of the fake helper. `types.TTSSynthesizer` remains the stable application interface. `TTSProviderInspector` is declared and matched by `helper.Client`, but has no production consumer.

### 4.2 Must remain untouched

- `internal/proc/job_windows.go` suspended-start/job ownership and `job_other.go` fallback.
- One `Wait()` observer in `Supervisor.observeExit`.
- Parent-owned stdin/stdout/stderr pipes; helper stdout is protocol only.
- Frozen wire framing, request-ID rules, JSON strictness, 4 MiB frame limit, and 64 KiB text limit.
- Parent acceptance barrier, cancellation precedence, retired IDs, duplicate-terminal handling, and terminal publication barrier.
- `output.go` exact-path WAV validation, hash validation, and exact-path cleanup.
- Fake helper's no-replace publication guard. B12a's hard-link versus atomic-rename decision remains deferred.
- `internal/audio/types/interfaces.go:127` `TTSSynthesizer`.
- Main module `CGO_ENABLED=0`; no helper import from the main dependency graph outside its own package boundary.
- Character voice resolution, cache key, and per-character engine precedence.
- Two-helper-process cap.

### 4.3 Must change

- Service capability floor and semantic handshake classification.
- Status projection and typed capability detail.
- Process/provider readiness truth and `Loaded` semantics.
- Explicit load and unavailable recovery gating.
- Identity capture, replacement ordering, restart budget/backoff, idle close, and exit classification.
- Production adapter and explicit load reachability.

### 4.4 New behavior required by §17.19

- Required service capabilities: `howl.tts.load`, `howl.tts.synthesize`.
- Optional service capabilities: `howl.tts.catalog`, `howl.tts.unload`, `howl.tts.cancel`.
- Missing required capability is terminal `protocol_incompatible`.
- A direct optional operation (`unload` or `Cancel`) whose capability is known absent returns `unsupported_capability` with no frame (A2-D1, A2-D7).
- The supervisor cancellation path is not a direct `Cancel` call. When `howl.tts.cancel` is known absent or unknown, that path records `cancel_requested`, skips the cancel frame, synthesizes no `unsupported_capability`, and follows the existing grace/termination rules (A2-D9). The path must use a path-scoped internal operation so the exemption cannot leak into the exported direct method.
- If the authoritative service-capability snapshot is unknown, no optional operation may start a helper merely to discover it; §17.7.2's lazy-start row limits starts to explicit TTS/catalog/status requests, not `unload` or `cancel`. The cancellation path still skips the frame and follows grace/termination. The direct `unload` and direct `Cancel` unknown-snapshot results are unresolved and recorded as **D10(a)** and **D10(b)** respectively; the supervisor cancellation path is A2-D9-settled and outside D10. C1 freezes no code for either direct branch. No separate unsupported-capability frame or wire message exists; capability absence is expressed, if at all, by the A2-D1 result code and never by a frame.
- `ServiceCapabilities` is process-scoped; `Capabilities` is provider-scoped.
- `Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown`, `TTSCapability.Supported`, and `TTSCapability.Detail`.
- Ordinary `unloaded` synthesis may ensure-load; `unavailable` synthesis must not recover.
- Explicit `load` is the only recovery trigger; §17.12.1 currently defines no production entry point for it (D1a/D1b).
- Same-identity recovery and A2-D7 replacement.
- Restart budget, backoff, idle close, intent-aware exits.

## 5. Detailed stage specifications

### C1 — Capability negotiation and semantic classification

**Scope.** Correct the service capability floor, scope validation, handshake error classification, and the distinct handling of direct optional operations versus the supervisor cancellation path when an optional capability is absent or unknown.

**Prerequisites.** None beyond design acceptance. D10(a) and D10(b) must both be dispositioned with a concrete caller-facing result before C1 acceptance; a deferral alone does not satisfy this gate. Until then C1 implements only the known-present and known-absent branches and freezes no unknown-snapshot result for either direct branch. This is the first code stage.

**Likely affected files.**

- `internal/audio/tts/helper/protocol.go`: replace the five mandatory set with required/optional policy; add A2-D3 semantic validation, including duplicate extension names and invalid handshake descriptors; do not reject wire capability arrays for order—only normalise them and order the `detail.missing` projection; make `validateHello`/`validateCapabilities` distinguish semantic incompatibility from frame violations; encode `detail.missing` with only required IDs.
- `internal/audio/tts/helper/client.go`: store the validated service-capability snapshot; `helloHandshake` returns `protocol_incompatible` for semantic failures; `Start` terminates/observes and sets `incompatible`; `requireServiceCapability` uses only the validated snapshot. `Unload` sends a frame only when the authoritative snapshot is known and advertises `howl.tts.unload`; when known absent it returns `unsupported_capability` and sends no frame; when unknown it sends no frame and does not spawn, and C1 does not choose a result code—D10(a) governs. The exported `Cancel` remains a direct optional operation: known absent returns `unsupported_capability` with no frame; when unknown it sends no frame and does not spawn, and D10(b) governs. The supervisor cancellation path uses a separate internal path-scoped operation; the existing cancellation watcher—the call site that currently invokes exported `Cancel`—must call that internal operation instead. When `cancel` is absent or unknown the path skips the frame, synthesizes no capability error, and continues the normal cancel grace/termination rules. No optional operation may start a process merely to discover a capability.
- `internal/audio/tts/helper/supervisor.go`: ensure `MarkHandshake` stores service and provider capabilities separately; no status expansion yet.
- `internal/audio/tts/helper/errors.go`: preserve precise code/message/field/detail through `Failure`.
- `internal/audio/tts/fakehelper/server.go` and `internal/audio/tts/fakehelper/protocol.go`: default service set is required two plus optional three; add test modes that omit each optional and each required capability; keep wire `{id,version}`.
- `internal/audio/tts/helper/protocol_test.go`, `client_test.go`, `golden_test.go`, `interop_test.go`, `fakehelper/server_test.go`, `testdata/wp01.helper.json` where fixture expectations change.

**Tests.**

- Required present, each optional present/absent independently and in combinations.
- Missing `load` and missing `synthesize`, both missing; `detail.missing` order has no spaces/duplicates and contains only missing required IDs.
- Unknown service ID; wrong version; duplicate capability; duplicate provider ID; duplicate service ID; duplicate extension name; invalid handshake descriptor; service-only capability in provider descriptor; provider-only feature in service array; shared `synthesize` in both arrays accepted.
- Semantic failures → `protocol_incompatible`, process terminated and `Wait()` observed, state `incompatible`, no `backoff`, no restart.
- Direct `unload` with `howl.tts.unload` advertised: exactly one unload frame and normal response.
- Direct `unload` with the capability known absent: `unsupported_capability`, no frame, no spawn (assert transport send count/bytes).
- Direct `unload` with unknown capability knowledge: C1 tests only the structural no-frame/no-spawn boundary; it does not assert a result code. The result is deferred to D10(a); C1 is not accepted until D10(a) has a concrete caller-facing disposition. A scheduling deferral does not satisfy that gate.
- Direct `Cancel` with `howl.tts.cancel` advertised: exactly one cancel frame and normal response.
- Direct `Cancel` with the capability known absent: `unsupported_capability`, no frame, no spawn.
- Direct `Cancel` with unknown capability knowledge: C1 tests only the structural no-frame/no-spawn boundary; it does not assert a result code. The result is deferred to D10(b); C1 is not accepted until D10(b) has a concrete caller-facing disposition. A scheduling deferral does not satisfy that gate.
- Supervisor cancellation path with `howl.tts.cancel` advertised: exactly one cancel frame and existing precedence preserved.
- Supervisor cancellation path with the capability known absent: no cancel frame, no `unsupported_capability`, normal cancel grace, and a missing terminal leads to `KillTree`/`Wait()` with the request ending `cancelled` or `timeout`.
- Supervisor cancellation path with unknown capability knowledge: no cancel frame, no capability error, no helper start solely to discover the capability; existing cancellation precedence remains authoritative. Any direct `Cancel` unknown-snapshot result is D10(b), not frozen here.
- Regression test: a `cancel`-less helper on the supervisor path produces no cancel frame and no `unsupported_capability` result, then exits through the existing grace/kill path; a direct `Cancel` on the same helper returns `unsupported_capability` with no frame.
- No path emits a frame merely to report unsupported capability; direct unload absence, direct `Cancel` absence, and path cancellation absence are asserted separately.
- Golden fixture still validates `{id,version}` with no `supported` field.

**Negative/failure injection.** Malformed but well-framed hello; duplicate extension name (must become `protocol_incompatible`, not `protocol_violation`); invalid handshake descriptor; missing capability; unknown capability; duplicated JSON keys (must remain `invalid_frame`, not `protocol_incompatible`); semantic failure when the helper remains alive; failure to terminate.

**Race/concurrency.** Handshake completion versus reader failure; capability snapshot read under lock; no frame send after optional-gate failure. Run focused `-race` tests on the helper package.

**Exclusions.** No status fields, recovery, identity, budget, idle, or production adapter.

**Report artifact.** C1 implementation report with capability matrix, exact tests, failure logs, and a statement that no production path imports the helper.

**Independent review.** Review the policy table against A2-D1–A2-D3, A2-D7, A2-D9, and A2-D10; inspect every path that can emit a service operation frame; confirm direct `unload`, direct `Cancel`, and the supervisor cancellation path are three separately implemented branches; confirm the path uses an internal path-scoped operation and direct `Cancel` still returns the direct-operation result; confirm no unknown-snapshot path spawns; separate D10(a)'s direct-unload result and D10(b)'s direct-`Cancel` result from the A2-D9-settled supervisor path; confirm duplicate-extension and invalid-descriptor cases are classified.

**Acceptance.** Required/optional cases pass; direct unload absence returns `unsupported_capability` with no frame; direct `Cancel` absence returns `unsupported_capability` with no frame; supervisor cancellation-path absence emits no frame and no `unsupported_capability`; no path emits a frame merely to report unsupported capability; duplicate extension names and invalid handshake descriptors are `protocol_incompatible`, not `protocol_violation`; direct `unload` and direct `Cancel` with unknown capability knowledge are tested only at the structural no-frame/no-spawn boundary, and no unknown-snapshot result code is claimed for D10(a) or D10(b) before a concrete disposition is recorded; semantic hello failures are terminal and non-restartable; existing framing/acceptance tests remain green. D10(a) and D10(b) must both be dispositioned with concrete caller-facing results before C1 acceptance; a deferral alone does not satisfy the gate; this list does not freeze an unknown-snapshot code.

**Runtime evidence.** Fake-helper real-process handshake tests are sufficient for C1. The statement that Windows job behavior is inherited applies to C1 only; it may not be reused for C4–C8, which change termination and replacement paths.

**Rollback.** Revert the protocol/client changes and fixture expectations; no production caller exists, so the blast radius is the helper package and fake helper.

---

### C2 — Status and capability projection

**Scope.** Add the adopted application status shape and project process, provider, service, and loaded state from authoritative snapshots without taking ownership of lifecycle transitions that arrive in C3–C6.

**Prerequisites.** C1 validated handshake snapshot and scope separation.

**Likely affected files.**

- `internal/audio/types/provider.go`: `TTSCapability{ID,Version,Supported,Detail}`; add `Loaded bool`, `ServiceCapabilities []TTSCapability`, `ServiceCapabilitiesKnown bool` to `TTSProviderStatus`.
- `internal/audio/tts/helper/supervisor.go`: store provider `Capabilities` and process `ServiceCapabilities`; track whether the snapshot is known; project `ProcessAlive` solely from `TTSProcessState`; project `Loaded` and `Ready` from the authoritative provider snapshot; project `Available` from configured-provider state and the two contract-defined `unavailable` cases, not process or aggregate state alone. C2 adds fields and projections only; C3–C6 own the transitions that write them.
- `internal/audio/tts/helper/client.go`: capture all five service IDs during handshake and project present/absent as `Supported`; retain provider capabilities separately; expose the snapshot-known projection. C2 does not implement stop/exit/backoff/replacement invalidation; C3–C6 own those transitions and update the authoritative snapshot.
- `internal/audio/tts/helper/supervisor_test.go`, `client_test.go`, `interop_test.go`: status matrix.

**Required projection rules.**

- `Capabilities` contains only provider-scoped advertised capabilities.
- `ServiceCapabilities` contains all five service IDs only when `ServiceCapabilitiesKnown=true`; absent optional IDs appear with `Supported=false`.
- Before handshake, or after an owning stage establishes exit/stop/incompatible/backoff/replacement: `ServiceCapabilitiesKnown=false` and `ServiceCapabilities` is empty. C2 projects this rule but does not cause those states; C3–C6 own the invalidation transitions. Do not synthesize five false values.
- `Loaded` is projected verbatim from the authoritative provider snapshot. It is never derived from `Available`, `Ready`, `State`, or `ProcessAlive`. C2 does not define or accept the future transition mechanics: C3 owns load/unload and ordinary-load transitions, C4 owns exit/stop, C5 owns budget/backoff, and C6 owns replacement. The A2-D5 rule remains the target invariant; each transition is accepted in its owning stage.
- `Ready` is true only when the projected `Loaded` is true and aggregate state is `ready`, `busy`, or `degraded`.
- `Available` is true for a configured, compatible provider in `unloaded`/`loading`/`ready`/`busy`/`degraded`, and for both legal `unavailable` cases below. It is false for `absent`, `starting`, `handshaking`, `stopping`, `crashed`, `incompatible`, and `backoff`, and for `unavailable` caused by a missing helper or invalid provider configuration. The two legal `unavailable` cases below remain true.
- The two `unavailable` cases are distinct and must never be collapsed:

| Case | `Available` | `ProcessAlive` | `Loaded` | `Ready` | `ServiceCapabilitiesKnown` |
|---|---:|---:|---:|---:|---:|
| (a) Asset-specific failure while the helper remains alive | true | true | false | false | true |
| (b) Restart budget exhausted after the process exited | true | false | false | false | false |

- `ProcessAlive` is true only for process states `starting`, `alive`, `stopping`; false for `absent`, `exited`. It is derived exclusively from actual process observation (`TTSProcessState`) and is never inferred from aggregate `State`. A live helper in asset-specific `unavailable` is a legal state; it is not suppressed or reclassified. Conversely, aggregate `unavailable` never proves that a process is alive.
- The C2 matrix must include at least `absent`, `starting`, `handshaking`, `stopping`, `crashed`, `incompatible`, and `backoff`; `Available` is false in every one, and the remaining flags follow A2-D5, the two `unavailable` rows, and the process rules above.

**Single source of truth.**

| Field | Authoritative source | C2 responsibility | Transition owner |
|---|---|---|---|
| `State` | Supervisor aggregate state | Read/project | C3–C6 |
| `ProcessState` | Process observer | Read/project | C4 |
| `ProcessAlive` | `ProcessState` only | Derive | C4 |
| `Capabilities` | Provider handshake descriptor | Project provider scope only | C1/C2 |
| `ServiceCapabilitiesKnown` | Authoritative current-instance service snapshot | Project | C3–C6 invalidate |
| `ServiceCapabilities` | Service snapshot when known | Project all five or empty | C1/C2 |
| `Available` | Configured-provider/compatibility snapshot plus the two `unavailable` cases | Project | C3–C6 |
| `Loaded` | Authoritative provider loaded flag | Project only; do not write | C3, C4, C5, C6 |
| `Ready` | Projected `Loaded` plus aggregate `State` | Derive | C3–C6 |
| `LastError` | Supervisor's last typed failure | Read/project | C1–C6 |

C2 must not claim `backoff` or identity-replacement behavior before C5/C6 implement those states. Those rows are forward contract requirements, not C2 acceptance.

**Tests.** Full matrix of process state, provider state, available, loaded, ready, and service-known for states reachable or snapshot-constructible at C2: handshake success, `absent`, `starting`, `handshaking`, `stopping`, `crashed`, `incompatible`, asset-specific unavailable with a live helper, budget-exhausted unavailable with `ProcessAlive=false`, and stopped. Assert both `unavailable` rows separately with their exact flag sets. Assert the full `Available=false` enumeration above, including the seven aggregate states and the missing-helper/configuration causes. Assert that `ProcessAlive` follows only `ProcessState`, including live-helper asset-specific `unavailable`. Assert provider/service list separation and no five-false default. `backoff` and identity-replacement transition tests are explicitly deferred to C5/C6 and are marked “not yet reachable” in the C2 matrix.

**Exclusions.** No load/recovery behavior change; no restart/idle; no backoff or identity-replacement transition ownership.

**Report artifact.** Status matrix and mapping table from state to fields.

**Independent review.** Confirm the renderer cannot infer readiness from capability presence; confirm no provider feature is projected as a service capability.

**Acceptance.** Every status field has the single documented source in the table above; both `unavailable` flag sets pass; a live helper may coexist with asset-specific `unavailable`; `ProcessAlive` is never inferred from aggregate `State`; all C2 matrix tests pass; backoff/replacement transition acceptance is deferred to C5/C6; `TTSSynthesizer` unchanged.

**Runtime evidence.** Unit and fake-helper status tests only.

**Rollback.** Revert type/status projection changes; downstream code does not yet consume the fields.

---

### C3 — Lifecycle state, load gating, and ordinary-load seam

**Scope.** Establish the state gate and, only if D1a permits, the ordinary-load seam; prevent synthesis from recovering `unavailable`.

**Prerequisites.** C2 status model; D1a disposition.

**Design decision.** The low-level `helper.Client` remains transport/state; it must not infer app asset policy. Ordinary ensure-load belongs in an adapter/session seam that owns the `LoadRequest`, because only the application knows provider, engine, model, assets, and cache directory. D1a must be dispositioned before C3 adds any such seam. If D1a permits it, C3 introduces the seam in a narrow, testable form and C7 makes it the production `TTSSynthesizer`; if D1a denies it, C3 is limited to state gating and an internal `Load` on the existing client, and no new adapter/session file is created. C2 projects `Loaded`; C3 owns the lifecycle transitions that set and clear it, with C4–C6 owning the exit, budget, and replacement transitions.

**Likely affected files.**

- `internal/audio/tts/helper/client.go`: expose/consume a narrow load-plan seam; state-aware `Synthesize` outcomes; do not auto-load from `unavailable`; capture `loadedIdentity`/`failedIdentity` on load success/failure; `Loaded` transitions.
- `internal/audio/tts/helper/supervisor.go`: lifecycle state transitions and `Loaded`.
- New internal adapter/session file in the helper package (or a small internal audio adapter file), only if D1a permits it, limited to `Start`, `Load`, `Synthesize`, `State`, and identity capture. No Wails or manager wiring.
- `internal/audio/types/provider.go` only if a small status helper is needed; no `TTSSynthesizer` change.
- Tests for the state gate and, only if D1a permits, the seam; fake-helper modes for unloaded/unavailable.

**Required behavior.**

- `Initialize`-equivalent stores config, no spawn/load.
- `Synthesize` from `unloaded` performs handshake if needed, then one load for the configured identity, then synthesis.
- `Synthesize` from `unavailable` returns `provider_not_loaded` and sends no load or synthesize frame.
- `Synthesize` from `incompatible` returns the incompatibility; no restart.
- `Load` is explicit and updates loaded/failed identity; no automatic retry.
- Load failure sets `unavailable`, `Loaded=false`, precise detail naming `provider_id`/`engine`/`model_identity`, and retains a usable live helper when the failure is asset-specific. The legal live case is `Available=true`, `ProcessAlive=true`, `Loaded=false`, `Ready=false`, `ServiceCapabilitiesKnown=true`; it must remain reachable.
- No background preload from status/catalog/health.
- `Loaded` becomes false when a load begins, fails, unloads, exits, or replacement begins; C3 owns the load/unload transitions, C4 the exit transition, and C6 the replacement transition.

**Tests.** Ordinary unloaded synthesis loads once; unavailable synthesis emits no load/synthesize frame; load failure detail; loaded true only after success; asset-specific failure preserves the live-helper status pair (`ProcessAlive=true`, `ServiceCapabilitiesKnown=true`); no status/catalog load; no automatic load after failure; identity capture on success/failure.

**Negative/failure injection.** Load timeout, load crash, missing asset, repeated load failure, service capability absent (load is required, so handshake should already reject absent load).

**Race/concurrency.** One active load; serialize with existing per-instance lock; no synthesis during loading.

**Exclusions.** Same-identity recovery is not yet enabled; different-identity replacement is C6; budget/backoff/idle is C5.

**Report artifact.** State/load table and the adapter seam contract.

**Independent review.** Confirm ordinary and unavailable paths are distinguishable in code, not only tests.

**Acceptance.** D1a is dispositioned before any seam is added; if permitted, unloaded synthesis succeeds with fake helper through the seam; if denied, the existing client's internal `Load` is the only surface and no new adapter/session file exists; unavailable synthesis fails without a load frame; no production import.

**Runtime evidence.** Real fake-helper process for unloaded success and unavailable no-frame behavior.

**Rollback.** Remove the new seam and state gate; existing client remains inert.

---

### C4 — Exit intent, shutdown, timeout, and cancellation consistency

**Scope.** Make process state authoritative and classify intentional termination versus unexpected exit.

**Prerequisites.** C2 status fields; C3 loaded transitions.

**Likely affected files.**

- `supervisor.go`: explicit exit-intent enum/field (`none`, `shutdown`, `replacement`, `incompatible`, `timeout-cancel`, `idle`); `observeExit` uses intent plus exit code; `SetState` no longer force-sets `ProcessAlive`; aggregate incompatible remains incompatible after expected termination; process state becomes exited; service capability snapshot invalidation.
- `client.go`: timeout/cancel force paths set `stopping`/terminal outcome and never `ready`; `Shutdown` classification; `cleanupRejected` ordering; active request cleared.
- `supervisor_test.go`, `client_test.go`, `interop_test.go`, `regression_test.go`.

**Required behavior.**

- Expected shutdown exit 0 → `absent`; nonzero → `crashed`; grace expiry/kill → `crashed`; no success on transport failure.
- Incompatible handshake termination: aggregate remains `incompatible`, process exited; no budget/restart.
- Unexpected exit from any live state → `crashed`; active request fails with `helper_crashed`.
- Timeout/cancel forced termination: transition through `stopping`; after kill/wait, state is `crashed` under the existing stopping/grace rule, not `ready`; request outcome remains `timeout`/`cancelled`.
- **D6 disposition (applied, not open):** §17.7.2 freezes `stopping` + grace expiry → `crashed`; A2-D9 governs the request outcome as `cancelled`/`timeout`; A2-D8 governs clean exit 0 → `absent`. The intent value is diagnostic and does not override those transitions.
- `ProcessAlive` reflects process state only.
- Service-capability snapshot invalidated on stop/exit/crash/incompatible/replacement.
- No `SetState` can overwrite process state with `alive` after exit.

**Tests.** Exit code 0/nonzero; stop grace; forced kill; incompatible termination; unexpected exit; timeout/cancel post-kill status; no ready after kill; stderr/last error preservation.

**Negative/failure injection.** Wait observer delayed; close/write failure; kill fail; process exits between state read and terminate.

**Race/concurrency.** Exit observer versus timeout/cancel cleanup; `SetState` versus `observeExit`; no double Wait; exact-path cleanup after exit.

**Exclusions.** No restart budget/backoff/idle timers yet.

**Report artifact.** Exit-intent/state table and command traces.

**Independent review.** Verify every termination path has an intent and every unexpected exit has one observer.

**Acceptance.** No path reports ready/alive after an observed exit; shutdown 0/nonzero classification correct; incompatible remains terminal.

**Runtime evidence.** Windows process boundary for job kill/wait and exit codes; required in C8.

**Rollback.** Restore prior observer/state logic; no production caller.

---

### C5 — Restart budget, backoff, and idle close

**Scope.** Add the provisional restart policy and idle shutdown without creating an autonomous restart loop.

**Prerequisites.** C4 intent-aware exits and shutdown classification; D9 observation-reset disposition before acceptance.

**Likely affected files.**

- `supervisor.go`: restart ledger per identity key; budget 3/10 min; backoff 1/2/5 s; `backoff` aggregate state; `ProcessAlive=false` while backoff; reset after 10 min stable readiness; idle timer 120 s.
- `client.go`: explicit request entry points consult budget/backoff; no restart from health/status/catalog observation; explicit status from `absent` remains the §17.7.2 C7 lazy start; idle close calls shutdown, not unload; post-exit classification.
- `interop_test.go`, `supervisor_test.go`, `client_test.go`; injectable clock/timer for tests.
- No production wiring.

**Required behavior.**

- Crash/backoff only on next explicit TTS request; no background loop.
- Restart budget per key counts process restarts. Identity replacement under A2-D7 is exempt. A load-only A2-D12 retry starts no process and therefore consumes no process-restart budget; any A2-D12 retry that restarts a process consumes the §17.11.5 budget and does not inherit the A2-D7 exemption.
- Backoff sequence 1/2/5 s; budget exhaustion → `unavailable`, process exited (`ProcessAlive=false`), `Available=true` while configured/compatible, `Loaded=false`, `Ready=false`, `ServiceCapabilitiesKnown=false`; no automatic restart.
- Stable readiness for 10 minutes resets the budget.
- Idle close: `ready`/`degraded`/`unloaded`, no active work, 120 s → `shutdown`; exit 0 absent, nonzero crashed; never `unloaded` merely because stopped; no `unload` frame.
- Idle close releases provider resources.
- Idle-deadline reset behavior is **D9** and is not frozen here. The contract supplies the 120 s trigger, but not whether status/catalog/health observation resets the deadline. D9 asks only that question. The separate question whether an explicit status request may start an `absent` helper is not open: §17.7.2 makes an explicit status request from `absent` a lazy-start trigger. If the product does not want status polling to start a helper, that is an application-layer polling policy—whether the shell issues status requests—not a C5 lifecycle exception. C5 must not assume an observation-reset answer; D9 must be dispositioned before C5 acceptance.

**Tests.** Three restarts then unavailable; 1/2/5 sequence with fake clock; budget reset after stable readiness; identity replacement does not consume budget; a load-only A2-D12 retry consumes no process-restart budget; an A2-D12 retry that restarts a process does consume it and has no replacement exemption; idle close at deadline, cancellation of timer on work, nonzero exit, no unload frame; process counts. D9 observation-reset tests are added only after D9 is dispositioned; status-from-`absent` lazy start is contract-owned and verified at C7, not tested as a D9 choice.

**Negative/failure injection.** Clock jumps, timer race with synthesis start, shutdown failure, repeated crash loop, budget exhausted then update-gated reset.

**Race/concurrency.** Idle timer versus workMu; backoff timer versus explicit request; one process start at a time.

**Exclusions.** No production adapter; no native measurements; numeric values remain provisional, not performance claims.

**Report artifact.** Budget/backoff table, timer test log, process-count evidence.

**Independent review.** Confirm no autonomous retry and no budget reset through status/catalog.

**Acceptance.** Budget permits three restarts within 10 minutes; a fourth attempt exhausts the budget; backoff sequence; idle shutdown classification; the D9 observation-reset disposition is recorded and tested; status-from-`absent` lazy start is verified at C7 as §17.7.2 behavior and is not gated on D9; A2-D7 exemption and A2-D12 non-exemption proven by tests.

**Runtime evidence.** Real fake-helper process with short internal overrides; Windows kill/wait boundary in C8.

**Rollback.** Remove timer/ledger; C4 exit semantics remain.

---

### C6 — Identity matching, unavailable recovery, and replacement

**Scope.** Implement the exact identity tuple, same-identity recovery, and A2-D7 replacement ordering.

**Prerequisites.** C3 load/state behavior; C4 termination; C5 budget.

**Likely affected files.**

- `client.go`: load identity fields; `Load` transition rules; state-gated explicit load; replacement sequence; no active synthesis check; `internal_error` on lock invariant violation.
- `supervisor.go`: identity/state reset, service capability invalidation at replacement, budget interaction.
- `interop_test.go`, `client_test.go`, `regression_test.go`; fake helper modes for shutdown errors, nonzero exit, transport failure.
- No production adapter.

**Required behavior.**

- Identity tuple exactly `provider_id`, `engine`, `model_identity`.
- `unavailable` + explicit `load`:
  - same identity, live helper, helper-side unloaded/accepting, valid config → one new load; failure remains unavailable; success normal ready.
  - different identity → A2-D7 stop-and-replace; never send differing load to old helper.
  - helper not live or config/helper invalid → no in-place recovery; update-gated reset only.
- `Synthesize` while unavailable → `provider_not_loaded`, no retry.
- A2-D7 sequence: assert no active synthesis (else `internal_error`, no cancel); shutdown; grace; close stdin; wait; KillTree if needed; observe Wait; start fresh only after exit.
- Shutdown error/response classification; transport failure/nonzero → no replacement.
- Replacement exempt from crash/restart budget; failed replacement start → `helper_start_failed`, no automatic backoff, budget unchanged.
- A2-D12 recovery does **not** inherit the replacement exemption. A same-identity retry that loads on a live helper starts no process and consumes no process-restart budget; a retry that restarts a process is governed by §17.11.5. **D7 is settled by §17.11.5:** after a successful in-place recovery, readiness resumes and the ledger resets only after ten continuous minutes of stable readiness; until then prior crashes remain counted. No owner decision or separate C6 gate remains. The contract's process-restart accounting is not reopened.
- `Loaded=false`, service caps unknown during replacement; old identity invalid at start; new identity only after load.

**Tests.** Same identity success/failure; different identity; missing unload; shutdown error; transport write failure; nonzero exit; active synthesis invariant; A2-D7 replacement exemption; load-only A2-D12 retry consumes no process-restart budget; A2-D12 process restart consumes it without exemption; failed replacement start; no load frame to old process; service snapshot invalidation.

**Negative/failure injection.** Shutdown hangs; exit race; kill failure; start failure after old exit; repeated same-identity failure; identity differs only in engine/model.

**Race/concurrency.** Per-instance lock; no active synthesis; wait observer; no start before exit; timer/budget lock ordering.

**Exclusions.** No production caller; no UI; no native provider.

**Report artifact.** Identity/replacement state table and process traces.

**Independent review.** Walk every stop/replace branch against A2-D7 and every recovery branch against A2-D12.

**Acceptance.** All identity/replacement cases pass with real fake helper; no differing-identity load to old process; budget behavior exactly as specified.

**Runtime evidence.** Windows process boundary for replacement shutdown/kill; C8.

**Rollback.** Remove identity/recovery rules; C4/C5 remain.

---

### C7 — Production adapter and explicit load reachability

**Scope.** Implement the production adapter without changing `TTSSynthesizer`. An internal explicit-load path is designed only if D1a permits it; otherwise C7 carries the adapter and records reachability as blocked rather than inventing a surface.

**Prerequisites.** C1–C6 stable; D1a already decided before C3; D1b, D2, D3, D10(a)/D10(b) with concrete results, and the remaining owner decisions.

**Likely affected files.**

- New adapter file in `internal/audio/tts/helper` or a small internal audio package: `TTSSynthesizer`, `TTSProviderInspector`, `GetConfig`.
- `internal/audio/pipeline.go`: construct adapter for a helper-backed engine; integrate with `ttsEngines` and `loadTTSForRequest`; keep cache key and character precedence.
- `internal/audio/manager_factory.go`: provider-neutral factory gate; no native changes.
- `app-audio.go`: invoke the explicit loader only from the existing explicit settings path; no UI.
- `app-tts-catalog.go`: static fallback; do not start helper from `GetTTSVoiceCatalog`.
- `internal/audio/types/provider.go`: optional internal loader interface only; no `TTSSynthesizer` change.
- Tests: adapter unit, pipeline integration, manager integration, cache/character compatibility; no native assets.

**Adapter behavior.**

- `Initialize`: validate/store config; no spawn/load.
- `Synthesize`: `Start`; from `unloaded` load with adapter-owned plan; from `unavailable` return `provider_not_loaded`; from ready/busy/degraded synth; validate output; preserve acceptance/cancellation.
- `SynthesizeStream`: `unsupported_capability`.
- `GetVoices`: last successful catalog snapshot; never spawn/block.
- `SupportsSSML`: false.
- `Release`: idempotent shutdown/kill; no new work.
- `GetConfig`: stored config for existing identity resolver.
- `ProviderStatus`: may start/handshake; no load.
- `ProviderCatalog`: may start/handshake/catalog; no load.
- `GetTTSVoiceCatalog`: static fallback; no helper start (O-1 remains open).

**Explicit load reachability (D1a and D1b, not decided here).**

- **D1a — surface permissibility.** The adopted §17.12.1 table has no `load`, `unload`, or `cancel` entry point, while A2-D12 requires an explicit `load` request for recovery. Determine whether a narrow internal, non-Wails loader seam outside that enumerated surface is permitted, or whether the adoption record needs a clarification. This proposal does not assume the answer and does not add a public API.
- **D1b — trigger and adapter lifetime.** If D1a permits the seam, identify the existing explicit product action that invokes it and whether the adapter survives that action for a genuine same-identity in-place retry. Do not invoke loading from `Synthesize`, status, health, catalog, `GetVoices`, or startup.
- If the existing Apply path rebuilds the generation or destroys the adapter, do not describe that as in-place recovery. Preserving the adapter or adding an internal retry command requires D1b authorization.
- No recovery UI and no public Wails binding are part of this design.

**Global/character ownership and process cap.**

- The cap counts **actual helper child processes**, not adapters, character voices, engine labels, or requests. Multiple requests or character voice configurations served by one helper process count once. The frozen cap remains two helper processes per app session (§17.8 line 1857; §17.19.1).
- The adapter does not own the count. A process registry in the supervisor/proc layer is keyed by the real process/job handle, not by adapter identity. An entry exists only while a child process is alive or its exit is unobserved; it is removed only after `Wait()` is observed. An adapter maps to zero or one process handle. Reusing an existing process adds no entry and starts no process. Admission is checked before `proc.Start`, and the cap is never derived from the adapter map, engine label, or character count.
- Character voice/engine remains request-scoped; explicit engine failure does not fall back.
- Helper-backed model identity must differ from the old Python bridge identity; cache key unchanged.
- A third concurrent helper-backed engine is rejected before process creation. The cap is frozen, but no §17.9 entry exactly defines this admission rejection: `busy` is wire-origin and means one active non-control operation; `helper_start_failed` means a process/job/pipe start failed; `unsupported_capability` is optional-capability-only. D3 therefore remains a narrow contract gap; do not invent or repurpose a code.

**Tests.** Interface conformance; Initialize no spawn; status/catalog no load; unloaded synth loads; unavailable synth no load; explicit loader same identity; different identity replacement; cache/character voice compatibility; no helper start from `GetTTSVoiceCatalog`; process cap counted over real child processes. Include: two requests sharing one adapter/process count once; two adapters sharing one process count once; two engine labels mapped to one process count once; an adapter with no process counts zero; an exited-and-reaped process frees its slot; a third real process is refused before start (error surface contingent on D3).

**Exclusions.** Native Kokoro/G2P/ORT/assets; UI; O-1/O-3/O-4 activation; B12a/B12b.

**Report artifact.** Adapter contract, caller map, process cap evidence, cache/identity compatibility report.

**Independent review.** Confirm no eager loading from inspection paths; confirm `TTSSynthesizer` unchanged; confirm recovery is explicit and distinguishable.

**Acceptance.** Pipeline can construct and exercise the adapter against fake helper; no production activation until D2; explicit loader exists only if D1a permits and is reachable only by the D1b intentional action; no status/catalog load; the cap is enforced by actual process count and cannot drift with adapters or character voices; no D3 error is claimed until the error-surface gap is dispositioned.

**Runtime evidence.** Fake-helper end-to-end; Windows process tests; native provider remains out.

**Rollback.** Keep adapter behind an internal default-off activation gate; revert wiring without touching C1–C6.

---

### C8 — Adopted-contract test closure and separately authorized runtime verification

**Scope.** Close the evidence matrix and run the checks that source inspection cannot establish. This stage is not implementation and requires separate execution authorization.

**Prerequisites.** C1–C7 complete and independently reviewed.

**Work.**

- Full adopted-contract test matrix; replace all superseded five-capability expectations.
- Focused race tests; full suite; known `internal/audio` cache flake isolated from helper claims.
- Real fake-helper process tests for handshake, direct-unload absence, cancellation absence (separately), unloaded load, unavailable no-recovery, identity replacement, budget/backoff/idle, shutdown exit codes, timeout/cancel status.
- Windows boundary: suspended start/job, kill/wait, final symlink/junction, shutdown nonzero, forced kill.
- Production adapter fake-helper integration; no native models.
- Evidence report with commands, exit codes, logs, hashes.

**Binding evidence rule.**

1. A row labeled Windows-specific — including suspended start/job ownership, `KillTree`/`Wait`, nonzero-exit classification, forced kill, and final symlink/junction handling — closes only through execution on Windows against a real child process.
2. A fake-helper, cross-platform, or CI result may support such a row but cannot close it. If the Windows run did not occur, record the row as **Unestablished**, not “indirect” and not “passed.”
3. Every C8 evidence row must name its environment (fake-helper, Windows, CI, or native-provider) and its process boundary (real child, in-process fake, or none), plus command and result. “Passed” without an environment and boundary is not evidence.
4. The C1 statement that Windows behavior is inherited is scoped to C1 only. It may not close C4–C8 rows.
5. Fake-helper protocol tests, Windows real-process tests, and native-provider tests remain separate evidence classes. Native-provider validation is out of scope; its absence is reported as unestablished, not replaced by fake-helper results.

**Exclusions.** No native provider, model download, licensing decision, or performance claim.

**Acceptance.** Every matrix row is direct or recorded as **Unestablished**; fake-helper or CI results may be recorded only as supporting evidence and never close a Windows-specific row; every row names its environment and process boundary; no contradictory test remains. C4, C5, and C6 produce the code; C8 produces the binding evidence.

## 6. Dependency graph and critical path

```text
C1 Capability
  |
  v
C2 Status projection
  |
  v
C3 Lifecycle/load seam
  |
  v
C4 Exit/shutdown/cancel truth
  |
  v
C5 Budget/backoff/idle
  |
  v
C6 Identity/recovery/replacement
  |
  v
C7 Production adapter + explicit load
  |
  v
C8 Contract closure + runtime verification
```

**Critical path:** C1 → C2 → C3 → C4 → C5 → C6 → C7 → C8. No later stage can be honestly accepted while an earlier contract invariant is still contradicted.

**Safe parallel work.** Documentation/test-fixture inventory for C2 may start after C1's type decisions. C4 test scaffolding may be prepared after C2, but C4 and C5 share the supervisor state model and should not merge. C7 adapter scaffolding can be reviewed after C3 but must not be wired before C6.

**Do not combine.**

- C1 and C2: wire semantics versus application projection; different failure modes.
- C2 and C3: status truth versus load behavior; a status-only patch can mask lifecycle bugs.
- C3 and C6: ordinary loading versus recovery/identity; A2-D12 specifically distinguishes them.
- C4 and C5: exit truth must precede timers/budgets; otherwise a timer can launder a crashed process into absent.
- C6 and C7: contract behavior must be testable before production reachability.
- C7 and C8: adapter wiring and runtime evidence are separate authorizations.

## 7. Affected source and interface inventory

| Area | File/type/function | Change class | Stage |
|---|---|---|---|
| Service policy | `helper/protocol.go` mandatory set, `validateHello`, `validateCapabilities` | Change | C1 |
| Handshake | `helper/client.go` `Start`, `helloHandshake`, `requireServiceCapability` | Change | C1 |
| Direct unload gating | `helper/client.go` `Unload` | Classification/no-frame only | C1 |
| Cancellation gating | `helper/client.go` supervisor cancellation path / exported `Cancel` | Path: skip-frame/no-error; direct: `unsupported_capability` | C1 |
| Capability wire | `helper/protocol.go` `capability` | Preserve `{id,version}` | C1 |
| Fake service set | `fakehelper/server.go` `serviceCapabilities`, `providerCapabilities` | Change | C1 |
| App status | `types/provider.go` `TTSCapability`, `TTSProviderStatus` | Change | C2 |
| Status projection | `helper/supervisor.go` `MarkHandshake`, `Status`, `capabilityList`, `processAlive`, `providerAvailable`, `providerReady` | Projection only in C2; transition truth in C4–C6 | C2, C4 |
| Loaded writer | `helper/supervisor.go`, load seam | C2 projects from snapshot; C3 owns load/unload transitions | C3, C4, C6 |
| Process truth | `helper/supervisor.go` `SetState`, `observeExit`, `Terminate` | Change | C4 |
| Load/state | `helper/client.go` fields, `Load`, `Synthesize`; new adapter/session seam only if D1a permits | Change/add | C3, C6 |
| Identity | `helper/client.go` load identity fields; adapter config | Add | C3, C6 |
| Budget/idle | `helper/supervisor.go`, `helper/client.go` | Add | C5 |
| Typed errors | `helper/errors.go` `Failure`, `failureFromWire` | Preserve/extend detail | C1, C2 |
| Adapter | new helper/audio adapter | Add | C7 |
| Process-cap registry | supervisor/proc process registry | Add; counts real child processes only | C7 |
| Explicit loader seam | internal non-Wails loader interface | D1a gate before C3; add only if permitted; D1b trigger | C3/C7 |
| Pipeline | `internal/audio/pipeline.go` `initTTS`, `loadTTSForRequest` | Change | C7 |
| Factory | `internal/audio/manager_factory.go` | Change | C7 |
| App explicit path | `app-audio.go` `applyAudioSettings` | Narrow hook | C7 |
| Static catalog | `app-tts-catalog.go` `GetTTSVoiceCatalog` | Preserve/guard | C7 |
| Stable interface | `types/interfaces.go` `TTSSynthesizer` | **Unchanged** | — |
| Inspector | `types/provider.go` `TTSProviderInspector` | **Unchanged shape** | C2/C7 |
| Process ownership | `internal/proc/job_windows.go` | **Unchanged** | — |
| Output ownership | `helper/output.go`, `fakehelper/publish.go` | **Unchanged** | — |
| Wire framing | `helper/protocol.go`, `fakehelper/protocol.go` | **Unchanged** | — |

## 8. Acceptance evidence matrix

Matrix cells record where a requirement is exercised; they are not closure claims. The C8 report supplies the required environment and process-boundary fields for every row.

| Requirement | C1 | C2 | C3 | C4 | C5 | C6 | C7 | C8 |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| Required-two/optional-three | Direct | — | — | — | — | — | — | Recheck |
| Missing required → incompatible | Direct | — | — | Recheck | — | — | — | Runtime |
| Direct unload absent → `unsupported_capability`, no frame | Direct | — | — | — | — | — | Recheck | Runtime |
| Direct `Cancel` absent → `unsupported_capability`, no frame | Direct | — | — | — | — | — | Recheck | Runtime |
| Cancellation path absent → no frame, no `unsupported_capability` | Direct | — | — | Recheck | — | — | Recheck | Runtime |
| Unknown capability knowledge → no optional-operation spawn; D10(a)/D10(b) require concrete results before acceptance | Direct | — | — | — | — | — | Recheck | Runtime |
| Duplicate extension / invalid descriptor → `protocol_incompatible` | Direct | — | — | — | — | — | Recheck | Runtime |
| Scope separation | Direct | Direct | — | — | — | — | — | Recheck |
| `detail.missing` | Direct | Direct | — | — | — | — | — | Recheck |
| Status fields and single-source table | — | Direct | — | — | — | — | Recheck | Runtime |
| Both `unavailable` flag sets | — | Direct | Recheck | — | Direct | Recheck | Recheck | Runtime |
| ProcessAlive from process state only | — | Direct | — | Direct | Direct | Direct | Recheck | Runtime |
| Loaded projection | — | Direct | — | — | — | — | Recheck | Runtime |
| Loaded transitions | — | — | Direct | Direct | — | Direct | Recheck | Runtime |
| Ordinary unloaded ensure-load | — | — | Direct | — | — | — | Direct | Runtime |
| Unavailable no recovery | — | — | Direct | — | — | Direct | Direct | Runtime |
| Identity tuple | — | — | Capture | — | — | Direct | Direct | Runtime |
| Same-identity recovery / D7 accounting | — | — | — | — | Budget | Direct | Direct | Runtime |
| A2-D7 replacement | — | — | — | Exit | Budget | Direct | Direct | Runtime |
| Budget/backoff | — | — | — | Intent | Direct | Direct | — | Runtime |
| Idle close | — | — | — | Shutdown | Direct | — | — | Runtime |
| Exit classification | — | — | — | Direct | Direct | Direct | — | Windows |
| Timeout/cancel status | — | — | — | Direct | — | — | Direct | Runtime |
| No production eager load | — | — | — | — | — | — | Direct | Runtime |
| `TTSSynthesizer` unchanged | — | — | — | — | — | — | Compile | Review |
| Process cap counts actual child processes | — | — | — | — | — | — | Direct | Runtime |
| Process/file ownership | — | — | — | Preserve | Preserve | Preserve | Recheck | Windows |
| Windows-specific row closed only on Windows real child | — | — | — | — | — | — | — | Windows |
| Evidence environment/process boundary recorded | — | — | — | — | — | — | — | Direct |
| No native provider | — | — | — | — | — | — | Excluded | Excluded |

## 9. Same-day execution proposal

This section describes what could be attempted **after separate implementation authorization**; it is not a claim that any of it is authorized or complete today.

### 9.1 Realistic same-day boundary

- **Can be implemented and fake-helper tested:** C1 and C2. They are the highest-value, narrowest stages, and they remove the two most direct contract contradictions without wiring production. C1 cannot be accepted until D10(a) and D10(b) each have a concrete caller-facing disposition; a deferral does not satisfy that gate; C2 has no open design gate.
- **Can begin only after C1/C2 review and D1a disposition:** C3. It is large enough that claiming it complete without review would weaken the acceptance boundary.
- **Cannot be honestly completed today:** C4–C6 as a unit. They require Windows process-boundary evidence, restart/idle timing, and identity replacement traces. Do not compress them into a “lifecycle” patch. C5 also waits for the D9 observation-reset disposition.
- **Cannot be completed today:** C7 production activation and C8 runtime closure. They depend on C1–C6, D1a already decided before C3, D1b/D2, and the concrete D3/D10(a)/D10(b) dispositions.

### 9.2 Work classification

| Work | Can use fake helper | Requires Windows boundary | Requires native model/assets |
|---|---:|---:|---:|
| C1 capability/classification | Yes | No (C1 only; inheritance does not extend to C4–C8) | No |
| C2 status projection | Yes | No | No |
| C3 load seam/state | Yes | No | No |
| C4 exit/shutdown | Yes | Yes for job/kill/exit | No |
| C5 budget/backoff/idle | Yes | Yes for real process/kill | No |
| C6 identity/replacement | Yes | Yes for replacement shutdown | No |
| C7 adapter/fake integration | Yes | Yes for child process | No |
| C8 closure | Yes | Yes | No (native excluded) |

### 9.3 Definition of “TTS helper integration complete”

It means:

- the production adapter implements `TTSSynthesizer` without changing that interface;
- manager/pipeline can select a helper-backed provider;
- ordinary `unloaded` synthesis loads through the adapter;
- `unavailable` does not recover on synthesis;
- one intentional explicit load/recovery path is reachable;
- identity/replacement, restart budget/backoff, idle close, shutdown, and status behavior satisfy C1–C7;
- fake-helper contract tests and Windows real-child process tests pass, with each Windows-specific row established by a Windows run rather than by inference.

It does **not** mean:

- Kokoro/G2P/ORT or any native provider is production-ready;
- O-5/O-6 licensing/assets/build are resolved;
- O-4 default backend switch is decided;
- TTS performance or model quality is established;
- B12a/B12b are resolved.

## 10. Unresolved decisions and out-of-scope items

The table distinguishes contract-settled rules, contract-surface gaps, and genuine product/architecture choices. D6 is applied in C4; D7 is settled by §17.11.5; D4 is folded into D2/D5; D1 is split into surface and trigger decisions; D9 records only the observation-driven idle-deadline reset decision; §17.7.2 settles status-from-`absent` lazy start; D10 explicitly records two separate unknown-snapshot direct-operation questions, D10(a) for direct `unload` and D10(b) for direct `Cancel`; each requires a concrete caller-facing result before C1 acceptance, and the supervisor cancellation path is A2-D9-settled and outside D10.

| ID | Classification and decision | Governing basis | Needed by |
|---|---|---|---|
| D1a | **Contract-surface gap.** May the implementation add a narrow internal, non-Wails loader seam outside §17.12.1's enumerated adapter behavior, or does explicit load reachability require a contract clarification/adoption-record decision? §17.12.1 also provides no direct `unload` or `cancel` entry point, so A2-D7's direct unload and A2-D9's cancellation path rely on the existing internal client surface. No public API is assumed. | §17.12.1 (no `load`, `unload`, or `cancel` entry point) and A2-D12 (explicit `load` is the only recovery trigger) | Before C3 introduces any loader seam; C3 is not authorizable until resolved |
| D1b | **Genuine product/architecture decision.** Which existing explicit product action invokes the permitted loader, and does the adapter survive that action for genuine same-identity in-place recovery? | A2-D12; the existing settings Apply/rebuild path | Before C7 production reachability |
| D2 | **Genuine product/architecture decision.** How a helper-backed engine is activated in settings without changing `TTSSynthesizer`, and its relationship to O-4. | §17.12.1; O-4 remains open | Before C7 production activation |
| D3 | **Frozen cap, open caller-facing error.** The cap is two actual helper processes, but §17.9 contains no entry that exactly defines admission rejection. `busy` is wire-origin and means one active non-control operation; `helper_start_failed` means process/job/pipe start failure; `unsupported_capability` is optional-capability-only. Do not invent or repurpose a code. The precise gap is: which existing code, if any, the adapter returns when rejecting a third process before `proc.Start`. | §17.8 line 1857 and §17.19.1 (frozen cap); §17.9 registry | Before C7 can claim a contract-compliant cap response; C7 may implement the counting/admission boundary without asserting the error |
| D4 | **Folded, not independent.** Whether the inspector has a production consumer is wiring implied by D2/D5; it is not a separate policy decision. | §17.12.1; O-1 static fallback | Resolved with D2/D5, no separate gate |
| D5 | **Genuine product/architecture decision.** Renderer exposure shape and O-1 picker behavior. Keep both static until decided. | §17.12; O-1 | Before any production UI integration |
| D6 | **Settled and applied.** Timeout/cancel forced termination is `stopping`, then grace expiry → `crashed`; request outcome is `cancelled`/`timeout`; clean shutdown exit 0 → `absent`. The exit intent is diagnostic only. | §17.7.2 (`stopping` + grace → `crashed`); A2-D9; A2-D8 | Already applied in C4; no open decision |
| D7 | **Settled by citation; no owner decision.** A load-only retry starts no process and consumes no process-restart budget. A retry that restarts a process consumes the §17.11.5 budget and never inherits A2-D7's exemption. After successful in-place recovery, readiness resumes; the ledger resets when ten continuous minutes of stable readiness have elapsed, so prior crashes remain counted until then. The former alternatives describe the same rule. | A2-D12; §17.11.5 (3/10 minutes; reset after 10 minutes stable readiness) | Closed; C5 implements the generic stable-readiness clock and C6 applies it |
| D8 | **Deferred contract issues.** B12a hard-link vs atomic rename; B12b extension-error wording. | A2-D11 | Not resolved in this plan |
| D9 | **Genuine product decision.** Does status/catalog/health observation reset the 120 s idle deadline? The contract supplies the trigger but not this reset rule. C5 must not assume either answer. The separate status-from-`absent` start behavior is not D9: §17.7.2 already makes an explicit status request a lazy-start trigger, and any desire to suppress poll-driven startup is an application-layer polling policy. | §17.8; §17.7.2 (settled status start, not part of D9); A2-D8 | Before C5 acceptance |
| D10 | **Owner decision / contract gap, explicitly expanded and split into two separately answerable direct branches.** **D10(a):** unknown-snapshot direct `unload` result. **D10(b):** unknown-snapshot direct `Cancel` result outside the supervisor path. D10 originated as the unload branch; direct `Cancel` is a separate direct-operation branch under A2-D1 and is not assumed to share unload's disposition. A2-D9 settles only the supervisor cancellation path, which is outside D10. The owner may select the same or different result for D10(a) and D10(b); neither is frozen. C1 implements known-present and known-absent only. Candidate directions are an explicit semantic extension of an existing code, treating an unverified capability as absent, or another concrete caller-facing result that preserves the no-frame/no-spawn rule; none is selected. Deferral may schedule when the owner answers, but it is not a behavioral resolution and cannot satisfy C1 acceptance: each branch needs a concrete caller-facing result. No optional operation may spawn. | §17.9; A2-D1; A2-D7; A2-D9 (supervisor path only, settled); §17.7.2 | Before C1 acceptance (both branches, concrete results) |
| O-1 | May the picker start a helper for catalog? | §17.12.1; O-1 remains open | Keep static fallback; no eager start |
| O-2 | WAV retention/cleanup ownership | §17.10; O-2 remains open | Do not add cache deletion |
| O-3 | `auto` fallback behavior | A2-D10; O-3 remains open | No silent fallback |
| O-4 | When Kokoro becomes helper-backed | §17.19; O-4 remains open | Native provider out of scope |
| O-5 | ORT version/DLL/CRT licensing | A2-D11; O-5 remains open | Native provider out of scope |
| O-6 | Kokoro/G2P assets/build/manifest | A2-D11; O-6 remains open | Native provider out of scope |
| O-7 | One binary vs separate provider executables | A2-D11; O-7 remains open | Post-v1; do not decide |

**Out of scope:** native provider implementation, G2P, Kokoro, ONNX, CGO, model downloads, packaging, UI recovery, new Wails public API, cache deletion, performance claims, contract edits.

## 11. Rollback and containment strategy

- **C1–C6:** changes remain inside `internal/audio/tts/helper`, `internal/audio/types`, and tests; no production import. Revert the stage's files and tests.
- **C2 type changes:** compile-time breakage is intentional and contained; no production consumer yet. Update literals/tests in the same stage.
- **C4/C5 state/timer changes:** keep old observer path until tests pass; use injectable clock/timer, no sleeps in production logic.
- **C6 replacement:** no production call path; fake-helper tests only.
- **C7:** adapter activation behind an internal default-off gate; no Wails setting until D2. Revert wiring without touching C1–C6.
- **C8:** tests/evidence only; no runtime artifacts in the repository.
- **No migrations or user data changes.** No cache key change. No output cleanup broadened. No process ownership change.

## 12. Design-only statement

This document is a design and sequencing proposal. It is **not** implementation authorization.

No production source, test, fixture, protocol, contract, amendment, prior report, or audit was modified. No helper was wired into the manager or pipeline. No recovery API or UI was added. No native provider, G2P, Kokoro, ONNX, CGO, or dependency was added. No build, test, replay, benchmark, process launch, asset download, commit, push, or stage was performed.

Stage 3 implementation remains unauthorized. Each stage requires explicit, stage-specific authorization, and C8 requires separate runtime-execution authorization.
