# INVESTIGATION — HOWL TTS Stage 3 readiness against adopted §17.19

**Date:** 2026-11-07  
**Mode:** READ-ONLY INVESTIGATION  
**Governing contract:** `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md`  
**Disposition:** **NOT READY** as a contract-compliant Stage 3 baseline.

This report compares the accepted Stage 2 implementation, as it exists in the current source tree, with the subsequently adopted governing contract, especially §17.19. Stage 2 acceptance is treated as evidence about a real-process foundation, not as proof of compliance with the later amendment.

No production source, test, fixture, protocol, contract, amendment, prior report, or integration file was changed. No build, test, process launch, replay, benchmark, asset download, commit, push, or stage was performed. Runtime behavior is therefore not claimed unless a cited prior execution record establishes it.

## 1. Executive disposition

### 1.1 Precise meaning of “not ready”

The accepted Stage 2 foundation is structurally sound at the process and protocol-mechanics level:

- It starts a parent-owned helper with Windows job-object containment and suspended startup.
- It has one blocking waiter and an exit observer.
- It has strict frame decoding, request/terminal validation, cancellation handling, and output containment/validation.
- The fake helper is a real child process, not an in-memory substitute.

It is **not** a Stage 3 baseline under adopted §17.19 because several current source behaviors directly contradict the adopted policy, not merely because tests are missing:

1. The helper handshake rejects absence of all three optional service capabilities, and the parent validates five mandatory service capabilities. A2-D1 requires exactly two required and three optional service capabilities.
2. Semantically invalid handshakes and capability descriptors are reported as `protocol_violation` and leave the child alive, rather than becoming terminal `protocol_incompatible` without restart.
3. Status types cannot represent the required `Supported`/`Detail`, `Loaded`, `ServiceCapabilities`, or `ServiceCapabilitiesKnown` projection.
4. `Synthesize` does not implement the ordinary unloaded ensure-handshake/load path; a real helper returns `provider_not_loaded`.
5. Recovery from `unavailable` is not state-gated, does not require matching `provider_id`/`engine`/`model_identity`, and is unreachable through any production caller.
6. There is no restart budget, no restart backoff, no idle close, and no `backoff` state assignment.
7. Shutdown classification can turn an unexpected nonzero exit into `absent`, and timeout/cancel cleanup can report a ready/live status after termination.
8. No production path imports or wires the helper. The helper does not implement `TTSSynthesizer`, and no explicit production `load` path exists.

These are source-level contract contradictions. Changing tests alone cannot make the implementation ready.

### 1.2 Readiness verdict by question

| Question | Verdict |
|---|---|
| Is the Stage 2 real-process foundation usable as a starting point? | Yes, conditionally; process ownership and protocol mechanics are substantial. |
| Does it satisfy adopted §17.19 as-is? | No. |
| Can Stage 3 be authorized on this baseline without remediation? | No. |
| Is there any production caller or end-to-end helper integration? | No. |
| Does `TTSSynthesizer` remain unchanged? | Yes. |
| Are the Stage 2 report’s five-mandatory-service statements still governing? | No; superseded by adopted A2-D1. |
| Is Stage 3 implemented or authorized by this report? | No. |

## 2. Verified governing contract and input hashes

All hashes below were recomputed from the current files with `sha256sum`. Hashes are calculated evidence; the contract/input texts are observed evidence.

| Artifact | Verified SHA-256 | Result |
|---|---|---|
| Governing architecture | `ec3c789e16c29adbfb35a269d9fa16e2020f87f06c909d9e4db6bba55c25b50f` | Matches authorization |
| Adoption record | `9da36691054df40d98c772f7ba9643bc959416e39a7f2d30681d9404cd16f043` | Matches authorization |
| Accepted amendment | `dc5852adbc70641dbe8f233851eaec44266330fce197b78a0ab2641e7f73f784` | Matches authorization |
| Stage 2 implementation report | `ec85b2b2c3760542699bcd406f1f3b6b3eb2cdc0165b9d6a3b21977a06741227` | Matches authorization |
| Stage 2 correction report | `c9f498f61e0444ab9a66b347f466da7ac456e5afb0a30534621ab453de83857f` | Matches authorization |
| Stage 2 independent audit | `e73cef36491d2d77235c31cc051c3e756ac73d78389ca2d0f90967b9dee18755` | Actual audit file hash |
| Stage 2 capability-correction audit | `d6b2349e030ac8867d236d571a5d3e309eb6c35dfde2a7e5b6586a51cc8df969` | Actual correction-audit hash |
| Stage 1 report | `1e8e6553a6e5c2f71f391940fe6ef670fbcf73d4c9381b23b066eeb639705d0a` | Historical input |
| Stage 1 correction | `5308ad5700b586ebc9c7764d654fa5bd770b128b52820ca8c9184a1b38a3e350` | Historical input |
| Stage 1 audit | `03ed0a4676089a0a3086bbbfeaf36f9159e6cca3bd9ec9b6d63121cfa68680f0` | Historical input |

The adversarial Stage 2 audit names subject hash `adb96b06…`, which is the pre-correction Stage 2 report and does not match the current corrected report `ec85b2b2…`. It is therefore a historical audit of an earlier report state, not current compliance evidence. The capability-correction audit is the audit that corresponds to the corrected report and current Stage 2 artifacts.

The Stage 2 report’s declared resulting hashes for all 17 manifest artifacts matched the current working tree. The Stage 1/Stage 2 baseline hashes for `supervisor.go`, `output.go`, `discovery.go`, helper path checks, and `job_windows.go` also matched. The corrected `client.go` and `protocol.go` hashes, and the fake-helper/subprocess artifacts, matched the correction report.

## 3. Scope and source files inspected

The inspection followed actual source and call paths rather than report claims.

### 3.1 Contract and review records

- `docs/review/DESIGN-howl-tts-helper-architecture-2026-11-06.md` — especially §17.1 item 11, §17.7.1–§17.7.4, §17.8, §17.9, §17.10.2/3/6, §17.11.1–§17.11.6, §17.12/§17.12.1, §17.17, §17.18, and §17.19 (A2-D1–A2-D12).
- `docs/review/ADOPTION-howl-tts-capability-policy-c1-17.19-2026-11-07.md`.
- `docs/review/AMENDMENT-DRAFT-howl-tts-capability-policy-c1-17.19-2026-11-07.md`.
- `docs/review/IMPLEMENTATION-howl-tts-helper-fake-process-stage2-2026-11-06.md`.
- `docs/review/CORRECTIVE-howl-tts-helper-capabilities-no-overwrite-2026-11-06.md`.
- `docs/review/AUDIT-howl-tts-helper-stage2-adversarial-2026-11-06.md`.
- `docs/review/AUDIT-howl-tts-helper-capability-correction-2026-11-06.md`.
- Stage 1 report, correction, and audit records identified in §2.

### 3.2 Source

- `internal/audio/tts/helper/protocol.go`
- `internal/audio/tts/helper/client.go`
- `internal/audio/tts/helper/supervisor.go`
- `internal/audio/tts/helper/output.go`
- `internal/audio/tts/helper/discovery.go`
- `internal/audio/tts/helper/pathcheck_windows.go`
- `internal/audio/tts/helper/pathcheck_other.go`
- `internal/audio/tts/fakehelper/server.go`
- `internal/audio/tts/fakehelper/publish.go`
- `internal/audio/tts/fakehelper/protocol.go`
- `cmd/howl-tts-helper-fake/main.go`
- `internal/audio/types/interfaces.go`
- `internal/audio/types/provider.go`
- `internal/audio/pipeline.go`
- `internal/audio/manager_factory.go`
- `app-tts-catalog.go`
- `internal/proc/job_windows.go`
- `internal/proc/job_other.go`
- relevant audio manager/catalog type definitions reached from those files.

### 3.3 Tests and fixtures

- `internal/audio/tts/helper/client_test.go` — 16 tests.
- `internal/audio/tts/helper/protocol_test.go` — 10 tests.
- `internal/audio/tts/helper/supervisor_test.go` — 5 tests.
- `internal/audio/tts/helper/discovery_test.go` — 4 tests.
- `internal/audio/tts/helper/golden_test.go` — 1 test plus golden fixture.
- `internal/audio/tts/helper/interop_test.go` — 12 tests.
- `internal/audio/tts/helper/interop_process_windows_test.go`
- `internal/audio/tts/helper/interop_process_other_test.go`
- `internal/audio/tts/helper/regression_test.go` — 3 tests.
- `internal/audio/tts/helper/reparse_windows_test.go` — 2 tests.
- `internal/audio/tts/fakehelper/server_test.go` — 11 tests.
- `internal/audio/tts/helper/testdata/wp01.helper.json`.

No additional source was needed for unrelated native TTS engines beyond the catalog/pipeline points that determine reachability.

## 4. Contract-to-code compliance matrix

Legend: **Compliant** = current source/test evidence supports the contract; **Partial** = some behavior exists but a required path/projection is absent or incorrect; **Noncompliant** = source contradicts the contract; **Missing** = no implementation; **Unestablished** = source cannot establish runtime behavior.

### A. Capability policy

| ID | Contract requirement | Current evidence | Grade | Result |
|---|---|---|---|---|
| A1 | Required service capabilities are exactly `howl.tts.load` and `howl.tts.synthesize`. | `protocol.go:53-55` defines `mandatoryServiceCapabilities` with `catalog`, `load`, `unload`, `synthesize`, `cancel`; `protocol.go:545-564` validates that exact list. | Observed | **Noncompliant** |
| A2 | `catalog`, `unload`, and `cancel` are optional. | Handshake validation rejects absence of any of the five; no optional concept exists in the type. | Observed / Source-derived | **Noncompliant** |
| A3 | Absence means unsupported. | `Client` gates speed/language provider features; catalog/load/synthesize have operation gates, but unload/cancel gate before process start. No explicit unsupported-set model. | Source-derived / Partial | **Partial** |
| A4 | Missing optional capability fails closed at the operation boundary without an unsupported frame. | `Catalog` checks catalog gate; `Load`/`Synthesize` check their gates; `Unload`/`Cancel` return `unsupported_capability` before `Start`. No test exercises all five optional-absence paths. | Source-derived / Partial | **Partial** |
| A5 | No unsupported operation frame is sent. | Gates precede writes for catalog/load/synthesize; unload/cancel never reach transport when unadvertised. | Source-derived | **Partial** (structural, not exhaustively tested) |
| A6 | Missing required capability produces `protocol_incompatible`. | Handshake missing/unknown/wrong-version/duplicate returns `protocol_violation`; discovery/start failures become `unavailable`. | Observed | **Noncompliant** |
| A7 | Unknown, duplicate, malformed, and wrong-version descriptors have adopted classification. | `validateCapabilities` returns `protocol_violation` for all; no process-scoped/provider-scoped distinction. | Observed | **Noncompliant** |
| A8 | Provider capabilities and process-scoped service capabilities stay separate. | One `Capabilities`/`ServiceCapabilities` validator and one set of capability IDs; no scope field or scope validation. | Observed / Source-derived | **Noncompliant** |
| A9 | `detail.missing` lists only missing required capabilities, in stable order. | Missing list is a comma join of all five mandatory IDs (`protocol.go:527-531`), including optional ones. | Observed | **Noncompliant** |
| A10 | `hasCapability` is the validated-handshake invariant for every frame. | Gates exist for operations, but `validateHello` does not enforce that the advertised set is a subset of known service capabilities; `Unload`/`Cancel` gate before handshake; no uniform invariant. | Source-derived | **Partial / Noncompliant** |
| A11 | No silent native fallback/downgrade. | No helper fallback to native provider found; helper failure is surfaced. | Source-derived | **Compliant structurally** |
| A12 | Wire capability shape remains `{id,version}`. | Fake protocol and helper protocol both use `{id,version}`; golden fixture preserves it. | Observed | **Compliant** |

### B. Lifecycle and recovery

| ID | Contract requirement | Current evidence | Grade | Result |
|---|---|---|---|---|
| B1 | Thirteen-state vocabulary is representable. | State constants cover absent/starting/handshaking/unloaded/loading/ready/busy/degraded/stopping/crashed/incompatible/unavailable/backoff. | Observed | **Compliant** |
| B2 | Successful handshake enters `unloaded`. | `MarkHandshake` sets `providerState=unloaded`. | Observed | **Compliant** |
| B3 | Semantic handshake failure becomes terminal `incompatible`. | Non-`ok` handshake returns error and `client.go:238-254` classifies transport/start as unavailable; `protocol_violation` remains; child is not explicitly terminated by the classification path. | Observed / Source-derived | **Noncompliant** |
| B4 | `incompatible` is terminal and is not restarted by ordinary calls. | No `incompatible` state is set by semantic failure; the generic explicit-call restart path (`client.go:201-212`) can restart a crashed or unavailable child without an incompatible gate. | Source-derived | **Noncompliant** |
| B5 | Process alive is determined solely by actual process state. | `processAlive` (`supervisor.go:351-356`) suppresses `crashed` and `unavailable` when `proc != nil`, so a live process can be reported through a non-alive state, and vice versa. | Observed | **Noncompliant** |
| B6 | `Available` means a configured usable provider. | `Available` excludes `unavailable` but does not distinguish configured provider from helper process or loaded model; no asset-specific availability. | Observed | **Partial / Noncompliant** |
| B7 | `Loaded` is projected. | `TTSProviderStatus` has no `Loaded`; nothing records it. | Observed | **Missing** |
| B8 | `Ready` implies loaded. | `Ready` is set by load success and synthesize success; no invariant tying it to `Loaded`; timeout cleanup can set `Ready` after kill. | Source-derived | **Partial** |
| B9 | Ordinary `unloaded` synthesis uses ensure-handshake/load. | `Synthesize` (`client.go:679-760`) calls `ensureStarted` but never `Load`; real fake helper returns `provider_not_loaded` (`fakehelper/server.go:495-501`). | Observed | **Noncompliant** |
| B10 | Synthesis while `State=unavailable` does not recover; explicit same-identity `load` is required. | `Synthesize` has no state gate and can call `ensureStarted`; no identity capture; no explicit recovery contract. | Source-derived | **Noncompliant** |
| B11 | In-place recovery requires matching `provider_id`, `engine`, `model_identity`. | No identity fields are captured or compared anywhere in `Client`/`Supervisor`. | Missing | **Missing** |
| B12 | Different identity follows A2-D7 replacement. | No identity divergence detection or replacement path. | Missing | **Missing** |
| B13 | Restart budget and backoff are applied. | `TTSStateBackoff` is declared but never assigned; no budget, no timer, no `backoff` transition. | Observed | **Missing** |
| B14 | A2-D12 correction does not inherit A2-D7 budget exemption. | No implementation of either rule; cannot satisfy. | Missing | **Missing** |
| B15 | Idle close exists. | No idle timer or idle shutdown path in `Client`/`Supervisor`. | Missing | **Missing** |
| B16 | Shutdown classifies expected versus unexpected exit. | `observeExit` (`supervisor.go:166-195`) maps any nonzero exit from `providerState=stopping` to `Absent`; unexpected nonzero becomes absent, not crashed. | Observed | **Noncompliant** |
| B17 | Crash/backoff/unavailable transitions are consistent. | `observeExit` sets crashed/unavailable only in selected branches; no backoff transition; `SetState` may overwrite process state. | Observed / Source-derived | **Noncompliant** |
| B18 | `ServiceCapabilitiesKnown` is projected. | No such field or flag. | Missing | **Missing** |

### C. Error handling and protocol correctness

| ID | Requirement | Current evidence | Grade | Result |
|---|---|---|---|---|
| C1 | `protocol_incompatible` for semantic handshake invalidity. | Only wire-code version errors use it; semantic hello errors use `protocol_violation`. | Observed | **Noncompliant** |
| C2 | `protocol_violation` for sequencing/framing violations, not capability semantics. | `validateHello` and `validateCapabilities` use it for missing/unknown/duplicate/wrong-version capability descriptors. | Observed | **Noncompliant** |
| C3 | `invalid_frame` for malformed/oversized/duplicate-key frames. | Strict decode and dedicated tests exist; error mapping is present. | Observed | **Compliant** |
| C4 | `unsupported_capability` at operation boundary. | Used by gates; exact message/code constructed at client boundary. | Source-derived | **Partial** (some paths predate handshake) |
| C5 | `provider_not_loaded` is surfaced. | Fake helper returns it when unloaded; parent maps helper error but may leave or set incorrect state. | Observed / Source-derived | **Partial** |
| C6 | `model_missing` is surfaced. | Wire code exists; no current test in the inspected set directly establishes parent preservation across a real process. | Unestablished | **Unestablished** |
| C7 | `helper_start_failed` is surfaced and stateful. | Start failure sets `unavailable`; error is returned. | Source-derived | **Compliant structurally** |
| C8 | Timeout and cancellation outcomes are correct. | `TestTimeoutForcesTermination` and cancellation tests pass historically; `client.go:766-840` timeout/cancel cleanup can set `Ready` after force-termination. | Observed / Source-derived | **Partial / Noncompliant** |
| C9 | `detail.missing` is bounded, ordered, and survives recovery failures. | No detail field on `TTSCapability`; missing list is five IDs; no recovery failure detail. | Missing / Observed | **Noncompliant** |
| C10 | Parent acceptance, duplicate terminal, retired IDs, cancellation precedence, publication ordering hold. | Dedicated tests cover acceptance-before-frame, duplicate terminal, retired ID, cancellation, terminal publication. | Observed | **Compliant** |
| C11 | `LastError` preserves diagnostic identity/detail. | Generic status error strings exist; no structured detail/identity projection. | Source-derived | **Partial** |

### D. Process and file ownership

| ID | Requirement | Current evidence | Grade | Result |
|---|---|---|---|---|
| D1 | Parent owns helper lifecycle. | `Supervisor` creates and terminates; `Client` owns supervisor. | Source-derived | **Compliant structurally** |
| D2 | Windows job object and suspended startup. | `job_windows.go:84-149` creates suspended process, assigns job, resumes. | Observed structural | **Compliant structurally**; OS boundary not re-exercised |
| D3 | Single waiter/exit observation. | `osProcess.Wait` called by `Supervisor.observeExit`; tests assert one waiter. | Observed / Source-derived | **Compliant structurally** |
| D4 | Kill/wait/cleanup ordering. | `Terminate` sets stopping, terminates, waits via observer; cleanup is bounded. Tests cover termination. | Source-derived | **Compliant structurally** |
| D5 | Output path containment and reparse protection. | Parent validates request path and helper output path; path checks per platform. | Observed | **Compliant structurally** |
| D6 | WAV validation and output ownership. | `output.go` checks exact path, WAV format, expected hash; cleanup targets exact paths. | Observed | **Compliant** |
| D7 | No-overwrite publication. | Fake helper uses hard-link no-replace publication; parent does not implement publication. Contract B12a divergence remains deferred. | Observed | **Partial** |
| D8 | Failure cleanup and stale output handling. | `cleanupRejected`, exact-path removal, force cleanup; tests cover stale output. | Source-derived | **Partial** (state side effects) |
| D9 | Shutdown and idle close. | Shutdown path exists; no idle path. | Observed | **Partial / Missing** |
| D10 | Unexpected nonzero exit is distinguishable. | Nonzero while stopping becomes absent. | Observed | **Noncompliant** |
| D11 | Final symlink/junction protections. | Path checks exist; Windows final-symlink test was historically skipped. | Unestablished | **Unestablished at OS boundary** |
| D12 | Real native helper publishes output. | No production native helper; only fake helper. | Missing | **Unestablished** |

### E. Production integration and reachability

| ID | Requirement | Current evidence | Grade | Result |
|---|---|---|---|---|
| E1 | Helpers can be reached from production. | Only import is `cmd/howl-tts-helper-fake/main.go:6`; no production import of helper/fakehelper. | Observed | **Missing** |
| E2 | `TTSSynthesizer` unchanged. | `internal/audio/types/interfaces.go:127` unchanged. | Observed | **Compliant as a constraint** |
| E3 | `TTSProviderInspector` implemented. | Declared at `types/provider.go:109-114`; `Client` has matching `ProviderStatus`/`ProviderCatalog` methods; no compile-time assertion or production consumer. | Observed | **Partial** |
| E4 | Production uses supervisor/capability policy. | No production path constructs `helper.Client` or `Supervisor`. | Observed | **Missing** |
| E5 | Explicit production `load` exists. | `Load` exists only on `helper.Client` (`client.go:551-595`); no caller. | Observed | **Missing / blocking prerequisite** |
| E6 | Recovery from `unavailable` is reachable without new API. | No production caller and no state-gated recovery path; adoption I1 acknowledges this. | Observed | **Missing** |
| E7 | Ordinary unloaded vs unavailable are distinguishable. | `Synthesize` does not branch on state; `ensureStarted` can restart any non-ready state. | Source-derived | **Noncompliant** |
| E8 | Adapter loading rules for `Initialize`/`Synthesize`/catalog/status. | `Initialize` would need to load under A2-D12, `Synthesize` needs unloaded ensure-load, catalog must not load, status must not mutate. No adapter exists. | Source-derived | **Missing** |
| E9 | Helper implements production synthesizer. | No `Initialize`, `GetVoices`, `SynthesizeStream`, `SupportsSSML`, `Release`; `Client.Synthesize` signature differs. | Observed | **Missing** |
| E10 | `GetTTSVoiceCatalog` stays static/read-only. | `app-tts-catalog.go:49-96` builds static profiles/native catalogs; no helper access. | Observed | **Compliant for current path** |
| E11 | `ProviderStatus`/`ProviderCatalog` may start but must not load. | `ProviderStatus` calls `Start`; `ProviderCatalog` calls `Catalog` which calls `Start`; `Catalog` does not load. No production use. | Observed / Source-derived | **Partial** |

## 5. Thirteen-state lifecycle comparison

The table compares the adopted lifecycle vocabulary and required transition semantics with the current implementation. A state being “representable” does not mean its adopted entry/exit conditions are implemented.

| State | Contract meaning / required transitions | Current implementation | Grade | Result |
|---|---|---|---|---|
| `absent` | No helper process; initial or updated/disabled provider. Explicit load may start if eligible. | Initial state; `setUnavailable` resets to absent; no provider-update transition gate. | Observed | **Partial** |
| `starting` | Process creation/job assignment. Failure to start → `unavailable`; success → `handshaking`. | `Start` sets starting then handshaking on successful factory handshake. | Observed | **Compliant** |
| `handshaking` | Wire exchange; valid hello → `unloaded`; semantic invalidity → `incompatible`; crash → `crashed`. | Valid hello → unloaded. Semantic invalidity → error/`protocol_violation`, not incompatible. Crash observer exists. | Observed | **Partial / Noncompliant** |
| `unloaded` | Handshake complete, no loaded model. Ordinary synthesize performs ensure-load; explicit load allowed. | Handshake sets it. `Synthesize` does not ensure-load. | Observed | **Noncompliant** |
| `loading` | Load in flight; success → `ready`; failure → `unavailable`; A2-D12 may settle `unloaded`; crash → crashed. | `Load` sets loading/ready, failure `unavailable`; no A2-D12 reverse transition. | Observed | **Partial** |
| `ready` | Loaded and usable; synthesize → busy; unload → unloaded. | Set by load and synthesize success; no loaded invariant; unload path exists. | Source-derived | **Partial** |
| `busy` | Active synthesize/cancel; success → ready; failure → degraded; shutdown → stopping. | Set by synthesize; success ready; error degraded; observer can overwrite. | Observed | **Partial** |
| `degraded` | Recoverable failed synthesis; can retry or stop. | Set on synthesis failure; no recovery/detail discipline. | Observed | **Partial** |
| `stopping` | Intentional termination; expected exit → absent; unexpected nonzero → crashed. | Set by `Terminate`; any nonzero exit while stopping → absent. | Observed | **Noncompliant** |
| `crashed` | Unexpected exit; restart budget/backoff then starting; exhausted → unavailable. | Observer can set crashed, but no budget/backoff and some exits are misclassified. | Source-derived | **Partial / Missing** |
| `incompatible` | Semantic handshake incompatibility; terminal; no ordinary restart. | Not reached; no terminal flag; generic restart path lacks incompatible gate. | Observed | **Missing / Noncompliant** |
| `unavailable` | Helper unavailable/failed; synthesize must not recover; explicit matching-identity load may recover. | Set on start/load errors; synthesize has no state gate; no identity. | Observed | **Noncompliant** |
| `backoff` | Restart delay with budget accounting. | Constant declared but never assigned; no timer/budget. | Observed | **Missing** |

### 5.1 Transition audit

| Transition | Required | Current |
|---|---|---|
| `absent → starting` | On explicit load. | `Start` can start from absent. |
| `absent → unavailable` | Discovery/start failure. | `Start` sets unavailable. |
| `starting → handshaking` | Successful process creation. | Implemented. |
| `starting → unavailable` | Start failure. | Implemented. |
| `handshaking → unloaded` | Valid hello. | Implemented. |
| `handshaking → incompatible` | Semantic hello failure. | Missing; `protocol_violation` and process remains usable. |
| `handshaking → crashed` | Unexpected exit. | Observer can set crashed; `SetState` may force alive. |
| `unloaded → loading` | Ordinary synthesize or explicit load. | Explicit `Load` only; `Synthesize` does not call it. |
| `loading → ready` | Load success. | Implemented. |
| `loading → unavailable` | Load failure. | Implemented; detail weak. |
| `loading → unloaded` | A2-D12 allowed reverse transition. | Missing. |
| `loading → crashed` | Unexpected exit. | Observer path exists; no backoff. |
| `unavailable → loading` | Explicit same-identity load. | No state/identity gate; no production caller. |
| `ready → busy` | Synthesis begins. | Implemented. |
| `busy → ready` | Synthesis succeeds. | Implemented. |
| `busy → degraded` | Synthesis fails. | Implemented, but exit observer may also write a terminal state. |
| `ready/degraded → busy` | Retry synthesis. | Implemented. |
| `ready/degraded/unloaded → stopping` | Shutdown/idle. | Shutdown sets stopping; no idle. |
| `busy → stopping` | Shutdown during work. | `Shutdown` calls transport before terminate; no explicit busy gate. |
| `stopping → absent` | Expected exit. | Any nonzero becomes absent. |
| `stopping → crashed` | Unexpected nonzero exit. | Missing. |
| `crashed → backoff → starting` | Budgeted restart. | Missing. |
| `crashed/backoff → unavailable` | Budget exhausted. | Missing. |
| `incompatible/unavailable → absent` | Provider update/disable. | No provider-update transition. |
| `any live → crashed` | Unexpected exit. | Observer can set crashed, but `SetState` and `processAlive` can contradict it. |

## 6. Capability and error classification comparison

### 6.1 Adopted classification

Under §17.19, a semantic handshake that cannot satisfy the service contract — including missing required capability, unknown service capability, wrong version, duplicate descriptor, malformed capability set, invalid provider capability, or process-scoped service capability appearing in a provider scope — is a terminal incompatibility: `protocol_incompatible`. It is not a recoverable `protocol_violation`.

`protocol_violation` remains appropriate for protocol sequencing and frame legality, such as duplicate terminal frames, retired request IDs, frames before acceptance, or terminal publication ordering violations. `invalid_frame` remains for decode-level malformation, duplicate keys, and oversize frames.

### 6.2 Observed implementation

| Case | Current classification/path | Grade | Result |
|---|---|---|---|
| Missing required service capability | `validateHello` → `protocol_violation`; client maps to an error, not `incompatible`. | Observed | Wrong |
| Missing optional service capability | Rejected during hello as part of five-mandatory set. | Observed | Wrong |
| Unknown service capability | `protocol_violation`. | Observed | Wrong |
| Wrong capability version | `protocol_violation`; wire-code mismatch uses `protocol_incompatible`. | Observed | Wrong |
| Duplicate capability descriptor | `protocol_violation`. | Observed | Wrong |
| Duplicate provider/process-scoped capability | No scope distinction; duplicate check only within the unified list. | Observed | Wrong |
| Invalid capability descriptor | `protocol_violation`. | Observed | Wrong |
| Malformed frame / duplicate JSON key | `invalid_frame`. | Observed | Correct |
| Oversize frame | `frame_too_large`/invalid frame path. | Observed | Correct |
| Duplicate terminal / retired request / pre-acceptance frame | `protocol_violation`. | Observed | Correct |
| Unsupported operation | `unsupported_capability` before frame. | Source-derived | Partial: unload/cancel gate before handshake; no exhaustive test |
| Provider not loaded | Helper sends `provider_not_loaded`; parent maps it. | Observed | Correct code, incorrect lifecycle state |
| Model missing | Wire error exists; real-process parent preservation not established in this pass. | Unestablished | Unestablished |
| Helper start failed | `unavailable` + `helper_start_failed` error. | Source-derived | Structurally correct |
| Timeout | `context.DeadlineExceeded`; cleanup force path. | Observed/tested | Correct error, wrong state after kill |
| Cancellation | Cancellation watcher, cancel frame, terminal precedence. | Observed/tested | Correct protocol, possible state corruption |

`detail.missing` is built at `protocol.go:527-531` by joining all five mandatory IDs. It is not limited to the two required IDs, is not tied to a typed `Detail` on a capability, and does not preserve provider-scoped detail. Under the adopted contract it is both semantically wrong and diagnostically incomplete.

## 7. Process and file ownership comparison

### 7.1 Structurally established

- **Parent-owned process:** `helper.Client` owns the supervisor; `Supervisor` creates the process. No production caller exists, but ownership code is coherent (`client.go`, `supervisor.go`).
- **Windows job object:** `internal/proc/job_windows.go:84-149` starts suspended, creates/attaches a job, and resumes only after assignment. Job setup failures fail startup rather than silently running uncontained.
- **Single waiter:** `Supervisor.observeExit` is the sole blocking wait site; tests assert the helper does not create a second waiter.
- **Kill/wait:** `Terminate` sets `stopping`, kills, and waits via the observed exit path. Timeout/cancel force paths call `cleanupRejected` and terminate.
- **Output containment:** parent validates requested output path; helper validates its output path; `output.go` validates exact path, WAV format, and expected hash; cleanup removes only exact owned paths.
- **Stale output:** tests cover pre-existing/stale output and force cleanup.
- **Fake publication:** `fakehelper/publish.go` uses a no-replace hard link. This is deliberately conservative but is not the B12a atomic-rename decision, which remains deferred.

### 7.2 Not established or not compliant

- **OS-boundary execution:** no process was launched in this investigation. Windows job ownership, handle inheritance, final symlink behavior, and forced-kill observation are documented from source and historical tests, not re-exercised here.
- **Final symlink/junction:** the Windows final-symlink test was historically skipped; source path checks are not proof of OS behavior.
- **Nonzero exit:** stopping + nonzero → `absent`, losing the unexpected-exit distinction required by the lifecycle.
- **State overwrite after cleanup:** timeout/cancel cleanup can set `Ready` after the transport is killed (`client.go:766-840`), so status can be internally inconsistent.
- **Real native helper publication:** only the fake helper exists; no real helper implementation establishes B12a/B12b behavior.
- **Idle close:** absent.

## 8. Production integration and caller reachability

### 8.1 Current reachability

- `rg`/code search finds only one import of the helper family: `cmd/howl-tts-helper-fake/main.go:6` imports `howl-chat/internal/audio/tts/fakehelper`. No production package imports `internal/audio/tts/helper` or `fakehelper`.
- `internal/audio/pipeline.go:101-185` initializes native TTS engines; `:418-455` loads native TTS for a request. It does not construct the helper.
- `internal/audio/manager_factory.go` does not construct the helper.
- `app-tts-catalog.go:49-96` uses static provider profiles and native catalog switches, not the helper.
- `TTSProviderInspector` is declared at `internal/audio/types/provider.go:109-114`. `helper.Client` has matching `ProviderStatus` and `ProviderCatalog` methods (`client.go:952-963`), but there is no compile-time assertion and no production consumer.
- `TTSSynthesizer` remains unchanged at `internal/audio/types/interfaces.go:127`. The helper has no `Initialize`, `GetVoices`, `SynthesizeStream`, `SupportsSSML`, or `Release` adapter. Its `Synthesize` signature does not match the production interface.

### 8.2 Loading-rule reachability

| Operation | Required behavior | Current reachable behavior |
|---|---|---|
| `GetTTSVoiceCatalog` | Do not start or load helper. | Static/native only; no helper. |
| `ProviderCatalog` | May start/handshake; must not load. | `Client.ProviderCatalog` → `Catalog` → `Start`; no load. No production caller. |
| `ProviderStatus` | May start/handshake; must not load. | `Client.ProviderStatus` → `Start`; no load. No production caller. |
| `Initialize` | Load under A2-D12; no helper adapter. | Missing. |
| `Synthesize` | From `unloaded`, ensure-load; from `unavailable`, do not recover. | `Synthesize` only `ensureStarted`, never `Load`; no state gate. |
| Explicit recovery `load` | Same-identity only; different identity uses A2-D7. | `Client.Load` exists; no caller, no identity, no budget. |

**Conclusion:** the capability policy is not enforced by any production path because no production path reaches the helper. Recovery from `unavailable` is not reachable without adding a caller/API path. Adoption record item I1 already acknowledges this; it is a prerequisite for Stage 3 integration, not a test-only issue.

## 9. Test and evidence coverage matrix

### 9.1 Strengths

- `client_test.go` has 16 tests covering lifecycle, catalog, provider feature gates, capability gates, timeouts, cancellation, and output/path behavior.
- `protocol_test.go` has 10 tests covering frame decode, oversize, duplicate keys, terminal ordering, and capability validation.
- `supervisor_test.go` has 5 tests covering state projection, waiter ownership, termination, and exit observation.
- `interop_test.go` has 12 stage tests, 11 of which launch a real fake-helper process; `TestStage2MissingHelperIsNonfatal` does not.
- `fakehelper/server_test.go` has 11 tests covering modes, no-overwrite publication, output validation, and helper protocol behavior.
- Acceptance barriers, duplicate terminal handling, retired request IDs, cancellation precedence, and terminal publication ordering are directly tested.
- Historical Stage 2 audit/correction records report real-process and race runs passing, with a pre-existing flaky `internal/audio` cache test in the full suite.

### 9.2 Coverage classification

| Requirement | Classification | Evidence |
|---|---|---|
| Required-two/optional-three policy | **Contradicted** | `protocol_test.go:81` expects missing `unload` to fail; `protocol_test.go:130` expects missing service caps to fail; `interop_test.go:164` and `fakehelper/server_test.go:180` assert five. |
| Optional absence fails closed | **Not tested** | No dedicated optional-absence test; `TestClientRejectsMissingServiceCapability` removes a required cap after handshake. |
| Missing required → `incompatible` | **Not tested / contradicted** | Tests expect `protocol_violation` for missing capability sets. |
| Semantic hello incompatibility | **Not tested** | No semantic-incompatible real-process test. |
| Process/provider scope separation | **Not tested** | No scope field or test. |
| Typed capability detail/`detail.missing` | **Not tested** | No status/detail assertions; type lacks fields. |
| `Loaded`/service capability status | **Not tested** | Status type lacks fields. |
| Ordinary unloaded synthesis | **Not tested in real process** | `TestStage2…` paths load explicitly; in-memory tests may pre-load. |
| `unavailable` non-recovery | **Not tested** | No state-gated synthesis test. |
| Explicit same-identity recovery | **Not tested** | No identity implementation. |
| A2-D7 different-identity replacement | **Not tested** | No implementation. |
| Restart budget/backoff | **Not tested** | No implementation. |
| Idle close | **Not tested** | No implementation. |
| Unexpected nonzero shutdown exit | **Not tested/contradicted** | `observeExit` maps to absent. |
| Timeout/cancel state consistency | **Indirectly tested** | Existing tests exercise errors, not post-kill status invariants. |
| Windows job/suspended startup | **Structurally established + historical real-process** | Source and Stage 2 records; not re-run here. |
| Final symlink/junction | **Unestablished** | Windows test historically skipped. |
| Real native helper publication | **Not tested** | No native helper. |
| Production integration | **Not tested / absent** | No production caller. |
| Race/full-suite stability | **Historical only** | Prior audit/correction; pre-existing flaky `internal/audio` cache test. |

### 9.3 What the test suite proves and does not prove

The suite proves substantial protocol and process mechanics for the fake helper under explicit load. It does **not** prove adopted capability policy, status projection, ordinary unloaded behavior, unavailable recovery, identity/replacement rules, restart budgeting, idle close, production adapter compatibility, or OS-boundary behavior after the amendment. Several tests encode the superseded five-service policy and will actively fail a correct implementation.

## 10. Findings

Severity definitions: **blocking** = must be resolved before Stage 3 authorization; **major** = contract-required and materially unsafe/incorrect, normally must precede Stage 3; **moderate** = important but can be explicitly carried into a Stage 3 slice; **minor** = localized; **informational** = no action required for authorization.

### S3R-01 — Five mandatory service capabilities contradict A2-D1

- **Contract:** §17.19 A2-D1–A2-D6; required exactly `load`, `synthesize`; optional `catalog`, `unload`, `cancel`.
- **Evidence:** `protocol.go:53-55`, `protocol.go:527-531`, `protocol.go:545-564`; `protocol_test.go:81`, `:130`; `interop_test.go:164`; `fakehelper/server_test.go:180`.
- **Grade:** Observed.
- **Severity:** **Blocking.**
- **Action:** Change the handshake floor to required-two/optional-three; reject only absence of required capabilities; keep optional absence as unsupported; update tests/fixtures.
- **Before Stage 3:** Yes.

### S3R-02 — Semantic handshake failure is not terminal `protocol_incompatible`

- **Contract:** §17.19 A2-D2/A2-D6; semantic invalidity is `protocol_incompatible`, terminal, no ordinary restart.
- **Evidence:** `protocol.go:502-543`; `client.go:238-254`; `client.go:201-212`; `protocol_test.go:130`.
- **Grade:** Observed / Source-derived.
- **Severity:** **Blocking.**
- **Action:** Classify semantic hello/capability failures as `protocol_incompatible`; terminate the child; prevent generic restart; preserve `detail`.
- **Before Stage 3:** Yes.

### S3R-03 — Service and provider capability scopes are conflated

- **Contract:** §17.19 process-scoped service capabilities versus provider capabilities; scope validation.
- **Evidence:** `protocol.go:545-564`; `supervisor.go:376-384`; `types/provider.go:49-52`.
- **Grade:** Observed.
- **Severity:** **Blocking.**
- **Action:** Add explicit scope validation and separate storage/projection; reject cross-scope descriptors as incompatible.
- **Before Stage 3:** Yes.

### S3R-04 — `detail.missing` is wrong and capability detail is unrepresentable

- **Contract:** §17.19 typed capability detail; `detail.missing` only missing required IDs, stable order/limits.
- **Evidence:** `protocol.go:527-531`; `types/provider.go:49-52`.
- **Grade:** Observed / Missing.
- **Severity:** **Major.**
- **Action:** Add `Supported`/`Detail`; compute missing required only; preserve diagnostic detail.
- **Before Stage 3:** Yes.

### S3R-05 — Status model cannot represent adopted §17.12

- **Contract:** `Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown`, typed capability `{ID,Version,Supported,Detail}`.
- **Evidence:** `types/provider.go:49-71`; `supervisor.go:306-331`; `supervisor.go:376-384`.
- **Grade:** Observed.
- **Severity:** **Blocking.**
- **Action:** Extend status types and projection; keep provider and service capabilities separate; add tests.
- **Before Stage 3:** Yes.

### S3R-06 — ProcessAlive/Available/Ready projections contradict lifecycle rules

- **Contract:** process state solely process ownership; `Available` configured provider; `Ready` loaded.
- **Evidence:** `supervisor.go:238-243`, `:351-356`, `:358-365`.
- **Grade:** Observed.
- **Severity:** **Blocking.**
- **Action:** Remove suppression of crashed/unavailable by live proc; derive process state solely from process; add loaded invariant.
- **Before Stage 3:** Yes.

### S3R-07 — Ordinary unloaded synthesis does not ensure-load

- **Contract:** §17.19 A2-D12 ordinary `unloaded` synthesize retains ensure-handshake/load.
- **Evidence:** `client.go:679-760`; `fakehelper/server.go:495-501`.
- **Grade:** Observed.
- **Severity:** **Blocking.**
- **Action:** Add state-aware ensure-load for `unloaded`; do not do so for `unavailable`.
- **Before Stage 3:** Yes.

### S3R-08 — `unavailable` recovery is ungated and unreachable

- **Contract:** synthesize while unavailable does not recover; explicit matching-identity load required.
- **Evidence:** `client.go:201-212`, `:551-595`, `:679-760`; no production import.
- **Grade:** Observed / Source-derived.
- **Severity:** **Blocking.**
- **Action:** Add state gate, explicit recovery path, and production reachability decision.
- **Before Stage 3:** Yes.

### S3R-09 — No provider/engine/model identity tracking

- **Contract:** A2-D7/A2-D12 identity matching and replacement rules.
- **Evidence:** no identity fields in `Client`/`Supervisor`; `Load` sends to current process.
- **Grade:** Missing.
- **Severity:** **Blocking.**
- **Action:** Capture identity from handshake/load; compare; implement replacement and budget rules.
- **Before Stage 3:** Yes.

### S3R-10 — No restart budget or backoff

- **Contract:** bounded restarts, delay, `backoff` state.
- **Evidence:** `TTSStateBackoff` declared but never assigned; no timer/budget.
- **Grade:** Observed.
- **Severity:** **Blocking.**
- **Action:** Implement budget/backoff transitions and exhaustion to unavailable.
- **Before Stage 3:** Yes.

### S3R-11 — No idle close

- **Contract:** idle close required.
- **Evidence:** no timer/path.
- **Grade:** Missing.
- **Severity:** **Blocking.**
- **Action:** Implement idle transition to stopping and process cleanup.
- **Before Stage 3:** Yes.

### S3R-12 — Shutdown misclassifies unexpected nonzero exit

- **Contract:** expected exit → absent; unexpected nonzero → crashed.
- **Evidence:** `supervisor.go:166-195`.
- **Grade:** Observed.
- **Severity:** **Blocking.**
- **Action:** Track expected stop; classify exit code accordingly; test both.
- **Before Stage 3:** Yes.

### S3R-13 — Timeout/cancel cleanup can report ready/live after kill

- **Contract:** health/status must reflect actual process/lifecycle.
- **Evidence:** `client.go:766-840` sets `Ready` after `cleanupRejected`/force termination.
- **Grade:** Source-derived.
- **Severity:** **Major.**
- **Action:** Set terminal/unavailable state after forced cleanup; add post-kill status assertions.
- **Before Stage 3:** Yes.

### S3R-14 — Unload/cancel gate before handshake

- **Contract:** optional absence fails closed at operation boundary; `hasCapability` validated-handshake invariant.
- **Evidence:** `client.go:599-605`, `:622-629`.
- **Grade:** Source-derived.
- **Severity:** **Major.**
- **Action:** Ensure handshake/validated capability set before evaluating optional operations; no frame when unsupported.
- **Before Stage 3:** Yes.

### S3R-15 — No production adapter or caller

- **Contract:** Stage 3 integration/reachability; ordinary ensure-load and explicit recovery must be usable.
- **Evidence:** only `cmd/howl-tts-helper-fake/main.go:6` imports fakehelper; no helper import; `TTSSynthesizer` unchanged; missing adapter methods.
- **Grade:** Observed.
- **Severity:** **Major** (blocking for Stage 3 integration, not for the Stage 2 milestone itself).
- **Action:** Define and implement the adapter/load/recovery path without changing `TTSSynthesizer`; add production caller tests.
- **Before Stage 3:** Yes, before integration.

### S3R-16 — Tests encode superseded policy and omit adopted paths

- **Contract:** adopted capability/lifecycle/status/recovery policy.
- **Evidence:** §9.2; `protocol_test.go:81`, `:130`; five-capability assertions.
- **Grade:** Observed.
- **Severity:** **Major.**
- **Action:** Replace contradictory expectations; add optional-absence, incompatible, status, unloaded, recovery, identity, budget, idle, shutdown-exit, production-adapter tests.
- **Before Stage 3:** Yes.

### S3R-17 — OS-boundary and native-helper evidence is incomplete

- **Contract:** process/file ownership and publication behavior.
- **Evidence:** Windows final-symlink test skipped; no native helper; no run in this pass.
- **Grade:** Unestablished.
- **Severity:** **Moderate.**
- **Action:** Execute separately under explicit authorization on Windows, including symlink/junction and forced-kill cases; implement native publication when authorized.
- **Before Stage 3:** No for parent Stage 3; yes before claiming native end-to-end readiness.

### S3R-18 — Historical Stage 2 audit is superseded for current report state

- **Contract:** evidence provenance.
- **Evidence:** adversarial audit subject `adb96b06…`; current corrected report `ec85b2b2…`.
- **Grade:** Superseded / Calculated.
- **Severity:** **Informational.**
- **Action:** Cite the correction audit for current artifacts; retain the adversarial audit as historical.
- **Before Stage 3:** No.

## 11. Required prerequisites before Stage 3 authorization

The following must be completed and independently verified before Stage 3 implementation is authorized. This report does not authorize any of them.

1. **Correct the capability floor and classification.** Implement A2-D1–A2-D6; required-two/optional-three; semantic invalid → terminal `protocol_incompatible`; strict scope separation; correct `detail.missing`; preserve wire shape.
2. **Correct the status model.** Add typed capability `Supported`/`Detail`, `Loaded`, `ServiceCapabilities`, `ServiceCapabilitiesKnown`; keep provider/process scopes separate; project process/availability/readiness correctly.
3. **Implement lifecycle/recovery.** Ordinary `unloaded` ensure-load; `unavailable` non-recovery; explicit same-identity load; A2-D7 replacement; A2-D12 reverse-settle; restart budget/backoff; idle close; expected/unexpected shutdown; incompatible terminality; consistent observers.
4. **Resolve integration reachability.** Provide the adapter/load/recovery path and a production caller without changing `TTSSynthesizer`; make ordinary and unavailable behavior distinguishable; do not mistake `ProviderStatus`/`Catalog` for loading APIs.
5. **Replace contradictory tests and add adopted-path coverage.** Optional absence, semantic incompatibility, scopes, detail, status, unloaded synthesis, unavailable/recovery, identity, replacement, budget, idle, shutdown exit, post-kill status, production adapter.
6. **Re-run authorized build/race/full-suite and Windows process checks.** This report does not run them. The pre-existing `internal/audio` cache-test flake must be separated from helper stability claims.
7. **Decide the deferred publication/protocol matters** (B12a/B12b and O-1–O-7) before native provider work; they need not all block the parent Stage 3 slice, but the Stage 3 plan must state their status.

Until these are addressed, a Stage 3 implementation would build on a baseline that contradicts its governing capability and lifecycle contract.

## 12. Explicitly deferred or unresolved matters

- **B12a/B12b:** atomic publication semantics and extension wording remain deferred. The fake helper’s hard-link no-replace behavior is not a resolution.
- **O-1–O-7:** unresolved owner/architecture items from the amendment remain unresolved; this investigation does not decide them.
- **I1:** adoption record’s integration/reachability gap is confirmed. No production caller and no explicit load path exist.
- **Restart numbers:** no budget/backoff values are implemented or validated; no provisional numbers are treated as accepted.
- **Windows final symlink/junction:** historically skipped; unestablished at the OS boundary.
- **Real native helper:** not present; native publication/exit behavior unestablished.
- **Race/full-suite stability:** historical only. The pre-existing `internal/audio` cache flake is not attributed to the helper.
- **Runtime behavior:** not re-executed in this investigation. All runtime claims are either source-derived or expressly tied to prior evidence.
- **Historical audits:** the adversarial Stage 2 audit is superseded for the current report hash; its findings remain useful historical evidence but are not current compliance proof.

## 13. Confirmation of no implementation changes

No implementation change was made.

- No production source file was modified.
- No test, fixture, golden file, protocol, contract, amendment, adoption record, prior report, or audit was modified.
- No helper was wired into the TTS manager or pipeline.
- No recovery API, UI, native provider, G2P, Kokoro, ONNX, CGO, or dependency was added.
- No build, test, process launch, replay, benchmark, or asset download was performed.
- No commit, push, or stage was performed.
- The only writes associated with this investigation are this report and one entry in `docs/review/README.md`.
- Stage 3 implementation remains unauthorized.
