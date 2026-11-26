# HRA-P1 — Implementation Plan

## Plan paper `HRA-P1-PLAN-2026-11-25-01`

| Field | Value |
|---|---|
| Date | **2026-11-25** |
| Type | Implementation planning record |
| Contract | `HRA-P1-RESOURCE-LIFECYCLE-CONTRACT-2026-11-25.md` |
| Characterization gate | **Closed** — D-1/D-9, D-3/D-6, D-4, and D-5 are graduated; 8 planned invariants remain |
| Status | **GATE 8 OPEN — voice-quality acceptance closed (human listening, 2026-11-25); ASR admission resolved, ASR/TTS demonstrated in both hosts, and the re-apply self-block fixed and re-verified in the packaged host (§7.6); image generation demonstrated in both hosts with real weights (§7.6.8); Gate 7 §15.1 attempted in the browser host — chat resident and clean shutdown pass, image leg blocked by a harness page-load failure, video impossible as configured (§7.6.9, §7.6.11); Kokoro/MeloTTS optional by owner guidance** |
| Production changes | **Gate 1: D-1 and D-9 in `internal/resource/admission.go`** |
| User-visible workflow demonstrated | **Audio yes, image yes, video no.** Voice quality accepted (operator heard local Piper); ASR transcription and Piper synthesis exercised in both supported hosts (§7.6.3, §7.6.5, §7.6.6); audio settings re-apply fixed and re-verified in the packaged host (§7.6.7); SDXL image generation produced real files in both hosts from `H:\AI GGUF\Image models` (§7.6.8). Video render is **not yet attempted** and the only video weight is incompatible as configured (§7.6.11); Kokoro/MeloTTS are uninstalled; the Gate 7 §15.1 browser-host run is **partial** — chat and shutdown pass, image leg blocked by a harness page hang (§7.6.9) — and has not been run in the packaged host |

This paper reconciles the closed characterization gate with HRA-P1. It does not add tests, production hooks, or instrumentation. Its purpose is to name the ownership changes production must make and the evidence each change must produce before the corresponding known failure or invariant is considered closed.

The characterization gate is frozen. Every test listed in contract §11 kept its status until the production change owning it was implemented. D-1 and D-9 graduated in Gate 1; D-3 and D-6 in Gate 2; D-4 and D-5 under their own behavioral records. Their original assertions are unchanged and run in the default suite. Source guards were retired in favor of behavioral tests and no tagged known-failure test remains.

---

## 1. Planning decision

**PLAN-01.** Phase-one implementation is organized around four ownership moves, not around defect-by-defect patches:

1. Registration ownership moves from callers to the resource authority.
2. Child-process ownership moves from engine fields to one process registry.
3. Execution ownership moves from implicit busy flags and callbacks to named execution scopes.
4. Capacity observation moves from registration-time release to pending/reflected/released transitions.

**PLAN-02.** A defect is closed only when its owning transition exists and its acceptance evidence passes in the default suite. Removing a known-failure build tag without the owning transition is prohibited.

**PLAN-03.** No test-only production seam is introduced. Where a defect is currently unobservable, the implementation must expose the real ownership boundary (authority request, process registry observation, execution scope) and test that boundary. A source-text guard may remain as a temporary tripwire, but never as the sole closure evidence.

**PLAN-04.** Existing image, chat, audio, and video behavior is preserved unless the work package explicitly replaces it. In particular, llama's single-child behavior, image generation quality, chat streaming, and workflow stage order do not change.

**PLAN-05.** Phase-one implementation does not include persistent calibration, learned placement, reserve optimization, UI redesign, or model-behavior changes.

---

## 2. Target ownership

| State | Current owner | Target owner | Production change |
|---|---|---|---|
| Admission and commitments | `resource.Authority` | Same, extended | Authority stores a commitment per lease with kind, owner, capacity state, and projected peak; release is reservation-only |
| Resident registrations | Callers calling `ResourceManager.RegisterModel` | Resource authority | `ResourceManager` becomes a read model; only the authority writes registrations |
| Image registration | `App.registerImageModel`, `Engine.Shutdown` | Resource authority | Engine reports process/readiness/exit; App requests transitions; no direct unregister |
| Chat registration | Load monitor | Resource authority | Monitor submits readiness plus process identity; authority verifies the current token |
| Audio registration | `Pipeline.registerAudio` | Resource authority | Pipeline requests registration with its load lease; no direct manager write |
| Image process | `sdcpp.Engine.serverProc` | Process registry | Registry owns start token, single `Wait`, exit observation, and stop |
| Video process | `sdcpp.Engine.videoProc` | Process registry | Registry owns `sd-cli` from before start until after exit |
| Chat process | `host.Processes` single slot | Process registry | Generalize to role + PID + start token while preserving llama stop semantics |
| Execution lifetime | HTTP callbacks, job callbacks, `imageStartMu` | Resource authority + engine scope | Named lease per chat turn, image job, video render, audio request, and workflow stage |
| Workflow claims | Implicit stage history | Workflow scope | Parent lease plus stable component claims; stage 4 verifies the claim |
| Capacity observation | Registration-time `Release()` | Observation layer + authority | Pending remains charged until a post-start observation reflects it, or it becomes `unknown` |
| Allocation advice | `resource.Allocate` | Resource authority | Projection from the same authority facts; second policy removed |
| Manual budget provenance | Settings/config paths | Configuration authority | Explicit source survives reload and workflow; authority consumes it as a ceiling |

The target has one writer per state. Direct callers may request a transition; they may not write the registry or declare capacity recovered.

---

## 3. Work packages

Each package is a production slice with its own acceptance gate. Packages may be developed in parallel only where their ownership boundaries do not overlap; the gate order in §5 is authoritative.

| Package | Ownership change | Closes | Primary evidence |
|---|---|---|---|
| `IP-1` Commitment correctness | Authority commits the peak it judged; RAM-only residents participate in blocking | D-1, D-9 | Converted known failures; accounting regression |
| `IP-2` Authority lease and registration contract | Authority becomes the only registration writer; lease kinds and capacity states exist | RLC-08, RLC-09, RLC-15, RLC-30, RLC-31, RLC-33 | Behavioral authority tests; public-bypass retirement |
| `IP-3` Process registry | One registry owns every child, one waiter, start tokens, exit events | RLC-10, RLC-34, RLC-35, RLC-37 | Registry tests; llama parity; image/video adoption |
| `IP-4` Image lifecycle transaction | Image load/reload/unload is one authority transaction; engine reports exit only | D-2, D-4, D-5; RLC-21, RLC-32, RLC-39 | Converted image tests; reload matrix |
| `IP-5` Video execution ownership | Video has a load/execution lease and a registry-owned `sd-cli`; stopping `sd-server` invalidates its registration | D-6, D-7; RLC-11, RLC-36 | Video lease/registration tests; single-Wait test |
| `IP-6` Chat and audio execution | Chat registration verifies process identity; chat and audio execution scopes carry leases | D-3; RLC-11a, RLC-43 | Behavioral chat unload race; audio lease test |
| `IP-7` Workflow scope and claims | Workflow parent, stage scopes, and component claims are explicit | D-8; RLC-12, RLC-22, RLC-40–RLC-43 | Workflow claim tests; visible stage evidence |
| `IP-8` Observation, recovery, and manual limits | Pending/reflected/released, recovery measurement, optional telemetry, provenance, shutdown join | RLC-15a–RLC-19, RLC-24–RLC-29, RLC-44–RLC-56 | Observation/recovery tests; cross-workload scenario |

**Gate 1 status.** `IP-1` is implemented. `grant` now commits `committedVRAMMB()` (the larger of `RequiredVRAMMB` and the normalized projected peak), and RAM-only residents participate in blocking when the request has a RAM requirement. A resident stage's projection includes its newly added adapter bytes before the comparison, so the resident-stage commitment is not just the fixed overhead. The D-1/D-9 tests graduated to `internal/resource/commitment_test.go`; additional tests pin full-load, resident-stage, resident-stage-with-adapter, replacement, and the RAM blocker reason. `IP-2`–`IP-8` remain unstarted.

### 3.1 Known-failure reconciliation

The seven known-failure tests were the closure targets for the first implementation slices. Gate 1 closed D-1 and D-9; Gate 2 closed D-3 and D-6; D-4 and D-5 graduated under their own behavioral tests. No tagged known-failure test remains. Assertions do not change; a build tag is removed only when the owning package makes the unchanged assertion pass.

| Defect / test | HRA-P1 clauses | Required ownership change | Acceptance evidence | Package |
|---|---|---|---|---|
| **D-1** `TestGrantReservesProjectedPeak` | RLC-14, RLC-15, RLC-45 | `grant` commits the need admission actually judged: `max(RequiredVRAMMB, ProjectedPeakVRAMMB)`. The commitment records `pending`; `AdmissionLease.Release()` remains reservation-only and does not imply ownership release. No caller recomputes the peak. | The unchanged test passes in the default suite. Add a resident-stage check that a second admission cannot spend the committed peak. `TestResidentExecutionDoesNotRechargeWeights` and the cross-workload scenario must still pass. | `IP-1` |
| **D-2** No dedicated known-failure test yet | RLC-15, RLC-15a, RLC-16, RLC-30, RLC-46 | Registration no longer releases the load commitment at registration time. The authority exposes distinct `ObserveAllocation`, `CommitResident`, and `ReleaseReservation` transitions. A post-start observation moves the commitment to `reflected` without closing the owner; failure leaves `pending` or `unknown`. | `TestPendingCommitmentSurvivesUntilObservation`; `TestPendingExecutionChargedUntilReflected`; `TestReflectedExecutionKeepsLeaseActive`; `TestLastProcessExitDoesNotFreeUntilMeasured`. These are new behavioral tests, not source guards. | `IP-2`, `IP-4`, `IP-8` |
| **D-3** `TestStaleChatProcessCannotRegisterAfterStop` / `TestCurrentChatProcessRegistersAfterHealth` | RLC-08, RLC-09, RLC-31, RLC-35 | The monitor captures the adopted handle (role + PID + start token) before the health check. The health-success registration route verifies that handle is still current before mutating the resource registry. No direct `RegisterModel` call remains in the monitor. | Stale registration is refused after the adopted process is stopped and leaves no resident; the current-process positive case still registers. Both pass in the default suite. | `IP-3`, `IP-6` |
| **D-4** `TestImageReloadFailureInvalidatesPreviousRegistration` / `TestImageReloadSuccessLeavesExactlyOneRegistration` | RLC-08, RLC-09, RLC-21, RLC-27, RLC-31, RLC-32, RLC-39 | A failed load/reload invalidates image registrations when the engine is no longer running; a successful replacement leaves exactly one. The current slice does not yet make reload a full authority transaction, so `registerImageModel`/`unregisterImageModels` remain the app-side route. | Both tests pass in the default suite: failed reload has zero registrations; successful reload has exactly one and the engine reports the new model. | `IP-4` |
| **D-5** `TestShutdownJoinsGeneration` and `TestEngineShutdownReportsExitWithoutUnregistering` | RLC-11, RLC-11a, RLC-13, RLC-21, RLC-27, RLC-28, RLC-31, RLC-39 | The engine owns an execution scope: cancel, join the active generation, stop and reap the process, report exit through its hook. Registration invalidation is an app/domain transition, never an engine loop. App shutdown reaches the join through `images.Close`. | Both named tests pass in the default suite, including the blocking transport; the engine leaves an image registration intact and invokes its exit hook once. `TestShutdownJoinsAllExecutionScopes` remains planned. | `IP-4`, `IP-8` |
| **D-6** `TestStopServerUsesOneRegisteredWaiter` / `TestStopServerRetainsOwnershipWhenRegisteredStopFails` | RLC-10, RLC-34, RLC-37 | The registry owns `sd-server` `Wait`. `stopServer` requests the registry stop and waits for the recorded waiter's exit signal; it never calls `Wait` on an `os.Process`. The adapter may keep a handle for cancellation, but not for waiting. | A recording registry proves one stop request and exactly one registered waiter call; a failed registry stop retains the process handle. Both pass in the default suite. | `IP-3`, `IP-5` |
| **D-7** `TestVideoExecutionHasLease` / `TestVideoLeavesNoImageRegistration` / `TestVideoProcessRegisteredAndReleased` | RLC-11, RLC-11a, RLC-15, RLC-17, RLC-36, RLC-39 | Video is a named execution scope, not an early lease release. The load/execution lease spans `SelectModel` and `GenerateVideo` until `sd-cli` exits. `sd-cli` is adopted immediately after start, before generation output is processed, and removed after its recorded waiter observes exit. Stopping `sd-server` for video reports the exit through the app hook, which removes the image registration. | The three named tests pass in the default suite. The lease is still open when generation starts and closes exactly once after it returns; normal exit and cancellation both release the `sd-cli` registration; no image registration remains after the server stops. The invalidation route is the app hook, not a full authority transaction. | `IP-5` |
| **D-8** `TestWorkflowRetainedUpscalerIsClaimed` / `TestSharedComponentRefcountAcrossStages` / `TestWorkflowReusesEngineAcrossStages` | RLC-12, RLC-22, RLC-40, RLC-41, RLC-42, RLC-43 | The workflow owns a `resource.WorkflowLease` parent scope. Stage 2 records an explicit upscaler claim under a stable key; stage 4 verifies that claim before reuse, and refuses when it is absent. Reclamation is refused while the parent scope or a component claim is active. The key includes engine instance and process identity, component role, normalized path, weight type, and the companion set. | The three named tests pass in the default suite. A claim is visible to stage 4, two stages share one component through a reference count, and companion order/path spelling do not change the key while a replacement engine process does. `TestWorkflowParentLeaseRefusesReclamation`, `TestWorkflowStage4RefusesAbsentUpscalerClaim`, and `TestWorkflowEngineReplacementInvalidatesScope` cover the refusal and stale-engine paths. | `IP-7` |
| **D-9** `TestRAMOnlyResidentBlocksRAMAdmission` | RLC-14, RLC-20 | `blockingModels` includes residents with `GPUMemoryMB <= 0` when they consume RAM, and sorts/blockers by the capacity that actually blocks. The admission reason names the RAM-only resident and the RAM numbers. | The unchanged test passes in the default suite. Add a regression where a RAM-only resident coexists with a GPU resident and the RAM request is refused with both relevant blockers. | `IP-1` |

**Closure rule.** A defect row is complete only when its named evidence passes in the default suite. Until then, the known-failure test stays tagged and failing, and no production behavior is described as fixed.

**Gate 1 closure record.** D-1 and D-9 graduated from `knownfailures` into `internal/resource/commitment_test.go`; their original assertions pass in the default resource suite. The new `TestCommitmentUsesComparableProjectedPeakAcrossPaths` covers full-load, resident-stage, resident-stage-with-adapter, and replacement accounting, and `TestRAMOnlyResidentReasonNamesRAM` covers the user-facing RAM blocker name. D-3, D-4, D-5, D-6, and the untested D-2/D-7/D-8 remain open.

### 3.2 Planned-test reconciliation

These invariants have no passing behavioral test today. They are not converted known failures; they need new evidence after the owning boundary exists. The Phase-one source guards have been retired.

| Planned test | HRA-P1 clauses | Required ownership change | Acceptance evidence | Package |
|---|---|---|---|---|
| `TestPendingCommitmentSurvivesUntilObservation` | RLC-15, RLC-15a, RLC-16, RLC-46 | Authority keeps a load/execution commitment `pending` across registration. A post-start snapshot transitions it to `reflected`; only a separate ownership transition closes the lease. Ambiguous observation becomes `unknown`, not released. | A fake observation source returns no movement, partial movement, then attributable movement. Pending remains deducted through the first two; reflected stops the extra deduction but keeps the owner; an ambiguous sample never reaches `released`. | `IP-2`, `IP-8` |
| `TestVideoExecutionHasLease` | RLC-11, RLC-11a, RLC-36, RLC-43 | Video admission creates a named execution scope before `sd-cli` starts. The lease is held through render and released exactly once after exit. The image load lease and the video execution lease are distinct. | A fake video adapter blocks at start; the test observes an active video lease while blocked and no open lease after exit. Cancellation also closes exactly once. | `IP-5` |
| `TestVideoLeavesNoImageRegistration` | RLC-09, RLC-31, RLC-39 | When video stops `sd-server` to free the GPU, the authority invalidates the image registration in the same transition. `SelectModel` cannot leave a resident image claim behind. | Register an image model, enter video mode, stop the server; the image registration is invalid/removed and no resident is reported. After video exits, no image registration reappears without a new load. | `IP-5` |
| `TestWorkflowReusesEngineAcrossStages` | RLC-13, RLC-38, RLC-39 | A workflow stage records the engine identity it uses. Same-model stages keep the same identity and engine instance; replacement invalidates the workflow scope rather than silently continuing on a stale pointer. | Run four stages with one engine; all stage records name the same engine key. Force a reconfigure between stages; the next stage receives an explicit invalidated result and does not run. | `IP-7` |
| `TestWorkflowRetainedUpscalerIsClaimed` | RLC-12, RLC-22, RLC-42 | Stage 2 creates a stable upscaler claim under the workflow parent. Stage 4 verifies the claim before reuse. Reclamation refuses while the claim is held. | Stage 2 increments the claim; stage 4 observes the same key; releasing the workflow drops it to zero. Remove the claim before stage 4 and the stage reloads or refuses, not silently reuses. | `IP-7` |
| `TestOptionalTelemetryAbsenceDoesNotRefuse` | RLC-25, RLC-53, RLC-54 | Optional observations are absent values, not zero-valued hard gates. Physical VRAM/RAM and explicit manual limits remain the only default gates. No new telemetry source is required for phase one; if fields do not exist yet, the test uses absent values and asserts no refusal. | Admission with no CPU utilization, temperature, or per-process VRAM yields a decision based only on physical memory and limits. A stale required physical value still refuses conservatively. | `IP-8` |
| `TestManualBudgetProvenanceSurvivesReload` | RLC-24, RLC-49, RLC-50, RLC-51, RLC-52 | Configuration authority records `manual`, `automatic`, `manual-clamped`, `live-free`, `negative-reserve`, or `unset`. The authority consumes the source as a ceiling; reload and workflow do not collapse provenance into a bare number. | Table-driven cases for manual, automatic, clamped manual, zero/live-free, negative reserve, and unset survive a settings reload and a workflow stage. Existing execution continues; new leases use the new ceiling; over-limit residents are reported, not silently re-admitted. | `IP-8` |
| `TestPublicRegistrationBypassesAuthority` | RLC-08, RLC-31, RLC-33 | Remove `App.RegisterLoadedModel` / `App.UnregisterLoadedModel` from the public mutation surface, or turn them into authority requests. `ResourceManager` mutation methods become package-private or authority-only. | The app no longer exposes a direct registry mutation; an authority request is the only mutation route. A supplementary source scan may assert no direct `RegisterModel`/`UnregisterModel` call outside the resource package, but behavioral authority tests are the acceptance evidence. | `IP-2` |

### 3.3 Remaining invariant coverage

These clauses have no dedicated named test in §3.1 or §3.2. Some are partially exercised by existing baseline tests; none is proven by the characterization gate. They are grouped by the ownership boundary that must exist before meaningful acceptance evidence is possible.

| Clauses | Untested invariant | Current evidence | Required ownership change | Acceptance evidence |
|---|---|---|---|---|
| RLC-17, RLC-19, RLC-30 | Exit does not prove recovery; replacement credit is conditional; the three release events are distinct | Replacement footprint is credited from registration; registration releases the lease immediately | Authority separates reservation release, ownership release, and capacity recovery; replacement credit is committed only after observed exit | `TestLastProcessExitDoesNotFreeUntilMeasured`; `TestReplacementCreditRequiresObservedRelease`; `TestReleaseVocabularyIsDistinct` |
| RLC-18, RLC-44 | Physical free is the hard gate; execution charges only its increment | Admission baseline and resident-stage characterization | Execution commitment derives from total workload peak minus resident footprint; managed budgets never replace physical free | `TestResidentExecutionDoesNotRechargeWeights`; `TestManagedBudgetCannotExceedPhysicalFree`; cross-workload scenario |
| RLC-21 | A successful load cannot become resident without a lease; a failed transfer leaves no dual lease | Known failures D-4/D-5 expose parts of the gap | Load, resident, and execution transitions are atomic authority operations | Reload success/failure matrix; `TestNoResidentWithoutResidentLease`; `TestFailedTransferLeavesNoDualLease` |
| RLC-23 | Workflow reserves only retained or explicitly committed capacity, not every future stage | `admitWorkflowDelta` admits per stage, but no workflow parent/commitment exists | Workflow parent holds only retained claims; stage commitments are created on advance | `TestWorkflowReserveOnlyCommittedComponents`; `TestWorkflowCancelBetweenStages` |
| RLC-26, RLC-47, RLC-48 | Historical observations never grant permission; predicted, observed, and confidence remain separate | `MeasuredResidentVRAMMB` reuses history for prediction, but no evidence model | Observation layer exposes age/confidence/completeness; authority may tighten on underprediction but never loosen on history alone | `TestPredictedObservedConfidenceRemainSeparate`; `TestHistoricalLowSampleCannotAdmit` |
| RLC-29 | Refusals are surfaced with blocker and numbers | Admission reasons exist; discarded-refusal paths are not enumerated | Every refusal returns through one result and user-visible event path; callbacks may not swallow it | `TestRefusalCarriesBlockerAndNumbers`; `TestDiscardedRefusalIsReported` |
| RLC-34, RLC-35 | Every child is registry-owned; identity is role + PID + start token | `host.Processes` has one slot and no start token; `sdcpp` keeps raw process fields | Generalize registry to multiple roles and tokenized identities while preserving llama behavior | Registry tests for llama, `sd-server`, and `sd-cli`; reused-PID test; `TestEveryChildHasOneOwner` |
| RLC-38, RLC-39 | Engine identity includes process/startup configuration; replacement invalidates scopes | `Engine.currentModel` and `LoadModel` replacement are not tokenized | Engine key is recorded on the scope; reconfigure/unload invalidates it before the new engine is visible | `TestWorkflowReusesEngineAcrossStages`; `TestReplacedEngineInvalidatesScope`; reload matrix |
| RLC-40, RLC-41, RLC-43 | Workflow parent, stable component key, one parent per batch | `imageStartMu` serializes but is not a lease; stage history is implicit | Workflow scope owns parent; batch items are child execution scopes; component key is explicit | `TestWorkflowParentClosesAfterBatch`; `TestWorkflowBatchItemsAreSeparateScopes`; `TestSharedComponentRefcountAcrossStages` |
| RLC-53, RLC-55, RLC-56 | Physical telemetry is the admission signal; stale/partial fails conservatively; a free increase is not recovery evidence | `TestStaleTelemetry` covers one refusal; no recovery observation test | Observation layer records source, age, completeness; authority rejects stale/contradictory physical values and keeps recovery `unknown` until measured | `TestStaleTelemetry`; `TestPartialPhysicalTelemetryFailsClosed`; `TestFreeIncreaseWithoutExitIsNotRecovery` |
| RLC-57 | Persistent calibration stays prohibited until phase one passes | No calibration path in this plan | No persistence or calibration work is permitted in `IP-1`–`IP-8`; a later gate must explicitly authorize it | Review gate: no calibration code or persisted observation store; phase-one acceptance suite green first |

---

## 4. Acceptance evidence rules

### 4.1 Evidence classes

| Class | Meaning | May close a defect alone? |
|---|---|---|
| `E1` Converted known failure | The existing `knownfailures` test, unchanged, moves into the default suite | Yes, when it exercises the real transition |
| `E2` Behavioral transition test | A test through the authority, registry, or engine boundary with fakes | Yes |
| `E3` Process-registry test | Exact wait count, start-token identity, stop/exit observation | Yes |
| `E4` Concurrency/race test | Repeated or `-race` execution of admission, reload, shutdown, or lease close | Yes, for the race it targets |
| `E5` Cross-workload scenario | The contract §12.1 sequence with chat, image, video, and workflow | Required for phase-one completion |
| `E6` Both-host workflow | Packaged desktop and browser dev host with isolated state root | Required only for user-visible completion |
| `E7` Source/API audit | Removal of a bypass or a forbidden call | Supplementary only; never the sole evidence |

A source scan is not `E2`. It may remain as a tripwire while a behavior is unobservable, but the defect stays open until `E2` or `E3` replaces it.

### 4.2 Evidence rules

- Every test runs with `-count=1` in the default package suite. The known-failure suite must shrink; it must never grow to accommodate new work.
- Assertions are not weakened, skipped, or made conditional on the current implementation. If an assertion is wrong, the contract must be amended first.
- A test may not reach into private state to prove a transition that the authority does not expose. The transition becomes part of the authority's real API, not a test hook.
- A process test uses a test double only at the process-registry boundary. It must still prove the production invariant: one waiter, one owner, one exit observation.
- A concurrency test must show the failed interleaving would have passed under the old behavior, so it is evidence of the fix rather than a tautology.
- No package is complete while two components can write the same state. A temporary adapter may translate calls, but it may not retain write authority.
- A package that changes user-visible behavior is not complete until the contract §15.1 both-host workflow is run. Until then it is source-complete, not user-visible complete.

### 4.3 What does not count

- A passing build or typecheck.
- A source guard that only checks for a forbidden string.
- A log line that shows the intended transition without an assertion.
- A test that calls `Release()` and assumes ownership ended.
- A UI screenshot without process and registration evidence.
- A successful generation while a lease or process remains open.

---

## 5. Implementation sequence and gates

### Gate 0 — Planning accepted

This paper is the gate record. No production code starts from a known-failure test alone; each starts from an owning package.

### Gate 1 — Commitment correctness (`IP-1`) — implemented

**Changes:** commit the projected peak; include RAM-only residents in blocking and reasons.

**Why first:** these are local authority corrections. They do not require the process registry or registration migration, and they remove two known failures without creating a second owner.

**Exit evidence:** `TestGrantReservesProjectedPeak` and `TestRAMOnlyResidentBlocksRAMAdmission` pass unchanged in the default suite; existing characterization tests still pass; no caller was given a new write path. `TestCommitmentUsesComparableProjectedPeakAcrossPaths`, `TestResidentStageCommitmentBlocksSecondAdmission`, `TestRAMOnlyResidentReasonNamesRAM`, and `TestRAMOnlyAndGPUResidentsBothAppear` pass in the same suite.

### Gate 2 — Authority contract and process registry (`IP-2`, `IP-3`)

**Changes:** introduce lease kinds and capacity states; make the authority the only registration writer; restrict the public registration surface; generalize the process registry to all children with role + PID + start token; preserve llama's stop and single-wait behavior.

**Why here:** every later package needs a truthful owner for registration and process lifetime. Migrating image or video before this gate would preserve the current split ownership.

**Exit evidence:** authority behavioral tests for registration requests and lease closure; registry tests for llama, `sd-server`, and `sd-cli`; reused-PID test; no direct public registration mutation; llama baseline tests unchanged.

**Implemented slice:** `internal/process` registry tests, host ownership tests, and llama/image/video adoption pass. Authority lease-kind and registration-writer work is not part of this slice.

### Gate 3 — Image and chat/audio ownership (`IP-4`, `IP-6`)

**Changes:** image load/reload/unload becomes one authority transaction; engine reports exit only; chat registration verifies process identity; chat and audio execution scopes carry leases; audio registration joins its load lease to the authority.

**Exit evidence:** D-3, D-4, and D-5 known failures pass unchanged; the chat unload race is behavioral; image success/failure matrix passes; audio execution lease closes exactly once; no registration outlives its process.

**Scope correction (2026-11-25).** The audio portion of this gate is not complete. Audio is promoted to the required Gate 8; the existing `TestASRStartConsultsAdmissionWhileChatLoaded` failure remains visible and is not evidence for this gate.

### Gate 4 — Video ownership (`IP-5`)

**Changes:** video gets a named execution scope and lease; `sd-cli` is registry-owned before start; `stopServer` requests stop and waits for registry exit; stopping `sd-server` invalidates its image registration.

**Exit evidence:** D-6 source guard replaced by a registry wait-count test; D-7 lease and registration tests pass; video cancellation and normal exit both close exactly once; no image registration remains after the server stops.

### Gate 5 — Workflow scope and claims (`IP-7`)

**Changes:** workflow parent lease; per-stage execution scopes; stable component keys; stage 2 claim and stage 4 verification; reclamation refuses while claimed; batch items are child scopes.

**Exit evidence:** D-8 tests pass; same-engine reuse stays stable; replacement invalidates the scope; cancel between stages starts no next stage; visible stage evidence names the claim and engine identity.

### Gate 6 — Observation, recovery, and manual limits (`IP-8`)

**Changes:** pending/reflected/released with post-start observation; recovery measurement; optional telemetry confidence; manual provenance; stale/partial physical telemetry refusals; shutdown joins all scopes.

**Exit evidence:** all contract §12 tests pass, including `TestCrossWorkloadLeaseSequence`; the known-failure suite is empty; no pending commitment survives the final reconciled measurement; no registration lacks a live process; no process remains alive after shutdown.

### Gate 7 — User-visible completion

Run the contract §15.1 workflow in both supported hosts with an isolated state root. Record native-surface differences and include process/registration evidence, not only screenshots. Until this gate runs, phase one is source-complete only.

### Gate 8 — Audio integration (required)

**Changes:** resolve the audio policy/test mismatch against the reserve contract; audit and unify ASR/TTS registration, execution leases, process ownership, cancellation, termination, and cleanup; add cross-workload evidence with chat, image, and video; exercise actual audio workflows in both supported hosts.

**Why this is a gate:** audio is a first-class workload. A red ASR test carried as a permanent exception, or an audio-only admission rule, would undermine the unified authority. The audio baseline is not optional cleanup after HRA-P1.

**Exit evidence:** the ASR test is resolved by a documented policy/test decision without weakening it; audio lifecycle and cross-workload tests pass; no audio registration or process outlives its execution; cancellation reaches the active audio operation and terminates its owned process; both hosts demonstrate ASR/TTS behavior with backend/resource evidence.

**Ordering:** after the current engine-lifecycle characterization and the remaining image/video cross-boundary work. Gate 8 may run alongside later lifecycle work, but HRA-P1 is not complete until it closes.

**Status (2026-11-25). Gate 8 remains open.** Voice quality is accepted — the operator heard local Piper in their running app and confirmed both read-aloud paths (`HRA-P1-TTS-LISTENING-ACCEPTANCE-2026-11-25.md` §5) — but that closes only the listening question. The gate closes when all four workstreams below carry their evidence:

1. **Cross-workload resource/lifecycle testing — backend complete, host image demonstrated, video unattempted.** Chat, image, video and audio against the one authority: audio residency appears as a blocker and disappears on release, no registration or child outlives its execution, cancellation reaches the active operation. **Backend criteria are met** (§7.6.4): `internal/audio/cross_workload_image_video_test.go` adds six default-suite tests covering audio→image refusal and release (including the reverse direction), a video render request observing audio residency and holding one lease until `sd-cli` exits, an audio child and a video `sd-cli` child coexisting in one registry so that stopping the video leaves audio intact, and cancellation confined to its own operation. Focused run, race run, `go vet`, and both builds are clean. The host run initially exercised the audio half only. **Correction (2026-11-25): the "no weights" finding was a lookup failure, not a machine limitation.** The scan covered `models/`, the configured `imageGen.modelsDir`, and `C:`/`D:`/`E:` model paths; it never looked at the owner's model drives. `H:\AI GGUF\Image models` holds 33 image models, 4 VAEs, 8 LoRAs and a Hunyuan video weight, each with a sidecar, and a packaged-host and browser-host generation both produced real 1024x1024 files (§7.6.8). The image half of the host evidence is therefore **demonstrated**; the video half remains unattempted, not blocked on assets. Sub-limitation unchanged: stale telemetry is proven against an injected snapshot, not observed against real Windows VRAM recovery.
2. **Packaged-host audio integration evidence — done (§7.6.3, §7.6.6).** The packaged WebView2 host demonstrated Piper synthesis with process and artefact evidence, and the operator heard it (§7.6.3). ASR was then enabled through the real settings UI, persisted `asrEnabled:true` / `qwen3-asr-1.7b-local`, spawned its owned `llama-server` child, reported **ready**, and transcribed the operator's live speech into the composer: “Testing. Okay, the quick brown fox jumps over the lazy dog. Package host speech recognition test.” (§7.6.6). Driving is coordinate-based because the accessibility tree exposes only the WebView2 host panes; that is a tooling cost, not a blocker. The re-apply self-block observed in that host was fixed the same day and re-verified there (§7.6.7).
3. **Kokoro/MeloTTS coverage — blocked on assets.** `models/kokoro` and `runtime/tts` are absent, only `scripts/kokoro-tts.py` exists, `python -c "import melo"` fails with `ModuleNotFoundError`, and no MeloTTS service listens on `127.0.0.1:8888`; the host capability records report both unavailable for those reasons (§7.6.5).   Neither provider has host evidence, and no asset was installed or service started. Closing this needs assets or a launcher, plus a run through the same ownership, residency and capability path Piper went through. **Owner guidance (2026-11-25): neither provider may become a release blocker merely for completeness, because Piper already provides accepted speech; whether they belong in the initial release is a deliberate release decision, not a Gate 8 prerequisite.** That decision is not recorded as made.
4. **ASR admission policy/test discrepancy — resolved 2026-11-25 (§7.6).** The test, not the policy, encoded the wrong expectation: at 4500 MB free the 4096 MB ASR request is admitted through the documented soft-reserve allowance, leaving 404 MB free above the 128 MB hard floor. The policy is unchanged and applies the same rule to every workload without an explicit reserve. The test keeps its refusal assertion against a genuinely refused fixture (4000 MB free) and asserts the admitted case explicitly with the exact accounting; no audio-only rule was added.

**Two default-suite failures surfaced during this workstream were repaired as test-only isolation fixes (2026-11-25).** They were characterized from their actual failure output first; neither was skipped, weakened, or attributed to the cross-workload implementation, and production code, resource policy and the ASR assertion were not touched.

- `TestBuildPiperVoiceCatalogMarksInstalledAndMissingVoices` (`app-tts-catalog_test.go:34`). `buildPiperVoiceCatalog` always merges `audio.TTSModelPathFor(audiotypes.TTSPiper)` into the catalog, so the real installed voice `en_US-lessac-medium` was advertised and won the default slot ahead of the fixture's `present` voice. The assertion was right; the fixture was wrong — it assumed an empty model directory and never controlled the default-voice input. **Fixed (test-only):** the fixture now points `HOWL_TTS_PIPER_MODEL` at a path inside its own `t.TempDir()`, so the default voice is an input the test owns. The availability and default assertions are unchanged.
- `TestMicrophoneServiceUsesWAVDataBindingAndComposer` (`app_audio_data_test.go:187`). Its source guard forbade the literal `createObjectURL` in `frontend/dist/js/audio-service.js`, which now uses it in `_ttsSourceFor` to convert the backend's base64 WAV data URL from `GetTTSAudioData` into a playable blob URL. The guard was written for the earlier defect where capture passed an object URL to `TranscribeAudio(tempPath)` as if it were a file path, so it conflated that handoff with legitimate playback. **Fixed (test-only):** the guard now forbids `saveBlobToTempFile` and `TranscribeAudio(tempPath)` — the actual browser-URL-as-filesystem-path handoff — while keeping the `MediaRecorder`, `#chat-input`, `AudioWav.encodeWav16`, `TranscribeAudioData`, `#message-input` and script-order assertions.

After those repairs the suite failed only on the preserved ASR test; the disposition in item 4 closed that too. `go test -count=1 -p 1 ./...` now passes every package (50 ok; `howl-chat` 17.1 s, `internal/audio` 13.6 s). The default parallel invocation hit a transient Go `cmd/link` panic on two different packages (`internal/image/deps`, then `internal/backend/loader`); each passed when rerun and the serialized run is clean, so this is a toolchain flake, not a test result.

HRA-P1 is not complete until these four close.

### 5.1 Dependency rules

- `IP-1` may start immediately and must not wait for the registry.
- `IP-2` and `IP-3` must agree on the process key before either is considered complete.
- `IP-4` and `IP-6` depend on `IP-2`/`IP-3`; do not migrate image or chat registration before the authority and registry can represent the transition.
- `IP-5` depends on `IP-3` and the image invalidation path from `IP-4`.
- `IP-7` depends on the authority lease model and the image/video execution scopes.
- `IP-8` spans all packages; its tests are the final gate, not a cleanup phase.
- Gate 8 (audio) depends on the Gate 3 audio lease work, the Gate 6 cross-workload model, and real host access. It cannot be skipped or replaced by a passing backend suite.

### 5.2 Stop conditions

Stop and revise the plan if any of these occurs:

- a package needs a test-only production hook to make its evidence possible;
- llama stop/restart behavior changes without a contract amendment;
- two components can still write a registration or declare recovery;
- a test can pass only by weakening an assertion or keeping a source guard;
- a package leaves a lease or process open on an error path;
- a change is justified only by a label such as “priority” rather than measured capacity and ownership.

---

## 6. Production surface changes

### 6.1 Stays

- `AdmissionRequest`, `AdmissionEstimate`, and `AdmissionDecision` remain the decision vocabulary.
- `AdmissionLease` remains the reservation handle; its release does not unregister or prove recovery.
- `ResourceManager` remains the read model for loaded models and hardware snapshots.
- `sdcpp.Engine`, llama-server, the audio loader, and the job manager remain execution adapters.
- `RunImageWorkflow` remains the sequential workflow owner; stage order and output behavior stay unchanged.
- llama's one-handle, one-waiter, ownership-safe clearing semantics are preserved in the generalized registry.

### 6.2 Removed or rerouted

| Current path | Disposition |
|---|---|
| `Engine.Shutdown` removing image registrations | Removed; engine reports exit, authority invalidates |
| `App.UnloadImageModel` removing registrations after shutdown | Replaced by authority unload transaction |
| `App.RegisterLoadedModel` / `App.UnregisterLoadedModel` | Removed from public mutation surface or converted to authority requests |
| `resource.GetManager().RegisterModel` / `UnregisterModel` from callers | Restricted to the authority package |
| `video.stopServer` calling `Wait()` | Replaced by registry stop + exit signal |
| `serverProc` / `videoProc` as lifetime authorities | Adapter state only; registry owns lifetime |
| `resource.Allocate` as a second policy | Removed or projected from authority facts |
| `imageReclaim` reserve policy | Reads authority facts; no independent reserve |
| `host.Processes` single-child slot | Generalized registry with roles and start tokens |

### 6.3 Compatibility constraints

- Chat streaming and cancellation behavior must remain unchanged.
- Image generation and workflow output must remain unchanged for the same inputs and placement.
- Audio ASR/TTS loading remains serialized by the audio manager; the resource lease joins that lifetime rather than replacing it.
- A failed unload still reports failure and keeps ownership; it may not claim capacity recovery.
- Manual/image budget clamping keeps its existing user-visible warnings until the provenance work explicitly replaces them.

---

## 7. Verification boundary

### 7.1 What this plan verifies

- Every known failure has one owning package, one set of HRA-P1 clauses, and acceptance evidence that cannot be satisfied by a source guard alone.
- Every planned test has a named ownership boundary and a package.
- Every HRA-P1 clause not covered by the closed gate appears in §3.3 or a named row in §3.1/§3.2.
- The removal list identifies the competing writers that must disappear.

### 7.2 What this plan does not verify

- Gate 1 changed `internal/resource/admission.go` and graduated two tests; no other production package has started.
- No process registry, lease kind, or authority API exists yet.
- No phase-one acceptance test was run.
- No user-visible workflow was run.
- No concrete Go signatures, package layout, migration size, or schedule are asserted here; those belong to the implementation change that owns the package.

### 7.3 Gate record

Gate 1 is implemented: D-1 and D-9 pass their original assertions in the default resource suite, and the new comparability/reason tests pass. At Gate 1 closure, five known-failure tests remained tagged and all eight planned invariants remained planned. This plan does not add production instrumentation to make them pass. The next production change was Gate 2 (`IP-2`/`IP-3`), with its own diff and evidence.

**Gate 2 process-registry slice record.** The process registry now exists in `internal/process`, keyed by role + PID + start token, with one recorded waiter per child and token-safe release. The host chat owner, `sd-server`, and `sd-cli` all adopt through it; the image domain receives the composition root's registry. `stopServer` delegates to the registry's recorded waiter and no longer calls `Wait`; the chat health-success route refuses registration when the adopted process is no longer current. D-3 and D-6 now have graduating behavioral tests: `TestStaleChatProcessCannotRegisterAfterStop`, `TestCurrentChatProcessRegistersAfterHealth`, `TestStopServerUsesOneRegisteredWaiter`, and `TestStopServerRetainsOwnershipWhenRegisteredStopFails`. D-4, D-5, and the `IP-2` authority contract remain open; no reserve policy or calibration changed. This is Gate 2 process-registry implementation and D-3/D-6 behavioral closure, not full lifecycle acceptance.

**Record correction.** An earlier interim note said no existing files changed. That is superseded: this slice modified production code, tests, and both architecture documents. The production-code scope remains the process-registry slice described above.

**Gate 2 checkpoint boundary.** This checkpoint closes only the process-registry slice and the D-3/D-6 behavioral tests. D-4 and D-5 remained separate lifecycle findings with their existing acceptance criteria; generation joining and engine-unregistration ownership are not claimed as solved here. Passing backend suites does not establish packaged-WebView behavior for shutdown, reload, or cancellation, which remains a separate acceptance layer. This checkpoint does not authorize a general resource-management rewrite.

**D-4 closure record (separate from Gate 2).** D-4 is behaviorally closed for the tested failed-load boundary. A failed image load/reload invalidates image registrations only when `IsRunning()` is false, so a process that survived a failed start keeps its registration. The existing D-4 failure test and the new success matrix pass in the default suite.

Remaining boundaries: full reload authority-transaction migration is not claimed; packaged-WebView behavior remains unverified; and the audio ASR failure remains the separately classified baseline issue. This evidence covers the tested failed-load path, not every possible concurrent reload/unload interleaving; that broader claim needs its own evidence if it becomes part of later acceptance. D-4 closes without reopening Gate 2 or expanding into a general lifecycle rewrite.

**D-5 closure record (separate from Gate 2 and D-4).** D-5's two named behavioral tests pass in the default suite. Shutdown cancels and joins the active generation scope before stopping the child; the HTTP cancellation helper makes that join effective even when a transport ignores the request context. The engine now reports exit through a hook and no longer removes image registrations itself; the app/domain performs invalidation.

**D-5 acceptance boundary.** Not yet proven: every execution scope is joined during shutdown; the full application shutdown matrix; concurrent reload/unload/shutdown interleavings; and packaged-WebView behavior. Releasing a local execution scope when a transport ignores cancellation does not prove the remote backend stopped processing; the report claims only that the local scope returned. Remaining work follows the planned acceptance matrix in contract §12, not another broad refactor. The audio admission baseline remains separate.

**D-7 closure record (separate from D-5).** Video now resolves the model and keeps the admission lease open across `GenerateVideo`; the synchronous path defers release, and the job path returns the release to the job manager, which runs it after every generation outcome. The `sd-server` stop before `sd-cli` reports process exit through the existing engine hook, so the app removes the now-stale image registration. `sd-cli` is adopted into the process registry immediately after start and removed by its deferred unregister after the recorded waiter observes exit; cancellation follows the same removal path.

Evidence: `TestVideoExecutionHasLease`, `TestVideoLeavesNoImageRegistration`, `TestVideoProcessRegisteredAndReleased`, and `TestVideoProcessReleasedOnCancel` pass in the default suite and in the race run for the touched packages. This closes the named D-7 boundary only. It does not claim a full authority registration transaction, a complete video shutdown matrix, or packaged-WebView behavior. The audio admission baseline remains separate.

**D-8 closure record (separate from D-7).** The workflow now opens one `resource.WorkflowLease` parent scope for the request and ends it before the terminal reclamation decision. Stage 2 acquires the retained upscaler's component claim; stage 4 verifies it and refuses when absent. The claim key includes engine instance and process identity, the `upscaler` role, the normalized path, weight type, and the companion set. Reclamation is refused while the parent lease is active, which preserves the component across batch item boundaries; component claims themselves are reference-counted and idempotently released.

Evidence: `TestWorkflowRetainedUpscalerIsClaimed`, `TestSharedComponentRefcountAcrossStages`, `TestWorkflowReusesEngineAcrossStages`, `TestWorkflowParentLeaseRefusesReclamation`, `TestWorkflowStage4RefusesAbsentUpscalerClaim`, and `TestWorkflowEngineReplacementInvalidatesScope` pass in the default suite. This closes the named D-8 boundary only. It is a lifecycle ownership key, not a new admission-capacity lease: stages still use the existing delta admission. The full workflow acceptance workflow has not been clicked through in either supported host; that remains Gate 7 evidence.

### 7.4 Cross-boundary acceptance baseline (frozen 2026-11-25)

This audit is the current acceptance baseline for the interaction layer. It is not a reopening of D-7 or D-8.

**D-8 remains valid.** The named ownership and claim tests establish what they claim: the parent scope exists, stage 2 acquires the component claim, stage 4 verifies it, keys are stable across path/companion spelling, and a replacement engine invalidates the scope. The findings below are missing integration guarantees, not failures of those tests. D-7's named video ownership tests remain behaviorally verified for the same reason.

**Cross-workload is not an admission test.** `TestCrossWorkloadLeaseSequence` must establish the actual sequence of ownership and capacity transitions, using a stateful fake runtime with physical free/free RAM, resident registrations, process identities, and pending/reflected commitments. A passing admission decision does not prove that video and a workflow cannot interfere with the same engine or process. The evidence must show which scope owns the engine at each boundary, when a process may be stopped, when a registration may be invalidated, and that every lease closes exactly once.

**Characterization order.**

1. Workflow cancellation between stages and during an active stage.
2. Engine unload/reconfigure during stage 1 and stage 3.
3. Shutdown joining active and between-stage workflows.
4. Video/workflow overlap, including server invalidation and admission transitions.

Engine lifecycle comes before shutdown integration deliberately: otherwise a shutdown failure cannot be attributed between cancellation, stale-engine use, and shutdown ownership.

**Acceptance status.**

| Boundary | Status |
|---|---|
| D-7 named video ownership tests | Behaviorally verified |
| D-8 named workflow ownership tests | Behaviorally verified |
| Structured reclamation-refusal diagnostics | Behaviorally verified |
| Workflow cancellation | Backend execution-scope implemented and unit-tested; host acceptance open |
| Engine lifecycle during workflow | Open |
| Shutdown joining all execution scopes | Open |
| D-7/D-8 cross-workload interaction | Open |
| Audio admission test/policy disposition | Resolved 2026-11-25 — the test treated the 10% reserve as hard; the documented soft-allowance policy is correct (§7.6) |
| Audio lifecycle characterization (owners, execution cancellation) | Characterized and unit-tested (§7.6.1); TTS child-process ownership gap open |
| Audio lifecycle integration (ASR/TTS ownership, leases, cancellation, termination) | Open — required HRA-P1 completion gate; the gaps are named in §7.6.1 |
| Audio cross-workload interaction with chat/image/video | Open — required HRA-P1 completion gate (§7.6) |
| Audio host acceptance | Both hosts demonstrated: browser dev host ASR transcription + Piper TTS (§7.6.5); packaged host ASR transcription + Piper TTS (§7.6.3, §7.6.6); re-apply idempotence fixed and re-verified (§7.6.7). Kokoro/MeloTTS uninstalled, optional by owner guidance |
| TTS end-to-end synthesis on this machine | Demonstrated: packaged host (§7.6.3) and browser host (§7.6.4, §7.6.5), with real WAV artefacts |
| Packaged WebView2/browser-host workflow | Open |
| HRA-P1 overall completion | Blocked by audio integration and remaining lifecycle work |

“Behaviorally verified” means the named tests pass in the default suite. It does not mean host-level or cross-boundary acceptance.

**Structured refusal implementation.** When reclamation is refused because a workflow parent lease is active, the guard now records an `imageReclaimEvent` with `Blocked`, the refusal reason, the workflow engine identity, and the active claim count. The reclaim report exposes the same fields and a blocked summary, so the UI no longer sees only a console line or a stale prior decision. `TestWorkflowParentLeaseRefusalIsStructured` passes in the default suite. This does not change D-8's claim semantics or reopen its named boundary.

The four cross-boundary interaction areas remain characterization work; no general resource-manager rewrite is implied.

**Diagnostics closure record.** The structured-reclamation-refusal slice is complete at its named boundary. `TestWorkflowParentLeaseRefusalIsStructured` passes in the default suite and under `-race`; it asserts the blocked flag, reason, workflow engine identity, active claim count, and blocked summary. The slice does not alter D-7/D-8 ownership semantics, admission behavior, or the reclamation thresholds. It is not host-level acceptance evidence.

### 7.5 Workflow execution-scope cancellation (backend, 2026-11-25)

**Previous characterization retained as evidence.** Before this slice, `RunImageWorkflow` took no context, every stage call used `context.Background()`, `DGLUpscale` and `NativeUpscale` ignored their context for ownership, workflow stages had an empty job identity, and the parent lease was ownership rather than a cancellation token. The characterization tests that fixed those facts remain, and `TestWorkflowStageHasNoCancellationIdentity` still records that an empty `CancelJob` cannot be the cancellation route.

**Implemented boundary**

- `RunImageWorkflow` creates one `imageWorkflowExecution` per request. Its context is the parent of every stage generation; the execution is registered for `CancelImageWorkflow` before the workflow parent opens.
- Stage 1 uses `GenerateTxt2ImgContext`, stage 3 uses `GenerateImg2ImgNativeControlContext`, and stages 2 and 4 pass the workflow context to `NativeUpscale`/`DGLUpscale`. `beginGenerationWithContext` makes that context the parent of the engine generation handle, so cancellation reaches the active request instead of being ignored.
- The workflow checks the execution context before every batch item and before starting stages 1, 2, 3 and 4, including after a model load. `acceptStage` is the lock-linearized gate: if cancellation wins, the stage is refused; if the stage wins, the cancellation is observed by that stage's execution context. Once cancellation is accepted, no later stage begins.
- `CancelImageWorkflow` accepts cancellation once, then waits for the execution `done` channel. The workflow marks the stage loop finished before cleanup, so a cancel racing normal completion is either accepted before that point or refused, and `done` closes once.
- Cleanup remains in the existing order: the parent `WorkflowLease` closes, component claims release through their existing defers, and one terminal reclamation decision runs. An accepted cancellation additionally calls `Engine.Shutdown`, which cancels and joins the active generation handle and kills/reaps the owned `sd-server` process before the result is published. If process termination cannot be established, the result reports that failure rather than claiming cancellation completed.
- The workflow panel exposes a cancel action while a run is outstanding and waits for that backend boundary before reporting.

**Evidence**

- `TestWorkflowActiveStageObservesCallerCancellation` — cancelling the caller context reaches the active DGL stage and makes it return.
- `TestWorkflowCancelledScopeRefusesNextStage` — a cancelled scope rejects the next stage through both the stage gate and the context, and names the cause.
- `TestCancelImageWorkflowWaitsForExecutionCleanup` — cancellation does not return before execution cleanup completes.
- `TestWorkflowCancellationCompletionRace` — cancellation and normal completion cannot accept cancellation more than once or leave `done` open.
- `TestWorkflowCancellationAfterCompletionIsRefused` and `TestWorkflowExecutionCompletionClosesOnce` — a late cancel is refused and completion closes exactly once.
- `TestWorkflowCancellationContextIsChildNotBackground` — a stage context is reachable from the workflow scope.
- `TestWorkflowCancelBetweenStages`, `TestWorkflowLeaseReleaseDoesNotStopActiveStage`, and `TestWorkflowParentReleaseIsIdempotent` remain valid parent/claim ownership evidence.

These pass in the default suite and in the focused `-race` run. This verifies the backend execution-scope boundary, not packaged-WebView host behavior. The fake engine test establishes local request cancellation; process termination is delegated to `Engine.Shutdown` and its existing join/reap tests rather than a real `sd-server` workflow run. D-7 and D-8 are unchanged.

### 7.6 Audio integration gate (required, 2026-11-25)

**Decision.** Audio is a first-class workload and a **required HRA-P1 completion gate**, not optional cleanup after HRA-P1. HRA-P1 is not complete until this gate closes. `internal/audio/TestASRStartConsultsAdmissionWhileChatLoaded` was resolved on 2026-11-25 by correcting the fixture's accounting, not by waiving, skipping, renaming, tagging or weakening its assertion; the decision is recorded below.

**Evidence standard.** This gate is not closed by turning `go test ./...` green. It closes when ASR and TTS are shown to participate in the same admission, ownership, cancellation, termination, cleanup and capacity-reporting guarantees as chat, image and video, with cross-workload and host evidence. Removing the ASR failure without a recorded reserve-contract decision is not progress toward this gate; neither is a passing audio unit suite without an owned process behind it.

**Policy/test disposition (resolved 2026-11-25).** The failure was a policy/test mismatch, not evidence that ASR bypasses admission. Under the authority's documented accounting at the fixture's own values — 4500 MB free, a 3000 MB resident chat model, and the 4096 MB Qwen ASR profile — the reserve is 450 MB (10% of effective free), the bounded soft allowance is 322 MB (reserve minus 500 MB, floored at 128 MB), 46 MB of it is used, and the start is admitted leaving 404 MB free above the 128 MB hard floor. `admission.go` documents the 10% reserve as a target rather than a hard rule and applies the soft allowance to every workload without an explicit reserve; contract §7.1 subtracts the reserve and then treats physical free as the hard gate. The test asserted refusal by treating the 10% as absolute, an expectation the frozen contract does not support, so **the test, not the policy, was wrong**. No policy, reserve, or production code changed.

`TestASRStartConsultsAdmissionWhileChatLoaded` keeps its refusal assertion and now uses a fixture that is genuinely refused (4000 MB free: 400 MB reserve, 272 MB allowance, 3872 MB usable, so 4096 MB does not fit), still naming the resident chat model. `TestASRAdmissionAdmitsThroughSoftReserveWhileChatLoaded` asserts the admitted case explicitly: 4500/450/4050 free/reserve/available, a 322 MB allowance, 46 MB used, 404 MB headroom, and the pipeline's own `admitAudio` accepting the same request. Neither test starts a recognizer; the admitted case stops at the admission boundary by design.

**Work package (ordered after the current engine-lifecycle characterization).**

| Work | Purpose | Required evidence |
|---|---|---|
| Audio baseline resolution | Decide whether the test or the policy is wrong against the intended reserve contract. | A recorded decision with the policy clause, fixture values, and expected outcome. If the test changes, its refusal intent is preserved or the admitted case is asserted explicitly. No audio-only admission rule. **Done 2026-11-25 — §7.6.** |
| Audio lifecycle audit | Trace ASR/TTS registration, execution leases, process ownership, cancellation, termination, and cleanup from the actual call sites. | A named owner gap list with competing writers; no silent second registration or process owner. **Done 2026-11-25 — §7.6.1 (F-A1…F-A5); gap closure still open.** |
| Cross-workload tests | Verify audio alongside chat, image, and video under shared resource pressure. | Stateful tests that assert admission, ownership, process identity, lease release, and capacity reporting—not just admission return values. |
| Host acceptance | Exercise actual audio workflows and resource transitions in the application. | Packaged WebView2 and browser dev host runs with isolated state roots, visible ASR/TTS results, and resource/process state read back from the backend or disk. |

**Constraints.**

- Audio admission goes through the same `resource.Authority` as chat, image, and video. No bypass, no audio-only reserve, and no workload-label exception.
- Audio load/execution ownership must have one authority and one mutation route. Existing serialization by the audio manager is preserved rather than replaced by a second admission policy.
- Cancellation and termination must reach the active audio execution and its owned process, with cleanup exactly once and refusals surfaced.
- Resource reporting must distinguish pending, resident, and executing audio work rather than relying on a registration-time release or an unmeasured assumption.
- The lifecycle audit must confirm the current owners and competing writers. The working assumption—`resource.Authority` for admission, the audio manager/pipeline for serialized ASR/TTS lifetime, and the process registry for child ownership—is not a substitute for that audit.

**Completion evidence.**

1. The ASR baseline is resolved by a recorded policy/test decision, and `go test ./...` is green without weakening the test. **Done 2026-11-25 — §7.6.**
2. The audio lifecycle and cross-workload tests pass with the same assertions that define the other workloads' boundaries. **Done 2026-11-25 — §7.6.4.**
3. Actual ASR and TTS workflows are exercised in both supported hosts with isolated roots; the visible result and backend resource state agree. **Done for ASR and TTS 2026-11-25 — §7.6.3, §7.6.5, §7.6.6:** both hosts transcribed speech and synthesized with Piper; Kokoro/MeloTTS remain blocked on assets.
4. Any native-surface difference in the browser host is recorded explicitly rather than treated as equivalent. **Done 2026-11-25 — §7.6.5.**

The audio admission disagreement is closed, ASR transcription plus Piper synthesis are demonstrated in both supported hosts, and the re-apply self-block recorded in §7.6.6 was fixed and re-verified (§7.6.7). Image generation is now demonstrated in both hosts with real weights (§7.6.8); video render remains unattempted. Kokoro/MeloTTS remain uninstalled, but by owner guidance (item 3) neither is a release blocker merely for completeness. This section records those facts.

#### 7.6.1 Ownership audit and characterization (2026-11-25)

**Scope.** This slice audits the audio surfaces from their real call sites and pins what the current code actually guarantees. It changes no admission policy, no provider behavior, and no test that defines another workload's boundary. Gate 8 is an integration gate: its evidence standard is the same ownership, cancellation, termination, cleanup and reporting standard the image and video slices had to meet, not a green suite.

**Ownership as implemented.**

| Surface | Call site | Who owns it | Termination guarantee |
|---|---|---|---|
| ASR decode settings and admission | `Pipeline.initASR` → `admitAudio` | Audio pipeline asks `resource.Authority`; no bypass | Refusal names the resident blocker; no provider is constructed on refusal |
| Local Qwen ASR server | `internal/audio/asr/llama/llama.go` (`proc.Start`, `stopLocked`) | ASR recognizer, itself owned by `runtime.ModelLoader`; server assigned to the application job before it is resumed | `KILL_ON_JOB_CLOSE` covers a crash; release kills the tree and waits |
| ASR/TTS registration | `Pipeline.registerAudio` | `resource.Manager` records the model, then the admission lease is released immediately | Registration outlives the lease; unmade only by `Pipeline.Shutdown` |
| TTS engines (Piper, Kokoro) | `initTTS`, `loadTTSForRequest`, `ensureTTSVoice` → `admitAudio`, then `ModelLoader.LoadTTS` | Audio pipeline owns the synthesizer object; each synthesis spawns a child | `exec.CommandContext` kills the direct child on cancel; not job-assigned, not tree-killed, `Release()` does not wait |
| MeloTTS | `internal/audio/tts/melotts/client.go` | An HTTP client to a service the application does not start | None: the process is external; admission charges a 1 GB VRAM / 800 MB RAM profile it cannot observe |
| Audio execution and shutdown | `App.beginManagedAudio`, `Manager.Close`, `Pipeline.Shutdown` | Audio manager operation + lease; app shutdown stage `audio manager shutdown` | Close cancels every operation context, retires the generation, drains bounded (5 s), then the loader releases providers |

**Findings (named gaps, in severity order).**

- **F-A1 — TTS children are outside the ownership model that ASR and image already use.** `internal/audio/tts/piper/local.go` and `internal/audio/tts/kokoro/local.go` spawn with `exec.CommandContext` instead of `proc.Start`, and their `Release()` is a no-op. A force-killed app can leave a `piper.exe`/`python.exe` behind, and cancellation terminates only the direct child. The job-owned supervisor that already exists at `internal/audio/tts/helper/supervisor.go` is not on these providers' path. Next evidence: a real child, job assignment before first instruction, kill-tree on cancel, and a test that observes the process disappear.
- **F-A2 — MeloTTS is an unowned process charged to a shared authority.** Admission reserves 1 GB VRAM, 500 MB RAM and the model footprint for a loopback service this application never starts, stops or measures. Until it is either adopted by a supervised owner or explicitly classified as an external service with its own disposition, audio cannot claim the same guarantee as image/video. Note: it is still local-only; `http://127.0.0.1:8888` is loopback, not a cloud service.
- **F-A3 — Audio registration outlives its admission lease without execution revalidation.** `registerAudio` registers and releases the lease immediately (the D-2 class of gap in the contract). Each distinct Piper voice path also adds a loader-tracked provider (`tts:piper:<n>`) that nothing unloads before pipeline shutdown; repeated character voices therefore accumulate registrations and resident claims that the manager never reconciles against measured memory.
- **F-A4 — Kokoro voice enumeration runs an uncancellable child while holding a lease.** `listVoices` uses `exec.Command` with no context and no timeout; it is reached from `GetVoices`, including from `App.GetTTSAvailableVoices`, which holds an audio-manager lease for the duration.
- **F-A5 — Audio shutdown is bounded but not forcing.** `Manager.Close` reports `ErrDrainTimeout` and leaves the manager closing when a generation will not drain. That is a correct refusal, but no audio surface yet proves a stuck provider is killed and reaped on host shutdown.

**Characterization tests added in this slice.**

- `TestTTSStartConsultsAdmissionWhileChatLoaded` (`internal/audio/tts_lifecycle_characterization_test.go`) — TTS asks the shared authority before constructing a provider, a refusal names the resident chat model, the capability stays not-ready, and no runtime model is tracked afterwards. The fixture is one the documented policy genuinely refuses (300 MB Piper RAM request against 180 MB available after the 10% reserve of 200 MB free); it does not touch the open ASR fixture.
- `TestManagedAudioCloseCancelsInFlightSynthesis` — a synthesis blocked in the provider observes cancellation through the manager operation context, and `Close` does not report completion before that provider call returns; leases and operations reach zero in-flight.

Both pass in the default suite and under `-race`. They characterize the execution-scope boundary only. They do **not** prove child-process termination; F-A1 is the reason.

**TTS runtime baseline on this machine (2026-11-25).**

- No Piper voice model is installed: `models/en_US-lessac-medium.onnx` does not exist, `models/` contains no `.onnx` files, and `bin/piper.exe` is absent. `config/settings.json` has `ttsEnabled: false`, `ttsEngine: "piper"`, `ttsVoice: "en_US-lessac-medium"`.
- No Kokoro asset set is installed (`models/kokoro`, `runtime/tts` are absent). No MeloTTS service is listening on the configured loopback URL.
- Therefore TTS has **never synthesized on this machine**. That is missing evidence, not a working capability and not proof of a defect.
- Every supported engine is local: Piper is a local executable, Kokoro is local ONNX behind a local Python bridge, MeloTTS is a loopback service. The only network-capable audio adapters in the build are the DashScope Qwen ASR clients (`qwen.NewClient`, key-gated), which automatic selection excludes because it only accepts `Local && ApprovedLocal` profiles. No audio inference in this application runs in the cloud.

**Host acceptance protocol (requires the user; the agent must not start or take over playback).**

1. Install a local voice once: `piper.exe` at `bin/piper.exe`, and `models/en_US-lessac-medium.onnx` plus its `.onnx.json` beside it (`models` is the project's model store). This is a one-time download; synthesis itself stays offline.
2. Launch the app, open Settings → Audio, enable TTS, select engine `piper` and voice `en_US-lessac-medium`, apply, and confirm the status reports TTS ready with no error.
3. Use the TTS button beside the chat input with a short sentence, or enable message autoplay and send one message. Playback must be started by the user.
4. Record: audible speech or the exact error; the generated WAV path/size; whether `piper.exe` exists during synthesis and is gone afterwards; the newest file under the TTS cache directory; and the audio status fields read back from the backend.
5. Repeat with a chat model loaded and ASR enabled to observe admission and capacity reporting together (cross-workload evidence), then repeat in the browser dev host with an isolated state root and record any native-surface difference.
6. Capture one case per failure mode: missing model, missing binary, and admission refusal (chat resident plus tight free memory), so the refusal text and the untouched process/registration state are visible.

**What this slice does not establish.** Audio is still an open gate: no host run, no real child termination evidence, no cross-workload test, no ASR baseline disposition. F-A1/F-A2 are the next bounded slice; they are integration work, not assertion edits.

#### 7.6.2 Ownership remediation implemented (2026-11-25)

**Scope.** This slice implements F-A1 through F-A5 inside the existing ownership architecture. It adds no second process manager and no audio-only admission rule; it extends `internal/process` with scoped child slots and a job-owned `Child` primitive, and routes the production TTS providers through it.

**Ownership after the change.**

| Surface | Owner | Termination guarantee |
|---|---|---|
| Piper/Kokoro synthesis | `process.Child` under `RoleTTSChild`, registered before the child's first instruction (`proc.StartWithOwnership`), one waiter | Cancellation escalates from stdin close to `KillTree` + `Process.Kill`, then waits for the registered waiter; unconfirmed exit returns `ErrTerminationUnconfirmed` and keeps the slot |
| Kokoro voice enumeration | Same primitive, its own cancellable operation scope (`VoicesContext`), execution lease held for the duration | Cancellation reaches the child and waits for cleanup |
| MeloTTS | `melotts.Service` under `RoleTTSService`, shared per loopback URL with a reference count | `Ensure` requires an owned launcher (`HOWL_MELOTTS_COMMAND` or `bin/melotts-serve.exe`), waits for readiness, and refuses an endpoint it did not start; the last reference terminates and reaps the process |
| Provider residency | `Pipeline.ttsProviders` + `resource.Manager` registration | Residency is the registration, not the initialization return; active executions hold a separate reference; idle providers are reclaimed only after a successful replacement load |
| Enumeration from the app | `App.GetTTSAvailableVoices` → manager operation + generation lease + `GetTTSVoicesContext` | Shutdown cancels it like any other audio operation; failures are logged and the status error count rises |
| Pipeline shutdown | `Pipeline.Shutdown` | Settles owned children first, then residency, then only the registrations this pipeline created; an executing or unconfirmed provider is an error and keeps its ownership |

**Named changes.** Process ownership: `internal/proc.StartWithOwnership`, `internal/process/child.go`, scoped slots in `internal/process/registry.go`. Providers: `piper/local.go`, `kokoro/local.go` (child + cancellable enumeration), `melotts/service.go` and `melotts/client.go` (managed service, reference counting, loopback-only). Pipeline: `pipeline.go` (`ttsProviders`, `globalTTSKey`, residency registration, shutdown order) and `tts_residency.go` (execution leases, reclamation guard, registration ownership, cancellable enumeration). Composition: `manager_factory.go` and `app-boot.go` inject the composition root's one process registry; `app-audio.go` runs enumeration as a managed operation.

**Residency and admission.** An admission lease becomes a resident registration and is released at that moment, so the same bytes are never counted as both a pending reservation and a resident model. Active synthesis takes an execution reference instead of a second reservation; a MeloTTS service this process already owns takes no second reservation. Registrations are removed by the pipeline that created them, never by a global sweep.

**Tests added (default suite).**

| File | Evidence |
|---|---|
| `internal/process/child_test.go` | Real child: registration before work, cancellation kills the tree (helper + grandchild heartbeat stops), one waiter, unconfirmed termination retains ownership, shutdown waits for confirmed exit, concurrent scoped slots all stop. |
| `internal/audio/tts/melotts/service_test.go` | Managed startup/readiness/registration/termination, failed startup is not ready, an unowned loopback endpoint is not adopted, non-loopback endpoints are refused. |
| `internal/audio/tts_ownership_test.go` | Residency survives initialization and shutdown is idempotent; reclamation refuses an executing provider and respects shared references; enumeration is cancellable and holds/releases its lease; shutdown stops a live owned child; shutdown does not claim success with an executing provider and keeps its registration; cross-workload capacity accounting (audio residency appears as a blocker, disappears when released); MeloTTS residency is not double-counted. |
| `internal/audio/pipeline_test.go` | `TestPipelineR8LazyTTSOwnership` graduated to the residency contract: the active engine and current voice stay loaded, the superseded voice is reclaimed exactly once, and shutdown releases the rest. |

**Preserved.** `TestASRStartConsultsAdmissionWhileChatLoaded` still fails at its existing assertion. No policy, fixture, or assertion was changed.

**MeloTTS readiness reporting.** A TCP probe no longer advertises a provider. `melotts.Readiness` answers "is this a ready Howl provider" by requiring (a) the endpoint to be loopback, (b) this process to own a live service at that endpoint, and (c) the owned service to answer; an unowned listener, a stopped owned service, and a remote endpoint are each reported as blocked with the reason. `audio.ProductionCapabilityRecords()` is the single capability contract used by the pipeline, the status payload, and the voice catalog, and it installs that probe; `capability.ProductionEnv` keeps only a reachability fallback for capability-package tests, documented as such. Behavioral coverage: `TestMeloTTSCapabilityRequiresOwnedService` (unowned listener → blocked and not selectable; owned running service → available; released service → blocked again) and `TestMeloTTSCapabilityRefusesNonLoopback`.

**Remaining gaps.** The earlier asset caveat no longer holds: `bin/piper.exe` and `models/en_US-lessac-medium.onnx` are installed and verified, synthesis was exercised through the production provider, and the operator's listening acceptance is recorded (2026-11-25). What remains is integration and lifecycle coverage: a packaged-host click-through of the audio surfaces, Kokoro and MeloTTS host evidence (neither has installed assets here), and the ASR baseline disposition. None of that re-opens the voice-quality question — it is about wiring and lifecycle, not whether Piper can speak.

#### 7.6.3 Packaged-host audio integration evidence (2026-11-25)

**Host and state.** `build/bin/howl-chat.exe` (`-tags "desktop,production"`, WebView2), launched with `HOWL_WEBVIEW_DATA_DIR=C:\HOWL-Test\g8-packaged-20261125` and the `.howl-isolated-root` marker. Every write stayed inside that root (`data\config\settings.json`, `data\cache\tts`, `webview2\`); no real user data was read or written. Window 1400x900 at 580,270 on a 2560x1440 screen.

**Driving method.** UIA exposes only the WebView2 host panes (9 elements, no `Document`), even with `WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS=--force-renderer-accessibility`, so driving is coordinate-based through `build/unit-c/tools/uia.ps1 -Action click -X -Y` with `build/unit-c/tools/shot.ps1` captures. Both scripts gained a `-ProcessName` parameter (default unchanged) to target any packaged build. **Correction for future runs:** `-Action key` (`SendKeys`) does not reach the page — a draft typed that way never appeared. The working sequence is a coordinate click on the composer, `Set-Clipboard` in PowerShell, then `-Action key -Text '^v'`.

**Workflow exercised, step by step** (captures in `build/tmp/g8/`):

| Step | Evidence |
|---|---|
| Settings menu opens in the packaged host | `packaged-02-settings-menu.png` |
| Audio & Voice Settings opens as an in-app modal | `packaged-03-audio-settings.png` |
| TTS tab shows **Piper (Local, fast)** and **en_US-lessac-medium (built-in)** | `packaged-04-tts-tab.png` |
| Save & Close applies; General tab reports **Status: ready** | `packaged-06-ready.png`; isolated `settings.json` records `ttsEnabled: true`, `ttsEngine: "piper"`, `ttsVoice: "en_US-lessac-medium"`, `ttsAutoplayMessage: false` |
| Draft pasted into the composer, speaker clicked | `packaged-08-draft2.png` (focused draft), `packaged-09-speaking.png` (speaker in its playing state) |
| After playback | `packaged-10-after.png`: draft still present, chat list still empty — no message created, no send attempted |

**Process and artefact evidence.** `piper.exe` was polled every 400 ms from the click: 1 child on the first poll, 0 on the following seven and afterwards. Artefact: `<isolation root>\data\cache\tts\tts_1795636237625064600.wav`, 94,368 B, RIFF/WAVE PCM mono 22,050 Hz 16-bit (~2.1 s).

**Shutdown.** The settings apply also initialized local ASR, which ran `llama-server.exe` (PID 10572, ~2.1 GB) for the session. `data\config\logs\startup.log` records `wails.Run returned`, and after the window closed neither `llama-server.exe` nor `piper.exe` was running — the session left no orphan, and no forced termination was needed.

**Human confirmation.** The operator heard this packaged-host composer run, which is the listening half of this step; the browser dev host had already demonstrated the same composer and message paths (`HRA-P1-TTS-LISTENING-ACCEPTANCE-2026-11-25.md`).

**Limits of this evidence.** ASR was demonstrated to readiness (local engine initialized), not to a transcription — that needs a microphone and a human speaker. The per-message speaker was not re-driven in the packaged host because the isolated run held no chat; it is proven in the browser host and shares the same player, voice resolution and cache path. Output device stayed at the Windows default, volume 0.8.

#### 7.6.4 Cross-workload resource and lifecycle validation (2026-11-25)

**Scope.** Extends audio-versus-chat to image and video against the one admission authority and the one process registry. No production code was changed: the added tests passed on first execution, so no ownership or accounting defect surfaced in this workstream. No policy, reserve formula, fixture or existing assertion was modified, and the ASR discrepancy is untouched (see below).

**Workload paths under test.** Image and video share one admission path in production: `ensureEngineReadyLocked -> admitImageLoad` for an image start, and `ensureVideoExecutionLocked -> admitImageLoad` for a video start, whose lease is returned to the caller and held until `sd-cli` exits. Video is therefore modelled as the same image-family request with `Mode: "video"`, not as a separate memory reservation. Audio residency is established through the pipeline's own `registerAudio` path with each provider profile's own measured request.

**Tests added (ordinary default suite, no build tags, file `internal/audio/cross_workload_image_video_test.go`).**

| Test | Evidence | Outcome |
|---|---|---|
| `TestAudioResidencyBlocksImageAdmissionThenReleases` | Image refused with the audio provider in `Blocking` exactly once and a non-zero estimate (projected peak vs available); the refusal names the workload; after `Pipeline.Shutdown` the registration is gone but the *stale* telemetry still refuses the image — a process exit is not treated as recovered capacity; after the measurement refresh the image is admitted; the admitted lease then blocks a second start until it is released, and the release is idempotent | PASS |
| `TestAudioExecutionLeaseSurvivesImageWorkAndReclaimsOnlyWhenIdle` | One registration for a shared provider with two execution leases; the image ledger shows the provider once (no double count); `unloadTTSProvider` refused while executing; releasing one of two leases unloads nothing; the last release allows reclamation and clears the registration; image admitted afterwards | PASS |
| `TestImageResidencyBlocksAudioAdmission` | Reverse direction: a resident image model appears in the audio refusal by name; after the owner unregisters and telemetry reports the freed card, the audio admission succeeds. The ledger entry alone never invents free VRAM | PASS |
| `TestVideoRenderReservationObservesAudioResidency` | Video request carries its render working set; refused beside audio residency naming it; refused again on stale telemetry after release; admitted after refresh; the held render reservation blocks a concurrent image start; released once (idempotent) and the image then admits | PASS |
| `TestVideoProcessStopLeavesAudioOwnerAndRegistrationIntact` | Real children in one registry: owned TTS child (`process.StartChild`, `RoleTTSChild`) and an `sd-cli` stand-in (`RoleVideoCLI`). Stopping the video slot leaves the audio child alive, keeps its slot and its TTS registration and provider; `Pipeline.Shutdown` then stops its own child, empties its slot, removes exactly its registration and does not touch the video slot | PASS |
| `TestAudioCancellationDoesNotReleaseImageExecutionLease` | Cancelling an audio enumeration reaches that operation and returns only after its cleanup, leaks no execution lease, and leaves the image workload's lease still counted; the image capacity returns only when the image owner releases it | PASS |

**Focused and full-suite results.**

| Command | Result |
|---|---|
| `go test ./internal/audio/ ./internal/resource/ ./internal/image/engine/sdcpp/ -count=1` (baseline, before edits) | `internal/resource` ok, `internal/image/engine/sdcpp` ok, `internal/audio` FAIL on `TestASRStartConsultsAdmissionWhileChatLoaded` only |
| `go test ./internal/audio/ -run '<the six tests above>' -count=1 -v` | all six PASS (0.09s, 0.00s, 0.00s, 0.00s, 0.56s, 0.00s) |
| `go test ./internal/audio/ -race -run 'CrossWorkload\|AudioResidency\|AudioExecutionLease\|ImageResidency\|VideoRender\|VideoProcessStop\|AudioCancellation' -count=1` | ok 3.620s, no race reports |
| `go test -count=1 ./...` | every package ok except: `howl-chat` FAIL (`TestBuildPiperVoiceCatalogMarksInstalledAndMissingVoices`, `TestMicrophoneServiceUsesWAVDataBindingAndComposer`), `howl-chat/internal/audio` FAIL (`TestASRStartConsultsAdmissionWhileChatLoaded`) |
| `go vet ./internal/audio/` | exit 0 |
| `go build -tags "desktop,production" -o build/bin/howl-chat-s1.exe .` | exit 0, 31,789,568 bytes |
| `go build -tags uiharness -o build/bin/howl-ui-harness.exe .` | exit 0, 21,977,088 bytes |

**The three suite failures are not from this workstream.** `TestASRStartConsultsAdmissionWhileChatLoaded` is the preserved ASR admission discrepancy (unchanged, still red, not skipped). `TestBuildPiperVoiceCatalogMarksInstalledAndMissingVoices` failed because the Gate 8 asset install put a real voice in `models/`, which changed the catalog default its fixture expected. `TestMicrophoneServiceUsesWAVDataBindingAndComposer` was a source-text guard rejecting `createObjectURL` in `audio-service.js`, which the accepted blob-playback rewrite now uses; it was a stale assertion from that workstream, not a lifecycle defect. All three are reported here. The two `howl-chat` failures were subsequently repaired as test-only isolation fixes, and the ASR discrepancy was resolved as a documented policy/test decision (plan §5 and §7.6); the default suite is green.

**Host run (browser dev host, real backend, isolated root `C:\HOWL-Test\cross-workload-20261125`, `.howl-isolated-root` = `howl-isolated-test`).**

| Step | Observed |
|---|---|
| Host start | `HOWL_WEBVIEW_DATA_DIR` set to the isolated root; `howl-ui-harness.exe` served `http://127.0.0.1:8799` and answered 200 on the first poll; startup log shows `[ImageEngine] Initialized: binary=C:\Users\StoneWolfPC\HOWL_Chat\sd-server.exe` and `[lifecycle] ready: 6 stages in 139ms` |
| Audio activation through the production binding | `ApplyAudioSettings(<settingsState audio>)` returned nil; `GetAudioStatus().tts_ready = true` |
| Real synthesis | `SynthesizeSpeech(...)` returned `...\data\cache\tts\tts_1795637943862442700.wav`, 552,328 bytes, RIFF/WAVE PCM mono 22,050 Hz 16-bit |
| Audio residency in the shared ledger | `GetLoadedModels()` returned `{type:"asr", name:"qwen3-asr-1.7b-local", gpu_memory_mb:4096, system_ram_mb:3840, priority:4}` and `{type:"tts", name:"piper", gpu_memory_mb:0, system_ram_mb:360, priority:5}` |
| Capacity the image workload would consult | `GetResourceSnapshot()`: free VRAM 4551 MB / total 8192 MB, used 3641 MB, free RAM 51333 MB, GPU "NVIDIA GeForce GTX 1080", `memory_tier` large |
| Shutdown | `POST /__shutdown` returned 200; harness exited within two 0.5s polls; `howl-ui-harness.exe`, `piper.exe`, `llama-server.exe`, `sd-server.exe`, `sd-cli.exe` all = 0 afterwards. The ASR child and the TTS child were reaped by the orderly shutdown |

**Host image and video workloads were recorded as blocked — that conclusion was wrong.** `sd-server.exe` (1,675,264 bytes) and `sd-cli.exe` (584,192 bytes) are present, but `models/` contains only the Piper voice, `imageGen.modelsDir` was empty, and no `.safetensors`/`.gguf`/`.ckpt` exists in the project (depth 3) or in `C:\models`, `D:\models`, `E:\models`. The scan never included the owner's model drives. `H:\AI GGUF\Image models` holds the real asset set (33 models, 4 VAEs, 8 LoRAs, plus `hunyuan-video-t2v-720p-Q3_K_S.gguf`), and both supported hosts generated a 1024x1024 SDXL image from it on 2026-11-25 (§7.6.8). The earlier block is retracted; video render is still unattempted.

**Screenshots.** The preview capture returned pixels but refused to save ("Use a new workspace-relative .png path"), so no screenshot file exists for this run; the host evidence above is the machine-readable binding output and the process table. No screenshots are claimed.

**Not reopened.** Voice quality and Piper listening acceptance are unchanged (§7.6.3 and the listening record). This subsection adds resource and lifecycle coverage only; Gate 8 stays open.

#### 7.6.5 Host ASR transcription and provider coverage (2026-11-25)

**Scope.** This slice exercises real ASR transcription in a supported host through the production bindings, under the same isolated-root discipline as §7.6.3/§7.6.4. It changes no production code. Kokoro and MeloTTS are reported from their host capability records and filesystem state; neither was installed or started.

**Environment.** Browser dev host (`-tags uiharness`), `HOWL_WEBVIEW_DATA_DIR=C:\HOWL-Test\asr-host-20261125` with `.howl-isolated-root`, harness on `http://127.0.0.1:8799`. Real backend, real local llama-server ASR, real Piper. The host opened no capture stream and started no playback.

| Step | Observed |
|---|---|
| Mic gate with the headset off | `GetAudioStatus().mic` = `{state:"no-device", device_count:0, detail:"GetDefaultAudioEndpoint(eCapture, eConsole) returned E_NOTFOUND: no default console capture endpoint exists"}`; `ApplyAudioSettings(asrEnabled:true, …)` **FAILED in 12 ms** with "Windows shows a mic isn't selected…"; `asr_ready:false`, `tts_ready:false`. A `TranscribeAudioData` attempt while ASR was off returned **"asr is not enabled"** in 5 ms, so the refusal is surfaced rather than silently satisfied. |
| Headset reconnected | The same status poll reported `mic.state:"ready"`, `device_count:1`, "capture was not opened (no stream, no recording)"; the retried `ApplyAudioSettings` returned nil in 2.35 s. |
| ASR and TTS ready | `asr_ready:true`, `asr_readiness:"local-initialized"`, `tts_ready:true`; `GetLoadedModels()` = `{type:"asr", name:"qwen3-asr-1.7b-local", gpu_memory_mb:4096, system_ram_mb:3840, priority:4}`. |
| Real synthesis used as the input | `SynthesizeSpeech("The quick brown fox jumps over the lazy dog. Howl audio integration gate, November twenty fifth.")` returned `data\cache\tts\tts_1795646079481228700.wav`, 256,788 bytes, PCM mono 22,050 Hz 16-bit, 128,372 frames (≈5.8 s), in 715 ms. |
| Real transcription through the binding | `TranscribeAudioData(data:audio/wav;base64,…)` returned **"The quick brown fox jumps over the lazy dog. Hello, audio integration gate. November twenty-fifth."** in 632 ms. |
| Process and capacity state | A single `llama-server.exe` (PID 17592) was present before and after the transcription; `GetResourceSnapshot()` reported free VRAM 4604 / total 8192, used 3588 MB, free RAM 51,144 MB, GTX 1080; `active_sessions:0`, `error_count:0`. |
| Shutdown | `GET /__shutdown` → `{"ok":true}`; `howl-ui-harness.exe`, `llama-server.exe`, `piper.exe`, `sd-server.exe`, `sd-cli.exe` all = 0 afterwards. |

**What the transcript proves.** The WAV the application wrote was base64-encoded and sent back through `TranscribeAudioData`; the local Qwen ASR child accepted the 22,050 Hz WAV and returned the spoken sentence. The only differences from the spoken text are "Howl" → "Hello" and "twenty fifth" → "twenty-fifth", both plausible confusions of the synthesized utterance. This is a genuine speech-to-text path over application-produced audio, not a canned fixture.

**Host scope and limitations.** This is the browser dev host. Settings were applied through the same binding the settings panel calls, not by clicking the panel; no microphone capture stream was opened (the input was a WAV the app itself produced) and no playback was started. The packaged host has TTS evidence (§7.6.3) but has not run ASR transcription. The mic-gate rows are recorded because they show ASR cannot claim readiness without a capture endpoint, and that the refusal is surfaced.

**Kokoro and MeloTTS are blocked, not failed.** The host capability records report `kokoro` with `model_state: unavailable — "no voice model path was resolved"` and `melotts` with `model_state: unavailable — "MeloTTS is not running as a Howl-managed local service at http://127.0.0.1:8888"`. On disk, `models/kokoro` and `runtime/tts` are absent, only `scripts/kokoro-tts.py` exists, `python -c "import melo"` fails with `ModuleNotFoundError: No module named 'melo'`, and nothing listens on the MeloTTS loopback URL. No asset was installed and no service was started, so their ownership, residency, cancellation and termination behavior remains **unverified**; Gate 8 stays open on that workstream.

#### 7.6.6 Packaged-host ASR attempt (2026-11-25)

**Scope.** Drives the packaged WebView2 host (`howl-chat-s1.exe`, `-tags "desktop,production"`) through the real settings UI on an isolated root, to extend audio host coverage from TTS (§7.6.3) to ASR. No production code was changed. Screenshots are in `build/tmp/g8-asr/`.

**Environment.** `HOWL_WEBVIEW_DATA_DIR=C:\HOWL-Test\asr-packaged-20261125` with `.howl-isolated-root`; the window was maximized (1400×900 default, then 2576×1456) and driven by `build/unit-c/tools/uia.ps1` coordinate clicks with `shot.ps1` captures, because the accessibility tree exposes only the WebView2 host panes (§7.6.3, 7 elements, no `Document`).

| Step | Observed |
|---|---|
| Settings → Audio & Voice Settings | `Settings` menu and the in-app modal opened (`packaged-04-settings-menu.png`, `packaged-07-audio-general.png`) |
| ASR enabled and engine chosen in the UI | General tab ASR toggle on; Speech Recognition tab engine set to **Qwen3 ASR 1.7B (local llama.cpp)** via the real `<select>`; Save & Close applied (`packaged-16-engine-picked.png`, `packaged-17-applied.png`) |
| Persisted and resident | `data\config\settings.json` records `asrEnabled: true`, `asrEngine: "qwen3-asr-1.7b-local"`, `pushToTalk: false`; `llama-server.exe` appeared (PID 11900, ~2.18 GB) as the packaged host's ASR child |
| Readiness in the packaged host | General tab reports **"Status: ready (local engine initialized; transcription not verified)"** (`packaged-23-status.png`) |
| First recording attempt | The mic button started recording (`Recording started…` toast, `packaged-34-start.png`) and stopped after 30 s; the toast read **"No speech detected"** (`packaged-35-stop2s.png`) and no ASR request reached the child. The device was not muted — the operator was away from the desk and did not speak, so the capture buffer was genuinely silent. |
| Second attempt (operator present) | Recording started (`packaged-42-record-start.png`); the operator spoke during the window; the composer received **"Testing. Okay, the quick brown fox jumps over the lazy dog. Package host speech recognition test."** (`packaged-44-final.png`, `crop-c-packaged-44-final.png`) |
| Process, GPU and temp state | `llama-server.exe` PID 12216 resident (~2.22 GB working set); `nvidia-smi` reported 4121 MiB / 8192 MiB while ASR was ready; `data\temp` held no leftover WAV after the request |
| Shutdown | Window closed normally; `data\config\logs\startup.log` records `wails.Run returned`; `howl-chat-s1.exe`, `llama-server.exe`, `piper.exe`, `sd-server.exe`, `sd-cli.exe` all = 0 afterwards |

**Packaged-host ASR transcription is demonstrated.** The real workflow ran: Settings UI enablement → persisted engine → owned `llama-server` child → real microphone capture → WAV encode → `TranscribeAudioData` → transcript inserted into the composer. The first attempt's silence was the operator being away, not a hardware or pipeline fault; the retry with a live speaker transcribed correctly. Together with §7.6.3's Piper run, ASR and TTS are now exercised in the packaged host as well as the browser dev host (§7.6.5).

**Lifecycle gap observed here (fixed 2026-11-25 — §7.6.7): re-applying settings while ASR is resident was refused as self-blocking.** Saving the Settings panel a second time while the same engine was already resident produced `ASR-MODEL-LOAD` at high severity (`data\config\logs\howl-errors-20261125.jsonl`):

> `apply ASR engine "qwen3-asr-1.7b-local": ASR: admit ASR qwen3-asr-1.7b-local: cannot start asr "qwen3-asr-1.7b-local": needs 4096 MB VRAM; 4184 MB free on device 0 (NVIDIA GeForce GTX 1080), 418 MB reserve (10%), 3766 MB available; the 290 MB soft reserve allowance is not enough (it would leave less than 128 MB free); blocking model: asr "qwen3-asr-1.7b-local" holds 4096 MB VRAM and 3840 MB RAM`

The blocking model is the audio pipeline's own registration, so the apply path re-admitted an engine it already owned instead of treating it as resident or replacing it. Readiness stayed `local-initialized` (the earlier engine kept running), but the renderer took the failed apply as `asrEnabled` off and the mic path refused until ASR was disabled and re-enabled cleanly. This was a genuine idempotency/lifecycle gap, distinct from the admission policy and from the ASR test discrepancy. **Resolved in §7.6.7; the records above are preserved as observed.**

### 7.6.7 Re-apply self-block — fixed and re-verified (2026-11-25)

**Scope and authority.** Owner-directed production fix, limited to the settings apply path. It implements the idempotence rule the design already specified (D5 in `docs/asr-repair-design.md`: "Apply is idempotent: if the normalized desired set equals the active set and the active pipeline is healthy, return current status without rebuilding"). No admission rule, reserve, provider behavior, test assertion or acceptance requirement was changed, and no T2.9.6 scope is claimed.

**Cause.** `Manager.Apply` advanced `DesiredGen` on every request, and `canSkipLocked` skipped only a generation that a caller ahead on `applyMu` had already installed (`appliedGen >= gen`). An identical save therefore built a new candidate, and that candidate's ASR admission request met the active generation's own registration — the same engine holding 4096 MB VRAM and 3840 MB RAM — so the resource authority refused a model the process already owned. The refusal was truthful; the rebuild was unnecessary. The renderer then treated ASR as off until it was disabled and re-enabled.

**Fix (`internal/audio/manager.go`).** `canSkipLocked` now returns true when the active generation's settings equal the desired settings and every requested capability is ready (`capabilitiesReadyLocked`). The skip path publishes the desired generation as applied so `DesiredGen`/`AppliedGen` stay converged for observers. Boundaries kept deliberately:

- **Genuine changes still rebuild.** A different engine, language, threads, beam/best-of, TTS voice/speed, or toggles all differ in `Settings`, so a new generation is built and admission runs unchanged. No release-before-admit was introduced, so a refused replacement still preserves the old generation.
- **A degraded install is not reused.** If a requested capability is not ready, an identical request may rebuild, so a transient failure cannot pin the capability until the settings change.
- **Disable still retires.** A request with both capabilities off never matches a running install and takes the existing retirement path.

**Tests (default suite, no skips).** `internal/audio/manager_idempotence_test.go`: `TestReapplyInstalledSettingsDoesNotRebuild` (one factory call, status ready, generations converged), `TestReapplyDegradedActivationRetries`, `TestChangedSettingsStillRebuild`, `TestIdempotenceDoesNotBlockRetirement`. Four existing lifecycle tests forced replacement by re-applying identical settings; under the corrected rule that request is a no-op, so their fixtures now request a genuine change (`settingsBothReplaced`) with every assertion unchanged: `TestManagerM2StrictReplacementKeepsOldGeneration`, `TestManagerM5ZeroReferenceRetirementViaApply`, `TestManagerM13PanicBoundary/retired_shutdown_panic`, `TestPipelineManagerR2StrictReplacement`. `TestA3TemporaryEnableBridgePreservesPreviousReadiness` changes the desired engine before its replacement enable for the same reason, keeping its bridge and old-readiness assertions.

**Checks run.** `go test -count=1 ./internal/audio/` pass; `-race` on the affected manager and pipeline tests pass; `go vet ./internal/audio/` clean; `go build ./...` exit 0; full `go test -count=1 -p 1 ./...` exit 0 with every package ok; packaged build `-tags "desktop,production"` rebuilt (31,789,568 bytes).

**Packaged-host verification (rebuilt binary, isolated root `C:\HOWL-Test\asr-packaged-20261125`).** ASR was not resident at launch (no `llama-server`). The first **Save & Close** loaded it through the real settings UI: `llama-server.exe` PID 14664, 2,130 MB working set, `nvidia-smi` 3,854 MiB / 8,192 MiB, and the reopened panel reported **"Status: ready (local engine installed, transcription not verified)"** (`reapply-04-audio-settings.png`, `reapply-06-panel-reopen.png`). The second **Save & Close** with unchanged values — the exact scenario that previously failed — returned silently with the panel closing normally: same PID 14664, same 2,130 MB, 3,853 MiB GPU, and no new `ASR-MODEL-LOAD` record (the error log still contains only the two records from 22:47:05 and 22:49:45 quoted above). `reapply-07-second-apply.png` is the post-save state. Closing the window recorded `wails.Run returned` in `startup.log`, and `howl-chat-s1`, `llama-server`, `piper`, `sd-server` and `sd-cli` were all 0 afterwards — the ownership behavior Gate 8 asks for.

**Not verified.** Transcription after the idempotent re-apply was not re-run: it needs a live speaker, and no audio was played through the speakers to simulate one. The resident instance after the re-apply is evidenced by the ready status, the unchanged PID, working set and GPU figure, and the clean shutdown — not by a new transcript.

### 7.6.8 Image generation with real weights in both hosts (2026-11-25)

**Correction to the earlier block.** §7.6.4 and §5 recorded that no image or video weights existed "anywhere" reachable. That was a lookup failure: the scan covered `models/`, the configured `imageGen.modelsDir`, and `C:`/`D:`/`E:` model paths, but never the owner's model drives. `H:\AI GGUF\Image models` holds the real asset set — SDXL checkpoints (`Illustrious-XL-v2.0`, `Juggernaut-XI`, `pony-v7-base`, `bbArtisticFusionXL_v1`, …), Flux/Qwen/HiDream/Chroma GGUFs, `t5xxl_fp8_e4m3fn`, `clip_l`, `ae`, `sdxl.vae`, 8 LoRAs, 4 VAEs, and `hunyuan-video-t2v-720p-Q3_K_S.gguf`, each with a `.howl.json` sidecar. The app reports **33 models · 4 VAE · 8 LoRA** from that folder. `G:\AI GGUF` also exists and holds the ASR/voice set used in §7.6.6.

**Isolated root and fixture.** `C:\HOWL-Test\hra-image-20261125`, marker content `howl-isolated-test`, scan dirs seeded to `H:\AI GGUF\Image models`. The first launch failed because that seed contained single-backslash JSON (`invalid escape sequence \A`); the loader refused and left the model registry nil, which the UI surfaced as "model registry not initialised" rather than substituting an empty list. The fixture was rewritten with valid JSON and the app relaunched; no production code was involved.

**First real attempt failed honestly.** `Pony-Realism-V23ULTRA-Q4_K_M.gguf` was refused by the engine: `sd-server txt2img error 500: line 45, column 18: cannot load component from 'H:\AI GGUF\Image models\Pony-Realism-V23ULTRA-Q4_K_M.gguf'`. The UI showed the error and marked the engine unavailable until restart, which is the intended refusal behavior; this file needs companion components it does not have. It is not a pass and is not counted as one.

**Successful workflow — `Illustrious-XL-v2.0` (SDXL, 6.94 GB safetensors).** The model preset set 1024x1024, 20 steps, Euler a, CFG 7, Karras (the preset overrode the 512 values entered before load — recorded as observed behavior).

| Host | Process / resource | Result |
|---|---|---|
| Browser dev host (harness `127.0.0.1:8799`, same isolated root) | `sd-server` PID 16936; GPU peaked **7716 MiB / 8192 MiB**; ~140 s | `data/images/gallery/img_20261125_233304_307700.png`, 1024x1024 RGB, 1,753,400 B, with JSON metadata (`model Illustrious-XL-v2.0`, `steps 20`, `seed 2126523921`, `sampler euler_a`) |
| Packaged WebView2 (`howl-chat-s1.exe`, maximize 2576x1456) | `sd-server` PID 7056; GPU 6.2–6.6 GB observed during decode | `img_20261125_233627_290100.png`, 1024x1024 RGB, 1,438,547 B; canvas shows the image with `Model: Illustrious-XL-v2.0 (sdxl)`, `Size 1024x1024`, `Steps 20`, `Seed 2126523921`, `Sampler euler_a`, `CFG 7`, `Mode txt2img`; gallery holds both images |

Both outputs were opened and visually inspected: a lighthouse on a rocky coast at dusk with warm lamp light, matching the prompt. Closing the packaged window recorded `wails.Run returned`; `howl-chat-s1`, `sd-server`, `sd-cli`, `llama-server` and `piper` were all 0 afterwards. Screenshots: `image-09-models.png` (61 assets discovered), `image-13-generated.png`, `image-20-pkg-generated.png`, `image-21-pkg-result.png`.

**Driving limitation recorded.** In the packaged host the native model `<select>` could not be opened by coordinate clicks, so the persisted model/prompt from the isolated root were used and generation was started with the page's documented **Ctrl+Enter** shortcut after focusing the prompt. The generation, the engine work and the visible result are real; the packaged *selection* step was not clicked. In the browser host the model value was set through the DOM change event and the Generate button was clicked with browser input.

**Not verified.** Video rendering: `sd-cli` was not driven. The Hunyuan weight's sidecar declares `family: flux`, so its classification needs checking before a video render can be attempted. Kokoro/MeloTTS unchanged (§7.6.5), and the Gate 7 §15.1 both-host workflow still has no record.

### 7.6.9 Gate 7 §15.1 in-session attempt — browser host (2026-11-25)

**Scope.** Contract §15.1 asks for one session per host demonstrating: chat resident, image generation, video admission/execution, image completion, and clean shutdown with no phantom registration or orphan process. This records the browser-host attempt and stops short of a pass.

**Setup.** Isolated root `C:\HOWL-Test\hra-image-20261125`; scan dirs `H:\AI GGUF\Image models` (image) and `H:\AI GGUF\GGUF` (llm). The app registered **27 chat models** and the **61 image assets** already used in §7.6.8.

**Chat resident — demonstrated.** `Qwen3.5-9B-Uncensored-HauhauCS-Aggressive-Q4_K_M` loaded through the real Model Browser: `llama-server` PID 17676, 5,517 MB working set, placement reported `0/33 layers vram=0MB ram=5356MB context=3328 matches_plan=true`, UI status **Ready**. A typed turn received a streamed reply containing "resident chat confirmed." (six chunks, 44 bytes; prompt 96 tokens at 39.8 tok/s, 10 predicted at 3.6 tok/s).

**Correction — that residency was CPU, not VRAM.** The spawned argv carried `-ngl 0 --fit off`, and admission logged `admitted: 0 MB VRAM required, 7224 MB available after 0 MB reserve`. The model was fully RAM-resident (5,356 MB) and used none of the free 7,224 MB, because the model-settings default is `gpu_layers: 0` (`frontend/dist/js/model-settings.js:17`, slider default 0) and chat admission treats that as a zero-VRAM request. §15.1's "chat resident" leg therefore did not demonstrate GPU-resident chat; see §7.6.10.

**Defect observed — conversation save refused mid-stream.** Immediately after streaming, `SaveChat` returned `FAILED`: `chat save refused: id=2026-11-25T23-46-25Z: refused_reentrant: the same operation entered the mutation boundary twice`. The refusal is surfaced in the backend log; the reply was not lost, because the shutdown path's `[FinalSave] committed history with 2 entries` wrote the conversation. This is a real mutation-boundary re-entrancy refusal during a normal send, not a transport error, and it needs a separate disposition.

**Image leg not completed in-session.** Opening `image-generator.html` in the same harness while chat was resident left the page's main thread unresponsive: `preview_evaluate` timed out repeatedly (even `1+1`), and the page issued no RPC after 17:46:54. Console showed `404` for `/wails/ipc.js` and `/wails/runtime.js` (refused for `text/plain` MIME) and `ReferenceError: rescanImageModelDirs is not defined`. The same page had loaded 33 models earlier in this harness when no chat model was resident, and it loaded normally in the packaged host, so this is recorded as a **harness page-load failure with an unresolved cause**, not as a proven application defect. No image was generated in this session; §7.6.8's two runs (which had no chat resident) remain the image evidence.

**Video leg not attempted.** See §7.6.11: the only video weight needs components that are not present.

**Clean shutdown — demonstrated.** `/__shutdown` returned, `[FinalSave]` committed, and `howl-ui-harness`, `llama-server`, `sd-server`, `sd-cli` and `howl-chat-s1` were all zero afterwards. No phantom registration or orphan process was left.

**Verdict.** §15.1 is **not complete** in the browser host: chat and shutdown pass, image and video do not, and one save defect was newly exposed. The packaged host's §15.1 run has not started.

### 7.6.10 Workload-switching audit — intended vs implemented (2026-11-25)

**Intended (owner's statement).** One primary GPU workload at a time when VRAM requires it: chat loads with automatic layer placement; starting an image generation unloads chat *entirely* from VRAM and loads the image model; returning to chat unloads the image model and reloads chat; audio prefers CPU/RAM and only uses spare VRAM, never at the primary workload's expense.

**Finding 1 — chat layer placement is not automatic at all.** `gpu_layers` defaults to `0` in `frontend/dist/js/model-settings.js:17` and the slider's default is `0` (`frontend/dist/index.html:705`). `chatAdmissionRequest` copies that value into the request (`app-model.go:654`), so admission sees a zero-VRAM requirement and grants it (`app-model.go:232` in the §15.1 log). `buildLlamaServerArgs` passes it straight to `-ngl` (`app-llama.go:111`), and `--fit off` disables llama.cpp's own placement. The only fitting that exists, `ClampToHardware` (`internal/backend/gguf/memory.go:118`), *reduces* the requested layer count and can never raise it. The 0/33 observed in §7.6.9 is therefore the default, not an image-residency effect.

**Finding 2 — with another model resident, chat is degraded to CPU instead of unloading it.** In `Admit` (`internal/resource/admission.go:479`), the fitted-plan branch requires `len(blocking) == 0`. When the image model is loaded, blocking is non-empty, so `LoadModel` (which always sets `AllowManual = true`, `app-model.go:285`) reaches `manualChat`, which calls `fitChatRequest(req, 0, availableRAM)` — a hard-coded **0 MB VRAM budget** (`internal/resource/admission.go:597`). That is precisely "push it to CPU" rather than "unload the other workload".

**Finding 3 — real switching exists, but only behind a confirmation dialog.** The image path runs `CheckImageFit` → `confirm(fit.message)` → `FreeVRAMForImage` (`frontend/dist/js/image-generator-core.js:690-712`); the chat path mirrors it with `FreeVRAMFor(kind, path, approved)` (`frontend/dist/js/chat-model.js:348`). The backend (`app-vram-arbitration.go`) stops the other workload outright (`UnloadModel` / `UnloadImageModel`, i.e. process teardown), then polls live telemetry until the same request is admitted (`waitForAdmissionFit`, 10 s). This is the mechanism the owner described — full unload, not layer demotion — but it is user-confirmed, not automatic.

**Finding 4 — a second, contradictory eviction policy exists and is advisory only.** `internal/resource/allocator.go:15` declares `PriorityLLM = 0 // Chat model — never evict automatically` and `PriorityImage = 8 // Image gen — evict first`, the opposite of the intended priority under contention. Its `Allocate` is reachable only from the advisory `RequestAllocation` binding (`app-resource.go:138`), never from a load. So the tree contains two conflicting policies, and the one that would auto-evict is unused.

**Finding 5 — audio is profile-driven, not "spare VRAM only".** ASR/TTS admission sets `RequireGPU` whenever the profile declares `VRAMBytes > 0` (`internal/audio/pipeline.go:808`, `:829`); only profiles that declare zero VRAM are CPU-only (Whisper Tiny). There is no "use VRAM if there is spare capacity, otherwise RAM" decision per load, and no rule preventing audio from taking space a primary workload would need beyond ordinary admission.

**Verdict.** The intended behavior is about half present: the unload/reload handoff is implemented and resource-truthful, but it is dialog-gated; automatic layer placement does not exist; and under contention the default path silently demotes the incoming chat model to CPU instead of switching workloads. This is a design gap in the product layer, not a defect in the lifecycle machinery audited in §7.6.

### 7.6.11 Hunyuan video weight — metadata inspection (2026-11-25)

**Actual metadata, read from the GGUF header.** `hunyuan-video-t2v-720p-Q3_K_S.gguf`: GGUF v3, 856 tensors, 3 key/value pairs, `general.architecture = hyvid`. It is a HunyuanVideo diffusion model, not a Flux model.

**The sidecar is wrong.** `hunyuan-video-t2v-720p-Q3_K_S.howl.json` declares `"family": "flux"`, `"type": "diffusion"`, with only a `model` entry and Flux-style recommendations. Loading it on that sidecar would present a Hunyuan weight to the engine as Flux, which is exactly the filename/classifier assumption this check exists to prevent. The scanner's filename classifier does map `hunyuan` to `hunyuanvideo` (`internal/image/scanner/classifier.go`), so the registry entry and the sidecar disagree; that conflict is recorded, not silently resolved.

**Engine expectations.** `internal/image/engine/sdcpp/video.go` treats Hunyuan as a split video architecture: `-M vid_gen` with `--diffusion-model` plus optional `--vae`, `--t5xxl`, `--clip_l`, `--clip_g`, `--clip_vision`, `--llm`, native size 848x480, with `--video-frames`/`--fps` and temporal VAE tiling. `sd-cli --help` on the shipped build (commit 3b12781) exposes those video flags, but `grep -a -i 'hyvid|hunyuan' sd-cli.exe` returns **no symbols**, so engine support for the `hyvid` architecture is not established by inspection.

**Companions present?** No. The folder has Flux/ SDXL VAEs (`ae.safetensors`, `sdxl.vae`, `qwen_image_vae`), `t5xxl_fp8_e4m3fn`, `clip_l`, and Qwen VL encoders — none of which is the HunyuanVideo VAE or its llava text encoder. A Hunyuan render would have to run without its required VAE/encoder pair.

**Decision.** No render was attempted and nothing was forced. Classification: **incompatible as configured** — wrong sidecar family, no Hunyuan companions present, and engine support for `hyvid` unproven. A future attempt needs: a corrected `.howl.json` (or a registry entry) declaring `hunyuanvideo`, the matching HunyuanVideo VAE and text encoder, and confirmation that the shipped `sd-cli` can load `hyvid`.

